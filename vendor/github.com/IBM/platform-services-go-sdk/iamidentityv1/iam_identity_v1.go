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
 * IBM OpenAPI SDK Code Generator Version: 99-SNAPSHOT-46891d34-20210426-162952
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKeyList)
	if err != nil {
		return
	}
	response.Result = result

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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
	if err != nil {
		return
	}
	response.Result = result

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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
	if err != nil {
		return
	}
	response.Result = result

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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
	if err != nil {
		return
	}
	response.Result = result

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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIKey)
	if err != nil {
		return
	}
	response.Result = result

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
// that are bound to an entity they have access to.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceIDList)
	if err != nil {
		return
	}
	response.Result = result

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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetServiceID : Get details of a service ID
// Returns the details of a service ID. Users can manage user API keys for themself, or service ID API keys for service
// IDs that are bound to an entity they have access to.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateServiceID : Update service ID
// Updates properties of a service ID. This does NOT affect existing access tokens. Their token content will stay
// unchanged until the access token is refreshed. To update a service ID, pass the property to be modified. To delete
// one property's value, pass the property with an empty value "".Users can manage user API keys for themself, or
// service ID API keys for service IDs that are bound to an entity they have access to.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceID)
	if err != nil {
		return
	}
	response.Result = result

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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsResponse)
	if err != nil {
		return
	}
	response.Result = result

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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsResponse)
	if err != nil {
		return
	}
	response.Result = result

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
	AccountSettingsResponseRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	AccountSettingsResponseRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the AccountSettingsResponse.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   * RESTRICTED - to apply access control
//   * NOT_RESTRICTED - to remove access control
//   * NOT_SET - to 'unset' a previous set value.
const (
	AccountSettingsResponseRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	AccountSettingsResponseRestrictCreatePlatformApikeyNotSetConst        = "NOT_SET"
	AccountSettingsResponseRestrictCreatePlatformApikeyRestrictedConst    = "RESTRICTED"
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
	AccountSettingsResponseMfaLevel1Const   = "LEVEL1"
	AccountSettingsResponseMfaLevel2Const   = "LEVEL2"
	AccountSettingsResponseMfaLevel3Const   = "LEVEL3"
	AccountSettingsResponseMfaNoneConst     = "NONE"
	AccountSettingsResponseMfaTotpConst     = "TOTP"
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
func (*IamIdentityV1) NewAPIKeyInsideCreateServiceIDRequest(name string) (model *APIKeyInsideCreateServiceIDRequest, err error) {
	model = &APIKeyInsideCreateServiceIDRequest{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(model, "required parameters")
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
	Name *string `validate:"required"`

	// The iam_id that this API key authenticates.
	IamID *string `validate:"required"`

	// The optional description of the API key. The 'description' property is only available if a description was provided
	// during a create of an API key.
	Description *string

	// The account ID of the API key.
	AccountID *string

	// You can optionally passthrough the API key value for this API key. If passed, NO validation of that apiKey value is
	// done, i.e. the value can be non-URL safe. If omitted, the API key management will create an URL safe opaque API key
	// value. The value of the API key is checked for uniqueness. Please ensure enough variations when passing in this
	// value.
	Apikey *string

	// Send true or false to set whether the API key value is retrievable in the future by using the Get details of an API
	// key request. If you create an API key for a user, you must specify `false` or omit the value. We don't allow storing
	// of API keys for users.
	StoreValue *bool

	// Indicates if the API key is locked for further write operations. False by default.
	EntityLock *string

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
func (options *CreateAPIKeyOptions) SetName(name string) *CreateAPIKeyOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetIamID : Allow user to set IamID
func (options *CreateAPIKeyOptions) SetIamID(iamID string) *CreateAPIKeyOptions {
	options.IamID = core.StringPtr(iamID)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateAPIKeyOptions) SetDescription(description string) *CreateAPIKeyOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *CreateAPIKeyOptions) SetAccountID(accountID string) *CreateAPIKeyOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetApikey : Allow user to set Apikey
func (options *CreateAPIKeyOptions) SetApikey(apikey string) *CreateAPIKeyOptions {
	options.Apikey = core.StringPtr(apikey)
	return options
}

// SetStoreValue : Allow user to set StoreValue
func (options *CreateAPIKeyOptions) SetStoreValue(storeValue bool) *CreateAPIKeyOptions {
	options.StoreValue = core.BoolPtr(storeValue)
	return options
}

// SetEntityLock : Allow user to set EntityLock
func (options *CreateAPIKeyOptions) SetEntityLock(entityLock string) *CreateAPIKeyOptions {
	options.EntityLock = core.StringPtr(entityLock)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAPIKeyOptions) SetHeaders(param map[string]string) *CreateAPIKeyOptions {
	options.Headers = param
	return options
}

// CreateServiceIDOptions : The CreateServiceID options.
type CreateServiceIDOptions struct {
	// ID of the account the service ID belongs to.
	AccountID *string `validate:"required"`

	// Name of the Service Id. The name is not checked for uniqueness. Therefore multiple names with the same value can
	// exist. Access is done via the UUID of the Service Id.
	Name *string `validate:"required"`

	// The optional description of the Service Id. The 'description' property is only available if a description was
	// provided during a create of a Service Id.
	Description *string

	// Optional list of CRNs (string array) which point to the services connected to the service ID.
	UniqueInstanceCrns []string

	// Parameters for the API key in the Create service Id V1 REST request.
	Apikey *APIKeyInsideCreateServiceIDRequest

	// Indicates if the service ID is locked for further write operations. False by default.
	EntityLock *string

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
func (options *CreateServiceIDOptions) SetAccountID(accountID string) *CreateServiceIDOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetName : Allow user to set Name
func (options *CreateServiceIDOptions) SetName(name string) *CreateServiceIDOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateServiceIDOptions) SetDescription(description string) *CreateServiceIDOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetUniqueInstanceCrns : Allow user to set UniqueInstanceCrns
func (options *CreateServiceIDOptions) SetUniqueInstanceCrns(uniqueInstanceCrns []string) *CreateServiceIDOptions {
	options.UniqueInstanceCrns = uniqueInstanceCrns
	return options
}

// SetApikey : Allow user to set Apikey
func (options *CreateServiceIDOptions) SetApikey(apikey *APIKeyInsideCreateServiceIDRequest) *CreateServiceIDOptions {
	options.Apikey = apikey
	return options
}

// SetEntityLock : Allow user to set EntityLock
func (options *CreateServiceIDOptions) SetEntityLock(entityLock string) *CreateServiceIDOptions {
	options.EntityLock = core.StringPtr(entityLock)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateServiceIDOptions) SetHeaders(param map[string]string) *CreateServiceIDOptions {
	options.Headers = param
	return options
}

// DeleteAPIKeyOptions : The DeleteAPIKey options.
type DeleteAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `validate:"required,ne="`

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
func (options *DeleteAPIKeyOptions) SetID(id string) *DeleteAPIKeyOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAPIKeyOptions) SetHeaders(param map[string]string) *DeleteAPIKeyOptions {
	options.Headers = param
	return options
}

// DeleteServiceIDOptions : The DeleteServiceID options.
type DeleteServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `validate:"required,ne="`

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
func (options *DeleteServiceIDOptions) SetID(id string) *DeleteServiceIDOptions {
	options.ID = core.StringPtr(id)
	return options
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
	AccountID *string `validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool

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
func (options *GetAccountSettingsOptions) SetAccountID(accountID string) *GetAccountSettingsOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (options *GetAccountSettingsOptions) SetIncludeHistory(includeHistory bool) *GetAccountSettingsOptions {
	options.IncludeHistory = core.BoolPtr(includeHistory)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountSettingsOptions) SetHeaders(param map[string]string) *GetAccountSettingsOptions {
	options.Headers = param
	return options
}

// GetAPIKeyOptions : The GetAPIKey options.
type GetAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool

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
func (options *GetAPIKeyOptions) SetID(id string) *GetAPIKeyOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (options *GetAPIKeyOptions) SetIncludeHistory(includeHistory bool) *GetAPIKeyOptions {
	options.IncludeHistory = core.BoolPtr(includeHistory)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetAPIKeyOptions) SetHeaders(param map[string]string) *GetAPIKeyOptions {
	options.Headers = param
	return options
}

// GetAPIKeysDetailsOptions : The GetAPIKeysDetails options.
type GetAPIKeysDetailsOptions struct {
	// API key value.
	IamAPIKey *string

	// Defines if the entity history is included in the response.
	IncludeHistory *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAPIKeysDetailsOptions : Instantiate GetAPIKeysDetailsOptions
func (*IamIdentityV1) NewGetAPIKeysDetailsOptions() *GetAPIKeysDetailsOptions {
	return &GetAPIKeysDetailsOptions{}
}

// SetIamAPIKey : Allow user to set IamAPIKey
func (options *GetAPIKeysDetailsOptions) SetIamAPIKey(iamAPIKey string) *GetAPIKeysDetailsOptions {
	options.IamAPIKey = core.StringPtr(iamAPIKey)
	return options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (options *GetAPIKeysDetailsOptions) SetIncludeHistory(includeHistory bool) *GetAPIKeysDetailsOptions {
	options.IncludeHistory = core.BoolPtr(includeHistory)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetAPIKeysDetailsOptions) SetHeaders(param map[string]string) *GetAPIKeysDetailsOptions {
	options.Headers = param
	return options
}

// GetServiceIDOptions : The GetServiceID options.
type GetServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `validate:"required,ne="`

	// Defines if the entity history is included in the response.
	IncludeHistory *bool

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
func (options *GetServiceIDOptions) SetID(id string) *GetServiceIDOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (options *GetServiceIDOptions) SetIncludeHistory(includeHistory bool) *GetServiceIDOptions {
	options.IncludeHistory = core.BoolPtr(includeHistory)
	return options
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
	AccountID *string

	// IAM ID of the API key(s) to be queried. The IAM ID may be that of a user or a service. For a user IAM ID iam_id must
	// match the Authorization token.
	IamID *string

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Pagesize *int64

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string

	// Optional parameter to define the scope of the queried API Keys. Can be 'entity' (default) or 'account'.
	Scope *string

	// Optional parameter to filter the type of the queried API Keys. Can be 'user' or 'serviceid'.
	Type *string

	// Optional sort property, valid values are name, description, created_at and created_by. If specified, the items are
	// sorted by the value of this property.
	Sort *string

	// Optional sort order, valid values are asc and desc. Default: asc.
	Order *string

	// Defines if the entity history is included in the response.
	IncludeHistory *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListAPIKeysOptions.Scope property.
// Optional parameter to define the scope of the queried API Keys. Can be 'entity' (default) or 'account'.
const (
	ListAPIKeysOptionsScopeAccountConst = "account"
	ListAPIKeysOptionsScopeEntityConst  = "entity"
)

// Constants associated with the ListAPIKeysOptions.Type property.
// Optional parameter to filter the type of the queried API Keys. Can be 'user' or 'serviceid'.
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
func (options *ListAPIKeysOptions) SetAccountID(accountID string) *ListAPIKeysOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetIamID : Allow user to set IamID
func (options *ListAPIKeysOptions) SetIamID(iamID string) *ListAPIKeysOptions {
	options.IamID = core.StringPtr(iamID)
	return options
}

// SetPagesize : Allow user to set Pagesize
func (options *ListAPIKeysOptions) SetPagesize(pagesize int64) *ListAPIKeysOptions {
	options.Pagesize = core.Int64Ptr(pagesize)
	return options
}

// SetPagetoken : Allow user to set Pagetoken
func (options *ListAPIKeysOptions) SetPagetoken(pagetoken string) *ListAPIKeysOptions {
	options.Pagetoken = core.StringPtr(pagetoken)
	return options
}

// SetScope : Allow user to set Scope
func (options *ListAPIKeysOptions) SetScope(scope string) *ListAPIKeysOptions {
	options.Scope = core.StringPtr(scope)
	return options
}

// SetType : Allow user to set Type
func (options *ListAPIKeysOptions) SetType(typeVar string) *ListAPIKeysOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListAPIKeysOptions) SetSort(sort string) *ListAPIKeysOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetOrder : Allow user to set Order
func (options *ListAPIKeysOptions) SetOrder(order string) *ListAPIKeysOptions {
	options.Order = core.StringPtr(order)
	return options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (options *ListAPIKeysOptions) SetIncludeHistory(includeHistory bool) *ListAPIKeysOptions {
	options.IncludeHistory = core.BoolPtr(includeHistory)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListAPIKeysOptions) SetHeaders(param map[string]string) *ListAPIKeysOptions {
	options.Headers = param
	return options
}

// ListServiceIdsOptions : The ListServiceIds options.
type ListServiceIdsOptions struct {
	// Account ID of the service ID(s) to query. This parameter is required (unless using a pagetoken).
	AccountID *string

	// Name of the service ID(s) to query. Optional.20 items per page. Valid range is 1 to 100.
	Name *string

	// Optional size of a single page. Default is 20 items per page. Valid range is 1 to 100.
	Pagesize *int64

	// Optional Prev or Next page token returned from a previous query execution. Default is start with first page.
	Pagetoken *string

	// Optional sort property, valid values are name, description, created_at and modified_at. If specified, the items are
	// sorted by the value of this property.
	Sort *string

	// Optional sort order, valid values are asc and desc. Default: asc.
	Order *string

	// Defines if the entity history is included in the response.
	IncludeHistory *bool

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
func (options *ListServiceIdsOptions) SetAccountID(accountID string) *ListServiceIdsOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetName : Allow user to set Name
func (options *ListServiceIdsOptions) SetName(name string) *ListServiceIdsOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetPagesize : Allow user to set Pagesize
func (options *ListServiceIdsOptions) SetPagesize(pagesize int64) *ListServiceIdsOptions {
	options.Pagesize = core.Int64Ptr(pagesize)
	return options
}

// SetPagetoken : Allow user to set Pagetoken
func (options *ListServiceIdsOptions) SetPagetoken(pagetoken string) *ListServiceIdsOptions {
	options.Pagetoken = core.StringPtr(pagetoken)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListServiceIdsOptions) SetSort(sort string) *ListServiceIdsOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetOrder : Allow user to set Order
func (options *ListServiceIdsOptions) SetOrder(order string) *ListServiceIdsOptions {
	options.Order = core.StringPtr(order)
	return options
}

// SetIncludeHistory : Allow user to set IncludeHistory
func (options *ListServiceIdsOptions) SetIncludeHistory(includeHistory bool) *ListServiceIdsOptions {
	options.IncludeHistory = core.BoolPtr(includeHistory)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListServiceIdsOptions) SetHeaders(param map[string]string) *ListServiceIdsOptions {
	options.Headers = param
	return options
}

// LockAPIKeyOptions : The LockAPIKey options.
type LockAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `validate:"required,ne="`

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
func (options *LockAPIKeyOptions) SetID(id string) *LockAPIKeyOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *LockAPIKeyOptions) SetHeaders(param map[string]string) *LockAPIKeyOptions {
	options.Headers = param
	return options
}

// LockServiceIDOptions : The LockServiceID options.
type LockServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `validate:"required,ne="`

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
func (options *LockServiceIDOptions) SetID(id string) *LockServiceIDOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *LockServiceIDOptions) SetHeaders(param map[string]string) *LockServiceIDOptions {
	options.Headers = param
	return options
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
	EntityTag *string `json:"entity_tag,omitempty"`

	// Cloud Resource Name of the item. Example Cloud Resource Name:
	// 'crn:v1:bluemix:public:iam-identity:us-south:a/myaccount::serviceid:1234-5678-9012'.
	CRN *string `json:"crn" validate:"required"`

	// The service ID cannot be changed if set to true.
	Locked *bool `json:"locked" validate:"required"`

	// If set contains a date time string of the creation date in ISO format.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// If set contains a date time string of the last modification date in ISO format.
	ModifiedAt *strfmt.DateTime `json:"modified_at,omitempty"`

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
	Apikey *APIKey `json:"apikey" validate:"required"`
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

// UnlockAPIKeyOptions : The UnlockAPIKey options.
type UnlockAPIKeyOptions struct {
	// Unique ID of the API key.
	ID *string `validate:"required,ne="`

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
func (options *UnlockAPIKeyOptions) SetID(id string) *UnlockAPIKeyOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UnlockAPIKeyOptions) SetHeaders(param map[string]string) *UnlockAPIKeyOptions {
	options.Headers = param
	return options
}

// UnlockServiceIDOptions : The UnlockServiceID options.
type UnlockServiceIDOptions struct {
	// Unique ID of the service ID.
	ID *string `validate:"required,ne="`

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
func (options *UnlockServiceIDOptions) SetID(id string) *UnlockServiceIDOptions {
	options.ID = core.StringPtr(id)
	return options
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
	IfMatch *string `validate:"required"`

	// The id of the account to update the settings for.
	AccountID *string `validate:"required,ne="`

	// Defines whether or not creating a Service Id is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
	//   * NOT_SET - to unset a previously set value.
	RestrictCreateServiceID *string

	// Defines whether or not creating platform API keys is access controlled. Valid values:
	//   * RESTRICTED - to apply access control
	//   * NOT_RESTRICTED - to remove access control
	//   * NOT_SET - to 'unset' a previous set value.
	RestrictCreatePlatformApikey *string

	// Defines the IP addresses and subnets from which IAM tokens can be created for the account.
	AllowedIPAddresses *string

	// Defines the MFA trait for the account. Valid values:
	//   * NONE - No MFA trait set
	//   * TOTP - For all non-federated IBMId users
	//   * TOTP4ALL - For all users
	//   * LEVEL1 - Email-based MFA for all users
	//   * LEVEL2 - TOTP-based MFA for all users
	//   * LEVEL3 - U2F MFA for all users.
	Mfa *string

	// Defines the session expiration in seconds for the account. Valid values:
	//   * Any whole number between between '900' and '86400'
	//   * NOT_SET - To unset account setting and use service default.
	SessionExpirationInSeconds *string

	// Defines the period of time in seconds in which a session will be invalidated due  to inactivity. Valid values:
	//   * Any whole number between '900' and '7200'
	//   * NOT_SET - To unset account setting and use service default.
	SessionInvalidationInSeconds *string

	// Defines the max allowed sessions per identity required by the account. Value values:
	//   * Any whole number greater than 0
	//   * NOT_SET - To unset account setting and use service default.
	MaxSessionsPerIdentity *string

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
	UpdateAccountSettingsOptionsRestrictCreateServiceIDNotSetConst        = "NOT_SET"
	UpdateAccountSettingsOptionsRestrictCreateServiceIDRestrictedConst    = "RESTRICTED"
)

// Constants associated with the UpdateAccountSettingsOptions.RestrictCreatePlatformApikey property.
// Defines whether or not creating platform API keys is access controlled. Valid values:
//   * RESTRICTED - to apply access control
//   * NOT_RESTRICTED - to remove access control
//   * NOT_SET - to 'unset' a previous set value.
const (
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyNotRestrictedConst = "NOT_RESTRICTED"
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyNotSetConst        = "NOT_SET"
	UpdateAccountSettingsOptionsRestrictCreatePlatformApikeyRestrictedConst    = "RESTRICTED"
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
	UpdateAccountSettingsOptionsMfaLevel1Const   = "LEVEL1"
	UpdateAccountSettingsOptionsMfaLevel2Const   = "LEVEL2"
	UpdateAccountSettingsOptionsMfaLevel3Const   = "LEVEL3"
	UpdateAccountSettingsOptionsMfaNoneConst     = "NONE"
	UpdateAccountSettingsOptionsMfaTotpConst     = "TOTP"
	UpdateAccountSettingsOptionsMfaTotp4allConst = "TOTP4ALL"
)

// NewUpdateAccountSettingsOptions : Instantiate UpdateAccountSettingsOptions
func (*IamIdentityV1) NewUpdateAccountSettingsOptions(ifMatch string, accountID string) *UpdateAccountSettingsOptions {
	return &UpdateAccountSettingsOptions{
		IfMatch:   core.StringPtr(ifMatch),
		AccountID: core.StringPtr(accountID),
	}
}

// SetIfMatch : Allow user to set IfMatch
func (options *UpdateAccountSettingsOptions) SetIfMatch(ifMatch string) *UpdateAccountSettingsOptions {
	options.IfMatch = core.StringPtr(ifMatch)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *UpdateAccountSettingsOptions) SetAccountID(accountID string) *UpdateAccountSettingsOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetRestrictCreateServiceID : Allow user to set RestrictCreateServiceID
func (options *UpdateAccountSettingsOptions) SetRestrictCreateServiceID(restrictCreateServiceID string) *UpdateAccountSettingsOptions {
	options.RestrictCreateServiceID = core.StringPtr(restrictCreateServiceID)
	return options
}

// SetRestrictCreatePlatformApikey : Allow user to set RestrictCreatePlatformApikey
func (options *UpdateAccountSettingsOptions) SetRestrictCreatePlatformApikey(restrictCreatePlatformApikey string) *UpdateAccountSettingsOptions {
	options.RestrictCreatePlatformApikey = core.StringPtr(restrictCreatePlatformApikey)
	return options
}

// SetAllowedIPAddresses : Allow user to set AllowedIPAddresses
func (options *UpdateAccountSettingsOptions) SetAllowedIPAddresses(allowedIPAddresses string) *UpdateAccountSettingsOptions {
	options.AllowedIPAddresses = core.StringPtr(allowedIPAddresses)
	return options
}

// SetMfa : Allow user to set Mfa
func (options *UpdateAccountSettingsOptions) SetMfa(mfa string) *UpdateAccountSettingsOptions {
	options.Mfa = core.StringPtr(mfa)
	return options
}

// SetSessionExpirationInSeconds : Allow user to set SessionExpirationInSeconds
func (options *UpdateAccountSettingsOptions) SetSessionExpirationInSeconds(sessionExpirationInSeconds string) *UpdateAccountSettingsOptions {
	options.SessionExpirationInSeconds = core.StringPtr(sessionExpirationInSeconds)
	return options
}

// SetSessionInvalidationInSeconds : Allow user to set SessionInvalidationInSeconds
func (options *UpdateAccountSettingsOptions) SetSessionInvalidationInSeconds(sessionInvalidationInSeconds string) *UpdateAccountSettingsOptions {
	options.SessionInvalidationInSeconds = core.StringPtr(sessionInvalidationInSeconds)
	return options
}

// SetMaxSessionsPerIdentity : Allow user to set MaxSessionsPerIdentity
func (options *UpdateAccountSettingsOptions) SetMaxSessionsPerIdentity(maxSessionsPerIdentity string) *UpdateAccountSettingsOptions {
	options.MaxSessionsPerIdentity = core.StringPtr(maxSessionsPerIdentity)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccountSettingsOptions) SetHeaders(param map[string]string) *UpdateAccountSettingsOptions {
	options.Headers = param
	return options
}

// UpdateAPIKeyOptions : The UpdateAPIKey options.
type UpdateAPIKeyOptions struct {
	// Unique ID of the API key to be updated.
	ID *string `validate:"required,ne="`

	// Version of the API key to be updated. Specify the version that you retrieved when reading the API key. This value
	// helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might result
	// in stale updates.
	IfMatch *string `validate:"required"`

	// The name of the API key to update. If specified in the request the parameter must not be empty. The name is not
	// checked for uniqueness. Failure to this will result in an Error condition.
	Name *string

	// The description of the API key to update. If specified an empty description will clear the description of the API
	// key. If a non empty value is provided the API key will be updated.
	Description *string

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
func (options *UpdateAPIKeyOptions) SetID(id string) *UpdateAPIKeyOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetIfMatch : Allow user to set IfMatch
func (options *UpdateAPIKeyOptions) SetIfMatch(ifMatch string) *UpdateAPIKeyOptions {
	options.IfMatch = core.StringPtr(ifMatch)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateAPIKeyOptions) SetName(name string) *UpdateAPIKeyOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateAPIKeyOptions) SetDescription(description string) *UpdateAPIKeyOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAPIKeyOptions) SetHeaders(param map[string]string) *UpdateAPIKeyOptions {
	options.Headers = param
	return options
}

// UpdateServiceIDOptions : The UpdateServiceID options.
type UpdateServiceIDOptions struct {
	// Unique ID of the service ID to be updated.
	ID *string `validate:"required,ne="`

	// Version of the service ID to be updated. Specify the version that you retrieved as entity_tag (ETag header) when
	// reading the service ID. This value helps identifying parallel usage of this API. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `validate:"required"`

	// The name of the service ID to update. If specified in the request the parameter must not be empty. The name is not
	// checked for uniqueness. Failure to this will result in an Error condition.
	Name *string

	// The description of the service ID to update. If specified an empty description will clear the description of the
	// service ID. If an non empty value is provided the service ID will be updated.
	Description *string

	// List of CRNs which point to the services connected to this service ID. If specified an empty list will clear all
	// existing unique instance crns of the service ID.
	UniqueInstanceCrns []string

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
func (options *UpdateServiceIDOptions) SetID(id string) *UpdateServiceIDOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetIfMatch : Allow user to set IfMatch
func (options *UpdateServiceIDOptions) SetIfMatch(ifMatch string) *UpdateServiceIDOptions {
	options.IfMatch = core.StringPtr(ifMatch)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateServiceIDOptions) SetName(name string) *UpdateServiceIDOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateServiceIDOptions) SetDescription(description string) *UpdateServiceIDOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetUniqueInstanceCrns : Allow user to set UniqueInstanceCrns
func (options *UpdateServiceIDOptions) SetUniqueInstanceCrns(uniqueInstanceCrns []string) *UpdateServiceIDOptions {
	options.UniqueInstanceCrns = uniqueInstanceCrns
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateServiceIDOptions) SetHeaders(param map[string]string) *UpdateServiceIDOptions {
	options.Headers = param
	return options
}
