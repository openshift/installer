/**
 * (C) Copyright IBM Corp. 2022.
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
 * IBM OpenAPI SDK Code Generator Version: 3.43.3-d49d4b21-20220104-223519
 */

// Package adminserviceapiv1 : Operations and models for the AdminServiceApiV1 service
package adminserviceapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/scc-go-sdk/v3/common"
)

// AdminServiceApiV1 : This is an API for the Admin Service
//
// API Version: 1.0.0
type AdminServiceApiV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us.compliance.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "admin_service_api"

// AdminServiceApiV1Options : Service options
type AdminServiceApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewAdminServiceApiV1UsingExternalConfig : constructs an instance of AdminServiceApiV1 with passed in options and external configuration.
func NewAdminServiceApiV1UsingExternalConfig(options *AdminServiceApiV1Options) (adminServiceApi *AdminServiceApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	adminServiceApi, err = NewAdminServiceApiV1(options)
	if err != nil {
		return
	}

	err = adminServiceApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = adminServiceApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewAdminServiceApiV1 : constructs an instance of AdminServiceApiV1 with passed in options.
func NewAdminServiceApiV1(options *AdminServiceApiV1Options) (service *AdminServiceApiV1, err error) {
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

	service = &AdminServiceApiV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south": "https://us.compliance.cloud.ibm.com",
		"us-east":  "https://us.compliance.cloud.ibm.com",
		"eu-de":    "https://eu.compliance.cloud.ibm.com",
		"eu-gb":    "https://uk.compliance.cloud.ibm.com",
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}

// Clone makes a copy of "adminServiceApi" suitable for processing requests.
func (adminServiceApi *AdminServiceApiV1) Clone() *AdminServiceApiV1 {
	if core.IsNil(adminServiceApi) {
		return nil
	}
	clone := *adminServiceApi
	clone.Service = adminServiceApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (adminServiceApi *AdminServiceApiV1) SetServiceURL(url string) error {
	return adminServiceApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (adminServiceApi *AdminServiceApiV1) GetServiceURL() string {
	return adminServiceApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (adminServiceApi *AdminServiceApiV1) SetDefaultHeaders(headers http.Header) {
	adminServiceApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (adminServiceApi *AdminServiceApiV1) SetEnableGzipCompression(enableGzip bool) {
	adminServiceApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (adminServiceApi *AdminServiceApiV1) GetEnableGzipCompression() bool {
	return adminServiceApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (adminServiceApi *AdminServiceApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	adminServiceApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (adminServiceApi *AdminServiceApiV1) DisableRetries() {
	adminServiceApi.Service.DisableRetries()
}

// GetSettings : View account settings
// View the current settings for a specific account.
func (adminServiceApi *AdminServiceApiV1) GetSettings(getSettingsOptions *GetSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	return adminServiceApi.GetSettingsWithContext(context.Background(), getSettingsOptions)
}

// GetSettingsWithContext is an alternate form of the GetSettings method which supports a Context parameter
func (adminServiceApi *AdminServiceApiV1) GetSettingsWithContext(ctx context.Context, getSettingsOptions *GetSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSettingsOptions, "getSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSettingsOptions, "getSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = adminServiceApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(adminServiceApi.Service.Options.URL, `/admin/v1/accounts/{account_id}/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("admin_service_api", "V1", "GetSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = adminServiceApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PatchAccountSettings : Update account settings
// Update the settings for a specific account.
func (adminServiceApi *AdminServiceApiV1) PatchAccountSettings(patchAccountSettingsOptions *PatchAccountSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	return adminServiceApi.PatchAccountSettingsWithContext(context.Background(), patchAccountSettingsOptions)
}

// PatchAccountSettingsWithContext is an alternate form of the PatchAccountSettings method which supports a Context parameter
func (adminServiceApi *AdminServiceApiV1) PatchAccountSettingsWithContext(ctx context.Context, patchAccountSettingsOptions *PatchAccountSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(patchAccountSettingsOptions, "patchAccountSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(patchAccountSettingsOptions, "patchAccountSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *patchAccountSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = adminServiceApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(adminServiceApi.Service.Options.URL, `/admin/v1/accounts/{account_id}/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range patchAccountSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("admin_service_api", "V1", "PatchAccountSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if patchAccountSettingsOptions.Location != nil {
		body["location"] = patchAccountSettingsOptions.Location
	}
	if patchAccountSettingsOptions.EventNotifications != nil {
		body["event_notifications"] = patchAccountSettingsOptions.EventNotifications
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
	response, err = adminServiceApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListLocations : View available locations
// View the available locations in which the data that is generated by the Security and Compliance Center can be
// managed.
func (adminServiceApi *AdminServiceApiV1) ListLocations(listLocationsOptions *ListLocationsOptions) (result *Locations, response *core.DetailedResponse, err error) {
	return adminServiceApi.ListLocationsWithContext(context.Background(), listLocationsOptions)
}

// ListLocationsWithContext is an alternate form of the ListLocations method which supports a Context parameter
func (adminServiceApi *AdminServiceApiV1) ListLocationsWithContext(ctx context.Context, listLocationsOptions *ListLocationsOptions) (result *Locations, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listLocationsOptions, "listLocationsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = adminServiceApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(adminServiceApi.Service.Options.URL, `/admin/v1/locations`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listLocationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("admin_service_api", "V1", "ListLocations")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = adminServiceApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLocations)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLocation : View the details of a location
// View the endpoints and regions that are available for a specific region.
func (adminServiceApi *AdminServiceApiV1) GetLocation(getLocationOptions *GetLocationOptions) (result *Location, response *core.DetailedResponse, err error) {
	return adminServiceApi.GetLocationWithContext(context.Background(), getLocationOptions)
}

// GetLocationWithContext is an alternate form of the GetLocation method which supports a Context parameter
func (adminServiceApi *AdminServiceApiV1) GetLocationWithContext(ctx context.Context, getLocationOptions *GetLocationOptions) (result *Location, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLocationOptions, "getLocationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLocationOptions, "getLocationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *getLocationOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = adminServiceApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(adminServiceApi.Service.Options.URL, `/admin/v1/locations/{location_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLocationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("admin_service_api", "V1", "GetLocation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = adminServiceApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLocation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SendTestEvent : Send test event
// Send a test event using your configured Event Notifications instance.
func (adminServiceApi *AdminServiceApiV1) SendTestEvent(sendTestEventOptions *SendTestEventOptions) (result *TestEvent, response *core.DetailedResponse, err error) {
	return adminServiceApi.SendTestEventWithContext(context.Background(), sendTestEventOptions)
}

// SendTestEventWithContext is an alternate form of the SendTestEvent method which supports a Context parameter
func (adminServiceApi *AdminServiceApiV1) SendTestEventWithContext(ctx context.Context, sendTestEventOptions *SendTestEventOptions) (result *TestEvent, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(sendTestEventOptions, "sendTestEventOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(sendTestEventOptions, "sendTestEventOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *sendTestEventOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = adminServiceApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(adminServiceApi.Service.Options.URL, `/admin/v1/accounts/{account_id}/test_event`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range sendTestEventOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("admin_service_api", "V1", "SendTestEvent")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = adminServiceApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTestEvent)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// AccountSettings : Account settings.
type AccountSettings struct {
	// Location settings.
	Location *LocationID `json:"location,omitempty"`

	// The Event Notification settings to register.
	EventNotifications *NotificationsRegistration `json:"event_notifications,omitempty"`
}

// UnmarshalAccountSettings unmarshals an instance of AccountSettings from the specified map of raw messages.
func UnmarshalAccountSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettings)
	err = core.UnmarshalModel(m, "location", &obj.Location, UnmarshalLocationID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "event_notifications", &obj.EventNotifications, UnmarshalNotificationsRegistration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetLocationOptions : The GetLocation options.
type GetLocationOptions struct {
	// The programatic ID of the location that you want to work in.
	LocationID *string `json:"location_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetLocationOptions.LocationID property.
// The programatic ID of the location that you want to work in.
const (
	GetLocationOptions_LocationID_Eu = "eu"
	GetLocationOptions_LocationID_Uk = "uk"
	GetLocationOptions_LocationID_Us = "us"
)

// NewGetLocationOptions : Instantiate GetLocationOptions
func (*AdminServiceApiV1) NewGetLocationOptions(locationID string) *GetLocationOptions {
	return &GetLocationOptions{
		LocationID: core.StringPtr(locationID),
	}
}

// SetLocationID : Allow user to set LocationID
func (_options *GetLocationOptions) SetLocationID(locationID string) *GetLocationOptions {
	_options.LocationID = core.StringPtr(locationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLocationOptions) SetHeaders(param map[string]string) *GetLocationOptions {
	options.Headers = param
	return options
}

// GetSettingsOptions : The GetSettings options.
type GetSettingsOptions struct {
	// The ID of the managing account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSettingsOptions : Instantiate GetSettingsOptions
func (*AdminServiceApiV1) NewGetSettingsOptions(accountID string) *GetSettingsOptions {
	return &GetSettingsOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetSettingsOptions) SetAccountID(accountID string) *GetSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSettingsOptions) SetHeaders(param map[string]string) *GetSettingsOptions {
	options.Headers = param
	return options
}

// ListLocationsOptions : The ListLocations options.
type ListLocationsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListLocationsOptions : Instantiate ListLocationsOptions
func (*AdminServiceApiV1) NewListLocationsOptions() *ListLocationsOptions {
	return &ListLocationsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListLocationsOptions) SetHeaders(param map[string]string) *ListLocationsOptions {
	options.Headers = param
	return options
}

// Location : The location that your account is current set to use.
type Location struct {
	// The programatic ID of the location that you want to work in.
	ID *string `json:"id,omitempty"`

	// The base URL for the service.
	MainEndpointURL *string `json:"main_endpoint_url,omitempty"`

	// The endpoint that is used to call the Configuration Governance APIs.
	GovernanceEndpointURL *string `json:"governance_endpoint_url,omitempty"`

	// The endpoint that is used to get the results for the Configuration Governance component.
	ResultsEndpointURL *string `json:"results_endpoint_url,omitempty"`

	// The endpoint that is used to call the Posture Management APIs.
	ComplianceEndpointURL *string `json:"compliance_endpoint_url,omitempty"`

	// The endpoint that is used to generate analytics for the Posture Management component.
	AnalyticsEndpointURL *string `json:"analytics_endpoint_url,omitempty"`

	// The endpoint that is used to call the Security Insights APIs.
	SiEndpointURL *string `json:"si_endpoint_url,omitempty"`

	Regions []Region `json:"regions,omitempty"`
}

// Constants associated with the Location.ID property.
// The programatic ID of the location that you want to work in.
const (
	Location_ID_Eu = "eu"
	Location_ID_Uk = "uk"
	Location_ID_Us = "us"
)

// UnmarshalLocation unmarshals an instance of Location from the specified map of raw messages.
func UnmarshalLocation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Location)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "main_endpoint_url", &obj.MainEndpointURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "governance_endpoint_url", &obj.GovernanceEndpointURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "results_endpoint_url", &obj.ResultsEndpointURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliance_endpoint_url", &obj.ComplianceEndpointURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "analytics_endpoint_url", &obj.AnalyticsEndpointURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "si_endpoint_url", &obj.SiEndpointURL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "regions", &obj.Regions, UnmarshalRegion)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LocationID : Location settings.
type LocationID struct {
	// The programatic ID of the location that you want to work in.
	ID *string `json:"id" validate:"required"`
}

// Constants associated with the LocationID.ID property.
// The programatic ID of the location that you want to work in.
const (
	LocationID_ID_Eu = "eu"
	LocationID_ID_Uk = "uk"
	LocationID_ID_Us = "us"
)

// NewLocationID : Instantiate LocationID (Generic Model Constructor)
func (*AdminServiceApiV1) NewLocationID(id string) (_model *LocationID, err error) {
	_model = &LocationID{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalLocationID unmarshals an instance of LocationID from the specified map of raw messages.
func UnmarshalLocationID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LocationID)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Locations : An array of available locations.
type Locations struct {
	Locations []Location `json:"locations" validate:"required"`
}

// UnmarshalLocations unmarshals an instance of Locations from the specified map of raw messages.
func UnmarshalLocations(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Locations)
	err = core.UnmarshalModel(m, "locations", &obj.Locations, UnmarshalLocation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NotificationsRegistration : The Event Notification settings to register.
type NotificationsRegistration struct {
	// The Cloud Resource Name (CRN) of the Event Notifications instance that you want to connect.
	InstanceCrn *string `json:"instance_crn" validate:"required"`

	// The name to register as a source in your Event Notifications instance.
	SourceName *string `json:"source_name,omitempty"`

	// An optional description for the source in your Event Notifications instance.
	SourceDescription *string `json:"source_description,omitempty"`
}

// NewNotificationsRegistration : Instantiate NotificationsRegistration (Generic Model Constructor)
func (*AdminServiceApiV1) NewNotificationsRegistration(instanceCrn string) (_model *NotificationsRegistration, err error) {
	_model = &NotificationsRegistration{
		InstanceCrn: core.StringPtr(instanceCrn),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalNotificationsRegistration unmarshals an instance of NotificationsRegistration from the specified map of raw messages.
func UnmarshalNotificationsRegistration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationsRegistration)
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_name", &obj.SourceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_description", &obj.SourceDescription)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PatchAccountSettingsOptions : The PatchAccountSettings options.
type PatchAccountSettingsOptions struct {
	// The ID of the managing account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Location settings.
	Location *LocationID `json:"location,omitempty"`

	// The Event Notification settings to register.
	EventNotifications *NotificationsRegistration `json:"event_notifications,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPatchAccountSettingsOptions : Instantiate PatchAccountSettingsOptions
func (*AdminServiceApiV1) NewPatchAccountSettingsOptions(accountID string) *PatchAccountSettingsOptions {
	return &PatchAccountSettingsOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *PatchAccountSettingsOptions) SetAccountID(accountID string) *PatchAccountSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *PatchAccountSettingsOptions) SetLocation(location *LocationID) *PatchAccountSettingsOptions {
	_options.Location = location
	return _options
}

// SetEventNotifications : Allow user to set EventNotifications
func (_options *PatchAccountSettingsOptions) SetEventNotifications(eventNotifications *NotificationsRegistration) *PatchAccountSettingsOptions {
	_options.EventNotifications = eventNotifications
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PatchAccountSettingsOptions) SetHeaders(param map[string]string) *PatchAccountSettingsOptions {
	options.Headers = param
	return options
}

// Region : The region or regions that are available for each location. Be sure to use the correct region ID when making your API
// call.
type Region struct {
	// The programatic ID of the available regions.
	ID *string `json:"id" validate:"required"`
}

// Constants associated with the Region.ID property.
// The programatic ID of the available regions.
const (
	Region_ID_Eu = "eu"
	Region_ID_Uk = "uk"
	Region_ID_Us = "us"
)

// UnmarshalRegion unmarshals an instance of Region from the specified map of raw messages.
func UnmarshalRegion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Region)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SendTestEventOptions : The SendTestEvent options.
type SendTestEventOptions struct {
	// The ID of the managing account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSendTestEventOptions : Instantiate SendTestEventOptions
func (*AdminServiceApiV1) NewSendTestEventOptions(accountID string) *SendTestEventOptions {
	return &SendTestEventOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *SendTestEventOptions) SetAccountID(accountID string) *SendTestEventOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SendTestEventOptions) SetHeaders(param map[string]string) *SendTestEventOptions {
	options.Headers = param
	return options
}

// TestEvent : The details of a test event response.
type TestEvent struct {
	// Indicates whether the event was received by Event Notifications.
	Success *bool `json:"success" validate:"required"`
}

// UnmarshalTestEvent unmarshals an instance of TestEvent from the specified map of raw messages.
func UnmarshalTestEvent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TestEvent)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
