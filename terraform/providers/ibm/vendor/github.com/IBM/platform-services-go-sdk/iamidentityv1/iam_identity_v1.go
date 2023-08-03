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
 * IBM OpenAPI SDK Code Generator Version: 3.72.0-5d70f2bb-20230511-203609
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
			return
		}
	}

	iamIdentity, err = NewIamIdentityV1(options)
	if err != nil {
		return
	}

	err = iamIdentity.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = iamIdentity.Service.SetServiceURL(options.URL)
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
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
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
	return "", fmt.Errorf("service does not support regional URLs")
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
	return iamIdentity.Service.SetServiceURL(url)
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
// for themself, or service ID API keys for service IDs that are bound to an entity they have access to. In case of
// service IDs and their API keys, a user must be either an account owner, a IBM Cloud org manager or IBM Cloud space
// developer in order to manage service IDs of the entity.
func (iamIdentity *IamIdentityV1) ListAPIKeys(listAPIKeysOptions *ListAPIKeysOptions) (result *APIKeyList, response *core.DetailedResponse, err error) {
	return iamIdentity.ListAPIKeysWithContext(context.Background(), listAPIKeysOptions)
}

// ListAPIKeysWithContext is an alternate form of the ListAPIKeys method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListAPIKeysWithContext(ctx context.Context, listAPIKeysOptions *ListAPIKeysOptions) (result *APIKeyList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAPIKeysOptions, "listAPIKeysOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys`, nil)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKeyList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAPIKey : Create an API key
// Creates an API key for a UserID or service ID. Users can manage user API keys for themself, or service ID API keys
// for service IDs that are bound to an entity they have access to.
func (iamIdentity *IamIdentityV1) CreateAPIKey(createAPIKeyOptions *CreateAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	return iamIdentity.CreateAPIKeyWithContext(context.Background(), createAPIKeyOptions)
}

// CreateAPIKeyWithContext is an alternate form of the CreateAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateAPIKeyWithContext(ctx context.Context, createAPIKeyOptions *CreateAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAPIKeyOptions, "createAPIKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createAPIKeyOptions, "createAPIKeyOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys`, nil)
	if err != nil {
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
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAPIKeysDetails : Get details of an API key by its value
// Returns the details of an API key by its value. Users can manage user API keys for themself, or service ID API keys
// for service IDs that are bound to an entity they have access to.
func (iamIdentity *IamIdentityV1) GetAPIKeysDetails(getAPIKeysDetailsOptions *GetAPIKeysDetailsOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	return iamIdentity.GetAPIKeysDetailsWithContext(context.Background(), getAPIKeysDetailsOptions)
}

// GetAPIKeysDetailsWithContext is an alternate form of the GetAPIKeysDetails method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetAPIKeysDetailsWithContext(ctx context.Context, getAPIKeysDetailsOptions *GetAPIKeysDetailsOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getAPIKeysDetailsOptions, "getAPIKeysDetailsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/apikeys/details`, nil)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAPIKey : Get details of an API key
// Returns the details of an API key. Users can manage user API keys for themself, or service ID API keys for service
// IDs that are bound to an entity they have access to. In case of service IDs and their API keys, a user must be either
// an account owner, a IBM Cloud org manager or IBM Cloud space developer in order to manage service IDs of the entity.
func (iamIdentity *IamIdentityV1) GetAPIKey(getAPIKeyOptions *GetAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	return iamIdentity.GetAPIKeyWithContext(context.Background(), getAPIKeyOptions)
}

// GetAPIKeyWithContext is an alternate form of the GetAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetAPIKeyWithContext(ctx context.Context, getAPIKeyOptions *GetAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAPIKeyOptions, "getAPIKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAPIKeyOptions, "getAPIKeyOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateAPIKey : Updates an API key
// Updates properties of an API key. This does NOT affect existing access tokens. Their token content will stay
// unchanged until the access token is refreshed. To update an API key, pass the property to be modified. To delete one
// property's value, pass the property with an empty value "".Users can manage user API keys for themself, or service ID
// API keys for service IDs that are bound to an entity they have access to.
func (iamIdentity *IamIdentityV1) UpdateAPIKey(updateAPIKeyOptions *UpdateAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	return iamIdentity.UpdateAPIKeyWithContext(context.Background(), updateAPIKeyOptions)
}

// UpdateAPIKeyWithContext is an alternate form of the UpdateAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateAPIKeyWithContext(ctx context.Context, updateAPIKeyOptions *UpdateAPIKeyOptions) (result *APIKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAPIKeyOptions, "updateAPIKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAPIKeyOptions, "updateAPIKeyOptions")
	if err != nil {
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
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteAPIKey : Deletes an API key
// Deletes an API key. Existing tokens will remain valid until expired. Users can manage user API keys for themself, or
// service ID API keys for service IDs that are bound to an entity they have access to.
func (iamIdentity *IamIdentityV1) DeleteAPIKey(deleteAPIKeyOptions *DeleteAPIKeyOptions) (response *core.DetailedResponse, err error) {
	return iamIdentity.DeleteAPIKeyWithContext(context.Background(), deleteAPIKeyOptions)
}

// DeleteAPIKeyWithContext is an alternate form of the DeleteAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteAPIKeyWithContext(ctx context.Context, deleteAPIKeyOptions *DeleteAPIKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAPIKeyOptions, "deleteAPIKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteAPIKeyOptions, "deleteAPIKeyOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// LockAPIKey : Lock the API key
// Locks an API key by ID. Users can manage user API keys for themself, or service ID API keys for service IDs that are
// bound to an entity they have access to. In case of service IDs and their API keys, a user must be either an account
// owner, a IBM Cloud org manager or IBM Cloud space developer in order to manage service IDs of the entity.
func (iamIdentity *IamIdentityV1) LockAPIKey(lockAPIKeyOptions *LockAPIKeyOptions) (response *core.DetailedResponse, err error) {
	return iamIdentity.LockAPIKeyWithContext(context.Background(), lockAPIKeyOptions)
}

// LockAPIKeyWithContext is an alternate form of the LockAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) LockAPIKeyWithContext(ctx context.Context, lockAPIKeyOptions *LockAPIKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(lockAPIKeyOptions, "lockAPIKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(lockAPIKeyOptions, "lockAPIKeyOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// UnlockAPIKey : Unlock the API key
// Unlocks an API key by ID. Users can manage user API keys for themself, or service ID API keys for service IDs that
// are bound to an entity they have access to. In case of service IDs and their API keys, a user must be either an
// account owner, a IBM Cloud org manager or IBM Cloud space developer in order to manage service IDs of the entity.
func (iamIdentity *IamIdentityV1) UnlockAPIKey(unlockAPIKeyOptions *UnlockAPIKeyOptions) (response *core.DetailedResponse, err error) {
	return iamIdentity.UnlockAPIKeyWithContext(context.Background(), unlockAPIKeyOptions)
}

// UnlockAPIKeyWithContext is an alternate form of the UnlockAPIKey method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UnlockAPIKeyWithContext(ctx context.Context, unlockAPIKeyOptions *UnlockAPIKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(unlockAPIKeyOptions, "unlockAPIKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(unlockAPIKeyOptions, "unlockAPIKeyOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// ListServiceIds : List service IDs
// Returns a list of service IDs. Users can manage user API keys for themself, or service ID API keys for service IDs
// that are bound to an entity they have access to. Note: apikey details are only included in the response when creating
// a Service ID with an api key.
func (iamIdentity *IamIdentityV1) ListServiceIds(listServiceIdsOptions *ListServiceIdsOptions) (result *ServiceIDList, response *core.DetailedResponse, err error) {
	return iamIdentity.ListServiceIdsWithContext(context.Background(), listServiceIdsOptions)
}

// ListServiceIdsWithContext is an alternate form of the ListServiceIds method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListServiceIdsWithContext(ctx context.Context, listServiceIdsOptions *ListServiceIdsOptions) (result *ServiceIDList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listServiceIdsOptions, "listServiceIdsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/`, nil)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceIDList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateServiceID : Create a service ID
// Creates a service ID for an IBM Cloud account. Users can manage user API keys for themself, or service ID API keys
// for service IDs that are bound to an entity they have access to.
func (iamIdentity *IamIdentityV1) CreateServiceID(createServiceIDOptions *CreateServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	return iamIdentity.CreateServiceIDWithContext(context.Background(), createServiceIDOptions)
}

// CreateServiceIDWithContext is an alternate form of the CreateServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateServiceIDWithContext(ctx context.Context, createServiceIDOptions *CreateServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createServiceIDOptions, "createServiceIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createServiceIDOptions, "createServiceIDOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/serviceids/`, nil)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetServiceID : Get details of a service ID
// Returns the details of a service ID. Users can manage user API keys for themself, or service ID API keys for service
// IDs that are bound to an entity they have access to. Note: apikey details are only included in the response when
// creating a Service ID with an api key.
func (iamIdentity *IamIdentityV1) GetServiceID(getServiceIDOptions *GetServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	return iamIdentity.GetServiceIDWithContext(context.Background(), getServiceIDOptions)
}

// GetServiceIDWithContext is an alternate form of the GetServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetServiceIDWithContext(ctx context.Context, getServiceIDOptions *GetServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getServiceIDOptions, "getServiceIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getServiceIDOptions, "getServiceIDOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
		if err != nil {
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
// service ID API keys for service IDs that are bound to an entity they have access to. Note: apikey details are only
// included in the response when creating a Service ID with an apikey.
func (iamIdentity *IamIdentityV1) UpdateServiceID(updateServiceIDOptions *UpdateServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	return iamIdentity.UpdateServiceIDWithContext(context.Background(), updateServiceIDOptions)
}

// UpdateServiceIDWithContext is an alternate form of the UpdateServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateServiceIDWithContext(ctx context.Context, updateServiceIDOptions *UpdateServiceIDOptions) (result *ServiceID, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateServiceIDOptions, "updateServiceIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateServiceIDOptions, "updateServiceIDOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
		if err != nil {
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
// service ID API keys for service IDs that are bound to an entity they have access to.
func (iamIdentity *IamIdentityV1) DeleteServiceID(deleteServiceIDOptions *DeleteServiceIDOptions) (response *core.DetailedResponse, err error) {
	return iamIdentity.DeleteServiceIDWithContext(context.Background(), deleteServiceIDOptions)
}

// DeleteServiceIDWithContext is an alternate form of the DeleteServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteServiceIDWithContext(ctx context.Context, deleteServiceIDOptions *DeleteServiceIDOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteServiceIDOptions, "deleteServiceIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteServiceIDOptions, "deleteServiceIDOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// LockServiceID : Lock the service ID
// Locks a service ID by ID. Users can manage user API keys for themself, or service ID API keys for service IDs that
// are bound to an entity they have access to. In case of service IDs and their API keys, a user must be either an
// account owner, a IBM Cloud org manager or IBM Cloud space developer in order to manage service IDs of the entity.
func (iamIdentity *IamIdentityV1) LockServiceID(lockServiceIDOptions *LockServiceIDOptions) (response *core.DetailedResponse, err error) {
	return iamIdentity.LockServiceIDWithContext(context.Background(), lockServiceIDOptions)
}

// LockServiceIDWithContext is an alternate form of the LockServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) LockServiceIDWithContext(ctx context.Context, lockServiceIDOptions *LockServiceIDOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(lockServiceIDOptions, "lockServiceIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(lockServiceIDOptions, "lockServiceIDOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// UnlockServiceID : Unlock the service ID
// Unlocks a service ID by ID. Users can manage user API keys for themself, or service ID API keys for service IDs that
// are bound to an entity they have access to. In case of service IDs and their API keys, a user must be either an
// account owner, a IBM Cloud org manager or IBM Cloud space developer in order to manage service IDs of the entity.
func (iamIdentity *IamIdentityV1) UnlockServiceID(unlockServiceIDOptions *UnlockServiceIDOptions) (response *core.DetailedResponse, err error) {
	return iamIdentity.UnlockServiceIDWithContext(context.Background(), unlockServiceIDOptions)
}

// UnlockServiceIDWithContext is an alternate form of the UnlockServiceID method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UnlockServiceIDWithContext(ctx context.Context, unlockServiceIDOptions *UnlockServiceIDOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(unlockServiceIDOptions, "unlockServiceIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(unlockServiceIDOptions, "unlockServiceIDOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// CreateProfile : Create a trusted profile
// Create a trusted profile for a given account ID.
func (iamIdentity *IamIdentityV1) CreateProfile(createProfileOptions *CreateProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	return iamIdentity.CreateProfileWithContext(context.Background(), createProfileOptions)
}

// CreateProfileWithContext is an alternate form of the CreateProfile method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateProfileWithContext(ctx context.Context, createProfileOptions *CreateProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
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
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles`, nil)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfile)
		if err != nil {
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
	return iamIdentity.ListProfilesWithContext(context.Background(), listProfilesOptions)
}

// ListProfilesWithContext is an alternate form of the ListProfiles method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListProfilesWithContext(ctx context.Context, listProfilesOptions *ListProfilesOptions) (result *TrustedProfilesList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProfilesOptions, "listProfilesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listProfilesOptions, "listProfilesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/profiles`, nil)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfilesList)
		if err != nil {
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
	return iamIdentity.GetProfileWithContext(context.Background(), getProfileOptions)
}

// GetProfileWithContext is an alternate form of the GetProfile method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetProfileWithContext(ctx context.Context, getProfileOptions *GetProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileOptions, "getProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileOptions, "getProfileOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateProfile : Update a trusted profile
// Update the name or description of an existing trusted profile.
func (iamIdentity *IamIdentityV1) UpdateProfile(updateProfileOptions *UpdateProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	return iamIdentity.UpdateProfileWithContext(context.Background(), updateProfileOptions)
}

// UpdateProfileWithContext is an alternate form of the UpdateProfile method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateProfileWithContext(ctx context.Context, updateProfileOptions *UpdateProfileOptions) (result *TrustedProfile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProfileOptions, "updateProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateProfileOptions, "updateProfileOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustedProfile)
		if err != nil {
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
	return iamIdentity.DeleteProfileWithContext(context.Background(), deleteProfileOptions)
}

// DeleteProfileWithContext is an alternate form of the DeleteProfile method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteProfileWithContext(ctx context.Context, deleteProfileOptions *DeleteProfileOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProfileOptions, "deleteProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProfileOptions, "deleteProfileOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// CreateClaimRule : Create claim rule for a trusted profile
// Create a claim rule for a trusted profile. There is a limit of 20 rules per trusted profile.
func (iamIdentity *IamIdentityV1) CreateClaimRule(createClaimRuleOptions *CreateClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	return iamIdentity.CreateClaimRuleWithContext(context.Background(), createClaimRuleOptions)
}

// CreateClaimRuleWithContext is an alternate form of the CreateClaimRule method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateClaimRuleWithContext(ctx context.Context, createClaimRuleOptions *CreateClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createClaimRuleOptions, "createClaimRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createClaimRuleOptions, "createClaimRuleOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileClaimRule)
		if err != nil {
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
	return iamIdentity.ListClaimRulesWithContext(context.Background(), listClaimRulesOptions)
}

// ListClaimRulesWithContext is an alternate form of the ListClaimRules method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListClaimRulesWithContext(ctx context.Context, listClaimRulesOptions *ListClaimRulesOptions) (result *ProfileClaimRuleList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listClaimRulesOptions, "listClaimRulesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listClaimRulesOptions, "listClaimRulesOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileClaimRuleList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetClaimRule : Get a claim rule for a trusted profile
// A specific claim rule can be fetched for a given trusted profile ID and rule ID.
func (iamIdentity *IamIdentityV1) GetClaimRule(getClaimRuleOptions *GetClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	return iamIdentity.GetClaimRuleWithContext(context.Background(), getClaimRuleOptions)
}

// GetClaimRuleWithContext is an alternate form of the GetClaimRule method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetClaimRuleWithContext(ctx context.Context, getClaimRuleOptions *GetClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getClaimRuleOptions, "getClaimRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getClaimRuleOptions, "getClaimRuleOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileClaimRule)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateClaimRule : Update claim rule for a trusted profile
// Update a specific claim rule for a given trusted profile ID and rule ID.
func (iamIdentity *IamIdentityV1) UpdateClaimRule(updateClaimRuleOptions *UpdateClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	return iamIdentity.UpdateClaimRuleWithContext(context.Background(), updateClaimRuleOptions)
}

// UpdateClaimRuleWithContext is an alternate form of the UpdateClaimRule method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateClaimRuleWithContext(ctx context.Context, updateClaimRuleOptions *UpdateClaimRuleOptions) (result *ProfileClaimRule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateClaimRuleOptions, "updateClaimRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateClaimRuleOptions, "updateClaimRuleOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileClaimRule)
		if err != nil {
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
	return iamIdentity.DeleteClaimRuleWithContext(context.Background(), deleteClaimRuleOptions)
}

// DeleteClaimRuleWithContext is an alternate form of the DeleteClaimRule method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteClaimRuleWithContext(ctx context.Context, deleteClaimRuleOptions *DeleteClaimRuleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteClaimRuleOptions, "deleteClaimRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteClaimRuleOptions, "deleteClaimRuleOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// CreateLink : Create link to a trusted profile
// Create a direct link between a specific compute resource and a trusted profile, rather than creating conditions that
// a compute resource must fulfill to apply a trusted profile.
func (iamIdentity *IamIdentityV1) CreateLink(createLinkOptions *CreateLinkOptions) (result *ProfileLink, response *core.DetailedResponse, err error) {
	return iamIdentity.CreateLinkWithContext(context.Background(), createLinkOptions)
}

// CreateLinkWithContext is an alternate form of the CreateLink method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateLinkWithContext(ctx context.Context, createLinkOptions *CreateLinkOptions) (result *ProfileLink, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createLinkOptions, "createLinkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createLinkOptions, "createLinkOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileLink)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListLinks : List links to a trusted profile
// Get a list of links to a trusted profile.
func (iamIdentity *IamIdentityV1) ListLinks(listLinksOptions *ListLinksOptions) (result *ProfileLinkList, response *core.DetailedResponse, err error) {
	return iamIdentity.ListLinksWithContext(context.Background(), listLinksOptions)
}

// ListLinksWithContext is an alternate form of the ListLinks method which supports a Context parameter
func (iamIdentity *IamIdentityV1) ListLinksWithContext(ctx context.Context, listLinksOptions *ListLinksOptions) (result *ProfileLinkList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listLinksOptions, "listLinksOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listLinksOptions, "listLinksOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileLinkList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLink : Get link to a trusted profile
// Get a specific link to a trusted profile by `link_id`.
func (iamIdentity *IamIdentityV1) GetLink(getLinkOptions *GetLinkOptions) (result *ProfileLink, response *core.DetailedResponse, err error) {
	return iamIdentity.GetLinkWithContext(context.Background(), getLinkOptions)
}

// GetLinkWithContext is an alternate form of the GetLink method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetLinkWithContext(ctx context.Context, getLinkOptions *GetLinkOptions) (result *ProfileLink, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLinkOptions, "getLinkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLinkOptions, "getLinkOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileLink)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteLink : Delete link to a trusted profile
// Delete a link between a compute resource and a trusted profile.
func (iamIdentity *IamIdentityV1) DeleteLink(deleteLinkOptions *DeleteLinkOptions) (response *core.DetailedResponse, err error) {
	return iamIdentity.DeleteLinkWithContext(context.Background(), deleteLinkOptions)
}

// DeleteLinkWithContext is an alternate form of the DeleteLink method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteLinkWithContext(ctx context.Context, deleteLinkOptions *DeleteLinkOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLinkOptions, "deleteLinkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLinkOptions, "deleteLinkOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// GetProfileIdentities : Get a list of identities that can assume the trusted profile
// Get a list of identities that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) GetProfileIdentities(getProfileIdentitiesOptions *GetProfileIdentitiesOptions) (result *ProfileIdentitiesResponse, response *core.DetailedResponse, err error) {
	return iamIdentity.GetProfileIdentitiesWithContext(context.Background(), getProfileIdentitiesOptions)
}

// GetProfileIdentitiesWithContext is an alternate form of the GetProfileIdentities method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetProfileIdentitiesWithContext(ctx context.Context, getProfileIdentitiesOptions *GetProfileIdentitiesOptions) (result *ProfileIdentitiesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileIdentitiesOptions, "getProfileIdentitiesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileIdentitiesOptions, "getProfileIdentitiesOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileIdentitiesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetProfileIdentities : Update the list of identities that can assume the trusted profile
// Update the list of identities that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) SetProfileIdentities(setProfileIdentitiesOptions *SetProfileIdentitiesOptions) (result *ProfileIdentitiesResponse, response *core.DetailedResponse, err error) {
	return iamIdentity.SetProfileIdentitiesWithContext(context.Background(), setProfileIdentitiesOptions)
}

// SetProfileIdentitiesWithContext is an alternate form of the SetProfileIdentities method which supports a Context parameter
func (iamIdentity *IamIdentityV1) SetProfileIdentitiesWithContext(ctx context.Context, setProfileIdentitiesOptions *SetProfileIdentitiesOptions) (result *ProfileIdentitiesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setProfileIdentitiesOptions, "setProfileIdentitiesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setProfileIdentitiesOptions, "setProfileIdentitiesOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileIdentitiesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetProfileIdentity : Add a specific identity that can assume the trusted profile
// Add a specific identity that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) SetProfileIdentity(setProfileIdentityOptions *SetProfileIdentityOptions) (result *ProfileIdentity, response *core.DetailedResponse, err error) {
	return iamIdentity.SetProfileIdentityWithContext(context.Background(), setProfileIdentityOptions)
}

// SetProfileIdentityWithContext is an alternate form of the SetProfileIdentity method which supports a Context parameter
func (iamIdentity *IamIdentityV1) SetProfileIdentityWithContext(ctx context.Context, setProfileIdentityOptions *SetProfileIdentityOptions) (result *ProfileIdentity, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setProfileIdentityOptions, "setProfileIdentityOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setProfileIdentityOptions, "setProfileIdentityOptions")
	if err != nil {
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
	if setProfileIdentityOptions.IamID != nil {
		body["iam_id"] = setProfileIdentityOptions.IamID
	}
	if setProfileIdentityOptions.Accounts != nil {
		body["accounts"] = setProfileIdentityOptions.Accounts
	}
	if setProfileIdentityOptions.Description != nil {
		body["description"] = setProfileIdentityOptions.Description
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
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileIdentity)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProfileIdentity : Get the identity that can assume the trusted profile
// Get the identity that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) GetProfileIdentity(getProfileIdentityOptions *GetProfileIdentityOptions) (result *ProfileIdentity, response *core.DetailedResponse, err error) {
	return iamIdentity.GetProfileIdentityWithContext(context.Background(), getProfileIdentityOptions)
}

// GetProfileIdentityWithContext is an alternate form of the GetProfileIdentity method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetProfileIdentityWithContext(ctx context.Context, getProfileIdentityOptions *GetProfileIdentityOptions) (result *ProfileIdentity, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileIdentityOptions, "getProfileIdentityOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileIdentityOptions, "getProfileIdentityOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileIdentity)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteProfileIdentity : Delete the identity that can assume the trusted profile
// Delete the identity that can assume the trusted profile.
func (iamIdentity *IamIdentityV1) DeleteProfileIdentity(deleteProfileIdentityOptions *DeleteProfileIdentityOptions) (response *core.DetailedResponse, err error) {
	return iamIdentity.DeleteProfileIdentityWithContext(context.Background(), deleteProfileIdentityOptions)
}

// DeleteProfileIdentityWithContext is an alternate form of the DeleteProfileIdentity method which supports a Context parameter
func (iamIdentity *IamIdentityV1) DeleteProfileIdentityWithContext(ctx context.Context, deleteProfileIdentityOptions *DeleteProfileIdentityOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProfileIdentityOptions, "deleteProfileIdentityOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProfileIdentityOptions, "deleteProfileIdentityOptions")
	if err != nil {
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
		return
	}

	response, err = iamIdentity.Service.Request(request, nil)

	return
}

// GetAccountSettings : Get account configurations
// Returns the details of an account's configuration.
func (iamIdentity *IamIdentityV1) GetAccountSettings(getAccountSettingsOptions *GetAccountSettingsOptions) (result *AccountSettingsResponse, response *core.DetailedResponse, err error) {
	return iamIdentity.GetAccountSettingsWithContext(context.Background(), getAccountSettingsOptions)
}

// GetAccountSettingsWithContext is an alternate form of the GetAccountSettings method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetAccountSettingsWithContext(ctx context.Context, getAccountSettingsOptions *GetAccountSettingsOptions) (result *AccountSettingsResponse, response *core.DetailedResponse, err error) {
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
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/accounts/{account_id}/settings/identity`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccountSettings : Update account configurations
// Allows a user to configure settings on their account with regards to MFA, MFA excemption list,  session lifetimes,
// access control for creating new identities, and enforcing IP restrictions on token creation.
func (iamIdentity *IamIdentityV1) UpdateAccountSettings(updateAccountSettingsOptions *UpdateAccountSettingsOptions) (result *AccountSettingsResponse, response *core.DetailedResponse, err error) {
	return iamIdentity.UpdateAccountSettingsWithContext(context.Background(), updateAccountSettingsOptions)
}

// UpdateAccountSettingsWithContext is an alternate form of the UpdateAccountSettings method which supports a Context parameter
func (iamIdentity *IamIdentityV1) UpdateAccountSettingsWithContext(ctx context.Context, updateAccountSettingsOptions *UpdateAccountSettingsOptions) (result *AccountSettingsResponse, response *core.DetailedResponse, err error) {
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

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamIdentity.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamIdentity.Service.Options.URL, `/v1/accounts/{account_id}/settings/identity`, pathParamsMap)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetMfaStatus : Get MFA enrollment status for a single user in the account
// Get MFA enrollment status for a single user in the account.
func (iamIdentity *IamIdentityV1) GetMfaStatus(getMfaStatusOptions *GetMfaStatusOptions) (result *UserMfaEnrollments, response *core.DetailedResponse, err error) {
	return iamIdentity.GetMfaStatusWithContext(context.Background(), getMfaStatusOptions)
}

// GetMfaStatusWithContext is an alternate form of the GetMfaStatus method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetMfaStatusWithContext(ctx context.Context, getMfaStatusOptions *GetMfaStatusOptions) (result *UserMfaEnrollments, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getMfaStatusOptions, "getMfaStatusOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getMfaStatusOptions, "getMfaStatusOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUserMfaEnrollments)
		if err != nil {
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
	return iamIdentity.CreateMfaReportWithContext(context.Background(), createMfaReportOptions)
}

// CreateMfaReportWithContext is an alternate form of the CreateMfaReport method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateMfaReportWithContext(ctx context.Context, createMfaReportOptions *CreateMfaReportOptions) (result *ReportReference, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createMfaReportOptions, "createMfaReportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createMfaReportOptions, "createMfaReportOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportReference)
		if err != nil {
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
	return iamIdentity.GetMfaReportWithContext(context.Background(), getMfaReportOptions)
}

// GetMfaReportWithContext is an alternate form of the GetMfaReport method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetMfaReportWithContext(ctx context.Context, getMfaReportOptions *GetMfaReportOptions) (result *ReportMfaEnrollmentStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getMfaReportOptions, "getMfaReportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getMfaReportOptions, "getMfaReportOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportMfaEnrollmentStatus)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateReport : Trigger activity report for the account
// Trigger activity report for the account by specifying the account ID. It can take a few minutes to generate the
// report for retrieval.
func (iamIdentity *IamIdentityV1) CreateReport(createReportOptions *CreateReportOptions) (result *ReportReference, response *core.DetailedResponse, err error) {
	return iamIdentity.CreateReportWithContext(context.Background(), createReportOptions)
}

// CreateReportWithContext is an alternate form of the CreateReport method which supports a Context parameter
func (iamIdentity *IamIdentityV1) CreateReportWithContext(ctx context.Context, createReportOptions *CreateReportOptions) (result *ReportReference, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createReportOptions, "createReportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createReportOptions, "createReportOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportReference)
		if err != nil {
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
	return iamIdentity.GetReportWithContext(context.Background(), getReportOptions)
}

// GetReportWithContext is an alternate form of the GetReport method which supports a Context parameter
func (iamIdentity *IamIdentityV1) GetReportWithContext(ctx context.Context, getReportOptions *GetReportOptions) (result *Report, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportOptions, "getReportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportOptions, "getReportOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamIdentity.Service.Request(request, &rawResponse)
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
		return
	}
	err = core.UnmarshalModel(m, "totp", &obj.Totp, UnmarshalMfaEnrollmentTypeStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "verisign", &obj.Verisign, UnmarshalMfaEnrollmentTypeStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "complies", &obj.Complies)
	if err != nil {
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

	// Defines whether or not creating a Service Id is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
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
// Defines whether or not creating a Service Id is access controlled. Valid values:
//   - RESTRICTED - to apply access control
//   - NOT_RESTRICTED - to remove access control
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
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_service_id", &obj.RestrictCreateServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "restrict_create_platform_apikey", &obj.RestrictCreatePlatformApikey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_ip_addresses", &obj.AllowedIPAddresses)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "user_mfa", &obj.UserMfa, UnmarshalAccountSettingsUserMfa)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "session_expiration_in_seconds", &obj.SessionExpirationInSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "session_invalidation_in_seconds", &obj.SessionInvalidationInSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_sessions_per_identity", &obj.MaxSessionsPerIdentity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "system_access_token_expiration_in_seconds", &obj.SystemAccessTokenExpirationInSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "system_refresh_token_expiration_in_seconds", &obj.SystemRefreshTokenExpirationInSeconds)
	if err != nil {
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
	return
}

// UnmarshalAccountSettingsUserMfa unmarshals an instance of AccountSettingsUserMfa from the specified map of raw messages.
func UnmarshalAccountSettingsUserMfa(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsUserMfa)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "authn_count", &obj.AuthnCount)
	if err != nil {
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

	// If set contains a date time string of the creation date in ISO format.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// IAM ID of the user or service which created the API key.
	CreatedBy *string `json:"created_by" validate:"required"`

	// If set contains a date time string of the last modification date in ISO format.
	ModifiedAt *strfmt.DateTime `json:"modified_at,omitempty"`

	// Name of the API key. The name is not checked for uniqueness. Therefore multiple names with the same value can exist.
	// Access is done via the UUID of the API key.
	Name *string `json:"name" validate:"required"`

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
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
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
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
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
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "apikey", &obj.Apikey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "activity", &obj.Activity, UnmarshalActivity)
	if err != nil {
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

	// You can optionally passthrough the API key value for this API key. If passed, NO validation of that apiKey value is
	// done, i.e. the value can be non-URL safe. If omitted, the API key management will create an URL safe opaque API key
	// value. The value of the API key is checked for uniqueness. Please ensure enough variations when passing in this
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
	return
}

// UnmarshalAPIKeyInsideCreateServiceIDRequest unmarshals an instance of APIKeyInsideCreateServiceIDRequest from the specified map of raw messages.
func UnmarshalAPIKeyInsideCreateServiceIDRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIKeyInsideCreateServiceIDRequest)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "apikey", &obj.Apikey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "store_value", &obj.StoreValue)
	if err != nil {
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
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "apikeys", &obj.Apikeys, UnmarshalAPIKey)
	if err != nil {
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
	err = core.UnmarshalModel(m, "serviceid", &obj.Serviceid, UnmarshalApikeyActivityServiceid)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "user", &obj.User, UnmarshalApikeyActivityUser)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_authn", &obj.LastAuthn)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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

	// You can optionally passthrough the API key value for this API key. If passed, NO validation of that apiKey value is
	// done, i.e. the value can be non-URL safe. If omitted, the API key management will create an URL safe opaque API key
	// value. The value of the API key is checked for uniqueness. Please ensure enough variations when passing in this
	// value.
	Apikey *string `json:"apikey,omitempty"`

	// Send true or false to set whether the API key value is retrievable in the future by using the Get details of an API
	// key request. If you create an API key for a user, you must specify `false` or omit the value. We don't allow storing
	// of API keys for users.
	StoreValue *bool `json:"store_value,omitempty"`

	// Indicates if the API key is locked for further write operations. False by default.
	EntityLock *string `json:"Entity-Lock,omitempty"`

	// Allows users to set headers on API requests
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

// SetEntityLock : Allow user to set EntityLock
func (_options *CreateAPIKeyOptions) SetEntityLock(entityLock string) *CreateAPIKeyOptions {
	_options.EntityLock = core.StringPtr(entityLock)
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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
	return
}

// UnmarshalCreateProfileLinkRequestLink unmarshals an instance of CreateProfileLinkRequestLink from the specified map of raw messages.
func UnmarshalCreateProfileLinkRequestLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateProfileLinkRequestLink)
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
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

	// Allows users to set headers on API requests
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

// CreateReportOptions : The CreateReport options.
type CreateReportOptions struct {
	// ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Optional report type. The supported value is 'inactive'. List all identities that have not authenticated within the
	// time indicated by duration.
	Type *string `json:"type,omitempty"`

	// Optional duration of the report. The supported unit of duration is hours.
	Duration *string `json:"duration,omitempty"`

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

// DeleteAPIKeyOptions : The DeleteAPIKey options.
type DeleteAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

// DeleteServiceIDOptions : The DeleteServiceID options.
type DeleteServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
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
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id_account", &obj.IamIDAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "params", &obj.Params)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_authn", &obj.LastAuthn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAccountSettingsOptions : The GetAccountSettings options.
type GetAccountSettingsOptions struct {
	// Unique ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Allows users to set headers on API requests
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

// GetAPIKeyOptions : The GetAPIKey options.
type GetAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"include_history,omitempty"`

	// Defines if the entity's activity is included in the response. Retrieving activity data is an expensive operation, so
	// please only request this when needed.
	IncludeActivity *bool `json:"include_activity,omitempty"`

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

// GetLinkOptions : The GetLink options.
type GetLinkOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"profile-id" validate:"required,ne="`

	// ID of the link.
	LinkID *string `json:"link-id" validate:"required,ne="`

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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
	// please only request this when needed.
	IncludeActivity *bool `json:"include_activity,omitempty"`

	// Allows users to set headers on API requests
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

// GetReportOptions : The GetReport options.
type GetReportOptions struct {
	// ID of the account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Reference for the report to be generated, You can use 'latest' to get the latest report for the given account.
	Reference *string `json:"reference" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// please only request this when needed.
	IncludeActivity *bool `json:"include_activity,omitempty"`

	// Allows users to set headers on API requests
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

// UnmarshalIDBasedMfaEnrollment unmarshals an instance of IDBasedMfaEnrollment from the specified map of raw messages.
func UnmarshalIDBasedMfaEnrollment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IDBasedMfaEnrollment)
	err = core.UnmarshalPrimitive(m, "trait_account_default", &obj.TraitAccountDefault)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trait_user_specific", &obj.TraitUserSpecific)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trait_effective", &obj.TraitEffective)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "complies", &obj.Complies)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

// LockAPIKeyOptions : The LockAPIKey options.
type LockAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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
		return
	}
	err = core.UnmarshalPrimitive(m, "enrolled", &obj.Enrolled)
	if err != nil {
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
		return
	}
	err = core.UnmarshalModel(m, "id_based_mfa", &obj.IDBasedMfa, UnmarshalIDBasedMfaEnrollment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account_based_mfa", &obj.AccountBasedMfa, UnmarshalAccountBasedMfaEnrollment)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
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
	err = core.UnmarshalPrimitive(m, "realm_name", &obj.RealmName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration", &obj.Expiration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cr_type", &obj.CrType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalProfileClaimRuleConditions)
	if err != nil {
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
	return
}

// UnmarshalProfileClaimRuleConditions unmarshals an instance of ProfileClaimRuleConditions from the specified map of raw messages.
func UnmarshalProfileClaimRuleConditions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileClaimRuleConditions)
	err = core.UnmarshalPrimitive(m, "claim", &obj.Claim)
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
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalProfileClaimRule)
	if err != nil {
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
	Identities []ProfileIdentity `json:"identities,omitempty"`
}

// UnmarshalProfileIdentitiesResponse unmarshals an instance of ProfileIdentitiesResponse from the specified map of raw messages.
func UnmarshalProfileIdentitiesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileIdentitiesResponse)
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "identities", &obj.Identities, UnmarshalProfileIdentity)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileIdentity : ProfileIdentity struct
type ProfileIdentity struct {
	// IAM ID of the identity.
	IamID *string `json:"iam_id,omitempty"`

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

// Constants associated with the ProfileIdentity.Type property.
// Type of the identity.
const (
	ProfileIdentityTypeCRNConst       = "crn"
	ProfileIdentityTypeServiceidConst = "serviceid"
	ProfileIdentityTypeUserConst      = "user"
)

// NewProfileIdentity : Instantiate ProfileIdentity (Generic Model Constructor)
func (*IamIdentityV1) NewProfileIdentity(identifier string, typeVar string) (_model *ProfileIdentity, err error) {
	_model = &ProfileIdentity{
		Identifier: core.StringPtr(identifier),
		Type:       core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalProfileIdentity unmarshals an instance of ProfileIdentity from the specified map of raw messages.
func UnmarshalProfileIdentity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileIdentity)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "identifier", &obj.Identifier)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "accounts", &obj.Accounts)
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
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cr_type", &obj.CrType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "link", &obj.Link, UnmarshalProfileLinkLink)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_duration", &obj.ReportDuration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_start_time", &obj.ReportStartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_end_time", &obj.ReportEndTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "users", &obj.Users, UnmarshalUserActivity)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "apikeys", &obj.Apikeys, UnmarshalApikeyActivity)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "serviceids", &obj.Serviceids, UnmarshalEntityActivity)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalEntityActivity)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_time", &obj.ReportTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ims_account_id", &obj.ImsAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "users", &obj.Users, UnmarshalUserReportMfaEnrollmentStatus)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "operation", &obj.Operation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_agent", &obj.UserAgent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "thread_id", &obj.ThreadID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "host", &obj.Host)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "end_time", &obj.EndTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "elapsed_time", &obj.ElapsedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_name", &obj.ClusterName)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
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
	err = core.UnmarshalPrimitive(m, "unique_instance_crns", &obj.UniqueInstanceCrns)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "apikey", &obj.Apikey, UnmarshalAPIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "activity", &obj.Activity, UnmarshalActivity)
	if err != nil {
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
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "serviceids", &obj.Serviceids, UnmarshalServiceID)
	if err != nil {
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
	Identities []ProfileIdentity `json:"identities,omitempty"`

	// Allows users to set headers on API requests
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
func (_options *SetProfileIdentitiesOptions) SetIdentities(identities []ProfileIdentity) *SetProfileIdentitiesOptions {
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

	// IAM ID of the identity.
	IamID *string `json:"iam_id,omitempty"`

	// Only valid for the type user. Accounts from which a user can assume the trusted profile.
	Accounts []string `json:"accounts,omitempty"`

	// Description of the identity that can assume the trusted profile. This is optional field for all the types of
	// identities. When this field is not set for the identity type 'serviceid' then the description of the service id is
	// used. Description is recommended for the identity type 'crn' E.g. 'Instance 1234 of IBM Cloud Service project'.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
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

// SetIamID : Allow user to set IamID
func (_options *SetProfileIdentityOptions) SetIamID(iamID string) *SetProfileIdentityOptions {
	_options.IamID = core.StringPtr(iamID)
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
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
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
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ims_account_id", &obj.ImsAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ims_user_id", &obj.ImsUserID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalEnityHistoryRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "activity", &obj.Activity, UnmarshalActivity)
	if err != nil {
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
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalTrustedProfile)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UnlockAPIKeyOptions : The UnlockAPIKey options.
type UnlockAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

// UpdateAccountSettingsOptions : The UpdateAccountSettings options.
type UpdateAccountSettingsOptions struct {
	// Version of the account settings to be updated. Specify the version that you retrieved as entity_tag (ETag header)
	// when reading the account. This value helps identifying parallel usage of this API. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The id of the account to update the settings for.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Defines whether or not creating a Service Id is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
	//   * NOT_SET - to unset a previously set value.
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateAccountSettingsOptions.RestrictCreateServiceID property.
// Defines whether or not creating a Service Id is access controlled. Valid values:
//   - RESTRICTED - to apply access control
//   - NOT_RESTRICTED - to remove access control
//   - NOT_SET - to unset a previously set value.
const (
	UpdateAccountSettingsOptionsRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	UpdateAccountSettingsOptionsRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	UpdateAccountSettingsOptionsRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the UpdateAccountSettingsOptions.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   - RESTRICTED - to apply access control
//   - NOT_RESTRICTED - to remove access control
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_authn", &obj.LastAuthn)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "effective_mfa_type", &obj.EffectiveMfaType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "id_based_mfa", &obj.IDBasedMfa, UnmarshalIDBasedMfaEnrollment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account_based_mfa", &obj.AccountBasedMfa, UnmarshalAccountBasedMfaEnrollment)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "enrollments", &obj.Enrollments, UnmarshalMfaEnrollments)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
