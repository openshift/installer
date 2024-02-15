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
 * IBM OpenAPI SDK Code Generator Version: 3.82.1-2082d402-20231115-195014
 */

// Package secretsmanagerv2 : Operations and models for the SecretsManagerV2 service
package secretsmanagerv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/secrets-manager-go-sdk/v2/common"
	"github.com/go-openapi/strfmt"
)

// SecretsManagerV2 : With IBM Cloud® Secrets Manager, you can create, lease, and centrally manage secrets that are used
// in IBM Cloud services or your custom-built applications.
//
// API Version: 2.0.0
// See: https://cloud.ibm.com/docs/secrets-manager
type SecretsManagerV2 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://provide-here-your-smgr-instanceuuid.us-south.secrets-manager.appdomain.cloud"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "secrets_manager"

const ParameterizedServiceURL = "https://{instance_id}.{region}.secrets-manager.appdomain.cloud"

var defaultUrlVariables = map[string]string{
	"instance_id": "provide-here-your-smgr-instanceuuid",
	"region":      "us-south",
}

// SecretsManagerV2Options : Service options
type SecretsManagerV2Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewSecretsManagerV2UsingExternalConfig : constructs an instance of SecretsManagerV2 with passed in options and external configuration.
func NewSecretsManagerV2UsingExternalConfig(options *SecretsManagerV2Options) (secretsManager *SecretsManagerV2, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	secretsManager, err = NewSecretsManagerV2(options)
	if err != nil {
		return
	}

	err = secretsManager.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = secretsManager.Service.SetServiceURL(options.URL)
	}
	return
}

// NewSecretsManagerV2 : constructs an instance of SecretsManagerV2 with passed in options.
func NewSecretsManagerV2(options *SecretsManagerV2Options) (service *SecretsManagerV2, err error) {
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

	service = &SecretsManagerV2{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "secretsManager" suitable for processing requests.
func (secretsManager *SecretsManagerV2) Clone() *SecretsManagerV2 {
	if core.IsNil(secretsManager) {
		return nil
	}
	clone := *secretsManager
	clone.Service = secretsManager.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (secretsManager *SecretsManagerV2) SetServiceURL(url string) error {
	return secretsManager.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (secretsManager *SecretsManagerV2) GetServiceURL() string {
	return secretsManager.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (secretsManager *SecretsManagerV2) SetDefaultHeaders(headers http.Header) {
	secretsManager.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (secretsManager *SecretsManagerV2) SetEnableGzipCompression(enableGzip bool) {
	secretsManager.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (secretsManager *SecretsManagerV2) GetEnableGzipCompression() bool {
	return secretsManager.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (secretsManager *SecretsManagerV2) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	secretsManager.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (secretsManager *SecretsManagerV2) DisableRetries() {
	secretsManager.Service.DisableRetries()
}

// CreateSecretGroup : Create a new secret group
// Create a secret group that you can use to organize secrets and control who can access them.
//
// A successful request returns the ID value of the secret group, along with other properties. To learn more about
// secret groups, check out the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-secret-groups).
func (secretsManager *SecretsManagerV2) CreateSecretGroup(createSecretGroupOptions *CreateSecretGroupOptions) (result *SecretGroup, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretGroupWithContext(context.Background(), createSecretGroupOptions)
}

// CreateSecretGroupWithContext is an alternate form of the CreateSecretGroup method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateSecretGroupWithContext(ctx context.Context, createSecretGroupOptions *CreateSecretGroupOptions) (result *SecretGroup, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretGroupOptions, "createSecretGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSecretGroupOptions, "createSecretGroupOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secret_groups`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateSecretGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createSecretGroupOptions.Name != nil {
		body["name"] = createSecretGroupOptions.Name
	}
	if createSecretGroupOptions.Description != nil {
		body["description"] = createSecretGroupOptions.Description
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
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroup)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSecretGroups : List secret groups
// List the secret groups that are available in your Secrets Manager instance.
func (secretsManager *SecretsManagerV2) ListSecretGroups(listSecretGroupsOptions *ListSecretGroupsOptions) (result *SecretGroupCollection, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretGroupsWithContext(context.Background(), listSecretGroupsOptions)
}

// ListSecretGroupsWithContext is an alternate form of the ListSecretGroups method which supports a Context parameter
func (secretsManager *SecretsManagerV2) ListSecretGroupsWithContext(ctx context.Context, listSecretGroupsOptions *ListSecretGroupsOptions) (result *SecretGroupCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSecretGroupsOptions, "listSecretGroupsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secret_groups`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretGroupsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "ListSecretGroups")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecretGroup : Get a secret group
// Get the properties of an existing secret group by specifying the ID of the group.
func (secretsManager *SecretsManagerV2) GetSecretGroup(getSecretGroupOptions *GetSecretGroupOptions) (result *SecretGroup, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretGroupWithContext(context.Background(), getSecretGroupOptions)
}

// GetSecretGroupWithContext is an alternate form of the GetSecretGroup method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetSecretGroupWithContext(ctx context.Context, getSecretGroupOptions *GetSecretGroupOptions) (result *SecretGroup, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretGroupOptions, "getSecretGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretGroupOptions, "getSecretGroupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getSecretGroupOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secret_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetSecretGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroup)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecretGroup : Update a secret group
// Update the properties of an existing secret group, such as its name or description.
func (secretsManager *SecretsManagerV2) UpdateSecretGroup(updateSecretGroupOptions *UpdateSecretGroupOptions) (result *SecretGroup, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateSecretGroupWithContext(context.Background(), updateSecretGroupOptions)
}

// UpdateSecretGroupWithContext is an alternate form of the UpdateSecretGroup method which supports a Context parameter
func (secretsManager *SecretsManagerV2) UpdateSecretGroupWithContext(ctx context.Context, updateSecretGroupOptions *UpdateSecretGroupOptions) (result *SecretGroup, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecretGroupOptions, "updateSecretGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSecretGroupOptions, "updateSecretGroupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateSecretGroupOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secret_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSecretGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "UpdateSecretGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	_, err = builder.SetBodyContentJSON(updateSecretGroupOptions.SecretGroupPatch)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroup)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteSecretGroup : Delete a secret group
// Delete a secret group by specifying the ID of the secret group.
//
// **Note:** To delete a secret group, it must be empty. If you need to remove a secret group that contains secrets, you
// must first [delete the secrets](#delete-secret) that are associated with the group.
func (secretsManager *SecretsManagerV2) DeleteSecretGroup(deleteSecretGroupOptions *DeleteSecretGroupOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteSecretGroupWithContext(context.Background(), deleteSecretGroupOptions)
}

// DeleteSecretGroupWithContext is an alternate form of the DeleteSecretGroup method which supports a Context parameter
func (secretsManager *SecretsManagerV2) DeleteSecretGroupWithContext(ctx context.Context, deleteSecretGroupOptions *DeleteSecretGroupOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSecretGroupOptions, "deleteSecretGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSecretGroupOptions, "deleteSecretGroupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteSecretGroupOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secret_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSecretGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "DeleteSecretGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = secretsManager.Service.Request(request, nil)

	return
}

// CreateSecret : Create a new secret
// Create a secret or import an existing value that you can use to access or authenticate to a protected resource.
//
// Use this operation to either generate or import an existing secret, such as a TLS certificate, that you can manage in
// your Secrets Manager service instance. A successful request stores the secret in your dedicated instance, based on
// the secret type and data that you specify. The response returns the ID value of the secret, along with other
// metadata.
//
// To learn more about the types of secrets that you can create with Secrets Manager, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-what-is-secret).
func (secretsManager *SecretsManagerV2) CreateSecret(createSecretOptions *CreateSecretOptions) (result SecretIntf, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretWithContext(context.Background(), createSecretOptions)
}

// CreateSecretWithContext is an alternate form of the CreateSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateSecretWithContext(ctx context.Context, createSecretOptions *CreateSecretOptions) (result SecretIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretOptions, "createSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSecretOptions, "createSecretOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(createSecretOptions.SecretPrototype)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecret)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSecrets : List secrets
// List the secrets that are available in your Secrets Manager instance.
func (secretsManager *SecretsManagerV2) ListSecrets(listSecretsOptions *ListSecretsOptions) (result *SecretMetadataPaginatedCollection, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretsWithContext(context.Background(), listSecretsOptions)
}

// ListSecretsWithContext is an alternate form of the ListSecrets method which supports a Context parameter
func (secretsManager *SecretsManagerV2) ListSecretsWithContext(ctx context.Context, listSecretsOptions *ListSecretsOptions) (result *SecretMetadataPaginatedCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSecretsOptions, "listSecretsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "ListSecrets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSecretsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listSecretsOptions.Offset))
	}
	if listSecretsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSecretsOptions.Limit))
	}
	if listSecretsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listSecretsOptions.Sort))
	}
	if listSecretsOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listSecretsOptions.Search))
	}
	if listSecretsOptions.Groups != nil {
		builder.AddQuery("groups", strings.Join(listSecretsOptions.Groups, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretMetadataPaginatedCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecret : Get a secret
// Get a secret and its details by specifying the ID of the secret.
//
// A successful request returns the secret data that is associated with your secret, along with other metadata. To view
// only the details of a specified secret without retrieving its value, use the [Get secret
// metadata](#get-secret-metadata) operation.
func (secretsManager *SecretsManagerV2) GetSecret(getSecretOptions *GetSecretOptions) (result SecretIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretWithContext(context.Background(), getSecretOptions)
}

// GetSecretWithContext is an alternate form of the GetSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetSecretWithContext(ctx context.Context, getSecretOptions *GetSecretOptions) (result SecretIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretOptions, "getSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretOptions, "getSecretOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getSecretOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecret)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteSecret : Delete a secret
// Delete a secret by specifying the ID of the secret.
func (secretsManager *SecretsManagerV2) DeleteSecret(deleteSecretOptions *DeleteSecretOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteSecretWithContext(context.Background(), deleteSecretOptions)
}

// DeleteSecretWithContext is an alternate form of the DeleteSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV2) DeleteSecretWithContext(ctx context.Context, deleteSecretOptions *DeleteSecretOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSecretOptions, "deleteSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSecretOptions, "deleteSecretOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteSecretOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "DeleteSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = secretsManager.Service.Request(request, nil)

	return
}

// GetSecretMetadata : Get the metadata of a secret
// Get the metadata of a secret by specifying the ID of the secret.
func (secretsManager *SecretsManagerV2) GetSecretMetadata(getSecretMetadataOptions *GetSecretMetadataOptions) (result SecretMetadataIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretMetadataWithContext(context.Background(), getSecretMetadataOptions)
}

// GetSecretMetadataWithContext is an alternate form of the GetSecretMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetSecretMetadataWithContext(ctx context.Context, getSecretMetadataOptions *GetSecretMetadataOptions) (result SecretMetadataIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretMetadataOptions, "getSecretMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretMetadataOptions, "getSecretMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getSecretMetadataOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{id}/metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetSecretMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretMetadata)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecretMetadata : Update the metadata of a secret
// Update the metadata of a secret, such as its name or description.
func (secretsManager *SecretsManagerV2) UpdateSecretMetadata(updateSecretMetadataOptions *UpdateSecretMetadataOptions) (result SecretMetadataIntf, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateSecretMetadataWithContext(context.Background(), updateSecretMetadataOptions)
}

// UpdateSecretMetadataWithContext is an alternate form of the UpdateSecretMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV2) UpdateSecretMetadataWithContext(ctx context.Context, updateSecretMetadataOptions *UpdateSecretMetadataOptions) (result SecretMetadataIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecretMetadataOptions, "updateSecretMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSecretMetadataOptions, "updateSecretMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateSecretMetadataOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{id}/metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSecretMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "UpdateSecretMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	_, err = builder.SetBodyContentJSON(updateSecretMetadataOptions.SecretMetadataPatch)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretMetadata)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateSecretAction : Create a secret action
// Create a secret action. This operation supports the following actions:.
func (secretsManager *SecretsManagerV2) CreateSecretAction(createSecretActionOptions *CreateSecretActionOptions) (result SecretActionIntf, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretActionWithContext(context.Background(), createSecretActionOptions)
}

// CreateSecretActionWithContext is an alternate form of the CreateSecretAction method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateSecretActionWithContext(ctx context.Context, createSecretActionOptions *CreateSecretActionOptions) (result SecretActionIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretActionOptions, "createSecretActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSecretActionOptions, "createSecretActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *createSecretActionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{id}/actions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateSecretAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(createSecretActionOptions.SecretActionPrototype)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretAction)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecretByNameType : Get a secret by name
// Get a secret and its details by specifying the Name and Type of the secret.
//
// A successful request returns the secret data that is associated with your secret, along with other metadata. To view
// only the details of a specified secret without retrieving its value, use the [Get secret
// metadata](#get-secret-metadata) operation.
func (secretsManager *SecretsManagerV2) GetSecretByNameType(getSecretByNameTypeOptions *GetSecretByNameTypeOptions) (result SecretIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretByNameTypeWithContext(context.Background(), getSecretByNameTypeOptions)
}

// GetSecretByNameTypeWithContext is an alternate form of the GetSecretByNameType method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetSecretByNameTypeWithContext(ctx context.Context, getSecretByNameTypeOptions *GetSecretByNameTypeOptions) (result SecretIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretByNameTypeOptions, "getSecretByNameTypeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretByNameTypeOptions, "getSecretByNameTypeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type":       *getSecretByNameTypeOptions.SecretType,
		"name":              *getSecretByNameTypeOptions.Name,
		"secret_group_name": *getSecretByNameTypeOptions.SecretGroupName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secret_groups/{secret_group_name}/secret_types/{secret_type}/secrets/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretByNameTypeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetSecretByNameType")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecret)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateSecretVersion : Create a new secret version
// Create a new secret version.
func (secretsManager *SecretsManagerV2) CreateSecretVersion(createSecretVersionOptions *CreateSecretVersionOptions) (result SecretVersionIntf, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretVersionWithContext(context.Background(), createSecretVersionOptions)
}

// CreateSecretVersionWithContext is an alternate form of the CreateSecretVersion method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateSecretVersionWithContext(ctx context.Context, createSecretVersionOptions *CreateSecretVersionOptions) (result SecretVersionIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretVersionOptions, "createSecretVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSecretVersionOptions, "createSecretVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *createSecretVersionOptions.SecretID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateSecretVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(createSecretVersionOptions.SecretVersionPrototype)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSecretVersions : List versions of a secret
// List the versions of a secret.
//
// A successful request returns the list of versions of a secret, along with the metadata of each version.
func (secretsManager *SecretsManagerV2) ListSecretVersions(listSecretVersionsOptions *ListSecretVersionsOptions) (result *SecretVersionMetadataCollection, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretVersionsWithContext(context.Background(), listSecretVersionsOptions)
}

// ListSecretVersionsWithContext is an alternate form of the ListSecretVersions method which supports a Context parameter
func (secretsManager *SecretsManagerV2) ListSecretVersionsWithContext(ctx context.Context, listSecretVersionsOptions *ListSecretVersionsOptions) (result *SecretVersionMetadataCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSecretVersionsOptions, "listSecretVersionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listSecretVersionsOptions, "listSecretVersionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *listSecretVersionsOptions.SecretID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "ListSecretVersions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretVersionMetadataCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecretVersion : Get a version of a secret
// Get a version of a secret by specifying the ID of the version. You can use the `current` or `previous` aliases to
// refer to the current or previous secret version.
//
// A successful request returns the secret data that is associated with the specified version of your secret, along with
// other metadata.
func (secretsManager *SecretsManagerV2) GetSecretVersion(getSecretVersionOptions *GetSecretVersionOptions) (result SecretVersionIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretVersionWithContext(context.Background(), getSecretVersionOptions)
}

// GetSecretVersionWithContext is an alternate form of the GetSecretVersion method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetSecretVersionWithContext(ctx context.Context, getSecretVersionOptions *GetSecretVersionOptions) (result SecretVersionIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretVersionOptions, "getSecretVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretVersionOptions, "getSecretVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *getSecretVersionOptions.SecretID,
		"id":        *getSecretVersionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetSecretVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteSecretVersionData : Delete the data of a secret version
// Delete the data of a secret version by specifying the ID of the version.
//
// This operation is available for secret type: iam_credentials current version.
func (secretsManager *SecretsManagerV2) DeleteSecretVersionData(deleteSecretVersionDataOptions *DeleteSecretVersionDataOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteSecretVersionDataWithContext(context.Background(), deleteSecretVersionDataOptions)
}

// DeleteSecretVersionDataWithContext is an alternate form of the DeleteSecretVersionData method which supports a Context parameter
func (secretsManager *SecretsManagerV2) DeleteSecretVersionDataWithContext(ctx context.Context, deleteSecretVersionDataOptions *DeleteSecretVersionDataOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSecretVersionDataOptions, "deleteSecretVersionDataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSecretVersionDataOptions, "deleteSecretVersionDataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *deleteSecretVersionDataOptions.SecretID,
		"id":        *deleteSecretVersionDataOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions/{id}/secret_data`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSecretVersionDataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "DeleteSecretVersionData")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = secretsManager.Service.Request(request, nil)

	return
}

// GetSecretVersionMetadata : Get the metadata of a secret version
// Get the metadata of a secret version by specifying the ID of the version. You can use the `current` or `previous`
// aliases to refer to the current or previous secret version.
//
// A successful request returns the metadata that is associated with the specified version of your secret.
func (secretsManager *SecretsManagerV2) GetSecretVersionMetadata(getSecretVersionMetadataOptions *GetSecretVersionMetadataOptions) (result SecretVersionMetadataIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretVersionMetadataWithContext(context.Background(), getSecretVersionMetadataOptions)
}

// GetSecretVersionMetadataWithContext is an alternate form of the GetSecretVersionMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetSecretVersionMetadataWithContext(ctx context.Context, getSecretVersionMetadataOptions *GetSecretVersionMetadataOptions) (result SecretVersionMetadataIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretVersionMetadataOptions, "getSecretVersionMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretVersionMetadataOptions, "getSecretVersionMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *getSecretVersionMetadataOptions.SecretID,
		"id":        *getSecretVersionMetadataOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions/{id}/metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretVersionMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetSecretVersionMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretVersionMetadata)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecretVersionMetadata : Update the metadata of a secret version
// Update the custom metadata of a secret version.
func (secretsManager *SecretsManagerV2) UpdateSecretVersionMetadata(updateSecretVersionMetadataOptions *UpdateSecretVersionMetadataOptions) (result SecretVersionMetadataIntf, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateSecretVersionMetadataWithContext(context.Background(), updateSecretVersionMetadataOptions)
}

// UpdateSecretVersionMetadataWithContext is an alternate form of the UpdateSecretVersionMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV2) UpdateSecretVersionMetadataWithContext(ctx context.Context, updateSecretVersionMetadataOptions *UpdateSecretVersionMetadataOptions) (result SecretVersionMetadataIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecretVersionMetadataOptions, "updateSecretVersionMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSecretVersionMetadataOptions, "updateSecretVersionMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *updateSecretVersionMetadataOptions.SecretID,
		"id":        *updateSecretVersionMetadataOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions/{id}/metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSecretVersionMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "UpdateSecretVersionMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	_, err = builder.SetBodyContentJSON(updateSecretVersionMetadataOptions.SecretVersionMetadataPatch)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretVersionMetadata)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateSecretVersionAction : Create a version action
// Create a secret version action. This operation supports the following actions:
//
// - `private_cert_action_revoke_certificate`: Revoke a version of a private certificate.
func (secretsManager *SecretsManagerV2) CreateSecretVersionAction(createSecretVersionActionOptions *CreateSecretVersionActionOptions) (result VersionActionIntf, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretVersionActionWithContext(context.Background(), createSecretVersionActionOptions)
}

// CreateSecretVersionActionWithContext is an alternate form of the CreateSecretVersionAction method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateSecretVersionActionWithContext(ctx context.Context, createSecretVersionActionOptions *CreateSecretVersionActionOptions) (result VersionActionIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretVersionActionOptions, "createSecretVersionActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSecretVersionActionOptions, "createSecretVersionActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *createSecretVersionActionOptions.SecretID,
		"id":        *createSecretVersionActionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions/{id}/actions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretVersionActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateSecretVersionAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(createSecretVersionActionOptions.SecretVersionActionPrototype)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVersionAction)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSecretsLocks : List secrets and their locks
// List the secrets and their locks in your Secrets Manager instance.
func (secretsManager *SecretsManagerV2) ListSecretsLocks(listSecretsLocksOptions *ListSecretsLocksOptions) (result *SecretsLocksPaginatedCollection, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretsLocksWithContext(context.Background(), listSecretsLocksOptions)
}

// ListSecretsLocksWithContext is an alternate form of the ListSecretsLocks method which supports a Context parameter
func (secretsManager *SecretsManagerV2) ListSecretsLocksWithContext(ctx context.Context, listSecretsLocksOptions *ListSecretsLocksOptions) (result *SecretsLocksPaginatedCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSecretsLocksOptions, "listSecretsLocksOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets_locks`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretsLocksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "ListSecretsLocks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSecretsLocksOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listSecretsLocksOptions.Offset))
	}
	if listSecretsLocksOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSecretsLocksOptions.Limit))
	}
	if listSecretsLocksOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listSecretsLocksOptions.Search))
	}
	if listSecretsLocksOptions.Groups != nil {
		builder.AddQuery("groups", strings.Join(listSecretsLocksOptions.Groups, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretsLocksPaginatedCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSecretLocks : List secret locks
// List the locks that are associated with a specified secret.
func (secretsManager *SecretsManagerV2) ListSecretLocks(listSecretLocksOptions *ListSecretLocksOptions) (result *SecretLocksPaginatedCollection, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretLocksWithContext(context.Background(), listSecretLocksOptions)
}

// ListSecretLocksWithContext is an alternate form of the ListSecretLocks method which supports a Context parameter
func (secretsManager *SecretsManagerV2) ListSecretLocksWithContext(ctx context.Context, listSecretLocksOptions *ListSecretLocksOptions) (result *SecretLocksPaginatedCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSecretLocksOptions, "listSecretLocksOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listSecretLocksOptions, "listSecretLocksOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *listSecretLocksOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{id}/locks`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretLocksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "ListSecretLocks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSecretLocksOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listSecretLocksOptions.Offset))
	}
	if listSecretLocksOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSecretLocksOptions.Limit))
	}
	if listSecretLocksOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listSecretLocksOptions.Sort))
	}
	if listSecretLocksOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listSecretLocksOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretLocksPaginatedCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateSecretLocksBulk : Create secret locks
// Create a lock on the current version of a secret.
//
// A lock can be used to prevent a secret from being deleted or modified while it's in use by your applications. A
// successful request attaches a new lock to your secret, or replaces a lock of the same name if it already exists.
// Additionally, you can use this operation to clear any matching locks on a secret by using one of the following
// optional lock modes:
//
// - `remove_previous`: Removes any other locks with matching names if they are found in the previous version of the
// secret.\n
// - `remove_previous_and_delete`: Carries out the same function as `remove_previous`, but also permanently deletes the
// data of the previous secret version if it doesn't have any locks.
func (secretsManager *SecretsManagerV2) CreateSecretLocksBulk(createSecretLocksBulkOptions *CreateSecretLocksBulkOptions) (result *SecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretLocksBulkWithContext(context.Background(), createSecretLocksBulkOptions)
}

// CreateSecretLocksBulkWithContext is an alternate form of the CreateSecretLocksBulk method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateSecretLocksBulkWithContext(ctx context.Context, createSecretLocksBulkOptions *CreateSecretLocksBulkOptions) (result *SecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretLocksBulkOptions, "createSecretLocksBulkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSecretLocksBulkOptions, "createSecretLocksBulkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *createSecretLocksBulkOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{id}/locks_bulk`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretLocksBulkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateSecretLocksBulk")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createSecretLocksBulkOptions.Mode != nil {
		builder.AddQuery("mode", fmt.Sprint(*createSecretLocksBulkOptions.Mode))
	}

	body := make(map[string]interface{})
	if createSecretLocksBulkOptions.Locks != nil {
		body["locks"] = createSecretLocksBulkOptions.Locks
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
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteSecretLocksBulk : Delete secret locks
// Delete all the locks or a subset of the locks that are associated with a version of a secret.
//
// To delete only a subset of the locks, add a query param with a comma to separate the list of lock names:
//
// Example: `?name=lock-example-1,lock-example-2`.
//
// **Note:** A secret is considered unlocked and able to be deleted only after you remove all of its locks. To determine
// whether a secret contains locks, check the `locks_total` field that is returned as part of the metadata of your
// secret.
func (secretsManager *SecretsManagerV2) DeleteSecretLocksBulk(deleteSecretLocksBulkOptions *DeleteSecretLocksBulkOptions) (result *SecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.DeleteSecretLocksBulkWithContext(context.Background(), deleteSecretLocksBulkOptions)
}

// DeleteSecretLocksBulkWithContext is an alternate form of the DeleteSecretLocksBulk method which supports a Context parameter
func (secretsManager *SecretsManagerV2) DeleteSecretLocksBulkWithContext(ctx context.Context, deleteSecretLocksBulkOptions *DeleteSecretLocksBulkOptions) (result *SecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSecretLocksBulkOptions, "deleteSecretLocksBulkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSecretLocksBulkOptions, "deleteSecretLocksBulkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteSecretLocksBulkOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{id}/locks_bulk`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSecretLocksBulkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "DeleteSecretLocksBulk")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if deleteSecretLocksBulkOptions.Name != nil {
		builder.AddQuery("name", strings.Join(deleteSecretLocksBulkOptions.Name, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSecretVersionLocks : List secret version locks
// List the locks that are associated with a specified secret version.
func (secretsManager *SecretsManagerV2) ListSecretVersionLocks(listSecretVersionLocksOptions *ListSecretVersionLocksOptions) (result *SecretVersionLocksPaginatedCollection, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretVersionLocksWithContext(context.Background(), listSecretVersionLocksOptions)
}

// ListSecretVersionLocksWithContext is an alternate form of the ListSecretVersionLocks method which supports a Context parameter
func (secretsManager *SecretsManagerV2) ListSecretVersionLocksWithContext(ctx context.Context, listSecretVersionLocksOptions *ListSecretVersionLocksOptions) (result *SecretVersionLocksPaginatedCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSecretVersionLocksOptions, "listSecretVersionLocksOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listSecretVersionLocksOptions, "listSecretVersionLocksOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *listSecretVersionLocksOptions.SecretID,
		"id":        *listSecretVersionLocksOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions/{id}/locks`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretVersionLocksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "ListSecretVersionLocks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSecretVersionLocksOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listSecretVersionLocksOptions.Offset))
	}
	if listSecretVersionLocksOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSecretVersionLocksOptions.Limit))
	}
	if listSecretVersionLocksOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listSecretVersionLocksOptions.Sort))
	}
	if listSecretVersionLocksOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listSecretVersionLocksOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretVersionLocksPaginatedCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateSecretVersionLocksBulk : Create secret version locks
// Create a lock on the specified version of a secret.
//
// A lock can be used to prevent a secret from being deleted or modified while it's in use by your applications. A
// successful request attaches a new lock to your secret, or replaces a lock of the same name if it already exists.
// Additionally, you can use this operation to clear any matching locks on a secret by using one of the following
// optional lock modes:
//
// - `remove_previous`: Removes any other locks with matching names if they are found in the previous version of the
// secret.
// - `remove_previous_and_delete`: Carries out the same function as `remove_previous`, but also permanently deletes the
// data of the previous secret version if it doesn't have any locks.
func (secretsManager *SecretsManagerV2) CreateSecretVersionLocksBulk(createSecretVersionLocksBulkOptions *CreateSecretVersionLocksBulkOptions) (result *SecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretVersionLocksBulkWithContext(context.Background(), createSecretVersionLocksBulkOptions)
}

// CreateSecretVersionLocksBulkWithContext is an alternate form of the CreateSecretVersionLocksBulk method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateSecretVersionLocksBulkWithContext(ctx context.Context, createSecretVersionLocksBulkOptions *CreateSecretVersionLocksBulkOptions) (result *SecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretVersionLocksBulkOptions, "createSecretVersionLocksBulkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSecretVersionLocksBulkOptions, "createSecretVersionLocksBulkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *createSecretVersionLocksBulkOptions.SecretID,
		"id":        *createSecretVersionLocksBulkOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions/{id}/locks_bulk`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretVersionLocksBulkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateSecretVersionLocksBulk")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createSecretVersionLocksBulkOptions.Mode != nil {
		builder.AddQuery("mode", fmt.Sprint(*createSecretVersionLocksBulkOptions.Mode))
	}

	body := make(map[string]interface{})
	if createSecretVersionLocksBulkOptions.Locks != nil {
		body["locks"] = createSecretVersionLocksBulkOptions.Locks
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
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteSecretVersionLocksBulk : Delete locks on a secret version
// Delete all the locks or a subset of the locks that are associated with the specified version of a secret.
//
// To delete only a subset of the locks, add a query param with a comma to separate the list of lock names:
//
// Example: `?name=lock-example-1,lock-example-2`.
//
// **Note:** A secret is considered unlocked and able to be deleted only after all of its locks are removed. To
// determine whether a secret contains locks, check the `locks_total` field that is returned as part of the metadata of
// your secret.
func (secretsManager *SecretsManagerV2) DeleteSecretVersionLocksBulk(deleteSecretVersionLocksBulkOptions *DeleteSecretVersionLocksBulkOptions) (result *SecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.DeleteSecretVersionLocksBulkWithContext(context.Background(), deleteSecretVersionLocksBulkOptions)
}

// DeleteSecretVersionLocksBulkWithContext is an alternate form of the DeleteSecretVersionLocksBulk method which supports a Context parameter
func (secretsManager *SecretsManagerV2) DeleteSecretVersionLocksBulkWithContext(ctx context.Context, deleteSecretVersionLocksBulkOptions *DeleteSecretVersionLocksBulkOptions) (result *SecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSecretVersionLocksBulkOptions, "deleteSecretVersionLocksBulkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSecretVersionLocksBulkOptions, "deleteSecretVersionLocksBulkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_id": *deleteSecretVersionLocksBulkOptions.SecretID,
		"id":        *deleteSecretVersionLocksBulkOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/secrets/{secret_id}/versions/{id}/locks_bulk`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSecretVersionLocksBulkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "DeleteSecretVersionLocksBulk")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if deleteSecretVersionLocksBulkOptions.Name != nil {
		builder.AddQuery("name", strings.Join(deleteSecretVersionLocksBulkOptions.Name, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateConfiguration : Create a new configuration
// Add a configuration to the specified secret type.
//
// Use this operation to define the configurations that are required to create public certificates (`public_cert`),
// private certificates (`private_cert`) and IAM Credentials secrets (`iam_credentials`).
//
// You can add multiple configurations for your instance as follows:
//
// - A single configuration for IAM Credentials.
// - Up to 10 CA configurations for public certificates.
// - Up to 10 DNS configurations for public certificates.
// - Up to 10 Root CA configurations for private certificates.
// - Up to 10 Intermediate CA configurations for private certificates.
// - Up to 10 Certificate Template configurations for private certificates.
func (secretsManager *SecretsManagerV2) CreateConfiguration(createConfigurationOptions *CreateConfigurationOptions) (result ConfigurationIntf, response *core.DetailedResponse, err error) {
	return secretsManager.CreateConfigurationWithContext(context.Background(), createConfigurationOptions)
}

// CreateConfigurationWithContext is an alternate form of the CreateConfiguration method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateConfigurationWithContext(ctx context.Context, createConfigurationOptions *CreateConfigurationOptions) (result ConfigurationIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createConfigurationOptions, "createConfigurationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createConfigurationOptions, "createConfigurationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/configurations`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(createConfigurationOptions.ConfigurationPrototype)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfiguration)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListConfigurations : List configurations
// List the configurations that are available in your Secrets Manager instance.
func (secretsManager *SecretsManagerV2) ListConfigurations(listConfigurationsOptions *ListConfigurationsOptions) (result *ConfigurationMetadataPaginatedCollection, response *core.DetailedResponse, err error) {
	return secretsManager.ListConfigurationsWithContext(context.Background(), listConfigurationsOptions)
}

// ListConfigurationsWithContext is an alternate form of the ListConfigurations method which supports a Context parameter
func (secretsManager *SecretsManagerV2) ListConfigurationsWithContext(ctx context.Context, listConfigurationsOptions *ListConfigurationsOptions) (result *ConfigurationMetadataPaginatedCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listConfigurationsOptions, "listConfigurationsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/configurations`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listConfigurationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "ListConfigurations")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listConfigurationsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listConfigurationsOptions.Offset))
	}
	if listConfigurationsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listConfigurationsOptions.Limit))
	}
	if listConfigurationsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listConfigurationsOptions.Sort))
	}
	if listConfigurationsOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listConfigurationsOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigurationMetadataPaginatedCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetConfiguration : Get a configuration
// Get a configuration by specifying its name.
//
// A successful request returns the details of your configuration.
func (secretsManager *SecretsManagerV2) GetConfiguration(getConfigurationOptions *GetConfigurationOptions) (result ConfigurationIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetConfigurationWithContext(context.Background(), getConfigurationOptions)
}

// GetConfigurationWithContext is an alternate form of the GetConfiguration method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetConfigurationWithContext(ctx context.Context, getConfigurationOptions *GetConfigurationOptions) (result ConfigurationIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigurationOptions, "getConfigurationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getConfigurationOptions, "getConfigurationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"name": *getConfigurationOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/configurations/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getConfigurationOptions.XSmAcceptConfigurationType != nil {
		builder.AddHeader("X-Sm-Accept-Configuration-Type", fmt.Sprint(*getConfigurationOptions.XSmAcceptConfigurationType))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfiguration)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateConfiguration : Update configuration
// Update a configuration.
func (secretsManager *SecretsManagerV2) UpdateConfiguration(updateConfigurationOptions *UpdateConfigurationOptions) (result ConfigurationIntf, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateConfigurationWithContext(context.Background(), updateConfigurationOptions)
}

// UpdateConfigurationWithContext is an alternate form of the UpdateConfiguration method which supports a Context parameter
func (secretsManager *SecretsManagerV2) UpdateConfigurationWithContext(ctx context.Context, updateConfigurationOptions *UpdateConfigurationOptions) (result ConfigurationIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateConfigurationOptions, "updateConfigurationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateConfigurationOptions, "updateConfigurationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"name": *updateConfigurationOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/configurations/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "UpdateConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateConfigurationOptions.XSmAcceptConfigurationType != nil {
		builder.AddHeader("X-Sm-Accept-Configuration-Type", fmt.Sprint(*updateConfigurationOptions.XSmAcceptConfigurationType))
	}

	_, err = builder.SetBodyContentJSON(updateConfigurationOptions.ConfigurationPatch)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfiguration)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteConfiguration : Delete a configuration
// Delete a configuration by specifying its name.
func (secretsManager *SecretsManagerV2) DeleteConfiguration(deleteConfigurationOptions *DeleteConfigurationOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteConfigurationWithContext(context.Background(), deleteConfigurationOptions)
}

// DeleteConfigurationWithContext is an alternate form of the DeleteConfiguration method which supports a Context parameter
func (secretsManager *SecretsManagerV2) DeleteConfigurationWithContext(ctx context.Context, deleteConfigurationOptions *DeleteConfigurationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteConfigurationOptions, "deleteConfigurationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteConfigurationOptions, "deleteConfigurationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"name": *deleteConfigurationOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/configurations/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "DeleteConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteConfigurationOptions.XSmAcceptConfigurationType != nil {
		builder.AddHeader("X-Sm-Accept-Configuration-Type", fmt.Sprint(*deleteConfigurationOptions.XSmAcceptConfigurationType))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = secretsManager.Service.Request(request, nil)

	return
}

// CreateConfigurationAction : Create a configuration action
// Create a configuration action. This operation supports the following actions:
//
// - `private_cert_configuration_action_sign_intermediate`: Sign an intermediate certificate authority.
// - `private_cert_configuration_action_sign_csr`: Sign a certificate signing request.
// - `private_cert_configuration_action_set_signed`: Set a signed intermediate certificate authority.
// - `private_cert_configuration_action_revoke_ca_certificate`: Revoke an internally signed intermediate certificate
// authority certificate.
// - `private_cert_configuration_action_rotate_crl`: Rotate the certificate revocation list (CRL) of an intermediate
// certificate authority.
func (secretsManager *SecretsManagerV2) CreateConfigurationAction(createConfigurationActionOptions *CreateConfigurationActionOptions) (result ConfigurationActionIntf, response *core.DetailedResponse, err error) {
	return secretsManager.CreateConfigurationActionWithContext(context.Background(), createConfigurationActionOptions)
}

// CreateConfigurationActionWithContext is an alternate form of the CreateConfigurationAction method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateConfigurationActionWithContext(ctx context.Context, createConfigurationActionOptions *CreateConfigurationActionOptions) (result ConfigurationActionIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createConfigurationActionOptions, "createConfigurationActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createConfigurationActionOptions, "createConfigurationActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"name": *createConfigurationActionOptions.Name,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/configurations/{name}/actions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createConfigurationActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateConfigurationAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createConfigurationActionOptions.XSmAcceptConfigurationType != nil {
		builder.AddHeader("X-Sm-Accept-Configuration-Type", fmt.Sprint(*createConfigurationActionOptions.XSmAcceptConfigurationType))
	}

	_, err = builder.SetBodyContentJSON(createConfigurationActionOptions.ConfigActionPrototype)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigurationAction)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateNotificationsRegistration : Register with Event Notifications instance
// Create a registration between a Secrets Manager instance and [Event
// Notifications](https://cloud.ibm.com/apidocs/event-notifications).
//
// A successful request adds Secrets Manager as a source that you can reference from your Event Notifications instance.
// For more information about enabling notifications for Secrets Manager, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-event-notifications).
func (secretsManager *SecretsManagerV2) CreateNotificationsRegistration(createNotificationsRegistrationOptions *CreateNotificationsRegistrationOptions) (result *NotificationsRegistration, response *core.DetailedResponse, err error) {
	return secretsManager.CreateNotificationsRegistrationWithContext(context.Background(), createNotificationsRegistrationOptions)
}

// CreateNotificationsRegistrationWithContext is an alternate form of the CreateNotificationsRegistration method which supports a Context parameter
func (secretsManager *SecretsManagerV2) CreateNotificationsRegistrationWithContext(ctx context.Context, createNotificationsRegistrationOptions *CreateNotificationsRegistrationOptions) (result *NotificationsRegistration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createNotificationsRegistrationOptions, "createNotificationsRegistrationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createNotificationsRegistrationOptions, "createNotificationsRegistrationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/notifications/registration`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createNotificationsRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "CreateNotificationsRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createNotificationsRegistrationOptions.EventNotificationsInstanceCrn != nil {
		body["event_notifications_instance_crn"] = createNotificationsRegistrationOptions.EventNotificationsInstanceCrn
	}
	if createNotificationsRegistrationOptions.EventNotificationsSourceName != nil {
		body["event_notifications_source_name"] = createNotificationsRegistrationOptions.EventNotificationsSourceName
	}
	if createNotificationsRegistrationOptions.EventNotificationsSourceDescription != nil {
		body["event_notifications_source_description"] = createNotificationsRegistrationOptions.EventNotificationsSourceDescription
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
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNotificationsRegistration)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetNotificationsRegistration : Get Event Notifications registration details
// Get the details of the registration between your Secrets Manager instance and Event Notifications.
func (secretsManager *SecretsManagerV2) GetNotificationsRegistration(getNotificationsRegistrationOptions *GetNotificationsRegistrationOptions) (result *NotificationsRegistration, response *core.DetailedResponse, err error) {
	return secretsManager.GetNotificationsRegistrationWithContext(context.Background(), getNotificationsRegistrationOptions)
}

// GetNotificationsRegistrationWithContext is an alternate form of the GetNotificationsRegistration method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetNotificationsRegistrationWithContext(ctx context.Context, getNotificationsRegistrationOptions *GetNotificationsRegistrationOptions) (result *NotificationsRegistration, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getNotificationsRegistrationOptions, "getNotificationsRegistrationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/notifications/registration`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getNotificationsRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetNotificationsRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNotificationsRegistration)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteNotificationsRegistration : Unregister from Event Notifications instance
// Delete the registration between your Secrets Manager instance and Event Notifications.
//
// A successful request removes your Secrets Manager instance as a source in Event Notifications.
func (secretsManager *SecretsManagerV2) DeleteNotificationsRegistration(deleteNotificationsRegistrationOptions *DeleteNotificationsRegistrationOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteNotificationsRegistrationWithContext(context.Background(), deleteNotificationsRegistrationOptions)
}

// DeleteNotificationsRegistrationWithContext is an alternate form of the DeleteNotificationsRegistration method which supports a Context parameter
func (secretsManager *SecretsManagerV2) DeleteNotificationsRegistrationWithContext(ctx context.Context, deleteNotificationsRegistrationOptions *DeleteNotificationsRegistrationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(deleteNotificationsRegistrationOptions, "deleteNotificationsRegistrationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/notifications/registration`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteNotificationsRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "DeleteNotificationsRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = secretsManager.Service.Request(request, nil)

	return
}

// GetNotificationsRegistrationTest : Send a test event for Event Notifications registrations
// Send a test event from a Secrets Manager instance to a configured [Event
// Notifications](https://cloud.ibm.com/apidocs/event-notifications) instance.
//
// A successful request sends a test event to the Event Notifications instance. For more information about enabling
// notifications for Secrets Manager, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-event-notifications).
func (secretsManager *SecretsManagerV2) GetNotificationsRegistrationTest(getNotificationsRegistrationTestOptions *GetNotificationsRegistrationTestOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.GetNotificationsRegistrationTestWithContext(context.Background(), getNotificationsRegistrationTestOptions)
}

// GetNotificationsRegistrationTestWithContext is an alternate form of the GetNotificationsRegistrationTest method which supports a Context parameter
func (secretsManager *SecretsManagerV2) GetNotificationsRegistrationTestWithContext(ctx context.Context, getNotificationsRegistrationTestOptions *GetNotificationsRegistrationTestOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getNotificationsRegistrationTestOptions, "getNotificationsRegistrationTestOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v2/notifications/registration/test`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getNotificationsRegistrationTestOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V2", "GetNotificationsRegistrationTest")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = secretsManager.Service.Request(request, nil)

	return
}

// CertificateIssuanceInfo : Issuance information that is associated with your certificate.
type CertificateIssuanceInfo struct {
	// This parameter indicates whether the issued certificate is configured with an automatic rotation policy.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The set of challenges. It is returned only when ordering public certificates by using manual DNS configuration.
	Challenges []ChallengeResource `json:"challenges,omitempty"`

	// The date that a user requests to validate DNS challenges for certificates that are ordered with a manual DNS
	// provider. The date format follows `RFC 3339`.
	DnsChallengeValidationTime *strfmt.DateTime `json:"dns_challenge_validation_time,omitempty"`

	// A code that identifies an issuance error.
	//
	// This field, along with `error_message`, is returned when Secrets Manager successfully processes your request, but
	// the certificate authority is unable to issue a certificate.
	ErrorCode *string `json:"error_code,omitempty"`

	// A human-readable message that provides details about the issuance error.
	ErrorMessage *string `json:"error_message,omitempty"`

	// The date when the certificate is ordered. The date format follows `RFC 3339`.
	OrderedOn *strfmt.DateTime `json:"ordered_on,omitempty"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`
}

// Constants associated with the CertificateIssuanceInfo.StateDescription property.
// A text representation of the secret state.
const (
	CertificateIssuanceInfo_StateDescription_Active        = "active"
	CertificateIssuanceInfo_StateDescription_Deactivated   = "deactivated"
	CertificateIssuanceInfo_StateDescription_Destroyed     = "destroyed"
	CertificateIssuanceInfo_StateDescription_PreActivation = "pre_activation"
	CertificateIssuanceInfo_StateDescription_Suspended     = "suspended"
)

// UnmarshalCertificateIssuanceInfo unmarshals an instance of CertificateIssuanceInfo from the specified map of raw messages.
func UnmarshalCertificateIssuanceInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateIssuanceInfo)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "challenges", &obj.Challenges, UnmarshalChallengeResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dns_challenge_validation_time", &obj.DnsChallengeValidationTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_code", &obj.ErrorCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_message", &obj.ErrorMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ordered_on", &obj.OrderedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CertificateValidity : The date and time that the certificate validity period begins and ends.
type CertificateValidity struct {
	// The date-time format follows `RFC 3339`.
	NotBefore *strfmt.DateTime `json:"not_before" validate:"required"`

	// The date-time format follows `RFC 3339`.
	NotAfter *strfmt.DateTime `json:"not_after" validate:"required"`
}

// UnmarshalCertificateValidity unmarshals an instance of CertificateValidity from the specified map of raw messages.
func UnmarshalCertificateValidity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateValidity)
	err = core.UnmarshalPrimitive(m, "not_before", &obj.NotBefore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_after", &obj.NotAfter)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ChallengeResource : Properties that describe a challenge.
type ChallengeResource struct {
	// The challenge domain.
	Domain *string `json:"domain,omitempty"`

	// The challenge expiration date. The date format follows `RFC 3339`.
	Expiration *strfmt.DateTime `json:"expiration,omitempty"`

	// The challenge status.
	Status *string `json:"status,omitempty"`

	// The TXT record name.
	TxtRecordName *string `json:"txt_record_name,omitempty"`

	// The TXT record value.
	TxtRecordValue *string `json:"txt_record_value,omitempty"`
}

// UnmarshalChallengeResource unmarshals an instance of ChallengeResource from the specified map of raw messages.
func UnmarshalChallengeResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ChallengeResource)
	err = core.UnmarshalPrimitive(m, "domain", &obj.Domain)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration", &obj.Expiration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "txt_record_name", &obj.TxtRecordName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "txt_record_value", &obj.TxtRecordValue)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Configuration : Your configuration.
// Models which "extend" this model:
// - PublicCertificateConfigurationCALetsEncrypt
// - PublicCertificateConfigurationDNSCloudInternetServices
// - PublicCertificateConfigurationDNSClassicInfrastructure
// - IAMCredentialsConfiguration
// - PrivateCertificateConfigurationRootCA
// - PrivateCertificateConfigurationIntermediateCA
// - PrivateCertificateConfigurationTemplate
type Configuration struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type,omitempty"`

	// The unique name of your configuration.
	Name *string `json:"name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The configuration of the Let's Encrypt CA environment.
	LetsEncryptEnvironment *string `json:"lets_encrypt_environment,omitempty"`

	// If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name. If
	// no match, the default offered chain will be used.
	LetsEncryptPreferredChain *string `json:"lets_encrypt_preferred_chain,omitempty"`

	// The PEM-encoded private key of your Let's Encrypt account. The data must be formatted on a single line with embedded
	// newline characters.
	LetsEncryptPrivateKey *string `json:"lets_encrypt_private_key,omitempty"`

	// An IBM Cloud API key that can to list domains in your Cloud Internet Services instance.
	//
	// To grant Secrets Manager the ability to view the Cloud Internet Services instance and all of its domains, the API
	// key must be assigned the Reader service role on Internet Services (`internet-svcs`).
	//
	// If you need to manage specific domains, you can assign the Manager role. For production environments, it is
	// recommended that you assign the Reader access role, and then use the
	// [IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific
	// domains. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	CloudInternetServicesApikey *string `json:"cloud_internet_services_apikey,omitempty"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	CloudInternetServicesCrn *string `json:"cloud_internet_services_crn,omitempty"`

	// The username that is associated with your classic infrastructure account.
	//
	// In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information,
	// see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructureUsername *string `json:"classic_infrastructure_username,omitempty"`

	// Your classic infrastructure API key.
	//
	// For information about viewing and accessing your classic infrastructure API key, see the
	// [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructurePassword *string `json:"classic_infrastructure_password,omitempty"`

	// An IBM Cloud API key that can create and manage service IDs. The API key must be assigned the Editor platform role
	// on the Access Groups Service and the Operator platform role on the IAM Identity Service.  For more information, see
	// the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).
	ApiKey *string `json:"api_key,omitempty"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA in seconds.
	MaxTtlSeconds *int64 `json:"max_ttl_seconds,omitempty"`

	// The time until the certificate revocation list (CRL) expires, in seconds.
	CrlExpirySeconds *int64 `json:"crl_expiry_seconds,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// he requested TTL, after which the certificate expires.
	TtlSeconds *int64 `json:"ttl_seconds,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The configuration data of your Private Certificate.
	Data PrivateCertificateCADataIntf `json:"data,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The signing method to use with this certificate authority to generate private certificates.
	//
	// You can choose between internal or externally signed options. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	SigningMethod *string `json:"signing_method,omitempty"`

	// The name of the intermediate certificate authority.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// This field scopes the creation of private certificates to only the secret groups that you specify.
	//
	// This field can be supplied as a comma-delimited list of secret group IDs.
	AllowedSecretGroups *string `json:"allowed_secret_groups,omitempty"`

	// This field indicates whether to allow `localhost` to be included as one of the requested common names.
	AllowLocalhost *bool `json:"allow_localhost,omitempty"`

	// The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and
	// `allow_subdomains` options.
	AllowedDomains []string `json:"allowed_domains,omitempty"`

	// This field indicates whether to allow the domains that are supplied in the `allowed_domains` field to contain access
	// control list (ACL) templates.
	AllowedDomainsTemplate *bool `json:"allowed_domains_template,omitempty"`

	// This field indicates whether to allow clients to request private certificates that match the value of the actual
	// domains on the final certificate.
	//
	// For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a
	// certificate that contains the name `example.com` as one of the DNS values on the final certificate.
	//
	// **Important:** In some scenarios, allowing bare domains can be considered a security risk.
	AllowBareDomains *bool `json:"allow_bare_domains,omitempty"`

	// This field indicates whether to allow clients to request private certificates with common names (CN) that are
	// subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.
	//
	// For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the
	// following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.
	//
	// **Note:** This field is redundant if you use the `allow_any_name` option.
	AllowSubdomains *bool `json:"allow_subdomains,omitempty"`

	// This field indicates whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are
	// specified in the `allowed_domains` field.
	//
	// If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	AllowGlobDomains *bool `json:"allow_glob_domains,omitempty"`

	// This field indicates whether to allow clients to request a private certificate that matches any common name.
	AllowAnyName *bool `json:"allow_any_name,omitempty"`

	// This field indicates whether to enforce only valid hostnames for common names, DNS Subject Alternative Names, and
	// the host section of email addresses.
	EnforceHostnames *bool `json:"enforce_hostnames,omitempty"`

	// This field indicates whether to allow clients to request a private certificate with IP Subject Alternative Names.
	AllowIpSans *bool `json:"allow_ip_sans,omitempty"`

	// The URI Subject Alternative Names to allow for private certificates.
	//
	// Values can contain glob patterns, for example `spiffe://hostname/_*`.
	AllowedUriSans []string `json:"allowed_uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private
	// certificates.
	//
	// The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type
	// is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any
	// `other_sans` input.
	AllowedOtherSans []string `json:"allowed_other_sans,omitempty"`

	// This field indicates whether private certificates are flagged for server use.
	ServerFlag *bool `json:"server_flag,omitempty"`

	// This field indicates whether private certificates are flagged for client use.
	ClientFlag *bool `json:"client_flag,omitempty"`

	// This field indicates whether private certificates are flagged for code signing use.
	CodeSigningFlag *bool `json:"code_signing_flag,omitempty"`

	// This field indicates whether private certificates are flagged for email protection use.
	EmailProtectionFlag *bool `json:"email_protection_flag,omitempty"`

	// The allowed key usage constraint to define for private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage). Omit the
	// `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to
	// an empty list.
	KeyUsage []string `json:"key_usage,omitempty"`

	// The allowed extended key usage constraint on private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage).
	// Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set
	// this field to an empty list.
	ExtKeyUsage []string `json:"ext_key_usage,omitempty"`

	// A list of extended key usage Object Identifiers (OIDs).
	ExtKeyUsageOids []string `json:"ext_key_usage_oids,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// common name (CN) from a certificate signing request (CSR) instead of the CN that is included in the data of the
	// certificate.
	//
	// Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include
	// the `use_csr_sans` property.
	UseCsrCommonName *bool `json:"use_csr_common_name,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// Subject Alternative Names
	// (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the data of the
	// certificate.
	//
	// This field does not include the common name in the CSR. To use the common name, include the `use_csr_common_name`
	// property.
	UseCsrSans *bool `json:"use_csr_sans,omitempty"`

	// This field indicates whether to require a common name to create a private certificate.
	//
	// By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the
	// `require_cn` option to `false`.
	RequireCn *bool `json:"require_cn,omitempty"`

	// A list of policy Object Identifiers (OIDs).
	PolicyIdentifiers []string `json:"policy_identifiers,omitempty"`

	// This field indicates whether to mark the Basic Constraints extension of an issued private certificate as valid for
	// non-CA certificates.
	BasicConstraintsValidForNonCa *bool `json:"basic_constraints_valid_for_non_ca,omitempty"`

	// The duration in seconds by which to backdate the `not_before` property of an issued private certificate.
	NotBeforeDurationSeconds *int64 `json:"not_before_duration_seconds,omitempty"`
}

// Constants associated with the Configuration.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	Configuration_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	Configuration_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	Configuration_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	Configuration_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	Configuration_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	Configuration_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	Configuration_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the Configuration.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	Configuration_SecretType_Arbitrary          = "arbitrary"
	Configuration_SecretType_IamCredentials     = "iam_credentials"
	Configuration_SecretType_ImportedCert       = "imported_cert"
	Configuration_SecretType_Kv                 = "kv"
	Configuration_SecretType_PrivateCert        = "private_cert"
	Configuration_SecretType_PublicCert         = "public_cert"
	Configuration_SecretType_ServiceCredentials = "service_credentials"
	Configuration_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the Configuration.LetsEncryptEnvironment property.
// The configuration of the Let's Encrypt CA environment.
const (
	Configuration_LetsEncryptEnvironment_Production = "production"
	Configuration_LetsEncryptEnvironment_Staging    = "staging"
)

// Constants associated with the Configuration.KeyType property.
// The type of private key to generate.
const (
	Configuration_KeyType_Ec  = "ec"
	Configuration_KeyType_Rsa = "rsa"
)

// Constants associated with the Configuration.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	Configuration_Status_CertificateTemplateRequired = "certificate_template_required"
	Configuration_Status_Configured                  = "configured"
	Configuration_Status_Expired                     = "expired"
	Configuration_Status_Revoked                     = "revoked"
	Configuration_Status_SignedCertificateRequired   = "signed_certificate_required"
	Configuration_Status_SigningRequired             = "signing_required"
)

// Constants associated with the Configuration.Format property.
// The format of the returned data.
const (
	Configuration_Format_Pem       = "pem"
	Configuration_Format_PemBundle = "pem_bundle"
)

// Constants associated with the Configuration.PrivateKeyFormat property.
// The format of the generated private key.
const (
	Configuration_PrivateKeyFormat_Der   = "der"
	Configuration_PrivateKeyFormat_Pkcs8 = "pkcs8"
)

// Constants associated with the Configuration.SigningMethod property.
// The signing method to use with this certificate authority to generate private certificates.
//
// You can choose between internal or externally signed options. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
const (
	Configuration_SigningMethod_External = "external"
	Configuration_SigningMethod_Internal = "internal"
)

func (*Configuration) isaConfiguration() bool {
	return true
}

type ConfigurationIntf interface {
	isaConfiguration() bool
}

// UnmarshalConfiguration unmarshals an instance of Configuration from the specified map of raw messages.
func UnmarshalConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "config_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'config_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'config_type' not found in JSON object")
		return
	}
	if discValue == "public_cert_configuration_ca_lets_encrypt" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationCALetsEncrypt)
	} else if discValue == "public_cert_configuration_dns_cloud_internet_services" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationDNSCloudInternetServices)
	} else if discValue == "public_cert_configuration_dns_classic_infrastructure" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationDNSClassicInfrastructure)
	} else if discValue == "iam_credentials_configuration" {
		err = core.UnmarshalModel(m, "", result, UnmarshalIAMCredentialsConfiguration)
	} else if discValue == "private_cert_configuration_root_ca" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationRootCA)
	} else if discValue == "private_cert_configuration_intermediate_ca" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationIntermediateCA)
	} else if discValue == "private_cert_configuration_template" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationTemplate)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'config_type': %s", discValue)
	}
	return
}

// ConfigurationAction : The response body to specify the properties of the action to create a configuration.
// Models which "extend" this model:
// - PrivateCertificateConfigurationActionRevoke
// - PrivateCertificateConfigurationActionSignCSR
// - PrivateCertificateConfigurationActionSignIntermediate
// - PrivateCertificateConfigurationActionSetSigned
// - PrivateCertificateConfigurationActionRotateCRL
type ConfigurationAction struct {
	// The type of configuration action.
	ActionType *string `json:"action_type,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '12h'. The value can't
	// exceed the `max_ttl` that is defined in the associated certificate template.
	TTL *string `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// This field indicates whether to use values from a certificate signing request (CSR) to complete a
	// `private_cert_configuration_action_sign_csr` action. If it is set to `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than by using the
	// values that are provided in the other parameters to this operation.
	//
	// 2) Any key usage, for example, non-repudiation, that is requested in the CSR are added to the basic set of key
	// usages used for CA certificates that are signed by the intermediate authority.
	//
	// 3) Extensions that are requested in the CSR are copied into the issued private certificate.
	UseCsrValues *bool `json:"use_csr_values,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The certificate signing request.
	Csr *string `json:"csr,omitempty"`

	// The data that is associated with the root certificate authority.
	Data *PrivateCertificateConfigurationCACertificate `json:"data,omitempty"`

	// The unique name of your configuration.
	IntermediateCertificateAuthority *string `json:"intermediate_certificate_authority,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// This field indicates whether the request to rotate the CRL for the private certificate configuration was successful.
	Success *bool `json:"success,omitempty"`
}

// Constants associated with the ConfigurationAction.ActionType property.
// The type of configuration action.
const (
	ConfigurationAction_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	ConfigurationAction_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	ConfigurationAction_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	ConfigurationAction_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	ConfigurationAction_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

// Constants associated with the ConfigurationAction.Format property.
// The format of the returned data.
const (
	ConfigurationAction_Format_Pem       = "pem"
	ConfigurationAction_Format_PemBundle = "pem_bundle"
)

func (*ConfigurationAction) isaConfigurationAction() bool {
	return true
}

type ConfigurationActionIntf interface {
	isaConfigurationAction() bool
}

// UnmarshalConfigurationAction unmarshals an instance of ConfigurationAction from the specified map of raw messages.
func UnmarshalConfigurationAction(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "action_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'action_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'action_type' not found in JSON object")
		return
	}
	if discValue == "private_cert_configuration_action_revoke_ca_certificate" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionRevoke)
	} else if discValue == "private_cert_configuration_action_sign_csr" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionSignCSR)
	} else if discValue == "private_cert_configuration_action_sign_intermediate" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionSignIntermediate)
	} else if discValue == "private_cert_configuration_action_set_signed" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionSetSigned)
	} else if discValue == "private_cert_configuration_action_rotate_crl" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionRotateCRL)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'action_type': %s", discValue)
	}
	return
}

// ConfigurationActionPrototype : The request body to specify the properties of the action to create a configuration.
// Models which "extend" this model:
// - PrivateCertificateConfigurationActionRotateCRLPrototype
// - PrivateCertificateConfigurationActionRevokePrototype
// - PrivateCertificateConfigurationActionSignCSRPrototype
// - PrivateCertificateConfigurationActionSignIntermediatePrototype
// - PrivateCertificateConfigurationActionSetSignedPrototype
type ConfigurationActionPrototype struct {
	// The type of configuration action.
	ActionType *string `json:"action_type,omitempty"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '12h'. The value can't
	// exceed the `max_ttl` that is defined in the associated certificate template.
	TTL *string `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// This field indicates whether to use values from a certificate signing request (CSR) to complete a
	// `private_cert_configuration_action_sign_csr` action. If it is set to `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than by using the
	// values that are provided in the other parameters to this operation.
	//
	// 2) Any key usage, for example, non-repudiation, that is requested in the CSR are added to the basic set of key
	// usages used for CA certificates that are signed by the intermediate authority.
	//
	// 3) Extensions that are requested in the CSR are copied into the issued private certificate.
	UseCsrValues *bool `json:"use_csr_values,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The certificate signing request.
	Csr *string `json:"csr,omitempty"`

	// The unique name of your configuration.
	IntermediateCertificateAuthority *string `json:"intermediate_certificate_authority,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`
}

// Constants associated with the ConfigurationActionPrototype.ActionType property.
// The type of configuration action.
const (
	ConfigurationActionPrototype_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	ConfigurationActionPrototype_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	ConfigurationActionPrototype_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	ConfigurationActionPrototype_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	ConfigurationActionPrototype_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

// Constants associated with the ConfigurationActionPrototype.Format property.
// The format of the returned data.
const (
	ConfigurationActionPrototype_Format_Pem       = "pem"
	ConfigurationActionPrototype_Format_PemBundle = "pem_bundle"
)

func (*ConfigurationActionPrototype) isaConfigurationActionPrototype() bool {
	return true
}

type ConfigurationActionPrototypeIntf interface {
	isaConfigurationActionPrototype() bool
}

// UnmarshalConfigurationActionPrototype unmarshals an instance of ConfigurationActionPrototype from the specified map of raw messages.
func UnmarshalConfigurationActionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "action_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'action_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'action_type' not found in JSON object")
		return
	}
	if discValue == "private_cert_configuration_action_rotate_crl" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionRotateCRLPrototype)
	} else if discValue == "private_cert_configuration_action_revoke_ca_certificate" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionRevokePrototype)
	} else if discValue == "private_cert_configuration_action_sign_csr" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionSignCSRPrototype)
	} else if discValue == "private_cert_configuration_action_sign_intermediate" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionSignIntermediatePrototype)
	} else if discValue == "private_cert_configuration_action_set_signed" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationActionSetSignedPrototype)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'action_type': %s", discValue)
	}
	return
}

// ConfigurationMetadata : Your configuration metadata properties.
// Models which "extend" this model:
// - IAMCredentialsConfigurationMetadata
// - PublicCertificateConfigurationCALetsEncryptMetadata
// - PublicCertificateConfigurationDNSCloudInternetServicesMetadata
// - PublicCertificateConfigurationDNSClassicInfrastructureMetadata
// - PrivateCertificateConfigurationRootCAMetadata
// - PrivateCertificateConfigurationIntermediateCAMetadata
// - PrivateCertificateConfigurationTemplateMetadata
type ConfigurationMetadata struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type,omitempty"`

	// The unique name of your configuration.
	Name *string `json:"name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The configuration of the Let's Encrypt CA environment.
	LetsEncryptEnvironment *string `json:"lets_encrypt_environment,omitempty"`

	// If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name. If
	// no match, the default offered chain will be used.
	LetsEncryptPreferredChain *string `json:"lets_encrypt_preferred_chain,omitempty"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The signing method to use with this certificate authority to generate private certificates.
	//
	// You can choose between internal or externally signed options. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	SigningMethod *string `json:"signing_method,omitempty"`

	// The name of the intermediate certificate authority.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`
}

// Constants associated with the ConfigurationMetadata.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	ConfigurationMetadata_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	ConfigurationMetadata_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	ConfigurationMetadata_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	ConfigurationMetadata_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	ConfigurationMetadata_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	ConfigurationMetadata_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	ConfigurationMetadata_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the ConfigurationMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ConfigurationMetadata_SecretType_Arbitrary          = "arbitrary"
	ConfigurationMetadata_SecretType_IamCredentials     = "iam_credentials"
	ConfigurationMetadata_SecretType_ImportedCert       = "imported_cert"
	ConfigurationMetadata_SecretType_Kv                 = "kv"
	ConfigurationMetadata_SecretType_PrivateCert        = "private_cert"
	ConfigurationMetadata_SecretType_PublicCert         = "public_cert"
	ConfigurationMetadata_SecretType_ServiceCredentials = "service_credentials"
	ConfigurationMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ConfigurationMetadata.LetsEncryptEnvironment property.
// The configuration of the Let's Encrypt CA environment.
const (
	ConfigurationMetadata_LetsEncryptEnvironment_Production = "production"
	ConfigurationMetadata_LetsEncryptEnvironment_Staging    = "staging"
)

// Constants associated with the ConfigurationMetadata.KeyType property.
// The type of private key to generate.
const (
	ConfigurationMetadata_KeyType_Ec  = "ec"
	ConfigurationMetadata_KeyType_Rsa = "rsa"
)

// Constants associated with the ConfigurationMetadata.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	ConfigurationMetadata_Status_CertificateTemplateRequired = "certificate_template_required"
	ConfigurationMetadata_Status_Configured                  = "configured"
	ConfigurationMetadata_Status_Expired                     = "expired"
	ConfigurationMetadata_Status_Revoked                     = "revoked"
	ConfigurationMetadata_Status_SignedCertificateRequired   = "signed_certificate_required"
	ConfigurationMetadata_Status_SigningRequired             = "signing_required"
)

// Constants associated with the ConfigurationMetadata.SigningMethod property.
// The signing method to use with this certificate authority to generate private certificates.
//
// You can choose between internal or externally signed options. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
const (
	ConfigurationMetadata_SigningMethod_External = "external"
	ConfigurationMetadata_SigningMethod_Internal = "internal"
)

func (*ConfigurationMetadata) isaConfigurationMetadata() bool {
	return true
}

type ConfigurationMetadataIntf interface {
	isaConfigurationMetadata() bool
}

// UnmarshalConfigurationMetadata unmarshals an instance of ConfigurationMetadata from the specified map of raw messages.
func UnmarshalConfigurationMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "config_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'config_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'config_type' not found in JSON object")
		return
	}
	if discValue == "iam_credentials_configuration" {
		err = core.UnmarshalModel(m, "", result, UnmarshalIAMCredentialsConfigurationMetadata)
	} else if discValue == "public_cert_configuration_ca_lets_encrypt" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationCALetsEncryptMetadata)
	} else if discValue == "public_cert_configuration_dns_cloud_internet_services" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationDNSCloudInternetServicesMetadata)
	} else if discValue == "public_cert_configuration_dns_classic_infrastructure" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationDNSClassicInfrastructureMetadata)
	} else if discValue == "private_cert_configuration_root_ca" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationRootCAMetadata)
	} else if discValue == "private_cert_configuration_intermediate_ca" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationIntermediateCAMetadata)
	} else if discValue == "private_cert_configuration_template" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationTemplateMetadata)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'config_type': %s", discValue)
	}
	return
}

// ConfigurationMetadataPaginatedCollection : Properties that describe a paginated collection of secret locks.
type ConfigurationMetadataPaginatedCollection struct {
	// The total number of resources in a collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of items that are retrieved in a collection.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of items that are skipped in a collection.
	Offset *int64 `json:"offset" validate:"required"`

	// A URL that points to the first page in a collection.
	First *PaginatedCollectionFirst `json:"first" validate:"required"`

	// A URL that points to the next page in a collection.
	Next *PaginatedCollectionNext `json:"next,omitempty"`

	// A URL that points to the previous page in a collection.
	Previous *PaginatedCollectionPrevious `json:"previous,omitempty"`

	// A URL that points to the last page in a collection.
	Last *PaginatedCollectionLast `json:"last" validate:"required"`

	// A collection of configuration metadata.
	Configurations []ConfigurationMetadataIntf `json:"configurations" validate:"required"`
}

// UnmarshalConfigurationMetadataPaginatedCollection unmarshals an instance of ConfigurationMetadataPaginatedCollection from the specified map of raw messages.
func UnmarshalConfigurationMetadataPaginatedCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigurationMetadataPaginatedCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
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
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedCollectionPrevious)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedCollectionLast)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "configurations", &obj.Configurations, UnmarshalConfigurationMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ConfigurationMetadataPaginatedCollection) GetNextOffset() (*int64, error) {
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

// ConfigurationPatch : Your configuration update data.
// Models which "extend" this model:
// - IAMCredentialsConfigurationPatch
// - PrivateCertificateConfigurationRootCAPatch
// - PrivateCertificateConfigurationIntermediateCAPatch
// - PrivateCertificateConfigurationTemplatePatch
// - PublicCertificateConfigurationCALetsEncryptPatch
// - PublicCertificateConfigurationDNSCloudInternetServicesPatch
// - PublicCertificateConfigurationDNSClassicInfrastructurePatch
type ConfigurationPatch struct {
	// An IBM Cloud API key that can create and manage service IDs. The API key must be assigned the Editor platform role
	// on the Access Groups Service and the Operator platform role on the IAM Identity Service.  For more information, see
	// the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).
	ApiKey *string `json:"api_key,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL *string `json:"max_ttl,omitempty"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry *string `json:"crl_expiry,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// This field scopes the creation of private certificates to only the secret groups that you specify.
	//
	// This field can be supplied as a comma-delimited list of secret group IDs.
	AllowedSecretGroups *string `json:"allowed_secret_groups,omitempty"`

	// The requested time-to-live (TTL) for certificates that are created by this CA. This field's value can't be longer
	// than the `max_ttl` limit.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	TTL *string `json:"ttl,omitempty"`

	// This field indicates whether to allow `localhost` to be included as one of the requested common names.
	AllowLocalhost *bool `json:"allow_localhost,omitempty"`

	// The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and
	// `allow_subdomains` options.
	AllowedDomains []string `json:"allowed_domains,omitempty"`

	// This field indicates whether to allow the domains that are supplied in the `allowed_domains` field to contain access
	// control list (ACL) templates.
	AllowedDomainsTemplate *bool `json:"allowed_domains_template,omitempty"`

	// This field indicates whether to allow clients to request private certificates that match the value of the actual
	// domains on the final certificate.
	//
	// For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a
	// certificate that contains the name `example.com` as one of the DNS values on the final certificate.
	//
	// **Important:** In some scenarios, allowing bare domains can be considered a security risk.
	AllowBareDomains *bool `json:"allow_bare_domains,omitempty"`

	// This field indicates whether to allow clients to request private certificates with common names (CN) that are
	// subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.
	//
	// For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the
	// following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.
	//
	// **Note:** This field is redundant if you use the `allow_any_name` option.
	AllowSubdomains *bool `json:"allow_subdomains,omitempty"`

	// This field indicates whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are
	// specified in the `allowed_domains` field.
	//
	// If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	AllowGlobDomains *bool `json:"allow_glob_domains,omitempty"`

	// This field indicates whether to allow clients to request a private certificate that matches any common name.
	AllowAnyName *bool `json:"allow_any_name,omitempty"`

	// This field indicates whether to enforce only valid hostnames for common names, DNS Subject Alternative Names, and
	// the host section of email addresses.
	EnforceHostnames *bool `json:"enforce_hostnames,omitempty"`

	// This field indicates whether to allow clients to request a private certificate with IP Subject Alternative Names.
	AllowIpSans *bool `json:"allow_ip_sans,omitempty"`

	// The URI Subject Alternative Names to allow for private certificates.
	//
	// Values can contain glob patterns, for example `spiffe://hostname/_*`.
	AllowedUriSans []string `json:"allowed_uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private
	// certificates.
	//
	// The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type
	// is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any
	// `other_sans` input.
	AllowedOtherSans []string `json:"allowed_other_sans,omitempty"`

	// This field indicates whether private certificates are flagged for server use.
	ServerFlag *bool `json:"server_flag,omitempty"`

	// This field indicates whether private certificates are flagged for client use.
	ClientFlag *bool `json:"client_flag,omitempty"`

	// This field indicates whether private certificates are flagged for code signing use.
	CodeSigningFlag *bool `json:"code_signing_flag,omitempty"`

	// This field indicates whether private certificates are flagged for email protection use.
	EmailProtectionFlag *bool `json:"email_protection_flag,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The allowed key usage constraint to define for private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage). Omit the
	// `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to
	// an empty list.
	KeyUsage []string `json:"key_usage,omitempty"`

	// The allowed extended key usage constraint on private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage).
	// Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set
	// this field to an empty list.
	ExtKeyUsage []string `json:"ext_key_usage,omitempty"`

	// A list of extended key usage Object Identifiers (OIDs).
	ExtKeyUsageOids []string `json:"ext_key_usage_oids,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// common name (CN) from a certificate signing request (CSR) instead of the CN that is included in the data of the
	// certificate.
	//
	// Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include
	// the `use_csr_sans` property.
	UseCsrCommonName *bool `json:"use_csr_common_name,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// Subject Alternative Names
	// (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the data of the
	// certificate.
	//
	// This field does not include the common name in the CSR. To use the common name, include the `use_csr_common_name`
	// property.
	UseCsrSans *bool `json:"use_csr_sans,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// This field is deprecated. You can ignore its value.
	SerialNumber *string `json:"serial_number,omitempty"`

	// This field indicates whether to require a common name to create a private certificate.
	//
	// By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the
	// `require_cn` option to `false`.
	RequireCn *bool `json:"require_cn,omitempty"`

	// A list of policy Object Identifiers (OIDs).
	PolicyIdentifiers []string `json:"policy_identifiers,omitempty"`

	// This field indicates whether to mark the Basic Constraints extension of an issued private certificate as valid for
	// non-CA certificates.
	BasicConstraintsValidForNonCa *bool `json:"basic_constraints_valid_for_non_ca,omitempty"`

	// The duration in seconds by which to backdate the `not_before` property of an issued private certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value
	// is returned in seconds (integer).
	NotBeforeDuration *string `json:"not_before_duration,omitempty"`

	// The configuration of the Let's Encrypt CA environment.
	LetsEncryptEnvironment *string `json:"lets_encrypt_environment,omitempty"`

	// The PEM-encoded private key of your Let's Encrypt account. The data must be formatted on a single line with embedded
	// newline characters.
	LetsEncryptPrivateKey *string `json:"lets_encrypt_private_key,omitempty"`

	// If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name. If
	// no match, the default offered chain will be used.
	LetsEncryptPreferredChain *string `json:"lets_encrypt_preferred_chain,omitempty"`

	// An IBM Cloud API key that can to list domains in your Cloud Internet Services instance.
	//
	// To grant Secrets Manager the ability to view the Cloud Internet Services instance and all of its domains, the API
	// key must be assigned the Reader service role on Internet Services (`internet-svcs`).
	//
	// If you need to manage specific domains, you can assign the Manager role. For production environments, it is
	// recommended that you assign the Reader access role, and then use the
	// [IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific
	// domains. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	CloudInternetServicesApikey *string `json:"cloud_internet_services_apikey,omitempty"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	CloudInternetServicesCrn *string `json:"cloud_internet_services_crn,omitempty"`

	// The username that is associated with your classic infrastructure account.
	//
	// In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information,
	// see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructureUsername *string `json:"classic_infrastructure_username,omitempty"`

	// Your classic infrastructure API key.
	//
	// For information about viewing and accessing your classic infrastructure API key, see the
	// [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructurePassword *string `json:"classic_infrastructure_password,omitempty"`
}

// Constants associated with the ConfigurationPatch.KeyType property.
// The type of private key to generate.
const (
	ConfigurationPatch_KeyType_Ec  = "ec"
	ConfigurationPatch_KeyType_Rsa = "rsa"
)

// Constants associated with the ConfigurationPatch.LetsEncryptEnvironment property.
// The configuration of the Let's Encrypt CA environment.
const (
	ConfigurationPatch_LetsEncryptEnvironment_Production = "production"
	ConfigurationPatch_LetsEncryptEnvironment_Staging    = "staging"
)

func (*ConfigurationPatch) isaConfigurationPatch() bool {
	return true
}

type ConfigurationPatchIntf interface {
	isaConfigurationPatch() bool
}

// UnmarshalConfigurationPatch unmarshals an instance of ConfigurationPatch from the specified map of raw messages.
func UnmarshalConfigurationPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigurationPatch)
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_ttl", &obj.MaxTTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_expiry", &obj.CrlExpiry)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_disable", &obj.CrlDisable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_certificates_urls_encoded", &obj.IssuingCertificatesUrlsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_secret_groups", &obj.AllowedSecretGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_localhost", &obj.AllowLocalhost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_domains", &obj.AllowedDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_domains_template", &obj.AllowedDomainsTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_bare_domains", &obj.AllowBareDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_subdomains", &obj.AllowSubdomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_glob_domains", &obj.AllowGlobDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_any_name", &obj.AllowAnyName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enforce_hostnames", &obj.EnforceHostnames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_ip_sans", &obj.AllowIpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_uri_sans", &obj.AllowedUriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_other_sans", &obj.AllowedOtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_flag", &obj.ServerFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_flag", &obj.ClientFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "code_signing_flag", &obj.CodeSigningFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email_protection_flag", &obj.EmailProtectionFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_usage", &obj.KeyUsage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ext_key_usage", &obj.ExtKeyUsage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ext_key_usage_oids", &obj.ExtKeyUsageOids)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_common_name", &obj.UseCsrCommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_sans", &obj.UseCsrSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "require_cn", &obj.RequireCn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_identifiers", &obj.PolicyIdentifiers)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "basic_constraints_valid_for_non_ca", &obj.BasicConstraintsValidForNonCa)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_before_duration", &obj.NotBeforeDuration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_environment", &obj.LetsEncryptEnvironment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_private_key", &obj.LetsEncryptPrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_preferred_chain", &obj.LetsEncryptPreferredChain)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_internet_services_apikey", &obj.CloudInternetServicesApikey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_internet_services_crn", &obj.CloudInternetServicesCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "classic_infrastructure_username", &obj.ClassicInfrastructureUsername)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "classic_infrastructure_password", &obj.ClassicInfrastructurePassword)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the ConfigurationPatch
func (configurationPatch *ConfigurationPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(configurationPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// ConfigurationPrototype : The details of your configuration.
// Models which "extend" this model:
// - PrivateCertificateConfigurationRootCAPrototype
// - PrivateCertificateConfigurationIntermediateCAPrototype
// - PrivateCertificateConfigurationTemplatePrototype
// - PublicCertificateConfigurationCALetsEncryptPrototype
// - PublicCertificateConfigurationDNSCloudInternetServicesPrototype
// - PublicCertificateConfigurationDNSClassicInfrastructurePrototype
// - IAMCredentialsConfigurationPrototype
type ConfigurationPrototype struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type,omitempty"`

	// A human-readable unique name to assign to your configuration.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.
	Name *string `json:"name,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL *string `json:"max_ttl,omitempty"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry *string `json:"crl_expiry,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The requested time-to-live (TTL) for certificates that are created by this CA. This field's value can't be longer
	// than the `max_ttl` limit.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	TTL *string `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The signing method to use with this certificate authority to generate private certificates.
	//
	// You can choose between internal or externally signed options. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	SigningMethod *string `json:"signing_method,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The name of the intermediate certificate authority.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// This field scopes the creation of private certificates to only the secret groups that you specify.
	//
	// This field can be supplied as a comma-delimited list of secret group IDs.
	AllowedSecretGroups *string `json:"allowed_secret_groups,omitempty"`

	// This field indicates whether to allow `localhost` to be included as one of the requested common names.
	AllowLocalhost *bool `json:"allow_localhost,omitempty"`

	// The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and
	// `allow_subdomains` options.
	AllowedDomains []string `json:"allowed_domains,omitempty"`

	// This field indicates whether to allow the domains that are supplied in the `allowed_domains` field to contain access
	// control list (ACL) templates.
	AllowedDomainsTemplate *bool `json:"allowed_domains_template,omitempty"`

	// This field indicates whether to allow clients to request private certificates that match the value of the actual
	// domains on the final certificate.
	//
	// For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a
	// certificate that contains the name `example.com` as one of the DNS values on the final certificate.
	//
	// **Important:** In some scenarios, allowing bare domains can be considered a security risk.
	AllowBareDomains *bool `json:"allow_bare_domains,omitempty"`

	// This field indicates whether to allow clients to request private certificates with common names (CN) that are
	// subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.
	//
	// For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the
	// following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.
	//
	// **Note:** This field is redundant if you use the `allow_any_name` option.
	AllowSubdomains *bool `json:"allow_subdomains,omitempty"`

	// This field indicates whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are
	// specified in the `allowed_domains` field.
	//
	// If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	AllowGlobDomains *bool `json:"allow_glob_domains,omitempty"`

	// This field indicates whether the issuance of certificates with RFC 6125 wildcards in the CN field.
	//
	// When set to false, this field prevents wildcards from being issued even if they can be allowed by an option
	// `allow_glob_domains`.
	AllowWildcardCertificates *bool `json:"allow_wildcard_certificates,omitempty"`

	// This field indicates whether to allow clients to request a private certificate that matches any common name.
	AllowAnyName *bool `json:"allow_any_name,omitempty"`

	// This field indicates whether to enforce only valid hostnames for common names, DNS Subject Alternative Names, and
	// the host section of email addresses.
	EnforceHostnames *bool `json:"enforce_hostnames,omitempty"`

	// This field indicates whether to allow clients to request a private certificate with IP Subject Alternative Names.
	AllowIpSans *bool `json:"allow_ip_sans,omitempty"`

	// The URI Subject Alternative Names to allow for private certificates.
	//
	// Values can contain glob patterns, for example `spiffe://hostname/_*`.
	AllowedUriSans []string `json:"allowed_uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private
	// certificates.
	//
	// The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type
	// is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any
	// `other_sans` input.
	AllowedOtherSans []string `json:"allowed_other_sans,omitempty"`

	// This field indicates whether private certificates are flagged for server use.
	ServerFlag *bool `json:"server_flag,omitempty"`

	// This field indicates whether private certificates are flagged for client use.
	ClientFlag *bool `json:"client_flag,omitempty"`

	// This field indicates whether private certificates are flagged for code signing use.
	CodeSigningFlag *bool `json:"code_signing_flag,omitempty"`

	// This field indicates whether private certificates are flagged for email protection use.
	EmailProtectionFlag *bool `json:"email_protection_flag,omitempty"`

	// The allowed key usage constraint to define for private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage). Omit the
	// `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to
	// an empty list.
	KeyUsage []string `json:"key_usage,omitempty"`

	// The allowed extended key usage constraint on private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage).
	// Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set
	// this field to an empty list.
	ExtKeyUsage []string `json:"ext_key_usage,omitempty"`

	// A list of extended key usage Object Identifiers (OIDs).
	ExtKeyUsageOids []string `json:"ext_key_usage_oids,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// common name (CN) from a certificate signing request (CSR) instead of the CN that is included in the data of the
	// certificate.
	//
	// Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include
	// the `use_csr_sans` property.
	UseCsrCommonName *bool `json:"use_csr_common_name,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// Subject Alternative Names
	// (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the data of the
	// certificate.
	//
	// This field does not include the common name in the CSR. To use the common name, include the `use_csr_common_name`
	// property.
	UseCsrSans *bool `json:"use_csr_sans,omitempty"`

	// This field indicates whether to require a common name to create a private certificate.
	//
	// By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the
	// `require_cn` option to `false`.
	RequireCn *bool `json:"require_cn,omitempty"`

	// A list of policy Object Identifiers (OIDs).
	PolicyIdentifiers []string `json:"policy_identifiers,omitempty"`

	// This field indicates whether to mark the Basic Constraints extension of an issued private certificate as valid for
	// non-CA certificates.
	BasicConstraintsValidForNonCa *bool `json:"basic_constraints_valid_for_non_ca,omitempty"`

	// The duration in seconds by which to backdate the `not_before` property of an issued private certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value
	// is returned in seconds (integer).
	NotBeforeDuration *string `json:"not_before_duration,omitempty"`

	// The configuration of the Let's Encrypt CA environment.
	LetsEncryptEnvironment *string `json:"lets_encrypt_environment,omitempty"`

	// The PEM-encoded private key of your Let's Encrypt account. The data must be formatted on a single line with embedded
	// newline characters.
	LetsEncryptPrivateKey *string `json:"lets_encrypt_private_key,omitempty"`

	// If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name. If
	// no match, the default offered chain will be used.
	LetsEncryptPreferredChain *string `json:"lets_encrypt_preferred_chain,omitempty"`

	// An IBM Cloud API key that can to list domains in your Cloud Internet Services instance.
	//
	// To grant Secrets Manager the ability to view the Cloud Internet Services instance and all of its domains, the API
	// key must be assigned the Reader service role on Internet Services (`internet-svcs`).
	//
	// If you need to manage specific domains, you can assign the Manager role. For production environments, it is
	// recommended that you assign the Reader access role, and then use the
	// [IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific
	// domains. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	CloudInternetServicesApikey *string `json:"cloud_internet_services_apikey,omitempty"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	CloudInternetServicesCrn *string `json:"cloud_internet_services_crn,omitempty"`

	// The username that is associated with your classic infrastructure account.
	//
	// In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information,
	// see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructureUsername *string `json:"classic_infrastructure_username,omitempty"`

	// Your classic infrastructure API key.
	//
	// For information about viewing and accessing your classic infrastructure API key, see the
	// [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructurePassword *string `json:"classic_infrastructure_password,omitempty"`

	// The API key that is used to set the iam_credentials engine.
	ApiKey *string `json:"api_key,omitempty"`
}

// Constants associated with the ConfigurationPrototype.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	ConfigurationPrototype_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	ConfigurationPrototype_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	ConfigurationPrototype_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	ConfigurationPrototype_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	ConfigurationPrototype_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	ConfigurationPrototype_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	ConfigurationPrototype_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the ConfigurationPrototype.Format property.
// The format of the returned data.
const (
	ConfigurationPrototype_Format_Pem       = "pem"
	ConfigurationPrototype_Format_PemBundle = "pem_bundle"
)

// Constants associated with the ConfigurationPrototype.PrivateKeyFormat property.
// The format of the generated private key.
const (
	ConfigurationPrototype_PrivateKeyFormat_Der   = "der"
	ConfigurationPrototype_PrivateKeyFormat_Pkcs8 = "pkcs8"
)

// Constants associated with the ConfigurationPrototype.KeyType property.
// The type of private key to generate.
const (
	ConfigurationPrototype_KeyType_Ec  = "ec"
	ConfigurationPrototype_KeyType_Rsa = "rsa"
)

// Constants associated with the ConfigurationPrototype.SigningMethod property.
// The signing method to use with this certificate authority to generate private certificates.
//
// You can choose between internal or externally signed options. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
const (
	ConfigurationPrototype_SigningMethod_External = "external"
	ConfigurationPrototype_SigningMethod_Internal = "internal"
)

// Constants associated with the ConfigurationPrototype.LetsEncryptEnvironment property.
// The configuration of the Let's Encrypt CA environment.
const (
	ConfigurationPrototype_LetsEncryptEnvironment_Production = "production"
	ConfigurationPrototype_LetsEncryptEnvironment_Staging    = "staging"
)

func (*ConfigurationPrototype) isaConfigurationPrototype() bool {
	return true
}

type ConfigurationPrototypeIntf interface {
	isaConfigurationPrototype() bool
}

// UnmarshalConfigurationPrototype unmarshals an instance of ConfigurationPrototype from the specified map of raw messages.
func UnmarshalConfigurationPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "config_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'config_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'config_type' not found in JSON object")
		return
	}
	if discValue == "private_cert_configuration_root_ca" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationRootCAPrototype)
	} else if discValue == "private_cert_configuration_intermediate_ca" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationIntermediateCAPrototype)
	} else if discValue == "private_cert_configuration_template" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateConfigurationTemplatePrototype)
	} else if discValue == "public_cert_configuration_ca_lets_encrypt" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationCALetsEncryptPrototype)
	} else if discValue == "public_cert_configuration_dns_cloud_internet_services" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationDNSCloudInternetServicesPrototype)
	} else if discValue == "public_cert_configuration_dns_classic_infrastructure" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateConfigurationDNSClassicInfrastructurePrototype)
	} else if discValue == "iam_credentials_configuration" {
		err = core.UnmarshalModel(m, "", result, UnmarshalIAMCredentialsConfigurationPrototype)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'config_type': %s", discValue)
	}
	return
}

// CosHmacKeys : The Cloud Object Storage HMAC keys that are returned after you create a service credentials secret.
type CosHmacKeys struct {
	// The access key ID for Cloud Object Storage HMAC credentials.
	AccessKeyID *string `json:"access_key_id,omitempty"`

	// The secret access key ID for Cloud Object Storage HMAC credentials.
	SecretAccessKey *string `json:"secret_access_key,omitempty"`
}

// UnmarshalCosHmacKeys unmarshals an instance of CosHmacKeys from the specified map of raw messages.
func UnmarshalCosHmacKeys(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CosHmacKeys)
	err = core.UnmarshalPrimitive(m, "access_key_id", &obj.AccessKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_access_key", &obj.SecretAccessKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateConfigurationActionOptions : The CreateConfigurationAction options.
type CreateConfigurationActionOptions struct {
	// The name that uniquely identifies a configuration.
	Name *string `json:"name" validate:"required,ne="`

	// The request body to specify the properties of the action to create a configuration.
	ConfigActionPrototype ConfigurationActionPrototypeIntf `json:"ConfigActionPrototype" validate:"required"`

	// The configuration type of this configuration - use this header to resolve 300 error responses.
	XSmAcceptConfigurationType *string `json:"X-Sm-Accept-Configuration-Type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateConfigurationActionOptions.XSmAcceptConfigurationType property.
// The configuration type of this configuration - use this header to resolve 300 error responses.
const (
	CreateConfigurationActionOptions_XSmAcceptConfigurationType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	CreateConfigurationActionOptions_XSmAcceptConfigurationType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	CreateConfigurationActionOptions_XSmAcceptConfigurationType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	CreateConfigurationActionOptions_XSmAcceptConfigurationType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	CreateConfigurationActionOptions_XSmAcceptConfigurationType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	CreateConfigurationActionOptions_XSmAcceptConfigurationType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	CreateConfigurationActionOptions_XSmAcceptConfigurationType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// NewCreateConfigurationActionOptions : Instantiate CreateConfigurationActionOptions
func (*SecretsManagerV2) NewCreateConfigurationActionOptions(name string, configActionPrototype ConfigurationActionPrototypeIntf) *CreateConfigurationActionOptions {
	return &CreateConfigurationActionOptions{
		Name:                  core.StringPtr(name),
		ConfigActionPrototype: configActionPrototype,
	}
}

// SetName : Allow user to set Name
func (_options *CreateConfigurationActionOptions) SetName(name string) *CreateConfigurationActionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetConfigActionPrototype : Allow user to set ConfigActionPrototype
func (_options *CreateConfigurationActionOptions) SetConfigActionPrototype(configActionPrototype ConfigurationActionPrototypeIntf) *CreateConfigurationActionOptions {
	_options.ConfigActionPrototype = configActionPrototype
	return _options
}

// SetXSmAcceptConfigurationType : Allow user to set XSmAcceptConfigurationType
func (_options *CreateConfigurationActionOptions) SetXSmAcceptConfigurationType(xSmAcceptConfigurationType string) *CreateConfigurationActionOptions {
	_options.XSmAcceptConfigurationType = core.StringPtr(xSmAcceptConfigurationType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateConfigurationActionOptions) SetHeaders(param map[string]string) *CreateConfigurationActionOptions {
	options.Headers = param
	return options
}

// CreateConfigurationOptions : The CreateConfiguration options.
type CreateConfigurationOptions struct {
	// The details of your configuration.
	ConfigurationPrototype ConfigurationPrototypeIntf `json:"ConfigurationPrototype" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateConfigurationOptions : Instantiate CreateConfigurationOptions
func (*SecretsManagerV2) NewCreateConfigurationOptions(configurationPrototype ConfigurationPrototypeIntf) *CreateConfigurationOptions {
	return &CreateConfigurationOptions{
		ConfigurationPrototype: configurationPrototype,
	}
}

// SetConfigurationPrototype : Allow user to set ConfigurationPrototype
func (_options *CreateConfigurationOptions) SetConfigurationPrototype(configurationPrototype ConfigurationPrototypeIntf) *CreateConfigurationOptions {
	_options.ConfigurationPrototype = configurationPrototype
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateConfigurationOptions) SetHeaders(param map[string]string) *CreateConfigurationOptions {
	options.Headers = param
	return options
}

// CreateNotificationsRegistrationOptions : The CreateNotificationsRegistration options.
type CreateNotificationsRegistrationOptions struct {
	// A CRN that uniquely identifies an IBM Cloud resource.
	EventNotificationsInstanceCrn *string `json:"event_notifications_instance_crn" validate:"required"`

	// The name that is displayed as a source that is in your Event Notifications instance.
	EventNotificationsSourceName *string `json:"event_notifications_source_name" validate:"required"`

	// An optional description for the source that is in your Event Notifications instance.
	EventNotificationsSourceDescription *string `json:"event_notifications_source_description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateNotificationsRegistrationOptions : Instantiate CreateNotificationsRegistrationOptions
func (*SecretsManagerV2) NewCreateNotificationsRegistrationOptions(eventNotificationsInstanceCrn string, eventNotificationsSourceName string) *CreateNotificationsRegistrationOptions {
	return &CreateNotificationsRegistrationOptions{
		EventNotificationsInstanceCrn: core.StringPtr(eventNotificationsInstanceCrn),
		EventNotificationsSourceName:  core.StringPtr(eventNotificationsSourceName),
	}
}

// SetEventNotificationsInstanceCrn : Allow user to set EventNotificationsInstanceCrn
func (_options *CreateNotificationsRegistrationOptions) SetEventNotificationsInstanceCrn(eventNotificationsInstanceCrn string) *CreateNotificationsRegistrationOptions {
	_options.EventNotificationsInstanceCrn = core.StringPtr(eventNotificationsInstanceCrn)
	return _options
}

// SetEventNotificationsSourceName : Allow user to set EventNotificationsSourceName
func (_options *CreateNotificationsRegistrationOptions) SetEventNotificationsSourceName(eventNotificationsSourceName string) *CreateNotificationsRegistrationOptions {
	_options.EventNotificationsSourceName = core.StringPtr(eventNotificationsSourceName)
	return _options
}

// SetEventNotificationsSourceDescription : Allow user to set EventNotificationsSourceDescription
func (_options *CreateNotificationsRegistrationOptions) SetEventNotificationsSourceDescription(eventNotificationsSourceDescription string) *CreateNotificationsRegistrationOptions {
	_options.EventNotificationsSourceDescription = core.StringPtr(eventNotificationsSourceDescription)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateNotificationsRegistrationOptions) SetHeaders(param map[string]string) *CreateNotificationsRegistrationOptions {
	options.Headers = param
	return options
}

// CreateSecretActionOptions : The CreateSecretAction options.
type CreateSecretActionOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	ID *string `json:"id" validate:"required,ne="`

	// The request body to specify the properties for your secret action.
	SecretActionPrototype SecretActionPrototypeIntf `json:"SecretActionPrototype" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSecretActionOptions : Instantiate CreateSecretActionOptions
func (*SecretsManagerV2) NewCreateSecretActionOptions(id string, secretActionPrototype SecretActionPrototypeIntf) *CreateSecretActionOptions {
	return &CreateSecretActionOptions{
		ID:                    core.StringPtr(id),
		SecretActionPrototype: secretActionPrototype,
	}
}

// SetID : Allow user to set ID
func (_options *CreateSecretActionOptions) SetID(id string) *CreateSecretActionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetSecretActionPrototype : Allow user to set SecretActionPrototype
func (_options *CreateSecretActionOptions) SetSecretActionPrototype(secretActionPrototype SecretActionPrototypeIntf) *CreateSecretActionOptions {
	_options.SecretActionPrototype = secretActionPrototype
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretActionOptions) SetHeaders(param map[string]string) *CreateSecretActionOptions {
	options.Headers = param
	return options
}

// CreateSecretGroupOptions : The CreateSecretGroup options.
type CreateSecretGroupOptions struct {
	// The name of your secret group.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret group.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSecretGroupOptions : Instantiate CreateSecretGroupOptions
func (*SecretsManagerV2) NewCreateSecretGroupOptions(name string) *CreateSecretGroupOptions {
	return &CreateSecretGroupOptions{
		Name: core.StringPtr(name),
	}
}

// SetName : Allow user to set Name
func (_options *CreateSecretGroupOptions) SetName(name string) *CreateSecretGroupOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateSecretGroupOptions) SetDescription(description string) *CreateSecretGroupOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretGroupOptions) SetHeaders(param map[string]string) *CreateSecretGroupOptions {
	options.Headers = param
	return options
}

// CreateSecretLocksBulkOptions : The CreateSecretLocksBulk options.
type CreateSecretLocksBulkOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	ID *string `json:"id" validate:"required,ne="`

	// The locks data to be attached to a secret version.
	Locks []SecretLockPrototype `json:"locks" validate:"required"`

	// An optional lock mode. When you create a lock, you can set one of the following modes to clear any matching locks on
	// a secret version.
	// - `remove_previous`: Removes any other locks with matching names if they are found in the previous version of the
	// secret. - `remove_previous_and_delete`: Completes the same action as `remove_previous`, but also permanently deletes
	// the data of the previous secret version if it doesn't have any locks.
	Mode *string `json:"mode,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateSecretLocksBulkOptions.Mode property.
// An optional lock mode. When you create a lock, you can set one of the following modes to clear any matching locks on
// a secret version.
// - `remove_previous`: Removes any other locks with matching names if they are found in the previous version of the
// secret. - `remove_previous_and_delete`: Completes the same action as `remove_previous`, but also permanently deletes
// the data of the previous secret version if it doesn't have any locks.
const (
	CreateSecretLocksBulkOptions_Mode_RemovePrevious          = "remove_previous"
	CreateSecretLocksBulkOptions_Mode_RemovePreviousAndDelete = "remove_previous_and_delete"
)

// NewCreateSecretLocksBulkOptions : Instantiate CreateSecretLocksBulkOptions
func (*SecretsManagerV2) NewCreateSecretLocksBulkOptions(id string, locks []SecretLockPrototype) *CreateSecretLocksBulkOptions {
	return &CreateSecretLocksBulkOptions{
		ID:    core.StringPtr(id),
		Locks: locks,
	}
}

// SetID : Allow user to set ID
func (_options *CreateSecretLocksBulkOptions) SetID(id string) *CreateSecretLocksBulkOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLocks : Allow user to set Locks
func (_options *CreateSecretLocksBulkOptions) SetLocks(locks []SecretLockPrototype) *CreateSecretLocksBulkOptions {
	_options.Locks = locks
	return _options
}

// SetMode : Allow user to set Mode
func (_options *CreateSecretLocksBulkOptions) SetMode(mode string) *CreateSecretLocksBulkOptions {
	_options.Mode = core.StringPtr(mode)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretLocksBulkOptions) SetHeaders(param map[string]string) *CreateSecretLocksBulkOptions {
	options.Headers = param
	return options
}

// CreateSecretOptions : The CreateSecret options.
type CreateSecretOptions struct {
	// Specify the properties for your secret.
	SecretPrototype SecretPrototypeIntf `json:"SecretPrototype" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSecretOptions : Instantiate CreateSecretOptions
func (*SecretsManagerV2) NewCreateSecretOptions(secretPrototype SecretPrototypeIntf) *CreateSecretOptions {
	return &CreateSecretOptions{
		SecretPrototype: secretPrototype,
	}
}

// SetSecretPrototype : Allow user to set SecretPrototype
func (_options *CreateSecretOptions) SetSecretPrototype(secretPrototype SecretPrototypeIntf) *CreateSecretOptions {
	_options.SecretPrototype = secretPrototype
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretOptions) SetHeaders(param map[string]string) *CreateSecretOptions {
	options.Headers = param
	return options
}

// CreateSecretVersionActionOptions : The CreateSecretVersionAction options.
type CreateSecretVersionActionOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies your secret version. You can use the `current` or `previous` aliases to refer
	// to the current or previous secret version.
	ID *string `json:"id" validate:"required,ne="`

	// The request body to specify the properties of the action to create a secret version.
	SecretVersionActionPrototype SecretVersionActionPrototypeIntf `json:"SecretVersionActionPrototype" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSecretVersionActionOptions : Instantiate CreateSecretVersionActionOptions
func (*SecretsManagerV2) NewCreateSecretVersionActionOptions(secretID string, id string, secretVersionActionPrototype SecretVersionActionPrototypeIntf) *CreateSecretVersionActionOptions {
	return &CreateSecretVersionActionOptions{
		SecretID:                     core.StringPtr(secretID),
		ID:                           core.StringPtr(id),
		SecretVersionActionPrototype: secretVersionActionPrototype,
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *CreateSecretVersionActionOptions) SetSecretID(secretID string) *CreateSecretVersionActionOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateSecretVersionActionOptions) SetID(id string) *CreateSecretVersionActionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetSecretVersionActionPrototype : Allow user to set SecretVersionActionPrototype
func (_options *CreateSecretVersionActionOptions) SetSecretVersionActionPrototype(secretVersionActionPrototype SecretVersionActionPrototypeIntf) *CreateSecretVersionActionOptions {
	_options.SecretVersionActionPrototype = secretVersionActionPrototype
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretVersionActionOptions) SetHeaders(param map[string]string) *CreateSecretVersionActionOptions {
	options.Headers = param
	return options
}

// CreateSecretVersionLocksBulkOptions : The CreateSecretVersionLocksBulk options.
type CreateSecretVersionLocksBulkOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies your secret version. You can use the `current` or `previous` aliases to refer
	// to the current or previous secret version.
	ID *string `json:"id" validate:"required,ne="`

	// The locks data to be attached to a secret version.
	Locks []SecretLockPrototype `json:"locks" validate:"required"`

	// An optional lock mode. When you create a lock, you can set one of the following modes to clear any matching locks on
	// a secret version.
	// - `remove_previous`: Removes any other locks with matching names if they are found in the previous version of the
	// secret. - `remove_previous_and_delete`: Completes the same action as `remove_previous`, but also permanently deletes
	// the data of the previous secret version if it doesn't have any locks.
	Mode *string `json:"mode,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateSecretVersionLocksBulkOptions.Mode property.
// An optional lock mode. When you create a lock, you can set one of the following modes to clear any matching locks on
// a secret version.
// - `remove_previous`: Removes any other locks with matching names if they are found in the previous version of the
// secret. - `remove_previous_and_delete`: Completes the same action as `remove_previous`, but also permanently deletes
// the data of the previous secret version if it doesn't have any locks.
const (
	CreateSecretVersionLocksBulkOptions_Mode_RemovePrevious          = "remove_previous"
	CreateSecretVersionLocksBulkOptions_Mode_RemovePreviousAndDelete = "remove_previous_and_delete"
)

// NewCreateSecretVersionLocksBulkOptions : Instantiate CreateSecretVersionLocksBulkOptions
func (*SecretsManagerV2) NewCreateSecretVersionLocksBulkOptions(secretID string, id string, locks []SecretLockPrototype) *CreateSecretVersionLocksBulkOptions {
	return &CreateSecretVersionLocksBulkOptions{
		SecretID: core.StringPtr(secretID),
		ID:       core.StringPtr(id),
		Locks:    locks,
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *CreateSecretVersionLocksBulkOptions) SetSecretID(secretID string) *CreateSecretVersionLocksBulkOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateSecretVersionLocksBulkOptions) SetID(id string) *CreateSecretVersionLocksBulkOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLocks : Allow user to set Locks
func (_options *CreateSecretVersionLocksBulkOptions) SetLocks(locks []SecretLockPrototype) *CreateSecretVersionLocksBulkOptions {
	_options.Locks = locks
	return _options
}

// SetMode : Allow user to set Mode
func (_options *CreateSecretVersionLocksBulkOptions) SetMode(mode string) *CreateSecretVersionLocksBulkOptions {
	_options.Mode = core.StringPtr(mode)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretVersionLocksBulkOptions) SetHeaders(param map[string]string) *CreateSecretVersionLocksBulkOptions {
	options.Headers = param
	return options
}

// CreateSecretVersionOptions : The CreateSecretVersion options.
type CreateSecretVersionOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// Specify the properties for your new secret version.
	SecretVersionPrototype SecretVersionPrototypeIntf `json:"SecretVersionPrototype" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSecretVersionOptions : Instantiate CreateSecretVersionOptions
func (*SecretsManagerV2) NewCreateSecretVersionOptions(secretID string, secretVersionPrototype SecretVersionPrototypeIntf) *CreateSecretVersionOptions {
	return &CreateSecretVersionOptions{
		SecretID:               core.StringPtr(secretID),
		SecretVersionPrototype: secretVersionPrototype,
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *CreateSecretVersionOptions) SetSecretID(secretID string) *CreateSecretVersionOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetSecretVersionPrototype : Allow user to set SecretVersionPrototype
func (_options *CreateSecretVersionOptions) SetSecretVersionPrototype(secretVersionPrototype SecretVersionPrototypeIntf) *CreateSecretVersionOptions {
	_options.SecretVersionPrototype = secretVersionPrototype
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretVersionOptions) SetHeaders(param map[string]string) *CreateSecretVersionOptions {
	options.Headers = param
	return options
}

// DeleteConfigurationOptions : The DeleteConfiguration options.
type DeleteConfigurationOptions struct {
	// The name that uniquely identifies a configuration.
	Name *string `json:"name" validate:"required,ne="`

	// The configuration type of this configuration - use this header to resolve 300 error responses.
	XSmAcceptConfigurationType *string `json:"X-Sm-Accept-Configuration-Type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteConfigurationOptions.XSmAcceptConfigurationType property.
// The configuration type of this configuration - use this header to resolve 300 error responses.
const (
	DeleteConfigurationOptions_XSmAcceptConfigurationType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	DeleteConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	DeleteConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	DeleteConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	DeleteConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	DeleteConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	DeleteConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// NewDeleteConfigurationOptions : Instantiate DeleteConfigurationOptions
func (*SecretsManagerV2) NewDeleteConfigurationOptions(name string) *DeleteConfigurationOptions {
	return &DeleteConfigurationOptions{
		Name: core.StringPtr(name),
	}
}

// SetName : Allow user to set Name
func (_options *DeleteConfigurationOptions) SetName(name string) *DeleteConfigurationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetXSmAcceptConfigurationType : Allow user to set XSmAcceptConfigurationType
func (_options *DeleteConfigurationOptions) SetXSmAcceptConfigurationType(xSmAcceptConfigurationType string) *DeleteConfigurationOptions {
	_options.XSmAcceptConfigurationType = core.StringPtr(xSmAcceptConfigurationType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteConfigurationOptions) SetHeaders(param map[string]string) *DeleteConfigurationOptions {
	options.Headers = param
	return options
}

// DeleteNotificationsRegistrationOptions : The DeleteNotificationsRegistration options.
type DeleteNotificationsRegistrationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteNotificationsRegistrationOptions : Instantiate DeleteNotificationsRegistrationOptions
func (*SecretsManagerV2) NewDeleteNotificationsRegistrationOptions() *DeleteNotificationsRegistrationOptions {
	return &DeleteNotificationsRegistrationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *DeleteNotificationsRegistrationOptions) SetHeaders(param map[string]string) *DeleteNotificationsRegistrationOptions {
	options.Headers = param
	return options
}

// DeleteSecretGroupOptions : The DeleteSecretGroup options.
type DeleteSecretGroupOptions struct {
	// The v4 UUID that uniquely identifies your secret group.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSecretGroupOptions : Instantiate DeleteSecretGroupOptions
func (*SecretsManagerV2) NewDeleteSecretGroupOptions(id string) *DeleteSecretGroupOptions {
	return &DeleteSecretGroupOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteSecretGroupOptions) SetID(id string) *DeleteSecretGroupOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecretGroupOptions) SetHeaders(param map[string]string) *DeleteSecretGroupOptions {
	options.Headers = param
	return options
}

// DeleteSecretLocksBulkOptions : The DeleteSecretLocksBulk options.
type DeleteSecretLocksBulkOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	ID *string `json:"id" validate:"required,ne="`

	// Specify the names of the secret locks to be deleted.
	Name []string `json:"name,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSecretLocksBulkOptions : Instantiate DeleteSecretLocksBulkOptions
func (*SecretsManagerV2) NewDeleteSecretLocksBulkOptions(id string) *DeleteSecretLocksBulkOptions {
	return &DeleteSecretLocksBulkOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteSecretLocksBulkOptions) SetID(id string) *DeleteSecretLocksBulkOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteSecretLocksBulkOptions) SetName(name []string) *DeleteSecretLocksBulkOptions {
	_options.Name = name
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecretLocksBulkOptions) SetHeaders(param map[string]string) *DeleteSecretLocksBulkOptions {
	options.Headers = param
	return options
}

// DeleteSecretOptions : The DeleteSecret options.
type DeleteSecretOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSecretOptions : Instantiate DeleteSecretOptions
func (*SecretsManagerV2) NewDeleteSecretOptions(id string) *DeleteSecretOptions {
	return &DeleteSecretOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteSecretOptions) SetID(id string) *DeleteSecretOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecretOptions) SetHeaders(param map[string]string) *DeleteSecretOptions {
	options.Headers = param
	return options
}

// DeleteSecretVersionDataOptions : The DeleteSecretVersionData options.
type DeleteSecretVersionDataOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies your secret version. You can use the `current` or `previous` aliases to refer
	// to the current or previous secret version.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSecretVersionDataOptions : Instantiate DeleteSecretVersionDataOptions
func (*SecretsManagerV2) NewDeleteSecretVersionDataOptions(secretID string, id string) *DeleteSecretVersionDataOptions {
	return &DeleteSecretVersionDataOptions{
		SecretID: core.StringPtr(secretID),
		ID:       core.StringPtr(id),
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *DeleteSecretVersionDataOptions) SetSecretID(secretID string) *DeleteSecretVersionDataOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteSecretVersionDataOptions) SetID(id string) *DeleteSecretVersionDataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecretVersionDataOptions) SetHeaders(param map[string]string) *DeleteSecretVersionDataOptions {
	options.Headers = param
	return options
}

// DeleteSecretVersionLocksBulkOptions : The DeleteSecretVersionLocksBulk options.
type DeleteSecretVersionLocksBulkOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies your secret version. You can use the `current` or `previous` aliases to refer
	// to the current or previous secret version.
	ID *string `json:"id" validate:"required,ne="`

	// Specify the names of the secret locks to be deleted.
	Name []string `json:"name,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSecretVersionLocksBulkOptions : Instantiate DeleteSecretVersionLocksBulkOptions
func (*SecretsManagerV2) NewDeleteSecretVersionLocksBulkOptions(secretID string, id string) *DeleteSecretVersionLocksBulkOptions {
	return &DeleteSecretVersionLocksBulkOptions{
		SecretID: core.StringPtr(secretID),
		ID:       core.StringPtr(id),
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *DeleteSecretVersionLocksBulkOptions) SetSecretID(secretID string) *DeleteSecretVersionLocksBulkOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteSecretVersionLocksBulkOptions) SetID(id string) *DeleteSecretVersionLocksBulkOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteSecretVersionLocksBulkOptions) SetName(name []string) *DeleteSecretVersionLocksBulkOptions {
	_options.Name = name
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecretVersionLocksBulkOptions) SetHeaders(param map[string]string) *DeleteSecretVersionLocksBulkOptions {
	options.Headers = param
	return options
}

// GetConfigurationOptions : The GetConfiguration options.
type GetConfigurationOptions struct {
	// The name that uniquely identifies a configuration.
	Name *string `json:"name" validate:"required,ne="`

	// The configuration type of this configuration - use this header to resolve 300 error responses.
	XSmAcceptConfigurationType *string `json:"X-Sm-Accept-Configuration-Type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetConfigurationOptions.XSmAcceptConfigurationType property.
// The configuration type of this configuration - use this header to resolve 300 error responses.
const (
	GetConfigurationOptions_XSmAcceptConfigurationType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	GetConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	GetConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	GetConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	GetConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	GetConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	GetConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// NewGetConfigurationOptions : Instantiate GetConfigurationOptions
func (*SecretsManagerV2) NewGetConfigurationOptions(name string) *GetConfigurationOptions {
	return &GetConfigurationOptions{
		Name: core.StringPtr(name),
	}
}

// SetName : Allow user to set Name
func (_options *GetConfigurationOptions) SetName(name string) *GetConfigurationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetXSmAcceptConfigurationType : Allow user to set XSmAcceptConfigurationType
func (_options *GetConfigurationOptions) SetXSmAcceptConfigurationType(xSmAcceptConfigurationType string) *GetConfigurationOptions {
	_options.XSmAcceptConfigurationType = core.StringPtr(xSmAcceptConfigurationType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigurationOptions) SetHeaders(param map[string]string) *GetConfigurationOptions {
	options.Headers = param
	return options
}

// GetNotificationsRegistrationOptions : The GetNotificationsRegistration options.
type GetNotificationsRegistrationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetNotificationsRegistrationOptions : Instantiate GetNotificationsRegistrationOptions
func (*SecretsManagerV2) NewGetNotificationsRegistrationOptions() *GetNotificationsRegistrationOptions {
	return &GetNotificationsRegistrationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetNotificationsRegistrationOptions) SetHeaders(param map[string]string) *GetNotificationsRegistrationOptions {
	options.Headers = param
	return options
}

// GetNotificationsRegistrationTestOptions : The GetNotificationsRegistrationTest options.
type GetNotificationsRegistrationTestOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetNotificationsRegistrationTestOptions : Instantiate GetNotificationsRegistrationTestOptions
func (*SecretsManagerV2) NewGetNotificationsRegistrationTestOptions() *GetNotificationsRegistrationTestOptions {
	return &GetNotificationsRegistrationTestOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetNotificationsRegistrationTestOptions) SetHeaders(param map[string]string) *GetNotificationsRegistrationTestOptions {
	options.Headers = param
	return options
}

// GetSecretByNameTypeOptions : The GetSecretByNameType options.
type GetSecretByNameTypeOptions struct {
	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// A human-readable name to assign to your secret. To protect your privacy, do not use personal data, such as your name
	// or location, as a name for your secret.
	Name *string `json:"name" validate:"required,ne="`

	// The name of your secret group.
	SecretGroupName *string `json:"secret_group_name" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSecretByNameTypeOptions.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	GetSecretByNameTypeOptions_SecretType_Arbitrary          = "arbitrary"
	GetSecretByNameTypeOptions_SecretType_IamCredentials     = "iam_credentials"
	GetSecretByNameTypeOptions_SecretType_ImportedCert       = "imported_cert"
	GetSecretByNameTypeOptions_SecretType_Kv                 = "kv"
	GetSecretByNameTypeOptions_SecretType_PrivateCert        = "private_cert"
	GetSecretByNameTypeOptions_SecretType_PublicCert         = "public_cert"
	GetSecretByNameTypeOptions_SecretType_ServiceCredentials = "service_credentials"
	GetSecretByNameTypeOptions_SecretType_UsernamePassword   = "username_password"
)

// NewGetSecretByNameTypeOptions : Instantiate GetSecretByNameTypeOptions
func (*SecretsManagerV2) NewGetSecretByNameTypeOptions(secretType string, name string, secretGroupName string) *GetSecretByNameTypeOptions {
	return &GetSecretByNameTypeOptions{
		SecretType:      core.StringPtr(secretType),
		Name:            core.StringPtr(name),
		SecretGroupName: core.StringPtr(secretGroupName),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetSecretByNameTypeOptions) SetSecretType(secretType string) *GetSecretByNameTypeOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetSecretByNameTypeOptions) SetName(name string) *GetSecretByNameTypeOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetSecretGroupName : Allow user to set SecretGroupName
func (_options *GetSecretByNameTypeOptions) SetSecretGroupName(secretGroupName string) *GetSecretByNameTypeOptions {
	_options.SecretGroupName = core.StringPtr(secretGroupName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretByNameTypeOptions) SetHeaders(param map[string]string) *GetSecretByNameTypeOptions {
	options.Headers = param
	return options
}

// GetSecretGroupOptions : The GetSecretGroup options.
type GetSecretGroupOptions struct {
	// The v4 UUID that uniquely identifies your secret group.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSecretGroupOptions : Instantiate GetSecretGroupOptions
func (*SecretsManagerV2) NewGetSecretGroupOptions(id string) *GetSecretGroupOptions {
	return &GetSecretGroupOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetSecretGroupOptions) SetID(id string) *GetSecretGroupOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretGroupOptions) SetHeaders(param map[string]string) *GetSecretGroupOptions {
	options.Headers = param
	return options
}

// GetSecretMetadataOptions : The GetSecretMetadata options.
type GetSecretMetadataOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSecretMetadataOptions : Instantiate GetSecretMetadataOptions
func (*SecretsManagerV2) NewGetSecretMetadataOptions(id string) *GetSecretMetadataOptions {
	return &GetSecretMetadataOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetSecretMetadataOptions) SetID(id string) *GetSecretMetadataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretMetadataOptions) SetHeaders(param map[string]string) *GetSecretMetadataOptions {
	options.Headers = param
	return options
}

// GetSecretOptions : The GetSecret options.
type GetSecretOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSecretOptions : Instantiate GetSecretOptions
func (*SecretsManagerV2) NewGetSecretOptions(id string) *GetSecretOptions {
	return &GetSecretOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetSecretOptions) SetID(id string) *GetSecretOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretOptions) SetHeaders(param map[string]string) *GetSecretOptions {
	options.Headers = param
	return options
}

// GetSecretVersionMetadataOptions : The GetSecretVersionMetadata options.
type GetSecretVersionMetadataOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies your secret version. You can use the `current` or `previous` aliases to refer
	// to the current or previous secret version.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSecretVersionMetadataOptions : Instantiate GetSecretVersionMetadataOptions
func (*SecretsManagerV2) NewGetSecretVersionMetadataOptions(secretID string, id string) *GetSecretVersionMetadataOptions {
	return &GetSecretVersionMetadataOptions{
		SecretID: core.StringPtr(secretID),
		ID:       core.StringPtr(id),
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *GetSecretVersionMetadataOptions) SetSecretID(secretID string) *GetSecretVersionMetadataOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetSecretVersionMetadataOptions) SetID(id string) *GetSecretVersionMetadataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretVersionMetadataOptions) SetHeaders(param map[string]string) *GetSecretVersionMetadataOptions {
	options.Headers = param
	return options
}

// GetSecretVersionOptions : The GetSecretVersion options.
type GetSecretVersionOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies your secret version. You can use the `current` or `previous` aliases to refer
	// to the current or previous secret version.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSecretVersionOptions : Instantiate GetSecretVersionOptions
func (*SecretsManagerV2) NewGetSecretVersionOptions(secretID string, id string) *GetSecretVersionOptions {
	return &GetSecretVersionOptions{
		SecretID: core.StringPtr(secretID),
		ID:       core.StringPtr(id),
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *GetSecretVersionOptions) SetSecretID(secretID string) *GetSecretVersionOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetSecretVersionOptions) SetID(id string) *GetSecretVersionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretVersionOptions) SetHeaders(param map[string]string) *GetSecretVersionOptions {
	options.Headers = param
	return options
}

// ListConfigurationsOptions : The ListConfigurations options.
type ListConfigurationsOptions struct {
	// The number of configurations to skip. By specifying `offset`, you retrieve a subset of items that starts with the
	// `offset` value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 configurations in your instance, and you want to retrieve configurations 26 through 50,
	// use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// The number of configurations to retrieve. By default, list operations return the first 200 items. To retrieve a
	// different set of items, use `limit` with `offset` to page through your available resources. Maximum limit allowed is
	// 1000 secrets.
	//
	// **Usage:** If you want to retrieve only the first 25 configurations in your instance, use
	// `..?limit=25`.
	Limit *int64 `json:"limit,omitempty"`

	// Sort a collection of configurations by the specified field in ascending order. To sort in descending order use the
	// `-` character
	//
	//
	// **Available values:**  config_type | secret_type | name
	//
	// **Usage:** To sort a list of configurations by their creation date, use
	// `../configurations?sort=config_type`.
	Sort *string `json:"sort,omitempty"`

	// Obtain a collection of configurations that contain the specified string in one or more of the fields: `name`,
	// `config_type`, `secret_type`.
	//
	// **Usage:** If you want to list only the configurations that contain the string `text`, use
	// `../configurations?search=text`.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListConfigurationsOptions : Instantiate ListConfigurationsOptions
func (*SecretsManagerV2) NewListConfigurationsOptions() *ListConfigurationsOptions {
	return &ListConfigurationsOptions{}
}

// SetOffset : Allow user to set Offset
func (_options *ListConfigurationsOptions) SetOffset(offset int64) *ListConfigurationsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListConfigurationsOptions) SetLimit(limit int64) *ListConfigurationsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListConfigurationsOptions) SetSort(sort string) *ListConfigurationsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListConfigurationsOptions) SetSearch(search string) *ListConfigurationsOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListConfigurationsOptions) SetHeaders(param map[string]string) *ListConfigurationsOptions {
	options.Headers = param
	return options
}

// ListSecretGroupsOptions : The ListSecretGroups options.
type ListSecretGroupsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSecretGroupsOptions : Instantiate ListSecretGroupsOptions
func (*SecretsManagerV2) NewListSecretGroupsOptions() *ListSecretGroupsOptions {
	return &ListSecretGroupsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretGroupsOptions) SetHeaders(param map[string]string) *ListSecretGroupsOptions {
	options.Headers = param
	return options
}

// ListSecretLocksOptions : The ListSecretLocks options.
type ListSecretLocksOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	ID *string `json:"id" validate:"required,ne="`

	// The number of locks to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 locks on your secret, and you want to retrieve locks 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// The number of locks with associated secret to retrieve. By default, list operations return the first 25 items. To
	// retrieve a different set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 secrets in your instance, and you want to retrieve only the first 5, use
	// `..?limit=5`.
	Limit *int64 `json:"limit,omitempty"`

	// Sort a collection of locks by the specified field in ascending order. To sort in descending order use the `-`
	// character
	//
	// **Available values:** created_at | updated_at | name
	//
	// **Usage:** To sort a list of locks by their creation date, use
	// `../locks?sort=created_at`.
	Sort *string `json:"sort,omitempty"`

	// Filter locks that contain the specified string in the field "name".
	//
	// **Usage:** If you want to list only the locks that contain the string "text" in the field "name", use
	// `..?search=text`.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSecretLocksOptions : Instantiate ListSecretLocksOptions
func (*SecretsManagerV2) NewListSecretLocksOptions(id string) *ListSecretLocksOptions {
	return &ListSecretLocksOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *ListSecretLocksOptions) SetID(id string) *ListSecretLocksOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListSecretLocksOptions) SetOffset(offset int64) *ListSecretLocksOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSecretLocksOptions) SetLimit(limit int64) *ListSecretLocksOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListSecretLocksOptions) SetSort(sort string) *ListSecretLocksOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListSecretLocksOptions) SetSearch(search string) *ListSecretLocksOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretLocksOptions) SetHeaders(param map[string]string) *ListSecretLocksOptions {
	options.Headers = param
	return options
}

// ListSecretVersionLocksOptions : The ListSecretVersionLocks options.
type ListSecretVersionLocksOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies your secret version. You can use the `current` or `previous` aliases to refer
	// to the current or previous secret version.
	ID *string `json:"id" validate:"required,ne="`

	// The number of locks to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 locks on your secret, and you want to retrieve locks 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// The number of locks with associated secret to retrieve. By default, list operations return the first 25 items. To
	// retrieve a different set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 secrets in your instance, and you want to retrieve only the first 5, use
	// `..?limit=5`.
	Limit *int64 `json:"limit,omitempty"`

	// Sort a collection of locks by the specified field in ascending order. To sort in descending order use the `-`
	// character
	//
	// **Available values:** created_at | updated_at | name
	//
	// **Usage:** To sort a list of locks by their creation date, use
	// `../locks?sort=created_at`.
	Sort *string `json:"sort,omitempty"`

	// Filter locks that contain the specified string in the field "name".
	//
	// **Usage:** If you want to list only the locks that contain the string "text" in the field "name", use
	// `..?search=text`.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSecretVersionLocksOptions : Instantiate ListSecretVersionLocksOptions
func (*SecretsManagerV2) NewListSecretVersionLocksOptions(secretID string, id string) *ListSecretVersionLocksOptions {
	return &ListSecretVersionLocksOptions{
		SecretID: core.StringPtr(secretID),
		ID:       core.StringPtr(id),
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *ListSecretVersionLocksOptions) SetSecretID(secretID string) *ListSecretVersionLocksOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ListSecretVersionLocksOptions) SetID(id string) *ListSecretVersionLocksOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListSecretVersionLocksOptions) SetOffset(offset int64) *ListSecretVersionLocksOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSecretVersionLocksOptions) SetLimit(limit int64) *ListSecretVersionLocksOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListSecretVersionLocksOptions) SetSort(sort string) *ListSecretVersionLocksOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListSecretVersionLocksOptions) SetSearch(search string) *ListSecretVersionLocksOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretVersionLocksOptions) SetHeaders(param map[string]string) *ListSecretVersionLocksOptions {
	options.Headers = param
	return options
}

// ListSecretVersionsOptions : The ListSecretVersions options.
type ListSecretVersionsOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSecretVersionsOptions : Instantiate ListSecretVersionsOptions
func (*SecretsManagerV2) NewListSecretVersionsOptions(secretID string) *ListSecretVersionsOptions {
	return &ListSecretVersionsOptions{
		SecretID: core.StringPtr(secretID),
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *ListSecretVersionsOptions) SetSecretID(secretID string) *ListSecretVersionsOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretVersionsOptions) SetHeaders(param map[string]string) *ListSecretVersionsOptions {
	options.Headers = param
	return options
}

// ListSecretsLocksOptions : The ListSecretsLocks options.
type ListSecretsLocksOptions struct {
	// The number of secrets to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 secrets in your instance, and you want to retrieve secrets 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// The number of secrets to retrieve. By default, list operations return the first 200 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources. Maximum limit allowed is 1000
	// secrets.
	//
	// **Usage:** If you want to retrieve only the first 25 secrets in your instance, use
	// `..?limit=25`.
	Limit *int64 `json:"limit,omitempty"`

	// Filter locks that contain the specified string in the field "name".
	//
	// **Usage:** If you want to list only the locks that contain the string "text" in the field "name", use
	// `..?search=text`.
	Search *string `json:"search,omitempty"`

	// Filter secrets by groups.
	//
	// You can apply multiple filters by using a comma-separated list of secret group IDs. If you need to filter secrets
	// that are in the default secret group, use the `default` keyword.
	//
	// **Usage:** To retrieve a list of secrets that are associated with an existing secret group or the default group, use
	// `..?groups={secret_group_ID},default`.
	Groups []string `json:"groups,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSecretsLocksOptions : Instantiate ListSecretsLocksOptions
func (*SecretsManagerV2) NewListSecretsLocksOptions() *ListSecretsLocksOptions {
	return &ListSecretsLocksOptions{}
}

// SetOffset : Allow user to set Offset
func (_options *ListSecretsLocksOptions) SetOffset(offset int64) *ListSecretsLocksOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSecretsLocksOptions) SetLimit(limit int64) *ListSecretsLocksOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListSecretsLocksOptions) SetSearch(search string) *ListSecretsLocksOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetGroups : Allow user to set Groups
func (_options *ListSecretsLocksOptions) SetGroups(groups []string) *ListSecretsLocksOptions {
	_options.Groups = groups
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretsLocksOptions) SetHeaders(param map[string]string) *ListSecretsLocksOptions {
	options.Headers = param
	return options
}

// ListSecretsOptions : The ListSecrets options.
type ListSecretsOptions struct {
	// The number of secrets to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 secrets in your instance, and you want to retrieve secrets 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// The number of secrets to retrieve. By default, list operations return the first 200 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources. Maximum limit allowed is 1000
	// secrets.
	//
	// **Usage:** If you want to retrieve only the first 25 secrets in your instance, use
	// `..?limit=25`.
	Limit *int64 `json:"limit,omitempty"`

	// Sort a collection of secrets by the specified field in ascending order. To sort in descending order use the `-`
	// character
	//
	//
	// **Available values:** id | created_at | updated_at | expiration_date | secret_type | name
	//
	// **Usage:** To sort a list of secrets by their creation date, use
	// `../secrets?sort=created_at`.
	Sort *string `json:"sort,omitempty"`

	// Obtain a collection of secrets that contain the specified string in one or more of the fields: `id`, `name`,
	// `description`,
	// `labels`, `secret_type`.
	//
	// **Usage:** If you want to list only the secrets that contain the string `text`, use
	// `../secrets?search=text`.
	Search *string `json:"search,omitempty"`

	// Filter secrets by groups.
	//
	// You can apply multiple filters by using a comma-separated list of secret group IDs. If you need to filter secrets
	// that are in the default secret group, use the `default` keyword.
	//
	// **Usage:** To retrieve a list of secrets that are associated with an existing secret group or the default group, use
	// `..?groups={secret_group_ID},default`.
	Groups []string `json:"groups,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSecretsOptions : Instantiate ListSecretsOptions
func (*SecretsManagerV2) NewListSecretsOptions() *ListSecretsOptions {
	return &ListSecretsOptions{}
}

// SetOffset : Allow user to set Offset
func (_options *ListSecretsOptions) SetOffset(offset int64) *ListSecretsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSecretsOptions) SetLimit(limit int64) *ListSecretsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListSecretsOptions) SetSort(sort string) *ListSecretsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListSecretsOptions) SetSearch(search string) *ListSecretsOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetGroups : Allow user to set Groups
func (_options *ListSecretsOptions) SetGroups(groups []string) *ListSecretsOptions {
	_options.Groups = groups
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretsOptions) SetHeaders(param map[string]string) *ListSecretsOptions {
	options.Headers = param
	return options
}

// NotificationsRegistration : The details of the Event Notifications registration.
type NotificationsRegistration struct {
	// A CRN that uniquely identifies an IBM Cloud resource.
	EventNotificationsInstanceCrn *string `json:"event_notifications_instance_crn" validate:"required"`
}

// UnmarshalNotificationsRegistration unmarshals an instance of NotificationsRegistration from the specified map of raw messages.
func UnmarshalNotificationsRegistration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationsRegistration)
	err = core.UnmarshalPrimitive(m, "event_notifications_instance_crn", &obj.EventNotificationsInstanceCrn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedCollectionFirst : A URL that points to the first page in a collection.
type PaginatedCollectionFirst struct {
	// A URL that points to a page in a collection.
	Href *string `json:"href" validate:"required"`
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

// PaginatedCollectionLast : A URL that points to the last page in a collection.
type PaginatedCollectionLast struct {
	// A URL that points to a page in a collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPaginatedCollectionLast unmarshals an instance of PaginatedCollectionLast from the specified map of raw messages.
func UnmarshalPaginatedCollectionLast(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedCollectionLast)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedCollectionNext : A URL that points to the next page in a collection.
type PaginatedCollectionNext struct {
	// A URL that points to a page in a collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPaginatedCollectionNext unmarshals an instance of PaginatedCollectionNext from the specified map of raw messages.
func UnmarshalPaginatedCollectionNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedCollectionNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedCollectionPrevious : A URL that points to the previous page in a collection.
type PaginatedCollectionPrevious struct {
	// A URL that points to a page in a collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPaginatedCollectionPrevious unmarshals an instance of PaginatedCollectionPrevious from the specified map of raw messages.
func UnmarshalPaginatedCollectionPrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedCollectionPrevious)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateCAData : The configuration data of your Private Certificate.
// Models which "extend" this model:
// - PrivateCertificateConfigurationIntermediateCACSR
// - PrivateCertificateConfigurationCACertificate
type PrivateCertificateCAData struct {
	// The certificate signing request.
	Csr *string `json:"csr,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The type of private key to generate.
	PrivateKeyType *string `json:"private_key_type,omitempty"`

	// The certificate expiration time.
	Expiration *int64 `json:"expiration,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The PEM-encoded certificate of the certificate authority that signed and issued this certificate.
	IssuingCa *string `json:"issuing_ca,omitempty"`

	// The chain of certificate authorities that are associated with the certificate.
	CaChain []string `json:"ca_chain,omitempty"`
}

// Constants associated with the PrivateCertificateCAData.PrivateKeyType property.
// The type of private key to generate.
const (
	PrivateCertificateCAData_PrivateKeyType_Ec  = "ec"
	PrivateCertificateCAData_PrivateKeyType_Rsa = "rsa"
)

func (*PrivateCertificateCAData) isaPrivateCertificateCAData() bool {
	return true
}

type PrivateCertificateCADataIntf interface {
	isaPrivateCertificateCAData() bool
}

// UnmarshalPrivateCertificateCAData unmarshals an instance of PrivateCertificateCAData from the specified map of raw messages.
func UnmarshalPrivateCertificateCAData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateCAData)
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_type", &obj.PrivateKeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration", &obj.Expiration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_ca", &obj.IssuingCa)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ca_chain", &obj.CaChain)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateRotationObject : Defines the rotation object that is used to manually rotate public certificates.
type PublicCertificateRotationObject struct {
	// This field indicates whether Secrets Manager rotates the private key for your public certificate automatically.
	//
	// The default is `false`. If it is set to `true`, the service generates and stores a new private key for your rotated
	// certificate.
	RotateKeys *bool `json:"rotate_keys,omitempty"`
}

// UnmarshalPublicCertificateRotationObject unmarshals an instance of PublicCertificateRotationObject from the specified map of raw messages.
func UnmarshalPublicCertificateRotationObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateRotationObject)
	err = core.UnmarshalPrimitive(m, "rotate_keys", &obj.RotateKeys)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RotationPolicy : This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
// username_password, private_cert, public_cert, iam_credentials.
// Models which "extend" this model:
// - CommonRotationPolicy
// - PublicCertificateRotationPolicy
type RotationPolicy struct {
	// This field indicates whether Secrets Manager rotates your secret automatically.
	//
	// The default is `false`. If `auto_rotate` is set to `true` the service rotates your secret based on the defined
	// interval.
	AutoRotate *bool `json:"auto_rotate,omitempty"`

	// The length of the secret rotation time interval.
	Interval *int64 `json:"interval,omitempty"`

	// The units for the secret rotation time interval.
	Unit *string `json:"unit,omitempty"`

	// This field indicates whether Secrets Manager rotates the private key for your public certificate automatically.
	//
	// The default is `false`. If it is set to `true`, the service generates and stores a new private key for your rotated
	// certificate.
	RotateKeys *bool `json:"rotate_keys,omitempty"`
}

// Constants associated with the RotationPolicy.Unit property.
// The units for the secret rotation time interval.
const (
	RotationPolicy_Unit_Day   = "day"
	RotationPolicy_Unit_Month = "month"
)

func (*RotationPolicy) isaRotationPolicy() bool {
	return true
}

type RotationPolicyIntf interface {
	isaRotationPolicy() bool
}

// UnmarshalRotationPolicy unmarshals an instance of RotationPolicy from the specified map of raw messages.
func UnmarshalRotationPolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotationPolicy)
	err = core.UnmarshalPrimitive(m, "auto_rotate", &obj.AutoRotate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unit", &obj.Unit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rotate_keys", &obj.RotateKeys)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Secret : Your secret.
// Models which "extend" this model:
// - ArbitrarySecret
// - IAMCredentialsSecret
// - ImportedCertificate
// - KVSecret
// - PrivateCertificate
// - PublicCertificate
// - ServiceCredentialsSecret
// - UsernamePasswordSecret
type Secret struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type,omitempty"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The secret data that is assigned to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`

	// Access Groups that you can use for an `iam_credentials` secret.
	//
	// Up to 10 Access Groups can be used for each secret.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The ID of the API key that is generated for this secret.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// Indicates whether an `iam_credentials` secret was created with a static service ID.
	//
	// If it is set to `true`, the service ID for the secret was provided by the user at secret creation. If it is set to
	// `false`, the service ID was generated by Secrets Manager.
	ServiceIdIsStatic *bool `json:"service_id_is_static,omitempty"`

	// (IAM credentials) This parameter indicates whether to reuse the service ID and API key for future read operations.
	//
	// If it is set to `true`, the service reuses the current credentials. If it is set to `false`, a new service ID and
	// API key are generated each time that the secret is read or accessed.
	ReuseApiKey *bool `json:"reuse_api_key,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease, the API key is deleted automatically. See the `time-to-live` field to
	// understand the duration of the lease. If you want to continue to use the same API key for future read operations,
	// see the `reuse_api_key` field.
	ApiKey *string `json:"api_key,omitempty"`

	// The identifier for the cryptographic algorithm that is used by the issuing certificate authority to sign a
	// certificate.
	SigningAlgorithm *string `json:"signing_algorithm,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The Common Name (CN) represents the server name protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// Indicates whether the certificate was imported with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The identifier for the cryptographic algorithm used to generate the public key that is associated with the
	// certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// Indicates whether the certificate was imported with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The payload data of a key-value secret.
	Data map[string]interface{} `json:"data,omitempty"`

	// The intermediate certificate authority that signed this certificate.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`

	// The date and time that the certificate was revoked. The date format follows `RFC 3339`.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`

	// The PEM-encoded certificate of the certificate authority that signed and issued this certificate.
	IssuingCa *string `json:"issuing_ca,omitempty"`

	// The chain of certificate authorities that are associated with the certificate.
	CaChain []string `json:"ca_chain,omitempty"`

	// Issuance information that is associated with your certificate.
	IssuanceInfo *CertificateIssuanceInfo `json:"issuance_info,omitempty"`

	// Indicates whether the issued certificate is bundled with intermediate certificates.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The name of the certificate authority configuration.
	Ca *string `json:"ca,omitempty"`

	// The name of the DNS provider configuration.
	Dns *string `json:"dns,omitempty"`

	// The properties that are required to create the service credentials for the specified source service instance.
	SourceService *ServiceCredentialsSecretSourceService `json:"source_service,omitempty"`

	// The properties of the service credentials secret payload.
	Credentials *ServiceCredentialsSecretCredentials `json:"credentials,omitempty"`

	// The username that is assigned to an `username_password` secret.
	Username *string `json:"username,omitempty"`

	// The password that is assigned to an `username_password` secret.
	Password *string `json:"password,omitempty"`
}

// Constants associated with the Secret.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	Secret_SecretType_Arbitrary          = "arbitrary"
	Secret_SecretType_IamCredentials     = "iam_credentials"
	Secret_SecretType_ImportedCert       = "imported_cert"
	Secret_SecretType_Kv                 = "kv"
	Secret_SecretType_PrivateCert        = "private_cert"
	Secret_SecretType_PublicCert         = "public_cert"
	Secret_SecretType_ServiceCredentials = "service_credentials"
	Secret_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the Secret.StateDescription property.
// A text representation of the secret state.
const (
	Secret_StateDescription_Active        = "active"
	Secret_StateDescription_Deactivated   = "deactivated"
	Secret_StateDescription_Destroyed     = "destroyed"
	Secret_StateDescription_PreActivation = "pre_activation"
	Secret_StateDescription_Suspended     = "suspended"
)

func (*Secret) isaSecret() bool {
	return true
}

type SecretIntf interface {
	isaSecret() bool
}

// UnmarshalSecret unmarshals an instance of Secret from the specified map of raw messages.
func UnmarshalSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "secret_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'secret_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'secret_type' not found in JSON object")
		return
	}
	if discValue == "arbitrary" {
		err = core.UnmarshalModel(m, "", result, UnmarshalArbitrarySecret)
	} else if discValue == "iam_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalIAMCredentialsSecret)
	} else if discValue == "imported_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalImportedCertificate)
	} else if discValue == "kv" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKVSecret)
	} else if discValue == "private_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificate)
	} else if discValue == "public_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificate)
	} else if discValue == "service_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalServiceCredentialsSecret)
	} else if discValue == "username_password" {
		err = core.UnmarshalModel(m, "", result, UnmarshalUsernamePasswordSecret)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'secret_type': %s", discValue)
	}
	return
}

// SecretAction : The response body to specify the properties of the action to create a secret.
// Models which "extend" this model:
// - PublicCertificateActionValidateManualDNS
// - PrivateCertificateActionRevoke
type SecretAction struct {
	// The type of secret action.
	ActionType *string `json:"action_type,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`
}

// Constants associated with the SecretAction.ActionType property.
// The type of secret action.
const (
	SecretAction_ActionType_PrivateCertActionRevokeCertificate   = "private_cert_action_revoke_certificate"
	SecretAction_ActionType_PublicCertActionValidateDnsChallenge = "public_cert_action_validate_dns_challenge"
)

func (*SecretAction) isaSecretAction() bool {
	return true
}

type SecretActionIntf interface {
	isaSecretAction() bool
}

// UnmarshalSecretAction unmarshals an instance of SecretAction from the specified map of raw messages.
func UnmarshalSecretAction(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "action_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'action_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'action_type' not found in JSON object")
		return
	}
	if discValue == "private_cert_action_revoke_certificate" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateActionRevoke)
	} else if discValue == "public_cert_action_validate_dns_challenge" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateActionValidateManualDNS)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'action_type': %s", discValue)
	}
	return
}

// SecretActionPrototype : The request body to specify the properties for your secret action.
// Models which "extend" this model:
// - PrivateCertificateActionRevokePrototype
// - PublicCertificateActionValidateManualDNSPrototype
type SecretActionPrototype struct {
	// The type of secret action.
	ActionType *string `json:"action_type,omitempty"`
}

// Constants associated with the SecretActionPrototype.ActionType property.
// The type of secret action.
const (
	SecretActionPrototype_ActionType_PrivateCertActionRevokeCertificate   = "private_cert_action_revoke_certificate"
	SecretActionPrototype_ActionType_PublicCertActionValidateDnsChallenge = "public_cert_action_validate_dns_challenge"
)

func (*SecretActionPrototype) isaSecretActionPrototype() bool {
	return true
}

type SecretActionPrototypeIntf interface {
	isaSecretActionPrototype() bool
}

// UnmarshalSecretActionPrototype unmarshals an instance of SecretActionPrototype from the specified map of raw messages.
func UnmarshalSecretActionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "action_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'action_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'action_type' not found in JSON object")
		return
	}
	if discValue == "public_cert_action_validate_dns_challenge" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateActionValidateManualDNSPrototype)
	} else if discValue == "private_cert_action_revoke_certificate" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateActionRevokePrototype)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'action_type': %s", discValue)
	}
	return
}

// SecretGroup : Properties that describe a secret group.
type SecretGroup struct {
	// A v4 UUID identifier, or `default` secret group.
	ID *string `json:"id" validate:"required"`

	// The name of your existing secret group.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret group.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`
}

// UnmarshalSecretGroup unmarshals an instance of SecretGroup from the specified map of raw messages.
func UnmarshalSecretGroup(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretGroup)
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
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
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

// SecretGroupCollection : Properties that describe a collection of secret groups.
type SecretGroupCollection struct {
	// A collection of secret groups.
	SecretGroups []SecretGroup `json:"secret_groups" validate:"required"`

	// The total number of resources in a collection.
	TotalCount *int64 `json:"total_count" validate:"required"`
}

// UnmarshalSecretGroupCollection unmarshals an instance of SecretGroupCollection from the specified map of raw messages.
func UnmarshalSecretGroupCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretGroupCollection)
	err = core.UnmarshalModel(m, "secret_groups", &obj.SecretGroups, UnmarshalSecretGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretGroupPatch : Update the name or description of your secret group.
type SecretGroupPatch struct {
	// The name of your secret group.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret group.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`
}

// UnmarshalSecretGroupPatch unmarshals an instance of SecretGroupPatch from the specified map of raw messages.
func UnmarshalSecretGroupPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretGroupPatch)
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

// AsPatch returns a generic map representation of the SecretGroupPatch
func (secretGroupPatch *SecretGroupPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(secretGroupPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// SecretLock : SecretLock struct
type SecretLock struct {
	// A human-readable name to assign to the lock. The lock name must be unique per secret version.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret lock.
	Name *string `json:"name" validate:"required"`

	// An extended description of the lock.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// lock.
	Description *string `json:"description,omitempty"`

	// Optional information to associate with a lock, such as resources CRNs to be used by automation.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// A v4 UUID identifier.
	SecretVersionID *string `json:"secret_version_id" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	SecretVersionAlias *string `json:"secret_version_alias" validate:"required"`
}

// Constants associated with the SecretLock.SecretVersionAlias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	SecretLock_SecretVersionAlias_Current  = "current"
	SecretLock_SecretVersionAlias_Previous = "previous"
)

// UnmarshalSecretLock unmarshals an instance of SecretLock from the specified map of raw messages.
func UnmarshalSecretLock(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretLock)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
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
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_version_id", &obj.SecretVersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_version_alias", &obj.SecretVersionAlias)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretLockPrototype : SecretLockPrototype struct
type SecretLockPrototype struct {
	// A human-readable name to assign to the lock. The lock name must be unique per secret version.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret lock.
	Name *string `json:"name" validate:"required"`

	// An extended description of the lock.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// lock.
	Description *string `json:"description,omitempty"`

	// Optional information to associate with a lock, such as resources CRNs to be used by automation.
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// NewSecretLockPrototype : Instantiate SecretLockPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewSecretLockPrototype(name string) (_model *SecretLockPrototype, err error) {
	_model = &SecretLockPrototype{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSecretLockPrototype unmarshals an instance of SecretLockPrototype from the specified map of raw messages.
func UnmarshalSecretLockPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretLockPrototype)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attributes", &obj.Attributes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretLocks : Create locks response body containing a collection of locks that are attached to a secret.
type SecretLocks struct {
	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type,omitempty"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// A collection of locks that are attached to a secret.
	Versions []SecretVersionLocks `json:"versions" validate:"required"`
}

// Constants associated with the SecretLocks.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	SecretLocks_SecretType_Arbitrary          = "arbitrary"
	SecretLocks_SecretType_IamCredentials     = "iam_credentials"
	SecretLocks_SecretType_ImportedCert       = "imported_cert"
	SecretLocks_SecretType_Kv                 = "kv"
	SecretLocks_SecretType_PrivateCert        = "private_cert"
	SecretLocks_SecretType_PublicCert         = "public_cert"
	SecretLocks_SecretType_ServiceCredentials = "service_credentials"
	SecretLocks_SecretType_UsernamePassword   = "username_password"
)

// UnmarshalSecretLocks unmarshals an instance of SecretLocks from the specified map of raw messages.
func UnmarshalSecretLocks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretLocks)
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalSecretVersionLocks)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretLocksPaginatedCollection : Properties that describe a paginated collection of your secret locks.
type SecretLocksPaginatedCollection struct {
	// The total number of resources in a collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of items that are retrieved in a collection.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of items that are skipped in a collection.
	Offset *int64 `json:"offset" validate:"required"`

	// A URL that points to the first page in a collection.
	First *PaginatedCollectionFirst `json:"first" validate:"required"`

	// A URL that points to the next page in a collection.
	Next *PaginatedCollectionNext `json:"next,omitempty"`

	// A URL that points to the previous page in a collection.
	Previous *PaginatedCollectionPrevious `json:"previous,omitempty"`

	// A URL that points to the last page in a collection.
	Last *PaginatedCollectionLast `json:"last" validate:"required"`

	// A collection of secret locks.
	Locks []SecretLock `json:"locks" validate:"required"`
}

// UnmarshalSecretLocksPaginatedCollection unmarshals an instance of SecretLocksPaginatedCollection from the specified map of raw messages.
func UnmarshalSecretLocksPaginatedCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretLocksPaginatedCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
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
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedCollectionPrevious)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedCollectionLast)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "locks", &obj.Locks, UnmarshalSecretLock)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SecretLocksPaginatedCollection) GetNextOffset() (*int64, error) {
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

// SecretMetadata : Properties of your secret metadata.
// Models which "extend" this model:
// - ArbitrarySecretMetadata
// - IAMCredentialsSecretMetadata
// - ImportedCertificateMetadata
// - KVSecretMetadata
// - PrivateCertificateMetadata
// - PublicCertificateMetadata
// - ServiceCredentialsSecretMetadata
// - UsernamePasswordSecretMetadata
type SecretMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type,omitempty"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`

	// Access Groups that you can use for an `iam_credentials` secret.
	//
	// Up to 10 Access Groups can be used for each secret.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The ID of the API key that is generated for this secret.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// Indicates whether an `iam_credentials` secret was created with a static service ID.
	//
	// If it is set to `true`, the service ID for the secret was provided by the user at secret creation. If it is set to
	// `false`, the service ID was generated by Secrets Manager.
	ServiceIdIsStatic *bool `json:"service_id_is_static,omitempty"`

	// (IAM credentials) This parameter indicates whether to reuse the service ID and API key for future read operations.
	//
	// If it is set to `true`, the service reuses the current credentials. If it is set to `false`, a new service ID and
	// API key are generated each time that the secret is read or accessed.
	ReuseApiKey *bool `json:"reuse_api_key,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// The identifier for the cryptographic algorithm that is used by the issuing certificate authority to sign a
	// certificate.
	SigningAlgorithm *string `json:"signing_algorithm,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The Common Name (CN) represents the server name protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// Indicates whether the certificate was imported with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The identifier for the cryptographic algorithm used to generate the public key that is associated with the
	// certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// Indicates whether the certificate was imported with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity,omitempty"`

	// The intermediate certificate authority that signed this certificate.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`

	// The date and time that the certificate was revoked. The date format follows `RFC 3339`.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`

	// Issuance information that is associated with your certificate.
	IssuanceInfo *CertificateIssuanceInfo `json:"issuance_info,omitempty"`

	// Indicates whether the issued certificate is bundled with intermediate certificates.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The name of the certificate authority configuration.
	Ca *string `json:"ca,omitempty"`

	// The name of the DNS provider configuration.
	Dns *string `json:"dns,omitempty"`

	// The properties that are required to create the service credentials for the specified source service instance.
	SourceService *ServiceCredentialsSecretSourceService `json:"source_service,omitempty"`
}

// Constants associated with the SecretMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	SecretMetadata_SecretType_Arbitrary          = "arbitrary"
	SecretMetadata_SecretType_IamCredentials     = "iam_credentials"
	SecretMetadata_SecretType_ImportedCert       = "imported_cert"
	SecretMetadata_SecretType_Kv                 = "kv"
	SecretMetadata_SecretType_PrivateCert        = "private_cert"
	SecretMetadata_SecretType_PublicCert         = "public_cert"
	SecretMetadata_SecretType_ServiceCredentials = "service_credentials"
	SecretMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the SecretMetadata.StateDescription property.
// A text representation of the secret state.
const (
	SecretMetadata_StateDescription_Active        = "active"
	SecretMetadata_StateDescription_Deactivated   = "deactivated"
	SecretMetadata_StateDescription_Destroyed     = "destroyed"
	SecretMetadata_StateDescription_PreActivation = "pre_activation"
	SecretMetadata_StateDescription_Suspended     = "suspended"
)

func (*SecretMetadata) isaSecretMetadata() bool {
	return true
}

type SecretMetadataIntf interface {
	isaSecretMetadata() bool
}

// UnmarshalSecretMetadata unmarshals an instance of SecretMetadata from the specified map of raw messages.
func UnmarshalSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "secret_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'secret_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'secret_type' not found in JSON object")
		return
	}
	if discValue == "arbitrary" {
		err = core.UnmarshalModel(m, "", result, UnmarshalArbitrarySecretMetadata)
	} else if discValue == "iam_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalIAMCredentialsSecretMetadata)
	} else if discValue == "imported_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalImportedCertificateMetadata)
	} else if discValue == "kv" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKVSecretMetadata)
	} else if discValue == "private_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateMetadata)
	} else if discValue == "public_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateMetadata)
	} else if discValue == "service_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalServiceCredentialsSecretMetadata)
	} else if discValue == "username_password" {
		err = core.UnmarshalModel(m, "", result, UnmarshalUsernamePasswordSecretMetadata)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'secret_type': %s", discValue)
	}
	return
}

// SecretMetadataPaginatedCollection : Properties that describe a paginated collection of your secret metadata.
type SecretMetadataPaginatedCollection struct {
	// The total number of resources in a collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of items that are retrieved in a collection.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of items that are skipped in a collection.
	Offset *int64 `json:"offset" validate:"required"`

	// A URL that points to the first page in a collection.
	First *PaginatedCollectionFirst `json:"first" validate:"required"`

	// A URL that points to the next page in a collection.
	Next *PaginatedCollectionNext `json:"next,omitempty"`

	// A URL that points to the previous page in a collection.
	Previous *PaginatedCollectionPrevious `json:"previous,omitempty"`

	// A URL that points to the last page in a collection.
	Last *PaginatedCollectionLast `json:"last" validate:"required"`

	// A collection of secret metadata.
	Secrets []SecretMetadataIntf `json:"secrets" validate:"required"`
}

// UnmarshalSecretMetadataPaginatedCollection unmarshals an instance of SecretMetadataPaginatedCollection from the specified map of raw messages.
func UnmarshalSecretMetadataPaginatedCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretMetadataPaginatedCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
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
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedCollectionPrevious)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedCollectionLast)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secrets", &obj.Secrets, UnmarshalSecretMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SecretMetadataPaginatedCollection) GetNextOffset() (*int64, error) {
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

// SecretMetadataPatch : Update your secret metadata.
// Models which "extend" this model:
// - ArbitrarySecretMetadataPatch
// - IAMCredentialsSecretMetadataPatch
// - ImportedCertificateMetadataPatch
// - KVSecretMetadataPatch
// - PrivateCertificateMetadataPatch
// - PublicCertificateMetadataPatch
// - ServiceCredentialsSecretMetadataPatch
// - UsernamePasswordSecretMetadataPatch
type SecretMetadataPatch struct {
	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`
}

func (*SecretMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

type SecretMetadataPatchIntf interface {
	isaSecretMetadataPatch() bool
}

// UnmarshalSecretMetadataPatch unmarshals an instance of SecretMetadataPatch from the specified map of raw messages.
func UnmarshalSecretMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretMetadataPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the SecretMetadataPatch
func (secretMetadataPatch *SecretMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(secretMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// SecretPrototype : Specify the properties for your secret.
// Models which "extend" this model:
// - ArbitrarySecretPrototype
// - IAMCredentialsSecretPrototype
// - ImportedCertificatePrototype
// - KVSecretPrototype
// - PrivateCertificatePrototype
// - PublicCertificatePrototype
// - ServiceCredentialsSecretPrototype
// - UsernamePasswordSecretPrototype
type SecretPrototype struct {
	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type,omitempty"`

	// The secret data that is assigned to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`

	// Access Groups that you can use for an `iam_credentials` secret.
	//
	// Up to 10 Access Groups can be used for each secret.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// (IAM credentials) This parameter indicates whether to reuse the service ID and API key for future read operations.
	//
	// If it is set to `true`, the service reuses the current credentials. If it is set to `false`, a new service ID and
	// API key are generated each time that the secret is read or accessed.
	ReuseApiKey *bool `json:"reuse_api_key,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The payload data of a key-value secret.
	Data map[string]interface{} `json:"data,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template,omitempty"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The certificate signing request.
	Csr *string `json:"csr,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The identifier for the cryptographic algorithm that is used to generate the public key that is associated with the
	// certificate.
	//
	// The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to
	// generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide
	// more encryption protection. Allowed values:  `RSA2048`, `RSA4096`, `ECDSA256`, and `ECDSA384`.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The name of the certificate authority configuration.
	Ca *string `json:"ca,omitempty"`

	// The name of the DNS provider configuration.
	Dns *string `json:"dns,omitempty"`

	// This field indicates whether your issued certificate is bundled with intermediate certificates. Set to `false` for
	// the certificate file to contain only the issued certificate.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The properties that are required to create the service credentials for the specified source service instance.
	SourceService *ServiceCredentialsSecretSourceService `json:"source_service,omitempty"`

	// The username that is assigned to an `username_password` secret.
	Username *string `json:"username,omitempty"`

	// The password that is assigned to an `username_password` secret.
	Password *string `json:"password,omitempty"`
}

// Constants associated with the SecretPrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	SecretPrototype_SecretType_Arbitrary          = "arbitrary"
	SecretPrototype_SecretType_IamCredentials     = "iam_credentials"
	SecretPrototype_SecretType_ImportedCert       = "imported_cert"
	SecretPrototype_SecretType_Kv                 = "kv"
	SecretPrototype_SecretType_PrivateCert        = "private_cert"
	SecretPrototype_SecretType_PublicCert         = "public_cert"
	SecretPrototype_SecretType_ServiceCredentials = "service_credentials"
	SecretPrototype_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the SecretPrototype.Format property.
// The format of the returned data.
const (
	SecretPrototype_Format_Pem       = "pem"
	SecretPrototype_Format_PemBundle = "pem_bundle"
)

// Constants associated with the SecretPrototype.PrivateKeyFormat property.
// The format of the generated private key.
const (
	SecretPrototype_PrivateKeyFormat_Der   = "der"
	SecretPrototype_PrivateKeyFormat_Pkcs8 = "pkcs8"
)

func (*SecretPrototype) isaSecretPrototype() bool {
	return true
}

type SecretPrototypeIntf interface {
	isaSecretPrototype() bool
}

// UnmarshalSecretPrototype unmarshals an instance of SecretPrototype from the specified map of raw messages.
func UnmarshalSecretPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "secret_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'secret_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'secret_type' not found in JSON object")
		return
	}
	if discValue == "arbitrary" {
		err = core.UnmarshalModel(m, "", result, UnmarshalArbitrarySecretPrototype)
	} else if discValue == "iam_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalIAMCredentialsSecretPrototype)
	} else if discValue == "imported_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalImportedCertificatePrototype)
	} else if discValue == "kv" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKVSecretPrototype)
	} else if discValue == "private_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificatePrototype)
	} else if discValue == "public_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificatePrototype)
	} else if discValue == "service_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalServiceCredentialsSecretPrototype)
	} else if discValue == "username_password" {
		err = core.UnmarshalModel(m, "", result, UnmarshalUsernamePasswordSecretPrototype)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'secret_type': %s", discValue)
	}
	return
}

// SecretVersion : Your secret version.
// Models which "extend" this model:
// - ArbitrarySecretVersion
// - IAMCredentialsSecretVersion
// - ImportedCertificateVersion
// - KVSecretVersion
// - PrivateCertificateVersion
// - PublicCertificateVersion
// - ServiceCredentialsSecretVersion
// - UsernamePasswordSecretVersion
type SecretVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id,omitempty"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The secret data that is assigned to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`

	// The ID of the API key that is generated for this secret.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease, the API key is deleted automatically. See the `time-to-live` field to
	// understand the duration of the lease. If you want to continue to use the same API key for future read operations,
	// see the `reuse_api_key` field.
	ApiKey *string `json:"api_key,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The payload data of a key-value secret.
	Data map[string]interface{} `json:"data,omitempty"`

	// The PEM-encoded certificate of the certificate authority that signed and issued this certificate.
	IssuingCa *string `json:"issuing_ca,omitempty"`

	// The chain of certificate authorities that are associated with the certificate.
	CaChain []string `json:"ca_chain,omitempty"`

	// The source service resource key data of the generated service credentials.
	ResourceKey *ServiceCredentialsResourceKey `json:"resource_key,omitempty"`

	// The properties of the service credentials secret payload.
	Credentials *ServiceCredentialsSecretCredentials `json:"credentials,omitempty"`

	// The username that is assigned to an `username_password` secret.
	Username *string `json:"username,omitempty"`

	// The password that is assigned to an `username_password` secret.
	Password *string `json:"password,omitempty"`
}

// Constants associated with the SecretVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	SecretVersion_SecretType_Arbitrary          = "arbitrary"
	SecretVersion_SecretType_IamCredentials     = "iam_credentials"
	SecretVersion_SecretType_ImportedCert       = "imported_cert"
	SecretVersion_SecretType_Kv                 = "kv"
	SecretVersion_SecretType_PrivateCert        = "private_cert"
	SecretVersion_SecretType_PublicCert         = "public_cert"
	SecretVersion_SecretType_ServiceCredentials = "service_credentials"
	SecretVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the SecretVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	SecretVersion_Alias_Current  = "current"
	SecretVersion_Alias_Previous = "previous"
)

func (*SecretVersion) isaSecretVersion() bool {
	return true
}

type SecretVersionIntf interface {
	isaSecretVersion() bool
}

// UnmarshalSecretVersion unmarshals an instance of SecretVersion from the specified map of raw messages.
func UnmarshalSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "secret_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'secret_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'secret_type' not found in JSON object")
		return
	}
	if discValue == "arbitrary" {
		err = core.UnmarshalModel(m, "", result, UnmarshalArbitrarySecretVersion)
	} else if discValue == "iam_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalIAMCredentialsSecretVersion)
	} else if discValue == "imported_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalImportedCertificateVersion)
	} else if discValue == "kv" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKVSecretVersion)
	} else if discValue == "private_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateVersion)
	} else if discValue == "public_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateVersion)
	} else if discValue == "service_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalServiceCredentialsSecretVersion)
	} else if discValue == "username_password" {
		err = core.UnmarshalModel(m, "", result, UnmarshalUsernamePasswordSecretVersion)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'secret_type': %s", discValue)
	}
	return
}

// SecretVersionActionPrototype : The request body to specify the properties of the action to create a secret version.
// Models which "extend" this model:
// - PrivateCertificateVersionActionRevokePrototype
type SecretVersionActionPrototype struct {
	// The type of secret version action.
	ActionType *string `json:"action_type,omitempty"`
}

// Constants associated with the SecretVersionActionPrototype.ActionType property.
// The type of secret version action.
const (
	SecretVersionActionPrototype_ActionType_PrivateCertActionRevokeCertificate = "private_cert_action_revoke_certificate"
)

func (*SecretVersionActionPrototype) isaSecretVersionActionPrototype() bool {
	return true
}

type SecretVersionActionPrototypeIntf interface {
	isaSecretVersionActionPrototype() bool
}

// UnmarshalSecretVersionActionPrototype unmarshals an instance of SecretVersionActionPrototype from the specified map of raw messages.
func UnmarshalSecretVersionActionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "action_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'action_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'action_type' not found in JSON object")
		return
	}
	if discValue == "private_cert_action_revoke_certificate" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateVersionActionRevokePrototype)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'action_type': %s", discValue)
	}
	return
}

// SecretVersionLocks : SecretVersionLocks struct
type SecretVersionLocks struct {
	// A v4 UUID identifier.
	VersionID *string `json:"version_id" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	VersionAlias *string `json:"version_alias" validate:"required"`

	// The names of all locks that are associated with this secret version.
	Locks []string `json:"locks" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available,omitempty"`
}

// Constants associated with the SecretVersionLocks.VersionAlias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	SecretVersionLocks_VersionAlias_Current  = "current"
	SecretVersionLocks_VersionAlias_Previous = "previous"
)

// UnmarshalSecretVersionLocks unmarshals an instance of SecretVersionLocks from the specified map of raw messages.
func UnmarshalSecretVersionLocks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersionLocks)
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_alias", &obj.VersionAlias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks", &obj.Locks)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretVersionLocksPaginatedCollection : Properties that describe a paginated collection of your secret version locks.
type SecretVersionLocksPaginatedCollection struct {
	// The total number of resources in a collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of items that are retrieved in a collection.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of items that are skipped in a collection.
	Offset *int64 `json:"offset" validate:"required"`

	// A URL that points to the first page in a collection.
	First *PaginatedCollectionFirst `json:"first" validate:"required"`

	// A URL that points to the next page in a collection.
	Next *PaginatedCollectionNext `json:"next,omitempty"`

	// A URL that points to the previous page in a collection.
	Previous *PaginatedCollectionPrevious `json:"previous,omitempty"`

	// A URL that points to the last page in a collection.
	Last *PaginatedCollectionLast `json:"last" validate:"required"`

	// A collection of secret version locks.
	Locks []SecretLock `json:"locks" validate:"required"`
}

// UnmarshalSecretVersionLocksPaginatedCollection unmarshals an instance of SecretVersionLocksPaginatedCollection from the specified map of raw messages.
func UnmarshalSecretVersionLocksPaginatedCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersionLocksPaginatedCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
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
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedCollectionPrevious)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedCollectionLast)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "locks", &obj.Locks, UnmarshalSecretLock)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SecretVersionLocksPaginatedCollection) GetNextOffset() (*int64, error) {
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

// SecretVersionMetadata : Properties of the version metadata of your secret.
// Models which "extend" this model:
// - ArbitrarySecretVersionMetadata
// - IAMCredentialsSecretVersionMetadata
// - ImportedCertificateVersionMetadata
// - KVSecretVersionMetadata
// - PrivateCertificateVersionMetadata
// - PublicCertificateVersionMetadata
// - ServiceCredentialsSecretVersionMetadata
// - UsernamePasswordSecretVersionMetadata
type SecretVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id,omitempty"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The ID of the API key that is generated for this secret.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity,omitempty"`

	// The source service resource key data of the generated service credentials.
	ResourceKey *ServiceCredentialsResourceKey `json:"resource_key,omitempty"`
}

// Constants associated with the SecretVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	SecretVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	SecretVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	SecretVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	SecretVersionMetadata_SecretType_Kv                 = "kv"
	SecretVersionMetadata_SecretType_PrivateCert        = "private_cert"
	SecretVersionMetadata_SecretType_PublicCert         = "public_cert"
	SecretVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	SecretVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the SecretVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	SecretVersionMetadata_Alias_Current  = "current"
	SecretVersionMetadata_Alias_Previous = "previous"
)

func (*SecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

type SecretVersionMetadataIntf interface {
	isaSecretVersionMetadata() bool
}

// UnmarshalSecretVersionMetadata unmarshals an instance of SecretVersionMetadata from the specified map of raw messages.
func UnmarshalSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "secret_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'secret_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'secret_type' not found in JSON object")
		return
	}
	if discValue == "arbitrary" {
		err = core.UnmarshalModel(m, "", result, UnmarshalArbitrarySecretVersionMetadata)
	} else if discValue == "iam_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalIAMCredentialsSecretVersionMetadata)
	} else if discValue == "imported_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalImportedCertificateVersionMetadata)
	} else if discValue == "kv" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKVSecretVersionMetadata)
	} else if discValue == "private_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateVersionMetadata)
	} else if discValue == "public_cert" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPublicCertificateVersionMetadata)
	} else if discValue == "service_credentials" {
		err = core.UnmarshalModel(m, "", result, UnmarshalServiceCredentialsSecretVersionMetadata)
	} else if discValue == "username_password" {
		err = core.UnmarshalModel(m, "", result, UnmarshalUsernamePasswordSecretVersionMetadata)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'secret_type': %s", discValue)
	}
	return
}

// SecretVersionMetadataCollection : Properties that describe a collection of your secret version metadata.
type SecretVersionMetadataCollection struct {
	// A collection of secret version metadata.
	Versions []SecretVersionMetadataIntf `json:"versions" validate:"required"`

	// The total number of resources in a collection.
	TotalCount *int64 `json:"total_count" validate:"required"`
}

// UnmarshalSecretVersionMetadataCollection unmarshals an instance of SecretVersionMetadataCollection from the specified map of raw messages.
func UnmarshalSecretVersionMetadataCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersionMetadataCollection)
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalSecretVersionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretVersionMetadataPatch : Update your secret version metadata.
type SecretVersionMetadataPatch struct {
	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// UnmarshalSecretVersionMetadataPatch unmarshals an instance of SecretVersionMetadataPatch from the specified map of raw messages.
func UnmarshalSecretVersionMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersionMetadataPatch)
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the SecretVersionMetadataPatch
func (secretVersionMetadataPatch *SecretVersionMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(secretVersionMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// SecretVersionPrototype : Specify the properties for your new secret version.
// Models which "extend" this model:
// - ArbitrarySecretVersionPrototype
// - IAMCredentialsSecretRestoreFromVersionPrototype
// - IAMCredentialsSecretVersionPrototype
// - ImportedCertificateVersionPrototype
// - KVSecretVersionPrototype
// - PrivateCertificateVersionPrototype
// - PublicCertificateVersionPrototype
// - ServiceCredentialsSecretVersionPrototype
// - UsernamePasswordSecretVersionPrototype
type SecretVersionPrototype struct {
	// The secret data that is assigned to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier, or `current` or `previous` secret version aliases.
	RestoreFromVersion *string `json:"restore_from_version,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The payload data of a key-value secret.
	Data map[string]interface{} `json:"data,omitempty"`

	// The certificate signing request.
	Csr *string `json:"csr,omitempty"`

	// Defines the rotation object that is used to manually rotate public certificates.
	Rotation *PublicCertificateRotationObject `json:"rotation,omitempty"`

	// The password that is assigned to an `username_password` secret.
	Password *string `json:"password,omitempty"`
}

func (*SecretVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

type SecretVersionPrototypeIntf interface {
	isaSecretVersionPrototype() bool
}

// UnmarshalSecretVersionPrototype unmarshals an instance of SecretVersionPrototype from the specified map of raw messages.
func UnmarshalSecretVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersionPrototype)
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "restore_from_version", &obj.RestoreFromVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalPublicCertificateRotationObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretsLocksPaginatedCollection : Properties that describe a paginated collection of your secrets locks.
type SecretsLocksPaginatedCollection struct {
	// The total number of resources in a collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of items that are retrieved in a collection.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of items that are skipped in a collection.
	Offset *int64 `json:"offset" validate:"required"`

	// A URL that points to the first page in a collection.
	First *PaginatedCollectionFirst `json:"first" validate:"required"`

	// A URL that points to the next page in a collection.
	Next *PaginatedCollectionNext `json:"next,omitempty"`

	// A URL that points to the previous page in a collection.
	Previous *PaginatedCollectionPrevious `json:"previous,omitempty"`

	// A URL that points to the last page in a collection.
	Last *PaginatedCollectionLast `json:"last" validate:"required"`

	// A collection of secrets and their locks.
	SecretsLocks []SecretLocks `json:"secrets_locks" validate:"required"`
}

// UnmarshalSecretsLocksPaginatedCollection unmarshals an instance of SecretsLocksPaginatedCollection from the specified map of raw messages.
func UnmarshalSecretsLocksPaginatedCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretsLocksPaginatedCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
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
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedCollectionPrevious)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedCollectionLast)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secrets_locks", &obj.SecretsLocks, UnmarshalSecretLocks)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SecretsLocksPaginatedCollection) GetNextOffset() (*int64, error) {
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

// ServiceCredentialsResourceKey : The source service resource key data of the generated service credentials.
type ServiceCredentialsResourceKey struct {
	// The resource key CRN of the generated service credentials.
	Crn *string `json:"crn,omitempty"`

	// The resource key name of the generated service credentials.
	Name *string `json:"name,omitempty"`
}

// UnmarshalServiceCredentialsResourceKey unmarshals an instance of ServiceCredentialsResourceKey from the specified map of raw messages.
func UnmarshalServiceCredentialsResourceKey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsResourceKey)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
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

// ServiceCredentialsSecretCredentials : The properties of the service credentials secret payload.
type ServiceCredentialsSecretCredentials struct {
	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease, the API key is deleted automatically. See the `time-to-live` field to
	// understand the duration of the lease.
	Apikey *string `json:"apikey,omitempty"`

	// The Cloud Object Storage HMAC keys that are returned after you create a service credentials secret.
	CosHmacKeys *CosHmacKeys `json:"cos_hmac_keys,omitempty"`

	// The endpoints that are returned after you create a service credentials secret.
	Endpoints *string `json:"endpoints,omitempty"`

	// The IAM API key description for the generated service credentials.
	IamApikeyDescription *string `json:"iam_apikey_description,omitempty"`

	// The IAM API key id for the generated service credentials.
	IamApikeyID *string `json:"iam_apikey_id,omitempty"`

	// The IAM API key name for the generated service credentials.
	IamApikeyName *string `json:"iam_apikey_name,omitempty"`

	// The IAM role CRN assigned to the generated service credentials.
	IamRoleCrn *string `json:"iam_role_crn,omitempty"`

	// The IAM Service ID CRN.
	IamServiceidCrn *string `json:"iam_serviceid_crn,omitempty"`

	// The resource instance CRN that is returned after you create a service credentials secret.
	ResourceInstanceID *string `json:"resource_instance_id,omitempty"`
}

// UnmarshalServiceCredentialsSecretCredentials unmarshals an instance of ServiceCredentialsSecretCredentials from the specified map of raw messages.
func UnmarshalServiceCredentialsSecretCredentials(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecretCredentials)
	err = core.UnmarshalPrimitive(m, "apikey", &obj.Apikey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "cos_hmac_keys", &obj.CosHmacKeys, UnmarshalCosHmacKeys)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "endpoints", &obj.Endpoints)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_apikey_description", &obj.IamApikeyDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_apikey_id", &obj.IamApikeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_apikey_name", &obj.IamApikeyName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_role_crn", &obj.IamRoleCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_serviceid_crn", &obj.IamServiceidCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_instance_id", &obj.ResourceInstanceID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSecretSourceService : The properties that are required to create the service credentials for the specified source service instance.
type ServiceCredentialsSecretSourceService struct {
	// The source service instance identifier.
	Instance *ServiceCredentialsSourceServiceInstance `json:"instance" validate:"required"`

	// Configuration options represented as key-value pairs. Service-defined options are used in the generation of
	// credentials for some services. For example, Cloud Object Storage accepts the optional boolean parameter HMAC for
	// creating specific kind of credentials.
	Parameters *ServiceCredentialsSourceServiceParameters `json:"parameters,omitempty"`

	// The service-specific custom role. CRN is accepted. The role is assigned as part of an access policy to any
	// auto-generated IAM service ID.  If you provide an existing service ID, it is added to the access policy for that ID.
	//  If a role is not provided, any new service IDs that are autogenerated, will not have an assigned access policy and
	// provided service IDs are not changed in any way.  Refer to the service documentation for supported roles.
	Role *ServiceCredentialsSourceServiceRole `json:"role,omitempty"`

	// The source service IAM data is returned in case IAM credentials where created for this secret.
	Iam *ServiceCredentialsSourceServiceIam `json:"iam,omitempty"`

	// The source service resource key data of the generated service credentials.
	ResourceKey *ServiceCredentialsResourceKey `json:"resource_key,omitempty"`
}

// NewServiceCredentialsSecretSourceService : Instantiate ServiceCredentialsSecretSourceService (Generic Model Constructor)
func (*SecretsManagerV2) NewServiceCredentialsSecretSourceService(instance *ServiceCredentialsSourceServiceInstance) (_model *ServiceCredentialsSecretSourceService, err error) {
	_model = &ServiceCredentialsSecretSourceService{
		Instance: instance,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalServiceCredentialsSecretSourceService unmarshals an instance of ServiceCredentialsSecretSourceService from the specified map of raw messages.
func UnmarshalServiceCredentialsSecretSourceService(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecretSourceService)
	err = core.UnmarshalModel(m, "instance", &obj.Instance, UnmarshalServiceCredentialsSourceServiceInstance)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "parameters", &obj.Parameters, UnmarshalServiceCredentialsSourceServiceParameters)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "role", &obj.Role, UnmarshalServiceCredentialsSourceServiceRole)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "iam", &obj.Iam, UnmarshalServiceCredentialsSourceServiceIam)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource_key", &obj.ResourceKey, UnmarshalServiceCredentialsResourceKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSourceServiceIam : The source service IAM data is returned in case IAM credentials where created for this secret.
type ServiceCredentialsSourceServiceIam struct {
	// The IAM apikey metadata for the IAM credentials that were generated.
	Apikey *ServiceCredentialsSourceServiceIamApikey `json:"apikey,omitempty"`

	// The IAM role for the generate service credentials.
	Role *ServiceCredentialsSourceServiceIamRole `json:"role,omitempty"`

	// The IAM serviceid for the generated service credentials.
	Serviceid *ServiceCredentialsSourceServiceIamServiceid `json:"serviceid,omitempty"`
}

// UnmarshalServiceCredentialsSourceServiceIam unmarshals an instance of ServiceCredentialsSourceServiceIam from the specified map of raw messages.
func UnmarshalServiceCredentialsSourceServiceIam(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSourceServiceIam)
	err = core.UnmarshalModel(m, "apikey", &obj.Apikey, UnmarshalServiceCredentialsSourceServiceIamApikey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "role", &obj.Role, UnmarshalServiceCredentialsSourceServiceIamRole)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "serviceid", &obj.Serviceid, UnmarshalServiceCredentialsSourceServiceIamServiceid)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSourceServiceIamApikey : The IAM apikey metadata for the IAM credentials that were generated.
type ServiceCredentialsSourceServiceIamApikey struct {
	// The IAM API key description for the generated service credentials.
	Description *string `json:"description,omitempty"`

	// The IAM API key id for the generated service credentials.
	ID *string `json:"id,omitempty"`

	// The IAM API key name for the generated service credentials.
	Name *string `json:"name,omitempty"`
}

// UnmarshalServiceCredentialsSourceServiceIamApikey unmarshals an instance of ServiceCredentialsSourceServiceIamApikey from the specified map of raw messages.
func UnmarshalServiceCredentialsSourceServiceIamApikey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSourceServiceIamApikey)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
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

// ServiceCredentialsSourceServiceIamRole : The IAM role for the generate service credentials.
type ServiceCredentialsSourceServiceIamRole struct {
	// The IAM role CRN assigned to the generated service credentials.
	Crn *string `json:"crn,omitempty"`
}

// UnmarshalServiceCredentialsSourceServiceIamRole unmarshals an instance of ServiceCredentialsSourceServiceIamRole from the specified map of raw messages.
func UnmarshalServiceCredentialsSourceServiceIamRole(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSourceServiceIamRole)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSourceServiceIamServiceid : The IAM serviceid for the generated service credentials.
type ServiceCredentialsSourceServiceIamServiceid struct {
	// The IAM Service ID CRN.
	Crn *string `json:"crn,omitempty"`
}

// UnmarshalServiceCredentialsSourceServiceIamServiceid unmarshals an instance of ServiceCredentialsSourceServiceIamServiceid from the specified map of raw messages.
func UnmarshalServiceCredentialsSourceServiceIamServiceid(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSourceServiceIamServiceid)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSourceServiceInstance : The source service instance identifier.
type ServiceCredentialsSourceServiceInstance struct {
	// A CRN that uniquely identifies a service credentials source.
	Crn *string `json:"crn,omitempty"`
}

// UnmarshalServiceCredentialsSourceServiceInstance unmarshals an instance of ServiceCredentialsSourceServiceInstance from the specified map of raw messages.
func UnmarshalServiceCredentialsSourceServiceInstance(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSourceServiceInstance)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSourceServiceParameters : Configuration options represented as key-value pairs. Service-defined options are used in the generation of
// credentials for some services. For example, Cloud Object Storage accepts the optional boolean parameter HMAC for
// creating specific kind of credentials.
type ServiceCredentialsSourceServiceParameters struct {
	// An optional platform defined option to reuse an existing IAM Service ID for the role assignment.
	ServiceidCrn *string `json:"serviceid_crn,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of ServiceCredentialsSourceServiceParameters
func (o *ServiceCredentialsSourceServiceParameters) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of ServiceCredentialsSourceServiceParameters
func (o *ServiceCredentialsSourceServiceParameters) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of ServiceCredentialsSourceServiceParameters
func (o *ServiceCredentialsSourceServiceParameters) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of ServiceCredentialsSourceServiceParameters
func (o *ServiceCredentialsSourceServiceParameters) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of ServiceCredentialsSourceServiceParameters
func (o *ServiceCredentialsSourceServiceParameters) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.ServiceidCrn != nil {
		m["serviceid_crn"] = o.ServiceidCrn
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalServiceCredentialsSourceServiceParameters unmarshals an instance of ServiceCredentialsSourceServiceParameters from the specified map of raw messages.
func UnmarshalServiceCredentialsSourceServiceParameters(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSourceServiceParameters)
	err = core.UnmarshalPrimitive(m, "serviceid_crn", &obj.ServiceidCrn)
	if err != nil {
		return
	}
	delete(m, "serviceid_crn")
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = e
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSourceServiceRole : The service-specific custom role. CRN is accepted. The role is assigned as part of an access policy to any
// auto-generated IAM service ID.  If you provide an existing service ID, it is added to the access policy for that ID.
// If a role is not provided, any new service IDs that are autogenerated, will not have an assigned access policy and
// provided service IDs are not changed in any way.  Refer to the service documentation for supported roles.
type ServiceCredentialsSourceServiceRole struct {
	// The service role CRN.
	Crn *string `json:"crn" validate:"required"`
}

// NewServiceCredentialsSourceServiceRole : Instantiate ServiceCredentialsSourceServiceRole (Generic Model Constructor)
func (*SecretsManagerV2) NewServiceCredentialsSourceServiceRole(crn string) (_model *ServiceCredentialsSourceServiceRole, err error) {
	_model = &ServiceCredentialsSourceServiceRole{
		Crn: core.StringPtr(crn),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalServiceCredentialsSourceServiceRole unmarshals an instance of ServiceCredentialsSourceServiceRole from the specified map of raw messages.
func UnmarshalServiceCredentialsSourceServiceRole(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSourceServiceRole)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateConfigurationOptions : The UpdateConfiguration options.
type UpdateConfigurationOptions struct {
	// The name that uniquely identifies a configuration.
	Name *string `json:"name" validate:"required,ne="`

	// JSON Merge-Patch content for update_configuration.
	ConfigurationPatch map[string]interface{} `json:"ConfigurationPatch" validate:"required"`

	// The configuration type of this configuration - use this header to resolve 300 error responses.
	XSmAcceptConfigurationType *string `json:"X-Sm-Accept-Configuration-Type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateConfigurationOptions.XSmAcceptConfigurationType property.
// The configuration type of this configuration - use this header to resolve 300 error responses.
const (
	UpdateConfigurationOptions_XSmAcceptConfigurationType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	UpdateConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	UpdateConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	UpdateConfigurationOptions_XSmAcceptConfigurationType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	UpdateConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	UpdateConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	UpdateConfigurationOptions_XSmAcceptConfigurationType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// NewUpdateConfigurationOptions : Instantiate UpdateConfigurationOptions
func (*SecretsManagerV2) NewUpdateConfigurationOptions(name string, configurationPatch map[string]interface{}) *UpdateConfigurationOptions {
	return &UpdateConfigurationOptions{
		Name:               core.StringPtr(name),
		ConfigurationPatch: configurationPatch,
	}
}

// SetName : Allow user to set Name
func (_options *UpdateConfigurationOptions) SetName(name string) *UpdateConfigurationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetConfigurationPatch : Allow user to set ConfigurationPatch
func (_options *UpdateConfigurationOptions) SetConfigurationPatch(configurationPatch map[string]interface{}) *UpdateConfigurationOptions {
	_options.ConfigurationPatch = configurationPatch
	return _options
}

// SetXSmAcceptConfigurationType : Allow user to set XSmAcceptConfigurationType
func (_options *UpdateConfigurationOptions) SetXSmAcceptConfigurationType(xSmAcceptConfigurationType string) *UpdateConfigurationOptions {
	_options.XSmAcceptConfigurationType = core.StringPtr(xSmAcceptConfigurationType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateConfigurationOptions) SetHeaders(param map[string]string) *UpdateConfigurationOptions {
	options.Headers = param
	return options
}

// UpdateSecretGroupOptions : The UpdateSecretGroup options.
type UpdateSecretGroupOptions struct {
	// The v4 UUID that uniquely identifies your secret group.
	ID *string `json:"id" validate:"required,ne="`

	// The request body to update a secret group.
	SecretGroupPatch map[string]interface{} `json:"SecretGroupPatch" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSecretGroupOptions : Instantiate UpdateSecretGroupOptions
func (*SecretsManagerV2) NewUpdateSecretGroupOptions(id string, secretGroupPatch map[string]interface{}) *UpdateSecretGroupOptions {
	return &UpdateSecretGroupOptions{
		ID:               core.StringPtr(id),
		SecretGroupPatch: secretGroupPatch,
	}
}

// SetID : Allow user to set ID
func (_options *UpdateSecretGroupOptions) SetID(id string) *UpdateSecretGroupOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetSecretGroupPatch : Allow user to set SecretGroupPatch
func (_options *UpdateSecretGroupOptions) SetSecretGroupPatch(secretGroupPatch map[string]interface{}) *UpdateSecretGroupOptions {
	_options.SecretGroupPatch = secretGroupPatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretGroupOptions) SetHeaders(param map[string]string) *UpdateSecretGroupOptions {
	options.Headers = param
	return options
}

// UpdateSecretMetadataOptions : The UpdateSecretMetadata options.
type UpdateSecretMetadataOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	ID *string `json:"id" validate:"required,ne="`

	// JSON Merge-Patch content for update_secret_metadata.
	SecretMetadataPatch map[string]interface{} `json:"SecretMetadataPatch" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSecretMetadataOptions : Instantiate UpdateSecretMetadataOptions
func (*SecretsManagerV2) NewUpdateSecretMetadataOptions(id string, secretMetadataPatch map[string]interface{}) *UpdateSecretMetadataOptions {
	return &UpdateSecretMetadataOptions{
		ID:                  core.StringPtr(id),
		SecretMetadataPatch: secretMetadataPatch,
	}
}

// SetID : Allow user to set ID
func (_options *UpdateSecretMetadataOptions) SetID(id string) *UpdateSecretMetadataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetSecretMetadataPatch : Allow user to set SecretMetadataPatch
func (_options *UpdateSecretMetadataOptions) SetSecretMetadataPatch(secretMetadataPatch map[string]interface{}) *UpdateSecretMetadataOptions {
	_options.SecretMetadataPatch = secretMetadataPatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretMetadataOptions) SetHeaders(param map[string]string) *UpdateSecretMetadataOptions {
	options.Headers = param
	return options
}

// UpdateSecretVersionMetadataOptions : The UpdateSecretVersionMetadata options.
type UpdateSecretVersionMetadataOptions struct {
	// The v4 UUID that uniquely identifies your secret.
	SecretID *string `json:"secret_id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies your secret version. You can use the `current` or `previous` aliases to refer
	// to the current or previous secret version.
	ID *string `json:"id" validate:"required,ne="`

	// JSON Merge-Patch content for update_secret_version_metadata.
	SecretVersionMetadataPatch map[string]interface{} `json:"SecretVersionMetadataPatch" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSecretVersionMetadataOptions : Instantiate UpdateSecretVersionMetadataOptions
func (*SecretsManagerV2) NewUpdateSecretVersionMetadataOptions(secretID string, id string, secretVersionMetadataPatch map[string]interface{}) *UpdateSecretVersionMetadataOptions {
	return &UpdateSecretVersionMetadataOptions{
		SecretID:                   core.StringPtr(secretID),
		ID:                         core.StringPtr(id),
		SecretVersionMetadataPatch: secretVersionMetadataPatch,
	}
}

// SetSecretID : Allow user to set SecretID
func (_options *UpdateSecretVersionMetadataOptions) SetSecretID(secretID string) *UpdateSecretVersionMetadataOptions {
	_options.SecretID = core.StringPtr(secretID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateSecretVersionMetadataOptions) SetID(id string) *UpdateSecretVersionMetadataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetSecretVersionMetadataPatch : Allow user to set SecretVersionMetadataPatch
func (_options *UpdateSecretVersionMetadataOptions) SetSecretVersionMetadataPatch(secretVersionMetadataPatch map[string]interface{}) *UpdateSecretVersionMetadataOptions {
	_options.SecretVersionMetadataPatch = secretVersionMetadataPatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretVersionMetadataOptions) SetHeaders(param map[string]string) *UpdateSecretVersionMetadataOptions {
	options.Headers = param
	return options
}

// VersionAction : The request body to specify the properties of the action to create a secret version.
// Models which "extend" this model:
// - PrivateCertificateVersionActionRevoke
type VersionAction struct {
	// The type of secret version action.
	ActionType *string `json:"action_type,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`
}

// Constants associated with the VersionAction.ActionType property.
// The type of secret version action.
const (
	VersionAction_ActionType_PrivateCertActionRevokeCertificate = "private_cert_action_revoke_certificate"
)

func (*VersionAction) isaVersionAction() bool {
	return true
}

type VersionActionIntf interface {
	isaVersionAction() bool
}

// UnmarshalVersionAction unmarshals an instance of VersionAction from the specified map of raw messages.
func UnmarshalVersionAction(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "action_type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'action_type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'action_type' not found in JSON object")
		return
	}
	if discValue == "private_cert_action_revoke_certificate" {
		err = core.UnmarshalModel(m, "", result, UnmarshalPrivateCertificateVersionActionRevoke)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'action_type': %s", discValue)
	}
	return
}

// ArbitrarySecret : Your arbitrary secret.
// This model "extends" Secret
type ArbitrarySecret struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The secret data that is assigned to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`
}

// Constants associated with the ArbitrarySecret.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ArbitrarySecret_SecretType_Arbitrary          = "arbitrary"
	ArbitrarySecret_SecretType_IamCredentials     = "iam_credentials"
	ArbitrarySecret_SecretType_ImportedCert       = "imported_cert"
	ArbitrarySecret_SecretType_Kv                 = "kv"
	ArbitrarySecret_SecretType_PrivateCert        = "private_cert"
	ArbitrarySecret_SecretType_PublicCert         = "public_cert"
	ArbitrarySecret_SecretType_ServiceCredentials = "service_credentials"
	ArbitrarySecret_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ArbitrarySecret.StateDescription property.
// A text representation of the secret state.
const (
	ArbitrarySecret_StateDescription_Active        = "active"
	ArbitrarySecret_StateDescription_Deactivated   = "deactivated"
	ArbitrarySecret_StateDescription_Destroyed     = "destroyed"
	ArbitrarySecret_StateDescription_PreActivation = "pre_activation"
	ArbitrarySecret_StateDescription_Suspended     = "suspended"
)

func (*ArbitrarySecret) isaSecret() bool {
	return true
}

// UnmarshalArbitrarySecret unmarshals an instance of ArbitrarySecret from the specified map of raw messages.
func UnmarshalArbitrarySecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecret)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ArbitrarySecretMetadata : Properties of the metadata of your arbitrary secret..
// This model "extends" SecretMetadata
type ArbitrarySecretMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`
}

// Constants associated with the ArbitrarySecretMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ArbitrarySecretMetadata_SecretType_Arbitrary          = "arbitrary"
	ArbitrarySecretMetadata_SecretType_IamCredentials     = "iam_credentials"
	ArbitrarySecretMetadata_SecretType_ImportedCert       = "imported_cert"
	ArbitrarySecretMetadata_SecretType_Kv                 = "kv"
	ArbitrarySecretMetadata_SecretType_PrivateCert        = "private_cert"
	ArbitrarySecretMetadata_SecretType_PublicCert         = "public_cert"
	ArbitrarySecretMetadata_SecretType_ServiceCredentials = "service_credentials"
	ArbitrarySecretMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ArbitrarySecretMetadata.StateDescription property.
// A text representation of the secret state.
const (
	ArbitrarySecretMetadata_StateDescription_Active        = "active"
	ArbitrarySecretMetadata_StateDescription_Deactivated   = "deactivated"
	ArbitrarySecretMetadata_StateDescription_Destroyed     = "destroyed"
	ArbitrarySecretMetadata_StateDescription_PreActivation = "pre_activation"
	ArbitrarySecretMetadata_StateDescription_Suspended     = "suspended"
)

func (*ArbitrarySecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalArbitrarySecretMetadata unmarshals an instance of ArbitrarySecretMetadata from the specified map of raw messages.
func UnmarshalArbitrarySecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretMetadata)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ArbitrarySecretMetadataPatch : ArbitrarySecretMetadataPatch struct
// This model "extends" SecretMetadataPatch
type ArbitrarySecretMetadataPatch struct {
	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`
}

func (*ArbitrarySecretMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

// UnmarshalArbitrarySecretMetadataPatch unmarshals an instance of ArbitrarySecretMetadataPatch from the specified map of raw messages.
func UnmarshalArbitrarySecretMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretMetadataPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the ArbitrarySecretMetadataPatch
func (arbitrarySecretMetadataPatch *ArbitrarySecretMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(arbitrarySecretMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// ArbitrarySecretPrototype : ArbitrarySecretPrototype struct
// This model "extends" SecretPrototype
type ArbitrarySecretPrototype struct {
	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret data that is assigned to an `arbitrary` secret.
	Payload *string `json:"payload" validate:"required"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// Constants associated with the ArbitrarySecretPrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ArbitrarySecretPrototype_SecretType_Arbitrary          = "arbitrary"
	ArbitrarySecretPrototype_SecretType_IamCredentials     = "iam_credentials"
	ArbitrarySecretPrototype_SecretType_ImportedCert       = "imported_cert"
	ArbitrarySecretPrototype_SecretType_Kv                 = "kv"
	ArbitrarySecretPrototype_SecretType_PrivateCert        = "private_cert"
	ArbitrarySecretPrototype_SecretType_PublicCert         = "public_cert"
	ArbitrarySecretPrototype_SecretType_ServiceCredentials = "service_credentials"
	ArbitrarySecretPrototype_SecretType_UsernamePassword   = "username_password"
)

// NewArbitrarySecretPrototype : Instantiate ArbitrarySecretPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewArbitrarySecretPrototype(name string, secretType string, payload string) (_model *ArbitrarySecretPrototype, err error) {
	_model = &ArbitrarySecretPrototype{
		Name:       core.StringPtr(name),
		SecretType: core.StringPtr(secretType),
		Payload:    core.StringPtr(payload),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ArbitrarySecretPrototype) isaSecretPrototype() bool {
	return true
}

// UnmarshalArbitrarySecretPrototype unmarshals an instance of ArbitrarySecretPrototype from the specified map of raw messages.
func UnmarshalArbitrarySecretPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretPrototype)
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ArbitrarySecretVersion : Your arbitrary secret version.
// This model "extends" SecretVersion
type ArbitrarySecretVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The secret data that is assigned to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`
}

// Constants associated with the ArbitrarySecretVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ArbitrarySecretVersion_SecretType_Arbitrary          = "arbitrary"
	ArbitrarySecretVersion_SecretType_IamCredentials     = "iam_credentials"
	ArbitrarySecretVersion_SecretType_ImportedCert       = "imported_cert"
	ArbitrarySecretVersion_SecretType_Kv                 = "kv"
	ArbitrarySecretVersion_SecretType_PrivateCert        = "private_cert"
	ArbitrarySecretVersion_SecretType_PublicCert         = "public_cert"
	ArbitrarySecretVersion_SecretType_ServiceCredentials = "service_credentials"
	ArbitrarySecretVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ArbitrarySecretVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	ArbitrarySecretVersion_Alias_Current  = "current"
	ArbitrarySecretVersion_Alias_Previous = "previous"
)

func (*ArbitrarySecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalArbitrarySecretVersion unmarshals an instance of ArbitrarySecretVersion from the specified map of raw messages.
func UnmarshalArbitrarySecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretVersion)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ArbitrarySecretVersionMetadata : Properties of the version metadata of your arbitrary secret.
// This model "extends" SecretVersionMetadata
type ArbitrarySecretVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`
}

// Constants associated with the ArbitrarySecretVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ArbitrarySecretVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	ArbitrarySecretVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	ArbitrarySecretVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	ArbitrarySecretVersionMetadata_SecretType_Kv                 = "kv"
	ArbitrarySecretVersionMetadata_SecretType_PrivateCert        = "private_cert"
	ArbitrarySecretVersionMetadata_SecretType_PublicCert         = "public_cert"
	ArbitrarySecretVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	ArbitrarySecretVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ArbitrarySecretVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	ArbitrarySecretVersionMetadata_Alias_Current  = "current"
	ArbitrarySecretVersionMetadata_Alias_Previous = "previous"
)

func (*ArbitrarySecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalArbitrarySecretVersionMetadata unmarshals an instance of ArbitrarySecretVersionMetadata from the specified map of raw messages.
func UnmarshalArbitrarySecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ArbitrarySecretVersionPrototype : ArbitrarySecretVersionPrototype struct
// This model "extends" SecretVersionPrototype
type ArbitrarySecretVersionPrototype struct {
	// The secret data that is assigned to an `arbitrary` secret.
	Payload *string `json:"payload" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewArbitrarySecretVersionPrototype : Instantiate ArbitrarySecretVersionPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewArbitrarySecretVersionPrototype(payload string) (_model *ArbitrarySecretVersionPrototype, err error) {
	_model = &ArbitrarySecretVersionPrototype{
		Payload: core.StringPtr(payload),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ArbitrarySecretVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalArbitrarySecretVersionPrototype unmarshals an instance of ArbitrarySecretVersionPrototype from the specified map of raw messages.
func UnmarshalArbitrarySecretVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretVersionPrototype)
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CommonRotationPolicy : This field indicates whether Secrets Manager rotates your secrets automatically.
// This model "extends" RotationPolicy
type CommonRotationPolicy struct {
	// This field indicates whether Secrets Manager rotates your secret automatically.
	//
	// The default is `false`. If `auto_rotate` is set to `true` the service rotates your secret based on the defined
	// interval.
	AutoRotate *bool `json:"auto_rotate" validate:"required"`

	// The length of the secret rotation time interval.
	Interval *int64 `json:"interval,omitempty"`

	// The units for the secret rotation time interval.
	Unit *string `json:"unit,omitempty"`
}

// Constants associated with the CommonRotationPolicy.Unit property.
// The units for the secret rotation time interval.
const (
	CommonRotationPolicy_Unit_Day   = "day"
	CommonRotationPolicy_Unit_Month = "month"
)

// NewCommonRotationPolicy : Instantiate CommonRotationPolicy (Generic Model Constructor)
func (*SecretsManagerV2) NewCommonRotationPolicy(autoRotate bool) (_model *CommonRotationPolicy, err error) {
	_model = &CommonRotationPolicy{
		AutoRotate: core.BoolPtr(autoRotate),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*CommonRotationPolicy) isaRotationPolicy() bool {
	return true
}

// UnmarshalCommonRotationPolicy unmarshals an instance of CommonRotationPolicy from the specified map of raw messages.
func UnmarshalCommonRotationPolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CommonRotationPolicy)
	err = core.UnmarshalPrimitive(m, "auto_rotate", &obj.AutoRotate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unit", &obj.Unit)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsConfiguration : Properties that describe a Classic Infrastructure DNS configuration.
// This model "extends" Configuration
type IAMCredentialsConfiguration struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// An IBM Cloud API key that can create and manage service IDs. The API key must be assigned the Editor platform role
	// on the Access Groups Service and the Operator platform role on the IAM Identity Service.  For more information, see
	// the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).
	ApiKey *string `json:"api_key,omitempty"`
}

// Constants associated with the IAMCredentialsConfiguration.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	IAMCredentialsConfiguration_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	IAMCredentialsConfiguration_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	IAMCredentialsConfiguration_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	IAMCredentialsConfiguration_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	IAMCredentialsConfiguration_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	IAMCredentialsConfiguration_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	IAMCredentialsConfiguration_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the IAMCredentialsConfiguration.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	IAMCredentialsConfiguration_SecretType_Arbitrary          = "arbitrary"
	IAMCredentialsConfiguration_SecretType_IamCredentials     = "iam_credentials"
	IAMCredentialsConfiguration_SecretType_ImportedCert       = "imported_cert"
	IAMCredentialsConfiguration_SecretType_Kv                 = "kv"
	IAMCredentialsConfiguration_SecretType_PrivateCert        = "private_cert"
	IAMCredentialsConfiguration_SecretType_PublicCert         = "public_cert"
	IAMCredentialsConfiguration_SecretType_ServiceCredentials = "service_credentials"
	IAMCredentialsConfiguration_SecretType_UsernamePassword   = "username_password"
)

func (*IAMCredentialsConfiguration) isaConfiguration() bool {
	return true
}

// UnmarshalIAMCredentialsConfiguration unmarshals an instance of IAMCredentialsConfiguration from the specified map of raw messages.
func UnmarshalIAMCredentialsConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsConfiguration)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsConfigurationMetadata : Your IAMCredentials Configuration metadata properties.
// This model "extends" ConfigurationMetadata
type IAMCredentialsConfigurationMetadata struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`
}

// Constants associated with the IAMCredentialsConfigurationMetadata.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	IAMCredentialsConfigurationMetadata_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	IAMCredentialsConfigurationMetadata_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	IAMCredentialsConfigurationMetadata_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	IAMCredentialsConfigurationMetadata_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	IAMCredentialsConfigurationMetadata_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	IAMCredentialsConfigurationMetadata_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	IAMCredentialsConfigurationMetadata_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the IAMCredentialsConfigurationMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	IAMCredentialsConfigurationMetadata_SecretType_Arbitrary          = "arbitrary"
	IAMCredentialsConfigurationMetadata_SecretType_IamCredentials     = "iam_credentials"
	IAMCredentialsConfigurationMetadata_SecretType_ImportedCert       = "imported_cert"
	IAMCredentialsConfigurationMetadata_SecretType_Kv                 = "kv"
	IAMCredentialsConfigurationMetadata_SecretType_PrivateCert        = "private_cert"
	IAMCredentialsConfigurationMetadata_SecretType_PublicCert         = "public_cert"
	IAMCredentialsConfigurationMetadata_SecretType_ServiceCredentials = "service_credentials"
	IAMCredentialsConfigurationMetadata_SecretType_UsernamePassword   = "username_password"
)

func (*IAMCredentialsConfigurationMetadata) isaConfigurationMetadata() bool {
	return true
}

// UnmarshalIAMCredentialsConfigurationMetadata unmarshals an instance of IAMCredentialsConfigurationMetadata from the specified map of raw messages.
func UnmarshalIAMCredentialsConfigurationMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsConfigurationMetadata)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsConfigurationPatch : The configuration update of the IAM Credentials engine.
// This model "extends" ConfigurationPatch
type IAMCredentialsConfigurationPatch struct {
	// An IBM Cloud API key that can create and manage service IDs. The API key must be assigned the Editor platform role
	// on the Access Groups Service and the Operator platform role on the IAM Identity Service.  For more information, see
	// the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).
	ApiKey *string `json:"api_key" validate:"required"`
}

// NewIAMCredentialsConfigurationPatch : Instantiate IAMCredentialsConfigurationPatch (Generic Model Constructor)
func (*SecretsManagerV2) NewIAMCredentialsConfigurationPatch(apiKey string) (_model *IAMCredentialsConfigurationPatch, err error) {
	_model = &IAMCredentialsConfigurationPatch{
		ApiKey: core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*IAMCredentialsConfigurationPatch) isaConfigurationPatch() bool {
	return true
}

// UnmarshalIAMCredentialsConfigurationPatch unmarshals an instance of IAMCredentialsConfigurationPatch from the specified map of raw messages.
func UnmarshalIAMCredentialsConfigurationPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsConfigurationPatch)
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the IAMCredentialsConfigurationPatch
func (iAMCredentialsConfigurationPatch *IAMCredentialsConfigurationPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(iAMCredentialsConfigurationPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// IAMCredentialsConfigurationPrototype : IAMCredentialsConfigurationPrototype struct
// This model "extends" ConfigurationPrototype
type IAMCredentialsConfigurationPrototype struct {
	// A human-readable unique name to assign to your configuration.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.
	Name *string `json:"name" validate:"required"`

	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The API key that is used to set the iam_credentials engine.
	ApiKey *string `json:"api_key" validate:"required"`
}

// Constants associated with the IAMCredentialsConfigurationPrototype.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	IAMCredentialsConfigurationPrototype_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	IAMCredentialsConfigurationPrototype_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	IAMCredentialsConfigurationPrototype_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	IAMCredentialsConfigurationPrototype_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	IAMCredentialsConfigurationPrototype_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	IAMCredentialsConfigurationPrototype_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	IAMCredentialsConfigurationPrototype_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// NewIAMCredentialsConfigurationPrototype : Instantiate IAMCredentialsConfigurationPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewIAMCredentialsConfigurationPrototype(name string, configType string, apiKey string) (_model *IAMCredentialsConfigurationPrototype, err error) {
	_model = &IAMCredentialsConfigurationPrototype{
		Name:       core.StringPtr(name),
		ConfigType: core.StringPtr(configType),
		ApiKey:     core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*IAMCredentialsConfigurationPrototype) isaConfigurationPrototype() bool {
	return true
}

// UnmarshalIAMCredentialsConfigurationPrototype unmarshals an instance of IAMCredentialsConfigurationPrototype from the specified map of raw messages.
func UnmarshalIAMCredentialsConfigurationPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsConfigurationPrototype)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsSecret : Your IAM credentials secret.
// This model "extends" Secret
type IAMCredentialsSecret struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl" validate:"required"`

	// Access Groups that you can use for an `iam_credentials` secret.
	//
	// Up to 10 Access Groups can be used for each secret.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The ID of the API key that is generated for this secret.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// Indicates whether an `iam_credentials` secret was created with a static service ID.
	//
	// If it is set to `true`, the service ID for the secret was provided by the user at secret creation. If it is set to
	// `false`, the service ID was generated by Secrets Manager.
	ServiceIdIsStatic *bool `json:"service_id_is_static,omitempty"`

	// (IAM credentials) This parameter indicates whether to reuse the service ID and API key for future read operations.
	//
	// If it is set to `true`, the service reuses the current credentials. If it is set to `false`, a new service ID and
	// API key are generated each time that the secret is read or accessed.
	ReuseApiKey *bool `json:"reuse_api_key" validate:"required"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease, the API key is deleted automatically. See the `time-to-live` field to
	// understand the duration of the lease. If you want to continue to use the same API key for future read operations,
	// see the `reuse_api_key` field.
	ApiKey *string `json:"api_key,omitempty"`
}

// Constants associated with the IAMCredentialsSecret.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	IAMCredentialsSecret_SecretType_Arbitrary          = "arbitrary"
	IAMCredentialsSecret_SecretType_IamCredentials     = "iam_credentials"
	IAMCredentialsSecret_SecretType_ImportedCert       = "imported_cert"
	IAMCredentialsSecret_SecretType_Kv                 = "kv"
	IAMCredentialsSecret_SecretType_PrivateCert        = "private_cert"
	IAMCredentialsSecret_SecretType_PublicCert         = "public_cert"
	IAMCredentialsSecret_SecretType_ServiceCredentials = "service_credentials"
	IAMCredentialsSecret_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the IAMCredentialsSecret.StateDescription property.
// A text representation of the secret state.
const (
	IAMCredentialsSecret_StateDescription_Active        = "active"
	IAMCredentialsSecret_StateDescription_Deactivated   = "deactivated"
	IAMCredentialsSecret_StateDescription_Destroyed     = "destroyed"
	IAMCredentialsSecret_StateDescription_PreActivation = "pre_activation"
	IAMCredentialsSecret_StateDescription_Suspended     = "suspended"
)

func (*IAMCredentialsSecret) isaSecret() bool {
	return true
}

// UnmarshalIAMCredentialsSecret unmarshals an instance of IAMCredentialsSecret from the specified map of raw messages.
func UnmarshalIAMCredentialsSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsSecret)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "access_groups", &obj.AccessGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.ApiKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id_is_static", &obj.ServiceIdIsStatic)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsSecretMetadata : Properties of the metadata of your IAM credentials secret.
// This model "extends" SecretMetadata
type IAMCredentialsSecretMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl" validate:"required"`

	// Access Groups that you can use for an `iam_credentials` secret.
	//
	// Up to 10 Access Groups can be used for each secret.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The ID of the API key that is generated for this secret.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// Indicates whether an `iam_credentials` secret was created with a static service ID.
	//
	// If it is set to `true`, the service ID for the secret was provided by the user at secret creation. If it is set to
	// `false`, the service ID was generated by Secrets Manager.
	ServiceIdIsStatic *bool `json:"service_id_is_static,omitempty"`

	// (IAM credentials) This parameter indicates whether to reuse the service ID and API key for future read operations.
	//
	// If it is set to `true`, the service reuses the current credentials. If it is set to `false`, a new service ID and
	// API key are generated each time that the secret is read or accessed.
	ReuseApiKey *bool `json:"reuse_api_key" validate:"required"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`
}

// Constants associated with the IAMCredentialsSecretMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	IAMCredentialsSecretMetadata_SecretType_Arbitrary          = "arbitrary"
	IAMCredentialsSecretMetadata_SecretType_IamCredentials     = "iam_credentials"
	IAMCredentialsSecretMetadata_SecretType_ImportedCert       = "imported_cert"
	IAMCredentialsSecretMetadata_SecretType_Kv                 = "kv"
	IAMCredentialsSecretMetadata_SecretType_PrivateCert        = "private_cert"
	IAMCredentialsSecretMetadata_SecretType_PublicCert         = "public_cert"
	IAMCredentialsSecretMetadata_SecretType_ServiceCredentials = "service_credentials"
	IAMCredentialsSecretMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the IAMCredentialsSecretMetadata.StateDescription property.
// A text representation of the secret state.
const (
	IAMCredentialsSecretMetadata_StateDescription_Active        = "active"
	IAMCredentialsSecretMetadata_StateDescription_Deactivated   = "deactivated"
	IAMCredentialsSecretMetadata_StateDescription_Destroyed     = "destroyed"
	IAMCredentialsSecretMetadata_StateDescription_PreActivation = "pre_activation"
	IAMCredentialsSecretMetadata_StateDescription_Suspended     = "suspended"
)

func (*IAMCredentialsSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalIAMCredentialsSecretMetadata unmarshals an instance of IAMCredentialsSecretMetadata from the specified map of raw messages.
func UnmarshalIAMCredentialsSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsSecretMetadata)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "access_groups", &obj.AccessGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.ApiKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id_is_static", &obj.ServiceIdIsStatic)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsSecretMetadataPatch : IAMCredentialsSecretMetadataPatch struct
// This model "extends" SecretMetadataPatch
type IAMCredentialsSecretMetadataPatch struct {
	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`
}

func (*IAMCredentialsSecretMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

// UnmarshalIAMCredentialsSecretMetadataPatch unmarshals an instance of IAMCredentialsSecretMetadataPatch from the specified map of raw messages.
func UnmarshalIAMCredentialsSecretMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsSecretMetadataPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the IAMCredentialsSecretMetadataPatch
func (iAMCredentialsSecretMetadataPatch *IAMCredentialsSecretMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(iAMCredentialsSecretMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// IAMCredentialsSecretPrototype : IAMCredentialsSecretPrototype struct
// This model "extends" SecretPrototype
type IAMCredentialsSecretPrototype struct {
	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl" validate:"required"`

	// Access Groups that you can use for an `iam_credentials` secret.
	//
	// Up to 10 Access Groups can be used for each secret.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// (IAM credentials) This parameter indicates whether to reuse the service ID and API key for future read operations.
	//
	// If it is set to `true`, the service reuses the current credentials. If it is set to `false`, a new service ID and
	// API key are generated each time that the secret is read or accessed.
	ReuseApiKey *bool `json:"reuse_api_key" validate:"required"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// Constants associated with the IAMCredentialsSecretPrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	IAMCredentialsSecretPrototype_SecretType_Arbitrary          = "arbitrary"
	IAMCredentialsSecretPrototype_SecretType_IamCredentials     = "iam_credentials"
	IAMCredentialsSecretPrototype_SecretType_ImportedCert       = "imported_cert"
	IAMCredentialsSecretPrototype_SecretType_Kv                 = "kv"
	IAMCredentialsSecretPrototype_SecretType_PrivateCert        = "private_cert"
	IAMCredentialsSecretPrototype_SecretType_PublicCert         = "public_cert"
	IAMCredentialsSecretPrototype_SecretType_ServiceCredentials = "service_credentials"
	IAMCredentialsSecretPrototype_SecretType_UsernamePassword   = "username_password"
)

// NewIAMCredentialsSecretPrototype : Instantiate IAMCredentialsSecretPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewIAMCredentialsSecretPrototype(secretType string, name string, ttl string, reuseApiKey bool) (_model *IAMCredentialsSecretPrototype, err error) {
	_model = &IAMCredentialsSecretPrototype{
		SecretType:  core.StringPtr(secretType),
		Name:        core.StringPtr(name),
		TTL:         core.StringPtr(ttl),
		ReuseApiKey: core.BoolPtr(reuseApiKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*IAMCredentialsSecretPrototype) isaSecretPrototype() bool {
	return true
}

// UnmarshalIAMCredentialsSecretPrototype unmarshals an instance of IAMCredentialsSecretPrototype from the specified map of raw messages.
func UnmarshalIAMCredentialsSecretPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsSecretPrototype)
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "access_groups", &obj.AccessGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsSecretRestoreFromVersionPrototype : IAMCredentialsSecretRestoreFromVersionPrototype struct
// This model "extends" SecretVersionPrototype
type IAMCredentialsSecretRestoreFromVersionPrototype struct {
	// A v4 UUID identifier, or `current` or `previous` secret version aliases.
	RestoreFromVersion *string `json:"restore_from_version" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewIAMCredentialsSecretRestoreFromVersionPrototype : Instantiate IAMCredentialsSecretRestoreFromVersionPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewIAMCredentialsSecretRestoreFromVersionPrototype(restoreFromVersion string) (_model *IAMCredentialsSecretRestoreFromVersionPrototype, err error) {
	_model = &IAMCredentialsSecretRestoreFromVersionPrototype{
		RestoreFromVersion: core.StringPtr(restoreFromVersion),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*IAMCredentialsSecretRestoreFromVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalIAMCredentialsSecretRestoreFromVersionPrototype unmarshals an instance of IAMCredentialsSecretRestoreFromVersionPrototype from the specified map of raw messages.
func UnmarshalIAMCredentialsSecretRestoreFromVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsSecretRestoreFromVersionPrototype)
	err = core.UnmarshalPrimitive(m, "restore_from_version", &obj.RestoreFromVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsSecretVersion : Your IAM credentials version.
// This model "extends" SecretVersion
type IAMCredentialsSecretVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The ID of the API key that is generated for this secret.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease, the API key is deleted automatically. See the `time-to-live` field to
	// understand the duration of the lease. If you want to continue to use the same API key for future read operations,
	// see the `reuse_api_key` field.
	ApiKey *string `json:"api_key,omitempty"`
}

// Constants associated with the IAMCredentialsSecretVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	IAMCredentialsSecretVersion_SecretType_Arbitrary          = "arbitrary"
	IAMCredentialsSecretVersion_SecretType_IamCredentials     = "iam_credentials"
	IAMCredentialsSecretVersion_SecretType_ImportedCert       = "imported_cert"
	IAMCredentialsSecretVersion_SecretType_Kv                 = "kv"
	IAMCredentialsSecretVersion_SecretType_PrivateCert        = "private_cert"
	IAMCredentialsSecretVersion_SecretType_PublicCert         = "public_cert"
	IAMCredentialsSecretVersion_SecretType_ServiceCredentials = "service_credentials"
	IAMCredentialsSecretVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the IAMCredentialsSecretVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	IAMCredentialsSecretVersion_Alias_Current  = "current"
	IAMCredentialsSecretVersion_Alias_Previous = "previous"
)

func (*IAMCredentialsSecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalIAMCredentialsSecretVersion unmarshals an instance of IAMCredentialsSecretVersion from the specified map of raw messages.
func UnmarshalIAMCredentialsSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsSecretVersion)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.ApiKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsSecretVersionMetadata : Properties of the version metadata of your IAM credentials secret.
// This model "extends" SecretVersionMetadata
type IAMCredentialsSecretVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The ID of the API key that is generated for this secret.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation, and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`
}

// Constants associated with the IAMCredentialsSecretVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	IAMCredentialsSecretVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	IAMCredentialsSecretVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	IAMCredentialsSecretVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	IAMCredentialsSecretVersionMetadata_SecretType_Kv                 = "kv"
	IAMCredentialsSecretVersionMetadata_SecretType_PrivateCert        = "private_cert"
	IAMCredentialsSecretVersionMetadata_SecretType_PublicCert         = "public_cert"
	IAMCredentialsSecretVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	IAMCredentialsSecretVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the IAMCredentialsSecretVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	IAMCredentialsSecretVersionMetadata_Alias_Current  = "current"
	IAMCredentialsSecretVersionMetadata_Alias_Previous = "previous"
)

func (*IAMCredentialsSecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalIAMCredentialsSecretVersionMetadata unmarshals an instance of IAMCredentialsSecretVersionMetadata from the specified map of raw messages.
func UnmarshalIAMCredentialsSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.ApiKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IAMCredentialsSecretVersionPrototype : IAMCredentialsSecretVersionPrototype struct
// This model "extends" SecretVersionPrototype
type IAMCredentialsSecretVersionPrototype struct {
	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

func (*IAMCredentialsSecretVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalIAMCredentialsSecretVersionPrototype unmarshals an instance of IAMCredentialsSecretVersionPrototype from the specified map of raw messages.
func UnmarshalIAMCredentialsSecretVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IAMCredentialsSecretVersionPrototype)
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportedCertificate : Your imported certificate.
// This model "extends" Secret
type ImportedCertificate struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The identifier for the cryptographic algorithm that is used by the issuing certificate authority to sign a
	// certificate.
	SigningAlgorithm *string `json:"signing_algorithm" validate:"required"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The Common Name (CN) represents the server name protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date" validate:"required"`

	// Indicates whether the certificate was imported with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included" validate:"required"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer" validate:"required"`

	// The identifier for the cryptographic algorithm used to generate the public key that is associated with the
	// certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// Indicates whether the certificate was imported with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included" validate:"required"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number" validate:"required"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity" validate:"required"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`
}

// Constants associated with the ImportedCertificate.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ImportedCertificate_SecretType_Arbitrary          = "arbitrary"
	ImportedCertificate_SecretType_IamCredentials     = "iam_credentials"
	ImportedCertificate_SecretType_ImportedCert       = "imported_cert"
	ImportedCertificate_SecretType_Kv                 = "kv"
	ImportedCertificate_SecretType_PrivateCert        = "private_cert"
	ImportedCertificate_SecretType_PublicCert         = "public_cert"
	ImportedCertificate_SecretType_ServiceCredentials = "service_credentials"
	ImportedCertificate_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ImportedCertificate.StateDescription property.
// A text representation of the secret state.
const (
	ImportedCertificate_StateDescription_Active        = "active"
	ImportedCertificate_StateDescription_Deactivated   = "deactivated"
	ImportedCertificate_StateDescription_Destroyed     = "destroyed"
	ImportedCertificate_StateDescription_PreActivation = "pre_activation"
	ImportedCertificate_StateDescription_Suspended     = "suspended"
)

func (*ImportedCertificate) isaSecret() bool {
	return true
}

// UnmarshalImportedCertificate unmarshals an instance of ImportedCertificate from the specified map of raw messages.
func UnmarshalImportedCertificate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportedCertificate)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_algorithm", &obj.SigningAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_included", &obj.IntermediateIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_included", &obj.PrivateKeyIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportedCertificateMetadata : Properties of the secret metadata of your imported certificate.
// This model "extends" SecretMetadata
type ImportedCertificateMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The identifier for the cryptographic algorithm that is used by the issuing certificate authority to sign a
	// certificate.
	SigningAlgorithm *string `json:"signing_algorithm" validate:"required"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The Common Name (CN) represents the server name protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date" validate:"required"`

	// Indicates whether the certificate was imported with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included" validate:"required"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer" validate:"required"`

	// The identifier for the cryptographic algorithm used to generate the public key that is associated with the
	// certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// Indicates whether the certificate was imported with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included" validate:"required"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number" validate:"required"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity" validate:"required"`
}

// Constants associated with the ImportedCertificateMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ImportedCertificateMetadata_SecretType_Arbitrary          = "arbitrary"
	ImportedCertificateMetadata_SecretType_IamCredentials     = "iam_credentials"
	ImportedCertificateMetadata_SecretType_ImportedCert       = "imported_cert"
	ImportedCertificateMetadata_SecretType_Kv                 = "kv"
	ImportedCertificateMetadata_SecretType_PrivateCert        = "private_cert"
	ImportedCertificateMetadata_SecretType_PublicCert         = "public_cert"
	ImportedCertificateMetadata_SecretType_ServiceCredentials = "service_credentials"
	ImportedCertificateMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ImportedCertificateMetadata.StateDescription property.
// A text representation of the secret state.
const (
	ImportedCertificateMetadata_StateDescription_Active        = "active"
	ImportedCertificateMetadata_StateDescription_Deactivated   = "deactivated"
	ImportedCertificateMetadata_StateDescription_Destroyed     = "destroyed"
	ImportedCertificateMetadata_StateDescription_PreActivation = "pre_activation"
	ImportedCertificateMetadata_StateDescription_Suspended     = "suspended"
)

func (*ImportedCertificateMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalImportedCertificateMetadata unmarshals an instance of ImportedCertificateMetadata from the specified map of raw messages.
func UnmarshalImportedCertificateMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportedCertificateMetadata)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_algorithm", &obj.SigningAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_included", &obj.IntermediateIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_included", &obj.PrivateKeyIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportedCertificateMetadataPatch : ImportedCertificateMetadataPatch struct
// This model "extends" SecretMetadataPatch
type ImportedCertificateMetadataPatch struct {
	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`
}

func (*ImportedCertificateMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

// UnmarshalImportedCertificateMetadataPatch unmarshals an instance of ImportedCertificateMetadataPatch from the specified map of raw messages.
func UnmarshalImportedCertificateMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportedCertificateMetadataPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the ImportedCertificateMetadataPatch
func (importedCertificateMetadataPatch *ImportedCertificateMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(importedCertificateMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// ImportedCertificatePrototype : ImportedCertificatePrototype struct
// This model "extends" SecretPrototype
type ImportedCertificatePrototype struct {
	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// Constants associated with the ImportedCertificatePrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ImportedCertificatePrototype_SecretType_Arbitrary          = "arbitrary"
	ImportedCertificatePrototype_SecretType_IamCredentials     = "iam_credentials"
	ImportedCertificatePrototype_SecretType_ImportedCert       = "imported_cert"
	ImportedCertificatePrototype_SecretType_Kv                 = "kv"
	ImportedCertificatePrototype_SecretType_PrivateCert        = "private_cert"
	ImportedCertificatePrototype_SecretType_PublicCert         = "public_cert"
	ImportedCertificatePrototype_SecretType_ServiceCredentials = "service_credentials"
	ImportedCertificatePrototype_SecretType_UsernamePassword   = "username_password"
)

// NewImportedCertificatePrototype : Instantiate ImportedCertificatePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewImportedCertificatePrototype(secretType string, name string, certificate string) (_model *ImportedCertificatePrototype, err error) {
	_model = &ImportedCertificatePrototype{
		SecretType:  core.StringPtr(secretType),
		Name:        core.StringPtr(name),
		Certificate: core.StringPtr(certificate),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ImportedCertificatePrototype) isaSecretPrototype() bool {
	return true
}

// UnmarshalImportedCertificatePrototype unmarshals an instance of ImportedCertificatePrototype from the specified map of raw messages.
func UnmarshalImportedCertificatePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportedCertificatePrototype)
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportedCertificateVersion : Versions of your imported certificate.
// This model "extends" SecretVersion
type ImportedCertificateVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date" validate:"required"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number" validate:"required"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity" validate:"required"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`
}

// Constants associated with the ImportedCertificateVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ImportedCertificateVersion_SecretType_Arbitrary          = "arbitrary"
	ImportedCertificateVersion_SecretType_IamCredentials     = "iam_credentials"
	ImportedCertificateVersion_SecretType_ImportedCert       = "imported_cert"
	ImportedCertificateVersion_SecretType_Kv                 = "kv"
	ImportedCertificateVersion_SecretType_PrivateCert        = "private_cert"
	ImportedCertificateVersion_SecretType_PublicCert         = "public_cert"
	ImportedCertificateVersion_SecretType_ServiceCredentials = "service_credentials"
	ImportedCertificateVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ImportedCertificateVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	ImportedCertificateVersion_Alias_Current  = "current"
	ImportedCertificateVersion_Alias_Previous = "previous"
)

func (*ImportedCertificateVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalImportedCertificateVersion unmarshals an instance of ImportedCertificateVersion from the specified map of raw messages.
func UnmarshalImportedCertificateVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportedCertificateVersion)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportedCertificateVersionMetadata : Properties of the version metadata of your imported certificate.
// This model "extends" SecretVersionMetadata
type ImportedCertificateVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date" validate:"required"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number" validate:"required"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity" validate:"required"`
}

// Constants associated with the ImportedCertificateVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ImportedCertificateVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	ImportedCertificateVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	ImportedCertificateVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	ImportedCertificateVersionMetadata_SecretType_Kv                 = "kv"
	ImportedCertificateVersionMetadata_SecretType_PrivateCert        = "private_cert"
	ImportedCertificateVersionMetadata_SecretType_PublicCert         = "public_cert"
	ImportedCertificateVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	ImportedCertificateVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ImportedCertificateVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	ImportedCertificateVersionMetadata_Alias_Current  = "current"
	ImportedCertificateVersionMetadata_Alias_Previous = "previous"
)

func (*ImportedCertificateVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalImportedCertificateVersionMetadata unmarshals an instance of ImportedCertificateVersionMetadata from the specified map of raw messages.
func UnmarshalImportedCertificateVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportedCertificateVersionMetadata)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportedCertificateVersionPrototype : ImportedCertificateVersionPrototype struct
// This model "extends" SecretVersionPrototype
type ImportedCertificateVersionPrototype struct {
	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewImportedCertificateVersionPrototype : Instantiate ImportedCertificateVersionPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewImportedCertificateVersionPrototype(certificate string) (_model *ImportedCertificateVersionPrototype, err error) {
	_model = &ImportedCertificateVersionPrototype{
		Certificate: core.StringPtr(certificate),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ImportedCertificateVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalImportedCertificateVersionPrototype unmarshals an instance of ImportedCertificateVersionPrototype from the specified map of raw messages.
func UnmarshalImportedCertificateVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportedCertificateVersionPrototype)
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KVSecret : Your key-value secret.
// This model "extends" Secret
type KVSecret struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The payload data of a key-value secret.
	Data map[string]interface{} `json:"data" validate:"required"`
}

// Constants associated with the KVSecret.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	KVSecret_SecretType_Arbitrary          = "arbitrary"
	KVSecret_SecretType_IamCredentials     = "iam_credentials"
	KVSecret_SecretType_ImportedCert       = "imported_cert"
	KVSecret_SecretType_Kv                 = "kv"
	KVSecret_SecretType_PrivateCert        = "private_cert"
	KVSecret_SecretType_PublicCert         = "public_cert"
	KVSecret_SecretType_ServiceCredentials = "service_credentials"
	KVSecret_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the KVSecret.StateDescription property.
// A text representation of the secret state.
const (
	KVSecret_StateDescription_Active        = "active"
	KVSecret_StateDescription_Deactivated   = "deactivated"
	KVSecret_StateDescription_Destroyed     = "destroyed"
	KVSecret_StateDescription_PreActivation = "pre_activation"
	KVSecret_StateDescription_Suspended     = "suspended"
)

func (*KVSecret) isaSecret() bool {
	return true
}

// UnmarshalKVSecret unmarshals an instance of KVSecret from the specified map of raw messages.
func UnmarshalKVSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KVSecret)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KVSecretMetadata : Properties of the metadata of your key-value secret metadata.
// This model "extends" SecretMetadata
type KVSecretMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`
}

// Constants associated with the KVSecretMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	KVSecretMetadata_SecretType_Arbitrary          = "arbitrary"
	KVSecretMetadata_SecretType_IamCredentials     = "iam_credentials"
	KVSecretMetadata_SecretType_ImportedCert       = "imported_cert"
	KVSecretMetadata_SecretType_Kv                 = "kv"
	KVSecretMetadata_SecretType_PrivateCert        = "private_cert"
	KVSecretMetadata_SecretType_PublicCert         = "public_cert"
	KVSecretMetadata_SecretType_ServiceCredentials = "service_credentials"
	KVSecretMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the KVSecretMetadata.StateDescription property.
// A text representation of the secret state.
const (
	KVSecretMetadata_StateDescription_Active        = "active"
	KVSecretMetadata_StateDescription_Deactivated   = "deactivated"
	KVSecretMetadata_StateDescription_Destroyed     = "destroyed"
	KVSecretMetadata_StateDescription_PreActivation = "pre_activation"
	KVSecretMetadata_StateDescription_Suspended     = "suspended"
)

func (*KVSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalKVSecretMetadata unmarshals an instance of KVSecretMetadata from the specified map of raw messages.
func UnmarshalKVSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KVSecretMetadata)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KVSecretMetadataPatch : KVSecretMetadataPatch struct
// This model "extends" SecretMetadataPatch
type KVSecretMetadataPatch struct {
	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`
}

func (*KVSecretMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

// UnmarshalKVSecretMetadataPatch unmarshals an instance of KVSecretMetadataPatch from the specified map of raw messages.
func UnmarshalKVSecretMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KVSecretMetadataPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the KVSecretMetadataPatch
func (kVSecretMetadataPatch *KVSecretMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(kVSecretMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// KVSecretPrototype : KVSecretPrototype struct
// This model "extends" SecretPrototype
type KVSecretPrototype struct {
	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The payload data of a key-value secret.
	Data map[string]interface{} `json:"data" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// Constants associated with the KVSecretPrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	KVSecretPrototype_SecretType_Arbitrary          = "arbitrary"
	KVSecretPrototype_SecretType_IamCredentials     = "iam_credentials"
	KVSecretPrototype_SecretType_ImportedCert       = "imported_cert"
	KVSecretPrototype_SecretType_Kv                 = "kv"
	KVSecretPrototype_SecretType_PrivateCert        = "private_cert"
	KVSecretPrototype_SecretType_PublicCert         = "public_cert"
	KVSecretPrototype_SecretType_ServiceCredentials = "service_credentials"
	KVSecretPrototype_SecretType_UsernamePassword   = "username_password"
)

// NewKVSecretPrototype : Instantiate KVSecretPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewKVSecretPrototype(secretType string, name string, data map[string]interface{}) (_model *KVSecretPrototype, err error) {
	_model = &KVSecretPrototype{
		SecretType: core.StringPtr(secretType),
		Name:       core.StringPtr(name),
		Data:       data,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KVSecretPrototype) isaSecretPrototype() bool {
	return true
}

// UnmarshalKVSecretPrototype unmarshals an instance of KVSecretPrototype from the specified map of raw messages.
func UnmarshalKVSecretPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KVSecretPrototype)
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KVSecretVersion : Your key-value secret version.
// This model "extends" SecretVersion
type KVSecretVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The payload data of a key-value secret.
	Data map[string]interface{} `json:"data" validate:"required"`
}

// Constants associated with the KVSecretVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	KVSecretVersion_SecretType_Arbitrary          = "arbitrary"
	KVSecretVersion_SecretType_IamCredentials     = "iam_credentials"
	KVSecretVersion_SecretType_ImportedCert       = "imported_cert"
	KVSecretVersion_SecretType_Kv                 = "kv"
	KVSecretVersion_SecretType_PrivateCert        = "private_cert"
	KVSecretVersion_SecretType_PublicCert         = "public_cert"
	KVSecretVersion_SecretType_ServiceCredentials = "service_credentials"
	KVSecretVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the KVSecretVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	KVSecretVersion_Alias_Current  = "current"
	KVSecretVersion_Alias_Previous = "previous"
)

func (*KVSecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalKVSecretVersion unmarshals an instance of KVSecretVersion from the specified map of raw messages.
func UnmarshalKVSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KVSecretVersion)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KVSecretVersionMetadata : Properties of the version metadata of your key-value secret.
// This model "extends" SecretVersionMetadata
type KVSecretVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`
}

// Constants associated with the KVSecretVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	KVSecretVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	KVSecretVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	KVSecretVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	KVSecretVersionMetadata_SecretType_Kv                 = "kv"
	KVSecretVersionMetadata_SecretType_PrivateCert        = "private_cert"
	KVSecretVersionMetadata_SecretType_PublicCert         = "public_cert"
	KVSecretVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	KVSecretVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the KVSecretVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	KVSecretVersionMetadata_Alias_Current  = "current"
	KVSecretVersionMetadata_Alias_Previous = "previous"
)

func (*KVSecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalKVSecretVersionMetadata unmarshals an instance of KVSecretVersionMetadata from the specified map of raw messages.
func UnmarshalKVSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KVSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KVSecretVersionPrototype : KVSecretVersionPrototype struct
// This model "extends" SecretVersionPrototype
type KVSecretVersionPrototype struct {
	// The payload data of a key-value secret.
	Data map[string]interface{} `json:"data" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewKVSecretVersionPrototype : Instantiate KVSecretVersionPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewKVSecretVersionPrototype(data map[string]interface{}) (_model *KVSecretVersionPrototype, err error) {
	_model = &KVSecretVersionPrototype{
		Data: data,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KVSecretVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalKVSecretVersionPrototype unmarshals an instance of KVSecretVersionPrototype from the specified map of raw messages.
func UnmarshalKVSecretVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KVSecretVersionPrototype)
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificate : Your private certificate.
// This model "extends" Secret
type PrivateCertificate struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The identifier for the cryptographic algorithm that is used by the issuing certificate authority to sign a
	// certificate.
	SigningAlgorithm *string `json:"signing_algorithm" validate:"required"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The intermediate certificate authority that signed this certificate.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template" validate:"required"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date" validate:"required"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer" validate:"required"`

	// The identifier for the cryptographic algorithm used to generate the public key that is associated with the
	// certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number" validate:"required"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity" validate:"required"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`

	// The date and time that the certificate was revoked. The date format follows `RFC 3339`.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key" validate:"required"`

	// The PEM-encoded certificate of the certificate authority that signed and issued this certificate.
	IssuingCa *string `json:"issuing_ca,omitempty"`

	// The chain of certificate authorities that are associated with the certificate.
	CaChain []string `json:"ca_chain,omitempty"`
}

// Constants associated with the PrivateCertificate.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificate_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificate_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificate_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificate_SecretType_Kv                 = "kv"
	PrivateCertificate_SecretType_PrivateCert        = "private_cert"
	PrivateCertificate_SecretType_PublicCert         = "public_cert"
	PrivateCertificate_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificate_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificate.StateDescription property.
// A text representation of the secret state.
const (
	PrivateCertificate_StateDescription_Active        = "active"
	PrivateCertificate_StateDescription_Deactivated   = "deactivated"
	PrivateCertificate_StateDescription_Destroyed     = "destroyed"
	PrivateCertificate_StateDescription_PreActivation = "pre_activation"
	PrivateCertificate_StateDescription_Suspended     = "suspended"
)

func (*PrivateCertificate) isaSecret() bool {
	return true
}

// UnmarshalPrivateCertificate unmarshals an instance of PrivateCertificate from the specified map of raw messages.
func UnmarshalPrivateCertificate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificate)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_algorithm", &obj.SigningAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_template", &obj.CertificateTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_seconds", &obj.RevocationTimeSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_rfc3339", &obj.RevocationTimeRfc3339)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_ca", &obj.IssuingCa)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ca_chain", &obj.CaChain)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateActionRevoke : The response body of the action to revoke the private certificate.
// This model "extends" SecretAction
type PrivateCertificateActionRevoke struct {
	// The type of secret action.
	ActionType *string `json:"action_type" validate:"required"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`
}

// Constants associated with the PrivateCertificateActionRevoke.ActionType property.
// The type of secret action.
const (
	PrivateCertificateActionRevoke_ActionType_PrivateCertActionRevokeCertificate   = "private_cert_action_revoke_certificate"
	PrivateCertificateActionRevoke_ActionType_PublicCertActionValidateDnsChallenge = "public_cert_action_validate_dns_challenge"
)

func (*PrivateCertificateActionRevoke) isaSecretAction() bool {
	return true
}

// UnmarshalPrivateCertificateActionRevoke unmarshals an instance of PrivateCertificateActionRevoke from the specified map of raw messages.
func UnmarshalPrivateCertificateActionRevoke(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateActionRevoke)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_seconds", &obj.RevocationTimeSeconds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateActionRevokePrototype : The request body to specify the properties of the action to revoke the private certificate.
// This model "extends" SecretActionPrototype
type PrivateCertificateActionRevokePrototype struct {
	// The type of secret action.
	ActionType *string `json:"action_type" validate:"required"`
}

// Constants associated with the PrivateCertificateActionRevokePrototype.ActionType property.
// The type of secret action.
const (
	PrivateCertificateActionRevokePrototype_ActionType_PrivateCertActionRevokeCertificate   = "private_cert_action_revoke_certificate"
	PrivateCertificateActionRevokePrototype_ActionType_PublicCertActionValidateDnsChallenge = "public_cert_action_validate_dns_challenge"
)

// NewPrivateCertificateActionRevokePrototype : Instantiate PrivateCertificateActionRevokePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateActionRevokePrototype(actionType string) (_model *PrivateCertificateActionRevokePrototype, err error) {
	_model = &PrivateCertificateActionRevokePrototype{
		ActionType: core.StringPtr(actionType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateActionRevokePrototype) isaSecretActionPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateActionRevokePrototype unmarshals an instance of PrivateCertificateActionRevokePrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateActionRevokePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateActionRevokePrototype)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionRevoke : The response body to specify the properties of the action to revoke the private certificate.
// This model "extends" ConfigurationAction
type PrivateCertificateConfigurationActionRevoke struct {
	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationActionRevoke.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionRevoke_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionRevoke_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionRevoke_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionRevoke_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionRevoke_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

func (*PrivateCertificateConfigurationActionRevoke) isaConfigurationAction() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionRevoke unmarshals an instance of PrivateCertificateConfigurationActionRevoke from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionRevoke(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionRevoke)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_seconds", &obj.RevocationTimeSeconds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionRevokePrototype : The request body to specify the properties of the action to revoke the private certificate configuration.
// This model "extends" ConfigurationActionPrototype
type PrivateCertificateConfigurationActionRevokePrototype struct {
	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationActionRevokePrototype.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionRevokePrototype_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionRevokePrototype_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionRevokePrototype_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionRevokePrototype_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionRevokePrototype_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

// NewPrivateCertificateConfigurationActionRevokePrototype : Instantiate PrivateCertificateConfigurationActionRevokePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateConfigurationActionRevokePrototype(actionType string) (_model *PrivateCertificateConfigurationActionRevokePrototype, err error) {
	_model = &PrivateCertificateConfigurationActionRevokePrototype{
		ActionType: core.StringPtr(actionType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateConfigurationActionRevokePrototype) isaConfigurationActionPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionRevokePrototype unmarshals an instance of PrivateCertificateConfigurationActionRevokePrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionRevokePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionRevokePrototype)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionRotateCRL : The response body of the action to rotate the CRL of an intermediate certificate authority for the private
// certificate configuration.
// This model "extends" ConfigurationAction
type PrivateCertificateConfigurationActionRotateCRL struct {
	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`

	// This field indicates whether the request to rotate the CRL for the private certificate configuration was successful.
	Success *bool `json:"success" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationActionRotateCRL.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionRotateCRL_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionRotateCRL_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionRotateCRL_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionRotateCRL_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionRotateCRL_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

func (*PrivateCertificateConfigurationActionRotateCRL) isaConfigurationAction() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionRotateCRL unmarshals an instance of PrivateCertificateConfigurationActionRotateCRL from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionRotateCRL(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionRotateCRL)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionRotateCRLPrototype : The request body of the action to rotate the CRL of an intermediate certificate authority for the private certificate
// configuration.
// This model "extends" ConfigurationActionPrototype
type PrivateCertificateConfigurationActionRotateCRLPrototype struct {
	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationActionRotateCRLPrototype.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionRotateCRLPrototype_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionRotateCRLPrototype_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionRotateCRLPrototype_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionRotateCRLPrototype_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionRotateCRLPrototype_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

// NewPrivateCertificateConfigurationActionRotateCRLPrototype : Instantiate PrivateCertificateConfigurationActionRotateCRLPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateConfigurationActionRotateCRLPrototype(actionType string) (_model *PrivateCertificateConfigurationActionRotateCRLPrototype, err error) {
	_model = &PrivateCertificateConfigurationActionRotateCRLPrototype{
		ActionType: core.StringPtr(actionType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateConfigurationActionRotateCRLPrototype) isaConfigurationActionPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionRotateCRLPrototype unmarshals an instance of PrivateCertificateConfigurationActionRotateCRLPrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionRotateCRLPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionRotateCRLPrototype)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionSetSigned : The response body of the action to set a signed intermediate certificate authority for the private certificate
// configuration.
// This model "extends" ConfigurationAction
type PrivateCertificateConfigurationActionSetSigned struct {
	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationActionSetSigned.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionSetSigned_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionSetSigned_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionSetSigned_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionSetSigned_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionSetSigned_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

func (*PrivateCertificateConfigurationActionSetSigned) isaConfigurationAction() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionSetSigned unmarshals an instance of PrivateCertificateConfigurationActionSetSigned from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionSetSigned(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionSetSigned)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionSetSignedPrototype : The request body of the action to set a signed intermediate certificate authority for the private certificate
// consideration.
// This model "extends" ConfigurationActionPrototype
type PrivateCertificateConfigurationActionSetSignedPrototype struct {
	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationActionSetSignedPrototype.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionSetSignedPrototype_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionSetSignedPrototype_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionSetSignedPrototype_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionSetSignedPrototype_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionSetSignedPrototype_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

// NewPrivateCertificateConfigurationActionSetSignedPrototype : Instantiate PrivateCertificateConfigurationActionSetSignedPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateConfigurationActionSetSignedPrototype(actionType string, certificate string) (_model *PrivateCertificateConfigurationActionSetSignedPrototype, err error) {
	_model = &PrivateCertificateConfigurationActionSetSignedPrototype{
		ActionType:  core.StringPtr(actionType),
		Certificate: core.StringPtr(certificate),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateConfigurationActionSetSignedPrototype) isaConfigurationActionPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionSetSignedPrototype unmarshals an instance of PrivateCertificateConfigurationActionSetSignedPrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionSetSignedPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionSetSignedPrototype)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionSignCSR : The response body of the action to sign the CSR for the private certificate configuration.
// This model "extends" ConfigurationAction
type PrivateCertificateConfigurationActionSignCSR struct {
	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '12h'. The value can't
	// exceed the `max_ttl` that is defined in the associated certificate template.
	TTL *string `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// This field indicates whether to use values from a certificate signing request (CSR) to complete a
	// `private_cert_configuration_action_sign_csr` action. If it is set to `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than by using the
	// values that are provided in the other parameters to this operation.
	//
	// 2) Any key usage, for example, non-repudiation, that is requested in the CSR are added to the basic set of key
	// usages used for CA certificates that are signed by the intermediate authority.
	//
	// 3) Extensions that are requested in the CSR are copied into the issued private certificate.
	UseCsrValues *bool `json:"use_csr_values,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`

	// The certificate signing request.
	Csr *string `json:"csr" validate:"required"`

	// The data that is associated with the root certificate authority.
	Data *PrivateCertificateConfigurationCACertificate `json:"data,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationActionSignCSR.Format property.
// The format of the returned data.
const (
	PrivateCertificateConfigurationActionSignCSR_Format_Pem       = "pem"
	PrivateCertificateConfigurationActionSignCSR_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificateConfigurationActionSignCSR.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionSignCSR_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionSignCSR_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionSignCSR_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionSignCSR_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionSignCSR_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

func (*PrivateCertificateConfigurationActionSignCSR) isaConfigurationAction() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionSignCSR unmarshals an instance of PrivateCertificateConfigurationActionSignCSR from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionSignCSR(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionSignCSR)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_path_length", &obj.MaxPathLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDnsDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_values", &obj.UseCsrValues)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "data", &obj.Data, UnmarshalPrivateCertificateConfigurationCACertificate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionSignCSRPrototype : The request body to specify the properties of the action to sign a CSR for the private certificate configuration.
// This model "extends" ConfigurationActionPrototype
type PrivateCertificateConfigurationActionSignCSRPrototype struct {
	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '12h'. The value can't
	// exceed the `max_ttl` that is defined in the associated certificate template.
	TTL *string `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// This field indicates whether to use values from a certificate signing request (CSR) to complete a
	// `private_cert_configuration_action_sign_csr` action. If it is set to `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than by using the
	// values that are provided in the other parameters to this operation.
	//
	// 2) Any key usage, for example, non-repudiation, that is requested in the CSR are added to the basic set of key
	// usages used for CA certificates that are signed by the intermediate authority.
	//
	// 3) Extensions that are requested in the CSR are copied into the issued private certificate.
	UseCsrValues *bool `json:"use_csr_values,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`

	// The certificate signing request.
	Csr *string `json:"csr" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationActionSignCSRPrototype.Format property.
// The format of the returned data.
const (
	PrivateCertificateConfigurationActionSignCSRPrototype_Format_Pem       = "pem"
	PrivateCertificateConfigurationActionSignCSRPrototype_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificateConfigurationActionSignCSRPrototype.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionSignCSRPrototype_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionSignCSRPrototype_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionSignCSRPrototype_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionSignCSRPrototype_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionSignCSRPrototype_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

// NewPrivateCertificateConfigurationActionSignCSRPrototype : Instantiate PrivateCertificateConfigurationActionSignCSRPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateConfigurationActionSignCSRPrototype(actionType string, csr string) (_model *PrivateCertificateConfigurationActionSignCSRPrototype, err error) {
	_model = &PrivateCertificateConfigurationActionSignCSRPrototype{
		ActionType: core.StringPtr(actionType),
		Csr:        core.StringPtr(csr),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateConfigurationActionSignCSRPrototype) isaConfigurationActionPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionSignCSRPrototype unmarshals an instance of PrivateCertificateConfigurationActionSignCSRPrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionSignCSRPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionSignCSRPrototype)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_path_length", &obj.MaxPathLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDnsDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_values", &obj.UseCsrValues)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionSignIntermediate : The response body of the action to sign the intermediate certificate authority for the private certificate
// configuration.
// This model "extends" ConfigurationAction
type PrivateCertificateConfigurationActionSignIntermediate struct {
	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '12h'. The value can't
	// exceed the `max_ttl` that is defined in the associated certificate template.
	TTL *string `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// This field indicates whether to use values from a certificate signing request (CSR) to complete a
	// `private_cert_configuration_action_sign_csr` action. If it is set to `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than by using the
	// values that are provided in the other parameters to this operation.
	//
	// 2) Any key usage, for example, non-repudiation, that is requested in the CSR are added to the basic set of key
	// usages used for CA certificates that are signed by the intermediate authority.
	//
	// 3) Extensions that are requested in the CSR are copied into the issued private certificate.
	UseCsrValues *bool `json:"use_csr_values,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`

	// The unique name of your configuration.
	IntermediateCertificateAuthority *string `json:"intermediate_certificate_authority" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationActionSignIntermediate.Format property.
// The format of the returned data.
const (
	PrivateCertificateConfigurationActionSignIntermediate_Format_Pem       = "pem"
	PrivateCertificateConfigurationActionSignIntermediate_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificateConfigurationActionSignIntermediate.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionSignIntermediate_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionSignIntermediate_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionSignIntermediate_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionSignIntermediate_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionSignIntermediate_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

func (*PrivateCertificateConfigurationActionSignIntermediate) isaConfigurationAction() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionSignIntermediate unmarshals an instance of PrivateCertificateConfigurationActionSignIntermediate from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionSignIntermediate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionSignIntermediate)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_path_length", &obj.MaxPathLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDnsDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_values", &obj.UseCsrValues)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_certificate_authority", &obj.IntermediateCertificateAuthority)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationActionSignIntermediatePrototype : The request body to specify the properties of the action to sign an intermediate certificate authority for the
// private certificate configuration.
// This model "extends" ConfigurationActionPrototype
type PrivateCertificateConfigurationActionSignIntermediatePrototype struct {
	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '12h'. The value can't
	// exceed the `max_ttl` that is defined in the associated certificate template.
	TTL *string `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// This field indicates whether to use values from a certificate signing request (CSR) to complete a
	// `private_cert_configuration_action_sign_csr` action. If it is set to `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than by using the
	// values that are provided in the other parameters to this operation.
	//
	// 2) Any key usage, for example, non-repudiation, that is requested in the CSR are added to the basic set of key
	// usages used for CA certificates that are signed by the intermediate authority.
	//
	// 3) Extensions that are requested in the CSR are copied into the issued private certificate.
	UseCsrValues *bool `json:"use_csr_values,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The type of configuration action.
	ActionType *string `json:"action_type" validate:"required"`

	// The unique name of your configuration.
	IntermediateCertificateAuthority *string `json:"intermediate_certificate_authority" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationActionSignIntermediatePrototype.Format property.
// The format of the returned data.
const (
	PrivateCertificateConfigurationActionSignIntermediatePrototype_Format_Pem       = "pem"
	PrivateCertificateConfigurationActionSignIntermediatePrototype_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificateConfigurationActionSignIntermediatePrototype.ActionType property.
// The type of configuration action.
const (
	PrivateCertificateConfigurationActionSignIntermediatePrototype_ActionType_PrivateCertConfigurationActionRevokeCaCertificate = "private_cert_configuration_action_revoke_ca_certificate"
	PrivateCertificateConfigurationActionSignIntermediatePrototype_ActionType_PrivateCertConfigurationActionRotateCrl           = "private_cert_configuration_action_rotate_crl"
	PrivateCertificateConfigurationActionSignIntermediatePrototype_ActionType_PrivateCertConfigurationActionSetSigned           = "private_cert_configuration_action_set_signed"
	PrivateCertificateConfigurationActionSignIntermediatePrototype_ActionType_PrivateCertConfigurationActionSignCsr             = "private_cert_configuration_action_sign_csr"
	PrivateCertificateConfigurationActionSignIntermediatePrototype_ActionType_PrivateCertConfigurationActionSignIntermediate    = "private_cert_configuration_action_sign_intermediate"
)

// NewPrivateCertificateConfigurationActionSignIntermediatePrototype : Instantiate PrivateCertificateConfigurationActionSignIntermediatePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateConfigurationActionSignIntermediatePrototype(actionType string, intermediateCertificateAuthority string) (_model *PrivateCertificateConfigurationActionSignIntermediatePrototype, err error) {
	_model = &PrivateCertificateConfigurationActionSignIntermediatePrototype{
		ActionType:                       core.StringPtr(actionType),
		IntermediateCertificateAuthority: core.StringPtr(intermediateCertificateAuthority),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateConfigurationActionSignIntermediatePrototype) isaConfigurationActionPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationActionSignIntermediatePrototype unmarshals an instance of PrivateCertificateConfigurationActionSignIntermediatePrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationActionSignIntermediatePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationActionSignIntermediatePrototype)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_path_length", &obj.MaxPathLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDnsDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_values", &obj.UseCsrValues)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_certificate_authority", &obj.IntermediateCertificateAuthority)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationCACertificate : The data that is associated with the root certificate authority.
// This model "extends" PrivateCertificateCAData
type PrivateCertificateConfigurationCACertificate struct {
	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`

	// The PEM-encoded certificate of the certificate authority that signed and issued this certificate.
	IssuingCa *string `json:"issuing_ca,omitempty"`

	// The chain of certificate authorities that are associated with the certificate.
	CaChain []string `json:"ca_chain,omitempty"`

	// The certificate expiration time.
	Expiration *int64 `json:"expiration,omitempty"`
}

func (*PrivateCertificateConfigurationCACertificate) isaPrivateCertificateCAData() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationCACertificate unmarshals an instance of PrivateCertificateConfigurationCACertificate from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationCACertificate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationCACertificate)
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_ca", &obj.IssuingCa)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ca_chain", &obj.CaChain)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration", &obj.Expiration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationIntermediateCA : The configuration of the root certificate authority.
// This model "extends" Configuration
type PrivateCertificateConfigurationIntermediateCA struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The signing method to use with this certificate authority to generate private certificates.
	//
	// You can choose between internal or externally signed options. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	SigningMethod *string `json:"signing_method" validate:"required"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA in seconds.
	MaxTtlSeconds *int64 `json:"max_ttl_seconds,omitempty"`

	// The time until the certificate revocation list (CRL) expires, in seconds.
	CrlExpirySeconds *int64 `json:"crl_expiry_seconds,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The configuration data of your Private Certificate.
	Data PrivateCertificateCADataIntf `json:"data,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationIntermediateCA.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationIntermediateCA_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationIntermediateCA_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationIntermediateCA_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationIntermediateCA_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationIntermediateCA_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationIntermediateCA_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationIntermediateCA_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCA.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateConfigurationIntermediateCA_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateConfigurationIntermediateCA_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateConfigurationIntermediateCA_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateConfigurationIntermediateCA_SecretType_Kv                 = "kv"
	PrivateCertificateConfigurationIntermediateCA_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateConfigurationIntermediateCA_SecretType_PublicCert         = "public_cert"
	PrivateCertificateConfigurationIntermediateCA_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateConfigurationIntermediateCA_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCA.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationIntermediateCA_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationIntermediateCA_KeyType_Rsa = "rsa"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCA.SigningMethod property.
// The signing method to use with this certificate authority to generate private certificates.
//
// You can choose between internal or externally signed options. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
const (
	PrivateCertificateConfigurationIntermediateCA_SigningMethod_External = "external"
	PrivateCertificateConfigurationIntermediateCA_SigningMethod_Internal = "internal"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCA.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	PrivateCertificateConfigurationIntermediateCA_Status_CertificateTemplateRequired = "certificate_template_required"
	PrivateCertificateConfigurationIntermediateCA_Status_Configured                  = "configured"
	PrivateCertificateConfigurationIntermediateCA_Status_Expired                     = "expired"
	PrivateCertificateConfigurationIntermediateCA_Status_Revoked                     = "revoked"
	PrivateCertificateConfigurationIntermediateCA_Status_SignedCertificateRequired   = "signed_certificate_required"
	PrivateCertificateConfigurationIntermediateCA_Status_SigningRequired             = "signing_required"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCA.Format property.
// The format of the returned data.
const (
	PrivateCertificateConfigurationIntermediateCA_Format_Pem       = "pem"
	PrivateCertificateConfigurationIntermediateCA_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCA.PrivateKeyFormat property.
// The format of the generated private key.
const (
	PrivateCertificateConfigurationIntermediateCA_PrivateKeyFormat_Der   = "der"
	PrivateCertificateConfigurationIntermediateCA_PrivateKeyFormat_Pkcs8 = "pkcs8"
)

func (*PrivateCertificateConfigurationIntermediateCA) isaConfiguration() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationIntermediateCA unmarshals an instance of PrivateCertificateConfigurationIntermediateCA from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationIntermediateCA(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationIntermediateCA)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_method", &obj.SigningMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_ttl_seconds", &obj.MaxTtlSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_expiry_seconds", &obj.CrlExpirySeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_disable", &obj.CrlDisable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_certificates_urls_encoded", &obj.IssuingCertificatesUrlsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_format", &obj.PrivateKeyFormat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "data", &obj.Data, UnmarshalPrivateCertificateCAData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationIntermediateCACSR : The data that is associated with the intermediate certificate authority.
// This model "extends" PrivateCertificateCAData
type PrivateCertificateConfigurationIntermediateCACSR struct {
	// The certificate signing request.
	Csr *string `json:"csr,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The type of private key to generate.
	PrivateKeyType *string `json:"private_key_type,omitempty"`

	// The certificate expiration time.
	Expiration *int64 `json:"expiration,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationIntermediateCACSR.PrivateKeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationIntermediateCACSR_PrivateKeyType_Ec  = "ec"
	PrivateCertificateConfigurationIntermediateCACSR_PrivateKeyType_Rsa = "rsa"
)

func (*PrivateCertificateConfigurationIntermediateCACSR) isaPrivateCertificateCAData() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationIntermediateCACSR unmarshals an instance of PrivateCertificateConfigurationIntermediateCACSR from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationIntermediateCACSR(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationIntermediateCACSR)
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_type", &obj.PrivateKeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration", &obj.Expiration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationIntermediateCAMetadata : The configuration of the metadata properties of the intermediate certificate authority.
// This model "extends" ConfigurationMetadata
type PrivateCertificateConfigurationIntermediateCAMetadata struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The signing method to use with this certificate authority to generate private certificates.
	//
	// You can choose between internal or externally signed options. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	SigningMethod *string `json:"signing_method" validate:"required"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationIntermediateCAMetadata.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationIntermediateCAMetadata_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationIntermediateCAMetadata_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationIntermediateCAMetadata_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationIntermediateCAMetadata_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationIntermediateCAMetadata_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationIntermediateCAMetadata_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationIntermediateCAMetadata_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCAMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateConfigurationIntermediateCAMetadata_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateConfigurationIntermediateCAMetadata_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateConfigurationIntermediateCAMetadata_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateConfigurationIntermediateCAMetadata_SecretType_Kv                 = "kv"
	PrivateCertificateConfigurationIntermediateCAMetadata_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateConfigurationIntermediateCAMetadata_SecretType_PublicCert         = "public_cert"
	PrivateCertificateConfigurationIntermediateCAMetadata_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateConfigurationIntermediateCAMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCAMetadata.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationIntermediateCAMetadata_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationIntermediateCAMetadata_KeyType_Rsa = "rsa"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCAMetadata.SigningMethod property.
// The signing method to use with this certificate authority to generate private certificates.
//
// You can choose between internal or externally signed options. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
const (
	PrivateCertificateConfigurationIntermediateCAMetadata_SigningMethod_External = "external"
	PrivateCertificateConfigurationIntermediateCAMetadata_SigningMethod_Internal = "internal"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCAMetadata.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	PrivateCertificateConfigurationIntermediateCAMetadata_Status_CertificateTemplateRequired = "certificate_template_required"
	PrivateCertificateConfigurationIntermediateCAMetadata_Status_Configured                  = "configured"
	PrivateCertificateConfigurationIntermediateCAMetadata_Status_Expired                     = "expired"
	PrivateCertificateConfigurationIntermediateCAMetadata_Status_Revoked                     = "revoked"
	PrivateCertificateConfigurationIntermediateCAMetadata_Status_SignedCertificateRequired   = "signed_certificate_required"
	PrivateCertificateConfigurationIntermediateCAMetadata_Status_SigningRequired             = "signing_required"
)

func (*PrivateCertificateConfigurationIntermediateCAMetadata) isaConfigurationMetadata() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationIntermediateCAMetadata unmarshals an instance of PrivateCertificateConfigurationIntermediateCAMetadata from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationIntermediateCAMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationIntermediateCAMetadata)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_method", &obj.SigningMethod)
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

// PrivateCertificateConfigurationIntermediateCAPatch : The configuration patch of the intermediate certificate authority.
// This model "extends" ConfigurationPatch
type PrivateCertificateConfigurationIntermediateCAPatch struct {
	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL *string `json:"max_ttl,omitempty"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry *string `json:"crl_expiry,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`
}

func (*PrivateCertificateConfigurationIntermediateCAPatch) isaConfigurationPatch() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationIntermediateCAPatch unmarshals an instance of PrivateCertificateConfigurationIntermediateCAPatch from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationIntermediateCAPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationIntermediateCAPatch)
	err = core.UnmarshalPrimitive(m, "max_ttl", &obj.MaxTTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_expiry", &obj.CrlExpiry)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_disable", &obj.CrlDisable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_certificates_urls_encoded", &obj.IssuingCertificatesUrlsEncoded)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the PrivateCertificateConfigurationIntermediateCAPatch
func (privateCertificateConfigurationIntermediateCAPatch *PrivateCertificateConfigurationIntermediateCAPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(privateCertificateConfigurationIntermediateCAPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// PrivateCertificateConfigurationIntermediateCAPrototype : The configuration of the intermediate certificate authority.
// This model "extends" ConfigurationPrototype
type PrivateCertificateConfigurationIntermediateCAPrototype struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// A human-readable unique name to assign to your configuration.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.
	Name *string `json:"name" validate:"required"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL *string `json:"max_ttl" validate:"required"`

	// The signing method to use with this certificate authority to generate private certificates.
	//
	// You can choose between internal or externally signed options. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	SigningMethod *string `json:"signing_method" validate:"required"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry *string `json:"crl_expiry,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationIntermediateCAPrototype.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationIntermediateCAPrototype_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationIntermediateCAPrototype_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationIntermediateCAPrototype_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationIntermediateCAPrototype_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationIntermediateCAPrototype_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationIntermediateCAPrototype_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationIntermediateCAPrototype_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCAPrototype.SigningMethod property.
// The signing method to use with this certificate authority to generate private certificates.
//
// You can choose between internal or externally signed options. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
const (
	PrivateCertificateConfigurationIntermediateCAPrototype_SigningMethod_External = "external"
	PrivateCertificateConfigurationIntermediateCAPrototype_SigningMethod_Internal = "internal"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCAPrototype.Format property.
// The format of the returned data.
const (
	PrivateCertificateConfigurationIntermediateCAPrototype_Format_Pem       = "pem"
	PrivateCertificateConfigurationIntermediateCAPrototype_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCAPrototype.PrivateKeyFormat property.
// The format of the generated private key.
const (
	PrivateCertificateConfigurationIntermediateCAPrototype_PrivateKeyFormat_Der   = "der"
	PrivateCertificateConfigurationIntermediateCAPrototype_PrivateKeyFormat_Pkcs8 = "pkcs8"
)

// Constants associated with the PrivateCertificateConfigurationIntermediateCAPrototype.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationIntermediateCAPrototype_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationIntermediateCAPrototype_KeyType_Rsa = "rsa"
)

// NewPrivateCertificateConfigurationIntermediateCAPrototype : Instantiate PrivateCertificateConfigurationIntermediateCAPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateConfigurationIntermediateCAPrototype(configType string, name string, maxTTL string, signingMethod string, commonName string) (_model *PrivateCertificateConfigurationIntermediateCAPrototype, err error) {
	_model = &PrivateCertificateConfigurationIntermediateCAPrototype{
		ConfigType:    core.StringPtr(configType),
		Name:          core.StringPtr(name),
		MaxTTL:        core.StringPtr(maxTTL),
		SigningMethod: core.StringPtr(signingMethod),
		CommonName:    core.StringPtr(commonName),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateConfigurationIntermediateCAPrototype) isaConfigurationPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationIntermediateCAPrototype unmarshals an instance of PrivateCertificateConfigurationIntermediateCAPrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationIntermediateCAPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationIntermediateCAPrototype)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_ttl", &obj.MaxTTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_method", &obj.SigningMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_expiry", &obj.CrlExpiry)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_disable", &obj.CrlDisable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_certificates_urls_encoded", &obj.IssuingCertificatesUrlsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_format", &obj.PrivateKeyFormat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationRootCA : The root certificate authority .
// This model "extends" Configuration
type PrivateCertificateConfigurationRootCA struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA in seconds.
	MaxTtlSeconds *int64 `json:"max_ttl_seconds,omitempty"`

	// The time until the certificate revocation list (CRL) expires, in seconds.
	CrlExpirySeconds *int64 `json:"crl_expiry_seconds,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// he requested TTL, after which the certificate expires.
	TtlSeconds *int64 `json:"ttl_seconds,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The configuration data of your Private Certificate.
	Data PrivateCertificateCADataIntf `json:"data,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationRootCA.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationRootCA_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationRootCA_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationRootCA_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationRootCA_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationRootCA_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationRootCA_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationRootCA_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationRootCA.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateConfigurationRootCA_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateConfigurationRootCA_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateConfigurationRootCA_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateConfigurationRootCA_SecretType_Kv                 = "kv"
	PrivateCertificateConfigurationRootCA_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateConfigurationRootCA_SecretType_PublicCert         = "public_cert"
	PrivateCertificateConfigurationRootCA_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateConfigurationRootCA_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificateConfigurationRootCA.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationRootCA_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationRootCA_KeyType_Rsa = "rsa"
)

// Constants associated with the PrivateCertificateConfigurationRootCA.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	PrivateCertificateConfigurationRootCA_Status_CertificateTemplateRequired = "certificate_template_required"
	PrivateCertificateConfigurationRootCA_Status_Configured                  = "configured"
	PrivateCertificateConfigurationRootCA_Status_Expired                     = "expired"
	PrivateCertificateConfigurationRootCA_Status_Revoked                     = "revoked"
	PrivateCertificateConfigurationRootCA_Status_SignedCertificateRequired   = "signed_certificate_required"
	PrivateCertificateConfigurationRootCA_Status_SigningRequired             = "signing_required"
)

// Constants associated with the PrivateCertificateConfigurationRootCA.Format property.
// The format of the returned data.
const (
	PrivateCertificateConfigurationRootCA_Format_Pem       = "pem"
	PrivateCertificateConfigurationRootCA_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificateConfigurationRootCA.PrivateKeyFormat property.
// The format of the generated private key.
const (
	PrivateCertificateConfigurationRootCA_PrivateKeyFormat_Der   = "der"
	PrivateCertificateConfigurationRootCA_PrivateKeyFormat_Pkcs8 = "pkcs8"
)

func (*PrivateCertificateConfigurationRootCA) isaConfiguration() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationRootCA unmarshals an instance of PrivateCertificateConfigurationRootCA from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationRootCA(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationRootCA)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_ttl_seconds", &obj.MaxTtlSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_expiry_seconds", &obj.CrlExpirySeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_disable", &obj.CrlDisable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_certificates_urls_encoded", &obj.IssuingCertificatesUrlsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl_seconds", &obj.TtlSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_format", &obj.PrivateKeyFormat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_path_length", &obj.MaxPathLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDnsDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "data", &obj.Data, UnmarshalPrivateCertificateCAData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationRootCAMetadata : The configuration of the metadata properties of the root certificate authority.
// This model "extends" ConfigurationMetadata
type PrivateCertificateConfigurationRootCAMetadata struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationRootCAMetadata.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationRootCAMetadata_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationRootCAMetadata_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationRootCAMetadata_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationRootCAMetadata_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationRootCAMetadata_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationRootCAMetadata_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationRootCAMetadata_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationRootCAMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateConfigurationRootCAMetadata_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateConfigurationRootCAMetadata_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateConfigurationRootCAMetadata_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateConfigurationRootCAMetadata_SecretType_Kv                 = "kv"
	PrivateCertificateConfigurationRootCAMetadata_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateConfigurationRootCAMetadata_SecretType_PublicCert         = "public_cert"
	PrivateCertificateConfigurationRootCAMetadata_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateConfigurationRootCAMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificateConfigurationRootCAMetadata.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationRootCAMetadata_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationRootCAMetadata_KeyType_Rsa = "rsa"
)

// Constants associated with the PrivateCertificateConfigurationRootCAMetadata.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	PrivateCertificateConfigurationRootCAMetadata_Status_CertificateTemplateRequired = "certificate_template_required"
	PrivateCertificateConfigurationRootCAMetadata_Status_Configured                  = "configured"
	PrivateCertificateConfigurationRootCAMetadata_Status_Expired                     = "expired"
	PrivateCertificateConfigurationRootCAMetadata_Status_Revoked                     = "revoked"
	PrivateCertificateConfigurationRootCAMetadata_Status_SignedCertificateRequired   = "signed_certificate_required"
	PrivateCertificateConfigurationRootCAMetadata_Status_SigningRequired             = "signing_required"
)

func (*PrivateCertificateConfigurationRootCAMetadata) isaConfigurationMetadata() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationRootCAMetadata unmarshals an instance of PrivateCertificateConfigurationRootCAMetadata from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationRootCAMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationRootCAMetadata)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
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

// PrivateCertificateConfigurationRootCAPatch : The configuration of the metadata patch for the root certificate authority.
// This model "extends" ConfigurationPatch
type PrivateCertificateConfigurationRootCAPatch struct {
	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL *string `json:"max_ttl,omitempty"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry *string `json:"crl_expiry,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`
}

func (*PrivateCertificateConfigurationRootCAPatch) isaConfigurationPatch() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationRootCAPatch unmarshals an instance of PrivateCertificateConfigurationRootCAPatch from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationRootCAPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationRootCAPatch)
	err = core.UnmarshalPrimitive(m, "max_ttl", &obj.MaxTTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_expiry", &obj.CrlExpiry)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_disable", &obj.CrlDisable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_certificates_urls_encoded", &obj.IssuingCertificatesUrlsEncoded)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the PrivateCertificateConfigurationRootCAPatch
func (privateCertificateConfigurationRootCAPatch *PrivateCertificateConfigurationRootCAPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(privateCertificateConfigurationRootCAPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// PrivateCertificateConfigurationRootCAPrototype : The configuration of the root certificate authority.
// This model "extends" ConfigurationPrototype
type PrivateCertificateConfigurationRootCAPrototype struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// A human-readable unique name to assign to your configuration.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.
	Name *string `json:"name" validate:"required"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL *string `json:"max_ttl" validate:"required"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry *string `json:"crl_expiry,omitempty"`

	// This field disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when you're downloading the CRL. If CRL
	// building is enabled, it rebuilds the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// This field determines whether to encode the certificate revocation list (CRL) distribution points in the
	// certificates that are issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// This field determines whether to encode the URL of the issuing certificate in the certificates that are issued by
	// this certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The requested time-to-live (TTL) for certificates that are created by this CA. This field's value can't be longer
	// than the `max_ttl` limit.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	TTL *string `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.
	PermittedDnsDomains []string `json:"permitted_dns_domains,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The requested value for the [`serialNumber`](https://datatracker.ietf.org/doc/html/rfc4519#section-2.31) attribute
	// that is in the certificate's distinguished name (DN).
	//
	// **Note:** This field is not related to the `serial_number` field that is returned in the API response. The
	// `serial_number` field represents the certificate's randomly assigned serial number.
	SerialNumber *string `json:"serial_number,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationRootCAPrototype.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationRootCAPrototype_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationRootCAPrototype_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationRootCAPrototype_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationRootCAPrototype_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationRootCAPrototype_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationRootCAPrototype_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationRootCAPrototype_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationRootCAPrototype.Format property.
// The format of the returned data.
const (
	PrivateCertificateConfigurationRootCAPrototype_Format_Pem       = "pem"
	PrivateCertificateConfigurationRootCAPrototype_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificateConfigurationRootCAPrototype.PrivateKeyFormat property.
// The format of the generated private key.
const (
	PrivateCertificateConfigurationRootCAPrototype_PrivateKeyFormat_Der   = "der"
	PrivateCertificateConfigurationRootCAPrototype_PrivateKeyFormat_Pkcs8 = "pkcs8"
)

// Constants associated with the PrivateCertificateConfigurationRootCAPrototype.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationRootCAPrototype_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationRootCAPrototype_KeyType_Rsa = "rsa"
)

// NewPrivateCertificateConfigurationRootCAPrototype : Instantiate PrivateCertificateConfigurationRootCAPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateConfigurationRootCAPrototype(configType string, name string, maxTTL string, commonName string) (_model *PrivateCertificateConfigurationRootCAPrototype, err error) {
	_model = &PrivateCertificateConfigurationRootCAPrototype{
		ConfigType: core.StringPtr(configType),
		Name:       core.StringPtr(name),
		MaxTTL:     core.StringPtr(maxTTL),
		CommonName: core.StringPtr(commonName),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateConfigurationRootCAPrototype) isaConfigurationPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationRootCAPrototype unmarshals an instance of PrivateCertificateConfigurationRootCAPrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationRootCAPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationRootCAPrototype)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_ttl", &obj.MaxTTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_expiry", &obj.CrlExpiry)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_disable", &obj.CrlDisable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crl_distribution_points_encoded", &obj.CrlDistributionPointsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_certificates_urls_encoded", &obj.IssuingCertificatesUrlsEncoded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_format", &obj.PrivateKeyFormat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_path_length", &obj.MaxPathLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDnsDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationTemplate : The configuration of the private certificate template.
// This model "extends" Configuration
type PrivateCertificateConfigurationTemplate struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The name of the intermediate certificate authority.
	CertificateAuthority *string `json:"certificate_authority" validate:"required"`

	// This field scopes the creation of private certificates to only the secret groups that you specify.
	//
	// This field can be supplied as a comma-delimited list of secret group IDs.
	AllowedSecretGroups *string `json:"allowed_secret_groups,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA in seconds.
	MaxTtlSeconds *int64 `json:"max_ttl_seconds,omitempty"`

	// he requested TTL, after which the certificate expires.
	TtlSeconds *int64 `json:"ttl_seconds,omitempty"`

	// This field indicates whether to allow `localhost` to be included as one of the requested common names.
	AllowLocalhost *bool `json:"allow_localhost,omitempty"`

	// The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and
	// `allow_subdomains` options.
	AllowedDomains []string `json:"allowed_domains,omitempty"`

	// This field indicates whether to allow the domains that are supplied in the `allowed_domains` field to contain access
	// control list (ACL) templates.
	AllowedDomainsTemplate *bool `json:"allowed_domains_template,omitempty"`

	// This field indicates whether to allow clients to request private certificates that match the value of the actual
	// domains on the final certificate.
	//
	// For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a
	// certificate that contains the name `example.com` as one of the DNS values on the final certificate.
	//
	// **Important:** In some scenarios, allowing bare domains can be considered a security risk.
	AllowBareDomains *bool `json:"allow_bare_domains,omitempty"`

	// This field indicates whether to allow clients to request private certificates with common names (CN) that are
	// subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.
	//
	// For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the
	// following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.
	//
	// **Note:** This field is redundant if you use the `allow_any_name` option.
	AllowSubdomains *bool `json:"allow_subdomains,omitempty"`

	// This field indicates whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are
	// specified in the `allowed_domains` field.
	//
	// If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	AllowGlobDomains *bool `json:"allow_glob_domains,omitempty"`

	// This field indicates whether to allow clients to request a private certificate that matches any common name.
	AllowAnyName *bool `json:"allow_any_name,omitempty"`

	// This field indicates whether to enforce only valid hostnames for common names, DNS Subject Alternative Names, and
	// the host section of email addresses.
	EnforceHostnames *bool `json:"enforce_hostnames,omitempty"`

	// This field indicates whether to allow clients to request a private certificate with IP Subject Alternative Names.
	AllowIpSans *bool `json:"allow_ip_sans,omitempty"`

	// The URI Subject Alternative Names to allow for private certificates.
	//
	// Values can contain glob patterns, for example `spiffe://hostname/_*`.
	AllowedUriSans []string `json:"allowed_uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private
	// certificates.
	//
	// The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type
	// is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any
	// `other_sans` input.
	AllowedOtherSans []string `json:"allowed_other_sans,omitempty"`

	// This field indicates whether private certificates are flagged for server use.
	ServerFlag *bool `json:"server_flag,omitempty"`

	// This field indicates whether private certificates are flagged for client use.
	ClientFlag *bool `json:"client_flag,omitempty"`

	// This field indicates whether private certificates are flagged for code signing use.
	CodeSigningFlag *bool `json:"code_signing_flag,omitempty"`

	// This field indicates whether private certificates are flagged for email protection use.
	EmailProtectionFlag *bool `json:"email_protection_flag,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The allowed key usage constraint to define for private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage). Omit the
	// `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to
	// an empty list.
	KeyUsage []string `json:"key_usage,omitempty"`

	// The allowed extended key usage constraint on private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage).
	// Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set
	// this field to an empty list.
	ExtKeyUsage []string `json:"ext_key_usage,omitempty"`

	// A list of extended key usage Object Identifiers (OIDs).
	ExtKeyUsageOids []string `json:"ext_key_usage_oids,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// common name (CN) from a certificate signing request (CSR) instead of the CN that is included in the data of the
	// certificate.
	//
	// Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include
	// the `use_csr_sans` property.
	UseCsrCommonName *bool `json:"use_csr_common_name,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// Subject Alternative Names
	// (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the data of the
	// certificate.
	//
	// This field does not include the common name in the CSR. To use the common name, include the `use_csr_common_name`
	// property.
	UseCsrSans *bool `json:"use_csr_sans,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// This field is deprecated. You can ignore its value.
	SerialNumber *string `json:"serial_number,omitempty"`

	// This field indicates whether to require a common name to create a private certificate.
	//
	// By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the
	// `require_cn` option to `false`.
	RequireCn *bool `json:"require_cn,omitempty"`

	// A list of policy Object Identifiers (OIDs).
	PolicyIdentifiers []string `json:"policy_identifiers,omitempty"`

	// This field indicates whether to mark the Basic Constraints extension of an issued private certificate as valid for
	// non-CA certificates.
	BasicConstraintsValidForNonCa *bool `json:"basic_constraints_valid_for_non_ca,omitempty"`

	// The duration in seconds by which to backdate the `not_before` property of an issued private certificate.
	NotBeforeDurationSeconds *int64 `json:"not_before_duration_seconds,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationTemplate.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationTemplate_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationTemplate_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationTemplate_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationTemplate_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationTemplate_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationTemplate_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationTemplate_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationTemplate.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateConfigurationTemplate_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateConfigurationTemplate_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateConfigurationTemplate_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateConfigurationTemplate_SecretType_Kv                 = "kv"
	PrivateCertificateConfigurationTemplate_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateConfigurationTemplate_SecretType_PublicCert         = "public_cert"
	PrivateCertificateConfigurationTemplate_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateConfigurationTemplate_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificateConfigurationTemplate.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationTemplate_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationTemplate_KeyType_Rsa = "rsa"
)

func (*PrivateCertificateConfigurationTemplate) isaConfiguration() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationTemplate unmarshals an instance of PrivateCertificateConfigurationTemplate from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationTemplate)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_secret_groups", &obj.AllowedSecretGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_ttl_seconds", &obj.MaxTtlSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl_seconds", &obj.TtlSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_localhost", &obj.AllowLocalhost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_domains", &obj.AllowedDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_domains_template", &obj.AllowedDomainsTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_bare_domains", &obj.AllowBareDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_subdomains", &obj.AllowSubdomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_glob_domains", &obj.AllowGlobDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_any_name", &obj.AllowAnyName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enforce_hostnames", &obj.EnforceHostnames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_ip_sans", &obj.AllowIpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_uri_sans", &obj.AllowedUriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_other_sans", &obj.AllowedOtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_flag", &obj.ServerFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_flag", &obj.ClientFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "code_signing_flag", &obj.CodeSigningFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email_protection_flag", &obj.EmailProtectionFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_usage", &obj.KeyUsage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ext_key_usage", &obj.ExtKeyUsage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ext_key_usage_oids", &obj.ExtKeyUsageOids)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_common_name", &obj.UseCsrCommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_sans", &obj.UseCsrSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "require_cn", &obj.RequireCn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_identifiers", &obj.PolicyIdentifiers)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "basic_constraints_valid_for_non_ca", &obj.BasicConstraintsValidForNonCa)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_before_duration_seconds", &obj.NotBeforeDurationSeconds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationTemplateMetadata : The metadata properties of the configuration of the private certificate template.
// This model "extends" ConfigurationMetadata
type PrivateCertificateConfigurationTemplateMetadata struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The name of the intermediate certificate authority.
	CertificateAuthority *string `json:"certificate_authority" validate:"required"`
}

// Constants associated with the PrivateCertificateConfigurationTemplateMetadata.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationTemplateMetadata_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationTemplateMetadata_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationTemplateMetadata_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationTemplateMetadata_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationTemplateMetadata_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationTemplateMetadata_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationTemplateMetadata_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationTemplateMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateConfigurationTemplateMetadata_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateConfigurationTemplateMetadata_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateConfigurationTemplateMetadata_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateConfigurationTemplateMetadata_SecretType_Kv                 = "kv"
	PrivateCertificateConfigurationTemplateMetadata_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateConfigurationTemplateMetadata_SecretType_PublicCert         = "public_cert"
	PrivateCertificateConfigurationTemplateMetadata_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateConfigurationTemplateMetadata_SecretType_UsernamePassword   = "username_password"
)

func (*PrivateCertificateConfigurationTemplateMetadata) isaConfigurationMetadata() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationTemplateMetadata unmarshals an instance of PrivateCertificateConfigurationTemplateMetadata from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationTemplateMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationTemplateMetadata)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateConfigurationTemplatePatch : Properties that describe a certificate template. You can use a certificate template to control the parameters that
// are applied to your issued private certificates. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-certificate-templates).
// This model "extends" ConfigurationPatch
type PrivateCertificateConfigurationTemplatePatch struct {
	// This field scopes the creation of private certificates to only the secret groups that you specify.
	//
	// This field can be supplied as a comma-delimited list of secret group IDs.
	AllowedSecretGroups *string `json:"allowed_secret_groups,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL *string `json:"max_ttl,omitempty"`

	// The requested time-to-live (TTL) for certificates that are created by this CA. This field's value can't be longer
	// than the `max_ttl` limit.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	TTL *string `json:"ttl,omitempty"`

	// This field indicates whether to allow `localhost` to be included as one of the requested common names.
	AllowLocalhost *bool `json:"allow_localhost,omitempty"`

	// The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and
	// `allow_subdomains` options.
	AllowedDomains []string `json:"allowed_domains,omitempty"`

	// This field indicates whether to allow the domains that are supplied in the `allowed_domains` field to contain access
	// control list (ACL) templates.
	AllowedDomainsTemplate *bool `json:"allowed_domains_template,omitempty"`

	// This field indicates whether to allow clients to request private certificates that match the value of the actual
	// domains on the final certificate.
	//
	// For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a
	// certificate that contains the name `example.com` as one of the DNS values on the final certificate.
	//
	// **Important:** In some scenarios, allowing bare domains can be considered a security risk.
	AllowBareDomains *bool `json:"allow_bare_domains,omitempty"`

	// This field indicates whether to allow clients to request private certificates with common names (CN) that are
	// subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.
	//
	// For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the
	// following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.
	//
	// **Note:** This field is redundant if you use the `allow_any_name` option.
	AllowSubdomains *bool `json:"allow_subdomains,omitempty"`

	// This field indicates whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are
	// specified in the `allowed_domains` field.
	//
	// If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	AllowGlobDomains *bool `json:"allow_glob_domains,omitempty"`

	// This field indicates whether to allow clients to request a private certificate that matches any common name.
	AllowAnyName *bool `json:"allow_any_name,omitempty"`

	// This field indicates whether to enforce only valid hostnames for common names, DNS Subject Alternative Names, and
	// the host section of email addresses.
	EnforceHostnames *bool `json:"enforce_hostnames,omitempty"`

	// This field indicates whether to allow clients to request a private certificate with IP Subject Alternative Names.
	AllowIpSans *bool `json:"allow_ip_sans,omitempty"`

	// The URI Subject Alternative Names to allow for private certificates.
	//
	// Values can contain glob patterns, for example `spiffe://hostname/_*`.
	AllowedUriSans []string `json:"allowed_uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private
	// certificates.
	//
	// The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type
	// is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any
	// `other_sans` input.
	AllowedOtherSans []string `json:"allowed_other_sans,omitempty"`

	// This field indicates whether private certificates are flagged for server use.
	ServerFlag *bool `json:"server_flag,omitempty"`

	// This field indicates whether private certificates are flagged for client use.
	ClientFlag *bool `json:"client_flag,omitempty"`

	// This field indicates whether private certificates are flagged for code signing use.
	CodeSigningFlag *bool `json:"code_signing_flag,omitempty"`

	// This field indicates whether private certificates are flagged for email protection use.
	EmailProtectionFlag *bool `json:"email_protection_flag,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The allowed key usage constraint to define for private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage). Omit the
	// `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to
	// an empty list.
	KeyUsage []string `json:"key_usage,omitempty"`

	// The allowed extended key usage constraint on private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage).
	// Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set
	// this field to an empty list.
	ExtKeyUsage []string `json:"ext_key_usage,omitempty"`

	// A list of extended key usage Object Identifiers (OIDs).
	ExtKeyUsageOids []string `json:"ext_key_usage_oids,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// common name (CN) from a certificate signing request (CSR) instead of the CN that is included in the data of the
	// certificate.
	//
	// Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include
	// the `use_csr_sans` property.
	UseCsrCommonName *bool `json:"use_csr_common_name,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// Subject Alternative Names
	// (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the data of the
	// certificate.
	//
	// This field does not include the common name in the CSR. To use the common name, include the `use_csr_common_name`
	// property.
	UseCsrSans *bool `json:"use_csr_sans,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// This field is deprecated. You can ignore its value.
	SerialNumber *string `json:"serial_number,omitempty"`

	// This field indicates whether to require a common name to create a private certificate.
	//
	// By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the
	// `require_cn` option to `false`.
	RequireCn *bool `json:"require_cn,omitempty"`

	// A list of policy Object Identifiers (OIDs).
	PolicyIdentifiers []string `json:"policy_identifiers,omitempty"`

	// This field indicates whether to mark the Basic Constraints extension of an issued private certificate as valid for
	// non-CA certificates.
	BasicConstraintsValidForNonCa *bool `json:"basic_constraints_valid_for_non_ca,omitempty"`

	// The duration in seconds by which to backdate the `not_before` property of an issued private certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value
	// is returned in seconds (integer).
	NotBeforeDuration *string `json:"not_before_duration,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationTemplatePatch.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationTemplatePatch_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationTemplatePatch_KeyType_Rsa = "rsa"
)

func (*PrivateCertificateConfigurationTemplatePatch) isaConfigurationPatch() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationTemplatePatch unmarshals an instance of PrivateCertificateConfigurationTemplatePatch from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationTemplatePatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationTemplatePatch)
	err = core.UnmarshalPrimitive(m, "allowed_secret_groups", &obj.AllowedSecretGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_ttl", &obj.MaxTTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_localhost", &obj.AllowLocalhost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_domains", &obj.AllowedDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_domains_template", &obj.AllowedDomainsTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_bare_domains", &obj.AllowBareDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_subdomains", &obj.AllowSubdomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_glob_domains", &obj.AllowGlobDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_any_name", &obj.AllowAnyName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enforce_hostnames", &obj.EnforceHostnames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_ip_sans", &obj.AllowIpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_uri_sans", &obj.AllowedUriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_other_sans", &obj.AllowedOtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_flag", &obj.ServerFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_flag", &obj.ClientFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "code_signing_flag", &obj.CodeSigningFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email_protection_flag", &obj.EmailProtectionFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_usage", &obj.KeyUsage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ext_key_usage", &obj.ExtKeyUsage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ext_key_usage_oids", &obj.ExtKeyUsageOids)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_common_name", &obj.UseCsrCommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_sans", &obj.UseCsrSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "require_cn", &obj.RequireCn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_identifiers", &obj.PolicyIdentifiers)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "basic_constraints_valid_for_non_ca", &obj.BasicConstraintsValidForNonCa)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_before_duration", &obj.NotBeforeDuration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the PrivateCertificateConfigurationTemplatePatch
func (privateCertificateConfigurationTemplatePatch *PrivateCertificateConfigurationTemplatePatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(privateCertificateConfigurationTemplatePatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// PrivateCertificateConfigurationTemplatePrototype : Properties that describe a certificate template. You can use a certificate template to control the parameters that
// are applied to your issued private certificates. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-certificate-templates).
// This model "extends" ConfigurationPrototype
type PrivateCertificateConfigurationTemplatePrototype struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// A human-readable unique name to assign to your configuration.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.
	Name *string `json:"name" validate:"required"`

	// The name of the intermediate certificate authority.
	CertificateAuthority *string `json:"certificate_authority" validate:"required"`

	// This field scopes the creation of private certificates to only the secret groups that you specify.
	//
	// This field can be supplied as a comma-delimited list of secret group IDs.
	AllowedSecretGroups *string `json:"allowed_secret_groups,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL *string `json:"max_ttl,omitempty"`

	// The requested time-to-live (TTL) for certificates that are created by this CA. This field's value can't be longer
	// than the `max_ttl` limit.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	TTL *string `json:"ttl,omitempty"`

	// This field indicates whether to allow `localhost` to be included as one of the requested common names.
	AllowLocalhost *bool `json:"allow_localhost,omitempty"`

	// The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and
	// `allow_subdomains` options.
	AllowedDomains []string `json:"allowed_domains,omitempty"`

	// This field indicates whether to allow the domains that are supplied in the `allowed_domains` field to contain access
	// control list (ACL) templates.
	AllowedDomainsTemplate *bool `json:"allowed_domains_template,omitempty"`

	// This field indicates whether to allow clients to request private certificates that match the value of the actual
	// domains on the final certificate.
	//
	// For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a
	// certificate that contains the name `example.com` as one of the DNS values on the final certificate.
	//
	// **Important:** In some scenarios, allowing bare domains can be considered a security risk.
	AllowBareDomains *bool `json:"allow_bare_domains,omitempty"`

	// This field indicates whether to allow clients to request private certificates with common names (CN) that are
	// subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.
	//
	// For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the
	// following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.
	//
	// **Note:** This field is redundant if you use the `allow_any_name` option.
	AllowSubdomains *bool `json:"allow_subdomains,omitempty"`

	// This field indicates whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are
	// specified in the `allowed_domains` field.
	//
	// If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	AllowGlobDomains *bool `json:"allow_glob_domains,omitempty"`

	// This field indicates whether the issuance of certificates with RFC 6125 wildcards in the CN field.
	//
	// When set to false, this field prevents wildcards from being issued even if they can be allowed by an option
	// `allow_glob_domains`.
	AllowWildcardCertificates *bool `json:"allow_wildcard_certificates,omitempty"`

	// This field indicates whether to allow clients to request a private certificate that matches any common name.
	AllowAnyName *bool `json:"allow_any_name,omitempty"`

	// This field indicates whether to enforce only valid hostnames for common names, DNS Subject Alternative Names, and
	// the host section of email addresses.
	EnforceHostnames *bool `json:"enforce_hostnames,omitempty"`

	// This field indicates whether to allow clients to request a private certificate with IP Subject Alternative Names.
	AllowIpSans *bool `json:"allow_ip_sans,omitempty"`

	// The URI Subject Alternative Names to allow for private certificates.
	//
	// Values can contain glob patterns, for example `spiffe://hostname/_*`.
	AllowedUriSans []string `json:"allowed_uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private
	// certificates.
	//
	// The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type
	// is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any
	// `other_sans` input.
	AllowedOtherSans []string `json:"allowed_other_sans,omitempty"`

	// This field indicates whether private certificates are flagged for server use.
	ServerFlag *bool `json:"server_flag,omitempty"`

	// This field indicates whether private certificates are flagged for client use.
	ClientFlag *bool `json:"client_flag,omitempty"`

	// This field indicates whether private certificates are flagged for code signing use.
	CodeSigningFlag *bool `json:"code_signing_flag,omitempty"`

	// This field indicates whether private certificates are flagged for email protection use.
	EmailProtectionFlag *bool `json:"email_protection_flag,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use to generate the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The allowed key usage constraint to define for private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage). Omit the
	// `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to
	// an empty list.
	KeyUsage []string `json:"key_usage,omitempty"`

	// The allowed extended key usage constraint on private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage).
	// Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set
	// this field to an empty list.
	ExtKeyUsage []string `json:"ext_key_usage,omitempty"`

	// A list of extended key usage Object Identifiers (OIDs).
	ExtKeyUsageOids []string `json:"ext_key_usage_oids,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// common name (CN) from a certificate signing request (CSR) instead of the CN that is included in the data of the
	// certificate.
	//
	// Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include
	// the `use_csr_sans` property.
	UseCsrCommonName *bool `json:"use_csr_common_name,omitempty"`

	// When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the
	// Subject Alternative Names
	// (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the data of the
	// certificate.
	//
	// This field does not include the common name in the CSR. To use the common name, include the `use_csr_common_name`
	// property.
	UseCsrSans *bool `json:"use_csr_sans,omitempty"`

	// The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	Ou []string `json:"ou,omitempty"`

	// The Organization (O) values to define in the subject field of the resulting certificate.
	Organization []string `json:"organization,omitempty"`

	// The Country (C) values to define in the subject field of the resulting certificate.
	Country []string `json:"country,omitempty"`

	// The Locality (L) values to define in the subject field of the resulting certificate.
	Locality []string `json:"locality,omitempty"`

	// The Province (ST) values to define in the subject field of the resulting certificate.
	Province []string `json:"province,omitempty"`

	// The street address values to define in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The postal code values to define in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// This field is deprecated. You can ignore its value.
	SerialNumber *string `json:"serial_number,omitempty"`

	// This field indicates whether to require a common name to create a private certificate.
	//
	// By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the
	// `require_cn` option to `false`.
	RequireCn *bool `json:"require_cn,omitempty"`

	// A list of policy Object Identifiers (OIDs).
	PolicyIdentifiers []string `json:"policy_identifiers,omitempty"`

	// This field indicates whether to mark the Basic Constraints extension of an issued private certificate as valid for
	// non-CA certificates.
	BasicConstraintsValidForNonCa *bool `json:"basic_constraints_valid_for_non_ca,omitempty"`

	// The duration in seconds by which to backdate the `not_before` property of an issued private certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value
	// is returned in seconds (integer).
	NotBeforeDuration *string `json:"not_before_duration,omitempty"`
}

// Constants associated with the PrivateCertificateConfigurationTemplatePrototype.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PrivateCertificateConfigurationTemplatePrototype_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PrivateCertificateConfigurationTemplatePrototype_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PrivateCertificateConfigurationTemplatePrototype_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PrivateCertificateConfigurationTemplatePrototype_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PrivateCertificateConfigurationTemplatePrototype_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PrivateCertificateConfigurationTemplatePrototype_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PrivateCertificateConfigurationTemplatePrototype_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PrivateCertificateConfigurationTemplatePrototype.KeyType property.
// The type of private key to generate.
const (
	PrivateCertificateConfigurationTemplatePrototype_KeyType_Ec  = "ec"
	PrivateCertificateConfigurationTemplatePrototype_KeyType_Rsa = "rsa"
)

// NewPrivateCertificateConfigurationTemplatePrototype : Instantiate PrivateCertificateConfigurationTemplatePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateConfigurationTemplatePrototype(configType string, name string, certificateAuthority string) (_model *PrivateCertificateConfigurationTemplatePrototype, err error) {
	_model = &PrivateCertificateConfigurationTemplatePrototype{
		ConfigType:           core.StringPtr(configType),
		Name:                 core.StringPtr(name),
		CertificateAuthority: core.StringPtr(certificateAuthority),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateConfigurationTemplatePrototype) isaConfigurationPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateConfigurationTemplatePrototype unmarshals an instance of PrivateCertificateConfigurationTemplatePrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateConfigurationTemplatePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateConfigurationTemplatePrototype)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_secret_groups", &obj.AllowedSecretGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_ttl", &obj.MaxTTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_localhost", &obj.AllowLocalhost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_domains", &obj.AllowedDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_domains_template", &obj.AllowedDomainsTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_bare_domains", &obj.AllowBareDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_subdomains", &obj.AllowSubdomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_glob_domains", &obj.AllowGlobDomains)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_wildcard_certificates", &obj.AllowWildcardCertificates)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_any_name", &obj.AllowAnyName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enforce_hostnames", &obj.EnforceHostnames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_ip_sans", &obj.AllowIpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_uri_sans", &obj.AllowedUriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_other_sans", &obj.AllowedOtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_flag", &obj.ServerFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_flag", &obj.ClientFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "code_signing_flag", &obj.CodeSigningFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email_protection_flag", &obj.EmailProtectionFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_type", &obj.KeyType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_bits", &obj.KeyBits)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_usage", &obj.KeyUsage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ext_key_usage", &obj.ExtKeyUsage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ext_key_usage_oids", &obj.ExtKeyUsageOids)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_common_name", &obj.UseCsrCommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_csr_sans", &obj.UseCsrSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ou", &obj.Ou)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "organization", &obj.Organization)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locality", &obj.Locality)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "province", &obj.Province)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "street_address", &obj.StreetAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "postal_code", &obj.PostalCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "require_cn", &obj.RequireCn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_identifiers", &obj.PolicyIdentifiers)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "basic_constraints_valid_for_non_ca", &obj.BasicConstraintsValidForNonCa)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_before_duration", &obj.NotBeforeDuration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateMetadata : Properties of the metadata of your private certificate.
// This model "extends" SecretMetadata
type PrivateCertificateMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The identifier for the cryptographic algorithm that is used by the issuing certificate authority to sign a
	// certificate.
	SigningAlgorithm *string `json:"signing_algorithm" validate:"required"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The intermediate certificate authority that signed this certificate.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template" validate:"required"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date" validate:"required"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer" validate:"required"`

	// The identifier for the cryptographic algorithm used to generate the public key that is associated with the
	// certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number" validate:"required"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity" validate:"required"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`

	// The date and time that the certificate was revoked. The date format follows `RFC 3339`.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`
}

// Constants associated with the PrivateCertificateMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateMetadata_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateMetadata_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateMetadata_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateMetadata_SecretType_Kv                 = "kv"
	PrivateCertificateMetadata_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateMetadata_SecretType_PublicCert         = "public_cert"
	PrivateCertificateMetadata_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificateMetadata.StateDescription property.
// A text representation of the secret state.
const (
	PrivateCertificateMetadata_StateDescription_Active        = "active"
	PrivateCertificateMetadata_StateDescription_Deactivated   = "deactivated"
	PrivateCertificateMetadata_StateDescription_Destroyed     = "destroyed"
	PrivateCertificateMetadata_StateDescription_PreActivation = "pre_activation"
	PrivateCertificateMetadata_StateDescription_Suspended     = "suspended"
)

func (*PrivateCertificateMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalPrivateCertificateMetadata unmarshals an instance of PrivateCertificateMetadata from the specified map of raw messages.
func UnmarshalPrivateCertificateMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateMetadata)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_algorithm", &obj.SigningAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_template", &obj.CertificateTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_seconds", &obj.RevocationTimeSeconds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_rfc3339", &obj.RevocationTimeRfc3339)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateMetadataPatch : PrivateCertificateMetadataPatch struct
// This model "extends" SecretMetadataPatch
type PrivateCertificateMetadataPatch struct {
	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`
}

func (*PrivateCertificateMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

// UnmarshalPrivateCertificateMetadataPatch unmarshals an instance of PrivateCertificateMetadataPatch from the specified map of raw messages.
func UnmarshalPrivateCertificateMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateMetadataPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the PrivateCertificateMetadataPatch
func (privateCertificateMetadataPatch *PrivateCertificateMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(privateCertificateMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// PrivateCertificatePrototype : PrivateCertificatePrototype struct
// This model "extends" SecretPrototype
type PrivateCertificatePrototype struct {
	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template" validate:"required"`

	// The Common Name (CN) represents the server name that is protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IpSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	UriSans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The certificate signing request.
	Csr *string `json:"csr,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If the common name is set to `true`, it is not included in DNS, or email SANs if they apply. This field can be
	// useful if the common name is a human-readable identifier, instead of a hostname or an email address.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '12h'. The value can't
	// exceed the `max_ttl` that is defined in the associated certificate template.
	TTL *string `json:"ttl,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// Constants associated with the PrivateCertificatePrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificatePrototype_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificatePrototype_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificatePrototype_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificatePrototype_SecretType_Kv                 = "kv"
	PrivateCertificatePrototype_SecretType_PrivateCert        = "private_cert"
	PrivateCertificatePrototype_SecretType_PublicCert         = "public_cert"
	PrivateCertificatePrototype_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificatePrototype_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificatePrototype.Format property.
// The format of the returned data.
const (
	PrivateCertificatePrototype_Format_Pem       = "pem"
	PrivateCertificatePrototype_Format_PemBundle = "pem_bundle"
)

// Constants associated with the PrivateCertificatePrototype.PrivateKeyFormat property.
// The format of the generated private key.
const (
	PrivateCertificatePrototype_PrivateKeyFormat_Der   = "der"
	PrivateCertificatePrototype_PrivateKeyFormat_Pkcs8 = "pkcs8"
)

// NewPrivateCertificatePrototype : Instantiate PrivateCertificatePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificatePrototype(secretType string, name string, certificateTemplate string, commonName string) (_model *PrivateCertificatePrototype, err error) {
	_model = &PrivateCertificatePrototype{
		SecretType:          core.StringPtr(secretType),
		Name:                core.StringPtr(name),
		CertificateTemplate: core.StringPtr(certificateTemplate),
		CommonName:          core.StringPtr(commonName),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificatePrototype) isaSecretPrototype() bool {
	return true
}

// UnmarshalPrivateCertificatePrototype unmarshals an instance of PrivateCertificatePrototype from the specified map of raw messages.
func UnmarshalPrivateCertificatePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificatePrototype)
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_template", &obj.CertificateTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IpSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.UriSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "other_sans", &obj.OtherSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_format", &obj.PrivateKeyFormat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateVersion : Your private certificate version.
// This model "extends" SecretVersion
type PrivateCertificateVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date" validate:"required"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number" validate:"required"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity" validate:"required"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate" validate:"required"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key" validate:"required"`

	// The PEM-encoded certificate of the certificate authority that signed and issued this certificate.
	IssuingCa *string `json:"issuing_ca,omitempty"`

	// The chain of certificate authorities that are associated with the certificate.
	CaChain []string `json:"ca_chain,omitempty"`
}

// Constants associated with the PrivateCertificateVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateVersion_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateVersion_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateVersion_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateVersion_SecretType_Kv                 = "kv"
	PrivateCertificateVersion_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateVersion_SecretType_PublicCert         = "public_cert"
	PrivateCertificateVersion_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificateVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	PrivateCertificateVersion_Alias_Current  = "current"
	PrivateCertificateVersion_Alias_Previous = "previous"
)

func (*PrivateCertificateVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalPrivateCertificateVersion unmarshals an instance of PrivateCertificateVersion from the specified map of raw messages.
func UnmarshalPrivateCertificateVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateVersion)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuing_ca", &obj.IssuingCa)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ca_chain", &obj.CaChain)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateVersionActionRevoke : The response body to specify the properties of the action to revoke the private certificate.
// This model "extends" VersionAction
type PrivateCertificateVersionActionRevoke struct {
	// The type of secret version action.
	ActionType *string `json:"action_type" validate:"required"`

	// The timestamp of the certificate revocation.
	RevocationTimeSeconds *int64 `json:"revocation_time_seconds,omitempty"`
}

// Constants associated with the PrivateCertificateVersionActionRevoke.ActionType property.
// The type of secret version action.
const (
	PrivateCertificateVersionActionRevoke_ActionType_PrivateCertActionRevokeCertificate = "private_cert_action_revoke_certificate"
)

func (*PrivateCertificateVersionActionRevoke) isaVersionAction() bool {
	return true
}

// UnmarshalPrivateCertificateVersionActionRevoke unmarshals an instance of PrivateCertificateVersionActionRevoke from the specified map of raw messages.
func UnmarshalPrivateCertificateVersionActionRevoke(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateVersionActionRevoke)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_seconds", &obj.RevocationTimeSeconds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateVersionActionRevokePrototype : The request body to specify the properties of the action to revoke the private certificate.
// This model "extends" SecretVersionActionPrototype
type PrivateCertificateVersionActionRevokePrototype struct {
	// The type of secret version action.
	ActionType *string `json:"action_type" validate:"required"`
}

// Constants associated with the PrivateCertificateVersionActionRevokePrototype.ActionType property.
// The type of secret version action.
const (
	PrivateCertificateVersionActionRevokePrototype_ActionType_PrivateCertActionRevokeCertificate = "private_cert_action_revoke_certificate"
)

// NewPrivateCertificateVersionActionRevokePrototype : Instantiate PrivateCertificateVersionActionRevokePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPrivateCertificateVersionActionRevokePrototype(actionType string) (_model *PrivateCertificateVersionActionRevokePrototype, err error) {
	_model = &PrivateCertificateVersionActionRevokePrototype{
		ActionType: core.StringPtr(actionType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateVersionActionRevokePrototype) isaSecretVersionActionPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateVersionActionRevokePrototype unmarshals an instance of PrivateCertificateVersionActionRevokePrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateVersionActionRevokePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateVersionActionRevokePrototype)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateVersionMetadata : Properties of the version metadata of your private certificate.
// This model "extends" SecretVersionMetadata
type PrivateCertificateVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date" validate:"required"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number" validate:"required"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity" validate:"required"`
}

// Constants associated with the PrivateCertificateVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PrivateCertificateVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	PrivateCertificateVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	PrivateCertificateVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	PrivateCertificateVersionMetadata_SecretType_Kv                 = "kv"
	PrivateCertificateVersionMetadata_SecretType_PrivateCert        = "private_cert"
	PrivateCertificateVersionMetadata_SecretType_PublicCert         = "public_cert"
	PrivateCertificateVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	PrivateCertificateVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PrivateCertificateVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	PrivateCertificateVersionMetadata_Alias_Current  = "current"
	PrivateCertificateVersionMetadata_Alias_Previous = "previous"
)

func (*PrivateCertificateVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalPrivateCertificateVersionMetadata unmarshals an instance of PrivateCertificateVersionMetadata from the specified map of raw messages.
func UnmarshalPrivateCertificateVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateVersionMetadata)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateVersionPrototype : PrivateCertificateVersionPrototype struct
// This model "extends" SecretVersionPrototype
type PrivateCertificateVersionPrototype struct {
	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The certificate signing request.
	Csr *string `json:"csr,omitempty"`
}

func (*PrivateCertificateVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalPrivateCertificateVersionPrototype unmarshals an instance of PrivateCertificateVersionPrototype from the specified map of raw messages.
func UnmarshalPrivateCertificateVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateVersionPrototype)
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificate : Your public certificate.
// This model "extends" Secret
type PublicCertificate struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The identifier for the cryptographic algorithm that is used by the issuing certificate authority to sign a
	// certificate.
	SigningAlgorithm *string `json:"signing_algorithm,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The Common Name (CN) represents the server name protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// Issuance information that is associated with your certificate.
	IssuanceInfo *CertificateIssuanceInfo `json:"issuance_info,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The identifier for the cryptographic algorithm that is used to generate the public key that is associated with the
	// certificate.
	//
	// The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to
	// generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide
	// more encryption protection. Allowed values:  `RSA2048`, `RSA4096`, `ECDSA256`, and `ECDSA384`.
	KeyAlgorithm *string `json:"key_algorithm" validate:"required"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation" validate:"required"`

	// Indicates whether the issued certificate is bundled with intermediate certificates.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The name of the certificate authority configuration.
	Ca *string `json:"ca,omitempty"`

	// The name of the DNS provider configuration.
	Dns *string `json:"dns,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`
}

// Constants associated with the PublicCertificate.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificate_SecretType_Arbitrary          = "arbitrary"
	PublicCertificate_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificate_SecretType_ImportedCert       = "imported_cert"
	PublicCertificate_SecretType_Kv                 = "kv"
	PublicCertificate_SecretType_PrivateCert        = "private_cert"
	PublicCertificate_SecretType_PublicCert         = "public_cert"
	PublicCertificate_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificate_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PublicCertificate.StateDescription property.
// A text representation of the secret state.
const (
	PublicCertificate_StateDescription_Active        = "active"
	PublicCertificate_StateDescription_Deactivated   = "deactivated"
	PublicCertificate_StateDescription_Destroyed     = "destroyed"
	PublicCertificate_StateDescription_PreActivation = "pre_activation"
	PublicCertificate_StateDescription_Suspended     = "suspended"
)

func (*PublicCertificate) isaSecret() bool {
	return true
}

// UnmarshalPublicCertificate unmarshals an instance of PublicCertificate from the specified map of raw messages.
func UnmarshalPublicCertificate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificate)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_algorithm", &obj.SigningAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "issuance_info", &obj.IssuanceInfo, UnmarshalCertificateIssuanceInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bundle_certs", &obj.BundleCerts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ca", &obj.Ca)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dns", &obj.Dns)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateActionValidateManualDNS : The response body of the action to validate manual DNS challenges for the public certificate.
// This model "extends" SecretAction
type PublicCertificateActionValidateManualDNS struct {
	// The type of secret action.
	ActionType *string `json:"action_type" validate:"required"`
}

// Constants associated with the PublicCertificateActionValidateManualDNS.ActionType property.
// The type of secret action.
const (
	PublicCertificateActionValidateManualDNS_ActionType_PrivateCertActionRevokeCertificate   = "private_cert_action_revoke_certificate"
	PublicCertificateActionValidateManualDNS_ActionType_PublicCertActionValidateDnsChallenge = "public_cert_action_validate_dns_challenge"
)

func (*PublicCertificateActionValidateManualDNS) isaSecretAction() bool {
	return true
}

// UnmarshalPublicCertificateActionValidateManualDNS unmarshals an instance of PublicCertificateActionValidateManualDNS from the specified map of raw messages.
func UnmarshalPublicCertificateActionValidateManualDNS(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateActionValidateManualDNS)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateActionValidateManualDNSPrototype : The request body to specify the properties of the action to validate manual DNS challenges for the public
// certificate.
// This model "extends" SecretActionPrototype
type PublicCertificateActionValidateManualDNSPrototype struct {
	// The type of secret action.
	ActionType *string `json:"action_type" validate:"required"`
}

// Constants associated with the PublicCertificateActionValidateManualDNSPrototype.ActionType property.
// The type of secret action.
const (
	PublicCertificateActionValidateManualDNSPrototype_ActionType_PrivateCertActionRevokeCertificate   = "private_cert_action_revoke_certificate"
	PublicCertificateActionValidateManualDNSPrototype_ActionType_PublicCertActionValidateDnsChallenge = "public_cert_action_validate_dns_challenge"
)

// NewPublicCertificateActionValidateManualDNSPrototype : Instantiate PublicCertificateActionValidateManualDNSPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificateActionValidateManualDNSPrototype(actionType string) (_model *PublicCertificateActionValidateManualDNSPrototype, err error) {
	_model = &PublicCertificateActionValidateManualDNSPrototype{
		ActionType: core.StringPtr(actionType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateActionValidateManualDNSPrototype) isaSecretActionPrototype() bool {
	return true
}

// UnmarshalPublicCertificateActionValidateManualDNSPrototype unmarshals an instance of PublicCertificateActionValidateManualDNSPrototype from the specified map of raw messages.
func UnmarshalPublicCertificateActionValidateManualDNSPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateActionValidateManualDNSPrototype)
	err = core.UnmarshalPrimitive(m, "action_type", &obj.ActionType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationCALetsEncrypt : Properties that describe a Let's Encrypt CA configuration.
// This model "extends" Configuration
type PublicCertificateConfigurationCALetsEncrypt struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The configuration of the Let's Encrypt CA environment.
	LetsEncryptEnvironment *string `json:"lets_encrypt_environment" validate:"required"`

	// If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name. If
	// no match, the default offered chain will be used.
	LetsEncryptPreferredChain *string `json:"lets_encrypt_preferred_chain,omitempty"`

	// The PEM-encoded private key of your Let's Encrypt account. The data must be formatted on a single line with embedded
	// newline characters.
	LetsEncryptPrivateKey *string `json:"lets_encrypt_private_key" validate:"required"`
}

// Constants associated with the PublicCertificateConfigurationCALetsEncrypt.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationCALetsEncrypt_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationCALetsEncrypt_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationCALetsEncrypt_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationCALetsEncrypt_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationCALetsEncrypt_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationCALetsEncrypt_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationCALetsEncrypt_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PublicCertificateConfigurationCALetsEncrypt.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateConfigurationCALetsEncrypt_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateConfigurationCALetsEncrypt_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateConfigurationCALetsEncrypt_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateConfigurationCALetsEncrypt_SecretType_Kv                 = "kv"
	PublicCertificateConfigurationCALetsEncrypt_SecretType_PrivateCert        = "private_cert"
	PublicCertificateConfigurationCALetsEncrypt_SecretType_PublicCert         = "public_cert"
	PublicCertificateConfigurationCALetsEncrypt_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateConfigurationCALetsEncrypt_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PublicCertificateConfigurationCALetsEncrypt.LetsEncryptEnvironment property.
// The configuration of the Let's Encrypt CA environment.
const (
	PublicCertificateConfigurationCALetsEncrypt_LetsEncryptEnvironment_Production = "production"
	PublicCertificateConfigurationCALetsEncrypt_LetsEncryptEnvironment_Staging    = "staging"
)

func (*PublicCertificateConfigurationCALetsEncrypt) isaConfiguration() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationCALetsEncrypt unmarshals an instance of PublicCertificateConfigurationCALetsEncrypt from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationCALetsEncrypt(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationCALetsEncrypt)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "lets_encrypt_environment", &obj.LetsEncryptEnvironment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_preferred_chain", &obj.LetsEncryptPreferredChain)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_private_key", &obj.LetsEncryptPrivateKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationCALetsEncryptMetadata : Your Let's Encrypt CA metadata properties.
// This model "extends" ConfigurationMetadata
type PublicCertificateConfigurationCALetsEncryptMetadata struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The configuration of the Let's Encrypt CA environment.
	LetsEncryptEnvironment *string `json:"lets_encrypt_environment" validate:"required"`

	// If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name. If
	// no match, the default offered chain will be used.
	LetsEncryptPreferredChain *string `json:"lets_encrypt_preferred_chain,omitempty"`
}

// Constants associated with the PublicCertificateConfigurationCALetsEncryptMetadata.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationCALetsEncryptMetadata_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationCALetsEncryptMetadata_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationCALetsEncryptMetadata_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationCALetsEncryptMetadata_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationCALetsEncryptMetadata_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationCALetsEncryptMetadata_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationCALetsEncryptMetadata_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PublicCertificateConfigurationCALetsEncryptMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateConfigurationCALetsEncryptMetadata_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateConfigurationCALetsEncryptMetadata_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateConfigurationCALetsEncryptMetadata_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateConfigurationCALetsEncryptMetadata_SecretType_Kv                 = "kv"
	PublicCertificateConfigurationCALetsEncryptMetadata_SecretType_PrivateCert        = "private_cert"
	PublicCertificateConfigurationCALetsEncryptMetadata_SecretType_PublicCert         = "public_cert"
	PublicCertificateConfigurationCALetsEncryptMetadata_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateConfigurationCALetsEncryptMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PublicCertificateConfigurationCALetsEncryptMetadata.LetsEncryptEnvironment property.
// The configuration of the Let's Encrypt CA environment.
const (
	PublicCertificateConfigurationCALetsEncryptMetadata_LetsEncryptEnvironment_Production = "production"
	PublicCertificateConfigurationCALetsEncryptMetadata_LetsEncryptEnvironment_Staging    = "staging"
)

func (*PublicCertificateConfigurationCALetsEncryptMetadata) isaConfigurationMetadata() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationCALetsEncryptMetadata unmarshals an instance of PublicCertificateConfigurationCALetsEncryptMetadata from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationCALetsEncryptMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationCALetsEncryptMetadata)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "lets_encrypt_environment", &obj.LetsEncryptEnvironment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_preferred_chain", &obj.LetsEncryptPreferredChain)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationCALetsEncryptPatch : The configuration update of the Let's Encrypt Certificate Authority.
// This model "extends" ConfigurationPatch
type PublicCertificateConfigurationCALetsEncryptPatch struct {
	// The configuration of the Let's Encrypt CA environment.
	LetsEncryptEnvironment *string `json:"lets_encrypt_environment" validate:"required"`

	// The PEM-encoded private key of your Let's Encrypt account. The data must be formatted on a single line with embedded
	// newline characters.
	LetsEncryptPrivateKey *string `json:"lets_encrypt_private_key,omitempty"`

	// If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name. If
	// no match, the default offered chain will be used.
	LetsEncryptPreferredChain *string `json:"lets_encrypt_preferred_chain,omitempty"`
}

// Constants associated with the PublicCertificateConfigurationCALetsEncryptPatch.LetsEncryptEnvironment property.
// The configuration of the Let's Encrypt CA environment.
const (
	PublicCertificateConfigurationCALetsEncryptPatch_LetsEncryptEnvironment_Production = "production"
	PublicCertificateConfigurationCALetsEncryptPatch_LetsEncryptEnvironment_Staging    = "staging"
)

// NewPublicCertificateConfigurationCALetsEncryptPatch : Instantiate PublicCertificateConfigurationCALetsEncryptPatch (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificateConfigurationCALetsEncryptPatch(letsEncryptEnvironment string) (_model *PublicCertificateConfigurationCALetsEncryptPatch, err error) {
	_model = &PublicCertificateConfigurationCALetsEncryptPatch{
		LetsEncryptEnvironment: core.StringPtr(letsEncryptEnvironment),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateConfigurationCALetsEncryptPatch) isaConfigurationPatch() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationCALetsEncryptPatch unmarshals an instance of PublicCertificateConfigurationCALetsEncryptPatch from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationCALetsEncryptPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationCALetsEncryptPatch)
	err = core.UnmarshalPrimitive(m, "lets_encrypt_environment", &obj.LetsEncryptEnvironment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_private_key", &obj.LetsEncryptPrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_preferred_chain", &obj.LetsEncryptPreferredChain)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the PublicCertificateConfigurationCALetsEncryptPatch
func (publicCertificateConfigurationCALetsEncryptPatch *PublicCertificateConfigurationCALetsEncryptPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(publicCertificateConfigurationCALetsEncryptPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// PublicCertificateConfigurationCALetsEncryptPrototype : The properties of the Let's Encrypt CA configuration.
// This model "extends" ConfigurationPrototype
type PublicCertificateConfigurationCALetsEncryptPrototype struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// A human-readable unique name to assign to your configuration.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.
	Name *string `json:"name" validate:"required"`

	// The configuration of the Let's Encrypt CA environment.
	LetsEncryptEnvironment *string `json:"lets_encrypt_environment" validate:"required"`

	// The PEM-encoded private key of your Let's Encrypt account. The data must be formatted on a single line with embedded
	// newline characters.
	LetsEncryptPrivateKey *string `json:"lets_encrypt_private_key" validate:"required"`

	// If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name. If
	// no match, the default offered chain will be used.
	LetsEncryptPreferredChain *string `json:"lets_encrypt_preferred_chain,omitempty"`
}

// Constants associated with the PublicCertificateConfigurationCALetsEncryptPrototype.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationCALetsEncryptPrototype_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationCALetsEncryptPrototype_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationCALetsEncryptPrototype_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationCALetsEncryptPrototype_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationCALetsEncryptPrototype_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationCALetsEncryptPrototype_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationCALetsEncryptPrototype_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PublicCertificateConfigurationCALetsEncryptPrototype.LetsEncryptEnvironment property.
// The configuration of the Let's Encrypt CA environment.
const (
	PublicCertificateConfigurationCALetsEncryptPrototype_LetsEncryptEnvironment_Production = "production"
	PublicCertificateConfigurationCALetsEncryptPrototype_LetsEncryptEnvironment_Staging    = "staging"
)

// NewPublicCertificateConfigurationCALetsEncryptPrototype : Instantiate PublicCertificateConfigurationCALetsEncryptPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificateConfigurationCALetsEncryptPrototype(configType string, name string, letsEncryptEnvironment string, letsEncryptPrivateKey string) (_model *PublicCertificateConfigurationCALetsEncryptPrototype, err error) {
	_model = &PublicCertificateConfigurationCALetsEncryptPrototype{
		ConfigType:             core.StringPtr(configType),
		Name:                   core.StringPtr(name),
		LetsEncryptEnvironment: core.StringPtr(letsEncryptEnvironment),
		LetsEncryptPrivateKey:  core.StringPtr(letsEncryptPrivateKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateConfigurationCALetsEncryptPrototype) isaConfigurationPrototype() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationCALetsEncryptPrototype unmarshals an instance of PublicCertificateConfigurationCALetsEncryptPrototype from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationCALetsEncryptPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationCALetsEncryptPrototype)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_environment", &obj.LetsEncryptEnvironment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_private_key", &obj.LetsEncryptPrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lets_encrypt_preferred_chain", &obj.LetsEncryptPreferredChain)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationDNSClassicInfrastructure : Properties that describe a Classic Infrastructure DNS configuration.
// This model "extends" Configuration
type PublicCertificateConfigurationDNSClassicInfrastructure struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The username that is associated with your classic infrastructure account.
	//
	// In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information,
	// see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructureUsername *string `json:"classic_infrastructure_username" validate:"required"`

	// Your classic infrastructure API key.
	//
	// For information about viewing and accessing your classic infrastructure API key, see the
	// [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructurePassword *string `json:"classic_infrastructure_password" validate:"required"`
}

// Constants associated with the PublicCertificateConfigurationDNSClassicInfrastructure.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationDNSClassicInfrastructure_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationDNSClassicInfrastructure_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationDNSClassicInfrastructure_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationDNSClassicInfrastructure_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationDNSClassicInfrastructure_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationDNSClassicInfrastructure_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationDNSClassicInfrastructure_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PublicCertificateConfigurationDNSClassicInfrastructure.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateConfigurationDNSClassicInfrastructure_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateConfigurationDNSClassicInfrastructure_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateConfigurationDNSClassicInfrastructure_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateConfigurationDNSClassicInfrastructure_SecretType_Kv                 = "kv"
	PublicCertificateConfigurationDNSClassicInfrastructure_SecretType_PrivateCert        = "private_cert"
	PublicCertificateConfigurationDNSClassicInfrastructure_SecretType_PublicCert         = "public_cert"
	PublicCertificateConfigurationDNSClassicInfrastructure_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateConfigurationDNSClassicInfrastructure_SecretType_UsernamePassword   = "username_password"
)

func (*PublicCertificateConfigurationDNSClassicInfrastructure) isaConfiguration() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationDNSClassicInfrastructure unmarshals an instance of PublicCertificateConfigurationDNSClassicInfrastructure from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationDNSClassicInfrastructure(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationDNSClassicInfrastructure)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "classic_infrastructure_username", &obj.ClassicInfrastructureUsername)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "classic_infrastructure_password", &obj.ClassicInfrastructurePassword)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationDNSClassicInfrastructureMetadata : Your Classic Infrastructure DNS metadata properties.
// This model "extends" ConfigurationMetadata
type PublicCertificateConfigurationDNSClassicInfrastructureMetadata struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`
}

// Constants associated with the PublicCertificateConfigurationDNSClassicInfrastructureMetadata.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PublicCertificateConfigurationDNSClassicInfrastructureMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_SecretType_Kv                 = "kv"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_SecretType_PrivateCert        = "private_cert"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_SecretType_PublicCert         = "public_cert"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateConfigurationDNSClassicInfrastructureMetadata_SecretType_UsernamePassword   = "username_password"
)

func (*PublicCertificateConfigurationDNSClassicInfrastructureMetadata) isaConfigurationMetadata() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationDNSClassicInfrastructureMetadata unmarshals an instance of PublicCertificateConfigurationDNSClassicInfrastructureMetadata from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationDNSClassicInfrastructureMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationDNSClassicInfrastructureMetadata)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationDNSClassicInfrastructurePatch : Properties that describe the configuration update of an IBM Cloud classic infrastructure (SoftLayer).
// This model "extends" ConfigurationPatch
type PublicCertificateConfigurationDNSClassicInfrastructurePatch struct {
	// The username that is associated with your classic infrastructure account.
	//
	// In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information,
	// see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructureUsername *string `json:"classic_infrastructure_username,omitempty"`

	// Your classic infrastructure API key.
	//
	// For information about viewing and accessing your classic infrastructure API key, see the
	// [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructurePassword *string `json:"classic_infrastructure_password,omitempty"`
}

func (*PublicCertificateConfigurationDNSClassicInfrastructurePatch) isaConfigurationPatch() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationDNSClassicInfrastructurePatch unmarshals an instance of PublicCertificateConfigurationDNSClassicInfrastructurePatch from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationDNSClassicInfrastructurePatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationDNSClassicInfrastructurePatch)
	err = core.UnmarshalPrimitive(m, "classic_infrastructure_username", &obj.ClassicInfrastructureUsername)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "classic_infrastructure_password", &obj.ClassicInfrastructurePassword)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the PublicCertificateConfigurationDNSClassicInfrastructurePatch
func (publicCertificateConfigurationDNSClassicInfrastructurePatch *PublicCertificateConfigurationDNSClassicInfrastructurePatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(publicCertificateConfigurationDNSClassicInfrastructurePatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// PublicCertificateConfigurationDNSClassicInfrastructurePrototype : PublicCertificateConfigurationDNSClassicInfrastructurePrototype struct
// This model "extends" ConfigurationPrototype
type PublicCertificateConfigurationDNSClassicInfrastructurePrototype struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// A human-readable unique name to assign to your configuration.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.
	Name *string `json:"name" validate:"required"`

	// The username that is associated with your classic infrastructure account.
	//
	// In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information,
	// see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructureUsername *string `json:"classic_infrastructure_username" validate:"required"`

	// Your classic infrastructure API key.
	//
	// For information about viewing and accessing your classic infrastructure API key, see the
	// [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	ClassicInfrastructurePassword *string `json:"classic_infrastructure_password" validate:"required"`
}

// Constants associated with the PublicCertificateConfigurationDNSClassicInfrastructurePrototype.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationDNSClassicInfrastructurePrototype_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationDNSClassicInfrastructurePrototype_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationDNSClassicInfrastructurePrototype_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationDNSClassicInfrastructurePrototype_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationDNSClassicInfrastructurePrototype_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationDNSClassicInfrastructurePrototype_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationDNSClassicInfrastructurePrototype_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// NewPublicCertificateConfigurationDNSClassicInfrastructurePrototype : Instantiate PublicCertificateConfigurationDNSClassicInfrastructurePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificateConfigurationDNSClassicInfrastructurePrototype(configType string, name string, classicInfrastructureUsername string, classicInfrastructurePassword string) (_model *PublicCertificateConfigurationDNSClassicInfrastructurePrototype, err error) {
	_model = &PublicCertificateConfigurationDNSClassicInfrastructurePrototype{
		ConfigType:                    core.StringPtr(configType),
		Name:                          core.StringPtr(name),
		ClassicInfrastructureUsername: core.StringPtr(classicInfrastructureUsername),
		ClassicInfrastructurePassword: core.StringPtr(classicInfrastructurePassword),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateConfigurationDNSClassicInfrastructurePrototype) isaConfigurationPrototype() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationDNSClassicInfrastructurePrototype unmarshals an instance of PublicCertificateConfigurationDNSClassicInfrastructurePrototype from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationDNSClassicInfrastructurePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationDNSClassicInfrastructurePrototype)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "classic_infrastructure_username", &obj.ClassicInfrastructureUsername)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "classic_infrastructure_password", &obj.ClassicInfrastructurePassword)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationDNSCloudInternetServices : Properties that describe a Cloud Internet Services DNS configuration.
// This model "extends" Configuration
type PublicCertificateConfigurationDNSCloudInternetServices struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// An IBM Cloud API key that can to list domains in your Cloud Internet Services instance.
	//
	// To grant Secrets Manager the ability to view the Cloud Internet Services instance and all of its domains, the API
	// key must be assigned the Reader service role on Internet Services (`internet-svcs`).
	//
	// If you need to manage specific domains, you can assign the Manager role. For production environments, it is
	// recommended that you assign the Reader access role, and then use the
	// [IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific
	// domains. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	CloudInternetServicesApikey *string `json:"cloud_internet_services_apikey,omitempty"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	CloudInternetServicesCrn *string `json:"cloud_internet_services_crn" validate:"required"`
}

// Constants associated with the PublicCertificateConfigurationDNSCloudInternetServices.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationDNSCloudInternetServices_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationDNSCloudInternetServices_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationDNSCloudInternetServices_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationDNSCloudInternetServices_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationDNSCloudInternetServices_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationDNSCloudInternetServices_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationDNSCloudInternetServices_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PublicCertificateConfigurationDNSCloudInternetServices.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateConfigurationDNSCloudInternetServices_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateConfigurationDNSCloudInternetServices_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateConfigurationDNSCloudInternetServices_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateConfigurationDNSCloudInternetServices_SecretType_Kv                 = "kv"
	PublicCertificateConfigurationDNSCloudInternetServices_SecretType_PrivateCert        = "private_cert"
	PublicCertificateConfigurationDNSCloudInternetServices_SecretType_PublicCert         = "public_cert"
	PublicCertificateConfigurationDNSCloudInternetServices_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateConfigurationDNSCloudInternetServices_SecretType_UsernamePassword   = "username_password"
)

func (*PublicCertificateConfigurationDNSCloudInternetServices) isaConfiguration() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationDNSCloudInternetServices unmarshals an instance of PublicCertificateConfigurationDNSCloudInternetServices from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationDNSCloudInternetServices(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationDNSCloudInternetServices)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "cloud_internet_services_apikey", &obj.CloudInternetServicesApikey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_internet_services_crn", &obj.CloudInternetServicesCrn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationDNSCloudInternetServicesMetadata : Your Cloud Internet Services DNS metadata properties.
// This model "extends" ConfigurationMetadata
type PublicCertificateConfigurationDNSCloudInternetServicesMetadata struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// The unique name of your configuration.
	Name *string `json:"name" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`
}

// Constants associated with the PublicCertificateConfigurationDNSCloudInternetServicesMetadata.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// Constants associated with the PublicCertificateConfigurationDNSCloudInternetServicesMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_SecretType_Kv                 = "kv"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_SecretType_PrivateCert        = "private_cert"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_SecretType_PublicCert         = "public_cert"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateConfigurationDNSCloudInternetServicesMetadata_SecretType_UsernamePassword   = "username_password"
)

func (*PublicCertificateConfigurationDNSCloudInternetServicesMetadata) isaConfigurationMetadata() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationDNSCloudInternetServicesMetadata unmarshals an instance of PublicCertificateConfigurationDNSCloudInternetServicesMetadata from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationDNSCloudInternetServicesMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationDNSCloudInternetServicesMetadata)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateConfigurationDNSCloudInternetServicesPatch : The configuration update of the Cloud Internet Services DNS.
// This model "extends" ConfigurationPatch
type PublicCertificateConfigurationDNSCloudInternetServicesPatch struct {
	// An IBM Cloud API key that can to list domains in your Cloud Internet Services instance.
	//
	// To grant Secrets Manager the ability to view the Cloud Internet Services instance and all of its domains, the API
	// key must be assigned the Reader service role on Internet Services (`internet-svcs`).
	//
	// If you need to manage specific domains, you can assign the Manager role. For production environments, it is
	// recommended that you assign the Reader access role, and then use the
	// [IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific
	// domains. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	CloudInternetServicesApikey *string `json:"cloud_internet_services_apikey" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	CloudInternetServicesCrn *string `json:"cloud_internet_services_crn,omitempty"`
}

// NewPublicCertificateConfigurationDNSCloudInternetServicesPatch : Instantiate PublicCertificateConfigurationDNSCloudInternetServicesPatch (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificateConfigurationDNSCloudInternetServicesPatch(cloudInternetServicesApikey string) (_model *PublicCertificateConfigurationDNSCloudInternetServicesPatch, err error) {
	_model = &PublicCertificateConfigurationDNSCloudInternetServicesPatch{
		CloudInternetServicesApikey: core.StringPtr(cloudInternetServicesApikey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateConfigurationDNSCloudInternetServicesPatch) isaConfigurationPatch() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationDNSCloudInternetServicesPatch unmarshals an instance of PublicCertificateConfigurationDNSCloudInternetServicesPatch from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationDNSCloudInternetServicesPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationDNSCloudInternetServicesPatch)
	err = core.UnmarshalPrimitive(m, "cloud_internet_services_apikey", &obj.CloudInternetServicesApikey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_internet_services_crn", &obj.CloudInternetServicesCrn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the PublicCertificateConfigurationDNSCloudInternetServicesPatch
func (publicCertificateConfigurationDNSCloudInternetServicesPatch *PublicCertificateConfigurationDNSCloudInternetServicesPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(publicCertificateConfigurationDNSCloudInternetServicesPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// PublicCertificateConfigurationDNSCloudInternetServicesPrototype : Specify the properties for Cloud Internet Services DNS configuration.
// This model "extends" ConfigurationPrototype
type PublicCertificateConfigurationDNSCloudInternetServicesPrototype struct {
	// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
	// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
	// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
	ConfigType *string `json:"config_type" validate:"required"`

	// A human-readable unique name to assign to your configuration.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.
	Name *string `json:"name" validate:"required"`

	// An IBM Cloud API key that can to list domains in your Cloud Internet Services instance.
	//
	// To grant Secrets Manager the ability to view the Cloud Internet Services instance and all of its domains, the API
	// key must be assigned the Reader service role on Internet Services (`internet-svcs`).
	//
	// If you need to manage specific domains, you can assign the Manager role. For production environments, it is
	// recommended that you assign the Reader access role, and then use the
	// [IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific
	// domains. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	CloudInternetServicesApikey *string `json:"cloud_internet_services_apikey,omitempty"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	CloudInternetServicesCrn *string `json:"cloud_internet_services_crn" validate:"required"`
}

// Constants associated with the PublicCertificateConfigurationDNSCloudInternetServicesPrototype.ConfigType property.
// The configuration type. Can be one of: iam_credentials_configuration, public_cert_configuration_ca_lets_encrypt,
// public_cert_configuration_dns_classic_infrastructure, public_cert_configuration_dns_cloud_internet_services,
// private_cert_configuration_root_ca, private_cert_configuration_intermediate_ca, private_cert_configuration_template.
const (
	PublicCertificateConfigurationDNSCloudInternetServicesPrototype_ConfigType_IamCredentialsConfiguration                     = "iam_credentials_configuration"
	PublicCertificateConfigurationDNSCloudInternetServicesPrototype_ConfigType_PrivateCertConfigurationIntermediateCa          = "private_cert_configuration_intermediate_ca"
	PublicCertificateConfigurationDNSCloudInternetServicesPrototype_ConfigType_PrivateCertConfigurationRootCa                  = "private_cert_configuration_root_ca"
	PublicCertificateConfigurationDNSCloudInternetServicesPrototype_ConfigType_PrivateCertConfigurationTemplate                = "private_cert_configuration_template"
	PublicCertificateConfigurationDNSCloudInternetServicesPrototype_ConfigType_PublicCertConfigurationCaLetsEncrypt            = "public_cert_configuration_ca_lets_encrypt"
	PublicCertificateConfigurationDNSCloudInternetServicesPrototype_ConfigType_PublicCertConfigurationDnsClassicInfrastructure = "public_cert_configuration_dns_classic_infrastructure"
	PublicCertificateConfigurationDNSCloudInternetServicesPrototype_ConfigType_PublicCertConfigurationDnsCloudInternetServices = "public_cert_configuration_dns_cloud_internet_services"
)

// NewPublicCertificateConfigurationDNSCloudInternetServicesPrototype : Instantiate PublicCertificateConfigurationDNSCloudInternetServicesPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificateConfigurationDNSCloudInternetServicesPrototype(configType string, name string, cloudInternetServicesCrn string) (_model *PublicCertificateConfigurationDNSCloudInternetServicesPrototype, err error) {
	_model = &PublicCertificateConfigurationDNSCloudInternetServicesPrototype{
		ConfigType:               core.StringPtr(configType),
		Name:                     core.StringPtr(name),
		CloudInternetServicesCrn: core.StringPtr(cloudInternetServicesCrn),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateConfigurationDNSCloudInternetServicesPrototype) isaConfigurationPrototype() bool {
	return true
}

// UnmarshalPublicCertificateConfigurationDNSCloudInternetServicesPrototype unmarshals an instance of PublicCertificateConfigurationDNSCloudInternetServicesPrototype from the specified map of raw messages.
func UnmarshalPublicCertificateConfigurationDNSCloudInternetServicesPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateConfigurationDNSCloudInternetServicesPrototype)
	err = core.UnmarshalPrimitive(m, "config_type", &obj.ConfigType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_internet_services_apikey", &obj.CloudInternetServicesApikey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_internet_services_crn", &obj.CloudInternetServicesCrn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateMetadata : Properties of the metadata of your public certificate.
// This model "extends" SecretMetadata
type PublicCertificateMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The identifier for the cryptographic algorithm that is used by the issuing certificate authority to sign a
	// certificate.
	SigningAlgorithm *string `json:"signing_algorithm,omitempty"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The Common Name (CN) represents the server name protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// Issuance information that is associated with your certificate.
	IssuanceInfo *CertificateIssuanceInfo `json:"issuance_info,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// The identifier for the cryptographic algorithm that is used to generate the public key that is associated with the
	// certificate.
	//
	// The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to
	// generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide
	// more encryption protection. Allowed values:  `RSA2048`, `RSA4096`, `ECDSA256`, and `ECDSA384`.
	KeyAlgorithm *string `json:"key_algorithm" validate:"required"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation" validate:"required"`

	// Indicates whether the issued certificate is bundled with intermediate certificates.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The name of the certificate authority configuration.
	Ca *string `json:"ca,omitempty"`

	// The name of the DNS provider configuration.
	Dns *string `json:"dns,omitempty"`
}

// Constants associated with the PublicCertificateMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateMetadata_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateMetadata_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateMetadata_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateMetadata_SecretType_Kv                 = "kv"
	PublicCertificateMetadata_SecretType_PrivateCert        = "private_cert"
	PublicCertificateMetadata_SecretType_PublicCert         = "public_cert"
	PublicCertificateMetadata_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PublicCertificateMetadata.StateDescription property.
// A text representation of the secret state.
const (
	PublicCertificateMetadata_StateDescription_Active        = "active"
	PublicCertificateMetadata_StateDescription_Deactivated   = "deactivated"
	PublicCertificateMetadata_StateDescription_Destroyed     = "destroyed"
	PublicCertificateMetadata_StateDescription_PreActivation = "pre_activation"
	PublicCertificateMetadata_StateDescription_Suspended     = "suspended"
)

func (*PublicCertificateMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalPublicCertificateMetadata unmarshals an instance of PublicCertificateMetadata from the specified map of raw messages.
func UnmarshalPublicCertificateMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateMetadata)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signing_algorithm", &obj.SigningAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "issuance_info", &obj.IssuanceInfo, UnmarshalCertificateIssuanceInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bundle_certs", &obj.BundleCerts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ca", &obj.Ca)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dns", &obj.Dns)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateMetadataPatch : PublicCertificateMetadataPatch struct
// This model "extends" SecretMetadataPatch
type PublicCertificateMetadataPatch struct {
	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`
}

func (*PublicCertificateMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

// UnmarshalPublicCertificateMetadataPatch unmarshals an instance of PublicCertificateMetadataPatch from the specified map of raw messages.
func UnmarshalPublicCertificateMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateMetadataPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the PublicCertificateMetadataPatch
func (publicCertificateMetadataPatch *PublicCertificateMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(publicCertificateMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// PublicCertificatePrototype : PublicCertificatePrototype struct
// This model "extends" SecretPrototype
type PublicCertificatePrototype struct {
	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The Common Name (CN) represents the server name protected by the SSL certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL
	// certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The identifier for the cryptographic algorithm that is used to generate the public key that is associated with the
	// certificate.
	//
	// The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to
	// generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide
	// more encryption protection. Allowed values:  `RSA2048`, `RSA4096`, `ECDSA256`, and `ECDSA384`.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The name of the certificate authority configuration.
	Ca *string `json:"ca" validate:"required"`

	// The name of the DNS provider configuration.
	Dns *string `json:"dns" validate:"required"`

	// This field indicates whether your issued certificate is bundled with intermediate certificates. Set to `false` for
	// the certificate file to contain only the issued certificate.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// Constants associated with the PublicCertificatePrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificatePrototype_SecretType_Arbitrary          = "arbitrary"
	PublicCertificatePrototype_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificatePrototype_SecretType_ImportedCert       = "imported_cert"
	PublicCertificatePrototype_SecretType_Kv                 = "kv"
	PublicCertificatePrototype_SecretType_PrivateCert        = "private_cert"
	PublicCertificatePrototype_SecretType_PublicCert         = "public_cert"
	PublicCertificatePrototype_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificatePrototype_SecretType_UsernamePassword   = "username_password"
)

// NewPublicCertificatePrototype : Instantiate PublicCertificatePrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificatePrototype(secretType string, name string, commonName string, ca string, dns string) (_model *PublicCertificatePrototype, err error) {
	_model = &PublicCertificatePrototype{
		SecretType: core.StringPtr(secretType),
		Name:       core.StringPtr(name),
		CommonName: core.StringPtr(commonName),
		Ca:         core.StringPtr(ca),
		Dns:        core.StringPtr(dns),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificatePrototype) isaSecretPrototype() bool {
	return true
}

// UnmarshalPublicCertificatePrototype unmarshals an instance of PublicCertificatePrototype from the specified map of raw messages.
func UnmarshalPublicCertificatePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificatePrototype)
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ca", &obj.Ca)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dns", &obj.Dns)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bundle_certs", &obj.BundleCerts)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateRotationPolicy : This field indicates whether Secrets Manager rotates your secrets automatically.
//
// For public certificates, if `auto_rotate` is set to `true`, the service reorders your certificate for 31 days, before
// it expires.
// This model "extends" RotationPolicy
type PublicCertificateRotationPolicy struct {
	// This field indicates whether Secrets Manager rotates your secret automatically.
	//
	// The default is `false`. If `auto_rotate` is set to `true` the service rotates your secret based on the defined
	// interval.
	AutoRotate *bool `json:"auto_rotate" validate:"required"`

	// This field indicates whether Secrets Manager rotates the private key for your public certificate automatically.
	//
	// The default is `false`. If it is set to `true`, the service generates and stores a new private key for your rotated
	// certificate.
	RotateKeys *bool `json:"rotate_keys" validate:"required"`
}

// NewPublicCertificateRotationPolicy : Instantiate PublicCertificateRotationPolicy (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificateRotationPolicy(autoRotate bool, rotateKeys bool) (_model *PublicCertificateRotationPolicy, err error) {
	_model = &PublicCertificateRotationPolicy{
		AutoRotate: core.BoolPtr(autoRotate),
		RotateKeys: core.BoolPtr(rotateKeys),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateRotationPolicy) isaRotationPolicy() bool {
	return true
}

// UnmarshalPublicCertificateRotationPolicy unmarshals an instance of PublicCertificateRotationPolicy from the specified map of raw messages.
func UnmarshalPublicCertificateRotationPolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateRotationPolicy)
	err = core.UnmarshalPrimitive(m, "auto_rotate", &obj.AutoRotate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rotate_keys", &obj.RotateKeys)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateVersion : Versions of your public certificate.
// This model "extends" SecretVersion
type PublicCertificateVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity,omitempty"`

	// Your PEM-encoded certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The PEM-encoded intermediate certificate that is associated with the root certificate. The data must be formatted on
	// a single line with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The PEM-encoded private key that is associated with the certificate. The data must be formatted on a single line
	// with embedded newline characters.
	PrivateKey *string `json:"private_key,omitempty"`
}

// Constants associated with the PublicCertificateVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateVersion_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateVersion_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateVersion_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateVersion_SecretType_Kv                 = "kv"
	PublicCertificateVersion_SecretType_PrivateCert        = "private_cert"
	PublicCertificateVersion_SecretType_PublicCert         = "public_cert"
	PublicCertificateVersion_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PublicCertificateVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	PublicCertificateVersion_Alias_Current  = "current"
	PublicCertificateVersion_Alias_Previous = "previous"
)

func (*PublicCertificateVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalPublicCertificateVersion unmarshals an instance of PublicCertificateVersion from the specified map of raw messages.
func UnmarshalPublicCertificateVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateVersion)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateVersionMetadata : Properties of the version metadata of your public certificate.
// This model "extends" SecretVersionMetadata
type PublicCertificateVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The unique serial number that was assigned to a certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date and time that the certificate validity period begins and ends.
	Validity *CertificateValidity `json:"validity,omitempty"`
}

// Constants associated with the PublicCertificateVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	PublicCertificateVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	PublicCertificateVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	PublicCertificateVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	PublicCertificateVersionMetadata_SecretType_Kv                 = "kv"
	PublicCertificateVersionMetadata_SecretType_PrivateCert        = "private_cert"
	PublicCertificateVersionMetadata_SecretType_PublicCert         = "public_cert"
	PublicCertificateVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	PublicCertificateVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the PublicCertificateVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	PublicCertificateVersionMetadata_Alias_Current  = "current"
	PublicCertificateVersionMetadata_Alias_Previous = "previous"
)

func (*PublicCertificateVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalPublicCertificateVersionMetadata unmarshals an instance of PublicCertificateVersionMetadata from the specified map of raw messages.
func UnmarshalPublicCertificateVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateVersionMetadata)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateVersionPrototype : PublicCertificateVersionPrototype struct
// This model "extends" SecretVersionPrototype
type PublicCertificateVersionPrototype struct {
	// Defines the rotation object that is used to manually rotate public certificates.
	Rotation *PublicCertificateRotationObject `json:"rotation" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewPublicCertificateVersionPrototype : Instantiate PublicCertificateVersionPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewPublicCertificateVersionPrototype(rotation *PublicCertificateRotationObject) (_model *PublicCertificateVersionPrototype, err error) {
	_model = &PublicCertificateVersionPrototype{
		Rotation: rotation,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalPublicCertificateVersionPrototype unmarshals an instance of PublicCertificateVersionPrototype from the specified map of raw messages.
func UnmarshalPublicCertificateVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateVersionPrototype)
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalPublicCertificateRotationObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSecret : Your service credentials secret.
// This model "extends" Secret
type ServiceCredentialsSecret struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`

	// The properties that are required to create the service credentials for the specified source service instance.
	SourceService *ServiceCredentialsSecretSourceService `json:"source_service" validate:"required"`

	// The properties of the service credentials secret payload.
	Credentials *ServiceCredentialsSecretCredentials `json:"credentials" validate:"required"`
}

// Constants associated with the ServiceCredentialsSecret.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ServiceCredentialsSecret_SecretType_Arbitrary          = "arbitrary"
	ServiceCredentialsSecret_SecretType_IamCredentials     = "iam_credentials"
	ServiceCredentialsSecret_SecretType_ImportedCert       = "imported_cert"
	ServiceCredentialsSecret_SecretType_Kv                 = "kv"
	ServiceCredentialsSecret_SecretType_PrivateCert        = "private_cert"
	ServiceCredentialsSecret_SecretType_PublicCert         = "public_cert"
	ServiceCredentialsSecret_SecretType_ServiceCredentials = "service_credentials"
	ServiceCredentialsSecret_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ServiceCredentialsSecret.StateDescription property.
// A text representation of the secret state.
const (
	ServiceCredentialsSecret_StateDescription_Active        = "active"
	ServiceCredentialsSecret_StateDescription_Deactivated   = "deactivated"
	ServiceCredentialsSecret_StateDescription_Destroyed     = "destroyed"
	ServiceCredentialsSecret_StateDescription_PreActivation = "pre_activation"
	ServiceCredentialsSecret_StateDescription_Suspended     = "suspended"
)

func (*ServiceCredentialsSecret) isaSecret() bool {
	return true
}

// UnmarshalServiceCredentialsSecret unmarshals an instance of ServiceCredentialsSecret from the specified map of raw messages.
func UnmarshalServiceCredentialsSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecret)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "source_service", &obj.SourceService, UnmarshalServiceCredentialsSecretSourceService)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "credentials", &obj.Credentials, UnmarshalServiceCredentialsSecretCredentials)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSecretMetadata : The metadata properties for your service credentials secret.
// This model "extends" SecretMetadata
type ServiceCredentialsSecretMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`

	// The properties that are required to create the service credentials for the specified source service instance.
	SourceService *ServiceCredentialsSecretSourceService `json:"source_service" validate:"required"`
}

// Constants associated with the ServiceCredentialsSecretMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ServiceCredentialsSecretMetadata_SecretType_Arbitrary          = "arbitrary"
	ServiceCredentialsSecretMetadata_SecretType_IamCredentials     = "iam_credentials"
	ServiceCredentialsSecretMetadata_SecretType_ImportedCert       = "imported_cert"
	ServiceCredentialsSecretMetadata_SecretType_Kv                 = "kv"
	ServiceCredentialsSecretMetadata_SecretType_PrivateCert        = "private_cert"
	ServiceCredentialsSecretMetadata_SecretType_PublicCert         = "public_cert"
	ServiceCredentialsSecretMetadata_SecretType_ServiceCredentials = "service_credentials"
	ServiceCredentialsSecretMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ServiceCredentialsSecretMetadata.StateDescription property.
// A text representation of the secret state.
const (
	ServiceCredentialsSecretMetadata_StateDescription_Active        = "active"
	ServiceCredentialsSecretMetadata_StateDescription_Deactivated   = "deactivated"
	ServiceCredentialsSecretMetadata_StateDescription_Destroyed     = "destroyed"
	ServiceCredentialsSecretMetadata_StateDescription_PreActivation = "pre_activation"
	ServiceCredentialsSecretMetadata_StateDescription_Suspended     = "suspended"
)

func (*ServiceCredentialsSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalServiceCredentialsSecretMetadata unmarshals an instance of ServiceCredentialsSecretMetadata from the specified map of raw messages.
func UnmarshalServiceCredentialsSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecretMetadata)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "source_service", &obj.SourceService, UnmarshalServiceCredentialsSecretSourceService)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSecretMetadataPatch : ServiceCredentialsSecretMetadataPatch struct
// This model "extends" SecretMetadataPatch
type ServiceCredentialsSecretMetadataPatch struct {
	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`
}

func (*ServiceCredentialsSecretMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

// UnmarshalServiceCredentialsSecretMetadataPatch unmarshals an instance of ServiceCredentialsSecretMetadataPatch from the specified map of raw messages.
func UnmarshalServiceCredentialsSecretMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecretMetadataPatch)
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the ServiceCredentialsSecretMetadataPatch
func (serviceCredentialsSecretMetadataPatch *ServiceCredentialsSecretMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(serviceCredentialsSecretMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// ServiceCredentialsSecretPrototype : ServiceCredentialsSecretPrototype struct
// This model "extends" SecretPrototype
type ServiceCredentialsSecretPrototype struct {
	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name" validate:"required"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The properties that are required to create the service credentials for the specified source service instance.
	SourceService *ServiceCredentialsSecretSourceService `json:"source_service" validate:"required"`

	// The time-to-live (TTL) or lease duration to assign to credentials that are generated. Supported secret types:
	// iam_credentials, service_credentials. The TTL defines how long generated credentials remain valid. The value can be
	// either an integer that specifies the number of seconds, or the string  representation of a duration, such as `1440m`
	// or `24h`. For the iam_credentials secret type, the TTL field is mandatory. The minimum duration is 1 minute. The
	// maximum is 90 days. For the service_credentials secret type, the TTL field is optional. If it is set the minimum
	// duration is 1 day. The maximum is 90 days. By default, the TTL is set to 0.
	TTL *string `json:"ttl,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// Constants associated with the ServiceCredentialsSecretPrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ServiceCredentialsSecretPrototype_SecretType_Arbitrary          = "arbitrary"
	ServiceCredentialsSecretPrototype_SecretType_IamCredentials     = "iam_credentials"
	ServiceCredentialsSecretPrototype_SecretType_ImportedCert       = "imported_cert"
	ServiceCredentialsSecretPrototype_SecretType_Kv                 = "kv"
	ServiceCredentialsSecretPrototype_SecretType_PrivateCert        = "private_cert"
	ServiceCredentialsSecretPrototype_SecretType_PublicCert         = "public_cert"
	ServiceCredentialsSecretPrototype_SecretType_ServiceCredentials = "service_credentials"
	ServiceCredentialsSecretPrototype_SecretType_UsernamePassword   = "username_password"
)

// NewServiceCredentialsSecretPrototype : Instantiate ServiceCredentialsSecretPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewServiceCredentialsSecretPrototype(name string, secretType string, sourceService *ServiceCredentialsSecretSourceService) (_model *ServiceCredentialsSecretPrototype, err error) {
	_model = &ServiceCredentialsSecretPrototype{
		Name:          core.StringPtr(name),
		SecretType:    core.StringPtr(secretType),
		SourceService: sourceService,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ServiceCredentialsSecretPrototype) isaSecretPrototype() bool {
	return true
}

// UnmarshalServiceCredentialsSecretPrototype unmarshals an instance of ServiceCredentialsSecretPrototype from the specified map of raw messages.
func UnmarshalServiceCredentialsSecretPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecretPrototype)
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "source_service", &obj.SourceService, UnmarshalServiceCredentialsSecretSourceService)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSecretVersion : Your service credentials secret version.
// This model "extends" SecretVersion
type ServiceCredentialsSecretVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The source service resource key data of the generated service credentials.
	ResourceKey *ServiceCredentialsResourceKey `json:"resource_key,omitempty"`

	// The properties of the service credentials secret payload.
	Credentials *ServiceCredentialsSecretCredentials `json:"credentials" validate:"required"`
}

// Constants associated with the ServiceCredentialsSecretVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ServiceCredentialsSecretVersion_SecretType_Arbitrary          = "arbitrary"
	ServiceCredentialsSecretVersion_SecretType_IamCredentials     = "iam_credentials"
	ServiceCredentialsSecretVersion_SecretType_ImportedCert       = "imported_cert"
	ServiceCredentialsSecretVersion_SecretType_Kv                 = "kv"
	ServiceCredentialsSecretVersion_SecretType_PrivateCert        = "private_cert"
	ServiceCredentialsSecretVersion_SecretType_PublicCert         = "public_cert"
	ServiceCredentialsSecretVersion_SecretType_ServiceCredentials = "service_credentials"
	ServiceCredentialsSecretVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ServiceCredentialsSecretVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	ServiceCredentialsSecretVersion_Alias_Current  = "current"
	ServiceCredentialsSecretVersion_Alias_Previous = "previous"
)

func (*ServiceCredentialsSecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalServiceCredentialsSecretVersion unmarshals an instance of ServiceCredentialsSecretVersion from the specified map of raw messages.
func UnmarshalServiceCredentialsSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecretVersion)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource_key", &obj.ResourceKey, UnmarshalServiceCredentialsResourceKey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "credentials", &obj.Credentials, UnmarshalServiceCredentialsSecretCredentials)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSecretVersionMetadata : The version metadata properties for your service credentials secret.
// This model "extends" SecretVersionMetadata
type ServiceCredentialsSecretVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The source service resource key data of the generated service credentials.
	ResourceKey *ServiceCredentialsResourceKey `json:"resource_key,omitempty"`
}

// Constants associated with the ServiceCredentialsSecretVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	ServiceCredentialsSecretVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	ServiceCredentialsSecretVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	ServiceCredentialsSecretVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	ServiceCredentialsSecretVersionMetadata_SecretType_Kv                 = "kv"
	ServiceCredentialsSecretVersionMetadata_SecretType_PrivateCert        = "private_cert"
	ServiceCredentialsSecretVersionMetadata_SecretType_PublicCert         = "public_cert"
	ServiceCredentialsSecretVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	ServiceCredentialsSecretVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the ServiceCredentialsSecretVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	ServiceCredentialsSecretVersionMetadata_Alias_Current  = "current"
	ServiceCredentialsSecretVersionMetadata_Alias_Previous = "previous"
)

func (*ServiceCredentialsSecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalServiceCredentialsSecretVersionMetadata unmarshals an instance of ServiceCredentialsSecretVersionMetadata from the specified map of raw messages.
func UnmarshalServiceCredentialsSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource_key", &obj.ResourceKey, UnmarshalServiceCredentialsResourceKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCredentialsSecretVersionPrototype : ServiceCredentialsSecretVersionPrototype struct
// This model "extends" SecretVersionPrototype
type ServiceCredentialsSecretVersionPrototype struct {
	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

func (*ServiceCredentialsSecretVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalServiceCredentialsSecretVersionPrototype unmarshals an instance of ServiceCredentialsSecretVersionPrototype from the specified map of raw messages.
func UnmarshalServiceCredentialsSecretVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCredentialsSecretVersionPrototype)
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsernamePasswordSecret : Your user credentials secret.
// This model "extends" Secret
type UsernamePasswordSecret struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// The username that is assigned to an `username_password` secret.
	Username *string `json:"username" validate:"required"`

	// The password that is assigned to an `username_password` secret.
	Password *string `json:"password" validate:"required"`
}

// Constants associated with the UsernamePasswordSecret.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	UsernamePasswordSecret_SecretType_Arbitrary          = "arbitrary"
	UsernamePasswordSecret_SecretType_IamCredentials     = "iam_credentials"
	UsernamePasswordSecret_SecretType_ImportedCert       = "imported_cert"
	UsernamePasswordSecret_SecretType_Kv                 = "kv"
	UsernamePasswordSecret_SecretType_PrivateCert        = "private_cert"
	UsernamePasswordSecret_SecretType_PublicCert         = "public_cert"
	UsernamePasswordSecret_SecretType_ServiceCredentials = "service_credentials"
	UsernamePasswordSecret_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the UsernamePasswordSecret.StateDescription property.
// A text representation of the secret state.
const (
	UsernamePasswordSecret_StateDescription_Active        = "active"
	UsernamePasswordSecret_StateDescription_Deactivated   = "deactivated"
	UsernamePasswordSecret_StateDescription_Destroyed     = "destroyed"
	UsernamePasswordSecret_StateDescription_PreActivation = "pre_activation"
	UsernamePasswordSecret_StateDescription_Suspended     = "suspended"
)

func (*UsernamePasswordSecret) isaSecret() bool {
	return true
}

// UnmarshalUsernamePasswordSecret unmarshals an instance of UsernamePasswordSecret from the specified map of raw messages.
func UnmarshalUsernamePasswordSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecret)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsernamePasswordSecretMetadata : Properties of the metadata of your user credentials secret.
// This model "extends" SecretMetadata
type UsernamePasswordSecretMetadata struct {
	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A CRN that uniquely identifies an IBM Cloud resource.
	Crn *string `json:"crn" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The number of locks of the secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The human-readable name of your secret.
	Name *string `json:"name,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// The secret state that is based on `NIST SP 800-57`. States are integers and correspond to the `Pre-activation = 0`,
	// `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The date when a resource was modified. The date format follows `RFC 3339`.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The number of versions of your secret.
	VersionsTotal *int64 `json:"versions_total" validate:"required"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`
}

// Constants associated with the UsernamePasswordSecretMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	UsernamePasswordSecretMetadata_SecretType_Arbitrary          = "arbitrary"
	UsernamePasswordSecretMetadata_SecretType_IamCredentials     = "iam_credentials"
	UsernamePasswordSecretMetadata_SecretType_ImportedCert       = "imported_cert"
	UsernamePasswordSecretMetadata_SecretType_Kv                 = "kv"
	UsernamePasswordSecretMetadata_SecretType_PrivateCert        = "private_cert"
	UsernamePasswordSecretMetadata_SecretType_PublicCert         = "public_cert"
	UsernamePasswordSecretMetadata_SecretType_ServiceCredentials = "service_credentials"
	UsernamePasswordSecretMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the UsernamePasswordSecretMetadata.StateDescription property.
// A text representation of the secret state.
const (
	UsernamePasswordSecretMetadata_StateDescription_Active        = "active"
	UsernamePasswordSecretMetadata_StateDescription_Deactivated   = "deactivated"
	UsernamePasswordSecretMetadata_StateDescription_Destroyed     = "destroyed"
	UsernamePasswordSecretMetadata_StateDescription_PreActivation = "pre_activation"
	UsernamePasswordSecretMetadata_StateDescription_Suspended     = "suspended"
)

func (*UsernamePasswordSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalUsernamePasswordSecretMetadata unmarshals an instance of UsernamePasswordSecretMetadata from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretMetadata)
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_description", &obj.StateDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsernamePasswordSecretMetadataPatch : UsernamePasswordSecretMetadataPatch struct
// This model "extends" SecretMetadataPatch
type UsernamePasswordSecretMetadataPatch struct {
	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`
}

func (*UsernamePasswordSecretMetadataPatch) isaSecretMetadataPatch() bool {
	return true
}

// UnmarshalUsernamePasswordSecretMetadataPatch unmarshals an instance of UsernamePasswordSecretMetadataPatch from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretMetadataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretMetadataPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the UsernamePasswordSecretMetadataPatch
func (usernamePasswordSecretMetadataPatch *UsernamePasswordSecretMetadataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(usernamePasswordSecretMetadataPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	return
}

// UsernamePasswordSecretPrototype : UsernamePasswordSecretPrototype struct
// This model "extends" SecretPrototype
type UsernamePasswordSecretPrototype struct {
	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A human-readable name to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to search secrets in your instance. Only 30 labels can be created.
	//
	// Label can be between 2-30 characters, including spaces.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The username that is assigned to an `username_password` secret.
	Username *string `json:"username" validate:"required"`

	// The password that is assigned to an `username_password` secret.
	Password *string `json:"password" validate:"required"`

	// The date when the secret material expires. The date format follows the `RFC 3339` format. Supported secret types:
	// Arbitrary, username_password.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// This field indicates whether Secrets Manager rotates your secrets automatically. Supported secret types:
	// username_password, private_cert, public_cert, iam_credentials.
	Rotation RotationPolicyIntf `json:"rotation,omitempty"`
}

// Constants associated with the UsernamePasswordSecretPrototype.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	UsernamePasswordSecretPrototype_SecretType_Arbitrary          = "arbitrary"
	UsernamePasswordSecretPrototype_SecretType_IamCredentials     = "iam_credentials"
	UsernamePasswordSecretPrototype_SecretType_ImportedCert       = "imported_cert"
	UsernamePasswordSecretPrototype_SecretType_Kv                 = "kv"
	UsernamePasswordSecretPrototype_SecretType_PrivateCert        = "private_cert"
	UsernamePasswordSecretPrototype_SecretType_PublicCert         = "public_cert"
	UsernamePasswordSecretPrototype_SecretType_ServiceCredentials = "service_credentials"
	UsernamePasswordSecretPrototype_SecretType_UsernamePassword   = "username_password"
)

// NewUsernamePasswordSecretPrototype : Instantiate UsernamePasswordSecretPrototype (Generic Model Constructor)
func (*SecretsManagerV2) NewUsernamePasswordSecretPrototype(secretType string, name string, username string, password string) (_model *UsernamePasswordSecretPrototype, err error) {
	_model = &UsernamePasswordSecretPrototype{
		SecretType: core.StringPtr(secretType),
		Name:       core.StringPtr(name),
		Username:   core.StringPtr(username),
		Password:   core.StringPtr(password),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*UsernamePasswordSecretPrototype) isaSecretPrototype() bool {
	return true
}

// UnmarshalUsernamePasswordSecretPrototype unmarshals an instance of UsernamePasswordSecretPrototype from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretPrototype)
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotationPolicy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsernamePasswordSecretVersion : Your user credentials secret version.
// This model "extends" SecretVersion
type UsernamePasswordSecretVersion struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`

	// The username that is assigned to an `username_password` secret.
	Username *string `json:"username" validate:"required"`

	// The password that is assigned to an `username_password` secret.
	Password *string `json:"password" validate:"required"`
}

// Constants associated with the UsernamePasswordSecretVersion.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	UsernamePasswordSecretVersion_SecretType_Arbitrary          = "arbitrary"
	UsernamePasswordSecretVersion_SecretType_IamCredentials     = "iam_credentials"
	UsernamePasswordSecretVersion_SecretType_ImportedCert       = "imported_cert"
	UsernamePasswordSecretVersion_SecretType_Kv                 = "kv"
	UsernamePasswordSecretVersion_SecretType_PrivateCert        = "private_cert"
	UsernamePasswordSecretVersion_SecretType_PublicCert         = "public_cert"
	UsernamePasswordSecretVersion_SecretType_ServiceCredentials = "service_credentials"
	UsernamePasswordSecretVersion_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the UsernamePasswordSecretVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	UsernamePasswordSecretVersion_Alias_Current  = "current"
	UsernamePasswordSecretVersion_Alias_Previous = "previous"
)

func (*UsernamePasswordSecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalUsernamePasswordSecretVersion unmarshals an instance of UsernamePasswordSecretVersion from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretVersion)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsernamePasswordSecretVersionMetadata : Properties of the version metadata of your user credentials secret.
// This model "extends" SecretVersionMetadata
type UsernamePasswordSecretVersionMetadata struct {
	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique identifier that is associated with the entity that created the secret.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the resource was created. The date format follows `RFC 3339`.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// This field indicates whether the secret data that is associated with a secret version was retrieved in a call to the
	// service API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// A v4 UUID identifier.
	ID *string `json:"id" validate:"required"`

	// The human-readable name of your secret.
	SecretName *string `json:"secret_name,omitempty"`

	// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
	// service_credentials, kv, and username_password.
	SecretType *string `json:"secret_type" validate:"required"`

	// A v4 UUID identifier, or `default` secret group.
	SecretGroupID *string `json:"secret_group_id" validate:"required"`

	// Indicates whether the secret payload is available in this secret version.
	PayloadAvailable *bool `json:"payload_available" validate:"required"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// A v4 UUID identifier.
	SecretID *string `json:"secret_id" validate:"required"`
}

// Constants associated with the UsernamePasswordSecretVersionMetadata.SecretType property.
// The secret type. Supported types are arbitrary, imported_cert, public_cert, private_cert, iam_credentials,
// service_credentials, kv, and username_password.
const (
	UsernamePasswordSecretVersionMetadata_SecretType_Arbitrary          = "arbitrary"
	UsernamePasswordSecretVersionMetadata_SecretType_IamCredentials     = "iam_credentials"
	UsernamePasswordSecretVersionMetadata_SecretType_ImportedCert       = "imported_cert"
	UsernamePasswordSecretVersionMetadata_SecretType_Kv                 = "kv"
	UsernamePasswordSecretVersionMetadata_SecretType_PrivateCert        = "private_cert"
	UsernamePasswordSecretVersionMetadata_SecretType_PublicCert         = "public_cert"
	UsernamePasswordSecretVersionMetadata_SecretType_ServiceCredentials = "service_credentials"
	UsernamePasswordSecretVersionMetadata_SecretType_UsernamePassword   = "username_password"
)

// Constants associated with the UsernamePasswordSecretVersionMetadata.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	UsernamePasswordSecretVersionMetadata_Alias_Current  = "current"
	UsernamePasswordSecretVersionMetadata_Alias_Previous = "previous"
)

func (*UsernamePasswordSecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalUsernamePasswordSecretVersionMetadata unmarshals an instance of UsernamePasswordSecretVersionMetadata from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
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
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsernamePasswordSecretVersionPrototype : UsernamePasswordSecretVersionPrototype struct
// This model "extends" SecretVersionPrototype
type UsernamePasswordSecretVersionPrototype struct {
	// The password that is assigned to an `username_password` secret.
	Password *string `json:"password,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

func (*UsernamePasswordSecretVersionPrototype) isaSecretVersionPrototype() bool {
	return true
}

// UnmarshalUsernamePasswordSecretVersionPrototype unmarshals an instance of UsernamePasswordSecretVersionPrototype from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretVersionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretVersionPrototype)
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretsPager can be used to simplify the use of the "ListSecrets" method.
type SecretsPager struct {
	hasNext     bool
	options     *ListSecretsOptions
	client      *SecretsManagerV2
	pageContext struct {
		next *int64
	}
}

// NewSecretsPager returns a new SecretsPager instance.
func (secretsManager *SecretsManagerV2) NewSecretsPager(options *ListSecretsOptions) (pager *SecretsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListSecretsOptions = *options
	pager = &SecretsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  secretsManager,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SecretsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SecretsPager) GetNextWithContext(ctx context.Context) (page []SecretMetadataIntf, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListSecretsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Secrets

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SecretsPager) GetAllWithContext(ctx context.Context) (allItems []SecretMetadataIntf, err error) {
	for pager.HasNext() {
		var nextPage []SecretMetadataIntf
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *SecretsPager) GetNext() (page []SecretMetadataIntf, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SecretsPager) GetAll() (allItems []SecretMetadataIntf, err error) {
	return pager.GetAllWithContext(context.Background())
}

// SecretsLocksPager can be used to simplify the use of the "ListSecretsLocks" method.
type SecretsLocksPager struct {
	hasNext     bool
	options     *ListSecretsLocksOptions
	client      *SecretsManagerV2
	pageContext struct {
		next *int64
	}
}

// NewSecretsLocksPager returns a new SecretsLocksPager instance.
func (secretsManager *SecretsManagerV2) NewSecretsLocksPager(options *ListSecretsLocksOptions) (pager *SecretsLocksPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListSecretsLocksOptions = *options
	pager = &SecretsLocksPager{
		hasNext: true,
		options: &optionsCopy,
		client:  secretsManager,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SecretsLocksPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SecretsLocksPager) GetNextWithContext(ctx context.Context) (page []SecretLocks, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListSecretsLocksWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.SecretsLocks

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SecretsLocksPager) GetAllWithContext(ctx context.Context) (allItems []SecretLocks, err error) {
	for pager.HasNext() {
		var nextPage []SecretLocks
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *SecretsLocksPager) GetNext() (page []SecretLocks, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SecretsLocksPager) GetAll() (allItems []SecretLocks, err error) {
	return pager.GetAllWithContext(context.Background())
}

// SecretLocksPager can be used to simplify the use of the "ListSecretLocks" method.
type SecretLocksPager struct {
	hasNext     bool
	options     *ListSecretLocksOptions
	client      *SecretsManagerV2
	pageContext struct {
		next *int64
	}
}

// NewSecretLocksPager returns a new SecretLocksPager instance.
func (secretsManager *SecretsManagerV2) NewSecretLocksPager(options *ListSecretLocksOptions) (pager *SecretLocksPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListSecretLocksOptions = *options
	pager = &SecretLocksPager{
		hasNext: true,
		options: &optionsCopy,
		client:  secretsManager,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SecretLocksPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SecretLocksPager) GetNextWithContext(ctx context.Context) (page []SecretLock, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListSecretLocksWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Locks

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SecretLocksPager) GetAllWithContext(ctx context.Context) (allItems []SecretLock, err error) {
	for pager.HasNext() {
		var nextPage []SecretLock
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *SecretLocksPager) GetNext() (page []SecretLock, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SecretLocksPager) GetAll() (allItems []SecretLock, err error) {
	return pager.GetAllWithContext(context.Background())
}

// SecretVersionLocksPager can be used to simplify the use of the "ListSecretVersionLocks" method.
type SecretVersionLocksPager struct {
	hasNext     bool
	options     *ListSecretVersionLocksOptions
	client      *SecretsManagerV2
	pageContext struct {
		next *int64
	}
}

// NewSecretVersionLocksPager returns a new SecretVersionLocksPager instance.
func (secretsManager *SecretsManagerV2) NewSecretVersionLocksPager(options *ListSecretVersionLocksOptions) (pager *SecretVersionLocksPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListSecretVersionLocksOptions = *options
	pager = &SecretVersionLocksPager{
		hasNext: true,
		options: &optionsCopy,
		client:  secretsManager,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SecretVersionLocksPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SecretVersionLocksPager) GetNextWithContext(ctx context.Context) (page []SecretLock, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListSecretVersionLocksWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Locks

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SecretVersionLocksPager) GetAllWithContext(ctx context.Context) (allItems []SecretLock, err error) {
	for pager.HasNext() {
		var nextPage []SecretLock
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *SecretVersionLocksPager) GetNext() (page []SecretLock, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SecretVersionLocksPager) GetAll() (allItems []SecretLock, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ConfigurationsPager can be used to simplify the use of the "ListConfigurations" method.
type ConfigurationsPager struct {
	hasNext     bool
	options     *ListConfigurationsOptions
	client      *SecretsManagerV2
	pageContext struct {
		next *int64
	}
}

// NewConfigurationsPager returns a new ConfigurationsPager instance.
func (secretsManager *SecretsManagerV2) NewConfigurationsPager(options *ListConfigurationsOptions) (pager *ConfigurationsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListConfigurationsOptions = *options
	pager = &ConfigurationsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  secretsManager,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ConfigurationsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ConfigurationsPager) GetNextWithContext(ctx context.Context) (page []ConfigurationMetadataIntf, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListConfigurationsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Configurations

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ConfigurationsPager) GetAllWithContext(ctx context.Context) (allItems []ConfigurationMetadataIntf, err error) {
	for pager.HasNext() {
		var nextPage []ConfigurationMetadataIntf
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ConfigurationsPager) GetNext() (page []ConfigurationMetadataIntf, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ConfigurationsPager) GetAll() (allItems []ConfigurationMetadataIntf, err error) {
	return pager.GetAllWithContext(context.Background())
}
