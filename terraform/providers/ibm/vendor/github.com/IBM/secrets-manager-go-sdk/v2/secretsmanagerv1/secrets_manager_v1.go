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
 * IBM OpenAPI SDK Code Generator Version: 3.60.2-95dc7721-20221102-203229
 */

// Package secretsmanagerv1 : Operations and models for the SecretsManagerV1 service
package secretsmanagerv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/secrets-manager-go-sdk/v2/common"
	"github.com/go-openapi/strfmt"
)

// SecretsManagerV1 : With IBM Cloud® Secrets Manager, you can create, lease, and centrally manage secrets that are used
// in IBM Cloud services or your custom-built applications. Secrets are stored in a dedicated instance of Secrets
// Manager, which is built on open source HashiCorp Vault.
//
// API Version: 1.0.0
// See: https://cloud.ibm.com/docs/secrets-manager
type SecretsManagerV1 struct {
	Service *core.BaseService
}

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "secrets_manager"

// SecretsManagerV1Options : Service options
type SecretsManagerV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewSecretsManagerV1UsingExternalConfig : constructs an instance of SecretsManagerV1 with passed in options and external configuration.
func NewSecretsManagerV1UsingExternalConfig(options *SecretsManagerV1Options) (secretsManager *SecretsManagerV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	secretsManager, err = NewSecretsManagerV1(options)
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

// NewSecretsManagerV1 : constructs an instance of SecretsManagerV1 with passed in options.
func NewSecretsManagerV1(options *SecretsManagerV1Options) (service *SecretsManagerV1, err error) {
	serviceOptions := &core.ServiceOptions{
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

	service = &SecretsManagerV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "secretsManager" suitable for processing requests.
func (secretsManager *SecretsManagerV1) Clone() *SecretsManagerV1 {
	if core.IsNil(secretsManager) {
		return nil
	}
	clone := *secretsManager
	clone.Service = secretsManager.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (secretsManager *SecretsManagerV1) SetServiceURL(url string) error {
	return secretsManager.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (secretsManager *SecretsManagerV1) GetServiceURL() string {
	return secretsManager.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (secretsManager *SecretsManagerV1) SetDefaultHeaders(headers http.Header) {
	secretsManager.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (secretsManager *SecretsManagerV1) SetEnableGzipCompression(enableGzip bool) {
	secretsManager.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (secretsManager *SecretsManagerV1) GetEnableGzipCompression() bool {
	return secretsManager.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (secretsManager *SecretsManagerV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	secretsManager.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (secretsManager *SecretsManagerV1) DisableRetries() {
	secretsManager.Service.DisableRetries()
}

// CreateSecretGroup : Create a secret group
// Create a secret group that you can use to organize secrets and control who on your team has access to them.
//
// A successful request returns the ID value of the secret group, along with other metadata. To learn more about secret
// groups, check out the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-secret-groups).
func (secretsManager *SecretsManagerV1) CreateSecretGroup(createSecretGroupOptions *CreateSecretGroupOptions) (result *SecretGroupDef, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretGroupWithContext(context.Background(), createSecretGroupOptions)
}

// CreateSecretGroupWithContext is an alternate form of the CreateSecretGroup method which supports a Context parameter
func (secretsManager *SecretsManagerV1) CreateSecretGroupWithContext(ctx context.Context, createSecretGroupOptions *CreateSecretGroupOptions) (result *SecretGroupDef, response *core.DetailedResponse, err error) {
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
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secret_groups`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "CreateSecretGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createSecretGroupOptions.Metadata != nil {
		body["metadata"] = createSecretGroupOptions.Metadata
	}
	if createSecretGroupOptions.Resources != nil {
		body["resources"] = createSecretGroupOptions.Resources
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupDef)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSecretGroups : List secret groups
// List the secret groups that are available in your Secrets Manager instance.
func (secretsManager *SecretsManagerV1) ListSecretGroups(listSecretGroupsOptions *ListSecretGroupsOptions) (result *SecretGroupDef, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretGroupsWithContext(context.Background(), listSecretGroupsOptions)
}

// ListSecretGroupsWithContext is an alternate form of the ListSecretGroups method which supports a Context parameter
func (secretsManager *SecretsManagerV1) ListSecretGroupsWithContext(ctx context.Context, listSecretGroupsOptions *ListSecretGroupsOptions) (result *SecretGroupDef, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSecretGroupsOptions, "listSecretGroupsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secret_groups`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretGroupsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "ListSecretGroups")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupDef)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecretGroup : Get a secret group
// Get the metadata of an existing secret group by specifying the ID of the group.
func (secretsManager *SecretsManagerV1) GetSecretGroup(getSecretGroupOptions *GetSecretGroupOptions) (result *SecretGroupDef, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretGroupWithContext(context.Background(), getSecretGroupOptions)
}

// GetSecretGroupWithContext is an alternate form of the GetSecretGroup method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetSecretGroupWithContext(ctx context.Context, getSecretGroupOptions *GetSecretGroupOptions) (result *SecretGroupDef, response *core.DetailedResponse, err error) {
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
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secret_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetSecretGroup")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupDef)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecretGroupMetadata : Update a secret group
// Update the metadata of an existing secret group, such as its name or description.
func (secretsManager *SecretsManagerV1) UpdateSecretGroupMetadata(updateSecretGroupMetadataOptions *UpdateSecretGroupMetadataOptions) (result *SecretGroupDef, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateSecretGroupMetadataWithContext(context.Background(), updateSecretGroupMetadataOptions)
}

// UpdateSecretGroupMetadataWithContext is an alternate form of the UpdateSecretGroupMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV1) UpdateSecretGroupMetadataWithContext(ctx context.Context, updateSecretGroupMetadataOptions *UpdateSecretGroupMetadataOptions) (result *SecretGroupDef, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecretGroupMetadataOptions, "updateSecretGroupMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSecretGroupMetadataOptions, "updateSecretGroupMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateSecretGroupMetadataOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secret_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSecretGroupMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "UpdateSecretGroupMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSecretGroupMetadataOptions.Metadata != nil {
		body["metadata"] = updateSecretGroupMetadataOptions.Metadata
	}
	if updateSecretGroupMetadataOptions.Resources != nil {
		body["resources"] = updateSecretGroupMetadataOptions.Resources
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupDef)
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
func (secretsManager *SecretsManagerV1) DeleteSecretGroup(deleteSecretGroupOptions *DeleteSecretGroupOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteSecretGroupWithContext(context.Background(), deleteSecretGroupOptions)
}

// DeleteSecretGroupWithContext is an alternate form of the DeleteSecretGroup method which supports a Context parameter
func (secretsManager *SecretsManagerV1) DeleteSecretGroupWithContext(ctx context.Context, deleteSecretGroupOptions *DeleteSecretGroupOptions) (response *core.DetailedResponse, err error) {
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
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secret_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSecretGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "DeleteSecretGroup")
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

// CreateSecret : Create a secret
// Create a secret or import an existing value that you can use to access or authenticate to a protected resource.
//
// Use this method to either generate or import an existing secret, such as an arbitrary value or a TLS certificate,
// that you can manage in your Secrets Manager service instance. A successful request stores the secret in your
// dedicated instance based on the secret type and data that you specify. The response returns the ID value of the
// secret, along with other metadata.
//
// To learn more about the types of secrets that you can create with Secrets Manager, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-what-is-secret).
func (secretsManager *SecretsManagerV1) CreateSecret(createSecretOptions *CreateSecretOptions) (result *CreateSecret, response *core.DetailedResponse, err error) {
	return secretsManager.CreateSecretWithContext(context.Background(), createSecretOptions)
}

// CreateSecretWithContext is an alternate form of the CreateSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV1) CreateSecretWithContext(ctx context.Context, createSecretOptions *CreateSecretOptions) (result *CreateSecret, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretOptions, "createSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSecretOptions, "createSecretOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *createSecretOptions.SecretType,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "CreateSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createSecretOptions.Metadata != nil {
		body["metadata"] = createSecretOptions.Metadata
	}
	if createSecretOptions.Resources != nil {
		body["resources"] = createSecretOptions.Resources
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateSecret)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSecrets : List secrets by type
// List the secrets in your Secrets Manager instance based on the type that you specify.
func (secretsManager *SecretsManagerV1) ListSecrets(listSecretsOptions *ListSecretsOptions) (result *ListSecrets, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretsWithContext(context.Background(), listSecretsOptions)
}

// ListSecretsWithContext is an alternate form of the ListSecrets method which supports a Context parameter
func (secretsManager *SecretsManagerV1) ListSecretsWithContext(ctx context.Context, listSecretsOptions *ListSecretsOptions) (result *ListSecrets, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSecretsOptions, "listSecretsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listSecretsOptions, "listSecretsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *listSecretsOptions.SecretType,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "ListSecrets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSecretsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSecretsOptions.Limit))
	}
	if listSecretsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listSecretsOptions.Offset))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListSecrets)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListAllSecrets : List all secrets
// List all of the secrets in your Secrets Manager instance.
func (secretsManager *SecretsManagerV1) ListAllSecrets(listAllSecretsOptions *ListAllSecretsOptions) (result *ListSecrets, response *core.DetailedResponse, err error) {
	return secretsManager.ListAllSecretsWithContext(context.Background(), listAllSecretsOptions)
}

// ListAllSecretsWithContext is an alternate form of the ListAllSecrets method which supports a Context parameter
func (secretsManager *SecretsManagerV1) ListAllSecretsWithContext(ctx context.Context, listAllSecretsOptions *ListAllSecretsOptions) (result *ListSecrets, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllSecretsOptions, "listAllSecretsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllSecretsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "ListAllSecrets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAllSecretsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAllSecretsOptions.Limit))
	}
	if listAllSecretsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listAllSecretsOptions.Offset))
	}
	if listAllSecretsOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listAllSecretsOptions.Search))
	}
	if listAllSecretsOptions.SortBy != nil {
		builder.AddQuery("sort_by", fmt.Sprint(*listAllSecretsOptions.SortBy))
	}
	if listAllSecretsOptions.Groups != nil {
		builder.AddQuery("groups", strings.Join(listAllSecretsOptions.Groups, ","))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListSecrets)
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
// metadata](#get-secret-metadata) method.
func (secretsManager *SecretsManagerV1) GetSecret(getSecretOptions *GetSecretOptions) (result *GetSecret, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretWithContext(context.Background(), getSecretOptions)
}

// GetSecretWithContext is an alternate form of the GetSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetSecretWithContext(ctx context.Context, getSecretOptions *GetSecretOptions) (result *GetSecret, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretOptions, "getSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretOptions, "getSecretOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *getSecretOptions.SecretType,
		"id":          *getSecretOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetSecret")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecret)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecret : Invoke an action on a secret
// Invoke an action on a specified secret. This method supports the following actions:
//
// - `rotate`: Replace the value of a secret.
// - `restore`: Restore a previous version of an `iam_credentials` secret.
// - `revoke`: Revoke a private certificate.
// - `delete_credentials`: Delete the API key that is associated with an `iam_credentials` secret.
// - `validate_dns_challenge`: Validate challenges for a public certificate that is ordered with a manual DNS provider.
func (secretsManager *SecretsManagerV1) UpdateSecret(updateSecretOptions *UpdateSecretOptions) (result *GetSecret, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateSecretWithContext(context.Background(), updateSecretOptions)
}

// UpdateSecretWithContext is an alternate form of the UpdateSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV1) UpdateSecretWithContext(ctx context.Context, updateSecretOptions *UpdateSecretOptions) (result *GetSecret, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecretOptions, "updateSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSecretOptions, "updateSecretOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *updateSecretOptions.SecretType,
		"id":          *updateSecretOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "UpdateSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	builder.AddQuery("action", fmt.Sprint(*updateSecretOptions.Action))

	if updateSecretOptions.SecretAction != nil {
		_, err = builder.SetBodyContentJSON(updateSecretOptions.SecretAction)
		if err != nil {
			return
		}
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecret)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteSecret : Delete a secret
// Delete a secret by specifying the ID of the secret.
func (secretsManager *SecretsManagerV1) DeleteSecret(deleteSecretOptions *DeleteSecretOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteSecretWithContext(context.Background(), deleteSecretOptions)
}

// DeleteSecretWithContext is an alternate form of the DeleteSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV1) DeleteSecretWithContext(ctx context.Context, deleteSecretOptions *DeleteSecretOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSecretOptions, "deleteSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSecretOptions, "deleteSecretOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *deleteSecretOptions.SecretType,
		"id":          *deleteSecretOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "DeleteSecret")
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

// ListSecretVersions : List versions of a secret
// List the versions of a secret.
//
// A successful request returns the list of the versions along with the metadata of each version.
func (secretsManager *SecretsManagerV1) ListSecretVersions(listSecretVersionsOptions *ListSecretVersionsOptions) (result *ListSecretVersions, response *core.DetailedResponse, err error) {
	return secretsManager.ListSecretVersionsWithContext(context.Background(), listSecretVersionsOptions)
}

// ListSecretVersionsWithContext is an alternate form of the ListSecretVersions method which supports a Context parameter
func (secretsManager *SecretsManagerV1) ListSecretVersionsWithContext(ctx context.Context, listSecretVersionsOptions *ListSecretVersionsOptions) (result *ListSecretVersions, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSecretVersionsOptions, "listSecretVersionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listSecretVersionsOptions, "listSecretVersionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *listSecretVersionsOptions.SecretType,
		"id":          *listSecretVersionsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/versions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSecretVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "ListSecretVersions")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListSecretVersions)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecretVersion : Get a version of a secret
// Get a version of a secret by specifying the ID of the version or the alias `previous`.
//
// A successful request returns the secret data that is associated with the specified version of your secret, along with
// other metadata.
func (secretsManager *SecretsManagerV1) GetSecretVersion(getSecretVersionOptions *GetSecretVersionOptions) (result *GetSecretVersion, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretVersionWithContext(context.Background(), getSecretVersionOptions)
}

// GetSecretVersionWithContext is an alternate form of the GetSecretVersion method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetSecretVersionWithContext(ctx context.Context, getSecretVersionOptions *GetSecretVersionOptions) (result *GetSecretVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretVersionOptions, "getSecretVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretVersionOptions, "getSecretVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *getSecretVersionOptions.SecretType,
		"id":          *getSecretVersionOptions.ID,
		"version_id":  *getSecretVersionOptions.VersionID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/versions/{version_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetSecretVersion")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecretVersion : Invoke an action on a version of a secret
// Invoke an action on a specified version of a secret. This method supports the following actions:
//
// - `revoke`: Revoke a version of a private certificate.
func (secretsManager *SecretsManagerV1) UpdateSecretVersion(updateSecretVersionOptions *UpdateSecretVersionOptions) (result *GetSecret, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateSecretVersionWithContext(context.Background(), updateSecretVersionOptions)
}

// UpdateSecretVersionWithContext is an alternate form of the UpdateSecretVersion method which supports a Context parameter
func (secretsManager *SecretsManagerV1) UpdateSecretVersionWithContext(ctx context.Context, updateSecretVersionOptions *UpdateSecretVersionOptions) (result *GetSecret, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecretVersionOptions, "updateSecretVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSecretVersionOptions, "updateSecretVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *updateSecretVersionOptions.SecretType,
		"id":          *updateSecretVersionOptions.ID,
		"version_id":  *updateSecretVersionOptions.VersionID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/versions/{version_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSecretVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "UpdateSecretVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("action", fmt.Sprint(*updateSecretVersionOptions.Action))

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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecret)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecretVersionMetadata : Get the metadata of a secret version
// Get the metadata of a secret version by specifying the ID of the version or the alias `previous`.
//
// A successful request returns the metadata that is associated with the specified version of your secret.
func (secretsManager *SecretsManagerV1) GetSecretVersionMetadata(getSecretVersionMetadataOptions *GetSecretVersionMetadataOptions) (result *GetSecretVersionMetadata, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretVersionMetadataWithContext(context.Background(), getSecretVersionMetadataOptions)
}

// GetSecretVersionMetadataWithContext is an alternate form of the GetSecretVersionMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetSecretVersionMetadataWithContext(ctx context.Context, getSecretVersionMetadataOptions *GetSecretVersionMetadataOptions) (result *GetSecretVersionMetadata, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretVersionMetadataOptions, "getSecretVersionMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretVersionMetadataOptions, "getSecretVersionMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *getSecretVersionMetadataOptions.SecretType,
		"id":          *getSecretVersionMetadataOptions.ID,
		"version_id":  *getSecretVersionMetadataOptions.VersionID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/versions/{version_id}/metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretVersionMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetSecretVersionMetadata")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretVersionMetadata)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecretVersionMetadata : Update the metadata of a secret version
// Update the metadata of a secret version, such as `version_custom_metadata`.
func (secretsManager *SecretsManagerV1) UpdateSecretVersionMetadata(updateSecretVersionMetadataOptions *UpdateSecretVersionMetadataOptions) (result *GetSecretVersionMetadata, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateSecretVersionMetadataWithContext(context.Background(), updateSecretVersionMetadataOptions)
}

// UpdateSecretVersionMetadataWithContext is an alternate form of the UpdateSecretVersionMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV1) UpdateSecretVersionMetadataWithContext(ctx context.Context, updateSecretVersionMetadataOptions *UpdateSecretVersionMetadataOptions) (result *GetSecretVersionMetadata, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecretVersionMetadataOptions, "updateSecretVersionMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSecretVersionMetadataOptions, "updateSecretVersionMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *updateSecretVersionMetadataOptions.SecretType,
		"id":          *updateSecretVersionMetadataOptions.ID,
		"version_id":  *updateSecretVersionMetadataOptions.VersionID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/versions/{version_id}/metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSecretVersionMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "UpdateSecretVersionMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSecretVersionMetadataOptions.Metadata != nil {
		body["metadata"] = updateSecretVersionMetadataOptions.Metadata
	}
	if updateSecretVersionMetadataOptions.Resources != nil {
		body["resources"] = updateSecretVersionMetadataOptions.Resources
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretVersionMetadata)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecretMetadata : Get the metadata of a secret
// Get the details of a secret by specifying its ID.
//
// A successful request returns only metadata about the secret, such as its name and creation date. To retrieve the
// value of a secret, use the [Get a secret](#get-secret) or [Get a version of a secret](#get-secret-version) methods.
func (secretsManager *SecretsManagerV1) GetSecretMetadata(getSecretMetadataOptions *GetSecretMetadataOptions) (result *SecretMetadataRequest, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretMetadataWithContext(context.Background(), getSecretMetadataOptions)
}

// GetSecretMetadataWithContext is an alternate form of the GetSecretMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetSecretMetadataWithContext(ctx context.Context, getSecretMetadataOptions *GetSecretMetadataOptions) (result *SecretMetadataRequest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretMetadataOptions, "getSecretMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretMetadataOptions, "getSecretMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *getSecretMetadataOptions.SecretType,
		"id":          *getSecretMetadataOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetSecretMetadata")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretMetadataRequest)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSecretMetadata : Update secret metadata
// Update the metadata of a secret, such as its name or description.
//
// To update the actual contents of a secret, rotate the secret by using the [Invoke an action on a
// secret](#update-secret) method.
func (secretsManager *SecretsManagerV1) UpdateSecretMetadata(updateSecretMetadataOptions *UpdateSecretMetadataOptions) (result *SecretMetadataRequest, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateSecretMetadataWithContext(context.Background(), updateSecretMetadataOptions)
}

// UpdateSecretMetadataWithContext is an alternate form of the UpdateSecretMetadata method which supports a Context parameter
func (secretsManager *SecretsManagerV1) UpdateSecretMetadataWithContext(ctx context.Context, updateSecretMetadataOptions *UpdateSecretMetadataOptions) (result *SecretMetadataRequest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSecretMetadataOptions, "updateSecretMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSecretMetadataOptions, "updateSecretMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *updateSecretMetadataOptions.SecretType,
		"id":          *updateSecretMetadataOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSecretMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "UpdateSecretMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSecretMetadataOptions.Metadata != nil {
		body["metadata"] = updateSecretMetadataOptions.Metadata
	}
	if updateSecretMetadataOptions.Resources != nil {
		body["resources"] = updateSecretMetadataOptions.Resources
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretMetadataRequest)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLocks : List secret locks
// List the locks that are associated with a specified secret.
func (secretsManager *SecretsManagerV1) GetLocks(getLocksOptions *GetLocksOptions) (result *ListSecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.GetLocksWithContext(context.Background(), getLocksOptions)
}

// GetLocksWithContext is an alternate form of the GetLocks method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetLocksWithContext(ctx context.Context, getLocksOptions *GetLocksOptions) (result *ListSecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLocksOptions, "getLocksOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLocksOptions, "getLocksOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *getLocksOptions.SecretType,
		"id":          *getLocksOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/locks/{secret_type}/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLocksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetLocks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getLocksOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getLocksOptions.Limit))
	}
	if getLocksOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*getLocksOptions.Offset))
	}
	if getLocksOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*getLocksOptions.Search))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// LockSecret : Lock a secret
// Create a lock on the current version of a secret.
//
// A lock can be used to prevent a secret from being deleted or modified while it's in use by your applications. A
// successful request attaches a new lock to your secret, or replaces a lock of the same name if it already exists.
// Additionally, you can use this method to clear any matching locks on a secret by using one of the following optional
// lock modes:
//
// - `exclusive`: Removes any other locks with matching names if they are found in the previous version of the secret.
// - `exclusive_delete`: Same as `exclusive`, but also permanently deletes the data of the previous secret version if it
// doesn't have any locks.
//
// For more information about locking secrets, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-secret-locks).
func (secretsManager *SecretsManagerV1) LockSecret(lockSecretOptions *LockSecretOptions) (result *GetSecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.LockSecretWithContext(context.Background(), lockSecretOptions)
}

// LockSecretWithContext is an alternate form of the LockSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV1) LockSecretWithContext(ctx context.Context, lockSecretOptions *LockSecretOptions) (result *GetSecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(lockSecretOptions, "lockSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(lockSecretOptions, "lockSecretOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *lockSecretOptions.SecretType,
		"id":          *lockSecretOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/locks/{secret_type}/{id}/lock`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range lockSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "LockSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if lockSecretOptions.Mode != nil {
		builder.AddQuery("mode", fmt.Sprint(*lockSecretOptions.Mode))
	}

	body := make(map[string]interface{})
	if lockSecretOptions.Locks != nil {
		body["locks"] = lockSecretOptions.Locks
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UnlockSecret : Unlock a secret
// Delete one or more locks that are associated with the current version of a secret.
//
// A successful request deletes the locks that you specify. To remove all locks, you can pass `{"locks": ["*"]}` in in
// the request body. Otherwise, specify the names of the locks that you want to delete. For example, `{"locks":
// ["lock1", "lock2"]}`.
//
// **Note:** A secret is considered unlocked and able to be revoked or deleted only after all of its locks are removed.
// To understand whether a secret contains locks, check the `locks_total` field that is returned as part of the metadata
// of your secret.
func (secretsManager *SecretsManagerV1) UnlockSecret(unlockSecretOptions *UnlockSecretOptions) (result *GetSecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.UnlockSecretWithContext(context.Background(), unlockSecretOptions)
}

// UnlockSecretWithContext is an alternate form of the UnlockSecret method which supports a Context parameter
func (secretsManager *SecretsManagerV1) UnlockSecretWithContext(ctx context.Context, unlockSecretOptions *UnlockSecretOptions) (result *GetSecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(unlockSecretOptions, "unlockSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(unlockSecretOptions, "unlockSecretOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *unlockSecretOptions.SecretType,
		"id":          *unlockSecretOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/locks/{secret_type}/{id}/unlock`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range unlockSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "UnlockSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if unlockSecretOptions.Locks != nil {
		body["locks"] = unlockSecretOptions.Locks
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecretVersionLocks : List secret version locks
// List the locks that are associated with a specified secret version.
func (secretsManager *SecretsManagerV1) GetSecretVersionLocks(getSecretVersionLocksOptions *GetSecretVersionLocksOptions) (result *ListSecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.GetSecretVersionLocksWithContext(context.Background(), getSecretVersionLocksOptions)
}

// GetSecretVersionLocksWithContext is an alternate form of the GetSecretVersionLocks method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetSecretVersionLocksWithContext(ctx context.Context, getSecretVersionLocksOptions *GetSecretVersionLocksOptions) (result *ListSecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretVersionLocksOptions, "getSecretVersionLocksOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecretVersionLocksOptions, "getSecretVersionLocksOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *getSecretVersionLocksOptions.SecretType,
		"id":          *getSecretVersionLocksOptions.ID,
		"version_id":  *getSecretVersionLocksOptions.VersionID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/locks/{secret_type}/{id}/versions/{version_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecretVersionLocksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetSecretVersionLocks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getSecretVersionLocksOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getSecretVersionLocksOptions.Limit))
	}
	if getSecretVersionLocksOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*getSecretVersionLocksOptions.Offset))
	}
	if getSecretVersionLocksOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*getSecretVersionLocksOptions.Search))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// LockSecretVersion : Lock a secret version
// Create a lock on the specified version of a secret.
//
// A lock can be used to prevent a secret from being deleted or modified while it's in use by your applications. A
// successful request attaches a new lock to the specified version, or replaces a lock of the same name if it already
// exists. Additionally, you can use this method to clear any matching locks on a secret version by using one of the
// following optional lock modes:
//
// - `exclusive`: Removes any other locks with matching names if they are found in the previous version of the secret.
// - `exclusive_delete`: Same as `exclusive`, but also permanently deletes the data of the previous secret version if it
// doesn't have any locks.
//
// For more information about locking secrets, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-secret-locks).
func (secretsManager *SecretsManagerV1) LockSecretVersion(lockSecretVersionOptions *LockSecretVersionOptions) (result *GetSecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.LockSecretVersionWithContext(context.Background(), lockSecretVersionOptions)
}

// LockSecretVersionWithContext is an alternate form of the LockSecretVersion method which supports a Context parameter
func (secretsManager *SecretsManagerV1) LockSecretVersionWithContext(ctx context.Context, lockSecretVersionOptions *LockSecretVersionOptions) (result *GetSecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(lockSecretVersionOptions, "lockSecretVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(lockSecretVersionOptions, "lockSecretVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *lockSecretVersionOptions.SecretType,
		"id":          *lockSecretVersionOptions.ID,
		"version_id":  *lockSecretVersionOptions.VersionID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/locks/{secret_type}/{id}/versions/{version_id}/lock`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range lockSecretVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "LockSecretVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if lockSecretVersionOptions.Mode != nil {
		builder.AddQuery("mode", fmt.Sprint(*lockSecretVersionOptions.Mode))
	}

	body := make(map[string]interface{})
	if lockSecretVersionOptions.Locks != nil {
		body["locks"] = lockSecretVersionOptions.Locks
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UnlockSecretVersion : Unlock a secret version
// Delete one or more locks that are associated with the specified secret version.
//
// A successful request deletes the locks that you specify. To remove all locks, you can pass `{"locks": ["*"]}` in in
// the request body. Otherwise, specify the names of the locks that you want to delete. For example, `{"locks":
// ["lock-1", "lock-2"]}`.
//
// **Note:** A secret is considered unlocked and able to be revoked or deleted only after all of its locks are removed.
// To understand whether a secret contains locks, check the `locks_total` field that is returned as part of the metadata
// of your secret.
func (secretsManager *SecretsManagerV1) UnlockSecretVersion(unlockSecretVersionOptions *UnlockSecretVersionOptions) (result *GetSecretLocks, response *core.DetailedResponse, err error) {
	return secretsManager.UnlockSecretVersionWithContext(context.Background(), unlockSecretVersionOptions)
}

// UnlockSecretVersionWithContext is an alternate form of the UnlockSecretVersion method which supports a Context parameter
func (secretsManager *SecretsManagerV1) UnlockSecretVersionWithContext(ctx context.Context, unlockSecretVersionOptions *UnlockSecretVersionOptions) (result *GetSecretLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(unlockSecretVersionOptions, "unlockSecretVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(unlockSecretVersionOptions, "unlockSecretVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *unlockSecretVersionOptions.SecretType,
		"id":          *unlockSecretVersionOptions.ID,
		"version_id":  *unlockSecretVersionOptions.VersionID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/locks/{secret_type}/{id}/versions/{version_id}/unlock`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range unlockSecretVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "UnlockSecretVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if unlockSecretVersionOptions.Locks != nil {
		body["locks"] = unlockSecretVersionOptions.Locks
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListInstanceSecretsLocks : List all secrets and locks
// List the lock details that are associated with all secrets in your Secrets Manager instance.
func (secretsManager *SecretsManagerV1) ListInstanceSecretsLocks(listInstanceSecretsLocksOptions *ListInstanceSecretsLocksOptions) (result *GetInstanceLocks, response *core.DetailedResponse, err error) {
	return secretsManager.ListInstanceSecretsLocksWithContext(context.Background(), listInstanceSecretsLocksOptions)
}

// ListInstanceSecretsLocksWithContext is an alternate form of the ListInstanceSecretsLocks method which supports a Context parameter
func (secretsManager *SecretsManagerV1) ListInstanceSecretsLocksWithContext(ctx context.Context, listInstanceSecretsLocksOptions *ListInstanceSecretsLocksOptions) (result *GetInstanceLocks, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listInstanceSecretsLocksOptions, "listInstanceSecretsLocksOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/locks`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listInstanceSecretsLocksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "ListInstanceSecretsLocks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listInstanceSecretsLocksOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listInstanceSecretsLocksOptions.Limit))
	}
	if listInstanceSecretsLocksOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listInstanceSecretsLocksOptions.Offset))
	}
	if listInstanceSecretsLocksOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listInstanceSecretsLocksOptions.Search))
	}
	if listInstanceSecretsLocksOptions.Groups != nil {
		builder.AddQuery("groups", strings.Join(listInstanceSecretsLocksOptions.Groups, ","))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetInstanceLocks)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutPolicy : Set secret policies
// Create or update one or more policies, such as an [automatic rotation
// policy](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-automatic-rotation), for the specified
// secret. To remove a policy, keep the resources block empty.
func (secretsManager *SecretsManagerV1) PutPolicy(putPolicyOptions *PutPolicyOptions) (result GetSecretPoliciesIntf, response *core.DetailedResponse, err error) {
	return secretsManager.PutPolicyWithContext(context.Background(), putPolicyOptions)
}

// PutPolicyWithContext is an alternate form of the PutPolicy method which supports a Context parameter
func (secretsManager *SecretsManagerV1) PutPolicyWithContext(ctx context.Context, putPolicyOptions *PutPolicyOptions) (result GetSecretPoliciesIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putPolicyOptions, "putPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putPolicyOptions, "putPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *putPolicyOptions.SecretType,
		"id":          *putPolicyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/policies`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "PutPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if putPolicyOptions.Policy != nil {
		builder.AddQuery("policy", fmt.Sprint(*putPolicyOptions.Policy))
	}

	body := make(map[string]interface{})
	if putPolicyOptions.Metadata != nil {
		body["metadata"] = putPolicyOptions.Metadata
	}
	if putPolicyOptions.Resources != nil {
		body["resources"] = putPolicyOptions.Resources
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretPolicies)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetPolicy : List secret policies
// List the rotation policies that are associated with a specified secret.
func (secretsManager *SecretsManagerV1) GetPolicy(getPolicyOptions *GetPolicyOptions) (result GetSecretPoliciesIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetPolicyWithContext(context.Background(), getPolicyOptions)
}

// GetPolicyWithContext is an alternate form of the GetPolicy method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetPolicyWithContext(ctx context.Context, getPolicyOptions *GetPolicyOptions) (result GetSecretPoliciesIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPolicyOptions, "getPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getPolicyOptions, "getPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *getPolicyOptions.SecretType,
		"id":          *getPolicyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/secrets/{secret_type}/{id}/policies`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getPolicyOptions.Policy != nil {
		builder.AddQuery("policy", fmt.Sprint(*getPolicyOptions.Policy))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretPolicies)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutConfig : Set the configuration of a secret type
// Set the configuration for the specified secret type.
//
// Use this method to configure the IAM credentials (`iam_credentials`) engine for your service instance. Looking to
// order or generate certificates? To configure the public certificates (`public_cert`) or  private certificates
// (`private_cert`) engines, use the [Add a configuration](#create_config_element) method.
func (secretsManager *SecretsManagerV1) PutConfig(putConfigOptions *PutConfigOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.PutConfigWithContext(context.Background(), putConfigOptions)
}

// PutConfigWithContext is an alternate form of the PutConfig method which supports a Context parameter
func (secretsManager *SecretsManagerV1) PutConfigWithContext(ctx context.Context, putConfigOptions *PutConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putConfigOptions, "putConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putConfigOptions, "putConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *putConfigOptions.SecretType,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/config/{secret_type}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "PutConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(putConfigOptions.EngineConfig)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = secretsManager.Service.Request(request, nil)

	return
}

// GetConfig : Get the configuration of a secret type
// Get the configuration that is associated with the specified secret type.
func (secretsManager *SecretsManagerV1) GetConfig(getConfigOptions *GetConfigOptions) (result *GetConfig, response *core.DetailedResponse, err error) {
	return secretsManager.GetConfigWithContext(context.Background(), getConfigOptions)
}

// GetConfigWithContext is an alternate form of the GetConfig method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetConfigWithContext(ctx context.Context, getConfigOptions *GetConfigOptions) (result *GetConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigOptions, "getConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getConfigOptions, "getConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type": *getConfigOptions.SecretType,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/config/{secret_type}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetConfig")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateConfigElement : Add a configuration
// Add a configuration element to the specified secret type.
//
// Use this method to define the configurations that are required to enable the public certificates (`public_cert`) and
// private certificates (`private_cert`) engines.
//
// You can add multiple configurations for your instance as follows:
//
// - Up to 10 public certificate authority configurations
// - Up to 10 DNS provider configurations
// - Up to 10 private root certificate authority configurations
// - Up to 10 private intermediate certificate authority configurations
// - Up to 10 certificate templates.
func (secretsManager *SecretsManagerV1) CreateConfigElement(createConfigElementOptions *CreateConfigElementOptions) (result *GetSingleConfigElement, response *core.DetailedResponse, err error) {
	return secretsManager.CreateConfigElementWithContext(context.Background(), createConfigElementOptions)
}

// CreateConfigElementWithContext is an alternate form of the CreateConfigElement method which supports a Context parameter
func (secretsManager *SecretsManagerV1) CreateConfigElementWithContext(ctx context.Context, createConfigElementOptions *CreateConfigElementOptions) (result *GetSingleConfigElement, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createConfigElementOptions, "createConfigElementOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createConfigElementOptions, "createConfigElementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type":    *createConfigElementOptions.SecretType,
		"config_element": *createConfigElementOptions.ConfigElement,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/config/{secret_type}/{config_element}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createConfigElementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "CreateConfigElement")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createConfigElementOptions.Name != nil {
		body["name"] = createConfigElementOptions.Name
	}
	if createConfigElementOptions.Type != nil {
		body["type"] = createConfigElementOptions.Type
	}
	if createConfigElementOptions.Config != nil {
		body["config"] = createConfigElementOptions.Config
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSingleConfigElement)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetConfigElements : List configurations
// List the configuration elements that are associated with a specified secret type.
func (secretsManager *SecretsManagerV1) GetConfigElements(getConfigElementsOptions *GetConfigElementsOptions) (result *GetConfigElements, response *core.DetailedResponse, err error) {
	return secretsManager.GetConfigElementsWithContext(context.Background(), getConfigElementsOptions)
}

// GetConfigElementsWithContext is an alternate form of the GetConfigElements method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetConfigElementsWithContext(ctx context.Context, getConfigElementsOptions *GetConfigElementsOptions) (result *GetConfigElements, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigElementsOptions, "getConfigElementsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getConfigElementsOptions, "getConfigElementsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type":    *getConfigElementsOptions.SecretType,
		"config_element": *getConfigElementsOptions.ConfigElement,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/config/{secret_type}/{config_element}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getConfigElementsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetConfigElements")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetConfigElements)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetConfigElement : Get a configuration
// Get the details of a specific configuration that is associated with a secret type.
func (secretsManager *SecretsManagerV1) GetConfigElement(getConfigElementOptions *GetConfigElementOptions) (result *GetSingleConfigElement, response *core.DetailedResponse, err error) {
	return secretsManager.GetConfigElementWithContext(context.Background(), getConfigElementOptions)
}

// GetConfigElementWithContext is an alternate form of the GetConfigElement method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetConfigElementWithContext(ctx context.Context, getConfigElementOptions *GetConfigElementOptions) (result *GetSingleConfigElement, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigElementOptions, "getConfigElementOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getConfigElementOptions, "getConfigElementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type":    *getConfigElementOptions.SecretType,
		"config_element": *getConfigElementOptions.ConfigElement,
		"config_name":    *getConfigElementOptions.ConfigName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/config/{secret_type}/{config_element}/{config_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getConfigElementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetConfigElement")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSingleConfigElement)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateConfigElement : Update a configuration
// Update a configuration element that is associated with the specified secret type.
func (secretsManager *SecretsManagerV1) UpdateConfigElement(updateConfigElementOptions *UpdateConfigElementOptions) (result *GetSingleConfigElement, response *core.DetailedResponse, err error) {
	return secretsManager.UpdateConfigElementWithContext(context.Background(), updateConfigElementOptions)
}

// UpdateConfigElementWithContext is an alternate form of the UpdateConfigElement method which supports a Context parameter
func (secretsManager *SecretsManagerV1) UpdateConfigElementWithContext(ctx context.Context, updateConfigElementOptions *UpdateConfigElementOptions) (result *GetSingleConfigElement, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateConfigElementOptions, "updateConfigElementOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateConfigElementOptions, "updateConfigElementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type":    *updateConfigElementOptions.SecretType,
		"config_element": *updateConfigElementOptions.ConfigElement,
		"config_name":    *updateConfigElementOptions.ConfigName,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/config/{secret_type}/{config_element}/{config_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateConfigElementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "UpdateConfigElement")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateConfigElementOptions.Type != nil {
		body["type"] = updateConfigElementOptions.Type
	}
	if updateConfigElementOptions.Config != nil {
		body["config"] = updateConfigElementOptions.Config
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSingleConfigElement)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ActionOnConfigElement : Invoke an action on a configuration
// Invoke an action on a specified configuration element. This method supports the following actions:
//
// - `sign_intermediate`: Sign an intermediate certificate authority.
// - `sign_csr`: Sign a certificate signing request.
// - `set_signed`: Set a signed intermediate certificate authority.
// - `revoke`: Revoke an internally signed intermediate certificate authority certificate.
// - `rotate_crl`: Rotate the certificate revocation list (CRL) of an intermediate certificate authority.
func (secretsManager *SecretsManagerV1) ActionOnConfigElement(actionOnConfigElementOptions *ActionOnConfigElementOptions) (result *ConfigElementActionResult, response *core.DetailedResponse, err error) {
	return secretsManager.ActionOnConfigElementWithContext(context.Background(), actionOnConfigElementOptions)
}

// ActionOnConfigElementWithContext is an alternate form of the ActionOnConfigElement method which supports a Context parameter
func (secretsManager *SecretsManagerV1) ActionOnConfigElementWithContext(ctx context.Context, actionOnConfigElementOptions *ActionOnConfigElementOptions) (result *ConfigElementActionResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(actionOnConfigElementOptions, "actionOnConfigElementOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(actionOnConfigElementOptions, "actionOnConfigElementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type":    *actionOnConfigElementOptions.SecretType,
		"config_element": *actionOnConfigElementOptions.ConfigElement,
		"config_name":    *actionOnConfigElementOptions.ConfigName,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/config/{secret_type}/{config_element}/{config_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range actionOnConfigElementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "ActionOnConfigElement")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	builder.AddQuery("action", fmt.Sprint(*actionOnConfigElementOptions.Action))

	body := make(map[string]interface{})
	if actionOnConfigElementOptions.Config != nil {
		body["config"] = actionOnConfigElementOptions.Config
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigElementActionResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteConfigElement : Delete a configuration
// Delete a configuration element from the specified secret type.
func (secretsManager *SecretsManagerV1) DeleteConfigElement(deleteConfigElementOptions *DeleteConfigElementOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteConfigElementWithContext(context.Background(), deleteConfigElementOptions)
}

// DeleteConfigElementWithContext is an alternate form of the DeleteConfigElement method which supports a Context parameter
func (secretsManager *SecretsManagerV1) DeleteConfigElementWithContext(ctx context.Context, deleteConfigElementOptions *DeleteConfigElementOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteConfigElementOptions, "deleteConfigElementOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteConfigElementOptions, "deleteConfigElementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"secret_type":    *deleteConfigElementOptions.SecretType,
		"config_element": *deleteConfigElementOptions.ConfigElement,
		"config_name":    *deleteConfigElementOptions.ConfigName,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/config/{secret_type}/{config_element}/{config_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteConfigElementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "DeleteConfigElement")
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

// CreateNotificationsRegistration : Register with Event Notifications
// Create a registration between a Secrets Manager instance and [Event
// Notifications](https://cloud.ibm.com/apidocs/event-notifications).
//
// A successful request adds Secrets Manager as a source that you can reference from your Event Notifications instance.
// For more information about enabling notifications for Secrets Manager, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-event-notifications).
func (secretsManager *SecretsManagerV1) CreateNotificationsRegistration(createNotificationsRegistrationOptions *CreateNotificationsRegistrationOptions) (result *GetNotificationsSettings, response *core.DetailedResponse, err error) {
	return secretsManager.CreateNotificationsRegistrationWithContext(context.Background(), createNotificationsRegistrationOptions)
}

// CreateNotificationsRegistrationWithContext is an alternate form of the CreateNotificationsRegistration method which supports a Context parameter
func (secretsManager *SecretsManagerV1) CreateNotificationsRegistrationWithContext(ctx context.Context, createNotificationsRegistrationOptions *CreateNotificationsRegistrationOptions) (result *GetNotificationsSettings, response *core.DetailedResponse, err error) {
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
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/notifications/registration`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createNotificationsRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "CreateNotificationsRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createNotificationsRegistrationOptions.EventNotificationsInstanceCRN != nil {
		body["event_notifications_instance_crn"] = createNotificationsRegistrationOptions.EventNotificationsInstanceCRN
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetNotificationsSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetNotificationsRegistration : Get Event Notifications registration details
// Get the details of an existing registration between a Secrets Manager instance and Event Notifications.
func (secretsManager *SecretsManagerV1) GetNotificationsRegistration(getNotificationsRegistrationOptions *GetNotificationsRegistrationOptions) (result *GetNotificationsSettings, response *core.DetailedResponse, err error) {
	return secretsManager.GetNotificationsRegistrationWithContext(context.Background(), getNotificationsRegistrationOptions)
}

// GetNotificationsRegistrationWithContext is an alternate form of the GetNotificationsRegistration method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetNotificationsRegistrationWithContext(ctx context.Context, getNotificationsRegistrationOptions *GetNotificationsRegistrationOptions) (result *GetNotificationsSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getNotificationsRegistrationOptions, "getNotificationsRegistrationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/notifications/registration`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getNotificationsRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "GetNotificationsRegistration")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetNotificationsSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteNotificationsRegistration : Unregister from Event Notifications
// Delete a registration between a Secrets Manager instance and Event Notifications.
//
// A successful request removes your Secrets Manager instance as a source in Event Notifications.
func (secretsManager *SecretsManagerV1) DeleteNotificationsRegistration(deleteNotificationsRegistrationOptions *DeleteNotificationsRegistrationOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.DeleteNotificationsRegistrationWithContext(context.Background(), deleteNotificationsRegistrationOptions)
}

// DeleteNotificationsRegistrationWithContext is an alternate form of the DeleteNotificationsRegistration method which supports a Context parameter
func (secretsManager *SecretsManagerV1) DeleteNotificationsRegistrationWithContext(ctx context.Context, deleteNotificationsRegistrationOptions *DeleteNotificationsRegistrationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(deleteNotificationsRegistrationOptions, "deleteNotificationsRegistrationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/notifications/registration`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteNotificationsRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "DeleteNotificationsRegistration")
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

// SendTestNotification : Send a test event
// Send a test event from a Secrets Manager instance to a configured [Event
// Notifications](https://cloud.ibm.com/apidocs/event-notifications) instance.
//
// A successful request sends a test event to the Event Notifications instance. For more information about enabling
// notifications for Secrets Manager, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-event-notifications).
func (secretsManager *SecretsManagerV1) SendTestNotification(sendTestNotificationOptions *SendTestNotificationOptions) (response *core.DetailedResponse, err error) {
	return secretsManager.SendTestNotificationWithContext(context.Background(), sendTestNotificationOptions)
}

// SendTestNotificationWithContext is an alternate form of the SendTestNotification method which supports a Context parameter
func (secretsManager *SecretsManagerV1) SendTestNotificationWithContext(ctx context.Context, sendTestNotificationOptions *SendTestNotificationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(sendTestNotificationOptions, "sendTestNotificationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = secretsManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(secretsManager.Service.Options.URL, `/api/v1/notifications/test`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range sendTestNotificationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("secrets_manager", "V1", "SendTestNotification")
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

// ActionOnConfigElementOptions : The ActionOnConfigElement options.
type ActionOnConfigElementOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The configuration element on which the action is applied.
	ConfigElement *string `json:"config_element" validate:"required,ne="`

	// The name of the certificate authority.
	ConfigName *string `json:"config_name" validate:"required,ne="`

	// The action to perform on the specified configuration element.
	Action *string `json:"action" validate:"required"`

	// Properties that describe an action on a configuration element.
	Config ConfigActionIntf `json:"config,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ActionOnConfigElementOptions.SecretType property.
// The secret type.
const (
	ActionOnConfigElementOptionsSecretTypePrivateCertConst = "private_cert"
)

// Constants associated with the ActionOnConfigElementOptions.ConfigElement property.
// The configuration element on which the action is applied.
const (
	ActionOnConfigElementOptionsConfigElementIntermediateCertificateAuthoritiesConst = "intermediate_certificate_authorities"
	ActionOnConfigElementOptionsConfigElementRootCertificateAuthoritiesConst         = "root_certificate_authorities"
)

// Constants associated with the ActionOnConfigElementOptions.Action property.
// The action to perform on the specified configuration element.
const (
	ActionOnConfigElementOptionsActionRevokeConst           = "revoke"
	ActionOnConfigElementOptionsActionRotateCrlConst        = "rotate_crl"
	ActionOnConfigElementOptionsActionSetSignedConst        = "set_signed"
	ActionOnConfigElementOptionsActionSignCsrConst          = "sign_csr"
	ActionOnConfigElementOptionsActionSignIntermediateConst = "sign_intermediate"
)

// NewActionOnConfigElementOptions : Instantiate ActionOnConfigElementOptions
func (*SecretsManagerV1) NewActionOnConfigElementOptions(secretType string, configElement string, configName string, action string) *ActionOnConfigElementOptions {
	return &ActionOnConfigElementOptions{
		SecretType:    core.StringPtr(secretType),
		ConfigElement: core.StringPtr(configElement),
		ConfigName:    core.StringPtr(configName),
		Action:        core.StringPtr(action),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *ActionOnConfigElementOptions) SetSecretType(secretType string) *ActionOnConfigElementOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetConfigElement : Allow user to set ConfigElement
func (_options *ActionOnConfigElementOptions) SetConfigElement(configElement string) *ActionOnConfigElementOptions {
	_options.ConfigElement = core.StringPtr(configElement)
	return _options
}

// SetConfigName : Allow user to set ConfigName
func (_options *ActionOnConfigElementOptions) SetConfigName(configName string) *ActionOnConfigElementOptions {
	_options.ConfigName = core.StringPtr(configName)
	return _options
}

// SetAction : Allow user to set Action
func (_options *ActionOnConfigElementOptions) SetAction(action string) *ActionOnConfigElementOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetConfig : Allow user to set Config
func (_options *ActionOnConfigElementOptions) SetConfig(config ConfigActionIntf) *ActionOnConfigElementOptions {
	_options.Config = config
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ActionOnConfigElementOptions) SetHeaders(param map[string]string) *ActionOnConfigElementOptions {
	options.Headers = param
	return options
}

// CertificateSecretData : The data that is associated with the secret version. The data object contains the following fields:
//
// - `certificate`: The contents of the certificate.
// - `private_key`: The private key that is associated with the certificate.
// - `intermediate`: The intermediate certificate that is associated with the certificate.
type CertificateSecretData struct {

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of CertificateSecretData
func (o *CertificateSecretData) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of CertificateSecretData
func (o *CertificateSecretData) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of CertificateSecretData
func (o *CertificateSecretData) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of CertificateSecretData
func (o *CertificateSecretData) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of CertificateSecretData
func (o *CertificateSecretData) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalCertificateSecretData unmarshals an instance of CertificateSecretData from the specified map of raw messages.
func UnmarshalCertificateSecretData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateSecretData)
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

// CertificateTemplatesConfigItem : Certificate templates configuration.
type CertificateTemplatesConfigItem struct {
	// The human-readable name to assign to your configuration.
	Name *string `json:"name" validate:"required"`

	// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
	Type *string `json:"type" validate:"required"`

	// Properties that describe a certificate template. You can use a certificate template to control the parameters that
	// are applied to your issued private certificates. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-certificate-templates).
	Config *CertificateTemplateConfig `json:"config,omitempty"`
}

// Constants associated with the CertificateTemplatesConfigItem.Type property.
// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
const (
	CertificateTemplatesConfigItemTypeCertificateTemplateConst              = "certificate_template"
	CertificateTemplatesConfigItemTypeCisConst                              = "cis"
	CertificateTemplatesConfigItemTypeClassicInfrastructureConst            = "classic_infrastructure"
	CertificateTemplatesConfigItemTypeIntermediateCertificateAuthorityConst = "intermediate_certificate_authority"
	CertificateTemplatesConfigItemTypeLetsencryptConst                      = "letsencrypt"
	CertificateTemplatesConfigItemTypeLetsencryptStageConst                 = "letsencrypt-stage"
	CertificateTemplatesConfigItemTypeRootCertificateAuthorityConst         = "root_certificate_authority"
)

// UnmarshalCertificateTemplatesConfigItem unmarshals an instance of CertificateTemplatesConfigItem from the specified map of raw messages.
func UnmarshalCertificateTemplatesConfigItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateTemplatesConfigItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalCertificateTemplateConfig)
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

	// The challenge expiration date. The date format follows RFC 3339.
	Expiration *strfmt.DateTime `json:"expiration,omitempty"`

	// The challenge status.
	Status *string `json:"status,omitempty"`

	// The txt_record_name.
	TxtRecordName *string `json:"txt_record_name,omitempty"`

	// The txt_record_value.
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

// CollectionMetadata : The metadata that describes the resource array.
type CollectionMetadata struct {
	// The type of resources in the resource array.
	CollectionType *string `json:"collection_type" validate:"required"`

	// The number of elements in the resource array.
	CollectionTotal *int64 `json:"collection_total" validate:"required"`
}

// Constants associated with the CollectionMetadata.CollectionType property.
// The type of resources in the resource array.
const (
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerConfigJSONConst        = "application/vnd.ibm.secrets-manager.config+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerErrorJSONConst         = "application/vnd.ibm.secrets-manager.error+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretGroupJSONConst   = "application/vnd.ibm.secrets-manager.secret.group+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretJSONConst        = "application/vnd.ibm.secrets-manager.secret+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretLockJSONConst    = "application/vnd.ibm.secrets-manager.secret.lock+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretPolicyJSONConst  = "application/vnd.ibm.secrets-manager.secret.policy+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretVersionJSONConst = "application/vnd.ibm.secrets-manager.secret.version+json"
)

// NewCollectionMetadata : Instantiate CollectionMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewCollectionMetadata(collectionType string, collectionTotal int64) (_model *CollectionMetadata, err error) {
	_model = &CollectionMetadata{
		CollectionType:  core.StringPtr(collectionType),
		CollectionTotal: core.Int64Ptr(collectionTotal),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalCollectionMetadata unmarshals an instance of CollectionMetadata from the specified map of raw messages.
func UnmarshalCollectionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CollectionMetadata)
	err = core.UnmarshalPrimitive(m, "collection_type", &obj.CollectionType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "collection_total", &obj.CollectionTotal)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigAction : Properties that describe an action on a configuration element.
// Models which "extend" this model:
// - SignCsrAction
// - SignIntermediateAction
// - SetSignedAction
// - RevokeAction
type ConfigAction struct {
	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `12h`. The value can't exceed
	// the `max_ttl` that is defined in the associated certificate template.
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	PermittedDNSDomains []string `json:"permitted_dns_domains,omitempty"`

	// Determines whether to use values from a certificate signing request (CSR) to complete a `sign_csr` action. If set to
	// `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than using the values
	// provided in the other parameters to this operation.
	//
	// 2) Any key usages (for example, non-repudiation) that are requested in the CSR are added to the basic set of key
	// usages used for CA certs signed by this intermediate authority.
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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The PEM-encoded certificate signing request (CSR). This field is required for the `sign_csr` action.
	Csr *string `json:"csr,omitempty"`

	// The intermediate certificate authority to be signed. The name must match one of the pre-configured intermediate
	// certificate authorities.
	IntermediateCertificateAuthority *string `json:"intermediate_certificate_authority,omitempty"`

	// The PEM-encoded certificate.
	Certificate *string `json:"certificate,omitempty"`
}

// Constants associated with the ConfigAction.Format property.
// The format of the returned data.
const (
	ConfigActionFormatPemConst       = "pem"
	ConfigActionFormatPemBundleConst = "pem_bundle"
)

func (*ConfigAction) isaConfigAction() bool {
	return true
}

type ConfigActionIntf interface {
	isaConfigAction() bool
}

// UnmarshalConfigAction unmarshals an instance of ConfigAction from the specified map of raw messages.
func UnmarshalConfigAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigAction)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDNSDomains)
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
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_certificate_authority", &obj.IntermediateCertificateAuthority)
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

// ConfigElementActionData : The configuration to add or update.
type ConfigElementActionData struct {
	// The human-readable name to assign to your configuration.
	Name *string `json:"name" validate:"required"`

	// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
	Type *string `json:"type" validate:"required"`

	Config ConfigElementActionResultConfigIntf `json:"config" validate:"required"`
}

// Constants associated with the ConfigElementActionData.Type property.
// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
const (
	ConfigElementActionDataTypeCertificateTemplateConst              = "certificate_template"
	ConfigElementActionDataTypeCisConst                              = "cis"
	ConfigElementActionDataTypeClassicInfrastructureConst            = "classic_infrastructure"
	ConfigElementActionDataTypeIntermediateCertificateAuthorityConst = "intermediate_certificate_authority"
	ConfigElementActionDataTypeLetsencryptConst                      = "letsencrypt"
	ConfigElementActionDataTypeLetsencryptStageConst                 = "letsencrypt-stage"
	ConfigElementActionDataTypeRootCertificateAuthorityConst         = "root_certificate_authority"
)

// UnmarshalConfigElementActionData unmarshals an instance of ConfigElementActionData from the specified map of raw messages.
func UnmarshalConfigElementActionData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementActionData)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalConfigElementActionResultConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigElementActionResult : Properties that describe an action on a configuration element.
type ConfigElementActionResult struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []ConfigElementActionData `json:"resources" validate:"required"`
}

// UnmarshalConfigElementActionResult unmarshals an instance of ConfigElementActionResult from the specified map of raw messages.
func UnmarshalConfigElementActionResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementActionResult)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalConfigElementActionData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigElementActionResultConfig : ConfigElementActionResultConfig struct
// Models which "extend" this model:
// - SignCsrActionResult
// - SignIntermediateActionResult
// - RotateCrlActionResult
// - SetSignedActionResult
// - RevokeActionResult
type ConfigElementActionResultConfig struct {
	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `12h`. The value can't exceed
	// the `max_ttl` that is defined in the associated certificate template.
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	PermittedDNSDomains []string `json:"permitted_dns_domains,omitempty"`

	// Determines whether to use values from a certificate signing request (CSR) to complete a `sign_csr` action. If set to
	// `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than using the values
	// provided in the other parameters to this operation.
	//
	// 2) Any key usages (for example, non-repudiation) that are requested in the CSR are added to the basic set of key
	// usages used for CA certs signed by this intermediate authority.
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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// Properties that are returned with a successful `sign` action.
	Data *SignActionResultData `json:"data,omitempty"`

	// The PEM-encoded certificate signing request (CSR).
	Csr *string `json:"csr,omitempty"`

	// The signed intermediate certificate authority.
	IntermediateCertificateAuthority *string `json:"intermediate_certificate_authority,omitempty"`

	// The time until the certificate authority is revoked.
	RevocationTime *int64 `json:"revocation_time,omitempty"`
}

// Constants associated with the ConfigElementActionResultConfig.Format property.
// The format of the returned data.
const (
	ConfigElementActionResultConfigFormatPemConst       = "pem"
	ConfigElementActionResultConfigFormatPemBundleConst = "pem_bundle"
)

func (*ConfigElementActionResultConfig) isaConfigElementActionResultConfig() bool {
	return true
}

type ConfigElementActionResultConfigIntf interface {
	isaConfigElementActionResultConfig() bool
}

// UnmarshalConfigElementActionResultConfig unmarshals an instance of ConfigElementActionResultConfig from the specified map of raw messages.
func UnmarshalConfigElementActionResultConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementActionResultConfig)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDNSDomains)
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
	err = core.UnmarshalModel(m, "data", &obj.Data, UnmarshalSignActionResultData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_certificate_authority", &obj.IntermediateCertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigElementDef : The configuration to add or update.
type ConfigElementDef struct {
	// The human-readable name to assign to your configuration.
	Name *string `json:"name" validate:"required"`

	// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
	Type *string `json:"type" validate:"required"`

	// The configuration to define for the specified secret type.
	Config ConfigElementDefConfigIntf `json:"config" validate:"required"`
}

// Constants associated with the ConfigElementDef.Type property.
// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
const (
	ConfigElementDefTypeCertificateTemplateConst              = "certificate_template"
	ConfigElementDefTypeCisConst                              = "cis"
	ConfigElementDefTypeClassicInfrastructureConst            = "classic_infrastructure"
	ConfigElementDefTypeIntermediateCertificateAuthorityConst = "intermediate_certificate_authority"
	ConfigElementDefTypeLetsencryptConst                      = "letsencrypt"
	ConfigElementDefTypeLetsencryptStageConst                 = "letsencrypt-stage"
	ConfigElementDefTypeRootCertificateAuthorityConst         = "root_certificate_authority"
)

// NewConfigElementDef : Instantiate ConfigElementDef (Generic Model Constructor)
func (*SecretsManagerV1) NewConfigElementDef(name string, typeVar string, config ConfigElementDefConfigIntf) (_model *ConfigElementDef, err error) {
	_model = &ConfigElementDef{
		Name:   core.StringPtr(name),
		Type:   core.StringPtr(typeVar),
		Config: config,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalConfigElementDef unmarshals an instance of ConfigElementDef from the specified map of raw messages.
func UnmarshalConfigElementDef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementDef)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalConfigElementDefConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigElementDefConfig : The configuration to define for the specified secret type.
// Models which "extend" this model:
// - ConfigElementDefConfigLetsEncryptConfig
// - ConfigElementDefConfigCloudInternetServicesConfig
// - ConfigElementDefConfigClassicInfrastructureConfig
// - RootCertificateAuthorityConfig
// - IntermediateCertificateAuthorityConfig
// - CertificateTemplateConfig
type ConfigElementDefConfig struct {
	// The private key that is associated with your Automatic Certificate Management Environment (ACME) account.
	//
	// If you have a working ACME client or account for Let's Encrypt, you can use the existing private key to enable
	// communications with Secrets Manager. If you don't have an account yet, you can create one. For more information, see
	// the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#create-acme-account).
	PrivateKey *string `json:"private_key,omitempty"`

	// The Cloud Resource Name (CRN) that is associated with the CIS instance.
	CisCRN *string `json:"cis_crn,omitempty"`

	// An IBM Cloud API key that can to list domains in your CIS instance.
	//
	// To grant Secrets Manager the ability to view the CIS instance and all of its domains, the API key must be assigned
	// the Reader service role on Internet Services (`internet-svcs`).
	//
	// If you need to manage specific domains, you can assign the Manager role. For production environments, it is
	// recommended that you assign the Reader access role, and then use the
	// [IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific
	// domains. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	CisApikey *string `json:"cis_apikey,omitempty"`

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

	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL interface{} `json:"max_ttl,omitempty"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry interface{} `json:"crl_expiry,omitempty"`

	// Disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is
	// enabled,  it will rebuild the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are
	// issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this
	// certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to this CA certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `12h`. The value can be supplied in
	// seconds (suffix `s`), minutes (suffix `m`), hours (suffix `h`) or days (suffix `d`). The value can't exceed the
	// `max_ttl` that is defined in the associated certificate template. In the API response, this value is returned in
	// seconds (integer).
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use when generating the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	PermittedDNSDomains []string `json:"permitted_dns_domains,omitempty"`

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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The data that is associated with the root certificate authority. The data object contains the following fields:
	//
	// - `certificate`: The root certificate content.
	// - `issuing_ca`: The certificate of the certificate authority that signed and issued this certificate.
	// - `serial_number`: The unique serial number of the root certificate.
	Data map[string]interface{} `json:"data,omitempty"`

	// The signing method to use with this certificate authority to generate private certificates.
	//
	// You can choose between internal or externally signed options. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	SigningMethod *string `json:"signing_method,omitempty"`

	// The certificate authority that signed and issued the certificate.
	//
	// If the certificate is signed internally, the `issuer` field is required and must match the name of a certificate
	// authority that is configured in the Secrets Manager service instance.
	Issuer *string `json:"issuer,omitempty"`

	// The name of the intermediate certificate authority.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// Scopes the creation of private certificates to only the secret groups that you specify.
	//
	// This field can be supplied as a comma-delimited list of secret group IDs.
	AllowedSecretGroups *string `json:"allowed_secret_groups,omitempty"`

	// Determines whether to allow `localhost` to be included as one of the requested common names.
	AllowLocalhost *bool `json:"allow_localhost,omitempty"`

	// The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and
	// `allow_subdomains` options.
	AllowedDomains []string `json:"allowed_domains,omitempty"`

	// Determines whether to allow the domains that are supplied in the `allowed_domains` field to contain access control
	// list (ACL) templates.
	AllowedDomainsTemplate *bool `json:"allowed_domains_template,omitempty"`

	// Determines whether to allow clients to request private certificates that match the value of the actual domains on
	// the final certificate.
	//
	// For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a
	// certificate that contains the name `example.com` as one of the DNS values on the final certificate.
	//
	// **Important:** In some scenarios, allowing bare domains can be considered a security risk.
	AllowBareDomains *bool `json:"allow_bare_domains,omitempty"`

	// Determines whether to allow clients to request private certificates with common names (CN) that are subdomains of
	// the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.
	//
	// For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the
	// following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.
	//
	// **Note:** This field is redundant if you use the `allow_any_name` option.
	AllowSubdomains *bool `json:"allow_subdomains,omitempty"`

	// Determines whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are specified in the
	// `allowed_domains` field.
	//
	// If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	AllowGlobDomains *bool `json:"allow_glob_domains,omitempty"`

	// Determines whether to allow clients to request a private certificate that matches any common name.
	AllowAnyName *bool `json:"allow_any_name,omitempty"`

	// Determines whether to enforce only valid host names for common names, DNS Subject Alternative Names, and the host
	// section of email addresses.
	EnforceHostnames *bool `json:"enforce_hostnames,omitempty"`

	// Determines whether to allow clients to request a private certificate with IP Subject Alternative Names.
	AllowIPSans *bool `json:"allow_ip_sans,omitempty"`

	// The URI Subject Alternative Names to allow for private certificates.
	//
	// Values can contain glob patterns, for example `spiffe://hostname/_*`.
	AllowedURISans []string `json:"allowed_uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private
	// certificates.
	//
	// The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type
	// is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any
	// `other_sans` input.
	AllowedOtherSans []string `json:"allowed_other_sans,omitempty"`

	// Determines whether private certificates are flagged for server use.
	ServerFlag *bool `json:"server_flag,omitempty"`

	// Determines whether private certificates are flagged for client use.
	ClientFlag *bool `json:"client_flag,omitempty"`

	// Determines whether private certificates are flagged for code signing use.
	CodeSigningFlag *bool `json:"code_signing_flag,omitempty"`

	// Determines whether private certificates are flagged for email protection use.
	EmailProtectionFlag *bool `json:"email_protection_flag,omitempty"`

	// The allowed key usage constraint to define for private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage).  Omit the
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

	// When used with the `sign_csr` action, this field determines whether to use the common name (CN) from a certificate
	// signing request (CSR) instead of the CN that's included in the JSON data of the certificate.
	//
	// Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include
	// the `use_csr_sans` property.
	UseCsrCommonName *bool `json:"use_csr_common_name,omitempty"`

	// When used with the `sign_csr` action, this field determines whether to use the Subject Alternative Names
	// (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the JSON data of the
	// certificate.
	//
	// Does not include the common name in the CSR. To use the common name, include the `use_csr_common_name` property.
	UseCsrSans *bool `json:"use_csr_sans,omitempty"`

	// Determines whether to require a common name to create a private certificate.
	//
	// By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the
	// `require_cn` option to `false`.
	RequireCn *bool `json:"require_cn,omitempty"`

	// A list of policy Object Identifiers (OIDs).
	PolicyIdentifiers []string `json:"policy_identifiers,omitempty"`

	// Determines whether to mark the Basic Constraints extension of an issued private certificate as valid for non-CA
	// certificates.
	BasicConstraintsValidForNonCa *bool `json:"basic_constraints_valid_for_non_ca,omitempty"`

	// The duration in seconds by which to backdate the `not_before` property of an issued private certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value
	// is returned in seconds (integer).
	NotBeforeDuration interface{} `json:"not_before_duration,omitempty"`
}

// Constants associated with the ConfigElementDefConfig.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	ConfigElementDefConfigStatusCertificateTemplateRequiredConst = "certificate_template_required"
	ConfigElementDefConfigStatusConfiguredConst                  = "configured"
	ConfigElementDefConfigStatusExpiredConst                     = "expired"
	ConfigElementDefConfigStatusRevokedConst                     = "revoked"
	ConfigElementDefConfigStatusSignedCertificateRequiredConst   = "signed_certificate_required"
	ConfigElementDefConfigStatusSigningRequiredConst             = "signing_required"
)

// Constants associated with the ConfigElementDefConfig.Format property.
// The format of the returned data.
const (
	ConfigElementDefConfigFormatPemConst       = "pem"
	ConfigElementDefConfigFormatPemBundleConst = "pem_bundle"
)

// Constants associated with the ConfigElementDefConfig.PrivateKeyFormat property.
// The format of the generated private key.
const (
	ConfigElementDefConfigPrivateKeyFormatDerConst   = "der"
	ConfigElementDefConfigPrivateKeyFormatPkcs8Const = "pkcs8"
)

// Constants associated with the ConfigElementDefConfig.KeyType property.
// The type of private key to generate.
const (
	ConfigElementDefConfigKeyTypeEcConst  = "ec"
	ConfigElementDefConfigKeyTypeRsaConst = "rsa"
)

// Constants associated with the ConfigElementDefConfig.SigningMethod property.
// The signing method to use with this certificate authority to generate private certificates.
//
// You can choose between internal or externally signed options. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
const (
	ConfigElementDefConfigSigningMethodExternalConst = "external"
	ConfigElementDefConfigSigningMethodInternalConst = "internal"
)

func (*ConfigElementDefConfig) isaConfigElementDefConfig() bool {
	return true
}

type ConfigElementDefConfigIntf interface {
	isaConfigElementDefConfig() bool
}

// UnmarshalConfigElementDefConfig unmarshals an instance of ConfigElementDefConfig from the specified map of raw messages.
func UnmarshalConfigElementDefConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementDefConfig)
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cis_crn", &obj.CisCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cis_apikey", &obj.CisApikey)
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
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDNSDomains)
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
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
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
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_secret_groups", &obj.AllowedSecretGroups)
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
	err = core.UnmarshalPrimitive(m, "allow_ip_sans", &obj.AllowIPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_uri_sans", &obj.AllowedURISans)
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

// ConfigElementMetadata : Properties that describe a configuration element.
type ConfigElementMetadata struct {
	// The human-readable name to assign to your configuration.
	Name *string `json:"name" validate:"required"`

	// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
	Type *string `json:"type" validate:"required"`
}

// Constants associated with the ConfigElementMetadata.Type property.
// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
const (
	ConfigElementMetadataTypeCertificateTemplateConst              = "certificate_template"
	ConfigElementMetadataTypeCisConst                              = "cis"
	ConfigElementMetadataTypeClassicInfrastructureConst            = "classic_infrastructure"
	ConfigElementMetadataTypeIntermediateCertificateAuthorityConst = "intermediate_certificate_authority"
	ConfigElementMetadataTypeLetsencryptConst                      = "letsencrypt"
	ConfigElementMetadataTypeLetsencryptStageConst                 = "letsencrypt-stage"
	ConfigElementMetadataTypeRootCertificateAuthorityConst         = "root_certificate_authority"
)

// UnmarshalConfigElementMetadata unmarshals an instance of ConfigElementMetadata from the specified map of raw messages.
func UnmarshalConfigElementMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementMetadata)
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

// CreateConfigElementOptions : The CreateConfigElement options.
type CreateConfigElementOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The configuration element to define or manage.
	ConfigElement *string `json:"config_element" validate:"required,ne="`

	// The human-readable name to assign to your configuration.
	Name *string `json:"name" validate:"required"`

	// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
	Type *string `json:"type" validate:"required"`

	// The configuration to define for the specified secret type.
	Config ConfigElementDefConfigIntf `json:"config" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateConfigElementOptions.SecretType property.
// The secret type.
const (
	CreateConfigElementOptionsSecretTypePrivateCertConst = "private_cert"
	CreateConfigElementOptionsSecretTypePublicCertConst  = "public_cert"
)

// Constants associated with the CreateConfigElementOptions.ConfigElement property.
// The configuration element to define or manage.
const (
	CreateConfigElementOptionsConfigElementCertificateAuthoritiesConst             = "certificate_authorities"
	CreateConfigElementOptionsConfigElementCertificateTemplatesConst               = "certificate_templates"
	CreateConfigElementOptionsConfigElementDNSProvidersConst                       = "dns_providers"
	CreateConfigElementOptionsConfigElementIntermediateCertificateAuthoritiesConst = "intermediate_certificate_authorities"
	CreateConfigElementOptionsConfigElementRootCertificateAuthoritiesConst         = "root_certificate_authorities"
)

// Constants associated with the CreateConfigElementOptions.Type property.
// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
const (
	CreateConfigElementOptionsTypeCertificateTemplateConst              = "certificate_template"
	CreateConfigElementOptionsTypeCisConst                              = "cis"
	CreateConfigElementOptionsTypeClassicInfrastructureConst            = "classic_infrastructure"
	CreateConfigElementOptionsTypeIntermediateCertificateAuthorityConst = "intermediate_certificate_authority"
	CreateConfigElementOptionsTypeLetsencryptConst                      = "letsencrypt"
	CreateConfigElementOptionsTypeLetsencryptStageConst                 = "letsencrypt-stage"
	CreateConfigElementOptionsTypeRootCertificateAuthorityConst         = "root_certificate_authority"
)

// NewCreateConfigElementOptions : Instantiate CreateConfigElementOptions
func (*SecretsManagerV1) NewCreateConfigElementOptions(secretType string, configElement string, name string, typeVar string, config ConfigElementDefConfigIntf) *CreateConfigElementOptions {
	return &CreateConfigElementOptions{
		SecretType:    core.StringPtr(secretType),
		ConfigElement: core.StringPtr(configElement),
		Name:          core.StringPtr(name),
		Type:          core.StringPtr(typeVar),
		Config:        config,
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *CreateConfigElementOptions) SetSecretType(secretType string) *CreateConfigElementOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetConfigElement : Allow user to set ConfigElement
func (_options *CreateConfigElementOptions) SetConfigElement(configElement string) *CreateConfigElementOptions {
	_options.ConfigElement = core.StringPtr(configElement)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateConfigElementOptions) SetName(name string) *CreateConfigElementOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateConfigElementOptions) SetType(typeVar string) *CreateConfigElementOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetConfig : Allow user to set Config
func (_options *CreateConfigElementOptions) SetConfig(config ConfigElementDefConfigIntf) *CreateConfigElementOptions {
	_options.Config = config
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateConfigElementOptions) SetHeaders(param map[string]string) *CreateConfigElementOptions {
	options.Headers = param
	return options
}

// CreateNotificationsRegistrationOptions : The CreateNotificationsRegistration options.
type CreateNotificationsRegistrationOptions struct {
	// The Cloud Resource Name (CRN) of the connected Event Notifications instance.
	EventNotificationsInstanceCRN *string `json:"event_notifications_instance_crn" validate:"required"`

	// The name that is displayed as a source in your Event Notifications instance.
	EventNotificationsSourceName *string `json:"event_notifications_source_name" validate:"required"`

	// An optional description for the source in your Event Notifications instance.
	EventNotificationsSourceDescription *string `json:"event_notifications_source_description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateNotificationsRegistrationOptions : Instantiate CreateNotificationsRegistrationOptions
func (*SecretsManagerV1) NewCreateNotificationsRegistrationOptions(eventNotificationsInstanceCRN string, eventNotificationsSourceName string) *CreateNotificationsRegistrationOptions {
	return &CreateNotificationsRegistrationOptions{
		EventNotificationsInstanceCRN: core.StringPtr(eventNotificationsInstanceCRN),
		EventNotificationsSourceName:  core.StringPtr(eventNotificationsSourceName),
	}
}

// SetEventNotificationsInstanceCRN : Allow user to set EventNotificationsInstanceCRN
func (_options *CreateNotificationsRegistrationOptions) SetEventNotificationsInstanceCRN(eventNotificationsInstanceCRN string) *CreateNotificationsRegistrationOptions {
	_options.EventNotificationsInstanceCRN = core.StringPtr(eventNotificationsInstanceCRN)
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

// CreateSecret : Properties that describe a secret.
type CreateSecret struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretResourceIntf `json:"resources" validate:"required"`
}

// NewCreateSecret : Instantiate CreateSecret (Generic Model Constructor)
func (*SecretsManagerV1) NewCreateSecret(metadata *CollectionMetadata, resources []SecretResourceIntf) (_model *CreateSecret, err error) {
	_model = &CreateSecret{
		Metadata:  metadata,
		Resources: resources,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalCreateSecret unmarshals an instance of CreateSecret from the specified map of raw messages.
func UnmarshalCreateSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateSecret)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateSecretGroupOptions : The CreateSecretGroup options.
type CreateSecretGroupOptions struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretGroupResource `json:"resources" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSecretGroupOptions : Instantiate CreateSecretGroupOptions
func (*SecretsManagerV1) NewCreateSecretGroupOptions(metadata *CollectionMetadata, resources []SecretGroupResource) *CreateSecretGroupOptions {
	return &CreateSecretGroupOptions{
		Metadata:  metadata,
		Resources: resources,
	}
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateSecretGroupOptions) SetMetadata(metadata *CollectionMetadata) *CreateSecretGroupOptions {
	_options.Metadata = metadata
	return _options
}

// SetResources : Allow user to set Resources
func (_options *CreateSecretGroupOptions) SetResources(resources []SecretGroupResource) *CreateSecretGroupOptions {
	_options.Resources = resources
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretGroupOptions) SetHeaders(param map[string]string) *CreateSecretGroupOptions {
	options.Headers = param
	return options
}

// CreateSecretOptions : The CreateSecret options.
type CreateSecretOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretResourceIntf `json:"resources" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateSecretOptions.SecretType property.
// The secret type.
const (
	CreateSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	CreateSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	CreateSecretOptionsSecretTypeImportedCertConst     = "imported_cert"
	CreateSecretOptionsSecretTypeKvConst               = "kv"
	CreateSecretOptionsSecretTypePrivateCertConst      = "private_cert"
	CreateSecretOptionsSecretTypePublicCertConst       = "public_cert"
	CreateSecretOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewCreateSecretOptions : Instantiate CreateSecretOptions
func (*SecretsManagerV1) NewCreateSecretOptions(secretType string, metadata *CollectionMetadata, resources []SecretResourceIntf) *CreateSecretOptions {
	return &CreateSecretOptions{
		SecretType: core.StringPtr(secretType),
		Metadata:   metadata,
		Resources:  resources,
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *CreateSecretOptions) SetSecretType(secretType string) *CreateSecretOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateSecretOptions) SetMetadata(metadata *CollectionMetadata) *CreateSecretOptions {
	_options.Metadata = metadata
	return _options
}

// SetResources : Allow user to set Resources
func (_options *CreateSecretOptions) SetResources(resources []SecretResourceIntf) *CreateSecretOptions {
	_options.Resources = resources
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretOptions) SetHeaders(param map[string]string) *CreateSecretOptions {
	options.Headers = param
	return options
}

// DeleteConfigElementOptions : The DeleteConfigElement options.
type DeleteConfigElementOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The configuration element to define or manage.
	ConfigElement *string `json:"config_element" validate:"required,ne="`

	// The name of your configuration.
	ConfigName *string `json:"config_name" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteConfigElementOptions.SecretType property.
// The secret type.
const (
	DeleteConfigElementOptionsSecretTypePrivateCertConst = "private_cert"
	DeleteConfigElementOptionsSecretTypePublicCertConst  = "public_cert"
)

// Constants associated with the DeleteConfigElementOptions.ConfigElement property.
// The configuration element to define or manage.
const (
	DeleteConfigElementOptionsConfigElementCertificateAuthoritiesConst             = "certificate_authorities"
	DeleteConfigElementOptionsConfigElementCertificateTemplatesConst               = "certificate_templates"
	DeleteConfigElementOptionsConfigElementDNSProvidersConst                       = "dns_providers"
	DeleteConfigElementOptionsConfigElementIntermediateCertificateAuthoritiesConst = "intermediate_certificate_authorities"
	DeleteConfigElementOptionsConfigElementRootCertificateAuthoritiesConst         = "root_certificate_authorities"
)

// NewDeleteConfigElementOptions : Instantiate DeleteConfigElementOptions
func (*SecretsManagerV1) NewDeleteConfigElementOptions(secretType string, configElement string, configName string) *DeleteConfigElementOptions {
	return &DeleteConfigElementOptions{
		SecretType:    core.StringPtr(secretType),
		ConfigElement: core.StringPtr(configElement),
		ConfigName:    core.StringPtr(configName),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *DeleteConfigElementOptions) SetSecretType(secretType string) *DeleteConfigElementOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetConfigElement : Allow user to set ConfigElement
func (_options *DeleteConfigElementOptions) SetConfigElement(configElement string) *DeleteConfigElementOptions {
	_options.ConfigElement = core.StringPtr(configElement)
	return _options
}

// SetConfigName : Allow user to set ConfigName
func (_options *DeleteConfigElementOptions) SetConfigName(configName string) *DeleteConfigElementOptions {
	_options.ConfigName = core.StringPtr(configName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteConfigElementOptions) SetHeaders(param map[string]string) *DeleteConfigElementOptions {
	options.Headers = param
	return options
}

// DeleteNotificationsRegistrationOptions : The DeleteNotificationsRegistration options.
type DeleteNotificationsRegistrationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteNotificationsRegistrationOptions : Instantiate DeleteNotificationsRegistrationOptions
func (*SecretsManagerV1) NewDeleteNotificationsRegistrationOptions() *DeleteNotificationsRegistrationOptions {
	return &DeleteNotificationsRegistrationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *DeleteNotificationsRegistrationOptions) SetHeaders(param map[string]string) *DeleteNotificationsRegistrationOptions {
	options.Headers = param
	return options
}

// DeleteSecretGroupOptions : The DeleteSecretGroup options.
type DeleteSecretGroupOptions struct {
	// The v4 UUID that uniquely identifies the secret group.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSecretGroupOptions : Instantiate DeleteSecretGroupOptions
func (*SecretsManagerV1) NewDeleteSecretGroupOptions(id string) *DeleteSecretGroupOptions {
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

// DeleteSecretOptions : The DeleteSecret options.
type DeleteSecretOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteSecretOptions.SecretType property.
// The secret type.
const (
	DeleteSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	DeleteSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	DeleteSecretOptionsSecretTypeImportedCertConst     = "imported_cert"
	DeleteSecretOptionsSecretTypeKvConst               = "kv"
	DeleteSecretOptionsSecretTypePrivateCertConst      = "private_cert"
	DeleteSecretOptionsSecretTypePublicCertConst       = "public_cert"
	DeleteSecretOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewDeleteSecretOptions : Instantiate DeleteSecretOptions
func (*SecretsManagerV1) NewDeleteSecretOptions(secretType string, id string) *DeleteSecretOptions {
	return &DeleteSecretOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *DeleteSecretOptions) SetSecretType(secretType string) *DeleteSecretOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
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

// EngineConfig : EngineConfig struct
// Models which "extend" this model:
// - CreateIamCredentialsSecretEngineRootConfig
type EngineConfig struct {
	// An IBM Cloud API key that can create and manage service IDs.
	//
	// The API key must be assigned the Editor platform role on the Access Groups Service and the Operator platform role on
	// the IAM Identity Service. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).
	APIKey *string `json:"api_key,omitempty"`

	// The hash value of the IBM Cloud API key that is used to create and manage service IDs.
	APIKeyHash *string `json:"api_key_hash,omitempty"`
}

func (*EngineConfig) isaEngineConfig() bool {
	return true
}

type EngineConfigIntf interface {
	isaEngineConfig() bool
}

// UnmarshalEngineConfig unmarshals an instance of EngineConfig from the specified map of raw messages.
func UnmarshalEngineConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EngineConfig)
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_hash", &obj.APIKeyHash)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetConfig : Configuration for the specified secret type.
type GetConfig struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []GetConfigResourcesItemIntf `json:"resources" validate:"required"`
}

// UnmarshalGetConfig unmarshals an instance of GetConfig from the specified map of raw messages.
func UnmarshalGetConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetConfig)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalGetConfigResourcesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetConfigElementOptions : The GetConfigElement options.
type GetConfigElementOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The configuration element to define or manage.
	ConfigElement *string `json:"config_element" validate:"required,ne="`

	// The name of your configuration.
	ConfigName *string `json:"config_name" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetConfigElementOptions.SecretType property.
// The secret type.
const (
	GetConfigElementOptionsSecretTypePrivateCertConst = "private_cert"
	GetConfigElementOptionsSecretTypePublicCertConst  = "public_cert"
)

// Constants associated with the GetConfigElementOptions.ConfigElement property.
// The configuration element to define or manage.
const (
	GetConfigElementOptionsConfigElementCertificateAuthoritiesConst             = "certificate_authorities"
	GetConfigElementOptionsConfigElementCertificateTemplatesConst               = "certificate_templates"
	GetConfigElementOptionsConfigElementDNSProvidersConst                       = "dns_providers"
	GetConfigElementOptionsConfigElementIntermediateCertificateAuthoritiesConst = "intermediate_certificate_authorities"
	GetConfigElementOptionsConfigElementRootCertificateAuthoritiesConst         = "root_certificate_authorities"
)

// NewGetConfigElementOptions : Instantiate GetConfigElementOptions
func (*SecretsManagerV1) NewGetConfigElementOptions(secretType string, configElement string, configName string) *GetConfigElementOptions {
	return &GetConfigElementOptions{
		SecretType:    core.StringPtr(secretType),
		ConfigElement: core.StringPtr(configElement),
		ConfigName:    core.StringPtr(configName),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetConfigElementOptions) SetSecretType(secretType string) *GetConfigElementOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetConfigElement : Allow user to set ConfigElement
func (_options *GetConfigElementOptions) SetConfigElement(configElement string) *GetConfigElementOptions {
	_options.ConfigElement = core.StringPtr(configElement)
	return _options
}

// SetConfigName : Allow user to set ConfigName
func (_options *GetConfigElementOptions) SetConfigName(configName string) *GetConfigElementOptions {
	_options.ConfigName = core.StringPtr(configName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigElementOptions) SetHeaders(param map[string]string) *GetConfigElementOptions {
	options.Headers = param
	return options
}

// GetConfigElements : Properties that describe a list of configurations.
type GetConfigElements struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []GetConfigElementsResourcesItemIntf `json:"resources" validate:"required"`
}

// UnmarshalGetConfigElements unmarshals an instance of GetConfigElements from the specified map of raw messages.
func UnmarshalGetConfigElements(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetConfigElements)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalGetConfigElementsResourcesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetConfigElementsOptions : The GetConfigElements options.
type GetConfigElementsOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The configuration element to define or manage.
	ConfigElement *string `json:"config_element" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetConfigElementsOptions.SecretType property.
// The secret type.
const (
	GetConfigElementsOptionsSecretTypePrivateCertConst = "private_cert"
	GetConfigElementsOptionsSecretTypePublicCertConst  = "public_cert"
)

// Constants associated with the GetConfigElementsOptions.ConfigElement property.
// The configuration element to define or manage.
const (
	GetConfigElementsOptionsConfigElementCertificateAuthoritiesConst             = "certificate_authorities"
	GetConfigElementsOptionsConfigElementCertificateTemplatesConst               = "certificate_templates"
	GetConfigElementsOptionsConfigElementDNSProvidersConst                       = "dns_providers"
	GetConfigElementsOptionsConfigElementIntermediateCertificateAuthoritiesConst = "intermediate_certificate_authorities"
	GetConfigElementsOptionsConfigElementRootCertificateAuthoritiesConst         = "root_certificate_authorities"
)

// NewGetConfigElementsOptions : Instantiate GetConfigElementsOptions
func (*SecretsManagerV1) NewGetConfigElementsOptions(secretType string, configElement string) *GetConfigElementsOptions {
	return &GetConfigElementsOptions{
		SecretType:    core.StringPtr(secretType),
		ConfigElement: core.StringPtr(configElement),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetConfigElementsOptions) SetSecretType(secretType string) *GetConfigElementsOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetConfigElement : Allow user to set ConfigElement
func (_options *GetConfigElementsOptions) SetConfigElement(configElement string) *GetConfigElementsOptions {
	_options.ConfigElement = core.StringPtr(configElement)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigElementsOptions) SetHeaders(param map[string]string) *GetConfigElementsOptions {
	options.Headers = param
	return options
}

// GetConfigElementsResourcesItem : GetConfigElementsResourcesItem struct
// Models which "extend" this model:
// - GetConfigElementsResourcesItemCertificateAuthoritiesConfig
// - GetConfigElementsResourcesItemDNSProvidersConfig
// - RootCertificateAuthoritiesConfig
// - IntermediateCertificateAuthoritiesConfig
// - CertificateTemplatesConfig
type GetConfigElementsResourcesItem struct {
	CertificateAuthorities []ConfigElementMetadata `json:"certificate_authorities,omitempty"`

	DNSProviders []ConfigElementMetadata `json:"dns_providers,omitempty"`

	RootCertificateAuthorities []RootCertificateAuthoritiesConfigItem `json:"root_certificate_authorities,omitempty"`

	IntermediateCertificateAuthorities []IntermediateCertificateAuthoritiesConfigItem `json:"intermediate_certificate_authorities,omitempty"`

	CertificateTemplates []CertificateTemplatesConfigItem `json:"certificate_templates,omitempty"`
}

func (*GetConfigElementsResourcesItem) isaGetConfigElementsResourcesItem() bool {
	return true
}

type GetConfigElementsResourcesItemIntf interface {
	isaGetConfigElementsResourcesItem() bool
}

// UnmarshalGetConfigElementsResourcesItem unmarshals an instance of GetConfigElementsResourcesItem from the specified map of raw messages.
func UnmarshalGetConfigElementsResourcesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetConfigElementsResourcesItem)
	err = core.UnmarshalModel(m, "certificate_authorities", &obj.CertificateAuthorities, UnmarshalConfigElementMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "dns_providers", &obj.DNSProviders, UnmarshalConfigElementMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "root_certificate_authorities", &obj.RootCertificateAuthorities, UnmarshalRootCertificateAuthoritiesConfigItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "intermediate_certificate_authorities", &obj.IntermediateCertificateAuthorities, UnmarshalIntermediateCertificateAuthoritiesConfigItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "certificate_templates", &obj.CertificateTemplates, UnmarshalCertificateTemplatesConfigItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetConfigOptions : The GetConfig options.
type GetConfigOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetConfigOptions.SecretType property.
// The secret type.
const (
	GetConfigOptionsSecretTypeIamCredentialsConst = "iam_credentials"
	GetConfigOptionsSecretTypePrivateCertConst    = "private_cert"
	GetConfigOptionsSecretTypePublicCertConst     = "public_cert"
)

// NewGetConfigOptions : Instantiate GetConfigOptions
func (*SecretsManagerV1) NewGetConfigOptions(secretType string) *GetConfigOptions {
	return &GetConfigOptions{
		SecretType: core.StringPtr(secretType),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetConfigOptions) SetSecretType(secretType string) *GetConfigOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigOptions) SetHeaders(param map[string]string) *GetConfigOptions {
	options.Headers = param
	return options
}

// GetConfigResourcesItem : GetConfigResourcesItem struct
// Models which "extend" this model:
// - PublicCertSecretEngineRootConfig
// - PrivateCertSecretEngineRootConfig
// - IamCredentialsSecretEngineRootConfig
type GetConfigResourcesItem struct {
	// The certificate authority configurations that are associated with your instance.
	CertificateAuthorities []ConfigElementMetadata `json:"certificate_authorities,omitempty"`

	// The DNS provider configurations that are associated with your instance.
	DNSProviders []ConfigElementMetadata `json:"dns_providers,omitempty"`

	// The root certificate authority configurations that are associated with your instance.
	RootCertificateAuthorities []RootCertificateAuthoritiesConfigItem `json:"root_certificate_authorities,omitempty"`

	// The intermediate certificate authority configurations that are associated with your instance.
	IntermediateCertificateAuthorities []IntermediateCertificateAuthoritiesConfigItem `json:"intermediate_certificate_authorities,omitempty"`

	// The certificate templates that are associated with your instance.
	CertificateTemplates []CertificateTemplatesConfigItem `json:"certificate_templates,omitempty"`

	// An IBM Cloud API key that can create and manage service IDs.
	//
	// The API key must be assigned the Editor platform role on the Access Groups Service and the Operator platform role on
	// the IAM Identity Service. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).
	APIKey *string `json:"api_key,omitempty"`

	// The hash value of the IBM Cloud API key that is used to create and manage service IDs.
	APIKeyHash *string `json:"api_key_hash,omitempty"`
}

func (*GetConfigResourcesItem) isaGetConfigResourcesItem() bool {
	return true
}

type GetConfigResourcesItemIntf interface {
	isaGetConfigResourcesItem() bool
}

// UnmarshalGetConfigResourcesItem unmarshals an instance of GetConfigResourcesItem from the specified map of raw messages.
func UnmarshalGetConfigResourcesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetConfigResourcesItem)
	err = core.UnmarshalModel(m, "certificate_authorities", &obj.CertificateAuthorities, UnmarshalConfigElementMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "dns_providers", &obj.DNSProviders, UnmarshalConfigElementMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "root_certificate_authorities", &obj.RootCertificateAuthorities, UnmarshalRootCertificateAuthoritiesConfigItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "intermediate_certificate_authorities", &obj.IntermediateCertificateAuthorities, UnmarshalIntermediateCertificateAuthoritiesConfigItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "certificate_templates", &obj.CertificateTemplates, UnmarshalCertificateTemplatesConfigItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_hash", &obj.APIKeyHash)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetInstanceLocks : Properties that describe the locks that are associated with an instance.
type GetInstanceLocks struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []InstanceSecretsLocks `json:"resources" validate:"required"`
}

// UnmarshalGetInstanceLocks unmarshals an instance of GetInstanceLocks from the specified map of raw messages.
func UnmarshalGetInstanceLocks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetInstanceLocks)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalInstanceSecretsLocks)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetLocksOptions : The GetLocks options.
type GetLocksOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The number of locks to retrieve. By default, list operations return the first 25 items. To retrieve a different set
	// of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 locks associated with your secret, and you want to retrieve only the first 5 locks, use
	// `..?limit=5`.
	Limit *int64 `json:"limit,omitempty"`

	// The number of locks to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 locks on your secret, and you want to retrieve locks 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// Filter locks that contain the specified string in the field "name".
	//
	// **Usage:** If you want to list only the locks that contain the string "text" in the field "name", use
	// `..?search=text`.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetLocksOptions.SecretType property.
// The secret type.
const (
	GetLocksOptionsSecretTypeArbitraryConst        = "arbitrary"
	GetLocksOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	GetLocksOptionsSecretTypeImportedCertConst     = "imported_cert"
	GetLocksOptionsSecretTypeKvConst               = "kv"
	GetLocksOptionsSecretTypePrivateCertConst      = "private_cert"
	GetLocksOptionsSecretTypePublicCertConst       = "public_cert"
	GetLocksOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewGetLocksOptions : Instantiate GetLocksOptions
func (*SecretsManagerV1) NewGetLocksOptions(secretType string, id string) *GetLocksOptions {
	return &GetLocksOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetLocksOptions) SetSecretType(secretType string) *GetLocksOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetLocksOptions) SetID(id string) *GetLocksOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetLocksOptions) SetLimit(limit int64) *GetLocksOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetLocksOptions) SetOffset(offset int64) *GetLocksOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *GetLocksOptions) SetSearch(search string) *GetLocksOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLocksOptions) SetHeaders(param map[string]string) *GetLocksOptions {
	options.Headers = param
	return options
}

// GetNotificationsRegistrationOptions : The GetNotificationsRegistration options.
type GetNotificationsRegistrationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetNotificationsRegistrationOptions : Instantiate GetNotificationsRegistrationOptions
func (*SecretsManagerV1) NewGetNotificationsRegistrationOptions() *GetNotificationsRegistrationOptions {
	return &GetNotificationsRegistrationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetNotificationsRegistrationOptions) SetHeaders(param map[string]string) *GetNotificationsRegistrationOptions {
	options.Headers = param
	return options
}

// GetNotificationsSettings : Properties that describe an existing registration with Event Notifications.
type GetNotificationsSettings struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []NotificationsSettings `json:"resources" validate:"required"`
}

// UnmarshalGetNotificationsSettings unmarshals an instance of GetNotificationsSettings from the specified map of raw messages.
func UnmarshalGetNotificationsSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetNotificationsSettings)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalNotificationsSettings)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetPolicyOptions : The GetPolicy options.
type GetPolicyOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The type of policy that is associated with the specified secret.
	Policy *string `json:"policy,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetPolicyOptions.SecretType property.
// The secret type.
const (
	GetPolicyOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	GetPolicyOptionsSecretTypePrivateCertConst      = "private_cert"
	GetPolicyOptionsSecretTypePublicCertConst       = "public_cert"
	GetPolicyOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the GetPolicyOptions.Policy property.
// The type of policy that is associated with the specified secret.
const (
	GetPolicyOptionsPolicyRotationConst = "rotation"
)

// NewGetPolicyOptions : Instantiate GetPolicyOptions
func (*SecretsManagerV1) NewGetPolicyOptions(secretType string, id string) *GetPolicyOptions {
	return &GetPolicyOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetPolicyOptions) SetSecretType(secretType string) *GetPolicyOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetPolicyOptions) SetID(id string) *GetPolicyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetPolicy : Allow user to set Policy
func (_options *GetPolicyOptions) SetPolicy(policy string) *GetPolicyOptions {
	_options.Policy = core.StringPtr(policy)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPolicyOptions) SetHeaders(param map[string]string) *GetPolicyOptions {
	options.Headers = param
	return options
}

// GetSecret : Properties that describe a secret.
type GetSecret struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretResourceIntf `json:"resources" validate:"required"`
}

// UnmarshalGetSecret unmarshals an instance of GetSecret from the specified map of raw messages.
func UnmarshalGetSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecret)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSecretGroupOptions : The GetSecretGroup options.
type GetSecretGroupOptions struct {
	// The v4 UUID that uniquely identifies the secret group.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSecretGroupOptions : Instantiate GetSecretGroupOptions
func (*SecretsManagerV1) NewGetSecretGroupOptions(id string) *GetSecretGroupOptions {
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

// GetSecretLocks : Properties that describe the lock of a secret or a secret version.
type GetSecretLocks struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretsLocks `json:"resources" validate:"required"`
}

// UnmarshalGetSecretLocks unmarshals an instance of GetSecretLocks from the specified map of raw messages.
func UnmarshalGetSecretLocks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretLocks)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretsLocks)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSecretMetadataOptions : The GetSecretMetadata options.
type GetSecretMetadataOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSecretMetadataOptions.SecretType property.
// The secret type.
const (
	GetSecretMetadataOptionsSecretTypeArbitraryConst        = "arbitrary"
	GetSecretMetadataOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	GetSecretMetadataOptionsSecretTypeImportedCertConst     = "imported_cert"
	GetSecretMetadataOptionsSecretTypeKvConst               = "kv"
	GetSecretMetadataOptionsSecretTypePrivateCertConst      = "private_cert"
	GetSecretMetadataOptionsSecretTypePublicCertConst       = "public_cert"
	GetSecretMetadataOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewGetSecretMetadataOptions : Instantiate GetSecretMetadataOptions
func (*SecretsManagerV1) NewGetSecretMetadataOptions(secretType string, id string) *GetSecretMetadataOptions {
	return &GetSecretMetadataOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetSecretMetadataOptions) SetSecretType(secretType string) *GetSecretMetadataOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
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
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSecretOptions.SecretType property.
// The secret type.
const (
	GetSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	GetSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	GetSecretOptionsSecretTypeImportedCertConst     = "imported_cert"
	GetSecretOptionsSecretTypeKvConst               = "kv"
	GetSecretOptionsSecretTypePrivateCertConst      = "private_cert"
	GetSecretOptionsSecretTypePublicCertConst       = "public_cert"
	GetSecretOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewGetSecretOptions : Instantiate GetSecretOptions
func (*SecretsManagerV1) NewGetSecretOptions(secretType string, id string) *GetSecretOptions {
	return &GetSecretOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetSecretOptions) SetSecretType(secretType string) *GetSecretOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
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

// GetSecretPolicies : GetSecretPolicies struct
// Models which "extend" this model:
// - GetSecretPolicyRotation
type GetSecretPolicies struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata,omitempty"`

	// A collection of resources.
	Resources []map[string]interface{} `json:"resources,omitempty"`
}

func (*GetSecretPolicies) isaGetSecretPolicies() bool {
	return true
}

type GetSecretPoliciesIntf interface {
	isaGetSecretPolicies() bool
}

// UnmarshalGetSecretPolicies unmarshals an instance of GetSecretPolicies from the specified map of raw messages.
func UnmarshalGetSecretPolicies(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretPolicies)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources", &obj.Resources)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSecretVersion : Properties that describe the version of a secret.
type GetSecretVersion struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretVersionIntf `json:"resources" validate:"required"`
}

// UnmarshalGetSecretVersion unmarshals an instance of GetSecretVersion from the specified map of raw messages.
func UnmarshalGetSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretVersion)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretVersion)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSecretVersionLocksOptions : The GetSecretVersionLocks options.
type GetSecretVersionLocksOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret version. You can also use `previous` to retrieve the previous
	// version.
	//
	// **Note:** To find the version ID of a secret, use the [Get secret metadata](#get-secret-metadata) method and check
	// the response details.
	VersionID *string `json:"version_id" validate:"required,ne="`

	// The number of locks to retrieve. By default, list operations return the first 25 items. To retrieve a different set
	// of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 locks associated with your secret, and you want to retrieve only the first 5 locks, use
	// `..?limit=5`.
	Limit *int64 `json:"limit,omitempty"`

	// The number of locks to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 locks on your secret, and you want to retrieve locks 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// Filter locks that contain the specified string in the field "name".
	//
	// **Usage:** If you want to list only the locks that contain the string "text" in the field "name", use
	// `..?search=text`.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSecretVersionLocksOptions.SecretType property.
// The secret type.
const (
	GetSecretVersionLocksOptionsSecretTypeArbitraryConst        = "arbitrary"
	GetSecretVersionLocksOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	GetSecretVersionLocksOptionsSecretTypeImportedCertConst     = "imported_cert"
	GetSecretVersionLocksOptionsSecretTypeKvConst               = "kv"
	GetSecretVersionLocksOptionsSecretTypePrivateCertConst      = "private_cert"
	GetSecretVersionLocksOptionsSecretTypePublicCertConst       = "public_cert"
	GetSecretVersionLocksOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewGetSecretVersionLocksOptions : Instantiate GetSecretVersionLocksOptions
func (*SecretsManagerV1) NewGetSecretVersionLocksOptions(secretType string, id string, versionID string) *GetSecretVersionLocksOptions {
	return &GetSecretVersionLocksOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		VersionID:  core.StringPtr(versionID),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetSecretVersionLocksOptions) SetSecretType(secretType string) *GetSecretVersionLocksOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetSecretVersionLocksOptions) SetID(id string) *GetSecretVersionLocksOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *GetSecretVersionLocksOptions) SetVersionID(versionID string) *GetSecretVersionLocksOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetSecretVersionLocksOptions) SetLimit(limit int64) *GetSecretVersionLocksOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetSecretVersionLocksOptions) SetOffset(offset int64) *GetSecretVersionLocksOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *GetSecretVersionLocksOptions) SetSearch(search string) *GetSecretVersionLocksOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretVersionLocksOptions) SetHeaders(param map[string]string) *GetSecretVersionLocksOptions {
	options.Headers = param
	return options
}

// GetSecretVersionMetadata : Properties that describe the version of a secret.
type GetSecretVersionMetadata struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretVersionMetadataIntf `json:"resources" validate:"required"`
}

// UnmarshalGetSecretVersionMetadata unmarshals an instance of GetSecretVersionMetadata from the specified map of raw messages.
func UnmarshalGetSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretVersionMetadata)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretVersionMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSecretVersionMetadataOptions : The GetSecretVersionMetadata options.
type GetSecretVersionMetadataOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret version. You can also use `previous` to retrieve the previous
	// version.
	//
	// **Note:** To find the version ID of a secret, use the [Get secret metadata](#get-secret-metadata) method and check
	// the response details.
	VersionID *string `json:"version_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSecretVersionMetadataOptions.SecretType property.
// The secret type.
const (
	GetSecretVersionMetadataOptionsSecretTypeArbitraryConst        = "arbitrary"
	GetSecretVersionMetadataOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	GetSecretVersionMetadataOptionsSecretTypeImportedCertConst     = "imported_cert"
	GetSecretVersionMetadataOptionsSecretTypeKvConst               = "kv"
	GetSecretVersionMetadataOptionsSecretTypePrivateCertConst      = "private_cert"
	GetSecretVersionMetadataOptionsSecretTypePublicCertConst       = "public_cert"
	GetSecretVersionMetadataOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewGetSecretVersionMetadataOptions : Instantiate GetSecretVersionMetadataOptions
func (*SecretsManagerV1) NewGetSecretVersionMetadataOptions(secretType string, id string, versionID string) *GetSecretVersionMetadataOptions {
	return &GetSecretVersionMetadataOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		VersionID:  core.StringPtr(versionID),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetSecretVersionMetadataOptions) SetSecretType(secretType string) *GetSecretVersionMetadataOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetSecretVersionMetadataOptions) SetID(id string) *GetSecretVersionMetadataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *GetSecretVersionMetadataOptions) SetVersionID(versionID string) *GetSecretVersionMetadataOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretVersionMetadataOptions) SetHeaders(param map[string]string) *GetSecretVersionMetadataOptions {
	options.Headers = param
	return options
}

// GetSecretVersionOptions : The GetSecretVersion options.
type GetSecretVersionOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret version. You can also use `previous` to retrieve the previous
	// version.
	//
	// **Note:** To find the version ID of a secret, use the [Get secret metadata](#get-secret-metadata) method and check
	// the response details.
	VersionID *string `json:"version_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSecretVersionOptions.SecretType property.
// The secret type.
const (
	GetSecretVersionOptionsSecretTypeArbitraryConst        = "arbitrary"
	GetSecretVersionOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	GetSecretVersionOptionsSecretTypeImportedCertConst     = "imported_cert"
	GetSecretVersionOptionsSecretTypeKvConst               = "kv"
	GetSecretVersionOptionsSecretTypePrivateCertConst      = "private_cert"
	GetSecretVersionOptionsSecretTypePublicCertConst       = "public_cert"
	GetSecretVersionOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewGetSecretVersionOptions : Instantiate GetSecretVersionOptions
func (*SecretsManagerV1) NewGetSecretVersionOptions(secretType string, id string, versionID string) *GetSecretVersionOptions {
	return &GetSecretVersionOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		VersionID:  core.StringPtr(versionID),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *GetSecretVersionOptions) SetSecretType(secretType string) *GetSecretVersionOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetSecretVersionOptions) SetID(id string) *GetSecretVersionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *GetSecretVersionOptions) SetVersionID(versionID string) *GetSecretVersionOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretVersionOptions) SetHeaders(param map[string]string) *GetSecretVersionOptions {
	options.Headers = param
	return options
}

// GetSingleConfigElement : Properties that describe a configuration.
type GetSingleConfigElement struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []ConfigElementDef `json:"resources" validate:"required"`
}

// UnmarshalGetSingleConfigElement unmarshals an instance of GetSingleConfigElement from the specified map of raw messages.
func UnmarshalGetSingleConfigElement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSingleConfigElement)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalConfigElementDef)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstanceSecretsLocks : Properties that describe the locks that are associated with an instance.
type InstanceSecretsLocks struct {
	// The unique ID of the secret.
	SecretID *string `json:"secret_id,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// A collection of locks that are attached to a secret version.
	Versions []SecretLockVersion `json:"versions,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// Constants associated with the InstanceSecretsLocks.SecretType property.
// The secret type.
const (
	InstanceSecretsLocksSecretTypeArbitraryConst        = "arbitrary"
	InstanceSecretsLocksSecretTypeIamCredentialsConst   = "iam_credentials"
	InstanceSecretsLocksSecretTypeImportedCertConst     = "imported_cert"
	InstanceSecretsLocksSecretTypeKvConst               = "kv"
	InstanceSecretsLocksSecretTypePrivateCertConst      = "private_cert"
	InstanceSecretsLocksSecretTypePublicCertConst       = "public_cert"
	InstanceSecretsLocksSecretTypeUsernamePasswordConst = "username_password"
)

// SetProperty allows the user to set an arbitrary property on an instance of InstanceSecretsLocks
func (o *InstanceSecretsLocks) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of InstanceSecretsLocks
func (o *InstanceSecretsLocks) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of InstanceSecretsLocks
func (o *InstanceSecretsLocks) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of InstanceSecretsLocks
func (o *InstanceSecretsLocks) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of InstanceSecretsLocks
func (o *InstanceSecretsLocks) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.SecretID != nil {
		m["secret_id"] = o.SecretID
	}
	if o.SecretGroupID != nil {
		m["secret_group_id"] = o.SecretGroupID
	}
	if o.SecretType != nil {
		m["secret_type"] = o.SecretType
	}
	if o.Versions != nil {
		m["versions"] = o.Versions
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalInstanceSecretsLocks unmarshals an instance of InstanceSecretsLocks from the specified map of raw messages.
func UnmarshalInstanceSecretsLocks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstanceSecretsLocks)
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	delete(m, "secret_id")
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	delete(m, "secret_group_id")
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	delete(m, "secret_type")
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalSecretLockVersion)
	if err != nil {
		return
	}
	delete(m, "versions")
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

// IntermediateCertificateAuthoritiesConfigItem : Intermediate certificate authorities configuration.
type IntermediateCertificateAuthoritiesConfigItem struct {
	// The human-readable name to assign to your configuration.
	Name *string `json:"name" validate:"required"`

	// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
	Type *string `json:"type" validate:"required"`

	// Intermediate certificate authority configuration.
	Config *IntermediateCertificateAuthorityConfig `json:"config,omitempty"`
}

// Constants associated with the IntermediateCertificateAuthoritiesConfigItem.Type property.
// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
const (
	IntermediateCertificateAuthoritiesConfigItemTypeCertificateTemplateConst              = "certificate_template"
	IntermediateCertificateAuthoritiesConfigItemTypeCisConst                              = "cis"
	IntermediateCertificateAuthoritiesConfigItemTypeClassicInfrastructureConst            = "classic_infrastructure"
	IntermediateCertificateAuthoritiesConfigItemTypeIntermediateCertificateAuthorityConst = "intermediate_certificate_authority"
	IntermediateCertificateAuthoritiesConfigItemTypeLetsencryptConst                      = "letsencrypt"
	IntermediateCertificateAuthoritiesConfigItemTypeLetsencryptStageConst                 = "letsencrypt-stage"
	IntermediateCertificateAuthoritiesConfigItemTypeRootCertificateAuthorityConst         = "root_certificate_authority"
)

// UnmarshalIntermediateCertificateAuthoritiesConfigItem unmarshals an instance of IntermediateCertificateAuthoritiesConfigItem from the specified map of raw messages.
func UnmarshalIntermediateCertificateAuthoritiesConfigItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IntermediateCertificateAuthoritiesConfigItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalIntermediateCertificateAuthorityConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IssuanceInfo : Issuance information that is associated with your certificate.
type IssuanceInfo struct {
	// The date the certificate was ordered. The date format follows RFC 3339.
	OrderedOn *strfmt.DateTime `json:"ordered_on,omitempty"`

	// A code that identifies an issuance error.
	//
	// This field, along with `error_message`, is returned when Secrets Manager successfully processes your request, but a
	// certificate is unable to be issued by the certificate authority.
	ErrorCode *string `json:"error_code,omitempty"`

	// A human-readable message that provides details about the issuance error.
	ErrorMessage *string `json:"error_message,omitempty"`

	// Indicates whether the issued certificate is bundled with intermediate certificates.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// Indicates whether the issued certificate is configured with an automatic rotation policy.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The name that was assigned to the certificate authority configuration.
	Ca *string `json:"ca,omitempty"`

	// The name that was assigned to the DNS provider configuration.
	DNS *string `json:"dns,omitempty"`

	// The set of challenges, will be returned only when ordering public certificate using manual DNS configuration.
	Challenges []ChallengeResource `json:"challenges,omitempty"`

	// The date a user called "validate dns challenges" for "manual" DNS provider. The date format follows RFC 3339.
	DNSChallengeValidationTime *strfmt.DateTime `json:"dns_challenge_validation_time,omitempty"`
}

// UnmarshalIssuanceInfo unmarshals an instance of IssuanceInfo from the specified map of raw messages.
func UnmarshalIssuanceInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IssuanceInfo)
	err = core.UnmarshalPrimitive(m, "ordered_on", &obj.OrderedOn)
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
	err = core.UnmarshalPrimitive(m, "bundle_certs", &obj.BundleCerts)
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
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ca", &obj.Ca)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dns", &obj.DNS)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "challenges", &obj.Challenges, UnmarshalChallengeResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dns_challenge_validation_time", &obj.DNSChallengeValidationTime)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAllSecretsOptions : The ListAllSecrets options.
type ListAllSecretsOptions struct {
	// The number of secrets to retrieve. By default, list operations return the first 200 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 secrets in your instance, and you want to retrieve only the first 5 secrets, use
	// `../secrets/{secret_type}?limit=5`.
	Limit *int64 `json:"limit,omitempty"`

	// The number of secrets to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 secrets in your instance, and you want to retrieve secrets 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// Filter secrets that contain the specified string. The fields that are searched include: id, name, description,
	// labels, secret_type.
	//
	// **Usage:** If you want to list only the secrets that contain the string "text", use
	// `../secrets/{secret_type}?search=text`.
	Search *string `json:"search,omitempty"`

	// Sort a list of secrets by the specified field.
	//
	// **Usage:** To sort a list of secrets by their creation date, use
	// `../secrets/{secret_type}?sort_by=creation_date`.
	SortBy *string `json:"sort_by,omitempty"`

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

// Constants associated with the ListAllSecretsOptions.SortBy property.
// Sort a list of secrets by the specified field.
//
// **Usage:** To sort a list of secrets by their creation date, use
// `../secrets/{secret_type}?sort_by=creation_date`.
const (
	ListAllSecretsOptionsSortByCreationDateConst   = "creation_date"
	ListAllSecretsOptionsSortByExpirationDateConst = "expiration_date"
	ListAllSecretsOptionsSortByIDConst             = "id"
	ListAllSecretsOptionsSortByNameConst           = "name"
	ListAllSecretsOptionsSortBySecretTypeConst     = "secret_type"
)

// NewListAllSecretsOptions : Instantiate ListAllSecretsOptions
func (*SecretsManagerV1) NewListAllSecretsOptions() *ListAllSecretsOptions {
	return &ListAllSecretsOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *ListAllSecretsOptions) SetLimit(limit int64) *ListAllSecretsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListAllSecretsOptions) SetOffset(offset int64) *ListAllSecretsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListAllSecretsOptions) SetSearch(search string) *ListAllSecretsOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetSortBy : Allow user to set SortBy
func (_options *ListAllSecretsOptions) SetSortBy(sortBy string) *ListAllSecretsOptions {
	_options.SortBy = core.StringPtr(sortBy)
	return _options
}

// SetGroups : Allow user to set Groups
func (_options *ListAllSecretsOptions) SetGroups(groups []string) *ListAllSecretsOptions {
	_options.Groups = groups
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllSecretsOptions) SetHeaders(param map[string]string) *ListAllSecretsOptions {
	options.Headers = param
	return options
}

// ListInstanceSecretsLocksOptions : The ListInstanceSecretsLocks options.
type ListInstanceSecretsLocksOptions struct {
	// The number of secrets with associated locks to retrieve. By default, list operations return the first 25 items. To
	// retrieve a different set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 secrets in your instance, and you want to retrieve only the first 5, use
	// `..?limit=5`.
	Limit *int64 `json:"limit,omitempty"`

	// The number of secrets to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 secrets in your instance, and you want to retrieve secrets 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

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

// NewListInstanceSecretsLocksOptions : Instantiate ListInstanceSecretsLocksOptions
func (*SecretsManagerV1) NewListInstanceSecretsLocksOptions() *ListInstanceSecretsLocksOptions {
	return &ListInstanceSecretsLocksOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *ListInstanceSecretsLocksOptions) SetLimit(limit int64) *ListInstanceSecretsLocksOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListInstanceSecretsLocksOptions) SetOffset(offset int64) *ListInstanceSecretsLocksOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListInstanceSecretsLocksOptions) SetSearch(search string) *ListInstanceSecretsLocksOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetGroups : Allow user to set Groups
func (_options *ListInstanceSecretsLocksOptions) SetGroups(groups []string) *ListInstanceSecretsLocksOptions {
	_options.Groups = groups
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListInstanceSecretsLocksOptions) SetHeaders(param map[string]string) *ListInstanceSecretsLocksOptions {
	options.Headers = param
	return options
}

// ListSecretGroupsOptions : The ListSecretGroups options.
type ListSecretGroupsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSecretGroupsOptions : Instantiate ListSecretGroupsOptions
func (*SecretsManagerV1) NewListSecretGroupsOptions() *ListSecretGroupsOptions {
	return &ListSecretGroupsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretGroupsOptions) SetHeaders(param map[string]string) *ListSecretGroupsOptions {
	options.Headers = param
	return options
}

// ListSecretLocks : Properties that describe the locks of a secret or a secret version.
type ListSecretLocks struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretLockData `json:"resources" validate:"required"`
}

// UnmarshalListSecretLocks unmarshals an instance of ListSecretLocks from the specified map of raw messages.
func UnmarshalListSecretLocks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListSecretLocks)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretLockData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListSecretVersions : Properties that describe a list of versions of a secret.
type ListSecretVersions struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretVersionInfoIntf `json:"resources,omitempty"`
}

// UnmarshalListSecretVersions unmarshals an instance of ListSecretVersions from the specified map of raw messages.
func UnmarshalListSecretVersions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListSecretVersions)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretVersionInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListSecretVersionsOptions : The ListSecretVersions options.
type ListSecretVersionsOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListSecretVersionsOptions.SecretType property.
// The secret type.
const (
	ListSecretVersionsOptionsSecretTypeArbitraryConst        = "arbitrary"
	ListSecretVersionsOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	ListSecretVersionsOptionsSecretTypeImportedCertConst     = "imported_cert"
	ListSecretVersionsOptionsSecretTypeKvConst               = "kv"
	ListSecretVersionsOptionsSecretTypePrivateCertConst      = "private_cert"
	ListSecretVersionsOptionsSecretTypePublicCertConst       = "public_cert"
	ListSecretVersionsOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewListSecretVersionsOptions : Instantiate ListSecretVersionsOptions
func (*SecretsManagerV1) NewListSecretVersionsOptions(secretType string, id string) *ListSecretVersionsOptions {
	return &ListSecretVersionsOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *ListSecretVersionsOptions) SetSecretType(secretType string) *ListSecretVersionsOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *ListSecretVersionsOptions) SetID(id string) *ListSecretVersionsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretVersionsOptions) SetHeaders(param map[string]string) *ListSecretVersionsOptions {
	options.Headers = param
	return options
}

// ListSecrets : Properties that describe a list of secrets.
type ListSecrets struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretResourceIntf `json:"resources,omitempty"`
}

// UnmarshalListSecrets unmarshals an instance of ListSecrets from the specified map of raw messages.
func UnmarshalListSecrets(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListSecrets)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListSecretsOptions : The ListSecrets options.
type ListSecretsOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The number of secrets to retrieve. By default, list operations return the first 200 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 secrets in your instance, and you want to retrieve only the first 5 secrets, use
	// `../secrets/{secret_type}?limit=5`.
	Limit *int64 `json:"limit,omitempty"`

	// The number of secrets to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 secrets in your instance, and you want to retrieve secrets 26 through 50, use
	// `..?offset=25&limit=25`.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListSecretsOptions.SecretType property.
// The secret type.
const (
	ListSecretsOptionsSecretTypeArbitraryConst        = "arbitrary"
	ListSecretsOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	ListSecretsOptionsSecretTypeImportedCertConst     = "imported_cert"
	ListSecretsOptionsSecretTypeKvConst               = "kv"
	ListSecretsOptionsSecretTypePrivateCertConst      = "private_cert"
	ListSecretsOptionsSecretTypePublicCertConst       = "public_cert"
	ListSecretsOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewListSecretsOptions : Instantiate ListSecretsOptions
func (*SecretsManagerV1) NewListSecretsOptions(secretType string) *ListSecretsOptions {
	return &ListSecretsOptions{
		SecretType: core.StringPtr(secretType),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *ListSecretsOptions) SetSecretType(secretType string) *ListSecretsOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSecretsOptions) SetLimit(limit int64) *ListSecretsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListSecretsOptions) SetOffset(offset int64) *ListSecretsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretsOptions) SetHeaders(param map[string]string) *ListSecretsOptions {
	options.Headers = param
	return options
}

// LockSecretBodyLocksItem : LockSecretBodyLocksItem struct
type LockSecretBodyLocksItem struct {
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

// NewLockSecretBodyLocksItem : Instantiate LockSecretBodyLocksItem (Generic Model Constructor)
func (*SecretsManagerV1) NewLockSecretBodyLocksItem(name string) (_model *LockSecretBodyLocksItem, err error) {
	_model = &LockSecretBodyLocksItem{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalLockSecretBodyLocksItem unmarshals an instance of LockSecretBodyLocksItem from the specified map of raw messages.
func UnmarshalLockSecretBodyLocksItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LockSecretBodyLocksItem)
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

// LockSecretOptions : The LockSecret options.
type LockSecretOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The lock data to be attached to a secret version.
	Locks []LockSecretBodyLocksItem `json:"locks,omitempty"`

	// An optional lock mode. At lock creation, you can set one of the following modes to clear any matching locks on a
	// secret version. Note: When you are locking the `previous` version, the mode parameter is ignored.
	//
	// - `exclusive`: Removes any other locks with matching names if they are found in the previous version of the secret.
	// - `exclusive_delete`: Same as `exclusive`, but also permanently deletes the data of the previous secret version if
	// it doesn't have any locks.
	Mode *string `json:"mode,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the LockSecretOptions.SecretType property.
// The secret type.
const (
	LockSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	LockSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	LockSecretOptionsSecretTypeImportedCertConst     = "imported_cert"
	LockSecretOptionsSecretTypeKvConst               = "kv"
	LockSecretOptionsSecretTypePrivateCertConst      = "private_cert"
	LockSecretOptionsSecretTypePublicCertConst       = "public_cert"
	LockSecretOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the LockSecretOptions.Mode property.
// An optional lock mode. At lock creation, you can set one of the following modes to clear any matching locks on a
// secret version. Note: When you are locking the `previous` version, the mode parameter is ignored.
//
// - `exclusive`: Removes any other locks with matching names if they are found in the previous version of the secret.
// - `exclusive_delete`: Same as `exclusive`, but also permanently deletes the data of the previous secret version if it
// doesn't have any locks.
const (
	LockSecretOptionsModeExclusiveConst       = "exclusive"
	LockSecretOptionsModeExclusiveDeleteConst = "exclusive_delete"
)

// NewLockSecretOptions : Instantiate LockSecretOptions
func (*SecretsManagerV1) NewLockSecretOptions(secretType string, id string) *LockSecretOptions {
	return &LockSecretOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *LockSecretOptions) SetSecretType(secretType string) *LockSecretOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *LockSecretOptions) SetID(id string) *LockSecretOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLocks : Allow user to set Locks
func (_options *LockSecretOptions) SetLocks(locks []LockSecretBodyLocksItem) *LockSecretOptions {
	_options.Locks = locks
	return _options
}

// SetMode : Allow user to set Mode
func (_options *LockSecretOptions) SetMode(mode string) *LockSecretOptions {
	_options.Mode = core.StringPtr(mode)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *LockSecretOptions) SetHeaders(param map[string]string) *LockSecretOptions {
	options.Headers = param
	return options
}

// LockSecretVersionOptions : The LockSecretVersion options.
type LockSecretVersionOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret version. You can also use `previous` to retrieve the previous
	// version.
	//
	// **Note:** To find the version ID of a secret, use the [Get secret metadata](#get-secret-metadata) method and check
	// the response details.
	VersionID *string `json:"version_id" validate:"required,ne="`

	// The lock data to be attached to a secret version.
	Locks []LockSecretBodyLocksItem `json:"locks,omitempty"`

	// An optional lock mode. At lock creation, you can set one of the following modes to clear any matching locks on a
	// secret version. Note: When you are locking the `previous` version, the mode parameter is ignored.
	//
	// - `exclusive`: Removes any other locks with matching names if they are found in the previous version of the secret.
	// - `exclusive_delete`: Same as `exclusive`, but also permanently deletes the data of the previous secret version if
	// it doesn't have any locks.
	Mode *string `json:"mode,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the LockSecretVersionOptions.SecretType property.
// The secret type.
const (
	LockSecretVersionOptionsSecretTypeArbitraryConst        = "arbitrary"
	LockSecretVersionOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	LockSecretVersionOptionsSecretTypeImportedCertConst     = "imported_cert"
	LockSecretVersionOptionsSecretTypeKvConst               = "kv"
	LockSecretVersionOptionsSecretTypePrivateCertConst      = "private_cert"
	LockSecretVersionOptionsSecretTypePublicCertConst       = "public_cert"
	LockSecretVersionOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the LockSecretVersionOptions.Mode property.
// An optional lock mode. At lock creation, you can set one of the following modes to clear any matching locks on a
// secret version. Note: When you are locking the `previous` version, the mode parameter is ignored.
//
// - `exclusive`: Removes any other locks with matching names if they are found in the previous version of the secret.
// - `exclusive_delete`: Same as `exclusive`, but also permanently deletes the data of the previous secret version if it
// doesn't have any locks.
const (
	LockSecretVersionOptionsModeExclusiveConst       = "exclusive"
	LockSecretVersionOptionsModeExclusiveDeleteConst = "exclusive_delete"
)

// NewLockSecretVersionOptions : Instantiate LockSecretVersionOptions
func (*SecretsManagerV1) NewLockSecretVersionOptions(secretType string, id string, versionID string) *LockSecretVersionOptions {
	return &LockSecretVersionOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		VersionID:  core.StringPtr(versionID),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *LockSecretVersionOptions) SetSecretType(secretType string) *LockSecretVersionOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *LockSecretVersionOptions) SetID(id string) *LockSecretVersionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *LockSecretVersionOptions) SetVersionID(versionID string) *LockSecretVersionOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetLocks : Allow user to set Locks
func (_options *LockSecretVersionOptions) SetLocks(locks []LockSecretBodyLocksItem) *LockSecretVersionOptions {
	_options.Locks = locks
	return _options
}

// SetMode : Allow user to set Mode
func (_options *LockSecretVersionOptions) SetMode(mode string) *LockSecretVersionOptions {
	_options.Mode = core.StringPtr(mode)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *LockSecretVersionOptions) SetHeaders(param map[string]string) *LockSecretVersionOptions {
	options.Headers = param
	return options
}

// NotificationsSettings : The Event Notifications details.
type NotificationsSettings struct {
	// The Cloud Resource Name (CRN) of the connected Event Notifications instance.
	EventNotificationsInstanceCRN *string `json:"event_notifications_instance_crn" validate:"required"`
}

// UnmarshalNotificationsSettings unmarshals an instance of NotificationsSettings from the specified map of raw messages.
func UnmarshalNotificationsSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationsSettings)
	err = core.UnmarshalPrimitive(m, "event_notifications_instance_crn", &obj.EventNotificationsInstanceCRN)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PutConfigOptions : The PutConfig options.
type PutConfigOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// Properties to update for a secrets engine.
	EngineConfig EngineConfigIntf `json:"EngineConfig" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutConfigOptions.SecretType property.
// The secret type.
const (
	PutConfigOptionsSecretTypeIamCredentialsConst = "iam_credentials"
)

// NewPutConfigOptions : Instantiate PutConfigOptions
func (*SecretsManagerV1) NewPutConfigOptions(secretType string, engineConfig EngineConfigIntf) *PutConfigOptions {
	return &PutConfigOptions{
		SecretType:   core.StringPtr(secretType),
		EngineConfig: engineConfig,
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *PutConfigOptions) SetSecretType(secretType string) *PutConfigOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetEngineConfig : Allow user to set EngineConfig
func (_options *PutConfigOptions) SetEngineConfig(engineConfig EngineConfigIntf) *PutConfigOptions {
	_options.EngineConfig = engineConfig
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutConfigOptions) SetHeaders(param map[string]string) *PutConfigOptions {
	options.Headers = param
	return options
}

// PutPolicyOptions : The PutPolicy options.
type PutPolicyOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretPolicyRotation `json:"resources" validate:"required"`

	// The type of policy that is associated with the specified secret.
	Policy *string `json:"policy,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutPolicyOptions.SecretType property.
// The secret type.
const (
	PutPolicyOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	PutPolicyOptionsSecretTypePrivateCertConst      = "private_cert"
	PutPolicyOptionsSecretTypePublicCertConst       = "public_cert"
	PutPolicyOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the PutPolicyOptions.Policy property.
// The type of policy that is associated with the specified secret.
const (
	PutPolicyOptionsPolicyRotationConst = "rotation"
)

// NewPutPolicyOptions : Instantiate PutPolicyOptions
func (*SecretsManagerV1) NewPutPolicyOptions(secretType string, id string, metadata *CollectionMetadata, resources []SecretPolicyRotation) *PutPolicyOptions {
	return &PutPolicyOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		Metadata:   metadata,
		Resources:  resources,
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *PutPolicyOptions) SetSecretType(secretType string) *PutPolicyOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *PutPolicyOptions) SetID(id string) *PutPolicyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *PutPolicyOptions) SetMetadata(metadata *CollectionMetadata) *PutPolicyOptions {
	_options.Metadata = metadata
	return _options
}

// SetResources : Allow user to set Resources
func (_options *PutPolicyOptions) SetResources(resources []SecretPolicyRotation) *PutPolicyOptions {
	_options.Resources = resources
	return _options
}

// SetPolicy : Allow user to set Policy
func (_options *PutPolicyOptions) SetPolicy(policy string) *PutPolicyOptions {
	_options.Policy = core.StringPtr(policy)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutPolicyOptions) SetHeaders(param map[string]string) *PutPolicyOptions {
	options.Headers = param
	return options
}

// RootCertificateAuthoritiesConfigItem : Root certificate authorities configuration.
type RootCertificateAuthoritiesConfigItem struct {
	// The human-readable name to assign to your configuration.
	Name *string `json:"name" validate:"required"`

	// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
	Type *string `json:"type" validate:"required"`

	// Root certificate authority configuration.
	Config *RootCertificateAuthorityConfig `json:"config,omitempty"`
}

// Constants associated with the RootCertificateAuthoritiesConfigItem.Type property.
// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
const (
	RootCertificateAuthoritiesConfigItemTypeCertificateTemplateConst              = "certificate_template"
	RootCertificateAuthoritiesConfigItemTypeCisConst                              = "cis"
	RootCertificateAuthoritiesConfigItemTypeClassicInfrastructureConst            = "classic_infrastructure"
	RootCertificateAuthoritiesConfigItemTypeIntermediateCertificateAuthorityConst = "intermediate_certificate_authority"
	RootCertificateAuthoritiesConfigItemTypeLetsencryptConst                      = "letsencrypt"
	RootCertificateAuthoritiesConfigItemTypeLetsencryptStageConst                 = "letsencrypt-stage"
	RootCertificateAuthoritiesConfigItemTypeRootCertificateAuthorityConst         = "root_certificate_authority"
)

// UnmarshalRootCertificateAuthoritiesConfigItem unmarshals an instance of RootCertificateAuthoritiesConfigItem from the specified map of raw messages.
func UnmarshalRootCertificateAuthoritiesConfigItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RootCertificateAuthoritiesConfigItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalRootCertificateAuthorityConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Rotation : Rotation struct
type Rotation struct {
	// Determines whether Secrets Manager rotates your certificate automatically.
	//
	// For public certificates, if `auto_rotate` is set to `true` the service reorders your certificate 31 days before it
	// expires. For private certificates, the certificate is rotated according to the time interval specified in the
	// `interval` and `unit` fields.
	//
	// To access the previous version of the certificate, you can use the
	// [Get a version of a secret](#get-secret-version) method.
	AutoRotate *bool `json:"auto_rotate,omitempty"`

	// Determines whether Secrets Manager rotates the private key for your certificate automatically.
	//
	// If set to `true`, the service generates and stores a new private key for your rotated certificate.
	//
	// **Note:** Use this field only for public certificates. It is ignored for private certificates.
	RotateKeys *bool `json:"rotate_keys,omitempty"`

	// Used together with the `unit` field to specify the rotation interval. The minimum interval is one day, and the
	// maximum interval is 3 years (1095 days). Required in case `auto_rotate` is set to `true`.
	//
	// **Note:** Use this field only for private certificates. It is ignored for public certificates.
	Interval *int64 `json:"interval,omitempty"`

	// The time unit of the rotation interval.
	//
	// **Note:** Use this field only for private certificates. It is ignored for public certificates.
	Unit *string `json:"unit,omitempty"`
}

// Constants associated with the Rotation.Unit property.
// The time unit of the rotation interval.
//
// **Note:** Use this field only for private certificates. It is ignored for public certificates.
const (
	RotationUnitDayConst   = "day"
	RotationUnitMonthConst = "month"
)

// UnmarshalRotation unmarshals an instance of Rotation from the specified map of raw messages.
func UnmarshalRotation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rotation)
	err = core.UnmarshalPrimitive(m, "auto_rotate", &obj.AutoRotate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rotate_keys", &obj.RotateKeys)
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

// SecretAction : SecretAction struct
// Models which "extend" this model:
// - RotateArbitrarySecretBody
// - RotatePublicCertBody
// - RotateUsernamePasswordSecretBody
// - RotateCertificateBody
// - RotatePrivateCertBody
// - RotatePrivateCertBodyWithCsr
// - RotatePrivateCertBodyWithVersionCustomMetadata
// - RestoreIamCredentialsSecretBody
// - DeleteCredentialsForIamCredentialsSecret
// - RotateKvSecretBody
type SecretAction struct {
	// The new secret data to assign to an `arbitrary` secret.
	Payload interface{} `json:"payload,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Determine whether keys must be rotated.
	RotateKeys *bool `json:"rotate_keys,omitempty"`

	// The new password to assign to a `username_password` secret.
	Password *string `json:"password,omitempty"`

	// The new data to associate with the certificate.
	Certificate *string `json:"certificate,omitempty"`

	// The new private key to associate with the certificate.
	PrivateKey *string `json:"private_key,omitempty"`

	// The new intermediate certificate to associate with the certificate.
	Intermediate *string `json:"intermediate,omitempty"`

	// The certificate signing request. If you provide a CSR, it is used for auto rotation and manual rotation requests
	// that do not include a CSR. If you don't include the CSR, the certificate is generated with the last CSR that you
	// provided to create the private certificate, or on a previous request to rotate the certificate. If no CSR was
	// provided in the past, the certificate is generated with a CSR that is created internally.
	Csr *string `json:"csr,omitempty"`

	// The ID of the target version or the alias `previous`.
	VersionID *string `json:"version_id,omitempty"`

	// The ID of the API key that you want to delete. If the secret was created with a static service ID, only the API key
	// is deleted. Otherwise, the service ID is deleted together with its API key.
	APIKeyID *string `json:"api_key_id,omitempty"`

	// The service ID that you want to delete. This property can be used instead of the `api_key_id` field, but only for
	// secrets that were created with a service ID that was generated by Secrets Manager.
	//
	// **Deprecated.** Use the `api_key_id` field instead.
	// Deprecated: this field is deprecated and may be removed in a future release.
	ServiceID *string `json:"service_id,omitempty"`
}

func (*SecretAction) isaSecretAction() bool {
	return true
}

type SecretActionIntf interface {
	isaSecretAction() bool
}

// UnmarshalSecretAction unmarshals an instance of SecretAction from the specified map of raw messages.
func UnmarshalSecretAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretAction)
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
	err = core.UnmarshalPrimitive(m, "rotate_keys", &obj.RotateKeys)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
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
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.APIKeyID)
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

// SecretGroupDef : Properties that describe a secret group.
type SecretGroupDef struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretGroupResource `json:"resources" validate:"required"`
}

// NewSecretGroupDef : Instantiate SecretGroupDef (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretGroupDef(metadata *CollectionMetadata, resources []SecretGroupResource) (_model *SecretGroupDef, err error) {
	_model = &SecretGroupDef{
		Metadata:  metadata,
		Resources: resources,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSecretGroupDef unmarshals an instance of SecretGroupDef from the specified map of raw messages.
func UnmarshalSecretGroupDef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretGroupDef)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretGroupResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretGroupMetadataUpdatable : Metadata properties to update for a secret group.
type SecretGroupMetadataUpdatable struct {
	// A human-readable name to assign to your secret group.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret group.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret group.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`
}

// UnmarshalSecretGroupMetadataUpdatable unmarshals an instance of SecretGroupMetadataUpdatable from the specified map of raw messages.
func UnmarshalSecretGroupMetadataUpdatable(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretGroupMetadataUpdatable)
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

// SecretGroupResource : Properties that describe a secret group.
type SecretGroupResource struct {
	// The v4 UUID that uniquely identifies the secret group.
	ID *string `json:"id,omitempty"`

	// A human-readable name to assign to your secret group.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for your secret group.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret group.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret
	// group.
	Description *string `json:"description,omitempty"`

	// The date the secret group was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// Updates when the metadata of the secret group is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The MIME type that represents the secret group.
	Type *string `json:"type,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of SecretGroupResource
func (o *SecretGroupResource) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretGroupResource
func (o *SecretGroupResource) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretGroupResource
func (o *SecretGroupResource) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretGroupResource
func (o *SecretGroupResource) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretGroupResource
func (o *SecretGroupResource) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.ID != nil {
		m["id"] = o.ID
	}
	if o.Name != nil {
		m["name"] = o.Name
	}
	if o.Description != nil {
		m["description"] = o.Description
	}
	if o.CreationDate != nil {
		m["creation_date"] = o.CreationDate
	}
	if o.LastUpdateDate != nil {
		m["last_update_date"] = o.LastUpdateDate
	}
	if o.Type != nil {
		m["type"] = o.Type
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalSecretGroupResource unmarshals an instance of SecretGroupResource from the specified map of raw messages.
func UnmarshalSecretGroupResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretGroupResource)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	delete(m, "id")
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	delete(m, "name")
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	delete(m, "description")
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	delete(m, "creation_date")
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	delete(m, "last_update_date")
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	delete(m, "type")
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

// SecretLockData : Properties that describe a lock.
type SecretLockData struct {
	// A human-readable name to assign to the secret lock.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a name for the secret lock.
	Name *string `json:"name,omitempty"`

	// An extended description of the secret lock.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for the secret
	// lock.
	Description *string `json:"description,omitempty"`

	// The date the secret lock was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret lock.
	CreatedBy *string `json:"created_by,omitempty"`

	// The information that is associated with a lock, such as resources CRNs to be used by automation.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// The v4 UUID that uniquely identifies the secret version.
	SecretVersionID *string `json:"secret_version_id,omitempty"`

	// The v4 UUID that uniquely identifies the secret.
	SecretID *string `json:"secret_id,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// A representation for the 2 last secret versions. Could be "current" for version (n) or "previous" for version (n-1).
	SecretVersionAlias *string `json:"secret_version_alias,omitempty"`
}

// UnmarshalSecretLockData unmarshals an instance of SecretLockData from the specified map of raw messages.
func UnmarshalSecretLockData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretLockData)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attributes", &obj.Attributes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_version_id", &obj.SecretVersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
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

// SecretLockVersion : Properties that describe the secret locks.
type SecretLockVersion struct {
	// The v4 UUID that uniquely identifies the lock.
	ID *string `json:"id,omitempty"`

	// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
	// for version `n-1`.
	Alias *string `json:"alias,omitempty"`

	// The names of all locks that are associated with this secret.
	Locks []string `json:"locks,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// Constants associated with the SecretLockVersion.Alias property.
// A human-readable alias that describes the secret version. 'Current' is used for version `n` and 'previous' is used
// for version `n-1`.
const (
	SecretLockVersionAliasCurrentConst  = "current"
	SecretLockVersionAliasPreviousConst = "previous"
)

// SetProperty allows the user to set an arbitrary property on an instance of SecretLockVersion
func (o *SecretLockVersion) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretLockVersion
func (o *SecretLockVersion) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretLockVersion
func (o *SecretLockVersion) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretLockVersion
func (o *SecretLockVersion) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretLockVersion
func (o *SecretLockVersion) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.ID != nil {
		m["id"] = o.ID
	}
	if o.Alias != nil {
		m["alias"] = o.Alias
	}
	if o.Locks != nil {
		m["locks"] = o.Locks
	}
	if o.PayloadAvailable != nil {
		m["payload_available"] = o.PayloadAvailable
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalSecretLockVersion unmarshals an instance of SecretLockVersion from the specified map of raw messages.
func UnmarshalSecretLockVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretLockVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	delete(m, "id")
	err = core.UnmarshalPrimitive(m, "alias", &obj.Alias)
	if err != nil {
		return
	}
	delete(m, "alias")
	err = core.UnmarshalPrimitive(m, "locks", &obj.Locks)
	if err != nil {
		return
	}
	delete(m, "locks")
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	delete(m, "payload_available")
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

// SecretMetadata : SecretMetadata struct
// Models which "extend" this model:
// - ArbitrarySecretMetadata
// - UsernamePasswordSecretMetadata
// - IamCredentialsSecretMetadata
// - CertificateSecretMetadata
// - PublicCertificateSecretMetadata
// - PrivateCertificateSecretMetadata
// - KvSecretMetadata
type SecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be in the range 2 - 30 characters, including spaces. Special characters
	// that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The time-to-live (TTL) or lease duration that is assigned to the secret. For `iam_credentials` secrets, the TTL
	// defines for how long each generated API key remains valid.
	TTL *string `json:"ttl,omitempty"`

	// Determines whether to use the same service ID and API key for future read operations on an
	// `iam_credentials` secret.
	//
	// If set to `true`, the service reuses the current credentials. If set to `false`, a new service ID and API key are
	// generated each time that the secret is read or accessed.
	ReuseAPIKey *bool `json:"reuse_api_key,omitempty"`

	// Indicates whether an `iam_credentials` secret was created with a static service ID.
	//
	// If the value is `true`, the service ID for the secret was provided by the user at secret creation. If the value is
	// `false`, the service ID was generated by Secrets Manager.
	ServiceIDIsStatic *bool `json:"service_id_is_static,omitempty"`

	// The service ID under which the API key is created. The service ID is included in the metadata only if the secret was
	// created with a static service ID.
	ServiceID *string `json:"service_id,omitempty"`

	// The access groups that define the capabilities of the service ID and API key that are generated for an
	// `iam_credentials` secret. The access groups are included in the metadata only if the secret was created with a
	// service ID that was generated by Secrets Manager.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the
	// certificate.
	Algorithm *string `json:"algorithm,omitempty"`

	// The identifier for the cryptographic algorithm that was used to generate the public and private keys that are
	// associated with the certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The fully qualified domain name or host domain name that is defined for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// Indicates whether the certificate was imported with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included,omitempty"`

	// Indicates whether the certificate was imported with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included,omitempty"`

	// The alternative names that are defined for the certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// Determines whether your issued certificate is bundled with intermediate certificates.
	//
	// Set to `false` for the certificate file to contain only the issued certificate.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	Rotation *Rotation `json:"rotation,omitempty"`

	// Issuance information that is associated with your certificate.
	IssuanceInfo *IssuanceInfo `json:"issuance_info,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template,omitempty"`

	// The intermediate certificate authority that signed this certificate.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`
}

// Constants associated with the SecretMetadata.SecretType property.
// The secret type.
const (
	SecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	SecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	SecretMetadataSecretTypeImportedCertConst     = "imported_cert"
	SecretMetadataSecretTypeKvConst               = "kv"
	SecretMetadataSecretTypePrivateCertConst      = "private_cert"
	SecretMetadataSecretTypePublicCertConst       = "public_cert"
	SecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

func (*SecretMetadata) isaSecretMetadata() bool {
	return true
}

type SecretMetadataIntf interface {
	isaSecretMetadata() bool
}

// UnmarshalSecretMetadata unmarshals an instance of SecretMetadata from the specified map of raw messages.
func UnmarshalSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseAPIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id_is_static", &obj.ServiceIDIsStatic)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "access_groups", &obj.AccessGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_included", &obj.IntermediateIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_included", &obj.PrivateKeyIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bundle_certs", &obj.BundleCerts)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotation)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "issuance_info", &obj.IssuanceInfo, UnmarshalIssuanceInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_template", &obj.CertificateTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
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

// SecretMetadataRequest : The metadata of a secret.
type SecretMetadataRequest struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretMetadataIntf `json:"resources" validate:"required"`
}

// NewSecretMetadataRequest : Instantiate SecretMetadataRequest (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretMetadataRequest(metadata *CollectionMetadata, resources []SecretMetadataIntf) (_model *SecretMetadataRequest, err error) {
	_model = &SecretMetadataRequest{
		Metadata:  metadata,
		Resources: resources,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSecretMetadataRequest unmarshals an instance of SecretMetadataRequest from the specified map of raw messages.
func UnmarshalSecretMetadataRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretMetadataRequest)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSecretMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretPolicyRotation : Properties that describe a rotation policy.
type SecretPolicyRotation struct {
	// The MIME type that represents the policy. Currently, only the default is supported.
	Type *string `json:"type" validate:"required"`

	Rotation SecretPolicyRotationRotationIntf `json:"rotation" validate:"required"`
}

// Constants associated with the SecretPolicyRotation.Type property.
// The MIME type that represents the policy. Currently, only the default is supported.
const (
	SecretPolicyRotationTypeApplicationVndIBMSecretsManagerSecretPolicyJSONConst = "application/vnd.ibm.secrets-manager.secret.policy+json"
)

// NewSecretPolicyRotation : Instantiate SecretPolicyRotation (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretPolicyRotation(typeVar string, rotation SecretPolicyRotationRotationIntf) (_model *SecretPolicyRotation, err error) {
	_model = &SecretPolicyRotation{
		Type:     core.StringPtr(typeVar),
		Rotation: rotation,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSecretPolicyRotation unmarshals an instance of SecretPolicyRotation from the specified map of raw messages.
func UnmarshalSecretPolicyRotation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretPolicyRotation)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalSecretPolicyRotationRotation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretPolicyRotationRotation : SecretPolicyRotationRotation struct
// Models which "extend" this model:
// - SecretPolicyRotationRotationPolicyRotation
// - SecretPolicyRotationRotationPublicCertPolicyRotation
// - PrivateCertPolicyRotation
type SecretPolicyRotationRotation struct {
	// The length of the secret rotation time interval.
	Interval *int64 `json:"interval,omitempty"`

	// The units for the secret rotation time interval.
	Unit *string `json:"unit,omitempty"`

	AutoRotate *bool `json:"auto_rotate,omitempty"`

	RotateKeys *bool `json:"rotate_keys,omitempty"`
}

// Constants associated with the SecretPolicyRotationRotation.Unit property.
// The units for the secret rotation time interval.
const (
	SecretPolicyRotationRotationUnitDayConst   = "day"
	SecretPolicyRotationRotationUnitMonthConst = "month"
)

func (*SecretPolicyRotationRotation) isaSecretPolicyRotationRotation() bool {
	return true
}

type SecretPolicyRotationRotationIntf interface {
	isaSecretPolicyRotationRotation() bool
}

// UnmarshalSecretPolicyRotationRotation unmarshals an instance of SecretPolicyRotationRotation from the specified map of raw messages.
func UnmarshalSecretPolicyRotationRotation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretPolicyRotationRotation)
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unit", &obj.Unit)
	if err != nil {
		return
	}
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

// SecretResource : SecretResource struct
// Models which "extend" this model:
// - ArbitrarySecretResource
// - UsernamePasswordSecretResource
// - IamCredentialsSecretResource
// - CertificateSecretResource
// - PublicCertificateSecretResource
// - PrivateCertificateSecretResource
// - KvSecretResource
type SecretResource struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name,omitempty"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// An array that contains metadata for each secret version. For more information on the metadata properties, see [Get
	// secret version metadata](#get-secret-version-metadata).
	Versions []map[string]interface{} `json:"versions,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The new secret data to assign to the secret.
	Payload interface{} `json:"payload,omitempty"`

	// The data that is associated with the secret version.
	//
	// The data object contains the field `payload`.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`

	// The username to assign to this secret.
	Username *string `json:"username,omitempty"`

	// The password to assign to this secret.
	Password *string `json:"password,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and have an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to generated credentials.
	//
	// For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be
	// either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m`
	// or `24h`.
	//
	// Minimum duration is 1 minute. Maximum is 90 days.
	TTL interface{} `json:"ttl,omitempty"`

	// The access groups that define the capabilities of the service ID and API key that are generated for an
	// `iam_credentials` secret. If you prefer to use an existing service ID that is already assigned the access policies
	// that you require, you can omit this parameter and use the `service_id` field instead.
	//
	// **Tip:** To list the access groups that are available in an account, you can use the [IAM Access Groups
	// API](https://cloud.ibm.com/apidocs/iam-access-groups#list-access-groups). To find the ID of an access group in the
	// console, go to **Manage > Access (IAM) > Access groups**. Select the access group to inspect, and click **Details**
	// to view its ID.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease (see the `ttl` field), the API key is deleted automatically. If you
	// want to continue to use the same API key for future read operations, see the `reuse_api_key` field.
	APIKey *string `json:"api_key,omitempty"`

	// The ID of the API key that is generated for this secret.
	APIKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// Indicates whether an `iam_credentials` secret was created with a static service ID.
	//
	// If `true`, the service ID for the secret was provided by the user at secret creation. If `false`, the service ID was
	// generated by Secrets Manager.
	ServiceIDIsStatic *bool `json:"service_id_is_static,omitempty"`

	// Determines whether to use the same service ID and API key for future read operations on an
	// `iam_credentials` secret.
	//
	// If set to `true`, the service reuses the current credentials. If set to `false`, a new service ID and API key are
	// generated each time that the secret is read or accessed.
	ReuseAPIKey *bool `json:"reuse_api_key,omitempty"`

	// The contents of your certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The private key to associate with the certificate. The data must be formatted on a single line with embedded newline
	// characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The intermediate certificate to associate with the root certificate. The data must be formatted on a single line
	// with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the
	// certificate.
	Algorithm *string `json:"algorithm,omitempty"`

	// The identifier for the cryptographic algorithm that was used to generate the public and private keys that are
	// associated with the certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The fully qualified domain name or host domain name that is defined for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// Indicates whether the certificate was imported with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included,omitempty"`

	// Indicates whether the certificate was imported with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included,omitempty"`

	// The alternative names that are defined for the certificate.
	//
	// For public certificates, this value is provided as an array of strings. For private certificates, this value is
	// provided as a comma-delimited list (string). In the API response, this value is returned as an array of strings for
	// all the types of certificate secrets.
	AltNames interface{} `json:"alt_names,omitempty"`

	// Determines whether your issued certificate is bundled with intermediate certificates.
	//
	// Set to `false` for the certificate file to contain only the issued certificate.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The name of the certificate authority configuration.
	//
	// To view a list of your configured authorities, use the [List configurations API](#get-secret-config-element).
	Ca *string `json:"ca,omitempty"`

	// The name of the DNS provider configuration.
	//
	// To view a list of your configured authorities, use the [List configurations API](#get-secret-config-element).
	DNS *string `json:"dns,omitempty"`

	Rotation *Rotation `json:"rotation,omitempty"`

	// Issuance information that is associated with your certificate.
	IssuanceInfo *IssuanceInfo `json:"issuance_info,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template,omitempty"`

	// The intermediate certificate authority that signed this certificate.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// The certificate signing request. If you don't include this parameter, the CSR that is used to generate the
	// certificate is created internally. If you provide a CSR, it is used also for auto rotation and manual rotation,
	// unless you provide another CSR in the manual rotation request.
	Csr *string `json:"csr,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

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

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`
}

// Constants associated with the SecretResource.SecretType property.
// The secret type.
const (
	SecretResourceSecretTypeArbitraryConst        = "arbitrary"
	SecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	SecretResourceSecretTypeImportedCertConst     = "imported_cert"
	SecretResourceSecretTypeKvConst               = "kv"
	SecretResourceSecretTypePrivateCertConst      = "private_cert"
	SecretResourceSecretTypePublicCertConst       = "public_cert"
	SecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the SecretResource.Format property.
// The format of the returned data.
const (
	SecretResourceFormatPemConst       = "pem"
	SecretResourceFormatPemBundleConst = "pem_bundle"
)

// Constants associated with the SecretResource.PrivateKeyFormat property.
// The format of the generated private key.
const (
	SecretResourcePrivateKeyFormatDerConst   = "der"
	SecretResourcePrivateKeyFormatPkcs8Const = "pkcs8"
)

func (*SecretResource) isaSecretResource() bool {
	return true
}

type SecretResourceIntf interface {
	isaSecretResource() bool
}

// UnmarshalSecretResource unmarshals an instance of SecretResource from the specified map of raw messages.
func UnmarshalSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretResource)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
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
	err = core.UnmarshalPrimitive(m, "next_rotation_date", &obj.NextRotationDate)
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
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.APIKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id_is_static", &obj.ServiceIDIsStatic)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseAPIKey)
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
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_included", &obj.IntermediateIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_included", &obj.PrivateKeyIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
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
	err = core.UnmarshalPrimitive(m, "dns", &obj.DNS)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotation)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "issuance_info", &obj.IssuanceInfo, UnmarshalIssuanceInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_template", &obj.CertificateTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
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

// SecretVersion : SecretVersion struct
// Models which "extend" this model:
// - ArbitrarySecretVersion
// - UsernamePasswordSecretVersion
// - IamCredentialsSecretVersion
// - CertificateSecretVersion
// - PrivateCertificateSecretVersion
type SecretVersion struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The data that is associated with the secret version.
	//
	// The data object contains the field `payload`.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`
}

func (*SecretVersion) isaSecretVersion() bool {
	return true
}

type SecretVersionIntf interface {
	isaSecretVersion() bool
}

// UnmarshalSecretVersion unmarshals an instance of SecretVersion from the specified map of raw messages.
func UnmarshalSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
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
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
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

// SecretVersionInfo : Properties that describe a secret version within a list of secret versions.
// Models which "extend" this model:
// - ArbitrarySecretVersionInfo
// - UsernamePasswordSecretVersionInfo
// - IamCredentialsSecretVersionInfo
// - CertificateSecretVersionInfo
// - PrivateCertificateSecretVersionInfo
type SecretVersionInfo struct {
	// The ID of the secret version.
	ID *string `json:"id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`
}

func (*SecretVersionInfo) isaSecretVersionInfo() bool {
	return true
}

type SecretVersionInfoIntf interface {
	isaSecretVersionInfo() bool
}

// UnmarshalSecretVersionInfo unmarshals an instance of SecretVersionInfo from the specified map of raw messages.
func UnmarshalSecretVersionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersionInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
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
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
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

// SecretVersionMetadata : SecretVersionMetadata struct
// Models which "extend" this model:
// - ArbitrarySecretVersionMetadata
// - UsernamePasswordSecretVersionMetadata
// - IamCredentialsSecretVersionMetadata
// - CertificateSecretVersionMetadata
// - PrivateCertificateSecretVersionMetadata
type SecretVersionMetadata struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`
}

func (*SecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

type SecretVersionMetadataIntf interface {
	isaSecretVersionMetadata() bool
}

// UnmarshalSecretVersionMetadata unmarshals an instance of SecretVersionMetadata from the specified map of raw messages.
func UnmarshalSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
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
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
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

// SecretsLocks : Properties that describe the secret locks.
type SecretsLocks struct {
	// The unique ID of the secret.
	SecretID *string `json:"secret_id,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// A collection of locks that are attached to a secret version.
	Versions []SecretLockVersion `json:"versions,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of SecretsLocks
func (o *SecretsLocks) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretsLocks
func (o *SecretsLocks) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretsLocks
func (o *SecretsLocks) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretsLocks
func (o *SecretsLocks) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretsLocks
func (o *SecretsLocks) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.SecretID != nil {
		m["secret_id"] = o.SecretID
	}
	if o.SecretGroupID != nil {
		m["secret_group_id"] = o.SecretGroupID
	}
	if o.Versions != nil {
		m["versions"] = o.Versions
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalSecretsLocks unmarshals an instance of SecretsLocks from the specified map of raw messages.
func UnmarshalSecretsLocks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretsLocks)
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		return
	}
	delete(m, "secret_id")
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	delete(m, "secret_group_id")
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalSecretLockVersion)
	if err != nil {
		return
	}
	delete(m, "versions")
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

// SendTestNotificationOptions : The SendTestNotification options.
type SendTestNotificationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSendTestNotificationOptions : Instantiate SendTestNotificationOptions
func (*SecretsManagerV1) NewSendTestNotificationOptions() *SendTestNotificationOptions {
	return &SendTestNotificationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *SendTestNotificationOptions) SetHeaders(param map[string]string) *SendTestNotificationOptions {
	options.Headers = param
	return options
}

// SignActionResultData : Properties that are returned with a successful `sign` action.
type SignActionResultData struct {
	// The PEM-encoded certificate.
	Certificate *string `json:"certificate,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The PEM-encoded certificate of the certificate authority that signed and issued this certificate.
	IssuingCa *string `json:"issuing_ca,omitempty"`

	// The chain of certificate authorities that are associated with the certificate.
	CaChain []string `json:"ca_chain,omitempty"`

	// The time until the certificate expires.
	Expiration *int64 `json:"expiration,omitempty"`
}

// UnmarshalSignActionResultData unmarshals an instance of SignActionResultData from the specified map of raw messages.
func UnmarshalSignActionResultData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SignActionResultData)
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
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

// SignIntermediateActionResultData : Properties that are returned with a successful `sign` action.
type SignIntermediateActionResultData struct {
	// The PEM-encoded certificate.
	Certificate *string `json:"certificate,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The PEM-encoded certificate of the certificate authority that signed and issued this certificate.
	IssuingCa *string `json:"issuing_ca,omitempty"`

	// The chain of certificate authorities that are associated with the certificate.
	CaChain []string `json:"ca_chain,omitempty"`

	// The time until the certificate expires.
	Expiration *int64 `json:"expiration,omitempty"`
}

// UnmarshalSignIntermediateActionResultData unmarshals an instance of SignIntermediateActionResultData from the specified map of raw messages.
func UnmarshalSignIntermediateActionResultData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SignIntermediateActionResultData)
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
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

// UnlockSecretOptions : The UnlockSecret options.
type UnlockSecretOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// A comma-separated list of locks to delete.
	Locks []string `json:"locks,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UnlockSecretOptions.SecretType property.
// The secret type.
const (
	UnlockSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	UnlockSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	UnlockSecretOptionsSecretTypeImportedCertConst     = "imported_cert"
	UnlockSecretOptionsSecretTypeKvConst               = "kv"
	UnlockSecretOptionsSecretTypePrivateCertConst      = "private_cert"
	UnlockSecretOptionsSecretTypePublicCertConst       = "public_cert"
	UnlockSecretOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewUnlockSecretOptions : Instantiate UnlockSecretOptions
func (*SecretsManagerV1) NewUnlockSecretOptions(secretType string, id string) *UnlockSecretOptions {
	return &UnlockSecretOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *UnlockSecretOptions) SetSecretType(secretType string) *UnlockSecretOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *UnlockSecretOptions) SetID(id string) *UnlockSecretOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLocks : Allow user to set Locks
func (_options *UnlockSecretOptions) SetLocks(locks []string) *UnlockSecretOptions {
	_options.Locks = locks
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UnlockSecretOptions) SetHeaders(param map[string]string) *UnlockSecretOptions {
	options.Headers = param
	return options
}

// UnlockSecretVersionOptions : The UnlockSecretVersion options.
type UnlockSecretVersionOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret version. You can also use `previous` to retrieve the previous
	// version.
	//
	// **Note:** To find the version ID of a secret, use the [Get secret metadata](#get-secret-metadata) method and check
	// the response details.
	VersionID *string `json:"version_id" validate:"required,ne="`

	// A comma-separated list of locks to delete.
	Locks []string `json:"locks,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UnlockSecretVersionOptions.SecretType property.
// The secret type.
const (
	UnlockSecretVersionOptionsSecretTypeArbitraryConst        = "arbitrary"
	UnlockSecretVersionOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	UnlockSecretVersionOptionsSecretTypeImportedCertConst     = "imported_cert"
	UnlockSecretVersionOptionsSecretTypeKvConst               = "kv"
	UnlockSecretVersionOptionsSecretTypePrivateCertConst      = "private_cert"
	UnlockSecretVersionOptionsSecretTypePublicCertConst       = "public_cert"
	UnlockSecretVersionOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewUnlockSecretVersionOptions : Instantiate UnlockSecretVersionOptions
func (*SecretsManagerV1) NewUnlockSecretVersionOptions(secretType string, id string, versionID string) *UnlockSecretVersionOptions {
	return &UnlockSecretVersionOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		VersionID:  core.StringPtr(versionID),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *UnlockSecretVersionOptions) SetSecretType(secretType string) *UnlockSecretVersionOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *UnlockSecretVersionOptions) SetID(id string) *UnlockSecretVersionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *UnlockSecretVersionOptions) SetVersionID(versionID string) *UnlockSecretVersionOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetLocks : Allow user to set Locks
func (_options *UnlockSecretVersionOptions) SetLocks(locks []string) *UnlockSecretVersionOptions {
	_options.Locks = locks
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UnlockSecretVersionOptions) SetHeaders(param map[string]string) *UnlockSecretVersionOptions {
	options.Headers = param
	return options
}

// UpdateConfigElementOptions : The UpdateConfigElement options.
type UpdateConfigElementOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The configuration element to define or manage.
	ConfigElement *string `json:"config_element" validate:"required,ne="`

	// The name of your configuration.
	ConfigName *string `json:"config_name" validate:"required,ne="`

	// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
	Type *string `json:"type" validate:"required"`

	// Properties that describe a configuration, which depends on type.
	Config map[string]interface{} `json:"config" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateConfigElementOptions.SecretType property.
// The secret type.
const (
	UpdateConfigElementOptionsSecretTypePrivateCertConst = "private_cert"
	UpdateConfigElementOptionsSecretTypePublicCertConst  = "public_cert"
)

// Constants associated with the UpdateConfigElementOptions.ConfigElement property.
// The configuration element to define or manage.
const (
	UpdateConfigElementOptionsConfigElementCertificateAuthoritiesConst             = "certificate_authorities"
	UpdateConfigElementOptionsConfigElementCertificateTemplatesConst               = "certificate_templates"
	UpdateConfigElementOptionsConfigElementDNSProvidersConst                       = "dns_providers"
	UpdateConfigElementOptionsConfigElementIntermediateCertificateAuthoritiesConst = "intermediate_certificate_authorities"
	UpdateConfigElementOptionsConfigElementRootCertificateAuthoritiesConst         = "root_certificate_authorities"
)

// Constants associated with the UpdateConfigElementOptions.Type property.
// The type of configuration. Value options differ depending on the `config_element` property that you want to define.
const (
	UpdateConfigElementOptionsTypeCertificateTemplateConst              = "certificate_template"
	UpdateConfigElementOptionsTypeCisConst                              = "cis"
	UpdateConfigElementOptionsTypeClassicInfrastructureConst            = "classic_infrastructure"
	UpdateConfigElementOptionsTypeIntermediateCertificateAuthorityConst = "intermediate_certificate_authority"
	UpdateConfigElementOptionsTypeLetsencryptConst                      = "letsencrypt"
	UpdateConfigElementOptionsTypeLetsencryptStageConst                 = "letsencrypt-stage"
	UpdateConfigElementOptionsTypeRootCertificateAuthorityConst         = "root_certificate_authority"
)

// NewUpdateConfigElementOptions : Instantiate UpdateConfigElementOptions
func (*SecretsManagerV1) NewUpdateConfigElementOptions(secretType string, configElement string, configName string, typeVar string, config map[string]interface{}) *UpdateConfigElementOptions {
	return &UpdateConfigElementOptions{
		SecretType:    core.StringPtr(secretType),
		ConfigElement: core.StringPtr(configElement),
		ConfigName:    core.StringPtr(configName),
		Type:          core.StringPtr(typeVar),
		Config:        config,
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *UpdateConfigElementOptions) SetSecretType(secretType string) *UpdateConfigElementOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetConfigElement : Allow user to set ConfigElement
func (_options *UpdateConfigElementOptions) SetConfigElement(configElement string) *UpdateConfigElementOptions {
	_options.ConfigElement = core.StringPtr(configElement)
	return _options
}

// SetConfigName : Allow user to set ConfigName
func (_options *UpdateConfigElementOptions) SetConfigName(configName string) *UpdateConfigElementOptions {
	_options.ConfigName = core.StringPtr(configName)
	return _options
}

// SetType : Allow user to set Type
func (_options *UpdateConfigElementOptions) SetType(typeVar string) *UpdateConfigElementOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetConfig : Allow user to set Config
func (_options *UpdateConfigElementOptions) SetConfig(config map[string]interface{}) *UpdateConfigElementOptions {
	_options.Config = config
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateConfigElementOptions) SetHeaders(param map[string]string) *UpdateConfigElementOptions {
	options.Headers = param
	return options
}

// UpdateSecretGroupMetadataOptions : The UpdateSecretGroupMetadata options.
type UpdateSecretGroupMetadataOptions struct {
	// The v4 UUID that uniquely identifies the secret group.
	ID *string `json:"id" validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretGroupMetadataUpdatable `json:"resources" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSecretGroupMetadataOptions : Instantiate UpdateSecretGroupMetadataOptions
func (*SecretsManagerV1) NewUpdateSecretGroupMetadataOptions(id string, metadata *CollectionMetadata, resources []SecretGroupMetadataUpdatable) *UpdateSecretGroupMetadataOptions {
	return &UpdateSecretGroupMetadataOptions{
		ID:        core.StringPtr(id),
		Metadata:  metadata,
		Resources: resources,
	}
}

// SetID : Allow user to set ID
func (_options *UpdateSecretGroupMetadataOptions) SetID(id string) *UpdateSecretGroupMetadataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *UpdateSecretGroupMetadataOptions) SetMetadata(metadata *CollectionMetadata) *UpdateSecretGroupMetadataOptions {
	_options.Metadata = metadata
	return _options
}

// SetResources : Allow user to set Resources
func (_options *UpdateSecretGroupMetadataOptions) SetResources(resources []SecretGroupMetadataUpdatable) *UpdateSecretGroupMetadataOptions {
	_options.Resources = resources
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretGroupMetadataOptions) SetHeaders(param map[string]string) *UpdateSecretGroupMetadataOptions {
	options.Headers = param
	return options
}

// UpdateSecretMetadataOptions : The UpdateSecretMetadata options.
type UpdateSecretMetadataOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretMetadataIntf `json:"resources" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateSecretMetadataOptions.SecretType property.
// The secret type.
const (
	UpdateSecretMetadataOptionsSecretTypeArbitraryConst        = "arbitrary"
	UpdateSecretMetadataOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	UpdateSecretMetadataOptionsSecretTypeImportedCertConst     = "imported_cert"
	UpdateSecretMetadataOptionsSecretTypeKvConst               = "kv"
	UpdateSecretMetadataOptionsSecretTypePrivateCertConst      = "private_cert"
	UpdateSecretMetadataOptionsSecretTypePublicCertConst       = "public_cert"
	UpdateSecretMetadataOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewUpdateSecretMetadataOptions : Instantiate UpdateSecretMetadataOptions
func (*SecretsManagerV1) NewUpdateSecretMetadataOptions(secretType string, id string, metadata *CollectionMetadata, resources []SecretMetadataIntf) *UpdateSecretMetadataOptions {
	return &UpdateSecretMetadataOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		Metadata:   metadata,
		Resources:  resources,
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *UpdateSecretMetadataOptions) SetSecretType(secretType string) *UpdateSecretMetadataOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateSecretMetadataOptions) SetID(id string) *UpdateSecretMetadataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *UpdateSecretMetadataOptions) SetMetadata(metadata *CollectionMetadata) *UpdateSecretMetadataOptions {
	_options.Metadata = metadata
	return _options
}

// SetResources : Allow user to set Resources
func (_options *UpdateSecretMetadataOptions) SetResources(resources []SecretMetadataIntf) *UpdateSecretMetadataOptions {
	_options.Resources = resources
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretMetadataOptions) SetHeaders(param map[string]string) *UpdateSecretMetadataOptions {
	options.Headers = param
	return options
}

// UpdateSecretOptions : The UpdateSecret options.
type UpdateSecretOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The action to perform on the specified secret.
	Action *string `json:"action" validate:"required"`

	// The properties to update for the secret.
	SecretAction SecretActionIntf `json:"SecretAction,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateSecretOptions.SecretType property.
// The secret type.
const (
	UpdateSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	UpdateSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	UpdateSecretOptionsSecretTypeImportedCertConst     = "imported_cert"
	UpdateSecretOptionsSecretTypeKvConst               = "kv"
	UpdateSecretOptionsSecretTypePrivateCertConst      = "private_cert"
	UpdateSecretOptionsSecretTypePublicCertConst       = "public_cert"
	UpdateSecretOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the UpdateSecretOptions.Action property.
// The action to perform on the specified secret.
const (
	UpdateSecretOptionsActionDeleteCredentialsConst    = "delete_credentials"
	UpdateSecretOptionsActionRestoreConst              = "restore"
	UpdateSecretOptionsActionRevokeConst               = "revoke"
	UpdateSecretOptionsActionRotateConst               = "rotate"
	UpdateSecretOptionsActionValidateDNSChallengeConst = "validate_dns_challenge"
)

// NewUpdateSecretOptions : Instantiate UpdateSecretOptions
func (*SecretsManagerV1) NewUpdateSecretOptions(secretType string, id string, action string) *UpdateSecretOptions {
	return &UpdateSecretOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		Action:     core.StringPtr(action),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *UpdateSecretOptions) SetSecretType(secretType string) *UpdateSecretOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateSecretOptions) SetID(id string) *UpdateSecretOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAction : Allow user to set Action
func (_options *UpdateSecretOptions) SetAction(action string) *UpdateSecretOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetSecretAction : Allow user to set SecretAction
func (_options *UpdateSecretOptions) SetSecretAction(secretAction SecretActionIntf) *UpdateSecretOptions {
	_options.SecretAction = secretAction
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretOptions) SetHeaders(param map[string]string) *UpdateSecretOptions {
	options.Headers = param
	return options
}

// UpdateSecretVersionMetadata : Properties that update the metadata of a secret version.
type UpdateSecretVersionMetadata struct {
	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// UnmarshalUpdateSecretVersionMetadata unmarshals an instance of UpdateSecretVersionMetadata from the specified map of raw messages.
func UnmarshalUpdateSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateSecretVersionMetadataOptions : The UpdateSecretVersionMetadata options.
type UpdateSecretVersionMetadataOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret version. You can also use `previous` to retrieve the previous
	// version.
	//
	// **Note:** To find the version ID of a secret, use the [Get secret metadata](#get-secret-metadata) method and check
	// the response details.
	VersionID *string `json:"version_id" validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []UpdateSecretVersionMetadata `json:"resources" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateSecretVersionMetadataOptions.SecretType property.
// The secret type.
const (
	UpdateSecretVersionMetadataOptionsSecretTypeArbitraryConst        = "arbitrary"
	UpdateSecretVersionMetadataOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	UpdateSecretVersionMetadataOptionsSecretTypeImportedCertConst     = "imported_cert"
	UpdateSecretVersionMetadataOptionsSecretTypeKvConst               = "kv"
	UpdateSecretVersionMetadataOptionsSecretTypePrivateCertConst      = "private_cert"
	UpdateSecretVersionMetadataOptionsSecretTypePublicCertConst       = "public_cert"
	UpdateSecretVersionMetadataOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewUpdateSecretVersionMetadataOptions : Instantiate UpdateSecretVersionMetadataOptions
func (*SecretsManagerV1) NewUpdateSecretVersionMetadataOptions(secretType string, id string, versionID string, metadata *CollectionMetadata, resources []UpdateSecretVersionMetadata) *UpdateSecretVersionMetadataOptions {
	return &UpdateSecretVersionMetadataOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		VersionID:  core.StringPtr(versionID),
		Metadata:   metadata,
		Resources:  resources,
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *UpdateSecretVersionMetadataOptions) SetSecretType(secretType string) *UpdateSecretVersionMetadataOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateSecretVersionMetadataOptions) SetID(id string) *UpdateSecretVersionMetadataOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *UpdateSecretVersionMetadataOptions) SetVersionID(versionID string) *UpdateSecretVersionMetadataOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *UpdateSecretVersionMetadataOptions) SetMetadata(metadata *CollectionMetadata) *UpdateSecretVersionMetadataOptions {
	_options.Metadata = metadata
	return _options
}

// SetResources : Allow user to set Resources
func (_options *UpdateSecretVersionMetadataOptions) SetResources(resources []UpdateSecretVersionMetadata) *UpdateSecretVersionMetadataOptions {
	_options.Resources = resources
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretVersionMetadataOptions) SetHeaders(param map[string]string) *UpdateSecretVersionMetadataOptions {
	options.Headers = param
	return options
}

// UpdateSecretVersionOptions : The UpdateSecretVersion options.
type UpdateSecretVersionOptions struct {
	// The secret type.
	SecretType *string `json:"secret_type" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id" validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret version. You can also use `previous` to retrieve the previous
	// version.
	//
	// **Note:** To find the version ID of a secret, use the [Get secret metadata](#get-secret-metadata) method and check
	// the response details.
	VersionID *string `json:"version_id" validate:"required,ne="`

	// The action to perform on the specified secret version.
	Action *string `json:"action" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateSecretVersionOptions.SecretType property.
// The secret type.
const (
	UpdateSecretVersionOptionsSecretTypePrivateCertConst = "private_cert"
)

// Constants associated with the UpdateSecretVersionOptions.Action property.
// The action to perform on the specified secret version.
const (
	UpdateSecretVersionOptionsActionRevokeConst = "revoke"
)

// NewUpdateSecretVersionOptions : Instantiate UpdateSecretVersionOptions
func (*SecretsManagerV1) NewUpdateSecretVersionOptions(secretType string, id string, versionID string, action string) *UpdateSecretVersionOptions {
	return &UpdateSecretVersionOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		VersionID:  core.StringPtr(versionID),
		Action:     core.StringPtr(action),
	}
}

// SetSecretType : Allow user to set SecretType
func (_options *UpdateSecretVersionOptions) SetSecretType(secretType string) *UpdateSecretVersionOptions {
	_options.SecretType = core.StringPtr(secretType)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateSecretVersionOptions) SetID(id string) *UpdateSecretVersionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *UpdateSecretVersionOptions) SetVersionID(versionID string) *UpdateSecretVersionOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetAction : Allow user to set Action
func (_options *UpdateSecretVersionOptions) SetAction(action string) *UpdateSecretVersionOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretVersionOptions) SetHeaders(param map[string]string) *UpdateSecretVersionOptions {
	options.Headers = param
	return options
}

// CertificateValidity : CertificateValidity struct
type CertificateValidity struct {
	// The date and time that the certificate validity period begins.
	NotBefore *strfmt.DateTime `json:"not_before,omitempty"`

	// The date and time that the certificate validity period ends.
	NotAfter *strfmt.DateTime `json:"not_after,omitempty"`
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

// ArbitrarySecretMetadata : Metadata properties that describe an arbitrary secret.
// This model "extends" SecretMetadata
type ArbitrarySecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be in the range 2 - 30 characters, including spaces. Special characters
	// that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`
}

// Constants associated with the ArbitrarySecretMetadata.SecretType property.
// The secret type.
const (
	ArbitrarySecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	ArbitrarySecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	ArbitrarySecretMetadataSecretTypeImportedCertConst     = "imported_cert"
	ArbitrarySecretMetadataSecretTypeKvConst               = "kv"
	ArbitrarySecretMetadataSecretTypePrivateCertConst      = "private_cert"
	ArbitrarySecretMetadataSecretTypePublicCertConst       = "public_cert"
	ArbitrarySecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

// NewArbitrarySecretMetadata : Instantiate ArbitrarySecretMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewArbitrarySecretMetadata(name string) (_model *ArbitrarySecretMetadata, err error) {
	_model = &ArbitrarySecretMetadata{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ArbitrarySecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalArbitrarySecretMetadata unmarshals an instance of ArbitrarySecretMetadata from the specified map of raw messages.
func UnmarshalArbitrarySecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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

// ArbitrarySecretResource : Properties that describe a secret.
// This model "extends" SecretResource
type ArbitrarySecretResource struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// An array that contains metadata for each secret version. For more information on the metadata properties, see [Get
	// secret version metadata](#get-secret-version-metadata).
	Versions []map[string]interface{} `json:"versions,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The new secret data to assign to the secret.
	Payload interface{} `json:"payload,omitempty"`

	// The data that is associated with the secret version.
	//
	// The data object contains the field `payload`.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`
}

// Constants associated with the ArbitrarySecretResource.SecretType property.
// The secret type.
const (
	ArbitrarySecretResourceSecretTypeArbitraryConst        = "arbitrary"
	ArbitrarySecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	ArbitrarySecretResourceSecretTypeImportedCertConst     = "imported_cert"
	ArbitrarySecretResourceSecretTypeKvConst               = "kv"
	ArbitrarySecretResourceSecretTypePrivateCertConst      = "private_cert"
	ArbitrarySecretResourceSecretTypePublicCertConst       = "public_cert"
	ArbitrarySecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// NewArbitrarySecretResource : Instantiate ArbitrarySecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewArbitrarySecretResource(name string) (_model *ArbitrarySecretResource, err error) {
	_model = &ArbitrarySecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ArbitrarySecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalArbitrarySecretResource unmarshals an instance of ArbitrarySecretResource from the specified map of raw messages.
func UnmarshalArbitrarySecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretResource)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ArbitrarySecretVersion : ArbitrarySecretVersion struct
// This model "extends" SecretVersion
type ArbitrarySecretVersion struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The data that is associated with the secret version.
	//
	// The data object contains the field `payload`.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`
}

func (*ArbitrarySecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalArbitrarySecretVersion unmarshals an instance of ArbitrarySecretVersion from the specified map of raw messages.
func UnmarshalArbitrarySecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ArbitrarySecretVersionInfo : ArbitrarySecretVersionInfo struct
// This model "extends" SecretVersionInfo
type ArbitrarySecretVersionInfo struct {
	// The ID of the secret version.
	ID *string `json:"id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

func (*ArbitrarySecretVersionInfo) isaSecretVersionInfo() bool {
	return true
}

// UnmarshalArbitrarySecretVersionInfo unmarshals an instance of ArbitrarySecretVersionInfo from the specified map of raw messages.
func UnmarshalArbitrarySecretVersionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretVersionInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
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

// ArbitrarySecretVersionMetadata : Properties that describe a secret version.
// This model "extends" SecretVersionMetadata
type ArbitrarySecretVersionMetadata struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

func (*ArbitrarySecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalArbitrarySecretVersionMetadata unmarshals an instance of ArbitrarySecretVersionMetadata from the specified map of raw messages.
func UnmarshalArbitrarySecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArbitrarySecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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

// CertificateSecretMetadata : Metadata properties that describe a certificate secret.
// This model "extends" SecretMetadata
type CertificateSecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be in the range 2 - 30 characters, including spaces. Special characters
	// that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the
	// certificate.
	Algorithm *string `json:"algorithm,omitempty"`

	// The identifier for the cryptographic algorithm that was used to generate the public and private keys that are
	// associated with the certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The fully qualified domain name or host domain name that is defined for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// Indicates whether the certificate was imported with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included,omitempty"`

	// Indicates whether the certificate was imported with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included,omitempty"`

	// The alternative names that are defined for the certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`
}

// Constants associated with the CertificateSecretMetadata.SecretType property.
// The secret type.
const (
	CertificateSecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	CertificateSecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	CertificateSecretMetadataSecretTypeImportedCertConst     = "imported_cert"
	CertificateSecretMetadataSecretTypeKvConst               = "kv"
	CertificateSecretMetadataSecretTypePrivateCertConst      = "private_cert"
	CertificateSecretMetadataSecretTypePublicCertConst       = "public_cert"
	CertificateSecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

// NewCertificateSecretMetadata : Instantiate CertificateSecretMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewCertificateSecretMetadata(name string) (_model *CertificateSecretMetadata, err error) {
	_model = &CertificateSecretMetadata{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*CertificateSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalCertificateSecretMetadata unmarshals an instance of CertificateSecretMetadata from the specified map of raw messages.
func UnmarshalCertificateSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateSecretMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_included", &obj.IntermediateIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_included", &obj.PrivateKeyIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
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

// CertificateSecretResource : Properties that describe a secret.
// This model "extends" SecretResource
type CertificateSecretResource struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// An array that contains metadata for each secret version. For more information on the metadata properties, see [Get
	// secret version metadata](#get-secret-version-metadata).
	Versions []map[string]interface{} `json:"versions,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The contents of your certificate. The data must be formatted on a single line with embedded newline characters.
	Certificate *string `json:"certificate,omitempty"`

	// The private key to associate with the certificate. The data must be formatted on a single line with embedded newline
	// characters.
	PrivateKey *string `json:"private_key,omitempty"`

	// The intermediate certificate to associate with the root certificate. The data must be formatted on a single line
	// with embedded newline characters.
	Intermediate *string `json:"intermediate,omitempty"`

	// The data that is associated with the secret. The data object contains the following fields:
	//
	// - `certificate`: The contents of the certificate.
	// - `private_key`: The private key that is associated with the certificate.
	// - `intermediate`: The intermediate certificate that is associated with the certificate.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the
	// certificate.
	Algorithm *string `json:"algorithm,omitempty"`

	// The identifier for the cryptographic algorithm that was used to generate the public and private keys that are
	// associated with the certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The fully qualified domain name or host domain name that is defined for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// Indicates whether the certificate was imported with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included,omitempty"`

	// Indicates whether the certificate was imported with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included,omitempty"`

	// The alternative names that are defined for the certificate.
	//
	// For public certificates, this value is provided as an array of strings. For private certificates, this value is
	// provided as a comma-delimited list (string). In the API response, this value is returned as an array of strings for
	// all the types of certificate secrets.
	AltNames interface{} `json:"alt_names,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`
}

// Constants associated with the CertificateSecretResource.SecretType property.
// The secret type.
const (
	CertificateSecretResourceSecretTypeArbitraryConst        = "arbitrary"
	CertificateSecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	CertificateSecretResourceSecretTypeImportedCertConst     = "imported_cert"
	CertificateSecretResourceSecretTypeKvConst               = "kv"
	CertificateSecretResourceSecretTypePrivateCertConst      = "private_cert"
	CertificateSecretResourceSecretTypePublicCertConst       = "public_cert"
	CertificateSecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// NewCertificateSecretResource : Instantiate CertificateSecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewCertificateSecretResource(name string) (_model *CertificateSecretResource, err error) {
	_model = &CertificateSecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*CertificateSecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalCertificateSecretResource unmarshals an instance of CertificateSecretResource from the specified map of raw messages.
func UnmarshalCertificateSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateSecretResource)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_included", &obj.IntermediateIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_included", &obj.PrivateKeyIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
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

// CertificateSecretVersion : CertificateSecretVersion struct
// This model "extends" SecretVersion
type CertificateSecretVersion struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The data that is associated with the secret version. The data object contains the following fields:
	//
	// - `certificate`: The contents of the certificate.
	// - `private_key`: The private key that is associated with the certificate.
	// - `intermediate`: The intermediate certificate that is associated with the certificate.
	SecretData *CertificateSecretData `json:"secret_data,omitempty"`
}

func (*CertificateSecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalCertificateSecretVersion unmarshals an instance of CertificateSecretVersion from the specified map of raw messages.
func UnmarshalCertificateSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateSecretVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secret_data", &obj.SecretData, UnmarshalCertificateSecretData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CertificateSecretVersionInfo : CertificateSecretVersionInfo struct
// This model "extends" SecretVersionInfo
type CertificateSecretVersionInfo struct {
	// The ID of the secret version.
	ID *string `json:"id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`
}

func (*CertificateSecretVersionInfo) isaSecretVersionInfo() bool {
	return true
}

// UnmarshalCertificateSecretVersionInfo unmarshals an instance of CertificateSecretVersionInfo from the specified map of raw messages.
func UnmarshalCertificateSecretVersionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateSecretVersionInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
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

// CertificateSecretVersionMetadata : Properties that describe a secret version.
// This model "extends" SecretVersionMetadata
type CertificateSecretVersionMetadata struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`
}

func (*CertificateSecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalCertificateSecretVersionMetadata unmarshals an instance of CertificateSecretVersionMetadata from the specified map of raw messages.
func UnmarshalCertificateSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
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

// CertificateTemplateConfig : Properties that describe a certificate template. You can use a certificate template to control the parameters that
// are applied to your issued private certificates. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-certificate-templates).
// This model "extends" ConfigElementDefConfig
type CertificateTemplateConfig struct {
	// The name of the intermediate certificate authority.
	CertificateAuthority *string `json:"certificate_authority" validate:"required"`

	// Scopes the creation of private certificates to only the secret groups that you specify.
	//
	// This field can be supplied as a comma-delimited list of secret group IDs.
	AllowedSecretGroups *string `json:"allowed_secret_groups,omitempty"`

	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL interface{} `json:"max_ttl,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `12h`. The value can be supplied in
	// seconds (suffix `s`), minutes (suffix `m`) or hours (suffix `h`). The value can't exceed the `max_ttl` that is
	// defined in the associated certificate template. In the API response, this value is returned in seconds (integer).
	TTL interface{} `json:"ttl,omitempty"`

	// Determines whether to allow `localhost` to be included as one of the requested common names.
	AllowLocalhost *bool `json:"allow_localhost,omitempty"`

	// The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and
	// `allow_subdomains` options.
	AllowedDomains []string `json:"allowed_domains,omitempty"`

	// Determines whether to allow the domains that are supplied in the `allowed_domains` field to contain access control
	// list (ACL) templates.
	AllowedDomainsTemplate *bool `json:"allowed_domains_template,omitempty"`

	// Determines whether to allow clients to request private certificates that match the value of the actual domains on
	// the final certificate.
	//
	// For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a
	// certificate that contains the name `example.com` as one of the DNS values on the final certificate.
	//
	// **Important:** In some scenarios, allowing bare domains can be considered a security risk.
	AllowBareDomains *bool `json:"allow_bare_domains,omitempty"`

	// Determines whether to allow clients to request private certificates with common names (CN) that are subdomains of
	// the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.
	//
	// For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the
	// following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.
	//
	// **Note:** This field is redundant if you use the `allow_any_name` option.
	AllowSubdomains *bool `json:"allow_subdomains,omitempty"`

	// Determines whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are specified in the
	// `allowed_domains` field.
	//
	// If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	AllowGlobDomains *bool `json:"allow_glob_domains,omitempty"`

	// Determines whether to allow clients to request a private certificate that matches any common name.
	AllowAnyName *bool `json:"allow_any_name,omitempty"`

	// Determines whether to enforce only valid host names for common names, DNS Subject Alternative Names, and the host
	// section of email addresses.
	EnforceHostnames *bool `json:"enforce_hostnames,omitempty"`

	// Determines whether to allow clients to request a private certificate with IP Subject Alternative Names.
	AllowIPSans *bool `json:"allow_ip_sans,omitempty"`

	// The URI Subject Alternative Names to allow for private certificates.
	//
	// Values can contain glob patterns, for example `spiffe://hostname/_*`.
	AllowedURISans []string `json:"allowed_uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private
	// certificates.
	//
	// The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type
	// is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any
	// `other_sans` input.
	AllowedOtherSans []string `json:"allowed_other_sans,omitempty"`

	// Determines whether private certificates are flagged for server use.
	ServerFlag *bool `json:"server_flag,omitempty"`

	// Determines whether private certificates are flagged for client use.
	ClientFlag *bool `json:"client_flag,omitempty"`

	// Determines whether private certificates are flagged for code signing use.
	CodeSigningFlag *bool `json:"code_signing_flag,omitempty"`

	// Determines whether private certificates are flagged for email protection use.
	EmailProtectionFlag *bool `json:"email_protection_flag,omitempty"`

	// The type of private key to generate for private certificates and the type of key that is expected for submitted
	// certificate signing requests (CSRs).
	//
	// Allowable values are: `rsa` and `ec`.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use when generating the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The allowed key usage constraint to define for private certificates.
	//
	// You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage).  Omit the
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

	// When used with the `sign_csr` action, this field determines whether to use the common name (CN) from a certificate
	// signing request (CSR) instead of the CN that's included in the JSON data of the certificate.
	//
	// Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include
	// the `use_csr_sans` property.
	UseCsrCommonName *bool `json:"use_csr_common_name,omitempty"`

	// When used with the `sign_csr` action, this field determines whether to use the Subject Alternative Names
	// (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the JSON data of the
	// certificate.
	//
	// Does not include the common name in the CSR. To use the common name, include the `use_csr_common_name` property.
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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// Determines whether to require a common name to create a private certificate.
	//
	// By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the
	// `require_cn` option to `false`.
	RequireCn *bool `json:"require_cn,omitempty"`

	// A list of policy Object Identifiers (OIDs).
	PolicyIdentifiers []string `json:"policy_identifiers,omitempty"`

	// Determines whether to mark the Basic Constraints extension of an issued private certificate as valid for non-CA
	// certificates.
	BasicConstraintsValidForNonCa *bool `json:"basic_constraints_valid_for_non_ca,omitempty"`

	// The duration in seconds by which to backdate the `not_before` property of an issued private certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value
	// is returned in seconds (integer).
	NotBeforeDuration interface{} `json:"not_before_duration,omitempty"`
}

// Constants associated with the CertificateTemplateConfig.KeyType property.
// The type of private key to generate for private certificates and the type of key that is expected for submitted
// certificate signing requests (CSRs).
//
// Allowable values are: `rsa` and `ec`.
const (
	CertificateTemplateConfigKeyTypeEcConst  = "ec"
	CertificateTemplateConfigKeyTypeRsaConst = "rsa"
)

// NewCertificateTemplateConfig : Instantiate CertificateTemplateConfig (Generic Model Constructor)
func (*SecretsManagerV1) NewCertificateTemplateConfig(certificateAuthority string) (_model *CertificateTemplateConfig, err error) {
	_model = &CertificateTemplateConfig{
		CertificateAuthority: core.StringPtr(certificateAuthority),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*CertificateTemplateConfig) isaConfigElementDefConfig() bool {
	return true
}

// UnmarshalCertificateTemplateConfig unmarshals an instance of CertificateTemplateConfig from the specified map of raw messages.
func UnmarshalCertificateTemplateConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateTemplateConfig)
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
	err = core.UnmarshalPrimitive(m, "allow_any_name", &obj.AllowAnyName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enforce_hostnames", &obj.EnforceHostnames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_ip_sans", &obj.AllowIPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_uri_sans", &obj.AllowedURISans)
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

// CertificateTemplatesConfig : Certificate templates configuration.
// This model "extends" GetConfigElementsResourcesItem
type CertificateTemplatesConfig struct {
	CertificateTemplates []CertificateTemplatesConfigItem `json:"certificate_templates" validate:"required"`
}

func (*CertificateTemplatesConfig) isaGetConfigElementsResourcesItem() bool {
	return true
}

// UnmarshalCertificateTemplatesConfig unmarshals an instance of CertificateTemplatesConfig from the specified map of raw messages.
func UnmarshalCertificateTemplatesConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateTemplatesConfig)
	err = core.UnmarshalModel(m, "certificate_templates", &obj.CertificateTemplates, UnmarshalCertificateTemplatesConfigItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigElementDefConfigClassicInfrastructureConfig : Properties that describe an IBM Cloud classic infrastructure (SoftLayer) configuration.
// This model "extends" ConfigElementDefConfig
type ConfigElementDefConfigClassicInfrastructureConfig struct {
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

// NewConfigElementDefConfigClassicInfrastructureConfig : Instantiate ConfigElementDefConfigClassicInfrastructureConfig (Generic Model Constructor)
func (*SecretsManagerV1) NewConfigElementDefConfigClassicInfrastructureConfig(classicInfrastructureUsername string, classicInfrastructurePassword string) (_model *ConfigElementDefConfigClassicInfrastructureConfig, err error) {
	_model = &ConfigElementDefConfigClassicInfrastructureConfig{
		ClassicInfrastructureUsername: core.StringPtr(classicInfrastructureUsername),
		ClassicInfrastructurePassword: core.StringPtr(classicInfrastructurePassword),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ConfigElementDefConfigClassicInfrastructureConfig) isaConfigElementDefConfig() bool {
	return true
}

// UnmarshalConfigElementDefConfigClassicInfrastructureConfig unmarshals an instance of ConfigElementDefConfigClassicInfrastructureConfig from the specified map of raw messages.
func UnmarshalConfigElementDefConfigClassicInfrastructureConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementDefConfigClassicInfrastructureConfig)
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

// ConfigElementDefConfigCloudInternetServicesConfig : Properties that describe an IBM Cloud Internet Services (CIS) configuration.
// This model "extends" ConfigElementDefConfig
type ConfigElementDefConfigCloudInternetServicesConfig struct {
	// The Cloud Resource Name (CRN) that is associated with the CIS instance.
	CisCRN *string `json:"cis_crn" validate:"required"`

	// An IBM Cloud API key that can to list domains in your CIS instance.
	//
	// To grant Secrets Manager the ability to view the CIS instance and all of its domains, the API key must be assigned
	// the Reader service role on Internet Services (`internet-svcs`).
	//
	// If you need to manage specific domains, you can assign the Manager role. For production environments, it is
	// recommended that you assign the Reader access role, and then use the
	// [IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific
	// domains. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	CisApikey *string `json:"cis_apikey,omitempty"`
}

// NewConfigElementDefConfigCloudInternetServicesConfig : Instantiate ConfigElementDefConfigCloudInternetServicesConfig (Generic Model Constructor)
func (*SecretsManagerV1) NewConfigElementDefConfigCloudInternetServicesConfig(cisCRN string) (_model *ConfigElementDefConfigCloudInternetServicesConfig, err error) {
	_model = &ConfigElementDefConfigCloudInternetServicesConfig{
		CisCRN: core.StringPtr(cisCRN),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ConfigElementDefConfigCloudInternetServicesConfig) isaConfigElementDefConfig() bool {
	return true
}

// UnmarshalConfigElementDefConfigCloudInternetServicesConfig unmarshals an instance of ConfigElementDefConfigCloudInternetServicesConfig from the specified map of raw messages.
func UnmarshalConfigElementDefConfigCloudInternetServicesConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementDefConfigCloudInternetServicesConfig)
	err = core.UnmarshalPrimitive(m, "cis_crn", &obj.CisCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cis_apikey", &obj.CisApikey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigElementDefConfigLetsEncryptConfig : Properties that describe a Let's Encrypt configuration.
// This model "extends" ConfigElementDefConfig
type ConfigElementDefConfigLetsEncryptConfig struct {
	// The private key that is associated with your Automatic Certificate Management Environment (ACME) account.
	//
	// If you have a working ACME client or account for Let's Encrypt, you can use the existing private key to enable
	// communications with Secrets Manager. If you don't have an account yet, you can create one. For more information, see
	// the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#create-acme-account).
	PrivateKey *string `json:"private_key" validate:"required"`
}

// NewConfigElementDefConfigLetsEncryptConfig : Instantiate ConfigElementDefConfigLetsEncryptConfig (Generic Model Constructor)
func (*SecretsManagerV1) NewConfigElementDefConfigLetsEncryptConfig(privateKey string) (_model *ConfigElementDefConfigLetsEncryptConfig, err error) {
	_model = &ConfigElementDefConfigLetsEncryptConfig{
		PrivateKey: core.StringPtr(privateKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ConfigElementDefConfigLetsEncryptConfig) isaConfigElementDefConfig() bool {
	return true
}

// UnmarshalConfigElementDefConfigLetsEncryptConfig unmarshals an instance of ConfigElementDefConfigLetsEncryptConfig from the specified map of raw messages.
func UnmarshalConfigElementDefConfigLetsEncryptConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigElementDefConfigLetsEncryptConfig)
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateIamCredentialsSecretEngineRootConfig : Configuration for the IAM credentials engine.
// This model "extends" EngineConfig
type CreateIamCredentialsSecretEngineRootConfig struct {
	// An IBM Cloud API key that can create and manage service IDs.
	//
	// The API key must be assigned the Editor platform role on the Access Groups Service and the Operator platform role on
	// the IAM Identity Service. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).
	APIKey *string `json:"api_key" validate:"required"`

	// The hash value of the IBM Cloud API key that is used to create and manage service IDs.
	APIKeyHash *string `json:"api_key_hash,omitempty"`
}

// NewCreateIamCredentialsSecretEngineRootConfig : Instantiate CreateIamCredentialsSecretEngineRootConfig (Generic Model Constructor)
func (*SecretsManagerV1) NewCreateIamCredentialsSecretEngineRootConfig(apiKey string) (_model *CreateIamCredentialsSecretEngineRootConfig, err error) {
	_model = &CreateIamCredentialsSecretEngineRootConfig{
		APIKey: core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*CreateIamCredentialsSecretEngineRootConfig) isaEngineConfig() bool {
	return true
}

// UnmarshalCreateIamCredentialsSecretEngineRootConfig unmarshals an instance of CreateIamCredentialsSecretEngineRootConfig from the specified map of raw messages.
func UnmarshalCreateIamCredentialsSecretEngineRootConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateIamCredentialsSecretEngineRootConfig)
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_hash", &obj.APIKeyHash)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteCredentialsForIamCredentialsSecret : Delete the credentials that are associated with an `iam_credentials` secret.
// This model "extends" SecretAction
type DeleteCredentialsForIamCredentialsSecret struct {
	// The ID of the API key that you want to delete. If the secret was created with a static service ID, only the API key
	// is deleted. Otherwise, the service ID is deleted together with its API key.
	APIKeyID *string `json:"api_key_id,omitempty"`

	// The service ID that you want to delete. This property can be used instead of the `api_key_id` field, but only for
	// secrets that were created with a service ID that was generated by Secrets Manager.
	//
	// **Deprecated.** Use the `api_key_id` field instead.
	// Deprecated: this field is deprecated and may be removed in a future release.
	ServiceID *string `json:"service_id,omitempty"`
}

func (*DeleteCredentialsForIamCredentialsSecret) isaSecretAction() bool {
	return true
}

// UnmarshalDeleteCredentialsForIamCredentialsSecret unmarshals an instance of DeleteCredentialsForIamCredentialsSecret from the specified map of raw messages.
func UnmarshalDeleteCredentialsForIamCredentialsSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteCredentialsForIamCredentialsSecret)
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.APIKeyID)
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

// GetConfigElementsResourcesItemCertificateAuthoritiesConfig : Certificate authorities configuration.
// This model "extends" GetConfigElementsResourcesItem
type GetConfigElementsResourcesItemCertificateAuthoritiesConfig struct {
	CertificateAuthorities []ConfigElementMetadata `json:"certificate_authorities" validate:"required"`
}

func (*GetConfigElementsResourcesItemCertificateAuthoritiesConfig) isaGetConfigElementsResourcesItem() bool {
	return true
}

// UnmarshalGetConfigElementsResourcesItemCertificateAuthoritiesConfig unmarshals an instance of GetConfigElementsResourcesItemCertificateAuthoritiesConfig from the specified map of raw messages.
func UnmarshalGetConfigElementsResourcesItemCertificateAuthoritiesConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetConfigElementsResourcesItemCertificateAuthoritiesConfig)
	err = core.UnmarshalModel(m, "certificate_authorities", &obj.CertificateAuthorities, UnmarshalConfigElementMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetConfigElementsResourcesItemDNSProvidersConfig : DNS providers configuration.
// This model "extends" GetConfigElementsResourcesItem
type GetConfigElementsResourcesItemDNSProvidersConfig struct {
	DNSProviders []ConfigElementMetadata `json:"dns_providers" validate:"required"`
}

func (*GetConfigElementsResourcesItemDNSProvidersConfig) isaGetConfigElementsResourcesItem() bool {
	return true
}

// UnmarshalGetConfigElementsResourcesItemDNSProvidersConfig unmarshals an instance of GetConfigElementsResourcesItemDNSProvidersConfig from the specified map of raw messages.
func UnmarshalGetConfigElementsResourcesItemDNSProvidersConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetConfigElementsResourcesItemDNSProvidersConfig)
	err = core.UnmarshalModel(m, "dns_providers", &obj.DNSProviders, UnmarshalConfigElementMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSecretPolicyRotation : Properties that describe a rotation policy.
// This model "extends" GetSecretPolicies
type GetSecretPolicyRotation struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []map[string]interface{} `json:"resources" validate:"required"`
}

func (*GetSecretPolicyRotation) isaGetSecretPolicies() bool {
	return true
}

// UnmarshalGetSecretPolicyRotation unmarshals an instance of GetSecretPolicyRotation from the specified map of raw messages.
func UnmarshalGetSecretPolicyRotation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretPolicyRotation)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources", &obj.Resources)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IamCredentialsSecretEngineRootConfig : Configuration for the IAM credentials engine.
// This model "extends" GetConfigResourcesItem
type IamCredentialsSecretEngineRootConfig struct {
	// An IBM Cloud API key that can create and manage service IDs.
	//
	// The API key must be assigned the Editor platform role on the Access Groups Service and the Operator platform role on
	// the IAM Identity Service. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).
	APIKey *string `json:"api_key" validate:"required"`

	// The hash value of the IBM Cloud API key that is used to create and manage service IDs.
	APIKeyHash *string `json:"api_key_hash,omitempty"`
}

func (*IamCredentialsSecretEngineRootConfig) isaGetConfigResourcesItem() bool {
	return true
}

// UnmarshalIamCredentialsSecretEngineRootConfig unmarshals an instance of IamCredentialsSecretEngineRootConfig from the specified map of raw messages.
func UnmarshalIamCredentialsSecretEngineRootConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamCredentialsSecretEngineRootConfig)
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_hash", &obj.APIKeyHash)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IamCredentialsSecretMetadata : Metadata properties that describe an `iam_credentials` secret.
// This model "extends" SecretMetadata
type IamCredentialsSecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be in the range 2 - 30 characters, including spaces. Special characters
	// that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The time-to-live (TTL) or lease duration that is assigned to the secret. For `iam_credentials` secrets, the TTL
	// defines for how long each generated API key remains valid.
	TTL *string `json:"ttl,omitempty"`

	// Determines whether to use the same service ID and API key for future read operations on an
	// `iam_credentials` secret.
	//
	// If set to `true`, the service reuses the current credentials. If set to `false`, a new service ID and API key are
	// generated each time that the secret is read or accessed.
	ReuseAPIKey *bool `json:"reuse_api_key,omitempty"`

	// Indicates whether an `iam_credentials` secret was created with a static service ID.
	//
	// If the value is `true`, the service ID for the secret was provided by the user at secret creation. If the value is
	// `false`, the service ID was generated by Secrets Manager.
	ServiceIDIsStatic *bool `json:"service_id_is_static,omitempty"`

	// The service ID under which the API key is created. The service ID is included in the metadata only if the secret was
	// created with a static service ID.
	ServiceID *string `json:"service_id,omitempty"`

	// The access groups that define the capabilities of the service ID and API key that are generated for an
	// `iam_credentials` secret. The access groups are included in the metadata only if the secret was created with a
	// service ID that was generated by Secrets Manager.
	AccessGroups []string `json:"access_groups,omitempty"`
}

// Constants associated with the IamCredentialsSecretMetadata.SecretType property.
// The secret type.
const (
	IamCredentialsSecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	IamCredentialsSecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	IamCredentialsSecretMetadataSecretTypeImportedCertConst     = "imported_cert"
	IamCredentialsSecretMetadataSecretTypeKvConst               = "kv"
	IamCredentialsSecretMetadataSecretTypePrivateCertConst      = "private_cert"
	IamCredentialsSecretMetadataSecretTypePublicCertConst       = "public_cert"
	IamCredentialsSecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

// NewIamCredentialsSecretMetadata : Instantiate IamCredentialsSecretMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewIamCredentialsSecretMetadata(name string) (_model *IamCredentialsSecretMetadata, err error) {
	_model = &IamCredentialsSecretMetadata{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*IamCredentialsSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalIamCredentialsSecretMetadata unmarshals an instance of IamCredentialsSecretMetadata from the specified map of raw messages.
func UnmarshalIamCredentialsSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamCredentialsSecretMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseAPIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id_is_static", &obj.ServiceIDIsStatic)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "access_groups", &obj.AccessGroups)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IamCredentialsSecretResource : Properties that describe a secret.
// This model "extends" SecretResource
type IamCredentialsSecretResource struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// An array that contains metadata for each secret version. For more information on the metadata properties, see [Get
	// secret version metadata](#get-secret-version-metadata).
	Versions []map[string]interface{} `json:"versions,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to generated credentials.
	//
	// For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be
	// either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m`
	// or `24h`.
	//
	// Minimum duration is 1 minute. Maximum is 90 days.
	TTL interface{} `json:"ttl,omitempty"`

	// The access groups that define the capabilities of the service ID and API key that are generated for an
	// `iam_credentials` secret. If you prefer to use an existing service ID that is already assigned the access policies
	// that you require, you can omit this parameter and use the `service_id` field instead.
	//
	// **Tip:** To list the access groups that are available in an account, you can use the [IAM Access Groups
	// API](https://cloud.ibm.com/apidocs/iam-access-groups#list-access-groups). To find the ID of an access group in the
	// console, go to **Manage > Access (IAM) > Access groups**. Select the access group to inspect, and click **Details**
	// to view its ID.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease (see the `ttl` field), the API key is deleted automatically. If you
	// want to continue to use the same API key for future read operations, see the `reuse_api_key` field.
	APIKey *string `json:"api_key,omitempty"`

	// The ID of the API key that is generated for this secret.
	APIKeyID *string `json:"api_key_id,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created.
	//
	// If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation and adds it
	// to the access groups that you assign.
	//
	// Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or
	// retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include
	// the `access_groups` parameter.
	ServiceID *string `json:"service_id,omitempty"`

	// Indicates whether an `iam_credentials` secret was created with a static service ID.
	//
	// If `true`, the service ID for the secret was provided by the user at secret creation. If `false`, the service ID was
	// generated by Secrets Manager.
	ServiceIDIsStatic *bool `json:"service_id_is_static,omitempty"`

	// Determines whether to use the same service ID and API key for future read operations on an
	// `iam_credentials` secret.
	//
	// If set to `true`, the service reuses the current credentials. If set to `false`, a new service ID and API key are
	// generated each time that the secret is read or accessed.
	ReuseAPIKey *bool `json:"reuse_api_key,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that have an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`
}

// Constants associated with the IamCredentialsSecretResource.SecretType property.
// The secret type.
const (
	IamCredentialsSecretResourceSecretTypeArbitraryConst        = "arbitrary"
	IamCredentialsSecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	IamCredentialsSecretResourceSecretTypeImportedCertConst     = "imported_cert"
	IamCredentialsSecretResourceSecretTypeKvConst               = "kv"
	IamCredentialsSecretResourceSecretTypePrivateCertConst      = "private_cert"
	IamCredentialsSecretResourceSecretTypePublicCertConst       = "public_cert"
	IamCredentialsSecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// NewIamCredentialsSecretResource : Instantiate IamCredentialsSecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewIamCredentialsSecretResource(name string) (_model *IamCredentialsSecretResource, err error) {
	_model = &IamCredentialsSecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*IamCredentialsSecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalIamCredentialsSecretResource unmarshals an instance of IamCredentialsSecretResource from the specified map of raw messages.
func UnmarshalIamCredentialsSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamCredentialsSecretResource)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "access_groups", &obj.AccessGroups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.APIKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id_is_static", &obj.ServiceIDIsStatic)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseAPIKey)
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

// IamCredentialsSecretVersion : IamCredentialsSecretVersion struct
// This model "extends" SecretVersion
type IamCredentialsSecretVersion struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The data that is associated with the secret version. The data object contains the following fields:
	//
	// - `api_key`: The API key that is generated for this secret.
	// - `api_key_id`: The ID of the API key that is generated for this secret.
	// - `service_id`: The service ID under which the API key is created.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`
}

func (*IamCredentialsSecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalIamCredentialsSecretVersion unmarshals an instance of IamCredentialsSecretVersion from the specified map of raw messages.
func UnmarshalIamCredentialsSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamCredentialsSecretVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IamCredentialsSecretVersionInfo : IamCredentialsSecretVersionInfo struct
// This model "extends" SecretVersionInfo
type IamCredentialsSecretVersionInfo struct {
	// The ID of the secret version.
	ID *string `json:"id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`
}

func (*IamCredentialsSecretVersionInfo) isaSecretVersionInfo() bool {
	return true
}

// UnmarshalIamCredentialsSecretVersionInfo unmarshals an instance of IamCredentialsSecretVersionInfo from the specified map of raw messages.
func UnmarshalIamCredentialsSecretVersionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamCredentialsSecretVersionInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IamCredentialsSecretVersionMetadata : Properties that describe a secret version.
// This model "extends" SecretVersionMetadata
type IamCredentialsSecretVersionMetadata struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`
}

func (*IamCredentialsSecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalIamCredentialsSecretVersionMetadata unmarshals an instance of IamCredentialsSecretVersionMetadata from the specified map of raw messages.
func UnmarshalIamCredentialsSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamCredentialsSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IntermediateCertificateAuthoritiesConfig : Intermediate certificate authorities configuration.
// This model "extends" GetConfigElementsResourcesItem
type IntermediateCertificateAuthoritiesConfig struct {
	IntermediateCertificateAuthorities []IntermediateCertificateAuthoritiesConfigItem `json:"intermediate_certificate_authorities" validate:"required"`
}

func (*IntermediateCertificateAuthoritiesConfig) isaGetConfigElementsResourcesItem() bool {
	return true
}

// UnmarshalIntermediateCertificateAuthoritiesConfig unmarshals an instance of IntermediateCertificateAuthoritiesConfig from the specified map of raw messages.
func UnmarshalIntermediateCertificateAuthoritiesConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IntermediateCertificateAuthoritiesConfig)
	err = core.UnmarshalModel(m, "intermediate_certificate_authorities", &obj.IntermediateCertificateAuthorities, UnmarshalIntermediateCertificateAuthoritiesConfigItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IntermediateCertificateAuthorityConfig : Intermediate certificate authority configuration.
// This model "extends" ConfigElementDefConfig
type IntermediateCertificateAuthorityConfig struct {
	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL interface{} `json:"max_ttl" validate:"required"`

	// The signing method to use with this certificate authority to generate private certificates.
	//
	// You can choose between internal or externally signed options. For more information, see the
	// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	SigningMethod *string `json:"signing_method" validate:"required"`

	// The certificate authority that signed and issued the certificate.
	//
	// If the certificate is signed internally, the `issuer` field is required and must match the name of a certificate
	// authority that is configured in the Secrets Manager service instance.
	Issuer *string `json:"issuer,omitempty"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry interface{} `json:"crl_expiry,omitempty"`

	// Disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is
	// enabled,  it will rebuild the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are
	// issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this
	// certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

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

	// The number of bits to use when generating the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The data that is associated with the intermediate certificate authority. The data object contains the
	//  following fields:
	//
	// - `csr`: The PEM-encoded certificate signing request.
	// - `private_key`: The private key.
	// - `private_key_type`: The type of private key, for example `rsa`.
	Data map[string]interface{} `json:"data,omitempty"`
}

// Constants associated with the IntermediateCertificateAuthorityConfig.SigningMethod property.
// The signing method to use with this certificate authority to generate private certificates.
//
// You can choose between internal or externally signed options. For more information, see the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
const (
	IntermediateCertificateAuthorityConfigSigningMethodExternalConst = "external"
	IntermediateCertificateAuthorityConfigSigningMethodInternalConst = "internal"
)

// Constants associated with the IntermediateCertificateAuthorityConfig.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	IntermediateCertificateAuthorityConfigStatusCertificateTemplateRequiredConst = "certificate_template_required"
	IntermediateCertificateAuthorityConfigStatusConfiguredConst                  = "configured"
	IntermediateCertificateAuthorityConfigStatusExpiredConst                     = "expired"
	IntermediateCertificateAuthorityConfigStatusRevokedConst                     = "revoked"
	IntermediateCertificateAuthorityConfigStatusSignedCertificateRequiredConst   = "signed_certificate_required"
	IntermediateCertificateAuthorityConfigStatusSigningRequiredConst             = "signing_required"
)

// Constants associated with the IntermediateCertificateAuthorityConfig.Format property.
// The format of the returned data.
const (
	IntermediateCertificateAuthorityConfigFormatPemConst       = "pem"
	IntermediateCertificateAuthorityConfigFormatPemBundleConst = "pem_bundle"
)

// Constants associated with the IntermediateCertificateAuthorityConfig.PrivateKeyFormat property.
// The format of the generated private key.
const (
	IntermediateCertificateAuthorityConfigPrivateKeyFormatDerConst   = "der"
	IntermediateCertificateAuthorityConfigPrivateKeyFormatPkcs8Const = "pkcs8"
)

// Constants associated with the IntermediateCertificateAuthorityConfig.KeyType property.
// The type of private key to generate.
const (
	IntermediateCertificateAuthorityConfigKeyTypeEcConst  = "ec"
	IntermediateCertificateAuthorityConfigKeyTypeRsaConst = "rsa"
)

// NewIntermediateCertificateAuthorityConfig : Instantiate IntermediateCertificateAuthorityConfig (Generic Model Constructor)
func (*SecretsManagerV1) NewIntermediateCertificateAuthorityConfig(maxTTL interface{}, signingMethod string, commonName string) (_model *IntermediateCertificateAuthorityConfig, err error) {
	_model = &IntermediateCertificateAuthorityConfig{
		MaxTTL:        maxTTL,
		SigningMethod: core.StringPtr(signingMethod),
		CommonName:    core.StringPtr(commonName),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*IntermediateCertificateAuthorityConfig) isaConfigElementDefConfig() bool {
	return true
}

// UnmarshalIntermediateCertificateAuthorityConfig unmarshals an instance of IntermediateCertificateAuthorityConfig from the specified map of raw messages.
func UnmarshalIntermediateCertificateAuthorityConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IntermediateCertificateAuthorityConfig)
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
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KvSecretMetadata : Metadata properties that describe a key-value secret.
// This model "extends" SecretMetadata
type KvSecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be in the range 2 - 30 characters, including spaces. Special characters
	// that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`
}

// Constants associated with the KvSecretMetadata.SecretType property.
// The secret type.
const (
	KvSecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	KvSecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	KvSecretMetadataSecretTypeImportedCertConst     = "imported_cert"
	KvSecretMetadataSecretTypeKvConst               = "kv"
	KvSecretMetadataSecretTypePrivateCertConst      = "private_cert"
	KvSecretMetadataSecretTypePublicCertConst       = "public_cert"
	KvSecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

// NewKvSecretMetadata : Instantiate KvSecretMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewKvSecretMetadata(name string) (_model *KvSecretMetadata, err error) {
	_model = &KvSecretMetadata{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KvSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalKvSecretMetadata unmarshals an instance of KvSecretMetadata from the specified map of raw messages.
func UnmarshalKvSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KvSecretMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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

// KvSecretResource : Properties that describe a secret.
// This model "extends" SecretResource
type KvSecretResource struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// An array that contains metadata for each secret version. For more information on the metadata properties, see [Get
	// secret version metadata](#get-secret-version-metadata).
	Versions []map[string]interface{} `json:"versions,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The new secret data to assign to the secret.
	Payload map[string]interface{} `json:"payload,omitempty"`

	// The data that is associated with the secret version.
	//
	// The data object contains the field `payload`.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`
}

// Constants associated with the KvSecretResource.SecretType property.
// The secret type.
const (
	KvSecretResourceSecretTypeArbitraryConst        = "arbitrary"
	KvSecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	KvSecretResourceSecretTypeImportedCertConst     = "imported_cert"
	KvSecretResourceSecretTypeKvConst               = "kv"
	KvSecretResourceSecretTypePrivateCertConst      = "private_cert"
	KvSecretResourceSecretTypePublicCertConst       = "public_cert"
	KvSecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// NewKvSecretResource : Instantiate KvSecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewKvSecretResource(name string) (_model *KvSecretResource, err error) {
	_model = &KvSecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KvSecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalKvSecretResource unmarshals an instance of KvSecretResource from the specified map of raw messages.
func UnmarshalKvSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KvSecretResource)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertPolicyRotation : The `private_cert` secret rotation policy.
// This model "extends" SecretPolicyRotationRotation
type PrivateCertPolicyRotation struct {
	AutoRotate *bool `json:"auto_rotate" validate:"required"`

	// The length of the secret rotation time interval.
	Interval *int64 `json:"interval,omitempty"`

	// The units for the secret rotation time interval.
	Unit *string `json:"unit,omitempty"`
}

// Constants associated with the PrivateCertPolicyRotation.Unit property.
// The units for the secret rotation time interval.
const (
	PrivateCertPolicyRotationUnitDayConst   = "day"
	PrivateCertPolicyRotationUnitMonthConst = "month"
)

// NewPrivateCertPolicyRotation : Instantiate PrivateCertPolicyRotation (Generic Model Constructor)
func (*SecretsManagerV1) NewPrivateCertPolicyRotation(autoRotate bool) (_model *PrivateCertPolicyRotation, err error) {
	_model = &PrivateCertPolicyRotation{
		AutoRotate: core.BoolPtr(autoRotate),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertPolicyRotation) isaSecretPolicyRotationRotation() bool {
	return true
}

// UnmarshalPrivateCertPolicyRotation unmarshals an instance of PrivateCertPolicyRotation from the specified map of raw messages.
func UnmarshalPrivateCertPolicyRotation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertPolicyRotation)
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

// PrivateCertSecretEngineRootConfig : Configuration for the private certificates engine.
// This model "extends" GetConfigResourcesItem
type PrivateCertSecretEngineRootConfig struct {
	// The root certificate authority configurations that are associated with your instance.
	RootCertificateAuthorities []RootCertificateAuthoritiesConfigItem `json:"root_certificate_authorities,omitempty"`

	// The intermediate certificate authority configurations that are associated with your instance.
	IntermediateCertificateAuthorities []IntermediateCertificateAuthoritiesConfigItem `json:"intermediate_certificate_authorities,omitempty"`

	// The certificate templates that are associated with your instance.
	CertificateTemplates []CertificateTemplatesConfigItem `json:"certificate_templates,omitempty"`
}

func (*PrivateCertSecretEngineRootConfig) isaGetConfigResourcesItem() bool {
	return true
}

// UnmarshalPrivateCertSecretEngineRootConfig unmarshals an instance of PrivateCertSecretEngineRootConfig from the specified map of raw messages.
func UnmarshalPrivateCertSecretEngineRootConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertSecretEngineRootConfig)
	err = core.UnmarshalModel(m, "root_certificate_authorities", &obj.RootCertificateAuthorities, UnmarshalRootCertificateAuthoritiesConfigItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "intermediate_certificate_authorities", &obj.IntermediateCertificateAuthorities, UnmarshalIntermediateCertificateAuthoritiesConfigItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "certificate_templates", &obj.CertificateTemplates, UnmarshalCertificateTemplatesConfigItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateSecretMetadata : Metadata properties that describe a private certificate secret.
// This model "extends" SecretMetadata
type PrivateCertificateSecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be in the range 2 - 30 characters, including spaces. Special characters
	// that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template,omitempty"`

	// The intermediate certificate authority that signed this certificate.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The alternative names that are defined for the certificate.
	AltNames []string `json:"alt_names,omitempty"`

	Rotation *Rotation `json:"rotation,omitempty"`

	// The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the
	// certificate.
	Algorithm *string `json:"algorithm,omitempty"`

	// The identifier for the cryptographic algorithm that was used to generate the public and private keys that are
	// associated with the certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The certificate authority that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`
}

// Constants associated with the PrivateCertificateSecretMetadata.SecretType property.
// The secret type.
const (
	PrivateCertificateSecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	PrivateCertificateSecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	PrivateCertificateSecretMetadataSecretTypeImportedCertConst     = "imported_cert"
	PrivateCertificateSecretMetadataSecretTypeKvConst               = "kv"
	PrivateCertificateSecretMetadataSecretTypePrivateCertConst      = "private_cert"
	PrivateCertificateSecretMetadataSecretTypePublicCertConst       = "public_cert"
	PrivateCertificateSecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

// NewPrivateCertificateSecretMetadata : Instantiate PrivateCertificateSecretMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewPrivateCertificateSecretMetadata(name string) (_model *PrivateCertificateSecretMetadata, err error) {
	_model = &PrivateCertificateSecretMetadata{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalPrivateCertificateSecretMetadata unmarshals an instance of PrivateCertificateSecretMetadata from the specified map of raw messages.
func UnmarshalPrivateCertificateSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateSecretMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_template", &obj.CertificateTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
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
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
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

// PrivateCertificateSecretResource : Properties that describe a secret.
// This model "extends" SecretResource
type PrivateCertificateSecretResource struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// An array that contains metadata for each secret version. For more information on the metadata properties, see [Get
	// secret version metadata](#get-secret-version-metadata).
	Versions []map[string]interface{} `json:"versions,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The name of the certificate template.
	CertificateTemplate *string `json:"certificate_template" validate:"required"`

	// The intermediate certificate authority that signed this certificate.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// The certificate signing request. If you don't include this parameter, the CSR that is used to generate the
	// certificate is created internally. If you provide a CSR, it is used also for auto rotation and manual rotation,
	// unless you provide another CSR in the manual rotation request.
	Csr *string `json:"csr,omitempty"`

	// The fully qualified domain name or host domain name for the certificate. If you provide a CSR that includes a common
	// name value, the certificate is generated with the common name that is provided in the CSR.
	CommonName *string `json:"common_name" validate:"required"`

	// The alternative names that are defined for the certificate.
	//
	// For public certificates, this value is provided as an array of strings. For private certificates, this value is
	// provided as a comma-delimited list (string). In the API response, this value is returned as an array of strings for
	// all the types of certificate secrets.
	AltNames interface{} `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

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
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	Rotation *Rotation `json:"rotation,omitempty"`

	// The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the
	// certificate.
	Algorithm *string `json:"algorithm,omitempty"`

	// The identifier for the cryptographic algorithm that was used to generate the public and private keys that are
	// associated with the certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The certificate authority that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`

	// The data that is associated with the secret. The data object contains the following fields:
	//
	// - `certificate`: The contents of the certificate.
	// - `private_key`: The private key that is associated with the certificate. If you provide a CSR in the request, the
	// private_key field is not included in the data.
	// - `issuing_ca`: The certificate of the certificate authority that signed and issued this certificate.
	// - `ca_chain`: The chain of certificate authorities that are associated with the certificate.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`
}

// Constants associated with the PrivateCertificateSecretResource.SecretType property.
// The secret type.
const (
	PrivateCertificateSecretResourceSecretTypeArbitraryConst        = "arbitrary"
	PrivateCertificateSecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	PrivateCertificateSecretResourceSecretTypeImportedCertConst     = "imported_cert"
	PrivateCertificateSecretResourceSecretTypeKvConst               = "kv"
	PrivateCertificateSecretResourceSecretTypePrivateCertConst      = "private_cert"
	PrivateCertificateSecretResourceSecretTypePublicCertConst       = "public_cert"
	PrivateCertificateSecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the PrivateCertificateSecretResource.Format property.
// The format of the returned data.
const (
	PrivateCertificateSecretResourceFormatPemConst       = "pem"
	PrivateCertificateSecretResourceFormatPemBundleConst = "pem_bundle"
)

// Constants associated with the PrivateCertificateSecretResource.PrivateKeyFormat property.
// The format of the generated private key.
const (
	PrivateCertificateSecretResourcePrivateKeyFormatDerConst   = "der"
	PrivateCertificateSecretResourcePrivateKeyFormatPkcs8Const = "pkcs8"
)

// NewPrivateCertificateSecretResource : Instantiate PrivateCertificateSecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewPrivateCertificateSecretResource(name string, certificateTemplate string, commonName string) (_model *PrivateCertificateSecretResource, err error) {
	_model = &PrivateCertificateSecretResource{
		Name:                core.StringPtr(name),
		CertificateTemplate: core.StringPtr(certificateTemplate),
		CommonName:          core.StringPtr(commonName),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PrivateCertificateSecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalPrivateCertificateSecretResource unmarshals an instance of PrivateCertificateSecretResource from the specified map of raw messages.
func UnmarshalPrivateCertificateSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateSecretResource)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "certificate_template", &obj.CertificateTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
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
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "exclude_cn_from_sans", &obj.ExcludeCnFromSans)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_rfc3339", &obj.RevocationTimeRfc3339)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateSecretVersion : PrivateCertificateSecretVersion struct
// This model "extends" SecretVersion
type PrivateCertificateSecretVersion struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The data that is associated with the secret version. The data object contains the following fields:
	//
	// - `certificate`: The contents of the certificate.
	// - `private_key`: The private key that is associated with the certificate.
	// - `intermediate`: The intermediate certificate that is associated with the certificate.
	SecretData *CertificateSecretData `json:"secret_data,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`
}

func (*PrivateCertificateSecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalPrivateCertificateSecretVersion unmarshals an instance of PrivateCertificateSecretVersion from the specified map of raw messages.
func UnmarshalPrivateCertificateSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateSecretVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secret_data", &obj.SecretData, UnmarshalCertificateSecretData)
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
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_rfc3339", &obj.RevocationTimeRfc3339)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateSecretVersionInfo : PrivateCertificateSecretVersionInfo struct
// This model "extends" SecretVersionInfo
type PrivateCertificateSecretVersionInfo struct {
	// The ID of the secret version.
	ID *string `json:"id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`
}

func (*PrivateCertificateSecretVersionInfo) isaSecretVersionInfo() bool {
	return true
}

// UnmarshalPrivateCertificateSecretVersionInfo unmarshals an instance of PrivateCertificateSecretVersionInfo from the specified map of raw messages.
func UnmarshalPrivateCertificateSecretVersionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateSecretVersionInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
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
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_rfc3339", &obj.RevocationTimeRfc3339)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrivateCertificateSecretVersionMetadata : Properties that describe a secret version.
// This model "extends" SecretVersionMetadata
type PrivateCertificateSecretVersionMetadata struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The timestamp of the certificate revocation.
	RevocationTime *int64 `json:"revocation_time,omitempty"`

	// The date and time that the certificate was revoked. The date format follows RFC 3339.
	RevocationTimeRfc3339 *strfmt.DateTime `json:"revocation_time_rfc3339,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`
}

func (*PrivateCertificateSecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalPrivateCertificateSecretVersionMetadata unmarshals an instance of PrivateCertificateSecretVersionMetadata from the specified map of raw messages.
func UnmarshalPrivateCertificateSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrivateCertificateSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
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
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revocation_time_rfc3339", &obj.RevocationTimeRfc3339)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertSecretEngineRootConfig : Configuration for the public certificates engine.
// This model "extends" GetConfigResourcesItem
type PublicCertSecretEngineRootConfig struct {
	// The certificate authority configurations that are associated with your instance.
	CertificateAuthorities []ConfigElementMetadata `json:"certificate_authorities,omitempty"`

	// The DNS provider configurations that are associated with your instance.
	DNSProviders []ConfigElementMetadata `json:"dns_providers,omitempty"`
}

func (*PublicCertSecretEngineRootConfig) isaGetConfigResourcesItem() bool {
	return true
}

// UnmarshalPublicCertSecretEngineRootConfig unmarshals an instance of PublicCertSecretEngineRootConfig from the specified map of raw messages.
func UnmarshalPublicCertSecretEngineRootConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertSecretEngineRootConfig)
	err = core.UnmarshalModel(m, "certificate_authorities", &obj.CertificateAuthorities, UnmarshalConfigElementMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "dns_providers", &obj.DNSProviders, UnmarshalConfigElementMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PublicCertificateSecretMetadata : Metadata properties that describe a public certificate secret.
// This model "extends" SecretMetadata
type PublicCertificateSecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be in the range 2 - 30 characters, including spaces. Special characters
	// that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// Determines whether your issued certificate is bundled with intermediate certificates.
	//
	// Set to `false` for the certificate file to contain only the issued certificate.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The identifier for the cryptographic algorithm to be used by the issuing certificate authority to sign the
	// certificate.
	Algorithm *string `json:"algorithm,omitempty"`

	// The identifier for the cryptographic algorithm to be used to generate the public key that is associated with the
	// certificate.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The alternative names that are defined for the certificate.
	AltNames []string `json:"alt_names,omitempty"`

	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// Indicates whether the certificate was ordered with an associated intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included,omitempty"`

	// Indicates whether the certificate was ordered with an associated private key.
	PrivateKeyIncluded *bool `json:"private_key_included,omitempty"`

	Rotation *Rotation `json:"rotation,omitempty"`

	// Issuance information that is associated with your certificate.
	IssuanceInfo *IssuanceInfo `json:"issuance_info,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`
}

// Constants associated with the PublicCertificateSecretMetadata.SecretType property.
// The secret type.
const (
	PublicCertificateSecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	PublicCertificateSecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	PublicCertificateSecretMetadataSecretTypeImportedCertConst     = "imported_cert"
	PublicCertificateSecretMetadataSecretTypeKvConst               = "kv"
	PublicCertificateSecretMetadataSecretTypePrivateCertConst      = "private_cert"
	PublicCertificateSecretMetadataSecretTypePublicCertConst       = "public_cert"
	PublicCertificateSecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the PublicCertificateSecretMetadata.KeyAlgorithm property.
// The identifier for the cryptographic algorithm to be used to generate the public key that is associated with the
// certificate.
const (
	PublicCertificateSecretMetadataKeyAlgorithmEc256Const   = "EC256"
	PublicCertificateSecretMetadataKeyAlgorithmEc384Const   = "EC384"
	PublicCertificateSecretMetadataKeyAlgorithmRsa2048Const = "RSA2048"
	PublicCertificateSecretMetadataKeyAlgorithmRsa4096Const = "RSA4096"
)

// NewPublicCertificateSecretMetadata : Instantiate PublicCertificateSecretMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewPublicCertificateSecretMetadata(name string) (_model *PublicCertificateSecretMetadata, err error) {
	_model = &PublicCertificateSecretMetadata{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalPublicCertificateSecretMetadata unmarshals an instance of PublicCertificateSecretMetadata from the specified map of raw messages.
func UnmarshalPublicCertificateSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateSecretMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_metadata", &obj.CustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bundle_certs", &obj.BundleCerts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
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
	err = core.UnmarshalPrimitive(m, "intermediate_included", &obj.IntermediateIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key_included", &obj.PrivateKeyIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotation)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "issuance_info", &obj.IssuanceInfo, UnmarshalIssuanceInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
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

// PublicCertificateSecretResource : Properties that describe a secret.
// This model "extends" SecretResource
type PublicCertificateSecretResource struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// An array that contains metadata for each secret version. For more information on the metadata properties, see [Get
	// secret version metadata](#get-secret-version-metadata).
	Versions []map[string]interface{} `json:"versions,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The distinguished name that identifies the entity that signed and issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// Determines whether your issued certificate is bundled with intermediate certificates.
	//
	// Set to `false` for the certificate file to contain only the issued certificate.
	BundleCerts *bool `json:"bundle_certs,omitempty"`

	// The name of the certificate authority configuration.
	//
	// To view a list of your configured authorities, use the [List configurations API](#get-secret-config-element).
	Ca *string `json:"ca,omitempty"`

	// The name of the DNS provider configuration.
	//
	// To view a list of your configured authorities, use the [List configurations API](#get-secret-config-element).
	DNS *string `json:"dns,omitempty"`

	// The identifier for the cryptographic algorithm to be used by the issuing certificate authority to sign the
	// certificate.
	Algorithm *string `json:"algorithm,omitempty"`

	// The identifier for the cryptographic algorithm to be used to generate the public key that is associated with the
	// certificate.
	//
	// The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to
	// generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide
	// more encryption protection.
	KeyAlgorithm *string `json:"key_algorithm,omitempty"`

	// The alternative names that are defined for the certificate.
	//
	// For public certificates, this value is provided as an array of strings. For private certificates, this value is
	// provided as a comma-delimited list (string). In the API response, this value is returned as an array of strings for
	// all the types of certificate secrets.
	AltNames interface{} `json:"alt_names,omitempty"`

	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// Indicates whether the issued certificate includes a private key.
	PrivateKeyIncluded *bool `json:"private_key_included,omitempty"`

	// Indicates whether the issued certificate includes an intermediate certificate.
	IntermediateIncluded *bool `json:"intermediate_included,omitempty"`

	Rotation *Rotation `json:"rotation,omitempty"`

	// Issuance information that is associated with your certificate.
	IssuanceInfo *IssuanceInfo `json:"issuance_info,omitempty"`

	Validity *CertificateValidity `json:"validity,omitempty"`

	// The unique serial number that was assigned to the certificate by the issuing certificate authority.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The data that is associated with the secret. The data object contains the following fields:
	//
	// - `certificate`: The contents of the certificate.
	// - `private_key`: The private key that is associated with the certificate.
	// - `intermediate`: The intermediate certificate that is associated with the certificate.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`
}

// Constants associated with the PublicCertificateSecretResource.SecretType property.
// The secret type.
const (
	PublicCertificateSecretResourceSecretTypeArbitraryConst        = "arbitrary"
	PublicCertificateSecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	PublicCertificateSecretResourceSecretTypeImportedCertConst     = "imported_cert"
	PublicCertificateSecretResourceSecretTypeKvConst               = "kv"
	PublicCertificateSecretResourceSecretTypePrivateCertConst      = "private_cert"
	PublicCertificateSecretResourceSecretTypePublicCertConst       = "public_cert"
	PublicCertificateSecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the PublicCertificateSecretResource.KeyAlgorithm property.
// The identifier for the cryptographic algorithm to be used to generate the public key that is associated with the
// certificate.
//
// The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to
// generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide
// more encryption protection.
const (
	PublicCertificateSecretResourceKeyAlgorithmEc256Const   = "EC256"
	PublicCertificateSecretResourceKeyAlgorithmEc384Const   = "EC384"
	PublicCertificateSecretResourceKeyAlgorithmRsa2048Const = "RSA2048"
	PublicCertificateSecretResourceKeyAlgorithmRsa4096Const = "RSA4096"
)

// NewPublicCertificateSecretResource : Instantiate PublicCertificateSecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewPublicCertificateSecretResource(name string) (_model *PublicCertificateSecretResource, err error) {
	_model = &PublicCertificateSecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PublicCertificateSecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalPublicCertificateSecretResource unmarshals an instance of PublicCertificateSecretResource from the specified map of raw messages.
func UnmarshalPublicCertificateSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublicCertificateSecretResource)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
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
	err = core.UnmarshalPrimitive(m, "dns", &obj.DNS)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_algorithm", &obj.KeyAlgorithm)
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
	err = core.UnmarshalPrimitive(m, "private_key_included", &obj.PrivateKeyIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate_included", &obj.IntermediateIncluded)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rotation", &obj.Rotation, UnmarshalRotation)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "issuance_info", &obj.IssuanceInfo, UnmarshalIssuanceInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validity", &obj.Validity, UnmarshalCertificateValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RestoreIamCredentialsSecretBody : The request body of a `restore` action.
// This model "extends" SecretAction
type RestoreIamCredentialsSecretBody struct {
	// The ID of the target version or the alias `previous`.
	VersionID *string `json:"version_id" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewRestoreIamCredentialsSecretBody : Instantiate RestoreIamCredentialsSecretBody (Generic Model Constructor)
func (*SecretsManagerV1) NewRestoreIamCredentialsSecretBody(versionID string) (_model *RestoreIamCredentialsSecretBody, err error) {
	_model = &RestoreIamCredentialsSecretBody{
		VersionID: core.StringPtr(versionID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RestoreIamCredentialsSecretBody) isaSecretAction() bool {
	return true
}

// UnmarshalRestoreIamCredentialsSecretBody unmarshals an instance of RestoreIamCredentialsSecretBody from the specified map of raw messages.
func UnmarshalRestoreIamCredentialsSecretBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RestoreIamCredentialsSecretBody)
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
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

// RevokeAction : A request to revoke the certificate of an internally signed intermediate certificate authority.
// This model "extends" ConfigAction
type RevokeAction struct {
	// The serial number of the certificate.
	SerialNumber *string `json:"serial_number" validate:"required"`
}

// NewRevokeAction : Instantiate RevokeAction (Generic Model Constructor)
func (*SecretsManagerV1) NewRevokeAction(serialNumber string) (_model *RevokeAction, err error) {
	_model = &RevokeAction{
		SerialNumber: core.StringPtr(serialNumber),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RevokeAction) isaConfigAction() bool {
	return true
}

// UnmarshalRevokeAction unmarshals an instance of RevokeAction from the specified map of raw messages.
func UnmarshalRevokeAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RevokeAction)
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RevokeActionResult : Properties that are returned with a successful `revoke` action.
// This model "extends" ConfigElementActionResultConfig
type RevokeActionResult struct {
	// The time until the certificate authority is revoked.
	RevocationTime *int64 `json:"revocation_time,omitempty"`
}

func (*RevokeActionResult) isaConfigElementActionResultConfig() bool {
	return true
}

// UnmarshalRevokeActionResult unmarshals an instance of RevokeActionResult from the specified map of raw messages.
func UnmarshalRevokeActionResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RevokeActionResult)
	err = core.UnmarshalPrimitive(m, "revocation_time", &obj.RevocationTime)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RootCertificateAuthoritiesConfig : Root certificate authorities configuration.
// This model "extends" GetConfigElementsResourcesItem
type RootCertificateAuthoritiesConfig struct {
	RootCertificateAuthorities []RootCertificateAuthoritiesConfigItem `json:"root_certificate_authorities" validate:"required"`
}

func (*RootCertificateAuthoritiesConfig) isaGetConfigElementsResourcesItem() bool {
	return true
}

// UnmarshalRootCertificateAuthoritiesConfig unmarshals an instance of RootCertificateAuthoritiesConfig from the specified map of raw messages.
func UnmarshalRootCertificateAuthoritiesConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RootCertificateAuthoritiesConfig)
	err = core.UnmarshalModel(m, "root_certificate_authorities", &obj.RootCertificateAuthorities, UnmarshalRootCertificateAuthoritiesConfigItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RootCertificateAuthorityConfig : Root certificate authority configuration.
// This model "extends" ConfigElementDefConfig
type RootCertificateAuthorityConfig struct {
	// The maximum time-to-live (TTL) for certificates that are created by this CA.
	//
	// The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API
	// response, this value is returned in seconds (integer).
	//
	// Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	MaxTTL interface{} `json:"max_ttl" validate:"required"`

	// The time until the certificate revocation list (CRL) expires.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours.
	// In the API response, this value is returned in seconds (integer).
	//
	// **Note:** The CRL is rotated automatically before it expires.
	CrlExpiry interface{} `json:"crl_expiry,omitempty"`

	// Disables or enables certificate revocation list (CRL) building.
	//
	// If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is
	// enabled,  it will rebuild the CRL.
	CrlDisable *bool `json:"crl_disable,omitempty"`

	// Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are
	// issued by this certificate authority.
	CrlDistributionPointsEncoded *bool `json:"crl_distribution_points_encoded,omitempty"`

	// Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this
	// certificate authority.
	IssuingCertificatesUrlsEncoded *bool `json:"issuing_certificates_urls_encoded,omitempty"`

	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name" validate:"required"`

	// The status of the certificate authority. The status of a root certificate authority is either `configured` or
	// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
	// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	Status *string `json:"status,omitempty"`

	// The date that the certificate expires. The date format follows RFC 3339.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to this CA certificate.
	//
	// The value can be supplied as a string representation of a duration, such as `12h`. The value can be supplied in
	// seconds (suffix `s`), minutes (suffix `m`), hours (suffix `h`) or days (suffix `d`). The value can't exceed the
	// `max_ttl` that is defined in the associated certificate template. In the API response, this value is returned in
	// seconds (integer).
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The format of the generated private key.
	PrivateKeyFormat *string `json:"private_key_format,omitempty"`

	// The type of private key to generate.
	KeyType *string `json:"key_type,omitempty"`

	// The number of bits to use when generating the private key.
	//
	// Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and
	// `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	KeyBits *int64 `json:"key_bits,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	PermittedDNSDomains []string `json:"permitted_dns_domains,omitempty"`

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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The data that is associated with the root certificate authority. The data object contains the following fields:
	//
	// - `certificate`: The root certificate content.
	// - `issuing_ca`: The certificate of the certificate authority that signed and issued this certificate.
	// - `serial_number`: The unique serial number of the root certificate.
	Data map[string]interface{} `json:"data,omitempty"`
}

// Constants associated with the RootCertificateAuthorityConfig.Status property.
// The status of the certificate authority. The status of a root certificate authority is either `configured` or
// `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,
// `signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
const (
	RootCertificateAuthorityConfigStatusCertificateTemplateRequiredConst = "certificate_template_required"
	RootCertificateAuthorityConfigStatusConfiguredConst                  = "configured"
	RootCertificateAuthorityConfigStatusExpiredConst                     = "expired"
	RootCertificateAuthorityConfigStatusRevokedConst                     = "revoked"
	RootCertificateAuthorityConfigStatusSignedCertificateRequiredConst   = "signed_certificate_required"
	RootCertificateAuthorityConfigStatusSigningRequiredConst             = "signing_required"
)

// Constants associated with the RootCertificateAuthorityConfig.Format property.
// The format of the returned data.
const (
	RootCertificateAuthorityConfigFormatPemConst       = "pem"
	RootCertificateAuthorityConfigFormatPemBundleConst = "pem_bundle"
)

// Constants associated with the RootCertificateAuthorityConfig.PrivateKeyFormat property.
// The format of the generated private key.
const (
	RootCertificateAuthorityConfigPrivateKeyFormatDerConst   = "der"
	RootCertificateAuthorityConfigPrivateKeyFormatPkcs8Const = "pkcs8"
)

// Constants associated with the RootCertificateAuthorityConfig.KeyType property.
// The type of private key to generate.
const (
	RootCertificateAuthorityConfigKeyTypeEcConst  = "ec"
	RootCertificateAuthorityConfigKeyTypeRsaConst = "rsa"
)

// NewRootCertificateAuthorityConfig : Instantiate RootCertificateAuthorityConfig (Generic Model Constructor)
func (*SecretsManagerV1) NewRootCertificateAuthorityConfig(maxTTL interface{}, commonName string) (_model *RootCertificateAuthorityConfig, err error) {
	_model = &RootCertificateAuthorityConfig{
		MaxTTL:     maxTTL,
		CommonName: core.StringPtr(commonName),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RootCertificateAuthorityConfig) isaConfigElementDefConfig() bool {
	return true
}

// UnmarshalRootCertificateAuthorityConfig unmarshals an instance of RootCertificateAuthorityConfig from the specified map of raw messages.
func UnmarshalRootCertificateAuthorityConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RootCertificateAuthorityConfig)
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
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDNSDomains)
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
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RotateArbitrarySecretBody : The request body of a `rotate` action.
// This model "extends" SecretAction
type RotateArbitrarySecretBody struct {
	// The new secret data to assign to an `arbitrary` secret.
	Payload interface{} `json:"payload" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewRotateArbitrarySecretBody : Instantiate RotateArbitrarySecretBody (Generic Model Constructor)
func (*SecretsManagerV1) NewRotateArbitrarySecretBody(payload string) (_model *RotateArbitrarySecretBody, err error) {
	_model = &RotateArbitrarySecretBody{
		Payload: core.StringPtr(payload),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RotateArbitrarySecretBody) isaSecretAction() bool {
	return true
}

// UnmarshalRotateArbitrarySecretBody unmarshals an instance of RotateArbitrarySecretBody from the specified map of raw messages.
func UnmarshalRotateArbitrarySecretBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotateArbitrarySecretBody)
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

// RotateCertificateBody : The request body of a rotate certificate action.
// This model "extends" SecretAction
type RotateCertificateBody struct {
	// The new data to associate with the certificate.
	Certificate *string `json:"certificate" validate:"required"`

	// The new private key to associate with the certificate.
	PrivateKey *string `json:"private_key,omitempty"`

	// The new intermediate certificate to associate with the certificate.
	Intermediate *string `json:"intermediate,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewRotateCertificateBody : Instantiate RotateCertificateBody (Generic Model Constructor)
func (*SecretsManagerV1) NewRotateCertificateBody(certificate string) (_model *RotateCertificateBody, err error) {
	_model = &RotateCertificateBody{
		Certificate: core.StringPtr(certificate),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RotateCertificateBody) isaSecretAction() bool {
	return true
}

// UnmarshalRotateCertificateBody unmarshals an instance of RotateCertificateBody from the specified map of raw messages.
func UnmarshalRotateCertificateBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotateCertificateBody)
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "intermediate", &obj.Intermediate)
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

// RotateCrlActionResult : Properties that are returned with a successful `rotate_crl` action.
// This model "extends" ConfigElementActionResultConfig
type RotateCrlActionResult struct {
}

func (*RotateCrlActionResult) isaConfigElementActionResultConfig() bool {
	return true
}

// UnmarshalRotateCrlActionResult unmarshals an instance of RotateCrlActionResult from the specified map of raw messages.
func UnmarshalRotateCrlActionResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotateCrlActionResult)
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RotateKvSecretBody : The request body of a `rotate` action.
// This model "extends" SecretAction
type RotateKvSecretBody struct {
	// The new secret data to assign to a key-value secret.
	Payload map[string]interface{} `json:"payload" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewRotateKvSecretBody : Instantiate RotateKvSecretBody (Generic Model Constructor)
func (*SecretsManagerV1) NewRotateKvSecretBody(payload map[string]interface{}) (_model *RotateKvSecretBody, err error) {
	_model = &RotateKvSecretBody{
		Payload: payload,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RotateKvSecretBody) isaSecretAction() bool {
	return true
}

// UnmarshalRotateKvSecretBody unmarshals an instance of RotateKvSecretBody from the specified map of raw messages.
func UnmarshalRotateKvSecretBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotateKvSecretBody)
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

// RotatePrivateCertBody : The request body of a rotate private certificate action.
// This model "extends" SecretAction
type RotatePrivateCertBody struct {
	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata" validate:"required"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewRotatePrivateCertBody : Instantiate RotatePrivateCertBody (Generic Model Constructor)
func (*SecretsManagerV1) NewRotatePrivateCertBody(customMetadata map[string]interface{}) (_model *RotatePrivateCertBody, err error) {
	_model = &RotatePrivateCertBody{
		CustomMetadata: customMetadata,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RotatePrivateCertBody) isaSecretAction() bool {
	return true
}

// UnmarshalRotatePrivateCertBody unmarshals an instance of RotatePrivateCertBody from the specified map of raw messages.
func UnmarshalRotatePrivateCertBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotatePrivateCertBody)
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

// RotatePrivateCertBodyWithCsr : The body of a request to rotate a private certificate.
// This model "extends" SecretAction
type RotatePrivateCertBodyWithCsr struct {
	// The certificate signing request. If you provide a CSR, it is used for auto rotation and manual rotation requests
	// that do not include a CSR. If you don't include the CSR, the certificate is generated with the last CSR that you
	// provided to create the private certificate, or on a previous request to rotate the certificate. If no CSR was
	// provided in the past, the certificate is generated with a CSR that is created internally.
	Csr *string `json:"csr" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewRotatePrivateCertBodyWithCsr : Instantiate RotatePrivateCertBodyWithCsr (Generic Model Constructor)
func (*SecretsManagerV1) NewRotatePrivateCertBodyWithCsr(csr string) (_model *RotatePrivateCertBodyWithCsr, err error) {
	_model = &RotatePrivateCertBodyWithCsr{
		Csr: core.StringPtr(csr),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RotatePrivateCertBodyWithCsr) isaSecretAction() bool {
	return true
}

// UnmarshalRotatePrivateCertBodyWithCsr unmarshals an instance of RotatePrivateCertBodyWithCsr from the specified map of raw messages.
func UnmarshalRotatePrivateCertBodyWithCsr(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotatePrivateCertBodyWithCsr)
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
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

// RotatePrivateCertBodyWithVersionCustomMetadata : The request body of a rotate private certificate action.
// This model "extends" SecretAction
type RotatePrivateCertBodyWithVersionCustomMetadata struct {
	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata" validate:"required"`
}

// NewRotatePrivateCertBodyWithVersionCustomMetadata : Instantiate RotatePrivateCertBodyWithVersionCustomMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewRotatePrivateCertBodyWithVersionCustomMetadata(versionCustomMetadata map[string]interface{}) (_model *RotatePrivateCertBodyWithVersionCustomMetadata, err error) {
	_model = &RotatePrivateCertBodyWithVersionCustomMetadata{
		VersionCustomMetadata: versionCustomMetadata,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RotatePrivateCertBodyWithVersionCustomMetadata) isaSecretAction() bool {
	return true
}

// UnmarshalRotatePrivateCertBodyWithVersionCustomMetadata unmarshals an instance of RotatePrivateCertBodyWithVersionCustomMetadata from the specified map of raw messages.
func UnmarshalRotatePrivateCertBodyWithVersionCustomMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotatePrivateCertBodyWithVersionCustomMetadata)
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RotatePublicCertBody : The request body of a `rotate` action.
// This model "extends" SecretAction
type RotatePublicCertBody struct {
	// Determine whether keys must be rotated.
	RotateKeys *bool `json:"rotate_keys" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewRotatePublicCertBody : Instantiate RotatePublicCertBody (Generic Model Constructor)
func (*SecretsManagerV1) NewRotatePublicCertBody(rotateKeys bool) (_model *RotatePublicCertBody, err error) {
	_model = &RotatePublicCertBody{
		RotateKeys: core.BoolPtr(rotateKeys),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RotatePublicCertBody) isaSecretAction() bool {
	return true
}

// UnmarshalRotatePublicCertBody unmarshals an instance of RotatePublicCertBody from the specified map of raw messages.
func UnmarshalRotatePublicCertBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotatePublicCertBody)
	err = core.UnmarshalPrimitive(m, "rotate_keys", &obj.RotateKeys)
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

// RotateUsernamePasswordSecretBody : The request body of a `rotate` action.
// This model "extends" SecretAction
type RotateUsernamePasswordSecretBody struct {
	// The new password to assign to a `username_password` secret.
	Password *string `json:"password" validate:"required"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`
}

// NewRotateUsernamePasswordSecretBody : Instantiate RotateUsernamePasswordSecretBody (Generic Model Constructor)
func (*SecretsManagerV1) NewRotateUsernamePasswordSecretBody(password string) (_model *RotateUsernamePasswordSecretBody, err error) {
	_model = &RotateUsernamePasswordSecretBody{
		Password: core.StringPtr(password),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RotateUsernamePasswordSecretBody) isaSecretAction() bool {
	return true
}

// UnmarshalRotateUsernamePasswordSecretBody unmarshals an instance of RotateUsernamePasswordSecretBody from the specified map of raw messages.
func UnmarshalRotateUsernamePasswordSecretBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RotateUsernamePasswordSecretBody)
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

// SecretPolicyRotationRotationPolicyRotation : The secret rotation time interval.
// This model "extends" SecretPolicyRotationRotation
type SecretPolicyRotationRotationPolicyRotation struct {
	// The length of the secret rotation time interval.
	Interval *int64 `json:"interval" validate:"required"`

	// The units for the secret rotation time interval.
	Unit *string `json:"unit" validate:"required"`
}

// Constants associated with the SecretPolicyRotationRotationPolicyRotation.Unit property.
// The units for the secret rotation time interval.
const (
	SecretPolicyRotationRotationPolicyRotationUnitDayConst   = "day"
	SecretPolicyRotationRotationPolicyRotationUnitMonthConst = "month"
)

// NewSecretPolicyRotationRotationPolicyRotation : Instantiate SecretPolicyRotationRotationPolicyRotation (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretPolicyRotationRotationPolicyRotation(interval int64, unit string) (_model *SecretPolicyRotationRotationPolicyRotation, err error) {
	_model = &SecretPolicyRotationRotationPolicyRotation{
		Interval: core.Int64Ptr(interval),
		Unit:     core.StringPtr(unit),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*SecretPolicyRotationRotationPolicyRotation) isaSecretPolicyRotationRotation() bool {
	return true
}

// UnmarshalSecretPolicyRotationRotationPolicyRotation unmarshals an instance of SecretPolicyRotationRotationPolicyRotation from the specified map of raw messages.
func UnmarshalSecretPolicyRotationRotationPolicyRotation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretPolicyRotationRotationPolicyRotation)
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

// SecretPolicyRotationRotationPublicCertPolicyRotation : The `public_cert` secret rotation policy.
// This model "extends" SecretPolicyRotationRotation
type SecretPolicyRotationRotationPublicCertPolicyRotation struct {
	AutoRotate *bool `json:"auto_rotate" validate:"required"`

	RotateKeys *bool `json:"rotate_keys" validate:"required"`
}

// NewSecretPolicyRotationRotationPublicCertPolicyRotation : Instantiate SecretPolicyRotationRotationPublicCertPolicyRotation (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretPolicyRotationRotationPublicCertPolicyRotation(autoRotate bool, rotateKeys bool) (_model *SecretPolicyRotationRotationPublicCertPolicyRotation, err error) {
	_model = &SecretPolicyRotationRotationPublicCertPolicyRotation{
		AutoRotate: core.BoolPtr(autoRotate),
		RotateKeys: core.BoolPtr(rotateKeys),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*SecretPolicyRotationRotationPublicCertPolicyRotation) isaSecretPolicyRotationRotation() bool {
	return true
}

// UnmarshalSecretPolicyRotationRotationPublicCertPolicyRotation unmarshals an instance of SecretPolicyRotationRotationPublicCertPolicyRotation from the specified map of raw messages.
func UnmarshalSecretPolicyRotationRotationPublicCertPolicyRotation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretPolicyRotationRotationPublicCertPolicyRotation)
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

// SetSignedAction : A request to set a signed certificate in an intermediate certificate authority.
// This model "extends" ConfigAction
type SetSignedAction struct {
	// The PEM-encoded certificate.
	Certificate *string `json:"certificate" validate:"required"`
}

// NewSetSignedAction : Instantiate SetSignedAction (Generic Model Constructor)
func (*SecretsManagerV1) NewSetSignedAction(certificate string) (_model *SetSignedAction, err error) {
	_model = &SetSignedAction{
		Certificate: core.StringPtr(certificate),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*SetSignedAction) isaConfigAction() bool {
	return true
}

// UnmarshalSetSignedAction unmarshals an instance of SetSignedAction from the specified map of raw messages.
func UnmarshalSetSignedAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SetSignedAction)
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SetSignedActionResult : Properties that are returned with a successful `set_signed` action.
// This model "extends" ConfigElementActionResultConfig
type SetSignedActionResult struct {
}

func (*SetSignedActionResult) isaConfigElementActionResultConfig() bool {
	return true
}

// UnmarshalSetSignedActionResult unmarshals an instance of SetSignedActionResult from the specified map of raw messages.
func UnmarshalSetSignedActionResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SetSignedActionResult)
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SignCsrAction : A request to sign a certificate signing request (CSR).
// This model "extends" ConfigAction
type SignCsrAction struct {
	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `12h`. The value can't exceed
	// the `max_ttl` that is defined in the associated certificate template.
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	PermittedDNSDomains []string `json:"permitted_dns_domains,omitempty"`

	// Determines whether to use values from a certificate signing request (CSR) to complete a `sign_csr` action. If set to
	// `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than using the values
	// provided in the other parameters to this operation.
	//
	// 2) Any key usages (for example, non-repudiation) that are requested in the CSR are added to the basic set of key
	// usages used for CA certs signed by this intermediate authority.
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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The PEM-encoded certificate signing request (CSR). This field is required for the `sign_csr` action.
	Csr *string `json:"csr" validate:"required"`
}

// Constants associated with the SignCsrAction.Format property.
// The format of the returned data.
const (
	SignCsrActionFormatPemConst       = "pem"
	SignCsrActionFormatPemBundleConst = "pem_bundle"
)

// NewSignCsrAction : Instantiate SignCsrAction (Generic Model Constructor)
func (*SecretsManagerV1) NewSignCsrAction(csr string) (_model *SignCsrAction, err error) {
	_model = &SignCsrAction{
		Csr: core.StringPtr(csr),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*SignCsrAction) isaConfigAction() bool {
	return true
}

// UnmarshalSignCsrAction unmarshals an instance of SignCsrAction from the specified map of raw messages.
func UnmarshalSignCsrAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SignCsrAction)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDNSDomains)
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
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SignCsrActionResult : Properties that are returned with a successful `sign_csr` action.
// This model "extends" ConfigElementActionResultConfig
type SignCsrActionResult struct {
	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `12h`. The value can't exceed
	// the `max_ttl` that is defined in the associated certificate template.
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	PermittedDNSDomains []string `json:"permitted_dns_domains,omitempty"`

	// Determines whether to use values from a certificate signing request (CSR) to complete a `sign_csr` action. If set to
	// `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than using the values
	// provided in the other parameters to this operation.
	//
	// 2) Any key usages (for example, non-repudiation) that are requested in the CSR are added to the basic set of key
	// usages used for CA certs signed by this intermediate authority.
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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// Properties that are returned with a successful `sign` action.
	Data *SignActionResultData `json:"data,omitempty"`

	// The PEM-encoded certificate signing request (CSR).
	Csr *string `json:"csr" validate:"required"`
}

// Constants associated with the SignCsrActionResult.Format property.
// The format of the returned data.
const (
	SignCsrActionResultFormatPemConst       = "pem"
	SignCsrActionResultFormatPemBundleConst = "pem_bundle"
)

func (*SignCsrActionResult) isaConfigElementActionResultConfig() bool {
	return true
}

// UnmarshalSignCsrActionResult unmarshals an instance of SignCsrActionResult from the specified map of raw messages.
func UnmarshalSignCsrActionResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SignCsrActionResult)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDNSDomains)
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
	err = core.UnmarshalModel(m, "data", &obj.Data, UnmarshalSignActionResultData)
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

// SignIntermediateAction : A request to sign an intermediate certificate authority.
// This model "extends" ConfigAction
type SignIntermediateAction struct {
	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `12h`. The value can't exceed
	// the `max_ttl` that is defined in the associated certificate template.
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	PermittedDNSDomains []string `json:"permitted_dns_domains,omitempty"`

	// Determines whether to use values from a certificate signing request (CSR) to complete a `sign_csr` action. If set to
	// `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than using the values
	// provided in the other parameters to this operation.
	//
	// 2) Any key usages (for example, non-repudiation) that are requested in the CSR are added to the basic set of key
	// usages used for CA certs signed by this intermediate authority.
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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// The intermediate certificate authority to be signed. The name must match one of the pre-configured intermediate
	// certificate authorities.
	IntermediateCertificateAuthority *string `json:"intermediate_certificate_authority" validate:"required"`
}

// Constants associated with the SignIntermediateAction.Format property.
// The format of the returned data.
const (
	SignIntermediateActionFormatPemConst       = "pem"
	SignIntermediateActionFormatPemBundleConst = "pem_bundle"
)

// NewSignIntermediateAction : Instantiate SignIntermediateAction (Generic Model Constructor)
func (*SecretsManagerV1) NewSignIntermediateAction(intermediateCertificateAuthority string) (_model *SignIntermediateAction, err error) {
	_model = &SignIntermediateAction{
		IntermediateCertificateAuthority: core.StringPtr(intermediateCertificateAuthority),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*SignIntermediateAction) isaConfigAction() bool {
	return true
}

// UnmarshalSignIntermediateAction unmarshals an instance of SignIntermediateAction from the specified map of raw messages.
func UnmarshalSignIntermediateAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SignIntermediateAction)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDNSDomains)
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
	err = core.UnmarshalPrimitive(m, "intermediate_certificate_authority", &obj.IntermediateCertificateAuthority)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SignIntermediateActionResult : Properties that are returned with a successful `sign_intermediate` action.
// This model "extends" ConfigElementActionResultConfig
type SignIntermediateActionResult struct {
	// The fully qualified domain name or host domain name for the certificate.
	CommonName *string `json:"common_name,omitempty"`

	// The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	//
	// The alternative names can be host names or email addresses.
	AltNames *string `json:"alt_names,omitempty"`

	// The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	IPSans *string `json:"ip_sans,omitempty"`

	// The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	URISans *string `json:"uri_sans,omitempty"`

	// The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.
	//
	// The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated
	// certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is
	// `UTF8`.
	OtherSans []string `json:"other_sans,omitempty"`

	// The time-to-live (TTL) to assign to a private certificate.
	//
	// The value can be supplied as a string representation of a duration in hours, such as `12h`. The value can't exceed
	// the `max_ttl` that is defined in the associated certificate template.
	TTL interface{} `json:"ttl,omitempty"`

	// The format of the returned data.
	Format *string `json:"format,omitempty"`

	// The maximum path length to encode in the generated certificate. `-1` means no limit.
	//
	// If the signing certificate has a maximum path length set, the path length is set to one less than that of the
	// signing certificate. A limit of `0` means a literal path length of zero.
	MaxPathLength *int64 `json:"max_path_length,omitempty"`

	// Controls whether the common name is excluded from Subject Alternative Names (SANs).
	//
	// If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if
	// the common name is not a hostname or an email address, but is instead a human-readable identifier.
	ExcludeCnFromSans *bool `json:"exclude_cn_from_sans,omitempty"`

	// The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	PermittedDNSDomains []string `json:"permitted_dns_domains,omitempty"`

	// Determines whether to use values from a certificate signing request (CSR) to complete a `sign_csr` action. If set to
	// `true`, then:
	//
	// 1) Subject information, including names and alternate names, are preserved from the CSR rather than using the values
	// provided in the other parameters to this operation.
	//
	// 2) Any key usages (for example, non-repudiation) that are requested in the CSR are added to the basic set of key
	// usages used for CA certs signed by this intermediate authority.
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

	// The Street Address values in the subject field of the resulting certificate.
	StreetAddress []string `json:"street_address,omitempty"`

	// The Postal Code values in the subject field of the resulting certificate.
	PostalCode []string `json:"postal_code,omitempty"`

	// The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	SerialNumber *string `json:"serial_number,omitempty"`

	// Properties that are returned with a successful `sign` action.
	Data *SignIntermediateActionResultData `json:"data,omitempty"`

	// The signed intermediate certificate authority.
	IntermediateCertificateAuthority *string `json:"intermediate_certificate_authority" validate:"required"`
}

// Constants associated with the SignIntermediateActionResult.Format property.
// The format of the returned data.
const (
	SignIntermediateActionResultFormatPemConst       = "pem"
	SignIntermediateActionResultFormatPemBundleConst = "pem_bundle"
)

func (*SignIntermediateActionResult) isaConfigElementActionResultConfig() bool {
	return true
}

// UnmarshalSignIntermediateActionResult unmarshals an instance of SignIntermediateActionResult from the specified map of raw messages.
func UnmarshalSignIntermediateActionResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SignIntermediateActionResult)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "alt_names", &obj.AltNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_sans", &obj.IPSans)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uri_sans", &obj.URISans)
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
	err = core.UnmarshalPrimitive(m, "permitted_dns_domains", &obj.PermittedDNSDomains)
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
	err = core.UnmarshalModel(m, "data", &obj.Data, UnmarshalSignIntermediateActionResultData)
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

// UsernamePasswordSecretMetadata : Metadata properties that describe a username_password secret.
// This model "extends" SecretMetadata
type UsernamePasswordSecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be in the range 2 - 30 characters, including spaces. Special characters
	// that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`
}

// Constants associated with the UsernamePasswordSecretMetadata.SecretType property.
// The secret type.
const (
	UsernamePasswordSecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	UsernamePasswordSecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	UsernamePasswordSecretMetadataSecretTypeImportedCertConst     = "imported_cert"
	UsernamePasswordSecretMetadataSecretTypeKvConst               = "kv"
	UsernamePasswordSecretMetadataSecretTypePrivateCertConst      = "private_cert"
	UsernamePasswordSecretMetadataSecretTypePublicCertConst       = "public_cert"
	UsernamePasswordSecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

// NewUsernamePasswordSecretMetadata : Instantiate UsernamePasswordSecretMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewUsernamePasswordSecretMetadata(name string) (_model *UsernamePasswordSecretMetadata, err error) {
	_model = &UsernamePasswordSecretMetadata{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*UsernamePasswordSecretMetadata) isaSecretMetadata() bool {
	return true
}

// UnmarshalUsernamePasswordSecretMetadata unmarshals an instance of UsernamePasswordSecretMetadata from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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

// UsernamePasswordSecretResource : Properties that describe a secret.
// This model "extends" SecretResource
type UsernamePasswordSecretResource struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// A human-readable alias to assign to your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	Name *string `json:"name" validate:"required"`

	// An extended description of your secret.
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	Description *string `json:"description,omitempty"`

	// The v4 UUID that uniquely identifies the secret group to assign to this secret.
	//
	// If you omit this parameter, your secret is assigned to the `default` secret group.
	SecretGroupID *string `json:"secret_group_id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
	//
	// To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	Labels []string `json:"labels,omitempty"`

	// The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,
	// Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	State *int64 `json:"state,omitempty"`

	// A text representation of the secret state.
	StateDescription *string `json:"state_description,omitempty"`

	// The secret type.
	SecretType *string `json:"secret_type,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the actual secret is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The number of versions that are associated with a secret.
	VersionsTotal *int64 `json:"versions_total,omitempty"`

	// An array that contains metadata for each secret version. For more information on the metadata properties, see [Get
	// secret version metadata](#get-secret-version-metadata).
	Versions []map[string]interface{} `json:"versions,omitempty"`

	// The number of locks that are associated with a secret.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret metadata that a user can customize.
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// The username to assign to this secret.
	Username *string `json:"username,omitempty"`

	// The password to assign to this secret.
	Password *string `json:"password,omitempty"`

	// The data that is associated with the secret version. The data object contains the following fields:
	//
	// - `username`: The username that is associated with the secret version.
	// - `password`: The password that is associated with the secret version.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The date that the secret is scheduled for automatic rotation.
	//
	// The service automatically creates a new version of the secret on its next rotation date. This field exists only for
	// secrets that can be auto-rotated and have an existing rotation policy.
	NextRotationDate *strfmt.DateTime `json:"next_rotation_date,omitempty"`
}

// Constants associated with the UsernamePasswordSecretResource.SecretType property.
// The secret type.
const (
	UsernamePasswordSecretResourceSecretTypeArbitraryConst        = "arbitrary"
	UsernamePasswordSecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	UsernamePasswordSecretResourceSecretTypeImportedCertConst     = "imported_cert"
	UsernamePasswordSecretResourceSecretTypeKvConst               = "kv"
	UsernamePasswordSecretResourceSecretTypePrivateCertConst      = "private_cert"
	UsernamePasswordSecretResourceSecretTypePublicCertConst       = "public_cert"
	UsernamePasswordSecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// NewUsernamePasswordSecretResource : Instantiate UsernamePasswordSecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewUsernamePasswordSecretResource(name string) (_model *UsernamePasswordSecretResource, err error) {
	_model = &UsernamePasswordSecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*UsernamePasswordSecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalUsernamePasswordSecretResource unmarshals an instance of UsernamePasswordSecretResource from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretResource)
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
	err = core.UnmarshalPrimitive(m, "secret_group_id", &obj.SecretGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
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
	err = core.UnmarshalPrimitive(m, "secret_type", &obj.SecretType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_update_date", &obj.LastUpdateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions_total", &obj.VersionsTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
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
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
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

// UsernamePasswordSecretVersion : UsernamePasswordSecretVersion struct
// This model "extends" SecretVersion
type UsernamePasswordSecretVersion struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`

	// The data that is associated with the secret version. The data object contains the following fields:
	//
	// - `username`: The username that is associated with the secret version.
	// - `password`: The password that is associated with the secret version.
	SecretData map[string]interface{} `json:"secret_data,omitempty"`
}

func (*UsernamePasswordSecretVersion) isaSecretVersion() bool {
	return true
}

// UnmarshalUsernamePasswordSecretVersion unmarshals an instance of UsernamePasswordSecretVersion from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_data", &obj.SecretData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsernamePasswordSecretVersionInfo : UsernamePasswordSecretVersionInfo struct
// This model "extends" SecretVersionInfo
type UsernamePasswordSecretVersionInfo struct {
	// The ID of the secret version.
	ID *string `json:"id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`
}

func (*UsernamePasswordSecretVersionInfo) isaSecretVersionInfo() bool {
	return true
}

// UnmarshalUsernamePasswordSecretVersionInfo unmarshals an instance of UsernamePasswordSecretVersionInfo from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretVersionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretVersionInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsernamePasswordSecretVersionMetadata : Properties that describe a secret version.
// This model "extends" SecretVersionMetadata
type UsernamePasswordSecretVersionMetadata struct {
	// The v4 UUID that uniquely identifies the secret.
	ID *string `json:"id,omitempty"`

	// The ID of the secret version.
	VersionID *string `json:"version_id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret version.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the payload for the secret version is stored and available.
	PayloadAvailable *bool `json:"payload_available,omitempty"`

	// Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service
	// API.
	Downloaded *bool `json:"downloaded,omitempty"`

	// The number of locks that are associated with a secret version.
	LocksTotal *int64 `json:"locks_total,omitempty"`

	// The secret version metadata that a user can customize.
	VersionCustomMetadata map[string]interface{} `json:"version_custom_metadata,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`
}

func (*UsernamePasswordSecretVersionMetadata) isaSecretVersionMetadata() bool {
	return true
}

// UnmarshalUsernamePasswordSecretVersionMetadata unmarshals an instance of UsernamePasswordSecretVersionMetadata from the specified map of raw messages.
func UnmarshalUsernamePasswordSecretVersionMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsernamePasswordSecretVersionMetadata)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payload_available", &obj.PayloadAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "downloaded", &obj.Downloaded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locks_total", &obj.LocksTotal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_custom_metadata", &obj.VersionCustomMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
