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
 * IBM OpenAPI SDK Code Generator Version: 3.43.0-49eab5c7-20211117-152138
 */

// Package alertsv1 : Operations and models for the AlertsV1 service
package alertsv1

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

// AlertsV1 : CIS Alert Policies
//
// API Version: 1.0.0
type AlertsV1 struct {
	Service *core.BaseService

	// Full url-encoded CRN of the service instance.
	Crn *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "alerts"

// AlertsV1Options : Service options
type AlertsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full url-encoded CRN of the service instance.
	Crn *string `validate:"required"`
}

// NewAlertsV1UsingExternalConfig : constructs an instance of AlertsV1 with passed in options and external configuration.
func NewAlertsV1UsingExternalConfig(options *AlertsV1Options) (alerts *AlertsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	alerts, err = NewAlertsV1(options)
	if err != nil {
		return
	}

	err = alerts.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = alerts.Service.SetServiceURL(options.URL)
	}
	return
}

// NewAlertsV1 : constructs an instance of AlertsV1 with passed in options.
func NewAlertsV1(options *AlertsV1Options) (service *AlertsV1, err error) {
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

	service = &AlertsV1{
		Service: baseService,
		Crn:     options.Crn,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "alerts" suitable for processing requests.
func (alerts *AlertsV1) Clone() *AlertsV1 {
	if core.IsNil(alerts) {
		return nil
	}
	clone := *alerts
	clone.Service = alerts.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (alerts *AlertsV1) SetServiceURL(url string) error {
	return alerts.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (alerts *AlertsV1) GetServiceURL() string {
	return alerts.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (alerts *AlertsV1) SetDefaultHeaders(headers http.Header) {
	alerts.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (alerts *AlertsV1) SetEnableGzipCompression(enableGzip bool) {
	alerts.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (alerts *AlertsV1) GetEnableGzipCompression() bool {
	return alerts.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (alerts *AlertsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	alerts.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (alerts *AlertsV1) DisableRetries() {
	alerts.Service.DisableRetries()
}

// GetAlertPolicies : List alert policies
// List configured alert policies for the CIS instance.
func (alerts *AlertsV1) GetAlertPolicies(getAlertPoliciesOptions *GetAlertPoliciesOptions) (result *ListAlertPoliciesResp, response *core.DetailedResponse, err error) {
	return alerts.GetAlertPoliciesWithContext(context.Background(), getAlertPoliciesOptions)
}

// GetAlertPoliciesWithContext is an alternate form of the GetAlertPolicies method which supports a Context parameter
func (alerts *AlertsV1) GetAlertPoliciesWithContext(ctx context.Context, getAlertPoliciesOptions *GetAlertPoliciesOptions) (result *ListAlertPoliciesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getAlertPoliciesOptions, "getAlertPoliciesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *alerts.Crn,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = alerts.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(alerts.Service.Options.URL, `/v1/{crn}/alerting/policies`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAlertPoliciesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("alerts", "V1", "GetAlertPolicies")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = alerts.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListAlertPoliciesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAlertPolicy : Create an alert policy
// Create a new alert policy for the CIS instance.
func (alerts *AlertsV1) CreateAlertPolicy(createAlertPolicyOptions *CreateAlertPolicyOptions) (result *AlertSuccessResp, response *core.DetailedResponse, err error) {
	return alerts.CreateAlertPolicyWithContext(context.Background(), createAlertPolicyOptions)
}

// CreateAlertPolicyWithContext is an alternate form of the CreateAlertPolicy method which supports a Context parameter
func (alerts *AlertsV1) CreateAlertPolicyWithContext(ctx context.Context, createAlertPolicyOptions *CreateAlertPolicyOptions) (result *AlertSuccessResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createAlertPolicyOptions, "createAlertPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *alerts.Crn,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = alerts.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(alerts.Service.Options.URL, `/v1/{crn}/alerting/policies`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createAlertPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("alerts", "V1", "CreateAlertPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAlertPolicyOptions.Name != nil {
		body["name"] = createAlertPolicyOptions.Name
	}
	if createAlertPolicyOptions.Description != nil {
		body["description"] = createAlertPolicyOptions.Description
	}
	if createAlertPolicyOptions.Enabled != nil {
		body["enabled"] = createAlertPolicyOptions.Enabled
	}
	if createAlertPolicyOptions.AlertType != nil {
		body["alert_type"] = createAlertPolicyOptions.AlertType
	}
	if createAlertPolicyOptions.Mechanisms != nil {
		body["mechanisms"] = createAlertPolicyOptions.Mechanisms
	}
	if createAlertPolicyOptions.Filters != nil {
		body["filters"] = createAlertPolicyOptions.Filters
	}
	if createAlertPolicyOptions.Conditions != nil {
		body["conditions"] = createAlertPolicyOptions.Conditions
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
	response, err = alerts.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAlertSuccessResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAlertPolicy : Get an alert policy
// Get an alert policy for the CIS instance.
func (alerts *AlertsV1) GetAlertPolicy(getAlertPolicyOptions *GetAlertPolicyOptions) (result *GetAlertPolicyResp, response *core.DetailedResponse, err error) {
	return alerts.GetAlertPolicyWithContext(context.Background(), getAlertPolicyOptions)
}

// GetAlertPolicyWithContext is an alternate form of the GetAlertPolicy method which supports a Context parameter
func (alerts *AlertsV1) GetAlertPolicyWithContext(ctx context.Context, getAlertPolicyOptions *GetAlertPolicyOptions) (result *GetAlertPolicyResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAlertPolicyOptions, "getAlertPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAlertPolicyOptions, "getAlertPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":       *alerts.Crn,
		"policy_id": *getAlertPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = alerts.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(alerts.Service.Options.URL, `/v1/{crn}/alerting/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAlertPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("alerts", "V1", "GetAlertPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = alerts.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetAlertPolicyResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateAlertPolicy : Update an alert policy
// Update an existing alert policy for the CIS instance.
func (alerts *AlertsV1) UpdateAlertPolicy(updateAlertPolicyOptions *UpdateAlertPolicyOptions) (result *AlertSuccessResp, response *core.DetailedResponse, err error) {
	return alerts.UpdateAlertPolicyWithContext(context.Background(), updateAlertPolicyOptions)
}

// UpdateAlertPolicyWithContext is an alternate form of the UpdateAlertPolicy method which supports a Context parameter
func (alerts *AlertsV1) UpdateAlertPolicyWithContext(ctx context.Context, updateAlertPolicyOptions *UpdateAlertPolicyOptions) (result *AlertSuccessResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAlertPolicyOptions, "updateAlertPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAlertPolicyOptions, "updateAlertPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":       *alerts.Crn,
		"policy_id": *updateAlertPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = alerts.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(alerts.Service.Options.URL, `/v1/{crn}/alerting/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateAlertPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("alerts", "V1", "UpdateAlertPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAlertPolicyOptions.Name != nil {
		body["name"] = updateAlertPolicyOptions.Name
	}
	if updateAlertPolicyOptions.Description != nil {
		body["description"] = updateAlertPolicyOptions.Description
	}
	if updateAlertPolicyOptions.Enabled != nil {
		body["enabled"] = updateAlertPolicyOptions.Enabled
	}
	if updateAlertPolicyOptions.AlertType != nil {
		body["alert_type"] = updateAlertPolicyOptions.AlertType
	}
	if updateAlertPolicyOptions.Mechanisms != nil {
		body["mechanisms"] = updateAlertPolicyOptions.Mechanisms
	}
	if updateAlertPolicyOptions.Conditions != nil {
		body["conditions"] = updateAlertPolicyOptions.Conditions
	}
	if updateAlertPolicyOptions.Filters != nil {
		body["filters"] = updateAlertPolicyOptions.Filters
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
	response, err = alerts.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAlertSuccessResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteAlertPolicy : Delete an alert policy
// Delete an alert policy for the CIS instance.
func (alerts *AlertsV1) DeleteAlertPolicy(deleteAlertPolicyOptions *DeleteAlertPolicyOptions) (result *AlertSuccessResp, response *core.DetailedResponse, err error) {
	return alerts.DeleteAlertPolicyWithContext(context.Background(), deleteAlertPolicyOptions)
}

// DeleteAlertPolicyWithContext is an alternate form of the DeleteAlertPolicy method which supports a Context parameter
func (alerts *AlertsV1) DeleteAlertPolicyWithContext(ctx context.Context, deleteAlertPolicyOptions *DeleteAlertPolicyOptions) (result *AlertSuccessResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAlertPolicyOptions, "deleteAlertPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteAlertPolicyOptions, "deleteAlertPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":       *alerts.Crn,
		"policy_id": *deleteAlertPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = alerts.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(alerts.Service.Options.URL, `/v1/{crn}/alerting/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteAlertPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("alerts", "V1", "DeleteAlertPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = alerts.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAlertSuccessResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// AlertSuccessRespErrorsItem : AlertSuccessRespErrorsItem struct
type AlertSuccessRespErrorsItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalAlertSuccessRespErrorsItem unmarshals an instance of AlertSuccessRespErrorsItem from the specified map of raw messages.
func UnmarshalAlertSuccessRespErrorsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AlertSuccessRespErrorsItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AlertSuccessRespMessagesItem : AlertSuccessRespMessagesItem struct
type AlertSuccessRespMessagesItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalAlertSuccessRespMessagesItem unmarshals an instance of AlertSuccessRespMessagesItem from the specified map of raw messages.
func UnmarshalAlertSuccessRespMessagesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AlertSuccessRespMessagesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AlertSuccessRespResult : Container for response information.
type AlertSuccessRespResult struct {
	// Policy ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalAlertSuccessRespResult unmarshals an instance of AlertSuccessRespResult from the specified map of raw messages.
func UnmarshalAlertSuccessRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AlertSuccessRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAlertPolicyInputMechanisms : Delivery mechanisms for the alert.
type CreateAlertPolicyInputMechanisms struct {
	Email []CreateAlertPolicyInputMechanismsEmailItem `json:"email,omitempty"`

	Webhooks []CreateAlertPolicyInputMechanismsWebhooksItem `json:"webhooks,omitempty"`
}

// UnmarshalCreateAlertPolicyInputMechanisms unmarshals an instance of CreateAlertPolicyInputMechanisms from the specified map of raw messages.
func UnmarshalCreateAlertPolicyInputMechanisms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAlertPolicyInputMechanisms)
	err = core.UnmarshalModel(m, "email", &obj.Email, UnmarshalCreateAlertPolicyInputMechanismsEmailItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "webhooks", &obj.Webhooks, UnmarshalCreateAlertPolicyInputMechanismsWebhooksItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAlertPolicyInputMechanismsEmailItem : CreateAlertPolicyInputMechanismsEmailItem struct
type CreateAlertPolicyInputMechanismsEmailItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalCreateAlertPolicyInputMechanismsEmailItem unmarshals an instance of CreateAlertPolicyInputMechanismsEmailItem from the specified map of raw messages.
func UnmarshalCreateAlertPolicyInputMechanismsEmailItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAlertPolicyInputMechanismsEmailItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAlertPolicyInputMechanismsWebhooksItem : CreateAlertPolicyInputMechanismsWebhooksItem struct
type CreateAlertPolicyInputMechanismsWebhooksItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalCreateAlertPolicyInputMechanismsWebhooksItem unmarshals an instance of CreateAlertPolicyInputMechanismsWebhooksItem from the specified map of raw messages.
func UnmarshalCreateAlertPolicyInputMechanismsWebhooksItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAlertPolicyInputMechanismsWebhooksItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAlertPolicyOptions : The CreateAlertPolicy options.
type CreateAlertPolicyOptions struct {
	// Policy name.
	Name *string `json:"name,omitempty"`

	// Policy description.
	Description *string `json:"description,omitempty"`

	// Is the alert policy active.
	Enabled *bool `json:"enabled,omitempty"`

	// Condition for the alert.
	AlertType *string `json:"alert_type,omitempty"`

	// Delivery mechanisms for the alert.
	Mechanisms *CreateAlertPolicyInputMechanisms `json:"mechanisms,omitempty"`

	// Optional filters depending for the alert type.
	Filters interface{} `json:"filters,omitempty"`

	// Conditions depending on the alert type. HTTP DDOS Attack Alerter does not have any conditions. The Load Balancing
	// Pool Enablement Alerter takes conditions that describe for all pools whether the pool is being enabled, disabled, or
	// both. This field is not required when creating a new alert.
	Conditions interface{} `json:"conditions,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateAlertPolicyOptions.AlertType property.
// Condition for the alert.
const (
	CreateAlertPolicyOptions_AlertType_ClickhouseAlertFwAnomaly    = "clickhouse_alert_fw_anomaly"
	CreateAlertPolicyOptions_AlertType_ClickhouseAlertFwEntAnomaly = "clickhouse_alert_fw_ent_anomaly"
	CreateAlertPolicyOptions_AlertType_DosAttackL7                 = "dos_attack_l7"
	CreateAlertPolicyOptions_AlertType_G6PoolToggleAlert           = "g6_pool_toggle_alert"
)

// NewCreateAlertPolicyOptions : Instantiate CreateAlertPolicyOptions
func (*AlertsV1) NewCreateAlertPolicyOptions() *CreateAlertPolicyOptions {
	return &CreateAlertPolicyOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateAlertPolicyOptions) SetName(name string) *CreateAlertPolicyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateAlertPolicyOptions) SetDescription(description string) *CreateAlertPolicyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateAlertPolicyOptions) SetEnabled(enabled bool) *CreateAlertPolicyOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetAlertType : Allow user to set AlertType
func (_options *CreateAlertPolicyOptions) SetAlertType(alertType string) *CreateAlertPolicyOptions {
	_options.AlertType = core.StringPtr(alertType)
	return _options
}

// SetMechanisms : Allow user to set Mechanisms
func (_options *CreateAlertPolicyOptions) SetMechanisms(mechanisms *CreateAlertPolicyInputMechanisms) *CreateAlertPolicyOptions {
	_options.Mechanisms = mechanisms
	return _options
}

// SetFilters : Allow user to set Filters
func (_options *CreateAlertPolicyOptions) SetFilters(filters interface{}) *CreateAlertPolicyOptions {
	_options.Filters = filters
	return _options
}

// SetConditions : Allow user to set Conditions
func (_options *CreateAlertPolicyOptions) SetConditions(conditions interface{}) *CreateAlertPolicyOptions {
	_options.Conditions = conditions
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAlertPolicyOptions) SetHeaders(param map[string]string) *CreateAlertPolicyOptions {
	options.Headers = param
	return options
}

// DeleteAlertPolicyOptions : The DeleteAlertPolicy options.
type DeleteAlertPolicyOptions struct {
	// Alert policy identifier.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteAlertPolicyOptions : Instantiate DeleteAlertPolicyOptions
func (*AlertsV1) NewDeleteAlertPolicyOptions(policyID string) *DeleteAlertPolicyOptions {
	return &DeleteAlertPolicyOptions{
		PolicyID: core.StringPtr(policyID),
	}
}

// SetPolicyID : Allow user to set PolicyID
func (_options *DeleteAlertPolicyOptions) SetPolicyID(policyID string) *DeleteAlertPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAlertPolicyOptions) SetHeaders(param map[string]string) *DeleteAlertPolicyOptions {
	options.Headers = param
	return options
}

// GetAlertPoliciesOptions : The GetAlertPolicies options.
type GetAlertPoliciesOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAlertPoliciesOptions : Instantiate GetAlertPoliciesOptions
func (*AlertsV1) NewGetAlertPoliciesOptions() *GetAlertPoliciesOptions {
	return &GetAlertPoliciesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetAlertPoliciesOptions) SetHeaders(param map[string]string) *GetAlertPoliciesOptions {
	options.Headers = param
	return options
}

// GetAlertPolicyOptions : The GetAlertPolicy options.
type GetAlertPolicyOptions struct {
	// Alert policy identifier.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAlertPolicyOptions : Instantiate GetAlertPolicyOptions
func (*AlertsV1) NewGetAlertPolicyOptions(policyID string) *GetAlertPolicyOptions {
	return &GetAlertPolicyOptions{
		PolicyID: core.StringPtr(policyID),
	}
}

// SetPolicyID : Allow user to set PolicyID
func (_options *GetAlertPolicyOptions) SetPolicyID(policyID string) *GetAlertPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAlertPolicyOptions) SetHeaders(param map[string]string) *GetAlertPolicyOptions {
	options.Headers = param
	return options
}

// GetAlertPolicyRespErrorsItem : GetAlertPolicyRespErrorsItem struct
type GetAlertPolicyRespErrorsItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalGetAlertPolicyRespErrorsItem unmarshals an instance of GetAlertPolicyRespErrorsItem from the specified map of raw messages.
func UnmarshalGetAlertPolicyRespErrorsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertPolicyRespErrorsItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAlertPolicyRespMessagesItem : GetAlertPolicyRespMessagesItem struct
type GetAlertPolicyRespMessagesItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalGetAlertPolicyRespMessagesItem unmarshals an instance of GetAlertPolicyRespMessagesItem from the specified map of raw messages.
func UnmarshalGetAlertPolicyRespMessagesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertPolicyRespMessagesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAlertPolicyRespResult : Container for response information.
type GetAlertPolicyRespResult struct {
	// Policy ID.
	ID *string `json:"id" validate:"required"`

	// Policy Name.
	Name *string `json:"name" validate:"required"`

	// Alert Policy description.
	Description *string `json:"description" validate:"required"`

	// Is the alert enabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// Condition for the alert.
	AlertType *string `json:"alert_type" validate:"required"`

	// Delivery mechanisms for the alert, can include an email, a webhook, or both.
	Mechanisms *GetAlertPolicyRespResultMechanisms `json:"mechanisms" validate:"required"`

	// When was the policy first created.
	Created *string `json:"created" validate:"required"`

	// When was the policy last modified.
	Modified *string `json:"modified" validate:"required"`

	// Optional conditions depending for the alert type.
	Conditions interface{} `json:"conditions" validate:"required"`

	// Optional filters depending for the alert type.
	Filters interface{} `json:"filters" validate:"required"`
}

// Constants associated with the GetAlertPolicyRespResult.AlertType property.
// Condition for the alert.
const (
	GetAlertPolicyRespResult_AlertType_ClickhouseAlertFwAnomaly    = "clickhouse_alert_fw_anomaly"
	GetAlertPolicyRespResult_AlertType_ClickhouseAlertFwEntAnomaly = "clickhouse_alert_fw_ent_anomaly"
	GetAlertPolicyRespResult_AlertType_DosAttackL7                 = "dos_attack_l7"
	GetAlertPolicyRespResult_AlertType_G6PoolToggleAlert           = "g6_pool_toggle_alert"
)

// UnmarshalGetAlertPolicyRespResult unmarshals an instance of GetAlertPolicyRespResult from the specified map of raw messages.
func UnmarshalGetAlertPolicyRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertPolicyRespResult)
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
	err = core.UnmarshalPrimitive(m, "alert_type", &obj.AlertType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "mechanisms", &obj.Mechanisms, UnmarshalGetAlertPolicyRespResultMechanisms)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified", &obj.Modified)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "conditions", &obj.Conditions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "filters", &obj.Filters)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAlertPolicyRespResultMechanisms : Delivery mechanisms for the alert, can include an email, a webhook, or both.
type GetAlertPolicyRespResultMechanisms struct {
	Email []GetAlertPolicyRespResultMechanismsEmailItem `json:"email,omitempty"`

	Webhooks []GetAlertPolicyRespResultMechanismsWebhooksItem `json:"webhooks,omitempty"`
}

// UnmarshalGetAlertPolicyRespResultMechanisms unmarshals an instance of GetAlertPolicyRespResultMechanisms from the specified map of raw messages.
func UnmarshalGetAlertPolicyRespResultMechanisms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertPolicyRespResultMechanisms)
	err = core.UnmarshalModel(m, "email", &obj.Email, UnmarshalGetAlertPolicyRespResultMechanismsEmailItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "webhooks", &obj.Webhooks, UnmarshalGetAlertPolicyRespResultMechanismsWebhooksItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAlertPolicyRespResultMechanismsEmailItem : GetAlertPolicyRespResultMechanismsEmailItem struct
type GetAlertPolicyRespResultMechanismsEmailItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalGetAlertPolicyRespResultMechanismsEmailItem unmarshals an instance of GetAlertPolicyRespResultMechanismsEmailItem from the specified map of raw messages.
func UnmarshalGetAlertPolicyRespResultMechanismsEmailItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertPolicyRespResultMechanismsEmailItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAlertPolicyRespResultMechanismsWebhooksItem : GetAlertPolicyRespResultMechanismsWebhooksItem struct
type GetAlertPolicyRespResultMechanismsWebhooksItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalGetAlertPolicyRespResultMechanismsWebhooksItem unmarshals an instance of GetAlertPolicyRespResultMechanismsWebhooksItem from the specified map of raw messages.
func UnmarshalGetAlertPolicyRespResultMechanismsWebhooksItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertPolicyRespResultMechanismsWebhooksItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAlertPoliciesRespErrorsItem : ListAlertPoliciesRespErrorsItem struct
type ListAlertPoliciesRespErrorsItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalListAlertPoliciesRespErrorsItem unmarshals an instance of ListAlertPoliciesRespErrorsItem from the specified map of raw messages.
func UnmarshalListAlertPoliciesRespErrorsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertPoliciesRespErrorsItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAlertPoliciesRespMessagesItem : ListAlertPoliciesRespMessagesItem struct
type ListAlertPoliciesRespMessagesItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalListAlertPoliciesRespMessagesItem unmarshals an instance of ListAlertPoliciesRespMessagesItem from the specified map of raw messages.
func UnmarshalListAlertPoliciesRespMessagesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertPoliciesRespMessagesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAlertPoliciesRespResultItem : ListAlertPoliciesRespResultItem struct
type ListAlertPoliciesRespResultItem struct {
	// Policy ID.
	ID *string `json:"id" validate:"required"`

	// Policy Name.
	Name *string `json:"name" validate:"required"`

	// Alert Policy description.
	Description *string `json:"description" validate:"required"`

	// Is the alert enabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// Condition for the alert.
	AlertType *string `json:"alert_type" validate:"required"`

	// Delivery mechanisms for the alert, can include an email, a webhook, or both.
	Mechanisms *ListAlertPoliciesRespResultItemMechanisms `json:"mechanisms" validate:"required"`

	// When was the policy first created.
	Created *string `json:"created" validate:"required"`

	// When was the policy last modified.
	Modified *string `json:"modified" validate:"required"`

	// Optional conditions depending for the alert type.
	Conditions interface{} `json:"conditions" validate:"required"`

	// Optional filters depending for the alert type.
	Filters interface{} `json:"filters" validate:"required"`
}

// Constants associated with the ListAlertPoliciesRespResultItem.AlertType property.
// Condition for the alert.
const (
	ListAlertPoliciesRespResultItem_AlertType_ClickhouseAlertFwAnomaly    = "clickhouse_alert_fw_anomaly"
	ListAlertPoliciesRespResultItem_AlertType_ClickhouseAlertFwEntAnomaly = "clickhouse_alert_fw_ent_anomaly"
	ListAlertPoliciesRespResultItem_AlertType_DosAttackL7                 = "dos_attack_l7"
	ListAlertPoliciesRespResultItem_AlertType_G6PoolToggleAlert           = "g6_pool_toggle_alert"
)

// UnmarshalListAlertPoliciesRespResultItem unmarshals an instance of ListAlertPoliciesRespResultItem from the specified map of raw messages.
func UnmarshalListAlertPoliciesRespResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertPoliciesRespResultItem)
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
	err = core.UnmarshalPrimitive(m, "alert_type", &obj.AlertType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "mechanisms", &obj.Mechanisms, UnmarshalListAlertPoliciesRespResultItemMechanisms)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified", &obj.Modified)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "conditions", &obj.Conditions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "filters", &obj.Filters)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAlertPoliciesRespResultItemMechanisms : Delivery mechanisms for the alert, can include an email, a webhook, or both.
type ListAlertPoliciesRespResultItemMechanisms struct {
	Email []ListAlertPoliciesRespResultItemMechanismsEmailItem `json:"email,omitempty"`

	Webhooks []ListAlertPoliciesRespResultItemMechanismsWebhooksItem `json:"webhooks,omitempty"`
}

// UnmarshalListAlertPoliciesRespResultItemMechanisms unmarshals an instance of ListAlertPoliciesRespResultItemMechanisms from the specified map of raw messages.
func UnmarshalListAlertPoliciesRespResultItemMechanisms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertPoliciesRespResultItemMechanisms)
	err = core.UnmarshalModel(m, "email", &obj.Email, UnmarshalListAlertPoliciesRespResultItemMechanismsEmailItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "webhooks", &obj.Webhooks, UnmarshalListAlertPoliciesRespResultItemMechanismsWebhooksItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAlertPoliciesRespResultItemMechanismsEmailItem : ListAlertPoliciesRespResultItemMechanismsEmailItem struct
type ListAlertPoliciesRespResultItemMechanismsEmailItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalListAlertPoliciesRespResultItemMechanismsEmailItem unmarshals an instance of ListAlertPoliciesRespResultItemMechanismsEmailItem from the specified map of raw messages.
func UnmarshalListAlertPoliciesRespResultItemMechanismsEmailItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertPoliciesRespResultItemMechanismsEmailItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAlertPoliciesRespResultItemMechanismsWebhooksItem : ListAlertPoliciesRespResultItemMechanismsWebhooksItem struct
type ListAlertPoliciesRespResultItemMechanismsWebhooksItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalListAlertPoliciesRespResultItemMechanismsWebhooksItem unmarshals an instance of ListAlertPoliciesRespResultItemMechanismsWebhooksItem from the specified map of raw messages.
func UnmarshalListAlertPoliciesRespResultItemMechanismsWebhooksItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertPoliciesRespResultItemMechanismsWebhooksItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateAlertPolicyInputMechanisms : Delivery mechanisms for the alert, can include an email, a webhook, or both.
type UpdateAlertPolicyInputMechanisms struct {
	Email []UpdateAlertPolicyInputMechanismsEmailItem `json:"email,omitempty"`

	Webhooks []UpdateAlertPolicyInputMechanismsWebhooksItem `json:"webhooks,omitempty"`
}

// UnmarshalUpdateAlertPolicyInputMechanisms unmarshals an instance of UpdateAlertPolicyInputMechanisms from the specified map of raw messages.
func UnmarshalUpdateAlertPolicyInputMechanisms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateAlertPolicyInputMechanisms)
	err = core.UnmarshalModel(m, "email", &obj.Email, UnmarshalUpdateAlertPolicyInputMechanismsEmailItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "webhooks", &obj.Webhooks, UnmarshalUpdateAlertPolicyInputMechanismsWebhooksItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateAlertPolicyInputMechanismsEmailItem : UpdateAlertPolicyInputMechanismsEmailItem struct
type UpdateAlertPolicyInputMechanismsEmailItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalUpdateAlertPolicyInputMechanismsEmailItem unmarshals an instance of UpdateAlertPolicyInputMechanismsEmailItem from the specified map of raw messages.
func UnmarshalUpdateAlertPolicyInputMechanismsEmailItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateAlertPolicyInputMechanismsEmailItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateAlertPolicyInputMechanismsWebhooksItem : UpdateAlertPolicyInputMechanismsWebhooksItem struct
type UpdateAlertPolicyInputMechanismsWebhooksItem struct {
	ID *string `json:"id,omitempty"`
}

// UnmarshalUpdateAlertPolicyInputMechanismsWebhooksItem unmarshals an instance of UpdateAlertPolicyInputMechanismsWebhooksItem from the specified map of raw messages.
func UnmarshalUpdateAlertPolicyInputMechanismsWebhooksItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateAlertPolicyInputMechanismsWebhooksItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateAlertPolicyOptions : The UpdateAlertPolicy options.
type UpdateAlertPolicyOptions struct {
	// Alert policy identifier.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Policy name.
	Name *string `json:"name,omitempty"`

	// Policy description.
	Description *string `json:"description,omitempty"`

	// Is the alert policy active.
	Enabled *bool `json:"enabled,omitempty"`

	// Condition for the alert. Use 'dos_attack_l7' to set up an HTTP DDOS Attack Alerter, use 'g6_pool_toggle_alert' to
	// set up Load Balancing Pool Enablement Alerter, use 'clickhouse_alert_fw_anomaly' to set up WAF Alerter and
	// 'clickhouse_alert_fw_ent_anomaly' to set up Advanced Security Alerter.
	AlertType *string `json:"alert_type,omitempty"`

	// Delivery mechanisms for the alert, can include an email, a webhook, or both.
	Mechanisms *UpdateAlertPolicyInputMechanisms `json:"mechanisms,omitempty"`

	// Conditions depending on the alert type. HTTP DDOS Attack Alerter does not have any conditions. The Load Balancing
	// Pool Enablement Alerter takes conditions that describe for all pools whether the pool is being enabled, disabled, or
	// both.
	Conditions interface{} `json:"conditions,omitempty"`

	// Optional filters depending for the alert type. HTTP DDOS Attack Alerter does not require any filters. The Load
	// Balancing Pool Enablement Alerter requires a list of IDs for the pools and their corresponding alert trigger (set
	// whether alerts are recieved on disablement, enablement, or both). The basic WAF Alerter requires a list of zones to
	// be monitored. The Advanced Security Alerter requires a list of zones to be monitored as well as a list of services
	// to monitor.
	Filters interface{} `json:"filters,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateAlertPolicyOptions.AlertType property.
// Condition for the alert. Use 'dos_attack_l7' to set up an HTTP DDOS Attack Alerter, use 'g6_pool_toggle_alert' to set
// up Load Balancing Pool Enablement Alerter, use 'clickhouse_alert_fw_anomaly' to set up WAF Alerter and
// 'clickhouse_alert_fw_ent_anomaly' to set up Advanced Security Alerter.
const (
	UpdateAlertPolicyOptions_AlertType_ClickhouseAlertFwAnomaly    = "clickhouse_alert_fw_anomaly"
	UpdateAlertPolicyOptions_AlertType_ClickhouseAlertFwEntAnomaly = "clickhouse_alert_fw_ent_anomaly"
	UpdateAlertPolicyOptions_AlertType_DosAttackL7                 = "dos_attack_l7"
	UpdateAlertPolicyOptions_AlertType_G6PoolToggleAlert           = "g6_pool_toggle_alert"
)

// NewUpdateAlertPolicyOptions : Instantiate UpdateAlertPolicyOptions
func (*AlertsV1) NewUpdateAlertPolicyOptions(policyID string) *UpdateAlertPolicyOptions {
	return &UpdateAlertPolicyOptions{
		PolicyID: core.StringPtr(policyID),
	}
}

// SetPolicyID : Allow user to set PolicyID
func (_options *UpdateAlertPolicyOptions) SetPolicyID(policyID string) *UpdateAlertPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAlertPolicyOptions) SetName(name string) *UpdateAlertPolicyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateAlertPolicyOptions) SetDescription(description string) *UpdateAlertPolicyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateAlertPolicyOptions) SetEnabled(enabled bool) *UpdateAlertPolicyOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetAlertType : Allow user to set AlertType
func (_options *UpdateAlertPolicyOptions) SetAlertType(alertType string) *UpdateAlertPolicyOptions {
	_options.AlertType = core.StringPtr(alertType)
	return _options
}

// SetMechanisms : Allow user to set Mechanisms
func (_options *UpdateAlertPolicyOptions) SetMechanisms(mechanisms *UpdateAlertPolicyInputMechanisms) *UpdateAlertPolicyOptions {
	_options.Mechanisms = mechanisms
	return _options
}

// SetConditions : Allow user to set Conditions
func (_options *UpdateAlertPolicyOptions) SetConditions(conditions interface{}) *UpdateAlertPolicyOptions {
	_options.Conditions = conditions
	return _options
}

// SetFilters : Allow user to set Filters
func (_options *UpdateAlertPolicyOptions) SetFilters(filters interface{}) *UpdateAlertPolicyOptions {
	_options.Filters = filters
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAlertPolicyOptions) SetHeaders(param map[string]string) *UpdateAlertPolicyOptions {
	options.Headers = param
	return options
}

// AlertSuccessResp : Alert Policies Response.
type AlertSuccessResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors []AlertSuccessRespErrorsItem `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages []AlertSuccessRespMessagesItem `json:"messages" validate:"required"`

	// Container for response information.
	Result *AlertSuccessRespResult `json:"result" validate:"required"`
}

// UnmarshalAlertSuccessResp unmarshals an instance of AlertSuccessResp from the specified map of raw messages.
func UnmarshalAlertSuccessResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AlertSuccessResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalAlertSuccessRespErrorsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalAlertSuccessRespMessagesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalAlertSuccessRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAlertPolicyResp : Get Alert Policies Response.
type GetAlertPolicyResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors []GetAlertPolicyRespErrorsItem `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages []GetAlertPolicyRespMessagesItem `json:"messages" validate:"required"`

	// Container for response information.
	Result *GetAlertPolicyRespResult `json:"result" validate:"required"`
}

// UnmarshalGetAlertPolicyResp unmarshals an instance of GetAlertPolicyResp from the specified map of raw messages.
func UnmarshalGetAlertPolicyResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertPolicyResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalGetAlertPolicyRespErrorsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalGetAlertPolicyRespMessagesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalGetAlertPolicyRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAlertPoliciesResp : List Alert Policies Response.
type ListAlertPoliciesResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors []ListAlertPoliciesRespErrorsItem `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages []ListAlertPoliciesRespMessagesItem `json:"messages" validate:"required"`

	// Container for response information.
	Result []ListAlertPoliciesRespResultItem `json:"result" validate:"required"`
}

// UnmarshalListAlertPoliciesResp unmarshals an instance of ListAlertPoliciesResp from the specified map of raw messages.
func UnmarshalListAlertPoliciesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertPoliciesResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalListAlertPoliciesRespErrorsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalListAlertPoliciesRespMessagesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalListAlertPoliciesRespResultItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
