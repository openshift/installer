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

// Package resourcerecordsv1 : Operations and models for the ResourceRecordsV1 service
// Deprecated: This module is deprecated. For more detail https://github.com/IBM/networking-go-sdk/blob/master/DEPRECATION-NOTICE.md
package resourcerecordsv1

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
)

// ResourceRecordsV1 : Resource Records
//
// Version: 1.0.0
type ResourceRecordsV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.dns-svcs.cloud.ibm.com/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "resource_records"

// ResourceRecordsV1Options : Service options
type ResourceRecordsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewResourceRecordsV1UsingExternalConfig : constructs an instance of ResourceRecordsV1 with passed in options and external configuration.
func NewResourceRecordsV1UsingExternalConfig(options *ResourceRecordsV1Options) (resourceRecords *ResourceRecordsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	resourceRecords, err = NewResourceRecordsV1(options)
	if err != nil {
		return
	}

	err = resourceRecords.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = resourceRecords.Service.SetServiceURL(options.URL)
	}
	return
}

// NewResourceRecordsV1 : constructs an instance of ResourceRecordsV1 with passed in options.
func NewResourceRecordsV1(options *ResourceRecordsV1Options) (service *ResourceRecordsV1, err error) {
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

	service = &ResourceRecordsV1{
		Service: baseService,
	}

	return
}

// SetServiceURL sets the service URL
func (resourceRecords *ResourceRecordsV1) SetServiceURL(url string) error {
	return resourceRecords.Service.SetServiceURL(url)
}

// ListResourceRecords : List Resource Records
// List the Resource Records for a given DNS zone.
func (resourceRecords *ResourceRecordsV1) ListResourceRecords(listResourceRecordsOptions *ListResourceRecordsOptions) (result *ListResourceRecords, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listResourceRecordsOptions, "listResourceRecordsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listResourceRecordsOptions, "listResourceRecordsOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones", "resource_records"}
	pathParameters := []string{*listResourceRecordsOptions.InstanceID, *listResourceRecordsOptions.DnszoneID}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(resourceRecords.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range listResourceRecordsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_records", "V1", "ListResourceRecords")
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
	response, err = resourceRecords.Service.Request(request, &rawResponse)
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
func (resourceRecords *ResourceRecordsV1) CreateResourceRecord(createResourceRecordOptions *CreateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createResourceRecordOptions, "createResourceRecordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createResourceRecordOptions, "createResourceRecordOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones", "resource_records"}
	pathParameters := []string{*createResourceRecordOptions.InstanceID, *createResourceRecordOptions.DnszoneID}

	builder := core.NewRequestBuilder(core.POST)
	_, err = builder.ConstructHTTPURL(resourceRecords.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range createResourceRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_records", "V1", "CreateResourceRecord")
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
	response, err = resourceRecords.Service.Request(request, &rawResponse)
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
func (resourceRecords *ResourceRecordsV1) DeleteResourceRecord(deleteResourceRecordOptions *DeleteResourceRecordOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteResourceRecordOptions, "deleteResourceRecordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteResourceRecordOptions, "deleteResourceRecordOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones", "resource_records"}
	pathParameters := []string{*deleteResourceRecordOptions.InstanceID, *deleteResourceRecordOptions.DnszoneID, *deleteResourceRecordOptions.RecordID}

	builder := core.NewRequestBuilder(core.DELETE)
	_, err = builder.ConstructHTTPURL(resourceRecords.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteResourceRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_records", "V1", "DeleteResourceRecord")
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

	response, err = resourceRecords.Service.Request(request, nil)

	return
}

// GetResourceRecord : Get a resource record
// Get details of a resource record.
func (resourceRecords *ResourceRecordsV1) GetResourceRecord(getResourceRecordOptions *GetResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceRecordOptions, "getResourceRecordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getResourceRecordOptions, "getResourceRecordOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones", "resource_records"}
	pathParameters := []string{*getResourceRecordOptions.InstanceID, *getResourceRecordOptions.DnszoneID, *getResourceRecordOptions.RecordID}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(resourceRecords.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getResourceRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_records", "V1", "GetResourceRecord")
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
	response, err = resourceRecords.Service.Request(request, &rawResponse)
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
func (resourceRecords *ResourceRecordsV1) UpdateResourceRecord(updateResourceRecordOptions *UpdateResourceRecordOptions) (result *ResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateResourceRecordOptions, "updateResourceRecordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateResourceRecordOptions, "updateResourceRecordOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones", "resource_records"}
	pathParameters := []string{*updateResourceRecordOptions.InstanceID, *updateResourceRecordOptions.DnszoneID, *updateResourceRecordOptions.RecordID}

	builder := core.NewRequestBuilder(core.PUT)
	_, err = builder.ConstructHTTPURL(resourceRecords.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateResourceRecordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_records", "V1", "UpdateResourceRecord")
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
	response, err = resourceRecords.Service.Request(request, &rawResponse)
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

// CreateResourceRecordOptions : The CreateResourceRecord options.
type CreateResourceRecordOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required"`

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
	CreateResourceRecordOptions_Type_A     = "A"
	CreateResourceRecordOptions_Type_Aaaa  = "AAAA"
	CreateResourceRecordOptions_Type_Cname = "CNAME"
	CreateResourceRecordOptions_Type_Mx    = "MX"
	CreateResourceRecordOptions_Type_Ptr   = "PTR"
	CreateResourceRecordOptions_Type_Srv   = "SRV"
	CreateResourceRecordOptions_Type_Txt   = "TXT"
)

// NewCreateResourceRecordOptions : Instantiate CreateResourceRecordOptions
func (*ResourceRecordsV1) NewCreateResourceRecordOptions(instanceID string, dnszoneID string) *CreateResourceRecordOptions {
	return &CreateResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
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

// DeleteResourceRecordOptions : The DeleteResourceRecord options.
type DeleteResourceRecordOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required"`

	// The unique identifier of a resource record.
	RecordID *string `json:"record_id" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteResourceRecordOptions : Instantiate DeleteResourceRecordOptions
func (*ResourceRecordsV1) NewDeleteResourceRecordOptions(instanceID string, dnszoneID string, recordID string) *DeleteResourceRecordOptions {
	return &DeleteResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		RecordID:   core.StringPtr(recordID),
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

// GetResourceRecordOptions : The GetResourceRecord options.
type GetResourceRecordOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required"`

	// The unique identifier of a resource record.
	RecordID *string `json:"record_id" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetResourceRecordOptions : Instantiate GetResourceRecordOptions
func (*ResourceRecordsV1) NewGetResourceRecordOptions(instanceID string, dnszoneID string, recordID string) *GetResourceRecordOptions {
	return &GetResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		RecordID:   core.StringPtr(recordID),
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

// ListResourceRecordsOptions : The ListResourceRecords options.
type ListResourceRecordsOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required"`

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
func (*ResourceRecordsV1) NewListResourceRecordsOptions(instanceID string, dnszoneID string) *ListResourceRecordsOptions {
	return &ListResourceRecordsOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
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

// UpdateResourceRecordOptions : The UpdateResourceRecord options.
type UpdateResourceRecordOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required"`

	// The unique identifier of a resource record.
	RecordID *string `json:"record_id" validate:"required"`

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
func (*ResourceRecordsV1) NewUpdateResourceRecordOptions(instanceID string, dnszoneID string, recordID string) *UpdateResourceRecordOptions {
	return &UpdateResourceRecordOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
		RecordID:   core.StringPtr(recordID),
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
func (*ResourceRecordsV1) NewResourceRecordInputRdataRdataARecord(ip string) (model *ResourceRecordInputRdataRdataARecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordInputRdataRdataAaaaRecord(ip string) (model *ResourceRecordInputRdataRdataAaaaRecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordInputRdataRdataCnameRecord(cname string) (model *ResourceRecordInputRdataRdataCnameRecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordInputRdataRdataMxRecord(exchange string, preference int64) (model *ResourceRecordInputRdataRdataMxRecord, err error) {
	model = &ResourceRecordInputRdataRdataMxRecord{
		Exchange:   core.StringPtr(exchange),
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
func (*ResourceRecordsV1) NewResourceRecordInputRdataRdataPtrRecord(ptrdname string) (model *ResourceRecordInputRdataRdataPtrRecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordInputRdataRdataSrvRecord(port int64, priority int64, target string, weight int64) (model *ResourceRecordInputRdataRdataSrvRecord, err error) {
	model = &ResourceRecordInputRdataRdataSrvRecord{
		Port:     core.Int64Ptr(port),
		Priority: core.Int64Ptr(priority),
		Target:   core.StringPtr(target),
		Weight:   core.Int64Ptr(weight),
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
func (*ResourceRecordsV1) NewResourceRecordInputRdataRdataTxtRecord(text string) (model *ResourceRecordInputRdataRdataTxtRecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordUpdateInputRdataRdataARecord(ip string) (model *ResourceRecordUpdateInputRdataRdataARecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordUpdateInputRdataRdataAaaaRecord(ip string) (model *ResourceRecordUpdateInputRdataRdataAaaaRecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordUpdateInputRdataRdataCnameRecord(cname string) (model *ResourceRecordUpdateInputRdataRdataCnameRecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordUpdateInputRdataRdataMxRecord(exchange string, preference int64) (model *ResourceRecordUpdateInputRdataRdataMxRecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataMxRecord{
		Exchange:   core.StringPtr(exchange),
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
func (*ResourceRecordsV1) NewResourceRecordUpdateInputRdataRdataPtrRecord(ptrdname string) (model *ResourceRecordUpdateInputRdataRdataPtrRecord, err error) {
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
func (*ResourceRecordsV1) NewResourceRecordUpdateInputRdataRdataSrvRecord(port int64, priority int64, target string, weight int64) (model *ResourceRecordUpdateInputRdataRdataSrvRecord, err error) {
	model = &ResourceRecordUpdateInputRdataRdataSrvRecord{
		Port:     core.Int64Ptr(port),
		Priority: core.Int64Ptr(priority),
		Target:   core.StringPtr(target),
		Weight:   core.Int64Ptr(weight),
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
func (*ResourceRecordsV1) NewResourceRecordUpdateInputRdataRdataTxtRecord(text string) (model *ResourceRecordUpdateInputRdataRdataTxtRecord, err error) {
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
