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
 * IBM OpenAPI SDK Code Generator Version: 3.62.2-e5d4c32b-20221214-193750
 */

// Package ukov4 : Operations and models for the UkoV4 service
package ukov4

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
	common "github.com/IBM/ibm-hpcs-uko-sdk/common"
	"github.com/go-openapi/strfmt"
)

// UkoV4 : API for UKO used for key management.
//
// API Version: 4.7.9
type UkoV4 struct {
	Service *core.BaseService
}

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "uko"

// UkoV4Options : Service options
type UkoV4Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewUkoV4UsingExternalConfig : constructs an instance of UkoV4 with passed in options and external configuration.
func NewUkoV4UsingExternalConfig(options *UkoV4Options) (uko *UkoV4, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	uko, err = NewUkoV4(options)
	if err != nil {
		return
	}

	err = uko.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = uko.Service.SetServiceURL(options.URL)
	}
	return
}

// NewUkoV4 : constructs an instance of UkoV4 with passed in options.
func NewUkoV4(options *UkoV4Options) (service *UkoV4, err error) {
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

	service = &UkoV4{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "uko" suitable for processing requests.
func (uko *UkoV4) Clone() *UkoV4 {
	if core.IsNil(uko) {
		return nil
	}
	clone := *uko
	clone.Service = uko.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (uko *UkoV4) SetServiceURL(url string) error {
	return uko.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (uko *UkoV4) GetServiceURL() string {
	return uko.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (uko *UkoV4) SetDefaultHeaders(headers http.Header) {
	uko.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (uko *UkoV4) SetEnableGzipCompression(enableGzip bool) {
	uko.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (uko *UkoV4) GetEnableGzipCompression() bool {
	return uko.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (uko *UkoV4) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	uko.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (uko *UkoV4) DisableRetries() {
	uko.Service.DisableRetries()
}

// ListManagedKeys : List managed keys
// List all managed keys in the instance. It is possible to sort by the following parameters: name, algorithm, state,
// activation_date, deactivation_date, created_at, updated_at, size, vault.id.
func (uko *UkoV4) ListManagedKeys(listManagedKeysOptions *ListManagedKeysOptions) (result *ManagedKeyList, response *core.DetailedResponse, err error) {
	return uko.ListManagedKeysWithContext(context.Background(), listManagedKeysOptions)
}

// ListManagedKeysWithContext is an alternate form of the ListManagedKeys method which supports a Context parameter
func (uko *UkoV4) ListManagedKeysWithContext(ctx context.Context, listManagedKeysOptions *ListManagedKeysOptions) (result *ManagedKeyList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listManagedKeysOptions, "listManagedKeysOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listManagedKeysOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "ListManagedKeys")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listManagedKeysOptions.VaultID != nil {
		builder.AddQuery("vault.id", strings.Join(listManagedKeysOptions.VaultID, ","))
	}
	if listManagedKeysOptions.Algorithm != nil {
		builder.AddQuery("algorithm", strings.Join(listManagedKeysOptions.Algorithm, ","))
	}
	if listManagedKeysOptions.State != nil {
		builder.AddQuery("state", strings.Join(listManagedKeysOptions.State, ","))
	}
	if listManagedKeysOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listManagedKeysOptions.Limit))
	}
	if listManagedKeysOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listManagedKeysOptions.Offset))
	}
	if listManagedKeysOptions.Sort != nil {
		builder.AddQuery("sort", strings.Join(listManagedKeysOptions.Sort, ","))
	}
	if listManagedKeysOptions.Label != nil {
		builder.AddQuery("label", fmt.Sprint(*listManagedKeysOptions.Label))
	}
	if listManagedKeysOptions.ActivationDate != nil {
		builder.AddQuery("activation_date", fmt.Sprint(*listManagedKeysOptions.ActivationDate))
	}
	if listManagedKeysOptions.ActivationDateMin != nil {
		builder.AddQuery("activation_date_min", fmt.Sprint(*listManagedKeysOptions.ActivationDateMin))
	}
	if listManagedKeysOptions.ActivationDateMax != nil {
		builder.AddQuery("activation_date_max", fmt.Sprint(*listManagedKeysOptions.ActivationDateMax))
	}
	if listManagedKeysOptions.DeactivationDate != nil {
		builder.AddQuery("deactivation_date", fmt.Sprint(*listManagedKeysOptions.DeactivationDate))
	}
	if listManagedKeysOptions.DeactivationDateMin != nil {
		builder.AddQuery("deactivation_date_min", fmt.Sprint(*listManagedKeysOptions.DeactivationDateMin))
	}
	if listManagedKeysOptions.DeactivationDateMax != nil {
		builder.AddQuery("deactivation_date_max", fmt.Sprint(*listManagedKeysOptions.DeactivationDateMax))
	}
	if listManagedKeysOptions.ExpirationDate != nil {
		builder.AddQuery("expiration_date", fmt.Sprint(*listManagedKeysOptions.ExpirationDate))
	}
	if listManagedKeysOptions.ExpirationDateMin != nil {
		builder.AddQuery("expiration_date_min", fmt.Sprint(*listManagedKeysOptions.ExpirationDateMin))
	}
	if listManagedKeysOptions.ExpirationDateMax != nil {
		builder.AddQuery("expiration_date_max", fmt.Sprint(*listManagedKeysOptions.ExpirationDateMax))
	}
	if listManagedKeysOptions.CreatedAt != nil {
		builder.AddQuery("created_at", fmt.Sprint(*listManagedKeysOptions.CreatedAt))
	}
	if listManagedKeysOptions.CreatedAtMin != nil {
		builder.AddQuery("created_at_min", fmt.Sprint(*listManagedKeysOptions.CreatedAtMin))
	}
	if listManagedKeysOptions.CreatedAtMax != nil {
		builder.AddQuery("created_at_max", fmt.Sprint(*listManagedKeysOptions.CreatedAtMax))
	}
	if listManagedKeysOptions.UpdatedAt != nil {
		builder.AddQuery("updated_at", fmt.Sprint(*listManagedKeysOptions.UpdatedAt))
	}
	if listManagedKeysOptions.UpdatedAtMin != nil {
		builder.AddQuery("updated_at_min", fmt.Sprint(*listManagedKeysOptions.UpdatedAtMin))
	}
	if listManagedKeysOptions.UpdatedAtMax != nil {
		builder.AddQuery("updated_at_max", fmt.Sprint(*listManagedKeysOptions.UpdatedAtMax))
	}
	if listManagedKeysOptions.Size != nil {
		builder.AddQuery("size", fmt.Sprint(*listManagedKeysOptions.Size))
	}
	if listManagedKeysOptions.SizeMin != nil {
		builder.AddQuery("size_min", fmt.Sprint(*listManagedKeysOptions.SizeMin))
	}
	if listManagedKeysOptions.SizeMax != nil {
		builder.AddQuery("size_max", fmt.Sprint(*listManagedKeysOptions.SizeMax))
	}
	if listManagedKeysOptions.ReferencedKeystoresType != nil {
		builder.AddQuery("referenced_keystores[].type", strings.Join(listManagedKeysOptions.ReferencedKeystoresType, ","))
	}
	if listManagedKeysOptions.ReferencedKeystoresName != nil {
		builder.AddQuery("referenced_keystores[].name", strings.Join(listManagedKeysOptions.ReferencedKeystoresName, ","))
	}
	if listManagedKeysOptions.InstancesKeystoreType != nil {
		builder.AddQuery("instances[].keystore.type", strings.Join(listManagedKeysOptions.InstancesKeystoreType, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKeyList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateManagedKey : Create a managed key
// Creates a new key based on the supplied template. The template must exist in the system prior to this call.
func (uko *UkoV4) CreateManagedKey(createManagedKeyOptions *CreateManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	return uko.CreateManagedKeyWithContext(context.Background(), createManagedKeyOptions)
}

// CreateManagedKeyWithContext is an alternate form of the CreateManagedKey method which supports a Context parameter
func (uko *UkoV4) CreateManagedKeyWithContext(ctx context.Context, createManagedKeyOptions *CreateManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createManagedKeyOptions, "createManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createManagedKeyOptions, "createManagedKeyOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "CreateManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*createManagedKeyOptions.UKOVault))
	}

	body := make(map[string]interface{})
	if createManagedKeyOptions.TemplateName != nil {
		body["template_name"] = createManagedKeyOptions.TemplateName
	}
	if createManagedKeyOptions.Vault != nil {
		body["vault"] = createManagedKeyOptions.Vault
	}
	if createManagedKeyOptions.Label != nil {
		body["label"] = createManagedKeyOptions.Label
	}
	if createManagedKeyOptions.Tags != nil {
		body["tags"] = createManagedKeyOptions.Tags
	}
	if createManagedKeyOptions.Description != nil {
		body["description"] = createManagedKeyOptions.Description
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
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteManagedKey : Delete a managed key
// Deletes a managed key by ID from the vault. A key must be in a 'destroyed' state for it to be eligible for deletion.
func (uko *UkoV4) DeleteManagedKey(deleteManagedKeyOptions *DeleteManagedKeyOptions) (response *core.DetailedResponse, err error) {
	return uko.DeleteManagedKeyWithContext(context.Background(), deleteManagedKeyOptions)
}

// DeleteManagedKeyWithContext is an alternate form of the DeleteManagedKey method which supports a Context parameter
func (uko *UkoV4) DeleteManagedKeyWithContext(ctx context.Context, deleteManagedKeyOptions *DeleteManagedKeyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteManagedKeyOptions, "deleteManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteManagedKeyOptions, "deleteManagedKeyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteManagedKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "DeleteManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*deleteManagedKeyOptions.UKOVault))
	}
	if deleteManagedKeyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deleteManagedKeyOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = uko.Service.Request(request, nil)

	return
}

// GetManagedKey : Retrieve a managed key
// Retrieves a managed key and its details by specifying the ID.
func (uko *UkoV4) GetManagedKey(getManagedKeyOptions *GetManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	return uko.GetManagedKeyWithContext(context.Background(), getManagedKeyOptions)
}

// GetManagedKeyWithContext is an alternate form of the GetManagedKey method which supports a Context parameter
func (uko *UkoV4) GetManagedKeyWithContext(ctx context.Context, getManagedKeyOptions *GetManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getManagedKeyOptions, "getManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getManagedKeyOptions, "getManagedKeyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getManagedKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "GetManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*getManagedKeyOptions.UKOVault))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateManagedKey : Update a managed key
// Updates attributes of a managed key. It is only possible to modify the key's state separately from other changes.
// Changing a key's state affects its availablity for crypto operations in keystores.
func (uko *UkoV4) UpdateManagedKey(updateManagedKeyOptions *UpdateManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	return uko.UpdateManagedKeyWithContext(context.Background(), updateManagedKeyOptions)
}

// UpdateManagedKeyWithContext is an alternate form of the UpdateManagedKey method which supports a Context parameter
func (uko *UkoV4) UpdateManagedKeyWithContext(ctx context.Context, updateManagedKeyOptions *UpdateManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateManagedKeyOptions, "updateManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateManagedKeyOptions, "updateManagedKeyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateManagedKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "UpdateManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*updateManagedKeyOptions.UKOVault))
	}
	if updateManagedKeyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateManagedKeyOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateManagedKeyOptions.Label != nil {
		body["label"] = updateManagedKeyOptions.Label
	}
	if updateManagedKeyOptions.ActivationDate != nil {
		body["activation_date"] = updateManagedKeyOptions.ActivationDate
	}
	if updateManagedKeyOptions.ExpirationDate != nil {
		body["expiration_date"] = updateManagedKeyOptions.ExpirationDate
	}
	if updateManagedKeyOptions.Tags != nil {
		body["tags"] = updateManagedKeyOptions.Tags
	}
	if updateManagedKeyOptions.Description != nil {
		body["description"] = updateManagedKeyOptions.Description
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
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListAssociatedResourcesForManagedKey : List associated resources for a managed key
// You can use this endpoint to obtain a list of resources associated with this managed key in IBM Cloud; which cloud
// resources are protected by the key you specify.
func (uko *UkoV4) ListAssociatedResourcesForManagedKey(listAssociatedResourcesForManagedKeyOptions *ListAssociatedResourcesForManagedKeyOptions) (result *AssociatedResourceList, response *core.DetailedResponse, err error) {
	return uko.ListAssociatedResourcesForManagedKeyWithContext(context.Background(), listAssociatedResourcesForManagedKeyOptions)
}

// ListAssociatedResourcesForManagedKeyWithContext is an alternate form of the ListAssociatedResourcesForManagedKey method which supports a Context parameter
func (uko *UkoV4) ListAssociatedResourcesForManagedKeyWithContext(ctx context.Context, listAssociatedResourcesForManagedKeyOptions *ListAssociatedResourcesForManagedKeyOptions) (result *AssociatedResourceList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAssociatedResourcesForManagedKeyOptions, "listAssociatedResourcesForManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listAssociatedResourcesForManagedKeyOptions, "listAssociatedResourcesForManagedKeyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *listAssociatedResourcesForManagedKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}/associated_resources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAssociatedResourcesForManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "ListAssociatedResourcesForManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listAssociatedResourcesForManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*listAssociatedResourcesForManagedKeyOptions.UKOVault))
	}

	if listAssociatedResourcesForManagedKeyOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAssociatedResourcesForManagedKeyOptions.Limit))
	}
	if listAssociatedResourcesForManagedKeyOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listAssociatedResourcesForManagedKeyOptions.Offset))
	}
	if listAssociatedResourcesForManagedKeyOptions.Sort != nil {
		builder.AddQuery("sort", strings.Join(listAssociatedResourcesForManagedKeyOptions.Sort, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAssociatedResourceList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetKeyDistributionStatusForKeystores : Retrieve distribution status for all keystores
// Return distribution status for all keystores for a key instance. If there's any problems reading the keystore status
// of the key instance, http code 200 will still be returned, and the error code will be returned alongside an 'error'
// keystore status.
func (uko *UkoV4) GetKeyDistributionStatusForKeystores(getKeyDistributionStatusForKeystoresOptions *GetKeyDistributionStatusForKeystoresOptions) (result *StatusInKeystores, response *core.DetailedResponse, err error) {
	return uko.GetKeyDistributionStatusForKeystoresWithContext(context.Background(), getKeyDistributionStatusForKeystoresOptions)
}

// GetKeyDistributionStatusForKeystoresWithContext is an alternate form of the GetKeyDistributionStatusForKeystores method which supports a Context parameter
func (uko *UkoV4) GetKeyDistributionStatusForKeystoresWithContext(ctx context.Context, getKeyDistributionStatusForKeystoresOptions *GetKeyDistributionStatusForKeystoresOptions) (result *StatusInKeystores, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getKeyDistributionStatusForKeystoresOptions, "getKeyDistributionStatusForKeystoresOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getKeyDistributionStatusForKeystoresOptions, "getKeyDistributionStatusForKeystoresOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getKeyDistributionStatusForKeystoresOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}/status_in_keystores`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getKeyDistributionStatusForKeystoresOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "GetKeyDistributionStatusForKeystores")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getKeyDistributionStatusForKeystoresOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*getKeyDistributionStatusForKeystoresOptions.UKOVault))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStatusInKeystores)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateManagedKeyFromTemplate : Update a managed key to match the key template
// Update a managed key to match the latest version of the associated key template. It will install, activate, or
// deactivate the key on target keystores in the group defined by the key template.
func (uko *UkoV4) UpdateManagedKeyFromTemplate(updateManagedKeyFromTemplateOptions *UpdateManagedKeyFromTemplateOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	return uko.UpdateManagedKeyFromTemplateWithContext(context.Background(), updateManagedKeyFromTemplateOptions)
}

// UpdateManagedKeyFromTemplateWithContext is an alternate form of the UpdateManagedKeyFromTemplate method which supports a Context parameter
func (uko *UkoV4) UpdateManagedKeyFromTemplateWithContext(ctx context.Context, updateManagedKeyFromTemplateOptions *UpdateManagedKeyFromTemplateOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateManagedKeyFromTemplateOptions, "updateManagedKeyFromTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateManagedKeyFromTemplateOptions, "updateManagedKeyFromTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateManagedKeyFromTemplateOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}/update_from_template`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateManagedKeyFromTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "UpdateManagedKeyFromTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if updateManagedKeyFromTemplateOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*updateManagedKeyFromTemplateOptions.UKOVault))
	}
	if updateManagedKeyFromTemplateOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateManagedKeyFromTemplateOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ActivateManagedKey : Activate a managed key
// Activates a managed key and performs key installation or activation operations on keystores in the keystore group
// associated with the managed key.
func (uko *UkoV4) ActivateManagedKey(activateManagedKeyOptions *ActivateManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	return uko.ActivateManagedKeyWithContext(context.Background(), activateManagedKeyOptions)
}

// ActivateManagedKeyWithContext is an alternate form of the ActivateManagedKey method which supports a Context parameter
func (uko *UkoV4) ActivateManagedKeyWithContext(ctx context.Context, activateManagedKeyOptions *ActivateManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(activateManagedKeyOptions, "activateManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(activateManagedKeyOptions, "activateManagedKeyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *activateManagedKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}/activate`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range activateManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "ActivateManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if activateManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*activateManagedKeyOptions.UKOVault))
	}
	if activateManagedKeyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*activateManagedKeyOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeactivateManagedKey : Deactivate a managed key
// Deactivates a managed key and performs key deactivation operations on keystores in the keystore group associated with
// the managed key.
func (uko *UkoV4) DeactivateManagedKey(deactivateManagedKeyOptions *DeactivateManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	return uko.DeactivateManagedKeyWithContext(context.Background(), deactivateManagedKeyOptions)
}

// DeactivateManagedKeyWithContext is an alternate form of the DeactivateManagedKey method which supports a Context parameter
func (uko *UkoV4) DeactivateManagedKeyWithContext(ctx context.Context, deactivateManagedKeyOptions *DeactivateManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deactivateManagedKeyOptions, "deactivateManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deactivateManagedKeyOptions, "deactivateManagedKeyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deactivateManagedKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}/deactivate`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deactivateManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "DeactivateManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deactivateManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*deactivateManagedKeyOptions.UKOVault))
	}
	if deactivateManagedKeyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deactivateManagedKeyOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DestroyManagedKey : Destroy a managed key
// Destroys a managed key and performs key destruction operations on keystores in the keystore group associated with the
// managed key. This operation cannot be undone. The managed key must be in a 'deactivated' state.
func (uko *UkoV4) DestroyManagedKey(destroyManagedKeyOptions *DestroyManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	return uko.DestroyManagedKeyWithContext(context.Background(), destroyManagedKeyOptions)
}

// DestroyManagedKeyWithContext is an alternate form of the DestroyManagedKey method which supports a Context parameter
func (uko *UkoV4) DestroyManagedKeyWithContext(ctx context.Context, destroyManagedKeyOptions *DestroyManagedKeyOptions) (result *ManagedKey, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(destroyManagedKeyOptions, "destroyManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(destroyManagedKeyOptions, "destroyManagedKeyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *destroyManagedKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}/destroy`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range destroyManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "DestroyManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if destroyManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*destroyManagedKeyOptions.UKOVault))
	}
	if destroyManagedKeyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*destroyManagedKeyOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKey)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SyncManagedKey : Sync a managed key in keystores
// Perform the synchronisation operation on a managed key to align the states in the associated keystores.
func (uko *UkoV4) SyncManagedKey(syncManagedKeyOptions *SyncManagedKeyOptions) (result *StatusInKeystores, response *core.DetailedResponse, err error) {
	return uko.SyncManagedKeyWithContext(context.Background(), syncManagedKeyOptions)
}

// SyncManagedKeyWithContext is an alternate form of the SyncManagedKey method which supports a Context parameter
func (uko *UkoV4) SyncManagedKeyWithContext(ctx context.Context, syncManagedKeyOptions *SyncManagedKeyOptions) (result *StatusInKeystores, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(syncManagedKeyOptions, "syncManagedKeyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(syncManagedKeyOptions, "syncManagedKeyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *syncManagedKeyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/managed_keys/{id}/sync_status_in_keystores`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range syncManagedKeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "SyncManagedKey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if syncManagedKeyOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*syncManagedKeyOptions.UKOVault))
	}
	if syncManagedKeyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*syncManagedKeyOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStatusInKeystores)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListKeyTemplates : List key templates
// List all key templates in the instance.
func (uko *UkoV4) ListKeyTemplates(listKeyTemplatesOptions *ListKeyTemplatesOptions) (result *TemplateList, response *core.DetailedResponse, err error) {
	return uko.ListKeyTemplatesWithContext(context.Background(), listKeyTemplatesOptions)
}

// ListKeyTemplatesWithContext is an alternate form of the ListKeyTemplates method which supports a Context parameter
func (uko *UkoV4) ListKeyTemplatesWithContext(ctx context.Context, listKeyTemplatesOptions *ListKeyTemplatesOptions) (result *TemplateList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listKeyTemplatesOptions, "listKeyTemplatesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/templates`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listKeyTemplatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "ListKeyTemplates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listKeyTemplatesOptions.VaultID != nil {
		builder.AddQuery("vault.id", strings.Join(listKeyTemplatesOptions.VaultID, ","))
	}
	if listKeyTemplatesOptions.KeyAlgorithm != nil {
		builder.AddQuery("key.algorithm", fmt.Sprint(*listKeyTemplatesOptions.KeyAlgorithm))
	}
	if listKeyTemplatesOptions.Sort != nil {
		builder.AddQuery("sort", strings.Join(listKeyTemplatesOptions.Sort, ","))
	}
	if listKeyTemplatesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listKeyTemplatesOptions.Limit))
	}
	if listKeyTemplatesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listKeyTemplatesOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateKeyTemplate : Create a key template
// Creates a new key template. Key templates are used to combine information necessary when creating a key that allow
// easy subsequent key creation, without needing to specify any of its details.
func (uko *UkoV4) CreateKeyTemplate(createKeyTemplateOptions *CreateKeyTemplateOptions) (result *Template, response *core.DetailedResponse, err error) {
	return uko.CreateKeyTemplateWithContext(context.Background(), createKeyTemplateOptions)
}

// CreateKeyTemplateWithContext is an alternate form of the CreateKeyTemplate method which supports a Context parameter
func (uko *UkoV4) CreateKeyTemplateWithContext(ctx context.Context, createKeyTemplateOptions *CreateKeyTemplateOptions) (result *Template, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createKeyTemplateOptions, "createKeyTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createKeyTemplateOptions, "createKeyTemplateOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/templates`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createKeyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "CreateKeyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createKeyTemplateOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*createKeyTemplateOptions.UKOVault))
	}

	body := make(map[string]interface{})
	if createKeyTemplateOptions.Vault != nil {
		body["vault"] = createKeyTemplateOptions.Vault
	}
	if createKeyTemplateOptions.Name != nil {
		body["name"] = createKeyTemplateOptions.Name
	}
	if createKeyTemplateOptions.Key != nil {
		body["key"] = createKeyTemplateOptions.Key
	}
	if createKeyTemplateOptions.Keystores != nil {
		body["keystores"] = createKeyTemplateOptions.Keystores
	}
	if createKeyTemplateOptions.Description != nil {
		body["description"] = createKeyTemplateOptions.Description
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
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplate)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteKeyTemplate : Delete a template
// Deletes a key template from the vault. The key template must not have any managed keys associated with it for it to
// be eligible for deletion.
func (uko *UkoV4) DeleteKeyTemplate(deleteKeyTemplateOptions *DeleteKeyTemplateOptions) (response *core.DetailedResponse, err error) {
	return uko.DeleteKeyTemplateWithContext(context.Background(), deleteKeyTemplateOptions)
}

// DeleteKeyTemplateWithContext is an alternate form of the DeleteKeyTemplate method which supports a Context parameter
func (uko *UkoV4) DeleteKeyTemplateWithContext(ctx context.Context, deleteKeyTemplateOptions *DeleteKeyTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteKeyTemplateOptions, "deleteKeyTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteKeyTemplateOptions, "deleteKeyTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteKeyTemplateOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/templates/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteKeyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "DeleteKeyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteKeyTemplateOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*deleteKeyTemplateOptions.UKOVault))
	}
	if deleteKeyTemplateOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deleteKeyTemplateOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = uko.Service.Request(request, nil)

	return
}

// GetKeyTemplate : Retrieve a key template
// Retrieves a key template and its details by specifying the ID.
func (uko *UkoV4) GetKeyTemplate(getKeyTemplateOptions *GetKeyTemplateOptions) (result *Template, response *core.DetailedResponse, err error) {
	return uko.GetKeyTemplateWithContext(context.Background(), getKeyTemplateOptions)
}

// GetKeyTemplateWithContext is an alternate form of the GetKeyTemplate method which supports a Context parameter
func (uko *UkoV4) GetKeyTemplateWithContext(ctx context.Context, getKeyTemplateOptions *GetKeyTemplateOptions) (result *Template, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getKeyTemplateOptions, "getKeyTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getKeyTemplateOptions, "getKeyTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getKeyTemplateOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/templates/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getKeyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "GetKeyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getKeyTemplateOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*getKeyTemplateOptions.UKOVault))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplate)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateKeyTemplate : Update a key template
// Updates attributes of a key template.
func (uko *UkoV4) UpdateKeyTemplate(updateKeyTemplateOptions *UpdateKeyTemplateOptions) (result *Template, response *core.DetailedResponse, err error) {
	return uko.UpdateKeyTemplateWithContext(context.Background(), updateKeyTemplateOptions)
}

// UpdateKeyTemplateWithContext is an alternate form of the UpdateKeyTemplate method which supports a Context parameter
func (uko *UkoV4) UpdateKeyTemplateWithContext(ctx context.Context, updateKeyTemplateOptions *UpdateKeyTemplateOptions) (result *Template, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateKeyTemplateOptions, "updateKeyTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateKeyTemplateOptions, "updateKeyTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateKeyTemplateOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/templates/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateKeyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "UpdateKeyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateKeyTemplateOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*updateKeyTemplateOptions.UKOVault))
	}
	if updateKeyTemplateOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateKeyTemplateOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateKeyTemplateOptions.Keystores != nil {
		body["keystores"] = updateKeyTemplateOptions.Keystores
	}
	if updateKeyTemplateOptions.Description != nil {
		body["description"] = updateKeyTemplateOptions.Description
	}
	if updateKeyTemplateOptions.Key != nil {
		body["key"] = updateKeyTemplateOptions.Key
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
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplate)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListKeystores : List all target keystores
// List all target keystores in the instance.
func (uko *UkoV4) ListKeystores(listKeystoresOptions *ListKeystoresOptions) (result *KeystoreList, response *core.DetailedResponse, err error) {
	return uko.ListKeystoresWithContext(context.Background(), listKeystoresOptions)
}

// ListKeystoresWithContext is an alternate form of the ListKeystores method which supports a Context parameter
func (uko *UkoV4) ListKeystoresWithContext(ctx context.Context, listKeystoresOptions *ListKeystoresOptions) (result *KeystoreList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listKeystoresOptions, "listKeystoresOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/keystores`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listKeystoresOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "ListKeystores")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listKeystoresOptions.Type != nil {
		builder.AddQuery("type", strings.Join(listKeystoresOptions.Type, ","))
	}
	if listKeystoresOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listKeystoresOptions.Name))
	}
	if listKeystoresOptions.Description != nil {
		builder.AddQuery("description", fmt.Sprint(*listKeystoresOptions.Description))
	}
	if listKeystoresOptions.Group != nil {
		builder.AddQuery("group", fmt.Sprint(*listKeystoresOptions.Group))
	}
	if listKeystoresOptions.Groups != nil {
		builder.AddQuery("groups[]", fmt.Sprint(*listKeystoresOptions.Groups))
	}
	if listKeystoresOptions.VaultID != nil {
		builder.AddQuery("vault.id", strings.Join(listKeystoresOptions.VaultID, ","))
	}
	if listKeystoresOptions.Location != nil {
		builder.AddQuery("location", fmt.Sprint(*listKeystoresOptions.Location))
	}
	if listKeystoresOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listKeystoresOptions.Limit))
	}
	if listKeystoresOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listKeystoresOptions.Offset))
	}
	if listKeystoresOptions.Sort != nil {
		builder.AddQuery("sort", strings.Join(listKeystoresOptions.Sort, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKeystoreList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateKeystore : Create an internal keystore or a keystore connection
// Creates a new internal keystore or a connection to an external keystore of the requested type.  If the `dry_run`
// query parameter is used, then a new keystore is not created in the database,  only a test is performed to verify if
// the connection information is correct. It is possible to sort by the following parameters: name, created_at,
// updated_at, vault.id.
func (uko *UkoV4) CreateKeystore(createKeystoreOptions *CreateKeystoreOptions) (result KeystoreIntf, response *core.DetailedResponse, err error) {
	return uko.CreateKeystoreWithContext(context.Background(), createKeystoreOptions)
}

// CreateKeystoreWithContext is an alternate form of the CreateKeystore method which supports a Context parameter
func (uko *UkoV4) CreateKeystoreWithContext(ctx context.Context, createKeystoreOptions *CreateKeystoreOptions) (result KeystoreIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createKeystoreOptions, "createKeystoreOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createKeystoreOptions, "createKeystoreOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/keystores`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createKeystoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "CreateKeystore")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createKeystoreOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*createKeystoreOptions.UKOVault))
	}

	if createKeystoreOptions.DryRun != nil {
		builder.AddQuery("dry_run", fmt.Sprint(*createKeystoreOptions.DryRun))
	}

	_, err = builder.SetBodyContentJSON(createKeystoreOptions.KeystoreBody)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKeystore)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteKeystore : Delete an internal keystore or a connection to an external keystore
// Deletes an internal keystore or a connection to an external keystore (in that case, the keystore on the remote system
// is unchanged).
func (uko *UkoV4) DeleteKeystore(deleteKeystoreOptions *DeleteKeystoreOptions) (response *core.DetailedResponse, err error) {
	return uko.DeleteKeystoreWithContext(context.Background(), deleteKeystoreOptions)
}

// DeleteKeystoreWithContext is an alternate form of the DeleteKeystore method which supports a Context parameter
func (uko *UkoV4) DeleteKeystoreWithContext(ctx context.Context, deleteKeystoreOptions *DeleteKeystoreOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteKeystoreOptions, "deleteKeystoreOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteKeystoreOptions, "deleteKeystoreOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteKeystoreOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/keystores/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteKeystoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "DeleteKeystore")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteKeystoreOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*deleteKeystoreOptions.UKOVault))
	}
	if deleteKeystoreOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deleteKeystoreOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = uko.Service.Request(request, nil)

	return
}

// GetKeystore : Retrieve a target keystore
// Retrieves a target keystore (either an internal keystore or a keystore connection) and its details by specifying the
// ID.
func (uko *UkoV4) GetKeystore(getKeystoreOptions *GetKeystoreOptions) (result KeystoreIntf, response *core.DetailedResponse, err error) {
	return uko.GetKeystoreWithContext(context.Background(), getKeystoreOptions)
}

// GetKeystoreWithContext is an alternate form of the GetKeystore method which supports a Context parameter
func (uko *UkoV4) GetKeystoreWithContext(ctx context.Context, getKeystoreOptions *GetKeystoreOptions) (result KeystoreIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getKeystoreOptions, "getKeystoreOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getKeystoreOptions, "getKeystoreOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getKeystoreOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/keystores/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getKeystoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "GetKeystore")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getKeystoreOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*getKeystoreOptions.UKOVault))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKeystore)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateKeystore : Update an internal keystore or a keystore connection
// Updates attributes of an internal keystore or a keystore connection.
func (uko *UkoV4) UpdateKeystore(updateKeystoreOptions *UpdateKeystoreOptions) (result KeystoreIntf, response *core.DetailedResponse, err error) {
	return uko.UpdateKeystoreWithContext(context.Background(), updateKeystoreOptions)
}

// UpdateKeystoreWithContext is an alternate form of the UpdateKeystore method which supports a Context parameter
func (uko *UkoV4) UpdateKeystoreWithContext(ctx context.Context, updateKeystoreOptions *UpdateKeystoreOptions) (result KeystoreIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateKeystoreOptions, "updateKeystoreOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateKeystoreOptions, "updateKeystoreOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateKeystoreOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/keystores/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateKeystoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "UpdateKeystore")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateKeystoreOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*updateKeystoreOptions.UKOVault))
	}
	if updateKeystoreOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateKeystoreOptions.IfMatch))
	}

	_, err = builder.SetBodyContentJSON(updateKeystoreOptions.KeystoreBody)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKeystore)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListAssociatedResourcesForTargetKeystore : List associated resources for a target keystore
// You can use this endpoint to obtain a list of resources associated with all keys referencing this keystore.
func (uko *UkoV4) ListAssociatedResourcesForTargetKeystore(listAssociatedResourcesForTargetKeystoreOptions *ListAssociatedResourcesForTargetKeystoreOptions) (result *AssociatedResourceList, response *core.DetailedResponse, err error) {
	return uko.ListAssociatedResourcesForTargetKeystoreWithContext(context.Background(), listAssociatedResourcesForTargetKeystoreOptions)
}

// ListAssociatedResourcesForTargetKeystoreWithContext is an alternate form of the ListAssociatedResourcesForTargetKeystore method which supports a Context parameter
func (uko *UkoV4) ListAssociatedResourcesForTargetKeystoreWithContext(ctx context.Context, listAssociatedResourcesForTargetKeystoreOptions *ListAssociatedResourcesForTargetKeystoreOptions) (result *AssociatedResourceList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAssociatedResourcesForTargetKeystoreOptions, "listAssociatedResourcesForTargetKeystoreOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listAssociatedResourcesForTargetKeystoreOptions, "listAssociatedResourcesForTargetKeystoreOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *listAssociatedResourcesForTargetKeystoreOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/keystores/{id}/associated_resources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAssociatedResourcesForTargetKeystoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "ListAssociatedResourcesForTargetKeystore")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listAssociatedResourcesForTargetKeystoreOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*listAssociatedResourcesForTargetKeystoreOptions.UKOVault))
	}

	if listAssociatedResourcesForTargetKeystoreOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAssociatedResourcesForTargetKeystoreOptions.Limit))
	}
	if listAssociatedResourcesForTargetKeystoreOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listAssociatedResourcesForTargetKeystoreOptions.Offset))
	}
	if listAssociatedResourcesForTargetKeystoreOptions.Sort != nil {
		builder.AddQuery("sort", strings.Join(listAssociatedResourcesForTargetKeystoreOptions.Sort, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAssociatedResourceList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetKeystoreStatus : Retrieve target keystore status
// Retrieves status of a single target keystore (either a keystore connection or an internal keystore).
func (uko *UkoV4) GetKeystoreStatus(getKeystoreStatusOptions *GetKeystoreStatusOptions) (result *KeystoreStatus, response *core.DetailedResponse, err error) {
	return uko.GetKeystoreStatusWithContext(context.Background(), getKeystoreStatusOptions)
}

// GetKeystoreStatusWithContext is an alternate form of the GetKeystoreStatus method which supports a Context parameter
func (uko *UkoV4) GetKeystoreStatusWithContext(ctx context.Context, getKeystoreStatusOptions *GetKeystoreStatusOptions) (result *KeystoreStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getKeystoreStatusOptions, "getKeystoreStatusOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getKeystoreStatusOptions, "getKeystoreStatusOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getKeystoreStatusOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/keystores/{id}/status`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getKeystoreStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "GetKeystoreStatus")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getKeystoreStatusOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*getKeystoreStatusOptions.UKOVault))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKeystoreStatus)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListManagedKeysFromKeystore : List managed keys on the target keystore
// Lists all managed keys installed on the target keystore (either a keystore connection or an internal keystore).
func (uko *UkoV4) ListManagedKeysFromKeystore(listManagedKeysFromKeystoreOptions *ListManagedKeysFromKeystoreOptions) (result *ManagedKeyList, response *core.DetailedResponse, err error) {
	return uko.ListManagedKeysFromKeystoreWithContext(context.Background(), listManagedKeysFromKeystoreOptions)
}

// ListManagedKeysFromKeystoreWithContext is an alternate form of the ListManagedKeysFromKeystore method which supports a Context parameter
func (uko *UkoV4) ListManagedKeysFromKeystoreWithContext(ctx context.Context, listManagedKeysFromKeystoreOptions *ListManagedKeysFromKeystoreOptions) (result *ManagedKeyList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listManagedKeysFromKeystoreOptions, "listManagedKeysFromKeystoreOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listManagedKeysFromKeystoreOptions, "listManagedKeysFromKeystoreOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *listManagedKeysFromKeystoreOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/keystores/{id}/managed_keys`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listManagedKeysFromKeystoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "ListManagedKeysFromKeystore")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listManagedKeysFromKeystoreOptions.UKOVault != nil {
		builder.AddHeader("UKO-Vault", fmt.Sprint(*listManagedKeysFromKeystoreOptions.UKOVault))
	}

	if listManagedKeysFromKeystoreOptions.Algorithm != nil {
		builder.AddQuery("algorithm", strings.Join(listManagedKeysFromKeystoreOptions.Algorithm, ","))
	}
	if listManagedKeysFromKeystoreOptions.State != nil {
		builder.AddQuery("state", strings.Join(listManagedKeysFromKeystoreOptions.State, ","))
	}
	if listManagedKeysFromKeystoreOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listManagedKeysFromKeystoreOptions.Limit))
	}
	if listManagedKeysFromKeystoreOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listManagedKeysFromKeystoreOptions.Offset))
	}
	if listManagedKeysFromKeystoreOptions.Sort != nil {
		builder.AddQuery("sort", strings.Join(listManagedKeysFromKeystoreOptions.Sort, ","))
	}
	if listManagedKeysFromKeystoreOptions.Label != nil {
		builder.AddQuery("label", fmt.Sprint(*listManagedKeysFromKeystoreOptions.Label))
	}
	if listManagedKeysFromKeystoreOptions.ActivationDate != nil {
		builder.AddQuery("activation_date", fmt.Sprint(*listManagedKeysFromKeystoreOptions.ActivationDate))
	}
	if listManagedKeysFromKeystoreOptions.ActivationDateMin != nil {
		builder.AddQuery("activation_date_min", fmt.Sprint(*listManagedKeysFromKeystoreOptions.ActivationDateMin))
	}
	if listManagedKeysFromKeystoreOptions.ActivationDateMax != nil {
		builder.AddQuery("activation_date_max", fmt.Sprint(*listManagedKeysFromKeystoreOptions.ActivationDateMax))
	}
	if listManagedKeysFromKeystoreOptions.DeactivationDate != nil {
		builder.AddQuery("deactivation_date", fmt.Sprint(*listManagedKeysFromKeystoreOptions.DeactivationDate))
	}
	if listManagedKeysFromKeystoreOptions.DeactivationDateMin != nil {
		builder.AddQuery("deactivation_date_min", fmt.Sprint(*listManagedKeysFromKeystoreOptions.DeactivationDateMin))
	}
	if listManagedKeysFromKeystoreOptions.DeactivationDateMax != nil {
		builder.AddQuery("deactivation_date_max", fmt.Sprint(*listManagedKeysFromKeystoreOptions.DeactivationDateMax))
	}
	if listManagedKeysFromKeystoreOptions.ExpirationDate != nil {
		builder.AddQuery("expiration_date", fmt.Sprint(*listManagedKeysFromKeystoreOptions.ExpirationDate))
	}
	if listManagedKeysFromKeystoreOptions.ExpirationDateMin != nil {
		builder.AddQuery("expiration_date_min", fmt.Sprint(*listManagedKeysFromKeystoreOptions.ExpirationDateMin))
	}
	if listManagedKeysFromKeystoreOptions.ExpirationDateMax != nil {
		builder.AddQuery("expiration_date_max", fmt.Sprint(*listManagedKeysFromKeystoreOptions.ExpirationDateMax))
	}
	if listManagedKeysFromKeystoreOptions.CreatedAt != nil {
		builder.AddQuery("created_at", fmt.Sprint(*listManagedKeysFromKeystoreOptions.CreatedAt))
	}
	if listManagedKeysFromKeystoreOptions.CreatedAtMin != nil {
		builder.AddQuery("created_at_min", fmt.Sprint(*listManagedKeysFromKeystoreOptions.CreatedAtMin))
	}
	if listManagedKeysFromKeystoreOptions.CreatedAtMax != nil {
		builder.AddQuery("created_at_max", fmt.Sprint(*listManagedKeysFromKeystoreOptions.CreatedAtMax))
	}
	if listManagedKeysFromKeystoreOptions.UpdatedAt != nil {
		builder.AddQuery("updated_at", fmt.Sprint(*listManagedKeysFromKeystoreOptions.UpdatedAt))
	}
	if listManagedKeysFromKeystoreOptions.UpdatedAtMin != nil {
		builder.AddQuery("updated_at_min", fmt.Sprint(*listManagedKeysFromKeystoreOptions.UpdatedAtMin))
	}
	if listManagedKeysFromKeystoreOptions.UpdatedAtMax != nil {
		builder.AddQuery("updated_at_max", fmt.Sprint(*listManagedKeysFromKeystoreOptions.UpdatedAtMax))
	}
	if listManagedKeysFromKeystoreOptions.Size != nil {
		builder.AddQuery("size", fmt.Sprint(*listManagedKeysFromKeystoreOptions.Size))
	}
	if listManagedKeysFromKeystoreOptions.SizeMin != nil {
		builder.AddQuery("size_min", fmt.Sprint(*listManagedKeysFromKeystoreOptions.SizeMin))
	}
	if listManagedKeysFromKeystoreOptions.SizeMax != nil {
		builder.AddQuery("size_max", fmt.Sprint(*listManagedKeysFromKeystoreOptions.SizeMax))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalManagedKeyList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListVaults : List all vaults
// List all vaults in the instance.
func (uko *UkoV4) ListVaults(listVaultsOptions *ListVaultsOptions) (result *VaultList, response *core.DetailedResponse, err error) {
	return uko.ListVaultsWithContext(context.Background(), listVaultsOptions)
}

// ListVaultsWithContext is an alternate form of the ListVaults method which supports a Context parameter
func (uko *UkoV4) ListVaultsWithContext(ctx context.Context, listVaultsOptions *ListVaultsOptions) (result *VaultList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listVaultsOptions, "listVaultsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/vaults`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listVaultsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "ListVaults")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listVaultsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listVaultsOptions.Limit))
	}
	if listVaultsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listVaultsOptions.Offset))
	}
	if listVaultsOptions.Sort != nil {
		builder.AddQuery("sort", strings.Join(listVaultsOptions.Sort, ","))
	}
	if listVaultsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listVaultsOptions.Name))
	}
	if listVaultsOptions.Description != nil {
		builder.AddQuery("description", fmt.Sprint(*listVaultsOptions.Description))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVaultList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateVault : Create a vault
// Creates a new vault in the instance with the specified name and description.
func (uko *UkoV4) CreateVault(createVaultOptions *CreateVaultOptions) (result *Vault, response *core.DetailedResponse, err error) {
	return uko.CreateVaultWithContext(context.Background(), createVaultOptions)
}

// CreateVaultWithContext is an alternate form of the CreateVault method which supports a Context parameter
func (uko *UkoV4) CreateVaultWithContext(ctx context.Context, createVaultOptions *CreateVaultOptions) (result *Vault, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createVaultOptions, "createVaultOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createVaultOptions, "createVaultOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/vaults`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createVaultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "CreateVault")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createVaultOptions.Name != nil {
		body["name"] = createVaultOptions.Name
	}
	if createVaultOptions.Description != nil {
		body["description"] = createVaultOptions.Description
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
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVault)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteVault : Delete an existing vault
// Deletes an existing vault from the system. A vault must be empty (that is, no managed keys or keystores remain in the
// vault) before the vault can be deleted.
func (uko *UkoV4) DeleteVault(deleteVaultOptions *DeleteVaultOptions) (response *core.DetailedResponse, err error) {
	return uko.DeleteVaultWithContext(context.Background(), deleteVaultOptions)
}

// DeleteVaultWithContext is an alternate form of the DeleteVault method which supports a Context parameter
func (uko *UkoV4) DeleteVaultWithContext(ctx context.Context, deleteVaultOptions *DeleteVaultOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteVaultOptions, "deleteVaultOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteVaultOptions, "deleteVaultOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteVaultOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/vaults/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteVaultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "DeleteVault")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteVaultOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deleteVaultOptions.IfMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = uko.Service.Request(request, nil)

	return
}

// GetVault : Retrieve a vault
// Retrieves a vault and its details by specifying the ID.
func (uko *UkoV4) GetVault(getVaultOptions *GetVaultOptions) (result *Vault, response *core.DetailedResponse, err error) {
	return uko.GetVaultWithContext(context.Background(), getVaultOptions)
}

// GetVaultWithContext is an alternate form of the GetVault method which supports a Context parameter
func (uko *UkoV4) GetVaultWithContext(ctx context.Context, getVaultOptions *GetVaultOptions) (result *Vault, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getVaultOptions, "getVaultOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getVaultOptions, "getVaultOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getVaultOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/vaults/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getVaultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "GetVault")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVault)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateVault : Update a vault
// Updates attributes of a vault.
func (uko *UkoV4) UpdateVault(updateVaultOptions *UpdateVaultOptions) (result *Vault, response *core.DetailedResponse, err error) {
	return uko.UpdateVaultWithContext(context.Background(), updateVaultOptions)
}

// UpdateVaultWithContext is an alternate form of the UpdateVault method which supports a Context parameter
func (uko *UkoV4) UpdateVaultWithContext(ctx context.Context, updateVaultOptions *UpdateVaultOptions) (result *Vault, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateVaultOptions, "updateVaultOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateVaultOptions, "updateVaultOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateVaultOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = uko.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(uko.Service.Options.URL, `/api/v4/vaults/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateVaultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("uko", "V4", "UpdateVault")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateVaultOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateVaultOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateVaultOptions.Name != nil {
		body["name"] = updateVaultOptions.Name
	}
	if updateVaultOptions.Description != nil {
		body["description"] = updateVaultOptions.Description
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
	response, err = uko.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVault)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ActivateManagedKeyOptions : The ActivateManagedKey options.
type ActivateManagedKeyOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewActivateManagedKeyOptions : Instantiate ActivateManagedKeyOptions
func (*UkoV4) NewActivateManagedKeyOptions(id string, uKOVault string, ifMatch string) *ActivateManagedKeyOptions {
	return &ActivateManagedKeyOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *ActivateManagedKeyOptions) SetID(id string) *ActivateManagedKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *ActivateManagedKeyOptions) SetUKOVault(uKOVault string) *ActivateManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ActivateManagedKeyOptions) SetIfMatch(ifMatch string) *ActivateManagedKeyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ActivateManagedKeyOptions) SetHeaders(param map[string]string) *ActivateManagedKeyOptions {
	options.Headers = param
	return options
}

// CreateKeyTemplateOptions : The CreateKeyTemplate options.
type CreateKeyTemplateOptions struct {
	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// ID of the Vault where the entity is to be created in.
	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Name of the template, it will be referenced when creating managed keys.
	Name *string `json:"name" validate:"required"`

	// Properties describing the properties of the managed key.
	Key *KeyProperties `json:"key" validate:"required"`

	// An array describing the type and group of target keystores the managed key is to be installed in.
	Keystores []KeystoresPropertiesCreateIntf `json:"keystores" validate:"required"`

	// Description of the key template.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateKeyTemplateOptions : Instantiate CreateKeyTemplateOptions
func (*UkoV4) NewCreateKeyTemplateOptions(uKOVault string, vault *VaultReferenceInCreationRequest, name string, key *KeyProperties, keystores []KeystoresPropertiesCreateIntf) *CreateKeyTemplateOptions {
	return &CreateKeyTemplateOptions{
		UKOVault: core.StringPtr(uKOVault),
		Vault: vault,
		Name: core.StringPtr(name),
		Key: key,
		Keystores: keystores,
	}
}

// SetUKOVault : Allow user to set UKOVault
func (_options *CreateKeyTemplateOptions) SetUKOVault(uKOVault string) *CreateKeyTemplateOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetVault : Allow user to set Vault
func (_options *CreateKeyTemplateOptions) SetVault(vault *VaultReferenceInCreationRequest) *CreateKeyTemplateOptions {
	_options.Vault = vault
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateKeyTemplateOptions) SetName(name string) *CreateKeyTemplateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetKey : Allow user to set Key
func (_options *CreateKeyTemplateOptions) SetKey(key *KeyProperties) *CreateKeyTemplateOptions {
	_options.Key = key
	return _options
}

// SetKeystores : Allow user to set Keystores
func (_options *CreateKeyTemplateOptions) SetKeystores(keystores []KeystoresPropertiesCreateIntf) *CreateKeyTemplateOptions {
	_options.Keystores = keystores
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateKeyTemplateOptions) SetDescription(description string) *CreateKeyTemplateOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateKeyTemplateOptions) SetHeaders(param map[string]string) *CreateKeyTemplateOptions {
	options.Headers = param
	return options
}

// CreateKeystoreOptions : The CreateKeystore options.
type CreateKeystoreOptions struct {
	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Keystore properties to update.
	KeystoreBody KeystoreCreationRequestIntf `json:"keystoreBody" validate:"required"`

	// Do not create a keystore, only verify if keystore created with given parameters can be communicated with
	// successfully.
	DryRun *bool `json:"dry_run,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateKeystoreOptions : Instantiate CreateKeystoreOptions
func (*UkoV4) NewCreateKeystoreOptions(uKOVault string, keystoreBody KeystoreCreationRequestIntf) *CreateKeystoreOptions {
	return &CreateKeystoreOptions{
		UKOVault: core.StringPtr(uKOVault),
		KeystoreBody: keystoreBody,
	}
}

// SetUKOVault : Allow user to set UKOVault
func (_options *CreateKeystoreOptions) SetUKOVault(uKOVault string) *CreateKeystoreOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetKeystoreBody : Allow user to set KeystoreBody
func (_options *CreateKeystoreOptions) SetKeystoreBody(keystoreBody KeystoreCreationRequestIntf) *CreateKeystoreOptions {
	_options.KeystoreBody = keystoreBody
	return _options
}

// SetDryRun : Allow user to set DryRun
func (_options *CreateKeystoreOptions) SetDryRun(dryRun bool) *CreateKeystoreOptions {
	_options.DryRun = core.BoolPtr(dryRun)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateKeystoreOptions) SetHeaders(param map[string]string) *CreateKeystoreOptions {
	options.Headers = param
	return options
}

// CreateManagedKeyOptions : The CreateManagedKey options.
type CreateManagedKeyOptions struct {
	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Name of the key template to use when creating a key.
	TemplateName *string `json:"template_name" validate:"required"`

	// ID of the Vault where the entity is to be created in.
	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// The label of the key.
	Label *string `json:"label" validate:"required"`

	// Key-value pairs associated with the key.
	Tags []Tag `json:"tags,omitempty"`

	// Description of the managed key.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateManagedKeyOptions : Instantiate CreateManagedKeyOptions
func (*UkoV4) NewCreateManagedKeyOptions(uKOVault string, templateName string, vault *VaultReferenceInCreationRequest, label string) *CreateManagedKeyOptions {
	return &CreateManagedKeyOptions{
		UKOVault: core.StringPtr(uKOVault),
		TemplateName: core.StringPtr(templateName),
		Vault: vault,
		Label: core.StringPtr(label),
	}
}

// SetUKOVault : Allow user to set UKOVault
func (_options *CreateManagedKeyOptions) SetUKOVault(uKOVault string) *CreateManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetTemplateName : Allow user to set TemplateName
func (_options *CreateManagedKeyOptions) SetTemplateName(templateName string) *CreateManagedKeyOptions {
	_options.TemplateName = core.StringPtr(templateName)
	return _options
}

// SetVault : Allow user to set Vault
func (_options *CreateManagedKeyOptions) SetVault(vault *VaultReferenceInCreationRequest) *CreateManagedKeyOptions {
	_options.Vault = vault
	return _options
}

// SetLabel : Allow user to set Label
func (_options *CreateManagedKeyOptions) SetLabel(label string) *CreateManagedKeyOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateManagedKeyOptions) SetTags(tags []Tag) *CreateManagedKeyOptions {
	_options.Tags = tags
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateManagedKeyOptions) SetDescription(description string) *CreateManagedKeyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateManagedKeyOptions) SetHeaders(param map[string]string) *CreateManagedKeyOptions {
	options.Headers = param
	return options
}

// CreateVaultOptions : The CreateVault options.
type CreateVaultOptions struct {
	// A human-readable name to assign to your vault. To protect your privacy, do not use personal data, such as your name
	// or location.
	Name *string `json:"name" validate:"required"`

	// Description of the vault.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateVaultOptions : Instantiate CreateVaultOptions
func (*UkoV4) NewCreateVaultOptions(name string) *CreateVaultOptions {
	return &CreateVaultOptions{
		Name: core.StringPtr(name),
	}
}

// SetName : Allow user to set Name
func (_options *CreateVaultOptions) SetName(name string) *CreateVaultOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateVaultOptions) SetDescription(description string) *CreateVaultOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateVaultOptions) SetHeaders(param map[string]string) *CreateVaultOptions {
	options.Headers = param
	return options
}

// DeactivateManagedKeyOptions : The DeactivateManagedKey options.
type DeactivateManagedKeyOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeactivateManagedKeyOptions : Instantiate DeactivateManagedKeyOptions
func (*UkoV4) NewDeactivateManagedKeyOptions(id string, uKOVault string, ifMatch string) *DeactivateManagedKeyOptions {
	return &DeactivateManagedKeyOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *DeactivateManagedKeyOptions) SetID(id string) *DeactivateManagedKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *DeactivateManagedKeyOptions) SetUKOVault(uKOVault string) *DeactivateManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeactivateManagedKeyOptions) SetIfMatch(ifMatch string) *DeactivateManagedKeyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeactivateManagedKeyOptions) SetHeaders(param map[string]string) *DeactivateManagedKeyOptions {
	options.Headers = param
	return options
}

// DeleteKeyTemplateOptions : The DeleteKeyTemplate options.
type DeleteKeyTemplateOptions struct {
	// UUID of the template.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteKeyTemplateOptions : Instantiate DeleteKeyTemplateOptions
func (*UkoV4) NewDeleteKeyTemplateOptions(id string, uKOVault string, ifMatch string) *DeleteKeyTemplateOptions {
	return &DeleteKeyTemplateOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteKeyTemplateOptions) SetID(id string) *DeleteKeyTemplateOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *DeleteKeyTemplateOptions) SetUKOVault(uKOVault string) *DeleteKeyTemplateOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeleteKeyTemplateOptions) SetIfMatch(ifMatch string) *DeleteKeyTemplateOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteKeyTemplateOptions) SetHeaders(param map[string]string) *DeleteKeyTemplateOptions {
	options.Headers = param
	return options
}

// DeleteKeystoreOptions : The DeleteKeystore options.
type DeleteKeystoreOptions struct {
	// UUID of the keystore.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteKeystoreOptions : Instantiate DeleteKeystoreOptions
func (*UkoV4) NewDeleteKeystoreOptions(id string, uKOVault string, ifMatch string) *DeleteKeystoreOptions {
	return &DeleteKeystoreOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteKeystoreOptions) SetID(id string) *DeleteKeystoreOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *DeleteKeystoreOptions) SetUKOVault(uKOVault string) *DeleteKeystoreOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeleteKeystoreOptions) SetIfMatch(ifMatch string) *DeleteKeystoreOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteKeystoreOptions) SetHeaders(param map[string]string) *DeleteKeystoreOptions {
	options.Headers = param
	return options
}

// DeleteManagedKeyOptions : The DeleteManagedKey options.
type DeleteManagedKeyOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteManagedKeyOptions : Instantiate DeleteManagedKeyOptions
func (*UkoV4) NewDeleteManagedKeyOptions(id string, uKOVault string, ifMatch string) *DeleteManagedKeyOptions {
	return &DeleteManagedKeyOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteManagedKeyOptions) SetID(id string) *DeleteManagedKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *DeleteManagedKeyOptions) SetUKOVault(uKOVault string) *DeleteManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeleteManagedKeyOptions) SetIfMatch(ifMatch string) *DeleteManagedKeyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteManagedKeyOptions) SetHeaders(param map[string]string) *DeleteManagedKeyOptions {
	options.Headers = param
	return options
}

// DeleteVaultOptions : The DeleteVault options.
type DeleteVaultOptions struct {
	// UUID of the vault.
	ID *string `json:"id" validate:"required,ne="`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteVaultOptions : Instantiate DeleteVaultOptions
func (*UkoV4) NewDeleteVaultOptions(id string, ifMatch string) *DeleteVaultOptions {
	return &DeleteVaultOptions{
		ID: core.StringPtr(id),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteVaultOptions) SetID(id string) *DeleteVaultOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeleteVaultOptions) SetIfMatch(ifMatch string) *DeleteVaultOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteVaultOptions) SetHeaders(param map[string]string) *DeleteVaultOptions {
	options.Headers = param
	return options
}

// DestroyManagedKeyOptions : The DestroyManagedKey options.
type DestroyManagedKeyOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDestroyManagedKeyOptions : Instantiate DestroyManagedKeyOptions
func (*UkoV4) NewDestroyManagedKeyOptions(id string, uKOVault string, ifMatch string) *DestroyManagedKeyOptions {
	return &DestroyManagedKeyOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *DestroyManagedKeyOptions) SetID(id string) *DestroyManagedKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *DestroyManagedKeyOptions) SetUKOVault(uKOVault string) *DestroyManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DestroyManagedKeyOptions) SetIfMatch(ifMatch string) *DestroyManagedKeyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DestroyManagedKeyOptions) SetHeaders(param map[string]string) *DestroyManagedKeyOptions {
	options.Headers = param
	return options
}

// GetKeyDistributionStatusForKeystoresOptions : The GetKeyDistributionStatusForKeystores options.
type GetKeyDistributionStatusForKeystoresOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetKeyDistributionStatusForKeystoresOptions : Instantiate GetKeyDistributionStatusForKeystoresOptions
func (*UkoV4) NewGetKeyDistributionStatusForKeystoresOptions(id string, uKOVault string) *GetKeyDistributionStatusForKeystoresOptions {
	return &GetKeyDistributionStatusForKeystoresOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
	}
}

// SetID : Allow user to set ID
func (_options *GetKeyDistributionStatusForKeystoresOptions) SetID(id string) *GetKeyDistributionStatusForKeystoresOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *GetKeyDistributionStatusForKeystoresOptions) SetUKOVault(uKOVault string) *GetKeyDistributionStatusForKeystoresOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetKeyDistributionStatusForKeystoresOptions) SetHeaders(param map[string]string) *GetKeyDistributionStatusForKeystoresOptions {
	options.Headers = param
	return options
}

// GetKeyTemplateOptions : The GetKeyTemplate options.
type GetKeyTemplateOptions struct {
	// UUID of the template.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetKeyTemplateOptions : Instantiate GetKeyTemplateOptions
func (*UkoV4) NewGetKeyTemplateOptions(id string, uKOVault string) *GetKeyTemplateOptions {
	return &GetKeyTemplateOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
	}
}

// SetID : Allow user to set ID
func (_options *GetKeyTemplateOptions) SetID(id string) *GetKeyTemplateOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *GetKeyTemplateOptions) SetUKOVault(uKOVault string) *GetKeyTemplateOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetKeyTemplateOptions) SetHeaders(param map[string]string) *GetKeyTemplateOptions {
	options.Headers = param
	return options
}

// GetKeystoreOptions : The GetKeystore options.
type GetKeystoreOptions struct {
	// UUID of the keystore.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetKeystoreOptions : Instantiate GetKeystoreOptions
func (*UkoV4) NewGetKeystoreOptions(id string, uKOVault string) *GetKeystoreOptions {
	return &GetKeystoreOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
	}
}

// SetID : Allow user to set ID
func (_options *GetKeystoreOptions) SetID(id string) *GetKeystoreOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *GetKeystoreOptions) SetUKOVault(uKOVault string) *GetKeystoreOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetKeystoreOptions) SetHeaders(param map[string]string) *GetKeystoreOptions {
	options.Headers = param
	return options
}

// GetKeystoreStatusOptions : The GetKeystoreStatus options.
type GetKeystoreStatusOptions struct {
	// UUID of the keystore.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetKeystoreStatusOptions : Instantiate GetKeystoreStatusOptions
func (*UkoV4) NewGetKeystoreStatusOptions(id string, uKOVault string) *GetKeystoreStatusOptions {
	return &GetKeystoreStatusOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
	}
}

// SetID : Allow user to set ID
func (_options *GetKeystoreStatusOptions) SetID(id string) *GetKeystoreStatusOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *GetKeystoreStatusOptions) SetUKOVault(uKOVault string) *GetKeystoreStatusOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetKeystoreStatusOptions) SetHeaders(param map[string]string) *GetKeystoreStatusOptions {
	options.Headers = param
	return options
}

// GetManagedKeyOptions : The GetManagedKey options.
type GetManagedKeyOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetManagedKeyOptions : Instantiate GetManagedKeyOptions
func (*UkoV4) NewGetManagedKeyOptions(id string, uKOVault string) *GetManagedKeyOptions {
	return &GetManagedKeyOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
	}
}

// SetID : Allow user to set ID
func (_options *GetManagedKeyOptions) SetID(id string) *GetManagedKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *GetManagedKeyOptions) SetUKOVault(uKOVault string) *GetManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetManagedKeyOptions) SetHeaders(param map[string]string) *GetManagedKeyOptions {
	options.Headers = param
	return options
}

// GetVaultOptions : The GetVault options.
type GetVaultOptions struct {
	// UUID of the vault.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetVaultOptions : Instantiate GetVaultOptions
func (*UkoV4) NewGetVaultOptions(id string) *GetVaultOptions {
	return &GetVaultOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetVaultOptions) SetID(id string) *GetVaultOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetVaultOptions) SetHeaders(param map[string]string) *GetVaultOptions {
	options.Headers = param
	return options
}

// ListAssociatedResourcesForManagedKeyOptions : The ListAssociatedResourcesForManagedKey options.
type ListAssociatedResourcesForManagedKeyOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// The number of resources to retrieve.
	Limit *int64 `json:"limit,omitempty"`

	// The number of resources to skip.
	Offset *int64 `json:"offset,omitempty"`

	// Define sorting order.
	Sort []string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAssociatedResourcesForManagedKeyOptions : Instantiate ListAssociatedResourcesForManagedKeyOptions
func (*UkoV4) NewListAssociatedResourcesForManagedKeyOptions(id string, uKOVault string) *ListAssociatedResourcesForManagedKeyOptions {
	return &ListAssociatedResourcesForManagedKeyOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
	}
}

// SetID : Allow user to set ID
func (_options *ListAssociatedResourcesForManagedKeyOptions) SetID(id string) *ListAssociatedResourcesForManagedKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *ListAssociatedResourcesForManagedKeyOptions) SetUKOVault(uKOVault string) *ListAssociatedResourcesForManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAssociatedResourcesForManagedKeyOptions) SetLimit(limit int64) *ListAssociatedResourcesForManagedKeyOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListAssociatedResourcesForManagedKeyOptions) SetOffset(offset int64) *ListAssociatedResourcesForManagedKeyOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListAssociatedResourcesForManagedKeyOptions) SetSort(sort []string) *ListAssociatedResourcesForManagedKeyOptions {
	_options.Sort = sort
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAssociatedResourcesForManagedKeyOptions) SetHeaders(param map[string]string) *ListAssociatedResourcesForManagedKeyOptions {
	options.Headers = param
	return options
}

// ListAssociatedResourcesForTargetKeystoreOptions : The ListAssociatedResourcesForTargetKeystore options.
type ListAssociatedResourcesForTargetKeystoreOptions struct {
	// UUID of the keystore.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// The number of resources to retrieve.
	Limit *int64 `json:"limit,omitempty"`

	// The number of resources to skip.
	Offset *int64 `json:"offset,omitempty"`

	// Define sorting order.
	Sort []string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAssociatedResourcesForTargetKeystoreOptions : Instantiate ListAssociatedResourcesForTargetKeystoreOptions
func (*UkoV4) NewListAssociatedResourcesForTargetKeystoreOptions(id string, uKOVault string) *ListAssociatedResourcesForTargetKeystoreOptions {
	return &ListAssociatedResourcesForTargetKeystoreOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
	}
}

// SetID : Allow user to set ID
func (_options *ListAssociatedResourcesForTargetKeystoreOptions) SetID(id string) *ListAssociatedResourcesForTargetKeystoreOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *ListAssociatedResourcesForTargetKeystoreOptions) SetUKOVault(uKOVault string) *ListAssociatedResourcesForTargetKeystoreOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAssociatedResourcesForTargetKeystoreOptions) SetLimit(limit int64) *ListAssociatedResourcesForTargetKeystoreOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListAssociatedResourcesForTargetKeystoreOptions) SetOffset(offset int64) *ListAssociatedResourcesForTargetKeystoreOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListAssociatedResourcesForTargetKeystoreOptions) SetSort(sort []string) *ListAssociatedResourcesForTargetKeystoreOptions {
	_options.Sort = sort
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAssociatedResourcesForTargetKeystoreOptions) SetHeaders(param map[string]string) *ListAssociatedResourcesForTargetKeystoreOptions {
	options.Headers = param
	return options
}

// ListKeyTemplatesOptions : The ListKeyTemplates options.
type ListKeyTemplatesOptions struct {
	// The UUID of the Vault.
	VaultID []string `json:"vault.id,omitempty"`

	// The algorithm of a returned key templates.
	KeyAlgorithm *string `json:"key.algorithm,omitempty"`

	// Define sorting order.
	Sort []string `json:"sort,omitempty"`

	// The number of resources to retrieve.
	Limit *int64 `json:"limit,omitempty"`

	// The number of resources to skip.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListKeyTemplatesOptions.KeyAlgorithm property.
// The algorithm of a returned key templates.
const (
	ListKeyTemplatesOptions_KeyAlgorithm_Aes = "aes"
	ListKeyTemplatesOptions_KeyAlgorithm_Ec = "ec"
	ListKeyTemplatesOptions_KeyAlgorithm_Hmac = "hmac"
	ListKeyTemplatesOptions_KeyAlgorithm_Rsa = "rsa"
)

// NewListKeyTemplatesOptions : Instantiate ListKeyTemplatesOptions
func (*UkoV4) NewListKeyTemplatesOptions() *ListKeyTemplatesOptions {
	return &ListKeyTemplatesOptions{}
}

// SetVaultID : Allow user to set VaultID
func (_options *ListKeyTemplatesOptions) SetVaultID(vaultID []string) *ListKeyTemplatesOptions {
	_options.VaultID = vaultID
	return _options
}

// SetKeyAlgorithm : Allow user to set KeyAlgorithm
func (_options *ListKeyTemplatesOptions) SetKeyAlgorithm(keyAlgorithm string) *ListKeyTemplatesOptions {
	_options.KeyAlgorithm = core.StringPtr(keyAlgorithm)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListKeyTemplatesOptions) SetSort(sort []string) *ListKeyTemplatesOptions {
	_options.Sort = sort
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListKeyTemplatesOptions) SetLimit(limit int64) *ListKeyTemplatesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListKeyTemplatesOptions) SetOffset(offset int64) *ListKeyTemplatesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListKeyTemplatesOptions) SetHeaders(param map[string]string) *ListKeyTemplatesOptions {
	options.Headers = param
	return options
}

// ListKeystoresOptions : The ListKeystores options.
type ListKeystoresOptions struct {
	// Keystore type.
	Type []string `json:"type,omitempty"`

	// Return only keystores whose name contains the string.
	Name *string `json:"name,omitempty"`

	// Return only keystores whose description contains the string.
	Description *string `json:"description,omitempty"`

	// A Keystore group. This query parameter cannot be used in conjunction with the 'groups[]' query parameter.
	Group *string `json:"group,omitempty"`

	// Keystore groups.
	Groups *string `json:"groups[],omitempty"`

	// The UUID of the Vault.
	VaultID []string `json:"vault.id,omitempty"`

	// Keystore location.
	Location *string `json:"location,omitempty"`

	// The number of resources to retrieve.
	Limit *int64 `json:"limit,omitempty"`

	// The number of resources to skip.
	Offset *int64 `json:"offset,omitempty"`

	// Define sorting order.
	Sort []string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListKeystoresOptions.Type property.
// Type of keystore.
const (
	ListKeystoresOptions_Type_AwsKms = "aws_kms"
	ListKeystoresOptions_Type_AzureKeyVault = "azure_key_vault"
	ListKeystoresOptions_Type_GoogleKms = "google_kms"
	ListKeystoresOptions_Type_IbmCloudKms = "ibm_cloud_kms"
)

// NewListKeystoresOptions : Instantiate ListKeystoresOptions
func (*UkoV4) NewListKeystoresOptions() *ListKeystoresOptions {
	return &ListKeystoresOptions{}
}

// SetType : Allow user to set Type
func (_options *ListKeystoresOptions) SetType(typeVar []string) *ListKeystoresOptions {
	_options.Type = typeVar
	return _options
}

// SetName : Allow user to set Name
func (_options *ListKeystoresOptions) SetName(name string) *ListKeystoresOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ListKeystoresOptions) SetDescription(description string) *ListKeystoresOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetGroup : Allow user to set Group
func (_options *ListKeystoresOptions) SetGroup(group string) *ListKeystoresOptions {
	_options.Group = core.StringPtr(group)
	return _options
}

// SetGroups : Allow user to set Groups
func (_options *ListKeystoresOptions) SetGroups(groups string) *ListKeystoresOptions {
	_options.Groups = core.StringPtr(groups)
	return _options
}

// SetVaultID : Allow user to set VaultID
func (_options *ListKeystoresOptions) SetVaultID(vaultID []string) *ListKeystoresOptions {
	_options.VaultID = vaultID
	return _options
}

// SetLocation : Allow user to set Location
func (_options *ListKeystoresOptions) SetLocation(location string) *ListKeystoresOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListKeystoresOptions) SetLimit(limit int64) *ListKeystoresOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListKeystoresOptions) SetOffset(offset int64) *ListKeystoresOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListKeystoresOptions) SetSort(sort []string) *ListKeystoresOptions {
	_options.Sort = sort
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListKeystoresOptions) SetHeaders(param map[string]string) *ListKeystoresOptions {
	options.Headers = param
	return options
}

// ListManagedKeysFromKeystoreOptions : The ListManagedKeysFromKeystore options.
type ListManagedKeysFromKeystoreOptions struct {
	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// UUID of the keystore.
	ID *string `json:"id" validate:"required,ne="`

	// The algorithm of a returned keys.
	Algorithm []string `json:"algorithm,omitempty"`

	// The state that returned keys are to be in.
	State []string `json:"state,omitempty"`

	// The number of resources to retrieve.
	Limit *int64 `json:"limit,omitempty"`

	// The number of resources to skip.
	Offset *int64 `json:"offset,omitempty"`

	// Define sorting order.
	Sort []string `json:"sort,omitempty"`

	// The label of the key.
	Label *string `json:"label,omitempty"`

	// Return only managed keys whose activation_date matches the parameter.
	ActivationDate *string `json:"activation_date,omitempty"`

	// Return only managed keys whose activation_date is at or after the parameter value. This query parameter cannot be
	// used in conjunction with the 'activation_date' query parameter.
	ActivationDateMin *string `json:"activation_date_min,omitempty"`

	// Return only managed keys whose activation_date is at or before the parameter value. This query parameter cannot be
	// used in conjunction with the 'activation_date' query parameter.
	ActivationDateMax *string `json:"activation_date_max,omitempty"`

	// Return only managed keys whose deactivation_date matches the parameter. This query parameter cannot be used in
	// conjunction with the 'expiration_date' query parameter.
	DeactivationDate *string `json:"deactivation_date,omitempty"`

	// Return only managed keys whose deactivation_date is at or after the parameter value. This query parameter cannot be
	// used in conjunction with the 'deactivation_date', 'expiration_date', 'expiration_date_min' and 'expiration_date_max'
	// query parameters.
	DeactivationDateMin *string `json:"deactivation_date_min,omitempty"`

	// Return only managed keys whose deactivation_date is at or before the parameter value. This query parameter cannot be
	// used in conjunction with the 'deactivation_date', 'expiration_date', 'expiration_date_min' and 'expiration_date_max'
	// query parameters.
	DeactivationDateMax *string `json:"deactivation_date_max,omitempty"`

	// Return only managed keys whose deactivation_date matches the parameter.
	ExpirationDate *string `json:"expiration_date,omitempty"`

	// Return only managed keys whose deactivation_date is at or after the parameter value. This query parameter cannot be
	// used in conjunction with the 'deactivation_date', 'expiration_date', 'deactivation_date_min' and
	// 'deactivation_date_max' query parameters.
	ExpirationDateMin *string `json:"expiration_date_min,omitempty"`

	// Return only managed keys whose deactivation_date is at or before the parameter value. This query parameter cannot be
	// used in conjunction with the 'deactivation_date', 'expiration_date', 'deactivation_date_min' and
	// 'deactivation_date_max' query parameters.
	ExpirationDateMax *string `json:"expiration_date_max,omitempty"`

	// Return only managed keys whose created_at matches the parameter.
	CreatedAt *string `json:"created_at,omitempty"`

	// Return only managed keys whose created_at is at or after the parameter value. This query parameter cannot be used in
	// conjunction with the 'created_at' query parameter.
	CreatedAtMin *string `json:"created_at_min,omitempty"`

	// Return only managed keys whose created_at is at or before the parameter value. This query parameter cannot be used
	// in conjunction with the 'created_at' query parameter.
	CreatedAtMax *string `json:"created_at_max,omitempty"`

	// Return only managed keys whose updated_at matches the parameter.
	UpdatedAt *string `json:"updated_at,omitempty"`

	// Return only managed keys whose updated_at is after the parameter value. This query parameter cannot be used in
	// conjunction with the 'updated_at' query parameter.
	UpdatedAtMin *string `json:"updated_at_min,omitempty"`

	// Return only managed keys whose updated_at is before the parameter value. This query parameter cannot be used in
	// conjunction with the 'updated_at' query parameter.
	UpdatedAtMax *string `json:"updated_at_max,omitempty"`

	// The size of the key.
	Size *int64 `json:"size,omitempty"`

	// The minimum size of the key. This query parameter cannot be used in conjunction with the 'size' query parameter.
	SizeMin *int64 `json:"size_min,omitempty"`

	// The maximum size of the key. This query parameter cannot be used in conjunction with the 'size' query parameter.
	SizeMax *int64 `json:"size_max,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListManagedKeysFromKeystoreOptions.Algorithm property.
// The algorithm of the key.
const (
	ListManagedKeysFromKeystoreOptions_Algorithm_Aes = "aes"
	ListManagedKeysFromKeystoreOptions_Algorithm_Ec = "ec"
	ListManagedKeysFromKeystoreOptions_Algorithm_Hmac = "hmac"
	ListManagedKeysFromKeystoreOptions_Algorithm_Rsa = "rsa"
)

// Constants associated with the ListManagedKeysFromKeystoreOptions.State property.
// The state of the key.
const (
	ListManagedKeysFromKeystoreOptions_State_Active = "active"
	ListManagedKeysFromKeystoreOptions_State_Deactivated = "deactivated"
	ListManagedKeysFromKeystoreOptions_State_Destroyed = "destroyed"
	ListManagedKeysFromKeystoreOptions_State_PreActivation = "pre_activation"
)

// NewListManagedKeysFromKeystoreOptions : Instantiate ListManagedKeysFromKeystoreOptions
func (*UkoV4) NewListManagedKeysFromKeystoreOptions(uKOVault string, id string) *ListManagedKeysFromKeystoreOptions {
	return &ListManagedKeysFromKeystoreOptions{
		UKOVault: core.StringPtr(uKOVault),
		ID: core.StringPtr(id),
	}
}

// SetUKOVault : Allow user to set UKOVault
func (_options *ListManagedKeysFromKeystoreOptions) SetUKOVault(uKOVault string) *ListManagedKeysFromKeystoreOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetID : Allow user to set ID
func (_options *ListManagedKeysFromKeystoreOptions) SetID(id string) *ListManagedKeysFromKeystoreOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAlgorithm : Allow user to set Algorithm
func (_options *ListManagedKeysFromKeystoreOptions) SetAlgorithm(algorithm []string) *ListManagedKeysFromKeystoreOptions {
	_options.Algorithm = algorithm
	return _options
}

// SetState : Allow user to set State
func (_options *ListManagedKeysFromKeystoreOptions) SetState(state []string) *ListManagedKeysFromKeystoreOptions {
	_options.State = state
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListManagedKeysFromKeystoreOptions) SetLimit(limit int64) *ListManagedKeysFromKeystoreOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListManagedKeysFromKeystoreOptions) SetOffset(offset int64) *ListManagedKeysFromKeystoreOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListManagedKeysFromKeystoreOptions) SetSort(sort []string) *ListManagedKeysFromKeystoreOptions {
	_options.Sort = sort
	return _options
}

// SetLabel : Allow user to set Label
func (_options *ListManagedKeysFromKeystoreOptions) SetLabel(label string) *ListManagedKeysFromKeystoreOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetActivationDate : Allow user to set ActivationDate
func (_options *ListManagedKeysFromKeystoreOptions) SetActivationDate(activationDate string) *ListManagedKeysFromKeystoreOptions {
	_options.ActivationDate = core.StringPtr(activationDate)
	return _options
}

// SetActivationDateMin : Allow user to set ActivationDateMin
func (_options *ListManagedKeysFromKeystoreOptions) SetActivationDateMin(activationDateMin string) *ListManagedKeysFromKeystoreOptions {
	_options.ActivationDateMin = core.StringPtr(activationDateMin)
	return _options
}

// SetActivationDateMax : Allow user to set ActivationDateMax
func (_options *ListManagedKeysFromKeystoreOptions) SetActivationDateMax(activationDateMax string) *ListManagedKeysFromKeystoreOptions {
	_options.ActivationDateMax = core.StringPtr(activationDateMax)
	return _options
}

// SetDeactivationDate : Allow user to set DeactivationDate
func (_options *ListManagedKeysFromKeystoreOptions) SetDeactivationDate(deactivationDate string) *ListManagedKeysFromKeystoreOptions {
	_options.DeactivationDate = core.StringPtr(deactivationDate)
	return _options
}

// SetDeactivationDateMin : Allow user to set DeactivationDateMin
func (_options *ListManagedKeysFromKeystoreOptions) SetDeactivationDateMin(deactivationDateMin string) *ListManagedKeysFromKeystoreOptions {
	_options.DeactivationDateMin = core.StringPtr(deactivationDateMin)
	return _options
}

// SetDeactivationDateMax : Allow user to set DeactivationDateMax
func (_options *ListManagedKeysFromKeystoreOptions) SetDeactivationDateMax(deactivationDateMax string) *ListManagedKeysFromKeystoreOptions {
	_options.DeactivationDateMax = core.StringPtr(deactivationDateMax)
	return _options
}

// SetExpirationDate : Allow user to set ExpirationDate
func (_options *ListManagedKeysFromKeystoreOptions) SetExpirationDate(expirationDate string) *ListManagedKeysFromKeystoreOptions {
	_options.ExpirationDate = core.StringPtr(expirationDate)
	return _options
}

// SetExpirationDateMin : Allow user to set ExpirationDateMin
func (_options *ListManagedKeysFromKeystoreOptions) SetExpirationDateMin(expirationDateMin string) *ListManagedKeysFromKeystoreOptions {
	_options.ExpirationDateMin = core.StringPtr(expirationDateMin)
	return _options
}

// SetExpirationDateMax : Allow user to set ExpirationDateMax
func (_options *ListManagedKeysFromKeystoreOptions) SetExpirationDateMax(expirationDateMax string) *ListManagedKeysFromKeystoreOptions {
	_options.ExpirationDateMax = core.StringPtr(expirationDateMax)
	return _options
}

// SetCreatedAt : Allow user to set CreatedAt
func (_options *ListManagedKeysFromKeystoreOptions) SetCreatedAt(createdAt string) *ListManagedKeysFromKeystoreOptions {
	_options.CreatedAt = core.StringPtr(createdAt)
	return _options
}

// SetCreatedAtMin : Allow user to set CreatedAtMin
func (_options *ListManagedKeysFromKeystoreOptions) SetCreatedAtMin(createdAtMin string) *ListManagedKeysFromKeystoreOptions {
	_options.CreatedAtMin = core.StringPtr(createdAtMin)
	return _options
}

// SetCreatedAtMax : Allow user to set CreatedAtMax
func (_options *ListManagedKeysFromKeystoreOptions) SetCreatedAtMax(createdAtMax string) *ListManagedKeysFromKeystoreOptions {
	_options.CreatedAtMax = core.StringPtr(createdAtMax)
	return _options
}

// SetUpdatedAt : Allow user to set UpdatedAt
func (_options *ListManagedKeysFromKeystoreOptions) SetUpdatedAt(updatedAt string) *ListManagedKeysFromKeystoreOptions {
	_options.UpdatedAt = core.StringPtr(updatedAt)
	return _options
}

// SetUpdatedAtMin : Allow user to set UpdatedAtMin
func (_options *ListManagedKeysFromKeystoreOptions) SetUpdatedAtMin(updatedAtMin string) *ListManagedKeysFromKeystoreOptions {
	_options.UpdatedAtMin = core.StringPtr(updatedAtMin)
	return _options
}

// SetUpdatedAtMax : Allow user to set UpdatedAtMax
func (_options *ListManagedKeysFromKeystoreOptions) SetUpdatedAtMax(updatedAtMax string) *ListManagedKeysFromKeystoreOptions {
	_options.UpdatedAtMax = core.StringPtr(updatedAtMax)
	return _options
}

// SetSize : Allow user to set Size
func (_options *ListManagedKeysFromKeystoreOptions) SetSize(size int64) *ListManagedKeysFromKeystoreOptions {
	_options.Size = core.Int64Ptr(size)
	return _options
}

// SetSizeMin : Allow user to set SizeMin
func (_options *ListManagedKeysFromKeystoreOptions) SetSizeMin(sizeMin int64) *ListManagedKeysFromKeystoreOptions {
	_options.SizeMin = core.Int64Ptr(sizeMin)
	return _options
}

// SetSizeMax : Allow user to set SizeMax
func (_options *ListManagedKeysFromKeystoreOptions) SetSizeMax(sizeMax int64) *ListManagedKeysFromKeystoreOptions {
	_options.SizeMax = core.Int64Ptr(sizeMax)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListManagedKeysFromKeystoreOptions) SetHeaders(param map[string]string) *ListManagedKeysFromKeystoreOptions {
	options.Headers = param
	return options
}

// ListManagedKeysOptions : The ListManagedKeys options.
type ListManagedKeysOptions struct {
	// The UUID of the Vault.
	VaultID []string `json:"vault.id,omitempty"`

	// The algorithm of a returned keys.
	Algorithm []string `json:"algorithm,omitempty"`

	// The state that returned keys are to be in.
	State []string `json:"state,omitempty"`

	// The number of resources to retrieve.
	Limit *int64 `json:"limit,omitempty"`

	// The number of resources to skip.
	Offset *int64 `json:"offset,omitempty"`

	// Define sorting order.
	Sort []string `json:"sort,omitempty"`

	// The label of the key.
	Label *string `json:"label,omitempty"`

	// Return only managed keys whose activation_date matches the parameter.
	ActivationDate *string `json:"activation_date,omitempty"`

	// Return only managed keys whose activation_date is at or after the parameter value. This query parameter cannot be
	// used in conjunction with the 'activation_date' query parameter.
	ActivationDateMin *string `json:"activation_date_min,omitempty"`

	// Return only managed keys whose activation_date is at or before the parameter value. This query parameter cannot be
	// used in conjunction with the 'activation_date' query parameter.
	ActivationDateMax *string `json:"activation_date_max,omitempty"`

	// Return only managed keys whose deactivation_date matches the parameter. This query parameter cannot be used in
	// conjunction with the 'expiration_date' query parameter.
	DeactivationDate *string `json:"deactivation_date,omitempty"`

	// Return only managed keys whose deactivation_date is at or after the parameter value. This query parameter cannot be
	// used in conjunction with the 'deactivation_date', 'expiration_date', 'expiration_date_min' and 'expiration_date_max'
	// query parameters.
	DeactivationDateMin *string `json:"deactivation_date_min,omitempty"`

	// Return only managed keys whose deactivation_date is at or before the parameter value. This query parameter cannot be
	// used in conjunction with the 'deactivation_date', 'expiration_date', 'expiration_date_min' and 'expiration_date_max'
	// query parameters.
	DeactivationDateMax *string `json:"deactivation_date_max,omitempty"`

	// Return only managed keys whose deactivation_date matches the parameter.
	ExpirationDate *string `json:"expiration_date,omitempty"`

	// Return only managed keys whose deactivation_date is at or after the parameter value. This query parameter cannot be
	// used in conjunction with the 'deactivation_date', 'expiration_date', 'deactivation_date_min' and
	// 'deactivation_date_max' query parameters.
	ExpirationDateMin *string `json:"expiration_date_min,omitempty"`

	// Return only managed keys whose deactivation_date is at or before the parameter value. This query parameter cannot be
	// used in conjunction with the 'deactivation_date', 'expiration_date', 'deactivation_date_min' and
	// 'deactivation_date_max' query parameters.
	ExpirationDateMax *string `json:"expiration_date_max,omitempty"`

	// Return only managed keys whose created_at matches the parameter.
	CreatedAt *string `json:"created_at,omitempty"`

	// Return only managed keys whose created_at is at or after the parameter value. This query parameter cannot be used in
	// conjunction with the 'created_at' query parameter.
	CreatedAtMin *string `json:"created_at_min,omitempty"`

	// Return only managed keys whose created_at is at or before the parameter value. This query parameter cannot be used
	// in conjunction with the 'created_at' query parameter.
	CreatedAtMax *string `json:"created_at_max,omitempty"`

	// Return only managed keys whose updated_at matches the parameter.
	UpdatedAt *string `json:"updated_at,omitempty"`

	// Return only managed keys whose updated_at is after the parameter value. This query parameter cannot be used in
	// conjunction with the 'updated_at' query parameter.
	UpdatedAtMin *string `json:"updated_at_min,omitempty"`

	// Return only managed keys whose updated_at is before the parameter value. This query parameter cannot be used in
	// conjunction with the 'updated_at' query parameter.
	UpdatedAtMax *string `json:"updated_at_max,omitempty"`

	// The size of the key.
	Size *int64 `json:"size,omitempty"`

	// The minimum size of the key. This query parameter cannot be used in conjunction with the 'size' query parameter.
	SizeMin *int64 `json:"size_min,omitempty"`

	// The maximum size of the key. This query parameter cannot be used in conjunction with the 'size' query parameter.
	SizeMax *int64 `json:"size_max,omitempty"`

	// Type of referenced keystore.
	ReferencedKeystoresType []string `json:"referenced_keystores[].type,omitempty"`

	// Name of referenced keystore.
	ReferencedKeystoresName []string `json:"referenced_keystores[].name,omitempty"`

	// Type of keystore supported by one of the instances.
	InstancesKeystoreType []string `json:"instances[].keystore.type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListManagedKeysOptions.Algorithm property.
// The algorithm of the key.
const (
	ListManagedKeysOptions_Algorithm_Aes = "aes"
	ListManagedKeysOptions_Algorithm_Ec = "ec"
	ListManagedKeysOptions_Algorithm_Hmac = "hmac"
	ListManagedKeysOptions_Algorithm_Rsa = "rsa"
)

// Constants associated with the ListManagedKeysOptions.State property.
// The state of the key.
const (
	ListManagedKeysOptions_State_Active = "active"
	ListManagedKeysOptions_State_Deactivated = "deactivated"
	ListManagedKeysOptions_State_Destroyed = "destroyed"
	ListManagedKeysOptions_State_PreActivation = "pre_activation"
)

// Constants associated with the ListManagedKeysOptions.ReferencedKeystoresType property.
// Type of keystore.
const (
	ListManagedKeysOptions_ReferencedKeystoresType_AwsKms = "aws_kms"
	ListManagedKeysOptions_ReferencedKeystoresType_AzureKeyVault = "azure_key_vault"
	ListManagedKeysOptions_ReferencedKeystoresType_GoogleKms = "google_kms"
	ListManagedKeysOptions_ReferencedKeystoresType_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the ListManagedKeysOptions.InstancesKeystoreType property.
// Type of keystore.
const (
	ListManagedKeysOptions_InstancesKeystoreType_AwsKms = "aws_kms"
	ListManagedKeysOptions_InstancesKeystoreType_AzureKeyVault = "azure_key_vault"
	ListManagedKeysOptions_InstancesKeystoreType_GoogleKms = "google_kms"
	ListManagedKeysOptions_InstancesKeystoreType_IbmCloudKms = "ibm_cloud_kms"
)

// NewListManagedKeysOptions : Instantiate ListManagedKeysOptions
func (*UkoV4) NewListManagedKeysOptions() *ListManagedKeysOptions {
	return &ListManagedKeysOptions{}
}

// SetVaultID : Allow user to set VaultID
func (_options *ListManagedKeysOptions) SetVaultID(vaultID []string) *ListManagedKeysOptions {
	_options.VaultID = vaultID
	return _options
}

// SetAlgorithm : Allow user to set Algorithm
func (_options *ListManagedKeysOptions) SetAlgorithm(algorithm []string) *ListManagedKeysOptions {
	_options.Algorithm = algorithm
	return _options
}

// SetState : Allow user to set State
func (_options *ListManagedKeysOptions) SetState(state []string) *ListManagedKeysOptions {
	_options.State = state
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListManagedKeysOptions) SetLimit(limit int64) *ListManagedKeysOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListManagedKeysOptions) SetOffset(offset int64) *ListManagedKeysOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListManagedKeysOptions) SetSort(sort []string) *ListManagedKeysOptions {
	_options.Sort = sort
	return _options
}

// SetLabel : Allow user to set Label
func (_options *ListManagedKeysOptions) SetLabel(label string) *ListManagedKeysOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetActivationDate : Allow user to set ActivationDate
func (_options *ListManagedKeysOptions) SetActivationDate(activationDate string) *ListManagedKeysOptions {
	_options.ActivationDate = core.StringPtr(activationDate)
	return _options
}

// SetActivationDateMin : Allow user to set ActivationDateMin
func (_options *ListManagedKeysOptions) SetActivationDateMin(activationDateMin string) *ListManagedKeysOptions {
	_options.ActivationDateMin = core.StringPtr(activationDateMin)
	return _options
}

// SetActivationDateMax : Allow user to set ActivationDateMax
func (_options *ListManagedKeysOptions) SetActivationDateMax(activationDateMax string) *ListManagedKeysOptions {
	_options.ActivationDateMax = core.StringPtr(activationDateMax)
	return _options
}

// SetDeactivationDate : Allow user to set DeactivationDate
func (_options *ListManagedKeysOptions) SetDeactivationDate(deactivationDate string) *ListManagedKeysOptions {
	_options.DeactivationDate = core.StringPtr(deactivationDate)
	return _options
}

// SetDeactivationDateMin : Allow user to set DeactivationDateMin
func (_options *ListManagedKeysOptions) SetDeactivationDateMin(deactivationDateMin string) *ListManagedKeysOptions {
	_options.DeactivationDateMin = core.StringPtr(deactivationDateMin)
	return _options
}

// SetDeactivationDateMax : Allow user to set DeactivationDateMax
func (_options *ListManagedKeysOptions) SetDeactivationDateMax(deactivationDateMax string) *ListManagedKeysOptions {
	_options.DeactivationDateMax = core.StringPtr(deactivationDateMax)
	return _options
}

// SetExpirationDate : Allow user to set ExpirationDate
func (_options *ListManagedKeysOptions) SetExpirationDate(expirationDate string) *ListManagedKeysOptions {
	_options.ExpirationDate = core.StringPtr(expirationDate)
	return _options
}

// SetExpirationDateMin : Allow user to set ExpirationDateMin
func (_options *ListManagedKeysOptions) SetExpirationDateMin(expirationDateMin string) *ListManagedKeysOptions {
	_options.ExpirationDateMin = core.StringPtr(expirationDateMin)
	return _options
}

// SetExpirationDateMax : Allow user to set ExpirationDateMax
func (_options *ListManagedKeysOptions) SetExpirationDateMax(expirationDateMax string) *ListManagedKeysOptions {
	_options.ExpirationDateMax = core.StringPtr(expirationDateMax)
	return _options
}

// SetCreatedAt : Allow user to set CreatedAt
func (_options *ListManagedKeysOptions) SetCreatedAt(createdAt string) *ListManagedKeysOptions {
	_options.CreatedAt = core.StringPtr(createdAt)
	return _options
}

// SetCreatedAtMin : Allow user to set CreatedAtMin
func (_options *ListManagedKeysOptions) SetCreatedAtMin(createdAtMin string) *ListManagedKeysOptions {
	_options.CreatedAtMin = core.StringPtr(createdAtMin)
	return _options
}

// SetCreatedAtMax : Allow user to set CreatedAtMax
func (_options *ListManagedKeysOptions) SetCreatedAtMax(createdAtMax string) *ListManagedKeysOptions {
	_options.CreatedAtMax = core.StringPtr(createdAtMax)
	return _options
}

// SetUpdatedAt : Allow user to set UpdatedAt
func (_options *ListManagedKeysOptions) SetUpdatedAt(updatedAt string) *ListManagedKeysOptions {
	_options.UpdatedAt = core.StringPtr(updatedAt)
	return _options
}

// SetUpdatedAtMin : Allow user to set UpdatedAtMin
func (_options *ListManagedKeysOptions) SetUpdatedAtMin(updatedAtMin string) *ListManagedKeysOptions {
	_options.UpdatedAtMin = core.StringPtr(updatedAtMin)
	return _options
}

// SetUpdatedAtMax : Allow user to set UpdatedAtMax
func (_options *ListManagedKeysOptions) SetUpdatedAtMax(updatedAtMax string) *ListManagedKeysOptions {
	_options.UpdatedAtMax = core.StringPtr(updatedAtMax)
	return _options
}

// SetSize : Allow user to set Size
func (_options *ListManagedKeysOptions) SetSize(size int64) *ListManagedKeysOptions {
	_options.Size = core.Int64Ptr(size)
	return _options
}

// SetSizeMin : Allow user to set SizeMin
func (_options *ListManagedKeysOptions) SetSizeMin(sizeMin int64) *ListManagedKeysOptions {
	_options.SizeMin = core.Int64Ptr(sizeMin)
	return _options
}

// SetSizeMax : Allow user to set SizeMax
func (_options *ListManagedKeysOptions) SetSizeMax(sizeMax int64) *ListManagedKeysOptions {
	_options.SizeMax = core.Int64Ptr(sizeMax)
	return _options
}

// SetReferencedKeystoresType : Allow user to set ReferencedKeystoresType
func (_options *ListManagedKeysOptions) SetReferencedKeystoresType(referencedKeystoresType []string) *ListManagedKeysOptions {
	_options.ReferencedKeystoresType = referencedKeystoresType
	return _options
}

// SetReferencedKeystoresName : Allow user to set ReferencedKeystoresName
func (_options *ListManagedKeysOptions) SetReferencedKeystoresName(referencedKeystoresName []string) *ListManagedKeysOptions {
	_options.ReferencedKeystoresName = referencedKeystoresName
	return _options
}

// SetInstancesKeystoreType : Allow user to set InstancesKeystoreType
func (_options *ListManagedKeysOptions) SetInstancesKeystoreType(instancesKeystoreType []string) *ListManagedKeysOptions {
	_options.InstancesKeystoreType = instancesKeystoreType
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListManagedKeysOptions) SetHeaders(param map[string]string) *ListManagedKeysOptions {
	options.Headers = param
	return options
}

// ListVaultsOptions : The ListVaults options.
type ListVaultsOptions struct {
	// The number of resources to retrieve.
	Limit *int64 `json:"limit,omitempty"`

	// The number of resources to skip.
	Offset *int64 `json:"offset,omitempty"`

	// Define sorting order.
	Sort []string `json:"sort,omitempty"`

	// Return only vaults whose names begin with the string.
	Name *string `json:"name,omitempty"`

	// Return only vaults whose description contains the string.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListVaultsOptions : Instantiate ListVaultsOptions
func (*UkoV4) NewListVaultsOptions() *ListVaultsOptions {
	return &ListVaultsOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *ListVaultsOptions) SetLimit(limit int64) *ListVaultsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListVaultsOptions) SetOffset(offset int64) *ListVaultsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListVaultsOptions) SetSort(sort []string) *ListVaultsOptions {
	_options.Sort = sort
	return _options
}

// SetName : Allow user to set Name
func (_options *ListVaultsOptions) SetName(name string) *ListVaultsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ListVaultsOptions) SetDescription(description string) *ListVaultsOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListVaultsOptions) SetHeaders(param map[string]string) *ListVaultsOptions {
	options.Headers = param
	return options
}

// SyncManagedKeyOptions : The SyncManagedKey options.
type SyncManagedKeyOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSyncManagedKeyOptions : Instantiate SyncManagedKeyOptions
func (*UkoV4) NewSyncManagedKeyOptions(id string, uKOVault string, ifMatch string) *SyncManagedKeyOptions {
	return &SyncManagedKeyOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *SyncManagedKeyOptions) SetID(id string) *SyncManagedKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *SyncManagedKeyOptions) SetUKOVault(uKOVault string) *SyncManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *SyncManagedKeyOptions) SetIfMatch(ifMatch string) *SyncManagedKeyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SyncManagedKeyOptions) SetHeaders(param map[string]string) *SyncManagedKeyOptions {
	options.Headers = param
	return options
}

// UpdateKeyTemplateOptions : The UpdateKeyTemplate options.
type UpdateKeyTemplateOptions struct {
	// UUID of the template.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Updated keystore related properties.
	Keystores []KeystoresPropertiesUpdateIntf `json:"keystores,omitempty"`

	// Updated description of the key template.
	Description *string `json:"description,omitempty"`

	// Updated key related properties.
	Key *KeyPropertiesUpdate `json:"key,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateKeyTemplateOptions : Instantiate UpdateKeyTemplateOptions
func (*UkoV4) NewUpdateKeyTemplateOptions(id string, uKOVault string, ifMatch string) *UpdateKeyTemplateOptions {
	return &UpdateKeyTemplateOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateKeyTemplateOptions) SetID(id string) *UpdateKeyTemplateOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *UpdateKeyTemplateOptions) SetUKOVault(uKOVault string) *UpdateKeyTemplateOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateKeyTemplateOptions) SetIfMatch(ifMatch string) *UpdateKeyTemplateOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetKeystores : Allow user to set Keystores
func (_options *UpdateKeyTemplateOptions) SetKeystores(keystores []KeystoresPropertiesUpdateIntf) *UpdateKeyTemplateOptions {
	_options.Keystores = keystores
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateKeyTemplateOptions) SetDescription(description string) *UpdateKeyTemplateOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetKey : Allow user to set Key
func (_options *UpdateKeyTemplateOptions) SetKey(key *KeyPropertiesUpdate) *UpdateKeyTemplateOptions {
	_options.Key = key
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateKeyTemplateOptions) SetHeaders(param map[string]string) *UpdateKeyTemplateOptions {
	options.Headers = param
	return options
}

// UpdateKeystoreOptions : The UpdateKeystore options.
type UpdateKeystoreOptions struct {
	// UUID of the keystore.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Keystore properties to update.
	KeystoreBody KeystoreUpdateRequestIntf `json:"keystoreBody" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateKeystoreOptions : Instantiate UpdateKeystoreOptions
func (*UkoV4) NewUpdateKeystoreOptions(id string, uKOVault string, ifMatch string, keystoreBody KeystoreUpdateRequestIntf) *UpdateKeystoreOptions {
	return &UpdateKeystoreOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
		KeystoreBody: keystoreBody,
	}
}

// SetID : Allow user to set ID
func (_options *UpdateKeystoreOptions) SetID(id string) *UpdateKeystoreOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *UpdateKeystoreOptions) SetUKOVault(uKOVault string) *UpdateKeystoreOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateKeystoreOptions) SetIfMatch(ifMatch string) *UpdateKeystoreOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetKeystoreBody : Allow user to set KeystoreBody
func (_options *UpdateKeystoreOptions) SetKeystoreBody(keystoreBody KeystoreUpdateRequestIntf) *UpdateKeystoreOptions {
	_options.KeystoreBody = keystoreBody
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateKeystoreOptions) SetHeaders(param map[string]string) *UpdateKeystoreOptions {
	options.Headers = param
	return options
}

// UpdateManagedKeyFromTemplateOptions : The UpdateManagedKeyFromTemplate options.
type UpdateManagedKeyFromTemplateOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateManagedKeyFromTemplateOptions : Instantiate UpdateManagedKeyFromTemplateOptions
func (*UkoV4) NewUpdateManagedKeyFromTemplateOptions(id string, uKOVault string, ifMatch string) *UpdateManagedKeyFromTemplateOptions {
	return &UpdateManagedKeyFromTemplateOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateManagedKeyFromTemplateOptions) SetID(id string) *UpdateManagedKeyFromTemplateOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *UpdateManagedKeyFromTemplateOptions) SetUKOVault(uKOVault string) *UpdateManagedKeyFromTemplateOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateManagedKeyFromTemplateOptions) SetIfMatch(ifMatch string) *UpdateManagedKeyFromTemplateOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateManagedKeyFromTemplateOptions) SetHeaders(param map[string]string) *UpdateManagedKeyFromTemplateOptions {
	options.Headers = param
	return options
}

// UpdateManagedKeyOptions : The UpdateManagedKey options.
type UpdateManagedKeyOptions struct {
	// UUID of the key.
	ID *string `json:"id" validate:"required,ne="`

	// The UUID of the Vault in which the update is to take place.
	UKOVault *string `json:"UKO-Vault" validate:"required"`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The label of the key.
	Label *string `json:"label,omitempty"`

	// Activation date provided in format: YYYY-MM-DD.
	ActivationDate *strfmt.Date `json:"activation_date,omitempty"`

	// Expiration date provided in format: YYYY-MM-DD.
	ExpirationDate *strfmt.Date `json:"expiration_date,omitempty"`

	// Key-value pairs associated with the key.
	Tags []Tag `json:"tags,omitempty"`

	// Updated description of the managed key.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateManagedKeyOptions : Instantiate UpdateManagedKeyOptions
func (*UkoV4) NewUpdateManagedKeyOptions(id string, uKOVault string, ifMatch string) *UpdateManagedKeyOptions {
	return &UpdateManagedKeyOptions{
		ID: core.StringPtr(id),
		UKOVault: core.StringPtr(uKOVault),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateManagedKeyOptions) SetID(id string) *UpdateManagedKeyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetUKOVault : Allow user to set UKOVault
func (_options *UpdateManagedKeyOptions) SetUKOVault(uKOVault string) *UpdateManagedKeyOptions {
	_options.UKOVault = core.StringPtr(uKOVault)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateManagedKeyOptions) SetIfMatch(ifMatch string) *UpdateManagedKeyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *UpdateManagedKeyOptions) SetLabel(label string) *UpdateManagedKeyOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetActivationDate : Allow user to set ActivationDate
func (_options *UpdateManagedKeyOptions) SetActivationDate(activationDate *strfmt.Date) *UpdateManagedKeyOptions {
	_options.ActivationDate = activationDate
	return _options
}

// SetExpirationDate : Allow user to set ExpirationDate
func (_options *UpdateManagedKeyOptions) SetExpirationDate(expirationDate *strfmt.Date) *UpdateManagedKeyOptions {
	_options.ExpirationDate = expirationDate
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateManagedKeyOptions) SetTags(tags []Tag) *UpdateManagedKeyOptions {
	_options.Tags = tags
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateManagedKeyOptions) SetDescription(description string) *UpdateManagedKeyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateManagedKeyOptions) SetHeaders(param map[string]string) *UpdateManagedKeyOptions {
	options.Headers = param
	return options
}

// UpdateVaultOptions : The UpdateVault options.
type UpdateVaultOptions struct {
	// UUID of the vault.
	ID *string `json:"id" validate:"required,ne="`

	// Precondition of the update; Value of the ETag from the header on a GET request.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Updated name of the vault.
	Name *string `json:"name,omitempty"`

	// Updated description of the vault.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateVaultOptions : Instantiate UpdateVaultOptions
func (*UkoV4) NewUpdateVaultOptions(id string, ifMatch string) *UpdateVaultOptions {
	return &UpdateVaultOptions{
		ID: core.StringPtr(id),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateVaultOptions) SetID(id string) *UpdateVaultOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateVaultOptions) SetIfMatch(ifMatch string) *UpdateVaultOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateVaultOptions) SetName(name string) *UpdateVaultOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateVaultOptions) SetDescription(description string) *UpdateVaultOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateVaultOptions) SetHeaders(param map[string]string) *UpdateVaultOptions {
	options.Headers = param
	return options
}

// ApiError : An error encountered while using the application.
type ApiError struct {
	// The HTTP status code used for the response.
	StatusCode *int64 `json:"status_code,omitempty"`

	// A unique identifier that is attached to the request and the message that allows to refer to the specific transaction
	// or the event chain.
	Trace *string `json:"trace,omitempty"`

	// A list of errors.
	Errors []ErrorModel `json:"errors" validate:"required"`
}

// UnmarshalApiError unmarshals an instance of ApiError from the specified map of raw messages.
func UnmarshalApiError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApiError)
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalErrorModel)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AssociatedResource : Associated resource is any object or entity that is using or referencing a Managed Key in any cloud.
type AssociatedResource struct {
	// An identifier uniquely identifing this associated resource.
	ID *string `json:"id" validate:"required"`

	// Reference to a vault.
	Vault *VaultReference `json:"vault,omitempty"`

	// Reference to a managed key.
	ManagedKey *ManagedKeyReference `json:"managed_key,omitempty"`

	// Reference to a target keystore.
	ReferencedKeystore *TargetKeystoreReference `json:"referenced_keystore,omitempty"`

	// ID of the key in this keystore. Every keystore may use different format such as: UUID, GUID, CRN, URI.
	KeyIdInKeystore *string `json:"key_id_in_keystore" validate:"required"`

	// Name of the associated resource.
	Name *string `json:"name" validate:"required"`

	// Type of the associated resource, in reverse domain name notation. Currently only 'com_ibm_cloud_kms_registration' is
	// supported.
	Type *string `json:"type" validate:"required"`

	// Properties of an associated resource of type IBM Cloud KMS Registration; com_ibm_cloud_kms_registration.
	ComIbmCloudKmsRegistration *IbmCloudKmsRegistration `json:"com_ibm_cloud_kms_registration,omitempty"`
}

// UnmarshalAssociatedResource unmarshals an instance of AssociatedResource from the specified map of raw messages.
func UnmarshalAssociatedResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AssociatedResource)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReference)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "managed_key", &obj.ManagedKey, UnmarshalManagedKeyReference)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "referenced_keystore", &obj.ReferencedKeystore, UnmarshalTargetKeystoreReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_id_in_keystore", &obj.KeyIdInKeystore)
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
	err = core.UnmarshalModel(m, "com_ibm_cloud_kms_registration", &obj.ComIbmCloudKmsRegistration, UnmarshalIbmCloudKmsRegistration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AssociatedResourceList : The base schema for listing associated resources.
type AssociatedResourceList struct {
	// The total count of all objects in the entire collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of resources that were requested in this query.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources that were skipped in this query.
	Offset *int64 `json:"offset" validate:"required"`

	// URL of a resource.
	First *HrefObject `json:"first,omitempty"`

	// URL of a resource.
	Last *HrefObject `json:"last,omitempty"`

	// URL of a resource.
	Previous *HrefObject `json:"previous,omitempty"`

	// URL of a resource.
	Next *HrefObject `json:"next,omitempty"`

	// A list of associated resources.
	AssociatedResources []AssociatedResource `json:"associated_resources" validate:"required"`
}

// UnmarshalAssociatedResourceList unmarshals an instance of AssociatedResourceList from the specified map of raw messages.
func UnmarshalAssociatedResourceList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AssociatedResourceList)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "associated_resources", &obj.AssociatedResources, UnmarshalAssociatedResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AssociatedResourceList) GetNextOffset() (*int64, error) {
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

// ErrorModel : error model.
type ErrorModel struct {
	// Identifier of the error.
	Code *string `json:"code" validate:"required"`

	// A message explaining the problem, with potential suggestions how to address them.
	Message *string `json:"message" validate:"required"`

	// A publicly-accessible URL where information about the error can be read in a web browser. Since more infomration is
	// not always available, this field is optional.
	MoreInfo *string `json:"more_info,omitempty"`

	// Parameters of the message that can be used e.g. for i18n purposes in conjunction with the code. Since the message
	// may not contain any parameters, this field is optional.
	MessageParams []string `json:"message_params,omitempty"`

	// Target of the error.
	Target *Target `json:"target,omitempty"`
}

// UnmarshalErrorModel unmarshals an instance of ErrorModel from the specified map of raw messages.
func UnmarshalErrorModel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ErrorModel)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "more_info", &obj.MoreInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message_params", &obj.MessageParams)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalTarget)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// HrefObject : URL of a resource.
type HrefObject struct {
	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`
}

// UnmarshalHrefObject unmarshals an instance of HrefObject from the specified map of raw messages.
func UnmarshalHrefObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HrefObject)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IbmCloudKmsRegistration : Properties of an associated resource of type IBM Cloud KMS Registration; com_ibm_cloud_kms_registration.
type IbmCloudKmsRegistration struct {
	// A boolean that determines whether Key Protect must prevent deletion of a root key.
	PreventsKeyDeletion *bool `json:"prevents_key_deletion" validate:"required"`

	// Name of the IBM Cloud service, derived from the CRN. It will be empty if UKO is unable to contact the resource
	// controller.
	ServiceName *string `json:"service_name" validate:"required"`

	// Name of the IBM Cloud service's instance, derived from the CRN. It will be empty if UKO is unable to contact the
	// resource controller.
	ServiceInstanceName *string `json:"service_instance_name" validate:"required"`

	// The Cloud Resource Name (CRN) that represents the cloud resource, such as a Cloud Object Storage bucket, that is
	// associated with the key.
	Crn *string `json:"crn" validate:"required"`

	// Description of the purpose of the registration.
	Description *string `json:"description" validate:"required"`
}

// UnmarshalIbmCloudKmsRegistration unmarshals an instance of IbmCloudKmsRegistration from the specified map of raw messages.
func UnmarshalIbmCloudKmsRegistration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IbmCloudKmsRegistration)
	err = core.UnmarshalPrimitive(m, "prevents_key_deletion", &obj.PreventsKeyDeletion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_name", &obj.ServiceInstanceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
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

// InstanceInKeystore : Description of properties of a key within the context of keystores.
type InstanceInKeystore struct {
	Group *string `json:"group" validate:"required"`

	// Type of keystore.
	Type *string `json:"type" validate:"required"`
}

// Constants associated with the InstanceInKeystore.Type property.
// Type of keystore.
const (
	InstanceInKeystore_Type_AwsKms = "aws_kms"
	InstanceInKeystore_Type_AzureKeyVault = "azure_key_vault"
	InstanceInKeystore_Type_GoogleKms = "google_kms"
	InstanceInKeystore_Type_IbmCloudKms = "ibm_cloud_kms"
)

// UnmarshalInstanceInKeystore unmarshals an instance of InstanceInKeystore from the specified map of raw messages.
func UnmarshalInstanceInKeystore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstanceInKeystore)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
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

// KeyInstance : KeyInstance struct
// Models which "extend" this model:
// - KeyInstanceGoogleKms
// - KeyInstanceAwsKms
// - KeyInstanceIbmCloudKms
// - KeyInstanceAzure
type KeyInstance struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// The label of the key.
	LabelInKeystore *string `json:"label_in_keystore,omitempty"`

	// Type of the key instance.
	Type *string `json:"type,omitempty"`

	// Description of properties of a key within the context of keystores.
	Keystore *InstanceInKeystore `json:"keystore,omitempty"`

	GoogleKeyProtectionLevel *string `json:"google_key_protection_level,omitempty"`

	GoogleKeyPurpose *string `json:"google_key_purpose,omitempty"`

	GoogleKmsAlgorithm *string `json:"google_kms_algorithm,omitempty"`
}

// Constants associated with the KeyInstance.Type property.
// Type of the key instance.
const (
	KeyInstance_Type_KeyPair = "key_pair"
	KeyInstance_Type_PrivateKey = "private_key"
	KeyInstance_Type_PublicKey = "public_key"
	KeyInstance_Type_SecretKey = "secret_key"
)

// Constants associated with the KeyInstance.GoogleKeyProtectionLevel property.
const (
	KeyInstance_GoogleKeyProtectionLevel_Hsm = "hsm"
	KeyInstance_GoogleKeyProtectionLevel_Software = "software"
)

// Constants associated with the KeyInstance.GoogleKeyPurpose property.
const (
	KeyInstance_GoogleKeyPurpose_AsymmetricDecrypt = "asymmetric_decrypt"
	KeyInstance_GoogleKeyPurpose_AsymmetricSign = "asymmetric_sign"
	KeyInstance_GoogleKeyPurpose_EncryptDecrypt = "encrypt_decrypt"
	KeyInstance_GoogleKeyPurpose_Mac = "mac"
)

// Constants associated with the KeyInstance.GoogleKmsAlgorithm property.
const (
	KeyInstance_GoogleKmsAlgorithm_EcSignP256Sha256 = "ec_sign_p256_sha256"
	KeyInstance_GoogleKmsAlgorithm_EcSignP384Sha384 = "ec_sign_p384_sha384"
	KeyInstance_GoogleKmsAlgorithm_EcSignSecp256k1Sha256 = "ec_sign_secp256k1_sha256"
	KeyInstance_GoogleKmsAlgorithm_GoogleSymmetricEncryption = "google_symmetric_encryption"
	KeyInstance_GoogleKmsAlgorithm_HmacSha256 = "hmac_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha1 = "rsa_decrypt_oaep_2048_sha1"
	KeyInstance_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha256 = "rsa_decrypt_oaep_2048_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha1 = "rsa_decrypt_oaep_3072_sha1"
	KeyInstance_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha256 = "rsa_decrypt_oaep_3072_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha1 = "rsa_decrypt_oaep_4096_sha1"
	KeyInstance_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha256 = "rsa_decrypt_oaep_4096_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha512 = "rsa_decrypt_oaep_4096_sha512"
	KeyInstance_GoogleKmsAlgorithm_RsaSignPkcs12048Sha256 = "rsa_sign_pkcs1_2048_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaSignPkcs13072Sha256 = "rsa_sign_pkcs1_3072_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaSignPkcs14096Sha256 = "rsa_sign_pkcs1_4096_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaSignPkcs14096Sha512 = "rsa_sign_pkcs1_4096_sha512"
	KeyInstance_GoogleKmsAlgorithm_RsaSignPss2048Sha256 = "rsa_sign_pss_2048_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaSignPss3072Sha256 = "rsa_sign_pss_3072_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaSignPss4096Sha256 = "rsa_sign_pss_4096_sha256"
	KeyInstance_GoogleKmsAlgorithm_RsaSignPss4096Sha512 = "rsa_sign_pss_4096_sha512"
	KeyInstance_GoogleKmsAlgorithm_RsaSignRawPkcs12048 = "rsa_sign_raw_pkcs1_2048"
	KeyInstance_GoogleKmsAlgorithm_RsaSignRawPkcs13072 = "rsa_sign_raw_pkcs1_3072"
	KeyInstance_GoogleKmsAlgorithm_RsaSignRawPkcs14096 = "rsa_sign_raw_pkcs1_4096"
)
func (*KeyInstance) isaKeyInstance() bool {
	return true
}

type KeyInstanceIntf interface {
	isaKeyInstance() bool
}

// UnmarshalKeyInstance unmarshals an instance of KeyInstance from the specified map of raw messages.
func UnmarshalKeyInstance(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyInstance)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label_in_keystore", &obj.LabelInKeystore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "keystore", &obj.Keystore, UnmarshalInstanceInKeystore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_protection_level", &obj.GoogleKeyProtectionLevel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_purpose", &obj.GoogleKeyPurpose)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_kms_algorithm", &obj.GoogleKmsAlgorithm)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeyProperties : Properties describing the properties of the managed key.
type KeyProperties struct {
	// The size of the underlying cryptographic key or key pair. E.g. "256" for AES keys, or "2048" for RSA.
	Size *string `json:"size" validate:"required"`

	// The algorithm of the key.
	Algorithm *string `json:"algorithm" validate:"required"`

	// Key activation date can be provided as a period definition (e.g. PY1 means 1 year).
	ActivationDate *string `json:"activation_date" validate:"required"`

	// Key expiration date can be provided as a period definition (e.g. PY1 means 1 year).
	ExpirationDate *string `json:"expiration_date" validate:"required"`

	// The state that the key will be in after generation.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the KeyProperties.Algorithm property.
// The algorithm of the key.
const (
	KeyProperties_Algorithm_Aes = "aes"
	KeyProperties_Algorithm_Ec = "ec"
	KeyProperties_Algorithm_Hmac = "hmac"
	KeyProperties_Algorithm_Rsa = "rsa"
)

// Constants associated with the KeyProperties.State property.
// The state that the key will be in after generation.
const (
	KeyProperties_State_Active = "active"
	KeyProperties_State_PreActivation = "pre_activation"
)

// NewKeyProperties : Instantiate KeyProperties (Generic Model Constructor)
func (*UkoV4) NewKeyProperties(size string, algorithm string, activationDate string, expirationDate string, state string) (_model *KeyProperties, err error) {
	_model = &KeyProperties{
		Size: core.StringPtr(size),
		Algorithm: core.StringPtr(algorithm),
		ActivationDate: core.StringPtr(activationDate),
		ExpirationDate: core.StringPtr(expirationDate),
		State: core.StringPtr(state),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalKeyProperties unmarshals an instance of KeyProperties from the specified map of raw messages.
func UnmarshalKeyProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyProperties)
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "activation_date", &obj.ActivationDate)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeyPropertiesUpdate : Updated key related properties.
type KeyPropertiesUpdate struct {
	// The size of the underlying cryptographic key or key pair. E.g. "256" for AES keys, or "2048" for RSA.
	Size *string `json:"size,omitempty"`

	// Key activation date can be provided as a period definition (e.g. PY1 means 1 year).
	ActivationDate *string `json:"activation_date,omitempty"`

	// Key expiration date can be provided as a period definition (e.g. PY1 means 1 year).
	ExpirationDate *string `json:"expiration_date,omitempty"`

	// The state that the key will be in after generation.
	State *string `json:"state,omitempty"`
}

// Constants associated with the KeyPropertiesUpdate.State property.
// The state that the key will be in after generation.
const (
	KeyPropertiesUpdate_State_Active = "active"
	KeyPropertiesUpdate_State_PreActivation = "pre_activation"
)

// UnmarshalKeyPropertiesUpdate unmarshals an instance of KeyPropertiesUpdate from the specified map of raw messages.
func UnmarshalKeyPropertiesUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyPropertiesUpdate)
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "activation_date", &obj.ActivationDate)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeyVerificationPattern : Key verification pattern is used to identify and distinguish cryptographic keys.
type KeyVerificationPattern struct {
	// The method used for calculating the verification pattern.
	Method *string `json:"method" validate:"required"`

	// The calculated value.
	Value *string `json:"value" validate:"required"`
}

// UnmarshalKeyVerificationPattern unmarshals an instance of KeyVerificationPattern from the specified map of raw messages.
func UnmarshalKeyVerificationPattern(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyVerificationPattern)
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
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

// Keystore : A target keystore is keystore that is assigned to a vault. If it is an internal keystore, it can be created only in a
// vault. If it is an external keystore, you need to assign the external keystore to a vault when you connect your
// service instance to it.
// Models which "extend" this model:
// - KeystoreTypeGoogleKms
// - KeystoreTypeAwsKms
// - KeystoreTypeAzure
// - KeystoreTypeIbmCloudKms
type Keystore struct {
	// Reference to a vault.
	Vault *VaultReference `json:"vault,omitempty"`

	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the target keystore. It can be changed in the future.
	Name *string `json:"name,omitempty"`

	// Geographic location of the keystore, if available.
	Location *string `json:"location,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// List of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// Type of keystore.
	Type *string `json:"type,omitempty"`

	// Date and time when the target keystore was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Date and time when the target keystore was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// ID of the user that created the key.
	CreatedBy *string `json:"created_by,omitempty"`

	// ID of the user that last updated the key.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`

	// The value of the JSON key represented in the Base64 format.
	GoogleCredentials *string `json:"google_credentials,omitempty"`

	// Location represents the geographical region where a Cloud KMS resource is stored and can be accessed. A key's
	// location impacts the performance of applications using the key.
	GoogleLocation *string `json:"google_location,omitempty"`

	// The project id associated with this keystore.
	GoogleProjectID *string `json:"google_project_id,omitempty"`

	// The private key id associated with this keystore.
	GooglePrivateKeyID *string `json:"google_private_key_id,omitempty"`

	// A key ring organizes keys in a specific Google Cloud location and allows you to manage access control on groups of
	// keys.
	GoogleKeyRing *string `json:"google_key_ring,omitempty"`

	// AWS Region.
	AwsRegion *string `json:"aws_region,omitempty"`

	// The access key id used for connecting to this instance of AWS KMS.
	AwsAccessKeyID *string `json:"aws_access_key_id,omitempty"`

	// The secret access key used for connecting to this instance of AWS KMS.
	AwsSecretAccessKey *string `json:"aws_secret_access_key,omitempty"`

	// Service name of the key vault instance from the Azure portal.
	AzureServiceName *string `json:"azure_service_name,omitempty"`

	// Resource group in Azure.
	AzureResourceGroup *string `json:"azure_resource_group,omitempty"`

	// Location of the Azure Key Vault.
	AzureLocation *string `json:"azure_location,omitempty"`

	// Azure service principal client ID.
	AzureServicePrincipalClientID *string `json:"azure_service_principal_client_id,omitempty"`

	// Azure service principal password.
	AzureServicePrincipalPassword *string `json:"azure_service_principal_password,omitempty"`

	// Azure tenant that the Key Vault is associated with,.
	AzureTenant *string `json:"azure_tenant,omitempty"`

	// Subscription ID in Azure.
	AzureSubscriptionID *string `json:"azure_subscription_id,omitempty"`

	// Azure environment, usually 'Azure'.
	AzureEnvironment *string `json:"azure_environment,omitempty"`

	// API endpoint of the IBM Cloud keystore.
	IbmApiEndpoint *string `json:"ibm_api_endpoint,omitempty"`

	// Endpoint of the IAM service for this IBM Cloud keystore.
	IbmIamEndpoint *string `json:"ibm_iam_endpoint,omitempty"`

	// The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.
	IbmApiKey *string `json:"ibm_api_key,omitempty"`

	// The instance ID of the IBM Cloud keystore.
	IbmInstanceID *string `json:"ibm_instance_id,omitempty"`

	// Possible IBM Cloud KMS variants.
	IbmVariant *string `json:"ibm_variant,omitempty"`

	// The key ring of an IBM Cloud KMS Keystore.
	IbmKeyRing *string `json:"ibm_key_ring,omitempty"`
}

// Constants associated with the Keystore.Type property.
// Type of keystore.
const (
	Keystore_Type_AwsKms = "aws_kms"
	Keystore_Type_AzureKeyVault = "azure_key_vault"
	Keystore_Type_GoogleKms = "google_kms"
	Keystore_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the Keystore.AwsRegion property.
// AWS Region.
const (
	Keystore_AwsRegion_AfSouth1 = "af_south_1"
	Keystore_AwsRegion_ApEast1 = "ap_east_1"
	Keystore_AwsRegion_ApNortheast1 = "ap_northeast_1"
	Keystore_AwsRegion_ApNortheast2 = "ap_northeast_2"
	Keystore_AwsRegion_ApSouth1 = "ap_south_1"
	Keystore_AwsRegion_ApSoutheast1 = "ap_southeast_1"
	Keystore_AwsRegion_ApSoutheast2 = "ap_southeast_2"
	Keystore_AwsRegion_AwsCnGlobal = "aws_cn_global"
	Keystore_AwsRegion_AwsGlobal = "aws_global"
	Keystore_AwsRegion_AwsIsoBGlobal = "aws_iso_b_global"
	Keystore_AwsRegion_AwsIsoGlobal = "aws_iso_global"
	Keystore_AwsRegion_AwsUsGovGlobal = "aws_us_gov_global"
	Keystore_AwsRegion_CaCentral1 = "ca_central_1"
	Keystore_AwsRegion_CnNorth1 = "cn_north_1"
	Keystore_AwsRegion_CnNorthwest1 = "cn_northwest_1"
	Keystore_AwsRegion_EuCentral1 = "eu_central_1"
	Keystore_AwsRegion_EuWest1 = "eu_west_1"
	Keystore_AwsRegion_EuWest2 = "eu_west_2"
	Keystore_AwsRegion_EuWest3 = "eu_west_3"
	Keystore_AwsRegion_MeSouth1 = "me_south_1"
	Keystore_AwsRegion_SaEast1 = "sa_east_1"
	Keystore_AwsRegion_UsEast1 = "us_east_1"
	Keystore_AwsRegion_UsEast2 = "us_east_2"
	Keystore_AwsRegion_UsGovEast1 = "us_gov_east_1"
	Keystore_AwsRegion_UsGovWest1 = "us_gov_west_1"
	Keystore_AwsRegion_UsIsoEast1 = "us_iso_east_1"
	Keystore_AwsRegion_UsIsobEast1 = "us_isob_east_1"
	Keystore_AwsRegion_UsWest1 = "us_west_1"
	Keystore_AwsRegion_UsWest2 = "us_west_2"
)

// Constants associated with the Keystore.AzureLocation property.
// Location of the Azure Key Vault.
const (
	Keystore_AzureLocation_AsiaEast = "asia_east"
	Keystore_AzureLocation_AsiaSoutheast = "asia_southeast"
	Keystore_AzureLocation_AustraliaCentral = "australia_central"
	Keystore_AzureLocation_AustraliaCentral2 = "australia_central_2"
	Keystore_AzureLocation_AustraliaEast = "australia_east"
	Keystore_AzureLocation_AustraliaSoutheast = "australia_southeast"
	Keystore_AzureLocation_BrazilSouth = "brazil_south"
	Keystore_AzureLocation_CanadaCentral = "canada_central"
	Keystore_AzureLocation_CanadaEast = "canada_east"
	Keystore_AzureLocation_ChinaEast = "china_east"
	Keystore_AzureLocation_ChinaEast2 = "china_east_2"
	Keystore_AzureLocation_ChinaNorth = "china_north"
	Keystore_AzureLocation_ChinaNorth2 = "china_north_2"
	Keystore_AzureLocation_EuropeNorth = "europe_north"
	Keystore_AzureLocation_EuropeWest = "europe_west"
	Keystore_AzureLocation_FranceCentral = "france_central"
	Keystore_AzureLocation_FranceSouth = "france_south"
	Keystore_AzureLocation_GermanyCentral = "germany_central"
	Keystore_AzureLocation_GermanyNortheast = "germany_northeast"
	Keystore_AzureLocation_IndiaCentral = "india_central"
	Keystore_AzureLocation_IndiaSouth = "india_south"
	Keystore_AzureLocation_IndiaWest = "india_west"
	Keystore_AzureLocation_JapanEast = "japan_east"
	Keystore_AzureLocation_JapanWest = "japan_west"
	Keystore_AzureLocation_KoreaCentral = "korea_central"
	Keystore_AzureLocation_KoreaSouth = "korea_south"
	Keystore_AzureLocation_SouthAfricaNorth = "south_africa_north"
	Keystore_AzureLocation_SouthAfricaWest = "south_africa_west"
	Keystore_AzureLocation_UkSouth = "uk_south"
	Keystore_AzureLocation_UkWest = "uk_west"
	Keystore_AzureLocation_UsCentral = "us_central"
	Keystore_AzureLocation_UsDodCentral = "us_dod_central"
	Keystore_AzureLocation_UsDodEast = "us_dod_east"
	Keystore_AzureLocation_UsEast = "us_east"
	Keystore_AzureLocation_UsEast2 = "us_east_2"
	Keystore_AzureLocation_UsGovArizona = "us_gov_arizona"
	Keystore_AzureLocation_UsGovIowa = "us_gov_iowa"
	Keystore_AzureLocation_UsGovTexas = "us_gov_texas"
	Keystore_AzureLocation_UsGovVirginia = "us_gov_virginia"
	Keystore_AzureLocation_UsNorthCentral = "us_north_central"
	Keystore_AzureLocation_UsSouthCentral = "us_south_central"
	Keystore_AzureLocation_UsWest = "us_west"
	Keystore_AzureLocation_UsWest2 = "us_west_2"
	Keystore_AzureLocation_UsWestCentral = "us_west_central"
)

// Constants associated with the Keystore.AzureEnvironment property.
// Azure environment, usually 'Azure'.
const (
	Keystore_AzureEnvironment_Azure = "azure"
	Keystore_AzureEnvironment_AzureChina = "azure_china"
	Keystore_AzureEnvironment_AzureGermany = "azure_germany"
	Keystore_AzureEnvironment_AzureUsGovernment = "azure_us_government"
)

// Constants associated with the Keystore.IbmVariant property.
// Possible IBM Cloud KMS variants.
const (
	Keystore_IbmVariant_Hpcs = "hpcs"
	Keystore_IbmVariant_Internal = "internal"
	Keystore_IbmVariant_KeyProtect = "key_protect"
)
func (*Keystore) isaKeystore() bool {
	return true
}

type KeystoreIntf interface {
	isaKeystore() bool
}

// UnmarshalKeystore unmarshals an instance of Keystore from the specified map of raw messages.
func UnmarshalKeystore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Keystore)
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReference)
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
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_credentials", &obj.GoogleCredentials)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_location", &obj.GoogleLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_project_id", &obj.GoogleProjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_private_key_id", &obj.GooglePrivateKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_ring", &obj.GoogleKeyRing)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_access_key_id", &obj.AwsAccessKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_secret_access_key", &obj.AwsSecretAccessKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_name", &obj.AzureServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_location", &obj.AzureLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_client_id", &obj.AzureServicePrincipalClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_password", &obj.AzureServicePrincipalPassword)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_tenant", &obj.AzureTenant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_environment", &obj.AzureEnvironment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_endpoint", &obj.IbmApiEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_iam_endpoint", &obj.IbmIamEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IbmApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_instance_id", &obj.IbmInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_variant", &obj.IbmVariant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_key_ring", &obj.IbmKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequest : Properties required to create a keystore.
// Models which "extend" this model:
// - KeystoreCreationRequestKeystoreTypeAwsKmsCreate
// - KeystoreCreationRequestKeystoreTypeGoogleKmsCreate
// - KeystoreCreationRequestKeystoreTypeAzureCreate
// - KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate
type KeystoreCreationRequest struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	// ID of the Vault where the entity is to be created in.
	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// AWS Region.
	AwsRegion *string `json:"aws_region,omitempty"`

	// The access key id used for connecting to this instance of AWS KMS.
	AwsAccessKeyID *string `json:"aws_access_key_id,omitempty"`

	// The secret access key used for connecting to this instance of AWS KMS.
	AwsSecretAccessKey *string `json:"aws_secret_access_key,omitempty"`

	// The value of the JSON key represented in the Base64 format.
	GoogleCredentials *string `json:"google_credentials,omitempty"`

	// Location represents the geographical region where a Cloud KMS resource is stored and can be accessed. A key's
	// location impacts the performance of applications using the key.
	GoogleLocation *string `json:"google_location,omitempty"`

	// The project id associated with this keystore.
	GoogleProjectID *string `json:"google_project_id,omitempty"`

	// The private key id associated with this keystore.
	GooglePrivateKeyID *string `json:"google_private_key_id,omitempty"`

	// A key ring organizes keys in a specific Google Cloud location and allows you to manage access control on groups of
	// keys.
	GoogleKeyRing *string `json:"google_key_ring,omitempty"`

	// Service name of the key vault instance from the Azure portal.
	AzureServiceName *string `json:"azure_service_name,omitempty"`

	// Resource group in Azure.
	AzureResourceGroup *string `json:"azure_resource_group,omitempty"`

	// Location of the Azure Key Vault.
	AzureLocation *string `json:"azure_location,omitempty"`

	// Azure service principal client ID.
	AzureServicePrincipalClientID *string `json:"azure_service_principal_client_id,omitempty"`

	// Azure service principal password.
	AzureServicePrincipalPassword *string `json:"azure_service_principal_password,omitempty"`

	// Azure tenant that the Key Vault is associated with,.
	AzureTenant *string `json:"azure_tenant,omitempty"`

	// Subscription ID in Azure.
	AzureSubscriptionID *string `json:"azure_subscription_id,omitempty"`

	// Azure environment, usually 'Azure'.
	AzureEnvironment *string `json:"azure_environment,omitempty"`

	// Possible IBM Cloud KMS variants.
	IbmVariant *string `json:"ibm_variant,omitempty"`

	// API endpoint of the IBM Cloud keystore.
	IbmApiEndpoint *string `json:"ibm_api_endpoint,omitempty"`

	// Endpoint of the IAM service for this IBM Cloud keystore.
	IbmIamEndpoint *string `json:"ibm_iam_endpoint,omitempty"`

	// The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.
	IbmApiKey *string `json:"ibm_api_key,omitempty"`

	// The instance ID of the IBM Cloud keystore.
	IbmInstanceID *string `json:"ibm_instance_id,omitempty"`

	// The key ring of an IBM Cloud KMS Keystore.
	IbmKeyRing *string `json:"ibm_key_ring,omitempty"`
}

// Constants associated with the KeystoreCreationRequest.Type property.
// Type of keystore.
const (
	KeystoreCreationRequest_Type_AwsKms = "aws_kms"
	KeystoreCreationRequest_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequest_Type_GoogleKms = "google_kms"
	KeystoreCreationRequest_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreCreationRequest.AwsRegion property.
// AWS Region.
const (
	KeystoreCreationRequest_AwsRegion_AfSouth1 = "af_south_1"
	KeystoreCreationRequest_AwsRegion_ApEast1 = "ap_east_1"
	KeystoreCreationRequest_AwsRegion_ApNortheast1 = "ap_northeast_1"
	KeystoreCreationRequest_AwsRegion_ApNortheast2 = "ap_northeast_2"
	KeystoreCreationRequest_AwsRegion_ApSouth1 = "ap_south_1"
	KeystoreCreationRequest_AwsRegion_ApSoutheast1 = "ap_southeast_1"
	KeystoreCreationRequest_AwsRegion_ApSoutheast2 = "ap_southeast_2"
	KeystoreCreationRequest_AwsRegion_AwsCnGlobal = "aws_cn_global"
	KeystoreCreationRequest_AwsRegion_AwsGlobal = "aws_global"
	KeystoreCreationRequest_AwsRegion_AwsIsoBGlobal = "aws_iso_b_global"
	KeystoreCreationRequest_AwsRegion_AwsIsoGlobal = "aws_iso_global"
	KeystoreCreationRequest_AwsRegion_AwsUsGovGlobal = "aws_us_gov_global"
	KeystoreCreationRequest_AwsRegion_CaCentral1 = "ca_central_1"
	KeystoreCreationRequest_AwsRegion_CnNorth1 = "cn_north_1"
	KeystoreCreationRequest_AwsRegion_CnNorthwest1 = "cn_northwest_1"
	KeystoreCreationRequest_AwsRegion_EuCentral1 = "eu_central_1"
	KeystoreCreationRequest_AwsRegion_EuWest1 = "eu_west_1"
	KeystoreCreationRequest_AwsRegion_EuWest2 = "eu_west_2"
	KeystoreCreationRequest_AwsRegion_EuWest3 = "eu_west_3"
	KeystoreCreationRequest_AwsRegion_MeSouth1 = "me_south_1"
	KeystoreCreationRequest_AwsRegion_SaEast1 = "sa_east_1"
	KeystoreCreationRequest_AwsRegion_UsEast1 = "us_east_1"
	KeystoreCreationRequest_AwsRegion_UsEast2 = "us_east_2"
	KeystoreCreationRequest_AwsRegion_UsGovEast1 = "us_gov_east_1"
	KeystoreCreationRequest_AwsRegion_UsGovWest1 = "us_gov_west_1"
	KeystoreCreationRequest_AwsRegion_UsIsoEast1 = "us_iso_east_1"
	KeystoreCreationRequest_AwsRegion_UsIsobEast1 = "us_isob_east_1"
	KeystoreCreationRequest_AwsRegion_UsWest1 = "us_west_1"
	KeystoreCreationRequest_AwsRegion_UsWest2 = "us_west_2"
)

// Constants associated with the KeystoreCreationRequest.AzureLocation property.
// Location of the Azure Key Vault.
const (
	KeystoreCreationRequest_AzureLocation_AsiaEast = "asia_east"
	KeystoreCreationRequest_AzureLocation_AsiaSoutheast = "asia_southeast"
	KeystoreCreationRequest_AzureLocation_AustraliaCentral = "australia_central"
	KeystoreCreationRequest_AzureLocation_AustraliaCentral2 = "australia_central_2"
	KeystoreCreationRequest_AzureLocation_AustraliaEast = "australia_east"
	KeystoreCreationRequest_AzureLocation_AustraliaSoutheast = "australia_southeast"
	KeystoreCreationRequest_AzureLocation_BrazilSouth = "brazil_south"
	KeystoreCreationRequest_AzureLocation_CanadaCentral = "canada_central"
	KeystoreCreationRequest_AzureLocation_CanadaEast = "canada_east"
	KeystoreCreationRequest_AzureLocation_ChinaEast = "china_east"
	KeystoreCreationRequest_AzureLocation_ChinaEast2 = "china_east_2"
	KeystoreCreationRequest_AzureLocation_ChinaNorth = "china_north"
	KeystoreCreationRequest_AzureLocation_ChinaNorth2 = "china_north_2"
	KeystoreCreationRequest_AzureLocation_EuropeNorth = "europe_north"
	KeystoreCreationRequest_AzureLocation_EuropeWest = "europe_west"
	KeystoreCreationRequest_AzureLocation_FranceCentral = "france_central"
	KeystoreCreationRequest_AzureLocation_FranceSouth = "france_south"
	KeystoreCreationRequest_AzureLocation_GermanyCentral = "germany_central"
	KeystoreCreationRequest_AzureLocation_GermanyNortheast = "germany_northeast"
	KeystoreCreationRequest_AzureLocation_IndiaCentral = "india_central"
	KeystoreCreationRequest_AzureLocation_IndiaSouth = "india_south"
	KeystoreCreationRequest_AzureLocation_IndiaWest = "india_west"
	KeystoreCreationRequest_AzureLocation_JapanEast = "japan_east"
	KeystoreCreationRequest_AzureLocation_JapanWest = "japan_west"
	KeystoreCreationRequest_AzureLocation_KoreaCentral = "korea_central"
	KeystoreCreationRequest_AzureLocation_KoreaSouth = "korea_south"
	KeystoreCreationRequest_AzureLocation_SouthAfricaNorth = "south_africa_north"
	KeystoreCreationRequest_AzureLocation_SouthAfricaWest = "south_africa_west"
	KeystoreCreationRequest_AzureLocation_UkSouth = "uk_south"
	KeystoreCreationRequest_AzureLocation_UkWest = "uk_west"
	KeystoreCreationRequest_AzureLocation_UsCentral = "us_central"
	KeystoreCreationRequest_AzureLocation_UsDodCentral = "us_dod_central"
	KeystoreCreationRequest_AzureLocation_UsDodEast = "us_dod_east"
	KeystoreCreationRequest_AzureLocation_UsEast = "us_east"
	KeystoreCreationRequest_AzureLocation_UsEast2 = "us_east_2"
	KeystoreCreationRequest_AzureLocation_UsGovArizona = "us_gov_arizona"
	KeystoreCreationRequest_AzureLocation_UsGovIowa = "us_gov_iowa"
	KeystoreCreationRequest_AzureLocation_UsGovTexas = "us_gov_texas"
	KeystoreCreationRequest_AzureLocation_UsGovVirginia = "us_gov_virginia"
	KeystoreCreationRequest_AzureLocation_UsNorthCentral = "us_north_central"
	KeystoreCreationRequest_AzureLocation_UsSouthCentral = "us_south_central"
	KeystoreCreationRequest_AzureLocation_UsWest = "us_west"
	KeystoreCreationRequest_AzureLocation_UsWest2 = "us_west_2"
	KeystoreCreationRequest_AzureLocation_UsWestCentral = "us_west_central"
)

// Constants associated with the KeystoreCreationRequest.AzureEnvironment property.
// Azure environment, usually 'Azure'.
const (
	KeystoreCreationRequest_AzureEnvironment_Azure = "azure"
	KeystoreCreationRequest_AzureEnvironment_AzureChina = "azure_china"
	KeystoreCreationRequest_AzureEnvironment_AzureGermany = "azure_germany"
	KeystoreCreationRequest_AzureEnvironment_AzureUsGovernment = "azure_us_government"
)

// Constants associated with the KeystoreCreationRequest.IbmVariant property.
// Possible IBM Cloud KMS variants.
const (
	KeystoreCreationRequest_IbmVariant_Hpcs = "hpcs"
	KeystoreCreationRequest_IbmVariant_Internal = "internal"
	KeystoreCreationRequest_IbmVariant_KeyProtect = "key_protect"
)
func (*KeystoreCreationRequest) isaKeystoreCreationRequest() bool {
	return true
}

type KeystoreCreationRequestIntf interface {
	isaKeystoreCreationRequest() bool
}

// UnmarshalKeystoreCreationRequest unmarshals an instance of KeystoreCreationRequest from the specified map of raw messages.
func UnmarshalKeystoreCreationRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "type", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'type': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'type' not found in JSON object")
		return
	}
	if discValue == "google_kms" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKeystoreCreationRequestKeystoreTypeGoogleKmsCreate)
	} else if discValue == "aws_kms" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKeystoreCreationRequestKeystoreTypeAwsKmsCreate)
	} else if discValue == "azure_key_vault" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKeystoreCreationRequestKeystoreTypeAzureCreate)
	} else if discValue == "ibm_cloud_kms" {
		err = core.UnmarshalModel(m, "", result, UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'type': %s", discValue)
	}
	return
}

// KeystoreList : The base schema for listing target keystores.
type KeystoreList struct {
	// The total count of all objects in the entire collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of resources that were requested in this query.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources that were skipped in this query.
	Offset *int64 `json:"offset" validate:"required"`

	// URL of a resource.
	First *HrefObject `json:"first,omitempty"`

	// URL of a resource.
	Last *HrefObject `json:"last,omitempty"`

	// URL of a resource.
	Previous *HrefObject `json:"previous,omitempty"`

	// URL of a resource.
	Next *HrefObject `json:"next,omitempty"`

	// A list of target keystores.
	Keystores []KeystoreIntf `json:"keystores" validate:"required"`
}

// UnmarshalKeystoreList unmarshals an instance of KeystoreList from the specified map of raw messages.
func UnmarshalKeystoreList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreList)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "keystores", &obj.Keystores, UnmarshalKeystore)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *KeystoreList) GetNextOffset() (*int64, error) {
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

// KeystoreStatus : The status of the connection to the keystore.
type KeystoreStatus struct {
	// Date of last successful communication with a keystore.
	LastHeartbeat *strfmt.DateTime `json:"last_heartbeat,omitempty"`

	// Possible states of a keystore.
	HealthStatus *string `json:"health_status,omitempty"`

	// Message returned with the status.
	Message *string `json:"message,omitempty"`
}

// Constants associated with the KeystoreStatus.HealthStatus property.
// Possible states of a keystore.
const (
	KeystoreStatus_HealthStatus_ConfigurationError = "configuration_error"
	KeystoreStatus_HealthStatus_NotResponding = "not_responding"
	KeystoreStatus_HealthStatus_Ok = "ok"
)

// UnmarshalKeystoreStatus unmarshals an instance of KeystoreStatus from the specified map of raw messages.
func UnmarshalKeystoreStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreStatus)
	err = core.UnmarshalPrimitive(m, "last_heartbeat", &obj.LastHeartbeat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "health_status", &obj.HealthStatus)
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

// KeystoreUpdateRequest : Properties of a keystore that can be updated.
// Models which "extend" this model:
// - KeystoreUpdateRequestKeystoreTypeGoogleKmsUpdate
// - KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate
// - KeystoreUpdateRequestKeystoreTypeAzureUpdate
// - KeystoreUpdateRequestKeystoreTypeIbmCloudKmsUpdate
// - KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate
type KeystoreUpdateRequest struct {
	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// The value of the JSON key represented in the Base64 format.
	GoogleCredentials *string `json:"google_credentials,omitempty"`

	// Location represents the geographical region where a Cloud KMS resource is stored and can be accessed. A key's
	// location impacts the performance of applications using the key.
	GoogleLocation *string `json:"google_location,omitempty"`

	// The project id associated with this keystore.
	GoogleProjectID *string `json:"google_project_id,omitempty"`

	// The private key id associated with this keystore.
	GooglePrivateKeyID *string `json:"google_private_key_id,omitempty"`

	// A key ring organizes keys in a specific Google Cloud location and allows you to manage access control on groups of
	// keys.
	GoogleKeyRing *string `json:"google_key_ring,omitempty"`

	// AWS Region.
	AwsRegion *string `json:"aws_region,omitempty"`

	// The access key id used for connecting to this instance of AWS KMS.
	AwsAccessKeyID *string `json:"aws_access_key_id,omitempty"`

	// The secret access key used for connecting to this instance of AWS KMS.
	AwsSecretAccessKey *string `json:"aws_secret_access_key,omitempty"`

	// Service name of the key vault instance from the Azure portal.
	AzureServiceName *string `json:"azure_service_name,omitempty"`

	// Resource group in Azure.
	AzureResourceGroup *string `json:"azure_resource_group,omitempty"`

	// Location of the Azure Key Vault.
	AzureLocation *string `json:"azure_location,omitempty"`

	// Azure service principal client ID.
	AzureServicePrincipalClientID *string `json:"azure_service_principal_client_id,omitempty"`

	// Azure service principal password.
	AzureServicePrincipalPassword *string `json:"azure_service_principal_password,omitempty"`

	// Azure tenant that the Key Vault is associated with,.
	AzureTenant *string `json:"azure_tenant,omitempty"`

	// Subscription ID in Azure.
	AzureSubscriptionID *string `json:"azure_subscription_id,omitempty"`

	// Azure environment, usually 'Azure'.
	AzureEnvironment *string `json:"azure_environment,omitempty"`

	// API endpoint of the IBM Cloud keystore.
	IbmApiEndpoint *string `json:"ibm_api_endpoint,omitempty"`

	// Endpoint of the IAM service for this IBM Cloud keystore.
	IbmIamEndpoint *string `json:"ibm_iam_endpoint,omitempty"`

	// The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.
	IbmApiKey *string `json:"ibm_api_key,omitempty"`

	// The instance ID of the IBM Cloud keystore.
	IbmInstanceID *string `json:"ibm_instance_id,omitempty"`

	// The key ring of an IBM Cloud KMS Keystore.
	IbmKeyRing *string `json:"ibm_key_ring,omitempty"`
}

// Constants associated with the KeystoreUpdateRequest.AwsRegion property.
// AWS Region.
const (
	KeystoreUpdateRequest_AwsRegion_AfSouth1 = "af_south_1"
	KeystoreUpdateRequest_AwsRegion_ApEast1 = "ap_east_1"
	KeystoreUpdateRequest_AwsRegion_ApNortheast1 = "ap_northeast_1"
	KeystoreUpdateRequest_AwsRegion_ApNortheast2 = "ap_northeast_2"
	KeystoreUpdateRequest_AwsRegion_ApSouth1 = "ap_south_1"
	KeystoreUpdateRequest_AwsRegion_ApSoutheast1 = "ap_southeast_1"
	KeystoreUpdateRequest_AwsRegion_ApSoutheast2 = "ap_southeast_2"
	KeystoreUpdateRequest_AwsRegion_AwsCnGlobal = "aws_cn_global"
	KeystoreUpdateRequest_AwsRegion_AwsGlobal = "aws_global"
	KeystoreUpdateRequest_AwsRegion_AwsIsoBGlobal = "aws_iso_b_global"
	KeystoreUpdateRequest_AwsRegion_AwsIsoGlobal = "aws_iso_global"
	KeystoreUpdateRequest_AwsRegion_AwsUsGovGlobal = "aws_us_gov_global"
	KeystoreUpdateRequest_AwsRegion_CaCentral1 = "ca_central_1"
	KeystoreUpdateRequest_AwsRegion_CnNorth1 = "cn_north_1"
	KeystoreUpdateRequest_AwsRegion_CnNorthwest1 = "cn_northwest_1"
	KeystoreUpdateRequest_AwsRegion_EuCentral1 = "eu_central_1"
	KeystoreUpdateRequest_AwsRegion_EuWest1 = "eu_west_1"
	KeystoreUpdateRequest_AwsRegion_EuWest2 = "eu_west_2"
	KeystoreUpdateRequest_AwsRegion_EuWest3 = "eu_west_3"
	KeystoreUpdateRequest_AwsRegion_MeSouth1 = "me_south_1"
	KeystoreUpdateRequest_AwsRegion_SaEast1 = "sa_east_1"
	KeystoreUpdateRequest_AwsRegion_UsEast1 = "us_east_1"
	KeystoreUpdateRequest_AwsRegion_UsEast2 = "us_east_2"
	KeystoreUpdateRequest_AwsRegion_UsGovEast1 = "us_gov_east_1"
	KeystoreUpdateRequest_AwsRegion_UsGovWest1 = "us_gov_west_1"
	KeystoreUpdateRequest_AwsRegion_UsIsoEast1 = "us_iso_east_1"
	KeystoreUpdateRequest_AwsRegion_UsIsobEast1 = "us_isob_east_1"
	KeystoreUpdateRequest_AwsRegion_UsWest1 = "us_west_1"
	KeystoreUpdateRequest_AwsRegion_UsWest2 = "us_west_2"
)

// Constants associated with the KeystoreUpdateRequest.AzureLocation property.
// Location of the Azure Key Vault.
const (
	KeystoreUpdateRequest_AzureLocation_AsiaEast = "asia_east"
	KeystoreUpdateRequest_AzureLocation_AsiaSoutheast = "asia_southeast"
	KeystoreUpdateRequest_AzureLocation_AustraliaCentral = "australia_central"
	KeystoreUpdateRequest_AzureLocation_AustraliaCentral2 = "australia_central_2"
	KeystoreUpdateRequest_AzureLocation_AustraliaEast = "australia_east"
	KeystoreUpdateRequest_AzureLocation_AustraliaSoutheast = "australia_southeast"
	KeystoreUpdateRequest_AzureLocation_BrazilSouth = "brazil_south"
	KeystoreUpdateRequest_AzureLocation_CanadaCentral = "canada_central"
	KeystoreUpdateRequest_AzureLocation_CanadaEast = "canada_east"
	KeystoreUpdateRequest_AzureLocation_ChinaEast = "china_east"
	KeystoreUpdateRequest_AzureLocation_ChinaEast2 = "china_east_2"
	KeystoreUpdateRequest_AzureLocation_ChinaNorth = "china_north"
	KeystoreUpdateRequest_AzureLocation_ChinaNorth2 = "china_north_2"
	KeystoreUpdateRequest_AzureLocation_EuropeNorth = "europe_north"
	KeystoreUpdateRequest_AzureLocation_EuropeWest = "europe_west"
	KeystoreUpdateRequest_AzureLocation_FranceCentral = "france_central"
	KeystoreUpdateRequest_AzureLocation_FranceSouth = "france_south"
	KeystoreUpdateRequest_AzureLocation_GermanyCentral = "germany_central"
	KeystoreUpdateRequest_AzureLocation_GermanyNortheast = "germany_northeast"
	KeystoreUpdateRequest_AzureLocation_IndiaCentral = "india_central"
	KeystoreUpdateRequest_AzureLocation_IndiaSouth = "india_south"
	KeystoreUpdateRequest_AzureLocation_IndiaWest = "india_west"
	KeystoreUpdateRequest_AzureLocation_JapanEast = "japan_east"
	KeystoreUpdateRequest_AzureLocation_JapanWest = "japan_west"
	KeystoreUpdateRequest_AzureLocation_KoreaCentral = "korea_central"
	KeystoreUpdateRequest_AzureLocation_KoreaSouth = "korea_south"
	KeystoreUpdateRequest_AzureLocation_SouthAfricaNorth = "south_africa_north"
	KeystoreUpdateRequest_AzureLocation_SouthAfricaWest = "south_africa_west"
	KeystoreUpdateRequest_AzureLocation_UkSouth = "uk_south"
	KeystoreUpdateRequest_AzureLocation_UkWest = "uk_west"
	KeystoreUpdateRequest_AzureLocation_UsCentral = "us_central"
	KeystoreUpdateRequest_AzureLocation_UsDodCentral = "us_dod_central"
	KeystoreUpdateRequest_AzureLocation_UsDodEast = "us_dod_east"
	KeystoreUpdateRequest_AzureLocation_UsEast = "us_east"
	KeystoreUpdateRequest_AzureLocation_UsEast2 = "us_east_2"
	KeystoreUpdateRequest_AzureLocation_UsGovArizona = "us_gov_arizona"
	KeystoreUpdateRequest_AzureLocation_UsGovIowa = "us_gov_iowa"
	KeystoreUpdateRequest_AzureLocation_UsGovTexas = "us_gov_texas"
	KeystoreUpdateRequest_AzureLocation_UsGovVirginia = "us_gov_virginia"
	KeystoreUpdateRequest_AzureLocation_UsNorthCentral = "us_north_central"
	KeystoreUpdateRequest_AzureLocation_UsSouthCentral = "us_south_central"
	KeystoreUpdateRequest_AzureLocation_UsWest = "us_west"
	KeystoreUpdateRequest_AzureLocation_UsWest2 = "us_west_2"
	KeystoreUpdateRequest_AzureLocation_UsWestCentral = "us_west_central"
)

// Constants associated with the KeystoreUpdateRequest.AzureEnvironment property.
// Azure environment, usually 'Azure'.
const (
	KeystoreUpdateRequest_AzureEnvironment_Azure = "azure"
	KeystoreUpdateRequest_AzureEnvironment_AzureChina = "azure_china"
	KeystoreUpdateRequest_AzureEnvironment_AzureGermany = "azure_germany"
	KeystoreUpdateRequest_AzureEnvironment_AzureUsGovernment = "azure_us_government"
)
func (*KeystoreUpdateRequest) isaKeystoreUpdateRequest() bool {
	return true
}

type KeystoreUpdateRequestIntf interface {
	isaKeystoreUpdateRequest() bool
}

// UnmarshalKeystoreUpdateRequest unmarshals an instance of KeystoreUpdateRequest from the specified map of raw messages.
func UnmarshalKeystoreUpdateRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreUpdateRequest)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_credentials", &obj.GoogleCredentials)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_location", &obj.GoogleLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_project_id", &obj.GoogleProjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_private_key_id", &obj.GooglePrivateKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_ring", &obj.GoogleKeyRing)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_access_key_id", &obj.AwsAccessKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_secret_access_key", &obj.AwsSecretAccessKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_name", &obj.AzureServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_location", &obj.AzureLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_client_id", &obj.AzureServicePrincipalClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_password", &obj.AzureServicePrincipalPassword)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_tenant", &obj.AzureTenant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_environment", &obj.AzureEnvironment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_endpoint", &obj.IbmApiEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_iam_endpoint", &obj.IbmIamEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IbmApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_instance_id", &obj.IbmInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_key_ring", &obj.IbmKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoresPropertiesCreate : KeystoresPropertiesCreate struct
// Models which "extend" this model:
// - KeystoresPropertiesCreateGoogleKms
// - KeystoresPropertiesCreateAwsKms
// - KeystoresPropertiesCreateIbmCloudKms
// - KeystoresPropertiesCreateAzure
type KeystoresPropertiesCreate struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`

	// Type of keystore.
	Type *string `json:"type,omitempty"`

	GoogleKeyProtectionLevel *string `json:"google_key_protection_level,omitempty"`

	GoogleKeyPurpose *string `json:"google_key_purpose,omitempty"`

	GoogleKmsAlgorithm *string `json:"google_kms_algorithm,omitempty"`
}

// Constants associated with the KeystoresPropertiesCreate.Type property.
// Type of keystore.
const (
	KeystoresPropertiesCreate_Type_AwsKms = "aws_kms"
	KeystoresPropertiesCreate_Type_AzureKeyVault = "azure_key_vault"
	KeystoresPropertiesCreate_Type_GoogleKms = "google_kms"
	KeystoresPropertiesCreate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoresPropertiesCreate.GoogleKeyProtectionLevel property.
const (
	KeystoresPropertiesCreate_GoogleKeyProtectionLevel_Hsm = "hsm"
	KeystoresPropertiesCreate_GoogleKeyProtectionLevel_Software = "software"
)

// Constants associated with the KeystoresPropertiesCreate.GoogleKeyPurpose property.
const (
	KeystoresPropertiesCreate_GoogleKeyPurpose_AsymmetricDecrypt = "asymmetric_decrypt"
	KeystoresPropertiesCreate_GoogleKeyPurpose_AsymmetricSign = "asymmetric_sign"
	KeystoresPropertiesCreate_GoogleKeyPurpose_EncryptDecrypt = "encrypt_decrypt"
	KeystoresPropertiesCreate_GoogleKeyPurpose_Mac = "mac"
)

// Constants associated with the KeystoresPropertiesCreate.GoogleKmsAlgorithm property.
const (
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_EcSignP256Sha256 = "ec_sign_p256_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_EcSignP384Sha384 = "ec_sign_p384_sha384"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_EcSignSecp256k1Sha256 = "ec_sign_secp256k1_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_GoogleSymmetricEncryption = "google_symmetric_encryption"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_HmacSha256 = "hmac_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha1 = "rsa_decrypt_oaep_2048_sha1"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha256 = "rsa_decrypt_oaep_2048_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha1 = "rsa_decrypt_oaep_3072_sha1"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha256 = "rsa_decrypt_oaep_3072_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha1 = "rsa_decrypt_oaep_4096_sha1"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha256 = "rsa_decrypt_oaep_4096_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha512 = "rsa_decrypt_oaep_4096_sha512"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignPkcs12048Sha256 = "rsa_sign_pkcs1_2048_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignPkcs13072Sha256 = "rsa_sign_pkcs1_3072_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignPkcs14096Sha256 = "rsa_sign_pkcs1_4096_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignPkcs14096Sha512 = "rsa_sign_pkcs1_4096_sha512"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignPss2048Sha256 = "rsa_sign_pss_2048_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignPss3072Sha256 = "rsa_sign_pss_3072_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignPss4096Sha256 = "rsa_sign_pss_4096_sha256"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignPss4096Sha512 = "rsa_sign_pss_4096_sha512"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignRawPkcs12048 = "rsa_sign_raw_pkcs1_2048"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignRawPkcs13072 = "rsa_sign_raw_pkcs1_3072"
	KeystoresPropertiesCreate_GoogleKmsAlgorithm_RsaSignRawPkcs14096 = "rsa_sign_raw_pkcs1_4096"
)
func (*KeystoresPropertiesCreate) isaKeystoresPropertiesCreate() bool {
	return true
}

type KeystoresPropertiesCreateIntf interface {
	isaKeystoresPropertiesCreate() bool
}

// UnmarshalKeystoresPropertiesCreate unmarshals an instance of KeystoresPropertiesCreate from the specified map of raw messages.
func UnmarshalKeystoresPropertiesCreate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesCreate)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_protection_level", &obj.GoogleKeyProtectionLevel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_purpose", &obj.GoogleKeyPurpose)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_kms_algorithm", &obj.GoogleKmsAlgorithm)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoresPropertiesUpdate : KeystoresPropertiesUpdate struct
// Models which "extend" this model:
// - KeystoresPropertiesUpdateGoogleKms
// - KeystoresPropertiesUpdateAwsKms
// - KeystoresPropertiesUpdateIbmCloudKms
// - KeystoresPropertiesUpdateAzure
type KeystoresPropertiesUpdate struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`

	GoogleKeyProtectionLevel *string `json:"google_key_protection_level,omitempty"`

	GoogleKeyPurpose *string `json:"google_key_purpose,omitempty"`

	GoogleKmsAlgorithm *string `json:"google_kms_algorithm,omitempty"`
}

// Constants associated with the KeystoresPropertiesUpdate.GoogleKeyProtectionLevel property.
const (
	KeystoresPropertiesUpdate_GoogleKeyProtectionLevel_Hsm = "hsm"
	KeystoresPropertiesUpdate_GoogleKeyProtectionLevel_Software = "software"
)

// Constants associated with the KeystoresPropertiesUpdate.GoogleKeyPurpose property.
const (
	KeystoresPropertiesUpdate_GoogleKeyPurpose_AsymmetricDecrypt = "asymmetric_decrypt"
	KeystoresPropertiesUpdate_GoogleKeyPurpose_AsymmetricSign = "asymmetric_sign"
	KeystoresPropertiesUpdate_GoogleKeyPurpose_EncryptDecrypt = "encrypt_decrypt"
	KeystoresPropertiesUpdate_GoogleKeyPurpose_Mac = "mac"
)

// Constants associated with the KeystoresPropertiesUpdate.GoogleKmsAlgorithm property.
const (
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_EcSignP256Sha256 = "ec_sign_p256_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_EcSignP384Sha384 = "ec_sign_p384_sha384"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_EcSignSecp256k1Sha256 = "ec_sign_secp256k1_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_GoogleSymmetricEncryption = "google_symmetric_encryption"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_HmacSha256 = "hmac_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha1 = "rsa_decrypt_oaep_2048_sha1"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha256 = "rsa_decrypt_oaep_2048_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha1 = "rsa_decrypt_oaep_3072_sha1"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha256 = "rsa_decrypt_oaep_3072_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha1 = "rsa_decrypt_oaep_4096_sha1"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha256 = "rsa_decrypt_oaep_4096_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha512 = "rsa_decrypt_oaep_4096_sha512"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignPkcs12048Sha256 = "rsa_sign_pkcs1_2048_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignPkcs13072Sha256 = "rsa_sign_pkcs1_3072_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignPkcs14096Sha256 = "rsa_sign_pkcs1_4096_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignPkcs14096Sha512 = "rsa_sign_pkcs1_4096_sha512"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignPss2048Sha256 = "rsa_sign_pss_2048_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignPss3072Sha256 = "rsa_sign_pss_3072_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignPss4096Sha256 = "rsa_sign_pss_4096_sha256"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignPss4096Sha512 = "rsa_sign_pss_4096_sha512"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignRawPkcs12048 = "rsa_sign_raw_pkcs1_2048"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignRawPkcs13072 = "rsa_sign_raw_pkcs1_3072"
	KeystoresPropertiesUpdate_GoogleKmsAlgorithm_RsaSignRawPkcs14096 = "rsa_sign_raw_pkcs1_4096"
)
func (*KeystoresPropertiesUpdate) isaKeystoresPropertiesUpdate() bool {
	return true
}

type KeystoresPropertiesUpdateIntf interface {
	isaKeystoresPropertiesUpdate() bool
}

// UnmarshalKeystoresPropertiesUpdate unmarshals an instance of KeystoresPropertiesUpdate from the specified map of raw messages.
func UnmarshalKeystoresPropertiesUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesUpdate)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_protection_level", &obj.GoogleKeyProtectionLevel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_purpose", &obj.GoogleKeyPurpose)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_kms_algorithm", &obj.GoogleKmsAlgorithm)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ManagedKey : A managed key is a key that is created in and assigned to a vault. You can manage the lifecycle of a managed key and
// install it to multiple keystores in the same vault. You can use a managed key for encryption and decryption only when
// it is installed in at least one target keystore. Installing a managed key in multiple keystores in the same vault
// enables key redundancy. To use a managed key for encryption and decryption, install in one or more keystores within
// the same vault first.
type ManagedKey struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Reference to a vault.
	Vault *VaultReference `json:"vault,omitempty"`

	// Reference to a key template.
	Template *TemplateReference `json:"template,omitempty"`

	// Description of the managed key.
	Description *string `json:"description,omitempty"`

	// The label of the key.
	Label *string `json:"label" validate:"required"`

	// The state of the key.
	State *string `json:"state" validate:"required"`

	// The size of the underlying cryptographic key or key pair. E.g. "256" for AES keys, or "2048" for RSA.
	Size *string `json:"size,omitempty"`

	// The algorithm of the key.
	Algorithm *string `json:"algorithm" validate:"required"`

	// A list of verification patterns of the key (e.g. public key hash for RSA keys).
	VerificationPatterns []KeyVerificationPattern `json:"verification_patterns,omitempty"`

	// First day when the key is active.
	ActivationDate *strfmt.Date `json:"activation_date,omitempty"`

	// Last day when the key is active.
	ExpirationDate *strfmt.Date `json:"expiration_date,omitempty"`

	// Key-value pairs associated with the key.
	Tags []Tag `json:"tags,omitempty"`

	// Date and time when the key was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Date and time when the key was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// ID of the user that created the key.
	CreatedBy *string `json:"created_by,omitempty"`

	// ID of the user that last updated the key.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// referenced keystores.
	ReferencedKeystores []TargetKeystoreReference `json:"referenced_keystores" validate:"required"`

	// key instances.
	Instances []KeyInstanceIntf `json:"instances" validate:"required"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the ManagedKey.State property.
// The state of the key.
const (
	ManagedKey_State_Active = "active"
	ManagedKey_State_Deactivated = "deactivated"
	ManagedKey_State_Destroyed = "destroyed"
	ManagedKey_State_PreActivation = "pre_activation"
)

// Constants associated with the ManagedKey.Algorithm property.
// The algorithm of the key.
const (
	ManagedKey_Algorithm_Aes = "aes"
	ManagedKey_Algorithm_Ec = "ec"
	ManagedKey_Algorithm_Hmac = "hmac"
	ManagedKey_Algorithm_Rsa = "rsa"
)

// UnmarshalManagedKey unmarshals an instance of ManagedKey from the specified map of raw messages.
func UnmarshalManagedKey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ManagedKey)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReference)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalTemplateReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "verification_patterns", &obj.VerificationPatterns, UnmarshalKeyVerificationPattern)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "activation_date", &obj.ActivationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_date", &obj.ExpirationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalTag)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "referenced_keystores", &obj.ReferencedKeystores, UnmarshalTargetKeystoreReference)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "instances", &obj.Instances, UnmarshalKeyInstance)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ManagedKeyList : The base schema for listing managed keys.
type ManagedKeyList struct {
	// The total count of all objects in the entire collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of resources that were requested in this query.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources that were skipped in this query.
	Offset *int64 `json:"offset" validate:"required"`

	// URL of a resource.
	First *HrefObject `json:"first,omitempty"`

	// URL of a resource.
	Last *HrefObject `json:"last,omitempty"`

	// URL of a resource.
	Previous *HrefObject `json:"previous,omitempty"`

	// URL of a resource.
	Next *HrefObject `json:"next,omitempty"`

	// A list of managed keys.
	ManagedKeys []ManagedKey `json:"managed_keys" validate:"required"`
}

// UnmarshalManagedKeyList unmarshals an instance of ManagedKeyList from the specified map of raw messages.
func UnmarshalManagedKeyList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ManagedKeyList)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "managed_keys", &obj.ManagedKeys, UnmarshalManagedKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ManagedKeyList) GetNextOffset() (*int64, error) {
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

// ManagedKeyReference : Reference to a managed key.
type ManagedKeyReference struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// The label of the key.
	Label *string `json:"label" validate:"required"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`
}

// UnmarshalManagedKeyReference unmarshals an instance of ManagedKeyReference from the specified map of raw messages.
func UnmarshalManagedKeyReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ManagedKeyReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StatusInKeystore : Describes the status of a key in a single keystore.
type StatusInKeystore struct {
	// Reference to a target keystore.
	Keystore *TargetKeystoreReference `json:"keystore,omitempty"`

	// Possible states of a key in keystore.  not_present: the key is not in the target keystore at all,  active: the key
	// is in the target keystore, and can be used for its intended purpose not_active: the key is in the target keystore,
	// but cannot be used for its intended purpose wrong_key: there is a key in the target keystore, but it doesn't have
	// the value that is expected error: there was an error checking the status of the key in the target keystore.
	Status *string `json:"status" validate:"required"`

	// Flag to represent synchronization status between UKO Managed Key and Target Keystore. Possible status flags. ok:
	// managed key state is the same as target keystore state, out_of_sync: managed key state is different than target
	// keystore state.
	KeystoreSyncFlag *string `json:"keystore_sync_flag" validate:"required"`

	// Detailed description to represents every possible state combination or mismatch between UKO Managed Key and Target
	// Keystore.
	KeystoreSyncFlagDetail *string `json:"keystore_sync_flag_detail" validate:"required"`

	// An error encountered while using the application.
	Error *ApiError `json:"error,omitempty"`

	// ID of the key in this keystore. Every keystore may use different format such as: UUID, GUID, CRN, URI.
	KeyIdInKeystore *string `json:"key_id_in_keystore,omitempty"`
}

// Constants associated with the StatusInKeystore.Status property.
// Possible states of a key in keystore.  not_present: the key is not in the target keystore at all,  active: the key is
// in the target keystore, and can be used for its intended purpose not_active: the key is in the target keystore, but
// cannot be used for its intended purpose wrong_key: there is a key in the target keystore, but it doesn't have the
// value that is expected error: there was an error checking the status of the key in the target keystore.
const (
	StatusInKeystore_Status_Active = "active"
	StatusInKeystore_Status_Error = "error"
	StatusInKeystore_Status_NotActive = "not_active"
	StatusInKeystore_Status_NotPresent = "not_present"
	StatusInKeystore_Status_WrongKey = "wrong_key"
)

// Constants associated with the StatusInKeystore.KeystoreSyncFlag property.
// Flag to represent synchronization status between UKO Managed Key and Target Keystore. Possible status flags. ok:
// managed key state is the same as target keystore state, out_of_sync: managed key state is different than target
// keystore state.
const (
	StatusInKeystore_KeystoreSyncFlag_Error = "error"
	StatusInKeystore_KeystoreSyncFlag_Ok = "ok"
	StatusInKeystore_KeystoreSyncFlag_OutOfSync = "out_of_sync"
)

// Constants associated with the StatusInKeystore.KeystoreSyncFlagDetail property.
// Detailed description to represents every possible state combination or mismatch between UKO Managed Key and Target
// Keystore.
const (
	StatusInKeystore_KeystoreSyncFlagDetail_ActiveKeyIsActiveInKeystore = "active_key_is_active_in_keystore"
	StatusInKeystore_KeystoreSyncFlagDetail_ActiveKeyIsNotActiveInKeystore = "active_key_is_not_active_in_keystore"
	StatusInKeystore_KeystoreSyncFlagDetail_ConnectionError = "connection_error"
	StatusInKeystore_KeystoreSyncFlagDetail_DeactivatedKeyIsDeactivatedInKeystore = "deactivated_key_is_deactivated_in_keystore"
	StatusInKeystore_KeystoreSyncFlagDetail_DeactivatedKeyIsNotDeactivatedInKeystore = "deactivated_key_is_not_deactivated_in_keystore"
	StatusInKeystore_KeystoreSyncFlagDetail_DestroyedKeyIsNotPresentInKeystore = "destroyed_key_is_not_present_in_keystore"
	StatusInKeystore_KeystoreSyncFlagDetail_DestroyedKeyIsPresentInKeystore = "destroyed_key_is_present_in_keystore"
	StatusInKeystore_KeystoreSyncFlagDetail_PreActiveKeyIsNotPresentInKeystore = "pre_active_key_is_not_present_in_keystore"
	StatusInKeystore_KeystoreSyncFlagDetail_PreActiveKeyIsPresentInKeystore = "pre_active_key_is_present_in_keystore"
	StatusInKeystore_KeystoreSyncFlagDetail_TargetKeystoreRemovedByUser = "target_keystore_removed_by_user"
	StatusInKeystore_KeystoreSyncFlagDetail_TargetKeystoreRemovedByUserContainsAnActiveKey = "target_keystore_removed_by_user_contains_an_active_key"
)

// UnmarshalStatusInKeystore unmarshals an instance of StatusInKeystore from the specified map of raw messages.
func UnmarshalStatusInKeystore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StatusInKeystore)
	err = core.UnmarshalModel(m, "keystore", &obj.Keystore, UnmarshalTargetKeystoreReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "keystore_sync_flag", &obj.KeystoreSyncFlag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "keystore_sync_flag_detail", &obj.KeystoreSyncFlagDetail)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "error", &obj.Error, UnmarshalApiError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_id_in_keystore", &obj.KeyIdInKeystore)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StatusInKeystores : Status of a key in keystores.
type StatusInKeystores struct {
	// list of key statuses in keystores.
	StatusInKeystores []StatusInKeystore `json:"status_in_keystores" validate:"required"`
}

// UnmarshalStatusInKeystores unmarshals an instance of StatusInKeystores from the specified map of raw messages.
func UnmarshalStatusInKeystores(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StatusInKeystores)
	err = core.UnmarshalModel(m, "status_in_keystores", &obj.StatusInKeystores, UnmarshalStatusInKeystore)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tag : A key tag, as used when creating keys and searching.
type Tag struct {
	// Name of a tag.
	Name *string `json:"name" validate:"required"`

	// Value of a tag.
	Value *string `json:"value" validate:"required"`
}

// NewTag : Instantiate Tag (Generic Model Constructor)
func (*UkoV4) NewTag(name string, value string) (_model *Tag, err error) {
	_model = &Tag{
		Name: core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTag unmarshals an instance of Tag from the specified map of raw messages.
func UnmarshalTag(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tag)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

// Target : Target of the error.
type Target struct {
	// Type, one of 'field', 'parameter', or 'header'.
	Type *string `json:"type,omitempty"`

	// Name of the field (with dot-syntax if necessary), query parameter, or header.
	Name *string `json:"name,omitempty"`
}

// UnmarshalTarget unmarshals an instance of Target from the specified map of raw messages.
func UnmarshalTarget(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Target)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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

// TargetKeystoreReference : Reference to a target keystore.
type TargetKeystoreReference struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the target keystore.
	Name *string `json:"name,omitempty"`

	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TargetKeystoreReference.Type property.
// Type of keystore.
const (
	TargetKeystoreReference_Type_AwsKms = "aws_kms"
	TargetKeystoreReference_Type_AzureKeyVault = "azure_key_vault"
	TargetKeystoreReference_Type_GoogleKms = "google_kms"
	TargetKeystoreReference_Type_IbmCloudKms = "ibm_cloud_kms"
)

// UnmarshalTargetKeystoreReference unmarshals an instance of TargetKeystoreReference from the specified map of raw messages.
func UnmarshalTargetKeystoreReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetKeystoreReference)
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
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Template : A template and all its properties.
type Template struct {
	// Reference to a vault.
	Vault *VaultReference `json:"vault,omitempty"`

	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Version of the key template. Every time the key template is updated, the version will be updated automatically.
	Version *int64 `json:"version,omitempty"`

	// Name of the key template.
	Name *string `json:"name,omitempty"`

	// Properties describing the properties of the managed key.
	Key *KeyProperties `json:"key" validate:"required"`

	// Description of the key template.
	Description *string `json:"description" validate:"required"`

	// Date and time when the key template was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Date and time when the key template was updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// ID of the user that created the key template.
	CreatedBy *string `json:"created_by,omitempty"`

	// ID of the user that updated the key.
	UpdatedBy *string `json:"updated_by,omitempty"`

	Keystores []KeystoresPropertiesCreateIntf `json:"keystores" validate:"required"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`
}

// UnmarshalTemplate unmarshals an instance of Template from the specified map of raw messages.
func UnmarshalTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Template)
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "key", &obj.Key, UnmarshalKeyProperties)
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
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "keystores", &obj.Keystores, UnmarshalKeystoresPropertiesCreate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateList : The base schema for listing key templates.
type TemplateList struct {
	// The total count of all objects in the entire collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of resources that were requested in this query.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources that were skipped in this query.
	Offset *int64 `json:"offset" validate:"required"`

	// URL of a resource.
	First *HrefObject `json:"first,omitempty"`

	// URL of a resource.
	Last *HrefObject `json:"last,omitempty"`

	// URL of a resource.
	Previous *HrefObject `json:"previous,omitempty"`

	// URL of a resource.
	Next *HrefObject `json:"next,omitempty"`

	// A list of key templates.
	Templates []Template `json:"templates" validate:"required"`
}

// UnmarshalTemplateList unmarshals an instance of TemplateList from the specified map of raw messages.
func UnmarshalTemplateList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateList)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "templates", &obj.Templates, UnmarshalTemplate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *TemplateList) GetNextOffset() (*int64, error) {
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

// TemplateReference : Reference to a key template.
type TemplateReference struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the key template.
	Name *string `json:"name,omitempty"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`
}

// UnmarshalTemplateReference unmarshals an instance of TemplateReference from the specified map of raw messages.
func UnmarshalTemplateReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Vault : Vaults are secure repositories for your cryptographic keys and keystores. A managed key or keystore can only be in
// one vault at a time.
type Vault struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the vault.
	Name *string `json:"name" validate:"required"`

	// Description of the vault.
	Description *string `json:"description" validate:"required"`

	// Date and time when the vault was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Date and time when the vault was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// ID of the user that created the vault.
	CreatedBy *string `json:"created_by,omitempty"`

	// ID of the user that last updated the vault.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`
}

// UnmarshalVault unmarshals an instance of Vault from the specified map of raw messages.
func UnmarshalVault(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Vault)
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
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VaultList : The base schema for listing vaults.
type VaultList struct {
	// The total count of all objects in the entire collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The number of resources that were requested in this query.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of resources that were skipped in this query.
	Offset *int64 `json:"offset" validate:"required"`

	// URL of a resource.
	First *HrefObject `json:"first,omitempty"`

	// URL of a resource.
	Last *HrefObject `json:"last,omitempty"`

	// URL of a resource.
	Previous *HrefObject `json:"previous,omitempty"`

	// URL of a resource.
	Next *HrefObject `json:"next,omitempty"`

	// A list of vaults.
	Vaults []Vault `json:"vaults" validate:"required"`
}

// UnmarshalVaultList unmarshals an instance of VaultList from the specified map of raw messages.
func UnmarshalVaultList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VaultList)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalHrefObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vaults", &obj.Vaults, UnmarshalVault)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *VaultList) GetNextOffset() (*int64, error) {
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

// VaultReference : Reference to a vault.
type VaultReference struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the referenced vault.
	Name *string `json:"name,omitempty"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`
}

// UnmarshalVaultReference unmarshals an instance of VaultReference from the specified map of raw messages.
func UnmarshalVaultReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VaultReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VaultReferenceInCreationRequest : ID of the Vault where the entity is to be created in.
type VaultReferenceInCreationRequest struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id" validate:"required"`
}

// NewVaultReferenceInCreationRequest : Instantiate VaultReferenceInCreationRequest (Generic Model Constructor)
func (*UkoV4) NewVaultReferenceInCreationRequest(id string) (_model *VaultReferenceInCreationRequest, err error) {
	_model = &VaultReferenceInCreationRequest{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalVaultReferenceInCreationRequest unmarshals an instance of VaultReferenceInCreationRequest from the specified map of raw messages.
func UnmarshalVaultReferenceInCreationRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VaultReferenceInCreationRequest)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeyInstanceAwsKms : The instance of a managed key for a specific keystore.
// This model "extends" KeyInstance
type KeyInstanceAwsKms struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id" validate:"required"`

	// The label of the key.
	LabelInKeystore *string `json:"label_in_keystore" validate:"required"`

	// Type of the key instance.
	Type *string `json:"type,omitempty"`

	// Description of properties of a key within the context of keystores.
	Keystore *InstanceInKeystore `json:"keystore" validate:"required"`
}

// Constants associated with the KeyInstanceAwsKms.Type property.
// Type of the key instance.
const (
	KeyInstanceAwsKms_Type_KeyPair = "key_pair"
	KeyInstanceAwsKms_Type_PrivateKey = "private_key"
	KeyInstanceAwsKms_Type_PublicKey = "public_key"
	KeyInstanceAwsKms_Type_SecretKey = "secret_key"
)

func (*KeyInstanceAwsKms) isaKeyInstance() bool {
	return true
}

// UnmarshalKeyInstanceAwsKms unmarshals an instance of KeyInstanceAwsKms from the specified map of raw messages.
func UnmarshalKeyInstanceAwsKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyInstanceAwsKms)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label_in_keystore", &obj.LabelInKeystore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "keystore", &obj.Keystore, UnmarshalInstanceInKeystore)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeyInstanceAzure : The instance of a managed key for a specific keystore.
// This model "extends" KeyInstance
type KeyInstanceAzure struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id" validate:"required"`

	// The label of the key.
	LabelInKeystore *string `json:"label_in_keystore" validate:"required"`

	// Type of the key instance.
	Type *string `json:"type,omitempty"`

	// Description of properties of a key within the context of keystores.
	Keystore *InstanceInKeystore `json:"keystore" validate:"required"`
}

// Constants associated with the KeyInstanceAzure.Type property.
// Type of the key instance.
const (
	KeyInstanceAzure_Type_KeyPair = "key_pair"
	KeyInstanceAzure_Type_PrivateKey = "private_key"
	KeyInstanceAzure_Type_PublicKey = "public_key"
	KeyInstanceAzure_Type_SecretKey = "secret_key"
)

func (*KeyInstanceAzure) isaKeyInstance() bool {
	return true
}

// UnmarshalKeyInstanceAzure unmarshals an instance of KeyInstanceAzure from the specified map of raw messages.
func UnmarshalKeyInstanceAzure(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyInstanceAzure)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label_in_keystore", &obj.LabelInKeystore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "keystore", &obj.Keystore, UnmarshalInstanceInKeystore)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeyInstanceGoogleKms : The instance of a managed key for a specific keystore.
// This model "extends" KeyInstance
type KeyInstanceGoogleKms struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id" validate:"required"`

	// The label of the key.
	LabelInKeystore *string `json:"label_in_keystore" validate:"required"`

	// Type of the key instance.
	Type *string `json:"type,omitempty"`

	// Description of properties of a key within the context of keystores.
	Keystore *InstanceInKeystore `json:"keystore" validate:"required"`

	GoogleKeyProtectionLevel *string `json:"google_key_protection_level,omitempty"`

	GoogleKeyPurpose *string `json:"google_key_purpose,omitempty"`

	GoogleKmsAlgorithm *string `json:"google_kms_algorithm,omitempty"`
}

// Constants associated with the KeyInstanceGoogleKms.Type property.
// Type of the key instance.
const (
	KeyInstanceGoogleKms_Type_KeyPair = "key_pair"
	KeyInstanceGoogleKms_Type_PrivateKey = "private_key"
	KeyInstanceGoogleKms_Type_PublicKey = "public_key"
	KeyInstanceGoogleKms_Type_SecretKey = "secret_key"
)

// Constants associated with the KeyInstanceGoogleKms.GoogleKeyProtectionLevel property.
const (
	KeyInstanceGoogleKms_GoogleKeyProtectionLevel_Hsm = "hsm"
	KeyInstanceGoogleKms_GoogleKeyProtectionLevel_Software = "software"
)

// Constants associated with the KeyInstanceGoogleKms.GoogleKeyPurpose property.
const (
	KeyInstanceGoogleKms_GoogleKeyPurpose_AsymmetricDecrypt = "asymmetric_decrypt"
	KeyInstanceGoogleKms_GoogleKeyPurpose_AsymmetricSign = "asymmetric_sign"
	KeyInstanceGoogleKms_GoogleKeyPurpose_EncryptDecrypt = "encrypt_decrypt"
	KeyInstanceGoogleKms_GoogleKeyPurpose_Mac = "mac"
)

// Constants associated with the KeyInstanceGoogleKms.GoogleKmsAlgorithm property.
const (
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_EcSignP256Sha256 = "ec_sign_p256_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_EcSignP384Sha384 = "ec_sign_p384_sha384"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_EcSignSecp256k1Sha256 = "ec_sign_secp256k1_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_GoogleSymmetricEncryption = "google_symmetric_encryption"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_HmacSha256 = "hmac_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha1 = "rsa_decrypt_oaep_2048_sha1"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha256 = "rsa_decrypt_oaep_2048_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha1 = "rsa_decrypt_oaep_3072_sha1"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha256 = "rsa_decrypt_oaep_3072_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha1 = "rsa_decrypt_oaep_4096_sha1"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha256 = "rsa_decrypt_oaep_4096_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha512 = "rsa_decrypt_oaep_4096_sha512"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs12048Sha256 = "rsa_sign_pkcs1_2048_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs13072Sha256 = "rsa_sign_pkcs1_3072_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs14096Sha256 = "rsa_sign_pkcs1_4096_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs14096Sha512 = "rsa_sign_pkcs1_4096_sha512"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignPss2048Sha256 = "rsa_sign_pss_2048_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignPss3072Sha256 = "rsa_sign_pss_3072_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignPss4096Sha256 = "rsa_sign_pss_4096_sha256"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignPss4096Sha512 = "rsa_sign_pss_4096_sha512"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs12048 = "rsa_sign_raw_pkcs1_2048"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs13072 = "rsa_sign_raw_pkcs1_3072"
	KeyInstanceGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs14096 = "rsa_sign_raw_pkcs1_4096"
)

func (*KeyInstanceGoogleKms) isaKeyInstance() bool {
	return true
}

// UnmarshalKeyInstanceGoogleKms unmarshals an instance of KeyInstanceGoogleKms from the specified map of raw messages.
func UnmarshalKeyInstanceGoogleKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyInstanceGoogleKms)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label_in_keystore", &obj.LabelInKeystore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "keystore", &obj.Keystore, UnmarshalInstanceInKeystore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_protection_level", &obj.GoogleKeyProtectionLevel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_purpose", &obj.GoogleKeyPurpose)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_kms_algorithm", &obj.GoogleKmsAlgorithm)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeyInstanceIbmCloudKms : The instance of a managed key for a specific keystore.
// This model "extends" KeyInstance
type KeyInstanceIbmCloudKms struct {
	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id" validate:"required"`

	// The label of the key.
	LabelInKeystore *string `json:"label_in_keystore" validate:"required"`

	// Type of the key instance.
	Type *string `json:"type,omitempty"`

	// Description of properties of a key within the context of keystores.
	Keystore *InstanceInKeystore `json:"keystore" validate:"required"`
}

// Constants associated with the KeyInstanceIbmCloudKms.Type property.
// Type of the key instance.
const (
	KeyInstanceIbmCloudKms_Type_KeyPair = "key_pair"
	KeyInstanceIbmCloudKms_Type_PrivateKey = "private_key"
	KeyInstanceIbmCloudKms_Type_PublicKey = "public_key"
	KeyInstanceIbmCloudKms_Type_SecretKey = "secret_key"
)

func (*KeyInstanceIbmCloudKms) isaKeyInstance() bool {
	return true
}

// UnmarshalKeyInstanceIbmCloudKms unmarshals an instance of KeyInstanceIbmCloudKms from the specified map of raw messages.
func UnmarshalKeyInstanceIbmCloudKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyInstanceIbmCloudKms)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label_in_keystore", &obj.LabelInKeystore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "keystore", &obj.Keystore, UnmarshalInstanceInKeystore)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequestKeystoreTypeAwsKmsCreate : AWS KMS is a managed service for you to create and manage cryptographic keys across a wide range of AWS services.
// This model "extends" KeystoreCreationRequest
type KeystoreCreationRequestKeystoreTypeAwsKmsCreate struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name" validate:"required"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// AWS Region.
	AwsRegion *string `json:"aws_region" validate:"required"`

	// The access key id used for connecting to this instance of AWS KMS.
	AwsAccessKeyID *string `json:"aws_access_key_id" validate:"required"`

	// The secret access key used for connecting to this instance of AWS KMS.
	AwsSecretAccessKey *string `json:"aws_secret_access_key" validate:"required"`
}

// Constants associated with the KeystoreCreationRequestKeystoreTypeAwsKmsCreate.Type property.
// Type of keystore.
const (
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_Type_AwsKms = "aws_kms"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_Type_GoogleKms = "google_kms"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreCreationRequestKeystoreTypeAwsKmsCreate.AwsRegion property.
// AWS Region.
const (
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_AfSouth1 = "af_south_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_ApEast1 = "ap_east_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_ApNortheast1 = "ap_northeast_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_ApNortheast2 = "ap_northeast_2"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_ApSouth1 = "ap_south_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_ApSoutheast1 = "ap_southeast_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_ApSoutheast2 = "ap_southeast_2"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_AwsCnGlobal = "aws_cn_global"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_AwsGlobal = "aws_global"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_AwsIsoBGlobal = "aws_iso_b_global"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_AwsIsoGlobal = "aws_iso_global"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_AwsUsGovGlobal = "aws_us_gov_global"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_CaCentral1 = "ca_central_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_CnNorth1 = "cn_north_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_CnNorthwest1 = "cn_northwest_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_EuCentral1 = "eu_central_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_EuWest1 = "eu_west_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_EuWest2 = "eu_west_2"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_EuWest3 = "eu_west_3"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_MeSouth1 = "me_south_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_SaEast1 = "sa_east_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_UsEast1 = "us_east_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_UsEast2 = "us_east_2"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_UsGovEast1 = "us_gov_east_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_UsGovWest1 = "us_gov_west_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_UsIsoEast1 = "us_iso_east_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_UsIsobEast1 = "us_isob_east_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_UsWest1 = "us_west_1"
	KeystoreCreationRequestKeystoreTypeAwsKmsCreate_AwsRegion_UsWest2 = "us_west_2"
)

// NewKeystoreCreationRequestKeystoreTypeAwsKmsCreate : Instantiate KeystoreCreationRequestKeystoreTypeAwsKmsCreate (Generic Model Constructor)
func (*UkoV4) NewKeystoreCreationRequestKeystoreTypeAwsKmsCreate(typeVar string, vault *VaultReferenceInCreationRequest, name string, awsRegion string, awsAccessKeyID string, awsSecretAccessKey string) (_model *KeystoreCreationRequestKeystoreTypeAwsKmsCreate, err error) {
	_model = &KeystoreCreationRequestKeystoreTypeAwsKmsCreate{
		Type: core.StringPtr(typeVar),
		Vault: vault,
		Name: core.StringPtr(name),
		AwsRegion: core.StringPtr(awsRegion),
		AwsAccessKeyID: core.StringPtr(awsAccessKeyID),
		AwsSecretAccessKey: core.StringPtr(awsSecretAccessKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KeystoreCreationRequestKeystoreTypeAwsKmsCreate) isaKeystoreCreationRequest() bool {
	return true
}

// UnmarshalKeystoreCreationRequestKeystoreTypeAwsKmsCreate unmarshals an instance of KeystoreCreationRequestKeystoreTypeAwsKmsCreate from the specified map of raw messages.
func UnmarshalKeystoreCreationRequestKeystoreTypeAwsKmsCreate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreCreationRequestKeystoreTypeAwsKmsCreate)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReferenceInCreationRequest)
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
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_access_key_id", &obj.AwsAccessKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_secret_access_key", &obj.AwsSecretAccessKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequestKeystoreTypeAzureCreate : Microsoft Azure Key Vault is a cloud service for you to create and manage cryptographic keys and other sensitive
// information.
// This model "extends" KeystoreCreationRequest
type KeystoreCreationRequestKeystoreTypeAzureCreate struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// Service name of the key vault instance from the Azure portal.
	AzureServiceName *string `json:"azure_service_name" validate:"required"`

	// Resource group in Azure.
	AzureResourceGroup *string `json:"azure_resource_group" validate:"required"`

	// Location of the Azure Key Vault.
	AzureLocation *string `json:"azure_location" validate:"required"`

	// Azure service principal client ID.
	AzureServicePrincipalClientID *string `json:"azure_service_principal_client_id" validate:"required"`

	// Azure service principal password.
	AzureServicePrincipalPassword *string `json:"azure_service_principal_password" validate:"required"`

	// Azure tenant that the Key Vault is associated with,.
	AzureTenant *string `json:"azure_tenant" validate:"required"`

	// Subscription ID in Azure.
	AzureSubscriptionID *string `json:"azure_subscription_id" validate:"required"`

	// Azure environment, usually 'Azure'.
	AzureEnvironment *string `json:"azure_environment" validate:"required"`
}

// Constants associated with the KeystoreCreationRequestKeystoreTypeAzureCreate.Type property.
// Type of keystore.
const (
	KeystoreCreationRequestKeystoreTypeAzureCreate_Type_AwsKms = "aws_kms"
	KeystoreCreationRequestKeystoreTypeAzureCreate_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequestKeystoreTypeAzureCreate_Type_GoogleKms = "google_kms"
	KeystoreCreationRequestKeystoreTypeAzureCreate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreCreationRequestKeystoreTypeAzureCreate.AzureLocation property.
// Location of the Azure Key Vault.
const (
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_AsiaEast = "asia_east"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_AsiaSoutheast = "asia_southeast"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_AustraliaCentral = "australia_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_AustraliaCentral2 = "australia_central_2"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_AustraliaEast = "australia_east"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_AustraliaSoutheast = "australia_southeast"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_BrazilSouth = "brazil_south"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_CanadaCentral = "canada_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_CanadaEast = "canada_east"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_ChinaEast = "china_east"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_ChinaEast2 = "china_east_2"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_ChinaNorth = "china_north"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_ChinaNorth2 = "china_north_2"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_EuropeNorth = "europe_north"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_EuropeWest = "europe_west"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_FranceCentral = "france_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_FranceSouth = "france_south"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_GermanyCentral = "germany_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_GermanyNortheast = "germany_northeast"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_IndiaCentral = "india_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_IndiaSouth = "india_south"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_IndiaWest = "india_west"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_JapanEast = "japan_east"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_JapanWest = "japan_west"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_KoreaCentral = "korea_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_KoreaSouth = "korea_south"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_SouthAfricaNorth = "south_africa_north"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_SouthAfricaWest = "south_africa_west"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UkSouth = "uk_south"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UkWest = "uk_west"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsCentral = "us_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsDodCentral = "us_dod_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsDodEast = "us_dod_east"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsEast = "us_east"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsEast2 = "us_east_2"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsGovArizona = "us_gov_arizona"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsGovIowa = "us_gov_iowa"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsGovTexas = "us_gov_texas"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsGovVirginia = "us_gov_virginia"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsNorthCentral = "us_north_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsSouthCentral = "us_south_central"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsWest = "us_west"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsWest2 = "us_west_2"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureLocation_UsWestCentral = "us_west_central"
)

// Constants associated with the KeystoreCreationRequestKeystoreTypeAzureCreate.AzureEnvironment property.
// Azure environment, usually 'Azure'.
const (
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureEnvironment_Azure = "azure"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureEnvironment_AzureChina = "azure_china"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureEnvironment_AzureGermany = "azure_germany"
	KeystoreCreationRequestKeystoreTypeAzureCreate_AzureEnvironment_AzureUsGovernment = "azure_us_government"
)

// NewKeystoreCreationRequestKeystoreTypeAzureCreate : Instantiate KeystoreCreationRequestKeystoreTypeAzureCreate (Generic Model Constructor)
func (*UkoV4) NewKeystoreCreationRequestKeystoreTypeAzureCreate(typeVar string, vault *VaultReferenceInCreationRequest, azureServiceName string, azureResourceGroup string, azureLocation string, azureServicePrincipalClientID string, azureServicePrincipalPassword string, azureTenant string, azureSubscriptionID string, azureEnvironment string) (_model *KeystoreCreationRequestKeystoreTypeAzureCreate, err error) {
	_model = &KeystoreCreationRequestKeystoreTypeAzureCreate{
		Type: core.StringPtr(typeVar),
		Vault: vault,
		AzureServiceName: core.StringPtr(azureServiceName),
		AzureResourceGroup: core.StringPtr(azureResourceGroup),
		AzureLocation: core.StringPtr(azureLocation),
		AzureServicePrincipalClientID: core.StringPtr(azureServicePrincipalClientID),
		AzureServicePrincipalPassword: core.StringPtr(azureServicePrincipalPassword),
		AzureTenant: core.StringPtr(azureTenant),
		AzureSubscriptionID: core.StringPtr(azureSubscriptionID),
		AzureEnvironment: core.StringPtr(azureEnvironment),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KeystoreCreationRequestKeystoreTypeAzureCreate) isaKeystoreCreationRequest() bool {
	return true
}

// UnmarshalKeystoreCreationRequestKeystoreTypeAzureCreate unmarshals an instance of KeystoreCreationRequestKeystoreTypeAzureCreate from the specified map of raw messages.
func UnmarshalKeystoreCreationRequestKeystoreTypeAzureCreate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreCreationRequestKeystoreTypeAzureCreate)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReferenceInCreationRequest)
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
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_name", &obj.AzureServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_location", &obj.AzureLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_client_id", &obj.AzureServicePrincipalClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_password", &obj.AzureServicePrincipalPassword)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_tenant", &obj.AzureTenant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_environment", &obj.AzureEnvironment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequestKeystoreTypeGoogleKmsCreate : Google Cloud KMS is a managed service for you to create and manage cryptographic keys across a wide range of Google
// Cloud services.
// This model "extends" KeystoreCreationRequest
type KeystoreCreationRequestKeystoreTypeGoogleKmsCreate struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name" validate:"required"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// The value of the JSON key represented in the Base64 format.
	GoogleCredentials *string `json:"google_credentials" validate:"required"`

	// Location represents the geographical region where a Cloud KMS resource is stored and can be accessed. A key's
	// location impacts the performance of applications using the key.
	GoogleLocation *string `json:"google_location,omitempty"`

	// The project id associated with this keystore.
	GoogleProjectID *string `json:"google_project_id,omitempty"`

	// The private key id associated with this keystore.
	GooglePrivateKeyID *string `json:"google_private_key_id,omitempty"`

	// A key ring organizes keys in a specific Google Cloud location and allows you to manage access control on groups of
	// keys.
	GoogleKeyRing *string `json:"google_key_ring,omitempty"`
}

// Constants associated with the KeystoreCreationRequestKeystoreTypeGoogleKmsCreate.Type property.
// Type of keystore.
const (
	KeystoreCreationRequestKeystoreTypeGoogleKmsCreate_Type_AwsKms = "aws_kms"
	KeystoreCreationRequestKeystoreTypeGoogleKmsCreate_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequestKeystoreTypeGoogleKmsCreate_Type_GoogleKms = "google_kms"
	KeystoreCreationRequestKeystoreTypeGoogleKmsCreate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// NewKeystoreCreationRequestKeystoreTypeGoogleKmsCreate : Instantiate KeystoreCreationRequestKeystoreTypeGoogleKmsCreate (Generic Model Constructor)
func (*UkoV4) NewKeystoreCreationRequestKeystoreTypeGoogleKmsCreate(typeVar string, vault *VaultReferenceInCreationRequest, name string, googleCredentials string) (_model *KeystoreCreationRequestKeystoreTypeGoogleKmsCreate, err error) {
	_model = &KeystoreCreationRequestKeystoreTypeGoogleKmsCreate{
		Type: core.StringPtr(typeVar),
		Vault: vault,
		Name: core.StringPtr(name),
		GoogleCredentials: core.StringPtr(googleCredentials),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KeystoreCreationRequestKeystoreTypeGoogleKmsCreate) isaKeystoreCreationRequest() bool {
	return true
}

// UnmarshalKeystoreCreationRequestKeystoreTypeGoogleKmsCreate unmarshals an instance of KeystoreCreationRequestKeystoreTypeGoogleKmsCreate from the specified map of raw messages.
func UnmarshalKeystoreCreationRequestKeystoreTypeGoogleKmsCreate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreCreationRequestKeystoreTypeGoogleKmsCreate)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReferenceInCreationRequest)
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
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_credentials", &obj.GoogleCredentials)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_location", &obj.GoogleLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_project_id", &obj.GoogleProjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_private_key_id", &obj.GooglePrivateKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_ring", &obj.GoogleKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate : Properties required to create an IBM Cloud keystore.
// Models which "extend" this model:
// - KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate
// - KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate
// This model "extends" KeystoreCreationRequest
type KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Possible IBM Cloud KMS variants.
	IbmVariant *string `json:"ibm_variant" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// API endpoint of the IBM Cloud keystore.
	IbmApiEndpoint *string `json:"ibm_api_endpoint,omitempty"`

	// Endpoint of the IAM service for this IBM Cloud keystore.
	IbmIamEndpoint *string `json:"ibm_iam_endpoint,omitempty"`

	// The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.
	IbmApiKey *string `json:"ibm_api_key,omitempty"`

	// The instance ID of the IBM Cloud keystore.
	IbmInstanceID *string `json:"ibm_instance_id,omitempty"`

	// The key ring of an IBM Cloud KMS Keystore.
	IbmKeyRing *string `json:"ibm_key_ring,omitempty"`
}

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate.Type property.
// Type of keystore.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate_Type_AwsKms = "aws_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate_Type_GoogleKms = "google_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate.IbmVariant property.
// Possible IBM Cloud KMS variants.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate_IbmVariant_Hpcs = "hpcs"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate_IbmVariant_Internal = "internal"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate_IbmVariant_KeyProtect = "key_protect"
)
func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate() bool {
	return true
}

type KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateIntf interface {
	KeystoreCreationRequestIntf
	isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate() bool
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate) isaKeystoreCreationRequest() bool {
	return true
}

// UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate unmarshals an instance of KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate from the specified map of raw messages.
func UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReferenceInCreationRequest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_variant", &obj.IbmVariant)
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
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_endpoint", &obj.IbmApiEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_iam_endpoint", &obj.IbmIamEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IbmApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_instance_id", &obj.IbmInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_key_ring", &obj.IbmKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreTypeAwsKms : Properties related to AWS KMS.
// This model "extends" Keystore
type KeystoreTypeAwsKms struct {
	// Reference to a vault.
	Vault *VaultReference `json:"vault,omitempty"`

	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the target keystore. It can be changed in the future.
	Name *string `json:"name" validate:"required"`

	// Geographic location of the keystore, if available.
	Location *string `json:"location" validate:"required"`

	// Description of the keystore.
	Description *string `json:"description" validate:"required"`

	// List of groups that this keystore belongs to.
	Groups []string `json:"groups" validate:"required"`

	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	// Date and time when the target keystore was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Date and time when the target keystore was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// ID of the user that created the key.
	CreatedBy *string `json:"created_by,omitempty"`

	// ID of the user that last updated the key.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`

	// AWS Region.
	AwsRegion *string `json:"aws_region" validate:"required"`

	// The access key id used for connecting to this instance of AWS KMS.
	AwsAccessKeyID *string `json:"aws_access_key_id" validate:"required"`

	// The secret access key used for connecting to this instance of AWS KMS.
	AwsSecretAccessKey *string `json:"aws_secret_access_key" validate:"required"`
}

// Constants associated with the KeystoreTypeAwsKms.Type property.
// Type of keystore.
const (
	KeystoreTypeAwsKms_Type_AwsKms = "aws_kms"
	KeystoreTypeAwsKms_Type_AzureKeyVault = "azure_key_vault"
	KeystoreTypeAwsKms_Type_GoogleKms = "google_kms"
	KeystoreTypeAwsKms_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreTypeAwsKms.AwsRegion property.
// AWS Region.
const (
	KeystoreTypeAwsKms_AwsRegion_AfSouth1 = "af_south_1"
	KeystoreTypeAwsKms_AwsRegion_ApEast1 = "ap_east_1"
	KeystoreTypeAwsKms_AwsRegion_ApNortheast1 = "ap_northeast_1"
	KeystoreTypeAwsKms_AwsRegion_ApNortheast2 = "ap_northeast_2"
	KeystoreTypeAwsKms_AwsRegion_ApSouth1 = "ap_south_1"
	KeystoreTypeAwsKms_AwsRegion_ApSoutheast1 = "ap_southeast_1"
	KeystoreTypeAwsKms_AwsRegion_ApSoutheast2 = "ap_southeast_2"
	KeystoreTypeAwsKms_AwsRegion_AwsCnGlobal = "aws_cn_global"
	KeystoreTypeAwsKms_AwsRegion_AwsGlobal = "aws_global"
	KeystoreTypeAwsKms_AwsRegion_AwsIsoBGlobal = "aws_iso_b_global"
	KeystoreTypeAwsKms_AwsRegion_AwsIsoGlobal = "aws_iso_global"
	KeystoreTypeAwsKms_AwsRegion_AwsUsGovGlobal = "aws_us_gov_global"
	KeystoreTypeAwsKms_AwsRegion_CaCentral1 = "ca_central_1"
	KeystoreTypeAwsKms_AwsRegion_CnNorth1 = "cn_north_1"
	KeystoreTypeAwsKms_AwsRegion_CnNorthwest1 = "cn_northwest_1"
	KeystoreTypeAwsKms_AwsRegion_EuCentral1 = "eu_central_1"
	KeystoreTypeAwsKms_AwsRegion_EuWest1 = "eu_west_1"
	KeystoreTypeAwsKms_AwsRegion_EuWest2 = "eu_west_2"
	KeystoreTypeAwsKms_AwsRegion_EuWest3 = "eu_west_3"
	KeystoreTypeAwsKms_AwsRegion_MeSouth1 = "me_south_1"
	KeystoreTypeAwsKms_AwsRegion_SaEast1 = "sa_east_1"
	KeystoreTypeAwsKms_AwsRegion_UsEast1 = "us_east_1"
	KeystoreTypeAwsKms_AwsRegion_UsEast2 = "us_east_2"
	KeystoreTypeAwsKms_AwsRegion_UsGovEast1 = "us_gov_east_1"
	KeystoreTypeAwsKms_AwsRegion_UsGovWest1 = "us_gov_west_1"
	KeystoreTypeAwsKms_AwsRegion_UsIsoEast1 = "us_iso_east_1"
	KeystoreTypeAwsKms_AwsRegion_UsIsobEast1 = "us_isob_east_1"
	KeystoreTypeAwsKms_AwsRegion_UsWest1 = "us_west_1"
	KeystoreTypeAwsKms_AwsRegion_UsWest2 = "us_west_2"
)

func (*KeystoreTypeAwsKms) isaKeystore() bool {
	return true
}

// UnmarshalKeystoreTypeAwsKms unmarshals an instance of KeystoreTypeAwsKms from the specified map of raw messages.
func UnmarshalKeystoreTypeAwsKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreTypeAwsKms)
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReference)
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
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_access_key_id", &obj.AwsAccessKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_secret_access_key", &obj.AwsSecretAccessKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreTypeAzure : Properties related to Azure Key Vaults.
// This model "extends" Keystore
type KeystoreTypeAzure struct {
	// Reference to a vault.
	Vault *VaultReference `json:"vault,omitempty"`

	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the target keystore. It can be changed in the future.
	Name *string `json:"name" validate:"required"`

	// Geographic location of the keystore, if available.
	Location *string `json:"location" validate:"required"`

	// Description of the keystore.
	Description *string `json:"description" validate:"required"`

	// List of groups that this keystore belongs to.
	Groups []string `json:"groups" validate:"required"`

	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	// Date and time when the target keystore was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Date and time when the target keystore was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// ID of the user that created the key.
	CreatedBy *string `json:"created_by,omitempty"`

	// ID of the user that last updated the key.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`

	// Service name of the key vault instance from the Azure portal.
	AzureServiceName *string `json:"azure_service_name" validate:"required"`

	// Resource group in Azure.
	AzureResourceGroup *string `json:"azure_resource_group" validate:"required"`

	// Location of the Azure Key Vault.
	AzureLocation *string `json:"azure_location" validate:"required"`

	// Azure service principal client ID.
	AzureServicePrincipalClientID *string `json:"azure_service_principal_client_id" validate:"required"`

	// Azure service principal password.
	AzureServicePrincipalPassword *string `json:"azure_service_principal_password" validate:"required"`

	// Azure tenant that the Key Vault is associated with,.
	AzureTenant *string `json:"azure_tenant" validate:"required"`

	// Subscription ID in Azure.
	AzureSubscriptionID *string `json:"azure_subscription_id" validate:"required"`

	// Azure environment, usually 'Azure'.
	AzureEnvironment *string `json:"azure_environment" validate:"required"`
}

// Constants associated with the KeystoreTypeAzure.Type property.
// Type of keystore.
const (
	KeystoreTypeAzure_Type_AwsKms = "aws_kms"
	KeystoreTypeAzure_Type_AzureKeyVault = "azure_key_vault"
	KeystoreTypeAzure_Type_GoogleKms = "google_kms"
	KeystoreTypeAzure_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreTypeAzure.AzureLocation property.
// Location of the Azure Key Vault.
const (
	KeystoreTypeAzure_AzureLocation_AsiaEast = "asia_east"
	KeystoreTypeAzure_AzureLocation_AsiaSoutheast = "asia_southeast"
	KeystoreTypeAzure_AzureLocation_AustraliaCentral = "australia_central"
	KeystoreTypeAzure_AzureLocation_AustraliaCentral2 = "australia_central_2"
	KeystoreTypeAzure_AzureLocation_AustraliaEast = "australia_east"
	KeystoreTypeAzure_AzureLocation_AustraliaSoutheast = "australia_southeast"
	KeystoreTypeAzure_AzureLocation_BrazilSouth = "brazil_south"
	KeystoreTypeAzure_AzureLocation_CanadaCentral = "canada_central"
	KeystoreTypeAzure_AzureLocation_CanadaEast = "canada_east"
	KeystoreTypeAzure_AzureLocation_ChinaEast = "china_east"
	KeystoreTypeAzure_AzureLocation_ChinaEast2 = "china_east_2"
	KeystoreTypeAzure_AzureLocation_ChinaNorth = "china_north"
	KeystoreTypeAzure_AzureLocation_ChinaNorth2 = "china_north_2"
	KeystoreTypeAzure_AzureLocation_EuropeNorth = "europe_north"
	KeystoreTypeAzure_AzureLocation_EuropeWest = "europe_west"
	KeystoreTypeAzure_AzureLocation_FranceCentral = "france_central"
	KeystoreTypeAzure_AzureLocation_FranceSouth = "france_south"
	KeystoreTypeAzure_AzureLocation_GermanyCentral = "germany_central"
	KeystoreTypeAzure_AzureLocation_GermanyNortheast = "germany_northeast"
	KeystoreTypeAzure_AzureLocation_IndiaCentral = "india_central"
	KeystoreTypeAzure_AzureLocation_IndiaSouth = "india_south"
	KeystoreTypeAzure_AzureLocation_IndiaWest = "india_west"
	KeystoreTypeAzure_AzureLocation_JapanEast = "japan_east"
	KeystoreTypeAzure_AzureLocation_JapanWest = "japan_west"
	KeystoreTypeAzure_AzureLocation_KoreaCentral = "korea_central"
	KeystoreTypeAzure_AzureLocation_KoreaSouth = "korea_south"
	KeystoreTypeAzure_AzureLocation_SouthAfricaNorth = "south_africa_north"
	KeystoreTypeAzure_AzureLocation_SouthAfricaWest = "south_africa_west"
	KeystoreTypeAzure_AzureLocation_UkSouth = "uk_south"
	KeystoreTypeAzure_AzureLocation_UkWest = "uk_west"
	KeystoreTypeAzure_AzureLocation_UsCentral = "us_central"
	KeystoreTypeAzure_AzureLocation_UsDodCentral = "us_dod_central"
	KeystoreTypeAzure_AzureLocation_UsDodEast = "us_dod_east"
	KeystoreTypeAzure_AzureLocation_UsEast = "us_east"
	KeystoreTypeAzure_AzureLocation_UsEast2 = "us_east_2"
	KeystoreTypeAzure_AzureLocation_UsGovArizona = "us_gov_arizona"
	KeystoreTypeAzure_AzureLocation_UsGovIowa = "us_gov_iowa"
	KeystoreTypeAzure_AzureLocation_UsGovTexas = "us_gov_texas"
	KeystoreTypeAzure_AzureLocation_UsGovVirginia = "us_gov_virginia"
	KeystoreTypeAzure_AzureLocation_UsNorthCentral = "us_north_central"
	KeystoreTypeAzure_AzureLocation_UsSouthCentral = "us_south_central"
	KeystoreTypeAzure_AzureLocation_UsWest = "us_west"
	KeystoreTypeAzure_AzureLocation_UsWest2 = "us_west_2"
	KeystoreTypeAzure_AzureLocation_UsWestCentral = "us_west_central"
)

// Constants associated with the KeystoreTypeAzure.AzureEnvironment property.
// Azure environment, usually 'Azure'.
const (
	KeystoreTypeAzure_AzureEnvironment_Azure = "azure"
	KeystoreTypeAzure_AzureEnvironment_AzureChina = "azure_china"
	KeystoreTypeAzure_AzureEnvironment_AzureGermany = "azure_germany"
	KeystoreTypeAzure_AzureEnvironment_AzureUsGovernment = "azure_us_government"
)

func (*KeystoreTypeAzure) isaKeystore() bool {
	return true
}

// UnmarshalKeystoreTypeAzure unmarshals an instance of KeystoreTypeAzure from the specified map of raw messages.
func UnmarshalKeystoreTypeAzure(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreTypeAzure)
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReference)
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
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_name", &obj.AzureServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_location", &obj.AzureLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_client_id", &obj.AzureServicePrincipalClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_password", &obj.AzureServicePrincipalPassword)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_tenant", &obj.AzureTenant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_environment", &obj.AzureEnvironment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreTypeGoogleKms : Properties related to Google Cloud KMS.
// This model "extends" Keystore
type KeystoreTypeGoogleKms struct {
	// Reference to a vault.
	Vault *VaultReference `json:"vault,omitempty"`

	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the target keystore. It can be changed in the future.
	Name *string `json:"name" validate:"required"`

	// Geographic location of the keystore, if available.
	Location *string `json:"location" validate:"required"`

	// Description of the keystore.
	Description *string `json:"description" validate:"required"`

	// List of groups that this keystore belongs to.
	Groups []string `json:"groups" validate:"required"`

	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	// Date and time when the target keystore was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Date and time when the target keystore was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// ID of the user that created the key.
	CreatedBy *string `json:"created_by,omitempty"`

	// ID of the user that last updated the key.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`

	// The value of the JSON key represented in the Base64 format.
	GoogleCredentials *string `json:"google_credentials" validate:"required"`

	// Location represents the geographical region where a Cloud KMS resource is stored and can be accessed. A key's
	// location impacts the performance of applications using the key.
	GoogleLocation *string `json:"google_location" validate:"required"`

	// The project id associated with this keystore.
	GoogleProjectID *string `json:"google_project_id,omitempty"`

	// The private key id associated with this keystore.
	GooglePrivateKeyID *string `json:"google_private_key_id,omitempty"`

	// A key ring organizes keys in a specific Google Cloud location and allows you to manage access control on groups of
	// keys.
	GoogleKeyRing *string `json:"google_key_ring" validate:"required"`
}

// Constants associated with the KeystoreTypeGoogleKms.Type property.
// Type of keystore.
const (
	KeystoreTypeGoogleKms_Type_AwsKms = "aws_kms"
	KeystoreTypeGoogleKms_Type_AzureKeyVault = "azure_key_vault"
	KeystoreTypeGoogleKms_Type_GoogleKms = "google_kms"
	KeystoreTypeGoogleKms_Type_IbmCloudKms = "ibm_cloud_kms"
)

func (*KeystoreTypeGoogleKms) isaKeystore() bool {
	return true
}

// UnmarshalKeystoreTypeGoogleKms unmarshals an instance of KeystoreTypeGoogleKms from the specified map of raw messages.
func UnmarshalKeystoreTypeGoogleKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreTypeGoogleKms)
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReference)
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
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_credentials", &obj.GoogleCredentials)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_location", &obj.GoogleLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_project_id", &obj.GoogleProjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_private_key_id", &obj.GooglePrivateKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_ring", &obj.GoogleKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreTypeIbmCloudKms : properties of a ibm cloud keystore.
// This model "extends" Keystore
type KeystoreTypeIbmCloudKms struct {
	// Reference to a vault.
	Vault *VaultReference `json:"vault,omitempty"`

	// The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.
	ID *string `json:"id,omitempty"`

	// Name of the target keystore. It can be changed in the future.
	Name *string `json:"name" validate:"required"`

	// Geographic location of the keystore, if available.
	Location *string `json:"location" validate:"required"`

	// Description of the keystore.
	Description *string `json:"description" validate:"required"`

	// List of groups that this keystore belongs to.
	Groups []string `json:"groups" validate:"required"`

	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	// Date and time when the target keystore was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Date and time when the target keystore was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// ID of the user that created the key.
	CreatedBy *string `json:"created_by,omitempty"`

	// ID of the user that last updated the key.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// A URL that uniquely identifies your cloud resource.
	Href *string `json:"href,omitempty"`

	// API endpoint of the IBM Cloud keystore.
	IbmApiEndpoint *string `json:"ibm_api_endpoint" validate:"required"`

	// Endpoint of the IAM service for this IBM Cloud keystore.
	IbmIamEndpoint *string `json:"ibm_iam_endpoint" validate:"required"`

	// The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.
	IbmApiKey *string `json:"ibm_api_key" validate:"required"`

	// The instance ID of the IBM Cloud keystore.
	IbmInstanceID *string `json:"ibm_instance_id" validate:"required"`

	// Possible IBM Cloud KMS variants.
	IbmVariant *string `json:"ibm_variant" validate:"required"`

	// The key ring of an IBM Cloud KMS Keystore.
	IbmKeyRing *string `json:"ibm_key_ring,omitempty"`
}

// Constants associated with the KeystoreTypeIbmCloudKms.Type property.
// Type of keystore.
const (
	KeystoreTypeIbmCloudKms_Type_AwsKms = "aws_kms"
	KeystoreTypeIbmCloudKms_Type_AzureKeyVault = "azure_key_vault"
	KeystoreTypeIbmCloudKms_Type_GoogleKms = "google_kms"
	KeystoreTypeIbmCloudKms_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreTypeIbmCloudKms.IbmVariant property.
// Possible IBM Cloud KMS variants.
const (
	KeystoreTypeIbmCloudKms_IbmVariant_Hpcs = "hpcs"
	KeystoreTypeIbmCloudKms_IbmVariant_Internal = "internal"
	KeystoreTypeIbmCloudKms_IbmVariant_KeyProtect = "key_protect"
)

func (*KeystoreTypeIbmCloudKms) isaKeystore() bool {
	return true
}

// UnmarshalKeystoreTypeIbmCloudKms unmarshals an instance of KeystoreTypeIbmCloudKms from the specified map of raw messages.
func UnmarshalKeystoreTypeIbmCloudKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreTypeIbmCloudKms)
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReference)
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
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_endpoint", &obj.IbmApiEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_iam_endpoint", &obj.IbmIamEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IbmApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_instance_id", &obj.IbmInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_variant", &obj.IbmVariant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_key_ring", &obj.IbmKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate : AWS KMS is a managed service for you to create and manage cryptographic keys across a wide range of AWS services.
// This model "extends" KeystoreUpdateRequest
type KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate struct {
	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// AWS Region.
	AwsRegion *string `json:"aws_region,omitempty"`

	// The access key id used for connecting to this instance of AWS KMS.
	AwsAccessKeyID *string `json:"aws_access_key_id,omitempty"`

	// The secret access key used for connecting to this instance of AWS KMS.
	AwsSecretAccessKey *string `json:"aws_secret_access_key,omitempty"`
}

// Constants associated with the KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate.AwsRegion property.
// AWS Region.
const (
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_AfSouth1 = "af_south_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_ApEast1 = "ap_east_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_ApNortheast1 = "ap_northeast_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_ApNortheast2 = "ap_northeast_2"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_ApSouth1 = "ap_south_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_ApSoutheast1 = "ap_southeast_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_ApSoutheast2 = "ap_southeast_2"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_AwsCnGlobal = "aws_cn_global"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_AwsGlobal = "aws_global"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_AwsIsoBGlobal = "aws_iso_b_global"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_AwsIsoGlobal = "aws_iso_global"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_AwsUsGovGlobal = "aws_us_gov_global"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_CaCentral1 = "ca_central_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_CnNorth1 = "cn_north_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_CnNorthwest1 = "cn_northwest_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_EuCentral1 = "eu_central_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_EuWest1 = "eu_west_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_EuWest2 = "eu_west_2"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_EuWest3 = "eu_west_3"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_MeSouth1 = "me_south_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_SaEast1 = "sa_east_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_UsEast1 = "us_east_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_UsEast2 = "us_east_2"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_UsGovEast1 = "us_gov_east_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_UsGovWest1 = "us_gov_west_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_UsIsoEast1 = "us_iso_east_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_UsIsobEast1 = "us_isob_east_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_UsWest1 = "us_west_1"
	KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate_AwsRegion_UsWest2 = "us_west_2"
)

func (*KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate) isaKeystoreUpdateRequest() bool {
	return true
}

// UnmarshalKeystoreUpdateRequestKeystoreTypeAwsKmsUpdate unmarshals an instance of KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate from the specified map of raw messages.
func UnmarshalKeystoreUpdateRequestKeystoreTypeAwsKmsUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreUpdateRequestKeystoreTypeAwsKmsUpdate)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_region", &obj.AwsRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_access_key_id", &obj.AwsAccessKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aws_secret_access_key", &obj.AwsSecretAccessKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreUpdateRequestKeystoreTypeAzureUpdate : Microsoft Azure Key Vault is a cloud service for you to create and manage cryptographic keys and other sensitive
// information.
// This model "extends" KeystoreUpdateRequest
type KeystoreUpdateRequestKeystoreTypeAzureUpdate struct {
	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// Service name of the key vault instance from the Azure portal.
	AzureServiceName *string `json:"azure_service_name,omitempty"`

	// Resource group in Azure.
	AzureResourceGroup *string `json:"azure_resource_group,omitempty"`

	// Location of the Azure Key Vault.
	AzureLocation *string `json:"azure_location,omitempty"`

	// Azure service principal client ID.
	AzureServicePrincipalClientID *string `json:"azure_service_principal_client_id,omitempty"`

	// Azure service principal password.
	AzureServicePrincipalPassword *string `json:"azure_service_principal_password,omitempty"`

	// Azure tenant that the Key Vault is associated with,.
	AzureTenant *string `json:"azure_tenant,omitempty"`

	// Subscription ID in Azure.
	AzureSubscriptionID *string `json:"azure_subscription_id,omitempty"`

	// Azure environment, usually 'Azure'.
	AzureEnvironment *string `json:"azure_environment,omitempty"`
}

// Constants associated with the KeystoreUpdateRequestKeystoreTypeAzureUpdate.AzureLocation property.
// Location of the Azure Key Vault.
const (
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_AsiaEast = "asia_east"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_AsiaSoutheast = "asia_southeast"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_AustraliaCentral = "australia_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_AustraliaCentral2 = "australia_central_2"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_AustraliaEast = "australia_east"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_AustraliaSoutheast = "australia_southeast"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_BrazilSouth = "brazil_south"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_CanadaCentral = "canada_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_CanadaEast = "canada_east"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_ChinaEast = "china_east"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_ChinaEast2 = "china_east_2"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_ChinaNorth = "china_north"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_ChinaNorth2 = "china_north_2"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_EuropeNorth = "europe_north"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_EuropeWest = "europe_west"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_FranceCentral = "france_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_FranceSouth = "france_south"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_GermanyCentral = "germany_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_GermanyNortheast = "germany_northeast"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_IndiaCentral = "india_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_IndiaSouth = "india_south"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_IndiaWest = "india_west"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_JapanEast = "japan_east"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_JapanWest = "japan_west"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_KoreaCentral = "korea_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_KoreaSouth = "korea_south"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_SouthAfricaNorth = "south_africa_north"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_SouthAfricaWest = "south_africa_west"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UkSouth = "uk_south"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UkWest = "uk_west"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsCentral = "us_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsDodCentral = "us_dod_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsDodEast = "us_dod_east"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsEast = "us_east"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsEast2 = "us_east_2"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsGovArizona = "us_gov_arizona"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsGovIowa = "us_gov_iowa"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsGovTexas = "us_gov_texas"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsGovVirginia = "us_gov_virginia"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsNorthCentral = "us_north_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsSouthCentral = "us_south_central"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsWest = "us_west"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsWest2 = "us_west_2"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureLocation_UsWestCentral = "us_west_central"
)

// Constants associated with the KeystoreUpdateRequestKeystoreTypeAzureUpdate.AzureEnvironment property.
// Azure environment, usually 'Azure'.
const (
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureEnvironment_Azure = "azure"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureEnvironment_AzureChina = "azure_china"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureEnvironment_AzureGermany = "azure_germany"
	KeystoreUpdateRequestKeystoreTypeAzureUpdate_AzureEnvironment_AzureUsGovernment = "azure_us_government"
)

func (*KeystoreUpdateRequestKeystoreTypeAzureUpdate) isaKeystoreUpdateRequest() bool {
	return true
}

// UnmarshalKeystoreUpdateRequestKeystoreTypeAzureUpdate unmarshals an instance of KeystoreUpdateRequestKeystoreTypeAzureUpdate from the specified map of raw messages.
func UnmarshalKeystoreUpdateRequestKeystoreTypeAzureUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreUpdateRequestKeystoreTypeAzureUpdate)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_name", &obj.AzureServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_resource_group", &obj.AzureResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_location", &obj.AzureLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_client_id", &obj.AzureServicePrincipalClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_service_principal_password", &obj.AzureServicePrincipalPassword)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_tenant", &obj.AzureTenant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_subscription_id", &obj.AzureSubscriptionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "azure_environment", &obj.AzureEnvironment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreUpdateRequestKeystoreTypeGoogleKmsUpdate : Google Cloud KMS is a managed service for you to create and manage cryptographic keys across a wide range of Google
// Cloud services.
// This model "extends" KeystoreUpdateRequest
type KeystoreUpdateRequestKeystoreTypeGoogleKmsUpdate struct {
	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// The value of the JSON key represented in the Base64 format.
	GoogleCredentials *string `json:"google_credentials,omitempty"`

	// Location represents the geographical region where a Cloud KMS resource is stored and can be accessed. A key's
	// location impacts the performance of applications using the key.
	GoogleLocation *string `json:"google_location,omitempty"`

	// The project id associated with this keystore.
	GoogleProjectID *string `json:"google_project_id,omitempty"`

	// The private key id associated with this keystore.
	GooglePrivateKeyID *string `json:"google_private_key_id,omitempty"`

	// A key ring organizes keys in a specific Google Cloud location and allows you to manage access control on groups of
	// keys.
	GoogleKeyRing *string `json:"google_key_ring,omitempty"`
}

func (*KeystoreUpdateRequestKeystoreTypeGoogleKmsUpdate) isaKeystoreUpdateRequest() bool {
	return true
}

// UnmarshalKeystoreUpdateRequestKeystoreTypeGoogleKmsUpdate unmarshals an instance of KeystoreUpdateRequestKeystoreTypeGoogleKmsUpdate from the specified map of raw messages.
func UnmarshalKeystoreUpdateRequestKeystoreTypeGoogleKmsUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreUpdateRequestKeystoreTypeGoogleKmsUpdate)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_credentials", &obj.GoogleCredentials)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_location", &obj.GoogleLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_project_id", &obj.GoogleProjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_private_key_id", &obj.GooglePrivateKeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_ring", &obj.GoogleKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate : KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate struct
// Models which "extend" this model:
// - KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate
// This model "extends" KeystoreUpdateRequest
type KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate struct {
	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`
}
func (*KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate) isaKeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate() bool {
	return true
}

type KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateIntf interface {
	KeystoreUpdateRequestIntf
	isaKeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate() bool
}

func (*KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate) isaKeystoreUpdateRequest() bool {
	return true
}

// UnmarshalKeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate unmarshals an instance of KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate from the specified map of raw messages.
func UnmarshalKeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreUpdateRequestKeystoreTypeIbmCloudKmsUpdate : You can connect your Hyper Protect Crypto Services instance to the keystores of another Hyper Protect Crypto Services
// or Key Protect instance, and manage its KMS keys using the current service instance.
// This model "extends" KeystoreUpdateRequest
type KeystoreUpdateRequestKeystoreTypeIbmCloudKmsUpdate struct {
	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// API endpoint of the IBM Cloud keystore.
	IbmApiEndpoint *string `json:"ibm_api_endpoint,omitempty"`

	// Endpoint of the IAM service for this IBM Cloud keystore.
	IbmIamEndpoint *string `json:"ibm_iam_endpoint,omitempty"`

	// The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.
	IbmApiKey *string `json:"ibm_api_key,omitempty"`

	// The instance ID of the IBM Cloud keystore.
	IbmInstanceID *string `json:"ibm_instance_id,omitempty"`

	// The key ring of an IBM Cloud KMS Keystore.
	IbmKeyRing *string `json:"ibm_key_ring,omitempty"`
}

func (*KeystoreUpdateRequestKeystoreTypeIbmCloudKmsUpdate) isaKeystoreUpdateRequest() bool {
	return true
}

// UnmarshalKeystoreUpdateRequestKeystoreTypeIbmCloudKmsUpdate unmarshals an instance of KeystoreUpdateRequestKeystoreTypeIbmCloudKmsUpdate from the specified map of raw messages.
func UnmarshalKeystoreUpdateRequestKeystoreTypeIbmCloudKmsUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreUpdateRequestKeystoreTypeIbmCloudKmsUpdate)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_endpoint", &obj.IbmApiEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_iam_endpoint", &obj.IbmIamEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IbmApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_instance_id", &obj.IbmInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_key_ring", &obj.IbmKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoresPropertiesCreateAwsKms : KeystoresPropertiesCreateAwsKms struct
// This model "extends" KeystoresPropertiesCreate
type KeystoresPropertiesCreateAwsKms struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`

	// Type of keystore.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the KeystoresPropertiesCreateAwsKms.Type property.
// Type of keystore.
const (
	KeystoresPropertiesCreateAwsKms_Type_AwsKms = "aws_kms"
	KeystoresPropertiesCreateAwsKms_Type_AzureKeyVault = "azure_key_vault"
	KeystoresPropertiesCreateAwsKms_Type_GoogleKms = "google_kms"
	KeystoresPropertiesCreateAwsKms_Type_IbmCloudKms = "ibm_cloud_kms"
)

func (*KeystoresPropertiesCreateAwsKms) isaKeystoresPropertiesCreate() bool {
	return true
}

// UnmarshalKeystoresPropertiesCreateAwsKms unmarshals an instance of KeystoresPropertiesCreateAwsKms from the specified map of raw messages.
func UnmarshalKeystoresPropertiesCreateAwsKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesCreateAwsKms)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
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

// KeystoresPropertiesCreateAzure : KeystoresPropertiesCreateAzure struct
// This model "extends" KeystoresPropertiesCreate
type KeystoresPropertiesCreateAzure struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`

	// Type of keystore.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the KeystoresPropertiesCreateAzure.Type property.
// Type of keystore.
const (
	KeystoresPropertiesCreateAzure_Type_AwsKms = "aws_kms"
	KeystoresPropertiesCreateAzure_Type_AzureKeyVault = "azure_key_vault"
	KeystoresPropertiesCreateAzure_Type_GoogleKms = "google_kms"
	KeystoresPropertiesCreateAzure_Type_IbmCloudKms = "ibm_cloud_kms"
)

func (*KeystoresPropertiesCreateAzure) isaKeystoresPropertiesCreate() bool {
	return true
}

// UnmarshalKeystoresPropertiesCreateAzure unmarshals an instance of KeystoresPropertiesCreateAzure from the specified map of raw messages.
func UnmarshalKeystoresPropertiesCreateAzure(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesCreateAzure)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
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

// KeystoresPropertiesCreateGoogleKms : KeystoresPropertiesCreateGoogleKms struct
// This model "extends" KeystoresPropertiesCreate
type KeystoresPropertiesCreateGoogleKms struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`

	// Type of keystore.
	Type *string `json:"type,omitempty"`

	GoogleKeyProtectionLevel *string `json:"google_key_protection_level,omitempty"`

	GoogleKeyPurpose *string `json:"google_key_purpose,omitempty"`

	GoogleKmsAlgorithm *string `json:"google_kms_algorithm,omitempty"`
}

// Constants associated with the KeystoresPropertiesCreateGoogleKms.Type property.
// Type of keystore.
const (
	KeystoresPropertiesCreateGoogleKms_Type_AwsKms = "aws_kms"
	KeystoresPropertiesCreateGoogleKms_Type_AzureKeyVault = "azure_key_vault"
	KeystoresPropertiesCreateGoogleKms_Type_GoogleKms = "google_kms"
	KeystoresPropertiesCreateGoogleKms_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoresPropertiesCreateGoogleKms.GoogleKeyProtectionLevel property.
const (
	KeystoresPropertiesCreateGoogleKms_GoogleKeyProtectionLevel_Hsm = "hsm"
	KeystoresPropertiesCreateGoogleKms_GoogleKeyProtectionLevel_Software = "software"
)

// Constants associated with the KeystoresPropertiesCreateGoogleKms.GoogleKeyPurpose property.
const (
	KeystoresPropertiesCreateGoogleKms_GoogleKeyPurpose_AsymmetricDecrypt = "asymmetric_decrypt"
	KeystoresPropertiesCreateGoogleKms_GoogleKeyPurpose_AsymmetricSign = "asymmetric_sign"
	KeystoresPropertiesCreateGoogleKms_GoogleKeyPurpose_EncryptDecrypt = "encrypt_decrypt"
	KeystoresPropertiesCreateGoogleKms_GoogleKeyPurpose_Mac = "mac"
)

// Constants associated with the KeystoresPropertiesCreateGoogleKms.GoogleKmsAlgorithm property.
const (
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_EcSignP256Sha256 = "ec_sign_p256_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_EcSignP384Sha384 = "ec_sign_p384_sha384"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_EcSignSecp256k1Sha256 = "ec_sign_secp256k1_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_GoogleSymmetricEncryption = "google_symmetric_encryption"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_HmacSha256 = "hmac_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha1 = "rsa_decrypt_oaep_2048_sha1"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha256 = "rsa_decrypt_oaep_2048_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha1 = "rsa_decrypt_oaep_3072_sha1"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha256 = "rsa_decrypt_oaep_3072_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha1 = "rsa_decrypt_oaep_4096_sha1"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha256 = "rsa_decrypt_oaep_4096_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha512 = "rsa_decrypt_oaep_4096_sha512"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs12048Sha256 = "rsa_sign_pkcs1_2048_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs13072Sha256 = "rsa_sign_pkcs1_3072_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs14096Sha256 = "rsa_sign_pkcs1_4096_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs14096Sha512 = "rsa_sign_pkcs1_4096_sha512"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignPss2048Sha256 = "rsa_sign_pss_2048_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignPss3072Sha256 = "rsa_sign_pss_3072_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignPss4096Sha256 = "rsa_sign_pss_4096_sha256"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignPss4096Sha512 = "rsa_sign_pss_4096_sha512"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs12048 = "rsa_sign_raw_pkcs1_2048"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs13072 = "rsa_sign_raw_pkcs1_3072"
	KeystoresPropertiesCreateGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs14096 = "rsa_sign_raw_pkcs1_4096"
)

func (*KeystoresPropertiesCreateGoogleKms) isaKeystoresPropertiesCreate() bool {
	return true
}

// UnmarshalKeystoresPropertiesCreateGoogleKms unmarshals an instance of KeystoresPropertiesCreateGoogleKms from the specified map of raw messages.
func UnmarshalKeystoresPropertiesCreateGoogleKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesCreateGoogleKms)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_protection_level", &obj.GoogleKeyProtectionLevel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_purpose", &obj.GoogleKeyPurpose)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_kms_algorithm", &obj.GoogleKmsAlgorithm)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoresPropertiesCreateIbmCloudKms : KeystoresPropertiesCreateIbmCloudKms struct
// This model "extends" KeystoresPropertiesCreate
type KeystoresPropertiesCreateIbmCloudKms struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`

	// Type of keystore.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the KeystoresPropertiesCreateIbmCloudKms.Type property.
// Type of keystore.
const (
	KeystoresPropertiesCreateIbmCloudKms_Type_AwsKms = "aws_kms"
	KeystoresPropertiesCreateIbmCloudKms_Type_AzureKeyVault = "azure_key_vault"
	KeystoresPropertiesCreateIbmCloudKms_Type_GoogleKms = "google_kms"
	KeystoresPropertiesCreateIbmCloudKms_Type_IbmCloudKms = "ibm_cloud_kms"
)

func (*KeystoresPropertiesCreateIbmCloudKms) isaKeystoresPropertiesCreate() bool {
	return true
}

// UnmarshalKeystoresPropertiesCreateIbmCloudKms unmarshals an instance of KeystoresPropertiesCreateIbmCloudKms from the specified map of raw messages.
func UnmarshalKeystoresPropertiesCreateIbmCloudKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesCreateIbmCloudKms)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
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

// KeystoresPropertiesUpdateAwsKms : KeystoresPropertiesUpdateAwsKms struct
// This model "extends" KeystoresPropertiesUpdate
type KeystoresPropertiesUpdateAwsKms struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`
}

func (*KeystoresPropertiesUpdateAwsKms) isaKeystoresPropertiesUpdate() bool {
	return true
}

// UnmarshalKeystoresPropertiesUpdateAwsKms unmarshals an instance of KeystoresPropertiesUpdateAwsKms from the specified map of raw messages.
func UnmarshalKeystoresPropertiesUpdateAwsKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesUpdateAwsKms)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoresPropertiesUpdateAzure : KeystoresPropertiesUpdateAzure struct
// This model "extends" KeystoresPropertiesUpdate
type KeystoresPropertiesUpdateAzure struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`
}

func (*KeystoresPropertiesUpdateAzure) isaKeystoresPropertiesUpdate() bool {
	return true
}

// UnmarshalKeystoresPropertiesUpdateAzure unmarshals an instance of KeystoresPropertiesUpdateAzure from the specified map of raw messages.
func UnmarshalKeystoresPropertiesUpdateAzure(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesUpdateAzure)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoresPropertiesUpdateGoogleKms : KeystoresPropertiesUpdateGoogleKms struct
// This model "extends" KeystoresPropertiesUpdate
type KeystoresPropertiesUpdateGoogleKms struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`

	GoogleKeyProtectionLevel *string `json:"google_key_protection_level,omitempty"`

	GoogleKeyPurpose *string `json:"google_key_purpose,omitempty"`

	GoogleKmsAlgorithm *string `json:"google_kms_algorithm,omitempty"`
}

// Constants associated with the KeystoresPropertiesUpdateGoogleKms.GoogleKeyProtectionLevel property.
const (
	KeystoresPropertiesUpdateGoogleKms_GoogleKeyProtectionLevel_Hsm = "hsm"
	KeystoresPropertiesUpdateGoogleKms_GoogleKeyProtectionLevel_Software = "software"
)

// Constants associated with the KeystoresPropertiesUpdateGoogleKms.GoogleKeyPurpose property.
const (
	KeystoresPropertiesUpdateGoogleKms_GoogleKeyPurpose_AsymmetricDecrypt = "asymmetric_decrypt"
	KeystoresPropertiesUpdateGoogleKms_GoogleKeyPurpose_AsymmetricSign = "asymmetric_sign"
	KeystoresPropertiesUpdateGoogleKms_GoogleKeyPurpose_EncryptDecrypt = "encrypt_decrypt"
	KeystoresPropertiesUpdateGoogleKms_GoogleKeyPurpose_Mac = "mac"
)

// Constants associated with the KeystoresPropertiesUpdateGoogleKms.GoogleKmsAlgorithm property.
const (
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_EcSignP256Sha256 = "ec_sign_p256_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_EcSignP384Sha384 = "ec_sign_p384_sha384"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_EcSignSecp256k1Sha256 = "ec_sign_secp256k1_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_GoogleSymmetricEncryption = "google_symmetric_encryption"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_HmacSha256 = "hmac_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha1 = "rsa_decrypt_oaep_2048_sha1"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep2048Sha256 = "rsa_decrypt_oaep_2048_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha1 = "rsa_decrypt_oaep_3072_sha1"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep3072Sha256 = "rsa_decrypt_oaep_3072_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha1 = "rsa_decrypt_oaep_4096_sha1"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha256 = "rsa_decrypt_oaep_4096_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaDecryptOaep4096Sha512 = "rsa_decrypt_oaep_4096_sha512"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs12048Sha256 = "rsa_sign_pkcs1_2048_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs13072Sha256 = "rsa_sign_pkcs1_3072_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs14096Sha256 = "rsa_sign_pkcs1_4096_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignPkcs14096Sha512 = "rsa_sign_pkcs1_4096_sha512"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignPss2048Sha256 = "rsa_sign_pss_2048_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignPss3072Sha256 = "rsa_sign_pss_3072_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignPss4096Sha256 = "rsa_sign_pss_4096_sha256"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignPss4096Sha512 = "rsa_sign_pss_4096_sha512"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs12048 = "rsa_sign_raw_pkcs1_2048"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs13072 = "rsa_sign_raw_pkcs1_3072"
	KeystoresPropertiesUpdateGoogleKms_GoogleKmsAlgorithm_RsaSignRawPkcs14096 = "rsa_sign_raw_pkcs1_4096"
)

func (*KeystoresPropertiesUpdateGoogleKms) isaKeystoresPropertiesUpdate() bool {
	return true
}

// UnmarshalKeystoresPropertiesUpdateGoogleKms unmarshals an instance of KeystoresPropertiesUpdateGoogleKms from the specified map of raw messages.
func UnmarshalKeystoresPropertiesUpdateGoogleKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesUpdateGoogleKms)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_protection_level", &obj.GoogleKeyProtectionLevel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_key_purpose", &obj.GoogleKeyPurpose)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "google_kms_algorithm", &obj.GoogleKmsAlgorithm)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoresPropertiesUpdateIbmCloudKms : KeystoresPropertiesUpdateIbmCloudKms struct
// This model "extends" KeystoresPropertiesUpdate
type KeystoresPropertiesUpdateIbmCloudKms struct {
	// Which keystore group to distribute the key to.
	Group *string `json:"group,omitempty"`
}

func (*KeystoresPropertiesUpdateIbmCloudKms) isaKeystoresPropertiesUpdate() bool {
	return true
}

// UnmarshalKeystoresPropertiesUpdateIbmCloudKms unmarshals an instance of KeystoresPropertiesUpdateIbmCloudKms from the specified map of raw messages.
func UnmarshalKeystoresPropertiesUpdateIbmCloudKms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoresPropertiesUpdateIbmCloudKms)
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate : You can connect your Hyper Protect Crypto Services instance to the keystores of another Hyper Protect Crypto Services
// or Key Protect instance, and manage its KMS keys using the current service instance.
// This model "extends" KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate
type KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name" validate:"required"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`

	// API endpoint of the IBM Cloud keystore.
	IbmApiEndpoint *string `json:"ibm_api_endpoint" validate:"required"`

	// Endpoint of the IAM service for this IBM Cloud keystore.
	IbmIamEndpoint *string `json:"ibm_iam_endpoint" validate:"required"`

	// The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.
	IbmApiKey *string `json:"ibm_api_key" validate:"required"`

	// The instance ID of the IBM Cloud keystore.
	IbmInstanceID *string `json:"ibm_instance_id" validate:"required"`

	// Possible IBM Cloud KMS variants.
	IbmVariant *string `json:"ibm_variant" validate:"required"`

	// The key ring of an IBM Cloud KMS Keystore.
	IbmKeyRing *string `json:"ibm_key_ring,omitempty"`
}

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate.Type property.
// Type of keystore.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate_Type_AwsKms = "aws_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate_Type_GoogleKms = "google_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate.IbmVariant property.
// Possible IBM Cloud KMS variants.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate_IbmVariant_Hpcs = "hpcs"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate_IbmVariant_Internal = "internal"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate_IbmVariant_KeyProtect = "key_protect"
)

// NewKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate : Instantiate KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate (Generic Model Constructor)
func (*UkoV4) NewKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate(typeVar string, vault *VaultReferenceInCreationRequest, name string, ibmApiEndpoint string, ibmIamEndpoint string, ibmApiKey string, ibmInstanceID string, ibmVariant string) (_model *KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate, err error) {
	_model = &KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate{
		Type: core.StringPtr(typeVar),
		Vault: vault,
		Name: core.StringPtr(name),
		IbmApiEndpoint: core.StringPtr(ibmApiEndpoint),
		IbmIamEndpoint: core.StringPtr(ibmIamEndpoint),
		IbmApiKey: core.StringPtr(ibmApiKey),
		IbmInstanceID: core.StringPtr(ibmInstanceID),
		IbmVariant: core.StringPtr(ibmVariant),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate() bool {
	return true
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate) isaKeystoreCreationRequest() bool {
	return true
}

// UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate unmarshals an instance of KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate from the specified map of raw messages.
func UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReferenceInCreationRequest)
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
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_endpoint", &obj.IbmApiEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_iam_endpoint", &obj.IbmIamEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_api_key", &obj.IbmApiKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_instance_id", &obj.IbmInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_variant", &obj.IbmVariant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_key_ring", &obj.IbmKeyRing)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate : An internal keystore generation request describing all information necessary to generate an internal keystore. It
// only require name.
// Models which "extend" this model:
// - KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate
// This model "extends" KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate
type KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Possible IBM Cloud KMS variants.
	IbmVariant *string `json:"ibm_variant" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name" validate:"required"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`
}

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate.Type property.
// Type of keystore.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate_Type_AwsKms = "aws_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate_Type_GoogleKms = "google_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate.IbmVariant property.
// Possible IBM Cloud KMS variants.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate_IbmVariant_Hpcs = "hpcs"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate_IbmVariant_Internal = "internal"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate_IbmVariant_KeyProtect = "key_protect"
)
func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate() bool {
	return true
}

type KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateIntf interface {
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateIntf
	isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate() bool
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate() bool {
	return true
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate) isaKeystoreCreationRequest() bool {
	return true
}

// UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate unmarshals an instance of KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate from the specified map of raw messages.
func UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReferenceInCreationRequest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_variant", &obj.IbmVariant)
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
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate : Base of a keystore update.
// This model "extends" KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate
type KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate struct {
	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`
}

func (*KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate) isaKeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdate() bool {
	return true
}

func (*KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate) isaKeystoreUpdateRequest() bool {
	return true
}

// UnmarshalKeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate unmarshals an instance of KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate from the specified map of raw messages.
func UnmarshalKeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreUpdateRequestKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate : KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate struct
// Models which "extend" this model:
// - KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate
// This model "extends" KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate
type KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Possible IBM Cloud KMS variants.
	IbmVariant *string `json:"ibm_variant" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`
}

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate.Type property.
// Type of keystore.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate_Type_AwsKms = "aws_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate_Type_GoogleKms = "google_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate.IbmVariant property.
// Possible IBM Cloud KMS variants.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate_IbmVariant_Hpcs = "hpcs"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate_IbmVariant_Internal = "internal"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate_IbmVariant_KeyProtect = "key_protect"
)
func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate() bool {
	return true
}

type KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateIntf interface {
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateIntf
	isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate() bool
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate() bool {
	return true
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate() bool {
	return true
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate) isaKeystoreCreationRequest() bool {
	return true
}

// UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate unmarshals an instance of KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate from the specified map of raw messages.
func UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReferenceInCreationRequest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_variant", &obj.IbmVariant)
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
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate : Base of a keystore update.
// This model "extends" KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate
type KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate struct {
	// Type of keystore.
	Type *string `json:"type" validate:"required"`

	Vault *VaultReferenceInCreationRequest `json:"vault" validate:"required"`

	// Possible IBM Cloud KMS variants.
	IbmVariant *string `json:"ibm_variant" validate:"required"`

	// Name of a target keystore.
	Name *string `json:"name,omitempty"`

	// Description of the keystore.
	Description *string `json:"description,omitempty"`

	// A list of groups that this keystore belongs to.
	Groups []string `json:"groups,omitempty"`
}

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate.Type property.
// Type of keystore.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate_Type_AwsKms = "aws_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate_Type_AzureKeyVault = "azure_key_vault"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate_Type_GoogleKms = "google_kms"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate_Type_IbmCloudKms = "ibm_cloud_kms"
)

// Constants associated with the KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate.IbmVariant property.
// Possible IBM Cloud KMS variants.
const (
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate_IbmVariant_Hpcs = "hpcs"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate_IbmVariant_Internal = "internal"
	KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate_IbmVariant_KeyProtect = "key_protect"
)

// NewKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate : Instantiate KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate (Generic Model Constructor)
func (*UkoV4) NewKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate(typeVar string, vault *VaultReferenceInCreationRequest, ibmVariant string) (_model *KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate, err error) {
	_model = &KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate{
		Type: core.StringPtr(typeVar),
		Vault: vault,
		IbmVariant: core.StringPtr(ibmVariant),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdate() bool {
	return true
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate() bool {
	return true
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate) isaKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate() bool {
	return true
}

func (*KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate) isaKeystoreCreationRequest() bool {
	return true
}

// UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate unmarshals an instance of KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate from the specified map of raw messages.
func UnmarshalKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vault", &obj.Vault, UnmarshalVaultReferenceInCreationRequest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_variant", &obj.IbmVariant)
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
	err = core.UnmarshalPrimitive(m, "groups", &obj.Groups)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

//
// ManagedKeysPager can be used to simplify the use of the "ListManagedKeys" method.
//
type ManagedKeysPager struct {
	hasNext bool
	options *ListManagedKeysOptions
	client  *UkoV4
	pageContext struct {
		next *int64
	}
}

// NewManagedKeysPager returns a new ManagedKeysPager instance.
func (uko *UkoV4) NewManagedKeysPager(options *ListManagedKeysOptions) (pager *ManagedKeysPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListManagedKeysOptions = *options
	pager = &ManagedKeysPager{
		hasNext: true,
		options: &optionsCopy,
		client:  uko,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ManagedKeysPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ManagedKeysPager) GetNextWithContext(ctx context.Context) (page []ManagedKey, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListManagedKeysWithContext(ctx, pager.options)
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
	page = result.ManagedKeys

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ManagedKeysPager) GetAllWithContext(ctx context.Context) (allItems []ManagedKey, err error) {
	for pager.HasNext() {
		var nextPage []ManagedKey
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ManagedKeysPager) GetNext() (page []ManagedKey, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ManagedKeysPager) GetAll() (allItems []ManagedKey, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// AssociatedResourcesForManagedKeyPager can be used to simplify the use of the "ListAssociatedResourcesForManagedKey" method.
//
type AssociatedResourcesForManagedKeyPager struct {
	hasNext bool
	options *ListAssociatedResourcesForManagedKeyOptions
	client  *UkoV4
	pageContext struct {
		next *int64
	}
}

// NewAssociatedResourcesForManagedKeyPager returns a new AssociatedResourcesForManagedKeyPager instance.
func (uko *UkoV4) NewAssociatedResourcesForManagedKeyPager(options *ListAssociatedResourcesForManagedKeyOptions) (pager *AssociatedResourcesForManagedKeyPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListAssociatedResourcesForManagedKeyOptions = *options
	pager = &AssociatedResourcesForManagedKeyPager{
		hasNext: true,
		options: &optionsCopy,
		client:  uko,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AssociatedResourcesForManagedKeyPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AssociatedResourcesForManagedKeyPager) GetNextWithContext(ctx context.Context) (page []AssociatedResource, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListAssociatedResourcesForManagedKeyWithContext(ctx, pager.options)
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
	page = result.AssociatedResources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AssociatedResourcesForManagedKeyPager) GetAllWithContext(ctx context.Context) (allItems []AssociatedResource, err error) {
	for pager.HasNext() {
		var nextPage []AssociatedResource
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *AssociatedResourcesForManagedKeyPager) GetNext() (page []AssociatedResource, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AssociatedResourcesForManagedKeyPager) GetAll() (allItems []AssociatedResource, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// KeyTemplatesPager can be used to simplify the use of the "ListKeyTemplates" method.
//
type KeyTemplatesPager struct {
	hasNext bool
	options *ListKeyTemplatesOptions
	client  *UkoV4
	pageContext struct {
		next *int64
	}
}

// NewKeyTemplatesPager returns a new KeyTemplatesPager instance.
func (uko *UkoV4) NewKeyTemplatesPager(options *ListKeyTemplatesOptions) (pager *KeyTemplatesPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListKeyTemplatesOptions = *options
	pager = &KeyTemplatesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  uko,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *KeyTemplatesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *KeyTemplatesPager) GetNextWithContext(ctx context.Context) (page []Template, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListKeyTemplatesWithContext(ctx, pager.options)
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
	page = result.Templates

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *KeyTemplatesPager) GetAllWithContext(ctx context.Context) (allItems []Template, err error) {
	for pager.HasNext() {
		var nextPage []Template
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *KeyTemplatesPager) GetNext() (page []Template, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *KeyTemplatesPager) GetAll() (allItems []Template, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// KeystoresPager can be used to simplify the use of the "ListKeystores" method.
//
type KeystoresPager struct {
	hasNext bool
	options *ListKeystoresOptions
	client  *UkoV4
	pageContext struct {
		next *int64
	}
}

// NewKeystoresPager returns a new KeystoresPager instance.
func (uko *UkoV4) NewKeystoresPager(options *ListKeystoresOptions) (pager *KeystoresPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListKeystoresOptions = *options
	pager = &KeystoresPager{
		hasNext: true,
		options: &optionsCopy,
		client:  uko,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *KeystoresPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *KeystoresPager) GetNextWithContext(ctx context.Context) (page []KeystoreIntf, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListKeystoresWithContext(ctx, pager.options)
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
	page = result.Keystores

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *KeystoresPager) GetAllWithContext(ctx context.Context) (allItems []KeystoreIntf, err error) {
	for pager.HasNext() {
		var nextPage []KeystoreIntf
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *KeystoresPager) GetNext() (page []KeystoreIntf, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *KeystoresPager) GetAll() (allItems []KeystoreIntf, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// AssociatedResourcesForTargetKeystorePager can be used to simplify the use of the "ListAssociatedResourcesForTargetKeystore" method.
//
type AssociatedResourcesForTargetKeystorePager struct {
	hasNext bool
	options *ListAssociatedResourcesForTargetKeystoreOptions
	client  *UkoV4
	pageContext struct {
		next *int64
	}
}

// NewAssociatedResourcesForTargetKeystorePager returns a new AssociatedResourcesForTargetKeystorePager instance.
func (uko *UkoV4) NewAssociatedResourcesForTargetKeystorePager(options *ListAssociatedResourcesForTargetKeystoreOptions) (pager *AssociatedResourcesForTargetKeystorePager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListAssociatedResourcesForTargetKeystoreOptions = *options
	pager = &AssociatedResourcesForTargetKeystorePager{
		hasNext: true,
		options: &optionsCopy,
		client:  uko,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AssociatedResourcesForTargetKeystorePager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AssociatedResourcesForTargetKeystorePager) GetNextWithContext(ctx context.Context) (page []AssociatedResource, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListAssociatedResourcesForTargetKeystoreWithContext(ctx, pager.options)
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
	page = result.AssociatedResources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AssociatedResourcesForTargetKeystorePager) GetAllWithContext(ctx context.Context) (allItems []AssociatedResource, err error) {
	for pager.HasNext() {
		var nextPage []AssociatedResource
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *AssociatedResourcesForTargetKeystorePager) GetNext() (page []AssociatedResource, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AssociatedResourcesForTargetKeystorePager) GetAll() (allItems []AssociatedResource, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// ManagedKeysFromKeystorePager can be used to simplify the use of the "ListManagedKeysFromKeystore" method.
//
type ManagedKeysFromKeystorePager struct {
	hasNext bool
	options *ListManagedKeysFromKeystoreOptions
	client  *UkoV4
	pageContext struct {
		next *int64
	}
}

// NewManagedKeysFromKeystorePager returns a new ManagedKeysFromKeystorePager instance.
func (uko *UkoV4) NewManagedKeysFromKeystorePager(options *ListManagedKeysFromKeystoreOptions) (pager *ManagedKeysFromKeystorePager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListManagedKeysFromKeystoreOptions = *options
	pager = &ManagedKeysFromKeystorePager{
		hasNext: true,
		options: &optionsCopy,
		client:  uko,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ManagedKeysFromKeystorePager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ManagedKeysFromKeystorePager) GetNextWithContext(ctx context.Context) (page []ManagedKey, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListManagedKeysFromKeystoreWithContext(ctx, pager.options)
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
	page = result.ManagedKeys

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ManagedKeysFromKeystorePager) GetAllWithContext(ctx context.Context) (allItems []ManagedKey, err error) {
	for pager.HasNext() {
		var nextPage []ManagedKey
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ManagedKeysFromKeystorePager) GetNext() (page []ManagedKey, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ManagedKeysFromKeystorePager) GetAll() (allItems []ManagedKey, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// VaultsPager can be used to simplify the use of the "ListVaults" method.
//
type VaultsPager struct {
	hasNext bool
	options *ListVaultsOptions
	client  *UkoV4
	pageContext struct {
		next *int64
	}
}

// NewVaultsPager returns a new VaultsPager instance.
func (uko *UkoV4) NewVaultsPager(options *ListVaultsOptions) (pager *VaultsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListVaultsOptions = *options
	pager = &VaultsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  uko,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *VaultsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *VaultsPager) GetNextWithContext(ctx context.Context) (page []Vault, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListVaultsWithContext(ctx, pager.options)
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
	page = result.Vaults

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *VaultsPager) GetAllWithContext(ctx context.Context) (allItems []Vault, err error) {
	for pager.HasNext() {
		var nextPage []Vault
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *VaultsPager) GetNext() (page []Vault, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *VaultsPager) GetAll() (allItems []Vault, err error) {
	return pager.GetAllWithContext(context.Background())
}
