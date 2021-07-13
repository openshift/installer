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
 * IBM OpenAPI SDK Code Generator Version: 3.31.0-902c9336-20210504-161156
 */

// Package catalogmanagementv1 : Operations and models for the CatalogManagementV1 service
package catalogmanagementv1

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

// CatalogManagementV1 : This is the API to use for managing private catalogs for IBM Cloud. Private catalogs provide a
// way to centrally manage access to products in the IBM Cloud catalog and your own catalogs.
//
// Version: 1.0
type CatalogManagementV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://cm.globalcatalog.cloud.ibm.com/api/v1-beta"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "catalog_management"

// CatalogManagementV1Options : Service options
type CatalogManagementV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewCatalogManagementV1UsingExternalConfig : constructs an instance of CatalogManagementV1 with passed in options and external configuration.
func NewCatalogManagementV1UsingExternalConfig(options *CatalogManagementV1Options) (catalogManagement *CatalogManagementV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	catalogManagement, err = NewCatalogManagementV1(options)
	if err != nil {
		return
	}

	err = catalogManagement.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = catalogManagement.Service.SetServiceURL(options.URL)
	}
	return
}

// NewCatalogManagementV1 : constructs an instance of CatalogManagementV1 with passed in options.
func NewCatalogManagementV1(options *CatalogManagementV1Options) (service *CatalogManagementV1, err error) {
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

	service = &CatalogManagementV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "catalogManagement" suitable for processing requests.
func (catalogManagement *CatalogManagementV1) Clone() *CatalogManagementV1 {
	if core.IsNil(catalogManagement) {
		return nil
	}
	clone := *catalogManagement
	clone.Service = catalogManagement.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (catalogManagement *CatalogManagementV1) SetServiceURL(url string) error {
	return catalogManagement.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (catalogManagement *CatalogManagementV1) GetServiceURL() string {
	return catalogManagement.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (catalogManagement *CatalogManagementV1) SetDefaultHeaders(headers http.Header) {
	catalogManagement.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (catalogManagement *CatalogManagementV1) SetEnableGzipCompression(enableGzip bool) {
	catalogManagement.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (catalogManagement *CatalogManagementV1) GetEnableGzipCompression() bool {
	return catalogManagement.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (catalogManagement *CatalogManagementV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	catalogManagement.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (catalogManagement *CatalogManagementV1) DisableRetries() {
	catalogManagement.Service.DisableRetries()
}

// GetCatalogAccount : Get catalog account settings
// Get the account level settings for the account for private catalog.
func (catalogManagement *CatalogManagementV1) GetCatalogAccount(getCatalogAccountOptions *GetCatalogAccountOptions) (result *Account, response *core.DetailedResponse, err error) {
	return catalogManagement.GetCatalogAccountWithContext(context.Background(), getCatalogAccountOptions)
}

// GetCatalogAccountWithContext is an alternate form of the GetCatalogAccount method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetCatalogAccountWithContext(ctx context.Context, getCatalogAccountOptions *GetCatalogAccountOptions) (result *Account, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getCatalogAccountOptions, "getCatalogAccountOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogaccount`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCatalogAccountOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetCatalogAccount")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccount)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateCatalogAccount : Update account settings
// Update the account level settings for the account for private catalog.
func (catalogManagement *CatalogManagementV1) UpdateCatalogAccount(updateCatalogAccountOptions *UpdateCatalogAccountOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.UpdateCatalogAccountWithContext(context.Background(), updateCatalogAccountOptions)
}

// UpdateCatalogAccountWithContext is an alternate form of the UpdateCatalogAccount method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) UpdateCatalogAccountWithContext(ctx context.Context, updateCatalogAccountOptions *UpdateCatalogAccountOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(updateCatalogAccountOptions, "updateCatalogAccountOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogaccount`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateCatalogAccountOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "UpdateCatalogAccount")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateCatalogAccountOptions.ID != nil {
		body["id"] = updateCatalogAccountOptions.ID
	}
	if updateCatalogAccountOptions.HideIBMCloudCatalog != nil {
		body["hide_IBM_cloud_catalog"] = updateCatalogAccountOptions.HideIBMCloudCatalog
	}
	if updateCatalogAccountOptions.AccountFilters != nil {
		body["account_filters"] = updateCatalogAccountOptions.AccountFilters
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetCatalogAccountAudit : Get catalog account audit log
// Get the audit log associated with a catalog account.
func (catalogManagement *CatalogManagementV1) GetCatalogAccountAudit(getCatalogAccountAuditOptions *GetCatalogAccountAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	return catalogManagement.GetCatalogAccountAuditWithContext(context.Background(), getCatalogAccountAuditOptions)
}

// GetCatalogAccountAuditWithContext is an alternate form of the GetCatalogAccountAudit method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetCatalogAccountAuditWithContext(ctx context.Context, getCatalogAccountAuditOptions *GetCatalogAccountAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getCatalogAccountAuditOptions, "getCatalogAccountAuditOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogaccount/audit`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCatalogAccountAuditOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetCatalogAccountAudit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCatalogAccountFilters : Get catalog account filters
// Get the accumulated filters of the account and of the catalogs you have access to.
func (catalogManagement *CatalogManagementV1) GetCatalogAccountFilters(getCatalogAccountFiltersOptions *GetCatalogAccountFiltersOptions) (result *AccumulatedFilters, response *core.DetailedResponse, err error) {
	return catalogManagement.GetCatalogAccountFiltersWithContext(context.Background(), getCatalogAccountFiltersOptions)
}

// GetCatalogAccountFiltersWithContext is an alternate form of the GetCatalogAccountFilters method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetCatalogAccountFiltersWithContext(ctx context.Context, getCatalogAccountFiltersOptions *GetCatalogAccountFiltersOptions) (result *AccumulatedFilters, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getCatalogAccountFiltersOptions, "getCatalogAccountFiltersOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogaccount/filters`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCatalogAccountFiltersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetCatalogAccountFilters")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getCatalogAccountFiltersOptions.Catalog != nil {
		builder.AddQuery("catalog", fmt.Sprint(*getCatalogAccountFiltersOptions.Catalog))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccumulatedFilters)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListCatalogs : Get list of catalogs
// Retrieves the available catalogs for a given account. This can be used by an unauthenticated user to retrieve the
// public catalog.
func (catalogManagement *CatalogManagementV1) ListCatalogs(listCatalogsOptions *ListCatalogsOptions) (result *CatalogSearchResult, response *core.DetailedResponse, err error) {
	return catalogManagement.ListCatalogsWithContext(context.Background(), listCatalogsOptions)
}

// ListCatalogsWithContext is an alternate form of the ListCatalogs method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListCatalogsWithContext(ctx context.Context, listCatalogsOptions *ListCatalogsOptions) (result *CatalogSearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCatalogsOptions, "listCatalogsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listCatalogsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListCatalogs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalogSearchResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateCatalog : Create a catalog
// Create a catalog for a given account.
func (catalogManagement *CatalogManagementV1) CreateCatalog(createCatalogOptions *CreateCatalogOptions) (result *Catalog, response *core.DetailedResponse, err error) {
	return catalogManagement.CreateCatalogWithContext(context.Background(), createCatalogOptions)
}

// CreateCatalogWithContext is an alternate form of the CreateCatalog method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) CreateCatalogWithContext(ctx context.Context, createCatalogOptions *CreateCatalogOptions) (result *Catalog, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createCatalogOptions, "createCatalogOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createCatalogOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "CreateCatalog")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createCatalogOptions.ID != nil {
		body["id"] = createCatalogOptions.ID
	}
	if createCatalogOptions.Rev != nil {
		body["_rev"] = createCatalogOptions.Rev
	}
	if createCatalogOptions.Label != nil {
		body["label"] = createCatalogOptions.Label
	}
	if createCatalogOptions.ShortDescription != nil {
		body["short_description"] = createCatalogOptions.ShortDescription
	}
	if createCatalogOptions.CatalogIconURL != nil {
		body["catalog_icon_url"] = createCatalogOptions.CatalogIconURL
	}
	if createCatalogOptions.Tags != nil {
		body["tags"] = createCatalogOptions.Tags
	}
	if createCatalogOptions.Features != nil {
		body["features"] = createCatalogOptions.Features
	}
	if createCatalogOptions.Disabled != nil {
		body["disabled"] = createCatalogOptions.Disabled
	}
	if createCatalogOptions.ResourceGroupID != nil {
		body["resource_group_id"] = createCatalogOptions.ResourceGroupID
	}
	if createCatalogOptions.OwningAccount != nil {
		body["owning_account"] = createCatalogOptions.OwningAccount
	}
	if createCatalogOptions.CatalogFilters != nil {
		body["catalog_filters"] = createCatalogOptions.CatalogFilters
	}
	if createCatalogOptions.SyndicationSettings != nil {
		body["syndication_settings"] = createCatalogOptions.SyndicationSettings
	}
	if createCatalogOptions.Kind != nil {
		body["kind"] = createCatalogOptions.Kind
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCatalog : Get catalog
// Get a catalog. This can also be used by an unauthenticated user to get the public catalog.
func (catalogManagement *CatalogManagementV1) GetCatalog(getCatalogOptions *GetCatalogOptions) (result *Catalog, response *core.DetailedResponse, err error) {
	return catalogManagement.GetCatalogWithContext(context.Background(), getCatalogOptions)
}

// GetCatalogWithContext is an alternate form of the GetCatalog method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetCatalogWithContext(ctx context.Context, getCatalogOptions *GetCatalogOptions) (result *Catalog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCatalogOptions, "getCatalogOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCatalogOptions, "getCatalogOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getCatalogOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCatalogOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetCatalog")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceCatalog : Update catalog
// Update a catalog.
func (catalogManagement *CatalogManagementV1) ReplaceCatalog(replaceCatalogOptions *ReplaceCatalogOptions) (result *Catalog, response *core.DetailedResponse, err error) {
	return catalogManagement.ReplaceCatalogWithContext(context.Background(), replaceCatalogOptions)
}

// ReplaceCatalogWithContext is an alternate form of the ReplaceCatalog method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ReplaceCatalogWithContext(ctx context.Context, replaceCatalogOptions *ReplaceCatalogOptions) (result *Catalog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceCatalogOptions, "replaceCatalogOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceCatalogOptions, "replaceCatalogOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *replaceCatalogOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceCatalogOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ReplaceCatalog")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceCatalogOptions.ID != nil {
		body["id"] = replaceCatalogOptions.ID
	}
	if replaceCatalogOptions.Rev != nil {
		body["_rev"] = replaceCatalogOptions.Rev
	}
	if replaceCatalogOptions.Label != nil {
		body["label"] = replaceCatalogOptions.Label
	}
	if replaceCatalogOptions.ShortDescription != nil {
		body["short_description"] = replaceCatalogOptions.ShortDescription
	}
	if replaceCatalogOptions.CatalogIconURL != nil {
		body["catalog_icon_url"] = replaceCatalogOptions.CatalogIconURL
	}
	if replaceCatalogOptions.Tags != nil {
		body["tags"] = replaceCatalogOptions.Tags
	}
	if replaceCatalogOptions.Features != nil {
		body["features"] = replaceCatalogOptions.Features
	}
	if replaceCatalogOptions.Disabled != nil {
		body["disabled"] = replaceCatalogOptions.Disabled
	}
	if replaceCatalogOptions.ResourceGroupID != nil {
		body["resource_group_id"] = replaceCatalogOptions.ResourceGroupID
	}
	if replaceCatalogOptions.OwningAccount != nil {
		body["owning_account"] = replaceCatalogOptions.OwningAccount
	}
	if replaceCatalogOptions.CatalogFilters != nil {
		body["catalog_filters"] = replaceCatalogOptions.CatalogFilters
	}
	if replaceCatalogOptions.SyndicationSettings != nil {
		body["syndication_settings"] = replaceCatalogOptions.SyndicationSettings
	}
	if replaceCatalogOptions.Kind != nil {
		body["kind"] = replaceCatalogOptions.Kind
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCatalog : Delete catalog
// Delete a catalog.
func (catalogManagement *CatalogManagementV1) DeleteCatalog(deleteCatalogOptions *DeleteCatalogOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteCatalogWithContext(context.Background(), deleteCatalogOptions)
}

// DeleteCatalogWithContext is an alternate form of the DeleteCatalog method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteCatalogWithContext(ctx context.Context, deleteCatalogOptions *DeleteCatalogOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCatalogOptions, "deleteCatalogOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCatalogOptions, "deleteCatalogOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *deleteCatalogOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCatalogOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteCatalog")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetCatalogAudit : Get catalog audit log
// Get the audit log associated with a catalog.
func (catalogManagement *CatalogManagementV1) GetCatalogAudit(getCatalogAuditOptions *GetCatalogAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	return catalogManagement.GetCatalogAuditWithContext(context.Background(), getCatalogAuditOptions)
}

// GetCatalogAuditWithContext is an alternate form of the GetCatalogAudit method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetCatalogAuditWithContext(ctx context.Context, getCatalogAuditOptions *GetCatalogAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCatalogAuditOptions, "getCatalogAuditOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCatalogAuditOptions, "getCatalogAuditOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getCatalogAuditOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/audit`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCatalogAuditOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetCatalogAudit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetConsumptionOfferings : Get consumption offerings
// Retrieve the available offerings from both public and from the account that currently scoped for consumption. These
// copies cannot be used for updating. They are not complete and only return what is visible to the caller. This can be
// used by an unauthenticated user to retreive publicly available offerings.
func (catalogManagement *CatalogManagementV1) GetConsumptionOfferings(getConsumptionOfferingsOptions *GetConsumptionOfferingsOptions) (result *OfferingSearchResult, response *core.DetailedResponse, err error) {
	return catalogManagement.GetConsumptionOfferingsWithContext(context.Background(), getConsumptionOfferingsOptions)
}

// GetConsumptionOfferingsWithContext is an alternate form of the GetConsumptionOfferings method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetConsumptionOfferingsWithContext(ctx context.Context, getConsumptionOfferingsOptions *GetConsumptionOfferingsOptions) (result *OfferingSearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getConsumptionOfferingsOptions, "getConsumptionOfferingsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/offerings`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getConsumptionOfferingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetConsumptionOfferings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getConsumptionOfferingsOptions.Digest != nil {
		builder.AddQuery("digest", fmt.Sprint(*getConsumptionOfferingsOptions.Digest))
	}
	if getConsumptionOfferingsOptions.Catalog != nil {
		builder.AddQuery("catalog", fmt.Sprint(*getConsumptionOfferingsOptions.Catalog))
	}
	if getConsumptionOfferingsOptions.Select != nil {
		builder.AddQuery("select", fmt.Sprint(*getConsumptionOfferingsOptions.Select))
	}
	if getConsumptionOfferingsOptions.IncludeHidden != nil {
		builder.AddQuery("includeHidden", fmt.Sprint(*getConsumptionOfferingsOptions.IncludeHidden))
	}
	if getConsumptionOfferingsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getConsumptionOfferingsOptions.Limit))
	}
	if getConsumptionOfferingsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*getConsumptionOfferingsOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOfferingSearchResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListOfferings : Get list of offerings
// Retrieve the available offerings in the specified catalog. This can also be used by an unauthenticated user to
// retreive publicly available offerings.
func (catalogManagement *CatalogManagementV1) ListOfferings(listOfferingsOptions *ListOfferingsOptions) (result *OfferingSearchResult, response *core.DetailedResponse, err error) {
	return catalogManagement.ListOfferingsWithContext(context.Background(), listOfferingsOptions)
}

// ListOfferingsWithContext is an alternate form of the ListOfferings method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListOfferingsWithContext(ctx context.Context, listOfferingsOptions *ListOfferingsOptions) (result *OfferingSearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listOfferingsOptions, "listOfferingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listOfferingsOptions, "listOfferingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *listOfferingsOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listOfferingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListOfferings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listOfferingsOptions.Digest != nil {
		builder.AddQuery("digest", fmt.Sprint(*listOfferingsOptions.Digest))
	}
	if listOfferingsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listOfferingsOptions.Limit))
	}
	if listOfferingsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listOfferingsOptions.Offset))
	}
	if listOfferingsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listOfferingsOptions.Name))
	}
	if listOfferingsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listOfferingsOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOfferingSearchResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateOffering : Create offering
// Create an offering.
func (catalogManagement *CatalogManagementV1) CreateOffering(createOfferingOptions *CreateOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.CreateOfferingWithContext(context.Background(), createOfferingOptions)
}

// CreateOfferingWithContext is an alternate form of the CreateOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) CreateOfferingWithContext(ctx context.Context, createOfferingOptions *CreateOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createOfferingOptions, "createOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createOfferingOptions, "createOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *createOfferingOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "CreateOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createOfferingOptions.ID != nil {
		body["id"] = createOfferingOptions.ID
	}
	if createOfferingOptions.Rev != nil {
		body["_rev"] = createOfferingOptions.Rev
	}
	if createOfferingOptions.URL != nil {
		body["url"] = createOfferingOptions.URL
	}
	if createOfferingOptions.CRN != nil {
		body["crn"] = createOfferingOptions.CRN
	}
	if createOfferingOptions.Label != nil {
		body["label"] = createOfferingOptions.Label
	}
	if createOfferingOptions.Name != nil {
		body["name"] = createOfferingOptions.Name
	}
	if createOfferingOptions.OfferingIconURL != nil {
		body["offering_icon_url"] = createOfferingOptions.OfferingIconURL
	}
	if createOfferingOptions.OfferingDocsURL != nil {
		body["offering_docs_url"] = createOfferingOptions.OfferingDocsURL
	}
	if createOfferingOptions.OfferingSupportURL != nil {
		body["offering_support_url"] = createOfferingOptions.OfferingSupportURL
	}
	if createOfferingOptions.Tags != nil {
		body["tags"] = createOfferingOptions.Tags
	}
	if createOfferingOptions.Keywords != nil {
		body["keywords"] = createOfferingOptions.Keywords
	}
	if createOfferingOptions.Rating != nil {
		body["rating"] = createOfferingOptions.Rating
	}
	if createOfferingOptions.Created != nil {
		body["created"] = createOfferingOptions.Created
	}
	if createOfferingOptions.Updated != nil {
		body["updated"] = createOfferingOptions.Updated
	}
	if createOfferingOptions.ShortDescription != nil {
		body["short_description"] = createOfferingOptions.ShortDescription
	}
	if createOfferingOptions.LongDescription != nil {
		body["long_description"] = createOfferingOptions.LongDescription
	}
	if createOfferingOptions.Features != nil {
		body["features"] = createOfferingOptions.Features
	}
	if createOfferingOptions.Kinds != nil {
		body["kinds"] = createOfferingOptions.Kinds
	}
	if createOfferingOptions.PermitRequestIBMPublicPublish != nil {
		body["permit_request_ibm_public_publish"] = createOfferingOptions.PermitRequestIBMPublicPublish
	}
	if createOfferingOptions.IBMPublishApproved != nil {
		body["ibm_publish_approved"] = createOfferingOptions.IBMPublishApproved
	}
	if createOfferingOptions.PublicPublishApproved != nil {
		body["public_publish_approved"] = createOfferingOptions.PublicPublishApproved
	}
	if createOfferingOptions.PublicOriginalCRN != nil {
		body["public_original_crn"] = createOfferingOptions.PublicOriginalCRN
	}
	if createOfferingOptions.PublishPublicCRN != nil {
		body["publish_public_crn"] = createOfferingOptions.PublishPublicCRN
	}
	if createOfferingOptions.PortalApprovalRecord != nil {
		body["portal_approval_record"] = createOfferingOptions.PortalApprovalRecord
	}
	if createOfferingOptions.PortalUIURL != nil {
		body["portal_ui_url"] = createOfferingOptions.PortalUIURL
	}
	if createOfferingOptions.CatalogID != nil {
		body["catalog_id"] = createOfferingOptions.CatalogID
	}
	if createOfferingOptions.CatalogName != nil {
		body["catalog_name"] = createOfferingOptions.CatalogName
	}
	if createOfferingOptions.Metadata != nil {
		body["metadata"] = createOfferingOptions.Metadata
	}
	if createOfferingOptions.Disclaimer != nil {
		body["disclaimer"] = createOfferingOptions.Disclaimer
	}
	if createOfferingOptions.Hidden != nil {
		body["hidden"] = createOfferingOptions.Hidden
	}
	if createOfferingOptions.Provider != nil {
		body["provider"] = createOfferingOptions.Provider
	}
	if createOfferingOptions.RepoInfo != nil {
		body["repo_info"] = createOfferingOptions.RepoInfo
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOffering)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ImportOfferingVersion : Import offering version
// Import new version to offering from a tgz.
func (catalogManagement *CatalogManagementV1) ImportOfferingVersion(importOfferingVersionOptions *ImportOfferingVersionOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.ImportOfferingVersionWithContext(context.Background(), importOfferingVersionOptions)
}

// ImportOfferingVersionWithContext is an alternate form of the ImportOfferingVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ImportOfferingVersionWithContext(ctx context.Context, importOfferingVersionOptions *ImportOfferingVersionOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(importOfferingVersionOptions, "importOfferingVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(importOfferingVersionOptions, "importOfferingVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *importOfferingVersionOptions.CatalogIdentifier,
		"offering_id": *importOfferingVersionOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/version`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range importOfferingVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ImportOfferingVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if importOfferingVersionOptions.Zipurl != nil {
		builder.AddQuery("zipurl", fmt.Sprint(*importOfferingVersionOptions.Zipurl))
	}
	if importOfferingVersionOptions.TargetVersion != nil {
		builder.AddQuery("targetVersion", fmt.Sprint(*importOfferingVersionOptions.TargetVersion))
	}
	if importOfferingVersionOptions.IncludeConfig != nil {
		builder.AddQuery("includeConfig", fmt.Sprint(*importOfferingVersionOptions.IncludeConfig))
	}
	if importOfferingVersionOptions.IsVsi != nil {
		builder.AddQuery("isVSI", fmt.Sprint(*importOfferingVersionOptions.IsVsi))
	}
	if importOfferingVersionOptions.RepoType != nil {
		builder.AddQuery("repoType", fmt.Sprint(*importOfferingVersionOptions.RepoType))
	}

	body := make(map[string]interface{})
	if importOfferingVersionOptions.Tags != nil {
		body["tags"] = importOfferingVersionOptions.Tags
	}
	if importOfferingVersionOptions.TargetKinds != nil {
		body["target_kinds"] = importOfferingVersionOptions.TargetKinds
	}
	if importOfferingVersionOptions.Content != nil {
		body["content"] = importOfferingVersionOptions.Content
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOffering)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ImportOffering : Import offering
// Import a new offering from a tgz.
func (catalogManagement *CatalogManagementV1) ImportOffering(importOfferingOptions *ImportOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.ImportOfferingWithContext(context.Background(), importOfferingOptions)
}

// ImportOfferingWithContext is an alternate form of the ImportOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ImportOfferingWithContext(ctx context.Context, importOfferingOptions *ImportOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(importOfferingOptions, "importOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(importOfferingOptions, "importOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *importOfferingOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/import/offerings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range importOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ImportOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if importOfferingOptions.XAuthToken != nil {
		builder.AddHeader("X-Auth-Token", fmt.Sprint(*importOfferingOptions.XAuthToken))
	}

	if importOfferingOptions.Zipurl != nil {
		builder.AddQuery("zipurl", fmt.Sprint(*importOfferingOptions.Zipurl))
	}
	if importOfferingOptions.OfferingID != nil {
		builder.AddQuery("offeringID", fmt.Sprint(*importOfferingOptions.OfferingID))
	}
	if importOfferingOptions.TargetVersion != nil {
		builder.AddQuery("targetVersion", fmt.Sprint(*importOfferingOptions.TargetVersion))
	}
	if importOfferingOptions.IncludeConfig != nil {
		builder.AddQuery("includeConfig", fmt.Sprint(*importOfferingOptions.IncludeConfig))
	}
	if importOfferingOptions.IsVsi != nil {
		builder.AddQuery("isVSI", fmt.Sprint(*importOfferingOptions.IsVsi))
	}
	if importOfferingOptions.RepoType != nil {
		builder.AddQuery("repoType", fmt.Sprint(*importOfferingOptions.RepoType))
	}

	body := make(map[string]interface{})
	if importOfferingOptions.Tags != nil {
		body["tags"] = importOfferingOptions.Tags
	}
	if importOfferingOptions.TargetKinds != nil {
		body["target_kinds"] = importOfferingOptions.TargetKinds
	}
	if importOfferingOptions.Content != nil {
		body["content"] = importOfferingOptions.Content
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOffering)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReloadOffering : Reload offering
// Reload an existing version in offering from a tgz.
func (catalogManagement *CatalogManagementV1) ReloadOffering(reloadOfferingOptions *ReloadOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.ReloadOfferingWithContext(context.Background(), reloadOfferingOptions)
}

// ReloadOfferingWithContext is an alternate form of the ReloadOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ReloadOfferingWithContext(ctx context.Context, reloadOfferingOptions *ReloadOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(reloadOfferingOptions, "reloadOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(reloadOfferingOptions, "reloadOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *reloadOfferingOptions.CatalogIdentifier,
		"offering_id": *reloadOfferingOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/reload`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range reloadOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ReloadOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	builder.AddQuery("targetVersion", fmt.Sprint(*reloadOfferingOptions.TargetVersion))
	if reloadOfferingOptions.Zipurl != nil {
		builder.AddQuery("zipurl", fmt.Sprint(*reloadOfferingOptions.Zipurl))
	}
	if reloadOfferingOptions.RepoType != nil {
		builder.AddQuery("repoType", fmt.Sprint(*reloadOfferingOptions.RepoType))
	}

	body := make(map[string]interface{})
	if reloadOfferingOptions.Tags != nil {
		body["tags"] = reloadOfferingOptions.Tags
	}
	if reloadOfferingOptions.TargetKinds != nil {
		body["target_kinds"] = reloadOfferingOptions.TargetKinds
	}
	if reloadOfferingOptions.Content != nil {
		body["content"] = reloadOfferingOptions.Content
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOffering)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOffering : Get offering
// Get an offering. This can be used by an unauthenticated user for publicly available offerings.
func (catalogManagement *CatalogManagementV1) GetOffering(getOfferingOptions *GetOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingWithContext(context.Background(), getOfferingOptions)
}

// GetOfferingWithContext is an alternate form of the GetOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingWithContext(ctx context.Context, getOfferingOptions *GetOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingOptions, "getOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingOptions, "getOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getOfferingOptions.CatalogIdentifier,
		"offering_id": *getOfferingOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOffering)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceOffering : Update offering
// Update an offering.
func (catalogManagement *CatalogManagementV1) ReplaceOffering(replaceOfferingOptions *ReplaceOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.ReplaceOfferingWithContext(context.Background(), replaceOfferingOptions)
}

// ReplaceOfferingWithContext is an alternate form of the ReplaceOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ReplaceOfferingWithContext(ctx context.Context, replaceOfferingOptions *ReplaceOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceOfferingOptions, "replaceOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceOfferingOptions, "replaceOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *replaceOfferingOptions.CatalogIdentifier,
		"offering_id": *replaceOfferingOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ReplaceOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceOfferingOptions.ID != nil {
		body["id"] = replaceOfferingOptions.ID
	}
	if replaceOfferingOptions.Rev != nil {
		body["_rev"] = replaceOfferingOptions.Rev
	}
	if replaceOfferingOptions.URL != nil {
		body["url"] = replaceOfferingOptions.URL
	}
	if replaceOfferingOptions.CRN != nil {
		body["crn"] = replaceOfferingOptions.CRN
	}
	if replaceOfferingOptions.Label != nil {
		body["label"] = replaceOfferingOptions.Label
	}
	if replaceOfferingOptions.Name != nil {
		body["name"] = replaceOfferingOptions.Name
	}
	if replaceOfferingOptions.OfferingIconURL != nil {
		body["offering_icon_url"] = replaceOfferingOptions.OfferingIconURL
	}
	if replaceOfferingOptions.OfferingDocsURL != nil {
		body["offering_docs_url"] = replaceOfferingOptions.OfferingDocsURL
	}
	if replaceOfferingOptions.OfferingSupportURL != nil {
		body["offering_support_url"] = replaceOfferingOptions.OfferingSupportURL
	}
	if replaceOfferingOptions.Tags != nil {
		body["tags"] = replaceOfferingOptions.Tags
	}
	if replaceOfferingOptions.Keywords != nil {
		body["keywords"] = replaceOfferingOptions.Keywords
	}
	if replaceOfferingOptions.Rating != nil {
		body["rating"] = replaceOfferingOptions.Rating
	}
	if replaceOfferingOptions.Created != nil {
		body["created"] = replaceOfferingOptions.Created
	}
	if replaceOfferingOptions.Updated != nil {
		body["updated"] = replaceOfferingOptions.Updated
	}
	if replaceOfferingOptions.ShortDescription != nil {
		body["short_description"] = replaceOfferingOptions.ShortDescription
	}
	if replaceOfferingOptions.LongDescription != nil {
		body["long_description"] = replaceOfferingOptions.LongDescription
	}
	if replaceOfferingOptions.Features != nil {
		body["features"] = replaceOfferingOptions.Features
	}
	if replaceOfferingOptions.Kinds != nil {
		body["kinds"] = replaceOfferingOptions.Kinds
	}
	if replaceOfferingOptions.PermitRequestIBMPublicPublish != nil {
		body["permit_request_ibm_public_publish"] = replaceOfferingOptions.PermitRequestIBMPublicPublish
	}
	if replaceOfferingOptions.IBMPublishApproved != nil {
		body["ibm_publish_approved"] = replaceOfferingOptions.IBMPublishApproved
	}
	if replaceOfferingOptions.PublicPublishApproved != nil {
		body["public_publish_approved"] = replaceOfferingOptions.PublicPublishApproved
	}
	if replaceOfferingOptions.PublicOriginalCRN != nil {
		body["public_original_crn"] = replaceOfferingOptions.PublicOriginalCRN
	}
	if replaceOfferingOptions.PublishPublicCRN != nil {
		body["publish_public_crn"] = replaceOfferingOptions.PublishPublicCRN
	}
	if replaceOfferingOptions.PortalApprovalRecord != nil {
		body["portal_approval_record"] = replaceOfferingOptions.PortalApprovalRecord
	}
	if replaceOfferingOptions.PortalUIURL != nil {
		body["portal_ui_url"] = replaceOfferingOptions.PortalUIURL
	}
	if replaceOfferingOptions.CatalogID != nil {
		body["catalog_id"] = replaceOfferingOptions.CatalogID
	}
	if replaceOfferingOptions.CatalogName != nil {
		body["catalog_name"] = replaceOfferingOptions.CatalogName
	}
	if replaceOfferingOptions.Metadata != nil {
		body["metadata"] = replaceOfferingOptions.Metadata
	}
	if replaceOfferingOptions.Disclaimer != nil {
		body["disclaimer"] = replaceOfferingOptions.Disclaimer
	}
	if replaceOfferingOptions.Hidden != nil {
		body["hidden"] = replaceOfferingOptions.Hidden
	}
	if replaceOfferingOptions.Provider != nil {
		body["provider"] = replaceOfferingOptions.Provider
	}
	if replaceOfferingOptions.RepoInfo != nil {
		body["repo_info"] = replaceOfferingOptions.RepoInfo
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOffering)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteOffering : Delete offering
// Delete an offering.
func (catalogManagement *CatalogManagementV1) DeleteOffering(deleteOfferingOptions *DeleteOfferingOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteOfferingWithContext(context.Background(), deleteOfferingOptions)
}

// DeleteOfferingWithContext is an alternate form of the DeleteOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteOfferingWithContext(ctx context.Context, deleteOfferingOptions *DeleteOfferingOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteOfferingOptions, "deleteOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteOfferingOptions, "deleteOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *deleteOfferingOptions.CatalogIdentifier,
		"offering_id": *deleteOfferingOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetOfferingAudit : Get offering audit log
// Get the audit log associated with an offering.
func (catalogManagement *CatalogManagementV1) GetOfferingAudit(getOfferingAuditOptions *GetOfferingAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingAuditWithContext(context.Background(), getOfferingAuditOptions)
}

// GetOfferingAuditWithContext is an alternate form of the GetOfferingAudit method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingAuditWithContext(ctx context.Context, getOfferingAuditOptions *GetOfferingAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingAuditOptions, "getOfferingAuditOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingAuditOptions, "getOfferingAuditOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getOfferingAuditOptions.CatalogIdentifier,
		"offering_id": *getOfferingAuditOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/audit`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingAuditOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingAudit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceOfferingIcon : Upload icon for offering
// Upload an icon file to be stored in GC. File is uploaded as a binary payload - not as a form.
func (catalogManagement *CatalogManagementV1) ReplaceOfferingIcon(replaceOfferingIconOptions *ReplaceOfferingIconOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.ReplaceOfferingIconWithContext(context.Background(), replaceOfferingIconOptions)
}

// ReplaceOfferingIconWithContext is an alternate form of the ReplaceOfferingIcon method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ReplaceOfferingIconWithContext(ctx context.Context, replaceOfferingIconOptions *ReplaceOfferingIconOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceOfferingIconOptions, "replaceOfferingIconOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceOfferingIconOptions, "replaceOfferingIconOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *replaceOfferingIconOptions.CatalogIdentifier,
		"offering_id": *replaceOfferingIconOptions.OfferingID,
		"file_name": *replaceOfferingIconOptions.FileName,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/icon/{file_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceOfferingIconOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ReplaceOfferingIcon")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOffering)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateOfferingIBM : Allow offering to be published
// Approve or disapprove the offering to be allowed to publish to the IBM Public Catalog. Options:
// * `allow_request` - (Allow requesting to publish to IBM)
// * `ibm` - (Allow publishing to be visible to IBM only)
// * `public` - (Allow publishing to be visible to everyone, including IBM)
//
// If disapprove `public`, then `ibm` approval will not  be changed. If disapprove `ibm` then `public` will
// automatically be disapproved. if disapprove `allow_request` then all rights to publish will be removed. This is
// because the process steps always go first through `allow` to `ibm` and then to `public`. `ibm` cannot be skipped.
// Only users with Approval IAM authority can use this. Approvers should use the catalog and offering id from the public
// catalog since they wouldn't have access to the private offering.'.
func (catalogManagement *CatalogManagementV1) UpdateOfferingIBM(updateOfferingIBMOptions *UpdateOfferingIBMOptions) (result *ApprovalResult, response *core.DetailedResponse, err error) {
	return catalogManagement.UpdateOfferingIBMWithContext(context.Background(), updateOfferingIBMOptions)
}

// UpdateOfferingIBMWithContext is an alternate form of the UpdateOfferingIBM method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) UpdateOfferingIBMWithContext(ctx context.Context, updateOfferingIBMOptions *UpdateOfferingIBMOptions) (result *ApprovalResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateOfferingIBMOptions, "updateOfferingIBMOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateOfferingIBMOptions, "updateOfferingIBMOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *updateOfferingIBMOptions.CatalogIdentifier,
		"offering_id": *updateOfferingIBMOptions.OfferingID,
		"approval_type": *updateOfferingIBMOptions.ApprovalType,
		"approved": *updateOfferingIBMOptions.Approved,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/publish/{approval_type}/{approved}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateOfferingIBMOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "UpdateOfferingIBM")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApprovalResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOfferingUpdates : Get version updates
// Get available updates for the specified version.
func (catalogManagement *CatalogManagementV1) GetOfferingUpdates(getOfferingUpdatesOptions *GetOfferingUpdatesOptions) (result []VersionUpdateDescriptor, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingUpdatesWithContext(context.Background(), getOfferingUpdatesOptions)
}

// GetOfferingUpdatesWithContext is an alternate form of the GetOfferingUpdates method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingUpdatesWithContext(ctx context.Context, getOfferingUpdatesOptions *GetOfferingUpdatesOptions) (result []VersionUpdateDescriptor, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingUpdatesOptions, "getOfferingUpdatesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingUpdatesOptions, "getOfferingUpdatesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getOfferingUpdatesOptions.CatalogIdentifier,
		"offering_id": *getOfferingUpdatesOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/updates`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingUpdatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingUpdates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("kind", fmt.Sprint(*getOfferingUpdatesOptions.Kind))
	if getOfferingUpdatesOptions.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*getOfferingUpdatesOptions.Version))
	}
	if getOfferingUpdatesOptions.ClusterID != nil {
		builder.AddQuery("cluster_id", fmt.Sprint(*getOfferingUpdatesOptions.ClusterID))
	}
	if getOfferingUpdatesOptions.Region != nil {
		builder.AddQuery("region", fmt.Sprint(*getOfferingUpdatesOptions.Region))
	}
	if getOfferingUpdatesOptions.ResourceGroupID != nil {
		builder.AddQuery("resource_group_id", fmt.Sprint(*getOfferingUpdatesOptions.ResourceGroupID))
	}
	if getOfferingUpdatesOptions.Namespace != nil {
		builder.AddQuery("namespace", fmt.Sprint(*getOfferingUpdatesOptions.Namespace))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVersionUpdateDescriptor)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOfferingAbout : Get version about information
// Get the about information, in markdown, for the current version.
func (catalogManagement *CatalogManagementV1) GetOfferingAbout(getOfferingAboutOptions *GetOfferingAboutOptions) (result *string, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingAboutWithContext(context.Background(), getOfferingAboutOptions)
}

// GetOfferingAboutWithContext is an alternate form of the GetOfferingAbout method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingAboutWithContext(ctx context.Context, getOfferingAboutOptions *GetOfferingAboutOptions) (result *string, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingAboutOptions, "getOfferingAboutOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingAboutOptions, "getOfferingAboutOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *getOfferingAboutOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/about`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingAboutOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingAbout")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "text/markdown")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, &result)

	return
}

// GetOfferingLicense : Get version license content
// Get the license content for the specified license ID in the specified version.
func (catalogManagement *CatalogManagementV1) GetOfferingLicense(getOfferingLicenseOptions *GetOfferingLicenseOptions) (result *string, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingLicenseWithContext(context.Background(), getOfferingLicenseOptions)
}

// GetOfferingLicenseWithContext is an alternate form of the GetOfferingLicense method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingLicenseWithContext(ctx context.Context, getOfferingLicenseOptions *GetOfferingLicenseOptions) (result *string, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingLicenseOptions, "getOfferingLicenseOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingLicenseOptions, "getOfferingLicenseOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *getOfferingLicenseOptions.VersionLocID,
		"license_id": *getOfferingLicenseOptions.LicenseID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/licenses/{license_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingLicenseOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingLicense")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "text/plain")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, &result)

	return
}

// GetOfferingContainerImages : Get version's container images
// Get the list of container images associated with the specified version. The "image_manifest_url" property of the
// version should be the URL for the image manifest, and the operation will return that content.
func (catalogManagement *CatalogManagementV1) GetOfferingContainerImages(getOfferingContainerImagesOptions *GetOfferingContainerImagesOptions) (result *ImageManifest, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingContainerImagesWithContext(context.Background(), getOfferingContainerImagesOptions)
}

// GetOfferingContainerImagesWithContext is an alternate form of the GetOfferingContainerImages method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingContainerImagesWithContext(ctx context.Context, getOfferingContainerImagesOptions *GetOfferingContainerImagesOptions) (result *ImageManifest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingContainerImagesOptions, "getOfferingContainerImagesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingContainerImagesOptions, "getOfferingContainerImagesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *getOfferingContainerImagesOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/containerImages`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingContainerImagesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingContainerImages")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImageManifest)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeprecateVersion : Deprecate version
// Deprecate the specified version.
func (catalogManagement *CatalogManagementV1) DeprecateVersion(deprecateVersionOptions *DeprecateVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeprecateVersionWithContext(context.Background(), deprecateVersionOptions)
}

// DeprecateVersionWithContext is an alternate form of the DeprecateVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeprecateVersionWithContext(ctx context.Context, deprecateVersionOptions *DeprecateVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deprecateVersionOptions, "deprecateVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deprecateVersionOptions, "deprecateVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *deprecateVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/deprecate`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deprecateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeprecateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// AccountPublishVersion : Publish version to account members
// Publish the specified version so it is viewable by account members.
func (catalogManagement *CatalogManagementV1) AccountPublishVersion(accountPublishVersionOptions *AccountPublishVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.AccountPublishVersionWithContext(context.Background(), accountPublishVersionOptions)
}

// AccountPublishVersionWithContext is an alternate form of the AccountPublishVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) AccountPublishVersionWithContext(ctx context.Context, accountPublishVersionOptions *AccountPublishVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(accountPublishVersionOptions, "accountPublishVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(accountPublishVersionOptions, "accountPublishVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *accountPublishVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/account-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range accountPublishVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "AccountPublishVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// IBMPublishVersion : Publish version to IBMers in public catalog
// Publish the specified version so that it is visible to IBMers in the public catalog.
func (catalogManagement *CatalogManagementV1) IBMPublishVersion(ibmPublishVersionOptions *IBMPublishVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.IBMPublishVersionWithContext(context.Background(), ibmPublishVersionOptions)
}

// IBMPublishVersionWithContext is an alternate form of the IBMPublishVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) IBMPublishVersionWithContext(ctx context.Context, ibmPublishVersionOptions *IBMPublishVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(ibmPublishVersionOptions, "ibmPublishVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(ibmPublishVersionOptions, "ibmPublishVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *ibmPublishVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/ibm-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range ibmPublishVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "IBMPublishVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// PublicPublishVersion : Publish version to all users in public catalog
// Publish the specified version so it is visible to all users in the public catalog.
func (catalogManagement *CatalogManagementV1) PublicPublishVersion(publicPublishVersionOptions *PublicPublishVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.PublicPublishVersionWithContext(context.Background(), publicPublishVersionOptions)
}

// PublicPublishVersionWithContext is an alternate form of the PublicPublishVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) PublicPublishVersionWithContext(ctx context.Context, publicPublishVersionOptions *PublicPublishVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(publicPublishVersionOptions, "publicPublishVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(publicPublishVersionOptions, "publicPublishVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *publicPublishVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/public-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range publicPublishVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "PublicPublishVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// CommitVersion : Commit version
// Commit a working copy of the specified version.
func (catalogManagement *CatalogManagementV1) CommitVersion(commitVersionOptions *CommitVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.CommitVersionWithContext(context.Background(), commitVersionOptions)
}

// CommitVersionWithContext is an alternate form of the CommitVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) CommitVersionWithContext(ctx context.Context, commitVersionOptions *CommitVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(commitVersionOptions, "commitVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(commitVersionOptions, "commitVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *commitVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/commit`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range commitVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "CommitVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// CopyVersion : Copy version to new target kind
// Copy the specified version to a new target kind within the same offering.
func (catalogManagement *CatalogManagementV1) CopyVersion(copyVersionOptions *CopyVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.CopyVersionWithContext(context.Background(), copyVersionOptions)
}

// CopyVersionWithContext is an alternate form of the CopyVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) CopyVersionWithContext(ctx context.Context, copyVersionOptions *CopyVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(copyVersionOptions, "copyVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(copyVersionOptions, "copyVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *copyVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/copy`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range copyVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "CopyVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if copyVersionOptions.Tags != nil {
		body["tags"] = copyVersionOptions.Tags
	}
	if copyVersionOptions.TargetKinds != nil {
		body["target_kinds"] = copyVersionOptions.TargetKinds
	}
	if copyVersionOptions.Content != nil {
		body["content"] = copyVersionOptions.Content
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetOfferingWorkingCopy : Create working copy of version
// Create a working copy of the specified version.
func (catalogManagement *CatalogManagementV1) GetOfferingWorkingCopy(getOfferingWorkingCopyOptions *GetOfferingWorkingCopyOptions) (result *Version, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingWorkingCopyWithContext(context.Background(), getOfferingWorkingCopyOptions)
}

// GetOfferingWorkingCopyWithContext is an alternate form of the GetOfferingWorkingCopy method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingWorkingCopyWithContext(ctx context.Context, getOfferingWorkingCopyOptions *GetOfferingWorkingCopyOptions) (result *Version, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingWorkingCopyOptions, "getOfferingWorkingCopyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingWorkingCopyOptions, "getOfferingWorkingCopyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *getOfferingWorkingCopyOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/workingcopy`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingWorkingCopyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingWorkingCopy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetVersion : Get offering/kind/version 'branch'
// Get the Offering/Kind/Version 'branch' for the specified locator ID.
func (catalogManagement *CatalogManagementV1) GetVersion(getVersionOptions *GetVersionOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.GetVersionWithContext(context.Background(), getVersionOptions)
}

// GetVersionWithContext is an alternate form of the GetVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetVersionWithContext(ctx context.Context, getVersionOptions *GetVersionOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getVersionOptions, "getVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getVersionOptions, "getVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *getVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOffering)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteVersion : Delete version
// Delete the specified version.  If the version is an active version with a working copy, the working copy will be
// deleted as well.
func (catalogManagement *CatalogManagementV1) DeleteVersion(deleteVersionOptions *DeleteVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteVersionWithContext(context.Background(), deleteVersionOptions)
}

// DeleteVersionWithContext is an alternate form of the DeleteVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteVersionWithContext(ctx context.Context, deleteVersionOptions *DeleteVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteVersionOptions, "deleteVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteVersionOptions, "deleteVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *deleteVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetCluster : Get kubernetes cluster
// Get the contents of the specified kubernetes cluster.
func (catalogManagement *CatalogManagementV1) GetCluster(getClusterOptions *GetClusterOptions) (result *ClusterInfo, response *core.DetailedResponse, err error) {
	return catalogManagement.GetClusterWithContext(context.Background(), getClusterOptions)
}

// GetClusterWithContext is an alternate form of the GetCluster method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetClusterWithContext(ctx context.Context, getClusterOptions *GetClusterOptions) (result *ClusterInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getClusterOptions, "getClusterOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getClusterOptions, "getClusterOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"cluster_id": *getClusterOptions.ClusterID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/deploy/kubernetes/clusters/{cluster_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getClusterOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetCluster")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getClusterOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*getClusterOptions.XAuthRefreshToken))
	}

	builder.AddQuery("region", fmt.Sprint(*getClusterOptions.Region))

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalClusterInfo)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetNamespaces : Get cluster namespaces
// Get the namespaces associated with the specified kubernetes cluster.
func (catalogManagement *CatalogManagementV1) GetNamespaces(getNamespacesOptions *GetNamespacesOptions) (result *NamespaceSearchResult, response *core.DetailedResponse, err error) {
	return catalogManagement.GetNamespacesWithContext(context.Background(), getNamespacesOptions)
}

// GetNamespacesWithContext is an alternate form of the GetNamespaces method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetNamespacesWithContext(ctx context.Context, getNamespacesOptions *GetNamespacesOptions) (result *NamespaceSearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getNamespacesOptions, "getNamespacesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getNamespacesOptions, "getNamespacesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"cluster_id": *getNamespacesOptions.ClusterID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/deploy/kubernetes/clusters/{cluster_id}/namespaces`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getNamespacesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetNamespaces")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getNamespacesOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*getNamespacesOptions.XAuthRefreshToken))
	}

	builder.AddQuery("region", fmt.Sprint(*getNamespacesOptions.Region))
	if getNamespacesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getNamespacesOptions.Limit))
	}
	if getNamespacesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*getNamespacesOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNamespaceSearchResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeployOperators : Deploy operators
// Deploy operators on a kubernetes cluster.
func (catalogManagement *CatalogManagementV1) DeployOperators(deployOperatorsOptions *DeployOperatorsOptions) (result []OperatorDeployResult, response *core.DetailedResponse, err error) {
	return catalogManagement.DeployOperatorsWithContext(context.Background(), deployOperatorsOptions)
}

// DeployOperatorsWithContext is an alternate form of the DeployOperators method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeployOperatorsWithContext(ctx context.Context, deployOperatorsOptions *DeployOperatorsOptions) (result []OperatorDeployResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deployOperatorsOptions, "deployOperatorsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deployOperatorsOptions, "deployOperatorsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/deploy/kubernetes/olm/operator`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range deployOperatorsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeployOperators")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if deployOperatorsOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*deployOperatorsOptions.XAuthRefreshToken))
	}

	body := make(map[string]interface{})
	if deployOperatorsOptions.ClusterID != nil {
		body["cluster_id"] = deployOperatorsOptions.ClusterID
	}
	if deployOperatorsOptions.Region != nil {
		body["region"] = deployOperatorsOptions.Region
	}
	if deployOperatorsOptions.Namespaces != nil {
		body["namespaces"] = deployOperatorsOptions.Namespaces
	}
	if deployOperatorsOptions.AllNamespaces != nil {
		body["all_namespaces"] = deployOperatorsOptions.AllNamespaces
	}
	if deployOperatorsOptions.VersionLocatorID != nil {
		body["version_locator_id"] = deployOperatorsOptions.VersionLocatorID
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOperatorDeployResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListOperators : List operators
// List the operators from a kubernetes cluster.
func (catalogManagement *CatalogManagementV1) ListOperators(listOperatorsOptions *ListOperatorsOptions) (result []OperatorDeployResult, response *core.DetailedResponse, err error) {
	return catalogManagement.ListOperatorsWithContext(context.Background(), listOperatorsOptions)
}

// ListOperatorsWithContext is an alternate form of the ListOperators method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListOperatorsWithContext(ctx context.Context, listOperatorsOptions *ListOperatorsOptions) (result []OperatorDeployResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listOperatorsOptions, "listOperatorsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listOperatorsOptions, "listOperatorsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/deploy/kubernetes/olm/operator`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listOperatorsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListOperators")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listOperatorsOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*listOperatorsOptions.XAuthRefreshToken))
	}

	builder.AddQuery("cluster_id", fmt.Sprint(*listOperatorsOptions.ClusterID))
	builder.AddQuery("region", fmt.Sprint(*listOperatorsOptions.Region))
	builder.AddQuery("version_locator_id", fmt.Sprint(*listOperatorsOptions.VersionLocatorID))

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOperatorDeployResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceOperators : Update operators
// Update the operators on a kubernetes cluster.
func (catalogManagement *CatalogManagementV1) ReplaceOperators(replaceOperatorsOptions *ReplaceOperatorsOptions) (result []OperatorDeployResult, response *core.DetailedResponse, err error) {
	return catalogManagement.ReplaceOperatorsWithContext(context.Background(), replaceOperatorsOptions)
}

// ReplaceOperatorsWithContext is an alternate form of the ReplaceOperators method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ReplaceOperatorsWithContext(ctx context.Context, replaceOperatorsOptions *ReplaceOperatorsOptions) (result []OperatorDeployResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceOperatorsOptions, "replaceOperatorsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceOperatorsOptions, "replaceOperatorsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/deploy/kubernetes/olm/operator`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceOperatorsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ReplaceOperators")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceOperatorsOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*replaceOperatorsOptions.XAuthRefreshToken))
	}

	body := make(map[string]interface{})
	if replaceOperatorsOptions.ClusterID != nil {
		body["cluster_id"] = replaceOperatorsOptions.ClusterID
	}
	if replaceOperatorsOptions.Region != nil {
		body["region"] = replaceOperatorsOptions.Region
	}
	if replaceOperatorsOptions.Namespaces != nil {
		body["namespaces"] = replaceOperatorsOptions.Namespaces
	}
	if replaceOperatorsOptions.AllNamespaces != nil {
		body["all_namespaces"] = replaceOperatorsOptions.AllNamespaces
	}
	if replaceOperatorsOptions.VersionLocatorID != nil {
		body["version_locator_id"] = replaceOperatorsOptions.VersionLocatorID
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOperatorDeployResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteOperators : Delete operators
// Delete operators from a kubernetes cluster.
func (catalogManagement *CatalogManagementV1) DeleteOperators(deleteOperatorsOptions *DeleteOperatorsOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteOperatorsWithContext(context.Background(), deleteOperatorsOptions)
}

// DeleteOperatorsWithContext is an alternate form of the DeleteOperators method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteOperatorsWithContext(ctx context.Context, deleteOperatorsOptions *DeleteOperatorsOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteOperatorsOptions, "deleteOperatorsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteOperatorsOptions, "deleteOperatorsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/deploy/kubernetes/olm/operator`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteOperatorsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteOperators")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteOperatorsOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*deleteOperatorsOptions.XAuthRefreshToken))
	}

	builder.AddQuery("cluster_id", fmt.Sprint(*deleteOperatorsOptions.ClusterID))
	builder.AddQuery("region", fmt.Sprint(*deleteOperatorsOptions.Region))
	builder.AddQuery("version_locator_id", fmt.Sprint(*deleteOperatorsOptions.VersionLocatorID))

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// InstallVersion : Install version
// Create an install for the specified version.
func (catalogManagement *CatalogManagementV1) InstallVersion(installVersionOptions *InstallVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.InstallVersionWithContext(context.Background(), installVersionOptions)
}

// InstallVersionWithContext is an alternate form of the InstallVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) InstallVersionWithContext(ctx context.Context, installVersionOptions *InstallVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(installVersionOptions, "installVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(installVersionOptions, "installVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *installVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/install`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range installVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "InstallVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")
	if installVersionOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*installVersionOptions.XAuthRefreshToken))
	}

	body := make(map[string]interface{})
	if installVersionOptions.ClusterID != nil {
		body["cluster_id"] = installVersionOptions.ClusterID
	}
	if installVersionOptions.Region != nil {
		body["region"] = installVersionOptions.Region
	}
	if installVersionOptions.Namespace != nil {
		body["namespace"] = installVersionOptions.Namespace
	}
	if installVersionOptions.OverrideValues != nil {
		body["override_values"] = installVersionOptions.OverrideValues
	}
	if installVersionOptions.EntitlementApikey != nil {
		body["entitlement_apikey"] = installVersionOptions.EntitlementApikey
	}
	if installVersionOptions.Schematics != nil {
		body["schematics"] = installVersionOptions.Schematics
	}
	if installVersionOptions.Script != nil {
		body["script"] = installVersionOptions.Script
	}
	if installVersionOptions.ScriptID != nil {
		body["script_id"] = installVersionOptions.ScriptID
	}
	if installVersionOptions.VersionLocatorID != nil {
		body["version_locator_id"] = installVersionOptions.VersionLocatorID
	}
	if installVersionOptions.VcenterID != nil {
		body["vcenter_id"] = installVersionOptions.VcenterID
	}
	if installVersionOptions.VcenterUser != nil {
		body["vcenter_user"] = installVersionOptions.VcenterUser
	}
	if installVersionOptions.VcenterPassword != nil {
		body["vcenter_password"] = installVersionOptions.VcenterPassword
	}
	if installVersionOptions.VcenterLocation != nil {
		body["vcenter_location"] = installVersionOptions.VcenterLocation
	}
	if installVersionOptions.VcenterDatastore != nil {
		body["vcenter_datastore"] = installVersionOptions.VcenterDatastore
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// PreinstallVersion : Pre-install version
// Create a pre-install for the specified version.
func (catalogManagement *CatalogManagementV1) PreinstallVersion(preinstallVersionOptions *PreinstallVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.PreinstallVersionWithContext(context.Background(), preinstallVersionOptions)
}

// PreinstallVersionWithContext is an alternate form of the PreinstallVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) PreinstallVersionWithContext(ctx context.Context, preinstallVersionOptions *PreinstallVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(preinstallVersionOptions, "preinstallVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(preinstallVersionOptions, "preinstallVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *preinstallVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/preinstall`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range preinstallVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "PreinstallVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")
	if preinstallVersionOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*preinstallVersionOptions.XAuthRefreshToken))
	}

	body := make(map[string]interface{})
	if preinstallVersionOptions.ClusterID != nil {
		body["cluster_id"] = preinstallVersionOptions.ClusterID
	}
	if preinstallVersionOptions.Region != nil {
		body["region"] = preinstallVersionOptions.Region
	}
	if preinstallVersionOptions.Namespace != nil {
		body["namespace"] = preinstallVersionOptions.Namespace
	}
	if preinstallVersionOptions.OverrideValues != nil {
		body["override_values"] = preinstallVersionOptions.OverrideValues
	}
	if preinstallVersionOptions.EntitlementApikey != nil {
		body["entitlement_apikey"] = preinstallVersionOptions.EntitlementApikey
	}
	if preinstallVersionOptions.Schematics != nil {
		body["schematics"] = preinstallVersionOptions.Schematics
	}
	if preinstallVersionOptions.Script != nil {
		body["script"] = preinstallVersionOptions.Script
	}
	if preinstallVersionOptions.ScriptID != nil {
		body["script_id"] = preinstallVersionOptions.ScriptID
	}
	if preinstallVersionOptions.VersionLocatorID != nil {
		body["version_locator_id"] = preinstallVersionOptions.VersionLocatorID
	}
	if preinstallVersionOptions.VcenterID != nil {
		body["vcenter_id"] = preinstallVersionOptions.VcenterID
	}
	if preinstallVersionOptions.VcenterUser != nil {
		body["vcenter_user"] = preinstallVersionOptions.VcenterUser
	}
	if preinstallVersionOptions.VcenterPassword != nil {
		body["vcenter_password"] = preinstallVersionOptions.VcenterPassword
	}
	if preinstallVersionOptions.VcenterLocation != nil {
		body["vcenter_location"] = preinstallVersionOptions.VcenterLocation
	}
	if preinstallVersionOptions.VcenterDatastore != nil {
		body["vcenter_datastore"] = preinstallVersionOptions.VcenterDatastore
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetPreinstall : Get version pre-install status
// Get the pre-install status for the specified version.
func (catalogManagement *CatalogManagementV1) GetPreinstall(getPreinstallOptions *GetPreinstallOptions) (result *InstallStatus, response *core.DetailedResponse, err error) {
	return catalogManagement.GetPreinstallWithContext(context.Background(), getPreinstallOptions)
}

// GetPreinstallWithContext is an alternate form of the GetPreinstall method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetPreinstallWithContext(ctx context.Context, getPreinstallOptions *GetPreinstallOptions) (result *InstallStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPreinstallOptions, "getPreinstallOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getPreinstallOptions, "getPreinstallOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *getPreinstallOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/preinstall`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getPreinstallOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetPreinstall")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getPreinstallOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*getPreinstallOptions.XAuthRefreshToken))
	}

	if getPreinstallOptions.ClusterID != nil {
		builder.AddQuery("cluster_id", fmt.Sprint(*getPreinstallOptions.ClusterID))
	}
	if getPreinstallOptions.Region != nil {
		builder.AddQuery("region", fmt.Sprint(*getPreinstallOptions.Region))
	}
	if getPreinstallOptions.Namespace != nil {
		builder.AddQuery("namespace", fmt.Sprint(*getPreinstallOptions.Namespace))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInstallStatus)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ValidateInstall : Validate offering
// Validate the offering associated with the specified version.
func (catalogManagement *CatalogManagementV1) ValidateInstall(validateInstallOptions *ValidateInstallOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.ValidateInstallWithContext(context.Background(), validateInstallOptions)
}

// ValidateInstallWithContext is an alternate form of the ValidateInstall method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ValidateInstallWithContext(ctx context.Context, validateInstallOptions *ValidateInstallOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(validateInstallOptions, "validateInstallOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(validateInstallOptions, "validateInstallOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *validateInstallOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/validation/install`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range validateInstallOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ValidateInstall")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")
	if validateInstallOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*validateInstallOptions.XAuthRefreshToken))
	}

	body := make(map[string]interface{})
	if validateInstallOptions.ClusterID != nil {
		body["cluster_id"] = validateInstallOptions.ClusterID
	}
	if validateInstallOptions.Region != nil {
		body["region"] = validateInstallOptions.Region
	}
	if validateInstallOptions.Namespace != nil {
		body["namespace"] = validateInstallOptions.Namespace
	}
	if validateInstallOptions.OverrideValues != nil {
		body["override_values"] = validateInstallOptions.OverrideValues
	}
	if validateInstallOptions.EntitlementApikey != nil {
		body["entitlement_apikey"] = validateInstallOptions.EntitlementApikey
	}
	if validateInstallOptions.Schematics != nil {
		body["schematics"] = validateInstallOptions.Schematics
	}
	if validateInstallOptions.Script != nil {
		body["script"] = validateInstallOptions.Script
	}
	if validateInstallOptions.ScriptID != nil {
		body["script_id"] = validateInstallOptions.ScriptID
	}
	if validateInstallOptions.VersionLocatorID != nil {
		body["version_locator_id"] = validateInstallOptions.VersionLocatorID
	}
	if validateInstallOptions.VcenterID != nil {
		body["vcenter_id"] = validateInstallOptions.VcenterID
	}
	if validateInstallOptions.VcenterUser != nil {
		body["vcenter_user"] = validateInstallOptions.VcenterUser
	}
	if validateInstallOptions.VcenterPassword != nil {
		body["vcenter_password"] = validateInstallOptions.VcenterPassword
	}
	if validateInstallOptions.VcenterLocation != nil {
		body["vcenter_location"] = validateInstallOptions.VcenterLocation
	}
	if validateInstallOptions.VcenterDatastore != nil {
		body["vcenter_datastore"] = validateInstallOptions.VcenterDatastore
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetValidationStatus : Get offering install status
// Returns the install status for the specified offering version.
func (catalogManagement *CatalogManagementV1) GetValidationStatus(getValidationStatusOptions *GetValidationStatusOptions) (result *Validation, response *core.DetailedResponse, err error) {
	return catalogManagement.GetValidationStatusWithContext(context.Background(), getValidationStatusOptions)
}

// GetValidationStatusWithContext is an alternate form of the GetValidationStatus method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetValidationStatusWithContext(ctx context.Context, getValidationStatusOptions *GetValidationStatusOptions) (result *Validation, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getValidationStatusOptions, "getValidationStatusOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getValidationStatusOptions, "getValidationStatusOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *getValidationStatusOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/validation/install`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getValidationStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetValidationStatus")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getValidationStatusOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*getValidationStatusOptions.XAuthRefreshToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalValidation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOverrideValues : Get override values
// Returns the override values that were used to validate the specified offering version.
func (catalogManagement *CatalogManagementV1) GetOverrideValues(getOverrideValuesOptions *GetOverrideValuesOptions) (result map[string]interface{}, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOverrideValuesWithContext(context.Background(), getOverrideValuesOptions)
}

// GetOverrideValuesWithContext is an alternate form of the GetOverrideValues method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOverrideValuesWithContext(ctx context.Context, getOverrideValuesOptions *GetOverrideValuesOptions) (result map[string]interface{}, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOverrideValuesOptions, "getOverrideValuesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOverrideValuesOptions, "getOverrideValuesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *getOverrideValuesOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/validation/overridevalues`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOverrideValuesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOverrideValues")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, &result)

	return
}

// SearchObjects : List objects across catalogs
// List the available objects from both public and private catalogs. These copies cannot be used for updating. They are
// not complete and only return what is visible to the caller.
func (catalogManagement *CatalogManagementV1) SearchObjects(searchObjectsOptions *SearchObjectsOptions) (result *ObjectSearchResult, response *core.DetailedResponse, err error) {
	return catalogManagement.SearchObjectsWithContext(context.Background(), searchObjectsOptions)
}

// SearchObjectsWithContext is an alternate form of the SearchObjects method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) SearchObjectsWithContext(ctx context.Context, searchObjectsOptions *SearchObjectsOptions) (result *ObjectSearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(searchObjectsOptions, "searchObjectsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(searchObjectsOptions, "searchObjectsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/objects`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range searchObjectsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "SearchObjects")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("query", fmt.Sprint(*searchObjectsOptions.Query))
	if searchObjectsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*searchObjectsOptions.Limit))
	}
	if searchObjectsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*searchObjectsOptions.Offset))
	}
	if searchObjectsOptions.Collapse != nil {
		builder.AddQuery("collapse", fmt.Sprint(*searchObjectsOptions.Collapse))
	}
	if searchObjectsOptions.Digest != nil {
		builder.AddQuery("digest", fmt.Sprint(*searchObjectsOptions.Digest))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalObjectSearchResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListObjects : List objects within a catalog
// List the available objects within the specified catalog.
func (catalogManagement *CatalogManagementV1) ListObjects(listObjectsOptions *ListObjectsOptions) (result *ObjectListResult, response *core.DetailedResponse, err error) {
	return catalogManagement.ListObjectsWithContext(context.Background(), listObjectsOptions)
}

// ListObjectsWithContext is an alternate form of the ListObjects method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListObjectsWithContext(ctx context.Context, listObjectsOptions *ListObjectsOptions) (result *ObjectListResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listObjectsOptions, "listObjectsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listObjectsOptions, "listObjectsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *listObjectsOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listObjectsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListObjects")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listObjectsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listObjectsOptions.Limit))
	}
	if listObjectsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listObjectsOptions.Offset))
	}
	if listObjectsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listObjectsOptions.Name))
	}
	if listObjectsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listObjectsOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalObjectListResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateObject : Create catalog object
// Create an object with a specific catalog.
func (catalogManagement *CatalogManagementV1) CreateObject(createObjectOptions *CreateObjectOptions) (result *CatalogObject, response *core.DetailedResponse, err error) {
	return catalogManagement.CreateObjectWithContext(context.Background(), createObjectOptions)
}

// CreateObjectWithContext is an alternate form of the CreateObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) CreateObjectWithContext(ctx context.Context, createObjectOptions *CreateObjectOptions) (result *CatalogObject, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createObjectOptions, "createObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createObjectOptions, "createObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *createObjectOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "CreateObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createObjectOptions.ID != nil {
		body["id"] = createObjectOptions.ID
	}
	if createObjectOptions.Name != nil {
		body["name"] = createObjectOptions.Name
	}
	if createObjectOptions.Rev != nil {
		body["_rev"] = createObjectOptions.Rev
	}
	if createObjectOptions.CRN != nil {
		body["crn"] = createObjectOptions.CRN
	}
	if createObjectOptions.URL != nil {
		body["url"] = createObjectOptions.URL
	}
	if createObjectOptions.ParentID != nil {
		body["parent_id"] = createObjectOptions.ParentID
	}
	if createObjectOptions.LabelI18n != nil {
		body["label_i18n"] = createObjectOptions.LabelI18n
	}
	if createObjectOptions.Label != nil {
		body["label"] = createObjectOptions.Label
	}
	if createObjectOptions.Tags != nil {
		body["tags"] = createObjectOptions.Tags
	}
	if createObjectOptions.Created != nil {
		body["created"] = createObjectOptions.Created
	}
	if createObjectOptions.Updated != nil {
		body["updated"] = createObjectOptions.Updated
	}
	if createObjectOptions.ShortDescription != nil {
		body["short_description"] = createObjectOptions.ShortDescription
	}
	if createObjectOptions.ShortDescriptionI18n != nil {
		body["short_description_i18n"] = createObjectOptions.ShortDescriptionI18n
	}
	if createObjectOptions.Kind != nil {
		body["kind"] = createObjectOptions.Kind
	}
	if createObjectOptions.Publish != nil {
		body["publish"] = createObjectOptions.Publish
	}
	if createObjectOptions.State != nil {
		body["state"] = createObjectOptions.State
	}
	if createObjectOptions.CatalogID != nil {
		body["catalog_id"] = createObjectOptions.CatalogID
	}
	if createObjectOptions.CatalogName != nil {
		body["catalog_name"] = createObjectOptions.CatalogName
	}
	if createObjectOptions.Data != nil {
		body["data"] = createObjectOptions.Data
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalogObject)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetObject : Get catalog object
// Get the specified object from within the specified catalog.
func (catalogManagement *CatalogManagementV1) GetObject(getObjectOptions *GetObjectOptions) (result *CatalogObject, response *core.DetailedResponse, err error) {
	return catalogManagement.GetObjectWithContext(context.Background(), getObjectOptions)
}

// GetObjectWithContext is an alternate form of the GetObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetObjectWithContext(ctx context.Context, getObjectOptions *GetObjectOptions) (result *CatalogObject, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getObjectOptions, "getObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getObjectOptions, "getObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getObjectOptions.CatalogIdentifier,
		"object_identifier": *getObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalogObject)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceObject : Update catalog object
// Update an object within a specific catalog.
func (catalogManagement *CatalogManagementV1) ReplaceObject(replaceObjectOptions *ReplaceObjectOptions) (result *CatalogObject, response *core.DetailedResponse, err error) {
	return catalogManagement.ReplaceObjectWithContext(context.Background(), replaceObjectOptions)
}

// ReplaceObjectWithContext is an alternate form of the ReplaceObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ReplaceObjectWithContext(ctx context.Context, replaceObjectOptions *ReplaceObjectOptions) (result *CatalogObject, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceObjectOptions, "replaceObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceObjectOptions, "replaceObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *replaceObjectOptions.CatalogIdentifier,
		"object_identifier": *replaceObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ReplaceObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceObjectOptions.ID != nil {
		body["id"] = replaceObjectOptions.ID
	}
	if replaceObjectOptions.Name != nil {
		body["name"] = replaceObjectOptions.Name
	}
	if replaceObjectOptions.Rev != nil {
		body["_rev"] = replaceObjectOptions.Rev
	}
	if replaceObjectOptions.CRN != nil {
		body["crn"] = replaceObjectOptions.CRN
	}
	if replaceObjectOptions.URL != nil {
		body["url"] = replaceObjectOptions.URL
	}
	if replaceObjectOptions.ParentID != nil {
		body["parent_id"] = replaceObjectOptions.ParentID
	}
	if replaceObjectOptions.LabelI18n != nil {
		body["label_i18n"] = replaceObjectOptions.LabelI18n
	}
	if replaceObjectOptions.Label != nil {
		body["label"] = replaceObjectOptions.Label
	}
	if replaceObjectOptions.Tags != nil {
		body["tags"] = replaceObjectOptions.Tags
	}
	if replaceObjectOptions.Created != nil {
		body["created"] = replaceObjectOptions.Created
	}
	if replaceObjectOptions.Updated != nil {
		body["updated"] = replaceObjectOptions.Updated
	}
	if replaceObjectOptions.ShortDescription != nil {
		body["short_description"] = replaceObjectOptions.ShortDescription
	}
	if replaceObjectOptions.ShortDescriptionI18n != nil {
		body["short_description_i18n"] = replaceObjectOptions.ShortDescriptionI18n
	}
	if replaceObjectOptions.Kind != nil {
		body["kind"] = replaceObjectOptions.Kind
	}
	if replaceObjectOptions.Publish != nil {
		body["publish"] = replaceObjectOptions.Publish
	}
	if replaceObjectOptions.State != nil {
		body["state"] = replaceObjectOptions.State
	}
	if replaceObjectOptions.CatalogID != nil {
		body["catalog_id"] = replaceObjectOptions.CatalogID
	}
	if replaceObjectOptions.CatalogName != nil {
		body["catalog_name"] = replaceObjectOptions.CatalogName
	}
	if replaceObjectOptions.Data != nil {
		body["data"] = replaceObjectOptions.Data
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCatalogObject)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteObject : Delete catalog object
// Delete a specific object within a specific catalog.
func (catalogManagement *CatalogManagementV1) DeleteObject(deleteObjectOptions *DeleteObjectOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteObjectWithContext(context.Background(), deleteObjectOptions)
}

// DeleteObjectWithContext is an alternate form of the DeleteObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteObjectWithContext(ctx context.Context, deleteObjectOptions *DeleteObjectOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteObjectOptions, "deleteObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteObjectOptions, "deleteObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *deleteObjectOptions.CatalogIdentifier,
		"object_identifier": *deleteObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetObjectAudit : Get catalog object audit log
// Get the audit log associated with a specific catalog object.
func (catalogManagement *CatalogManagementV1) GetObjectAudit(getObjectAuditOptions *GetObjectAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	return catalogManagement.GetObjectAuditWithContext(context.Background(), getObjectAuditOptions)
}

// GetObjectAuditWithContext is an alternate form of the GetObjectAudit method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetObjectAuditWithContext(ctx context.Context, getObjectAuditOptions *GetObjectAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getObjectAuditOptions, "getObjectAuditOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getObjectAuditOptions, "getObjectAuditOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getObjectAuditOptions.CatalogIdentifier,
		"object_identifier": *getObjectAuditOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/audit`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getObjectAuditOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetObjectAudit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// AccountPublishObject : Publish object to account
// Publish a catalog object to account.
func (catalogManagement *CatalogManagementV1) AccountPublishObject(accountPublishObjectOptions *AccountPublishObjectOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.AccountPublishObjectWithContext(context.Background(), accountPublishObjectOptions)
}

// AccountPublishObjectWithContext is an alternate form of the AccountPublishObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) AccountPublishObjectWithContext(ctx context.Context, accountPublishObjectOptions *AccountPublishObjectOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(accountPublishObjectOptions, "accountPublishObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(accountPublishObjectOptions, "accountPublishObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *accountPublishObjectOptions.CatalogIdentifier,
		"object_identifier": *accountPublishObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/account-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range accountPublishObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "AccountPublishObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// SharedPublishObject : Publish object to share with allow list
// Publish the specified object so that it is visible to those in the allow list.
func (catalogManagement *CatalogManagementV1) SharedPublishObject(sharedPublishObjectOptions *SharedPublishObjectOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.SharedPublishObjectWithContext(context.Background(), sharedPublishObjectOptions)
}

// SharedPublishObjectWithContext is an alternate form of the SharedPublishObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) SharedPublishObjectWithContext(ctx context.Context, sharedPublishObjectOptions *SharedPublishObjectOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(sharedPublishObjectOptions, "sharedPublishObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(sharedPublishObjectOptions, "sharedPublishObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *sharedPublishObjectOptions.CatalogIdentifier,
		"object_identifier": *sharedPublishObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/shared-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range sharedPublishObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "SharedPublishObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// IBMPublishObject : Publish object to share with IBMers
// Publish the specified object so that it is visible to IBMers in the public catalog.
func (catalogManagement *CatalogManagementV1) IBMPublishObject(ibmPublishObjectOptions *IBMPublishObjectOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.IBMPublishObjectWithContext(context.Background(), ibmPublishObjectOptions)
}

// IBMPublishObjectWithContext is an alternate form of the IBMPublishObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) IBMPublishObjectWithContext(ctx context.Context, ibmPublishObjectOptions *IBMPublishObjectOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(ibmPublishObjectOptions, "ibmPublishObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(ibmPublishObjectOptions, "ibmPublishObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *ibmPublishObjectOptions.CatalogIdentifier,
		"object_identifier": *ibmPublishObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/ibm-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range ibmPublishObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "IBMPublishObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// PublicPublishObject : Publish object to share with all users
// Publish the specified object so it is visible to all users in the public catalog.
func (catalogManagement *CatalogManagementV1) PublicPublishObject(publicPublishObjectOptions *PublicPublishObjectOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.PublicPublishObjectWithContext(context.Background(), publicPublishObjectOptions)
}

// PublicPublishObjectWithContext is an alternate form of the PublicPublishObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) PublicPublishObjectWithContext(ctx context.Context, publicPublishObjectOptions *PublicPublishObjectOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(publicPublishObjectOptions, "publicPublishObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(publicPublishObjectOptions, "publicPublishObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *publicPublishObjectOptions.CatalogIdentifier,
		"object_identifier": *publicPublishObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/public-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range publicPublishObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "PublicPublishObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// CreateObjectAccess : Add account ID to object access list
// Add an account ID to an object's access list.
func (catalogManagement *CatalogManagementV1) CreateObjectAccess(createObjectAccessOptions *CreateObjectAccessOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.CreateObjectAccessWithContext(context.Background(), createObjectAccessOptions)
}

// CreateObjectAccessWithContext is an alternate form of the CreateObjectAccess method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) CreateObjectAccessWithContext(ctx context.Context, createObjectAccessOptions *CreateObjectAccessOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createObjectAccessOptions, "createObjectAccessOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createObjectAccessOptions, "createObjectAccessOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *createObjectAccessOptions.CatalogIdentifier,
		"object_identifier": *createObjectAccessOptions.ObjectIdentifier,
		"account_identifier": *createObjectAccessOptions.AccountIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access/{account_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createObjectAccessOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "CreateObjectAccess")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetObjectAccess : Check for account ID in object access list
// Determine if an account ID is in an object's access list.
func (catalogManagement *CatalogManagementV1) GetObjectAccess(getObjectAccessOptions *GetObjectAccessOptions) (result *ObjectAccess, response *core.DetailedResponse, err error) {
	return catalogManagement.GetObjectAccessWithContext(context.Background(), getObjectAccessOptions)
}

// GetObjectAccessWithContext is an alternate form of the GetObjectAccess method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetObjectAccessWithContext(ctx context.Context, getObjectAccessOptions *GetObjectAccessOptions) (result *ObjectAccess, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getObjectAccessOptions, "getObjectAccessOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getObjectAccessOptions, "getObjectAccessOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getObjectAccessOptions.CatalogIdentifier,
		"object_identifier": *getObjectAccessOptions.ObjectIdentifier,
		"account_identifier": *getObjectAccessOptions.AccountIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access/{account_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getObjectAccessOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetObjectAccess")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalObjectAccess)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteObjectAccess : Remove account ID from object access list
// Delete the specified account ID from the specified object's access list.
func (catalogManagement *CatalogManagementV1) DeleteObjectAccess(deleteObjectAccessOptions *DeleteObjectAccessOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteObjectAccessWithContext(context.Background(), deleteObjectAccessOptions)
}

// DeleteObjectAccessWithContext is an alternate form of the DeleteObjectAccess method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteObjectAccessWithContext(ctx context.Context, deleteObjectAccessOptions *DeleteObjectAccessOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteObjectAccessOptions, "deleteObjectAccessOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteObjectAccessOptions, "deleteObjectAccessOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *deleteObjectAccessOptions.CatalogIdentifier,
		"object_identifier": *deleteObjectAccessOptions.ObjectIdentifier,
		"account_identifier": *deleteObjectAccessOptions.AccountIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access/{account_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteObjectAccessOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteObjectAccess")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// GetObjectAccessList : Get object access list
// Get the access list associated with the specified object.
func (catalogManagement *CatalogManagementV1) GetObjectAccessList(getObjectAccessListOptions *GetObjectAccessListOptions) (result *ObjectAccessListResult, response *core.DetailedResponse, err error) {
	return catalogManagement.GetObjectAccessListWithContext(context.Background(), getObjectAccessListOptions)
}

// GetObjectAccessListWithContext is an alternate form of the GetObjectAccessList method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetObjectAccessListWithContext(ctx context.Context, getObjectAccessListOptions *GetObjectAccessListOptions) (result *ObjectAccessListResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getObjectAccessListOptions, "getObjectAccessListOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getObjectAccessListOptions, "getObjectAccessListOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getObjectAccessListOptions.CatalogIdentifier,
		"object_identifier": *getObjectAccessListOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getObjectAccessListOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetObjectAccessList")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getObjectAccessListOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getObjectAccessListOptions.Limit))
	}
	if getObjectAccessListOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*getObjectAccessListOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalObjectAccessListResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteObjectAccessList : Delete accounts from object access list
// Delete all or a set of accounts from an object's access list.
func (catalogManagement *CatalogManagementV1) DeleteObjectAccessList(deleteObjectAccessListOptions *DeleteObjectAccessListOptions) (result *AccessListBulkResponse, response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteObjectAccessListWithContext(context.Background(), deleteObjectAccessListOptions)
}

// DeleteObjectAccessListWithContext is an alternate form of the DeleteObjectAccessList method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteObjectAccessListWithContext(ctx context.Context, deleteObjectAccessListOptions *DeleteObjectAccessListOptions) (result *AccessListBulkResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteObjectAccessListOptions, "deleteObjectAccessListOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteObjectAccessListOptions, "deleteObjectAccessListOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *deleteObjectAccessListOptions.CatalogIdentifier,
		"object_identifier": *deleteObjectAccessListOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteObjectAccessListOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteObjectAccessList")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(deleteObjectAccessListOptions.Accounts)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessListBulkResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// AddObjectAccessList : Add accounts to object access list
// Add one or more accounts to the specified object's access list.
func (catalogManagement *CatalogManagementV1) AddObjectAccessList(addObjectAccessListOptions *AddObjectAccessListOptions) (result *AccessListBulkResponse, response *core.DetailedResponse, err error) {
	return catalogManagement.AddObjectAccessListWithContext(context.Background(), addObjectAccessListOptions)
}

// AddObjectAccessListWithContext is an alternate form of the AddObjectAccessList method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) AddObjectAccessListWithContext(ctx context.Context, addObjectAccessListOptions *AddObjectAccessListOptions) (result *AccessListBulkResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(addObjectAccessListOptions, "addObjectAccessListOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(addObjectAccessListOptions, "addObjectAccessListOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *addObjectAccessListOptions.CatalogIdentifier,
		"object_identifier": *addObjectAccessListOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range addObjectAccessListOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "AddObjectAccessList")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(addObjectAccessListOptions.Accounts)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessListBulkResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateOfferingInstance : Create an offering resource instance
// Provision a new offering in a given account, and return its resource instance.
func (catalogManagement *CatalogManagementV1) CreateOfferingInstance(createOfferingInstanceOptions *CreateOfferingInstanceOptions) (result *OfferingInstance, response *core.DetailedResponse, err error) {
	return catalogManagement.CreateOfferingInstanceWithContext(context.Background(), createOfferingInstanceOptions)
}

// CreateOfferingInstanceWithContext is an alternate form of the CreateOfferingInstance method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) CreateOfferingInstanceWithContext(ctx context.Context, createOfferingInstanceOptions *CreateOfferingInstanceOptions) (result *OfferingInstance, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createOfferingInstanceOptions, "createOfferingInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createOfferingInstanceOptions, "createOfferingInstanceOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/instances/offerings`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createOfferingInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "CreateOfferingInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createOfferingInstanceOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*createOfferingInstanceOptions.XAuthRefreshToken))
	}

	body := make(map[string]interface{})
	if createOfferingInstanceOptions.ID != nil {
		body["id"] = createOfferingInstanceOptions.ID
	}
	if createOfferingInstanceOptions.Rev != nil {
		body["_rev"] = createOfferingInstanceOptions.Rev
	}
	if createOfferingInstanceOptions.URL != nil {
		body["url"] = createOfferingInstanceOptions.URL
	}
	if createOfferingInstanceOptions.CRN != nil {
		body["crn"] = createOfferingInstanceOptions.CRN
	}
	if createOfferingInstanceOptions.Label != nil {
		body["label"] = createOfferingInstanceOptions.Label
	}
	if createOfferingInstanceOptions.CatalogID != nil {
		body["catalog_id"] = createOfferingInstanceOptions.CatalogID
	}
	if createOfferingInstanceOptions.OfferingID != nil {
		body["offering_id"] = createOfferingInstanceOptions.OfferingID
	}
	if createOfferingInstanceOptions.KindFormat != nil {
		body["kind_format"] = createOfferingInstanceOptions.KindFormat
	}
	if createOfferingInstanceOptions.Version != nil {
		body["version"] = createOfferingInstanceOptions.Version
	}
	if createOfferingInstanceOptions.ClusterID != nil {
		body["cluster_id"] = createOfferingInstanceOptions.ClusterID
	}
	if createOfferingInstanceOptions.ClusterRegion != nil {
		body["cluster_region"] = createOfferingInstanceOptions.ClusterRegion
	}
	if createOfferingInstanceOptions.ClusterNamespaces != nil {
		body["cluster_namespaces"] = createOfferingInstanceOptions.ClusterNamespaces
	}
	if createOfferingInstanceOptions.ClusterAllNamespaces != nil {
		body["cluster_all_namespaces"] = createOfferingInstanceOptions.ClusterAllNamespaces
	}
	if createOfferingInstanceOptions.SchematicsWorkspaceID != nil {
		body["schematics_workspace_id"] = createOfferingInstanceOptions.SchematicsWorkspaceID
	}
	if createOfferingInstanceOptions.ResourceGroupID != nil {
		body["resource_group_id"] = createOfferingInstanceOptions.ResourceGroupID
	}
	if createOfferingInstanceOptions.InstallPlan != nil {
		body["install_plan"] = createOfferingInstanceOptions.InstallPlan
	}
	if createOfferingInstanceOptions.Channel != nil {
		body["channel"] = createOfferingInstanceOptions.Channel
	}
	if createOfferingInstanceOptions.Metadata != nil {
		body["metadata"] = createOfferingInstanceOptions.Metadata
	}
	if createOfferingInstanceOptions.LastOperation != nil {
		body["last_operation"] = createOfferingInstanceOptions.LastOperation
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOfferingInstance)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOfferingInstance : Get Offering Instance
// Get the resource associated with an installed offering instance.
func (catalogManagement *CatalogManagementV1) GetOfferingInstance(getOfferingInstanceOptions *GetOfferingInstanceOptions) (result *OfferingInstance, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingInstanceWithContext(context.Background(), getOfferingInstanceOptions)
}

// GetOfferingInstanceWithContext is an alternate form of the GetOfferingInstance method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingInstanceWithContext(ctx context.Context, getOfferingInstanceOptions *GetOfferingInstanceOptions) (result *OfferingInstance, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingInstanceOptions, "getOfferingInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingInstanceOptions, "getOfferingInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_identifier": *getOfferingInstanceOptions.InstanceIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/instances/offerings/{instance_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOfferingInstance)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutOfferingInstance : Update Offering Instance
// Update an installed offering instance.
func (catalogManagement *CatalogManagementV1) PutOfferingInstance(putOfferingInstanceOptions *PutOfferingInstanceOptions) (result *OfferingInstance, response *core.DetailedResponse, err error) {
	return catalogManagement.PutOfferingInstanceWithContext(context.Background(), putOfferingInstanceOptions)
}

// PutOfferingInstanceWithContext is an alternate form of the PutOfferingInstance method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) PutOfferingInstanceWithContext(ctx context.Context, putOfferingInstanceOptions *PutOfferingInstanceOptions) (result *OfferingInstance, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putOfferingInstanceOptions, "putOfferingInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putOfferingInstanceOptions, "putOfferingInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_identifier": *putOfferingInstanceOptions.InstanceIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/instances/offerings/{instance_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putOfferingInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "PutOfferingInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if putOfferingInstanceOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*putOfferingInstanceOptions.XAuthRefreshToken))
	}

	body := make(map[string]interface{})
	if putOfferingInstanceOptions.ID != nil {
		body["id"] = putOfferingInstanceOptions.ID
	}
	if putOfferingInstanceOptions.Rev != nil {
		body["_rev"] = putOfferingInstanceOptions.Rev
	}
	if putOfferingInstanceOptions.URL != nil {
		body["url"] = putOfferingInstanceOptions.URL
	}
	if putOfferingInstanceOptions.CRN != nil {
		body["crn"] = putOfferingInstanceOptions.CRN
	}
	if putOfferingInstanceOptions.Label != nil {
		body["label"] = putOfferingInstanceOptions.Label
	}
	if putOfferingInstanceOptions.CatalogID != nil {
		body["catalog_id"] = putOfferingInstanceOptions.CatalogID
	}
	if putOfferingInstanceOptions.OfferingID != nil {
		body["offering_id"] = putOfferingInstanceOptions.OfferingID
	}
	if putOfferingInstanceOptions.KindFormat != nil {
		body["kind_format"] = putOfferingInstanceOptions.KindFormat
	}
	if putOfferingInstanceOptions.Version != nil {
		body["version"] = putOfferingInstanceOptions.Version
	}
	if putOfferingInstanceOptions.ClusterID != nil {
		body["cluster_id"] = putOfferingInstanceOptions.ClusterID
	}
	if putOfferingInstanceOptions.ClusterRegion != nil {
		body["cluster_region"] = putOfferingInstanceOptions.ClusterRegion
	}
	if putOfferingInstanceOptions.ClusterNamespaces != nil {
		body["cluster_namespaces"] = putOfferingInstanceOptions.ClusterNamespaces
	}
	if putOfferingInstanceOptions.ClusterAllNamespaces != nil {
		body["cluster_all_namespaces"] = putOfferingInstanceOptions.ClusterAllNamespaces
	}
	if putOfferingInstanceOptions.SchematicsWorkspaceID != nil {
		body["schematics_workspace_id"] = putOfferingInstanceOptions.SchematicsWorkspaceID
	}
	if putOfferingInstanceOptions.ResourceGroupID != nil {
		body["resource_group_id"] = putOfferingInstanceOptions.ResourceGroupID
	}
	if putOfferingInstanceOptions.InstallPlan != nil {
		body["install_plan"] = putOfferingInstanceOptions.InstallPlan
	}
	if putOfferingInstanceOptions.Channel != nil {
		body["channel"] = putOfferingInstanceOptions.Channel
	}
	if putOfferingInstanceOptions.Metadata != nil {
		body["metadata"] = putOfferingInstanceOptions.Metadata
	}
	if putOfferingInstanceOptions.LastOperation != nil {
		body["last_operation"] = putOfferingInstanceOptions.LastOperation
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOfferingInstance)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteOfferingInstance : Delete a version instance
// Delete and instance deployed out of a product version.
func (catalogManagement *CatalogManagementV1) DeleteOfferingInstance(deleteOfferingInstanceOptions *DeleteOfferingInstanceOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteOfferingInstanceWithContext(context.Background(), deleteOfferingInstanceOptions)
}

// DeleteOfferingInstanceWithContext is an alternate form of the DeleteOfferingInstance method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteOfferingInstanceWithContext(ctx context.Context, deleteOfferingInstanceOptions *DeleteOfferingInstanceOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteOfferingInstanceOptions, "deleteOfferingInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteOfferingInstanceOptions, "deleteOfferingInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_identifier": *deleteOfferingInstanceOptions.InstanceIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/instances/offerings/{instance_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteOfferingInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteOfferingInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteOfferingInstanceOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*deleteOfferingInstanceOptions.XAuthRefreshToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, nil)

	return
}

// AccessListBulkResponse : Access List Add/Remove result.
type AccessListBulkResponse struct {
	// in the case of error on an account add/remove - account: error.
	Errors map[string]string `json:"errors,omitempty"`
}

// UnmarshalAccessListBulkResponse unmarshals an instance of AccessListBulkResponse from the specified map of raw messages.
func UnmarshalAccessListBulkResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessListBulkResponse)
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Account : Account information.
type Account struct {
	// Account identification.
	ID *string `json:"id,omitempty"`

	// Hide the public catalog in this account.
	HideIBMCloudCatalog *bool `json:"hide_IBM_cloud_catalog,omitempty"`

	// Filters for account and catalog filters.
	AccountFilters *Filters `json:"account_filters,omitempty"`
}

// UnmarshalAccount unmarshals an instance of Account from the specified map of raw messages.
func UnmarshalAccount(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Account)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hide_IBM_cloud_catalog", &obj.HideIBMCloudCatalog)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account_filters", &obj.AccountFilters, UnmarshalFilters)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountPublishObjectOptions : The AccountPublishObject options.
type AccountPublishObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewAccountPublishObjectOptions : Instantiate AccountPublishObjectOptions
func (*CatalogManagementV1) NewAccountPublishObjectOptions(catalogIdentifier string, objectIdentifier string) *AccountPublishObjectOptions {
	return &AccountPublishObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *AccountPublishObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *AccountPublishObjectOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *AccountPublishObjectOptions) SetObjectIdentifier(objectIdentifier string) *AccountPublishObjectOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *AccountPublishObjectOptions) SetHeaders(param map[string]string) *AccountPublishObjectOptions {
	options.Headers = param
	return options
}

// AccountPublishVersionOptions : The AccountPublishVersion options.
type AccountPublishVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewAccountPublishVersionOptions : Instantiate AccountPublishVersionOptions
func (*CatalogManagementV1) NewAccountPublishVersionOptions(versionLocID string) *AccountPublishVersionOptions {
	return &AccountPublishVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *AccountPublishVersionOptions) SetVersionLocID(versionLocID string) *AccountPublishVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *AccountPublishVersionOptions) SetHeaders(param map[string]string) *AccountPublishVersionOptions {
	options.Headers = param
	return options
}

// AccumulatedFilters : The accumulated filters for an account. This will return the account filters plus a filter for each catalog the user
// has access to.
type AccumulatedFilters struct {
	// Filters for accounts (at this time this will always be just one item array).
	AccountFilters []Filters `json:"account_filters,omitempty"`

	// The filters for all of the accessible catalogs.
	CatalogFilters []AccumulatedFiltersCatalogFiltersItem `json:"catalog_filters,omitempty"`
}

// UnmarshalAccumulatedFilters unmarshals an instance of AccumulatedFilters from the specified map of raw messages.
func UnmarshalAccumulatedFilters(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccumulatedFilters)
	err = core.UnmarshalModel(m, "account_filters", &obj.AccountFilters, UnmarshalFilters)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "catalog_filters", &obj.CatalogFilters, UnmarshalAccumulatedFiltersCatalogFiltersItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccumulatedFiltersCatalogFiltersItem : AccumulatedFiltersCatalogFiltersItem struct
type AccumulatedFiltersCatalogFiltersItem struct {
	// Filters for catalog.
	Catalog *AccumulatedFiltersCatalogFiltersItemCatalog `json:"catalog,omitempty"`

	// Filters for account and catalog filters.
	Filters *Filters `json:"filters,omitempty"`
}

// UnmarshalAccumulatedFiltersCatalogFiltersItem unmarshals an instance of AccumulatedFiltersCatalogFiltersItem from the specified map of raw messages.
func UnmarshalAccumulatedFiltersCatalogFiltersItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccumulatedFiltersCatalogFiltersItem)
	err = core.UnmarshalModel(m, "catalog", &obj.Catalog, UnmarshalAccumulatedFiltersCatalogFiltersItemCatalog)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "filters", &obj.Filters, UnmarshalFilters)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccumulatedFiltersCatalogFiltersItemCatalog : Filters for catalog.
type AccumulatedFiltersCatalogFiltersItemCatalog struct {
	// The ID of the catalog.
	ID *string `json:"id,omitempty"`

	// The name of the catalog.
	Name *string `json:"name,omitempty"`
}

// UnmarshalAccumulatedFiltersCatalogFiltersItemCatalog unmarshals an instance of AccumulatedFiltersCatalogFiltersItemCatalog from the specified map of raw messages.
func UnmarshalAccumulatedFiltersCatalogFiltersItemCatalog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccumulatedFiltersCatalogFiltersItemCatalog)
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

// AddObjectAccessListOptions : The AddObjectAccessList options.
type AddObjectAccessListOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// A list of accounts to add.
	Accounts []string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewAddObjectAccessListOptions : Instantiate AddObjectAccessListOptions
func (*CatalogManagementV1) NewAddObjectAccessListOptions(catalogIdentifier string, objectIdentifier string, accounts []string) *AddObjectAccessListOptions {
	return &AddObjectAccessListOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		Accounts: accounts,
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *AddObjectAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *AddObjectAccessListOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *AddObjectAccessListOptions) SetObjectIdentifier(objectIdentifier string) *AddObjectAccessListOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetAccounts : Allow user to set Accounts
func (options *AddObjectAccessListOptions) SetAccounts(accounts []string) *AddObjectAccessListOptions {
	options.Accounts = accounts
	return options
}

// SetHeaders : Allow user to set Headers
func (options *AddObjectAccessListOptions) SetHeaders(param map[string]string) *AddObjectAccessListOptions {
	options.Headers = param
	return options
}

// ApprovalResult : Result of approval.
type ApprovalResult struct {
	// Allowed to request to publish.
	AllowRequest *bool `json:"allow_request,omitempty"`

	// Visible to IBM.
	IBM *bool `json:"ibm,omitempty"`

	// Visible to everyone.
	Public *bool `json:"public,omitempty"`

	// Denotes whether approval has changed.
	Changed *bool `json:"changed,omitempty"`
}

// UnmarshalApprovalResult unmarshals an instance of ApprovalResult from the specified map of raw messages.
func UnmarshalApprovalResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApprovalResult)
	err = core.UnmarshalPrimitive(m, "allow_request", &obj.AllowRequest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm", &obj.IBM)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "public", &obj.Public)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "changed", &obj.Changed)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AuditLog : A collection of audit records.
type AuditLog struct {
	// A list of audit records.
	List []AuditRecord `json:"list,omitempty"`
}

// UnmarshalAuditLog unmarshals an instance of AuditLog from the specified map of raw messages.
func UnmarshalAuditLog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AuditLog)
	err = core.UnmarshalModel(m, "list", &obj.List, UnmarshalAuditRecord)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AuditRecord : An audit record which describes a change made to a catalog or associated resource.
type AuditRecord struct {
	// The identifier of the audit record.
	ID *string `json:"id,omitempty"`

	// The time at which the change was made.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The type of change described by the audit record.
	ChangeType *string `json:"change_type,omitempty"`

	// The resource type associated with the change.
	TargetType *string `json:"target_type,omitempty"`

	// The identifier of the resource that was changed.
	TargetID *string `json:"target_id,omitempty"`

	// The email address of the user that made the change.
	WhoDelegateEmail *string `json:"who_delegate_email,omitempty"`

	// A message which describes the change.
	Message *string `json:"message,omitempty"`
}

// UnmarshalAuditRecord unmarshals an instance of AuditRecord from the specified map of raw messages.
func UnmarshalAuditRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AuditRecord)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "change_type", &obj.ChangeType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_type", &obj.TargetType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_id", &obj.TargetID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "who_delegate_email", &obj.WhoDelegateEmail)
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

// Catalog : Catalog information.
type Catalog struct {
	// Unique ID.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// Display Name in the requested language.
	Label *string `json:"label,omitempty"`

	// Description in the requested language.
	ShortDescription *string `json:"short_description,omitempty"`

	// URL for an icon associated with this catalog.
	CatalogIconURL *string `json:"catalog_icon_url,omitempty"`

	// List of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// The url for this specific catalog.
	URL *string `json:"url,omitempty"`

	// CRN associated with the catalog.
	CRN *string `json:"crn,omitempty"`

	// URL path to offerings.
	OfferingsURL *string `json:"offerings_url,omitempty"`

	// List of features associated with this catalog.
	Features []Feature `json:"features,omitempty"`

	// Denotes whether a catalog is disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// The date-time this catalog was created.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The date-time this catalog was last updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Resource group id the catalog is owned by.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// Account that owns catalog.
	OwningAccount *string `json:"owning_account,omitempty"`

	// Filters for account and catalog filters.
	CatalogFilters *Filters `json:"catalog_filters,omitempty"`

	// Feature information.
	SyndicationSettings *SyndicationResource `json:"syndication_settings,omitempty"`

	// Kind of catalog. Supported kinds are offering and vpe.
	Kind *string `json:"kind,omitempty"`
}

// UnmarshalCatalog unmarshals an instance of Catalog from the specified map of raw messages.
func UnmarshalCatalog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Catalog)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "short_description", &obj.ShortDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_icon_url", &obj.CatalogIconURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offerings_url", &obj.OfferingsURL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeature)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "owning_account", &obj.OwningAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "catalog_filters", &obj.CatalogFilters, UnmarshalFilters)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "syndication_settings", &obj.SyndicationSettings, UnmarshalSyndicationResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CatalogObject : object information.
type CatalogObject struct {
	// unique id.
	ID *string `json:"id,omitempty"`

	// The programmatic name of this offering.
	Name *string `json:"name,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// The crn for this specific object.
	CRN *string `json:"crn,omitempty"`

	// The url for this specific object.
	URL *string `json:"url,omitempty"`

	// The parent for this specific object.
	ParentID *string `json:"parent_id,omitempty"`

	// Translated display name in the requested language.
	LabelI18n *string `json:"label_i18n,omitempty"`

	// Display name in the requested language.
	Label *string `json:"label,omitempty"`

	// List of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// The date and time this catalog was created.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The date and time this catalog was last updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Short description in the requested language.
	ShortDescription *string `json:"short_description,omitempty"`

	// Short description translation.
	ShortDescriptionI18n *string `json:"short_description_i18n,omitempty"`

	// Kind of object.
	Kind *string `json:"kind,omitempty"`

	// Publish information.
	Publish *PublishObject `json:"publish,omitempty"`

	// Offering state.
	State *State `json:"state,omitempty"`

	// The id of the catalog containing this offering.
	CatalogID *string `json:"catalog_id,omitempty"`

	// The name of the catalog.
	CatalogName *string `json:"catalog_name,omitempty"`

	// Map of data values for this object.
	Data map[string]interface{} `json:"data,omitempty"`
}

// UnmarshalCatalogObject unmarshals an instance of CatalogObject from the specified map of raw messages.
func UnmarshalCatalogObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogObject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parent_id", &obj.ParentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label_i18n", &obj.LabelI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "short_description", &obj.ShortDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "short_description_i18n", &obj.ShortDescriptionI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "publish", &obj.Publish, UnmarshalPublishObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "state", &obj.State, UnmarshalState)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_name", &obj.CatalogName)
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

// CatalogSearchResult : Paginated catalog search result.
type CatalogSearchResult struct {
	// The overall total number of resources in the search result set.
	TotalCount *int64 `json:"total_count,omitempty"`

	// Resulting objects.
	Resources []Catalog `json:"resources,omitempty"`
}

// UnmarshalCatalogSearchResult unmarshals an instance of CatalogSearchResult from the specified map of raw messages.
func UnmarshalCatalogSearchResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogSearchResult)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalCatalog)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CategoryFilter : Filter on a category. The filter will match against the values of the given category with include or exclude.
type CategoryFilter struct {
	// -> true - This is an include filter, false - this is an exclude filter.
	Include *bool `json:"include,omitempty"`

	// Offering filter terms.
	Filter *FilterTerms `json:"filter,omitempty"`
}

// UnmarshalCategoryFilter unmarshals an instance of CategoryFilter from the specified map of raw messages.
func UnmarshalCategoryFilter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CategoryFilter)
	err = core.UnmarshalPrimitive(m, "include", &obj.Include)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "filter", &obj.Filter, UnmarshalFilterTerms)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ClusterInfo : Cluster information.
type ClusterInfo struct {
	// Resource Group ID.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// Resource Group name.
	ResourceGroupName *string `json:"resource_group_name,omitempty"`

	// Cluster ID.
	ID *string `json:"id,omitempty"`

	// Cluster name.
	Name *string `json:"name,omitempty"`

	// Cluster region.
	Region *string `json:"region,omitempty"`
}

// UnmarshalClusterInfo unmarshals an instance of ClusterInfo from the specified map of raw messages.
func UnmarshalClusterInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ClusterInfo)
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_name", &obj.ResourceGroupName)
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
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CommitVersionOptions : The CommitVersion options.
type CommitVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCommitVersionOptions : Instantiate CommitVersionOptions
func (*CatalogManagementV1) NewCommitVersionOptions(versionLocID string) *CommitVersionOptions {
	return &CommitVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *CommitVersionOptions) SetVersionLocID(versionLocID string) *CommitVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CommitVersionOptions) SetHeaders(param map[string]string) *CommitVersionOptions {
	options.Headers = param
	return options
}

// Configuration : Configuration description.
type Configuration struct {
	// Configuration key.
	Key *string `json:"key,omitempty"`

	// Value type (string, boolean, int).
	Type *string `json:"type,omitempty"`

	// The default value.  To use a secret when the type is password, specify a JSON encoded value of
	// $ref:#/components/schemas/SecretInstance, prefixed with `cmsm_v1:`.
	DefaultValue interface{} `json:"default_value,omitempty"`

	// Constraint associated with value, e.g., for string type - regx:[a-z].
	ValueConstraint *string `json:"value_constraint,omitempty"`

	// Key description.
	Description *string `json:"description,omitempty"`

	// Is key required to install.
	Required *bool `json:"required,omitempty"`

	// List of options of type.
	Options []interface{} `json:"options,omitempty"`

	// Hide values.
	Hidden *bool `json:"hidden,omitempty"`
}

// UnmarshalConfiguration unmarshals an instance of Configuration from the specified map of raw messages.
func UnmarshalConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Configuration)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_value", &obj.DefaultValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value_constraint", &obj.ValueConstraint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "required", &obj.Required)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "options", &obj.Options)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CopyVersionOptions : The CopyVersion options.
type CopyVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Tags array.
	Tags []string

	// Target kinds.  Current valid values are 'iks', 'roks', 'vcenter', and 'terraform'.
	TargetKinds []string

	// byte array representing the content to be imported.  Only supported for OVA images at this time.
	Content *[]byte

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCopyVersionOptions : Instantiate CopyVersionOptions
func (*CatalogManagementV1) NewCopyVersionOptions(versionLocID string) *CopyVersionOptions {
	return &CopyVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *CopyVersionOptions) SetVersionLocID(versionLocID string) *CopyVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetTags : Allow user to set Tags
func (options *CopyVersionOptions) SetTags(tags []string) *CopyVersionOptions {
	options.Tags = tags
	return options
}

// SetTargetKinds : Allow user to set TargetKinds
func (options *CopyVersionOptions) SetTargetKinds(targetKinds []string) *CopyVersionOptions {
	options.TargetKinds = targetKinds
	return options
}

// SetContent : Allow user to set Content
func (options *CopyVersionOptions) SetContent(content []byte) *CopyVersionOptions {
	options.Content = &content
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CopyVersionOptions) SetHeaders(param map[string]string) *CopyVersionOptions {
	options.Headers = param
	return options
}

// CreateCatalogOptions : The CreateCatalog options.
type CreateCatalogOptions struct {
	// Unique ID.
	ID *string

	// Cloudant revision.
	Rev *string

	// Display Name in the requested language.
	Label *string

	// Description in the requested language.
	ShortDescription *string

	// URL for an icon associated with this catalog.
	CatalogIconURL *string

	// List of tags associated with this catalog.
	Tags []string

	// List of features associated with this catalog.
	Features []Feature

	// Denotes whether a catalog is disabled.
	Disabled *bool

	// Resource group id the catalog is owned by.
	ResourceGroupID *string

	// Account that owns catalog.
	OwningAccount *string

	// Filters for account and catalog filters.
	CatalogFilters *Filters

	// Feature information.
	SyndicationSettings *SyndicationResource

	// Kind of catalog. Supported kinds are offering and vpe.
	Kind *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateCatalogOptions : Instantiate CreateCatalogOptions
func (*CatalogManagementV1) NewCreateCatalogOptions() *CreateCatalogOptions {
	return &CreateCatalogOptions{}
}

// SetID : Allow user to set ID
func (options *CreateCatalogOptions) SetID(id string) *CreateCatalogOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetRev : Allow user to set Rev
func (options *CreateCatalogOptions) SetRev(rev string) *CreateCatalogOptions {
	options.Rev = core.StringPtr(rev)
	return options
}

// SetLabel : Allow user to set Label
func (options *CreateCatalogOptions) SetLabel(label string) *CreateCatalogOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetShortDescription : Allow user to set ShortDescription
func (options *CreateCatalogOptions) SetShortDescription(shortDescription string) *CreateCatalogOptions {
	options.ShortDescription = core.StringPtr(shortDescription)
	return options
}

// SetCatalogIconURL : Allow user to set CatalogIconURL
func (options *CreateCatalogOptions) SetCatalogIconURL(catalogIconURL string) *CreateCatalogOptions {
	options.CatalogIconURL = core.StringPtr(catalogIconURL)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateCatalogOptions) SetTags(tags []string) *CreateCatalogOptions {
	options.Tags = tags
	return options
}

// SetFeatures : Allow user to set Features
func (options *CreateCatalogOptions) SetFeatures(features []Feature) *CreateCatalogOptions {
	options.Features = features
	return options
}

// SetDisabled : Allow user to set Disabled
func (options *CreateCatalogOptions) SetDisabled(disabled bool) *CreateCatalogOptions {
	options.Disabled = core.BoolPtr(disabled)
	return options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (options *CreateCatalogOptions) SetResourceGroupID(resourceGroupID string) *CreateCatalogOptions {
	options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return options
}

// SetOwningAccount : Allow user to set OwningAccount
func (options *CreateCatalogOptions) SetOwningAccount(owningAccount string) *CreateCatalogOptions {
	options.OwningAccount = core.StringPtr(owningAccount)
	return options
}

// SetCatalogFilters : Allow user to set CatalogFilters
func (options *CreateCatalogOptions) SetCatalogFilters(catalogFilters *Filters) *CreateCatalogOptions {
	options.CatalogFilters = catalogFilters
	return options
}

// SetSyndicationSettings : Allow user to set SyndicationSettings
func (options *CreateCatalogOptions) SetSyndicationSettings(syndicationSettings *SyndicationResource) *CreateCatalogOptions {
	options.SyndicationSettings = syndicationSettings
	return options
}

// SetKind : Allow user to set Kind
func (options *CreateCatalogOptions) SetKind(kind string) *CreateCatalogOptions {
	options.Kind = core.StringPtr(kind)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCatalogOptions) SetHeaders(param map[string]string) *CreateCatalogOptions {
	options.Headers = param
	return options
}

// CreateObjectAccessOptions : The CreateObjectAccess options.
type CreateObjectAccessOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Account identifier.
	AccountIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateObjectAccessOptions : Instantiate CreateObjectAccessOptions
func (*CatalogManagementV1) NewCreateObjectAccessOptions(catalogIdentifier string, objectIdentifier string, accountIdentifier string) *CreateObjectAccessOptions {
	return &CreateObjectAccessOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		AccountIdentifier: core.StringPtr(accountIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *CreateObjectAccessOptions) SetCatalogIdentifier(catalogIdentifier string) *CreateObjectAccessOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *CreateObjectAccessOptions) SetObjectIdentifier(objectIdentifier string) *CreateObjectAccessOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetAccountIdentifier : Allow user to set AccountIdentifier
func (options *CreateObjectAccessOptions) SetAccountIdentifier(accountIdentifier string) *CreateObjectAccessOptions {
	options.AccountIdentifier = core.StringPtr(accountIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateObjectAccessOptions) SetHeaders(param map[string]string) *CreateObjectAccessOptions {
	options.Headers = param
	return options
}

// CreateObjectOptions : The CreateObject options.
type CreateObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// unique id.
	ID *string

	// The programmatic name of this offering.
	Name *string

	// Cloudant revision.
	Rev *string

	// The crn for this specific object.
	CRN *string

	// The url for this specific object.
	URL *string

	// The parent for this specific object.
	ParentID *string

	// Translated display name in the requested language.
	LabelI18n *string

	// Display name in the requested language.
	Label *string

	// List of tags associated with this catalog.
	Tags []string

	// The date and time this catalog was created.
	Created *strfmt.DateTime

	// The date and time this catalog was last updated.
	Updated *strfmt.DateTime

	// Short description in the requested language.
	ShortDescription *string

	// Short description translation.
	ShortDescriptionI18n *string

	// Kind of object.
	Kind *string

	// Publish information.
	Publish *PublishObject

	// Offering state.
	State *State

	// The id of the catalog containing this offering.
	CatalogID *string

	// The name of the catalog.
	CatalogName *string

	// Map of data values for this object.
	Data map[string]interface{}

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateObjectOptions : Instantiate CreateObjectOptions
func (*CatalogManagementV1) NewCreateObjectOptions(catalogIdentifier string) *CreateObjectOptions {
	return &CreateObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *CreateObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *CreateObjectOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetID : Allow user to set ID
func (options *CreateObjectOptions) SetID(id string) *CreateObjectOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetName : Allow user to set Name
func (options *CreateObjectOptions) SetName(name string) *CreateObjectOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetRev : Allow user to set Rev
func (options *CreateObjectOptions) SetRev(rev string) *CreateObjectOptions {
	options.Rev = core.StringPtr(rev)
	return options
}

// SetCRN : Allow user to set CRN
func (options *CreateObjectOptions) SetCRN(crn string) *CreateObjectOptions {
	options.CRN = core.StringPtr(crn)
	return options
}

// SetURL : Allow user to set URL
func (options *CreateObjectOptions) SetURL(url string) *CreateObjectOptions {
	options.URL = core.StringPtr(url)
	return options
}

// SetParentID : Allow user to set ParentID
func (options *CreateObjectOptions) SetParentID(parentID string) *CreateObjectOptions {
	options.ParentID = core.StringPtr(parentID)
	return options
}

// SetLabelI18n : Allow user to set LabelI18n
func (options *CreateObjectOptions) SetLabelI18n(labelI18n string) *CreateObjectOptions {
	options.LabelI18n = core.StringPtr(labelI18n)
	return options
}

// SetLabel : Allow user to set Label
func (options *CreateObjectOptions) SetLabel(label string) *CreateObjectOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateObjectOptions) SetTags(tags []string) *CreateObjectOptions {
	options.Tags = tags
	return options
}

// SetCreated : Allow user to set Created
func (options *CreateObjectOptions) SetCreated(created *strfmt.DateTime) *CreateObjectOptions {
	options.Created = created
	return options
}

// SetUpdated : Allow user to set Updated
func (options *CreateObjectOptions) SetUpdated(updated *strfmt.DateTime) *CreateObjectOptions {
	options.Updated = updated
	return options
}

// SetShortDescription : Allow user to set ShortDescription
func (options *CreateObjectOptions) SetShortDescription(shortDescription string) *CreateObjectOptions {
	options.ShortDescription = core.StringPtr(shortDescription)
	return options
}

// SetShortDescriptionI18n : Allow user to set ShortDescriptionI18n
func (options *CreateObjectOptions) SetShortDescriptionI18n(shortDescriptionI18n string) *CreateObjectOptions {
	options.ShortDescriptionI18n = core.StringPtr(shortDescriptionI18n)
	return options
}

// SetKind : Allow user to set Kind
func (options *CreateObjectOptions) SetKind(kind string) *CreateObjectOptions {
	options.Kind = core.StringPtr(kind)
	return options
}

// SetPublish : Allow user to set Publish
func (options *CreateObjectOptions) SetPublish(publish *PublishObject) *CreateObjectOptions {
	options.Publish = publish
	return options
}

// SetState : Allow user to set State
func (options *CreateObjectOptions) SetState(state *State) *CreateObjectOptions {
	options.State = state
	return options
}

// SetCatalogID : Allow user to set CatalogID
func (options *CreateObjectOptions) SetCatalogID(catalogID string) *CreateObjectOptions {
	options.CatalogID = core.StringPtr(catalogID)
	return options
}

// SetCatalogName : Allow user to set CatalogName
func (options *CreateObjectOptions) SetCatalogName(catalogName string) *CreateObjectOptions {
	options.CatalogName = core.StringPtr(catalogName)
	return options
}

// SetData : Allow user to set Data
func (options *CreateObjectOptions) SetData(data map[string]interface{}) *CreateObjectOptions {
	options.Data = data
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateObjectOptions) SetHeaders(param map[string]string) *CreateObjectOptions {
	options.Headers = param
	return options
}

// CreateOfferingInstanceOptions : The CreateOfferingInstance options.
type CreateOfferingInstanceOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// provisioned instance ID (part of the CRN).
	ID *string

	// Cloudant revision.
	Rev *string

	// url reference to this object.
	URL *string

	// platform CRN for this instance.
	CRN *string

	// the label for this instance.
	Label *string

	// Catalog ID this instance was created from.
	CatalogID *string

	// Offering ID this instance was created from.
	OfferingID *string

	// the format this instance has (helm, operator, ova...).
	KindFormat *string

	// The version this instance was installed from (not version id).
	Version *string

	// Cluster ID.
	ClusterID *string

	// Cluster region (e.g., us-south).
	ClusterRegion *string

	// List of target namespaces to install into.
	ClusterNamespaces []string

	// designate to install into all namespaces.
	ClusterAllNamespaces *bool

	// Id of the schematics workspace, for offering instances provisioned through schematics.
	SchematicsWorkspaceID *string

	// Id of the resource group to provision the offering instance into.
	ResourceGroupID *string

	// Type of install plan (also known as approval strategy) for operator subscriptions. Can be either automatic, which
	// automatically upgrades operators to the latest in a channel, or manual, which requires approval on the cluster.
	InstallPlan *string

	// Channel to pin the operator subscription to.
	Channel *string

	// Map of metadata values for this offering instance.
	Metadata map[string]interface{}

	// the last operation performed and status.
	LastOperation *OfferingInstanceLastOperation

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateOfferingInstanceOptions : Instantiate CreateOfferingInstanceOptions
func (*CatalogManagementV1) NewCreateOfferingInstanceOptions(xAuthRefreshToken string) *CreateOfferingInstanceOptions {
	return &CreateOfferingInstanceOptions{
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *CreateOfferingInstanceOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *CreateOfferingInstanceOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetID : Allow user to set ID
func (options *CreateOfferingInstanceOptions) SetID(id string) *CreateOfferingInstanceOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetRev : Allow user to set Rev
func (options *CreateOfferingInstanceOptions) SetRev(rev string) *CreateOfferingInstanceOptions {
	options.Rev = core.StringPtr(rev)
	return options
}

// SetURL : Allow user to set URL
func (options *CreateOfferingInstanceOptions) SetURL(url string) *CreateOfferingInstanceOptions {
	options.URL = core.StringPtr(url)
	return options
}

// SetCRN : Allow user to set CRN
func (options *CreateOfferingInstanceOptions) SetCRN(crn string) *CreateOfferingInstanceOptions {
	options.CRN = core.StringPtr(crn)
	return options
}

// SetLabel : Allow user to set Label
func (options *CreateOfferingInstanceOptions) SetLabel(label string) *CreateOfferingInstanceOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetCatalogID : Allow user to set CatalogID
func (options *CreateOfferingInstanceOptions) SetCatalogID(catalogID string) *CreateOfferingInstanceOptions {
	options.CatalogID = core.StringPtr(catalogID)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *CreateOfferingInstanceOptions) SetOfferingID(offeringID string) *CreateOfferingInstanceOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetKindFormat : Allow user to set KindFormat
func (options *CreateOfferingInstanceOptions) SetKindFormat(kindFormat string) *CreateOfferingInstanceOptions {
	options.KindFormat = core.StringPtr(kindFormat)
	return options
}

// SetVersion : Allow user to set Version
func (options *CreateOfferingInstanceOptions) SetVersion(version string) *CreateOfferingInstanceOptions {
	options.Version = core.StringPtr(version)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *CreateOfferingInstanceOptions) SetClusterID(clusterID string) *CreateOfferingInstanceOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetClusterRegion : Allow user to set ClusterRegion
func (options *CreateOfferingInstanceOptions) SetClusterRegion(clusterRegion string) *CreateOfferingInstanceOptions {
	options.ClusterRegion = core.StringPtr(clusterRegion)
	return options
}

// SetClusterNamespaces : Allow user to set ClusterNamespaces
func (options *CreateOfferingInstanceOptions) SetClusterNamespaces(clusterNamespaces []string) *CreateOfferingInstanceOptions {
	options.ClusterNamespaces = clusterNamespaces
	return options
}

// SetClusterAllNamespaces : Allow user to set ClusterAllNamespaces
func (options *CreateOfferingInstanceOptions) SetClusterAllNamespaces(clusterAllNamespaces bool) *CreateOfferingInstanceOptions {
	options.ClusterAllNamespaces = core.BoolPtr(clusterAllNamespaces)
	return options
}

// SetSchematicsWorkspaceID : Allow user to set SchematicsWorkspaceID
func (options *CreateOfferingInstanceOptions) SetSchematicsWorkspaceID(schematicsWorkspaceID string) *CreateOfferingInstanceOptions {
	options.SchematicsWorkspaceID = core.StringPtr(schematicsWorkspaceID)
	return options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (options *CreateOfferingInstanceOptions) SetResourceGroupID(resourceGroupID string) *CreateOfferingInstanceOptions {
	options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return options
}

// SetInstallPlan : Allow user to set InstallPlan
func (options *CreateOfferingInstanceOptions) SetInstallPlan(installPlan string) *CreateOfferingInstanceOptions {
	options.InstallPlan = core.StringPtr(installPlan)
	return options
}

// SetChannel : Allow user to set Channel
func (options *CreateOfferingInstanceOptions) SetChannel(channel string) *CreateOfferingInstanceOptions {
	options.Channel = core.StringPtr(channel)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *CreateOfferingInstanceOptions) SetMetadata(metadata map[string]interface{}) *CreateOfferingInstanceOptions {
	options.Metadata = metadata
	return options
}

// SetLastOperation : Allow user to set LastOperation
func (options *CreateOfferingInstanceOptions) SetLastOperation(lastOperation *OfferingInstanceLastOperation) *CreateOfferingInstanceOptions {
	options.LastOperation = lastOperation
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateOfferingInstanceOptions) SetHeaders(param map[string]string) *CreateOfferingInstanceOptions {
	options.Headers = param
	return options
}

// CreateOfferingOptions : The CreateOffering options.
type CreateOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// unique id.
	ID *string

	// Cloudant revision.
	Rev *string

	// The url for this specific offering.
	URL *string

	// The crn for this specific offering.
	CRN *string

	// Display Name in the requested language.
	Label *string

	// The programmatic name of this offering.
	Name *string

	// URL for an icon associated with this offering.
	OfferingIconURL *string

	// URL for an additional docs with this offering.
	OfferingDocsURL *string

	// URL to be displayed in the Consumption UI for getting support on this offering.
	OfferingSupportURL *string

	// List of tags associated with this catalog.
	Tags []string

	// List of keywords associated with offering, typically used to search for it.
	Keywords []string

	// Repository info for offerings.
	Rating *Rating

	// The date and time this catalog was created.
	Created *strfmt.DateTime

	// The date and time this catalog was last updated.
	Updated *strfmt.DateTime

	// Short description in the requested language.
	ShortDescription *string

	// Long description in the requested language.
	LongDescription *string

	// list of features associated with this offering.
	Features []Feature

	// Array of kind.
	Kinds []Kind

	// Is it permitted to request publishing to IBM or Public.
	PermitRequestIBMPublicPublish *bool

	// Indicates if this offering has been approved for use by all IBMers.
	IBMPublishApproved *bool

	// Indicates if this offering has been approved for use by all IBM Cloud users.
	PublicPublishApproved *bool

	// The original offering CRN that this publish entry came from.
	PublicOriginalCRN *string

	// The crn of the public catalog entry of this offering.
	PublishPublicCRN *string

	// The portal's approval record ID.
	PortalApprovalRecord *string

	// The portal UI URL.
	PortalUIURL *string

	// The id of the catalog containing this offering.
	CatalogID *string

	// The name of the catalog.
	CatalogName *string

	// Map of metadata values for this offering.
	Metadata map[string]interface{}

	// A disclaimer for this offering.
	Disclaimer *string

	// Determine if this offering should be displayed in the Consumption UI.
	Hidden *bool

	// Provider of this offering.
	Provider *string

	// Repository info for offerings.
	RepoInfo *RepoInfo

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateOfferingOptions : Instantiate CreateOfferingOptions
func (*CatalogManagementV1) NewCreateOfferingOptions(catalogIdentifier string) *CreateOfferingOptions {
	return &CreateOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *CreateOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *CreateOfferingOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetID : Allow user to set ID
func (options *CreateOfferingOptions) SetID(id string) *CreateOfferingOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetRev : Allow user to set Rev
func (options *CreateOfferingOptions) SetRev(rev string) *CreateOfferingOptions {
	options.Rev = core.StringPtr(rev)
	return options
}

// SetURL : Allow user to set URL
func (options *CreateOfferingOptions) SetURL(url string) *CreateOfferingOptions {
	options.URL = core.StringPtr(url)
	return options
}

// SetCRN : Allow user to set CRN
func (options *CreateOfferingOptions) SetCRN(crn string) *CreateOfferingOptions {
	options.CRN = core.StringPtr(crn)
	return options
}

// SetLabel : Allow user to set Label
func (options *CreateOfferingOptions) SetLabel(label string) *CreateOfferingOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetName : Allow user to set Name
func (options *CreateOfferingOptions) SetName(name string) *CreateOfferingOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetOfferingIconURL : Allow user to set OfferingIconURL
func (options *CreateOfferingOptions) SetOfferingIconURL(offeringIconURL string) *CreateOfferingOptions {
	options.OfferingIconURL = core.StringPtr(offeringIconURL)
	return options
}

// SetOfferingDocsURL : Allow user to set OfferingDocsURL
func (options *CreateOfferingOptions) SetOfferingDocsURL(offeringDocsURL string) *CreateOfferingOptions {
	options.OfferingDocsURL = core.StringPtr(offeringDocsURL)
	return options
}

// SetOfferingSupportURL : Allow user to set OfferingSupportURL
func (options *CreateOfferingOptions) SetOfferingSupportURL(offeringSupportURL string) *CreateOfferingOptions {
	options.OfferingSupportURL = core.StringPtr(offeringSupportURL)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateOfferingOptions) SetTags(tags []string) *CreateOfferingOptions {
	options.Tags = tags
	return options
}

// SetKeywords : Allow user to set Keywords
func (options *CreateOfferingOptions) SetKeywords(keywords []string) *CreateOfferingOptions {
	options.Keywords = keywords
	return options
}

// SetRating : Allow user to set Rating
func (options *CreateOfferingOptions) SetRating(rating *Rating) *CreateOfferingOptions {
	options.Rating = rating
	return options
}

// SetCreated : Allow user to set Created
func (options *CreateOfferingOptions) SetCreated(created *strfmt.DateTime) *CreateOfferingOptions {
	options.Created = created
	return options
}

// SetUpdated : Allow user to set Updated
func (options *CreateOfferingOptions) SetUpdated(updated *strfmt.DateTime) *CreateOfferingOptions {
	options.Updated = updated
	return options
}

// SetShortDescription : Allow user to set ShortDescription
func (options *CreateOfferingOptions) SetShortDescription(shortDescription string) *CreateOfferingOptions {
	options.ShortDescription = core.StringPtr(shortDescription)
	return options
}

// SetLongDescription : Allow user to set LongDescription
func (options *CreateOfferingOptions) SetLongDescription(longDescription string) *CreateOfferingOptions {
	options.LongDescription = core.StringPtr(longDescription)
	return options
}

// SetFeatures : Allow user to set Features
func (options *CreateOfferingOptions) SetFeatures(features []Feature) *CreateOfferingOptions {
	options.Features = features
	return options
}

// SetKinds : Allow user to set Kinds
func (options *CreateOfferingOptions) SetKinds(kinds []Kind) *CreateOfferingOptions {
	options.Kinds = kinds
	return options
}

// SetPermitRequestIBMPublicPublish : Allow user to set PermitRequestIBMPublicPublish
func (options *CreateOfferingOptions) SetPermitRequestIBMPublicPublish(permitRequestIBMPublicPublish bool) *CreateOfferingOptions {
	options.PermitRequestIBMPublicPublish = core.BoolPtr(permitRequestIBMPublicPublish)
	return options
}

// SetIBMPublishApproved : Allow user to set IBMPublishApproved
func (options *CreateOfferingOptions) SetIBMPublishApproved(ibmPublishApproved bool) *CreateOfferingOptions {
	options.IBMPublishApproved = core.BoolPtr(ibmPublishApproved)
	return options
}

// SetPublicPublishApproved : Allow user to set PublicPublishApproved
func (options *CreateOfferingOptions) SetPublicPublishApproved(publicPublishApproved bool) *CreateOfferingOptions {
	options.PublicPublishApproved = core.BoolPtr(publicPublishApproved)
	return options
}

// SetPublicOriginalCRN : Allow user to set PublicOriginalCRN
func (options *CreateOfferingOptions) SetPublicOriginalCRN(publicOriginalCRN string) *CreateOfferingOptions {
	options.PublicOriginalCRN = core.StringPtr(publicOriginalCRN)
	return options
}

// SetPublishPublicCRN : Allow user to set PublishPublicCRN
func (options *CreateOfferingOptions) SetPublishPublicCRN(publishPublicCRN string) *CreateOfferingOptions {
	options.PublishPublicCRN = core.StringPtr(publishPublicCRN)
	return options
}

// SetPortalApprovalRecord : Allow user to set PortalApprovalRecord
func (options *CreateOfferingOptions) SetPortalApprovalRecord(portalApprovalRecord string) *CreateOfferingOptions {
	options.PortalApprovalRecord = core.StringPtr(portalApprovalRecord)
	return options
}

// SetPortalUIURL : Allow user to set PortalUIURL
func (options *CreateOfferingOptions) SetPortalUIURL(portalUIURL string) *CreateOfferingOptions {
	options.PortalUIURL = core.StringPtr(portalUIURL)
	return options
}

// SetCatalogID : Allow user to set CatalogID
func (options *CreateOfferingOptions) SetCatalogID(catalogID string) *CreateOfferingOptions {
	options.CatalogID = core.StringPtr(catalogID)
	return options
}

// SetCatalogName : Allow user to set CatalogName
func (options *CreateOfferingOptions) SetCatalogName(catalogName string) *CreateOfferingOptions {
	options.CatalogName = core.StringPtr(catalogName)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *CreateOfferingOptions) SetMetadata(metadata map[string]interface{}) *CreateOfferingOptions {
	options.Metadata = metadata
	return options
}

// SetDisclaimer : Allow user to set Disclaimer
func (options *CreateOfferingOptions) SetDisclaimer(disclaimer string) *CreateOfferingOptions {
	options.Disclaimer = core.StringPtr(disclaimer)
	return options
}

// SetHidden : Allow user to set Hidden
func (options *CreateOfferingOptions) SetHidden(hidden bool) *CreateOfferingOptions {
	options.Hidden = core.BoolPtr(hidden)
	return options
}

// SetProvider : Allow user to set Provider
func (options *CreateOfferingOptions) SetProvider(provider string) *CreateOfferingOptions {
	options.Provider = core.StringPtr(provider)
	return options
}

// SetRepoInfo : Allow user to set RepoInfo
func (options *CreateOfferingOptions) SetRepoInfo(repoInfo *RepoInfo) *CreateOfferingOptions {
	options.RepoInfo = repoInfo
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateOfferingOptions) SetHeaders(param map[string]string) *CreateOfferingOptions {
	options.Headers = param
	return options
}

// DeleteCatalogOptions : The DeleteCatalog options.
type DeleteCatalogOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCatalogOptions : Instantiate DeleteCatalogOptions
func (*CatalogManagementV1) NewDeleteCatalogOptions(catalogIdentifier string) *DeleteCatalogOptions {
	return &DeleteCatalogOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *DeleteCatalogOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteCatalogOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCatalogOptions) SetHeaders(param map[string]string) *DeleteCatalogOptions {
	options.Headers = param
	return options
}

// DeleteObjectAccessListOptions : The DeleteObjectAccessList options.
type DeleteObjectAccessListOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// A list of accounts to delete.  An entry with star["*"] will remove all accounts.
	Accounts []string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteObjectAccessListOptions : Instantiate DeleteObjectAccessListOptions
func (*CatalogManagementV1) NewDeleteObjectAccessListOptions(catalogIdentifier string, objectIdentifier string, accounts []string) *DeleteObjectAccessListOptions {
	return &DeleteObjectAccessListOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		Accounts: accounts,
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *DeleteObjectAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteObjectAccessListOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *DeleteObjectAccessListOptions) SetObjectIdentifier(objectIdentifier string) *DeleteObjectAccessListOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetAccounts : Allow user to set Accounts
func (options *DeleteObjectAccessListOptions) SetAccounts(accounts []string) *DeleteObjectAccessListOptions {
	options.Accounts = accounts
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteObjectAccessListOptions) SetHeaders(param map[string]string) *DeleteObjectAccessListOptions {
	options.Headers = param
	return options
}

// DeleteObjectAccessOptions : The DeleteObjectAccess options.
type DeleteObjectAccessOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Account identifier.
	AccountIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteObjectAccessOptions : Instantiate DeleteObjectAccessOptions
func (*CatalogManagementV1) NewDeleteObjectAccessOptions(catalogIdentifier string, objectIdentifier string, accountIdentifier string) *DeleteObjectAccessOptions {
	return &DeleteObjectAccessOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		AccountIdentifier: core.StringPtr(accountIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *DeleteObjectAccessOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteObjectAccessOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *DeleteObjectAccessOptions) SetObjectIdentifier(objectIdentifier string) *DeleteObjectAccessOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetAccountIdentifier : Allow user to set AccountIdentifier
func (options *DeleteObjectAccessOptions) SetAccountIdentifier(accountIdentifier string) *DeleteObjectAccessOptions {
	options.AccountIdentifier = core.StringPtr(accountIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteObjectAccessOptions) SetHeaders(param map[string]string) *DeleteObjectAccessOptions {
	options.Headers = param
	return options
}

// DeleteObjectOptions : The DeleteObject options.
type DeleteObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteObjectOptions : Instantiate DeleteObjectOptions
func (*CatalogManagementV1) NewDeleteObjectOptions(catalogIdentifier string, objectIdentifier string) *DeleteObjectOptions {
	return &DeleteObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *DeleteObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteObjectOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *DeleteObjectOptions) SetObjectIdentifier(objectIdentifier string) *DeleteObjectOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteObjectOptions) SetHeaders(param map[string]string) *DeleteObjectOptions {
	options.Headers = param
	return options
}

// DeleteOfferingInstanceOptions : The DeleteOfferingInstance options.
type DeleteOfferingInstanceOptions struct {
	// Version Instance identifier.
	InstanceIdentifier *string `validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteOfferingInstanceOptions : Instantiate DeleteOfferingInstanceOptions
func (*CatalogManagementV1) NewDeleteOfferingInstanceOptions(instanceIdentifier string, xAuthRefreshToken string) *DeleteOfferingInstanceOptions {
	return &DeleteOfferingInstanceOptions{
		InstanceIdentifier: core.StringPtr(instanceIdentifier),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetInstanceIdentifier : Allow user to set InstanceIdentifier
func (options *DeleteOfferingInstanceOptions) SetInstanceIdentifier(instanceIdentifier string) *DeleteOfferingInstanceOptions {
	options.InstanceIdentifier = core.StringPtr(instanceIdentifier)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *DeleteOfferingInstanceOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *DeleteOfferingInstanceOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOfferingInstanceOptions) SetHeaders(param map[string]string) *DeleteOfferingInstanceOptions {
	options.Headers = param
	return options
}

// DeleteOfferingOptions : The DeleteOffering options.
type DeleteOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identification.
	OfferingID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteOfferingOptions : Instantiate DeleteOfferingOptions
func (*CatalogManagementV1) NewDeleteOfferingOptions(catalogIdentifier string, offeringID string) *DeleteOfferingOptions {
	return &DeleteOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *DeleteOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteOfferingOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *DeleteOfferingOptions) SetOfferingID(offeringID string) *DeleteOfferingOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOfferingOptions) SetHeaders(param map[string]string) *DeleteOfferingOptions {
	options.Headers = param
	return options
}

// DeleteOperatorsOptions : The DeleteOperators options.
type DeleteOperatorsOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Cluster identification.
	ClusterID *string `validate:"required"`

	// Cluster region.
	Region *string `validate:"required"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteOperatorsOptions : Instantiate DeleteOperatorsOptions
func (*CatalogManagementV1) NewDeleteOperatorsOptions(xAuthRefreshToken string, clusterID string, region string, versionLocatorID string) *DeleteOperatorsOptions {
	return &DeleteOperatorsOptions{
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
		ClusterID: core.StringPtr(clusterID),
		Region: core.StringPtr(region),
		VersionLocatorID: core.StringPtr(versionLocatorID),
	}
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *DeleteOperatorsOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *DeleteOperatorsOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *DeleteOperatorsOptions) SetClusterID(clusterID string) *DeleteOperatorsOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *DeleteOperatorsOptions) SetRegion(region string) *DeleteOperatorsOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (options *DeleteOperatorsOptions) SetVersionLocatorID(versionLocatorID string) *DeleteOperatorsOptions {
	options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOperatorsOptions) SetHeaders(param map[string]string) *DeleteOperatorsOptions {
	options.Headers = param
	return options
}

// DeleteVersionOptions : The DeleteVersion options.
type DeleteVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteVersionOptions : Instantiate DeleteVersionOptions
func (*CatalogManagementV1) NewDeleteVersionOptions(versionLocID string) *DeleteVersionOptions {
	return &DeleteVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *DeleteVersionOptions) SetVersionLocID(versionLocID string) *DeleteVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteVersionOptions) SetHeaders(param map[string]string) *DeleteVersionOptions {
	options.Headers = param
	return options
}

// DeployOperatorsOptions : The DeployOperators options.
type DeployOperatorsOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Cluster ID.
	ClusterID *string

	// Cluster region.
	Region *string

	// Kube namespaces to deploy Operator(s) to.
	Namespaces []string

	// Denotes whether to install Operator(s) globally.
	AllNamespaces *bool

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeployOperatorsOptions : Instantiate DeployOperatorsOptions
func (*CatalogManagementV1) NewDeployOperatorsOptions(xAuthRefreshToken string) *DeployOperatorsOptions {
	return &DeployOperatorsOptions{
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *DeployOperatorsOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *DeployOperatorsOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *DeployOperatorsOptions) SetClusterID(clusterID string) *DeployOperatorsOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *DeployOperatorsOptions) SetRegion(region string) *DeployOperatorsOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetNamespaces : Allow user to set Namespaces
func (options *DeployOperatorsOptions) SetNamespaces(namespaces []string) *DeployOperatorsOptions {
	options.Namespaces = namespaces
	return options
}

// SetAllNamespaces : Allow user to set AllNamespaces
func (options *DeployOperatorsOptions) SetAllNamespaces(allNamespaces bool) *DeployOperatorsOptions {
	options.AllNamespaces = core.BoolPtr(allNamespaces)
	return options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (options *DeployOperatorsOptions) SetVersionLocatorID(versionLocatorID string) *DeployOperatorsOptions {
	options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeployOperatorsOptions) SetHeaders(param map[string]string) *DeployOperatorsOptions {
	options.Headers = param
	return options
}

// DeployRequestBodySchematics : Schematics workspace configuration.
type DeployRequestBodySchematics struct {
	// Schematics workspace name.
	Name *string `json:"name,omitempty"`

	// Schematics workspace description.
	Description *string `json:"description,omitempty"`

	// Schematics workspace tags.
	Tags []string `json:"tags,omitempty"`

	// Resource group to use when creating the schematics workspace.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`
}

// UnmarshalDeployRequestBodySchematics unmarshals an instance of DeployRequestBodySchematics from the specified map of raw messages.
func UnmarshalDeployRequestBodySchematics(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeployRequestBodySchematics)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Deployment : Deployment for offering.
type Deployment struct {
	// unique id.
	ID *string `json:"id,omitempty"`

	// Display Name in the requested language.
	Label *string `json:"label,omitempty"`

	// The programmatic name of this offering.
	Name *string `json:"name,omitempty"`

	// Short description in the requested language.
	ShortDescription *string `json:"short_description,omitempty"`

	// Long description in the requested language.
	LongDescription *string `json:"long_description,omitempty"`

	// open ended metadata information.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// list of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// the date'time this catalog was created.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// the date'time this catalog was last updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`
}

// UnmarshalDeployment unmarshals an instance of Deployment from the specified map of raw messages.
func UnmarshalDeployment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Deployment)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "short_description", &obj.ShortDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeprecateVersionOptions : The DeprecateVersion options.
type DeprecateVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeprecateVersionOptions : Instantiate DeprecateVersionOptions
func (*CatalogManagementV1) NewDeprecateVersionOptions(versionLocID string) *DeprecateVersionOptions {
	return &DeprecateVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *DeprecateVersionOptions) SetVersionLocID(versionLocID string) *DeprecateVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeprecateVersionOptions) SetHeaders(param map[string]string) *DeprecateVersionOptions {
	options.Headers = param
	return options
}

// Feature : Feature information.
type Feature struct {
	// Heading.
	Title *string `json:"title,omitempty"`

	// Feature description.
	Description *string `json:"description,omitempty"`
}

// UnmarshalFeature unmarshals an instance of Feature from the specified map of raw messages.
func UnmarshalFeature(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Feature)
	err = core.UnmarshalPrimitive(m, "title", &obj.Title)
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

// FilterTerms : Offering filter terms.
type FilterTerms struct {
	// List of values to match against. If include is true, then if the offering has one of the values then the offering is
	// included. If include is false, then if the offering has one of the values then the offering is excluded.
	FilterTerms []string `json:"filter_terms,omitempty"`
}

// UnmarshalFilterTerms unmarshals an instance of FilterTerms from the specified map of raw messages.
func UnmarshalFilterTerms(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FilterTerms)
	err = core.UnmarshalPrimitive(m, "filter_terms", &obj.FilterTerms)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Filters : Filters for account and catalog filters.
type Filters struct {
	// -> true - Include all of the public catalog when filtering. Further settings will specifically exclude some
	// offerings. false - Exclude all of the public catalog when filtering. Further settings will specifically include some
	// offerings.
	IncludeAll *bool `json:"include_all,omitempty"`

	// Filter against offering properties.
	CategoryFilters map[string]CategoryFilter `json:"category_filters,omitempty"`

	// Filter on offering ID's. There is an include filter and an exclule filter. Both can be set.
	IDFilters *IDFilter `json:"id_filters,omitempty"`
}

// UnmarshalFilters unmarshals an instance of Filters from the specified map of raw messages.
func UnmarshalFilters(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Filters)
	err = core.UnmarshalPrimitive(m, "include_all", &obj.IncludeAll)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "category_filters", &obj.CategoryFilters, UnmarshalCategoryFilter)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "id_filters", &obj.IDFilters, UnmarshalIDFilter)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetCatalogAccountAuditOptions : The GetCatalogAccountAudit options.
type GetCatalogAccountAuditOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCatalogAccountAuditOptions : Instantiate GetCatalogAccountAuditOptions
func (*CatalogManagementV1) NewGetCatalogAccountAuditOptions() *GetCatalogAccountAuditOptions {
	return &GetCatalogAccountAuditOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogAccountAuditOptions) SetHeaders(param map[string]string) *GetCatalogAccountAuditOptions {
	options.Headers = param
	return options
}

// GetCatalogAccountFiltersOptions : The GetCatalogAccountFilters options.
type GetCatalogAccountFiltersOptions struct {
	// catalog id. Narrow down filters to the account and just the one catalog.
	Catalog *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCatalogAccountFiltersOptions : Instantiate GetCatalogAccountFiltersOptions
func (*CatalogManagementV1) NewGetCatalogAccountFiltersOptions() *GetCatalogAccountFiltersOptions {
	return &GetCatalogAccountFiltersOptions{}
}

// SetCatalog : Allow user to set Catalog
func (options *GetCatalogAccountFiltersOptions) SetCatalog(catalog string) *GetCatalogAccountFiltersOptions {
	options.Catalog = core.StringPtr(catalog)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogAccountFiltersOptions) SetHeaders(param map[string]string) *GetCatalogAccountFiltersOptions {
	options.Headers = param
	return options
}

// GetCatalogAccountOptions : The GetCatalogAccount options.
type GetCatalogAccountOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCatalogAccountOptions : Instantiate GetCatalogAccountOptions
func (*CatalogManagementV1) NewGetCatalogAccountOptions() *GetCatalogAccountOptions {
	return &GetCatalogAccountOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogAccountOptions) SetHeaders(param map[string]string) *GetCatalogAccountOptions {
	options.Headers = param
	return options
}

// GetCatalogAuditOptions : The GetCatalogAudit options.
type GetCatalogAuditOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCatalogAuditOptions : Instantiate GetCatalogAuditOptions
func (*CatalogManagementV1) NewGetCatalogAuditOptions(catalogIdentifier string) *GetCatalogAuditOptions {
	return &GetCatalogAuditOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetCatalogAuditOptions) SetCatalogIdentifier(catalogIdentifier string) *GetCatalogAuditOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogAuditOptions) SetHeaders(param map[string]string) *GetCatalogAuditOptions {
	options.Headers = param
	return options
}

// GetCatalogOptions : The GetCatalog options.
type GetCatalogOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCatalogOptions : Instantiate GetCatalogOptions
func (*CatalogManagementV1) NewGetCatalogOptions(catalogIdentifier string) *GetCatalogOptions {
	return &GetCatalogOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetCatalogOptions) SetCatalogIdentifier(catalogIdentifier string) *GetCatalogOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogOptions) SetHeaders(param map[string]string) *GetCatalogOptions {
	options.Headers = param
	return options
}

// GetClusterOptions : The GetCluster options.
type GetClusterOptions struct {
	// ID of the cluster.
	ClusterID *string `validate:"required,ne="`

	// Region of the cluster.
	Region *string `validate:"required"`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetClusterOptions : Instantiate GetClusterOptions
func (*CatalogManagementV1) NewGetClusterOptions(clusterID string, region string, xAuthRefreshToken string) *GetClusterOptions {
	return &GetClusterOptions{
		ClusterID: core.StringPtr(clusterID),
		Region: core.StringPtr(region),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetClusterID : Allow user to set ClusterID
func (options *GetClusterOptions) SetClusterID(clusterID string) *GetClusterOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *GetClusterOptions) SetRegion(region string) *GetClusterOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *GetClusterOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetClusterOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetClusterOptions) SetHeaders(param map[string]string) *GetClusterOptions {
	options.Headers = param
	return options
}

// GetConsumptionOfferingsOptions : The GetConsumptionOfferings options.
type GetConsumptionOfferingsOptions struct {
	// true - Strip down the content of what is returned. For example don't return the readme. Makes the result much
	// smaller. Defaults to false.
	Digest *bool

	// catalog id. Narrow search down to just a particular catalog. It will apply the catalog's public filters to the
	// public catalog offerings on the result.
	Catalog *string

	// What should be selected. Default is 'all' which will return both public and private offerings. 'public' returns only
	// the public offerings and 'private' returns only the private offerings.
	Select *string

	// true - include offerings which have been marked as hidden. The default is false and hidden offerings are not
	// returned.
	IncludeHidden *bool

	// number or results to return.
	Limit *int64

	// number of results to skip before returning values.
	Offset *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetConsumptionOfferingsOptions.Select property.
// What should be selected. Default is 'all' which will return both public and private offerings. 'public' returns only
// the public offerings and 'private' returns only the private offerings.
const (
	GetConsumptionOfferingsOptionsSelectAllConst = "all"
	GetConsumptionOfferingsOptionsSelectPrivateConst = "private"
	GetConsumptionOfferingsOptionsSelectPublicConst = "public"
)

// NewGetConsumptionOfferingsOptions : Instantiate GetConsumptionOfferingsOptions
func (*CatalogManagementV1) NewGetConsumptionOfferingsOptions() *GetConsumptionOfferingsOptions {
	return &GetConsumptionOfferingsOptions{}
}

// SetDigest : Allow user to set Digest
func (options *GetConsumptionOfferingsOptions) SetDigest(digest bool) *GetConsumptionOfferingsOptions {
	options.Digest = core.BoolPtr(digest)
	return options
}

// SetCatalog : Allow user to set Catalog
func (options *GetConsumptionOfferingsOptions) SetCatalog(catalog string) *GetConsumptionOfferingsOptions {
	options.Catalog = core.StringPtr(catalog)
	return options
}

// SetSelect : Allow user to set Select
func (options *GetConsumptionOfferingsOptions) SetSelect(selectVar string) *GetConsumptionOfferingsOptions {
	options.Select = core.StringPtr(selectVar)
	return options
}

// SetIncludeHidden : Allow user to set IncludeHidden
func (options *GetConsumptionOfferingsOptions) SetIncludeHidden(includeHidden bool) *GetConsumptionOfferingsOptions {
	options.IncludeHidden = core.BoolPtr(includeHidden)
	return options
}

// SetLimit : Allow user to set Limit
func (options *GetConsumptionOfferingsOptions) SetLimit(limit int64) *GetConsumptionOfferingsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *GetConsumptionOfferingsOptions) SetOffset(offset int64) *GetConsumptionOfferingsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetConsumptionOfferingsOptions) SetHeaders(param map[string]string) *GetConsumptionOfferingsOptions {
	options.Headers = param
	return options
}

// GetNamespacesOptions : The GetNamespaces options.
type GetNamespacesOptions struct {
	// ID of the cluster.
	ClusterID *string `validate:"required,ne="`

	// Cluster region.
	Region *string `validate:"required"`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// The maximum number of results to return.
	Limit *int64

	// The number of results to skip before returning values.
	Offset *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetNamespacesOptions : Instantiate GetNamespacesOptions
func (*CatalogManagementV1) NewGetNamespacesOptions(clusterID string, region string, xAuthRefreshToken string) *GetNamespacesOptions {
	return &GetNamespacesOptions{
		ClusterID: core.StringPtr(clusterID),
		Region: core.StringPtr(region),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetClusterID : Allow user to set ClusterID
func (options *GetNamespacesOptions) SetClusterID(clusterID string) *GetNamespacesOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *GetNamespacesOptions) SetRegion(region string) *GetNamespacesOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *GetNamespacesOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetNamespacesOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetLimit : Allow user to set Limit
func (options *GetNamespacesOptions) SetLimit(limit int64) *GetNamespacesOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *GetNamespacesOptions) SetOffset(offset int64) *GetNamespacesOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetNamespacesOptions) SetHeaders(param map[string]string) *GetNamespacesOptions {
	options.Headers = param
	return options
}

// GetObjectAccessListOptions : The GetObjectAccessList options.
type GetObjectAccessListOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// The maximum number of results to return.
	Limit *int64

	// The number of results to skip before returning values.
	Offset *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetObjectAccessListOptions : Instantiate GetObjectAccessListOptions
func (*CatalogManagementV1) NewGetObjectAccessListOptions(catalogIdentifier string, objectIdentifier string) *GetObjectAccessListOptions {
	return &GetObjectAccessListOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetObjectAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectAccessListOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *GetObjectAccessListOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectAccessListOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetLimit : Allow user to set Limit
func (options *GetObjectAccessListOptions) SetLimit(limit int64) *GetObjectAccessListOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *GetObjectAccessListOptions) SetOffset(offset int64) *GetObjectAccessListOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectAccessListOptions) SetHeaders(param map[string]string) *GetObjectAccessListOptions {
	options.Headers = param
	return options
}

// GetObjectAccessOptions : The GetObjectAccess options.
type GetObjectAccessOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Account identifier.
	AccountIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetObjectAccessOptions : Instantiate GetObjectAccessOptions
func (*CatalogManagementV1) NewGetObjectAccessOptions(catalogIdentifier string, objectIdentifier string, accountIdentifier string) *GetObjectAccessOptions {
	return &GetObjectAccessOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		AccountIdentifier: core.StringPtr(accountIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetObjectAccessOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectAccessOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *GetObjectAccessOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectAccessOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetAccountIdentifier : Allow user to set AccountIdentifier
func (options *GetObjectAccessOptions) SetAccountIdentifier(accountIdentifier string) *GetObjectAccessOptions {
	options.AccountIdentifier = core.StringPtr(accountIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectAccessOptions) SetHeaders(param map[string]string) *GetObjectAccessOptions {
	options.Headers = param
	return options
}

// GetObjectAuditOptions : The GetObjectAudit options.
type GetObjectAuditOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetObjectAuditOptions : Instantiate GetObjectAuditOptions
func (*CatalogManagementV1) NewGetObjectAuditOptions(catalogIdentifier string, objectIdentifier string) *GetObjectAuditOptions {
	return &GetObjectAuditOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetObjectAuditOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectAuditOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *GetObjectAuditOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectAuditOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectAuditOptions) SetHeaders(param map[string]string) *GetObjectAuditOptions {
	options.Headers = param
	return options
}

// GetObjectOptions : The GetObject options.
type GetObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetObjectOptions : Instantiate GetObjectOptions
func (*CatalogManagementV1) NewGetObjectOptions(catalogIdentifier string, objectIdentifier string) *GetObjectOptions {
	return &GetObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *GetObjectOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectOptions) SetHeaders(param map[string]string) *GetObjectOptions {
	options.Headers = param
	return options
}

// GetOfferingAboutOptions : The GetOfferingAbout options.
type GetOfferingAboutOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingAboutOptions : Instantiate GetOfferingAboutOptions
func (*CatalogManagementV1) NewGetOfferingAboutOptions(versionLocID string) *GetOfferingAboutOptions {
	return &GetOfferingAboutOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *GetOfferingAboutOptions) SetVersionLocID(versionLocID string) *GetOfferingAboutOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingAboutOptions) SetHeaders(param map[string]string) *GetOfferingAboutOptions {
	options.Headers = param
	return options
}

// GetOfferingAuditOptions : The GetOfferingAudit options.
type GetOfferingAuditOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identifier.
	OfferingID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingAuditOptions : Instantiate GetOfferingAuditOptions
func (*CatalogManagementV1) NewGetOfferingAuditOptions(catalogIdentifier string, offeringID string) *GetOfferingAuditOptions {
	return &GetOfferingAuditOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetOfferingAuditOptions) SetCatalogIdentifier(catalogIdentifier string) *GetOfferingAuditOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *GetOfferingAuditOptions) SetOfferingID(offeringID string) *GetOfferingAuditOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingAuditOptions) SetHeaders(param map[string]string) *GetOfferingAuditOptions {
	options.Headers = param
	return options
}

// GetOfferingContainerImagesOptions : The GetOfferingContainerImages options.
type GetOfferingContainerImagesOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingContainerImagesOptions : Instantiate GetOfferingContainerImagesOptions
func (*CatalogManagementV1) NewGetOfferingContainerImagesOptions(versionLocID string) *GetOfferingContainerImagesOptions {
	return &GetOfferingContainerImagesOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *GetOfferingContainerImagesOptions) SetVersionLocID(versionLocID string) *GetOfferingContainerImagesOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingContainerImagesOptions) SetHeaders(param map[string]string) *GetOfferingContainerImagesOptions {
	options.Headers = param
	return options
}

// GetOfferingInstanceOptions : The GetOfferingInstance options.
type GetOfferingInstanceOptions struct {
	// Version Instance identifier.
	InstanceIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingInstanceOptions : Instantiate GetOfferingInstanceOptions
func (*CatalogManagementV1) NewGetOfferingInstanceOptions(instanceIdentifier string) *GetOfferingInstanceOptions {
	return &GetOfferingInstanceOptions{
		InstanceIdentifier: core.StringPtr(instanceIdentifier),
	}
}

// SetInstanceIdentifier : Allow user to set InstanceIdentifier
func (options *GetOfferingInstanceOptions) SetInstanceIdentifier(instanceIdentifier string) *GetOfferingInstanceOptions {
	options.InstanceIdentifier = core.StringPtr(instanceIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingInstanceOptions) SetHeaders(param map[string]string) *GetOfferingInstanceOptions {
	options.Headers = param
	return options
}

// GetOfferingLicenseOptions : The GetOfferingLicense options.
type GetOfferingLicenseOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// The ID of the license, which maps to the file name in the 'licenses' directory of this verions tgz file.
	LicenseID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingLicenseOptions : Instantiate GetOfferingLicenseOptions
func (*CatalogManagementV1) NewGetOfferingLicenseOptions(versionLocID string, licenseID string) *GetOfferingLicenseOptions {
	return &GetOfferingLicenseOptions{
		VersionLocID: core.StringPtr(versionLocID),
		LicenseID: core.StringPtr(licenseID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *GetOfferingLicenseOptions) SetVersionLocID(versionLocID string) *GetOfferingLicenseOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetLicenseID : Allow user to set LicenseID
func (options *GetOfferingLicenseOptions) SetLicenseID(licenseID string) *GetOfferingLicenseOptions {
	options.LicenseID = core.StringPtr(licenseID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingLicenseOptions) SetHeaders(param map[string]string) *GetOfferingLicenseOptions {
	options.Headers = param
	return options
}

// GetOfferingOptions : The GetOffering options.
type GetOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identification.
	OfferingID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingOptions : Instantiate GetOfferingOptions
func (*CatalogManagementV1) NewGetOfferingOptions(catalogIdentifier string, offeringID string) *GetOfferingOptions {
	return &GetOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *GetOfferingOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *GetOfferingOptions) SetOfferingID(offeringID string) *GetOfferingOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingOptions) SetHeaders(param map[string]string) *GetOfferingOptions {
	options.Headers = param
	return options
}

// GetOfferingUpdatesOptions : The GetOfferingUpdates options.
type GetOfferingUpdatesOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identification.
	OfferingID *string `validate:"required,ne="`

	// The kind of offering (e.g, helm, ova, terraform ...).
	Kind *string `validate:"required"`

	// optionaly provide an existing version to check updates for if one is not given, all version will be returned.
	Version *string

	// The id of the cluster where this version was installed.
	ClusterID *string

	// The region of the cluster where this version was installed.
	Region *string

	// The resource group id of the cluster where this version was installed.
	ResourceGroupID *string

	// The namespace of the cluster where this version was installed.
	Namespace *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingUpdatesOptions : Instantiate GetOfferingUpdatesOptions
func (*CatalogManagementV1) NewGetOfferingUpdatesOptions(catalogIdentifier string, offeringID string, kind string) *GetOfferingUpdatesOptions {
	return &GetOfferingUpdatesOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		Kind: core.StringPtr(kind),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *GetOfferingUpdatesOptions) SetCatalogIdentifier(catalogIdentifier string) *GetOfferingUpdatesOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *GetOfferingUpdatesOptions) SetOfferingID(offeringID string) *GetOfferingUpdatesOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetKind : Allow user to set Kind
func (options *GetOfferingUpdatesOptions) SetKind(kind string) *GetOfferingUpdatesOptions {
	options.Kind = core.StringPtr(kind)
	return options
}

// SetVersion : Allow user to set Version
func (options *GetOfferingUpdatesOptions) SetVersion(version string) *GetOfferingUpdatesOptions {
	options.Version = core.StringPtr(version)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *GetOfferingUpdatesOptions) SetClusterID(clusterID string) *GetOfferingUpdatesOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *GetOfferingUpdatesOptions) SetRegion(region string) *GetOfferingUpdatesOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (options *GetOfferingUpdatesOptions) SetResourceGroupID(resourceGroupID string) *GetOfferingUpdatesOptions {
	options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return options
}

// SetNamespace : Allow user to set Namespace
func (options *GetOfferingUpdatesOptions) SetNamespace(namespace string) *GetOfferingUpdatesOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingUpdatesOptions) SetHeaders(param map[string]string) *GetOfferingUpdatesOptions {
	options.Headers = param
	return options
}

// GetOfferingWorkingCopyOptions : The GetOfferingWorkingCopy options.
type GetOfferingWorkingCopyOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingWorkingCopyOptions : Instantiate GetOfferingWorkingCopyOptions
func (*CatalogManagementV1) NewGetOfferingWorkingCopyOptions(versionLocID string) *GetOfferingWorkingCopyOptions {
	return &GetOfferingWorkingCopyOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *GetOfferingWorkingCopyOptions) SetVersionLocID(versionLocID string) *GetOfferingWorkingCopyOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingWorkingCopyOptions) SetHeaders(param map[string]string) *GetOfferingWorkingCopyOptions {
	options.Headers = param
	return options
}

// GetOverrideValuesOptions : The GetOverrideValues options.
type GetOverrideValuesOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOverrideValuesOptions : Instantiate GetOverrideValuesOptions
func (*CatalogManagementV1) NewGetOverrideValuesOptions(versionLocID string) *GetOverrideValuesOptions {
	return &GetOverrideValuesOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *GetOverrideValuesOptions) SetVersionLocID(versionLocID string) *GetOverrideValuesOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetOverrideValuesOptions) SetHeaders(param map[string]string) *GetOverrideValuesOptions {
	options.Headers = param
	return options
}

// GetPreinstallOptions : The GetPreinstall options.
type GetPreinstallOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// ID of the cluster.
	ClusterID *string

	// Cluster region.
	Region *string

	// Required if the version's pre-install scope is `namespace`.
	Namespace *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetPreinstallOptions : Instantiate GetPreinstallOptions
func (*CatalogManagementV1) NewGetPreinstallOptions(versionLocID string, xAuthRefreshToken string) *GetPreinstallOptions {
	return &GetPreinstallOptions{
		VersionLocID: core.StringPtr(versionLocID),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *GetPreinstallOptions) SetVersionLocID(versionLocID string) *GetPreinstallOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *GetPreinstallOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetPreinstallOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *GetPreinstallOptions) SetClusterID(clusterID string) *GetPreinstallOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *GetPreinstallOptions) SetRegion(region string) *GetPreinstallOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetNamespace : Allow user to set Namespace
func (options *GetPreinstallOptions) SetNamespace(namespace string) *GetPreinstallOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetPreinstallOptions) SetHeaders(param map[string]string) *GetPreinstallOptions {
	options.Headers = param
	return options
}

// GetValidationStatusOptions : The GetValidationStatus options.
type GetValidationStatusOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetValidationStatusOptions : Instantiate GetValidationStatusOptions
func (*CatalogManagementV1) NewGetValidationStatusOptions(versionLocID string, xAuthRefreshToken string) *GetValidationStatusOptions {
	return &GetValidationStatusOptions{
		VersionLocID: core.StringPtr(versionLocID),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *GetValidationStatusOptions) SetVersionLocID(versionLocID string) *GetValidationStatusOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *GetValidationStatusOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetValidationStatusOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetValidationStatusOptions) SetHeaders(param map[string]string) *GetValidationStatusOptions {
	options.Headers = param
	return options
}

// GetVersionOptions : The GetVersion options.
type GetVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetVersionOptions : Instantiate GetVersionOptions
func (*CatalogManagementV1) NewGetVersionOptions(versionLocID string) *GetVersionOptions {
	return &GetVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *GetVersionOptions) SetVersionLocID(versionLocID string) *GetVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetVersionOptions) SetHeaders(param map[string]string) *GetVersionOptions {
	options.Headers = param
	return options
}

// IDFilter : Filter on offering ID's. There is an include filter and an exclule filter. Both can be set.
type IDFilter struct {
	// Offering filter terms.
	Include *FilterTerms `json:"include,omitempty"`

	// Offering filter terms.
	Exclude *FilterTerms `json:"exclude,omitempty"`
}

// UnmarshalIDFilter unmarshals an instance of IDFilter from the specified map of raw messages.
func UnmarshalIDFilter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IDFilter)
	err = core.UnmarshalModel(m, "include", &obj.Include, UnmarshalFilterTerms)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "exclude", &obj.Exclude, UnmarshalFilterTerms)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IBMPublishObjectOptions : The IBMPublishObject options.
type IBMPublishObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewIBMPublishObjectOptions : Instantiate IBMPublishObjectOptions
func (*CatalogManagementV1) NewIBMPublishObjectOptions(catalogIdentifier string, objectIdentifier string) *IBMPublishObjectOptions {
	return &IBMPublishObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *IBMPublishObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *IBMPublishObjectOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *IBMPublishObjectOptions) SetObjectIdentifier(objectIdentifier string) *IBMPublishObjectOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *IBMPublishObjectOptions) SetHeaders(param map[string]string) *IBMPublishObjectOptions {
	options.Headers = param
	return options
}

// IBMPublishVersionOptions : The IBMPublishVersion options.
type IBMPublishVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewIBMPublishVersionOptions : Instantiate IBMPublishVersionOptions
func (*CatalogManagementV1) NewIBMPublishVersionOptions(versionLocID string) *IBMPublishVersionOptions {
	return &IBMPublishVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *IBMPublishVersionOptions) SetVersionLocID(versionLocID string) *IBMPublishVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *IBMPublishVersionOptions) SetHeaders(param map[string]string) *IBMPublishVersionOptions {
	options.Headers = param
	return options
}

// Image : Image.
type Image struct {
	// Image.
	Image *string `json:"image,omitempty"`
}

// UnmarshalImage unmarshals an instance of Image from the specified map of raw messages.
func UnmarshalImage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Image)
	err = core.UnmarshalPrimitive(m, "image", &obj.Image)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImageManifest : Image Manifest.
type ImageManifest struct {
	// Image manifest description.
	Description *string `json:"description,omitempty"`

	// List of images.
	Images []Image `json:"images,omitempty"`
}

// UnmarshalImageManifest unmarshals an instance of ImageManifest from the specified map of raw messages.
func UnmarshalImageManifest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImageManifest)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "images", &obj.Images, UnmarshalImage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportOfferingOptions : The ImportOffering options.
type ImportOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Tags array.
	Tags []string

	// Target kinds.  Current valid values are 'iks', 'roks', 'vcenter', and 'terraform'.
	TargetKinds []string

	// byte array representing the content to be imported.  Only supported for OVA images at this time.
	Content *[]byte

	// URL path to zip location.  If not specified, must provide content in this post body.
	Zipurl *string

	// Re-use the specified offeringID during import.
	OfferingID *string

	// The semver value for this new version.
	TargetVersion *string

	// Add all possible configuration items when creating this version.
	IncludeConfig *bool

	// Indicates that the current terraform template is used to install a VSI Image.
	IsVsi *bool

	// The type of repository containing this version.  Valid values are 'public_git' or 'enterprise_git'.
	RepoType *string

	// Authentication token used to access the specified zip file.
	XAuthToken *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewImportOfferingOptions : Instantiate ImportOfferingOptions
func (*CatalogManagementV1) NewImportOfferingOptions(catalogIdentifier string) *ImportOfferingOptions {
	return &ImportOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ImportOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *ImportOfferingOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetTags : Allow user to set Tags
func (options *ImportOfferingOptions) SetTags(tags []string) *ImportOfferingOptions {
	options.Tags = tags
	return options
}

// SetTargetKinds : Allow user to set TargetKinds
func (options *ImportOfferingOptions) SetTargetKinds(targetKinds []string) *ImportOfferingOptions {
	options.TargetKinds = targetKinds
	return options
}

// SetContent : Allow user to set Content
func (options *ImportOfferingOptions) SetContent(content []byte) *ImportOfferingOptions {
	options.Content = &content
	return options
}

// SetZipurl : Allow user to set Zipurl
func (options *ImportOfferingOptions) SetZipurl(zipurl string) *ImportOfferingOptions {
	options.Zipurl = core.StringPtr(zipurl)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *ImportOfferingOptions) SetOfferingID(offeringID string) *ImportOfferingOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetTargetVersion : Allow user to set TargetVersion
func (options *ImportOfferingOptions) SetTargetVersion(targetVersion string) *ImportOfferingOptions {
	options.TargetVersion = core.StringPtr(targetVersion)
	return options
}

// SetIncludeConfig : Allow user to set IncludeConfig
func (options *ImportOfferingOptions) SetIncludeConfig(includeConfig bool) *ImportOfferingOptions {
	options.IncludeConfig = core.BoolPtr(includeConfig)
	return options
}

// SetIsVsi : Allow user to set IsVsi
func (options *ImportOfferingOptions) SetIsVsi(isVsi bool) *ImportOfferingOptions {
	options.IsVsi = core.BoolPtr(isVsi)
	return options
}

// SetRepoType : Allow user to set RepoType
func (options *ImportOfferingOptions) SetRepoType(repoType string) *ImportOfferingOptions {
	options.RepoType = core.StringPtr(repoType)
	return options
}

// SetXAuthToken : Allow user to set XAuthToken
func (options *ImportOfferingOptions) SetXAuthToken(xAuthToken string) *ImportOfferingOptions {
	options.XAuthToken = core.StringPtr(xAuthToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ImportOfferingOptions) SetHeaders(param map[string]string) *ImportOfferingOptions {
	options.Headers = param
	return options
}

// ImportOfferingVersionOptions : The ImportOfferingVersion options.
type ImportOfferingVersionOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identification.
	OfferingID *string `validate:"required,ne="`

	// Tags array.
	Tags []string

	// Target kinds.  Current valid values are 'iks', 'roks', 'vcenter', and 'terraform'.
	TargetKinds []string

	// byte array representing the content to be imported.  Only supported for OVA images at this time.
	Content *[]byte

	// URL path to zip location.  If not specified, must provide content in the body of this call.
	Zipurl *string

	// The semver value for this new version, if not found in the zip url package content.
	TargetVersion *string

	// Add all possible configuration values to this version when importing.
	IncludeConfig *bool

	// Indicates that the current terraform template is used to install a VSI Image.
	IsVsi *bool

	// The type of repository containing this version.  Valid values are 'public_git' or 'enterprise_git'.
	RepoType *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewImportOfferingVersionOptions : Instantiate ImportOfferingVersionOptions
func (*CatalogManagementV1) NewImportOfferingVersionOptions(catalogIdentifier string, offeringID string) *ImportOfferingVersionOptions {
	return &ImportOfferingVersionOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ImportOfferingVersionOptions) SetCatalogIdentifier(catalogIdentifier string) *ImportOfferingVersionOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *ImportOfferingVersionOptions) SetOfferingID(offeringID string) *ImportOfferingVersionOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetTags : Allow user to set Tags
func (options *ImportOfferingVersionOptions) SetTags(tags []string) *ImportOfferingVersionOptions {
	options.Tags = tags
	return options
}

// SetTargetKinds : Allow user to set TargetKinds
func (options *ImportOfferingVersionOptions) SetTargetKinds(targetKinds []string) *ImportOfferingVersionOptions {
	options.TargetKinds = targetKinds
	return options
}

// SetContent : Allow user to set Content
func (options *ImportOfferingVersionOptions) SetContent(content []byte) *ImportOfferingVersionOptions {
	options.Content = &content
	return options
}

// SetZipurl : Allow user to set Zipurl
func (options *ImportOfferingVersionOptions) SetZipurl(zipurl string) *ImportOfferingVersionOptions {
	options.Zipurl = core.StringPtr(zipurl)
	return options
}

// SetTargetVersion : Allow user to set TargetVersion
func (options *ImportOfferingVersionOptions) SetTargetVersion(targetVersion string) *ImportOfferingVersionOptions {
	options.TargetVersion = core.StringPtr(targetVersion)
	return options
}

// SetIncludeConfig : Allow user to set IncludeConfig
func (options *ImportOfferingVersionOptions) SetIncludeConfig(includeConfig bool) *ImportOfferingVersionOptions {
	options.IncludeConfig = core.BoolPtr(includeConfig)
	return options
}

// SetIsVsi : Allow user to set IsVsi
func (options *ImportOfferingVersionOptions) SetIsVsi(isVsi bool) *ImportOfferingVersionOptions {
	options.IsVsi = core.BoolPtr(isVsi)
	return options
}

// SetRepoType : Allow user to set RepoType
func (options *ImportOfferingVersionOptions) SetRepoType(repoType string) *ImportOfferingVersionOptions {
	options.RepoType = core.StringPtr(repoType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ImportOfferingVersionOptions) SetHeaders(param map[string]string) *ImportOfferingVersionOptions {
	options.Headers = param
	return options
}

// InstallStatus : Installation status.
type InstallStatus struct {
	// Installation status metadata.
	Metadata *InstallStatusMetadata `json:"metadata,omitempty"`

	// Release information.
	Release *InstallStatusRelease `json:"release,omitempty"`

	// Content management information.
	ContentMgmt *InstallStatusContentMgmt `json:"content_mgmt,omitempty"`
}

// UnmarshalInstallStatus unmarshals an instance of InstallStatus from the specified map of raw messages.
func UnmarshalInstallStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstallStatus)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalInstallStatusMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "release", &obj.Release, UnmarshalInstallStatusRelease)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "content_mgmt", &obj.ContentMgmt, UnmarshalInstallStatusContentMgmt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstallStatusContentMgmt : Content management information.
type InstallStatusContentMgmt struct {
	// Pods.
	Pods []map[string]string `json:"pods,omitempty"`

	// Errors.
	Errors []map[string]string `json:"errors,omitempty"`
}

// UnmarshalInstallStatusContentMgmt unmarshals an instance of InstallStatusContentMgmt from the specified map of raw messages.
func UnmarshalInstallStatusContentMgmt(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstallStatusContentMgmt)
	err = core.UnmarshalPrimitive(m, "pods", &obj.Pods)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstallStatusMetadata : Installation status metadata.
type InstallStatusMetadata struct {
	// Cluster ID.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster region.
	Region *string `json:"region,omitempty"`

	// Cluster namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Workspace ID.
	WorkspaceID *string `json:"workspace_id,omitempty"`

	// Workspace name.
	WorkspaceName *string `json:"workspace_name,omitempty"`
}

// UnmarshalInstallStatusMetadata unmarshals an instance of InstallStatusMetadata from the specified map of raw messages.
func UnmarshalInstallStatusMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstallStatusMetadata)
	err = core.UnmarshalPrimitive(m, "cluster_id", &obj.ClusterID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "workspace_id", &obj.WorkspaceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "workspace_name", &obj.WorkspaceName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstallStatusRelease : Release information.
type InstallStatusRelease struct {
	// Kube deployments.
	Deployments []map[string]interface{} `json:"deployments,omitempty"`

	// Kube replica sets.
	Replicasets []map[string]interface{} `json:"replicasets,omitempty"`

	// Kube stateful sets.
	Statefulsets []map[string]interface{} `json:"statefulsets,omitempty"`

	// Kube pods.
	Pods []map[string]interface{} `json:"pods,omitempty"`

	// Kube errors.
	Errors []map[string]string `json:"errors,omitempty"`
}

// UnmarshalInstallStatusRelease unmarshals an instance of InstallStatusRelease from the specified map of raw messages.
func UnmarshalInstallStatusRelease(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstallStatusRelease)
	err = core.UnmarshalPrimitive(m, "deployments", &obj.Deployments)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "replicasets", &obj.Replicasets)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "statefulsets", &obj.Statefulsets)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pods", &obj.Pods)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstallVersionOptions : The InstallVersion options.
type InstallVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Cluster ID.
	ClusterID *string

	// Cluster region.
	Region *string

	// Kube namespace.
	Namespace *string

	// Object containing Helm chart override values.  To use a secret for items of type password, specify a JSON encoded
	// value of $ref:#/components/schemas/SecretInstance, prefixed with `cmsm_v1:`.
	OverrideValues map[string]interface{}

	// Entitlement API Key for this offering.
	EntitlementApikey *string

	// Schematics workspace configuration.
	Schematics *DeployRequestBodySchematics

	// Script.
	Script *string

	// Script ID.
	ScriptID *string

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string

	// VCenter ID.
	VcenterID *string

	// VCenter User.
	VcenterUser *string

	// VCenter Password.
	VcenterPassword *string

	// VCenter Location.
	VcenterLocation *string

	// VCenter Datastore.
	VcenterDatastore *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewInstallVersionOptions : Instantiate InstallVersionOptions
func (*CatalogManagementV1) NewInstallVersionOptions(versionLocID string, xAuthRefreshToken string) *InstallVersionOptions {
	return &InstallVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *InstallVersionOptions) SetVersionLocID(versionLocID string) *InstallVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *InstallVersionOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *InstallVersionOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *InstallVersionOptions) SetClusterID(clusterID string) *InstallVersionOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *InstallVersionOptions) SetRegion(region string) *InstallVersionOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetNamespace : Allow user to set Namespace
func (options *InstallVersionOptions) SetNamespace(namespace string) *InstallVersionOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetOverrideValues : Allow user to set OverrideValues
func (options *InstallVersionOptions) SetOverrideValues(overrideValues map[string]interface{}) *InstallVersionOptions {
	options.OverrideValues = overrideValues
	return options
}

// SetEntitlementApikey : Allow user to set EntitlementApikey
func (options *InstallVersionOptions) SetEntitlementApikey(entitlementApikey string) *InstallVersionOptions {
	options.EntitlementApikey = core.StringPtr(entitlementApikey)
	return options
}

// SetSchematics : Allow user to set Schematics
func (options *InstallVersionOptions) SetSchematics(schematics *DeployRequestBodySchematics) *InstallVersionOptions {
	options.Schematics = schematics
	return options
}

// SetScript : Allow user to set Script
func (options *InstallVersionOptions) SetScript(script string) *InstallVersionOptions {
	options.Script = core.StringPtr(script)
	return options
}

// SetScriptID : Allow user to set ScriptID
func (options *InstallVersionOptions) SetScriptID(scriptID string) *InstallVersionOptions {
	options.ScriptID = core.StringPtr(scriptID)
	return options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (options *InstallVersionOptions) SetVersionLocatorID(versionLocatorID string) *InstallVersionOptions {
	options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return options
}

// SetVcenterID : Allow user to set VcenterID
func (options *InstallVersionOptions) SetVcenterID(vcenterID string) *InstallVersionOptions {
	options.VcenterID = core.StringPtr(vcenterID)
	return options
}

// SetVcenterUser : Allow user to set VcenterUser
func (options *InstallVersionOptions) SetVcenterUser(vcenterUser string) *InstallVersionOptions {
	options.VcenterUser = core.StringPtr(vcenterUser)
	return options
}

// SetVcenterPassword : Allow user to set VcenterPassword
func (options *InstallVersionOptions) SetVcenterPassword(vcenterPassword string) *InstallVersionOptions {
	options.VcenterPassword = core.StringPtr(vcenterPassword)
	return options
}

// SetVcenterLocation : Allow user to set VcenterLocation
func (options *InstallVersionOptions) SetVcenterLocation(vcenterLocation string) *InstallVersionOptions {
	options.VcenterLocation = core.StringPtr(vcenterLocation)
	return options
}

// SetVcenterDatastore : Allow user to set VcenterDatastore
func (options *InstallVersionOptions) SetVcenterDatastore(vcenterDatastore string) *InstallVersionOptions {
	options.VcenterDatastore = core.StringPtr(vcenterDatastore)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *InstallVersionOptions) SetHeaders(param map[string]string) *InstallVersionOptions {
	options.Headers = param
	return options
}

// Kind : Offering kind.
type Kind struct {
	// Unique ID.
	ID *string `json:"id,omitempty"`

	// content kind, e.g., helm, vm image.
	FormatKind *string `json:"format_kind,omitempty"`

	// target cloud to install, e.g., iks, open_shift_iks.
	TargetKind *string `json:"target_kind,omitempty"`

	// Open ended metadata information.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Installation instruction.
	InstallDescription *string `json:"install_description,omitempty"`

	// List of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// List of features associated with this offering.
	AdditionalFeatures []Feature `json:"additional_features,omitempty"`

	// The date and time this catalog was created.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The date and time this catalog was last updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// list of versions.
	Versions []Version `json:"versions,omitempty"`

	// list of plans.
	Plans []Plan `json:"plans,omitempty"`
}

// UnmarshalKind unmarshals an instance of Kind from the specified map of raw messages.
func UnmarshalKind(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Kind)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format_kind", &obj.FormatKind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_kind", &obj.TargetKind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "install_description", &obj.InstallDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "additional_features", &obj.AdditionalFeatures, UnmarshalFeature)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "plans", &obj.Plans, UnmarshalPlan)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// License : BSS license.
type License struct {
	// License ID.
	ID *string `json:"id,omitempty"`

	// license name.
	Name *string `json:"name,omitempty"`

	// type of license e.g., Apache xxx.
	Type *string `json:"type,omitempty"`

	// URL for the license text.
	URL *string `json:"url,omitempty"`

	// License description.
	Description *string `json:"description,omitempty"`
}

// UnmarshalLicense unmarshals an instance of License from the specified map of raw messages.
func UnmarshalLicense(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(License)
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
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
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

// ListCatalogsOptions : The ListCatalogs options.
type ListCatalogsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListCatalogsOptions : Instantiate ListCatalogsOptions
func (*CatalogManagementV1) NewListCatalogsOptions() *ListCatalogsOptions {
	return &ListCatalogsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListCatalogsOptions) SetHeaders(param map[string]string) *ListCatalogsOptions {
	options.Headers = param
	return options
}

// ListObjectsOptions : The ListObjects options.
type ListObjectsOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// The number of results to return.
	Limit *int64

	// The number of results to skip before returning values.
	Offset *int64

	// Only return results that contain the specified string.
	Name *string

	// The field on which the output is sorted. Sorts by default by **label** property. Available fields are **name**,
	// **label**, **created**, and **updated**. By adding **-** (i.e. **-label**) in front of the query string, you can
	// specify descending order. Default is ascending order.
	Sort *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListObjectsOptions : Instantiate ListObjectsOptions
func (*CatalogManagementV1) NewListObjectsOptions(catalogIdentifier string) *ListObjectsOptions {
	return &ListObjectsOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ListObjectsOptions) SetCatalogIdentifier(catalogIdentifier string) *ListObjectsOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListObjectsOptions) SetLimit(limit int64) *ListObjectsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListObjectsOptions) SetOffset(offset int64) *ListObjectsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetName : Allow user to set Name
func (options *ListObjectsOptions) SetName(name string) *ListObjectsOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListObjectsOptions) SetSort(sort string) *ListObjectsOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListObjectsOptions) SetHeaders(param map[string]string) *ListObjectsOptions {
	options.Headers = param
	return options
}

// ListOfferingsOptions : The ListOfferings options.
type ListOfferingsOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// true - Strip down the content of what is returned. For example don't return the readme. Makes the result much
	// smaller. Defaults to false.
	Digest *bool

	// The maximum number of results to return.
	Limit *int64

	// The number of results to skip before returning values.
	Offset *int64

	// Only return results that contain the specified string.
	Name *string

	// The field on which the output is sorted. Sorts by default by **label** property. Available fields are **name**,
	// **label**, **created**, and **updated**. By adding **-** (i.e. **-label**) in front of the query string, you can
	// specify descending order. Default is ascending order.
	Sort *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListOfferingsOptions : Instantiate ListOfferingsOptions
func (*CatalogManagementV1) NewListOfferingsOptions(catalogIdentifier string) *ListOfferingsOptions {
	return &ListOfferingsOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ListOfferingsOptions) SetCatalogIdentifier(catalogIdentifier string) *ListOfferingsOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetDigest : Allow user to set Digest
func (options *ListOfferingsOptions) SetDigest(digest bool) *ListOfferingsOptions {
	options.Digest = core.BoolPtr(digest)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListOfferingsOptions) SetLimit(limit int64) *ListOfferingsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListOfferingsOptions) SetOffset(offset int64) *ListOfferingsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetName : Allow user to set Name
func (options *ListOfferingsOptions) SetName(name string) *ListOfferingsOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListOfferingsOptions) SetSort(sort string) *ListOfferingsOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListOfferingsOptions) SetHeaders(param map[string]string) *ListOfferingsOptions {
	options.Headers = param
	return options
}

// ListOperatorsOptions : The ListOperators options.
type ListOperatorsOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Cluster identification.
	ClusterID *string `validate:"required"`

	// Cluster region.
	Region *string `validate:"required"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListOperatorsOptions : Instantiate ListOperatorsOptions
func (*CatalogManagementV1) NewListOperatorsOptions(xAuthRefreshToken string, clusterID string, region string, versionLocatorID string) *ListOperatorsOptions {
	return &ListOperatorsOptions{
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
		ClusterID: core.StringPtr(clusterID),
		Region: core.StringPtr(region),
		VersionLocatorID: core.StringPtr(versionLocatorID),
	}
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *ListOperatorsOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *ListOperatorsOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *ListOperatorsOptions) SetClusterID(clusterID string) *ListOperatorsOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *ListOperatorsOptions) SetRegion(region string) *ListOperatorsOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (options *ListOperatorsOptions) SetVersionLocatorID(versionLocatorID string) *ListOperatorsOptions {
	options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListOperatorsOptions) SetHeaders(param map[string]string) *ListOperatorsOptions {
	options.Headers = param
	return options
}

// NamespaceSearchResult : Paginated list of namespace search results.
type NamespaceSearchResult struct {
	// The offset (origin 0) of the first resource in this page of search results.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources returned in each page of search results.
	Limit *int64 `json:"limit" validate:"required"`

	// The overall total number of resources in the search result set.
	TotalCount *int64 `json:"total_count,omitempty"`

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

	// Resulting objects.
	Resources []string `json:"resources,omitempty"`
}

// UnmarshalNamespaceSearchResult unmarshals an instance of NamespaceSearchResult from the specified map of raw messages.
func UnmarshalNamespaceSearchResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NamespaceSearchResult)
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
	err = core.UnmarshalPrimitive(m, "resource_count", &obj.ResourceCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last", &obj.Last)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "prev", &obj.Prev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
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

// ObjectAccess : object access.
type ObjectAccess struct {
	// unique id.
	ID *string `json:"id,omitempty"`

	// account id.
	Account *string `json:"account,omitempty"`

	// unique id.
	CatalogID *string `json:"catalog_id,omitempty"`

	// object id.
	TargetID *string `json:"target_id,omitempty"`

	// date and time create.
	Create *strfmt.DateTime `json:"create,omitempty"`
}

// UnmarshalObjectAccess unmarshals an instance of ObjectAccess from the specified map of raw messages.
func UnmarshalObjectAccess(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ObjectAccess)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_id", &obj.TargetID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "create", &obj.Create)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ObjectAccessListResult : Paginated object search result.
type ObjectAccessListResult struct {
	// The offset (origin 0) of the first resource in this page of search results.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources returned in each page of search results.
	Limit *int64 `json:"limit" validate:"required"`

	// The overall total number of resources in the search result set.
	TotalCount *int64 `json:"total_count,omitempty"`

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

	// Resulting objects.
	Resources []ObjectAccess `json:"resources,omitempty"`
}

// UnmarshalObjectAccessListResult unmarshals an instance of ObjectAccessListResult from the specified map of raw messages.
func UnmarshalObjectAccessListResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ObjectAccessListResult)
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
	err = core.UnmarshalPrimitive(m, "resource_count", &obj.ResourceCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last", &obj.Last)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "prev", &obj.Prev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalObjectAccess)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ObjectListResult : Paginated object search result.
type ObjectListResult struct {
	// The offset (origin 0) of the first resource in this page of search results.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources returned in each page of search results.
	Limit *int64 `json:"limit" validate:"required"`

	// The overall total number of resources in the search result set.
	TotalCount *int64 `json:"total_count,omitempty"`

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

	// Resulting objects.
	Resources []CatalogObject `json:"resources,omitempty"`
}

// UnmarshalObjectListResult unmarshals an instance of ObjectListResult from the specified map of raw messages.
func UnmarshalObjectListResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ObjectListResult)
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
	err = core.UnmarshalPrimitive(m, "resource_count", &obj.ResourceCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last", &obj.Last)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "prev", &obj.Prev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalCatalogObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ObjectSearchResult : Paginated object search result.
type ObjectSearchResult struct {
	// The offset (origin 0) of the first resource in this page of search results.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources returned in each page of search results.
	Limit *int64 `json:"limit" validate:"required"`

	// The overall total number of resources in the search result set.
	TotalCount *int64 `json:"total_count,omitempty"`

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

	// Resulting objects.
	Resources []CatalogObject `json:"resources,omitempty"`
}

// UnmarshalObjectSearchResult unmarshals an instance of ObjectSearchResult from the specified map of raw messages.
func UnmarshalObjectSearchResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ObjectSearchResult)
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
	err = core.UnmarshalPrimitive(m, "resource_count", &obj.ResourceCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last", &obj.Last)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "prev", &obj.Prev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalCatalogObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Offering : Offering information.
type Offering struct {
	// unique id.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// The url for this specific offering.
	URL *string `json:"url,omitempty"`

	// The crn for this specific offering.
	CRN *string `json:"crn,omitempty"`

	// Display Name in the requested language.
	Label *string `json:"label,omitempty"`

	// The programmatic name of this offering.
	Name *string `json:"name,omitempty"`

	// URL for an icon associated with this offering.
	OfferingIconURL *string `json:"offering_icon_url,omitempty"`

	// URL for an additional docs with this offering.
	OfferingDocsURL *string `json:"offering_docs_url,omitempty"`

	// URL to be displayed in the Consumption UI for getting support on this offering.
	OfferingSupportURL *string `json:"offering_support_url,omitempty"`

	// List of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// List of keywords associated with offering, typically used to search for it.
	Keywords []string `json:"keywords,omitempty"`

	// Repository info for offerings.
	Rating *Rating `json:"rating,omitempty"`

	// The date and time this catalog was created.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The date and time this catalog was last updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Short description in the requested language.
	ShortDescription *string `json:"short_description,omitempty"`

	// Long description in the requested language.
	LongDescription *string `json:"long_description,omitempty"`

	// list of features associated with this offering.
	Features []Feature `json:"features,omitempty"`

	// Array of kind.
	Kinds []Kind `json:"kinds,omitempty"`

	// Is it permitted to request publishing to IBM or Public.
	PermitRequestIBMPublicPublish *bool `json:"permit_request_ibm_public_publish,omitempty"`

	// Indicates if this offering has been approved for use by all IBMers.
	IBMPublishApproved *bool `json:"ibm_publish_approved,omitempty"`

	// Indicates if this offering has been approved for use by all IBM Cloud users.
	PublicPublishApproved *bool `json:"public_publish_approved,omitempty"`

	// The original offering CRN that this publish entry came from.
	PublicOriginalCRN *string `json:"public_original_crn,omitempty"`

	// The crn of the public catalog entry of this offering.
	PublishPublicCRN *string `json:"publish_public_crn,omitempty"`

	// The portal's approval record ID.
	PortalApprovalRecord *string `json:"portal_approval_record,omitempty"`

	// The portal UI URL.
	PortalUIURL *string `json:"portal_ui_url,omitempty"`

	// The id of the catalog containing this offering.
	CatalogID *string `json:"catalog_id,omitempty"`

	// The name of the catalog.
	CatalogName *string `json:"catalog_name,omitempty"`

	// Map of metadata values for this offering.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// A disclaimer for this offering.
	Disclaimer *string `json:"disclaimer,omitempty"`

	// Determine if this offering should be displayed in the Consumption UI.
	Hidden *bool `json:"hidden,omitempty"`

	// Provider of this offering.
	Provider *string `json:"provider,omitempty"`

	// Repository info for offerings.
	RepoInfo *RepoInfo `json:"repo_info,omitempty"`
}

// UnmarshalOffering unmarshals an instance of Offering from the specified map of raw messages.
func UnmarshalOffering(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Offering)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_icon_url", &obj.OfferingIconURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_docs_url", &obj.OfferingDocsURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_support_url", &obj.OfferingSupportURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "keywords", &obj.Keywords)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rating", &obj.Rating, UnmarshalRating)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "short_description", &obj.ShortDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeature)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "kinds", &obj.Kinds, UnmarshalKind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "permit_request_ibm_public_publish", &obj.PermitRequestIBMPublicPublish)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_publish_approved", &obj.IBMPublishApproved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "public_publish_approved", &obj.PublicPublishApproved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "public_original_crn", &obj.PublicOriginalCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "publish_public_crn", &obj.PublishPublicCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "portal_approval_record", &obj.PortalApprovalRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "portal_ui_url", &obj.PortalUIURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_name", &obj.CatalogName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disclaimer", &obj.Disclaimer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "repo_info", &obj.RepoInfo, UnmarshalRepoInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OfferingInstance : A offering instance resource (provision instance of a catalog offering).
type OfferingInstance struct {
	// provisioned instance ID (part of the CRN).
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// url reference to this object.
	URL *string `json:"url,omitempty"`

	// platform CRN for this instance.
	CRN *string `json:"crn,omitempty"`

	// the label for this instance.
	Label *string `json:"label,omitempty"`

	// Catalog ID this instance was created from.
	CatalogID *string `json:"catalog_id,omitempty"`

	// Offering ID this instance was created from.
	OfferingID *string `json:"offering_id,omitempty"`

	// the format this instance has (helm, operator, ova...).
	KindFormat *string `json:"kind_format,omitempty"`

	// The version this instance was installed from (not version id).
	Version *string `json:"version,omitempty"`

	// Cluster ID.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster region (e.g., us-south).
	ClusterRegion *string `json:"cluster_region,omitempty"`

	// List of target namespaces to install into.
	ClusterNamespaces []string `json:"cluster_namespaces,omitempty"`

	// designate to install into all namespaces.
	ClusterAllNamespaces *bool `json:"cluster_all_namespaces,omitempty"`

	// Id of the schematics workspace, for offering instances provisioned through schematics.
	SchematicsWorkspaceID *string `json:"schematics_workspace_id,omitempty"`

	// Id of the resource group to provision the offering instance into.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// Type of install plan (also known as approval strategy) for operator subscriptions. Can be either automatic, which
	// automatically upgrades operators to the latest in a channel, or manual, which requires approval on the cluster.
	InstallPlan *string `json:"install_plan,omitempty"`

	// Channel to pin the operator subscription to.
	Channel *string `json:"channel,omitempty"`

	// Map of metadata values for this offering instance.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// the last operation performed and status.
	LastOperation *OfferingInstanceLastOperation `json:"last_operation,omitempty"`
}

// UnmarshalOfferingInstance unmarshals an instance of OfferingInstance from the specified map of raw messages.
func UnmarshalOfferingInstance(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OfferingInstance)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_id", &obj.OfferingID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind_format", &obj.KindFormat)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_id", &obj.ClusterID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_region", &obj.ClusterRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_namespaces", &obj.ClusterNamespaces)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_all_namespaces", &obj.ClusterAllNamespaces)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schematics_workspace_id", &obj.SchematicsWorkspaceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "install_plan", &obj.InstallPlan)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "channel", &obj.Channel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_operation", &obj.LastOperation, UnmarshalOfferingInstanceLastOperation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OfferingInstanceLastOperation : the last operation performed and status.
type OfferingInstanceLastOperation struct {
	// last operation performed.
	Operation *string `json:"operation,omitempty"`

	// state after the last operation performed.
	State *string `json:"state,omitempty"`

	// additional information about the last operation.
	Message *string `json:"message,omitempty"`

	// transaction id from the last operation.
	TransactionID *string `json:"transaction_id,omitempty"`

	// Date and time last updated.
	Updated *string `json:"updated,omitempty"`
}

// UnmarshalOfferingInstanceLastOperation unmarshals an instance of OfferingInstanceLastOperation from the specified map of raw messages.
func UnmarshalOfferingInstanceLastOperation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OfferingInstanceLastOperation)
	err = core.UnmarshalPrimitive(m, "operation", &obj.Operation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "transaction_id", &obj.TransactionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OfferingSearchResult : Paginated offering search result.
type OfferingSearchResult struct {
	// The offset (origin 0) of the first resource in this page of search results.
	Offset *int64 `json:"offset" validate:"required"`

	// The maximum number of resources returned in each page of search results.
	Limit *int64 `json:"limit" validate:"required"`

	// The overall total number of resources in the search result set.
	TotalCount *int64 `json:"total_count,omitempty"`

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

	// Resulting objects.
	Resources []Offering `json:"resources,omitempty"`
}

// UnmarshalOfferingSearchResult unmarshals an instance of OfferingSearchResult from the specified map of raw messages.
func UnmarshalOfferingSearchResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OfferingSearchResult)
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
	err = core.UnmarshalPrimitive(m, "resource_count", &obj.ResourceCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "first", &obj.First)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last", &obj.Last)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "prev", &obj.Prev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next", &obj.Next)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalOffering)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OperatorDeployResult : Operator deploy result.
type OperatorDeployResult struct {
	// Status phase.
	Phase *string `json:"phase,omitempty"`

	// Status message.
	Message *string `json:"message,omitempty"`

	// Operator API path.
	Link *string `json:"link,omitempty"`

	// Name of Operator.
	Name *string `json:"name,omitempty"`

	// Operator version.
	Version *string `json:"version,omitempty"`

	// Kube namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Package Operator exists in.
	PackageName *string `json:"package_name,omitempty"`

	// Catalog identification.
	CatalogID *string `json:"catalog_id,omitempty"`
}

// UnmarshalOperatorDeployResult unmarshals an instance of OperatorDeployResult from the specified map of raw messages.
func UnmarshalOperatorDeployResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OperatorDeployResult)
	err = core.UnmarshalPrimitive(m, "phase", &obj.Phase)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "link", &obj.Link)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "package_name", &obj.PackageName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Plan : Offering plan.
type Plan struct {
	// unique id.
	ID *string `json:"id,omitempty"`

	// Display Name in the requested language.
	Label *string `json:"label,omitempty"`

	// The programmatic name of this offering.
	Name *string `json:"name,omitempty"`

	// Short description in the requested language.
	ShortDescription *string `json:"short_description,omitempty"`

	// Long description in the requested language.
	LongDescription *string `json:"long_description,omitempty"`

	// open ended metadata information.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// list of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// list of features associated with this offering.
	AdditionalFeatures []Feature `json:"additional_features,omitempty"`

	// the date'time this catalog was created.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// the date'time this catalog was last updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// list of deployments.
	Deployments []Deployment `json:"deployments,omitempty"`
}

// UnmarshalPlan unmarshals an instance of Plan from the specified map of raw messages.
func UnmarshalPlan(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Plan)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "short_description", &obj.ShortDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "additional_features", &obj.AdditionalFeatures, UnmarshalFeature)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "deployments", &obj.Deployments, UnmarshalDeployment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PreinstallVersionOptions : The PreinstallVersion options.
type PreinstallVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Cluster ID.
	ClusterID *string

	// Cluster region.
	Region *string

	// Kube namespace.
	Namespace *string

	// Object containing Helm chart override values.  To use a secret for items of type password, specify a JSON encoded
	// value of $ref:#/components/schemas/SecretInstance, prefixed with `cmsm_v1:`.
	OverrideValues map[string]interface{}

	// Entitlement API Key for this offering.
	EntitlementApikey *string

	// Schematics workspace configuration.
	Schematics *DeployRequestBodySchematics

	// Script.
	Script *string

	// Script ID.
	ScriptID *string

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string

	// VCenter ID.
	VcenterID *string

	// VCenter User.
	VcenterUser *string

	// VCenter Password.
	VcenterPassword *string

	// VCenter Location.
	VcenterLocation *string

	// VCenter Datastore.
	VcenterDatastore *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPreinstallVersionOptions : Instantiate PreinstallVersionOptions
func (*CatalogManagementV1) NewPreinstallVersionOptions(versionLocID string, xAuthRefreshToken string) *PreinstallVersionOptions {
	return &PreinstallVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *PreinstallVersionOptions) SetVersionLocID(versionLocID string) *PreinstallVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *PreinstallVersionOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *PreinstallVersionOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *PreinstallVersionOptions) SetClusterID(clusterID string) *PreinstallVersionOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *PreinstallVersionOptions) SetRegion(region string) *PreinstallVersionOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetNamespace : Allow user to set Namespace
func (options *PreinstallVersionOptions) SetNamespace(namespace string) *PreinstallVersionOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetOverrideValues : Allow user to set OverrideValues
func (options *PreinstallVersionOptions) SetOverrideValues(overrideValues map[string]interface{}) *PreinstallVersionOptions {
	options.OverrideValues = overrideValues
	return options
}

// SetEntitlementApikey : Allow user to set EntitlementApikey
func (options *PreinstallVersionOptions) SetEntitlementApikey(entitlementApikey string) *PreinstallVersionOptions {
	options.EntitlementApikey = core.StringPtr(entitlementApikey)
	return options
}

// SetSchematics : Allow user to set Schematics
func (options *PreinstallVersionOptions) SetSchematics(schematics *DeployRequestBodySchematics) *PreinstallVersionOptions {
	options.Schematics = schematics
	return options
}

// SetScript : Allow user to set Script
func (options *PreinstallVersionOptions) SetScript(script string) *PreinstallVersionOptions {
	options.Script = core.StringPtr(script)
	return options
}

// SetScriptID : Allow user to set ScriptID
func (options *PreinstallVersionOptions) SetScriptID(scriptID string) *PreinstallVersionOptions {
	options.ScriptID = core.StringPtr(scriptID)
	return options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (options *PreinstallVersionOptions) SetVersionLocatorID(versionLocatorID string) *PreinstallVersionOptions {
	options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return options
}

// SetVcenterID : Allow user to set VcenterID
func (options *PreinstallVersionOptions) SetVcenterID(vcenterID string) *PreinstallVersionOptions {
	options.VcenterID = core.StringPtr(vcenterID)
	return options
}

// SetVcenterUser : Allow user to set VcenterUser
func (options *PreinstallVersionOptions) SetVcenterUser(vcenterUser string) *PreinstallVersionOptions {
	options.VcenterUser = core.StringPtr(vcenterUser)
	return options
}

// SetVcenterPassword : Allow user to set VcenterPassword
func (options *PreinstallVersionOptions) SetVcenterPassword(vcenterPassword string) *PreinstallVersionOptions {
	options.VcenterPassword = core.StringPtr(vcenterPassword)
	return options
}

// SetVcenterLocation : Allow user to set VcenterLocation
func (options *PreinstallVersionOptions) SetVcenterLocation(vcenterLocation string) *PreinstallVersionOptions {
	options.VcenterLocation = core.StringPtr(vcenterLocation)
	return options
}

// SetVcenterDatastore : Allow user to set VcenterDatastore
func (options *PreinstallVersionOptions) SetVcenterDatastore(vcenterDatastore string) *PreinstallVersionOptions {
	options.VcenterDatastore = core.StringPtr(vcenterDatastore)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PreinstallVersionOptions) SetHeaders(param map[string]string) *PreinstallVersionOptions {
	options.Headers = param
	return options
}

// PublicPublishObjectOptions : The PublicPublishObject options.
type PublicPublishObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPublicPublishObjectOptions : Instantiate PublicPublishObjectOptions
func (*CatalogManagementV1) NewPublicPublishObjectOptions(catalogIdentifier string, objectIdentifier string) *PublicPublishObjectOptions {
	return &PublicPublishObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *PublicPublishObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *PublicPublishObjectOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *PublicPublishObjectOptions) SetObjectIdentifier(objectIdentifier string) *PublicPublishObjectOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PublicPublishObjectOptions) SetHeaders(param map[string]string) *PublicPublishObjectOptions {
	options.Headers = param
	return options
}

// PublicPublishVersionOptions : The PublicPublishVersion options.
type PublicPublishVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPublicPublishVersionOptions : Instantiate PublicPublishVersionOptions
func (*CatalogManagementV1) NewPublicPublishVersionOptions(versionLocID string) *PublicPublishVersionOptions {
	return &PublicPublishVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *PublicPublishVersionOptions) SetVersionLocID(versionLocID string) *PublicPublishVersionOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PublicPublishVersionOptions) SetHeaders(param map[string]string) *PublicPublishVersionOptions {
	options.Headers = param
	return options
}

// PublishObject : Publish information.
type PublishObject struct {
	// Is it permitted to request publishing to IBM or Public.
	PermitIBMPublicPublish *bool `json:"permit_ibm_public_publish,omitempty"`

	// Indicates if this offering has been approved for use by all IBMers.
	IBMApproved *bool `json:"ibm_approved,omitempty"`

	// Indicates if this offering has been approved for use by all IBM Cloud users.
	PublicApproved *bool `json:"public_approved,omitempty"`

	// The portal's approval record ID.
	PortalApprovalRecord *string `json:"portal_approval_record,omitempty"`

	// The portal UI URL.
	PortalURL *string `json:"portal_url,omitempty"`
}

// UnmarshalPublishObject unmarshals an instance of PublishObject from the specified map of raw messages.
func UnmarshalPublishObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PublishObject)
	err = core.UnmarshalPrimitive(m, "permit_ibm_public_publish", &obj.PermitIBMPublicPublish)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibm_approved", &obj.IBMApproved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "public_approved", &obj.PublicApproved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "portal_approval_record", &obj.PortalApprovalRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "portal_url", &obj.PortalURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PutOfferingInstanceOptions : The PutOfferingInstance options.
type PutOfferingInstanceOptions struct {
	// Version Instance identifier.
	InstanceIdentifier *string `validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// provisioned instance ID (part of the CRN).
	ID *string

	// Cloudant revision.
	Rev *string

	// url reference to this object.
	URL *string

	// platform CRN for this instance.
	CRN *string

	// the label for this instance.
	Label *string

	// Catalog ID this instance was created from.
	CatalogID *string

	// Offering ID this instance was created from.
	OfferingID *string

	// the format this instance has (helm, operator, ova...).
	KindFormat *string

	// The version this instance was installed from (not version id).
	Version *string

	// Cluster ID.
	ClusterID *string

	// Cluster region (e.g., us-south).
	ClusterRegion *string

	// List of target namespaces to install into.
	ClusterNamespaces []string

	// designate to install into all namespaces.
	ClusterAllNamespaces *bool

	// Id of the schematics workspace, for offering instances provisioned through schematics.
	SchematicsWorkspaceID *string

	// Id of the resource group to provision the offering instance into.
	ResourceGroupID *string

	// Type of install plan (also known as approval strategy) for operator subscriptions. Can be either automatic, which
	// automatically upgrades operators to the latest in a channel, or manual, which requires approval on the cluster.
	InstallPlan *string

	// Channel to pin the operator subscription to.
	Channel *string

	// Map of metadata values for this offering instance.
	Metadata map[string]interface{}

	// the last operation performed and status.
	LastOperation *OfferingInstanceLastOperation

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPutOfferingInstanceOptions : Instantiate PutOfferingInstanceOptions
func (*CatalogManagementV1) NewPutOfferingInstanceOptions(instanceIdentifier string, xAuthRefreshToken string) *PutOfferingInstanceOptions {
	return &PutOfferingInstanceOptions{
		InstanceIdentifier: core.StringPtr(instanceIdentifier),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetInstanceIdentifier : Allow user to set InstanceIdentifier
func (options *PutOfferingInstanceOptions) SetInstanceIdentifier(instanceIdentifier string) *PutOfferingInstanceOptions {
	options.InstanceIdentifier = core.StringPtr(instanceIdentifier)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *PutOfferingInstanceOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *PutOfferingInstanceOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetID : Allow user to set ID
func (options *PutOfferingInstanceOptions) SetID(id string) *PutOfferingInstanceOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetRev : Allow user to set Rev
func (options *PutOfferingInstanceOptions) SetRev(rev string) *PutOfferingInstanceOptions {
	options.Rev = core.StringPtr(rev)
	return options
}

// SetURL : Allow user to set URL
func (options *PutOfferingInstanceOptions) SetURL(url string) *PutOfferingInstanceOptions {
	options.URL = core.StringPtr(url)
	return options
}

// SetCRN : Allow user to set CRN
func (options *PutOfferingInstanceOptions) SetCRN(crn string) *PutOfferingInstanceOptions {
	options.CRN = core.StringPtr(crn)
	return options
}

// SetLabel : Allow user to set Label
func (options *PutOfferingInstanceOptions) SetLabel(label string) *PutOfferingInstanceOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetCatalogID : Allow user to set CatalogID
func (options *PutOfferingInstanceOptions) SetCatalogID(catalogID string) *PutOfferingInstanceOptions {
	options.CatalogID = core.StringPtr(catalogID)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *PutOfferingInstanceOptions) SetOfferingID(offeringID string) *PutOfferingInstanceOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetKindFormat : Allow user to set KindFormat
func (options *PutOfferingInstanceOptions) SetKindFormat(kindFormat string) *PutOfferingInstanceOptions {
	options.KindFormat = core.StringPtr(kindFormat)
	return options
}

// SetVersion : Allow user to set Version
func (options *PutOfferingInstanceOptions) SetVersion(version string) *PutOfferingInstanceOptions {
	options.Version = core.StringPtr(version)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *PutOfferingInstanceOptions) SetClusterID(clusterID string) *PutOfferingInstanceOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetClusterRegion : Allow user to set ClusterRegion
func (options *PutOfferingInstanceOptions) SetClusterRegion(clusterRegion string) *PutOfferingInstanceOptions {
	options.ClusterRegion = core.StringPtr(clusterRegion)
	return options
}

// SetClusterNamespaces : Allow user to set ClusterNamespaces
func (options *PutOfferingInstanceOptions) SetClusterNamespaces(clusterNamespaces []string) *PutOfferingInstanceOptions {
	options.ClusterNamespaces = clusterNamespaces
	return options
}

// SetClusterAllNamespaces : Allow user to set ClusterAllNamespaces
func (options *PutOfferingInstanceOptions) SetClusterAllNamespaces(clusterAllNamespaces bool) *PutOfferingInstanceOptions {
	options.ClusterAllNamespaces = core.BoolPtr(clusterAllNamespaces)
	return options
}

// SetSchematicsWorkspaceID : Allow user to set SchematicsWorkspaceID
func (options *PutOfferingInstanceOptions) SetSchematicsWorkspaceID(schematicsWorkspaceID string) *PutOfferingInstanceOptions {
	options.SchematicsWorkspaceID = core.StringPtr(schematicsWorkspaceID)
	return options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (options *PutOfferingInstanceOptions) SetResourceGroupID(resourceGroupID string) *PutOfferingInstanceOptions {
	options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return options
}

// SetInstallPlan : Allow user to set InstallPlan
func (options *PutOfferingInstanceOptions) SetInstallPlan(installPlan string) *PutOfferingInstanceOptions {
	options.InstallPlan = core.StringPtr(installPlan)
	return options
}

// SetChannel : Allow user to set Channel
func (options *PutOfferingInstanceOptions) SetChannel(channel string) *PutOfferingInstanceOptions {
	options.Channel = core.StringPtr(channel)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *PutOfferingInstanceOptions) SetMetadata(metadata map[string]interface{}) *PutOfferingInstanceOptions {
	options.Metadata = metadata
	return options
}

// SetLastOperation : Allow user to set LastOperation
func (options *PutOfferingInstanceOptions) SetLastOperation(lastOperation *OfferingInstanceLastOperation) *PutOfferingInstanceOptions {
	options.LastOperation = lastOperation
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PutOfferingInstanceOptions) SetHeaders(param map[string]string) *PutOfferingInstanceOptions {
	options.Headers = param
	return options
}

// Rating : Repository info for offerings.
type Rating struct {
	// One start rating.
	OneStarCount *int64 `json:"one_star_count,omitempty"`

	// Two start rating.
	TwoStarCount *int64 `json:"two_star_count,omitempty"`

	// Three start rating.
	ThreeStarCount *int64 `json:"three_star_count,omitempty"`

	// Four start rating.
	FourStarCount *int64 `json:"four_star_count,omitempty"`
}

// UnmarshalRating unmarshals an instance of Rating from the specified map of raw messages.
func UnmarshalRating(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rating)
	err = core.UnmarshalPrimitive(m, "one_star_count", &obj.OneStarCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "two_star_count", &obj.TwoStarCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "three_star_count", &obj.ThreeStarCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "four_star_count", &obj.FourStarCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReloadOfferingOptions : The ReloadOffering options.
type ReloadOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identification.
	OfferingID *string `validate:"required,ne="`

	// The semver value for this new version.
	TargetVersion *string `validate:"required"`

	// Tags array.
	Tags []string

	// Target kinds.  Current valid values are 'iks', 'roks', 'vcenter', and 'terraform'.
	TargetKinds []string

	// byte array representing the content to be imported.  Only supported for OVA images at this time.
	Content *[]byte

	// URL path to zip location.  If not specified, must provide content in this post body.
	Zipurl *string

	// The type of repository containing this version.  Valid values are 'public_git' or 'enterprise_git'.
	RepoType *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReloadOfferingOptions : Instantiate ReloadOfferingOptions
func (*CatalogManagementV1) NewReloadOfferingOptions(catalogIdentifier string, offeringID string, targetVersion string) *ReloadOfferingOptions {
	return &ReloadOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		TargetVersion: core.StringPtr(targetVersion),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ReloadOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *ReloadOfferingOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *ReloadOfferingOptions) SetOfferingID(offeringID string) *ReloadOfferingOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetTargetVersion : Allow user to set TargetVersion
func (options *ReloadOfferingOptions) SetTargetVersion(targetVersion string) *ReloadOfferingOptions {
	options.TargetVersion = core.StringPtr(targetVersion)
	return options
}

// SetTags : Allow user to set Tags
func (options *ReloadOfferingOptions) SetTags(tags []string) *ReloadOfferingOptions {
	options.Tags = tags
	return options
}

// SetTargetKinds : Allow user to set TargetKinds
func (options *ReloadOfferingOptions) SetTargetKinds(targetKinds []string) *ReloadOfferingOptions {
	options.TargetKinds = targetKinds
	return options
}

// SetContent : Allow user to set Content
func (options *ReloadOfferingOptions) SetContent(content []byte) *ReloadOfferingOptions {
	options.Content = &content
	return options
}

// SetZipurl : Allow user to set Zipurl
func (options *ReloadOfferingOptions) SetZipurl(zipurl string) *ReloadOfferingOptions {
	options.Zipurl = core.StringPtr(zipurl)
	return options
}

// SetRepoType : Allow user to set RepoType
func (options *ReloadOfferingOptions) SetRepoType(repoType string) *ReloadOfferingOptions {
	options.RepoType = core.StringPtr(repoType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReloadOfferingOptions) SetHeaders(param map[string]string) *ReloadOfferingOptions {
	options.Headers = param
	return options
}

// ReplaceCatalogOptions : The ReplaceCatalog options.
type ReplaceCatalogOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Unique ID.
	ID *string

	// Cloudant revision.
	Rev *string

	// Display Name in the requested language.
	Label *string

	// Description in the requested language.
	ShortDescription *string

	// URL for an icon associated with this catalog.
	CatalogIconURL *string

	// List of tags associated with this catalog.
	Tags []string

	// List of features associated with this catalog.
	Features []Feature

	// Denotes whether a catalog is disabled.
	Disabled *bool

	// Resource group id the catalog is owned by.
	ResourceGroupID *string

	// Account that owns catalog.
	OwningAccount *string

	// Filters for account and catalog filters.
	CatalogFilters *Filters

	// Feature information.
	SyndicationSettings *SyndicationResource

	// Kind of catalog. Supported kinds are offering and vpe.
	Kind *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceCatalogOptions : Instantiate ReplaceCatalogOptions
func (*CatalogManagementV1) NewReplaceCatalogOptions(catalogIdentifier string) *ReplaceCatalogOptions {
	return &ReplaceCatalogOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ReplaceCatalogOptions) SetCatalogIdentifier(catalogIdentifier string) *ReplaceCatalogOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetID : Allow user to set ID
func (options *ReplaceCatalogOptions) SetID(id string) *ReplaceCatalogOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetRev : Allow user to set Rev
func (options *ReplaceCatalogOptions) SetRev(rev string) *ReplaceCatalogOptions {
	options.Rev = core.StringPtr(rev)
	return options
}

// SetLabel : Allow user to set Label
func (options *ReplaceCatalogOptions) SetLabel(label string) *ReplaceCatalogOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetShortDescription : Allow user to set ShortDescription
func (options *ReplaceCatalogOptions) SetShortDescription(shortDescription string) *ReplaceCatalogOptions {
	options.ShortDescription = core.StringPtr(shortDescription)
	return options
}

// SetCatalogIconURL : Allow user to set CatalogIconURL
func (options *ReplaceCatalogOptions) SetCatalogIconURL(catalogIconURL string) *ReplaceCatalogOptions {
	options.CatalogIconURL = core.StringPtr(catalogIconURL)
	return options
}

// SetTags : Allow user to set Tags
func (options *ReplaceCatalogOptions) SetTags(tags []string) *ReplaceCatalogOptions {
	options.Tags = tags
	return options
}

// SetFeatures : Allow user to set Features
func (options *ReplaceCatalogOptions) SetFeatures(features []Feature) *ReplaceCatalogOptions {
	options.Features = features
	return options
}

// SetDisabled : Allow user to set Disabled
func (options *ReplaceCatalogOptions) SetDisabled(disabled bool) *ReplaceCatalogOptions {
	options.Disabled = core.BoolPtr(disabled)
	return options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (options *ReplaceCatalogOptions) SetResourceGroupID(resourceGroupID string) *ReplaceCatalogOptions {
	options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return options
}

// SetOwningAccount : Allow user to set OwningAccount
func (options *ReplaceCatalogOptions) SetOwningAccount(owningAccount string) *ReplaceCatalogOptions {
	options.OwningAccount = core.StringPtr(owningAccount)
	return options
}

// SetCatalogFilters : Allow user to set CatalogFilters
func (options *ReplaceCatalogOptions) SetCatalogFilters(catalogFilters *Filters) *ReplaceCatalogOptions {
	options.CatalogFilters = catalogFilters
	return options
}

// SetSyndicationSettings : Allow user to set SyndicationSettings
func (options *ReplaceCatalogOptions) SetSyndicationSettings(syndicationSettings *SyndicationResource) *ReplaceCatalogOptions {
	options.SyndicationSettings = syndicationSettings
	return options
}

// SetKind : Allow user to set Kind
func (options *ReplaceCatalogOptions) SetKind(kind string) *ReplaceCatalogOptions {
	options.Kind = core.StringPtr(kind)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceCatalogOptions) SetHeaders(param map[string]string) *ReplaceCatalogOptions {
	options.Headers = param
	return options
}

// ReplaceObjectOptions : The ReplaceObject options.
type ReplaceObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// unique id.
	ID *string

	// The programmatic name of this offering.
	Name *string

	// Cloudant revision.
	Rev *string

	// The crn for this specific object.
	CRN *string

	// The url for this specific object.
	URL *string

	// The parent for this specific object.
	ParentID *string

	// Translated display name in the requested language.
	LabelI18n *string

	// Display name in the requested language.
	Label *string

	// List of tags associated with this catalog.
	Tags []string

	// The date and time this catalog was created.
	Created *strfmt.DateTime

	// The date and time this catalog was last updated.
	Updated *strfmt.DateTime

	// Short description in the requested language.
	ShortDescription *string

	// Short description translation.
	ShortDescriptionI18n *string

	// Kind of object.
	Kind *string

	// Publish information.
	Publish *PublishObject

	// Offering state.
	State *State

	// The id of the catalog containing this offering.
	CatalogID *string

	// The name of the catalog.
	CatalogName *string

	// Map of data values for this object.
	Data map[string]interface{}

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceObjectOptions : Instantiate ReplaceObjectOptions
func (*CatalogManagementV1) NewReplaceObjectOptions(catalogIdentifier string, objectIdentifier string) *ReplaceObjectOptions {
	return &ReplaceObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ReplaceObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *ReplaceObjectOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *ReplaceObjectOptions) SetObjectIdentifier(objectIdentifier string) *ReplaceObjectOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetID : Allow user to set ID
func (options *ReplaceObjectOptions) SetID(id string) *ReplaceObjectOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetName : Allow user to set Name
func (options *ReplaceObjectOptions) SetName(name string) *ReplaceObjectOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetRev : Allow user to set Rev
func (options *ReplaceObjectOptions) SetRev(rev string) *ReplaceObjectOptions {
	options.Rev = core.StringPtr(rev)
	return options
}

// SetCRN : Allow user to set CRN
func (options *ReplaceObjectOptions) SetCRN(crn string) *ReplaceObjectOptions {
	options.CRN = core.StringPtr(crn)
	return options
}

// SetURL : Allow user to set URL
func (options *ReplaceObjectOptions) SetURL(url string) *ReplaceObjectOptions {
	options.URL = core.StringPtr(url)
	return options
}

// SetParentID : Allow user to set ParentID
func (options *ReplaceObjectOptions) SetParentID(parentID string) *ReplaceObjectOptions {
	options.ParentID = core.StringPtr(parentID)
	return options
}

// SetLabelI18n : Allow user to set LabelI18n
func (options *ReplaceObjectOptions) SetLabelI18n(labelI18n string) *ReplaceObjectOptions {
	options.LabelI18n = core.StringPtr(labelI18n)
	return options
}

// SetLabel : Allow user to set Label
func (options *ReplaceObjectOptions) SetLabel(label string) *ReplaceObjectOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetTags : Allow user to set Tags
func (options *ReplaceObjectOptions) SetTags(tags []string) *ReplaceObjectOptions {
	options.Tags = tags
	return options
}

// SetCreated : Allow user to set Created
func (options *ReplaceObjectOptions) SetCreated(created *strfmt.DateTime) *ReplaceObjectOptions {
	options.Created = created
	return options
}

// SetUpdated : Allow user to set Updated
func (options *ReplaceObjectOptions) SetUpdated(updated *strfmt.DateTime) *ReplaceObjectOptions {
	options.Updated = updated
	return options
}

// SetShortDescription : Allow user to set ShortDescription
func (options *ReplaceObjectOptions) SetShortDescription(shortDescription string) *ReplaceObjectOptions {
	options.ShortDescription = core.StringPtr(shortDescription)
	return options
}

// SetShortDescriptionI18n : Allow user to set ShortDescriptionI18n
func (options *ReplaceObjectOptions) SetShortDescriptionI18n(shortDescriptionI18n string) *ReplaceObjectOptions {
	options.ShortDescriptionI18n = core.StringPtr(shortDescriptionI18n)
	return options
}

// SetKind : Allow user to set Kind
func (options *ReplaceObjectOptions) SetKind(kind string) *ReplaceObjectOptions {
	options.Kind = core.StringPtr(kind)
	return options
}

// SetPublish : Allow user to set Publish
func (options *ReplaceObjectOptions) SetPublish(publish *PublishObject) *ReplaceObjectOptions {
	options.Publish = publish
	return options
}

// SetState : Allow user to set State
func (options *ReplaceObjectOptions) SetState(state *State) *ReplaceObjectOptions {
	options.State = state
	return options
}

// SetCatalogID : Allow user to set CatalogID
func (options *ReplaceObjectOptions) SetCatalogID(catalogID string) *ReplaceObjectOptions {
	options.CatalogID = core.StringPtr(catalogID)
	return options
}

// SetCatalogName : Allow user to set CatalogName
func (options *ReplaceObjectOptions) SetCatalogName(catalogName string) *ReplaceObjectOptions {
	options.CatalogName = core.StringPtr(catalogName)
	return options
}

// SetData : Allow user to set Data
func (options *ReplaceObjectOptions) SetData(data map[string]interface{}) *ReplaceObjectOptions {
	options.Data = data
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceObjectOptions) SetHeaders(param map[string]string) *ReplaceObjectOptions {
	options.Headers = param
	return options
}

// ReplaceOfferingIconOptions : The ReplaceOfferingIcon options.
type ReplaceOfferingIconOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identification.
	OfferingID *string `validate:"required,ne="`

	// Name of the file name that is being uploaded.
	FileName *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceOfferingIconOptions : Instantiate ReplaceOfferingIconOptions
func (*CatalogManagementV1) NewReplaceOfferingIconOptions(catalogIdentifier string, offeringID string, fileName string) *ReplaceOfferingIconOptions {
	return &ReplaceOfferingIconOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		FileName: core.StringPtr(fileName),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ReplaceOfferingIconOptions) SetCatalogIdentifier(catalogIdentifier string) *ReplaceOfferingIconOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *ReplaceOfferingIconOptions) SetOfferingID(offeringID string) *ReplaceOfferingIconOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetFileName : Allow user to set FileName
func (options *ReplaceOfferingIconOptions) SetFileName(fileName string) *ReplaceOfferingIconOptions {
	options.FileName = core.StringPtr(fileName)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceOfferingIconOptions) SetHeaders(param map[string]string) *ReplaceOfferingIconOptions {
	options.Headers = param
	return options
}

// ReplaceOfferingOptions : The ReplaceOffering options.
type ReplaceOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identification.
	OfferingID *string `validate:"required,ne="`

	// unique id.
	ID *string

	// Cloudant revision.
	Rev *string

	// The url for this specific offering.
	URL *string

	// The crn for this specific offering.
	CRN *string

	// Display Name in the requested language.
	Label *string

	// The programmatic name of this offering.
	Name *string

	// URL for an icon associated with this offering.
	OfferingIconURL *string

	// URL for an additional docs with this offering.
	OfferingDocsURL *string

	// URL to be displayed in the Consumption UI for getting support on this offering.
	OfferingSupportURL *string

	// List of tags associated with this catalog.
	Tags []string

	// List of keywords associated with offering, typically used to search for it.
	Keywords []string

	// Repository info for offerings.
	Rating *Rating

	// The date and time this catalog was created.
	Created *strfmt.DateTime

	// The date and time this catalog was last updated.
	Updated *strfmt.DateTime

	// Short description in the requested language.
	ShortDescription *string

	// Long description in the requested language.
	LongDescription *string

	// list of features associated with this offering.
	Features []Feature

	// Array of kind.
	Kinds []Kind

	// Is it permitted to request publishing to IBM or Public.
	PermitRequestIBMPublicPublish *bool

	// Indicates if this offering has been approved for use by all IBMers.
	IBMPublishApproved *bool

	// Indicates if this offering has been approved for use by all IBM Cloud users.
	PublicPublishApproved *bool

	// The original offering CRN that this publish entry came from.
	PublicOriginalCRN *string

	// The crn of the public catalog entry of this offering.
	PublishPublicCRN *string

	// The portal's approval record ID.
	PortalApprovalRecord *string

	// The portal UI URL.
	PortalUIURL *string

	// The id of the catalog containing this offering.
	CatalogID *string

	// The name of the catalog.
	CatalogName *string

	// Map of metadata values for this offering.
	Metadata map[string]interface{}

	// A disclaimer for this offering.
	Disclaimer *string

	// Determine if this offering should be displayed in the Consumption UI.
	Hidden *bool

	// Provider of this offering.
	Provider *string

	// Repository info for offerings.
	RepoInfo *RepoInfo

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceOfferingOptions : Instantiate ReplaceOfferingOptions
func (*CatalogManagementV1) NewReplaceOfferingOptions(catalogIdentifier string, offeringID string) *ReplaceOfferingOptions {
	return &ReplaceOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *ReplaceOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *ReplaceOfferingOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *ReplaceOfferingOptions) SetOfferingID(offeringID string) *ReplaceOfferingOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetID : Allow user to set ID
func (options *ReplaceOfferingOptions) SetID(id string) *ReplaceOfferingOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetRev : Allow user to set Rev
func (options *ReplaceOfferingOptions) SetRev(rev string) *ReplaceOfferingOptions {
	options.Rev = core.StringPtr(rev)
	return options
}

// SetURL : Allow user to set URL
func (options *ReplaceOfferingOptions) SetURL(url string) *ReplaceOfferingOptions {
	options.URL = core.StringPtr(url)
	return options
}

// SetCRN : Allow user to set CRN
func (options *ReplaceOfferingOptions) SetCRN(crn string) *ReplaceOfferingOptions {
	options.CRN = core.StringPtr(crn)
	return options
}

// SetLabel : Allow user to set Label
func (options *ReplaceOfferingOptions) SetLabel(label string) *ReplaceOfferingOptions {
	options.Label = core.StringPtr(label)
	return options
}

// SetName : Allow user to set Name
func (options *ReplaceOfferingOptions) SetName(name string) *ReplaceOfferingOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetOfferingIconURL : Allow user to set OfferingIconURL
func (options *ReplaceOfferingOptions) SetOfferingIconURL(offeringIconURL string) *ReplaceOfferingOptions {
	options.OfferingIconURL = core.StringPtr(offeringIconURL)
	return options
}

// SetOfferingDocsURL : Allow user to set OfferingDocsURL
func (options *ReplaceOfferingOptions) SetOfferingDocsURL(offeringDocsURL string) *ReplaceOfferingOptions {
	options.OfferingDocsURL = core.StringPtr(offeringDocsURL)
	return options
}

// SetOfferingSupportURL : Allow user to set OfferingSupportURL
func (options *ReplaceOfferingOptions) SetOfferingSupportURL(offeringSupportURL string) *ReplaceOfferingOptions {
	options.OfferingSupportURL = core.StringPtr(offeringSupportURL)
	return options
}

// SetTags : Allow user to set Tags
func (options *ReplaceOfferingOptions) SetTags(tags []string) *ReplaceOfferingOptions {
	options.Tags = tags
	return options
}

// SetKeywords : Allow user to set Keywords
func (options *ReplaceOfferingOptions) SetKeywords(keywords []string) *ReplaceOfferingOptions {
	options.Keywords = keywords
	return options
}

// SetRating : Allow user to set Rating
func (options *ReplaceOfferingOptions) SetRating(rating *Rating) *ReplaceOfferingOptions {
	options.Rating = rating
	return options
}

// SetCreated : Allow user to set Created
func (options *ReplaceOfferingOptions) SetCreated(created *strfmt.DateTime) *ReplaceOfferingOptions {
	options.Created = created
	return options
}

// SetUpdated : Allow user to set Updated
func (options *ReplaceOfferingOptions) SetUpdated(updated *strfmt.DateTime) *ReplaceOfferingOptions {
	options.Updated = updated
	return options
}

// SetShortDescription : Allow user to set ShortDescription
func (options *ReplaceOfferingOptions) SetShortDescription(shortDescription string) *ReplaceOfferingOptions {
	options.ShortDescription = core.StringPtr(shortDescription)
	return options
}

// SetLongDescription : Allow user to set LongDescription
func (options *ReplaceOfferingOptions) SetLongDescription(longDescription string) *ReplaceOfferingOptions {
	options.LongDescription = core.StringPtr(longDescription)
	return options
}

// SetFeatures : Allow user to set Features
func (options *ReplaceOfferingOptions) SetFeatures(features []Feature) *ReplaceOfferingOptions {
	options.Features = features
	return options
}

// SetKinds : Allow user to set Kinds
func (options *ReplaceOfferingOptions) SetKinds(kinds []Kind) *ReplaceOfferingOptions {
	options.Kinds = kinds
	return options
}

// SetPermitRequestIBMPublicPublish : Allow user to set PermitRequestIBMPublicPublish
func (options *ReplaceOfferingOptions) SetPermitRequestIBMPublicPublish(permitRequestIBMPublicPublish bool) *ReplaceOfferingOptions {
	options.PermitRequestIBMPublicPublish = core.BoolPtr(permitRequestIBMPublicPublish)
	return options
}

// SetIBMPublishApproved : Allow user to set IBMPublishApproved
func (options *ReplaceOfferingOptions) SetIBMPublishApproved(ibmPublishApproved bool) *ReplaceOfferingOptions {
	options.IBMPublishApproved = core.BoolPtr(ibmPublishApproved)
	return options
}

// SetPublicPublishApproved : Allow user to set PublicPublishApproved
func (options *ReplaceOfferingOptions) SetPublicPublishApproved(publicPublishApproved bool) *ReplaceOfferingOptions {
	options.PublicPublishApproved = core.BoolPtr(publicPublishApproved)
	return options
}

// SetPublicOriginalCRN : Allow user to set PublicOriginalCRN
func (options *ReplaceOfferingOptions) SetPublicOriginalCRN(publicOriginalCRN string) *ReplaceOfferingOptions {
	options.PublicOriginalCRN = core.StringPtr(publicOriginalCRN)
	return options
}

// SetPublishPublicCRN : Allow user to set PublishPublicCRN
func (options *ReplaceOfferingOptions) SetPublishPublicCRN(publishPublicCRN string) *ReplaceOfferingOptions {
	options.PublishPublicCRN = core.StringPtr(publishPublicCRN)
	return options
}

// SetPortalApprovalRecord : Allow user to set PortalApprovalRecord
func (options *ReplaceOfferingOptions) SetPortalApprovalRecord(portalApprovalRecord string) *ReplaceOfferingOptions {
	options.PortalApprovalRecord = core.StringPtr(portalApprovalRecord)
	return options
}

// SetPortalUIURL : Allow user to set PortalUIURL
func (options *ReplaceOfferingOptions) SetPortalUIURL(portalUIURL string) *ReplaceOfferingOptions {
	options.PortalUIURL = core.StringPtr(portalUIURL)
	return options
}

// SetCatalogID : Allow user to set CatalogID
func (options *ReplaceOfferingOptions) SetCatalogID(catalogID string) *ReplaceOfferingOptions {
	options.CatalogID = core.StringPtr(catalogID)
	return options
}

// SetCatalogName : Allow user to set CatalogName
func (options *ReplaceOfferingOptions) SetCatalogName(catalogName string) *ReplaceOfferingOptions {
	options.CatalogName = core.StringPtr(catalogName)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *ReplaceOfferingOptions) SetMetadata(metadata map[string]interface{}) *ReplaceOfferingOptions {
	options.Metadata = metadata
	return options
}

// SetDisclaimer : Allow user to set Disclaimer
func (options *ReplaceOfferingOptions) SetDisclaimer(disclaimer string) *ReplaceOfferingOptions {
	options.Disclaimer = core.StringPtr(disclaimer)
	return options
}

// SetHidden : Allow user to set Hidden
func (options *ReplaceOfferingOptions) SetHidden(hidden bool) *ReplaceOfferingOptions {
	options.Hidden = core.BoolPtr(hidden)
	return options
}

// SetProvider : Allow user to set Provider
func (options *ReplaceOfferingOptions) SetProvider(provider string) *ReplaceOfferingOptions {
	options.Provider = core.StringPtr(provider)
	return options
}

// SetRepoInfo : Allow user to set RepoInfo
func (options *ReplaceOfferingOptions) SetRepoInfo(repoInfo *RepoInfo) *ReplaceOfferingOptions {
	options.RepoInfo = repoInfo
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceOfferingOptions) SetHeaders(param map[string]string) *ReplaceOfferingOptions {
	options.Headers = param
	return options
}

// ReplaceOperatorsOptions : The ReplaceOperators options.
type ReplaceOperatorsOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Cluster ID.
	ClusterID *string

	// Cluster region.
	Region *string

	// Kube namespaces to deploy Operator(s) to.
	Namespaces []string

	// Denotes whether to install Operator(s) globally.
	AllNamespaces *bool

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceOperatorsOptions : Instantiate ReplaceOperatorsOptions
func (*CatalogManagementV1) NewReplaceOperatorsOptions(xAuthRefreshToken string) *ReplaceOperatorsOptions {
	return &ReplaceOperatorsOptions{
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *ReplaceOperatorsOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *ReplaceOperatorsOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *ReplaceOperatorsOptions) SetClusterID(clusterID string) *ReplaceOperatorsOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *ReplaceOperatorsOptions) SetRegion(region string) *ReplaceOperatorsOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetNamespaces : Allow user to set Namespaces
func (options *ReplaceOperatorsOptions) SetNamespaces(namespaces []string) *ReplaceOperatorsOptions {
	options.Namespaces = namespaces
	return options
}

// SetAllNamespaces : Allow user to set AllNamespaces
func (options *ReplaceOperatorsOptions) SetAllNamespaces(allNamespaces bool) *ReplaceOperatorsOptions {
	options.AllNamespaces = core.BoolPtr(allNamespaces)
	return options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (options *ReplaceOperatorsOptions) SetVersionLocatorID(versionLocatorID string) *ReplaceOperatorsOptions {
	options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceOperatorsOptions) SetHeaders(param map[string]string) *ReplaceOperatorsOptions {
	options.Headers = param
	return options
}

// RepoInfo : Repository info for offerings.
type RepoInfo struct {
	// Token for private repos.
	Token *string `json:"token,omitempty"`

	// Public or enterprise GitHub.
	Type *string `json:"type,omitempty"`
}

// UnmarshalRepoInfo unmarshals an instance of RepoInfo from the specified map of raw messages.
func UnmarshalRepoInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RepoInfo)
	err = core.UnmarshalPrimitive(m, "token", &obj.Token)
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

// Resource : Resource requirements.
type Resource struct {
	// Type of requirement.
	Type *string `json:"type,omitempty"`

	// mem, disk, cores, and nodes can be parsed as an int.  targetVersion will be a semver range value.
	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the Resource.Type property.
// Type of requirement.
const (
	ResourceTypeCoresConst = "cores"
	ResourceTypeDiskConst = "disk"
	ResourceTypeMemConst = "mem"
	ResourceTypeNodesConst = "nodes"
	ResourceTypeTargetversionConst = "targetVersion"
)

// UnmarshalResource unmarshals an instance of Resource from the specified map of raw messages.
func UnmarshalResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Resource)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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

// Script : Script information.
type Script struct {
	// Instruction on step and by whom (role) that are needed to take place to prepare the target for installing this
	// version.
	Instructions *string `json:"instructions,omitempty"`

	// Optional script that needs to be run post any pre-condition script.
	Script *string `json:"script,omitempty"`

	// Optional iam permissions that are required on the target cluster to run this script.
	ScriptPermission *string `json:"script_permission,omitempty"`

	// Optional script that if run will remove the installed version.
	DeleteScript *string `json:"delete_script,omitempty"`

	// Optional value indicating if this script is scoped to a namespace or the entire cluster.
	Scope *string `json:"scope,omitempty"`
}

// UnmarshalScript unmarshals an instance of Script from the specified map of raw messages.
func UnmarshalScript(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Script)
	err = core.UnmarshalPrimitive(m, "instructions", &obj.Instructions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "script", &obj.Script)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "script_permission", &obj.ScriptPermission)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "delete_script", &obj.DeleteScript)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope", &obj.Scope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SearchObjectsOptions : The SearchObjects options.
type SearchObjectsOptions struct {
	// Lucene query string.
	Query *string `validate:"required"`

	// The maximum number of results to return.
	Limit *int64

	// The number of results to skip before returning values.
	Offset *int64

	// When true, hide private objects that correspond to public or IBM published objects.
	Collapse *bool

	// Display a digests of search results, has default value of true.
	Digest *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSearchObjectsOptions : Instantiate SearchObjectsOptions
func (*CatalogManagementV1) NewSearchObjectsOptions(query string) *SearchObjectsOptions {
	return &SearchObjectsOptions{
		Query: core.StringPtr(query),
	}
}

// SetQuery : Allow user to set Query
func (options *SearchObjectsOptions) SetQuery(query string) *SearchObjectsOptions {
	options.Query = core.StringPtr(query)
	return options
}

// SetLimit : Allow user to set Limit
func (options *SearchObjectsOptions) SetLimit(limit int64) *SearchObjectsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *SearchObjectsOptions) SetOffset(offset int64) *SearchObjectsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetCollapse : Allow user to set Collapse
func (options *SearchObjectsOptions) SetCollapse(collapse bool) *SearchObjectsOptions {
	options.Collapse = core.BoolPtr(collapse)
	return options
}

// SetDigest : Allow user to set Digest
func (options *SearchObjectsOptions) SetDigest(digest bool) *SearchObjectsOptions {
	options.Digest = core.BoolPtr(digest)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SearchObjectsOptions) SetHeaders(param map[string]string) *SearchObjectsOptions {
	options.Headers = param
	return options
}

// SharedPublishObjectOptions : The SharedPublishObject options.
type SharedPublishObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSharedPublishObjectOptions : Instantiate SharedPublishObjectOptions
func (*CatalogManagementV1) NewSharedPublishObjectOptions(catalogIdentifier string, objectIdentifier string) *SharedPublishObjectOptions {
	return &SharedPublishObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *SharedPublishObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *SharedPublishObjectOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (options *SharedPublishObjectOptions) SetObjectIdentifier(objectIdentifier string) *SharedPublishObjectOptions {
	options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SharedPublishObjectOptions) SetHeaders(param map[string]string) *SharedPublishObjectOptions {
	options.Headers = param
	return options
}

// State : Offering state.
type State struct {
	// one of: new, validated, account-published, ibm-published, public-published.
	Current *string `json:"current,omitempty"`

	// Date and time of current request.
	CurrentEntered *strfmt.DateTime `json:"current_entered,omitempty"`

	// one of: new, validated, account-published, ibm-published, public-published.
	Pending *string `json:"pending,omitempty"`

	// Date and time of pending request.
	PendingRequested *strfmt.DateTime `json:"pending_requested,omitempty"`

	// one of: new, validated, account-published, ibm-published, public-published.
	Previous *string `json:"previous,omitempty"`
}

// UnmarshalState unmarshals an instance of State from the specified map of raw messages.
func UnmarshalState(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(State)
	err = core.UnmarshalPrimitive(m, "current", &obj.Current)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "current_entered", &obj.CurrentEntered)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pending", &obj.Pending)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pending_requested", &obj.PendingRequested)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "previous", &obj.Previous)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SyndicationAuthorization : Feature information.
type SyndicationAuthorization struct {
	// Array of syndicated namespaces.
	Token *string `json:"token,omitempty"`

	// Date and time last updated.
	LastRun *strfmt.DateTime `json:"last_run,omitempty"`
}

// UnmarshalSyndicationAuthorization unmarshals an instance of SyndicationAuthorization from the specified map of raw messages.
func UnmarshalSyndicationAuthorization(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SyndicationAuthorization)
	err = core.UnmarshalPrimitive(m, "token", &obj.Token)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_run", &obj.LastRun)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SyndicationCluster : Feature information.
type SyndicationCluster struct {
	// Cluster region.
	Region *string `json:"region,omitempty"`

	// Cluster ID.
	ID *string `json:"id,omitempty"`

	// Cluster name.
	Name *string `json:"name,omitempty"`

	// Resource group ID.
	ResourceGroupName *string `json:"resource_group_name,omitempty"`

	// Syndication type.
	Type *string `json:"type,omitempty"`

	// Syndicated namespaces.
	Namespaces []string `json:"namespaces,omitempty"`

	// Syndicated to all namespaces on cluster.
	AllNamespaces *bool `json:"all_namespaces,omitempty"`
}

// UnmarshalSyndicationCluster unmarshals an instance of SyndicationCluster from the specified map of raw messages.
func UnmarshalSyndicationCluster(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SyndicationCluster)
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
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
	err = core.UnmarshalPrimitive(m, "resource_group_name", &obj.ResourceGroupName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "namespaces", &obj.Namespaces)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "all_namespaces", &obj.AllNamespaces)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SyndicationHistory : Feature information.
type SyndicationHistory struct {
	// Array of syndicated namespaces.
	Namespaces []string `json:"namespaces,omitempty"`

	// Array of syndicated namespaces.
	Clusters []SyndicationCluster `json:"clusters,omitempty"`

	// Date and time last syndicated.
	LastRun *strfmt.DateTime `json:"last_run,omitempty"`
}

// UnmarshalSyndicationHistory unmarshals an instance of SyndicationHistory from the specified map of raw messages.
func UnmarshalSyndicationHistory(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SyndicationHistory)
	err = core.UnmarshalPrimitive(m, "namespaces", &obj.Namespaces)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "clusters", &obj.Clusters, UnmarshalSyndicationCluster)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_run", &obj.LastRun)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SyndicationResource : Feature information.
type SyndicationResource struct {
	// Remove related components.
	RemoveRelatedComponents *bool `json:"remove_related_components,omitempty"`

	// Syndication clusters.
	Clusters []SyndicationCluster `json:"clusters,omitempty"`

	// Feature information.
	History *SyndicationHistory `json:"history,omitempty"`

	// Feature information.
	Authorization *SyndicationAuthorization `json:"authorization,omitempty"`
}

// UnmarshalSyndicationResource unmarshals an instance of SyndicationResource from the specified map of raw messages.
func UnmarshalSyndicationResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SyndicationResource)
	err = core.UnmarshalPrimitive(m, "remove_related_components", &obj.RemoveRelatedComponents)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "clusters", &obj.Clusters, UnmarshalSyndicationCluster)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalSyndicationHistory)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorization", &obj.Authorization, UnmarshalSyndicationAuthorization)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateCatalogAccountOptions : The UpdateCatalogAccount options.
type UpdateCatalogAccountOptions struct {
	// Account identification.
	ID *string

	// Hide the public catalog in this account.
	HideIBMCloudCatalog *bool

	// Filters for account and catalog filters.
	AccountFilters *Filters

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateCatalogAccountOptions : Instantiate UpdateCatalogAccountOptions
func (*CatalogManagementV1) NewUpdateCatalogAccountOptions() *UpdateCatalogAccountOptions {
	return &UpdateCatalogAccountOptions{}
}

// SetID : Allow user to set ID
func (options *UpdateCatalogAccountOptions) SetID(id string) *UpdateCatalogAccountOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHideIBMCloudCatalog : Allow user to set HideIBMCloudCatalog
func (options *UpdateCatalogAccountOptions) SetHideIBMCloudCatalog(hideIBMCloudCatalog bool) *UpdateCatalogAccountOptions {
	options.HideIBMCloudCatalog = core.BoolPtr(hideIBMCloudCatalog)
	return options
}

// SetAccountFilters : Allow user to set AccountFilters
func (options *UpdateCatalogAccountOptions) SetAccountFilters(accountFilters *Filters) *UpdateCatalogAccountOptions {
	options.AccountFilters = accountFilters
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCatalogAccountOptions) SetHeaders(param map[string]string) *UpdateCatalogAccountOptions {
	options.Headers = param
	return options
}

// UpdateOfferingIBMOptions : The UpdateOfferingIBM options.
type UpdateOfferingIBMOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `validate:"required,ne="`

	// Offering identification.
	OfferingID *string `validate:"required,ne="`

	// Type of approval, ibm or public.
	ApprovalType *string `validate:"required,ne="`

	// Approve (true) or disapprove (false).
	Approved *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateOfferingIBMOptions.ApprovalType property.
// Type of approval, ibm or public.
const (
	UpdateOfferingIBMOptionsApprovalTypeAllowRequestConst = "allow_request"
	UpdateOfferingIBMOptionsApprovalTypeIBMConst = "ibm"
	UpdateOfferingIBMOptionsApprovalTypePublicConst = "public"
)

// Constants associated with the UpdateOfferingIBMOptions.Approved property.
// Approve (true) or disapprove (false).
const (
	UpdateOfferingIBMOptionsApprovedFalseConst = "false"
	UpdateOfferingIBMOptionsApprovedTrueConst = "true"
)

// NewUpdateOfferingIBMOptions : Instantiate UpdateOfferingIBMOptions
func (*CatalogManagementV1) NewUpdateOfferingIBMOptions(catalogIdentifier string, offeringID string, approvalType string, approved string) *UpdateOfferingIBMOptions {
	return &UpdateOfferingIBMOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		ApprovalType: core.StringPtr(approvalType),
		Approved: core.StringPtr(approved),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (options *UpdateOfferingIBMOptions) SetCatalogIdentifier(catalogIdentifier string) *UpdateOfferingIBMOptions {
	options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return options
}

// SetOfferingID : Allow user to set OfferingID
func (options *UpdateOfferingIBMOptions) SetOfferingID(offeringID string) *UpdateOfferingIBMOptions {
	options.OfferingID = core.StringPtr(offeringID)
	return options
}

// SetApprovalType : Allow user to set ApprovalType
func (options *UpdateOfferingIBMOptions) SetApprovalType(approvalType string) *UpdateOfferingIBMOptions {
	options.ApprovalType = core.StringPtr(approvalType)
	return options
}

// SetApproved : Allow user to set Approved
func (options *UpdateOfferingIBMOptions) SetApproved(approved string) *UpdateOfferingIBMOptions {
	options.Approved = core.StringPtr(approved)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateOfferingIBMOptions) SetHeaders(param map[string]string) *UpdateOfferingIBMOptions {
	options.Headers = param
	return options
}

// ValidateInstallOptions : The ValidateInstall options.
type ValidateInstallOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `validate:"required"`

	// Cluster ID.
	ClusterID *string

	// Cluster region.
	Region *string

	// Kube namespace.
	Namespace *string

	// Object containing Helm chart override values.  To use a secret for items of type password, specify a JSON encoded
	// value of $ref:#/components/schemas/SecretInstance, prefixed with `cmsm_v1:`.
	OverrideValues map[string]interface{}

	// Entitlement API Key for this offering.
	EntitlementApikey *string

	// Schematics workspace configuration.
	Schematics *DeployRequestBodySchematics

	// Script.
	Script *string

	// Script ID.
	ScriptID *string

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string

	// VCenter ID.
	VcenterID *string

	// VCenter User.
	VcenterUser *string

	// VCenter Password.
	VcenterPassword *string

	// VCenter Location.
	VcenterLocation *string

	// VCenter Datastore.
	VcenterDatastore *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewValidateInstallOptions : Instantiate ValidateInstallOptions
func (*CatalogManagementV1) NewValidateInstallOptions(versionLocID string, xAuthRefreshToken string) *ValidateInstallOptions {
	return &ValidateInstallOptions{
		VersionLocID: core.StringPtr(versionLocID),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (options *ValidateInstallOptions) SetVersionLocID(versionLocID string) *ValidateInstallOptions {
	options.VersionLocID = core.StringPtr(versionLocID)
	return options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (options *ValidateInstallOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *ValidateInstallOptions {
	options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return options
}

// SetClusterID : Allow user to set ClusterID
func (options *ValidateInstallOptions) SetClusterID(clusterID string) *ValidateInstallOptions {
	options.ClusterID = core.StringPtr(clusterID)
	return options
}

// SetRegion : Allow user to set Region
func (options *ValidateInstallOptions) SetRegion(region string) *ValidateInstallOptions {
	options.Region = core.StringPtr(region)
	return options
}

// SetNamespace : Allow user to set Namespace
func (options *ValidateInstallOptions) SetNamespace(namespace string) *ValidateInstallOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetOverrideValues : Allow user to set OverrideValues
func (options *ValidateInstallOptions) SetOverrideValues(overrideValues map[string]interface{}) *ValidateInstallOptions {
	options.OverrideValues = overrideValues
	return options
}

// SetEntitlementApikey : Allow user to set EntitlementApikey
func (options *ValidateInstallOptions) SetEntitlementApikey(entitlementApikey string) *ValidateInstallOptions {
	options.EntitlementApikey = core.StringPtr(entitlementApikey)
	return options
}

// SetSchematics : Allow user to set Schematics
func (options *ValidateInstallOptions) SetSchematics(schematics *DeployRequestBodySchematics) *ValidateInstallOptions {
	options.Schematics = schematics
	return options
}

// SetScript : Allow user to set Script
func (options *ValidateInstallOptions) SetScript(script string) *ValidateInstallOptions {
	options.Script = core.StringPtr(script)
	return options
}

// SetScriptID : Allow user to set ScriptID
func (options *ValidateInstallOptions) SetScriptID(scriptID string) *ValidateInstallOptions {
	options.ScriptID = core.StringPtr(scriptID)
	return options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (options *ValidateInstallOptions) SetVersionLocatorID(versionLocatorID string) *ValidateInstallOptions {
	options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return options
}

// SetVcenterID : Allow user to set VcenterID
func (options *ValidateInstallOptions) SetVcenterID(vcenterID string) *ValidateInstallOptions {
	options.VcenterID = core.StringPtr(vcenterID)
	return options
}

// SetVcenterUser : Allow user to set VcenterUser
func (options *ValidateInstallOptions) SetVcenterUser(vcenterUser string) *ValidateInstallOptions {
	options.VcenterUser = core.StringPtr(vcenterUser)
	return options
}

// SetVcenterPassword : Allow user to set VcenterPassword
func (options *ValidateInstallOptions) SetVcenterPassword(vcenterPassword string) *ValidateInstallOptions {
	options.VcenterPassword = core.StringPtr(vcenterPassword)
	return options
}

// SetVcenterLocation : Allow user to set VcenterLocation
func (options *ValidateInstallOptions) SetVcenterLocation(vcenterLocation string) *ValidateInstallOptions {
	options.VcenterLocation = core.StringPtr(vcenterLocation)
	return options
}

// SetVcenterDatastore : Allow user to set VcenterDatastore
func (options *ValidateInstallOptions) SetVcenterDatastore(vcenterDatastore string) *ValidateInstallOptions {
	options.VcenterDatastore = core.StringPtr(vcenterDatastore)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ValidateInstallOptions) SetHeaders(param map[string]string) *ValidateInstallOptions {
	options.Headers = param
	return options
}

// Validation : Validation response.
type Validation struct {
	// Date and time of last successful validation.
	Validated *strfmt.DateTime `json:"validated,omitempty"`

	// Date and time of last validation was requested.
	Requested *strfmt.DateTime `json:"requested,omitempty"`

	// Current validation state - <empty>, in_progress, valid, invalid, expired.
	State *string `json:"state,omitempty"`

	// Last operation (e.g. submit_deployment, generate_installer, install_offering.
	LastOperation *string `json:"last_operation,omitempty"`

	// Validation target information (e.g. cluster_id, region, namespace, etc).  Values will vary by Content type.
	Target map[string]interface{} `json:"target,omitempty"`
}

// UnmarshalValidation unmarshals an instance of Validation from the specified map of raw messages.
func UnmarshalValidation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Validation)
	err = core.UnmarshalPrimitive(m, "validated", &obj.Validated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "requested", &obj.Requested)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_operation", &obj.LastOperation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Version : Offering version information.
type Version struct {
	// Unique ID.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// Version's CRN.
	CRN *string `json:"crn,omitempty"`

	// Version of content type.
	Version *string `json:"version,omitempty"`

	// hash of the content.
	Sha *string `json:"sha,omitempty"`

	// The date and time this version was created.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The date and time this version was last updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Offering ID.
	OfferingID *string `json:"offering_id,omitempty"`

	// Catalog ID.
	CatalogID *string `json:"catalog_id,omitempty"`

	// Kind ID.
	KindID *string `json:"kind_id,omitempty"`

	// List of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// Content's repo URL.
	RepoURL *string `json:"repo_url,omitempty"`

	// Content's source URL (e.g git repo).
	SourceURL *string `json:"source_url,omitempty"`

	// File used to on-board this version.
	TgzURL *string `json:"tgz_url,omitempty"`

	// List of user solicited overrides.
	Configuration []Configuration `json:"configuration,omitempty"`

	// Open ended metadata information.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Validation response.
	Validation *Validation `json:"validation,omitempty"`

	// Resource requirments for installation.
	RequiredResources []Resource `json:"required_resources,omitempty"`

	// Denotes if single instance can be deployed to a given cluster.
	SingleInstance *bool `json:"single_instance,omitempty"`

	// Script information.
	Install *Script `json:"install,omitempty"`

	// Optional pre-install instructions.
	PreInstall []Script `json:"pre_install,omitempty"`

	// Entitlement license info.
	Entitlement *VersionEntitlement `json:"entitlement,omitempty"`

	// List of licenses the product was built with.
	Licenses []License `json:"licenses,omitempty"`

	// If set, denotes a url to a YAML file with list of container images used by this version.
	ImageManifestURL *string `json:"image_manifest_url,omitempty"`

	// read only field, indicating if this version is deprecated.
	Deprecated *bool `json:"deprecated,omitempty"`

	// Version of the package used to create this version.
	PackageVersion *string `json:"package_version,omitempty"`

	// Offering state.
	State *State `json:"state,omitempty"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocator *string `json:"version_locator,omitempty"`

	// Console URL.
	ConsoleURL *string `json:"console_url,omitempty"`

	// Long description for version.
	LongDescription *string `json:"long_description,omitempty"`

	// Whitelisted accounts for version.
	WhitelistedAccounts []string `json:"whitelisted_accounts,omitempty"`
}

// UnmarshalVersion unmarshals an instance of Version from the specified map of raw messages.
func UnmarshalVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Version)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sha", &obj.Sha)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_id", &obj.OfferingID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind_id", &obj.KindID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "repo_url", &obj.RepoURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_url", &obj.SourceURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tgz_url", &obj.TgzURL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "configuration", &obj.Configuration, UnmarshalConfiguration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validation", &obj.Validation, UnmarshalValidation)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "required_resources", &obj.RequiredResources, UnmarshalResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "single_instance", &obj.SingleInstance)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "install", &obj.Install, UnmarshalScript)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "pre_install", &obj.PreInstall, UnmarshalScript)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "entitlement", &obj.Entitlement, UnmarshalVersionEntitlement)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "licenses", &obj.Licenses, UnmarshalLicense)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "image_manifest_url", &obj.ImageManifestURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "deprecated", &obj.Deprecated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "package_version", &obj.PackageVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "state", &obj.State, UnmarshalState)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_locator", &obj.VersionLocator)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "console_url", &obj.ConsoleURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "whitelisted_accounts", &obj.WhitelistedAccounts)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VersionEntitlement : Entitlement license info.
type VersionEntitlement struct {
	// Provider name.
	ProviderName *string `json:"provider_name,omitempty"`

	// Provider ID.
	ProviderID *string `json:"provider_id,omitempty"`

	// Product ID.
	ProductID *string `json:"product_id,omitempty"`

	// list of license entitlement part numbers, eg. D1YGZLL,D1ZXILL.
	PartNumbers []string `json:"part_numbers,omitempty"`

	// Image repository name.
	ImageRepoName *string `json:"image_repo_name,omitempty"`
}

// UnmarshalVersionEntitlement unmarshals an instance of VersionEntitlement from the specified map of raw messages.
func UnmarshalVersionEntitlement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VersionEntitlement)
	err = core.UnmarshalPrimitive(m, "provider_name", &obj.ProviderName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "provider_id", &obj.ProviderID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "product_id", &obj.ProductID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "part_numbers", &obj.PartNumbers)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "image_repo_name", &obj.ImageRepoName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VersionUpdateDescriptor : Indicates if the current version can be upgraded to the version identified by the descriptor.
type VersionUpdateDescriptor struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocator *string `json:"version_locator,omitempty"`

	// the version number of this version.
	Version *string `json:"version,omitempty"`

	// Offering state.
	State *State `json:"state,omitempty"`

	// Resource requirments for installation.
	RequiredResources []Resource `json:"required_resources,omitempty"`

	// Version of package.
	PackageVersion *string `json:"package_version,omitempty"`

	// true if the current version can be upgraded to this version, false otherwise.
	CanUpdate *bool `json:"can_update,omitempty"`

	// If can_update is false, this map will contain messages for each failed check, otherwise it will be omitted.
	// Possible keys include nodes, cores, mem, disk, targetVersion, and install-permission-check.
	Messages map[string]string `json:"messages,omitempty"`
}

// UnmarshalVersionUpdateDescriptor unmarshals an instance of VersionUpdateDescriptor from the specified map of raw messages.
func UnmarshalVersionUpdateDescriptor(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VersionUpdateDescriptor)
	err = core.UnmarshalPrimitive(m, "version_locator", &obj.VersionLocator)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "state", &obj.State, UnmarshalState)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "required_resources", &obj.RequiredResources, UnmarshalResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "package_version", &obj.PackageVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "can_update", &obj.CanUpdate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
