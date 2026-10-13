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

// Package dnszonesv1 : Operations and models for the DnsZonesV1 service
package dnszonesv1

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
)

// DnsZonesV1 : DNS Zones
//
// Version: 1.0.0
type DnsZonesV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.dns-svcs.cloud.ibm.com/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "dns_zones"

// DnsZonesV1Options : Service options
type DnsZonesV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewDnsZonesV1UsingExternalConfig : constructs an instance of DnsZonesV1 with passed in options and external configuration.
func NewDnsZonesV1UsingExternalConfig(options *DnsZonesV1Options) (dnsZones *DnsZonesV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	dnsZones, err = NewDnsZonesV1(options)
	if err != nil {
		return
	}

	err = dnsZones.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = dnsZones.Service.SetServiceURL(options.URL)
	}
	return
}

// NewDnsZonesV1 : constructs an instance of DnsZonesV1 with passed in options.
func NewDnsZonesV1(options *DnsZonesV1Options) (service *DnsZonesV1, err error) {
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

	service = &DnsZonesV1{
		Service: baseService,
	}

	return
}

// SetServiceURL sets the service URL
func (dnsZones *DnsZonesV1) SetServiceURL(url string) error {
	return dnsZones.Service.SetServiceURL(url)
}

// ListDnszones : List DNS zones
// List the DNS zones for a given service instance.
func (dnsZones *DnsZonesV1) ListDnszones(listDnszonesOptions *ListDnszonesOptions) (result *ListDnszones, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listDnszonesOptions, "listDnszonesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listDnszonesOptions, "listDnszonesOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones"}
	pathParameters := []string{*listDnszonesOptions.InstanceID}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(dnsZones.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range listDnszonesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_zones", "V1", "ListDnszones")
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = dnsZones.Service.Request(request, &rawResponse)
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
func (dnsZones *DnsZonesV1) CreateDnszone(createDnszoneOptions *CreateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDnszoneOptions, "createDnszoneOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createDnszoneOptions, "createDnszoneOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones"}
	pathParameters := []string{*createDnszoneOptions.InstanceID}

	builder := core.NewRequestBuilder(core.POST)
	_, err = builder.ConstructHTTPURL(dnsZones.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range createDnszoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_zones", "V1", "CreateDnszone")
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
	response, err = dnsZones.Service.Request(request, &rawResponse)
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
func (dnsZones *DnsZonesV1) DeleteDnszone(deleteDnszoneOptions *DeleteDnszoneOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDnszoneOptions, "deleteDnszoneOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteDnszoneOptions, "deleteDnszoneOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones"}
	pathParameters := []string{*deleteDnszoneOptions.InstanceID, *deleteDnszoneOptions.DnszoneID}

	builder := core.NewRequestBuilder(core.DELETE)
	_, err = builder.ConstructHTTPURL(dnsZones.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteDnszoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_zones", "V1", "DeleteDnszone")
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

	response, err = dnsZones.Service.Request(request, nil)

	return
}

// GetDnszone : Get a DNS zone
// Get details of a DNS zone.
func (dnsZones *DnsZonesV1) GetDnszone(getDnszoneOptions *GetDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDnszoneOptions, "getDnszoneOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDnszoneOptions, "getDnszoneOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones"}
	pathParameters := []string{*getDnszoneOptions.InstanceID, *getDnszoneOptions.DnszoneID}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(dnsZones.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDnszoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_zones", "V1", "GetDnszone")
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
	response, err = dnsZones.Service.Request(request, &rawResponse)
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
func (dnsZones *DnsZonesV1) UpdateDnszone(updateDnszoneOptions *UpdateDnszoneOptions) (result *Dnszone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDnszoneOptions, "updateDnszoneOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateDnszoneOptions, "updateDnszoneOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"instances", "dnszones"}
	pathParameters := []string{*updateDnszoneOptions.InstanceID, *updateDnszoneOptions.DnszoneID}

	builder := core.NewRequestBuilder(core.PATCH)
	_, err = builder.ConstructHTTPURL(dnsZones.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateDnszoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("dns_zones", "V1", "UpdateDnszone")
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
	response, err = dnsZones.Service.Request(request, &rawResponse)
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

// CreateDnszoneOptions : The CreateDnszone options.
type CreateDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

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
func (*DnsZonesV1) NewCreateDnszoneOptions(instanceID string) *CreateDnszoneOptions {
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

// DeleteDnszoneOptions : The DeleteDnszone options.
type DeleteDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteDnszoneOptions : Instantiate DeleteDnszoneOptions
func (*DnsZonesV1) NewDeleteDnszoneOptions(instanceID string, dnszoneID string) *DeleteDnszoneOptions {
	return &DeleteDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
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

// GetDnszoneOptions : The GetDnszone options.
type GetDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDnszoneOptions : Instantiate GetDnszoneOptions
func (*DnsZonesV1) NewGetDnszoneOptions(instanceID string, dnszoneID string) *GetDnszoneOptions {
	return &GetDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
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

// ListDnszonesOptions : The ListDnszones options.
type ListDnszonesOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Specify how many DNS zones to skip over, the default value is 0.
	Offset *int64 `json:"offset,omitempty"`

	// Specify how many DNS zones are returned, the default value is 10.
	Limit *int64 `json:"limit,omitempty"`

	// Specify the VPC instance id.
	VpcID *string `json:"vpc_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListDnszonesOptions : Instantiate ListDnszonesOptions
func (*DnsZonesV1) NewListDnszonesOptions(instanceID string) *ListDnszonesOptions {
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

// SetVpcID : Allow user to set VpcID
func (options *ListDnszonesOptions) SetVpcID(vpcID string) *ListDnszonesOptions {
	options.VpcID = core.StringPtr(vpcID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListDnszonesOptions) SetHeaders(param map[string]string) *ListDnszonesOptions {
	options.Headers = param
	return options
}

// UpdateDnszoneOptions : The UpdateDnszone options.
type UpdateDnszoneOptions struct {
	// The unique identifier of a service instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The unique identifier of a DNS zone.
	DnszoneID *string `json:"dnszone_id" validate:"required"`

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
func (*DnsZonesV1) NewUpdateDnszoneOptions(instanceID string, dnszoneID string) *UpdateDnszoneOptions {
	return &UpdateDnszoneOptions{
		InstanceID: core.StringPtr(instanceID),
		DnszoneID:  core.StringPtr(dnszoneID),
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
