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
 * IBM OpenAPI SDK Code Generator Version: 3.92.0-af5c89a5-20240617-153232
 */

// Package configurationaggregatorv1 : Operations and models for the ConfigurationAggregatorV1 service
package configurationaggregatorv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	common "github.com/IBM/configuration-aggregator-go-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
)

// ConfigurationAggregatorV1 : Configuration Aggregator
//
// API Version: 1.0.0
// See: https://cloud.ibm.com/docs/app-configuration
type ConfigurationAggregatorV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.apprapp.cloud.ibm.com/apprapp/config_aggregator/v1/instances/provide-here-your-appconfig-instance-uuid"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "configuration_aggregator"

const ParameterizedServiceURL = "https://{region}.apprapp.cloud.ibm.com/apprapp/config_aggregator/v1/instances/{instance_id}"

var defaultUrlVariables = map[string]string{
	"region":      "us-south",
	"instance_id": "provide-here-your-appconfig-instance-uuid",
}

// ConfigurationAggregatorV1Options : Service options
type ConfigurationAggregatorV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewConfigurationAggregatorV1UsingExternalConfig : constructs an instance of ConfigurationAggregatorV1 with passed in options and external configuration.
func NewConfigurationAggregatorV1UsingExternalConfig(options *ConfigurationAggregatorV1Options) (configurationAggregator *ConfigurationAggregatorV1, err error) {
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

	configurationAggregator, err = NewConfigurationAggregatorV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = configurationAggregator.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = configurationAggregator.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewConfigurationAggregatorV1 : constructs an instance of ConfigurationAggregatorV1 with passed in options.
func NewConfigurationAggregatorV1(options *ConfigurationAggregatorV1Options) (service *ConfigurationAggregatorV1, err error) {
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

	service = &ConfigurationAggregatorV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "configurationAggregator" suitable for processing requests.
func (configurationAggregator *ConfigurationAggregatorV1) Clone() *ConfigurationAggregatorV1 {
	if core.IsNil(configurationAggregator) {
		return nil
	}
	clone := *configurationAggregator
	clone.Service = configurationAggregator.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (configurationAggregator *ConfigurationAggregatorV1) SetServiceURL(url string) error {
	err := configurationAggregator.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (configurationAggregator *ConfigurationAggregatorV1) GetServiceURL() string {
	return configurationAggregator.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (configurationAggregator *ConfigurationAggregatorV1) SetDefaultHeaders(headers http.Header) {
	configurationAggregator.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (configurationAggregator *ConfigurationAggregatorV1) SetEnableGzipCompression(enableGzip bool) {
	configurationAggregator.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (configurationAggregator *ConfigurationAggregatorV1) GetEnableGzipCompression() bool {
	return configurationAggregator.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (configurationAggregator *ConfigurationAggregatorV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	configurationAggregator.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (configurationAggregator *ConfigurationAggregatorV1) DisableRetries() {
	configurationAggregator.Service.DisableRetries()
}

// ListConfigs : Get the list of configurations of the resources
// Retrieve the list of resource configurations collected as part of Configuration Aggregator.
func (configurationAggregator *ConfigurationAggregatorV1) ListConfigs(listConfigsOptions *ListConfigsOptions) (result *ListConfigsResponse, response *core.DetailedResponse, err error) {
	result, response, err = configurationAggregator.ListConfigsWithContext(context.Background(), listConfigsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListConfigsWithContext is an alternate form of the ListConfigs method which supports a Context parameter
func (configurationAggregator *ConfigurationAggregatorV1) ListConfigsWithContext(ctx context.Context, listConfigsOptions *ListConfigsOptions) (result *ListConfigsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listConfigsOptions, "listConfigsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationAggregator.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationAggregator.Service.Options.URL, `/configs`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listConfigsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_aggregator", "V1", "ListConfigs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listConfigsOptions.ConfigType != nil {
		builder.AddQuery("config_type", fmt.Sprint(*listConfigsOptions.ConfigType))
	}
	if listConfigsOptions.ServiceName != nil {
		builder.AddQuery("service_name", fmt.Sprint(*listConfigsOptions.ServiceName))
	}
	if listConfigsOptions.ResourceGroupID != nil {
		builder.AddQuery("resource_group_id", fmt.Sprint(*listConfigsOptions.ResourceGroupID))
	}
	if listConfigsOptions.Location != nil {
		builder.AddQuery("location", fmt.Sprint(*listConfigsOptions.Location))
	}
	if listConfigsOptions.ResourceCrn != nil {
		builder.AddQuery("resource_crn", fmt.Sprint(*listConfigsOptions.ResourceCrn))
	}
	if listConfigsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listConfigsOptions.Limit))
	}
	if listConfigsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listConfigsOptions.Start))
	}
	if listConfigsOptions.SubAccount != nil {
		builder.AddQuery("sub_account", fmt.Sprint(*listConfigsOptions.SubAccount))
	}
	if listConfigsOptions.AccessTags != nil {
		builder.AddQuery("access_tags", fmt.Sprint(*listConfigsOptions.AccessTags))
	}
	if listConfigsOptions.UserTags != nil {
		builder.AddQuery("user_tags", fmt.Sprint(*listConfigsOptions.UserTags))
	}
	if listConfigsOptions.ServiceTags != nil {
		builder.AddQuery("service_tags", fmt.Sprint(*listConfigsOptions.ServiceTags))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationAggregator.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_configs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListConfigsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceSettings : Replace the settings for Configuration Aggregator
// Replace the settings for resource collection as part of the Configuration Aggregator feature.
func (configurationAggregator *ConfigurationAggregatorV1) ReplaceSettings(replaceSettingsOptions *ReplaceSettingsOptions) (result *SettingsResponse, response *core.DetailedResponse, err error) {
	result, response, err = configurationAggregator.ReplaceSettingsWithContext(context.Background(), replaceSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceSettingsWithContext is an alternate form of the ReplaceSettings method which supports a Context parameter
func (configurationAggregator *ConfigurationAggregatorV1) ReplaceSettingsWithContext(ctx context.Context, replaceSettingsOptions *ReplaceSettingsOptions) (result *SettingsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceSettingsOptions, "replaceSettingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceSettingsOptions, "replaceSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationAggregator.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationAggregator.Service.Options.URL, `/settings`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_aggregator", "V1", "ReplaceSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceSettingsOptions.ResourceCollectionEnabled != nil {
		body["resource_collection_enabled"] = replaceSettingsOptions.ResourceCollectionEnabled
	}
	if replaceSettingsOptions.TrustedProfileID != nil {
		body["trusted_profile_id"] = replaceSettingsOptions.TrustedProfileID
	}
	if replaceSettingsOptions.Regions != nil {
		body["regions"] = replaceSettingsOptions.Regions
	}
	if replaceSettingsOptions.AdditionalScope != nil {
		body["additional_scope"] = replaceSettingsOptions.AdditionalScope
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
	response, err = configurationAggregator.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_settings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSettingsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetSettings : Retrieve the settings for Configuration Aggregator feature
// Retrieve settings for resource collection in Configuration Aggregator.
func (configurationAggregator *ConfigurationAggregatorV1) GetSettings(getSettingsOptions *GetSettingsOptions) (result *SettingsResponse, response *core.DetailedResponse, err error) {
	result, response, err = configurationAggregator.GetSettingsWithContext(context.Background(), getSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetSettingsWithContext is an alternate form of the GetSettings method which supports a Context parameter
func (configurationAggregator *ConfigurationAggregatorV1) GetSettingsWithContext(ctx context.Context, getSettingsOptions *GetSettingsOptions) (result *SettingsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSettingsOptions, "getSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationAggregator.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationAggregator.Service.Options.URL, `/settings`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_aggregator", "V1", "GetSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationAggregator.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_settings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSettingsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetResourceCollectionStatus : Retrieve status for resource collection in Configuration Aggregator
// Retrieve the status of the resource collection as part of Configuration Aggregator.
func (configurationAggregator *ConfigurationAggregatorV1) GetResourceCollectionStatus(getResourceCollectionStatusOptions *GetResourceCollectionStatusOptions) (result *StatusResponse, response *core.DetailedResponse, err error) {
	result, response, err = configurationAggregator.GetResourceCollectionStatusWithContext(context.Background(), getResourceCollectionStatusOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetResourceCollectionStatusWithContext is an alternate form of the GetResourceCollectionStatus method which supports a Context parameter
func (configurationAggregator *ConfigurationAggregatorV1) GetResourceCollectionStatusWithContext(ctx context.Context, getResourceCollectionStatusOptions *GetResourceCollectionStatusOptions) (result *StatusResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getResourceCollectionStatusOptions, "getResourceCollectionStatusOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationAggregator.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationAggregator.Service.Options.URL, `/resource_collection_status`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getResourceCollectionStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_aggregator", "V1", "GetResourceCollectionStatus")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationAggregator.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_resource_collection_status", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStatusResponse)
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

// About : The basic metadata fetched from the query API.
type About struct {
	// The account ID in which the resource exists.
	AccountID *string `json:"account_id" validate:"required"`

	// The type of configuration of the retrieved resource.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique CRN of the IBM Cloud resource.
	ResourceCrn *string `json:"resource_crn" validate:"required"`

	// The account ID.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// The name of the service to which the resources belongs.
	ServiceName *string `json:"service_name" validate:"required"`

	// User defined name of the resource.
	ResourceName *string `json:"resource_name" validate:"required"`

	// Date/time stamp identifying when the information was last collected. Must be in the RFC 3339 format.
	LastConfigRefreshTime *strfmt.DateTime `json:"last_config_refresh_time" validate:"required"`

	// Location of the resource specified.
	Location *string `json:"location" validate:"required"`

	// Access tags specified by the user for the resource. For more information, see
	// https://cloud.ibm.com/docs/account?topic=account-tag&interface=ui#tag-types.
	AccessTags []string `json:"access_tags,omitempty"`

	// User tags specified by the user for the resource. For more information, see
	// https://cloud.ibm.com/docs/account?topic=account-tag&interface=ui#tag-types.
	UserTags []string `json:"user_tags,omitempty"`

	// Tags attached to resources or service IDs by an authorized user in the account. For more information, see
	// https://cloud.ibm.com/docs/account?topic=account-tag&interface=ui#tag-types.
	ServiceTags []string `json:"service_tags,omitempty"`
}

// UnmarshalAbout unmarshals an instance of About from the specified map of raw messages.
func UnmarshalAbout(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(About)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		err = core.SDKErrorf(err, "", "config_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crn", &obj.ResourceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_config_refresh_time", &obj.LastConfigRefreshTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_config_refresh_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "access_tags", &obj.AccessTags)
	if err != nil {
		err = core.SDKErrorf(err, "", "access_tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_tags", &obj.UserTags)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_tags", &obj.ServiceTags)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_tags-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalScope : The additional scope that enables resource collection for Enterprise acccounts.
type AdditionalScope struct {
	// The type of scope. Currently allowed value is Enterprise.
	Type *string `json:"type,omitempty"`

	// The Enterprise ID.
	EnterpriseID *string `json:"enterprise_id,omitempty"`

	// The Profile Template details applied on the enterprise account.
	ProfileTemplate *ProfileTemplate `json:"profile_template,omitempty"`
}

// UnmarshalAdditionalScope unmarshals an instance of AdditionalScope from the specified map of raw messages.
func UnmarshalAdditionalScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalScope)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_id", &obj.EnterpriseID)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "profile_template", &obj.ProfileTemplate, UnmarshalProfileTemplate)
	if err != nil {
		err = core.SDKErrorf(err, "", "profile_template-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Config : Configuration of each individual resource.
type Config struct {
	// The basic metadata fetched from the query API.
	About *About `json:"about" validate:"required"`

	// The configuration of the resource.
	Config *Configuration `json:"config" validate:"required"`
}

// UnmarshalConfig unmarshals an instance of Config from the specified map of raw messages.
func UnmarshalConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Config)
	err = core.UnmarshalModel(m, "about", &obj.About, UnmarshalAbout)
	if err != nil {
		err = core.SDKErrorf(err, "", "about-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalConfiguration)
	if err != nil {
		err = core.SDKErrorf(err, "", "config-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Configuration : The configuration of the resource.
type Configuration struct {

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of Configuration
func (o *Configuration) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of Configuration
func (o *Configuration) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of Configuration
func (o *Configuration) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of Configuration
func (o *Configuration) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of Configuration
func (o *Configuration) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalConfiguration unmarshals an instance of Configuration from the specified map of raw messages.
func UnmarshalConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Configuration)
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetResourceCollectionStatusOptions : The GetResourceCollectionStatus options.
type GetResourceCollectionStatusOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetResourceCollectionStatusOptions : Instantiate GetResourceCollectionStatusOptions
func (*ConfigurationAggregatorV1) NewGetResourceCollectionStatusOptions() *GetResourceCollectionStatusOptions {
	return &GetResourceCollectionStatusOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceCollectionStatusOptions) SetHeaders(param map[string]string) *GetResourceCollectionStatusOptions {
	options.Headers = param
	return options
}

// GetSettingsOptions : The GetSettings options.
type GetSettingsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSettingsOptions : Instantiate GetSettingsOptions
func (*ConfigurationAggregatorV1) NewGetSettingsOptions() *GetSettingsOptions {
	return &GetSettingsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetSettingsOptions) SetHeaders(param map[string]string) *GetSettingsOptions {
	options.Headers = param
	return options
}

// ListConfigsOptions : The ListConfigs options.
type ListConfigsOptions struct {
	// The type of resource configuration that are to be retrieved.
	ConfigType *string `json:"config_type,omitempty"`

	// The name of the IBM Cloud service for which resources are to be retrieved.
	ServiceName *string `json:"service_name,omitempty"`

	// The resource group id of the resources.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// The location or region in which the resources are created.
	Location *string `json:"location,omitempty"`

	// The crn of the resource.
	ResourceCrn *string `json:"resource_crn,omitempty"`

	// The number of resources for which the configuration can be fetched.
	Limit *int64 `json:"limit,omitempty"`

	// The start string to fetch the resource.
	Start *string `json:"start,omitempty"`

	// Filter the resource configurations from the specified sub-account in an enterprise hierarchy.
	SubAccount *string `json:"sub_account,omitempty"`

	// Filter the resource configurations attached with the specified access tags.
	AccessTags *string `json:"access_tags,omitempty"`

	// Filter the resource configurations attached with the specified user tags.
	UserTags *string `json:"user_tags,omitempty"`

	// Filter the resource configurations attached with the specified service tags.
	ServiceTags *string `json:"service_tags,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListConfigsOptions : Instantiate ListConfigsOptions
func (*ConfigurationAggregatorV1) NewListConfigsOptions() *ListConfigsOptions {
	return &ListConfigsOptions{}
}

// SetConfigType : Allow user to set ConfigType
func (_options *ListConfigsOptions) SetConfigType(configType string) *ListConfigsOptions {
	_options.ConfigType = core.StringPtr(configType)
	return _options
}

// SetServiceName : Allow user to set ServiceName
func (_options *ListConfigsOptions) SetServiceName(serviceName string) *ListConfigsOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *ListConfigsOptions) SetResourceGroupID(resourceGroupID string) *ListConfigsOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *ListConfigsOptions) SetLocation(location string) *ListConfigsOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetResourceCrn : Allow user to set ResourceCrn
func (_options *ListConfigsOptions) SetResourceCrn(resourceCrn string) *ListConfigsOptions {
	_options.ResourceCrn = core.StringPtr(resourceCrn)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListConfigsOptions) SetLimit(limit int64) *ListConfigsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListConfigsOptions) SetStart(start string) *ListConfigsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetSubAccount : Allow user to set SubAccount
func (_options *ListConfigsOptions) SetSubAccount(subAccount string) *ListConfigsOptions {
	_options.SubAccount = core.StringPtr(subAccount)
	return _options
}

// SetAccessTags : Allow user to set AccessTags
func (_options *ListConfigsOptions) SetAccessTags(accessTags string) *ListConfigsOptions {
	_options.AccessTags = core.StringPtr(accessTags)
	return _options
}

// SetUserTags : Allow user to set UserTags
func (_options *ListConfigsOptions) SetUserTags(userTags string) *ListConfigsOptions {
	_options.UserTags = core.StringPtr(userTags)
	return _options
}

// SetServiceTags : Allow user to set ServiceTags
func (_options *ListConfigsOptions) SetServiceTags(serviceTags string) *ListConfigsOptions {
	_options.ServiceTags = core.StringPtr(serviceTags)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListConfigsOptions) SetHeaders(param map[string]string) *ListConfigsOptions {
	options.Headers = param
	return options
}

// ListConfigsResponse : List configs api response.
type ListConfigsResponse struct {
	// The total number of resources present.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The maximum number of resources retrieved per page.
	Limit *int64 `json:"limit,omitempty"`

	// The reference to the first page of entries.
	First *PaginatedFirst `json:"first,omitempty"`

	// The reference to the previous page of entries.
	Prev *PaginatedPrevious `json:"prev,omitempty"`

	// The reference to the next page of entries.
	Next *PaginatedNext `json:"next,omitempty"`

	// Array of resource configurations.
	Configs []Config `json:"configs,omitempty"`
}

// UnmarshalListConfigsResponse unmarshals an instance of ListConfigsResponse from the specified map of raw messages.
func UnmarshalListConfigsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListConfigsResponse)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "prev", &obj.Prev, UnmarshalPaginatedPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "prev-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "configs", &obj.Configs, UnmarshalConfig)
	if err != nil {
		err = core.SDKErrorf(err, "", "configs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListConfigsResponse) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// PaginatedFirst : The reference to the first page of entries.
type PaginatedFirst struct {
	// The reference to the first page of entries.
	Href *string `json:"href,omitempty"`
}

// UnmarshalPaginatedFirst unmarshals an instance of PaginatedFirst from the specified map of raw messages.
func UnmarshalPaginatedFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedNext : The reference to the next page of entries.
type PaginatedNext struct {
	// The reference to the next page of entries.
	Href *string `json:"href,omitempty"`

	// the start string for the query to view the page.
	Start *string `json:"start,omitempty"`
}

// UnmarshalPaginatedNext unmarshals an instance of PaginatedNext from the specified map of raw messages.
func UnmarshalPaginatedNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedPrevious : The reference to the previous page of entries.
type PaginatedPrevious struct {
	// The reference to the previous page of entries.
	Href *string `json:"href,omitempty"`

	// the start string for the query to view the page.
	Start *string `json:"start,omitempty"`
}

// UnmarshalPaginatedPrevious unmarshals an instance of PaginatedPrevious from the specified map of raw messages.
func UnmarshalPaginatedPrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedPrevious)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileTemplate : The Profile Template details applied on the enterprise account.
type ProfileTemplate struct {
	// The Profile Template ID created in the enterprise account that provides access to App Configuration instance for
	// resource collection.
	ID *string `json:"id,omitempty"`

	// The trusted profile ID that provides access to App Configuration instance to retrieve template information.
	TrustedProfileID *string `json:"trusted_profile_id,omitempty"`
}

// UnmarshalProfileTemplate unmarshals an instance of ProfileTemplate from the specified map of raw messages.
func UnmarshalProfileTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileTemplate)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "trusted_profile_id", &obj.TrustedProfileID)
	if err != nil {
		err = core.SDKErrorf(err, "", "trusted_profile_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceSettingsOptions : The ReplaceSettings options.
type ReplaceSettingsOptions struct {
	// The field denoting if the resource collection is enabled.
	ResourceCollectionEnabled *bool `json:"resource_collection_enabled,omitempty"`

	// The trusted profile id that provides Reader access to the App Configuration instance to collect resource metadata.
	TrustedProfileID *string `json:"trusted_profile_id,omitempty"`

	// The list of regions across which the resource collection is enabled.
	Regions []string `json:"regions,omitempty"`

	// The additional scope that enables resource collection for Enterprise acccounts.
	AdditionalScope []AdditionalScope `json:"additional_scope,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceSettingsOptions : Instantiate ReplaceSettingsOptions
func (*ConfigurationAggregatorV1) NewReplaceSettingsOptions() *ReplaceSettingsOptions {
	return &ReplaceSettingsOptions{}
}

// SetResourceCollectionEnabled : Allow user to set ResourceCollectionEnabled
func (_options *ReplaceSettingsOptions) SetResourceCollectionEnabled(resourceCollectionEnabled bool) *ReplaceSettingsOptions {
	_options.ResourceCollectionEnabled = core.BoolPtr(resourceCollectionEnabled)
	return _options
}

// SetTrustedProfileID : Allow user to set TrustedProfileID
func (_options *ReplaceSettingsOptions) SetTrustedProfileID(trustedProfileID string) *ReplaceSettingsOptions {
	_options.TrustedProfileID = core.StringPtr(trustedProfileID)
	return _options
}

// SetRegions : Allow user to set Regions
func (_options *ReplaceSettingsOptions) SetRegions(regions []string) *ReplaceSettingsOptions {
	_options.Regions = regions
	return _options
}

// SetAdditionalScope : Allow user to set AdditionalScope
func (_options *ReplaceSettingsOptions) SetAdditionalScope(additionalScope []AdditionalScope) *ReplaceSettingsOptions {
	_options.AdditionalScope = additionalScope
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceSettingsOptions) SetHeaders(param map[string]string) *ReplaceSettingsOptions {
	options.Headers = param
	return options
}

// SettingsResponse : Settings API response.
type SettingsResponse struct {
	// The field to check if the resource collection is enabled.
	ResourceCollectionEnabled *bool `json:"resource_collection_enabled,omitempty"`

	// The trusted profile ID that provides access to App Configuration instance to retrieve resource metadata.
	TrustedProfileID *string `json:"trusted_profile_id,omitempty"`

	// The last time the settings was last updated.
	LastUpdated *strfmt.DateTime `json:"last_updated,omitempty"`

	// Regions for which the resource collection is enabled.
	Regions []string `json:"regions,omitempty"`

	// The additional scope that enables resource collection for Enterprise acccounts.
	AdditionalScope []AdditionalScope `json:"additional_scope,omitempty"`
}

// UnmarshalSettingsResponse unmarshals an instance of SettingsResponse from the specified map of raw messages.
func UnmarshalSettingsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SettingsResponse)
	err = core.UnmarshalPrimitive(m, "resource_collection_enabled", &obj.ResourceCollectionEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_collection_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "trusted_profile_id", &obj.TrustedProfileID)
	if err != nil {
		err = core.SDKErrorf(err, "", "trusted_profile_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_updated", &obj.LastUpdated)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_updated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "regions", &obj.Regions)
	if err != nil {
		err = core.SDKErrorf(err, "", "regions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "additional_scope", &obj.AdditionalScope, UnmarshalAdditionalScope)
	if err != nil {
		err = core.SDKErrorf(err, "", "additional_scope-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StatusResponse : The Status response.
type StatusResponse struct {
	// The timestamp at which the configuration was last refreshed.
	LastConfigRefreshTime *strfmt.DateTime `json:"last_config_refresh_time,omitempty"`

	// Status of the resource collection.
	Status *string `json:"status,omitempty"`
}

// Constants associated with the StatusResponse.Status property.
// Status of the resource collection.
const (
	StatusResponse_Status_Complete   = "complete"
	StatusResponse_Status_Initiated  = "initiated"
	StatusResponse_Status_Inprogress = "inprogress"
)

// UnmarshalStatusResponse unmarshals an instance of StatusResponse from the specified map of raw messages.
func UnmarshalStatusResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StatusResponse)
	err = core.UnmarshalPrimitive(m, "last_config_refresh_time", &obj.LastConfigRefreshTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_config_refresh_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigsPager can be used to simplify the use of the "ListConfigs" method.
type ConfigsPager struct {
	hasNext     bool
	options     *ListConfigsOptions
	client      *ConfigurationAggregatorV1
	pageContext struct {
		next *string
	}
}

// NewConfigsPager returns a new ConfigsPager instance.
func (configurationAggregator *ConfigurationAggregatorV1) NewConfigsPager(options *ListConfigsOptions) (pager *ConfigsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListConfigsOptions = *options
	pager = &ConfigsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  configurationAggregator,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ConfigsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ConfigsPager) GetNextWithContext(ctx context.Context) (page []Config, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListConfigsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Configs

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ConfigsPager) GetAllWithContext(ctx context.Context) (allItems []Config, err error) {
	for pager.HasNext() {
		var nextPage []Config
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
func (pager *ConfigsPager) GetNext() (page []Config, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ConfigsPager) GetAll() (allItems []Config, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
