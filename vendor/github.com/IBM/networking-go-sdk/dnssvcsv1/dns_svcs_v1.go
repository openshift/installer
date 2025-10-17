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
 * IBM OpenAPI SDK Code Generator Version: 3.94.1-71478489-20240820-161623
 */

// Package dnssvcsv1 : Operations and models for the DnsSvcsV1 service
package dnssvcsv1

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
	common "github.com/IBM/networking-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// DnsSvcsV1 : DNS Services API
//
// API Version: 1.0.0
type DnsSvcsV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.dns-svcs.cloud.ibm.com/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "dns_svcs"

// DnsSvcsV1Options : Service options
type DnsSvcsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewDnsSvcsV1UsingExternalConfig : constructs an instance of DnsSvcsV1 with passed in options and external configuration.
func NewDnsSvcsV1UsingExternalConfig(options *DnsSvcsV1Options) (dnsSvcs *DnsSvcsV1, err error) {
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

	dnsSvcs, err = NewDnsSvcsV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = dnsSvcs.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = dnsSvcs.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewDnsSvcsV1 : constructs an instance of DnsSvcsV1 with passed in options.
func NewDnsSvcsV1(options *DnsSvcsV1Options) (service *DnsSvcsV1, err error) {
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

	service = &DnsSvcsV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "dnsSvcs" suitable for processing requests.
func (dnsSvcs *DnsSvcsV1) Clone() *DnsSvcsV1 {
	if core.IsNil(dnsSvcs) {
		return nil
	}
	clone := *dnsSvcs
	clone.Service = dnsSvcs.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (dnsSvcs *DnsSvcsV1) SetServiceURL(url string) error {
	err := dnsSvcs.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (dnsSvcs *DnsSvcsV1) GetServiceURL() string {
	return dnsSvcs.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (dnsSvcs *DnsSvcsV1) SetDefaultHeaders(headers http.Header) {
	dnsSvcs.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (dnsSvcs *DnsSvcsV1) SetEnableGzipCompression(enableGzip bool) {
	dnsSvcs.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (dnsSvcs *DnsSvcsV1) GetEnableGzipCompression() bool {
	return dnsSvcs.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (dnsSvcs *DnsSvcsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	dnsSvcs.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (dnsSvcs *DnsSvcsV1) DisableRetries() {
	dnsSvcs.Service.DisableRetries()
}

// ListDnszones : List DNS zones
// List the DNS zones for a given service instance.
func (dnsSvcs *DnsSvcsV1) ListDnszones(listDnszonesOptions *ListDnszonesOptions) (result *ListDnszones, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListDnszonesWithContext(context.Background(), listDnszonesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListDnszonesWithContext is an alternate form of the ListDnszones method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListDnszonesWithContext(ctx context.Context, listDnszonesOptions *ListDnszonesOptions) (result *ListDnszones, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listDnszonesOptions, "listDnszonesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listDnszonesOptions, "listDnszonesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listDnszonesOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listDnszonesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListDnszones")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listDnszonesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listDnszonesOptions.XCorrelationID))
	}

	if listDnszonesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listDnszonesOptions.Offset))
	}
	if listDnszonesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listDnszonesOptions.Limit))
	}
	if listDnszonesOptions.VpcID != nil {
		builder.AddQuery("vpc_id", fmt.Sprint(*listDnszonesOptions.VpcID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_dnszones", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListDnszones)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateDnszone : Create DNS zone
// Create a DNS zone for a given service instance.
func (dnsSvcs *DnsSvcsV1) CreateDnszone(createDnszoneOptions *CreateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateDnszoneWithContext(context.Background(), createDnszoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateDnszoneWithContext is an alternate form of the CreateDnszone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateDnszoneWithContext(ctx context.Context, createDnszoneOptions *CreateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDnszoneOptions, "createDnszoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createDnszoneOptions, "createDnszoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createDnszoneOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createDnszoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateDnszone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createDnszoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createDnszoneOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createDnszoneOptions.Name != nil {
		body["name"] = createDnszoneOptions.Name
	}
	if createDnszoneOptions.Description != nil {
		body["description"] = createDnszoneOptions.Description
	}
	if createDnszoneOptions.Label != nil {
		body["label"] = createDnszoneOptions.Label
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_dnszone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnszone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteDnszone : Delete DNS zone
// Delete a DNS zone.
func (dnsSvcs *DnsSvcsV1) DeleteDnszone(deleteDnszoneOptions *DeleteDnszoneOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteDnszoneWithContext(context.Background(), deleteDnszoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteDnszoneWithContext is an alternate form of the DeleteDnszone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteDnszoneWithContext(ctx context.Context, deleteDnszoneOptions *DeleteDnszoneOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDnszoneOptions, "deleteDnszoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteDnszoneOptions, "deleteDnszoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteDnszoneOptions.InstanceID,
		"dnszone_id":  *deleteDnszoneOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteDnszoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteDnszone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteDnszoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteDnszoneOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_dnszone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetDnszone : Get DNS zone
// Get details of a DNS zone.
func (dnsSvcs *DnsSvcsV1) GetDnszone(getDnszoneOptions *GetDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetDnszoneWithContext(context.Background(), getDnszoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDnszoneWithContext is an alternate form of the GetDnszone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetDnszoneWithContext(ctx context.Context, getDnszoneOptions *GetDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDnszoneOptions, "getDnszoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDnszoneOptions, "getDnszoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getDnszoneOptions.InstanceID,
		"dnszone_id":  *getDnszoneOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDnszoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetDnszone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDnszoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getDnszoneOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_dnszone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnszone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateDnszone : Update DNS zone
// Update the properties of a DNS zone.
func (dnsSvcs *DnsSvcsV1) UpdateDnszone(updateDnszoneOptions *UpdateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateDnszoneWithContext(context.Background(), updateDnszoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateDnszoneWithContext is an alternate form of the UpdateDnszone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateDnszoneWithContext(ctx context.Context, updateDnszoneOptions *UpdateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDnszoneOptions, "updateDnszoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateDnszoneOptions, "updateDnszoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateDnszoneOptions.InstanceID,
		"dnszone_id":  *updateDnszoneOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateDnszoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateDnszone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateDnszoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateDnszoneOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateDnszoneOptions.Description != nil {
		body["description"] = updateDnszoneOptions.Description
	}
	if updateDnszoneOptions.Label != nil {
		body["label"] = updateDnszoneOptions.Label
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_dnszone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnszone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListResourceRecords : List resource records
// List the Resource Records for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) ListResourceRecords(listResourceRecordsOptions *ListResourceRecordsOptions) (result *ListResourceRecords, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListResourceRecordsWithContext(context.Background(), listResourceRecordsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListResourceRecordsWithContext is an alternate form of the ListResourceRecords method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListResourceRecordsWithContext(ctx context.Context, listResourceRecordsOptions *ListResourceRecordsOptions) (result *ListResourceRecords, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listResourceRecordsOptions, "listResourceRecordsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listResourceRecordsOptions, "listResourceRecordsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listResourceRecordsOptions.InstanceID,
		"dnszone_id":  *listResourceRecordsOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listResourceRecordsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListResourceRecords")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listResourceRecordsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listResourceRecordsOptions.XCorrelationID))
	}

	if listResourceRecordsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listResourceRecordsOptions.Offset))
	}
	if listResourceRecordsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listResourceRecordsOptions.Limit))
	}
	if listResourceRecordsOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listResourceRecordsOptions.Type))
	}
	if listResourceRecordsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listResourceRecordsOptions.Name))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_resource_records", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListResourceRecords)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateResourceRecord : Create resource record
// Create a resource record for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) CreateResourceRecord(createResourceRecordOptions *CreateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateResourceRecordWithContext(context.Background(), createResourceRecordOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateResourceRecordWithContext is an alternate form of the CreateResourceRecord method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateResourceRecordWithContext(ctx context.Context, createResourceRecordOptions *CreateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createResourceRecordOptions, "createResourceRecordOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createResourceRecordOptions, "createResourceRecordOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createResourceRecordOptions.InstanceID,
		"dnszone_id":  *createResourceRecordOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createResourceRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateResourceRecord")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createResourceRecordOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createResourceRecordOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createResourceRecordOptions.Type != nil {
		body["type"] = createResourceRecordOptions.Type
	}
	if createResourceRecordOptions.Name != nil {
		body["name"] = createResourceRecordOptions.Name
	}
	if createResourceRecordOptions.Rdata != nil {
		body["rdata"] = createResourceRecordOptions.Rdata
	}
	if createResourceRecordOptions.TTL != nil {
		body["ttl"] = createResourceRecordOptions.TTL
	}
	if createResourceRecordOptions.Service != nil {
		body["service"] = createResourceRecordOptions.Service
	}
	if createResourceRecordOptions.Protocol != nil {
		body["protocol"] = createResourceRecordOptions.Protocol
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_resource_record", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceRecord)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteResourceRecord : Delete resource record
// Delete a resource record.
func (dnsSvcs *DnsSvcsV1) DeleteResourceRecord(deleteResourceRecordOptions *DeleteResourceRecordOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteResourceRecordWithContext(context.Background(), deleteResourceRecordOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteResourceRecordWithContext is an alternate form of the DeleteResourceRecord method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteResourceRecordWithContext(ctx context.Context, deleteResourceRecordOptions *DeleteResourceRecordOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteResourceRecordOptions, "deleteResourceRecordOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteResourceRecordOptions, "deleteResourceRecordOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteResourceRecordOptions.InstanceID,
		"dnszone_id":  *deleteResourceRecordOptions.DnszoneID,
		"record_id":   *deleteResourceRecordOptions.RecordID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records/{record_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteResourceRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteResourceRecord")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteResourceRecordOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteResourceRecordOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_resource_record", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetResourceRecord : Get resource record
// Get details of a resource record.
func (dnsSvcs *DnsSvcsV1) GetResourceRecord(getResourceRecordOptions *GetResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetResourceRecordWithContext(context.Background(), getResourceRecordOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetResourceRecordWithContext is an alternate form of the GetResourceRecord method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetResourceRecordWithContext(ctx context.Context, getResourceRecordOptions *GetResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceRecordOptions, "getResourceRecordOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getResourceRecordOptions, "getResourceRecordOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getResourceRecordOptions.InstanceID,
		"dnszone_id":  *getResourceRecordOptions.DnszoneID,
		"record_id":   *getResourceRecordOptions.RecordID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records/{record_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getResourceRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetResourceRecord")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getResourceRecordOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getResourceRecordOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_resource_record", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceRecord)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateResourceRecord : Update resource record
// Update the properties of a resource record.
func (dnsSvcs *DnsSvcsV1) UpdateResourceRecord(updateResourceRecordOptions *UpdateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateResourceRecordWithContext(context.Background(), updateResourceRecordOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateResourceRecordWithContext is an alternate form of the UpdateResourceRecord method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateResourceRecordWithContext(ctx context.Context, updateResourceRecordOptions *UpdateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateResourceRecordOptions, "updateResourceRecordOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateResourceRecordOptions, "updateResourceRecordOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateResourceRecordOptions.InstanceID,
		"dnszone_id":  *updateResourceRecordOptions.DnszoneID,
		"record_id":   *updateResourceRecordOptions.RecordID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records/{record_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateResourceRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateResourceRecord")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateResourceRecordOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateResourceRecordOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateResourceRecordOptions.Name != nil {
		body["name"] = updateResourceRecordOptions.Name
	}
	if updateResourceRecordOptions.Rdata != nil {
		body["rdata"] = updateResourceRecordOptions.Rdata
	}
	if updateResourceRecordOptions.TTL != nil {
		body["ttl"] = updateResourceRecordOptions.TTL
	}
	if updateResourceRecordOptions.Service != nil {
		body["service"] = updateResourceRecordOptions.Service
	}
	if updateResourceRecordOptions.Protocol != nil {
		body["protocol"] = updateResourceRecordOptions.Protocol
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_resource_record", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceRecord)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ExportResourceRecords : Export resource records to a zone file
// Export resource records to a zone file.
func (dnsSvcs *DnsSvcsV1) ExportResourceRecords(exportResourceRecordsOptions *ExportResourceRecordsOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ExportResourceRecordsWithContext(context.Background(), exportResourceRecordsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ExportResourceRecordsWithContext is an alternate form of the ExportResourceRecords method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ExportResourceRecordsWithContext(ctx context.Context, exportResourceRecordsOptions *ExportResourceRecordsOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(exportResourceRecordsOptions, "exportResourceRecordsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(exportResourceRecordsOptions, "exportResourceRecordsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *exportResourceRecordsOptions.InstanceID,
		"dnszone_id":  *exportResourceRecordsOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/export_resource_records`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range exportResourceRecordsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ExportResourceRecords")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "text/plain; charset=utf-8")
	if exportResourceRecordsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*exportResourceRecordsOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, &result)
	if err != nil {
		core.EnrichHTTPProblem(err, "export_resource_records", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ImportResourceRecords : Import resource records from a zone file
// Import resource records from a zone file. The maximum size of a zone file is 8MB.
func (dnsSvcs *DnsSvcsV1) ImportResourceRecords(importResourceRecordsOptions *ImportResourceRecordsOptions) (result *ImportResourceRecordsResp, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ImportResourceRecordsWithContext(context.Background(), importResourceRecordsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ImportResourceRecordsWithContext is an alternate form of the ImportResourceRecords method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ImportResourceRecordsWithContext(ctx context.Context, importResourceRecordsOptions *ImportResourceRecordsOptions) (result *ImportResourceRecordsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(importResourceRecordsOptions, "importResourceRecordsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(importResourceRecordsOptions, "importResourceRecordsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}
	if importResourceRecordsOptions.File == nil {
		err = core.SDKErrorf(nil, "file must be supplied", "condition-not-met", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *importResourceRecordsOptions.InstanceID,
		"dnszone_id":  *importResourceRecordsOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/import_resource_records`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range importResourceRecordsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ImportResourceRecords")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if importResourceRecordsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*importResourceRecordsOptions.XCorrelationID))
	}

	if importResourceRecordsOptions.File != nil {
		builder.AddFormData("file", "filename",
			core.StringNilMapper(importResourceRecordsOptions.FileContentType), importResourceRecordsOptions.File)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "import_resource_records", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImportResourceRecordsResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListPermittedNetworks : List permitted networks
// List the permitted networks for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) ListPermittedNetworks(listPermittedNetworksOptions *ListPermittedNetworksOptions) (result *ListPermittedNetworks, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListPermittedNetworksWithContext(context.Background(), listPermittedNetworksOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListPermittedNetworksWithContext is an alternate form of the ListPermittedNetworks method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListPermittedNetworksWithContext(ctx context.Context, listPermittedNetworksOptions *ListPermittedNetworksOptions) (result *ListPermittedNetworks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPermittedNetworksOptions, "listPermittedNetworksOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listPermittedNetworksOptions, "listPermittedNetworksOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listPermittedNetworksOptions.InstanceID,
		"dnszone_id":  *listPermittedNetworksOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/permitted_networks`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listPermittedNetworksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListPermittedNetworks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listPermittedNetworksOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listPermittedNetworksOptions.XCorrelationID))
	}

	if listPermittedNetworksOptions.Accounts != nil {
		builder.AddQuery("accounts", fmt.Sprint(*listPermittedNetworksOptions.Accounts))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_permitted_networks", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListPermittedNetworks)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreatePermittedNetwork : Create permitted network
// Create a permitted network for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) CreatePermittedNetwork(createPermittedNetworkOptions *CreatePermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreatePermittedNetworkWithContext(context.Background(), createPermittedNetworkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreatePermittedNetworkWithContext is an alternate form of the CreatePermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreatePermittedNetworkWithContext(ctx context.Context, createPermittedNetworkOptions *CreatePermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPermittedNetworkOptions, "createPermittedNetworkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createPermittedNetworkOptions, "createPermittedNetworkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createPermittedNetworkOptions.InstanceID,
		"dnszone_id":  *createPermittedNetworkOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/permitted_networks`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createPermittedNetworkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreatePermittedNetwork")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createPermittedNetworkOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createPermittedNetworkOptions.XCorrelationID))
	}

	if createPermittedNetworkOptions.Accounts != nil {
		builder.AddQuery("accounts", fmt.Sprint(*createPermittedNetworkOptions.Accounts))
	}

	body := make(map[string]interface{})
	if createPermittedNetworkOptions.Type != nil {
		body["type"] = createPermittedNetworkOptions.Type
	}
	if createPermittedNetworkOptions.PermittedNetwork != nil {
		body["permitted_network"] = createPermittedNetworkOptions.PermittedNetwork
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_permitted_network", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeletePermittedNetwork : Remove permitted network
// Remove a permitted network.
func (dnsSvcs *DnsSvcsV1) DeletePermittedNetwork(deletePermittedNetworkOptions *DeletePermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.DeletePermittedNetworkWithContext(context.Background(), deletePermittedNetworkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeletePermittedNetworkWithContext is an alternate form of the DeletePermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeletePermittedNetworkWithContext(ctx context.Context, deletePermittedNetworkOptions *DeletePermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePermittedNetworkOptions, "deletePermittedNetworkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deletePermittedNetworkOptions, "deletePermittedNetworkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":          *deletePermittedNetworkOptions.InstanceID,
		"dnszone_id":           *deletePermittedNetworkOptions.DnszoneID,
		"permitted_network_id": *deletePermittedNetworkOptions.PermittedNetworkID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/permitted_networks/{permitted_network_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deletePermittedNetworkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeletePermittedNetwork")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deletePermittedNetworkOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deletePermittedNetworkOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_permitted_network", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetPermittedNetwork : Get permitted network
// Get details of a permitted network.
func (dnsSvcs *DnsSvcsV1) GetPermittedNetwork(getPermittedNetworkOptions *GetPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetPermittedNetworkWithContext(context.Background(), getPermittedNetworkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPermittedNetworkWithContext is an alternate form of the GetPermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetPermittedNetworkWithContext(ctx context.Context, getPermittedNetworkOptions *GetPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPermittedNetworkOptions, "getPermittedNetworkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPermittedNetworkOptions, "getPermittedNetworkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":          *getPermittedNetworkOptions.InstanceID,
		"dnszone_id":           *getPermittedNetworkOptions.DnszoneID,
		"permitted_network_id": *getPermittedNetworkOptions.PermittedNetworkID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/permitted_networks/{permitted_network_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getPermittedNetworkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetPermittedNetwork")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getPermittedNetworkOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getPermittedNetworkOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_permitted_network", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListLoadBalancers : List load balancers
// List the Global Load Balancers for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) ListLoadBalancers(listLoadBalancersOptions *ListLoadBalancersOptions) (result *ListLoadBalancers, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListLoadBalancersWithContext(context.Background(), listLoadBalancersOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListLoadBalancersWithContext is an alternate form of the ListLoadBalancers method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListLoadBalancersWithContext(ctx context.Context, listLoadBalancersOptions *ListLoadBalancersOptions) (result *ListLoadBalancers, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listLoadBalancersOptions, "listLoadBalancersOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listLoadBalancersOptions, "listLoadBalancersOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listLoadBalancersOptions.InstanceID,
		"dnszone_id":  *listLoadBalancersOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listLoadBalancersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListLoadBalancers")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listLoadBalancersOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listLoadBalancersOptions.XCorrelationID))
	}

	if listLoadBalancersOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listLoadBalancersOptions.Offset))
	}
	if listLoadBalancersOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listLoadBalancersOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_load_balancers", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListLoadBalancers)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateLoadBalancer : Create load balancer
// Create a load balancer for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) CreateLoadBalancer(createLoadBalancerOptions *CreateLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateLoadBalancerWithContext(context.Background(), createLoadBalancerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateLoadBalancerWithContext is an alternate form of the CreateLoadBalancer method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateLoadBalancerWithContext(ctx context.Context, createLoadBalancerOptions *CreateLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createLoadBalancerOptions, "createLoadBalancerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createLoadBalancerOptions, "createLoadBalancerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createLoadBalancerOptions.InstanceID,
		"dnszone_id":  *createLoadBalancerOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createLoadBalancerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateLoadBalancer")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createLoadBalancerOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createLoadBalancerOptions.XCorrelationID))
	}

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
	if createLoadBalancerOptions.Enabled != nil {
		body["enabled"] = createLoadBalancerOptions.Enabled
	}
	if createLoadBalancerOptions.TTL != nil {
		body["ttl"] = createLoadBalancerOptions.TTL
	}
	if createLoadBalancerOptions.AzPools != nil {
		body["az_pools"] = createLoadBalancerOptions.AzPools
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_load_balancer", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancer)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteLoadBalancer : Delete load balancer
// Delete a load balancer.
func (dnsSvcs *DnsSvcsV1) DeleteLoadBalancer(deleteLoadBalancerOptions *DeleteLoadBalancerOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteLoadBalancerWithContext(context.Background(), deleteLoadBalancerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteLoadBalancerWithContext is an alternate form of the DeleteLoadBalancer method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteLoadBalancerWithContext(ctx context.Context, deleteLoadBalancerOptions *DeleteLoadBalancerOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLoadBalancerOptions, "deleteLoadBalancerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteLoadBalancerOptions, "deleteLoadBalancerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteLoadBalancerOptions.InstanceID,
		"dnszone_id":  *deleteLoadBalancerOptions.DnszoneID,
		"lb_id":       *deleteLoadBalancerOptions.LbID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers/{lb_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteLoadBalancerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteLoadBalancer")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteLoadBalancerOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteLoadBalancerOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_load_balancer", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetLoadBalancer : Get load balancer
// Get details of a load balancer.
func (dnsSvcs *DnsSvcsV1) GetLoadBalancer(getLoadBalancerOptions *GetLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetLoadBalancerWithContext(context.Background(), getLoadBalancerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLoadBalancerWithContext is an alternate form of the GetLoadBalancer method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetLoadBalancerWithContext(ctx context.Context, getLoadBalancerOptions *GetLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLoadBalancerOptions, "getLoadBalancerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getLoadBalancerOptions, "getLoadBalancerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getLoadBalancerOptions.InstanceID,
		"dnszone_id":  *getLoadBalancerOptions.DnszoneID,
		"lb_id":       *getLoadBalancerOptions.LbID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers/{lb_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLoadBalancerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetLoadBalancer")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getLoadBalancerOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getLoadBalancerOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_load_balancer", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancer)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateLoadBalancer : Update load balancer
// Update the properties of a load balancer.
func (dnsSvcs *DnsSvcsV1) UpdateLoadBalancer(updateLoadBalancerOptions *UpdateLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateLoadBalancerWithContext(context.Background(), updateLoadBalancerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateLoadBalancerWithContext is an alternate form of the UpdateLoadBalancer method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateLoadBalancerWithContext(ctx context.Context, updateLoadBalancerOptions *UpdateLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLoadBalancerOptions, "updateLoadBalancerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateLoadBalancerOptions, "updateLoadBalancerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateLoadBalancerOptions.InstanceID,
		"dnszone_id":  *updateLoadBalancerOptions.DnszoneID,
		"lb_id":       *updateLoadBalancerOptions.LbID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers/{lb_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateLoadBalancerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateLoadBalancer")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateLoadBalancerOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateLoadBalancerOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateLoadBalancerOptions.Name != nil {
		body["name"] = updateLoadBalancerOptions.Name
	}
	if updateLoadBalancerOptions.Description != nil {
		body["description"] = updateLoadBalancerOptions.Description
	}
	if updateLoadBalancerOptions.Enabled != nil {
		body["enabled"] = updateLoadBalancerOptions.Enabled
	}
	if updateLoadBalancerOptions.TTL != nil {
		body["ttl"] = updateLoadBalancerOptions.TTL
	}
	if updateLoadBalancerOptions.FallbackPool != nil {
		body["fallback_pool"] = updateLoadBalancerOptions.FallbackPool
	}
	if updateLoadBalancerOptions.DefaultPools != nil {
		body["default_pools"] = updateLoadBalancerOptions.DefaultPools
	}
	if updateLoadBalancerOptions.AzPools != nil {
		body["az_pools"] = updateLoadBalancerOptions.AzPools
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_load_balancer", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancer)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListPools : List load balancer pools
// List the load balancer pools.
func (dnsSvcs *DnsSvcsV1) ListPools(listPoolsOptions *ListPoolsOptions) (result *ListPools, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListPoolsWithContext(context.Background(), listPoolsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListPoolsWithContext is an alternate form of the ListPools method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListPoolsWithContext(ctx context.Context, listPoolsOptions *ListPoolsOptions) (result *ListPools, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPoolsOptions, "listPoolsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listPoolsOptions, "listPoolsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listPoolsOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/pools`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listPoolsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListPools")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listPoolsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listPoolsOptions.XCorrelationID))
	}

	if listPoolsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listPoolsOptions.Offset))
	}
	if listPoolsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listPoolsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_pools", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListPools)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreatePool : Create load balancer pool
// Create a load balancer pool.
func (dnsSvcs *DnsSvcsV1) CreatePool(createPoolOptions *CreatePoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreatePoolWithContext(context.Background(), createPoolOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreatePoolWithContext is an alternate form of the CreatePool method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreatePoolWithContext(ctx context.Context, createPoolOptions *CreatePoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPoolOptions, "createPoolOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createPoolOptions, "createPoolOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createPoolOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/pools`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createPoolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreatePool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createPoolOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createPoolOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createPoolOptions.Name != nil {
		body["name"] = createPoolOptions.Name
	}
	if createPoolOptions.Origins != nil {
		body["origins"] = createPoolOptions.Origins
	}
	if createPoolOptions.Description != nil {
		body["description"] = createPoolOptions.Description
	}
	if createPoolOptions.Enabled != nil {
		body["enabled"] = createPoolOptions.Enabled
	}
	if createPoolOptions.HealthyOriginsThreshold != nil {
		body["healthy_origins_threshold"] = createPoolOptions.HealthyOriginsThreshold
	}
	if createPoolOptions.Monitor != nil {
		body["monitor"] = createPoolOptions.Monitor
	}
	if createPoolOptions.NotificationChannel != nil {
		body["notification_channel"] = createPoolOptions.NotificationChannel
	}
	if createPoolOptions.HealthcheckRegion != nil {
		body["healthcheck_region"] = createPoolOptions.HealthcheckRegion
	}
	if createPoolOptions.HealthcheckSubnets != nil {
		body["healthcheck_subnets"] = createPoolOptions.HealthcheckSubnets
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_pool", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPool)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeletePool : Delete load balancer pool
// Delete a load balancer pool.
func (dnsSvcs *DnsSvcsV1) DeletePool(deletePoolOptions *DeletePoolOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeletePoolWithContext(context.Background(), deletePoolOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeletePoolWithContext is an alternate form of the DeletePool method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeletePoolWithContext(ctx context.Context, deletePoolOptions *DeletePoolOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePoolOptions, "deletePoolOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deletePoolOptions, "deletePoolOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deletePoolOptions.InstanceID,
		"pool_id":     *deletePoolOptions.PoolID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/pools/{pool_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deletePoolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeletePool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deletePoolOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deletePoolOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_pool", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetPool : Get load balancer pool
// Get details of a load balancer pool.
func (dnsSvcs *DnsSvcsV1) GetPool(getPoolOptions *GetPoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetPoolWithContext(context.Background(), getPoolOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPoolWithContext is an alternate form of the GetPool method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetPoolWithContext(ctx context.Context, getPoolOptions *GetPoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPoolOptions, "getPoolOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPoolOptions, "getPoolOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getPoolOptions.InstanceID,
		"pool_id":     *getPoolOptions.PoolID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/pools/{pool_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getPoolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetPool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getPoolOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getPoolOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_pool", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPool)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdatePool : Update load balancer pool
// Update the properties of a load balancer pool.
func (dnsSvcs *DnsSvcsV1) UpdatePool(updatePoolOptions *UpdatePoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdatePoolWithContext(context.Background(), updatePoolOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdatePoolWithContext is an alternate form of the UpdatePool method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdatePoolWithContext(ctx context.Context, updatePoolOptions *UpdatePoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePoolOptions, "updatePoolOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updatePoolOptions, "updatePoolOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updatePoolOptions.InstanceID,
		"pool_id":     *updatePoolOptions.PoolID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/pools/{pool_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updatePoolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdatePool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updatePoolOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updatePoolOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updatePoolOptions.Name != nil {
		body["name"] = updatePoolOptions.Name
	}
	if updatePoolOptions.Description != nil {
		body["description"] = updatePoolOptions.Description
	}
	if updatePoolOptions.Enabled != nil {
		body["enabled"] = updatePoolOptions.Enabled
	}
	if updatePoolOptions.HealthyOriginsThreshold != nil {
		body["healthy_origins_threshold"] = updatePoolOptions.HealthyOriginsThreshold
	}
	if updatePoolOptions.Origins != nil {
		body["origins"] = updatePoolOptions.Origins
	}
	if updatePoolOptions.Monitor != nil {
		body["monitor"] = updatePoolOptions.Monitor
	}
	if updatePoolOptions.NotificationChannel != nil {
		body["notification_channel"] = updatePoolOptions.NotificationChannel
	}
	if updatePoolOptions.HealthcheckRegion != nil {
		body["healthcheck_region"] = updatePoolOptions.HealthcheckRegion
	}
	if updatePoolOptions.HealthcheckSubnets != nil {
		body["healthcheck_subnets"] = updatePoolOptions.HealthcheckSubnets
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_pool", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPool)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListMonitors : List load balancer monitors
// List the load balancer monitors.
func (dnsSvcs *DnsSvcsV1) ListMonitors(listMonitorsOptions *ListMonitorsOptions) (result *ListMonitors, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListMonitorsWithContext(context.Background(), listMonitorsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListMonitorsWithContext is an alternate form of the ListMonitors method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListMonitorsWithContext(ctx context.Context, listMonitorsOptions *ListMonitorsOptions) (result *ListMonitors, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listMonitorsOptions, "listMonitorsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listMonitorsOptions, "listMonitorsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listMonitorsOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/monitors`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listMonitorsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListMonitors")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listMonitorsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listMonitorsOptions.XCorrelationID))
	}

	if listMonitorsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listMonitorsOptions.Offset))
	}
	if listMonitorsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listMonitorsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_monitors", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListMonitors)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateMonitor : Create load balancer monitor
// Create a load balancer monitor.
func (dnsSvcs *DnsSvcsV1) CreateMonitor(createMonitorOptions *CreateMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateMonitorWithContext(context.Background(), createMonitorOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateMonitorWithContext is an alternate form of the CreateMonitor method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateMonitorWithContext(ctx context.Context, createMonitorOptions *CreateMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createMonitorOptions, "createMonitorOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createMonitorOptions, "createMonitorOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createMonitorOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/monitors`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createMonitorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateMonitor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createMonitorOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createMonitorOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createMonitorOptions.Name != nil {
		body["name"] = createMonitorOptions.Name
	}
	if createMonitorOptions.Type != nil {
		body["type"] = createMonitorOptions.Type
	}
	if createMonitorOptions.Description != nil {
		body["description"] = createMonitorOptions.Description
	}
	if createMonitorOptions.Port != nil {
		body["port"] = createMonitorOptions.Port
	}
	if createMonitorOptions.Interval != nil {
		body["interval"] = createMonitorOptions.Interval
	}
	if createMonitorOptions.Retries != nil {
		body["retries"] = createMonitorOptions.Retries
	}
	if createMonitorOptions.Timeout != nil {
		body["timeout"] = createMonitorOptions.Timeout
	}
	if createMonitorOptions.Method != nil {
		body["method"] = createMonitorOptions.Method
	}
	if createMonitorOptions.Path != nil {
		body["path"] = createMonitorOptions.Path
	}
	if createMonitorOptions.HeadersVar != nil {
		body["headers"] = createMonitorOptions.HeadersVar
	}
	if createMonitorOptions.AllowInsecure != nil {
		body["allow_insecure"] = createMonitorOptions.AllowInsecure
	}
	if createMonitorOptions.ExpectedCodes != nil {
		body["expected_codes"] = createMonitorOptions.ExpectedCodes
	}
	if createMonitorOptions.ExpectedBody != nil {
		body["expected_body"] = createMonitorOptions.ExpectedBody
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_monitor", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitor)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteMonitor : Delete load balancer monitor
// Delete a load balancer monitor.
func (dnsSvcs *DnsSvcsV1) DeleteMonitor(deleteMonitorOptions *DeleteMonitorOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteMonitorWithContext(context.Background(), deleteMonitorOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteMonitorWithContext is an alternate form of the DeleteMonitor method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteMonitorWithContext(ctx context.Context, deleteMonitorOptions *DeleteMonitorOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteMonitorOptions, "deleteMonitorOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteMonitorOptions, "deleteMonitorOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteMonitorOptions.InstanceID,
		"monitor_id":  *deleteMonitorOptions.MonitorID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/monitors/{monitor_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteMonitorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteMonitor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteMonitorOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteMonitorOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_monitor", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetMonitor : Get load balancer monitor
// Get details of a load balancer monitor.
func (dnsSvcs *DnsSvcsV1) GetMonitor(getMonitorOptions *GetMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetMonitorWithContext(context.Background(), getMonitorOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetMonitorWithContext is an alternate form of the GetMonitor method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetMonitorWithContext(ctx context.Context, getMonitorOptions *GetMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getMonitorOptions, "getMonitorOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getMonitorOptions, "getMonitorOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getMonitorOptions.InstanceID,
		"monitor_id":  *getMonitorOptions.MonitorID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/monitors/{monitor_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getMonitorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetMonitor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getMonitorOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getMonitorOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_monitor", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitor)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateMonitor : Update load balancer monitor
// Update the properties of a load balancer monitor.
func (dnsSvcs *DnsSvcsV1) UpdateMonitor(updateMonitorOptions *UpdateMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateMonitorWithContext(context.Background(), updateMonitorOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateMonitorWithContext is an alternate form of the UpdateMonitor method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateMonitorWithContext(ctx context.Context, updateMonitorOptions *UpdateMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateMonitorOptions, "updateMonitorOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateMonitorOptions, "updateMonitorOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateMonitorOptions.InstanceID,
		"monitor_id":  *updateMonitorOptions.MonitorID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/monitors/{monitor_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateMonitorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateMonitor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateMonitorOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateMonitorOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateMonitorOptions.Name != nil {
		body["name"] = updateMonitorOptions.Name
	}
	if updateMonitorOptions.Description != nil {
		body["description"] = updateMonitorOptions.Description
	}
	if updateMonitorOptions.Type != nil {
		body["type"] = updateMonitorOptions.Type
	}
	if updateMonitorOptions.Port != nil {
		body["port"] = updateMonitorOptions.Port
	}
	if updateMonitorOptions.Interval != nil {
		body["interval"] = updateMonitorOptions.Interval
	}
	if updateMonitorOptions.Retries != nil {
		body["retries"] = updateMonitorOptions.Retries
	}
	if updateMonitorOptions.Timeout != nil {
		body["timeout"] = updateMonitorOptions.Timeout
	}
	if updateMonitorOptions.Method != nil {
		body["method"] = updateMonitorOptions.Method
	}
	if updateMonitorOptions.Path != nil {
		body["path"] = updateMonitorOptions.Path
	}
	if updateMonitorOptions.HeadersVar != nil {
		body["headers"] = updateMonitorOptions.HeadersVar
	}
	if updateMonitorOptions.AllowInsecure != nil {
		body["allow_insecure"] = updateMonitorOptions.AllowInsecure
	}
	if updateMonitorOptions.ExpectedCodes != nil {
		body["expected_codes"] = updateMonitorOptions.ExpectedCodes
	}
	if updateMonitorOptions.ExpectedBody != nil {
		body["expected_body"] = updateMonitorOptions.ExpectedBody
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_monitor", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitor)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListCustomResolvers : List custom resolvers
// List the custom resolvers.
func (dnsSvcs *DnsSvcsV1) ListCustomResolvers(listCustomResolversOptions *ListCustomResolversOptions) (result *CustomResolverList, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListCustomResolversWithContext(context.Background(), listCustomResolversOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListCustomResolversWithContext is an alternate form of the ListCustomResolvers method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListCustomResolversWithContext(ctx context.Context, listCustomResolversOptions *ListCustomResolversOptions) (result *CustomResolverList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listCustomResolversOptions, "listCustomResolversOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listCustomResolversOptions, "listCustomResolversOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listCustomResolversOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listCustomResolversOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListCustomResolvers")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listCustomResolversOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listCustomResolversOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_custom_resolvers", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomResolverList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateCustomResolver : Create a custom resolver
// Create a custom resolver.
func (dnsSvcs *DnsSvcsV1) CreateCustomResolver(createCustomResolverOptions *CreateCustomResolverOptions) (result *CustomResolver, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateCustomResolverWithContext(context.Background(), createCustomResolverOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateCustomResolverWithContext is an alternate form of the CreateCustomResolver method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateCustomResolverWithContext(ctx context.Context, createCustomResolverOptions *CreateCustomResolverOptions) (result *CustomResolver, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCustomResolverOptions, "createCustomResolverOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createCustomResolverOptions, "createCustomResolverOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createCustomResolverOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createCustomResolverOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateCustomResolver")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createCustomResolverOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createCustomResolverOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createCustomResolverOptions.Name != nil {
		body["name"] = createCustomResolverOptions.Name
	}
	if createCustomResolverOptions.Description != nil {
		body["description"] = createCustomResolverOptions.Description
	}
	if createCustomResolverOptions.Locations != nil {
		body["locations"] = createCustomResolverOptions.Locations
	}
	if createCustomResolverOptions.Profile != nil {
		body["profile"] = createCustomResolverOptions.Profile
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_custom_resolver", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomResolver)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteCustomResolver : Delete a custom resolver
// Delete a custom resolver.
func (dnsSvcs *DnsSvcsV1) DeleteCustomResolver(deleteCustomResolverOptions *DeleteCustomResolverOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteCustomResolverWithContext(context.Background(), deleteCustomResolverOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteCustomResolverWithContext is an alternate form of the DeleteCustomResolver method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteCustomResolverWithContext(ctx context.Context, deleteCustomResolverOptions *DeleteCustomResolverOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCustomResolverOptions, "deleteCustomResolverOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteCustomResolverOptions, "deleteCustomResolverOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteCustomResolverOptions.InstanceID,
		"resolver_id": *deleteCustomResolverOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteCustomResolverOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteCustomResolver")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteCustomResolverOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteCustomResolverOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_custom_resolver", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetCustomResolver : Get a custom resolver
// Get details of a custom resolver.
func (dnsSvcs *DnsSvcsV1) GetCustomResolver(getCustomResolverOptions *GetCustomResolverOptions) (result *CustomResolver, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetCustomResolverWithContext(context.Background(), getCustomResolverOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetCustomResolverWithContext is an alternate form of the GetCustomResolver method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetCustomResolverWithContext(ctx context.Context, getCustomResolverOptions *GetCustomResolverOptions) (result *CustomResolver, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCustomResolverOptions, "getCustomResolverOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getCustomResolverOptions, "getCustomResolverOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getCustomResolverOptions.InstanceID,
		"resolver_id": *getCustomResolverOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getCustomResolverOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetCustomResolver")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getCustomResolverOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getCustomResolverOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_custom_resolver", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomResolver)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateCustomResolver : Update a custom resolver
// Update the properties of a custom resolver.
func (dnsSvcs *DnsSvcsV1) UpdateCustomResolver(updateCustomResolverOptions *UpdateCustomResolverOptions) (result *CustomResolver, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateCustomResolverWithContext(context.Background(), updateCustomResolverOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateCustomResolverWithContext is an alternate form of the UpdateCustomResolver method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateCustomResolverWithContext(ctx context.Context, updateCustomResolverOptions *UpdateCustomResolverOptions) (result *CustomResolver, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCustomResolverOptions, "updateCustomResolverOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateCustomResolverOptions, "updateCustomResolverOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateCustomResolverOptions.InstanceID,
		"resolver_id": *updateCustomResolverOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateCustomResolverOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateCustomResolver")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateCustomResolverOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateCustomResolverOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateCustomResolverOptions.Name != nil {
		body["name"] = updateCustomResolverOptions.Name
	}
	if updateCustomResolverOptions.Description != nil {
		body["description"] = updateCustomResolverOptions.Description
	}
	if updateCustomResolverOptions.Enabled != nil {
		body["enabled"] = updateCustomResolverOptions.Enabled
	}
	if updateCustomResolverOptions.Profile != nil {
		body["profile"] = updateCustomResolverOptions.Profile
	}
	if updateCustomResolverOptions.AllowDisruptiveUpdates != nil {
		body["allow_disruptive_updates"] = updateCustomResolverOptions.AllowDisruptiveUpdates
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_custom_resolver", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomResolver)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateCrLocationsOrder : Update the locations order of a custom resolver (DEPRECATED)
// Update the locations order of a custom resolver.
func (dnsSvcs *DnsSvcsV1) UpdateCrLocationsOrder(updateCrLocationsOrderOptions *UpdateCrLocationsOrderOptions) (result *CustomResolver, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateCrLocationsOrderWithContext(context.Background(), updateCrLocationsOrderOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateCrLocationsOrderWithContext is an alternate form of the UpdateCrLocationsOrder method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateCrLocationsOrderWithContext(ctx context.Context, updateCrLocationsOrderOptions *UpdateCrLocationsOrderOptions) (result *CustomResolver, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCrLocationsOrderOptions, "updateCrLocationsOrderOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateCrLocationsOrderOptions, "updateCrLocationsOrderOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateCrLocationsOrderOptions.InstanceID,
		"resolver_id": *updateCrLocationsOrderOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/locations_order`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateCrLocationsOrderOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateCrLocationsOrder")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateCrLocationsOrderOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateCrLocationsOrderOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateCrLocationsOrderOptions.Locations != nil {
		body["locations"] = updateCrLocationsOrderOptions.Locations
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_cr_locations_order", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomResolver)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// AddCustomResolverLocation : Add custom resolver location
// Add custom resolver location.
func (dnsSvcs *DnsSvcsV1) AddCustomResolverLocation(addCustomResolverLocationOptions *AddCustomResolverLocationOptions) (result *Location, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.AddCustomResolverLocationWithContext(context.Background(), addCustomResolverLocationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// AddCustomResolverLocationWithContext is an alternate form of the AddCustomResolverLocation method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) AddCustomResolverLocationWithContext(ctx context.Context, addCustomResolverLocationOptions *AddCustomResolverLocationOptions) (result *Location, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(addCustomResolverLocationOptions, "addCustomResolverLocationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(addCustomResolverLocationOptions, "addCustomResolverLocationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *addCustomResolverLocationOptions.InstanceID,
		"resolver_id": *addCustomResolverLocationOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/locations`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range addCustomResolverLocationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "AddCustomResolverLocation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if addCustomResolverLocationOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*addCustomResolverLocationOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if addCustomResolverLocationOptions.SubnetCrn != nil {
		body["subnet_crn"] = addCustomResolverLocationOptions.SubnetCrn
	}
	if addCustomResolverLocationOptions.Enabled != nil {
		body["enabled"] = addCustomResolverLocationOptions.Enabled
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "add_custom_resolver_location", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLocation)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateCustomResolverLocation : Update custom resolver location
// Update custom resolver location.
func (dnsSvcs *DnsSvcsV1) UpdateCustomResolverLocation(updateCustomResolverLocationOptions *UpdateCustomResolverLocationOptions) (result *Location, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateCustomResolverLocationWithContext(context.Background(), updateCustomResolverLocationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateCustomResolverLocationWithContext is an alternate form of the UpdateCustomResolverLocation method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateCustomResolverLocationWithContext(ctx context.Context, updateCustomResolverLocationOptions *UpdateCustomResolverLocationOptions) (result *Location, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCustomResolverLocationOptions, "updateCustomResolverLocationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateCustomResolverLocationOptions, "updateCustomResolverLocationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateCustomResolverLocationOptions.InstanceID,
		"resolver_id": *updateCustomResolverLocationOptions.ResolverID,
		"location_id": *updateCustomResolverLocationOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/locations/{location_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateCustomResolverLocationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateCustomResolverLocation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateCustomResolverLocationOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateCustomResolverLocationOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateCustomResolverLocationOptions.Enabled != nil {
		body["enabled"] = updateCustomResolverLocationOptions.Enabled
	}
	if updateCustomResolverLocationOptions.SubnetCrn != nil {
		body["subnet_crn"] = updateCustomResolverLocationOptions.SubnetCrn
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_custom_resolver_location", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLocation)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteCustomResolverLocation : Delete custom resolver location
// Delete custom resolver location.
func (dnsSvcs *DnsSvcsV1) DeleteCustomResolverLocation(deleteCustomResolverLocationOptions *DeleteCustomResolverLocationOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteCustomResolverLocationWithContext(context.Background(), deleteCustomResolverLocationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteCustomResolverLocationWithContext is an alternate form of the DeleteCustomResolverLocation method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteCustomResolverLocationWithContext(ctx context.Context, deleteCustomResolverLocationOptions *DeleteCustomResolverLocationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCustomResolverLocationOptions, "deleteCustomResolverLocationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteCustomResolverLocationOptions, "deleteCustomResolverLocationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteCustomResolverLocationOptions.InstanceID,
		"resolver_id": *deleteCustomResolverLocationOptions.ResolverID,
		"location_id": *deleteCustomResolverLocationOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/locations/{location_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteCustomResolverLocationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteCustomResolverLocation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteCustomResolverLocationOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteCustomResolverLocationOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_custom_resolver_location", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListForwardingRules : List forwarding rules
// List the forwarding rules of the given custom resolver.
func (dnsSvcs *DnsSvcsV1) ListForwardingRules(listForwardingRulesOptions *ListForwardingRulesOptions) (result *ForwardingRuleList, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListForwardingRulesWithContext(context.Background(), listForwardingRulesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListForwardingRulesWithContext is an alternate form of the ListForwardingRules method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListForwardingRulesWithContext(ctx context.Context, listForwardingRulesOptions *ListForwardingRulesOptions) (result *ForwardingRuleList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listForwardingRulesOptions, "listForwardingRulesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listForwardingRulesOptions, "listForwardingRulesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listForwardingRulesOptions.InstanceID,
		"resolver_id": *listForwardingRulesOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/forwarding_rules`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listForwardingRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListForwardingRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listForwardingRulesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listForwardingRulesOptions.XCorrelationID))
	}

	if listForwardingRulesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listForwardingRulesOptions.Offset))
	}
	if listForwardingRulesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listForwardingRulesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_forwarding_rules", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalForwardingRuleList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateForwardingRule : Create a forwarding rule
// Create a forwarding rule for the given custom resolver. The maximum number limit of forwarding rules depends on the
// custom resolver profile.
//
// * Essential profiles can have a maximum of 10 forwarding rules.
// * Advanced profiles can have a maximum of 50 forwarding rules.
// * Premier profiles can have a maximum of 100 forwarding rules.
func (dnsSvcs *DnsSvcsV1) CreateForwardingRule(createForwardingRuleOptions *CreateForwardingRuleOptions) (result *ForwardingRule, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateForwardingRuleWithContext(context.Background(), createForwardingRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateForwardingRuleWithContext is an alternate form of the CreateForwardingRule method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateForwardingRuleWithContext(ctx context.Context, createForwardingRuleOptions *CreateForwardingRuleOptions) (result *ForwardingRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createForwardingRuleOptions, "createForwardingRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createForwardingRuleOptions, "createForwardingRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createForwardingRuleOptions.InstanceID,
		"resolver_id": *createForwardingRuleOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/forwarding_rules`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createForwardingRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateForwardingRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createForwardingRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createForwardingRuleOptions.XCorrelationID))
	}

	_, err = builder.SetBodyContentJSON(createForwardingRuleOptions.ForwardingRuleInput)
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_forwarding_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalForwardingRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteForwardingRule : Delete a forwarding rule
// Delete a forwarding rule on the given custom resolver.
func (dnsSvcs *DnsSvcsV1) DeleteForwardingRule(deleteForwardingRuleOptions *DeleteForwardingRuleOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteForwardingRuleWithContext(context.Background(), deleteForwardingRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteForwardingRuleWithContext is an alternate form of the DeleteForwardingRule method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteForwardingRuleWithContext(ctx context.Context, deleteForwardingRuleOptions *DeleteForwardingRuleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteForwardingRuleOptions, "deleteForwardingRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteForwardingRuleOptions, "deleteForwardingRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteForwardingRuleOptions.InstanceID,
		"resolver_id": *deleteForwardingRuleOptions.ResolverID,
		"rule_id":     *deleteForwardingRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/forwarding_rules/{rule_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteForwardingRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteForwardingRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteForwardingRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteForwardingRuleOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_forwarding_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetForwardingRule : Get a forwarding rule
// Get details of a forwarding rule on the given custom resolver.
func (dnsSvcs *DnsSvcsV1) GetForwardingRule(getForwardingRuleOptions *GetForwardingRuleOptions) (result *ForwardingRule, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetForwardingRuleWithContext(context.Background(), getForwardingRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetForwardingRuleWithContext is an alternate form of the GetForwardingRule method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetForwardingRuleWithContext(ctx context.Context, getForwardingRuleOptions *GetForwardingRuleOptions) (result *ForwardingRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getForwardingRuleOptions, "getForwardingRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getForwardingRuleOptions, "getForwardingRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getForwardingRuleOptions.InstanceID,
		"resolver_id": *getForwardingRuleOptions.ResolverID,
		"rule_id":     *getForwardingRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/forwarding_rules/{rule_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getForwardingRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetForwardingRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getForwardingRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getForwardingRuleOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_forwarding_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalForwardingRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateForwardingRule : Update a forwarding rule
// Update the properties of a forwarding rule on the given custom resolver.
func (dnsSvcs *DnsSvcsV1) UpdateForwardingRule(updateForwardingRuleOptions *UpdateForwardingRuleOptions) (result *ForwardingRule, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateForwardingRuleWithContext(context.Background(), updateForwardingRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateForwardingRuleWithContext is an alternate form of the UpdateForwardingRule method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateForwardingRuleWithContext(ctx context.Context, updateForwardingRuleOptions *UpdateForwardingRuleOptions) (result *ForwardingRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateForwardingRuleOptions, "updateForwardingRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateForwardingRuleOptions, "updateForwardingRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateForwardingRuleOptions.InstanceID,
		"resolver_id": *updateForwardingRuleOptions.ResolverID,
		"rule_id":     *updateForwardingRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/forwarding_rules/{rule_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateForwardingRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateForwardingRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateForwardingRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateForwardingRuleOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateForwardingRuleOptions.Description != nil {
		body["description"] = updateForwardingRuleOptions.Description
	}
	if updateForwardingRuleOptions.Match != nil {
		body["match"] = updateForwardingRuleOptions.Match
	}
	if updateForwardingRuleOptions.ForwardTo != nil {
		body["forward_to"] = updateForwardingRuleOptions.ForwardTo
	}
	if updateForwardingRuleOptions.Views != nil {
		body["views"] = updateForwardingRuleOptions.Views
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_forwarding_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalForwardingRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateSecondaryZone : Create a secondary zone
// Create a secondary zone for a given custom resolver. The maximum number limit of secondary zones depends on the
// custom resolver profile.
//
// * Essential profiles can have a maximum of 10 secondary zones.
// * Advanced profiles can have a maximum of 50 secondary zones.
// * Premier profiles can have a maximum of 100 secondary zones.
func (dnsSvcs *DnsSvcsV1) CreateSecondaryZone(createSecondaryZoneOptions *CreateSecondaryZoneOptions) (result *SecondaryZone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateSecondaryZoneWithContext(context.Background(), createSecondaryZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateSecondaryZoneWithContext is an alternate form of the CreateSecondaryZone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateSecondaryZoneWithContext(ctx context.Context, createSecondaryZoneOptions *CreateSecondaryZoneOptions) (result *SecondaryZone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecondaryZoneOptions, "createSecondaryZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createSecondaryZoneOptions, "createSecondaryZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createSecondaryZoneOptions.InstanceID,
		"resolver_id": *createSecondaryZoneOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/secondary_zones`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createSecondaryZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateSecondaryZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createSecondaryZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createSecondaryZoneOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createSecondaryZoneOptions.Zone != nil {
		body["zone"] = createSecondaryZoneOptions.Zone
	}
	if createSecondaryZoneOptions.TransferFrom != nil {
		body["transfer_from"] = createSecondaryZoneOptions.TransferFrom
	}
	if createSecondaryZoneOptions.Description != nil {
		body["description"] = createSecondaryZoneOptions.Description
	}
	if createSecondaryZoneOptions.Enabled != nil {
		body["enabled"] = createSecondaryZoneOptions.Enabled
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_secondary_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecondaryZone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListSecondaryZones : List secondary zones
// List secondary zones for the custom resolver.
func (dnsSvcs *DnsSvcsV1) ListSecondaryZones(listSecondaryZonesOptions *ListSecondaryZonesOptions) (result *SecondaryZoneList, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListSecondaryZonesWithContext(context.Background(), listSecondaryZonesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListSecondaryZonesWithContext is an alternate form of the ListSecondaryZones method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListSecondaryZonesWithContext(ctx context.Context, listSecondaryZonesOptions *ListSecondaryZonesOptions) (result *SecondaryZoneList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSecondaryZonesOptions, "listSecondaryZonesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listSecondaryZonesOptions, "listSecondaryZonesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listSecondaryZonesOptions.InstanceID,
		"resolver_id": *listSecondaryZonesOptions.ResolverID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/secondary_zones`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listSecondaryZonesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListSecondaryZones")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listSecondaryZonesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listSecondaryZonesOptions.XCorrelationID))
	}

	if listSecondaryZonesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listSecondaryZonesOptions.Offset))
	}
	if listSecondaryZonesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSecondaryZonesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_secondary_zones", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecondaryZoneList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetSecondaryZone : Get a secondary zone
// Get details of a secondary zone for the custom resolver.
func (dnsSvcs *DnsSvcsV1) GetSecondaryZone(getSecondaryZoneOptions *GetSecondaryZoneOptions) (result *SecondaryZone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetSecondaryZoneWithContext(context.Background(), getSecondaryZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetSecondaryZoneWithContext is an alternate form of the GetSecondaryZone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetSecondaryZoneWithContext(ctx context.Context, getSecondaryZoneOptions *GetSecondaryZoneOptions) (result *SecondaryZone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecondaryZoneOptions, "getSecondaryZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getSecondaryZoneOptions, "getSecondaryZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":       *getSecondaryZoneOptions.InstanceID,
		"resolver_id":       *getSecondaryZoneOptions.ResolverID,
		"secondary_zone_id": *getSecondaryZoneOptions.SecondaryZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/secondary_zones/{secondary_zone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getSecondaryZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetSecondaryZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getSecondaryZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getSecondaryZoneOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_secondary_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecondaryZone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecondaryZone : Update a secondary zone
// Update a secondary zone for the custom resolver.
func (dnsSvcs *DnsSvcsV1) UpdateSecondaryZone(updateSecondaryZoneOptions *UpdateSecondaryZoneOptions) (result *SecondaryZone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateSecondaryZoneWithContext(context.Background(), updateSecondaryZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateSecondaryZoneWithContext is an alternate form of the UpdateSecondaryZone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateSecondaryZoneWithContext(ctx context.Context, updateSecondaryZoneOptions *UpdateSecondaryZoneOptions) (result *SecondaryZone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecondaryZoneOptions, "updateSecondaryZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateSecondaryZoneOptions, "updateSecondaryZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":       *updateSecondaryZoneOptions.InstanceID,
		"resolver_id":       *updateSecondaryZoneOptions.ResolverID,
		"secondary_zone_id": *updateSecondaryZoneOptions.SecondaryZoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/secondary_zones/{secondary_zone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateSecondaryZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateSecondaryZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateSecondaryZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateSecondaryZoneOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateSecondaryZoneOptions.Description != nil {
		body["description"] = updateSecondaryZoneOptions.Description
	}
	if updateSecondaryZoneOptions.Enabled != nil {
		body["enabled"] = updateSecondaryZoneOptions.Enabled
	}
	if updateSecondaryZoneOptions.TransferFrom != nil {
		body["transfer_from"] = updateSecondaryZoneOptions.TransferFrom
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_secondary_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecondaryZone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteSecondaryZone : Delete a secondary zone
// Delete a secondary zone for the custom resolver.
func (dnsSvcs *DnsSvcsV1) DeleteSecondaryZone(deleteSecondaryZoneOptions *DeleteSecondaryZoneOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteSecondaryZoneWithContext(context.Background(), deleteSecondaryZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteSecondaryZoneWithContext is an alternate form of the DeleteSecondaryZone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteSecondaryZoneWithContext(ctx context.Context, deleteSecondaryZoneOptions *DeleteSecondaryZoneOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSecondaryZoneOptions, "deleteSecondaryZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteSecondaryZoneOptions, "deleteSecondaryZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":       *deleteSecondaryZoneOptions.InstanceID,
		"resolver_id":       *deleteSecondaryZoneOptions.ResolverID,
		"secondary_zone_id": *deleteSecondaryZoneOptions.SecondaryZoneID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/custom_resolvers/{resolver_id}/secondary_zones/{secondary_zone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteSecondaryZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteSecondaryZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteSecondaryZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteSecondaryZoneOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_secondary_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListLinkedZones : List linked zones
// List linked zones in requestor's instance.
func (dnsSvcs *DnsSvcsV1) ListLinkedZones(listLinkedZonesOptions *ListLinkedZonesOptions) (result *LinkedDnszonesList, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListLinkedZonesWithContext(context.Background(), listLinkedZonesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListLinkedZonesWithContext is an alternate form of the ListLinkedZones method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListLinkedZonesWithContext(ctx context.Context, listLinkedZonesOptions *ListLinkedZonesOptions) (result *LinkedDnszonesList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listLinkedZonesOptions, "listLinkedZonesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listLinkedZonesOptions, "listLinkedZonesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listLinkedZonesOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listLinkedZonesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListLinkedZones")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listLinkedZonesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listLinkedZonesOptions.XCorrelationID))
	}

	if listLinkedZonesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listLinkedZonesOptions.Offset))
	}
	if listLinkedZonesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listLinkedZonesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_linked_zones", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLinkedDnszonesList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateLinkedZone : Create a linked zone
// Create a linked zone.
func (dnsSvcs *DnsSvcsV1) CreateLinkedZone(createLinkedZoneOptions *CreateLinkedZoneOptions) (result *LinkedDnszone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateLinkedZoneWithContext(context.Background(), createLinkedZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateLinkedZoneWithContext is an alternate form of the CreateLinkedZone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateLinkedZoneWithContext(ctx context.Context, createLinkedZoneOptions *CreateLinkedZoneOptions) (result *LinkedDnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createLinkedZoneOptions, "createLinkedZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createLinkedZoneOptions, "createLinkedZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createLinkedZoneOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createLinkedZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateLinkedZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createLinkedZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createLinkedZoneOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createLinkedZoneOptions.OwnerInstanceID != nil {
		body["owner_instance_id"] = createLinkedZoneOptions.OwnerInstanceID
	}
	if createLinkedZoneOptions.OwnerZoneID != nil {
		body["owner_zone_id"] = createLinkedZoneOptions.OwnerZoneID
	}
	if createLinkedZoneOptions.Description != nil {
		body["description"] = createLinkedZoneOptions.Description
	}
	if createLinkedZoneOptions.Label != nil {
		body["label"] = createLinkedZoneOptions.Label
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_linked_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLinkedDnszone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetLinkedZone : Get a linked zone
// Get details of a linked zone.
func (dnsSvcs *DnsSvcsV1) GetLinkedZone(getLinkedZoneOptions *GetLinkedZoneOptions) (result *LinkedDnszone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetLinkedZoneWithContext(context.Background(), getLinkedZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLinkedZoneWithContext is an alternate form of the GetLinkedZone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetLinkedZoneWithContext(ctx context.Context, getLinkedZoneOptions *GetLinkedZoneOptions) (result *LinkedDnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLinkedZoneOptions, "getLinkedZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getLinkedZoneOptions, "getLinkedZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":       *getLinkedZoneOptions.InstanceID,
		"linked_dnszone_id": *getLinkedZoneOptions.LinkedDnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones/{linked_dnszone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLinkedZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetLinkedZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getLinkedZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getLinkedZoneOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_linked_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLinkedDnszone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateLinkedZone : Update the properties of a linked zone
// Update the properties of a linked zone.
func (dnsSvcs *DnsSvcsV1) UpdateLinkedZone(updateLinkedZoneOptions *UpdateLinkedZoneOptions) (result *LinkedDnszone, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateLinkedZoneWithContext(context.Background(), updateLinkedZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateLinkedZoneWithContext is an alternate form of the UpdateLinkedZone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateLinkedZoneWithContext(ctx context.Context, updateLinkedZoneOptions *UpdateLinkedZoneOptions) (result *LinkedDnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLinkedZoneOptions, "updateLinkedZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateLinkedZoneOptions, "updateLinkedZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":       *updateLinkedZoneOptions.InstanceID,
		"linked_dnszone_id": *updateLinkedZoneOptions.LinkedDnszoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones/{linked_dnszone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateLinkedZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateLinkedZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateLinkedZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateLinkedZoneOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateLinkedZoneOptions.Description != nil {
		body["description"] = updateLinkedZoneOptions.Description
	}
	if updateLinkedZoneOptions.Label != nil {
		body["label"] = updateLinkedZoneOptions.Label
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_linked_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLinkedDnszone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteLinkedZone : Delete a linked zone
// Delete a linked zone.
func (dnsSvcs *DnsSvcsV1) DeleteLinkedZone(deleteLinkedZoneOptions *DeleteLinkedZoneOptions) (response *core.DetailedResponse, err error) {
	response, err = dnsSvcs.DeleteLinkedZoneWithContext(context.Background(), deleteLinkedZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteLinkedZoneWithContext is an alternate form of the DeleteLinkedZone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteLinkedZoneWithContext(ctx context.Context, deleteLinkedZoneOptions *DeleteLinkedZoneOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLinkedZoneOptions, "deleteLinkedZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteLinkedZoneOptions, "deleteLinkedZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":       *deleteLinkedZoneOptions.InstanceID,
		"linked_dnszone_id": *deleteLinkedZoneOptions.LinkedDnszoneID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones/{linked_dnszone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteLinkedZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteLinkedZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteLinkedZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteLinkedZoneOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_linked_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListDnszoneAccessRequests : List Access Requests
// List access requests in owner's instance.
func (dnsSvcs *DnsSvcsV1) ListDnszoneAccessRequests(listDnszoneAccessRequestsOptions *ListDnszoneAccessRequestsOptions) (result *AccessRequestsList, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListDnszoneAccessRequestsWithContext(context.Background(), listDnszoneAccessRequestsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListDnszoneAccessRequestsWithContext is an alternate form of the ListDnszoneAccessRequests method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListDnszoneAccessRequestsWithContext(ctx context.Context, listDnszoneAccessRequestsOptions *ListDnszoneAccessRequestsOptions) (result *AccessRequestsList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listDnszoneAccessRequestsOptions, "listDnszoneAccessRequestsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listDnszoneAccessRequestsOptions, "listDnszoneAccessRequestsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listDnszoneAccessRequestsOptions.InstanceID,
		"dnszone_id":  *listDnszoneAccessRequestsOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/access_requests`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listDnszoneAccessRequestsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListDnszoneAccessRequests")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listDnszoneAccessRequestsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listDnszoneAccessRequestsOptions.XCorrelationID))
	}

	if listDnszoneAccessRequestsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listDnszoneAccessRequestsOptions.Offset))
	}
	if listDnszoneAccessRequestsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listDnszoneAccessRequestsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_dnszone_access_requests", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessRequestsList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDnszoneAccessRequest : Get an access request
// Get details of an access request.
func (dnsSvcs *DnsSvcsV1) GetDnszoneAccessRequest(getDnszoneAccessRequestOptions *GetDnszoneAccessRequestOptions) (result *AccessRequest, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetDnszoneAccessRequestWithContext(context.Background(), getDnszoneAccessRequestOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDnszoneAccessRequestWithContext is an alternate form of the GetDnszoneAccessRequest method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetDnszoneAccessRequestWithContext(ctx context.Context, getDnszoneAccessRequestOptions *GetDnszoneAccessRequestOptions) (result *AccessRequest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDnszoneAccessRequestOptions, "getDnszoneAccessRequestOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDnszoneAccessRequestOptions, "getDnszoneAccessRequestOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getDnszoneAccessRequestOptions.InstanceID,
		"dnszone_id":  *getDnszoneAccessRequestOptions.DnszoneID,
		"request_id":  *getDnszoneAccessRequestOptions.RequestID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/access_requests/{request_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDnszoneAccessRequestOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetDnszoneAccessRequest")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDnszoneAccessRequestOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getDnszoneAccessRequestOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_dnszone_access_request", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessRequest)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateDnszoneAccessRequest : Update an access request
// Update the state of an access request.
func (dnsSvcs *DnsSvcsV1) UpdateDnszoneAccessRequest(updateDnszoneAccessRequestOptions *UpdateDnszoneAccessRequestOptions) (result *AccessRequest, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.UpdateDnszoneAccessRequestWithContext(context.Background(), updateDnszoneAccessRequestOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateDnszoneAccessRequestWithContext is an alternate form of the UpdateDnszoneAccessRequest method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateDnszoneAccessRequestWithContext(ctx context.Context, updateDnszoneAccessRequestOptions *UpdateDnszoneAccessRequestOptions) (result *AccessRequest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDnszoneAccessRequestOptions, "updateDnszoneAccessRequestOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateDnszoneAccessRequestOptions, "updateDnszoneAccessRequestOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateDnszoneAccessRequestOptions.InstanceID,
		"dnszone_id":  *updateDnszoneAccessRequestOptions.DnszoneID,
		"request_id":  *updateDnszoneAccessRequestOptions.RequestID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/access_requests/{request_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateDnszoneAccessRequestOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "UpdateDnszoneAccessRequest")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateDnszoneAccessRequestOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateDnszoneAccessRequestOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateDnszoneAccessRequestOptions.Action != nil {
		body["action"] = updateDnszoneAccessRequestOptions.Action
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_dnszone_access_request", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessRequest)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListLinkedPermittedNetworks : List permitted networks
// List the permitted networks for a linked zone.
func (dnsSvcs *DnsSvcsV1) ListLinkedPermittedNetworks(listLinkedPermittedNetworksOptions *ListLinkedPermittedNetworksOptions) (result *ListPermittedNetworks, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.ListLinkedPermittedNetworksWithContext(context.Background(), listLinkedPermittedNetworksOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListLinkedPermittedNetworksWithContext is an alternate form of the ListLinkedPermittedNetworks method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListLinkedPermittedNetworksWithContext(ctx context.Context, listLinkedPermittedNetworksOptions *ListLinkedPermittedNetworksOptions) (result *ListPermittedNetworks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listLinkedPermittedNetworksOptions, "listLinkedPermittedNetworksOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listLinkedPermittedNetworksOptions, "listLinkedPermittedNetworksOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":       *listLinkedPermittedNetworksOptions.InstanceID,
		"linked_dnszone_id": *listLinkedPermittedNetworksOptions.LinkedDnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones/{linked_dnszone_id}/permitted_networks`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listLinkedPermittedNetworksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "ListLinkedPermittedNetworks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listLinkedPermittedNetworksOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listLinkedPermittedNetworksOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_linked_permitted_networks", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListPermittedNetworks)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateLzPermittedNetwork : Create a permitted network
// Create a permitted network for a linked zone.
func (dnsSvcs *DnsSvcsV1) CreateLzPermittedNetwork(createLzPermittedNetworkOptions *CreateLzPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.CreateLzPermittedNetworkWithContext(context.Background(), createLzPermittedNetworkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateLzPermittedNetworkWithContext is an alternate form of the CreateLzPermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateLzPermittedNetworkWithContext(ctx context.Context, createLzPermittedNetworkOptions *CreateLzPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createLzPermittedNetworkOptions, "createLzPermittedNetworkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createLzPermittedNetworkOptions, "createLzPermittedNetworkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":       *createLzPermittedNetworkOptions.InstanceID,
		"linked_dnszone_id": *createLzPermittedNetworkOptions.LinkedDnszoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones/{linked_dnszone_id}/permitted_networks`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createLzPermittedNetworkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "CreateLzPermittedNetwork")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createLzPermittedNetworkOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createLzPermittedNetworkOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createLzPermittedNetworkOptions.Type != nil {
		body["type"] = createLzPermittedNetworkOptions.Type
	}
	if createLzPermittedNetworkOptions.PermittedNetwork != nil {
		body["permitted_network"] = createLzPermittedNetworkOptions.PermittedNetwork
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_lz_permitted_network", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteLzPermittedNetwork : Remove a permitted network
// Remove a permitted network from a linked zone.
func (dnsSvcs *DnsSvcsV1) DeleteLzPermittedNetwork(deleteLzPermittedNetworkOptions *DeleteLzPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.DeleteLzPermittedNetworkWithContext(context.Background(), deleteLzPermittedNetworkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteLzPermittedNetworkWithContext is an alternate form of the DeleteLzPermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteLzPermittedNetworkWithContext(ctx context.Context, deleteLzPermittedNetworkOptions *DeleteLzPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLzPermittedNetworkOptions, "deleteLzPermittedNetworkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteLzPermittedNetworkOptions, "deleteLzPermittedNetworkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":          *deleteLzPermittedNetworkOptions.InstanceID,
		"linked_dnszone_id":    *deleteLzPermittedNetworkOptions.LinkedDnszoneID,
		"permitted_network_id": *deleteLzPermittedNetworkOptions.PermittedNetworkID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones/{linked_dnszone_id}/permitted_networks/{permitted_network_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteLzPermittedNetworkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "DeleteLzPermittedNetwork")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteLzPermittedNetworkOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteLzPermittedNetworkOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_lz_permitted_network", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetLinkedPermittedNetwork : Get a permitted network
// Get a permitted network of a linked zone.
func (dnsSvcs *DnsSvcsV1) GetLinkedPermittedNetwork(getLinkedPermittedNetworkOptions *GetLinkedPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	result, response, err = dnsSvcs.GetLinkedPermittedNetworkWithContext(context.Background(), getLinkedPermittedNetworkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLinkedPermittedNetworkWithContext is an alternate form of the GetLinkedPermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetLinkedPermittedNetworkWithContext(ctx context.Context, getLinkedPermittedNetworkOptions *GetLinkedPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLinkedPermittedNetworkOptions, "getLinkedPermittedNetworkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getLinkedPermittedNetworkOptions, "getLinkedPermittedNetworkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":          *getLinkedPermittedNetworkOptions.InstanceID,
		"linked_dnszone_id":    *getLinkedPermittedNetworkOptions.LinkedDnszoneID,
		"permitted_network_id": *getLinkedPermittedNetworkOptions.PermittedNetworkID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/linked_dnszones/{linked_dnszone_id}/permitted_networks/{permitted_network_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLinkedPermittedNetworkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_svcs", "V1", "GetLinkedPermittedNetwork")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getLinkedPermittedNetworkOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getLinkedPermittedNetworkOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_linked_permitted_network", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.0")
}

// AccessRequestRequestor : The information of requestor.
type AccessRequestRequestor struct {
	// The account ID of requestor.
	AccountID *string `json:"account_id,omitempty"`

	// The requestor's DNS service instance ID.
	InstanceID *string `json:"instance_id,omitempty"`

	// The requestor's linked zone ID.
	LinkedZoneID *string `json:"linked_zone_id,omitempty"`
}

// UnmarshalAccessRequestRequestor unmarshals an instance of AccessRequestRequestor from the specified map of raw messages.
func UnmarshalAccessRequestRequestor(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessRequestRequestor)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "linked_zone_id", &obj.LinkedZoneID)
	if err != nil {
		err = core.SDKErrorf(err, "", "linked_zone_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AddCustomResolverLocationOptions : The AddCustomResolverLocation options.
type AddCustomResolverLocationOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Custom resolver location, subnet CRN.
	SubnetCrn *string `json:"subnet_crn" validate:"required"`

	// Enable/Disable custom resolver location.
	Enabled *bool `json:"enabled,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewAddCustomResolverLocationOptions : Instantiate AddCustomResolverLocationOptions
func (*DnsSvcsV1) NewAddCustomResolverLocationOptions(instanceID string, resolverID string, subnetCrn string) *AddCustomResolverLocationOptions {
	return &AddCustomResolverLocationOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
		SubnetCrn:  core.StringPtr(subnetCrn),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *AddCustomResolverLocationOptions) SetInstanceID(instanceID string) *AddCustomResolverLocationOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *AddCustomResolverLocationOptions) SetResolverID(resolverID string) *AddCustomResolverLocationOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetSubnetCrn : Allow user to set SubnetCrn
func (_options *AddCustomResolverLocationOptions) SetSubnetCrn(subnetCrn string) *AddCustomResolverLocationOptions {
	_options.SubnetCrn = core.StringPtr(subnetCrn)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *AddCustomResolverLocationOptions) SetEnabled(enabled bool) *AddCustomResolverLocationOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *AddCustomResolverLocationOptions) SetXCorrelationID(xCorrelationID string) *AddCustomResolverLocationOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *AddCustomResolverLocationOptions) SetHeaders(param map[string]string) *AddCustomResolverLocationOptions {
	options.Headers = param
	return options
}

// AzPoolsItem : AzPoolsItem struct
type AzPoolsItem struct {
	// Availability zone.
	AvailabilityZone *string `json:"availability_zone,omitempty"`

	// List of load balancer pools.
	Pools []string `json:"pools,omitempty"`
}

// UnmarshalAzPoolsItem unmarshals an instance of AzPoolsItem from the specified map of raw messages.
func UnmarshalAzPoolsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AzPoolsItem)
	err = core.UnmarshalPrimitive(m, "availability_zone", &obj.AvailabilityZone)
	if err != nil {
		err = core.SDKErrorf(err, "", "availability_zone-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pools", &obj.Pools)
	if err != nil {
		err = core.SDKErrorf(err, "", "pools-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCustomResolverOptions : The CreateCustomResolver options.
type CreateCustomResolverOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Name of the custom resolver.
	Name *string `json:"name" validate:"required"`

	// Descriptive text of the custom resolver.
	Description *string `json:"description,omitempty"`

	// Locations on which the custom resolver will be running.
	Locations []LocationInput `json:"locations,omitempty"`

	// The profile name of a custom resolver.
	Profile *string `json:"profile,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateCustomResolverOptions.Profile property.
// The profile name of a custom resolver.
const (
	CreateCustomResolverOptions_Profile_Advanced  = "advanced"
	CreateCustomResolverOptions_Profile_Essential = "essential"
	CreateCustomResolverOptions_Profile_Premier   = "premier"
)

// NewCreateCustomResolverOptions : Instantiate CreateCustomResolverOptions
func (*DnsSvcsV1) NewCreateCustomResolverOptions(instanceID string, name string) *CreateCustomResolverOptions {
	return &CreateCustomResolverOptions{
		InstanceID: core.StringPtr(instanceID),
		Name:       core.StringPtr(name),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateCustomResolverOptions) SetInstanceID(instanceID string) *CreateCustomResolverOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateCustomResolverOptions) SetName(name string) *CreateCustomResolverOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateCustomResolverOptions) SetDescription(description string) *CreateCustomResolverOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLocations : Allow user to set Locations
func (_options *CreateCustomResolverOptions) SetLocations(locations []LocationInput) *CreateCustomResolverOptions {
	_options.Locations = locations
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *CreateCustomResolverOptions) SetProfile(profile string) *CreateCustomResolverOptions {
	_options.Profile = core.StringPtr(profile)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateCustomResolverOptions) SetXCorrelationID(xCorrelationID string) *CreateCustomResolverOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCustomResolverOptions) SetHeaders(param map[string]string) *CreateCustomResolverOptions {
	options.Headers = param
	return options
}

// CreateDnszoneOptions : The CreateDnszone options.
type CreateDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Name of DNS zone.
	Name *string `json:"name" validate:"required"`

	// The text describing the purpose of a DNS zone.
	Description *string `json:"description,omitempty"`

	// The label of a DNS zone.
	Label *string `json:"label,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateDnszoneOptions : Instantiate CreateDnszoneOptions
func (*DnsSvcsV1) NewCreateDnszoneOptions(instanceID string, name string) *CreateDnszoneOptions {
	return &CreateDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		Name:       core.StringPtr(name),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateDnszoneOptions) SetInstanceID(instanceID string) *CreateDnszoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateDnszoneOptions) SetName(name string) *CreateDnszoneOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateDnszoneOptions) SetDescription(description string) *CreateDnszoneOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *CreateDnszoneOptions) SetLabel(label string) *CreateDnszoneOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateDnszoneOptions) SetXCorrelationID(xCorrelationID string) *CreateDnszoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDnszoneOptions) SetHeaders(param map[string]string) *CreateDnszoneOptions {
	options.Headers = param
	return options
}

// CreateForwardingRuleOptions : The CreateForwardingRule options.
type CreateForwardingRuleOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Create a forwarding rule.
	ForwardingRuleInput ForwardingRuleInputIntf `json:"forwarding_rule_input" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateForwardingRuleOptions : Instantiate CreateForwardingRuleOptions
func (*DnsSvcsV1) NewCreateForwardingRuleOptions(instanceID string, resolverID string, forwardingRuleInput ForwardingRuleInputIntf) *CreateForwardingRuleOptions {
	return &CreateForwardingRuleOptions{
		InstanceID:          core.StringPtr(instanceID),
		ResolverID:          core.StringPtr(resolverID),
		ForwardingRuleInput: forwardingRuleInput,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateForwardingRuleOptions) SetInstanceID(instanceID string) *CreateForwardingRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *CreateForwardingRuleOptions) SetResolverID(resolverID string) *CreateForwardingRuleOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetForwardingRuleInput : Allow user to set ForwardingRuleInput
func (_options *CreateForwardingRuleOptions) SetForwardingRuleInput(forwardingRuleInput ForwardingRuleInputIntf) *CreateForwardingRuleOptions {
	_options.ForwardingRuleInput = forwardingRuleInput
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateForwardingRuleOptions) SetXCorrelationID(xCorrelationID string) *CreateForwardingRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateForwardingRuleOptions) SetHeaders(param map[string]string) *CreateForwardingRuleOptions {
	options.Headers = param
	return options
}

// CreateLinkedZoneOptions : The CreateLinkedZone options.
type CreateLinkedZoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Owner's instance ID.
	OwnerInstanceID *string `json:"owner_instance_id" validate:"required"`

	// Owner's DNS zone ID.
	OwnerZoneID *string `json:"owner_zone_id" validate:"required"`

	// Descriptive text of the linked zone.
	Description *string `json:"description,omitempty"`

	// The label of linked zone.
	Label *string `json:"label,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateLinkedZoneOptions : Instantiate CreateLinkedZoneOptions
func (*DnsSvcsV1) NewCreateLinkedZoneOptions(instanceID string, ownerInstanceID string, ownerZoneID string) *CreateLinkedZoneOptions {
	return &CreateLinkedZoneOptions{
		InstanceID:      core.StringPtr(instanceID),
		OwnerInstanceID: core.StringPtr(ownerInstanceID),
		OwnerZoneID:     core.StringPtr(ownerZoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateLinkedZoneOptions) SetInstanceID(instanceID string) *CreateLinkedZoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetOwnerInstanceID : Allow user to set OwnerInstanceID
func (_options *CreateLinkedZoneOptions) SetOwnerInstanceID(ownerInstanceID string) *CreateLinkedZoneOptions {
	_options.OwnerInstanceID = core.StringPtr(ownerInstanceID)
	return _options
}

// SetOwnerZoneID : Allow user to set OwnerZoneID
func (_options *CreateLinkedZoneOptions) SetOwnerZoneID(ownerZoneID string) *CreateLinkedZoneOptions {
	_options.OwnerZoneID = core.StringPtr(ownerZoneID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateLinkedZoneOptions) SetDescription(description string) *CreateLinkedZoneOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *CreateLinkedZoneOptions) SetLabel(label string) *CreateLinkedZoneOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateLinkedZoneOptions) SetXCorrelationID(xCorrelationID string) *CreateLinkedZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLinkedZoneOptions) SetHeaders(param map[string]string) *CreateLinkedZoneOptions {
	options.Headers = param
	return options
}

// CreateLoadBalancerOptions : The CreateLoadBalancer options.
type CreateLoadBalancerOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// Name of the load balancer.
	Name *string `json:"name" validate:"required"`

	// The pool ID to use when all other pools are detected as unhealthy.
	FallbackPool *string `json:"fallback_pool" validate:"required"`

	// A list of pool IDs ordered by their failover priority. Pools defined here are used by default, or when region_pools
	// are not configured for a given region.
	DefaultPools []string `json:"default_pools" validate:"required"`

	// Descriptive text of the load balancer.
	Description *string `json:"description,omitempty"`

	// Whether the load balancer is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Time to live in second.
	TTL *int64 `json:"ttl,omitempty"`

	// Map availability zones to pool IDs.
	AzPools []AzPoolsItem `json:"az_pools,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateLoadBalancerOptions : Instantiate CreateLoadBalancerOptions
func (*DnsSvcsV1) NewCreateLoadBalancerOptions(instanceID string, dnszoneID string, name string, fallbackPool string, defaultPools []string) *CreateLoadBalancerOptions {
	return &CreateLoadBalancerOptions{
		InstanceID:   core.StringPtr(instanceID),
		DnszoneID:    core.StringPtr(dnszoneID),
		Name:         core.StringPtr(name),
		FallbackPool: core.StringPtr(fallbackPool),
		DefaultPools: defaultPools,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateLoadBalancerOptions) SetInstanceID(instanceID string) *CreateLoadBalancerOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *CreateLoadBalancerOptions) SetDnszoneID(dnszoneID string) *CreateLoadBalancerOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateLoadBalancerOptions) SetName(name string) *CreateLoadBalancerOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetFallbackPool : Allow user to set FallbackPool
func (_options *CreateLoadBalancerOptions) SetFallbackPool(fallbackPool string) *CreateLoadBalancerOptions {
	_options.FallbackPool = core.StringPtr(fallbackPool)
	return _options
}

// SetDefaultPools : Allow user to set DefaultPools
func (_options *CreateLoadBalancerOptions) SetDefaultPools(defaultPools []string) *CreateLoadBalancerOptions {
	_options.DefaultPools = defaultPools
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateLoadBalancerOptions) SetDescription(description string) *CreateLoadBalancerOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateLoadBalancerOptions) SetEnabled(enabled bool) *CreateLoadBalancerOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetTTL : Allow user to set TTL
func (_options *CreateLoadBalancerOptions) SetTTL(ttl int64) *CreateLoadBalancerOptions {
	_options.TTL = core.Int64Ptr(ttl)
	return _options
}

// SetAzPools : Allow user to set AzPools
func (_options *CreateLoadBalancerOptions) SetAzPools(azPools []AzPoolsItem) *CreateLoadBalancerOptions {
	_options.AzPools = azPools
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateLoadBalancerOptions) SetXCorrelationID(xCorrelationID string) *CreateLoadBalancerOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLoadBalancerOptions) SetHeaders(param map[string]string) *CreateLoadBalancerOptions {
	options.Headers = param
	return options
}

// CreateLzPermittedNetworkOptions : The CreateLzPermittedNetwork options.
type CreateLzPermittedNetworkOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a linked zone.
	LinkedDnszoneID *string `json:"linked_dnszone_id" validate:"required,ne="`

	// The type of a permitted network.
	Type *string `json:"type" validate:"required"`

	// Permitted network data for VPC.
	PermittedNetwork *PermittedNetworkVpc `json:"permitted_network" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateLzPermittedNetworkOptions.Type property.
// The type of a permitted network.
const (
	CreateLzPermittedNetworkOptions_Type_Vpc = "vpc"
)

// NewCreateLzPermittedNetworkOptions : Instantiate CreateLzPermittedNetworkOptions
func (*DnsSvcsV1) NewCreateLzPermittedNetworkOptions(instanceID string, linkedDnszoneID string, typeVar string, permittedNetwork *PermittedNetworkVpc) *CreateLzPermittedNetworkOptions {
	return &CreateLzPermittedNetworkOptions{
		InstanceID:       core.StringPtr(instanceID),
		LinkedDnszoneID:  core.StringPtr(linkedDnszoneID),
		Type:             core.StringPtr(typeVar),
		PermittedNetwork: permittedNetwork,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateLzPermittedNetworkOptions) SetInstanceID(instanceID string) *CreateLzPermittedNetworkOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLinkedDnszoneID : Allow user to set LinkedDnszoneID
func (_options *CreateLzPermittedNetworkOptions) SetLinkedDnszoneID(linkedDnszoneID string) *CreateLzPermittedNetworkOptions {
	_options.LinkedDnszoneID = core.StringPtr(linkedDnszoneID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateLzPermittedNetworkOptions) SetType(typeVar string) *CreateLzPermittedNetworkOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetPermittedNetwork : Allow user to set PermittedNetwork
func (_options *CreateLzPermittedNetworkOptions) SetPermittedNetwork(permittedNetwork *PermittedNetworkVpc) *CreateLzPermittedNetworkOptions {
	_options.PermittedNetwork = permittedNetwork
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateLzPermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *CreateLzPermittedNetworkOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLzPermittedNetworkOptions) SetHeaders(param map[string]string) *CreateLzPermittedNetworkOptions {
	options.Headers = param
	return options
}

// CreateMonitorOptions : The CreateMonitor options.
type CreateMonitorOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The name of the load balancer monitor.
	Name *string `json:"name" validate:"required"`

	// The protocol to use for the health check. Currently supported protocols are 'HTTP','HTTPS' and 'TCP'.
	Type *string `json:"type" validate:"required"`

	// Descriptive text of the load balancer monitor.
	Description *string `json:"description,omitempty"`

	// Port number to connect to for the health check. Required for TCP checks. HTTP and HTTPS checks should only define
	// the port when using a non-standard port (HTTP: default 80, HTTPS: default 443).
	Port *int64 `json:"port,omitempty"`

	// The interval between each health check. Shorter intervals may improve failover time, but will increase load on the
	// origins as we check from multiple locations.
	Interval *int64 `json:"interval,omitempty"`

	// The number of retries to attempt in case of a timeout before marking the origin as unhealthy. Retries are attempted
	// immediately.
	Retries *int64 `json:"retries,omitempty"`

	// The timeout (in seconds) before marking the health check as failed.
	Timeout *int64 `json:"timeout,omitempty"`

	// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
	Method *string `json:"method,omitempty"`

	// The endpoint path to health check against. This parameter is only valid for HTTP and HTTPS monitors.
	Path *string `json:"path,omitempty"`

	// The HTTP request headers to send in the health check. It is recommended you set a Host header by default. The
	// User-Agent header cannot be overridden. This parameter is only valid for HTTP and HTTPS monitors.
	HeadersVar []HealthcheckHeader `json:"headers,omitempty"`

	// Do not validate the certificate when monitor use HTTPS. This parameter is currently only valid for HTTPS monitors.
	AllowInsecure *bool `json:"allow_insecure,omitempty"`

	// The expected HTTP response code or code range of the health check. This parameter is only valid for HTTP and HTTPS
	// monitors.
	ExpectedCodes *string `json:"expected_codes,omitempty"`

	// A case-insensitive sub-string to look for in the response body. If this string is not found, the origin will be
	// marked as unhealthy. This parameter is only valid for HTTP and HTTPS monitors.
	ExpectedBody *string `json:"expected_body,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateMonitorOptions.Type property.
// The protocol to use for the health check. Currently supported protocols are 'HTTP','HTTPS' and 'TCP'.
const (
	CreateMonitorOptions_Type_Http  = "HTTP"
	CreateMonitorOptions_Type_Https = "HTTPS"
	CreateMonitorOptions_Type_Tcp   = "TCP"
)

// Constants associated with the CreateMonitorOptions.Method property.
// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
const (
	CreateMonitorOptions_Method_Get  = "GET"
	CreateMonitorOptions_Method_Head = "HEAD"
)

// NewCreateMonitorOptions : Instantiate CreateMonitorOptions
func (*DnsSvcsV1) NewCreateMonitorOptions(instanceID string, name string, typeVar string) *CreateMonitorOptions {
	return &CreateMonitorOptions{
		InstanceID: core.StringPtr(instanceID),
		Name:       core.StringPtr(name),
		Type:       core.StringPtr(typeVar),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateMonitorOptions) SetInstanceID(instanceID string) *CreateMonitorOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateMonitorOptions) SetName(name string) *CreateMonitorOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateMonitorOptions) SetType(typeVar string) *CreateMonitorOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateMonitorOptions) SetDescription(description string) *CreateMonitorOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetPort : Allow user to set Port
func (_options *CreateMonitorOptions) SetPort(port int64) *CreateMonitorOptions {
	_options.Port = core.Int64Ptr(port)
	return _options
}

// SetInterval : Allow user to set Interval
func (_options *CreateMonitorOptions) SetInterval(interval int64) *CreateMonitorOptions {
	_options.Interval = core.Int64Ptr(interval)
	return _options
}

// SetRetries : Allow user to set Retries
func (_options *CreateMonitorOptions) SetRetries(retries int64) *CreateMonitorOptions {
	_options.Retries = core.Int64Ptr(retries)
	return _options
}

// SetTimeout : Allow user to set Timeout
func (_options *CreateMonitorOptions) SetTimeout(timeout int64) *CreateMonitorOptions {
	_options.Timeout = core.Int64Ptr(timeout)
	return _options
}

// SetMethod : Allow user to set Method
func (_options *CreateMonitorOptions) SetMethod(method string) *CreateMonitorOptions {
	_options.Method = core.StringPtr(method)
	return _options
}

// SetPath : Allow user to set Path
func (_options *CreateMonitorOptions) SetPath(path string) *CreateMonitorOptions {
	_options.Path = core.StringPtr(path)
	return _options
}

// SetHeadersVar : Allow user to set HeadersVar
func (_options *CreateMonitorOptions) SetHeadersVar(headersVar []HealthcheckHeader) *CreateMonitorOptions {
	_options.HeadersVar = headersVar
	return _options
}

// SetAllowInsecure : Allow user to set AllowInsecure
func (_options *CreateMonitorOptions) SetAllowInsecure(allowInsecure bool) *CreateMonitorOptions {
	_options.AllowInsecure = core.BoolPtr(allowInsecure)
	return _options
}

// SetExpectedCodes : Allow user to set ExpectedCodes
func (_options *CreateMonitorOptions) SetExpectedCodes(expectedCodes string) *CreateMonitorOptions {
	_options.ExpectedCodes = core.StringPtr(expectedCodes)
	return _options
}

// SetExpectedBody : Allow user to set ExpectedBody
func (_options *CreateMonitorOptions) SetExpectedBody(expectedBody string) *CreateMonitorOptions {
	_options.ExpectedBody = core.StringPtr(expectedBody)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateMonitorOptions) SetXCorrelationID(xCorrelationID string) *CreateMonitorOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateMonitorOptions) SetHeaders(param map[string]string) *CreateMonitorOptions {
	options.Headers = param
	return options
}

// CreatePermittedNetworkOptions : The CreatePermittedNetwork options.
type CreatePermittedNetworkOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The type of a permitted network.
	Type *string `json:"type" validate:"required"`

	// Permitted network data for VPC.
	PermittedNetwork *PermittedNetworkVpc `json:"permitted_network" validate:"required"`

	// The account identifiers of the owner zone and linked zones in the format of "?account=account1,account2,account3".
	// Maximum 5 accounts are allowed.
	Accounts *string `json:"accounts,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreatePermittedNetworkOptions.Type property.
// The type of a permitted network.
const (
	CreatePermittedNetworkOptions_Type_Vpc = "vpc"
)

// NewCreatePermittedNetworkOptions : Instantiate CreatePermittedNetworkOptions
func (*DnsSvcsV1) NewCreatePermittedNetworkOptions(instanceID string, dnszoneID string, typeVar string, permittedNetwork *PermittedNetworkVpc) *CreatePermittedNetworkOptions {
	return &CreatePermittedNetworkOptions{
		InstanceID:       core.StringPtr(instanceID),
		DnszoneID:        core.StringPtr(dnszoneID),
		Type:             core.StringPtr(typeVar),
		PermittedNetwork: permittedNetwork,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreatePermittedNetworkOptions) SetInstanceID(instanceID string) *CreatePermittedNetworkOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *CreatePermittedNetworkOptions) SetDnszoneID(dnszoneID string) *CreatePermittedNetworkOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreatePermittedNetworkOptions) SetType(typeVar string) *CreatePermittedNetworkOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetPermittedNetwork : Allow user to set PermittedNetwork
func (_options *CreatePermittedNetworkOptions) SetPermittedNetwork(permittedNetwork *PermittedNetworkVpc) *CreatePermittedNetworkOptions {
	_options.PermittedNetwork = permittedNetwork
	return _options
}

// SetAccounts : Allow user to set Accounts
func (_options *CreatePermittedNetworkOptions) SetAccounts(accounts string) *CreatePermittedNetworkOptions {
	_options.Accounts = core.StringPtr(accounts)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreatePermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *CreatePermittedNetworkOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePermittedNetworkOptions) SetHeaders(param map[string]string) *CreatePermittedNetworkOptions {
	options.Headers = param
	return options
}

// CreatePoolOptions : The CreatePool options.
type CreatePoolOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Name of the load balancer pool.
	Name *string `json:"name" validate:"required"`

	// The list of origins within this pool. Traffic directed at this pool is balanced across all currently healthy
	// origins, provided the pool itself is healthy.
	Origins []OriginInput `json:"origins" validate:"required"`

	// Descriptive text of the load balancer pool.
	Description *string `json:"description,omitempty"`

	// Whether the load balancer pool is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The minimum number of origins that must be healthy for this pool to serve traffic. If the number of healthy origins
	// falls below this number, the pool will be marked unhealthy and we will failover to the next available pool.
	HealthyOriginsThreshold *int64 `json:"healthy_origins_threshold,omitempty"`

	// The ID of the load balancer monitor to be associated to this pool.
	Monitor *string `json:"monitor,omitempty"`

	// The notification channel.
	NotificationChannel *string `json:"notification_channel,omitempty"`

	// Health check region of VSIs.
	HealthcheckRegion *string `json:"healthcheck_region,omitempty"`

	// Health check subnet CRN.
	HealthcheckSubnets []string `json:"healthcheck_subnets,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreatePoolOptions.HealthcheckRegion property.
// Health check region of VSIs.
const (
	CreatePoolOptions_HealthcheckRegion_AuSyd   = "au-syd"
	CreatePoolOptions_HealthcheckRegion_EuDu    = "eu-du"
	CreatePoolOptions_HealthcheckRegion_EuGb    = "eu-gb"
	CreatePoolOptions_HealthcheckRegion_JpTok   = "jp-tok"
	CreatePoolOptions_HealthcheckRegion_UsEast  = "us-east"
	CreatePoolOptions_HealthcheckRegion_UsSouth = "us-south"
)

// NewCreatePoolOptions : Instantiate CreatePoolOptions
func (*DnsSvcsV1) NewCreatePoolOptions(instanceID string, name string, origins []OriginInput) *CreatePoolOptions {
	return &CreatePoolOptions{
		InstanceID: core.StringPtr(instanceID),
		Name:       core.StringPtr(name),
		Origins:    origins,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreatePoolOptions) SetInstanceID(instanceID string) *CreatePoolOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreatePoolOptions) SetName(name string) *CreatePoolOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetOrigins : Allow user to set Origins
func (_options *CreatePoolOptions) SetOrigins(origins []OriginInput) *CreatePoolOptions {
	_options.Origins = origins
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreatePoolOptions) SetDescription(description string) *CreatePoolOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreatePoolOptions) SetEnabled(enabled bool) *CreatePoolOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetHealthyOriginsThreshold : Allow user to set HealthyOriginsThreshold
func (_options *CreatePoolOptions) SetHealthyOriginsThreshold(healthyOriginsThreshold int64) *CreatePoolOptions {
	_options.HealthyOriginsThreshold = core.Int64Ptr(healthyOriginsThreshold)
	return _options
}

// SetMonitor : Allow user to set Monitor
func (_options *CreatePoolOptions) SetMonitor(monitor string) *CreatePoolOptions {
	_options.Monitor = core.StringPtr(monitor)
	return _options
}

// SetNotificationChannel : Allow user to set NotificationChannel
func (_options *CreatePoolOptions) SetNotificationChannel(notificationChannel string) *CreatePoolOptions {
	_options.NotificationChannel = core.StringPtr(notificationChannel)
	return _options
}

// SetHealthcheckRegion : Allow user to set HealthcheckRegion
func (_options *CreatePoolOptions) SetHealthcheckRegion(healthcheckRegion string) *CreatePoolOptions {
	_options.HealthcheckRegion = core.StringPtr(healthcheckRegion)
	return _options
}

// SetHealthcheckSubnets : Allow user to set HealthcheckSubnets
func (_options *CreatePoolOptions) SetHealthcheckSubnets(healthcheckSubnets []string) *CreatePoolOptions {
	_options.HealthcheckSubnets = healthcheckSubnets
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreatePoolOptions) SetXCorrelationID(xCorrelationID string) *CreatePoolOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePoolOptions) SetHeaders(param map[string]string) *CreatePoolOptions {
	options.Headers = param
	return options
}

// CreateResourceRecordOptions : The CreateResourceRecord options.
type CreateResourceRecordOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// Type of the resource record.
	Type *string `json:"type" validate:"required"`

	// Name of the resource record.
	Name *string `json:"name,omitempty"`

	// Content of the resource record.
	Rdata ResourceRecordInputRdataIntf `json:"rdata,omitempty"`

	// Time to live in second.
	TTL *int64 `json:"ttl,omitempty"`

	// Only used for SRV record.
	Service *string `json:"service,omitempty"`

	// Only used for SRV record.
	Protocol *string `json:"protocol,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateResourceRecordOptions.Type property.
// Type of the resource record.
const (
	CreateResourceRecordOptions_Type_A     = "A"
	CreateResourceRecordOptions_Type_Aaaa  = "AAAA"
	CreateResourceRecordOptions_Type_Cname = "CNAME"
	CreateResourceRecordOptions_Type_Mx    = "MX"
	CreateResourceRecordOptions_Type_Ptr   = "PTR"
	CreateResourceRecordOptions_Type_Srv   = "SRV"
	CreateResourceRecordOptions_Type_Txt   = "TXT"
)

// NewCreateResourceRecordOptions : Instantiate CreateResourceRecordOptions
func (*DnsSvcsV1) NewCreateResourceRecordOptions(instanceID string, dnszoneID string, typeVar string) *CreateResourceRecordOptions {
	return &CreateResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		Type:       core.StringPtr(typeVar),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateResourceRecordOptions) SetInstanceID(instanceID string) *CreateResourceRecordOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *CreateResourceRecordOptions) SetDnszoneID(dnszoneID string) *CreateResourceRecordOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateResourceRecordOptions) SetType(typeVar string) *CreateResourceRecordOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateResourceRecordOptions) SetName(name string) *CreateResourceRecordOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetRdata : Allow user to set Rdata
func (_options *CreateResourceRecordOptions) SetRdata(rdata ResourceRecordInputRdataIntf) *CreateResourceRecordOptions {
	_options.Rdata = rdata
	return _options
}

// SetTTL : Allow user to set TTL
func (_options *CreateResourceRecordOptions) SetTTL(ttl int64) *CreateResourceRecordOptions {
	_options.TTL = core.Int64Ptr(ttl)
	return _options
}

// SetService : Allow user to set Service
func (_options *CreateResourceRecordOptions) SetService(service string) *CreateResourceRecordOptions {
	_options.Service = core.StringPtr(service)
	return _options
}

// SetProtocol : Allow user to set Protocol
func (_options *CreateResourceRecordOptions) SetProtocol(protocol string) *CreateResourceRecordOptions {
	_options.Protocol = core.StringPtr(protocol)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateResourceRecordOptions) SetXCorrelationID(xCorrelationID string) *CreateResourceRecordOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateResourceRecordOptions) SetHeaders(param map[string]string) *CreateResourceRecordOptions {
	options.Headers = param
	return options
}

// CreateSecondaryZoneOptions : The CreateSecondaryZone options.
type CreateSecondaryZoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// zone name.
	Zone *string `json:"zone" validate:"required"`

	// The addresses of DNS servers where the secondary zone data should be transferred from.
	TransferFrom []string `json:"transfer_from" validate:"required"`

	// Descriptive text of the secondary zone.
	Description *string `json:"description,omitempty"`

	// Enable/Disable the secondary zone.
	Enabled *bool `json:"enabled,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateSecondaryZoneOptions : Instantiate CreateSecondaryZoneOptions
func (*DnsSvcsV1) NewCreateSecondaryZoneOptions(instanceID string, resolverID string, zone string, transferFrom []string) *CreateSecondaryZoneOptions {
	return &CreateSecondaryZoneOptions{
		InstanceID:   core.StringPtr(instanceID),
		ResolverID:   core.StringPtr(resolverID),
		Zone:         core.StringPtr(zone),
		TransferFrom: transferFrom,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateSecondaryZoneOptions) SetInstanceID(instanceID string) *CreateSecondaryZoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *CreateSecondaryZoneOptions) SetResolverID(resolverID string) *CreateSecondaryZoneOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetZone : Allow user to set Zone
func (_options *CreateSecondaryZoneOptions) SetZone(zone string) *CreateSecondaryZoneOptions {
	_options.Zone = core.StringPtr(zone)
	return _options
}

// SetTransferFrom : Allow user to set TransferFrom
func (_options *CreateSecondaryZoneOptions) SetTransferFrom(transferFrom []string) *CreateSecondaryZoneOptions {
	_options.TransferFrom = transferFrom
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateSecondaryZoneOptions) SetDescription(description string) *CreateSecondaryZoneOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateSecondaryZoneOptions) SetEnabled(enabled bool) *CreateSecondaryZoneOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateSecondaryZoneOptions) SetXCorrelationID(xCorrelationID string) *CreateSecondaryZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecondaryZoneOptions) SetHeaders(param map[string]string) *CreateSecondaryZoneOptions {
	options.Headers = param
	return options
}

// DeleteCustomResolverLocationOptions : The DeleteCustomResolverLocation options.
type DeleteCustomResolverLocationOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Custom resolver location ID.
	LocationID *string `json:"location_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteCustomResolverLocationOptions : Instantiate DeleteCustomResolverLocationOptions
func (*DnsSvcsV1) NewDeleteCustomResolverLocationOptions(instanceID string, resolverID string, locationID string) *DeleteCustomResolverLocationOptions {
	return &DeleteCustomResolverLocationOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
		LocationID: core.StringPtr(locationID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteCustomResolverLocationOptions) SetInstanceID(instanceID string) *DeleteCustomResolverLocationOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *DeleteCustomResolverLocationOptions) SetResolverID(resolverID string) *DeleteCustomResolverLocationOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetLocationID : Allow user to set LocationID
func (_options *DeleteCustomResolverLocationOptions) SetLocationID(locationID string) *DeleteCustomResolverLocationOptions {
	_options.LocationID = core.StringPtr(locationID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteCustomResolverLocationOptions) SetXCorrelationID(xCorrelationID string) *DeleteCustomResolverLocationOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCustomResolverLocationOptions) SetHeaders(param map[string]string) *DeleteCustomResolverLocationOptions {
	options.Headers = param
	return options
}

// DeleteCustomResolverOptions : The DeleteCustomResolver options.
type DeleteCustomResolverOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteCustomResolverOptions : Instantiate DeleteCustomResolverOptions
func (*DnsSvcsV1) NewDeleteCustomResolverOptions(instanceID string, resolverID string) *DeleteCustomResolverOptions {
	return &DeleteCustomResolverOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteCustomResolverOptions) SetInstanceID(instanceID string) *DeleteCustomResolverOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *DeleteCustomResolverOptions) SetResolverID(resolverID string) *DeleteCustomResolverOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteCustomResolverOptions) SetXCorrelationID(xCorrelationID string) *DeleteCustomResolverOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCustomResolverOptions) SetHeaders(param map[string]string) *DeleteCustomResolverOptions {
	options.Headers = param
	return options
}

// DeleteDnszoneOptions : The DeleteDnszone options.
type DeleteDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteDnszoneOptions : Instantiate DeleteDnszoneOptions
func (*DnsSvcsV1) NewDeleteDnszoneOptions(instanceID string, dnszoneID string) *DeleteDnszoneOptions {
	return &DeleteDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteDnszoneOptions) SetInstanceID(instanceID string) *DeleteDnszoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *DeleteDnszoneOptions) SetDnszoneID(dnszoneID string) *DeleteDnszoneOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteDnszoneOptions) SetXCorrelationID(xCorrelationID string) *DeleteDnszoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDnszoneOptions) SetHeaders(param map[string]string) *DeleteDnszoneOptions {
	options.Headers = param
	return options
}

// DeleteForwardingRuleOptions : The DeleteForwardingRule options.
type DeleteForwardingRuleOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// The unique identifier of a rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteForwardingRuleOptions : Instantiate DeleteForwardingRuleOptions
func (*DnsSvcsV1) NewDeleteForwardingRuleOptions(instanceID string, resolverID string, ruleID string) *DeleteForwardingRuleOptions {
	return &DeleteForwardingRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteForwardingRuleOptions) SetInstanceID(instanceID string) *DeleteForwardingRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *DeleteForwardingRuleOptions) SetResolverID(resolverID string) *DeleteForwardingRuleOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteForwardingRuleOptions) SetRuleID(ruleID string) *DeleteForwardingRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteForwardingRuleOptions) SetXCorrelationID(xCorrelationID string) *DeleteForwardingRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteForwardingRuleOptions) SetHeaders(param map[string]string) *DeleteForwardingRuleOptions {
	options.Headers = param
	return options
}

// DeleteLinkedZoneOptions : The DeleteLinkedZone options.
type DeleteLinkedZoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a linked zone.
	LinkedDnszoneID *string `json:"linked_dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteLinkedZoneOptions : Instantiate DeleteLinkedZoneOptions
func (*DnsSvcsV1) NewDeleteLinkedZoneOptions(instanceID string, linkedDnszoneID string) *DeleteLinkedZoneOptions {
	return &DeleteLinkedZoneOptions{
		InstanceID:      core.StringPtr(instanceID),
		LinkedDnszoneID: core.StringPtr(linkedDnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteLinkedZoneOptions) SetInstanceID(instanceID string) *DeleteLinkedZoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLinkedDnszoneID : Allow user to set LinkedDnszoneID
func (_options *DeleteLinkedZoneOptions) SetLinkedDnszoneID(linkedDnszoneID string) *DeleteLinkedZoneOptions {
	_options.LinkedDnszoneID = core.StringPtr(linkedDnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteLinkedZoneOptions) SetXCorrelationID(xCorrelationID string) *DeleteLinkedZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLinkedZoneOptions) SetHeaders(param map[string]string) *DeleteLinkedZoneOptions {
	options.Headers = param
	return options
}

// DeleteLoadBalancerOptions : The DeleteLoadBalancer options.
type DeleteLoadBalancerOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of a load balancer.
	LbID *string `json:"lb_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteLoadBalancerOptions : Instantiate DeleteLoadBalancerOptions
func (*DnsSvcsV1) NewDeleteLoadBalancerOptions(instanceID string, dnszoneID string, lbID string) *DeleteLoadBalancerOptions {
	return &DeleteLoadBalancerOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		LbID:       core.StringPtr(lbID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteLoadBalancerOptions) SetInstanceID(instanceID string) *DeleteLoadBalancerOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *DeleteLoadBalancerOptions) SetDnszoneID(dnszoneID string) *DeleteLoadBalancerOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetLbID : Allow user to set LbID
func (_options *DeleteLoadBalancerOptions) SetLbID(lbID string) *DeleteLoadBalancerOptions {
	_options.LbID = core.StringPtr(lbID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteLoadBalancerOptions) SetXCorrelationID(xCorrelationID string) *DeleteLoadBalancerOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLoadBalancerOptions) SetHeaders(param map[string]string) *DeleteLoadBalancerOptions {
	options.Headers = param
	return options
}

// DeleteLzPermittedNetworkOptions : The DeleteLzPermittedNetwork options.
type DeleteLzPermittedNetworkOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a linked zone.
	LinkedDnszoneID *string `json:"linked_dnszone_id" validate:"required,ne="`

	// The unique identifier of a permitted network.
	PermittedNetworkID *string `json:"permitted_network_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteLzPermittedNetworkOptions : Instantiate DeleteLzPermittedNetworkOptions
func (*DnsSvcsV1) NewDeleteLzPermittedNetworkOptions(instanceID string, linkedDnszoneID string, permittedNetworkID string) *DeleteLzPermittedNetworkOptions {
	return &DeleteLzPermittedNetworkOptions{
		InstanceID:         core.StringPtr(instanceID),
		LinkedDnszoneID:    core.StringPtr(linkedDnszoneID),
		PermittedNetworkID: core.StringPtr(permittedNetworkID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteLzPermittedNetworkOptions) SetInstanceID(instanceID string) *DeleteLzPermittedNetworkOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLinkedDnszoneID : Allow user to set LinkedDnszoneID
func (_options *DeleteLzPermittedNetworkOptions) SetLinkedDnszoneID(linkedDnszoneID string) *DeleteLzPermittedNetworkOptions {
	_options.LinkedDnszoneID = core.StringPtr(linkedDnszoneID)
	return _options
}

// SetPermittedNetworkID : Allow user to set PermittedNetworkID
func (_options *DeleteLzPermittedNetworkOptions) SetPermittedNetworkID(permittedNetworkID string) *DeleteLzPermittedNetworkOptions {
	_options.PermittedNetworkID = core.StringPtr(permittedNetworkID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteLzPermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *DeleteLzPermittedNetworkOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLzPermittedNetworkOptions) SetHeaders(param map[string]string) *DeleteLzPermittedNetworkOptions {
	options.Headers = param
	return options
}

// DeleteMonitorOptions : The DeleteMonitor options.
type DeleteMonitorOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a load balancer monitor.
	MonitorID *string `json:"monitor_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteMonitorOptions : Instantiate DeleteMonitorOptions
func (*DnsSvcsV1) NewDeleteMonitorOptions(instanceID string, monitorID string) *DeleteMonitorOptions {
	return &DeleteMonitorOptions{
		InstanceID: core.StringPtr(instanceID),
		MonitorID:  core.StringPtr(monitorID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteMonitorOptions) SetInstanceID(instanceID string) *DeleteMonitorOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetMonitorID : Allow user to set MonitorID
func (_options *DeleteMonitorOptions) SetMonitorID(monitorID string) *DeleteMonitorOptions {
	_options.MonitorID = core.StringPtr(monitorID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteMonitorOptions) SetXCorrelationID(xCorrelationID string) *DeleteMonitorOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteMonitorOptions) SetHeaders(param map[string]string) *DeleteMonitorOptions {
	options.Headers = param
	return options
}

// DeletePermittedNetworkOptions : The DeletePermittedNetwork options.
type DeletePermittedNetworkOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of a permitted network.
	PermittedNetworkID *string `json:"permitted_network_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeletePermittedNetworkOptions : Instantiate DeletePermittedNetworkOptions
func (*DnsSvcsV1) NewDeletePermittedNetworkOptions(instanceID string, dnszoneID string, permittedNetworkID string) *DeletePermittedNetworkOptions {
	return &DeletePermittedNetworkOptions{
		InstanceID:         core.StringPtr(instanceID),
		DnszoneID:          core.StringPtr(dnszoneID),
		PermittedNetworkID: core.StringPtr(permittedNetworkID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeletePermittedNetworkOptions) SetInstanceID(instanceID string) *DeletePermittedNetworkOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *DeletePermittedNetworkOptions) SetDnszoneID(dnszoneID string) *DeletePermittedNetworkOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetPermittedNetworkID : Allow user to set PermittedNetworkID
func (_options *DeletePermittedNetworkOptions) SetPermittedNetworkID(permittedNetworkID string) *DeletePermittedNetworkOptions {
	_options.PermittedNetworkID = core.StringPtr(permittedNetworkID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeletePermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *DeletePermittedNetworkOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeletePermittedNetworkOptions) SetHeaders(param map[string]string) *DeletePermittedNetworkOptions {
	options.Headers = param
	return options
}

// DeletePoolOptions : The DeletePool options.
type DeletePoolOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a load balancer pool.
	PoolID *string `json:"pool_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeletePoolOptions : Instantiate DeletePoolOptions
func (*DnsSvcsV1) NewDeletePoolOptions(instanceID string, poolID string) *DeletePoolOptions {
	return &DeletePoolOptions{
		InstanceID: core.StringPtr(instanceID),
		PoolID:     core.StringPtr(poolID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeletePoolOptions) SetInstanceID(instanceID string) *DeletePoolOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetPoolID : Allow user to set PoolID
func (_options *DeletePoolOptions) SetPoolID(poolID string) *DeletePoolOptions {
	_options.PoolID = core.StringPtr(poolID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeletePoolOptions) SetXCorrelationID(xCorrelationID string) *DeletePoolOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeletePoolOptions) SetHeaders(param map[string]string) *DeletePoolOptions {
	options.Headers = param
	return options
}

// DeleteResourceRecordOptions : The DeleteResourceRecord options.
type DeleteResourceRecordOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of a resource record.
	RecordID *string `json:"record_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteResourceRecordOptions : Instantiate DeleteResourceRecordOptions
func (*DnsSvcsV1) NewDeleteResourceRecordOptions(instanceID string, dnszoneID string, recordID string) *DeleteResourceRecordOptions {
	return &DeleteResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		RecordID:   core.StringPtr(recordID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteResourceRecordOptions) SetInstanceID(instanceID string) *DeleteResourceRecordOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *DeleteResourceRecordOptions) SetDnszoneID(dnszoneID string) *DeleteResourceRecordOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetRecordID : Allow user to set RecordID
func (_options *DeleteResourceRecordOptions) SetRecordID(recordID string) *DeleteResourceRecordOptions {
	_options.RecordID = core.StringPtr(recordID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteResourceRecordOptions) SetXCorrelationID(xCorrelationID string) *DeleteResourceRecordOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteResourceRecordOptions) SetHeaders(param map[string]string) *DeleteResourceRecordOptions {
	options.Headers = param
	return options
}

// DeleteSecondaryZoneOptions : The DeleteSecondaryZone options.
type DeleteSecondaryZoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// The unique identifier of a secondary zone.
	SecondaryZoneID *string `json:"secondary_zone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteSecondaryZoneOptions : Instantiate DeleteSecondaryZoneOptions
func (*DnsSvcsV1) NewDeleteSecondaryZoneOptions(instanceID string, resolverID string, secondaryZoneID string) *DeleteSecondaryZoneOptions {
	return &DeleteSecondaryZoneOptions{
		InstanceID:      core.StringPtr(instanceID),
		ResolverID:      core.StringPtr(resolverID),
		SecondaryZoneID: core.StringPtr(secondaryZoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteSecondaryZoneOptions) SetInstanceID(instanceID string) *DeleteSecondaryZoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *DeleteSecondaryZoneOptions) SetResolverID(resolverID string) *DeleteSecondaryZoneOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetSecondaryZoneID : Allow user to set SecondaryZoneID
func (_options *DeleteSecondaryZoneOptions) SetSecondaryZoneID(secondaryZoneID string) *DeleteSecondaryZoneOptions {
	_options.SecondaryZoneID = core.StringPtr(secondaryZoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteSecondaryZoneOptions) SetXCorrelationID(xCorrelationID string) *DeleteSecondaryZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecondaryZoneOptions) SetHeaders(param map[string]string) *DeleteSecondaryZoneOptions {
	options.Headers = param
	return options
}

// ExportResourceRecordsOptions : The ExportResourceRecords options.
type ExportResourceRecordsOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewExportResourceRecordsOptions : Instantiate ExportResourceRecordsOptions
func (*DnsSvcsV1) NewExportResourceRecordsOptions(instanceID string, dnszoneID string) *ExportResourceRecordsOptions {
	return &ExportResourceRecordsOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ExportResourceRecordsOptions) SetInstanceID(instanceID string) *ExportResourceRecordsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *ExportResourceRecordsOptions) SetDnszoneID(dnszoneID string) *ExportResourceRecordsOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ExportResourceRecordsOptions) SetXCorrelationID(xCorrelationID string) *ExportResourceRecordsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ExportResourceRecordsOptions) SetHeaders(param map[string]string) *ExportResourceRecordsOptions {
	options.Headers = param
	return options
}

// GetCustomResolverOptions : The GetCustomResolver options.
type GetCustomResolverOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetCustomResolverOptions : Instantiate GetCustomResolverOptions
func (*DnsSvcsV1) NewGetCustomResolverOptions(instanceID string, resolverID string) *GetCustomResolverOptions {
	return &GetCustomResolverOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetCustomResolverOptions) SetInstanceID(instanceID string) *GetCustomResolverOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *GetCustomResolverOptions) SetResolverID(resolverID string) *GetCustomResolverOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetCustomResolverOptions) SetXCorrelationID(xCorrelationID string) *GetCustomResolverOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCustomResolverOptions) SetHeaders(param map[string]string) *GetCustomResolverOptions {
	options.Headers = param
	return options
}

// GetDnszoneAccessRequestOptions : The GetDnszoneAccessRequest options.
type GetDnszoneAccessRequestOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of an access request.
	RequestID *string `json:"request_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDnszoneAccessRequestOptions : Instantiate GetDnszoneAccessRequestOptions
func (*DnsSvcsV1) NewGetDnszoneAccessRequestOptions(instanceID string, dnszoneID string, requestID string) *GetDnszoneAccessRequestOptions {
	return &GetDnszoneAccessRequestOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		RequestID:  core.StringPtr(requestID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetDnszoneAccessRequestOptions) SetInstanceID(instanceID string) *GetDnszoneAccessRequestOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *GetDnszoneAccessRequestOptions) SetDnszoneID(dnszoneID string) *GetDnszoneAccessRequestOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetRequestID : Allow user to set RequestID
func (_options *GetDnszoneAccessRequestOptions) SetRequestID(requestID string) *GetDnszoneAccessRequestOptions {
	_options.RequestID = core.StringPtr(requestID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetDnszoneAccessRequestOptions) SetXCorrelationID(xCorrelationID string) *GetDnszoneAccessRequestOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDnszoneAccessRequestOptions) SetHeaders(param map[string]string) *GetDnszoneAccessRequestOptions {
	options.Headers = param
	return options
}

// GetDnszoneOptions : The GetDnszone options.
type GetDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDnszoneOptions : Instantiate GetDnszoneOptions
func (*DnsSvcsV1) NewGetDnszoneOptions(instanceID string, dnszoneID string) *GetDnszoneOptions {
	return &GetDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetDnszoneOptions) SetInstanceID(instanceID string) *GetDnszoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *GetDnszoneOptions) SetDnszoneID(dnszoneID string) *GetDnszoneOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetDnszoneOptions) SetXCorrelationID(xCorrelationID string) *GetDnszoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDnszoneOptions) SetHeaders(param map[string]string) *GetDnszoneOptions {
	options.Headers = param
	return options
}

// GetForwardingRuleOptions : The GetForwardingRule options.
type GetForwardingRuleOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// The unique identifier of a rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetForwardingRuleOptions : Instantiate GetForwardingRuleOptions
func (*DnsSvcsV1) NewGetForwardingRuleOptions(instanceID string, resolverID string, ruleID string) *GetForwardingRuleOptions {
	return &GetForwardingRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetForwardingRuleOptions) SetInstanceID(instanceID string) *GetForwardingRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *GetForwardingRuleOptions) SetResolverID(resolverID string) *GetForwardingRuleOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *GetForwardingRuleOptions) SetRuleID(ruleID string) *GetForwardingRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetForwardingRuleOptions) SetXCorrelationID(xCorrelationID string) *GetForwardingRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetForwardingRuleOptions) SetHeaders(param map[string]string) *GetForwardingRuleOptions {
	options.Headers = param
	return options
}

// GetLinkedPermittedNetworkOptions : The GetLinkedPermittedNetwork options.
type GetLinkedPermittedNetworkOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a linked zone.
	LinkedDnszoneID *string `json:"linked_dnszone_id" validate:"required,ne="`

	// The unique identifier of a permitted network.
	PermittedNetworkID *string `json:"permitted_network_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLinkedPermittedNetworkOptions : Instantiate GetLinkedPermittedNetworkOptions
func (*DnsSvcsV1) NewGetLinkedPermittedNetworkOptions(instanceID string, linkedDnszoneID string, permittedNetworkID string) *GetLinkedPermittedNetworkOptions {
	return &GetLinkedPermittedNetworkOptions{
		InstanceID:         core.StringPtr(instanceID),
		LinkedDnszoneID:    core.StringPtr(linkedDnszoneID),
		PermittedNetworkID: core.StringPtr(permittedNetworkID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetLinkedPermittedNetworkOptions) SetInstanceID(instanceID string) *GetLinkedPermittedNetworkOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLinkedDnszoneID : Allow user to set LinkedDnszoneID
func (_options *GetLinkedPermittedNetworkOptions) SetLinkedDnszoneID(linkedDnszoneID string) *GetLinkedPermittedNetworkOptions {
	_options.LinkedDnszoneID = core.StringPtr(linkedDnszoneID)
	return _options
}

// SetPermittedNetworkID : Allow user to set PermittedNetworkID
func (_options *GetLinkedPermittedNetworkOptions) SetPermittedNetworkID(permittedNetworkID string) *GetLinkedPermittedNetworkOptions {
	_options.PermittedNetworkID = core.StringPtr(permittedNetworkID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetLinkedPermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *GetLinkedPermittedNetworkOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLinkedPermittedNetworkOptions) SetHeaders(param map[string]string) *GetLinkedPermittedNetworkOptions {
	options.Headers = param
	return options
}

// GetLinkedZoneOptions : The GetLinkedZone options.
type GetLinkedZoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a linked zone.
	LinkedDnszoneID *string `json:"linked_dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLinkedZoneOptions : Instantiate GetLinkedZoneOptions
func (*DnsSvcsV1) NewGetLinkedZoneOptions(instanceID string, linkedDnszoneID string) *GetLinkedZoneOptions {
	return &GetLinkedZoneOptions{
		InstanceID:      core.StringPtr(instanceID),
		LinkedDnszoneID: core.StringPtr(linkedDnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetLinkedZoneOptions) SetInstanceID(instanceID string) *GetLinkedZoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLinkedDnszoneID : Allow user to set LinkedDnszoneID
func (_options *GetLinkedZoneOptions) SetLinkedDnszoneID(linkedDnszoneID string) *GetLinkedZoneOptions {
	_options.LinkedDnszoneID = core.StringPtr(linkedDnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetLinkedZoneOptions) SetXCorrelationID(xCorrelationID string) *GetLinkedZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLinkedZoneOptions) SetHeaders(param map[string]string) *GetLinkedZoneOptions {
	options.Headers = param
	return options
}

// GetLoadBalancerOptions : The GetLoadBalancer options.
type GetLoadBalancerOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of a load balancer.
	LbID *string `json:"lb_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLoadBalancerOptions : Instantiate GetLoadBalancerOptions
func (*DnsSvcsV1) NewGetLoadBalancerOptions(instanceID string, dnszoneID string, lbID string) *GetLoadBalancerOptions {
	return &GetLoadBalancerOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		LbID:       core.StringPtr(lbID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetLoadBalancerOptions) SetInstanceID(instanceID string) *GetLoadBalancerOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *GetLoadBalancerOptions) SetDnszoneID(dnszoneID string) *GetLoadBalancerOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetLbID : Allow user to set LbID
func (_options *GetLoadBalancerOptions) SetLbID(lbID string) *GetLoadBalancerOptions {
	_options.LbID = core.StringPtr(lbID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetLoadBalancerOptions) SetXCorrelationID(xCorrelationID string) *GetLoadBalancerOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLoadBalancerOptions) SetHeaders(param map[string]string) *GetLoadBalancerOptions {
	options.Headers = param
	return options
}

// GetMonitorOptions : The GetMonitor options.
type GetMonitorOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a load balancer monitor.
	MonitorID *string `json:"monitor_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetMonitorOptions : Instantiate GetMonitorOptions
func (*DnsSvcsV1) NewGetMonitorOptions(instanceID string, monitorID string) *GetMonitorOptions {
	return &GetMonitorOptions{
		InstanceID: core.StringPtr(instanceID),
		MonitorID:  core.StringPtr(monitorID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetMonitorOptions) SetInstanceID(instanceID string) *GetMonitorOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetMonitorID : Allow user to set MonitorID
func (_options *GetMonitorOptions) SetMonitorID(monitorID string) *GetMonitorOptions {
	_options.MonitorID = core.StringPtr(monitorID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetMonitorOptions) SetXCorrelationID(xCorrelationID string) *GetMonitorOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetMonitorOptions) SetHeaders(param map[string]string) *GetMonitorOptions {
	options.Headers = param
	return options
}

// GetPermittedNetworkOptions : The GetPermittedNetwork options.
type GetPermittedNetworkOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of a permitted network.
	PermittedNetworkID *string `json:"permitted_network_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetPermittedNetworkOptions : Instantiate GetPermittedNetworkOptions
func (*DnsSvcsV1) NewGetPermittedNetworkOptions(instanceID string, dnszoneID string, permittedNetworkID string) *GetPermittedNetworkOptions {
	return &GetPermittedNetworkOptions{
		InstanceID:         core.StringPtr(instanceID),
		DnszoneID:          core.StringPtr(dnszoneID),
		PermittedNetworkID: core.StringPtr(permittedNetworkID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetPermittedNetworkOptions) SetInstanceID(instanceID string) *GetPermittedNetworkOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *GetPermittedNetworkOptions) SetDnszoneID(dnszoneID string) *GetPermittedNetworkOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetPermittedNetworkID : Allow user to set PermittedNetworkID
func (_options *GetPermittedNetworkOptions) SetPermittedNetworkID(permittedNetworkID string) *GetPermittedNetworkOptions {
	_options.PermittedNetworkID = core.StringPtr(permittedNetworkID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetPermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *GetPermittedNetworkOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPermittedNetworkOptions) SetHeaders(param map[string]string) *GetPermittedNetworkOptions {
	options.Headers = param
	return options
}

// GetPoolOptions : The GetPool options.
type GetPoolOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a load balancer pool.
	PoolID *string `json:"pool_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetPoolOptions : Instantiate GetPoolOptions
func (*DnsSvcsV1) NewGetPoolOptions(instanceID string, poolID string) *GetPoolOptions {
	return &GetPoolOptions{
		InstanceID: core.StringPtr(instanceID),
		PoolID:     core.StringPtr(poolID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetPoolOptions) SetInstanceID(instanceID string) *GetPoolOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetPoolID : Allow user to set PoolID
func (_options *GetPoolOptions) SetPoolID(poolID string) *GetPoolOptions {
	_options.PoolID = core.StringPtr(poolID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetPoolOptions) SetXCorrelationID(xCorrelationID string) *GetPoolOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPoolOptions) SetHeaders(param map[string]string) *GetPoolOptions {
	options.Headers = param
	return options
}

// GetResourceRecordOptions : The GetResourceRecord options.
type GetResourceRecordOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of a resource record.
	RecordID *string `json:"record_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetResourceRecordOptions : Instantiate GetResourceRecordOptions
func (*DnsSvcsV1) NewGetResourceRecordOptions(instanceID string, dnszoneID string, recordID string) *GetResourceRecordOptions {
	return &GetResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		RecordID:   core.StringPtr(recordID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetResourceRecordOptions) SetInstanceID(instanceID string) *GetResourceRecordOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *GetResourceRecordOptions) SetDnszoneID(dnszoneID string) *GetResourceRecordOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetRecordID : Allow user to set RecordID
func (_options *GetResourceRecordOptions) SetRecordID(recordID string) *GetResourceRecordOptions {
	_options.RecordID = core.StringPtr(recordID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetResourceRecordOptions) SetXCorrelationID(xCorrelationID string) *GetResourceRecordOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceRecordOptions) SetHeaders(param map[string]string) *GetResourceRecordOptions {
	options.Headers = param
	return options
}

// GetSecondaryZoneOptions : The GetSecondaryZone options.
type GetSecondaryZoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// The unique identifier of a secondary zone.
	SecondaryZoneID *string `json:"secondary_zone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetSecondaryZoneOptions : Instantiate GetSecondaryZoneOptions
func (*DnsSvcsV1) NewGetSecondaryZoneOptions(instanceID string, resolverID string, secondaryZoneID string) *GetSecondaryZoneOptions {
	return &GetSecondaryZoneOptions{
		InstanceID:      core.StringPtr(instanceID),
		ResolverID:      core.StringPtr(resolverID),
		SecondaryZoneID: core.StringPtr(secondaryZoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetSecondaryZoneOptions) SetInstanceID(instanceID string) *GetSecondaryZoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *GetSecondaryZoneOptions) SetResolverID(resolverID string) *GetSecondaryZoneOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetSecondaryZoneID : Allow user to set SecondaryZoneID
func (_options *GetSecondaryZoneOptions) SetSecondaryZoneID(secondaryZoneID string) *GetSecondaryZoneOptions {
	_options.SecondaryZoneID = core.StringPtr(secondaryZoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetSecondaryZoneOptions) SetXCorrelationID(xCorrelationID string) *GetSecondaryZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecondaryZoneOptions) SetHeaders(param map[string]string) *GetSecondaryZoneOptions {
	options.Headers = param
	return options
}

// ImportResourceRecordsOptions : The ImportResourceRecords options.
type ImportResourceRecordsOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// file to upload.
	File io.ReadCloser `json:"file,omitempty"`

	// The content type of file.
	FileContentType *string `json:"file_content_type,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewImportResourceRecordsOptions : Instantiate ImportResourceRecordsOptions
func (*DnsSvcsV1) NewImportResourceRecordsOptions(instanceID string, dnszoneID string) *ImportResourceRecordsOptions {
	return &ImportResourceRecordsOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ImportResourceRecordsOptions) SetInstanceID(instanceID string) *ImportResourceRecordsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *ImportResourceRecordsOptions) SetDnszoneID(dnszoneID string) *ImportResourceRecordsOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetFile : Allow user to set File
func (_options *ImportResourceRecordsOptions) SetFile(file io.ReadCloser) *ImportResourceRecordsOptions {
	_options.File = file
	return _options
}

// SetFileContentType : Allow user to set FileContentType
func (_options *ImportResourceRecordsOptions) SetFileContentType(fileContentType string) *ImportResourceRecordsOptions {
	_options.FileContentType = core.StringPtr(fileContentType)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ImportResourceRecordsOptions) SetXCorrelationID(xCorrelationID string) *ImportResourceRecordsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ImportResourceRecordsOptions) SetHeaders(param map[string]string) *ImportResourceRecordsOptions {
	options.Headers = param
	return options
}

// LinkedDnszoneLinkedTo : The owner's instance and zone that the zone is linked to.
type LinkedDnszoneLinkedTo struct {
	// The owner's instance CRN.
	InstanceCrn *string `json:"instance_crn,omitempty"`

	// The owner's DNS zone.
	ZoneID *string `json:"zone_id,omitempty"`
}

// UnmarshalLinkedDnszoneLinkedTo unmarshals an instance of LinkedDnszoneLinkedTo from the specified map of raw messages.
func UnmarshalLinkedDnszoneLinkedTo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LinkedDnszoneLinkedTo)
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zone_id", &obj.ZoneID)
	if err != nil {
		err = core.SDKErrorf(err, "", "zone_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListCustomResolversOptions : The ListCustomResolvers options.
type ListCustomResolversOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListCustomResolversOptions : Instantiate ListCustomResolversOptions
func (*DnsSvcsV1) NewListCustomResolversOptions(instanceID string) *ListCustomResolversOptions {
	return &ListCustomResolversOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListCustomResolversOptions) SetInstanceID(instanceID string) *ListCustomResolversOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListCustomResolversOptions) SetXCorrelationID(xCorrelationID string) *ListCustomResolversOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListCustomResolversOptions) SetHeaders(param map[string]string) *ListCustomResolversOptions {
	options.Headers = param
	return options
}

// ListDnszoneAccessRequestsOptions : The ListDnszoneAccessRequests options.
type ListDnszoneAccessRequestsOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListDnszoneAccessRequestsOptions : Instantiate ListDnszoneAccessRequestsOptions
func (*DnsSvcsV1) NewListDnszoneAccessRequestsOptions(instanceID string, dnszoneID string) *ListDnszoneAccessRequestsOptions {
	return &ListDnszoneAccessRequestsOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListDnszoneAccessRequestsOptions) SetInstanceID(instanceID string) *ListDnszoneAccessRequestsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *ListDnszoneAccessRequestsOptions) SetDnszoneID(dnszoneID string) *ListDnszoneAccessRequestsOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListDnszoneAccessRequestsOptions) SetXCorrelationID(xCorrelationID string) *ListDnszoneAccessRequestsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListDnszoneAccessRequestsOptions) SetOffset(offset int64) *ListDnszoneAccessRequestsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListDnszoneAccessRequestsOptions) SetLimit(limit int64) *ListDnszoneAccessRequestsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListDnszoneAccessRequestsOptions) SetHeaders(param map[string]string) *ListDnszoneAccessRequestsOptions {
	options.Headers = param
	return options
}

// ListDnszonesOptions : The ListDnszones options.
type ListDnszonesOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Specify the VPC ID.
	VpcID *string `json:"vpc_id,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListDnszonesOptions : Instantiate ListDnszonesOptions
func (*DnsSvcsV1) NewListDnszonesOptions(instanceID string) *ListDnszonesOptions {
	return &ListDnszonesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListDnszonesOptions) SetInstanceID(instanceID string) *ListDnszonesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListDnszonesOptions) SetXCorrelationID(xCorrelationID string) *ListDnszonesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListDnszonesOptions) SetOffset(offset int64) *ListDnszonesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListDnszonesOptions) SetLimit(limit int64) *ListDnszonesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetVpcID : Allow user to set VpcID
func (_options *ListDnszonesOptions) SetVpcID(vpcID string) *ListDnszonesOptions {
	_options.VpcID = core.StringPtr(vpcID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListDnszonesOptions) SetHeaders(param map[string]string) *ListDnszonesOptions {
	options.Headers = param
	return options
}

// ListForwardingRulesOptions : The ListForwardingRules options.
type ListForwardingRulesOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListForwardingRulesOptions : Instantiate ListForwardingRulesOptions
func (*DnsSvcsV1) NewListForwardingRulesOptions(instanceID string, resolverID string) *ListForwardingRulesOptions {
	return &ListForwardingRulesOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListForwardingRulesOptions) SetInstanceID(instanceID string) *ListForwardingRulesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *ListForwardingRulesOptions) SetResolverID(resolverID string) *ListForwardingRulesOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListForwardingRulesOptions) SetXCorrelationID(xCorrelationID string) *ListForwardingRulesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListForwardingRulesOptions) SetOffset(offset int64) *ListForwardingRulesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListForwardingRulesOptions) SetLimit(limit int64) *ListForwardingRulesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListForwardingRulesOptions) SetHeaders(param map[string]string) *ListForwardingRulesOptions {
	options.Headers = param
	return options
}

// ListLinkedPermittedNetworksOptions : The ListLinkedPermittedNetworks options.
type ListLinkedPermittedNetworksOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a linked zone.
	LinkedDnszoneID *string `json:"linked_dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListLinkedPermittedNetworksOptions : Instantiate ListLinkedPermittedNetworksOptions
func (*DnsSvcsV1) NewListLinkedPermittedNetworksOptions(instanceID string, linkedDnszoneID string) *ListLinkedPermittedNetworksOptions {
	return &ListLinkedPermittedNetworksOptions{
		InstanceID:      core.StringPtr(instanceID),
		LinkedDnszoneID: core.StringPtr(linkedDnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListLinkedPermittedNetworksOptions) SetInstanceID(instanceID string) *ListLinkedPermittedNetworksOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLinkedDnszoneID : Allow user to set LinkedDnszoneID
func (_options *ListLinkedPermittedNetworksOptions) SetLinkedDnszoneID(linkedDnszoneID string) *ListLinkedPermittedNetworksOptions {
	_options.LinkedDnszoneID = core.StringPtr(linkedDnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListLinkedPermittedNetworksOptions) SetXCorrelationID(xCorrelationID string) *ListLinkedPermittedNetworksOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListLinkedPermittedNetworksOptions) SetHeaders(param map[string]string) *ListLinkedPermittedNetworksOptions {
	options.Headers = param
	return options
}

// ListLinkedZonesOptions : The ListLinkedZones options.
type ListLinkedZonesOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListLinkedZonesOptions : Instantiate ListLinkedZonesOptions
func (*DnsSvcsV1) NewListLinkedZonesOptions(instanceID string) *ListLinkedZonesOptions {
	return &ListLinkedZonesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListLinkedZonesOptions) SetInstanceID(instanceID string) *ListLinkedZonesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListLinkedZonesOptions) SetXCorrelationID(xCorrelationID string) *ListLinkedZonesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListLinkedZonesOptions) SetOffset(offset int64) *ListLinkedZonesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListLinkedZonesOptions) SetLimit(limit int64) *ListLinkedZonesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListLinkedZonesOptions) SetHeaders(param map[string]string) *ListLinkedZonesOptions {
	options.Headers = param
	return options
}

// ListLoadBalancersOptions : The ListLoadBalancers options.
type ListLoadBalancersOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListLoadBalancersOptions : Instantiate ListLoadBalancersOptions
func (*DnsSvcsV1) NewListLoadBalancersOptions(instanceID string, dnszoneID string) *ListLoadBalancersOptions {
	return &ListLoadBalancersOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListLoadBalancersOptions) SetInstanceID(instanceID string) *ListLoadBalancersOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *ListLoadBalancersOptions) SetDnszoneID(dnszoneID string) *ListLoadBalancersOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListLoadBalancersOptions) SetXCorrelationID(xCorrelationID string) *ListLoadBalancersOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListLoadBalancersOptions) SetOffset(offset int64) *ListLoadBalancersOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListLoadBalancersOptions) SetLimit(limit int64) *ListLoadBalancersOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListLoadBalancersOptions) SetHeaders(param map[string]string) *ListLoadBalancersOptions {
	options.Headers = param
	return options
}

// ListMonitorsOptions : The ListMonitors options.
type ListMonitorsOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListMonitorsOptions : Instantiate ListMonitorsOptions
func (*DnsSvcsV1) NewListMonitorsOptions(instanceID string) *ListMonitorsOptions {
	return &ListMonitorsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListMonitorsOptions) SetInstanceID(instanceID string) *ListMonitorsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListMonitorsOptions) SetXCorrelationID(xCorrelationID string) *ListMonitorsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListMonitorsOptions) SetOffset(offset int64) *ListMonitorsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListMonitorsOptions) SetLimit(limit int64) *ListMonitorsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListMonitorsOptions) SetHeaders(param map[string]string) *ListMonitorsOptions {
	options.Headers = param
	return options
}

// ListPermittedNetworksOptions : The ListPermittedNetworks options.
type ListPermittedNetworksOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The account identifiers of the owner zone and linked zones in the format of "?account=account1,account2,account3".
	// Maximum 5 accounts are allowed.
	Accounts *string `json:"accounts,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListPermittedNetworksOptions : Instantiate ListPermittedNetworksOptions
func (*DnsSvcsV1) NewListPermittedNetworksOptions(instanceID string, dnszoneID string) *ListPermittedNetworksOptions {
	return &ListPermittedNetworksOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListPermittedNetworksOptions) SetInstanceID(instanceID string) *ListPermittedNetworksOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *ListPermittedNetworksOptions) SetDnszoneID(dnszoneID string) *ListPermittedNetworksOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetAccounts : Allow user to set Accounts
func (_options *ListPermittedNetworksOptions) SetAccounts(accounts string) *ListPermittedNetworksOptions {
	_options.Accounts = core.StringPtr(accounts)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListPermittedNetworksOptions) SetXCorrelationID(xCorrelationID string) *ListPermittedNetworksOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListPermittedNetworksOptions) SetHeaders(param map[string]string) *ListPermittedNetworksOptions {
	options.Headers = param
	return options
}

// ListPoolsOptions : The ListPools options.
type ListPoolsOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListPoolsOptions : Instantiate ListPoolsOptions
func (*DnsSvcsV1) NewListPoolsOptions(instanceID string) *ListPoolsOptions {
	return &ListPoolsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListPoolsOptions) SetInstanceID(instanceID string) *ListPoolsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListPoolsOptions) SetXCorrelationID(xCorrelationID string) *ListPoolsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListPoolsOptions) SetOffset(offset int64) *ListPoolsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListPoolsOptions) SetLimit(limit int64) *ListPoolsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListPoolsOptions) SetHeaders(param map[string]string) *ListPoolsOptions {
	options.Headers = param
	return options
}

// ListResourceRecordsOptions : The ListResourceRecords options.
type ListResourceRecordsOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Specify the type of resource record to query.
	Type *string `json:"type,omitempty"`

	// Specify the name of resource record to query.
	Name *string `json:"name,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListResourceRecordsOptions : Instantiate ListResourceRecordsOptions
func (*DnsSvcsV1) NewListResourceRecordsOptions(instanceID string, dnszoneID string) *ListResourceRecordsOptions {
	return &ListResourceRecordsOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListResourceRecordsOptions) SetInstanceID(instanceID string) *ListResourceRecordsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *ListResourceRecordsOptions) SetDnszoneID(dnszoneID string) *ListResourceRecordsOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListResourceRecordsOptions) SetXCorrelationID(xCorrelationID string) *ListResourceRecordsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListResourceRecordsOptions) SetOffset(offset int64) *ListResourceRecordsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListResourceRecordsOptions) SetLimit(limit int64) *ListResourceRecordsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListResourceRecordsOptions) SetType(typeVar string) *ListResourceRecordsOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListResourceRecordsOptions) SetName(name string) *ListResourceRecordsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListResourceRecordsOptions) SetHeaders(param map[string]string) *ListResourceRecordsOptions {
	options.Headers = param
	return options
}

// ListSecondaryZonesOptions : The ListSecondaryZones options.
type ListSecondaryZonesOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resources to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify maximum resources might be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListSecondaryZonesOptions : Instantiate ListSecondaryZonesOptions
func (*DnsSvcsV1) NewListSecondaryZonesOptions(instanceID string, resolverID string) *ListSecondaryZonesOptions {
	return &ListSecondaryZonesOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListSecondaryZonesOptions) SetInstanceID(instanceID string) *ListSecondaryZonesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *ListSecondaryZonesOptions) SetResolverID(resolverID string) *ListSecondaryZonesOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListSecondaryZonesOptions) SetXCorrelationID(xCorrelationID string) *ListSecondaryZonesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListSecondaryZonesOptions) SetOffset(offset int64) *ListSecondaryZonesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSecondaryZonesOptions) SetLimit(limit int64) *ListSecondaryZonesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecondaryZonesOptions) SetHeaders(param map[string]string) *ListSecondaryZonesOptions {
	options.Headers = param
	return options
}

// PoolHealthcheckVsisItem : PoolHealthcheckVsisItem struct
type PoolHealthcheckVsisItem struct {
	// Health check VSI subnet CRN.
	Subnet *string `json:"subnet,omitempty"`

	// healthcheck VSI ip address.
	Ipv4Address *string `json:"ipv4_address,omitempty"`

	// ipv4 cidr block.
	Ipv4CidrBlock *string `json:"ipv4_cidr_block,omitempty"`

	// vpc crn.
	Vpc *string `json:"vpc,omitempty"`
}

// UnmarshalPoolHealthcheckVsisItem unmarshals an instance of PoolHealthcheckVsisItem from the specified map of raw messages.
func UnmarshalPoolHealthcheckVsisItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PoolHealthcheckVsisItem)
	err = core.UnmarshalPrimitive(m, "subnet", &obj.Subnet)
	if err != nil {
		err = core.SDKErrorf(err, "", "subnet-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ipv4_address", &obj.Ipv4Address)
	if err != nil {
		err = core.SDKErrorf(err, "", "ipv4_address-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ipv4_cidr_block", &obj.Ipv4CidrBlock)
	if err != nil {
		err = core.SDKErrorf(err, "", "ipv4_cidr_block-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "vpc", &obj.Vpc)
	if err != nil {
		err = core.SDKErrorf(err, "", "vpc-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RecordsImportErrorModelError : Error container.
type RecordsImportErrorModelError struct {
	// Internal service error when DNS resource created fails by internal error.
	Code *string `json:"code" validate:"required"`

	// An internal error occurred. Try again later.
	Message *string `json:"message" validate:"required"`
}

// UnmarshalRecordsImportErrorModelError unmarshals an instance of RecordsImportErrorModelError from the specified map of raw messages.
func UnmarshalRecordsImportErrorModelError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RecordsImportErrorModelError)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		err = core.SDKErrorf(err, "", "code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordInputRdata : Content of the resource record.
// Models which "extend" this model:
// - ResourceRecordInputRdataRdataARecord
// - ResourceRecordInputRdataRdataAaaaRecord
// - ResourceRecordInputRdataRdataCnameRecord
// - ResourceRecordInputRdataRdataMxRecord
// - ResourceRecordInputRdataRdataSrvRecord
// - ResourceRecordInputRdataRdataTxtRecord
// - ResourceRecordInputRdataRdataPtrRecord
type ResourceRecordInputRdata struct {
	// IPv4 address.
	Ip *string `json:"ip,omitempty"`

	// Canonical name.
	Cname *string `json:"cname,omitempty"`

	// Hostname of Exchange server.
	Exchange *string `json:"exchange,omitempty"`

	// Preference of the MX record.
	Preference *int64 `json:"preference,omitempty"`

	// Port number of the target server.
	Port *int64 `json:"port,omitempty"`

	// Priority of the SRV record.
	Priority *int64 `json:"priority,omitempty"`

	// Hostname of the target server.
	Target *string `json:"target,omitempty"`

	// Weight of distributing queries among multiple target servers.
	Weight *int64 `json:"weight,omitempty"`

	// Human readable text.
	Text *string `json:"text,omitempty"`

	// Hostname of the relevant A or AAAA record.
	Ptrdname *string `json:"ptrdname,omitempty"`
}

func (*ResourceRecordInputRdata) isaResourceRecordInputRdata() bool {
	return true
}

type ResourceRecordInputRdataIntf interface {
	isaResourceRecordInputRdata() bool
}

// UnmarshalResourceRecordInputRdata unmarshals an instance of ResourceRecordInputRdata from the specified map of raw messages.
func UnmarshalResourceRecordInputRdata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordInputRdata)
	err = core.UnmarshalPrimitive(m, "ip", &obj.Ip)
	if err != nil {
		err = core.SDKErrorf(err, "", "ip-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cname", &obj.Cname)
	if err != nil {
		err = core.SDKErrorf(err, "", "cname-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "exchange", &obj.Exchange)
	if err != nil {
		err = core.SDKErrorf(err, "", "exchange-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "preference", &obj.Preference)
	if err != nil {
		err = core.SDKErrorf(err, "", "preference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "weight", &obj.Weight)
	if err != nil {
		err = core.SDKErrorf(err, "", "weight-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		err = core.SDKErrorf(err, "", "text-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ptrdname", &obj.Ptrdname)
	if err != nil {
		err = core.SDKErrorf(err, "", "ptrdname-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordUpdateInputRdata : Content of the resource record.
// Models which "extend" this model:
// - ResourceRecordUpdateInputRdataRdataARecord
// - ResourceRecordUpdateInputRdataRdataAaaaRecord
// - ResourceRecordUpdateInputRdataRdataCnameRecord
// - ResourceRecordUpdateInputRdataRdataMxRecord
// - ResourceRecordUpdateInputRdataRdataSrvRecord
// - ResourceRecordUpdateInputRdataRdataTxtRecord
// - ResourceRecordUpdateInputRdataRdataPtrRecord
type ResourceRecordUpdateInputRdata struct {
	// IPv4 address.
	Ip *string `json:"ip,omitempty"`

	// Canonical name.
	Cname *string `json:"cname,omitempty"`

	// Hostname of Exchange server.
	Exchange *string `json:"exchange,omitempty"`

	// Preference of the MX record.
	Preference *int64 `json:"preference,omitempty"`

	// Port number of the target server.
	Port *int64 `json:"port,omitempty"`

	// Priority of the SRV record.
	Priority *int64 `json:"priority,omitempty"`

	// Hostname of the target server.
	Target *string `json:"target,omitempty"`

	// Weight of distributing queries among multiple target servers.
	Weight *int64 `json:"weight,omitempty"`

	// Human readable text.
	Text *string `json:"text,omitempty"`

	// Hostname of the relevant A or AAAA record.
	Ptrdname *string `json:"ptrdname,omitempty"`
}

func (*ResourceRecordUpdateInputRdata) isaResourceRecordUpdateInputRdata() bool {
	return true
}

type ResourceRecordUpdateInputRdataIntf interface {
	isaResourceRecordUpdateInputRdata() bool
}

// UnmarshalResourceRecordUpdateInputRdata unmarshals an instance of ResourceRecordUpdateInputRdata from the specified map of raw messages.
func UnmarshalResourceRecordUpdateInputRdata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordUpdateInputRdata)
	err = core.UnmarshalPrimitive(m, "ip", &obj.Ip)
	if err != nil {
		err = core.SDKErrorf(err, "", "ip-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cname", &obj.Cname)
	if err != nil {
		err = core.SDKErrorf(err, "", "cname-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "exchange", &obj.Exchange)
	if err != nil {
		err = core.SDKErrorf(err, "", "exchange-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "preference", &obj.Preference)
	if err != nil {
		err = core.SDKErrorf(err, "", "preference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "weight", &obj.Weight)
	if err != nil {
		err = core.SDKErrorf(err, "", "weight-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		err = core.SDKErrorf(err, "", "text-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ptrdname", &obj.Ptrdname)
	if err != nil {
		err = core.SDKErrorf(err, "", "ptrdname-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateCrLocationsOrderOptions : The UpdateCrLocationsOrder options.
type UpdateCrLocationsOrderOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Array of custom resolver location ID.
	Locations []string `json:"locations" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateCrLocationsOrderOptions : Instantiate UpdateCrLocationsOrderOptions
func (*DnsSvcsV1) NewUpdateCrLocationsOrderOptions(instanceID string, resolverID string, locations []string) *UpdateCrLocationsOrderOptions {
	return &UpdateCrLocationsOrderOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
		Locations:  locations,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateCrLocationsOrderOptions) SetInstanceID(instanceID string) *UpdateCrLocationsOrderOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *UpdateCrLocationsOrderOptions) SetResolverID(resolverID string) *UpdateCrLocationsOrderOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetLocations : Allow user to set Locations
func (_options *UpdateCrLocationsOrderOptions) SetLocations(locations []string) *UpdateCrLocationsOrderOptions {
	_options.Locations = locations
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateCrLocationsOrderOptions) SetXCorrelationID(xCorrelationID string) *UpdateCrLocationsOrderOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCrLocationsOrderOptions) SetHeaders(param map[string]string) *UpdateCrLocationsOrderOptions {
	options.Headers = param
	return options
}

// UpdateCustomResolverLocationOptions : The UpdateCustomResolverLocation options.
type UpdateCustomResolverLocationOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Custom resolver location ID.
	LocationID *string `json:"location_id" validate:"required,ne="`

	// Enable/Disable custom resolver location.
	Enabled *bool `json:"enabled,omitempty"`

	// Subnet CRN.
	SubnetCrn *string `json:"subnet_crn,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateCustomResolverLocationOptions : Instantiate UpdateCustomResolverLocationOptions
func (*DnsSvcsV1) NewUpdateCustomResolverLocationOptions(instanceID string, resolverID string, locationID string) *UpdateCustomResolverLocationOptions {
	return &UpdateCustomResolverLocationOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
		LocationID: core.StringPtr(locationID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateCustomResolverLocationOptions) SetInstanceID(instanceID string) *UpdateCustomResolverLocationOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *UpdateCustomResolverLocationOptions) SetResolverID(resolverID string) *UpdateCustomResolverLocationOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetLocationID : Allow user to set LocationID
func (_options *UpdateCustomResolverLocationOptions) SetLocationID(locationID string) *UpdateCustomResolverLocationOptions {
	_options.LocationID = core.StringPtr(locationID)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateCustomResolverLocationOptions) SetEnabled(enabled bool) *UpdateCustomResolverLocationOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetSubnetCrn : Allow user to set SubnetCrn
func (_options *UpdateCustomResolverLocationOptions) SetSubnetCrn(subnetCrn string) *UpdateCustomResolverLocationOptions {
	_options.SubnetCrn = core.StringPtr(subnetCrn)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateCustomResolverLocationOptions) SetXCorrelationID(xCorrelationID string) *UpdateCustomResolverLocationOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCustomResolverLocationOptions) SetHeaders(param map[string]string) *UpdateCustomResolverLocationOptions {
	options.Headers = param
	return options
}

// UpdateCustomResolverOptions : The UpdateCustomResolver options.
type UpdateCustomResolverOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// Name of the custom resolver.
	Name *string `json:"name,omitempty"`

	// Descriptive text of the custom resolver.
	Description *string `json:"description,omitempty"`

	// Whether the custom resolver is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The profile name of a custom resolver.
	Profile *string `json:"profile,omitempty"`

	// Whether a disruptive update is allowed for the custom resolver.
	AllowDisruptiveUpdates *bool `json:"allow_disruptive_updates,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the UpdateCustomResolverOptions.Profile property.
// The profile name of a custom resolver.
const (
	UpdateCustomResolverOptions_Profile_Advanced  = "advanced"
	UpdateCustomResolverOptions_Profile_Essential = "essential"
	UpdateCustomResolverOptions_Profile_Premier   = "premier"
)

// NewUpdateCustomResolverOptions : Instantiate UpdateCustomResolverOptions
func (*DnsSvcsV1) NewUpdateCustomResolverOptions(instanceID string, resolverID string) *UpdateCustomResolverOptions {
	return &UpdateCustomResolverOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateCustomResolverOptions) SetInstanceID(instanceID string) *UpdateCustomResolverOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *UpdateCustomResolverOptions) SetResolverID(resolverID string) *UpdateCustomResolverOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateCustomResolverOptions) SetName(name string) *UpdateCustomResolverOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateCustomResolverOptions) SetDescription(description string) *UpdateCustomResolverOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateCustomResolverOptions) SetEnabled(enabled bool) *UpdateCustomResolverOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *UpdateCustomResolverOptions) SetProfile(profile string) *UpdateCustomResolverOptions {
	_options.Profile = core.StringPtr(profile)
	return _options
}

// SetAllowDisruptiveUpdates : Allow user to set AllowDisruptiveUpdates
func (_options *UpdateCustomResolverOptions) SetAllowDisruptiveUpdates(allowDisruptiveUpdates bool) *UpdateCustomResolverOptions {
	_options.AllowDisruptiveUpdates = core.BoolPtr(allowDisruptiveUpdates)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateCustomResolverOptions) SetXCorrelationID(xCorrelationID string) *UpdateCustomResolverOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCustomResolverOptions) SetHeaders(param map[string]string) *UpdateCustomResolverOptions {
	options.Headers = param
	return options
}

// UpdateDnszoneAccessRequestOptions : The UpdateDnszoneAccessRequest options.
type UpdateDnszoneAccessRequestOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of an access request.
	RequestID *string `json:"request_id" validate:"required,ne="`

	// The action applies to the access request.
	Action *string `json:"action" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the UpdateDnszoneAccessRequestOptions.Action property.
// The action applies to the access request.
const (
	UpdateDnszoneAccessRequestOptions_Action_Approve = "APPROVE"
	UpdateDnszoneAccessRequestOptions_Action_Reject  = "REJECT"
	UpdateDnszoneAccessRequestOptions_Action_Revoke  = "REVOKE"
)

// NewUpdateDnszoneAccessRequestOptions : Instantiate UpdateDnszoneAccessRequestOptions
func (*DnsSvcsV1) NewUpdateDnszoneAccessRequestOptions(instanceID string, dnszoneID string, requestID string, action string) *UpdateDnszoneAccessRequestOptions {
	return &UpdateDnszoneAccessRequestOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		RequestID:  core.StringPtr(requestID),
		Action:     core.StringPtr(action),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateDnszoneAccessRequestOptions) SetInstanceID(instanceID string) *UpdateDnszoneAccessRequestOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *UpdateDnszoneAccessRequestOptions) SetDnszoneID(dnszoneID string) *UpdateDnszoneAccessRequestOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetRequestID : Allow user to set RequestID
func (_options *UpdateDnszoneAccessRequestOptions) SetRequestID(requestID string) *UpdateDnszoneAccessRequestOptions {
	_options.RequestID = core.StringPtr(requestID)
	return _options
}

// SetAction : Allow user to set Action
func (_options *UpdateDnszoneAccessRequestOptions) SetAction(action string) *UpdateDnszoneAccessRequestOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateDnszoneAccessRequestOptions) SetXCorrelationID(xCorrelationID string) *UpdateDnszoneAccessRequestOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDnszoneAccessRequestOptions) SetHeaders(param map[string]string) *UpdateDnszoneAccessRequestOptions {
	options.Headers = param
	return options
}

// UpdateDnszoneOptions : The UpdateDnszone options.
type UpdateDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The text describing the purpose of a DNS zone.
	Description *string `json:"description,omitempty"`

	// The label of a DNS zone.
	Label *string `json:"label,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateDnszoneOptions : Instantiate UpdateDnszoneOptions
func (*DnsSvcsV1) NewUpdateDnszoneOptions(instanceID string, dnszoneID string) *UpdateDnszoneOptions {
	return &UpdateDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateDnszoneOptions) SetInstanceID(instanceID string) *UpdateDnszoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *UpdateDnszoneOptions) SetDnszoneID(dnszoneID string) *UpdateDnszoneOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateDnszoneOptions) SetDescription(description string) *UpdateDnszoneOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *UpdateDnszoneOptions) SetLabel(label string) *UpdateDnszoneOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateDnszoneOptions) SetXCorrelationID(xCorrelationID string) *UpdateDnszoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDnszoneOptions) SetHeaders(param map[string]string) *UpdateDnszoneOptions {
	options.Headers = param
	return options
}

// UpdateForwardingRuleOptions : The UpdateForwardingRule options.
type UpdateForwardingRuleOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// The unique identifier of a rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Descriptive text of the forwarding rule.
	Description *string `json:"description,omitempty"`

	// The matching zone or hostname. For the default forwarding rule, the match must always be the wildcard '*', and can
	// not be changed.
	Match *string `json:"match,omitempty"`

	// The upstream DNS servers that the DNS queries will be forwarded to.
	ForwardTo []string `json:"forward_to,omitempty"`

	// An array of views.
	Views []ViewConfig `json:"views,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateForwardingRuleOptions : Instantiate UpdateForwardingRuleOptions
func (*DnsSvcsV1) NewUpdateForwardingRuleOptions(instanceID string, resolverID string, ruleID string) *UpdateForwardingRuleOptions {
	return &UpdateForwardingRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		ResolverID: core.StringPtr(resolverID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateForwardingRuleOptions) SetInstanceID(instanceID string) *UpdateForwardingRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *UpdateForwardingRuleOptions) SetResolverID(resolverID string) *UpdateForwardingRuleOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *UpdateForwardingRuleOptions) SetRuleID(ruleID string) *UpdateForwardingRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateForwardingRuleOptions) SetDescription(description string) *UpdateForwardingRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetMatch : Allow user to set Match
func (_options *UpdateForwardingRuleOptions) SetMatch(match string) *UpdateForwardingRuleOptions {
	_options.Match = core.StringPtr(match)
	return _options
}

// SetForwardTo : Allow user to set ForwardTo
func (_options *UpdateForwardingRuleOptions) SetForwardTo(forwardTo []string) *UpdateForwardingRuleOptions {
	_options.ForwardTo = forwardTo
	return _options
}

// SetViews : Allow user to set Views
func (_options *UpdateForwardingRuleOptions) SetViews(views []ViewConfig) *UpdateForwardingRuleOptions {
	_options.Views = views
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateForwardingRuleOptions) SetXCorrelationID(xCorrelationID string) *UpdateForwardingRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateForwardingRuleOptions) SetHeaders(param map[string]string) *UpdateForwardingRuleOptions {
	options.Headers = param
	return options
}

// UpdateLinkedZoneOptions : The UpdateLinkedZone options.
type UpdateLinkedZoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a linked zone.
	LinkedDnszoneID *string `json:"linked_dnszone_id" validate:"required,ne="`

	// Descriptive text of the linked zone.
	Description *string `json:"description,omitempty"`

	// The label of linked zone.
	Label *string `json:"label,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateLinkedZoneOptions : Instantiate UpdateLinkedZoneOptions
func (*DnsSvcsV1) NewUpdateLinkedZoneOptions(instanceID string, linkedDnszoneID string) *UpdateLinkedZoneOptions {
	return &UpdateLinkedZoneOptions{
		InstanceID:      core.StringPtr(instanceID),
		LinkedDnszoneID: core.StringPtr(linkedDnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateLinkedZoneOptions) SetInstanceID(instanceID string) *UpdateLinkedZoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLinkedDnszoneID : Allow user to set LinkedDnszoneID
func (_options *UpdateLinkedZoneOptions) SetLinkedDnszoneID(linkedDnszoneID string) *UpdateLinkedZoneOptions {
	_options.LinkedDnszoneID = core.StringPtr(linkedDnszoneID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateLinkedZoneOptions) SetDescription(description string) *UpdateLinkedZoneOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *UpdateLinkedZoneOptions) SetLabel(label string) *UpdateLinkedZoneOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateLinkedZoneOptions) SetXCorrelationID(xCorrelationID string) *UpdateLinkedZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateLinkedZoneOptions) SetHeaders(param map[string]string) *UpdateLinkedZoneOptions {
	options.Headers = param
	return options
}

// UpdateLoadBalancerOptions : The UpdateLoadBalancer options.
type UpdateLoadBalancerOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of a load balancer.
	LbID *string `json:"lb_id" validate:"required,ne="`

	// Name of the load balancer.
	Name *string `json:"name,omitempty"`

	// Descriptive text of the load balancer.
	Description *string `json:"description,omitempty"`

	// Whether the load balancer is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Time to live in second.
	TTL *int64 `json:"ttl,omitempty"`

	// The pool ID to use when all other pools are detected as unhealthy.
	FallbackPool *string `json:"fallback_pool,omitempty"`

	// A list of pool IDs ordered by their failover priority. Pools defined here are used by default, or when region_pools
	// are not configured for a given region.
	DefaultPools []string `json:"default_pools,omitempty"`

	// Map availability zones to pool IDs.
	AzPools []AzPoolsItem `json:"az_pools,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateLoadBalancerOptions : Instantiate UpdateLoadBalancerOptions
func (*DnsSvcsV1) NewUpdateLoadBalancerOptions(instanceID string, dnszoneID string, lbID string) *UpdateLoadBalancerOptions {
	return &UpdateLoadBalancerOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		LbID:       core.StringPtr(lbID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateLoadBalancerOptions) SetInstanceID(instanceID string) *UpdateLoadBalancerOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *UpdateLoadBalancerOptions) SetDnszoneID(dnszoneID string) *UpdateLoadBalancerOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetLbID : Allow user to set LbID
func (_options *UpdateLoadBalancerOptions) SetLbID(lbID string) *UpdateLoadBalancerOptions {
	_options.LbID = core.StringPtr(lbID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateLoadBalancerOptions) SetName(name string) *UpdateLoadBalancerOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateLoadBalancerOptions) SetDescription(description string) *UpdateLoadBalancerOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateLoadBalancerOptions) SetEnabled(enabled bool) *UpdateLoadBalancerOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetTTL : Allow user to set TTL
func (_options *UpdateLoadBalancerOptions) SetTTL(ttl int64) *UpdateLoadBalancerOptions {
	_options.TTL = core.Int64Ptr(ttl)
	return _options
}

// SetFallbackPool : Allow user to set FallbackPool
func (_options *UpdateLoadBalancerOptions) SetFallbackPool(fallbackPool string) *UpdateLoadBalancerOptions {
	_options.FallbackPool = core.StringPtr(fallbackPool)
	return _options
}

// SetDefaultPools : Allow user to set DefaultPools
func (_options *UpdateLoadBalancerOptions) SetDefaultPools(defaultPools []string) *UpdateLoadBalancerOptions {
	_options.DefaultPools = defaultPools
	return _options
}

// SetAzPools : Allow user to set AzPools
func (_options *UpdateLoadBalancerOptions) SetAzPools(azPools []AzPoolsItem) *UpdateLoadBalancerOptions {
	_options.AzPools = azPools
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateLoadBalancerOptions) SetXCorrelationID(xCorrelationID string) *UpdateLoadBalancerOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateLoadBalancerOptions) SetHeaders(param map[string]string) *UpdateLoadBalancerOptions {
	options.Headers = param
	return options
}

// UpdateMonitorOptions : The UpdateMonitor options.
type UpdateMonitorOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a load balancer monitor.
	MonitorID *string `json:"monitor_id" validate:"required,ne="`

	// The name of the load balancer monitor.
	Name *string `json:"name,omitempty"`

	// Descriptive text of the load balancer monitor.
	Description *string `json:"description,omitempty"`

	// The protocol to use for the health check. Currently supported protocols are 'HTTP','HTTPS' and 'TCP'.
	Type *string `json:"type,omitempty"`

	// Port number to connect to for the health check. Required for TCP checks. HTTP and HTTPS checks should only define
	// the port when using a non-standard port (HTTP: default 80, HTTPS: default 443).
	Port *int64 `json:"port,omitempty"`

	// The interval between each health check. Shorter intervals may improve failover time, but will increase load on the
	// origins as we check from multiple locations.
	Interval *int64 `json:"interval,omitempty"`

	// The number of retries to attempt in case of a timeout before marking the origin as unhealthy. Retries are attempted
	// immediately.
	Retries *int64 `json:"retries,omitempty"`

	// The timeout (in seconds) before marking the health check as failed.
	Timeout *int64 `json:"timeout,omitempty"`

	// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
	Method *string `json:"method,omitempty"`

	// The endpoint path to health check against. This parameter is only valid for HTTP and HTTPS monitors.
	Path *string `json:"path,omitempty"`

	// The HTTP request headers to send in the health check. It is recommended you set a Host header by default. The
	// User-Agent header cannot be overridden. This parameter is only valid for HTTP and HTTPS monitors.
	HeadersVar []HealthcheckHeader `json:"headers,omitempty"`

	// Do not validate the certificate when monitor use HTTPS. This parameter is currently only valid for HTTP and HTTPS
	// monitors.
	AllowInsecure *bool `json:"allow_insecure,omitempty"`

	// The expected HTTP response code or code range of the health check. This parameter is only valid for HTTP and HTTPS
	// monitors.
	ExpectedCodes *string `json:"expected_codes,omitempty"`

	// A case-insensitive sub-string to look for in the response body. If this string is not found, the origin will be
	// marked as unhealthy. This parameter is only valid for HTTP and HTTPS monitors.
	ExpectedBody *string `json:"expected_body,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the UpdateMonitorOptions.Type property.
// The protocol to use for the health check. Currently supported protocols are 'HTTP','HTTPS' and 'TCP'.
const (
	UpdateMonitorOptions_Type_Http  = "HTTP"
	UpdateMonitorOptions_Type_Https = "HTTPS"
	UpdateMonitorOptions_Type_Tcp   = "TCP"
)

// Constants associated with the UpdateMonitorOptions.Method property.
// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
const (
	UpdateMonitorOptions_Method_Get  = "GET"
	UpdateMonitorOptions_Method_Head = "HEAD"
)

// NewUpdateMonitorOptions : Instantiate UpdateMonitorOptions
func (*DnsSvcsV1) NewUpdateMonitorOptions(instanceID string, monitorID string) *UpdateMonitorOptions {
	return &UpdateMonitorOptions{
		InstanceID: core.StringPtr(instanceID),
		MonitorID:  core.StringPtr(monitorID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateMonitorOptions) SetInstanceID(instanceID string) *UpdateMonitorOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetMonitorID : Allow user to set MonitorID
func (_options *UpdateMonitorOptions) SetMonitorID(monitorID string) *UpdateMonitorOptions {
	_options.MonitorID = core.StringPtr(monitorID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateMonitorOptions) SetName(name string) *UpdateMonitorOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateMonitorOptions) SetDescription(description string) *UpdateMonitorOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetType : Allow user to set Type
func (_options *UpdateMonitorOptions) SetType(typeVar string) *UpdateMonitorOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetPort : Allow user to set Port
func (_options *UpdateMonitorOptions) SetPort(port int64) *UpdateMonitorOptions {
	_options.Port = core.Int64Ptr(port)
	return _options
}

// SetInterval : Allow user to set Interval
func (_options *UpdateMonitorOptions) SetInterval(interval int64) *UpdateMonitorOptions {
	_options.Interval = core.Int64Ptr(interval)
	return _options
}

// SetRetries : Allow user to set Retries
func (_options *UpdateMonitorOptions) SetRetries(retries int64) *UpdateMonitorOptions {
	_options.Retries = core.Int64Ptr(retries)
	return _options
}

// SetTimeout : Allow user to set Timeout
func (_options *UpdateMonitorOptions) SetTimeout(timeout int64) *UpdateMonitorOptions {
	_options.Timeout = core.Int64Ptr(timeout)
	return _options
}

// SetMethod : Allow user to set Method
func (_options *UpdateMonitorOptions) SetMethod(method string) *UpdateMonitorOptions {
	_options.Method = core.StringPtr(method)
	return _options
}

// SetPath : Allow user to set Path
func (_options *UpdateMonitorOptions) SetPath(path string) *UpdateMonitorOptions {
	_options.Path = core.StringPtr(path)
	return _options
}

// SetHeadersVar : Allow user to set HeadersVar
func (_options *UpdateMonitorOptions) SetHeadersVar(headersVar []HealthcheckHeader) *UpdateMonitorOptions {
	_options.HeadersVar = headersVar
	return _options
}

// SetAllowInsecure : Allow user to set AllowInsecure
func (_options *UpdateMonitorOptions) SetAllowInsecure(allowInsecure bool) *UpdateMonitorOptions {
	_options.AllowInsecure = core.BoolPtr(allowInsecure)
	return _options
}

// SetExpectedCodes : Allow user to set ExpectedCodes
func (_options *UpdateMonitorOptions) SetExpectedCodes(expectedCodes string) *UpdateMonitorOptions {
	_options.ExpectedCodes = core.StringPtr(expectedCodes)
	return _options
}

// SetExpectedBody : Allow user to set ExpectedBody
func (_options *UpdateMonitorOptions) SetExpectedBody(expectedBody string) *UpdateMonitorOptions {
	_options.ExpectedBody = core.StringPtr(expectedBody)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateMonitorOptions) SetXCorrelationID(xCorrelationID string) *UpdateMonitorOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateMonitorOptions) SetHeaders(param map[string]string) *UpdateMonitorOptions {
	options.Headers = param
	return options
}

// UpdatePoolOptions : The UpdatePool options.
type UpdatePoolOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a load balancer pool.
	PoolID *string `json:"pool_id" validate:"required,ne="`

	// Name of the load balancer pool.
	Name *string `json:"name,omitempty"`

	// Descriptive text of the load balancer pool.
	Description *string `json:"description,omitempty"`

	// Whether the load balancer pool is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The minimum number of origins that must be healthy for this pool to serve traffic. If the number of healthy origins
	// falls below this number, the pool will be marked unhealthy and we will failover to the next available pool.
	HealthyOriginsThreshold *int64 `json:"healthy_origins_threshold,omitempty"`

	// The list of origins within this pool. Traffic directed at this pool is balanced across all currently healthy
	// origins, provided the pool itself is healthy.
	Origins []OriginInput `json:"origins,omitempty"`

	// The ID of the load balancer monitor to be associated to this pool.
	Monitor *string `json:"monitor,omitempty"`

	// The notification channel.
	NotificationChannel *string `json:"notification_channel,omitempty"`

	// Health check region of VSIs.
	HealthcheckRegion *string `json:"healthcheck_region,omitempty"`

	// Health check subnet CRNs.
	HealthcheckSubnets []string `json:"healthcheck_subnets,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the UpdatePoolOptions.HealthcheckRegion property.
// Health check region of VSIs.
const (
	UpdatePoolOptions_HealthcheckRegion_AuSyd   = "au-syd"
	UpdatePoolOptions_HealthcheckRegion_EuDu    = "eu-du"
	UpdatePoolOptions_HealthcheckRegion_EuGb    = "eu-gb"
	UpdatePoolOptions_HealthcheckRegion_JpTok   = "jp-tok"
	UpdatePoolOptions_HealthcheckRegion_UsEast  = "us-east"
	UpdatePoolOptions_HealthcheckRegion_UsSouth = "us-south"
)

// NewUpdatePoolOptions : Instantiate UpdatePoolOptions
func (*DnsSvcsV1) NewUpdatePoolOptions(instanceID string, poolID string) *UpdatePoolOptions {
	return &UpdatePoolOptions{
		InstanceID: core.StringPtr(instanceID),
		PoolID:     core.StringPtr(poolID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdatePoolOptions) SetInstanceID(instanceID string) *UpdatePoolOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetPoolID : Allow user to set PoolID
func (_options *UpdatePoolOptions) SetPoolID(poolID string) *UpdatePoolOptions {
	_options.PoolID = core.StringPtr(poolID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdatePoolOptions) SetName(name string) *UpdatePoolOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdatePoolOptions) SetDescription(description string) *UpdatePoolOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdatePoolOptions) SetEnabled(enabled bool) *UpdatePoolOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetHealthyOriginsThreshold : Allow user to set HealthyOriginsThreshold
func (_options *UpdatePoolOptions) SetHealthyOriginsThreshold(healthyOriginsThreshold int64) *UpdatePoolOptions {
	_options.HealthyOriginsThreshold = core.Int64Ptr(healthyOriginsThreshold)
	return _options
}

// SetOrigins : Allow user to set Origins
func (_options *UpdatePoolOptions) SetOrigins(origins []OriginInput) *UpdatePoolOptions {
	_options.Origins = origins
	return _options
}

// SetMonitor : Allow user to set Monitor
func (_options *UpdatePoolOptions) SetMonitor(monitor string) *UpdatePoolOptions {
	_options.Monitor = core.StringPtr(monitor)
	return _options
}

// SetNotificationChannel : Allow user to set NotificationChannel
func (_options *UpdatePoolOptions) SetNotificationChannel(notificationChannel string) *UpdatePoolOptions {
	_options.NotificationChannel = core.StringPtr(notificationChannel)
	return _options
}

// SetHealthcheckRegion : Allow user to set HealthcheckRegion
func (_options *UpdatePoolOptions) SetHealthcheckRegion(healthcheckRegion string) *UpdatePoolOptions {
	_options.HealthcheckRegion = core.StringPtr(healthcheckRegion)
	return _options
}

// SetHealthcheckSubnets : Allow user to set HealthcheckSubnets
func (_options *UpdatePoolOptions) SetHealthcheckSubnets(healthcheckSubnets []string) *UpdatePoolOptions {
	_options.HealthcheckSubnets = healthcheckSubnets
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdatePoolOptions) SetXCorrelationID(xCorrelationID string) *UpdatePoolOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdatePoolOptions) SetHeaders(param map[string]string) *UpdatePoolOptions {
	options.Headers = param
	return options
}

// UpdateResourceRecordOptions : The UpdateResourceRecord options.
type UpdateResourceRecordOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required,ne="`

	// The unique identifier of a resource record.
	RecordID *string `json:"record_id" validate:"required,ne="`

	// Name of the resource record.
	Name *string `json:"name" validate:"required"`

	// Content of the resource record.
	Rdata ResourceRecordUpdateInputRdataIntf `json:"rdata" validate:"required"`

	// Time to live in second.
	TTL *int64 `json:"ttl,omitempty"`

	// Only used for SRV record.
	Service *string `json:"service,omitempty"`

	// Only used for SRV record.
	Protocol *string `json:"protocol,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateResourceRecordOptions : Instantiate UpdateResourceRecordOptions
func (*DnsSvcsV1) NewUpdateResourceRecordOptions(instanceID string, dnszoneID string, recordID string, name string, rdata ResourceRecordUpdateInputRdataIntf) *UpdateResourceRecordOptions {
	return &UpdateResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		RecordID:   core.StringPtr(recordID),
		Name:       core.StringPtr(name),
		Rdata:      rdata,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateResourceRecordOptions) SetInstanceID(instanceID string) *UpdateResourceRecordOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDnszoneID : Allow user to set DnszoneID
func (_options *UpdateResourceRecordOptions) SetDnszoneID(dnszoneID string) *UpdateResourceRecordOptions {
	_options.DnszoneID = core.StringPtr(dnszoneID)
	return _options
}

// SetRecordID : Allow user to set RecordID
func (_options *UpdateResourceRecordOptions) SetRecordID(recordID string) *UpdateResourceRecordOptions {
	_options.RecordID = core.StringPtr(recordID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateResourceRecordOptions) SetName(name string) *UpdateResourceRecordOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetRdata : Allow user to set Rdata
func (_options *UpdateResourceRecordOptions) SetRdata(rdata ResourceRecordUpdateInputRdataIntf) *UpdateResourceRecordOptions {
	_options.Rdata = rdata
	return _options
}

// SetTTL : Allow user to set TTL
func (_options *UpdateResourceRecordOptions) SetTTL(ttl int64) *UpdateResourceRecordOptions {
	_options.TTL = core.Int64Ptr(ttl)
	return _options
}

// SetService : Allow user to set Service
func (_options *UpdateResourceRecordOptions) SetService(service string) *UpdateResourceRecordOptions {
	_options.Service = core.StringPtr(service)
	return _options
}

// SetProtocol : Allow user to set Protocol
func (_options *UpdateResourceRecordOptions) SetProtocol(protocol string) *UpdateResourceRecordOptions {
	_options.Protocol = core.StringPtr(protocol)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateResourceRecordOptions) SetXCorrelationID(xCorrelationID string) *UpdateResourceRecordOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateResourceRecordOptions) SetHeaders(param map[string]string) *UpdateResourceRecordOptions {
	options.Headers = param
	return options
}

// UpdateSecondaryZoneOptions : The UpdateSecondaryZone options.
type UpdateSecondaryZoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The unique identifier of a custom resolver.
	ResolverID *string `json:"resolver_id" validate:"required,ne="`

	// The unique identifier of a secondary zone.
	SecondaryZoneID *string `json:"secondary_zone_id" validate:"required,ne="`

	// Descriptive text of the secondary zone.
	Description *string `json:"description,omitempty"`

	// Enable/Disable the secondary zone.
	Enabled *bool `json:"enabled,omitempty"`

	// The addresses of DNS servers where the secondary zone data should be transferred from.
	TransferFrom []string `json:"transfer_from,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateSecondaryZoneOptions : Instantiate UpdateSecondaryZoneOptions
func (*DnsSvcsV1) NewUpdateSecondaryZoneOptions(instanceID string, resolverID string, secondaryZoneID string) *UpdateSecondaryZoneOptions {
	return &UpdateSecondaryZoneOptions{
		InstanceID:      core.StringPtr(instanceID),
		ResolverID:      core.StringPtr(resolverID),
		SecondaryZoneID: core.StringPtr(secondaryZoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateSecondaryZoneOptions) SetInstanceID(instanceID string) *UpdateSecondaryZoneOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetResolverID : Allow user to set ResolverID
func (_options *UpdateSecondaryZoneOptions) SetResolverID(resolverID string) *UpdateSecondaryZoneOptions {
	_options.ResolverID = core.StringPtr(resolverID)
	return _options
}

// SetSecondaryZoneID : Allow user to set SecondaryZoneID
func (_options *UpdateSecondaryZoneOptions) SetSecondaryZoneID(secondaryZoneID string) *UpdateSecondaryZoneOptions {
	_options.SecondaryZoneID = core.StringPtr(secondaryZoneID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateSecondaryZoneOptions) SetDescription(description string) *UpdateSecondaryZoneOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateSecondaryZoneOptions) SetEnabled(enabled bool) *UpdateSecondaryZoneOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetTransferFrom : Allow user to set TransferFrom
func (_options *UpdateSecondaryZoneOptions) SetTransferFrom(transferFrom []string) *UpdateSecondaryZoneOptions {
	_options.TransferFrom = transferFrom
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateSecondaryZoneOptions) SetXCorrelationID(xCorrelationID string) *UpdateSecondaryZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecondaryZoneOptions) SetHeaders(param map[string]string) *UpdateSecondaryZoneOptions {
	options.Headers = param
	return options
}

// AccessRequest : Access request.
type AccessRequest struct {
	// Access request ID.
	ID *string `json:"id" validate:"required"`

	// The information of requestor.
	Requestor *AccessRequestRequestor `json:"requestor" validate:"required"`

	// The zone ID that requestor requests access for.
	ZoneID *string `json:"zone_id" validate:"required"`

	// The zone name that requestor requests access for.
	ZoneName *string `json:"zone_name" validate:"required"`

	// The state of the access request.
	State *string `json:"state" validate:"required"`

	// The expired time of access request with state `pending`.
	PendingExpiresAt *strfmt.DateTime `json:"pending_expires_at,omitempty"`

	// The time when the linked zone is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when the linked zone is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// Constants associated with the AccessRequest.State property.
// The state of the access request.
const (
	AccessRequest_State_Approved = "APPROVED"
	AccessRequest_State_Pending  = "PENDING"
	AccessRequest_State_Rejected = "REJECTED"
	AccessRequest_State_Revoked  = "REVOKED"
	AccessRequest_State_Timedout = "TIMEDOUT"
)

// UnmarshalAccessRequest unmarshals an instance of AccessRequest from the specified map of raw messages.
func UnmarshalAccessRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessRequest)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "requestor", &obj.Requestor, UnmarshalAccessRequestRequestor)
	if err != nil {
		err = core.SDKErrorf(err, "", "requestor-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zone_id", &obj.ZoneID)
	if err != nil {
		err = core.SDKErrorf(err, "", "zone_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zone_name", &obj.ZoneName)
	if err != nil {
		err = core.SDKErrorf(err, "", "zone_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pending_expires_at", &obj.PendingExpiresAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "pending_expires_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccessRequestsList : The list of access requests.
type AccessRequestsList struct {
	// The list of access requests.
	AccessRequests []AccessRequest `json:"access_requests" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first,omitempty"`

	// href.
	Last *PaginationRef `json:"last,omitempty"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalAccessRequestsList unmarshals an instance of AccessRequestsList from the specified map of raw messages.
func UnmarshalAccessRequestsList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessRequestsList)
	err = core.UnmarshalModel(m, "access_requests", &obj.AccessRequests, UnmarshalAccessRequest)
	if err != nil {
		err = core.SDKErrorf(err, "", "access_requests-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AccessRequestsList) GetNextOffset() (*int64, error) {
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

// ActiveDelta : When `state` is `updating`, the current values of properties that are being updated.
type ActiveDelta struct {
	// The profile name of a custom resolver.
	Profile *string `json:"profile,omitempty"`
}

// Constants associated with the ActiveDelta.Profile property.
// The profile name of a custom resolver.
const (
	ActiveDelta_Profile_Advanced  = "advanced"
	ActiveDelta_Profile_Essential = "essential"
	ActiveDelta_Profile_Premier   = "premier"
)

// UnmarshalActiveDelta unmarshals an instance of ActiveDelta from the specified map of raw messages.
func UnmarshalActiveDelta(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActiveDelta)
	err = core.UnmarshalPrimitive(m, "profile", &obj.Profile)
	if err != nil {
		err = core.SDKErrorf(err, "", "profile-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CustomResolver : custom resolver details.
type CustomResolver struct {
	// Identifier of the custom resolver.
	ID *string `json:"id,omitempty"`

	// Name of the custom resolver.
	Name *string `json:"name,omitempty"`

	// Descriptive text of the custom resolver.
	Description *string `json:"description,omitempty"`

	// Whether the custom resolver is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Healthy state of the custom resolver.
	Health *string `json:"health,omitempty"`

	// Locations on which the custom resolver will be running.
	Locations []Location `json:"locations,omitempty"`

	// The profile name of a custom resolver.
	Profile *string `json:"profile,omitempty"`

	// Whether a disruptive update is allowed for the custom resolver.
	AllowDisruptiveUpdates *bool `json:"allow_disruptive_updates,omitempty"`

	// The lifecycle of a custom resolver.
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`

	// The time when a custom resolver is created, in RFC3339 format.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when a custom resolver is modified, in RFC3339 format.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// Constants associated with the CustomResolver.Health property.
// Healthy state of the custom resolver.
const (
	CustomResolver_Health_Critical = "CRITICAL"
	CustomResolver_Health_Degraded = "DEGRADED"
	CustomResolver_Health_Healthy  = "HEALTHY"
)

// Constants associated with the CustomResolver.Profile property.
// The profile name of a custom resolver.
const (
	CustomResolver_Profile_Advanced  = "advanced"
	CustomResolver_Profile_Essential = "essential"
	CustomResolver_Profile_Premier   = "premier"
)

// UnmarshalCustomResolver unmarshals an instance of CustomResolver from the specified map of raw messages.
func UnmarshalCustomResolver(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomResolver)
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
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "health", &obj.Health)
	if err != nil {
		err = core.SDKErrorf(err, "", "health-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "locations", &obj.Locations, UnmarshalLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "locations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "profile", &obj.Profile)
	if err != nil {
		err = core.SDKErrorf(err, "", "profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_disruptive_updates", &obj.AllowDisruptiveUpdates)
	if err != nil {
		err = core.SDKErrorf(err, "", "allow_disruptive_updates-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "lifecycle", &obj.Lifecycle, UnmarshalLifecycle)
	if err != nil {
		err = core.SDKErrorf(err, "", "lifecycle-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CustomResolverList : List custom resolvers response.
type CustomResolverList struct {
	// An array of custom resolvers.
	CustomResolvers []CustomResolver `json:"custom_resolvers,omitempty"`
}

// UnmarshalCustomResolverList unmarshals an instance of CustomResolverList from the specified map of raw messages.
func UnmarshalCustomResolverList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomResolverList)
	err = core.UnmarshalModel(m, "custom_resolvers", &obj.CustomResolvers, UnmarshalCustomResolver)
	if err != nil {
		err = core.SDKErrorf(err, "", "custom_resolvers-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Dnszone : DNS zone details.
type Dnszone struct {
	// Unique identifier of a DNS zone.
	ID *string `json:"id,omitempty"`

	// The time when a DNS zone is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when a DNS zone is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`

	// Unique identifier of a service instance.
	InstanceID *string `json:"instance_id,omitempty"`

	// Name of the DNS zone. Must be a fully qualified domain name.
	Name *string `json:"name,omitempty"`

	// The text describing the purpose of a DNS zone.
	Description *string `json:"description,omitempty"`

	// State of DNS zone.
	State *string `json:"state,omitempty"`

	// The label of a DNS zone.
	Label *string `json:"label,omitempty"`
}

// Constants associated with the Dnszone.State property.
// State of DNS zone.
const (
	Dnszone_State_Active            = "active"
	Dnszone_State_Deleted           = "deleted"
	Dnszone_State_Disabled          = "disabled"
	Dnszone_State_PendingDelete     = "pending_delete"
	Dnszone_State_PendingNetworkAdd = "pending_network_add"
)

// UnmarshalDnszone unmarshals an instance of Dnszone from the specified map of raw messages.
func UnmarshalDnszone(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Dnszone)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_id-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		err = core.SDKErrorf(err, "", "label-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ForwardingRule : forwarding rule details.
type ForwardingRule struct {
	// Identifier of the forwarding rule.
	ID *string `json:"id,omitempty"`

	// Descriptive text of the forwarding rule.
	Description *string `json:"description,omitempty"`

	// Type of the forwarding rule.
	Type *string `json:"type,omitempty"`

	// The matching zone or hostname.
	Match *string `json:"match,omitempty"`

	// The upstream DNS servers that the DNS queries will be forwarded to.
	ForwardTo []string `json:"forward_to,omitempty"`

	// An array of views.
	Views []ViewConfig `json:"views,omitempty"`

	// The time when a forwarding rule is created, in RFC3339 format.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when a forwarding rule is modified, in RFC3339 format.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// Constants associated with the ForwardingRule.Type property.
// Type of the forwarding rule.
const (
	ForwardingRule_Type_Default = "default"
	ForwardingRule_Type_Zone    = "zone"
)

// UnmarshalForwardingRule unmarshals an instance of ForwardingRule from the specified map of raw messages.
func UnmarshalForwardingRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ForwardingRule)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "match", &obj.Match)
	if err != nil {
		err = core.SDKErrorf(err, "", "match-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "forward_to", &obj.ForwardTo)
	if err != nil {
		err = core.SDKErrorf(err, "", "forward_to-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "views", &obj.Views, UnmarshalViewConfig)
	if err != nil {
		err = core.SDKErrorf(err, "", "views-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ForwardingRuleInput : forwarding rule details.
// Models which "extend" this model:
// - ForwardingRuleInputForwardingRuleBoth
// - ForwardingRuleInputForwardingRuleOnlyView
// - ForwardingRuleInputForwardingRuleOnlyForward
type ForwardingRuleInput struct {
	// Descriptive text of the forwarding rule.
	Description *string `json:"description,omitempty"`

	// Type of the forwarding rule.
	Type *string `json:"type" validate:"required"`

	// The matching zone or hostname.
	Match *string `json:"match" validate:"required"`

	// The upstream DNS servers that the DNS queries will be forwarded to.
	ForwardTo []string `json:"forward_to,omitempty"`

	// An array of views.
	Views []ViewConfig `json:"views,omitempty"`
}

// Constants associated with the ForwardingRuleInput.Type property.
// Type of the forwarding rule.
const (
	ForwardingRuleInput_Type_Zone = "zone"
)

func (*ForwardingRuleInput) isaForwardingRuleInput() bool {
	return true
}

type ForwardingRuleInputIntf interface {
	isaForwardingRuleInput() bool
}

// UnmarshalForwardingRuleInput unmarshals an instance of ForwardingRuleInput from the specified map of raw messages.
func UnmarshalForwardingRuleInput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ForwardingRuleInput)
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
	err = core.UnmarshalPrimitive(m, "match", &obj.Match)
	if err != nil {
		err = core.SDKErrorf(err, "", "match-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "forward_to", &obj.ForwardTo)
	if err != nil {
		err = core.SDKErrorf(err, "", "forward_to-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "views", &obj.Views, UnmarshalViewConfig)
	if err != nil {
		err = core.SDKErrorf(err, "", "views-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ForwardingRuleList : List of forwarding rules.
type ForwardingRuleList struct {
	// An array of forwarding rules.
	ForwardingRules []ForwardingRule `json:"forwarding_rules" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first" validate:"required"`

	// href.
	Last *PaginationRef `json:"last" validate:"required"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalForwardingRuleList unmarshals an instance of ForwardingRuleList from the specified map of raw messages.
func UnmarshalForwardingRuleList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ForwardingRuleList)
	err = core.UnmarshalModel(m, "forwarding_rules", &obj.ForwardingRules, UnmarshalForwardingRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "forwarding_rules-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ForwardingRuleList) GetNextOffset() (*int64, error) {
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

// HealthcheckHeader : The HTTP header of health check request.
type HealthcheckHeader struct {
	// The name of HTTP request header.
	Name *string `json:"name" validate:"required"`

	// The value of HTTP request header.
	Value []string `json:"value" validate:"required"`
}

// NewHealthcheckHeader : Instantiate HealthcheckHeader (Generic Model Constructor)
func (*DnsSvcsV1) NewHealthcheckHeader(name string, value []string) (_model *HealthcheckHeader, err error) {
	_model = &HealthcheckHeader{
		Name:  core.StringPtr(name),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalHealthcheckHeader unmarshals an instance of HealthcheckHeader from the specified map of raw messages.
func UnmarshalHealthcheckHeader(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HealthcheckHeader)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportResourceRecordsResp : Import DNS records response.
type ImportResourceRecordsResp struct {
	// Number of records parsed from the zone file.
	TotalRecordsParsed *int64 `json:"total_records_parsed" validate:"required"`

	// Number of records imported successfully.
	RecordsAdded *int64 `json:"records_added" validate:"required"`

	// Number of records failed import.
	RecordsFailed *int64 `json:"records_failed" validate:"required"`

	// Number of records classified by type.
	RecordsAddedByType *RecordStatsByType `json:"records_added_by_type" validate:"required"`

	// Number of records classified by type.
	RecordsFailedByType *RecordStatsByType `json:"records_failed_by_type" validate:"required"`

	// Error messages.
	Messages []RecordsImportMessageModel `json:"messages,omitempty"`

	// Number of records parsed from the zone file.
	Errors []RecordsImportErrorModel `json:"errors,omitempty"`
}

// UnmarshalImportResourceRecordsResp unmarshals an instance of ImportResourceRecordsResp from the specified map of raw messages.
func UnmarshalImportResourceRecordsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportResourceRecordsResp)
	err = core.UnmarshalPrimitive(m, "total_records_parsed", &obj.TotalRecordsParsed)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_records_parsed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "records_added", &obj.RecordsAdded)
	if err != nil {
		err = core.SDKErrorf(err, "", "records_added-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "records_failed", &obj.RecordsFailed)
	if err != nil {
		err = core.SDKErrorf(err, "", "records_failed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "records_added_by_type", &obj.RecordsAddedByType, UnmarshalRecordStatsByType)
	if err != nil {
		err = core.SDKErrorf(err, "", "records_added_by_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "records_failed_by_type", &obj.RecordsFailedByType, UnmarshalRecordStatsByType)
	if err != nil {
		err = core.SDKErrorf(err, "", "records_failed_by_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalRecordsImportMessageModel)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalRecordsImportErrorModel)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Lifecycle : The lifecycle of a custom resolver.
type Lifecycle struct {
	// The lifecycle state of a custom resolver.
	State *string `json:"state" validate:"required"`

	// When `state` is `updating`, the current values of properties that are being updated.
	ActiveDelta *ActiveDelta `json:"active_delta,omitempty"`

	// The reasons why a lifecycle state is not stable.
	Reasons []LifecycleReason `json:"reasons" validate:"required"`
}

// Constants associated with the Lifecycle.State property.
// The lifecycle state of a custom resolver.
const (
	Lifecycle_State_Failed   = "failed"
	Lifecycle_State_Pending  = "pending"
	Lifecycle_State_Stable   = "stable"
	Lifecycle_State_Updating = "updating"
)

// UnmarshalLifecycle unmarshals an instance of Lifecycle from the specified map of raw messages.
func UnmarshalLifecycle(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Lifecycle)
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "active_delta", &obj.ActiveDelta, UnmarshalActiveDelta)
	if err != nil {
		err = core.SDKErrorf(err, "", "active_delta-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "reasons", &obj.Reasons, UnmarshalLifecycleReason)
	if err != nil {
		err = core.SDKErrorf(err, "", "reasons-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LifecycleReason : LifecycleReason struct
type LifecycleReason struct {
	// A reason code for a lifecycle state that is not stable.
	Code *string `json:"code,omitempty"`

	// An explanation of the reason for this lifecycle state.
	Message *string `json:"message,omitempty"`

	// Link to documentation about the reason for this lifecycle state.
	MoreInfo *string `json:"more_info,omitempty"`
}

// Constants associated with the LifecycleReason.Code property.
// A reason code for a lifecycle state that is not stable.
const (
	LifecycleReason_Code_InternalError = "internal_error"
)

// UnmarshalLifecycleReason unmarshals an instance of LifecycleReason from the specified map of raw messages.
func UnmarshalLifecycleReason(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LifecycleReason)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		err = core.SDKErrorf(err, "", "code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "more_info", &obj.MoreInfo)
	if err != nil {
		err = core.SDKErrorf(err, "", "more_info-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LinkedDnszone : linked zone details.
type LinkedDnszone struct {
	// Identifier of the linked zone.
	ID *string `json:"id" validate:"required"`

	// Unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// Name of owner's DNS zone.
	Name *string `json:"name" validate:"required"`

	// Descriptive text of the linked zone.
	Description *string `json:"description,omitempty"`

	// The owner's instance and zone that the zone is linked to.
	LinkedTo *LinkedDnszoneLinkedTo `json:"linked_to" validate:"required"`

	// The state of linked zone.
	State *string `json:"state" validate:"required"`

	// The label of linked zone.
	Label *string `json:"label,omitempty"`

	// The expired time of linked zone with state `approval pending`.
	ApprovalRequiredBefore *strfmt.DateTime `json:"approval_required_before,omitempty"`

	// The time when the linked zone is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when the linked zone is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// Constants associated with the LinkedDnszone.State property.
// The state of linked zone.
const (
	LinkedDnszone_State_Active            = "ACTIVE"
	LinkedDnszone_State_ApprovalRejected  = "APPROVAL_REJECTED"
	LinkedDnszone_State_ApprovalRevoked   = "APPROVAL_REVOKED"
	LinkedDnszone_State_ApprovalTimedout  = "APPROVAL_TIMEDOUT"
	LinkedDnszone_State_PendingApproval   = "PENDING_APPROVAL"
	LinkedDnszone_State_PendingNetworkAdd = "PENDING_NETWORK_ADD"
)

// UnmarshalLinkedDnszone unmarshals an instance of LinkedDnszone from the specified map of raw messages.
func UnmarshalLinkedDnszone(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LinkedDnszone)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_id-error", common.GetComponentInfo())
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
	err = core.UnmarshalModel(m, "linked_to", &obj.LinkedTo, UnmarshalLinkedDnszoneLinkedTo)
	if err != nil {
		err = core.SDKErrorf(err, "", "linked_to-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		err = core.SDKErrorf(err, "", "label-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_required_before", &obj.ApprovalRequiredBefore)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_required_before-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LinkedDnszonesList : The list of linked zones.
type LinkedDnszonesList struct {
	// The list of linked zones.
	LinkedDnszones []LinkedDnszone `json:"linked_dnszones" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first,omitempty"`

	// href.
	Last *PaginationRef `json:"last,omitempty"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalLinkedDnszonesList unmarshals an instance of LinkedDnszonesList from the specified map of raw messages.
func UnmarshalLinkedDnszonesList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LinkedDnszonesList)
	err = core.UnmarshalModel(m, "linked_dnszones", &obj.LinkedDnszones, UnmarshalLinkedDnszone)
	if err != nil {
		err = core.SDKErrorf(err, "", "linked_dnszones-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *LinkedDnszonesList) GetNextOffset() (*int64, error) {
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

// ListDnszones : List DNS zones response.
type ListDnszones struct {
	// An array of DNS zones.
	Dnszones []Dnszone `json:"dnszones" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first" validate:"required"`

	// href.
	Last *PaginationRef `json:"last" validate:"required"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalListDnszones unmarshals an instance of ListDnszones from the specified map of raw messages.
func UnmarshalListDnszones(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListDnszones)
	err = core.UnmarshalModel(m, "dnszones", &obj.Dnszones, UnmarshalDnszone)
	if err != nil {
		err = core.SDKErrorf(err, "", "dnszones-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListDnszones) GetNextOffset() (*int64, error) {
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

// ListLoadBalancers : List Global Load Balancers response.
type ListLoadBalancers struct {
	// An array of Global Load Balancers.
	LoadBalancers []LoadBalancer `json:"load_balancers" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first" validate:"required"`

	// href.
	Last *PaginationRef `json:"last" validate:"required"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalListLoadBalancers unmarshals an instance of ListLoadBalancers from the specified map of raw messages.
func UnmarshalListLoadBalancers(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListLoadBalancers)
	err = core.UnmarshalModel(m, "load_balancers", &obj.LoadBalancers, UnmarshalLoadBalancer)
	if err != nil {
		err = core.SDKErrorf(err, "", "load_balancers-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListLoadBalancers) GetNextOffset() (*int64, error) {
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

// ListMonitors : List load balancer monitors response.
type ListMonitors struct {
	// An array of load balancer monitors.
	Monitors []Monitor `json:"monitors" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first" validate:"required"`

	// href.
	Last *PaginationRef `json:"last" validate:"required"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalListMonitors unmarshals an instance of ListMonitors from the specified map of raw messages.
func UnmarshalListMonitors(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListMonitors)
	err = core.UnmarshalModel(m, "monitors", &obj.Monitors, UnmarshalMonitor)
	if err != nil {
		err = core.SDKErrorf(err, "", "monitors-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListMonitors) GetNextOffset() (*int64, error) {
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

// ListPermittedNetworks : List permitted networks response.
type ListPermittedNetworks struct {
	// An array of permitted networks.
	PermittedNetworks []PermittedNetwork `json:"permitted_networks" validate:"required"`
}

// UnmarshalListPermittedNetworks unmarshals an instance of ListPermittedNetworks from the specified map of raw messages.
func UnmarshalListPermittedNetworks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListPermittedNetworks)
	err = core.UnmarshalModel(m, "permitted_networks", &obj.PermittedNetworks, UnmarshalPermittedNetwork)
	if err != nil {
		err = core.SDKErrorf(err, "", "permitted_networks-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListPools : List load balancer pools response.
type ListPools struct {
	// An array of load balancer pools.
	Pools []Pool `json:"pools" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first" validate:"required"`

	// href.
	Last *PaginationRef `json:"last" validate:"required"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalListPools unmarshals an instance of ListPools from the specified map of raw messages.
func UnmarshalListPools(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListPools)
	err = core.UnmarshalModel(m, "pools", &obj.Pools, UnmarshalPool)
	if err != nil {
		err = core.SDKErrorf(err, "", "pools-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListPools) GetNextOffset() (*int64, error) {
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

// ListResourceRecords : List Resource Records response.
type ListResourceRecords struct {
	// An array of resource records.
	ResourceRecords []ResourceRecord `json:"resource_records" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first" validate:"required"`

	// href.
	Last *PaginationRef `json:"last" validate:"required"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalListResourceRecords unmarshals an instance of ListResourceRecords from the specified map of raw messages.
func UnmarshalListResourceRecords(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListResourceRecords)
	err = core.UnmarshalModel(m, "resource_records", &obj.ResourceRecords, UnmarshalResourceRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_records-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListResourceRecords) GetNextOffset() (*int64, error) {
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

// LoadBalancer : Load balancer details.
type LoadBalancer struct {
	// Identifier of the load balancer.
	ID *string `json:"id,omitempty"`

	// Name of the load balancer.
	Name *string `json:"name,omitempty"`

	// Descriptive text of the load balancer.
	Description *string `json:"description,omitempty"`

	// Whether the load balancer is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Time to live in second.
	TTL *int64 `json:"ttl,omitempty"`

	// Healthy state of the load balancer.
	Health *string `json:"health,omitempty"`

	// The pool ID to use when all other pools are detected as unhealthy.
	FallbackPool *string `json:"fallback_pool,omitempty"`

	// A list of pool IDs ordered by their failover priority. Pools defined here are used by default, or when region_pools
	// are not configured for a given region.
	DefaultPools []string `json:"default_pools,omitempty"`

	// Map availability zones to pool IDs.
	AzPools []AzPoolsItem `json:"az_pools,omitempty"`

	// The time when a load balancer is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when a load balancer is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// Constants associated with the LoadBalancer.Health property.
// Healthy state of the load balancer.
const (
	LoadBalancer_Health_Critical = "CRITICAL"
	LoadBalancer_Health_Degraded = "DEGRADED"
	LoadBalancer_Health_Healthy  = "HEALTHY"
)

// UnmarshalLoadBalancer unmarshals an instance of LoadBalancer from the specified map of raw messages.
func UnmarshalLoadBalancer(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancer)
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
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		err = core.SDKErrorf(err, "", "ttl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "health", &obj.Health)
	if err != nil {
		err = core.SDKErrorf(err, "", "health-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "fallback_pool", &obj.FallbackPool)
	if err != nil {
		err = core.SDKErrorf(err, "", "fallback_pool-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "default_pools", &obj.DefaultPools)
	if err != nil {
		err = core.SDKErrorf(err, "", "default_pools-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "az_pools", &obj.AzPools, UnmarshalAzPoolsItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "az_pools-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Location : Custom resolver location.
type Location struct {
	// Location ID.
	ID *string `json:"id,omitempty"`

	// Subnet CRN.
	SubnetCrn *string `json:"subnet_crn,omitempty"`

	// Whether the location is enabled for the custom resolver.
	Enabled *bool `json:"enabled,omitempty"`

	// Whether the DNS server in this location is healthy or not.
	Healthy *bool `json:"healthy,omitempty"`

	// The ip address of this dns server.
	DnsServerIp *string `json:"dns_server_ip,omitempty"`
}

// UnmarshalLocation unmarshals an instance of Location from the specified map of raw messages.
func UnmarshalLocation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Location)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "subnet_crn", &obj.SubnetCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "subnet_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "healthy", &obj.Healthy)
	if err != nil {
		err = core.SDKErrorf(err, "", "healthy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dns_server_ip", &obj.DnsServerIp)
	if err != nil {
		err = core.SDKErrorf(err, "", "dns_server_ip-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LocationInput : Request to add custom resolver location.
type LocationInput struct {
	// Custom resolver location, subnet CRN.
	SubnetCrn *string `json:"subnet_crn" validate:"required"`

	// Enable/Disable custom resolver location.
	Enabled *bool `json:"enabled,omitempty"`
}

// NewLocationInput : Instantiate LocationInput (Generic Model Constructor)
func (*DnsSvcsV1) NewLocationInput(subnetCrn string) (_model *LocationInput, err error) {
	_model = &LocationInput{
		SubnetCrn: core.StringPtr(subnetCrn),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalLocationInput unmarshals an instance of LocationInput from the specified map of raw messages.
func UnmarshalLocationInput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LocationInput)
	err = core.UnmarshalPrimitive(m, "subnet_crn", &obj.SubnetCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "subnet_crn-error", common.GetComponentInfo())
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

// Monitor : Load balancer monitor details.
type Monitor struct {
	// Identifier of the load balancer monitor.
	ID *string `json:"id,omitempty"`

	// The name of the load balancer monitor.
	Name *string `json:"name,omitempty"`

	// Descriptive text of the load balancer monitor.
	Description *string `json:"description,omitempty"`

	// The protocol to use for the health check. Currently supported protocols are 'HTTP','HTTPS' and 'TCP'.
	Type *string `json:"type,omitempty"`

	// Port number to connect to for the health check. Required for TCP checks. HTTP and HTTPS checks should only define
	// the port when using a non-standard port (HTTP: default 80, HTTPS: default 443).
	Port *int64 `json:"port,omitempty"`

	// The interval between each health check. Shorter intervals may improve failover time, but will increase load on the
	// origins as we check from multiple locations.
	Interval *int64 `json:"interval,omitempty"`

	// The number of retries to attempt in case of a timeout before marking the origin as unhealthy. Retries are attempted
	// immediately.
	Retries *int64 `json:"retries,omitempty"`

	// The timeout (in seconds) before marking the health check as failed.
	Timeout *int64 `json:"timeout,omitempty"`

	// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
	Method *string `json:"method,omitempty"`

	// The endpoint path to health check against. This parameter is only valid for HTTP and HTTPS monitors.
	Path *string `json:"path,omitempty"`

	// The HTTP request headers to send in the health check. It is recommended you set a Host header by default. The
	// User-Agent header cannot be overridden. This parameter is only valid for HTTP and HTTPS monitors.
	HeadersVar []HealthcheckHeader `json:"headers,omitempty"`

	// Do not validate the certificate when monitor use HTTPS. This parameter is currently only valid for HTTPS monitors.
	AllowInsecure *bool `json:"allow_insecure,omitempty"`

	// The expected HTTP response code or code range of the health check. This parameter is only valid for HTTP and HTTPS
	// monitors.
	ExpectedCodes *string `json:"expected_codes,omitempty"`

	// A case-insensitive sub-string to look for in the response body. If this string is not found, the origin will be
	// marked as unhealthy. This parameter is only valid for HTTP and HTTPS monitors.
	ExpectedBody *string `json:"expected_body,omitempty"`

	// the time when a load balancer monitor is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// the recent time when a load balancer monitor is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// Constants associated with the Monitor.Method property.
// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
const (
	Monitor_Method_Get  = "GET"
	Monitor_Method_Head = "HEAD"
)

// UnmarshalMonitor unmarshals an instance of Monitor from the specified map of raw messages.
func UnmarshalMonitor(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Monitor)
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
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		err = core.SDKErrorf(err, "", "interval-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "retries", &obj.Retries)
	if err != nil {
		err = core.SDKErrorf(err, "", "retries-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "timeout-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		err = core.SDKErrorf(err, "", "method-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		err = core.SDKErrorf(err, "", "path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "headers", &obj.HeadersVar, UnmarshalHealthcheckHeader)
	if err != nil {
		err = core.SDKErrorf(err, "", "headers-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_insecure", &obj.AllowInsecure)
	if err != nil {
		err = core.SDKErrorf(err, "", "allow_insecure-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_codes", &obj.ExpectedCodes)
	if err != nil {
		err = core.SDKErrorf(err, "", "expected_codes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_body", &obj.ExpectedBody)
	if err != nil {
		err = core.SDKErrorf(err, "", "expected_body-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Origin : Origin server.
type Origin struct {
	// The name of the origin server.
	Name *string `json:"name,omitempty"`

	// Description of the origin server.
	Description *string `json:"description,omitempty"`

	// The address of the origin server. It can be a hostname or an IP address.
	Address *string `json:"address,omitempty"`

	// Whether the origin server is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The health state of the origin server.
	Health *bool `json:"health,omitempty"`

	// The failure reason of the origin server if it is unhealthy.
	HealthFailureReason *string `json:"health_failure_reason,omitempty"`
}

// UnmarshalOrigin unmarshals an instance of Origin from the specified map of raw messages.
func UnmarshalOrigin(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Origin)
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
	err = core.UnmarshalPrimitive(m, "address", &obj.Address)
	if err != nil {
		err = core.SDKErrorf(err, "", "address-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "health", &obj.Health)
	if err != nil {
		err = core.SDKErrorf(err, "", "health-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "health_failure_reason", &obj.HealthFailureReason)
	if err != nil {
		err = core.SDKErrorf(err, "", "health_failure_reason-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OriginInput : The request data of origin server.
type OriginInput struct {
	// The name of the origin server.
	Name *string `json:"name,omitempty"`

	// Description of the origin server.
	Description *string `json:"description,omitempty"`

	// The address of the origin server. It can be a hostname or an IP address.
	Address *string `json:"address,omitempty"`

	// Whether the origin server is enabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// UnmarshalOriginInput unmarshals an instance of OriginInput from the specified map of raw messages.
func UnmarshalOriginInput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OriginInput)
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
	err = core.UnmarshalPrimitive(m, "address", &obj.Address)
	if err != nil {
		err = core.SDKErrorf(err, "", "address-error", common.GetComponentInfo())
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

// PaginationRef : href.
type PaginationRef struct {
	// href.
	Href *string `json:"href,omitempty"`
}

// UnmarshalPaginationRef unmarshals an instance of PaginationRef from the specified map of raw messages.
func UnmarshalPaginationRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginationRef)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PermittedNetwork : Permitted network details.
type PermittedNetwork struct {
	// Unique identifier of a permitted network.
	ID *string `json:"id,omitempty"`

	// The time when a permitted network is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when a permitted network is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`

	// Permitted network data for VPC.
	PermittedNetwork *PermittedNetworkVpc `json:"permitted_network,omitempty"`

	// The type of a permitted network.
	Type *string `json:"type,omitempty"`

	// The state of a permitted network.
	State *string `json:"state,omitempty"`

	// Unique identifier of a linked zone through which the permitted network was added.
	LinkedZoneID *string `json:"linked_zone_id,omitempty"`
}

// Constants associated with the PermittedNetwork.Type property.
// The type of a permitted network.
const (
	PermittedNetwork_Type_Vpc = "vpc"
)

// Constants associated with the PermittedNetwork.State property.
// The state of a permitted network.
const (
	PermittedNetwork_State_Active            = "ACTIVE"
	PermittedNetwork_State_RemovalInProgress = "REMOVAL_IN_PROGRESS"
)

// UnmarshalPermittedNetwork unmarshals an instance of PermittedNetwork from the specified map of raw messages.
func UnmarshalPermittedNetwork(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PermittedNetwork)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "permitted_network", &obj.PermittedNetwork, UnmarshalPermittedNetworkVpc)
	if err != nil {
		err = core.SDKErrorf(err, "", "permitted_network-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "linked_zone_id", &obj.LinkedZoneID)
	if err != nil {
		err = core.SDKErrorf(err, "", "linked_zone_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PermittedNetworkVpc : Permitted network data for VPC.
type PermittedNetworkVpc struct {
	// CRN string uniquely identifies a VPC.
	VpcCrn *string `json:"vpc_crn" validate:"required"`
}

// NewPermittedNetworkVpc : Instantiate PermittedNetworkVpc (Generic Model Constructor)
func (*DnsSvcsV1) NewPermittedNetworkVpc(vpcCrn string) (_model *PermittedNetworkVpc, err error) {
	_model = &PermittedNetworkVpc{
		VpcCrn: core.StringPtr(vpcCrn),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalPermittedNetworkVpc unmarshals an instance of PermittedNetworkVpc from the specified map of raw messages.
func UnmarshalPermittedNetworkVpc(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PermittedNetworkVpc)
	err = core.UnmarshalPrimitive(m, "vpc_crn", &obj.VpcCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "vpc_crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Pool : Load balancer pool details.
type Pool struct {
	// Identifier of the load balancer pool.
	ID *string `json:"id,omitempty"`

	// Name of the load balancer pool.
	Name *string `json:"name,omitempty"`

	// Descriptive text of the load balancer pool.
	Description *string `json:"description,omitempty"`

	// Whether the load balancer pool is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The minimum number of origins that must be healthy for this pool to serve traffic. If the number of healthy origins
	// falls below this number, the pool will be marked unhealthy and we will failover to the next available pool.
	HealthyOriginsThreshold *int64 `json:"healthy_origins_threshold,omitempty"`

	// The list of origins within this pool. Traffic directed at this pool is balanced across all currently healthy
	// origins, provided the pool itself is healthy.
	Origins []Origin `json:"origins,omitempty"`

	// The ID of the load balancer monitor to be associated to this pool.
	Monitor *string `json:"monitor,omitempty"`

	// The notification channel.
	NotificationChannel *string `json:"notification_channel,omitempty"`

	// Healthy state of the load balancer pool.
	Health *string `json:"health,omitempty"`

	// Health check region of VSIs.
	HealthcheckRegion *string `json:"healthcheck_region,omitempty"`

	// Health check subnet CRNs.
	HealthcheckSubnets []string `json:"healthcheck_subnets,omitempty"`

	// Health check VSI information.
	HealthcheckVsis []PoolHealthcheckVsisItem `json:"healthcheck_vsis,omitempty"`

	// the time when a load balancer pool is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// the recent time when a load balancer pool is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// Constants associated with the Pool.Health property.
// Healthy state of the load balancer pool.
const (
	Pool_Health_Critical = "CRITICAL"
	Pool_Health_Degraded = "DEGRADED"
	Pool_Health_Healthy  = "HEALTHY"
)

// Constants associated with the Pool.HealthcheckRegion property.
// Health check region of VSIs.
const (
	Pool_HealthcheckRegion_AuSyd   = "au-syd"
	Pool_HealthcheckRegion_EuDu    = "eu-du"
	Pool_HealthcheckRegion_EuGb    = "eu-gb"
	Pool_HealthcheckRegion_JpTok   = "jp-tok"
	Pool_HealthcheckRegion_UsEast  = "us-east"
	Pool_HealthcheckRegion_UsSouth = "us-south"
)

// UnmarshalPool unmarshals an instance of Pool from the specified map of raw messages.
func UnmarshalPool(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Pool)
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
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "healthy_origins_threshold", &obj.HealthyOriginsThreshold)
	if err != nil {
		err = core.SDKErrorf(err, "", "healthy_origins_threshold-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "origins", &obj.Origins, UnmarshalOrigin)
	if err != nil {
		err = core.SDKErrorf(err, "", "origins-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "monitor", &obj.Monitor)
	if err != nil {
		err = core.SDKErrorf(err, "", "monitor-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "notification_channel", &obj.NotificationChannel)
	if err != nil {
		err = core.SDKErrorf(err, "", "notification_channel-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "health", &obj.Health)
	if err != nil {
		err = core.SDKErrorf(err, "", "health-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "healthcheck_region", &obj.HealthcheckRegion)
	if err != nil {
		err = core.SDKErrorf(err, "", "healthcheck_region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "healthcheck_subnets", &obj.HealthcheckSubnets)
	if err != nil {
		err = core.SDKErrorf(err, "", "healthcheck_subnets-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "healthcheck_vsis", &obj.HealthcheckVsis, UnmarshalPoolHealthcheckVsisItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "healthcheck_vsis-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RecordStatsByType : Number of records classified by type.
type RecordStatsByType struct {
	// Number of records, type A.
	A *int64 `json:"A" validate:"required"`

	// Number of records, type AAAA.
	AAAA *int64 `json:"AAAA" validate:"required"`

	// Number of records, type CNAME.
	CNAME *int64 `json:"CNAME" validate:"required"`

	// Number of records, type SRV.
	SRV *int64 `json:"SRV" validate:"required"`

	// Number of records, type TXT.
	TXT *int64 `json:"TXT" validate:"required"`

	// Number of records, type MX.
	MX *int64 `json:"MX" validate:"required"`

	// Number of records, type PTR.
	PTR *int64 `json:"PTR" validate:"required"`
}

// UnmarshalRecordStatsByType unmarshals an instance of RecordStatsByType from the specified map of raw messages.
func UnmarshalRecordStatsByType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RecordStatsByType)
	err = core.UnmarshalPrimitive(m, "A", &obj.A)
	if err != nil {
		err = core.SDKErrorf(err, "", "A-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AAAA", &obj.AAAA)
	if err != nil {
		err = core.SDKErrorf(err, "", "AAAA-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CNAME", &obj.CNAME)
	if err != nil {
		err = core.SDKErrorf(err, "", "CNAME-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SRV", &obj.SRV)
	if err != nil {
		err = core.SDKErrorf(err, "", "SRV-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "TXT", &obj.TXT)
	if err != nil {
		err = core.SDKErrorf(err, "", "TXT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MX", &obj.MX)
	if err != nil {
		err = core.SDKErrorf(err, "", "MX-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PTR", &obj.PTR)
	if err != nil {
		err = core.SDKErrorf(err, "", "PTR-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RecordsImportErrorModel : RecordsImportErrorModel struct
type RecordsImportErrorModel struct {
	// resource record content in zone file.
	ResourceRecord *string `json:"resource_record" validate:"required"`

	// Error container.
	Error *RecordsImportErrorModelError `json:"error" validate:"required"`
}

// UnmarshalRecordsImportErrorModel unmarshals an instance of RecordsImportErrorModel from the specified map of raw messages.
func UnmarshalRecordsImportErrorModel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RecordsImportErrorModel)
	err = core.UnmarshalPrimitive(m, "resource_record", &obj.ResourceRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_record-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "error", &obj.Error, UnmarshalRecordsImportErrorModelError)
	if err != nil {
		err = core.SDKErrorf(err, "", "error-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RecordsImportMessageModel : RecordsImportMessageModel struct
type RecordsImportMessageModel struct {
	// Code to classify import DNS records error.
	Code *string `json:"code" validate:"required"`

	// Message to describe import DNS records error.
	Message *string `json:"message" validate:"required"`
}

// UnmarshalRecordsImportMessageModel unmarshals an instance of RecordsImportMessageModel from the specified map of raw messages.
func UnmarshalRecordsImportMessageModel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RecordsImportMessageModel)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		err = core.SDKErrorf(err, "", "code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecord : Resource record details.
type ResourceRecord struct {
	// Identifier of the resource record.
	ID *string `json:"id,omitempty"`

	// The time when a resource record is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when a resource record is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`

	// Name of the resource record.
	Name *string `json:"name,omitempty"`

	// Type of the resource record.
	Type *string `json:"type,omitempty"`

	// Time to live in second.
	TTL *int64 `json:"ttl,omitempty"`

	// Content of the resource record.
	Rdata map[string]interface{} `json:"rdata,omitempty"`

	// Only used for SRV record.
	Service *string `json:"service,omitempty"`

	// Only used for SRV record.
	Protocol *string `json:"protocol,omitempty"`
}

// Constants associated with the ResourceRecord.Type property.
// Type of the resource record.
const (
	ResourceRecord_Type_A     = "A"
	ResourceRecord_Type_Aaaa  = "AAAA"
	ResourceRecord_Type_Cname = "CNAME"
	ResourceRecord_Type_Mx    = "MX"
	ResourceRecord_Type_Ptr   = "PTR"
	ResourceRecord_Type_Srv   = "SRV"
	ResourceRecord_Type_Txt   = "TXT"
)

// UnmarshalResourceRecord unmarshals an instance of ResourceRecord from the specified map of raw messages.
func UnmarshalResourceRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecord)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		err = core.SDKErrorf(err, "", "ttl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rdata", &obj.Rdata)
	if err != nil {
		err = core.SDKErrorf(err, "", "rdata-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service", &obj.Service)
	if err != nil {
		err = core.SDKErrorf(err, "", "service-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "protocol", &obj.Protocol)
	if err != nil {
		err = core.SDKErrorf(err, "", "protocol-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecondaryZone : Secondary zone details.
type SecondaryZone struct {
	// Identifier of the secondary zone.
	ID *string `json:"id" validate:"required"`

	// Descriptive text of the secondary zone.
	Description *string `json:"description,omitempty"`

	// zone name.
	Zone *string `json:"zone" validate:"required"`

	// Enable/Disable the secondary zone.
	Enabled *bool `json:"enabled" validate:"required"`

	// The addresses of DNS servers where the secondary zone data should be transferred from.
	TransferFrom []string `json:"transfer_from" validate:"required"`

	// The time when a secondary zone is created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The recent time when a secondary zone is modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// UnmarshalSecondaryZone unmarshals an instance of SecondaryZone from the specified map of raw messages.
func UnmarshalSecondaryZone(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecondaryZone)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zone", &obj.Zone)
	if err != nil {
		err = core.SDKErrorf(err, "", "zone-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "transfer_from", &obj.TransferFrom)
	if err != nil {
		err = core.SDKErrorf(err, "", "transfer_from-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_on-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecondaryZoneList : List of secondary zones.
type SecondaryZoneList struct {
	// Secondary zones.
	SecondaryZones []SecondaryZone `json:"secondary_zones" validate:"required"`

	// The number of resources to skip over.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources might be returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources are returned.
	Count *int64 `json:"count" validate:"required"`

	// Total number of resources.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *PaginationRef `json:"first,omitempty"`

	// href.
	Last *PaginationRef `json:"last,omitempty"`

	// href.
	Previous *PaginationRef `json:"previous,omitempty"`

	// href.
	Next *PaginationRef `json:"next,omitempty"`
}

// UnmarshalSecondaryZoneList unmarshals an instance of SecondaryZoneList from the specified map of raw messages.
func UnmarshalSecondaryZoneList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecondaryZoneList)
	err = core.UnmarshalModel(m, "secondary_zones", &obj.SecondaryZones, UnmarshalSecondaryZone)
	if err != nil {
		err = core.SDKErrorf(err, "", "secondary_zones-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SecondaryZoneList) GetNextOffset() (*int64, error) {
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

// ViewConfig : View configuration describes the DNS forwarding policy such that DNS queries matching an expression should be
// forwarded to the target DNS servers.
type ViewConfig struct {
	// The view name, which must be unique in a forwarding rule.
	Name *string `json:"name" validate:"required"`

	// Descriptive text of the view.
	Description *string `json:"description,omitempty"`

	// The custom resolver will only apply the rule to the incoming DNS queries if the expression evaluates to true. The
	// expression follows
	// [Common Expression Language](https://github.com/google/cel-spec/blob/master/doc/langdef.md), but does not support
	// CEL build-in functions and macros. Currently, the expression only supports the following custom functions, variables
	// and operators:
	//
	// - Functions:
	//   - `ipInRange(ip, cidr)`: Return boolean value indicating whether
	//   the `ip` address is in the `cidr` range.
	//
	// - Variables:
	//   - `source`: Client information for the DNS query.
	//     - `ip`: The client's IP address.
	//
	// - Operators:
	//   - `||`: Logical OR
	//   - `&&`: Logical AND
	//   - `!`: Logical NOT
	//   - `==`: Logical Equals
	//   - `!=`: Logical NotEquals
	//   - `?:`: Conditional
	//
	// The supported Common Expression Language features can change without advanced notification.
	Expression *string `json:"expression" validate:"required"`

	// The target DNS servers that the matching DNS queries are forwarded to.
	ForwardTo []string `json:"forward_to" validate:"required"`
}

// NewViewConfig : Instantiate ViewConfig (Generic Model Constructor)
func (*DnsSvcsV1) NewViewConfig(name string, expression string, forwardTo []string) (_model *ViewConfig, err error) {
	_model = &ViewConfig{
		Name:       core.StringPtr(name),
		Expression: core.StringPtr(expression),
		ForwardTo:  forwardTo,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalViewConfig unmarshals an instance of ViewConfig from the specified map of raw messages.
func UnmarshalViewConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ViewConfig)
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
	err = core.UnmarshalPrimitive(m, "expression", &obj.Expression)
	if err != nil {
		err = core.SDKErrorf(err, "", "expression-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "forward_to", &obj.ForwardTo)
	if err != nil {
		err = core.SDKErrorf(err, "", "forward_to-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordInputRdataRdataARecord : The content of type-A resource record.
// This model "extends" ResourceRecordInputRdata
type ResourceRecordInputRdataRdataARecord struct {
	// IPv4 address.
	Ip *string `json:"ip" validate:"required"`
}

// NewResourceRecordInputRdataRdataARecord : Instantiate ResourceRecordInputRdataRdataARecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataARecord(ip string) (_model *ResourceRecordInputRdataRdataARecord, err error) {
	_model = &ResourceRecordInputRdataRdataARecord{
		Ip: core.StringPtr(ip),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordInputRdataRdataARecord) isaResourceRecordInputRdata() bool {
	return true
}

// UnmarshalResourceRecordInputRdataRdataARecord unmarshals an instance of ResourceRecordInputRdataRdataARecord from the specified map of raw messages.
func UnmarshalResourceRecordInputRdataRdataARecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordInputRdataRdataARecord)
	err = core.UnmarshalPrimitive(m, "ip", &obj.Ip)
	if err != nil {
		err = core.SDKErrorf(err, "", "ip-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordInputRdataRdataAaaaRecord : The content of type-AAAA resource record.
// This model "extends" ResourceRecordInputRdata
type ResourceRecordInputRdataRdataAaaaRecord struct {
	// IPv6 address.
	Ip *string `json:"ip" validate:"required"`
}

// NewResourceRecordInputRdataRdataAaaaRecord : Instantiate ResourceRecordInputRdataRdataAaaaRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataAaaaRecord(ip string) (_model *ResourceRecordInputRdataRdataAaaaRecord, err error) {
	_model = &ResourceRecordInputRdataRdataAaaaRecord{
		Ip: core.StringPtr(ip),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordInputRdataRdataAaaaRecord) isaResourceRecordInputRdata() bool {
	return true
}

// UnmarshalResourceRecordInputRdataRdataAaaaRecord unmarshals an instance of ResourceRecordInputRdataRdataAaaaRecord from the specified map of raw messages.
func UnmarshalResourceRecordInputRdataRdataAaaaRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordInputRdataRdataAaaaRecord)
	err = core.UnmarshalPrimitive(m, "ip", &obj.Ip)
	if err != nil {
		err = core.SDKErrorf(err, "", "ip-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordInputRdataRdataCnameRecord : The content of type-CNAME resource record.
// This model "extends" ResourceRecordInputRdata
type ResourceRecordInputRdataRdataCnameRecord struct {
	// Canonical name.
	Cname *string `json:"cname" validate:"required"`
}

// NewResourceRecordInputRdataRdataCnameRecord : Instantiate ResourceRecordInputRdataRdataCnameRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataCnameRecord(cname string) (_model *ResourceRecordInputRdataRdataCnameRecord, err error) {
	_model = &ResourceRecordInputRdataRdataCnameRecord{
		Cname: core.StringPtr(cname),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordInputRdataRdataCnameRecord) isaResourceRecordInputRdata() bool {
	return true
}

// UnmarshalResourceRecordInputRdataRdataCnameRecord unmarshals an instance of ResourceRecordInputRdataRdataCnameRecord from the specified map of raw messages.
func UnmarshalResourceRecordInputRdataRdataCnameRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordInputRdataRdataCnameRecord)
	err = core.UnmarshalPrimitive(m, "cname", &obj.Cname)
	if err != nil {
		err = core.SDKErrorf(err, "", "cname-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordInputRdataRdataMxRecord : The content of type-MX resource record.
// This model "extends" ResourceRecordInputRdata
type ResourceRecordInputRdataRdataMxRecord struct {
	// Hostname of Exchange server.
	Exchange *string `json:"exchange" validate:"required"`

	// Preference of the MX record.
	Preference *int64 `json:"preference" validate:"required"`
}

// NewResourceRecordInputRdataRdataMxRecord : Instantiate ResourceRecordInputRdataRdataMxRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataMxRecord(exchange string, preference int64) (_model *ResourceRecordInputRdataRdataMxRecord, err error) {
	_model = &ResourceRecordInputRdataRdataMxRecord{
		Exchange:   core.StringPtr(exchange),
		Preference: core.Int64Ptr(preference),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordInputRdataRdataMxRecord) isaResourceRecordInputRdata() bool {
	return true
}

// UnmarshalResourceRecordInputRdataRdataMxRecord unmarshals an instance of ResourceRecordInputRdataRdataMxRecord from the specified map of raw messages.
func UnmarshalResourceRecordInputRdataRdataMxRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordInputRdataRdataMxRecord)
	err = core.UnmarshalPrimitive(m, "exchange", &obj.Exchange)
	if err != nil {
		err = core.SDKErrorf(err, "", "exchange-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "preference", &obj.Preference)
	if err != nil {
		err = core.SDKErrorf(err, "", "preference-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordInputRdataRdataPtrRecord : The content of type-PTR resource record.
// This model "extends" ResourceRecordInputRdata
type ResourceRecordInputRdataRdataPtrRecord struct {
	// Hostname of the relevant A or AAAA record.
	Ptrdname *string `json:"ptrdname" validate:"required"`
}

// NewResourceRecordInputRdataRdataPtrRecord : Instantiate ResourceRecordInputRdataRdataPtrRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataPtrRecord(ptrdname string) (_model *ResourceRecordInputRdataRdataPtrRecord, err error) {
	_model = &ResourceRecordInputRdataRdataPtrRecord{
		Ptrdname: core.StringPtr(ptrdname),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordInputRdataRdataPtrRecord) isaResourceRecordInputRdata() bool {
	return true
}

// UnmarshalResourceRecordInputRdataRdataPtrRecord unmarshals an instance of ResourceRecordInputRdataRdataPtrRecord from the specified map of raw messages.
func UnmarshalResourceRecordInputRdataRdataPtrRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordInputRdataRdataPtrRecord)
	err = core.UnmarshalPrimitive(m, "ptrdname", &obj.Ptrdname)
	if err != nil {
		err = core.SDKErrorf(err, "", "ptrdname-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordInputRdataRdataSrvRecord : The content of type-SRV resource record.
// This model "extends" ResourceRecordInputRdata
type ResourceRecordInputRdataRdataSrvRecord struct {
	// Port number of the target server.
	Port *int64 `json:"port" validate:"required"`

	// Priority of the SRV record.
	Priority *int64 `json:"priority" validate:"required"`

	// Hostname of the target server.
	Target *string `json:"target" validate:"required"`

	// Weight of distributing queries among multiple target servers.
	Weight *int64 `json:"weight" validate:"required"`
}

// NewResourceRecordInputRdataRdataSrvRecord : Instantiate ResourceRecordInputRdataRdataSrvRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataSrvRecord(port int64, priority int64, target string, weight int64) (_model *ResourceRecordInputRdataRdataSrvRecord, err error) {
	_model = &ResourceRecordInputRdataRdataSrvRecord{
		Port:     core.Int64Ptr(port),
		Priority: core.Int64Ptr(priority),
		Target:   core.StringPtr(target),
		Weight:   core.Int64Ptr(weight),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordInputRdataRdataSrvRecord) isaResourceRecordInputRdata() bool {
	return true
}

// UnmarshalResourceRecordInputRdataRdataSrvRecord unmarshals an instance of ResourceRecordInputRdataRdataSrvRecord from the specified map of raw messages.
func UnmarshalResourceRecordInputRdataRdataSrvRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordInputRdataRdataSrvRecord)
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "weight", &obj.Weight)
	if err != nil {
		err = core.SDKErrorf(err, "", "weight-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordInputRdataRdataTxtRecord : The content of type-TXT resource record.
// This model "extends" ResourceRecordInputRdata
type ResourceRecordInputRdataRdataTxtRecord struct {
	// Human readable text.
	Text *string `json:"text" validate:"required"`
}

// NewResourceRecordInputRdataRdataTxtRecord : Instantiate ResourceRecordInputRdataRdataTxtRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataTxtRecord(text string) (_model *ResourceRecordInputRdataRdataTxtRecord, err error) {
	_model = &ResourceRecordInputRdataRdataTxtRecord{
		Text: core.StringPtr(text),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordInputRdataRdataTxtRecord) isaResourceRecordInputRdata() bool {
	return true
}

// UnmarshalResourceRecordInputRdataRdataTxtRecord unmarshals an instance of ResourceRecordInputRdataRdataTxtRecord from the specified map of raw messages.
func UnmarshalResourceRecordInputRdataRdataTxtRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordInputRdataRdataTxtRecord)
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		err = core.SDKErrorf(err, "", "text-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordUpdateInputRdataRdataARecord : The content of type-A resource record.
// This model "extends" ResourceRecordUpdateInputRdata
type ResourceRecordUpdateInputRdataRdataARecord struct {
	// IPv4 address.
	Ip *string `json:"ip" validate:"required"`
}

// NewResourceRecordUpdateInputRdataRdataARecord : Instantiate ResourceRecordUpdateInputRdataRdataARecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataARecord(ip string) (_model *ResourceRecordUpdateInputRdataRdataARecord, err error) {
	_model = &ResourceRecordUpdateInputRdataRdataARecord{
		Ip: core.StringPtr(ip),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordUpdateInputRdataRdataARecord) isaResourceRecordUpdateInputRdata() bool {
	return true
}

// UnmarshalResourceRecordUpdateInputRdataRdataARecord unmarshals an instance of ResourceRecordUpdateInputRdataRdataARecord from the specified map of raw messages.
func UnmarshalResourceRecordUpdateInputRdataRdataARecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordUpdateInputRdataRdataARecord)
	err = core.UnmarshalPrimitive(m, "ip", &obj.Ip)
	if err != nil {
		err = core.SDKErrorf(err, "", "ip-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordUpdateInputRdataRdataAaaaRecord : The content of type-AAAA resource record.
// This model "extends" ResourceRecordUpdateInputRdata
type ResourceRecordUpdateInputRdataRdataAaaaRecord struct {
	// IPv6 address.
	Ip *string `json:"ip" validate:"required"`
}

// NewResourceRecordUpdateInputRdataRdataAaaaRecord : Instantiate ResourceRecordUpdateInputRdataRdataAaaaRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataAaaaRecord(ip string) (_model *ResourceRecordUpdateInputRdataRdataAaaaRecord, err error) {
	_model = &ResourceRecordUpdateInputRdataRdataAaaaRecord{
		Ip: core.StringPtr(ip),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordUpdateInputRdataRdataAaaaRecord) isaResourceRecordUpdateInputRdata() bool {
	return true
}

// UnmarshalResourceRecordUpdateInputRdataRdataAaaaRecord unmarshals an instance of ResourceRecordUpdateInputRdataRdataAaaaRecord from the specified map of raw messages.
func UnmarshalResourceRecordUpdateInputRdataRdataAaaaRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordUpdateInputRdataRdataAaaaRecord)
	err = core.UnmarshalPrimitive(m, "ip", &obj.Ip)
	if err != nil {
		err = core.SDKErrorf(err, "", "ip-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordUpdateInputRdataRdataCnameRecord : The content of type-CNAME resource record.
// This model "extends" ResourceRecordUpdateInputRdata
type ResourceRecordUpdateInputRdataRdataCnameRecord struct {
	// Canonical name.
	Cname *string `json:"cname" validate:"required"`
}

// NewResourceRecordUpdateInputRdataRdataCnameRecord : Instantiate ResourceRecordUpdateInputRdataRdataCnameRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataCnameRecord(cname string) (_model *ResourceRecordUpdateInputRdataRdataCnameRecord, err error) {
	_model = &ResourceRecordUpdateInputRdataRdataCnameRecord{
		Cname: core.StringPtr(cname),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordUpdateInputRdataRdataCnameRecord) isaResourceRecordUpdateInputRdata() bool {
	return true
}

// UnmarshalResourceRecordUpdateInputRdataRdataCnameRecord unmarshals an instance of ResourceRecordUpdateInputRdataRdataCnameRecord from the specified map of raw messages.
func UnmarshalResourceRecordUpdateInputRdataRdataCnameRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordUpdateInputRdataRdataCnameRecord)
	err = core.UnmarshalPrimitive(m, "cname", &obj.Cname)
	if err != nil {
		err = core.SDKErrorf(err, "", "cname-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordUpdateInputRdataRdataMxRecord : The content of type-MX resource record.
// This model "extends" ResourceRecordUpdateInputRdata
type ResourceRecordUpdateInputRdataRdataMxRecord struct {
	// Hostname of Exchange server.
	Exchange *string `json:"exchange" validate:"required"`

	// Preference of the MX record.
	Preference *int64 `json:"preference" validate:"required"`
}

// NewResourceRecordUpdateInputRdataRdataMxRecord : Instantiate ResourceRecordUpdateInputRdataRdataMxRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataMxRecord(exchange string, preference int64) (_model *ResourceRecordUpdateInputRdataRdataMxRecord, err error) {
	_model = &ResourceRecordUpdateInputRdataRdataMxRecord{
		Exchange:   core.StringPtr(exchange),
		Preference: core.Int64Ptr(preference),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordUpdateInputRdataRdataMxRecord) isaResourceRecordUpdateInputRdata() bool {
	return true
}

// UnmarshalResourceRecordUpdateInputRdataRdataMxRecord unmarshals an instance of ResourceRecordUpdateInputRdataRdataMxRecord from the specified map of raw messages.
func UnmarshalResourceRecordUpdateInputRdataRdataMxRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordUpdateInputRdataRdataMxRecord)
	err = core.UnmarshalPrimitive(m, "exchange", &obj.Exchange)
	if err != nil {
		err = core.SDKErrorf(err, "", "exchange-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "preference", &obj.Preference)
	if err != nil {
		err = core.SDKErrorf(err, "", "preference-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordUpdateInputRdataRdataPtrRecord : The content of type-PTR resource record.
// This model "extends" ResourceRecordUpdateInputRdata
type ResourceRecordUpdateInputRdataRdataPtrRecord struct {
	// Hostname of the relevant A or AAAA record.
	Ptrdname *string `json:"ptrdname" validate:"required"`
}

// NewResourceRecordUpdateInputRdataRdataPtrRecord : Instantiate ResourceRecordUpdateInputRdataRdataPtrRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataPtrRecord(ptrdname string) (_model *ResourceRecordUpdateInputRdataRdataPtrRecord, err error) {
	_model = &ResourceRecordUpdateInputRdataRdataPtrRecord{
		Ptrdname: core.StringPtr(ptrdname),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordUpdateInputRdataRdataPtrRecord) isaResourceRecordUpdateInputRdata() bool {
	return true
}

// UnmarshalResourceRecordUpdateInputRdataRdataPtrRecord unmarshals an instance of ResourceRecordUpdateInputRdataRdataPtrRecord from the specified map of raw messages.
func UnmarshalResourceRecordUpdateInputRdataRdataPtrRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordUpdateInputRdataRdataPtrRecord)
	err = core.UnmarshalPrimitive(m, "ptrdname", &obj.Ptrdname)
	if err != nil {
		err = core.SDKErrorf(err, "", "ptrdname-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordUpdateInputRdataRdataSrvRecord : The content of type-SRV resource record.
// This model "extends" ResourceRecordUpdateInputRdata
type ResourceRecordUpdateInputRdataRdataSrvRecord struct {
	// Port number of the target server.
	Port *int64 `json:"port" validate:"required"`

	// Priority of the SRV record.
	Priority *int64 `json:"priority" validate:"required"`

	// Hostname of the target server.
	Target *string `json:"target" validate:"required"`

	// Weight of distributing queries among multiple target servers.
	Weight *int64 `json:"weight" validate:"required"`
}

// NewResourceRecordUpdateInputRdataRdataSrvRecord : Instantiate ResourceRecordUpdateInputRdataRdataSrvRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataSrvRecord(port int64, priority int64, target string, weight int64) (_model *ResourceRecordUpdateInputRdataRdataSrvRecord, err error) {
	_model = &ResourceRecordUpdateInputRdataRdataSrvRecord{
		Port:     core.Int64Ptr(port),
		Priority: core.Int64Ptr(priority),
		Target:   core.StringPtr(target),
		Weight:   core.Int64Ptr(weight),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordUpdateInputRdataRdataSrvRecord) isaResourceRecordUpdateInputRdata() bool {
	return true
}

// UnmarshalResourceRecordUpdateInputRdataRdataSrvRecord unmarshals an instance of ResourceRecordUpdateInputRdataRdataSrvRecord from the specified map of raw messages.
func UnmarshalResourceRecordUpdateInputRdataRdataSrvRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordUpdateInputRdataRdataSrvRecord)
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "weight", &obj.Weight)
	if err != nil {
		err = core.SDKErrorf(err, "", "weight-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceRecordUpdateInputRdataRdataTxtRecord : The content of type-TXT resource record.
// This model "extends" ResourceRecordUpdateInputRdata
type ResourceRecordUpdateInputRdataRdataTxtRecord struct {
	// Human readable text.
	Text *string `json:"text" validate:"required"`
}

// NewResourceRecordUpdateInputRdataRdataTxtRecord : Instantiate ResourceRecordUpdateInputRdataRdataTxtRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataTxtRecord(text string) (_model *ResourceRecordUpdateInputRdataRdataTxtRecord, err error) {
	_model = &ResourceRecordUpdateInputRdataRdataTxtRecord{
		Text: core.StringPtr(text),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ResourceRecordUpdateInputRdataRdataTxtRecord) isaResourceRecordUpdateInputRdata() bool {
	return true
}

// UnmarshalResourceRecordUpdateInputRdataRdataTxtRecord unmarshals an instance of ResourceRecordUpdateInputRdataRdataTxtRecord from the specified map of raw messages.
func UnmarshalResourceRecordUpdateInputRdataRdataTxtRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecordUpdateInputRdataRdataTxtRecord)
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		err = core.SDKErrorf(err, "", "text-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ForwardingRuleInputForwardingRuleBoth : Forwarding rule with both forward_to and views.
// This model "extends" ForwardingRuleInput
type ForwardingRuleInputForwardingRuleBoth struct {
	// Descriptive text of the forwarding rule.
	Description *string `json:"description,omitempty"`

	// Type of the forwarding rule.
	Type *string `json:"type" validate:"required"`

	// The matching zone or hostname.
	Match *string `json:"match" validate:"required"`

	// The upstream DNS servers that the DNS queries will be forwarded to.
	ForwardTo []string `json:"forward_to" validate:"required"`

	// An array of views.
	Views []ViewConfig `json:"views" validate:"required"`
}

// Constants associated with the ForwardingRuleInputForwardingRuleBoth.Type property.
// Type of the forwarding rule.
const (
	ForwardingRuleInputForwardingRuleBoth_Type_Zone = "zone"
)

// NewForwardingRuleInputForwardingRuleBoth : Instantiate ForwardingRuleInputForwardingRuleBoth (Generic Model Constructor)
func (*DnsSvcsV1) NewForwardingRuleInputForwardingRuleBoth(typeVar string, match string, forwardTo []string, views []ViewConfig) (_model *ForwardingRuleInputForwardingRuleBoth, err error) {
	_model = &ForwardingRuleInputForwardingRuleBoth{
		Type:      core.StringPtr(typeVar),
		Match:     core.StringPtr(match),
		ForwardTo: forwardTo,
		Views:     views,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ForwardingRuleInputForwardingRuleBoth) isaForwardingRuleInput() bool {
	return true
}

// UnmarshalForwardingRuleInputForwardingRuleBoth unmarshals an instance of ForwardingRuleInputForwardingRuleBoth from the specified map of raw messages.
func UnmarshalForwardingRuleInputForwardingRuleBoth(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ForwardingRuleInputForwardingRuleBoth)
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
	err = core.UnmarshalPrimitive(m, "match", &obj.Match)
	if err != nil {
		err = core.SDKErrorf(err, "", "match-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "forward_to", &obj.ForwardTo)
	if err != nil {
		err = core.SDKErrorf(err, "", "forward_to-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "views", &obj.Views, UnmarshalViewConfig)
	if err != nil {
		err = core.SDKErrorf(err, "", "views-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ForwardingRuleInputForwardingRuleOnlyForward : Forwarding rule without views.
// This model "extends" ForwardingRuleInput
type ForwardingRuleInputForwardingRuleOnlyForward struct {
	// Descriptive text of the forwarding rule.
	Description *string `json:"description,omitempty"`

	// Type of the forwarding rule.
	Type *string `json:"type" validate:"required"`

	// The matching zone or hostname.
	Match *string `json:"match" validate:"required"`

	// The upstream DNS servers that the DNS queries will be forwarded to.
	ForwardTo []string `json:"forward_to" validate:"required"`
}

// Constants associated with the ForwardingRuleInputForwardingRuleOnlyForward.Type property.
// Type of the forwarding rule.
const (
	ForwardingRuleInputForwardingRuleOnlyForward_Type_Zone = "zone"
)

// NewForwardingRuleInputForwardingRuleOnlyForward : Instantiate ForwardingRuleInputForwardingRuleOnlyForward (Generic Model Constructor)
func (*DnsSvcsV1) NewForwardingRuleInputForwardingRuleOnlyForward(typeVar string, match string, forwardTo []string) (_model *ForwardingRuleInputForwardingRuleOnlyForward, err error) {
	_model = &ForwardingRuleInputForwardingRuleOnlyForward{
		Type:      core.StringPtr(typeVar),
		Match:     core.StringPtr(match),
		ForwardTo: forwardTo,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ForwardingRuleInputForwardingRuleOnlyForward) isaForwardingRuleInput() bool {
	return true
}

// UnmarshalForwardingRuleInputForwardingRuleOnlyForward unmarshals an instance of ForwardingRuleInputForwardingRuleOnlyForward from the specified map of raw messages.
func UnmarshalForwardingRuleInputForwardingRuleOnlyForward(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ForwardingRuleInputForwardingRuleOnlyForward)
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
	err = core.UnmarshalPrimitive(m, "match", &obj.Match)
	if err != nil {
		err = core.SDKErrorf(err, "", "match-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "forward_to", &obj.ForwardTo)
	if err != nil {
		err = core.SDKErrorf(err, "", "forward_to-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ForwardingRuleInputForwardingRuleOnlyView : Forwarding rule with only views.
// This model "extends" ForwardingRuleInput
type ForwardingRuleInputForwardingRuleOnlyView struct {
	// Descriptive text of the forwarding rule.
	Description *string `json:"description,omitempty"`

	// Type of the forwarding rule.
	Type *string `json:"type" validate:"required"`

	// The matching zone or hostname.
	Match *string `json:"match" validate:"required"`

	// An array of views.
	Views []ViewConfig `json:"views" validate:"required"`
}

// Constants associated with the ForwardingRuleInputForwardingRuleOnlyView.Type property.
// Type of the forwarding rule.
const (
	ForwardingRuleInputForwardingRuleOnlyView_Type_Zone = "zone"
)

// NewForwardingRuleInputForwardingRuleOnlyView : Instantiate ForwardingRuleInputForwardingRuleOnlyView (Generic Model Constructor)
func (*DnsSvcsV1) NewForwardingRuleInputForwardingRuleOnlyView(typeVar string, match string, views []ViewConfig) (_model *ForwardingRuleInputForwardingRuleOnlyView, err error) {
	_model = &ForwardingRuleInputForwardingRuleOnlyView{
		Type:  core.StringPtr(typeVar),
		Match: core.StringPtr(match),
		Views: views,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ForwardingRuleInputForwardingRuleOnlyView) isaForwardingRuleInput() bool {
	return true
}

// UnmarshalForwardingRuleInputForwardingRuleOnlyView unmarshals an instance of ForwardingRuleInputForwardingRuleOnlyView from the specified map of raw messages.
func UnmarshalForwardingRuleInputForwardingRuleOnlyView(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ForwardingRuleInputForwardingRuleOnlyView)
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
	err = core.UnmarshalPrimitive(m, "match", &obj.Match)
	if err != nil {
		err = core.SDKErrorf(err, "", "match-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "views", &obj.Views, UnmarshalViewConfig)
	if err != nil {
		err = core.SDKErrorf(err, "", "views-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DnszonesPager can be used to simplify the use of the "ListDnszones" method.
type DnszonesPager struct {
	hasNext     bool
	options     *ListDnszonesOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewDnszonesPager returns a new DnszonesPager instance.
func (dnsSvcs *DnsSvcsV1) NewDnszonesPager(options *ListDnszonesOptions) (pager *DnszonesPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListDnszonesOptions = *options
	pager = &DnszonesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *DnszonesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *DnszonesPager) GetNextWithContext(ctx context.Context) (page []Dnszone, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListDnszonesWithContext(ctx, pager.options)
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
	page = result.Dnszones

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *DnszonesPager) GetAllWithContext(ctx context.Context) (allItems []Dnszone, err error) {
	for pager.HasNext() {
		var nextPage []Dnszone
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
func (pager *DnszonesPager) GetNext() (page []Dnszone, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *DnszonesPager) GetAll() (allItems []Dnszone, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ResourceRecordsPager can be used to simplify the use of the "ListResourceRecords" method.
type ResourceRecordsPager struct {
	hasNext     bool
	options     *ListResourceRecordsOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewResourceRecordsPager returns a new ResourceRecordsPager instance.
func (dnsSvcs *DnsSvcsV1) NewResourceRecordsPager(options *ListResourceRecordsOptions) (pager *ResourceRecordsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListResourceRecordsOptions = *options
	pager = &ResourceRecordsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ResourceRecordsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ResourceRecordsPager) GetNextWithContext(ctx context.Context) (page []ResourceRecord, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListResourceRecordsWithContext(ctx, pager.options)
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
	page = result.ResourceRecords

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ResourceRecordsPager) GetAllWithContext(ctx context.Context) (allItems []ResourceRecord, err error) {
	for pager.HasNext() {
		var nextPage []ResourceRecord
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
func (pager *ResourceRecordsPager) GetNext() (page []ResourceRecord, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ResourceRecordsPager) GetAll() (allItems []ResourceRecord, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// LoadBalancersPager can be used to simplify the use of the "ListLoadBalancers" method.
type LoadBalancersPager struct {
	hasNext     bool
	options     *ListLoadBalancersOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewLoadBalancersPager returns a new LoadBalancersPager instance.
func (dnsSvcs *DnsSvcsV1) NewLoadBalancersPager(options *ListLoadBalancersOptions) (pager *LoadBalancersPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListLoadBalancersOptions = *options
	pager = &LoadBalancersPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *LoadBalancersPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *LoadBalancersPager) GetNextWithContext(ctx context.Context) (page []LoadBalancer, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListLoadBalancersWithContext(ctx, pager.options)
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
	page = result.LoadBalancers

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *LoadBalancersPager) GetAllWithContext(ctx context.Context) (allItems []LoadBalancer, err error) {
	for pager.HasNext() {
		var nextPage []LoadBalancer
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
func (pager *LoadBalancersPager) GetNext() (page []LoadBalancer, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *LoadBalancersPager) GetAll() (allItems []LoadBalancer, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PoolsPager can be used to simplify the use of the "ListPools" method.
type PoolsPager struct {
	hasNext     bool
	options     *ListPoolsOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewPoolsPager returns a new PoolsPager instance.
func (dnsSvcs *DnsSvcsV1) NewPoolsPager(options *ListPoolsOptions) (pager *PoolsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListPoolsOptions = *options
	pager = &PoolsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *PoolsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *PoolsPager) GetNextWithContext(ctx context.Context) (page []Pool, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListPoolsWithContext(ctx, pager.options)
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
	page = result.Pools

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *PoolsPager) GetAllWithContext(ctx context.Context) (allItems []Pool, err error) {
	for pager.HasNext() {
		var nextPage []Pool
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
func (pager *PoolsPager) GetNext() (page []Pool, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *PoolsPager) GetAll() (allItems []Pool, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// MonitorsPager can be used to simplify the use of the "ListMonitors" method.
type MonitorsPager struct {
	hasNext     bool
	options     *ListMonitorsOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewMonitorsPager returns a new MonitorsPager instance.
func (dnsSvcs *DnsSvcsV1) NewMonitorsPager(options *ListMonitorsOptions) (pager *MonitorsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListMonitorsOptions = *options
	pager = &MonitorsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *MonitorsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *MonitorsPager) GetNextWithContext(ctx context.Context) (page []Monitor, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListMonitorsWithContext(ctx, pager.options)
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
	page = result.Monitors

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *MonitorsPager) GetAllWithContext(ctx context.Context) (allItems []Monitor, err error) {
	for pager.HasNext() {
		var nextPage []Monitor
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
func (pager *MonitorsPager) GetNext() (page []Monitor, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *MonitorsPager) GetAll() (allItems []Monitor, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ForwardingRulesPager can be used to simplify the use of the "ListForwardingRules" method.
type ForwardingRulesPager struct {
	hasNext     bool
	options     *ListForwardingRulesOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewForwardingRulesPager returns a new ForwardingRulesPager instance.
func (dnsSvcs *DnsSvcsV1) NewForwardingRulesPager(options *ListForwardingRulesOptions) (pager *ForwardingRulesPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListForwardingRulesOptions = *options
	pager = &ForwardingRulesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ForwardingRulesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ForwardingRulesPager) GetNextWithContext(ctx context.Context) (page []ForwardingRule, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListForwardingRulesWithContext(ctx, pager.options)
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
	page = result.ForwardingRules

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ForwardingRulesPager) GetAllWithContext(ctx context.Context) (allItems []ForwardingRule, err error) {
	for pager.HasNext() {
		var nextPage []ForwardingRule
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
func (pager *ForwardingRulesPager) GetNext() (page []ForwardingRule, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ForwardingRulesPager) GetAll() (allItems []ForwardingRule, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SecondaryZonesPager can be used to simplify the use of the "ListSecondaryZones" method.
type SecondaryZonesPager struct {
	hasNext     bool
	options     *ListSecondaryZonesOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewSecondaryZonesPager returns a new SecondaryZonesPager instance.
func (dnsSvcs *DnsSvcsV1) NewSecondaryZonesPager(options *ListSecondaryZonesOptions) (pager *SecondaryZonesPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListSecondaryZonesOptions = *options
	pager = &SecondaryZonesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SecondaryZonesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SecondaryZonesPager) GetNextWithContext(ctx context.Context) (page []SecondaryZone, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListSecondaryZonesWithContext(ctx, pager.options)
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
	page = result.SecondaryZones

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SecondaryZonesPager) GetAllWithContext(ctx context.Context) (allItems []SecondaryZone, err error) {
	for pager.HasNext() {
		var nextPage []SecondaryZone
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
func (pager *SecondaryZonesPager) GetNext() (page []SecondaryZone, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SecondaryZonesPager) GetAll() (allItems []SecondaryZone, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// LinkedZonesPager can be used to simplify the use of the "ListLinkedZones" method.
type LinkedZonesPager struct {
	hasNext     bool
	options     *ListLinkedZonesOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewLinkedZonesPager returns a new LinkedZonesPager instance.
func (dnsSvcs *DnsSvcsV1) NewLinkedZonesPager(options *ListLinkedZonesOptions) (pager *LinkedZonesPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListLinkedZonesOptions = *options
	pager = &LinkedZonesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *LinkedZonesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *LinkedZonesPager) GetNextWithContext(ctx context.Context) (page []LinkedDnszone, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListLinkedZonesWithContext(ctx, pager.options)
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
	page = result.LinkedDnszones

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *LinkedZonesPager) GetAllWithContext(ctx context.Context) (allItems []LinkedDnszone, err error) {
	for pager.HasNext() {
		var nextPage []LinkedDnszone
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
func (pager *LinkedZonesPager) GetNext() (page []LinkedDnszone, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *LinkedZonesPager) GetAll() (allItems []LinkedDnszone, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DnszoneAccessRequestsPager can be used to simplify the use of the "ListDnszoneAccessRequests" method.
type DnszoneAccessRequestsPager struct {
	hasNext     bool
	options     *ListDnszoneAccessRequestsOptions
	client      *DnsSvcsV1
	pageContext struct {
		next *int64
	}
}

// NewDnszoneAccessRequestsPager returns a new DnszoneAccessRequestsPager instance.
func (dnsSvcs *DnsSvcsV1) NewDnszoneAccessRequestsPager(options *ListDnszoneAccessRequestsOptions) (pager *DnszoneAccessRequestsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListDnszoneAccessRequestsOptions = *options
	pager = &DnszoneAccessRequestsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  dnsSvcs,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *DnszoneAccessRequestsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *DnszoneAccessRequestsPager) GetNextWithContext(ctx context.Context) (page []AccessRequest, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListDnszoneAccessRequestsWithContext(ctx, pager.options)
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
	page = result.AccessRequests

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *DnszoneAccessRequestsPager) GetAllWithContext(ctx context.Context) (allItems []AccessRequest, err error) {
	for pager.HasNext() {
		var nextPage []AccessRequest
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
func (pager *DnszoneAccessRequestsPager) GetNext() (page []AccessRequest, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *DnszoneAccessRequestsPager) GetAll() (allItems []AccessRequest, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
