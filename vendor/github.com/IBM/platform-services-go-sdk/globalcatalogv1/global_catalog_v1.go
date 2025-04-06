/**
 * (C) Copyright IBM Corp. 2025.
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
 * IBM OpenAPI SDK Code Generator Version: 3.102.0-615ec964-20250307-203034
 */

// Package globalcatalogv1 : Operations and models for the GlobalCatalogV1 service
package globalcatalogv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// GlobalCatalogV1 : The catalog service manages offerings across geographies as the system of record. The catalog
// supports a RESTful API where users can retrieve information about existing offerings and create, manage, and delete
// their offerings. Start with the base URL and use the endpoints to retrieve metadata about services in the catalog and
// manage service visbility. Depending on the kind of object, the metadata can include information about pricing,
// provisioning, regions, and more. For more information, see the [catalog
// documentation](https://cloud.ibm.com/docs/overview/catalog.html#global-catalog-overview).
//
// API Version: 1.0.3
type GlobalCatalogV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://globalcatalog.cloud.ibm.com/api/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "global_catalog"

// GlobalCatalogV1Options : Service options
type GlobalCatalogV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewGlobalCatalogV1UsingExternalConfig : constructs an instance of GlobalCatalogV1 with passed in options and external configuration.
func NewGlobalCatalogV1UsingExternalConfig(options *GlobalCatalogV1Options) (globalCatalog *GlobalCatalogV1, err error) {
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

	globalCatalog, err = NewGlobalCatalogV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = globalCatalog.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = globalCatalog.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewGlobalCatalogV1 : constructs an instance of GlobalCatalogV1 with passed in options.
func NewGlobalCatalogV1(options *GlobalCatalogV1Options) (service *GlobalCatalogV1, err error) {
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

	service = &GlobalCatalogV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "globalCatalog" suitable for processing requests.
func (globalCatalog *GlobalCatalogV1) Clone() *GlobalCatalogV1 {
	if core.IsNil(globalCatalog) {
		return nil
	}
	clone := *globalCatalog
	clone.Service = globalCatalog.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (globalCatalog *GlobalCatalogV1) SetServiceURL(url string) error {
	err := globalCatalog.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (globalCatalog *GlobalCatalogV1) GetServiceURL() string {
	return globalCatalog.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (globalCatalog *GlobalCatalogV1) SetDefaultHeaders(headers http.Header) {
	globalCatalog.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (globalCatalog *GlobalCatalogV1) SetEnableGzipCompression(enableGzip bool) {
	globalCatalog.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (globalCatalog *GlobalCatalogV1) GetEnableGzipCompression() bool {
	return globalCatalog.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (globalCatalog *GlobalCatalogV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	globalCatalog.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (globalCatalog *GlobalCatalogV1) DisableRetries() {
	globalCatalog.Service.DisableRetries()
}

// ListCatalogEntries : Returns parent catalog entries
// Includes key information, such as ID, name, kind, CRN, tags, and provider. This endpoint is ETag enabled.
func (globalCatalog *GlobalCatalogV1) ListCatalogEntries(listCatalogEntriesOptions *ListCatalogEntriesOptions) (result *EntrySearchResult, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.ListCatalogEntriesWithContext(context.Background(), listCatalogEntriesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListCatalogEntriesWithContext is an alternate form of the ListCatalogEntries method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) ListCatalogEntriesWithContext(ctx context.Context, listCatalogEntriesOptions *ListCatalogEntriesOptions) (result *EntrySearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCatalogEntriesOptions, "listCatalogEntriesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listCatalogEntriesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "ListCatalogEntries")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listCatalogEntriesOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*listCatalogEntriesOptions.Account))
	}
	if listCatalogEntriesOptions.Include != nil {
		builder.AddQuery("include", fmt.Sprint(*listCatalogEntriesOptions.Include))
	}
	if listCatalogEntriesOptions.Q != nil {
		builder.AddQuery("q", fmt.Sprint(*listCatalogEntriesOptions.Q))
	}
	if listCatalogEntriesOptions.SortBy != nil {
		builder.AddQuery("sort-by", fmt.Sprint(*listCatalogEntriesOptions.SortBy))
	}
	if listCatalogEntriesOptions.Descending != nil {
		builder.AddQuery("descending", fmt.Sprint(*listCatalogEntriesOptions.Descending))
	}
	if listCatalogEntriesOptions.Languages != nil {
		builder.AddQuery("languages", fmt.Sprint(*listCatalogEntriesOptions.Languages))
	}
	if listCatalogEntriesOptions.Catalog != nil {
		builder.AddQuery("catalog", fmt.Sprint(*listCatalogEntriesOptions.Catalog))
	}
	if listCatalogEntriesOptions.Complete != nil {
		builder.AddQuery("complete", fmt.Sprint(*listCatalogEntriesOptions.Complete))
	}
	if listCatalogEntriesOptions.Offset != nil {
		builder.AddQuery("_offset", fmt.Sprint(*listCatalogEntriesOptions.Offset))
	}
	if listCatalogEntriesOptions.Limit != nil {
		builder.AddQuery("_limit", fmt.Sprint(*listCatalogEntriesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_catalog_entries", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEntrySearchResult)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateCatalogEntry : Create a catalog entry
// The created catalog entry is restricted by default. You must have an administrator or editor role in the scope of the
// provided token. This API will return an ETag that can be used for standard ETag processing, except when depth query
// is used.
func (globalCatalog *GlobalCatalogV1) CreateCatalogEntry(createCatalogEntryOptions *CreateCatalogEntryOptions) (result *CatalogEntry, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.CreateCatalogEntryWithContext(context.Background(), createCatalogEntryOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateCatalogEntryWithContext is an alternate form of the CreateCatalogEntry method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) CreateCatalogEntryWithContext(ctx context.Context, createCatalogEntryOptions *CreateCatalogEntryOptions) (result *CatalogEntry, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCatalogEntryOptions, "createCatalogEntryOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createCatalogEntryOptions, "createCatalogEntryOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createCatalogEntryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "CreateCatalogEntry")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createCatalogEntryOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*createCatalogEntryOptions.Account))
	}

	body := make(map[string]interface{})
	if createCatalogEntryOptions.Name != nil {
		body["name"] = createCatalogEntryOptions.Name
	}
	if createCatalogEntryOptions.Kind != nil {
		body["kind"] = createCatalogEntryOptions.Kind
	}
	if createCatalogEntryOptions.OverviewUI != nil {
		body["overview_ui"] = createCatalogEntryOptions.OverviewUI
	}
	if createCatalogEntryOptions.Images != nil {
		body["images"] = createCatalogEntryOptions.Images
	}
	if createCatalogEntryOptions.Disabled != nil {
		body["disabled"] = createCatalogEntryOptions.Disabled
	}
	if createCatalogEntryOptions.Tags != nil {
		body["tags"] = createCatalogEntryOptions.Tags
	}
	if createCatalogEntryOptions.Provider != nil {
		body["provider"] = createCatalogEntryOptions.Provider
	}
	if createCatalogEntryOptions.ID != nil {
		body["id"] = createCatalogEntryOptions.ID
	}
	if createCatalogEntryOptions.ParentID != nil {
		body["parent_id"] = createCatalogEntryOptions.ParentID
	}
	if createCatalogEntryOptions.Group != nil {
		body["group"] = createCatalogEntryOptions.Group
	}
	if createCatalogEntryOptions.Active != nil {
		body["active"] = createCatalogEntryOptions.Active
	}
	if createCatalogEntryOptions.URL != nil {
		body["url"] = createCatalogEntryOptions.URL
	}
	if createCatalogEntryOptions.Metadata != nil {
		body["metadata"] = createCatalogEntryOptions.Metadata
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
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_catalog_entry", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalogEntry)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetCatalogEntry : Get a specific catalog object
// This endpoint returns a specific catalog entry using the object's unique identifier, for example
// `/_*service_name*?complete=true`. This endpoint is ETag enabled. This can be used by an unauthenticated user for
// publicly available services.
func (globalCatalog *GlobalCatalogV1) GetCatalogEntry(getCatalogEntryOptions *GetCatalogEntryOptions) (result *CatalogEntry, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.GetCatalogEntryWithContext(context.Background(), getCatalogEntryOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetCatalogEntryWithContext is an alternate form of the GetCatalogEntry method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) GetCatalogEntryWithContext(ctx context.Context, getCatalogEntryOptions *GetCatalogEntryOptions) (result *CatalogEntry, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCatalogEntryOptions, "getCatalogEntryOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getCatalogEntryOptions, "getCatalogEntryOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getCatalogEntryOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getCatalogEntryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "GetCatalogEntry")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getCatalogEntryOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*getCatalogEntryOptions.Account))
	}
	if getCatalogEntryOptions.Include != nil {
		builder.AddQuery("include", fmt.Sprint(*getCatalogEntryOptions.Include))
	}
	if getCatalogEntryOptions.Languages != nil {
		builder.AddQuery("languages", fmt.Sprint(*getCatalogEntryOptions.Languages))
	}
	if getCatalogEntryOptions.Complete != nil {
		builder.AddQuery("complete", fmt.Sprint(*getCatalogEntryOptions.Complete))
	}
	if getCatalogEntryOptions.Depth != nil {
		builder.AddQuery("depth", fmt.Sprint(*getCatalogEntryOptions.Depth))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_catalog_entry", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalogEntry)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateCatalogEntry : Update a catalog entry
// Update a catalog entry. The visibility of the catalog entry cannot be modified with this endpoint. You must be an
// administrator or editor in the scope of the provided token. This endpoint is ETag enabled.
func (globalCatalog *GlobalCatalogV1) UpdateCatalogEntry(updateCatalogEntryOptions *UpdateCatalogEntryOptions) (result *CatalogEntry, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.UpdateCatalogEntryWithContext(context.Background(), updateCatalogEntryOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateCatalogEntryWithContext is an alternate form of the UpdateCatalogEntry method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) UpdateCatalogEntryWithContext(ctx context.Context, updateCatalogEntryOptions *UpdateCatalogEntryOptions) (result *CatalogEntry, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCatalogEntryOptions, "updateCatalogEntryOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateCatalogEntryOptions, "updateCatalogEntryOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateCatalogEntryOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateCatalogEntryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "UpdateCatalogEntry")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if updateCatalogEntryOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*updateCatalogEntryOptions.Account))
	}
	if updateCatalogEntryOptions.Move != nil {
		builder.AddQuery("move", fmt.Sprint(*updateCatalogEntryOptions.Move))
	}

	body := make(map[string]interface{})
	if updateCatalogEntryOptions.Name != nil {
		body["name"] = updateCatalogEntryOptions.Name
	}
	if updateCatalogEntryOptions.Kind != nil {
		body["kind"] = updateCatalogEntryOptions.Kind
	}
	if updateCatalogEntryOptions.OverviewUI != nil {
		body["overview_ui"] = updateCatalogEntryOptions.OverviewUI
	}
	if updateCatalogEntryOptions.Images != nil {
		body["images"] = updateCatalogEntryOptions.Images
	}
	if updateCatalogEntryOptions.Disabled != nil {
		body["disabled"] = updateCatalogEntryOptions.Disabled
	}
	if updateCatalogEntryOptions.Tags != nil {
		body["tags"] = updateCatalogEntryOptions.Tags
	}
	if updateCatalogEntryOptions.Provider != nil {
		body["provider"] = updateCatalogEntryOptions.Provider
	}
	if updateCatalogEntryOptions.ParentID != nil {
		body["parent_id"] = updateCatalogEntryOptions.ParentID
	}
	if updateCatalogEntryOptions.Group != nil {
		body["group"] = updateCatalogEntryOptions.Group
	}
	if updateCatalogEntryOptions.Active != nil {
		body["active"] = updateCatalogEntryOptions.Active
	}
	if updateCatalogEntryOptions.URL != nil {
		body["url"] = updateCatalogEntryOptions.URL
	}
	if updateCatalogEntryOptions.Metadata != nil {
		body["metadata"] = updateCatalogEntryOptions.Metadata
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
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_catalog_entry", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalogEntry)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteCatalogEntry : Delete a catalog entry
// Delete a catalog entry. This will archive the catalog entry for a minimum of two weeks. While archived, it can be
// restored using the PUT restore API. After two weeks, it will be deleted and cannot be restored. You must have
// administrator role in the scope of the provided token to modify it. This endpoint is ETag enabled.
func (globalCatalog *GlobalCatalogV1) DeleteCatalogEntry(deleteCatalogEntryOptions *DeleteCatalogEntryOptions) (response *core.DetailedResponse, err error) {
	response, err = globalCatalog.DeleteCatalogEntryWithContext(context.Background(), deleteCatalogEntryOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteCatalogEntryWithContext is an alternate form of the DeleteCatalogEntry method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) DeleteCatalogEntryWithContext(ctx context.Context, deleteCatalogEntryOptions *DeleteCatalogEntryOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCatalogEntryOptions, "deleteCatalogEntryOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteCatalogEntryOptions, "deleteCatalogEntryOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteCatalogEntryOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteCatalogEntryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "DeleteCatalogEntry")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteCatalogEntryOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*deleteCatalogEntryOptions.Account))
	}
	if deleteCatalogEntryOptions.Force != nil {
		builder.AddQuery("force", fmt.Sprint(*deleteCatalogEntryOptions.Force))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = globalCatalog.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_catalog_entry", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetChildObjects : Get child catalog entries of a specific kind
// Fetch child catalog entries for a catalog entry with a specific id. This endpoint is ETag enabled. This can be used
// by an unauthenticated user for publicly available services.
func (globalCatalog *GlobalCatalogV1) GetChildObjects(getChildObjectsOptions *GetChildObjectsOptions) (result *EntrySearchResult, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.GetChildObjectsWithContext(context.Background(), getChildObjectsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetChildObjectsWithContext is an alternate form of the GetChildObjects method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) GetChildObjectsWithContext(ctx context.Context, getChildObjectsOptions *GetChildObjectsOptions) (result *EntrySearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getChildObjectsOptions, "getChildObjectsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getChildObjectsOptions, "getChildObjectsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id":   *getChildObjectsOptions.ID,
		"kind": *getChildObjectsOptions.Kind,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}/{kind}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getChildObjectsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "GetChildObjects")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getChildObjectsOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*getChildObjectsOptions.Account))
	}
	if getChildObjectsOptions.Include != nil {
		builder.AddQuery("include", fmt.Sprint(*getChildObjectsOptions.Include))
	}
	if getChildObjectsOptions.Q != nil {
		builder.AddQuery("q", fmt.Sprint(*getChildObjectsOptions.Q))
	}
	if getChildObjectsOptions.SortBy != nil {
		builder.AddQuery("sort-by", fmt.Sprint(*getChildObjectsOptions.SortBy))
	}
	if getChildObjectsOptions.Descending != nil {
		builder.AddQuery("descending", fmt.Sprint(*getChildObjectsOptions.Descending))
	}
	if getChildObjectsOptions.Languages != nil {
		builder.AddQuery("languages", fmt.Sprint(*getChildObjectsOptions.Languages))
	}
	if getChildObjectsOptions.Complete != nil {
		builder.AddQuery("complete", fmt.Sprint(*getChildObjectsOptions.Complete))
	}
	if getChildObjectsOptions.Offset != nil {
		builder.AddQuery("_offset", fmt.Sprint(*getChildObjectsOptions.Offset))
	}
	if getChildObjectsOptions.Limit != nil {
		builder.AddQuery("_limit", fmt.Sprint(*getChildObjectsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_child_objects", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEntrySearchResult)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// RestoreCatalogEntry : Restore archived catalog entry
// Restore an archived catalog entry. You must have an administrator role in the scope of the provided token.
func (globalCatalog *GlobalCatalogV1) RestoreCatalogEntry(restoreCatalogEntryOptions *RestoreCatalogEntryOptions) (response *core.DetailedResponse, err error) {
	response, err = globalCatalog.RestoreCatalogEntryWithContext(context.Background(), restoreCatalogEntryOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// RestoreCatalogEntryWithContext is an alternate form of the RestoreCatalogEntry method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) RestoreCatalogEntryWithContext(ctx context.Context, restoreCatalogEntryOptions *RestoreCatalogEntryOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(restoreCatalogEntryOptions, "restoreCatalogEntryOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(restoreCatalogEntryOptions, "restoreCatalogEntryOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *restoreCatalogEntryOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}/restore`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range restoreCatalogEntryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "RestoreCatalogEntry")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if restoreCatalogEntryOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*restoreCatalogEntryOptions.Account))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = globalCatalog.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "restore_catalog_entry", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetVisibility : Get the visibility constraints for an object
// This endpoint returns the visibility rules for this object. Overall visibility is determined by the parent objects
// and any further restrictions on this object. You must have an administrator role in the scope of the provided token.
// This endpoint is ETag enabled.
func (globalCatalog *GlobalCatalogV1) GetVisibility(getVisibilityOptions *GetVisibilityOptions) (result *Visibility, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.GetVisibilityWithContext(context.Background(), getVisibilityOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetVisibilityWithContext is an alternate form of the GetVisibility method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) GetVisibilityWithContext(ctx context.Context, getVisibilityOptions *GetVisibilityOptions) (result *Visibility, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getVisibilityOptions, "getVisibilityOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getVisibilityOptions, "getVisibilityOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getVisibilityOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}/visibility`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getVisibilityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "GetVisibility")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getVisibilityOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*getVisibilityOptions.Account))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_visibility", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVisibility)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateVisibility : Update visibility
// Update an Object's Visibility. You must have an administrator role in the scope of the provided token. This endpoint
// is ETag enabled.
func (globalCatalog *GlobalCatalogV1) UpdateVisibility(updateVisibilityOptions *UpdateVisibilityOptions) (response *core.DetailedResponse, err error) {
	response, err = globalCatalog.UpdateVisibilityWithContext(context.Background(), updateVisibilityOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateVisibilityWithContext is an alternate form of the UpdateVisibility method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) UpdateVisibilityWithContext(ctx context.Context, updateVisibilityOptions *UpdateVisibilityOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateVisibilityOptions, "updateVisibilityOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateVisibilityOptions, "updateVisibilityOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateVisibilityOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}/visibility`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateVisibilityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "UpdateVisibility")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	if updateVisibilityOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*updateVisibilityOptions.Account))
	}

	body := make(map[string]interface{})
	if updateVisibilityOptions.Extendable != nil {
		body["extendable"] = updateVisibilityOptions.Extendable
	}
	if updateVisibilityOptions.Include != nil {
		body["include"] = updateVisibilityOptions.Include
	}
	if updateVisibilityOptions.Exclude != nil {
		body["exclude"] = updateVisibilityOptions.Exclude
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

	response, err = globalCatalog.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_visibility", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetPricing : Get the pricing for an object
// This endpoint returns the pricing for an object. Static pricing is defined in the catalog. Dynamic pricing is stored
// in IBM Cloud Pricing Catalog. This can be used by an unauthenticated user for publicly available services.
func (globalCatalog *GlobalCatalogV1) GetPricing(getPricingOptions *GetPricingOptions) (result *PricingGet, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.GetPricingWithContext(context.Background(), getPricingOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPricingWithContext is an alternate form of the GetPricing method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) GetPricingWithContext(ctx context.Context, getPricingOptions *GetPricingOptions) (result *PricingGet, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPricingOptions, "getPricingOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPricingOptions, "getPricingOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getPricingOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}/pricing`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getPricingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "GetPricing")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getPricingOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*getPricingOptions.Account))
	}
	if getPricingOptions.DeploymentRegion != nil {
		builder.AddQuery("deployment_region", fmt.Sprint(*getPricingOptions.DeploymentRegion))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_pricing", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPricingGet)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetPricingDeployments : Get the pricing deployments for a plan
// This endpoint returns the deployment pricing for a plan. For a plan it returns a pricing for each visible child
// deployment object. Static pricing is defined in the catalog. Dynamic pricing is stored in IBM Cloud Pricing Catalog.
// This can be used by an unauthenticated user for publicly available services.
func (globalCatalog *GlobalCatalogV1) GetPricingDeployments(getPricingDeploymentsOptions *GetPricingDeploymentsOptions) (result *PricingSearchResult, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.GetPricingDeploymentsWithContext(context.Background(), getPricingDeploymentsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPricingDeploymentsWithContext is an alternate form of the GetPricingDeployments method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) GetPricingDeploymentsWithContext(ctx context.Context, getPricingDeploymentsOptions *GetPricingDeploymentsOptions) (result *PricingSearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPricingDeploymentsOptions, "getPricingDeploymentsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPricingDeploymentsOptions, "getPricingDeploymentsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getPricingDeploymentsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}/pricing/deployment`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getPricingDeploymentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "GetPricingDeployments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getPricingDeploymentsOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*getPricingDeploymentsOptions.Account))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_pricing_deployments", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPricingSearchResult)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAuditLogs : Get the audit logs for an object
// This endpoint returns the audit logs for an object. Only administrators and editors can get logs.
func (globalCatalog *GlobalCatalogV1) GetAuditLogs(getAuditLogsOptions *GetAuditLogsOptions) (result *AuditSearchResult, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.GetAuditLogsWithContext(context.Background(), getAuditLogsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAuditLogsWithContext is an alternate form of the GetAuditLogs method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) GetAuditLogsWithContext(ctx context.Context, getAuditLogsOptions *GetAuditLogsOptions) (result *AuditSearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAuditLogsOptions, "getAuditLogsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAuditLogsOptions, "getAuditLogsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getAuditLogsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{id}/logs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAuditLogsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "GetAuditLogs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getAuditLogsOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*getAuditLogsOptions.Account))
	}
	if getAuditLogsOptions.Ascending != nil {
		builder.AddQuery("ascending", fmt.Sprint(*getAuditLogsOptions.Ascending))
	}
	if getAuditLogsOptions.Startat != nil {
		builder.AddQuery("startat", fmt.Sprint(*getAuditLogsOptions.Startat))
	}
	if getAuditLogsOptions.Offset != nil {
		builder.AddQuery("_offset", fmt.Sprint(*getAuditLogsOptions.Offset))
	}
	if getAuditLogsOptions.Limit != nil {
		builder.AddQuery("_limit", fmt.Sprint(*getAuditLogsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_audit_logs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditSearchResult)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListArtifacts : Get artifacts
// This endpoint returns a list of artifacts for an object.
func (globalCatalog *GlobalCatalogV1) ListArtifacts(listArtifactsOptions *ListArtifactsOptions) (result *Artifacts, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.ListArtifactsWithContext(context.Background(), listArtifactsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListArtifactsWithContext is an alternate form of the ListArtifacts method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) ListArtifactsWithContext(ctx context.Context, listArtifactsOptions *ListArtifactsOptions) (result *Artifacts, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listArtifactsOptions, "listArtifactsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listArtifactsOptions, "listArtifactsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"object_id": *listArtifactsOptions.ObjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{object_id}/artifacts`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listArtifactsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "ListArtifacts")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listArtifactsOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*listArtifactsOptions.Account))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalCatalog.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_artifacts", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalArtifacts)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetArtifact : Get artifact
// This endpoint returns the binary of an artifact.
func (globalCatalog *GlobalCatalogV1) GetArtifact(getArtifactOptions *GetArtifactOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	result, response, err = globalCatalog.GetArtifactWithContext(context.Background(), getArtifactOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetArtifactWithContext is an alternate form of the GetArtifact method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) GetArtifactWithContext(ctx context.Context, getArtifactOptions *GetArtifactOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getArtifactOptions, "getArtifactOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getArtifactOptions, "getArtifactOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"object_id":   *getArtifactOptions.ObjectID,
		"artifact_id": *getArtifactOptions.ArtifactID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{object_id}/artifacts/{artifact_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getArtifactOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "GetArtifact")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "*/*")
	if getArtifactOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*getArtifactOptions.Accept))
	}

	if getArtifactOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*getArtifactOptions.Account))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = globalCatalog.Service.Request(request, &result)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_artifact", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UploadArtifact : Upload artifact
// This endpoint uploads the binary for an artifact. Only administrators and editors can upload artifacts.
func (globalCatalog *GlobalCatalogV1) UploadArtifact(uploadArtifactOptions *UploadArtifactOptions) (response *core.DetailedResponse, err error) {
	response, err = globalCatalog.UploadArtifactWithContext(context.Background(), uploadArtifactOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UploadArtifactWithContext is an alternate form of the UploadArtifact method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) UploadArtifactWithContext(ctx context.Context, uploadArtifactOptions *UploadArtifactOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(uploadArtifactOptions, "uploadArtifactOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(uploadArtifactOptions, "uploadArtifactOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"object_id":   *uploadArtifactOptions.ObjectID,
		"artifact_id": *uploadArtifactOptions.ArtifactID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{object_id}/artifacts/{artifact_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range uploadArtifactOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "UploadArtifact")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if uploadArtifactOptions.ContentType != nil {
		builder.AddHeader("Content-Type", fmt.Sprint(*uploadArtifactOptions.ContentType))
	}

	if uploadArtifactOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*uploadArtifactOptions.Account))
	}

	_, err = builder.SetBodyContent(core.StringNilMapper(uploadArtifactOptions.ContentType), nil, nil, uploadArtifactOptions.Artifact)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-body-content-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = globalCatalog.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "upload_artifact", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// DeleteArtifact : Delete artifact
// This endpoint deletes an artifact. Only administrators and editors can delete artifacts.
func (globalCatalog *GlobalCatalogV1) DeleteArtifact(deleteArtifactOptions *DeleteArtifactOptions) (response *core.DetailedResponse, err error) {
	response, err = globalCatalog.DeleteArtifactWithContext(context.Background(), deleteArtifactOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteArtifactWithContext is an alternate form of the DeleteArtifact method which supports a Context parameter
func (globalCatalog *GlobalCatalogV1) DeleteArtifactWithContext(ctx context.Context, deleteArtifactOptions *DeleteArtifactOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteArtifactOptions, "deleteArtifactOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteArtifactOptions, "deleteArtifactOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"object_id":   *deleteArtifactOptions.ObjectID,
		"artifact_id": *deleteArtifactOptions.ArtifactID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalCatalog.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalCatalog.Service.Options.URL, `/{object_id}/artifacts/{artifact_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteArtifactOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_catalog", "V1", "DeleteArtifact")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteArtifactOptions.Account != nil {
		builder.AddQuery("account", fmt.Sprint(*deleteArtifactOptions.Account))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = globalCatalog.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_artifact", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.3")
}

// AliasMetaData : Alias-related metadata.
type AliasMetaData struct {
	// Type of alias.
	Type *string `json:"type,omitempty"`

	// Points to the plan that this object is an alias for.
	PlanID *string `json:"plan_id,omitempty"`
}

// UnmarshalAliasMetaData unmarshals an instance of AliasMetaData from the specified map of raw messages.
func UnmarshalAliasMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AliasMetaData)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "plan_id", &obj.PlanID)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Amount : Country-specific pricing information.
type Amount struct {
	// Country.
	Country *string `json:"country,omitempty"`

	// Currency.
	Currency *string `json:"currency,omitempty"`

	// See Price for nested fields.
	Prices []Price `json:"prices,omitempty"`
}

// UnmarshalAmount unmarshals an instance of Amount from the specified map of raw messages.
func UnmarshalAmount(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Amount)
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		err = core.SDKErrorf(err, "", "country-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency", &obj.Currency)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "prices", &obj.Prices, UnmarshalPrice)
	if err != nil {
		err = core.SDKErrorf(err, "", "prices-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Artifact : Artifact Details.
type Artifact struct {
	// The name of the artifact.
	Name *string `json:"name,omitempty"`

	// The timestamp of the last update to the artifact.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// The url for the artifact.
	URL *string `json:"url,omitempty"`

	// The etag of the artifact.
	Etag *string `json:"etag,omitempty"`

	// The content length of the artifact.
	Size *int64 `json:"size,omitempty"`
}

// UnmarshalArtifact unmarshals an instance of Artifact from the specified map of raw messages.
func UnmarshalArtifact(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Artifact)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "etag", &obj.Etag)
	if err != nil {
		err = core.SDKErrorf(err, "", "etag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		err = core.SDKErrorf(err, "", "size-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Artifacts : Artifacts List.
type Artifacts struct {
	// The total number of artifacts.
	Count *int64 `json:"count,omitempty"`

	// The list of artifacts.
	Resources []Artifact `json:"resources,omitempty"`
}

// UnmarshalArtifacts unmarshals an instance of Artifacts from the specified map of raw messages.
func UnmarshalArtifacts(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Artifacts)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalArtifact)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AuditSearchResult : A paginated search result containing audit logs.
type AuditSearchResult struct {
	// The offset (origin 0) of the first resource in this page of search results.
	Offset *int64 `json:"offset,omitempty"`

	// The maximum number of resources returned in each page of search results.
	Limit *int64 `json:"limit,omitempty"`

	// The overall total number of resources in the search result set.
	Count *int64 `json:"count,omitempty"`

	// The number of resources returned in this page of search results.
	ResourceCount *int64 `json:"resource_count,omitempty"`

	// A URL for retrieving the first page of search results.
	First *string `json:"first,omitempty"`

	// A URL for retrieving the last page of search results.
	Last *string `json:"last,omitempty"`

	// A URL for retrieving the previous page of search results.
	Prev *string `json:"prev,omitempty"`

	// A URL for retrieving the next page of search results.
	Next *string `json:"next,omitempty"`

	// The resources (audit messages) contained in this page of search results.
	Resources []Message `json:"resources,omitempty"`
}

// UnmarshalAuditSearchResult unmarshals an instance of AuditSearchResult from the specified map of raw messages.
func UnmarshalAuditSearchResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AuditSearchResult)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_count", &obj.ResourceCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last", &obj.Last)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "prev", &obj.Prev)
	if err != nil {
		err = core.SDKErrorf(err, "", "prev-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Broker : The broker associated with a catalog entry.
type Broker struct {
	// Broker name.
	Name *string `json:"name,omitempty"`

	// Broker guid.
	GUID *string `json:"guid,omitempty"`
}

// UnmarshalBroker unmarshals an instance of Broker from the specified map of raw messages.
func UnmarshalBroker(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Broker)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "guid", &obj.GUID)
	if err != nil {
		err = core.SDKErrorf(err, "", "guid-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Bullets : Information related to list delimiters.
type Bullets struct {
	// The bullet title.
	Title *string `json:"title,omitempty"`

	// The bullet description.
	Description *string `json:"description,omitempty"`

	// The icon to use for rendering the bullet.
	Icon *string `json:"icon,omitempty"`

	// The bullet quantity.
	Quantity *int64 `json:"quantity,omitempty"`
}

// UnmarshalBullets unmarshals an instance of Bullets from the specified map of raw messages.
func UnmarshalBullets(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Bullets)
	err = core.UnmarshalPrimitive(m, "title", &obj.Title)
	if err != nil {
		err = core.SDKErrorf(err, "", "title-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "icon", &obj.Icon)
	if err != nil {
		err = core.SDKErrorf(err, "", "icon-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "quantity", &obj.Quantity)
	if err != nil {
		err = core.SDKErrorf(err, "", "quantity-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CfMetaData : Service-related metadata.
type CfMetaData struct {
	// Type of service.
	Type *string `json:"type,omitempty"`

	// Boolean value that describes whether the service is compatible with Identity and Access Management.
	IamCompatible *bool `json:"iam_compatible,omitempty"`

	// Boolean value that describes whether the service has a unique API key.
	UniqueAPIKey *bool `json:"unique_api_key,omitempty"`

	// Boolean value that describes whether the service is provisionable or not. You may need sales or support to create
	// this service.
	Provisionable *bool `json:"provisionable,omitempty"`

	// Boolean value that describes whether you can create bindings for this service.
	Bindable *bool `json:"bindable,omitempty"`

	// Boolean value that describes whether the service supports asynchronous provisioning.
	AsyncProvisioningSupported *bool `json:"async_provisioning_supported,omitempty"`

	// Boolean value that describes whether the service supports asynchronous unprovisioning.
	AsyncUnprovisioningSupported *bool `json:"async_unprovisioning_supported,omitempty"`

	// Service dependencies.
	Requires []string `json:"requires,omitempty"`

	// Boolean value that describes whether the service supports upgrade or downgrade for some plans.
	PlanUpdateable *bool `json:"plan_updateable,omitempty"`

	// String that describes whether the service is active or inactive.
	State *string `json:"state,omitempty"`

	// Boolean value that describes whether the service check is enabled.
	ServiceCheckEnabled *bool `json:"service_check_enabled,omitempty"`

	// Test check interval.
	TestCheckInterval *int64 `json:"test_check_interval,omitempty"`

	// Boolean value that describes whether the service supports service keys.
	ServiceKeySupported *bool `json:"service_key_supported,omitempty"`

	// If the field is imported from Cloud Foundry, the Cloud Foundry region's GUID. This is a required field. For example,
	// `us-south=123`.
	CfGUID map[string]string `json:"cf_guid,omitempty"`
}

// UnmarshalCfMetaData unmarshals an instance of CfMetaData from the specified map of raw messages.
func UnmarshalCfMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CfMetaData)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_compatible", &obj.IamCompatible)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_compatible-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "unique_api_key", &obj.UniqueAPIKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "unique_api_key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provisionable", &obj.Provisionable)
	if err != nil {
		err = core.SDKErrorf(err, "", "provisionable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bindable", &obj.Bindable)
	if err != nil {
		err = core.SDKErrorf(err, "", "bindable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "async_provisioning_supported", &obj.AsyncProvisioningSupported)
	if err != nil {
		err = core.SDKErrorf(err, "", "async_provisioning_supported-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "async_unprovisioning_supported", &obj.AsyncUnprovisioningSupported)
	if err != nil {
		err = core.SDKErrorf(err, "", "async_unprovisioning_supported-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "requires", &obj.Requires)
	if err != nil {
		err = core.SDKErrorf(err, "", "requires-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "plan_updateable", &obj.PlanUpdateable)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_updateable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_check_enabled", &obj.ServiceCheckEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_check_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "test_check_interval", &obj.TestCheckInterval)
	if err != nil {
		err = core.SDKErrorf(err, "", "test_check_interval-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_key_supported", &obj.ServiceKeySupported)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_key_supported-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cf_guid", &obj.CfGUID)
	if err != nil {
		err = core.SDKErrorf(err, "", "cf_guid-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Callbacks : Callback-related information associated with a catalog entry.
type Callbacks struct {
	// The URL of the deployment controller.
	ControllerURL *string `json:"controller_url,omitempty"`

	// The URL of the deployment broker.
	BrokerURL *string `json:"broker_url,omitempty"`

	// The URL of the deployment broker SC proxy.
	BrokerProxyURL *string `json:"broker_proxy_url,omitempty"`

	// The URL of dashboard callback.
	DashboardURL *string `json:"dashboard_url,omitempty"`

	// The URL of dashboard data.
	DashboardDataURL *string `json:"dashboard_data_url,omitempty"`

	// The URL of the dashboard detail tab.
	DashboardDetailTabURL *string `json:"dashboard_detail_tab_url,omitempty"`

	// The URL of the dashboard detail tab extension.
	DashboardDetailTabExtURL *string `json:"dashboard_detail_tab_ext_url,omitempty"`

	// Service monitor API URL.
	ServiceMonitorAPI *string `json:"service_monitor_api,omitempty"`

	// Service monitor app URL.
	ServiceMonitorApp *string `json:"service_monitor_app,omitempty"`

	// API endpoint.
	APIEndpoint map[string]string `json:"api_endpoint,omitempty"`
}

// UnmarshalCallbacks unmarshals an instance of Callbacks from the specified map of raw messages.
func UnmarshalCallbacks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Callbacks)
	err = core.UnmarshalPrimitive(m, "controller_url", &obj.ControllerURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "controller_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "broker_url", &obj.BrokerURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "broker_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "broker_proxy_url", &obj.BrokerProxyURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "broker_proxy_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dashboard_url", &obj.DashboardURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "dashboard_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dashboard_data_url", &obj.DashboardDataURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "dashboard_data_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dashboard_detail_tab_url", &obj.DashboardDetailTabURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "dashboard_detail_tab_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dashboard_detail_tab_ext_url", &obj.DashboardDetailTabExtURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "dashboard_detail_tab_ext_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_monitor_api", &obj.ServiceMonitorAPI)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_monitor_api-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_monitor_app", &obj.ServiceMonitorApp)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_monitor_app-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_endpoint", &obj.APIEndpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_endpoint-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CatalogEntry : An entry in the global catalog.
type CatalogEntry struct {
	// Programmatic name for this catalog entry, which must be formatted like a CRN segment. See the display name in
	// OverviewUI for a user-readable name.
	Name *string `json:"name" validate:"required"`

	// The type of catalog entry, **service**, **template**, **dashboard**, which determines the type and shape of the
	// object.
	Kind *string `json:"kind" validate:"required"`

	// Overview is nested in the top level. The key value pair is `[_language_]overview_ui`.
	OverviewUI map[string]Overview `json:"overview_ui" validate:"required"`

	// Image annotation for this catalog entry. The image is a URL.
	Images *Image `json:"images" validate:"required"`

	// The ID of the parent catalog entry if it exists.
	ParentID *string `json:"parent_id,omitempty"`

	// Boolean value that determines the global visibility for the catalog entry, and its children. If it is not enabled,
	// all plans are disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// A list of tags. For example, IBM, 3rd Party, Beta, GA, and Single Tenant.
	Tags []string `json:"tags" validate:"required"`

	// Boolean value that determines whether the catalog entry is a group.
	Group *bool `json:"group,omitempty"`

	// Information related to the provider associated with a catalog entry.
	Provider *Provider `json:"provider" validate:"required"`

	// Boolean value that describes whether the service is active.
	Active *bool `json:"active,omitempty"`

	// URL to get details about this object.
	URL *string `json:"url,omitempty"`

	// Model used to describe metadata object returned.
	Metadata *CatalogEntryMetadata `json:"metadata,omitempty"`

	// Catalog entry's unique ID. It's the same across all catalog instances.
	ID *string `json:"id,omitempty"`

	// The CRN associated with the catalog entry.
	CatalogCRN *string `json:"catalog_crn,omitempty"`

	// URL to get details about children of this object.
	ChildrenURL *string `json:"children_url,omitempty"`

	// tags to indicate the locations this service is deployable to.
	GeoTags []string `json:"geo_tags,omitempty"`

	// tags to indicate the type of pricing plans this service supports.
	PricingTags []string `json:"pricing_tags,omitempty"`

	// Date created.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// Date last updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`
}

// Constants associated with the CatalogEntry.Kind property.
// The type of catalog entry, **service**, **template**, **dashboard**, which determines the type and shape of the
// object.
const (
	CatalogEntryKindDashboardConst = "dashboard"
	CatalogEntryKindServiceConst   = "service"
	CatalogEntryKindTemplateConst  = "template"
)

// UnmarshalCatalogEntry unmarshals an instance of CatalogEntry from the specified map of raw messages.
func UnmarshalCatalogEntry(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogEntry)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		err = core.SDKErrorf(err, "", "kind-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "overview_ui", &obj.OverviewUI, UnmarshalOverview)
	if err != nil {
		err = core.SDKErrorf(err, "", "overview_ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "images", &obj.Images, UnmarshalImage)
	if err != nil {
		err = core.SDKErrorf(err, "", "images-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parent_id", &obj.ParentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "parent_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		err = core.SDKErrorf(err, "", "group-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "provider", &obj.Provider, UnmarshalProvider)
	if err != nil {
		err = core.SDKErrorf(err, "", "provider-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		err = core.SDKErrorf(err, "", "active-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalCatalogEntryMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_crn", &obj.CatalogCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "children_url", &obj.ChildrenURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "children_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "geo_tags", &obj.GeoTags)
	if err != nil {
		err = core.SDKErrorf(err, "", "geo_tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_tags", &obj.PricingTags)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		err = core.SDKErrorf(err, "", "created-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CatalogEntryMetadata : Model used to describe metadata object returned.
type CatalogEntryMetadata struct {
	// Boolean value that describes whether the service is compatible with the Resource Controller.
	RcCompatible *bool `json:"rc_compatible,omitempty"`

	// Service-related metadata.
	Service *CfMetaData `json:"service,omitempty"`

	// Plan-related metadata.
	Plan *PlanMetaData `json:"plan,omitempty"`

	// Alias-related metadata.
	Alias *AliasMetaData `json:"alias,omitempty"`

	// Template-related metadata.
	Template *TemplateMetaData `json:"template,omitempty"`

	// Information related to the UI presentation associated with a catalog entry.
	UI *UIMetaData `json:"ui,omitempty"`

	// Compliance information for HIPAA and PCI.
	Compliance []string `json:"compliance,omitempty"`

	// Service Level Agreement related metadata.
	SLA *SLAMetaData `json:"sla,omitempty"`

	// Callback-related information associated with a catalog entry.
	Callbacks *Callbacks `json:"callbacks,omitempty"`

	// The original name of the object.
	OriginalName *string `json:"original_name,omitempty"`

	// Optional version of the object.
	Version *string `json:"version,omitempty"`

	// Additional information.
	Other map[string]interface{} `json:"other,omitempty"`

	// Pricing-related information.
	Pricing *CatalogEntryMetadataPricing `json:"pricing,omitempty"`

	// Deployment-related metadata.
	Deployment *CatalogEntryMetadataDeployment `json:"deployment,omitempty"`
}

// UnmarshalCatalogEntryMetadata unmarshals an instance of CatalogEntryMetadata from the specified map of raw messages.
func UnmarshalCatalogEntryMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogEntryMetadata)
	err = core.UnmarshalPrimitive(m, "rc_compatible", &obj.RcCompatible)
	if err != nil {
		err = core.SDKErrorf(err, "", "rc_compatible-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service", &obj.Service, UnmarshalCfMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "service-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "plan", &obj.Plan, UnmarshalPlanMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "alias", &obj.Alias, UnmarshalAliasMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "alias-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalTemplateMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "template-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ui", &obj.UI, UnmarshalUIMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "compliance", &obj.Compliance)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "sla", &obj.SLA, UnmarshalSLAMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "sla-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "callbacks", &obj.Callbacks, UnmarshalCallbacks)
	if err != nil {
		err = core.SDKErrorf(err, "", "callbacks-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "original_name", &obj.OriginalName)
	if err != nil {
		err = core.SDKErrorf(err, "", "original_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "other", &obj.Other)
	if err != nil {
		err = core.SDKErrorf(err, "", "other-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pricing", &obj.Pricing, UnmarshalCatalogEntryMetadataPricing)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "deployment", &obj.Deployment, UnmarshalCatalogEntryMetadataDeployment)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CatalogEntryMetadataDeployment : Deployment-related metadata.
type CatalogEntryMetadataDeployment struct {
	// Describes the region where the service is located.
	Location *string `json:"location,omitempty"`

	// Pointer to the location resource in the catalog.
	LocationURL *string `json:"location_url,omitempty"`

	// Original service location.
	OriginalLocation *string `json:"original_location,omitempty"`

	// A CRN that describes the deployment. crn:v1:[cname]:[ctype]:[location]:[scope]::[resource-type]:[resource].
	TargetCRN *string `json:"target_crn,omitempty"`

	// CRN for the service.
	ServiceCRN *string `json:"service_crn,omitempty"`

	// ID for MCCP.
	MccpID *string `json:"mccp_id,omitempty"`

	// The broker associated with a catalog entry.
	Broker *Broker `json:"broker,omitempty"`

	// This deployment not only supports RC but is ready to migrate and support the RC broker for a location.
	SupportsRcMigration *bool `json:"supports_rc_migration,omitempty"`

	// network to use during deployment.
	TargetNetwork *string `json:"target_network,omitempty"`
}

// UnmarshalCatalogEntryMetadataDeployment unmarshals an instance of CatalogEntryMetadataDeployment from the specified map of raw messages.
func UnmarshalCatalogEntryMetadataDeployment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogEntryMetadataDeployment)
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location_url", &obj.LocationURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "location_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "original_location", &obj.OriginalLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "original_location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mccp_id", &obj.MccpID)
	if err != nil {
		err = core.SDKErrorf(err, "", "mccp_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "broker", &obj.Broker, UnmarshalBroker)
	if err != nil {
		err = core.SDKErrorf(err, "", "broker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "supports_rc_migration", &obj.SupportsRcMigration)
	if err != nil {
		err = core.SDKErrorf(err, "", "supports_rc_migration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_network", &obj.TargetNetwork)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_network-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CatalogEntryMetadataPricing : Pricing-related information.
type CatalogEntryMetadataPricing struct {
	// Type of plan. Valid values are `free`, `trial`, `paygo`, `bluemix-subscription`, and `ibm-subscription`.
	Type *string `json:"type,omitempty"`

	// Defines where the pricing originates.
	Origin *string `json:"origin,omitempty"`

	// Plan-specific starting price information.
	StartingPrice *StartingPrice `json:"starting_price,omitempty"`

	// The deployment object id this pricing is from. Only set if object kind is deployment.
	DeploymentID *string `json:"deployment_id,omitempty"`

	// The deployment location this pricing is from. Only set if object kind is deployment.
	DeploymentLocation *string `json:"deployment_location,omitempty"`

	// Is the location price not available. Only set in api /pricing/deployment and only set if true. This means for the
	// given deployment object there was no pricing set in pricing catalog.
	DeploymentLocationNoPriceAvailable *bool `json:"deployment_location_no_price_available,omitempty"`

	// Plan-specific cost metric structure.
	Metrics []Metrics `json:"metrics,omitempty"`

	// List of regions where region pricing is available. Only set on global deployments if enabled by owner.
	DeploymentRegions []string `json:"deployment_regions,omitempty"`
}

// UnmarshalCatalogEntryMetadataPricing unmarshals an instance of CatalogEntryMetadataPricing from the specified map of raw messages.
func UnmarshalCatalogEntryMetadataPricing(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogEntryMetadataPricing)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "origin", &obj.Origin)
	if err != nil {
		err = core.SDKErrorf(err, "", "origin-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "starting_price", &obj.StartingPrice, UnmarshalStartingPrice)
	if err != nil {
		err = core.SDKErrorf(err, "", "starting_price-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_id", &obj.DeploymentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_location", &obj.DeploymentLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_location_no_price_available", &obj.DeploymentLocationNoPriceAvailable)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_location_no_price_available-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metrics", &obj.Metrics, UnmarshalMetrics)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_regions", &obj.DeploymentRegions)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_regions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCatalogEntryOptions : The CreateCatalogEntry options.
type CreateCatalogEntryOptions struct {
	// Programmatic name for this catalog entry, which must be formatted like a CRN segment. See the display name in
	// OverviewUI for a user-readable name.
	Name *string `json:"name" validate:"required"`

	// The type of catalog entry, **service**, **template**, **dashboard**, which determines the type and shape of the
	// object.
	Kind *string `json:"kind" validate:"required"`

	// Overview is nested in the top level. The key value pair is `[_language_]overview_ui`.
	OverviewUI map[string]Overview `json:"overview_ui" validate:"required"`

	// Image annotation for this catalog entry. The image is a URL.
	Images *Image `json:"images" validate:"required"`

	// Boolean value that determines the global visibility for the catalog entry, and its children. If it is not enabled,
	// all plans are disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// A list of tags. For example, IBM, 3rd Party, Beta, GA, and Single Tenant.
	Tags []string `json:"tags" validate:"required"`

	// Information related to the provider associated with a catalog entry.
	Provider *Provider `json:"provider" validate:"required"`

	// Catalog entry's unique ID. It's the same across all catalog instances.
	ID *string `json:"id" validate:"required"`

	// The ID of the parent catalog entry if it exists.
	ParentID *string `json:"parent_id,omitempty"`

	// Boolean value that determines whether the catalog entry is a group.
	Group *bool `json:"group,omitempty"`

	// Boolean value that describes whether the service is active.
	Active *bool `json:"active,omitempty"`

	// Url of the object.
	URL *string `json:"url,omitempty"`

	// Model used to describe metadata object that can be set.
	Metadata *ObjectMetadataSet `json:"metadata,omitempty"`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateCatalogEntryOptions.Kind property.
// The type of catalog entry, **service**, **template**, **dashboard**, which determines the type and shape of the
// object.
const (
	CreateCatalogEntryOptionsKindDashboardConst = "dashboard"
	CreateCatalogEntryOptionsKindServiceConst   = "service"
	CreateCatalogEntryOptionsKindTemplateConst  = "template"
)

// NewCreateCatalogEntryOptions : Instantiate CreateCatalogEntryOptions
func (*GlobalCatalogV1) NewCreateCatalogEntryOptions(name string, kind string, overviewUI map[string]Overview, images *Image, disabled bool, tags []string, provider *Provider, id string) *CreateCatalogEntryOptions {
	return &CreateCatalogEntryOptions{
		Name:       core.StringPtr(name),
		Kind:       core.StringPtr(kind),
		OverviewUI: overviewUI,
		Images:     images,
		Disabled:   core.BoolPtr(disabled),
		Tags:       tags,
		Provider:   provider,
		ID:         core.StringPtr(id),
	}
}

// SetName : Allow user to set Name
func (_options *CreateCatalogEntryOptions) SetName(name string) *CreateCatalogEntryOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *CreateCatalogEntryOptions) SetKind(kind string) *CreateCatalogEntryOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetOverviewUI : Allow user to set OverviewUI
func (_options *CreateCatalogEntryOptions) SetOverviewUI(overviewUI map[string]Overview) *CreateCatalogEntryOptions {
	_options.OverviewUI = overviewUI
	return _options
}

// SetImages : Allow user to set Images
func (_options *CreateCatalogEntryOptions) SetImages(images *Image) *CreateCatalogEntryOptions {
	_options.Images = images
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *CreateCatalogEntryOptions) SetDisabled(disabled bool) *CreateCatalogEntryOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateCatalogEntryOptions) SetTags(tags []string) *CreateCatalogEntryOptions {
	_options.Tags = tags
	return _options
}

// SetProvider : Allow user to set Provider
func (_options *CreateCatalogEntryOptions) SetProvider(provider *Provider) *CreateCatalogEntryOptions {
	_options.Provider = provider
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateCatalogEntryOptions) SetID(id string) *CreateCatalogEntryOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetParentID : Allow user to set ParentID
func (_options *CreateCatalogEntryOptions) SetParentID(parentID string) *CreateCatalogEntryOptions {
	_options.ParentID = core.StringPtr(parentID)
	return _options
}

// SetGroup : Allow user to set Group
func (_options *CreateCatalogEntryOptions) SetGroup(group bool) *CreateCatalogEntryOptions {
	_options.Group = core.BoolPtr(group)
	return _options
}

// SetActive : Allow user to set Active
func (_options *CreateCatalogEntryOptions) SetActive(active bool) *CreateCatalogEntryOptions {
	_options.Active = core.BoolPtr(active)
	return _options
}

// SetURL : Allow user to set URL
func (_options *CreateCatalogEntryOptions) SetURL(url string) *CreateCatalogEntryOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateCatalogEntryOptions) SetMetadata(metadata *ObjectMetadataSet) *CreateCatalogEntryOptions {
	_options.Metadata = metadata
	return _options
}

// SetAccount : Allow user to set Account
func (_options *CreateCatalogEntryOptions) SetAccount(account string) *CreateCatalogEntryOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCatalogEntryOptions) SetHeaders(param map[string]string) *CreateCatalogEntryOptions {
	options.Headers = param
	return options
}

// DrMetaData : SLA Disaster Recovery-related metadata.
type DrMetaData struct {
	// Required boolean value that describes whether disaster recovery is on.
	Dr *bool `json:"dr,omitempty"`

	// Description of the disaster recovery implementation.
	Description *string `json:"description,omitempty"`
}

// UnmarshalDrMetaData unmarshals an instance of DrMetaData from the specified map of raw messages.
func UnmarshalDrMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DrMetaData)
	err = core.UnmarshalPrimitive(m, "dr", &obj.Dr)
	if err != nil {
		err = core.SDKErrorf(err, "", "dr-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteArtifactOptions : The DeleteArtifact options.
type DeleteArtifactOptions struct {
	// The object's unique ID.
	ObjectID *string `json:"object_id" validate:"required,ne="`

	// The artifact's ID.
	ArtifactID *string `json:"artifact_id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteArtifactOptions : Instantiate DeleteArtifactOptions
func (*GlobalCatalogV1) NewDeleteArtifactOptions(objectID string, artifactID string) *DeleteArtifactOptions {
	return &DeleteArtifactOptions{
		ObjectID:   core.StringPtr(objectID),
		ArtifactID: core.StringPtr(artifactID),
	}
}

// SetObjectID : Allow user to set ObjectID
func (_options *DeleteArtifactOptions) SetObjectID(objectID string) *DeleteArtifactOptions {
	_options.ObjectID = core.StringPtr(objectID)
	return _options
}

// SetArtifactID : Allow user to set ArtifactID
func (_options *DeleteArtifactOptions) SetArtifactID(artifactID string) *DeleteArtifactOptions {
	_options.ArtifactID = core.StringPtr(artifactID)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *DeleteArtifactOptions) SetAccount(account string) *DeleteArtifactOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteArtifactOptions) SetHeaders(param map[string]string) *DeleteArtifactOptions {
	options.Headers = param
	return options
}

// DeleteCatalogEntryOptions : The DeleteCatalogEntry options.
type DeleteCatalogEntryOptions struct {
	// The object's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// This will cause entry to be deleted fully. By default it is archived for two weeks, so that it can be restored if
	// necessary.
	Force *bool `json:"force,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteCatalogEntryOptions : Instantiate DeleteCatalogEntryOptions
func (*GlobalCatalogV1) NewDeleteCatalogEntryOptions(id string) *DeleteCatalogEntryOptions {
	return &DeleteCatalogEntryOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteCatalogEntryOptions) SetID(id string) *DeleteCatalogEntryOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *DeleteCatalogEntryOptions) SetAccount(account string) *DeleteCatalogEntryOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetForce : Allow user to set Force
func (_options *DeleteCatalogEntryOptions) SetForce(force bool) *DeleteCatalogEntryOptions {
	_options.Force = core.BoolPtr(force)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCatalogEntryOptions) SetHeaders(param map[string]string) *DeleteCatalogEntryOptions {
	options.Headers = param
	return options
}

// DeploymentBase : Deployment-related metadata.
type DeploymentBase struct {
	// Describes the region where the service is located.
	Location *string `json:"location,omitempty"`

	// URL of deployment.
	LocationURL *string `json:"location_url,omitempty"`

	// Original service location.
	OriginalLocation *string `json:"original_location,omitempty"`

	// A CRN that describes the deployment. crn:v1:[cname]:[ctype]:[location]:[scope]::[resource-type]:[resource].
	TargetCRN *string `json:"target_crn,omitempty"`

	// CRN for the service.
	ServiceCRN *string `json:"service_crn,omitempty"`

	// ID for MCCP.
	MccpID *string `json:"mccp_id,omitempty"`

	// The broker associated with a catalog entry.
	Broker *Broker `json:"broker,omitempty"`

	// This deployment not only supports RC but is ready to migrate and support the RC broker for a location.
	SupportsRcMigration *bool `json:"supports_rc_migration,omitempty"`

	// network to use during deployment.
	TargetNetwork *string `json:"target_network,omitempty"`
}

// UnmarshalDeploymentBase unmarshals an instance of DeploymentBase from the specified map of raw messages.
func UnmarshalDeploymentBase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeploymentBase)
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location_url", &obj.LocationURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "location_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "original_location", &obj.OriginalLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "original_location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "mccp_id", &obj.MccpID)
	if err != nil {
		err = core.SDKErrorf(err, "", "mccp_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "broker", &obj.Broker, UnmarshalBroker)
	if err != nil {
		err = core.SDKErrorf(err, "", "broker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "supports_rc_migration", &obj.SupportsRcMigration)
	if err != nil {
		err = core.SDKErrorf(err, "", "supports_rc_migration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_network", &obj.TargetNetwork)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_network-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EntrySearchResult : A paginated search result containing catalog entries.
type EntrySearchResult struct {
	// The offset (origin 0) of the first resource in this page of search results.
	Offset *int64 `json:"offset,omitempty"`

	// The maximum number of resources returned in each page of search results.
	Limit *int64 `json:"limit,omitempty"`

	// The overall total number of resources in the search result set.
	Count *int64 `json:"count,omitempty"`

	// The number of resources returned in this page of search results.
	ResourceCount *int64 `json:"resource_count,omitempty"`

	// A URL for retrieving the first page of search results.
	First *string `json:"first,omitempty"`

	// A URL for retrieving the last page of search results.
	Last *string `json:"last,omitempty"`

	// A URL for retrieving the previous page of search results.
	Prev *string `json:"prev,omitempty"`

	// A URL for retrieving the next page of search results.
	Next *string `json:"next,omitempty"`

	// The resources (catalog entries) contained in this page of search results.
	Resources []CatalogEntry `json:"resources,omitempty"`
}

// UnmarshalEntrySearchResult unmarshals an instance of EntrySearchResult from the specified map of raw messages.
func UnmarshalEntrySearchResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EntrySearchResult)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_count", &obj.ResourceCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last", &obj.Last)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "prev", &obj.Prev)
	if err != nil {
		err = core.SDKErrorf(err, "", "prev-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalCatalogEntry)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetArtifactOptions : The GetArtifact options.
type GetArtifactOptions struct {
	// The object's unique ID.
	ObjectID *string `json:"object_id" validate:"required,ne="`

	// The artifact's ID.
	ArtifactID *string `json:"artifact_id" validate:"required,ne="`

	// The type of the response: *_/_*.
	Accept *string `json:"Accept,omitempty"`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetArtifactOptions : Instantiate GetArtifactOptions
func (*GlobalCatalogV1) NewGetArtifactOptions(objectID string, artifactID string) *GetArtifactOptions {
	return &GetArtifactOptions{
		ObjectID:   core.StringPtr(objectID),
		ArtifactID: core.StringPtr(artifactID),
	}
}

// SetObjectID : Allow user to set ObjectID
func (_options *GetArtifactOptions) SetObjectID(objectID string) *GetArtifactOptions {
	_options.ObjectID = core.StringPtr(objectID)
	return _options
}

// SetArtifactID : Allow user to set ArtifactID
func (_options *GetArtifactOptions) SetArtifactID(artifactID string) *GetArtifactOptions {
	_options.ArtifactID = core.StringPtr(artifactID)
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *GetArtifactOptions) SetAccept(accept string) *GetArtifactOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *GetArtifactOptions) SetAccount(account string) *GetArtifactOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetArtifactOptions) SetHeaders(param map[string]string) *GetArtifactOptions {
	options.Headers = param
	return options
}

// GetAuditLogsOptions : The GetAuditLogs options.
type GetAuditLogsOptions struct {
	// The object's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Sets the sort order. False is descending.
	Ascending *string `json:"ascending,omitempty"`

	// Starting time for the logs. If it's descending then the entries will be equal or earlier. The default is latest. For
	// ascending it will entries equal or later. The default is earliest. It can be either a number or a string. If a
	// number then it is in the format of Unix timestamps. If it is a string then it is a date in the format
	// YYYY-MM-DDTHH:MM:SSZ  and the time is UTC. The T and the Z are required. For example: 2017-12-24T12:00:00Z for Noon
	// UTC on Dec 24, 2017.
	Startat *string `json:"startat,omitempty"`

	// Count of number of log entries to skip before returning logs. The default is zero.
	Offset *int64 `json:"_offset,omitempty"`

	// Count of number of entries to return. The default is fifty. The maximum value is two hundred.
	Limit *int64 `json:"_limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAuditLogsOptions : Instantiate GetAuditLogsOptions
func (*GlobalCatalogV1) NewGetAuditLogsOptions(id string) *GetAuditLogsOptions {
	return &GetAuditLogsOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetAuditLogsOptions) SetID(id string) *GetAuditLogsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *GetAuditLogsOptions) SetAccount(account string) *GetAuditLogsOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetAscending : Allow user to set Ascending
func (_options *GetAuditLogsOptions) SetAscending(ascending string) *GetAuditLogsOptions {
	_options.Ascending = core.StringPtr(ascending)
	return _options
}

// SetStartat : Allow user to set Startat
func (_options *GetAuditLogsOptions) SetStartat(startat string) *GetAuditLogsOptions {
	_options.Startat = core.StringPtr(startat)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetAuditLogsOptions) SetOffset(offset int64) *GetAuditLogsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetAuditLogsOptions) SetLimit(limit int64) *GetAuditLogsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAuditLogsOptions) SetHeaders(param map[string]string) *GetAuditLogsOptions {
	options.Headers = param
	return options
}

// GetCatalogEntryOptions : The GetCatalogEntry options.
type GetCatalogEntryOptions struct {
	// The catalog entry's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// A GET call by default returns a basic set of properties. To include other properties, you must add this parameter. A
	// wildcard (`*`) includes all properties for an object, for example `GET /id?include=*`. To include specific metadata
	// fields, separate each field with a colon (:), for example `GET /id?include=metadata.ui:metadata.pricing`.
	Include *string `json:"include,omitempty"`

	// Return the data strings in the specified language. By default the strings returned are of the language preferred by
	// your browser through the Accept-Language header, which allows an override of the header. Languages are specified in
	// standard form, such as `en-us`. To include all languages use a wildcard (*).
	Languages *string `json:"languages,omitempty"`

	// Returns all available fields for all languages. Use the value `?complete=true` as shortcut for
	// ?include=*&languages=*.
	Complete *bool `json:"complete,omitempty"`

	// Return the children down to the requested depth. Use * to include the entire children tree. If there are more
	// children than the maximum permitted an error will be returned. Be judicious with this as it can cause a large number
	// of database accesses and can result in a large amount of data returned.
	Depth *int64 `json:"depth,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetCatalogEntryOptions : Instantiate GetCatalogEntryOptions
func (*GlobalCatalogV1) NewGetCatalogEntryOptions(id string) *GetCatalogEntryOptions {
	return &GetCatalogEntryOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetCatalogEntryOptions) SetID(id string) *GetCatalogEntryOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *GetCatalogEntryOptions) SetAccount(account string) *GetCatalogEntryOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *GetCatalogEntryOptions) SetInclude(include string) *GetCatalogEntryOptions {
	_options.Include = core.StringPtr(include)
	return _options
}

// SetLanguages : Allow user to set Languages
func (_options *GetCatalogEntryOptions) SetLanguages(languages string) *GetCatalogEntryOptions {
	_options.Languages = core.StringPtr(languages)
	return _options
}

// SetComplete : Allow user to set Complete
func (_options *GetCatalogEntryOptions) SetComplete(complete bool) *GetCatalogEntryOptions {
	_options.Complete = core.BoolPtr(complete)
	return _options
}

// SetDepth : Allow user to set Depth
func (_options *GetCatalogEntryOptions) SetDepth(depth int64) *GetCatalogEntryOptions {
	_options.Depth = core.Int64Ptr(depth)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogEntryOptions) SetHeaders(param map[string]string) *GetCatalogEntryOptions {
	options.Headers = param
	return options
}

// GetChildObjectsOptions : The GetChildObjects options.
type GetChildObjectsOptions struct {
	// The parent catalog entry's ID.
	ID *string `json:"id" validate:"required,ne="`

	// The **kind** of child catalog entries to search for. A wildcard (*) includes all child catalog entries for all
	// kinds, for example `GET /service_name/_*`.
	Kind *string `json:"kind" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// A colon (:) separated list of properties to include. A GET call by defaults return a limited set of properties. To
	// include other properties, you must add the include parameter.  A wildcard (*) includes all properties.
	Include *string `json:"include,omitempty"`

	// A query filter, for example, `q=kind:iaas IBM`  will filter on entries of **kind** iaas that has `IBM` in their
	// name, display name, or description.
	Q *string `json:"q,omitempty"`

	// The field on which to sort the output. By default by name. Available fields are **name**, **kind**, and
	// **provider**.
	SortBy *string `json:"sort-by,omitempty"`

	// The sort order. The default is false, which is ascending.
	Descending *string `json:"descending,omitempty"`

	// Return the data strings in the specified language. By default the strings returned are of the language preferred by
	// your browser through the Accept-Language header. This allows an override of the header. Languages are specified in
	// standard form, such as `en-us`. To include all languages use the wildcard (*).
	Languages *string `json:"languages,omitempty"`

	// Use the value `?complete=true` as shortcut for ?include=*&languages=*.
	Complete *bool `json:"complete,omitempty"`

	// Useful for pagination, specifies index (origin 0) of first item to return in response.
	Offset *int64 `json:"_offset,omitempty"`

	// Useful for pagination, specifies the maximum number of items to return in the response.
	Limit *int64 `json:"_limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetChildObjectsOptions : Instantiate GetChildObjectsOptions
func (*GlobalCatalogV1) NewGetChildObjectsOptions(id string, kind string) *GetChildObjectsOptions {
	return &GetChildObjectsOptions{
		ID:   core.StringPtr(id),
		Kind: core.StringPtr(kind),
	}
}

// SetID : Allow user to set ID
func (_options *GetChildObjectsOptions) SetID(id string) *GetChildObjectsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *GetChildObjectsOptions) SetKind(kind string) *GetChildObjectsOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *GetChildObjectsOptions) SetAccount(account string) *GetChildObjectsOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *GetChildObjectsOptions) SetInclude(include string) *GetChildObjectsOptions {
	_options.Include = core.StringPtr(include)
	return _options
}

// SetQ : Allow user to set Q
func (_options *GetChildObjectsOptions) SetQ(q string) *GetChildObjectsOptions {
	_options.Q = core.StringPtr(q)
	return _options
}

// SetSortBy : Allow user to set SortBy
func (_options *GetChildObjectsOptions) SetSortBy(sortBy string) *GetChildObjectsOptions {
	_options.SortBy = core.StringPtr(sortBy)
	return _options
}

// SetDescending : Allow user to set Descending
func (_options *GetChildObjectsOptions) SetDescending(descending string) *GetChildObjectsOptions {
	_options.Descending = core.StringPtr(descending)
	return _options
}

// SetLanguages : Allow user to set Languages
func (_options *GetChildObjectsOptions) SetLanguages(languages string) *GetChildObjectsOptions {
	_options.Languages = core.StringPtr(languages)
	return _options
}

// SetComplete : Allow user to set Complete
func (_options *GetChildObjectsOptions) SetComplete(complete bool) *GetChildObjectsOptions {
	_options.Complete = core.BoolPtr(complete)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetChildObjectsOptions) SetOffset(offset int64) *GetChildObjectsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetChildObjectsOptions) SetLimit(limit int64) *GetChildObjectsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetChildObjectsOptions) SetHeaders(param map[string]string) *GetChildObjectsOptions {
	options.Headers = param
	return options
}

// GetPricingDeploymentsOptions : The GetPricingDeployments options.
type GetPricingDeploymentsOptions struct {
	// The object's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetPricingDeploymentsOptions : Instantiate GetPricingDeploymentsOptions
func (*GlobalCatalogV1) NewGetPricingDeploymentsOptions(id string) *GetPricingDeploymentsOptions {
	return &GetPricingDeploymentsOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetPricingDeploymentsOptions) SetID(id string) *GetPricingDeploymentsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *GetPricingDeploymentsOptions) SetAccount(account string) *GetPricingDeploymentsOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPricingDeploymentsOptions) SetHeaders(param map[string]string) *GetPricingDeploymentsOptions {
	options.Headers = param
	return options
}

// GetPricingOptions : The GetPricing options.
type GetPricingOptions struct {
	// The object's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Specify a region to retrieve plan pricing for a global deployment. The value must match an entry in the
	// `deployment_regions` list.
	DeploymentRegion *string `json:"deployment_region,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetPricingOptions : Instantiate GetPricingOptions
func (*GlobalCatalogV1) NewGetPricingOptions(id string) *GetPricingOptions {
	return &GetPricingOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetPricingOptions) SetID(id string) *GetPricingOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *GetPricingOptions) SetAccount(account string) *GetPricingOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetDeploymentRegion : Allow user to set DeploymentRegion
func (_options *GetPricingOptions) SetDeploymentRegion(deploymentRegion string) *GetPricingOptions {
	_options.DeploymentRegion = core.StringPtr(deploymentRegion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPricingOptions) SetHeaders(param map[string]string) *GetPricingOptions {
	options.Headers = param
	return options
}

// GetVisibilityOptions : The GetVisibility options.
type GetVisibilityOptions struct {
	// The object's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetVisibilityOptions : Instantiate GetVisibilityOptions
func (*GlobalCatalogV1) NewGetVisibilityOptions(id string) *GetVisibilityOptions {
	return &GetVisibilityOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetVisibilityOptions) SetID(id string) *GetVisibilityOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *GetVisibilityOptions) SetAccount(account string) *GetVisibilityOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetVisibilityOptions) SetHeaders(param map[string]string) *GetVisibilityOptions {
	options.Headers = param
	return options
}

// Image : Image annotation for this catalog entry. The image is a URL.
type Image struct {
	// URL for the large, default image.
	Image *string `json:"image" validate:"required"`

	// URL for a small image.
	SmallImage *string `json:"small_image,omitempty"`

	// URL for a medium image.
	MediumImage *string `json:"medium_image,omitempty"`

	// URL for a featured image.
	FeatureImage *string `json:"feature_image,omitempty"`
}

// NewImage : Instantiate Image (Generic Model Constructor)
func (*GlobalCatalogV1) NewImage(image string) (_model *Image, err error) {
	_model = &Image{
		Image: core.StringPtr(image),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalImage unmarshals an instance of Image from the specified map of raw messages.
func UnmarshalImage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Image)
	err = core.UnmarshalPrimitive(m, "image", &obj.Image)
	if err != nil {
		err = core.SDKErrorf(err, "", "image-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "small_image", &obj.SmallImage)
	if err != nil {
		err = core.SDKErrorf(err, "", "small_image-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "medium_image", &obj.MediumImage)
	if err != nil {
		err = core.SDKErrorf(err, "", "medium_image-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "feature_image", &obj.FeatureImage)
	if err != nil {
		err = core.SDKErrorf(err, "", "feature_image-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListArtifactsOptions : The ListArtifacts options.
type ListArtifactsOptions struct {
	// The object's unique ID.
	ObjectID *string `json:"object_id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListArtifactsOptions : Instantiate ListArtifactsOptions
func (*GlobalCatalogV1) NewListArtifactsOptions(objectID string) *ListArtifactsOptions {
	return &ListArtifactsOptions{
		ObjectID: core.StringPtr(objectID),
	}
}

// SetObjectID : Allow user to set ObjectID
func (_options *ListArtifactsOptions) SetObjectID(objectID string) *ListArtifactsOptions {
	_options.ObjectID = core.StringPtr(objectID)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *ListArtifactsOptions) SetAccount(account string) *ListArtifactsOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListArtifactsOptions) SetHeaders(param map[string]string) *ListArtifactsOptions {
	options.Headers = param
	return options
}

// ListCatalogEntriesOptions : The ListCatalogEntries options.
type ListCatalogEntriesOptions struct {
	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// A GET call by default returns a basic set of properties. To include other properties, you must add this parameter. A
	// wildcard (`*`) includes all properties for an object, for example `GET /?include=*`. To include specific metadata
	// fields, separate each field with a colon (:), for example `GET /?include=metadata.ui:metadata.pricing`.
	Include *string `json:"include,omitempty"`

	// Searches the catalog entries for keywords. Add filters to refine your search. A query filter, for example,
	// `q=kind:iaas service_name rc:true`, filters entries of kind iaas with metadata.service.rc_compatible set to true and
	//  have a service name is in their name, display name, or description.  Valid tags are **kind**:<string>,
	// **tag**:<strging>, **rc**:[true|false], **iam**:[true|false], **active**:[true|false], **geo**:<string>, and
	// **price**:<string>.
	Q *string `json:"q,omitempty"`

	// The field on which the output is sorted. Sorts by default by **name** property. Available fields are **name**,
	// **displayname** (overview_ui.display_name), **kind**, **provider** (provider.name), **sbsindex**
	// (metadata.ui.side_by_side_index), and the time **created**, and **updated**.
	SortBy *string `json:"sort-by,omitempty"`

	// Sets the sort order. The default is false, which is ascending.
	Descending *string `json:"descending,omitempty"`

	// Return the data strings in a specified language. By default, the strings returned are of the language preferred by
	// your browser through the Accept-Language header, which allows an override of the header. Languages are specified in
	// standard form, such as `en-us`. To include all languages use a wildcard (*).
	Languages *string `json:"languages,omitempty"`

	// Checks to see if a catalog's object is visible, or if it's filtered by service, plan, deployment, or region. Use the
	// value `?catalog=true`. If a `200` code is returned, the object is visible. If a `403` code is returned, the object
	// is not visible for the user.
	Catalog *bool `json:"catalog,omitempty"`

	// Returns all available fields for all languages. Use the value `?complete=true` as shortcut for
	// ?include=*&languages=*.
	Complete *bool `json:"complete,omitempty"`

	// Useful for pagination, specifies index (origin 0) of first item to return in response.
	Offset *int64 `json:"_offset,omitempty"`

	// Useful for pagination, specifies the maximum number of items to return in the response.
	Limit *int64 `json:"_limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListCatalogEntriesOptions : Instantiate ListCatalogEntriesOptions
func (*GlobalCatalogV1) NewListCatalogEntriesOptions() *ListCatalogEntriesOptions {
	return &ListCatalogEntriesOptions{}
}

// SetAccount : Allow user to set Account
func (_options *ListCatalogEntriesOptions) SetAccount(account string) *ListCatalogEntriesOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *ListCatalogEntriesOptions) SetInclude(include string) *ListCatalogEntriesOptions {
	_options.Include = core.StringPtr(include)
	return _options
}

// SetQ : Allow user to set Q
func (_options *ListCatalogEntriesOptions) SetQ(q string) *ListCatalogEntriesOptions {
	_options.Q = core.StringPtr(q)
	return _options
}

// SetSortBy : Allow user to set SortBy
func (_options *ListCatalogEntriesOptions) SetSortBy(sortBy string) *ListCatalogEntriesOptions {
	_options.SortBy = core.StringPtr(sortBy)
	return _options
}

// SetDescending : Allow user to set Descending
func (_options *ListCatalogEntriesOptions) SetDescending(descending string) *ListCatalogEntriesOptions {
	_options.Descending = core.StringPtr(descending)
	return _options
}

// SetLanguages : Allow user to set Languages
func (_options *ListCatalogEntriesOptions) SetLanguages(languages string) *ListCatalogEntriesOptions {
	_options.Languages = core.StringPtr(languages)
	return _options
}

// SetCatalog : Allow user to set Catalog
func (_options *ListCatalogEntriesOptions) SetCatalog(catalog bool) *ListCatalogEntriesOptions {
	_options.Catalog = core.BoolPtr(catalog)
	return _options
}

// SetComplete : Allow user to set Complete
func (_options *ListCatalogEntriesOptions) SetComplete(complete bool) *ListCatalogEntriesOptions {
	_options.Complete = core.BoolPtr(complete)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListCatalogEntriesOptions) SetOffset(offset int64) *ListCatalogEntriesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListCatalogEntriesOptions) SetLimit(limit int64) *ListCatalogEntriesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListCatalogEntriesOptions) SetHeaders(param map[string]string) *ListCatalogEntriesOptions {
	options.Headers = param
	return options
}

// Message : log object describing who did what.
type Message struct {
	// id of catalog entry.
	ID *string `json:"id,omitempty"`

	// Information related to the visibility of a catalog entry.
	Effective *Visibility `json:"effective,omitempty"`

	// time of action.
	Time *strfmt.DateTime `json:"time,omitempty"`

	// user ID of person who did action.
	WhoID *string `json:"who_id,omitempty"`

	// name of person who did action.
	WhoName *string `json:"who_name,omitempty"`

	// user email of person who did action.
	WhoEmail *string `json:"who_email,omitempty"`

	// Global catalog instance where this occured.
	Instance *string `json:"instance,omitempty"`

	// transaction id associatd with action.
	Gid *string `json:"gid,omitempty"`

	// type of action taken.
	Type *string `json:"type,omitempty"`

	// message describing action.
	Message *string `json:"message,omitempty"`

	// An object containing details on changes made to object data.
	Data map[string]interface{} `json:"data,omitempty"`
}

// UnmarshalMessage unmarshals an instance of Message from the specified map of raw messages.
func UnmarshalMessage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Message)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "effective", &obj.Effective, UnmarshalVisibility)
	if err != nil {
		err = core.SDKErrorf(err, "", "effective-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "time", &obj.Time)
	if err != nil {
		err = core.SDKErrorf(err, "", "time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "who_id", &obj.WhoID)
	if err != nil {
		err = core.SDKErrorf(err, "", "who_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "who_name", &obj.WhoName)
	if err != nil {
		err = core.SDKErrorf(err, "", "who_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "who_email", &obj.WhoEmail)
	if err != nil {
		err = core.SDKErrorf(err, "", "who_email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instance", &obj.Instance)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "gid", &obj.Gid)
	if err != nil {
		err = core.SDKErrorf(err, "", "gid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Metrics : Plan-specific cost metrics information.
type Metrics struct {
	// The part reference.
	PartRef *string `json:"part_ref,omitempty"`

	// The metric ID or part number.
	MetricID *string `json:"metric_id,omitempty"`

	// The tier model.
	TierModel *string `json:"tier_model,omitempty"`

	// The unit to charge.
	ChargeUnit *string `json:"charge_unit,omitempty"`

	// The charge unit name.
	ChargeUnitName *string `json:"charge_unit_name,omitempty"`

	// The charge unit quantity.
	ChargeUnitQuantity *int64 `json:"charge_unit_quantity,omitempty"`

	// Display name of the resource.
	ResourceDisplayName *string `json:"resource_display_name,omitempty"`

	// Display name of the charge unit.
	ChargeUnitDisplayName *string `json:"charge_unit_display_name,omitempty"`

	// Usage limit for the metric.
	UsageCapQty *int64 `json:"usage_cap_qty,omitempty"`

	// Display capacity.
	DisplayCap *int64 `json:"display_cap,omitempty"`

	// Effective from time.
	EffectiveFrom *strfmt.DateTime `json:"effective_from,omitempty"`

	// Effective until time.
	EffectiveUntil *strfmt.DateTime `json:"effective_until,omitempty"`

	// The pricing per metric by country and currency.
	Amounts []Amount `json:"amounts,omitempty"`
}

// UnmarshalMetrics unmarshals an instance of Metrics from the specified map of raw messages.
func UnmarshalMetrics(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Metrics)
	err = core.UnmarshalPrimitive(m, "part_ref", &obj.PartRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "part_ref-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metric_id", &obj.MetricID)
	if err != nil {
		err = core.SDKErrorf(err, "", "metric_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tier_model", &obj.TierModel)
	if err != nil {
		err = core.SDKErrorf(err, "", "tier_model-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "charge_unit", &obj.ChargeUnit)
	if err != nil {
		err = core.SDKErrorf(err, "", "charge_unit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "charge_unit_name", &obj.ChargeUnitName)
	if err != nil {
		err = core.SDKErrorf(err, "", "charge_unit_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "charge_unit_quantity", &obj.ChargeUnitQuantity)
	if err != nil {
		err = core.SDKErrorf(err, "", "charge_unit_quantity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_display_name", &obj.ResourceDisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "charge_unit_display_name", &obj.ChargeUnitDisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "charge_unit_display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "usage_cap_qty", &obj.UsageCapQty)
	if err != nil {
		err = core.SDKErrorf(err, "", "usage_cap_qty-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "display_cap", &obj.DisplayCap)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_cap-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "effective_from", &obj.EffectiveFrom)
	if err != nil {
		err = core.SDKErrorf(err, "", "effective_from-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "effective_until", &obj.EffectiveUntil)
	if err != nil {
		err = core.SDKErrorf(err, "", "effective_until-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "amounts", &obj.Amounts, UnmarshalAmount)
	if err != nil {
		err = core.SDKErrorf(err, "", "amounts-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ObjectMetadataSet : Model used to describe metadata object that can be set.
type ObjectMetadataSet struct {
	// Boolean value that describes whether the service is compatible with the Resource Controller.
	RcCompatible *bool `json:"rc_compatible,omitempty"`

	// Service-related metadata.
	Service *CfMetaData `json:"service,omitempty"`

	// Plan-related metadata.
	Plan *PlanMetaData `json:"plan,omitempty"`

	// Alias-related metadata.
	Alias *AliasMetaData `json:"alias,omitempty"`

	// Template-related metadata.
	Template *TemplateMetaData `json:"template,omitempty"`

	// Information related to the UI presentation associated with a catalog entry.
	UI *UIMetaData `json:"ui,omitempty"`

	// Compliance information for HIPAA and PCI.
	Compliance []string `json:"compliance,omitempty"`

	// Service Level Agreement related metadata.
	SLA *SLAMetaData `json:"sla,omitempty"`

	// Callback-related information associated with a catalog entry.
	Callbacks *Callbacks `json:"callbacks,omitempty"`

	// The original name of the object.
	OriginalName *string `json:"original_name,omitempty"`

	// Optional version of the object.
	Version *string `json:"version,omitempty"`

	// Additional information.
	Other map[string]interface{} `json:"other,omitempty"`

	// Pricing-related information.
	Pricing *PricingSet `json:"pricing,omitempty"`

	// Deployment-related metadata.
	Deployment *DeploymentBase `json:"deployment,omitempty"`
}

// UnmarshalObjectMetadataSet unmarshals an instance of ObjectMetadataSet from the specified map of raw messages.
func UnmarshalObjectMetadataSet(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ObjectMetadataSet)
	err = core.UnmarshalPrimitive(m, "rc_compatible", &obj.RcCompatible)
	if err != nil {
		err = core.SDKErrorf(err, "", "rc_compatible-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service", &obj.Service, UnmarshalCfMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "service-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "plan", &obj.Plan, UnmarshalPlanMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "alias", &obj.Alias, UnmarshalAliasMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "alias-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalTemplateMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "template-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ui", &obj.UI, UnmarshalUIMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "compliance", &obj.Compliance)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "sla", &obj.SLA, UnmarshalSLAMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "sla-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "callbacks", &obj.Callbacks, UnmarshalCallbacks)
	if err != nil {
		err = core.SDKErrorf(err, "", "callbacks-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "original_name", &obj.OriginalName)
	if err != nil {
		err = core.SDKErrorf(err, "", "original_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "other", &obj.Other)
	if err != nil {
		err = core.SDKErrorf(err, "", "other-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pricing", &obj.Pricing, UnmarshalPricingSet)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "deployment", &obj.Deployment, UnmarshalDeploymentBase)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Overview : Overview is nested in the top level. The key value pair is `[_language_]overview_ui`.
type Overview struct {
	// The translated display name.
	DisplayName *string `json:"display_name" validate:"required"`

	// The translated long description.
	LongDescription *string `json:"long_description" validate:"required"`

	// The translated description.
	Description *string `json:"description" validate:"required"`

	// The translated description that will be featured.
	FeaturedDescription *string `json:"featured_description,omitempty"`
}

// NewOverview : Instantiate Overview (Generic Model Constructor)
func (*GlobalCatalogV1) NewOverview(displayName string, longDescription string, description string) (_model *Overview, err error) {
	_model = &Overview{
		DisplayName:     core.StringPtr(displayName),
		LongDescription: core.StringPtr(longDescription),
		Description:     core.StringPtr(description),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalOverview unmarshals an instance of Overview from the specified map of raw messages.
func UnmarshalOverview(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Overview)
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		err = core.SDKErrorf(err, "", "long_description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "featured_description", &obj.FeaturedDescription)
	if err != nil {
		err = core.SDKErrorf(err, "", "featured_description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PlanMetaData : Plan-related metadata.
type PlanMetaData struct {
	// Boolean value that describes whether the service can be bound to an application.
	Bindable *bool `json:"bindable,omitempty"`

	// Boolean value that describes whether the service can be reserved.
	Reservable *bool `json:"reservable,omitempty"`

	// Boolean value that describes whether the service can be used internally.
	AllowInternalUsers *bool `json:"allow_internal_users,omitempty"`

	// Boolean value that describes whether the service can be provisioned asynchronously.
	AsyncProvisioningSupported *bool `json:"async_provisioning_supported,omitempty"`

	// Boolean value that describes whether the service can be unprovisioned asynchronously.
	AsyncUnprovisioningSupported *bool `json:"async_unprovisioning_supported,omitempty"`

	// Test check interval.
	TestCheckInterval *int64 `json:"test_check_interval,omitempty"`

	// Single scope instance.
	SingleScopeInstance *string `json:"single_scope_instance,omitempty"`

	// Boolean value that describes whether the service check is enabled.
	ServiceCheckEnabled *bool `json:"service_check_enabled,omitempty"`

	// If the field is imported from Cloud Foundry, the Cloud Foundry region's GUID. This is a required field. For example,
	// `us-south=123`.
	CfGUID map[string]string `json:"cf_guid,omitempty"`
}

// UnmarshalPlanMetaData unmarshals an instance of PlanMetaData from the specified map of raw messages.
func UnmarshalPlanMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PlanMetaData)
	err = core.UnmarshalPrimitive(m, "bindable", &obj.Bindable)
	if err != nil {
		err = core.SDKErrorf(err, "", "bindable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reservable", &obj.Reservable)
	if err != nil {
		err = core.SDKErrorf(err, "", "reservable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_internal_users", &obj.AllowInternalUsers)
	if err != nil {
		err = core.SDKErrorf(err, "", "allow_internal_users-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "async_provisioning_supported", &obj.AsyncProvisioningSupported)
	if err != nil {
		err = core.SDKErrorf(err, "", "async_provisioning_supported-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "async_unprovisioning_supported", &obj.AsyncUnprovisioningSupported)
	if err != nil {
		err = core.SDKErrorf(err, "", "async_unprovisioning_supported-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "test_check_interval", &obj.TestCheckInterval)
	if err != nil {
		err = core.SDKErrorf(err, "", "test_check_interval-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "single_scope_instance", &obj.SingleScopeInstance)
	if err != nil {
		err = core.SDKErrorf(err, "", "single_scope_instance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_check_enabled", &obj.ServiceCheckEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_check_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cf_guid", &obj.CfGUID)
	if err != nil {
		err = core.SDKErrorf(err, "", "cf_guid-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Price : Pricing-related information.
type Price struct {
	// Pricing tier.
	QuantityTier *int64 `json:"quantity_tier,omitempty"`

	// Price in the selected currency.
	Price *float64 `json:"price,omitempty"`
}

// UnmarshalPrice unmarshals an instance of Price from the specified map of raw messages.
func UnmarshalPrice(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Price)
	err = core.UnmarshalPrimitive(m, "quantity_tier", &obj.QuantityTier)
	if err != nil {
		err = core.SDKErrorf(err, "", "quantity_tier-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "price", &obj.Price)
	if err != nil {
		err = core.SDKErrorf(err, "", "price-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PricingGet : Pricing-related information.
type PricingGet struct {
	// The deployment object id this pricing is from. Only set if object kind is deployment.
	DeploymentID *string `json:"deployment_id,omitempty"`

	// The deployment location this pricing is from. Only set if object kind is deployment.
	DeploymentLocation *string `json:"deployment_location,omitempty"`

	// Is the location price not available. Only set in api /pricing/deployment and only set if true. This means for the
	// given deployment object there was no pricing set in pricing catalog.
	DeploymentLocationNoPriceAvailable *bool `json:"deployment_location_no_price_available,omitempty"`

	// Type of plan. Valid values are `free`, `trial`, `paygo`, `bluemix-subscription`, and `ibm-subscription`.
	Type *string `json:"type,omitempty"`

	// Defines where the pricing originates.
	Origin *string `json:"origin,omitempty"`

	// Plan-specific starting price information.
	StartingPrice *StartingPrice `json:"starting_price,omitempty"`

	// Plan-specific cost metric structure.
	Metrics []Metrics `json:"metrics,omitempty"`

	// List of regions where region pricing is available. Only set on global deployments if enabled by owner.
	DeploymentRegions []string `json:"deployment_regions,omitempty"`
}

// UnmarshalPricingGet unmarshals an instance of PricingGet from the specified map of raw messages.
func UnmarshalPricingGet(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PricingGet)
	err = core.UnmarshalPrimitive(m, "deployment_id", &obj.DeploymentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_location", &obj.DeploymentLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_location_no_price_available", &obj.DeploymentLocationNoPriceAvailable)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_location_no_price_available-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "origin", &obj.Origin)
	if err != nil {
		err = core.SDKErrorf(err, "", "origin-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "starting_price", &obj.StartingPrice, UnmarshalStartingPrice)
	if err != nil {
		err = core.SDKErrorf(err, "", "starting_price-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metrics", &obj.Metrics, UnmarshalMetrics)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_regions", &obj.DeploymentRegions)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_regions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PricingSearchResult : A paginated result containing pricing entries.
type PricingSearchResult struct {
	// The offset (origin 0) of the first resource in this page of search results.
	Offset *int64 `json:"offset,omitempty"`

	// The maximum number of resources returned in each page of search results.
	Limit *int64 `json:"limit,omitempty"`

	// The overall total number of resources in the search result set.
	Count *int64 `json:"count,omitempty"`

	// The number of resources returned in this page of search results.
	ResourceCount *int64 `json:"resource_count,omitempty"`

	// A URL for retrieving the first page of search results.
	First *string `json:"first,omitempty"`

	// A URL for retrieving the last page of search results.
	Last *string `json:"last,omitempty"`

	// A URL for retrieving the previous page of search results.
	Prev *string `json:"prev,omitempty"`

	// A URL for retrieving the next page of search results.
	Next *string `json:"next,omitempty"`

	// The resources (prices) contained in this page of search results.
	Resources []PricingGet `json:"resources,omitempty"`
}

// UnmarshalPricingSearchResult unmarshals an instance of PricingSearchResult from the specified map of raw messages.
func UnmarshalPricingSearchResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PricingSearchResult)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_count", &obj.ResourceCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last", &obj.Last)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "prev", &obj.Prev)
	if err != nil {
		err = core.SDKErrorf(err, "", "prev-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalPricingGet)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PricingSet : Pricing-related information.
type PricingSet struct {
	// Type of plan. Valid values are `free`, `trial`, `paygo`, `bluemix-subscription`, and `ibm-subscription`.
	Type *string `json:"type,omitempty"`

	// Defines where the pricing originates.
	Origin *string `json:"origin,omitempty"`

	// Plan-specific starting price information.
	StartingPrice *StartingPrice `json:"starting_price,omitempty"`
}

// UnmarshalPricingSet unmarshals an instance of PricingSet from the specified map of raw messages.
func UnmarshalPricingSet(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PricingSet)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "origin", &obj.Origin)
	if err != nil {
		err = core.SDKErrorf(err, "", "origin-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "starting_price", &obj.StartingPrice, UnmarshalStartingPrice)
	if err != nil {
		err = core.SDKErrorf(err, "", "starting_price-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Provider : Information related to the provider associated with a catalog entry.
type Provider struct {
	// Provider's email address for this catalog entry.
	Email *string `json:"email" validate:"required"`

	// Provider's name, for example, IBM.
	Name *string `json:"name" validate:"required"`

	// Provider's contact name.
	Contact *string `json:"contact,omitempty"`

	// Provider's support email.
	SupportEmail *string `json:"support_email,omitempty"`

	// Provider's contact phone.
	Phone *string `json:"phone,omitempty"`
}

// NewProvider : Instantiate Provider (Generic Model Constructor)
func (*GlobalCatalogV1) NewProvider(email string, name string) (_model *Provider, err error) {
	_model = &Provider{
		Email: core.StringPtr(email),
		Name:  core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalProvider unmarshals an instance of Provider from the specified map of raw messages.
func UnmarshalProvider(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Provider)
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "contact", &obj.Contact)
	if err != nil {
		err = core.SDKErrorf(err, "", "contact-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "support_email", &obj.SupportEmail)
	if err != nil {
		err = core.SDKErrorf(err, "", "support_email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "phone", &obj.Phone)
	if err != nil {
		err = core.SDKErrorf(err, "", "phone-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RestoreCatalogEntryOptions : The RestoreCatalogEntry options.
type RestoreCatalogEntryOptions struct {
	// The catalog entry's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewRestoreCatalogEntryOptions : Instantiate RestoreCatalogEntryOptions
func (*GlobalCatalogV1) NewRestoreCatalogEntryOptions(id string) *RestoreCatalogEntryOptions {
	return &RestoreCatalogEntryOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *RestoreCatalogEntryOptions) SetID(id string) *RestoreCatalogEntryOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *RestoreCatalogEntryOptions) SetAccount(account string) *RestoreCatalogEntryOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *RestoreCatalogEntryOptions) SetHeaders(param map[string]string) *RestoreCatalogEntryOptions {
	options.Headers = param
	return options
}

// SLAMetaData : Service Level Agreement related metadata.
type SLAMetaData struct {
	// Required Service License Agreement Terms of Use.
	Terms *string `json:"terms,omitempty"`

	// Required deployment type. Valid values are dedicated, local, or public. It can be Single or Multi tennancy, more
	// specifically on a Server, VM, Physical, or Pod.
	Tenancy *string `json:"tenancy,omitempty"`

	// Provisioning reliability, for example, 99.95.
	Provisioning *float64 `json:"provisioning,omitempty"`

	// Uptime reliability of the service, for example, 99.95.
	Responsiveness *float64 `json:"responsiveness,omitempty"`

	// SLA Disaster Recovery-related metadata.
	Dr *DrMetaData `json:"dr,omitempty"`
}

// UnmarshalSLAMetaData unmarshals an instance of SLAMetaData from the specified map of raw messages.
func UnmarshalSLAMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SLAMetaData)
	err = core.UnmarshalPrimitive(m, "terms", &obj.Terms)
	if err != nil {
		err = core.SDKErrorf(err, "", "terms-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tenancy", &obj.Tenancy)
	if err != nil {
		err = core.SDKErrorf(err, "", "tenancy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provisioning", &obj.Provisioning)
	if err != nil {
		err = core.SDKErrorf(err, "", "provisioning-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "responsiveness", &obj.Responsiveness)
	if err != nil {
		err = core.SDKErrorf(err, "", "responsiveness-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "dr", &obj.Dr, UnmarshalDrMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "dr-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SourceMetaData : Location of your applications source files.
type SourceMetaData struct {
	// Path to your application.
	Path *string `json:"path,omitempty"`

	// Type of source, for example, git.
	Type *string `json:"type,omitempty"`

	// URL to source.
	URL *string `json:"url,omitempty"`
}

// UnmarshalSourceMetaData unmarshals an instance of SourceMetaData from the specified map of raw messages.
func UnmarshalSourceMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SourceMetaData)
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		err = core.SDKErrorf(err, "", "path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StartingPrice : Plan-specific starting price information.
type StartingPrice struct {
	// ID of the plan the starting price is calculated.
	PlanID *string `json:"plan_id,omitempty"`

	// ID of the deployment the starting price is calculated.
	DeploymentID *string `json:"deployment_id,omitempty"`

	// Pricing unit.
	Unit *string `json:"unit,omitempty"`

	// The pricing per metric by country and currency.
	Amount []Amount `json:"amount,omitempty"`
}

// UnmarshalStartingPrice unmarshals an instance of StartingPrice from the specified map of raw messages.
func UnmarshalStartingPrice(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StartingPrice)
	err = core.UnmarshalPrimitive(m, "plan_id", &obj.PlanID)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_id", &obj.DeploymentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "unit", &obj.Unit)
	if err != nil {
		err = core.SDKErrorf(err, "", "unit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "amount", &obj.Amount, UnmarshalAmount)
	if err != nil {
		err = core.SDKErrorf(err, "", "amount-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Strings : Information related to a translated text message.
type Strings struct {
	// Presentation information related to list delimiters.
	Bullets []Bullets `json:"bullets,omitempty"`

	// Media-related metadata.
	Media []UIMetaMedia `json:"media,omitempty"`

	// Warning that a message is not creatable.
	NotCreatableMsg *string `json:"not_creatable_msg,omitempty"`

	// Warning that a robot message is not creatable.
	NotCreatableRobotMsg *string `json:"not_creatable__robot_msg,omitempty"`

	// Warning for deprecation.
	DeprecationWarning *string `json:"deprecation_warning,omitempty"`

	// Popup warning message.
	PopupWarningMessage *string `json:"popup_warning_message,omitempty"`

	// Instructions for UI strings.
	Instruction *string `json:"instruction,omitempty"`
}

// UnmarshalStrings unmarshals an instance of Strings from the specified map of raw messages.
func UnmarshalStrings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Strings)
	err = core.UnmarshalModel(m, "bullets", &obj.Bullets, UnmarshalBullets)
	if err != nil {
		err = core.SDKErrorf(err, "", "bullets-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "media", &obj.Media, UnmarshalUIMetaMedia)
	if err != nil {
		err = core.SDKErrorf(err, "", "media-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "not_creatable_msg", &obj.NotCreatableMsg)
	if err != nil {
		err = core.SDKErrorf(err, "", "not_creatable_msg-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "not_creatable__robot_msg", &obj.NotCreatableRobotMsg)
	if err != nil {
		err = core.SDKErrorf(err, "", "not_creatable__robot_msg-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deprecation_warning", &obj.DeprecationWarning)
	if err != nil {
		err = core.SDKErrorf(err, "", "deprecation_warning-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "popup_warning_message", &obj.PopupWarningMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "popup_warning_message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instruction", &obj.Instruction)
	if err != nil {
		err = core.SDKErrorf(err, "", "instruction-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateMetaData : Template-related metadata.
type TemplateMetaData struct {
	// List of required offering or plan IDs.
	Services []string `json:"services,omitempty"`

	// Cloud Foundry instance memory value.
	DefaultMemory *int64 `json:"default_memory,omitempty"`

	// Start Command.
	StartCmd *string `json:"start_cmd,omitempty"`

	// Location of your applications source files.
	Source *SourceMetaData `json:"source,omitempty"`

	// ID of the runtime.
	RuntimeCatalogID *string `json:"runtime_catalog_id,omitempty"`

	// ID of the Cloud Foundry runtime.
	CfRuntimeID *string `json:"cf_runtime_id,omitempty"`

	// ID of the boilerplate or template.
	TemplateID *string `json:"template_id,omitempty"`

	// File path to the executable file for the template.
	ExecutableFile *string `json:"executable_file,omitempty"`

	// ID of the buildpack used by the template.
	Buildpack *string `json:"buildpack,omitempty"`

	// Environment variables (key/value pairs) for the template.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
}

// UnmarshalTemplateMetaData unmarshals an instance of TemplateMetaData from the specified map of raw messages.
func UnmarshalTemplateMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateMetaData)
	err = core.UnmarshalPrimitive(m, "services", &obj.Services)
	if err != nil {
		err = core.SDKErrorf(err, "", "services-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "default_memory", &obj.DefaultMemory)
	if err != nil {
		err = core.SDKErrorf(err, "", "default_memory-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start_cmd", &obj.StartCmd)
	if err != nil {
		err = core.SDKErrorf(err, "", "start_cmd-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalSourceMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "source-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "runtime_catalog_id", &obj.RuntimeCatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "runtime_catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cf_runtime_id", &obj.CfRuntimeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "cf_runtime_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "executable_file", &obj.ExecutableFile)
	if err != nil {
		err = core.SDKErrorf(err, "", "executable_file-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "buildpack", &obj.Buildpack)
	if err != nil {
		err = core.SDKErrorf(err, "", "buildpack-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_variables", &obj.EnvironmentVariables)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_variables-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UIMediaSourceMetaData : Location of your applications media source files.
type UIMediaSourceMetaData struct {
	// Type of source, for example, git.
	Type *string `json:"type,omitempty"`

	// URL to source.
	URL *string `json:"url,omitempty"`
}

// UnmarshalUIMediaSourceMetaData unmarshals an instance of UIMediaSourceMetaData from the specified map of raw messages.
func UnmarshalUIMediaSourceMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UIMediaSourceMetaData)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UIMetaData : Information related to the UI presentation associated with a catalog entry.
type UIMetaData struct {
	// Language specific translation of translation properties, like label and description.
	Strings map[string]Strings `json:"strings,omitempty"`

	// UI based URLs.
	Urls *Urls `json:"urls,omitempty"`

	// Describes how the embeddable dashboard is rendered.
	EmbeddableDashboard *string `json:"embeddable_dashboard,omitempty"`

	// Describes whether the embeddable dashboard is rendered at the full width.
	EmbeddableDashboardFullWidth *bool `json:"embeddable_dashboard_full_width,omitempty"`

	// Defines the order of information presented.
	NavigationOrder []string `json:"navigation_order,omitempty"`

	// Describes whether this entry is able to be created from the UI element or CLI.
	NotCreatable *bool `json:"not_creatable,omitempty"`

	// ID of the primary offering for a group.
	PrimaryOfferingID *string `json:"primary_offering_id,omitempty"`

	// Alert to ACE to allow instance UI to be accessible while the provisioning state of instance is in progress.
	AccessibleDuringProvision *bool `json:"accessible_during_provision,omitempty"`

	// Specifies a side by side ordering weight to the UI.
	SideBySideIndex *int64 `json:"side_by_side_index,omitempty"`

	// Date and time the service will no longer be available.
	EndOfServiceTime *strfmt.DateTime `json:"end_of_service_time,omitempty"`

	// Denotes visibility.
	Hidden *bool `json:"hidden,omitempty"`

	// Denotes lite metering visibility.
	HideLiteMetering *bool `json:"hide_lite_metering,omitempty"`

	// Denotes whether an upgrade should occur.
	NoUpgradeNextStep *bool `json:"no_upgrade_next_step,omitempty"`
}

// UnmarshalUIMetaData unmarshals an instance of UIMetaData from the specified map of raw messages.
func UnmarshalUIMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UIMetaData)
	err = core.UnmarshalModel(m, "strings", &obj.Strings, UnmarshalStrings)
	if err != nil {
		err = core.SDKErrorf(err, "", "strings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "urls", &obj.Urls, UnmarshalUrls)
	if err != nil {
		err = core.SDKErrorf(err, "", "urls-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "embeddable_dashboard", &obj.EmbeddableDashboard)
	if err != nil {
		err = core.SDKErrorf(err, "", "embeddable_dashboard-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "embeddable_dashboard_full_width", &obj.EmbeddableDashboardFullWidth)
	if err != nil {
		err = core.SDKErrorf(err, "", "embeddable_dashboard_full_width-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "navigation_order", &obj.NavigationOrder)
	if err != nil {
		err = core.SDKErrorf(err, "", "navigation_order-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "not_creatable", &obj.NotCreatable)
	if err != nil {
		err = core.SDKErrorf(err, "", "not_creatable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "primary_offering_id", &obj.PrimaryOfferingID)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_offering_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "accessible_during_provision", &obj.AccessibleDuringProvision)
	if err != nil {
		err = core.SDKErrorf(err, "", "accessible_during_provision-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "side_by_side_index", &obj.SideBySideIndex)
	if err != nil {
		err = core.SDKErrorf(err, "", "side_by_side_index-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "end_of_service_time", &obj.EndOfServiceTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "end_of_service_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		err = core.SDKErrorf(err, "", "hidden-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hide_lite_metering", &obj.HideLiteMetering)
	if err != nil {
		err = core.SDKErrorf(err, "", "hide_lite_metering-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "no_upgrade_next_step", &obj.NoUpgradeNextStep)
	if err != nil {
		err = core.SDKErrorf(err, "", "no_upgrade_next_step-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UIMetaMedia : Media-related metadata.
type UIMetaMedia struct {
	// Caption for an image.
	Caption *string `json:"caption,omitempty"`

	// URL for thumbnail image.
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`

	// Type of media.
	Type *string `json:"type,omitempty"`

	// URL for media.
	URL *string `json:"URL,omitempty"`

	// UI media source data for for UI media data.
	Source []UIMediaSourceMetaData `json:"source,omitempty"`
}

// UnmarshalUIMetaMedia unmarshals an instance of UIMetaMedia from the specified map of raw messages.
func UnmarshalUIMetaMedia(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UIMetaMedia)
	err = core.UnmarshalPrimitive(m, "caption", &obj.Caption)
	if err != nil {
		err = core.SDKErrorf(err, "", "caption-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "thumbnail_url", &obj.ThumbnailURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "thumbnail_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "URL", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "URL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalUIMediaSourceMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "source-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Urls : UI based URLs.
type Urls struct {
	// URL for documentation.
	DocURL *string `json:"doc_url,omitempty"`

	// URL for usage instructions.
	InstructionsURL *string `json:"instructions_url,omitempty"`

	// API URL.
	APIURL *string `json:"api_url,omitempty"`

	// URL Creation UI / API.
	CreateURL *string `json:"create_url,omitempty"`

	// URL to downlaod an SDK.
	SdkDownloadURL *string `json:"sdk_download_url,omitempty"`

	// URL to the terms of use for your service.
	TermsURL *string `json:"terms_url,omitempty"`

	// URL to the custom create page for your serivce.
	CustomCreatePageURL *string `json:"custom_create_page_url,omitempty"`

	// URL to the catalog details page for your serivce.
	CatalogDetailsURL *string `json:"catalog_details_url,omitempty"`

	// URL for deprecation documentation.
	DeprecationDocURL *string `json:"deprecation_doc_url,omitempty"`

	// URL for dashboard.
	DashboardURL *string `json:"dashboard_url,omitempty"`

	// URL for registration.
	RegistrationURL *string `json:"registration_url,omitempty"`

	// URL for API documentation.
	Apidocsurl *string `json:"apidocsurl,omitempty"`
}

// UnmarshalUrls unmarshals an instance of Urls from the specified map of raw messages.
func UnmarshalUrls(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Urls)
	err = core.UnmarshalPrimitive(m, "doc_url", &obj.DocURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "doc_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instructions_url", &obj.InstructionsURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "instructions_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_url", &obj.APIURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "create_url", &obj.CreateURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "create_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sdk_download_url", &obj.SdkDownloadURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "sdk_download_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "terms_url", &obj.TermsURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "terms_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_create_page_url", &obj.CustomCreatePageURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "custom_create_page_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_details_url", &obj.CatalogDetailsURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_details_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deprecation_doc_url", &obj.DeprecationDocURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "deprecation_doc_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dashboard_url", &obj.DashboardURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "dashboard_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "registration_url", &obj.RegistrationURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "registration_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "apidocsurl", &obj.Apidocsurl)
	if err != nil {
		err = core.SDKErrorf(err, "", "apidocsurl-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateCatalogEntryOptions : The UpdateCatalogEntry options.
type UpdateCatalogEntryOptions struct {
	// The object's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// Programmatic name for this catalog entry, which must be formatted like a CRN segment. See the display name in
	// OverviewUI for a user-readable name.
	Name *string `json:"name" validate:"required"`

	// The type of catalog entry, **service**, **template**, **dashboard**, which determines the type and shape of the
	// object.
	Kind *string `json:"kind" validate:"required"`

	// Overview is nested in the top level. The key value pair is `[_language_]overview_ui`.
	OverviewUI map[string]Overview `json:"overview_ui" validate:"required"`

	// Image annotation for this catalog entry. The image is a URL.
	Images *Image `json:"images" validate:"required"`

	// Boolean value that determines the global visibility for the catalog entry, and its children. If it is not enabled,
	// all plans are disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// A list of tags. For example, IBM, 3rd Party, Beta, GA, and Single Tenant.
	Tags []string `json:"tags" validate:"required"`

	// Information related to the provider associated with a catalog entry.
	Provider *Provider `json:"provider" validate:"required"`

	// The ID of the parent catalog entry if it exists.
	ParentID *string `json:"parent_id,omitempty"`

	// Boolean value that determines whether the catalog entry is a group.
	Group *bool `json:"group,omitempty"`

	// Boolean value that describes whether the service is active.
	Active *bool `json:"active,omitempty"`

	// Url of the object.
	URL *string `json:"url,omitempty"`

	// Model used to describe metadata object that can be set.
	Metadata *ObjectMetadataSet `json:"metadata,omitempty"`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Reparenting object. In the body set the parent_id to a different parent. Or remove the parent_id field to reparent
	// to the root of the catalog. If this is not set to 'true' then changing the parent_id in the body of the request will
	// not be permitted. If this is 'true' and no change to parent_id then this is also error. This is to prevent
	// accidental changing of parent.
	Move *string `json:"move,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the UpdateCatalogEntryOptions.Kind property.
// The type of catalog entry, **service**, **template**, **dashboard**, which determines the type and shape of the
// object.
const (
	UpdateCatalogEntryOptionsKindDashboardConst = "dashboard"
	UpdateCatalogEntryOptionsKindServiceConst   = "service"
	UpdateCatalogEntryOptionsKindTemplateConst  = "template"
)

// NewUpdateCatalogEntryOptions : Instantiate UpdateCatalogEntryOptions
func (*GlobalCatalogV1) NewUpdateCatalogEntryOptions(id string, name string, kind string, overviewUI map[string]Overview, images *Image, disabled bool, tags []string, provider *Provider) *UpdateCatalogEntryOptions {
	return &UpdateCatalogEntryOptions{
		ID:         core.StringPtr(id),
		Name:       core.StringPtr(name),
		Kind:       core.StringPtr(kind),
		OverviewUI: overviewUI,
		Images:     images,
		Disabled:   core.BoolPtr(disabled),
		Tags:       tags,
		Provider:   provider,
	}
}

// SetID : Allow user to set ID
func (_options *UpdateCatalogEntryOptions) SetID(id string) *UpdateCatalogEntryOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateCatalogEntryOptions) SetName(name string) *UpdateCatalogEntryOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *UpdateCatalogEntryOptions) SetKind(kind string) *UpdateCatalogEntryOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetOverviewUI : Allow user to set OverviewUI
func (_options *UpdateCatalogEntryOptions) SetOverviewUI(overviewUI map[string]Overview) *UpdateCatalogEntryOptions {
	_options.OverviewUI = overviewUI
	return _options
}

// SetImages : Allow user to set Images
func (_options *UpdateCatalogEntryOptions) SetImages(images *Image) *UpdateCatalogEntryOptions {
	_options.Images = images
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *UpdateCatalogEntryOptions) SetDisabled(disabled bool) *UpdateCatalogEntryOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateCatalogEntryOptions) SetTags(tags []string) *UpdateCatalogEntryOptions {
	_options.Tags = tags
	return _options
}

// SetProvider : Allow user to set Provider
func (_options *UpdateCatalogEntryOptions) SetProvider(provider *Provider) *UpdateCatalogEntryOptions {
	_options.Provider = provider
	return _options
}

// SetParentID : Allow user to set ParentID
func (_options *UpdateCatalogEntryOptions) SetParentID(parentID string) *UpdateCatalogEntryOptions {
	_options.ParentID = core.StringPtr(parentID)
	return _options
}

// SetGroup : Allow user to set Group
func (_options *UpdateCatalogEntryOptions) SetGroup(group bool) *UpdateCatalogEntryOptions {
	_options.Group = core.BoolPtr(group)
	return _options
}

// SetActive : Allow user to set Active
func (_options *UpdateCatalogEntryOptions) SetActive(active bool) *UpdateCatalogEntryOptions {
	_options.Active = core.BoolPtr(active)
	return _options
}

// SetURL : Allow user to set URL
func (_options *UpdateCatalogEntryOptions) SetURL(url string) *UpdateCatalogEntryOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *UpdateCatalogEntryOptions) SetMetadata(metadata *ObjectMetadataSet) *UpdateCatalogEntryOptions {
	_options.Metadata = metadata
	return _options
}

// SetAccount : Allow user to set Account
func (_options *UpdateCatalogEntryOptions) SetAccount(account string) *UpdateCatalogEntryOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetMove : Allow user to set Move
func (_options *UpdateCatalogEntryOptions) SetMove(move string) *UpdateCatalogEntryOptions {
	_options.Move = core.StringPtr(move)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCatalogEntryOptions) SetHeaders(param map[string]string) *UpdateCatalogEntryOptions {
	options.Headers = param
	return options
}

// UpdateVisibilityOptions : The UpdateVisibility options.
type UpdateVisibilityOptions struct {
	// The object's unique ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows the visibility to be extenable.
	Extendable *bool `json:"extendable,omitempty"`

	// Visibility details related to a catalog entry.
	Include *VisibilityDetail `json:"include,omitempty"`

	// Visibility details related to a catalog entry.
	Exclude *VisibilityDetail `json:"exclude,omitempty"`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateVisibilityOptions : Instantiate UpdateVisibilityOptions
func (*GlobalCatalogV1) NewUpdateVisibilityOptions(id string) *UpdateVisibilityOptions {
	return &UpdateVisibilityOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateVisibilityOptions) SetID(id string) *UpdateVisibilityOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetExtendable : Allow user to set Extendable
func (_options *UpdateVisibilityOptions) SetExtendable(extendable bool) *UpdateVisibilityOptions {
	_options.Extendable = core.BoolPtr(extendable)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *UpdateVisibilityOptions) SetInclude(include *VisibilityDetail) *UpdateVisibilityOptions {
	_options.Include = include
	return _options
}

// SetExclude : Allow user to set Exclude
func (_options *UpdateVisibilityOptions) SetExclude(exclude *VisibilityDetail) *UpdateVisibilityOptions {
	_options.Exclude = exclude
	return _options
}

// SetAccount : Allow user to set Account
func (_options *UpdateVisibilityOptions) SetAccount(account string) *UpdateVisibilityOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateVisibilityOptions) SetHeaders(param map[string]string) *UpdateVisibilityOptions {
	options.Headers = param
	return options
}

// UploadArtifactOptions : The UploadArtifact options.
type UploadArtifactOptions struct {
	// The object's unique ID.
	ObjectID *string `json:"object_id" validate:"required,ne="`

	// The artifact's ID.
	ArtifactID *string `json:"artifact_id" validate:"required,ne="`

	Artifact io.ReadCloser `json:"artifact,omitempty"`

	// The type of the input.
	ContentType *string `json:"Content-Type,omitempty"`

	// This changes the scope of the request regardless of the authorization header. Example scopes are `account` and
	// `global`. `account=global` is reqired if operating with a service ID that has a global admin policy, for example
	// `GET /?account=global`.
	Account *string `json:"account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUploadArtifactOptions : Instantiate UploadArtifactOptions
func (*GlobalCatalogV1) NewUploadArtifactOptions(objectID string, artifactID string) *UploadArtifactOptions {
	return &UploadArtifactOptions{
		ObjectID:   core.StringPtr(objectID),
		ArtifactID: core.StringPtr(artifactID),
	}
}

// SetObjectID : Allow user to set ObjectID
func (_options *UploadArtifactOptions) SetObjectID(objectID string) *UploadArtifactOptions {
	_options.ObjectID = core.StringPtr(objectID)
	return _options
}

// SetArtifactID : Allow user to set ArtifactID
func (_options *UploadArtifactOptions) SetArtifactID(artifactID string) *UploadArtifactOptions {
	_options.ArtifactID = core.StringPtr(artifactID)
	return _options
}

// SetArtifact : Allow user to set Artifact
func (_options *UploadArtifactOptions) SetArtifact(artifact io.ReadCloser) *UploadArtifactOptions {
	_options.Artifact = artifact
	return _options
}

// SetContentType : Allow user to set ContentType
func (_options *UploadArtifactOptions) SetContentType(contentType string) *UploadArtifactOptions {
	_options.ContentType = core.StringPtr(contentType)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *UploadArtifactOptions) SetAccount(account string) *UploadArtifactOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UploadArtifactOptions) SetHeaders(param map[string]string) *UploadArtifactOptions {
	options.Headers = param
	return options
}

// Visibility : Information related to the visibility of a catalog entry.
type Visibility struct {
	// This controls the overall visibility. It is an enum of *public*, *ibm_only*, and *private*. public means it is
	// visible to all. ibm_only means it is visible to all IBM unless their account is explicitly excluded. private means
	// it is visible only to the included accounts.
	Restrictions *string `json:"restrictions,omitempty"`

	// IAM Scope-related information associated with a catalog entry.
	Owner *string `json:"owner,omitempty"`

	// Allows the visibility to be extenable.
	Extendable *bool `json:"extendable,omitempty"`

	// Visibility details related to a catalog entry.
	Include *VisibilityDetail `json:"include,omitempty"`

	// Visibility details related to a catalog entry.
	Exclude *VisibilityDetail `json:"exclude,omitempty"`

	// Determines whether the owning account has full control over the visibility of the entry such as adding non-IBM
	// accounts to the whitelist and making entries `private`, `ibm_only` or `public`.
	Approved *bool `json:"approved,omitempty"`
}

// UnmarshalVisibility unmarshals an instance of Visibility from the specified map of raw messages.
func UnmarshalVisibility(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Visibility)
	err = core.UnmarshalPrimitive(m, "restrictions", &obj.Restrictions)
	if err != nil {
		err = core.SDKErrorf(err, "", "restrictions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "owner", &obj.Owner)
	if err != nil {
		err = core.SDKErrorf(err, "", "owner-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "extendable", &obj.Extendable)
	if err != nil {
		err = core.SDKErrorf(err, "", "extendable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "include", &obj.Include, UnmarshalVisibilityDetail)
	if err != nil {
		err = core.SDKErrorf(err, "", "include-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "exclude", &obj.Exclude, UnmarshalVisibilityDetail)
	if err != nil {
		err = core.SDKErrorf(err, "", "exclude-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approved", &obj.Approved)
	if err != nil {
		err = core.SDKErrorf(err, "", "approved-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VisibilityDetail : Visibility details related to a catalog entry.
type VisibilityDetail struct {
	// Information related to the accounts for which a catalog entry is visible.
	Accounts *VisibilityDetailAccounts `json:"accounts" validate:"required"`
}

// NewVisibilityDetail : Instantiate VisibilityDetail (Generic Model Constructor)
func (*GlobalCatalogV1) NewVisibilityDetail(accounts *VisibilityDetailAccounts) (_model *VisibilityDetail, err error) {
	_model = &VisibilityDetail{
		Accounts: accounts,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalVisibilityDetail unmarshals an instance of VisibilityDetail from the specified map of raw messages.
func UnmarshalVisibilityDetail(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VisibilityDetail)
	err = core.UnmarshalModel(m, "accounts", &obj.Accounts, UnmarshalVisibilityDetailAccounts)
	if err != nil {
		err = core.SDKErrorf(err, "", "accounts-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VisibilityDetailAccounts : Information related to the accounts for which a catalog entry is visible.
type VisibilityDetailAccounts struct {
	// (_accountid_) is the GUID of the account and the value is the scope of who set it. For setting visibility use "" as
	// the value. It is replaced with the owner scope when saved.
	Accountid *string `json:"_accountid_,omitempty"`
}

// UnmarshalVisibilityDetailAccounts unmarshals an instance of VisibilityDetailAccounts from the specified map of raw messages.
func UnmarshalVisibilityDetailAccounts(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VisibilityDetailAccounts)
	err = core.UnmarshalPrimitive(m, "_accountid_", &obj.Accountid)
	if err != nil {
		err = core.SDKErrorf(err, "", "_accountid_-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
