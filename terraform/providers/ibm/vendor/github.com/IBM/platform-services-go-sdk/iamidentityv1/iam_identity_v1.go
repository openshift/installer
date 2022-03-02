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
 * IBM OpenAPI SDK Code Generator Version: 3.37.0-a85661cd-20210802-190136
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
// Version: 1.0.0
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
// for themself, or service ID API keys for  service IDs that are bound to an entity they have access to. In case of
// service IDs and their API keys, a user must be either an account owner,  a IBM Cloud org manager or IBM Cloud space
// developer in order to manage  service IDs of the entity.
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
// for  service IDs that are bound to an entity they have access to.
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
// Returns the details of an API key. Users can manage user API keys for themself, or service ID API keys for  service
// IDs that are bound to an entity they have access to. In case of  service IDs and their API keys, a user must be
// either an account owner,  a IBM Cloud org manager or IBM Cloud space developer in order to manage  service IDs of the
// entity.
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
// service ID API  keys for service IDs that are bound to an entity they have access  to.
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
// that are bound to an entity they have access to. Note: apikey details are only included in the response when
// creating a Service ID with an apikey.
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
// creating a Service ID with an apikey.
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
// service ID API keys for service IDs that are bound to an entity they have access to.   Note: apikey details are only
// included in the response when creating a  Service ID with an apikey.
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
		"rule-id": *getClaimRuleOptions.RuleID,
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
		"rule-id": *updateClaimRuleOptions.RuleID,
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
		"rule-id": *deleteClaimRuleOptions.RuleID,
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
		"link-id": *getLinkOptions.LinkID,
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
		"link-id": *deleteLinkOptions.LinkID,
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
// Allows a user to configure settings on their account with regards to MFA, session lifetimes,  access control for
// creating new identities, and enforcing IP restrictions on  token creation.
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
	if updateAccountSettingsOptions.SessionExpirationInSeconds != nil {
		body["session_expiration_in_seconds"] = updateAccountSettingsOptions.SessionExpirationInSeconds
	}
	if updateAccountSettingsOptions.SessionInvalidationInSeconds != nil {
		body["session_invalidation_in_seconds"] = updateAccountSettingsOptions.SessionInvalidationInSeconds
	}
	if updateAccountSettingsOptions.MaxSessionsPerIdentity != nil {
		body["max_sessions_per_identity"] = updateAccountSettingsOptions.MaxSessionsPerIdentity
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
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa" validate:"required"`

	// History of the Account Settings.
	History []EnityHistoryRecord `json:"history,omitempty"`

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string `json:"session_expiration_in_seconds" validate:"required"`

	// Defines the period of time in seconds in which a session will be invalidated due  to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string `json:"session_invalidation_in_seconds" validate:"required"`

	// Defines the max allowed sessions per identity required by the account. Valid values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string `json:"max_sessions_per_identity" validate:"required"`
}

// Constants associated with the AccountSettingsResponse.RestrictCreateServiceID property.
// Defines whether or not creating a Service Id is access controlled. Valid values:
//   * RESTRICTED - to apply access control
//   * NOT_RESTRICTED - to remove access control
//   * NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsResponseRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsResponseRestrictCreateServiceIDNotSetConst = "NOT_SET"
	AccountSettingsResponseRestrictCreateServiceIDRestrictedConst = "RESTRICTED"
)

// Constants associated with the AccountSettingsResponse.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   * RESTRICTED - to apply access control
//   * NOT_RESTRICTED - to remove access control
//   * NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsResponseRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsResponseRestrictCreatePlatformApikeyNotSetConst = "NOT_SET"
	AccountSettingsResponseRestrictCreatePlatformApikeyRestrictedConst = "RESTRICTED"
)

// Constants associated with the AccountSettingsResponse.Mfa property.
// Defines the MFA trait for the account. Valid values:
//   * NONE - No MFA trait set
//   * TOTP - For all non-federated IBMId users
//   * TOTP4ALL - For all users
//   * LEVEL1 - Email-based MFA for all users
//   * LEVEL2 - TOTP-based MFA for all users
//   * LEVEL3 - U2F MFA for all users.
const (
	AccountSettingsResponseMfaLevel1Const = "LEVEL1"
	AccountSettingsResponseMfaLevel2Const = "LEVEL2"
	AccountSettingsResponseMfaLevel3Const = "LEVEL3"
	AccountSettingsResponseMfaNoneConst = "NONE"
	AccountSettingsResponseMfaTotpConst = "TOTP"
	AccountSettingsResponseMfaTotp4allConst = "TOTP4ALL"
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
	EntityLock *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateAPIKeyOptions : Instantiate CreateAPIKeyOptions
func (*IamIdentityV1) NewCreateAPIKeyOptions(name string, iamID string) *CreateAPIKeyOptions {
	return &CreateAPIKeyOptions{
		Name: core.StringPtr(name),
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
	ProfileID *string `json:"-" validate:"required,ne="`

	// Type of the calim rule, either 'Profile-SAML' or 'Profile-CR'.
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
		ProfileID: core.StringPtr(profileID),
		Type: core.StringPtr(typeVar),
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
	ProfileID *string `json:"-" validate:"required,ne="`

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
		CrType: core.StringPtr(crType),
		Link: link,
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
		CRN: core.StringPtr(crn),
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
		Name: core.StringPtr(name),
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
	EntityLock *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateServiceIDOptions : Instantiate CreateServiceIDOptions
func (*IamIdentityV1) NewCreateServiceIDOptions(accountID string, name string) *CreateServiceIDOptions {
	return &CreateServiceIDOptions{
		AccountID: core.StringPtr(accountID),
		Name: core.StringPtr(name),
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
	ID *string `json:"-" validate:"required,ne="`

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
	ProfileID *string `json:"-" validate:"required,ne="`

	// ID of the claim rule to delete.
	RuleID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteClaimRuleOptions : Instantiate DeleteClaimRuleOptions
func (*IamIdentityV1) NewDeleteClaimRuleOptions(profileID string, ruleID string) *DeleteClaimRuleOptions {
	return &DeleteClaimRuleOptions{
		ProfileID: core.StringPtr(profileID),
		RuleID: core.StringPtr(ruleID),
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
	ProfileID *string `json:"-" validate:"required,ne="`

	// ID of the link.
	LinkID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteLinkOptions : Instantiate DeleteLinkOptions
func (*IamIdentityV1) NewDeleteLinkOptions(profileID string, linkID string) *DeleteLinkOptions {
	return &DeleteLinkOptions{
		ProfileID: core.StringPtr(profileID),
		LinkID: core.StringPtr(linkID),
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

// DeleteProfileOptions : The DeleteProfile options.
type DeleteProfileOptions struct {
	// ID of the trusted profile.
	ProfileID *string `json:"-" validate:"required,ne="`

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
	ID *string `json:"-" validate:"required,ne="`

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

// GetAccountSettingsOptions : The GetAccountSettings options.
type GetAccountSettingsOptions struct {
	// Unique ID of the account.
	AccountID *string `json:"-" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"-"`

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
	ID *string `json:"-" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"-"`

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

// SetHeaders : Allow user to set Headers
func (options *GetAPIKeyOptions) SetHeaders(param map[string]string) *GetAPIKeyOptions {
	options.Headers = param
	return options
}

// GetAPIKeysDetailsOptions : The GetAPIKeysDetails options.
type GetAPIKeysDetailsOptions struct {
	// API key value.
	IamAPIKey *string `json:"-"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"-"`

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
	ProfileID *string `json:"-" validate:"required,ne="`

	// ID of the claim rule to get.
	RuleID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetClaimRuleOptions : Instantiate GetClaimRuleOptions
func (*IamIdentityV1) NewGetClaimRuleOptions(profileID string, ruleID string) *GetClaimRuleOptions {
	return &GetClaimRuleOptions{
		ProfileID: core.StringPtr(profileID),
		RuleID: core.StringPtr(ruleID),
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
	ProfileID *string `json:"-" validate:"required,ne="`

	// ID of the link.
	LinkID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLinkOptions : Instantiate GetLinkOptions
func (*IamIdentityV1) NewGetLinkOptions(profileID string, linkID string) *GetLinkOptions {
	return &GetLinkOptions{
		ProfileID: core.StringPtr(profileID),
		LinkID: core.StringPtr(linkID),
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

// GetProfileOptions : The GetProfile options.
type GetProfileOptions struct {
	// ID of the trusted profile to get.
	ProfileID *string `json:"-" validate:"required,ne="`

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

// SetHeaders : Allow user to set Headers
func (options *GetProfileOptions) SetHeaders(param map[string]string) *GetProfileOptions {
	options.Headers = param
	return options
}

// GetServiceIDOptions : The GetServiceID options.
type GetServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `json:"-" validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"-"`

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

// SetHeaders : Allow user to set Headers
func (options *GetServiceIDOptions) SetHeaders(param map[string]string) *GetServiceIDOptions {
	options.Headers = param
	return options
}

// ListAPIKeysOptions : The ListAPIKeys options.
type ListAPIKeysOptions struct {
	// Account ID of the API keys(s) to query. If a service IAM ID is specified in iam_id then account_id must match the
	// account of the IAM ID. If a user IAM ID is specified in iam_id then then account_id must match the account of the
	// Authorization token.
	AccountID *string `json:"-"`

	// IAM ID of the API key(s) to be queried. The IAM ID may be that of a user or a service. For a user IAM ID iam_id must
	// match the Authorization token.
	IamID *string `json:"-"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Pagesize *int64 `json:"-"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"-"`

	// Optional parameter to define the scope of the queried API Keys. Can be 'entity' (default) or 'account'.
	Scope *string `json:"-"`

	// Optional parameter to filter the type of the queried API Keys. Can be 'user' or 'serviceid'.
	Type *string `json:"-"`

	// Optional sort property, valid values are name, description, created_at and created_by. If specified, the items are
	// sorted by the value of this property.
	Sort *string `json:"-"`

	// Optional sort order, valid values are asc and desc. Default: asc.
	Order *string `json:"-"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListAPIKeysOptions.Scope property.
// Optional parameter to define the scope of the queried API Keys. Can be 'entity' (default) or 'account'.
const (
	ListAPIKeysOptionsScopeAccountConst = "account"
	ListAPIKeysOptionsScopeEntityConst = "entity"
)

// Constants associated with the ListAPIKeysOptions.Type property.
// Optional parameter to filter the type of the queried API Keys. Can be 'user' or 'serviceid'.
const (
	ListAPIKeysOptionsTypeServiceidConst = "serviceid"
	ListAPIKeysOptionsTypeUserConst = "user"
)

// Constants associated with the ListAPIKeysOptions.Order property.
// Optional sort order, valid values are asc and desc. Default: asc.
const (
	ListAPIKeysOptionsOrderAscConst = "asc"
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
	ProfileID *string `json:"-" validate:"required,ne="`

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
	ProfileID *string `json:"-" validate:"required,ne="`

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
	AccountID *string `json:"-" validate:"required"`

	// Name of the trusted profile to query.
	Name *string `json:"-"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Pagesize *int64 `json:"-"`

	// Optional sort property, valid values are name, description, created_at and modified_at. If specified, the items are
	// sorted by the value of this property.
	Sort *string `json:"-"`

	// Optional sort order, valid values are asc and desc. Default: asc.
	Order *string `json:"-"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"-"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListProfilesOptions.Order property.
// Optional sort order, valid values are asc and desc. Default: asc.
const (
	ListProfilesOptionsOrderAscConst = "asc"
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
	AccountID *string `json:"-"`

	// Name of the service ID(s) to query. Optional.20 items per page. Valid range is 1 to 100.
	Name *string `json:"-"`

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Pagesize *int64 `json:"-"`

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string `json:"-"`

	// Optional sort property, valid values are name, description, created_at and modified_at. If specified, the items are
	// sorted by the value of this property.
	Sort *string `json:"-"`

	// Optional sort order, valid values are asc and desc. Default: asc.
	Order *string `json:"-"`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListServiceIdsOptions.Order property.
// Optional sort order, valid values are asc and desc. Default: asc.
const (
	ListServiceIdsOptionsOrderAscConst = "asc"
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
	ID *string `json:"-" validate:"required,ne="`

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
	ID *string `json:"-" validate:"required,ne="`

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

	// Type of the Calim rule, either 'Profile-SAML' or 'Profile-CR'.
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
	// The claim to evaluate against.
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
		Claim: core.StringPtr(claim),
		Operator: core.StringPtr(operator),
		Value: core.StringPtr(value),
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

// ProfileLink : Link details.
type ProfileLink struct {
	// the unique identifier of the claim rule.
	ID *string `json:"id" validate:"required"`

	// version of the claim rule.
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
	ID *string `json:"-" validate:"required,ne="`

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
	ID *string `json:"-" validate:"required,ne="`

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
	// Version of the account settings to be updated. Specify the version that you  retrieved as entity_tag (ETag header)
	// when reading the account. This value helps  identifying parallel usage of this API. Pass * to indicate to update any
	// version  available. This might result in stale updates.
	IfMatch *string `json:"-" validate:"required"`

	// The id of the account to update the settings for.
	AccountID *string `json:"-" validate:"required,ne="`

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
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string `json:"mfa,omitempty"`

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string `json:"session_expiration_in_seconds,omitempty"`

	// Defines the period of time in seconds in which a session will be invalidated due  to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string `json:"session_invalidation_in_seconds,omitempty"`

	// Defines the max allowed sessions per identity required by the account. Value values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string `json:"max_sessions_per_identity,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateAccountSettingsOptions.RestrictCreateServiceID property.
// Defines whether or not creating a Service Id is access controlled. Valid values:
//   * RESTRICTED - to apply access control
//   * NOT_RESTRICTED - to remove access control
//   * NOT_SET - to unset a previously set value.
const (
	UpdateAccountSettingsOptionsRestrictCreateServiceIDNotRestrictedConst = "NOT_RESTRICTED"
	UpdateAccountSettingsOptionsRestrictCreateServiceIDNotSetConst = "NOT_SET"
	UpdateAccountSettingsOptionsRestrictCreateServiceIDRestrictedConst = "RESTRICTED"
)

// Constants associated with the UpdateAccountSettingsOptions.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   * RESTRICTED - to apply access control
//   * NOT_RESTRICTED - to remove access control
//   * NOT_SET - to 'unset' a previous set value.
const (
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyNotSetConst = "NOT_SET"
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyRestrictedConst = "RESTRICTED"
)

// Constants associated with the UpdateAccountSettingsOptions.Mfa property.
// Defines the MFA trait for the account. Valid values:
//   * NONE - No MFA trait set
//   * TOTP - For all non-federated IBMId users
//   * TOTP4ALL - For all users
//   * LEVEL1 - Email-based MFA for all users
//   * LEVEL2 - TOTP-based MFA for all users
//   * LEVEL3 - U2F MFA for all users.
const (
	UpdateAccountSettingsOptionsMfaLevel1Const = "LEVEL1"
	UpdateAccountSettingsOptionsMfaLevel2Const = "LEVEL2"
	UpdateAccountSettingsOptionsMfaLevel3Const = "LEVEL3"
	UpdateAccountSettingsOptionsMfaNoneConst = "NONE"
	UpdateAccountSettingsOptionsMfaTotpConst = "TOTP"
	UpdateAccountSettingsOptionsMfaTotp4allConst = "TOTP4ALL"
)

// NewUpdateAccountSettingsOptions : Instantiate UpdateAccountSettingsOptions
func (*IamIdentityV1) NewUpdateAccountSettingsOptions(ifMatch string, accountID string) *UpdateAccountSettingsOptions {
	return &UpdateAccountSettingsOptions{
		IfMatch: core.StringPtr(ifMatch),
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

// SetHeaders : Allow user to set Headers
func (options *UpdateAccountSettingsOptions) SetHeaders(param map[string]string) *UpdateAccountSettingsOptions {
	options.Headers = param
	return options
}

// UpdateAPIKeyOptions : The UpdateAPIKey options.
type UpdateAPIKeyOptions struct {
	// Unique ID of the API key to be updated.
	ID *string `json:"-" validate:"required,ne="`

	// Version of the API key to be updated. Specify the version that you retrieved when reading the API key. This value
	// helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might result
	// in stale updates.
	IfMatch *string `json:"-" validate:"required"`

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
		ID: core.StringPtr(id),
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
	ProfileID *string `json:"-" validate:"required,ne="`

	// ID of the claim rule to update.
	RuleID *string `json:"-" validate:"required,ne="`

	// Version of the claim rule to be updated.  Specify the version that you retrived when reading list of claim rules.
	// This value helps to identify any parallel usage of claim rule. Pass * to indicate to update any version available.
	// This might result in stale updates.
	IfMatch *string `json:"-" validate:"required"`

	// Type of the calim rule, either 'Profile-SAML' or 'Profile-CR'.
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
		ProfileID: core.StringPtr(profileID),
		RuleID: core.StringPtr(ruleID),
		IfMatch: core.StringPtr(ifMatch),
		Type: core.StringPtr(typeVar),
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
	ProfileID *string `json:"-" validate:"required,ne="`

	// Version of the trusted profile to be updated.  Specify the version that you retrived when reading list of trusted
	// profiles. This value helps to identify any parallel usage of trusted profile. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"-" validate:"required"`

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
		IfMatch: core.StringPtr(ifMatch),
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
	ID *string `json:"-" validate:"required,ne="`

	// Version of the service ID to be updated. Specify the version that you retrieved as entity_tag (ETag header) when
	// reading the service ID. This value helps identifying parallel usage of this API. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"-" validate:"required"`

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
		ID: core.StringPtr(id),
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
