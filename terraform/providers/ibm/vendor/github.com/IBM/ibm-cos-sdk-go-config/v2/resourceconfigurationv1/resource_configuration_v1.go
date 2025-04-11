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
 * IBM OpenAPI SDK Code Generator Version: 3.96.1-5136e54a-20241108-203028
 */

// Package resourceconfigurationv1 : Operations and models for the ResourceConfigurationV1 service
package resourceconfigurationv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/ibm-cos-sdk-go-config/v2/common"
	"github.com/go-openapi/strfmt"
)

// ResourceConfigurationV1 : REST API used to configure Cloud Object Storage buckets.
//
// API Version: 1.0.0
type ResourceConfigurationV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://config.cloud-object-storage.cloud.ibm.com/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "resource_configuration"

// ResourceConfigurationV1Options : Service options
type ResourceConfigurationV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewResourceConfigurationV1UsingExternalConfig : constructs an instance of ResourceConfigurationV1 with passed in options and external configuration.
func NewResourceConfigurationV1UsingExternalConfig(options *ResourceConfigurationV1Options) (resourceConfiguration *ResourceConfigurationV1, err error) {
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

	resourceConfiguration, err = NewResourceConfigurationV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = resourceConfiguration.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = resourceConfiguration.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewResourceConfigurationV1 : constructs an instance of ResourceConfigurationV1 with passed in options.
func NewResourceConfigurationV1(options *ResourceConfigurationV1Options) (service *ResourceConfigurationV1, err error) {
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

	service = &ResourceConfigurationV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "resourceConfiguration" suitable for processing requests.
func (resourceConfiguration *ResourceConfigurationV1) Clone() *ResourceConfigurationV1 {
	if core.IsNil(resourceConfiguration) {
		return nil
	}
	clone := *resourceConfiguration
	clone.Service = resourceConfiguration.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (resourceConfiguration *ResourceConfigurationV1) SetServiceURL(url string) error {
	err := resourceConfiguration.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (resourceConfiguration *ResourceConfigurationV1) GetServiceURL() string {
	return resourceConfiguration.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (resourceConfiguration *ResourceConfigurationV1) SetDefaultHeaders(headers http.Header) {
	resourceConfiguration.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (resourceConfiguration *ResourceConfigurationV1) SetEnableGzipCompression(enableGzip bool) {
	resourceConfiguration.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (resourceConfiguration *ResourceConfigurationV1) GetEnableGzipCompression() bool {
	return resourceConfiguration.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (resourceConfiguration *ResourceConfigurationV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	resourceConfiguration.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (resourceConfiguration *ResourceConfigurationV1) DisableRetries() {
	resourceConfiguration.Service.DisableRetries()
}

// CreateBackupPolicy : Add a new backup policy to the COS Bucket
// Attach a new Backup Policy on a bucket.
//
// This request results in the creation of a single, new RecoveryRange on the destination BackupVault.
//
// Deletion and re-creation of a BackupPolicy to the same BackupVault destination will generate a new RecoveryRange.
//
// The following shall be validated. Any failure to validate shall cause a HTTP 400 to be returned.
//
//   * the user has `cloud-object-storage.bucket.post_backup_policy` permissions on the source-bucket
//   * the source-bucket must have `cloud-object-storage.backup_vault.sync` permissions on the Backup Vault
//   * the source-bucket must have versioning-on
//   * the Backup Vault must exist and be able to be contacted by the source-bucket
//   * the source-bucket must not have an existing BackupPolicy targeting the Backup Vault
//   * the source-bucket must not have a BackupPolicy with the same policy_name
//   * the source-bucket must have fewer than 3 total BackupPolicies
//
// This request generates the "cloud-object-storage.bucket-backup-policy.create" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) CreateBackupPolicy(createBackupPolicyOptions *CreateBackupPolicyOptions) (result *BackupPolicy, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.CreateBackupPolicyWithContext(context.Background(), createBackupPolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateBackupPolicyWithContext is an alternate form of the CreateBackupPolicy method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) CreateBackupPolicyWithContext(ctx context.Context, createBackupPolicyOptions *CreateBackupPolicyOptions) (result *BackupPolicy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createBackupPolicyOptions, "createBackupPolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createBackupPolicyOptions, "createBackupPolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"bucket": *createBackupPolicyOptions.Bucket,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/buckets/{bucket}/backup_policies`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createBackupPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "CreateBackupPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createBackupPolicyOptions.MD5 != nil {
		builder.AddHeader("MD5", fmt.Sprint(*createBackupPolicyOptions.MD5))
	}

	body := make(map[string]interface{})
	if createBackupPolicyOptions.PolicyName != nil {
		body["policy_name"] = createBackupPolicyOptions.PolicyName
	}
	if createBackupPolicyOptions.TargetBackupVaultCrn != nil {
		body["target_backup_vault_crn"] = createBackupPolicyOptions.TargetBackupVaultCrn
	}
	if createBackupPolicyOptions.BackupType != nil {
		body["backup_type"] = createBackupPolicyOptions.BackupType
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_backup_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBackupPolicy)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListBackupPolicies : List BackupPolicies
// Get all backup policies on a bucket.
//
// Requires that the user has `cloud-object-storage.bucket.list_backup_policies` permissions on the source bucket.
//
// This request generates the "cloud-object-storage.bucket-backup-policy.list" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) ListBackupPolicies(listBackupPoliciesOptions *ListBackupPoliciesOptions) (result *BackupPolicyCollection, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.ListBackupPoliciesWithContext(context.Background(), listBackupPoliciesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListBackupPoliciesWithContext is an alternate form of the ListBackupPolicies method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) ListBackupPoliciesWithContext(ctx context.Context, listBackupPoliciesOptions *ListBackupPoliciesOptions) (result *BackupPolicyCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listBackupPoliciesOptions, "listBackupPoliciesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listBackupPoliciesOptions, "listBackupPoliciesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"bucket": *listBackupPoliciesOptions.Bucket,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/buckets/{bucket}/backup_policies`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listBackupPoliciesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "ListBackupPolicies")
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_backup_policies", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBackupPolicyCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetBackupPolicy : Get BackupPolicy
// Read a specific backup policy on a bucket.
//
// Requires that the user has `cloud-object-storage.bucket.get_backup_policy` permissions on the bucket.
//
// This request generates the "cloud-object-storage.bucket-backup-policy.read" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) GetBackupPolicy(getBackupPolicyOptions *GetBackupPolicyOptions) (result *BackupPolicy, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.GetBackupPolicyWithContext(context.Background(), getBackupPolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetBackupPolicyWithContext is an alternate form of the GetBackupPolicy method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) GetBackupPolicyWithContext(ctx context.Context, getBackupPolicyOptions *GetBackupPolicyOptions) (result *BackupPolicy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBackupPolicyOptions, "getBackupPolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getBackupPolicyOptions, "getBackupPolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"bucket": *getBackupPolicyOptions.Bucket,
		"policy_id": *getBackupPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/buckets/{bucket}/backup_policies/{policy_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getBackupPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "GetBackupPolicy")
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_backup_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBackupPolicy)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteBackupPolicy : Delete a BackupPolicy
// Delete a specific BackupPolicy.
//
// Requires that the user has `cloud-object-storage.bucket.delete_backup_policy` permissions on the bucket.
//
// This request generates the "cloud-object-storage.bucket-backup-policy.delete" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) DeleteBackupPolicy(deleteBackupPolicyOptions *DeleteBackupPolicyOptions) (response *core.DetailedResponse, err error) {
	response, err = resourceConfiguration.DeleteBackupPolicyWithContext(context.Background(), deleteBackupPolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteBackupPolicyWithContext is an alternate form of the DeleteBackupPolicy method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) DeleteBackupPolicyWithContext(ctx context.Context, deleteBackupPolicyOptions *DeleteBackupPolicyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteBackupPolicyOptions, "deleteBackupPolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteBackupPolicyOptions, "deleteBackupPolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"bucket": *deleteBackupPolicyOptions.Bucket,
		"policy_id": *deleteBackupPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/buckets/{bucket}/backup_policies/{policy_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteBackupPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "DeleteBackupPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = resourceConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_backup_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListBackupVaults : list BackupVaults
// Returns a list of BackupVault CRNs owned by the account.
//
// Requires that the user has `cloud-object-storage.backup_vault.list_account_backup_vaults` permissions for the
// account.
//
// This request generates the "cloud-object-storage.backup-vault.list" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) ListBackupVaults(listBackupVaultsOptions *ListBackupVaultsOptions) (result *BackupVaultCollection, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.ListBackupVaultsWithContext(context.Background(), listBackupVaultsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListBackupVaultsWithContext is an alternate form of the ListBackupVaults method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) ListBackupVaultsWithContext(ctx context.Context, listBackupVaultsOptions *ListBackupVaultsOptions) (result *BackupVaultCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listBackupVaultsOptions, "listBackupVaultsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listBackupVaultsOptions, "listBackupVaultsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listBackupVaultsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "ListBackupVaults")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("service_instance_id", fmt.Sprint(*listBackupVaultsOptions.ServiceInstanceID))
	if listBackupVaultsOptions.Token != nil {
		builder.AddQuery("token", fmt.Sprint(*listBackupVaultsOptions.Token))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_backup_vaults", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBackupVaultCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateBackupVault : create a BackupVault
// Creates a BackupVault.
//
// Requires that the user has `cloud-object-storage.backup_vault.post_backup_vault` permissions for the account.
//
// Certain fields will be returned only if the user has specific permissions:
//   - `activity_tracking` requires `cloud-object-storage.backup_vault.put_activity_tracking`
//   - `metrics_monitoring` requires `cloud-object-storage.backup_vault.put_metrics_monitoring`
//
// This request generates the "cloud-object-storage.backup-vault.create" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) CreateBackupVault(createBackupVaultOptions *CreateBackupVaultOptions) (result *BackupVault, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.CreateBackupVaultWithContext(context.Background(), createBackupVaultOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateBackupVaultWithContext is an alternate form of the CreateBackupVault method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) CreateBackupVaultWithContext(ctx context.Context, createBackupVaultOptions *CreateBackupVaultOptions) (result *BackupVault, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createBackupVaultOptions, "createBackupVaultOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createBackupVaultOptions, "createBackupVaultOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createBackupVaultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "CreateBackupVault")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	builder.AddQuery("service_instance_id", fmt.Sprint(*createBackupVaultOptions.ServiceInstanceID))

	body := make(map[string]interface{})
	if createBackupVaultOptions.BackupVaultName != nil {
		body["backup_vault_name"] = createBackupVaultOptions.BackupVaultName
	}
	if createBackupVaultOptions.Region != nil {
		body["region"] = createBackupVaultOptions.Region
	}
	if createBackupVaultOptions.ActivityTracking != nil {
		body["activity_tracking"] = createBackupVaultOptions.ActivityTracking
	}
	if createBackupVaultOptions.MetricsMonitoring != nil {
		body["metrics_monitoring"] = createBackupVaultOptions.MetricsMonitoring
	}
	if createBackupVaultOptions.SseKpCustomerRootKeyCrn != nil {
		body["sse_kp_customer_root_key_crn"] = createBackupVaultOptions.SseKpCustomerRootKeyCrn
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_backup_vault", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBackupVault)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetBackupVault : get the config for a Backup Vault
// Gets configuration information for a Backup Vault.
//
// Requires that the user has `cloud-object-storage.backup_vault.get_basic` permissions on the backup vault.
//
// Certain fields will be returned only if the user has specific permissions:
//   - `activity_tracking` requires `cloud-object-storage.backup_vault.get_activity_tracking`
//   - `metrics_monitoring` requires `cloud-object-storage.backup_vault.get_metrics_monitoring`
//   - `sse_kp_customer_root_key_crn` requires `cloud-object-storage.backup_vault.get_crk_id`
//
// This request generates the "cloud-object-storage.backup-vault-configuration.read" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) GetBackupVault(getBackupVaultOptions *GetBackupVaultOptions) (result *BackupVault, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.GetBackupVaultWithContext(context.Background(), getBackupVaultOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetBackupVaultWithContext is an alternate form of the GetBackupVault method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) GetBackupVaultWithContext(ctx context.Context, getBackupVaultOptions *GetBackupVaultOptions) (result *BackupVault, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBackupVaultOptions, "getBackupVaultOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getBackupVaultOptions, "getBackupVaultOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"backup_vault_name": *getBackupVaultOptions.BackupVaultName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults/{backup_vault_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getBackupVaultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "GetBackupVault")
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_backup_vault", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBackupVault)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateBackupVault : Update the config on a Backup Vault
// Update the Backup Vault config via JSON Merge Patch update semantics.
//
// In particular, note that providing an empty object (`{}`) to either field in the request body will remove any
// existing configuration.
//
// Requires that the user has `cloud-object-storage.backup_vault.get_basic` permissions on the backup vault.
//
// Certain fields can be modified only if the user has specific permissions:
//   - `activity_tracking` requires `cloud-object-storage.backup_vault.put_activity_tracking`
//   - `metrics_monitoring` requires `cloud-object-storage.backup_vault.put_metrics_monitoring`
//
// This request generates the "cloud-object-storage.backup-vault-configuration.update" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) UpdateBackupVault(updateBackupVaultOptions *UpdateBackupVaultOptions) (result *BackupVault, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.UpdateBackupVaultWithContext(context.Background(), updateBackupVaultOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateBackupVaultWithContext is an alternate form of the UpdateBackupVault method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) UpdateBackupVaultWithContext(ctx context.Context, updateBackupVaultOptions *UpdateBackupVaultOptions) (result *BackupVault, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateBackupVaultOptions, "updateBackupVaultOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateBackupVaultOptions, "updateBackupVaultOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"backup_vault_name": *updateBackupVaultOptions.BackupVaultName,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults/{backup_vault_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateBackupVaultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "UpdateBackupVault")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateBackupVaultOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateBackupVaultOptions.IfMatch))
	}

	_, err = builder.SetBodyContentJSON(updateBackupVaultOptions.BackupVaultPatch)
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_backup_vault", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBackupVault)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteBackupVault : Delete an empty Backup Vault
// Delete the Backup Vault.
//
// Requires that the BackupVault not contain any RecoveryRanges.  Requires that the user has
// `cloud-object-storage.backup_vault.delete_backup_vault` permissions for the account.
//
// This request generates the "cloud-object-storage.backup-vault.delete" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) DeleteBackupVault(deleteBackupVaultOptions *DeleteBackupVaultOptions) (response *core.DetailedResponse, err error) {
	response, err = resourceConfiguration.DeleteBackupVaultWithContext(context.Background(), deleteBackupVaultOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteBackupVaultWithContext is an alternate form of the DeleteBackupVault method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) DeleteBackupVaultWithContext(ctx context.Context, deleteBackupVaultOptions *DeleteBackupVaultOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteBackupVaultOptions, "deleteBackupVaultOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteBackupVaultOptions, "deleteBackupVaultOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"backup_vault_name": *deleteBackupVaultOptions.BackupVaultName,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults/{backup_vault_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteBackupVaultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "DeleteBackupVault")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = resourceConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_backup_vault", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetBucketConfig : Returns metadata for the specified bucket
// Returns metadata for the specified bucket.
func (resourceConfiguration *ResourceConfigurationV1) GetBucketConfig(getBucketConfigOptions *GetBucketConfigOptions) (result *Bucket, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.GetBucketConfigWithContext(context.Background(), getBucketConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetBucketConfigWithContext is an alternate form of the GetBucketConfig method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) GetBucketConfigWithContext(ctx context.Context, getBucketConfigOptions *GetBucketConfigOptions) (result *Bucket, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBucketConfigOptions, "getBucketConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getBucketConfigOptions, "getBucketConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"bucket": *getBucketConfigOptions.Bucket,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/b/{bucket}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getBucketConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "GetBucketConfig")
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "getBucketConfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBucket)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateBucketConfig : Make changes to a bucket's configuration
// Updates a bucket using [JSON Merge Patch](https://tools.ietf.org/html/rfc7396). This request is used to add
// functionality (like an IP access filter) or to update existing parameters.  **Primitives are overwritten and replaced
// in their entirety. It is not possible to append a new (or to delete a specific) value to an array.**  Arrays can be
// cleared by updating the parameter with an empty array `[]`. A `PATCH` operation only updates specified mutable
// fields. Please don't use `PATCH` trying to update the number of objects in a bucket, any timestamps, or other
// non-mutable fields.
func (resourceConfiguration *ResourceConfigurationV1) UpdateBucketConfig(updateBucketConfigOptions *UpdateBucketConfigOptions) (response *core.DetailedResponse, err error) {
	response, err = resourceConfiguration.UpdateBucketConfigWithContext(context.Background(), updateBucketConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateBucketConfigWithContext is an alternate form of the UpdateBucketConfig method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) UpdateBucketConfigWithContext(ctx context.Context, updateBucketConfigOptions *UpdateBucketConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateBucketConfigOptions, "updateBucketConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateBucketConfigOptions, "updateBucketConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"bucket": *updateBucketConfigOptions.Bucket,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/b/{bucket}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateBucketConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "UpdateBucketConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateBucketConfigOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateBucketConfigOptions.IfMatch))
	}

	if updateBucketConfigOptions.BucketPatch != nil {
		_, err = builder.SetBodyContentJSON(updateBucketConfigOptions.BucketPatch)
		if err != nil {
			err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = resourceConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "updateBucketConfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListRecoveryRanges : List RecoveryRanges on a backup vault
// List RecoveryRanges on a backup vault. Lists all available ranges for all source resources by default. The
// `?source_resource_crn` query parameter will limit the list to only ranges for the specified resource.
//
//
// Requires the user have `cloud-object-storage.backup_vault.list_recovery_ranges` permissions to the Backup Vault.
//
// This request generates the "cloud-object-storage.backup-recovery-range.list" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) ListRecoveryRanges(listRecoveryRangesOptions *ListRecoveryRangesOptions) (result *RecoveryRangeCollection, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.ListRecoveryRangesWithContext(context.Background(), listRecoveryRangesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListRecoveryRangesWithContext is an alternate form of the ListRecoveryRanges method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) ListRecoveryRangesWithContext(ctx context.Context, listRecoveryRangesOptions *ListRecoveryRangesOptions) (result *RecoveryRangeCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listRecoveryRangesOptions, "listRecoveryRangesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listRecoveryRangesOptions, "listRecoveryRangesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"backup_vault_name": *listRecoveryRangesOptions.BackupVaultName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults/{backup_vault_name}/recovery_ranges`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listRecoveryRangesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "ListRecoveryRanges")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listRecoveryRangesOptions.SourceResourceCrn != nil {
		builder.AddQuery("source_resource_crn", fmt.Sprint(*listRecoveryRangesOptions.SourceResourceCrn))
	}
	if listRecoveryRangesOptions.Token != nil {
		builder.AddQuery("token", fmt.Sprint(*listRecoveryRangesOptions.Token))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_recovery_ranges", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRecoveryRangeCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetSourceResourceRecoveryRange : get RecoveryRange info
// Get info for a specific RecoveryRange.
//
// Requires the user have `cloud-object-storage.backup_vault.get_recovery_range` permissions to the Backup Vault.
//
// This request generates the "cloud-object-storage.backup-recovery-range.read" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) GetSourceResourceRecoveryRange(getSourceResourceRecoveryRangeOptions *GetSourceResourceRecoveryRangeOptions) (result *RecoveryRange, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.GetSourceResourceRecoveryRangeWithContext(context.Background(), getSourceResourceRecoveryRangeOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetSourceResourceRecoveryRangeWithContext is an alternate form of the GetSourceResourceRecoveryRange method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) GetSourceResourceRecoveryRangeWithContext(ctx context.Context, getSourceResourceRecoveryRangeOptions *GetSourceResourceRecoveryRangeOptions) (result *RecoveryRange, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSourceResourceRecoveryRangeOptions, "getSourceResourceRecoveryRangeOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getSourceResourceRecoveryRangeOptions, "getSourceResourceRecoveryRangeOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"backup_vault_name": *getSourceResourceRecoveryRangeOptions.BackupVaultName,
		"recovery_range_id": *getSourceResourceRecoveryRangeOptions.RecoveryRangeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults/{backup_vault_name}/recovery_ranges/{recovery_range_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getSourceResourceRecoveryRangeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "GetSourceResourceRecoveryRange")
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_source_resource_recovery_range", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRecoveryRange)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateRestore : Initiate a Restore
// Initiates a restore operation against some RecoveryRange to some destination bucket.
//
// The following shall be validated. Any failure to validate shall cause a HTTP 400 to be returned.
//
//   * The specified RecoveryRange must exist
//   * The restore time must be within the RecoveryRange
//   * the user has `cloud-object-storage.backup-vault.post_restore` permissions on the backup-vault
//   * the target-bucket must exist and be able to be contacted by the Backup Vault
//   * target-bucket must have versioning-on
//   * the Backup Vault must have `cloud-object-storage.bucket.restore_sync` permissions on the target-bucket
//
// This request generates the "cloud-object-storage.backup-restore.create" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) CreateRestore(createRestoreOptions *CreateRestoreOptions) (result *Restore, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.CreateRestoreWithContext(context.Background(), createRestoreOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateRestoreWithContext is an alternate form of the CreateRestore method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) CreateRestoreWithContext(ctx context.Context, createRestoreOptions *CreateRestoreOptions) (result *Restore, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRestoreOptions, "createRestoreOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createRestoreOptions, "createRestoreOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"backup_vault_name": *createRestoreOptions.BackupVaultName,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults/{backup_vault_name}/restores`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createRestoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "CreateRestore")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createRestoreOptions.RecoveryRangeID != nil {
		body["recovery_range_id"] = createRestoreOptions.RecoveryRangeID
	}
	if createRestoreOptions.RestoreType != nil {
		body["restore_type"] = createRestoreOptions.RestoreType
	}
	if createRestoreOptions.RestorePointInTime != nil {
		body["restore_point_in_time"] = createRestoreOptions.RestorePointInTime
	}
	if createRestoreOptions.TargetResourceCrn != nil {
		body["target_resource_crn"] = createRestoreOptions.TargetResourceCrn
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_restore", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRestore)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListRestores : List Restores
// List all current and complete restores.
//
// Requires that the user have `cloud-object-storage.backup_vault.list_restores` permission on the backup vault.
//
// This request generates the "cloud-object-storage.backup-restore.list" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) ListRestores(listRestoresOptions *ListRestoresOptions) (result *RestoreCollection, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.ListRestoresWithContext(context.Background(), listRestoresOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListRestoresWithContext is an alternate form of the ListRestores method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) ListRestoresWithContext(ctx context.Context, listRestoresOptions *ListRestoresOptions) (result *RestoreCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listRestoresOptions, "listRestoresOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listRestoresOptions, "listRestoresOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"backup_vault_name": *listRestoresOptions.BackupVaultName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults/{backup_vault_name}/restores`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listRestoresOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "ListRestores")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listRestoresOptions.Token != nil {
		builder.AddQuery("token", fmt.Sprint(*listRestoresOptions.Token))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_restores", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRestoreCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetRestore : Get Restore
// Introspect on a specific restore.
//
// Requires that the user have `cloud-object-storage.backup_vault.get_restore` permission on the backup vault.
//
// This request generates the "cloud-object-storage.backup-restore.read" ActivityTracking event.
func (resourceConfiguration *ResourceConfigurationV1) GetRestore(getRestoreOptions *GetRestoreOptions) (result *Restore, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.GetRestoreWithContext(context.Background(), getRestoreOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetRestoreWithContext is an alternate form of the GetRestore method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) GetRestoreWithContext(ctx context.Context, getRestoreOptions *GetRestoreOptions) (result *Restore, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRestoreOptions, "getRestoreOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getRestoreOptions, "getRestoreOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"backup_vault_name": *getRestoreOptions.BackupVaultName,
		"restore_id": *getRestoreOptions.RestoreID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/backup_vaults/{backup_vault_name}/restores/{restore_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getRestoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "GetRestore")
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
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_restore", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRestore)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.0")
}

// ActivityTracking : Enables sending log data to IBM Cloud Activity Tracker Event Routing to provide visibility into bucket management,
// object read and write events. (Recommended) When the `activity_tracker_crn` is not populated, then enabled events are
// sent to the Activity Tracker Event Routing instance at the container's location unless otherwise specified in the
// Activity Tracker Event Routing Event Routing service configuration. (Legacy) When the `activity_tracker_crn` is
// populated, then enabled events are sent to the Activity Tracker Event Routing instance specified.
type ActivityTracking struct {
	// If set to `true`, all object read events (i.e. downloads) will be sent to Activity Tracker Event Routing.
	ReadDataEvents *bool `json:"read_data_events,omitempty"`

	// If set to `true`, all object write events (i.e. uploads) will be sent to Activity Tracker Event Routing.
	WriteDataEvents *bool `json:"write_data_events,omitempty"`

	// When the `activity_tracker_crn` is not populated, then enabled events are sent to the Activity Tracker Event Routing
	// instance associated to the container's location unless otherwise specified in the Activity Tracker Event Routing
	// Event Routing service configuration. If `activity_tracker_crn` is populated, then enabled events are sent to the
	// Activity Tracker Event Routing instance specified and bucket management events are always enabled.
	ActivityTrackerCrn *string `json:"activity_tracker_crn,omitempty"`

	// This field only applies if `activity_tracker_crn` is not populated. If set to `true`, all bucket management events
	// will be sent to Activity Tracker Event Routing.
	ManagementEvents *bool `json:"management_events,omitempty"`
}

// UnmarshalActivityTracking unmarshals an instance of ActivityTracking from the specified map of raw messages.
func UnmarshalActivityTracking(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActivityTracking)
	err = core.UnmarshalPrimitive(m, "read_data_events", &obj.ReadDataEvents)
	if err != nil {
		err = core.SDKErrorf(err, "", "read_data_events-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "write_data_events", &obj.WriteDataEvents)
	if err != nil {
		err = core.SDKErrorf(err, "", "write_data_events-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "activity_tracker_crn", &obj.ActivityTrackerCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity_tracker_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "management_events", &obj.ManagementEvents)
	if err != nil {
		err = core.SDKErrorf(err, "", "management_events-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the ActivityTracking
func (activityTracking *ActivityTracking) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(activityTracking.ReadDataEvents) {
		_patch["read_data_events"] = activityTracking.ReadDataEvents
	}
	if !core.IsNil(activityTracking.WriteDataEvents) {
		_patch["write_data_events"] = activityTracking.WriteDataEvents
	}
	if !core.IsNil(activityTracking.ActivityTrackerCrn) {
		_patch["activity_tracker_crn"] = activityTracking.ActivityTrackerCrn
	}
	if !core.IsNil(activityTracking.ManagementEvents) {
		_patch["management_events"] = activityTracking.ManagementEvents
	}

	return
}

// BackupPolicy : The current backup coverage for a COS Bucket.
type BackupPolicy struct {
	// The name granted to the policy. Validation :
	//   * chars limited to alphanumeric, underscore, hyphen and period.
	PolicyName *string `json:"policy_name" validate:"required"`

	// The CRN for a COS BackupVault.
	TargetBackupVaultCrn *string `json:"target_backup_vault_crn" validate:"required"`

	// The type of backup to support. For LA+GA this is limited to "continuous".
	BackupType *string `json:"backup_type" validate:"required"`

	// A UUID that uniquely identifies a resource.
	PolicyID *string `json:"policy_id" validate:"required"`

	// The current status of the backup policy.
	//
	// pending : the policy has been received and has begun processing. initializing : pre-existing objects are being sync
	// to the backup vault. active : the policy is active and healthy. action_needed : the policy is unhealthy and requires
	// some intervention to recover degraded : the policy is unhealthy failed : the policy has failed unrecoverably.
	PolicyStatus *string `json:"policy_status" validate:"required"`

	// Reports percent-doneness of init.
	//
	// Only present when policy_status=INITIALING.
	InitialSyncProgress *int64 `json:"initial_sync_progress,omitempty"`

	// reports error cause. Only present when policy_status=ERROR/FAILED.
	ErrorCause *string `json:"error_cause,omitempty"`
}

// Constants associated with the BackupPolicy.BackupType property.
// The type of backup to support. For LA+GA this is limited to "continuous".
const (
	BackupPolicy_BackupType_Continuous = "continuous"
)

// Constants associated with the BackupPolicy.PolicyStatus property.
// The current status of the backup policy.
//
// pending : the policy has been received and has begun processing. initializing : pre-existing objects are being sync
// to the backup vault. active : the policy is active and healthy. action_needed : the policy is unhealthy and requires
// some intervention to recover degraded : the policy is unhealthy failed : the policy has failed unrecoverably.
const (
	BackupPolicy_PolicyStatus_ActionNeeded = "action_needed"
	BackupPolicy_PolicyStatus_Active = "active"
	BackupPolicy_PolicyStatus_Degraded = "degraded"
	BackupPolicy_PolicyStatus_Failed = "failed"
	BackupPolicy_PolicyStatus_Initializing = "initializing"
	BackupPolicy_PolicyStatus_Pending = "pending"
)

// UnmarshalBackupPolicy unmarshals an instance of BackupPolicy from the specified map of raw messages.
func UnmarshalBackupPolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BackupPolicy)
	err = core.UnmarshalPrimitive(m, "policy_name", &obj.PolicyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_backup_vault_crn", &obj.TargetBackupVaultCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_backup_vault_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "backup_type", &obj.BackupType)
	if err != nil {
		err = core.SDKErrorf(err, "", "backup_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_id", &obj.PolicyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_status", &obj.PolicyStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "initial_sync_progress", &obj.InitialSyncProgress)
	if err != nil {
		err = core.SDKErrorf(err, "", "initial_sync_progress-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "error_cause", &obj.ErrorCause)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_cause-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BackupPolicyCollection : A collection of backup policies.
type BackupPolicyCollection struct {
	// A collection of backup policies.
	BackupPolicies []BackupPolicy `json:"backup_policies,omitempty"`
}

// UnmarshalBackupPolicyCollection unmarshals an instance of BackupPolicyCollection from the specified map of raw messages.
func UnmarshalBackupPolicyCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BackupPolicyCollection)
	err = core.UnmarshalModel(m, "backup_policies", &obj.BackupPolicies, UnmarshalBackupPolicy)
	if err != nil {
		err = core.SDKErrorf(err, "", "backup_policies-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BackupVault : Metadata associated with a backup vault.
type BackupVault struct {
	// Activity Tracking configuration. An empty object (`{}`) indicates no configuration, and no events will be sent (This
	// is the same behavior as `{"management_events":false}`). Note that read/write events cannot be enabled, and events
	// cannot be routed to a non-default Activity Tracker instance.
	ActivityTracking *BackupVaultActivityTracking `json:"activity_tracking,omitempty"`

	// Metrics Monitoring configuration. An empty object (`{}`) indicates no configuration, and no metrics will be
	// collected (This is the same behavior as `{"usage_metrics_enabled":false}`). Note that request metrics cannot be
	// enabled, and metrics cannot be routed to a non-default metrics router instance.
	MetricsMonitoring *BackupVaultMetricsMonitoring `json:"metrics_monitoring,omitempty"`

	// The name given to a Bucket.
	//
	// Bucket names must be between 3 and 63 characters long must be made of lowercase letters, numbers, dots (periods),
	// and dashes (hyphens). Bucket names must begin and end with a lowercase letter or number. Bucket names canâ€t contain
	// consecutive dots or dashes. Bucket names that resemble IP addresses are not allowed.
	//
	// Bucket and BackupVault names exist in a global global namespace and therefore must be unique.
	BackupVaultName *string `json:"backup_vault_name" validate:"required"`

	// the region in which this backup-vault should be created within.
	Region *string `json:"region" validate:"required"`

	// The CRN for a KeyProtect root key.
	SseKpCustomerRootKeyCrn *string `json:"sse_kp_customer_root_key_crn,omitempty"`

	// CRN of the backup-vault.
	Crn *string `json:"crn,omitempty"`

	// A COS ServiceInstance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn,omitempty"`

	// creation time of the backup-vault. Returns "YYYY-MM-DDTHH:mm:ss.sssZ" timestamp format.
	TimeCreated *strfmt.DateTime `json:"time_created,omitempty"`

	// time of last update to the backup-vault Returns "YYYY-MM-DDTHH:mm:ss.sssZ" timestamp format.
	TimeUpdated *strfmt.DateTime `json:"time_updated,omitempty"`

	// byte useage of the backup-vault. This should include all usage, including non-current versions.
	BytesUsed *int64 `json:"bytes_used,omitempty"`
}

// UnmarshalBackupVault unmarshals an instance of BackupVault from the specified map of raw messages.
func UnmarshalBackupVault(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BackupVault)
	err = core.UnmarshalModel(m, "activity_tracking", &obj.ActivityTracking, UnmarshalBackupVaultActivityTracking)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity_tracking-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metrics_monitoring", &obj.MetricsMonitoring, UnmarshalBackupVaultMetricsMonitoring)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics_monitoring-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "backup_vault_name", &obj.BackupVaultName)
	if err != nil {
		err = core.SDKErrorf(err, "", "backup_vault_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sse_kp_customer_root_key_crn", &obj.SseKpCustomerRootKeyCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sse_kp_customer_root_key_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_crn", &obj.ServiceInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "time_created", &obj.TimeCreated)
	if err != nil {
		err = core.SDKErrorf(err, "", "time_created-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "time_updated", &obj.TimeUpdated)
	if err != nil {
		err = core.SDKErrorf(err, "", "time_updated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bytes_used", &obj.BytesUsed)
	if err != nil {
		err = core.SDKErrorf(err, "", "bytes_used-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BackupVaultActivityTracking : Activity Tracking configuration. An empty object (`{}`) indicates no configuration, and no events will be sent (This
// is the same behavior as `{"management_events":false}`). Note that read/write events cannot be enabled, and events
// cannot be routed to a non-default Activity Tracker instance.
type BackupVaultActivityTracking struct {
	// Whether to send notifications for management events on the BackupVault.
	ManagementEvents *bool `json:"management_events,omitempty"`
}

// UnmarshalBackupVaultActivityTracking unmarshals an instance of BackupVaultActivityTracking from the specified map of raw messages.
func UnmarshalBackupVaultActivityTracking(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BackupVaultActivityTracking)
	err = core.UnmarshalPrimitive(m, "management_events", &obj.ManagementEvents)
	if err != nil {
		err = core.SDKErrorf(err, "", "management_events-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the BackupVaultActivityTracking
func (backupVaultActivityTracking *BackupVaultActivityTracking) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(backupVaultActivityTracking.ManagementEvents) {
		_patch["management_events"] = backupVaultActivityTracking.ManagementEvents
	}

	return
}

// BackupVaultCollection : A listing of backup vaults.
type BackupVaultCollection struct {
	// Pagination response body.
	Next *NextPagination `json:"next,omitempty"`

	// List of Backup Vaults. If no Backup Vaults exist, this array will be empty.
	BackupVaults []string `json:"backup_vaults,omitempty"`
}

// UnmarshalBackupVaultCollection unmarshals an instance of BackupVaultCollection from the specified map of raw messages.
func UnmarshalBackupVaultCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BackupVaultCollection)
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextPagination)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "backup_vaults", &obj.BackupVaults)
	if err != nil {
		err = core.SDKErrorf(err, "", "backup_vaults-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *BackupVaultCollection) GetNextToken() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Token, nil
}

// BackupVaultMetricsMonitoring : Metrics Monitoring configuration. An empty object (`{}`) indicates no configuration, and no metrics will be collected
// (This is the same behavior as `{"usage_metrics_enabled":false}`). Note that request metrics cannot be enabled, and
// metrics cannot be routed to a non-default metrics router instance.
type BackupVaultMetricsMonitoring struct {
	// Whether usage metrics are collected for this BackupVault.
	UsageMetricsEnabled *bool `json:"usage_metrics_enabled,omitempty"`
}

// UnmarshalBackupVaultMetricsMonitoring unmarshals an instance of BackupVaultMetricsMonitoring from the specified map of raw messages.
func UnmarshalBackupVaultMetricsMonitoring(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BackupVaultMetricsMonitoring)
	err = core.UnmarshalPrimitive(m, "usage_metrics_enabled", &obj.UsageMetricsEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "usage_metrics_enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the BackupVaultMetricsMonitoring
func (backupVaultMetricsMonitoring *BackupVaultMetricsMonitoring) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(backupVaultMetricsMonitoring.UsageMetricsEnabled) {
		_patch["usage_metrics_enabled"] = backupVaultMetricsMonitoring.UsageMetricsEnabled
	}

	return
}

// BackupVaultPatch : Metadata elements on a backup vault that can be updated.
type BackupVaultPatch struct {
	// Activity Tracking configuration. An empty object (`{}`) indicates no configuration, and no events will be sent (This
	// is the same behavior as `{"management_events":false}`). Note that read/write events cannot be enabled, and events
	// cannot be routed to a non-default Activity Tracker instance.
	ActivityTracking *BackupVaultActivityTracking `json:"activity_tracking,omitempty"`

	// Metrics Monitoring configuration. An empty object (`{}`) indicates no configuration, and no metrics will be
	// collected (This is the same behavior as `{"usage_metrics_enabled":false}`). Note that request metrics cannot be
	// enabled, and metrics cannot be routed to a non-default metrics router instance.
	MetricsMonitoring *BackupVaultMetricsMonitoring `json:"metrics_monitoring,omitempty"`
}

// UnmarshalBackupVaultPatch unmarshals an instance of BackupVaultPatch from the specified map of raw messages.
func UnmarshalBackupVaultPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BackupVaultPatch)
	err = core.UnmarshalModel(m, "activity_tracking", &obj.ActivityTracking, UnmarshalBackupVaultActivityTracking)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity_tracking-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metrics_monitoring", &obj.MetricsMonitoring, UnmarshalBackupVaultMetricsMonitoring)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics_monitoring-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the BackupVaultPatch
func (backupVaultPatch *BackupVaultPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(backupVaultPatch.ActivityTracking) {
		_patch["activity_tracking"] = backupVaultPatch.ActivityTracking.asPatch()
	}
	if !core.IsNil(backupVaultPatch.MetricsMonitoring) {
		_patch["metrics_monitoring"] = backupVaultPatch.MetricsMonitoring.asPatch()
	}

	return
}

// Bucket : A bucket.
type Bucket struct {
	// The name of the bucket. Non-mutable.
	Name *string `json:"name,omitempty"`

	// The service instance that holds the bucket. Non-mutable.
	Crn *string `json:"crn,omitempty"`

	// The service instance that holds the bucket. Non-mutable.
	ServiceInstanceID *string `json:"service_instance_id,omitempty"`

	// The service instance that holds the bucket. Non-mutable.
	ServiceInstanceCrn *string `json:"service_instance_crn,omitempty"`

	// The creation time of the bucket in RFC 3339 format. Non-mutable.
	TimeCreated *strfmt.DateTime `json:"time_created,omitempty"`

	// The modification time of the bucket in RFC 3339 format. Non-mutable.
	TimeUpdated *strfmt.DateTime `json:"time_updated,omitempty"`

	// Total number of objects in the bucket. Non-mutable.
	ObjectCount *int64 `json:"object_count,omitempty"`

	// Total size of all objects in the bucket. Non-mutable.
	BytesUsed *int64 `json:"bytes_used,omitempty"`

	// Number of non-current object versions in the bucket. Non-mutable.
	NoncurrentObjectCount *int64 `json:"noncurrent_object_count,omitempty"`

	// Total size of all non-current object versions in the bucket. Non-mutable.
	NoncurrentBytesUsed *int64 `json:"noncurrent_bytes_used,omitempty"`

	// Total number of delete markers in the bucket. Non-mutable.
	DeleteMarkerCount *int64 `json:"delete_marker_count,omitempty"`

	// An access control mechanism based on the network (IP address) where request originated. Requests not originating
	// from IP addresses listed in the `allowed_ip` field will be denied regardless of any access policies (including
	// public access) that might otherwise permit the request.  Viewing or updating the `Firewall` element requires the
	// requester to have the `manager` role.
	Firewall *Firewall `json:"firewall,omitempty"`

	// Enables sending log data to IBM Cloud Activity Tracker Event Routing to provide visibility into bucket management,
	// object read and write events. (Recommended) When the `activity_tracker_crn` is not populated, then enabled events
	// are sent to the Activity Tracker Event Routing instance at the container's location unless otherwise specified in
	// the Activity Tracker Event Routing Event Routing service configuration. (Legacy) When the `activity_tracker_crn` is
	// populated, then enabled events are sent to the Activity Tracker Event Routing instance specified.
	ActivityTracking *ActivityTracking `json:"activity_tracking,omitempty"`

	// Enables sending metrics to IBM Cloud Monitoring.  All metrics are opt-in. (Recommended) When the
	// `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the Monitoring instance at the
	// container's location unless otherwise specified in the Metrics Router service configuration. (Legacy) When the
	// `metrics_monitoring_crn` is populated, then enabled metrics are sent to the Monitoring instance defined in the
	// `metrics_monitoring_crn` field.
	MetricsMonitoring *MetricsMonitoring `json:"metrics_monitoring,omitempty"`

	// Maximum bytes for this bucket.
	HardQuota *int64 `json:"hard_quota,omitempty"`

	// Data structure holding protection management response.
	ProtectionManagement *ProtectionManagementResponse `json:"protection_management,omitempty"`
}

// UnmarshalBucket unmarshals an instance of Bucket from the specified map of raw messages.
func UnmarshalBucket(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Bucket)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_id", &obj.ServiceInstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_crn", &obj.ServiceInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "time_created", &obj.TimeCreated)
	if err != nil {
		err = core.SDKErrorf(err, "", "time_created-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "time_updated", &obj.TimeUpdated)
	if err != nil {
		err = core.SDKErrorf(err, "", "time_updated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "object_count", &obj.ObjectCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bytes_used", &obj.BytesUsed)
	if err != nil {
		err = core.SDKErrorf(err, "", "bytes_used-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "noncurrent_object_count", &obj.NoncurrentObjectCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "noncurrent_object_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "noncurrent_bytes_used", &obj.NoncurrentBytesUsed)
	if err != nil {
		err = core.SDKErrorf(err, "", "noncurrent_bytes_used-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "delete_marker_count", &obj.DeleteMarkerCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "delete_marker_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "firewall", &obj.Firewall, UnmarshalFirewall)
	if err != nil {
		err = core.SDKErrorf(err, "", "firewall-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "activity_tracking", &obj.ActivityTracking, UnmarshalActivityTracking)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity_tracking-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metrics_monitoring", &obj.MetricsMonitoring, UnmarshalMetricsMonitoring)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics_monitoring-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hard_quota", &obj.HardQuota)
	if err != nil {
		err = core.SDKErrorf(err, "", "hard_quota-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "protection_management", &obj.ProtectionManagement, UnmarshalProtectionManagementResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "protection_management-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BucketPatch : An object containing new bucket metadata.
type BucketPatch struct {
	// An access control mechanism based on the network (IP address) where request originated. Requests not originating
	// from IP addresses listed in the `allowed_ip` field will be denied regardless of any access policies (including
	// public access) that might otherwise permit the request.  Viewing or updating the `Firewall` element requires the
	// requester to have the `manager` role.
	Firewall *Firewall `json:"firewall,omitempty"`

	// Enables sending log data to IBM Cloud Activity Tracker Event Routing to provide visibility into bucket management,
	// object read and write events. (Recommended) When the `activity_tracker_crn` is not populated, then enabled events
	// are sent to the Activity Tracker Event Routing instance at the container's location unless otherwise specified in
	// the Activity Tracker Event Routing Event Routing service configuration. (Legacy) When the `activity_tracker_crn` is
	// populated, then enabled events are sent to the Activity Tracker Event Routing instance specified.
	ActivityTracking *ActivityTracking `json:"activity_tracking,omitempty"`

	// Enables sending metrics to IBM Cloud Monitoring.  All metrics are opt-in. (Recommended) When the
	// `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the Monitoring instance at the
	// container's location unless otherwise specified in the Metrics Router service configuration. (Legacy) When the
	// `metrics_monitoring_crn` is populated, then enabled metrics are sent to the Monitoring instance defined in the
	// `metrics_monitoring_crn` field.
	MetricsMonitoring *MetricsMonitoring `json:"metrics_monitoring,omitempty"`

	// Maximum bytes for this bucket.
	HardQuota *int64 `json:"hard_quota,omitempty"`

	// Data structure holding protection management operations.
	ProtectionManagement *ProtectionManagement `json:"protection_management,omitempty"`
}

// UnmarshalBucketPatch unmarshals an instance of BucketPatch from the specified map of raw messages.
func UnmarshalBucketPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BucketPatch)
	err = core.UnmarshalModel(m, "firewall", &obj.Firewall, UnmarshalFirewall)
	if err != nil {
		err = core.SDKErrorf(err, "", "firewall-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "activity_tracking", &obj.ActivityTracking, UnmarshalActivityTracking)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity_tracking-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metrics_monitoring", &obj.MetricsMonitoring, UnmarshalMetricsMonitoring)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics_monitoring-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hard_quota", &obj.HardQuota)
	if err != nil {
		err = core.SDKErrorf(err, "", "hard_quota-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "protection_management", &obj.ProtectionManagement, UnmarshalProtectionManagement)
	if err != nil {
		err = core.SDKErrorf(err, "", "protection_management-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the BucketPatch
func (bucketPatch *BucketPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(bucketPatch.Firewall) {
		_patch["firewall"] = bucketPatch.Firewall.asPatch()
	}
	if !core.IsNil(bucketPatch.ActivityTracking) {
		_patch["activity_tracking"] = bucketPatch.ActivityTracking.asPatch()
	}
	if !core.IsNil(bucketPatch.MetricsMonitoring) {
		_patch["metrics_monitoring"] = bucketPatch.MetricsMonitoring.asPatch()
	}
	if !core.IsNil(bucketPatch.HardQuota) {
		_patch["hard_quota"] = bucketPatch.HardQuota
	}
	if !core.IsNil(bucketPatch.ProtectionManagement) {
		_patch["protection_management"] = bucketPatch.ProtectionManagement.asPatch()
	}

	return
}

// CreateBackupPolicyOptions : The CreateBackupPolicy options.
type CreateBackupPolicyOptions struct {
	// Name of the COS Bucket name.
	Bucket *string `json:"bucket" validate:"required,ne="`

	// The name granted to the policy. Validation :
	//   * chars limited to alphanumeric, underscore, hyphen and period.
	PolicyName *string `json:"policy_name" validate:"required"`

	// The CRN for a COS BackupVault.
	TargetBackupVaultCrn *string `json:"target_backup_vault_crn" validate:"required"`

	// The type of backup to support. For LA+GA this is limited to "continuous".
	BackupType *string `json:"backup_type" validate:"required"`

	// MD5 hash of content. If provided, the hash of the request must match.
	MD5 *string `json:"MD5,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateBackupPolicyOptions.BackupType property.
// The type of backup to support. For LA+GA this is limited to "continuous".
const (
	CreateBackupPolicyOptions_BackupType_Continuous = "continuous"
)

// NewCreateBackupPolicyOptions : Instantiate CreateBackupPolicyOptions
func (*ResourceConfigurationV1) NewCreateBackupPolicyOptions(bucket string, policyName string, targetBackupVaultCrn string, backupType string) *CreateBackupPolicyOptions {
	return &CreateBackupPolicyOptions{
		Bucket: core.StringPtr(bucket),
		PolicyName: core.StringPtr(policyName),
		TargetBackupVaultCrn: core.StringPtr(targetBackupVaultCrn),
		BackupType: core.StringPtr(backupType),
	}
}

// SetBucket : Allow user to set Bucket
func (_options *CreateBackupPolicyOptions) SetBucket(bucket string) *CreateBackupPolicyOptions {
	_options.Bucket = core.StringPtr(bucket)
	return _options
}

// SetPolicyName : Allow user to set PolicyName
func (_options *CreateBackupPolicyOptions) SetPolicyName(policyName string) *CreateBackupPolicyOptions {
	_options.PolicyName = core.StringPtr(policyName)
	return _options
}

// SetTargetBackupVaultCrn : Allow user to set TargetBackupVaultCrn
func (_options *CreateBackupPolicyOptions) SetTargetBackupVaultCrn(targetBackupVaultCrn string) *CreateBackupPolicyOptions {
	_options.TargetBackupVaultCrn = core.StringPtr(targetBackupVaultCrn)
	return _options
}

// SetBackupType : Allow user to set BackupType
func (_options *CreateBackupPolicyOptions) SetBackupType(backupType string) *CreateBackupPolicyOptions {
	_options.BackupType = core.StringPtr(backupType)
	return _options
}

// SetMD5 : Allow user to set MD5
func (_options *CreateBackupPolicyOptions) SetMD5(mD5 string) *CreateBackupPolicyOptions {
	_options.MD5 = core.StringPtr(mD5)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateBackupPolicyOptions) SetHeaders(param map[string]string) *CreateBackupPolicyOptions {
	options.Headers = param
	return options
}

// CreateBackupVaultOptions : The CreateBackupVault options.
type CreateBackupVaultOptions struct {
	// Name of the service_instance to list BackupVaults for.
	ServiceInstanceID *string `json:"service_instance_id" validate:"required"`

	// The name given to a Bucket.
	//
	// Bucket names must be between 3 and 63 characters long must be made of lowercase letters, numbers, dots (periods),
	// and dashes (hyphens). Bucket names must begin and end with a lowercase letter or number. Bucket names canâ€t contain
	// consecutive dots or dashes. Bucket names that resemble IP addresses are not allowed.
	//
	// Bucket and BackupVault names exist in a global global namespace and therefore must be unique.
	BackupVaultName *string `json:"backup_vault_name" validate:"required"`

	// the region in which this backup-vault should be created within.
	Region *string `json:"region" validate:"required"`

	// Activity Tracking configuration. An empty object (`{}`) indicates no configuration, and no events will be sent (This
	// is the same behavior as `{"management_events":false}`). Note that read/write events cannot be enabled, and events
	// cannot be routed to a non-default Activity Tracker instance.
	ActivityTracking *BackupVaultActivityTracking `json:"activity_tracking,omitempty"`

	// Metrics Monitoring configuration. An empty object (`{}`) indicates no configuration, and no metrics will be
	// collected (This is the same behavior as `{"usage_metrics_enabled":false}`). Note that request metrics cannot be
	// enabled, and metrics cannot be routed to a non-default metrics router instance.
	MetricsMonitoring *BackupVaultMetricsMonitoring `json:"metrics_monitoring,omitempty"`

	// The CRN for a KeyProtect root key.
	SseKpCustomerRootKeyCrn *string `json:"sse_kp_customer_root_key_crn,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateBackupVaultOptions : Instantiate CreateBackupVaultOptions
func (*ResourceConfigurationV1) NewCreateBackupVaultOptions(serviceInstanceID string, backupVaultName string, region string) *CreateBackupVaultOptions {
	return &CreateBackupVaultOptions{
		ServiceInstanceID: core.StringPtr(serviceInstanceID),
		BackupVaultName: core.StringPtr(backupVaultName),
		Region: core.StringPtr(region),
	}
}

// SetServiceInstanceID : Allow user to set ServiceInstanceID
func (_options *CreateBackupVaultOptions) SetServiceInstanceID(serviceInstanceID string) *CreateBackupVaultOptions {
	_options.ServiceInstanceID = core.StringPtr(serviceInstanceID)
	return _options
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *CreateBackupVaultOptions) SetBackupVaultName(backupVaultName string) *CreateBackupVaultOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *CreateBackupVaultOptions) SetRegion(region string) *CreateBackupVaultOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetActivityTracking : Allow user to set ActivityTracking
func (_options *CreateBackupVaultOptions) SetActivityTracking(activityTracking *BackupVaultActivityTracking) *CreateBackupVaultOptions {
	_options.ActivityTracking = activityTracking
	return _options
}

// SetMetricsMonitoring : Allow user to set MetricsMonitoring
func (_options *CreateBackupVaultOptions) SetMetricsMonitoring(metricsMonitoring *BackupVaultMetricsMonitoring) *CreateBackupVaultOptions {
	_options.MetricsMonitoring = metricsMonitoring
	return _options
}

// SetSseKpCustomerRootKeyCrn : Allow user to set SseKpCustomerRootKeyCrn
func (_options *CreateBackupVaultOptions) SetSseKpCustomerRootKeyCrn(sseKpCustomerRootKeyCrn string) *CreateBackupVaultOptions {
	_options.SseKpCustomerRootKeyCrn = core.StringPtr(sseKpCustomerRootKeyCrn)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateBackupVaultOptions) SetHeaders(param map[string]string) *CreateBackupVaultOptions {
	options.Headers = param
	return options
}

// CreateRestoreOptions : The CreateRestore options.
type CreateRestoreOptions struct {
	// name of BackupVault to restore from.
	BackupVaultName *string `json:"backup_vault_name" validate:"required,ne="`

	// A UUID that uniquely identifies a resource.
	RecoveryRangeID *string `json:"recovery_range_id" validate:"required"`

	// The type of restore to support. More options will be available in the future.
	RestoreType *string `json:"restore_type" validate:"required"`

	// Timestamp format used throughout the API.
	//
	// Accepts the following formats:
	//
	// YYYY-MM-DDTHH:mm:ssZ YYYY-MM-DDTHH:mm:ss YYYY-MM-DDTHH:mm:ss-hh:mm YYYY-MM-DDTHH:mm:ss+hh:mm
	// YYYY-MM-DDTHH:mm:ss.sssZ YYYY-MM-DDTHH:mm:ss.sss YYYY-MM-DDTHH:mm:ss.sss-hh:mm YYYY-MM-DDTHH:mm:ss.sss+hh:mm.
	RestorePointInTime *strfmt.DateTime `json:"restore_point_in_time" validate:"required"`

	// The CRN for a COS Bucket.
	//
	// Note that Softlayer CRNs do not contain dashes within the service_instance_id, whereas regular CRNs do. Although
	// bucket backup is not supported for softlayer accounts, this need not be enforced at the CRN parsing level.
	TargetResourceCrn *string `json:"target_resource_crn" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateRestoreOptions.RestoreType property.
// The type of restore to support. More options will be available in the future.
const (
	CreateRestoreOptions_RestoreType_InPlace = "in_place"
)

// NewCreateRestoreOptions : Instantiate CreateRestoreOptions
func (*ResourceConfigurationV1) NewCreateRestoreOptions(backupVaultName string, recoveryRangeID string, restoreType string, restorePointInTime *strfmt.DateTime, targetResourceCrn string) *CreateRestoreOptions {
	return &CreateRestoreOptions{
		BackupVaultName: core.StringPtr(backupVaultName),
		RecoveryRangeID: core.StringPtr(recoveryRangeID),
		RestoreType: core.StringPtr(restoreType),
		RestorePointInTime: restorePointInTime,
		TargetResourceCrn: core.StringPtr(targetResourceCrn),
	}
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *CreateRestoreOptions) SetBackupVaultName(backupVaultName string) *CreateRestoreOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetRecoveryRangeID : Allow user to set RecoveryRangeID
func (_options *CreateRestoreOptions) SetRecoveryRangeID(recoveryRangeID string) *CreateRestoreOptions {
	_options.RecoveryRangeID = core.StringPtr(recoveryRangeID)
	return _options
}

// SetRestoreType : Allow user to set RestoreType
func (_options *CreateRestoreOptions) SetRestoreType(restoreType string) *CreateRestoreOptions {
	_options.RestoreType = core.StringPtr(restoreType)
	return _options
}

// SetRestorePointInTime : Allow user to set RestorePointInTime
func (_options *CreateRestoreOptions) SetRestorePointInTime(restorePointInTime *strfmt.DateTime) *CreateRestoreOptions {
	_options.RestorePointInTime = restorePointInTime
	return _options
}

// SetTargetResourceCrn : Allow user to set TargetResourceCrn
func (_options *CreateRestoreOptions) SetTargetResourceCrn(targetResourceCrn string) *CreateRestoreOptions {
	_options.TargetResourceCrn = core.StringPtr(targetResourceCrn)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRestoreOptions) SetHeaders(param map[string]string) *CreateRestoreOptions {
	options.Headers = param
	return options
}

// DeleteBackupPolicyOptions : The DeleteBackupPolicy options.
type DeleteBackupPolicyOptions struct {
	// name of the bucket affected.
	Bucket *string `json:"bucket" validate:"required,ne="`

	// uuid of the BackupPolicy.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteBackupPolicyOptions : Instantiate DeleteBackupPolicyOptions
func (*ResourceConfigurationV1) NewDeleteBackupPolicyOptions(bucket string, policyID string) *DeleteBackupPolicyOptions {
	return &DeleteBackupPolicyOptions{
		Bucket: core.StringPtr(bucket),
		PolicyID: core.StringPtr(policyID),
	}
}

// SetBucket : Allow user to set Bucket
func (_options *DeleteBackupPolicyOptions) SetBucket(bucket string) *DeleteBackupPolicyOptions {
	_options.Bucket = core.StringPtr(bucket)
	return _options
}

// SetPolicyID : Allow user to set PolicyID
func (_options *DeleteBackupPolicyOptions) SetPolicyID(policyID string) *DeleteBackupPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteBackupPolicyOptions) SetHeaders(param map[string]string) *DeleteBackupPolicyOptions {
	options.Headers = param
	return options
}

// DeleteBackupVaultOptions : The DeleteBackupVault options.
type DeleteBackupVaultOptions struct {
	// Name of the backup-vault to create or update.
	BackupVaultName *string `json:"backup_vault_name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteBackupVaultOptions : Instantiate DeleteBackupVaultOptions
func (*ResourceConfigurationV1) NewDeleteBackupVaultOptions(backupVaultName string) *DeleteBackupVaultOptions {
	return &DeleteBackupVaultOptions{
		BackupVaultName: core.StringPtr(backupVaultName),
	}
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *DeleteBackupVaultOptions) SetBackupVaultName(backupVaultName string) *DeleteBackupVaultOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteBackupVaultOptions) SetHeaders(param map[string]string) *DeleteBackupVaultOptions {
	options.Headers = param
	return options
}

// Firewall : An access control mechanism based on the network (IP address) where request originated. Requests not originating from
// IP addresses listed in the `allowed_ip` field will be denied regardless of any access policies (including public
// access) that might otherwise permit the request.  Viewing or updating the `Firewall` element requires the requester
// to have the `manager` role.
type Firewall struct {
	// List of IPv4 or IPv6 addresses in CIDR notation to be affected by firewall in CIDR notation is supported. Passing an
	// empty array will lift the IP address filter.  The `allowed_ip` array can contain a maximum of 1000 items.
	AllowedIp []string `json:"allowed_ip"`
}

// Constants associated with the Firewall.AllowedNetworkType property.
// May contain `public`, `private`, and/or `direct` elements. Setting `allowed_network_type` to only `private` will
// prevent access to object storage from outside of the IBM Cloud.  The entire array will be overwritten in a `PATCH`
// operation. For more information on network types, [see the
// documentation](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-endpoints#advanced-endpoint-types).
const (
	Firewall_AllowedNetworkType_Direct = "direct"
	Firewall_AllowedNetworkType_Private = "private"
	Firewall_AllowedNetworkType_Public = "public"
)

// UnmarshalFirewall unmarshals an instance of Firewall from the specified map of raw messages.
func UnmarshalFirewall(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Firewall)
	err = core.UnmarshalPrimitive(m, "allowed_ip", &obj.AllowedIp)
	if err != nil {
		err = core.SDKErrorf(err, "", "allowed_ip-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the Firewall
func (firewall *Firewall) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(firewall.AllowedIp) {
		_patch["allowed_ip"] = firewall.AllowedIp
	}

	return
}

// GetBackupPolicyOptions : The GetBackupPolicy options.
type GetBackupPolicyOptions struct {
	// name of the bucket affected.
	Bucket *string `json:"bucket" validate:"required,ne="`

	// uuid of the BackupPolicy.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetBackupPolicyOptions : Instantiate GetBackupPolicyOptions
func (*ResourceConfigurationV1) NewGetBackupPolicyOptions(bucket string, policyID string) *GetBackupPolicyOptions {
	return &GetBackupPolicyOptions{
		Bucket: core.StringPtr(bucket),
		PolicyID: core.StringPtr(policyID),
	}
}

// SetBucket : Allow user to set Bucket
func (_options *GetBackupPolicyOptions) SetBucket(bucket string) *GetBackupPolicyOptions {
	_options.Bucket = core.StringPtr(bucket)
	return _options
}

// SetPolicyID : Allow user to set PolicyID
func (_options *GetBackupPolicyOptions) SetPolicyID(policyID string) *GetBackupPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBackupPolicyOptions) SetHeaders(param map[string]string) *GetBackupPolicyOptions {
	options.Headers = param
	return options
}

// GetBackupVaultOptions : The GetBackupVault options.
type GetBackupVaultOptions struct {
	// Name of the backup-vault to create or update.
	BackupVaultName *string `json:"backup_vault_name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetBackupVaultOptions : Instantiate GetBackupVaultOptions
func (*ResourceConfigurationV1) NewGetBackupVaultOptions(backupVaultName string) *GetBackupVaultOptions {
	return &GetBackupVaultOptions{
		BackupVaultName: core.StringPtr(backupVaultName),
	}
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *GetBackupVaultOptions) SetBackupVaultName(backupVaultName string) *GetBackupVaultOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBackupVaultOptions) SetHeaders(param map[string]string) *GetBackupVaultOptions {
	options.Headers = param
	return options
}

// GetBucketConfigOptions : The GetBucketConfig options.
type GetBucketConfigOptions struct {
	// Name of a bucket.
	Bucket *string `json:"bucket" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetBucketConfigOptions : Instantiate GetBucketConfigOptions
func (*ResourceConfigurationV1) NewGetBucketConfigOptions(bucket string) *GetBucketConfigOptions {
	return &GetBucketConfigOptions{
		Bucket: core.StringPtr(bucket),
	}
}

// SetBucket : Allow user to set Bucket
func (_options *GetBucketConfigOptions) SetBucket(bucket string) *GetBucketConfigOptions {
	_options.Bucket = core.StringPtr(bucket)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBucketConfigOptions) SetHeaders(param map[string]string) *GetBucketConfigOptions {
	options.Headers = param
	return options
}

// GetRestoreOptions : The GetRestore options.
type GetRestoreOptions struct {
	// name of BackupVault that the restore occured on.
	BackupVaultName *string `json:"backup_vault_name" validate:"required,ne="`

	// id of the restore to introspect on.
	RestoreID *string `json:"restore_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetRestoreOptions : Instantiate GetRestoreOptions
func (*ResourceConfigurationV1) NewGetRestoreOptions(backupVaultName string, restoreID string) *GetRestoreOptions {
	return &GetRestoreOptions{
		BackupVaultName: core.StringPtr(backupVaultName),
		RestoreID: core.StringPtr(restoreID),
	}
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *GetRestoreOptions) SetBackupVaultName(backupVaultName string) *GetRestoreOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetRestoreID : Allow user to set RestoreID
func (_options *GetRestoreOptions) SetRestoreID(restoreID string) *GetRestoreOptions {
	_options.RestoreID = core.StringPtr(restoreID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRestoreOptions) SetHeaders(param map[string]string) *GetRestoreOptions {
	options.Headers = param
	return options
}

// GetSourceResourceRecoveryRangeOptions : The GetSourceResourceRecoveryRange options.
type GetSourceResourceRecoveryRangeOptions struct {
	// name of BackupVault to update.
	BackupVaultName *string `json:"backup_vault_name" validate:"required,ne="`

	// ID of the RecoveryRange to update.
	RecoveryRangeID *string `json:"recovery_range_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetSourceResourceRecoveryRangeOptions : Instantiate GetSourceResourceRecoveryRangeOptions
func (*ResourceConfigurationV1) NewGetSourceResourceRecoveryRangeOptions(backupVaultName string, recoveryRangeID string) *GetSourceResourceRecoveryRangeOptions {
	return &GetSourceResourceRecoveryRangeOptions{
		BackupVaultName: core.StringPtr(backupVaultName),
		RecoveryRangeID: core.StringPtr(recoveryRangeID),
	}
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *GetSourceResourceRecoveryRangeOptions) SetBackupVaultName(backupVaultName string) *GetSourceResourceRecoveryRangeOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetRecoveryRangeID : Allow user to set RecoveryRangeID
func (_options *GetSourceResourceRecoveryRangeOptions) SetRecoveryRangeID(recoveryRangeID string) *GetSourceResourceRecoveryRangeOptions {
	_options.RecoveryRangeID = core.StringPtr(recoveryRangeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSourceResourceRecoveryRangeOptions) SetHeaders(param map[string]string) *GetSourceResourceRecoveryRangeOptions {
	options.Headers = param
	return options
}

// ListBackupPoliciesOptions : The ListBackupPolicies options.
type ListBackupPoliciesOptions struct {
	// Name of the COS Bucket name.
	Bucket *string `json:"bucket" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListBackupPoliciesOptions : Instantiate ListBackupPoliciesOptions
func (*ResourceConfigurationV1) NewListBackupPoliciesOptions(bucket string) *ListBackupPoliciesOptions {
	return &ListBackupPoliciesOptions{
		Bucket: core.StringPtr(bucket),
	}
}

// SetBucket : Allow user to set Bucket
func (_options *ListBackupPoliciesOptions) SetBucket(bucket string) *ListBackupPoliciesOptions {
	_options.Bucket = core.StringPtr(bucket)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListBackupPoliciesOptions) SetHeaders(param map[string]string) *ListBackupPoliciesOptions {
	options.Headers = param
	return options
}

// ListBackupVaultsOptions : The ListBackupVaults options.
type ListBackupVaultsOptions struct {
	// Name of the service_instance to list BackupVaults for.
	ServiceInstanceID *string `json:"service_instance_id" validate:"required"`

	// the continuation token for controlling pagination.
	Token *string `json:"token,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListBackupVaultsOptions : Instantiate ListBackupVaultsOptions
func (*ResourceConfigurationV1) NewListBackupVaultsOptions(serviceInstanceID string) *ListBackupVaultsOptions {
	return &ListBackupVaultsOptions{
		ServiceInstanceID: core.StringPtr(serviceInstanceID),
	}
}

// SetServiceInstanceID : Allow user to set ServiceInstanceID
func (_options *ListBackupVaultsOptions) SetServiceInstanceID(serviceInstanceID string) *ListBackupVaultsOptions {
	_options.ServiceInstanceID = core.StringPtr(serviceInstanceID)
	return _options
}

// SetToken : Allow user to set Token
func (_options *ListBackupVaultsOptions) SetToken(token string) *ListBackupVaultsOptions {
	_options.Token = core.StringPtr(token)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListBackupVaultsOptions) SetHeaders(param map[string]string) *ListBackupVaultsOptions {
	options.Headers = param
	return options
}

// ListRecoveryRangesOptions : The ListRecoveryRanges options.
type ListRecoveryRangesOptions struct {
	// name of BackupVault.
	BackupVaultName *string `json:"backup_vault_name" validate:"required,ne="`

	// CRN of source resource to filter on. This limits ranges returned to only ranges where the source_resource_crn
	// matches the parameter value.
	SourceResourceCrn *string `json:"source_resource_crn,omitempty"`

	// the continuation token for controlling pagination.
	Token *string `json:"token,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListRecoveryRangesOptions : Instantiate ListRecoveryRangesOptions
func (*ResourceConfigurationV1) NewListRecoveryRangesOptions(backupVaultName string) *ListRecoveryRangesOptions {
	return &ListRecoveryRangesOptions{
		BackupVaultName: core.StringPtr(backupVaultName),
	}
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *ListRecoveryRangesOptions) SetBackupVaultName(backupVaultName string) *ListRecoveryRangesOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetSourceResourceCrn : Allow user to set SourceResourceCrn
func (_options *ListRecoveryRangesOptions) SetSourceResourceCrn(sourceResourceCrn string) *ListRecoveryRangesOptions {
	_options.SourceResourceCrn = core.StringPtr(sourceResourceCrn)
	return _options
}

// SetToken : Allow user to set Token
func (_options *ListRecoveryRangesOptions) SetToken(token string) *ListRecoveryRangesOptions {
	_options.Token = core.StringPtr(token)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListRecoveryRangesOptions) SetHeaders(param map[string]string) *ListRecoveryRangesOptions {
	options.Headers = param
	return options
}

// ListRestoresOptions : The ListRestores options.
type ListRestoresOptions struct {
	// name of BackupVault to restore from.
	BackupVaultName *string `json:"backup_vault_name" validate:"required,ne="`

	// the continuation token for controlling pagination.
	Token *string `json:"token,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListRestoresOptions : Instantiate ListRestoresOptions
func (*ResourceConfigurationV1) NewListRestoresOptions(backupVaultName string) *ListRestoresOptions {
	return &ListRestoresOptions{
		BackupVaultName: core.StringPtr(backupVaultName),
	}
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *ListRestoresOptions) SetBackupVaultName(backupVaultName string) *ListRestoresOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetToken : Allow user to set Token
func (_options *ListRestoresOptions) SetToken(token string) *ListRestoresOptions {
	_options.Token = core.StringPtr(token)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListRestoresOptions) SetHeaders(param map[string]string) *ListRestoresOptions {
	options.Headers = param
	return options
}

// MetricsMonitoring : Enables sending metrics to IBM Cloud Monitoring.  All metrics are opt-in. (Recommended) When the
// `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the Monitoring instance at the
// container's location unless otherwise specified in the Metrics Router service configuration. (Legacy) When the
// `metrics_monitoring_crn` is populated, then enabled metrics are sent to the Monitoring instance defined in the
// `metrics_monitoring_crn` field.
type MetricsMonitoring struct {
	// If set to `true`, all usage metrics (i.e. `bytes_used`) will be sent to the monitoring service.
	UsageMetricsEnabled *bool `json:"usage_metrics_enabled,omitempty"`

	// If set to `true`, all request metrics (i.e. `rest.object.head`) will be sent to the monitoring service.
	RequestMetricsEnabled *bool `json:"request_metrics_enabled,omitempty"`

	// When the `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the monitoring instance
	// associated to the container's location unless otherwise specified in the Metrics Router service configuration. If
	// `metrics_monitoring_crn` is populated, then enabled events are sent to the Metrics Monitoring instance specified.
	MetricsMonitoringCrn *string `json:"metrics_monitoring_crn,omitempty"`
}

// UnmarshalMetricsMonitoring unmarshals an instance of MetricsMonitoring from the specified map of raw messages.
func UnmarshalMetricsMonitoring(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MetricsMonitoring)
	err = core.UnmarshalPrimitive(m, "usage_metrics_enabled", &obj.UsageMetricsEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "usage_metrics_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "request_metrics_enabled", &obj.RequestMetricsEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "request_metrics_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metrics_monitoring_crn", &obj.MetricsMonitoringCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics_monitoring_crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the MetricsMonitoring
func (metricsMonitoring *MetricsMonitoring) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(metricsMonitoring.UsageMetricsEnabled) {
		_patch["usage_metrics_enabled"] = metricsMonitoring.UsageMetricsEnabled
	}
	if !core.IsNil(metricsMonitoring.RequestMetricsEnabled) {
		_patch["request_metrics_enabled"] = metricsMonitoring.RequestMetricsEnabled
	}
	if !core.IsNil(metricsMonitoring.MetricsMonitoringCrn) {
		_patch["metrics_monitoring_crn"] = metricsMonitoring.MetricsMonitoringCrn
	}

	return
}

// NextPagination : Pagination response body.
type NextPagination struct {
	// A URL to the continuation of results.
	Href *string `json:"href" validate:"required"`

	// The continutation token utilized for paginated results.
	Token *string `json:"token" validate:"required"`
}

// UnmarshalNextPagination unmarshals an instance of NextPagination from the specified map of raw messages.
func UnmarshalNextPagination(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NextPagination)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "token", &obj.Token)
	if err != nil {
		err = core.SDKErrorf(err, "", "token-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProtectionManagement : Data structure holding protection management operations.
type ProtectionManagement struct {
	// If set to `activate`, protection management action on the bucket is being activated.
	RequestedState *string `json:"requested_state,omitempty"`

	// This field is required when using requested_state\:`activate` and holds a JWT that is provided by the Cloud
	// Operator. This should be the encoded JWT.
	ProtectionManagementToken *string `json:"protection_management_token,omitempty"`
}

// Constants associated with the ProtectionManagement.RequestedState property.
// If set to `activate`, protection management action on the bucket is being activated.
const (
	ProtectionManagement_RequestedState_Activate = "activate"
	ProtectionManagement_RequestedState_Deactivate = "deactivate"
)

// UnmarshalProtectionManagement unmarshals an instance of ProtectionManagement from the specified map of raw messages.
func UnmarshalProtectionManagement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProtectionManagement)
	err = core.UnmarshalPrimitive(m, "requested_state", &obj.RequestedState)
	if err != nil {
		err = core.SDKErrorf(err, "", "requested_state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "protection_management_token", &obj.ProtectionManagementToken)
	if err != nil {
		err = core.SDKErrorf(err, "", "protection_management_token-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the ProtectionManagement
func (protectionManagement *ProtectionManagement) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(protectionManagement.RequestedState) {
		_patch["requested_state"] = protectionManagement.RequestedState
	}
	if !core.IsNil(protectionManagement.ProtectionManagementToken) {
		_patch["protection_management_token"] = protectionManagement.ProtectionManagementToken
	}

	return
}

// ProtectionManagementResponse : Data structure holding protection management response.
type ProtectionManagementResponse struct {
	// Indicates the X number of protection management tokens that have been applied to the bucket in its lifetime.
	TokenAppliedCounter *string `json:"token_applied_counter,omitempty"`

	// The 'protection management token list' holding a recent list of applied tokens. This list may contain a subset of
	// all tokens applied to the bucket, as indicated by the counter.
	TokenEntries []ProtectionManagementResponseTokenEntry `json:"token_entries,omitempty"`
}

// UnmarshalProtectionManagementResponse unmarshals an instance of ProtectionManagementResponse from the specified map of raw messages.
func UnmarshalProtectionManagementResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProtectionManagementResponse)
	err = core.UnmarshalPrimitive(m, "token_applied_counter", &obj.TokenAppliedCounter)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_applied_counter-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "token_entries", &obj.TokenEntries, UnmarshalProtectionManagementResponseTokenEntry)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_entries-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProtectionManagementResponseTokenEntry : Data structure holding protection management token.
type ProtectionManagementResponseTokenEntry struct {
	TokenID *string `json:"token_id,omitempty"`

	TokenExpirationTime *string `json:"token_expiration_time,omitempty"`

	TokenReferenceID *string `json:"token_reference_id,omitempty"`

	AppliedTime *string `json:"applied_time,omitempty"`

	InvalidatedTime *string `json:"invalidated_time,omitempty"`

	ExpirationTime *string `json:"expiration_time,omitempty"`

	ShortenRetentionFlag *bool `json:"shorten_retention_flag,omitempty"`
}

// UnmarshalProtectionManagementResponseTokenEntry unmarshals an instance of ProtectionManagementResponseTokenEntry from the specified map of raw messages.
func UnmarshalProtectionManagementResponseTokenEntry(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProtectionManagementResponseTokenEntry)
	err = core.UnmarshalPrimitive(m, "token_id", &obj.TokenID)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "token_expiration_time", &obj.TokenExpirationTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_expiration_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "token_reference_id", &obj.TokenReferenceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_reference_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "applied_time", &obj.AppliedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "applied_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "invalidated_time", &obj.InvalidatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "invalidated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_time", &obj.ExpirationTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "expiration_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "shorten_retention_flag", &obj.ShortenRetentionFlag)
	if err != nil {
		err = core.SDKErrorf(err, "", "shorten_retention_flag-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RecoveryRange : Metadata associated with a recovery range.
type RecoveryRange struct {
	// The CRN of the sourceResource backed up by this RecoveryRange.
	SourceResourceCrn *string `json:"source_resource_crn,omitempty"`

	// The name of the backupPolicy that triggered the creation of this RecoveryRange.
	BackupPolicyName *string `json:"backup_policy_name,omitempty"`

	// The point in time at which backup coverage of the sourceResource begins.
	//
	// Returns "YYYY-MM-DDTHH:mm:ss.sssZ" timestamp format.
	RangeStartTime *strfmt.DateTime `json:"range_start_time,omitempty"`

	// the point in time at which backup coverage of the sourceResource ends. Returns "YYYY-MM-DDTHH:mm:ss.sssZ" timestamp
	// format.
	RangeEndTime *strfmt.DateTime `json:"range_end_time,omitempty"`

	// The time at which this recoveryRange was initially created.
	//
	// Returns "YYYY-MM-DDTHH:mm:ss.sssZ" timestamp format
	//
	// NOTE : this can be before the start-time.
	RangeCreateTime *strfmt.DateTime `json:"range_create_time,omitempty"`

	// A UUID that uniquely identifies a resource.
	RecoveryRangeID *string `json:"recovery_range_id,omitempty"`
}

// UnmarshalRecoveryRange unmarshals an instance of RecoveryRange from the specified map of raw messages.
func UnmarshalRecoveryRange(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RecoveryRange)
	err = core.UnmarshalPrimitive(m, "source_resource_crn", &obj.SourceResourceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_resource_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "backup_policy_name", &obj.BackupPolicyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "backup_policy_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "range_start_time", &obj.RangeStartTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "range_start_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "range_end_time", &obj.RangeEndTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "range_end_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "range_create_time", &obj.RangeCreateTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "range_create_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "recovery_range_id", &obj.RecoveryRangeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "recovery_range_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RecoveryRangeCollection : A collection of recovery ranges.
type RecoveryRangeCollection struct {
	// Pagination response body.
	Next *NextPagination `json:"next,omitempty"`

	// A list of recovery ranges.
	RecoveryRanges []RecoveryRange `json:"recovery_ranges,omitempty"`
}

// UnmarshalRecoveryRangeCollection unmarshals an instance of RecoveryRangeCollection from the specified map of raw messages.
func UnmarshalRecoveryRangeCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RecoveryRangeCollection)
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextPagination)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "recovery_ranges", &obj.RecoveryRanges, UnmarshalRecoveryRange)
	if err != nil {
		err = core.SDKErrorf(err, "", "recovery_ranges-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *RecoveryRangeCollection) GetNextToken() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Token, nil
}

// Restore : Metadata associated with a requested restore operation.
type Restore struct {
	// A UUID that uniquely identifies a resource.
	RecoveryRangeID *string `json:"recovery_range_id" validate:"required"`

	// The type of restore to support. More options will be available in the future.
	RestoreType *string `json:"restore_type" validate:"required"`

	// Timestamp format used throughout the API.
	//
	// Accepts the following formats:
	//
	// YYYY-MM-DDTHH:mm:ssZ YYYY-MM-DDTHH:mm:ss YYYY-MM-DDTHH:mm:ss-hh:mm YYYY-MM-DDTHH:mm:ss+hh:mm
	// YYYY-MM-DDTHH:mm:ss.sssZ YYYY-MM-DDTHH:mm:ss.sss YYYY-MM-DDTHH:mm:ss.sss-hh:mm YYYY-MM-DDTHH:mm:ss.sss+hh:mm.
	RestorePointInTime *strfmt.DateTime `json:"restore_point_in_time" validate:"required"`

	// The CRN for a COS Bucket.
	//
	// Note that Softlayer CRNs do not contain dashes within the service_instance_id, whereas regular CRNs do. Although
	// bucket backup is not supported for softlayer accounts, this need not be enforced at the CRN parsing level.
	TargetResourceCrn *string `json:"target_resource_crn" validate:"required"`

	// The name of the source resource that is being restored from.
	SourceResourceCrn *string `json:"source_resource_crn,omitempty"`

	// Unique system-defined UUID for this restore operation.
	RestoreID *string `json:"restore_id,omitempty"`

	// The current status for this restore operation.
	//
	// initializing: The operation is initializing. Do not expect to see restored objects on the target bucket.  running :
	// The operation is ongoing. Expect to see some restored objects on the target bucket.  complete: The operation has
	// completed successfully.  failed: The operation has completed unsuccessfully.
	RestoreStatus *string `json:"restore_status,omitempty"`

	// The time at which this restore was initiated Returns "YYYY-MM-DDTHH:mm:ss.sssZ" timestamp format.
	InitTime *strfmt.DateTime `json:"init_time,omitempty"`

	// The time at which this restore ended (in both success and error cases) Returns "YYYY-MM-DDTHH:mm:ss.sssZ" timestamp
	// format.
	CompleteTime *strfmt.DateTime `json:"complete_time,omitempty"`

	// reports percent-doneness of init. Only present when restore_status=running.
	RestorePercentProgress *int64 `json:"restore_percent_progress,omitempty"`

	// Only present when restore_status=running.
	ErrorCause *string `json:"error_cause,omitempty"`
}

// Constants associated with the Restore.RestoreType property.
// The type of restore to support. More options will be available in the future.
const (
	Restore_RestoreType_InPlace = "in_place"
)

// Constants associated with the Restore.RestoreStatus property.
// The current status for this restore operation.
//
// initializing: The operation is initializing. Do not expect to see restored objects on the target bucket.  running :
// The operation is ongoing. Expect to see some restored objects on the target bucket.  complete: The operation has
// completed successfully.  failed: The operation has completed unsuccessfully.
const (
	Restore_RestoreStatus_Complete = "complete"
	Restore_RestoreStatus_Failed = "failed"
	Restore_RestoreStatus_Initializing = "initializing"
	Restore_RestoreStatus_Running = "running"
)

// UnmarshalRestore unmarshals an instance of Restore from the specified map of raw messages.
func UnmarshalRestore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Restore)
	err = core.UnmarshalPrimitive(m, "recovery_range_id", &obj.RecoveryRangeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "recovery_range_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restore_type", &obj.RestoreType)
	if err != nil {
		err = core.SDKErrorf(err, "", "restore_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restore_point_in_time", &obj.RestorePointInTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "restore_point_in_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_resource_crn", &obj.TargetResourceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_resource_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_resource_crn", &obj.SourceResourceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_resource_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restore_id", &obj.RestoreID)
	if err != nil {
		err = core.SDKErrorf(err, "", "restore_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restore_status", &obj.RestoreStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "restore_status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "init_time", &obj.InitTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "init_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "complete_time", &obj.CompleteTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "complete_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restore_percent_progress", &obj.RestorePercentProgress)
	if err != nil {
		err = core.SDKErrorf(err, "", "restore_percent_progress-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "error_cause", &obj.ErrorCause)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_cause-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RestoreCollection : A list of restore operations.
type RestoreCollection struct {
	// Pagination response body.
	Next *NextPagination `json:"next,omitempty"`

	// A collection of active and completed restore operations.
	Restores []Restore `json:"restores,omitempty"`
}

// UnmarshalRestoreCollection unmarshals an instance of RestoreCollection from the specified map of raw messages.
func UnmarshalRestoreCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RestoreCollection)
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNextPagination)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "restores", &obj.Restores, UnmarshalRestore)
	if err != nil {
		err = core.SDKErrorf(err, "", "restores-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *RestoreCollection) GetNextToken() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Token, nil
}

// UpdateBackupVaultOptions : The UpdateBackupVault options.
type UpdateBackupVaultOptions struct {
	// Name of the backup-vault to create or update.
	BackupVaultName *string `json:"backup_vault_name" validate:"required,ne="`

	// A Backup Vault config object containing changes to apply to the existing Backup Vault config.
	BackupVaultPatch map[string]interface{} `json:"BackupVault_patch" validate:"required"`

	// Conditionally update the Backup Vault config if and only if the ETag of the existing config exactly matches the
	// provided If-Match MD5.
	IfMatch *string `json:"If-Match,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateBackupVaultOptions : Instantiate UpdateBackupVaultOptions
func (*ResourceConfigurationV1) NewUpdateBackupVaultOptions(backupVaultName string, backupVaultPatch map[string]interface{}) *UpdateBackupVaultOptions {
	return &UpdateBackupVaultOptions{
		BackupVaultName: core.StringPtr(backupVaultName),
		BackupVaultPatch: backupVaultPatch,
	}
}

// SetBackupVaultName : Allow user to set BackupVaultName
func (_options *UpdateBackupVaultOptions) SetBackupVaultName(backupVaultName string) *UpdateBackupVaultOptions {
	_options.BackupVaultName = core.StringPtr(backupVaultName)
	return _options
}

// SetBackupVaultPatch : Allow user to set BackupVaultPatch
func (_options *UpdateBackupVaultOptions) SetBackupVaultPatch(backupVaultPatch map[string]interface{}) *UpdateBackupVaultOptions {
	_options.BackupVaultPatch = backupVaultPatch
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateBackupVaultOptions) SetIfMatch(ifMatch string) *UpdateBackupVaultOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateBackupVaultOptions) SetHeaders(param map[string]string) *UpdateBackupVaultOptions {
	options.Headers = param
	return options
}

// UpdateBucketConfigOptions : The UpdateBucketConfig options.
type UpdateBucketConfigOptions struct {
	// Name of a bucket.
	Bucket *string `json:"bucket" validate:"required,ne="`

	// An object containing new configuration metadata.
	BucketPatch map[string]interface{} `json:"Bucket_patch,omitempty"`

	// An Etag previously returned in a header when fetching or updating a bucket's metadata. If this value does not match
	// the active Etag, the request will fail.
	IfMatch *string `json:"If-Match,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateBucketConfigOptions : Instantiate UpdateBucketConfigOptions
func (*ResourceConfigurationV1) NewUpdateBucketConfigOptions(bucket string) *UpdateBucketConfigOptions {
	return &UpdateBucketConfigOptions{
		Bucket: core.StringPtr(bucket),
	}
}

// SetBucket : Allow user to set Bucket
func (_options *UpdateBucketConfigOptions) SetBucket(bucket string) *UpdateBucketConfigOptions {
	_options.Bucket = core.StringPtr(bucket)
	return _options
}

// SetBucketPatch : Allow user to set BucketPatch
func (_options *UpdateBucketConfigOptions) SetBucketPatch(bucketPatch map[string]interface{}) *UpdateBucketConfigOptions {
	_options.BucketPatch = bucketPatch
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateBucketConfigOptions) SetIfMatch(ifMatch string) *UpdateBucketConfigOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateBucketConfigOptions) SetHeaders(param map[string]string) *UpdateBucketConfigOptions {
	options.Headers = param
	return options
}

//
// BackupVaultsPager can be used to simplify the use of the "ListBackupVaults" method.
//
type BackupVaultsPager struct {
	hasNext bool
	options *ListBackupVaultsOptions
	client  *ResourceConfigurationV1
	pageContext struct {
		next *string
	}
}

// NewBackupVaultsPager returns a new BackupVaultsPager instance.
func (resourceConfiguration *ResourceConfigurationV1) NewBackupVaultsPager(options *ListBackupVaultsOptions) (pager *BackupVaultsPager, err error) {
	if options.Token != nil && *options.Token != "" {
		err = core.SDKErrorf(nil, "the 'options.Token' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListBackupVaultsOptions = *options
	pager = &BackupVaultsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  resourceConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *BackupVaultsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *BackupVaultsPager) GetNextWithContext(ctx context.Context) (page []string, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Token = pager.pageContext.next

	result, _, err := pager.client.ListBackupVaultsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Token
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.BackupVaults

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *BackupVaultsPager) GetAllWithContext(ctx context.Context) (allItems []string, err error) {
	for pager.HasNext() {
		var nextPage []string
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			err = core.RepurposeSDKProblem(err, "error-getting-next-page")
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *BackupVaultsPager) GetNext() (page []string, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *BackupVaultsPager) GetAll() (allItems []string, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// RecoveryRangesPager can be used to simplify the use of the "ListRecoveryRanges" method.
//
type RecoveryRangesPager struct {
	hasNext bool
	options *ListRecoveryRangesOptions
	client  *ResourceConfigurationV1
	pageContext struct {
		next *string
	}
}

// NewRecoveryRangesPager returns a new RecoveryRangesPager instance.
func (resourceConfiguration *ResourceConfigurationV1) NewRecoveryRangesPager(options *ListRecoveryRangesOptions) (pager *RecoveryRangesPager, err error) {
	if options.Token != nil && *options.Token != "" {
		err = core.SDKErrorf(nil, "the 'options.Token' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListRecoveryRangesOptions = *options
	pager = &RecoveryRangesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  resourceConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *RecoveryRangesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *RecoveryRangesPager) GetNextWithContext(ctx context.Context) (page []RecoveryRange, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Token = pager.pageContext.next

	result, _, err := pager.client.ListRecoveryRangesWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Token
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.RecoveryRanges

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *RecoveryRangesPager) GetAllWithContext(ctx context.Context) (allItems []RecoveryRange, err error) {
	for pager.HasNext() {
		var nextPage []RecoveryRange
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			err = core.RepurposeSDKProblem(err, "error-getting-next-page")
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *RecoveryRangesPager) GetNext() (page []RecoveryRange, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *RecoveryRangesPager) GetAll() (allItems []RecoveryRange, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// RestoresPager can be used to simplify the use of the "ListRestores" method.
//
type RestoresPager struct {
	hasNext bool
	options *ListRestoresOptions
	client  *ResourceConfigurationV1
	pageContext struct {
		next *string
	}
}

// NewRestoresPager returns a new RestoresPager instance.
func (resourceConfiguration *ResourceConfigurationV1) NewRestoresPager(options *ListRestoresOptions) (pager *RestoresPager, err error) {
	if options.Token != nil && *options.Token != "" {
		err = core.SDKErrorf(nil, "the 'options.Token' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListRestoresOptions = *options
	pager = &RestoresPager{
		hasNext: true,
		options: &optionsCopy,
		client:  resourceConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *RestoresPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *RestoresPager) GetNextWithContext(ctx context.Context) (page []Restore, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Token = pager.pageContext.next

	result, _, err := pager.client.ListRestoresWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Token
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Restores

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *RestoresPager) GetAllWithContext(ctx context.Context) (allItems []Restore, err error) {
	for pager.HasNext() {
		var nextPage []Restore
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			err = core.RepurposeSDKProblem(err, "error-getting-next-page")
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *RestoresPager) GetNext() (page []Restore, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *RestoresPager) GetAll() (allItems []Restore, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
