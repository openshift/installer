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
 * IBM OpenAPI SDK Code Generator Version: 3.93.0-c40121e6-20240729-182103
 */

// Package iamidentityv1 : Operations and models for the IamIdentityV1 service
package iamidentityv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// IamIdentityV1 : The IAM Identity Service API allows for the management of Account Settings and Identities (Service
// IDs, ApiKeys).
//
// API Version: 1.0.0
type IamIdentityV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://iam.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "iam_identity"

// IamIdentityV1Options : Service options
type IamIdentityV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewIamIdentityV1UsingExternalConfig : constructs an instance of IamIdentityV1 with passed in options and external configuration.
func NewIamIdentityV1UsingExternalConfig(options *IamIdentityV1Options) (iamIdentity *IamIdentityV1, err error) {
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

	iamIdentity, err = NewIamIdentityV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = iamIdentity.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = iamIdentity.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewIamIdentityV1 : constructs an instance of IamIdentityV1 with passed in options.
func NewIamIdentityV1(options *IamIdentityV1Options) (service *IamIdentityV1, err error) {
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

	service = &IamIdentityV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "iamIdentity" suitable for processing requests.
func (iamIdentity *IamIdentityV1) Clone() *IamIdentityV1 {
	if core.IsNil(iamIdentity) {
		return nil
	}
	clone := *iamIdentity
	clone.Service = iamIdentity.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (iamIdentity *IamIdentityV1) SetServiceURL(url string) error {
	err := iamIdentity.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (iamIdentity *IamIdentityV1) GetServiceURL() string {
	return iamIdentity.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (iamIdentity *IamIdentityV1) SetDefaultHeaders(headers http.Header) {
	iamIdentity.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (iamIdentity *IamIdentityV1) SetEnableGzipCompression(enableGzip bool) {
	iamIdentity.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (iamIdentity *IamIdentityV1) GetEnableGzipCompression() bool {
	return iamIdentity.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (iamIdentity *IamIdentityV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	iamIdentity.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (iamIdentity *IamIdentityV1) DisableRetries() {
	iamIdentity.Service.DisableRetries()
}

// ListAPIKeys : Get API keys for a given service or user IAM ID and account ID
// Returns the list of API key details for a given service or user IAM ID and account ID. Users can manage user API keys
// for themself, or service ID API keys for service IDs they have access to.
func (iamIdentity *IamIdentityV1) ListAPIKeys(listAPIKeysOptions *ListAPIKeysOptions) (result *APIKeyList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListAPIKeysWithContext(context.Background(), listAPIKeysOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAPIKeysWithContext is an alternate form of the ListAPIKeys method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListAPIKeysWithContext(ctx context.Context, listAPIKeysOptions *ListAPIKeysOptions) (result *APIKeyList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAPIKeysOptions, "listAPIKeysOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAPIKeysOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListAPIKeys")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAPIKeysOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listAPIKeysOptions.AccountID))
	}
	if listAPIKeysOptions.IamID != nil {
		builder.AddQuery("iam_id", fmt.Sprint(*listAPIKeysOptions.IamID))
	}
	if listAPIKeysOptions.Pagesize != nil {
		builder.AddQuery("pagesize", fmt.Sprint(*listAPIKeysOptions.Pagesize))
	}
	if listAPIKeysOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listAPIKeysOptions.Pagetoken))
	}
	if listAPIKeysOptions.Scope != nil {
		builder.AddQuery("scope", fmt.Sprint(*listAPIKeysOptions.Scope))
	}
	if listAPIKeysOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listAPIKeysOptions.Type))
	}
	if listAPIKeysOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listAPIKeysOptions.Sort))
	}
	if listAPIKeysOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listAPIKeysOptions.Order))
	}
	if listAPIKeysOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listAPIKeysOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_api_keys", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKeyList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateAPIKey : Create an API key
// Creates an API key for a UserID or service ID. Users can manage user API keys for themself, or service ID API keys
// for service IDs they have access to.
func (iamIdentity *IamIdentityV1) CreateAPIKey(createAPIKeyOptions *CreateAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateAPIKeyWithContext(context.Background(), createAPIKeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateAPIKeyWithContext is an alternate form of the CreateAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateAPIKeyWithContext(ctx context.Context, createAPIKeyOptions *CreateAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAPIKeyOptions, "createAPIKeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createAPIKeyOptions, "createAPIKeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createAPIKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateAPIKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createAPIKeyOptions.EntityLock != nil {
		builder.AddHeader("Entity-Lock", fmt.Sprint(*createAPIKeyOptions.EntityLock))
	}
	if createAPIKeyOptions.EntityDisable != nil {
		builder.AddHeader("Entity-Disable", fmt.Sprint(*createAPIKeyOptions.EntityDisable))
	}

	body := make(map[string]interface{})
	if createAPIKeyOptions.Name != nil {
		body["name"] = createAPIKeyOptions.Name
	}
	if createAPIKeyOptions.IamID != nil {
		body["iam_id"] = createAPIKeyOptions.IamID
	}
	if createAPIKeyOptions.Description != nil {
		body["description"] = createAPIKeyOptions.Description
	}
	if createAPIKeyOptions.AccountID != nil {
		body["account_id"] = createAPIKeyOptions.AccountID
	}
	if createAPIKeyOptions.Apikey != nil {
		body["apikey"] = createAPIKeyOptions.Apikey
	}
	if createAPIKeyOptions.StoreValue != nil {
		body["store_value"] = createAPIKeyOptions.StoreValue
	}
	if createAPIKeyOptions.SupportSessions != nil {
		body["support_sessions"] = createAPIKeyOptions.SupportSessions
	}
	if createAPIKeyOptions.ActionWhenLeaked != nil {
		body["action_when_leaked"] = createAPIKeyOptions.ActionWhenLeaked
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_api_key", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAPIKeysDetails : Get details of an API key by its value
// Returns the details of an API key by its value. Users can manage user API keys for themself, or service ID API keys
// for service IDs they have access to.
func (iamIdentity *IamIdentityV1) GetAPIKeysDetails(getAPIKeysDetailsOptions *GetAPIKeysDetailsOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetAPIKeysDetailsWithContext(context.Background(), getAPIKeysDetailsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAPIKeysDetailsWithContext is an alternate form of the GetAPIKeysDetails method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetAPIKeysDetailsWithContext(ctx context.Context, getAPIKeysDetailsOptions *GetAPIKeysDetailsOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getAPIKeysDetailsOptions, "getAPIKeysDetailsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/details`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAPIKeysDetailsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetAPIKeysDetails")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getAPIKeysDetailsOptions.IamAPIKey != nil {
		builder.AddHeader("IAM-ApiKey", fmt.Sprint(*getAPIKeysDetailsOptions.IamAPIKey))
	}

	if getAPIKeysDetailsOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getAPIKeysDetailsOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_api_keys_details", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAPIKey : Get details of an API key
// Returns the details of an API key. Users can manage user API keys for themself, or service ID API keys for service
// IDs they have access to.
func (iamIdentity *IamIdentityV1) GetAPIKey(getAPIKeyOptions *GetAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetAPIKeyWithContext(context.Background(), getAPIKeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAPIKeyWithContext is an alternate form of the GetAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetAPIKeyWithContext(ctx context.Context, getAPIKeyOptions *GetAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAPIKeyOptions, "getAPIKeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAPIKeyOptions, "getAPIKeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getAPIKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAPIKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetAPIKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getAPIKeyOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getAPIKeyOptions.IncludeHistory))
	}
	if getAPIKeyOptions.IncludeActivity != nil {
		builder.AddQuery("include_activity", fmt.Sprint(*getAPIKeyOptions.IncludeActivity))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_api_key", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateAPIKey : Updates an API key
// Updates properties of an API key. This does NOT affect existing access tokens. Their token content will stay
// unchanged until the access token is refreshed. To update an API key, pass the property to be modified. To delete one
// property's value, pass the property with an empty value "". Users can manage user API keys for themself, or service
// ID API keys for service IDs they have access to.
func (iamIdentity *IamIdentityV1) UpdateAPIKey(updateAPIKeyOptions *UpdateAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateAPIKeyWithContext(context.Background(), updateAPIKeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateAPIKeyWithContext is an alternate form of the UpdateAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateAPIKeyWithContext(ctx context.Context, updateAPIKeyOptions *UpdateAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAPIKeyOptions, "updateAPIKeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateAPIKeyOptions, "updateAPIKeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateAPIKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateAPIKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateAPIKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateAPIKeyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateAPIKeyOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateAPIKeyOptions.Name != nil {
		body["name"] = updateAPIKeyOptions.Name
	}
	if updateAPIKeyOptions.Description != nil {
		body["description"] = updateAPIKeyOptions.Description
	}
	if updateAPIKeyOptions.SupportSessions != nil {
		body["support_sessions"] = updateAPIKeyOptions.SupportSessions
	}
	if updateAPIKeyOptions.ActionWhenLeaked != nil {
		body["action_when_leaked"] = updateAPIKeyOptions.ActionWhenLeaked
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_api_key", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteAPIKey : Deletes an API key
// Deletes an API key. Existing tokens will remain valid until expired. Users can manage user API keys for themself, or
// service ID API keys for service IDs they have access to.
func (iamIdentity *IamIdentityV1) DeleteAPIKey(deleteAPIKeyOptions *DeleteAPIKeyOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteAPIKeyWithContext(context.Background(), deleteAPIKeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAPIKeyWithContext is an alternate form of the DeleteAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteAPIKeyWithContext(ctx context.Context, deleteAPIKeyOptions *DeleteAPIKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAPIKeyOptions, "deleteAPIKeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAPIKeyOptions, "deleteAPIKeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteAPIKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAPIKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteAPIKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_api_key", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// LockAPIKey : Lock the API key
// Locks an API key by ID. Users can manage user API keys for themself, or service ID API keys for service IDs they have
// access to.
func (iamIdentity *IamIdentityV1) LockAPIKey(lockAPIKeyOptions *LockAPIKeyOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.LockAPIKeyWithContext(context.Background(), lockAPIKeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// LockAPIKeyWithContext is an alternate form of the LockAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) LockAPIKeyWithContext(ctx context.Context, lockAPIKeyOptions *LockAPIKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(lockAPIKeyOptions, "lockAPIKeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(lockAPIKeyOptions, "lockAPIKeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *lockAPIKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/{id}/lock`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range lockAPIKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "LockAPIKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "lock_api_key", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UnlockAPIKey : Unlock the API key
// Unlocks an API key by ID. Users can manage user API keys for themself, or service ID API keys for service IDs they
// have access to.
func (iamIdentity *IamIdentityV1) UnlockAPIKey(unlockAPIKeyOptions *UnlockAPIKeyOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.UnlockAPIKeyWithContext(context.Background(), unlockAPIKeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UnlockAPIKeyWithContext is an alternate form of the UnlockAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UnlockAPIKeyWithContext(ctx context.Context, unlockAPIKeyOptions *UnlockAPIKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(unlockAPIKeyOptions, "unlockAPIKeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(unlockAPIKeyOptions, "unlockAPIKeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *unlockAPIKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/{id}/lock`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range unlockAPIKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UnlockAPIKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "unlock_api_key", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// DisableAPIKey : Disable the API key
// Disable an API key. Users can manage user API keys for themself, or service ID API keys for service IDs they have
// access to.
func (iamIdentity *IamIdentityV1) DisableAPIKey(disableAPIKeyOptions *DisableAPIKeyOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DisableAPIKeyWithContext(context.Background(), disableAPIKeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DisableAPIKeyWithContext is an alternate form of the DisableAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DisableAPIKeyWithContext(ctx context.Context, disableAPIKeyOptions *DisableAPIKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(disableAPIKeyOptions, "disableAPIKeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(disableAPIKeyOptions, "disableAPIKeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *disableAPIKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/{id}/disable`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range disableAPIKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DisableAPIKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "disable_api_key", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// EnableAPIKey : Enable the API key
// Enable an API key. Users can manage user API keys for themself, or service ID API keys for service IDs they have
// access to.
func (iamIdentity *IamIdentityV1) EnableAPIKey(enableAPIKeyOptions *EnableAPIKeyOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.EnableAPIKeyWithContext(context.Background(), enableAPIKeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// EnableAPIKeyWithContext is an alternate form of the EnableAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) EnableAPIKeyWithContext(ctx context.Context, enableAPIKeyOptions *EnableAPIKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(enableAPIKeyOptions, "enableAPIKeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(enableAPIKeyOptions, "enableAPIKeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *enableAPIKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/{id}/disable`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range enableAPIKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "EnableAPIKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "enable_api_key", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListServiceIds : List service IDs
// Returns a list of service IDs. Users can manage user API keys for themself, or service ID API keys for service IDs
// they have access to. Note: apikey details are only included in the response when creating a Service ID with an api
// key.
func (iamIdentity *IamIdentityV1) ListServiceIds(listServiceIdsOptions *ListServiceIdsOptions) (result *ServiceIDList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListServiceIdsWithContext(context.Background(), listServiceIdsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListServiceIdsWithContext is an alternate form of the ListServiceIds method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListServiceIdsWithContext(ctx context.Context, listServiceIdsOptions *ListServiceIdsOptions) (result *ServiceIDList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listServiceIdsOptions, "listServiceIdsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listServiceIdsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListServiceIds")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listServiceIdsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listServiceIdsOptions.AccountID))
	}
	if listServiceIdsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listServiceIdsOptions.Name))
	}
	if listServiceIdsOptions.Pagesize != nil {
		builder.AddQuery("pagesize", fmt.Sprint(*listServiceIdsOptions.Pagesize))
	}
	if listServiceIdsOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listServiceIdsOptions.Pagetoken))
	}
	if listServiceIdsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listServiceIdsOptions.Sort))
	}
	if listServiceIdsOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listServiceIdsOptions.Order))
	}
	if listServiceIdsOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listServiceIdsOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_service_ids", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceIDList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateServiceID : Create a service ID
// Creates a service ID for an IBM Cloud account. Users can manage user API keys for themself, or service ID API keys
// for service IDs they have access to.
func (iamIdentity *IamIdentityV1) CreateServiceID(createServiceIDOptions *CreateServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateServiceIDWithContext(context.Background(), createServiceIDOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateServiceIDWithContext is an alternate form of the CreateServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateServiceIDWithContext(ctx context.Context, createServiceIDOptions *CreateServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createServiceIDOptions, "createServiceIDOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createServiceIDOptions, "createServiceIDOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createServiceIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateServiceID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createServiceIDOptions.EntityLock != nil {
		builder.AddHeader("Entity-Lock", fmt.Sprint(*createServiceIDOptions.EntityLock))
	}

	body := make(map[string]interface{})
	if createServiceIDOptions.AccountID != nil {
		body["account_id"] = createServiceIDOptions.AccountID
	}
	if createServiceIDOptions.Name != nil {
		body["name"] = createServiceIDOptions.Name
	}
	if createServiceIDOptions.Description != nil {
		body["description"] = createServiceIDOptions.Description
	}
	if createServiceIDOptions.UniqueInstanceCrns != nil {
		body["unique_instance_crns"] = createServiceIDOptions.UniqueInstanceCrns
	}
	if createServiceIDOptions.Apikey != nil {
		body["apikey"] = createServiceIDOptions.Apikey
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_service_id", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetServiceID : Get details of a service ID
// Returns the details of a service ID. Users can manage user API keys for themself, or service ID API keys for service
// IDs they have access to. Note: apikey details are only included in the response when creating a Service ID with an
// api key.
func (iamIdentity *IamIdentityV1) GetServiceID(getServiceIDOptions *GetServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetServiceIDWithContext(context.Background(), getServiceIDOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetServiceIDWithContext is an alternate form of the GetServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetServiceIDWithContext(ctx context.Context, getServiceIDOptions *GetServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getServiceIDOptions, "getServiceIDOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getServiceIDOptions, "getServiceIDOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getServiceIDOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getServiceIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetServiceID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getServiceIDOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getServiceIDOptions.IncludeHistory))
	}
	if getServiceIDOptions.IncludeActivity != nil {
		builder.AddQuery("include_activity", fmt.Sprint(*getServiceIDOptions.IncludeActivity))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_service_id", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateServiceID : Update service ID
// Updates properties of a service ID. This does NOT affect existing access tokens. Their token content will stay
// unchanged until the access token is refreshed. To update a service ID, pass the property to be modified. To delete
// one property's value, pass the property with an empty value "".Users can manage user API keys for themself, or
// service ID API keys for service IDs they have access to. Note: apikey details are only included in the response when
// creating a Service ID with an apikey.
func (iamIdentity *IamIdentityV1) UpdateServiceID(updateServiceIDOptions *UpdateServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateServiceIDWithContext(context.Background(), updateServiceIDOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateServiceIDWithContext is an alternate form of the UpdateServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateServiceIDWithContext(ctx context.Context, updateServiceIDOptions *UpdateServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateServiceIDOptions, "updateServiceIDOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateServiceIDOptions, "updateServiceIDOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateServiceIDOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateServiceIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateServiceID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateServiceIDOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateServiceIDOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateServiceIDOptions.Name != nil {
		body["name"] = updateServiceIDOptions.Name
	}
	if updateServiceIDOptions.Description != nil {
		body["description"] = updateServiceIDOptions.Description
	}
	if updateServiceIDOptions.UniqueInstanceCrns != nil {
		body["unique_instance_crns"] = updateServiceIDOptions.UniqueInstanceCrns
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_service_id", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteServiceID : Deletes a service ID and associated API keys
// Deletes a service ID and all API keys associated to it. Before deleting the service ID, all associated API keys are
// deleted. In case a Delete Conflict (status code 409) a retry of the request may help as the service ID is only
// deleted if the associated API keys were successfully deleted before. Users can manage user API keys for themself, or
// service ID API keys for service IDs they have access to.
func (iamIdentity *IamIdentityV1) DeleteServiceID(deleteServiceIDOptions *DeleteServiceIDOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteServiceIDWithContext(context.Background(), deleteServiceIDOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteServiceIDWithContext is an alternate form of the DeleteServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteServiceIDWithContext(ctx context.Context, deleteServiceIDOptions *DeleteServiceIDOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteServiceIDOptions, "deleteServiceIDOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteServiceIDOptions, "deleteServiceIDOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteServiceIDOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteServiceIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteServiceID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_service_id", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// LockServiceID : Lock the service ID
// Locks a service ID by ID. Users can manage user API keys for themself, or service ID API keys for service IDs they
// have access to.
func (iamIdentity *IamIdentityV1) LockServiceID(lockServiceIDOptions *LockServiceIDOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.LockServiceIDWithContext(context.Background(), lockServiceIDOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// LockServiceIDWithContext is an alternate form of the LockServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) LockServiceIDWithContext(ctx context.Context, lockServiceIDOptions *LockServiceIDOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(lockServiceIDOptions, "lockServiceIDOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(lockServiceIDOptions, "lockServiceIDOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *lockServiceIDOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/{id}/lock`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range lockServiceIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "LockServiceID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "lock_service_id", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UnlockServiceID : Unlock the service ID
// Unlocks a service ID by ID. Users can manage user API keys for themself, or service ID API keys for service IDs they
// have access to.
func (iamIdentity *IamIdentityV1) UnlockServiceID(unlockServiceIDOptions *UnlockServiceIDOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.UnlockServiceIDWithContext(context.Background(), unlockServiceIDOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UnlockServiceIDWithContext is an alternate form of the UnlockServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UnlockServiceIDWithContext(ctx context.Context, unlockServiceIDOptions *UnlockServiceIDOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(unlockServiceIDOptions, "unlockServiceIDOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(unlockServiceIDOptions, "unlockServiceIDOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *unlockServiceIDOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/{id}/lock`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range unlockServiceIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UnlockServiceID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "unlock_service_id", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateProfile : Create a trusted profile
// Create a trusted profile for a given account ID.
func (iamIdentity *IamIdentityV1) CreateProfile(createProfileOptions *CreateProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateProfileWithContext(context.Background(), createProfileOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateProfileWithContext is an alternate form of the CreateProfile method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateProfileWithContext(ctx context.Context, createProfileOptions *CreateProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProfileOptions, "createProfileOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createProfileOptions, "createProfileOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createProfileOptions.Name != nil {
		body["name"] = createProfileOptions.Name
	}
	if createProfileOptions.AccountID != nil {
		body["account_id"] = createProfileOptions.AccountID
	}
	if createProfileOptions.Description != nil {
		body["description"] = createProfileOptions.Description
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_profile", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfile)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListProfiles : List trusted profiles
// List the trusted profiles in an account. The `account_id` query parameter determines the account from which to
// retrieve the list of trusted profiles.
func (iamIdentity *IamIdentityV1) ListProfiles(listProfilesOptions *ListProfilesOptions) (result *TrustedProfilesList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListProfilesWithContext(context.Background(), listProfilesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListProfilesWithContext is an alternate form of the ListProfiles method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListProfilesWithContext(ctx context.Context, listProfilesOptions *ListProfilesOptions) (result *TrustedProfilesList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProfilesOptions, "listProfilesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listProfilesOptions, "listProfilesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listProfilesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListProfiles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("account_id", fmt.Sprint(*listProfilesOptions.AccountID))
	if listProfilesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listProfilesOptions.Name))
	}
	if listProfilesOptions.Pagesize != nil {
		builder.AddQuery("pagesize", fmt.Sprint(*listProfilesOptions.Pagesize))
	}
	if listProfilesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listProfilesOptions.Sort))
	}
	if listProfilesOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listProfilesOptions.Order))
	}
	if listProfilesOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listProfilesOptions.IncludeHistory))
	}
	if listProfilesOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listProfilesOptions.Pagetoken))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_profiles", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfilesList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProfile : Get a trusted profile
// Retrieve a trusted profile by its `profile-id`. Only the trusted profile's data is returned (`name`, `description`,
// `iam_id`, etc.), not the federated users or compute resources that qualify to apply the trusted profile.
func (iamIdentity *IamIdentityV1) GetProfile(getProfileOptions *GetProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetProfileWithContext(context.Background(), getProfileOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProfileWithContext is an alternate form of the GetProfile method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetProfileWithContext(ctx context.Context, getProfileOptions *GetProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileOptions, "getProfileOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProfileOptions, "getProfileOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *getProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getProfileOptions.IncludeActivity != nil {
		builder.AddQuery("include_activity", fmt.Sprint(*getProfileOptions.IncludeActivity))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_profile", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfile)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateProfile : Update a trusted profile
// Update the name or description of an existing trusted profile.
func (iamIdentity *IamIdentityV1) UpdateProfile(updateProfileOptions *UpdateProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateProfileWithContext(context.Background(), updateProfileOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateProfileWithContext is an alternate form of the UpdateProfile method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateProfileWithContext(ctx context.Context, updateProfileOptions *UpdateProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProfileOptions, "updateProfileOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateProfileOptions, "updateProfileOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *updateProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateProfileOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateProfileOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateProfileOptions.Name != nil {
		body["name"] = updateProfileOptions.Name
	}
	if updateProfileOptions.Description != nil {
		body["description"] = updateProfileOptions.Description
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_profile", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfile)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteProfile : Delete a trusted profile
// Delete a trusted profile. When you delete trusted profile, compute resources and federated users are unlinked from
// the profile and can no longer apply the trusted profile identity.
func (iamIdentity *IamIdentityV1) DeleteProfile(deleteProfileOptions *DeleteProfileOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteProfileWithContext(context.Background(), deleteProfileOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteProfileWithContext is an alternate form of the DeleteProfile method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteProfileWithContext(ctx context.Context, deleteProfileOptions *DeleteProfileOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProfileOptions, "deleteProfileOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteProfileOptions, "deleteProfileOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *deleteProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_profile", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateClaimRule : Create claim rule for a trusted profile
// Create a claim rule for a trusted profile. There is a limit of 20 rules per trusted profile.
func (iamIdentity *IamIdentityV1) CreateClaimRule(createClaimRuleOptions *CreateClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateClaimRuleWithContext(context.Background(), createClaimRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateClaimRuleWithContext is an alternate form of the CreateClaimRule method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateClaimRuleWithContext(ctx context.Context, createClaimRuleOptions *CreateClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createClaimRuleOptions, "createClaimRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createClaimRuleOptions, "createClaimRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *createClaimRuleOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/rules`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createClaimRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateClaimRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createClaimRuleOptions.Type != nil {
		body["type"] = createClaimRuleOptions.Type
	}
	if createClaimRuleOptions.Conditions != nil {
		body["conditions"] = createClaimRuleOptions.Conditions
	}
	if createClaimRuleOptions.Context != nil {
		body["context"] = createClaimRuleOptions.Context
	}
	if createClaimRuleOptions.Name != nil {
		body["name"] = createClaimRuleOptions.Name
	}
	if createClaimRuleOptions.RealmName != nil {
		body["realm_name"] = createClaimRuleOptions.RealmName
	}
	if createClaimRuleOptions.CrType != nil {
		body["cr_type"] = createClaimRuleOptions.CrType
	}
	if createClaimRuleOptions.Expiration != nil {
		body["expiration"] = createClaimRuleOptions.Expiration
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_claim_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileClaimRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListClaimRules : List claim rules for a trusted profile
// Get a list of all claim rules for a trusted profile. The `profile-id` query parameter determines the profile from
// which to retrieve the list of claim rules.
func (iamIdentity *IamIdentityV1) ListClaimRules(listClaimRulesOptions *ListClaimRulesOptions) (result *ProfileClaimRuleList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListClaimRulesWithContext(context.Background(), listClaimRulesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListClaimRulesWithContext is an alternate form of the ListClaimRules method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListClaimRulesWithContext(ctx context.Context, listClaimRulesOptions *ListClaimRulesOptions) (result *ProfileClaimRuleList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listClaimRulesOptions, "listClaimRulesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listClaimRulesOptions, "listClaimRulesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *listClaimRulesOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/rules`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listClaimRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListClaimRules")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_claim_rules", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileClaimRuleList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetClaimRule : Get a claim rule for a trusted profile
// A specific claim rule can be fetched for a given trusted profile ID and rule ID.
func (iamIdentity *IamIdentityV1) GetClaimRule(getClaimRuleOptions *GetClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetClaimRuleWithContext(context.Background(), getClaimRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetClaimRuleWithContext is an alternate form of the GetClaimRule method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetClaimRuleWithContext(ctx context.Context, getClaimRuleOptions *GetClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getClaimRuleOptions, "getClaimRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getClaimRuleOptions, "getClaimRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *getClaimRuleOptions.ProfileID,
		"rule-id":    *getClaimRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/rules/{rule-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getClaimRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetClaimRule")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_claim_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileClaimRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateClaimRule : Update claim rule for a trusted profile
// Update a specific claim rule for a given trusted profile ID and rule ID.
func (iamIdentity *IamIdentityV1) UpdateClaimRule(updateClaimRuleOptions *UpdateClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateClaimRuleWithContext(context.Background(), updateClaimRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateClaimRuleWithContext is an alternate form of the UpdateClaimRule method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateClaimRuleWithContext(ctx context.Context, updateClaimRuleOptions *UpdateClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateClaimRuleOptions, "updateClaimRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateClaimRuleOptions, "updateClaimRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *updateClaimRuleOptions.ProfileID,
		"rule-id":    *updateClaimRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/rules/{rule-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateClaimRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateClaimRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateClaimRuleOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateClaimRuleOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateClaimRuleOptions.Type != nil {
		body["type"] = updateClaimRuleOptions.Type
	}
	if updateClaimRuleOptions.Conditions != nil {
		body["conditions"] = updateClaimRuleOptions.Conditions
	}
	if updateClaimRuleOptions.Context != nil {
		body["context"] = updateClaimRuleOptions.Context
	}
	if updateClaimRuleOptions.Name != nil {
		body["name"] = updateClaimRuleOptions.Name
	}
	if updateClaimRuleOptions.RealmName != nil {
		body["realm_name"] = updateClaimRuleOptions.RealmName
	}
	if updateClaimRuleOptions.CrType != nil {
		body["cr_type"] = updateClaimRuleOptions.CrType
	}
	if updateClaimRuleOptions.Expiration != nil {
		body["expiration"] = updateClaimRuleOptions.Expiration
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_claim_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileClaimRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteClaimRule : Delete a claim rule
// Delete a claim rule. When you delete a claim rule, federated user or compute resources are no longer required to meet
// the conditions of the claim rule in order to apply the trusted profile.
func (iamIdentity *IamIdentityV1) DeleteClaimRule(deleteClaimRuleOptions *DeleteClaimRuleOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteClaimRuleWithContext(context.Background(), deleteClaimRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteClaimRuleWithContext is an alternate form of the DeleteClaimRule method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteClaimRuleWithContext(ctx context.Context, deleteClaimRuleOptions *DeleteClaimRuleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteClaimRuleOptions, "deleteClaimRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteClaimRuleOptions, "deleteClaimRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *deleteClaimRuleOptions.ProfileID,
		"rule-id":    *deleteClaimRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/rules/{rule-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteClaimRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteClaimRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_claim_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateLink : Create link to a trusted profile
// Create a direct link between a specific compute resource and a trusted profile, rather than creating conditions that
// a compute resource must fulfill to apply a trusted profile.
func (iamIdentity *IamIdentityV1) CreateLink(createLinkOptions *CreateLinkOptions) (result *ProfileLink, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateLinkWithContext(context.Background(), createLinkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateLinkWithContext is an alternate form of the CreateLink method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateLinkWithContext(ctx context.Context, createLinkOptions *CreateLinkOptions) (result *ProfileLink, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createLinkOptions, "createLinkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createLinkOptions, "createLinkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *createLinkOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/links`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createLinkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateLink")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createLinkOptions.CrType != nil {
		body["cr_type"] = createLinkOptions.CrType
	}
	if createLinkOptions.Link != nil {
		body["link"] = createLinkOptions.Link
	}
	if createLinkOptions.Name != nil {
		body["name"] = createLinkOptions.Name
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_link", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileLink)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListLinks : List links to a trusted profile
// Get a list of links to a trusted profile.
func (iamIdentity *IamIdentityV1) ListLinks(listLinksOptions *ListLinksOptions) (result *ProfileLinkList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListLinksWithContext(context.Background(), listLinksOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListLinksWithContext is an alternate form of the ListLinks method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListLinksWithContext(ctx context.Context, listLinksOptions *ListLinksOptions) (result *ProfileLinkList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listLinksOptions, "listLinksOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listLinksOptions, "listLinksOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *listLinksOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/links`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listLinksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListLinks")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_links", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileLinkList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetLink : Get link to a trusted profile
// Get a specific link to a trusted profile by `link_id`.
func (iamIdentity *IamIdentityV1) GetLink(getLinkOptions *GetLinkOptions) (result *ProfileLink, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetLinkWithContext(context.Background(), getLinkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLinkWithContext is an alternate form of the GetLink method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetLinkWithContext(ctx context.Context, getLinkOptions *GetLinkOptions) (result *ProfileLink, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLinkOptions, "getLinkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getLinkOptions, "getLinkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *getLinkOptions.ProfileID,
		"link-id":    *getLinkOptions.LinkID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/links/{link-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLinkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetLink")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_link", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileLink)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteLink : Delete link to a trusted profile
// Delete a link between a compute resource and a trusted profile.
func (iamIdentity *IamIdentityV1) DeleteLink(deleteLinkOptions *DeleteLinkOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteLinkWithContext(context.Background(), deleteLinkOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteLinkWithContext is an alternate form of the DeleteLink method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteLinkWithContext(ctx context.Context, deleteLinkOptions *DeleteLinkOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLinkOptions, "deleteLinkOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteLinkOptions, "deleteLinkOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *deleteLinkOptions.ProfileID,
		"link-id":    *deleteLinkOptions.LinkID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/links/{link-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteLinkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteLink")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_link", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetProfileIdentities : Get a list of identities that can assume the trusted profile
// Get a list of identities that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) GetProfileIdentities(getProfileIdentitiesOptions *GetProfileIdentitiesOptions) (result *ProfileIdentitiesResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetProfileIdentitiesWithContext(context.Background(), getProfileIdentitiesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProfileIdentitiesWithContext is an alternate form of the GetProfileIdentities method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetProfileIdentitiesWithContext(ctx context.Context, getProfileIdentitiesOptions *GetProfileIdentitiesOptions) (result *ProfileIdentitiesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileIdentitiesOptions, "getProfileIdentitiesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProfileIdentitiesOptions, "getProfileIdentitiesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *getProfileIdentitiesOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/identities`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getProfileIdentitiesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetProfileIdentities")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_profile_identities", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileIdentitiesResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// SetProfileIdentities : Update the list of identities that can assume the trusted profile
// Update the list of identities that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) SetProfileIdentities(setProfileIdentitiesOptions *SetProfileIdentitiesOptions) (result *ProfileIdentitiesResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.SetProfileIdentitiesWithContext(context.Background(), setProfileIdentitiesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SetProfileIdentitiesWithContext is an alternate form of the SetProfileIdentities method which supports a Context parameter
func (iamIdentity *IamIdentityV1) SetProfileIdentitiesWithContext(ctx context.Context, setProfileIdentitiesOptions *SetProfileIdentitiesOptions) (result *ProfileIdentitiesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setProfileIdentitiesOptions, "setProfileIdentitiesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(setProfileIdentitiesOptions, "setProfileIdentitiesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id": *setProfileIdentitiesOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/identities`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range setProfileIdentitiesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "SetProfileIdentities")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if setProfileIdentitiesOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*setProfileIdentitiesOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if setProfileIdentitiesOptions.Identities != nil {
		body["identities"] = setProfileIdentitiesOptions.Identities
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "set_profile_identities", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileIdentitiesResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// SetProfileIdentity : Add a specific identity that can assume the trusted profile
// Add a specific identity that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) SetProfileIdentity(setProfileIdentityOptions *SetProfileIdentityOptions) (result *ProfileIdentityResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.SetProfileIdentityWithContext(context.Background(), setProfileIdentityOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SetProfileIdentityWithContext is an alternate form of the SetProfileIdentity method which supports a Context parameter
func (iamIdentity *IamIdentityV1) SetProfileIdentityWithContext(ctx context.Context, setProfileIdentityOptions *SetProfileIdentityOptions) (result *ProfileIdentityResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setProfileIdentityOptions, "setProfileIdentityOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(setProfileIdentityOptions, "setProfileIdentityOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id":    *setProfileIdentityOptions.ProfileID,
		"identity-type": *setProfileIdentityOptions.IdentityType,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/identities/{identity-type}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range setProfileIdentityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "SetProfileIdentity")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setProfileIdentityOptions.Identifier != nil {
		body["identifier"] = setProfileIdentityOptions.Identifier
	}
	if setProfileIdentityOptions.Type != nil {
		body["type"] = setProfileIdentityOptions.Type
	}
	if setProfileIdentityOptions.Accounts != nil {
		body["accounts"] = setProfileIdentityOptions.Accounts
	}
	if setProfileIdentityOptions.Description != nil {
		body["description"] = setProfileIdentityOptions.Description
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "set_profile_identity", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileIdentityResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProfileIdentity : Get the identity that can assume the trusted profile
// Get the identity that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) GetProfileIdentity(getProfileIdentityOptions *GetProfileIdentityOptions) (result *ProfileIdentityResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetProfileIdentityWithContext(context.Background(), getProfileIdentityOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProfileIdentityWithContext is an alternate form of the GetProfileIdentity method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetProfileIdentityWithContext(ctx context.Context, getProfileIdentityOptions *GetProfileIdentityOptions) (result *ProfileIdentityResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileIdentityOptions, "getProfileIdentityOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProfileIdentityOptions, "getProfileIdentityOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id":    *getProfileIdentityOptions.ProfileID,
		"identity-type": *getProfileIdentityOptions.IdentityType,
		"identifier-id": *getProfileIdentityOptions.IdentifierID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/identities/{identity-type}/{identifier-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getProfileIdentityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetProfileIdentity")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_profile_identity", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileIdentityResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteProfileIdentity : Delete the identity that can assume the trusted profile
// Delete the identity that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) DeleteProfileIdentity(deleteProfileIdentityOptions *DeleteProfileIdentityOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteProfileIdentityWithContext(context.Background(), deleteProfileIdentityOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteProfileIdentityWithContext is an alternate form of the DeleteProfileIdentity method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteProfileIdentityWithContext(ctx context.Context, deleteProfileIdentityOptions *DeleteProfileIdentityOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProfileIdentityOptions, "deleteProfileIdentityOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteProfileIdentityOptions, "deleteProfileIdentityOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"profile-id":    *deleteProfileIdentityOptions.ProfileID,
		"identity-type": *deleteProfileIdentityOptions.IdentityType,
		"identifier-id": *deleteProfileIdentityOptions.IdentifierID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles/{profile-id}/identities/{identity-type}/{identifier-id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteProfileIdentityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteProfileIdentity")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_profile_identity", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetAccountSettings : Get account configurations
// Returns the details of an account's configuration.
func (iamIdentity *IamIdentityV1) GetAccountSettings(getAccountSettingsOptions *GetAccountSettingsOptions) (result *AccountSettingsResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetAccountSettingsWithContext(context.Background(), getAccountSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAccountSettingsWithContext is an alternate form of the GetAccountSettings method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetAccountSettingsWithContext(ctx context.Context, getAccountSettingsOptions *GetAccountSettingsOptions) (result *AccountSettingsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountSettingsOptions, "getAccountSettingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAccountSettingsOptions, "getAccountSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getAccountSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/accounts/{account_id}/settings/identity`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAccountSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetAccountSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getAccountSettingsOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getAccountSettingsOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "getAccountSettings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccountSettings : Update account configurations
// Allows a user to configure settings on their account with regards to MFA, MFA excemption list, session lifetimes,
// access control for creating new identities, and enforcing IP restrictions on token creation.
func (iamIdentity *IamIdentityV1) UpdateAccountSettings(updateAccountSettingsOptions *UpdateAccountSettingsOptions) (result *AccountSettingsResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateAccountSettingsWithContext(context.Background(), updateAccountSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateAccountSettingsWithContext is an alternate form of the UpdateAccountSettings method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateAccountSettingsWithContext(ctx context.Context, updateAccountSettingsOptions *UpdateAccountSettingsOptions) (result *AccountSettingsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccountSettingsOptions, "updateAccountSettingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateAccountSettingsOptions, "updateAccountSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *updateAccountSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/accounts/{account_id}/settings/identity`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateAccountSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateAccountSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateAccountSettingsOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateAccountSettingsOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateAccountSettingsOptions.RestrictCreateServiceID != nil {
		body["restrict_create_service_id"] = updateAccountSettingsOptions.RestrictCreateServiceID
	}
	if updateAccountSettingsOptions.RestrictCreatePlatformApikey != nil {
		body["restrict_create_platform_apikey"] = updateAccountSettingsOptions.RestrictCreatePlatformApikey
	}
	if updateAccountSettingsOptions.AllowedIPAddresses != nil {
		body["allowed_ip_addresses"] = updateAccountSettingsOptions.AllowedIPAddresses
	}
	if updateAccountSettingsOptions.Mfa != nil {
		body["mfa"] = updateAccountSettingsOptions.Mfa
	}
	if updateAccountSettingsOptions.UserMfa != nil {
		body["user_mfa"] = updateAccountSettingsOptions.UserMfa
	}
	if updateAccountSettingsOptions.SessionExpirationInSeconds != nil {
		body["session_expiration_in_seconds"] = updateAccountSettingsOptions.SessionExpirationInSeconds
	}
	if updateAccountSettingsOptions.SessionInvalidationInSeconds != nil {
		body["session_invalidation_in_seconds"] = updateAccountSettingsOptions.SessionInvalidationInSeconds
	}
	if updateAccountSettingsOptions.MaxSessionsPerIdentity != nil {
		body["max_sessions_per_identity"] = updateAccountSettingsOptions.MaxSessionsPerIdentity
	}
	if updateAccountSettingsOptions.SystemAccessTokenExpirationInSeconds != nil {
		body["system_access_token_expiration_in_seconds"] = updateAccountSettingsOptions.SystemAccessTokenExpirationInSeconds
	}
	if updateAccountSettingsOptions.SystemRefreshTokenExpirationInSeconds != nil {
		body["system_refresh_token_expiration_in_seconds"] = updateAccountSettingsOptions.SystemRefreshTokenExpirationInSeconds
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "updateAccountSettings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetMfaStatus : Get MFA enrollment status for a single user in the account
// Get MFA enrollment status for a single user in the account.
func (iamIdentity *IamIdentityV1) GetMfaStatus(getMfaStatusOptions *GetMfaStatusOptions) (result *UserMfaEnrollments, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetMfaStatusWithContext(context.Background(), getMfaStatusOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetMfaStatusWithContext is an alternate form of the GetMfaStatus method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetMfaStatusWithContext(ctx context.Context, getMfaStatusOptions *GetMfaStatusOptions) (result *UserMfaEnrollments, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getMfaStatusOptions, "getMfaStatusOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getMfaStatusOptions, "getMfaStatusOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getMfaStatusOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/mfa/accounts/{account_id}/status`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getMfaStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetMfaStatus")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("iam_id", fmt.Sprint(*getMfaStatusOptions.IamID))

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_mfa_status", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUserMfaEnrollments)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateMfaReport : Trigger MFA enrollment status report for the account
// Trigger MFA enrollment status report for the account by specifying the account ID. It can take a few minutes to
// generate the report for retrieval.
func (iamIdentity *IamIdentityV1) CreateMfaReport(createMfaReportOptions *CreateMfaReportOptions) (result *ReportReference, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateMfaReportWithContext(context.Background(), createMfaReportOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateMfaReportWithContext is an alternate form of the CreateMfaReport method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateMfaReportWithContext(ctx context.Context, createMfaReportOptions *CreateMfaReportOptions) (result *ReportReference, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createMfaReportOptions, "createMfaReportOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createMfaReportOptions, "createMfaReportOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *createMfaReportOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/mfa/accounts/{account_id}/report`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createMfaReportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateMfaReport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if createMfaReportOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*createMfaReportOptions.Type))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_mfa_report", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportReference)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetMfaReport : Get MFA enrollment status report for the account
// Get MFA enrollment status report for the account by specifying the account ID and the reference that is generated by
// triggering the report. Reports older than a day are deleted when generating a new report.
func (iamIdentity *IamIdentityV1) GetMfaReport(getMfaReportOptions *GetMfaReportOptions) (result *ReportMfaEnrollmentStatus, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetMfaReportWithContext(context.Background(), getMfaReportOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetMfaReportWithContext is an alternate form of the GetMfaReport method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetMfaReportWithContext(ctx context.Context, getMfaReportOptions *GetMfaReportOptions) (result *ReportMfaEnrollmentStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getMfaReportOptions, "getMfaReportOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getMfaReportOptions, "getMfaReportOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getMfaReportOptions.AccountID,
		"reference":  *getMfaReportOptions.Reference,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/mfa/accounts/{account_id}/report/{reference}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getMfaReportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetMfaReport")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_mfa_report", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportMfaEnrollmentStatus)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListAccountSettingsAssignments : List assignments
// List account settings assignments.
func (iamIdentity *IamIdentityV1) ListAccountSettingsAssignments(listAccountSettingsAssignmentsOptions *ListAccountSettingsAssignmentsOptions) (result *TemplateAssignmentListResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListAccountSettingsAssignmentsWithContext(context.Background(), listAccountSettingsAssignmentsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAccountSettingsAssignmentsWithContext is an alternate form of the ListAccountSettingsAssignments method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListAccountSettingsAssignmentsWithContext(ctx context.Context, listAccountSettingsAssignmentsOptions *ListAccountSettingsAssignmentsOptions) (result *TemplateAssignmentListResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAccountSettingsAssignmentsOptions, "listAccountSettingsAssignmentsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_assignments/`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAccountSettingsAssignmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListAccountSettingsAssignments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAccountSettingsAssignmentsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listAccountSettingsAssignmentsOptions.AccountID))
	}
	if listAccountSettingsAssignmentsOptions.TemplateID != nil {
		builder.AddQuery("template_id", fmt.Sprint(*listAccountSettingsAssignmentsOptions.TemplateID))
	}
	if listAccountSettingsAssignmentsOptions.TemplateVersion != nil {
		builder.AddQuery("template_version", fmt.Sprint(*listAccountSettingsAssignmentsOptions.TemplateVersion))
	}
	if listAccountSettingsAssignmentsOptions.Target != nil {
		builder.AddQuery("target", fmt.Sprint(*listAccountSettingsAssignmentsOptions.Target))
	}
	if listAccountSettingsAssignmentsOptions.TargetType != nil {
		builder.AddQuery("target_type", fmt.Sprint(*listAccountSettingsAssignmentsOptions.TargetType))
	}
	if listAccountSettingsAssignmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAccountSettingsAssignmentsOptions.Limit))
	}
	if listAccountSettingsAssignmentsOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listAccountSettingsAssignmentsOptions.Pagetoken))
	}
	if listAccountSettingsAssignmentsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listAccountSettingsAssignmentsOptions.Sort))
	}
	if listAccountSettingsAssignmentsOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listAccountSettingsAssignmentsOptions.Order))
	}
	if listAccountSettingsAssignmentsOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listAccountSettingsAssignmentsOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_account_settings_assignments", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAssignmentListResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateAccountSettingsAssignment : Create assignment
// Create an assigment for an account settings template.
func (iamIdentity *IamIdentityV1) CreateAccountSettingsAssignment(createAccountSettingsAssignmentOptions *CreateAccountSettingsAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateAccountSettingsAssignmentWithContext(context.Background(), createAccountSettingsAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateAccountSettingsAssignmentWithContext is an alternate form of the CreateAccountSettingsAssignment method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateAccountSettingsAssignmentWithContext(ctx context.Context, createAccountSettingsAssignmentOptions *CreateAccountSettingsAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAccountSettingsAssignmentOptions, "createAccountSettingsAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createAccountSettingsAssignmentOptions, "createAccountSettingsAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_assignments/`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createAccountSettingsAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateAccountSettingsAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccountSettingsAssignmentOptions.TemplateID != nil {
		body["template_id"] = createAccountSettingsAssignmentOptions.TemplateID
	}
	if createAccountSettingsAssignmentOptions.TemplateVersion != nil {
		body["template_version"] = createAccountSettingsAssignmentOptions.TemplateVersion
	}
	if createAccountSettingsAssignmentOptions.TargetType != nil {
		body["target_type"] = createAccountSettingsAssignmentOptions.TargetType
	}
	if createAccountSettingsAssignmentOptions.Target != nil {
		body["target"] = createAccountSettingsAssignmentOptions.Target
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_account_settings_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAssignmentResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAccountSettingsAssignment : Get assignment
// Get an assigment for an account settings template.
func (iamIdentity *IamIdentityV1) GetAccountSettingsAssignment(getAccountSettingsAssignmentOptions *GetAccountSettingsAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetAccountSettingsAssignmentWithContext(context.Background(), getAccountSettingsAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAccountSettingsAssignmentWithContext is an alternate form of the GetAccountSettingsAssignment method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetAccountSettingsAssignmentWithContext(ctx context.Context, getAccountSettingsAssignmentOptions *GetAccountSettingsAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountSettingsAssignmentOptions, "getAccountSettingsAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAccountSettingsAssignmentOptions, "getAccountSettingsAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *getAccountSettingsAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAccountSettingsAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetAccountSettingsAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getAccountSettingsAssignmentOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getAccountSettingsAssignmentOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_account_settings_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAssignmentResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteAccountSettingsAssignment : Delete assignment
// Delete an account settings template assignment. This removes any IAM resources created by this assignment in child
// accounts.
func (iamIdentity *IamIdentityV1) DeleteAccountSettingsAssignment(deleteAccountSettingsAssignmentOptions *DeleteAccountSettingsAssignmentOptions) (result *ExceptionResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.DeleteAccountSettingsAssignmentWithContext(context.Background(), deleteAccountSettingsAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAccountSettingsAssignmentWithContext is an alternate form of the DeleteAccountSettingsAssignment method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteAccountSettingsAssignmentWithContext(ctx context.Context, deleteAccountSettingsAssignmentOptions *DeleteAccountSettingsAssignmentOptions) (result *ExceptionResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAccountSettingsAssignmentOptions, "deleteAccountSettingsAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAccountSettingsAssignmentOptions, "deleteAccountSettingsAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *deleteAccountSettingsAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAccountSettingsAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteAccountSettingsAssignment")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_account_settings_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExceptionResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccountSettingsAssignment : Update assignment
// Update an account settings assignment. Call this method to retry failed assignments or migrate the settings in child
// accounts to a new version.
func (iamIdentity *IamIdentityV1) UpdateAccountSettingsAssignment(updateAccountSettingsAssignmentOptions *UpdateAccountSettingsAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateAccountSettingsAssignmentWithContext(context.Background(), updateAccountSettingsAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateAccountSettingsAssignmentWithContext is an alternate form of the UpdateAccountSettingsAssignment method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateAccountSettingsAssignmentWithContext(ctx context.Context, updateAccountSettingsAssignmentOptions *UpdateAccountSettingsAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccountSettingsAssignmentOptions, "updateAccountSettingsAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateAccountSettingsAssignmentOptions, "updateAccountSettingsAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *updateAccountSettingsAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateAccountSettingsAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateAccountSettingsAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateAccountSettingsAssignmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateAccountSettingsAssignmentOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateAccountSettingsAssignmentOptions.TemplateVersion != nil {
		body["template_version"] = updateAccountSettingsAssignmentOptions.TemplateVersion
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_account_settings_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAssignmentResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListAccountSettingsTemplates : List account settings templates
// List account settings templates in an enterprise account.
func (iamIdentity *IamIdentityV1) ListAccountSettingsTemplates(listAccountSettingsTemplatesOptions *ListAccountSettingsTemplatesOptions) (result *AccountSettingsTemplateList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListAccountSettingsTemplatesWithContext(context.Background(), listAccountSettingsTemplatesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAccountSettingsTemplatesWithContext is an alternate form of the ListAccountSettingsTemplates method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListAccountSettingsTemplatesWithContext(ctx context.Context, listAccountSettingsTemplatesOptions *ListAccountSettingsTemplatesOptions) (result *AccountSettingsTemplateList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAccountSettingsTemplatesOptions, "listAccountSettingsTemplatesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAccountSettingsTemplatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListAccountSettingsTemplates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAccountSettingsTemplatesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listAccountSettingsTemplatesOptions.AccountID))
	}
	if listAccountSettingsTemplatesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAccountSettingsTemplatesOptions.Limit))
	}
	if listAccountSettingsTemplatesOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listAccountSettingsTemplatesOptions.Pagetoken))
	}
	if listAccountSettingsTemplatesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listAccountSettingsTemplatesOptions.Sort))
	}
	if listAccountSettingsTemplatesOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listAccountSettingsTemplatesOptions.Order))
	}
	if listAccountSettingsTemplatesOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listAccountSettingsTemplatesOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_account_settings_templates", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsTemplateList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateAccountSettingsTemplate : Create an account settings template
// Create a new account settings template in an enterprise account.
func (iamIdentity *IamIdentityV1) CreateAccountSettingsTemplate(createAccountSettingsTemplateOptions *CreateAccountSettingsTemplateOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateAccountSettingsTemplateWithContext(context.Background(), createAccountSettingsTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateAccountSettingsTemplateWithContext is an alternate form of the CreateAccountSettingsTemplate method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateAccountSettingsTemplateWithContext(ctx context.Context, createAccountSettingsTemplateOptions *CreateAccountSettingsTemplateOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAccountSettingsTemplateOptions, "createAccountSettingsTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createAccountSettingsTemplateOptions, "createAccountSettingsTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createAccountSettingsTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateAccountSettingsTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccountSettingsTemplateOptions.AccountID != nil {
		body["account_id"] = createAccountSettingsTemplateOptions.AccountID
	}
	if createAccountSettingsTemplateOptions.Name != nil {
		body["name"] = createAccountSettingsTemplateOptions.Name
	}
	if createAccountSettingsTemplateOptions.Description != nil {
		body["description"] = createAccountSettingsTemplateOptions.Description
	}
	if createAccountSettingsTemplateOptions.AccountSettings != nil {
		body["account_settings"] = createAccountSettingsTemplateOptions.AccountSettings
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_account_settings_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetLatestAccountSettingsTemplateVersion : Get latest version of an account settings template
// Get the latest version of a specific account settings template in an enterprise account.
func (iamIdentity *IamIdentityV1) GetLatestAccountSettingsTemplateVersion(getLatestAccountSettingsTemplateVersionOptions *GetLatestAccountSettingsTemplateVersionOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetLatestAccountSettingsTemplateVersionWithContext(context.Background(), getLatestAccountSettingsTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLatestAccountSettingsTemplateVersionWithContext is an alternate form of the GetLatestAccountSettingsTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetLatestAccountSettingsTemplateVersionWithContext(ctx context.Context, getLatestAccountSettingsTemplateVersionOptions *GetLatestAccountSettingsTemplateVersionOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLatestAccountSettingsTemplateVersionOptions, "getLatestAccountSettingsTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getLatestAccountSettingsTemplateVersionOptions, "getLatestAccountSettingsTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *getLatestAccountSettingsTemplateVersionOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates/{template_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLatestAccountSettingsTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetLatestAccountSettingsTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getLatestAccountSettingsTemplateVersionOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getLatestAccountSettingsTemplateVersionOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_latest_account_settings_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteAllVersionsOfAccountSettingsTemplate : Delete all versions of an account settings template
// Delete all versions of an account settings template in an enterprise account. If any version is assigned to child
// accounts, you must first delete the assignment.
func (iamIdentity *IamIdentityV1) DeleteAllVersionsOfAccountSettingsTemplate(deleteAllVersionsOfAccountSettingsTemplateOptions *DeleteAllVersionsOfAccountSettingsTemplateOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteAllVersionsOfAccountSettingsTemplateWithContext(context.Background(), deleteAllVersionsOfAccountSettingsTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAllVersionsOfAccountSettingsTemplateWithContext is an alternate form of the DeleteAllVersionsOfAccountSettingsTemplate method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteAllVersionsOfAccountSettingsTemplateWithContext(ctx context.Context, deleteAllVersionsOfAccountSettingsTemplateOptions *DeleteAllVersionsOfAccountSettingsTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAllVersionsOfAccountSettingsTemplateOptions, "deleteAllVersionsOfAccountSettingsTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAllVersionsOfAccountSettingsTemplateOptions, "deleteAllVersionsOfAccountSettingsTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *deleteAllVersionsOfAccountSettingsTemplateOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates/{template_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAllVersionsOfAccountSettingsTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteAllVersionsOfAccountSettingsTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_all_versions_of_account_settings_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListVersionsOfAccountSettingsTemplate : List account settings template versions
// List the versions of a specific account settings template in an enterprise account.
func (iamIdentity *IamIdentityV1) ListVersionsOfAccountSettingsTemplate(listVersionsOfAccountSettingsTemplateOptions *ListVersionsOfAccountSettingsTemplateOptions) (result *AccountSettingsTemplateList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListVersionsOfAccountSettingsTemplateWithContext(context.Background(), listVersionsOfAccountSettingsTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListVersionsOfAccountSettingsTemplateWithContext is an alternate form of the ListVersionsOfAccountSettingsTemplate method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListVersionsOfAccountSettingsTemplateWithContext(ctx context.Context, listVersionsOfAccountSettingsTemplateOptions *ListVersionsOfAccountSettingsTemplateOptions) (result *AccountSettingsTemplateList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listVersionsOfAccountSettingsTemplateOptions, "listVersionsOfAccountSettingsTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listVersionsOfAccountSettingsTemplateOptions, "listVersionsOfAccountSettingsTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *listVersionsOfAccountSettingsTemplateOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates/{template_id}/versions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listVersionsOfAccountSettingsTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListVersionsOfAccountSettingsTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listVersionsOfAccountSettingsTemplateOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listVersionsOfAccountSettingsTemplateOptions.Limit))
	}
	if listVersionsOfAccountSettingsTemplateOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listVersionsOfAccountSettingsTemplateOptions.Pagetoken))
	}
	if listVersionsOfAccountSettingsTemplateOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listVersionsOfAccountSettingsTemplateOptions.Sort))
	}
	if listVersionsOfAccountSettingsTemplateOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listVersionsOfAccountSettingsTemplateOptions.Order))
	}
	if listVersionsOfAccountSettingsTemplateOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listVersionsOfAccountSettingsTemplateOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_versions_of_account_settings_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsTemplateList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateAccountSettingsTemplateVersion : Create a new version of an account settings template
// Create a new version of an account settings template in an Enterprise Account.
func (iamIdentity *IamIdentityV1) CreateAccountSettingsTemplateVersion(createAccountSettingsTemplateVersionOptions *CreateAccountSettingsTemplateVersionOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateAccountSettingsTemplateVersionWithContext(context.Background(), createAccountSettingsTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateAccountSettingsTemplateVersionWithContext is an alternate form of the CreateAccountSettingsTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateAccountSettingsTemplateVersionWithContext(ctx context.Context, createAccountSettingsTemplateVersionOptions *CreateAccountSettingsTemplateVersionOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAccountSettingsTemplateVersionOptions, "createAccountSettingsTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createAccountSettingsTemplateVersionOptions, "createAccountSettingsTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *createAccountSettingsTemplateVersionOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates/{template_id}/versions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createAccountSettingsTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateAccountSettingsTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccountSettingsTemplateVersionOptions.AccountID != nil {
		body["account_id"] = createAccountSettingsTemplateVersionOptions.AccountID
	}
	if createAccountSettingsTemplateVersionOptions.Name != nil {
		body["name"] = createAccountSettingsTemplateVersionOptions.Name
	}
	if createAccountSettingsTemplateVersionOptions.Description != nil {
		body["description"] = createAccountSettingsTemplateVersionOptions.Description
	}
	if createAccountSettingsTemplateVersionOptions.AccountSettings != nil {
		body["account_settings"] = createAccountSettingsTemplateVersionOptions.AccountSettings
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_account_settings_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAccountSettingsTemplateVersion : Get version of an account settings template
// Get a specific version of an account settings template in an Enterprise Account.
func (iamIdentity *IamIdentityV1) GetAccountSettingsTemplateVersion(getAccountSettingsTemplateVersionOptions *GetAccountSettingsTemplateVersionOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetAccountSettingsTemplateVersionWithContext(context.Background(), getAccountSettingsTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAccountSettingsTemplateVersionWithContext is an alternate form of the GetAccountSettingsTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetAccountSettingsTemplateVersionWithContext(ctx context.Context, getAccountSettingsTemplateVersionOptions *GetAccountSettingsTemplateVersionOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountSettingsTemplateVersionOptions, "getAccountSettingsTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAccountSettingsTemplateVersionOptions, "getAccountSettingsTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *getAccountSettingsTemplateVersionOptions.TemplateID,
		"version":     *getAccountSettingsTemplateVersionOptions.Version,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates/{template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAccountSettingsTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetAccountSettingsTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getAccountSettingsTemplateVersionOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getAccountSettingsTemplateVersionOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_account_settings_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccountSettingsTemplateVersion : Update version of an account settings template
// Update a specific version of an account settings template in an Enterprise Account.
func (iamIdentity *IamIdentityV1) UpdateAccountSettingsTemplateVersion(updateAccountSettingsTemplateVersionOptions *UpdateAccountSettingsTemplateVersionOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateAccountSettingsTemplateVersionWithContext(context.Background(), updateAccountSettingsTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateAccountSettingsTemplateVersionWithContext is an alternate form of the UpdateAccountSettingsTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateAccountSettingsTemplateVersionWithContext(ctx context.Context, updateAccountSettingsTemplateVersionOptions *UpdateAccountSettingsTemplateVersionOptions) (result *AccountSettingsTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccountSettingsTemplateVersionOptions, "updateAccountSettingsTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateAccountSettingsTemplateVersionOptions, "updateAccountSettingsTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *updateAccountSettingsTemplateVersionOptions.TemplateID,
		"version":     *updateAccountSettingsTemplateVersionOptions.Version,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates/{template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateAccountSettingsTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateAccountSettingsTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateAccountSettingsTemplateVersionOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateAccountSettingsTemplateVersionOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateAccountSettingsTemplateVersionOptions.AccountID != nil {
		body["account_id"] = updateAccountSettingsTemplateVersionOptions.AccountID
	}
	if updateAccountSettingsTemplateVersionOptions.Name != nil {
		body["name"] = updateAccountSettingsTemplateVersionOptions.Name
	}
	if updateAccountSettingsTemplateVersionOptions.Description != nil {
		body["description"] = updateAccountSettingsTemplateVersionOptions.Description
	}
	if updateAccountSettingsTemplateVersionOptions.AccountSettings != nil {
		body["account_settings"] = updateAccountSettingsTemplateVersionOptions.AccountSettings
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_account_settings_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteAccountSettingsTemplateVersion : Delete version of an account settings template
// Delete a specific version of an account settings template in an Enterprise Account.
func (iamIdentity *IamIdentityV1) DeleteAccountSettingsTemplateVersion(deleteAccountSettingsTemplateVersionOptions *DeleteAccountSettingsTemplateVersionOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteAccountSettingsTemplateVersionWithContext(context.Background(), deleteAccountSettingsTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAccountSettingsTemplateVersionWithContext is an alternate form of the DeleteAccountSettingsTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteAccountSettingsTemplateVersionWithContext(ctx context.Context, deleteAccountSettingsTemplateVersionOptions *DeleteAccountSettingsTemplateVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAccountSettingsTemplateVersionOptions, "deleteAccountSettingsTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAccountSettingsTemplateVersionOptions, "deleteAccountSettingsTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *deleteAccountSettingsTemplateVersionOptions.TemplateID,
		"version":     *deleteAccountSettingsTemplateVersionOptions.Version,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates/{template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAccountSettingsTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteAccountSettingsTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_account_settings_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CommitAccountSettingsTemplate : Commit a template version
// Commit a specific version of an account settings template in an Enterprise Account. A Template must be committed
// before being assigned, and once committed, can no longer be modified.
func (iamIdentity *IamIdentityV1) CommitAccountSettingsTemplate(commitAccountSettingsTemplateOptions *CommitAccountSettingsTemplateOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.CommitAccountSettingsTemplateWithContext(context.Background(), commitAccountSettingsTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CommitAccountSettingsTemplateWithContext is an alternate form of the CommitAccountSettingsTemplate method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CommitAccountSettingsTemplateWithContext(ctx context.Context, commitAccountSettingsTemplateOptions *CommitAccountSettingsTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(commitAccountSettingsTemplateOptions, "commitAccountSettingsTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(commitAccountSettingsTemplateOptions, "commitAccountSettingsTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *commitAccountSettingsTemplateOptions.TemplateID,
		"version":     *commitAccountSettingsTemplateOptions.Version,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/account_settings_templates/{template_id}/versions/{version}/commit`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range commitAccountSettingsTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CommitAccountSettingsTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "commit_account_settings_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateReport : Trigger activity report for the account
// Trigger activity report for the account by specifying the account ID. It can take a few minutes to generate the
// report for retrieval.
func (iamIdentity *IamIdentityV1) CreateReport(createReportOptions *CreateReportOptions) (result *ReportReference, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateReportWithContext(context.Background(), createReportOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateReportWithContext is an alternate form of the CreateReport method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateReportWithContext(ctx context.Context, createReportOptions *CreateReportOptions) (result *ReportReference, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createReportOptions, "createReportOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createReportOptions, "createReportOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *createReportOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/activity/accounts/{account_id}/report`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createReportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateReport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if createReportOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*createReportOptions.Type))
	}
	if createReportOptions.Duration != nil {
		builder.AddQuery("duration", fmt.Sprint(*createReportOptions.Duration))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_report", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportReference)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetReport : Get activity report for the account
// Get activity report for the account by specifying the account ID and the reference that is generated by triggering
// the report. Reports older than a day are deleted when generating a new report.
func (iamIdentity *IamIdentityV1) GetReport(getReportOptions *GetReportOptions) (result *Report, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetReportWithContext(context.Background(), getReportOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetReportWithContext is an alternate form of the GetReport method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetReportWithContext(ctx context.Context, getReportOptions *GetReportOptions) (result *Report, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportOptions, "getReportOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getReportOptions, "getReportOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getReportOptions.AccountID,
		"reference":  *getReportOptions.Reference,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/activity/accounts/{account_id}/report/{reference}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getReportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetReport")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_report", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReport)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetEffectiveAccountSettings : Get effective account settings configuration
// Returns effective account settings for given account ID.
func (iamIdentity *IamIdentityV1) GetEffectiveAccountSettings(getEffectiveAccountSettingsOptions *GetEffectiveAccountSettingsOptions) (result *EffectiveAccountSettingsResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetEffectiveAccountSettingsWithContext(context.Background(), getEffectiveAccountSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetEffectiveAccountSettingsWithContext is an alternate form of the GetEffectiveAccountSettings method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetEffectiveAccountSettingsWithContext(ctx context.Context, getEffectiveAccountSettingsOptions *GetEffectiveAccountSettingsOptions) (result *EffectiveAccountSettingsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEffectiveAccountSettingsOptions, "getEffectiveAccountSettingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getEffectiveAccountSettingsOptions, "getEffectiveAccountSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getEffectiveAccountSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/accounts/{account_id}/effective_settings/identity`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getEffectiveAccountSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetEffectiveAccountSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getEffectiveAccountSettingsOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getEffectiveAccountSettingsOptions.IncludeHistory))
	}
	if getEffectiveAccountSettingsOptions.ResolveUserMfa != nil {
		builder.AddQuery("resolve_user_mfa", fmt.Sprint(*getEffectiveAccountSettingsOptions.ResolveUserMfa))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_effective_account_settings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEffectiveAccountSettingsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListTrustedProfileAssignments : List assignments
// List trusted profile template assignments.
func (iamIdentity *IamIdentityV1) ListTrustedProfileAssignments(listTrustedProfileAssignmentsOptions *ListTrustedProfileAssignmentsOptions) (result *TemplateAssignmentListResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListTrustedProfileAssignmentsWithContext(context.Background(), listTrustedProfileAssignmentsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTrustedProfileAssignmentsWithContext is an alternate form of the ListTrustedProfileAssignments method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListTrustedProfileAssignmentsWithContext(ctx context.Context, listTrustedProfileAssignmentsOptions *ListTrustedProfileAssignmentsOptions) (result *TemplateAssignmentListResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listTrustedProfileAssignmentsOptions, "listTrustedProfileAssignmentsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_assignments/`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTrustedProfileAssignmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListTrustedProfileAssignments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listTrustedProfileAssignmentsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listTrustedProfileAssignmentsOptions.AccountID))
	}
	if listTrustedProfileAssignmentsOptions.TemplateID != nil {
		builder.AddQuery("template_id", fmt.Sprint(*listTrustedProfileAssignmentsOptions.TemplateID))
	}
	if listTrustedProfileAssignmentsOptions.TemplateVersion != nil {
		builder.AddQuery("template_version", fmt.Sprint(*listTrustedProfileAssignmentsOptions.TemplateVersion))
	}
	if listTrustedProfileAssignmentsOptions.Target != nil {
		builder.AddQuery("target", fmt.Sprint(*listTrustedProfileAssignmentsOptions.Target))
	}
	if listTrustedProfileAssignmentsOptions.TargetType != nil {
		builder.AddQuery("target_type", fmt.Sprint(*listTrustedProfileAssignmentsOptions.TargetType))
	}
	if listTrustedProfileAssignmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listTrustedProfileAssignmentsOptions.Limit))
	}
	if listTrustedProfileAssignmentsOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listTrustedProfileAssignmentsOptions.Pagetoken))
	}
	if listTrustedProfileAssignmentsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listTrustedProfileAssignmentsOptions.Sort))
	}
	if listTrustedProfileAssignmentsOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listTrustedProfileAssignmentsOptions.Order))
	}
	if listTrustedProfileAssignmentsOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listTrustedProfileAssignmentsOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_trusted_profile_assignments", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAssignmentListResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTrustedProfileAssignment : Create assignment
// Create an assigment for a trusted profile template.
func (iamIdentity *IamIdentityV1) CreateTrustedProfileAssignment(createTrustedProfileAssignmentOptions *CreateTrustedProfileAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateTrustedProfileAssignmentWithContext(context.Background(), createTrustedProfileAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTrustedProfileAssignmentWithContext is an alternate form of the CreateTrustedProfileAssignment method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateTrustedProfileAssignmentWithContext(ctx context.Context, createTrustedProfileAssignmentOptions *CreateTrustedProfileAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTrustedProfileAssignmentOptions, "createTrustedProfileAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTrustedProfileAssignmentOptions, "createTrustedProfileAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_assignments/`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTrustedProfileAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateTrustedProfileAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTrustedProfileAssignmentOptions.TemplateID != nil {
		body["template_id"] = createTrustedProfileAssignmentOptions.TemplateID
	}
	if createTrustedProfileAssignmentOptions.TemplateVersion != nil {
		body["template_version"] = createTrustedProfileAssignmentOptions.TemplateVersion
	}
	if createTrustedProfileAssignmentOptions.TargetType != nil {
		body["target_type"] = createTrustedProfileAssignmentOptions.TargetType
	}
	if createTrustedProfileAssignmentOptions.Target != nil {
		body["target"] = createTrustedProfileAssignmentOptions.Target
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_trusted_profile_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAssignmentResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTrustedProfileAssignment : Get assignment
// Get an assigment for a trusted profile template.
func (iamIdentity *IamIdentityV1) GetTrustedProfileAssignment(getTrustedProfileAssignmentOptions *GetTrustedProfileAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetTrustedProfileAssignmentWithContext(context.Background(), getTrustedProfileAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTrustedProfileAssignmentWithContext is an alternate form of the GetTrustedProfileAssignment method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetTrustedProfileAssignmentWithContext(ctx context.Context, getTrustedProfileAssignmentOptions *GetTrustedProfileAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTrustedProfileAssignmentOptions, "getTrustedProfileAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTrustedProfileAssignmentOptions, "getTrustedProfileAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *getTrustedProfileAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTrustedProfileAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetTrustedProfileAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getTrustedProfileAssignmentOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getTrustedProfileAssignmentOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_trusted_profile_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAssignmentResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTrustedProfileAssignment : Delete assignment
// Delete a trusted profile assignment. This removes any IAM resources created by this assignment in child accounts.
func (iamIdentity *IamIdentityV1) DeleteTrustedProfileAssignment(deleteTrustedProfileAssignmentOptions *DeleteTrustedProfileAssignmentOptions) (result *ExceptionResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.DeleteTrustedProfileAssignmentWithContext(context.Background(), deleteTrustedProfileAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTrustedProfileAssignmentWithContext is an alternate form of the DeleteTrustedProfileAssignment method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteTrustedProfileAssignmentWithContext(ctx context.Context, deleteTrustedProfileAssignmentOptions *DeleteTrustedProfileAssignmentOptions) (result *ExceptionResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTrustedProfileAssignmentOptions, "deleteTrustedProfileAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTrustedProfileAssignmentOptions, "deleteTrustedProfileAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *deleteTrustedProfileAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTrustedProfileAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteTrustedProfileAssignment")
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_trusted_profile_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExceptionResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateTrustedProfileAssignment : Update assignment
// Update a trusted profile assignment. Call this method to retry failed assignments or migrate the trusted profile in
// child accounts to a new version.
func (iamIdentity *IamIdentityV1) UpdateTrustedProfileAssignment(updateTrustedProfileAssignmentOptions *UpdateTrustedProfileAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateTrustedProfileAssignmentWithContext(context.Background(), updateTrustedProfileAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateTrustedProfileAssignmentWithContext is an alternate form of the UpdateTrustedProfileAssignment method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateTrustedProfileAssignmentWithContext(ctx context.Context, updateTrustedProfileAssignmentOptions *UpdateTrustedProfileAssignmentOptions) (result *TemplateAssignmentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateTrustedProfileAssignmentOptions, "updateTrustedProfileAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateTrustedProfileAssignmentOptions, "updateTrustedProfileAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *updateTrustedProfileAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateTrustedProfileAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateTrustedProfileAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateTrustedProfileAssignmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateTrustedProfileAssignmentOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateTrustedProfileAssignmentOptions.TemplateVersion != nil {
		body["template_version"] = updateTrustedProfileAssignmentOptions.TemplateVersion
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_trusted_profile_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAssignmentResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListProfileTemplates : List trusted profile templates
// List the trusted profile templates in an enterprise account.
func (iamIdentity *IamIdentityV1) ListProfileTemplates(listProfileTemplatesOptions *ListProfileTemplatesOptions) (result *TrustedProfileTemplateList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListProfileTemplatesWithContext(context.Background(), listProfileTemplatesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListProfileTemplatesWithContext is an alternate form of the ListProfileTemplates method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListProfileTemplatesWithContext(ctx context.Context, listProfileTemplatesOptions *ListProfileTemplatesOptions) (result *TrustedProfileTemplateList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProfileTemplatesOptions, "listProfileTemplatesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listProfileTemplatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListProfileTemplates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listProfileTemplatesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listProfileTemplatesOptions.AccountID))
	}
	if listProfileTemplatesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProfileTemplatesOptions.Limit))
	}
	if listProfileTemplatesOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listProfileTemplatesOptions.Pagetoken))
	}
	if listProfileTemplatesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listProfileTemplatesOptions.Sort))
	}
	if listProfileTemplatesOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listProfileTemplatesOptions.Order))
	}
	if listProfileTemplatesOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listProfileTemplatesOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_profile_templates", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfileTemplateList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateProfileTemplate : Create a trusted profile template
// Create a new trusted profile template in an enterprise account.
func (iamIdentity *IamIdentityV1) CreateProfileTemplate(createProfileTemplateOptions *CreateProfileTemplateOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateProfileTemplateWithContext(context.Background(), createProfileTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateProfileTemplateWithContext is an alternate form of the CreateProfileTemplate method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateProfileTemplateWithContext(ctx context.Context, createProfileTemplateOptions *CreateProfileTemplateOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProfileTemplateOptions, "createProfileTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createProfileTemplateOptions, "createProfileTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createProfileTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateProfileTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createProfileTemplateOptions.AccountID != nil {
		body["account_id"] = createProfileTemplateOptions.AccountID
	}
	if createProfileTemplateOptions.Name != nil {
		body["name"] = createProfileTemplateOptions.Name
	}
	if createProfileTemplateOptions.Description != nil {
		body["description"] = createProfileTemplateOptions.Description
	}
	if createProfileTemplateOptions.Profile != nil {
		body["profile"] = createProfileTemplateOptions.Profile
	}
	if createProfileTemplateOptions.PolicyTemplateReferences != nil {
		body["policy_template_references"] = createProfileTemplateOptions.PolicyTemplateReferences
	}
	if createProfileTemplateOptions.ActionControls != nil {
		body["action_controls"] = createProfileTemplateOptions.ActionControls
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_profile_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfileTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetLatestProfileTemplateVersion : Get latest version of a trusted profile template
// Get the latest version of a trusted profile template in an enterprise account.
func (iamIdentity *IamIdentityV1) GetLatestProfileTemplateVersion(getLatestProfileTemplateVersionOptions *GetLatestProfileTemplateVersionOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetLatestProfileTemplateVersionWithContext(context.Background(), getLatestProfileTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLatestProfileTemplateVersionWithContext is an alternate form of the GetLatestProfileTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetLatestProfileTemplateVersionWithContext(ctx context.Context, getLatestProfileTemplateVersionOptions *GetLatestProfileTemplateVersionOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLatestProfileTemplateVersionOptions, "getLatestProfileTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getLatestProfileTemplateVersionOptions, "getLatestProfileTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *getLatestProfileTemplateVersionOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates/{template_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLatestProfileTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetLatestProfileTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getLatestProfileTemplateVersionOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getLatestProfileTemplateVersionOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_latest_profile_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfileTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteAllVersionsOfProfileTemplate : Delete all versions of a trusted profile template
// Delete all versions of a trusted profile template in an enterprise account. If any version is assigned to child
// accounts, you must first delete the assignment.
func (iamIdentity *IamIdentityV1) DeleteAllVersionsOfProfileTemplate(deleteAllVersionsOfProfileTemplateOptions *DeleteAllVersionsOfProfileTemplateOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteAllVersionsOfProfileTemplateWithContext(context.Background(), deleteAllVersionsOfProfileTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAllVersionsOfProfileTemplateWithContext is an alternate form of the DeleteAllVersionsOfProfileTemplate method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteAllVersionsOfProfileTemplateWithContext(ctx context.Context, deleteAllVersionsOfProfileTemplateOptions *DeleteAllVersionsOfProfileTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAllVersionsOfProfileTemplateOptions, "deleteAllVersionsOfProfileTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAllVersionsOfProfileTemplateOptions, "deleteAllVersionsOfProfileTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *deleteAllVersionsOfProfileTemplateOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates/{template_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAllVersionsOfProfileTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteAllVersionsOfProfileTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_all_versions_of_profile_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListVersionsOfProfileTemplate : List trusted profile template versions
// List the versions of a trusted profile template in an enterprise account.
func (iamIdentity *IamIdentityV1) ListVersionsOfProfileTemplate(listVersionsOfProfileTemplateOptions *ListVersionsOfProfileTemplateOptions) (result *TrustedProfileTemplateList, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.ListVersionsOfProfileTemplateWithContext(context.Background(), listVersionsOfProfileTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListVersionsOfProfileTemplateWithContext is an alternate form of the ListVersionsOfProfileTemplate method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListVersionsOfProfileTemplateWithContext(ctx context.Context, listVersionsOfProfileTemplateOptions *ListVersionsOfProfileTemplateOptions) (result *TrustedProfileTemplateList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listVersionsOfProfileTemplateOptions, "listVersionsOfProfileTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listVersionsOfProfileTemplateOptions, "listVersionsOfProfileTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *listVersionsOfProfileTemplateOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates/{template_id}/versions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listVersionsOfProfileTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "ListVersionsOfProfileTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listVersionsOfProfileTemplateOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listVersionsOfProfileTemplateOptions.Limit))
	}
	if listVersionsOfProfileTemplateOptions.Pagetoken != nil {
		builder.AddQuery("pagetoken", fmt.Sprint(*listVersionsOfProfileTemplateOptions.Pagetoken))
	}
	if listVersionsOfProfileTemplateOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listVersionsOfProfileTemplateOptions.Sort))
	}
	if listVersionsOfProfileTemplateOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listVersionsOfProfileTemplateOptions.Order))
	}
	if listVersionsOfProfileTemplateOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*listVersionsOfProfileTemplateOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_versions_of_profile_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfileTemplateList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateProfileTemplateVersion : Create new version of a trusted profile template
// Create a new version of a trusted profile template in an enterprise account.
func (iamIdentity *IamIdentityV1) CreateProfileTemplateVersion(createProfileTemplateVersionOptions *CreateProfileTemplateVersionOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.CreateProfileTemplateVersionWithContext(context.Background(), createProfileTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateProfileTemplateVersionWithContext is an alternate form of the CreateProfileTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateProfileTemplateVersionWithContext(ctx context.Context, createProfileTemplateVersionOptions *CreateProfileTemplateVersionOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProfileTemplateVersionOptions, "createProfileTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createProfileTemplateVersionOptions, "createProfileTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *createProfileTemplateVersionOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates/{template_id}/versions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createProfileTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CreateProfileTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createProfileTemplateVersionOptions.AccountID != nil {
		body["account_id"] = createProfileTemplateVersionOptions.AccountID
	}
	if createProfileTemplateVersionOptions.Name != nil {
		body["name"] = createProfileTemplateVersionOptions.Name
	}
	if createProfileTemplateVersionOptions.Description != nil {
		body["description"] = createProfileTemplateVersionOptions.Description
	}
	if createProfileTemplateVersionOptions.Profile != nil {
		body["profile"] = createProfileTemplateVersionOptions.Profile
	}
	if createProfileTemplateVersionOptions.PolicyTemplateReferences != nil {
		body["policy_template_references"] = createProfileTemplateVersionOptions.PolicyTemplateReferences
	}
	if createProfileTemplateVersionOptions.ActionControls != nil {
		body["action_controls"] = createProfileTemplateVersionOptions.ActionControls
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_profile_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfileTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProfileTemplateVersion : Get version of trusted profile template
// Get a specific version of a trusted profile template in an enterprise account.
func (iamIdentity *IamIdentityV1) GetProfileTemplateVersion(getProfileTemplateVersionOptions *GetProfileTemplateVersionOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.GetProfileTemplateVersionWithContext(context.Background(), getProfileTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProfileTemplateVersionWithContext is an alternate form of the GetProfileTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetProfileTemplateVersionWithContext(ctx context.Context, getProfileTemplateVersionOptions *GetProfileTemplateVersionOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileTemplateVersionOptions, "getProfileTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProfileTemplateVersionOptions, "getProfileTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *getProfileTemplateVersionOptions.TemplateID,
		"version":     *getProfileTemplateVersionOptions.Version,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates/{template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getProfileTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "GetProfileTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getProfileTemplateVersionOptions.IncludeHistory != nil {
		builder.AddQuery("include_history", fmt.Sprint(*getProfileTemplateVersionOptions.IncludeHistory))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_profile_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfileTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateProfileTemplateVersion : Update version of trusted profile template
// Update a specific version of a trusted profile template in an enterprise account.
func (iamIdentity *IamIdentityV1) UpdateProfileTemplateVersion(updateProfileTemplateVersionOptions *UpdateProfileTemplateVersionOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	result, response, err = iamIdentity.UpdateProfileTemplateVersionWithContext(context.Background(), updateProfileTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateProfileTemplateVersionWithContext is an alternate form of the UpdateProfileTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateProfileTemplateVersionWithContext(ctx context.Context, updateProfileTemplateVersionOptions *UpdateProfileTemplateVersionOptions) (result *TrustedProfileTemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProfileTemplateVersionOptions, "updateProfileTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateProfileTemplateVersionOptions, "updateProfileTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *updateProfileTemplateVersionOptions.TemplateID,
		"version":     *updateProfileTemplateVersionOptions.Version,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates/{template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateProfileTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "UpdateProfileTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateProfileTemplateVersionOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateProfileTemplateVersionOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateProfileTemplateVersionOptions.AccountID != nil {
		body["account_id"] = updateProfileTemplateVersionOptions.AccountID
	}
	if updateProfileTemplateVersionOptions.Name != nil {
		body["name"] = updateProfileTemplateVersionOptions.Name
	}
	if updateProfileTemplateVersionOptions.Description != nil {
		body["description"] = updateProfileTemplateVersionOptions.Description
	}
	if updateProfileTemplateVersionOptions.Profile != nil {
		body["profile"] = updateProfileTemplateVersionOptions.Profile
	}
	if updateProfileTemplateVersionOptions.PolicyTemplateReferences != nil {
		body["policy_template_references"] = updateProfileTemplateVersionOptions.PolicyTemplateReferences
	}
	if updateProfileTemplateVersionOptions.ActionControls != nil {
		body["action_controls"] = updateProfileTemplateVersionOptions.ActionControls
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_profile_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfileTemplateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteProfileTemplateVersion : Delete version of trusted profile template
// Delete a specific version of a trusted profile template in an enterprise account. If the version is assigned to child
// accounts, you must first delete the assignment.
func (iamIdentity *IamIdentityV1) DeleteProfileTemplateVersion(deleteProfileTemplateVersionOptions *DeleteProfileTemplateVersionOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.DeleteProfileTemplateVersionWithContext(context.Background(), deleteProfileTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteProfileTemplateVersionWithContext is an alternate form of the DeleteProfileTemplateVersion method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteProfileTemplateVersionWithContext(ctx context.Context, deleteProfileTemplateVersionOptions *DeleteProfileTemplateVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProfileTemplateVersionOptions, "deleteProfileTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteProfileTemplateVersionOptions, "deleteProfileTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *deleteProfileTemplateVersionOptions.TemplateID,
		"version":     *deleteProfileTemplateVersionOptions.Version,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates/{template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteProfileTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "DeleteProfileTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_profile_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CommitProfileTemplate : Commit a template version
// Commit a specific version of a trusted profile template in an enterprise account. You must commit a template before
// you can assign it to child accounts. Once a template is committed, you can no longer modify the template.
func (iamIdentity *IamIdentityV1) CommitProfileTemplate(commitProfileTemplateOptions *CommitProfileTemplateOptions) (response *core.DetailedResponse, err error) {
	response, err = iamIdentity.CommitProfileTemplateWithContext(context.Background(), commitProfileTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CommitProfileTemplateWithContext is an alternate form of the CommitProfileTemplate method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CommitProfileTemplateWithContext(ctx context.Context, commitProfileTemplateOptions *CommitProfileTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(commitProfileTemplateOptions, "commitProfileTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(commitProfileTemplateOptions, "commitProfileTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *commitProfileTemplateOptions.TemplateID,
		"version":     *commitProfileTemplateOptions.Version,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profile_templates/{template_id}/versions/{version}/commit`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range commitProfileTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_identity", "V1", "CommitProfileTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "commit_profile_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.0")
}

// AccountBasedMfaEnrollment : AccountBasedMfaEnrollment struct
type AccountBasedMfaEnrollment struct {
	SecurityQuestions *MfaEnrollmentTypeStatus `json:"security_questions" validate:"required"`

	Totp *MfaEnrollmentTypeStatus `json:"totp" validate:"required"`

	Verisign *MfaEnrollmentTypeStatus `json:"verisign" validate:"required"`

	// The enrollment complies to the effective requirement.
	Complies *bool `json:"complies" validate:"required"`
}

// UnmarshalAccountBasedMfaEnrollment unmarshals an instance of AccountBasedMfaEnrollment from the specified map of raw messages.
func UnmarshalAccountBasedMfaEnrollment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountBasedMfaEnrollment)
	err = core.UnmarshalModel(m, "security_questions", &obj.SecurityQuestions, UnmarshalMfaEnrollmentTypeStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "security_questions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "totp", &obj.Totp, UnmarshalMfaEnrollmentTypeStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "totp-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "verisign", &obj.Verisign, UnmarshalMfaEnrollmentTypeStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "verisign-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "complies", &obj.Complies)
	if err != nil {
		err = core.SDKErrorf(err, "", "complies-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettingsAccountSection : AccountSettingsAccountSection struct
type AccountSettingsAccountSection struct {
	// Unique ID of the account.
	AccountID *string `json:"account_id,omitempty"`

	// Defines whether or not creating a service ID is access controlled. Valid values:
	//   * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
	// IDs, including the account owner
	//   * NOT_RESTRICTED - all members of an account can create service IDs
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreateServiceID *string `json:"restrict_create_service_id,omitempty"`

	// Defines whether or not creating platform API keys is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreatePlatformApikey *string `json:"restrict_create_platform_apikey,omitempty"`

	// Defines the IP addresses and subnets from which IAM tokens can be created for the account.
	AllowedIPAddresses *string `json:"allowed_ip_addresses,omitempty"`

	// Defines the MFA requirement for the user. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa,omitempty"`

	// List of users that are exempted from the MFA requirement of the account.
	UserMfa []EffectiveAccountSettingsUserMfa `json:"user_mfa,omitempty"`

	// History of the Account Settings.
	History []EnityHistoryRecord `json:"history,omitempty"`

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string `json:"session_expiration_in_seconds,omitempty"`

	// Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string `json:"session_invalidation_in_seconds,omitempty"`

	// Defines the max allowed sessions per identity required by the account. Valid values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string `json:"max_sessions_per_identity,omitempty"`

	// Defines the access token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '3600'
	//   * NOT_SET - To unset account setting and use service default.
	SystemAccessTokenExpirationInSeconds *string `json:"system_access_token_expiration_in_seconds,omitempty"`

	// Defines the refresh token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '259200'
	//   * NOT_SET - To unset account setting and use service default.
	SystemRefreshTokenExpirationInSeconds *string `json:"system_refresh_token_expiration_in_seconds,omitempty"`
}

// Constants associated with the AccountSettingsAccountSection.RestrictCreateServiceID property.
// Defines whether or not creating a service ID is access controlled. Valid values:
//   - RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
//
// IDs, including the account owner
//   - NOT_RESTRICTED - all members of an account can create service IDs
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsAccountSectionRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsAccountSectionRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	AccountSettingsAccountSectionRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsAccountSection.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   - RESTRICTED - to apply access control
//   - NOT_RESTRICTED - to remove access control
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsAccountSectionRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsAccountSectionRestrictCreatePlatformApikeyNotSetConst        = "NOT_SET"
	AccountSettingsAccountSectionRestrictCreatePlatformApikeyRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsAccountSection.Mfa property.
// Defines the MFA requirement for the user. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	AccountSettingsAccountSectionMfaLevel1Const     = "LEVEL1"
	AccountSettingsAccountSectionMfaLevel2Const     = "LEVEL2"
	AccountSettingsAccountSectionMfaLevel3Const     = "LEVEL3"
	AccountSettingsAccountSectionMfaNoneConst       = "NONE"
	AccountSettingsAccountSectionMfaNoneNoRopcConst = "NONE_NO_ROPC"
	AccountSettingsAccountSectionMfaTotpConst       = "TOTP"
	AccountSettingsAccountSectionMfaTotp4allConst   = "TOTP4ALL"
)

// UnmarshalAccountSettingsAccountSection unmarshals an instance of AccountSettingsAccountSection from the specified map of raw messages.
func UnmarshalAccountSettingsAccountSection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsAccountSection)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_service_id", &obj.RestrictCreateServiceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_service_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_platform_apikey", &obj.RestrictCreatePlatformApikey)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_platform_apikey-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_ip_addresses", &obj.AllowedIPAddresses)
	if err != nil {
		err = core.SDKErrorf(err, "", "allowed_ip_addresses-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "user_mfa", &obj.UserMfa, UnmarshalEffectiveAccountSettingsUserMfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_expiration_in_seconds", &obj.SessionExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_invalidation_in_seconds", &obj.SessionInvalidationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_invalidation_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_sessions_per_identity", &obj.MaxSessionsPerIdentity)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_sessions_per_identity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_access_token_expiration_in_seconds", &obj.SystemAccessTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_access_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_refresh_token_expiration_in_seconds", &obj.SystemRefreshTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_refresh_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettingsAssignedTemplatesSection : AccountSettingsAssignedTemplatesSection struct
type AccountSettingsAssignedTemplatesSection struct {
	// Template Id.
	TemplateID *string `json:"template_id,omitempty"`

	// Template version.
	TemplateVersion *int64 `json:"template_version,omitempty"`

	// Template name.
	TemplateName *string `json:"template_name,omitempty"`

	// Defines whether or not creating a service ID is access controlled. Valid values:
	//   * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
	// IDs, including the account owner
	//   * NOT_RESTRICTED - all members of an account can create service IDs
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreateServiceID *string `json:"restrict_create_service_id,omitempty"`

	// Defines whether or not creating platform API keys is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreatePlatformApikey *string `json:"restrict_create_platform_apikey,omitempty"`

	// Defines the IP addresses and subnets from which IAM tokens can be created for the account.
	AllowedIPAddresses *string `json:"allowed_ip_addresses,omitempty"`

	// Defines the MFA requirement for the user. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa,omitempty"`

	// List of users that are exempted from the MFA requirement of the account.
	UserMfa []EffectiveAccountSettingsUserMfa `json:"user_mfa,omitempty"`

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string `json:"session_expiration_in_seconds,omitempty"`

	// Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string `json:"session_invalidation_in_seconds,omitempty"`

	// Defines the max allowed sessions per identity required by the account. Valid values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string `json:"max_sessions_per_identity,omitempty"`

	// Defines the access token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '3600'
	//   * NOT_SET - To unset account setting and use service default.
	SystemAccessTokenExpirationInSeconds *string `json:"system_access_token_expiration_in_seconds,omitempty"`

	// Defines the refresh token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '259200'
	//   * NOT_SET - To unset account setting and use service default.
	SystemRefreshTokenExpirationInSeconds *string `json:"system_refresh_token_expiration_in_seconds,omitempty"`
}

// Constants associated with the AccountSettingsAssignedTemplatesSection.RestrictCreateServiceID property.
// Defines whether or not creating a service ID is access controlled. Valid values:
//   - RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
//
// IDs, including the account owner
//   - NOT_RESTRICTED - all members of an account can create service IDs
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsAssignedTemplatesSectionRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsAssignedTemplatesSectionRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	AccountSettingsAssignedTemplatesSectionRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsAssignedTemplatesSection.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   - RESTRICTED - to apply access control
//   - NOT_RESTRICTED - to remove access control
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsAssignedTemplatesSectionRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsAssignedTemplatesSectionRestrictCreatePlatformApikeyNotSetConst        = "NOT_SET"
	AccountSettingsAssignedTemplatesSectionRestrictCreatePlatformApikeyRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsAssignedTemplatesSection.Mfa property.
// Defines the MFA requirement for the user. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	AccountSettingsAssignedTemplatesSectionMfaLevel1Const     = "LEVEL1"
	AccountSettingsAssignedTemplatesSectionMfaLevel2Const     = "LEVEL2"
	AccountSettingsAssignedTemplatesSectionMfaLevel3Const     = "LEVEL3"
	AccountSettingsAssignedTemplatesSectionMfaNoneConst       = "NONE"
	AccountSettingsAssignedTemplatesSectionMfaNoneNoRopcConst = "NONE_NO_ROPC"
	AccountSettingsAssignedTemplatesSectionMfaTotpConst       = "TOTP"
	AccountSettingsAssignedTemplatesSectionMfaTotp4allConst   = "TOTP4ALL"
)

// UnmarshalAccountSettingsAssignedTemplatesSection unmarshals an instance of AccountSettingsAssignedTemplatesSection from the specified map of raw messages.
func UnmarshalAccountSettingsAssignedTemplatesSection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsAssignedTemplatesSection)
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_version", &obj.TemplateVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_name", &obj.TemplateName)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_service_id", &obj.RestrictCreateServiceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_service_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_platform_apikey", &obj.RestrictCreatePlatformApikey)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_platform_apikey-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_ip_addresses", &obj.AllowedIPAddresses)
	if err != nil {
		err = core.SDKErrorf(err, "", "allowed_ip_addresses-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "user_mfa", &obj.UserMfa, UnmarshalEffectiveAccountSettingsUserMfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_expiration_in_seconds", &obj.SessionExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_invalidation_in_seconds", &obj.SessionInvalidationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_invalidation_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_sessions_per_identity", &obj.MaxSessionsPerIdentity)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_sessions_per_identity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_access_token_expiration_in_seconds", &obj.SystemAccessTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_access_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_refresh_token_expiration_in_seconds", &obj.SystemRefreshTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_refresh_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettingsComponent : AccountSettingsComponent struct
type AccountSettingsComponent struct {
	// Defines whether or not creating a service ID is access controlled. Valid values:
	//   * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
	// IDs, including the account owner
	//   * NOT_RESTRICTED - all members of an account can create service IDs
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreateServiceID *string `json:"restrict_create_service_id,omitempty"`

	// Defines whether or not creating platform API keys is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreatePlatformApikey *string `json:"restrict_create_platform_apikey,omitempty"`

	// Defines the IP addresses and subnets from which IAM tokens can be created for the account.
	AllowedIPAddresses *string `json:"allowed_ip_addresses,omitempty"`

	// Defines the MFA trait for the account. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa,omitempty"`

	// List of users that are exempted from the MFA requirement of the account.
	UserMfa []AccountSettingsUserMfa `json:"user_mfa,omitempty"`

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string `json:"session_expiration_in_seconds,omitempty"`

	// Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string `json:"session_invalidation_in_seconds,omitempty"`

	// Defines the max allowed sessions per identity required by the account. Valid values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string `json:"max_sessions_per_identity,omitempty"`

	// Defines the access token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '3600'
	//   * NOT_SET - To unset account setting and use service default.
	SystemAccessTokenExpirationInSeconds *string `json:"system_access_token_expiration_in_seconds,omitempty"`

	// Defines the refresh token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '259200'
	//   * NOT_SET - To unset account setting and use service default.
	SystemRefreshTokenExpirationInSeconds *string `json:"system_refresh_token_expiration_in_seconds,omitempty"`
}

// Constants associated with the AccountSettingsComponent.RestrictCreateServiceID property.
// Defines whether or not creating a service ID is access controlled. Valid values:
//   - RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
//
// IDs, including the account owner
//   - NOT_RESTRICTED - all members of an account can create service IDs
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsComponentRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsComponentRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	AccountSettingsComponentRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsComponent.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   - RESTRICTED - to apply access control
//   - NOT_RESTRICTED - to remove access control
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsComponentRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsComponentRestrictCreatePlatformApikeyNotSetConst        = "NOT_SET"
	AccountSettingsComponentRestrictCreatePlatformApikeyRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsComponent.Mfa property.
// Defines the MFA trait for the account. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	AccountSettingsComponentMfaLevel1Const     = "LEVEL1"
	AccountSettingsComponentMfaLevel2Const     = "LEVEL2"
	AccountSettingsComponentMfaLevel3Const     = "LEVEL3"
	AccountSettingsComponentMfaNoneConst       = "NONE"
	AccountSettingsComponentMfaNoneNoRopcConst = "NONE_NO_ROPC"
	AccountSettingsComponentMfaTotpConst       = "TOTP"
	AccountSettingsComponentMfaTotp4allConst   = "TOTP4ALL"
)

// UnmarshalAccountSettingsComponent unmarshals an instance of AccountSettingsComponent from the specified map of raw messages.
func UnmarshalAccountSettingsComponent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsComponent)
	err = core.UnmarshalPrimitive(m, "restrict_create_service_id", &obj.RestrictCreateServiceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_service_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_platform_apikey", &obj.RestrictCreatePlatformApikey)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_platform_apikey-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_ip_addresses", &obj.AllowedIPAddresses)
	if err != nil {
		err = core.SDKErrorf(err, "", "allowed_ip_addresses-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "user_mfa", &obj.UserMfa, UnmarshalAccountSettingsUserMfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_expiration_in_seconds", &obj.SessionExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_invalidation_in_seconds", &obj.SessionInvalidationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_invalidation_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_sessions_per_identity", &obj.MaxSessionsPerIdentity)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_sessions_per_identity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_access_token_expiration_in_seconds", &obj.SystemAccessTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_access_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_refresh_token_expiration_in_seconds", &obj.SystemRefreshTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_refresh_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettingsEffectiveSection : AccountSettingsEffectiveSection struct
type AccountSettingsEffectiveSection struct {
	// Defines whether or not creating a service ID is access controlled. Valid values:
	//   * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
	// IDs, including the account owner
	//   * NOT_RESTRICTED - all members of an account can create service IDs
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreateServiceID *string `json:"restrict_create_service_id,omitempty"`

	// Defines whether or not creating platform API keys is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreatePlatformApikey *string `json:"restrict_create_platform_apikey,omitempty"`

	// Defines the IP addresses and subnets from which IAM tokens can be created for the account.
	AllowedIPAddresses *string `json:"allowed_ip_addresses,omitempty"`

	// Defines the MFA requirement for the user. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa,omitempty"`

	// List of users that are exempted from the MFA requirement of the account.
	UserMfa []EffectiveAccountSettingsUserMfa `json:"user_mfa,omitempty"`

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string `json:"session_expiration_in_seconds,omitempty"`

	// Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string `json:"session_invalidation_in_seconds,omitempty"`

	// Defines the max allowed sessions per identity required by the account. Valid values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string `json:"max_sessions_per_identity,omitempty"`

	// Defines the access token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '3600'
	//   * NOT_SET - To unset account setting and use service default.
	SystemAccessTokenExpirationInSeconds *string `json:"system_access_token_expiration_in_seconds,omitempty"`

	// Defines the refresh token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '259200'
	//   * NOT_SET - To unset account setting and use service default.
	SystemRefreshTokenExpirationInSeconds *string `json:"system_refresh_token_expiration_in_seconds,omitempty"`
}

// Constants associated with the AccountSettingsEffectiveSection.RestrictCreateServiceID property.
// Defines whether or not creating a service ID is access controlled. Valid values:
//   - RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
//
// IDs, including the account owner
//   - NOT_RESTRICTED - all members of an account can create service IDs
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsEffectiveSectionRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsEffectiveSectionRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	AccountSettingsEffectiveSectionRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsEffectiveSection.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   - RESTRICTED - to apply access control
//   - NOT_RESTRICTED - to remove access control
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsEffectiveSectionRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsEffectiveSectionRestrictCreatePlatformApikeyNotSetConst        = "NOT_SET"
	AccountSettingsEffectiveSectionRestrictCreatePlatformApikeyRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsEffectiveSection.Mfa property.
// Defines the MFA requirement for the user. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	AccountSettingsEffectiveSectionMfaLevel1Const     = "LEVEL1"
	AccountSettingsEffectiveSectionMfaLevel2Const     = "LEVEL2"
	AccountSettingsEffectiveSectionMfaLevel3Const     = "LEVEL3"
	AccountSettingsEffectiveSectionMfaNoneConst       = "NONE"
	AccountSettingsEffectiveSectionMfaNoneNoRopcConst = "NONE_NO_ROPC"
	AccountSettingsEffectiveSectionMfaTotpConst       = "TOTP"
	AccountSettingsEffectiveSectionMfaTotp4allConst   = "TOTP4ALL"
)

// UnmarshalAccountSettingsEffectiveSection unmarshals an instance of AccountSettingsEffectiveSection from the specified map of raw messages.
func UnmarshalAccountSettingsEffectiveSection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsEffectiveSection)
	err = core.UnmarshalPrimitive(m, "restrict_create_service_id", &obj.RestrictCreateServiceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_service_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_platform_apikey", &obj.RestrictCreatePlatformApikey)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_platform_apikey-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_ip_addresses", &obj.AllowedIPAddresses)
	if err != nil {
		err = core.SDKErrorf(err, "", "allowed_ip_addresses-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "user_mfa", &obj.UserMfa, UnmarshalEffectiveAccountSettingsUserMfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_expiration_in_seconds", &obj.SessionExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_invalidation_in_seconds", &obj.SessionInvalidationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_invalidation_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_sessions_per_identity", &obj.MaxSessionsPerIdentity)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_sessions_per_identity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_access_token_expiration_in_seconds", &obj.SystemAccessTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_access_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_refresh_token_expiration_in_seconds", &obj.SystemRefreshTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_refresh_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettingsResponse : Response body format for Account Settings REST requests.
type AccountSettingsResponse struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// Unique ID of the account.
	AccountID *string `json:"account_id" validate:"required"`

	// Defines whether or not creating a service ID is access controlled. Valid values:
	//   * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
	// IDs, including the account owner
	//   * NOT_RESTRICTED - all members of an account can create service IDs
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreateServiceID *string `json:"restrict_create_service_id" validate:"required"`

	// Defines whether or not creating platform API keys is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreatePlatformApikey *string `json:"restrict_create_platform_apikey" validate:"required"`

	// Defines the IP addresses and subnets from which IAM tokens can be created for the account.
	AllowedIPAddresses *string `json:"allowed_ip_addresses" validate:"required"`

	// Version of the account settings.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// Defines the MFA trait for the account. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa" validate:"required"`

	// List of users that are exempted from the MFA requirement of the account.
	UserMfa []AccountSettingsUserMfa `json:"user_mfa" validate:"required"`

	// History of the Account Settings.
	History []EnityHistoryRecord `json:"history,omitempty"`

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string `json:"session_expiration_in_seconds" validate:"required"`

	// Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string `json:"session_invalidation_in_seconds" validate:"required"`

	// Defines the max allowed sessions per identity required by the account. Valid values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string `json:"max_sessions_per_identity" validate:"required"`

	// Defines the access token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '3600'
	//   * NOT_SET - To unset account setting and use service default.
	SystemAccessTokenExpirationInSeconds *string `json:"system_access_token_expiration_in_seconds" validate:"required"`

	// Defines the refresh token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '259200'
	//   * NOT_SET - To unset account setting and use service default.
	SystemRefreshTokenExpirationInSeconds *string `json:"system_refresh_token_expiration_in_seconds" validate:"required"`
}

// Constants associated with the AccountSettingsResponse.RestrictCreateServiceID property.
// Defines whether or not creating a service ID is access controlled. Valid values:
//   - RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
//
// IDs, including the account owner
//   - NOT_RESTRICTED - all members of an account can create service IDs
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsResponseRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsResponseRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	AccountSettingsResponseRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsResponse.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   - RESTRICTED - to apply access control
//   - NOT_RESTRICTED - to remove access control
//   - NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsResponseRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsResponseRestrictCreatePlatformApikeyNotSetConst        = "NOT_SET"
	AccountSettingsResponseRestrictCreatePlatformApikeyRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsResponse.Mfa property.
// Defines the MFA trait for the account. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	AccountSettingsResponseMfaLevel1Const     = "LEVEL1"
	AccountSettingsResponseMfaLevel2Const     = "LEVEL2"
	AccountSettingsResponseMfaLevel3Const     = "LEVEL3"
	AccountSettingsResponseMfaNoneConst       = "NONE"
	AccountSettingsResponseMfaNoneNoRopcConst = "NONE_NO_ROPC"
	AccountSettingsResponseMfaTotpConst       = "TOTP"
	AccountSettingsResponseMfaTotp4allConst   = "TOTP4ALL"
)

// UnmarshalAccountSettingsResponse unmarshals an instance of AccountSettingsResponse from the specified map of raw messages.
func UnmarshalAccountSettingsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsResponse)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_service_id", &obj.RestrictCreateServiceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_service_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_platform_apikey", &obj.RestrictCreatePlatformApikey)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrict_create_platform_apikey-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_ip_addresses", &obj.AllowedIPAddresses)
	if err != nil {
		err = core.SDKErrorf(err, "", "allowed_ip_addresses-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "user_mfa", &obj.UserMfa, UnmarshalAccountSettingsUserMfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_expiration_in_seconds", &obj.SessionExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "session_invalidation_in_seconds", &obj.SessionInvalidationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "session_invalidation_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_sessions_per_identity", &obj.MaxSessionsPerIdentity)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_sessions_per_identity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_access_token_expiration_in_seconds", &obj.SystemAccessTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_access_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "system_refresh_token_expiration_in_seconds", &obj.SystemRefreshTokenExpirationInSeconds)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_refresh_token_expiration_in_seconds-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettingsTemplateList : AccountSettingsTemplateList struct
type AccountSettingsTemplateList struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// The offset of the current page.
	Offset *int64 `json:"offset,omitempty"`

	// Optional size of a single page.
	Limit *int64 `json:"limit,omitempty"`

	// Link to the first page.
	First *string `json:"first,omitempty"`

	// Link to the previous available page. If 'previous' property is not part of the response no previous page is
	// available.
	Previous *string `json:"previous,omitempty"`

	// Link to the next available page. If 'next' property is not part of the response no next page is available.
	Next *string `json:"next,omitempty"`

	// List of account settings templates based on the query paramters and the page size. The account_settings_templates
	// array is always part of the response but might be empty depending on the query parameter values provided.
	AccountSettingsTemplates []AccountSettingsTemplateResponse `json:"account_settings_templates" validate:"required"`
}

// UnmarshalAccountSettingsTemplateList unmarshals an instance of AccountSettingsTemplateList from the specified map of raw messages.
func UnmarshalAccountSettingsTemplateList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsTemplateList)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "account_settings_templates", &obj.AccountSettingsTemplates, UnmarshalAccountSettingsTemplateResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_settings_templates-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettingsTemplateResponse : Response body format for account settings template REST requests.
type AccountSettingsTemplateResponse struct {
	// ID of the the template.
	ID *string `json:"id" validate:"required"`

	// Version of the the template.
	Version *int64 `json:"version" validate:"required"`

	// ID of the account where the template resides.
	AccountID *string `json:"account_id" validate:"required"`

	// The name of the trusted profile template. This is visible only in the enterprise account.
	Name *string `json:"name" validate:"required"`

	// The description of the trusted profile template. Describe the template for enterprise account users.
	Description *string `json:"description,omitempty"`

	// Committed flag determines if the template is ready for assignment.
	Committed *bool `json:"committed" validate:"required"`

	AccountSettings *AccountSettingsComponent `json:"account_settings" validate:"required"`

	// History of the Template.
	History []EnityHistoryRecord `json:"history,omitempty"`

	// Entity tag for this templateId-version combination.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// Cloud resource name.
	CRN *string `json:"crn" validate:"required"`

	// Template Created At.
	CreatedAt *string `json:"created_at,omitempty"`

	// IAMid of the creator.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// Template last modified at.
	LastModifiedAt *string `json:"last_modified_at,omitempty"`

	// IAMid of the identity that made the latest modification.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`
}

// UnmarshalAccountSettingsTemplateResponse unmarshals an instance of AccountSettingsTemplateResponse from the specified map of raw messages.
func UnmarshalAccountSettingsTemplateResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsTemplateResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "committed", &obj.Committed)
	if err != nil {
		err = core.SDKErrorf(err, "", "committed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "account_settings", &obj.AccountSettings, UnmarshalAccountSettingsComponent)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_settings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettingsUserMfa : AccountSettingsUserMfa struct
type AccountSettingsUserMfa struct {
	// The iam_id of the user.
	IamID *string `json:"iam_id" validate:"required"`

	// Defines the MFA requirement for the user. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa" validate:"required"`
}

// Constants associated with the AccountSettingsUserMfa.Mfa property.
// Defines the MFA requirement for the user. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	AccountSettingsUserMfaMfaLevel1Const     = "LEVEL1"
	AccountSettingsUserMfaMfaLevel2Const     = "LEVEL2"
	AccountSettingsUserMfaMfaLevel3Const     = "LEVEL3"
	AccountSettingsUserMfaMfaNoneConst       = "NONE"
	AccountSettingsUserMfaMfaNoneNoRopcConst = "NONE_NO_ROPC"
	AccountSettingsUserMfaMfaTotpConst       = "TOTP"
	AccountSettingsUserMfaMfaTotp4allConst   = "TOTP4ALL"
)

// NewAccountSettingsUserMfa : Instantiate AccountSettingsUserMfa (Generic Model Constructor)
func (*IamIdentityV1) NewAccountSettingsUserMfa(iamID string, mfa string) (_model *AccountSettingsUserMfa, err error) {
	_model = &AccountSettingsUserMfa{
		IamID: core.StringPtr(iamID),
		Mfa:   core.StringPtr(mfa),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalAccountSettingsUserMfa unmarshals an instance of AccountSettingsUserMfa from the specified map of raw messages.
func UnmarshalAccountSettingsUserMfa(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsUserMfa)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "mfa-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionControls : ActionControls struct
type ActionControls struct {
	Identities *ActionControlsIdentities `json:"identities,omitempty"`

	Rules *ActionControlsRules `json:"rules" validate:"required"`

	Policies *ActionControlsPolicies `json:"policies" validate:"required"`
}

// NewActionControls : Instantiate ActionControls (Generic Model Constructor)
func (*IamIdentityV1) NewActionControls(rules *ActionControlsRules, policies *ActionControlsPolicies) (_model *ActionControls, err error) {
	_model = &ActionControls{
		Rules:    rules,
		Policies: policies,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalActionControls unmarshals an instance of ActionControls from the specified map of raw messages.
func UnmarshalActionControls(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionControls)
	err = core.UnmarshalModel(m, "identities", &obj.Identities, UnmarshalActionControlsIdentities)
	if err != nil {
		err = core.SDKErrorf(err, "", "identities-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalActionControlsRules)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policies", &obj.Policies, UnmarshalActionControlsPolicies)
	if err != nil {
		err = core.SDKErrorf(err, "", "policies-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionControlsIdentities : ActionControlsIdentities struct
type ActionControlsIdentities struct {
	Add *bool `json:"add" validate:"required"`

	Remove *bool `json:"remove" validate:"required"`
}

// NewActionControlsIdentities : Instantiate ActionControlsIdentities (Generic Model Constructor)
func (*IamIdentityV1) NewActionControlsIdentities(add bool, remove bool) (_model *ActionControlsIdentities, err error) {
	_model = &ActionControlsIdentities{
		Add:    core.BoolPtr(add),
		Remove: core.BoolPtr(remove),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalActionControlsIdentities unmarshals an instance of ActionControlsIdentities from the specified map of raw messages.
func UnmarshalActionControlsIdentities(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionControlsIdentities)
	err = core.UnmarshalPrimitive(m, "add", &obj.Add)
	if err != nil {
		err = core.SDKErrorf(err, "", "add-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "remove", &obj.Remove)
	if err != nil {
		err = core.SDKErrorf(err, "", "remove-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionControlsPolicies : ActionControlsPolicies struct
type ActionControlsPolicies struct {
	Add *bool `json:"add" validate:"required"`

	Remove *bool `json:"remove" validate:"required"`
}

// NewActionControlsPolicies : Instantiate ActionControlsPolicies (Generic Model Constructor)
func (*IamIdentityV1) NewActionControlsPolicies(add bool, remove bool) (_model *ActionControlsPolicies, err error) {
	_model = &ActionControlsPolicies{
		Add:    core.BoolPtr(add),
		Remove: core.BoolPtr(remove),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalActionControlsPolicies unmarshals an instance of ActionControlsPolicies from the specified map of raw messages.
func UnmarshalActionControlsPolicies(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionControlsPolicies)
	err = core.UnmarshalPrimitive(m, "add", &obj.Add)
	if err != nil {
		err = core.SDKErrorf(err, "", "add-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "remove", &obj.Remove)
	if err != nil {
		err = core.SDKErrorf(err, "", "remove-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionControlsRules : ActionControlsRules struct
type ActionControlsRules struct {
	Add *bool `json:"add" validate:"required"`

	Remove *bool `json:"remove" validate:"required"`
}

// NewActionControlsRules : Instantiate ActionControlsRules (Generic Model Constructor)
func (*IamIdentityV1) NewActionControlsRules(add bool, remove bool) (_model *ActionControlsRules, err error) {
	_model = &ActionControlsRules{
		Add:    core.BoolPtr(add),
		Remove: core.BoolPtr(remove),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalActionControlsRules unmarshals an instance of ActionControlsRules from the specified map of raw messages.
func UnmarshalActionControlsRules(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionControlsRules)
	err = core.UnmarshalPrimitive(m, "add", &obj.Add)
	if err != nil {
		err = core.SDKErrorf(err, "", "add-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "remove", &obj.Remove)
	if err != nil {
		err = core.SDKErrorf(err, "", "remove-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Activity : Activity struct
type Activity struct {
	// Time when the entity was last authenticated.
	LastAuthn *string `json:"last_authn,omitempty"`

	// Authentication count, number of times the entity was authenticated.
	AuthnCount *int64 `json:"authn_count" validate:"required"`
}

// UnmarshalActivity unmarshals an instance of Activity from the specified map of raw messages.
func UnmarshalActivity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Activity)
	err = core.UnmarshalPrimitive(m, "last_authn", &obj.LastAuthn)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_authn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "authn_count", &obj.AuthnCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "authn_count-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APIKey : Response body format for API key V1 REST requests.
type APIKey struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// Unique identifier of this API Key.
	ID *string `json:"id" validate:"required"`

	// Version of the API Key details object. You need to specify this value when updating the API key to avoid stale
	// updates.
	EntityTag *string `json:"entity_tag,omitempty"`

	// Cloud Resource Name of the item. Example Cloud Resource Name:
	// 'crn:v1:bluemix:public:iam-identity:us-south:a/myaccount::apikey:1234-9012-5678'.
	CRN *string `json:"crn" validate:"required"`

	// The API key cannot be changed if set to true.
	Locked *bool `json:"locked" validate:"required"`

	// Defines if API key is disabled, API key cannot be used if 'disabled' is set to true.
	Disabled *bool `json:"disabled,omitempty"`

	// If set contains a date time string of the creation date in ISO format.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// IAM ID of the user or service which created the API key.
	CreatedBy *string `json:"created_by" validate:"required"`

	// If set contains a date time string of the last modification date in ISO format.
	ModifiedAt *strfmt.DateTime `json:"modified_at,omitempty"`

	// Name of the API key. The name is not checked for uniqueness. Therefore multiple names with the same value can exist.
	// Access is done via the UUID of the API key.
	Name *string `json:"name" validate:"required"`

	// Defines if the API key supports sessions. Sessions are only supported for user apikeys.
	SupportSessions *bool `json:"support_sessions,omitempty"`

	// Defines the action to take when API key is leaked, valid values are 'none', 'disable' and 'delete'.
	ActionWhenLeaked *string `json:"action_when_leaked,omitempty"`

	// The optional description of the API key. The 'description' property is only available if a description was provided
	// during a create of an API key.
	Description *string `json:"description,omitempty"`

	// The iam_id that this API key authenticates.
	IamID *string `json:"iam_id" validate:"required"`

	// ID of the account that this API key authenticates for.
	AccountID *string `json:"account_id" validate:"required"`

	// The API key value. This property only contains the API key value for the following cases: create an API key, update
	// a service ID API key that stores the API key value as retrievable, or get a service ID API key that stores the API
	// key value as retrievable. All other operations don't return the API key value, for example all user API key related
	// operations, except for create, don't contain the API key value.
	Apikey *string `json:"apikey" validate:"required"`

	// History of the API key.
	History []EnityHistoryRecord `json:"history,omitempty"`

	Activity *Activity `json:"activity,omitempty"`
}

// UnmarshalAPIKey unmarshals an instance of APIKey from the specified map of raw messages.
func UnmarshalAPIKey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIKey)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		err = core.SDKErrorf(err, "", "locked-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "support_sessions", &obj.SupportSessions)
	if err != nil {
		err = core.SDKErrorf(err, "", "support_sessions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "action_when_leaked", &obj.ActionWhenLeaked)
	if err != nil {
		err = core.SDKErrorf(err, "", "action_when_leaked-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "apikey", &obj.Apikey)
	if err != nil {
		err = core.SDKErrorf(err, "", "apikey-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "activity", &obj.Activity, UnmarshalActivity)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APIKeyInsideCreateServiceIDRequest : Parameters for the API key in the Create service Id V1 REST request.
type APIKeyInsideCreateServiceIDRequest struct {
	// Name of the API key. The name is not checked for uniqueness. Therefore multiple names with the same value can exist.
	// Access is done via the UUID of the API key.
	Name *string `json:"name" validate:"required"`

	// The optional description of the API key. The 'description' property is only available if a description was provided
	// during a create of an API key.
	Description *string `json:"description,omitempty"`

	// You can optionally passthrough the API key value for this API key. If passed, a minimum length validation of 32
	// characters for that apiKey value is done, i.e. the value can contain any characters and can even be non-URL safe,
	// but the minimum length requirement must be met. If omitted, the API key management will create an URL safe opaque
	// API key value. The value of the API key is checked for uniqueness. Ensure enough variations when passing in this
	// value.
	Apikey *string `json:"apikey,omitempty"`

	// Send true or false to set whether the API key value is retrievable in the future by using the Get details of an API
	// key request. If you create an API key for a user, you must specify `false` or omit the value. We don't allow storing
	// of API keys for users.
	StoreValue *bool `json:"store_value,omitempty"`
}

// NewAPIKeyInsideCreateServiceIDRequest : Instantiate APIKeyInsideCreateServiceIDRequest (Generic Model Constructor)
func (*IamIdentityV1) NewAPIKeyInsideCreateServiceIDRequest(name string) (_model *APIKeyInsideCreateServiceIDRequest, err error) {
	_model = &APIKeyInsideCreateServiceIDRequest{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalAPIKeyInsideCreateServiceIDRequest unmarshals an instance of APIKeyInsideCreateServiceIDRequest from the specified map of raw messages.
func UnmarshalAPIKeyInsideCreateServiceIDRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIKeyInsideCreateServiceIDRequest)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "apikey", &obj.Apikey)
	if err != nil {
		err = core.SDKErrorf(err, "", "apikey-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "store_value", &obj.StoreValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "store_value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APIKeyList : Response body format for the List API keys V1 REST request.
type APIKeyList struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// The offset of the current page.
	Offset *int64 `json:"offset,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Limit *int64 `json:"limit,omitempty"`

	// Link to the first page.
	First *string `json:"first,omitempty"`

	// Link to the previous available page. If 'previous' property is not part of the response no previous page is
	// available.
	Previous *string `json:"previous,omitempty"`

	// Link to the next available page. If 'next' property is not part of the response no next page is available.
	Next *string `json:"next,omitempty"`

	// List of API keys based on the query paramters and the page size. The apikeys array is always part of the response
	// but might be empty depending on the query parameters values provided.
	Apikeys []APIKey `json:"apikeys" validate:"required"`
}

// UnmarshalAPIKeyList unmarshals an instance of APIKeyList from the specified map of raw messages.
func UnmarshalAPIKeyList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIKeyList)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "apikeys", &obj.Apikeys, UnmarshalAPIKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "apikeys-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApikeyActivity : Apikeys activity details.
type ApikeyActivity struct {
	// Unique id of the apikey.
	ID *string `json:"id" validate:"required"`

	// Name provided during creation of the apikey.
	Name *string `json:"name,omitempty"`

	// Type of the apikey. Supported values are `serviceid` and `user`.
	Type *string `json:"type" validate:"required"`

	// serviceid details will be present if type is `serviceid`.
	Serviceid *ApikeyActivityServiceid `json:"serviceid,omitempty"`

	// user details will be present if type is `user`.
	User *ApikeyActivityUser `json:"user,omitempty"`

	// Time when the apikey was last authenticated.
	LastAuthn *string `json:"last_authn,omitempty"`
}

// UnmarshalApikeyActivity unmarshals an instance of ApikeyActivity from the specified map of raw messages.
func UnmarshalApikeyActivity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApikeyActivity)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "serviceid", &obj.Serviceid, UnmarshalApikeyActivityServiceid)
	if err != nil {
		err = core.SDKErrorf(err, "", "serviceid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "user", &obj.User, UnmarshalApikeyActivityUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_authn", &obj.LastAuthn)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_authn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApikeyActivityServiceid : serviceid details will be present if type is `serviceid`.
type ApikeyActivityServiceid struct {
	// Unique identifier of this Service Id.
	ID *string `json:"id,omitempty"`

	// Name provided during creation of the serviceid.
	Name *string `json:"name,omitempty"`
}

// UnmarshalApikeyActivityServiceid unmarshals an instance of ApikeyActivityServiceid from the specified map of raw messages.
func UnmarshalApikeyActivityServiceid(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApikeyActivityServiceid)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApikeyActivityUser : user details will be present if type is `user`.
type ApikeyActivityUser struct {
	// IAMid of the user.
	IamID *string `json:"iam_id,omitempty"`

	// Name of the user.
	Name *string `json:"name,omitempty"`

	// Username of the user.
	Username *string `json:"username,omitempty"`

	// Email of the user.
	Email *string `json:"email,omitempty"`
}

// UnmarshalApikeyActivityUser unmarshals an instance of ApikeyActivityUser from the specified map of raw messages.
func UnmarshalApikeyActivityUser(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApikeyActivityUser)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		err = core.SDKErrorf(err, "", "username-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CommitAccountSettingsTemplateOptions : The CommitAccountSettingsTemplate options.
type CommitAccountSettingsTemplateOptions struct {
	// ID of the account settings template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Version of the account settings template.
	Version *string `json:"version" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCommitAccountSettingsTemplateOptions : Instantiate CommitAccountSettingsTemplateOptions
func (*IamIdentityV1) NewCommitAccountSettingsTemplateOptions(templateID string, version string) *CommitAccountSettingsTemplateOptions {
	return &CommitAccountSettingsTemplateOptions{
		TemplateID: core.StringPtr(templateID),
		Version:    core.StringPtr(version),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *CommitAccountSettingsTemplateOptions) SetTemplateID(templateID string) *CommitAccountSettingsTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *CommitAccountSettingsTemplateOptions) SetVersion(version string) *CommitAccountSettingsTemplateOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CommitAccountSettingsTemplateOptions) SetHeaders(param map[string]string) *CommitAccountSettingsTemplateOptions {
	options.Headers = param
	return options
}

// CommitProfileTemplateOptions : The CommitProfileTemplate options.
type CommitProfileTemplateOptions struct {
	// ID of the trusted profile template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Version of the Profile Template.
	Version *string `json:"version" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCommitProfileTemplateOptions : Instantiate CommitProfileTemplateOptions
func (*IamIdentityV1) NewCommitProfileTemplateOptions(templateID string, version string) *CommitProfileTemplateOptions {
	return &CommitProfileTemplateOptions{
		TemplateID: core.StringPtr(templateID),
		Version:    core.StringPtr(version),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *CommitProfileTemplateOptions) SetTemplateID(templateID string) *CommitProfileTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *CommitProfileTemplateOptions) SetVersion(version string) *CommitProfileTemplateOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CommitProfileTemplateOptions) SetHeaders(param map[string]string) *CommitProfileTemplateOptions {
	options.Headers = param
	return options
}

// CreateAccountSettingsAssignmentOptions : The CreateAccountSettingsAssignment options.
type CreateAccountSettingsAssignmentOptions struct {
	// ID of the template to assign.
	TemplateID *string `json:"template_id" validate:"required"`

	// Version of the template to assign.
	TemplateVersion *int64 `json:"template_version" validate:"required"`

	// Type of target to deploy to.
	TargetType *string `json:"target_type" validate:"required"`

	// Identifier of target to deploy to.
	Target *string `json:"target" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateAccountSettingsAssignmentOptions.TargetType property.
// Type of target to deploy to.
const (
	CreateAccountSettingsAssignmentOptionsTargetTypeAccountConst      = "Account"
	CreateAccountSettingsAssignmentOptionsTargetTypeAccountgroupConst = "AccountGroup"
)

// NewCreateAccountSettingsAssignmentOptions : Instantiate CreateAccountSettingsAssignmentOptions
func (*IamIdentityV1) NewCreateAccountSettingsAssignmentOptions(templateID string, templateVersion int64, targetType string, target string) *CreateAccountSettingsAssignmentOptions {
	return &CreateAccountSettingsAssignmentOptions{
		TemplateID:      core.StringPtr(templateID),
		TemplateVersion: core.Int64Ptr(templateVersion),
		TargetType:      core.StringPtr(targetType),
		Target:          core.StringPtr(target),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *CreateAccountSettingsAssignmentOptions) SetTemplateID(templateID string) *CreateAccountSettingsAssignmentOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetTemplateVersion : Allow user to set TemplateVersion
func (_options *CreateAccountSettingsAssignmentOptions) SetTemplateVersion(templateVersion int64) *CreateAccountSettingsAssignmentOptions {
	_options.TemplateVersion = core.Int64Ptr(templateVersion)
	return _options
}

// SetTargetType : Allow user to set TargetType
func (_options *CreateAccountSettingsAssignmentOptions) SetTargetType(targetType string) *CreateAccountSettingsAssignmentOptions {
	_options.TargetType = core.StringPtr(targetType)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *CreateAccountSettingsAssignmentOptions) SetTarget(target string) *CreateAccountSettingsAssignmentOptions {
	_options.Target = core.StringPtr(target)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccountSettingsAssignmentOptions) SetHeaders(param map[string]string) *CreateAccountSettingsAssignmentOptions {
	options.Headers = param
	return options
}

// CreateAccountSettingsTemplateOptions : The CreateAccountSettingsTemplate options.
type CreateAccountSettingsTemplateOptions struct {
	// ID of the account where the template resides.
	AccountID *string `json:"account_id,omitempty"`

	// The name of the trusted profile template. This is visible only in the enterprise account.
	Name *string `json:"name,omitempty"`

	// The description of the trusted profile template. Describe the template for enterprise account users.
	Description *string `json:"description,omitempty"`

	AccountSettings *AccountSettingsComponent `json:"account_settings,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateAccountSettingsTemplateOptions : Instantiate CreateAccountSettingsTemplateOptions
func (*IamIdentityV1) NewCreateAccountSettingsTemplateOptions() *CreateAccountSettingsTemplateOptions {
	return &CreateAccountSettingsTemplateOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateAccountSettingsTemplateOptions) SetAccountID(accountID string) *CreateAccountSettingsTemplateOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateAccountSettingsTemplateOptions) SetName(name string) *CreateAccountSettingsTemplateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateAccountSettingsTemplateOptions) SetDescription(description string) *CreateAccountSettingsTemplateOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAccountSettings : Allow user to set AccountSettings
func (_options *CreateAccountSettingsTemplateOptions) SetAccountSettings(accountSettings *AccountSettingsComponent) *CreateAccountSettingsTemplateOptions {
	_options.AccountSettings = accountSettings
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccountSettingsTemplateOptions) SetHeaders(param map[string]string) *CreateAccountSettingsTemplateOptions {
	options.Headers = param
	return options
}

// CreateAccountSettingsTemplateVersionOptions : The CreateAccountSettingsTemplateVersion options.
type CreateAccountSettingsTemplateVersionOptions struct {
	// ID of the account settings template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// ID of the account where the template resides.
	AccountID *string `json:"account_id,omitempty"`

	// The name of the trusted profile template. This is visible only in the enterprise account.
	Name *string `json:"name,omitempty"`

	// The description of the trusted profile template. Describe the template for enterprise account users.
	Description *string `json:"description,omitempty"`

	AccountSettings *AccountSettingsComponent `json:"account_settings,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateAccountSettingsTemplateVersionOptions : Instantiate CreateAccountSettingsTemplateVersionOptions
func (*IamIdentityV1) NewCreateAccountSettingsTemplateVersionOptions(templateID string) *CreateAccountSettingsTemplateVersionOptions {
	return &CreateAccountSettingsTemplateVersionOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *CreateAccountSettingsTemplateVersionOptions) SetTemplateID(templateID string) *CreateAccountSettingsTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateAccountSettingsTemplateVersionOptions) SetAccountID(accountID string) *CreateAccountSettingsTemplateVersionOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateAccountSettingsTemplateVersionOptions) SetName(name string) *CreateAccountSettingsTemplateVersionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateAccountSettingsTemplateVersionOptions) SetDescription(description string) *CreateAccountSettingsTemplateVersionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAccountSettings : Allow user to set AccountSettings
func (_options *CreateAccountSettingsTemplateVersionOptions) SetAccountSettings(accountSettings *AccountSettingsComponent) *CreateAccountSettingsTemplateVersionOptions {
	_options.AccountSettings = accountSettings
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccountSettingsTemplateVersionOptions) SetHeaders(param map[string]string) *CreateAccountSettingsTemplateVersionOptions {
	options.Headers = param
	return options
}

// CreateAPIKeyOptions : The CreateAPIKey options.
type CreateAPIKeyOptions struct {
	// Name of the API key. The name is not checked for uniqueness. Therefore multiple names with the same value can exist.
	// Access is done via the UUID of the API key.
	Name *string `json:"name" validate:"required"`

	// The iam_id that this API key authenticates.
	IamID *string `json:"iam_id" validate:"required"`

	// The optional description of the API key. The 'description' property is only available if a description was provided
	// during a create of an API key.
	Description *string `json:"description,omitempty"`

	// The account ID of the API key.
	AccountID *string `json:"account_id,omitempty"`

	// You can optionally passthrough the API key value for this API key. If passed, a minimum length validation of 32
	// characters for that apiKey value is done, i.e. the value can contain any characters and can even be non-URL safe,
	// but the minimum length requirement must be met. If omitted, the API key management will create an URL safe opaque
	// API key value. The value of the API key is checked for uniqueness. Ensure enough variations when passing in this
	// value.
	Apikey *string `json:"apikey,omitempty"`

	// Send true or false to set whether the API key value is retrievable in the future by using the Get details of an API
	// key request. If you create an API key for a user, you must specify `false` or omit the value. We don't allow storing
	// of API keys for users.
	StoreValue *bool `json:"store_value,omitempty"`

	// Defines if the API key supports sessions. Sessions are only supported for user apikeys.
	SupportSessions *bool `json:"support_sessions,omitempty"`

	// Defines the action to take when API key is leaked, valid values are 'none', 'disable' and 'delete'.
	ActionWhenLeaked *string `json:"action_when_leaked,omitempty"`

	// Indicates if the API key is locked for further write operations. False by default.
	EntityLock *string `json:"Entity-Lock,omitempty"`

	// Indicates if the API key is disabled. False by default.
	EntityDisable *string `json:"Entity-Disable,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateAPIKeyOptions : Instantiate CreateAPIKeyOptions
func (*IamIdentityV1) NewCreateAPIKeyOptions(name string, iamID string) *CreateAPIKeyOptions {
	return &CreateAPIKeyOptions{
		Name:  core.StringPtr(name),
		IamID: core.StringPtr(iamID),
	}
}

// SetName : Allow user to set Name
func (_options *CreateAPIKeyOptions) SetName(name string) *CreateAPIKeyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIamID : Allow user to set IamID
func (_options *CreateAPIKeyOptions) SetIamID(iamID string) *CreateAPIKeyOptions {
	_options.IamID = core.StringPtr(iamID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateAPIKeyOptions) SetDescription(description string) *CreateAPIKeyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateAPIKeyOptions) SetAccountID(accountID string) *CreateAPIKeyOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetApikey : Allow user to set Apikey
func (_options *CreateAPIKeyOptions) SetApikey(apikey string) *CreateAPIKeyOptions {
	_options.Apikey = core.StringPtr(apikey)
	return _options
}

// SetStoreValue : Allow user to set StoreValue
func (_options *CreateAPIKeyOptions) SetStoreValue(storeValue bool) *CreateAPIKeyOptions {
	_options.StoreValue = core.BoolPtr(storeValue)
	return _options
}

// SetSupportSessions : Allow user to set SupportSessions
func (_options *CreateAPIKeyOptions) SetSupportSessions(supportSessions bool) *CreateAPIKeyOptions {
	_options.SupportSessions = core.BoolPtr(supportSessions)
	return _options
}

// SetActionWhenLeaked : Allow user to set ActionWhenLeaked
func (_options *CreateAPIKeyOptions) SetActionWhenLeaked(actionWhenLeaked string) *CreateAPIKeyOptions {
	_options.ActionWhenLeaked = core.StringPtr(actionWhenLeaked)
	return _options
}

// SetEntityLock : Allow user to set EntityLock
func (_options *CreateAPIKeyOptions) SetEntityLock(entityLock string) *CreateAPIKeyOptions {
	_options.EntityLock = core.StringPtr(entityLock)
	return _options
}

// SetEntityDisable : Allow user to set EntityDisable
func (_options *CreateAPIKeyOptions) SetEntityDisable(entityDisable string) *CreateAPIKeyOptions {
	_options.EntityDisable = core.StringPtr(entityDisable)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAPIKeyOptions) SetHeaders(param map[string]string) *CreateAPIKeyOptions {
	options.Headers = param
	return options
}

// CreateClaimRuleOptions : The CreateClaimRule options.
type CreateClaimRuleOptions struct {
	// ID of the trusted profile to create a claim rule.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Type of the claim rule, either 'Profile-SAML' or 'Profile-CR'.
	Type *string `json:"type" validate:"required"`

	// Conditions of this claim rule.
	Conditions []ProfileClaimRuleConditions `json:"conditions" validate:"required"`

	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// Name of the claim rule to be created or updated.
	Name *string `json:"name,omitempty"`

	// The realm name of the Idp this claim rule applies to. This field is required only if the type is specified as
	// 'Profile-SAML'.
	RealmName *string `json:"realm_name,omitempty"`

	// The compute resource type the rule applies to, required only if type is specified as 'Profile-CR'. Valid values are
	// VSI, IKS_SA, ROKS_SA.
	CrType *string `json:"cr_type,omitempty"`

	// Session expiration in seconds, only required if type is 'Profile-SAML'.
	Expiration *int64 `json:"expiration,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateClaimRuleOptions : Instantiate CreateClaimRuleOptions
func (*IamIdentityV1) NewCreateClaimRuleOptions(profileID string, typeVar string, conditions []ProfileClaimRuleConditions) *CreateClaimRuleOptions {
	return &CreateClaimRuleOptions{
		ProfileID:  core.StringPtr(profileID),
		Type:       core.StringPtr(typeVar),
		Conditions: conditions,
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *CreateClaimRuleOptions) SetProfileID(profileID string) *CreateClaimRuleOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateClaimRuleOptions) SetType(typeVar string) *CreateClaimRuleOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetConditions : Allow user to set Conditions
func (_options *CreateClaimRuleOptions) SetConditions(conditions []ProfileClaimRuleConditions) *CreateClaimRuleOptions {
	_options.Conditions = conditions
	return _options
}

// SetContext : Allow user to set Context
func (_options *CreateClaimRuleOptions) SetContext(context *ResponseContext) *CreateClaimRuleOptions {
	_options.Context = context
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateClaimRuleOptions) SetName(name string) *CreateClaimRuleOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetRealmName : Allow user to set RealmName
func (_options *CreateClaimRuleOptions) SetRealmName(realmName string) *CreateClaimRuleOptions {
	_options.RealmName = core.StringPtr(realmName)
	return _options
}

// SetCrType : Allow user to set CrType
func (_options *CreateClaimRuleOptions) SetCrType(crType string) *CreateClaimRuleOptions {
	_options.CrType = core.StringPtr(crType)
	return _options
}

// SetExpiration : Allow user to set Expiration
func (_options *CreateClaimRuleOptions) SetExpiration(expiration int64) *CreateClaimRuleOptions {
	_options.Expiration = core.Int64Ptr(expiration)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateClaimRuleOptions) SetHeaders(param map[string]string) *CreateClaimRuleOptions {
	options.Headers = param
	return options
}

// CreateLinkOptions : The CreateLink options.
type CreateLinkOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// The compute resource type. Valid values are VSI, IKS_SA, ROKS_SA.
	CrType *string `json:"cr_type" validate:"required"`

	// Link details.
	Link *CreateProfileLinkRequestLink `json:"link" validate:"required"`

	// Optional name of the Link.
	Name *string `json:"name,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateLinkOptions : Instantiate CreateLinkOptions
func (*IamIdentityV1) NewCreateLinkOptions(profileID string, crType string, link *CreateProfileLinkRequestLink) *CreateLinkOptions {
	return &CreateLinkOptions{
		ProfileID: core.StringPtr(profileID),
		CrType:    core.StringPtr(crType),
		Link:      link,
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *CreateLinkOptions) SetProfileID(profileID string) *CreateLinkOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetCrType : Allow user to set CrType
func (_options *CreateLinkOptions) SetCrType(crType string) *CreateLinkOptions {
	_options.CrType = core.StringPtr(crType)
	return _options
}

// SetLink : Allow user to set Link
func (_options *CreateLinkOptions) SetLink(link *CreateProfileLinkRequestLink) *CreateLinkOptions {
	_options.Link = link
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateLinkOptions) SetName(name string) *CreateLinkOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLinkOptions) SetHeaders(param map[string]string) *CreateLinkOptions {
	options.Headers = param
	return options
}

// CreateMfaReportOptions : The CreateMfaReport options.
type CreateMfaReportOptions struct {
	// ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Optional report type. The supported value is 'mfa_status'. List MFA enrollment status for all the identities.
	Type *string `json:"type,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateMfaReportOptions : Instantiate CreateMfaReportOptions
func (*IamIdentityV1) NewCreateMfaReportOptions(accountID string) *CreateMfaReportOptions {
	return &CreateMfaReportOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateMfaReportOptions) SetAccountID(accountID string) *CreateMfaReportOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateMfaReportOptions) SetType(typeVar string) *CreateMfaReportOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateMfaReportOptions) SetHeaders(param map[string]string) *CreateMfaReportOptions {
	options.Headers = param
	return options
}

// CreateProfileLinkRequestLink : Link details.
type CreateProfileLinkRequestLink struct {
	// The CRN of the compute resource.
	CRN *string `json:"crn" validate:"required"`

	// The compute resource namespace, only required if cr_type is IKS_SA or ROKS_SA.
	Namespace *string `json:"namespace" validate:"required"`

	// Name of the compute resource, only required if cr_type is IKS_SA or ROKS_SA.
	Name *string `json:"name,omitempty"`
}

// NewCreateProfileLinkRequestLink : Instantiate CreateProfileLinkRequestLink (Generic Model Constructor)
func (*IamIdentityV1) NewCreateProfileLinkRequestLink(crn string, namespace string) (_model *CreateProfileLinkRequestLink, err error) {
	_model = &CreateProfileLinkRequestLink{
		CRN:       core.StringPtr(crn),
		Namespace: core.StringPtr(namespace),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalCreateProfileLinkRequestLink unmarshals an instance of CreateProfileLinkRequestLink from the specified map of raw messages.
func UnmarshalCreateProfileLinkRequestLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateProfileLinkRequestLink)
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		err = core.SDKErrorf(err, "", "namespace-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateProfileOptions : The CreateProfile options.
type CreateProfileOptions struct {
	// Name of the trusted profile. The name is checked for uniqueness. Therefore trusted profiles with the same names can
	// not exist in the same account.
	Name *string `json:"name" validate:"required"`

	// The account ID of the trusted profile.
	AccountID *string `json:"account_id" validate:"required"`

	// The optional description of the trusted profile. The 'description' property is only available if a description was
	// provided during creation of trusted profile.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateProfileOptions : Instantiate CreateProfileOptions
func (*IamIdentityV1) NewCreateProfileOptions(name string, accountID string) *CreateProfileOptions {
	return &CreateProfileOptions{
		Name:      core.StringPtr(name),
		AccountID: core.StringPtr(accountID),
	}
}

// SetName : Allow user to set Name
func (_options *CreateProfileOptions) SetName(name string) *CreateProfileOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateProfileOptions) SetAccountID(accountID string) *CreateProfileOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateProfileOptions) SetDescription(description string) *CreateProfileOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProfileOptions) SetHeaders(param map[string]string) *CreateProfileOptions {
	options.Headers = param
	return options
}

// CreateProfileTemplateOptions : The CreateProfileTemplate options.
type CreateProfileTemplateOptions struct {
	// ID of the account where the template resides.
	AccountID *string `json:"account_id,omitempty"`

	// The name of the trusted profile template. This is visible only in the enterprise account. Required field when
	// creating a new template. Otherwise this field is optional. If the field is included it will change the name value
	// for all existing versions of the template.
	Name *string `json:"name,omitempty"`

	// The description of the trusted profile template. Describe the template for enterprise account users.
	Description *string `json:"description,omitempty"`

	// Input body parameters for the TemplateProfileComponent.
	Profile *TemplateProfileComponentRequest `json:"profile,omitempty"`

	// Existing policy templates that you can reference to assign access in the trusted profile component.
	PolicyTemplateReferences []PolicyTemplateReference `json:"policy_template_references,omitempty"`

	ActionControls *ActionControls `json:"action_controls,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateProfileTemplateOptions : Instantiate CreateProfileTemplateOptions
func (*IamIdentityV1) NewCreateProfileTemplateOptions() *CreateProfileTemplateOptions {
	return &CreateProfileTemplateOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateProfileTemplateOptions) SetAccountID(accountID string) *CreateProfileTemplateOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateProfileTemplateOptions) SetName(name string) *CreateProfileTemplateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateProfileTemplateOptions) SetDescription(description string) *CreateProfileTemplateOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *CreateProfileTemplateOptions) SetProfile(profile *TemplateProfileComponentRequest) *CreateProfileTemplateOptions {
	_options.Profile = profile
	return _options
}

// SetPolicyTemplateReferences : Allow user to set PolicyTemplateReferences
func (_options *CreateProfileTemplateOptions) SetPolicyTemplateReferences(policyTemplateReferences []PolicyTemplateReference) *CreateProfileTemplateOptions {
	_options.PolicyTemplateReferences = policyTemplateReferences
	return _options
}

// SetActionControls : Allow user to set ActionControls
func (_options *CreateProfileTemplateOptions) SetActionControls(actionControls *ActionControls) *CreateProfileTemplateOptions {
	_options.ActionControls = actionControls
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProfileTemplateOptions) SetHeaders(param map[string]string) *CreateProfileTemplateOptions {
	options.Headers = param
	return options
}

// CreateProfileTemplateVersionOptions : The CreateProfileTemplateVersion options.
type CreateProfileTemplateVersionOptions struct {
	// ID of the trusted profile template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// ID of the account where the template resides.
	AccountID *string `json:"account_id,omitempty"`

	// The name of the trusted profile template. This is visible only in the enterprise account. Required field when
	// creating a new template. Otherwise this field is optional. If the field is included it will change the name value
	// for all existing versions of the template.
	Name *string `json:"name,omitempty"`

	// The description of the trusted profile template. Describe the template for enterprise account users.
	Description *string `json:"description,omitempty"`

	// Input body parameters for the TemplateProfileComponent.
	Profile *TemplateProfileComponentRequest `json:"profile,omitempty"`

	// Existing policy templates that you can reference to assign access in the trusted profile component.
	PolicyTemplateReferences []PolicyTemplateReference `json:"policy_template_references,omitempty"`

	ActionControls *ActionControls `json:"action_controls,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateProfileTemplateVersionOptions : Instantiate CreateProfileTemplateVersionOptions
func (*IamIdentityV1) NewCreateProfileTemplateVersionOptions(templateID string) *CreateProfileTemplateVersionOptions {
	return &CreateProfileTemplateVersionOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *CreateProfileTemplateVersionOptions) SetTemplateID(templateID string) *CreateProfileTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateProfileTemplateVersionOptions) SetAccountID(accountID string) *CreateProfileTemplateVersionOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateProfileTemplateVersionOptions) SetName(name string) *CreateProfileTemplateVersionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateProfileTemplateVersionOptions) SetDescription(description string) *CreateProfileTemplateVersionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *CreateProfileTemplateVersionOptions) SetProfile(profile *TemplateProfileComponentRequest) *CreateProfileTemplateVersionOptions {
	_options.Profile = profile
	return _options
}

// SetPolicyTemplateReferences : Allow user to set PolicyTemplateReferences
func (_options *CreateProfileTemplateVersionOptions) SetPolicyTemplateReferences(policyTemplateReferences []PolicyTemplateReference) *CreateProfileTemplateVersionOptions {
	_options.PolicyTemplateReferences = policyTemplateReferences
	return _options
}

// SetActionControls : Allow user to set ActionControls
func (_options *CreateProfileTemplateVersionOptions) SetActionControls(actionControls *ActionControls) *CreateProfileTemplateVersionOptions {
	_options.ActionControls = actionControls
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProfileTemplateVersionOptions) SetHeaders(param map[string]string) *CreateProfileTemplateVersionOptions {
	options.Headers = param
	return options
}

// CreateReportOptions : The CreateReport options.
type CreateReportOptions struct {
	// ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Optional report type. The supported value is 'inactive'. List all identities that have not authenticated within the
	// time indicated by duration.
	Type *string `json:"type,omitempty"`

	// Optional duration of the report. The supported unit of duration is hours.
	Duration *string `json:"duration,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateReportOptions : Instantiate CreateReportOptions
func (*IamIdentityV1) NewCreateReportOptions(accountID string) *CreateReportOptions {
	return &CreateReportOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateReportOptions) SetAccountID(accountID string) *CreateReportOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateReportOptions) SetType(typeVar string) *CreateReportOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetDuration : Allow user to set Duration
func (_options *CreateReportOptions) SetDuration(duration string) *CreateReportOptions {
	_options.Duration = core.StringPtr(duration)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateReportOptions) SetHeaders(param map[string]string) *CreateReportOptions {
	options.Headers = param
	return options
}

// CreateServiceIDOptions : The CreateServiceID options.
type CreateServiceIDOptions struct {
	// ID of the account the service ID belongs to.
	AccountID *string `json:"account_id" validate:"required"`

	// Name of the Service Id. The name is not checked for uniqueness. Therefore multiple names with the same value can
	// exist. Access is done via the UUID of the Service Id.
	Name *string `json:"name" validate:"required"`

	// The optional description of the Service Id. The 'description' property is only available if a description was
	// provided during a create of a Service Id.
	Description *string `json:"description,omitempty"`

	// Optional list of CRNs (string array) which point to the services connected to the service ID.
	UniqueInstanceCrns []string `json:"unique_instance_crns,omitempty"`

	// Parameters for the API key in the Create service Id V1 REST request.
	Apikey *APIKeyInsideCreateServiceIDRequest `json:"apikey,omitempty"`

	// Indicates if the service ID is locked for further write operations. False by default.
	EntityLock *string `json:"Entity-Lock,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateServiceIDOptions : Instantiate CreateServiceIDOptions
func (*IamIdentityV1) NewCreateServiceIDOptions(accountID string, name string) *CreateServiceIDOptions {
	return &CreateServiceIDOptions{
		AccountID: core.StringPtr(accountID),
		Name:      core.StringPtr(name),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateServiceIDOptions) SetAccountID(accountID string) *CreateServiceIDOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateServiceIDOptions) SetName(name string) *CreateServiceIDOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateServiceIDOptions) SetDescription(description string) *CreateServiceIDOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetUniqueInstanceCrns : Allow user to set UniqueInstanceCrns
func (_options *CreateServiceIDOptions) SetUniqueInstanceCrns(uniqueInstanceCrns []string) *CreateServiceIDOptions {
	_options.UniqueInstanceCrns = uniqueInstanceCrns
	return _options
}

// SetApikey : Allow user to set Apikey
func (_options *CreateServiceIDOptions) SetApikey(apikey *APIKeyInsideCreateServiceIDRequest) *CreateServiceIDOptions {
	_options.Apikey = apikey
	return _options
}

// SetEntityLock : Allow user to set EntityLock
func (_options *CreateServiceIDOptions) SetEntityLock(entityLock string) *CreateServiceIDOptions {
	_options.EntityLock = core.StringPtr(entityLock)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateServiceIDOptions) SetHeaders(param map[string]string) *CreateServiceIDOptions {
	options.Headers = param
	return options
}

// CreateTrustedProfileAssignmentOptions : The CreateTrustedProfileAssignment options.
type CreateTrustedProfileAssignmentOptions struct {
	// ID of the template to assign.
	TemplateID *string `json:"template_id" validate:"required"`

	// Version of the template to assign.
	TemplateVersion *int64 `json:"template_version" validate:"required"`

	// Type of target to deploy to.
	TargetType *string `json:"target_type" validate:"required"`

	// Identifier of target to deploy to.
	Target *string `json:"target" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateTrustedProfileAssignmentOptions.TargetType property.
// Type of target to deploy to.
const (
	CreateTrustedProfileAssignmentOptionsTargetTypeAccountConst      = "Account"
	CreateTrustedProfileAssignmentOptionsTargetTypeAccountgroupConst = "AccountGroup"
)

// NewCreateTrustedProfileAssignmentOptions : Instantiate CreateTrustedProfileAssignmentOptions
func (*IamIdentityV1) NewCreateTrustedProfileAssignmentOptions(templateID string, templateVersion int64, targetType string, target string) *CreateTrustedProfileAssignmentOptions {
	return &CreateTrustedProfileAssignmentOptions{
		TemplateID:      core.StringPtr(templateID),
		TemplateVersion: core.Int64Ptr(templateVersion),
		TargetType:      core.StringPtr(targetType),
		Target:          core.StringPtr(target),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *CreateTrustedProfileAssignmentOptions) SetTemplateID(templateID string) *CreateTrustedProfileAssignmentOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetTemplateVersion : Allow user to set TemplateVersion
func (_options *CreateTrustedProfileAssignmentOptions) SetTemplateVersion(templateVersion int64) *CreateTrustedProfileAssignmentOptions {
	_options.TemplateVersion = core.Int64Ptr(templateVersion)
	return _options
}

// SetTargetType : Allow user to set TargetType
func (_options *CreateTrustedProfileAssignmentOptions) SetTargetType(targetType string) *CreateTrustedProfileAssignmentOptions {
	_options.TargetType = core.StringPtr(targetType)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *CreateTrustedProfileAssignmentOptions) SetTarget(target string) *CreateTrustedProfileAssignmentOptions {
	_options.Target = core.StringPtr(target)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTrustedProfileAssignmentOptions) SetHeaders(param map[string]string) *CreateTrustedProfileAssignmentOptions {
	options.Headers = param
	return options
}

// DeleteAccountSettingsAssignmentOptions : The DeleteAccountSettingsAssignment options.
type DeleteAccountSettingsAssignmentOptions struct {
	// ID of the Assignment Record.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteAccountSettingsAssignmentOptions : Instantiate DeleteAccountSettingsAssignmentOptions
func (*IamIdentityV1) NewDeleteAccountSettingsAssignmentOptions(assignmentID string) *DeleteAccountSettingsAssignmentOptions {
	return &DeleteAccountSettingsAssignmentOptions{
		AssignmentID: core.StringPtr(assignmentID),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *DeleteAccountSettingsAssignmentOptions) SetAssignmentID(assignmentID string) *DeleteAccountSettingsAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAccountSettingsAssignmentOptions) SetHeaders(param map[string]string) *DeleteAccountSettingsAssignmentOptions {
	options.Headers = param
	return options
}

// DeleteAccountSettingsTemplateVersionOptions : The DeleteAccountSettingsTemplateVersion options.
type DeleteAccountSettingsTemplateVersionOptions struct {
	// ID of the account settings template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Version of the account settings template.
	Version *string `json:"version" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteAccountSettingsTemplateVersionOptions : Instantiate DeleteAccountSettingsTemplateVersionOptions
func (*IamIdentityV1) NewDeleteAccountSettingsTemplateVersionOptions(templateID string, version string) *DeleteAccountSettingsTemplateVersionOptions {
	return &DeleteAccountSettingsTemplateVersionOptions{
		TemplateID: core.StringPtr(templateID),
		Version:    core.StringPtr(version),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *DeleteAccountSettingsTemplateVersionOptions) SetTemplateID(templateID string) *DeleteAccountSettingsTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *DeleteAccountSettingsTemplateVersionOptions) SetVersion(version string) *DeleteAccountSettingsTemplateVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAccountSettingsTemplateVersionOptions) SetHeaders(param map[string]string) *DeleteAccountSettingsTemplateVersionOptions {
	options.Headers = param
	return options
}

// DeleteAllVersionsOfAccountSettingsTemplateOptions : The DeleteAllVersionsOfAccountSettingsTemplate options.
type DeleteAllVersionsOfAccountSettingsTemplateOptions struct {
	// ID of the account settings template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteAllVersionsOfAccountSettingsTemplateOptions : Instantiate DeleteAllVersionsOfAccountSettingsTemplateOptions
func (*IamIdentityV1) NewDeleteAllVersionsOfAccountSettingsTemplateOptions(templateID string) *DeleteAllVersionsOfAccountSettingsTemplateOptions {
	return &DeleteAllVersionsOfAccountSettingsTemplateOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *DeleteAllVersionsOfAccountSettingsTemplateOptions) SetTemplateID(templateID string) *DeleteAllVersionsOfAccountSettingsTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAllVersionsOfAccountSettingsTemplateOptions) SetHeaders(param map[string]string) *DeleteAllVersionsOfAccountSettingsTemplateOptions {
	options.Headers = param
	return options
}

// DeleteAllVersionsOfProfileTemplateOptions : The DeleteAllVersionsOfProfileTemplate options.
type DeleteAllVersionsOfProfileTemplateOptions struct {
	// ID of the trusted profile template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteAllVersionsOfProfileTemplateOptions : Instantiate DeleteAllVersionsOfProfileTemplateOptions
func (*IamIdentityV1) NewDeleteAllVersionsOfProfileTemplateOptions(templateID string) *DeleteAllVersionsOfProfileTemplateOptions {
	return &DeleteAllVersionsOfProfileTemplateOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *DeleteAllVersionsOfProfileTemplateOptions) SetTemplateID(templateID string) *DeleteAllVersionsOfProfileTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAllVersionsOfProfileTemplateOptions) SetHeaders(param map[string]string) *DeleteAllVersionsOfProfileTemplateOptions {
	options.Headers = param
	return options
}

// DeleteAPIKeyOptions : The DeleteAPIKey options.
type DeleteAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteAPIKeyOptions : Instantiate DeleteAPIKeyOptions
func (*IamIdentityV1) NewDeleteAPIKeyOptions(id string) *DeleteAPIKeyOptions {
	return &DeleteAPIKeyOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteAPIKeyOptions) SetID(id string) *DeleteAPIKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAPIKeyOptions) SetHeaders(param map[string]string) *DeleteAPIKeyOptions {
	options.Headers = param
	return options
}

// DeleteClaimRuleOptions : The DeleteClaimRule options.
type DeleteClaimRuleOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// ID of the claim rule to delete.
	RuleID *string `json:"rule-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteClaimRuleOptions : Instantiate DeleteClaimRuleOptions
func (*IamIdentityV1) NewDeleteClaimRuleOptions(profileID string, ruleID string) *DeleteClaimRuleOptions {
	return &DeleteClaimRuleOptions{
		ProfileID: core.StringPtr(profileID),
		RuleID:    core.StringPtr(ruleID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *DeleteClaimRuleOptions) SetProfileID(profileID string) *DeleteClaimRuleOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteClaimRuleOptions) SetRuleID(ruleID string) *DeleteClaimRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteClaimRuleOptions) SetHeaders(param map[string]string) *DeleteClaimRuleOptions {
	options.Headers = param
	return options
}

// DeleteLinkOptions : The DeleteLink options.
type DeleteLinkOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// ID of the link.
	LinkID *string `json:"link-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteLinkOptions : Instantiate DeleteLinkOptions
func (*IamIdentityV1) NewDeleteLinkOptions(profileID string, linkID string) *DeleteLinkOptions {
	return &DeleteLinkOptions{
		ProfileID: core.StringPtr(profileID),
		LinkID:    core.StringPtr(linkID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *DeleteLinkOptions) SetProfileID(profileID string) *DeleteLinkOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetLinkID : Allow user to set LinkID
func (_options *DeleteLinkOptions) SetLinkID(linkID string) *DeleteLinkOptions {
	_options.LinkID = core.StringPtr(linkID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLinkOptions) SetHeaders(param map[string]string) *DeleteLinkOptions {
	options.Headers = param
	return options
}

// DeleteProfileIdentityOptions : The DeleteProfileIdentity options.
type DeleteProfileIdentityOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Type of the identity.
	IdentityType *string `json:"identity-type" validate:"required,ne="`

	// Identifier of the identity that can assume the trusted profiles.
	IdentifierID *string `json:"identifier-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the DeleteProfileIdentityOptions.IdentityType property.
// Type of the identity.
const (
	DeleteProfileIdentityOptionsIdentityTypeCRNConst       = "crn"
	DeleteProfileIdentityOptionsIdentityTypeServiceidConst = "serviceid"
	DeleteProfileIdentityOptionsIdentityTypeUserConst      = "user"
)

// NewDeleteProfileIdentityOptions : Instantiate DeleteProfileIdentityOptions
func (*IamIdentityV1) NewDeleteProfileIdentityOptions(profileID string, identityType string, identifierID string) *DeleteProfileIdentityOptions {
	return &DeleteProfileIdentityOptions{
		ProfileID:    core.StringPtr(profileID),
		IdentityType: core.StringPtr(identityType),
		IdentifierID: core.StringPtr(identifierID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *DeleteProfileIdentityOptions) SetProfileID(profileID string) *DeleteProfileIdentityOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetIdentityType : Allow user to set IdentityType
func (_options *DeleteProfileIdentityOptions) SetIdentityType(identityType string) *DeleteProfileIdentityOptions {
	_options.IdentityType = core.StringPtr(identityType)
	return _options
}

// SetIdentifierID : Allow user to set IdentifierID
func (_options *DeleteProfileIdentityOptions) SetIdentifierID(identifierID string) *DeleteProfileIdentityOptions {
	_options.IdentifierID = core.StringPtr(identifierID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProfileIdentityOptions) SetHeaders(param map[string]string) *DeleteProfileIdentityOptions {
	options.Headers = param
	return options
}

// DeleteProfileOptions : The DeleteProfile options.
type DeleteProfileOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteProfileOptions : Instantiate DeleteProfileOptions
func (*IamIdentityV1) NewDeleteProfileOptions(profileID string) *DeleteProfileOptions {
	return &DeleteProfileOptions{
		ProfileID: core.StringPtr(profileID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *DeleteProfileOptions) SetProfileID(profileID string) *DeleteProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProfileOptions) SetHeaders(param map[string]string) *DeleteProfileOptions {
	options.Headers = param
	return options
}

// DeleteProfileTemplateVersionOptions : The DeleteProfileTemplateVersion options.
type DeleteProfileTemplateVersionOptions struct {
	// ID of the trusted profile template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Version of the Profile Template.
	Version *string `json:"version" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteProfileTemplateVersionOptions : Instantiate DeleteProfileTemplateVersionOptions
func (*IamIdentityV1) NewDeleteProfileTemplateVersionOptions(templateID string, version string) *DeleteProfileTemplateVersionOptions {
	return &DeleteProfileTemplateVersionOptions{
		TemplateID: core.StringPtr(templateID),
		Version:    core.StringPtr(version),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *DeleteProfileTemplateVersionOptions) SetTemplateID(templateID string) *DeleteProfileTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *DeleteProfileTemplateVersionOptions) SetVersion(version string) *DeleteProfileTemplateVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProfileTemplateVersionOptions) SetHeaders(param map[string]string) *DeleteProfileTemplateVersionOptions {
	options.Headers = param
	return options
}

// DeleteServiceIDOptions : The DeleteServiceID options.
type DeleteServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteServiceIDOptions : Instantiate DeleteServiceIDOptions
func (*IamIdentityV1) NewDeleteServiceIDOptions(id string) *DeleteServiceIDOptions {
	return &DeleteServiceIDOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteServiceIDOptions) SetID(id string) *DeleteServiceIDOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteServiceIDOptions) SetHeaders(param map[string]string) *DeleteServiceIDOptions {
	options.Headers = param
	return options
}

// DeleteTrustedProfileAssignmentOptions : The DeleteTrustedProfileAssignment options.
type DeleteTrustedProfileAssignmentOptions struct {
	// ID of the Assignment Record.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTrustedProfileAssignmentOptions : Instantiate DeleteTrustedProfileAssignmentOptions
func (*IamIdentityV1) NewDeleteTrustedProfileAssignmentOptions(assignmentID string) *DeleteTrustedProfileAssignmentOptions {
	return &DeleteTrustedProfileAssignmentOptions{
		AssignmentID: core.StringPtr(assignmentID),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *DeleteTrustedProfileAssignmentOptions) SetAssignmentID(assignmentID string) *DeleteTrustedProfileAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTrustedProfileAssignmentOptions) SetHeaders(param map[string]string) *DeleteTrustedProfileAssignmentOptions {
	options.Headers = param
	return options
}

// DisableAPIKeyOptions : The DisableAPIKey options.
type DisableAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDisableAPIKeyOptions : Instantiate DisableAPIKeyOptions
func (*IamIdentityV1) NewDisableAPIKeyOptions(id string) *DisableAPIKeyOptions {
	return &DisableAPIKeyOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DisableAPIKeyOptions) SetID(id string) *DisableAPIKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DisableAPIKeyOptions) SetHeaders(param map[string]string) *DisableAPIKeyOptions {
	options.Headers = param
	return options
}

// EffectiveAccountSettingsResponse : Response body format for Account Settings REST requests.
type EffectiveAccountSettingsResponse struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// Unique ID of the account.
	AccountID *string `json:"account_id" validate:"required"`

	Effective *AccountSettingsEffectiveSection `json:"effective" validate:"required"`

	Account *AccountSettingsAccountSection `json:"account" validate:"required"`

	// assigned template section.
	AssignedTemplates []AccountSettingsAssignedTemplatesSection `json:"assigned_templates,omitempty"`
}

// UnmarshalEffectiveAccountSettingsResponse unmarshals an instance of EffectiveAccountSettingsResponse from the specified map of raw messages.
func UnmarshalEffectiveAccountSettingsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EffectiveAccountSettingsResponse)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "effective", &obj.Effective, UnmarshalAccountSettingsEffectiveSection)
	if err != nil {
		err = core.SDKErrorf(err, "", "effective-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "account", &obj.Account, UnmarshalAccountSettingsAccountSection)
	if err != nil {
		err = core.SDKErrorf(err, "", "account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "assigned_templates", &obj.AssignedTemplates, UnmarshalAccountSettingsAssignedTemplatesSection)
	if err != nil {
		err = core.SDKErrorf(err, "", "assigned_templates-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EffectiveAccountSettingsUserMfa : EffectiveAccountSettingsUserMfa struct
type EffectiveAccountSettingsUserMfa struct {
	// The iam_id of the user.
	IamID *string `json:"iam_id" validate:"required"`

	// Defines the MFA requirement for the user. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa" validate:"required"`

	// name of the user account.
	Name *string `json:"name,omitempty"`

	// userName of the user.
	UserName *string `json:"userName,omitempty"`

	// email of the user.
	Email *string `json:"email,omitempty"`

	// optional description.
	Description *string `json:"description,omitempty"`
}

// Constants associated with the EffectiveAccountSettingsUserMfa.Mfa property.
// Defines the MFA requirement for the user. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	EffectiveAccountSettingsUserMfaMfaLevel1Const     = "LEVEL1"
	EffectiveAccountSettingsUserMfaMfaLevel2Const     = "LEVEL2"
	EffectiveAccountSettingsUserMfaMfaLevel3Const     = "LEVEL3"
	EffectiveAccountSettingsUserMfaMfaNoneConst       = "NONE"
	EffectiveAccountSettingsUserMfaMfaNoneNoRopcConst = "NONE_NO_ROPC"
	EffectiveAccountSettingsUserMfaMfaTotpConst       = "TOTP"
	EffectiveAccountSettingsUserMfaMfaTotp4allConst   = "TOTP4ALL"
)

// UnmarshalEffectiveAccountSettingsUserMfa unmarshals an instance of EffectiveAccountSettingsUserMfa from the specified map of raw messages.
func UnmarshalEffectiveAccountSettingsUserMfa(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EffectiveAccountSettingsUserMfa)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "userName", &obj.UserName)
	if err != nil {
		err = core.SDKErrorf(err, "", "userName-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnableAPIKeyOptions : The EnableAPIKey options.
type EnableAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewEnableAPIKeyOptions : Instantiate EnableAPIKeyOptions
func (*IamIdentityV1) NewEnableAPIKeyOptions(id string) *EnableAPIKeyOptions {
	return &EnableAPIKeyOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *EnableAPIKeyOptions) SetID(id string) *EnableAPIKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *EnableAPIKeyOptions) SetHeaders(param map[string]string) *EnableAPIKeyOptions {
	options.Headers = param
	return options
}

// EnityHistoryRecord : Response body format for an entity history record.
type EnityHistoryRecord struct {
	// Timestamp when the action was triggered.
	Timestamp *string `json:"timestamp" validate:"required"`

	// IAM ID of the identity which triggered the action.
	IamID *string `json:"iam_id" validate:"required"`

	// Account of the identity which triggered the action.
	IamIDAccount *string `json:"iam_id_account" validate:"required"`

	// Action of the history entry.
	Action *string `json:"action" validate:"required"`

	// Params of the history entry.
	Params []string `json:"params" validate:"required"`

	// Message which summarizes the executed action.
	Message *string `json:"message" validate:"required"`
}

// UnmarshalEnityHistoryRecord unmarshals an instance of EnityHistoryRecord from the specified map of raw messages.
func UnmarshalEnityHistoryRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnityHistoryRecord)
	err = core.UnmarshalPrimitive(m, "timestamp", &obj.Timestamp)
	if err != nil {
		err = core.SDKErrorf(err, "", "timestamp-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id_account", &obj.IamIDAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		err = core.SDKErrorf(err, "", "action-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "params", &obj.Params)
	if err != nil {
		err = core.SDKErrorf(err, "", "params-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EntityActivity : EntityActivity struct
type EntityActivity struct {
	// Unique id of the entity.
	ID *string `json:"id" validate:"required"`

	// Name provided during creation of the entity.
	Name *string `json:"name,omitempty"`

	// Time when the entity was last authenticated.
	LastAuthn *string `json:"last_authn,omitempty"`
}

// UnmarshalEntityActivity unmarshals an instance of EntityActivity from the specified map of raw messages.
func UnmarshalEntityActivity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EntityActivity)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_authn", &obj.LastAuthn)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_authn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Error : Error information.
type Error struct {
	// Error code of the REST Exception.
	Code *string `json:"code" validate:"required"`

	// Error message code of the REST Exception.
	MessageCode *string `json:"message_code" validate:"required"`

	// Error message of the REST Exception. Error messages are derived base on the input locale of the REST request and the
	// available Message catalogs. Dynamic fallback to 'us-english' is happening if no message catalog is available for the
	// provided input locale.
	Message *string `json:"message" validate:"required"`

	// Error details of the REST Exception.
	Details *string `json:"details,omitempty"`
}

// UnmarshalError unmarshals an instance of Error from the specified map of raw messages.
func UnmarshalError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Error)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		err = core.SDKErrorf(err, "", "code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message_code", &obj.MessageCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "message_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "details", &obj.Details)
	if err != nil {
		err = core.SDKErrorf(err, "", "details-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExceptionResponse : Response body parameters in case of error situations.
type ExceptionResponse struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// Error message code of the REST Exception.
	StatusCode *string `json:"status_code" validate:"required"`

	// List of errors that occured.
	Errors []Error `json:"errors" validate:"required"`

	// Unique ID of the requst.
	Trace *string `json:"trace,omitempty"`
}

// UnmarshalExceptionResponse unmarshals an instance of ExceptionResponse from the specified map of raw messages.
func UnmarshalExceptionResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExceptionResponse)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalError)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		err = core.SDKErrorf(err, "", "trace-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAccountSettingsAssignmentOptions : The GetAccountSettingsAssignment options.
type GetAccountSettingsAssignmentOptions struct {
	// ID of the Assignment Record.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAccountSettingsAssignmentOptions : Instantiate GetAccountSettingsAssignmentOptions
func (*IamIdentityV1) NewGetAccountSettingsAssignmentOptions(assignmentID string) *GetAccountSettingsAssignmentOptions {
	return &GetAccountSettingsAssignmentOptions{
		AssignmentID: core.StringPtr(assignmentID),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *GetAccountSettingsAssignmentOptions) SetAssignmentID(assignmentID string) *GetAccountSettingsAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetAccountSettingsAssignmentOptions) SetIncludeHistory(includeHistory bool) *GetAccountSettingsAssignmentOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountSettingsAssignmentOptions) SetHeaders(param map[string]string) *GetAccountSettingsAssignmentOptions {
	options.Headers = param
	return options
}

// GetAccountSettingsOptions : The GetAccountSettings options.
type GetAccountSettingsOptions struct {
	// Unique ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAccountSettingsOptions : Instantiate GetAccountSettingsOptions
func (*IamIdentityV1) NewGetAccountSettingsOptions(accountID string) *GetAccountSettingsOptions {
	return &GetAccountSettingsOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetAccountSettingsOptions) SetAccountID(accountID string) *GetAccountSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetAccountSettingsOptions) SetIncludeHistory(includeHistory bool) *GetAccountSettingsOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountSettingsOptions) SetHeaders(param map[string]string) *GetAccountSettingsOptions {
	options.Headers = param
	return options
}

// GetAccountSettingsTemplateVersionOptions : The GetAccountSettingsTemplateVersion options.
type GetAccountSettingsTemplateVersionOptions struct {
	// ID of the account settings template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Version of the account settings template.
	Version *string `json:"version" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAccountSettingsTemplateVersionOptions : Instantiate GetAccountSettingsTemplateVersionOptions
func (*IamIdentityV1) NewGetAccountSettingsTemplateVersionOptions(templateID string, version string) *GetAccountSettingsTemplateVersionOptions {
	return &GetAccountSettingsTemplateVersionOptions{
		TemplateID: core.StringPtr(templateID),
		Version:    core.StringPtr(version),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *GetAccountSettingsTemplateVersionOptions) SetTemplateID(templateID string) *GetAccountSettingsTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *GetAccountSettingsTemplateVersionOptions) SetVersion(version string) *GetAccountSettingsTemplateVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetAccountSettingsTemplateVersionOptions) SetIncludeHistory(includeHistory bool) *GetAccountSettingsTemplateVersionOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountSettingsTemplateVersionOptions) SetHeaders(param map[string]string) *GetAccountSettingsTemplateVersionOptions {
	options.Headers = param
	return options
}

// GetAPIKeyOptions : The GetAPIKey options.
type GetAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Defines if the entity's activity is included in the response. Retrieving activity data is an expensive operation, so
	// only request this when needed.
	IncludeActivity *bool `json:"include_activity,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAPIKeyOptions : Instantiate GetAPIKeyOptions
func (*IamIdentityV1) NewGetAPIKeyOptions(id string) *GetAPIKeyOptions {
	return &GetAPIKeyOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetAPIKeyOptions) SetID(id string) *GetAPIKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetAPIKeyOptions) SetIncludeHistory(includeHistory bool) *GetAPIKeyOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetIncludeActivity : Allow user to set IncludeActivity
func (_options *GetAPIKeyOptions) SetIncludeActivity(includeActivity bool) *GetAPIKeyOptions {
	_options.IncludeActivity = core.BoolPtr(includeActivity)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAPIKeyOptions) SetHeaders(param map[string]string) *GetAPIKeyOptions {
	options.Headers = param
	return options
}

// GetAPIKeysDetailsOptions : The GetAPIKeysDetails options.
type GetAPIKeysDetailsOptions struct {
	// API key value.
	IamAPIKey *string `json:"IAM-ApiKey,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAPIKeysDetailsOptions : Instantiate GetAPIKeysDetailsOptions
func (*IamIdentityV1) NewGetAPIKeysDetailsOptions() *GetAPIKeysDetailsOptions {
	return &GetAPIKeysDetailsOptions{}
}

// SetIamAPIKey : Allow user to set IamAPIKey
func (_options *GetAPIKeysDetailsOptions) SetIamAPIKey(iamAPIKey string) *GetAPIKeysDetailsOptions {
	_options.IamAPIKey = core.StringPtr(iamAPIKey)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetAPIKeysDetailsOptions) SetIncludeHistory(includeHistory bool) *GetAPIKeysDetailsOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAPIKeysDetailsOptions) SetHeaders(param map[string]string) *GetAPIKeysDetailsOptions {
	options.Headers = param
	return options
}

// GetClaimRuleOptions : The GetClaimRule options.
type GetClaimRuleOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// ID of the claim rule to get.
	RuleID *string `json:"rule-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetClaimRuleOptions : Instantiate GetClaimRuleOptions
func (*IamIdentityV1) NewGetClaimRuleOptions(profileID string, ruleID string) *GetClaimRuleOptions {
	return &GetClaimRuleOptions{
		ProfileID: core.StringPtr(profileID),
		RuleID:    core.StringPtr(ruleID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetClaimRuleOptions) SetProfileID(profileID string) *GetClaimRuleOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *GetClaimRuleOptions) SetRuleID(ruleID string) *GetClaimRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetClaimRuleOptions) SetHeaders(param map[string]string) *GetClaimRuleOptions {
	options.Headers = param
	return options
}

// GetEffectiveAccountSettingsOptions : The GetEffectiveAccountSettings options.
type GetEffectiveAccountSettingsOptions struct {
	// Unique ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Enrich MFA exemptions with user information.
	ResolveUserMfa *bool `json:"resolve_user_mfa,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetEffectiveAccountSettingsOptions : Instantiate GetEffectiveAccountSettingsOptions
func (*IamIdentityV1) NewGetEffectiveAccountSettingsOptions(accountID string) *GetEffectiveAccountSettingsOptions {
	return &GetEffectiveAccountSettingsOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetEffectiveAccountSettingsOptions) SetAccountID(accountID string) *GetEffectiveAccountSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetEffectiveAccountSettingsOptions) SetIncludeHistory(includeHistory bool) *GetEffectiveAccountSettingsOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetResolveUserMfa : Allow user to set ResolveUserMfa
func (_options *GetEffectiveAccountSettingsOptions) SetResolveUserMfa(resolveUserMfa bool) *GetEffectiveAccountSettingsOptions {
	_options.ResolveUserMfa = core.BoolPtr(resolveUserMfa)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetEffectiveAccountSettingsOptions) SetHeaders(param map[string]string) *GetEffectiveAccountSettingsOptions {
	options.Headers = param
	return options
}

// GetLatestAccountSettingsTemplateVersionOptions : The GetLatestAccountSettingsTemplateVersion options.
type GetLatestAccountSettingsTemplateVersionOptions struct {
	// ID of the account settings template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLatestAccountSettingsTemplateVersionOptions : Instantiate GetLatestAccountSettingsTemplateVersionOptions
func (*IamIdentityV1) NewGetLatestAccountSettingsTemplateVersionOptions(templateID string) *GetLatestAccountSettingsTemplateVersionOptions {
	return &GetLatestAccountSettingsTemplateVersionOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *GetLatestAccountSettingsTemplateVersionOptions) SetTemplateID(templateID string) *GetLatestAccountSettingsTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetLatestAccountSettingsTemplateVersionOptions) SetIncludeHistory(includeHistory bool) *GetLatestAccountSettingsTemplateVersionOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLatestAccountSettingsTemplateVersionOptions) SetHeaders(param map[string]string) *GetLatestAccountSettingsTemplateVersionOptions {
	options.Headers = param
	return options
}

// GetLatestProfileTemplateVersionOptions : The GetLatestProfileTemplateVersion options.
type GetLatestProfileTemplateVersionOptions struct {
	// ID of the trusted profile template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLatestProfileTemplateVersionOptions : Instantiate GetLatestProfileTemplateVersionOptions
func (*IamIdentityV1) NewGetLatestProfileTemplateVersionOptions(templateID string) *GetLatestProfileTemplateVersionOptions {
	return &GetLatestProfileTemplateVersionOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *GetLatestProfileTemplateVersionOptions) SetTemplateID(templateID string) *GetLatestProfileTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetLatestProfileTemplateVersionOptions) SetIncludeHistory(includeHistory bool) *GetLatestProfileTemplateVersionOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLatestProfileTemplateVersionOptions) SetHeaders(param map[string]string) *GetLatestProfileTemplateVersionOptions {
	options.Headers = param
	return options
}

// GetLinkOptions : The GetLink options.
type GetLinkOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// ID of the link.
	LinkID *string `json:"link-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLinkOptions : Instantiate GetLinkOptions
func (*IamIdentityV1) NewGetLinkOptions(profileID string, linkID string) *GetLinkOptions {
	return &GetLinkOptions{
		ProfileID: core.StringPtr(profileID),
		LinkID:    core.StringPtr(linkID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetLinkOptions) SetProfileID(profileID string) *GetLinkOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetLinkID : Allow user to set LinkID
func (_options *GetLinkOptions) SetLinkID(linkID string) *GetLinkOptions {
	_options.LinkID = core.StringPtr(linkID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLinkOptions) SetHeaders(param map[string]string) *GetLinkOptions {
	options.Headers = param
	return options
}

// GetMfaReportOptions : The GetMfaReport options.
type GetMfaReportOptions struct {
	// ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Reference for the report to be generated, You can use 'latest' to get the latest report for the given account.
	Reference *string `json:"reference" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetMfaReportOptions : Instantiate GetMfaReportOptions
func (*IamIdentityV1) NewGetMfaReportOptions(accountID string, reference string) *GetMfaReportOptions {
	return &GetMfaReportOptions{
		AccountID: core.StringPtr(accountID),
		Reference: core.StringPtr(reference),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetMfaReportOptions) SetAccountID(accountID string) *GetMfaReportOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetReference : Allow user to set Reference
func (_options *GetMfaReportOptions) SetReference(reference string) *GetMfaReportOptions {
	_options.Reference = core.StringPtr(reference)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetMfaReportOptions) SetHeaders(param map[string]string) *GetMfaReportOptions {
	options.Headers = param
	return options
}

// GetMfaStatusOptions : The GetMfaStatus options.
type GetMfaStatusOptions struct {
	// ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// iam_id of the user. This user must be the member of the account.
	IamID *string `json:"iam_id" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetMfaStatusOptions : Instantiate GetMfaStatusOptions
func (*IamIdentityV1) NewGetMfaStatusOptions(accountID string, iamID string) *GetMfaStatusOptions {
	return &GetMfaStatusOptions{
		AccountID: core.StringPtr(accountID),
		IamID:     core.StringPtr(iamID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetMfaStatusOptions) SetAccountID(accountID string) *GetMfaStatusOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetIamID : Allow user to set IamID
func (_options *GetMfaStatusOptions) SetIamID(iamID string) *GetMfaStatusOptions {
	_options.IamID = core.StringPtr(iamID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetMfaStatusOptions) SetHeaders(param map[string]string) *GetMfaStatusOptions {
	options.Headers = param
	return options
}

// GetProfileIdentitiesOptions : The GetProfileIdentities options.
type GetProfileIdentitiesOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetProfileIdentitiesOptions : Instantiate GetProfileIdentitiesOptions
func (*IamIdentityV1) NewGetProfileIdentitiesOptions(profileID string) *GetProfileIdentitiesOptions {
	return &GetProfileIdentitiesOptions{
		ProfileID: core.StringPtr(profileID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetProfileIdentitiesOptions) SetProfileID(profileID string) *GetProfileIdentitiesOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileIdentitiesOptions) SetHeaders(param map[string]string) *GetProfileIdentitiesOptions {
	options.Headers = param
	return options
}

// GetProfileIdentityOptions : The GetProfileIdentity options.
type GetProfileIdentityOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Type of the identity.
	IdentityType *string `json:"identity-type" validate:"required,ne="`

	// Identifier of the identity that can assume the trusted profiles.
	IdentifierID *string `json:"identifier-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the GetProfileIdentityOptions.IdentityType property.
// Type of the identity.
const (
	GetProfileIdentityOptionsIdentityTypeCRNConst       = "crn"
	GetProfileIdentityOptionsIdentityTypeServiceidConst = "serviceid"
	GetProfileIdentityOptionsIdentityTypeUserConst      = "user"
)

// NewGetProfileIdentityOptions : Instantiate GetProfileIdentityOptions
func (*IamIdentityV1) NewGetProfileIdentityOptions(profileID string, identityType string, identifierID string) *GetProfileIdentityOptions {
	return &GetProfileIdentityOptions{
		ProfileID:    core.StringPtr(profileID),
		IdentityType: core.StringPtr(identityType),
		IdentifierID: core.StringPtr(identifierID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetProfileIdentityOptions) SetProfileID(profileID string) *GetProfileIdentityOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetIdentityType : Allow user to set IdentityType
func (_options *GetProfileIdentityOptions) SetIdentityType(identityType string) *GetProfileIdentityOptions {
	_options.IdentityType = core.StringPtr(identityType)
	return _options
}

// SetIdentifierID : Allow user to set IdentifierID
func (_options *GetProfileIdentityOptions) SetIdentifierID(identifierID string) *GetProfileIdentityOptions {
	_options.IdentifierID = core.StringPtr(identifierID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileIdentityOptions) SetHeaders(param map[string]string) *GetProfileIdentityOptions {
	options.Headers = param
	return options
}

// GetProfileOptions : The GetProfile options.
type GetProfileOptions struct {
	// ID of the trusted profile to get.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Defines if the entity's activity is included in the response. Retrieving activity data is an expensive operation, so
	// only request this when needed.
	IncludeActivity *bool `json:"include_activity,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetProfileOptions : Instantiate GetProfileOptions
func (*IamIdentityV1) NewGetProfileOptions(profileID string) *GetProfileOptions {
	return &GetProfileOptions{
		ProfileID: core.StringPtr(profileID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetProfileOptions) SetProfileID(profileID string) *GetProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetIncludeActivity : Allow user to set IncludeActivity
func (_options *GetProfileOptions) SetIncludeActivity(includeActivity bool) *GetProfileOptions {
	_options.IncludeActivity = core.BoolPtr(includeActivity)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileOptions) SetHeaders(param map[string]string) *GetProfileOptions {
	options.Headers = param
	return options
}

// GetProfileTemplateVersionOptions : The GetProfileTemplateVersion options.
type GetProfileTemplateVersionOptions struct {
	// ID of the trusted profile template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Version of the Profile Template.
	Version *string `json:"version" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetProfileTemplateVersionOptions : Instantiate GetProfileTemplateVersionOptions
func (*IamIdentityV1) NewGetProfileTemplateVersionOptions(templateID string, version string) *GetProfileTemplateVersionOptions {
	return &GetProfileTemplateVersionOptions{
		TemplateID: core.StringPtr(templateID),
		Version:    core.StringPtr(version),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *GetProfileTemplateVersionOptions) SetTemplateID(templateID string) *GetProfileTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *GetProfileTemplateVersionOptions) SetVersion(version string) *GetProfileTemplateVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetProfileTemplateVersionOptions) SetIncludeHistory(includeHistory bool) *GetProfileTemplateVersionOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileTemplateVersionOptions) SetHeaders(param map[string]string) *GetProfileTemplateVersionOptions {
	options.Headers = param
	return options
}

// GetReportOptions : The GetReport options.
type GetReportOptions struct {
	// ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Reference for the report to be generated, You can use 'latest' to get the latest report for the given account.
	Reference *string `json:"reference" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetReportOptions : Instantiate GetReportOptions
func (*IamIdentityV1) NewGetReportOptions(accountID string, reference string) *GetReportOptions {
	return &GetReportOptions{
		AccountID: core.StringPtr(accountID),
		Reference: core.StringPtr(reference),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetReportOptions) SetAccountID(accountID string) *GetReportOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetReference : Allow user to set Reference
func (_options *GetReportOptions) SetReference(reference string) *GetReportOptions {
	_options.Reference = core.StringPtr(reference)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportOptions) SetHeaders(param map[string]string) *GetReportOptions {
	options.Headers = param
	return options
}

// GetServiceIDOptions : The GetServiceID options.
type GetServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `json:"id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Defines if the entity's activity is included in the response. Retrieving activity data is an expensive operation, so
	// only request this when needed.
	IncludeActivity *bool `json:"include_activity,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetServiceIDOptions : Instantiate GetServiceIDOptions
func (*IamIdentityV1) NewGetServiceIDOptions(id string) *GetServiceIDOptions {
	return &GetServiceIDOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetServiceIDOptions) SetID(id string) *GetServiceIDOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetServiceIDOptions) SetIncludeHistory(includeHistory bool) *GetServiceIDOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetIncludeActivity : Allow user to set IncludeActivity
func (_options *GetServiceIDOptions) SetIncludeActivity(includeActivity bool) *GetServiceIDOptions {
	_options.IncludeActivity = core.BoolPtr(includeActivity)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetServiceIDOptions) SetHeaders(param map[string]string) *GetServiceIDOptions {
	options.Headers = param
	return options
}

// GetTrustedProfileAssignmentOptions : The GetTrustedProfileAssignment options.
type GetTrustedProfileAssignmentOptions struct {
	// ID of the Assignment Record.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTrustedProfileAssignmentOptions : Instantiate GetTrustedProfileAssignmentOptions
func (*IamIdentityV1) NewGetTrustedProfileAssignmentOptions(assignmentID string) *GetTrustedProfileAssignmentOptions {
	return &GetTrustedProfileAssignmentOptions{
		AssignmentID: core.StringPtr(assignmentID),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *GetTrustedProfileAssignmentOptions) SetAssignmentID(assignmentID string) *GetTrustedProfileAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *GetTrustedProfileAssignmentOptions) SetIncludeHistory(includeHistory bool) *GetTrustedProfileAssignmentOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTrustedProfileAssignmentOptions) SetHeaders(param map[string]string) *GetTrustedProfileAssignmentOptions {
	options.Headers = param
	return options
}

// IDBasedMfaEnrollment : IDBasedMfaEnrollment struct
type IDBasedMfaEnrollment struct {
	// Defines the MFA trait for the account. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	TraitAccountDefault *string `json:"trait_account_default" validate:"required"`

	// Defines the MFA trait for the account. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	TraitUserSpecific *string `json:"trait_user_specific,omitempty"`

	// Defines the MFA trait for the account. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	TraitEffective *string `json:"trait_effective" validate:"required"`

	// The enrollment complies to the effective requirement.
	Complies *bool `json:"complies" validate:"required"`

	// Defines comply state for the account. Valid values:
	//   * NO - User does not comply in the given account.
	//   * ACCOUNT- User complies in the given account, but does not comply in at least one of the other account
	// memberships.
	//   * CROSS_ACCOUNT - User complies in the given account and across all other account memberships.
	ComplyState *string `json:"comply_state,omitempty"`
}

// Constants associated with the IDBasedMfaEnrollment.TraitAccountDefault property.
// Defines the MFA trait for the account. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	IDBasedMfaEnrollmentTraitAccountDefaultLevel1Const     = "LEVEL1"
	IDBasedMfaEnrollmentTraitAccountDefaultLevel2Const     = "LEVEL2"
	IDBasedMfaEnrollmentTraitAccountDefaultLevel3Const     = "LEVEL3"
	IDBasedMfaEnrollmentTraitAccountDefaultNoneConst       = "NONE"
	IDBasedMfaEnrollmentTraitAccountDefaultNoneNoRopcConst = "NONE_NO_ROPC"
	IDBasedMfaEnrollmentTraitAccountDefaultTotpConst       = "TOTP"
	IDBasedMfaEnrollmentTraitAccountDefaultTotp4allConst   = "TOTP4ALL"
)

// Constants associated with the IDBasedMfaEnrollment.TraitUserSpecific property.
// Defines the MFA trait for the account. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	IDBasedMfaEnrollmentTraitUserSpecificLevel1Const     = "LEVEL1"
	IDBasedMfaEnrollmentTraitUserSpecificLevel2Const     = "LEVEL2"
	IDBasedMfaEnrollmentTraitUserSpecificLevel3Const     = "LEVEL3"
	IDBasedMfaEnrollmentTraitUserSpecificNoneConst       = "NONE"
	IDBasedMfaEnrollmentTraitUserSpecificNoneNoRopcConst = "NONE_NO_ROPC"
	IDBasedMfaEnrollmentTraitUserSpecificTotpConst       = "TOTP"
	IDBasedMfaEnrollmentTraitUserSpecificTotp4allConst   = "TOTP4ALL"
)

// Constants associated with the IDBasedMfaEnrollment.TraitEffective property.
// Defines the MFA trait for the account. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	IDBasedMfaEnrollmentTraitEffectiveLevel1Const     = "LEVEL1"
	IDBasedMfaEnrollmentTraitEffectiveLevel2Const     = "LEVEL2"
	IDBasedMfaEnrollmentTraitEffectiveLevel3Const     = "LEVEL3"
	IDBasedMfaEnrollmentTraitEffectiveNoneConst       = "NONE"
	IDBasedMfaEnrollmentTraitEffectiveNoneNoRopcConst = "NONE_NO_ROPC"
	IDBasedMfaEnrollmentTraitEffectiveTotpConst       = "TOTP"
	IDBasedMfaEnrollmentTraitEffectiveTotp4allConst   = "TOTP4ALL"
)

// Constants associated with the IDBasedMfaEnrollment.ComplyState property.
// Defines comply state for the account. Valid values:
//   - NO - User does not comply in the given account.
//   - ACCOUNT- User complies in the given account, but does not comply in at least one of the other account
//
// memberships.
//   - CROSS_ACCOUNT - User complies in the given account and across all other account memberships.
const (
	IDBasedMfaEnrollmentComplyStateAccountConst      = "ACCOUNT"
	IDBasedMfaEnrollmentComplyStateCrossAccountConst = "CROSS_ACCOUNT"
	IDBasedMfaEnrollmentComplyStateNoConst           = "NO"
)

// UnmarshalIDBasedMfaEnrollment unmarshals an instance of IDBasedMfaEnrollment from the specified map of raw messages.
func UnmarshalIDBasedMfaEnrollment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IDBasedMfaEnrollment)
	err = core.UnmarshalPrimitive(m, "trait_account_default", &obj.TraitAccountDefault)
	if err != nil {
		err = core.SDKErrorf(err, "", "trait_account_default-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "trait_user_specific", &obj.TraitUserSpecific)
	if err != nil {
		err = core.SDKErrorf(err, "", "trait_user_specific-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "trait_effective", &obj.TraitEffective)
	if err != nil {
		err = core.SDKErrorf(err, "", "trait_effective-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "complies", &obj.Complies)
	if err != nil {
		err = core.SDKErrorf(err, "", "complies-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "comply_state", &obj.ComplyState)
	if err != nil {
		err = core.SDKErrorf(err, "", "comply_state-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAccountSettingsAssignmentsOptions : The ListAccountSettingsAssignments options.
type ListAccountSettingsAssignmentsOptions struct {
	// Account ID of the Assignments to query. This parameter is required unless using a pagetoken.
	AccountID *string `json:"account_id,omitempty"`

	// Filter results by Template Id.
	TemplateID *string `json:"template_id,omitempty"`

	// Filter results Template Version.
	TemplateVersion *string `json:"template_version,omitempty"`

	// Filter results by the assignment target.
	Target *string `json:"target,omitempty"`

	// Filter results by the assignment's target type.
	TargetType *string `json:"target_type,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Limit *int64 `json:"limit,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// If specified, the items are sorted by the value of this property.
	Sort *string `json:"sort,omitempty"`

	// Sort order.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListAccountSettingsAssignmentsOptions.TargetType property.
// Filter results by the assignment's target type.
const (
	ListAccountSettingsAssignmentsOptionsTargetTypeAccountConst      = "Account"
	ListAccountSettingsAssignmentsOptionsTargetTypeAccountgroupConst = "AccountGroup"
)

// Constants associated with the ListAccountSettingsAssignmentsOptions.Sort property.
// If specified, the items are sorted by the value of this property.
const (
	ListAccountSettingsAssignmentsOptionsSortCreatedAtConst      = "created_at"
	ListAccountSettingsAssignmentsOptionsSortLastModifiedAtConst = "last_modified_at"
	ListAccountSettingsAssignmentsOptionsSortTemplateIDConst     = "template_id"
)

// Constants associated with the ListAccountSettingsAssignmentsOptions.Order property.
// Sort order.
const (
	ListAccountSettingsAssignmentsOptionsOrderAscConst  = "asc"
	ListAccountSettingsAssignmentsOptionsOrderDescConst = "desc"
)

// NewListAccountSettingsAssignmentsOptions : Instantiate ListAccountSettingsAssignmentsOptions
func (*IamIdentityV1) NewListAccountSettingsAssignmentsOptions() *ListAccountSettingsAssignmentsOptions {
	return &ListAccountSettingsAssignmentsOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListAccountSettingsAssignmentsOptions) SetAccountID(accountID string) *ListAccountSettingsAssignmentsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTemplateID : Allow user to set TemplateID
func (_options *ListAccountSettingsAssignmentsOptions) SetTemplateID(templateID string) *ListAccountSettingsAssignmentsOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetTemplateVersion : Allow user to set TemplateVersion
func (_options *ListAccountSettingsAssignmentsOptions) SetTemplateVersion(templateVersion string) *ListAccountSettingsAssignmentsOptions {
	_options.TemplateVersion = core.StringPtr(templateVersion)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *ListAccountSettingsAssignmentsOptions) SetTarget(target string) *ListAccountSettingsAssignmentsOptions {
	_options.Target = core.StringPtr(target)
	return _options
}

// SetTargetType : Allow user to set TargetType
func (_options *ListAccountSettingsAssignmentsOptions) SetTargetType(targetType string) *ListAccountSettingsAssignmentsOptions {
	_options.TargetType = core.StringPtr(targetType)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAccountSettingsAssignmentsOptions) SetLimit(limit int64) *ListAccountSettingsAssignmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListAccountSettingsAssignmentsOptions) SetPagetoken(pagetoken string) *ListAccountSettingsAssignmentsOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListAccountSettingsAssignmentsOptions) SetSort(sort string) *ListAccountSettingsAssignmentsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListAccountSettingsAssignmentsOptions) SetOrder(order string) *ListAccountSettingsAssignmentsOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListAccountSettingsAssignmentsOptions) SetIncludeHistory(includeHistory bool) *ListAccountSettingsAssignmentsOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAccountSettingsAssignmentsOptions) SetHeaders(param map[string]string) *ListAccountSettingsAssignmentsOptions {
	options.Headers = param
	return options
}

// ListAccountSettingsTemplatesOptions : The ListAccountSettingsTemplates options.
type ListAccountSettingsTemplatesOptions struct {
	// Account ID of the account settings templates to query. This parameter is required unless using a pagetoken.
	AccountID *string `json:"account_id,omitempty"`

	// Optional size of a single page.
	Limit *string `json:"limit,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// Optional sort property. If specified, the returned templated are sorted according to this property.
	Sort *string `json:"sort,omitempty"`

	// Optional sort order.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *string `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListAccountSettingsTemplatesOptions.Sort property.
// Optional sort property. If specified, the returned templated are sorted according to this property.
const (
	ListAccountSettingsTemplatesOptionsSortCreatedAtConst      = "created_at"
	ListAccountSettingsTemplatesOptionsSortLastModifiedAtConst = "last_modified_at"
	ListAccountSettingsTemplatesOptionsSortNameConst           = "name"
)

// Constants associated with the ListAccountSettingsTemplatesOptions.Order property.
// Optional sort order.
const (
	ListAccountSettingsTemplatesOptionsOrderAscConst  = "asc"
	ListAccountSettingsTemplatesOptionsOrderDescConst = "desc"
)

// NewListAccountSettingsTemplatesOptions : Instantiate ListAccountSettingsTemplatesOptions
func (*IamIdentityV1) NewListAccountSettingsTemplatesOptions() *ListAccountSettingsTemplatesOptions {
	return &ListAccountSettingsTemplatesOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListAccountSettingsTemplatesOptions) SetAccountID(accountID string) *ListAccountSettingsTemplatesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAccountSettingsTemplatesOptions) SetLimit(limit string) *ListAccountSettingsTemplatesOptions {
	_options.Limit = core.StringPtr(limit)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListAccountSettingsTemplatesOptions) SetPagetoken(pagetoken string) *ListAccountSettingsTemplatesOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListAccountSettingsTemplatesOptions) SetSort(sort string) *ListAccountSettingsTemplatesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListAccountSettingsTemplatesOptions) SetOrder(order string) *ListAccountSettingsTemplatesOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListAccountSettingsTemplatesOptions) SetIncludeHistory(includeHistory string) *ListAccountSettingsTemplatesOptions {
	_options.IncludeHistory = core.StringPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAccountSettingsTemplatesOptions) SetHeaders(param map[string]string) *ListAccountSettingsTemplatesOptions {
	options.Headers = param
	return options
}

// ListAPIKeysOptions : The ListAPIKeys options.
type ListAPIKeysOptions struct {
	// Account ID of the API keys to query. If a service IAM ID is specified in iam_id then account_id must match the
	// account of the IAM ID. If a user IAM ID is specified in iam_id then then account_id must match the account of the
	// Authorization token.
	AccountID *string `json:"account_id,omitempty"`

	// IAM ID of the API keys to be queried. The IAM ID may be that of a user or a service. For a user IAM ID iam_id must
	// match the Authorization token.
	IamID *string `json:"iam_id,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Pagesize *int64 `json:"pagesize,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// Optional parameter to define the scope of the queried API keys. Can be 'entity' (default) or 'account'.
	Scope *string `json:"scope,omitempty"`

	// Optional parameter to filter the type of the queried API keys. Can be 'user' or 'serviceid'.
	Type *string `json:"type,omitempty"`

	// Optional sort property, valid values are name, description, created_at and created_by. If specified, the items are
	// sorted by the value of this property.
	Sort *string `json:"sort,omitempty"`

	// Optional sort order, valid values are asc and desc. Default: asc.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListAPIKeysOptions.Scope property.
// Optional parameter to define the scope of the queried API keys. Can be 'entity' (default) or 'account'.
const (
	ListAPIKeysOptionsScopeAccountConst = "account"
	ListAPIKeysOptionsScopeEntityConst  = "entity"
)

// Constants associated with the ListAPIKeysOptions.Type property.
// Optional parameter to filter the type of the queried API keys. Can be 'user' or 'serviceid'.
const (
	ListAPIKeysOptionsTypeServiceidConst = "serviceid"
	ListAPIKeysOptionsTypeUserConst      = "user"
)

// Constants associated with the ListAPIKeysOptions.Order property.
// Optional sort order, valid values are asc and desc. Default: asc.
const (
	ListAPIKeysOptionsOrderAscConst  = "asc"
	ListAPIKeysOptionsOrderDescConst = "desc"
)

// NewListAPIKeysOptions : Instantiate ListAPIKeysOptions
func (*IamIdentityV1) NewListAPIKeysOptions() *ListAPIKeysOptions {
	return &ListAPIKeysOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListAPIKeysOptions) SetAccountID(accountID string) *ListAPIKeysOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetIamID : Allow user to set IamID
func (_options *ListAPIKeysOptions) SetIamID(iamID string) *ListAPIKeysOptions {
	_options.IamID = core.StringPtr(iamID)
	return _options
}

// SetPagesize : Allow user to set Pagesize
func (_options *ListAPIKeysOptions) SetPagesize(pagesize int64) *ListAPIKeysOptions {
	_options.Pagesize = core.Int64Ptr(pagesize)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListAPIKeysOptions) SetPagetoken(pagetoken string) *ListAPIKeysOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetScope : Allow user to set Scope
func (_options *ListAPIKeysOptions) SetScope(scope string) *ListAPIKeysOptions {
	_options.Scope = core.StringPtr(scope)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListAPIKeysOptions) SetType(typeVar string) *ListAPIKeysOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListAPIKeysOptions) SetSort(sort string) *ListAPIKeysOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListAPIKeysOptions) SetOrder(order string) *ListAPIKeysOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListAPIKeysOptions) SetIncludeHistory(includeHistory bool) *ListAPIKeysOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAPIKeysOptions) SetHeaders(param map[string]string) *ListAPIKeysOptions {
	options.Headers = param
	return options
}

// ListClaimRulesOptions : The ListClaimRules options.
type ListClaimRulesOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListClaimRulesOptions : Instantiate ListClaimRulesOptions
func (*IamIdentityV1) NewListClaimRulesOptions(profileID string) *ListClaimRulesOptions {
	return &ListClaimRulesOptions{
		ProfileID: core.StringPtr(profileID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *ListClaimRulesOptions) SetProfileID(profileID string) *ListClaimRulesOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListClaimRulesOptions) SetHeaders(param map[string]string) *ListClaimRulesOptions {
	options.Headers = param
	return options
}

// ListLinksOptions : The ListLinks options.
type ListLinksOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListLinksOptions : Instantiate ListLinksOptions
func (*IamIdentityV1) NewListLinksOptions(profileID string) *ListLinksOptions {
	return &ListLinksOptions{
		ProfileID: core.StringPtr(profileID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *ListLinksOptions) SetProfileID(profileID string) *ListLinksOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListLinksOptions) SetHeaders(param map[string]string) *ListLinksOptions {
	options.Headers = param
	return options
}

// ListProfileTemplatesOptions : The ListProfileTemplates options.
type ListProfileTemplatesOptions struct {
	// Account ID of the trusted profile templates to query. This parameter is required unless using a pagetoken.
	AccountID *string `json:"account_id,omitempty"`

	// Optional size of a single page.
	Limit *string `json:"limit,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// Optional sort property. If specified, the returned templates are sorted according to this property.
	Sort *string `json:"sort,omitempty"`

	// Optional sort order.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *string `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListProfileTemplatesOptions.Sort property.
// Optional sort property. If specified, the returned templates are sorted according to this property.
const (
	ListProfileTemplatesOptionsSortCreatedAtConst      = "created_at"
	ListProfileTemplatesOptionsSortLastModifiedAtConst = "last_modified_at"
	ListProfileTemplatesOptionsSortNameConst           = "name"
)

// Constants associated with the ListProfileTemplatesOptions.Order property.
// Optional sort order.
const (
	ListProfileTemplatesOptionsOrderAscConst  = "asc"
	ListProfileTemplatesOptionsOrderDescConst = "desc"
)

// NewListProfileTemplatesOptions : Instantiate ListProfileTemplatesOptions
func (*IamIdentityV1) NewListProfileTemplatesOptions() *ListProfileTemplatesOptions {
	return &ListProfileTemplatesOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListProfileTemplatesOptions) SetAccountID(accountID string) *ListProfileTemplatesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListProfileTemplatesOptions) SetLimit(limit string) *ListProfileTemplatesOptions {
	_options.Limit = core.StringPtr(limit)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListProfileTemplatesOptions) SetPagetoken(pagetoken string) *ListProfileTemplatesOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListProfileTemplatesOptions) SetSort(sort string) *ListProfileTemplatesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListProfileTemplatesOptions) SetOrder(order string) *ListProfileTemplatesOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListProfileTemplatesOptions) SetIncludeHistory(includeHistory string) *ListProfileTemplatesOptions {
	_options.IncludeHistory = core.StringPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProfileTemplatesOptions) SetHeaders(param map[string]string) *ListProfileTemplatesOptions {
	options.Headers = param
	return options
}

// ListProfilesOptions : The ListProfiles options.
type ListProfilesOptions struct {
	// Account ID to query for trusted profiles.
	AccountID *string `json:"account_id" validate:"required"`

	// Name of the trusted profile to query.
	Name *string `json:"name,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Pagesize *int64 `json:"pagesize,omitempty"`

	// Optional sort property, valid values are name, description, created_at and modified_at. If specified, the items are
	// sorted by the value of this property.
	Sort *string `json:"sort,omitempty"`

	// Optional sort order, valid values are asc and desc. Default: asc.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListProfilesOptions.Order property.
// Optional sort order, valid values are asc and desc. Default: asc.
const (
	ListProfilesOptionsOrderAscConst  = "asc"
	ListProfilesOptionsOrderDescConst = "desc"
)

// NewListProfilesOptions : Instantiate ListProfilesOptions
func (*IamIdentityV1) NewListProfilesOptions(accountID string) *ListProfilesOptions {
	return &ListProfilesOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListProfilesOptions) SetAccountID(accountID string) *ListProfilesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListProfilesOptions) SetName(name string) *ListProfilesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPagesize : Allow user to set Pagesize
func (_options *ListProfilesOptions) SetPagesize(pagesize int64) *ListProfilesOptions {
	_options.Pagesize = core.Int64Ptr(pagesize)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListProfilesOptions) SetSort(sort string) *ListProfilesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListProfilesOptions) SetOrder(order string) *ListProfilesOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListProfilesOptions) SetIncludeHistory(includeHistory bool) *ListProfilesOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListProfilesOptions) SetPagetoken(pagetoken string) *ListProfilesOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProfilesOptions) SetHeaders(param map[string]string) *ListProfilesOptions {
	options.Headers = param
	return options
}

// ListServiceIdsOptions : The ListServiceIds options.
type ListServiceIdsOptions struct {
	// Account ID of the service ID(s) to query. This parameter is required (unless using a pagetoken).
	AccountID *string `json:"account_id,omitempty"`

	// Name of the service ID(s) to query. Optional.20 items per page. Valid range is 1 to 100.
	Name *string `json:"name,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Pagesize *int64 `json:"pagesize,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// Optional sort property, valid values are name, description, created_at and modified_at. If specified, the items are
	// sorted by the value of this property.
	Sort *string `json:"sort,omitempty"`

	// Optional sort order, valid values are asc and desc. Default: asc.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListServiceIdsOptions.Order property.
// Optional sort order, valid values are asc and desc. Default: asc.
const (
	ListServiceIdsOptionsOrderAscConst  = "asc"
	ListServiceIdsOptionsOrderDescConst = "desc"
)

// NewListServiceIdsOptions : Instantiate ListServiceIdsOptions
func (*IamIdentityV1) NewListServiceIdsOptions() *ListServiceIdsOptions {
	return &ListServiceIdsOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListServiceIdsOptions) SetAccountID(accountID string) *ListServiceIdsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListServiceIdsOptions) SetName(name string) *ListServiceIdsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPagesize : Allow user to set Pagesize
func (_options *ListServiceIdsOptions) SetPagesize(pagesize int64) *ListServiceIdsOptions {
	_options.Pagesize = core.Int64Ptr(pagesize)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListServiceIdsOptions) SetPagetoken(pagetoken string) *ListServiceIdsOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListServiceIdsOptions) SetSort(sort string) *ListServiceIdsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListServiceIdsOptions) SetOrder(order string) *ListServiceIdsOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListServiceIdsOptions) SetIncludeHistory(includeHistory bool) *ListServiceIdsOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListServiceIdsOptions) SetHeaders(param map[string]string) *ListServiceIdsOptions {
	options.Headers = param
	return options
}

// ListTrustedProfileAssignmentsOptions : The ListTrustedProfileAssignments options.
type ListTrustedProfileAssignmentsOptions struct {
	// Account ID of the Assignments to query. This parameter is required unless using a pagetoken.
	AccountID *string `json:"account_id,omitempty"`

	// Filter results by Template Id.
	TemplateID *string `json:"template_id,omitempty"`

	// Filter results Template Version.
	TemplateVersion *string `json:"template_version,omitempty"`

	// Filter results by the assignment target.
	Target *string `json:"target,omitempty"`

	// Filter results by the assignment's target type.
	TargetType *string `json:"target_type,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Limit *int64 `json:"limit,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// If specified, the items are sorted by the value of this property.
	Sort *string `json:"sort,omitempty"`

	// Sort order.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListTrustedProfileAssignmentsOptions.TargetType property.
// Filter results by the assignment's target type.
const (
	ListTrustedProfileAssignmentsOptionsTargetTypeAccountConst      = "Account"
	ListTrustedProfileAssignmentsOptionsTargetTypeAccountgroupConst = "AccountGroup"
)

// Constants associated with the ListTrustedProfileAssignmentsOptions.Sort property.
// If specified, the items are sorted by the value of this property.
const (
	ListTrustedProfileAssignmentsOptionsSortCreatedAtConst      = "created_at"
	ListTrustedProfileAssignmentsOptionsSortLastModifiedAtConst = "last_modified_at"
	ListTrustedProfileAssignmentsOptionsSortTemplateIDConst     = "template_id"
)

// Constants associated with the ListTrustedProfileAssignmentsOptions.Order property.
// Sort order.
const (
	ListTrustedProfileAssignmentsOptionsOrderAscConst  = "asc"
	ListTrustedProfileAssignmentsOptionsOrderDescConst = "desc"
)

// NewListTrustedProfileAssignmentsOptions : Instantiate ListTrustedProfileAssignmentsOptions
func (*IamIdentityV1) NewListTrustedProfileAssignmentsOptions() *ListTrustedProfileAssignmentsOptions {
	return &ListTrustedProfileAssignmentsOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListTrustedProfileAssignmentsOptions) SetAccountID(accountID string) *ListTrustedProfileAssignmentsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTemplateID : Allow user to set TemplateID
func (_options *ListTrustedProfileAssignmentsOptions) SetTemplateID(templateID string) *ListTrustedProfileAssignmentsOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetTemplateVersion : Allow user to set TemplateVersion
func (_options *ListTrustedProfileAssignmentsOptions) SetTemplateVersion(templateVersion string) *ListTrustedProfileAssignmentsOptions {
	_options.TemplateVersion = core.StringPtr(templateVersion)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *ListTrustedProfileAssignmentsOptions) SetTarget(target string) *ListTrustedProfileAssignmentsOptions {
	_options.Target = core.StringPtr(target)
	return _options
}

// SetTargetType : Allow user to set TargetType
func (_options *ListTrustedProfileAssignmentsOptions) SetTargetType(targetType string) *ListTrustedProfileAssignmentsOptions {
	_options.TargetType = core.StringPtr(targetType)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListTrustedProfileAssignmentsOptions) SetLimit(limit int64) *ListTrustedProfileAssignmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListTrustedProfileAssignmentsOptions) SetPagetoken(pagetoken string) *ListTrustedProfileAssignmentsOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListTrustedProfileAssignmentsOptions) SetSort(sort string) *ListTrustedProfileAssignmentsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListTrustedProfileAssignmentsOptions) SetOrder(order string) *ListTrustedProfileAssignmentsOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListTrustedProfileAssignmentsOptions) SetIncludeHistory(includeHistory bool) *ListTrustedProfileAssignmentsOptions {
	_options.IncludeHistory = core.BoolPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTrustedProfileAssignmentsOptions) SetHeaders(param map[string]string) *ListTrustedProfileAssignmentsOptions {
	options.Headers = param
	return options
}

// ListVersionsOfAccountSettingsTemplateOptions : The ListVersionsOfAccountSettingsTemplate options.
type ListVersionsOfAccountSettingsTemplateOptions struct {
	// ID of the account settings template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Optional size of a single page.
	Limit *string `json:"limit,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// Optional sort property. If specified, the returned templated are sorted according to this property.
	Sort *string `json:"sort,omitempty"`

	// Optional sort order.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *string `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListVersionsOfAccountSettingsTemplateOptions.Sort property.
// Optional sort property. If specified, the returned templated are sorted according to this property.
const (
	ListVersionsOfAccountSettingsTemplateOptionsSortCreatedAtConst      = "created_at"
	ListVersionsOfAccountSettingsTemplateOptionsSortLastModifiedAtConst = "last_modified_at"
	ListVersionsOfAccountSettingsTemplateOptionsSortNameConst           = "name"
)

// Constants associated with the ListVersionsOfAccountSettingsTemplateOptions.Order property.
// Optional sort order.
const (
	ListVersionsOfAccountSettingsTemplateOptionsOrderAscConst  = "asc"
	ListVersionsOfAccountSettingsTemplateOptionsOrderDescConst = "desc"
)

// NewListVersionsOfAccountSettingsTemplateOptions : Instantiate ListVersionsOfAccountSettingsTemplateOptions
func (*IamIdentityV1) NewListVersionsOfAccountSettingsTemplateOptions(templateID string) *ListVersionsOfAccountSettingsTemplateOptions {
	return &ListVersionsOfAccountSettingsTemplateOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *ListVersionsOfAccountSettingsTemplateOptions) SetTemplateID(templateID string) *ListVersionsOfAccountSettingsTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListVersionsOfAccountSettingsTemplateOptions) SetLimit(limit string) *ListVersionsOfAccountSettingsTemplateOptions {
	_options.Limit = core.StringPtr(limit)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListVersionsOfAccountSettingsTemplateOptions) SetPagetoken(pagetoken string) *ListVersionsOfAccountSettingsTemplateOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListVersionsOfAccountSettingsTemplateOptions) SetSort(sort string) *ListVersionsOfAccountSettingsTemplateOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListVersionsOfAccountSettingsTemplateOptions) SetOrder(order string) *ListVersionsOfAccountSettingsTemplateOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListVersionsOfAccountSettingsTemplateOptions) SetIncludeHistory(includeHistory string) *ListVersionsOfAccountSettingsTemplateOptions {
	_options.IncludeHistory = core.StringPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListVersionsOfAccountSettingsTemplateOptions) SetHeaders(param map[string]string) *ListVersionsOfAccountSettingsTemplateOptions {
	options.Headers = param
	return options
}

// ListVersionsOfProfileTemplateOptions : The ListVersionsOfProfileTemplate options.
type ListVersionsOfProfileTemplateOptions struct {
	// ID of the trusted profile template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Optional size of a single page.
	Limit *string `json:"limit,omitempty"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"pagetoken,omitempty"`

	// Optional sort property. If specified, the returned templated are sorted according to this property.
	Sort *string `json:"sort,omitempty"`

	// Optional sort order.
	Order *string `json:"order,omitempty"`

	// Defines if the entity history is included in the response.
	IncludeHistory *string `json:"include_history,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListVersionsOfProfileTemplateOptions.Sort property.
// Optional sort property. If specified, the returned templated are sorted according to this property.
const (
	ListVersionsOfProfileTemplateOptionsSortCreatedAtConst      = "created_at"
	ListVersionsOfProfileTemplateOptionsSortLastModifiedAtConst = "last_modified_at"
	ListVersionsOfProfileTemplateOptionsSortNameConst           = "name"
)

// Constants associated with the ListVersionsOfProfileTemplateOptions.Order property.
// Optional sort order.
const (
	ListVersionsOfProfileTemplateOptionsOrderAscConst  = "asc"
	ListVersionsOfProfileTemplateOptionsOrderDescConst = "desc"
)

// NewListVersionsOfProfileTemplateOptions : Instantiate ListVersionsOfProfileTemplateOptions
func (*IamIdentityV1) NewListVersionsOfProfileTemplateOptions(templateID string) *ListVersionsOfProfileTemplateOptions {
	return &ListVersionsOfProfileTemplateOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *ListVersionsOfProfileTemplateOptions) SetTemplateID(templateID string) *ListVersionsOfProfileTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListVersionsOfProfileTemplateOptions) SetLimit(limit string) *ListVersionsOfProfileTemplateOptions {
	_options.Limit = core.StringPtr(limit)
	return _options
}

// SetPagetoken : Allow user to set Pagetoken
func (_options *ListVersionsOfProfileTemplateOptions) SetPagetoken(pagetoken string) *ListVersionsOfProfileTemplateOptions {
	_options.Pagetoken = core.StringPtr(pagetoken)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListVersionsOfProfileTemplateOptions) SetSort(sort string) *ListVersionsOfProfileTemplateOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetOrder : Allow user to set Order
func (_options *ListVersionsOfProfileTemplateOptions) SetOrder(order string) *ListVersionsOfProfileTemplateOptions {
	_options.Order = core.StringPtr(order)
	return _options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (_options *ListVersionsOfProfileTemplateOptions) SetIncludeHistory(includeHistory string) *ListVersionsOfProfileTemplateOptions {
	_options.IncludeHistory = core.StringPtr(includeHistory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListVersionsOfProfileTemplateOptions) SetHeaders(param map[string]string) *ListVersionsOfProfileTemplateOptions {
	options.Headers = param
	return options
}

// LockAPIKeyOptions : The LockAPIKey options.
type LockAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewLockAPIKeyOptions : Instantiate LockAPIKeyOptions
func (*IamIdentityV1) NewLockAPIKeyOptions(id string) *LockAPIKeyOptions {
	return &LockAPIKeyOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *LockAPIKeyOptions) SetID(id string) *LockAPIKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *LockAPIKeyOptions) SetHeaders(param map[string]string) *LockAPIKeyOptions {
	options.Headers = param
	return options
}

// LockServiceIDOptions : The LockServiceID options.
type LockServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewLockServiceIDOptions : Instantiate LockServiceIDOptions
func (*IamIdentityV1) NewLockServiceIDOptions(id string) *LockServiceIDOptions {
	return &LockServiceIDOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *LockServiceIDOptions) SetID(id string) *LockServiceIDOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *LockServiceIDOptions) SetHeaders(param map[string]string) *LockServiceIDOptions {
	options.Headers = param
	return options
}

// MfaEnrollmentTypeStatus : MfaEnrollmentTypeStatus struct
type MfaEnrollmentTypeStatus struct {
	// Describes whether the enrollment type is required.
	Required *bool `json:"required" validate:"required"`

	// Describes whether the enrollment type is enrolled.
	Enrolled *bool `json:"enrolled" validate:"required"`
}

// UnmarshalMfaEnrollmentTypeStatus unmarshals an instance of MfaEnrollmentTypeStatus from the specified map of raw messages.
func UnmarshalMfaEnrollmentTypeStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MfaEnrollmentTypeStatus)
	err = core.UnmarshalPrimitive(m, "required", &obj.Required)
	if err != nil {
		err = core.SDKErrorf(err, "", "required-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enrolled", &obj.Enrolled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enrolled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// MfaEnrollments : MfaEnrollments struct
type MfaEnrollments struct {
	// currently effective mfa type i.e. id_based_mfa or account_based_mfa.
	EffectiveMfaType *string `json:"effective_mfa_type" validate:"required"`

	IDBasedMfa *IDBasedMfaEnrollment `json:"id_based_mfa,omitempty"`

	AccountBasedMfa *AccountBasedMfaEnrollment `json:"account_based_mfa,omitempty"`
}

// UnmarshalMfaEnrollments unmarshals an instance of MfaEnrollments from the specified map of raw messages.
func UnmarshalMfaEnrollments(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MfaEnrollments)
	err = core.UnmarshalPrimitive(m, "effective_mfa_type", &obj.EffectiveMfaType)
	if err != nil {
		err = core.SDKErrorf(err, "", "effective_mfa_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "id_based_mfa", &obj.IDBasedMfa, UnmarshalIDBasedMfaEnrollment)
	if err != nil {
		err = core.SDKErrorf(err, "", "id_based_mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "account_based_mfa", &obj.AccountBasedMfa, UnmarshalAccountBasedMfaEnrollment)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_based_mfa-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyTemplateReference : Metadata for external access policy.
type PolicyTemplateReference struct {
	// ID of Access Policy Template.
	ID *string `json:"id" validate:"required"`

	// Version of Access Policy Template.
	Version *string `json:"version" validate:"required"`
}

// NewPolicyTemplateReference : Instantiate PolicyTemplateReference (Generic Model Constructor)
func (*IamIdentityV1) NewPolicyTemplateReference(id string, version string) (_model *PolicyTemplateReference, err error) {
	_model = &PolicyTemplateReference{
		ID:      core.StringPtr(id),
		Version: core.StringPtr(version),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalPolicyTemplateReference unmarshals an instance of PolicyTemplateReference from the specified map of raw messages.
func UnmarshalPolicyTemplateReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileClaimRule : ProfileClaimRule struct
type ProfileClaimRule struct {
	// the unique identifier of the claim rule.
	ID *string `json:"id" validate:"required"`

	// version of the claim rule.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// If set contains a date time string of the creation date in ISO format.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// If set contains a date time string of the last modification date in ISO format.
	ModifiedAt *strfmt.DateTime `json:"modified_at,omitempty"`

	// The optional claim rule name.
	Name *string `json:"name,omitempty"`

	// Type of the claim rule, either 'Profile-SAML' or 'Profile-CR'.
	Type *string `json:"type" validate:"required"`

	// The realm name of the Idp this claim rule applies to.
	RealmName *string `json:"realm_name,omitempty"`

	// Session expiration in seconds.
	Expiration *int64 `json:"expiration" validate:"required"`

	// The compute resource type. Not required if type is Profile-SAML. Valid values are VSI, IKS_SA, ROKS_SA.
	CrType *string `json:"cr_type,omitempty"`

	// Conditions of this claim rule.
	Conditions []ProfileClaimRuleConditions `json:"conditions" validate:"required"`
}

// UnmarshalProfileClaimRule unmarshals an instance of ProfileClaimRule from the specified map of raw messages.
func UnmarshalProfileClaimRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileClaimRule)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "realm_name", &obj.RealmName)
	if err != nil {
		err = core.SDKErrorf(err, "", "realm_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration", &obj.Expiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cr_type", &obj.CrType)
	if err != nil {
		err = core.SDKErrorf(err, "", "cr_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalProfileClaimRuleConditions)
	if err != nil {
		err = core.SDKErrorf(err, "", "conditions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileClaimRuleConditions : ProfileClaimRuleConditions struct
type ProfileClaimRuleConditions struct {
	// The claim to evaluate against. [Learn
	// more](/docs/account?topic=account-iam-condition-properties&interface=ui#cr-attribute-names).
	Claim *string `json:"claim" validate:"required"`

	// The operation to perform on the claim. valid values are EQUALS, NOT_EQUALS, EQUALS_IGNORE_CASE,
	// NOT_EQUALS_IGNORE_CASE, CONTAINS, IN.
	Operator *string `json:"operator" validate:"required"`

	// The stringified JSON value that the claim is compared to using the operator.
	Value *string `json:"value" validate:"required"`
}

// NewProfileClaimRuleConditions : Instantiate ProfileClaimRuleConditions (Generic Model Constructor)
func (*IamIdentityV1) NewProfileClaimRuleConditions(claim string, operator string, value string) (_model *ProfileClaimRuleConditions, err error) {
	_model = &ProfileClaimRuleConditions{
		Claim:    core.StringPtr(claim),
		Operator: core.StringPtr(operator),
		Value:    core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalProfileClaimRuleConditions unmarshals an instance of ProfileClaimRuleConditions from the specified map of raw messages.
func UnmarshalProfileClaimRuleConditions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileClaimRuleConditions)
	err = core.UnmarshalPrimitive(m, "claim", &obj.Claim)
	if err != nil {
		err = core.SDKErrorf(err, "", "claim-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileClaimRuleList : ProfileClaimRuleList struct
type ProfileClaimRuleList struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// List of claim rules.
	Rules []ProfileClaimRule `json:"rules" validate:"required"`
}

// UnmarshalProfileClaimRuleList unmarshals an instance of ProfileClaimRuleList from the specified map of raw messages.
func UnmarshalProfileClaimRuleList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileClaimRuleList)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalProfileClaimRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileIdentitiesResponse : ProfileIdentitiesResponse struct
type ProfileIdentitiesResponse struct {
	// Entity tag of the profile identities response.
	EntityTag *string `json:"entity_tag,omitempty"`

	// List of identities.
	Identities []ProfileIdentityResponse `json:"identities,omitempty"`
}

// UnmarshalProfileIdentitiesResponse unmarshals an instance of ProfileIdentitiesResponse from the specified map of raw messages.
func UnmarshalProfileIdentitiesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileIdentitiesResponse)
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "identities", &obj.Identities, UnmarshalProfileIdentityResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "identities-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileIdentityRequest : ProfileIdentityRequest struct
type ProfileIdentityRequest struct {
	// Identifier of the identity that can assume the trusted profiles. This can be a user identifier (IAM id), serviceid
	// or crn. Internally it uses account id of the service id for the identifier 'serviceid' and for the identifier 'crn'
	// it uses account id contained in the CRN.
	Identifier *string `json:"identifier" validate:"required"`

	// Type of the identity.
	Type *string `json:"type" validate:"required"`

	// Only valid for the type user. Accounts from which a user can assume the trusted profile.
	Accounts []string `json:"accounts,omitempty"`

	// Description of the identity that can assume the trusted profile. This is optional field for all the types of
	// identities. When this field is not set for the identity type 'serviceid' then the description of the service id is
	// used. Description is recommended for the identity type 'crn' E.g. 'Instance 1234 of IBM Cloud Service project'.
	Description *string `json:"description,omitempty"`
}

// Constants associated with the ProfileIdentityRequest.Type property.
// Type of the identity.
const (
	ProfileIdentityRequestTypeCRNConst       = "crn"
	ProfileIdentityRequestTypeServiceidConst = "serviceid"
	ProfileIdentityRequestTypeUserConst      = "user"
)

// NewProfileIdentityRequest : Instantiate ProfileIdentityRequest (Generic Model Constructor)
func (*IamIdentityV1) NewProfileIdentityRequest(identifier string, typeVar string) (_model *ProfileIdentityRequest, err error) {
	_model = &ProfileIdentityRequest{
		Identifier: core.StringPtr(identifier),
		Type:       core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalProfileIdentityRequest unmarshals an instance of ProfileIdentityRequest from the specified map of raw messages.
func UnmarshalProfileIdentityRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileIdentityRequest)
	err = core.UnmarshalPrimitive(m, "identifier", &obj.Identifier)
	if err != nil {
		err = core.SDKErrorf(err, "", "identifier-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "accounts", &obj.Accounts)
	if err != nil {
		err = core.SDKErrorf(err, "", "accounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileIdentityResponse : ProfileIdentityResponse struct
type ProfileIdentityResponse struct {
	// IAM ID of the identity.
	IamID *string `json:"iam_id" validate:"required"`

	// Identifier of the identity that can assume the trusted profiles. This can be a user identifier (IAM id), serviceid
	// or crn. Internally it uses account id of the service id for the identifier 'serviceid' and for the identifier 'crn'
	// it uses account id contained in the CRN.
	Identifier *string `json:"identifier" validate:"required"`

	// Type of the identity.
	Type *string `json:"type" validate:"required"`

	// Only valid for the type user. Accounts from which a user can assume the trusted profile.
	Accounts []string `json:"accounts,omitempty"`

	// Description of the identity that can assume the trusted profile. This is optional field for all the types of
	// identities. When this field is not set for the identity type 'serviceid' then the description of the service id is
	// used. Description is recommended for the identity type 'crn' E.g. 'Instance 1234 of IBM Cloud Service project'.
	Description *string `json:"description,omitempty"`
}

// Constants associated with the ProfileIdentityResponse.Type property.
// Type of the identity.
const (
	ProfileIdentityResponseTypeCRNConst       = "crn"
	ProfileIdentityResponseTypeServiceidConst = "serviceid"
	ProfileIdentityResponseTypeUserConst      = "user"
)

// UnmarshalProfileIdentityResponse unmarshals an instance of ProfileIdentityResponse from the specified map of raw messages.
func UnmarshalProfileIdentityResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileIdentityResponse)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "identifier", &obj.Identifier)
	if err != nil {
		err = core.SDKErrorf(err, "", "identifier-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "accounts", &obj.Accounts)
	if err != nil {
		err = core.SDKErrorf(err, "", "accounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileLink : Link details.
type ProfileLink struct {
	// the unique identifier of the link.
	ID *string `json:"id" validate:"required"`

	// version of the link.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// If set contains a date time string of the creation date in ISO format.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// If set contains a date time string of the last modification date in ISO format.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// Optional name of the Link.
	Name *string `json:"name,omitempty"`

	// The compute resource type. Valid values are VSI, IKS_SA, ROKS_SA.
	CrType *string `json:"cr_type" validate:"required"`

	Link *ProfileLinkLink `json:"link" validate:"required"`
}

// UnmarshalProfileLink unmarshals an instance of ProfileLink from the specified map of raw messages.
func UnmarshalProfileLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileLink)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cr_type", &obj.CrType)
	if err != nil {
		err = core.SDKErrorf(err, "", "cr_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "link", &obj.Link, UnmarshalProfileLinkLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "link-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileLinkLink : ProfileLinkLink struct
type ProfileLinkLink struct {
	// The CRN of the compute resource.
	CRN *string `json:"crn,omitempty"`

	// The compute resource namespace, only required if cr_type is IKS_SA or ROKS_SA.
	Namespace *string `json:"namespace,omitempty"`

	// Name of the compute resource, only required if cr_type is IKS_SA or ROKS_SA.
	Name *string `json:"name,omitempty"`
}

// UnmarshalProfileLinkLink unmarshals an instance of ProfileLinkLink from the specified map of raw messages.
func UnmarshalProfileLinkLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileLinkLink)
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		err = core.SDKErrorf(err, "", "namespace-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileLinkList : ProfileLinkList struct
type ProfileLinkList struct {
	// List of links to a trusted profile.
	Links []ProfileLink `json:"links" validate:"required"`
}

// UnmarshalProfileLinkList unmarshals an instance of ProfileLinkList from the specified map of raw messages.
func UnmarshalProfileLinkList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileLinkList)
	err = core.UnmarshalModel(m, "links", &obj.Links, UnmarshalProfileLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "links-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Report : Report struct
type Report struct {
	// IAMid of the user who triggered the report.
	CreatedBy *string `json:"created_by" validate:"required"`

	// Unique reference used to generate the report.
	Reference *string `json:"reference" validate:"required"`

	// Duration in hours for which the report is generated.
	ReportDuration *string `json:"report_duration" validate:"required"`

	// Start time of the report.
	ReportStartTime *string `json:"report_start_time" validate:"required"`

	// End time of the report.
	ReportEndTime *string `json:"report_end_time" validate:"required"`

	// List of users.
	Users []UserActivity `json:"users,omitempty"`

	// List of apikeys.
	Apikeys []ApikeyActivity `json:"apikeys,omitempty"`

	// List of serviceids.
	Serviceids []EntityActivity `json:"serviceids,omitempty"`

	// List of profiles.
	Profiles []EntityActivity `json:"profiles,omitempty"`
}

// UnmarshalReport unmarshals an instance of Report from the specified map of raw messages.
func UnmarshalReport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Report)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		err = core.SDKErrorf(err, "", "reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "report_duration", &obj.ReportDuration)
	if err != nil {
		err = core.SDKErrorf(err, "", "report_duration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "report_start_time", &obj.ReportStartTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "report_start_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "report_end_time", &obj.ReportEndTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "report_end_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "users", &obj.Users, UnmarshalUserActivity)
	if err != nil {
		err = core.SDKErrorf(err, "", "users-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "apikeys", &obj.Apikeys, UnmarshalApikeyActivity)
	if err != nil {
		err = core.SDKErrorf(err, "", "apikeys-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "serviceids", &obj.Serviceids, UnmarshalEntityActivity)
	if err != nil {
		err = core.SDKErrorf(err, "", "serviceids-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalEntityActivity)
	if err != nil {
		err = core.SDKErrorf(err, "", "profiles-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportMfaEnrollmentStatus : ReportMfaEnrollmentStatus struct
type ReportMfaEnrollmentStatus struct {
	// IAMid of the user who triggered the report.
	CreatedBy *string `json:"created_by" validate:"required"`

	// Unique reference used to generate the report.
	Reference *string `json:"reference" validate:"required"`

	// Date time at which report is generated. Date is in ISO format.
	ReportTime *string `json:"report_time" validate:"required"`

	// BSS account id of the user who triggered the report.
	AccountID *string `json:"account_id" validate:"required"`

	// IMS account id of the user who triggered the report.
	ImsAccountID *string `json:"ims_account_id,omitempty"`

	// List of users.
	Users []UserReportMfaEnrollmentStatus `json:"users,omitempty"`
}

// UnmarshalReportMfaEnrollmentStatus unmarshals an instance of ReportMfaEnrollmentStatus from the specified map of raw messages.
func UnmarshalReportMfaEnrollmentStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportMfaEnrollmentStatus)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		err = core.SDKErrorf(err, "", "reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "report_time", &obj.ReportTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "report_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ims_account_id", &obj.ImsAccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "ims_account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "users", &obj.Users, UnmarshalUserReportMfaEnrollmentStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "users-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportReference : ReportReference struct
type ReportReference struct {
	// Reference for the report to be generated.
	Reference *string `json:"reference" validate:"required"`
}

// UnmarshalReportReference unmarshals an instance of ReportReference from the specified map of raw messages.
func UnmarshalReportReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportReference)
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		err = core.SDKErrorf(err, "", "reference-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResponseContext : Context with key properties for problem determination.
type ResponseContext struct {
	// The transaction ID of the inbound REST request.
	TransactionID *string `json:"transaction_id,omitempty"`

	// The operation of the inbound REST request.
	Operation *string `json:"operation,omitempty"`

	// The user agent of the inbound REST request.
	UserAgent *string `json:"user_agent,omitempty"`

	// The URL of that cluster.
	URL *string `json:"url,omitempty"`

	// The instance ID of the server instance processing the request.
	InstanceID *string `json:"instance_id,omitempty"`

	// The thread ID of the server instance processing the request.
	ThreadID *string `json:"thread_id,omitempty"`

	// The host of the server instance processing the request.
	Host *string `json:"host,omitempty"`

	// The start time of the request.
	StartTime *string `json:"start_time,omitempty"`

	// The finish time of the request.
	EndTime *string `json:"end_time,omitempty"`

	// The elapsed time in msec.
	ElapsedTime *string `json:"elapsed_time,omitempty"`

	// The cluster name.
	ClusterName *string `json:"cluster_name,omitempty"`
}

// UnmarshalResponseContext unmarshals an instance of ResponseContext from the specified map of raw messages.
func UnmarshalResponseContext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResponseContext)
	err = core.UnmarshalPrimitive(m, "transaction_id", &obj.TransactionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "transaction_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operation", &obj.Operation)
	if err != nil {
		err = core.SDKErrorf(err, "", "operation-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_agent", &obj.UserAgent)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_agent-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "thread_id", &obj.ThreadID)
	if err != nil {
		err = core.SDKErrorf(err, "", "thread_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host", &obj.Host)
	if err != nil {
		err = core.SDKErrorf(err, "", "host-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "start_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "end_time", &obj.EndTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "end_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "elapsed_time", &obj.ElapsedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "elapsed_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_name", &obj.ClusterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "cluster_name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceID : Response body format for service ID V1 REST requests.
type ServiceID struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// Unique identifier of this Service Id.
	ID *string `json:"id" validate:"required"`

	// Cloud wide identifier for identities of this service ID.
	IamID *string `json:"iam_id" validate:"required"`

	// Version of the service ID details object. You need to specify this value when updating the service ID to avoid stale
	// updates.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// Cloud Resource Name of the item. Example Cloud Resource Name:
	// 'crn:v1:bluemix:public:iam-identity:us-south:a/myaccount::serviceid:1234-5678-9012'.
	CRN *string `json:"crn" validate:"required"`

	// The service ID cannot be changed if set to true.
	Locked *bool `json:"locked" validate:"required"`

	// If set contains a date time string of the creation date in ISO format.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// If set contains a date time string of the last modification date in ISO format.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// ID of the account the service ID belongs to.
	AccountID *string `json:"account_id" validate:"required"`

	// Name of the Service Id. The name is not checked for uniqueness. Therefore multiple names with the same value can
	// exist. Access is done via the UUID of the Service Id.
	Name *string `json:"name" validate:"required"`

	// The optional description of the Service Id. The 'description' property is only available if a description was
	// provided during a create of a Service Id.
	Description *string `json:"description,omitempty"`

	// Optional list of CRNs (string array) which point to the services connected to the service ID.
	UniqueInstanceCrns []string `json:"unique_instance_crns,omitempty"`

	// History of the Service ID.
	History []EnityHistoryRecord `json:"history,omitempty"`

	// Response body format for API key V1 REST requests.
	Apikey *APIKey `json:"apikey,omitempty"`

	Activity *Activity `json:"activity,omitempty"`
}

// UnmarshalServiceID unmarshals an instance of ServiceID from the specified map of raw messages.
func UnmarshalServiceID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceID)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		err = core.SDKErrorf(err, "", "locked-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "unique_instance_crns", &obj.UniqueInstanceCrns)
	if err != nil {
		err = core.SDKErrorf(err, "", "unique_instance_crns-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "apikey", &obj.Apikey, UnmarshalAPIKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "apikey-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "activity", &obj.Activity, UnmarshalActivity)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceIDList : Response body format for the list service ID V1 REST request.
type ServiceIDList struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// The offset of the current page.
	Offset *int64 `json:"offset,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Limit *int64 `json:"limit,omitempty"`

	// Link to the first page.
	First *string `json:"first,omitempty"`

	// Link to the previous available page. If 'previous' property is not part of the response no previous page is
	// available.
	Previous *string `json:"previous,omitempty"`

	// Link to the next available page. If 'next' property is not part of the response no next page is available.
	Next *string `json:"next,omitempty"`

	// List of service IDs based on the query paramters and the page size. The service IDs array is always part of the
	// response but might be empty depending on the query parameter values provided.
	Serviceids []ServiceID `json:"serviceids" validate:"required"`
}

// UnmarshalServiceIDList unmarshals an instance of ServiceIDList from the specified map of raw messages.
func UnmarshalServiceIDList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceIDList)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "serviceids", &obj.Serviceids, UnmarshalServiceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "serviceids-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SetProfileIdentitiesOptions : The SetProfileIdentities options.
type SetProfileIdentitiesOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Entity tag of the Identities to be updated. Specify the tag that you retrieved when reading the Profile Identities.
	// This value helps identify parallel usage of this API. Pass * to indicate updating any available version, which may
	// result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// List of identities that can assume the trusted profile.
	Identities []ProfileIdentityRequest `json:"identities,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewSetProfileIdentitiesOptions : Instantiate SetProfileIdentitiesOptions
func (*IamIdentityV1) NewSetProfileIdentitiesOptions(profileID string, ifMatch string) *SetProfileIdentitiesOptions {
	return &SetProfileIdentitiesOptions{
		ProfileID: core.StringPtr(profileID),
		IfMatch:   core.StringPtr(ifMatch),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *SetProfileIdentitiesOptions) SetProfileID(profileID string) *SetProfileIdentitiesOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *SetProfileIdentitiesOptions) SetIfMatch(ifMatch string) *SetProfileIdentitiesOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetIdentities : Allow user to set Identities
func (_options *SetProfileIdentitiesOptions) SetIdentities(identities []ProfileIdentityRequest) *SetProfileIdentitiesOptions {
	_options.Identities = identities
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetProfileIdentitiesOptions) SetHeaders(param map[string]string) *SetProfileIdentitiesOptions {
	options.Headers = param
	return options
}

// SetProfileIdentityOptions : The SetProfileIdentity options.
type SetProfileIdentityOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Type of the identity.
	IdentityType *string `json:"identity-type" validate:"required,ne="`

	// Identifier of the identity that can assume the trusted profiles. This can be a user identifier (IAM id), serviceid
	// or crn. Internally it uses account id of the service id for the identifier 'serviceid' and for the identifier 'crn'
	// it uses account id contained in the CRN.
	Identifier *string `json:"identifier" validate:"required"`

	// Type of the identity.
	Type *string `json:"type" validate:"required"`

	// Only valid for the type user. Accounts from which a user can assume the trusted profile.
	Accounts []string `json:"accounts,omitempty"`

	// Description of the identity that can assume the trusted profile. This is optional field for all the types of
	// identities. When this field is not set for the identity type 'serviceid' then the description of the service id is
	// used. Description is recommended for the identity type 'crn' E.g. 'Instance 1234 of IBM Cloud Service project'.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the SetProfileIdentityOptions.IdentityType property.
// Type of the identity.
const (
	SetProfileIdentityOptionsIdentityTypeCRNConst       = "crn"
	SetProfileIdentityOptionsIdentityTypeServiceidConst = "serviceid"
	SetProfileIdentityOptionsIdentityTypeUserConst      = "user"
)

// Constants associated with the SetProfileIdentityOptions.Type property.
// Type of the identity.
const (
	SetProfileIdentityOptionsTypeCRNConst       = "crn"
	SetProfileIdentityOptionsTypeServiceidConst = "serviceid"
	SetProfileIdentityOptionsTypeUserConst      = "user"
)

// NewSetProfileIdentityOptions : Instantiate SetProfileIdentityOptions
func (*IamIdentityV1) NewSetProfileIdentityOptions(profileID string, identityType string, identifier string, typeVar string) *SetProfileIdentityOptions {
	return &SetProfileIdentityOptions{
		ProfileID:    core.StringPtr(profileID),
		IdentityType: core.StringPtr(identityType),
		Identifier:   core.StringPtr(identifier),
		Type:         core.StringPtr(typeVar),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *SetProfileIdentityOptions) SetProfileID(profileID string) *SetProfileIdentityOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetIdentityType : Allow user to set IdentityType
func (_options *SetProfileIdentityOptions) SetIdentityType(identityType string) *SetProfileIdentityOptions {
	_options.IdentityType = core.StringPtr(identityType)
	return _options
}

// SetIdentifier : Allow user to set Identifier
func (_options *SetProfileIdentityOptions) SetIdentifier(identifier string) *SetProfileIdentityOptions {
	_options.Identifier = core.StringPtr(identifier)
	return _options
}

// SetType : Allow user to set Type
func (_options *SetProfileIdentityOptions) SetType(typeVar string) *SetProfileIdentityOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetAccounts : Allow user to set Accounts
func (_options *SetProfileIdentityOptions) SetAccounts(accounts []string) *SetProfileIdentityOptions {
	_options.Accounts = accounts
	return _options
}

// SetDescription : Allow user to set Description
func (_options *SetProfileIdentityOptions) SetDescription(description string) *SetProfileIdentityOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetProfileIdentityOptions) SetHeaders(param map[string]string) *SetProfileIdentityOptions {
	options.Headers = param
	return options
}

// TemplateAssignmentListResponse : List Response body format for Template Assignments Records.
type TemplateAssignmentListResponse struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// The offset of the current page.
	Offset *int64 `json:"offset,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Limit *int64 `json:"limit,omitempty"`

	// Link to the first page.
	First *string `json:"first,omitempty"`

	// Link to the previous available page. If 'previous' property is not part of the response no previous page is
	// available.
	Previous *string `json:"previous,omitempty"`

	// Link to the next available page. If 'next' property is not part of the response no next page is available.
	Next *string `json:"next,omitempty"`

	// List of Assignments based on the query paramters and the page size. The assignments array is always part of the
	// response but might be empty depending on the query parameter values provided.
	Assignments []TemplateAssignmentResponse `json:"assignments" validate:"required"`
}

// UnmarshalTemplateAssignmentListResponse unmarshals an instance of TemplateAssignmentListResponse from the specified map of raw messages.
func UnmarshalTemplateAssignmentListResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAssignmentListResponse)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "assignments", &obj.Assignments, UnmarshalTemplateAssignmentResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "assignments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateAssignmentResource : Body parameters for created resource.
type TemplateAssignmentResource struct {
	// Id of the created resource.
	ID *string `json:"id,omitempty"`
}

// UnmarshalTemplateAssignmentResource unmarshals an instance of TemplateAssignmentResource from the specified map of raw messages.
func UnmarshalTemplateAssignmentResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAssignmentResource)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateAssignmentResourceError : Body parameters for assignment error.
type TemplateAssignmentResourceError struct {
	// Name of the error.
	Name *string `json:"name,omitempty"`

	// Internal error code.
	ErrorCode *string `json:"errorCode,omitempty"`

	// Error message detailing the nature of the error.
	Message *string `json:"message,omitempty"`

	// Internal status code for the error.
	StatusCode *string `json:"statusCode,omitempty"`
}

// UnmarshalTemplateAssignmentResourceError unmarshals an instance of TemplateAssignmentResourceError from the specified map of raw messages.
func UnmarshalTemplateAssignmentResourceError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAssignmentResourceError)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errorCode", &obj.ErrorCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "errorCode-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "statusCode", &obj.StatusCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "statusCode-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateAssignmentResponse : Response body format for Template Assignment Record.
type TemplateAssignmentResponse struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// Assignment record Id.
	ID *string `json:"id" validate:"required"`

	// Enterprise account Id.
	AccountID *string `json:"account_id" validate:"required"`

	// Template Id.
	TemplateID *string `json:"template_id" validate:"required"`

	// Template version.
	TemplateVersion *int64 `json:"template_version" validate:"required"`

	// Assignment target type.
	TargetType *string `json:"target_type" validate:"required"`

	// Assignment target.
	Target *string `json:"target" validate:"required"`

	// Assignment status.
	Status *string `json:"status" validate:"required"`

	// Status breakdown per target account of IAM resources created or errors encountered in attempting to create those IAM
	// resources. IAM resources are only included in the response providing the assignment is not in progress. IAM
	// resources are also only included when getting a single assignment, and excluded by list APIs.
	Resources []TemplateAssignmentResponseResource `json:"resources,omitempty"`

	// Assignment history.
	History []EnityHistoryRecord `json:"history,omitempty"`

	// Href.
	Href *string `json:"href,omitempty"`

	// Assignment created at.
	CreatedAt *string `json:"created_at" validate:"required"`

	// IAMid of the identity that created the assignment.
	CreatedByID *string `json:"created_by_id" validate:"required"`

	// Assignment modified at.
	LastModifiedAt *string `json:"last_modified_at" validate:"required"`

	// IAMid of the identity that last modified the assignment.
	LastModifiedByID *string `json:"last_modified_by_id" validate:"required"`

	// Entity tag for this assignment record.
	EntityTag *string `json:"entity_tag" validate:"required"`
}

// UnmarshalTemplateAssignmentResponse unmarshals an instance of TemplateAssignmentResponse from the specified map of raw messages.
func UnmarshalTemplateAssignmentResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAssignmentResponse)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_version", &obj.TemplateVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_type", &obj.TargetType)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalTemplateAssignmentResponseResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateAssignmentResponseResource : Overview of resources assignment per target account.
type TemplateAssignmentResponseResource struct {
	// Target account where the IAM resource is created.
	Target *string `json:"target" validate:"required"`

	Profile *TemplateAssignmentResponseResourceDetail `json:"profile,omitempty"`

	AccountSettings *TemplateAssignmentResponseResourceDetail `json:"account_settings,omitempty"`

	// Policy resource(s) included only for trusted profile assignments with policy references.
	PolicyTemplateRefs []TemplateAssignmentResponseResourceDetail `json:"policy_template_refs,omitempty"`
}

// UnmarshalTemplateAssignmentResponseResource unmarshals an instance of TemplateAssignmentResponseResource from the specified map of raw messages.
func UnmarshalTemplateAssignmentResponseResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAssignmentResponseResource)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "profile", &obj.Profile, UnmarshalTemplateAssignmentResponseResourceDetail)
	if err != nil {
		err = core.SDKErrorf(err, "", "profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "account_settings", &obj.AccountSettings, UnmarshalTemplateAssignmentResponseResourceDetail)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_settings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policy_template_refs", &obj.PolicyTemplateRefs, UnmarshalTemplateAssignmentResponseResourceDetail)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_template_refs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateAssignmentResponseResourceDetail : TemplateAssignmentResponseResourceDetail struct
type TemplateAssignmentResponseResourceDetail struct {
	// Policy Template Id, only returned for a profile assignment with policy references.
	ID *string `json:"id,omitempty"`

	// Policy version, only returned for a profile assignment with policy references.
	Version *string `json:"version,omitempty"`

	// Body parameters for created resource.
	ResourceCreated *TemplateAssignmentResource `json:"resource_created,omitempty"`

	// Body parameters for assignment error.
	ErrorMessage *TemplateAssignmentResourceError `json:"error_message,omitempty"`

	// Status for the target account's assignment.
	Status *string `json:"status" validate:"required"`
}

// UnmarshalTemplateAssignmentResponseResourceDetail unmarshals an instance of TemplateAssignmentResponseResourceDetail from the specified map of raw messages.
func UnmarshalTemplateAssignmentResponseResourceDetail(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAssignmentResponseResourceDetail)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resource_created", &obj.ResourceCreated, UnmarshalTemplateAssignmentResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_created-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "error_message", &obj.ErrorMessage, UnmarshalTemplateAssignmentResourceError)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_message-error", common.GetComponentInfo())
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

// TemplateProfileComponentRequest : Input body parameters for the TemplateProfileComponent.
type TemplateProfileComponentRequest struct {
	// Name of the Profile.
	Name *string `json:"name" validate:"required"`

	// Description of the Profile.
	Description *string `json:"description,omitempty"`

	// Rules for the Profile.
	Rules []TrustedProfileTemplateClaimRule `json:"rules,omitempty"`

	// Identities for the Profile.
	Identities []ProfileIdentityRequest `json:"identities,omitempty"`
}

// NewTemplateProfileComponentRequest : Instantiate TemplateProfileComponentRequest (Generic Model Constructor)
func (*IamIdentityV1) NewTemplateProfileComponentRequest(name string) (_model *TemplateProfileComponentRequest, err error) {
	_model = &TemplateProfileComponentRequest{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalTemplateProfileComponentRequest unmarshals an instance of TemplateProfileComponentRequest from the specified map of raw messages.
func UnmarshalTemplateProfileComponentRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateProfileComponentRequest)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalTrustedProfileTemplateClaimRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "identities", &obj.Identities, UnmarshalProfileIdentityRequest)
	if err != nil {
		err = core.SDKErrorf(err, "", "identities-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateProfileComponentResponse : Input body parameters for the TemplateProfileComponent.
type TemplateProfileComponentResponse struct {
	// Name of the Profile.
	Name *string `json:"name" validate:"required"`

	// Description of the Profile.
	Description *string `json:"description,omitempty"`

	// Rules for the Profile.
	Rules []TrustedProfileTemplateClaimRule `json:"rules,omitempty"`

	// Identities for the Profile.
	Identities []ProfileIdentityResponse `json:"identities,omitempty"`
}

// UnmarshalTemplateProfileComponentResponse unmarshals an instance of TemplateProfileComponentResponse from the specified map of raw messages.
func UnmarshalTemplateProfileComponentResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateProfileComponentResponse)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalTrustedProfileTemplateClaimRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "identities", &obj.Identities, UnmarshalProfileIdentityResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "identities-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TrustedProfile : Response body format for trusted profile V1 REST requests.
type TrustedProfile struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// the unique identifier of the trusted profile. Example:'Profile-94497d0d-2ac3-41bf-a993-a49d1b14627c'.
	ID *string `json:"id" validate:"required"`

	// Version of the trusted profile details object. You need to specify this value when updating the trusted profile to
	// avoid stale updates.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// Cloud Resource Name of the item. Example Cloud Resource Name:
	// 'crn:v1:bluemix:public:iam-identity:us-south:a/myaccount::profile:Profile-94497d0d-2ac3-41bf-a993-a49d1b14627c'.
	CRN *string `json:"crn" validate:"required"`

	// Name of the trusted profile. The name is checked for uniqueness. Therefore trusted profiles with the same names can
	// not exist in the same account.
	Name *string `json:"name" validate:"required"`

	// The optional description of the trusted profile. The 'description' property is only available if a description was
	// provided during a create of a trusted profile.
	Description *string `json:"description,omitempty"`

	// If set contains a date time string of the creation date in ISO format.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// If set contains a date time string of the last modification date in ISO format.
	ModifiedAt *strfmt.DateTime `json:"modified_at,omitempty"`

	// The iam_id of this trusted profile.
	IamID *string `json:"iam_id" validate:"required"`

	// ID of the account that this trusted profile belong to.
	AccountID *string `json:"account_id" validate:"required"`

	// ID of the IAM template that was used to create an enterprise-managed trusted profile in your account. When returned,
	// this indicates that the trusted profile is created from and managed by a template in the root enterprise account.
	TemplateID *string `json:"template_id,omitempty"`

	// ID of the assignment that was used to create an enterprise-managed trusted profile in your account. When returned,
	// this indicates that the trusted profile is created from and managed by a template in the root enterprise account.
	AssignmentID *string `json:"assignment_id,omitempty"`

	// IMS acount ID of the trusted profile.
	ImsAccountID *int64 `json:"ims_account_id,omitempty"`

	// IMS user ID of the trusted profile.
	ImsUserID *int64 `json:"ims_user_id,omitempty"`

	// History of the trusted profile.
	History []EnityHistoryRecord `json:"history,omitempty"`

	Activity *Activity `json:"activity,omitempty"`
}

// UnmarshalTrustedProfile unmarshals an instance of TrustedProfile from the specified map of raw messages.
func UnmarshalTrustedProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TrustedProfile)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "assignment_id", &obj.AssignmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "assignment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ims_account_id", &obj.ImsAccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "ims_account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ims_user_id", &obj.ImsUserID)
	if err != nil {
		err = core.SDKErrorf(err, "", "ims_user_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "activity", &obj.Activity, UnmarshalActivity)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TrustedProfileTemplateClaimRule : TrustedProfileTemplateClaimRule struct
type TrustedProfileTemplateClaimRule struct {
	// Name of the claim rule to be created or updated.
	Name *string `json:"name,omitempty"`

	// Type of the claim rule.
	Type *string `json:"type" validate:"required"`

	// The realm name of the Idp this claim rule applies to. This field is required only if the type is specified as
	// 'Profile-SAML'.
	RealmName *string `json:"realm_name,omitempty"`

	// Session expiration in seconds, only required if type is 'Profile-SAML'.
	Expiration *int64 `json:"expiration,omitempty"`

	// Conditions of this claim rule.
	Conditions []ProfileClaimRuleConditions `json:"conditions" validate:"required"`
}

// Constants associated with the TrustedProfileTemplateClaimRule.Type property.
// Type of the claim rule.
const (
	TrustedProfileTemplateClaimRuleTypeProfileSamlConst = "Profile-SAML"
)

// NewTrustedProfileTemplateClaimRule : Instantiate TrustedProfileTemplateClaimRule (Generic Model Constructor)
func (*IamIdentityV1) NewTrustedProfileTemplateClaimRule(typeVar string, conditions []ProfileClaimRuleConditions) (_model *TrustedProfileTemplateClaimRule, err error) {
	_model = &TrustedProfileTemplateClaimRule{
		Type:       core.StringPtr(typeVar),
		Conditions: conditions,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalTrustedProfileTemplateClaimRule unmarshals an instance of TrustedProfileTemplateClaimRule from the specified map of raw messages.
func UnmarshalTrustedProfileTemplateClaimRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TrustedProfileTemplateClaimRule)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "realm_name", &obj.RealmName)
	if err != nil {
		err = core.SDKErrorf(err, "", "realm_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration", &obj.Expiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalProfileClaimRuleConditions)
	if err != nil {
		err = core.SDKErrorf(err, "", "conditions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TrustedProfileTemplateList : TrustedProfileTemplateList struct
type TrustedProfileTemplateList struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// The offset of the current page.
	Offset *int64 `json:"offset,omitempty"`

	// Optional size of a single page.
	Limit *int64 `json:"limit,omitempty"`

	// Link to the first page.
	First *string `json:"first,omitempty"`

	// Link to the previous available page. If 'previous' property is not part of the response no previous page is
	// available.
	Previous *string `json:"previous,omitempty"`

	// Link to the next available page. If 'next' property is not part of the response no next page is available.
	Next *string `json:"next,omitempty"`

	// List of Profile Templates based on the query paramters and the page size. The profile_templates array is always part
	// of the response but might be empty depending on the query parameter values provided.
	ProfileTemplates []TrustedProfileTemplateResponse `json:"profile_templates" validate:"required"`
}

// UnmarshalTrustedProfileTemplateList unmarshals an instance of TrustedProfileTemplateList from the specified map of raw messages.
func UnmarshalTrustedProfileTemplateList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TrustedProfileTemplateList)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "profile_templates", &obj.ProfileTemplates, UnmarshalTrustedProfileTemplateResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "profile_templates-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TrustedProfileTemplateResponse : Response body format for Trusted Profile Template REST requests.
type TrustedProfileTemplateResponse struct {
	// ID of the the template.
	ID *string `json:"id" validate:"required"`

	// Version of the the template.
	Version *int64 `json:"version" validate:"required"`

	// ID of the account where the template resides.
	AccountID *string `json:"account_id" validate:"required"`

	// The name of the trusted profile template. This is visible only in the enterprise account.
	Name *string `json:"name" validate:"required"`

	// The description of the trusted profile template. Describe the template for enterprise account users.
	Description *string `json:"description,omitempty"`

	// Committed flag determines if the template is ready for assignment.
	Committed *bool `json:"committed,omitempty"`

	// Input body parameters for the TemplateProfileComponent.
	Profile *TemplateProfileComponentResponse `json:"profile,omitempty"`

	// Existing policy templates that you can reference to assign access in the trusted profile component.
	PolicyTemplateReferences []PolicyTemplateReference `json:"policy_template_references,omitempty"`

	ActionControls *ActionControls `json:"action_controls,omitempty"`

	// History of the trusted profile template.
	History []EnityHistoryRecord `json:"history,omitempty"`

	// Entity tag for this templateId-version combination.
	EntityTag *string `json:"entity_tag,omitempty"`

	// Cloud resource name.
	CRN *string `json:"crn,omitempty"`

	// Timestamp of when the template was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// IAMid of the creator.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// Timestamp of when the template was last modified.
	LastModifiedAt *string `json:"last_modified_at,omitempty"`

	// IAMid of the identity that made the latest modification.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`
}

// UnmarshalTrustedProfileTemplateResponse unmarshals an instance of TrustedProfileTemplateResponse from the specified map of raw messages.
func UnmarshalTrustedProfileTemplateResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TrustedProfileTemplateResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "committed", &obj.Committed)
	if err != nil {
		err = core.SDKErrorf(err, "", "committed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "profile", &obj.Profile, UnmarshalTemplateProfileComponentResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policy_template_references", &obj.PolicyTemplateReferences, UnmarshalPolicyTemplateReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_template_references-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "action_controls", &obj.ActionControls, UnmarshalActionControls)
	if err != nil {
		err = core.SDKErrorf(err, "", "action_controls-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TrustedProfilesList : Response body format for the List trusted profiles V1 REST request.
type TrustedProfilesList struct {
	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// The offset of the current page.
	Offset *int64 `json:"offset,omitempty"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Limit *int64 `json:"limit,omitempty"`

	// Link to the first page.
	First *string `json:"first,omitempty"`

	// Link to the previous available page. If 'previous' property is not part of the response no previous page is
	// available.
	Previous *string `json:"previous,omitempty"`

	// Link to the next available page. If 'next' property is not part of the response no next page is available.
	Next *string `json:"next,omitempty"`

	// List of trusted profiles.
	Profiles []TrustedProfile `json:"profiles" validate:"required"`
}

// UnmarshalTrustedProfilesList unmarshals an instance of TrustedProfilesList from the specified map of raw messages.
func UnmarshalTrustedProfilesList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TrustedProfilesList)
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalResponseContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "context-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalTrustedProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "profiles-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UnlockAPIKeyOptions : The UnlockAPIKey options.
type UnlockAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUnlockAPIKeyOptions : Instantiate UnlockAPIKeyOptions
func (*IamIdentityV1) NewUnlockAPIKeyOptions(id string) *UnlockAPIKeyOptions {
	return &UnlockAPIKeyOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UnlockAPIKeyOptions) SetID(id string) *UnlockAPIKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UnlockAPIKeyOptions) SetHeaders(param map[string]string) *UnlockAPIKeyOptions {
	options.Headers = param
	return options
}

// UnlockServiceIDOptions : The UnlockServiceID options.
type UnlockServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUnlockServiceIDOptions : Instantiate UnlockServiceIDOptions
func (*IamIdentityV1) NewUnlockServiceIDOptions(id string) *UnlockServiceIDOptions {
	return &UnlockServiceIDOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UnlockServiceIDOptions) SetID(id string) *UnlockServiceIDOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UnlockServiceIDOptions) SetHeaders(param map[string]string) *UnlockServiceIDOptions {
	options.Headers = param
	return options
}

// UpdateAccountSettingsAssignmentOptions : The UpdateAccountSettingsAssignment options.
type UpdateAccountSettingsAssignmentOptions struct {
	// ID of the Assignment Record.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// Version of the assignment to be updated. Specify the version that you retrieved when reading the assignment. This
	// value  helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might
	// result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Template version to be applied to the assignment. To retry all failed assignments, provide the existing version. To
	// migrate to a different version, provide the new version number.
	TemplateVersion *int64 `json:"template_version" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateAccountSettingsAssignmentOptions : Instantiate UpdateAccountSettingsAssignmentOptions
func (*IamIdentityV1) NewUpdateAccountSettingsAssignmentOptions(assignmentID string, ifMatch string, templateVersion int64) *UpdateAccountSettingsAssignmentOptions {
	return &UpdateAccountSettingsAssignmentOptions{
		AssignmentID:    core.StringPtr(assignmentID),
		IfMatch:         core.StringPtr(ifMatch),
		TemplateVersion: core.Int64Ptr(templateVersion),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *UpdateAccountSettingsAssignmentOptions) SetAssignmentID(assignmentID string) *UpdateAccountSettingsAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateAccountSettingsAssignmentOptions) SetIfMatch(ifMatch string) *UpdateAccountSettingsAssignmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetTemplateVersion : Allow user to set TemplateVersion
func (_options *UpdateAccountSettingsAssignmentOptions) SetTemplateVersion(templateVersion int64) *UpdateAccountSettingsAssignmentOptions {
	_options.TemplateVersion = core.Int64Ptr(templateVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccountSettingsAssignmentOptions) SetHeaders(param map[string]string) *UpdateAccountSettingsAssignmentOptions {
	options.Headers = param
	return options
}

// UpdateAccountSettingsOptions : The UpdateAccountSettings options.
type UpdateAccountSettingsOptions struct {
	// Version of the account settings to be updated. Specify the version that you retrieved as entity_tag (ETag header)
	// when reading the account. This value helps identifying parallel usage of this API. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The id of the account to update the settings for.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Defines whether or not creating a service ID is access controlled. Valid values:
	//   * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
	// IDs, including the account owner
	//   * NOT_RESTRICTED - all members of an account can create service IDs
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreateServiceID *string `json:"restrict_create_service_id,omitempty"`

	// Defines whether or not creating platform API keys is access controlled. Valid values:
	//   * RESTRICTED - only users assigned the 'User API key creator' role on the IAM Identity Service can create API
	// keys, including the account owner
	//   * NOT_RESTRICTED - all members of an account can create platform API keys
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreatePlatformApikey *string `json:"restrict_create_platform_apikey,omitempty"`

	// Defines the IP addresses and subnets from which IAM tokens can be created for the account.
	AllowedIPAddresses *string `json:"allowed_ip_addresses,omitempty"`

	// Defines the MFA trait for the account. Valid values:
	//   * NONE - No MFA trait set
	//   * NONE_NO_ROPC- No MFA, disable CLI logins with only a password
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa,omitempty"`

	// List of users that are exempted from the MFA requirement of the account.
	UserMfa []AccountSettingsUserMfa `json:"user_mfa,omitempty"`

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string `json:"session_expiration_in_seconds,omitempty"`

	// Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string `json:"session_invalidation_in_seconds,omitempty"`

	// Defines the max allowed sessions per identity required by the account. Value values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string `json:"max_sessions_per_identity,omitempty"`

	// Defines the access token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '3600'
	//   * NOT_SET - To unset account setting and use service default.
	SystemAccessTokenExpirationInSeconds *string `json:"system_access_token_expiration_in_seconds,omitempty"`

	// Defines the refresh token expiration in seconds. Valid values:
	//   * Any whole number between '900' and '259200'
	//   * NOT_SET - To unset account setting and use service default.
	SystemRefreshTokenExpirationInSeconds *string `json:"system_refresh_token_expiration_in_seconds,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the UpdateAccountSettingsOptions.RestrictCreateServiceID property.
// Defines whether or not creating a service ID is access controlled. Valid values:
//   - RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service
//
// IDs, including the account owner
//   - NOT_RESTRICTED - all members of an account can create service IDs
//   - NOT_SET - to 'unset' a previous set value.
const (
	UpdateAccountSettingsOptionsRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	UpdateAccountSettingsOptionsRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	UpdateAccountSettingsOptionsRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the UpdateAccountSettingsOptions.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   - RESTRICTED - only users assigned the 'User API key creator' role on the IAM Identity Service can create API keys,
//
// including the account owner
//   - NOT_RESTRICTED - all members of an account can create platform API keys
//   - NOT_SET - to 'unset' a previous set value.
const (
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyNotSetConst        = "NOT_SET"
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyRestrictedConst    = "RESTRICTED"
)

// Constants associated with the UpdateAccountSettingsOptions.Mfa property.
// Defines the MFA trait for the account. Valid values:
//   - NONE - No MFA trait set
//   - NONE_NO_ROPC- No MFA, disable CLI logins with only a password
//   - TOTP - For all non-federated IBMId users
//   - TOTP4ALL - For all users
//   - LEVEL1 - Email-based MFA for all users
//   - LEVEL2 - TOTP-based MFA for all users
//   - LEVEL3 - U2F MFA for all users.
const (
	UpdateAccountSettingsOptionsMfaLevel1Const     = "LEVEL1"
	UpdateAccountSettingsOptionsMfaLevel2Const     = "LEVEL2"
	UpdateAccountSettingsOptionsMfaLevel3Const     = "LEVEL3"
	UpdateAccountSettingsOptionsMfaNoneConst       = "NONE"
	UpdateAccountSettingsOptionsMfaNoneNoRopcConst = "NONE_NO_ROPC"
	UpdateAccountSettingsOptionsMfaTotpConst       = "TOTP"
	UpdateAccountSettingsOptionsMfaTotp4allConst   = "TOTP4ALL"
)

// NewUpdateAccountSettingsOptions : Instantiate UpdateAccountSettingsOptions
func (*IamIdentityV1) NewUpdateAccountSettingsOptions(ifMatch string, accountID string) *UpdateAccountSettingsOptions {
	return &UpdateAccountSettingsOptions{
		IfMatch:   core.StringPtr(ifMatch),
		AccountID: core.StringPtr(accountID),
	}
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateAccountSettingsOptions) SetIfMatch(ifMatch string) *UpdateAccountSettingsOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateAccountSettingsOptions) SetAccountID(accountID string) *UpdateAccountSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetRestrictCreateServiceID : Allow user to set RestrictCreateServiceID
func (_options *UpdateAccountSettingsOptions) SetRestrictCreateServiceID(restrictCreateServiceID string) *UpdateAccountSettingsOptions {
	_options.RestrictCreateServiceID = core.StringPtr(restrictCreateServiceID)
	return _options
}

// SetRestrictCreatePlatformApikey : Allow user to set RestrictCreatePlatformApikey
func (_options *UpdateAccountSettingsOptions) SetRestrictCreatePlatformApikey(restrictCreatePlatformApikey string) *UpdateAccountSettingsOptions {
	_options.RestrictCreatePlatformApikey = core.StringPtr(restrictCreatePlatformApikey)
	return _options
}

// SetAllowedIPAddresses : Allow user to set AllowedIPAddresses
func (_options *UpdateAccountSettingsOptions) SetAllowedIPAddresses(allowedIPAddresses string) *UpdateAccountSettingsOptions {
	_options.AllowedIPAddresses = core.StringPtr(allowedIPAddresses)
	return _options
}

// SetMfa : Allow user to set Mfa
func (_options *UpdateAccountSettingsOptions) SetMfa(mfa string) *UpdateAccountSettingsOptions {
	_options.Mfa = core.StringPtr(mfa)
	return _options
}

// SetUserMfa : Allow user to set UserMfa
func (_options *UpdateAccountSettingsOptions) SetUserMfa(userMfa []AccountSettingsUserMfa) *UpdateAccountSettingsOptions {
	_options.UserMfa = userMfa
	return _options
}

// SetSessionExpirationInSeconds : Allow user to set SessionExpirationInSeconds
func (_options *UpdateAccountSettingsOptions) SetSessionExpirationInSeconds(sessionExpirationInSeconds string) *UpdateAccountSettingsOptions {
	_options.SessionExpirationInSeconds = core.StringPtr(sessionExpirationInSeconds)
	return _options
}

// SetSessionInvalidationInSeconds : Allow user to set SessionInvalidationInSeconds
func (_options *UpdateAccountSettingsOptions) SetSessionInvalidationInSeconds(sessionInvalidationInSeconds string) *UpdateAccountSettingsOptions {
	_options.SessionInvalidationInSeconds = core.StringPtr(sessionInvalidationInSeconds)
	return _options
}

// SetMaxSessionsPerIdentity : Allow user to set MaxSessionsPerIdentity
func (_options *UpdateAccountSettingsOptions) SetMaxSessionsPerIdentity(maxSessionsPerIdentity string) *UpdateAccountSettingsOptions {
	_options.MaxSessionsPerIdentity = core.StringPtr(maxSessionsPerIdentity)
	return _options
}

// SetSystemAccessTokenExpirationInSeconds : Allow user to set SystemAccessTokenExpirationInSeconds
func (_options *UpdateAccountSettingsOptions) SetSystemAccessTokenExpirationInSeconds(systemAccessTokenExpirationInSeconds string) *UpdateAccountSettingsOptions {
	_options.SystemAccessTokenExpirationInSeconds = core.StringPtr(systemAccessTokenExpirationInSeconds)
	return _options
}

// SetSystemRefreshTokenExpirationInSeconds : Allow user to set SystemRefreshTokenExpirationInSeconds
func (_options *UpdateAccountSettingsOptions) SetSystemRefreshTokenExpirationInSeconds(systemRefreshTokenExpirationInSeconds string) *UpdateAccountSettingsOptions {
	_options.SystemRefreshTokenExpirationInSeconds = core.StringPtr(systemRefreshTokenExpirationInSeconds)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccountSettingsOptions) SetHeaders(param map[string]string) *UpdateAccountSettingsOptions {
	options.Headers = param
	return options
}

// UpdateAccountSettingsTemplateVersionOptions : The UpdateAccountSettingsTemplateVersion options.
type UpdateAccountSettingsTemplateVersionOptions struct {
	// Entity tag of the Template to be updated. Specify the tag that you retrieved when reading the account settings
	// template. This value helps identifying parallel usage of this API. Pass * to indicate to update any version
	// available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// ID of the account settings template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Version of the account settings template.
	Version *string `json:"version" validate:"required,ne="`

	// ID of the account where the template resides.
	AccountID *string `json:"account_id,omitempty"`

	// The name of the trusted profile template. This is visible only in the enterprise account.
	Name *string `json:"name,omitempty"`

	// The description of the trusted profile template. Describe the template for enterprise account users.
	Description *string `json:"description,omitempty"`

	AccountSettings *AccountSettingsComponent `json:"account_settings,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateAccountSettingsTemplateVersionOptions : Instantiate UpdateAccountSettingsTemplateVersionOptions
func (*IamIdentityV1) NewUpdateAccountSettingsTemplateVersionOptions(ifMatch string, templateID string, version string) *UpdateAccountSettingsTemplateVersionOptions {
	return &UpdateAccountSettingsTemplateVersionOptions{
		IfMatch:    core.StringPtr(ifMatch),
		TemplateID: core.StringPtr(templateID),
		Version:    core.StringPtr(version),
	}
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateAccountSettingsTemplateVersionOptions) SetIfMatch(ifMatch string) *UpdateAccountSettingsTemplateVersionOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetTemplateID : Allow user to set TemplateID
func (_options *UpdateAccountSettingsTemplateVersionOptions) SetTemplateID(templateID string) *UpdateAccountSettingsTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *UpdateAccountSettingsTemplateVersionOptions) SetVersion(version string) *UpdateAccountSettingsTemplateVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateAccountSettingsTemplateVersionOptions) SetAccountID(accountID string) *UpdateAccountSettingsTemplateVersionOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAccountSettingsTemplateVersionOptions) SetName(name string) *UpdateAccountSettingsTemplateVersionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateAccountSettingsTemplateVersionOptions) SetDescription(description string) *UpdateAccountSettingsTemplateVersionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAccountSettings : Allow user to set AccountSettings
func (_options *UpdateAccountSettingsTemplateVersionOptions) SetAccountSettings(accountSettings *AccountSettingsComponent) *UpdateAccountSettingsTemplateVersionOptions {
	_options.AccountSettings = accountSettings
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccountSettingsTemplateVersionOptions) SetHeaders(param map[string]string) *UpdateAccountSettingsTemplateVersionOptions {
	options.Headers = param
	return options
}

// UpdateAPIKeyOptions : The UpdateAPIKey options.
type UpdateAPIKeyOptions struct {
	// Unique ID of the API key to be updated.
	ID *string `json:"id" validate:"required,ne="`

	// Version of the API key to be updated. Specify the version that you retrieved when reading the API key. This value
	// helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might result
	// in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The name of the API key to update. If specified in the request the parameter must not be empty. The name is not
	// checked for uniqueness. Failure to this will result in an Error condition.
	Name *string `json:"name,omitempty"`

	// The description of the API key to update. If specified an empty description will clear the description of the API
	// key. If a non empty value is provided the API key will be updated.
	Description *string `json:"description,omitempty"`

	// Defines if the API key supports sessions. Sessions are only supported for user apikeys.
	SupportSessions *bool `json:"support_sessions,omitempty"`

	// Defines the action to take when API key is leaked, valid values are 'none', 'disable' and 'delete'.
	ActionWhenLeaked *string `json:"action_when_leaked,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateAPIKeyOptions : Instantiate UpdateAPIKeyOptions
func (*IamIdentityV1) NewUpdateAPIKeyOptions(id string, ifMatch string) *UpdateAPIKeyOptions {
	return &UpdateAPIKeyOptions{
		ID:      core.StringPtr(id),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateAPIKeyOptions) SetID(id string) *UpdateAPIKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateAPIKeyOptions) SetIfMatch(ifMatch string) *UpdateAPIKeyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAPIKeyOptions) SetName(name string) *UpdateAPIKeyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateAPIKeyOptions) SetDescription(description string) *UpdateAPIKeyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetSupportSessions : Allow user to set SupportSessions
func (_options *UpdateAPIKeyOptions) SetSupportSessions(supportSessions bool) *UpdateAPIKeyOptions {
	_options.SupportSessions = core.BoolPtr(supportSessions)
	return _options
}

// SetActionWhenLeaked : Allow user to set ActionWhenLeaked
func (_options *UpdateAPIKeyOptions) SetActionWhenLeaked(actionWhenLeaked string) *UpdateAPIKeyOptions {
	_options.ActionWhenLeaked = core.StringPtr(actionWhenLeaked)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAPIKeyOptions) SetHeaders(param map[string]string) *UpdateAPIKeyOptions {
	options.Headers = param
	return options
}

// UpdateClaimRuleOptions : The UpdateClaimRule options.
type UpdateClaimRuleOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// ID of the claim rule to update.
	RuleID *string `json:"rule-id" validate:"required,ne="`

	// Version of the claim rule to be updated. Specify the version that you retrived when reading list of claim rules.
	// This value helps to identify any parallel usage of claim rule. Pass * to indicate to update any version available.
	// This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Type of the claim rule, either 'Profile-SAML' or 'Profile-CR'.
	Type *string `json:"type" validate:"required"`

	// Conditions of this claim rule.
	Conditions []ProfileClaimRuleConditions `json:"conditions" validate:"required"`

	// Context with key properties for problem determination.
	Context *ResponseContext `json:"context,omitempty"`

	// Name of the claim rule to be created or updated.
	Name *string `json:"name,omitempty"`

	// The realm name of the Idp this claim rule applies to. This field is required only if the type is specified as
	// 'Profile-SAML'.
	RealmName *string `json:"realm_name,omitempty"`

	// The compute resource type the rule applies to, required only if type is specified as 'Profile-CR'. Valid values are
	// VSI, IKS_SA, ROKS_SA.
	CrType *string `json:"cr_type,omitempty"`

	// Session expiration in seconds, only required if type is 'Profile-SAML'.
	Expiration *int64 `json:"expiration,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateClaimRuleOptions : Instantiate UpdateClaimRuleOptions
func (*IamIdentityV1) NewUpdateClaimRuleOptions(profileID string, ruleID string, ifMatch string, typeVar string, conditions []ProfileClaimRuleConditions) *UpdateClaimRuleOptions {
	return &UpdateClaimRuleOptions{
		ProfileID:  core.StringPtr(profileID),
		RuleID:     core.StringPtr(ruleID),
		IfMatch:    core.StringPtr(ifMatch),
		Type:       core.StringPtr(typeVar),
		Conditions: conditions,
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *UpdateClaimRuleOptions) SetProfileID(profileID string) *UpdateClaimRuleOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *UpdateClaimRuleOptions) SetRuleID(ruleID string) *UpdateClaimRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateClaimRuleOptions) SetIfMatch(ifMatch string) *UpdateClaimRuleOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetType : Allow user to set Type
func (_options *UpdateClaimRuleOptions) SetType(typeVar string) *UpdateClaimRuleOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetConditions : Allow user to set Conditions
func (_options *UpdateClaimRuleOptions) SetConditions(conditions []ProfileClaimRuleConditions) *UpdateClaimRuleOptions {
	_options.Conditions = conditions
	return _options
}

// SetContext : Allow user to set Context
func (_options *UpdateClaimRuleOptions) SetContext(context *ResponseContext) *UpdateClaimRuleOptions {
	_options.Context = context
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateClaimRuleOptions) SetName(name string) *UpdateClaimRuleOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetRealmName : Allow user to set RealmName
func (_options *UpdateClaimRuleOptions) SetRealmName(realmName string) *UpdateClaimRuleOptions {
	_options.RealmName = core.StringPtr(realmName)
	return _options
}

// SetCrType : Allow user to set CrType
func (_options *UpdateClaimRuleOptions) SetCrType(crType string) *UpdateClaimRuleOptions {
	_options.CrType = core.StringPtr(crType)
	return _options
}

// SetExpiration : Allow user to set Expiration
func (_options *UpdateClaimRuleOptions) SetExpiration(expiration int64) *UpdateClaimRuleOptions {
	_options.Expiration = core.Int64Ptr(expiration)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateClaimRuleOptions) SetHeaders(param map[string]string) *UpdateClaimRuleOptions {
	options.Headers = param
	return options
}

// UpdateProfileOptions : The UpdateProfile options.
type UpdateProfileOptions struct {
	// ID of the trusted profile to be updated.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// Version of the trusted profile to be updated. Specify the version that you retrived when reading list of trusted
	// profiles. This value helps to identify any parallel usage of trusted profile. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The name of the trusted profile to update. If specified in the request the parameter must not be empty. The name is
	// checked for uniqueness. Failure to this will result in an Error condition.
	Name *string `json:"name,omitempty"`

	// The description of the trusted profile to update. If specified an empty description will clear the description of
	// the trusted profile. If a non empty value is provided the trusted profile will be updated.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateProfileOptions : Instantiate UpdateProfileOptions
func (*IamIdentityV1) NewUpdateProfileOptions(profileID string, ifMatch string) *UpdateProfileOptions {
	return &UpdateProfileOptions{
		ProfileID: core.StringPtr(profileID),
		IfMatch:   core.StringPtr(ifMatch),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *UpdateProfileOptions) SetProfileID(profileID string) *UpdateProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateProfileOptions) SetIfMatch(ifMatch string) *UpdateProfileOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateProfileOptions) SetName(name string) *UpdateProfileOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateProfileOptions) SetDescription(description string) *UpdateProfileOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateProfileOptions) SetHeaders(param map[string]string) *UpdateProfileOptions {
	options.Headers = param
	return options
}

// UpdateProfileTemplateVersionOptions : The UpdateProfileTemplateVersion options.
type UpdateProfileTemplateVersionOptions struct {
	// Entity tag of the Template to be updated. Specify the tag that you retrieved when reading the Profile Template. This
	// value helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might
	// result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// ID of the trusted profile template.
	TemplateID *string `json:"template_id" validate:"required,ne="`

	// Version of the Profile Template.
	Version *string `json:"version" validate:"required,ne="`

	// ID of the account where the template resides.
	AccountID *string `json:"account_id,omitempty"`

	// The name of the trusted profile template. This is visible only in the enterprise account. Required field when
	// creating a new template. Otherwise this field is optional. If the field is included it will change the name value
	// for all existing versions of the template.
	Name *string `json:"name,omitempty"`

	// The description of the trusted profile template. Describe the template for enterprise account users.
	Description *string `json:"description,omitempty"`

	// Input body parameters for the TemplateProfileComponent.
	Profile *TemplateProfileComponentRequest `json:"profile,omitempty"`

	// Existing policy templates that you can reference to assign access in the trusted profile component.
	PolicyTemplateReferences []PolicyTemplateReference `json:"policy_template_references,omitempty"`

	ActionControls *ActionControls `json:"action_controls,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateProfileTemplateVersionOptions : Instantiate UpdateProfileTemplateVersionOptions
func (*IamIdentityV1) NewUpdateProfileTemplateVersionOptions(ifMatch string, templateID string, version string) *UpdateProfileTemplateVersionOptions {
	return &UpdateProfileTemplateVersionOptions{
		IfMatch:    core.StringPtr(ifMatch),
		TemplateID: core.StringPtr(templateID),
		Version:    core.StringPtr(version),
	}
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateProfileTemplateVersionOptions) SetIfMatch(ifMatch string) *UpdateProfileTemplateVersionOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetTemplateID : Allow user to set TemplateID
func (_options *UpdateProfileTemplateVersionOptions) SetTemplateID(templateID string) *UpdateProfileTemplateVersionOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *UpdateProfileTemplateVersionOptions) SetVersion(version string) *UpdateProfileTemplateVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateProfileTemplateVersionOptions) SetAccountID(accountID string) *UpdateProfileTemplateVersionOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateProfileTemplateVersionOptions) SetName(name string) *UpdateProfileTemplateVersionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateProfileTemplateVersionOptions) SetDescription(description string) *UpdateProfileTemplateVersionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *UpdateProfileTemplateVersionOptions) SetProfile(profile *TemplateProfileComponentRequest) *UpdateProfileTemplateVersionOptions {
	_options.Profile = profile
	return _options
}

// SetPolicyTemplateReferences : Allow user to set PolicyTemplateReferences
func (_options *UpdateProfileTemplateVersionOptions) SetPolicyTemplateReferences(policyTemplateReferences []PolicyTemplateReference) *UpdateProfileTemplateVersionOptions {
	_options.PolicyTemplateReferences = policyTemplateReferences
	return _options
}

// SetActionControls : Allow user to set ActionControls
func (_options *UpdateProfileTemplateVersionOptions) SetActionControls(actionControls *ActionControls) *UpdateProfileTemplateVersionOptions {
	_options.ActionControls = actionControls
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateProfileTemplateVersionOptions) SetHeaders(param map[string]string) *UpdateProfileTemplateVersionOptions {
	options.Headers = param
	return options
}

// UpdateServiceIDOptions : The UpdateServiceID options.
type UpdateServiceIDOptions struct {
	// Unique ID of the service ID to be updated.
	ID *string `json:"id" validate:"required,ne="`

	// Version of the service ID to be updated. Specify the version that you retrieved as entity_tag (ETag header) when
	// reading the service ID. This value helps identifying parallel usage of this API. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The name of the service ID to update. If specified in the request the parameter must not be empty. The name is not
	// checked for uniqueness. Failure to this will result in an Error condition.
	Name *string `json:"name,omitempty"`

	// The description of the service ID to update. If specified an empty description will clear the description of the
	// service ID. If an non empty value is provided the service ID will be updated.
	Description *string `json:"description,omitempty"`

	// List of CRNs which point to the services connected to this service ID. If specified an empty list will clear all
	// existing unique instance crns of the service ID.
	UniqueInstanceCrns []string `json:"unique_instance_crns,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateServiceIDOptions : Instantiate UpdateServiceIDOptions
func (*IamIdentityV1) NewUpdateServiceIDOptions(id string, ifMatch string) *UpdateServiceIDOptions {
	return &UpdateServiceIDOptions{
		ID:      core.StringPtr(id),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateServiceIDOptions) SetID(id string) *UpdateServiceIDOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateServiceIDOptions) SetIfMatch(ifMatch string) *UpdateServiceIDOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateServiceIDOptions) SetName(name string) *UpdateServiceIDOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateServiceIDOptions) SetDescription(description string) *UpdateServiceIDOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetUniqueInstanceCrns : Allow user to set UniqueInstanceCrns
func (_options *UpdateServiceIDOptions) SetUniqueInstanceCrns(uniqueInstanceCrns []string) *UpdateServiceIDOptions {
	_options.UniqueInstanceCrns = uniqueInstanceCrns
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateServiceIDOptions) SetHeaders(param map[string]string) *UpdateServiceIDOptions {
	options.Headers = param
	return options
}

// UpdateTrustedProfileAssignmentOptions : The UpdateTrustedProfileAssignment options.
type UpdateTrustedProfileAssignmentOptions struct {
	// ID of the Assignment Record.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// Version of the Assignment to be updated. Specify the version that you retrieved when reading the Assignment. This
	// value  helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might
	// result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Template version to be applied to the assignment. To retry all failed assignments, provide the existing version. To
	// migrate to a different version, provide the new version number.
	TemplateVersion *int64 `json:"template_version" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateTrustedProfileAssignmentOptions : Instantiate UpdateTrustedProfileAssignmentOptions
func (*IamIdentityV1) NewUpdateTrustedProfileAssignmentOptions(assignmentID string, ifMatch string, templateVersion int64) *UpdateTrustedProfileAssignmentOptions {
	return &UpdateTrustedProfileAssignmentOptions{
		AssignmentID:    core.StringPtr(assignmentID),
		IfMatch:         core.StringPtr(ifMatch),
		TemplateVersion: core.Int64Ptr(templateVersion),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *UpdateTrustedProfileAssignmentOptions) SetAssignmentID(assignmentID string) *UpdateTrustedProfileAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateTrustedProfileAssignmentOptions) SetIfMatch(ifMatch string) *UpdateTrustedProfileAssignmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetTemplateVersion : Allow user to set TemplateVersion
func (_options *UpdateTrustedProfileAssignmentOptions) SetTemplateVersion(templateVersion int64) *UpdateTrustedProfileAssignmentOptions {
	_options.TemplateVersion = core.Int64Ptr(templateVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateTrustedProfileAssignmentOptions) SetHeaders(param map[string]string) *UpdateTrustedProfileAssignmentOptions {
	options.Headers = param
	return options
}

// UserActivity : UserActivity struct
type UserActivity struct {
	// IAMid of the user.
	IamID *string `json:"iam_id" validate:"required"`

	// Name of the user.
	Name *string `json:"name,omitempty"`

	// Username of the user.
	Username *string `json:"username" validate:"required"`

	// Email of the user.
	Email *string `json:"email,omitempty"`

	// Time when the user was last authenticated.
	LastAuthn *string `json:"last_authn,omitempty"`
}

// UnmarshalUserActivity unmarshals an instance of UserActivity from the specified map of raw messages.
func UnmarshalUserActivity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserActivity)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		err = core.SDKErrorf(err, "", "username-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_authn", &obj.LastAuthn)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_authn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UserMfaEnrollments : UserMfaEnrollments struct
type UserMfaEnrollments struct {
	// IAMid of the user.
	IamID *string `json:"iam_id" validate:"required"`

	// currently effective mfa type i.e. id_based_mfa or account_based_mfa.
	EffectiveMfaType *string `json:"effective_mfa_type,omitempty"`

	IDBasedMfa *IDBasedMfaEnrollment `json:"id_based_mfa,omitempty"`

	AccountBasedMfa *AccountBasedMfaEnrollment `json:"account_based_mfa,omitempty"`
}

// UnmarshalUserMfaEnrollments unmarshals an instance of UserMfaEnrollments from the specified map of raw messages.
func UnmarshalUserMfaEnrollments(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserMfaEnrollments)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "effective_mfa_type", &obj.EffectiveMfaType)
	if err != nil {
		err = core.SDKErrorf(err, "", "effective_mfa_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "id_based_mfa", &obj.IDBasedMfa, UnmarshalIDBasedMfaEnrollment)
	if err != nil {
		err = core.SDKErrorf(err, "", "id_based_mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "account_based_mfa", &obj.AccountBasedMfa, UnmarshalAccountBasedMfaEnrollment)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_based_mfa-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UserReportMfaEnrollmentStatus : UserReportMfaEnrollmentStatus struct
type UserReportMfaEnrollmentStatus struct {
	// IAMid of the user.
	IamID *string `json:"iam_id" validate:"required"`

	// Name of the user.
	Name *string `json:"name,omitempty"`

	// Username of the user.
	Username *string `json:"username" validate:"required"`

	// Email of the user.
	Email *string `json:"email,omitempty"`

	Enrollments *MfaEnrollments `json:"enrollments" validate:"required"`
}

// UnmarshalUserReportMfaEnrollmentStatus unmarshals an instance of UserReportMfaEnrollmentStatus from the specified map of raw messages.
func UnmarshalUserReportMfaEnrollmentStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserReportMfaEnrollmentStatus)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		err = core.SDKErrorf(err, "", "username-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "enrollments", &obj.Enrollments, UnmarshalMfaEnrollments)
	if err != nil {
		err = core.SDKErrorf(err, "", "enrollments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
