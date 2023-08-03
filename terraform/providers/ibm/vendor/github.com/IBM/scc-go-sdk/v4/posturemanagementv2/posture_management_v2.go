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
 * IBM OpenAPI SDK Code Generator Version: 3.38.1-1037b405-20210908-184149
 */

// Package posturemanagementv2 : Operations and models for the PostureManagementV2 service
package posturemanagementv2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/scc-go-sdk/v4/common"
	"github.com/go-openapi/strfmt"
)

// PostureManagementV2 : With IBM Cloud® Security and Compliance Center, you can embed checks into your every day
// workflows to help manage your current security and compliance posture. By monitoring for risks, you can identify
// security vulnerabilities and quickly work to mitigate the impact.
//
// API Version: 2.0.0
type PostureManagementV2 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us.compliance.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "posture_management"

const ParameterizedServiceURL = "https://{environment}.cloud.ibm.com"

var defaultUrlVariables = map[string]string{
	"environment": "us.compliance",
}

// PostureManagementV2Options : Service options
type PostureManagementV2Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewPostureManagementV2UsingExternalConfig : constructs an instance of PostureManagementV2 with passed in options and external configuration.
func NewPostureManagementV2UsingExternalConfig(options *PostureManagementV2Options) (postureManagement *PostureManagementV2, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	postureManagement, err = NewPostureManagementV2(options)
	if err != nil {
		return
	}

	err = postureManagement.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = postureManagement.Service.SetServiceURL(options.URL)
	}
	return
}

// NewPostureManagementV2 : constructs an instance of PostureManagementV2 with passed in options.
func NewPostureManagementV2(options *PostureManagementV2Options) (service *PostureManagementV2, err error) {
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

	service = &PostureManagementV2{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south": "https://us.compliance.cloud.ibm.com",
		"us-east":  "https://us.compliance.cloud.ibm.com",
		"eu-de":    "https://eu.compliance.cloud.ibm.com",
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "postureManagement" suitable for processing requests.
func (postureManagement *PostureManagementV2) Clone() *PostureManagementV2 {
	if core.IsNil(postureManagement) {
		return nil
	}
	clone := *postureManagement
	clone.Service = postureManagement.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (postureManagement *PostureManagementV2) SetServiceURL(url string) error {
	return postureManagement.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (postureManagement *PostureManagementV2) GetServiceURL() string {
	return postureManagement.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (postureManagement *PostureManagementV2) SetDefaultHeaders(headers http.Header) {
	postureManagement.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (postureManagement *PostureManagementV2) SetEnableGzipCompression(enableGzip bool) {
	postureManagement.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (postureManagement *PostureManagementV2) GetEnableGzipCompression() bool {
	return postureManagement.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (postureManagement *PostureManagementV2) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	postureManagement.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (postureManagement *PostureManagementV2) DisableRetries() {
	postureManagement.Service.DisableRetries()
}

// CreateCredential : Add a credential
// Add an existing credential that can be used by a collector to access your resources to gather information about your
// configurations, validate them, and initiate any remediation where possible.
func (postureManagement *PostureManagementV2) CreateCredential(createCredentialOptions *CreateCredentialOptions) (result *Credential, response *core.DetailedResponse, err error) {
	return postureManagement.CreateCredentialWithContext(context.Background(), createCredentialOptions)
}

// CreateCredentialWithContext is an alternate form of the CreateCredential method which supports a Context parameter
func (postureManagement *PostureManagementV2) CreateCredentialWithContext(ctx context.Context, createCredentialOptions *CreateCredentialOptions) (result *Credential, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCredentialOptions, "createCredentialOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createCredentialOptions, "createCredentialOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/credentials`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createCredentialOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "CreateCredential")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createCredentialOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createCredentialOptions.TransactionID))
	}

	if createCredentialOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createCredentialOptions.AccountID))
	}

	body := make(map[string]interface{})
	if createCredentialOptions.Enabled != nil {
		body["enabled"] = createCredentialOptions.Enabled
	}
	if createCredentialOptions.Type != nil {
		body["type"] = createCredentialOptions.Type
	}
	if createCredentialOptions.Name != nil {
		body["name"] = createCredentialOptions.Name
	}
	if createCredentialOptions.Description != nil {
		body["description"] = createCredentialOptions.Description
	}
	if createCredentialOptions.DisplayFields != nil {
		body["display_fields"] = createCredentialOptions.DisplayFields
	}
	if createCredentialOptions.Purpose != nil {
		body["purpose"] = createCredentialOptions.Purpose
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCredential)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListCredentials : List credentials
// List the credentials that were previously added to the Security and Compliance Center.
func (postureManagement *PostureManagementV2) ListCredentials(listCredentialsOptions *ListCredentialsOptions) (result *CredentialList, response *core.DetailedResponse, err error) {
	return postureManagement.ListCredentialsWithContext(context.Background(), listCredentialsOptions)
}

// ListCredentialsWithContext is an alternate form of the ListCredentials method which supports a Context parameter
func (postureManagement *PostureManagementV2) ListCredentialsWithContext(ctx context.Context, listCredentialsOptions *ListCredentialsOptions) (result *CredentialList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCredentialsOptions, "listCredentialsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/credentials`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listCredentialsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ListCredentials")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listCredentialsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listCredentialsOptions.TransactionID))
	}

	if listCredentialsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listCredentialsOptions.AccountID))
	}
	if listCredentialsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listCredentialsOptions.Offset))
	}
	if listCredentialsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listCredentialsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCredentialList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCredential : View credential details
// View the details of a stored credential, which include its name, type, and secret.
func (postureManagement *PostureManagementV2) GetCredential(getCredentialOptions *GetCredentialOptions) (result *Credential, response *core.DetailedResponse, err error) {
	return postureManagement.GetCredentialWithContext(context.Background(), getCredentialOptions)
}

// GetCredentialWithContext is an alternate form of the GetCredential method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetCredentialWithContext(ctx context.Context, getCredentialOptions *GetCredentialOptions) (result *Credential, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCredentialOptions, "getCredentialOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCredentialOptions, "getCredentialOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getCredentialOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/credentials/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCredentialOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetCredential")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getCredentialOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getCredentialOptions.TransactionID))
	}

	if getCredentialOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getCredentialOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCredential)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateCredential : Update a credential
// Update the way that a credential is stored in the Security and Compliance Center, or update the credential itself.
func (postureManagement *PostureManagementV2) UpdateCredential(updateCredentialOptions *UpdateCredentialOptions) (result *Credential, response *core.DetailedResponse, err error) {
	return postureManagement.UpdateCredentialWithContext(context.Background(), updateCredentialOptions)
}

// UpdateCredentialWithContext is an alternate form of the UpdateCredential method which supports a Context parameter
func (postureManagement *PostureManagementV2) UpdateCredentialWithContext(ctx context.Context, updateCredentialOptions *UpdateCredentialOptions) (result *Credential, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCredentialOptions, "updateCredentialOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateCredentialOptions, "updateCredentialOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateCredentialOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/credentials/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateCredentialOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "UpdateCredential")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateCredentialOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateCredentialOptions.TransactionID))
	}

	if updateCredentialOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*updateCredentialOptions.AccountID))
	}

	body := make(map[string]interface{})
	if updateCredentialOptions.Enabled != nil {
		body["enabled"] = updateCredentialOptions.Enabled
	}
	if updateCredentialOptions.Type != nil {
		body["type"] = updateCredentialOptions.Type
	}
	if updateCredentialOptions.Name != nil {
		body["name"] = updateCredentialOptions.Name
	}
	if updateCredentialOptions.Description != nil {
		body["description"] = updateCredentialOptions.Description
	}
	if updateCredentialOptions.DisplayFields != nil {
		body["display_fields"] = updateCredentialOptions.DisplayFields
	}
	if updateCredentialOptions.Purpose != nil {
		body["purpose"] = updateCredentialOptions.Purpose
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCredential)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCredential : Delete a credential
// If you no longer need to use a credential, you can remove it from the Security and Compliance Center.
func (postureManagement *PostureManagementV2) DeleteCredential(deleteCredentialOptions *DeleteCredentialOptions) (response *core.DetailedResponse, err error) {
	return postureManagement.DeleteCredentialWithContext(context.Background(), deleteCredentialOptions)
}

// DeleteCredentialWithContext is an alternate form of the DeleteCredential method which supports a Context parameter
func (postureManagement *PostureManagementV2) DeleteCredentialWithContext(ctx context.Context, deleteCredentialOptions *DeleteCredentialOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCredentialOptions, "deleteCredentialOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCredentialOptions, "deleteCredentialOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteCredentialOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/credentials/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCredentialOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "DeleteCredential")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteCredentialOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteCredentialOptions.TransactionID))
	}

	if deleteCredentialOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteCredentialOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = postureManagement.Service.Request(request, nil)

	return
}

// CreateCollector : Create a collector
// Create a collector to fetch the configuration information of your resources and then validate it by using a specified
// profile.
func (postureManagement *PostureManagementV2) CreateCollector(createCollectorOptions *CreateCollectorOptions) (result *Collector, response *core.DetailedResponse, err error) {
	return postureManagement.CreateCollectorWithContext(context.Background(), createCollectorOptions)
}

// CreateCollectorWithContext is an alternate form of the CreateCollector method which supports a Context parameter
func (postureManagement *PostureManagementV2) CreateCollectorWithContext(ctx context.Context, createCollectorOptions *CreateCollectorOptions) (result *Collector, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCollectorOptions, "createCollectorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createCollectorOptions, "createCollectorOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/collectors`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createCollectorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "CreateCollector")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createCollectorOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createCollectorOptions.TransactionID))
	}

	if createCollectorOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createCollectorOptions.AccountID))
	}

	body := make(map[string]interface{})
	if createCollectorOptions.Name != nil {
		body["name"] = createCollectorOptions.Name
	}
	if createCollectorOptions.IsPublic != nil {
		body["is_public"] = createCollectorOptions.IsPublic
	}
	if createCollectorOptions.ManagedBy != nil {
		body["managed_by"] = createCollectorOptions.ManagedBy
	}
	if createCollectorOptions.Description != nil {
		body["description"] = createCollectorOptions.Description
	}
	if createCollectorOptions.Passphrase != nil {
		body["passphrase"] = createCollectorOptions.Passphrase
	}
	if createCollectorOptions.IsUbiImage != nil {
		body["is_ubi_image"] = createCollectorOptions.IsUbiImage
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollector)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListCollectors : List collectors
// View a list of all of the collectors that are available in your account and their status.
func (postureManagement *PostureManagementV2) ListCollectors(listCollectorsOptions *ListCollectorsOptions) (result *CollectorList, response *core.DetailedResponse, err error) {
	return postureManagement.ListCollectorsWithContext(context.Background(), listCollectorsOptions)
}

// ListCollectorsWithContext is an alternate form of the ListCollectors method which supports a Context parameter
func (postureManagement *PostureManagementV2) ListCollectorsWithContext(ctx context.Context, listCollectorsOptions *ListCollectorsOptions) (result *CollectorList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCollectorsOptions, "listCollectorsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/collectors`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listCollectorsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ListCollectors")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listCollectorsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listCollectorsOptions.TransactionID))
	}

	if listCollectorsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listCollectorsOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollectorList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCollector : View collector details
// View collector details .
func (postureManagement *PostureManagementV2) GetCollector(getCollectorOptions *GetCollectorOptions) (result *Collector, response *core.DetailedResponse, err error) {
	return postureManagement.GetCollectorWithContext(context.Background(), getCollectorOptions)
}

// GetCollectorWithContext is an alternate form of the GetCollector method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetCollectorWithContext(ctx context.Context, getCollectorOptions *GetCollectorOptions) (result *Collector, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCollectorOptions, "getCollectorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCollectorOptions, "getCollectorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getCollectorOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/collectors/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCollectorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetCollector")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getCollectorOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getCollectorOptions.TransactionID))
	}

	if getCollectorOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getCollectorOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollector)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateCollector : Update a collector
// Update a collector.
func (postureManagement *PostureManagementV2) UpdateCollector(updateCollectorOptions *UpdateCollectorOptions) (result *Collector, response *core.DetailedResponse, err error) {
	return postureManagement.UpdateCollectorWithContext(context.Background(), updateCollectorOptions)
}

// UpdateCollectorWithContext is an alternate form of the UpdateCollector method which supports a Context parameter
func (postureManagement *PostureManagementV2) UpdateCollectorWithContext(ctx context.Context, updateCollectorOptions *UpdateCollectorOptions) (result *Collector, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCollectorOptions, "updateCollectorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateCollectorOptions, "updateCollectorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateCollectorOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/collectors/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateCollectorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "UpdateCollector")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateCollectorOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateCollectorOptions.TransactionID))
	}

	if updateCollectorOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*updateCollectorOptions.AccountID))
	}

	_, err = builder.SetBodyContentJSON(updateCollectorOptions.Collector)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollector)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCollector : Delete a collector
// Delete a collector from the Security and Compliance Center that you no longer need.
func (postureManagement *PostureManagementV2) DeleteCollector(deleteCollectorOptions *DeleteCollectorOptions) (response *core.DetailedResponse, err error) {
	return postureManagement.DeleteCollectorWithContext(context.Background(), deleteCollectorOptions)
}

// DeleteCollectorWithContext is an alternate form of the DeleteCollector method which supports a Context parameter
func (postureManagement *PostureManagementV2) DeleteCollectorWithContext(ctx context.Context, deleteCollectorOptions *DeleteCollectorOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCollectorOptions, "deleteCollectorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCollectorOptions, "deleteCollectorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteCollectorOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/collectors/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCollectorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "DeleteCollector")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteCollectorOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteCollectorOptions.TransactionID))
	}

	if deleteCollectorOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteCollectorOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = postureManagement.Service.Request(request, nil)

	return
}

// ImportProfiles : Import profile
// Import a profile that you formatted locally. For more information about how your profile must be formatted, see [the
// docs](/docs/security-compliance?topic=security-compliance-custom-profiles#CSV-format).
func (postureManagement *PostureManagementV2) ImportProfiles(importProfilesOptions *ImportProfilesOptions) (result *BasicResult, response *core.DetailedResponse, err error) {
	return postureManagement.ImportProfilesWithContext(context.Background(), importProfilesOptions)
}

// ImportProfilesWithContext is an alternate form of the ImportProfiles method which supports a Context parameter
func (postureManagement *PostureManagementV2) ImportProfilesWithContext(ctx context.Context, importProfilesOptions *ImportProfilesOptions) (result *BasicResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(importProfilesOptions, "importProfilesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(importProfilesOptions, "importProfilesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/profiles/import`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range importProfilesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ImportProfiles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if importProfilesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*importProfilesOptions.TransactionID))
	}

	if importProfilesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*importProfilesOptions.AccountID))
	}

	builder.AddFormData("file", "filename",
		"text/csv", importProfilesOptions.File)

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBasicResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProfiles : List profiles
// List all of the profiles that are available in your account. To view a specific profile, you can filter by name.
func (postureManagement *PostureManagementV2) ListProfiles(listProfilesOptions *ListProfilesOptions) (result *ProfileList, response *core.DetailedResponse, err error) {
	return postureManagement.ListProfilesWithContext(context.Background(), listProfilesOptions)
}

// ListProfilesWithContext is an alternate form of the ListProfiles method which supports a Context parameter
func (postureManagement *PostureManagementV2) ListProfilesWithContext(ctx context.Context, listProfilesOptions *ListProfilesOptions) (result *ProfileList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProfilesOptions, "listProfilesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/profiles`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProfilesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ListProfiles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listProfilesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listProfilesOptions.TransactionID))
	}

	if listProfilesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listProfilesOptions.AccountID))
	}
	if listProfilesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listProfilesOptions.Offset))
	}
	if listProfilesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProfilesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProfile : View profile details
// View profile details.
func (postureManagement *PostureManagementV2) GetProfile(getProfileOptions *GetProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return postureManagement.GetProfileWithContext(context.Background(), getProfileOptions)
}

// GetProfileWithContext is an alternate form of the GetProfile method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetProfileWithContext(ctx context.Context, getProfileOptions *GetProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileOptions, "getProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileOptions, "getProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getProfileOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/profiles/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getProfileOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getProfileOptions.TransactionID))
	}

	builder.AddQuery("profile_type", fmt.Sprint(*getProfileOptions.ProfileType))
	if getProfileOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getProfileOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
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

// UpdateProfiles : Update a profile
// Update a profile. Set the enable field to false to mark the profile as deleted.
func (postureManagement *PostureManagementV2) UpdateProfiles(updateProfilesOptions *UpdateProfilesOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return postureManagement.UpdateProfilesWithContext(context.Background(), updateProfilesOptions)
}

// UpdateProfilesWithContext is an alternate form of the UpdateProfiles method which supports a Context parameter
func (postureManagement *PostureManagementV2) UpdateProfilesWithContext(ctx context.Context, updateProfilesOptions *UpdateProfilesOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProfilesOptions, "updateProfilesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateProfilesOptions, "updateProfilesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateProfilesOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/profiles/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateProfilesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "UpdateProfiles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateProfilesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateProfilesOptions.TransactionID))
	}

	if updateProfilesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*updateProfilesOptions.AccountID))
	}

	body := make(map[string]interface{})
	if updateProfilesOptions.Name != nil {
		body["name"] = updateProfilesOptions.Name
	}
	if updateProfilesOptions.Description != nil {
		body["description"] = updateProfilesOptions.Description
	}
	if updateProfilesOptions.BaseProfile != nil {
		body["base_profile"] = updateProfilesOptions.BaseProfile
	}
	if updateProfilesOptions.Type != nil {
		body["type"] = updateProfilesOptions.Type
	}
	if updateProfilesOptions.IsEnabled != nil {
		body["is_enabled"] = updateProfilesOptions.IsEnabled
	}
	if updateProfilesOptions.ControlIds != nil {
		body["control_ids"] = updateProfilesOptions.ControlIds
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
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

// DeleteProfile : Delete a profile
// Delete a custom profile that was previously created in your account.
func (postureManagement *PostureManagementV2) DeleteProfile(deleteProfileOptions *DeleteProfileOptions) (response *core.DetailedResponse, err error) {
	return postureManagement.DeleteProfileWithContext(context.Background(), deleteProfileOptions)
}

// DeleteProfileWithContext is an alternate form of the DeleteProfile method which supports a Context parameter
func (postureManagement *PostureManagementV2) DeleteProfileWithContext(ctx context.Context, deleteProfileOptions *DeleteProfileOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProfileOptions, "deleteProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProfileOptions, "deleteProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteProfileOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/profiles/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "DeleteProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteProfileOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteProfileOptions.TransactionID))
	}

	if deleteProfileOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteProfileOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = postureManagement.Service.Request(request, nil)

	return
}

// GetProfileControls : View profile controls
// View a list of the controls and their associated goals for a specified profile.
func (postureManagement *PostureManagementV2) GetProfileControls(getProfileControlsOptions *GetProfileControlsOptions) (result *ControlList, response *core.DetailedResponse, err error) {
	return postureManagement.GetProfileControlsWithContext(context.Background(), getProfileControlsOptions)
}

// GetProfileControlsWithContext is an alternate form of the GetProfileControls method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetProfileControlsWithContext(ctx context.Context, getProfileControlsOptions *GetProfileControlsOptions) (result *ControlList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileControlsOptions, "getProfileControlsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileControlsOptions, "getProfileControlsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"profile_id": *getProfileControlsOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/profiles/{profile_id}/controls`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProfileControlsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetProfileControls")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getProfileControlsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getProfileControlsOptions.TransactionID))
	}

	if getProfileControlsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getProfileControlsOptions.AccountID))
	}
	if getProfileControlsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*getProfileControlsOptions.Offset))
	}
	if getProfileControlsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getProfileControlsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetGroupProfileControls : View group profile controls
// View the controls and goals that are available as part of a profile group.
func (postureManagement *PostureManagementV2) GetGroupProfileControls(getGroupProfileControlsOptions *GetGroupProfileControlsOptions) (result *ControlList, response *core.DetailedResponse, err error) {
	return postureManagement.GetGroupProfileControlsWithContext(context.Background(), getGroupProfileControlsOptions)
}

// GetGroupProfileControlsWithContext is an alternate form of the GetGroupProfileControls method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetGroupProfileControlsWithContext(ctx context.Context, getGroupProfileControlsOptions *GetGroupProfileControlsOptions) (result *ControlList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getGroupProfileControlsOptions, "getGroupProfileControlsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getGroupProfileControlsOptions, "getGroupProfileControlsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"group_id": *getGroupProfileControlsOptions.GroupID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/profiles/groups/{group_id}/controls`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getGroupProfileControlsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetGroupProfileControls")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getGroupProfileControlsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getGroupProfileControlsOptions.TransactionID))
	}

	if getGroupProfileControlsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getGroupProfileControlsOptions.AccountID))
	}
	if getGroupProfileControlsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*getGroupProfileControlsOptions.Offset))
	}
	if getGroupProfileControlsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getGroupProfileControlsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateScope : Create a scope
// A scope is the selection of resources that you want to validate the configuration of.
func (postureManagement *PostureManagementV2) CreateScope(createScopeOptions *CreateScopeOptions) (result *Scope, response *core.DetailedResponse, err error) {
	return postureManagement.CreateScopeWithContext(context.Background(), createScopeOptions)
}

// CreateScopeWithContext is an alternate form of the CreateScope method which supports a Context parameter
func (postureManagement *PostureManagementV2) CreateScopeWithContext(ctx context.Context, createScopeOptions *CreateScopeOptions) (result *Scope, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createScopeOptions, "createScopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createScopeOptions, "createScopeOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createScopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "CreateScope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createScopeOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createScopeOptions.TransactionID))
	}

	if createScopeOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createScopeOptions.AccountID))
	}

	body := make(map[string]interface{})
	if createScopeOptions.Name != nil {
		body["name"] = createScopeOptions.Name
	}
	if createScopeOptions.Description != nil {
		body["description"] = createScopeOptions.Description
	}
	if createScopeOptions.CollectorIds != nil {
		body["collector_ids"] = createScopeOptions.CollectorIds
	}
	if createScopeOptions.CredentialID != nil {
		body["credential_id"] = createScopeOptions.CredentialID
	}
	if createScopeOptions.CredentialType != nil {
		body["credential_type"] = createScopeOptions.CredentialType
	}
	if createScopeOptions.Interval != nil {
		body["interval"] = createScopeOptions.Interval
	}
	if createScopeOptions.IsDiscoveryScheduled != nil {
		body["is_discovery_scheduled"] = createScopeOptions.IsDiscoveryScheduled
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScope)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListScopes : List scopes
// List all of the scopes that are available in your account. To view a specific scope, you can filter by name.
func (postureManagement *PostureManagementV2) ListScopes(listScopesOptions *ListScopesOptions) (result *ScopeList, response *core.DetailedResponse, err error) {
	return postureManagement.ListScopesWithContext(context.Background(), listScopesOptions)
}

// ListScopesWithContext is an alternate form of the ListScopes method which supports a Context parameter
func (postureManagement *PostureManagementV2) ListScopesWithContext(ctx context.Context, listScopesOptions *ListScopesOptions) (result *ScopeList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listScopesOptions, "listScopesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listScopesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ListScopes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listScopesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listScopesOptions.TransactionID))
	}

	if listScopesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listScopesOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScopeList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetScopeDetails : View scope details
// View the details of a specific scope.
func (postureManagement *PostureManagementV2) GetScopeDetails(getScopeDetailsOptions *GetScopeDetailsOptions) (result *Scope, response *core.DetailedResponse, err error) {
	return postureManagement.GetScopeDetailsWithContext(context.Background(), getScopeDetailsOptions)
}

// GetScopeDetailsWithContext is an alternate form of the GetScopeDetails method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetScopeDetailsWithContext(ctx context.Context, getScopeDetailsOptions *GetScopeDetailsOptions) (result *Scope, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getScopeDetailsOptions, "getScopeDetailsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getScopeDetailsOptions, "getScopeDetailsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getScopeDetailsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getScopeDetailsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetScopeDetails")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getScopeDetailsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getScopeDetailsOptions.TransactionID))
	}

	if getScopeDetailsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getScopeDetailsOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScope)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateScopeDetails : Update scope
// Update a scope's details.
func (postureManagement *PostureManagementV2) UpdateScopeDetails(updateScopeDetailsOptions *UpdateScopeDetailsOptions) (result *Scope, response *core.DetailedResponse, err error) {
	return postureManagement.UpdateScopeDetailsWithContext(context.Background(), updateScopeDetailsOptions)
}

// UpdateScopeDetailsWithContext is an alternate form of the UpdateScopeDetails method which supports a Context parameter
func (postureManagement *PostureManagementV2) UpdateScopeDetailsWithContext(ctx context.Context, updateScopeDetailsOptions *UpdateScopeDetailsOptions) (result *Scope, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateScopeDetailsOptions, "updateScopeDetailsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateScopeDetailsOptions, "updateScopeDetailsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateScopeDetailsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateScopeDetailsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "UpdateScopeDetails")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateScopeDetailsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateScopeDetailsOptions.TransactionID))
	}

	if updateScopeDetailsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*updateScopeDetailsOptions.AccountID))
	}

	body := make(map[string]interface{})
	if updateScopeDetailsOptions.Name != nil {
		body["name"] = updateScopeDetailsOptions.Name
	}
	if updateScopeDetailsOptions.Description != nil {
		body["description"] = updateScopeDetailsOptions.Description
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScope)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteScope : Delete a scope
// If you no longer need to target a specific scope with your scan, you can delete it.
func (postureManagement *PostureManagementV2) DeleteScope(deleteScopeOptions *DeleteScopeOptions) (response *core.DetailedResponse, err error) {
	return postureManagement.DeleteScopeWithContext(context.Background(), deleteScopeOptions)
}

// DeleteScopeWithContext is an alternate form of the DeleteScope method which supports a Context parameter
func (postureManagement *PostureManagementV2) DeleteScopeWithContext(ctx context.Context, deleteScopeOptions *DeleteScopeOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteScopeOptions, "deleteScopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteScopeOptions, "deleteScopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteScopeOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteScopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "DeleteScope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteScopeOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteScopeOptions.TransactionID))
	}

	if deleteScopeOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteScopeOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = postureManagement.Service.Request(request, nil)

	return
}

// GetScopeTimeline : Get scope event history
// Get the list of events for a scope such as the last discovery or fact collection.
func (postureManagement *PostureManagementV2) GetScopeTimeline(getScopeTimelineOptions *GetScopeTimelineOptions) (result *EventList, response *core.DetailedResponse, err error) {
	return postureManagement.GetScopeTimelineWithContext(context.Background(), getScopeTimelineOptions)
}

// GetScopeTimelineWithContext is an alternate form of the GetScopeTimeline method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetScopeTimelineWithContext(ctx context.Context, getScopeTimelineOptions *GetScopeTimelineOptions) (result *EventList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getScopeTimelineOptions, "getScopeTimelineOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getScopeTimelineOptions, "getScopeTimelineOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"scope_id": *getScopeTimelineOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes/{scope_id}/events`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getScopeTimelineOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetScopeTimeline")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getScopeTimelineOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getScopeTimelineOptions.TransactionID))
	}

	if getScopeTimelineOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getScopeTimelineOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEventList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetScopeDetailsCredentials : Get a scope's credentials
// Get the credentials that are associated with a scope.
func (postureManagement *PostureManagementV2) GetScopeDetailsCredentials(getScopeDetailsCredentialsOptions *GetScopeDetailsCredentialsOptions) (result *ScopeCredential, response *core.DetailedResponse, err error) {
	return postureManagement.GetScopeDetailsCredentialsWithContext(context.Background(), getScopeDetailsCredentialsOptions)
}

// GetScopeDetailsCredentialsWithContext is an alternate form of the GetScopeDetailsCredentials method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetScopeDetailsCredentialsWithContext(ctx context.Context, getScopeDetailsCredentialsOptions *GetScopeDetailsCredentialsOptions) (result *ScopeCredential, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getScopeDetailsCredentialsOptions, "getScopeDetailsCredentialsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getScopeDetailsCredentialsOptions, "getScopeDetailsCredentialsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"scope_id": *getScopeDetailsCredentialsOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes/{scope_id}/credentials`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getScopeDetailsCredentialsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetScopeDetailsCredentials")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getScopeDetailsCredentialsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getScopeDetailsCredentialsOptions.TransactionID))
	}

	if getScopeDetailsCredentialsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getScopeDetailsCredentialsOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScopeCredential)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceScopeDetailsCredentials : Update a scope's credentials
// Update the credentials that are associated with a scope.
func (postureManagement *PostureManagementV2) ReplaceScopeDetailsCredentials(replaceScopeDetailsCredentialsOptions *ReplaceScopeDetailsCredentialsOptions) (result *ScopeCredential, response *core.DetailedResponse, err error) {
	return postureManagement.ReplaceScopeDetailsCredentialsWithContext(context.Background(), replaceScopeDetailsCredentialsOptions)
}

// ReplaceScopeDetailsCredentialsWithContext is an alternate form of the ReplaceScopeDetailsCredentials method which supports a Context parameter
func (postureManagement *PostureManagementV2) ReplaceScopeDetailsCredentialsWithContext(ctx context.Context, replaceScopeDetailsCredentialsOptions *ReplaceScopeDetailsCredentialsOptions) (result *ScopeCredential, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceScopeDetailsCredentialsOptions, "replaceScopeDetailsCredentialsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceScopeDetailsCredentialsOptions, "replaceScopeDetailsCredentialsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"scope_id": *replaceScopeDetailsCredentialsOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes/{scope_id}/credentials`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceScopeDetailsCredentialsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ReplaceScopeDetailsCredentials")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceScopeDetailsCredentialsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*replaceScopeDetailsCredentialsOptions.TransactionID))
	}

	if replaceScopeDetailsCredentialsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*replaceScopeDetailsCredentialsOptions.AccountID))
	}

	body := make(map[string]interface{})
	if replaceScopeDetailsCredentialsOptions.CredentialID != nil {
		body["credential_id"] = replaceScopeDetailsCredentialsOptions.CredentialID
	}
	if replaceScopeDetailsCredentialsOptions.CredentialAttribute != nil {
		body["credential_attribute"] = replaceScopeDetailsCredentialsOptions.CredentialAttribute
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScopeCredential)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetScopeDetailsCollector : Get a scope's collector
// Get the collector that is associated with a scope.
func (postureManagement *PostureManagementV2) GetScopeDetailsCollector(getScopeDetailsCollectorOptions *GetScopeDetailsCollectorOptions) (result *ScopeCollector, response *core.DetailedResponse, err error) {
	return postureManagement.GetScopeDetailsCollectorWithContext(context.Background(), getScopeDetailsCollectorOptions)
}

// GetScopeDetailsCollectorWithContext is an alternate form of the GetScopeDetailsCollector method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetScopeDetailsCollectorWithContext(ctx context.Context, getScopeDetailsCollectorOptions *GetScopeDetailsCollectorOptions) (result *ScopeCollector, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getScopeDetailsCollectorOptions, "getScopeDetailsCollectorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getScopeDetailsCollectorOptions, "getScopeDetailsCollectorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"scope_id": *getScopeDetailsCollectorOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes/{scope_id}/collectors`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getScopeDetailsCollectorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetScopeDetailsCollector")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getScopeDetailsCollectorOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getScopeDetailsCollectorOptions.TransactionID))
	}

	if getScopeDetailsCollectorOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getScopeDetailsCollectorOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScopeCollector)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceScopeDetailsCollector : Update a scope's collector
// Update the collector that is associated with a scope.
func (postureManagement *PostureManagementV2) ReplaceScopeDetailsCollector(replaceScopeDetailsCollectorOptions *ReplaceScopeDetailsCollectorOptions) (result *ScopeCollector, response *core.DetailedResponse, err error) {
	return postureManagement.ReplaceScopeDetailsCollectorWithContext(context.Background(), replaceScopeDetailsCollectorOptions)
}

// ReplaceScopeDetailsCollectorWithContext is an alternate form of the ReplaceScopeDetailsCollector method which supports a Context parameter
func (postureManagement *PostureManagementV2) ReplaceScopeDetailsCollectorWithContext(ctx context.Context, replaceScopeDetailsCollectorOptions *ReplaceScopeDetailsCollectorOptions) (result *ScopeCollector, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceScopeDetailsCollectorOptions, "replaceScopeDetailsCollectorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceScopeDetailsCollectorOptions, "replaceScopeDetailsCollectorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"scope_id": *replaceScopeDetailsCollectorOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scopes/{scope_id}/collectors`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceScopeDetailsCollectorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ReplaceScopeDetailsCollector")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceScopeDetailsCollectorOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*replaceScopeDetailsCollectorOptions.TransactionID))
	}

	if replaceScopeDetailsCollectorOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*replaceScopeDetailsCollectorOptions.AccountID))
	}

	body := make(map[string]interface{})
	if replaceScopeDetailsCollectorOptions.CollectorIds != nil {
		body["collector_ids"] = replaceScopeDetailsCollectorOptions.CollectorIds
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScopeCollector)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCorrelationID : Get status of a scope
// Get the status of any task. By providing a correlation ID, you can track the status of any task that is running or
// completed. So, for example, you can check the status of an initial discovery after scope creation or the status of
// validation after a scan is triggered.
func (postureManagement *PostureManagementV2) GetCorrelationID(getCorrelationIDOptions *GetCorrelationIDOptions) (result *ScopeTaskStatus, response *core.DetailedResponse, err error) {
	return postureManagement.GetCorrelationIDWithContext(context.Background(), getCorrelationIDOptions)
}

// GetCorrelationIDWithContext is an alternate form of the GetCorrelationID method which supports a Context parameter
func (postureManagement *PostureManagementV2) GetCorrelationIDWithContext(ctx context.Context, getCorrelationIDOptions *GetCorrelationIDOptions) (result *ScopeTaskStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCorrelationIDOptions, "getCorrelationIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCorrelationIDOptions, "getCorrelationIDOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"correlation_id": *getCorrelationIDOptions.CorrelationID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scope/status/{correlation_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCorrelationIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "GetCorrelationID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getCorrelationIDOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getCorrelationIDOptions.TransactionID))
	}

	if getCorrelationIDOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getCorrelationIDOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScopeTaskStatus)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListLatestScans : List latest scans
// List the last scan results that are available in your account for each profile and scope combination.
func (postureManagement *PostureManagementV2) ListLatestScans(listLatestScansOptions *ListLatestScansOptions) (result *ScanList, response *core.DetailedResponse, err error) {
	return postureManagement.ListLatestScansWithContext(context.Background(), listLatestScansOptions)
}

// ListLatestScansWithContext is an alternate form of the ListLatestScans method which supports a Context parameter
func (postureManagement *PostureManagementV2) ListLatestScansWithContext(ctx context.Context, listLatestScansOptions *ListLatestScansOptions) (result *ScanList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listLatestScansOptions, "listLatestScansOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scans/validations/latest_scans`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listLatestScansOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ListLatestScans")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listLatestScansOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listLatestScansOptions.TransactionID))
	}

	if listLatestScansOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listLatestScansOptions.AccountID))
	}
	if listLatestScansOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listLatestScansOptions.Offset))
	}
	if listLatestScansOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listLatestScansOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScanList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateValidation : Initiate a validation scan
// Validation scans determine a specified scope's adherence to regulatory controls by validating the configuration of
// the resources in your scope to the attached profile. To initiate a scan, you must configure a collector, provided
// credentials, and completed both a fact collection and discovery scan. [Learn
// more](/docs/security-compliance?topic=security-compliance-schedule-scan).
func (postureManagement *PostureManagementV2) CreateValidation(createValidationOptions *CreateValidationOptions) (result *Result, response *core.DetailedResponse, err error) {
	return postureManagement.CreateValidationWithContext(context.Background(), createValidationOptions)
}

// CreateValidationWithContext is an alternate form of the CreateValidation method which supports a Context parameter
func (postureManagement *PostureManagementV2) CreateValidationWithContext(ctx context.Context, createValidationOptions *CreateValidationOptions) (result *Result, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createValidationOptions, "createValidationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createValidationOptions, "createValidationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scans/validations`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createValidationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "CreateValidation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createValidationOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createValidationOptions.TransactionID))
	}

	if createValidationOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createValidationOptions.AccountID))
	}

	body := make(map[string]interface{})
	if createValidationOptions.ScopeID != nil {
		body["scope_id"] = createValidationOptions.ScopeID
	}
	if createValidationOptions.ProfileID != nil {
		body["profile_id"] = createValidationOptions.ProfileID
	}
	if createValidationOptions.GroupProfileID != nil {
		body["group_profile_id"] = createValidationOptions.GroupProfileID
	}
	if createValidationOptions.Name != nil {
		body["name"] = createValidationOptions.Name
	}
	if createValidationOptions.Description != nil {
		body["description"] = createValidationOptions.Description
	}
	if createValidationOptions.Frequency != nil {
		body["frequency"] = createValidationOptions.Frequency
	}
	if createValidationOptions.NoOfOccurrences != nil {
		body["no_of_occurrences"] = createValidationOptions.NoOfOccurrences
	}
	if createValidationOptions.EndTime != nil {
		body["end_time"] = createValidationOptions.EndTime
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
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ScansSummary : View a specified scan
// Retrieve the results summary of a validation scan by specifying a scan and profile ID combination. To obtain your
// profile ID and scan ID for your most recent scan, make a GET request to the
// "/posture/v2/scans/validations/latest_scans" endpoint.
func (postureManagement *PostureManagementV2) ScansSummary(scansSummaryOptions *ScansSummaryOptions) (result *Summary, response *core.DetailedResponse, err error) {
	return postureManagement.ScansSummaryWithContext(context.Background(), scansSummaryOptions)
}

// ScansSummaryWithContext is an alternate form of the ScansSummary method which supports a Context parameter
func (postureManagement *PostureManagementV2) ScansSummaryWithContext(ctx context.Context, scansSummaryOptions *ScansSummaryOptions) (result *Summary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(scansSummaryOptions, "scansSummaryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(scansSummaryOptions, "scansSummaryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"scan_id": *scansSummaryOptions.ScanID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scans/validations/{scan_id}/summary`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range scansSummaryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ScansSummary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if scansSummaryOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*scansSummaryOptions.TransactionID))
	}

	builder.AddQuery("profile_id", fmt.Sprint(*scansSummaryOptions.ProfileID))
	if scansSummaryOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*scansSummaryOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSummary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ScanSummaries : View scan summaries
// List all of the previous and current validation summaries for a specific scan.
func (postureManagement *PostureManagementV2) ScanSummaries(scanSummariesOptions *ScanSummariesOptions) (result *SummaryList, response *core.DetailedResponse, err error) {
	return postureManagement.ScanSummariesWithContext(context.Background(), scanSummariesOptions)
}

// ScanSummariesWithContext is an alternate form of the ScanSummaries method which supports a Context parameter
func (postureManagement *PostureManagementV2) ScanSummariesWithContext(ctx context.Context, scanSummariesOptions *ScanSummariesOptions) (result *SummaryList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(scanSummariesOptions, "scanSummariesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(scanSummariesOptions, "scanSummariesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = postureManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(postureManagement.Service.Options.URL, `/posture/v2/scans/validations/summaries`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range scanSummariesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("posture_management", "V2", "ScanSummaries")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if scanSummariesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*scanSummariesOptions.TransactionID))
	}

	builder.AddQuery("report_setting_id", fmt.Sprint(*scanSummariesOptions.ReportSettingID))
	if scanSummariesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*scanSummariesOptions.AccountID))
	}
	if scanSummariesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*scanSummariesOptions.Offset))
	}
	if scanSummariesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*scanSummariesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = postureManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSummaryList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ApplicabilityCriteria : The criteria that defines how a profile applies.
type ApplicabilityCriteria struct {
	// A list of environments that a profile can be applied to.
	Environment []string `json:"environment,omitempty"`

	// A list of resources that a profile can be used with.
	Resource []string `json:"resource,omitempty"`

	// The type of environment that a profile is able to be applied to.
	EnvironmentCategory []string `json:"environment_category,omitempty"`

	// The type of resource that a profile is able to be applied to.
	ResourceCategory []string `json:"resource_category,omitempty"`

	// The resource type that the profile applies to.
	ResourceType []string `json:"resource_type,omitempty"`

	// The software that the profile applies to.
	SoftwareDetails interface{} `json:"software_details,omitempty"`

	// The operating system that the profile applies to.
	OsDetails interface{} `json:"os_details,omitempty"`

	// Any additional details about the profile.
	AdditionalDetails interface{} `json:"additional_details,omitempty"`

	// The type of environment that your scope is targeted to.
	EnvironmentCategoryDescription map[string]string `json:"environment_category_description,omitempty"`

	// The environment that your scope is targeted to.
	EnvironmentDescription map[string]string `json:"environment_description,omitempty"`

	// The type of resource that your scope is targeted to.
	ResourceCategoryDescription map[string]string `json:"resource_category_description,omitempty"`

	// A further classification of the type of resource that your scope is targeted to.
	ResourceTypeDescription map[string]string `json:"resource_type_description,omitempty"`

	// The resource that is scanned as part of your scope.
	ResourceDescription map[string]string `json:"resource_description,omitempty"`
}

// UnmarshalApplicabilityCriteria unmarshals an instance of ApplicabilityCriteria from the specified map of raw messages.
func UnmarshalApplicabilityCriteria(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApplicabilityCriteria)
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource", &obj.Resource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_category", &obj.EnvironmentCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_category", &obj.ResourceCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "software_details", &obj.SoftwareDetails)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "os_details", &obj.OsDetails)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "additional_details", &obj.AdditionalDetails)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_category_description", &obj.EnvironmentCategoryDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_description", &obj.EnvironmentDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_category_description", &obj.ResourceCategoryDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type_description", &obj.ResourceTypeDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_description", &obj.ResourceDescription)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BasicResult : A basic response.
type BasicResult struct {
	// A message.
	Message *string `json:"message" validate:"required"`

	// The result of the operation.
	Result *bool `json:"result" validate:"required"`

	// Id of created Profile.
	ProfileID *string `json:"profile_id,omitempty"`
}

// UnmarshalBasicResult unmarshals an instance of BasicResult from the specified map of raw messages.
func UnmarshalBasicResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BasicResult)
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_id", &obj.ProfileID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Collector : The details of a collector.
type Collector struct {
	// The ID of the collector.
	ID *string `json:"id" validate:"required"`

	// The user-friendly name of the collector.
	DisplayName *string `json:"display_name" validate:"required"`

	// The name of the collector.
	Name *string `json:"name" validate:"required"`

	// The public key of the collector. The key is used for SSL communication between collector and orchestrator. This
	// property is populated when the collector is installed.
	PublicKey *string `json:"public_key,omitempty"`

	// The heartbeat time of the controller. This value exists when the collector is installed and running.
	LastHeartbeat *strfmt.DateTime `json:"last_heartbeat,omitempty"`

	// The status of the collector.
	Status *string `json:"status" validate:"required"`

	// The collector version. This field is populated when the collector is installed.
	CollectorVersion *string `json:"collector_version,omitempty"`

	// The image version of the collector. This field is populated when the collector is installed.".
	ImageVersion *string `json:"image_version,omitempty"`

	// The description of the collector.
	Description *string `json:"description" validate:"required"`

	// The ID of the user that created the collector.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The ISO date and time when the collector was created.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The ID of the user that modified the collector.
	UpdatedBy *string `json:"updated_by" validate:"required"`

	// The ISO date and time when the collector was modified.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The indication of whether the collector is enabled or not(deleted).
	Enabled *bool `json:"enabled" validate:"required"`

	// The registration code of the collector. The code is used for initial authentication during the installation of the
	// collector.
	RegistrationCode *string `json:"registration_code" validate:"required"`

	// The type of the collector.
	Type *string `json:"type" validate:"required"`

	// The credential's public key.
	CredentialPublicKey *string `json:"credential_public_key,omitempty"`

	// The number of times that the collector failed.
	FailureCount *int64 `json:"failure_count" validate:"required"`

	// The approved local gateway IP of the collector. The IP is populated only when the collector is installed.
	ApprovedLocalGatewayIP *string `json:"approved_local_gateway_ip,omitempty"`

	// The approved internet gateway IP of the collector. The IP is populated only when the collector is installed.
	ApprovedInternetGatewayIP *string `json:"approved_internet_gateway_ip,omitempty"`

	// The failed local gateway IP. The IP is populated only when the collector is installed.
	LastFailedLocalGatewayIP *string `json:"last_failed_local_gateway_ip,omitempty"`

	// The reason for the collector reset. User resets the collector with a reason for reset. The reason that is entered by
	// the user is saved in this field .
	ResetReason *string `json:"reset_reason,omitempty"`

	// The collector hostname. The hostname is populated when the collector is installed. The fully qualified domain name
	// is included.
	Hostname *string `json:"hostname,omitempty"`

	// The installation path of the collector. This field is populated when the collector is installed. The value includes
	// the folder path.
	InstallPath *string `json:"install_path,omitempty"`

	// The indication of whether the collector uses a public or private endpoint. This value is generated based on the
	// `is_public` field value during collector creation. If `is_public` is set to true, the `use_private_endpoint` value
	// is false.
	UsePrivateEndpoint *bool `json:"use_private_endpoint" validate:"required"`

	// The entity that manages the collector.
	ManagedBy *string `json:"managed_by" validate:"required"`

	// The trial expiry indicates the expiry date of `registration_code`. This field is populated when the collector is
	// installed.
	TrialExpiry *strfmt.DateTime `json:"trial_expiry,omitempty"`

	// The failed internet gateway IP of the collector.
	LastFailedInternetGatewayIP *string `json:"last_failed_internet_gateway_ip,omitempty"`

	// The collector status.
	StatusDescription *string `json:"status_description" validate:"required"`

	// The ISO date and time of the collector reset. This value is populated when a collector is reset. The data-time when
	// the reset event occurs is captured in this field.
	ResetTime *strfmt.DateTime `json:"reset_time,omitempty"`

	// An indication of whether the collector endpoint is accessible on a public network. If set to `true`, the collector
	// connects to resources in your account over a public network. If set to `false`, the collector connects to resources
	// by using a private IP that is accessible only through the IBM Cloud private network.
	IsPublic *bool `json:"is_public" validate:"required"`

	// An indication of whether the collector has a UBI image.
	IsUbiImage *bool `json:"is_ubi_image,omitempty"`
}

// Constants associated with the Collector.Status property.
// The status of the collector.
const (
	CollectorStatusActiveConst                        = "active"
	CollectorStatusApprovalRequiredConst              = "approval_required"
	CollectorStatusApprovedDownloadInProgressConst    = "approved_download_in_progress"
	CollectorStatusApprovedInstallInProgressConst     = "approved_install_in_progress"
	CollectorStatusCoreDownloadedConst                = "core_downloaded"
	CollectorStatusInstallInProgressConst             = "install_in_progress"
	CollectorStatusInstallationFailedConst            = "installation_failed"
	CollectorStatusInstalledConst                     = "installed"
	CollectorStatusInstalledAssigningCredentialsConst = "installed_assigning_credentials"
	CollectorStatusInstalledCredentialsRequiredConst  = "installed_credentials_required"
	CollectorStatusReadyToInstallConst                = "ready_to_install"
	CollectorStatusSuspendedConst                     = "suspended"
	CollectorStatusUnableToConnectConst               = "unable_to_connect"
	CollectorStatusWaitingForUpgradeConst             = "waiting_for_upgrade"
)

// Constants associated with the Collector.Type property.
// The type of the collector.
const (
	CollectorTypeRestrictedConst   = "restricted"
	CollectorTypeUnrestrictedConst = "unrestricted"
)

// Constants associated with the Collector.ManagedBy property.
// The entity that manages the collector.
const (
	CollectorManagedByCustomerConst = "customer"
	CollectorManagedByIBMConst      = "ibm"
)

// UnmarshalCollector unmarshals an instance of Collector from the specified map of raw messages.
func UnmarshalCollector(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Collector)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "public_key", &obj.PublicKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_heartbeat", &obj.LastHeartbeat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "collector_version", &obj.CollectorVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "image_version", &obj.ImageVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "registration_code", &obj.RegistrationCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "credential_public_key", &obj.CredentialPublicKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_count", &obj.FailureCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "approved_local_gateway_ip", &obj.ApprovedLocalGatewayIP)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "approved_internet_gateway_ip", &obj.ApprovedInternetGatewayIP)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_failed_local_gateway_ip", &obj.LastFailedLocalGatewayIP)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reset_reason", &obj.ResetReason)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hostname", &obj.Hostname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "install_path", &obj.InstallPath)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_private_endpoint", &obj.UsePrivateEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "managed_by", &obj.ManagedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trial_expiry", &obj.TrialExpiry)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_failed_internet_gateway_ip", &obj.LastFailedInternetGatewayIP)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_description", &obj.StatusDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reset_time", &obj.ResetTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_public", &obj.IsPublic)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_ubi_image", &obj.IsUbiImage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CollectorList : The response to a request to list collectors.
type CollectorList struct {
	// The offset from the start of the list (0-based).
	Offset *int64 `json:"offset" validate:"required"`

	// The number of items that are returned.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of items that are in the list. The field's value is 0 when no collectors are available and the
	// details are not populated in that case.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The URL of a page.
	First *PageLink `json:"first" validate:"required"`

	// The URL of a page.
	Last *PageLink `json:"last" validate:"required"`

	// The URL of a page.
	Next *PageLink `json:"next,omitempty"`

	// The URL of a page.
	Previous *PageLink `json:"previous,omitempty"`

	// The array of items that are returned.
	Collectors []Collector `json:"collectors" validate:"required"`
}

// UnmarshalCollectorList unmarshals an instance of CollectorList from the specified map of raw messages.
func UnmarshalCollectorList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CollectorList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "collectors", &obj.Collectors, UnmarshalCollector)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CollectorUpdate : The instance of the collector update.
type CollectorUpdate struct {
	// The display name of the collector.
	DisplayName *string `json:"display_name,omitempty"`

	// The description of the collector.
	Description *string `json:"description,omitempty"`
}

// UnmarshalCollectorUpdate unmarshals an instance of CollectorUpdate from the specified map of raw messages.
func UnmarshalCollectorUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CollectorUpdate)
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
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

// AsPatch returns a generic map representation of the CollectorUpdate
func (collectorUpdate *CollectorUpdate) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(collectorUpdate)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// Control : A scan's summary controls.
type Control struct {
	// The scan summary control ID.
	ID *string `json:"id,omitempty"`

	// The control status.
	Status *string `json:"status,omitempty"`

	// The external control ID.
	ExternalControlID *string `json:"external_control_id,omitempty"`

	// The scan profile name.
	Description *string `json:"description,omitempty"`

	// The list of goals that are on the control.
	Goals []Goal `json:"goals,omitempty"`

	// A scan's summary controls.
	ResourceStatistics *ResourceStatistics `json:"resource_statistics,omitempty"`
}

// Constants associated with the Control.Status property.
// The control status.
const (
	ControlStatusPassConst            = "pass"
	ControlStatusUnableToPerformConst = "unable_to_perform"
)

// UnmarshalControl unmarshals an instance of Control from the specified map of raw messages.
func UnmarshalControl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Control)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "external_control_id", &obj.ExternalControlID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "goals", &obj.Goals, UnmarshalGoal)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource_statistics", &obj.ResourceStatistics, UnmarshalResourceStatistics)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlItem : The details of the profile.
type ControlItem struct {
	// The identifier number of the control.
	ID *string `json:"id" validate:"required"`

	// The description of the control.
	Description *string `json:"description" validate:"required"`

	// The external identifier number of the control.
	ExternalControlID *string `json:"external_control_id" validate:"required"`

	// The goals that are mapped against the control identifier.
	Goals []GoalItem `json:"goals" validate:"required"`
}

// UnmarshalControlItem unmarshals an instance of ControlItem from the specified map of raw messages.
func UnmarshalControlItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "external_control_id", &obj.ExternalControlID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "goals", &obj.Goals, UnmarshalGoalItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlList : The details of the controls for the profile.
type ControlList struct {
	// The offset of the page.
	Offset *int64 `json:"offset" validate:"required"`

	// The number of profiles that are displayed per page.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of profiles. If no profiles are available, the count is 0 and the detail fields are not populated.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The URL of a page.
	First *PageLink `json:"first" validate:"required"`

	// The URL of a page.
	Last *PageLink `json:"last" validate:"required"`

	// The URL of a page.
	Previous *PageLink `json:"previous,omitempty"`

	// The URL of a page.
	Next *PageLink `json:"next,omitempty"`

	// Profiles array.
	Controls []ControlItem `json:"controls" validate:"required"`
}

// UnmarshalControlList unmarshals an instance of ControlList from the specified map of raw messages.
func UnmarshalControlList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalControlItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ControlList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil || offset == nil {
		return nil, err
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// CreateCollectorOptions : The CreateCollector options.
type CreateCollectorOptions struct {
	// A unique name for your collector.
	Name *string `json:"name" validate:"required"`

	// The parameter `is_public` determines whether the collector endpoint is accessible on a public network. If set to
	// `true`, the collector connects to resources that are in your account over a public network. If set to `false`, the
	// collector connects to your resources by using a private IP that is accessible only through the IBM Cloud private
	// network.
	IsPublic *bool `json:"is_public" validate:"required"`

	// The parameter `managed_by` determines whether the collector is an IBM or customer-managed virtual machine. Use `ibm`
	// to allow Security and Compliance Center to create, install, and manage the collector on your behalf. The collector
	// is installed in an Red Hat OpenShift cluster and approved automatically for use. Use `customer` if you would like to
	// install the collector by using your own virtual machine. For more information, check out the
	// [docs](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-collector).
	ManagedBy *string `json:"managed_by" validate:"required"`

	// A detailed description of the collector.
	Description *string `json:"description,omitempty"`

	// To protect the credentials that you add to the service, a passphrase is used to generate a data encryption key. The
	// key is used to securely store your credentials and prevent anyone from accessing them.
	Passphrase *string `json:"passphrase,omitempty"`

	// The parameter `is_ubi_image` determines whether the collector has a UBI image.
	IsUbiImage *bool `json:"is_ubi_image,omitempty"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateCollectorOptions.ManagedBy property.
// The parameter `managed_by` determines whether the collector is an IBM or customer-managed virtual machine. Use `ibm`
// to allow Security and Compliance Center to create, install, and manage the collector on your behalf. The collector is
// installed in an Red Hat OpenShift cluster and approved automatically for use. Use `customer` if you would like to
// install the collector by using your own virtual machine. For more information, check out the
// [docs](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-collector).
const (
	CreateCollectorOptionsManagedByCustomerConst = "customer"
	CreateCollectorOptionsManagedByIBMConst      = "ibm"
)

// NewCreateCollectorOptions : Instantiate CreateCollectorOptions
func (*PostureManagementV2) NewCreateCollectorOptions(name string, isPublic bool, managedBy string) *CreateCollectorOptions {
	return &CreateCollectorOptions{
		Name:      core.StringPtr(name),
		IsPublic:  core.BoolPtr(isPublic),
		ManagedBy: core.StringPtr(managedBy),
	}
}

// SetName : Allow user to set Name
func (_options *CreateCollectorOptions) SetName(name string) *CreateCollectorOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIsPublic : Allow user to set IsPublic
func (_options *CreateCollectorOptions) SetIsPublic(isPublic bool) *CreateCollectorOptions {
	_options.IsPublic = core.BoolPtr(isPublic)
	return _options
}

// SetManagedBy : Allow user to set ManagedBy
func (_options *CreateCollectorOptions) SetManagedBy(managedBy string) *CreateCollectorOptions {
	_options.ManagedBy = core.StringPtr(managedBy)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateCollectorOptions) SetDescription(description string) *CreateCollectorOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetPassphrase : Allow user to set Passphrase
func (_options *CreateCollectorOptions) SetPassphrase(passphrase string) *CreateCollectorOptions {
	_options.Passphrase = core.StringPtr(passphrase)
	return _options
}

// SetIsUbiImage : Allow user to set IsUbiImage
func (_options *CreateCollectorOptions) SetIsUbiImage(isUbiImage bool) *CreateCollectorOptions {
	_options.IsUbiImage = core.BoolPtr(isUbiImage)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateCollectorOptions) SetAccountID(accountID string) *CreateCollectorOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateCollectorOptions) SetTransactionID(transactionID string) *CreateCollectorOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCollectorOptions) SetHeaders(param map[string]string) *CreateCollectorOptions {
	options.Headers = param
	return options
}

// CreateCredentialOptions : The CreateCredential options.
type CreateCredentialOptions struct {
	// The status of credentials is enabled or disabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// The type of credential.
	Type *string `json:"type" validate:"required"`

	// The name of the credential.
	Name *string `json:"name" validate:"required"`

	// The description of the credential.
	Description *string `json:"description" validate:"required"`

	// The details of the credential. The details change as the selected credential type varies.
	DisplayFields *NewCredentialDisplayFields `json:"display_fields" validate:"required"`

	// The purpose for which the credential is created.
	Purpose *string `json:"purpose" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateCredentialOptions.Type property.
// The type of credential.
const (
	CreateCredentialOptionsTypeAwsCloudConst         = "aws_cloud"
	CreateCredentialOptionsTypeAzureCloudConst       = "azure_cloud"
	CreateCredentialOptionsTypeDatabaseConst         = "database"
	CreateCredentialOptionsTypeIBMCloudConst         = "ibm_cloud"
	CreateCredentialOptionsTypeKerberosWindowsConst  = "kerberos_windows"
	CreateCredentialOptionsTypeMs365Const            = "ms_365"
	CreateCredentialOptionsTypeOpenstackCloudConst   = "openstack_cloud"
	CreateCredentialOptionsTypeUserNamePemConst      = "user_name_pem"
	CreateCredentialOptionsTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the CreateCredentialOptions.Purpose property.
// The purpose for which the credential is created.
const (
	CreateCredentialOptionsPurposeDiscoveryCollectionConst                = "discovery_collection"
	CreateCredentialOptionsPurposeDiscoveryCollectionRemediationConst     = "discovery_collection_remediation"
	CreateCredentialOptionsPurposeDiscoveryFactCollectionConst            = "discovery_fact_collection"
	CreateCredentialOptionsPurposeDiscoveryFactCollectionRemediationConst = "discovery_fact_collection_remediation"
	CreateCredentialOptionsPurposeRemediationConst                        = "remediation"
)

// NewCreateCredentialOptions : Instantiate CreateCredentialOptions
func (*PostureManagementV2) NewCreateCredentialOptions(enabled bool, typeVar string, name string, description string, displayFields *NewCredentialDisplayFields, purpose string) *CreateCredentialOptions {
	return &CreateCredentialOptions{
		Enabled:       core.BoolPtr(enabled),
		Type:          core.StringPtr(typeVar),
		Name:          core.StringPtr(name),
		Description:   core.StringPtr(description),
		DisplayFields: displayFields,
		Purpose:       core.StringPtr(purpose),
	}
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateCredentialOptions) SetEnabled(enabled bool) *CreateCredentialOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateCredentialOptions) SetType(typeVar string) *CreateCredentialOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateCredentialOptions) SetName(name string) *CreateCredentialOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateCredentialOptions) SetDescription(description string) *CreateCredentialOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetDisplayFields : Allow user to set DisplayFields
func (_options *CreateCredentialOptions) SetDisplayFields(displayFields *NewCredentialDisplayFields) *CreateCredentialOptions {
	_options.DisplayFields = displayFields
	return _options
}

// SetPurpose : Allow user to set Purpose
func (_options *CreateCredentialOptions) SetPurpose(purpose string) *CreateCredentialOptions {
	_options.Purpose = core.StringPtr(purpose)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateCredentialOptions) SetAccountID(accountID string) *CreateCredentialOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateCredentialOptions) SetTransactionID(transactionID string) *CreateCredentialOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCredentialOptions) SetHeaders(param map[string]string) *CreateCredentialOptions {
	options.Headers = param
	return options
}

// CreateScopeOptions : The CreateScope options.
type CreateScopeOptions struct {
	// A unique name for your scope.
	Name *string `json:"name" validate:"required"`

	// A detailed description of the scope.
	Description *string `json:"description" validate:"required"`

	// The unique IDs of the collectors that are attached to the scope.
	CollectorIds []string `json:"collector_ids" validate:"required"`

	// The unique identifier of the credential.
	CredentialID *string `json:"credential_id" validate:"required"`

	// The environment that the scope is targeted to.
	CredentialType *string `json:"credential_type" validate:"required"`

	// The frequency of the scope. `interval` is used with on-prem scope if the user wants to schedule a discovery task.
	// The unit is seconds. For example, if a user wants to trigger discovery every hour, this value is set to 3600.
	Interval *int64 `json:"interval,omitempty"`

	// The discovery scheduled for the scope. `is_discovery_scheduled` is used with on-prem scope if the user wants to
	// schedule a discovery task.
	IsDiscoveryScheduled *bool `json:"is_discovery_scheduled,omitempty"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateScopeOptions.CredentialType property.
// The environment that the scope is targeted to.
const (
	CreateScopeOptionsCredentialTypeAwsConst       = "aws"
	CreateScopeOptionsCredentialTypeAzureConst     = "azure"
	CreateScopeOptionsCredentialTypeGcpConst       = "gcp"
	CreateScopeOptionsCredentialTypeHostedConst    = "hosted"
	CreateScopeOptionsCredentialTypeIBMConst       = "ibm"
	CreateScopeOptionsCredentialTypeOnPremiseConst = "on_premise"
	CreateScopeOptionsCredentialTypeOpenstackConst = "openstack"
	CreateScopeOptionsCredentialTypeServicesConst  = "services"
)

// NewCreateScopeOptions : Instantiate CreateScopeOptions
func (*PostureManagementV2) NewCreateScopeOptions(name string, description string, collectorIds []string, credentialID string, credentialType string) *CreateScopeOptions {
	return &CreateScopeOptions{
		Name:           core.StringPtr(name),
		Description:    core.StringPtr(description),
		CollectorIds:   collectorIds,
		CredentialID:   core.StringPtr(credentialID),
		CredentialType: core.StringPtr(credentialType),
	}
}

// SetName : Allow user to set Name
func (_options *CreateScopeOptions) SetName(name string) *CreateScopeOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateScopeOptions) SetDescription(description string) *CreateScopeOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetCollectorIds : Allow user to set CollectorIds
func (_options *CreateScopeOptions) SetCollectorIds(collectorIds []string) *CreateScopeOptions {
	_options.CollectorIds = collectorIds
	return _options
}

// SetCredentialID : Allow user to set CredentialID
func (_options *CreateScopeOptions) SetCredentialID(credentialID string) *CreateScopeOptions {
	_options.CredentialID = core.StringPtr(credentialID)
	return _options
}

// SetCredentialType : Allow user to set CredentialType
func (_options *CreateScopeOptions) SetCredentialType(credentialType string) *CreateScopeOptions {
	_options.CredentialType = core.StringPtr(credentialType)
	return _options
}

// SetInterval : Allow user to set Interval
func (_options *CreateScopeOptions) SetInterval(interval int64) *CreateScopeOptions {
	_options.Interval = core.Int64Ptr(interval)
	return _options
}

// SetIsDiscoveryScheduled : Allow user to set IsDiscoveryScheduled
func (_options *CreateScopeOptions) SetIsDiscoveryScheduled(isDiscoveryScheduled bool) *CreateScopeOptions {
	_options.IsDiscoveryScheduled = core.BoolPtr(isDiscoveryScheduled)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateScopeOptions) SetAccountID(accountID string) *CreateScopeOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateScopeOptions) SetTransactionID(transactionID string) *CreateScopeOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateScopeOptions) SetHeaders(param map[string]string) *CreateScopeOptions {
	options.Headers = param
	return options
}

// CreateValidationOptions : The CreateValidation options.
type CreateValidationOptions struct {
	// The unique ID of the scope.
	ScopeID *string `json:"scope_id" validate:"required"`

	// The unique ID of the profile.
	ProfileID *string `json:"profile_id" validate:"required"`

	// The ID of the profile group.
	GroupProfileID *string `json:"group_profile_id,omitempty"`

	// The name of a scheduled scan.This is mandatory when scheduled scan is initiated.
	Name *string `json:"name,omitempty"`

	// The description of a scheduled scan.
	Description *string `json:"description,omitempty"`

	// The frequency of a scheduled scan in milliseconds.
	Frequency *int64 `json:"frequency,omitempty"`

	// The no_of_occurrences of a scheduled scan.
	NoOfOccurrences *int64 `json:"no_of_occurrences,omitempty"`

	// The end date-time of a scheduled scan in UTC.
	EndTime *strfmt.DateTime `json:"end_time,omitempty"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateValidationOptions : Instantiate CreateValidationOptions
func (*PostureManagementV2) NewCreateValidationOptions(scopeID string, profileID string) *CreateValidationOptions {
	return &CreateValidationOptions{
		ScopeID:   core.StringPtr(scopeID),
		ProfileID: core.StringPtr(profileID),
	}
}

// SetScopeID : Allow user to set ScopeID
func (_options *CreateValidationOptions) SetScopeID(scopeID string) *CreateValidationOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *CreateValidationOptions) SetProfileID(profileID string) *CreateValidationOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetGroupProfileID : Allow user to set GroupProfileID
func (_options *CreateValidationOptions) SetGroupProfileID(groupProfileID string) *CreateValidationOptions {
	_options.GroupProfileID = core.StringPtr(groupProfileID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateValidationOptions) SetName(name string) *CreateValidationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateValidationOptions) SetDescription(description string) *CreateValidationOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetFrequency : Allow user to set Frequency
func (_options *CreateValidationOptions) SetFrequency(frequency int64) *CreateValidationOptions {
	_options.Frequency = core.Int64Ptr(frequency)
	return _options
}

// SetNoOfOccurrences : Allow user to set NoOfOccurrences
func (_options *CreateValidationOptions) SetNoOfOccurrences(noOfOccurrences int64) *CreateValidationOptions {
	_options.NoOfOccurrences = core.Int64Ptr(noOfOccurrences)
	return _options
}

// SetEndTime : Allow user to set EndTime
func (_options *CreateValidationOptions) SetEndTime(endTime *strfmt.DateTime) *CreateValidationOptions {
	_options.EndTime = endTime
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateValidationOptions) SetAccountID(accountID string) *CreateValidationOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateValidationOptions) SetTransactionID(transactionID string) *CreateValidationOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateValidationOptions) SetHeaders(param map[string]string) *CreateValidationOptions {
	options.Headers = param
	return options
}

// Credential : Get the credential details.
type Credential struct {
	// The status of the credential is enabled or disabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// The credential's ID.
	ID *string `json:"id" validate:"required"`

	// The credential's type.
	Type *string `json:"type" validate:"required"`

	// The credential's name.
	Name *string `json:"name" validate:"required"`

	// The credential's description.
	Description *string `json:"description" validate:"required"`

	// The details of the credential. The details change as the selected credential type changes.
	DisplayFields *CredentialDisplayFields `json:"display_fields" validate:"required"`

	// The ID of the user who created the credentials.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The time of creation of the credentials in UTC.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The modified time of the credentials in UTC.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The ID of the user who modified the credentials.
	UpdatedBy *string `json:"updated_by" validate:"required"`

	// The purpose for which the credential is created.
	Purpose *string `json:"purpose" validate:"required"`
}

// Constants associated with the Credential.Type property.
// The credential's type.
const (
	CredentialTypeAwsCloudConst         = "aws_cloud"
	CredentialTypeAzureCloudConst       = "azure_cloud"
	CredentialTypeDatabaseConst         = "database"
	CredentialTypeIBMCloudConst         = "ibm_cloud"
	CredentialTypeKerberosWindowsConst  = "kerberos_windows"
	CredentialTypeMs365Const            = "ms_365"
	CredentialTypeOpenstackCloudConst   = "openstack_cloud"
	CredentialTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the Credential.Purpose property.
// The purpose for which the credential is created.
const (
	CredentialPurposeDiscoveryCollectionConst                = "discovery_collection"
	CredentialPurposeDiscoveryCollectionRemediationConst     = "discovery_collection_remediation"
	CredentialPurposeDiscoveryFactCollectionConst            = "discovery_fact_collection"
	CredentialPurposeDiscoveryFactCollectionRemediationConst = "discovery_fact_collection_remediation"
	CredentialPurposeRemediationConst                        = "remediation"
)

// UnmarshalCredential unmarshals an instance of Credential from the specified map of raw messages.
func UnmarshalCredential(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Credential)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
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
	err = core.UnmarshalModel(m, "display_fields", &obj.DisplayFields, UnmarshalCredentialDisplayFields)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}

	err = core.UnmarshalPrimitive(m, "purpose", &obj.Purpose)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CredentialDisplayFields : The details of the credential. The details change as the selected credential type changes.
type CredentialDisplayFields struct {
	// The IBM Cloud API key. The API key is mandatory when IBM is the credential type.
	IBMAPIKey *string `json:"ibm_api_key,omitempty"`

	// The Amazon Web Services client ID. The client ID is mandatory when AWS Cloud is the credential type.
	AwsClientID *string `json:"aws_client_id,omitempty"`

	// The Amazon Web Services client secret. The secret is mandatory when AWS Cloud is the credential type.
	AwsClientSecret *string `json:"aws_client_secret,omitempty"`

	// The Amazon Web Services region. The region is mandatory when AWS Cloud is the credential type.
	AwsRegion *string `json:"aws_region,omitempty"`

	// AWS arn value.
	AwsArn *string `json:"aws_arn,omitempty"`

	// The username of the user. The username is mandatory when the credential type is DataBase, Kerberos, or OpenStack.
	Username *string `json:"username,omitempty"`

	// The password of the user. The password is mandatory when the credential type is DataBase, Kerberos, or OpenStack.
	Password *string `json:"password,omitempty"`

	// The Microsoft Azure client ID. The client ID is mandatory when Azure is the credential type.
	AzureClientID *string `json:"azure_client_id,omitempty"`

	// The Microsoft Azure client secret. The secret is mandatory when Azure is the credential type.
	AzureClientSecret *string `json:"azure_client_secret,omitempty"`

	// The Microsoft Azure subscription ID. The subscription ID is mandatory when Azure is the credential type.
	AzureSubscriptionID *string `json:"azure_subscription_id,omitempty"`

	// The Microsoft Azure resource group. The resource group is mandatory when Azure is the credential type.
	AzureResourceGroup *string `json:"azure_resource_group,omitempty"`

	// The Database name. The name is mandatory when Database is the credential type.
	DatabaseName *string `json:"database_name,omitempty"`

	// The Kerberos Windows auth type. The auth type is mandatory when Kerberos is the credential type.
	WinrmAuthtype *string `json:"winrm_authtype,omitempty"`

	// The Kerberos Windows SSL. The SSL is mandatory when Kerberos is the credential type.
	WinrmUsessl *string `json:"winrm_usessl,omitempty"`

	// The Kerberos Windows port. The port is mandatory when Kerberos is the credential type.
	WinrmPort *string `json:"winrm_port,omitempty"`

	// The Microsoft 365 client ID. The client ID is mandatory when Microsoft 365 is the credential type.
	Ms365ClientID *string `json:"ms_365_client_id,omitempty"`

	// The Microsoft 365 client secret. The secret is mandatory when Microsoft 365 is the credential type.
	Ms365ClientSecret *string `json:"ms_365_client_secret,omitempty"`

	// The Microsoft 365 tenant ID. The tenant ID is mandatory when Microsoft 365 is the credential type.
	Ms365TenantID *string `json:"ms_365_tenant_id,omitempty"`

	// The auth url of the OpenStack cloud. The auth url is mandatory when OpenStack is the credential type.
	AuthURL *string `json:"auth_url,omitempty"`

	// The project name of the OpenStack cloud. The project name is mandatory when OpenStack is the credential type.
	ProjectName *string `json:"project_name,omitempty"`

	// The user domain name of the OpenStack cloud. The domain name is mandatory when OpenStack is the credential type.
	UserDomainName *string `json:"user_domain_name,omitempty"`

	// The project domain name of the OpenStack cloud. The project domain name is mandatory when OpenStack is the
	// credential type.
	ProjectDomainName *string `json:"project_domain_name,omitempty"`

	// The user pem file name.
	PemFileName *string `json:"pem_file_name,omitempty"`

	// The base64 encoded form of pem.Will be displayed a xxxxxx.
	PemData *string `json:"pem_data,omitempty"`
}

// UnmarshalCredentialDisplayFields unmarshals an instance of CredentialDisplayFields from the specified map of raw messages.
func UnmarshalCredentialDisplayFields(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CredentialDisplayFields)
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IBMAPIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_client_id", &obj.AwsClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_client_secret", &obj.AwsClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_arn", &obj.AwsArn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_client_id", &obj.AzureClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_client_secret", &obj.AzureClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "database_name", &obj.DatabaseName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_authtype", &obj.WinrmAuthtype)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_usessl", &obj.WinrmUsessl)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_port", &obj.WinrmPort)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_client_id", &obj.Ms365ClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_client_secret", &obj.Ms365ClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_tenant_id", &obj.Ms365TenantID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_url", &obj.AuthURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "project_name", &obj.ProjectName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_domain_name", &obj.UserDomainName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "project_domain_name", &obj.ProjectDomainName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pem_file_name", &obj.PemFileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pem_data", &obj.PemData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CredentialList : A list of credentials.
type CredentialList struct {
	// The offset of the page.
	Offset *int64 `json:"offset" validate:"required"`

	// The number of credentials that are displayed per page.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of credentials that are in the list. The number is 0 if no credentials are available and the detail
	// fields of the credentials are not populated in that case.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The URL of a page.
	First *PageLink `json:"first" validate:"required"`

	// The URL of a page.
	Last *PageLink `json:"last" validate:"required"`

	// The URL of a page.
	Previous *PageLink `json:"previous,omitempty"`

	// The URL of a page.
	Next *PageLink `json:"next,omitempty"`

	// The details of a credential.
	Credentials []Credential `json:"credentials" validate:"required"`
}

// UnmarshalCredentialList unmarshals an instance of CredentialList from the specified map of raw messages.
func UnmarshalCredentialList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CredentialList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "credentials", &obj.Credentials, UnmarshalCredential)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *CredentialList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil || offset == nil {
		return nil, err
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// DeleteCollectorOptions : The DeleteCollector options.
type DeleteCollectorOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCollectorOptions : Instantiate DeleteCollectorOptions
func (*PostureManagementV2) NewDeleteCollectorOptions(id string) *DeleteCollectorOptions {
	return &DeleteCollectorOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteCollectorOptions) SetID(id string) *DeleteCollectorOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteCollectorOptions) SetAccountID(accountID string) *DeleteCollectorOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteCollectorOptions) SetTransactionID(transactionID string) *DeleteCollectorOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCollectorOptions) SetHeaders(param map[string]string) *DeleteCollectorOptions {
	options.Headers = param
	return options
}

// DeleteCredentialOptions : The DeleteCredential options.
type DeleteCredentialOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCredentialOptions : Instantiate DeleteCredentialOptions
func (*PostureManagementV2) NewDeleteCredentialOptions(id string) *DeleteCredentialOptions {
	return &DeleteCredentialOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteCredentialOptions) SetID(id string) *DeleteCredentialOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteCredentialOptions) SetAccountID(accountID string) *DeleteCredentialOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteCredentialOptions) SetTransactionID(transactionID string) *DeleteCredentialOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCredentialOptions) SetHeaders(param map[string]string) *DeleteCredentialOptions {
	options.Headers = param
	return options
}

// DeleteProfileOptions : The DeleteProfile options.
type DeleteProfileOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteProfileOptions : Instantiate DeleteProfileOptions
func (*PostureManagementV2) NewDeleteProfileOptions(id string) *DeleteProfileOptions {
	return &DeleteProfileOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteProfileOptions) SetID(id string) *DeleteProfileOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteProfileOptions) SetAccountID(accountID string) *DeleteProfileOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteProfileOptions) SetTransactionID(transactionID string) *DeleteProfileOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProfileOptions) SetHeaders(param map[string]string) *DeleteProfileOptions {
	options.Headers = param
	return options
}

// DeleteScopeOptions : The DeleteScope options.
type DeleteScopeOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteScopeOptions : Instantiate DeleteScopeOptions
func (*PostureManagementV2) NewDeleteScopeOptions(id string) *DeleteScopeOptions {
	return &DeleteScopeOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteScopeOptions) SetID(id string) *DeleteScopeOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteScopeOptions) SetAccountID(accountID string) *DeleteScopeOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteScopeOptions) SetTransactionID(transactionID string) *DeleteScopeOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteScopeOptions) SetHeaders(param map[string]string) *DeleteScopeOptions {
	options.Headers = param
	return options
}

// EventItem : Th event details.
type EventItem struct {
	// The event ID for the scope.
	ID *string `json:"id" validate:"required"`

	// The time that the event was created in UTC for this scope.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The time that the event was last updated in UTC for this scope.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The task type. The task type has two values (fact_collection and discovery).
	TaskType *string `json:"task_type" validate:"required"`

	// The status of the event.
	Status *string `json:"status" validate:"required"`

	// Data is available for this event.
	DataAvailable *bool `json:"data_available" validate:"required"`

	// The status of the event.
	StatusMessage *string `json:"status_message" validate:"required"`
}

// Constants associated with the EventItem.TaskType property.
// The task type. The task type has two values (fact_collection and discovery).
const (
	EventItemTaskTypeDiscoveryConst      = "discovery"
	EventItemTaskTypeFactCollectionConst = "fact_collection"
)

// Constants associated with the EventItem.Status property.
// The status of the event.
const (
	EventItemStatusAbortTaskRequestCompletedConst       = "abort_task_request_completed"
	EventItemStatusAbortTaskRequestFailedConst          = "abort_task_request_failed"
	EventItemStatusAbortTaskRequestReceivedConst        = "abort_task_request_received"
	EventItemStatusCertRegularValidationCompletedConst  = "cert_regular_validation_completed"
	EventItemStatusCertRegularValidationErrorConst      = "cert_regular_validation_error"
	EventItemStatusCertRegularValidationStartedConst    = "cert_regular_validation_started"
	EventItemStatusCertValidationCompletedConst         = "cert_validation_completed"
	EventItemStatusCertValidationErrorConst             = "cert_validation_error"
	EventItemStatusCertValidationStartedConst           = "cert_validation_started"
	EventItemStatusControllerAbortedConst               = "controller_aborted"
	EventItemStatusCveRegularValidationCompletedConst   = "cve_regular_validation_completed"
	EventItemStatusCveRegularValidationErrorConst       = "cve_regular_validation_error"
	EventItemStatusCveRegularValidationStartedConst     = "cve_regular_validation_started"
	EventItemStatusCveValidationCompletedConst          = "cve_validation_completed"
	EventItemStatusCveValidationErrorConst              = "cve_validation_error"
	EventItemStatusCveValidationStartedConst            = "cve_validation_started"
	EventItemStatusDiscoveryCompletedConst              = "discovery_completed"
	EventItemStatusDiscoveryInProgressConst             = "discovery_in_progress"
	EventItemStatusDiscoveryResultPostedNoErrorConst    = "discovery_result_posted_no_error"
	EventItemStatusDiscoveryResultPostedWithErrorConst  = "discovery_result_posted_with_error"
	EventItemStatusEolRegularValidationCompletedConst   = "eol_regular_validation_completed"
	EventItemStatusEolRegularValidationErrorConst       = "eol_regular_validation_error"
	EventItemStatusEolRegularValidationStartedConst     = "eol_regular_validation_started"
	EventItemStatusEolValidationCompletedConst          = "eol_validation_completed"
	EventItemStatusEolValidationErrorConst              = "eol_validation_error"
	EventItemStatusEolValidationStartedConst            = "eol_validation_started"
	EventItemStatusErrorInAbortTaskRequestConst         = "error_in_abort_task_request"
	EventItemStatusErrorInDiscoveryConst                = "error_in_discovery"
	EventItemStatusErrorInFactCollectionConst           = "error_in_fact_collection"
	EventItemStatusErrorInFactValidationConst           = "error_in_fact_validation"
	EventItemStatusErrorInInventoryConst                = "error_in_inventory"
	EventItemStatusErrorInRemediationConst              = "error_in_remediation"
	EventItemStatusErrorInValidationConst               = "error_in_validation"
	EventItemStatusFactCollectionCompletedConst         = "fact_collection_completed"
	EventItemStatusFactCollectionInProgressConst        = "fact_collection_in_progress"
	EventItemStatusFactCollectionStartedConst           = "fact_collection_started"
	EventItemStatusFactValidationCompletedConst         = "fact_validation_completed"
	EventItemStatusFactValidationInProgressConst        = "fact_validation_in_progress"
	EventItemStatusFactValidationStartedConst           = "fact_validation_started"
	EventItemStatusGatewayAbortedConst                  = "gateway_aborted"
	EventItemStatusInventoryCompletedConst              = "inventory_completed"
	EventItemStatusInventoryCompletedWithErrorConst     = "inventory_completed_with_error"
	EventItemStatusInventoryInProgressConst             = "inventory_in_progress"
	EventItemStatusInventoryStartedConst                = "inventory_started"
	EventItemStatusLocationChangeAbortedConst           = "location_change_aborted"
	EventItemStatusNotAcceptedConst                     = "not_accepted"
	EventItemStatusPendingConst                         = "pending"
	EventItemStatusRemediationCompletedConst            = "remediation_completed"
	EventItemStatusRemediationInProgressConst           = "remediation_in_progress"
	EventItemStatusRemediationStartedConst              = "remediation_started"
	EventItemStatusSentToCollectorConst                 = "sent_to_collector"
	EventItemStatusUserAbortedConst                     = "user_aborted"
	EventItemStatusValidationCompletedConst             = "validation_completed"
	EventItemStatusValidationInProgressConst            = "validation_in_progress"
	EventItemStatusValidationResultPostedNoErrorConst   = "validation_result_posted_no_error"
	EventItemStatusValidationResultPostedWithErrorConst = "validation_result_posted_with_error"
	EventItemStatusValidationStartedConst               = "validation_started"
	EventItemStatusWaitingForRefineConst                = "waiting_for_refine"
)

// UnmarshalEventItem unmarshals an instance of EventItem from the specified map of raw messages.
func UnmarshalEventItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EventItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "task_type", &obj.TaskType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data_available", &obj.DataAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EventList : The list of events.
type EventList struct {
	// The events for a given scope.
	Events []EventItem `json:"events,omitempty"`
}

// UnmarshalEventList unmarshals an instance of EventList from the specified map of raw messages.
func UnmarshalEventList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EventList)
	err = core.UnmarshalModel(m, "events", &obj.Events, UnmarshalEventItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetCollectorOptions : The GetCollector options.
type GetCollectorOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCollectorOptions : Instantiate GetCollectorOptions
func (*PostureManagementV2) NewGetCollectorOptions(id string) *GetCollectorOptions {
	return &GetCollectorOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetCollectorOptions) SetID(id string) *GetCollectorOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetCollectorOptions) SetAccountID(accountID string) *GetCollectorOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetCollectorOptions) SetTransactionID(transactionID string) *GetCollectorOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCollectorOptions) SetHeaders(param map[string]string) *GetCollectorOptions {
	options.Headers = param
	return options
}

// GetCorrelationIDOptions : The GetCorrelationID options.
type GetCorrelationIDOptions struct {
	// Get the status of a task such as discovery or validation. A correlation ID is created when a scope is created and
	// discovery or validation is triggered for a scope.
	CorrelationID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCorrelationIDOptions : Instantiate GetCorrelationIDOptions
func (*PostureManagementV2) NewGetCorrelationIDOptions(correlationID string) *GetCorrelationIDOptions {
	return &GetCorrelationIDOptions{
		CorrelationID: core.StringPtr(correlationID),
	}
}

// SetCorrelationID : Allow user to set CorrelationID
func (_options *GetCorrelationIDOptions) SetCorrelationID(correlationID string) *GetCorrelationIDOptions {
	_options.CorrelationID = core.StringPtr(correlationID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetCorrelationIDOptions) SetAccountID(accountID string) *GetCorrelationIDOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetCorrelationIDOptions) SetTransactionID(transactionID string) *GetCorrelationIDOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCorrelationIDOptions) SetHeaders(param map[string]string) *GetCorrelationIDOptions {
	options.Headers = param
	return options
}

// GetCredentialOptions : The GetCredential options.
type GetCredentialOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCredentialOptions : Instantiate GetCredentialOptions
func (*PostureManagementV2) NewGetCredentialOptions(id string) *GetCredentialOptions {
	return &GetCredentialOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetCredentialOptions) SetID(id string) *GetCredentialOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetCredentialOptions) SetAccountID(accountID string) *GetCredentialOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetCredentialOptions) SetTransactionID(transactionID string) *GetCredentialOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCredentialOptions) SetHeaders(param map[string]string) *GetCredentialOptions {
	options.Headers = param
	return options
}

// GetGroupProfileControlsOptions : The GetGroupProfileControls options.
type GetGroupProfileControlsOptions struct {
	// The group ID. The ID can be obtained from the profile list API call. In the profile list API call, the records that
	// have type='profile_group' are the groups. The ID of that object displays group_id.
	GroupID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// The offset of the profiles.
	Offset *int64 `json:"-"`

	// The number of profiles that are included per page.
	Limit *int64 `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetGroupProfileControlsOptions : Instantiate GetGroupProfileControlsOptions
func (*PostureManagementV2) NewGetGroupProfileControlsOptions(groupID string) *GetGroupProfileControlsOptions {
	return &GetGroupProfileControlsOptions{
		GroupID: core.StringPtr(groupID),
	}
}

// SetGroupID : Allow user to set GroupID
func (_options *GetGroupProfileControlsOptions) SetGroupID(groupID string) *GetGroupProfileControlsOptions {
	_options.GroupID = core.StringPtr(groupID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetGroupProfileControlsOptions) SetAccountID(accountID string) *GetGroupProfileControlsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetGroupProfileControlsOptions) SetTransactionID(transactionID string) *GetGroupProfileControlsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetGroupProfileControlsOptions) SetOffset(offset int64) *GetGroupProfileControlsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetGroupProfileControlsOptions) SetLimit(limit int64) *GetGroupProfileControlsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetGroupProfileControlsOptions) SetHeaders(param map[string]string) *GetGroupProfileControlsOptions {
	options.Headers = param
	return options
}

// GetProfileControlsOptions : The GetProfileControls options.
type GetProfileControlsOptions struct {
	// The profile ID. The ID can be obtained from the Security and Compliance Center UI by clicking the profile name. The
	// URL contains the ID.
	ProfileID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// The offset of the profiles.
	Offset *int64 `json:"-"`

	// The number of profiles that are included per page.
	Limit *int64 `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProfileControlsOptions : Instantiate GetProfileControlsOptions
func (*PostureManagementV2) NewGetProfileControlsOptions(profileID string) *GetProfileControlsOptions {
	return &GetProfileControlsOptions{
		ProfileID: core.StringPtr(profileID),
	}
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetProfileControlsOptions) SetProfileID(profileID string) *GetProfileControlsOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetProfileControlsOptions) SetAccountID(accountID string) *GetProfileControlsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetProfileControlsOptions) SetTransactionID(transactionID string) *GetProfileControlsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetProfileControlsOptions) SetOffset(offset int64) *GetProfileControlsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetProfileControlsOptions) SetLimit(limit int64) *GetProfileControlsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileControlsOptions) SetHeaders(param map[string]string) *GetProfileControlsOptions {
	options.Headers = param
	return options
}

// GetProfileOptions : The GetProfile options.
type GetProfileOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// The profile type name. The name can be authored/custom/predefined for profiles and profile_group for group profiles.
	ProfileType *string `json:"-" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetProfileOptions.ProfileType property.
// The profile type name. The name can be authored/custom/predefined for profiles and profile_group for group profiles.
const (
	GetProfileOptionsProfileTypeAuthoredConst     = "authored"
	GetProfileOptionsProfileTypeCustomConst       = "custom"
	GetProfileOptionsProfileTypePredefinedConst   = "predefined"
	GetProfileOptionsProfileTypeProfileGroupConst = "profile_group"
)

// NewGetProfileOptions : Instantiate GetProfileOptions
func (*PostureManagementV2) NewGetProfileOptions(id string, profileType string) *GetProfileOptions {
	return &GetProfileOptions{
		ID:          core.StringPtr(id),
		ProfileType: core.StringPtr(profileType),
	}
}

// SetID : Allow user to set ID
func (_options *GetProfileOptions) SetID(id string) *GetProfileOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetProfileType : Allow user to set ProfileType
func (_options *GetProfileOptions) SetProfileType(profileType string) *GetProfileOptions {
	_options.ProfileType = core.StringPtr(profileType)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetProfileOptions) SetAccountID(accountID string) *GetProfileOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetProfileOptions) SetTransactionID(transactionID string) *GetProfileOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileOptions) SetHeaders(param map[string]string) *GetProfileOptions {
	options.Headers = param
	return options
}

// GetScopeDetailsCollectorOptions : The GetScopeDetailsCollector options.
type GetScopeDetailsCollectorOptions struct {
	// The unique identifier that is used to trace an entire Scope request.
	ScopeID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetScopeDetailsCollectorOptions : Instantiate GetScopeDetailsCollectorOptions
func (*PostureManagementV2) NewGetScopeDetailsCollectorOptions(scopeID string) *GetScopeDetailsCollectorOptions {
	return &GetScopeDetailsCollectorOptions{
		ScopeID: core.StringPtr(scopeID),
	}
}

// SetScopeID : Allow user to set ScopeID
func (_options *GetScopeDetailsCollectorOptions) SetScopeID(scopeID string) *GetScopeDetailsCollectorOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetScopeDetailsCollectorOptions) SetAccountID(accountID string) *GetScopeDetailsCollectorOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetScopeDetailsCollectorOptions) SetTransactionID(transactionID string) *GetScopeDetailsCollectorOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetScopeDetailsCollectorOptions) SetHeaders(param map[string]string) *GetScopeDetailsCollectorOptions {
	options.Headers = param
	return options
}

// GetScopeDetailsCredentialsOptions : The GetScopeDetailsCredentials options.
type GetScopeDetailsCredentialsOptions struct {
	// The unique identifier that is used to trace an entire Scope request.
	ScopeID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetScopeDetailsCredentialsOptions : Instantiate GetScopeDetailsCredentialsOptions
func (*PostureManagementV2) NewGetScopeDetailsCredentialsOptions(scopeID string) *GetScopeDetailsCredentialsOptions {
	return &GetScopeDetailsCredentialsOptions{
		ScopeID: core.StringPtr(scopeID),
	}
}

// SetScopeID : Allow user to set ScopeID
func (_options *GetScopeDetailsCredentialsOptions) SetScopeID(scopeID string) *GetScopeDetailsCredentialsOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetScopeDetailsCredentialsOptions) SetAccountID(accountID string) *GetScopeDetailsCredentialsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetScopeDetailsCredentialsOptions) SetTransactionID(transactionID string) *GetScopeDetailsCredentialsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetScopeDetailsCredentialsOptions) SetHeaders(param map[string]string) *GetScopeDetailsCredentialsOptions {
	options.Headers = param
	return options
}

// GetScopeDetailsOptions : The GetScopeDetails options.
type GetScopeDetailsOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetScopeDetailsOptions : Instantiate GetScopeDetailsOptions
func (*PostureManagementV2) NewGetScopeDetailsOptions(id string) *GetScopeDetailsOptions {
	return &GetScopeDetailsOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetScopeDetailsOptions) SetID(id string) *GetScopeDetailsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetScopeDetailsOptions) SetAccountID(accountID string) *GetScopeDetailsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetScopeDetailsOptions) SetTransactionID(transactionID string) *GetScopeDetailsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetScopeDetailsOptions) SetHeaders(param map[string]string) *GetScopeDetailsOptions {
	options.Headers = param
	return options
}

// GetScopeTimelineOptions : The GetScopeTimeline options.
type GetScopeTimelineOptions struct {
	// The unique identifier that is used to trace an entire Scope request.
	ScopeID *string `json:"-" validate:"required,ne="`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetScopeTimelineOptions : Instantiate GetScopeTimelineOptions
func (*PostureManagementV2) NewGetScopeTimelineOptions(scopeID string) *GetScopeTimelineOptions {
	return &GetScopeTimelineOptions{
		ScopeID: core.StringPtr(scopeID),
	}
}

// SetScopeID : Allow user to set ScopeID
func (_options *GetScopeTimelineOptions) SetScopeID(scopeID string) *GetScopeTimelineOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetScopeTimelineOptions) SetAccountID(accountID string) *GetScopeTimelineOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetScopeTimelineOptions) SetTransactionID(transactionID string) *GetScopeTimelineOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetScopeTimelineOptions) SetHeaders(param map[string]string) *GetScopeTimelineOptions {
	options.Headers = param
	return options
}

// Goal : The goals that are on the goals list.
type Goal struct {
	// The description of the goal.
	Description *string `json:"description,omitempty"`

	// The goal ID.
	ID *string `json:"id,omitempty"`

	// The goal status.
	Status *string `json:"status,omitempty"`

	// The severity of the goal.
	Severity *string `json:"severity,omitempty"`

	// The report's time of completion.
	CompletedTime *strfmt.DateTime `json:"completed_time,omitempty"`

	// The error that occurred on goal validation.
	Error *string `json:"error,omitempty"`

	// The list of resource results.
	ResourceResult []ResourceResult `json:"resource_result,omitempty"`

	// The criteria that defines how a profile applies.
	ApplicabilityCriteria *GoalApplicabilityCriteria `json:"applicability_criteria,omitempty"`
}

// Constants associated with the Goal.Status property.
// The goal status.
const (
	GoalStatusFailConst = "fail"
	GoalStatusPassConst = "pass"
)

// UnmarshalGoal unmarshals an instance of Goal from the specified map of raw messages.
func UnmarshalGoal(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Goal)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "severity", &obj.Severity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "completed_time", &obj.CompletedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource_result", &obj.ResourceResult, UnmarshalResourceResult)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "applicability_criteria", &obj.ApplicabilityCriteria, UnmarshalGoalApplicabilityCriteria)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GoalApplicabilityCriteria : The criteria that defines how a profile applies.
type GoalApplicabilityCriteria struct {
	// A list of environments that a profile can be applied to.
	Environment []string `json:"environment,omitempty"`

	// A list of resources that a profile can be used with.
	Resource []string `json:"resource,omitempty"`

	// The type of environment that a profile can be applied to.
	EnvironmentCategory []string `json:"environment_category,omitempty"`

	// The type of resource that a profile can be applied to.
	ResourceCategory []string `json:"resource_category,omitempty"`

	// The resource type that the profile applies to.
	ResourceType []string `json:"resource_type,omitempty"`

	// The software that the profile applies to.
	SoftwareDetails interface{} `json:"software_details,omitempty"`

	// The operating system that the profile applies to.
	OsDetails interface{} `json:"os_details,omitempty"`

	// Any additional details about the profile.
	AdditionalDetails interface{} `json:"additional_details,omitempty"`

	// The type of environment that your scope is targeted to.
	EnvironmentCategoryDescription map[string]string `json:"environment_category_description,omitempty"`

	// The environment that your scope is targeted to.
	EnvironmentDescription map[string]string `json:"environment_description,omitempty"`

	// The type of resource that your scope is targeted to.
	ResourceCategoryDescription map[string]string `json:"resource_category_description,omitempty"`

	// The type of resource that your scope is targeted to.
	ResourceTypeDescription map[string]string `json:"resource_type_description,omitempty"`

	// The resource that is scanned as part of your scope.
	ResourceDescription map[string]string `json:"resource_description,omitempty"`
}

// UnmarshalGoalApplicabilityCriteria unmarshals an instance of GoalApplicabilityCriteria from the specified map of raw messages.
func UnmarshalGoalApplicabilityCriteria(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GoalApplicabilityCriteria)
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource", &obj.Resource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_category", &obj.EnvironmentCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_category", &obj.ResourceCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "software_details", &obj.SoftwareDetails)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "os_details", &obj.OsDetails)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "additional_details", &obj.AdditionalDetails)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_category_description", &obj.EnvironmentCategoryDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_description", &obj.EnvironmentDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_category_description", &obj.ResourceCategoryDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type_description", &obj.ResourceTypeDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_description", &obj.ResourceDescription)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GoalItem : The details of the goal.
type GoalItem struct {
	// The description of the goal.
	Description *string `json:"description" validate:"required"`

	// The goal ID.
	ID *string `json:"id" validate:"required"`

	// The severity of the goal.
	Severity *string `json:"severity" validate:"required"`

	// The goal is manually checked.
	IsManual *bool `json:"is_manual" validate:"required"`

	// The goal is remediable or not.
	IsRemediable *bool `json:"is_remediable" validate:"required"`

	// The goal is reversible or not.
	IsReversible *bool `json:"is_reversible" validate:"required"`

	// The goal is automatable or not.
	IsAutomatable *bool `json:"is_automatable" validate:"required"`

	// The goal is autoremediable or not.
	IsAutoRemediable *bool `json:"is_auto_remediable" validate:"required"`
}

// UnmarshalGoalItem unmarshals an instance of GoalItem from the specified map of raw messages.
func UnmarshalGoalItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GoalItem)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "severity", &obj.Severity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_manual", &obj.IsManual)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_remediable", &obj.IsRemediable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_reversible", &obj.IsReversible)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_automatable", &obj.IsAutomatable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_auto_remediable", &obj.IsAutoRemediable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportProfilesOptions : The ImportProfiles options.
type ImportProfilesOptions struct {
	// The import data file that you want to use to import a profile.
	File io.ReadCloser `json:"-" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewImportProfilesOptions : Instantiate ImportProfilesOptions
func (*PostureManagementV2) NewImportProfilesOptions(file io.ReadCloser) *ImportProfilesOptions {
	return &ImportProfilesOptions{
		File: file,
	}
}

// SetFile : Allow user to set File
func (_options *ImportProfilesOptions) SetFile(file io.ReadCloser) *ImportProfilesOptions {
	_options.File = file
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ImportProfilesOptions) SetAccountID(accountID string) *ImportProfilesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ImportProfilesOptions) SetTransactionID(transactionID string) *ImportProfilesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ImportProfilesOptions) SetHeaders(param map[string]string) *ImportProfilesOptions {
	options.Headers = param
	return options
}

// ListCollectorsOptions : The ListCollectors options.
type ListCollectorsOptions struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListCollectorsOptions : Instantiate ListCollectorsOptions
func (*PostureManagementV2) NewListCollectorsOptions() *ListCollectorsOptions {
	return &ListCollectorsOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListCollectorsOptions) SetAccountID(accountID string) *ListCollectorsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListCollectorsOptions) SetTransactionID(transactionID string) *ListCollectorsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListCollectorsOptions) SetHeaders(param map[string]string) *ListCollectorsOptions {
	options.Headers = param
	return options
}

// ListCredentialsOptions : The ListCredentials options.
type ListCredentialsOptions struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The offset of the profiles.
	Offset *int64 `json:"-"`

	// The number of profiles that are included per page.
	Limit *int64 `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListCredentialsOptions : Instantiate ListCredentialsOptions
func (*PostureManagementV2) NewListCredentialsOptions() *ListCredentialsOptions {
	return &ListCredentialsOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListCredentialsOptions) SetAccountID(accountID string) *ListCredentialsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListCredentialsOptions) SetOffset(offset int64) *ListCredentialsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListCredentialsOptions) SetLimit(limit int64) *ListCredentialsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListCredentialsOptions) SetTransactionID(transactionID string) *ListCredentialsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListCredentialsOptions) SetHeaders(param map[string]string) *ListCredentialsOptions {
	options.Headers = param
	return options
}

// ListLatestScansOptions : The ListLatestScans options.
type ListLatestScansOptions struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// The offset of the profiles.
	Offset *int64 `json:"-"`

	// The number of profiles that are included per page.
	Limit *int64 `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListLatestScansOptions : Instantiate ListLatestScansOptions
func (*PostureManagementV2) NewListLatestScansOptions() *ListLatestScansOptions {
	return &ListLatestScansOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListLatestScansOptions) SetAccountID(accountID string) *ListLatestScansOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListLatestScansOptions) SetTransactionID(transactionID string) *ListLatestScansOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListLatestScansOptions) SetOffset(offset int64) *ListLatestScansOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListLatestScansOptions) SetLimit(limit int64) *ListLatestScansOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListLatestScansOptions) SetHeaders(param map[string]string) *ListLatestScansOptions {
	options.Headers = param
	return options
}

// ListProfilesOptions : The ListProfiles options.
type ListProfilesOptions struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// The offset of the profiles.
	Offset *int64 `json:"-"`

	// The number of profiles that are included per page.
	Limit *int64 `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProfilesOptions : Instantiate ListProfilesOptions
func (*PostureManagementV2) NewListProfilesOptions() *ListProfilesOptions {
	return &ListProfilesOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListProfilesOptions) SetAccountID(accountID string) *ListProfilesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListProfilesOptions) SetTransactionID(transactionID string) *ListProfilesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListProfilesOptions) SetOffset(offset int64) *ListProfilesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListProfilesOptions) SetLimit(limit int64) *ListProfilesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProfilesOptions) SetHeaders(param map[string]string) *ListProfilesOptions {
	options.Headers = param
	return options
}

// ListScopesOptions : The ListScopes options.
type ListScopesOptions struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListScopesOptions : Instantiate ListScopesOptions
func (*PostureManagementV2) NewListScopesOptions() *ListScopesOptions {
	return &ListScopesOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListScopesOptions) SetAccountID(accountID string) *ListScopesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListScopesOptions) SetTransactionID(transactionID string) *ListScopesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListScopesOptions) SetHeaders(param map[string]string) *ListScopesOptions {
	options.Headers = param
	return options
}

// NewCredentialDisplayFields : The details of the credential. The details change as the selected credential type varies.
type NewCredentialDisplayFields struct {
	// The IBM Cloud API Key. The API key is mandatory when IBM is selected as the credential type.
	IBMAPIKey *string `json:"ibm_api_key,omitempty"`

	// The Amazon Web Services client ID. The client ID is mandatory when AWS Cloud is selected as the credential type.
	AwsClientID *string `json:"aws_client_id,omitempty"`

	// The Amazon Web Services client secret. The client secret is mandatory when AWS Cloud is selected as the credential
	// type.
	AwsClientSecret *string `json:"aws_client_secret,omitempty"`

	// The Amazon Web Services region. The region is used when AWS Cloud is selected as the credential type.
	AwsRegion *string `json:"aws_region,omitempty"`

	// The Amazon Web Services arn value. The arn value is used when AWS Cloud is selected as the credential type.
	AwsArn *string `json:"aws_arn,omitempty"`

	// The username of the user. The username is mandatory when the credential type is DataBase, Kerberos, OpenStack, and
	// Username-Password.
	Username *string `json:"username,omitempty"`

	// The password of the user. The password is mandatory when the credential type is DataBase, Kerberos, OpenStack, and
	// Username-Password.
	Password *string `json:"password,omitempty"`

	// The Microsoft Azure client ID. The client ID is mandatory when Azure is selected as the credential type.
	AzureClientID *string `json:"azure_client_id,omitempty"`

	// The Microsoft Azure client secret. The secret is mandatory when the type of credential is set to Azure.
	AzureClientSecret *string `json:"azure_client_secret,omitempty"`

	// The Microsoft Azure subscription ID. The subscription ID is mandatory when the type of credential is set to Azure.
	AzureSubscriptionID *string `json:"azure_subscription_id,omitempty"`

	// The Microsoft Azure resource group. The resource group is used when Azure is the credential type.
	AzureResourceGroup *string `json:"azure_resource_group,omitempty"`

	// The database name. The database name is mandatory when Database is the credential type.
	DatabaseName *string `json:"database_name,omitempty"`

	// The Kerberos Windows authentication type. The authentication type is mandatory when the credential type is Kerberos
	// Windows.
	WinrmAuthtype *string `json:"winrm_authtype,omitempty"`

	// The Kerberos Windows SSL. The SSL is mandatory when the credential type is Kerberos Windows.
	WinrmUsessl *string `json:"winrm_usessl,omitempty"`

	// The Kerberos Windows port. The port is mandatory When Kerberos Windows is the credential type.
	WinrmPort *string `json:"winrm_port,omitempty"`

	// The Microsoft 365 client ID. The client ID is mandatory when Microsoft 365 is the credential type.
	Ms365ClientID *string `json:"ms_365_client_id,omitempty"`

	// The Microsoft 365 client secret. The secret is mandatory when Microsoft 365 is the credential type.
	Ms365ClientSecret *string `json:"ms_365_client_secret,omitempty"`

	// The Microsoft 365 tenant ID. The tenant ID is mandatory when Microsoft 365 is the credential type.
	Ms365TenantID *string `json:"ms_365_tenant_id,omitempty"`

	// The auth url of the OpenStack cloud. The auth url is mandatory when OpenStack is the credential type.
	AuthURL *string `json:"auth_url,omitempty"`

	// The project name of the OpenStack cloud. The project name is mandatory when OpenStack is the credential type.
	ProjectName *string `json:"project_name,omitempty"`

	// The user domain name of the OpenStack cloud. The domain name is mandatory when OpenStack is the credential type.
	UserDomainName *string `json:"user_domain_name,omitempty"`

	// The project domain name of the OpenStack cloud. The project domain name is mandatory when OpenStack is the
	// credential type.
	ProjectDomainName *string `json:"project_domain_name,omitempty"`

	// The user pem file name.
	PemFileName *string `json:"pem_file_name,omitempty"`

	// The base64 encoded form of pem.
	PemData *string `json:"pem_data,omitempty"`
}

// UnmarshalNewCredentialDisplayFields unmarshals an instance of NewCredentialDisplayFields from the specified map of raw messages.
func UnmarshalNewCredentialDisplayFields(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NewCredentialDisplayFields)
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IBMAPIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_client_id", &obj.AwsClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_client_secret", &obj.AwsClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_arn", &obj.AwsArn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_client_id", &obj.AzureClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_client_secret", &obj.AzureClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "database_name", &obj.DatabaseName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_authtype", &obj.WinrmAuthtype)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_usessl", &obj.WinrmUsessl)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_port", &obj.WinrmPort)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_client_id", &obj.Ms365ClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_client_secret", &obj.Ms365ClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_tenant_id", &obj.Ms365TenantID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_url", &obj.AuthURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "project_name", &obj.ProjectName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_domain_name", &obj.UserDomainName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "project_domain_name", &obj.ProjectDomainName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pem_file_name", &obj.PemFileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pem_data", &obj.PemData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageLink : The URL of a page.
type PageLink struct {
	// The URL of a page.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPageLink unmarshals an instance of PageLink from the specified map of raw messages.
func UnmarshalPageLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageLink)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Profile : Profile.
type Profile struct {
	// The name of the profile.
	Name *string `json:"name" validate:"required"`

	// A description of the profile.
	Description *string `json:"description" validate:"required"`

	// The version of the profile.
	Version *int64 `json:"version" validate:"required"`

	// The user who created the profile.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The user who last modified the profile.
	ModifiedBy *string `json:"modified_by" validate:"required"`

	// A reason why you want to delete a profile.
	ReasonForDelete *string `json:"reason_for_delete" validate:"required"`

	// The criteria that defines how a profile applies.
	ApplicabilityCriteria *ApplicabilityCriteria `json:"applicability_criteria" validate:"required"`

	// An auto-generated unique identifying number of the profile.
	ID *string `json:"id" validate:"required"`

	// The base profile that the controls are pulled from.
	BaseProfile *string `json:"base_profile" validate:"required"`

	// The type of profile.
	Type *string `json:"type" validate:"required"`

	// no of Controls.
	NoOfControls *int64 `json:"no_of_controls" validate:"required"`

	// The time that the profile was created in UTC.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The time that the profile was most recently modified in UTC.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The profile status. If the profile is enabled, the value is true. If the profile is disabled, the value is false.
	Enabled *bool `json:"enabled" validate:"required"`
}

// Constants associated with the Profile.Type property.
// The type of profile.
const (
	ProfileTypeAuthoredConst      = "authored"
	ProfileTypeCustomConst        = "custom"
	ProfileTypePredefinedConst    = "predefined"
	ProfileTypeTemplateGroupConst = "template_group"
)

// UnmarshalProfile unmarshals an instance of Profile from the specified map of raw messages.
func UnmarshalProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Profile)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_by", &obj.ModifiedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reason_for_delete", &obj.ReasonForDelete)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "applicability_criteria", &obj.ApplicabilityCriteria, UnmarshalApplicabilityCriteria)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "base_profile", &obj.BaseProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "no_of_controls", &obj.NoOfControls)
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
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileItem : Profile details.
type ProfileItem struct {
	// The name of the profile.
	Name *string `json:"name" validate:"required"`

	// An auto-generated unique identifier for the scope.
	ID *string `json:"id" validate:"required"`

	// The type of profile.
	Type *string `json:"type" validate:"required"`
}

// Constants associated with the ProfileItem.Type property.
// The type of profile.
const (
	ProfileItemTypeCustomConst        = "custom"
	ProfileItemTypePredefinedConst    = "predefined"
	ProfileItemTypeTemplateGroupConst = "template_group"
)

// UnmarshalProfileItem unmarshals an instance of ProfileItem from the specified map of raw messages.
func UnmarshalProfileItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// ProfileList : A list of profiles.
type ProfileList struct {
	// The offset of the page.
	Offset *int64 `json:"offset" validate:"required"`

	// The number of profiles displayed per page.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of profiles. If no profiles are available, the value of this field is 0.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The URL of a page.
	First *PageLink `json:"first" validate:"required"`

	// The URL of a page.
	Last *PageLink `json:"last" validate:"required"`

	// The URL of a page.
	Previous *PageLink `json:"previous,omitempty"`

	// The URL of a page.
	Next *PageLink `json:"next,omitempty"`

	// Profiles.
	Profiles []Profile `json:"profiles" validate:"required"`
}

// UnmarshalProfileList unmarshals an instance of ProfileList from the specified map of raw messages.
func UnmarshalProfileList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalProfile)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ProfileList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil || offset == nil {
		return nil, err
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// ProfileResult : The result of a profile.
type ProfileResult struct {
	// The ID of the profile.
	ID *string `json:"id" validate:"required"`

	// The name of the profile.
	Name *string `json:"name" validate:"required"`

	// The type of profile. To learn more about profile types, check out the
	// [docs](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-profiles).
	Type *string `json:"type" validate:"required"`

	// The result of a scan. The controls values are not available if no scopes are available.
	ValidationResult *ScanResult `json:"validation_result" validate:"required"`
}

// Constants associated with the ProfileResult.Type property.
// The type of profile. To learn more about profile types, check out the
// [docs](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-profiles).
const (
	ProfileResultTypeAuthoredConst            = "authored"
	ProfileResultTypeCustomConst              = "custom"
	ProfileResultTypePredefinedConst          = "predefined"
	ProfileResultTypeStandardConst            = "standard"
	ProfileResultTypeStandardCertificateConst = "standard_certificate"
	ProfileResultTypeStandardCvConst          = "standard_cv"
	ProfileResultTypeTemmplategroupConst      = "temmplategroup"
)

// UnmarshalProfileResult unmarshals an instance of ProfileResult from the specified map of raw messages.
func UnmarshalProfileResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileResult)
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
	err = core.UnmarshalModel(m, "validation_result", &obj.ValidationResult, UnmarshalScanResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceScopeDetailsCollectorOptions : The ReplaceScopeDetailsCollector options.
type ReplaceScopeDetailsCollectorOptions struct {
	// The unique identifier that is used to trace an entire Scope request.
	ScopeID *string `json:"-" validate:"required,ne="`

	// The collector IDs of the scope.
	CollectorIds []string `json:"collector_ids" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceScopeDetailsCollectorOptions : Instantiate ReplaceScopeDetailsCollectorOptions
func (*PostureManagementV2) NewReplaceScopeDetailsCollectorOptions(scopeID string, collectorIds []string) *ReplaceScopeDetailsCollectorOptions {
	return &ReplaceScopeDetailsCollectorOptions{
		ScopeID:      core.StringPtr(scopeID),
		CollectorIds: collectorIds,
	}
}

// SetScopeID : Allow user to set ScopeID
func (_options *ReplaceScopeDetailsCollectorOptions) SetScopeID(scopeID string) *ReplaceScopeDetailsCollectorOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetCollectorIds : Allow user to set CollectorIds
func (_options *ReplaceScopeDetailsCollectorOptions) SetCollectorIds(collectorIds []string) *ReplaceScopeDetailsCollectorOptions {
	_options.CollectorIds = collectorIds
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceScopeDetailsCollectorOptions) SetAccountID(accountID string) *ReplaceScopeDetailsCollectorOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ReplaceScopeDetailsCollectorOptions) SetTransactionID(transactionID string) *ReplaceScopeDetailsCollectorOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceScopeDetailsCollectorOptions) SetHeaders(param map[string]string) *ReplaceScopeDetailsCollectorOptions {
	options.Headers = param
	return options
}

// ReplaceScopeDetailsCredentialsOptions : The ReplaceScopeDetailsCredentials options.
type ReplaceScopeDetailsCredentialsOptions struct {
	// The unique identifier that is used to trace an entire Scope request.
	ScopeID *string `json:"-" validate:"required,ne="`

	// The credential ID of the scope.
	CredentialID *string `json:"credential_id" validate:"required"`

	// The credential attribute of the scope.
	CredentialAttribute *string `json:"credential_attribute,omitempty"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceScopeDetailsCredentialsOptions : Instantiate ReplaceScopeDetailsCredentialsOptions
func (*PostureManagementV2) NewReplaceScopeDetailsCredentialsOptions(scopeID string, credentialID string) *ReplaceScopeDetailsCredentialsOptions {
	return &ReplaceScopeDetailsCredentialsOptions{
		ScopeID:      core.StringPtr(scopeID),
		CredentialID: core.StringPtr(credentialID),
	}
}

// SetScopeID : Allow user to set ScopeID
func (_options *ReplaceScopeDetailsCredentialsOptions) SetScopeID(scopeID string) *ReplaceScopeDetailsCredentialsOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetCredentialID : Allow user to set CredentialID
func (_options *ReplaceScopeDetailsCredentialsOptions) SetCredentialID(credentialID string) *ReplaceScopeDetailsCredentialsOptions {
	_options.CredentialID = core.StringPtr(credentialID)
	return _options
}

// SetCredentialAttribute : Allow user to set CredentialAttribute
func (_options *ReplaceScopeDetailsCredentialsOptions) SetCredentialAttribute(credentialAttribute string) *ReplaceScopeDetailsCredentialsOptions {
	_options.CredentialAttribute = core.StringPtr(credentialAttribute)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceScopeDetailsCredentialsOptions) SetAccountID(accountID string) *ReplaceScopeDetailsCredentialsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ReplaceScopeDetailsCredentialsOptions) SetTransactionID(transactionID string) *ReplaceScopeDetailsCredentialsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceScopeDetailsCredentialsOptions) SetHeaders(param map[string]string) *ReplaceScopeDetailsCredentialsOptions {
	options.Headers = param
	return options
}

// ResourceResult : The resource results.
type ResourceResult struct {
	// The resource name.
	Name *string `json:"name,omitempty"`

	// The resource type.
	Types *string `json:"types,omitempty"`

	// The result status of resource control.
	Status *string `json:"status,omitempty"`

	// The expected results of a resource.
	DisplayExpectedValue *string `json:"display_expected_value,omitempty"`

	// The actual results of a resource.
	ActualValue *string `json:"actual_value,omitempty"`

	// The results information.
	ResultsInfo *string `json:"results_info,omitempty"`

	// The reason why a goal is not applicable to a resource.
	NotApplicableReason *string `json:"not_applicable_reason,omitempty"`
}

// Constants associated with the ResourceResult.Status property.
// The result status of resource control.
const (
	ResourceResultStatusPassConst            = "pass"
	ResourceResultStatusUnableToPerformConst = "unable_to_perform"
)

// UnmarshalResourceResult unmarshals an instance of ResourceResult from the specified map of raw messages.
func UnmarshalResourceResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceResult)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "types", &obj.Types)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_expected_value", &obj.DisplayExpectedValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "actual_value", &obj.ActualValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "results_info", &obj.ResultsInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_applicable_reason", &obj.NotApplicableReason)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceStatistics : A scan's summary controls.
type ResourceStatistics struct {
	// The resource count of pass controls.
	PassCount *int64 `json:"pass_count,omitempty"`

	// The resource count of fail controls.
	FailCount *int64 `json:"fail_count,omitempty"`

	// The number of resources that were unable to be scanned against a control.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The resource count of controls that are not applicable.
	NotApplicableCount *int64 `json:"not_applicable_count,omitempty"`
}

// UnmarshalResourceStatistics unmarshals an instance of ResourceStatistics from the specified map of raw messages.
func UnmarshalResourceStatistics(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceStatistics)
	err = core.UnmarshalPrimitive(m, "pass_count", &obj.PassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "fail_count", &obj.FailCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_applicable_count", &obj.NotApplicableCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Result : Result.
type Result struct {
	// Result.
	Result *bool `json:"result" validate:"required"`

	// A message is returned.
	Message *string `json:"message" validate:"required"`
}

// UnmarshalResult unmarshals an instance of Result from the specified map of raw messages.
func UnmarshalResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Result)
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
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

// ScanItem : The details of a scan.
type ScanItem struct {
	// The ID of the scan.
	ScanID *string `json:"scan_id" validate:"required"`

	// A system-generated name that is the combination of 12 characters in the scope name and 12 characters of a profile
	// name.
	ScanName *string `json:"scan_name" validate:"required"`

	// The scope ID of the scan.
	ScopeID *string `json:"scope_id" validate:"required"`

	// The name of the scope.
	ScopeName *string `json:"scope_name" validate:"required"`

	// Profiles array.
	Profiles []ProfileItem `json:"profiles,omitempty"`

	// The group ID of profile.
	GroupProfileID *string `json:"group_profile_id" validate:"required"`

	// The group name of the profile.
	GroupProfileName *string `json:"group_profile_name" validate:"required"`

	// The entity that ran the report.
	ReportRunBy *string `json:"report_run_by" validate:"required"`

	// The date and time the scan was run.
	StartTime *strfmt.DateTime `json:"start_time" validate:"required"`

	// The unique ID for the scan that is created.
	ReportSettingID *string `json:"report_setting_id,omitempty"`

	// The date and time when the scan completed.
	EndTime *strfmt.DateTime `json:"end_time" validate:"required"`

	// The result of a scan. The controls values are not available if no scopes are available.
	Result *ScanResult `json:"result" validate:"required"`
}

// UnmarshalScanItem unmarshals an instance of ScanItem from the specified map of raw messages.
func UnmarshalScanItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScanItem)
	err = core.UnmarshalPrimitive(m, "scan_id", &obj.ScanID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scan_name", &obj.ScanName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_id", &obj.ScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_name", &obj.ScopeName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalProfileItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "group_profile_id", &obj.GroupProfileID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "group_profile_name", &obj.GroupProfileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_run_by", &obj.ReportRunBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_setting_id", &obj.ReportSettingID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "end_time", &obj.EndTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalScanResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScanList : A list of scans.
type ScanList struct {
	// The offset of the page.
	Offset *int64 `json:"offset" validate:"required"`

	// The number of scans that are displayed per page.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of scans that are in the list. This value is 0 when no scans are available and the detail fields
	// are not displayed in that case.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The URL of a page.
	First *PageLink `json:"first" validate:"required"`

	// The URL of a page.
	Last *PageLink `json:"last" validate:"required"`

	// The URL of a page.
	Previous *PageLink `json:"previous,omitempty"`

	// The URL of a page.
	Next *PageLink `json:"next,omitempty"`

	// The details of a scan.
	LatestScans []ScanItem `json:"latest_scans" validate:"required"`
}

// UnmarshalScanList unmarshals an instance of ScanList from the specified map of raw messages.
func UnmarshalScanList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScanList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "latest_scans", &obj.LatestScans, UnmarshalScanItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ScanList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil || offset == nil {
		return nil, err
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// ScanResult : The result of a scan. The controls values are not available if no scopes are available.
type ScanResult struct {
	// The number of goals that passed the scan.
	GoalsPassCount *int64 `json:"goals_pass_count" validate:"required"`

	// The number of goals that can't be validated. A control is listed as 'Unable to perform' when information about its
	// associated resource can't be collected.
	GoalsUnableToPerformCount *int64 `json:"goals_unable_to_perform_count" validate:"required"`

	// The number of goals that are not relevant to the current scan. A scan is listed as 'Not applicable' when information
	// about its associated resource can't be found.
	GoalsNotApplicableCount *int64 `json:"goals_not_applicable_count" validate:"required"`

	// The number of goals that failed the scan.
	GoalsFailCount *int64 `json:"goals_fail_count" validate:"required"`

	// The total number of goals that were included in the scan.
	GoalsTotalCount *int64 `json:"goals_total_count" validate:"required"`

	// The number of controls that passed the scan.
	ControlsPassCount *int64 `json:"controls_pass_count" validate:"required"`

	// The number of controls that failed the scan.
	ControlsFailCount *int64 `json:"controls_fail_count" validate:"required"`

	// The number of controls that are not relevant to the current scan. A scan is listed as 'Not applicable' when
	// information about its associated resource can't be found.
	ControlsNotApplicableCount *int64 `json:"controls_not_applicable_count" validate:"required"`

	// The number of controls that can't be validated. A control is listed as 'Unable to perform' when information about
	// its associated resource can't be collected.
	ControlsUnableToPerformCount *int64 `json:"controls_unable_to_perform_count" validate:"required"`

	// The total number of controls that are included in the scan.
	ControlsTotalCount *int64 `json:"controls_total_count" validate:"required"`
}

// UnmarshalScanResult unmarshals an instance of ScanResult from the specified map of raw messages.
func UnmarshalScanResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScanResult)
	err = core.UnmarshalPrimitive(m, "goals_pass_count", &obj.GoalsPassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "goals_unable_to_perform_count", &obj.GoalsUnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "goals_not_applicable_count", &obj.GoalsNotApplicableCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "goals_fail_count", &obj.GoalsFailCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "goals_total_count", &obj.GoalsTotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_pass_count", &obj.ControlsPassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_fail_count", &obj.ControlsFailCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_not_applicable_count", &obj.ControlsNotApplicableCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_unable_to_perform_count", &obj.ControlsUnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_total_count", &obj.ControlsTotalCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScanSummariesOptions : The ScanSummaries options.
type ScanSummariesOptions struct {
	// The report setting ID. The ID can be obtained from the /validations/latest_scans API call.
	ReportSettingID *string `json:"-" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// The offset of the profiles.
	Offset *int64 `json:"-"`

	// The number of profiles that are included per page.
	Limit *int64 `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewScanSummariesOptions : Instantiate ScanSummariesOptions
func (*PostureManagementV2) NewScanSummariesOptions(reportSettingID string) *ScanSummariesOptions {
	return &ScanSummariesOptions{
		ReportSettingID: core.StringPtr(reportSettingID),
	}
}

// SetReportSettingID : Allow user to set ReportSettingID
func (_options *ScanSummariesOptions) SetReportSettingID(reportSettingID string) *ScanSummariesOptions {
	_options.ReportSettingID = core.StringPtr(reportSettingID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ScanSummariesOptions) SetAccountID(accountID string) *ScanSummariesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ScanSummariesOptions) SetTransactionID(transactionID string) *ScanSummariesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ScanSummariesOptions) SetOffset(offset int64) *ScanSummariesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ScanSummariesOptions) SetLimit(limit int64) *ScanSummariesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ScanSummariesOptions) SetHeaders(param map[string]string) *ScanSummariesOptions {
	options.Headers = param
	return options
}

// ScansSummaryOptions : The ScansSummary options.
type ScansSummaryOptions struct {
	// Your Scan ID.
	ScanID *string `json:"-" validate:"required,ne="`

	// The profile ID. The ID can be obtained from the Security and Compliance Center UI by clicking the profile name. The
	// URL contains the ID.
	ProfileID *string `json:"-" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewScansSummaryOptions : Instantiate ScansSummaryOptions
func (*PostureManagementV2) NewScansSummaryOptions(scanID string, profileID string) *ScansSummaryOptions {
	return &ScansSummaryOptions{
		ScanID:    core.StringPtr(scanID),
		ProfileID: core.StringPtr(profileID),
	}
}

// SetScanID : Allow user to set ScanID
func (_options *ScansSummaryOptions) SetScanID(scanID string) *ScansSummaryOptions {
	_options.ScanID = core.StringPtr(scanID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ScansSummaryOptions) SetProfileID(profileID string) *ScansSummaryOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ScansSummaryOptions) SetAccountID(accountID string) *ScansSummaryOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ScansSummaryOptions) SetTransactionID(transactionID string) *ScansSummaryOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ScansSummaryOptions) SetHeaders(param map[string]string) *ScansSummaryOptions {
	options.Headers = param
	return options
}

// Scope : Get the details of the scope.
type Scope struct {
	// The ID of the scope.
	ID *string `json:"id" validate:"required"`

	// The name of the scope.
	Name *string `json:"name" validate:"required"`

	// The UUID of the scope. The unique user ID is displayed only when the value exists.
	UUID *string `json:"uuid,omitempty"`

	// The partner_uuid of the scope. The ID is displayed only when the value exists.
	PartnerUUID *string `json:"partner_uuid,omitempty"`

	// The description of the scope. The description is displayed only when the value exists.
	Description *string `json:"description,omitempty"`

	// The organization ID of the scope. The organization ID is displayed only when the value exists.
	OrgID *int64 `json:"org_id,omitempty"`

	// The cloud type ID of the scope. The ID is displayed only when the value exists.
	CloudTypeID *int64 `json:"cloud_type_id,omitempty"`

	// The credential ID of the scope. The ID is displayed only when the value exists.
	TldCredentialID *int64 `json:"tld_credential_id,omitempty"`

	// The status of the scope. The status is displayed only when the value exists.
	Status *string `json:"status,omitempty"`

	// The status message of the scope. The message is displayed only when the value exists.
	StatusMsg *string `json:"status_msg,omitempty"`

	// The subset-selected field of the scope. The field is displayed only when the value exists.
	SubsetSelected *bool `json:"subset_selected,omitempty"`

	// The enabled field of the scope. The field is displayed only when the value exists.
	Enabled *bool `json:"enabled,omitempty"`

	// The last discover start time of the scope. The time is displayed only when the value exists.
	LastDiscoverStartTime *string `json:"last_discover_start_time,omitempty"`

	// The last discover completed time of the scope. The time is displayed only when the value exists.
	LastDiscoverCompletedTime *string `json:"last_discover_completed_time,omitempty"`

	// The last successful discover start time of the scope. The time is displayed only when the value exists.
	LastSuccessfulDiscoverStartTime *string `json:"last_successful_discover_start_time,omitempty"`

	// The last successful discover completed time of the scope. The time is displayed only when the value exists.
	LastSuccessfulDiscoverCompletedTime *string `json:"last_successful_discover_completed_time,omitempty"`

	// The task type of the scope. The task type is displayed only when the value exists.
	TaskType *string `json:"task_type,omitempty"`

	// The tasks of the scope. The tasks are displayed only when the value exists.
	Tasks []ScopeDetailsGatewayTask `json:"tasks,omitempty"`

	// The status updated time of the scope. The time is displayed only when the value exists.
	StatusUpdatedTime *string `json:"status_updated_time,omitempty"`

	// The collectors by type of the scope. The collectors are displayed only when the values exist.
	CollectorsByType map[string][]Collector `json:"collectors_by_type,omitempty"`

	// The credentials by type of the scope. The credentials are displayed only when the values exist.
	CredentialsByType map[string][]ScopeDetailsCredential `json:"credentials_by_type,omitempty"`

	// The credentials by sub category type of the scope. The credentials are displayed only when the values exist.
	CredentialsBySubCategeoryType map[string][]ScopeDetailsCredential `json:"credentials_by_sub_categeory_type,omitempty"`

	// The sub categories by type of the scope. The categories are displayed only when the values exist.
	SubCategoriesByType map[string][]string `json:"sub_categories_by_type,omitempty"`

	// The resource groups of the scope. The resource groups are displayed only when the values exist.
	ResourceGroups *string `json:"resource_groups,omitempty"`

	// The region names of the scope. The names are displayed only when the values exist.
	RegionNames *string `json:"region_names,omitempty"`

	// The cloud type of the scope. The cloud type is displayed only when the value exists.
	CloudType *string `json:"cloud_type,omitempty"`

	// The env sub category of the scope. The category is displayed only when the value exists.
	EnvSubCategory *string `json:"env_sub_category,omitempty"`

	// The credential details of the scope.
	TldCredentail *ScopeDetailsCredential `json:"tld_credentail,omitempty"`

	// The collectors of the scope. The collectors are displayed only when the values exist.
	Collectors []Collector `json:"collectors,omitempty"`

	// The first-level scoped data of the scope. The data is displayed only when the value exists.
	FirstLevelScopedData []ScopeDetailsAssetData `json:"first_level_scoped_data,omitempty"`

	// The discovery methods of the scope. The methods are displayed only when the values exist.
	DiscoveryMethods []string `json:"discovery_methods,omitempty"`

	// The discovery method of the scope. The method is displayed only when the value exists.
	DiscoveryMethod *string `json:"discovery_method,omitempty"`

	// The file type of the scope. The type is displayed only when the value exists.
	FileType *string `json:"file_type,omitempty"`

	// The file format of the scope. The format is displayed only when the value exists.
	FileFormat *string `json:"file_format,omitempty"`

	// The user who created the scope. The user is displayed only when the value exists.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date of creation of the scope. The date is displayed only when the value exists.
	CreatedAt *string `json:"created_at,omitempty"`

	// The user who modified the scope. The user is displayed only when the value exists.
	ModifiedBy *string `json:"modified_by,omitempty"`

	// The modified date of the scope. The date is displayed only when the value exists.
	ModifiedAt *string `json:"modified_at,omitempty"`

	// The scheduled discovery configuration of the scope. The data is displayed only when the value exists.
	IsDiscoveryScheduled *bool `json:"is_discovery_scheduled,omitempty"`

	// The interval of the scope. The interval is displayed only when the value exists.
	Interval *int64 `json:"interval,omitempty"`

	// The discovery setting ID of the scope. The ID is displayed only when the value exists.
	DiscoverySettingID *int64 `json:"discovery_setting_id,omitempty"`

	// The `include_new_eagerly` of the scope. The data is displayed only when the value exists.
	IncludeNewEagerly *bool `json:"include_new_eagerly,omitempty"`

	// The type of the scope. The type is displayed only when the value exists.
	Type *string `json:"type,omitempty"`

	// Get the status of a task such as discovery or validation. A correlation ID is created when a scope is created and
	// discovery or validation is triggered for a scope.
	CorrelationID *string `json:"correlation_id,omitempty"`

	// The credential attributes of the scope. The attributes are displayed only when the value exists.
	CredentialAttributes *string `json:"credential_attributes,omitempty"`
}

// Constants associated with the Scope.Status property.
// The status of the scope. The status is displayed only when the value exists.
const (
	ScopeStatusAbortTaskRequestCompletedConst       = "abort_task_request_completed"
	ScopeStatusAbortTaskRequestFailedConst          = "abort_task_request_failed"
	ScopeStatusAbortTaskRequestReceivedConst        = "abort_task_request_received"
	ScopeStatusCertRegularValidationCompletedConst  = "cert_regular_validation_completed"
	ScopeStatusCertRegularValidationErrorConst      = "cert_regular_validation_error"
	ScopeStatusCertRegularValidationStartedConst    = "cert_regular_validation_started"
	ScopeStatusCertValidationCompletedConst         = "cert_validation_completed"
	ScopeStatusCertValidationErrorConst             = "cert_validation_error"
	ScopeStatusCertValidationStartedConst           = "cert_validation_started"
	ScopeStatusControllerAbortedConst               = "controller_aborted"
	ScopeStatusCveRegularValidationCompletedConst   = "cve_regular_validation_completed"
	ScopeStatusCveRegularValidationErrorConst       = "cve_regular_validation_error"
	ScopeStatusCveRegularValidationStartedConst     = "cve_regular_validation_started"
	ScopeStatusCveValidationCompletedConst          = "cve_validation_completed"
	ScopeStatusCveValidationErrorConst              = "cve_validation_error"
	ScopeStatusCveValidationStartedConst            = "cve_validation_started"
	ScopeStatusDiscoveryCompletedConst              = "discovery_completed"
	ScopeStatusDiscoveryInProgressConst             = "discovery_in_progress"
	ScopeStatusDiscoveryResultPostedNoErrorConst    = "discovery_result_posted_no_error"
	ScopeStatusDiscoveryResultPostedWithErrorConst  = "discovery_result_posted_with_error"
	ScopeStatusDiscoveryStartedConst                = "discovery_started"
	ScopeStatusEolRegularValidationCompletedConst   = "eol_regular_validation_completed"
	ScopeStatusEolRegularValidationErrorConst       = "eol_regular_validation_error"
	ScopeStatusEolRegularValidationStartedConst     = "eol_regular_validation_started"
	ScopeStatusEolValidationCompletedConst          = "eol_validation_completed"
	ScopeStatusEolValidationErrorConst              = "eol_validation_error"
	ScopeStatusEolValidationStartedConst            = "eol_validation_started"
	ScopeStatusErrorInAbortTaskRequestConst         = "error_in_abort_task_request"
	ScopeStatusErrorInDiscoverConst                 = "error_in_discover"
	ScopeStatusErrorInFactCollectionConst           = "error_in_fact_collection"
	ScopeStatusErrorInFactValidationConst           = "error_in_fact_validation"
	ScopeStatusErrorInInventoryConst                = "error_in_inventory"
	ScopeStatusErrorInRemediationConst              = "error_in_remediation"
	ScopeStatusErrorInValidationConst               = "error_in_validation"
	ScopeStatusFactCollectionCompletedConst         = "fact_collection_completed"
	ScopeStatusFactCollectionInProgressConst        = "fact_collection_in_progress"
	ScopeStatusFactCollectionStartedConst           = "fact_collection_started"
	ScopeStatusFactValidationCompletedConst         = "fact_validation_completed"
	ScopeStatusFactValidationInProgressConst        = "fact_validation_in_progress"
	ScopeStatusFactValidationStartedConst           = "fact_validation_started"
	ScopeStatusGatewayAbortedConst                  = "gateway_aborted"
	ScopeStatusInventoryCompletedConst              = "inventory_completed"
	ScopeStatusInventoryCompletedWithErrorConst     = "inventory_completed_with_error"
	ScopeStatusInventoryInProgressConst             = "inventory_in_progress"
	ScopeStatusInventoryStartedConst                = "inventory_started"
	ScopeStatusLocationChangeAbortedConst           = "location_change_aborted"
	ScopeStatusNotAcceptedConst                     = "not_accepted"
	ScopeStatusPendingConst                         = "pending"
	ScopeStatusRemediationCompletedConst            = "remediation_completed"
	ScopeStatusRemediationInProgressConst           = "remediation_in_progress"
	ScopeStatusRemediationStartedConst              = "remediation_started"
	ScopeStatusSentToCollectorConst                 = "sent_to_collector"
	ScopeStatusUserAbortedConst                     = "user_aborted"
	ScopeStatusValidationCompletedConst             = "validation_completed"
	ScopeStatusValidationInProgressConst            = "validation_in_progress"
	ScopeStatusValidationResultPostedNoErrorConst   = "validation_result_posted_no_error"
	ScopeStatusValidationResultPostedWithErrorConst = "validation_result_posted_with_error"
	ScopeStatusValidationStartedConst               = "validation_started"
	ScopeStatusWaitingForRefineConst                = "waiting_for_refine"
)

// Constants associated with the Scope.TaskType property.
// The task type of the scope. The task type is displayed only when the value exists.
const (
	ScopeTaskTypeAborttasksConst            = "aborttasks"
	ScopeTaskTypeCertRegularValidationConst = "cert_regular_validation"
	ScopeTaskTypeCertValidationConst        = "cert_validation"
	ScopeTaskTypeCveRegularValidationConst  = "cve_regular_validation"
	ScopeTaskTypeCveValidationConst         = "cve_validation"
	ScopeTaskTypeDiscoverConst              = "discover"
	ScopeTaskTypeEolRegularValidationConst  = "eol_regular_validation"
	ScopeTaskTypeEolValidationConst         = "eol_validation"
	ScopeTaskTypeEvidenceConst              = "evidence"
	ScopeTaskTypeFactcollectionConst        = "factcollection"
	ScopeTaskTypeFactvalidationConst        = "factvalidation"
	ScopeTaskTypeInventoryConst             = "inventory"
	ScopeTaskTypeNopConst                   = "nop"
	ScopeTaskTypeRemediationConst           = "remediation"
	ScopeTaskTypeScriptConst                = "script"
	ScopeTaskTypeSubsetvalidateConst        = "subsetvalidate"
	ScopeTaskTypeTldiscoverConst            = "tldiscover"
)

// Constants associated with the Scope.SubCategoriesByType property.
// The sub categories by type of the scope. The categories are displayed only when the values exist.
const (
	ScopeSubCategoriesByTypeMs365Const = "ms_365"
)

// Constants associated with the Scope.Type property.
// The type of the scope. The type is displayed only when the value exists.
const (
	ScopeTypeInventoryConst  = "inventory"
	ScopeTypeValidationConst = "validation"
)

// UnmarshalScope unmarshals an instance of Scope from the specified map of raw messages.
func UnmarshalScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Scope)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uuid", &obj.UUID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "partner_uuid", &obj.PartnerUUID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "org_id", &obj.OrgID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_type_id", &obj.CloudTypeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tld_credential_id", &obj.TldCredentialID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_msg", &obj.StatusMsg)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "subset_selected", &obj.SubsetSelected)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_discover_start_time", &obj.LastDiscoverStartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_discover_completed_time", &obj.LastDiscoverCompletedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_successful_discover_start_time", &obj.LastSuccessfulDiscoverStartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_successful_discover_completed_time", &obj.LastSuccessfulDiscoverCompletedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_type", &obj.TaskType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tasks", &obj.Tasks, UnmarshalScopeDetailsGatewayTask)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_updated_time", &obj.StatusUpdatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "collectors_by_type", &obj.CollectorsByType, UnmarshalCollector)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "credentials_by_type", &obj.CredentialsByType, UnmarshalScopeDetailsCredential)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "credentials_by_sub_categeory_type", &obj.CredentialsBySubCategeoryType, UnmarshalScopeDetailsCredential)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sub_categories_by_type", &obj.SubCategoriesByType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_groups", &obj.ResourceGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "region_names", &obj.RegionNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_type", &obj.CloudType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "env_sub_category", &obj.EnvSubCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tld_credentail", &obj.TldCredentail, UnmarshalScopeDetailsCredential)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "collectors", &obj.Collectors, UnmarshalCollector)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first_level_scoped_data", &obj.FirstLevelScopedData, UnmarshalScopeDetailsAssetData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "discovery_methods", &obj.DiscoveryMethods)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "discovery_method", &obj.DiscoveryMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_type", &obj.FileType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_format", &obj.FileFormat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_by", &obj.ModifiedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_discovery_scheduled", &obj.IsDiscoveryScheduled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "discovery_setting_id", &obj.DiscoverySettingID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "include_new_eagerly", &obj.IncludeNewEagerly)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "correlation_id", &obj.CorrelationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "credential_attributes", &obj.CredentialAttributes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeCollector : The collector details of the scope.
type ScopeCollector struct {
	// The collector IDs of the scope.
	CollectorIds []string `json:"collector_ids" validate:"required"`
}

// NewScopeCollector : Instantiate ScopeCollector (Generic Model Constructor)
func (*PostureManagementV2) NewScopeCollector(collectorIds []string) (_model *ScopeCollector, err error) {
	_model = &ScopeCollector{
		CollectorIds: collectorIds,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalScopeCollector unmarshals an instance of ScopeCollector from the specified map of raw messages.
func UnmarshalScopeCollector(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeCollector)
	err = core.UnmarshalPrimitive(m, "collector_ids", &obj.CollectorIds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeCredential : The details of the scope credential.
type ScopeCredential struct {
	// The credential attribute of the scope.
	CredentialAttribute *string `json:"credential_attribute,omitempty"`

	// The credential ID of the scope.
	CredentialID *string `json:"credential_id" validate:"required"`
}

// NewScopeCredential : Instantiate ScopeCredential (Generic Model Constructor)
func (*PostureManagementV2) NewScopeCredential(credentialID string) (_model *ScopeCredential, err error) {
	_model = &ScopeCredential{
		CredentialID: core.StringPtr(credentialID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalScopeCredential unmarshals an instance of ScopeCredential from the specified map of raw messages.
func UnmarshalScopeCredential(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeCredential)
	err = core.UnmarshalPrimitive(m, "credential_attribute", &obj.CredentialAttribute)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "credential_id", &obj.CredentialID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeDetailsAssetData : The asset data details of the scope.
type ScopeDetailsAssetData struct {
	// The asset object of the scope.
	ScopeObject *string `json:"scope_object,omitempty"`

	// The initial value of the scope.
	ScopeInitScope *string `json:"scope_init_scope,omitempty"`

	// The asset of the scope.
	Scope *string `json:"scope,omitempty"`

	// The changed value of the scope.
	ScopeChanged *bool `json:"scope_changed,omitempty"`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The properties of the scope.
	ScopeProperties *string `json:"scope_properties,omitempty"`

	// The overlay of the scope.
	ScopeOverlay *string `json:"scope_overlay,omitempty"`

	// The newfound value of the scope.
	ScopeNewFound *bool `json:"scope_new_found,omitempty"`

	// The discovery status of the scope.
	ScopeDiscoveryStatus interface{} `json:"scope_discovery_status,omitempty"`

	// The fact status of the scope.
	ScopeFactStatus interface{} `json:"scope_fact_status,omitempty"`

	// The facts of the scope.
	ScopeFacts *string `json:"scope_facts,omitempty"`

	// The list members of the scope.
	ScopeListMembers interface{} `json:"scope_list_members,omitempty"`

	// The children of the scope.
	ScopeChildren interface{} `json:"scope_children,omitempty"`

	// The resource category of the scope.
	ScopeResourceCategory *string `json:"scope_resource_category,omitempty"`

	// The resource type of the scope.
	ScopeResourceType *string `json:"scope_resource_type,omitempty"`

	// The resource of the scope.
	ScopeResource *string `json:"scope_resource,omitempty"`

	// The resource attributes of the scope.
	ScopeResourceAttributes interface{} `json:"scope_resource_attributes,omitempty"`

	// The drift of the scope.
	ScopeDrift *string `json:"scope_drift,omitempty"`

	// The parse status of the scope.
	ScopeParseStatus *string `json:"scope_parse_status,omitempty"`

	// The transformed facts of the scope.
	ScopeTransformedFacts interface{} `json:"scope_transformed_facts,omitempty"`

	// The collector ID of the scope.
	ScopeCollectorID *int64 `json:"scope_collector_id,omitempty"`
}

// UnmarshalScopeDetailsAssetData unmarshals an instance of ScopeDetailsAssetData from the specified map of raw messages.
func UnmarshalScopeDetailsAssetData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeDetailsAssetData)
	err = core.UnmarshalPrimitive(m, "scope_object", &obj.ScopeObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_init_scope", &obj.ScopeInitScope)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope", &obj.Scope)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_changed", &obj.ScopeChanged)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_id", &obj.ScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_properties", &obj.ScopeProperties)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_overlay", &obj.ScopeOverlay)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_new_found", &obj.ScopeNewFound)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_discovery_status", &obj.ScopeDiscoveryStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_fact_status", &obj.ScopeFactStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_facts", &obj.ScopeFacts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_list_members", &obj.ScopeListMembers)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_children", &obj.ScopeChildren)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_resource_category", &obj.ScopeResourceCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_resource_type", &obj.ScopeResourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_resource", &obj.ScopeResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_resource_attributes", &obj.ScopeResourceAttributes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_drift", &obj.ScopeDrift)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_parse_status", &obj.ScopeParseStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_transformed_facts", &obj.ScopeTransformedFacts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_collector_id", &obj.ScopeCollectorID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeDetailsCredential : The credential details of the scope.
type ScopeDetailsCredential struct {
	// The credential ID of the scope.
	ID *string `json:"id,omitempty"`

	// The credential name of the scope.
	Name *string `json:"name,omitempty"`

	// The credential uuid of the scope.
	UUID *string `json:"uuid,omitempty"`

	// The credential type of the scope.
	Type *string `json:"type,omitempty"`

	// The credential data of the scope.
	Data interface{} `json:"data,omitempty"`

	// The display fields of the credential. The fields change based on the selected credential type.
	DisplayFields *ScopeDetailsCredentialDisplayFields `json:"display_fields,omitempty"`

	// The credential version timestamp of the scope.
	VersionTimestamp interface{} `json:"version_timestamp,omitempty"`

	// The credential description of the scope.
	Description *string `json:"description,omitempty"`

	// The configuration of whether the scope's credential is enabled.
	IsEnabled *bool `json:"is_enabled,omitempty"`

	// The credential gateway key of the scope.
	GatewayKey *string `json:"gateway_key,omitempty"`

	// The credential purpose of the scope.
	Purpose *string `json:"purpose,omitempty"`
}

// UnmarshalScopeDetailsCredential unmarshals an instance of ScopeDetailsCredential from the specified map of raw messages.
func UnmarshalScopeDetailsCredential(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeDetailsCredential)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uuid", &obj.UUID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "display_fields", &obj.DisplayFields, UnmarshalScopeDetailsCredentialDisplayFields)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_timestamp", &obj.VersionTimestamp)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_enabled", &obj.IsEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "gateway_key", &obj.GatewayKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "purpose", &obj.Purpose)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeDetailsCredentialDisplayFields : The display fields of the credential. The fields change based on the selected credential type.
type ScopeDetailsCredentialDisplayFields struct {
	// The IBM Cloud API Key. The API key is mandatory for the IBM Credential type.
	IBMAPIKey *string `json:"ibm_api_key,omitempty"`

	// The Amazon Web Services client ID. The client ID is mandatory when AWS Cloud is selected as the credential type.
	AwsClientID *string `json:"aws_client_id,omitempty"`

	// The Amazon Web Services client secret. The client secret is mandatory when AWS Cloud is selected as the credential
	// type.
	AwsClientSecret *string `json:"aws_client_secret,omitempty"`

	// The Amazon Web Services region. The region is used when AWS Cloud is selected as the credential type.
	AwsRegion *string `json:"aws_region,omitempty"`

	// The Amazon Web Services arn value. The arn value is used when AWS Cloud is selected as the credential type.
	AwsArn *string `json:"aws_arn,omitempty"`

	// The username of the user. The username is mandatory when the credential type is DataBase, Kerberos, OpenStack, and
	// Username-Password.
	Username *string `json:"username,omitempty"`

	// The password of the user. The password is mandatory when the credential type is DataBase, Kerberos, OpenStack, and
	// Username-Password.
	Password *string `json:"password,omitempty"`

	// The Microsoft Azure client ID. The client ID is mandatory when Azure is selected as the credential type.
	AzureClientID *string `json:"azure_client_id,omitempty"`

	// The Microsoft Azure client secret. The secret is mandatory when the type of credential is set to Azure.
	AzureClientSecret *string `json:"azure_client_secret,omitempty"`

	// The Microsoft Azure subscription ID. The subscription ID is mandatory when the type of credential is set to Azure.
	AzureSubscriptionID *string `json:"azure_subscription_id,omitempty"`

	// The Microsoft Azure resource group. The resource group is used when Azure is the credential type.
	AzureResourceGroup *string `json:"azure_resource_group,omitempty"`

	// The database name. The database name is mandatory when Database is the credential type.
	DatabaseName *string `json:"database_name,omitempty"`

	// The Kerberos Windows authentication type. The authentication type is mandatory when the credential type is Kerberos
	// Windows.
	WinrmAuthtype *string `json:"winrm_authtype,omitempty"`

	// The Kerberos Windows SSL. The SSL is mandatory when the credential type is Kerberos Windows.
	WinrmUsessl *string `json:"winrm_usessl,omitempty"`

	// The Kerberos Windows port. The port is mandatory When Kerberos Windows is the credential type.
	WinrmPort *string `json:"winrm_port,omitempty"`

	// The Microsoft 365 client ID. The client ID is mandatory when Microsoft 365 is the credential type.
	Ms365ClientID *string `json:"ms_365_client_id,omitempty"`

	// The Microsoft 365 client secret. The secret is mandatory when Microsoft 365 is the credential type.
	Ms365ClientSecret *string `json:"ms_365_client_secret,omitempty"`

	// The Microsoft 365 tenant ID. The tenant ID is mandatory when Microsoft 365 is the credential type.
	Ms365TenantID *string `json:"ms_365_tenant_id,omitempty"`

	// The auth url of the OpenStack cloud. The auth url is mandatory when OpenStack is the credential type.
	AuthURL *string `json:"auth_url,omitempty"`

	// The project name of the OpenStack cloud. The project name is mandatory when OpenStack is the credential type.
	ProjectName *string `json:"project_name,omitempty"`

	// The user domain name of the OpenStack cloud. The domain name is mandatory when OpenStack is the credential type.
	UserDomainName *string `json:"user_domain_name,omitempty"`

	// The project domain name of the OpenStack cloud. The project domain name is mandatory when OpenStack is the
	// credential type.
	ProjectDomainName *string `json:"project_domain_name,omitempty"`
}

// UnmarshalScopeDetailsCredentialDisplayFields unmarshals an instance of ScopeDetailsCredentialDisplayFields from the specified map of raw messages.
func UnmarshalScopeDetailsCredentialDisplayFields(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeDetailsCredentialDisplayFields)
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IBMAPIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_client_id", &obj.AwsClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_client_secret", &obj.AwsClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_arn", &obj.AwsArn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_client_id", &obj.AzureClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_client_secret", &obj.AzureClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "database_name", &obj.DatabaseName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_authtype", &obj.WinrmAuthtype)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_usessl", &obj.WinrmUsessl)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_port", &obj.WinrmPort)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_client_id", &obj.Ms365ClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_client_secret", &obj.Ms365ClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_tenant_id", &obj.Ms365TenantID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_url", &obj.AuthURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "project_name", &obj.ProjectName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_domain_name", &obj.UserDomainName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "project_domain_name", &obj.ProjectDomainName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeDetailsGatewayTask : The gateway task details of the scope.
type ScopeDetailsGatewayTask struct {
	// The task logs of the scope.
	TaskLogs []TaskLogs `json:"task_logs,omitempty"`

	// The task ID of the scope.
	TaskID *int64 `json:"task_id,omitempty"`

	// The task gateway ID of the scope.
	TaskGatewayID *int64 `json:"task_gateway_id,omitempty"`

	// The task gateway name of the scope.
	TaskGatewayName *string `json:"task_gateway_name,omitempty"`

	// The task type of the scope.
	TaskTaskType *string `json:"task_task_type,omitempty"`

	// The task gateway schema ID of the scope.
	TaskGatewaySchemaID *int64 `json:"task_gateway_schema_id,omitempty"`

	// The task schema name of the scope.
	TaskSchemaName *string `json:"task_schema_name,omitempty"`

	// The task-discover ID of the scope.
	TaskDiscoverID *int64 `json:"task_discover_id,omitempty"`

	// The task status of the scope.
	TaskStatus *string `json:"task_status,omitempty"`

	// The task status message of the scope.
	TaskStatusMsg *string `json:"task_status_msg,omitempty"`

	// The task start time of the scope.
	TaskStartTime *int64 `json:"task_start_time,omitempty"`

	// The task updated time of the scope.
	TaskUpdatedTime *int64 `json:"task_updated_time,omitempty"`

	// The task derived status of the scope.
	TaskDerivedStatus *string `json:"task_derived_status,omitempty"`

	// The user who created the task that is associated with scope.
	TaskCreatedBy *string `json:"task_created_by,omitempty"`
}

// Constants associated with the ScopeDetailsGatewayTask.TaskTaskType property.
// The task type of the scope.
const (
	ScopeDetailsGatewayTaskTaskTaskTypeAborttasksConst            = "aborttasks"
	ScopeDetailsGatewayTaskTaskTaskTypeCertRegularValidationConst = "cert_regular_validation"
	ScopeDetailsGatewayTaskTaskTaskTypeCertValidationConst        = "cert_validation"
	ScopeDetailsGatewayTaskTaskTaskTypeCveRegularValidationConst  = "cve_regular_validation"
	ScopeDetailsGatewayTaskTaskTaskTypeCveValidationConst         = "cve_validation"
	ScopeDetailsGatewayTaskTaskTaskTypeDiscoverConst              = "discover"
	ScopeDetailsGatewayTaskTaskTaskTypeEolRegularValidationConst  = "eol_regular_validation"
	ScopeDetailsGatewayTaskTaskTaskTypeEolValidationConst         = "eol_validation"
	ScopeDetailsGatewayTaskTaskTaskTypeEvidenceConst              = "evidence"
	ScopeDetailsGatewayTaskTaskTaskTypeFactcollectionConst        = "factcollection"
	ScopeDetailsGatewayTaskTaskTaskTypeFactvalidationConst        = "factvalidation"
	ScopeDetailsGatewayTaskTaskTaskTypeInventoryConst             = "inventory"
	ScopeDetailsGatewayTaskTaskTaskTypeNopConst                   = "nop"
	ScopeDetailsGatewayTaskTaskTaskTypeRemediationConst           = "remediation"
	ScopeDetailsGatewayTaskTaskTaskTypeScriptConst                = "script"
	ScopeDetailsGatewayTaskTaskTaskTypeSubsetvalidateConst        = "subsetvalidate"
	ScopeDetailsGatewayTaskTaskTaskTypeTldiscoverConst            = "tldiscover"
)

// Constants associated with the ScopeDetailsGatewayTask.TaskStatus property.
// The task status of the scope.
const (
	ScopeDetailsGatewayTaskTaskStatusAbortTaskRequestCompletedConst       = "abort_task_request_completed"
	ScopeDetailsGatewayTaskTaskStatusAbortTaskRequestFailedConst          = "abort_task_request_failed"
	ScopeDetailsGatewayTaskTaskStatusAbortTaskRequestReceivedConst        = "abort_task_request_received"
	ScopeDetailsGatewayTaskTaskStatusCertRegularValidationCompletedConst  = "cert_regular_validation_completed"
	ScopeDetailsGatewayTaskTaskStatusCertRegularValidationErrorConst      = "cert_regular_validation_error"
	ScopeDetailsGatewayTaskTaskStatusCertRegularValidationStartedConst    = "cert_regular_validation_started"
	ScopeDetailsGatewayTaskTaskStatusCertValidationCompletedConst         = "cert_validation_completed"
	ScopeDetailsGatewayTaskTaskStatusCertValidationErrorConst             = "cert_validation_error"
	ScopeDetailsGatewayTaskTaskStatusCertValidationStartedConst           = "cert_validation_started"
	ScopeDetailsGatewayTaskTaskStatusControllerAbortedConst               = "controller_aborted"
	ScopeDetailsGatewayTaskTaskStatusCveRegularValidationCompletedConst   = "cve_regular_validation_completed"
	ScopeDetailsGatewayTaskTaskStatusCveRegularValidationErrorConst       = "cve_regular_validation_error"
	ScopeDetailsGatewayTaskTaskStatusCveRegularValidationStartedConst     = "cve_regular_validation_started"
	ScopeDetailsGatewayTaskTaskStatusCveValidationCompletedConst          = "cve_validation_completed"
	ScopeDetailsGatewayTaskTaskStatusCveValidationErrorConst              = "cve_validation_error"
	ScopeDetailsGatewayTaskTaskStatusCveValidationStartedConst            = "cve_validation_started"
	ScopeDetailsGatewayTaskTaskStatusDiscoveryCompletedConst              = "discovery_completed"
	ScopeDetailsGatewayTaskTaskStatusDiscoveryInProgressConst             = "discovery_in_progress"
	ScopeDetailsGatewayTaskTaskStatusDiscoveryResultPostedNoErrorConst    = "discovery_result_posted_no_error"
	ScopeDetailsGatewayTaskTaskStatusDiscoveryResultPostedWithErrorConst  = "discovery_result_posted_with_error"
	ScopeDetailsGatewayTaskTaskStatusDiscoveryStartedConst                = "discovery_started"
	ScopeDetailsGatewayTaskTaskStatusEolRegularValidationCompletedConst   = "eol_regular_validation_completed"
	ScopeDetailsGatewayTaskTaskStatusEolRegularValidationErrorConst       = "eol_regular_validation_error"
	ScopeDetailsGatewayTaskTaskStatusEolRegularValidationStartedConst     = "eol_regular_validation_started"
	ScopeDetailsGatewayTaskTaskStatusEolValidationCompletedConst          = "eol_validation_completed"
	ScopeDetailsGatewayTaskTaskStatusEolValidationErrorConst              = "eol_validation_error"
	ScopeDetailsGatewayTaskTaskStatusEolValidationStartedConst            = "eol_validation_started"
	ScopeDetailsGatewayTaskTaskStatusErrorInAbortTaskRequestConst         = "error_in_abort_task_request"
	ScopeDetailsGatewayTaskTaskStatusErrorInDiscoverConst                 = "error_in_discover"
	ScopeDetailsGatewayTaskTaskStatusErrorInFactCollectionConst           = "error_in_fact_collection"
	ScopeDetailsGatewayTaskTaskStatusErrorInFactValidationConst           = "error_in_fact_validation"
	ScopeDetailsGatewayTaskTaskStatusErrorInInventoryConst                = "error_in_inventory"
	ScopeDetailsGatewayTaskTaskStatusErrorInRemediationConst              = "error_in_remediation"
	ScopeDetailsGatewayTaskTaskStatusErrorInValidationConst               = "error_in_validation"
	ScopeDetailsGatewayTaskTaskStatusFactCollectionCompletedConst         = "fact_collection_completed"
	ScopeDetailsGatewayTaskTaskStatusFactCollectionInProgressConst        = "fact_collection_in_progress"
	ScopeDetailsGatewayTaskTaskStatusFactCollectionStartedConst           = "fact_collection_started"
	ScopeDetailsGatewayTaskTaskStatusFactValidationCompletedConst         = "fact_validation_completed"
	ScopeDetailsGatewayTaskTaskStatusFactValidationInProgressConst        = "fact_validation_in_progress"
	ScopeDetailsGatewayTaskTaskStatusFactValidationStartedConst           = "fact_validation_started"
	ScopeDetailsGatewayTaskTaskStatusGatewayAbortedConst                  = "gateway_aborted"
	ScopeDetailsGatewayTaskTaskStatusInventoryCompletedConst              = "inventory_completed"
	ScopeDetailsGatewayTaskTaskStatusInventoryCompletedWithErrorConst     = "inventory_completed_with_error"
	ScopeDetailsGatewayTaskTaskStatusInventoryInProgressConst             = "inventory_in_progress"
	ScopeDetailsGatewayTaskTaskStatusInventoryStartedConst                = "inventory_started"
	ScopeDetailsGatewayTaskTaskStatusLocationChangeAbortedConst           = "location_change_aborted"
	ScopeDetailsGatewayTaskTaskStatusNotAcceptedConst                     = "not_accepted"
	ScopeDetailsGatewayTaskTaskStatusPendingConst                         = "pending"
	ScopeDetailsGatewayTaskTaskStatusRemediationCompletedConst            = "remediation_completed"
	ScopeDetailsGatewayTaskTaskStatusRemediationInProgressConst           = "remediation_in_progress"
	ScopeDetailsGatewayTaskTaskStatusRemediationStartedConst              = "remediation_started"
	ScopeDetailsGatewayTaskTaskStatusSentToCollectorConst                 = "sent_to_collector"
	ScopeDetailsGatewayTaskTaskStatusUserAbortedConst                     = "user_aborted"
	ScopeDetailsGatewayTaskTaskStatusValidationCompletedConst             = "validation_completed"
	ScopeDetailsGatewayTaskTaskStatusValidationInProgressConst            = "validation_in_progress"
	ScopeDetailsGatewayTaskTaskStatusValidationResultPostedNoErrorConst   = "validation_result_posted_no_error"
	ScopeDetailsGatewayTaskTaskStatusValidationResultPostedWithErrorConst = "validation_result_posted_with_error"
	ScopeDetailsGatewayTaskTaskStatusValidationStartedConst               = "validation_started"
	ScopeDetailsGatewayTaskTaskStatusWaitingForRefineConst                = "waiting_for_refine"
)

// Constants associated with the ScopeDetailsGatewayTask.TaskDerivedStatus property.
// The task derived status of the scope.
const (
	ScopeDetailsGatewayTaskTaskDerivedStatusAbortTaskRequestCompletedConst       = "abort_task_request_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusAbortTaskRequestFailedConst          = "abort_task_request_failed"
	ScopeDetailsGatewayTaskTaskDerivedStatusAbortTaskRequestReceivedConst        = "abort_task_request_received"
	ScopeDetailsGatewayTaskTaskDerivedStatusCertRegularValidationCompletedConst  = "cert_regular_validation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusCertRegularValidationErrorConst      = "cert_regular_validation_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusCertRegularValidationStartedConst    = "cert_regular_validation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusCertValidationCompletedConst         = "cert_validation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusCertValidationErrorConst             = "cert_validation_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusCertValidationStartedConst           = "cert_validation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusControllerAbortedConst               = "controller_aborted"
	ScopeDetailsGatewayTaskTaskDerivedStatusCveRegularValidationCompletedConst   = "cve_regular_validation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusCveRegularValidationErrorConst       = "cve_regular_validation_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusCveRegularValidationStartedConst     = "cve_regular_validation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusCveValidationCompletedConst          = "cve_validation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusCveValidationErrorConst              = "cve_validation_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusCveValidationStartedConst            = "cve_validation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusDiscoveryCompletedConst              = "discovery_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusDiscoveryInProgressConst             = "discovery_in_progress"
	ScopeDetailsGatewayTaskTaskDerivedStatusDiscoveryResultPostedNoErrorConst    = "discovery_result_posted_no_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusDiscoveryResultPostedWithErrorConst  = "discovery_result_posted_with_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusDiscoveryStartedConst                = "discovery_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusEolRegularValidationCompletedConst   = "eol_regular_validation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusEolRegularValidationErrorConst       = "eol_regular_validation_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusEolRegularValidationStartedConst     = "eol_regular_validation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusEolValidationCompletedConst          = "eol_validation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusEolValidationErrorConst              = "eol_validation_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusEolValidationStartedConst            = "eol_validation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusErrorInAbortTaskRequestConst         = "error_in_abort_task_request"
	ScopeDetailsGatewayTaskTaskDerivedStatusErrorInDiscoverConst                 = "error_in_discover"
	ScopeDetailsGatewayTaskTaskDerivedStatusErrorInFactCollectionConst           = "error_in_fact_collection"
	ScopeDetailsGatewayTaskTaskDerivedStatusErrorInFactValidationConst           = "error_in_fact_validation"
	ScopeDetailsGatewayTaskTaskDerivedStatusErrorInInventoryConst                = "error_in_inventory"
	ScopeDetailsGatewayTaskTaskDerivedStatusErrorInRemediationConst              = "error_in_remediation"
	ScopeDetailsGatewayTaskTaskDerivedStatusErrorInValidationConst               = "error_in_validation"
	ScopeDetailsGatewayTaskTaskDerivedStatusFactCollectionCompletedConst         = "fact_collection_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusFactCollectionInProgressConst        = "fact_collection_in_progress"
	ScopeDetailsGatewayTaskTaskDerivedStatusFactCollectionStartedConst           = "fact_collection_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusFactValidationCompletedConst         = "fact_validation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusFactValidationInProgressConst        = "fact_validation_in_progress"
	ScopeDetailsGatewayTaskTaskDerivedStatusFactValidationStartedConst           = "fact_validation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusGatewayAbortedConst                  = "gateway_aborted"
	ScopeDetailsGatewayTaskTaskDerivedStatusInventoryCompletedConst              = "inventory_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusInventoryCompletedWithErrorConst     = "inventory_completed_with_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusInventoryInProgressConst             = "inventory_in_progress"
	ScopeDetailsGatewayTaskTaskDerivedStatusInventoryStartedConst                = "inventory_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusLocationChangeAbortedConst           = "location_change_aborted"
	ScopeDetailsGatewayTaskTaskDerivedStatusNotAcceptedConst                     = "not_accepted"
	ScopeDetailsGatewayTaskTaskDerivedStatusPendingConst                         = "pending"
	ScopeDetailsGatewayTaskTaskDerivedStatusRemediationCompletedConst            = "remediation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusRemediationInProgressConst           = "remediation_in_progress"
	ScopeDetailsGatewayTaskTaskDerivedStatusRemediationStartedConst              = "remediation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusSentToCollectorConst                 = "sent_to_collector"
	ScopeDetailsGatewayTaskTaskDerivedStatusUserAbortedConst                     = "user_aborted"
	ScopeDetailsGatewayTaskTaskDerivedStatusValidationCompletedConst             = "validation_completed"
	ScopeDetailsGatewayTaskTaskDerivedStatusValidationInProgressConst            = "validation_in_progress"
	ScopeDetailsGatewayTaskTaskDerivedStatusValidationResultPostedNoErrorConst   = "validation_result_posted_no_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusValidationResultPostedWithErrorConst = "validation_result_posted_with_error"
	ScopeDetailsGatewayTaskTaskDerivedStatusValidationStartedConst               = "validation_started"
	ScopeDetailsGatewayTaskTaskDerivedStatusWaitingForRefineConst                = "waiting_for_refine"
)

// UnmarshalScopeDetailsGatewayTask unmarshals an instance of ScopeDetailsGatewayTask from the specified map of raw messages.
func UnmarshalScopeDetailsGatewayTask(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeDetailsGatewayTask)
	err = core.UnmarshalModel(m, "task_logs", &obj.TaskLogs, UnmarshalTaskLogs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_id", &obj.TaskID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_gateway_id", &obj.TaskGatewayID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_gateway_name", &obj.TaskGatewayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_task_type", &obj.TaskTaskType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_gateway_schema_id", &obj.TaskGatewaySchemaID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_schema_name", &obj.TaskSchemaName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_discover_id", &obj.TaskDiscoverID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_status", &obj.TaskStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_status_msg", &obj.TaskStatusMsg)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_start_time", &obj.TaskStartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_updated_time", &obj.TaskUpdatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_derived_status", &obj.TaskDerivedStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_created_by", &obj.TaskCreatedBy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeItem : Scope.
type ScopeItem struct {
	// A detailed description of the scope.
	Description *string `json:"description" validate:"required"`

	// The user who created the scope.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The user who most recently modified the scope.
	ModifiedBy *string `json:"modified_by" validate:"required"`

	// An auto-generated unique identifier for the scope.
	ID *string `json:"id" validate:"required"`

	// The unique user ID of the scope.
	UUID *string `json:"uuid" validate:"required"`

	// A unique name for your scope.
	Name *string `json:"name" validate:"required"`

	// The scope is enabled or disabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// The environment that the scope is targeted to.
	CredentialType *string `json:"credential_type" validate:"required"`

	// The time that the scope was created in UTC.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The time that the scope was last modified in UTC.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The collectors of the scope. The collectors are displayed only when the values exist.
	Collectors []Collector `json:"collectors" validate:"required"`
}

// Constants associated with the ScopeItem.CredentialType property.
// The environment that the scope is targeted to.
const (
	ScopeItemCredentialTypeAwsConst       = "aws"
	ScopeItemCredentialTypeAzureConst     = "azure"
	ScopeItemCredentialTypeGcpConst       = "gcp"
	ScopeItemCredentialTypeHostedConst    = "hosted"
	ScopeItemCredentialTypeIBMConst       = "ibm"
	ScopeItemCredentialTypeOnPremiseConst = "on_premise"
	ScopeItemCredentialTypeOpenstackConst = "openstack"
	ScopeItemCredentialTypeServicesConst  = "services"
)

// UnmarshalScopeItem unmarshals an instance of ScopeItem from the specified map of raw messages.
func UnmarshalScopeItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeItem)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_by", &obj.ModifiedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uuid", &obj.UUID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "credential_type", &obj.CredentialType)
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
	err = core.UnmarshalModel(m, "collectors", &obj.Collectors, UnmarshalCollector)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeList : Scopes list.
type ScopeList struct {
	// The offset of the page.
	Offset *int64 `json:"offset" validate:"required"`

	// The number of scopes that are displayed per page.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of scopes. The value is 0 if no scopes are available. The detail fields are not available in that
	// case.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The URL of a page.
	First *PageLink `json:"first" validate:"required"`

	// The URL of a page.
	Last *PageLink `json:"last" validate:"required"`

	// The URL of a page.
	Previous *PageLink `json:"previous,omitempty"`

	// The URL of a page.
	Next *PageLink `json:"next,omitempty"`

	// Scopes.
	Scopes []ScopeItem `json:"scopes" validate:"required"`
}

// UnmarshalScopeList unmarshals an instance of ScopeList from the specified map of raw messages.
func UnmarshalScopeList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scopes", &obj.Scopes, UnmarshalScopeItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeTaskStatus : Get the current task list for the collectors that are attached to the scope.
type ScopeTaskStatus struct {
	// The correlation ID.
	CorrelationID *string `json:"correlation_id" validate:"required"`

	// The status of a task.
	Status *string `json:"status" validate:"required"`

	// The time that the task started.
	StartTime *string `json:"start_time" validate:"required"`

	// The time that the scope was last updated. This value exists when a collector is installed and running.
	LastHeartbeat *strfmt.DateTime `json:"last_heartbeat" validate:"required"`
}

// UnmarshalScopeTaskStatus unmarshals an instance of ScopeTaskStatus from the specified map of raw messages.
func UnmarshalScopeTaskStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeTaskStatus)
	err = core.UnmarshalPrimitive(m, "correlation_id", &obj.CorrelationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_heartbeat", &obj.LastHeartbeat)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Summary : A list of scans summary.
type Summary struct {
	// The scan ID.
	ID *string `json:"id" validate:"required"`

	// The scan discovery ID.
	DiscoverID *string `json:"discover_id" validate:"required"`

	// The scan profile ID.
	ProfileID *string `json:"profile_id" validate:"required"`

	// The scan profile name.
	ProfileName *string `json:"profile_name" validate:"required"`

	// The scan summary scope ID.
	ScopeID *string `json:"scope_id" validate:"required"`

	// The list of controls that are on the scan summary.
	Controls []Control `json:"controls" validate:"required"`
}

// UnmarshalSummary unmarshals an instance of Summary from the specified map of raw messages.
func UnmarshalSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Summary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "discover_id", &obj.DiscoverID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_id", &obj.ProfileID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_name", &obj.ProfileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_id", &obj.ScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalControl)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SummaryItem : The result of a scan summary.
type SummaryItem struct {
	// The ID of the scan.
	ID *string `json:"id" validate:"required"`

	// A system-generated name that is the combination of 12 characters in the scope name and 12 characters in a profile
	// name.
	Name *string `json:"name" validate:"required"`

	// The ID of the scope.
	ScopeID *string `json:"scope_id" validate:"required"`

	// The name of the scope.
	ScopeName *string `json:"scope_name" validate:"required"`

	// The entity that ran the report.
	ReportRunBy *string `json:"report_run_by" validate:"required"`

	// The date and time that the scan was run.
	StartTime *strfmt.DateTime `json:"start_time" validate:"required"`

	// The date and time that the scan completed.
	EndTime *strfmt.DateTime `json:"end_time" validate:"required"`

	// The status of the collector as it completes a scan.
	Status *string `json:"status" validate:"required"`

	// The list of profiles.
	Profiles []ProfileResult `json:"profiles" validate:"required"`

	// The list of group profiles.
	GroupProfiles []ProfileResult `json:"group_profiles" validate:"required"`
}

// Constants associated with the SummaryItem.Status property.
// The status of the collector as it completes a scan.
const (
	SummaryItemStatusAbortTaskRequestCompletedConst       = "abort_task_request_completed"
	SummaryItemStatusAbortTaskRequestFailedConst          = "abort_task_request_failed"
	SummaryItemStatusAbortTaskRequestReceivedConst        = "abort_task_request_received"
	SummaryItemStatusControllerAbortedConst               = "controller_aborted"
	SummaryItemStatusDiscoveryCompletedConst              = "discovery_completed"
	SummaryItemStatusDiscoveryInProgressConst             = "discovery_in_progress"
	SummaryItemStatusDiscoveryResultPostedNoErrorConst    = "discovery_result_posted_no_error"
	SummaryItemStatusDiscoveryResultPostedWithErrorConst  = "discovery_result_posted_with_error"
	SummaryItemStatusDiscoveryStartedConst                = "discovery_started"
	SummaryItemStatusErrorInAbortTaskRequestConst         = "error_in_abort_task_request"
	SummaryItemStatusErrorInDiscoveryConst                = "error_in_discovery"
	SummaryItemStatusErrorInFactCollectionConst           = "error_in_fact_collection"
	SummaryItemStatusErrorInFactValidationConst           = "error_in_fact_validation"
	SummaryItemStatusErrorInInventoryConst                = "error_in_inventory"
	SummaryItemStatusErrorInRemediationConst              = "error_in_remediation"
	SummaryItemStatusErrorInValidationConst               = "error_in_validation"
	SummaryItemStatusFactCollectionCompletedConst         = "fact_collection_completed"
	SummaryItemStatusFactCollectionInProgressConst        = "fact_collection_in_progress"
	SummaryItemStatusFactCollectionStartedConst           = "fact_collection_started"
	SummaryItemStatusFactValidationCompletedConst         = "fact_validation_completed"
	SummaryItemStatusFactValidationInProgressConst        = "fact_validation_in_progress"
	SummaryItemStatusFactValidationStartedConst           = "fact_validation_started"
	SummaryItemStatusGatewayAbortedConst                  = "gateway_aborted"
	SummaryItemStatusInventoryCompletedConst              = "inventory_completed"
	SummaryItemStatusInventoryCompletedWithErrorConst     = "inventory_completed_with_error"
	SummaryItemStatusInventoryInProgressConst             = "inventory_in_progress"
	SummaryItemStatusInventoryStartedConst                = "inventory_started"
	SummaryItemStatusNotAcceptedConst                     = "not_accepted"
	SummaryItemStatusPendingConst                         = "pending"
	SummaryItemStatusRemediationCompletedConst            = "remediation_completed"
	SummaryItemStatusRemediationInProgressConst           = "remediation_in_progress"
	SummaryItemStatusRemediationStartedConst              = "remediation_started"
	SummaryItemStatusSentToCollectorConst                 = "sent_to_collector"
	SummaryItemStatusUserAbortedConst                     = "user_aborted"
	SummaryItemStatusValidationCompletedConst             = "validation_completed"
	SummaryItemStatusValidationInProgressConst            = "validation_in_progress"
	SummaryItemStatusValidationResultPostedNoErrorConst   = "validation_result_posted_no_error"
	SummaryItemStatusValidationResultPostedWithErrorConst = "validation_result_posted_with_error"
	SummaryItemStatusValidationStartedConst               = "validation_started"
	SummaryItemStatusWaitingForRefineConst                = "waiting_for_refine"
)

// UnmarshalSummaryItem unmarshals an instance of SummaryItem from the specified map of raw messages.
func UnmarshalSummaryItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SummaryItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_id", &obj.ScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_name", &obj.ScopeName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_run_by", &obj.ReportRunBy)
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
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalProfileResult)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "group_profiles", &obj.GroupProfiles, UnmarshalProfileResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SummaryList : A list of scan summaries.
type SummaryList struct {
	// The offset of the page.
	Offset *int64 `json:"offset" validate:"required"`

	// The number of scans that are displayed per page.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of scans that are available in the list of summaries.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The URL of a page.
	First *PageLink `json:"first" validate:"required"`

	// The URL of a page.
	Last *PageLink `json:"last" validate:"required"`

	// The URL of a page.
	Previous *PageLink `json:"previous,omitempty"`

	// The URL of a page.
	Next *PageLink `json:"next,omitempty"`

	// Summaries.
	Summaries []SummaryItem `json:"summaries" validate:"required"`
}

// UnmarshalSummaryList unmarshals an instance of SummaryList from the specified map of raw messages.
func UnmarshalSummaryList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SummaryList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "summaries", &obj.Summaries, UnmarshalSummaryItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SummaryList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil || offset == nil {
		return nil, err
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// TaskLogs : The logs for the tasks that you ran.
type TaskLogs struct {
}

// UnmarshalTaskLogs unmarshals an instance of TaskLogs from the specified map of raw messages.
func UnmarshalTaskLogs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TaskLogs)
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateCollectorOptions : The UpdateCollector options.
type UpdateCollectorOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// JSON Merge-Patch content for update_collector.
	Collector map[string]interface{} `json:"collector" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateCollectorOptions : Instantiate UpdateCollectorOptions
func (*PostureManagementV2) NewUpdateCollectorOptions(id string, collector map[string]interface{}) *UpdateCollectorOptions {
	return &UpdateCollectorOptions{
		ID:        core.StringPtr(id),
		Collector: collector,
	}
}

// SetID : Allow user to set ID
func (_options *UpdateCollectorOptions) SetID(id string) *UpdateCollectorOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetCollector : Allow user to set Collector
func (_options *UpdateCollectorOptions) SetCollector(collector map[string]interface{}) *UpdateCollectorOptions {
	_options.Collector = collector
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateCollectorOptions) SetAccountID(accountID string) *UpdateCollectorOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateCollectorOptions) SetTransactionID(transactionID string) *UpdateCollectorOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCollectorOptions) SetHeaders(param map[string]string) *UpdateCollectorOptions {
	options.Headers = param
	return options
}

// UpdateCredentialDisplayFields : The details of the credential. The details change as the selected credential type varies.
type UpdateCredentialDisplayFields struct {
	// The IBM Cloud API key. The API key is mandatory when IBM is the credential type.
	IBMAPIKey *string `json:"ibm_api_key,omitempty"`

	// The Amazon Web Services client ID. The client ID is mandatory when AWS Cloud is the credential type.
	AwsClientID *string `json:"aws_client_id,omitempty"`

	// The Amazon Web Services client secret. The secret is mandatory when AWS Cloud is the credential type.
	AwsClientSecret *string `json:"aws_client_secret,omitempty"`

	// The Amazon Web Services region. The region is mandatory when AWS Cloud is the credential type.
	AwsRegion *string `json:"aws_region,omitempty"`

	// The Amazon Web Services arn value. The arn value is mandatory when AWS Cloud is the credential type.
	AwsArn *string `json:"aws_arn,omitempty"`

	// The username of the user. The username is mandatory when the credential type is DataBase, Kerberos, or OpenStack.
	Username *string `json:"username,omitempty"`

	// The password of the user. The password is mandatory for DataBase, Kerberos, OpenStack credentials.
	Password *string `json:"password,omitempty"`

	// The Microsoft Azure client ID. The client ID is mandatory when Azure is the credential type.
	AzureClientID *string `json:"azure_client_id,omitempty"`

	// The Microsoft Azure client secret. The secret is mandatory when Azure is the credential type.
	AzureClientSecret *string `json:"azure_client_secret,omitempty"`

	// The Microsoft Azure subscription ID. The subscription ID is mandatory when Azure is the credential type.
	AzureSubscriptionID *string `json:"azure_subscription_id,omitempty"`

	// The Microsoft Azure resource group. The resource group is mandatory when Azure is the credential type.
	AzureResourceGroup *string `json:"azure_resource_group,omitempty"`

	// The Database name. The name is mandatory when Database is the credential type.
	DatabaseName *string `json:"database_name,omitempty"`

	// The Kerberos Windows auth type. The auth type is mandatory when Kerberos is the credential type.
	WinrmAuthtype *string `json:"winrm_authtype,omitempty"`

	// The Kerberos Windows SSL. The SSL is mandatory when Kerberos is the credential type.
	WinrmUsessl *string `json:"winrm_usessl,omitempty"`

	// The Kerberos Windows port. The port is mandatory when Kerberos is the credential type.
	WinrmPort *string `json:"winrm_port,omitempty"`

	// The Microsoft 365 client ID. The client ID is mandatory when Microsoft 365 is the credential type.
	Ms365ClientID *string `json:"ms_365_client_id,omitempty"`

	// The Microsoft 365 client secret. The secret is mandatory when Microsoft 365 is the credential type.
	Ms365ClientSecret *string `json:"ms_365_client_secret,omitempty"`

	// The Microsoft 365 tenant ID. The tenant ID is mandatory when Microsoft 365 is the credential type.
	Ms365TenantID *string `json:"ms_365_tenant_id,omitempty"`

	// The auth url of the OpenStack cloud. The auth url is mandatory when OpenStack is the credential type.
	AuthURL *string `json:"auth_url,omitempty"`

	// The project name of the OpenStack cloud. The project name is mandatory when OpenStack is the credential type.
	ProjectName *string `json:"project_name,omitempty"`

	// The user domain name of the OpenStack cloud. The domain name is mandatory when OpenStack is the credential type.
	UserDomainName *string `json:"user_domain_name,omitempty"`

	// The project domain name of the OpenStack cloud. The project domain name is mandatory when OpenStack is the
	// credential type.
	ProjectDomainName *string `json:"project_domain_name,omitempty"`

	// The user pem file name.
	PemFileName *string `json:"pem_file_name,omitempty"`

	// The base64 encoded form of pem.
	PemData *string `json:"pem_data,omitempty"`
}

// UnmarshalUpdateCredentialDisplayFields unmarshals an instance of UpdateCredentialDisplayFields from the specified map of raw messages.
func UnmarshalUpdateCredentialDisplayFields(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateCredentialDisplayFields)
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IBMAPIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_client_id", &obj.AwsClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_client_secret", &obj.AwsClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_arn", &obj.AwsArn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_client_id", &obj.AzureClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_client_secret", &obj.AzureClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "database_name", &obj.DatabaseName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_authtype", &obj.WinrmAuthtype)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_usessl", &obj.WinrmUsessl)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "winrm_port", &obj.WinrmPort)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_client_id", &obj.Ms365ClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_client_secret", &obj.Ms365ClientSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ms_365_tenant_id", &obj.Ms365TenantID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_url", &obj.AuthURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "project_name", &obj.ProjectName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_domain_name", &obj.UserDomainName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "project_domain_name", &obj.ProjectDomainName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pem_file_name", &obj.PemFileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pem_data", &obj.PemData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateCredentialOptions : The UpdateCredential options.
type UpdateCredentialOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// The status of the credential is enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The credential's type.
	Type *string `json:"type,omitempty"`

	// The credential's name.
	Name *string `json:"name,omitempty"`

	// The credential's description.
	Description *string `json:"description,omitempty"`

	// The details of the credential. The details change as the selected credential type varies.
	DisplayFields *UpdateCredentialDisplayFields `json:"display_fields,omitempty"`

	// The purpose for which the credential is created.
	Purpose *string `json:"purpose,omitempty"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateCredentialOptions.Type property.
// The credential's type.
const (
	UpdateCredentialOptionsTypeAwsCloudConst         = "aws_cloud"
	UpdateCredentialOptionsTypeAzureCloudConst       = "azure_cloud"
	UpdateCredentialOptionsTypeDatabaseConst         = "database"
	UpdateCredentialOptionsTypeIBMCloudConst         = "ibm_cloud"
	UpdateCredentialOptionsTypeKerberosWindowsConst  = "kerberos_windows"
	UpdateCredentialOptionsTypeMs365Const            = "ms_365"
	UpdateCredentialOptionsTypeOpenstackCloudConst   = "openstack_cloud"
	UpdateCredentialOptionsTypeUserNamePemConst      = "user_name_pem"
	UpdateCredentialOptionsTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the UpdateCredentialOptions.Purpose property.
// The purpose for which the credential is created.
const (
	UpdateCredentialOptionsPurposeDiscoveryCollectionConst                = "discovery_collection"
	UpdateCredentialOptionsPurposeDiscoveryCollectionRemediationConst     = "discovery_collection_remediation"
	UpdateCredentialOptionsPurposeDiscoveryFactCollectionConst            = "discovery_fact_collection"
	UpdateCredentialOptionsPurposeDiscoveryFactCollectionRemediationConst = "discovery_fact_collection_remediation"
	UpdateCredentialOptionsPurposeRemediationConst                        = "remediation"
)

// NewUpdateCredentialOptions : Instantiate UpdateCredentialOptions
func (*PostureManagementV2) NewUpdateCredentialOptions(id string) *UpdateCredentialOptions {
	return &UpdateCredentialOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateCredentialOptions) SetID(id string) *UpdateCredentialOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateCredentialOptions) SetEnabled(enabled bool) *UpdateCredentialOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetType : Allow user to set Type
func (_options *UpdateCredentialOptions) SetType(typeVar string) *UpdateCredentialOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateCredentialOptions) SetName(name string) *UpdateCredentialOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateCredentialOptions) SetDescription(description string) *UpdateCredentialOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetDisplayFields : Allow user to set DisplayFields
func (_options *UpdateCredentialOptions) SetDisplayFields(displayFields *UpdateCredentialDisplayFields) *UpdateCredentialOptions {
	_options.DisplayFields = displayFields
	return _options
}

// SetPurpose : Allow user to set Purpose
func (_options *UpdateCredentialOptions) SetPurpose(purpose string) *UpdateCredentialOptions {
	_options.Purpose = core.StringPtr(purpose)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateCredentialOptions) SetAccountID(accountID string) *UpdateCredentialOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateCredentialOptions) SetTransactionID(transactionID string) *UpdateCredentialOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCredentialOptions) SetHeaders(param map[string]string) *UpdateCredentialOptions {
	options.Headers = param
	return options
}

// UpdateProfilesOptions : The UpdateProfiles options.
type UpdateProfilesOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// The name of the profile.
	Name *string `json:"name,omitempty"`

	// A description of the profile.
	Description *string `json:"description,omitempty"`

	// The base profile that the controls are pulled from.
	BaseProfile *string `json:"base_profile,omitempty"`

	// The type of profile. Seed profiles have the type set as 'predefined' and user-generated profiles have the type set
	// as 'custom'.
	Type *string `json:"type,omitempty"`

	// The profile status. If the profile is enabled, the value is true. If the profile is disabled, the value is false.
	IsEnabled *bool `json:"is_enabled,omitempty"`

	// A list of goal and control IDs that needs to be updated in the profile. These values can be retrieved from the
	// profiles/{profile_id}/controls API. The `profile_id` of the `base_profile` must be provided.
	ControlIds []string `json:"control_ids,omitempty"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateProfilesOptions.Type property.
// The type of profile. Seed profiles have the type set as 'predefined' and user-generated profiles have the type set as
// 'custom'.
const (
	UpdateProfilesOptionsTypeCustomConst     = "custom"
	UpdateProfilesOptionsTypePredefinedConst = "predefined"
)

// NewUpdateProfilesOptions : Instantiate UpdateProfilesOptions
func (*PostureManagementV2) NewUpdateProfilesOptions(id string) *UpdateProfilesOptions {
	return &UpdateProfilesOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateProfilesOptions) SetID(id string) *UpdateProfilesOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateProfilesOptions) SetName(name string) *UpdateProfilesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateProfilesOptions) SetDescription(description string) *UpdateProfilesOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetBaseProfile : Allow user to set BaseProfile
func (_options *UpdateProfilesOptions) SetBaseProfile(baseProfile string) *UpdateProfilesOptions {
	_options.BaseProfile = core.StringPtr(baseProfile)
	return _options
}

// SetType : Allow user to set Type
func (_options *UpdateProfilesOptions) SetType(typeVar string) *UpdateProfilesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetIsEnabled : Allow user to set IsEnabled
func (_options *UpdateProfilesOptions) SetIsEnabled(isEnabled bool) *UpdateProfilesOptions {
	_options.IsEnabled = core.BoolPtr(isEnabled)
	return _options
}

// SetControlIds : Allow user to set ControlIds
func (_options *UpdateProfilesOptions) SetControlIds(controlIds []string) *UpdateProfilesOptions {
	_options.ControlIds = controlIds
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateProfilesOptions) SetAccountID(accountID string) *UpdateProfilesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateProfilesOptions) SetTransactionID(transactionID string) *UpdateProfilesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateProfilesOptions) SetHeaders(param map[string]string) *UpdateProfilesOptions {
	options.Headers = param
	return options
}

// UpdateScopeDetailsOptions : The UpdateScopeDetails options.
type UpdateScopeDetailsOptions struct {
	// The ID for the API.
	ID *string `json:"-" validate:"required,ne="`

	// The name of the scope.
	Name *string `json:"name,omitempty"`

	// The description of the scope.
	Description *string `json:"description,omitempty"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"-"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request.
	TransactionID *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateScopeDetailsOptions : Instantiate UpdateScopeDetailsOptions
func (*PostureManagementV2) NewUpdateScopeDetailsOptions(id string) *UpdateScopeDetailsOptions {
	return &UpdateScopeDetailsOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateScopeDetailsOptions) SetID(id string) *UpdateScopeDetailsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateScopeDetailsOptions) SetName(name string) *UpdateScopeDetailsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateScopeDetailsOptions) SetDescription(description string) *UpdateScopeDetailsOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateScopeDetailsOptions) SetAccountID(accountID string) *UpdateScopeDetailsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateScopeDetailsOptions) SetTransactionID(transactionID string) *UpdateScopeDetailsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateScopeDetailsOptions) SetHeaders(param map[string]string) *UpdateScopeDetailsOptions {
	options.Headers = param
	return options
}
