/**
 * (C) Copyright IBM Corp. 2026.
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
 * IBM OpenAPI SDK Code Generator Version: 3.106.0-09823488-20250707-071701
 */

// Package dnsrecordsv1 : Operations and models for the DnsRecordsV1 service
package dnsrecordsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
)

// DnsRecordsV1 : DNS records
//
// API Version: 1.0.1
type DnsRecordsV1 struct {
	Service *core.BaseService

	// Full crn of the service instance.
	Crn *string

	// Zone identifier (zone id).
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "dns_records"

// DnsRecordsV1Options : Service options
type DnsRecordsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full crn of the service instance.
	Crn *string `validate:"required"`

	// Zone identifier (zone id).
	ZoneIdentifier *string `validate:"required"`
}

// NewDnsRecordsV1UsingExternalConfig : constructs an instance of DnsRecordsV1 with passed in options and external configuration.
func NewDnsRecordsV1UsingExternalConfig(options *DnsRecordsV1Options) (dnsRecords *DnsRecordsV1, err error) {
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

	dnsRecords, err = NewDnsRecordsV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = dnsRecords.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = dnsRecords.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewDnsRecordsV1 : constructs an instance of DnsRecordsV1 with passed in options.
func NewDnsRecordsV1(options *DnsRecordsV1Options) (service *DnsRecordsV1, err error) {
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

	service = &DnsRecordsV1{
		Service:        baseService,
		Crn:            options.Crn,
		ZoneIdentifier: options.ZoneIdentifier,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "dnsRecords" suitable for processing requests.
func (dnsRecords *DnsRecordsV1) Clone() *DnsRecordsV1 {
	if core.IsNil(dnsRecords) {
		return nil
	}
	clone := *dnsRecords
	clone.Service = dnsRecords.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (dnsRecords *DnsRecordsV1) SetServiceURL(url string) error {
	err := dnsRecords.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (dnsRecords *DnsRecordsV1) GetServiceURL() string {
	return dnsRecords.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (dnsRecords *DnsRecordsV1) SetDefaultHeaders(headers http.Header) {
	dnsRecords.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (dnsRecords *DnsRecordsV1) SetEnableGzipCompression(enableGzip bool) {
	dnsRecords.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (dnsRecords *DnsRecordsV1) GetEnableGzipCompression() bool {
	return dnsRecords.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (dnsRecords *DnsRecordsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	dnsRecords.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (dnsRecords *DnsRecordsV1) DisableRetries() {
	dnsRecords.Service.DisableRetries()
}

// ListAllDnsRecords : List all DNS records
// List all DNS records for a given zone of a service instance.
func (dnsRecords *DnsRecordsV1) ListAllDnsRecords(listAllDnsRecordsOptions *ListAllDnsRecordsOptions) (result *ListDnsrecordsResp, response *core.DetailedResponse, err error) {
	result, response, err = dnsRecords.ListAllDnsRecordsWithContext(context.Background(), listAllDnsRecordsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAllDnsRecordsWithContext is an alternate form of the ListAllDnsRecords method which supports a Context parameter
func (dnsRecords *DnsRecordsV1) ListAllDnsRecordsWithContext(ctx context.Context, listAllDnsRecordsOptions *ListAllDnsRecordsOptions) (result *ListDnsrecordsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllDnsRecordsOptions, "listAllDnsRecordsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *dnsRecords.Crn,
		"zone_identifier": *dnsRecords.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsRecords.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsRecords.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/dns_records`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	sdkHeaders := common.GetSdkHeaders("dns_records", "V1", "ListAllDnsRecords")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	for headerName, headerValue := range listAllDnsRecordsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAllDnsRecordsOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listAllDnsRecordsOptions.Type))
	}
	if listAllDnsRecordsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listAllDnsRecordsOptions.Name))
	}
	if listAllDnsRecordsOptions.Content != nil {
		builder.AddQuery("content", fmt.Sprint(*listAllDnsRecordsOptions.Content))
	}
	if listAllDnsRecordsOptions.Page != nil {
		builder.AddQuery("page", fmt.Sprint(*listAllDnsRecordsOptions.Page))
	}
	if listAllDnsRecordsOptions.PerPage != nil {
		builder.AddQuery("per_page", fmt.Sprint(*listAllDnsRecordsOptions.PerPage))
	}
	if listAllDnsRecordsOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listAllDnsRecordsOptions.Order))
	}
	if listAllDnsRecordsOptions.Direction != nil {
		builder.AddQuery("direction", fmt.Sprint(*listAllDnsRecordsOptions.Direction))
	}
	if listAllDnsRecordsOptions.Match != nil {
		builder.AddQuery("match", fmt.Sprint(*listAllDnsRecordsOptions.Match))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsRecords.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_all_dns_records", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListDnsrecordsResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateDnsRecord : Create DNS record
// Add a new DNS record for a given zone for a given service instance.
func (dnsRecords *DnsRecordsV1) CreateDnsRecord(createDnsRecordOptions *CreateDnsRecordOptions) (result *DnsrecordResp, response *core.DetailedResponse, err error) {
	result, response, err = dnsRecords.CreateDnsRecordWithContext(context.Background(), createDnsRecordOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateDnsRecordWithContext is an alternate form of the CreateDnsRecord method which supports a Context parameter
func (dnsRecords *DnsRecordsV1) CreateDnsRecordWithContext(ctx context.Context, createDnsRecordOptions *CreateDnsRecordOptions) (result *DnsrecordResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createDnsRecordOptions, "createDnsRecordOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *dnsRecords.Crn,
		"zone_identifier": *dnsRecords.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsRecords.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsRecords.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/dns_records`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	sdkHeaders := common.GetSdkHeaders("dns_records", "V1", "CreateDnsRecord")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	for headerName, headerValue := range createDnsRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createDnsRecordOptions.Name != nil {
		body["name"] = createDnsRecordOptions.Name
	}
	if createDnsRecordOptions.Type != nil {
		body["type"] = createDnsRecordOptions.Type
	}
	if createDnsRecordOptions.TTL != nil {
		body["ttl"] = createDnsRecordOptions.TTL
	}
	if createDnsRecordOptions.Content != nil {
		body["content"] = createDnsRecordOptions.Content
	}
	if createDnsRecordOptions.Priority != nil {
		body["priority"] = createDnsRecordOptions.Priority
	}
	if createDnsRecordOptions.Proxied != nil {
		body["proxied"] = createDnsRecordOptions.Proxied
	}
	if createDnsRecordOptions.Data != nil {
		body["data"] = createDnsRecordOptions.Data
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
	response, err = dnsRecords.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_dns_record", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnsrecordResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteDnsRecord : Delete DNS record
// Delete a DNS record given its id.
func (dnsRecords *DnsRecordsV1) DeleteDnsRecord(deleteDnsRecordOptions *DeleteDnsRecordOptions) (result *DeleteDnsrecordResp, response *core.DetailedResponse, err error) {
	result, response, err = dnsRecords.DeleteDnsRecordWithContext(context.Background(), deleteDnsRecordOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteDnsRecordWithContext is an alternate form of the DeleteDnsRecord method which supports a Context parameter
func (dnsRecords *DnsRecordsV1) DeleteDnsRecordWithContext(ctx context.Context, deleteDnsRecordOptions *DeleteDnsRecordOptions) (result *DeleteDnsrecordResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDnsRecordOptions, "deleteDnsRecordOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteDnsRecordOptions, "deleteDnsRecordOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":                  *dnsRecords.Crn,
		"zone_identifier":      *dnsRecords.ZoneIdentifier,
		"dnsrecord_identifier": *deleteDnsRecordOptions.DnsrecordIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsRecords.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsRecords.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/dns_records/{dnsrecord_identifier}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	sdkHeaders := common.GetSdkHeaders("dns_records", "V1", "DeleteDnsRecord")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	for headerName, headerValue := range deleteDnsRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsRecords.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_dns_record", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteDnsrecordResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDnsRecord : Get DNS record
// Get the details of a DNS record for a given zone under a given service instance.
func (dnsRecords *DnsRecordsV1) GetDnsRecord(getDnsRecordOptions *GetDnsRecordOptions) (result *DnsrecordResp, response *core.DetailedResponse, err error) {
	result, response, err = dnsRecords.GetDnsRecordWithContext(context.Background(), getDnsRecordOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDnsRecordWithContext is an alternate form of the GetDnsRecord method which supports a Context parameter
func (dnsRecords *DnsRecordsV1) GetDnsRecordWithContext(ctx context.Context, getDnsRecordOptions *GetDnsRecordOptions) (result *DnsrecordResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDnsRecordOptions, "getDnsRecordOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDnsRecordOptions, "getDnsRecordOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":                  *dnsRecords.Crn,
		"zone_identifier":      *dnsRecords.ZoneIdentifier,
		"dnsrecord_identifier": *getDnsRecordOptions.DnsrecordIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsRecords.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsRecords.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/dns_records/{dnsrecord_identifier}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	sdkHeaders := common.GetSdkHeaders("dns_records", "V1", "GetDnsRecord")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	for headerName, headerValue := range getDnsRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsRecords.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_dns_record", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnsrecordResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateDnsRecord : Update DNS record
// Update an existing DNS record for a given zone under a given service instance.
func (dnsRecords *DnsRecordsV1) UpdateDnsRecord(updateDnsRecordOptions *UpdateDnsRecordOptions) (result *DnsrecordResp, response *core.DetailedResponse, err error) {
	result, response, err = dnsRecords.UpdateDnsRecordWithContext(context.Background(), updateDnsRecordOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateDnsRecordWithContext is an alternate form of the UpdateDnsRecord method which supports a Context parameter
func (dnsRecords *DnsRecordsV1) UpdateDnsRecordWithContext(ctx context.Context, updateDnsRecordOptions *UpdateDnsRecordOptions) (result *DnsrecordResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDnsRecordOptions, "updateDnsRecordOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateDnsRecordOptions, "updateDnsRecordOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":                  *dnsRecords.Crn,
		"zone_identifier":      *dnsRecords.ZoneIdentifier,
		"dnsrecord_identifier": *updateDnsRecordOptions.DnsrecordIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsRecords.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsRecords.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/dns_records/{dnsrecord_identifier}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	sdkHeaders := common.GetSdkHeaders("dns_records", "V1", "UpdateDnsRecord")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	for headerName, headerValue := range updateDnsRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateDnsRecordOptions.Name != nil {
		body["name"] = updateDnsRecordOptions.Name
	}
	if updateDnsRecordOptions.Type != nil {
		body["type"] = updateDnsRecordOptions.Type
	}
	if updateDnsRecordOptions.TTL != nil {
		body["ttl"] = updateDnsRecordOptions.TTL
	}
	if updateDnsRecordOptions.Content != nil {
		body["content"] = updateDnsRecordOptions.Content
	}
	if updateDnsRecordOptions.Priority != nil {
		body["priority"] = updateDnsRecordOptions.Priority
	}
	if updateDnsRecordOptions.Proxied != nil {
		body["proxied"] = updateDnsRecordOptions.Proxied
	}
	if updateDnsRecordOptions.Data != nil {
		body["data"] = updateDnsRecordOptions.Data
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
	response, err = dnsRecords.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_dns_record", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnsrecordResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// BatchDnsRecords : Batch DNS records
// Send a Batch of DNS Record API calls to be executed together. The operations you specify within the /batch request
// body are always executed in the following order: deletes, patches, puts, posts.
func (dnsRecords *DnsRecordsV1) BatchDnsRecords(batchDnsRecordsOptions *BatchDnsRecordsOptions) (result *BatchDnsRecordsResponse, response *core.DetailedResponse, err error) {
	result, response, err = dnsRecords.BatchDnsRecordsWithContext(context.Background(), batchDnsRecordsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// BatchDnsRecordsWithContext is an alternate form of the BatchDnsRecords method which supports a Context parameter
func (dnsRecords *DnsRecordsV1) BatchDnsRecordsWithContext(ctx context.Context, batchDnsRecordsOptions *BatchDnsRecordsOptions) (result *BatchDnsRecordsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(batchDnsRecordsOptions, "batchDnsRecordsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *dnsRecords.Crn,
		"zone_identifier": *dnsRecords.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsRecords.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsRecords.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/dns_records/batch`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	sdkHeaders := common.GetSdkHeaders("dns_records", "V1", "BatchDnsRecords")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	for headerName, headerValue := range batchDnsRecordsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if batchDnsRecordsOptions.Deletes != nil {
		body["deletes"] = batchDnsRecordsOptions.Deletes
	}
	if batchDnsRecordsOptions.Patches != nil {
		body["patches"] = batchDnsRecordsOptions.Patches
	}
	if batchDnsRecordsOptions.Posts != nil {
		body["posts"] = batchDnsRecordsOptions.Posts
	}
	if batchDnsRecordsOptions.Puts != nil {
		body["puts"] = batchDnsRecordsOptions.Puts
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
	response, err = dnsRecords.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "batch_dns_records", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBatchDnsRecordsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.1")
}

// BatchDnsRecordsOptions : The BatchDnsRecords options.
type BatchDnsRecordsOptions struct {
	Deletes []BatchDnsRecordsRequestDeletesItem `json:"deletes,omitempty"`

	Patches []BatchDnsRecordsRequestPatchesItem `json:"patches,omitempty"`

	Posts []DnsrecordInput `json:"posts,omitempty"`

	Puts []BatchDnsRecordsRequestPutsItem `json:"puts,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewBatchDnsRecordsOptions : Instantiate BatchDnsRecordsOptions
func (*DnsRecordsV1) NewBatchDnsRecordsOptions() *BatchDnsRecordsOptions {
	return &BatchDnsRecordsOptions{}
}

// SetDeletes : Allow user to set Deletes
func (_options *BatchDnsRecordsOptions) SetDeletes(deletes []BatchDnsRecordsRequestDeletesItem) *BatchDnsRecordsOptions {
	_options.Deletes = deletes
	return _options
}

// SetPatches : Allow user to set Patches
func (_options *BatchDnsRecordsOptions) SetPatches(patches []BatchDnsRecordsRequestPatchesItem) *BatchDnsRecordsOptions {
	_options.Patches = patches
	return _options
}

// SetPosts : Allow user to set Posts
func (_options *BatchDnsRecordsOptions) SetPosts(posts []DnsrecordInput) *BatchDnsRecordsOptions {
	_options.Posts = posts
	return _options
}

// SetPuts : Allow user to set Puts
func (_options *BatchDnsRecordsOptions) SetPuts(puts []BatchDnsRecordsRequestPutsItem) *BatchDnsRecordsOptions {
	_options.Puts = puts
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *BatchDnsRecordsOptions) SetHeaders(param map[string]string) *BatchDnsRecordsOptions {
	options.Headers = param
	return options
}

// BatchDnsRecordsRequestDeletesItem : BatchDnsRecordsRequestDeletesItem struct
type BatchDnsRecordsRequestDeletesItem struct {
	// DNS record ID to delete.
	ID *string `json:"id" validate:"required"`
}

// NewBatchDnsRecordsRequestDeletesItem : Instantiate BatchDnsRecordsRequestDeletesItem (Generic Model Constructor)
func (*DnsRecordsV1) NewBatchDnsRecordsRequestDeletesItem(id string) (_model *BatchDnsRecordsRequestDeletesItem, err error) {
	_model = &BatchDnsRecordsRequestDeletesItem{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalBatchDnsRecordsRequestDeletesItem unmarshals an instance of BatchDnsRecordsRequestDeletesItem from the specified map of raw messages.
func UnmarshalBatchDnsRecordsRequestDeletesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BatchDnsRecordsRequestDeletesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BatchDnsRecordsRequestPatchesItem : BatchDnsRecordsRequestPatchesItem struct
type BatchDnsRecordsRequestPatchesItem struct {
	// DNS record ID to patch.
	ID *string `json:"id" validate:"required"`

	// Required for all record types except SRV.
	Name *string `json:"name,omitempty"`

	// dns record type.
	Type *string `json:"type,omitempty"`

	// dns record ttl value.
	TTL *int64 `json:"ttl,omitempty"`

	// content of dns record.
	Content *string `json:"content,omitempty"`

	// For MX records only.
	Priority *int64 `json:"priority,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// For LOC, SRV, CAA, DS records only.
	Data map[string]interface{} `json:"data,omitempty"`
}

// Constants associated with the BatchDnsRecordsRequestPatchesItem.Type property.
// dns record type.
const (
	BatchDnsRecordsRequestPatchesItem_Type_A     = "A"
	BatchDnsRecordsRequestPatchesItem_Type_Aaaa  = "AAAA"
	BatchDnsRecordsRequestPatchesItem_Type_Caa   = "CAA"
	BatchDnsRecordsRequestPatchesItem_Type_Cname = "CNAME"
	BatchDnsRecordsRequestPatchesItem_Type_Ds    = "DS"
	BatchDnsRecordsRequestPatchesItem_Type_Loc   = "LOC"
	BatchDnsRecordsRequestPatchesItem_Type_Mx    = "MX"
	BatchDnsRecordsRequestPatchesItem_Type_Ns    = "NS"
	BatchDnsRecordsRequestPatchesItem_Type_Ptr   = "PTR"
	BatchDnsRecordsRequestPatchesItem_Type_Srv   = "SRV"
	BatchDnsRecordsRequestPatchesItem_Type_Txt   = "TXT"
)

// NewBatchDnsRecordsRequestPatchesItem : Instantiate BatchDnsRecordsRequestPatchesItem (Generic Model Constructor)
func (*DnsRecordsV1) NewBatchDnsRecordsRequestPatchesItem(id string) (_model *BatchDnsRecordsRequestPatchesItem, err error) {
	_model = &BatchDnsRecordsRequestPatchesItem{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalBatchDnsRecordsRequestPatchesItem unmarshals an instance of BatchDnsRecordsRequestPatchesItem from the specified map of raw messages.
func UnmarshalBatchDnsRecordsRequestPatchesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BatchDnsRecordsRequestPatchesItem)
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
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		err = core.SDKErrorf(err, "", "content-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "proxied", &obj.Proxied)
	if err != nil {
		err = core.SDKErrorf(err, "", "proxied-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BatchDnsRecordsRequestPutsItem : BatchDnsRecordsRequestPutsItem struct
type BatchDnsRecordsRequestPutsItem struct {
	// DNS record ID to update.
	ID *string `json:"id" validate:"required"`

	// Required for all record types except SRV.
	Name *string `json:"name" validate:"required"`

	// dns record type.
	Type *string `json:"type" validate:"required"`

	// dns record ttl value.
	TTL *int64 `json:"ttl" validate:"required"`

	// dns record content.
	Content *string `json:"content" validate:"required"`

	// For MX records only.
	Priority *int64 `json:"priority,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// For LOC, SRV, CAA, DS records only.
	Data map[string]interface{} `json:"data,omitempty"`
}

// Constants associated with the BatchDnsRecordsRequestPutsItem.Type property.
// dns record type.
const (
	BatchDnsRecordsRequestPutsItem_Type_A     = "A"
	BatchDnsRecordsRequestPutsItem_Type_Aaaa  = "AAAA"
	BatchDnsRecordsRequestPutsItem_Type_Caa   = "CAA"
	BatchDnsRecordsRequestPutsItem_Type_Cname = "CNAME"
	BatchDnsRecordsRequestPutsItem_Type_Ds    = "DS"
	BatchDnsRecordsRequestPutsItem_Type_Loc   = "LOC"
	BatchDnsRecordsRequestPutsItem_Type_Mx    = "MX"
	BatchDnsRecordsRequestPutsItem_Type_Ns    = "NS"
	BatchDnsRecordsRequestPutsItem_Type_Ptr   = "PTR"
	BatchDnsRecordsRequestPutsItem_Type_Srv   = "SRV"
	BatchDnsRecordsRequestPutsItem_Type_Txt   = "TXT"
)

// NewBatchDnsRecordsRequestPutsItem : Instantiate BatchDnsRecordsRequestPutsItem (Generic Model Constructor)
func (*DnsRecordsV1) NewBatchDnsRecordsRequestPutsItem(id string, name string, typeVar string, ttl int64, content string) (_model *BatchDnsRecordsRequestPutsItem, err error) {
	_model = &BatchDnsRecordsRequestPutsItem{
		ID:      core.StringPtr(id),
		Name:    core.StringPtr(name),
		Type:    core.StringPtr(typeVar),
		TTL:     core.Int64Ptr(ttl),
		Content: core.StringPtr(content),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalBatchDnsRecordsRequestPutsItem unmarshals an instance of BatchDnsRecordsRequestPutsItem from the specified map of raw messages.
func UnmarshalBatchDnsRecordsRequestPutsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BatchDnsRecordsRequestPutsItem)
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
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		err = core.SDKErrorf(err, "", "content-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "proxied", &obj.Proxied)
	if err != nil {
		err = core.SDKErrorf(err, "", "proxied-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BatchDnsRecordsResponseResult : BatchDnsRecordsResponseResult struct
type BatchDnsRecordsResponseResult struct {
	Deletes []BatchDnsRecordDetails `json:"deletes,omitempty"`

	Patches []BatchDnsRecordDetails `json:"patches,omitempty"`

	Posts []BatchDnsRecordDetails `json:"posts,omitempty"`

	Puts []BatchDnsRecordDetails `json:"puts,omitempty"`
}

// UnmarshalBatchDnsRecordsResponseResult unmarshals an instance of BatchDnsRecordsResponseResult from the specified map of raw messages.
func UnmarshalBatchDnsRecordsResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BatchDnsRecordsResponseResult)
	err = core.UnmarshalModel(m, "deletes", &obj.Deletes, UnmarshalBatchDnsRecordDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "deletes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "patches", &obj.Patches, UnmarshalBatchDnsRecordDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "patches-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "posts", &obj.Posts, UnmarshalBatchDnsRecordDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "posts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "puts", &obj.Puts, UnmarshalBatchDnsRecordDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "puts-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateDnsRecordOptions : The CreateDnsRecord options.
type CreateDnsRecordOptions struct {
	// Required for all record types except SRV.
	Name *string `json:"name,omitempty"`

	// dns record type.
	Type *string `json:"type,omitempty"`

	// dns record ttl value.
	TTL *int64 `json:"ttl,omitempty"`

	// dns record content.
	Content *string `json:"content,omitempty"`

	// For MX records only.
	Priority *int64 `json:"priority,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// For LOC, SRV, CAA, DS records only.
	Data map[string]interface{} `json:"data,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateDnsRecordOptions.Type property.
// dns record type.
const (
	CreateDnsRecordOptions_Type_A     = "A"
	CreateDnsRecordOptions_Type_Aaaa  = "AAAA"
	CreateDnsRecordOptions_Type_Caa   = "CAA"
	CreateDnsRecordOptions_Type_Cname = "CNAME"
	CreateDnsRecordOptions_Type_Ds    = "DS"
	CreateDnsRecordOptions_Type_Loc   = "LOC"
	CreateDnsRecordOptions_Type_Mx    = "MX"
	CreateDnsRecordOptions_Type_Ns    = "NS"
	CreateDnsRecordOptions_Type_Ptr   = "PTR"
	CreateDnsRecordOptions_Type_Srv   = "SRV"
	CreateDnsRecordOptions_Type_Txt   = "TXT"
)

// NewCreateDnsRecordOptions : Instantiate CreateDnsRecordOptions
func (*DnsRecordsV1) NewCreateDnsRecordOptions() *CreateDnsRecordOptions {
	return &CreateDnsRecordOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateDnsRecordOptions) SetName(name string) *CreateDnsRecordOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateDnsRecordOptions) SetType(typeVar string) *CreateDnsRecordOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetTTL : Allow user to set TTL
func (_options *CreateDnsRecordOptions) SetTTL(ttl int64) *CreateDnsRecordOptions {
	_options.TTL = core.Int64Ptr(ttl)
	return _options
}

// SetContent : Allow user to set Content
func (_options *CreateDnsRecordOptions) SetContent(content string) *CreateDnsRecordOptions {
	_options.Content = core.StringPtr(content)
	return _options
}

// SetPriority : Allow user to set Priority
func (_options *CreateDnsRecordOptions) SetPriority(priority int64) *CreateDnsRecordOptions {
	_options.Priority = core.Int64Ptr(priority)
	return _options
}

// SetProxied : Allow user to set Proxied
func (_options *CreateDnsRecordOptions) SetProxied(proxied bool) *CreateDnsRecordOptions {
	_options.Proxied = core.BoolPtr(proxied)
	return _options
}

// SetData : Allow user to set Data
func (_options *CreateDnsRecordOptions) SetData(data map[string]interface{}) *CreateDnsRecordOptions {
	_options.Data = data
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDnsRecordOptions) SetHeaders(param map[string]string) *CreateDnsRecordOptions {
	options.Headers = param
	return options
}

// DeleteDnsRecordOptions : The DeleteDnsRecord options.
type DeleteDnsRecordOptions struct {
	// Identifier of DNS record.
	DnsrecordIdentifier *string `json:"dnsrecord_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteDnsRecordOptions : Instantiate DeleteDnsRecordOptions
func (*DnsRecordsV1) NewDeleteDnsRecordOptions(dnsrecordIdentifier string) *DeleteDnsRecordOptions {
	return &DeleteDnsRecordOptions{
		DnsrecordIdentifier: core.StringPtr(dnsrecordIdentifier),
	}
}

// SetDnsrecordIdentifier : Allow user to set DnsrecordIdentifier
func (_options *DeleteDnsRecordOptions) SetDnsrecordIdentifier(dnsrecordIdentifier string) *DeleteDnsRecordOptions {
	_options.DnsrecordIdentifier = core.StringPtr(dnsrecordIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDnsRecordOptions) SetHeaders(param map[string]string) *DeleteDnsRecordOptions {
	options.Headers = param
	return options
}

// DeleteDnsrecordRespResult : result.
type DeleteDnsrecordRespResult struct {
	// dns record id.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalDeleteDnsrecordRespResult unmarshals an instance of DeleteDnsrecordRespResult from the specified map of raw messages.
func UnmarshalDeleteDnsrecordRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteDnsrecordRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetDnsRecordOptions : The GetDnsRecord options.
type GetDnsRecordOptions struct {
	// Identifier of DNS record.
	DnsrecordIdentifier *string `json:"dnsrecord_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDnsRecordOptions : Instantiate GetDnsRecordOptions
func (*DnsRecordsV1) NewGetDnsRecordOptions(dnsrecordIdentifier string) *GetDnsRecordOptions {
	return &GetDnsRecordOptions{
		DnsrecordIdentifier: core.StringPtr(dnsrecordIdentifier),
	}
}

// SetDnsrecordIdentifier : Allow user to set DnsrecordIdentifier
func (_options *GetDnsRecordOptions) SetDnsrecordIdentifier(dnsrecordIdentifier string) *GetDnsRecordOptions {
	_options.DnsrecordIdentifier = core.StringPtr(dnsrecordIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDnsRecordOptions) SetHeaders(param map[string]string) *GetDnsRecordOptions {
	options.Headers = param
	return options
}

// ListAllDnsRecordsOptions : The ListAllDnsRecords options.
type ListAllDnsRecordsOptions struct {
	// Type of DNS records to display.
	Type *string `json:"type,omitempty"`

	// Value of name field to filter by.
	Name *string `json:"name,omitempty"`

	// Value of content field to filter by.
	Content *string `json:"content,omitempty"`

	// Page number of paginated results.
	Page *int64 `json:"page,omitempty"`

	// Maximum number of DNS records per page.
	PerPage *int64 `json:"per_page,omitempty"`

	// Field by which to order list of DNS records.
	Order *string `json:"order,omitempty"`

	// Direction in which to order results [ascending/descending order].
	Direction *string `json:"direction,omitempty"`

	// Whether to match all (all) or atleast one search parameter (any).
	Match *string `json:"match,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListAllDnsRecordsOptions.Order property.
// Field by which to order list of DNS records.
const (
	ListAllDnsRecordsOptions_Order_Content = "content"
	ListAllDnsRecordsOptions_Order_Name    = "name"
	ListAllDnsRecordsOptions_Order_Proxied = "proxied"
	ListAllDnsRecordsOptions_Order_TTL     = "ttl"
	ListAllDnsRecordsOptions_Order_Type    = "type"
)

// Constants associated with the ListAllDnsRecordsOptions.Direction property.
// Direction in which to order results [ascending/descending order].
const (
	ListAllDnsRecordsOptions_Direction_Asc  = "asc"
	ListAllDnsRecordsOptions_Direction_Desc = "desc"
)

// Constants associated with the ListAllDnsRecordsOptions.Match property.
// Whether to match all (all) or atleast one search parameter (any).
const (
	ListAllDnsRecordsOptions_Match_All = "all"
	ListAllDnsRecordsOptions_Match_Any = "any"
)

// NewListAllDnsRecordsOptions : Instantiate ListAllDnsRecordsOptions
func (*DnsRecordsV1) NewListAllDnsRecordsOptions() *ListAllDnsRecordsOptions {
	return &ListAllDnsRecordsOptions{}
}

// SetType : Allow user to set Type
func (_options *ListAllDnsRecordsOptions) SetType(typeVar string) *ListAllDnsRecordsOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListAllDnsRecordsOptions) SetName(name string) *ListAllDnsRecordsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetContent : Allow user to set Content
func (_options *ListAllDnsRecordsOptions) SetContent(content string) *ListAllDnsRecordsOptions {
	_options.Content = core.StringPtr(content)
	return _options
}

// SetPage : Allow user to set Page
func (_options *ListAllDnsRecordsOptions) SetPage(page int64) *ListAllDnsRecordsOptions {
	_options.Page = core.Int64Ptr(page)
	return _options
}

// SetPerPage : Allow user to set PerPage
func (_options *ListAllDnsRecordsOptions) SetPerPage(perPage int64) *ListAllDnsRecordsOptions {
	_options.PerPage = core.Int64Ptr(perPage)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListAllDnsRecordsOptions) SetOrder(order string) *ListAllDnsRecordsOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetDirection : Allow user to set Direction
func (_options *ListAllDnsRecordsOptions) SetDirection(direction string) *ListAllDnsRecordsOptions {
	_options.Direction = core.StringPtr(direction)
	return _options
}

// SetMatch : Allow user to set Match
func (_options *ListAllDnsRecordsOptions) SetMatch(match string) *ListAllDnsRecordsOptions {
	_options.Match = core.StringPtr(match)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllDnsRecordsOptions) SetHeaders(param map[string]string) *ListAllDnsRecordsOptions {
	options.Headers = param
	return options
}

// UpdateDnsRecordOptions : The UpdateDnsRecord options.
type UpdateDnsRecordOptions struct {
	// Identifier of DNS record.
	DnsrecordIdentifier *string `json:"dnsrecord_identifier" validate:"required,ne="`

	// Required for all record types except SRV.
	Name *string `json:"name,omitempty"`

	// dns record type.
	Type *string `json:"type,omitempty"`

	// dns record ttl value.
	TTL *int64 `json:"ttl,omitempty"`

	// content of dns record.
	Content *string `json:"content,omitempty"`

	// For MX records only.
	Priority *int64 `json:"priority,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// For LOC, SRV, CAA, DS records only.
	Data map[string]interface{} `json:"data,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the UpdateDnsRecordOptions.Type property.
// dns record type.
const (
	UpdateDnsRecordOptions_Type_A     = "A"
	UpdateDnsRecordOptions_Type_Aaaa  = "AAAA"
	UpdateDnsRecordOptions_Type_Caa   = "CAA"
	UpdateDnsRecordOptions_Type_Cname = "CNAME"
	UpdateDnsRecordOptions_Type_Ds    = "DS"
	UpdateDnsRecordOptions_Type_Loc   = "LOC"
	UpdateDnsRecordOptions_Type_Mx    = "MX"
	UpdateDnsRecordOptions_Type_Ns    = "NS"
	UpdateDnsRecordOptions_Type_Ptr   = "PTR"
	UpdateDnsRecordOptions_Type_Srv   = "SRV"
	UpdateDnsRecordOptions_Type_Txt   = "TXT"
)

// NewUpdateDnsRecordOptions : Instantiate UpdateDnsRecordOptions
func (*DnsRecordsV1) NewUpdateDnsRecordOptions(dnsrecordIdentifier string) *UpdateDnsRecordOptions {
	return &UpdateDnsRecordOptions{
		DnsrecordIdentifier: core.StringPtr(dnsrecordIdentifier),
	}
}

// SetDnsrecordIdentifier : Allow user to set DnsrecordIdentifier
func (_options *UpdateDnsRecordOptions) SetDnsrecordIdentifier(dnsrecordIdentifier string) *UpdateDnsRecordOptions {
	_options.DnsrecordIdentifier = core.StringPtr(dnsrecordIdentifier)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateDnsRecordOptions) SetName(name string) *UpdateDnsRecordOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *UpdateDnsRecordOptions) SetType(typeVar string) *UpdateDnsRecordOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetTTL : Allow user to set TTL
func (_options *UpdateDnsRecordOptions) SetTTL(ttl int64) *UpdateDnsRecordOptions {
	_options.TTL = core.Int64Ptr(ttl)
	return _options
}

// SetContent : Allow user to set Content
func (_options *UpdateDnsRecordOptions) SetContent(content string) *UpdateDnsRecordOptions {
	_options.Content = core.StringPtr(content)
	return _options
}

// SetPriority : Allow user to set Priority
func (_options *UpdateDnsRecordOptions) SetPriority(priority int64) *UpdateDnsRecordOptions {
	_options.Priority = core.Int64Ptr(priority)
	return _options
}

// SetProxied : Allow user to set Proxied
func (_options *UpdateDnsRecordOptions) SetProxied(proxied bool) *UpdateDnsRecordOptions {
	_options.Proxied = core.BoolPtr(proxied)
	return _options
}

// SetData : Allow user to set Data
func (_options *UpdateDnsRecordOptions) SetData(data map[string]interface{}) *UpdateDnsRecordOptions {
	_options.Data = data
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDnsRecordOptions) SetHeaders(param map[string]string) *UpdateDnsRecordOptions {
	options.Headers = param
	return options
}

// BatchDnsRecordDetails : dns record details as returned by the batch API.
type BatchDnsRecordDetails struct {
	// dns record identifier.
	ID *string `json:"id,omitempty"`

	// created on.
	CreatedOn *string `json:"created_on,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`

	// dns record name.
	Name *string `json:"name,omitempty"`

	// dns record type.
	Type *string `json:"type,omitempty"`

	// dns record content.
	Content *string `json:"content,omitempty"`

	// proxiable.
	Proxiable *bool `json:"proxiable,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// dns record ttl value.
	TTL *int64 `json:"ttl,omitempty"`

	// Relevant only to MX type records.
	Priority *int64 `json:"priority,omitempty"`

	// Data details for the DNS record. Only for LOC, SRV, CAA records.
	Data map[string]interface{} `json:"data,omitempty"`

	// DNS record settings.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// DNS record metadata.
	Meta map[string]interface{} `json:"meta,omitempty"`

	// Optional comment for the DNS record.
	Comment *string `json:"comment,omitempty"`

	// Tags associated with the DNS record.
	Tags []string `json:"tags,omitempty"`
}

// Constants associated with the BatchDnsRecordDetails.Type property.
// dns record type.
const (
	BatchDnsRecordDetails_Type_A     = "A"
	BatchDnsRecordDetails_Type_Aaaa  = "AAAA"
	BatchDnsRecordDetails_Type_Caa   = "CAA"
	BatchDnsRecordDetails_Type_Cname = "CNAME"
	BatchDnsRecordDetails_Type_Ds    = "DS"
	BatchDnsRecordDetails_Type_Loc   = "LOC"
	BatchDnsRecordDetails_Type_Mx    = "MX"
	BatchDnsRecordDetails_Type_Ns    = "NS"
	BatchDnsRecordDetails_Type_Ptr   = "PTR"
	BatchDnsRecordDetails_Type_Srv   = "SRV"
	BatchDnsRecordDetails_Type_Txt   = "TXT"
)

// UnmarshalBatchDnsRecordDetails unmarshals an instance of BatchDnsRecordDetails from the specified map of raw messages.
func UnmarshalBatchDnsRecordDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BatchDnsRecordDetails)
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
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		err = core.SDKErrorf(err, "", "content-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "proxiable", &obj.Proxiable)
	if err != nil {
		err = core.SDKErrorf(err, "", "proxiable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "proxied", &obj.Proxied)
	if err != nil {
		err = core.SDKErrorf(err, "", "proxied-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		err = core.SDKErrorf(err, "", "ttl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "meta", &obj.Meta)
	if err != nil {
		err = core.SDKErrorf(err, "", "meta-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "comment", &obj.Comment)
	if err != nil {
		err = core.SDKErrorf(err, "", "comment-error", common.GetComponentInfo())
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

// BatchDnsRecordsResponse : Batch DNS records response.
type BatchDnsRecordsResponse struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	Result *BatchDnsRecordsResponseResult `json:"result" validate:"required"`
}

// UnmarshalBatchDnsRecordsResponse unmarshals an instance of BatchDnsRecordsResponse from the specified map of raw messages.
func UnmarshalBatchDnsRecordsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BatchDnsRecordsResponse)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalBatchDnsRecordsResponseResult)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteDnsrecordResp : dns record delete response.
type DeleteDnsrecordResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result *DeleteDnsrecordRespResult `json:"result" validate:"required"`
}

// UnmarshalDeleteDnsrecordResp unmarshals an instance of DeleteDnsrecordResp from the specified map of raw messages.
func UnmarshalDeleteDnsrecordResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteDnsrecordResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteDnsrecordRespResult)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DnsrecordDetails : dns record details.
type DnsrecordDetails struct {
	// dns record identifier.
	ID *string `json:"id,omitempty"`

	// created on.
	CreatedOn *string `json:"created_on,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`

	// dns record name.
	Name *string `json:"name,omitempty"`

	// dns record type.
	Type *string `json:"type,omitempty"`

	// dns record content.
	Content *string `json:"content,omitempty"`

	// zone identifier.
	ZoneID *string `json:"zone_id,omitempty"`

	// zone name.
	ZoneName *string `json:"zone_name,omitempty"`

	// proxiable.
	Proxiable *bool `json:"proxiable,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// dns record ttl value.
	TTL *int64 `json:"ttl,omitempty"`

	// Relevant only to MX type records.
	Priority *int64 `json:"priority,omitempty"`

	// Data details for the DNS record. Only for LOC, SRV, CAA records.
	Data map[string]interface{} `json:"data,omitempty"`
}

// Constants associated with the DnsrecordDetails.Type property.
// dns record type.
const (
	DnsrecordDetails_Type_A     = "A"
	DnsrecordDetails_Type_Aaaa  = "AAAA"
	DnsrecordDetails_Type_Caa   = "CAA"
	DnsrecordDetails_Type_Cname = "CNAME"
	DnsrecordDetails_Type_Ds    = "DS"
	DnsrecordDetails_Type_Loc   = "LOC"
	DnsrecordDetails_Type_Mx    = "MX"
	DnsrecordDetails_Type_Ns    = "NS"
	DnsrecordDetails_Type_Ptr   = "PTR"
	DnsrecordDetails_Type_Srv   = "SRV"
	DnsrecordDetails_Type_Txt   = "TXT"
)

// UnmarshalDnsrecordDetails unmarshals an instance of DnsrecordDetails from the specified map of raw messages.
func UnmarshalDnsrecordDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DnsrecordDetails)
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
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		err = core.SDKErrorf(err, "", "content-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "proxiable", &obj.Proxiable)
	if err != nil {
		err = core.SDKErrorf(err, "", "proxiable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "proxied", &obj.Proxied)
	if err != nil {
		err = core.SDKErrorf(err, "", "proxied-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		err = core.SDKErrorf(err, "", "ttl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DnsrecordInput : dns record input.
type DnsrecordInput struct {
	// Required for all record types except SRV.
	Name *string `json:"name,omitempty"`

	// dns record type.
	Type *string `json:"type" validate:"required"`

	// dns record ttl value.
	TTL *int64 `json:"ttl,omitempty"`

	// dns record content.
	Content *string `json:"content,omitempty"`

	// For MX records only.
	Priority *int64 `json:"priority,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// For LOC, SRV, CAA, DS records only.
	Data map[string]interface{} `json:"data,omitempty"`
}

// Constants associated with the DnsrecordInput.Type property.
// dns record type.
const (
	DnsrecordInput_Type_A     = "A"
	DnsrecordInput_Type_Aaaa  = "AAAA"
	DnsrecordInput_Type_Caa   = "CAA"
	DnsrecordInput_Type_Cname = "CNAME"
	DnsrecordInput_Type_Ds    = "DS"
	DnsrecordInput_Type_Loc   = "LOC"
	DnsrecordInput_Type_Mx    = "MX"
	DnsrecordInput_Type_Ns    = "NS"
	DnsrecordInput_Type_Ptr   = "PTR"
	DnsrecordInput_Type_Srv   = "SRV"
	DnsrecordInput_Type_Txt   = "TXT"
)

// NewDnsrecordInput : Instantiate DnsrecordInput (Generic Model Constructor)
func (*DnsRecordsV1) NewDnsrecordInput(typeVar string) (_model *DnsrecordInput, err error) {
	_model = &DnsrecordInput{
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalDnsrecordInput unmarshals an instance of DnsrecordInput from the specified map of raw messages.
func UnmarshalDnsrecordInput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DnsrecordInput)
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
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		err = core.SDKErrorf(err, "", "content-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		err = core.SDKErrorf(err, "", "priority-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "proxied", &obj.Proxied)
	if err != nil {
		err = core.SDKErrorf(err, "", "proxied-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DnsrecordResp : dns record response.
type DnsrecordResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// dns record details.
	Result *DnsrecordDetails `json:"result" validate:"required"`
}

// UnmarshalDnsrecordResp unmarshals an instance of DnsrecordResp from the specified map of raw messages.
func UnmarshalDnsrecordResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DnsrecordResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDnsrecordDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListDnsrecordsResp : dns records list response.
type ListDnsrecordsResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// dns record list.
	Result []DnsrecordDetails `json:"result" validate:"required"`

	// result information.
	ResultInfo *ResultInfo `json:"result_info" validate:"required"`
}

// UnmarshalListDnsrecordsResp unmarshals an instance of ListDnsrecordsResp from the specified map of raw messages.
func UnmarshalListDnsrecordsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListDnsrecordsResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDnsrecordDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalResultInfo)
	if err != nil {
		err = core.SDKErrorf(err, "", "result_info-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResultInfo : result information.
type ResultInfo struct {
	// page.
	Page *int64 `json:"page" validate:"required"`

	// per page.
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
		err = core.SDKErrorf(err, "", "page-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "per_page", &obj.PerPage)
	if err != nil {
		err = core.SDKErrorf(err, "", "per_page-error", common.GetComponentInfo())
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
