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
 * IBM OpenAPI SDK Code Generator Version: 3.29.0-cd9ba74f-20210305-183535
 */

// Package secretsmanagerv1 : Operations and models for the SecretsManagerV1 service
package secretsmanagerv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/secrets-manager-go-sdk/common"
	"github.com/go-openapi/strfmt"
	"net/http"
	"reflect"
	"time"
)

// SecretsManagerV1 : With IBM Cloud® Secrets Manager, you can create, lease, and centrally manage secrets that are used
// in IBM Cloud services or your custom-built applications. Secrets are stored in a dedicated instance of Secrets
// Manager, built on open source HashiCorp Vault.
//
// Version: 1.0.0
// See: https://cloud.ibm.com/docs/secrets-manager
type SecretsManagerV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://secrets-manager.cloud.ibm.com"

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

// PutConfig : Configure secrets of a given type
// Updates the configuration for the given secret type.
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

	_, err = builder.SetBodyContentJSON(putConfigOptions.EngineConfigOneOf)
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

// GetConfig : Get the configuration for a secret type
// Retrieves the configuration that is associated with the given secret type.
func (secretsManager *SecretsManagerV1) GetConfig(getConfigOptions *GetConfigOptions) (result EngineConfigOneOfIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetConfigWithContext(context.Background(), getConfigOptions)
}

// GetConfigWithContext is an alternate form of the GetConfig method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetConfigWithContext(ctx context.Context, getConfigOptions *GetConfigOptions) (result EngineConfigOneOfIntf, response *core.DetailedResponse, err error) {
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEngineConfigOneOf)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// PutPolicy : Set secret policies
// Creates or updates one or more policies, such as an [automatic rotation
// policy](http://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-rotate-secrets#auto-rotate-secret), for the
// specified secret.
func (secretsManager *SecretsManagerV1) PutPolicy(putPolicyOptions *PutPolicyOptions) (result GetSecretPoliciesOneOfIntf, response *core.DetailedResponse, err error) {
	return secretsManager.PutPolicyWithContext(context.Background(), putPolicyOptions)
}

// PutPolicyWithContext is an alternate form of the PutPolicy method which supports a Context parameter
func (secretsManager *SecretsManagerV1) PutPolicyWithContext(ctx context.Context, putPolicyOptions *PutPolicyOptions) (result GetSecretPoliciesOneOfIntf, response *core.DetailedResponse, err error) {
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretPoliciesOneOf)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetPolicy : List secret policies
// Retrieves a list of policies that are associated with a specified secret.
func (secretsManager *SecretsManagerV1) GetPolicy(getPolicyOptions *GetPolicyOptions) (result GetSecretPoliciesOneOfIntf, response *core.DetailedResponse, err error) {
	return secretsManager.GetPolicyWithContext(context.Background(), getPolicyOptions)
}

// GetPolicyWithContext is an alternate form of the GetPolicy method which supports a Context parameter
func (secretsManager *SecretsManagerV1) GetPolicyWithContext(ctx context.Context, getPolicyOptions *GetPolicyOptions) (result GetSecretPoliciesOneOfIntf, response *core.DetailedResponse, err error) {
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecretPoliciesOneOf)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateSecretGroup : Create a secret group
// Creates a secret group that you can use to organize secrets and control who on your team has access to them.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupDef)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListSecretGroups : List secret groups
// Retrieves the list of secret groups that are available in your Secrets Manager instance.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupDef)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetSecretGroup : Get a secret group
// Retrieves the metadata of an existing secret group by specifying the ID of the group.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupDef)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateSecretGroupMetadata : Update a secret group
// Updates the metadata of an existing secret group, such as its name or description.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretGroupDef)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteSecretGroup : Delete a secret group
// Deletes a secret group by specifying the ID of the secret group.
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
// Creates a secret that you can use to access or authenticate to a protected resource.
//
// A successful request stores the secret in your dedicated instance based on the secret type and data that you specify.
// The response returns the ID value of the secret, along with other metadata.
//
// To learn more about the types of secrets that you can create with Secrets Manager, check out the
// [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-secret-basics).
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateSecret)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListSecrets : List secrets by type
// Retrieves a list of secrets based on the type that you specify.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListSecrets)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListAllSecrets : List all secrets
// Retrieves a list of all secrets in your Secrets Manager instance.
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = secretsManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListSecrets)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetSecret : Get a secret
// Retrieves a secret and its details by specifying the ID of the secret.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecret)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateSecret : Invoke an action on a secret
// Invokes an action on a specified secret. This method supports the following actions:
//
// - `rotate`: Replace the value of an `arbitrary` or `username_password` secret.
// - `delete_credentials`: Delete the API key that is associated with an `iam_credentials` secret.
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

	_, err = builder.SetBodyContentJSON(updateSecretOptions.SecretActionOneOf)
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSecret)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteSecret : Delete a secret
// Deletes a secret by specifying the ID of the secret.
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

// GetSecretMetadata : Get secret metadata
// Retrieves the details of a secret by specifying the ID.
//
// A successful request returns only metadata about the secret, such as its name and creation date. To retrieve the
// value of a secret, use the [Get a secret](#get-secret) method.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretMetadataRequest)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateSecretMetadata : Update secret metadata
// Updates the metadata of a secret, such as its name or description.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretMetadataRequest)
	if err != nil {
		return
	}
	response.Result = result

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
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerErrorJSONConst         = "application/vnd.ibm.secrets-manager.error+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretGroupJSONConst   = "application/vnd.ibm.secrets-manager.secret.group+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretJSONConst        = "application/vnd.ibm.secrets-manager.secret+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretPolicyJSONConst  = "application/vnd.ibm.secrets-manager.secret.policy+json"
	CollectionMetadataCollectionTypeApplicationVndIBMSecretsManagerSecretVersionJSONConst = "application/vnd.ibm.secrets-manager.secret.version+json"
)

// NewCollectionMetadata : Instantiate CollectionMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewCollectionMetadata(collectionType string, collectionTotal int64) (model *CollectionMetadata, err error) {
	model = &CollectionMetadata{
		CollectionType:  core.StringPtr(collectionType),
		CollectionTotal: core.Int64Ptr(collectionTotal),
	}
	err = core.ValidateStruct(model, "required parameters")
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

// CreateSecret : The base schema for creating secrets.
type CreateSecret struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretResourceIntf `json:"resources" validate:"required"`
}

// NewCreateSecret : Instantiate CreateSecret (Generic Model Constructor)
func (*SecretsManagerV1) NewCreateSecret(metadata *CollectionMetadata, resources []SecretResourceIntf) (model *CreateSecret, err error) {
	model = &CreateSecret{
		Metadata:  metadata,
		Resources: resources,
	}
	err = core.ValidateStruct(model, "required parameters")
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
	Metadata *CollectionMetadata `validate:"required"`

	// A collection of resources.
	Resources []SecretGroupResource `validate:"required"`

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
func (options *CreateSecretGroupOptions) SetMetadata(metadata *CollectionMetadata) *CreateSecretGroupOptions {
	options.Metadata = metadata
	return options
}

// SetResources : Allow user to set Resources
func (options *CreateSecretGroupOptions) SetResources(resources []SecretGroupResource) *CreateSecretGroupOptions {
	options.Resources = resources
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretGroupOptions) SetHeaders(param map[string]string) *CreateSecretGroupOptions {
	options.Headers = param
	return options
}

// CreateSecretOptions : The CreateSecret options.
type CreateSecretOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `validate:"required"`

	// A collection of resources.
	Resources []SecretResourceIntf `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateSecretOptions.SecretType property.
// The secret type.
const (
	CreateSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	CreateSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
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
func (options *CreateSecretOptions) SetSecretType(secretType string) *CreateSecretOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *CreateSecretOptions) SetMetadata(metadata *CollectionMetadata) *CreateSecretOptions {
	options.Metadata = metadata
	return options
}

// SetResources : Allow user to set Resources
func (options *CreateSecretOptions) SetResources(resources []SecretResourceIntf) *CreateSecretOptions {
	options.Resources = resources
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretOptions) SetHeaders(param map[string]string) *CreateSecretOptions {
	options.Headers = param
	return options
}

// DeleteSecretGroupOptions : The DeleteSecretGroup options.
type DeleteSecretGroupOptions struct {
	// The v4 UUID that uniquely identifies the secret group.
	ID *string `validate:"required,ne="`

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
func (options *DeleteSecretGroupOptions) SetID(id string) *DeleteSecretGroupOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecretGroupOptions) SetHeaders(param map[string]string) *DeleteSecretGroupOptions {
	options.Headers = param
	return options
}

// DeleteSecretOptions : The DeleteSecret options.
type DeleteSecretOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteSecretOptions.SecretType property.
// The secret type.
const (
	DeleteSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	DeleteSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
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
func (options *DeleteSecretOptions) SetSecretType(secretType string) *DeleteSecretOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetID : Allow user to set ID
func (options *DeleteSecretOptions) SetID(id string) *DeleteSecretOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecretOptions) SetHeaders(param map[string]string) *DeleteSecretOptions {
	options.Headers = param
	return options
}

// EngineConfigOneOf : EngineConfigOneOf struct
// Models which "extend" this model:
// - EngineConfigOneOfIamSecretEngineRootConfig
type EngineConfigOneOf struct {
	// An IBM Cloud API key that has the capability to create and manage service IDs.
	//
	// The API key must be assigned the Editor platform role on the Access Groups Service and the Operator platform role on
	// the IAM Identity Service. For more information, see [Enabling the IAM secrets
	// engine](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-secret-engines#configure-iam-engine).
	APIKey *string `json:"api_key,omitempty"`

	// The hash value of the IBM Cloud API key that is used to create and manage service IDs.
	APIKeyHash *string `json:"api_key_hash,omitempty"`
}

func (*EngineConfigOneOf) isaEngineConfigOneOf() bool {
	return true
}

type EngineConfigOneOfIntf interface {
	isaEngineConfigOneOf() bool
}

// UnmarshalEngineConfigOneOf unmarshals an instance of EngineConfigOneOf from the specified map of raw messages.
func UnmarshalEngineConfigOneOf(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EngineConfigOneOf)
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

// GetConfigOptions : The GetConfig options.
type GetConfigOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetConfigOptions.SecretType property.
// The secret type.
const (
	GetConfigOptionsSecretTypeIamCredentialsConst = "iam_credentials"
)

// NewGetConfigOptions : Instantiate GetConfigOptions
func (*SecretsManagerV1) NewGetConfigOptions(secretType string) *GetConfigOptions {
	return &GetConfigOptions{
		SecretType: core.StringPtr(secretType),
	}
}

// SetSecretType : Allow user to set SecretType
func (options *GetConfigOptions) SetSecretType(secretType string) *GetConfigOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigOptions) SetHeaders(param map[string]string) *GetConfigOptions {
	options.Headers = param
	return options
}

// GetPolicyOptions : The GetPolicy options.
type GetPolicyOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `validate:"required,ne="`

	// The type of policy that is associated with the specified secret.
	Policy *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetPolicyOptions.SecretType property.
// The secret type.
const (
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
func (options *GetPolicyOptions) SetSecretType(secretType string) *GetPolicyOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetID : Allow user to set ID
func (options *GetPolicyOptions) SetID(id string) *GetPolicyOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetPolicy : Allow user to set Policy
func (options *GetPolicyOptions) SetPolicy(policy string) *GetPolicyOptions {
	options.Policy = core.StringPtr(policy)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetPolicyOptions) SetHeaders(param map[string]string) *GetPolicyOptions {
	options.Headers = param
	return options
}

// GetSecret : The base schema for retrieving a secret.
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
	ID *string `validate:"required,ne="`

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
func (options *GetSecretGroupOptions) SetID(id string) *GetSecretGroupOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretGroupOptions) SetHeaders(param map[string]string) *GetSecretGroupOptions {
	options.Headers = param
	return options
}

// GetSecretMetadataOptions : The GetSecretMetadata options.
type GetSecretMetadataOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSecretMetadataOptions.SecretType property.
// The secret type.
const (
	GetSecretMetadataOptionsSecretTypeArbitraryConst        = "arbitrary"
	GetSecretMetadataOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
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
func (options *GetSecretMetadataOptions) SetSecretType(secretType string) *GetSecretMetadataOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetID : Allow user to set ID
func (options *GetSecretMetadataOptions) SetID(id string) *GetSecretMetadataOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretMetadataOptions) SetHeaders(param map[string]string) *GetSecretMetadataOptions {
	options.Headers = param
	return options
}

// GetSecretOptions : The GetSecret options.
type GetSecretOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSecretOptions.SecretType property.
// The secret type.
const (
	GetSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	GetSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
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
func (options *GetSecretOptions) SetSecretType(secretType string) *GetSecretOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetID : Allow user to set ID
func (options *GetSecretOptions) SetID(id string) *GetSecretOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretOptions) SetHeaders(param map[string]string) *GetSecretOptions {
	options.Headers = param
	return options
}

// GetSecretPoliciesOneOf : GetSecretPoliciesOneOf struct
// Models which "extend" this model:
// - GetSecretPoliciesOneOfGetSecretPolicyRotation
type GetSecretPoliciesOneOf struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata,omitempty"`

	// A collection of resources.
	Resources []GetSecretPoliciesOneOfResourcesItem `json:"resources,omitempty"`
}

func (*GetSecretPoliciesOneOf) isaGetSecretPoliciesOneOf() bool {
	return true
}

type GetSecretPoliciesOneOfIntf interface {
	isaGetSecretPoliciesOneOf() bool
}

// UnmarshalGetSecretPoliciesOneOf unmarshals an instance of GetSecretPoliciesOneOf from the specified map of raw messages.
func UnmarshalGetSecretPoliciesOneOf(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretPoliciesOneOf)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalGetSecretPoliciesOneOfResourcesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem : Properties that are associated with a rotation policy.
type GetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem struct {
	// The v4 UUID that uniquely identifies the policy.
	ID *string `json:"id" validate:"required"`

	// The Cloud Resource Name (CRN) that uniquely identifies your cloud resources.
	CRN *string `json:"crn,omitempty"`

	// The date the policy was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the policy.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the policy is replaced or modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The unique identifier for the entity that updated the policy.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The MIME type that represents the policy. Currently, only the default is supported.
	Type *string `json:"type" validate:"required"`

	// The secret rotation time interval.
	Rotation *SecretPolicyRotationRotation `json:"rotation" validate:"required"`
}

// Constants associated with the GetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem.Type property.
// The MIME type that represents the policy. Currently, only the default is supported.
const (
	GetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItemTypeApplicationVndIBMSecretsManagerSecretPolicyJSONConst = "application/vnd.ibm.secrets-manager.secret.policy+json"
)

// UnmarshalGetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem unmarshals an instance of GetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem from the specified map of raw messages.
func UnmarshalGetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
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

// GetSecretPoliciesOneOfResourcesItem : Properties that are associated with a rotation policy.
type GetSecretPoliciesOneOfResourcesItem struct {
	// The v4 UUID that uniquely identifies the policy.
	ID *string `json:"id" validate:"required"`

	// The Cloud Resource Name (CRN) that uniquely identifies your cloud resources.
	CRN *string `json:"crn,omitempty"`

	// The date the policy was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the policy.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when the policy is replaced or modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`

	// The unique identifier for the entity that updated the policy.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The MIME type that represents the policy. Currently, only the default is supported.
	Type *string `json:"type" validate:"required"`

	// The secret rotation time interval.
	Rotation *SecretPolicyRotationRotation `json:"rotation" validate:"required"`
}

// Constants associated with the GetSecretPoliciesOneOfResourcesItem.Type property.
// The MIME type that represents the policy. Currently, only the default is supported.
const (
	GetSecretPoliciesOneOfResourcesItemTypeApplicationVndIBMSecretsManagerSecretPolicyJSONConst = "application/vnd.ibm.secrets-manager.secret.policy+json"
)

// UnmarshalGetSecretPoliciesOneOfResourcesItem unmarshals an instance of GetSecretPoliciesOneOfResourcesItem from the specified map of raw messages.
func UnmarshalGetSecretPoliciesOneOfResourcesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretPoliciesOneOfResourcesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
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

// ListAllSecretsOptions : The ListAllSecrets options.
type ListAllSecretsOptions struct {
	// The number of secrets to retrieve. By default, list operations return the first 200 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 secrets in your instance, and you want to retrieve only the first 5 secrets, use
	// `../secrets/{secret-type}?limit=5`.
	Limit *int64

	// The number of secrets to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 secrets in your instance, and you want to retrieve secrets 26 through 50, use
	// `../secrets/{secret-type}?offset=25&limit=25`.
	Offset *int64

	// Filter secrets that contain the specified string. The fields that are searched include: id, name, description,
	// labels, secret_type.
	//
	// **Usage:** If you want to list only the secrets that contain the string "text", use
	// `../secrets/{secret-type}?search=text`.
	Search *string

	// Sort a list of secrets by the specified field.
	//
	// **Usage:** To sort a list of secrets by their creation date, use
	// `../secrets/{secret-type}?sort_by=creation_date`.
	SortBy *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListAllSecretsOptions.SortBy property.
// Sort a list of secrets by the specified field.
//
// **Usage:** To sort a list of secrets by their creation date, use
// `../secrets/{secret-type}?sort_by=creation_date`.
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
func (options *ListAllSecretsOptions) SetLimit(limit int64) *ListAllSecretsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListAllSecretsOptions) SetOffset(offset int64) *ListAllSecretsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetSearch : Allow user to set Search
func (options *ListAllSecretsOptions) SetSearch(search string) *ListAllSecretsOptions {
	options.Search = core.StringPtr(search)
	return options
}

// SetSortBy : Allow user to set SortBy
func (options *ListAllSecretsOptions) SetSortBy(sortBy string) *ListAllSecretsOptions {
	options.SortBy = core.StringPtr(sortBy)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllSecretsOptions) SetHeaders(param map[string]string) *ListAllSecretsOptions {
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

// ListSecrets : The base schema for listing secrets.
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
	SecretType *string `validate:"required,ne="`

	// The number of secrets to retrieve. By default, list operations return the first 200 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 secrets in your instance, and you want to retrieve only the first 5 secrets, use
	// `../secrets/{secret-type}?limit=5`.
	Limit *int64

	// The number of secrets to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 secrets in your instance, and you want to retrieve secrets 26 through 50, use
	// `../secrets/{secret-type}?offset=25&limit=25`.
	Offset *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListSecretsOptions.SecretType property.
// The secret type.
const (
	ListSecretsOptionsSecretTypeArbitraryConst        = "arbitrary"
	ListSecretsOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	ListSecretsOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewListSecretsOptions : Instantiate ListSecretsOptions
func (*SecretsManagerV1) NewListSecretsOptions(secretType string) *ListSecretsOptions {
	return &ListSecretsOptions{
		SecretType: core.StringPtr(secretType),
	}
}

// SetSecretType : Allow user to set SecretType
func (options *ListSecretsOptions) SetSecretType(secretType string) *ListSecretsOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListSecretsOptions) SetLimit(limit int64) *ListSecretsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListSecretsOptions) SetOffset(offset int64) *ListSecretsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretsOptions) SetHeaders(param map[string]string) *ListSecretsOptions {
	options.Headers = param
	return options
}

// PutConfigOptions : The PutConfig options.
type PutConfigOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The base request for setting secret engine configuration.
	EngineConfigOneOf EngineConfigOneOfIntf `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutConfigOptions.SecretType property.
// The secret type.
const (
	PutConfigOptionsSecretTypeIamCredentialsConst = "iam_credentials"
)

// NewPutConfigOptions : Instantiate PutConfigOptions
func (*SecretsManagerV1) NewPutConfigOptions(secretType string, engineConfigOneOf EngineConfigOneOfIntf) *PutConfigOptions {
	return &PutConfigOptions{
		SecretType:        core.StringPtr(secretType),
		EngineConfigOneOf: engineConfigOneOf,
	}
}

// SetSecretType : Allow user to set SecretType
func (options *PutConfigOptions) SetSecretType(secretType string) *PutConfigOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetEngineConfigOneOf : Allow user to set EngineConfigOneOf
func (options *PutConfigOptions) SetEngineConfigOneOf(engineConfigOneOf EngineConfigOneOfIntf) *PutConfigOptions {
	options.EngineConfigOneOf = engineConfigOneOf
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PutConfigOptions) SetHeaders(param map[string]string) *PutConfigOptions {
	options.Headers = param
	return options
}

// PutPolicyOptions : The PutPolicy options.
type PutPolicyOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `validate:"required"`

	// A collection of resources.
	Resources []SecretPolicyRotation `validate:"required"`

	// The type of policy that is associated with the specified secret.
	Policy *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutPolicyOptions.SecretType property.
// The secret type.
const (
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
func (options *PutPolicyOptions) SetSecretType(secretType string) *PutPolicyOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetID : Allow user to set ID
func (options *PutPolicyOptions) SetID(id string) *PutPolicyOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *PutPolicyOptions) SetMetadata(metadata *CollectionMetadata) *PutPolicyOptions {
	options.Metadata = metadata
	return options
}

// SetResources : Allow user to set Resources
func (options *PutPolicyOptions) SetResources(resources []SecretPolicyRotation) *PutPolicyOptions {
	options.Resources = resources
	return options
}

// SetPolicy : Allow user to set Policy
func (options *PutPolicyOptions) SetPolicy(policy string) *PutPolicyOptions {
	options.Policy = core.StringPtr(policy)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PutPolicyOptions) SetHeaders(param map[string]string) *PutPolicyOptions {
	options.Headers = param
	return options
}

// SecretActionOneOf : SecretActionOneOf struct
// Models which "extend" this model:
// - SecretActionOneOfRotateArbitrarySecretBody
// - SecretActionOneOfRotateUsernamePasswordSecretBody
// - SecretActionOneOfDeleteCredentialsForIamSecret
type SecretActionOneOf struct {
	// The new secret data to assign to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`

	// The new password to assign to a `username_password` secret.
	Password *string `json:"password,omitempty"`

	// The service ID that you want to delete. It is deleted together with its API key.
	ServiceID *string `json:"service_id,omitempty"`
}

func (*SecretActionOneOf) isaSecretActionOneOf() bool {
	return true
}

type SecretActionOneOfIntf interface {
	isaSecretActionOneOf() bool
}

// UnmarshalSecretActionOneOf unmarshals an instance of SecretActionOneOf from the specified map of raw messages.
func UnmarshalSecretActionOneOf(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretActionOneOf)
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
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

// SecretGroupDef : The base schema definition for a secret group.
type SecretGroupDef struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretGroupResource `json:"resources" validate:"required"`
}

// NewSecretGroupDef : Instantiate SecretGroupDef (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretGroupDef(metadata *CollectionMetadata, resources []SecretGroupResource) (model *SecretGroupDef, err error) {
	model = &SecretGroupDef{
		Metadata:  metadata,
		Resources: resources,
	}
	err = core.ValidateStruct(model, "required parameters")
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

// SecretGroupMetadataUpdatable : Metadata properties that describe a secret group.
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

// SecretMetadata : Metadata properties that describe a secret.
type SecretMetadata struct {
	// The unique ID of the secret.
	ID *string `json:"id,omitempty"`

	// Labels that you can use to filter for secrets in your instance.
	//
	// Up to 30 labels can be created. Labels can be between 2-30 characters, including spaces. Special characters not
	// permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).
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

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to generated credentials.
	//
	// For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be
	// either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m`
	// or `24h`.
	TTL interface{} `json:"ttl,omitempty"`

	// The Cloud Resource Name (CRN) that uniquely identifies the resource.
	CRN *string `json:"crn,omitempty"`

	// The date the secret was created. The date format follows RFC 3339.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Updates when any part of the secret metadata is modified. The date format follows RFC 3339.
	LastUpdateDate *strfmt.DateTime `json:"last_update_date,omitempty"`
}

// Constants associated with the SecretMetadata.SecretType property.
// The secret type.
const (
	SecretMetadataSecretTypeArbitraryConst        = "arbitrary"
	SecretMetadataSecretTypeIamCredentialsConst   = "iam_credentials"
	SecretMetadataSecretTypeUsernamePasswordConst = "username_password"
)

// NewSecretMetadata : Instantiate SecretMetadata (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretMetadata(name string) (model *SecretMetadata, err error) {
	model = &SecretMetadata{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
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
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretMetadataRequest : The metadata of a secret.
type SecretMetadataRequest struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []SecretMetadata `json:"resources" validate:"required"`
}

// NewSecretMetadataRequest : Instantiate SecretMetadataRequest (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretMetadataRequest(metadata *CollectionMetadata, resources []SecretMetadata) (model *SecretMetadataRequest, err error) {
	model = &SecretMetadataRequest{
		Metadata:  metadata,
		Resources: resources,
	}
	err = core.ValidateStruct(model, "required parameters")
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

// SecretPolicyRotation : Properties that are associated with a rotation policy.
type SecretPolicyRotation struct {
	// The MIME type that represents the policy. Currently, only the default is supported.
	Type *string `json:"type" validate:"required"`

	// The secret rotation time interval.
	Rotation *SecretPolicyRotationRotation `json:"rotation" validate:"required"`
}

// Constants associated with the SecretPolicyRotation.Type property.
// The MIME type that represents the policy. Currently, only the default is supported.
const (
	SecretPolicyRotationTypeApplicationVndIBMSecretsManagerSecretPolicyJSONConst = "application/vnd.ibm.secrets-manager.secret.policy+json"
)

// NewSecretPolicyRotation : Instantiate SecretPolicyRotation (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretPolicyRotation(typeVar string, rotation *SecretPolicyRotationRotation) (model *SecretPolicyRotation, err error) {
	model = &SecretPolicyRotation{
		Type:     core.StringPtr(typeVar),
		Rotation: rotation,
	}
	err = core.ValidateStruct(model, "required parameters")
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

// SecretPolicyRotationRotation : The secret rotation time interval.
type SecretPolicyRotationRotation struct {
	// Specifies the length of the secret rotation time interval.
	Interval *int64 `json:"interval" validate:"required"`

	// Specifies the units for the secret rotation time interval.
	Unit *string `json:"unit" validate:"required"`
}

// Constants associated with the SecretPolicyRotationRotation.Unit property.
// Specifies the units for the secret rotation time interval.
const (
	SecretPolicyRotationRotationUnitDayConst   = "day"
	SecretPolicyRotationRotationUnitMonthConst = "month"
)

// NewSecretPolicyRotationRotation : Instantiate SecretPolicyRotationRotation (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretPolicyRotationRotation(interval int64, unit string) (model *SecretPolicyRotationRotation, err error) {
	model = &SecretPolicyRotationRotation{
		Interval: core.Int64Ptr(interval),
		Unit:     core.StringPtr(unit),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretResource : SecretResource struct
// Models which "extend" this model:
// - SecretResourceArbitrarySecretResource
// - SecretResourceUsernamePasswordSecretResource
// - SecretResourceIamSecretResource
type SecretResource struct {
	// The MIME type that represents the secret.
	Type *string `json:"type,omitempty"`

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
	// Up to 30 labels can be created. Labels can be between 2-30 characters, including spaces. Special characters not
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

	// An array that contains metadata for each secret version.
	Versions []SecretVersion `json:"versions,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The new secret data to assign to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`

	SecretData interface{} `json:"secret_data,omitempty"`

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
	TTL interface{} `json:"ttl,omitempty"`

	// The access groups that define the capabilities of the service ID and API key that are generated for an
	// `iam_credentials` secret.
	//
	// **Tip:** To find the ID of an access group, go to **Manage > Access (IAM) > Access groups** in the IBM Cloud
	// console. Select the access group to inspect, and click **Details** to view its ID.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease (see the `ttl` field), the API key is deleted automatically. If you
	// want to continue to use the same API key for future read operations, see the `reuse_api_key` field.
	APIKey *string `json:"api_key,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created. This service ID is added to the access
	// groups that you assign for this secret.
	ServiceID *string `json:"service_id,omitempty"`

	// Set to `true` to reuse the service ID and API key for this secret.
	//
	// Use this field to control whether to use the same service ID and API key for future read operations on this secret.
	// If set to `true`, the service reuses the current credentials. If set to `false`, a new service ID and API key is
	// generated each time that the secret is read or accessed.
	ReuseAPIKey *bool `json:"reuse_api_key,omitempty"`
}

// Constants associated with the SecretResource.SecretType property.
// The secret type.
const (
	SecretResourceSecretTypeArbitraryConst        = "arbitrary"
	SecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	SecretResourceSecretTypeUsernamePasswordConst = "username_password"
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
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalSecretVersion)
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
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseAPIKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretVersion : Properties that are associated with a specific secret version.
type SecretVersion struct {
	// The ID of the secret version.
	ID *string `json:"id,omitempty"`

	// The date that the version of the secret was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the entity that created the secret.
	CreatedBy *string `json:"created_by,omitempty"`

	// Indicates whether the version of the secret was created by automatic rotation.
	AutoRotated *bool `json:"auto_rotated,omitempty"`
}

// UnmarshalSecretVersion unmarshals an instance of SecretVersion from the specified map of raw messages.
func UnmarshalSecretVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretVersion)
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
	err = core.UnmarshalPrimitive(m, "auto_rotated", &obj.AutoRotated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateSecretGroupMetadataOptions : The UpdateSecretGroupMetadata options.
type UpdateSecretGroupMetadataOptions struct {
	// The v4 UUID that uniquely identifies the secret group.
	ID *string `validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `validate:"required"`

	// A collection of resources.
	Resources []SecretGroupMetadataUpdatable `validate:"required"`

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
func (options *UpdateSecretGroupMetadataOptions) SetID(id string) *UpdateSecretGroupMetadataOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *UpdateSecretGroupMetadataOptions) SetMetadata(metadata *CollectionMetadata) *UpdateSecretGroupMetadataOptions {
	options.Metadata = metadata
	return options
}

// SetResources : Allow user to set Resources
func (options *UpdateSecretGroupMetadataOptions) SetResources(resources []SecretGroupMetadataUpdatable) *UpdateSecretGroupMetadataOptions {
	options.Resources = resources
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretGroupMetadataOptions) SetHeaders(param map[string]string) *UpdateSecretGroupMetadataOptions {
	options.Headers = param
	return options
}

// UpdateSecretMetadataOptions : The UpdateSecretMetadata options.
type UpdateSecretMetadataOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `validate:"required,ne="`

	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `validate:"required"`

	// A collection of resources.
	Resources []SecretMetadata `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateSecretMetadataOptions.SecretType property.
// The secret type.
const (
	UpdateSecretMetadataOptionsSecretTypeArbitraryConst        = "arbitrary"
	UpdateSecretMetadataOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	UpdateSecretMetadataOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// NewUpdateSecretMetadataOptions : Instantiate UpdateSecretMetadataOptions
func (*SecretsManagerV1) NewUpdateSecretMetadataOptions(secretType string, id string, metadata *CollectionMetadata, resources []SecretMetadata) *UpdateSecretMetadataOptions {
	return &UpdateSecretMetadataOptions{
		SecretType: core.StringPtr(secretType),
		ID:         core.StringPtr(id),
		Metadata:   metadata,
		Resources:  resources,
	}
}

// SetSecretType : Allow user to set SecretType
func (options *UpdateSecretMetadataOptions) SetSecretType(secretType string) *UpdateSecretMetadataOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetID : Allow user to set ID
func (options *UpdateSecretMetadataOptions) SetID(id string) *UpdateSecretMetadataOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *UpdateSecretMetadataOptions) SetMetadata(metadata *CollectionMetadata) *UpdateSecretMetadataOptions {
	options.Metadata = metadata
	return options
}

// SetResources : Allow user to set Resources
func (options *UpdateSecretMetadataOptions) SetResources(resources []SecretMetadata) *UpdateSecretMetadataOptions {
	options.Resources = resources
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretMetadataOptions) SetHeaders(param map[string]string) *UpdateSecretMetadataOptions {
	options.Headers = param
	return options
}

// UpdateSecretOptions : The UpdateSecret options.
type UpdateSecretOptions struct {
	// The secret type.
	SecretType *string `validate:"required,ne="`

	// The v4 UUID that uniquely identifies the secret.
	ID *string `validate:"required,ne="`

	// The action to perform on the specified secret.
	Action *string `validate:"required"`

	// The base request for invoking an action on a secret.
	SecretActionOneOf SecretActionOneOfIntf `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateSecretOptions.SecretType property.
// The secret type.
const (
	UpdateSecretOptionsSecretTypeArbitraryConst        = "arbitrary"
	UpdateSecretOptionsSecretTypeIamCredentialsConst   = "iam_credentials"
	UpdateSecretOptionsSecretTypeUsernamePasswordConst = "username_password"
)

// Constants associated with the UpdateSecretOptions.Action property.
// The action to perform on the specified secret.
const (
	UpdateSecretOptionsActionDeleteCredentialsConst = "delete_credentials"
	UpdateSecretOptionsActionRotateConst            = "rotate"
)

// NewUpdateSecretOptions : Instantiate UpdateSecretOptions
func (*SecretsManagerV1) NewUpdateSecretOptions(secretType string, id string, action string, secretActionOneOf SecretActionOneOfIntf) *UpdateSecretOptions {
	return &UpdateSecretOptions{
		SecretType:        core.StringPtr(secretType),
		ID:                core.StringPtr(id),
		Action:            core.StringPtr(action),
		SecretActionOneOf: secretActionOneOf,
	}
}

// SetSecretType : Allow user to set SecretType
func (options *UpdateSecretOptions) SetSecretType(secretType string) *UpdateSecretOptions {
	options.SecretType = core.StringPtr(secretType)
	return options
}

// SetID : Allow user to set ID
func (options *UpdateSecretOptions) SetID(id string) *UpdateSecretOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetAction : Allow user to set Action
func (options *UpdateSecretOptions) SetAction(action string) *UpdateSecretOptions {
	options.Action = core.StringPtr(action)
	return options
}

// SetSecretActionOneOf : Allow user to set SecretActionOneOf
func (options *UpdateSecretOptions) SetSecretActionOneOf(secretActionOneOf SecretActionOneOfIntf) *UpdateSecretOptions {
	options.SecretActionOneOf = secretActionOneOf
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSecretOptions) SetHeaders(param map[string]string) *UpdateSecretOptions {
	options.Headers = param
	return options
}

// EngineConfigOneOfIamSecretEngineRootConfig : Configuration that is used to generate IAM credentials.
// This model "extends" EngineConfigOneOf
type EngineConfigOneOfIamSecretEngineRootConfig struct {
	// An IBM Cloud API key that has the capability to create and manage service IDs.
	//
	// The API key must be assigned the Editor platform role on the Access Groups Service and the Operator platform role on
	// the IAM Identity Service. For more information, see [Enabling the IAM secrets
	// engine](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-secret-engines#configure-iam-engine).
	APIKey *string `json:"api_key" validate:"required"`

	// The hash value of the IBM Cloud API key that is used to create and manage service IDs.
	APIKeyHash *string `json:"api_key_hash,omitempty"`
}

// NewEngineConfigOneOfIamSecretEngineRootConfig : Instantiate EngineConfigOneOfIamSecretEngineRootConfig (Generic Model Constructor)
func (*SecretsManagerV1) NewEngineConfigOneOfIamSecretEngineRootConfig(apiKey string) (model *EngineConfigOneOfIamSecretEngineRootConfig, err error) {
	model = &EngineConfigOneOfIamSecretEngineRootConfig{
		APIKey: core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*EngineConfigOneOfIamSecretEngineRootConfig) isaEngineConfigOneOf() bool {
	return true
}

// UnmarshalEngineConfigOneOfIamSecretEngineRootConfig unmarshals an instance of EngineConfigOneOfIamSecretEngineRootConfig from the specified map of raw messages.
func UnmarshalEngineConfigOneOfIamSecretEngineRootConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EngineConfigOneOfIamSecretEngineRootConfig)
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

// GetSecretPoliciesOneOfGetSecretPolicyRotation : The base schema for retrieving a policy that is associated with a secret.
// This model "extends" GetSecretPoliciesOneOf
type GetSecretPoliciesOneOfGetSecretPolicyRotation struct {
	// The metadata that describes the resource array.
	Metadata *CollectionMetadata `json:"metadata" validate:"required"`

	// A collection of resources.
	Resources []GetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem `json:"resources" validate:"required"`
}

func (*GetSecretPoliciesOneOfGetSecretPolicyRotation) isaGetSecretPoliciesOneOf() bool {
	return true
}

// UnmarshalGetSecretPoliciesOneOfGetSecretPolicyRotation unmarshals an instance of GetSecretPoliciesOneOfGetSecretPolicyRotation from the specified map of raw messages.
func UnmarshalGetSecretPoliciesOneOfGetSecretPolicyRotation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSecretPoliciesOneOfGetSecretPolicyRotation)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCollectionMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalGetSecretPoliciesOneOfGetSecretPolicyRotationResourcesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretActionOneOfDeleteCredentialsForIamSecret : Delete the credentials that are associated with an `iam_credentials` secret.
// This model "extends" SecretActionOneOf
type SecretActionOneOfDeleteCredentialsForIamSecret struct {
	// The service ID that you want to delete. It is deleted together with its API key.
	ServiceID *string `json:"service_id" validate:"required"`
}

// NewSecretActionOneOfDeleteCredentialsForIamSecret : Instantiate SecretActionOneOfDeleteCredentialsForIamSecret (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretActionOneOfDeleteCredentialsForIamSecret(serviceID string) (model *SecretActionOneOfDeleteCredentialsForIamSecret, err error) {
	model = &SecretActionOneOfDeleteCredentialsForIamSecret{
		ServiceID: core.StringPtr(serviceID),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*SecretActionOneOfDeleteCredentialsForIamSecret) isaSecretActionOneOf() bool {
	return true
}

// UnmarshalSecretActionOneOfDeleteCredentialsForIamSecret unmarshals an instance of SecretActionOneOfDeleteCredentialsForIamSecret from the specified map of raw messages.
func UnmarshalSecretActionOneOfDeleteCredentialsForIamSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretActionOneOfDeleteCredentialsForIamSecret)
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretActionOneOfRotateArbitrarySecretBody : The request body of a `rotate` action.
// This model "extends" SecretActionOneOf
type SecretActionOneOfRotateArbitrarySecretBody struct {
	// The new secret data to assign to an `arbitrary` secret.
	Payload *string `json:"payload" validate:"required"`
}

// NewSecretActionOneOfRotateArbitrarySecretBody : Instantiate SecretActionOneOfRotateArbitrarySecretBody (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretActionOneOfRotateArbitrarySecretBody(payload string) (model *SecretActionOneOfRotateArbitrarySecretBody, err error) {
	model = &SecretActionOneOfRotateArbitrarySecretBody{
		Payload: core.StringPtr(payload),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*SecretActionOneOfRotateArbitrarySecretBody) isaSecretActionOneOf() bool {
	return true
}

// UnmarshalSecretActionOneOfRotateArbitrarySecretBody unmarshals an instance of SecretActionOneOfRotateArbitrarySecretBody from the specified map of raw messages.
func UnmarshalSecretActionOneOfRotateArbitrarySecretBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretActionOneOfRotateArbitrarySecretBody)
	err = core.UnmarshalPrimitive(m, "payload", &obj.Payload)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretActionOneOfRotateUsernamePasswordSecretBody : The request body of a `rotate` action.
// This model "extends" SecretActionOneOf
type SecretActionOneOfRotateUsernamePasswordSecretBody struct {
	// The new password to assign to a `username_password` secret.
	Password *string `json:"password" validate:"required"`
}

// NewSecretActionOneOfRotateUsernamePasswordSecretBody : Instantiate SecretActionOneOfRotateUsernamePasswordSecretBody (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretActionOneOfRotateUsernamePasswordSecretBody(password string) (model *SecretActionOneOfRotateUsernamePasswordSecretBody, err error) {
	model = &SecretActionOneOfRotateUsernamePasswordSecretBody{
		Password: core.StringPtr(password),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*SecretActionOneOfRotateUsernamePasswordSecretBody) isaSecretActionOneOf() bool {
	return true
}

// UnmarshalSecretActionOneOfRotateUsernamePasswordSecretBody unmarshals an instance of SecretActionOneOfRotateUsernamePasswordSecretBody from the specified map of raw messages.
func UnmarshalSecretActionOneOfRotateUsernamePasswordSecretBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretActionOneOfRotateUsernamePasswordSecretBody)
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretResourceArbitrarySecretResource : The base schema for secrets.
// This model "extends" SecretResource
type SecretResourceArbitrarySecretResource struct {
	// The MIME type that represents the secret.
	Type *string `json:"type,omitempty"`

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
	// Up to 30 labels can be created. Labels can be between 2-30 characters, including spaces. Special characters not
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

	// An array that contains metadata for each secret version.
	Versions []SecretVersion `json:"versions,omitempty"`

	// The date the secret material expires. The date format follows RFC 3339.
	//
	// You can set an expiration date on supported secret types at their creation. If you create a secret without
	// specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the
	// following secret types:
	//
	// - `arbitrary`
	// - `username_password`.
	ExpirationDate *strfmt.DateTime `json:"expiration_date,omitempty"`

	// The new secret data to assign to an `arbitrary` secret.
	Payload *string `json:"payload,omitempty"`

	SecretData interface{} `json:"secret_data,omitempty"`
}

// Constants associated with the SecretResourceArbitrarySecretResource.SecretType property.
// The secret type.
const (
	SecretResourceArbitrarySecretResourceSecretTypeArbitraryConst        = "arbitrary"
	SecretResourceArbitrarySecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	SecretResourceArbitrarySecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// NewSecretResourceArbitrarySecretResource : Instantiate SecretResourceArbitrarySecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretResourceArbitrarySecretResource(name string) (model *SecretResourceArbitrarySecretResource, err error) {
	model = &SecretResourceArbitrarySecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*SecretResourceArbitrarySecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalSecretResourceArbitrarySecretResource unmarshals an instance of SecretResourceArbitrarySecretResource from the specified map of raw messages.
func UnmarshalSecretResourceArbitrarySecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretResourceArbitrarySecretResource)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalSecretVersion)
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

// SecretResourceIamSecretResource : The base schema for secrets.
// This model "extends" SecretResource
type SecretResourceIamSecretResource struct {
	// The MIME type that represents the secret.
	Type *string `json:"type,omitempty"`

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
	// Up to 30 labels can be created. Labels can be between 2-30 characters, including spaces. Special characters not
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

	// An array that contains metadata for each secret version.
	Versions []SecretVersion `json:"versions,omitempty"`

	// The time-to-live (TTL) or lease duration to assign to generated credentials.
	//
	// For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be
	// either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m`
	// or `24h`.
	TTL interface{} `json:"ttl,omitempty"`

	// The access groups that define the capabilities of the service ID and API key that are generated for an
	// `iam_credentials` secret.
	//
	// **Tip:** To find the ID of an access group, go to **Manage > Access (IAM) > Access groups** in the IBM Cloud
	// console. Select the access group to inspect, and click **Details** to view its ID.
	AccessGroups []string `json:"access_groups,omitempty"`

	// The API key that is generated for this secret.
	//
	// After the secret reaches the end of its lease (see the `ttl` field), the API key is deleted automatically. If you
	// want to continue to use the same API key for future read operations, see the `reuse_api_key` field.
	APIKey *string `json:"api_key,omitempty"`

	// The service ID under which the API key (see the `api_key` field) is created. This service ID is added to the access
	// groups that you assign for this secret.
	ServiceID *string `json:"service_id,omitempty"`

	// Set to `true` to reuse the service ID and API key for this secret.
	//
	// Use this field to control whether to use the same service ID and API key for future read operations on this secret.
	// If set to `true`, the service reuses the current credentials. If set to `false`, a new service ID and API key is
	// generated each time that the secret is read or accessed.
	ReuseAPIKey *bool `json:"reuse_api_key,omitempty"`
}

// Constants associated with the SecretResourceIamSecretResource.SecretType property.
// The secret type.
const (
	SecretResourceIamSecretResourceSecretTypeArbitraryConst        = "arbitrary"
	SecretResourceIamSecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	SecretResourceIamSecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// NewSecretResourceIamSecretResource : Instantiate SecretResourceIamSecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretResourceIamSecretResource(name string) (model *SecretResourceIamSecretResource, err error) {
	model = &SecretResourceIamSecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*SecretResourceIamSecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalSecretResourceIamSecretResource unmarshals an instance of SecretResourceIamSecretResource from the specified map of raw messages.
func UnmarshalSecretResourceIamSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretResourceIamSecretResource)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalSecretVersion)
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
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reuse_api_key", &obj.ReuseAPIKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretResourceUsernamePasswordSecretResource : The base schema for secrets.
// This model "extends" SecretResource
type SecretResourceUsernamePasswordSecretResource struct {
	// The MIME type that represents the secret.
	Type *string `json:"type,omitempty"`

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
	// Up to 30 labels can be created. Labels can be between 2-30 characters, including spaces. Special characters not
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

	// An array that contains metadata for each secret version.
	Versions []SecretVersion `json:"versions,omitempty"`

	// The username to assign to this secret.
	Username *string `json:"username,omitempty"`

	// The password to assign to this secret.
	Password *string `json:"password,omitempty"`

	SecretData interface{} `json:"secret_data,omitempty"`

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

// Constants associated with the SecretResourceUsernamePasswordSecretResource.SecretType property.
// The secret type.
const (
	SecretResourceUsernamePasswordSecretResourceSecretTypeArbitraryConst        = "arbitrary"
	SecretResourceUsernamePasswordSecretResourceSecretTypeIamCredentialsConst   = "iam_credentials"
	SecretResourceUsernamePasswordSecretResourceSecretTypeUsernamePasswordConst = "username_password"
)

// NewSecretResourceUsernamePasswordSecretResource : Instantiate SecretResourceUsernamePasswordSecretResource (Generic Model Constructor)
func (*SecretsManagerV1) NewSecretResourceUsernamePasswordSecretResource(name string) (model *SecretResourceUsernamePasswordSecretResource, err error) {
	model = &SecretResourceUsernamePasswordSecretResource{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*SecretResourceUsernamePasswordSecretResource) isaSecretResource() bool {
	return true
}

// UnmarshalSecretResourceUsernamePasswordSecretResource unmarshals an instance of SecretResourceUsernamePasswordSecretResource from the specified map of raw messages.
func UnmarshalSecretResourceUsernamePasswordSecretResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretResourceUsernamePasswordSecretResource)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalSecretVersion)
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
