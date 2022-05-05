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
 * IBM OpenAPI SDK Code Generator Version: 3.43.0-49eab5c7-20211117-152138
 */

// Package webhooksv1 : Operations and models for the WebhooksV1 service
package webhooksv1

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

// WebhooksV1 : CIS Alert Webhooks
//
// API Version: 1.0.0
type WebhooksV1 struct {
	Service *core.BaseService

	// Full url-encoded CRN of the service instance.
	Crn *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "webhooks"

// WebhooksV1Options : Service options
type WebhooksV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full url-encoded CRN of the service instance.
	Crn *string `validate:"required"`
}

// NewWebhooksV1UsingExternalConfig : constructs an instance of WebhooksV1 with passed in options and external configuration.
func NewWebhooksV1UsingExternalConfig(options *WebhooksV1Options) (webhooks *WebhooksV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	webhooks, err = NewWebhooksV1(options)
	if err != nil {
		return
	}

	err = webhooks.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = webhooks.Service.SetServiceURL(options.URL)
	}
	return
}

// NewWebhooksV1 : constructs an instance of WebhooksV1 with passed in options.
func NewWebhooksV1(options *WebhooksV1Options) (service *WebhooksV1, err error) {
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

	service = &WebhooksV1{
		Service: baseService,
		Crn: options.Crn,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "webhooks" suitable for processing requests.
func (webhooks *WebhooksV1) Clone() *WebhooksV1 {
	if core.IsNil(webhooks) {
		return nil
	}
	clone := *webhooks
	clone.Service = webhooks.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (webhooks *WebhooksV1) SetServiceURL(url string) error {
	return webhooks.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (webhooks *WebhooksV1) GetServiceURL() string {
	return webhooks.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (webhooks *WebhooksV1) SetDefaultHeaders(headers http.Header) {
	webhooks.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (webhooks *WebhooksV1) SetEnableGzipCompression(enableGzip bool) {
	webhooks.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (webhooks *WebhooksV1) GetEnableGzipCompression() bool {
	return webhooks.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (webhooks *WebhooksV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	webhooks.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (webhooks *WebhooksV1) DisableRetries() {
	webhooks.Service.DisableRetries()
}

// ListWebhooks : List alert webhooks
// List configured alert webhooks for the CIS instance.
func (webhooks *WebhooksV1) ListWebhooks(listWebhooksOptions *ListWebhooksOptions) (result *ListAlertWebhooksResp, response *core.DetailedResponse, err error) {
	return webhooks.ListWebhooksWithContext(context.Background(), listWebhooksOptions)
}

// ListWebhooksWithContext is an alternate form of the ListWebhooks method which supports a Context parameter
func (webhooks *WebhooksV1) ListWebhooksWithContext(ctx context.Context, listWebhooksOptions *ListWebhooksOptions) (result *ListAlertWebhooksResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listWebhooksOptions, "listWebhooksOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *webhooks.Crn,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = webhooks.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(webhooks.Service.Options.URL, `/v1/{crn}/alerting/destinations/webhooks`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listWebhooksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("webhooks", "V1", "ListWebhooks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = webhooks.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListAlertWebhooksResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAlertWebhook : Create an alert webhook
// Create a new alert webhook for the CIS instance.
func (webhooks *WebhooksV1) CreateAlertWebhook(createAlertWebhookOptions *CreateAlertWebhookOptions) (result *WebhookSuccessResp, response *core.DetailedResponse, err error) {
	return webhooks.CreateAlertWebhookWithContext(context.Background(), createAlertWebhookOptions)
}

// CreateAlertWebhookWithContext is an alternate form of the CreateAlertWebhook method which supports a Context parameter
func (webhooks *WebhooksV1) CreateAlertWebhookWithContext(ctx context.Context, createAlertWebhookOptions *CreateAlertWebhookOptions) (result *WebhookSuccessResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createAlertWebhookOptions, "createAlertWebhookOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *webhooks.Crn,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = webhooks.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(webhooks.Service.Options.URL, `/v1/{crn}/alerting/destinations/webhooks`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createAlertWebhookOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("webhooks", "V1", "CreateAlertWebhook")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAlertWebhookOptions.Name != nil {
		body["name"] = createAlertWebhookOptions.Name
	}
	if createAlertWebhookOptions.URL != nil {
		body["url"] = createAlertWebhookOptions.URL
	}
	if createAlertWebhookOptions.Secret != nil {
		body["secret"] = createAlertWebhookOptions.Secret
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
	response, err = webhooks.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWebhookSuccessResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWebhook : Get an alert webhook
// Get an alert webhook for the CIS instance.
func (webhooks *WebhooksV1) GetWebhook(getWebhookOptions *GetWebhookOptions) (result *GetAlertWebhookResp, response *core.DetailedResponse, err error) {
	return webhooks.GetWebhookWithContext(context.Background(), getWebhookOptions)
}

// GetWebhookWithContext is an alternate form of the GetWebhook method which supports a Context parameter
func (webhooks *WebhooksV1) GetWebhookWithContext(ctx context.Context, getWebhookOptions *GetWebhookOptions) (result *GetAlertWebhookResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWebhookOptions, "getWebhookOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWebhookOptions, "getWebhookOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *webhooks.Crn,
		"webhook_id": *getWebhookOptions.WebhookID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = webhooks.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(webhooks.Service.Options.URL, `/v1/{crn}/alerting/destinations/webhooks/{webhook_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWebhookOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("webhooks", "V1", "GetWebhook")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = webhooks.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetAlertWebhookResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateAlertWebhook : Update an alert webhook
// Update an existing alert webhook for the CIS instance.
func (webhooks *WebhooksV1) UpdateAlertWebhook(updateAlertWebhookOptions *UpdateAlertWebhookOptions) (result *WebhookSuccessResp, response *core.DetailedResponse, err error) {
	return webhooks.UpdateAlertWebhookWithContext(context.Background(), updateAlertWebhookOptions)
}

// UpdateAlertWebhookWithContext is an alternate form of the UpdateAlertWebhook method which supports a Context parameter
func (webhooks *WebhooksV1) UpdateAlertWebhookWithContext(ctx context.Context, updateAlertWebhookOptions *UpdateAlertWebhookOptions) (result *WebhookSuccessResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAlertWebhookOptions, "updateAlertWebhookOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAlertWebhookOptions, "updateAlertWebhookOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *webhooks.Crn,
		"webhook_id": *updateAlertWebhookOptions.WebhookID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = webhooks.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(webhooks.Service.Options.URL, `/v1/{crn}/alerting/destinations/webhooks/{webhook_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateAlertWebhookOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("webhooks", "V1", "UpdateAlertWebhook")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAlertWebhookOptions.Name != nil {
		body["name"] = updateAlertWebhookOptions.Name
	}
	if updateAlertWebhookOptions.URL != nil {
		body["url"] = updateAlertWebhookOptions.URL
	}
	if updateAlertWebhookOptions.Secret != nil {
		body["secret"] = updateAlertWebhookOptions.Secret
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
	response, err = webhooks.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWebhookSuccessResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteWebhook : Delete an alert webhook
// Delete an alert webhook for the CIS instance.
func (webhooks *WebhooksV1) DeleteWebhook(deleteWebhookOptions *DeleteWebhookOptions) (result *WebhookSuccessResp, response *core.DetailedResponse, err error) {
	return webhooks.DeleteWebhookWithContext(context.Background(), deleteWebhookOptions)
}

// DeleteWebhookWithContext is an alternate form of the DeleteWebhook method which supports a Context parameter
func (webhooks *WebhooksV1) DeleteWebhookWithContext(ctx context.Context, deleteWebhookOptions *DeleteWebhookOptions) (result *WebhookSuccessResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteWebhookOptions, "deleteWebhookOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteWebhookOptions, "deleteWebhookOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *webhooks.Crn,
		"webhook_id": *deleteWebhookOptions.WebhookID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = webhooks.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(webhooks.Service.Options.URL, `/v1/{crn}/alerting/destinations/webhooks/{webhook_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteWebhookOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("webhooks", "V1", "DeleteWebhook")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = webhooks.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWebhookSuccessResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAlertWebhookOptions : The CreateAlertWebhook options.
type CreateAlertWebhookOptions struct {
	// Webhook Name.
	Name *string `json:"name,omitempty"`

	// Webhook url.
	URL *string `json:"url,omitempty"`

	// The optional secret or API key needed to use the webhook.
	Secret *string `json:"secret,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateAlertWebhookOptions : Instantiate CreateAlertWebhookOptions
func (*WebhooksV1) NewCreateAlertWebhookOptions() *CreateAlertWebhookOptions {
	return &CreateAlertWebhookOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateAlertWebhookOptions) SetName(name string) *CreateAlertWebhookOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetURL : Allow user to set URL
func (_options *CreateAlertWebhookOptions) SetURL(url string) *CreateAlertWebhookOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetSecret : Allow user to set Secret
func (_options *CreateAlertWebhookOptions) SetSecret(secret string) *CreateAlertWebhookOptions {
	_options.Secret = core.StringPtr(secret)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAlertWebhookOptions) SetHeaders(param map[string]string) *CreateAlertWebhookOptions {
	options.Headers = param
	return options
}

// DeleteWebhookOptions : The DeleteWebhook options.
type DeleteWebhookOptions struct {
	// Alert webhook identifier.
	WebhookID *string `json:"webhook_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteWebhookOptions : Instantiate DeleteWebhookOptions
func (*WebhooksV1) NewDeleteWebhookOptions(webhookID string) *DeleteWebhookOptions {
	return &DeleteWebhookOptions{
		WebhookID: core.StringPtr(webhookID),
	}
}

// SetWebhookID : Allow user to set WebhookID
func (_options *DeleteWebhookOptions) SetWebhookID(webhookID string) *DeleteWebhookOptions {
	_options.WebhookID = core.StringPtr(webhookID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteWebhookOptions) SetHeaders(param map[string]string) *DeleteWebhookOptions {
	options.Headers = param
	return options
}

// GetAlertWebhookRespResult : Container for response information.
type GetAlertWebhookRespResult struct {
	// Webhook ID.
	ID *string `json:"id" validate:"required"`

	// Webhook Name.
	Name *string `json:"name" validate:"required"`

	// Webhook url.
	URL *string `json:"url" validate:"required"`

	// Webhook type.
	Type *string `json:"type" validate:"required"`

	// When was the webhook created.
	CreatedAt *string `json:"created_at" validate:"required"`

	// When was the webhook last used successfully.
	LastSuccess *string `json:"last_success" validate:"required"`

	// When was the webhook last used and failed.
	LastFailure *string `json:"last_failure" validate:"required"`
}

// UnmarshalGetAlertWebhookRespResult unmarshals an instance of GetAlertWebhookRespResult from the specified map of raw messages.
func UnmarshalGetAlertWebhookRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertWebhookRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_success", &obj.LastSuccess)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_failure", &obj.LastFailure)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetWebhookOptions : The GetWebhook options.
type GetWebhookOptions struct {
	// Alert webhook identifier.
	WebhookID *string `json:"webhook_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWebhookOptions : Instantiate GetWebhookOptions
func (*WebhooksV1) NewGetWebhookOptions(webhookID string) *GetWebhookOptions {
	return &GetWebhookOptions{
		WebhookID: core.StringPtr(webhookID),
	}
}

// SetWebhookID : Allow user to set WebhookID
func (_options *GetWebhookOptions) SetWebhookID(webhookID string) *GetWebhookOptions {
	_options.WebhookID = core.StringPtr(webhookID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWebhookOptions) SetHeaders(param map[string]string) *GetWebhookOptions {
	options.Headers = param
	return options
}

// ListAlertWebhooksRespResultItem : ListAlertWebhooksRespResultItem struct
type ListAlertWebhooksRespResultItem struct {
	// Webhook ID.
	ID *string `json:"id" validate:"required"`

	// Webhook Name.
	Name *string `json:"name" validate:"required"`

	// Webhook url.
	URL *string `json:"url" validate:"required"`

	// Webhook type.
	Type *string `json:"type" validate:"required"`

	// When was the webhook created.
	CreatedAt *string `json:"created_at" validate:"required"`

	// When was the webhook last used successfully.
	LastSuccess *string `json:"last_success" validate:"required"`

	// When was the webhook last used and failed.
	LastFailure *string `json:"last_failure" validate:"required"`
}

// UnmarshalListAlertWebhooksRespResultItem unmarshals an instance of ListAlertWebhooksRespResultItem from the specified map of raw messages.
func UnmarshalListAlertWebhooksRespResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertWebhooksRespResultItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_success", &obj.LastSuccess)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_failure", &obj.LastFailure)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListWebhooksOptions : The ListWebhooks options.
type ListWebhooksOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListWebhooksOptions : Instantiate ListWebhooksOptions
func (*WebhooksV1) NewListWebhooksOptions() *ListWebhooksOptions {
	return &ListWebhooksOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListWebhooksOptions) SetHeaders(param map[string]string) *ListWebhooksOptions {
	options.Headers = param
	return options
}

// UpdateAlertWebhookOptions : The UpdateAlertWebhook options.
type UpdateAlertWebhookOptions struct {
	// Alert webhook identifier.
	WebhookID *string `json:"webhook_id" validate:"required,ne="`

	// Webhook Name.
	Name *string `json:"name,omitempty"`

	// Webhook url.
	URL *string `json:"url,omitempty"`

	// The optional secret or API key needed to use the webhook.
	Secret *string `json:"secret,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateAlertWebhookOptions : Instantiate UpdateAlertWebhookOptions
func (*WebhooksV1) NewUpdateAlertWebhookOptions(webhookID string) *UpdateAlertWebhookOptions {
	return &UpdateAlertWebhookOptions{
		WebhookID: core.StringPtr(webhookID),
	}
}

// SetWebhookID : Allow user to set WebhookID
func (_options *UpdateAlertWebhookOptions) SetWebhookID(webhookID string) *UpdateAlertWebhookOptions {
	_options.WebhookID = core.StringPtr(webhookID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAlertWebhookOptions) SetName(name string) *UpdateAlertWebhookOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetURL : Allow user to set URL
func (_options *UpdateAlertWebhookOptions) SetURL(url string) *UpdateAlertWebhookOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetSecret : Allow user to set Secret
func (_options *UpdateAlertWebhookOptions) SetSecret(secret string) *UpdateAlertWebhookOptions {
	_options.Secret = core.StringPtr(secret)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAlertWebhookOptions) SetHeaders(param map[string]string) *UpdateAlertWebhookOptions {
	options.Headers = param
	return options
}

// WebhookSuccessRespResult : Container for response information.
type WebhookSuccessRespResult struct {
	// Webhook ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalWebhookSuccessRespResult unmarshals an instance of WebhookSuccessRespResult from the specified map of raw messages.
func UnmarshalWebhookSuccessRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WebhookSuccessRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAlertWebhookResp : Get Alert Webhooks Response.
type GetAlertWebhookResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *GetAlertWebhookRespResult `json:"result" validate:"required"`
}

// UnmarshalGetAlertWebhookResp unmarshals an instance of GetAlertWebhookResp from the specified map of raw messages.
func UnmarshalGetAlertWebhookResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAlertWebhookResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalGetAlertWebhookRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAlertWebhooksResp : List Alert Webhooks Response.
type ListAlertWebhooksResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []ListAlertWebhooksRespResultItem `json:"result" validate:"required"`
}

// UnmarshalListAlertWebhooksResp unmarshals an instance of ListAlertWebhooksResp from the specified map of raw messages.
func UnmarshalListAlertWebhooksResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAlertWebhooksResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalListAlertWebhooksRespResultItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WebhookSuccessResp : Alert Webhooks Response.
type WebhookSuccessResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *WebhookSuccessRespResult `json:"result" validate:"required"`
}

// UnmarshalWebhookSuccessResp unmarshals an instance of WebhookSuccessResp from the specified map of raw messages.
func UnmarshalWebhookSuccessResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WebhookSuccessResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalWebhookSuccessRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
