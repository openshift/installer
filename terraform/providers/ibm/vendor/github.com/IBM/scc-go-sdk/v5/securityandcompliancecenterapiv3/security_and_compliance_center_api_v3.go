/**
 * (C) Copyright IBM Corp. 2023.
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
 * IBM OpenAPI SDK Code Generator Version: 3.75.0-726bc7e3-20230713-221716
 */

// Package securityandcompliancecenterapiv3 : Operations and models for the SecurityAndComplianceCenterApiV3 service
package securityandcompliancecenterapiv3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/scc-go-sdk/v5/common"
	"github.com/go-openapi/strfmt"
)

// SecurityAndComplianceCenterApiV3 : Security and Compliance Center API
//
// API Version: 3.0.0
type SecurityAndComplianceCenterApiV3 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.compliance.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "security_and_compliance_center_api"

const ParameterizedServiceURL = "https://{region}.compliance.cloud.ibm.com"

var defaultUrlVariables = map[string]string{
	"region": "us-south",
}

// SecurityAndComplianceCenterApiV3Options : Service options
type SecurityAndComplianceCenterApiV3Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewSecurityAndComplianceCenterApiV3UsingExternalConfig : constructs an instance of SecurityAndComplianceCenterApiV3 with passed in options and external configuration.
func NewSecurityAndComplianceCenterApiV3UsingExternalConfig(options *SecurityAndComplianceCenterApiV3Options) (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	securityAndComplianceCenterApi, err = NewSecurityAndComplianceCenterApiV3(options)
	if err != nil {
		return
	}

	err = securityAndComplianceCenterApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = securityAndComplianceCenterApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewSecurityAndComplianceCenterApiV3 : constructs an instance of SecurityAndComplianceCenterApiV3 with passed in options.
func NewSecurityAndComplianceCenterApiV3(options *SecurityAndComplianceCenterApiV3Options) (service *SecurityAndComplianceCenterApiV3, err error) {
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

	service = &SecurityAndComplianceCenterApiV3{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	endpoints := map[string]string{
		"au-syd":   "https://au-syd.compliance.cloud.ibm.com",   // The API endpoint in the au-syd region
		"ca-tor":   "https://ca-tor.compliance.cloud.ibm.com",   // The API endpoint in the ca-tor region
		"eu-de":    "https://eu-de.compliance.cloud.ibm.com",    // The API endpoint in the eu-de region
		"eu-fr2":   "https://eu-fr2.compliance.cloud.ibm.com",   // The API endpoint in the eu-fr2 region
		"us-south": "https://us-south.compliance.cloud.ibm.com", // The API endpoint in the us-south region
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}

// Clone makes a copy of "securityAndComplianceCenterApi" suitable for processing requests.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) Clone() *SecurityAndComplianceCenterApiV3 {
	if core.IsNil(securityAndComplianceCenterApi) {
		return nil
	}
	clone := *securityAndComplianceCenterApi
	clone.Service = securityAndComplianceCenterApi.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) SetServiceURL(url string) error {
	return securityAndComplianceCenterApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetServiceURL() string {
	return securityAndComplianceCenterApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) SetDefaultHeaders(headers http.Header) {
	securityAndComplianceCenterApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) SetEnableGzipCompression(enableGzip bool) {
	securityAndComplianceCenterApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetEnableGzipCompression() bool {
	return securityAndComplianceCenterApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	securityAndComplianceCenterApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DisableRetries() {
	securityAndComplianceCenterApi.Service.DisableRetries()
}

// GetSettings : Get settings
// Retrieve the settings of your service instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetSettings(getSettingsOptions *GetSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetSettingsWithContext(context.Background(), getSettingsOptions)
}

// GetSettingsWithContext is an alternate form of the GetSettings method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetSettingsWithContext(ctx context.Context, getSettingsOptions *GetSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSettingsOptions, "getSettingsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getSettingsOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/settings`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getSettingsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*getSettingsOptions.XCorrelationID))
	}
	if getSettingsOptions.XRequestID != nil {
		builder.AddHeader("X-Request-Id", fmt.Sprint(*getSettingsOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSettings : Update settings
// Update the settings of your service instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) UpdateSettings(updateSettingsOptions *UpdateSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.UpdateSettingsWithContext(context.Background(), updateSettingsOptions)
}

// UpdateSettingsWithContext is an alternate form of the UpdateSettings method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) UpdateSettingsWithContext(ctx context.Context, updateSettingsOptions *UpdateSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSettingsOptions, "updateSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSettingsOptions, "updateSettingsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *updateSettingsOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/settings`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "UpdateSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateSettingsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*updateSettingsOptions.XCorrelationID))
	}
	if updateSettingsOptions.XRequestID != nil {
		builder.AddHeader("X-Request-Id", fmt.Sprint(*updateSettingsOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if updateSettingsOptions.EventNotifications != nil {
		body["event_notifications"] = updateSettingsOptions.EventNotifications
	}
	if updateSettingsOptions.ObjectStorage != nil {
		body["object_storage"] = updateSettingsOptions.ObjectStorage
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostTestEvent : Create a test event
// Send a test event to your Event Notifications instance to ensure that the events that are generated  by Security and
// Compliance Center are being forwarded to Event Notifications. For more information, see [Enabling event
// notifications](/docs/security-compliance?topic=security-compliance-event-notifications#event-notifications-test-api).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) PostTestEvent(postTestEventOptions *PostTestEventOptions) (result *TestEvent, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.PostTestEventWithContext(context.Background(), postTestEventOptions)
}

// PostTestEventWithContext is an alternate form of the PostTestEvent method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) PostTestEventWithContext(ctx context.Context, postTestEventOptions *PostTestEventOptions) (result *TestEvent, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(postTestEventOptions, "postTestEventOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *postTestEventOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/test_event`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range postTestEventOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "PostTestEvent")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if postTestEventOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*postTestEventOptions.XCorrelationID))
	}
	if postTestEventOptions.XRequestID != nil {
		builder.AddHeader("X-Request-Id", fmt.Sprint(*postTestEventOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
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

// ListControlLibraries : Get control libraries
// Retrieve all of the control libraries that are available in your account, including predefined, and custom libraries.
//
// With Security and Compliance Center, you can create a custom control library that is specific to your organization's
// needs.  You define the controls and specifications before you map previously created assessments. Each control has
// several specifications  and assessments that are mapped to it. A specification is a defined requirement that is
// specific to a component. An assessment, or several,  are mapped to each specification with a detailed evaluation that
// is done to check whether the specification is compliant. For more information, see [Creating custom
// libraries](/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListControlLibraries(listControlLibrariesOptions *ListControlLibrariesOptions) (result *ControlLibraryCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListControlLibrariesWithContext(context.Background(), listControlLibrariesOptions)
}

// ListControlLibrariesWithContext is an alternate form of the ListControlLibraries method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListControlLibrariesWithContext(ctx context.Context, listControlLibrariesOptions *ListControlLibrariesOptions) (result *ControlLibraryCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listControlLibrariesOptions, "listControlLibrariesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listControlLibrariesOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/control_libraries`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listControlLibrariesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListControlLibraries")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listControlLibrariesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listControlLibrariesOptions.XCorrelationID))
	}
	if listControlLibrariesOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listControlLibrariesOptions.XRequestID))
	}

	if listControlLibrariesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listControlLibrariesOptions.Limit))
	}
	if listControlLibrariesOptions.ControlLibraryType != nil {
		builder.AddQuery("control_library_type", fmt.Sprint(*listControlLibrariesOptions.ControlLibraryType))
	}
	if listControlLibrariesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listControlLibrariesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibraryCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateCustomControlLibrary : Create a custom control library
// Create a custom control library that is specific to your organization's needs.
//
// With Security and Compliance Center, you can create a custom control library that is specific to your organization's
// needs.  You define the controls and specifications before you map previously created assessments. Each control has
// several specifications  and assessments that are mapped to it. A specification is a defined requirement that is
// specific to a component. An assessment, or several,  are mapped to each specification with a detailed evaluation that
// is done to check whether the specification is compliant. For more information, see [Creating custom
// libraries](/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateCustomControlLibrary(createCustomControlLibraryOptions *CreateCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.CreateCustomControlLibraryWithContext(context.Background(), createCustomControlLibraryOptions)
}

// CreateCustomControlLibraryWithContext is an alternate form of the CreateCustomControlLibrary method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateCustomControlLibraryWithContext(ctx context.Context, createCustomControlLibraryOptions *CreateCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCustomControlLibraryOptions, "createCustomControlLibraryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createCustomControlLibraryOptions, "createCustomControlLibraryOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *createCustomControlLibraryOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/control_libraries`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createCustomControlLibraryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "CreateCustomControlLibrary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createCustomControlLibraryOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createCustomControlLibraryOptions.XCorrelationID))
	}
	if createCustomControlLibraryOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*createCustomControlLibraryOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if createCustomControlLibraryOptions.ControlLibraryName != nil {
		body["control_library_name"] = createCustomControlLibraryOptions.ControlLibraryName
	}
	if createCustomControlLibraryOptions.ControlLibraryDescription != nil {
		body["control_library_description"] = createCustomControlLibraryOptions.ControlLibraryDescription
	}
	if createCustomControlLibraryOptions.ControlLibraryType != nil {
		body["control_library_type"] = createCustomControlLibraryOptions.ControlLibraryType
	}
	if createCustomControlLibraryOptions.Controls != nil {
		body["controls"] = createCustomControlLibraryOptions.Controls
	}
	if createCustomControlLibraryOptions.VersionGroupLabel != nil {
		body["version_group_label"] = createCustomControlLibraryOptions.VersionGroupLabel
	}
	if createCustomControlLibraryOptions.ControlLibraryVersion != nil {
		body["control_library_version"] = createCustomControlLibraryOptions.ControlLibraryVersion
	}
	if createCustomControlLibraryOptions.Latest != nil {
		body["latest"] = createCustomControlLibraryOptions.Latest
	}
	if createCustomControlLibraryOptions.ControlsCount != nil {
		body["controls_count"] = createCustomControlLibraryOptions.ControlsCount
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibrary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCustomControlLibrary : Delete a control library
// Delete a custom control library by providing the control library ID.  You can find this ID by looking in the Security
// and Compliance Center UI.
//
// With Security and Compliance Center, you can manage a custom control library  that is specific to your organization's
// needs. Each control has several specifications  and assessments that are mapped to it.  For more information, see
// [Creating custom libraries](/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteCustomControlLibrary(deleteCustomControlLibraryOptions *DeleteCustomControlLibraryOptions) (result *ControlLibraryDelete, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.DeleteCustomControlLibraryWithContext(context.Background(), deleteCustomControlLibraryOptions)
}

// DeleteCustomControlLibraryWithContext is an alternate form of the DeleteCustomControlLibrary method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteCustomControlLibraryWithContext(ctx context.Context, deleteCustomControlLibraryOptions *DeleteCustomControlLibraryOptions) (result *ControlLibraryDelete, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCustomControlLibraryOptions, "deleteCustomControlLibraryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCustomControlLibraryOptions, "deleteCustomControlLibraryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"control_libraries_id": *deleteCustomControlLibraryOptions.ControlLibrariesID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *deleteCustomControlLibraryOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/control_libraries/{control_libraries_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCustomControlLibraryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "DeleteCustomControlLibrary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteCustomControlLibraryOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteCustomControlLibraryOptions.XCorrelationID))
	}
	if deleteCustomControlLibraryOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*deleteCustomControlLibraryOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibraryDelete)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetControlLibrary : Get a control library
// View the details of a control library by specifying its ID.
//
// With Security and Compliance Center, you can create a custom control library that is specific to your organization's
// needs.  You define the controls and specifications before you map previously created assessments. Each control has
// several specifications  and assessments that are mapped to it. A specification is a defined requirement that is
// specific to a component. An assessment, or several,  are mapped to each specification with a detailed evaluation that
// is done to check whether the specification is compliant. For more information, see [Creating custom
// libraries](/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetControlLibrary(getControlLibraryOptions *GetControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetControlLibraryWithContext(context.Background(), getControlLibraryOptions)
}

// GetControlLibraryWithContext is an alternate form of the GetControlLibrary method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetControlLibraryWithContext(ctx context.Context, getControlLibraryOptions *GetControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getControlLibraryOptions, "getControlLibraryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getControlLibraryOptions, "getControlLibraryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"control_libraries_id": *getControlLibraryOptions.ControlLibrariesID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getControlLibraryOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/control_libraries/{control_libraries_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getControlLibraryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetControlLibrary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getControlLibraryOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getControlLibraryOptions.XCorrelationID))
	}
	if getControlLibraryOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getControlLibraryOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibrary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceCustomControlLibrary : Update a control library
// Update a custom control library by providing the control library ID. You can find this ID in the Security and
// Compliance Center UI.
//
// With Security and Compliance Center, you can create and update a custom control library that is specific to your
// organization's needs.  You define the controls and specifications before you map previously created assessments. Each
// control has several specifications  and assessments that are mapped to it. For more information, see [Creating custom
// libraries](/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ReplaceCustomControlLibrary(replaceCustomControlLibraryOptions *ReplaceCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ReplaceCustomControlLibraryWithContext(context.Background(), replaceCustomControlLibraryOptions)
}

// ReplaceCustomControlLibraryWithContext is an alternate form of the ReplaceCustomControlLibrary method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ReplaceCustomControlLibraryWithContext(ctx context.Context, replaceCustomControlLibraryOptions *ReplaceCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceCustomControlLibraryOptions, "replaceCustomControlLibraryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceCustomControlLibraryOptions, "replaceCustomControlLibraryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"control_libraries_id": *replaceCustomControlLibraryOptions.ControlLibrariesID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *replaceCustomControlLibraryOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/control_libraries/{control_libraries_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceCustomControlLibraryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ReplaceCustomControlLibrary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceCustomControlLibraryOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*replaceCustomControlLibraryOptions.XCorrelationID))
	}
	if replaceCustomControlLibraryOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*replaceCustomControlLibraryOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if replaceCustomControlLibraryOptions.ID != nil {
		body["id"] = replaceCustomControlLibraryOptions.ID
	}
	if replaceCustomControlLibraryOptions.AccountID != nil {
		body["account_id"] = replaceCustomControlLibraryOptions.AccountID
	}
	if replaceCustomControlLibraryOptions.ControlLibraryName != nil {
		body["control_library_name"] = replaceCustomControlLibraryOptions.ControlLibraryName
	}
	if replaceCustomControlLibraryOptions.ControlLibraryDescription != nil {
		body["control_library_description"] = replaceCustomControlLibraryOptions.ControlLibraryDescription
	}
	if replaceCustomControlLibraryOptions.ControlLibraryType != nil {
		body["control_library_type"] = replaceCustomControlLibraryOptions.ControlLibraryType
	}
	if replaceCustomControlLibraryOptions.VersionGroupLabel != nil {
		body["version_group_label"] = replaceCustomControlLibraryOptions.VersionGroupLabel
	}
	if replaceCustomControlLibraryOptions.ControlLibraryVersion != nil {
		body["control_library_version"] = replaceCustomControlLibraryOptions.ControlLibraryVersion
	}
	if replaceCustomControlLibraryOptions.CreatedOn != nil {
		body["created_on"] = replaceCustomControlLibraryOptions.CreatedOn
	}
	if replaceCustomControlLibraryOptions.CreatedBy != nil {
		body["created_by"] = replaceCustomControlLibraryOptions.CreatedBy
	}
	if replaceCustomControlLibraryOptions.UpdatedOn != nil {
		body["updated_on"] = replaceCustomControlLibraryOptions.UpdatedOn
	}
	if replaceCustomControlLibraryOptions.UpdatedBy != nil {
		body["updated_by"] = replaceCustomControlLibraryOptions.UpdatedBy
	}
	if replaceCustomControlLibraryOptions.Latest != nil {
		body["latest"] = replaceCustomControlLibraryOptions.Latest
	}
	if replaceCustomControlLibraryOptions.HierarchyEnabled != nil {
		body["hierarchy_enabled"] = replaceCustomControlLibraryOptions.HierarchyEnabled
	}
	if replaceCustomControlLibraryOptions.ControlsCount != nil {
		body["controls_count"] = replaceCustomControlLibraryOptions.ControlsCount
	}
	if replaceCustomControlLibraryOptions.ControlParentsCount != nil {
		body["control_parents_count"] = replaceCustomControlLibraryOptions.ControlParentsCount
	}
	if replaceCustomControlLibraryOptions.Controls != nil {
		body["controls"] = replaceCustomControlLibraryOptions.Controls
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibrary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProfiles : Get all profiles
// View all of the predefined and custom profiles that are available in your account.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListProfiles(listProfilesOptions *ListProfilesOptions) (result *ProfileCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListProfilesWithContext(context.Background(), listProfilesOptions)
}

// ListProfilesWithContext is an alternate form of the ListProfiles method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListProfilesWithContext(ctx context.Context, listProfilesOptions *ListProfilesOptions) (result *ProfileCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProfilesOptions, "listProfilesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listProfilesOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProfilesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListProfiles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listProfilesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listProfilesOptions.XCorrelationID))
	}
	if listProfilesOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listProfilesOptions.XRequestID))
	}

	if listProfilesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProfilesOptions.Limit))
	}
	if listProfilesOptions.ProfileType != nil {
		builder.AddQuery("profile_type", fmt.Sprint(*listProfilesOptions.ProfileType))
	}
	if listProfilesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listProfilesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateProfile : Create a custom profile
// Create a custom profile that is specific to your usecase, by using an existing library as a starting point.  For more
// information, see [Building custom
// profiles](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateProfile(createProfileOptions *CreateProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.CreateProfileWithContext(context.Background(), createProfileOptions)
}

// CreateProfileWithContext is an alternate form of the CreateProfile method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateProfileWithContext(ctx context.Context, createProfileOptions *CreateProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProfileOptions, "createProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createProfileOptions, "createProfileOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *createProfileOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles`, nil)
	if err != nil {
		return
	}
	for headerName, headerValue := range createProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "CreateProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createProfileOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createProfileOptions.XCorrelationID))
	}
	if createProfileOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*createProfileOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if createProfileOptions.ProfileName != nil {
		body["profile_name"] = createProfileOptions.ProfileName
	}
	if createProfileOptions.ProfileDescription != nil {
		body["profile_description"] = createProfileOptions.ProfileDescription
	}
	if createProfileOptions.ProfileType != nil {
		body["profile_type"] = createProfileOptions.ProfileType
	}
	if createProfileOptions.Controls != nil {
		body["controls"] = createProfileOptions.Controls
	}
	if createProfileOptions.DefaultParameters != nil {
		body["default_parameters"] = createProfileOptions.DefaultParameters
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCustomProfile : Delete a custom profile
// Delete a custom profile by providing the profile ID.  You can find the ID in the Security and Compliance Center UI.
// For more information about managing your custom profiles, see [Building custom
// profiles](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteCustomProfile(deleteCustomProfileOptions *DeleteCustomProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.DeleteCustomProfileWithContext(context.Background(), deleteCustomProfileOptions)
}

// DeleteCustomProfileWithContext is an alternate form of the DeleteCustomProfile method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteCustomProfileWithContext(ctx context.Context, deleteCustomProfileOptions *DeleteCustomProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCustomProfileOptions, "deleteCustomProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCustomProfileOptions, "deleteCustomProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"profile_id": *deleteCustomProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *deleteCustomProfileOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles/{profile_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCustomProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "DeleteCustomProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteCustomProfileOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteCustomProfileOptions.XCorrelationID))
	}
	if deleteCustomProfileOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*deleteCustomProfileOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProfile : Get a profile
// View the details of a profile by providing the profile ID.  You can find the profile ID in the Security and
// Compliance Center UI. For more information, see [Building custom
// profiles](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProfile(getProfileOptions *GetProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetProfileWithContext(context.Background(), getProfileOptions)
}

// GetProfileWithContext is an alternate form of the GetProfile method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProfileWithContext(ctx context.Context, getProfileOptions *GetProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileOptions, "getProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileOptions, "getProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"profile_id": *getProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getProfileOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles/{profile_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getProfileOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getProfileOptions.XCorrelationID))
	}
	if getProfileOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getProfileOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceProfile : Update a profile
// Update the details of a custom profile. With Security and Compliance Center, you can manage  a profile that is
// specific to your usecase, by using an existing library as a starting point.  For more information, see [Building
// custom
// profiles](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ReplaceProfile(replaceProfileOptions *ReplaceProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ReplaceProfileWithContext(context.Background(), replaceProfileOptions)
}

// ReplaceProfileWithContext is an alternate form of the ReplaceProfile method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ReplaceProfileWithContext(ctx context.Context, replaceProfileOptions *ReplaceProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceProfileOptions, "replaceProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceProfileOptions, "replaceProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"profile_id": *replaceProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *replaceProfileOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles/{profile_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ReplaceProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceProfileOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*replaceProfileOptions.XCorrelationID))
	}
	if replaceProfileOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*replaceProfileOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if replaceProfileOptions.ProfileName != nil {
		body["profile_name"] = replaceProfileOptions.ProfileName
	}
	if replaceProfileOptions.ProfileDescription != nil {
		body["profile_description"] = replaceProfileOptions.ProfileDescription
	}
	if replaceProfileOptions.ProfileType != nil {
		body["profile_type"] = replaceProfileOptions.ProfileType
	}
	if replaceProfileOptions.Controls != nil {
		body["controls"] = replaceProfileOptions.Controls
	}
	if replaceProfileOptions.DefaultParameters != nil {
		body["default_parameters"] = replaceProfileOptions.DefaultParameters
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListRules : List all rules
// Retrieve all the rules that you use to target the exact configuration properties  that you need to ensure are
// compliant. For more information, see [Defining custom
// rules](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListRules(listRulesOptions *ListRulesOptions) (result *RulesPageBase, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListRulesWithContext(context.Background(), listRulesOptions)
}

// ListRulesWithContext is an alternate form of the ListRules method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListRulesWithContext(ctx context.Context, listRulesOptions *ListRulesOptions) (result *RulesPageBase, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listRulesOptions, "listRulesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listRulesOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/rules`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listRulesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*listRulesOptions.XCorrelationID))
	}
	if listRulesOptions.XRequestID != nil {
		builder.AddHeader("X-Request-Id", fmt.Sprint(*listRulesOptions.XRequestID))
	}

	if listRulesOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listRulesOptions.Type))
	}
	if listRulesOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listRulesOptions.Search))
	}
	if listRulesOptions.ServiceName != nil {
		builder.AddQuery("service_name", fmt.Sprint(*listRulesOptions.ServiceName))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesPageBase)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateRule : Create a custom rule
// Create a custom rule to to target the exact configuration properties  that you need to evaluate your resources for
// compliance. For more information, see [Defining custom
// rules](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateRule(createRuleOptions *CreateRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.CreateRuleWithContext(context.Background(), createRuleOptions)
}

// CreateRuleWithContext is an alternate form of the CreateRule method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateRuleWithContext(ctx context.Context, createRuleOptions *CreateRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRuleOptions, "createRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createRuleOptions, "createRuleOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *createRuleOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/rules`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "CreateRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*createRuleOptions.XCorrelationID))
	}
	if createRuleOptions.XRequestID != nil {
		builder.AddHeader("X-Request-Id", fmt.Sprint(*createRuleOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if createRuleOptions.Description != nil {
		body["description"] = createRuleOptions.Description
	}
	if createRuleOptions.Target != nil {
		body["target"] = createRuleOptions.Target
	}
	if createRuleOptions.RequiredConfig != nil {
		body["required_config"] = createRuleOptions.RequiredConfig
	}
	if createRuleOptions.Type != nil {
		body["type"] = createRuleOptions.Type
	}
	if createRuleOptions.Version != nil {
		body["version"] = createRuleOptions.Version
	}
	if createRuleOptions.Import != nil {
		body["import"] = createRuleOptions.Import
	}
	if createRuleOptions.Labels != nil {
		body["labels"] = createRuleOptions.Labels
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRule)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteRule : Delete a custom rule
// Delete a custom rule that you no longer require to evaluate your resources. For more information, see [Defining
// custom rules](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteRule(deleteRuleOptions *DeleteRuleOptions) (response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.DeleteRuleWithContext(context.Background(), deleteRuleOptions)
}

// DeleteRuleWithContext is an alternate form of the DeleteRule method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteRuleWithContext(ctx context.Context, deleteRuleOptions *DeleteRuleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRuleOptions, "deleteRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteRuleOptions, "deleteRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *deleteRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *deleteRuleOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "DeleteRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*deleteRuleOptions.XCorrelationID))
	}
	if deleteRuleOptions.XRequestID != nil {
		builder.AddHeader("X-Request-Id", fmt.Sprint(*deleteRuleOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenterApi.Service.Request(request, nil)

	return
}

// GetRule : Get a custom rule
// Retrieve a rule that you created to evaluate your resources.  For more information, see [Defining custom
// rules](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetRule(getRuleOptions *GetRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetRuleWithContext(context.Background(), getRuleOptions)
}

// GetRuleWithContext is an alternate form of the GetRule method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetRuleWithContext(ctx context.Context, getRuleOptions *GetRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRuleOptions, "getRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRuleOptions, "getRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *getRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getRuleOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*getRuleOptions.XCorrelationID))
	}
	if getRuleOptions.XRequestID != nil {
		builder.AddHeader("X-Request-Id", fmt.Sprint(*getRuleOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRule)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceRule : Update a custom rule
// Update a custom rule that you use to target the exact configuration properties  that you need to evaluate your
// resources for compliance. For more information, see [Defining custom
// rules](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ReplaceRule(replaceRuleOptions *ReplaceRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ReplaceRuleWithContext(context.Background(), replaceRuleOptions)
}

// ReplaceRuleWithContext is an alternate form of the ReplaceRule method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ReplaceRuleWithContext(ctx context.Context, replaceRuleOptions *ReplaceRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceRuleOptions, "replaceRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceRuleOptions, "replaceRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *replaceRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *replaceRuleOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ReplaceRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceRuleOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replaceRuleOptions.IfMatch))
	}
	if replaceRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*replaceRuleOptions.XCorrelationID))
	}
	if replaceRuleOptions.XRequestID != nil {
		builder.AddHeader("X-Request-Id", fmt.Sprint(*replaceRuleOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if replaceRuleOptions.Description != nil {
		body["description"] = replaceRuleOptions.Description
	}
	if replaceRuleOptions.Target != nil {
		body["target"] = replaceRuleOptions.Target
	}
	if replaceRuleOptions.RequiredConfig != nil {
		body["required_config"] = replaceRuleOptions.RequiredConfig
	}
	if replaceRuleOptions.Type != nil {
		body["type"] = replaceRuleOptions.Type
	}
	if replaceRuleOptions.Version != nil {
		body["version"] = replaceRuleOptions.Version
	}
	if replaceRuleOptions.Import != nil {
		body["import"] = replaceRuleOptions.Import
	}
	if replaceRuleOptions.Labels != nil {
		body["labels"] = replaceRuleOptions.Labels
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRule)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListAttachments : Get all attachments linked to a specific profile
// View all of the attachments that are linked to a specific profile.  An attachment is the association between the set
// of resources that you want to evaluate  and a profile that contains the specific controls that you want to use. For
// more information, see [Running an evaluation for IBM
// Cloud](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListAttachments(listAttachmentsOptions *ListAttachmentsOptions) (result *AttachmentCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListAttachmentsWithContext(context.Background(), listAttachmentsOptions)
}

// ListAttachmentsWithContext is an alternate form of the ListAttachments method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListAttachmentsWithContext(ctx context.Context, listAttachmentsOptions *ListAttachmentsOptions) (result *AttachmentCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAttachmentsOptions, "listAttachmentsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listAttachmentsOptions, "listAttachmentsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"profile_id": *listAttachmentsOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listAttachmentsOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles/{profile_id}/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAttachmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListAttachments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listAttachmentsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listAttachmentsOptions.XCorrelationID))
	}
	if listAttachmentsOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listAttachmentsOptions.XRequestID))
	}

	if listAttachmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAttachmentsOptions.Limit))
	}
	if listAttachmentsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listAttachmentsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAttachmentCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAttachment : Create an attachment
// Create an attachment to link to a profile to schedule evaluations  of your resources on a recurring schedule, or
// on-demand. For more information, see [Running an evaluation for IBM
// Cloud](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateAttachment(createAttachmentOptions *CreateAttachmentOptions) (result *AttachmentPrototype, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.CreateAttachmentWithContext(context.Background(), createAttachmentOptions)
}

// CreateAttachmentWithContext is an alternate form of the CreateAttachment method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateAttachmentWithContext(ctx context.Context, createAttachmentOptions *CreateAttachmentOptions) (result *AttachmentPrototype, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAttachmentOptions, "createAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createAttachmentOptions, "createAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"profile_id": *createAttachmentOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *createAttachmentOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles/{profile_id}/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "CreateAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createAttachmentOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createAttachmentOptions.XCorrelationID))
	}
	if createAttachmentOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*createAttachmentOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if createAttachmentOptions.Attachments != nil {
		body["attachments"] = createAttachmentOptions.Attachments
	}
	if createAttachmentOptions.ProfileID != nil {
		body["profile_id"] = createAttachmentOptions.ProfileID
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAttachmentPrototype)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteProfileAttachment : Delete an attachment
// Delete an attachment. Alternatively, if you think that you might need  this configuration in the future, you can
// pause an attachment to stop being charged. For more information, see [Running an evaluation for IBM
// Cloud](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteProfileAttachment(deleteProfileAttachmentOptions *DeleteProfileAttachmentOptions) (result *AttachmentItem, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.DeleteProfileAttachmentWithContext(context.Background(), deleteProfileAttachmentOptions)
}

// DeleteProfileAttachmentWithContext is an alternate form of the DeleteProfileAttachment method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteProfileAttachmentWithContext(ctx context.Context, deleteProfileAttachmentOptions *DeleteProfileAttachmentOptions) (result *AttachmentItem, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProfileAttachmentOptions, "deleteProfileAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProfileAttachmentOptions, "deleteProfileAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"attachment_id": *deleteProfileAttachmentOptions.AttachmentID,
		"profile_id":    *deleteProfileAttachmentOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *deleteProfileAttachmentOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles/{profile_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteProfileAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "DeleteProfileAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteProfileAttachmentOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteProfileAttachmentOptions.XCorrelationID))
	}
	if deleteProfileAttachmentOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*deleteProfileAttachmentOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAttachmentItem)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProfileAttachment : Get an attachment
// View the details of an attachment a profile by providing the attachment ID.  You can find this value in the Security
// and Compliance Center UI. For more information, see [Running an evaluation for IBM
// Cloud](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProfileAttachment(getProfileAttachmentOptions *GetProfileAttachmentOptions) (result *AttachmentItem, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetProfileAttachmentWithContext(context.Background(), getProfileAttachmentOptions)
}

// GetProfileAttachmentWithContext is an alternate form of the GetProfileAttachment method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProfileAttachmentWithContext(ctx context.Context, getProfileAttachmentOptions *GetProfileAttachmentOptions) (result *AttachmentItem, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileAttachmentOptions, "getProfileAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileAttachmentOptions, "getProfileAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"attachment_id": *getProfileAttachmentOptions.AttachmentID,
		"profile_id":    *getProfileAttachmentOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getProfileAttachmentOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles/{profile_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProfileAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetProfileAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getProfileAttachmentOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getProfileAttachmentOptions.XCorrelationID))
	}
	if getProfileAttachmentOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getProfileAttachmentOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAttachmentItem)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceProfileAttachment : Update an attachment
// Update an attachment that is linked to a profile to evaluate your resources  on a recurring schedule, or on-demand.
// For more information, see [Running an evaluation for IBM
// Cloud](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ReplaceProfileAttachment(replaceProfileAttachmentOptions *ReplaceProfileAttachmentOptions) (result *AttachmentItem, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ReplaceProfileAttachmentWithContext(context.Background(), replaceProfileAttachmentOptions)
}

// ReplaceProfileAttachmentWithContext is an alternate form of the ReplaceProfileAttachment method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ReplaceProfileAttachmentWithContext(ctx context.Context, replaceProfileAttachmentOptions *ReplaceProfileAttachmentOptions) (result *AttachmentItem, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceProfileAttachmentOptions, "replaceProfileAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceProfileAttachmentOptions, "replaceProfileAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"attachment_id": *replaceProfileAttachmentOptions.AttachmentID,
		"profile_id":    *replaceProfileAttachmentOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *replaceProfileAttachmentOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/profiles/{profile_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceProfileAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ReplaceProfileAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceProfileAttachmentOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*replaceProfileAttachmentOptions.XCorrelationID))
	}
	if replaceProfileAttachmentOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*replaceProfileAttachmentOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if replaceProfileAttachmentOptions.ID != nil {
		body["id"] = replaceProfileAttachmentOptions.ID
	}
	if replaceProfileAttachmentOptions.ProfileID != nil {
		body["profile_id"] = replaceProfileAttachmentOptions.ProfileID
	}
	if replaceProfileAttachmentOptions.AccountID != nil {
		body["account_id"] = replaceProfileAttachmentOptions.AccountID
	}
	if replaceProfileAttachmentOptions.InstanceID != nil {
		body["instance_id"] = replaceProfileAttachmentOptions.InstanceID
	}
	if replaceProfileAttachmentOptions.Scope != nil {
		body["scope"] = replaceProfileAttachmentOptions.Scope
	}
	if replaceProfileAttachmentOptions.CreatedOn != nil {
		body["created_on"] = replaceProfileAttachmentOptions.CreatedOn
	}
	if replaceProfileAttachmentOptions.CreatedBy != nil {
		body["created_by"] = replaceProfileAttachmentOptions.CreatedBy
	}
	if replaceProfileAttachmentOptions.UpdatedOn != nil {
		body["updated_on"] = replaceProfileAttachmentOptions.UpdatedOn
	}
	if replaceProfileAttachmentOptions.UpdatedBy != nil {
		body["updated_by"] = replaceProfileAttachmentOptions.UpdatedBy
	}
	if replaceProfileAttachmentOptions.Status != nil {
		body["status"] = replaceProfileAttachmentOptions.Status
	}
	if replaceProfileAttachmentOptions.Schedule != nil {
		body["schedule"] = replaceProfileAttachmentOptions.Schedule
	}
	if replaceProfileAttachmentOptions.Notifications != nil {
		body["notifications"] = replaceProfileAttachmentOptions.Notifications
	}
	if replaceProfileAttachmentOptions.AttachmentParameters != nil {
		body["attachment_parameters"] = replaceProfileAttachmentOptions.AttachmentParameters
	}
	if replaceProfileAttachmentOptions.LastScan != nil {
		body["last_scan"] = replaceProfileAttachmentOptions.LastScan
	}
	if replaceProfileAttachmentOptions.NextScanTime != nil {
		body["next_scan_time"] = replaceProfileAttachmentOptions.NextScanTime
	}
	if replaceProfileAttachmentOptions.Name != nil {
		body["name"] = replaceProfileAttachmentOptions.Name
	}
	if replaceProfileAttachmentOptions.Description != nil {
		body["description"] = replaceProfileAttachmentOptions.Description
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAttachmentItem)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateScan : Create a scan
// Create a scan to evaluate your resources on a recurring basis or on demand.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateScan(createScanOptions *CreateScanOptions) (result *Scan, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.CreateScanWithContext(context.Background(), createScanOptions)
}

// CreateScanWithContext is an alternate form of the CreateScan method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateScanWithContext(ctx context.Context, createScanOptions *CreateScanOptions) (result *Scan, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createScanOptions, "createScanOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createScanOptions, "createScanOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *createScanOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/scans`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createScanOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "CreateScan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createScanOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createScanOptions.XCorrelationID))
	}
	if createScanOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*createScanOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if createScanOptions.AttachmentID != nil {
		body["attachment_id"] = createScanOptions.AttachmentID
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScan)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListAttachmentsAccount : Get all attachments in an instance
// View all of the attachments that are linked to an account. An attachment is the association between the set of
// resources that you want to evaluate  and a profile that contains the specific controls that you want to use. For more
// information, see [Running an evaluation for IBM
// Cloud](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListAttachmentsAccount(listAttachmentsAccountOptions *ListAttachmentsAccountOptions) (result *AttachmentCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListAttachmentsAccountWithContext(context.Background(), listAttachmentsAccountOptions)
}

// ListAttachmentsAccountWithContext is an alternate form of the ListAttachmentsAccount method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListAttachmentsAccountWithContext(ctx context.Context, listAttachmentsAccountOptions *ListAttachmentsAccountOptions) (result *AttachmentCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAttachmentsAccountOptions, "listAttachmentsAccountOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listAttachmentsAccountOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/attachments`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAttachmentsAccountOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListAttachmentsAccount")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listAttachmentsAccountOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listAttachmentsAccountOptions.XCorrelationID))
	}
	if listAttachmentsAccountOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listAttachmentsAccountOptions.XRequestID))
	}

	if listAttachmentsAccountOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAttachmentsAccountOptions.Limit))
	}
	if listAttachmentsAccountOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listAttachmentsAccountOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAttachmentCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLatestReports : Get the latest reports
// Retrieve the latest reports, which are grouped by profile ID, scope ID, and attachment ID. For more information, see
// [Viewing results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetLatestReports(getLatestReportsOptions *GetLatestReportsOptions) (result *ReportLatest, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetLatestReportsWithContext(context.Background(), getLatestReportsOptions)
}

// GetLatestReportsWithContext is an alternate form of the GetLatestReports method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetLatestReportsWithContext(ctx context.Context, getLatestReportsOptions *GetLatestReportsOptions) (result *ReportLatest, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getLatestReportsOptions, "getLatestReportsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getLatestReportsOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/latest`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLatestReportsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetLatestReports")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getLatestReportsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getLatestReportsOptions.XCorrelationID))
	}
	if getLatestReportsOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getLatestReportsOptions.XRequestID))
	}

	if getLatestReportsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*getLatestReportsOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportLatest)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListReports : List reports
// Retrieve a page of reports that are filtered by the specified parameters. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListReports(listReportsOptions *ListReportsOptions) (result *ReportPage, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListReportsWithContext(context.Background(), listReportsOptions)
}

// ListReportsWithContext is an alternate form of the ListReports method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListReportsWithContext(ctx context.Context, listReportsOptions *ListReportsOptions) (result *ReportPage, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listReportsOptions, "listReportsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listReportsOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listReportsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListReports")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listReportsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listReportsOptions.XCorrelationID))
	}
	if listReportsOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listReportsOptions.XRequestID))
	}

	if listReportsOptions.AttachmentID != nil {
		builder.AddQuery("attachment_id", fmt.Sprint(*listReportsOptions.AttachmentID))
	}
	if listReportsOptions.GroupID != nil {
		builder.AddQuery("group_id", fmt.Sprint(*listReportsOptions.GroupID))
	}
	if listReportsOptions.ProfileID != nil {
		builder.AddQuery("profile_id", fmt.Sprint(*listReportsOptions.ProfileID))
	}
	if listReportsOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listReportsOptions.Type))
	}
	if listReportsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listReportsOptions.Start))
	}
	if listReportsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listReportsOptions.Limit))
	}
	if listReportsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listReportsOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportPage)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReport : Get a report
// Retrieve a report by specifying its ID. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReport(getReportOptions *GetReportOptions) (result *Report, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetReportWithContext(context.Background(), getReportOptions)
}

// GetReportWithContext is an alternate form of the GetReport method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportWithContext(ctx context.Context, getReportOptions *GetReportOptions) (result *Report, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportOptions, "getReportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportOptions, "getReportOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *getReportOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getReportOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetReport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getReportOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getReportOptions.XCorrelationID))
	}
	if getReportOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getReportOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReport)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportSummary : Get a report summary
// Retrieve the complete summarized information for a report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportSummary(getReportSummaryOptions *GetReportSummaryOptions) (result *ReportSummary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetReportSummaryWithContext(context.Background(), getReportSummaryOptions)
}

// GetReportSummaryWithContext is an alternate form of the GetReportSummary method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportSummaryWithContext(ctx context.Context, getReportSummaryOptions *GetReportSummaryOptions) (result *ReportSummary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportSummaryOptions, "getReportSummaryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportSummaryOptions, "getReportSummaryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *getReportSummaryOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getReportSummaryOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}/summary`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportSummaryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetReportSummary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getReportSummaryOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getReportSummaryOptions.XCorrelationID))
	}
	if getReportSummaryOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getReportSummaryOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportSummary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportEvaluation : Get report evaluation details
// Retrieve the evaluation details of a report by specifying the report ID. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportEvaluation(getReportEvaluationOptions *GetReportEvaluationOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetReportEvaluationWithContext(context.Background(), getReportEvaluationOptions)
}

// GetReportEvaluationWithContext is an alternate form of the GetReportEvaluation method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportEvaluationWithContext(ctx context.Context, getReportEvaluationOptions *GetReportEvaluationOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportEvaluationOptions, "getReportEvaluationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportEvaluationOptions, "getReportEvaluationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *getReportEvaluationOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getReportEvaluationOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}/download`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportEvaluationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetReportEvaluation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/csv")
	if getReportEvaluationOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getReportEvaluationOptions.XCorrelationID))
	}
	if getReportEvaluationOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getReportEvaluationOptions.XRequestID))
	}

	if getReportEvaluationOptions.ExcludeSummary != nil {
		builder.AddQuery("exclude_summary", fmt.Sprint(*getReportEvaluationOptions.ExcludeSummary))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenterApi.Service.Request(request, &result)

	return
}

// GetReportControls : Get report controls
// Retrieve a sorted and filtered list of controls for the specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportControls(getReportControlsOptions *GetReportControlsOptions) (result *ReportControls, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetReportControlsWithContext(context.Background(), getReportControlsOptions)
}

// GetReportControlsWithContext is an alternate form of the GetReportControls method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportControlsWithContext(ctx context.Context, getReportControlsOptions *GetReportControlsOptions) (result *ReportControls, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportControlsOptions, "getReportControlsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportControlsOptions, "getReportControlsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *getReportControlsOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getReportControlsOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}/controls`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportControlsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetReportControls")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getReportControlsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getReportControlsOptions.XCorrelationID))
	}
	if getReportControlsOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getReportControlsOptions.XRequestID))
	}

	if getReportControlsOptions.ControlID != nil {
		builder.AddQuery("control_id", fmt.Sprint(*getReportControlsOptions.ControlID))
	}
	if getReportControlsOptions.ControlName != nil {
		builder.AddQuery("control_name", fmt.Sprint(*getReportControlsOptions.ControlName))
	}
	if getReportControlsOptions.ControlDescription != nil {
		builder.AddQuery("control_description", fmt.Sprint(*getReportControlsOptions.ControlDescription))
	}
	if getReportControlsOptions.ControlCategory != nil {
		builder.AddQuery("control_category", fmt.Sprint(*getReportControlsOptions.ControlCategory))
	}
	if getReportControlsOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*getReportControlsOptions.Status))
	}
	if getReportControlsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*getReportControlsOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportControls)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportRule : Get a report rule
// Retrieve the rule by specifying the report ID and rule ID. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportRule(getReportRuleOptions *GetReportRuleOptions) (result *RuleInfo, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetReportRuleWithContext(context.Background(), getReportRuleOptions)
}

// GetReportRuleWithContext is an alternate form of the GetReportRule method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportRuleWithContext(ctx context.Context, getReportRuleOptions *GetReportRuleOptions) (result *RuleInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportRuleOptions, "getReportRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportRuleOptions, "getReportRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *getReportRuleOptions.ReportID,
		"rule_id":   *getReportRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getReportRuleOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetReportRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getReportRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getReportRuleOptions.XCorrelationID))
	}
	if getReportRuleOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getReportRuleOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleInfo)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListReportEvaluations : List report evaluations
// Get a paginated list of evaluations for the specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListReportEvaluations(listReportEvaluationsOptions *ListReportEvaluationsOptions) (result *EvaluationPage, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListReportEvaluationsWithContext(context.Background(), listReportEvaluationsOptions)
}

// ListReportEvaluationsWithContext is an alternate form of the ListReportEvaluations method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListReportEvaluationsWithContext(ctx context.Context, listReportEvaluationsOptions *ListReportEvaluationsOptions) (result *EvaluationPage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listReportEvaluationsOptions, "listReportEvaluationsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listReportEvaluationsOptions, "listReportEvaluationsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *listReportEvaluationsOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listReportEvaluationsOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}/evaluations`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listReportEvaluationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListReportEvaluations")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listReportEvaluationsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listReportEvaluationsOptions.XCorrelationID))
	}
	if listReportEvaluationsOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listReportEvaluationsOptions.XRequestID))
	}

	if listReportEvaluationsOptions.AssessmentID != nil {
		builder.AddQuery("assessment_id", fmt.Sprint(*listReportEvaluationsOptions.AssessmentID))
	}
	if listReportEvaluationsOptions.ComponentID != nil {
		builder.AddQuery("component_id", fmt.Sprint(*listReportEvaluationsOptions.ComponentID))
	}
	if listReportEvaluationsOptions.TargetID != nil {
		builder.AddQuery("target_id", fmt.Sprint(*listReportEvaluationsOptions.TargetID))
	}
	if listReportEvaluationsOptions.TargetName != nil {
		builder.AddQuery("target_name", fmt.Sprint(*listReportEvaluationsOptions.TargetName))
	}
	if listReportEvaluationsOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*listReportEvaluationsOptions.Status))
	}
	if listReportEvaluationsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listReportEvaluationsOptions.Start))
	}
	if listReportEvaluationsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listReportEvaluationsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEvaluationPage)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListReportResources : List report resources
// Get a paginated list of resources for the specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListReportResources(listReportResourcesOptions *ListReportResourcesOptions) (result *ResourcePage, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListReportResourcesWithContext(context.Background(), listReportResourcesOptions)
}

// ListReportResourcesWithContext is an alternate form of the ListReportResources method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListReportResourcesWithContext(ctx context.Context, listReportResourcesOptions *ListReportResourcesOptions) (result *ResourcePage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listReportResourcesOptions, "listReportResourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listReportResourcesOptions, "listReportResourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *listReportResourcesOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listReportResourcesOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}/resources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listReportResourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListReportResources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listReportResourcesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listReportResourcesOptions.XCorrelationID))
	}
	if listReportResourcesOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listReportResourcesOptions.XRequestID))
	}

	if listReportResourcesOptions.ID != nil {
		builder.AddQuery("id", fmt.Sprint(*listReportResourcesOptions.ID))
	}
	if listReportResourcesOptions.ResourceName != nil {
		builder.AddQuery("resource_name", fmt.Sprint(*listReportResourcesOptions.ResourceName))
	}
	if listReportResourcesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listReportResourcesOptions.AccountID))
	}
	if listReportResourcesOptions.ComponentID != nil {
		builder.AddQuery("component_id", fmt.Sprint(*listReportResourcesOptions.ComponentID))
	}
	if listReportResourcesOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*listReportResourcesOptions.Status))
	}
	if listReportResourcesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listReportResourcesOptions.Sort))
	}
	if listReportResourcesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listReportResourcesOptions.Start))
	}
	if listReportResourcesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listReportResourcesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourcePage)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportTags : Get report tags
// Retrieve a list of tags for the specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportTags(getReportTagsOptions *GetReportTagsOptions) (result *ReportTags, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetReportTagsWithContext(context.Background(), getReportTagsOptions)
}

// GetReportTagsWithContext is an alternate form of the GetReportTags method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportTagsWithContext(ctx context.Context, getReportTagsOptions *GetReportTagsOptions) (result *ReportTags, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportTagsOptions, "getReportTagsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportTagsOptions, "getReportTagsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *getReportTagsOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getReportTagsOptions.InstanceID)
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}/tags`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportTagsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetReportTags")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getReportTagsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getReportTagsOptions.XCorrelationID))
	}
	if getReportTagsOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getReportTagsOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportTags)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportViolationsDrift : Get report violations drift
// Get a list of report violation data points for the specified report and time frame. For more information, see
// [Viewing results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportViolationsDrift(getReportViolationsDriftOptions *GetReportViolationsDriftOptions) (result *ReportViolationsDrift, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetReportViolationsDriftWithContext(context.Background(), getReportViolationsDriftOptions)
}

// GetReportViolationsDriftWithContext is an alternate form of the GetReportViolationsDrift method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetReportViolationsDriftWithContext(ctx context.Context, getReportViolationsDriftOptions *GetReportViolationsDriftOptions) (result *ReportViolationsDrift, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportViolationsDriftOptions, "getReportViolationsDriftOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportViolationsDriftOptions, "getReportViolationsDriftOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id": *getReportViolationsDriftOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getReportViolationsDriftOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/reports/{report_id}/violations_drift`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportViolationsDriftOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetReportViolationsDrift")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getReportViolationsDriftOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getReportViolationsDriftOptions.XCorrelationID))
	}
	if getReportViolationsDriftOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getReportViolationsDriftOptions.XRequestID))
	}

	if getReportViolationsDriftOptions.ScanTimeDuration != nil {
		builder.AddQuery("scan_time_duration", fmt.Sprint(*getReportViolationsDriftOptions.ScanTimeDuration))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportViolationsDrift)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProviderTypes : List all provider types
// List all the registered provider types. For more information about connecting Workload Protection with the Security
// and Compliance Center, see [Connecting Workload
// Protection](/docs/security-compliance?topic=security-compliance-setup-workload-protection&interface=api#wp-register).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListProviderTypes(listProviderTypesOptions *ListProviderTypesOptions) (result *ProviderTypesCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListProviderTypesWithContext(context.Background(), listProviderTypesOptions)
}

// ListProviderTypesWithContext is an alternate form of the ListProviderTypes method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListProviderTypesWithContext(ctx context.Context, listProviderTypesOptions *ListProviderTypesOptions) (result *ProviderTypesCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProviderTypesOptions, "listProviderTypesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listProviderTypesOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/provider_types`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProviderTypesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListProviderTypes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listProviderTypesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listProviderTypesOptions.XCorrelationID))
	}
	if listProviderTypesOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listProviderTypesOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypesCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProviderTypeByID : Get a provider type
// Retrieve a provider type by specifying its ID. For more information about integrations, see [Connecting Workload
// Protection](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProviderTypeByID(getProviderTypeByIdOptions *GetProviderTypeByIdOptions) (result *ProviderTypeItem, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetProviderTypeByIDWithContext(context.Background(), getProviderTypeByIdOptions)
}

// GetProviderTypeByIDWithContext is an alternate form of the GetProviderTypeByID method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProviderTypeByIDWithContext(ctx context.Context, getProviderTypeByIdOptions *GetProviderTypeByIdOptions) (result *ProviderTypeItem, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProviderTypeByIdOptions, "getProviderTypeByIdOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProviderTypeByIdOptions, "getProviderTypeByIdOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"provider_type_id": *getProviderTypeByIdOptions.ProviderTypeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenterApi.Service.Options.URL, `/v3/provider_types/{provider_type_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProviderTypeByIdOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetProviderTypeByID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getProviderTypeByIdOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getProviderTypeByIdOptions.XCorrelationID))
	}
	if getProviderTypeByIdOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getProviderTypeByIdOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeItem)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProviderTypeInstances : List all provider type instances
// Retrieve all instances of provider type. For more information about integrations, see [Connecting Workload
// Protection](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListProviderTypeInstances(listProviderTypeInstancesOptions *ListProviderTypeInstancesOptions) (result *ProviderTypeInstancesResponse, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.ListProviderTypeInstancesWithContext(context.Background(), listProviderTypeInstancesOptions)
}

// ListProviderTypeInstancesWithContext is an alternate form of the ListProviderTypeInstances method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) ListProviderTypeInstancesWithContext(ctx context.Context, listProviderTypeInstancesOptions *ListProviderTypeInstancesOptions) (result *ProviderTypeInstancesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProviderTypeInstancesOptions, "listProviderTypeInstancesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listProviderTypeInstancesOptions, "listProviderTypeInstancesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"provider_type_id": *listProviderTypeInstancesOptions.ProviderTypeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *listProviderTypeInstancesOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/provider_types/{provider_type_id}/provider_type_instances`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProviderTypeInstancesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "ListProviderTypeInstances")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listProviderTypeInstancesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listProviderTypeInstancesOptions.XCorrelationID))
	}
	if listProviderTypeInstancesOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*listProviderTypeInstancesOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeInstancesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateProviderTypeInstance : Create a provider type instance
// Create an instance of a provider type. For more information about integrations, see [Connecting Workload
// Protection](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateProviderTypeInstance(createProviderTypeInstanceOptions *CreateProviderTypeInstanceOptions) (result *ProviderTypeInstanceItem, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.CreateProviderTypeInstanceWithContext(context.Background(), createProviderTypeInstanceOptions)
}

// CreateProviderTypeInstanceWithContext is an alternate form of the CreateProviderTypeInstance method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) CreateProviderTypeInstanceWithContext(ctx context.Context, createProviderTypeInstanceOptions *CreateProviderTypeInstanceOptions) (result *ProviderTypeInstanceItem, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProviderTypeInstanceOptions, "createProviderTypeInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createProviderTypeInstanceOptions, "createProviderTypeInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"provider_type_id": *createProviderTypeInstanceOptions.ProviderTypeID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *createProviderTypeInstanceOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/provider_types/{provider_type_id}/provider_type_instances`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createProviderTypeInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "CreateProviderTypeInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createProviderTypeInstanceOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createProviderTypeInstanceOptions.XCorrelationID))
	}
	if createProviderTypeInstanceOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*createProviderTypeInstanceOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if createProviderTypeInstanceOptions.Name != nil {
		body["name"] = createProviderTypeInstanceOptions.Name
	}
	if createProviderTypeInstanceOptions.Attributes != nil {
		body["attributes"] = createProviderTypeInstanceOptions.Attributes
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeInstanceItem)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteProviderTypeInstance : Remove a specific instance of a provider type
// Remove a specific instance of a provider type.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteProviderTypeInstance(deleteProviderTypeInstanceOptions *DeleteProviderTypeInstanceOptions) (response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.DeleteProviderTypeInstanceWithContext(context.Background(), deleteProviderTypeInstanceOptions)
}

// DeleteProviderTypeInstanceWithContext is an alternate form of the DeleteProviderTypeInstance method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) DeleteProviderTypeInstanceWithContext(ctx context.Context, deleteProviderTypeInstanceOptions *DeleteProviderTypeInstanceOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProviderTypeInstanceOptions, "deleteProviderTypeInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProviderTypeInstanceOptions, "deleteProviderTypeInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"provider_type_id":          *deleteProviderTypeInstanceOptions.ProviderTypeID,
		"provider_type_instance_id": *deleteProviderTypeInstanceOptions.ProviderTypeInstanceID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *deleteProviderTypeInstanceOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/provider_types/{provider_type_id}/provider_type_instances/{provider_type_instance_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteProviderTypeInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "DeleteProviderTypeInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteProviderTypeInstanceOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteProviderTypeInstanceOptions.XCorrelationID))
	}
	if deleteProviderTypeInstanceOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*deleteProviderTypeInstanceOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenterApi.Service.Request(request, nil)

	return
}

// GetProviderTypeInstance : List a provider type instance
// Retrieve a provider type instance by specifying the provider type ID, and Security and Compliance Center instance ID.
// For more information about integrations, see [Connecting Workload
// Protection](https://test.cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProviderTypeInstance(getProviderTypeInstanceOptions *GetProviderTypeInstanceOptions) (result *ProviderTypeInstanceItem, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetProviderTypeInstanceWithContext(context.Background(), getProviderTypeInstanceOptions)
}

// GetProviderTypeInstanceWithContext is an alternate form of the GetProviderTypeInstance method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProviderTypeInstanceWithContext(ctx context.Context, getProviderTypeInstanceOptions *GetProviderTypeInstanceOptions) (result *ProviderTypeInstanceItem, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProviderTypeInstanceOptions, "getProviderTypeInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProviderTypeInstanceOptions, "getProviderTypeInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"provider_type_id":          *getProviderTypeInstanceOptions.ProviderTypeID,
		"provider_type_instance_id": *getProviderTypeInstanceOptions.ProviderTypeInstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getProviderTypeInstanceOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/provider_types/{provider_type_id}/provider_type_instances/{provider_type_instance_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProviderTypeInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetProviderTypeInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getProviderTypeInstanceOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getProviderTypeInstanceOptions.XCorrelationID))
	}
	if getProviderTypeInstanceOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getProviderTypeInstanceOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeInstanceItem)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateProviderTypeInstance : Patch a specific instance of a provider type
// Patch a specific instance of a provider type.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) UpdateProviderTypeInstance(updateProviderTypeInstanceOptions *UpdateProviderTypeInstanceOptions) (result *ProviderTypeInstanceItem, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.UpdateProviderTypeInstanceWithContext(context.Background(), updateProviderTypeInstanceOptions)
}

// UpdateProviderTypeInstanceWithContext is an alternate form of the UpdateProviderTypeInstance method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) UpdateProviderTypeInstanceWithContext(ctx context.Context, updateProviderTypeInstanceOptions *UpdateProviderTypeInstanceOptions) (result *ProviderTypeInstanceItem, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProviderTypeInstanceOptions, "updateProviderTypeInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateProviderTypeInstanceOptions, "updateProviderTypeInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"provider_type_id":          *updateProviderTypeInstanceOptions.ProviderTypeID,
		"provider_type_instance_id": *updateProviderTypeInstanceOptions.ProviderTypeInstanceID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *updateProviderTypeInstanceOptions.InstanceID)
	if err != nil {
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/provider_types/{provider_type_id}/provider_type_instances/{provider_type_instance_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateProviderTypeInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "UpdateProviderTypeInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateProviderTypeInstanceOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateProviderTypeInstanceOptions.XCorrelationID))
	}
	if updateProviderTypeInstanceOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*updateProviderTypeInstanceOptions.XRequestID))
	}

	body := make(map[string]interface{})
	if updateProviderTypeInstanceOptions.Name != nil {
		body["name"] = updateProviderTypeInstanceOptions.Name
	}
	if updateProviderTypeInstanceOptions.Attributes != nil {
		body["attributes"] = updateProviderTypeInstanceOptions.Attributes
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
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeInstanceItem)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProviderTypesInstances : Get a list of instances for all provider types
// Get a list of instances for all provider types.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProviderTypesInstances(getProviderTypesInstancesOptions *GetProviderTypesInstancesOptions) (result *ProviderTypesInstancesResponse, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenterApi.GetProviderTypesInstancesWithContext(context.Background(), getProviderTypesInstancesOptions)
}

// GetProviderTypesInstancesWithContext is an alternate form of the GetProviderTypesInstances method which supports a Context parameter
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) GetProviderTypesInstancesWithContext(ctx context.Context, getProviderTypesInstancesOptions *GetProviderTypesInstancesOptions) (result *ProviderTypesInstancesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getProviderTypesInstancesOptions, "getProviderTypesInstancesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenterApi.GetEnableGzipCompression()
	instanceURL, err := getInstanceBasedURL(securityAndComplianceCenterApi.Service.Options.URL, *getProviderTypesInstancesOptions.InstanceID)
	if err != nil {
		return
	}
	_, err = builder.ResolveRequestURL(instanceURL, `/provider_types_instances`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProviderTypesInstancesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center_api", "V3", "GetProviderTypesInstances")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getProviderTypesInstancesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getProviderTypesInstancesOptions.XCorrelationID))
	}
	if getProviderTypesInstancesOptions.XRequestID != nil {
		builder.AddHeader("X-Request-ID", fmt.Sprint(*getProviderTypesInstancesOptions.XRequestID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenterApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypesInstancesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// Account : The account that is associated with a report.
type Account struct {
	// The account ID.
	ID *string `json:"id,omitempty"`

	// The account name.
	Name *string `json:"name,omitempty"`

	// The account type.
	Type *string `json:"type,omitempty"`
}

// UnmarshalAccount unmarshals an instance of Account from the specified map of raw messages.
func UnmarshalAccount(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Account)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalProperty : AdditionalProperty struct
type AdditionalProperty struct {
	// An additional property that indicates the type of the attribute in various formats (text, url, secret, label,
	// masked).
	Type *string `json:"type" validate:"required"`

	// The name of the attribute that is displayed in the UI.
	DisplayName *string `json:"display_name" validate:"required"`
}

// Constants associated with the AdditionalProperty.Type property.
// An additional property that indicates the type of the attribute in various formats (text, url, secret, label,
// masked).
const (
	AdditionalProperty_Type_Label  = "label"
	AdditionalProperty_Type_Masked = "masked"
	AdditionalProperty_Type_Secret = "secret"
	AdditionalProperty_Type_Text   = "text"
	AdditionalProperty_Type_URL    = "url"
)

// UnmarshalAdditionalProperty unmarshals an instance of AdditionalProperty from the specified map of raw messages.
func UnmarshalAdditionalProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalProperty)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalTargetAttribute : The additional target attribute of the service.
type AdditionalTargetAttribute struct {
	// The additional target attribute name.
	Name *string `json:"name,omitempty"`

	// The operator.
	Operator *string `json:"operator,omitempty"`

	// The value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the AdditionalTargetAttribute.Operator property.
// The operator.
const (
	AdditionalTargetAttribute_Operator_DaysLessThan         = "days_less_than"
	AdditionalTargetAttribute_Operator_IpsEquals            = "ips_equals"
	AdditionalTargetAttribute_Operator_IpsInRange           = "ips_in_range"
	AdditionalTargetAttribute_Operator_IpsNotEquals         = "ips_not_equals"
	AdditionalTargetAttribute_Operator_IsEmpty              = "is_empty"
	AdditionalTargetAttribute_Operator_IsFalse              = "is_false"
	AdditionalTargetAttribute_Operator_IsNotEmpty           = "is_not_empty"
	AdditionalTargetAttribute_Operator_IsTrue               = "is_true"
	AdditionalTargetAttribute_Operator_NumEquals            = "num_equals"
	AdditionalTargetAttribute_Operator_NumGreaterThan       = "num_greater_than"
	AdditionalTargetAttribute_Operator_NumGreaterThanEquals = "num_greater_than_equals"
	AdditionalTargetAttribute_Operator_NumLessThan          = "num_less_than"
	AdditionalTargetAttribute_Operator_NumLessThanEquals    = "num_less_than_equals"
	AdditionalTargetAttribute_Operator_NumNotEquals         = "num_not_equals"
	AdditionalTargetAttribute_Operator_StringContains       = "string_contains"
	AdditionalTargetAttribute_Operator_StringEquals         = "string_equals"
	AdditionalTargetAttribute_Operator_StringMatch          = "string_match"
	AdditionalTargetAttribute_Operator_StringNotContains    = "string_not_contains"
	AdditionalTargetAttribute_Operator_StringNotEquals      = "string_not_equals"
	AdditionalTargetAttribute_Operator_StringNotMatch       = "string_not_match"
	AdditionalTargetAttribute_Operator_StringsAllowed       = "strings_allowed"
	AdditionalTargetAttribute_Operator_StringsInList        = "strings_in_list"
	AdditionalTargetAttribute_Operator_StringsRequired      = "strings_required"
)

// UnmarshalAdditionalTargetAttribute unmarshals an instance of AdditionalTargetAttribute from the specified map of raw messages.
func UnmarshalAdditionalTargetAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalTargetAttribute)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// Assessment : The control specification assessment.
type Assessment struct {
	// The assessment ID.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The assessment type.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The assessment method.
	AssessmentMethod *string `json:"assessment_method,omitempty"`

	// The assessment description.
	AssessmentDescription *string `json:"assessment_description,omitempty"`

	// The number of parameters of this assessment.
	ParameterCount *int64 `json:"parameter_count,omitempty"`

	// The list of parameters of this assessment.
	Parameters []ParameterInfo `json:"parameters,omitempty"`
}

// UnmarshalAssessment unmarshals an instance of Assessment from the specified map of raw messages.
func UnmarshalAssessment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Assessment)
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_method", &obj.AssessmentMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_description", &obj.AssessmentDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_count", &obj.ParameterCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "parameters", &obj.Parameters, UnmarshalParameterInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Attachment : The attachment that is associated with a report.
type Attachment struct {
	// The attachment ID.
	ID *string `json:"id,omitempty"`

	// The name of the attachment.
	Name *string `json:"name,omitempty"`

	// The description of the attachment.
	Description *string `json:"description,omitempty"`

	// The attachment schedule.
	Schedule *string `json:"schedule,omitempty"`

	// The scope of the attachment.
	Scope []AttachmentScope `json:"scope,omitempty"`
}

// UnmarshalAttachment unmarshals an instance of Attachment from the specified map of raw messages.
func UnmarshalAttachment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Attachment)
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
	err = core.UnmarshalPrimitive(m, "schedule", &obj.Schedule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scope", &obj.Scope, UnmarshalAttachmentScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AttachmentCollection : The response body of an attachment.
type AttachmentCollection struct {
	// The number of attachments.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The limit of attachments per request.
	Limit *int64 `json:"limit" validate:"required"`

	// The reference to the first page of entries.
	First *PaginatedCollectionFirst `json:"first" validate:"required"`

	// The reference URL for the next few entries.
	Next *PaginatedCollectionNext `json:"next" validate:"required"`

	// The list of attachments.
	Attachments []AttachmentItem `json:"attachments" validate:"required"`
}

// UnmarshalAttachmentCollection unmarshals an instance of AttachmentCollection from the specified map of raw messages.
func UnmarshalAttachmentCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedCollectionFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedCollectionNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachments", &obj.Attachments, UnmarshalAttachmentItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AttachmentCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// AttachmentItem : The request payload of the attachments parameter.
type AttachmentItem struct {
	// The ID of the attachment.
	ID *string `json:"id,omitempty"`

	// The ID of the profile that is specified in the attachment.
	ProfileID *string `json:"profile_id,omitempty"`

	// The account ID that is associated to the attachment.
	AccountID *string `json:"account_id,omitempty"`

	// The instance ID of the account that is associated to the attachment.
	InstanceID *string `json:"instance_id,omitempty"`

	// The scope payload for the multi cloud feature.
	Scope []MultiCloudScope `json:"scope,omitempty"`

	// The date when the attachment was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The user who created the attachment.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the attachment was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The user who updated the attachment.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The status of an attachment evaluation.
	Status *string `json:"status,omitempty"`

	// The schedule of an attachment evaluation.
	Schedule *string `json:"schedule,omitempty"`

	// The request payload of the attachment notifications.
	Notifications *AttachmentsNotificationsPrototype `json:"notifications,omitempty"`

	// The profile parameters for the attachment.
	AttachmentParameters []AttachmentParameterPrototype `json:"attachment_parameters,omitempty"`

	// The details of the last scan of an attachment.
	LastScan *LastScan `json:"last_scan,omitempty"`

	// The start time of the next scan.
	NextScanTime *strfmt.DateTime `json:"next_scan_time,omitempty"`

	// The name of the attachment.
	Name *string `json:"name,omitempty"`

	// The description for the attachment.
	Description *string `json:"description,omitempty"`
}

// Constants associated with the AttachmentItem.Status property.
// The status of an attachment evaluation.
const (
	AttachmentItem_Status_Disabled = "disabled"
	AttachmentItem_Status_Enabled  = "enabled"
)

// Constants associated with the AttachmentItem.Schedule property.
// The schedule of an attachment evaluation.
const (
	AttachmentItem_Schedule_Daily       = "daily"
	AttachmentItem_Schedule_Every30Days = "every_30_days"
	AttachmentItem_Schedule_Every7Days  = "every_7_days"
)

// UnmarshalAttachmentItem unmarshals an instance of AttachmentItem from the specified map of raw messages.
func UnmarshalAttachmentItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_id", &obj.ProfileID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scope", &obj.Scope, UnmarshalMultiCloudScope)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schedule", &obj.Schedule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "notifications", &obj.Notifications, UnmarshalAttachmentsNotificationsPrototype)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachment_parameters", &obj.AttachmentParameters, UnmarshalAttachmentParameterPrototype)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_scan", &obj.LastScan, UnmarshalLastScan)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_scan_time", &obj.NextScanTime)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AttachmentParameterPrototype : The parameters related to the Attachment.
type AttachmentParameterPrototype struct {
	// The type of the implementation.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The implementation ID of the parameter.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The parameter name.
	ParameterName *string `json:"parameter_name,omitempty"`

	// The value of the parameter.
	ParameterValue *string `json:"parameter_value,omitempty"`

	// The parameter display name.
	ParameterDisplayName *string `json:"parameter_display_name,omitempty"`

	// The parameter type.
	ParameterType *string `json:"parameter_type,omitempty"`
}

// Constants associated with the AttachmentParameterPrototype.ParameterType property.
// The parameter type.
const (
	AttachmentParameterPrototype_ParameterType_Boolean    = "boolean"
	AttachmentParameterPrototype_ParameterType_General    = "general"
	AttachmentParameterPrototype_ParameterType_IpList     = "ip_list"
	AttachmentParameterPrototype_ParameterType_Numeric    = "numeric"
	AttachmentParameterPrototype_ParameterType_String     = "string"
	AttachmentParameterPrototype_ParameterType_StringList = "string_list"
	AttachmentParameterPrototype_ParameterType_Timestamp  = "timestamp"
)

// UnmarshalAttachmentParameterPrototype unmarshals an instance of AttachmentParameterPrototype from the specified map of raw messages.
func UnmarshalAttachmentParameterPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentParameterPrototype)
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_name", &obj.ParameterName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_value", &obj.ParameterValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_display_name", &obj.ParameterDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_type", &obj.ParameterType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AttachmentPrototype : The request body of getting an attachment that is associated with your account.
type AttachmentPrototype struct {
	// The ID of the profile that is specified in the attachment.
	ProfileID *string `json:"profile_id,omitempty"`

	// The array that displays all of the available attachments.
	Attachments []AttachmentsPrototype `json:"attachments" validate:"required"`
}

// NewAttachmentPrototype : Instantiate AttachmentPrototype (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewAttachmentPrototype(attachments []AttachmentsPrototype) (_model *AttachmentPrototype, err error) {
	_model = &AttachmentPrototype{
		Attachments: attachments,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalAttachmentPrototype unmarshals an instance of AttachmentPrototype from the specified map of raw messages.
func UnmarshalAttachmentPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentPrototype)
	err = core.UnmarshalPrimitive(m, "profile_id", &obj.ProfileID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachments", &obj.Attachments, UnmarshalAttachmentsPrototype)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AttachmentScope : A scope of the attachment.
type AttachmentScope struct {
	// The unique identifier for this scope.
	ID *string `json:"id,omitempty"`

	// The environment that relates to this scope.
	Environment *string `json:"environment,omitempty"`

	// The properties that are supported for scoping by this environment.
	Properties []ScopeProperty `json:"properties,omitempty"`
}

// UnmarshalAttachmentScope unmarshals an instance of AttachmentScope from the specified map of raw messages.
func UnmarshalAttachmentScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentScope)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalScopeProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AttachmentsNotificationsPrototype : The request payload of the attachment notifications.
type AttachmentsNotificationsPrototype struct {
	// enabled notifications.
	Enabled *bool `json:"enabled" validate:"required"`

	// The failed controls.
	Controls *FailedControls `json:"controls" validate:"required"`
}

// NewAttachmentsNotificationsPrototype : Instantiate AttachmentsNotificationsPrototype (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewAttachmentsNotificationsPrototype(enabled bool, controls *FailedControls) (_model *AttachmentsNotificationsPrototype, err error) {
	_model = &AttachmentsNotificationsPrototype{
		Enabled:  core.BoolPtr(enabled),
		Controls: controls,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalAttachmentsNotificationsPrototype unmarshals an instance of AttachmentsNotificationsPrototype from the specified map of raw messages.
func UnmarshalAttachmentsNotificationsPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentsNotificationsPrototype)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalFailedControls)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AttachmentsPrototype : The request payload of getting all of the attachments that are associated with the account.
type AttachmentsPrototype struct {
	// The id that is generated from the scope type and ID.
	ID *string `json:"id,omitempty"`

	// The name that is generated from the scope type and ID.
	Name *string `json:"name" validate:"required"`

	// The description for the attachment.
	Description *string `json:"description,omitempty"`

	// The scope payload for the multi cloud feature.
	Scope []MultiCloudScope `json:"scope" validate:"required"`

	// The status of the scan of an attachment.
	Status *string `json:"status" validate:"required"`

	// The schedule of an attachment evaluation.
	Schedule *string `json:"schedule" validate:"required"`

	// The request payload of the attachment notifications.
	Notifications *AttachmentsNotificationsPrototype `json:"notifications,omitempty"`

	// The profile parameters for the attachment.
	AttachmentParameters []AttachmentParameterPrototype `json:"attachment_parameters" validate:"required"`
}

// Constants associated with the AttachmentsPrototype.Status property.
// The status of the scan of an attachment.
const (
	AttachmentsPrototype_Status_Disabled = "disabled"
	AttachmentsPrototype_Status_Enabled  = "enabled"
)

// Constants associated with the AttachmentsPrototype.Schedule property.
// The schedule of an attachment evaluation.
const (
	AttachmentsPrototype_Schedule_Daily       = "daily"
	AttachmentsPrototype_Schedule_Every30Days = "every_30_days"
	AttachmentsPrototype_Schedule_Every7Days  = "every_7_days"
)

// NewAttachmentsPrototype : Instantiate AttachmentsPrototype (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewAttachmentsPrototype(name string, scope []MultiCloudScope, status string, schedule string, attachmentParameters []AttachmentParameterPrototype) (_model *AttachmentsPrototype, err error) {
	_model = &AttachmentsPrototype{
		Name:                 core.StringPtr(name),
		Scope:                scope,
		Status:               core.StringPtr(status),
		Schedule:             core.StringPtr(schedule),
		AttachmentParameters: attachmentParameters,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalAttachmentsPrototype unmarshals an instance of AttachmentsPrototype from the specified map of raw messages.
func UnmarshalAttachmentsPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentsPrototype)
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
	err = core.UnmarshalModel(m, "scope", &obj.Scope, UnmarshalMultiCloudScope)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schedule", &obj.Schedule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "notifications", &obj.Notifications, UnmarshalAttachmentsNotificationsPrototype)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachment_parameters", &obj.AttachmentParameters, UnmarshalAttachmentParameterPrototype)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComplianceScore : The compliance score.
type ComplianceScore struct {
	// The number of successful evaluations.
	Passed *int64 `json:"passed,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The percentage of successful evaluations.
	Percent *int64 `json:"percent,omitempty"`
}

// UnmarshalComplianceScore unmarshals an instance of ComplianceScore from the specified map of raw messages.
func UnmarshalComplianceScore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceScore)
	err = core.UnmarshalPrimitive(m, "passed", &obj.Passed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "percent", &obj.Percent)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComplianceStats : The compliance stats.
type ComplianceStats struct {
	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`
}

// Constants associated with the ComplianceStats.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ComplianceStats_Status_Compliant              = "compliant"
	ComplianceStats_Status_NotCompliant           = "not_compliant"
	ComplianceStats_Status_UnableToPerform        = "unable_to_perform"
	ComplianceStats_Status_UserEvaluationRequired = "user_evaluation_required"
)

// UnmarshalComplianceStats unmarshals an instance of ComplianceStats from the specified map of raw messages.
func UnmarshalComplianceStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceStats)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlDocs : The control documentation.
type ControlDocs struct {
	// The ID of the control documentation.
	ControlDocsID *string `json:"control_docs_id,omitempty"`

	// The type of control documentation.
	ControlDocsType *string `json:"control_docs_type,omitempty"`
}

// UnmarshalControlDocs unmarshals an instance of ControlDocs from the specified map of raw messages.
func UnmarshalControlDocs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlDocs)
	err = core.UnmarshalPrimitive(m, "control_docs_id", &obj.ControlDocsID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_docs_type", &obj.ControlDocsType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlLibrary : The request payload of the control library.
type ControlLibrary struct {
	// The control library ID.
	ID *string `json:"id,omitempty"`

	// The account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The control library name.
	ControlLibraryName *string `json:"control_library_name,omitempty"`

	// The control library description.
	ControlLibraryDescription *string `json:"control_library_description,omitempty"`

	// The control library type.
	ControlLibraryType *string `json:"control_library_type,omitempty"`

	// The version group label.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// The control library version.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The date when the control library was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The user who created the control library.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the control library was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The user who updated the control library.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The latest version of the control library.
	Latest *bool `json:"latest,omitempty"`

	// The indication of whether hierarchy is enabled for the control library.
	HierarchyEnabled *bool `json:"hierarchy_enabled,omitempty"`

	// The number of controls.
	ControlsCount *int64 `json:"controls_count,omitempty"`

	// The number of parent controls in the control library.
	ControlParentsCount *int64 `json:"control_parents_count,omitempty"`

	// The list of controls in a control library.
	Controls []ControlsInControlLib `json:"controls,omitempty"`
}

// Constants associated with the ControlLibrary.ControlLibraryType property.
// The control library type.
const (
	ControlLibrary_ControlLibraryType_Custom     = "custom"
	ControlLibrary_ControlLibraryType_Predefined = "predefined"
)

// UnmarshalControlLibrary unmarshals an instance of ControlLibrary from the specified map of raw messages.
func UnmarshalControlLibrary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlLibrary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_name", &obj.ControlLibraryName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_description", &obj.ControlLibraryDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_type", &obj.ControlLibraryType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_group_label", &obj.VersionGroupLabel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_version", &obj.ControlLibraryVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "latest", &obj.Latest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hierarchy_enabled", &obj.HierarchyEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_count", &obj.ControlsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_parents_count", &obj.ControlParentsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalControlsInControlLib)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlLibraryCollection : The response body of control libraries.
type ControlLibraryCollection struct {
	// The number of control libraries.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The reference to the first page of entries.
	First *PaginatedCollectionFirst `json:"first" validate:"required"`

	// The reference URL for the next few entries.
	Next *PaginatedCollectionNext `json:"next" validate:"required"`

	// The control libraries.
	ControlLibraries []ControlLibraryItem `json:"control_libraries" validate:"required"`
}

// UnmarshalControlLibraryCollection unmarshals an instance of ControlLibraryCollection from the specified map of raw messages.
func UnmarshalControlLibraryCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlLibraryCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedCollectionFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedCollectionNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_libraries", &obj.ControlLibraries, UnmarshalControlLibraryItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ControlLibraryCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ControlLibraryDelete : The response body of deleting of a control library.
type ControlLibraryDelete struct {
	// The delete message of a control library.
	Deleted *string `json:"deleted,omitempty"`
}

// UnmarshalControlLibraryDelete unmarshals an instance of ControlLibraryDelete from the specified map of raw messages.
func UnmarshalControlLibraryDelete(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlLibraryDelete)
	err = core.UnmarshalPrimitive(m, "deleted", &obj.Deleted)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlLibraryItem : ControlLibraryItem struct
type ControlLibraryItem struct {
	// The ID of the control library.
	ID *string `json:"id,omitempty"`

	// The Account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The control library name.
	ControlLibraryName *string `json:"control_library_name,omitempty"`

	// The control library description.
	ControlLibraryDescription *string `json:"control_library_description,omitempty"`

	// The control library type.
	ControlLibraryType *string `json:"control_library_type,omitempty"`

	// The date when the control library was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The user who created the control library.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the control library was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The use who updated the control library.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The version group label.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// The control library version.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The latest control library version.
	Latest *bool `json:"latest,omitempty"`

	// The number of controls.
	ControlsCount *int64 `json:"controls_count,omitempty"`
}

// UnmarshalControlLibraryItem unmarshals an instance of ControlLibraryItem from the specified map of raw messages.
func UnmarshalControlLibraryItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlLibraryItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_name", &obj.ControlLibraryName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_description", &obj.ControlLibraryDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_type", &obj.ControlLibraryType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_group_label", &obj.VersionGroupLabel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_version", &obj.ControlLibraryVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "latest", &obj.Latest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_count", &obj.ControlsCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlSpecificationWithStats : The control specification with compliance stats.
type ControlSpecificationWithStats struct {
	// The control specification ID.
	ControlSpecificationID *string `json:"control_specification_id,omitempty"`

	// The component ID.
	ComponentID *string `json:"component_id,omitempty"`

	// The component description.
	ControlSpecificationDescription *string `json:"control_specification_description,omitempty"`

	// The environment.
	Environment *string `json:"environment,omitempty"`

	// The responsibility for managing control specifications.
	Responsibility *string `json:"responsibility,omitempty"`

	// The list of assessments.
	Assessments []Assessment `json:"assessments,omitempty"`

	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`
}

// Constants associated with the ControlSpecificationWithStats.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ControlSpecificationWithStats_Status_Compliant              = "compliant"
	ControlSpecificationWithStats_Status_NotCompliant           = "not_compliant"
	ControlSpecificationWithStats_Status_UnableToPerform        = "unable_to_perform"
	ControlSpecificationWithStats_Status_UserEvaluationRequired = "user_evaluation_required"
)

// UnmarshalControlSpecificationWithStats unmarshals an instance of ControlSpecificationWithStats from the specified map of raw messages.
func UnmarshalControlSpecificationWithStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlSpecificationWithStats)
	err = core.UnmarshalPrimitive(m, "control_specification_id", &obj.ControlSpecificationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_id", &obj.ComponentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_specification_description", &obj.ControlSpecificationDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "responsibility", &obj.Responsibility)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "assessments", &obj.Assessments, UnmarshalAssessment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlSpecifications : The control specifications of a control library.
type ControlSpecifications struct {
	// The control specification ID.
	ControlSpecificationID *string `json:"control_specification_id,omitempty"`

	// The responsibility for managing the control.
	Responsibility *string `json:"responsibility,omitempty"`

	// The component ID.
	ComponentID *string `json:"component_id,omitempty"`

	// The component name.
	ComponentName *string `json:"component_name,omitempty"`

	// The control specifications environment.
	Environment *string `json:"environment,omitempty"`

	// The control specifications description.
	ControlSpecificationDescription *string `json:"control_specification_description,omitempty"`

	// The number of assessments.
	AssessmentsCount *int64 `json:"assessments_count,omitempty"`

	// The assessments.
	Assessments []Implementation `json:"assessments,omitempty"`
}

// Constants associated with the ControlSpecifications.Responsibility property.
// The responsibility for managing the control.
const (
	ControlSpecifications_Responsibility_User = "user"
)

// UnmarshalControlSpecifications unmarshals an instance of ControlSpecifications from the specified map of raw messages.
func UnmarshalControlSpecifications(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlSpecifications)
	err = core.UnmarshalPrimitive(m, "control_specification_id", &obj.ControlSpecificationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "responsibility", &obj.Responsibility)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_id", &obj.ComponentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_name", &obj.ComponentName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_specification_description", &obj.ControlSpecificationDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessments_count", &obj.AssessmentsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "assessments", &obj.Assessments, UnmarshalImplementation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlWithStats : The control with compliance stats.
type ControlWithStats struct {
	// The control ID.
	ID *string `json:"id,omitempty"`

	// The control library ID.
	ControlLibraryID *string `json:"control_library_id,omitempty"`

	// The control library version.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The control name.
	ControlName *string `json:"control_name,omitempty"`

	// The control description.
	ControlDescription *string `json:"control_description,omitempty"`

	// The control category.
	ControlCategory *string `json:"control_category,omitempty"`

	// The control path.
	ControlPath *string `json:"control_path,omitempty"`

	// The list of specifications that are on the page.
	ControlSpecifications []ControlSpecificationWithStats `json:"control_specifications,omitempty"`

	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`
}

// Constants associated with the ControlWithStats.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ControlWithStats_Status_Compliant              = "compliant"
	ControlWithStats_Status_NotCompliant           = "not_compliant"
	ControlWithStats_Status_UnableToPerform        = "unable_to_perform"
	ControlWithStats_Status_UserEvaluationRequired = "user_evaluation_required"
)

// UnmarshalControlWithStats unmarshals an instance of ControlWithStats from the specified map of raw messages.
func UnmarshalControlWithStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlWithStats)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_id", &obj.ControlLibraryID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_version", &obj.ControlLibraryVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_name", &obj.ControlName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_description", &obj.ControlDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_category", &obj.ControlCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_path", &obj.ControlPath)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_specifications", &obj.ControlSpecifications, UnmarshalControlSpecificationWithStats)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlsInControlLib : The control details of a control library.
type ControlsInControlLib struct {
	// The ID of the control library that contains the profile.
	ControlName *string `json:"control_name,omitempty"`

	// The control name.
	ControlID *string `json:"control_id,omitempty"`

	// The control description.
	ControlDescription *string `json:"control_description,omitempty"`

	// The control category.
	ControlCategory *string `json:"control_category,omitempty"`

	// The parent control.
	ControlParent *string `json:"control_parent,omitempty"`

	// The control tags.
	ControlTags []string `json:"control_tags,omitempty"`

	// The control specifications.
	ControlSpecifications []ControlSpecifications `json:"control_specifications,omitempty"`

	// The control documentation.
	ControlDocs *ControlDocs `json:"control_docs,omitempty"`

	// Is this a control that can be automated or manually evaluated.
	ControlRequirement *bool `json:"control_requirement,omitempty"`

	// The control status.
	Status *string `json:"status,omitempty"`
}

// Constants associated with the ControlsInControlLib.Status property.
// The control status.
const (
	ControlsInControlLib_Status_Disabled = "disabled"
	ControlsInControlLib_Status_Enabled  = "enabled"
)

// UnmarshalControlsInControlLib unmarshals an instance of ControlsInControlLib from the specified map of raw messages.
func UnmarshalControlsInControlLib(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlsInControlLib)
	err = core.UnmarshalPrimitive(m, "control_name", &obj.ControlName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_id", &obj.ControlID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_description", &obj.ControlDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_category", &obj.ControlCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_parent", &obj.ControlParent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_tags", &obj.ControlTags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_specifications", &obj.ControlSpecifications, UnmarshalControlSpecifications)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_docs", &obj.ControlDocs, UnmarshalControlDocs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_requirement", &obj.ControlRequirement)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAttachmentOptions : The CreateAttachment options.
type CreateAttachmentOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The array that displays all of the available attachments.
	Attachments []AttachmentsPrototype `json:"attachments" validate:"required"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateAttachmentOptions : Instantiate CreateAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateAttachmentOptions(instanceID string, profileID string, attachments []AttachmentsPrototype) *CreateAttachmentOptions {
	return &CreateAttachmentOptions{
		InstanceID:  core.StringPtr(instanceID),
		ProfileID:   core.StringPtr(profileID),
		Attachments: attachments,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateAttachmentOptions) SetInstanceID(instanceID string) *CreateAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *CreateAttachmentOptions) SetProfileID(profileID string) *CreateAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *CreateAttachmentOptions) SetAttachments(attachments []AttachmentsPrototype) *CreateAttachmentOptions {
	_options.Attachments = attachments
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateAttachmentOptions) SetXCorrelationID(xCorrelationID string) *CreateAttachmentOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *CreateAttachmentOptions) SetXRequestID(xRequestID string) *CreateAttachmentOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAttachmentOptions) SetHeaders(param map[string]string) *CreateAttachmentOptions {
	options.Headers = param
	return options
}

// CreateCustomControlLibraryOptions : The CreateCustomControlLibrary options.
type CreateCustomControlLibraryOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The control library name.
	ControlLibraryName *string `json:"control_library_name" validate:"required"`

	// The control library description.
	ControlLibraryDescription *string `json:"control_library_description" validate:"required"`

	// The control library type.
	ControlLibraryType *string `json:"control_library_type" validate:"required"`

	// The controls.
	Controls []ControlsInControlLib `json:"controls" validate:"required"`

	// The version group label.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// The control library version.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The latest control library version.
	Latest *bool `json:"latest,omitempty"`

	// The number of controls.
	ControlsCount *int64 `json:"controls_count,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateCustomControlLibraryOptions.ControlLibraryType property.
// The control library type.
const (
	CreateCustomControlLibraryOptions_ControlLibraryType_Custom     = "custom"
	CreateCustomControlLibraryOptions_ControlLibraryType_Predefined = "predefined"
)

// NewCreateCustomControlLibraryOptions : Instantiate CreateCustomControlLibraryOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateCustomControlLibraryOptions(instanceID string, controlLibraryName string, controlLibraryDescription string, controlLibraryType string, controls []ControlsInControlLib) *CreateCustomControlLibraryOptions {
	return &CreateCustomControlLibraryOptions{
		ControlLibraryName:        core.StringPtr(controlLibraryName),
		ControlLibraryDescription: core.StringPtr(controlLibraryDescription),
		ControlLibraryType:        core.StringPtr(controlLibraryType),
		Controls:                  controls,
		InstanceID:                core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateCustomControlLibraryOptions) SetInstanceID(instanceID string) *CreateCustomControlLibraryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetControlLibraryName : Allow user to set ControlLibraryName
func (_options *CreateCustomControlLibraryOptions) SetControlLibraryName(controlLibraryName string) *CreateCustomControlLibraryOptions {
	_options.ControlLibraryName = core.StringPtr(controlLibraryName)
	return _options
}

// SetControlLibraryDescription : Allow user to set ControlLibraryDescription
func (_options *CreateCustomControlLibraryOptions) SetControlLibraryDescription(controlLibraryDescription string) *CreateCustomControlLibraryOptions {
	_options.ControlLibraryDescription = core.StringPtr(controlLibraryDescription)
	return _options
}

// SetControlLibraryType : Allow user to set ControlLibraryType
func (_options *CreateCustomControlLibraryOptions) SetControlLibraryType(controlLibraryType string) *CreateCustomControlLibraryOptions {
	_options.ControlLibraryType = core.StringPtr(controlLibraryType)
	return _options
}

// SetControls : Allow user to set Controls
func (_options *CreateCustomControlLibraryOptions) SetControls(controls []ControlsInControlLib) *CreateCustomControlLibraryOptions {
	_options.Controls = controls
	return _options
}

// SetVersionGroupLabel : Allow user to set VersionGroupLabel
func (_options *CreateCustomControlLibraryOptions) SetVersionGroupLabel(versionGroupLabel string) *CreateCustomControlLibraryOptions {
	_options.VersionGroupLabel = core.StringPtr(versionGroupLabel)
	return _options
}

// SetControlLibraryVersion : Allow user to set ControlLibraryVersion
func (_options *CreateCustomControlLibraryOptions) SetControlLibraryVersion(controlLibraryVersion string) *CreateCustomControlLibraryOptions {
	_options.ControlLibraryVersion = core.StringPtr(controlLibraryVersion)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *CreateCustomControlLibraryOptions) SetLatest(latest bool) *CreateCustomControlLibraryOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetControlsCount : Allow user to set ControlsCount
func (_options *CreateCustomControlLibraryOptions) SetControlsCount(controlsCount int64) *CreateCustomControlLibraryOptions {
	_options.ControlsCount = core.Int64Ptr(controlsCount)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateCustomControlLibraryOptions) SetXCorrelationID(xCorrelationID string) *CreateCustomControlLibraryOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *CreateCustomControlLibraryOptions) SetXRequestID(xRequestID string) *CreateCustomControlLibraryOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCustomControlLibraryOptions) SetHeaders(param map[string]string) *CreateCustomControlLibraryOptions {
	options.Headers = param
	return options
}

// CreateProfileOptions : The CreateProfile options.
type CreateProfileOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The name of the profile.
	ProfileName *string `json:"profile_name" validate:"required"`

	// The description of the profile.
	ProfileDescription *string `json:"profile_description" validate:"required"`

	// The profile type.
	ProfileType *string `json:"profile_type" validate:"required"`

	// The controls that are in the profile.
	Controls []ProfileControlsPrototype `json:"controls" validate:"required"`

	// The default parameters of the profile.
	DefaultParameters []DefaultParametersPrototype `json:"default_parameters" validate:"required"`

	// The supplied or generated value of this header is logged for a request, and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests, and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateProfileOptions.ProfileType property.
// The profile type.
const (
	CreateProfileOptions_ProfileType_Custom     = "custom"
	CreateProfileOptions_ProfileType_Predefined = "predefined"
)

// NewCreateProfileOptions : Instantiate CreateProfileOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateProfileOptions(instanceID string, profileName string, profileDescription string, profileType string, controls []ProfileControlsPrototype, defaultParameters []DefaultParametersPrototype) *CreateProfileOptions {
	return &CreateProfileOptions{
		ProfileName:        core.StringPtr(profileName),
		ProfileDescription: core.StringPtr(profileDescription),
		ProfileType:        core.StringPtr(profileType),
		InstanceID:         core.StringPtr(instanceID),
		Controls:           controls,
		DefaultParameters:  defaultParameters,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateProfileOptions) SetInstanceID(instanceID string) *CreateProfileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileName : Allow user to set ProfileName
func (_options *CreateProfileOptions) SetProfileName(profileName string) *CreateProfileOptions {
	_options.ProfileName = core.StringPtr(profileName)
	return _options
}

// SetProfileDescription : Allow user to set ProfileDescription
func (_options *CreateProfileOptions) SetProfileDescription(profileDescription string) *CreateProfileOptions {
	_options.ProfileDescription = core.StringPtr(profileDescription)
	return _options
}

// SetProfileType : Allow user to set ProfileType
func (_options *CreateProfileOptions) SetProfileType(profileType string) *CreateProfileOptions {
	_options.ProfileType = core.StringPtr(profileType)
	return _options
}

// SetControls : Allow user to set Controls
func (_options *CreateProfileOptions) SetControls(controls []ProfileControlsPrototype) *CreateProfileOptions {
	_options.Controls = controls
	return _options
}

// SetDefaultParameters : Allow user to set DefaultParameters
func (_options *CreateProfileOptions) SetDefaultParameters(defaultParameters []DefaultParametersPrototype) *CreateProfileOptions {
	_options.DefaultParameters = defaultParameters
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateProfileOptions) SetXCorrelationID(xCorrelationID string) *CreateProfileOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *CreateProfileOptions) SetXRequestID(xRequestID string) *CreateProfileOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProfileOptions) SetHeaders(param map[string]string) *CreateProfileOptions {
	options.Headers = param
	return options
}

// CreateProviderTypeInstanceOptions : The CreateProviderTypeInstance options.
type CreateProviderTypeInstanceOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The provider type instance name.
	Name *string `json:"name,omitempty"`

	// The attributes for connecting to the provider type instance.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateProviderTypeInstanceOptions : Instantiate CreateProviderTypeInstanceOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateProviderTypeInstanceOptions(instanceID string, providerTypeID string) *CreateProviderTypeInstanceOptions {
	return &CreateProviderTypeInstanceOptions{
		InstanceID:     core.StringPtr(instanceID),
		ProviderTypeID: core.StringPtr(providerTypeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateProviderTypeInstanceOptions) SetInstanceID(instanceID string) *CreateProviderTypeInstanceOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *CreateProviderTypeInstanceOptions) SetProviderTypeID(providerTypeID string) *CreateProviderTypeInstanceOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateProviderTypeInstanceOptions) SetName(name string) *CreateProviderTypeInstanceOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAttributes : Allow user to set Attributes
func (_options *CreateProviderTypeInstanceOptions) SetAttributes(attributes map[string]interface{}) *CreateProviderTypeInstanceOptions {
	_options.Attributes = attributes
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateProviderTypeInstanceOptions) SetXCorrelationID(xCorrelationID string) *CreateProviderTypeInstanceOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *CreateProviderTypeInstanceOptions) SetXRequestID(xRequestID string) *CreateProviderTypeInstanceOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProviderTypeInstanceOptions) SetHeaders(param map[string]string) *CreateProviderTypeInstanceOptions {
	options.Headers = param
	return options
}

// CreateRuleOptions : The CreateRule options.
type CreateRuleOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The rule description.
	Description *string `json:"description" validate:"required"`

	// The rule target.
	Target *Target `json:"target" validate:"required"`

	// The required configurations.
	RequiredConfig RequiredConfigIntf `json:"required_config" validate:"required"`

	// The rule type (user_defined or system_defined).
	Type *string `json:"type,omitempty"`

	// The rule version number.
	Version *string `json:"version,omitempty"`

	// The collection of import parameters.
	Import *Import `json:"import,omitempty"`

	// The list of labels that correspond to a rule.
	Labels []string `json:"labels,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateRuleOptions.Type property.
// The rule type (user_defined or system_defined).
const (
	CreateRuleOptions_Type_SystemDefined = "system_defined"
	CreateRuleOptions_Type_UserDefined   = "user_defined"
)

// NewCreateRuleOptions : Instantiate CreateRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateRuleOptions(instanceID string, description string, target *Target, requiredConfig RequiredConfigIntf) *CreateRuleOptions {
	return &CreateRuleOptions{
		InstanceID:     core.StringPtr(instanceID),
		Description:    core.StringPtr(description),
		Target:         target,
		RequiredConfig: requiredConfig,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateRuleOptions) SetInstanceID(instanceID string) *CreateRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateRuleOptions) SetDescription(description string) *CreateRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *CreateRuleOptions) SetTarget(target *Target) *CreateRuleOptions {
	_options.Target = target
	return _options
}

// SetRequiredConfig : Allow user to set RequiredConfig
func (_options *CreateRuleOptions) SetRequiredConfig(requiredConfig RequiredConfigIntf) *CreateRuleOptions {
	_options.RequiredConfig = requiredConfig
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateRuleOptions) SetType(typeVar string) *CreateRuleOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *CreateRuleOptions) SetVersion(version string) *CreateRuleOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetImport : Allow user to set Import
func (_options *CreateRuleOptions) SetImport(importVar *Import) *CreateRuleOptions {
	_options.Import = importVar
	return _options
}

// SetLabels : Allow user to set Labels
func (_options *CreateRuleOptions) SetLabels(labels []string) *CreateRuleOptions {
	_options.Labels = labels
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateRuleOptions) SetXCorrelationID(xCorrelationID string) *CreateRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *CreateRuleOptions) SetXRequestID(xRequestID string) *CreateRuleOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRuleOptions) SetHeaders(param map[string]string) *CreateRuleOptions {
	options.Headers = param
	return options
}

// CreateScanOptions : The CreateScan options.
type CreateScanOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The attachment ID of a profile.
	AttachmentID *string `json:"attachment_id" validate:"required"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateScanOptions : Instantiate CreateScanOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateScanOptions(instanceID string, attachmentID string) *CreateScanOptions {
	return &CreateScanOptions{
		InstanceID:   core.StringPtr(instanceID),
		AttachmentID: core.StringPtr(attachmentID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateScanOptions) SetInstanceID(instanceID string) *CreateScanOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *CreateScanOptions) SetAttachmentID(attachmentID string) *CreateScanOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateScanOptions) SetXCorrelationID(xCorrelationID string) *CreateScanOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *CreateScanOptions) SetXRequestID(xRequestID string) *CreateScanOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateScanOptions) SetHeaders(param map[string]string) *CreateScanOptions {
	options.Headers = param
	return options
}

// DefaultParametersPrototype : The control details of a profile.
type DefaultParametersPrototype struct {
	// The type of the implementation.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The implementation ID of the parameter.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The parameter name.
	ParameterName *string `json:"parameter_name,omitempty"`

	// The default value of the parameter.
	ParameterDefaultValue *string `json:"parameter_default_value,omitempty"`

	// The parameter display name.
	ParameterDisplayName *string `json:"parameter_display_name,omitempty"`

	// The parameter type.
	ParameterType *string `json:"parameter_type,omitempty"`
}

// Constants associated with the DefaultParametersPrototype.ParameterType property.
// The parameter type.
const (
	DefaultParametersPrototype_ParameterType_Boolean    = "boolean"
	DefaultParametersPrototype_ParameterType_General    = "general"
	DefaultParametersPrototype_ParameterType_IpList     = "ip_list"
	DefaultParametersPrototype_ParameterType_Numeric    = "numeric"
	DefaultParametersPrototype_ParameterType_String     = "string"
	DefaultParametersPrototype_ParameterType_StringList = "string_list"
	DefaultParametersPrototype_ParameterType_Timestamp  = "timestamp"
)

// UnmarshalDefaultParametersPrototype unmarshals an instance of DefaultParametersPrototype from the specified map of raw messages.
func UnmarshalDefaultParametersPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefaultParametersPrototype)
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_name", &obj.ParameterName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_default_value", &obj.ParameterDefaultValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_display_name", &obj.ParameterDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_type", &obj.ParameterType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteCustomControlLibraryOptions : The DeleteCustomControlLibrary options.
type DeleteCustomControlLibraryOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The control library ID.
	ControlLibrariesID *string `json:"control_libraries_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCustomControlLibraryOptions : Instantiate DeleteCustomControlLibraryOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteCustomControlLibraryOptions(instanceID string, controlLibrariesID string) *DeleteCustomControlLibraryOptions {
	return &DeleteCustomControlLibraryOptions{
		InstanceID:         core.StringPtr(instanceID),
		ControlLibrariesID: core.StringPtr(controlLibrariesID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteCustomControlLibraryOptions) SetInstanceID(instanceID string) *DeleteCustomControlLibraryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetControlLibrariesID : Allow user to set ControlLibrariesID
func (_options *DeleteCustomControlLibraryOptions) SetControlLibrariesID(controlLibrariesID string) *DeleteCustomControlLibraryOptions {
	_options.ControlLibrariesID = core.StringPtr(controlLibrariesID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteCustomControlLibraryOptions) SetXCorrelationID(xCorrelationID string) *DeleteCustomControlLibraryOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *DeleteCustomControlLibraryOptions) SetXRequestID(xRequestID string) *DeleteCustomControlLibraryOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCustomControlLibraryOptions) SetHeaders(param map[string]string) *DeleteCustomControlLibraryOptions {
	options.Headers = param
	return options
}

// DeleteCustomProfileOptions : The DeleteCustomProfile options.
type DeleteCustomProfileOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCustomProfileOptions : Instantiate DeleteCustomProfileOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteCustomProfileOptions(instanceID string, profileID string) *DeleteCustomProfileOptions {
	return &DeleteCustomProfileOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteCustomProfileOptions) SetInstanceID(instanceID string) *DeleteCustomProfileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *DeleteCustomProfileOptions) SetProfileID(profileID string) *DeleteCustomProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteCustomProfileOptions) SetXCorrelationID(xCorrelationID string) *DeleteCustomProfileOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *DeleteCustomProfileOptions) SetXRequestID(xRequestID string) *DeleteCustomProfileOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCustomProfileOptions) SetHeaders(param map[string]string) *DeleteCustomProfileOptions {
	options.Headers = param
	return options
}

// DeleteProfileAttachmentOptions : The DeleteProfileAttachment options.
type DeleteProfileAttachmentOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The attachment ID.
	AttachmentID *string `json:"attachment_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteProfileAttachmentOptions : Instantiate DeleteProfileAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteProfileAttachmentOptions(instanceID string, attachmentID string, profileID string) *DeleteProfileAttachmentOptions {
	return &DeleteProfileAttachmentOptions{
		InstanceID:   core.StringPtr(instanceID),
		AttachmentID: core.StringPtr(attachmentID),
		ProfileID:    core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteProfileAttachmentOptions) SetInstanceID(instanceID string) *DeleteProfileAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *DeleteProfileAttachmentOptions) SetAttachmentID(attachmentID string) *DeleteProfileAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *DeleteProfileAttachmentOptions) SetProfileID(profileID string) *DeleteProfileAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteProfileAttachmentOptions) SetXCorrelationID(xCorrelationID string) *DeleteProfileAttachmentOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *DeleteProfileAttachmentOptions) SetXRequestID(xRequestID string) *DeleteProfileAttachmentOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProfileAttachmentOptions) SetHeaders(param map[string]string) *DeleteProfileAttachmentOptions {
	options.Headers = param
	return options
}

// DeleteProviderTypeInstanceOptions : The DeleteProviderTypeInstance options.
type DeleteProviderTypeInstanceOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The provider type instance ID.
	ProviderTypeInstanceID *string `json:"provider_type_instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteProviderTypeInstanceOptions : Instantiate DeleteProviderTypeInstanceOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteProviderTypeInstanceOptions(instanceID string, providerTypeID string, providerTypeInstanceID string) *DeleteProviderTypeInstanceOptions {
	return &DeleteProviderTypeInstanceOptions{
		InstanceID:             core.StringPtr(instanceID),
		ProviderTypeID:         core.StringPtr(providerTypeID),
		ProviderTypeInstanceID: core.StringPtr(providerTypeInstanceID),
	}
}

// // SetInstanceID : Allow user to set InstanceID
func (_options *DeleteProviderTypeInstanceOptions) SetInstanceID(instanceID string) *DeleteProviderTypeInstanceOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *DeleteProviderTypeInstanceOptions) SetProviderTypeID(providerTypeID string) *DeleteProviderTypeInstanceOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetProviderTypeInstanceID : Allow user to set ProviderTypeInstanceID
func (_options *DeleteProviderTypeInstanceOptions) SetProviderTypeInstanceID(providerTypeInstanceID string) *DeleteProviderTypeInstanceOptions {
	_options.ProviderTypeInstanceID = core.StringPtr(providerTypeInstanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteProviderTypeInstanceOptions) SetXCorrelationID(xCorrelationID string) *DeleteProviderTypeInstanceOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *DeleteProviderTypeInstanceOptions) SetXRequestID(xRequestID string) *DeleteProviderTypeInstanceOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProviderTypeInstanceOptions) SetHeaders(param map[string]string) *DeleteProviderTypeInstanceOptions {
	options.Headers = param
	return options
}

// DeleteRuleOptions : The DeleteRule options.
type DeleteRuleOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the corresponding rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRuleOptions : Instantiate DeleteRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteRuleOptions(instanceID string, ruleID string) *DeleteRuleOptions {
	return &DeleteRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteRuleOptions) SetInstanceID(instanceID string) *DeleteRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteRuleOptions) SetRuleID(ruleID string) *DeleteRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteRuleOptions) SetXCorrelationID(xCorrelationID string) *DeleteRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *DeleteRuleOptions) SetXRequestID(xRequestID string) *DeleteRuleOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRuleOptions) SetHeaders(param map[string]string) *DeleteRuleOptions {
	options.Headers = param
	return options
}

// EvalDetails : The evaluation details.
type EvalDetails struct {
	// The evaluation properties.
	Properties []Property `json:"properties,omitempty"`
}

// UnmarshalEvalDetails unmarshals an instance of EvalDetails from the specified map of raw messages.
func UnmarshalEvalDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EvalDetails)
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EvalStats : The evaluation stats.
type EvalStats struct {
	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of passed evaluations.
	PassCount *int64 `json:"pass_count,omitempty"`

	// The number of failed evaluations.
	FailureCount *int64 `json:"failure_count,omitempty"`

	// The number of evaluations that started, but did not finish, and ended with errors.
	ErrorCount *int64 `json:"error_count,omitempty"`

	// The total number of completed evaluations.
	CompletedCount *int64 `json:"completed_count,omitempty"`
}

// Constants associated with the EvalStats.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	EvalStats_Status_Compliant              = "compliant"
	EvalStats_Status_NotCompliant           = "not_compliant"
	EvalStats_Status_UnableToPerform        = "unable_to_perform"
	EvalStats_Status_UserEvaluationRequired = "user_evaluation_required"
)

// UnmarshalEvalStats unmarshals an instance of EvalStats from the specified map of raw messages.
func UnmarshalEvalStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EvalStats)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pass_count", &obj.PassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_count", &obj.FailureCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "completed_count", &obj.CompletedCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Evaluation : The evaluation of a control specification assessment.
type Evaluation struct {
	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The ID of the report that is associated to the evaluation.
	ReportID *string `json:"report_id,omitempty"`

	// The control ID.
	ControlID *string `json:"control_id,omitempty"`

	// The component ID.
	ComponentID *string `json:"component_id,omitempty"`

	// The control specification assessment.
	Assessment *Assessment `json:"assessment,omitempty"`

	// The time when the evaluation was made.
	EvaluateTime *string `json:"evaluate_time,omitempty"`

	// The evaluation target.
	Target *TargetInfo `json:"target,omitempty"`

	// The allowed values of an evaluation status.
	Status *string `json:"status,omitempty"`

	// The reason for the evaluation failure.
	Reason *string `json:"reason,omitempty"`

	// The evaluation details.
	Details *EvalDetails `json:"details,omitempty"`
}

// Constants associated with the Evaluation.Status property.
// The allowed values of an evaluation status.
const (
	Evaluation_Status_Error   = "error"
	Evaluation_Status_Failure = "failure"
	Evaluation_Status_Pass    = "pass"
	Evaluation_Status_Skipped = "skipped"
)

// UnmarshalEvaluation unmarshals an instance of Evaluation from the specified map of raw messages.
func UnmarshalEvaluation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Evaluation)
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_id", &obj.ControlID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_id", &obj.ComponentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "assessment", &obj.Assessment, UnmarshalAssessment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "evaluate_time", &obj.EvaluateTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalTargetInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "details", &obj.Details, UnmarshalEvalDetails)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EvaluationPage : The page of assessment evaluations.
type EvaluationPage struct {
	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The requested page limi.t.
	Limit *int64 `json:"limit" validate:"required"`

	// The token of the next page, when it's present.
	Start *string `json:"start,omitempty"`

	// The page reference.
	First *PageHRef `json:"first" validate:"required"`

	// The page reference.
	Next *PageHRef `json:"next,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The list of evaluations that are on the page.
	Evaluations []Evaluation `json:"evaluations,omitempty"`
}

// UnmarshalEvaluationPage unmarshals an instance of EvaluationPage from the specified map of raw messages.
func UnmarshalEvaluationPage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EvaluationPage)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRef)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRef)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "evaluations", &obj.Evaluations, UnmarshalEvaluation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *EvaluationPage) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	start, err := core.GetQueryParam(resp.Next.Href, "start")
	if err != nil || start == nil {
		return nil, err
	}
	return start, nil
}

// EventNotifications : The Event Notifications settings.
type EventNotifications struct {
	// The Event Notifications instance CRN.
	InstanceCrn *string `json:"instance_crn,omitempty"`

	// The date when the Event Notifications connection was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The connected Security and Compliance Center instance CRN.
	SourceID *string `json:"source_id,omitempty"`

	// The description of the source of the Event Notifications.
	SourceDescription *string `json:"source_description,omitempty"`

	// The name of the source of the Event Notifications.
	SourceName *string `json:"source_name,omitempty"`
}

// UnmarshalEventNotifications unmarshals an instance of EventNotifications from the specified map of raw messages.
func UnmarshalEventNotifications(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EventNotifications)
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_id", &obj.SourceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_description", &obj.SourceDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_name", &obj.SourceName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FailedControls : The failed controls.
type FailedControls struct {
	// The threshold limit.
	ThresholdLimit *int64 `json:"threshold_limit,omitempty"`

	// The failed control IDs.
	FailedControlIds []string `json:"failed_control_ids,omitempty"`
}

// UnmarshalFailedControls unmarshals an instance of FailedControls from the specified map of raw messages.
func UnmarshalFailedControls(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FailedControls)
	err = core.UnmarshalPrimitive(m, "threshold_limit", &obj.ThresholdLimit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed_control_ids", &obj.FailedControlIds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetControlLibraryOptions : The GetControlLibrary options.
type GetControlLibraryOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The control library ID.
	ControlLibrariesID *string `json:"control_libraries_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetControlLibraryOptions : Instantiate GetControlLibraryOptions
func (*SecurityAndComplianceCenterApiV3) NewGetControlLibraryOptions(instanceID string, controlLibrariesID string) *GetControlLibraryOptions {
	return &GetControlLibraryOptions{
		InstanceID:         core.StringPtr(instanceID),
		ControlLibrariesID: core.StringPtr(controlLibrariesID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetControlLibraryOptions) SetInstanceID(instanceID string) *GetControlLibraryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetControlLibrariesID : Allow user to set ControlLibrariesID
func (_options *GetControlLibraryOptions) SetControlLibrariesID(controlLibrariesID string) *GetControlLibraryOptions {
	_options.ControlLibrariesID = core.StringPtr(controlLibrariesID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetControlLibraryOptions) SetXCorrelationID(xCorrelationID string) *GetControlLibraryOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetControlLibraryOptions) SetXRequestID(xRequestID string) *GetControlLibraryOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetControlLibraryOptions) SetHeaders(param map[string]string) *GetControlLibraryOptions {
	options.Headers = param
	return options
}

// GetLatestReportsOptions : The GetLatestReports options.
type GetLatestReportsOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// This field sorts results by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLatestReportsOptions : Instantiate GetLatestReportsOptions
func (*SecurityAndComplianceCenterApiV3) NewGetLatestReportsOptions(instanceID string) *GetLatestReportsOptions {
	return &GetLatestReportsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetLatestReportsOptions) SetInstanceID(instanceID string) *GetLatestReportsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetLatestReportsOptions) SetXCorrelationID(xCorrelationID string) *GetLatestReportsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetLatestReportsOptions) SetXRequestID(xRequestID string) *GetLatestReportsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *GetLatestReportsOptions) SetSort(sort string) *GetLatestReportsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLatestReportsOptions) SetHeaders(param map[string]string) *GetLatestReportsOptions {
	options.Headers = param
	return options
}

// GetProfileAttachmentOptions : The GetProfileAttachment options.
type GetProfileAttachmentOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The attachment ID.
	AttachmentID *string `json:"attachment_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProfileAttachmentOptions : Instantiate GetProfileAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProfileAttachmentOptions(instanceID string, attachmentID string, profileID string) *GetProfileAttachmentOptions {
	return &GetProfileAttachmentOptions{
		InstanceID:   core.StringPtr(instanceID),
		AttachmentID: core.StringPtr(attachmentID),
		ProfileID:    core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetProfileAttachmentOptions) SetInstanceID(instanceID string) *GetProfileAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *GetProfileAttachmentOptions) SetAttachmentID(attachmentID string) *GetProfileAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetProfileAttachmentOptions) SetProfileID(profileID string) *GetProfileAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetProfileAttachmentOptions) SetXCorrelationID(xCorrelationID string) *GetProfileAttachmentOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetProfileAttachmentOptions) SetXRequestID(xRequestID string) *GetProfileAttachmentOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileAttachmentOptions) SetHeaders(param map[string]string) *GetProfileAttachmentOptions {
	options.Headers = param
	return options
}

// GetProfileOptions : The GetProfile options.
type GetProfileOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProfileOptions : Instantiate GetProfileOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProfileOptions(instanceID string, profileID string) *GetProfileOptions {
	return &GetProfileOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetProfileOptions) SetInstanceID(instanceID string) *GetProfileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetProfileOptions) SetProfileID(profileID string) *GetProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetProfileOptions) SetXCorrelationID(xCorrelationID string) *GetProfileOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetProfileOptions) SetXRequestID(xRequestID string) *GetProfileOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileOptions) SetHeaders(param map[string]string) *GetProfileOptions {
	options.Headers = param
	return options
}

// GetProviderTypeByIdOptions : The GetProviderTypeByID options.
type GetProviderTypeByIdOptions struct {
	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProviderTypeByIdOptions : Instantiate GetProviderTypeByIdOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProviderTypeByIdOptions(providerTypeID string) *GetProviderTypeByIdOptions {
	return &GetProviderTypeByIdOptions{
		ProviderTypeID: core.StringPtr(providerTypeID),
	}
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *GetProviderTypeByIdOptions) SetProviderTypeID(providerTypeID string) *GetProviderTypeByIdOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetProviderTypeByIdOptions) SetXCorrelationID(xCorrelationID string) *GetProviderTypeByIdOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetProviderTypeByIdOptions) SetXRequestID(xRequestID string) *GetProviderTypeByIdOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProviderTypeByIdOptions) SetHeaders(param map[string]string) *GetProviderTypeByIdOptions {
	options.Headers = param
	return options
}

// GetProviderTypeInstanceOptions : The GetProviderTypeInstance options.
type GetProviderTypeInstanceOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The provider type instance ID.
	ProviderTypeInstanceID *string `json:"provider_type_instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProviderTypeInstanceOptions : Instantiate GetProviderTypeInstanceOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProviderTypeInstanceOptions(instanceID string, providerTypeID string, providerTypeInstanceID string) *GetProviderTypeInstanceOptions {
	return &GetProviderTypeInstanceOptions{
		InstanceID:             core.StringPtr(instanceID),
		ProviderTypeID:         core.StringPtr(providerTypeID),
		ProviderTypeInstanceID: core.StringPtr(providerTypeInstanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetProviderTypeInstanceOptions) SetInstanceID(instanceID string) *GetProviderTypeInstanceOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *GetProviderTypeInstanceOptions) SetProviderTypeID(providerTypeID string) *GetProviderTypeInstanceOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetProviderTypeInstanceID : Allow user to set ProviderTypeInstanceID
func (_options *GetProviderTypeInstanceOptions) SetProviderTypeInstanceID(providerTypeInstanceID string) *GetProviderTypeInstanceOptions {
	_options.ProviderTypeInstanceID = core.StringPtr(providerTypeInstanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetProviderTypeInstanceOptions) SetXCorrelationID(xCorrelationID string) *GetProviderTypeInstanceOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetProviderTypeInstanceOptions) SetXRequestID(xRequestID string) *GetProviderTypeInstanceOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProviderTypeInstanceOptions) SetHeaders(param map[string]string) *GetProviderTypeInstanceOptions {
	options.Headers = param
	return options
}

// GetProviderTypesInstancesOptions : The GetProviderTypesInstances options.
type GetProviderTypesInstancesOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`
	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProviderTypesInstancesOptions : Instantiate GetProviderTypesInstancesOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProviderTypesInstancesOptions(instanceID string) *GetProviderTypesInstancesOptions {
	return &GetProviderTypesInstancesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstance : Allow user to set Instantiate
func (_options *GetProviderTypesInstancesOptions) SetInstanceID(instanceID string) *GetProviderTypesInstancesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetProviderTypesInstancesOptions) SetXCorrelationID(xCorrelationID string) *GetProviderTypesInstancesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetProviderTypesInstancesOptions) SetXRequestID(xRequestID string) *GetProviderTypesInstancesOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProviderTypesInstancesOptions) SetHeaders(param map[string]string) *GetProviderTypesInstancesOptions {
	options.Headers = param
	return options
}

// GetReportControlsOptions : The GetReportControls options.
type GetReportControlsOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The ID of the control.
	ControlID *string `json:"control_id,omitempty"`

	// The name of the control.
	ControlName *string `json:"control_name,omitempty"`

	// The description of the control.
	ControlDescription *string `json:"control_description,omitempty"`

	// A control category value.
	ControlCategory *string `json:"control_category,omitempty"`

	// The compliance status value.
	Status *string `json:"status,omitempty"`

	// This field sorts controls by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetReportControlsOptions.Status property.
// The compliance status value.
const (
	GetReportControlsOptions_Status_Compliant              = "compliant"
	GetReportControlsOptions_Status_NotCompliant           = "not_compliant"
	GetReportControlsOptions_Status_UnableToPerform        = "unable_to_perform"
	GetReportControlsOptions_Status_UserEvaluationRequired = "user_evaluation_required"
)

// Constants associated with the GetReportControlsOptions.Sort property.
// This field sorts controls by using a valid sort field. To learn more, see
// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
const (
	GetReportControlsOptions_Sort_ControlCategory = "control_category"
	GetReportControlsOptions_Sort_ControlName     = "control_name"
	GetReportControlsOptions_Sort_Status          = "status"
)

// NewGetReportControlsOptions : Instantiate GetReportControlsOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportControlsOptions(instanceID string, reportID string) *GetReportControlsOptions {
	return &GetReportControlsOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportControlsOptions) SetInstanceID(instanceID string) *GetReportControlsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportControlsOptions) SetReportID(reportID string) *GetReportControlsOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetReportControlsOptions) SetXCorrelationID(xCorrelationID string) *GetReportControlsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetReportControlsOptions) SetXRequestID(xRequestID string) *GetReportControlsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetControlID : Allow user to set ControlID
func (_options *GetReportControlsOptions) SetControlID(controlID string) *GetReportControlsOptions {
	_options.ControlID = core.StringPtr(controlID)
	return _options
}

// SetControlName : Allow user to set ControlName
func (_options *GetReportControlsOptions) SetControlName(controlName string) *GetReportControlsOptions {
	_options.ControlName = core.StringPtr(controlName)
	return _options
}

// SetControlDescription : Allow user to set ControlDescription
func (_options *GetReportControlsOptions) SetControlDescription(controlDescription string) *GetReportControlsOptions {
	_options.ControlDescription = core.StringPtr(controlDescription)
	return _options
}

// SetControlCategory : Allow user to set ControlCategory
func (_options *GetReportControlsOptions) SetControlCategory(controlCategory string) *GetReportControlsOptions {
	_options.ControlCategory = core.StringPtr(controlCategory)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *GetReportControlsOptions) SetStatus(status string) *GetReportControlsOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *GetReportControlsOptions) SetSort(sort string) *GetReportControlsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportControlsOptions) SetHeaders(param map[string]string) *GetReportControlsOptions {
	options.Headers = param
	return options
}

// GetReportEvaluationOptions : The GetReportEvaluation options.
type GetReportEvaluationOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The indication of whether report summary metadata must be excluded.
	ExcludeSummary *bool `json:"exclude_summary,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportEvaluationOptions : Instantiate GetReportEvaluationOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportEvaluationOptions(instanceID string, reportID string) *GetReportEvaluationOptions {
	return &GetReportEvaluationOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportEvaluationOptions) SetInstanceID(instanceID string) *GetReportEvaluationOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportEvaluationOptions) SetReportID(reportID string) *GetReportEvaluationOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetReportEvaluationOptions) SetXCorrelationID(xCorrelationID string) *GetReportEvaluationOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetReportEvaluationOptions) SetXRequestID(xRequestID string) *GetReportEvaluationOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetExcludeSummary : Allow user to set ExcludeSummary
func (_options *GetReportEvaluationOptions) SetExcludeSummary(excludeSummary bool) *GetReportEvaluationOptions {
	_options.ExcludeSummary = core.BoolPtr(excludeSummary)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportEvaluationOptions) SetHeaders(param map[string]string) *GetReportEvaluationOptions {
	options.Headers = param
	return options
}

// GetReportOptions : The GetReport options.
type GetReportOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportOptions : Instantiate GetReportOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportOptions(instanceID string, reportID string) *GetReportOptions {
	return &GetReportOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportOptions) SetInstanceID(instanceID string) *GetReportOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportOptions) SetReportID(reportID string) *GetReportOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetReportOptions) SetXCorrelationID(xCorrelationID string) *GetReportOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetReportOptions) SetXRequestID(xRequestID string) *GetReportOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportOptions) SetHeaders(param map[string]string) *GetReportOptions {
	options.Headers = param
	return options
}

// GetReportRuleOptions : The GetReportRule options.
type GetReportRuleOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of a rule in a report.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportRuleOptions : Instantiate GetReportRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportRuleOptions(instanceID string, reportID string, ruleID string) *GetReportRuleOptions {
	return &GetReportRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportRuleOptions) SetInstanceID(instanceID string) *GetReportRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportRuleOptions) SetReportID(reportID string) *GetReportRuleOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *GetReportRuleOptions) SetRuleID(ruleID string) *GetReportRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetReportRuleOptions) SetXCorrelationID(xCorrelationID string) *GetReportRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetReportRuleOptions) SetXRequestID(xRequestID string) *GetReportRuleOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportRuleOptions) SetHeaders(param map[string]string) *GetReportRuleOptions {
	options.Headers = param
	return options
}

// GetReportSummaryOptions : The GetReportSummary options.
type GetReportSummaryOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportSummaryOptions : Instantiate GetReportSummaryOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportSummaryOptions(instanceID string, reportID string) *GetReportSummaryOptions {
	return &GetReportSummaryOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportSummaryOptions) SetInstanceID(instanceID string) *GetReportSummaryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportSummaryOptions) SetReportID(reportID string) *GetReportSummaryOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetReportSummaryOptions) SetXCorrelationID(xCorrelationID string) *GetReportSummaryOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetReportSummaryOptions) SetXRequestID(xRequestID string) *GetReportSummaryOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportSummaryOptions) SetHeaders(param map[string]string) *GetReportSummaryOptions {
	options.Headers = param
	return options
}

// GetReportTagsOptions : The GetReportTags options.
type GetReportTagsOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportTagsOptions : Instantiate GetReportTagsOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportTagsOptions(instanceID string, reportID string) *GetReportTagsOptions {
	return &GetReportTagsOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportTagsOptions) SetInstanceID(instanceID string) *GetReportTagsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportTagsOptions) SetReportID(reportID string) *GetReportTagsOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetReportTagsOptions) SetXCorrelationID(xCorrelationID string) *GetReportTagsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetReportTagsOptions) SetXRequestID(xRequestID string) *GetReportTagsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportTagsOptions) SetHeaders(param map[string]string) *GetReportTagsOptions {
	options.Headers = param
	return options
}

// GetReportViolationsDriftOptions : The GetReportViolationsDrift options.
type GetReportViolationsDriftOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The duration of the `scan_time` timestamp in number of days.
	ScanTimeDuration *int64 `json:"scan_time_duration,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportViolationsDriftOptions : Instantiate GetReportViolationsDriftOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportViolationsDriftOptions(instanceID string, reportID string) *GetReportViolationsDriftOptions {
	return &GetReportViolationsDriftOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportViolationsDriftOptions) SetInstanceID(instanceID string) *GetReportViolationsDriftOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportViolationsDriftOptions) SetReportID(reportID string) *GetReportViolationsDriftOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetReportViolationsDriftOptions) SetXCorrelationID(xCorrelationID string) *GetReportViolationsDriftOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetReportViolationsDriftOptions) SetXRequestID(xRequestID string) *GetReportViolationsDriftOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetScanTimeDuration : Allow user to set ScanTimeDuration
func (_options *GetReportViolationsDriftOptions) SetScanTimeDuration(scanTimeDuration int64) *GetReportViolationsDriftOptions {
	_options.ScanTimeDuration = core.Int64Ptr(scanTimeDuration)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportViolationsDriftOptions) SetHeaders(param map[string]string) *GetReportViolationsDriftOptions {
	options.Headers = param
	return options
}

// GetRuleOptions : The GetRule options.
type GetRuleOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the corresponding rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRuleOptions : Instantiate GetRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewGetRuleOptions(instanceID string, ruleID string) *GetRuleOptions {
	return &GetRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetRuleOptions) SetInstanceID(instanceID string) *GetRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *GetRuleOptions) SetRuleID(ruleID string) *GetRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetRuleOptions) SetXCorrelationID(xCorrelationID string) *GetRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetRuleOptions) SetXRequestID(xRequestID string) *GetRuleOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRuleOptions) SetHeaders(param map[string]string) *GetRuleOptions {
	options.Headers = param
	return options
}

// GetSettingsOptions : The GetSettings options.
type GetSettingsOptions struct {
	// The ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request, and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The supplied or generated value of this header is logged for a request, and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSettingsOptions : Instantiate GetSettingsOptions
func (*SecurityAndComplianceCenterApiV3) NewGetSettingsOptions(instanceID string) *GetSettingsOptions {
	return &GetSettingsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetSettingsOptions) SetInstanceID(instanceID string) *GetSettingsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetSettingsOptions) SetXCorrelationID(xCorrelationID string) *GetSettingsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *GetSettingsOptions) SetXRequestID(xRequestID string) *GetSettingsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSettingsOptions) SetHeaders(param map[string]string) *GetSettingsOptions {
	options.Headers = param
	return options
}

// Implementation : The implementation details of a control library.
type Implementation struct {
	// The assessment ID.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The assessment method.
	AssessmentMethod *string `json:"assessment_method,omitempty"`

	// The assessment type.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The assessment description.
	AssessmentDescription *string `json:"assessment_description,omitempty"`

	// The parameter count.
	ParameterCount *int64 `json:"parameter_count,omitempty"`

	// The parameters.
	Parameters []ParameterInfo `json:"parameters,omitempty"`
}

// UnmarshalImplementation unmarshals an instance of Implementation from the specified map of raw messages.
func UnmarshalImplementation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Implementation)
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_method", &obj.AssessmentMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_description", &obj.AssessmentDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_count", &obj.ParameterCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "parameters", &obj.Parameters, UnmarshalParameterInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Import : The collection of import parameters.
type Import struct {
	// The list of import parameters.
	Parameters []Parameter `json:"parameters,omitempty"`
}

// UnmarshalImport unmarshals an instance of Import from the specified map of raw messages.
func UnmarshalImport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Import)
	err = core.UnmarshalModel(m, "parameters", &obj.Parameters, UnmarshalParameter)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LabelType : The label that is associated with the provider type.
type LabelType struct {
	// The text of the label.
	Text *string `json:"text,omitempty"`

	// The text to be shown when user hover overs the label.
	Tip *string `json:"tip,omitempty"`
}

// UnmarshalLabelType unmarshals an instance of LabelType from the specified map of raw messages.
func UnmarshalLabelType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LabelType)
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tip", &obj.Tip)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LastScan : The details of the last scan of an attachment.
type LastScan struct {
	// The ID of the last scan of an attachment.
	ID *string `json:"id,omitempty"`

	// The status of the last scan of an attachment.
	Status *string `json:"status,omitempty"`

	// The time when the last scan started.
	Time *strfmt.DateTime `json:"time,omitempty"`
}

// Constants associated with the LastScan.Status property.
// The status of the last scan of an attachment.
const (
	LastScan_Status_Completed  = "completed"
	LastScan_Status_InProgress = "in_progress"
)

// UnmarshalLastScan unmarshals an instance of LastScan from the specified map of raw messages.
func UnmarshalLastScan(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LastScan)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "time", &obj.Time)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAttachmentsAccountOptions : The ListAttachmentsAccount options.
type ListAttachmentsAccountOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The indication of how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAttachmentsAccountOptions : Instantiate ListAttachmentsAccountOptions
func (*SecurityAndComplianceCenterApiV3) NewListAttachmentsAccountOptions(instanceID string) *ListAttachmentsAccountOptions {
	return &ListAttachmentsAccountOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListAttachmentsAccountOptions) SetXCorrelationID(xCorrelationID string) *ListAttachmentsAccountOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListAttachmentsAccountOptions) SetXRequestID(xRequestID string) *ListAttachmentsAccountOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAttachmentsAccountOptions) SetLimit(limit int64) *ListAttachmentsAccountOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListAttachmentsAccountOptions) SetStart(start string) *ListAttachmentsAccountOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAttachmentsAccountOptions) SetHeaders(param map[string]string) *ListAttachmentsAccountOptions {
	options.Headers = param
	return options
}

// ListAttachmentsOptions : The ListAttachments options.
type ListAttachmentsOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The indication of how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAttachmentsOptions : Instantiate ListAttachmentsOptions
func (*SecurityAndComplianceCenterApiV3) NewListAttachmentsOptions(instanceID string, profileID string) *ListAttachmentsOptions {
	return &ListAttachmentsOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListAttachmentsOptions) SetInstanceID(instanceID string) *ListAttachmentsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ListAttachmentsOptions) SetProfileID(profileID string) *ListAttachmentsOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListAttachmentsOptions) SetXCorrelationID(xCorrelationID string) *ListAttachmentsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListAttachmentsOptions) SetXRequestID(xRequestID string) *ListAttachmentsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAttachmentsOptions) SetLimit(limit int64) *ListAttachmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListAttachmentsOptions) SetStart(start string) *ListAttachmentsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAttachmentsOptions) SetHeaders(param map[string]string) *ListAttachmentsOptions {
	options.Headers = param
	return options
}

// ListControlLibrariesOptions : The ListControlLibraries options.
type ListControlLibrariesOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The field that indicates how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// The field that indicate how you want the resources to be filtered by.
	ControlLibraryType *string `json:"control_library_type,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListControlLibrariesOptions : Instantiate ListControlLibrariesOptions
func (*SecurityAndComplianceCenterApiV3) NewListControlLibrariesOptions(instanceID string) *ListControlLibrariesOptions {
	return &ListControlLibrariesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListControlLibrariesOptions) SetInstanceID(instanceID string) *ListControlLibrariesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListControlLibrariesOptions) SetXCorrelationID(xCorrelationID string) *ListControlLibrariesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListControlLibrariesOptions) SetXRequestID(xRequestID string) *ListControlLibrariesOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListControlLibrariesOptions) SetLimit(limit int64) *ListControlLibrariesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetControlLibraryType : Allow user to set ControlLibraryType
func (_options *ListControlLibrariesOptions) SetControlLibraryType(controlLibraryType string) *ListControlLibrariesOptions {
	_options.ControlLibraryType = core.StringPtr(controlLibraryType)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListControlLibrariesOptions) SetStart(start string) *ListControlLibrariesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListControlLibrariesOptions) SetHeaders(param map[string]string) *ListControlLibrariesOptions {
	options.Headers = param
	return options
}

// ListProfilesOptions : The ListProfiles options.
type ListProfilesOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request, and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests, and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The indication of how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// The field that indicate how you want the resources to be filtered by.
	ProfileType *string `json:"profile_type,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProfilesOptions : Instantiate ListProfilesOptions
func (*SecurityAndComplianceCenterApiV3) NewListProfilesOptions(instanceID string) *ListProfilesOptions {
	return &ListProfilesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListProfilesOptions) SetInstanceID(instanceID string) *ListProfilesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListProfilesOptions) SetXCorrelationID(xCorrelationID string) *ListProfilesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListProfilesOptions) SetXRequestID(xRequestID string) *ListProfilesOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListProfilesOptions) SetLimit(limit int64) *ListProfilesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetProfileType : Allow user to set ProfileType
func (_options *ListProfilesOptions) SetProfileType(profileType string) *ListProfilesOptions {
	_options.ProfileType = core.StringPtr(profileType)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListProfilesOptions) SetStart(start string) *ListProfilesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProfilesOptions) SetHeaders(param map[string]string) *ListProfilesOptions {
	options.Headers = param
	return options
}

// ListProviderTypeInstancesOptions : The ListProviderTypeInstances options.
type ListProviderTypeInstancesOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProviderTypeInstancesOptions : Instantiate ListProviderTypeInstancesOptions
func (*SecurityAndComplianceCenterApiV3) NewListProviderTypeInstancesOptions(instanceID string, providerTypeID string) *ListProviderTypeInstancesOptions {
	return &ListProviderTypeInstancesOptions{
		InstanceID:     core.StringPtr(instanceID),
		ProviderTypeID: core.StringPtr(providerTypeID),
	}
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *ListProviderTypeInstancesOptions) SetProviderTypeID(providerTypeID string) *ListProviderTypeInstancesOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListProviderTypeInstancesOptions) SetXCorrelationID(xCorrelationID string) *ListProviderTypeInstancesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListProviderTypeInstancesOptions) SetXRequestID(xRequestID string) *ListProviderTypeInstancesOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProviderTypeInstancesOptions) SetHeaders(param map[string]string) *ListProviderTypeInstancesOptions {
	options.Headers = param
	return options
}

// ListProviderTypesOptions : The ListProviderTypes options.
type ListProviderTypesOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProviderTypesOptions : Instantiate ListProviderTypesOptions
func (*SecurityAndComplianceCenterApiV3) NewListProviderTypesOptions() *ListProviderTypesOptions {
	return &ListProviderTypesOptions{}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListProviderTypesOptions) SetInstanceID(instanceID string) *ListProviderTypesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListProviderTypesOptions) SetXCorrelationID(xCorrelationID string) *ListProviderTypesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListProviderTypesOptions) SetXRequestID(xRequestID string) *ListProviderTypesOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProviderTypesOptions) SetHeaders(param map[string]string) *ListProviderTypesOptions {
	options.Headers = param
	return options
}

// ListReportEvaluationsOptions : The ListReportEvaluations options.
type ListReportEvaluationsOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The ID of the assessment.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The ID of component.
	ComponentID *string `json:"component_id,omitempty"`

	// The ID of the evaluation target.
	TargetID *string `json:"target_id,omitempty"`

	// The name of the evaluation target.
	TargetName *string `json:"target_name,omitempty"`

	// The evaluation status value.
	Status *string `json:"status,omitempty"`

	// The indication of what resource to start the page on.
	Start *string `json:"start,omitempty"`

	// The indication of many resources to return, unless the response is  the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListReportEvaluationsOptions.Status property.
// The evaluation status value.
const (
	ListReportEvaluationsOptions_Status_Error   = "error"
	ListReportEvaluationsOptions_Status_Failure = "failure"
	ListReportEvaluationsOptions_Status_Pass    = "pass"
	ListReportEvaluationsOptions_Status_Skipped = "skipped"
)

// NewListReportEvaluationsOptions : Instantiate ListReportEvaluationsOptions
func (*SecurityAndComplianceCenterApiV3) NewListReportEvaluationsOptions(instanceID string, reportID string) *ListReportEvaluationsOptions {
	return &ListReportEvaluationsOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListReportEvaluationsOptions) SetInstanceID(instanceID string) *ListReportEvaluationsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *ListReportEvaluationsOptions) SetReportID(reportID string) *ListReportEvaluationsOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListReportEvaluationsOptions) SetXCorrelationID(xCorrelationID string) *ListReportEvaluationsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListReportEvaluationsOptions) SetXRequestID(xRequestID string) *ListReportEvaluationsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetAssessmentID : Allow user to set AssessmentID
func (_options *ListReportEvaluationsOptions) SetAssessmentID(assessmentID string) *ListReportEvaluationsOptions {
	_options.AssessmentID = core.StringPtr(assessmentID)
	return _options
}

// SetComponentID : Allow user to set ComponentID
func (_options *ListReportEvaluationsOptions) SetComponentID(componentID string) *ListReportEvaluationsOptions {
	_options.ComponentID = core.StringPtr(componentID)
	return _options
}

// SetTargetID : Allow user to set TargetID
func (_options *ListReportEvaluationsOptions) SetTargetID(targetID string) *ListReportEvaluationsOptions {
	_options.TargetID = core.StringPtr(targetID)
	return _options
}

// SetTargetName : Allow user to set TargetName
func (_options *ListReportEvaluationsOptions) SetTargetName(targetName string) *ListReportEvaluationsOptions {
	_options.TargetName = core.StringPtr(targetName)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *ListReportEvaluationsOptions) SetStatus(status string) *ListReportEvaluationsOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListReportEvaluationsOptions) SetStart(start string) *ListReportEvaluationsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListReportEvaluationsOptions) SetLimit(limit int64) *ListReportEvaluationsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListReportEvaluationsOptions) SetHeaders(param map[string]string) *ListReportEvaluationsOptions {
	options.Headers = param
	return options
}

// ListReportResourcesOptions : The ListReportResources options.
type ListReportResourcesOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The ID of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the resource.
	ResourceName *string `json:"resource_name,omitempty"`

	// The ID of the account owning a resource.
	AccountID *string `json:"account_id,omitempty"`

	// The ID of component.
	ComponentID *string `json:"component_id,omitempty"`

	// The compliance status value.
	Status *string `json:"status,omitempty"`

	// This field sorts resources by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// The indication of what resource to start the page on.
	Start *string `json:"start,omitempty"`

	// The indication of many resources to return, unless the response is  the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListReportResourcesOptions.Status property.
// The compliance status value.
const (
	ListReportResourcesOptions_Status_Compliant              = "compliant"
	ListReportResourcesOptions_Status_NotCompliant           = "not_compliant"
	ListReportResourcesOptions_Status_UnableToPerform        = "unable_to_perform"
	ListReportResourcesOptions_Status_UserEvaluationRequired = "user_evaluation_required"
)

// Constants associated with the ListReportResourcesOptions.Sort property.
// This field sorts resources by using a valid sort field. To learn more, see
// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
const (
	ListReportResourcesOptions_Sort_AccountID    = "account_id"
	ListReportResourcesOptions_Sort_ComponentID  = "component_id"
	ListReportResourcesOptions_Sort_ResourceName = "resource_name"
	ListReportResourcesOptions_Sort_Status       = "status"
)

// NewListReportResourcesOptions : Instantiate ListReportResourcesOptions
func (*SecurityAndComplianceCenterApiV3) NewListReportResourcesOptions(instanceID string, reportID string) *ListReportResourcesOptions {
	return &ListReportResourcesOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListReportResourcesOptions) SetInstanceID(instanceID string) *ListReportResourcesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *ListReportResourcesOptions) SetReportID(reportID string) *ListReportResourcesOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListReportResourcesOptions) SetXCorrelationID(xCorrelationID string) *ListReportResourcesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListReportResourcesOptions) SetXRequestID(xRequestID string) *ListReportResourcesOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ListReportResourcesOptions) SetID(id string) *ListReportResourcesOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetResourceName : Allow user to set ResourceName
func (_options *ListReportResourcesOptions) SetResourceName(resourceName string) *ListReportResourcesOptions {
	_options.ResourceName = core.StringPtr(resourceName)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListReportResourcesOptions) SetAccountID(accountID string) *ListReportResourcesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetComponentID : Allow user to set ComponentID
func (_options *ListReportResourcesOptions) SetComponentID(componentID string) *ListReportResourcesOptions {
	_options.ComponentID = core.StringPtr(componentID)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *ListReportResourcesOptions) SetStatus(status string) *ListReportResourcesOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListReportResourcesOptions) SetSort(sort string) *ListReportResourcesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListReportResourcesOptions) SetStart(start string) *ListReportResourcesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListReportResourcesOptions) SetLimit(limit int64) *ListReportResourcesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListReportResourcesOptions) SetHeaders(param map[string]string) *ListReportResourcesOptions {
	options.Headers = param
	return options
}

// ListReportsOptions : The ListReports options.
type ListReportsOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// The ID of the attachment.
	AttachmentID *string `json:"attachment_id,omitempty"`

	// The report group ID.
	GroupID *string `json:"group_id,omitempty"`

	// The ID of the profile.
	ProfileID *string `json:"profile_id,omitempty"`

	// The type of the scan.
	Type *string `json:"type,omitempty"`

	// The indication of what resource to start the page on.
	Start *string `json:"start,omitempty"`

	// The indication of many resources to return, unless the response is  the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// This field sorts results by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListReportsOptions.Type property.
// The type of the scan.
const (
	ListReportsOptions_Type_Ondemand  = "ondemand"
	ListReportsOptions_Type_Scheduled = "scheduled"
)

// NewListReportsOptions : Instantiate ListReportsOptions
func (*SecurityAndComplianceCenterApiV3) NewListReportsOptions(instanceID string) *ListReportsOptions {
	return &ListReportsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListReportsOptions) SetInstanceID(instanceID string) *ListReportsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListReportsOptions) SetXCorrelationID(xCorrelationID string) *ListReportsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListReportsOptions) SetXRequestID(xRequestID string) *ListReportsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *ListReportsOptions) SetAttachmentID(attachmentID string) *ListReportsOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetGroupID : Allow user to set GroupID
func (_options *ListReportsOptions) SetGroupID(groupID string) *ListReportsOptions {
	_options.GroupID = core.StringPtr(groupID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ListReportsOptions) SetProfileID(profileID string) *ListReportsOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListReportsOptions) SetType(typeVar string) *ListReportsOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListReportsOptions) SetStart(start string) *ListReportsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListReportsOptions) SetLimit(limit int64) *ListReportsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListReportsOptions) SetSort(sort string) *ListReportsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListReportsOptions) SetHeaders(param map[string]string) *ListReportsOptions {
	options.Headers = param
	return options
}

// ListRulesOptions : The ListRules options.
type ListRulesOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-Id,omitempty"`

	// The list of only user-defined, or system-defined rules.
	Type *string `json:"type,omitempty"`

	// The indication of whether to search for rules with a specific string string in the name, description, or labels.
	Search *string `json:"search,omitempty"`

	// Searches for rules targeting corresponding service.
	ServiceName *string `json:"service_name,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListRulesOptions : Instantiate ListRulesOptions
func (*SecurityAndComplianceCenterApiV3) NewListRulesOptions(instanceID string) *ListRulesOptions {
	return &ListRulesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListRulesOptions) SetInstanceID(instanceID string) *ListRulesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListRulesOptions) SetXCorrelationID(xCorrelationID string) *ListRulesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListRulesOptions) SetXRequestID(xRequestID string) *ListRulesOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListRulesOptions) SetType(typeVar string) *ListRulesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListRulesOptions) SetSearch(search string) *ListRulesOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetServiceName : Allow user to set ServiceName
func (_options *ListRulesOptions) SetServiceName(serviceName string) *ListRulesOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListRulesOptions) SetHeaders(param map[string]string) *ListRulesOptions {
	options.Headers = param
	return options
}

// MultiCloudScope : The scope payload for the multi cloud feature.
type MultiCloudScope struct {
	// The environment that relates to this scope.
	Environment *string `json:"environment" validate:"required"`

	// The properties supported for scoping by this environment.
	Properties []PropertyItem `json:"properties" validate:"required"`
}

// NewMultiCloudScope : Instantiate MultiCloudScope (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewMultiCloudScope(environment string, properties []PropertyItem) (_model *MultiCloudScope, err error) {
	_model = &MultiCloudScope{
		Environment: core.StringPtr(environment),
		Properties:  properties,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalMultiCloudScope unmarshals an instance of MultiCloudScope from the specified map of raw messages.
func UnmarshalMultiCloudScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MultiCloudScope)
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalPropertyItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ObjectStorage : The Cloud Object Storage settings.
type ObjectStorage struct {
	// The connected Cloud Object Storage instance CRN.
	InstanceCrn *string `json:"instance_crn,omitempty"`

	// The connected Cloud Object Storage bucket name.
	Bucket *string `json:"bucket,omitempty"`

	// The connected Cloud Object Storage bucket location.
	BucketLocation *string `json:"bucket_location,omitempty"`

	// The connected Cloud Object Storage bucket endpoint.
	BucketEndpoint *string `json:"bucket_endpoint,omitempty"`

	// The date when the bucket connection was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`
}

// UnmarshalObjectStorage unmarshals an instance of ObjectStorage from the specified map of raw messages.
func UnmarshalObjectStorage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ObjectStorage)
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket", &obj.Bucket)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket_location", &obj.BucketLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket_endpoint", &obj.BucketEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageHRef : The page reference.
type PageHRef struct {
	// The URL for the first and next page.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPageHRef unmarshals an instance of PageHRef from the specified map of raw messages.
func UnmarshalPageHRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageHRef)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageHRefFirst : A page reference.
type PageHRefFirst struct {
	// A URL for the first and next page.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPageHRefFirst unmarshals an instance of PageHRefFirst from the specified map of raw messages.
func UnmarshalPageHRefFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageHRefFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageHRefNext : A page reference.
type PageHRefNext struct {
	// A URL for the first and next page.
	Href *string `json:"href" validate:"required"`

	// The token of the next page when present.
	Start *string `json:"start,omitempty"`
}

// UnmarshalPageHRefNext unmarshals an instance of PageHRefNext from the specified map of raw messages.
func UnmarshalPageHRefNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageHRefNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedCollectionFirst : The reference to the first page of entries.
type PaginatedCollectionFirst struct {
	// The reference URL for the first few entries.
	Href *string `json:"href,omitempty"`
}

// UnmarshalPaginatedCollectionFirst unmarshals an instance of PaginatedCollectionFirst from the specified map of raw messages.
func UnmarshalPaginatedCollectionFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedCollectionFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedCollectionNext : The reference URL for the next few entries.
type PaginatedCollectionNext struct {
	// The reference URL for the entries.
	Href *string `json:"href,omitempty"`

	// The reference to the start of the list of entries.
	Start *string `json:"start,omitempty"`
}

// UnmarshalPaginatedCollectionNext unmarshals an instance of PaginatedCollectionNext from the specified map of raw messages.
func UnmarshalPaginatedCollectionNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedCollectionNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Parameter : The rule import parameter.
type Parameter struct {
	// The import parameter name.
	Name *string `json:"name,omitempty"`

	// The display name of the property.
	DisplayName *string `json:"display_name,omitempty"`

	// The propery description.
	Description *string `json:"description,omitempty"`

	// The property type.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the Parameter.Type property.
// The property type.
const (
	Parameter_Type_Boolean    = "boolean"
	Parameter_Type_General    = "general"
	Parameter_Type_IpList     = "ip_list"
	Parameter_Type_Numeric    = "numeric"
	Parameter_Type_String     = "string"
	Parameter_Type_StringList = "string_list"
	Parameter_Type_Timestamp  = "timestamp"
)

// UnmarshalParameter unmarshals an instance of Parameter from the specified map of raw messages.
func UnmarshalParameter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Parameter)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ParameterInfo : The parameter details.
type ParameterInfo struct {
	// The parameter name.
	ParameterName *string `json:"parameter_name,omitempty"`

	// The parameter display name.
	ParameterDisplayName *string `json:"parameter_display_name,omitempty"`

	// The parameter type.
	ParameterType *string `json:"parameter_type,omitempty"`

	// The property value.
	ParameterValue interface{} `json:"parameter_value,omitempty"`
}

// Constants associated with the ParameterInfo.ParameterType property.
// The parameter type.
const (
	ParameterInfo_ParameterType_Boolean    = "boolean"
	ParameterInfo_ParameterType_General    = "general"
	ParameterInfo_ParameterType_IpList     = "ip_list"
	ParameterInfo_ParameterType_Numeric    = "numeric"
	ParameterInfo_ParameterType_String     = "string"
	ParameterInfo_ParameterType_StringList = "string_list"
	ParameterInfo_ParameterType_Timestamp  = "timestamp"
)

// UnmarshalParameterInfo unmarshals an instance of ParameterInfo from the specified map of raw messages.
func UnmarshalParameterInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ParameterInfo)
	err = core.UnmarshalPrimitive(m, "parameter_name", &obj.ParameterName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_display_name", &obj.ParameterDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_type", &obj.ParameterType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_value", &obj.ParameterValue)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PostTestEventOptions : The PostTestEvent options.
type PostTestEventOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request, and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The supplied or generated value of this header is logged for a request, and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostTestEventOptions : Instantiate PostTestEventOptions
func (*SecurityAndComplianceCenterApiV3) NewPostTestEventOptions(instanceID string) *PostTestEventOptions {
	return &PostTestEventOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *PostTestEventOptions) SetInstanceID(instanceID string) *PostTestEventOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *PostTestEventOptions) SetXCorrelationID(xCorrelationID string) *PostTestEventOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *PostTestEventOptions) SetXRequestID(xRequestID string) *PostTestEventOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostTestEventOptions) SetHeaders(param map[string]string) *PostTestEventOptions {
	options.Headers = param
	return options
}

// Profile : The response body of the profile.
type Profile struct {
	// The unique ID of the profile.
	ID *string `json:"id,omitempty"`

	// The profile name.
	ProfileName *string `json:"profile_name,omitempty"`

	// The profile description.
	ProfileDescription *string `json:"profile_description,omitempty"`

	// The profile type, such as custom or predefined.
	ProfileType *string `json:"profile_type,omitempty"`

	// The version status of the profile.
	ProfileVersion *string `json:"profile_version,omitempty"`

	// The version group label of the profile.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// The instance ID.
	InstanceID *string `json:"instance_id,omitempty"`

	// The latest version of the profile.
	Latest *bool `json:"latest,omitempty"`

	// The indication of whether hierarchy is enabled for the profile.
	HierarchyEnabled *bool `json:"hierarchy_enabled,omitempty"`

	// The user who created the profile.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the profile was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The user who updated the profile.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The date when the profile was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The number of controls for the profile.
	ControlsCount *int64 `json:"controls_count,omitempty"`

	// The number of parent controls for the profile.
	ControlParentsCount *int64 `json:"control_parents_count,omitempty"`

	// The number of attachments related to this profile.
	AttachmentsCount *int64 `json:"attachments_count,omitempty"`

	// The array of controls that are used to create the profile.
	Controls []ProfileControls `json:"controls,omitempty"`

	// The default parameters of the profile.
	DefaultParameters []DefaultParametersPrototype `json:"default_parameters,omitempty"`
}

// Constants associated with the Profile.ProfileType property.
// The profile type, such as custom or predefined.
const (
	Profile_ProfileType_Custom     = "custom"
	Profile_ProfileType_Predefined = "predefined"
)

// UnmarshalProfile unmarshals an instance of Profile from the specified map of raw messages.
func UnmarshalProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Profile)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_name", &obj.ProfileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_description", &obj.ProfileDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_type", &obj.ProfileType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_version", &obj.ProfileVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_group_label", &obj.VersionGroupLabel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "latest", &obj.Latest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hierarchy_enabled", &obj.HierarchyEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_count", &obj.ControlsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_parents_count", &obj.ControlParentsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachments_count", &obj.AttachmentsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalProfileControls)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "default_parameters", &obj.DefaultParameters, UnmarshalDefaultParametersPrototype)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileCollection : The response body to get all profiles that are linked to your account.
type ProfileCollection struct {
	// The number of profiles.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The limit of profiles that can be created.
	Limit *int64 `json:"limit" validate:"required"`

	// The reference to the first page of entries.
	First *PaginatedCollectionFirst `json:"first" validate:"required"`

	// The reference URL for the next few entries.
	Next *PaginatedCollectionNext `json:"next" validate:"required"`

	// The profiles.
	Profiles []ProfileItem `json:"profiles" validate:"required"`
}

// UnmarshalProfileCollection unmarshals an instance of ProfileCollection from the specified map of raw messages.
func UnmarshalProfileCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedCollectionFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedCollectionNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalProfileItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ProfileCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ProfileControls : The control details for the profile.
type ProfileControls struct {
	// The ID of the control library that contains the profile.
	ControlLibraryID *string `json:"control_library_id,omitempty"`

	// The unique ID of the control library that contains the profile.
	ControlID *string `json:"control_id,omitempty"`

	// The most recent version of the control library.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The control name.
	ControlName *string `json:"control_name,omitempty"`

	// The control description.
	ControlDescription *string `json:"control_description,omitempty"`

	// The control category.
	ControlCategory *string `json:"control_category,omitempty"`

	// The parent control.
	ControlParent *string `json:"control_parent,omitempty"`

	// Is this a control that can be automated or manually evaluated.
	ControlRequirement *bool `json:"control_requirement,omitempty"`

	// The control documentation.
	ControlDocs *ControlDocs `json:"control_docs,omitempty"`

	// The number of control specifications.
	ControlSpecificationsCount *int64 `json:"control_specifications_count,omitempty"`

	// The control specifications.
	ControlSpecifications []ControlSpecifications `json:"control_specifications,omitempty"`
}

// UnmarshalProfileControls unmarshals an instance of ProfileControls from the specified map of raw messages.
func UnmarshalProfileControls(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileControls)
	err = core.UnmarshalPrimitive(m, "control_library_id", &obj.ControlLibraryID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_id", &obj.ControlID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_version", &obj.ControlLibraryVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_name", &obj.ControlName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_description", &obj.ControlDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_category", &obj.ControlCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_parent", &obj.ControlParent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_requirement", &obj.ControlRequirement)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_docs", &obj.ControlDocs, UnmarshalControlDocs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_specifications_count", &obj.ControlSpecificationsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_specifications", &obj.ControlSpecifications, UnmarshalControlSpecifications)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileControlsPrototype : The control details of a profile.
type ProfileControlsPrototype struct {
	// The ID of the control library that contains the profile.
	ControlLibraryID *string `json:"control_library_id,omitempty"`

	// The control ID.
	ControlID *string `json:"control_id,omitempty"`
}

// UnmarshalProfileControlsPrototype unmarshals an instance of ProfileControlsPrototype from the specified map of raw messages.
func UnmarshalProfileControlsPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileControlsPrototype)
	err = core.UnmarshalPrimitive(m, "control_library_id", &obj.ControlLibraryID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_id", &obj.ControlID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileInfo : The profile information.
type ProfileInfo struct {
	// The profile ID.
	ID *string `json:"id,omitempty"`

	// The profile name.
	Name *string `json:"name,omitempty"`

	// The profile version.
	Version *string `json:"version,omitempty"`
}

// UnmarshalProfileInfo unmarshals an instance of ProfileInfo from the specified map of raw messages.
func UnmarshalProfileInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileItem : ProfileItem struct
type ProfileItem struct {
	// The profile ID.
	ID *string `json:"id,omitempty"`

	// The profile name.
	ProfileName *string `json:"profile_name,omitempty"`

	// The profile description.
	ProfileDescription *string `json:"profile_description,omitempty"`

	// The profile type.
	ProfileType *string `json:"profile_type,omitempty"`

	// The profile version.
	ProfileVersion *string `json:"profile_version,omitempty"`

	// The version group label.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// The latest profile.
	Latest *bool `json:"latest,omitempty"`

	// The user who created the profile.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the profile was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The user who updated the profile.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The date when the profile was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The number of controls.
	ControlsCount *int64 `json:"controls_count,omitempty"`

	// The number of attachments.
	AttachmentsCount *int64 `json:"attachments_count,omitempty"`
}

// UnmarshalProfileItem unmarshals an instance of ProfileItem from the specified map of raw messages.
func UnmarshalProfileItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_name", &obj.ProfileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_description", &obj.ProfileDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_type", &obj.ProfileType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_version", &obj.ProfileVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_group_label", &obj.VersionGroupLabel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "latest", &obj.Latest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_count", &obj.ControlsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachments_count", &obj.AttachmentsCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Property : The property.
type Property struct {
	// The property name.
	Property *string `json:"property,omitempty"`

	// The property description.
	PropertyDescription *string `json:"property_description,omitempty"`

	// The property operator.
	Operator *string `json:"operator,omitempty"`

	// The property value.
	ExpectedValue interface{} `json:"expected_value,omitempty"`

	// The property value.
	FoundValue interface{} `json:"found_value,omitempty"`
}

// UnmarshalProperty unmarshals an instance of Property from the specified map of raw messages.
func UnmarshalProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Property)
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property_description", &obj.PropertyDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_value", &obj.ExpectedValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "found_value", &obj.FoundValue)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PropertyItem : The properties supported for scoping by this environment.
type PropertyItem struct {
	// The name of the property.
	Name *string `json:"name,omitempty"`

	// The value of the property.
	Value interface{} `json:"value,omitempty"`
}

// UnmarshalPropertyItem unmarshals an instance of PropertyItem from the specified map of raw messages.
func UnmarshalPropertyItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PropertyItem)
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

// ProviderTypeInstanceItem : A provider type instance.
type ProviderTypeInstanceItem struct {
	// The unique identifier of the provider type instance.
	ID *string `json:"id,omitempty"`

	// The type of the provider type.
	Type *string `json:"type,omitempty"`

	// The name of the provider type instance.
	Name *string `json:"name,omitempty"`

	// The attributes for connecting to the provider type instance.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// Time at which resource was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Time at which resource was updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalProviderTypeInstanceItem unmarshals an instance of ProviderTypeInstanceItem from the specified map of raw messages.
func UnmarshalProviderTypeInstanceItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderTypeInstanceItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attributes", &obj.Attributes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderTypeInstancesResponse : Provider type instances response.
type ProviderTypeInstancesResponse struct {
	// The array of instances for a provider type.
	ProviderTypeInstances []ProviderTypeInstanceItem `json:"provider_type_instances,omitempty"`
}

// UnmarshalProviderTypeInstancesResponse unmarshals an instance of ProviderTypeInstancesResponse from the specified map of raw messages.
func UnmarshalProviderTypeInstancesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderTypeInstancesResponse)
	err = core.UnmarshalModel(m, "provider_type_instances", &obj.ProviderTypeInstances, UnmarshalProviderTypeInstanceItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderTypeItem : The provider type item.
type ProviderTypeItem struct {
	// The unique identifier of the provider type.
	ID *string `json:"id" validate:"required"`

	// The type of the provider type.
	Type *string `json:"type" validate:"required"`

	// The name of the provider type.
	Name *string `json:"name" validate:"required"`

	// The provider type description.
	Description *string `json:"description" validate:"required"`

	// A boolean that indicates whether the provider type is s2s-enabled.
	S2sEnabled *bool `json:"s2s_enabled" validate:"required"`

	// The maximum number of instances that can be created for the provider type.
	InstanceLimit *int64 `json:"instance_limit" validate:"required"`

	// The mode that is used to get results from provider (`PUSH` or `PULL`).
	Mode *string `json:"mode" validate:"required"`

	// The format of the results that a provider supports.
	DataType *string `json:"data_type" validate:"required"`

	// The icon of a provider in .svg format that is encoded as a base64 string.
	Icon *string `json:"icon" validate:"required"`

	// The label that is associated with the provider type.
	Label *LabelType `json:"label,omitempty"`

	// The attributes that are required when you're creating an instance of a provider type. The attributes field can have
	// multiple  keys in its value. Each of those keys has a value  object that includes the type, and display name as
	// keys. For example, `{type:"", display_name:""}`.
	// **NOTE;** If the provider type is s2s-enabled, which means that if the `s2s_enabled` field is set to `true`, then a
	// CRN field of type text is required in the attributes value object.
	Attributes map[string]AdditionalProperty `json:"attributes" validate:"required"`

	// Time at which resource was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Time at which resource was updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalProviderTypeItem unmarshals an instance of ProviderTypeItem from the specified map of raw messages.
func UnmarshalProviderTypeItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderTypeItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalPrimitive(m, "s2s_enabled", &obj.S2sEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_limit", &obj.InstanceLimit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "mode", &obj.Mode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data_type", &obj.DataType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "icon", &obj.Icon)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "label", &obj.Label, UnmarshalLabelType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalAdditionalProperty)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderTypesCollection : The provider types collection.
type ProviderTypesCollection struct {
	// The array of provder type.
	ProviderTypes []ProviderTypeItem `json:"provider_types,omitempty"`
}

// UnmarshalProviderTypesCollection unmarshals an instance of ProviderTypesCollection from the specified map of raw messages.
func UnmarshalProviderTypesCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderTypesCollection)
	err = core.UnmarshalModel(m, "provider_types", &obj.ProviderTypes, UnmarshalProviderTypeItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderTypesInstancesResponse : Provider types instances response.
type ProviderTypesInstancesResponse struct {
	// The array of instances for all provider types.
	ProviderTypesInstances []ProviderTypeInstanceItem `json:"provider_types_instances,omitempty"`
}

// UnmarshalProviderTypesInstancesResponse unmarshals an instance of ProviderTypesInstancesResponse from the specified map of raw messages.
func UnmarshalProviderTypesInstancesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderTypesInstancesResponse)
	err = core.UnmarshalModel(m, "provider_types_instances", &obj.ProviderTypesInstances, UnmarshalProviderTypeInstanceItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceCustomControlLibraryOptions : The ReplaceCustomControlLibrary options.
type ReplaceCustomControlLibraryOptions struct {
	// The control library ID.
	ControlLibrariesID *string `json:"control_libraries_id" validate:"required,ne="`

	// The control library ID.
	ID *string `json:"id,omitempty"`

	// The account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The control library name.
	ControlLibraryName *string `json:"control_library_name,omitempty"`

	// The control library description.
	ControlLibraryDescription *string `json:"control_library_description,omitempty"`

	// The control library type.
	ControlLibraryType *string `json:"control_library_type,omitempty"`

	// The version group label.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// The control library version.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The date when the control library was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The user who created the control library.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the control library was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The user who updated the control library.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The latest version of the control library.
	Latest *bool `json:"latest,omitempty"`

	// The indication of whether hierarchy is enabled for the control library.
	HierarchyEnabled *bool `json:"hierarchy_enabled,omitempty"`

	// The number of controls.
	ControlsCount *int64 `json:"controls_count,omitempty"`

	// The number of parent controls in the control library.
	ControlParentsCount *int64 `json:"control_parents_count,omitempty"`

	// The list of controls in a control library.
	Controls []ControlsInControlLib `json:"controls,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string

	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`
}

// Constants associated with the ReplaceCustomControlLibraryOptions.ControlLibraryType property.
// The control library type.
const (
	ReplaceCustomControlLibraryOptions_ControlLibraryType_Custom     = "custom"
	ReplaceCustomControlLibraryOptions_ControlLibraryType_Predefined = "predefined"
)

// NewReplaceCustomControlLibraryOptions : Instantiate ReplaceCustomControlLibraryOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceCustomControlLibraryOptions(instanceID string, controlLibrariesID string) *ReplaceCustomControlLibraryOptions {
	return &ReplaceCustomControlLibraryOptions{
		InstanceID:         core.StringPtr(instanceID),
		ControlLibrariesID: core.StringPtr(controlLibrariesID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceCustomControlLibraryOptions) SetInstanceID(instanceID string) *ReplaceCustomControlLibraryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetControlLibrariesID : Allow user to set ControlLibrariesID
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibrariesID(controlLibrariesID string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibrariesID = core.StringPtr(controlLibrariesID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ReplaceCustomControlLibraryOptions) SetID(id string) *ReplaceCustomControlLibraryOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceCustomControlLibraryOptions) SetAccountID(accountID string) *ReplaceCustomControlLibraryOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetControlLibraryName : Allow user to set ControlLibraryName
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryName(controlLibraryName string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryName = core.StringPtr(controlLibraryName)
	return _options
}

// SetControlLibraryDescription : Allow user to set ControlLibraryDescription
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryDescription(controlLibraryDescription string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryDescription = core.StringPtr(controlLibraryDescription)
	return _options
}

// SetControlLibraryType : Allow user to set ControlLibraryType
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryType(controlLibraryType string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryType = core.StringPtr(controlLibraryType)
	return _options
}

// SetVersionGroupLabel : Allow user to set VersionGroupLabel
func (_options *ReplaceCustomControlLibraryOptions) SetVersionGroupLabel(versionGroupLabel string) *ReplaceCustomControlLibraryOptions {
	_options.VersionGroupLabel = core.StringPtr(versionGroupLabel)
	return _options
}

// SetControlLibraryVersion : Allow user to set ControlLibraryVersion
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryVersion(controlLibraryVersion string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryVersion = core.StringPtr(controlLibraryVersion)
	return _options
}

// SetCreatedOn : Allow user to set CreatedOn
func (_options *ReplaceCustomControlLibraryOptions) SetCreatedOn(createdOn *strfmt.DateTime) *ReplaceCustomControlLibraryOptions {
	_options.CreatedOn = createdOn
	return _options
}

// SetCreatedBy : Allow user to set CreatedBy
func (_options *ReplaceCustomControlLibraryOptions) SetCreatedBy(createdBy string) *ReplaceCustomControlLibraryOptions {
	_options.CreatedBy = core.StringPtr(createdBy)
	return _options
}

// SetUpdatedOn : Allow user to set UpdatedOn
func (_options *ReplaceCustomControlLibraryOptions) SetUpdatedOn(updatedOn *strfmt.DateTime) *ReplaceCustomControlLibraryOptions {
	_options.UpdatedOn = updatedOn
	return _options
}

// SetUpdatedBy : Allow user to set UpdatedBy
func (_options *ReplaceCustomControlLibraryOptions) SetUpdatedBy(updatedBy string) *ReplaceCustomControlLibraryOptions {
	_options.UpdatedBy = core.StringPtr(updatedBy)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *ReplaceCustomControlLibraryOptions) SetLatest(latest bool) *ReplaceCustomControlLibraryOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetHierarchyEnabled : Allow user to set HierarchyEnabled
func (_options *ReplaceCustomControlLibraryOptions) SetHierarchyEnabled(hierarchyEnabled bool) *ReplaceCustomControlLibraryOptions {
	_options.HierarchyEnabled = core.BoolPtr(hierarchyEnabled)
	return _options
}

// SetControlsCount : Allow user to set ControlsCount
func (_options *ReplaceCustomControlLibraryOptions) SetControlsCount(controlsCount int64) *ReplaceCustomControlLibraryOptions {
	_options.ControlsCount = core.Int64Ptr(controlsCount)
	return _options
}

// SetControlParentsCount : Allow user to set ControlParentsCount
func (_options *ReplaceCustomControlLibraryOptions) SetControlParentsCount(controlParentsCount int64) *ReplaceCustomControlLibraryOptions {
	_options.ControlParentsCount = core.Int64Ptr(controlParentsCount)
	return _options
}

// SetControls : Allow user to set Controls
func (_options *ReplaceCustomControlLibraryOptions) SetControls(controls []ControlsInControlLib) *ReplaceCustomControlLibraryOptions {
	_options.Controls = controls
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ReplaceCustomControlLibraryOptions) SetXCorrelationID(xCorrelationID string) *ReplaceCustomControlLibraryOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ReplaceCustomControlLibraryOptions) SetXRequestID(xRequestID string) *ReplaceCustomControlLibraryOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceCustomControlLibraryOptions) SetHeaders(param map[string]string) *ReplaceCustomControlLibraryOptions {
	options.Headers = param
	return options
}

// ReplaceProfileAttachmentOptions : The ReplaceProfileAttachment options.
type ReplaceProfileAttachmentOptions struct {
	// The attachment ID.
	AttachmentID *string `json:"attachment_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The ID of the attachment.
	ID *string `json:"id,omitempty"`

	// The account ID that is associated to the attachment.
	AccountID *string `json:"account_id,omitempty"`

	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The scope payload for the multi cloud feature.
	Scope []MultiCloudScope `json:"scope,omitempty"`

	// The date when the attachment was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The user who created the attachment.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the attachment was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The user who updated the attachment.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The status of an attachment evaluation.
	Status *string `json:"status,omitempty"`

	// The schedule of an attachment evaluation.
	Schedule *string `json:"schedule,omitempty"`

	// The request payload of the attachment notifications.
	Notifications *AttachmentsNotificationsPrototype `json:"notifications,omitempty"`

	// The profile parameters for the attachment.
	AttachmentParameters []AttachmentParameterPrototype `json:"attachment_parameters,omitempty"`

	// The details of the last scan of an attachment.
	LastScan *LastScan `json:"last_scan,omitempty"`

	// The start time of the next scan.
	NextScanTime *strfmt.DateTime `json:"next_scan_time,omitempty"`

	// The name of the attachment.
	Name *string `json:"name,omitempty"`

	// The description for the attachment.
	Description *string `json:"description,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceProfileAttachmentOptions.Status property.
// The status of an attachment evaluation.
const (
	ReplaceProfileAttachmentOptions_Status_Disabled = "disabled"
	ReplaceProfileAttachmentOptions_Status_Enabled  = "enabled"
)

// Constants associated with the ReplaceProfileAttachmentOptions.Schedule property.
// The schedule of an attachment evaluation.
const (
	ReplaceProfileAttachmentOptions_Schedule_Daily       = "daily"
	ReplaceProfileAttachmentOptions_Schedule_Every30Days = "every_30_days"
	ReplaceProfileAttachmentOptions_Schedule_Every7Days  = "every_7_days"
)

// NewReplaceProfileAttachmentOptions : Instantiate ReplaceProfileAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceProfileAttachmentOptions(instanceID string, attachmentID string, profileID string) *ReplaceProfileAttachmentOptions {
	return &ReplaceProfileAttachmentOptions{
		InstanceID:   core.StringPtr(instanceID),
		AttachmentID: core.StringPtr(attachmentID),
		ProfileID:    core.StringPtr(profileID),
	}
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *ReplaceProfileAttachmentOptions) SetAttachmentID(attachmentID string) *ReplaceProfileAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ReplaceProfileAttachmentOptions) SetProfileID(profileID string) *ReplaceProfileAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ReplaceProfileAttachmentOptions) SetID(id string) *ReplaceProfileAttachmentOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceProfileAttachmentOptions) SetAccountID(accountID string) *ReplaceProfileAttachmentOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceProfileAttachmentOptions) SetInstanceID(instanceID string) *ReplaceProfileAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScope : Allow user to set Scope
func (_options *ReplaceProfileAttachmentOptions) SetScope(scope []MultiCloudScope) *ReplaceProfileAttachmentOptions {
	_options.Scope = scope
	return _options
}

// SetCreatedOn : Allow user to set CreatedOn
func (_options *ReplaceProfileAttachmentOptions) SetCreatedOn(createdOn *strfmt.DateTime) *ReplaceProfileAttachmentOptions {
	_options.CreatedOn = createdOn
	return _options
}

// SetCreatedBy : Allow user to set CreatedBy
func (_options *ReplaceProfileAttachmentOptions) SetCreatedBy(createdBy string) *ReplaceProfileAttachmentOptions {
	_options.CreatedBy = core.StringPtr(createdBy)
	return _options
}

// SetUpdatedOn : Allow user to set UpdatedOn
func (_options *ReplaceProfileAttachmentOptions) SetUpdatedOn(updatedOn *strfmt.DateTime) *ReplaceProfileAttachmentOptions {
	_options.UpdatedOn = updatedOn
	return _options
}

// SetUpdatedBy : Allow user to set UpdatedBy
func (_options *ReplaceProfileAttachmentOptions) SetUpdatedBy(updatedBy string) *ReplaceProfileAttachmentOptions {
	_options.UpdatedBy = core.StringPtr(updatedBy)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *ReplaceProfileAttachmentOptions) SetStatus(status string) *ReplaceProfileAttachmentOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetSchedule : Allow user to set Schedule
func (_options *ReplaceProfileAttachmentOptions) SetSchedule(schedule string) *ReplaceProfileAttachmentOptions {
	_options.Schedule = core.StringPtr(schedule)
	return _options
}

// SetNotifications : Allow user to set Notifications
func (_options *ReplaceProfileAttachmentOptions) SetNotifications(notifications *AttachmentsNotificationsPrototype) *ReplaceProfileAttachmentOptions {
	_options.Notifications = notifications
	return _options
}

// SetAttachmentParameters : Allow user to set AttachmentParameters
func (_options *ReplaceProfileAttachmentOptions) SetAttachmentParameters(attachmentParameters []AttachmentParameterPrototype) *ReplaceProfileAttachmentOptions {
	_options.AttachmentParameters = attachmentParameters
	return _options
}

// SetLastScan : Allow user to set LastScan
func (_options *ReplaceProfileAttachmentOptions) SetLastScan(lastScan *LastScan) *ReplaceProfileAttachmentOptions {
	_options.LastScan = lastScan
	return _options
}

// SetNextScanTime : Allow user to set NextScanTime
func (_options *ReplaceProfileAttachmentOptions) SetNextScanTime(nextScanTime *strfmt.DateTime) *ReplaceProfileAttachmentOptions {
	_options.NextScanTime = nextScanTime
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceProfileAttachmentOptions) SetName(name string) *ReplaceProfileAttachmentOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceProfileAttachmentOptions) SetDescription(description string) *ReplaceProfileAttachmentOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ReplaceProfileAttachmentOptions) SetXCorrelationID(xCorrelationID string) *ReplaceProfileAttachmentOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ReplaceProfileAttachmentOptions) SetXRequestID(xRequestID string) *ReplaceProfileAttachmentOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceProfileAttachmentOptions) SetHeaders(param map[string]string) *ReplaceProfileAttachmentOptions {
	options.Headers = param
	return options
}

// ReplaceProfileOptions : The ReplaceProfile options.
type ReplaceProfileOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The name of the profile.
	ProfileName *string `json:"profile_name" validate:"required"`

	// The description of the profile.
	ProfileDescription *string `json:"profile_description" validate:"required"`

	// The profile type.
	ProfileType *string `json:"profile_type" validate:"required"`

	// The version status of the profile.
	ProfileVersion *string `json:"profile_version,omitempty"`

	// The controls that are in the profile.
	Controls []ProfileControlsPrototype `json:"controls" validate:"required"`

	// The default parameters of the profile.
	DefaultParameters []DefaultParametersPrototype `json:"default_parameters" validate:"required"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is not used for downstream requests and retries of those requests. If a value
	// of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceProfileOptions.ProfileType property.
// The profile type.
const (
	ReplaceProfileOptions_ProfileType_Custom     = "custom"
	ReplaceProfileOptions_ProfileType_Predefined = "predefined"
)

// NewReplaceProfileOptions : Instantiate ReplaceProfileOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceProfileOptions(instanceID string, profileID string, profileName string, profileDescription string, profileType string, controls []ProfileControlsPrototype, defaultParameters []DefaultParametersPrototype) *ReplaceProfileOptions {
	return &ReplaceProfileOptions{
		InstanceID:         core.StringPtr(instanceID),
		ProfileID:          core.StringPtr(profileID),
		ProfileName:        core.StringPtr(profileName),
		ProfileDescription: core.StringPtr(profileDescription),
		ProfileType:        core.StringPtr(profileType),
		Controls:           controls,
		DefaultParameters:  defaultParameters,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceProfileOptions) SetInstanceID(instanceID string) *ReplaceProfileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ReplaceProfileOptions) SetProfileID(profileID string) *ReplaceProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetProfileName : Allow user to set ProfileName
func (_options *ReplaceProfileOptions) SetProfileName(profileName string) *ReplaceProfileOptions {
	_options.ProfileName = core.StringPtr(profileName)
	return _options
}

// SetProfileDescription : Allow user to set ProfileDescription
func (_options *ReplaceProfileOptions) SetProfileDescription(profileDescription string) *ReplaceProfileOptions {
	_options.ProfileDescription = core.StringPtr(profileDescription)
	return _options
}

// SetProfileType : Allow user to set ProfileType
func (_options *ReplaceProfileOptions) SetProfileType(profileType string) *ReplaceProfileOptions {
	_options.ProfileType = core.StringPtr(profileType)
	return _options
}

// SetProfileVersion : Allow user to set ProfileType
func (_options *ReplaceProfileOptions) SetProfileVersion(profileVersion string) *ReplaceProfileOptions {
	_options.ProfileVersion = core.StringPtr(profileVersion)
	return _options
}

// SetControls : Allow user to set Controls
func (_options *ReplaceProfileOptions) SetControls(controls []ProfileControlsPrototype) *ReplaceProfileOptions {
	_options.Controls = controls
	return _options
}

// SetDefaultParameters : Allow user to set DefaultParameters
func (_options *ReplaceProfileOptions) SetDefaultParameters(defaultParameters []DefaultParametersPrototype) *ReplaceProfileOptions {
	_options.DefaultParameters = defaultParameters
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ReplaceProfileOptions) SetXCorrelationID(xCorrelationID string) *ReplaceProfileOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ReplaceProfileOptions) SetXRequestID(xRequestID string) *ReplaceProfileOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceProfileOptions) SetHeaders(param map[string]string) *ReplaceProfileOptions {
	options.Headers = param
	return options
}

// ReplaceRuleOptions : The ReplaceRule options.
type ReplaceRuleOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the corresponding rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// This field compares a supplied `Etag` value with the version that is stored for the requested resource. If the
	// values match, the server allows the request method to continue.
	//
	// To find the `Etag` value, run a GET request on the resource that you want to modify, and check the response headers.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The rule description.
	Description *string `json:"description" validate:"required"`

	// The rule target.
	Target *Target `json:"target" validate:"required"`

	// The required configurations.
	RequiredConfig RequiredConfigIntf `json:"required_config" validate:"required"`

	// The rule type (user_defined or system_defined).
	Type *string `json:"type,omitempty"`

	// The rule version number.
	Version *string `json:"version,omitempty"`

	// The collection of import parameters.
	Import *Import `json:"import,omitempty"`

	// The list of labels that correspond to a rule.
	Labels []string `json:"labels,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceRuleOptions.Type property.
// The rule type (user_defined or system_defined).
const (
	ReplaceRuleOptions_Type_SystemDefined = "system_defined"
	ReplaceRuleOptions_Type_UserDefined   = "user_defined"
)

// NewReplaceRuleOptions : Instantiate ReplaceRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceRuleOptions(instanceID string, ruleID string, ifMatch string, description string, target *Target, requiredConfig RequiredConfigIntf) *ReplaceRuleOptions {
	return &ReplaceRuleOptions{
		InstanceID:     core.StringPtr(instanceID),
		RuleID:         core.StringPtr(ruleID),
		IfMatch:        core.StringPtr(ifMatch),
		Description:    core.StringPtr(description),
		Target:         target,
		RequiredConfig: requiredConfig,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceRuleOptions) SetInstanceID(instanceID string) *ReplaceRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *ReplaceRuleOptions) SetRuleID(ruleID string) *ReplaceRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplaceRuleOptions) SetIfMatch(ifMatch string) *ReplaceRuleOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceRuleOptions) SetDescription(description string) *ReplaceRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *ReplaceRuleOptions) SetTarget(target *Target) *ReplaceRuleOptions {
	_options.Target = target
	return _options
}

// SetRequiredConfig : Allow user to set RequiredConfig
func (_options *ReplaceRuleOptions) SetRequiredConfig(requiredConfig RequiredConfigIntf) *ReplaceRuleOptions {
	_options.RequiredConfig = requiredConfig
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplaceRuleOptions) SetType(typeVar string) *ReplaceRuleOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *ReplaceRuleOptions) SetVersion(version string) *ReplaceRuleOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetImport : Allow user to set Import
func (_options *ReplaceRuleOptions) SetImport(importVar *Import) *ReplaceRuleOptions {
	_options.Import = importVar
	return _options
}

// SetLabels : Allow user to set Labels
func (_options *ReplaceRuleOptions) SetLabels(labels []string) *ReplaceRuleOptions {
	_options.Labels = labels
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ReplaceRuleOptions) SetXCorrelationID(xCorrelationID string) *ReplaceRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ReplaceRuleOptions) SetXRequestID(xRequestID string) *ReplaceRuleOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceRuleOptions) SetHeaders(param map[string]string) *ReplaceRuleOptions {
	options.Headers = param
	return options
}

// Report : The report.
type Report struct {
	// The ID of the report.
	ID *string `json:"id,omitempty"`

	// The group ID that is associated with the report. The group ID combines profile, scope, and attachment IDs.
	GroupID *string `json:"group_id,omitempty"`

	// The date when the report was created.
	CreatedOn *string `json:"created_on,omitempty"`

	// The date when the scan was run.
	ScanTime *string `json:"scan_time,omitempty"`

	// The type of the scan.
	Type *string `json:"type,omitempty"`

	// The Cloud Object Storage object that is associated with the report.
	CosObject *string `json:"cos_object,omitempty"`

	// Instance ID.
	InstanceID *string `json:"instance_id,omitempty"`

	// The account that is associated with a report.
	Account *Account `json:"account,omitempty"`

	// The profile information.
	Profile *ProfileInfo `json:"profile,omitempty"`

	// The attachment that is associated with a report.
	Attachment *Attachment `json:"attachment,omitempty"`
}

// UnmarshalReport unmarshals an instance of Report from the specified map of raw messages.
func UnmarshalReport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Report)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "group_id", &obj.GroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scan_time", &obj.ScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_object", &obj.CosObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account", &obj.Account, UnmarshalAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profile", &obj.Profile, UnmarshalProfileInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachment", &obj.Attachment, UnmarshalAttachment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportControls : The list of controls.
type ReportControls struct {
	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The list of controls that are in the report.
	Controls []ControlWithStats `json:"controls,omitempty"`
}

// Constants associated with the ReportControls.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ReportControls_Status_Compliant              = "compliant"
	ReportControls_Status_NotCompliant           = "not_compliant"
	ReportControls_Status_UnableToPerform        = "unable_to_perform"
	ReportControls_Status_UserEvaluationRequired = "user_evaluation_required"
)

// UnmarshalReportControls unmarshals an instance of ReportControls from the specified map of raw messages.
func UnmarshalReportControls(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportControls)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalControlWithStats)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportLatest : The response body of the `get_latest_reports` operation.
type ReportLatest struct {
	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The compliance stats.
	ControlsSummary *ComplianceStats `json:"controls_summary,omitempty"`

	// The evaluation stats.
	EvaluationsSummary *EvalStats `json:"evaluations_summary,omitempty"`

	// The compliance score.
	Score *ComplianceScore `json:"score,omitempty"`

	// The list of reports.
	Reports []Report `json:"reports,omitempty"`
}

// UnmarshalReportLatest unmarshals an instance of ReportLatest from the specified map of raw messages.
func UnmarshalReportLatest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportLatest)
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls_summary", &obj.ControlsSummary, UnmarshalComplianceStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "evaluations_summary", &obj.EvaluationsSummary, UnmarshalEvalStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "score", &obj.Score, UnmarshalComplianceScore)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "reports", &obj.Reports, UnmarshalReport)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportPage : The page of reports.
type ReportPage struct {
	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The requested page limi.t.
	Limit *int64 `json:"limit" validate:"required"`

	// The token of the next page, when it's present.
	Start *string `json:"start,omitempty"`

	// The page reference.
	First *PageHRef `json:"first" validate:"required"`

	// The page reference.
	Next *PageHRef `json:"next,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The list of reports that are on the page.
	Reports []Report `json:"reports,omitempty"`
}

// UnmarshalReportPage unmarshals an instance of ReportPage from the specified map of raw messages.
func UnmarshalReportPage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportPage)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRef)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRef)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "reports", &obj.Reports, UnmarshalReport)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ReportPage) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	start, err := core.GetQueryParam(resp.Next.Href, "start")
	if err != nil || start == nil {
		return nil, err
	}
	return start, nil
}

// ReportSummary : The report summary.
type ReportSummary struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// Instance ID.
	IsntanceID *string `json:"isntance_id,omitempty"`

	// The account that is associated with a report.
	Account *Account `json:"account,omitempty"`

	// The compliance score.
	Score *ComplianceScore `json:"score,omitempty"`

	// The compliance stats.
	Controls *ComplianceStats `json:"controls,omitempty"`

	// The evaluation stats.
	Evaluations *EvalStats `json:"evaluations,omitempty"`

	// The resource summary.
	Resources *ResourceSummary `json:"resources,omitempty"`
}

// UnmarshalReportSummary unmarshals an instance of ReportSummary from the specified map of raw messages.
func UnmarshalReportSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportSummary)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "isntance_id", &obj.IsntanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account", &obj.Account, UnmarshalAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "score", &obj.Score, UnmarshalComplianceScore)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalComplianceStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "evaluations", &obj.Evaluations, UnmarshalEvalStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResourceSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportTags : The response body of the `get_tags` operation.
type ReportTags struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The collection of different types of tags.
	Tags *Tags `json:"tags,omitempty"`
}

// UnmarshalReportTags unmarshals an instance of ReportTags from the specified map of raw messages.
func UnmarshalReportTags(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportTags)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalTags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportViolationDataPoint : The report violation data point.
type ReportViolationDataPoint struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The group ID that is associated with the report. The group ID combines profile, scope, and attachment IDs.
	ReportGroupID *string `json:"report_group_id,omitempty"`

	// The date when the scan was run.
	ScanTime *string `json:"scan_time,omitempty"`

	// The compliance stats.
	Controls *ComplianceStats `json:"controls,omitempty"`
}

// UnmarshalReportViolationDataPoint unmarshals an instance of ReportViolationDataPoint from the specified map of raw messages.
func UnmarshalReportViolationDataPoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportViolationDataPoint)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_group_id", &obj.ReportGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scan_time", &obj.ScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalComplianceStats)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportViolationsDrift : The response body of the `get_report_violations_drift` operation.
type ReportViolationsDrift struct {
	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The list of report violations data points.
	DataPoints []ReportViolationDataPoint `json:"data_points,omitempty"`
}

// UnmarshalReportViolationsDrift unmarshals an instance of ReportViolationsDrift from the specified map of raw messages.
func UnmarshalReportViolationsDrift(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportViolationsDrift)
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "data_points", &obj.DataPoints, UnmarshalReportViolationDataPoint)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfig : The required configurations.
// Models which "extend" this model:
// - RequiredConfigRequiredConfigAnd
// - RequiredConfigRequiredConfigOr
// - RequiredConfigRequiredConfigBase
type RequiredConfig struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The `AND` required configurations.
	And []RequiredConfigItemsIntf `json:"and,omitempty"`

	// The `OR` required configurations.
	Or []RequiredConfigItemsIntf `json:"or,omitempty"`

	// The property.
	Property *string `json:"property,omitempty"`

	// The operator.
	Operator *string `json:"operator,omitempty"`

	// Schema for any JSON type.
	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the RequiredConfig.Operator property.
// The operator.
const (
	RequiredConfig_Operator_DaysLessThan         = "days_less_than"
	RequiredConfig_Operator_IpsEquals            = "ips_equals"
	RequiredConfig_Operator_IpsInRange           = "ips_in_range"
	RequiredConfig_Operator_IpsNotEquals         = "ips_not_equals"
	RequiredConfig_Operator_IsEmpty              = "is_empty"
	RequiredConfig_Operator_IsFalse              = "is_false"
	RequiredConfig_Operator_IsNotEmpty           = "is_not_empty"
	RequiredConfig_Operator_IsTrue               = "is_true"
	RequiredConfig_Operator_NumEquals            = "num_equals"
	RequiredConfig_Operator_NumGreaterThan       = "num_greater_than"
	RequiredConfig_Operator_NumGreaterThanEquals = "num_greater_than_equals"
	RequiredConfig_Operator_NumLessThan          = "num_less_than"
	RequiredConfig_Operator_NumLessThanEquals    = "num_less_than_equals"
	RequiredConfig_Operator_NumNotEquals         = "num_not_equals"
	RequiredConfig_Operator_StringContains       = "string_contains"
	RequiredConfig_Operator_StringEquals         = "string_equals"
	RequiredConfig_Operator_StringMatch          = "string_match"
	RequiredConfig_Operator_StringNotContains    = "string_not_contains"
	RequiredConfig_Operator_StringNotEquals      = "string_not_equals"
	RequiredConfig_Operator_StringNotMatch       = "string_not_match"
	RequiredConfig_Operator_StringsAllowed       = "strings_allowed"
	RequiredConfig_Operator_StringsInList        = "strings_in_list"
	RequiredConfig_Operator_StringsRequired      = "strings_required"
)

func (*RequiredConfig) isaRequiredConfig() bool {
	return true
}

type RequiredConfigIntf interface {
	isaRequiredConfig() bool
}

// UnmarshalRequiredConfig unmarshals an instance of RequiredConfig from the specified map of raw messages.
func UnmarshalRequiredConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfig)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRequiredConfigItems)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRequiredConfigItems)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// RequiredConfigItems : RequiredConfigItems struct
// Models which "extend" this model:
// - RequiredConfigItemsRequiredConfigOr
// - RequiredConfigItemsRequiredConfigAnd
// - RequiredConfigItemsRequiredConfigBase
type RequiredConfigItems struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The `OR` required configurations.
	Or []RequiredConfigItemsIntf `json:"or,omitempty"`

	// The `AND` required configurations.
	And []RequiredConfigItemsIntf `json:"and,omitempty"`

	// The property.
	Property *string `json:"property,omitempty"`

	// The operator.
	Operator *string `json:"operator,omitempty"`

	// Schema for any JSON type.
	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the RequiredConfigItems.Operator property.
// The operator.
const (
	RequiredConfigItems_Operator_DaysLessThan         = "days_less_than"
	RequiredConfigItems_Operator_IpsEquals            = "ips_equals"
	RequiredConfigItems_Operator_IpsInRange           = "ips_in_range"
	RequiredConfigItems_Operator_IpsNotEquals         = "ips_not_equals"
	RequiredConfigItems_Operator_IsEmpty              = "is_empty"
	RequiredConfigItems_Operator_IsFalse              = "is_false"
	RequiredConfigItems_Operator_IsNotEmpty           = "is_not_empty"
	RequiredConfigItems_Operator_IsTrue               = "is_true"
	RequiredConfigItems_Operator_NumEquals            = "num_equals"
	RequiredConfigItems_Operator_NumGreaterThan       = "num_greater_than"
	RequiredConfigItems_Operator_NumGreaterThanEquals = "num_greater_than_equals"
	RequiredConfigItems_Operator_NumLessThan          = "num_less_than"
	RequiredConfigItems_Operator_NumLessThanEquals    = "num_less_than_equals"
	RequiredConfigItems_Operator_NumNotEquals         = "num_not_equals"
	RequiredConfigItems_Operator_StringContains       = "string_contains"
	RequiredConfigItems_Operator_StringEquals         = "string_equals"
	RequiredConfigItems_Operator_StringMatch          = "string_match"
	RequiredConfigItems_Operator_StringNotContains    = "string_not_contains"
	RequiredConfigItems_Operator_StringNotEquals      = "string_not_equals"
	RequiredConfigItems_Operator_StringNotMatch       = "string_not_match"
	RequiredConfigItems_Operator_StringsAllowed       = "strings_allowed"
	RequiredConfigItems_Operator_StringsInList        = "strings_in_list"
	RequiredConfigItems_Operator_StringsRequired      = "strings_required"
)

func (*RequiredConfigItems) isaRequiredConfigItems() bool {
	return true
}

type RequiredConfigItemsIntf interface {
	isaRequiredConfigItems() bool
}

// UnmarshalRequiredConfigItems unmarshals an instance of RequiredConfigItems from the specified map of raw messages.
func UnmarshalRequiredConfigItems(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigItems)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRequiredConfigItems)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRequiredConfigItems)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// Resource : The resource.
type Resource struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The resource CRN.
	ID *string `json:"id,omitempty"`

	// The resource name.
	ResourceName *string `json:"resource_name,omitempty"`

	// The ID of the component.
	ComponentID *string `json:"component_id,omitempty"`

	// The environment.
	Environment *string `json:"environment,omitempty"`

	// The account that is associated with a report.
	Account *Account `json:"account,omitempty"`

	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of passed evaluations.
	PassCount *int64 `json:"pass_count,omitempty"`

	// The number of failed evaluations.
	FailureCount *int64 `json:"failure_count,omitempty"`

	// The number of evaluations that started, but did not finish, and ended with errors.
	ErrorCount *int64 `json:"error_count,omitempty"`

	// The total number of completed evaluations.
	CompletedCount *int64 `json:"completed_count,omitempty"`
}

// Constants associated with the Resource.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	Resource_Status_Compliant              = "compliant"
	Resource_Status_NotCompliant           = "not_compliant"
	Resource_Status_UnableToPerform        = "unable_to_perform"
	Resource_Status_UserEvaluationRequired = "user_evaluation_required"
)

// UnmarshalResource unmarshals an instance of Resource from the specified map of raw messages.
func UnmarshalResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Resource)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_id", &obj.ComponentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account", &obj.Account, UnmarshalAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pass_count", &obj.PassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_count", &obj.FailureCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "completed_count", &obj.CompletedCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourcePage : The page of resource evaluation summaries.
type ResourcePage struct {
	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The requested page limi.t.
	Limit *int64 `json:"limit" validate:"required"`

	// The token of the next page, when it's present.
	Start *string `json:"start,omitempty"`

	// The page reference.
	First *PageHRef `json:"first" validate:"required"`

	// The page reference.
	Next *PageHRef `json:"next,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The list of resource evaluation summaries that are on the page.
	Resources []Resource `json:"resources,omitempty"`
}

// UnmarshalResourcePage unmarshals an instance of ResourcePage from the specified map of raw messages.
func UnmarshalResourcePage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourcePage)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRef)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRef)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ResourcePage) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	start, err := core.GetQueryParam(resp.Next.Href, "start")
	if err != nil || start == nil {
		return nil, err
	}
	return start, nil
}

// ResourceSummary : The resource summary.
type ResourceSummary struct {
	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`

	// The top 10 resources that have the most failures.
	TopFailed []ResourceSummaryItem `json:"top_failed,omitempty"`
}

// Constants associated with the ResourceSummary.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ResourceSummary_Status_Compliant              = "compliant"
	ResourceSummary_Status_NotCompliant           = "not_compliant"
	ResourceSummary_Status_UnableToPerform        = "unable_to_perform"
	ResourceSummary_Status_UserEvaluationRequired = "user_evaluation_required"
)

// UnmarshalResourceSummary unmarshals an instance of ResourceSummary from the specified map of raw messages.
func UnmarshalResourceSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceSummary)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "top_failed", &obj.TopFailed, UnmarshalResourceSummaryItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceSummaryItem : The resource summary item.
type ResourceSummaryItem struct {
	// The resource name.
	Name *string `json:"name,omitempty"`

	// The resource ID.
	ID *string `json:"id,omitempty"`

	// The service that is managing the resource.
	Service *string `json:"service,omitempty"`

	// The collection of different types of tags.
	Tags *Tags `json:"tags,omitempty"`

	// The account that owns the resource.
	Account *string `json:"account,omitempty"`

	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of passed evaluations.
	PassCount *int64 `json:"pass_count,omitempty"`

	// The number of failed evaluations.
	FailureCount *int64 `json:"failure_count,omitempty"`

	// The number of evaluations that started, but did not finish, and ended with errors.
	ErrorCount *int64 `json:"error_count,omitempty"`

	// The total number of completed evaluations.
	CompletedCount *int64 `json:"completed_count,omitempty"`
}

// Constants associated with the ResourceSummaryItem.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ResourceSummaryItem_Status_Compliant              = "compliant"
	ResourceSummaryItem_Status_NotCompliant           = "not_compliant"
	ResourceSummaryItem_Status_UnableToPerform        = "unable_to_perform"
	ResourceSummaryItem_Status_UserEvaluationRequired = "user_evaluation_required"
)

// UnmarshalResourceSummaryItem unmarshals an instance of ResourceSummaryItem from the specified map of raw messages.
func UnmarshalResourceSummaryItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceSummaryItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service", &obj.Service)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalTags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pass_count", &obj.PassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_count", &obj.FailureCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "completed_count", &obj.CompletedCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Rule : The rule response that corresponds to an account instance.
type Rule struct {
	// The date when the rule was created.
	CreatedOn *strfmt.DateTime `json:"created_on" validate:"required"`

	// The user who created the rule.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the rule was modified.
	UpdatedOn *strfmt.DateTime `json:"updated_on" validate:"required"`

	// The user who modified the rule.
	UpdatedBy *string `json:"updated_by" validate:"required"`

	// The rule ID.
	ID *string `json:"id" validate:"required"`

	// The account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The details of a rule's response.
	Description *string `json:"description" validate:"required"`

	// The rule type (allowable values are `user_defined` or `system_defined`).
	Type *string `json:"type" validate:"required"`

	// The version number of a rule.
	Version *string `json:"version" validate:"required"`

	// The collection of import parameters.
	Import *Import `json:"import,omitempty"`

	// The rule target.
	Target *Target `json:"target" validate:"required"`

	// The required configurations.
	RequiredConfig RequiredConfigIntf `json:"required_config" validate:"required"`

	// The list of labels.
	Labels []string `json:"labels" validate:"required"`
}

// Constants associated with the Rule.Type property.
// The rule type (allowable values are `user_defined` or `system_defined`).
const (
	Rule_Type_SystemDefined = "system_defined"
	Rule_Type_UserDefined   = "user_defined"
)

// UnmarshalRule unmarshals an instance of Rule from the specified map of raw messages.
func UnmarshalRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rule)
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
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
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "import", &obj.Import, UnmarshalImport)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalTarget)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "required_config", &obj.RequiredConfig, UnmarshalRequiredConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleInfo : The rule.
type RuleInfo struct {
	// The rule ID.
	ID *string `json:"id,omitempty"`

	// The rule type.
	Type *string `json:"type,omitempty"`

	// The rule description.
	Description *string `json:"description,omitempty"`

	// The rule version.
	Version *string `json:"version,omitempty"`

	// The rule account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The date when the rule was created.
	CreatedOn *string `json:"created_on,omitempty"`

	// The ID of the user who created the rule.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the rule was updated.
	UpdatedOn *string `json:"updated_on,omitempty"`

	// The ID of the user who updated the rule.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The rule labels.
	Labels []string `json:"labels,omitempty"`
}

// UnmarshalRuleInfo unmarshals an instance of RuleInfo from the specified map of raw messages.
func UnmarshalRuleInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RulesPageBase : Page common fields.
type RulesPageBase struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first" validate:"required"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The collection of rules that correspond to an account instance. Maximum of 100/500 custom rules per
	// stand-alone/enterprise account.
	Rules []Rule `json:"rules,omitempty"`
}

// UnmarshalRulesPageBase unmarshals an instance of RulesPageBase from the specified map of raw messages.
func UnmarshalRulesPageBase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RulesPageBase)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Scan : The response schema for creating a scan.
type Scan struct {
	// The scan ID.
	ID *string `json:"id,omitempty"`

	// The account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The attachment ID of a profile.
	AttachmentID *string `json:"attachment_id,omitempty"`

	// The report ID.
	ReportID *string `json:"report_id,omitempty"`

	// The status of the scan.
	Status *string `json:"status,omitempty"`

	// The last scan time.
	LastScanTime *string `json:"last_scan_time,omitempty"`

	// The next scan time.
	NextScanTime *string `json:"next_scan_time,omitempty"`

	// The type of scan.
	ScanType *string `json:"scan_type,omitempty"`

	// The occurrence of the scan.
	Occurence *int64 `json:"occurence,omitempty"`
}

// Constants associated with the Scan.Status property.
// The status of the scan.
const (
	Scan_Status_Completed  = "completed"
	Scan_Status_InProgress = "in_progress"
)

// Constants associated with the Scan.ScanType property.
// The type of scan.
const (
	Scan_ScanType_Ondemand  = "ondemand"
	Scan_ScanType_Scheduled = "scheduled"
)

// UnmarshalScan unmarshals an instance of Scan from the specified map of raw messages.
func UnmarshalScan(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Scan)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachment_id", &obj.AttachmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_scan_time", &obj.LastScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_scan_time", &obj.NextScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scan_type", &obj.ScanType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "occurence", &obj.Occurence)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeProperty : The properties that are supported for scoping by this attachment.
type ScopeProperty struct {
	// The property name.
	Name *string `json:"name,omitempty"`

	// The property value.
	Value *string `json:"value,omitempty"`
}

// UnmarshalScopeProperty unmarshals an instance of ScopeProperty from the specified map of raw messages.
func UnmarshalScopeProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeProperty)
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

// Settings : The settings.
type Settings struct {
	// The Event Notifications settings.
	EventNotifications *EventNotifications `json:"event_notifications,omitempty"`

	// The Cloud Object Storage settings.
	ObjectStorage *ObjectStorage `json:"object_storage,omitempty"`
}

// UnmarshalSettings unmarshals an instance of Settings from the specified map of raw messages.
func UnmarshalSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Settings)
	err = core.UnmarshalModel(m, "event_notifications", &obj.EventNotifications, UnmarshalEventNotifications)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "object_storage", &obj.ObjectStorage, UnmarshalObjectStorage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tags : The collection of different types of tags.
type Tags struct {
	// The collection of user tags.
	User []string `json:"user,omitempty"`

	// The collection of access tags.
	Access []string `json:"access,omitempty"`

	// The collection of service tags.
	Service []string `json:"service,omitempty"`
}

// UnmarshalTags unmarshals an instance of Tags from the specified map of raw messages.
func UnmarshalTags(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tags)
	err = core.UnmarshalPrimitive(m, "user", &obj.User)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "access", &obj.Access)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service", &obj.Service)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Target : The rule target.
type Target struct {
	// The target service name.
	ServiceName *string `json:"service_name" validate:"required"`

	// The display name of the target service.
	ServiceDisplayName *string `json:"service_display_name,omitempty"`

	// The target resource kind.
	ResourceKind *string `json:"resource_kind" validate:"required"`

	// The list of targets supported properties.
	AdditionalTargetAttributes []AdditionalTargetAttribute `json:"additional_target_attributes,omitempty"`
}

// NewTarget : Instantiate Target (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewTarget(serviceName string, resourceKind string) (_model *Target, err error) {
	_model = &Target{
		ServiceName:  core.StringPtr(serviceName),
		ResourceKind: core.StringPtr(resourceKind),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTarget unmarshals an instance of Target from the specified map of raw messages.
func UnmarshalTarget(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Target)
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_display_name", &obj.ServiceDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_kind", &obj.ResourceKind)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "additional_target_attributes", &obj.AdditionalTargetAttributes, UnmarshalAdditionalTargetAttribute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TargetInfo : The evaluation target.
type TargetInfo struct {
	// The target ID.
	ID *string `json:"id,omitempty"`

	// The target account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The target resource CRN.
	ResourceCrn *string `json:"resource_crn,omitempty"`

	// The target resource name.
	ResourceName *string `json:"resource_name,omitempty"`

	// The target service name.
	ServiceName *string `json:"service_name,omitempty"`
}

// UnmarshalTargetInfo unmarshals an instance of TargetInfo from the specified map of raw messages.
func UnmarshalTargetInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crn", &obj.ResourceCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TestEvent : The details of a test event response.
type TestEvent struct {
	// The indication of whether the event was received by Event Notifications.
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

// UpdateProviderTypeInstanceOptions : The UpdateProviderTypeInstance options.
type UpdateProviderTypeInstanceOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The provider type instance ID.
	ProviderTypeInstanceID *string `json:"provider_type_instance_id" validate:"required,ne="`

	// The provider type instance name.
	Name *string `json:"name,omitempty"`

	// The attributes for connecting to the provider type instance.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateProviderTypeInstanceOptions : Instantiate UpdateProviderTypeInstanceOptions
func (*SecurityAndComplianceCenterApiV3) NewUpdateProviderTypeInstanceOptions(instanceID string, providerTypeID string, providerTypeInstanceID string) *UpdateProviderTypeInstanceOptions {
	return &UpdateProviderTypeInstanceOptions{
		InstanceID:             core.StringPtr(instanceID),
		ProviderTypeID:         core.StringPtr(providerTypeID),
		ProviderTypeInstanceID: core.StringPtr(providerTypeInstanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateProviderTypeInstanceOptions) SetInstanceID(instanceID string) *UpdateProviderTypeInstanceOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *UpdateProviderTypeInstanceOptions) SetProviderTypeID(providerTypeID string) *UpdateProviderTypeInstanceOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetProviderTypeInstanceID : Allow user to set ProviderTypeInstanceID
func (_options *UpdateProviderTypeInstanceOptions) SetProviderTypeInstanceID(providerTypeInstanceID string) *UpdateProviderTypeInstanceOptions {
	_options.ProviderTypeInstanceID = core.StringPtr(providerTypeInstanceID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateProviderTypeInstanceOptions) SetName(name string) *UpdateProviderTypeInstanceOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAttributes : Allow user to set Attributes
func (_options *UpdateProviderTypeInstanceOptions) SetAttributes(attributes map[string]interface{}) *UpdateProviderTypeInstanceOptions {
	_options.Attributes = attributes
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateProviderTypeInstanceOptions) SetXCorrelationID(xCorrelationID string) *UpdateProviderTypeInstanceOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *UpdateProviderTypeInstanceOptions) SetXRequestID(xRequestID string) *UpdateProviderTypeInstanceOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateProviderTypeInstanceOptions) SetHeaders(param map[string]string) *UpdateProviderTypeInstanceOptions {
	options.Headers = param
	return options
}

// UpdateSettingsOptions : The UpdateSettings options.
type UpdateSettingsOptions struct {
	// ID of the instance
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The Event Notifications settings.
	EventNotifications *EventNotifications `json:"event_notifications,omitempty"`

	// The Cloud Object Storage settings.
	ObjectStorage *ObjectStorage `json:"object_storage,omitempty"`

	// The supplied or generated value of this header is logged for a request, and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this header is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The supplied or generated value of this header is logged for a request, and repeated in a response header  for the
	// corresponding response.  The same value is not used for downstream requests and retries of those requests.  If a
	// value of this header is not supplied in a request, the service generates a random (version 4) UUID.
	XRequestID *string `json:"X-Request-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSettingsOptions : Instantiate UpdateSettingsOptions
func (*SecurityAndComplianceCenterApiV3) NewUpdateSettingsOptions(instanceID string) *UpdateSettingsOptions {
	return &UpdateSettingsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateSettingsOptions) SetInstanceID(instanceID string) *UpdateSettingsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetEventNotifications : Allow user to set EventNotifications
func (_options *UpdateSettingsOptions) SetEventNotifications(eventNotifications *EventNotifications) *UpdateSettingsOptions {
	_options.EventNotifications = eventNotifications
	return _options
}

// SetObjectStorage : Allow user to set ObjectStorage
func (_options *UpdateSettingsOptions) SetObjectStorage(objectStorage *ObjectStorage) *UpdateSettingsOptions {
	_options.ObjectStorage = objectStorage
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UpdateSettingsOptions) SetXCorrelationID(xCorrelationID string) *UpdateSettingsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *UpdateSettingsOptions) SetXRequestID(xRequestID string) *UpdateSettingsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSettingsOptions) SetHeaders(param map[string]string) *UpdateSettingsOptions {
	options.Headers = param
	return options
}

// RequiredConfigItemsRequiredConfigAnd : RequiredConfigItemsRequiredConfigAnd struct
// This model "extends" RequiredConfigItems
type RequiredConfigItemsRequiredConfigAnd struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The `AND` required configurations.
	And []RequiredConfigItemsIntf `json:"and,omitempty"`
}

func (*RequiredConfigItemsRequiredConfigAnd) isaRequiredConfigItems() bool {
	return true
}

// UnmarshalRequiredConfigItemsRequiredConfigAnd unmarshals an instance of RequiredConfigItemsRequiredConfigAnd from the specified map of raw messages.
func UnmarshalRequiredConfigItemsRequiredConfigAnd(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigItemsRequiredConfigAnd)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRequiredConfigItems)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigItemsRequiredConfigBase : The required configuration base object.
// This model "extends" RequiredConfigItems
type RequiredConfigItemsRequiredConfigBase struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The property.
	Property *string `json:"property" validate:"required"`

	// The operator.
	Operator *string `json:"operator" validate:"required"`

	// Schema for any JSON type.
	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the RequiredConfigItemsRequiredConfigBase.Operator property.
// The operator.
const (
	RequiredConfigItemsRequiredConfigBase_Operator_DaysLessThan         = "days_less_than"
	RequiredConfigItemsRequiredConfigBase_Operator_IpsEquals            = "ips_equals"
	RequiredConfigItemsRequiredConfigBase_Operator_IpsInRange           = "ips_in_range"
	RequiredConfigItemsRequiredConfigBase_Operator_IpsNotEquals         = "ips_not_equals"
	RequiredConfigItemsRequiredConfigBase_Operator_IsEmpty              = "is_empty"
	RequiredConfigItemsRequiredConfigBase_Operator_IsFalse              = "is_false"
	RequiredConfigItemsRequiredConfigBase_Operator_IsNotEmpty           = "is_not_empty"
	RequiredConfigItemsRequiredConfigBase_Operator_IsTrue               = "is_true"
	RequiredConfigItemsRequiredConfigBase_Operator_NumEquals            = "num_equals"
	RequiredConfigItemsRequiredConfigBase_Operator_NumGreaterThan       = "num_greater_than"
	RequiredConfigItemsRequiredConfigBase_Operator_NumGreaterThanEquals = "num_greater_than_equals"
	RequiredConfigItemsRequiredConfigBase_Operator_NumLessThan          = "num_less_than"
	RequiredConfigItemsRequiredConfigBase_Operator_NumLessThanEquals    = "num_less_than_equals"
	RequiredConfigItemsRequiredConfigBase_Operator_NumNotEquals         = "num_not_equals"
	RequiredConfigItemsRequiredConfigBase_Operator_StringContains       = "string_contains"
	RequiredConfigItemsRequiredConfigBase_Operator_StringEquals         = "string_equals"
	RequiredConfigItemsRequiredConfigBase_Operator_StringMatch          = "string_match"
	RequiredConfigItemsRequiredConfigBase_Operator_StringNotContains    = "string_not_contains"
	RequiredConfigItemsRequiredConfigBase_Operator_StringNotEquals      = "string_not_equals"
	RequiredConfigItemsRequiredConfigBase_Operator_StringNotMatch       = "string_not_match"
	RequiredConfigItemsRequiredConfigBase_Operator_StringsAllowed       = "strings_allowed"
	RequiredConfigItemsRequiredConfigBase_Operator_StringsInList        = "strings_in_list"
	RequiredConfigItemsRequiredConfigBase_Operator_StringsRequired      = "strings_required"
)

// NewRequiredConfigItemsRequiredConfigBase : Instantiate RequiredConfigItemsRequiredConfigBase (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewRequiredConfigItemsRequiredConfigBase(property string, operator string) (_model *RequiredConfigItemsRequiredConfigBase, err error) {
	_model = &RequiredConfigItemsRequiredConfigBase{
		Property: core.StringPtr(property),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RequiredConfigItemsRequiredConfigBase) isaRequiredConfigItems() bool {
	return true
}

// UnmarshalRequiredConfigItemsRequiredConfigBase unmarshals an instance of RequiredConfigItemsRequiredConfigBase from the specified map of raw messages.
func UnmarshalRequiredConfigItemsRequiredConfigBase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigItemsRequiredConfigBase)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// RequiredConfigItemsRequiredConfigOr : The `OR` required configurations.
// This model "extends" RequiredConfigItems
type RequiredConfigItemsRequiredConfigOr struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The `OR` required configurations.
	Or []RequiredConfigItemsIntf `json:"or,omitempty"`
}

func (*RequiredConfigItemsRequiredConfigOr) isaRequiredConfigItems() bool {
	return true
}

// UnmarshalRequiredConfigItemsRequiredConfigOr unmarshals an instance of RequiredConfigItemsRequiredConfigOr from the specified map of raw messages.
func UnmarshalRequiredConfigItemsRequiredConfigOr(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigItemsRequiredConfigOr)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRequiredConfigItems)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigRequiredConfigAnd : RequiredConfigRequiredConfigAnd struct
// This model "extends" RequiredConfig
type RequiredConfigRequiredConfigAnd struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The `AND` required configurations.
	And []RequiredConfigItemsIntf `json:"and,omitempty"`
}

func (*RequiredConfigRequiredConfigAnd) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigRequiredConfigAnd unmarshals an instance of RequiredConfigRequiredConfigAnd from the specified map of raw messages.
func UnmarshalRequiredConfigRequiredConfigAnd(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigRequiredConfigAnd)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRequiredConfigItems)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigRequiredConfigBase : The required configuration base object.
// This model "extends" RequiredConfig
type RequiredConfigRequiredConfigBase struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The property.
	Property *string `json:"property" validate:"required"`

	// The operator.
	Operator *string `json:"operator" validate:"required"`

	// Schema for any JSON type.
	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the RequiredConfigRequiredConfigBase.Operator property.
// The operator.
const (
	RequiredConfigRequiredConfigBase_Operator_DaysLessThan         = "days_less_than"
	RequiredConfigRequiredConfigBase_Operator_IpsEquals            = "ips_equals"
	RequiredConfigRequiredConfigBase_Operator_IpsInRange           = "ips_in_range"
	RequiredConfigRequiredConfigBase_Operator_IpsNotEquals         = "ips_not_equals"
	RequiredConfigRequiredConfigBase_Operator_IsEmpty              = "is_empty"
	RequiredConfigRequiredConfigBase_Operator_IsFalse              = "is_false"
	RequiredConfigRequiredConfigBase_Operator_IsNotEmpty           = "is_not_empty"
	RequiredConfigRequiredConfigBase_Operator_IsTrue               = "is_true"
	RequiredConfigRequiredConfigBase_Operator_NumEquals            = "num_equals"
	RequiredConfigRequiredConfigBase_Operator_NumGreaterThan       = "num_greater_than"
	RequiredConfigRequiredConfigBase_Operator_NumGreaterThanEquals = "num_greater_than_equals"
	RequiredConfigRequiredConfigBase_Operator_NumLessThan          = "num_less_than"
	RequiredConfigRequiredConfigBase_Operator_NumLessThanEquals    = "num_less_than_equals"
	RequiredConfigRequiredConfigBase_Operator_NumNotEquals         = "num_not_equals"
	RequiredConfigRequiredConfigBase_Operator_StringContains       = "string_contains"
	RequiredConfigRequiredConfigBase_Operator_StringEquals         = "string_equals"
	RequiredConfigRequiredConfigBase_Operator_StringMatch          = "string_match"
	RequiredConfigRequiredConfigBase_Operator_StringNotContains    = "string_not_contains"
	RequiredConfigRequiredConfigBase_Operator_StringNotEquals      = "string_not_equals"
	RequiredConfigRequiredConfigBase_Operator_StringNotMatch       = "string_not_match"
	RequiredConfigRequiredConfigBase_Operator_StringsAllowed       = "strings_allowed"
	RequiredConfigRequiredConfigBase_Operator_StringsInList        = "strings_in_list"
	RequiredConfigRequiredConfigBase_Operator_StringsRequired      = "strings_required"
)

// NewRequiredConfigRequiredConfigBase : Instantiate RequiredConfigRequiredConfigBase (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewRequiredConfigRequiredConfigBase(property string, operator string) (_model *RequiredConfigRequiredConfigBase, err error) {
	_model = &RequiredConfigRequiredConfigBase{
		Property: core.StringPtr(property),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RequiredConfigRequiredConfigBase) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigRequiredConfigBase unmarshals an instance of RequiredConfigRequiredConfigBase from the specified map of raw messages.
func UnmarshalRequiredConfigRequiredConfigBase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigRequiredConfigBase)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// RequiredConfigRequiredConfigOr : The `OR` required configurations.
// This model "extends" RequiredConfig
type RequiredConfigRequiredConfigOr struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The `OR` required configurations.
	Or []RequiredConfigItemsIntf `json:"or,omitempty"`
}

func (*RequiredConfigRequiredConfigOr) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigRequiredConfigOr unmarshals an instance of RequiredConfigRequiredConfigOr from the specified map of raw messages.
func UnmarshalRequiredConfigRequiredConfigOr(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigRequiredConfigOr)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRequiredConfigItems)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlLibrariesPager can be used to simplify the use of the "ListControlLibraries" method.
type ControlLibrariesPager struct {
	hasNext     bool
	options     *ListControlLibrariesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewControlLibrariesPager returns a new ControlLibrariesPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewControlLibrariesPager(options *ListControlLibrariesOptions) (pager *ControlLibrariesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListControlLibrariesOptions = *options
	pager = &ControlLibrariesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ControlLibrariesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ControlLibrariesPager) GetNextWithContext(ctx context.Context) (page []ControlLibraryItem, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListControlLibrariesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.ControlLibraries

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ControlLibrariesPager) GetAllWithContext(ctx context.Context) (allItems []ControlLibraryItem, err error) {
	for pager.HasNext() {
		var nextPage []ControlLibraryItem
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ControlLibrariesPager) GetNext() (page []ControlLibraryItem, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ControlLibrariesPager) GetAll() (allItems []ControlLibraryItem, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ProfilePager can be used to simplify the use of the "ListProfile" method.
type ProfilesPager struct {
	hasNext     bool
	options     *ListProfilesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewProfilesPager returns a new ProfilesPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewProfilesPager(options *ListProfilesOptions) (pager *ProfilesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListProfilesOptions = *options
	pager = &ProfilesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ProfilesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ProfilesPager) GetNextWithContext(ctx context.Context) (page []ProfileItem, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListProfilesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Profiles

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ProfilesPager) GetAllWithContext(ctx context.Context) (allItems []ProfileItem, err error) {
	for pager.HasNext() {
		var nextPage []ProfileItem
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ProfilesPager) GetNext() (page []ProfileItem, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ProfilesPager) GetAll() (allItems []ProfileItem, err error) {
	return pager.GetAllWithContext(context.Background())
}

// AttachmentsPager can be used to simplify the use of the "ListAttachments" method.
type AttachmentsPager struct {
	hasNext     bool
	options     *ListAttachmentsOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewAttachmentsPager returns a new AttachmentsPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewAttachmentsPager(options *ListAttachmentsOptions) (pager *AttachmentsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListAttachmentsOptions = *options
	pager = &AttachmentsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AttachmentsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AttachmentsPager) GetNextWithContext(ctx context.Context) (page []AttachmentItem, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListAttachmentsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Attachments

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AttachmentsPager) GetAllWithContext(ctx context.Context) (allItems []AttachmentItem, err error) {
	for pager.HasNext() {
		var nextPage []AttachmentItem
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *AttachmentsPager) GetNext() (page []AttachmentItem, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AttachmentsPager) GetAll() (allItems []AttachmentItem, err error) {
	return pager.GetAllWithContext(context.Background())
}

// AttachmentsAccountPager can be used to simplify the use of the "ListAttachmentsAccount" method.
type AttachmentsAccountPager struct {
	hasNext     bool
	options     *ListAttachmentsAccountOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewAttachmentsAccountPager returns a new AttachmentsAccountPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewAttachmentsAccountPager(options *ListAttachmentsAccountOptions) (pager *AttachmentsAccountPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListAttachmentsAccountOptions = *options
	pager = &AttachmentsAccountPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AttachmentsAccountPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AttachmentsAccountPager) GetNextWithContext(ctx context.Context) (page []AttachmentItem, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListAttachmentsAccountWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Attachments

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AttachmentsAccountPager) GetAllWithContext(ctx context.Context) (allItems []AttachmentItem, err error) {
	for pager.HasNext() {
		var nextPage []AttachmentItem
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *AttachmentsAccountPager) GetNext() (page []AttachmentItem, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AttachmentsAccountPager) GetAll() (allItems []AttachmentItem, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ReportsPager can be used to simplify the use of the "ListReports" method.
type ReportsPager struct {
	hasNext     bool
	options     *ListReportsOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewReportsPager returns a new ReportsPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewReportsPager(options *ListReportsOptions) (pager *ReportsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListReportsOptions = *options
	pager = &ReportsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ReportsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ReportsPager) GetNextWithContext(ctx context.Context) (page []Report, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListReportsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		var start *string
		start, err = core.GetQueryParam(result.Next.Href, "start")
		if err != nil {
			err = fmt.Errorf("error retrieving 'start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Reports

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ReportsPager) GetAllWithContext(ctx context.Context) (allItems []Report, err error) {
	for pager.HasNext() {
		var nextPage []Report
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ReportsPager) GetNext() (page []Report, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ReportsPager) GetAll() (allItems []Report, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ReportEvaluationsPager can be used to simplify the use of the "ListReportEvaluations" method.
type ReportEvaluationsPager struct {
	hasNext     bool
	options     *ListReportEvaluationsOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewReportEvaluationsPager returns a new ReportEvaluationsPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewReportEvaluationsPager(options *ListReportEvaluationsOptions) (pager *ReportEvaluationsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListReportEvaluationsOptions = *options
	pager = &ReportEvaluationsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ReportEvaluationsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ReportEvaluationsPager) GetNextWithContext(ctx context.Context) (page []Evaluation, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListReportEvaluationsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		var start *string
		start, err = core.GetQueryParam(result.Next.Href, "start")
		if err != nil {
			err = fmt.Errorf("error retrieving 'start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Evaluations

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ReportEvaluationsPager) GetAllWithContext(ctx context.Context) (allItems []Evaluation, err error) {
	for pager.HasNext() {
		var nextPage []Evaluation
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ReportEvaluationsPager) GetNext() (page []Evaluation, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ReportEvaluationsPager) GetAll() (allItems []Evaluation, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ReportResourcesPager can be used to simplify the use of the "ListReportResources" method.
type ReportResourcesPager struct {
	hasNext     bool
	options     *ListReportResourcesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewReportResourcesPager returns a new ReportResourcesPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewReportResourcesPager(options *ListReportResourcesOptions) (pager *ReportResourcesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListReportResourcesOptions = *options
	pager = &ReportResourcesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ReportResourcesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ReportResourcesPager) GetNextWithContext(ctx context.Context) (page []Resource, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListReportResourcesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		var start *string
		start, err = core.GetQueryParam(result.Next.Href, "start")
		if err != nil {
			err = fmt.Errorf("error retrieving 'start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ReportResourcesPager) GetAllWithContext(ctx context.Context) (allItems []Resource, err error) {
	for pager.HasNext() {
		var nextPage []Resource
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ReportResourcesPager) GetNext() (page []Resource, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ReportResourcesPager) GetAll() (allItems []Resource, err error) {
	return pager.GetAllWithContext(context.Background())
}

// getInstanceBasedURL uses the baseURL and instanceID to return the correct endpoint URL
func getInstanceBasedURL(baseURL, instanceID string) (string, error) {
	if baseURL == "" {
		return "", fmt.Errorf("service URL is empty")
	}
	return fmt.Sprintf("%s/instances/%s/v3", baseURL, instanceID), nil
}
