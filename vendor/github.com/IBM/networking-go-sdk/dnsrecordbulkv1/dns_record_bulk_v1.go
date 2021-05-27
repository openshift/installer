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
 

// Package dnsrecordbulkv1 : Operations and models for the DnsRecordBulkV1 service
package dnsrecordbulkv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/networking-go-sdk/common"
	"io"
	"net/http"
	"reflect"
	"time"
)

// DnsRecordBulkV1 : Import/Export zone files
//
// Version: 1.0.0
type DnsRecordBulkV1 struct {
	Service *core.BaseService

	// Full url-encoded CRN of the service instance.
	Crn *string

	// Identifier of zone.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "dns_record_bulk"

// DnsRecordBulkV1Options : Service options
type DnsRecordBulkV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full url-encoded CRN of the service instance.
	Crn *string `validate:"required"`

	// Identifier of zone.
	ZoneIdentifier *string `validate:"required"`
}

// NewDnsRecordBulkV1UsingExternalConfig : constructs an instance of DnsRecordBulkV1 with passed in options and external configuration.
func NewDnsRecordBulkV1UsingExternalConfig(options *DnsRecordBulkV1Options) (dnsRecordBulk *DnsRecordBulkV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	dnsRecordBulk, err = NewDnsRecordBulkV1(options)
	if err != nil {
		return
	}

	err = dnsRecordBulk.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = dnsRecordBulk.Service.SetServiceURL(options.URL)
	}
	return
}

// NewDnsRecordBulkV1 : constructs an instance of DnsRecordBulkV1 with passed in options.
func NewDnsRecordBulkV1(options *DnsRecordBulkV1Options) (service *DnsRecordBulkV1, err error) {
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

	service = &DnsRecordBulkV1{
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

// Clone makes a copy of "dnsRecordBulk" suitable for processing requests.
func (dnsRecordBulk *DnsRecordBulkV1) Clone() *DnsRecordBulkV1 {
	if core.IsNil(dnsRecordBulk) {
		return nil
	}
	clone := *dnsRecordBulk
	clone.Service = dnsRecordBulk.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (dnsRecordBulk *DnsRecordBulkV1) SetServiceURL(url string) error {
	return dnsRecordBulk.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (dnsRecordBulk *DnsRecordBulkV1) GetServiceURL() string {
	return dnsRecordBulk.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (dnsRecordBulk *DnsRecordBulkV1) SetDefaultHeaders(headers http.Header) {
	dnsRecordBulk.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (dnsRecordBulk *DnsRecordBulkV1) SetEnableGzipCompression(enableGzip bool) {
	dnsRecordBulk.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (dnsRecordBulk *DnsRecordBulkV1) GetEnableGzipCompression() bool {
	return dnsRecordBulk.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (dnsRecordBulk *DnsRecordBulkV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	dnsRecordBulk.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (dnsRecordBulk *DnsRecordBulkV1) DisableRetries() {
	dnsRecordBulk.Service.DisableRetries()
}

// GetDnsRecordsBulk : Export zone file
// Export zone file.
func (dnsRecordBulk *DnsRecordBulkV1) GetDnsRecordsBulk(getDnsRecordsBulkOptions *GetDnsRecordsBulkOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return dnsRecordBulk.GetDnsRecordsBulkWithContext(context.Background(), getDnsRecordsBulkOptions)
}

// GetDnsRecordsBulkWithContext is an alternate form of the GetDnsRecordsBulk method which supports a Context parameter
func (dnsRecordBulk *DnsRecordBulkV1) GetDnsRecordsBulkWithContext(ctx context.Context, getDnsRecordsBulkOptions *GetDnsRecordsBulkOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getDnsRecordsBulkOptions, "getDnsRecordsBulkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *dnsRecordBulk.Crn,
		"zone_identifier": *dnsRecordBulk.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsRecordBulk.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsRecordBulk.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/dns_records_bulk`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDnsRecordsBulkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_record_bulk", "V1", "GetDnsRecordsBulk")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "text/plain; charset=utf-8")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = dnsRecordBulk.Service.Request(request, &result)

	return
}

// PostDnsRecordsBulk : Import zone file
// Import zone file.
func (dnsRecordBulk *DnsRecordBulkV1) PostDnsRecordsBulk(postDnsRecordsBulkOptions *PostDnsRecordsBulkOptions) (result *DnsRecordsObject, response *core.DetailedResponse, err error) {
	return dnsRecordBulk.PostDnsRecordsBulkWithContext(context.Background(), postDnsRecordsBulkOptions)
}

// PostDnsRecordsBulkWithContext is an alternate form of the PostDnsRecordsBulk method which supports a Context parameter
func (dnsRecordBulk *DnsRecordBulkV1) PostDnsRecordsBulkWithContext(ctx context.Context, postDnsRecordsBulkOptions *PostDnsRecordsBulkOptions) (result *DnsRecordsObject, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDnsRecordsBulkOptions, "postDnsRecordsBulkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postDnsRecordsBulkOptions, "postDnsRecordsBulkOptions")
	if err != nil {
		return
	}
	if (postDnsRecordsBulkOptions.File == nil) {
		err = fmt.Errorf("at least one of  or file must be supplied")
		return
	}

	pathParamsMap := map[string]string{
		"crn": *dnsRecordBulk.Crn,
		"zone_identifier": *dnsRecordBulk.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = dnsRecordBulk.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(dnsRecordBulk.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/dns_records_bulk`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postDnsRecordsBulkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_record_bulk", "V1", "PostDnsRecordsBulk")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if postDnsRecordsBulkOptions.File != nil {
		builder.AddFormData("file", "filename",
			core.StringNilMapper(postDnsRecordsBulkOptions.FileContentType), postDnsRecordsBulkOptions.File)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsRecordBulk.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDnsRecordsObject)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DnsRecordsObjectMessagesItem : DnsRecordsObjectMessagesItem struct
type DnsRecordsObjectMessagesItem struct {
	// Message code.
	Code *int64 `json:"code,omitempty"`

	// Message corresponding to the code.
	Message *string `json:"message,omitempty"`
}


// UnmarshalDnsRecordsObjectMessagesItem unmarshals an instance of DnsRecordsObjectMessagesItem from the specified map of raw messages.
func UnmarshalDnsRecordsObjectMessagesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DnsRecordsObjectMessagesItem)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DnsRecordsObjectResult : DNS record.
type DnsRecordsObjectResult struct {
	// total records added.
	RecsAdded *int64 `json:"recs_added" validate:"required"`

	// total records parsed.
	TotalRecordsParsed *int64 `json:"total_records_parsed" validate:"required"`
}


// UnmarshalDnsRecordsObjectResult unmarshals an instance of DnsRecordsObjectResult from the specified map of raw messages.
func UnmarshalDnsRecordsObjectResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DnsRecordsObjectResult)
	err = core.UnmarshalPrimitive(m, "recs_added", &obj.RecsAdded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_records_parsed", &obj.TotalRecordsParsed)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DnsRecordsObjectTiming : timing object.
type DnsRecordsObjectTiming struct {
	// start time.
	StartTime *string `json:"start_time,omitempty"`

	// end time.
	EndTime *string `json:"end_time,omitempty"`

	// process time.
	ProcessTime *int64 `json:"process_time,omitempty"`
}


// UnmarshalDnsRecordsObjectTiming unmarshals an instance of DnsRecordsObjectTiming from the specified map of raw messages.
func UnmarshalDnsRecordsObjectTiming(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DnsRecordsObjectTiming)
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "end_time", &obj.EndTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "process_time", &obj.ProcessTime)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetDnsRecordsBulkOptions : The GetDnsRecordsBulk options.
type GetDnsRecordsBulkOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDnsRecordsBulkOptions : Instantiate GetDnsRecordsBulkOptions
func (*DnsRecordBulkV1) NewGetDnsRecordsBulkOptions() *GetDnsRecordsBulkOptions {
	return &GetDnsRecordsBulkOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetDnsRecordsBulkOptions) SetHeaders(param map[string]string) *GetDnsRecordsBulkOptions {
	options.Headers = param
	return options
}

// PostDnsRecordsBulkOptions : The PostDnsRecordsBulk options.
type PostDnsRecordsBulkOptions struct {
	// file to upload.
	File io.ReadCloser `json:"file,omitempty"`

	// The content type of file.
	FileContentType *string `json:"file_content_type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostDnsRecordsBulkOptions : Instantiate PostDnsRecordsBulkOptions
func (*DnsRecordBulkV1) NewPostDnsRecordsBulkOptions() *PostDnsRecordsBulkOptions {
	return &PostDnsRecordsBulkOptions{}
}

// SetFile : Allow user to set File
func (options *PostDnsRecordsBulkOptions) SetFile(file io.ReadCloser) *PostDnsRecordsBulkOptions {
	options.File = file
	return options
}

// SetFileContentType : Allow user to set FileContentType
func (options *PostDnsRecordsBulkOptions) SetFileContentType(fileContentType string) *PostDnsRecordsBulkOptions {
	options.FileContentType = core.StringPtr(fileContentType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PostDnsRecordsBulkOptions) SetHeaders(param map[string]string) *PostDnsRecordsBulkOptions {
	options.Headers = param
	return options
}

// DnsRecordsObject : dns records objects.
type DnsRecordsObject struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages []DnsRecordsObjectMessagesItem `json:"messages" validate:"required"`

	// DNS record.
	Result *DnsRecordsObjectResult `json:"result" validate:"required"`

	// timing object.
	Timing *DnsRecordsObjectTiming `json:"timing,omitempty"`
}


// UnmarshalDnsRecordsObject unmarshals an instance of DnsRecordsObject from the specified map of raw messages.
func UnmarshalDnsRecordsObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DnsRecordsObject)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalDnsRecordsObjectMessagesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDnsRecordsObjectResult)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "timing", &obj.Timing, UnmarshalDnsRecordsObjectTiming)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
