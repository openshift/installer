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
 * IBM OpenAPI SDK Code Generator Version: 3.34.0-e2a502a2-20210616-185634
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
	common "github.com/IBM/scc-go-sdk/common"
)

// AdminServiceApiV1 : This is an API for the Admin Service
//
// Version: 1.0.0
type AdminServiceApiV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://admin-service-api.cloud.ibm.com"

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
	return "", fmt.Errorf("service does not support regional URLs")
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

// AccountSettings : Account settings.
type AccountSettings struct {
	// Location settings.
	Location *LocationID `json:"location" validate:"required"`
}

// NewAccountSettings : Instantiate AccountSettings (Generic Model Constructor)
func (*AdminServiceApiV1) NewAccountSettings(location *LocationID) (_model *AccountSettings, err error) {
	_model = &AccountSettings{
		Location: location,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalAccountSettings unmarshals an instance of AccountSettings from the specified map of raw messages.
func UnmarshalAccountSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettings)
	err = core.UnmarshalModel(m, "location", &obj.Location, UnmarshalLocationID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetLocationOptions : The GetLocation options.
type GetLocationOptions struct {
	// The programatic ID of the location that you want to work in.
	LocationID *string `validate:"required,ne="`

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
	AccountID *string `validate:"required,ne="`

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

	MainEndpointURL *string `json:"main_endpoint_url,omitempty"`

	GovernanceEndpointURL *string `json:"governance_endpoint_url,omitempty"`

	ResultsEndpointURL *string `json:"results_endpoint_url,omitempty"`

	ComplianceEndpointURL *string `json:"compliance_endpoint_url,omitempty"`

	AnalyticsEndpointURL *string `json:"analytics_endpoint_url,omitempty"`

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

// PatchAccountSettingsOptions : The PatchAccountSettings options.
type PatchAccountSettingsOptions struct {
	// The ID of the managing account.
	AccountID *string `validate:"required,ne="`

	// Location settings.
	Location *LocationID `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPatchAccountSettingsOptions : Instantiate PatchAccountSettingsOptions
func (*AdminServiceApiV1) NewPatchAccountSettingsOptions(accountID string, location *LocationID) *PatchAccountSettingsOptions {
	return &PatchAccountSettingsOptions{
		AccountID: core.StringPtr(accountID),
		Location:  location,
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

// SetHeaders : Allow user to set Headers
func (options *PatchAccountSettingsOptions) SetHeaders(param map[string]string) *PatchAccountSettingsOptions {
	options.Headers = param
	return options
}

// Region : Region.
type Region struct {
	ID *string `json:"id" validate:"required"`
}

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
