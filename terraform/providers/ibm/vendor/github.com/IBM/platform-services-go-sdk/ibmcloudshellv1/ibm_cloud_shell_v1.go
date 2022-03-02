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
 * IBM OpenAPI SDK Code Generator Version: 3.33.0-caf29bd0-20210603-225214
 */

// Package ibmcloudshellv1 : Operations and models for the IBMCloudShellV1 service
package ibmcloudshellv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
)

// IBMCloudShellV1 : API docs for IBM Cloud Shell repository
//
// Version: 1.0
type IBMCloudShellV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.shell.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "ibm_cloud_shell"

// IBMCloudShellV1Options : Service options
type IBMCloudShellV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewIBMCloudShellV1UsingExternalConfig : constructs an instance of IBMCloudShellV1 with passed in options and external configuration.
func NewIBMCloudShellV1UsingExternalConfig(options *IBMCloudShellV1Options) (ibmCloudShell *IBMCloudShellV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	ibmCloudShell, err = NewIBMCloudShellV1(options)
	if err != nil {
		return
	}

	err = ibmCloudShell.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = ibmCloudShell.Service.SetServiceURL(options.URL)
	}
	return
}

// NewIBMCloudShellV1 : constructs an instance of IBMCloudShellV1 with passed in options.
func NewIBMCloudShellV1(options *IBMCloudShellV1Options) (service *IBMCloudShellV1, err error) {
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

	service = &IBMCloudShellV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "ibmCloudShell" suitable for processing requests.
func (ibmCloudShell *IBMCloudShellV1) Clone() *IBMCloudShellV1 {
	if core.IsNil(ibmCloudShell) {
		return nil
	}
	clone := *ibmCloudShell
	clone.Service = ibmCloudShell.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (ibmCloudShell *IBMCloudShellV1) SetServiceURL(url string) error {
	return ibmCloudShell.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (ibmCloudShell *IBMCloudShellV1) GetServiceURL() string {
	return ibmCloudShell.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (ibmCloudShell *IBMCloudShellV1) SetDefaultHeaders(headers http.Header) {
	ibmCloudShell.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (ibmCloudShell *IBMCloudShellV1) SetEnableGzipCompression(enableGzip bool) {
	ibmCloudShell.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (ibmCloudShell *IBMCloudShellV1) GetEnableGzipCompression() bool {
	return ibmCloudShell.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (ibmCloudShell *IBMCloudShellV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	ibmCloudShell.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (ibmCloudShell *IBMCloudShellV1) DisableRetries() {
	ibmCloudShell.Service.DisableRetries()
}

// GetAccountSettings : Get account settings
// Retrieve account settings for the given account ID. Call this method to get details about a particular account
// setting, whether Cloud Shell is enabled, the list of enabled regions and the list of enabled features. Users need to
// be an account owner or users need to be assigned an IAM policy with the Administrator role for the Cloud Shell
// account management service.
func (ibmCloudShell *IBMCloudShellV1) GetAccountSettings(getAccountSettingsOptions *GetAccountSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	return ibmCloudShell.GetAccountSettingsWithContext(context.Background(), getAccountSettingsOptions)
}

// GetAccountSettingsWithContext is an alternate form of the GetAccountSettings method which supports a Context parameter
func (ibmCloudShell *IBMCloudShellV1) GetAccountSettingsWithContext(ctx context.Context, getAccountSettingsOptions *GetAccountSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountSettingsOptions, "getAccountSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAccountSettingsOptions, "getAccountSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getAccountSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = ibmCloudShell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(ibmCloudShell.Service.Options.URL, `/api/v1/user/accounts/{account_id}/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAccountSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ibm_cloud_shell", "V1", "GetAccountSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = ibmCloudShell.Service.Request(request, &rawResponse)
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

// UpdateAccountSettings : Update account settings
// Update account settings for the given account ID. Call this method to update account settings configuration, you can
// enable or disable Cloud Shell, enable or disable available regions and enable and disable features. To update account
// settings, users need to be an account owner or users need to be assigned an IAM policy with the Administrator role
// for the Cloud Shell account management service.
func (ibmCloudShell *IBMCloudShellV1) UpdateAccountSettings(updateAccountSettingsOptions *UpdateAccountSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	return ibmCloudShell.UpdateAccountSettingsWithContext(context.Background(), updateAccountSettingsOptions)
}

// UpdateAccountSettingsWithContext is an alternate form of the UpdateAccountSettings method which supports a Context parameter
func (ibmCloudShell *IBMCloudShellV1) UpdateAccountSettingsWithContext(ctx context.Context, updateAccountSettingsOptions *UpdateAccountSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccountSettingsOptions, "updateAccountSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAccountSettingsOptions, "updateAccountSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *updateAccountSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = ibmCloudShell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(ibmCloudShell.Service.Options.URL, `/api/v1/user/accounts/{account_id}/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateAccountSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ibm_cloud_shell", "V1", "UpdateAccountSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAccountSettingsOptions.Rev != nil {
		body["_rev"] = updateAccountSettingsOptions.Rev
	}
	if updateAccountSettingsOptions.DefaultEnableNewFeatures != nil {
		body["default_enable_new_features"] = updateAccountSettingsOptions.DefaultEnableNewFeatures
	}
	if updateAccountSettingsOptions.DefaultEnableNewRegions != nil {
		body["default_enable_new_regions"] = updateAccountSettingsOptions.DefaultEnableNewRegions
	}
	if updateAccountSettingsOptions.Enabled != nil {
		body["enabled"] = updateAccountSettingsOptions.Enabled
	}
	if updateAccountSettingsOptions.Features != nil {
		body["features"] = updateAccountSettingsOptions.Features
	}
	if updateAccountSettingsOptions.Regions != nil {
		body["regions"] = updateAccountSettingsOptions.Regions
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
	response, err = ibmCloudShell.Service.Request(request, &rawResponse)
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

// AccountSettings : Definition of Cloud Shell account settings.
type AccountSettings struct {
	// Unique id of the settings object.
	ID *string `json:"_id,omitempty"`

	// Unique revision number for the settings object.
	Rev *string `json:"_rev,omitempty"`

	// The id of the account the settings belong to.
	AccountID *string `json:"account_id,omitempty"`

	// Creation timestamp in Unix epoch time.
	CreatedAt *int64 `json:"created_at,omitempty"`

	// IAM ID of creator.
	CreatedBy *string `json:"created_by,omitempty"`

	// You can choose which Cloud Shell features are available in the account and whether any new features are enabled as
	// they become available. The feature settings apply only to the enabled Cloud Shell locations.
	DefaultEnableNewFeatures *bool `json:"default_enable_new_features,omitempty"`

	// Set whether Cloud Shell is enabled in a specific location for the account. The location determines where user and
	// session data are stored. By default, users are routed to the nearest available location.
	DefaultEnableNewRegions *bool `json:"default_enable_new_regions,omitempty"`

	// When enabled, Cloud Shell is available to all users in the account.
	Enabled *bool `json:"enabled,omitempty"`

	// List of Cloud Shell features.
	Features []Feature `json:"features,omitempty"`

	// List of Cloud Shell region settings.
	Regions []RegionSetting `json:"regions,omitempty"`

	// Type of api response object.
	Type *string `json:"type,omitempty"`

	// Timestamp of last update in Unix epoch time.
	UpdatedAt *int64 `json:"updated_at,omitempty"`

	// IAM ID of last updater.
	UpdatedBy *string `json:"updated_by,omitempty"`
}

// UnmarshalAccountSettings unmarshals an instance of AccountSettings from the specified map of raw messages.
func UnmarshalAccountSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettings)
	err = core.UnmarshalPrimitive(m, "_id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_enable_new_features", &obj.DefaultEnableNewFeatures)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_enable_new_regions", &obj.DefaultEnableNewRegions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeature)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "regions", &obj.Regions, UnmarshalRegionSetting)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Feature : Describes a Cloud Shell feature.
type Feature struct {
	// State of the feature.
	Enabled *bool `json:"enabled,omitempty"`

	// Name of the feature.
	Key *string `json:"key,omitempty"`
}

// UnmarshalFeature unmarshals an instance of Feature from the specified map of raw messages.
func UnmarshalFeature(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Feature)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAccountSettingsOptions : The GetAccountSettings options.
type GetAccountSettingsOptions struct {
	// The account ID in which the account settings belong to.
	AccountID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAccountSettingsOptions : Instantiate GetAccountSettingsOptions
func (*IBMCloudShellV1) NewGetAccountSettingsOptions(accountID string) *GetAccountSettingsOptions {
	return &GetAccountSettingsOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetAccountSettingsOptions) SetAccountID(accountID string) *GetAccountSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountSettingsOptions) SetHeaders(param map[string]string) *GetAccountSettingsOptions {
	options.Headers = param
	return options
}

// RegionSetting : Describes a Cloud Shell region setting.
type RegionSetting struct {
	// State of the region.
	Enabled *bool `json:"enabled,omitempty"`

	// Name of the region.
	Key *string `json:"key,omitempty"`
}

// UnmarshalRegionSetting unmarshals an instance of RegionSetting from the specified map of raw messages.
func UnmarshalRegionSetting(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RegionSetting)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateAccountSettingsOptions : The UpdateAccountSettings options.
type UpdateAccountSettingsOptions struct {
	// The account ID in which the account settings belong to.
	AccountID *string `validate:"required,ne="`

	// Unique revision number for the settings object.
	Rev *string

	// You can choose which Cloud Shell features are available in the account and whether any new features are enabled as
	// they become available. The feature settings apply only to the enabled Cloud Shell locations.
	DefaultEnableNewFeatures *bool

	// Set whether Cloud Shell is enabled in a specific location for the account. The location determines where user and
	// session data are stored. By default, users are routed to the nearest available location.
	DefaultEnableNewRegions *bool

	// When enabled, Cloud Shell is available to all users in the account.
	Enabled *bool

	// List of Cloud Shell features.
	Features []Feature

	// List of Cloud Shell region settings.
	Regions []RegionSetting

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateAccountSettingsOptions : Instantiate UpdateAccountSettingsOptions
func (*IBMCloudShellV1) NewUpdateAccountSettingsOptions(accountID string) *UpdateAccountSettingsOptions {
	return &UpdateAccountSettingsOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateAccountSettingsOptions) SetAccountID(accountID string) *UpdateAccountSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *UpdateAccountSettingsOptions) SetRev(rev string) *UpdateAccountSettingsOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetDefaultEnableNewFeatures : Allow user to set DefaultEnableNewFeatures
func (_options *UpdateAccountSettingsOptions) SetDefaultEnableNewFeatures(defaultEnableNewFeatures bool) *UpdateAccountSettingsOptions {
	_options.DefaultEnableNewFeatures = core.BoolPtr(defaultEnableNewFeatures)
	return _options
}

// SetDefaultEnableNewRegions : Allow user to set DefaultEnableNewRegions
func (_options *UpdateAccountSettingsOptions) SetDefaultEnableNewRegions(defaultEnableNewRegions bool) *UpdateAccountSettingsOptions {
	_options.DefaultEnableNewRegions = core.BoolPtr(defaultEnableNewRegions)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateAccountSettingsOptions) SetEnabled(enabled bool) *UpdateAccountSettingsOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetFeatures : Allow user to set Features
func (_options *UpdateAccountSettingsOptions) SetFeatures(features []Feature) *UpdateAccountSettingsOptions {
	_options.Features = features
	return _options
}

// SetRegions : Allow user to set Regions
func (_options *UpdateAccountSettingsOptions) SetRegions(regions []RegionSetting) *UpdateAccountSettingsOptions {
	_options.Regions = regions
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccountSettingsOptions) SetHeaders(param map[string]string) *UpdateAccountSettingsOptions {
	options.Headers = param
	return options
}
