/**
 * (C) Copyright IBM Corp. 2024-2025.
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
 * IBM OpenAPI SDK Code Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

// Package sdsaasv1 : Operations and models for the SdsaasV1 service
package sdsaasv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/sds-go-sdk/common"
)

// SdsaasV1 : OpenAPI definition for SDSaaS
//
// API Version: 1.0.0
type SdsaasV1 struct {
	Service *core.BaseService
}

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "sdsaas"

const ParameterizedServiceURL = "{url}"

var defaultUrlVariables = map[string]string{
	"url": "{url}",
}

// SdsaasV1Options : Service options
type SdsaasV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewSdsaasV1UsingExternalConfig : constructs an instance of SdsaasV1 with passed in options and external configuration.
func NewSdsaasV1UsingExternalConfig(options *SdsaasV1Options) (sdsaas *SdsaasV1, err error) {
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

	sdsaas, err = NewSdsaasV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = sdsaas.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = sdsaas.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewSdsaasV1 : constructs an instance of SdsaasV1 with passed in options.
func NewSdsaasV1(options *SdsaasV1Options) (service *SdsaasV1, err error) {
	serviceOptions := &core.ServiceOptions{
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

	service = &SdsaasV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "sdsaas" suitable for processing requests.
func (sdsaas *SdsaasV1) Clone() *SdsaasV1 {
	if core.IsNil(sdsaas) {
		return nil
	}
	clone := *sdsaas
	clone.Service = sdsaas.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (sdsaas *SdsaasV1) SetServiceURL(url string) error {
	err := sdsaas.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (sdsaas *SdsaasV1) GetServiceURL() string {
	return sdsaas.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (sdsaas *SdsaasV1) SetDefaultHeaders(headers http.Header) {
	sdsaas.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (sdsaas *SdsaasV1) SetEnableGzipCompression(enableGzip bool) {
	sdsaas.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (sdsaas *SdsaasV1) GetEnableGzipCompression() bool {
	return sdsaas.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (sdsaas *SdsaasV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	sdsaas.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (sdsaas *SdsaasV1) DisableRetries() {
	sdsaas.Service.DisableRetries()
}

// Volumes : This request lists all volumes in the region
// Volumes are network-connected block storage devices that may be attached to one or more instances in the same region.
func (sdsaas *SdsaasV1) Volumes(volumesOptions *VolumesOptions) (result *VolumeCollection, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.VolumesWithContext(context.Background(), volumesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// VolumesWithContext is an alternate form of the Volumes method which supports a Context parameter
func (sdsaas *SdsaasV1) VolumesWithContext(ctx context.Context, volumesOptions *VolumesOptions) (result *VolumeCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(volumesOptions, "volumesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/volumes`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range volumesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "Volumes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if volumesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*volumesOptions.Limit))
	}
	if volumesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*volumesOptions.Name))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "volumes", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVolumeCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// VolumeCreate : Create a new volume
// Create a volume.
func (sdsaas *SdsaasV1) VolumeCreate(volumeCreateOptions *VolumeCreateOptions) (result *Volume, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.VolumeCreateWithContext(context.Background(), volumeCreateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// VolumeCreateWithContext is an alternate form of the VolumeCreate method which supports a Context parameter
func (sdsaas *SdsaasV1) VolumeCreateWithContext(ctx context.Context, volumeCreateOptions *VolumeCreateOptions) (result *Volume, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(volumeCreateOptions, "volumeCreateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(volumeCreateOptions, "volumeCreateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/volumes`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range volumeCreateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "VolumeCreate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if volumeCreateOptions.Hostnqnstring != nil {
		builder.AddQuery("hostnqnstring", fmt.Sprint(*volumeCreateOptions.Hostnqnstring))
	}

	body := make(map[string]interface{})
	if volumeCreateOptions.Capacity != nil {
		body["capacity"] = volumeCreateOptions.Capacity
	}
	if volumeCreateOptions.Name != nil {
		body["name"] = volumeCreateOptions.Name
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
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "volume_create", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVolume)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// Volume : Retrieve a volume profile
// This request retrieves a single volume profile specified by ID.
func (sdsaas *SdsaasV1) Volume(volumeOptions *VolumeOptions) (result *Volume, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.VolumeWithContext(context.Background(), volumeOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// VolumeWithContext is an alternate form of the Volume method which supports a Context parameter
func (sdsaas *SdsaasV1) VolumeWithContext(ctx context.Context, volumeOptions *VolumeOptions) (result *Volume, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(volumeOptions, "volumeOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(volumeOptions, "volumeOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"volume_id": *volumeOptions.VolumeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/volumes/{volume_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range volumeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "Volume")
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
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "volume", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVolume)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// VolumeDelete : Delete a volume
// This request deletes a single volume profile based on the name.
func (sdsaas *SdsaasV1) VolumeDelete(volumeDeleteOptions *VolumeDeleteOptions) (response *core.DetailedResponse, err error) {
	response, err = sdsaas.VolumeDeleteWithContext(context.Background(), volumeDeleteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// VolumeDeleteWithContext is an alternate form of the VolumeDelete method which supports a Context parameter
func (sdsaas *SdsaasV1) VolumeDeleteWithContext(ctx context.Context, volumeDeleteOptions *VolumeDeleteOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(volumeDeleteOptions, "volumeDeleteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(volumeDeleteOptions, "volumeDeleteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"volume_id": *volumeDeleteOptions.VolumeID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/volumes/{volume_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range volumeDeleteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "VolumeDelete")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = sdsaas.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "volume_delete", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// VolumeUpdate : Update a volume
// This request updates a volume with the information in a provided volume patch.
func (sdsaas *SdsaasV1) VolumeUpdate(volumeUpdateOptions *VolumeUpdateOptions) (result *Volume, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.VolumeUpdateWithContext(context.Background(), volumeUpdateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// VolumeUpdateWithContext is an alternate form of the VolumeUpdate method which supports a Context parameter
func (sdsaas *SdsaasV1) VolumeUpdateWithContext(ctx context.Context, volumeUpdateOptions *VolumeUpdateOptions) (result *Volume, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(volumeUpdateOptions, "volumeUpdateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(volumeUpdateOptions, "volumeUpdateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"volume_id": *volumeUpdateOptions.VolumeID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/volumes/{volume_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range volumeUpdateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "VolumeUpdate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if volumeUpdateOptions.VolumePatch != nil {
		_, err = builder.SetBodyContentJSON(volumeUpdateOptions.VolumePatch)
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

	var rawResponse map[string]json.RawMessage
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "volume_update", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVolume)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// Creds : List storage account credentials
// This request retrieves credentials for a specific storage account.
func (sdsaas *SdsaasV1) Creds(credsOptions *CredsOptions) (result *CredentialsFound, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.CredsWithContext(context.Background(), credsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CredsWithContext is an alternate form of the Creds method which supports a Context parameter
func (sdsaas *SdsaasV1) CredsWithContext(ctx context.Context, credsOptions *CredsOptions) (result *CredentialsFound, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(credsOptions, "credsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/v1/object/workspace/credentials`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range credsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "Creds")
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
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "creds", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCredentialsFound)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CredCreate : Create or modify storage account credentials
// Updates credentials for a storage account or creates them if they do not exist.
func (sdsaas *SdsaasV1) CredCreate(credCreateOptions *CredCreateOptions) (result *CredentialsUpdated, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.CredCreateWithContext(context.Background(), credCreateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CredCreateWithContext is an alternate form of the CredCreate method which supports a Context parameter
func (sdsaas *SdsaasV1) CredCreateWithContext(ctx context.Context, credCreateOptions *CredCreateOptions) (result *CredentialsUpdated, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(credCreateOptions, "credCreateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(credCreateOptions, "credCreateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/v1/object/workspace/credentials`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range credCreateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "CredCreate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("access_key", fmt.Sprint(*credCreateOptions.AccessKey))

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "cred_create", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCredentialsUpdated)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CredDelete : Delete storage account credentials
// Deletes specific credentials for a storage account.
func (sdsaas *SdsaasV1) CredDelete(credDeleteOptions *CredDeleteOptions) (response *core.DetailedResponse, err error) {
	response, err = sdsaas.CredDeleteWithContext(context.Background(), credDeleteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CredDeleteWithContext is an alternate form of the CredDelete method which supports a Context parameter
func (sdsaas *SdsaasV1) CredDeleteWithContext(ctx context.Context, credDeleteOptions *CredDeleteOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(credDeleteOptions, "credDeleteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(credDeleteOptions, "credDeleteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/v1/object/workspace/credentials`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range credDeleteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "CredDelete")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddQuery("access_key", fmt.Sprint(*credDeleteOptions.AccessKey))

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = sdsaas.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "cred_delete", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// Cert : Retrieves the S3 SSL certificate expiration date and status
// This request retrieves the S3 SSL certificate expiration date and status.
func (sdsaas *SdsaasV1) Cert(certOptions *CertOptions) (result *CertificateFound, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.CertWithContext(context.Background(), certOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CertWithContext is an alternate form of the Cert method which supports a Context parameter
func (sdsaas *SdsaasV1) CertWithContext(ctx context.Context, certOptions *CertOptions) (result *CertificateFound, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(certOptions, "certOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/v1/object/certificate/s3`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range certOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "Cert")
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
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "cert", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCertificateFound)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CertUpload : Creates/updates the S3 SSL Certificates
// Updates the S3 SSL Certificates or creates them if they do not exist.
func (sdsaas *SdsaasV1) CertUpload(certUploadOptions *CertUploadOptions) (result *CertificateUpdated, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.CertUploadWithContext(context.Background(), certUploadOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CertUploadWithContext is an alternate form of the CertUpload method which supports a Context parameter
func (sdsaas *SdsaasV1) CertUploadWithContext(ctx context.Context, certUploadOptions *CertUploadOptions) (result *CertificateUpdated, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(certUploadOptions, "certUploadOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(certUploadOptions, "certUploadOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/v1/object/certificate/s3`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range certUploadOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "CertUpload")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/octect-stream")

	_, err = builder.SetBodyContent("application/octect-stream", nil, nil, certUploadOptions.Body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "cert_upload", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCertificateUpdated)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// Hosts : Lists all hosts and all host IDs
// This request lists all hosts and host IDs.
func (sdsaas *SdsaasV1) Hosts(hostsOptions *HostsOptions) (result *HostCollection, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.HostsWithContext(context.Background(), hostsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// HostsWithContext is an alternate form of the Hosts method which supports a Context parameter
func (sdsaas *SdsaasV1) HostsWithContext(ctx context.Context, hostsOptions *HostsOptions) (result *HostCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(hostsOptions, "hostsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/hosts`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range hostsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "Hosts")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if hostsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*hostsOptions.Limit))
	}
	if hostsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*hostsOptions.Name))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "hosts", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalHostCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// HostCreate : Creates a host
// This request creates a new host from a host template object.
func (sdsaas *SdsaasV1) HostCreate(hostCreateOptions *HostCreateOptions) (result *Host, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.HostCreateWithContext(context.Background(), hostCreateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// HostCreateWithContext is an alternate form of the HostCreate method which supports a Context parameter
func (sdsaas *SdsaasV1) HostCreateWithContext(ctx context.Context, hostCreateOptions *HostCreateOptions) (result *Host, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(hostCreateOptions, "hostCreateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(hostCreateOptions, "hostCreateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/hosts`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range hostCreateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "HostCreate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if hostCreateOptions.Nqn != nil {
		body["nqn"] = hostCreateOptions.Nqn
	}
	if hostCreateOptions.Name != nil {
		body["name"] = hostCreateOptions.Name
	}
	if hostCreateOptions.Volumes != nil {
		body["volumes"] = hostCreateOptions.Volumes
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
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "host_create", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalHost)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// Host : Retrieve a host by ID
// This request retrieves a host specified by the host ID.
func (sdsaas *SdsaasV1) Host(hostOptions *HostOptions) (result *Host, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.HostWithContext(context.Background(), hostOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// HostWithContext is an alternate form of the Host method which supports a Context parameter
func (sdsaas *SdsaasV1) HostWithContext(ctx context.Context, hostOptions *HostOptions) (result *Host, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(hostOptions, "hostOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(hostOptions, "hostOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"host_id": *hostOptions.HostID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/hosts/{host_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range hostOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "Host")
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
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "host", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalHost)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// HostUpdate : Update a host
// This request updates a host with the information in a provided host patch.
func (sdsaas *SdsaasV1) HostUpdate(hostUpdateOptions *HostUpdateOptions) (result *Host, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.HostUpdateWithContext(context.Background(), hostUpdateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// HostUpdateWithContext is an alternate form of the HostUpdate method which supports a Context parameter
func (sdsaas *SdsaasV1) HostUpdateWithContext(ctx context.Context, hostUpdateOptions *HostUpdateOptions) (result *Host, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(hostUpdateOptions, "hostUpdateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(hostUpdateOptions, "hostUpdateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"host_id": *hostUpdateOptions.HostID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/hosts/{host_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range hostUpdateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "HostUpdate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if hostUpdateOptions.HostPatch != nil {
		_, err = builder.SetBodyContentJSON(hostUpdateOptions.HostPatch)
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

	var rawResponse map[string]json.RawMessage
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "host_update", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalHost)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// HostDelete : Delete a specific host
// This request deletes a host using the ID.
func (sdsaas *SdsaasV1) HostDelete(hostDeleteOptions *HostDeleteOptions) (response *core.DetailedResponse, err error) {
	response, err = sdsaas.HostDeleteWithContext(context.Background(), hostDeleteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// HostDeleteWithContext is an alternate form of the HostDelete method which supports a Context parameter
func (sdsaas *SdsaasV1) HostDeleteWithContext(ctx context.Context, hostDeleteOptions *HostDeleteOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(hostDeleteOptions, "hostDeleteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(hostDeleteOptions, "hostDeleteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"host_id": *hostDeleteOptions.HostID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/hosts/{host_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range hostDeleteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "HostDelete")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = sdsaas.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "host_delete", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// HostVolDeleteall : Deletes all the volume mappings for a given host
// This request deletes all volume mappings associated with a specific host ID.
func (sdsaas *SdsaasV1) HostVolDeleteall(hostVolDeleteallOptions *HostVolDeleteallOptions) (response *core.DetailedResponse, err error) {
	response, err = sdsaas.HostVolDeleteallWithContext(context.Background(), hostVolDeleteallOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// HostVolDeleteallWithContext is an alternate form of the HostVolDeleteall method which supports a Context parameter
func (sdsaas *SdsaasV1) HostVolDeleteallWithContext(ctx context.Context, hostVolDeleteallOptions *HostVolDeleteallOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(hostVolDeleteallOptions, "hostVolDeleteallOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(hostVolDeleteallOptions, "hostVolDeleteallOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"host_id": *hostVolDeleteallOptions.HostID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/hosts/{host_id}/volumes`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range hostVolDeleteallOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "HostVolDeleteall")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = sdsaas.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "host_vol_deleteall", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// HostVolDelete : Deletes the given volume mapping for a specific host
// This request deletes a particular volume mapped from the host.
func (sdsaas *SdsaasV1) HostVolDelete(hostVolDeleteOptions *HostVolDeleteOptions) (response *core.DetailedResponse, err error) {
	response, err = sdsaas.HostVolDeleteWithContext(context.Background(), hostVolDeleteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// HostVolDeleteWithContext is an alternate form of the HostVolDelete method which supports a Context parameter
func (sdsaas *SdsaasV1) HostVolDeleteWithContext(ctx context.Context, hostVolDeleteOptions *HostVolDeleteOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(hostVolDeleteOptions, "hostVolDeleteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(hostVolDeleteOptions, "hostVolDeleteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"host_id":   *hostVolDeleteOptions.HostID,
		"volume_id": *hostVolDeleteOptions.VolumeID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/hosts/{host_id}/volumes/{volume_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range hostVolDeleteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "HostVolDelete")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = sdsaas.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "host_vol_delete", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// HostVolUpdate : Maps the given volume to the given host
// This request creates a volume mapping to the given host ID.
func (sdsaas *SdsaasV1) HostVolUpdate(hostVolUpdateOptions *HostVolUpdateOptions) (result *Host, response *core.DetailedResponse, err error) {
	result, response, err = sdsaas.HostVolUpdateWithContext(context.Background(), hostVolUpdateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// HostVolUpdateWithContext is an alternate form of the HostVolUpdate method which supports a Context parameter
func (sdsaas *SdsaasV1) HostVolUpdateWithContext(ctx context.Context, hostVolUpdateOptions *HostVolUpdateOptions) (result *Host, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(hostVolUpdateOptions, "hostVolUpdateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(hostVolUpdateOptions, "hostVolUpdateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"host_id":   *hostVolUpdateOptions.HostID,
		"volume_id": *hostVolUpdateOptions.VolumeID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sdsaas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sdsaas.Service.Options.URL, `/hosts/{host_id}/volumes/{volume_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range hostVolUpdateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("sdsaas", "V1", "HostVolUpdate")
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
	response, err = sdsaas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "host_vol_update", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalHost)
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

// CertOptions : The Cert options.
type CertOptions struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCertOptions : Instantiate CertOptions
func (*SdsaasV1) NewCertOptions() *CertOptions {
	return &CertOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *CertOptions) SetHeaders(param map[string]string) *CertOptions {
	options.Headers = param
	return options
}

// CertUploadOptions : The CertUpload options.
type CertUploadOptions struct {
	// The request body containing the S3 TLS certificate. The CLI will accept certificate files of any type, but they must
	// be in proper .pem format.
	Body io.ReadCloser `json:"body" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCertUploadOptions : Instantiate CertUploadOptions
func (*SdsaasV1) NewCertUploadOptions(body io.ReadCloser) *CertUploadOptions {
	return &CertUploadOptions{
		Body: body,
	}
}

// SetBody : Allow user to set Body
func (_options *CertUploadOptions) SetBody(body io.ReadCloser) *CertUploadOptions {
	_options.Body = body
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CertUploadOptions) SetHeaders(param map[string]string) *CertUploadOptions {
	options.Headers = param
	return options
}

// CertificateFound : The responese object for certificate GET operations.
type CertificateFound struct {
	// The expiration date of the certificate.
	ExpirationDate *string `json:"expiration_date,omitempty"`

	// The boolean value of the expiration status.
	Expired *bool `json:"expired,omitempty"`
}

// UnmarshalCertificateFound unmarshals an instance of CertificateFound from the specified map of raw messages.
func UnmarshalCertificateFound(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateFound)
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		err = core.SDKErrorf(err, "", "expiration_date-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expired", &obj.Expired)
	if err != nil {
		err = core.SDKErrorf(err, "", "expired-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CertificateUpdated : The response object for certificate POST operations.
type CertificateUpdated struct {
	// An array of certificate error codes and their descriptions.
	Error []map[string]string `json:"error,omitempty"`

	// The boolean valid status of the certificate.
	ValidCertificate *bool `json:"valid_certificate,omitempty"`

	// The boolean valid status of the key.
	ValidKey *bool `json:"valid_key,omitempty"`
}

// UnmarshalCertificateUpdated unmarshals an instance of CertificateUpdated from the specified map of raw messages.
func UnmarshalCertificateUpdated(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateUpdated)
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		err = core.SDKErrorf(err, "", "error-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "valid_certificate", &obj.ValidCertificate)
	if err != nil {
		err = core.SDKErrorf(err, "", "valid_certificate-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "valid_key", &obj.ValidKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "valid_key-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CredCreateOptions : The CredCreate options.
type CredCreateOptions struct {
	// Access key to update or set.
	AccessKey *string `json:"access_key" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCredCreateOptions : Instantiate CredCreateOptions
func (*SdsaasV1) NewCredCreateOptions(accessKey string) *CredCreateOptions {
	return &CredCreateOptions{
		AccessKey: core.StringPtr(accessKey),
	}
}

// SetAccessKey : Allow user to set AccessKey
func (_options *CredCreateOptions) SetAccessKey(accessKey string) *CredCreateOptions {
	_options.AccessKey = core.StringPtr(accessKey)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CredCreateOptions) SetHeaders(param map[string]string) *CredCreateOptions {
	options.Headers = param
	return options
}

// CredDeleteOptions : The CredDelete options.
type CredDeleteOptions struct {
	// Access key to update or set.
	AccessKey *string `json:"access_key" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCredDeleteOptions : Instantiate CredDeleteOptions
func (*SdsaasV1) NewCredDeleteOptions(accessKey string) *CredDeleteOptions {
	return &CredDeleteOptions{
		AccessKey: core.StringPtr(accessKey),
	}
}

// SetAccessKey : Allow user to set AccessKey
func (_options *CredDeleteOptions) SetAccessKey(accessKey string) *CredDeleteOptions {
	_options.AccessKey = core.StringPtr(accessKey)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CredDeleteOptions) SetHeaders(param map[string]string) *CredDeleteOptions {
	options.Headers = param
	return options
}

// CredentialsFound : The response object for credential GET operations.
type CredentialsFound struct {
	// Collection of access keys.
	AccessKeys []string `json:"access_keys,omitempty"`
}

// UnmarshalCredentialsFound unmarshals an instance of CredentialsFound from the specified map of raw messages.
func UnmarshalCredentialsFound(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CredentialsFound)
	err = core.UnmarshalPrimitive(m, "access_keys", &obj.AccessKeys)
	if err != nil {
		err = core.SDKErrorf(err, "", "access_keys-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CredentialsUpdated : The response object for credential POST operations.
type CredentialsUpdated struct {
	// The user created access key.
	AccessKey *string `json:"access_key,omitempty"`

	// The key material associated with and access key.
	SecretKey *string `json:"secret_key,omitempty"`
}

// UnmarshalCredentialsUpdated unmarshals an instance of CredentialsUpdated from the specified map of raw messages.
func UnmarshalCredentialsUpdated(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CredentialsUpdated)
	err = core.UnmarshalPrimitive(m, "access_key", &obj.AccessKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "access_key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_key", &obj.SecretKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_key-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CredsOptions : The Creds options.
type CredsOptions struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCredsOptions : Instantiate CredsOptions
func (*SdsaasV1) NewCredsOptions() *CredsOptions {
	return &CredsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *CredsOptions) SetHeaders(param map[string]string) *CredsOptions {
	options.Headers = param
	return options
}

// Host : The host object.
type Host struct {
	// The date and time that the host was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// The unique identifier for this host.
	ID *string `json:"id,omitempty"`

	// The name for this host. The name must not be used by another host.  If unspecified, the name will be a hyphenated
	// list of randomly-selected words.
	Name *string `json:"name,omitempty"`

	// The NQN of the host configured in customer's environment.
	Nqn *string `json:"nqn" validate:"required"`

	// The host-to-volume map.
	Volumes []VolumeMappingReference `json:"volumes,omitempty"`
}

// UnmarshalHost unmarshals an instance of Host from the specified map of raw messages.
func UnmarshalHost(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Host)
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "nqn", &obj.Nqn)
	if err != nil {
		err = core.SDKErrorf(err, "", "nqn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "volumes", &obj.Volumes, UnmarshalVolumeMappingReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "volumes-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// HostCollection : A collection of hosts at a particular endpoint with the total number found.  Any hosts beyond the return limit are
// found in the Next link.
type HostCollection struct {
	// A link to the first page of resources.
	First *PageLink `json:"first" validate:"required"`

	// Collection of hosts.
	Hosts []Host `json:"hosts" validate:"required"`

	// The maximum number of resources that can be returned by the request.
	Limit *int64 `json:"limit,omitempty"`

	// A link to the next page of resources. This property is present for all pages except the last page.
	Next *PageLink `json:"next,omitempty"`

	// The total number of resources across all pages.
	TotalCount *int64 `json:"total_count" validate:"required"`
}

// UnmarshalHostCollection unmarshals an instance of HostCollection from the specified map of raw messages.
func UnmarshalHostCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HostCollection)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "hosts", &obj.Hosts, UnmarshalHost)
	if err != nil {
		err = core.SDKErrorf(err, "", "hosts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// HostCreateOptions : The HostCreate options.
type HostCreateOptions struct {
	// The NQN of the host configured in customer's environment.
	Nqn *string `json:"nqn" validate:"required"`

	// The name for this host. The name must not be used by another host.  If unspecified, the name will be a hyphenated
	// list of randomly-selected words.
	Name *string `json:"name,omitempty"`

	// The unique identifier of the volume to be mapped to this host.  Must be in the form '['volume_id':
	// '1a6b7274-678d-4dfb-8981-c71dd9d4daa5']'.  If curly braces {} are used to separate volumes, double quotes must be
	// used instead.
	Volumes []VolumeMappingIdentity `json:"volumes,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewHostCreateOptions : Instantiate HostCreateOptions
func (*SdsaasV1) NewHostCreateOptions(nqn string) *HostCreateOptions {
	return &HostCreateOptions{
		Nqn: core.StringPtr(nqn),
	}
}

// SetNqn : Allow user to set Nqn
func (_options *HostCreateOptions) SetNqn(nqn string) *HostCreateOptions {
	_options.Nqn = core.StringPtr(nqn)
	return _options
}

// SetName : Allow user to set Name
func (_options *HostCreateOptions) SetName(name string) *HostCreateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetVolumes : Allow user to set Volumes
func (_options *HostCreateOptions) SetVolumes(volumes []VolumeMappingIdentity) *HostCreateOptions {
	_options.Volumes = volumes
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HostCreateOptions) SetHeaders(param map[string]string) *HostCreateOptions {
	options.Headers = param
	return options
}

// HostDeleteOptions : The HostDelete options.
type HostDeleteOptions struct {
	// A unique host ID.
	HostID *string `json:"host_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewHostDeleteOptions : Instantiate HostDeleteOptions
func (*SdsaasV1) NewHostDeleteOptions(hostID string) *HostDeleteOptions {
	return &HostDeleteOptions{
		HostID: core.StringPtr(hostID),
	}
}

// SetHostID : Allow user to set HostID
func (_options *HostDeleteOptions) SetHostID(hostID string) *HostDeleteOptions {
	_options.HostID = core.StringPtr(hostID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HostDeleteOptions) SetHeaders(param map[string]string) *HostDeleteOptions {
	options.Headers = param
	return options
}

// HostMapping : HostMapping struct
type HostMapping struct {
	// Unique identifer of the host.
	HostID *string `json:"host_id,omitempty"`

	// Unique name of the host.
	HostName *string `json:"host_name,omitempty"`

	// The NQN of the host configured in customer's environment.
	HostNqn *string `json:"host_nqn,omitempty"`
}

// UnmarshalHostMapping unmarshals an instance of HostMapping from the specified map of raw messages.
func UnmarshalHostMapping(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HostMapping)
	err = core.UnmarshalPrimitive(m, "host_id", &obj.HostID)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_name", &obj.HostName)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_nqn", &obj.HostNqn)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_nqn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// HostOptions : The Host options.
type HostOptions struct {
	// A unique host ID.
	HostID *string `json:"host_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewHostOptions : Instantiate HostOptions
func (*SdsaasV1) NewHostOptions(hostID string) *HostOptions {
	return &HostOptions{
		HostID: core.StringPtr(hostID),
	}
}

// SetHostID : Allow user to set HostID
func (_options *HostOptions) SetHostID(hostID string) *HostOptions {
	_options.HostID = core.StringPtr(hostID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HostOptions) SetHeaders(param map[string]string) *HostOptions {
	options.Headers = param
	return options
}

// HostPatch : The host PATCH request body.
type HostPatch struct {
	// The name for this Host. The name must not be used by another host.
	Name *string `json:"name,omitempty"`
}

// UnmarshalHostPatch unmarshals an instance of HostPatch from the specified map of raw messages.
func UnmarshalHostPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HostPatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the HostPatch
func (hostPatch *HostPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(hostPatch.Name) {
		_patch["name"] = hostPatch.Name
	}

	return
}

// HostUpdateOptions : The HostUpdate options.
type HostUpdateOptions struct {
	// A unique host ID.
	HostID *string `json:"host_id" validate:"required,ne="`

	// JSON Merge-Patch content for host_update.
	HostPatch map[string]interface{} `json:"Host_patch,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewHostUpdateOptions : Instantiate HostUpdateOptions
func (*SdsaasV1) NewHostUpdateOptions(hostID string) *HostUpdateOptions {
	return &HostUpdateOptions{
		HostID: core.StringPtr(hostID),
	}
}

// SetHostID : Allow user to set HostID
func (_options *HostUpdateOptions) SetHostID(hostID string) *HostUpdateOptions {
	_options.HostID = core.StringPtr(hostID)
	return _options
}

// SetHostPatch : Allow user to set HostPatch
func (_options *HostUpdateOptions) SetHostPatch(hostPatch map[string]interface{}) *HostUpdateOptions {
	_options.HostPatch = hostPatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HostUpdateOptions) SetHeaders(param map[string]string) *HostUpdateOptions {
	options.Headers = param
	return options
}

// HostVolDeleteOptions : The HostVolDelete options.
type HostVolDeleteOptions struct {
	// A unique host ID.
	HostID *string `json:"host_id" validate:"required,ne="`

	// A unique volume ID.
	VolumeID *string `json:"volume_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewHostVolDeleteOptions : Instantiate HostVolDeleteOptions
func (*SdsaasV1) NewHostVolDeleteOptions(hostID string, volumeID string) *HostVolDeleteOptions {
	return &HostVolDeleteOptions{
		HostID:   core.StringPtr(hostID),
		VolumeID: core.StringPtr(volumeID),
	}
}

// SetHostID : Allow user to set HostID
func (_options *HostVolDeleteOptions) SetHostID(hostID string) *HostVolDeleteOptions {
	_options.HostID = core.StringPtr(hostID)
	return _options
}

// SetVolumeID : Allow user to set VolumeID
func (_options *HostVolDeleteOptions) SetVolumeID(volumeID string) *HostVolDeleteOptions {
	_options.VolumeID = core.StringPtr(volumeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HostVolDeleteOptions) SetHeaders(param map[string]string) *HostVolDeleteOptions {
	options.Headers = param
	return options
}

// HostVolDeleteallOptions : The HostVolDeleteall options.
type HostVolDeleteallOptions struct {
	// A unique host ID.
	HostID *string `json:"host_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewHostVolDeleteallOptions : Instantiate HostVolDeleteallOptions
func (*SdsaasV1) NewHostVolDeleteallOptions(hostID string) *HostVolDeleteallOptions {
	return &HostVolDeleteallOptions{
		HostID: core.StringPtr(hostID),
	}
}

// SetHostID : Allow user to set HostID
func (_options *HostVolDeleteallOptions) SetHostID(hostID string) *HostVolDeleteallOptions {
	_options.HostID = core.StringPtr(hostID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HostVolDeleteallOptions) SetHeaders(param map[string]string) *HostVolDeleteallOptions {
	options.Headers = param
	return options
}

// HostVolUpdateOptions : The HostVolUpdate options.
type HostVolUpdateOptions struct {
	// A unique host ID.
	HostID *string `json:"host_id" validate:"required,ne="`

	// A unique volume ID.
	VolumeID *string `json:"volume_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewHostVolUpdateOptions : Instantiate HostVolUpdateOptions
func (*SdsaasV1) NewHostVolUpdateOptions(hostID string, volumeID string) *HostVolUpdateOptions {
	return &HostVolUpdateOptions{
		HostID:   core.StringPtr(hostID),
		VolumeID: core.StringPtr(volumeID),
	}
}

// SetHostID : Allow user to set HostID
func (_options *HostVolUpdateOptions) SetHostID(hostID string) *HostVolUpdateOptions {
	_options.HostID = core.StringPtr(hostID)
	return _options
}

// SetVolumeID : Allow user to set VolumeID
func (_options *HostVolUpdateOptions) SetVolumeID(volumeID string) *HostVolUpdateOptions {
	_options.VolumeID = core.StringPtr(volumeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HostVolUpdateOptions) SetHeaders(param map[string]string) *HostVolUpdateOptions {
	options.Headers = param
	return options
}

// HostsOptions : The Hosts options.
type HostsOptions struct {
	// The number of resources to return on a page.
	Limit *int64 `json:"limit,omitempty"`

	// Filters the collection of resources by name.
	Name *string `json:"name,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewHostsOptions : Instantiate HostsOptions
func (*SdsaasV1) NewHostsOptions() *HostsOptions {
	return &HostsOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *HostsOptions) SetLimit(limit int64) *HostsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetName : Allow user to set Name
func (_options *HostsOptions) SetName(name string) *HostsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HostsOptions) SetHeaders(param map[string]string) *HostsOptions {
	options.Headers = param
	return options
}

// NetworkInfoReference : NetworkInfoReference struct
type NetworkInfoReference struct {
	// Network information for volume/host mappings.
	GatewayIP *string `json:"gateway_ip,omitempty"`

	// Network information for volume/host mappings.
	Port *int64 `json:"port,omitempty"`
}

// UnmarshalNetworkInfoReference unmarshals an instance of NetworkInfoReference from the specified map of raw messages.
func UnmarshalNetworkInfoReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NetworkInfoReference)
	err = core.UnmarshalPrimitive(m, "gateway_ip", &obj.GatewayIP)
	if err != nil {
		err = core.SDKErrorf(err, "", "gateway_ip-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageLink : PageLink struct
type PageLink struct {
	// The URL for a page of resources.
	Href *string `json:"href,omitempty"`
}

// UnmarshalPageLink unmarshals an instance of PageLink from the specified map of raw messages.
func UnmarshalPageLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageLink)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StorageIdentifiersReference : Storage network and ID information associated with a volume/host mapping.
type StorageIdentifiersReference struct {
	// The storage ID associated with a volume/host mapping.
	ID *string `json:"id,omitempty"`

	// The namespace ID associated with a volume/host mapping.
	NamespaceID *int64 `json:"namespace_id,omitempty"`

	// The namespace UUID associated with a volume/host mapping.
	NamespaceUUID *string `json:"namespace_uuid,omitempty"`

	// The IP and port for volume/host mappings.
	NetworkInfo []NetworkInfoReference `json:"network_info,omitempty"`
}

// UnmarshalStorageIdentifiersReference unmarshals an instance of StorageIdentifiersReference from the specified map of raw messages.
func UnmarshalStorageIdentifiersReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StorageIdentifiersReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace_id", &obj.NamespaceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "namespace_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace_uuid", &obj.NamespaceUUID)
	if err != nil {
		err = core.SDKErrorf(err, "", "namespace_uuid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "network_info", &obj.NetworkInfo, UnmarshalNetworkInfoReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "network_info-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Volume : The volume metadata prototype.
type Volume struct {
	// The maximum bandwidth (in megabits per second) for the volume.
	Bandwidth *int64 `json:"bandwidth,omitempty"`

	// The capacity of the volume (in gigabytes).
	Capacity *int64 `json:"capacity,omitempty"`

	// The date and time that the volume was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// List of host details that volume is mapped to.
	Hosts []HostMapping `json:"hosts,omitempty"`

	// The volume profile id.
	ID *string `json:"id,omitempty"`

	// Iops The maximum I/O operations per second (IOPS) for this volume.
	Iops *int64 `json:"iops,omitempty"`

	// The name of the volume.
	Name *string `json:"name,omitempty"`

	// The resource type of the volume.
	ResourceType *string `json:"resource_type,omitempty"`

	// The current status of the volume.
	Status *string `json:"status,omitempty"`

	// Reasons for the current status of the volume.
	StatusReasons []string `json:"status_reasons,omitempty"`
}

// UnmarshalVolume unmarshals an instance of Volume from the specified map of raw messages.
func UnmarshalVolume(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Volume)
	err = core.UnmarshalPrimitive(m, "bandwidth", &obj.Bandwidth)
	if err != nil {
		err = core.SDKErrorf(err, "", "bandwidth-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "capacity", &obj.Capacity)
	if err != nil {
		err = core.SDKErrorf(err, "", "capacity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "hosts", &obj.Hosts, UnmarshalHostMapping)
	if err != nil {
		err = core.SDKErrorf(err, "", "hosts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iops", &obj.Iops)
	if err != nil {
		err = core.SDKErrorf(err, "", "iops-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status_reasons", &obj.StatusReasons)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_reasons-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VolumeCollection : Volume object showing the results of the GET volumes operation.
type VolumeCollection struct {
	// The first page of volume objects.
	First *PageLink `json:"first,omitempty"`

	// The maximum number of volumes retrieved with the volumes command.
	Limit *int64 `json:"limit,omitempty"`

	// The next page of volume objects.
	Next *PageLink `json:"next,omitempty"`

	// Total number of volumes retrieved.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// List of volumes retrieved.
	Volumes []Volume `json:"volumes,omitempty"`
}

// UnmarshalVolumeCollection unmarshals an instance of VolumeCollection from the specified map of raw messages.
func UnmarshalVolumeCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VolumeCollection)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "volumes", &obj.Volumes, UnmarshalVolume)
	if err != nil {
		err = core.SDKErrorf(err, "", "volumes-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VolumeCreateOptions : The VolumeCreate options.
type VolumeCreateOptions struct {
	// The capacity to use for the volume (in gigabytes). The specified value must be within the capacity range of the
	// volume's profile.
	Capacity *int64 `json:"capacity" validate:"required"`

	// The name for this volume. The name must not be used by another volume. If unspecified, the name will be a hyphenated
	// list of randomly-selected words.
	Name *string `json:"name,omitempty"`

	// The host nqn.
	Hostnqnstring *string `json:"hostnqnstring,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewVolumeCreateOptions : Instantiate VolumeCreateOptions
func (*SdsaasV1) NewVolumeCreateOptions(capacity int64) *VolumeCreateOptions {
	return &VolumeCreateOptions{
		Capacity: core.Int64Ptr(capacity),
	}
}

// SetCapacity : Allow user to set Capacity
func (_options *VolumeCreateOptions) SetCapacity(capacity int64) *VolumeCreateOptions {
	_options.Capacity = core.Int64Ptr(capacity)
	return _options
}

// SetName : Allow user to set Name
func (_options *VolumeCreateOptions) SetName(name string) *VolumeCreateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHostnqnstring : Allow user to set Hostnqnstring
func (_options *VolumeCreateOptions) SetHostnqnstring(hostnqnstring string) *VolumeCreateOptions {
	_options.Hostnqnstring = core.StringPtr(hostnqnstring)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *VolumeCreateOptions) SetHeaders(param map[string]string) *VolumeCreateOptions {
	options.Headers = param
	return options
}

// VolumeDeleteOptions : The VolumeDelete options.
type VolumeDeleteOptions struct {
	// The volume profile id.
	VolumeID *string `json:"volume_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewVolumeDeleteOptions : Instantiate VolumeDeleteOptions
func (*SdsaasV1) NewVolumeDeleteOptions(volumeID string) *VolumeDeleteOptions {
	return &VolumeDeleteOptions{
		VolumeID: core.StringPtr(volumeID),
	}
}

// SetVolumeID : Allow user to set VolumeID
func (_options *VolumeDeleteOptions) SetVolumeID(volumeID string) *VolumeDeleteOptions {
	_options.VolumeID = core.StringPtr(volumeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *VolumeDeleteOptions) SetHeaders(param map[string]string) *VolumeDeleteOptions {
	options.Headers = param
	return options
}

// VolumeMappingIdentity : VolumeMappingIdentity struct
type VolumeMappingIdentity struct {
	// The volume ID that needs to be mapped with a host.
	VolumeID *string `json:"volume_id" validate:"required"`
}

// NewVolumeMappingIdentity : Instantiate VolumeMappingIdentity (Generic Model Constructor)
func (*SdsaasV1) NewVolumeMappingIdentity(volumeID string) (_model *VolumeMappingIdentity, err error) {
	_model = &VolumeMappingIdentity{
		VolumeID: core.StringPtr(volumeID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalVolumeMappingIdentity unmarshals an instance of VolumeMappingIdentity from the specified map of raw messages.
func UnmarshalVolumeMappingIdentity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VolumeMappingIdentity)
	err = core.UnmarshalPrimitive(m, "volume_id", &obj.VolumeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "volume_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VolumeMappingReference : VolumeMappingReference struct
type VolumeMappingReference struct {
	// The current status of a volume/host mapping attempt.
	Status *string `json:"status,omitempty"`

	// The volume ID that needs to be mapped with a host.
	VolumeID *string `json:"volume_id" validate:"required"`

	// The volume name.
	VolumeName *string `json:"volume_name" validate:"required"`

	// Storage network and ID information associated with a volume/host mapping.
	StorageIdentifiers *StorageIdentifiersReference `json:"storage_identifiers,omitempty"`
}

// UnmarshalVolumeMappingReference unmarshals an instance of VolumeMappingReference from the specified map of raw messages.
func UnmarshalVolumeMappingReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VolumeMappingReference)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "volume_id", &obj.VolumeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "volume_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "volume_name", &obj.VolumeName)
	if err != nil {
		err = core.SDKErrorf(err, "", "volume_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "storage_identifiers", &obj.StorageIdentifiers, UnmarshalStorageIdentifiersReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "storage_identifiers-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VolumeOptions : The Volume options.
type VolumeOptions struct {
	// The volume profile id.
	VolumeID *string `json:"volume_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewVolumeOptions : Instantiate VolumeOptions
func (*SdsaasV1) NewVolumeOptions(volumeID string) *VolumeOptions {
	return &VolumeOptions{
		VolumeID: core.StringPtr(volumeID),
	}
}

// SetVolumeID : Allow user to set VolumeID
func (_options *VolumeOptions) SetVolumeID(volumeID string) *VolumeOptions {
	_options.VolumeID = core.StringPtr(volumeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *VolumeOptions) SetHeaders(param map[string]string) *VolumeOptions {
	options.Headers = param
	return options
}

// VolumePatch : Volume update request metadata.
type VolumePatch struct {
	// The capacity of the volume (in gigabytes).
	Capacity *int64 `json:"capacity,omitempty"`

	// The name of the volume.
	Name *string `json:"name,omitempty"`
}

// UnmarshalVolumePatch unmarshals an instance of VolumePatch from the specified map of raw messages.
func UnmarshalVolumePatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VolumePatch)
	err = core.UnmarshalPrimitive(m, "capacity", &obj.Capacity)
	if err != nil {
		err = core.SDKErrorf(err, "", "capacity-error", common.GetComponentInfo())
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

// AsPatch returns a generic map representation of the VolumePatch
func (volumePatch *VolumePatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(volumePatch.Capacity) {
		_patch["capacity"] = volumePatch.Capacity
	}
	if !core.IsNil(volumePatch.Name) {
		_patch["name"] = volumePatch.Name
	}

	return
}

// VolumeUpdateOptions : The VolumeUpdate options.
type VolumeUpdateOptions struct {
	// The volume profile id.
	VolumeID *string `json:"volume_id" validate:"required,ne="`

	// A JSON object containing volume information.
	VolumePatch map[string]interface{} `json:"Volume_patch,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewVolumeUpdateOptions : Instantiate VolumeUpdateOptions
func (*SdsaasV1) NewVolumeUpdateOptions(volumeID string) *VolumeUpdateOptions {
	return &VolumeUpdateOptions{
		VolumeID: core.StringPtr(volumeID),
	}
}

// SetVolumeID : Allow user to set VolumeID
func (_options *VolumeUpdateOptions) SetVolumeID(volumeID string) *VolumeUpdateOptions {
	_options.VolumeID = core.StringPtr(volumeID)
	return _options
}

// SetVolumePatch : Allow user to set VolumePatch
func (_options *VolumeUpdateOptions) SetVolumePatch(volumePatch map[string]interface{}) *VolumeUpdateOptions {
	_options.VolumePatch = volumePatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *VolumeUpdateOptions) SetHeaders(param map[string]string) *VolumeUpdateOptions {
	options.Headers = param
	return options
}

// VolumesOptions : The Volumes options.
type VolumesOptions struct {
	// The number of resources to return on a page.
	Limit *int64 `json:"limit,omitempty"`

	// Filters the collection of resources by name.
	Name *string `json:"name,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewVolumesOptions : Instantiate VolumesOptions
func (*SdsaasV1) NewVolumesOptions() *VolumesOptions {
	return &VolumesOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *VolumesOptions) SetLimit(limit int64) *VolumesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetName : Allow user to set Name
func (_options *VolumesOptions) SetName(name string) *VolumesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *VolumesOptions) SetHeaders(param map[string]string) *VolumesOptions {
	options.Headers = param
	return options
}
