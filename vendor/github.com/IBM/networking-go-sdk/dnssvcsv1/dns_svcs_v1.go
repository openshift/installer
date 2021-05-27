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
 

// Package dnssvcsv1 : Operations and models for the DnsSvcsV1 service
package dnssvcsv1

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

// DnsSvcsV1 : DNS Services API
//
// Version: 1.0.0
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
			return
		}
	}

	dnsSvcs, err = NewDnsSvcsV1(options)
	if err != nil {
		return
	}

	err = dnsSvcs.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = dnsSvcs.Service.SetServiceURL(options.URL)
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
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
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
	return "", fmt.Errorf("service does not support regional URLs")
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
	return dnsSvcs.Service.SetServiceURL(url)
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
	return dnsSvcs.ListDnszonesWithContext(context.Background(), listDnszonesOptions)
}

// ListDnszonesWithContext is an alternate form of the ListDnszones method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListDnszonesWithContext(ctx context.Context, listDnszonesOptions *ListDnszonesOptions) (result *ListDnszones, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listDnszonesOptions, "listDnszonesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listDnszonesOptions, "listDnszonesOptions")
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListDnszones)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateDnszone : Create a DNS zone
// Create a DNS zone for a given service instance.
func (dnsSvcs *DnsSvcsV1) CreateDnszone(createDnszoneOptions *CreateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	return dnsSvcs.CreateDnszoneWithContext(context.Background(), createDnszoneOptions)
}

// CreateDnszoneWithContext is an alternate form of the CreateDnszone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateDnszoneWithContext(ctx context.Context, createDnszoneOptions *CreateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDnszoneOptions, "createDnszoneOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createDnszoneOptions, "createDnszoneOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnszone)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteDnszone : Delete a DNS zone
// Delete a DNS zone.
func (dnsSvcs *DnsSvcsV1) DeleteDnszone(deleteDnszoneOptions *DeleteDnszoneOptions) (response *core.DetailedResponse, err error) {
	return dnsSvcs.DeleteDnszoneWithContext(context.Background(), deleteDnszoneOptions)
}

// DeleteDnszoneWithContext is an alternate form of the DeleteDnszone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteDnszoneWithContext(ctx context.Context, deleteDnszoneOptions *DeleteDnszoneOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDnszoneOptions, "deleteDnszoneOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteDnszoneOptions, "deleteDnszoneOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteDnszoneOptions.InstanceID,
		"dnszone_id": *deleteDnszoneOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)

	return
}

// GetDnszone : Get a DNS zone
// Get details of a DNS zone.
func (dnsSvcs *DnsSvcsV1) GetDnszone(getDnszoneOptions *GetDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	return dnsSvcs.GetDnszoneWithContext(context.Background(), getDnszoneOptions)
}

// GetDnszoneWithContext is an alternate form of the GetDnszone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetDnszoneWithContext(ctx context.Context, getDnszoneOptions *GetDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDnszoneOptions, "getDnszoneOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDnszoneOptions, "getDnszoneOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getDnszoneOptions.InstanceID,
		"dnszone_id": *getDnszoneOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnszone)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateDnszone : Update the properties of a DNS zone
// Update the properties of a DNS zone.
func (dnsSvcs *DnsSvcsV1) UpdateDnszone(updateDnszoneOptions *UpdateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	return dnsSvcs.UpdateDnszoneWithContext(context.Background(), updateDnszoneOptions)
}

// UpdateDnszoneWithContext is an alternate form of the UpdateDnszone method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateDnszoneWithContext(ctx context.Context, updateDnszoneOptions *UpdateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDnszoneOptions, "updateDnszoneOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateDnszoneOptions, "updateDnszoneOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateDnszoneOptions.InstanceID,
		"dnszone_id": *updateDnszoneOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnszone)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListResourceRecords : List Resource Records
// List the Resource Records for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) ListResourceRecords(listResourceRecordsOptions *ListResourceRecordsOptions) (result *ListResourceRecords, response *core.DetailedResponse, err error) {
	return dnsSvcs.ListResourceRecordsWithContext(context.Background(), listResourceRecordsOptions)
}

// ListResourceRecordsWithContext is an alternate form of the ListResourceRecords method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListResourceRecordsWithContext(ctx context.Context, listResourceRecordsOptions *ListResourceRecordsOptions) (result *ListResourceRecords, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listResourceRecordsOptions, "listResourceRecordsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listResourceRecordsOptions, "listResourceRecordsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listResourceRecordsOptions.InstanceID,
		"dnszone_id": *listResourceRecordsOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records`, pathParamsMap)
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListResourceRecords)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateResourceRecord : Create a resource record
// Create a resource record for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) CreateResourceRecord(createResourceRecordOptions *CreateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	return dnsSvcs.CreateResourceRecordWithContext(context.Background(), createResourceRecordOptions)
}

// CreateResourceRecordWithContext is an alternate form of the CreateResourceRecord method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateResourceRecordWithContext(ctx context.Context, createResourceRecordOptions *CreateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createResourceRecordOptions, "createResourceRecordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createResourceRecordOptions, "createResourceRecordOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createResourceRecordOptions.InstanceID,
		"dnszone_id": *createResourceRecordOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records`, pathParamsMap)
	if err != nil {
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
	if createResourceRecordOptions.Name != nil {
		body["name"] = createResourceRecordOptions.Name
	}
	if createResourceRecordOptions.Type != nil {
		body["type"] = createResourceRecordOptions.Type
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceRecord)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteResourceRecord : Delete a resource record
// Delete a resource record.
func (dnsSvcs *DnsSvcsV1) DeleteResourceRecord(deleteResourceRecordOptions *DeleteResourceRecordOptions) (response *core.DetailedResponse, err error) {
	return dnsSvcs.DeleteResourceRecordWithContext(context.Background(), deleteResourceRecordOptions)
}

// DeleteResourceRecordWithContext is an alternate form of the DeleteResourceRecord method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteResourceRecordWithContext(ctx context.Context, deleteResourceRecordOptions *DeleteResourceRecordOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteResourceRecordOptions, "deleteResourceRecordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteResourceRecordOptions, "deleteResourceRecordOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteResourceRecordOptions.InstanceID,
		"dnszone_id": *deleteResourceRecordOptions.DnszoneID,
		"record_id": *deleteResourceRecordOptions.RecordID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records/{record_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)

	return
}

// GetResourceRecord : Get a resource record
// Get details of a resource record.
func (dnsSvcs *DnsSvcsV1) GetResourceRecord(getResourceRecordOptions *GetResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	return dnsSvcs.GetResourceRecordWithContext(context.Background(), getResourceRecordOptions)
}

// GetResourceRecordWithContext is an alternate form of the GetResourceRecord method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetResourceRecordWithContext(ctx context.Context, getResourceRecordOptions *GetResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceRecordOptions, "getResourceRecordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getResourceRecordOptions, "getResourceRecordOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getResourceRecordOptions.InstanceID,
		"dnszone_id": *getResourceRecordOptions.DnszoneID,
		"record_id": *getResourceRecordOptions.RecordID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records/{record_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceRecord)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateResourceRecord : Update the properties of a resource record
// Update the properties of a resource record.
func (dnsSvcs *DnsSvcsV1) UpdateResourceRecord(updateResourceRecordOptions *UpdateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	return dnsSvcs.UpdateResourceRecordWithContext(context.Background(), updateResourceRecordOptions)
}

// UpdateResourceRecordWithContext is an alternate form of the UpdateResourceRecord method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateResourceRecordWithContext(ctx context.Context, updateResourceRecordOptions *UpdateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateResourceRecordOptions, "updateResourceRecordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateResourceRecordOptions, "updateResourceRecordOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateResourceRecordOptions.InstanceID,
		"dnszone_id": *updateResourceRecordOptions.DnszoneID,
		"record_id": *updateResourceRecordOptions.RecordID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/resource_records/{record_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceRecord)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListPermittedNetworks : List permitted networks
// List the permitted networks for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) ListPermittedNetworks(listPermittedNetworksOptions *ListPermittedNetworksOptions) (result *ListPermittedNetworks, response *core.DetailedResponse, err error) {
	return dnsSvcs.ListPermittedNetworksWithContext(context.Background(), listPermittedNetworksOptions)
}

// ListPermittedNetworksWithContext is an alternate form of the ListPermittedNetworks method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListPermittedNetworksWithContext(ctx context.Context, listPermittedNetworksOptions *ListPermittedNetworksOptions) (result *ListPermittedNetworks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPermittedNetworksOptions, "listPermittedNetworksOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listPermittedNetworksOptions, "listPermittedNetworksOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listPermittedNetworksOptions.InstanceID,
		"dnszone_id": *listPermittedNetworksOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/permitted_networks`, pathParamsMap)
	if err != nil {
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

	if listPermittedNetworksOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listPermittedNetworksOptions.Offset))
	}
	if listPermittedNetworksOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listPermittedNetworksOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListPermittedNetworks)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreatePermittedNetwork : Create a permitted network
// Create a permitted network for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) CreatePermittedNetwork(createPermittedNetworkOptions *CreatePermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	return dnsSvcs.CreatePermittedNetworkWithContext(context.Background(), createPermittedNetworkOptions)
}

// CreatePermittedNetworkWithContext is an alternate form of the CreatePermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreatePermittedNetworkWithContext(ctx context.Context, createPermittedNetworkOptions *CreatePermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPermittedNetworkOptions, "createPermittedNetworkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createPermittedNetworkOptions, "createPermittedNetworkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createPermittedNetworkOptions.InstanceID,
		"dnszone_id": *createPermittedNetworkOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/permitted_networks`, pathParamsMap)
	if err != nil {
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

	body := make(map[string]interface{})
	if createPermittedNetworkOptions.Type != nil {
		body["type"] = createPermittedNetworkOptions.Type
	}
	if createPermittedNetworkOptions.PermittedNetwork != nil {
		body["permitted_network"] = createPermittedNetworkOptions.PermittedNetwork
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeletePermittedNetwork : Remove a permitted network
// Remove a permitted network.
func (dnsSvcs *DnsSvcsV1) DeletePermittedNetwork(deletePermittedNetworkOptions *DeletePermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	return dnsSvcs.DeletePermittedNetworkWithContext(context.Background(), deletePermittedNetworkOptions)
}

// DeletePermittedNetworkWithContext is an alternate form of the DeletePermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeletePermittedNetworkWithContext(ctx context.Context, deletePermittedNetworkOptions *DeletePermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePermittedNetworkOptions, "deletePermittedNetworkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deletePermittedNetworkOptions, "deletePermittedNetworkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deletePermittedNetworkOptions.InstanceID,
		"dnszone_id": *deletePermittedNetworkOptions.DnszoneID,
		"permitted_network_id": *deletePermittedNetworkOptions.PermittedNetworkID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/permitted_networks/{permitted_network_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetPermittedNetwork : Get a permitted network
// Get details of a permitted network.
func (dnsSvcs *DnsSvcsV1) GetPermittedNetwork(getPermittedNetworkOptions *GetPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	return dnsSvcs.GetPermittedNetworkWithContext(context.Background(), getPermittedNetworkOptions)
}

// GetPermittedNetworkWithContext is an alternate form of the GetPermittedNetwork method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetPermittedNetworkWithContext(ctx context.Context, getPermittedNetworkOptions *GetPermittedNetworkOptions) (result *PermittedNetwork, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPermittedNetworkOptions, "getPermittedNetworkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getPermittedNetworkOptions, "getPermittedNetworkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getPermittedNetworkOptions.InstanceID,
		"dnszone_id": *getPermittedNetworkOptions.DnszoneID,
		"permitted_network_id": *getPermittedNetworkOptions.PermittedNetworkID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/permitted_networks/{permitted_network_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPermittedNetwork)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListLoadBalancers : List load balancers
// List the Global Load Balancers for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) ListLoadBalancers(listLoadBalancersOptions *ListLoadBalancersOptions) (result *ListLoadBalancers, response *core.DetailedResponse, err error) {
	return dnsSvcs.ListLoadBalancersWithContext(context.Background(), listLoadBalancersOptions)
}

// ListLoadBalancersWithContext is an alternate form of the ListLoadBalancers method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListLoadBalancersWithContext(ctx context.Context, listLoadBalancersOptions *ListLoadBalancersOptions) (result *ListLoadBalancers, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listLoadBalancersOptions, "listLoadBalancersOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listLoadBalancersOptions, "listLoadBalancersOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listLoadBalancersOptions.InstanceID,
		"dnszone_id": *listLoadBalancersOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers`, pathParamsMap)
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListLoadBalancers)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateLoadBalancer : Create a load balancer
// Create a load balancer for a given DNS zone.
func (dnsSvcs *DnsSvcsV1) CreateLoadBalancer(createLoadBalancerOptions *CreateLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	return dnsSvcs.CreateLoadBalancerWithContext(context.Background(), createLoadBalancerOptions)
}

// CreateLoadBalancerWithContext is an alternate form of the CreateLoadBalancer method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateLoadBalancerWithContext(ctx context.Context, createLoadBalancerOptions *CreateLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createLoadBalancerOptions, "createLoadBalancerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createLoadBalancerOptions, "createLoadBalancerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createLoadBalancerOptions.InstanceID,
		"dnszone_id": *createLoadBalancerOptions.DnszoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers`, pathParamsMap)
	if err != nil {
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
	if createLoadBalancerOptions.Description != nil {
		body["description"] = createLoadBalancerOptions.Description
	}
	if createLoadBalancerOptions.Enabled != nil {
		body["enabled"] = createLoadBalancerOptions.Enabled
	}
	if createLoadBalancerOptions.TTL != nil {
		body["ttl"] = createLoadBalancerOptions.TTL
	}
	if createLoadBalancerOptions.FallbackPool != nil {
		body["fallback_pool"] = createLoadBalancerOptions.FallbackPool
	}
	if createLoadBalancerOptions.DefaultPools != nil {
		body["default_pools"] = createLoadBalancerOptions.DefaultPools
	}
	if createLoadBalancerOptions.AzPools != nil {
		body["az_pools"] = createLoadBalancerOptions.AzPools
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
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancer)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteLoadBalancer : Delete a load balancer
// Delete a load balancer.
func (dnsSvcs *DnsSvcsV1) DeleteLoadBalancer(deleteLoadBalancerOptions *DeleteLoadBalancerOptions) (response *core.DetailedResponse, err error) {
	return dnsSvcs.DeleteLoadBalancerWithContext(context.Background(), deleteLoadBalancerOptions)
}

// DeleteLoadBalancerWithContext is an alternate form of the DeleteLoadBalancer method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteLoadBalancerWithContext(ctx context.Context, deleteLoadBalancerOptions *DeleteLoadBalancerOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLoadBalancerOptions, "deleteLoadBalancerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLoadBalancerOptions, "deleteLoadBalancerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteLoadBalancerOptions.InstanceID,
		"dnszone_id": *deleteLoadBalancerOptions.DnszoneID,
		"lb_id": *deleteLoadBalancerOptions.LbID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers/{lb_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)

	return
}

// GetLoadBalancer : Get a load balancer
// Get details of a load balancer.
func (dnsSvcs *DnsSvcsV1) GetLoadBalancer(getLoadBalancerOptions *GetLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	return dnsSvcs.GetLoadBalancerWithContext(context.Background(), getLoadBalancerOptions)
}

// GetLoadBalancerWithContext is an alternate form of the GetLoadBalancer method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetLoadBalancerWithContext(ctx context.Context, getLoadBalancerOptions *GetLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLoadBalancerOptions, "getLoadBalancerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLoadBalancerOptions, "getLoadBalancerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getLoadBalancerOptions.InstanceID,
		"dnszone_id": *getLoadBalancerOptions.DnszoneID,
		"lb_id": *getLoadBalancerOptions.LbID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers/{lb_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancer)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateLoadBalancer : Update the properties of a load balancer
// Update the properties of a load balancer.
func (dnsSvcs *DnsSvcsV1) UpdateLoadBalancer(updateLoadBalancerOptions *UpdateLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	return dnsSvcs.UpdateLoadBalancerWithContext(context.Background(), updateLoadBalancerOptions)
}

// UpdateLoadBalancerWithContext is an alternate form of the UpdateLoadBalancer method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateLoadBalancerWithContext(ctx context.Context, updateLoadBalancerOptions *UpdateLoadBalancerOptions) (result *LoadBalancer, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLoadBalancerOptions, "updateLoadBalancerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateLoadBalancerOptions, "updateLoadBalancerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateLoadBalancerOptions.InstanceID,
		"dnszone_id": *updateLoadBalancerOptions.DnszoneID,
		"lb_id": *updateLoadBalancerOptions.LbID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/dnszones/{dnszone_id}/load_balancers/{lb_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancer)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListPools : List load balancer pools
// List the load balancer pools.
func (dnsSvcs *DnsSvcsV1) ListPools(listPoolsOptions *ListPoolsOptions) (result *ListPools, response *core.DetailedResponse, err error) {
	return dnsSvcs.ListPoolsWithContext(context.Background(), listPoolsOptions)
}

// ListPoolsWithContext is an alternate form of the ListPools method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListPoolsWithContext(ctx context.Context, listPoolsOptions *ListPoolsOptions) (result *ListPools, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPoolsOptions, "listPoolsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listPoolsOptions, "listPoolsOptions")
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListPools)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreatePool : Create a load balancer pool
// Create a load balancer pool.
func (dnsSvcs *DnsSvcsV1) CreatePool(createPoolOptions *CreatePoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	return dnsSvcs.CreatePoolWithContext(context.Background(), createPoolOptions)
}

// CreatePoolWithContext is an alternate form of the CreatePool method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreatePoolWithContext(ctx context.Context, createPoolOptions *CreatePoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPoolOptions, "createPoolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createPoolOptions, "createPoolOptions")
	if err != nil {
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
	if createPoolOptions.Description != nil {
		body["description"] = createPoolOptions.Description
	}
	if createPoolOptions.Enabled != nil {
		body["enabled"] = createPoolOptions.Enabled
	}
	if createPoolOptions.HealthyOriginsThreshold != nil {
		body["healthy_origins_threshold"] = createPoolOptions.HealthyOriginsThreshold
	}
	if createPoolOptions.Origins != nil {
		body["origins"] = createPoolOptions.Origins
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPool)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeletePool : Delete a load balancer pool
// Delete a load balancer pool.
func (dnsSvcs *DnsSvcsV1) DeletePool(deletePoolOptions *DeletePoolOptions) (response *core.DetailedResponse, err error) {
	return dnsSvcs.DeletePoolWithContext(context.Background(), deletePoolOptions)
}

// DeletePoolWithContext is an alternate form of the DeletePool method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeletePoolWithContext(ctx context.Context, deletePoolOptions *DeletePoolOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePoolOptions, "deletePoolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deletePoolOptions, "deletePoolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deletePoolOptions.InstanceID,
		"pool_id": *deletePoolOptions.PoolID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/pools/{pool_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)

	return
}

// GetPool : Get a load balancer pool
// Get details of a load balancer pool.
func (dnsSvcs *DnsSvcsV1) GetPool(getPoolOptions *GetPoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	return dnsSvcs.GetPoolWithContext(context.Background(), getPoolOptions)
}

// GetPoolWithContext is an alternate form of the GetPool method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetPoolWithContext(ctx context.Context, getPoolOptions *GetPoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPoolOptions, "getPoolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getPoolOptions, "getPoolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getPoolOptions.InstanceID,
		"pool_id": *getPoolOptions.PoolID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/pools/{pool_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPool)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdatePool : Update the properties of a load balancer pool
// Update the properties of a load balancer pool.
func (dnsSvcs *DnsSvcsV1) UpdatePool(updatePoolOptions *UpdatePoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	return dnsSvcs.UpdatePoolWithContext(context.Background(), updatePoolOptions)
}

// UpdatePoolWithContext is an alternate form of the UpdatePool method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdatePoolWithContext(ctx context.Context, updatePoolOptions *UpdatePoolOptions) (result *Pool, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePoolOptions, "updatePoolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updatePoolOptions, "updatePoolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updatePoolOptions.InstanceID,
		"pool_id": *updatePoolOptions.PoolID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/pools/{pool_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPool)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListMonitors : List load balancer monitors
// List the load balancer monitors.
func (dnsSvcs *DnsSvcsV1) ListMonitors(listMonitorsOptions *ListMonitorsOptions) (result *ListMonitors, response *core.DetailedResponse, err error) {
	return dnsSvcs.ListMonitorsWithContext(context.Background(), listMonitorsOptions)
}

// ListMonitorsWithContext is an alternate form of the ListMonitors method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) ListMonitorsWithContext(ctx context.Context, listMonitorsOptions *ListMonitorsOptions) (result *ListMonitors, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listMonitorsOptions, "listMonitorsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listMonitorsOptions, "listMonitorsOptions")
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListMonitors)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateMonitor : Create a load balancer monitor
// Create a load balancer monitor.
func (dnsSvcs *DnsSvcsV1) CreateMonitor(createMonitorOptions *CreateMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	return dnsSvcs.CreateMonitorWithContext(context.Background(), createMonitorOptions)
}

// CreateMonitorWithContext is an alternate form of the CreateMonitor method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) CreateMonitorWithContext(ctx context.Context, createMonitorOptions *CreateMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createMonitorOptions, "createMonitorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createMonitorOptions, "createMonitorOptions")
	if err != nil {
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
	if createMonitorOptions.Description != nil {
		body["description"] = createMonitorOptions.Description
	}
	if createMonitorOptions.Type != nil {
		body["type"] = createMonitorOptions.Type
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitor)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteMonitor : Delete a load balancer monitor
// Delete a load balancer monitor.
func (dnsSvcs *DnsSvcsV1) DeleteMonitor(deleteMonitorOptions *DeleteMonitorOptions) (response *core.DetailedResponse, err error) {
	return dnsSvcs.DeleteMonitorWithContext(context.Background(), deleteMonitorOptions)
}

// DeleteMonitorWithContext is an alternate form of the DeleteMonitor method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) DeleteMonitorWithContext(ctx context.Context, deleteMonitorOptions *DeleteMonitorOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteMonitorOptions, "deleteMonitorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteMonitorOptions, "deleteMonitorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteMonitorOptions.InstanceID,
		"monitor_id": *deleteMonitorOptions.MonitorID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/monitors/{monitor_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	response, err = dnsSvcs.Service.Request(request, nil)

	return
}

// GetMonitor : Get a load balancer monitor
// Get details of a load balancer monitor.
func (dnsSvcs *DnsSvcsV1) GetMonitor(getMonitorOptions *GetMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	return dnsSvcs.GetMonitorWithContext(context.Background(), getMonitorOptions)
}

// GetMonitorWithContext is an alternate form of the GetMonitor method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) GetMonitorWithContext(ctx context.Context, getMonitorOptions *GetMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getMonitorOptions, "getMonitorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getMonitorOptions, "getMonitorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getMonitorOptions.InstanceID,
		"monitor_id": *getMonitorOptions.MonitorID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/monitors/{monitor_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitor)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateMonitor : Update the properties of a load balancer monitor
// Update the properties of a load balancer monitor.
func (dnsSvcs *DnsSvcsV1) UpdateMonitor(updateMonitorOptions *UpdateMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	return dnsSvcs.UpdateMonitorWithContext(context.Background(), updateMonitorOptions)
}

// UpdateMonitorWithContext is an alternate form of the UpdateMonitor method which supports a Context parameter
func (dnsSvcs *DnsSvcsV1) UpdateMonitorWithContext(ctx context.Context, updateMonitorOptions *UpdateMonitorOptions) (result *Monitor, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateMonitorOptions, "updateMonitorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateMonitorOptions, "updateMonitorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateMonitorOptions.InstanceID,
		"monitor_id": *updateMonitorOptions.MonitorID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsSvcs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsSvcs.Service.Options.URL, `/instances/{instance_id}/monitors/{monitor_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsSvcs.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitor)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateDnszoneOptions : The CreateDnszone options.
type CreateDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Name of DNS zone.
	Name *string `json:"name,omitempty"`

	// The text describing the purpose of a DNS zone.
	Description *string `json:"description,omitempty"`

	// The label of a DNS zone.
	Label *string `json:"label,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateDnszoneOptions : Instantiate CreateDnszoneOptions
func (*DnsSvcsV1) NewCreateDnszoneOptions(instanceID string) *CreateDnszoneOptions {
	return &CreateDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *CreateDnszoneOptions) SetInstanceID(instanceID string) *CreateDnszoneOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetName : Allow user to set Name
func (options *CreateDnszoneOptions) SetName(name string) *CreateDnszoneOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateDnszoneOptions) SetDescription(description string) *CreateDnszoneOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetLabel : Allow user to set Label
func (options *CreateDnszoneOptions) SetLabel(label string) *CreateDnszoneOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *CreateDnszoneOptions) SetXCorrelationID(xCorrelationID string) *CreateDnszoneOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDnszoneOptions) SetHeaders(param map[string]string) *CreateDnszoneOptions {
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
	AzPools []LoadBalancerAzPoolsItem `json:"az_pools,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateLoadBalancerOptions : Instantiate CreateLoadBalancerOptions
func (*DnsSvcsV1) NewCreateLoadBalancerOptions(instanceID string, dnszoneID string) *CreateLoadBalancerOptions {
	return &CreateLoadBalancerOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *CreateLoadBalancerOptions) SetInstanceID(instanceID string) *CreateLoadBalancerOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *CreateLoadBalancerOptions) SetDnszoneID(dnszoneID string) *CreateLoadBalancerOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetName : Allow user to set Name
func (options *CreateLoadBalancerOptions) SetName(name string) *CreateLoadBalancerOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateLoadBalancerOptions) SetDescription(description string) *CreateLoadBalancerOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *CreateLoadBalancerOptions) SetEnabled(enabled bool) *CreateLoadBalancerOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetTTL : Allow user to set TTL
func (options *CreateLoadBalancerOptions) SetTTL(ttl int64) *CreateLoadBalancerOptions {
	options.TTL = core.Int64Ptr(ttl)
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

// SetAzPools : Allow user to set AzPools
func (options *CreateLoadBalancerOptions) SetAzPools(azPools []LoadBalancerAzPoolsItem) *CreateLoadBalancerOptions {
	options.AzPools = azPools
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *CreateLoadBalancerOptions) SetXCorrelationID(xCorrelationID string) *CreateLoadBalancerOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLoadBalancerOptions) SetHeaders(param map[string]string) *CreateLoadBalancerOptions {
	options.Headers = param
	return options
}

// CreateMonitorOptions : The CreateMonitor options.
type CreateMonitorOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

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

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateMonitorOptions.Type property.
// The protocol to use for the health check. Currently supported protocols are 'HTTP','HTTPS' and 'TCP'.
const (
	CreateMonitorOptions_Type_Http = "HTTP"
	CreateMonitorOptions_Type_Https = "HTTPS"
	CreateMonitorOptions_Type_Tcp = "TCP"
)

// Constants associated with the CreateMonitorOptions.Method property.
// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
const (
	CreateMonitorOptions_Method_Get = "GET"
	CreateMonitorOptions_Method_Head = "HEAD"
)

// NewCreateMonitorOptions : Instantiate CreateMonitorOptions
func (*DnsSvcsV1) NewCreateMonitorOptions(instanceID string) *CreateMonitorOptions {
	return &CreateMonitorOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *CreateMonitorOptions) SetInstanceID(instanceID string) *CreateMonitorOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetName : Allow user to set Name
func (options *CreateMonitorOptions) SetName(name string) *CreateMonitorOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateMonitorOptions) SetDescription(description string) *CreateMonitorOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetType : Allow user to set Type
func (options *CreateMonitorOptions) SetType(typeVar string) *CreateMonitorOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetPort : Allow user to set Port
func (options *CreateMonitorOptions) SetPort(port int64) *CreateMonitorOptions {
	options.Port = core.Int64Ptr(port)
	return options
}

// SetInterval : Allow user to set Interval
func (options *CreateMonitorOptions) SetInterval(interval int64) *CreateMonitorOptions {
	options.Interval = core.Int64Ptr(interval)
	return options
}

// SetRetries : Allow user to set Retries
func (options *CreateMonitorOptions) SetRetries(retries int64) *CreateMonitorOptions {
	options.Retries = core.Int64Ptr(retries)
	return options
}

// SetTimeout : Allow user to set Timeout
func (options *CreateMonitorOptions) SetTimeout(timeout int64) *CreateMonitorOptions {
	options.Timeout = core.Int64Ptr(timeout)
	return options
}

// SetMethod : Allow user to set Method
func (options *CreateMonitorOptions) SetMethod(method string) *CreateMonitorOptions {
	options.Method = core.StringPtr(method)
	return options
}

// SetPath : Allow user to set Path
func (options *CreateMonitorOptions) SetPath(path string) *CreateMonitorOptions {
	options.Path = core.StringPtr(path)
	return options
}

// SetHeadersVar : Allow user to set HeadersVar
func (options *CreateMonitorOptions) SetHeadersVar(headersVar []HealthcheckHeader) *CreateMonitorOptions {
	options.HeadersVar = headersVar
	return options
}

// SetAllowInsecure : Allow user to set AllowInsecure
func (options *CreateMonitorOptions) SetAllowInsecure(allowInsecure bool) *CreateMonitorOptions {
	options.AllowInsecure = core.BoolPtr(allowInsecure)
	return options
}

// SetExpectedCodes : Allow user to set ExpectedCodes
func (options *CreateMonitorOptions) SetExpectedCodes(expectedCodes string) *CreateMonitorOptions {
	options.ExpectedCodes = core.StringPtr(expectedCodes)
	return options
}

// SetExpectedBody : Allow user to set ExpectedBody
func (options *CreateMonitorOptions) SetExpectedBody(expectedBody string) *CreateMonitorOptions {
	options.ExpectedBody = core.StringPtr(expectedBody)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *CreateMonitorOptions) SetXCorrelationID(xCorrelationID string) *CreateMonitorOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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
	Type *string `json:"type,omitempty"`

	// Permitted network data for VPC.
	PermittedNetwork *PermittedNetworkVpc `json:"permitted_network,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreatePermittedNetworkOptions.Type property.
// The type of a permitted network.
const (
	CreatePermittedNetworkOptions_Type_Vpc = "vpc"
)

// NewCreatePermittedNetworkOptions : Instantiate CreatePermittedNetworkOptions
func (*DnsSvcsV1) NewCreatePermittedNetworkOptions(instanceID string, dnszoneID string) *CreatePermittedNetworkOptions {
	return &CreatePermittedNetworkOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *CreatePermittedNetworkOptions) SetInstanceID(instanceID string) *CreatePermittedNetworkOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *CreatePermittedNetworkOptions) SetDnszoneID(dnszoneID string) *CreatePermittedNetworkOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetType : Allow user to set Type
func (options *CreatePermittedNetworkOptions) SetType(typeVar string) *CreatePermittedNetworkOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetPermittedNetwork : Allow user to set PermittedNetwork
func (options *CreatePermittedNetworkOptions) SetPermittedNetwork(permittedNetwork *PermittedNetworkVpc) *CreatePermittedNetworkOptions {
	options.PermittedNetwork = permittedNetwork
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *CreatePermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *CreatePermittedNetworkOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Health check subnet CRN.
	HealthcheckSubnets []string `json:"healthcheck_subnets,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreatePoolOptions.HealthcheckRegion property.
// Health check region of VSIs.
const (
	CreatePoolOptions_HealthcheckRegion_AuSyd = "au-syd"
	CreatePoolOptions_HealthcheckRegion_EuDu = "eu-du"
	CreatePoolOptions_HealthcheckRegion_EuGb = "eu-gb"
	CreatePoolOptions_HealthcheckRegion_JpTok = "jp-tok"
	CreatePoolOptions_HealthcheckRegion_UsEast = "us-east"
	CreatePoolOptions_HealthcheckRegion_UsSouth = "us-south"
)

// NewCreatePoolOptions : Instantiate CreatePoolOptions
func (*DnsSvcsV1) NewCreatePoolOptions(instanceID string) *CreatePoolOptions {
	return &CreatePoolOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *CreatePoolOptions) SetInstanceID(instanceID string) *CreatePoolOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetName : Allow user to set Name
func (options *CreatePoolOptions) SetName(name string) *CreatePoolOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreatePoolOptions) SetDescription(description string) *CreatePoolOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *CreatePoolOptions) SetEnabled(enabled bool) *CreatePoolOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetHealthyOriginsThreshold : Allow user to set HealthyOriginsThreshold
func (options *CreatePoolOptions) SetHealthyOriginsThreshold(healthyOriginsThreshold int64) *CreatePoolOptions {
	options.HealthyOriginsThreshold = core.Int64Ptr(healthyOriginsThreshold)
	return options
}

// SetOrigins : Allow user to set Origins
func (options *CreatePoolOptions) SetOrigins(origins []OriginInput) *CreatePoolOptions {
	options.Origins = origins
	return options
}

// SetMonitor : Allow user to set Monitor
func (options *CreatePoolOptions) SetMonitor(monitor string) *CreatePoolOptions {
	options.Monitor = core.StringPtr(monitor)
	return options
}

// SetNotificationChannel : Allow user to set NotificationChannel
func (options *CreatePoolOptions) SetNotificationChannel(notificationChannel string) *CreatePoolOptions {
	options.NotificationChannel = core.StringPtr(notificationChannel)
	return options
}

// SetHealthcheckRegion : Allow user to set HealthcheckRegion
func (options *CreatePoolOptions) SetHealthcheckRegion(healthcheckRegion string) *CreatePoolOptions {
	options.HealthcheckRegion = core.StringPtr(healthcheckRegion)
	return options
}

// SetHealthcheckSubnets : Allow user to set HealthcheckSubnets
func (options *CreatePoolOptions) SetHealthcheckSubnets(healthcheckSubnets []string) *CreatePoolOptions {
	options.HealthcheckSubnets = healthcheckSubnets
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *CreatePoolOptions) SetXCorrelationID(xCorrelationID string) *CreatePoolOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Name of the resource record.
	Name *string `json:"name,omitempty"`

	// Type of the resource record.
	Type *string `json:"type,omitempty"`

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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateResourceRecordOptions.Type property.
// Type of the resource record.
const (
	CreateResourceRecordOptions_Type_A = "A"
	CreateResourceRecordOptions_Type_Aaaa = "AAAA"
	CreateResourceRecordOptions_Type_Cname = "CNAME"
	CreateResourceRecordOptions_Type_Mx = "MX"
	CreateResourceRecordOptions_Type_Ptr = "PTR"
	CreateResourceRecordOptions_Type_Srv = "SRV"
	CreateResourceRecordOptions_Type_Txt = "TXT"
)

// NewCreateResourceRecordOptions : Instantiate CreateResourceRecordOptions
func (*DnsSvcsV1) NewCreateResourceRecordOptions(instanceID string, dnszoneID string) *CreateResourceRecordOptions {
	return &CreateResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *CreateResourceRecordOptions) SetInstanceID(instanceID string) *CreateResourceRecordOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *CreateResourceRecordOptions) SetDnszoneID(dnszoneID string) *CreateResourceRecordOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetName : Allow user to set Name
func (options *CreateResourceRecordOptions) SetName(name string) *CreateResourceRecordOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetType : Allow user to set Type
func (options *CreateResourceRecordOptions) SetType(typeVar string) *CreateResourceRecordOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetRdata : Allow user to set Rdata
func (options *CreateResourceRecordOptions) SetRdata(rdata ResourceRecordInputRdataIntf) *CreateResourceRecordOptions {
	options.Rdata = rdata
	return options
}

// SetTTL : Allow user to set TTL
func (options *CreateResourceRecordOptions) SetTTL(ttl int64) *CreateResourceRecordOptions {
	options.TTL = core.Int64Ptr(ttl)
	return options
}

// SetService : Allow user to set Service
func (options *CreateResourceRecordOptions) SetService(service string) *CreateResourceRecordOptions {
	options.Service = core.StringPtr(service)
	return options
}

// SetProtocol : Allow user to set Protocol
func (options *CreateResourceRecordOptions) SetProtocol(protocol string) *CreateResourceRecordOptions {
	options.Protocol = core.StringPtr(protocol)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *CreateResourceRecordOptions) SetXCorrelationID(xCorrelationID string) *CreateResourceRecordOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateResourceRecordOptions) SetHeaders(param map[string]string) *CreateResourceRecordOptions {
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteDnszoneOptions : Instantiate DeleteDnszoneOptions
func (*DnsSvcsV1) NewDeleteDnszoneOptions(instanceID string, dnszoneID string) *DeleteDnszoneOptions {
	return &DeleteDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *DeleteDnszoneOptions) SetInstanceID(instanceID string) *DeleteDnszoneOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *DeleteDnszoneOptions) SetDnszoneID(dnszoneID string) *DeleteDnszoneOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeleteDnszoneOptions) SetXCorrelationID(xCorrelationID string) *DeleteDnszoneOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDnszoneOptions) SetHeaders(param map[string]string) *DeleteDnszoneOptions {
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteLoadBalancerOptions : Instantiate DeleteLoadBalancerOptions
func (*DnsSvcsV1) NewDeleteLoadBalancerOptions(instanceID string, dnszoneID string, lbID string) *DeleteLoadBalancerOptions {
	return &DeleteLoadBalancerOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
		LbID: core.StringPtr(lbID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *DeleteLoadBalancerOptions) SetInstanceID(instanceID string) *DeleteLoadBalancerOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *DeleteLoadBalancerOptions) SetDnszoneID(dnszoneID string) *DeleteLoadBalancerOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetLbID : Allow user to set LbID
func (options *DeleteLoadBalancerOptions) SetLbID(lbID string) *DeleteLoadBalancerOptions {
	options.LbID = core.StringPtr(lbID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeleteLoadBalancerOptions) SetXCorrelationID(xCorrelationID string) *DeleteLoadBalancerOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLoadBalancerOptions) SetHeaders(param map[string]string) *DeleteLoadBalancerOptions {
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteMonitorOptions : Instantiate DeleteMonitorOptions
func (*DnsSvcsV1) NewDeleteMonitorOptions(instanceID string, monitorID string) *DeleteMonitorOptions {
	return &DeleteMonitorOptions{
		InstanceID: core.StringPtr(instanceID),
		MonitorID: core.StringPtr(monitorID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *DeleteMonitorOptions) SetInstanceID(instanceID string) *DeleteMonitorOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetMonitorID : Allow user to set MonitorID
func (options *DeleteMonitorOptions) SetMonitorID(monitorID string) *DeleteMonitorOptions {
	options.MonitorID = core.StringPtr(monitorID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeleteMonitorOptions) SetXCorrelationID(xCorrelationID string) *DeleteMonitorOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeletePermittedNetworkOptions : Instantiate DeletePermittedNetworkOptions
func (*DnsSvcsV1) NewDeletePermittedNetworkOptions(instanceID string, dnszoneID string, permittedNetworkID string) *DeletePermittedNetworkOptions {
	return &DeletePermittedNetworkOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
		PermittedNetworkID: core.StringPtr(permittedNetworkID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *DeletePermittedNetworkOptions) SetInstanceID(instanceID string) *DeletePermittedNetworkOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *DeletePermittedNetworkOptions) SetDnszoneID(dnszoneID string) *DeletePermittedNetworkOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetPermittedNetworkID : Allow user to set PermittedNetworkID
func (options *DeletePermittedNetworkOptions) SetPermittedNetworkID(permittedNetworkID string) *DeletePermittedNetworkOptions {
	options.PermittedNetworkID = core.StringPtr(permittedNetworkID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeletePermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *DeletePermittedNetworkOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeletePoolOptions : Instantiate DeletePoolOptions
func (*DnsSvcsV1) NewDeletePoolOptions(instanceID string, poolID string) *DeletePoolOptions {
	return &DeletePoolOptions{
		InstanceID: core.StringPtr(instanceID),
		PoolID: core.StringPtr(poolID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *DeletePoolOptions) SetInstanceID(instanceID string) *DeletePoolOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetPoolID : Allow user to set PoolID
func (options *DeletePoolOptions) SetPoolID(poolID string) *DeletePoolOptions {
	options.PoolID = core.StringPtr(poolID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeletePoolOptions) SetXCorrelationID(xCorrelationID string) *DeletePoolOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteResourceRecordOptions : Instantiate DeleteResourceRecordOptions
func (*DnsSvcsV1) NewDeleteResourceRecordOptions(instanceID string, dnszoneID string, recordID string) *DeleteResourceRecordOptions {
	return &DeleteResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
		RecordID: core.StringPtr(recordID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *DeleteResourceRecordOptions) SetInstanceID(instanceID string) *DeleteResourceRecordOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *DeleteResourceRecordOptions) SetDnszoneID(dnszoneID string) *DeleteResourceRecordOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetRecordID : Allow user to set RecordID
func (options *DeleteResourceRecordOptions) SetRecordID(recordID string) *DeleteResourceRecordOptions {
	options.RecordID = core.StringPtr(recordID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeleteResourceRecordOptions) SetXCorrelationID(xCorrelationID string) *DeleteResourceRecordOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteResourceRecordOptions) SetHeaders(param map[string]string) *DeleteResourceRecordOptions {
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDnszoneOptions : Instantiate GetDnszoneOptions
func (*DnsSvcsV1) NewGetDnszoneOptions(instanceID string, dnszoneID string) *GetDnszoneOptions {
	return &GetDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *GetDnszoneOptions) SetInstanceID(instanceID string) *GetDnszoneOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *GetDnszoneOptions) SetDnszoneID(dnszoneID string) *GetDnszoneOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *GetDnszoneOptions) SetXCorrelationID(xCorrelationID string) *GetDnszoneOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetDnszoneOptions) SetHeaders(param map[string]string) *GetDnszoneOptions {
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLoadBalancerOptions : Instantiate GetLoadBalancerOptions
func (*DnsSvcsV1) NewGetLoadBalancerOptions(instanceID string, dnszoneID string, lbID string) *GetLoadBalancerOptions {
	return &GetLoadBalancerOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
		LbID: core.StringPtr(lbID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *GetLoadBalancerOptions) SetInstanceID(instanceID string) *GetLoadBalancerOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *GetLoadBalancerOptions) SetDnszoneID(dnszoneID string) *GetLoadBalancerOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetLbID : Allow user to set LbID
func (options *GetLoadBalancerOptions) SetLbID(lbID string) *GetLoadBalancerOptions {
	options.LbID = core.StringPtr(lbID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *GetLoadBalancerOptions) SetXCorrelationID(xCorrelationID string) *GetLoadBalancerOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetMonitorOptions : Instantiate GetMonitorOptions
func (*DnsSvcsV1) NewGetMonitorOptions(instanceID string, monitorID string) *GetMonitorOptions {
	return &GetMonitorOptions{
		InstanceID: core.StringPtr(instanceID),
		MonitorID: core.StringPtr(monitorID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *GetMonitorOptions) SetInstanceID(instanceID string) *GetMonitorOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetMonitorID : Allow user to set MonitorID
func (options *GetMonitorOptions) SetMonitorID(monitorID string) *GetMonitorOptions {
	options.MonitorID = core.StringPtr(monitorID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *GetMonitorOptions) SetXCorrelationID(xCorrelationID string) *GetMonitorOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetPermittedNetworkOptions : Instantiate GetPermittedNetworkOptions
func (*DnsSvcsV1) NewGetPermittedNetworkOptions(instanceID string, dnszoneID string, permittedNetworkID string) *GetPermittedNetworkOptions {
	return &GetPermittedNetworkOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
		PermittedNetworkID: core.StringPtr(permittedNetworkID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *GetPermittedNetworkOptions) SetInstanceID(instanceID string) *GetPermittedNetworkOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *GetPermittedNetworkOptions) SetDnszoneID(dnszoneID string) *GetPermittedNetworkOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetPermittedNetworkID : Allow user to set PermittedNetworkID
func (options *GetPermittedNetworkOptions) SetPermittedNetworkID(permittedNetworkID string) *GetPermittedNetworkOptions {
	options.PermittedNetworkID = core.StringPtr(permittedNetworkID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *GetPermittedNetworkOptions) SetXCorrelationID(xCorrelationID string) *GetPermittedNetworkOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetPoolOptions : Instantiate GetPoolOptions
func (*DnsSvcsV1) NewGetPoolOptions(instanceID string, poolID string) *GetPoolOptions {
	return &GetPoolOptions{
		InstanceID: core.StringPtr(instanceID),
		PoolID: core.StringPtr(poolID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *GetPoolOptions) SetInstanceID(instanceID string) *GetPoolOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetPoolID : Allow user to set PoolID
func (options *GetPoolOptions) SetPoolID(poolID string) *GetPoolOptions {
	options.PoolID = core.StringPtr(poolID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *GetPoolOptions) SetXCorrelationID(xCorrelationID string) *GetPoolOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetResourceRecordOptions : Instantiate GetResourceRecordOptions
func (*DnsSvcsV1) NewGetResourceRecordOptions(instanceID string, dnszoneID string, recordID string) *GetResourceRecordOptions {
	return &GetResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
		RecordID: core.StringPtr(recordID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *GetResourceRecordOptions) SetInstanceID(instanceID string) *GetResourceRecordOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *GetResourceRecordOptions) SetDnszoneID(dnszoneID string) *GetResourceRecordOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetRecordID : Allow user to set RecordID
func (options *GetResourceRecordOptions) SetRecordID(recordID string) *GetResourceRecordOptions {
	options.RecordID = core.StringPtr(recordID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *GetResourceRecordOptions) SetXCorrelationID(xCorrelationID string) *GetResourceRecordOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceRecordOptions) SetHeaders(param map[string]string) *GetResourceRecordOptions {
	options.Headers = param
	return options
}

// ListDnszonesOptions : The ListDnszones options.
type ListDnszonesOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resource records to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify how many resource records are returned, the default value is 200.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListDnszonesOptions : Instantiate ListDnszonesOptions
func (*DnsSvcsV1) NewListDnszonesOptions(instanceID string) *ListDnszonesOptions {
	return &ListDnszonesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *ListDnszonesOptions) SetInstanceID(instanceID string) *ListDnszonesOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListDnszonesOptions) SetXCorrelationID(xCorrelationID string) *ListDnszonesOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListDnszonesOptions) SetOffset(offset int64) *ListDnszonesOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListDnszonesOptions) SetLimit(limit int64) *ListDnszonesOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListDnszonesOptions) SetHeaders(param map[string]string) *ListDnszonesOptions {
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListLoadBalancersOptions : Instantiate ListLoadBalancersOptions
func (*DnsSvcsV1) NewListLoadBalancersOptions(instanceID string, dnszoneID string) *ListLoadBalancersOptions {
	return &ListLoadBalancersOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *ListLoadBalancersOptions) SetInstanceID(instanceID string) *ListLoadBalancersOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *ListLoadBalancersOptions) SetDnszoneID(dnszoneID string) *ListLoadBalancersOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListLoadBalancersOptions) SetXCorrelationID(xCorrelationID string) *ListLoadBalancersOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListMonitorsOptions : Instantiate ListMonitorsOptions
func (*DnsSvcsV1) NewListMonitorsOptions(instanceID string) *ListMonitorsOptions {
	return &ListMonitorsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *ListMonitorsOptions) SetInstanceID(instanceID string) *ListMonitorsOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListMonitorsOptions) SetXCorrelationID(xCorrelationID string) *ListMonitorsOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many resource records to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify how many resource records are returned, the default value is 200.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListPermittedNetworksOptions : Instantiate ListPermittedNetworksOptions
func (*DnsSvcsV1) NewListPermittedNetworksOptions(instanceID string, dnszoneID string) *ListPermittedNetworksOptions {
	return &ListPermittedNetworksOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *ListPermittedNetworksOptions) SetInstanceID(instanceID string) *ListPermittedNetworksOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *ListPermittedNetworksOptions) SetDnszoneID(dnszoneID string) *ListPermittedNetworksOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListPermittedNetworksOptions) SetXCorrelationID(xCorrelationID string) *ListPermittedNetworksOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListPermittedNetworksOptions) SetOffset(offset int64) *ListPermittedNetworksOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListPermittedNetworksOptions) SetLimit(limit int64) *ListPermittedNetworksOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListPoolsOptions : Instantiate ListPoolsOptions
func (*DnsSvcsV1) NewListPoolsOptions(instanceID string) *ListPoolsOptions {
	return &ListPoolsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *ListPoolsOptions) SetInstanceID(instanceID string) *ListPoolsOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListPoolsOptions) SetXCorrelationID(xCorrelationID string) *ListPoolsOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Specify how many resource records to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify how many resource records are returned, the default value is 200.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListResourceRecordsOptions : Instantiate ListResourceRecordsOptions
func (*DnsSvcsV1) NewListResourceRecordsOptions(instanceID string, dnszoneID string) *ListResourceRecordsOptions {
	return &ListResourceRecordsOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *ListResourceRecordsOptions) SetInstanceID(instanceID string) *ListResourceRecordsOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *ListResourceRecordsOptions) SetDnszoneID(dnszoneID string) *ListResourceRecordsOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListResourceRecordsOptions) SetXCorrelationID(xCorrelationID string) *ListResourceRecordsOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListResourceRecordsOptions) SetOffset(offset int64) *ListResourceRecordsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListResourceRecordsOptions) SetLimit(limit int64) *ListResourceRecordsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListResourceRecordsOptions) SetHeaders(param map[string]string) *ListResourceRecordsOptions {
	options.Headers = param
	return options
}

// LoadBalancerAzPoolsItem : LoadBalancerAzPoolsItem struct
type LoadBalancerAzPoolsItem struct {
	// Availability zone.
	AvailabilityZone *string `json:"availability_zone,omitempty"`

	// List of load balancer pools.
	Pools []string `json:"pools,omitempty"`
}


// UnmarshalLoadBalancerAzPoolsItem unmarshals an instance of LoadBalancerAzPoolsItem from the specified map of raw messages.
func UnmarshalLoadBalancerAzPoolsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancerAzPoolsItem)
	err = core.UnmarshalPrimitive(m, "availability_zone", &obj.AvailabilityZone)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pools", &obj.Pools)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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
		return
	}
	err = core.UnmarshalPrimitive(m, "ipv4_address", &obj.Ipv4Address)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ipv4_cidr_block", &obj.Ipv4CidrBlock)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "vpc", &obj.Vpc)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "cname", &obj.Cname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exchange", &obj.Exchange)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "preference", &obj.Preference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "weight", &obj.Weight)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ptrdname", &obj.Ptrdname)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "cname", &obj.Cname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exchange", &obj.Exchange)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "preference", &obj.Preference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "weight", &obj.Weight)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ptrdname", &obj.Ptrdname)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateDnszoneOptions : Instantiate UpdateDnszoneOptions
func (*DnsSvcsV1) NewUpdateDnszoneOptions(instanceID string, dnszoneID string) *UpdateDnszoneOptions {
	return &UpdateDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *UpdateDnszoneOptions) SetInstanceID(instanceID string) *UpdateDnszoneOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *UpdateDnszoneOptions) SetDnszoneID(dnszoneID string) *UpdateDnszoneOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateDnszoneOptions) SetDescription(description string) *UpdateDnszoneOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetLabel : Allow user to set Label
func (options *UpdateDnszoneOptions) SetLabel(label string) *UpdateDnszoneOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *UpdateDnszoneOptions) SetXCorrelationID(xCorrelationID string) *UpdateDnszoneOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDnszoneOptions) SetHeaders(param map[string]string) *UpdateDnszoneOptions {
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
	AzPools []LoadBalancerAzPoolsItem `json:"az_pools,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateLoadBalancerOptions : Instantiate UpdateLoadBalancerOptions
func (*DnsSvcsV1) NewUpdateLoadBalancerOptions(instanceID string, dnszoneID string, lbID string) *UpdateLoadBalancerOptions {
	return &UpdateLoadBalancerOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
		LbID: core.StringPtr(lbID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *UpdateLoadBalancerOptions) SetInstanceID(instanceID string) *UpdateLoadBalancerOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *UpdateLoadBalancerOptions) SetDnszoneID(dnszoneID string) *UpdateLoadBalancerOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetLbID : Allow user to set LbID
func (options *UpdateLoadBalancerOptions) SetLbID(lbID string) *UpdateLoadBalancerOptions {
	options.LbID = core.StringPtr(lbID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateLoadBalancerOptions) SetName(name string) *UpdateLoadBalancerOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateLoadBalancerOptions) SetDescription(description string) *UpdateLoadBalancerOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *UpdateLoadBalancerOptions) SetEnabled(enabled bool) *UpdateLoadBalancerOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetTTL : Allow user to set TTL
func (options *UpdateLoadBalancerOptions) SetTTL(ttl int64) *UpdateLoadBalancerOptions {
	options.TTL = core.Int64Ptr(ttl)
	return options
}

// SetFallbackPool : Allow user to set FallbackPool
func (options *UpdateLoadBalancerOptions) SetFallbackPool(fallbackPool string) *UpdateLoadBalancerOptions {
	options.FallbackPool = core.StringPtr(fallbackPool)
	return options
}

// SetDefaultPools : Allow user to set DefaultPools
func (options *UpdateLoadBalancerOptions) SetDefaultPools(defaultPools []string) *UpdateLoadBalancerOptions {
	options.DefaultPools = defaultPools
	return options
}

// SetAzPools : Allow user to set AzPools
func (options *UpdateLoadBalancerOptions) SetAzPools(azPools []LoadBalancerAzPoolsItem) *UpdateLoadBalancerOptions {
	options.AzPools = azPools
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *UpdateLoadBalancerOptions) SetXCorrelationID(xCorrelationID string) *UpdateLoadBalancerOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateMonitorOptions.Type property.
// The protocol to use for the health check. Currently supported protocols are 'HTTP','HTTPS' and 'TCP'.
const (
	UpdateMonitorOptions_Type_Http = "HTTP"
	UpdateMonitorOptions_Type_Https = "HTTPS"
	UpdateMonitorOptions_Type_Tcp = "TCP"
)

// Constants associated with the UpdateMonitorOptions.Method property.
// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
const (
	UpdateMonitorOptions_Method_Get = "GET"
	UpdateMonitorOptions_Method_Head = "HEAD"
)

// NewUpdateMonitorOptions : Instantiate UpdateMonitorOptions
func (*DnsSvcsV1) NewUpdateMonitorOptions(instanceID string, monitorID string) *UpdateMonitorOptions {
	return &UpdateMonitorOptions{
		InstanceID: core.StringPtr(instanceID),
		MonitorID: core.StringPtr(monitorID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *UpdateMonitorOptions) SetInstanceID(instanceID string) *UpdateMonitorOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetMonitorID : Allow user to set MonitorID
func (options *UpdateMonitorOptions) SetMonitorID(monitorID string) *UpdateMonitorOptions {
	options.MonitorID = core.StringPtr(monitorID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateMonitorOptions) SetName(name string) *UpdateMonitorOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateMonitorOptions) SetDescription(description string) *UpdateMonitorOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetType : Allow user to set Type
func (options *UpdateMonitorOptions) SetType(typeVar string) *UpdateMonitorOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetPort : Allow user to set Port
func (options *UpdateMonitorOptions) SetPort(port int64) *UpdateMonitorOptions {
	options.Port = core.Int64Ptr(port)
	return options
}

// SetInterval : Allow user to set Interval
func (options *UpdateMonitorOptions) SetInterval(interval int64) *UpdateMonitorOptions {
	options.Interval = core.Int64Ptr(interval)
	return options
}

// SetRetries : Allow user to set Retries
func (options *UpdateMonitorOptions) SetRetries(retries int64) *UpdateMonitorOptions {
	options.Retries = core.Int64Ptr(retries)
	return options
}

// SetTimeout : Allow user to set Timeout
func (options *UpdateMonitorOptions) SetTimeout(timeout int64) *UpdateMonitorOptions {
	options.Timeout = core.Int64Ptr(timeout)
	return options
}

// SetMethod : Allow user to set Method
func (options *UpdateMonitorOptions) SetMethod(method string) *UpdateMonitorOptions {
	options.Method = core.StringPtr(method)
	return options
}

// SetPath : Allow user to set Path
func (options *UpdateMonitorOptions) SetPath(path string) *UpdateMonitorOptions {
	options.Path = core.StringPtr(path)
	return options
}

// SetHeadersVar : Allow user to set HeadersVar
func (options *UpdateMonitorOptions) SetHeadersVar(headersVar []HealthcheckHeader) *UpdateMonitorOptions {
	options.HeadersVar = headersVar
	return options
}

// SetAllowInsecure : Allow user to set AllowInsecure
func (options *UpdateMonitorOptions) SetAllowInsecure(allowInsecure bool) *UpdateMonitorOptions {
	options.AllowInsecure = core.BoolPtr(allowInsecure)
	return options
}

// SetExpectedCodes : Allow user to set ExpectedCodes
func (options *UpdateMonitorOptions) SetExpectedCodes(expectedCodes string) *UpdateMonitorOptions {
	options.ExpectedCodes = core.StringPtr(expectedCodes)
	return options
}

// SetExpectedBody : Allow user to set ExpectedBody
func (options *UpdateMonitorOptions) SetExpectedBody(expectedBody string) *UpdateMonitorOptions {
	options.ExpectedBody = core.StringPtr(expectedBody)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *UpdateMonitorOptions) SetXCorrelationID(xCorrelationID string) *UpdateMonitorOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdatePoolOptions.HealthcheckRegion property.
// Health check region of VSIs.
const (
	UpdatePoolOptions_HealthcheckRegion_AuSyd = "au-syd"
	UpdatePoolOptions_HealthcheckRegion_EuDu = "eu-du"
	UpdatePoolOptions_HealthcheckRegion_EuGb = "eu-gb"
	UpdatePoolOptions_HealthcheckRegion_JpTok = "jp-tok"
	UpdatePoolOptions_HealthcheckRegion_UsEast = "us-east"
	UpdatePoolOptions_HealthcheckRegion_UsSouth = "us-south"
)

// NewUpdatePoolOptions : Instantiate UpdatePoolOptions
func (*DnsSvcsV1) NewUpdatePoolOptions(instanceID string, poolID string) *UpdatePoolOptions {
	return &UpdatePoolOptions{
		InstanceID: core.StringPtr(instanceID),
		PoolID: core.StringPtr(poolID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *UpdatePoolOptions) SetInstanceID(instanceID string) *UpdatePoolOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetPoolID : Allow user to set PoolID
func (options *UpdatePoolOptions) SetPoolID(poolID string) *UpdatePoolOptions {
	options.PoolID = core.StringPtr(poolID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdatePoolOptions) SetName(name string) *UpdatePoolOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdatePoolOptions) SetDescription(description string) *UpdatePoolOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *UpdatePoolOptions) SetEnabled(enabled bool) *UpdatePoolOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetHealthyOriginsThreshold : Allow user to set HealthyOriginsThreshold
func (options *UpdatePoolOptions) SetHealthyOriginsThreshold(healthyOriginsThreshold int64) *UpdatePoolOptions {
	options.HealthyOriginsThreshold = core.Int64Ptr(healthyOriginsThreshold)
	return options
}

// SetOrigins : Allow user to set Origins
func (options *UpdatePoolOptions) SetOrigins(origins []OriginInput) *UpdatePoolOptions {
	options.Origins = origins
	return options
}

// SetMonitor : Allow user to set Monitor
func (options *UpdatePoolOptions) SetMonitor(monitor string) *UpdatePoolOptions {
	options.Monitor = core.StringPtr(monitor)
	return options
}

// SetNotificationChannel : Allow user to set NotificationChannel
func (options *UpdatePoolOptions) SetNotificationChannel(notificationChannel string) *UpdatePoolOptions {
	options.NotificationChannel = core.StringPtr(notificationChannel)
	return options
}

// SetHealthcheckRegion : Allow user to set HealthcheckRegion
func (options *UpdatePoolOptions) SetHealthcheckRegion(healthcheckRegion string) *UpdatePoolOptions {
	options.HealthcheckRegion = core.StringPtr(healthcheckRegion)
	return options
}

// SetHealthcheckSubnets : Allow user to set HealthcheckSubnets
func (options *UpdatePoolOptions) SetHealthcheckSubnets(healthcheckSubnets []string) *UpdatePoolOptions {
	options.HealthcheckSubnets = healthcheckSubnets
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *UpdatePoolOptions) SetXCorrelationID(xCorrelationID string) *UpdatePoolOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
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
	Name *string `json:"name,omitempty"`

	// Content of the resource record.
	Rdata ResourceRecordUpdateInputRdataIntf `json:"rdata,omitempty"`

	// Time to live in second.
	TTL *int64 `json:"ttl,omitempty"`

	// Only used for SRV record.
	Service *string `json:"service,omitempty"`

	// Only used for SRV record.
	Protocol *string `json:"protocol,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateResourceRecordOptions : Instantiate UpdateResourceRecordOptions
func (*DnsSvcsV1) NewUpdateResourceRecordOptions(instanceID string, dnszoneID string, recordID string) *UpdateResourceRecordOptions {
	return &UpdateResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID: core.StringPtr(dnszoneID),
		RecordID: core.StringPtr(recordID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (options *UpdateResourceRecordOptions) SetInstanceID(instanceID string) *UpdateResourceRecordOptions {
	options.InstanceID = core.StringPtr(instanceID)
	return options
}

// SetDnszoneID : Allow user to set DnszoneID
func (options *UpdateResourceRecordOptions) SetDnszoneID(dnszoneID string) *UpdateResourceRecordOptions {
	options.DnszoneID = core.StringPtr(dnszoneID)
	return options
}

// SetRecordID : Allow user to set RecordID
func (options *UpdateResourceRecordOptions) SetRecordID(recordID string) *UpdateResourceRecordOptions {
	options.RecordID = core.StringPtr(recordID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateResourceRecordOptions) SetName(name string) *UpdateResourceRecordOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetRdata : Allow user to set Rdata
func (options *UpdateResourceRecordOptions) SetRdata(rdata ResourceRecordUpdateInputRdataIntf) *UpdateResourceRecordOptions {
	options.Rdata = rdata
	return options
}

// SetTTL : Allow user to set TTL
func (options *UpdateResourceRecordOptions) SetTTL(ttl int64) *UpdateResourceRecordOptions {
	options.TTL = core.Int64Ptr(ttl)
	return options
}

// SetService : Allow user to set Service
func (options *UpdateResourceRecordOptions) SetService(service string) *UpdateResourceRecordOptions {
	options.Service = core.StringPtr(service)
	return options
}

// SetProtocol : Allow user to set Protocol
func (options *UpdateResourceRecordOptions) SetProtocol(protocol string) *UpdateResourceRecordOptions {
	options.Protocol = core.StringPtr(protocol)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *UpdateResourceRecordOptions) SetXCorrelationID(xCorrelationID string) *UpdateResourceRecordOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateResourceRecordOptions) SetHeaders(param map[string]string) *UpdateResourceRecordOptions {
	options.Headers = param
	return options
}

// Dnszone : DNS zone details.
type Dnszone struct {
	// Unique identifier of a DNS zone.
	ID *string `json:"id,omitempty"`

	// the time when a DNS zone is created.
	CreatedOn *string `json:"created_on,omitempty"`

	// the recent time when a DNS zone is modified.
	ModifiedOn *string `json:"modified_on,omitempty"`

	// Unique identifier of a service instance.
	InstanceID *string `json:"instance_id,omitempty"`

	// Name of DNS zone.
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
	Dnszone_State_Active = "active"
	Dnszone_State_Deleted = "deleted"
	Dnszone_State_Disabled = "disabled"
	Dnszone_State_PendingDelete = "pending_delete"
	Dnszone_State_PendingNetworkAdd = "pending_network_add"
)


// UnmarshalDnszone unmarshals an instance of Dnszone from the specified map of raw messages.
func UnmarshalDnszone(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Dnszone)
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
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirstHref : href.
type FirstHref struct {
	// href.
	Href *string `json:"href,omitempty"`
}


// UnmarshalFirstHref unmarshals an instance of FirstHref from the specified map of raw messages.
func UnmarshalFirstHref(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirstHref)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// HealthcheckHeader : The HTTP header of health check request.
type HealthcheckHeader struct {
	// The name of HTTP request header.
	Name *string `json:"name" validate:"required"`

	// The value of HTTP request header.
	Value []string `json:"value" validate:"required"`
}


// NewHealthcheckHeader : Instantiate HealthcheckHeader (Generic Model Constructor)
func (*DnsSvcsV1) NewHealthcheckHeader(name string, value []string) (model *HealthcheckHeader, err error) {
	model = &HealthcheckHeader{
		Name: core.StringPtr(name),
		Value: value,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalHealthcheckHeader unmarshals an instance of HealthcheckHeader from the specified map of raw messages.
func UnmarshalHealthcheckHeader(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HealthcheckHeader)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListDnszones : List DNS zones response.
type ListDnszones struct {
	// An array of DNS zones.
	Dnszones []Dnszone `json:"dnszones" validate:"required"`

	// Specify how many DNS zones to skip over, the default value is 0.
	Offset *int64 `json:"offset" validate:"required"`

	// Specify how many DNS zones are returned, the default value is 10.
	Limit *int64 `json:"limit" validate:"required"`

	// Total number of DNS zones.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *FirstHref `json:"first" validate:"required"`

	// href.
	Next *NextHref `json:"next,omitempty"`
}


// UnmarshalListDnszones unmarshals an instance of ListDnszones from the specified map of raw messages.
func UnmarshalListDnszones(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListDnszones)
	err = core.UnmarshalModel(m, "dnszones", &obj.Dnszones, UnmarshalDnszone)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirstHref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListLoadBalancers : List Global Load Balancers response.
type ListLoadBalancers struct {
	// An array of Global Load Balancers.
	LoadBalancers []LoadBalancer `json:"load_balancers" validate:"required"`

	// Page number.
	Offset *int64 `json:"offset" validate:"required"`

	// Number of Global Load Balancers per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Number of Global Load Balancers.
	Count *int64 `json:"count" validate:"required"`

	// Total number of Global Load Balancers.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *FirstHref `json:"first" validate:"required"`

	// href.
	Next *NextHref `json:"next" validate:"required"`
}


// UnmarshalListLoadBalancers unmarshals an instance of ListLoadBalancers from the specified map of raw messages.
func UnmarshalListLoadBalancers(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListLoadBalancers)
	err = core.UnmarshalModel(m, "load_balancers", &obj.LoadBalancers, UnmarshalLoadBalancer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirstHref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListMonitors : List load balancer monitors response.
type ListMonitors struct {
	// An array of load balancer monitors.
	Monitors []Monitor `json:"monitors" validate:"required"`

	// Page number.
	Offset *int64 `json:"offset" validate:"required"`

	// Number of load balancer monitors per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Number of load balancers.
	Count *int64 `json:"count" validate:"required"`

	// Total number of load balancers.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *FirstHref `json:"first" validate:"required"`

	// href.
	Next *NextHref `json:"next" validate:"required"`
}


// UnmarshalListMonitors unmarshals an instance of ListMonitors from the specified map of raw messages.
func UnmarshalListMonitors(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListMonitors)
	err = core.UnmarshalModel(m, "monitors", &obj.Monitors, UnmarshalMonitor)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirstHref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListPermittedNetworks : List permitted networks response.
type ListPermittedNetworks struct {
	// An array of permitted networks.
	PermittedNetworks []PermittedNetwork `json:"permitted_networks" validate:"required"`

	// Specify how many permitted networks to skip over, the default value is 0.
	Offset *int64 `json:"offset" validate:"required"`

	// Specify how many permitted networks are returned, the default value is 10.
	Limit *int64 `json:"limit" validate:"required"`

	// Total number of permitted networks.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *FirstHref `json:"first" validate:"required"`

	// href.
	Next *NextHref `json:"next,omitempty"`
}


// UnmarshalListPermittedNetworks unmarshals an instance of ListPermittedNetworks from the specified map of raw messages.
func UnmarshalListPermittedNetworks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListPermittedNetworks)
	err = core.UnmarshalModel(m, "permitted_networks", &obj.PermittedNetworks, UnmarshalPermittedNetwork)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirstHref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListPools : List load balancer pools response.
type ListPools struct {
	// An array of load balancer pools.
	Pools []Pool `json:"pools" validate:"required"`

	// Page number.
	Offset *int64 `json:"offset" validate:"required"`

	// Number of load balancer pools per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Number of load balancers.
	Count *int64 `json:"count" validate:"required"`

	// Total number of load balancers.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *FirstHref `json:"first" validate:"required"`

	// href.
	Next *NextHref `json:"next" validate:"required"`
}


// UnmarshalListPools unmarshals an instance of ListPools from the specified map of raw messages.
func UnmarshalListPools(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListPools)
	err = core.UnmarshalModel(m, "pools", &obj.Pools, UnmarshalPool)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirstHref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListResourceRecords : List Resource Records response.
type ListResourceRecords struct {
	// An array of resource records.
	ResourceRecords []ResourceRecord `json:"resource_records" validate:"required"`

	// Specify how many resource records to skip over, the default value is 0.
	Offset *int64 `json:"offset" validate:"required"`

	// Specify how many resource records are returned, the default value is 20.
	Limit *int64 `json:"limit" validate:"required"`

	// Total number of resource records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// href.
	First *FirstHref `json:"first" validate:"required"`

	// href.
	Next *NextHref `json:"next,omitempty"`
}


// UnmarshalListResourceRecords unmarshals an instance of ListResourceRecords from the specified map of raw messages.
func UnmarshalListResourceRecords(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListResourceRecords)
	err = core.UnmarshalModel(m, "resource_records", &obj.ResourceRecords, UnmarshalResourceRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirstHref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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
	AzPools []LoadBalancerAzPoolsItem `json:"az_pools,omitempty"`

	// The time when a load balancer is created.
	CreatedOn *string `json:"created_on,omitempty"`

	// The recent time when a load balancer is modified.
	ModifiedOn *string `json:"modified_on,omitempty"`
}

// Constants associated with the LoadBalancer.Health property.
// Healthy state of the load balancer.
const (
	LoadBalancer_Health_Critical = "CRITICAL"
	LoadBalancer_Health_Degraded = "DEGRADED"
	LoadBalancer_Health_Healthy = "HEALTHY"
)


// UnmarshalLoadBalancer unmarshals an instance of LoadBalancer from the specified map of raw messages.
func UnmarshalLoadBalancer(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancer)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "health", &obj.Health)
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
	err = core.UnmarshalModel(m, "az_pools", &obj.AzPools, UnmarshalLoadBalancerAzPoolsItem)
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
	CreatedOn *string `json:"created_on,omitempty"`

	// the recent time when a load balancer monitor is modified.
	ModifiedOn *string `json:"modified_on,omitempty"`
}

// Constants associated with the Monitor.Method property.
// The method to use for the health check applicable to HTTP/HTTPS based checks, the default value is 'GET'.
const (
	Monitor_Method_Get = "GET"
	Monitor_Method_Head = "HEAD"
)


// UnmarshalMonitor unmarshals an instance of Monitor from the specified map of raw messages.
func UnmarshalMonitor(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Monitor)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "retries", &obj.Retries)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "headers", &obj.HeadersVar, UnmarshalHealthcheckHeader)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_insecure", &obj.AllowInsecure)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_codes", &obj.ExpectedCodes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_body", &obj.ExpectedBody)
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

// NextHref : href.
type NextHref struct {
	// href.
	Href *string `json:"href,omitempty"`
}


// UnmarshalNextHref unmarshals an instance of NextHref from the specified map of raw messages.
func UnmarshalNextHref(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NextHref)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
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
	err = core.UnmarshalPrimitive(m, "health", &obj.Health)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "health_failure_reason", &obj.HealthFailureReason)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PermittedNetwork : Permitted network details.
type PermittedNetwork struct {
	// Unique identifier of a permitted network.
	ID *string `json:"id,omitempty"`

	// The time when a permitted network is created.
	CreatedOn *string `json:"created_on,omitempty"`

	// The recent time when a permitted network is modified.
	ModifiedOn *string `json:"modified_on,omitempty"`

	// Permitted network data for VPC.
	PermittedNetwork *PermittedNetworkVpc `json:"permitted_network,omitempty"`

	// The type of a permitted network.
	Type *string `json:"type,omitempty"`

	// The state of a permitted network.
	State *string `json:"state,omitempty"`
}

// Constants associated with the PermittedNetwork.Type property.
// The type of a permitted network.
const (
	PermittedNetwork_Type_Vpc = "vpc"
)

// Constants associated with the PermittedNetwork.State property.
// The state of a permitted network.
const (
	PermittedNetwork_State_Active = "ACTIVE"
	PermittedNetwork_State_RemovalInProgress = "REMOVAL_IN_PROGRESS"
)


// UnmarshalPermittedNetwork unmarshals an instance of PermittedNetwork from the specified map of raw messages.
func UnmarshalPermittedNetwork(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PermittedNetwork)
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
	err = core.UnmarshalModel(m, "permitted_network", &obj.PermittedNetwork, UnmarshalPermittedNetworkVpc)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
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
func (*DnsSvcsV1) NewPermittedNetworkVpc(vpcCrn string) (model *PermittedNetworkVpc, err error) {
	model = &PermittedNetworkVpc{
		VpcCrn: core.StringPtr(vpcCrn),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalPermittedNetworkVpc unmarshals an instance of PermittedNetworkVpc from the specified map of raw messages.
func UnmarshalPermittedNetworkVpc(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PermittedNetworkVpc)
	err = core.UnmarshalPrimitive(m, "vpc_crn", &obj.VpcCrn)
	if err != nil {
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
	CreatedOn *string `json:"created_on,omitempty"`

	// the recent time when a load balancer pool is modified.
	ModifiedOn *string `json:"modified_on,omitempty"`
}

// Constants associated with the Pool.Health property.
// Healthy state of the load balancer pool.
const (
	Pool_Health_Critical = "CRITICAL"
	Pool_Health_Degraded = "DEGRADED"
	Pool_Health_Healthy = "HEALTHY"
)

// Constants associated with the Pool.HealthcheckRegion property.
// Health check region of VSIs.
const (
	Pool_HealthcheckRegion_AuSyd = "au-syd"
	Pool_HealthcheckRegion_EuDu = "eu-du"
	Pool_HealthcheckRegion_EuGb = "eu-gb"
	Pool_HealthcheckRegion_JpTok = "jp-tok"
	Pool_HealthcheckRegion_UsEast = "us-east"
	Pool_HealthcheckRegion_UsSouth = "us-south"
)


// UnmarshalPool unmarshals an instance of Pool from the specified map of raw messages.
func UnmarshalPool(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Pool)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "healthy_origins_threshold", &obj.HealthyOriginsThreshold)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "origins", &obj.Origins, UnmarshalOrigin)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "monitor", &obj.Monitor)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "notification_channel", &obj.NotificationChannel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "health", &obj.Health)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "healthcheck_region", &obj.HealthcheckRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "healthcheck_subnets", &obj.HealthcheckSubnets)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "healthcheck_vsis", &obj.HealthcheckVsis, UnmarshalPoolHealthcheckVsisItem)
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

// ResourceRecord : Resource record details.
type ResourceRecord struct {
	// Identifier of the resource record.
	ID *string `json:"id,omitempty"`

	// the time when a resource record is created.
	CreatedOn *string `json:"created_on,omitempty"`

	// the recent time when a resource record is modified.
	ModifiedOn *string `json:"modified_on,omitempty"`

	// Name of the resource record.
	Name *string `json:"name,omitempty"`

	// Type of the resource record.
	Type *string `json:"type,omitempty"`

	// Time to live in second.
	TTL *int64 `json:"ttl,omitempty"`

	// Content of the resource record.
	Rdata interface{} `json:"rdata,omitempty"`

	// Only used for SRV record.
	Service *string `json:"service,omitempty"`

	// Only used for SRV record.
	Protocol *string `json:"protocol,omitempty"`
}

// Constants associated with the ResourceRecord.Type property.
// Type of the resource record.
const (
	ResourceRecord_Type_A = "A"
	ResourceRecord_Type_Aaaa = "AAAA"
	ResourceRecord_Type_Cname = "CNAME"
	ResourceRecord_Type_Mx = "MX"
	ResourceRecord_Type_Ptr = "PTR"
	ResourceRecord_Type_Srv = "SRV"
	ResourceRecord_Type_Txt = "TXT"
)


// UnmarshalResourceRecord unmarshals an instance of ResourceRecord from the specified map of raw messages.
func UnmarshalResourceRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceRecord)
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rdata", &obj.Rdata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service", &obj.Service)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "protocol", &obj.Protocol)
	if err != nil {
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
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataARecord(ip string) (model *ResourceRecordInputRdataRdataARecord, err error) {
	model = &ResourceRecordInputRdataRdataARecord{
		Ip: core.StringPtr(ip),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataAaaaRecord(ip string) (model *ResourceRecordInputRdataRdataAaaaRecord, err error) {
	model = &ResourceRecordInputRdataRdataAaaaRecord{
		Ip: core.StringPtr(ip),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataCnameRecord(cname string) (model *ResourceRecordInputRdataRdataCnameRecord, err error) {
	model = &ResourceRecordInputRdataRdataCnameRecord{
		Cname: core.StringPtr(cname),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataMxRecord(exchange string, preference int64) (model *ResourceRecordInputRdataRdataMxRecord, err error) {
	model = &ResourceRecordInputRdataRdataMxRecord{
		Exchange: core.StringPtr(exchange),
		Preference: core.Int64Ptr(preference),
	}
	err = core.ValidateStruct(model, "required parameters")
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
		return
	}
	err = core.UnmarshalPrimitive(m, "preference", &obj.Preference)
	if err != nil {
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
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataPtrRecord(ptrdname string) (model *ResourceRecordInputRdataRdataPtrRecord, err error) {
	model = &ResourceRecordInputRdataRdataPtrRecord{
		Ptrdname: core.StringPtr(ptrdname),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataSrvRecord(port int64, priority int64, target string, weight int64) (model *ResourceRecordInputRdataRdataSrvRecord, err error) {
	model = &ResourceRecordInputRdataRdataSrvRecord{
		Port: core.Int64Ptr(port),
		Priority: core.Int64Ptr(priority),
		Target: core.StringPtr(target),
		Weight: core.Int64Ptr(weight),
	}
	err = core.ValidateStruct(model, "required parameters")
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
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
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

// ResourceRecordInputRdataRdataTxtRecord : The content of type-TXT resource record.
// This model "extends" ResourceRecordInputRdata
type ResourceRecordInputRdataRdataTxtRecord struct {
	// Human readable text.
	Text *string `json:"text" validate:"required"`
}


// NewResourceRecordInputRdataRdataTxtRecord : Instantiate ResourceRecordInputRdataRdataTxtRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordInputRdataRdataTxtRecord(text string) (model *ResourceRecordInputRdataRdataTxtRecord, err error) {
	model = &ResourceRecordInputRdataRdataTxtRecord{
		Text: core.StringPtr(text),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataARecord(ip string) (model *ResourceRecordUpdateInputRdataRdataARecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataARecord{
		Ip: core.StringPtr(ip),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataAaaaRecord(ip string) (model *ResourceRecordUpdateInputRdataRdataAaaaRecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataAaaaRecord{
		Ip: core.StringPtr(ip),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataCnameRecord(cname string) (model *ResourceRecordUpdateInputRdataRdataCnameRecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataCnameRecord{
		Cname: core.StringPtr(cname),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataMxRecord(exchange string, preference int64) (model *ResourceRecordUpdateInputRdataRdataMxRecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataMxRecord{
		Exchange: core.StringPtr(exchange),
		Preference: core.Int64Ptr(preference),
	}
	err = core.ValidateStruct(model, "required parameters")
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
		return
	}
	err = core.UnmarshalPrimitive(m, "preference", &obj.Preference)
	if err != nil {
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
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataPtrRecord(ptrdname string) (model *ResourceRecordUpdateInputRdataRdataPtrRecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataPtrRecord{
		Ptrdname: core.StringPtr(ptrdname),
	}
	err = core.ValidateStruct(model, "required parameters")
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
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataSrvRecord(port int64, priority int64, target string, weight int64) (model *ResourceRecordUpdateInputRdataRdataSrvRecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataSrvRecord{
		Port: core.Int64Ptr(port),
		Priority: core.Int64Ptr(priority),
		Target: core.StringPtr(target),
		Weight: core.Int64Ptr(weight),
	}
	err = core.ValidateStruct(model, "required parameters")
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
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
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

// ResourceRecordUpdateInputRdataRdataTxtRecord : The content of type-TXT resource record.
// This model "extends" ResourceRecordUpdateInputRdata
type ResourceRecordUpdateInputRdataRdataTxtRecord struct {
	// Human readable text.
	Text *string `json:"text" validate:"required"`
}


// NewResourceRecordUpdateInputRdataRdataTxtRecord : Instantiate ResourceRecordUpdateInputRdataRdataTxtRecord (Generic Model Constructor)
func (*DnsSvcsV1) NewResourceRecordUpdateInputRdataRdataTxtRecord(text string) (model *ResourceRecordUpdateInputRdataRdataTxtRecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataTxtRecord{
		Text: core.StringPtr(text),
	}
	err = core.ValidateStruct(model, "required parameters")
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
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
