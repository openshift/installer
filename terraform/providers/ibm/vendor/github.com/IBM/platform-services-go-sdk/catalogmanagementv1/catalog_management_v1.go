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
 * IBM OpenAPI SDK Code Generator Version: 3.68.2-ac7def68-20230310-195410
 */

// Package catalogmanagementv1 : Operations and models for the CatalogManagementV1 service
package catalogmanagementv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// CatalogManagementV1 : This is the API to use for managing private catalogs for IBM Cloud. Private catalogs provide a
// way to centrally manage access to products in the IBM Cloud catalog and your own catalogs.
//
// API Version: 1.0
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
func (catalogManagement *CatalogManagementV1) UpdateCatalogAccount(updateCatalogAccountOptions *UpdateCatalogAccountOptions) (result *Account, response *core.DetailedResponse, err error) {
	return catalogManagement.UpdateCatalogAccountWithContext(context.Background(), updateCatalogAccountOptions)
}

// UpdateCatalogAccountWithContext is an alternate form of the UpdateCatalogAccount method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) UpdateCatalogAccountWithContext(ctx context.Context, updateCatalogAccountOptions *UpdateCatalogAccountOptions) (result *Account, response *core.DetailedResponse, err error) {
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
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateCatalogAccountOptions.ID != nil {
		body["id"] = updateCatalogAccountOptions.ID
	}
	if updateCatalogAccountOptions.Rev != nil {
		body["_rev"] = updateCatalogAccountOptions.Rev
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

// ListCatalogAccountAudits : Get catalog account audit logs
// Get the audit logs associated with a catalog account.
func (catalogManagement *CatalogManagementV1) ListCatalogAccountAudits(listCatalogAccountAuditsOptions *ListCatalogAccountAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	return catalogManagement.ListCatalogAccountAuditsWithContext(context.Background(), listCatalogAccountAuditsOptions)
}

// ListCatalogAccountAuditsWithContext is an alternate form of the ListCatalogAccountAudits method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListCatalogAccountAuditsWithContext(ctx context.Context, listCatalogAccountAuditsOptions *ListCatalogAccountAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCatalogAccountAuditsOptions, "listCatalogAccountAuditsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogaccount/audits`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listCatalogAccountAuditsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListCatalogAccountAudits")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listCatalogAccountAuditsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listCatalogAccountAuditsOptions.Start))
	}
	if listCatalogAccountAuditsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listCatalogAccountAuditsOptions.Limit))
	}
	if listCatalogAccountAuditsOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*listCatalogAccountAuditsOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLogs)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCatalogAccountAudit : Get a catalog account audit log entry
// Get the full audit log entry associated with a catalog account.
func (catalogManagement *CatalogManagementV1) GetCatalogAccountAudit(getCatalogAccountAuditOptions *GetCatalogAccountAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	return catalogManagement.GetCatalogAccountAuditWithContext(context.Background(), getCatalogAccountAuditOptions)
}

// GetCatalogAccountAuditWithContext is an alternate form of the GetCatalogAccountAudit method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetCatalogAccountAuditWithContext(ctx context.Context, getCatalogAccountAuditOptions *GetCatalogAccountAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCatalogAccountAuditOptions, "getCatalogAccountAuditOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCatalogAccountAuditOptions, "getCatalogAccountAuditOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"auditlog_identifier": *getCatalogAccountAuditOptions.AuditlogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogaccount/audits/{auditlog_identifier}`, pathParamsMap)
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

	if getCatalogAccountAuditOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*getCatalogAccountAuditOptions.Lookupnames))
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
	if createCatalogOptions.Label != nil {
		body["label"] = createCatalogOptions.Label
	}
	if createCatalogOptions.LabelI18n != nil {
		body["label_i18n"] = createCatalogOptions.LabelI18n
	}
	if createCatalogOptions.ShortDescription != nil {
		body["short_description"] = createCatalogOptions.ShortDescription
	}
	if createCatalogOptions.ShortDescriptionI18n != nil {
		body["short_description_i18n"] = createCatalogOptions.ShortDescriptionI18n
	}
	if createCatalogOptions.CatalogIconURL != nil {
		body["catalog_icon_url"] = createCatalogOptions.CatalogIconURL
	}
	if createCatalogOptions.CatalogBannerURL != nil {
		body["catalog_banner_url"] = createCatalogOptions.CatalogBannerURL
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
	if createCatalogOptions.Metadata != nil {
		body["metadata"] = createCatalogOptions.Metadata
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
	if replaceCatalogOptions.LabelI18n != nil {
		body["label_i18n"] = replaceCatalogOptions.LabelI18n
	}
	if replaceCatalogOptions.ShortDescription != nil {
		body["short_description"] = replaceCatalogOptions.ShortDescription
	}
	if replaceCatalogOptions.ShortDescriptionI18n != nil {
		body["short_description_i18n"] = replaceCatalogOptions.ShortDescriptionI18n
	}
	if replaceCatalogOptions.CatalogIconURL != nil {
		body["catalog_icon_url"] = replaceCatalogOptions.CatalogIconURL
	}
	if replaceCatalogOptions.CatalogBannerURL != nil {
		body["catalog_banner_url"] = replaceCatalogOptions.CatalogBannerURL
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
	if replaceCatalogOptions.Metadata != nil {
		body["metadata"] = replaceCatalogOptions.Metadata
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

// ListCatalogAudits : Get catalog audit logs
// Get the audit logs associated with a catalog.
func (catalogManagement *CatalogManagementV1) ListCatalogAudits(listCatalogAuditsOptions *ListCatalogAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	return catalogManagement.ListCatalogAuditsWithContext(context.Background(), listCatalogAuditsOptions)
}

// ListCatalogAuditsWithContext is an alternate form of the ListCatalogAudits method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListCatalogAuditsWithContext(ctx context.Context, listCatalogAuditsOptions *ListCatalogAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listCatalogAuditsOptions, "listCatalogAuditsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listCatalogAuditsOptions, "listCatalogAuditsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *listCatalogAuditsOptions.CatalogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/audits`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listCatalogAuditsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListCatalogAudits")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listCatalogAuditsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listCatalogAuditsOptions.Start))
	}
	if listCatalogAuditsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listCatalogAuditsOptions.Limit))
	}
	if listCatalogAuditsOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*listCatalogAuditsOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLogs)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCatalogAudit : Get a catalog audit log entry
// Get the full audit log entry associated with a catalog.
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
		"auditlog_identifier": *getCatalogAuditOptions.AuditlogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/audits/{auditlog_identifier}`, pathParamsMap)
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

	if getCatalogAuditOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*getCatalogAuditOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListEnterpriseAudits : Get enterprise audit logs
// Get the audit logs associated with an enterprise.
func (catalogManagement *CatalogManagementV1) ListEnterpriseAudits(listEnterpriseAuditsOptions *ListEnterpriseAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	return catalogManagement.ListEnterpriseAuditsWithContext(context.Background(), listEnterpriseAuditsOptions)
}

// ListEnterpriseAuditsWithContext is an alternate form of the ListEnterpriseAudits method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListEnterpriseAuditsWithContext(ctx context.Context, listEnterpriseAuditsOptions *ListEnterpriseAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listEnterpriseAuditsOptions, "listEnterpriseAuditsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listEnterpriseAuditsOptions, "listEnterpriseAuditsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"enterprise_identifier": *listEnterpriseAuditsOptions.EnterpriseIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/enterprises/{enterprise_identifier}/audits`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listEnterpriseAuditsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListEnterpriseAudits")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listEnterpriseAuditsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listEnterpriseAuditsOptions.Start))
	}
	if listEnterpriseAuditsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listEnterpriseAuditsOptions.Limit))
	}
	if listEnterpriseAuditsOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*listEnterpriseAuditsOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLogs)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetEnterpriseAudit : Get an enterprise audit log entry
// Get the full audit log entry associated with an enterprise.
func (catalogManagement *CatalogManagementV1) GetEnterpriseAudit(getEnterpriseAuditOptions *GetEnterpriseAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	return catalogManagement.GetEnterpriseAuditWithContext(context.Background(), getEnterpriseAuditOptions)
}

// GetEnterpriseAuditWithContext is an alternate form of the GetEnterpriseAudit method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetEnterpriseAuditWithContext(ctx context.Context, getEnterpriseAuditOptions *GetEnterpriseAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEnterpriseAuditOptions, "getEnterpriseAuditOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEnterpriseAuditOptions, "getEnterpriseAuditOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"enterprise_identifier": *getEnterpriseAuditOptions.EnterpriseIdentifier,
		"auditlog_identifier": *getEnterpriseAuditOptions.AuditlogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/enterprises/{enterprise_identifier}/audits/{auditlog_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEnterpriseAuditOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetEnterpriseAudit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getEnterpriseAuditOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*getEnterpriseAuditOptions.Lookupnames))
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
	if listOfferingsOptions.IncludeHidden != nil {
		builder.AddQuery("includeHidden", fmt.Sprint(*listOfferingsOptions.IncludeHidden))
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
	if createOfferingOptions.URL != nil {
		body["url"] = createOfferingOptions.URL
	}
	if createOfferingOptions.CRN != nil {
		body["crn"] = createOfferingOptions.CRN
	}
	if createOfferingOptions.Label != nil {
		body["label"] = createOfferingOptions.Label
	}
	if createOfferingOptions.LabelI18n != nil {
		body["label_i18n"] = createOfferingOptions.LabelI18n
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
	if createOfferingOptions.ShortDescriptionI18n != nil {
		body["short_description_i18n"] = createOfferingOptions.ShortDescriptionI18n
	}
	if createOfferingOptions.LongDescription != nil {
		body["long_description"] = createOfferingOptions.LongDescription
	}
	if createOfferingOptions.LongDescriptionI18n != nil {
		body["long_description_i18n"] = createOfferingOptions.LongDescriptionI18n
	}
	if createOfferingOptions.Features != nil {
		body["features"] = createOfferingOptions.Features
	}
	if createOfferingOptions.Kinds != nil {
		body["kinds"] = createOfferingOptions.Kinds
	}
	if createOfferingOptions.PcManaged != nil {
		body["pc_managed"] = createOfferingOptions.PcManaged
	}
	if createOfferingOptions.PublishApproved != nil {
		body["publish_approved"] = createOfferingOptions.PublishApproved
	}
	if createOfferingOptions.ShareWithAll != nil {
		body["share_with_all"] = createOfferingOptions.ShareWithAll
	}
	if createOfferingOptions.ShareWithIBM != nil {
		body["share_with_ibm"] = createOfferingOptions.ShareWithIBM
	}
	if createOfferingOptions.ShareEnabled != nil {
		body["share_enabled"] = createOfferingOptions.ShareEnabled
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
	if createOfferingOptions.ProviderInfo != nil {
		body["provider_info"] = createOfferingOptions.ProviderInfo
	}
	if createOfferingOptions.RepoInfo != nil {
		body["repo_info"] = createOfferingOptions.RepoInfo
	}
	if createOfferingOptions.ImagePullKeys != nil {
		body["image_pull_keys"] = createOfferingOptions.ImagePullKeys
	}
	if createOfferingOptions.Support != nil {
		body["support"] = createOfferingOptions.Support
	}
	if createOfferingOptions.Media != nil {
		body["media"] = createOfferingOptions.Media
	}
	if createOfferingOptions.DeprecatePending != nil {
		body["deprecate_pending"] = createOfferingOptions.DeprecatePending
	}
	if createOfferingOptions.ProductKind != nil {
		body["product_kind"] = createOfferingOptions.ProductKind
	}
	if createOfferingOptions.Badges != nil {
		body["badges"] = createOfferingOptions.Badges
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
// Import new version to an offering.
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
	if importOfferingVersionOptions.XAuthToken != nil {
		builder.AddHeader("X-Auth-Token", fmt.Sprint(*importOfferingVersionOptions.XAuthToken))
	}

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
	if importOfferingVersionOptions.Repotype != nil {
		builder.AddQuery("repotype", fmt.Sprint(*importOfferingVersionOptions.Repotype))
	}

	body := make(map[string]interface{})
	if importOfferingVersionOptions.Tags != nil {
		body["tags"] = importOfferingVersionOptions.Tags
	}
	if importOfferingVersionOptions.Content != nil {
		body["content"] = importOfferingVersionOptions.Content
	}
	if importOfferingVersionOptions.Name != nil {
		body["name"] = importOfferingVersionOptions.Name
	}
	if importOfferingVersionOptions.Label != nil {
		body["label"] = importOfferingVersionOptions.Label
	}
	if importOfferingVersionOptions.InstallKind != nil {
		body["install_kind"] = importOfferingVersionOptions.InstallKind
	}
	if importOfferingVersionOptions.TargetKinds != nil {
		body["target_kinds"] = importOfferingVersionOptions.TargetKinds
	}
	if importOfferingVersionOptions.FormatKind != nil {
		body["format_kind"] = importOfferingVersionOptions.FormatKind
	}
	if importOfferingVersionOptions.ProductKind != nil {
		body["product_kind"] = importOfferingVersionOptions.ProductKind
	}
	if importOfferingVersionOptions.Sha != nil {
		body["sha"] = importOfferingVersionOptions.Sha
	}
	if importOfferingVersionOptions.Version != nil {
		body["version"] = importOfferingVersionOptions.Version
	}
	if importOfferingVersionOptions.Flavor != nil {
		body["flavor"] = importOfferingVersionOptions.Flavor
	}
	if importOfferingVersionOptions.Metadata != nil {
		body["metadata"] = importOfferingVersionOptions.Metadata
	}
	if importOfferingVersionOptions.WorkingDirectory != nil {
		body["working_directory"] = importOfferingVersionOptions.WorkingDirectory
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
// Import a new offering.
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
	if importOfferingOptions.Repotype != nil {
		builder.AddQuery("repotype", fmt.Sprint(*importOfferingOptions.Repotype))
	}

	body := make(map[string]interface{})
	if importOfferingOptions.Tags != nil {
		body["tags"] = importOfferingOptions.Tags
	}
	if importOfferingOptions.Content != nil {
		body["content"] = importOfferingOptions.Content
	}
	if importOfferingOptions.Name != nil {
		body["name"] = importOfferingOptions.Name
	}
	if importOfferingOptions.Label != nil {
		body["label"] = importOfferingOptions.Label
	}
	if importOfferingOptions.InstallKind != nil {
		body["install_kind"] = importOfferingOptions.InstallKind
	}
	if importOfferingOptions.TargetKinds != nil {
		body["target_kinds"] = importOfferingOptions.TargetKinds
	}
	if importOfferingOptions.FormatKind != nil {
		body["format_kind"] = importOfferingOptions.FormatKind
	}
	if importOfferingOptions.ProductKind != nil {
		body["product_kind"] = importOfferingOptions.ProductKind
	}
	if importOfferingOptions.Sha != nil {
		body["sha"] = importOfferingOptions.Sha
	}
	if importOfferingOptions.Version != nil {
		body["version"] = importOfferingOptions.Version
	}
	if importOfferingOptions.Flavor != nil {
		body["flavor"] = importOfferingOptions.Flavor
	}
	if importOfferingOptions.Metadata != nil {
		body["metadata"] = importOfferingOptions.Metadata
	}
	if importOfferingOptions.WorkingDirectory != nil {
		body["working_directory"] = importOfferingOptions.WorkingDirectory
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
	if reloadOfferingOptions.Content != nil {
		body["content"] = reloadOfferingOptions.Content
	}
	if reloadOfferingOptions.TargetKinds != nil {
		body["target_kinds"] = reloadOfferingOptions.TargetKinds
	}
	if reloadOfferingOptions.FormatKind != nil {
		body["format_kind"] = reloadOfferingOptions.FormatKind
	}
	if reloadOfferingOptions.Flavor != nil {
		body["flavor"] = reloadOfferingOptions.Flavor
	}
	if reloadOfferingOptions.WorkingDirectory != nil {
		body["working_directory"] = reloadOfferingOptions.WorkingDirectory
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

	if getOfferingOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*getOfferingOptions.Type))
	}
	if getOfferingOptions.Digest != nil {
		builder.AddQuery("digest", fmt.Sprint(*getOfferingOptions.Digest))
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
	if replaceOfferingOptions.LabelI18n != nil {
		body["label_i18n"] = replaceOfferingOptions.LabelI18n
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
	if replaceOfferingOptions.ShortDescriptionI18n != nil {
		body["short_description_i18n"] = replaceOfferingOptions.ShortDescriptionI18n
	}
	if replaceOfferingOptions.LongDescription != nil {
		body["long_description"] = replaceOfferingOptions.LongDescription
	}
	if replaceOfferingOptions.LongDescriptionI18n != nil {
		body["long_description_i18n"] = replaceOfferingOptions.LongDescriptionI18n
	}
	if replaceOfferingOptions.Features != nil {
		body["features"] = replaceOfferingOptions.Features
	}
	if replaceOfferingOptions.Kinds != nil {
		body["kinds"] = replaceOfferingOptions.Kinds
	}
	if replaceOfferingOptions.PcManaged != nil {
		body["pc_managed"] = replaceOfferingOptions.PcManaged
	}
	if replaceOfferingOptions.PublishApproved != nil {
		body["publish_approved"] = replaceOfferingOptions.PublishApproved
	}
	if replaceOfferingOptions.ShareWithAll != nil {
		body["share_with_all"] = replaceOfferingOptions.ShareWithAll
	}
	if replaceOfferingOptions.ShareWithIBM != nil {
		body["share_with_ibm"] = replaceOfferingOptions.ShareWithIBM
	}
	if replaceOfferingOptions.ShareEnabled != nil {
		body["share_enabled"] = replaceOfferingOptions.ShareEnabled
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
	if replaceOfferingOptions.ProviderInfo != nil {
		body["provider_info"] = replaceOfferingOptions.ProviderInfo
	}
	if replaceOfferingOptions.RepoInfo != nil {
		body["repo_info"] = replaceOfferingOptions.RepoInfo
	}
	if replaceOfferingOptions.ImagePullKeys != nil {
		body["image_pull_keys"] = replaceOfferingOptions.ImagePullKeys
	}
	if replaceOfferingOptions.Support != nil {
		body["support"] = replaceOfferingOptions.Support
	}
	if replaceOfferingOptions.Media != nil {
		body["media"] = replaceOfferingOptions.Media
	}
	if replaceOfferingOptions.DeprecatePending != nil {
		body["deprecate_pending"] = replaceOfferingOptions.DeprecatePending
	}
	if replaceOfferingOptions.ProductKind != nil {
		body["product_kind"] = replaceOfferingOptions.ProductKind
	}
	if replaceOfferingOptions.Badges != nil {
		body["badges"] = replaceOfferingOptions.Badges
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

// UpdateOffering : Update offering
// Update an offering.
func (catalogManagement *CatalogManagementV1) UpdateOffering(updateOfferingOptions *UpdateOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	return catalogManagement.UpdateOfferingWithContext(context.Background(), updateOfferingOptions)
}

// UpdateOfferingWithContext is an alternate form of the UpdateOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) UpdateOfferingWithContext(ctx context.Context, updateOfferingOptions *UpdateOfferingOptions) (result *Offering, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateOfferingOptions, "updateOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateOfferingOptions, "updateOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *updateOfferingOptions.CatalogIdentifier,
		"offering_id": *updateOfferingOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "UpdateOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json-patch+json")
	if updateOfferingOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateOfferingOptions.IfMatch))
	}

	_, err = builder.SetBodyContentJSON(updateOfferingOptions.Updates)
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

// ListOfferingAudits : Get offering audit logs
// Get the audit logs associated with an offering.
func (catalogManagement *CatalogManagementV1) ListOfferingAudits(listOfferingAuditsOptions *ListOfferingAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	return catalogManagement.ListOfferingAuditsWithContext(context.Background(), listOfferingAuditsOptions)
}

// ListOfferingAuditsWithContext is an alternate form of the ListOfferingAudits method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListOfferingAuditsWithContext(ctx context.Context, listOfferingAuditsOptions *ListOfferingAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listOfferingAuditsOptions, "listOfferingAuditsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listOfferingAuditsOptions, "listOfferingAuditsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *listOfferingAuditsOptions.CatalogIdentifier,
		"offering_id": *listOfferingAuditsOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/audits`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listOfferingAuditsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListOfferingAudits")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listOfferingAuditsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listOfferingAuditsOptions.Start))
	}
	if listOfferingAuditsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listOfferingAuditsOptions.Limit))
	}
	if listOfferingAuditsOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*listOfferingAuditsOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLogs)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOfferingAudit : Get an offering audit log entry
// Get the full audit log entry associated with an offering.
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
		"auditlog_identifier": *getOfferingAuditOptions.AuditlogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/audits/{auditlog_identifier}`, pathParamsMap)
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

	if getOfferingAuditOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*getOfferingAuditOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetOfferingPublish : Set offering publish approval settings
// Approve or disapprove the offering to be allowed to publish to the IBM Public Catalog. This is used only by Partner
// Center. Only users with Approval IAM authority can use this. Approvers should use the catalog and offering id from
// the public catalog since they wouldn't have access to the private offering.
func (catalogManagement *CatalogManagementV1) SetOfferingPublish(setOfferingPublishOptions *SetOfferingPublishOptions) (result *ApprovalResult, response *core.DetailedResponse, err error) {
	return catalogManagement.SetOfferingPublishWithContext(context.Background(), setOfferingPublishOptions)
}

// SetOfferingPublishWithContext is an alternate form of the SetOfferingPublish method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) SetOfferingPublishWithContext(ctx context.Context, setOfferingPublishOptions *SetOfferingPublishOptions) (result *ApprovalResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setOfferingPublishOptions, "setOfferingPublishOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setOfferingPublishOptions, "setOfferingPublishOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *setOfferingPublishOptions.CatalogIdentifier,
		"offering_id": *setOfferingPublishOptions.OfferingID,
		"approval_type": *setOfferingPublishOptions.ApprovalType,
		"approved": *setOfferingPublishOptions.Approved,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/publish/{approval_type}/{approved}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setOfferingPublishOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "SetOfferingPublish")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if setOfferingPublishOptions.XApproverToken != nil {
		builder.AddHeader("X-Approver-Token", fmt.Sprint(*setOfferingPublishOptions.XApproverToken))
	}

	if setOfferingPublishOptions.PortalRecord != nil {
		builder.AddQuery("portal_record", fmt.Sprint(*setOfferingPublishOptions.PortalRecord))
	}
	if setOfferingPublishOptions.PortalURL != nil {
		builder.AddQuery("portal_url", fmt.Sprint(*setOfferingPublishOptions.PortalURL))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApprovalResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeprecateOffering : Allows offering to be deprecated
// Approve or disapprove the offering to be deprecated.
func (catalogManagement *CatalogManagementV1) DeprecateOffering(deprecateOfferingOptions *DeprecateOfferingOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeprecateOfferingWithContext(context.Background(), deprecateOfferingOptions)
}

// DeprecateOfferingWithContext is an alternate form of the DeprecateOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeprecateOfferingWithContext(ctx context.Context, deprecateOfferingOptions *DeprecateOfferingOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deprecateOfferingOptions, "deprecateOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deprecateOfferingOptions, "deprecateOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *deprecateOfferingOptions.CatalogIdentifier,
		"offering_id": *deprecateOfferingOptions.OfferingID,
		"setting": *deprecateOfferingOptions.Setting,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/deprecate/{setting}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deprecateOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeprecateOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if deprecateOfferingOptions.Description != nil {
		body["description"] = deprecateOfferingOptions.Description
	}
	if deprecateOfferingOptions.DaysUntilDeprecate != nil {
		body["days_until_deprecate"] = deprecateOfferingOptions.DaysUntilDeprecate
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

// ShareOffering : Allows offering to be shared
// Set the share options on an offering.
func (catalogManagement *CatalogManagementV1) ShareOffering(shareOfferingOptions *ShareOfferingOptions) (result *ShareSetting, response *core.DetailedResponse, err error) {
	return catalogManagement.ShareOfferingWithContext(context.Background(), shareOfferingOptions)
}

// ShareOfferingWithContext is an alternate form of the ShareOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ShareOfferingWithContext(ctx context.Context, shareOfferingOptions *ShareOfferingOptions) (result *ShareSetting, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(shareOfferingOptions, "shareOfferingOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(shareOfferingOptions, "shareOfferingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *shareOfferingOptions.CatalogIdentifier,
		"offering_id": *shareOfferingOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/share`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range shareOfferingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ShareOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if shareOfferingOptions.IBM != nil {
		body["ibm"] = shareOfferingOptions.IBM
	}
	if shareOfferingOptions.Public != nil {
		body["public"] = shareOfferingOptions.Public
	}
	if shareOfferingOptions.Enabled != nil {
		body["enabled"] = shareOfferingOptions.Enabled
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalShareSetting)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOfferingAccess : Check for account ID in offering access list
// Determine if an account ID is in an offering's access list.
func (catalogManagement *CatalogManagementV1) GetOfferingAccess(getOfferingAccessOptions *GetOfferingAccessOptions) (result *Access, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingAccessWithContext(context.Background(), getOfferingAccessOptions)
}

// GetOfferingAccessWithContext is an alternate form of the GetOfferingAccess method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingAccessWithContext(ctx context.Context, getOfferingAccessOptions *GetOfferingAccessOptions) (result *Access, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingAccessOptions, "getOfferingAccessOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingAccessOptions, "getOfferingAccessOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getOfferingAccessOptions.CatalogIdentifier,
		"offering_id": *getOfferingAccessOptions.OfferingID,
		"access_identifier": *getOfferingAccessOptions.AccessIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/access/{access_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingAccessOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingAccess")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccess)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOfferingAccessList : Get offering access list
// Get the access list associated with the specified offering.
func (catalogManagement *CatalogManagementV1) GetOfferingAccessList(getOfferingAccessListOptions *GetOfferingAccessListOptions) (result *AccessListResult, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingAccessListWithContext(context.Background(), getOfferingAccessListOptions)
}

// GetOfferingAccessListWithContext is an alternate form of the GetOfferingAccessList method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingAccessListWithContext(ctx context.Context, getOfferingAccessListOptions *GetOfferingAccessListOptions) (result *AccessListResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingAccessListOptions, "getOfferingAccessListOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingAccessListOptions, "getOfferingAccessListOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getOfferingAccessListOptions.CatalogIdentifier,
		"offering_id": *getOfferingAccessListOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/access`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingAccessListOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingAccessList")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getOfferingAccessListOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*getOfferingAccessListOptions.Start))
	}
	if getOfferingAccessListOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getOfferingAccessListOptions.Limit))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessListResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteOfferingAccessList : Delete accesses from offering access list
// Delete all or a set of accesses from an offering's access list.
func (catalogManagement *CatalogManagementV1) DeleteOfferingAccessList(deleteOfferingAccessListOptions *DeleteOfferingAccessListOptions) (result *AccessListBulkResponse, response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteOfferingAccessListWithContext(context.Background(), deleteOfferingAccessListOptions)
}

// DeleteOfferingAccessListWithContext is an alternate form of the DeleteOfferingAccessList method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) DeleteOfferingAccessListWithContext(ctx context.Context, deleteOfferingAccessListOptions *DeleteOfferingAccessListOptions) (result *AccessListBulkResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteOfferingAccessListOptions, "deleteOfferingAccessListOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteOfferingAccessListOptions, "deleteOfferingAccessListOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *deleteOfferingAccessListOptions.CatalogIdentifier,
		"offering_id": *deleteOfferingAccessListOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/access`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteOfferingAccessListOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "DeleteOfferingAccessList")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(deleteOfferingAccessListOptions.Accesses)
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

// AddOfferingAccessList : Add accesses to offering access list
// Add one or more accesses to the specified offering's access list.
func (catalogManagement *CatalogManagementV1) AddOfferingAccessList(addOfferingAccessListOptions *AddOfferingAccessListOptions) (result *AccessListResult, response *core.DetailedResponse, err error) {
	return catalogManagement.AddOfferingAccessListWithContext(context.Background(), addOfferingAccessListOptions)
}

// AddOfferingAccessListWithContext is an alternate form of the AddOfferingAccessList method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) AddOfferingAccessListWithContext(ctx context.Context, addOfferingAccessListOptions *AddOfferingAccessListOptions) (result *AccessListResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(addOfferingAccessListOptions, "addOfferingAccessListOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(addOfferingAccessListOptions, "addOfferingAccessListOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *addOfferingAccessListOptions.CatalogIdentifier,
		"offering_id": *addOfferingAccessListOptions.OfferingID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/offerings/{offering_id}/access`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range addOfferingAccessListOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "AddOfferingAccessList")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(addOfferingAccessListOptions.Accesses)
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessListResult)
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
	if getOfferingUpdatesOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*getOfferingUpdatesOptions.XAuthRefreshToken))
	}

	builder.AddQuery("kind", fmt.Sprint(*getOfferingUpdatesOptions.Kind))
	if getOfferingUpdatesOptions.Target != nil {
		builder.AddQuery("target", fmt.Sprint(*getOfferingUpdatesOptions.Target))
	}
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
	if getOfferingUpdatesOptions.Sha != nil {
		builder.AddQuery("sha", fmt.Sprint(*getOfferingUpdatesOptions.Sha))
	}
	if getOfferingUpdatesOptions.Channel != nil {
		builder.AddQuery("channel", fmt.Sprint(*getOfferingUpdatesOptions.Channel))
	}
	if getOfferingUpdatesOptions.Namespaces != nil {
		builder.AddQuery("namespaces", strings.Join(getOfferingUpdatesOptions.Namespaces, ","))
	}
	if getOfferingUpdatesOptions.AllNamespaces != nil {
		builder.AddQuery("all_namespaces", fmt.Sprint(*getOfferingUpdatesOptions.AllNamespaces))
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

// GetOfferingSource : Get offering source
// Get an offering's source.  This request requires authorization, even for public offerings.
func (catalogManagement *CatalogManagementV1) GetOfferingSource(getOfferingSourceOptions *GetOfferingSourceOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingSourceWithContext(context.Background(), getOfferingSourceOptions)
}

// GetOfferingSourceWithContext is an alternate form of the GetOfferingSource method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingSourceWithContext(ctx context.Context, getOfferingSourceOptions *GetOfferingSourceOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingSourceOptions, "getOfferingSourceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingSourceOptions, "getOfferingSourceOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/offering/source`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingSourceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingSource")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/yaml")
	if getOfferingSourceOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*getOfferingSourceOptions.Accept))
	}

	builder.AddQuery("version", fmt.Sprint(*getOfferingSourceOptions.Version))
	if getOfferingSourceOptions.CatalogID != nil {
		builder.AddQuery("catalogID", fmt.Sprint(*getOfferingSourceOptions.CatalogID))
	}
	if getOfferingSourceOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*getOfferingSourceOptions.Name))
	}
	if getOfferingSourceOptions.ID != nil {
		builder.AddQuery("id", fmt.Sprint(*getOfferingSourceOptions.ID))
	}
	if getOfferingSourceOptions.Kind != nil {
		builder.AddQuery("kind", fmt.Sprint(*getOfferingSourceOptions.Kind))
	}
	if getOfferingSourceOptions.Channel != nil {
		builder.AddQuery("channel", fmt.Sprint(*getOfferingSourceOptions.Channel))
	}
	if getOfferingSourceOptions.Flavor != nil {
		builder.AddQuery("flavor", fmt.Sprint(*getOfferingSourceOptions.Flavor))
	}
	if getOfferingSourceOptions.AsIs != nil {
		builder.AddQuery("asIs", fmt.Sprint(*getOfferingSourceOptions.AsIs))
	}
	if getOfferingSourceOptions.InstallType != nil {
		builder.AddQuery("installType", fmt.Sprint(*getOfferingSourceOptions.InstallType))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, &result)

	return
}

// GetOfferingSourceURL : Get offering source URL
// Get an offering's private source image.
func (catalogManagement *CatalogManagementV1) GetOfferingSourceURL(getOfferingSourceURLOptions *GetOfferingSourceURLOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingSourceURLWithContext(context.Background(), getOfferingSourceURLOptions)
}

// GetOfferingSourceURLWithContext is an alternate form of the GetOfferingSourceURL method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingSourceURLWithContext(ctx context.Context, getOfferingSourceURLOptions *GetOfferingSourceURLOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingSourceURLOptions, "getOfferingSourceURLOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingSourceURLOptions, "getOfferingSourceURLOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"key_identifier": *getOfferingSourceURLOptions.KeyIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/offering/source/url/{key_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingSourceURLOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingSourceURL")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/yaml")
	if getOfferingSourceURLOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*getOfferingSourceURLOptions.Accept))
	}

	if getOfferingSourceURLOptions.CatalogID != nil {
		builder.AddQuery("catalogID", fmt.Sprint(*getOfferingSourceURLOptions.CatalogID))
	}
	if getOfferingSourceURLOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*getOfferingSourceURLOptions.Name))
	}
	if getOfferingSourceURLOptions.ID != nil {
		builder.AddQuery("id", fmt.Sprint(*getOfferingSourceURLOptions.ID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = catalogManagement.Service.Request(request, &result)

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

// ArchiveVersion : Archive version immediately
// Archive the specified version.
func (catalogManagement *CatalogManagementV1) ArchiveVersion(archiveVersionOptions *ArchiveVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.ArchiveVersionWithContext(context.Background(), archiveVersionOptions)
}

// ArchiveVersionWithContext is an alternate form of the ArchiveVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ArchiveVersionWithContext(ctx context.Context, archiveVersionOptions *ArchiveVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(archiveVersionOptions, "archiveVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(archiveVersionOptions, "archiveVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *archiveVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/archive`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range archiveVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ArchiveVersion")
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

// SetDeprecateVersion : Sets version to be deprecated in a certain time period
// Set or cancel the version to be deprecated.
func (catalogManagement *CatalogManagementV1) SetDeprecateVersion(setDeprecateVersionOptions *SetDeprecateVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.SetDeprecateVersionWithContext(context.Background(), setDeprecateVersionOptions)
}

// SetDeprecateVersionWithContext is an alternate form of the SetDeprecateVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) SetDeprecateVersionWithContext(ctx context.Context, setDeprecateVersionOptions *SetDeprecateVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setDeprecateVersionOptions, "setDeprecateVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setDeprecateVersionOptions, "setDeprecateVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *setDeprecateVersionOptions.VersionLocID,
		"setting": *setDeprecateVersionOptions.Setting,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/deprecate/{setting}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setDeprecateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "SetDeprecateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setDeprecateVersionOptions.Description != nil {
		body["description"] = setDeprecateVersionOptions.Description
	}
	if setDeprecateVersionOptions.DaysUntilDeprecate != nil {
		body["days_until_deprecate"] = setDeprecateVersionOptions.DaysUntilDeprecate
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

// ConsumableVersion : Make version consumable for sharing
// Set the version as consumable in order to inherit the offering sharing permissions.
func (catalogManagement *CatalogManagementV1) ConsumableVersion(consumableVersionOptions *ConsumableVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.ConsumableVersionWithContext(context.Background(), consumableVersionOptions)
}

// ConsumableVersionWithContext is an alternate form of the ConsumableVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ConsumableVersionWithContext(ctx context.Context, consumableVersionOptions *ConsumableVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(consumableVersionOptions, "consumableVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(consumableVersionOptions, "consumableVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *consumableVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/consume-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range consumableVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ConsumableVersion")
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

// SuspendVersion : Suspend a version
// Limits the visibility of a version by moving a version state from consumable back to validated.
func (catalogManagement *CatalogManagementV1) SuspendVersion(suspendVersionOptions *SuspendVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.SuspendVersionWithContext(context.Background(), suspendVersionOptions)
}

// SuspendVersionWithContext is an alternate form of the SuspendVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) SuspendVersionWithContext(ctx context.Context, suspendVersionOptions *SuspendVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(suspendVersionOptions, "suspendVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(suspendVersionOptions, "suspendVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *suspendVersionOptions.VersionLocID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/suspend`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range suspendVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "SuspendVersion")
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
	if copyVersionOptions.Content != nil {
		body["content"] = copyVersionOptions.Content
	}
	if copyVersionOptions.TargetKinds != nil {
		body["target_kinds"] = copyVersionOptions.TargetKinds
	}
	if copyVersionOptions.FormatKind != nil {
		body["format_kind"] = copyVersionOptions.FormatKind
	}
	if copyVersionOptions.Flavor != nil {
		body["flavor"] = copyVersionOptions.Flavor
	}
	if copyVersionOptions.WorkingDirectory != nil {
		body["working_directory"] = copyVersionOptions.WorkingDirectory
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

// CopyFromPreviousVersion : Copy values from a previous version
// Copy values from a specified previous version.
func (catalogManagement *CatalogManagementV1) CopyFromPreviousVersion(copyFromPreviousVersionOptions *CopyFromPreviousVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.CopyFromPreviousVersionWithContext(context.Background(), copyFromPreviousVersionOptions)
}

// CopyFromPreviousVersionWithContext is an alternate form of the CopyFromPreviousVersion method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) CopyFromPreviousVersionWithContext(ctx context.Context, copyFromPreviousVersionOptions *CopyFromPreviousVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(copyFromPreviousVersionOptions, "copyFromPreviousVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(copyFromPreviousVersionOptions, "copyFromPreviousVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"version_loc_id": *copyFromPreviousVersionOptions.VersionLocID,
		"type": *copyFromPreviousVersionOptions.Type,
		"version_loc_id_to_copy_from": *copyFromPreviousVersionOptions.VersionLocIDToCopyFrom,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/versions/{version_loc_id}/copy/{type}/{version_loc_id_to_copy_from}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range copyFromPreviousVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "CopyFromPreviousVersion")
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

// DeprecateVersion : Deprecate version immediately - use /archive instead
// Deprecate the specified version.
// Deprecated: this method is deprecated and may be removed in a future release.
func (catalogManagement *CatalogManagementV1) DeprecateVersion(deprecateVersionOptions *DeprecateVersionOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeprecateVersionWithContext(context.Background(), deprecateVersionOptions)
}

// DeprecateVersionWithContext is an alternate form of the DeprecateVersion method which supports a Context parameter
// Deprecated: this method is deprecated and may be removed in a future release.
func (catalogManagement *CatalogManagementV1) DeprecateVersionWithContext(ctx context.Context, deprecateVersionOptions *DeprecateVersionOptions) (response *core.DetailedResponse, err error) {
	core.GetLogger().Warn("A deprecated operation has been invoked: DeprecateVersion")
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
	if deployOperatorsOptions.Channel != nil {
		body["channel"] = deployOperatorsOptions.Channel
	}
	if deployOperatorsOptions.InstallPlan != nil {
		body["install_plan"] = deployOperatorsOptions.InstallPlan
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
	if replaceOperatorsOptions.Channel != nil {
		body["channel"] = replaceOperatorsOptions.Channel
	}
	if replaceOperatorsOptions.InstallPlan != nil {
		body["install_plan"] = replaceOperatorsOptions.InstallPlan
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
	if installVersionOptions.EnvironmentVariables != nil {
		body["environment_variables"] = installVersionOptions.EnvironmentVariables
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
	if installVersionOptions.VcenterLocation != nil {
		body["vcenter_location"] = installVersionOptions.VcenterLocation
	}
	if installVersionOptions.VcenterUser != nil {
		body["vcenter_user"] = installVersionOptions.VcenterUser
	}
	if installVersionOptions.VcenterPassword != nil {
		body["vcenter_password"] = installVersionOptions.VcenterPassword
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
	if preinstallVersionOptions.EnvironmentVariables != nil {
		body["environment_variables"] = preinstallVersionOptions.EnvironmentVariables
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
	if preinstallVersionOptions.VcenterLocation != nil {
		body["vcenter_location"] = preinstallVersionOptions.VcenterLocation
	}
	if preinstallVersionOptions.VcenterUser != nil {
		body["vcenter_user"] = preinstallVersionOptions.VcenterUser
	}
	if preinstallVersionOptions.VcenterPassword != nil {
		body["vcenter_password"] = preinstallVersionOptions.VcenterPassword
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
	if validateInstallOptions.EnvironmentVariables != nil {
		body["environment_variables"] = validateInstallOptions.EnvironmentVariables
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
	if validateInstallOptions.VcenterLocation != nil {
		body["vcenter_location"] = validateInstallOptions.VcenterLocation
	}
	if validateInstallOptions.VcenterUser != nil {
		body["vcenter_user"] = validateInstallOptions.VcenterUser
	}
	if validateInstallOptions.VcenterPassword != nil {
		body["vcenter_password"] = validateInstallOptions.VcenterPassword
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
	if searchObjectsOptions.Kind != nil {
		builder.AddQuery("kind", fmt.Sprint(*searchObjectsOptions.Kind))
	}
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
	if createObjectOptions.Name != nil {
		body["name"] = createObjectOptions.Name
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
	if replaceObjectOptions.Rev != nil {
		body["_rev"] = replaceObjectOptions.Rev
	}
	if replaceObjectOptions.Name != nil {
		body["name"] = replaceObjectOptions.Name
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

// ListObjectAudits : Get object audit logs
// Get the audit logs associated with an object.
func (catalogManagement *CatalogManagementV1) ListObjectAudits(listObjectAuditsOptions *ListObjectAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	return catalogManagement.ListObjectAuditsWithContext(context.Background(), listObjectAuditsOptions)
}

// ListObjectAuditsWithContext is an alternate form of the ListObjectAudits method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListObjectAuditsWithContext(ctx context.Context, listObjectAuditsOptions *ListObjectAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listObjectAuditsOptions, "listObjectAuditsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listObjectAuditsOptions, "listObjectAuditsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *listObjectAuditsOptions.CatalogIdentifier,
		"object_identifier": *listObjectAuditsOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/audits`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listObjectAuditsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListObjectAudits")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listObjectAuditsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listObjectAuditsOptions.Start))
	}
	if listObjectAuditsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listObjectAuditsOptions.Limit))
	}
	if listObjectAuditsOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*listObjectAuditsOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLogs)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetObjectAudit : Get an object audit log entry
// Get the full audit log entry associated with an object.
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
		"auditlog_identifier": *getObjectAuditOptions.AuditlogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/audits/{auditlog_identifier}`, pathParamsMap)
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

	if getObjectAuditOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*getObjectAuditOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ConsumableShareObject : Make object consumable for sharing
// Set the object as consumable in order to use the object sharing permissions.
func (catalogManagement *CatalogManagementV1) ConsumableShareObject(consumableShareObjectOptions *ConsumableShareObjectOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.ConsumableShareObjectWithContext(context.Background(), consumableShareObjectOptions)
}

// ConsumableShareObjectWithContext is an alternate form of the ConsumableShareObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ConsumableShareObjectWithContext(ctx context.Context, consumableShareObjectOptions *ConsumableShareObjectOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(consumableShareObjectOptions, "consumableShareObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(consumableShareObjectOptions, "consumableShareObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *consumableShareObjectOptions.CatalogIdentifier,
		"object_identifier": *consumableShareObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/consume-publish`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range consumableShareObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ConsumableShareObject")
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

// ShareObject : Allows object to be shared
// Set the share options on an object.
func (catalogManagement *CatalogManagementV1) ShareObject(shareObjectOptions *ShareObjectOptions) (result *ShareSetting, response *core.DetailedResponse, err error) {
	return catalogManagement.ShareObjectWithContext(context.Background(), shareObjectOptions)
}

// ShareObjectWithContext is an alternate form of the ShareObject method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ShareObjectWithContext(ctx context.Context, shareObjectOptions *ShareObjectOptions) (result *ShareSetting, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(shareObjectOptions, "shareObjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(shareObjectOptions, "shareObjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *shareObjectOptions.CatalogIdentifier,
		"object_identifier": *shareObjectOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/share`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range shareObjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ShareObject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if shareObjectOptions.IBM != nil {
		body["ibm"] = shareObjectOptions.IBM
	}
	if shareObjectOptions.Public != nil {
		body["public"] = shareObjectOptions.Public
	}
	if shareObjectOptions.Enabled != nil {
		body["enabled"] = shareObjectOptions.Enabled
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalShareSetting)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetObjectAccessList : Get object access list
// Get the access list associated with the specified object.
func (catalogManagement *CatalogManagementV1) GetObjectAccessList(getObjectAccessListOptions *GetObjectAccessListOptions) (result *AccessListResult, response *core.DetailedResponse, err error) {
	return catalogManagement.GetObjectAccessListWithContext(context.Background(), getObjectAccessListOptions)
}

// GetObjectAccessListWithContext is an alternate form of the GetObjectAccessList method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetObjectAccessListWithContext(ctx context.Context, getObjectAccessListOptions *GetObjectAccessListOptions) (result *AccessListResult, response *core.DetailedResponse, err error) {
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
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/accessv1`, pathParamsMap)
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

	if getObjectAccessListOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*getObjectAccessListOptions.Start))
	}
	if getObjectAccessListOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getObjectAccessListOptions.Limit))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessListResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetObjectAccess : Check for account ID in object access list
// Determine if an account ID is in an object's access list.
func (catalogManagement *CatalogManagementV1) GetObjectAccess(getObjectAccessOptions *GetObjectAccessOptions) (result *Access, response *core.DetailedResponse, err error) {
	return catalogManagement.GetObjectAccessWithContext(context.Background(), getObjectAccessOptions)
}

// GetObjectAccessWithContext is an alternate form of the GetObjectAccess method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetObjectAccessWithContext(ctx context.Context, getObjectAccessOptions *GetObjectAccessOptions) (result *Access, response *core.DetailedResponse, err error) {
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
		"access_identifier": *getObjectAccessOptions.AccessIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access/{access_identifier}`, pathParamsMap)
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccess)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateObjectAccess : Add account ID to object access list
// Add an account ID to an object's access list.
// Deprecated: this method is deprecated and may be removed in a future release.
func (catalogManagement *CatalogManagementV1) CreateObjectAccess(createObjectAccessOptions *CreateObjectAccessOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.CreateObjectAccessWithContext(context.Background(), createObjectAccessOptions)
}

// CreateObjectAccessWithContext is an alternate form of the CreateObjectAccess method which supports a Context parameter
// Deprecated: this method is deprecated and may be removed in a future release.
func (catalogManagement *CatalogManagementV1) CreateObjectAccessWithContext(ctx context.Context, createObjectAccessOptions *CreateObjectAccessOptions) (response *core.DetailedResponse, err error) {
	core.GetLogger().Warn("A deprecated operation has been invoked: CreateObjectAccess")
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
		"access_identifier": *createObjectAccessOptions.AccessIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access/{access_identifier}`, pathParamsMap)
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

// DeleteObjectAccess : Remove account ID from object access list
// Delete the specified account ID from the specified object's access list.
// Deprecated: this method is deprecated and may be removed in a future release.
func (catalogManagement *CatalogManagementV1) DeleteObjectAccess(deleteObjectAccessOptions *DeleteObjectAccessOptions) (response *core.DetailedResponse, err error) {
	return catalogManagement.DeleteObjectAccessWithContext(context.Background(), deleteObjectAccessOptions)
}

// DeleteObjectAccessWithContext is an alternate form of the DeleteObjectAccess method which supports a Context parameter
// Deprecated: this method is deprecated and may be removed in a future release.
func (catalogManagement *CatalogManagementV1) DeleteObjectAccessWithContext(ctx context.Context, deleteObjectAccessOptions *DeleteObjectAccessOptions) (response *core.DetailedResponse, err error) {
	core.GetLogger().Warn("A deprecated operation has been invoked: DeleteObjectAccess")
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
		"access_identifier": *deleteObjectAccessOptions.AccessIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access/{access_identifier}`, pathParamsMap)
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

// GetObjectAccessListDeprecated : Get object access list
// Deprecated - use /accessv1 instead.
// Deprecated: this method is deprecated and may be removed in a future release.
func (catalogManagement *CatalogManagementV1) GetObjectAccessListDeprecated(getObjectAccessListDeprecatedOptions *GetObjectAccessListDeprecatedOptions) (result *ObjectAccessListResult, response *core.DetailedResponse, err error) {
	return catalogManagement.GetObjectAccessListDeprecatedWithContext(context.Background(), getObjectAccessListDeprecatedOptions)
}

// GetObjectAccessListDeprecatedWithContext is an alternate form of the GetObjectAccessListDeprecated method which supports a Context parameter
// Deprecated: this method is deprecated and may be removed in a future release.
func (catalogManagement *CatalogManagementV1) GetObjectAccessListDeprecatedWithContext(ctx context.Context, getObjectAccessListDeprecatedOptions *GetObjectAccessListDeprecatedOptions) (result *ObjectAccessListResult, response *core.DetailedResponse, err error) {
	core.GetLogger().Warn("A deprecated operation has been invoked: GetObjectAccessListDeprecated")
	err = core.ValidateNotNil(getObjectAccessListDeprecatedOptions, "getObjectAccessListDeprecatedOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getObjectAccessListDeprecatedOptions, "getObjectAccessListDeprecatedOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"catalog_identifier": *getObjectAccessListDeprecatedOptions.CatalogIdentifier,
		"object_identifier": *getObjectAccessListDeprecatedOptions.ObjectIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalog_identifier}/objects/{object_identifier}/access`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getObjectAccessListDeprecatedOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetObjectAccessListDeprecated")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getObjectAccessListDeprecatedOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getObjectAccessListDeprecatedOptions.Limit))
	}
	if getObjectAccessListDeprecatedOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*getObjectAccessListDeprecatedOptions.Offset))
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

// DeleteObjectAccessList : Delete accesses from object access list
// Delete all or a set of accesses from an object's access list.
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

	_, err = builder.SetBodyContentJSON(deleteObjectAccessListOptions.Accesses)
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

// AddObjectAccessList : Add accesses to object access list
// Add one or more accesses to the specified object's access list.
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

	_, err = builder.SetBodyContentJSON(addObjectAccessListOptions.Accesses)
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
	if createOfferingInstanceOptions.VersionID != nil {
		body["version_id"] = createOfferingInstanceOptions.VersionID
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
	if createOfferingInstanceOptions.InstallPlan != nil {
		body["install_plan"] = createOfferingInstanceOptions.InstallPlan
	}
	if createOfferingInstanceOptions.Channel != nil {
		body["channel"] = createOfferingInstanceOptions.Channel
	}
	if createOfferingInstanceOptions.Created != nil {
		body["created"] = createOfferingInstanceOptions.Created
	}
	if createOfferingInstanceOptions.Updated != nil {
		body["updated"] = createOfferingInstanceOptions.Updated
	}
	if createOfferingInstanceOptions.Metadata != nil {
		body["metadata"] = createOfferingInstanceOptions.Metadata
	}
	if createOfferingInstanceOptions.ResourceGroupID != nil {
		body["resource_group_id"] = createOfferingInstanceOptions.ResourceGroupID
	}
	if createOfferingInstanceOptions.Location != nil {
		body["location"] = createOfferingInstanceOptions.Location
	}
	if createOfferingInstanceOptions.Disabled != nil {
		body["disabled"] = createOfferingInstanceOptions.Disabled
	}
	if createOfferingInstanceOptions.Account != nil {
		body["account"] = createOfferingInstanceOptions.Account
	}
	if createOfferingInstanceOptions.LastOperation != nil {
		body["last_operation"] = createOfferingInstanceOptions.LastOperation
	}
	if createOfferingInstanceOptions.KindTarget != nil {
		body["kind_target"] = createOfferingInstanceOptions.KindTarget
	}
	if createOfferingInstanceOptions.Sha != nil {
		body["sha"] = createOfferingInstanceOptions.Sha
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
	if putOfferingInstanceOptions.VersionID != nil {
		body["version_id"] = putOfferingInstanceOptions.VersionID
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
	if putOfferingInstanceOptions.InstallPlan != nil {
		body["install_plan"] = putOfferingInstanceOptions.InstallPlan
	}
	if putOfferingInstanceOptions.Channel != nil {
		body["channel"] = putOfferingInstanceOptions.Channel
	}
	if putOfferingInstanceOptions.Created != nil {
		body["created"] = putOfferingInstanceOptions.Created
	}
	if putOfferingInstanceOptions.Updated != nil {
		body["updated"] = putOfferingInstanceOptions.Updated
	}
	if putOfferingInstanceOptions.Metadata != nil {
		body["metadata"] = putOfferingInstanceOptions.Metadata
	}
	if putOfferingInstanceOptions.ResourceGroupID != nil {
		body["resource_group_id"] = putOfferingInstanceOptions.ResourceGroupID
	}
	if putOfferingInstanceOptions.Location != nil {
		body["location"] = putOfferingInstanceOptions.Location
	}
	if putOfferingInstanceOptions.Disabled != nil {
		body["disabled"] = putOfferingInstanceOptions.Disabled
	}
	if putOfferingInstanceOptions.Account != nil {
		body["account"] = putOfferingInstanceOptions.Account
	}
	if putOfferingInstanceOptions.LastOperation != nil {
		body["last_operation"] = putOfferingInstanceOptions.LastOperation
	}
	if putOfferingInstanceOptions.KindTarget != nil {
		body["kind_target"] = putOfferingInstanceOptions.KindTarget
	}
	if putOfferingInstanceOptions.Sha != nil {
		body["sha"] = putOfferingInstanceOptions.Sha
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

// ListOfferingInstanceAudits : Get offering instance audit logs
// Get the audit logs associated with an offering instance.
func (catalogManagement *CatalogManagementV1) ListOfferingInstanceAudits(listOfferingInstanceAuditsOptions *ListOfferingInstanceAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	return catalogManagement.ListOfferingInstanceAuditsWithContext(context.Background(), listOfferingInstanceAuditsOptions)
}

// ListOfferingInstanceAuditsWithContext is an alternate form of the ListOfferingInstanceAudits method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) ListOfferingInstanceAuditsWithContext(ctx context.Context, listOfferingInstanceAuditsOptions *ListOfferingInstanceAuditsOptions) (result *AuditLogs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listOfferingInstanceAuditsOptions, "listOfferingInstanceAuditsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listOfferingInstanceAuditsOptions, "listOfferingInstanceAuditsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_identifier": *listOfferingInstanceAuditsOptions.InstanceIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/instances/offerings/{instance_identifier}/audits`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listOfferingInstanceAuditsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "ListOfferingInstanceAudits")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listOfferingInstanceAuditsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listOfferingInstanceAuditsOptions.Start))
	}
	if listOfferingInstanceAuditsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listOfferingInstanceAuditsOptions.Limit))
	}
	if listOfferingInstanceAuditsOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*listOfferingInstanceAuditsOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLogs)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOfferingInstanceAudit : Get an offering instance audit log entry
// Get the full audit log entry associated with an offering instance.
func (catalogManagement *CatalogManagementV1) GetOfferingInstanceAudit(getOfferingInstanceAuditOptions *GetOfferingInstanceAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	return catalogManagement.GetOfferingInstanceAuditWithContext(context.Background(), getOfferingInstanceAuditOptions)
}

// GetOfferingInstanceAuditWithContext is an alternate form of the GetOfferingInstanceAudit method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) GetOfferingInstanceAuditWithContext(ctx context.Context, getOfferingInstanceAuditOptions *GetOfferingInstanceAuditOptions) (result *AuditLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOfferingInstanceAuditOptions, "getOfferingInstanceAuditOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOfferingInstanceAuditOptions, "getOfferingInstanceAuditOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_identifier": *getOfferingInstanceAuditOptions.InstanceIdentifier,
		"auditlog_identifier": *getOfferingInstanceAuditOptions.AuditlogIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/instances/offerings/{instance_identifier}/audits/{auditlog_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOfferingInstanceAuditOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "GetOfferingInstanceAudit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getOfferingInstanceAuditOptions.Lookupnames != nil {
		builder.AddQuery("lookupnames", fmt.Sprint(*getOfferingInstanceAuditOptions.Lookupnames))
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuditLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// Access : access.
type Access struct {
	// unique id.
	ID *string `json:"id,omitempty"`

	// account id.
	Account *string `json:"account,omitempty"`

	// Normal account or enterprise.
	AccountType *int64 `json:"account_type,omitempty"`

	// unique id.
	CatalogID *string `json:"catalog_id,omitempty"`

	// object ID.
	TargetID *string `json:"target_id,omitempty"`

	// object's owner's account.
	TargetAccount *string `json:"target_account,omitempty"`

	// entity type.
	TargetKind *string `json:"target_kind,omitempty"`

	// accessible to the private object.
	PrivateAccessible *bool `json:"private_accessible,omitempty"`

	// date and time create.
	Created *strfmt.DateTime `json:"created,omitempty"`
}

// UnmarshalAccess unmarshals an instance of Access from the specified map of raw messages.
func UnmarshalAccess(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Access)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_type", &obj.AccountType)
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
	err = core.UnmarshalPrimitive(m, "target_account", &obj.TargetAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_kind", &obj.TargetKind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_accessible", &obj.PrivateAccessible)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
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

// AccessListResult : Paginated Offering search result.
type AccessListResult struct {
	// The start token used for this response.
	Start *string `json:"start,omitempty"`

	// The limit that was applied to this response. It may be smaller than in the request because that was too large.
	Limit *int64 `json:"limit" validate:"required"`

	// The total count of resources in the system that matches the request.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of resources returned in this response.
	ResourceCount *int64 `json:"resource_count" validate:"required"`

	// Link response on a token paginated query.
	First *PaginationTokenLink `json:"first" validate:"required"`

	// Link response on a token paginated query.
	Next *PaginationTokenLink `json:"next,omitempty"`

	// Link response on a token paginated query.
	Prev *PaginationTokenLink `json:"prev,omitempty"`

	// Link response on a token paginated query.
	Last *PaginationTokenLink `json:"last,omitempty"`

	// A list of access records.
	Resources []Access `json:"resources" validate:"required"`
}

// UnmarshalAccessListResult unmarshals an instance of AccessListResult from the specified map of raw messages.
func UnmarshalAccessListResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessListResult)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationTokenLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationTokenLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "prev", &obj.Prev, UnmarshalPaginationTokenLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationTokenLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalAccess)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AccessListResult) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// Account : Account information.
type Account struct {
	// Account identification.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

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
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
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

// AccumulatedFilters : The accumulated filters for an account. This will return the account filters plus a filter for each catalog the user
// has access to.
type AccumulatedFilters struct {
	// Hide the public catalog in this account.
	HideIBMCloudCatalog *bool `json:"hide_IBM_cloud_catalog,omitempty"`

	// Filters for accounts (at this time this will always be just one item array).
	AccountFilters []Filters `json:"account_filters,omitempty"`

	// The filters for all of the accessible catalogs.
	CatalogFilters []AccumulatedFiltersCatalogFiltersItem `json:"catalog_filters,omitempty"`
}

// UnmarshalAccumulatedFilters unmarshals an instance of AccumulatedFilters from the specified map of raw messages.
func UnmarshalAccumulatedFilters(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccumulatedFilters)
	err = core.UnmarshalPrimitive(m, "hide_IBM_cloud_catalog", &obj.HideIBMCloudCatalog)
	if err != nil {
		return
	}
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
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// A list of accesses to add.
	Accesses []string `json:"accesses" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewAddObjectAccessListOptions : Instantiate AddObjectAccessListOptions
func (*CatalogManagementV1) NewAddObjectAccessListOptions(catalogIdentifier string, objectIdentifier string, accesses []string) *AddObjectAccessListOptions {
	return &AddObjectAccessListOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		Accesses: accesses,
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *AddObjectAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *AddObjectAccessListOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *AddObjectAccessListOptions) SetObjectIdentifier(objectIdentifier string) *AddObjectAccessListOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetAccesses : Allow user to set Accesses
func (_options *AddObjectAccessListOptions) SetAccesses(accesses []string) *AddObjectAccessListOptions {
	_options.Accesses = accesses
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *AddObjectAccessListOptions) SetHeaders(param map[string]string) *AddObjectAccessListOptions {
	options.Headers = param
	return options
}

// AddOfferingAccessListOptions : The AddOfferingAccessList options.
type AddOfferingAccessListOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// A list of accesses to add.
	Accesses []string `json:"accesses" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewAddOfferingAccessListOptions : Instantiate AddOfferingAccessListOptions
func (*CatalogManagementV1) NewAddOfferingAccessListOptions(catalogIdentifier string, offeringID string, accesses []string) *AddOfferingAccessListOptions {
	return &AddOfferingAccessListOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		Accesses: accesses,
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *AddOfferingAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *AddOfferingAccessListOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *AddOfferingAccessListOptions) SetOfferingID(offeringID string) *AddOfferingAccessListOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetAccesses : Allow user to set Accesses
func (_options *AddOfferingAccessListOptions) SetAccesses(accesses []string) *AddOfferingAccessListOptions {
	_options.Accesses = accesses
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *AddOfferingAccessListOptions) SetHeaders(param map[string]string) *AddOfferingAccessListOptions {
	options.Headers = param
	return options
}

// ApprovalResult : Result of approval.
type ApprovalResult struct {
	// Shared - object is shared using access list - not set when using PC Managed objects.
	// Deprecated: this field is deprecated and may be removed in a future release.
	Shared *bool `json:"shared,omitempty"`

	// Shared with IBM only - access list is also applicable - not set when using PC Managed objects.
	// Deprecated: this field is deprecated and may be removed in a future release.
	IBM *bool `json:"ibm,omitempty"`

	// Shared with everyone - not set when using PC Managed objects.
	// Deprecated: this field is deprecated and may be removed in a future release.
	Public *bool `json:"public,omitempty"`

	// Published to Partner Center (pc_managed) or for objects, allowed to request publishing.
	AllowRequest *bool `json:"allow_request,omitempty"`

	// Approvers have approved publishing to public catalog.
	Approved *bool `json:"approved,omitempty"`

	// Partner Center document ID.
	PortalRecord *string `json:"portal_record,omitempty"`

	// Partner Center URL for this product.
	PortalURL *string `json:"portal_url,omitempty"`

	// Denotes whether approvals have changed.
	Changed *bool `json:"changed,omitempty"`
}

// UnmarshalApprovalResult unmarshals an instance of ApprovalResult from the specified map of raw messages.
func UnmarshalApprovalResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApprovalResult)
	err = core.UnmarshalPrimitive(m, "shared", &obj.Shared)
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
	err = core.UnmarshalPrimitive(m, "allow_request", &obj.AllowRequest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "approved", &obj.Approved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "portal_record", &obj.PortalRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "portal_url", &obj.PortalURL)
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

// ArchitectureDiagram : An Architecture Diagram.
type ArchitectureDiagram struct {
	// Offering Media information.
	Diagram *MediaItem `json:"diagram,omitempty"`

	// Description of this diagram.
	Description *string `json:"description,omitempty"`

	// A map of translated strings, by language code.
	DescriptionI18n map[string]string `json:"description_i18n,omitempty"`
}

// UnmarshalArchitectureDiagram unmarshals an instance of ArchitectureDiagram from the specified map of raw messages.
func UnmarshalArchitectureDiagram(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ArchitectureDiagram)
	err = core.UnmarshalModel(m, "diagram", &obj.Diagram, UnmarshalMediaItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description_i18n", &obj.DescriptionI18n)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ArchiveVersionOptions : The ArchiveVersion options.
type ArchiveVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewArchiveVersionOptions : Instantiate ArchiveVersionOptions
func (*CatalogManagementV1) NewArchiveVersionOptions(versionLocID string) *ArchiveVersionOptions {
	return &ArchiveVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (_options *ArchiveVersionOptions) SetVersionLocID(versionLocID string) *ArchiveVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ArchiveVersionOptions) SetHeaders(param map[string]string) *ArchiveVersionOptions {
	options.Headers = param
	return options
}

// AuditLog : An audit log which describes a change made to a catalog or associated resource.
type AuditLog struct {
	// The identifier of the audit record.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// The time at which the change was made.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The type of change described by the audit record.
	ChangeType *string `json:"change_type,omitempty"`

	// The resource type associated with the change.
	TargetType *string `json:"target_type,omitempty"`

	// The identifier of the resource that was changed.
	TargetID *string `json:"target_id,omitempty"`

	// The email address of the user that made the change.
	WhoEmail *string `json:"who_email,omitempty"`

	// The email address of the delegate user that made the change. This happens when a service makes a change onbehalf of
	// the user.
	WhoDelegateEmail *string `json:"who_delegate_email,omitempty"`

	// A message which describes the change.
	Message *string `json:"message,omitempty"`

	// Transaction id for this change.
	Gid *string `json:"gid,omitempty"`

	// IAM identifier of the user who made the change.
	WhoID *string `json:"who_id,omitempty"`

	// Name of the user who made the change.
	WhoName *string `json:"who_name,omitempty"`

	// IAM identifier of the delegate user who made the change.
	WhoDelegateID *string `json:"who_delegate_id,omitempty"`

	// Name of the delegate user who made the change.
	WhoDelegateName *string `json:"who_delegate_name,omitempty"`

	// Data about the change. Usually a change log of what was changed, both before and after. Can be of any type.
	Data interface{} `json:"data,omitempty"`
}

// UnmarshalAuditLog unmarshals an instance of AuditLog from the specified map of raw messages.
func UnmarshalAuditLog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AuditLog)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
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
	err = core.UnmarshalPrimitive(m, "who_email", &obj.WhoEmail)
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
	err = core.UnmarshalPrimitive(m, "gid", &obj.Gid)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "who_id", &obj.WhoID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "who_name", &obj.WhoName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "who_delegate_id", &obj.WhoDelegateID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "who_delegate_name", &obj.WhoDelegateName)
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

// AuditLogDigest : An reduced audit log which describes a change made to a catalog or associated resource.
type AuditLogDigest struct {
	// The identifier of the audit record.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// The time at which the change was made.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The type of change described by the audit record.
	ChangeType *string `json:"change_type,omitempty"`

	// The resource type associated with the change.
	TargetType *string `json:"target_type,omitempty"`

	// The identifier of the resource that was changed.
	TargetID *string `json:"target_id,omitempty"`

	// The email address of the user that made the change.
	WhoEmail *string `json:"who_email,omitempty"`

	// The email address of the delegate user that made the change. This happens when a service makes a change onbehalf of
	// the user.
	WhoDelegateEmail *string `json:"who_delegate_email,omitempty"`

	// A message which describes the change.
	Message *string `json:"message,omitempty"`
}

// UnmarshalAuditLogDigest unmarshals an instance of AuditLogDigest from the specified map of raw messages.
func UnmarshalAuditLogDigest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AuditLogDigest)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
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
	err = core.UnmarshalPrimitive(m, "who_email", &obj.WhoEmail)
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

// AuditLogs : A collection of audit records.
type AuditLogs struct {
	// The start token used for this response.
	Start *string `json:"start,omitempty"`

	// The limit that was applied to this response. It may be smaller than in the request because that was too large.
	Limit *int64 `json:"limit" validate:"required"`

	// The total count of resources in the system that matches the request.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of resources returned in this response.
	ResourceCount *int64 `json:"resource_count" validate:"required"`

	// Link response on a token paginated query.
	First *PaginationTokenLink `json:"first" validate:"required"`

	// Link response on a token paginated query.
	Next *PaginationTokenLink `json:"next,omitempty"`

	// Link response on a token paginated query.
	Prev *PaginationTokenLink `json:"prev,omitempty"`

	// Link response on a token paginated query.
	Last *PaginationTokenLink `json:"last,omitempty"`

	// A list of audit records.
	Audits []AuditLogDigest `json:"audits" validate:"required"`
}

// UnmarshalAuditLogs unmarshals an instance of AuditLogs from the specified map of raw messages.
func UnmarshalAuditLogs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AuditLogs)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationTokenLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationTokenLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "prev", &obj.Prev, UnmarshalPaginationTokenLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationTokenLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "audits", &obj.Audits, UnmarshalAuditLogDigest)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AuditLogs) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// Badge : Badge information.
type Badge struct {
	// ID of the current badge.
	ID *string `json:"id,omitempty"`

	// Display name for the current badge.
	Label *string `json:"label,omitempty"`

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

	// Description of the current badge.
	Description *string `json:"description,omitempty"`

	// A map of translated strings, by language code.
	DescriptionI18n map[string]string `json:"description_i18n,omitempty"`

	// Icon for the current badge.
	Icon *string `json:"icon,omitempty"`

	// Authority for the current badge.
	Authority *string `json:"authority,omitempty"`

	// Tag for the current badge.
	Tag *string `json:"tag,omitempty"`

	// Learn more links for a badge.
	LearnMoreLinks *LearnMoreLinks `json:"learn_more_links,omitempty"`

	// An optional set of constraints indicating which versions in an Offering have this particular badge.
	Constraints []Constraint `json:"constraints,omitempty"`
}

// UnmarshalBadge unmarshals an instance of Badge from the specified map of raw messages.
func UnmarshalBadge(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Badge)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label_i18n", &obj.LabelI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description_i18n", &obj.DescriptionI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "icon", &obj.Icon)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "authority", &obj.Authority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tag", &obj.Tag)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "learn_more_links", &obj.LearnMoreLinks, UnmarshalLearnMoreLinks)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "constraints", &obj.Constraints, UnmarshalConstraint)
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

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

	// Description in the requested language.
	ShortDescription *string `json:"short_description,omitempty"`

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

	// URL for an icon associated with this catalog.
	CatalogIconURL *string `json:"catalog_icon_url,omitempty"`

	// URL for a banner image for this catalog.
	CatalogBannerURL *string `json:"catalog_banner_url,omitempty"`

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

	// Catalog specific metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "label_i18n", &obj.LabelI18n)
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
	err = core.UnmarshalPrimitive(m, "catalog_icon_url", &obj.CatalogIconURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_banner_url", &obj.CatalogBannerURL)
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
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
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

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// The programmatic name of this object.
	Name *string `json:"name,omitempty"`

	// The crn for this specific object.
	CRN *string `json:"crn,omitempty"`

	// The url for this specific object.
	URL *string `json:"url,omitempty"`

	// The parent for this specific object.
	ParentID *string `json:"parent_id,omitempty"`

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

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

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

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
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

	// Cluster Ingress hostname.
	IngressHostname *string `json:"ingress_hostname,omitempty"`

	// Cluster provider.
	Provider *string `json:"provider,omitempty"`

	// Cluster status.
	Status *string `json:"status,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "ingress_hostname", &obj.IngressHostname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
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

// CommitVersionOptions : The CommitVersion options.
type CommitVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

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
func (_options *CommitVersionOptions) SetVersionLocID(versionLocID string) *CommitVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CommitVersionOptions) SetHeaders(param map[string]string) *CommitVersionOptions {
	options.Headers = param
	return options
}

// ComplianceControl : Control that can be added to a version.
type ComplianceControl struct {
	// SCC profile.
	SccProfile *ComplianceControlSccProfile `json:"scc_profile,omitempty"`

	// Control family.
	Family *ComplianceControlFamily `json:"family,omitempty"`

	// Control goals.
	Goals []Goal `json:"goals,omitempty"`

	// Control validation.
	Validation *ComplianceControlValidation `json:"validation,omitempty"`
}

// UnmarshalComplianceControl unmarshals an instance of ComplianceControl from the specified map of raw messages.
func UnmarshalComplianceControl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceControl)
	err = core.UnmarshalModel(m, "scc_profile", &obj.SccProfile, UnmarshalComplianceControlSccProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "family", &obj.Family, UnmarshalComplianceControlFamily)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "goals", &obj.Goals, UnmarshalGoal)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validation", &obj.Validation, UnmarshalComplianceControlValidation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComplianceControlFamily : Control family.
type ComplianceControlFamily struct {
	// ID.
	ID *string `json:"id,omitempty"`

	// External ID.
	ExternalID *string `json:"external_id,omitempty"`

	// Description.
	Description *string `json:"description,omitempty"`

	// UI href.
	UIHref *string `json:"ui_href,omitempty"`
}

// UnmarshalComplianceControlFamily unmarshals an instance of ComplianceControlFamily from the specified map of raw messages.
func UnmarshalComplianceControlFamily(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceControlFamily)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "external_id", &obj.ExternalID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComplianceControlSccProfile : SCC profile.
type ComplianceControlSccProfile struct {
	// Profile type.
	Type *string `json:"type,omitempty"`
}

// UnmarshalComplianceControlSccProfile unmarshals an instance of ComplianceControlSccProfile from the specified map of raw messages.
func UnmarshalComplianceControlSccProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceControlSccProfile)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComplianceControlValidation : Control validation.
type ComplianceControlValidation struct {
	// Validation certified bool.
	Certified *bool `json:"certified,omitempty"`

	// Map of validation results.
	Results map[string]interface{} `json:"results,omitempty"`
}

// UnmarshalComplianceControlValidation unmarshals an instance of ComplianceControlValidation from the specified map of raw messages.
func UnmarshalComplianceControlValidation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceControlValidation)
	err = core.UnmarshalPrimitive(m, "certified", &obj.Certified)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "results", &obj.Results)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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

	// Display name for configuration type.
	DisplayName *string `json:"display_name,omitempty"`

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

	// Render type.
	CustomConfig *RenderType `json:"custom_config,omitempty"`

	// The original type, as found in the source being onboarded.
	TypeMetadata *string `json:"type_metadata,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
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
	err = core.UnmarshalModel(m, "custom_config", &obj.CustomConfig, UnmarshalRenderType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type_metadata", &obj.TypeMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Constraint : Constraint information.
type Constraint struct {
	// Type of the current constraint.
	Type *string `json:"type,omitempty"`

	// Rule for the current constraint.
	Rule interface{} `json:"rule,omitempty"`
}

// UnmarshalConstraint unmarshals an instance of Constraint from the specified map of raw messages.
func UnmarshalConstraint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Constraint)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rule", &obj.Rule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConsumableShareObjectOptions : The ConsumableShareObject options.
type ConsumableShareObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewConsumableShareObjectOptions : Instantiate ConsumableShareObjectOptions
func (*CatalogManagementV1) NewConsumableShareObjectOptions(catalogIdentifier string, objectIdentifier string) *ConsumableShareObjectOptions {
	return &ConsumableShareObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *ConsumableShareObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *ConsumableShareObjectOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *ConsumableShareObjectOptions) SetObjectIdentifier(objectIdentifier string) *ConsumableShareObjectOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ConsumableShareObjectOptions) SetHeaders(param map[string]string) *ConsumableShareObjectOptions {
	options.Headers = param
	return options
}

// ConsumableVersionOptions : The ConsumableVersion options.
type ConsumableVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewConsumableVersionOptions : Instantiate ConsumableVersionOptions
func (*CatalogManagementV1) NewConsumableVersionOptions(versionLocID string) *ConsumableVersionOptions {
	return &ConsumableVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (_options *ConsumableVersionOptions) SetVersionLocID(versionLocID string) *ConsumableVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ConsumableVersionOptions) SetHeaders(param map[string]string) *ConsumableVersionOptions {
	options.Headers = param
	return options
}

// CopyFromPreviousVersionOptions : The CopyFromPreviousVersion options.
type CopyFromPreviousVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// The type of data you would like to copy from a previous version. Valid values are 'configuration' or 'licenses'.
	Type *string `json:"type" validate:"required,ne="`

	// The version locator id of the version you wish to copy data from.
	VersionLocIDToCopyFrom *string `json:"version_loc_id_to_copy_from" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCopyFromPreviousVersionOptions : Instantiate CopyFromPreviousVersionOptions
func (*CatalogManagementV1) NewCopyFromPreviousVersionOptions(versionLocID string, typeVar string, versionLocIDToCopyFrom string) *CopyFromPreviousVersionOptions {
	return &CopyFromPreviousVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
		Type: core.StringPtr(typeVar),
		VersionLocIDToCopyFrom: core.StringPtr(versionLocIDToCopyFrom),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (_options *CopyFromPreviousVersionOptions) SetVersionLocID(versionLocID string) *CopyFromPreviousVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CopyFromPreviousVersionOptions) SetType(typeVar string) *CopyFromPreviousVersionOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetVersionLocIDToCopyFrom : Allow user to set VersionLocIDToCopyFrom
func (_options *CopyFromPreviousVersionOptions) SetVersionLocIDToCopyFrom(versionLocIDToCopyFrom string) *CopyFromPreviousVersionOptions {
	_options.VersionLocIDToCopyFrom = core.StringPtr(versionLocIDToCopyFrom)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CopyFromPreviousVersionOptions) SetHeaders(param map[string]string) *CopyFromPreviousVersionOptions {
	options.Headers = param
	return options
}

// CopyVersionOptions : The CopyVersion options.
type CopyVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// Tags array.
	Tags []string `json:"tags,omitempty"`

	// byte array representing the content to be imported.  Only supported for OVA images at this time.
	Content *[]byte `json:"content,omitempty"`

	// Target kinds.  Current valid values are 'iks', 'roks', 'vcenter', 'power-iaas', and 'terraform'.
	TargetKinds []string `json:"target_kinds,omitempty"`

	// Format of content being onboarded. Example: vsi-image. Required for virtual server image for VPC.
	FormatKind *string `json:"format_kind,omitempty"`

	// Version Flavor Information.  Only supported for Product kind Solution.
	Flavor *Flavor `json:"flavor,omitempty"`

	// Optional - The sub-folder within the specified tgz file that contains the software being onboarded.
	WorkingDirectory *string `json:"working_directory,omitempty"`

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
func (_options *CopyVersionOptions) SetVersionLocID(versionLocID string) *CopyVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CopyVersionOptions) SetTags(tags []string) *CopyVersionOptions {
	_options.Tags = tags
	return _options
}

// SetContent : Allow user to set Content
func (_options *CopyVersionOptions) SetContent(content []byte) *CopyVersionOptions {
	_options.Content = &content
	return _options
}

// SetTargetKinds : Allow user to set TargetKinds
func (_options *CopyVersionOptions) SetTargetKinds(targetKinds []string) *CopyVersionOptions {
	_options.TargetKinds = targetKinds
	return _options
}

// SetFormatKind : Allow user to set FormatKind
func (_options *CopyVersionOptions) SetFormatKind(formatKind string) *CopyVersionOptions {
	_options.FormatKind = core.StringPtr(formatKind)
	return _options
}

// SetFlavor : Allow user to set Flavor
func (_options *CopyVersionOptions) SetFlavor(flavor *Flavor) *CopyVersionOptions {
	_options.Flavor = flavor
	return _options
}

// SetWorkingDirectory : Allow user to set WorkingDirectory
func (_options *CopyVersionOptions) SetWorkingDirectory(workingDirectory string) *CopyVersionOptions {
	_options.WorkingDirectory = core.StringPtr(workingDirectory)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CopyVersionOptions) SetHeaders(param map[string]string) *CopyVersionOptions {
	options.Headers = param
	return options
}

// CostBreakdown : Cost breakdown definition.
type CostBreakdown struct {
	// Total hourly cost.
	TotalHourlyCost *string `json:"totalHourlyCost,omitempty"`

	// Total monthly cost.
	TotalMonthlyCost *string `json:"totalMonthlyCost,omitempty"`

	// Resources.
	Resources []CostResource `json:"resources,omitempty"`
}

// UnmarshalCostBreakdown unmarshals an instance of CostBreakdown from the specified map of raw messages.
func UnmarshalCostBreakdown(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CostBreakdown)
	err = core.UnmarshalPrimitive(m, "totalHourlyCost", &obj.TotalHourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalMonthlyCost", &obj.TotalMonthlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalCostResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CostComponent : Cost component definition.
type CostComponent struct {
	// Cost component name.
	Name *string `json:"name,omitempty"`

	// Cost component unit.
	Unit *string `json:"unit,omitempty"`

	// Cost component hourly quantity.
	HourlyQuantity *string `json:"hourlyQuantity,omitempty"`

	// Cost component monthly quantity.
	MonthlyQuantity *string `json:"monthlyQuantity,omitempty"`

	// Cost component price.
	Price *string `json:"price,omitempty"`

	// Cost component hourly cost.
	HourlyCost *string `json:"hourlyCost,omitempty"`

	// Cost component monthly cist.
	MonthlyCost *string `json:"monthlyCost,omitempty"`
}

// UnmarshalCostComponent unmarshals an instance of CostComponent from the specified map of raw messages.
func UnmarshalCostComponent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CostComponent)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unit", &obj.Unit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hourlyQuantity", &obj.HourlyQuantity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "monthlyQuantity", &obj.MonthlyQuantity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "price", &obj.Price)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hourlyCost", &obj.HourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "monthlyCost", &obj.MonthlyCost)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CostEstimate : Cost estimate definition.
type CostEstimate struct {
	// Cost estimate version.
	Version *string `json:"version,omitempty"`

	// Cost estimate currency.
	Currency *string `json:"currency,omitempty"`

	// Cost estimate projects.
	Projects []Project `json:"projects,omitempty"`

	// Cost summary definition.
	Summary *CostSummary `json:"summary,omitempty"`

	// Total hourly cost.
	TotalHourlyCost *string `json:"totalHourlyCost,omitempty"`

	// Total monthly cost.
	TotalMonthlyCost *string `json:"totalMonthlyCost,omitempty"`

	// Past total hourly cost.
	PastTotalHourlyCost *string `json:"pastTotalHourlyCost,omitempty"`

	// Past total monthly cost.
	PastTotalMonthlyCost *string `json:"pastTotalMonthlyCost,omitempty"`

	// Difference in total hourly cost.
	DiffTotalHourlyCost *string `json:"diffTotalHourlyCost,omitempty"`

	// Difference in total monthly cost.
	DiffTotalMonthlyCost *string `json:"diffTotalMonthlyCost,omitempty"`

	// When this estimate was generated.
	TimeGenerated *strfmt.DateTime `json:"timeGenerated,omitempty"`
}

// UnmarshalCostEstimate unmarshals an instance of CostEstimate from the specified map of raw messages.
func UnmarshalCostEstimate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CostEstimate)
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "currency", &obj.Currency)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "projects", &obj.Projects, UnmarshalProject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalCostSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalHourlyCost", &obj.TotalHourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalMonthlyCost", &obj.TotalMonthlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pastTotalHourlyCost", &obj.PastTotalHourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pastTotalMonthlyCost", &obj.PastTotalMonthlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "diffTotalHourlyCost", &obj.DiffTotalHourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "diffTotalMonthlyCost", &obj.DiffTotalMonthlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timeGenerated", &obj.TimeGenerated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CostResource : Cost resource definition.
type CostResource struct {
	// Resource name.
	Name *string `json:"name,omitempty"`

	// Resource metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Hourly cost.
	HourlyCost *string `json:"hourlyCost,omitempty"`

	// Monthly cost.
	MonthlyCost *string `json:"monthlyCost,omitempty"`

	// Cost components.
	CostComponents []CostComponent `json:"costComponents,omitempty"`
}

// UnmarshalCostResource unmarshals an instance of CostResource from the specified map of raw messages.
func UnmarshalCostResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CostResource)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hourlyCost", &obj.HourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "monthlyCost", &obj.MonthlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "costComponents", &obj.CostComponents, UnmarshalCostComponent)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CostSummary : Cost summary definition.
type CostSummary struct {
	// Total detected resources.
	TotalDetectedResources *int64 `json:"totalDetectedResources,omitempty"`

	// Total supported resources.
	TotalSupportedResources *int64 `json:"totalSupportedResources,omitempty"`

	// Total unsupported resources.
	TotalUnsupportedResources *int64 `json:"totalUnsupportedResources,omitempty"`

	// Total usage based resources.
	TotalUsageBasedResources *int64 `json:"totalUsageBasedResources,omitempty"`

	// Total no price resources.
	TotalNoPriceResources *int64 `json:"totalNoPriceResources,omitempty"`

	// Unsupported resource counts.
	UnsupportedResourceCounts map[string]int64 `json:"unsupportedResourceCounts,omitempty"`

	// No price resource counts.
	NoPriceResourceCounts map[string]int64 `json:"noPriceResourceCounts,omitempty"`
}

// UnmarshalCostSummary unmarshals an instance of CostSummary from the specified map of raw messages.
func UnmarshalCostSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CostSummary)
	err = core.UnmarshalPrimitive(m, "totalDetectedResources", &obj.TotalDetectedResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalSupportedResources", &obj.TotalSupportedResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalUnsupportedResources", &obj.TotalUnsupportedResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalUsageBasedResources", &obj.TotalUsageBasedResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalNoPriceResources", &obj.TotalNoPriceResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unsupportedResourceCounts", &obj.UnsupportedResourceCounts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "noPriceResourceCounts", &obj.NoPriceResourceCounts)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCatalogOptions : The CreateCatalog options.
type CreateCatalogOptions struct {
	// Display Name in the requested language.
	Label *string `json:"label,omitempty"`

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

	// Description in the requested language.
	ShortDescription *string `json:"short_description,omitempty"`

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

	// URL for an icon associated with this catalog.
	CatalogIconURL *string `json:"catalog_icon_url,omitempty"`

	// URL for a banner image for this catalog.
	CatalogBannerURL *string `json:"catalog_banner_url,omitempty"`

	// List of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// List of features associated with this catalog.
	Features []Feature `json:"features,omitempty"`

	// Denotes whether a catalog is disabled.
	Disabled *bool `json:"disabled,omitempty"`

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

	// Catalog specific metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateCatalogOptions : Instantiate CreateCatalogOptions
func (*CatalogManagementV1) NewCreateCatalogOptions() *CreateCatalogOptions {
	return &CreateCatalogOptions{}
}

// SetLabel : Allow user to set Label
func (_options *CreateCatalogOptions) SetLabel(label string) *CreateCatalogOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetLabelI18n : Allow user to set LabelI18n
func (_options *CreateCatalogOptions) SetLabelI18n(labelI18n map[string]string) *CreateCatalogOptions {
	_options.LabelI18n = labelI18n
	return _options
}

// SetShortDescription : Allow user to set ShortDescription
func (_options *CreateCatalogOptions) SetShortDescription(shortDescription string) *CreateCatalogOptions {
	_options.ShortDescription = core.StringPtr(shortDescription)
	return _options
}

// SetShortDescriptionI18n : Allow user to set ShortDescriptionI18n
func (_options *CreateCatalogOptions) SetShortDescriptionI18n(shortDescriptionI18n map[string]string) *CreateCatalogOptions {
	_options.ShortDescriptionI18n = shortDescriptionI18n
	return _options
}

// SetCatalogIconURL : Allow user to set CatalogIconURL
func (_options *CreateCatalogOptions) SetCatalogIconURL(catalogIconURL string) *CreateCatalogOptions {
	_options.CatalogIconURL = core.StringPtr(catalogIconURL)
	return _options
}

// SetCatalogBannerURL : Allow user to set CatalogBannerURL
func (_options *CreateCatalogOptions) SetCatalogBannerURL(catalogBannerURL string) *CreateCatalogOptions {
	_options.CatalogBannerURL = core.StringPtr(catalogBannerURL)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateCatalogOptions) SetTags(tags []string) *CreateCatalogOptions {
	_options.Tags = tags
	return _options
}

// SetFeatures : Allow user to set Features
func (_options *CreateCatalogOptions) SetFeatures(features []Feature) *CreateCatalogOptions {
	_options.Features = features
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *CreateCatalogOptions) SetDisabled(disabled bool) *CreateCatalogOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *CreateCatalogOptions) SetResourceGroupID(resourceGroupID string) *CreateCatalogOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetOwningAccount : Allow user to set OwningAccount
func (_options *CreateCatalogOptions) SetOwningAccount(owningAccount string) *CreateCatalogOptions {
	_options.OwningAccount = core.StringPtr(owningAccount)
	return _options
}

// SetCatalogFilters : Allow user to set CatalogFilters
func (_options *CreateCatalogOptions) SetCatalogFilters(catalogFilters *Filters) *CreateCatalogOptions {
	_options.CatalogFilters = catalogFilters
	return _options
}

// SetSyndicationSettings : Allow user to set SyndicationSettings
func (_options *CreateCatalogOptions) SetSyndicationSettings(syndicationSettings *SyndicationResource) *CreateCatalogOptions {
	_options.SyndicationSettings = syndicationSettings
	return _options
}

// SetKind : Allow user to set Kind
func (_options *CreateCatalogOptions) SetKind(kind string) *CreateCatalogOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateCatalogOptions) SetMetadata(metadata map[string]interface{}) *CreateCatalogOptions {
	_options.Metadata = metadata
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCatalogOptions) SetHeaders(param map[string]string) *CreateCatalogOptions {
	options.Headers = param
	return options
}

// CreateObjectAccessOptions : The CreateObjectAccess options.
type CreateObjectAccessOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// Identifier for access. Use 'accountId' for an account, '-ent-enterpriseid' for an enterprise, and
	// '-entgroup-enterprisegroupid' for an enterprise group.
	AccessIdentifier *string `json:"access_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateObjectAccessOptions : Instantiate CreateObjectAccessOptions
func (*CatalogManagementV1) NewCreateObjectAccessOptions(catalogIdentifier string, objectIdentifier string, accessIdentifier string) *CreateObjectAccessOptions {
	return &CreateObjectAccessOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		AccessIdentifier: core.StringPtr(accessIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *CreateObjectAccessOptions) SetCatalogIdentifier(catalogIdentifier string) *CreateObjectAccessOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *CreateObjectAccessOptions) SetObjectIdentifier(objectIdentifier string) *CreateObjectAccessOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetAccessIdentifier : Allow user to set AccessIdentifier
func (_options *CreateObjectAccessOptions) SetAccessIdentifier(accessIdentifier string) *CreateObjectAccessOptions {
	_options.AccessIdentifier = core.StringPtr(accessIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateObjectAccessOptions) SetHeaders(param map[string]string) *CreateObjectAccessOptions {
	options.Headers = param
	return options
}

// CreateObjectOptions : The CreateObject options.
type CreateObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// The programmatic name of this object.
	Name *string `json:"name,omitempty"`

	// The crn for this specific object.
	CRN *string `json:"crn,omitempty"`

	// The url for this specific object.
	URL *string `json:"url,omitempty"`

	// The parent for this specific object.
	ParentID *string `json:"parent_id,omitempty"`

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

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

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

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
func (_options *CreateObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *CreateObjectOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateObjectOptions) SetName(name string) *CreateObjectOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetCRN : Allow user to set CRN
func (_options *CreateObjectOptions) SetCRN(crn string) *CreateObjectOptions {
	_options.CRN = core.StringPtr(crn)
	return _options
}

// SetURL : Allow user to set URL
func (_options *CreateObjectOptions) SetURL(url string) *CreateObjectOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetParentID : Allow user to set ParentID
func (_options *CreateObjectOptions) SetParentID(parentID string) *CreateObjectOptions {
	_options.ParentID = core.StringPtr(parentID)
	return _options
}

// SetLabelI18n : Allow user to set LabelI18n
func (_options *CreateObjectOptions) SetLabelI18n(labelI18n map[string]string) *CreateObjectOptions {
	_options.LabelI18n = labelI18n
	return _options
}

// SetLabel : Allow user to set Label
func (_options *CreateObjectOptions) SetLabel(label string) *CreateObjectOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateObjectOptions) SetTags(tags []string) *CreateObjectOptions {
	_options.Tags = tags
	return _options
}

// SetCreated : Allow user to set Created
func (_options *CreateObjectOptions) SetCreated(created *strfmt.DateTime) *CreateObjectOptions {
	_options.Created = created
	return _options
}

// SetUpdated : Allow user to set Updated
func (_options *CreateObjectOptions) SetUpdated(updated *strfmt.DateTime) *CreateObjectOptions {
	_options.Updated = updated
	return _options
}

// SetShortDescription : Allow user to set ShortDescription
func (_options *CreateObjectOptions) SetShortDescription(shortDescription string) *CreateObjectOptions {
	_options.ShortDescription = core.StringPtr(shortDescription)
	return _options
}

// SetShortDescriptionI18n : Allow user to set ShortDescriptionI18n
func (_options *CreateObjectOptions) SetShortDescriptionI18n(shortDescriptionI18n map[string]string) *CreateObjectOptions {
	_options.ShortDescriptionI18n = shortDescriptionI18n
	return _options
}

// SetKind : Allow user to set Kind
func (_options *CreateObjectOptions) SetKind(kind string) *CreateObjectOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetPublish : Allow user to set Publish
func (_options *CreateObjectOptions) SetPublish(publish *PublishObject) *CreateObjectOptions {
	_options.Publish = publish
	return _options
}

// SetState : Allow user to set State
func (_options *CreateObjectOptions) SetState(state *State) *CreateObjectOptions {
	_options.State = state
	return _options
}

// SetCatalogID : Allow user to set CatalogID
func (_options *CreateObjectOptions) SetCatalogID(catalogID string) *CreateObjectOptions {
	_options.CatalogID = core.StringPtr(catalogID)
	return _options
}

// SetCatalogName : Allow user to set CatalogName
func (_options *CreateObjectOptions) SetCatalogName(catalogName string) *CreateObjectOptions {
	_options.CatalogName = core.StringPtr(catalogName)
	return _options
}

// SetData : Allow user to set Data
func (_options *CreateObjectOptions) SetData(data map[string]interface{}) *CreateObjectOptions {
	_options.Data = data
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateObjectOptions) SetHeaders(param map[string]string) *CreateObjectOptions {
	options.Headers = param
	return options
}

// CreateOfferingInstanceOptions : The CreateOfferingInstance options.
type CreateOfferingInstanceOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

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

	// The version this instance was installed from (semver - not version id).
	Version *string `json:"version,omitempty"`

	// The version id this instance was installed from (version id - not semver).
	VersionID *string `json:"version_id,omitempty"`

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

	// Type of install plan (also known as approval strategy) for operator subscriptions. Can be either automatic, which
	// automatically upgrades operators to the latest in a channel, or manual, which requires approval on the cluster.
	InstallPlan *string `json:"install_plan,omitempty"`

	// Channel to pin the operator subscription to.
	Channel *string `json:"channel,omitempty"`

	// date and time create.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// date and time updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Map of metadata values for this offering instance.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Id of the resource group to provision the offering instance into.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// String location of OfferingInstance deployment.
	Location *string `json:"location,omitempty"`

	// Indicates if Resource Controller has disabled this instance.
	Disabled *bool `json:"disabled,omitempty"`

	// The account this instance is owned by.
	Account *string `json:"account,omitempty"`

	// the last operation performed and status.
	LastOperation *OfferingInstanceLastOperation `json:"last_operation,omitempty"`

	// The target kind for the installed software version.
	KindTarget *string `json:"kind_target,omitempty"`

	// The digest value of the installed software version.
	Sha *string `json:"sha,omitempty"`

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
func (_options *CreateOfferingInstanceOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *CreateOfferingInstanceOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateOfferingInstanceOptions) SetID(id string) *CreateOfferingInstanceOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *CreateOfferingInstanceOptions) SetRev(rev string) *CreateOfferingInstanceOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetURL : Allow user to set URL
func (_options *CreateOfferingInstanceOptions) SetURL(url string) *CreateOfferingInstanceOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetCRN : Allow user to set CRN
func (_options *CreateOfferingInstanceOptions) SetCRN(crn string) *CreateOfferingInstanceOptions {
	_options.CRN = core.StringPtr(crn)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *CreateOfferingInstanceOptions) SetLabel(label string) *CreateOfferingInstanceOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetCatalogID : Allow user to set CatalogID
func (_options *CreateOfferingInstanceOptions) SetCatalogID(catalogID string) *CreateOfferingInstanceOptions {
	_options.CatalogID = core.StringPtr(catalogID)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *CreateOfferingInstanceOptions) SetOfferingID(offeringID string) *CreateOfferingInstanceOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetKindFormat : Allow user to set KindFormat
func (_options *CreateOfferingInstanceOptions) SetKindFormat(kindFormat string) *CreateOfferingInstanceOptions {
	_options.KindFormat = core.StringPtr(kindFormat)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *CreateOfferingInstanceOptions) SetVersion(version string) *CreateOfferingInstanceOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *CreateOfferingInstanceOptions) SetVersionID(versionID string) *CreateOfferingInstanceOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *CreateOfferingInstanceOptions) SetClusterID(clusterID string) *CreateOfferingInstanceOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetClusterRegion : Allow user to set ClusterRegion
func (_options *CreateOfferingInstanceOptions) SetClusterRegion(clusterRegion string) *CreateOfferingInstanceOptions {
	_options.ClusterRegion = core.StringPtr(clusterRegion)
	return _options
}

// SetClusterNamespaces : Allow user to set ClusterNamespaces
func (_options *CreateOfferingInstanceOptions) SetClusterNamespaces(clusterNamespaces []string) *CreateOfferingInstanceOptions {
	_options.ClusterNamespaces = clusterNamespaces
	return _options
}

// SetClusterAllNamespaces : Allow user to set ClusterAllNamespaces
func (_options *CreateOfferingInstanceOptions) SetClusterAllNamespaces(clusterAllNamespaces bool) *CreateOfferingInstanceOptions {
	_options.ClusterAllNamespaces = core.BoolPtr(clusterAllNamespaces)
	return _options
}

// SetSchematicsWorkspaceID : Allow user to set SchematicsWorkspaceID
func (_options *CreateOfferingInstanceOptions) SetSchematicsWorkspaceID(schematicsWorkspaceID string) *CreateOfferingInstanceOptions {
	_options.SchematicsWorkspaceID = core.StringPtr(schematicsWorkspaceID)
	return _options
}

// SetInstallPlan : Allow user to set InstallPlan
func (_options *CreateOfferingInstanceOptions) SetInstallPlan(installPlan string) *CreateOfferingInstanceOptions {
	_options.InstallPlan = core.StringPtr(installPlan)
	return _options
}

// SetChannel : Allow user to set Channel
func (_options *CreateOfferingInstanceOptions) SetChannel(channel string) *CreateOfferingInstanceOptions {
	_options.Channel = core.StringPtr(channel)
	return _options
}

// SetCreated : Allow user to set Created
func (_options *CreateOfferingInstanceOptions) SetCreated(created *strfmt.DateTime) *CreateOfferingInstanceOptions {
	_options.Created = created
	return _options
}

// SetUpdated : Allow user to set Updated
func (_options *CreateOfferingInstanceOptions) SetUpdated(updated *strfmt.DateTime) *CreateOfferingInstanceOptions {
	_options.Updated = updated
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateOfferingInstanceOptions) SetMetadata(metadata map[string]interface{}) *CreateOfferingInstanceOptions {
	_options.Metadata = metadata
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *CreateOfferingInstanceOptions) SetResourceGroupID(resourceGroupID string) *CreateOfferingInstanceOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *CreateOfferingInstanceOptions) SetLocation(location string) *CreateOfferingInstanceOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *CreateOfferingInstanceOptions) SetDisabled(disabled bool) *CreateOfferingInstanceOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *CreateOfferingInstanceOptions) SetAccount(account string) *CreateOfferingInstanceOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetLastOperation : Allow user to set LastOperation
func (_options *CreateOfferingInstanceOptions) SetLastOperation(lastOperation *OfferingInstanceLastOperation) *CreateOfferingInstanceOptions {
	_options.LastOperation = lastOperation
	return _options
}

// SetKindTarget : Allow user to set KindTarget
func (_options *CreateOfferingInstanceOptions) SetKindTarget(kindTarget string) *CreateOfferingInstanceOptions {
	_options.KindTarget = core.StringPtr(kindTarget)
	return _options
}

// SetSha : Allow user to set Sha
func (_options *CreateOfferingInstanceOptions) SetSha(sha string) *CreateOfferingInstanceOptions {
	_options.Sha = core.StringPtr(sha)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateOfferingInstanceOptions) SetHeaders(param map[string]string) *CreateOfferingInstanceOptions {
	options.Headers = param
	return options
}

// CreateOfferingOptions : The CreateOffering options.
type CreateOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// The url for this specific offering.
	URL *string `json:"url,omitempty"`

	// The crn for this specific offering.
	CRN *string `json:"crn,omitempty"`

	// Display Name in the requested language.
	Label *string `json:"label,omitempty"`

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

	// The programmatic name of this offering.
	Name *string `json:"name,omitempty"`

	// URL for an icon associated with this offering.
	OfferingIconURL *string `json:"offering_icon_url,omitempty"`

	// URL for an additional docs with this offering.
	OfferingDocsURL *string `json:"offering_docs_url,omitempty"`

	// [deprecated] - Use offering.support instead.  URL to be displayed in the Consumption UI for getting support on this
	// offering.
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

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

	// Long description in the requested language.
	LongDescription *string `json:"long_description,omitempty"`

	// A map of translated strings, by language code.
	LongDescriptionI18n map[string]string `json:"long_description_i18n,omitempty"`

	// list of features associated with this offering.
	Features []Feature `json:"features,omitempty"`

	// Array of kind.
	Kinds []Kind `json:"kinds,omitempty"`

	// Offering is managed by Partner Center.
	PcManaged *bool `json:"pc_managed,omitempty"`

	// Offering has been approved to publish to permitted to IBM or Public Catalog.
	PublishApproved *bool `json:"publish_approved,omitempty"`

	// Denotes public availability of an Offering - if share_enabled is true.
	ShareWithAll *bool `json:"share_with_all,omitempty"`

	// Denotes IBM employee availability of an Offering - if share_enabled is true.
	ShareWithIBM *bool `json:"share_with_ibm,omitempty"`

	// Denotes sharing including access list availability of an Offering is enabled.
	ShareEnabled *bool `json:"share_enabled,omitempty"`

	// Is it permitted to request publishing to IBM or Public.
	// Deprecated: this field is deprecated and may be removed in a future release.
	PermitRequestIBMPublicPublish *bool `json:"permit_request_ibm_public_publish,omitempty"`

	// Indicates if this offering has been approved for use by all IBMers.
	// Deprecated: this field is deprecated and may be removed in a future release.
	IBMPublishApproved *bool `json:"ibm_publish_approved,omitempty"`

	// Indicates if this offering has been approved for use by all IBM Cloud users.
	// Deprecated: this field is deprecated and may be removed in a future release.
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

	// Deprecated - Provider of this offering.
	// Deprecated: this field is deprecated and may be removed in a future release.
	Provider *string `json:"provider,omitempty"`

	// Information on the provider for this offering, or omitted if no provider information is given.
	ProviderInfo *ProviderInfo `json:"provider_info,omitempty"`

	// Repository info for offerings.
	RepoInfo *RepoInfo `json:"repo_info,omitempty"`

	// Image pull keys for this offering.
	ImagePullKeys []ImagePullKey `json:"image_pull_keys,omitempty"`

	// Offering Support information.
	Support *Support `json:"support,omitempty"`

	// A list of media items related to this offering.
	Media []MediaItem `json:"media,omitempty"`

	// Deprecation information for an Offering.
	DeprecatePending *DeprecatePending `json:"deprecate_pending,omitempty"`

	// The product kind.  Valid values are module, solution, or empty string.
	ProductKind *string `json:"product_kind,omitempty"`

	// A list of badges for this offering.
	Badges []Badge `json:"badges,omitempty"`

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
func (_options *CreateOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *CreateOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetURL : Allow user to set URL
func (_options *CreateOfferingOptions) SetURL(url string) *CreateOfferingOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetCRN : Allow user to set CRN
func (_options *CreateOfferingOptions) SetCRN(crn string) *CreateOfferingOptions {
	_options.CRN = core.StringPtr(crn)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *CreateOfferingOptions) SetLabel(label string) *CreateOfferingOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetLabelI18n : Allow user to set LabelI18n
func (_options *CreateOfferingOptions) SetLabelI18n(labelI18n map[string]string) *CreateOfferingOptions {
	_options.LabelI18n = labelI18n
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateOfferingOptions) SetName(name string) *CreateOfferingOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetOfferingIconURL : Allow user to set OfferingIconURL
func (_options *CreateOfferingOptions) SetOfferingIconURL(offeringIconURL string) *CreateOfferingOptions {
	_options.OfferingIconURL = core.StringPtr(offeringIconURL)
	return _options
}

// SetOfferingDocsURL : Allow user to set OfferingDocsURL
func (_options *CreateOfferingOptions) SetOfferingDocsURL(offeringDocsURL string) *CreateOfferingOptions {
	_options.OfferingDocsURL = core.StringPtr(offeringDocsURL)
	return _options
}

// SetOfferingSupportURL : Allow user to set OfferingSupportURL
func (_options *CreateOfferingOptions) SetOfferingSupportURL(offeringSupportURL string) *CreateOfferingOptions {
	_options.OfferingSupportURL = core.StringPtr(offeringSupportURL)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateOfferingOptions) SetTags(tags []string) *CreateOfferingOptions {
	_options.Tags = tags
	return _options
}

// SetKeywords : Allow user to set Keywords
func (_options *CreateOfferingOptions) SetKeywords(keywords []string) *CreateOfferingOptions {
	_options.Keywords = keywords
	return _options
}

// SetRating : Allow user to set Rating
func (_options *CreateOfferingOptions) SetRating(rating *Rating) *CreateOfferingOptions {
	_options.Rating = rating
	return _options
}

// SetCreated : Allow user to set Created
func (_options *CreateOfferingOptions) SetCreated(created *strfmt.DateTime) *CreateOfferingOptions {
	_options.Created = created
	return _options
}

// SetUpdated : Allow user to set Updated
func (_options *CreateOfferingOptions) SetUpdated(updated *strfmt.DateTime) *CreateOfferingOptions {
	_options.Updated = updated
	return _options
}

// SetShortDescription : Allow user to set ShortDescription
func (_options *CreateOfferingOptions) SetShortDescription(shortDescription string) *CreateOfferingOptions {
	_options.ShortDescription = core.StringPtr(shortDescription)
	return _options
}

// SetShortDescriptionI18n : Allow user to set ShortDescriptionI18n
func (_options *CreateOfferingOptions) SetShortDescriptionI18n(shortDescriptionI18n map[string]string) *CreateOfferingOptions {
	_options.ShortDescriptionI18n = shortDescriptionI18n
	return _options
}

// SetLongDescription : Allow user to set LongDescription
func (_options *CreateOfferingOptions) SetLongDescription(longDescription string) *CreateOfferingOptions {
	_options.LongDescription = core.StringPtr(longDescription)
	return _options
}

// SetLongDescriptionI18n : Allow user to set LongDescriptionI18n
func (_options *CreateOfferingOptions) SetLongDescriptionI18n(longDescriptionI18n map[string]string) *CreateOfferingOptions {
	_options.LongDescriptionI18n = longDescriptionI18n
	return _options
}

// SetFeatures : Allow user to set Features
func (_options *CreateOfferingOptions) SetFeatures(features []Feature) *CreateOfferingOptions {
	_options.Features = features
	return _options
}

// SetKinds : Allow user to set Kinds
func (_options *CreateOfferingOptions) SetKinds(kinds []Kind) *CreateOfferingOptions {
	_options.Kinds = kinds
	return _options
}

// SetPcManaged : Allow user to set PcManaged
func (_options *CreateOfferingOptions) SetPcManaged(pcManaged bool) *CreateOfferingOptions {
	_options.PcManaged = core.BoolPtr(pcManaged)
	return _options
}

// SetPublishApproved : Allow user to set PublishApproved
func (_options *CreateOfferingOptions) SetPublishApproved(publishApproved bool) *CreateOfferingOptions {
	_options.PublishApproved = core.BoolPtr(publishApproved)
	return _options
}

// SetShareWithAll : Allow user to set ShareWithAll
func (_options *CreateOfferingOptions) SetShareWithAll(shareWithAll bool) *CreateOfferingOptions {
	_options.ShareWithAll = core.BoolPtr(shareWithAll)
	return _options
}

// SetShareWithIBM : Allow user to set ShareWithIBM
func (_options *CreateOfferingOptions) SetShareWithIBM(shareWithIBM bool) *CreateOfferingOptions {
	_options.ShareWithIBM = core.BoolPtr(shareWithIBM)
	return _options
}

// SetShareEnabled : Allow user to set ShareEnabled
func (_options *CreateOfferingOptions) SetShareEnabled(shareEnabled bool) *CreateOfferingOptions {
	_options.ShareEnabled = core.BoolPtr(shareEnabled)
	return _options
}

// SetPermitRequestIBMPublicPublish : Allow user to set PermitRequestIBMPublicPublish
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *CreateOfferingOptions) SetPermitRequestIBMPublicPublish(permitRequestIBMPublicPublish bool) *CreateOfferingOptions {
	_options.PermitRequestIBMPublicPublish = core.BoolPtr(permitRequestIBMPublicPublish)
	return _options
}

// SetIBMPublishApproved : Allow user to set IBMPublishApproved
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *CreateOfferingOptions) SetIBMPublishApproved(ibmPublishApproved bool) *CreateOfferingOptions {
	_options.IBMPublishApproved = core.BoolPtr(ibmPublishApproved)
	return _options
}

// SetPublicPublishApproved : Allow user to set PublicPublishApproved
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *CreateOfferingOptions) SetPublicPublishApproved(publicPublishApproved bool) *CreateOfferingOptions {
	_options.PublicPublishApproved = core.BoolPtr(publicPublishApproved)
	return _options
}

// SetPublicOriginalCRN : Allow user to set PublicOriginalCRN
func (_options *CreateOfferingOptions) SetPublicOriginalCRN(publicOriginalCRN string) *CreateOfferingOptions {
	_options.PublicOriginalCRN = core.StringPtr(publicOriginalCRN)
	return _options
}

// SetPublishPublicCRN : Allow user to set PublishPublicCRN
func (_options *CreateOfferingOptions) SetPublishPublicCRN(publishPublicCRN string) *CreateOfferingOptions {
	_options.PublishPublicCRN = core.StringPtr(publishPublicCRN)
	return _options
}

// SetPortalApprovalRecord : Allow user to set PortalApprovalRecord
func (_options *CreateOfferingOptions) SetPortalApprovalRecord(portalApprovalRecord string) *CreateOfferingOptions {
	_options.PortalApprovalRecord = core.StringPtr(portalApprovalRecord)
	return _options
}

// SetPortalUIURL : Allow user to set PortalUIURL
func (_options *CreateOfferingOptions) SetPortalUIURL(portalUIURL string) *CreateOfferingOptions {
	_options.PortalUIURL = core.StringPtr(portalUIURL)
	return _options
}

// SetCatalogID : Allow user to set CatalogID
func (_options *CreateOfferingOptions) SetCatalogID(catalogID string) *CreateOfferingOptions {
	_options.CatalogID = core.StringPtr(catalogID)
	return _options
}

// SetCatalogName : Allow user to set CatalogName
func (_options *CreateOfferingOptions) SetCatalogName(catalogName string) *CreateOfferingOptions {
	_options.CatalogName = core.StringPtr(catalogName)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateOfferingOptions) SetMetadata(metadata map[string]interface{}) *CreateOfferingOptions {
	_options.Metadata = metadata
	return _options
}

// SetDisclaimer : Allow user to set Disclaimer
func (_options *CreateOfferingOptions) SetDisclaimer(disclaimer string) *CreateOfferingOptions {
	_options.Disclaimer = core.StringPtr(disclaimer)
	return _options
}

// SetHidden : Allow user to set Hidden
func (_options *CreateOfferingOptions) SetHidden(hidden bool) *CreateOfferingOptions {
	_options.Hidden = core.BoolPtr(hidden)
	return _options
}

// SetProvider : Allow user to set Provider
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *CreateOfferingOptions) SetProvider(provider string) *CreateOfferingOptions {
	_options.Provider = core.StringPtr(provider)
	return _options
}

// SetProviderInfo : Allow user to set ProviderInfo
func (_options *CreateOfferingOptions) SetProviderInfo(providerInfo *ProviderInfo) *CreateOfferingOptions {
	_options.ProviderInfo = providerInfo
	return _options
}

// SetRepoInfo : Allow user to set RepoInfo
func (_options *CreateOfferingOptions) SetRepoInfo(repoInfo *RepoInfo) *CreateOfferingOptions {
	_options.RepoInfo = repoInfo
	return _options
}

// SetImagePullKeys : Allow user to set ImagePullKeys
func (_options *CreateOfferingOptions) SetImagePullKeys(imagePullKeys []ImagePullKey) *CreateOfferingOptions {
	_options.ImagePullKeys = imagePullKeys
	return _options
}

// SetSupport : Allow user to set Support
func (_options *CreateOfferingOptions) SetSupport(support *Support) *CreateOfferingOptions {
	_options.Support = support
	return _options
}

// SetMedia : Allow user to set Media
func (_options *CreateOfferingOptions) SetMedia(media []MediaItem) *CreateOfferingOptions {
	_options.Media = media
	return _options
}

// SetDeprecatePending : Allow user to set DeprecatePending
func (_options *CreateOfferingOptions) SetDeprecatePending(deprecatePending *DeprecatePending) *CreateOfferingOptions {
	_options.DeprecatePending = deprecatePending
	return _options
}

// SetProductKind : Allow user to set ProductKind
func (_options *CreateOfferingOptions) SetProductKind(productKind string) *CreateOfferingOptions {
	_options.ProductKind = core.StringPtr(productKind)
	return _options
}

// SetBadges : Allow user to set Badges
func (_options *CreateOfferingOptions) SetBadges(badges []Badge) *CreateOfferingOptions {
	_options.Badges = badges
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateOfferingOptions) SetHeaders(param map[string]string) *CreateOfferingOptions {
	options.Headers = param
	return options
}

// DeleteCatalogOptions : The DeleteCatalog options.
type DeleteCatalogOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

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
func (_options *DeleteCatalogOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteCatalogOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCatalogOptions) SetHeaders(param map[string]string) *DeleteCatalogOptions {
	options.Headers = param
	return options
}

// DeleteObjectAccessListOptions : The DeleteObjectAccessList options.
type DeleteObjectAccessListOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// A list of accesses to delete.  An entry with star["*"] will remove all accesses.
	Accesses []string `json:"accesses" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteObjectAccessListOptions : Instantiate DeleteObjectAccessListOptions
func (*CatalogManagementV1) NewDeleteObjectAccessListOptions(catalogIdentifier string, objectIdentifier string, accesses []string) *DeleteObjectAccessListOptions {
	return &DeleteObjectAccessListOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		Accesses: accesses,
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *DeleteObjectAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteObjectAccessListOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *DeleteObjectAccessListOptions) SetObjectIdentifier(objectIdentifier string) *DeleteObjectAccessListOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetAccesses : Allow user to set Accesses
func (_options *DeleteObjectAccessListOptions) SetAccesses(accesses []string) *DeleteObjectAccessListOptions {
	_options.Accesses = accesses
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteObjectAccessListOptions) SetHeaders(param map[string]string) *DeleteObjectAccessListOptions {
	options.Headers = param
	return options
}

// DeleteObjectAccessOptions : The DeleteObjectAccess options.
type DeleteObjectAccessOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// Identifier for access. Use 'accountId' for an account, '-ent-enterpriseid' for an enterprise, and
	// '-entgroup-enterprisegroupid' for an enterprise group.
	AccessIdentifier *string `json:"access_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteObjectAccessOptions : Instantiate DeleteObjectAccessOptions
func (*CatalogManagementV1) NewDeleteObjectAccessOptions(catalogIdentifier string, objectIdentifier string, accessIdentifier string) *DeleteObjectAccessOptions {
	return &DeleteObjectAccessOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		AccessIdentifier: core.StringPtr(accessIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *DeleteObjectAccessOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteObjectAccessOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *DeleteObjectAccessOptions) SetObjectIdentifier(objectIdentifier string) *DeleteObjectAccessOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetAccessIdentifier : Allow user to set AccessIdentifier
func (_options *DeleteObjectAccessOptions) SetAccessIdentifier(accessIdentifier string) *DeleteObjectAccessOptions {
	_options.AccessIdentifier = core.StringPtr(accessIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteObjectAccessOptions) SetHeaders(param map[string]string) *DeleteObjectAccessOptions {
	options.Headers = param
	return options
}

// DeleteObjectOptions : The DeleteObject options.
type DeleteObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

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
func (_options *DeleteObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteObjectOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *DeleteObjectOptions) SetObjectIdentifier(objectIdentifier string) *DeleteObjectOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteObjectOptions) SetHeaders(param map[string]string) *DeleteObjectOptions {
	options.Headers = param
	return options
}

// DeleteOfferingAccessListOptions : The DeleteOfferingAccessList options.
type DeleteOfferingAccessListOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// A list of accesses to delete.  An entry with star["*"] will remove all accesses.
	Accesses []string `json:"accesses" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteOfferingAccessListOptions : Instantiate DeleteOfferingAccessListOptions
func (*CatalogManagementV1) NewDeleteOfferingAccessListOptions(catalogIdentifier string, offeringID string, accesses []string) *DeleteOfferingAccessListOptions {
	return &DeleteOfferingAccessListOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		Accesses: accesses,
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *DeleteOfferingAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteOfferingAccessListOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *DeleteOfferingAccessListOptions) SetOfferingID(offeringID string) *DeleteOfferingAccessListOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetAccesses : Allow user to set Accesses
func (_options *DeleteOfferingAccessListOptions) SetAccesses(accesses []string) *DeleteOfferingAccessListOptions {
	_options.Accesses = accesses
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOfferingAccessListOptions) SetHeaders(param map[string]string) *DeleteOfferingAccessListOptions {
	options.Headers = param
	return options
}

// DeleteOfferingInstanceOptions : The DeleteOfferingInstance options.
type DeleteOfferingInstanceOptions struct {
	// Version Instance identifier.
	InstanceIdentifier *string `json:"instance_identifier" validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

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
func (_options *DeleteOfferingInstanceOptions) SetInstanceIdentifier(instanceIdentifier string) *DeleteOfferingInstanceOptions {
	_options.InstanceIdentifier = core.StringPtr(instanceIdentifier)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *DeleteOfferingInstanceOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *DeleteOfferingInstanceOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOfferingInstanceOptions) SetHeaders(param map[string]string) *DeleteOfferingInstanceOptions {
	options.Headers = param
	return options
}

// DeleteOfferingOptions : The DeleteOffering options.
type DeleteOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

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
func (_options *DeleteOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *DeleteOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *DeleteOfferingOptions) SetOfferingID(offeringID string) *DeleteOfferingOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOfferingOptions) SetHeaders(param map[string]string) *DeleteOfferingOptions {
	options.Headers = param
	return options
}

// DeleteOperatorsOptions : The DeleteOperators options.
type DeleteOperatorsOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// Cluster identification.
	ClusterID *string `json:"cluster_id" validate:"required"`

	// Cluster region.
	Region *string `json:"region" validate:"required"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `json:"version_locator_id" validate:"required"`

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
func (_options *DeleteOperatorsOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *DeleteOperatorsOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *DeleteOperatorsOptions) SetClusterID(clusterID string) *DeleteOperatorsOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *DeleteOperatorsOptions) SetRegion(region string) *DeleteOperatorsOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (_options *DeleteOperatorsOptions) SetVersionLocatorID(versionLocatorID string) *DeleteOperatorsOptions {
	_options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOperatorsOptions) SetHeaders(param map[string]string) *DeleteOperatorsOptions {
	options.Headers = param
	return options
}

// DeleteVersionOptions : The DeleteVersion options.
type DeleteVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

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
func (_options *DeleteVersionOptions) SetVersionLocID(versionLocID string) *DeleteVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteVersionOptions) SetHeaders(param map[string]string) *DeleteVersionOptions {
	options.Headers = param
	return options
}

// DeployOperatorsOptions : The DeployOperators options.
type DeployOperatorsOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// Cluster ID.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster region.
	Region *string `json:"region,omitempty"`

	// Kube namespaces to deploy Operator(s) to.
	Namespaces []string `json:"namespaces,omitempty"`

	// Denotes whether to install Operator(s) globally.
	AllNamespaces *bool `json:"all_namespaces,omitempty"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `json:"version_locator_id,omitempty"`

	// Operator channel.
	Channel *string `json:"channel,omitempty"`

	// Plan.
	InstallPlan *string `json:"install_plan,omitempty"`

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
func (_options *DeployOperatorsOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *DeployOperatorsOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *DeployOperatorsOptions) SetClusterID(clusterID string) *DeployOperatorsOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *DeployOperatorsOptions) SetRegion(region string) *DeployOperatorsOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetNamespaces : Allow user to set Namespaces
func (_options *DeployOperatorsOptions) SetNamespaces(namespaces []string) *DeployOperatorsOptions {
	_options.Namespaces = namespaces
	return _options
}

// SetAllNamespaces : Allow user to set AllNamespaces
func (_options *DeployOperatorsOptions) SetAllNamespaces(allNamespaces bool) *DeployOperatorsOptions {
	_options.AllNamespaces = core.BoolPtr(allNamespaces)
	return _options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (_options *DeployOperatorsOptions) SetVersionLocatorID(versionLocatorID string) *DeployOperatorsOptions {
	_options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return _options
}

// SetChannel : Allow user to set Channel
func (_options *DeployOperatorsOptions) SetChannel(channel string) *DeployOperatorsOptions {
	_options.Channel = core.StringPtr(channel)
	return _options
}

// SetInstallPlan : Allow user to set InstallPlan
func (_options *DeployOperatorsOptions) SetInstallPlan(installPlan string) *DeployOperatorsOptions {
	_options.InstallPlan = core.StringPtr(installPlan)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeployOperatorsOptions) SetHeaders(param map[string]string) *DeployOperatorsOptions {
	options.Headers = param
	return options
}

// DeployRequestBodyEnvironmentVariablesItem : DeployRequestBodyEnvironmentVariablesItem struct
type DeployRequestBodyEnvironmentVariablesItem struct {
	// Variable name.
	Name *string `json:"name,omitempty"`

	// Variable value.
	Value interface{} `json:"value,omitempty"`

	// Does this variable contain a secure value.
	Secure *bool `json:"secure,omitempty"`
}

// UnmarshalDeployRequestBodyEnvironmentVariablesItem unmarshals an instance of DeployRequestBodyEnvironmentVariablesItem from the specified map of raw messages.
func UnmarshalDeployRequestBodyEnvironmentVariablesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeployRequestBodyEnvironmentVariablesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secure", &obj.Secure)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeployRequestBodyOverrideValues : Validation override values. Required for virtual server image for VPC.
type DeployRequestBodyOverrideValues struct {
	// Name of virtual server image instance to create. Required for virtual server image for VPC.
	VsiInstanceName *string `json:"vsi_instance_name,omitempty"`

	// Profile to use when validating virtual server image. Required for virtual server image for VPC.
	VPCProfile *string `json:"vpc_profile,omitempty"`

	// ID of subnet to use when validating virtual server image. Required for virtual server image for VPC.
	SubnetID *string `json:"subnet_id,omitempty"`

	// ID of VPC to use when validating virtual server image. Required for virtual server image for VPC.
	VPCID *string `json:"vpc_id,omitempty"`

	// Zone of subnet to use when validating virtual server image. Required for virtual server image for VPC.
	SubnetZone *string `json:"subnet_zone,omitempty"`

	// ID off SSH key to use when validating virtual server image. Required for virtual server image for VPC.
	SSHKeyID *string `json:"ssh_key_id,omitempty"`

	// Region virtual server image exists in. Required for virtual server image for VPC.
	VPCRegion *string `json:"vpc_region,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of DeployRequestBodyOverrideValues
func (o *DeployRequestBodyOverrideValues) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of DeployRequestBodyOverrideValues
func (o *DeployRequestBodyOverrideValues) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of DeployRequestBodyOverrideValues
func (o *DeployRequestBodyOverrideValues) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of DeployRequestBodyOverrideValues
func (o *DeployRequestBodyOverrideValues) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of DeployRequestBodyOverrideValues
func (o *DeployRequestBodyOverrideValues) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.VsiInstanceName != nil {
		m["vsi_instance_name"] = o.VsiInstanceName
	}
	if o.VPCProfile != nil {
		m["vpc_profile"] = o.VPCProfile
	}
	if o.SubnetID != nil {
		m["subnet_id"] = o.SubnetID
	}
	if o.VPCID != nil {
		m["vpc_id"] = o.VPCID
	}
	if o.SubnetZone != nil {
		m["subnet_zone"] = o.SubnetZone
	}
	if o.SSHKeyID != nil {
		m["ssh_key_id"] = o.SSHKeyID
	}
	if o.VPCRegion != nil {
		m["vpc_region"] = o.VPCRegion
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalDeployRequestBodyOverrideValues unmarshals an instance of DeployRequestBodyOverrideValues from the specified map of raw messages.
func UnmarshalDeployRequestBodyOverrideValues(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeployRequestBodyOverrideValues)
	err = core.UnmarshalPrimitive(m, "vsi_instance_name", &obj.VsiInstanceName)
	if err != nil {
		return
	}
	delete(m, "vsi_instance_name")
	err = core.UnmarshalPrimitive(m, "vpc_profile", &obj.VPCProfile)
	if err != nil {
		return
	}
	delete(m, "vpc_profile")
	err = core.UnmarshalPrimitive(m, "subnet_id", &obj.SubnetID)
	if err != nil {
		return
	}
	delete(m, "subnet_id")
	err = core.UnmarshalPrimitive(m, "vpc_id", &obj.VPCID)
	if err != nil {
		return
	}
	delete(m, "vpc_id")
	err = core.UnmarshalPrimitive(m, "subnet_zone", &obj.SubnetZone)
	if err != nil {
		return
	}
	delete(m, "subnet_zone")
	err = core.UnmarshalPrimitive(m, "ssh_key_id", &obj.SSHKeyID)
	if err != nil {
		return
	}
	delete(m, "ssh_key_id")
	err = core.UnmarshalPrimitive(m, "vpc_region", &obj.VPCRegion)
	if err != nil {
		return
	}
	delete(m, "vpc_region")
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

	// Terraform version override.
	TerraformVersion *string `json:"terraform_version,omitempty"`

	// Schematics workspace region.
	Region *string `json:"region,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "terraform_version", &obj.TerraformVersion)
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

// DeprecateOfferingOptions : The DeprecateOffering options.
type DeprecateOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Set deprecation (true) or cancel deprecation (false).
	Setting *string `json:"setting" validate:"required,ne="`

	// Additional information that users can provide to be displayed in deprecation notification.
	Description *string `json:"description,omitempty"`

	// Specifies the amount of days until product is not available in catalog.
	DaysUntilDeprecate *int64 `json:"days_until_deprecate,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeprecateOfferingOptions.Setting property.
// Set deprecation (true) or cancel deprecation (false).
const (
	DeprecateOfferingOptionsSettingFalseConst = "false"
	DeprecateOfferingOptionsSettingTrueConst = "true"
)

// NewDeprecateOfferingOptions : Instantiate DeprecateOfferingOptions
func (*CatalogManagementV1) NewDeprecateOfferingOptions(catalogIdentifier string, offeringID string, setting string) *DeprecateOfferingOptions {
	return &DeprecateOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		Setting: core.StringPtr(setting),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *DeprecateOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *DeprecateOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *DeprecateOfferingOptions) SetOfferingID(offeringID string) *DeprecateOfferingOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetSetting : Allow user to set Setting
func (_options *DeprecateOfferingOptions) SetSetting(setting string) *DeprecateOfferingOptions {
	_options.Setting = core.StringPtr(setting)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *DeprecateOfferingOptions) SetDescription(description string) *DeprecateOfferingOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetDaysUntilDeprecate : Allow user to set DaysUntilDeprecate
func (_options *DeprecateOfferingOptions) SetDaysUntilDeprecate(daysUntilDeprecate int64) *DeprecateOfferingOptions {
	_options.DaysUntilDeprecate = core.Int64Ptr(daysUntilDeprecate)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeprecateOfferingOptions) SetHeaders(param map[string]string) *DeprecateOfferingOptions {
	options.Headers = param
	return options
}

// DeprecatePending : Deprecation information for an Offering.
type DeprecatePending struct {
	// Date of deprecation.
	DeprecateDate *strfmt.DateTime `json:"deprecate_date,omitempty"`

	// Deprecation state.
	DeprecateState *string `json:"deprecate_state,omitempty"`

	Description *string `json:"description,omitempty"`
}

// UnmarshalDeprecatePending unmarshals an instance of DeprecatePending from the specified map of raw messages.
func UnmarshalDeprecatePending(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeprecatePending)
	err = core.UnmarshalPrimitive(m, "deprecate_date", &obj.DeprecateDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "deprecate_state", &obj.DeprecateState)
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

// DeprecateVersionOptions : The DeprecateVersion options.
type DeprecateVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

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
func (_options *DeprecateVersionOptions) SetVersionLocID(versionLocID string) *DeprecateVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
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

	// A map of translated strings, by language code.
	TitleI18n map[string]string `json:"title_i18n,omitempty"`

	// Feature description.
	Description *string `json:"description,omitempty"`

	// A map of translated strings, by language code.
	DescriptionI18n map[string]string `json:"description_i18n,omitempty"`
}

// UnmarshalFeature unmarshals an instance of Feature from the specified map of raw messages.
func UnmarshalFeature(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Feature)
	err = core.UnmarshalPrimitive(m, "title", &obj.Title)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "title_i18n", &obj.TitleI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description_i18n", &obj.DescriptionI18n)
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

// Flavor : Version Flavor Information.  Only supported for Product kind Solution.
type Flavor struct {
	// Programmatic name for this flavor.
	Name *string `json:"name,omitempty"`

	// Label for this flavor.
	Label *string `json:"label,omitempty"`

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

	// Order that this flavor should appear when listed for a single version.
	Index *int64 `json:"index,omitempty"`
}

// UnmarshalFlavor unmarshals an instance of Flavor from the specified map of raw messages.
func UnmarshalFlavor(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Flavor)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "label_i18n", &obj.LabelI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "index", &obj.Index)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetCatalogAccountAuditOptions : The GetCatalogAccountAudit options.
type GetCatalogAccountAuditOptions struct {
	// Auditlog ID.
	AuditlogIdentifier *string `json:"auditlog_identifier" validate:"required,ne="`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCatalogAccountAuditOptions : Instantiate GetCatalogAccountAuditOptions
func (*CatalogManagementV1) NewGetCatalogAccountAuditOptions(auditlogIdentifier string) *GetCatalogAccountAuditOptions {
	return &GetCatalogAccountAuditOptions{
		AuditlogIdentifier: core.StringPtr(auditlogIdentifier),
	}
}

// SetAuditlogIdentifier : Allow user to set AuditlogIdentifier
func (_options *GetCatalogAccountAuditOptions) SetAuditlogIdentifier(auditlogIdentifier string) *GetCatalogAccountAuditOptions {
	_options.AuditlogIdentifier = core.StringPtr(auditlogIdentifier)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *GetCatalogAccountAuditOptions) SetLookupnames(lookupnames bool) *GetCatalogAccountAuditOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogAccountAuditOptions) SetHeaders(param map[string]string) *GetCatalogAccountAuditOptions {
	options.Headers = param
	return options
}

// GetCatalogAccountFiltersOptions : The GetCatalogAccountFilters options.
type GetCatalogAccountFiltersOptions struct {
	// catalog id. Narrow down filters to the account and just the one catalog.
	Catalog *string `json:"catalog,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCatalogAccountFiltersOptions : Instantiate GetCatalogAccountFiltersOptions
func (*CatalogManagementV1) NewGetCatalogAccountFiltersOptions() *GetCatalogAccountFiltersOptions {
	return &GetCatalogAccountFiltersOptions{}
}

// SetCatalog : Allow user to set Catalog
func (_options *GetCatalogAccountFiltersOptions) SetCatalog(catalog string) *GetCatalogAccountFiltersOptions {
	_options.Catalog = core.StringPtr(catalog)
	return _options
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
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Auditlog ID.
	AuditlogIdentifier *string `json:"auditlog_identifier" validate:"required,ne="`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCatalogAuditOptions : Instantiate GetCatalogAuditOptions
func (*CatalogManagementV1) NewGetCatalogAuditOptions(catalogIdentifier string, auditlogIdentifier string) *GetCatalogAuditOptions {
	return &GetCatalogAuditOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		AuditlogIdentifier: core.StringPtr(auditlogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *GetCatalogAuditOptions) SetCatalogIdentifier(catalogIdentifier string) *GetCatalogAuditOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetAuditlogIdentifier : Allow user to set AuditlogIdentifier
func (_options *GetCatalogAuditOptions) SetAuditlogIdentifier(auditlogIdentifier string) *GetCatalogAuditOptions {
	_options.AuditlogIdentifier = core.StringPtr(auditlogIdentifier)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *GetCatalogAuditOptions) SetLookupnames(lookupnames bool) *GetCatalogAuditOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogAuditOptions) SetHeaders(param map[string]string) *GetCatalogAuditOptions {
	options.Headers = param
	return options
}

// GetCatalogOptions : The GetCatalog options.
type GetCatalogOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

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
func (_options *GetCatalogOptions) SetCatalogIdentifier(catalogIdentifier string) *GetCatalogOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogOptions) SetHeaders(param map[string]string) *GetCatalogOptions {
	options.Headers = param
	return options
}

// GetClusterOptions : The GetCluster options.
type GetClusterOptions struct {
	// ID of the cluster.
	ClusterID *string `json:"cluster_id" validate:"required,ne="`

	// Region of the cluster.
	Region *string `json:"region" validate:"required"`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

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
func (_options *GetClusterOptions) SetClusterID(clusterID string) *GetClusterOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *GetClusterOptions) SetRegion(region string) *GetClusterOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *GetClusterOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetClusterOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
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
	Digest *bool `json:"digest,omitempty"`

	// catalog id. Narrow search down to just a particular catalog. It will apply the catalog's public filters to the
	// public catalog offerings on the result.
	Catalog *string `json:"catalog,omitempty"`

	// What should be selected. Default is 'all' which will return both public and private offerings. 'public' returns only
	// the public offerings and 'private' returns only the private offerings.
	Select *string `json:"select,omitempty"`

	// true - include offerings which have been marked as hidden. The default is false and hidden offerings are not
	// returned.
	IncludeHidden *bool `json:"includeHidden,omitempty"`

	// number or results to return.
	Limit *int64 `json:"limit,omitempty"`

	// number of results to skip before returning values.
	Offset *int64 `json:"offset,omitempty"`

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
func (_options *GetConsumptionOfferingsOptions) SetDigest(digest bool) *GetConsumptionOfferingsOptions {
	_options.Digest = core.BoolPtr(digest)
	return _options
}

// SetCatalog : Allow user to set Catalog
func (_options *GetConsumptionOfferingsOptions) SetCatalog(catalog string) *GetConsumptionOfferingsOptions {
	_options.Catalog = core.StringPtr(catalog)
	return _options
}

// SetSelect : Allow user to set Select
func (_options *GetConsumptionOfferingsOptions) SetSelect(selectVar string) *GetConsumptionOfferingsOptions {
	_options.Select = core.StringPtr(selectVar)
	return _options
}

// SetIncludeHidden : Allow user to set IncludeHidden
func (_options *GetConsumptionOfferingsOptions) SetIncludeHidden(includeHidden bool) *GetConsumptionOfferingsOptions {
	_options.IncludeHidden = core.BoolPtr(includeHidden)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetConsumptionOfferingsOptions) SetLimit(limit int64) *GetConsumptionOfferingsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetConsumptionOfferingsOptions) SetOffset(offset int64) *GetConsumptionOfferingsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConsumptionOfferingsOptions) SetHeaders(param map[string]string) *GetConsumptionOfferingsOptions {
	options.Headers = param
	return options
}

// GetEnterpriseAuditOptions : The GetEnterpriseAudit options.
type GetEnterpriseAuditOptions struct {
	// Enterprise ID.
	EnterpriseIdentifier *string `json:"enterprise_identifier" validate:"required,ne="`

	// Auditlog ID.
	AuditlogIdentifier *string `json:"auditlog_identifier" validate:"required,ne="`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEnterpriseAuditOptions : Instantiate GetEnterpriseAuditOptions
func (*CatalogManagementV1) NewGetEnterpriseAuditOptions(enterpriseIdentifier string, auditlogIdentifier string) *GetEnterpriseAuditOptions {
	return &GetEnterpriseAuditOptions{
		EnterpriseIdentifier: core.StringPtr(enterpriseIdentifier),
		AuditlogIdentifier: core.StringPtr(auditlogIdentifier),
	}
}

// SetEnterpriseIdentifier : Allow user to set EnterpriseIdentifier
func (_options *GetEnterpriseAuditOptions) SetEnterpriseIdentifier(enterpriseIdentifier string) *GetEnterpriseAuditOptions {
	_options.EnterpriseIdentifier = core.StringPtr(enterpriseIdentifier)
	return _options
}

// SetAuditlogIdentifier : Allow user to set AuditlogIdentifier
func (_options *GetEnterpriseAuditOptions) SetAuditlogIdentifier(auditlogIdentifier string) *GetEnterpriseAuditOptions {
	_options.AuditlogIdentifier = core.StringPtr(auditlogIdentifier)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *GetEnterpriseAuditOptions) SetLookupnames(lookupnames bool) *GetEnterpriseAuditOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetEnterpriseAuditOptions) SetHeaders(param map[string]string) *GetEnterpriseAuditOptions {
	options.Headers = param
	return options
}

// GetNamespacesOptions : The GetNamespaces options.
type GetNamespacesOptions struct {
	// ID of the cluster.
	ClusterID *string `json:"cluster_id" validate:"required,ne="`

	// Cluster region.
	Region *string `json:"region" validate:"required"`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// The maximum number of results to return.
	Limit *int64 `json:"limit,omitempty"`

	// The number of results to skip before returning values.
	Offset *int64 `json:"offset,omitempty"`

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
func (_options *GetNamespacesOptions) SetClusterID(clusterID string) *GetNamespacesOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *GetNamespacesOptions) SetRegion(region string) *GetNamespacesOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *GetNamespacesOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetNamespacesOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetNamespacesOptions) SetLimit(limit int64) *GetNamespacesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetNamespacesOptions) SetOffset(offset int64) *GetNamespacesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetNamespacesOptions) SetHeaders(param map[string]string) *GetNamespacesOptions {
	options.Headers = param
	return options
}

// GetObjectAccessListDeprecatedOptions : The GetObjectAccessListDeprecated options.
type GetObjectAccessListDeprecatedOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// The maximum number of results to return.
	Limit *int64 `json:"limit,omitempty"`

	// The number of results to skip before returning values.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetObjectAccessListDeprecatedOptions : Instantiate GetObjectAccessListDeprecatedOptions
func (*CatalogManagementV1) NewGetObjectAccessListDeprecatedOptions(catalogIdentifier string, objectIdentifier string) *GetObjectAccessListDeprecatedOptions {
	return &GetObjectAccessListDeprecatedOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *GetObjectAccessListDeprecatedOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectAccessListDeprecatedOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *GetObjectAccessListDeprecatedOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectAccessListDeprecatedOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetObjectAccessListDeprecatedOptions) SetLimit(limit int64) *GetObjectAccessListDeprecatedOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *GetObjectAccessListDeprecatedOptions) SetOffset(offset int64) *GetObjectAccessListDeprecatedOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectAccessListDeprecatedOptions) SetHeaders(param map[string]string) *GetObjectAccessListDeprecatedOptions {
	options.Headers = param
	return options
}

// GetObjectAccessListOptions : The GetObjectAccessList options.
type GetObjectAccessListOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// Start token for a query.
	Start *string `json:"start,omitempty"`

	// number or results to return in the query.
	Limit *int64 `json:"limit,omitempty"`

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
func (_options *GetObjectAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectAccessListOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *GetObjectAccessListOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectAccessListOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetStart : Allow user to set Start
func (_options *GetObjectAccessListOptions) SetStart(start string) *GetObjectAccessListOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetObjectAccessListOptions) SetLimit(limit int64) *GetObjectAccessListOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectAccessListOptions) SetHeaders(param map[string]string) *GetObjectAccessListOptions {
	options.Headers = param
	return options
}

// GetObjectAccessOptions : The GetObjectAccess options.
type GetObjectAccessOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// Identifier for access. Use 'accountId' for an account, '-ent-enterpriseid' for an enterprise, and
	// '-entgroup-enterprisegroupid' for an enterprise group.
	AccessIdentifier *string `json:"access_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetObjectAccessOptions : Instantiate GetObjectAccessOptions
func (*CatalogManagementV1) NewGetObjectAccessOptions(catalogIdentifier string, objectIdentifier string, accessIdentifier string) *GetObjectAccessOptions {
	return &GetObjectAccessOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		AccessIdentifier: core.StringPtr(accessIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *GetObjectAccessOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectAccessOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *GetObjectAccessOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectAccessOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetAccessIdentifier : Allow user to set AccessIdentifier
func (_options *GetObjectAccessOptions) SetAccessIdentifier(accessIdentifier string) *GetObjectAccessOptions {
	_options.AccessIdentifier = core.StringPtr(accessIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectAccessOptions) SetHeaders(param map[string]string) *GetObjectAccessOptions {
	options.Headers = param
	return options
}

// GetObjectAuditOptions : The GetObjectAudit options.
type GetObjectAuditOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// Auditlog ID.
	AuditlogIdentifier *string `json:"auditlog_identifier" validate:"required,ne="`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetObjectAuditOptions : Instantiate GetObjectAuditOptions
func (*CatalogManagementV1) NewGetObjectAuditOptions(catalogIdentifier string, objectIdentifier string, auditlogIdentifier string) *GetObjectAuditOptions {
	return &GetObjectAuditOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
		AuditlogIdentifier: core.StringPtr(auditlogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *GetObjectAuditOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectAuditOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *GetObjectAuditOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectAuditOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetAuditlogIdentifier : Allow user to set AuditlogIdentifier
func (_options *GetObjectAuditOptions) SetAuditlogIdentifier(auditlogIdentifier string) *GetObjectAuditOptions {
	_options.AuditlogIdentifier = core.StringPtr(auditlogIdentifier)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *GetObjectAuditOptions) SetLookupnames(lookupnames bool) *GetObjectAuditOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectAuditOptions) SetHeaders(param map[string]string) *GetObjectAuditOptions {
	options.Headers = param
	return options
}

// GetObjectOptions : The GetObject options.
type GetObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

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
func (_options *GetObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *GetObjectOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *GetObjectOptions) SetObjectIdentifier(objectIdentifier string) *GetObjectOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetObjectOptions) SetHeaders(param map[string]string) *GetObjectOptions {
	options.Headers = param
	return options
}

// GetOfferingAboutOptions : The GetOfferingAbout options.
type GetOfferingAboutOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

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
func (_options *GetOfferingAboutOptions) SetVersionLocID(versionLocID string) *GetOfferingAboutOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingAboutOptions) SetHeaders(param map[string]string) *GetOfferingAboutOptions {
	options.Headers = param
	return options
}

// GetOfferingAccessListOptions : The GetOfferingAccessList options.
type GetOfferingAccessListOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Start token for a query.
	Start *string `json:"start,omitempty"`

	// number or results to return in the query.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingAccessListOptions : Instantiate GetOfferingAccessListOptions
func (*CatalogManagementV1) NewGetOfferingAccessListOptions(catalogIdentifier string, offeringID string) *GetOfferingAccessListOptions {
	return &GetOfferingAccessListOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *GetOfferingAccessListOptions) SetCatalogIdentifier(catalogIdentifier string) *GetOfferingAccessListOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *GetOfferingAccessListOptions) SetOfferingID(offeringID string) *GetOfferingAccessListOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetStart : Allow user to set Start
func (_options *GetOfferingAccessListOptions) SetStart(start string) *GetOfferingAccessListOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetOfferingAccessListOptions) SetLimit(limit int64) *GetOfferingAccessListOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingAccessListOptions) SetHeaders(param map[string]string) *GetOfferingAccessListOptions {
	options.Headers = param
	return options
}

// GetOfferingAccessOptions : The GetOfferingAccess options.
type GetOfferingAccessOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Identifier for access. Use 'accountId' for an account, '-ent-enterpriseid' for an enterprise, and
	// '-entgroup-enterprisegroupid' for an enterprise group.
	AccessIdentifier *string `json:"access_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingAccessOptions : Instantiate GetOfferingAccessOptions
func (*CatalogManagementV1) NewGetOfferingAccessOptions(catalogIdentifier string, offeringID string, accessIdentifier string) *GetOfferingAccessOptions {
	return &GetOfferingAccessOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		AccessIdentifier: core.StringPtr(accessIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *GetOfferingAccessOptions) SetCatalogIdentifier(catalogIdentifier string) *GetOfferingAccessOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *GetOfferingAccessOptions) SetOfferingID(offeringID string) *GetOfferingAccessOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetAccessIdentifier : Allow user to set AccessIdentifier
func (_options *GetOfferingAccessOptions) SetAccessIdentifier(accessIdentifier string) *GetOfferingAccessOptions {
	_options.AccessIdentifier = core.StringPtr(accessIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingAccessOptions) SetHeaders(param map[string]string) *GetOfferingAccessOptions {
	options.Headers = param
	return options
}

// GetOfferingAuditOptions : The GetOfferingAudit options.
type GetOfferingAuditOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Auditlog ID.
	AuditlogIdentifier *string `json:"auditlog_identifier" validate:"required,ne="`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingAuditOptions : Instantiate GetOfferingAuditOptions
func (*CatalogManagementV1) NewGetOfferingAuditOptions(catalogIdentifier string, offeringID string, auditlogIdentifier string) *GetOfferingAuditOptions {
	return &GetOfferingAuditOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		AuditlogIdentifier: core.StringPtr(auditlogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *GetOfferingAuditOptions) SetCatalogIdentifier(catalogIdentifier string) *GetOfferingAuditOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *GetOfferingAuditOptions) SetOfferingID(offeringID string) *GetOfferingAuditOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetAuditlogIdentifier : Allow user to set AuditlogIdentifier
func (_options *GetOfferingAuditOptions) SetAuditlogIdentifier(auditlogIdentifier string) *GetOfferingAuditOptions {
	_options.AuditlogIdentifier = core.StringPtr(auditlogIdentifier)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *GetOfferingAuditOptions) SetLookupnames(lookupnames bool) *GetOfferingAuditOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingAuditOptions) SetHeaders(param map[string]string) *GetOfferingAuditOptions {
	options.Headers = param
	return options
}

// GetOfferingContainerImagesOptions : The GetOfferingContainerImages options.
type GetOfferingContainerImagesOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

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
func (_options *GetOfferingContainerImagesOptions) SetVersionLocID(versionLocID string) *GetOfferingContainerImagesOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingContainerImagesOptions) SetHeaders(param map[string]string) *GetOfferingContainerImagesOptions {
	options.Headers = param
	return options
}

// GetOfferingInstanceAuditOptions : The GetOfferingInstanceAudit options.
type GetOfferingInstanceAuditOptions struct {
	// Version Instance identifier.
	InstanceIdentifier *string `json:"instance_identifier" validate:"required,ne="`

	// Auditlog ID.
	AuditlogIdentifier *string `json:"auditlog_identifier" validate:"required,ne="`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingInstanceAuditOptions : Instantiate GetOfferingInstanceAuditOptions
func (*CatalogManagementV1) NewGetOfferingInstanceAuditOptions(instanceIdentifier string, auditlogIdentifier string) *GetOfferingInstanceAuditOptions {
	return &GetOfferingInstanceAuditOptions{
		InstanceIdentifier: core.StringPtr(instanceIdentifier),
		AuditlogIdentifier: core.StringPtr(auditlogIdentifier),
	}
}

// SetInstanceIdentifier : Allow user to set InstanceIdentifier
func (_options *GetOfferingInstanceAuditOptions) SetInstanceIdentifier(instanceIdentifier string) *GetOfferingInstanceAuditOptions {
	_options.InstanceIdentifier = core.StringPtr(instanceIdentifier)
	return _options
}

// SetAuditlogIdentifier : Allow user to set AuditlogIdentifier
func (_options *GetOfferingInstanceAuditOptions) SetAuditlogIdentifier(auditlogIdentifier string) *GetOfferingInstanceAuditOptions {
	_options.AuditlogIdentifier = core.StringPtr(auditlogIdentifier)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *GetOfferingInstanceAuditOptions) SetLookupnames(lookupnames bool) *GetOfferingInstanceAuditOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingInstanceAuditOptions) SetHeaders(param map[string]string) *GetOfferingInstanceAuditOptions {
	options.Headers = param
	return options
}

// GetOfferingInstanceOptions : The GetOfferingInstance options.
type GetOfferingInstanceOptions struct {
	// Version Instance identifier.
	InstanceIdentifier *string `json:"instance_identifier" validate:"required,ne="`

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
func (_options *GetOfferingInstanceOptions) SetInstanceIdentifier(instanceIdentifier string) *GetOfferingInstanceOptions {
	_options.InstanceIdentifier = core.StringPtr(instanceIdentifier)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingInstanceOptions) SetHeaders(param map[string]string) *GetOfferingInstanceOptions {
	options.Headers = param
	return options
}

// GetOfferingLicenseOptions : The GetOfferingLicense options.
type GetOfferingLicenseOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// The ID of the license, which maps to the file name in the 'licenses' directory of this verions tgz file.
	LicenseID *string `json:"license_id" validate:"required,ne="`

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
func (_options *GetOfferingLicenseOptions) SetVersionLocID(versionLocID string) *GetOfferingLicenseOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetLicenseID : Allow user to set LicenseID
func (_options *GetOfferingLicenseOptions) SetLicenseID(licenseID string) *GetOfferingLicenseOptions {
	_options.LicenseID = core.StringPtr(licenseID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingLicenseOptions) SetHeaders(param map[string]string) *GetOfferingLicenseOptions {
	options.Headers = param
	return options
}

// GetOfferingOptions : The GetOffering options.
type GetOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Offering Parameter Type.  Valid values are 'name' or 'id'.  Default is 'id'.
	Type *string `json:"type,omitempty"`

	// Return the digest format of the specified offering.  Default is false.
	Digest *bool `json:"digest,omitempty"`

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
func (_options *GetOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *GetOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *GetOfferingOptions) SetOfferingID(offeringID string) *GetOfferingOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetType : Allow user to set Type
func (_options *GetOfferingOptions) SetType(typeVar string) *GetOfferingOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetDigest : Allow user to set Digest
func (_options *GetOfferingOptions) SetDigest(digest bool) *GetOfferingOptions {
	_options.Digest = core.BoolPtr(digest)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingOptions) SetHeaders(param map[string]string) *GetOfferingOptions {
	options.Headers = param
	return options
}

// GetOfferingSourceOptions : The GetOfferingSource options.
type GetOfferingSourceOptions struct {
	// The version being requested.
	Version *string `json:"version" validate:"required"`

	// The type of the response: application/yaml, application/json, or application/x-gzip.
	Accept *string `json:"Accept,omitempty"`

	// Catalog ID.  If not specified, this value will default to the public catalog.
	CatalogID *string `json:"catalogID,omitempty"`

	// Offering name.  An offering name or ID must be specified.
	Name *string `json:"name,omitempty"`

	// Offering id.  An offering name or ID must be specified.
	ID *string `json:"id,omitempty"`

	// The kind of offering (e.g. helm, ova, terraform...).
	Kind *string `json:"kind,omitempty"`

	// The channel value of the specified version.
	Channel *string `json:"channel,omitempty"`

	// The programmatic flavor name of the specified version.
	Flavor *string `json:"flavor,omitempty"`

	// If false (the default), the root folder from the original onboarded tgz file is removed.  If true, the root folder
	// is returned.
	AsIs *bool `json:"asIs,omitempty"`

	// The install type of the specified version.
	InstallType *string `json:"installType,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingSourceOptions : Instantiate GetOfferingSourceOptions
func (*CatalogManagementV1) NewGetOfferingSourceOptions(version string) *GetOfferingSourceOptions {
	return &GetOfferingSourceOptions{
		Version: core.StringPtr(version),
	}
}

// SetVersion : Allow user to set Version
func (_options *GetOfferingSourceOptions) SetVersion(version string) *GetOfferingSourceOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *GetOfferingSourceOptions) SetAccept(accept string) *GetOfferingSourceOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetCatalogID : Allow user to set CatalogID
func (_options *GetOfferingSourceOptions) SetCatalogID(catalogID string) *GetOfferingSourceOptions {
	_options.CatalogID = core.StringPtr(catalogID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetOfferingSourceOptions) SetName(name string) *GetOfferingSourceOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetOfferingSourceOptions) SetID(id string) *GetOfferingSourceOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *GetOfferingSourceOptions) SetKind(kind string) *GetOfferingSourceOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetChannel : Allow user to set Channel
func (_options *GetOfferingSourceOptions) SetChannel(channel string) *GetOfferingSourceOptions {
	_options.Channel = core.StringPtr(channel)
	return _options
}

// SetFlavor : Allow user to set Flavor
func (_options *GetOfferingSourceOptions) SetFlavor(flavor string) *GetOfferingSourceOptions {
	_options.Flavor = core.StringPtr(flavor)
	return _options
}

// SetAsIs : Allow user to set AsIs
func (_options *GetOfferingSourceOptions) SetAsIs(asIs bool) *GetOfferingSourceOptions {
	_options.AsIs = core.BoolPtr(asIs)
	return _options
}

// SetInstallType : Allow user to set InstallType
func (_options *GetOfferingSourceOptions) SetInstallType(installType string) *GetOfferingSourceOptions {
	_options.InstallType = core.StringPtr(installType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingSourceOptions) SetHeaders(param map[string]string) *GetOfferingSourceOptions {
	options.Headers = param
	return options
}

// GetOfferingSourceURLOptions : The GetOfferingSourceURL options.
type GetOfferingSourceURLOptions struct {
	// Unique key identifying an image.
	KeyIdentifier *string `json:"key_identifier" validate:"required,ne="`

	// The type of the response: application/yaml, application/json, or application/x-gzip.
	Accept *string `json:"Accept,omitempty"`

	// Catalog ID. If not specified, this value will default to the public catalog.
	CatalogID *string `json:"catalogID,omitempty"`

	// Offering name. An offering name or ID must be specified.
	Name *string `json:"name,omitempty"`

	// Offering id. An offering name or ID must be specified.
	ID *string `json:"id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingSourceURLOptions : Instantiate GetOfferingSourceURLOptions
func (*CatalogManagementV1) NewGetOfferingSourceURLOptions(keyIdentifier string) *GetOfferingSourceURLOptions {
	return &GetOfferingSourceURLOptions{
		KeyIdentifier: core.StringPtr(keyIdentifier),
	}
}

// SetKeyIdentifier : Allow user to set KeyIdentifier
func (_options *GetOfferingSourceURLOptions) SetKeyIdentifier(keyIdentifier string) *GetOfferingSourceURLOptions {
	_options.KeyIdentifier = core.StringPtr(keyIdentifier)
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *GetOfferingSourceURLOptions) SetAccept(accept string) *GetOfferingSourceURLOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetCatalogID : Allow user to set CatalogID
func (_options *GetOfferingSourceURLOptions) SetCatalogID(catalogID string) *GetOfferingSourceURLOptions {
	_options.CatalogID = core.StringPtr(catalogID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetOfferingSourceURLOptions) SetName(name string) *GetOfferingSourceURLOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetOfferingSourceURLOptions) SetID(id string) *GetOfferingSourceURLOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingSourceURLOptions) SetHeaders(param map[string]string) *GetOfferingSourceURLOptions {
	options.Headers = param
	return options
}

// GetOfferingUpdatesOptions : The GetOfferingUpdates options.
type GetOfferingUpdatesOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// The kind of offering (e.g, helm, ova, terraform ...).
	Kind *string `json:"kind" validate:"required"`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// The target kind of the currently installed version (e.g. iks, roks, etc).
	Target *string `json:"target,omitempty"`

	// optionaly provide an existing version to check updates for if one is not given, all version will be returned.
	Version *string `json:"version,omitempty"`

	// The id of the cluster where this version was installed.
	ClusterID *string `json:"cluster_id,omitempty"`

	// The region of the cluster where this version was installed.
	Region *string `json:"region,omitempty"`

	// The resource group id of the cluster where this version was installed.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// The namespace of the cluster where this version was installed.
	Namespace *string `json:"namespace,omitempty"`

	// The sha value of the currently installed version.
	Sha *string `json:"sha,omitempty"`

	// Optionally provide the channel value of the currently installed version.
	Channel *string `json:"channel,omitempty"`

	// Optionally provide a list of namespaces used for the currently installed version.
	Namespaces []string `json:"namespaces,omitempty"`

	// Optionally indicate that the current version was installed in all namespaces.
	AllNamespaces *bool `json:"all_namespaces,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOfferingUpdatesOptions : Instantiate GetOfferingUpdatesOptions
func (*CatalogManagementV1) NewGetOfferingUpdatesOptions(catalogIdentifier string, offeringID string, kind string, xAuthRefreshToken string) *GetOfferingUpdatesOptions {
	return &GetOfferingUpdatesOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		Kind: core.StringPtr(kind),
		XAuthRefreshToken: core.StringPtr(xAuthRefreshToken),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *GetOfferingUpdatesOptions) SetCatalogIdentifier(catalogIdentifier string) *GetOfferingUpdatesOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *GetOfferingUpdatesOptions) SetOfferingID(offeringID string) *GetOfferingUpdatesOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *GetOfferingUpdatesOptions) SetKind(kind string) *GetOfferingUpdatesOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *GetOfferingUpdatesOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetOfferingUpdatesOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *GetOfferingUpdatesOptions) SetTarget(target string) *GetOfferingUpdatesOptions {
	_options.Target = core.StringPtr(target)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *GetOfferingUpdatesOptions) SetVersion(version string) *GetOfferingUpdatesOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *GetOfferingUpdatesOptions) SetClusterID(clusterID string) *GetOfferingUpdatesOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *GetOfferingUpdatesOptions) SetRegion(region string) *GetOfferingUpdatesOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *GetOfferingUpdatesOptions) SetResourceGroupID(resourceGroupID string) *GetOfferingUpdatesOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetNamespace : Allow user to set Namespace
func (_options *GetOfferingUpdatesOptions) SetNamespace(namespace string) *GetOfferingUpdatesOptions {
	_options.Namespace = core.StringPtr(namespace)
	return _options
}

// SetSha : Allow user to set Sha
func (_options *GetOfferingUpdatesOptions) SetSha(sha string) *GetOfferingUpdatesOptions {
	_options.Sha = core.StringPtr(sha)
	return _options
}

// SetChannel : Allow user to set Channel
func (_options *GetOfferingUpdatesOptions) SetChannel(channel string) *GetOfferingUpdatesOptions {
	_options.Channel = core.StringPtr(channel)
	return _options
}

// SetNamespaces : Allow user to set Namespaces
func (_options *GetOfferingUpdatesOptions) SetNamespaces(namespaces []string) *GetOfferingUpdatesOptions {
	_options.Namespaces = namespaces
	return _options
}

// SetAllNamespaces : Allow user to set AllNamespaces
func (_options *GetOfferingUpdatesOptions) SetAllNamespaces(allNamespaces bool) *GetOfferingUpdatesOptions {
	_options.AllNamespaces = core.BoolPtr(allNamespaces)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingUpdatesOptions) SetHeaders(param map[string]string) *GetOfferingUpdatesOptions {
	options.Headers = param
	return options
}

// GetOfferingWorkingCopyOptions : The GetOfferingWorkingCopy options.
type GetOfferingWorkingCopyOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

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
func (_options *GetOfferingWorkingCopyOptions) SetVersionLocID(versionLocID string) *GetOfferingWorkingCopyOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOfferingWorkingCopyOptions) SetHeaders(param map[string]string) *GetOfferingWorkingCopyOptions {
	options.Headers = param
	return options
}

// GetOverrideValuesOptions : The GetOverrideValues options.
type GetOverrideValuesOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

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
func (_options *GetOverrideValuesOptions) SetVersionLocID(versionLocID string) *GetOverrideValuesOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOverrideValuesOptions) SetHeaders(param map[string]string) *GetOverrideValuesOptions {
	options.Headers = param
	return options
}

// GetPreinstallOptions : The GetPreinstall options.
type GetPreinstallOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// ID of the cluster.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster region.
	Region *string `json:"region,omitempty"`

	// Required if the version's pre-install scope is `namespace`.
	Namespace *string `json:"namespace,omitempty"`

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
func (_options *GetPreinstallOptions) SetVersionLocID(versionLocID string) *GetPreinstallOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *GetPreinstallOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetPreinstallOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *GetPreinstallOptions) SetClusterID(clusterID string) *GetPreinstallOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *GetPreinstallOptions) SetRegion(region string) *GetPreinstallOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetNamespace : Allow user to set Namespace
func (_options *GetPreinstallOptions) SetNamespace(namespace string) *GetPreinstallOptions {
	_options.Namespace = core.StringPtr(namespace)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPreinstallOptions) SetHeaders(param map[string]string) *GetPreinstallOptions {
	options.Headers = param
	return options
}

// GetValidationStatusOptions : The GetValidationStatus options.
type GetValidationStatusOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

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
func (_options *GetValidationStatusOptions) SetVersionLocID(versionLocID string) *GetValidationStatusOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *GetValidationStatusOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *GetValidationStatusOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetValidationStatusOptions) SetHeaders(param map[string]string) *GetValidationStatusOptions {
	options.Headers = param
	return options
}

// GetVersionOptions : The GetVersion options.
type GetVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

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
func (_options *GetVersionOptions) SetVersionLocID(versionLocID string) *GetVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetVersionOptions) SetHeaders(param map[string]string) *GetVersionOptions {
	options.Headers = param
	return options
}

// Goal : Compliance control goal.
type Goal struct {
	// Goal ID.
	ID *string `json:"id,omitempty"`

	// Goal description.
	Description *string `json:"description,omitempty"`

	// Goal UI href.
	UIHref *string `json:"ui_href,omitempty"`
}

// UnmarshalGoal unmarshals an instance of Goal from the specified map of raw messages.
func UnmarshalGoal(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Goal)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IamPermission : IAM Permission definition.
type IamPermission struct {
	// Service name.
	ServiceName *string `json:"service_name,omitempty"`

	// Role CRNs for this permission.
	RoleCrns []string `json:"role_crns,omitempty"`

	// Resources for this permission.
	Resources []IamResource `json:"resources,omitempty"`
}

// UnmarshalIamPermission unmarshals an instance of IamPermission from the specified map of raw messages.
func UnmarshalIamPermission(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamPermission)
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "role_crns", &obj.RoleCrns)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalIamResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IamResource : IAM Resource definition.
type IamResource struct {
	// Resource name.
	Name *string `json:"name,omitempty"`

	// Resource description.
	Description *string `json:"description,omitempty"`

	// Role CRNs for this permission.
	RoleCrns []string `json:"role_crns,omitempty"`
}

// UnmarshalIamResource unmarshals an instance of IamResource from the specified map of raw messages.
func UnmarshalIamResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamResource)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "role_crns", &obj.RoleCrns)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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

// ImagePullKey : Image pull keys for an offering.
type ImagePullKey struct {
	// Key name.
	Name *string `json:"name,omitempty"`

	// Key value.
	Value *string `json:"value,omitempty"`

	// Key description.
	Description *string `json:"description,omitempty"`
}

// UnmarshalImagePullKey unmarshals an instance of ImagePullKey from the specified map of raw messages.
func UnmarshalImagePullKey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImagePullKey)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
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

// ImportOfferingBodyMetadata : Generic data to be included with content being onboarded. Required for virtual server image for VPC.
type ImportOfferingBodyMetadata struct {
	// Operating system included in this image. Required for virtual server image for VPC.
	OperatingSystem *ImportOfferingBodyMetadataOperatingSystem `json:"operating_system,omitempty"`

	// Details for the stored image file. Required for virtual server image for VPC.
	File *ImportOfferingBodyMetadataFile `json:"file,omitempty"`

	// Minimum size (in gigabytes) of a volume onto which this image may be provisioned. Required for virtual server image
	// for VPC.
	MinimumProvisionedSize *int64 `json:"minimum_provisioned_size,omitempty"`

	// Image operating system. Required for virtual server image for VPC.
	Images []ImportOfferingBodyMetadataImagesItem `json:"images,omitempty"`
}

// UnmarshalImportOfferingBodyMetadata unmarshals an instance of ImportOfferingBodyMetadata from the specified map of raw messages.
func UnmarshalImportOfferingBodyMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportOfferingBodyMetadata)
	err = core.UnmarshalModel(m, "operating_system", &obj.OperatingSystem, UnmarshalImportOfferingBodyMetadataOperatingSystem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "file", &obj.File, UnmarshalImportOfferingBodyMetadataFile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "minimum_provisioned_size", &obj.MinimumProvisionedSize)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "images", &obj.Images, UnmarshalImportOfferingBodyMetadataImagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportOfferingBodyMetadataFile : Details for the stored image file. Required for virtual server image for VPC.
type ImportOfferingBodyMetadataFile struct {
	// Size of the stored image file rounded up to the next gigabyte. Required for virtual server image for VPC.
	Size *int64 `json:"size,omitempty"`
}

// UnmarshalImportOfferingBodyMetadataFile unmarshals an instance of ImportOfferingBodyMetadataFile from the specified map of raw messages.
func UnmarshalImportOfferingBodyMetadataFile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportOfferingBodyMetadataFile)
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportOfferingBodyMetadataImagesItem : A list of details that identify a virtual server image. Required for virtual server image for VPC.
type ImportOfferingBodyMetadataImagesItem struct {
	// Programmatic ID of virtual server image. Required for virtual server image for VPC.
	ID *string `json:"id,omitempty"`

	// Programmatic name of virtual server image. Required for virtual server image for VPC.
	Name *string `json:"name,omitempty"`

	// Region the virtual server image is available in. Required for virtual server image for VPC.
	Region *string `json:"region,omitempty"`
}

// UnmarshalImportOfferingBodyMetadataImagesItem unmarshals an instance of ImportOfferingBodyMetadataImagesItem from the specified map of raw messages.
func UnmarshalImportOfferingBodyMetadataImagesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportOfferingBodyMetadataImagesItem)
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

// ImportOfferingBodyMetadataOperatingSystem : Operating system included in this image. Required for virtual server image for VPC.
type ImportOfferingBodyMetadataOperatingSystem struct {
	// Images with this operating system can only be used on dedicated hosts or dedicated host groups. Required for virtual
	// server image for VPC.
	DedicatedHostOnly *bool `json:"dedicated_host_only,omitempty"`

	// Vendor of the operating system. Required for virtual server image for VPC.
	Vendor *string `json:"vendor,omitempty"`

	// Globally unique name for this operating system Required for virtual server image for VPC.
	Name *string `json:"name,omitempty"`

	// URL for this operating system. Required for virtual server image for VPC.
	Href *string `json:"href,omitempty"`

	// Unique, display-friendly name for the operating system. Required for virtual server image for VPC.
	DisplayName *string `json:"display_name,omitempty"`

	// Software family for this operating system. Required for virtual server image for VPC.
	Family *string `json:"family,omitempty"`

	// Major release version of this operating system. Required for virtual server image for VPC.
	Version *string `json:"version,omitempty"`

	// Operating system architecture. Required for virtual server image for VPC.
	Architecture *string `json:"architecture,omitempty"`
}

// UnmarshalImportOfferingBodyMetadataOperatingSystem unmarshals an instance of ImportOfferingBodyMetadataOperatingSystem from the specified map of raw messages.
func UnmarshalImportOfferingBodyMetadataOperatingSystem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportOfferingBodyMetadataOperatingSystem)
	err = core.UnmarshalPrimitive(m, "dedicated_host_only", &obj.DedicatedHostOnly)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "vendor", &obj.Vendor)
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
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "family", &obj.Family)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "architecture", &obj.Architecture)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportOfferingOptions : The ImportOffering options.
type ImportOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Tags array.
	Tags []string `json:"tags,omitempty"`

	// Byte array representing the content to be imported. Only supported for OVA images at this time.
	Content *[]byte `json:"content,omitempty"`

	// Name of version. Required for virtual server image for VPC.
	Name *string `json:"name,omitempty"`

	// Display name of version. Required for virtual server image for VPC.
	Label *string `json:"label,omitempty"`

	// Install type. Example: instance. Required for virtual server image for VPC.
	InstallKind *string `json:"install_kind,omitempty"`

	// Deployment target of the content being onboarded. Current valid values are iks, roks, vcenter, power-iaas,
	// terraform, and vpc-x86. Required for virtual server image for VPC.
	TargetKinds []string `json:"target_kinds,omitempty"`

	// Format of content being onboarded. Example: vsi-image. Required for virtual server image for VPC.
	FormatKind *string `json:"format_kind,omitempty"`

	// Optional product kind for the software being onboarded.  Valid values are software, module, or solution.  Default
	// value is software.
	ProductKind *string `json:"product_kind,omitempty"`

	// SHA256 fingerprint of the image file. Required for virtual server image for VPC.
	Sha *string `json:"sha,omitempty"`

	// Semantic version of the software being onboarded. Required for virtual server image for VPC.
	Version *string `json:"version,omitempty"`

	// Version Flavor Information.  Only supported for Product kind Solution.
	Flavor *Flavor `json:"flavor,omitempty"`

	// Generic data to be included with content being onboarded. Required for virtual server image for VPC.
	Metadata *ImportOfferingBodyMetadata `json:"metadata,omitempty"`

	// Optional - The sub-folder within the specified tgz file that contains the software being onboarded.
	WorkingDirectory *string `json:"working_directory,omitempty"`

	// URL path to zip location.  If not specified, must provide content in this post body.
	Zipurl *string `json:"zipurl,omitempty"`

	// Re-use the specified offeringID during import.
	OfferingID *string `json:"offeringID,omitempty"`

	// The semver value for this new version.
	TargetVersion *string `json:"targetVersion,omitempty"`

	// Add all possible configuration items when creating this version.
	IncludeConfig *bool `json:"includeConfig,omitempty"`

	// Indicates that the current terraform template is used to install a virtual server image.
	IsVsi *bool `json:"isVSI,omitempty"`

	// The type of repository containing this version.  Valid values are 'public_git' or 'enterprise_git'.
	Repotype *string `json:"repotype,omitempty"`

	// Authentication token used to access the specified zip file.
	XAuthToken *string `json:"X-Auth-Token,omitempty"`

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
func (_options *ImportOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *ImportOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ImportOfferingOptions) SetTags(tags []string) *ImportOfferingOptions {
	_options.Tags = tags
	return _options
}

// SetContent : Allow user to set Content
func (_options *ImportOfferingOptions) SetContent(content []byte) *ImportOfferingOptions {
	_options.Content = &content
	return _options
}

// SetName : Allow user to set Name
func (_options *ImportOfferingOptions) SetName(name string) *ImportOfferingOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *ImportOfferingOptions) SetLabel(label string) *ImportOfferingOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetInstallKind : Allow user to set InstallKind
func (_options *ImportOfferingOptions) SetInstallKind(installKind string) *ImportOfferingOptions {
	_options.InstallKind = core.StringPtr(installKind)
	return _options
}

// SetTargetKinds : Allow user to set TargetKinds
func (_options *ImportOfferingOptions) SetTargetKinds(targetKinds []string) *ImportOfferingOptions {
	_options.TargetKinds = targetKinds
	return _options
}

// SetFormatKind : Allow user to set FormatKind
func (_options *ImportOfferingOptions) SetFormatKind(formatKind string) *ImportOfferingOptions {
	_options.FormatKind = core.StringPtr(formatKind)
	return _options
}

// SetProductKind : Allow user to set ProductKind
func (_options *ImportOfferingOptions) SetProductKind(productKind string) *ImportOfferingOptions {
	_options.ProductKind = core.StringPtr(productKind)
	return _options
}

// SetSha : Allow user to set Sha
func (_options *ImportOfferingOptions) SetSha(sha string) *ImportOfferingOptions {
	_options.Sha = core.StringPtr(sha)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *ImportOfferingOptions) SetVersion(version string) *ImportOfferingOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetFlavor : Allow user to set Flavor
func (_options *ImportOfferingOptions) SetFlavor(flavor *Flavor) *ImportOfferingOptions {
	_options.Flavor = flavor
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *ImportOfferingOptions) SetMetadata(metadata *ImportOfferingBodyMetadata) *ImportOfferingOptions {
	_options.Metadata = metadata
	return _options
}

// SetWorkingDirectory : Allow user to set WorkingDirectory
func (_options *ImportOfferingOptions) SetWorkingDirectory(workingDirectory string) *ImportOfferingOptions {
	_options.WorkingDirectory = core.StringPtr(workingDirectory)
	return _options
}

// SetZipurl : Allow user to set Zipurl
func (_options *ImportOfferingOptions) SetZipurl(zipurl string) *ImportOfferingOptions {
	_options.Zipurl = core.StringPtr(zipurl)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *ImportOfferingOptions) SetOfferingID(offeringID string) *ImportOfferingOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetTargetVersion : Allow user to set TargetVersion
func (_options *ImportOfferingOptions) SetTargetVersion(targetVersion string) *ImportOfferingOptions {
	_options.TargetVersion = core.StringPtr(targetVersion)
	return _options
}

// SetIncludeConfig : Allow user to set IncludeConfig
func (_options *ImportOfferingOptions) SetIncludeConfig(includeConfig bool) *ImportOfferingOptions {
	_options.IncludeConfig = core.BoolPtr(includeConfig)
	return _options
}

// SetIsVsi : Allow user to set IsVsi
func (_options *ImportOfferingOptions) SetIsVsi(isVsi bool) *ImportOfferingOptions {
	_options.IsVsi = core.BoolPtr(isVsi)
	return _options
}

// SetRepotype : Allow user to set Repotype
func (_options *ImportOfferingOptions) SetRepotype(repotype string) *ImportOfferingOptions {
	_options.Repotype = core.StringPtr(repotype)
	return _options
}

// SetXAuthToken : Allow user to set XAuthToken
func (_options *ImportOfferingOptions) SetXAuthToken(xAuthToken string) *ImportOfferingOptions {
	_options.XAuthToken = core.StringPtr(xAuthToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ImportOfferingOptions) SetHeaders(param map[string]string) *ImportOfferingOptions {
	options.Headers = param
	return options
}

// ImportOfferingVersionOptions : The ImportOfferingVersion options.
type ImportOfferingVersionOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Tags array.
	Tags []string `json:"tags,omitempty"`

	// Byte array representing the content to be imported. Only supported for OVA images at this time.
	Content *[]byte `json:"content,omitempty"`

	// Name of version. Required for virtual server image for VPC.
	Name *string `json:"name,omitempty"`

	// Display name of version. Required for virtual server image for VPC.
	Label *string `json:"label,omitempty"`

	// Install type. Example: instance. Required for virtual server image for VPC.
	InstallKind *string `json:"install_kind,omitempty"`

	// Deployment target of the content being onboarded. Current valid values are iks, roks, vcenter, power-iaas,
	// terraform, and vpc-x86. Required for virtual server image for VPC.
	TargetKinds []string `json:"target_kinds,omitempty"`

	// Format of content being onboarded. Example: vsi-image. Required for virtual server image for VPC.
	FormatKind *string `json:"format_kind,omitempty"`

	// Optional product kind for the software being onboarded.  Valid values are software, module, or solution.  Default
	// value is software.
	ProductKind *string `json:"product_kind,omitempty"`

	// SHA256 fingerprint of the image file. Required for virtual server image for VPC.
	Sha *string `json:"sha,omitempty"`

	// Semantic version of the software being onboarded. Required for virtual server image for VPC.
	Version *string `json:"version,omitempty"`

	// Version Flavor Information.  Only supported for Product kind Solution.
	Flavor *Flavor `json:"flavor,omitempty"`

	// Generic data to be included with content being onboarded. Required for virtual server image for VPC.
	Metadata *ImportOfferingBodyMetadata `json:"metadata,omitempty"`

	// Optional - The sub-folder within the specified tgz file that contains the software being onboarded.
	WorkingDirectory *string `json:"working_directory,omitempty"`

	// URL path to zip location.  If not specified, must provide content in the body of this call.
	Zipurl *string `json:"zipurl,omitempty"`

	// The semver value for this new version, if not found in the zip url package content.
	TargetVersion *string `json:"targetVersion,omitempty"`

	// Add all possible configuration values to this version when importing.
	IncludeConfig *bool `json:"includeConfig,omitempty"`

	// Indicates that the current terraform template is used to install a virtual server image.
	IsVsi *bool `json:"isVSI,omitempty"`

	// The type of repository containing this version.  Valid values are 'public_git' or 'enterprise_git'.
	Repotype *string `json:"repotype,omitempty"`

	// Authentication token used to access the specified zip file.
	XAuthToken *string `json:"X-Auth-Token,omitempty"`

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
func (_options *ImportOfferingVersionOptions) SetCatalogIdentifier(catalogIdentifier string) *ImportOfferingVersionOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *ImportOfferingVersionOptions) SetOfferingID(offeringID string) *ImportOfferingVersionOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ImportOfferingVersionOptions) SetTags(tags []string) *ImportOfferingVersionOptions {
	_options.Tags = tags
	return _options
}

// SetContent : Allow user to set Content
func (_options *ImportOfferingVersionOptions) SetContent(content []byte) *ImportOfferingVersionOptions {
	_options.Content = &content
	return _options
}

// SetName : Allow user to set Name
func (_options *ImportOfferingVersionOptions) SetName(name string) *ImportOfferingVersionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *ImportOfferingVersionOptions) SetLabel(label string) *ImportOfferingVersionOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetInstallKind : Allow user to set InstallKind
func (_options *ImportOfferingVersionOptions) SetInstallKind(installKind string) *ImportOfferingVersionOptions {
	_options.InstallKind = core.StringPtr(installKind)
	return _options
}

// SetTargetKinds : Allow user to set TargetKinds
func (_options *ImportOfferingVersionOptions) SetTargetKinds(targetKinds []string) *ImportOfferingVersionOptions {
	_options.TargetKinds = targetKinds
	return _options
}

// SetFormatKind : Allow user to set FormatKind
func (_options *ImportOfferingVersionOptions) SetFormatKind(formatKind string) *ImportOfferingVersionOptions {
	_options.FormatKind = core.StringPtr(formatKind)
	return _options
}

// SetProductKind : Allow user to set ProductKind
func (_options *ImportOfferingVersionOptions) SetProductKind(productKind string) *ImportOfferingVersionOptions {
	_options.ProductKind = core.StringPtr(productKind)
	return _options
}

// SetSha : Allow user to set Sha
func (_options *ImportOfferingVersionOptions) SetSha(sha string) *ImportOfferingVersionOptions {
	_options.Sha = core.StringPtr(sha)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *ImportOfferingVersionOptions) SetVersion(version string) *ImportOfferingVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetFlavor : Allow user to set Flavor
func (_options *ImportOfferingVersionOptions) SetFlavor(flavor *Flavor) *ImportOfferingVersionOptions {
	_options.Flavor = flavor
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *ImportOfferingVersionOptions) SetMetadata(metadata *ImportOfferingBodyMetadata) *ImportOfferingVersionOptions {
	_options.Metadata = metadata
	return _options
}

// SetWorkingDirectory : Allow user to set WorkingDirectory
func (_options *ImportOfferingVersionOptions) SetWorkingDirectory(workingDirectory string) *ImportOfferingVersionOptions {
	_options.WorkingDirectory = core.StringPtr(workingDirectory)
	return _options
}

// SetZipurl : Allow user to set Zipurl
func (_options *ImportOfferingVersionOptions) SetZipurl(zipurl string) *ImportOfferingVersionOptions {
	_options.Zipurl = core.StringPtr(zipurl)
	return _options
}

// SetTargetVersion : Allow user to set TargetVersion
func (_options *ImportOfferingVersionOptions) SetTargetVersion(targetVersion string) *ImportOfferingVersionOptions {
	_options.TargetVersion = core.StringPtr(targetVersion)
	return _options
}

// SetIncludeConfig : Allow user to set IncludeConfig
func (_options *ImportOfferingVersionOptions) SetIncludeConfig(includeConfig bool) *ImportOfferingVersionOptions {
	_options.IncludeConfig = core.BoolPtr(includeConfig)
	return _options
}

// SetIsVsi : Allow user to set IsVsi
func (_options *ImportOfferingVersionOptions) SetIsVsi(isVsi bool) *ImportOfferingVersionOptions {
	_options.IsVsi = core.BoolPtr(isVsi)
	return _options
}

// SetRepotype : Allow user to set Repotype
func (_options *ImportOfferingVersionOptions) SetRepotype(repotype string) *ImportOfferingVersionOptions {
	_options.Repotype = core.StringPtr(repotype)
	return _options
}

// SetXAuthToken : Allow user to set XAuthToken
func (_options *ImportOfferingVersionOptions) SetXAuthToken(xAuthToken string) *ImportOfferingVersionOptions {
	_options.XAuthToken = core.StringPtr(xAuthToken)
	return _options
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
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// Cluster ID.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster region.
	Region *string `json:"region,omitempty"`

	// Kube namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Validation override values. Required for virtual server image for VPC.
	OverrideValues *DeployRequestBodyOverrideValues `json:"override_values,omitempty"`

	// Schematics environment variables to use with this workspace.
	EnvironmentVariables []DeployRequestBodyEnvironmentVariablesItem `json:"environment_variables,omitempty"`

	// Entitlement API Key for this offering.
	EntitlementApikey *string `json:"entitlement_apikey,omitempty"`

	// Schematics workspace configuration.
	Schematics *DeployRequestBodySchematics `json:"schematics,omitempty"`

	// Script.
	Script *string `json:"script,omitempty"`

	// Script ID.
	ScriptID *string `json:"script_id,omitempty"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `json:"version_locator_id,omitempty"`

	// VCenter ID.
	VcenterID *string `json:"vcenter_id,omitempty"`

	// VCenter Location.
	VcenterLocation *string `json:"vcenter_location,omitempty"`

	// VCenter User.
	VcenterUser *string `json:"vcenter_user,omitempty"`

	// VCenter Password.
	VcenterPassword *string `json:"vcenter_password,omitempty"`

	// VCenter Datastore.
	VcenterDatastore *string `json:"vcenter_datastore,omitempty"`

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
func (_options *InstallVersionOptions) SetVersionLocID(versionLocID string) *InstallVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *InstallVersionOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *InstallVersionOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *InstallVersionOptions) SetClusterID(clusterID string) *InstallVersionOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *InstallVersionOptions) SetRegion(region string) *InstallVersionOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetNamespace : Allow user to set Namespace
func (_options *InstallVersionOptions) SetNamespace(namespace string) *InstallVersionOptions {
	_options.Namespace = core.StringPtr(namespace)
	return _options
}

// SetOverrideValues : Allow user to set OverrideValues
func (_options *InstallVersionOptions) SetOverrideValues(overrideValues *DeployRequestBodyOverrideValues) *InstallVersionOptions {
	_options.OverrideValues = overrideValues
	return _options
}

// SetEnvironmentVariables : Allow user to set EnvironmentVariables
func (_options *InstallVersionOptions) SetEnvironmentVariables(environmentVariables []DeployRequestBodyEnvironmentVariablesItem) *InstallVersionOptions {
	_options.EnvironmentVariables = environmentVariables
	return _options
}

// SetEntitlementApikey : Allow user to set EntitlementApikey
func (_options *InstallVersionOptions) SetEntitlementApikey(entitlementApikey string) *InstallVersionOptions {
	_options.EntitlementApikey = core.StringPtr(entitlementApikey)
	return _options
}

// SetSchematics : Allow user to set Schematics
func (_options *InstallVersionOptions) SetSchematics(schematics *DeployRequestBodySchematics) *InstallVersionOptions {
	_options.Schematics = schematics
	return _options
}

// SetScript : Allow user to set Script
func (_options *InstallVersionOptions) SetScript(script string) *InstallVersionOptions {
	_options.Script = core.StringPtr(script)
	return _options
}

// SetScriptID : Allow user to set ScriptID
func (_options *InstallVersionOptions) SetScriptID(scriptID string) *InstallVersionOptions {
	_options.ScriptID = core.StringPtr(scriptID)
	return _options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (_options *InstallVersionOptions) SetVersionLocatorID(versionLocatorID string) *InstallVersionOptions {
	_options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return _options
}

// SetVcenterID : Allow user to set VcenterID
func (_options *InstallVersionOptions) SetVcenterID(vcenterID string) *InstallVersionOptions {
	_options.VcenterID = core.StringPtr(vcenterID)
	return _options
}

// SetVcenterLocation : Allow user to set VcenterLocation
func (_options *InstallVersionOptions) SetVcenterLocation(vcenterLocation string) *InstallVersionOptions {
	_options.VcenterLocation = core.StringPtr(vcenterLocation)
	return _options
}

// SetVcenterUser : Allow user to set VcenterUser
func (_options *InstallVersionOptions) SetVcenterUser(vcenterUser string) *InstallVersionOptions {
	_options.VcenterUser = core.StringPtr(vcenterUser)
	return _options
}

// SetVcenterPassword : Allow user to set VcenterPassword
func (_options *InstallVersionOptions) SetVcenterPassword(vcenterPassword string) *InstallVersionOptions {
	_options.VcenterPassword = core.StringPtr(vcenterPassword)
	return _options
}

// SetVcenterDatastore : Allow user to set VcenterDatastore
func (_options *InstallVersionOptions) SetVcenterDatastore(vcenterDatastore string) *InstallVersionOptions {
	_options.VcenterDatastore = core.StringPtr(vcenterDatastore)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *InstallVersionOptions) SetHeaders(param map[string]string) *InstallVersionOptions {
	options.Headers = param
	return options
}

// JSONPatchOperation : A JSONPatch document as defined by RFC 6902.
type JSONPatchOperation struct {
	// The operation to be performed.
	Op *string `json:"op" validate:"required"`

	// A JSON-Pointer.
	Path *string `json:"path" validate:"required"`

	// The value to be used within the operations.
	Value interface{} `json:"value,omitempty"`

	// A string containing a JSON Pointer value.
	From *string `json:"from,omitempty"`
}

// Constants associated with the JSONPatchOperation.Op property.
// The operation to be performed.
const (
	JSONPatchOperationOpAddConst = "add"
	JSONPatchOperationOpCopyConst = "copy"
	JSONPatchOperationOpMoveConst = "move"
	JSONPatchOperationOpRemoveConst = "remove"
	JSONPatchOperationOpReplaceConst = "replace"
	JSONPatchOperationOpTestConst = "test"
)

// NewJSONPatchOperation : Instantiate JSONPatchOperation (Generic Model Constructor)
func (*CatalogManagementV1) NewJSONPatchOperation(op string, path string) (_model *JSONPatchOperation, err error) {
	_model = &JSONPatchOperation{
		Op: core.StringPtr(op),
		Path: core.StringPtr(path),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalJSONPatchOperation unmarshals an instance of JSONPatchOperation from the specified map of raw messages.
func UnmarshalJSONPatchOperation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JSONPatchOperation)
	err = core.UnmarshalPrimitive(m, "op", &obj.Op)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "from", &obj.From)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Kind : Offering kind.
type Kind struct {
	// Unique ID.
	ID *string `json:"id,omitempty"`

	// content kind, e.g., helm, vm image.
	FormatKind *string `json:"format_kind,omitempty"`

	// install kind, e.g., helm, operator, terraform.
	InstallKind *string `json:"install_kind,omitempty"`

	// target cloud to install, e.g., iks, open_shift_iks.
	TargetKind *string `json:"target_kind,omitempty"`

	// Open ended metadata information.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

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
	err = core.UnmarshalPrimitive(m, "install_kind", &obj.InstallKind)
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

// LearnMoreLinks : Learn more links for a badge.
type LearnMoreLinks struct {
	// First party link.
	FirstParty *string `json:"first_party,omitempty"`

	// Third party link.
	ThirdParty *string `json:"third_party,omitempty"`
}

// UnmarshalLearnMoreLinks unmarshals an instance of LearnMoreLinks from the specified map of raw messages.
func UnmarshalLearnMoreLinks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LearnMoreLinks)
	err = core.UnmarshalPrimitive(m, "first_party", &obj.FirstParty)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "third_party", &obj.ThirdParty)
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

// ListCatalogAccountAuditsOptions : The ListCatalogAccountAudits options.
type ListCatalogAccountAuditsOptions struct {
	// Start token for a query.
	Start *string `json:"start,omitempty"`

	// number or results to return in the query.
	Limit *int64 `json:"limit,omitempty"`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListCatalogAccountAuditsOptions : Instantiate ListCatalogAccountAuditsOptions
func (*CatalogManagementV1) NewListCatalogAccountAuditsOptions() *ListCatalogAccountAuditsOptions {
	return &ListCatalogAccountAuditsOptions{}
}

// SetStart : Allow user to set Start
func (_options *ListCatalogAccountAuditsOptions) SetStart(start string) *ListCatalogAccountAuditsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListCatalogAccountAuditsOptions) SetLimit(limit int64) *ListCatalogAccountAuditsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *ListCatalogAccountAuditsOptions) SetLookupnames(lookupnames bool) *ListCatalogAccountAuditsOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListCatalogAccountAuditsOptions) SetHeaders(param map[string]string) *ListCatalogAccountAuditsOptions {
	options.Headers = param
	return options
}

// ListCatalogAuditsOptions : The ListCatalogAudits options.
type ListCatalogAuditsOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Start token for a query.
	Start *string `json:"start,omitempty"`

	// number or results to return in the query.
	Limit *int64 `json:"limit,omitempty"`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListCatalogAuditsOptions : Instantiate ListCatalogAuditsOptions
func (*CatalogManagementV1) NewListCatalogAuditsOptions(catalogIdentifier string) *ListCatalogAuditsOptions {
	return &ListCatalogAuditsOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *ListCatalogAuditsOptions) SetCatalogIdentifier(catalogIdentifier string) *ListCatalogAuditsOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListCatalogAuditsOptions) SetStart(start string) *ListCatalogAuditsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListCatalogAuditsOptions) SetLimit(limit int64) *ListCatalogAuditsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *ListCatalogAuditsOptions) SetLookupnames(lookupnames bool) *ListCatalogAuditsOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListCatalogAuditsOptions) SetHeaders(param map[string]string) *ListCatalogAuditsOptions {
	options.Headers = param
	return options
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

// ListEnterpriseAuditsOptions : The ListEnterpriseAudits options.
type ListEnterpriseAuditsOptions struct {
	// Enterprise ID.
	EnterpriseIdentifier *string `json:"enterprise_identifier" validate:"required,ne="`

	// Start token for a query.
	Start *string `json:"start,omitempty"`

	// number or results to return in the query.
	Limit *int64 `json:"limit,omitempty"`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListEnterpriseAuditsOptions : Instantiate ListEnterpriseAuditsOptions
func (*CatalogManagementV1) NewListEnterpriseAuditsOptions(enterpriseIdentifier string) *ListEnterpriseAuditsOptions {
	return &ListEnterpriseAuditsOptions{
		EnterpriseIdentifier: core.StringPtr(enterpriseIdentifier),
	}
}

// SetEnterpriseIdentifier : Allow user to set EnterpriseIdentifier
func (_options *ListEnterpriseAuditsOptions) SetEnterpriseIdentifier(enterpriseIdentifier string) *ListEnterpriseAuditsOptions {
	_options.EnterpriseIdentifier = core.StringPtr(enterpriseIdentifier)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListEnterpriseAuditsOptions) SetStart(start string) *ListEnterpriseAuditsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListEnterpriseAuditsOptions) SetLimit(limit int64) *ListEnterpriseAuditsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *ListEnterpriseAuditsOptions) SetLookupnames(lookupnames bool) *ListEnterpriseAuditsOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListEnterpriseAuditsOptions) SetHeaders(param map[string]string) *ListEnterpriseAuditsOptions {
	options.Headers = param
	return options
}

// ListObjectAuditsOptions : The ListObjectAudits options.
type ListObjectAuditsOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// Start token for a query.
	Start *string `json:"start,omitempty"`

	// number or results to return in the query.
	Limit *int64 `json:"limit,omitempty"`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListObjectAuditsOptions : Instantiate ListObjectAuditsOptions
func (*CatalogManagementV1) NewListObjectAuditsOptions(catalogIdentifier string, objectIdentifier string) *ListObjectAuditsOptions {
	return &ListObjectAuditsOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *ListObjectAuditsOptions) SetCatalogIdentifier(catalogIdentifier string) *ListObjectAuditsOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *ListObjectAuditsOptions) SetObjectIdentifier(objectIdentifier string) *ListObjectAuditsOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListObjectAuditsOptions) SetStart(start string) *ListObjectAuditsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListObjectAuditsOptions) SetLimit(limit int64) *ListObjectAuditsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *ListObjectAuditsOptions) SetLookupnames(lookupnames bool) *ListObjectAuditsOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListObjectAuditsOptions) SetHeaders(param map[string]string) *ListObjectAuditsOptions {
	options.Headers = param
	return options
}

// ListObjectsOptions : The ListObjects options.
type ListObjectsOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// The number of results to return.
	Limit *int64 `json:"limit,omitempty"`

	// The number of results to skip before returning values.
	Offset *int64 `json:"offset,omitempty"`

	// Only return results that contain the specified string.
	Name *string `json:"name,omitempty"`

	// The field on which the output is sorted. Sorts by default by **label** property. Available fields are **name**,
	// **label**, **created**, and **updated**. By adding **-** (i.e. **-label**) in front of the query string, you can
	// specify descending order. Default is ascending order.
	Sort *string `json:"sort,omitempty"`

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
func (_options *ListObjectsOptions) SetCatalogIdentifier(catalogIdentifier string) *ListObjectsOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListObjectsOptions) SetLimit(limit int64) *ListObjectsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListObjectsOptions) SetOffset(offset int64) *ListObjectsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListObjectsOptions) SetName(name string) *ListObjectsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListObjectsOptions) SetSort(sort string) *ListObjectsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListObjectsOptions) SetHeaders(param map[string]string) *ListObjectsOptions {
	options.Headers = param
	return options
}

// ListOfferingAuditsOptions : The ListOfferingAudits options.
type ListOfferingAuditsOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Start token for a query.
	Start *string `json:"start,omitempty"`

	// number or results to return in the query.
	Limit *int64 `json:"limit,omitempty"`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListOfferingAuditsOptions : Instantiate ListOfferingAuditsOptions
func (*CatalogManagementV1) NewListOfferingAuditsOptions(catalogIdentifier string, offeringID string) *ListOfferingAuditsOptions {
	return &ListOfferingAuditsOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *ListOfferingAuditsOptions) SetCatalogIdentifier(catalogIdentifier string) *ListOfferingAuditsOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *ListOfferingAuditsOptions) SetOfferingID(offeringID string) *ListOfferingAuditsOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListOfferingAuditsOptions) SetStart(start string) *ListOfferingAuditsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListOfferingAuditsOptions) SetLimit(limit int64) *ListOfferingAuditsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *ListOfferingAuditsOptions) SetLookupnames(lookupnames bool) *ListOfferingAuditsOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListOfferingAuditsOptions) SetHeaders(param map[string]string) *ListOfferingAuditsOptions {
	options.Headers = param
	return options
}

// ListOfferingInstanceAuditsOptions : The ListOfferingInstanceAudits options.
type ListOfferingInstanceAuditsOptions struct {
	// Version Instance identifier.
	InstanceIdentifier *string `json:"instance_identifier" validate:"required,ne="`

	// Start token for a query.
	Start *string `json:"start,omitempty"`

	// number or results to return in the query.
	Limit *int64 `json:"limit,omitempty"`

	// Auditlog Lookup Names - by default names are not returned in auditlog.
	Lookupnames *bool `json:"lookupnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListOfferingInstanceAuditsOptions : Instantiate ListOfferingInstanceAuditsOptions
func (*CatalogManagementV1) NewListOfferingInstanceAuditsOptions(instanceIdentifier string) *ListOfferingInstanceAuditsOptions {
	return &ListOfferingInstanceAuditsOptions{
		InstanceIdentifier: core.StringPtr(instanceIdentifier),
	}
}

// SetInstanceIdentifier : Allow user to set InstanceIdentifier
func (_options *ListOfferingInstanceAuditsOptions) SetInstanceIdentifier(instanceIdentifier string) *ListOfferingInstanceAuditsOptions {
	_options.InstanceIdentifier = core.StringPtr(instanceIdentifier)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListOfferingInstanceAuditsOptions) SetStart(start string) *ListOfferingInstanceAuditsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListOfferingInstanceAuditsOptions) SetLimit(limit int64) *ListOfferingInstanceAuditsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetLookupnames : Allow user to set Lookupnames
func (_options *ListOfferingInstanceAuditsOptions) SetLookupnames(lookupnames bool) *ListOfferingInstanceAuditsOptions {
	_options.Lookupnames = core.BoolPtr(lookupnames)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListOfferingInstanceAuditsOptions) SetHeaders(param map[string]string) *ListOfferingInstanceAuditsOptions {
	options.Headers = param
	return options
}

// ListOfferingsOptions : The ListOfferings options.
type ListOfferingsOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// true - Strip down the content of what is returned. For example don't return the readme. Makes the result much
	// smaller. Defaults to false.
	Digest *bool `json:"digest,omitempty"`

	// The maximum number of results to return.
	Limit *int64 `json:"limit,omitempty"`

	// The number of results to skip before returning values.
	Offset *int64 `json:"offset,omitempty"`

	// Only return results that contain the specified string.
	Name *string `json:"name,omitempty"`

	// The field on which the output is sorted. Sorts by default by **label** property. Available fields are **name**,
	// **label**, **created**, and **updated**. By adding **-** (i.e. **-label**) in front of the query string, you can
	// specify descending order. Default is ascending order.
	Sort *string `json:"sort,omitempty"`

	// true - include offerings which have been marked as hidden. The default is true. To not return hidden offerings false
	// must be explicitly set.
	IncludeHidden *bool `json:"includeHidden,omitempty"`

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
func (_options *ListOfferingsOptions) SetCatalogIdentifier(catalogIdentifier string) *ListOfferingsOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetDigest : Allow user to set Digest
func (_options *ListOfferingsOptions) SetDigest(digest bool) *ListOfferingsOptions {
	_options.Digest = core.BoolPtr(digest)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListOfferingsOptions) SetLimit(limit int64) *ListOfferingsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListOfferingsOptions) SetOffset(offset int64) *ListOfferingsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListOfferingsOptions) SetName(name string) *ListOfferingsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListOfferingsOptions) SetSort(sort string) *ListOfferingsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetIncludeHidden : Allow user to set IncludeHidden
func (_options *ListOfferingsOptions) SetIncludeHidden(includeHidden bool) *ListOfferingsOptions {
	_options.IncludeHidden = core.BoolPtr(includeHidden)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListOfferingsOptions) SetHeaders(param map[string]string) *ListOfferingsOptions {
	options.Headers = param
	return options
}

// ListOperatorsOptions : The ListOperators options.
type ListOperatorsOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// Cluster identification.
	ClusterID *string `json:"cluster_id" validate:"required"`

	// Cluster region.
	Region *string `json:"region" validate:"required"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `json:"version_locator_id" validate:"required"`

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
func (_options *ListOperatorsOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *ListOperatorsOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *ListOperatorsOptions) SetClusterID(clusterID string) *ListOperatorsOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *ListOperatorsOptions) SetRegion(region string) *ListOperatorsOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (_options *ListOperatorsOptions) SetVersionLocatorID(versionLocatorID string) *ListOperatorsOptions {
	_options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListOperatorsOptions) SetHeaders(param map[string]string) *ListOperatorsOptions {
	options.Headers = param
	return options
}

// MediaItem : Offering Media information.
type MediaItem struct {
	// URL of the specified media item.
	URL *string `json:"url,omitempty"`

	// CM API specific URL of the specified media item.
	APIURL *string `json:"api_url,omitempty"`

	// Offering URL proxy information.
	URLProxy *URLProxy `json:"url_proxy,omitempty"`

	// Caption for this media item.
	Caption *string `json:"caption,omitempty"`

	// A map of translated strings, by language code.
	CaptionI18n map[string]string `json:"caption_i18n,omitempty"`

	// Type of this media item.
	Type *string `json:"type,omitempty"`

	// Thumbnail URL for this media item.
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
}

// UnmarshalMediaItem unmarshals an instance of MediaItem from the specified map of raw messages.
func UnmarshalMediaItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MediaItem)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_url", &obj.APIURL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "url_proxy", &obj.URLProxy, UnmarshalURLProxy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "caption", &obj.Caption)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "caption_i18n", &obj.CaptionI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "thumbnail_url", &obj.ThumbnailURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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

// Retrieve the value to be passed to a request to access the next page of results
func (resp *NamespaceSearchResult) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next, "offset")
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
	Resources []Access `json:"resources,omitempty"`
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
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalAccess)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ObjectAccessListResult) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next, "offset")
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

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ObjectListResult) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next, "offset")
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

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ObjectSearchResult) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next, "offset")
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

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

	// The programmatic name of this offering.
	Name *string `json:"name,omitempty"`

	// URL for an icon associated with this offering.
	OfferingIconURL *string `json:"offering_icon_url,omitempty"`

	// URL for an additional docs with this offering.
	OfferingDocsURL *string `json:"offering_docs_url,omitempty"`

	// [deprecated] - Use offering.support instead.  URL to be displayed in the Consumption UI for getting support on this
	// offering.
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

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

	// Long description in the requested language.
	LongDescription *string `json:"long_description,omitempty"`

	// A map of translated strings, by language code.
	LongDescriptionI18n map[string]string `json:"long_description_i18n,omitempty"`

	// list of features associated with this offering.
	Features []Feature `json:"features,omitempty"`

	// Array of kind.
	Kinds []Kind `json:"kinds,omitempty"`

	// Offering is managed by Partner Center.
	PcManaged *bool `json:"pc_managed,omitempty"`

	// Offering has been approved to publish to permitted to IBM or Public Catalog.
	PublishApproved *bool `json:"publish_approved,omitempty"`

	// Denotes public availability of an Offering - if share_enabled is true.
	ShareWithAll *bool `json:"share_with_all,omitempty"`

	// Denotes IBM employee availability of an Offering - if share_enabled is true.
	ShareWithIBM *bool `json:"share_with_ibm,omitempty"`

	// Denotes sharing including access list availability of an Offering is enabled.
	ShareEnabled *bool `json:"share_enabled,omitempty"`

	// Is it permitted to request publishing to IBM or Public.
	// Deprecated: this field is deprecated and may be removed in a future release.
	PermitRequestIBMPublicPublish *bool `json:"permit_request_ibm_public_publish,omitempty"`

	// Indicates if this offering has been approved for use by all IBMers.
	// Deprecated: this field is deprecated and may be removed in a future release.
	IBMPublishApproved *bool `json:"ibm_publish_approved,omitempty"`

	// Indicates if this offering has been approved for use by all IBM Cloud users.
	// Deprecated: this field is deprecated and may be removed in a future release.
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

	// Deprecated - Provider of this offering.
	// Deprecated: this field is deprecated and may be removed in a future release.
	Provider *string `json:"provider,omitempty"`

	// Information on the provider for this offering, or omitted if no provider information is given.
	ProviderInfo *ProviderInfo `json:"provider_info,omitempty"`

	// Repository info for offerings.
	RepoInfo *RepoInfo `json:"repo_info,omitempty"`

	// Image pull keys for this offering.
	ImagePullKeys []ImagePullKey `json:"image_pull_keys,omitempty"`

	// Offering Support information.
	Support *Support `json:"support,omitempty"`

	// A list of media items related to this offering.
	Media []MediaItem `json:"media,omitempty"`

	// Deprecation information for an Offering.
	DeprecatePending *DeprecatePending `json:"deprecate_pending,omitempty"`

	// The product kind.  Valid values are module, solution, or empty string.
	ProductKind *string `json:"product_kind,omitempty"`

	// A list of badges for this offering.
	Badges []Badge `json:"badges,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "label_i18n", &obj.LabelI18n)
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
	err = core.UnmarshalPrimitive(m, "short_description_i18n", &obj.ShortDescriptionI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description_i18n", &obj.LongDescriptionI18n)
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
	err = core.UnmarshalPrimitive(m, "pc_managed", &obj.PcManaged)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "publish_approved", &obj.PublishApproved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "share_with_all", &obj.ShareWithAll)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "share_with_ibm", &obj.ShareWithIBM)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "share_enabled", &obj.ShareEnabled)
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
	err = core.UnmarshalModel(m, "provider_info", &obj.ProviderInfo, UnmarshalProviderInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "repo_info", &obj.RepoInfo, UnmarshalRepoInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "image_pull_keys", &obj.ImagePullKeys, UnmarshalImagePullKey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "support", &obj.Support, UnmarshalSupport)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "media", &obj.Media, UnmarshalMediaItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "deprecate_pending", &obj.DeprecatePending, UnmarshalDeprecatePending)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "product_kind", &obj.ProductKind)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "badges", &obj.Badges, UnmarshalBadge)
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

	// The version this instance was installed from (semver - not version id).
	Version *string `json:"version,omitempty"`

	// The version id this instance was installed from (version id - not semver).
	VersionID *string `json:"version_id,omitempty"`

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

	// Type of install plan (also known as approval strategy) for operator subscriptions. Can be either automatic, which
	// automatically upgrades operators to the latest in a channel, or manual, which requires approval on the cluster.
	InstallPlan *string `json:"install_plan,omitempty"`

	// Channel to pin the operator subscription to.
	Channel *string `json:"channel,omitempty"`

	// date and time create.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// date and time updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Map of metadata values for this offering instance.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Id of the resource group to provision the offering instance into.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// String location of OfferingInstance deployment.
	Location *string `json:"location,omitempty"`

	// Indicates if Resource Controller has disabled this instance.
	Disabled *bool `json:"disabled,omitempty"`

	// The account this instance is owned by.
	Account *string `json:"account,omitempty"`

	// the last operation performed and status.
	LastOperation *OfferingInstanceLastOperation `json:"last_operation,omitempty"`

	// The target kind for the installed software version.
	KindTarget *string `json:"kind_target,omitempty"`

	// The digest value of the installed software version.
	Sha *string `json:"sha,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "version_id", &obj.VersionID)
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
	err = core.UnmarshalPrimitive(m, "install_plan", &obj.InstallPlan)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "channel", &obj.Channel)
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
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_operation", &obj.LastOperation, UnmarshalOfferingInstanceLastOperation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind_target", &obj.KindTarget)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sha", &obj.Sha)
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
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Error code from the last operation, if applicable.
	Code *string `json:"code,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OfferingReference : Offering reference definition.
type OfferingReference struct {
	// Optional - If not specified, assumes the Public Catalog.
	CatalogID *string `json:"catalog_id,omitempty"`

	// Optional - Offering ID - not required if name is set.
	ID *string `json:"id,omitempty"`

	// Optional - Programmatic Offering name.
	Name *string `json:"name,omitempty"`

	// Format kind.
	Kind *string `json:"kind,omitempty"`

	// Required - Semver value or range.
	Version *string `json:"version,omitempty"`

	// Optional - List of dependent flavors in the specified range.
	Flavors []string `json:"flavors,omitempty"`
}

// UnmarshalOfferingReference unmarshals an instance of OfferingReference from the specified map of raw messages.
func UnmarshalOfferingReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OfferingReference)
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
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
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "flavors", &obj.Flavors)
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

// Retrieve the value to be passed to a request to access the next page of results
func (resp *OfferingSearchResult) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next, "offset")
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

// Output : Outputs for a version.
type Output struct {
	// Output key.
	Key *string `json:"key,omitempty"`

	// Output description.
	Description *string `json:"description,omitempty"`
}

// UnmarshalOutput unmarshals an instance of Output from the specified map of raw messages.
func UnmarshalOutput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Output)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
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

// PaginationTokenLink : Link response on a token paginated query.
type PaginationTokenLink struct {
	// The href to the linked response.
	Href *string `json:"href" validate:"required"`

	// The start token used in this link. Will not be returned on First links.
	Start *string `json:"start,omitempty"`
}

// UnmarshalPaginationTokenLink unmarshals an instance of PaginationTokenLink from the specified map of raw messages.
func UnmarshalPaginationTokenLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginationTokenLink)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
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
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// Cluster ID.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster region.
	Region *string `json:"region,omitempty"`

	// Kube namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Validation override values. Required for virtual server image for VPC.
	OverrideValues *DeployRequestBodyOverrideValues `json:"override_values,omitempty"`

	// Schematics environment variables to use with this workspace.
	EnvironmentVariables []DeployRequestBodyEnvironmentVariablesItem `json:"environment_variables,omitempty"`

	// Entitlement API Key for this offering.
	EntitlementApikey *string `json:"entitlement_apikey,omitempty"`

	// Schematics workspace configuration.
	Schematics *DeployRequestBodySchematics `json:"schematics,omitempty"`

	// Script.
	Script *string `json:"script,omitempty"`

	// Script ID.
	ScriptID *string `json:"script_id,omitempty"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `json:"version_locator_id,omitempty"`

	// VCenter ID.
	VcenterID *string `json:"vcenter_id,omitempty"`

	// VCenter Location.
	VcenterLocation *string `json:"vcenter_location,omitempty"`

	// VCenter User.
	VcenterUser *string `json:"vcenter_user,omitempty"`

	// VCenter Password.
	VcenterPassword *string `json:"vcenter_password,omitempty"`

	// VCenter Datastore.
	VcenterDatastore *string `json:"vcenter_datastore,omitempty"`

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
func (_options *PreinstallVersionOptions) SetVersionLocID(versionLocID string) *PreinstallVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *PreinstallVersionOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *PreinstallVersionOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *PreinstallVersionOptions) SetClusterID(clusterID string) *PreinstallVersionOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *PreinstallVersionOptions) SetRegion(region string) *PreinstallVersionOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetNamespace : Allow user to set Namespace
func (_options *PreinstallVersionOptions) SetNamespace(namespace string) *PreinstallVersionOptions {
	_options.Namespace = core.StringPtr(namespace)
	return _options
}

// SetOverrideValues : Allow user to set OverrideValues
func (_options *PreinstallVersionOptions) SetOverrideValues(overrideValues *DeployRequestBodyOverrideValues) *PreinstallVersionOptions {
	_options.OverrideValues = overrideValues
	return _options
}

// SetEnvironmentVariables : Allow user to set EnvironmentVariables
func (_options *PreinstallVersionOptions) SetEnvironmentVariables(environmentVariables []DeployRequestBodyEnvironmentVariablesItem) *PreinstallVersionOptions {
	_options.EnvironmentVariables = environmentVariables
	return _options
}

// SetEntitlementApikey : Allow user to set EntitlementApikey
func (_options *PreinstallVersionOptions) SetEntitlementApikey(entitlementApikey string) *PreinstallVersionOptions {
	_options.EntitlementApikey = core.StringPtr(entitlementApikey)
	return _options
}

// SetSchematics : Allow user to set Schematics
func (_options *PreinstallVersionOptions) SetSchematics(schematics *DeployRequestBodySchematics) *PreinstallVersionOptions {
	_options.Schematics = schematics
	return _options
}

// SetScript : Allow user to set Script
func (_options *PreinstallVersionOptions) SetScript(script string) *PreinstallVersionOptions {
	_options.Script = core.StringPtr(script)
	return _options
}

// SetScriptID : Allow user to set ScriptID
func (_options *PreinstallVersionOptions) SetScriptID(scriptID string) *PreinstallVersionOptions {
	_options.ScriptID = core.StringPtr(scriptID)
	return _options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (_options *PreinstallVersionOptions) SetVersionLocatorID(versionLocatorID string) *PreinstallVersionOptions {
	_options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return _options
}

// SetVcenterID : Allow user to set VcenterID
func (_options *PreinstallVersionOptions) SetVcenterID(vcenterID string) *PreinstallVersionOptions {
	_options.VcenterID = core.StringPtr(vcenterID)
	return _options
}

// SetVcenterLocation : Allow user to set VcenterLocation
func (_options *PreinstallVersionOptions) SetVcenterLocation(vcenterLocation string) *PreinstallVersionOptions {
	_options.VcenterLocation = core.StringPtr(vcenterLocation)
	return _options
}

// SetVcenterUser : Allow user to set VcenterUser
func (_options *PreinstallVersionOptions) SetVcenterUser(vcenterUser string) *PreinstallVersionOptions {
	_options.VcenterUser = core.StringPtr(vcenterUser)
	return _options
}

// SetVcenterPassword : Allow user to set VcenterPassword
func (_options *PreinstallVersionOptions) SetVcenterPassword(vcenterPassword string) *PreinstallVersionOptions {
	_options.VcenterPassword = core.StringPtr(vcenterPassword)
	return _options
}

// SetVcenterDatastore : Allow user to set VcenterDatastore
func (_options *PreinstallVersionOptions) SetVcenterDatastore(vcenterDatastore string) *PreinstallVersionOptions {
	_options.VcenterDatastore = core.StringPtr(vcenterDatastore)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PreinstallVersionOptions) SetHeaders(param map[string]string) *PreinstallVersionOptions {
	options.Headers = param
	return options
}

// Project : Cost estimate project definition.
type Project struct {
	// Project name.
	Name *string `json:"name,omitempty"`

	// Project metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Cost breakdown definition.
	PastBreakdown *CostBreakdown `json:"pastBreakdown,omitempty"`

	// Cost breakdown definition.
	Breakdown *CostBreakdown `json:"breakdown,omitempty"`

	// Cost breakdown definition.
	Diff *CostBreakdown `json:"diff,omitempty"`

	// Cost summary definition.
	Summary *CostSummary `json:"summary,omitempty"`
}

// UnmarshalProject unmarshals an instance of Project from the specified map of raw messages.
func UnmarshalProject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Project)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "pastBreakdown", &obj.PastBreakdown, UnmarshalCostBreakdown)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "breakdown", &obj.Breakdown, UnmarshalCostBreakdown)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "diff", &obj.Diff, UnmarshalCostBreakdown)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalCostSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderInfo : Information on the provider for this offering, or omitted if no provider information is given.
type ProviderInfo struct {
	// The id of this provider.
	ID *string `json:"id,omitempty"`

	// The name of this provider.
	Name *string `json:"name,omitempty"`
}

// UnmarshalProviderInfo unmarshals an instance of ProviderInfo from the specified map of raw messages.
func UnmarshalProviderInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderInfo)
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
	InstanceIdentifier *string `json:"instance_identifier" validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

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

	// The version this instance was installed from (semver - not version id).
	Version *string `json:"version,omitempty"`

	// The version id this instance was installed from (version id - not semver).
	VersionID *string `json:"version_id,omitempty"`

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

	// Type of install plan (also known as approval strategy) for operator subscriptions. Can be either automatic, which
	// automatically upgrades operators to the latest in a channel, or manual, which requires approval on the cluster.
	InstallPlan *string `json:"install_plan,omitempty"`

	// Channel to pin the operator subscription to.
	Channel *string `json:"channel,omitempty"`

	// date and time create.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// date and time updated.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Map of metadata values for this offering instance.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Id of the resource group to provision the offering instance into.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// String location of OfferingInstance deployment.
	Location *string `json:"location,omitempty"`

	// Indicates if Resource Controller has disabled this instance.
	Disabled *bool `json:"disabled,omitempty"`

	// The account this instance is owned by.
	Account *string `json:"account,omitempty"`

	// the last operation performed and status.
	LastOperation *OfferingInstanceLastOperation `json:"last_operation,omitempty"`

	// The target kind for the installed software version.
	KindTarget *string `json:"kind_target,omitempty"`

	// The digest value of the installed software version.
	Sha *string `json:"sha,omitempty"`

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
func (_options *PutOfferingInstanceOptions) SetInstanceIdentifier(instanceIdentifier string) *PutOfferingInstanceOptions {
	_options.InstanceIdentifier = core.StringPtr(instanceIdentifier)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *PutOfferingInstanceOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *PutOfferingInstanceOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetID : Allow user to set ID
func (_options *PutOfferingInstanceOptions) SetID(id string) *PutOfferingInstanceOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *PutOfferingInstanceOptions) SetRev(rev string) *PutOfferingInstanceOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetURL : Allow user to set URL
func (_options *PutOfferingInstanceOptions) SetURL(url string) *PutOfferingInstanceOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetCRN : Allow user to set CRN
func (_options *PutOfferingInstanceOptions) SetCRN(crn string) *PutOfferingInstanceOptions {
	_options.CRN = core.StringPtr(crn)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *PutOfferingInstanceOptions) SetLabel(label string) *PutOfferingInstanceOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetCatalogID : Allow user to set CatalogID
func (_options *PutOfferingInstanceOptions) SetCatalogID(catalogID string) *PutOfferingInstanceOptions {
	_options.CatalogID = core.StringPtr(catalogID)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *PutOfferingInstanceOptions) SetOfferingID(offeringID string) *PutOfferingInstanceOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetKindFormat : Allow user to set KindFormat
func (_options *PutOfferingInstanceOptions) SetKindFormat(kindFormat string) *PutOfferingInstanceOptions {
	_options.KindFormat = core.StringPtr(kindFormat)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *PutOfferingInstanceOptions) SetVersion(version string) *PutOfferingInstanceOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetVersionID : Allow user to set VersionID
func (_options *PutOfferingInstanceOptions) SetVersionID(versionID string) *PutOfferingInstanceOptions {
	_options.VersionID = core.StringPtr(versionID)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *PutOfferingInstanceOptions) SetClusterID(clusterID string) *PutOfferingInstanceOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetClusterRegion : Allow user to set ClusterRegion
func (_options *PutOfferingInstanceOptions) SetClusterRegion(clusterRegion string) *PutOfferingInstanceOptions {
	_options.ClusterRegion = core.StringPtr(clusterRegion)
	return _options
}

// SetClusterNamespaces : Allow user to set ClusterNamespaces
func (_options *PutOfferingInstanceOptions) SetClusterNamespaces(clusterNamespaces []string) *PutOfferingInstanceOptions {
	_options.ClusterNamespaces = clusterNamespaces
	return _options
}

// SetClusterAllNamespaces : Allow user to set ClusterAllNamespaces
func (_options *PutOfferingInstanceOptions) SetClusterAllNamespaces(clusterAllNamespaces bool) *PutOfferingInstanceOptions {
	_options.ClusterAllNamespaces = core.BoolPtr(clusterAllNamespaces)
	return _options
}

// SetSchematicsWorkspaceID : Allow user to set SchematicsWorkspaceID
func (_options *PutOfferingInstanceOptions) SetSchematicsWorkspaceID(schematicsWorkspaceID string) *PutOfferingInstanceOptions {
	_options.SchematicsWorkspaceID = core.StringPtr(schematicsWorkspaceID)
	return _options
}

// SetInstallPlan : Allow user to set InstallPlan
func (_options *PutOfferingInstanceOptions) SetInstallPlan(installPlan string) *PutOfferingInstanceOptions {
	_options.InstallPlan = core.StringPtr(installPlan)
	return _options
}

// SetChannel : Allow user to set Channel
func (_options *PutOfferingInstanceOptions) SetChannel(channel string) *PutOfferingInstanceOptions {
	_options.Channel = core.StringPtr(channel)
	return _options
}

// SetCreated : Allow user to set Created
func (_options *PutOfferingInstanceOptions) SetCreated(created *strfmt.DateTime) *PutOfferingInstanceOptions {
	_options.Created = created
	return _options
}

// SetUpdated : Allow user to set Updated
func (_options *PutOfferingInstanceOptions) SetUpdated(updated *strfmt.DateTime) *PutOfferingInstanceOptions {
	_options.Updated = updated
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *PutOfferingInstanceOptions) SetMetadata(metadata map[string]interface{}) *PutOfferingInstanceOptions {
	_options.Metadata = metadata
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *PutOfferingInstanceOptions) SetResourceGroupID(resourceGroupID string) *PutOfferingInstanceOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *PutOfferingInstanceOptions) SetLocation(location string) *PutOfferingInstanceOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *PutOfferingInstanceOptions) SetDisabled(disabled bool) *PutOfferingInstanceOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetAccount : Allow user to set Account
func (_options *PutOfferingInstanceOptions) SetAccount(account string) *PutOfferingInstanceOptions {
	_options.Account = core.StringPtr(account)
	return _options
}

// SetLastOperation : Allow user to set LastOperation
func (_options *PutOfferingInstanceOptions) SetLastOperation(lastOperation *OfferingInstanceLastOperation) *PutOfferingInstanceOptions {
	_options.LastOperation = lastOperation
	return _options
}

// SetKindTarget : Allow user to set KindTarget
func (_options *PutOfferingInstanceOptions) SetKindTarget(kindTarget string) *PutOfferingInstanceOptions {
	_options.KindTarget = core.StringPtr(kindTarget)
	return _options
}

// SetSha : Allow user to set Sha
func (_options *PutOfferingInstanceOptions) SetSha(sha string) *PutOfferingInstanceOptions {
	_options.Sha = core.StringPtr(sha)
	return _options
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
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// The semver value for this new version.
	TargetVersion *string `json:"targetVersion" validate:"required"`

	// Tags array.
	Tags []string `json:"tags,omitempty"`

	// byte array representing the content to be imported.  Only supported for OVA images at this time.
	Content *[]byte `json:"content,omitempty"`

	// Target kinds.  Current valid values are 'iks', 'roks', 'vcenter', 'power-iaas', and 'terraform'.
	TargetKinds []string `json:"target_kinds,omitempty"`

	// Format of content being onboarded. Example: vsi-image. Required for virtual server image for VPC.
	FormatKind *string `json:"format_kind,omitempty"`

	// Version Flavor Information.  Only supported for Product kind Solution.
	Flavor *Flavor `json:"flavor,omitempty"`

	// Optional - The sub-folder within the specified tgz file that contains the software being onboarded.
	WorkingDirectory *string `json:"working_directory,omitempty"`

	// URL path to zip location.  If not specified, must provide content in this post body.
	Zipurl *string `json:"zipurl,omitempty"`

	// The type of repository containing this version.  Valid values are 'public_git' or 'enterprise_git'.
	RepoType *string `json:"repoType,omitempty"`

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
func (_options *ReloadOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *ReloadOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *ReloadOfferingOptions) SetOfferingID(offeringID string) *ReloadOfferingOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetTargetVersion : Allow user to set TargetVersion
func (_options *ReloadOfferingOptions) SetTargetVersion(targetVersion string) *ReloadOfferingOptions {
	_options.TargetVersion = core.StringPtr(targetVersion)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ReloadOfferingOptions) SetTags(tags []string) *ReloadOfferingOptions {
	_options.Tags = tags
	return _options
}

// SetContent : Allow user to set Content
func (_options *ReloadOfferingOptions) SetContent(content []byte) *ReloadOfferingOptions {
	_options.Content = &content
	return _options
}

// SetTargetKinds : Allow user to set TargetKinds
func (_options *ReloadOfferingOptions) SetTargetKinds(targetKinds []string) *ReloadOfferingOptions {
	_options.TargetKinds = targetKinds
	return _options
}

// SetFormatKind : Allow user to set FormatKind
func (_options *ReloadOfferingOptions) SetFormatKind(formatKind string) *ReloadOfferingOptions {
	_options.FormatKind = core.StringPtr(formatKind)
	return _options
}

// SetFlavor : Allow user to set Flavor
func (_options *ReloadOfferingOptions) SetFlavor(flavor *Flavor) *ReloadOfferingOptions {
	_options.Flavor = flavor
	return _options
}

// SetWorkingDirectory : Allow user to set WorkingDirectory
func (_options *ReloadOfferingOptions) SetWorkingDirectory(workingDirectory string) *ReloadOfferingOptions {
	_options.WorkingDirectory = core.StringPtr(workingDirectory)
	return _options
}

// SetZipurl : Allow user to set Zipurl
func (_options *ReloadOfferingOptions) SetZipurl(zipurl string) *ReloadOfferingOptions {
	_options.Zipurl = core.StringPtr(zipurl)
	return _options
}

// SetRepoType : Allow user to set RepoType
func (_options *ReloadOfferingOptions) SetRepoType(repoType string) *ReloadOfferingOptions {
	_options.RepoType = core.StringPtr(repoType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReloadOfferingOptions) SetHeaders(param map[string]string) *ReloadOfferingOptions {
	options.Headers = param
	return options
}

// RenderType : Render type.
type RenderType struct {
	// ID of the widget type.
	Type *string `json:"type,omitempty"`

	// Determines where this configuration type is rendered (3 sections today - Target, Resource, and Deployment).
	Grouping *string `json:"grouping,omitempty"`

	// Original grouping type for this configuration (3 types - Target, Resource, and Deployment).
	OriginalGrouping *string `json:"original_grouping,omitempty"`

	// Determines the order that this configuration item shows in that particular grouping.
	GroupingIndex *int64 `json:"grouping_index,omitempty"`

	// Map of constraint parameters that will be passed to the custom widget.
	ConfigConstraints map[string]interface{} `json:"config_constraints,omitempty"`

	// List of parameters that are associated with this configuration.
	Associations *RenderTypeAssociations `json:"associations,omitempty"`
}

// UnmarshalRenderType unmarshals an instance of RenderType from the specified map of raw messages.
func UnmarshalRenderType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RenderType)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "grouping", &obj.Grouping)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "original_grouping", &obj.OriginalGrouping)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "grouping_index", &obj.GroupingIndex)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "config_constraints", &obj.ConfigConstraints)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "associations", &obj.Associations, UnmarshalRenderTypeAssociations)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RenderTypeAssociations : List of parameters that are associated with this configuration.
type RenderTypeAssociations struct {
	// Parameters for this association.
	Parameters []RenderTypeAssociationsParametersItem `json:"parameters,omitempty"`
}

// UnmarshalRenderTypeAssociations unmarshals an instance of RenderTypeAssociations from the specified map of raw messages.
func UnmarshalRenderTypeAssociations(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RenderTypeAssociations)
	err = core.UnmarshalModel(m, "parameters", &obj.Parameters, UnmarshalRenderTypeAssociationsParametersItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RenderTypeAssociationsParametersItem : RenderTypeAssociationsParametersItem struct
type RenderTypeAssociationsParametersItem struct {
	// Name of this parameter.
	Name *string `json:"name,omitempty"`

	// Refresh options.
	OptionsRefresh *bool `json:"optionsRefresh,omitempty"`
}

// UnmarshalRenderTypeAssociationsParametersItem unmarshals an instance of RenderTypeAssociationsParametersItem from the specified map of raw messages.
func UnmarshalRenderTypeAssociationsParametersItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RenderTypeAssociationsParametersItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "optionsRefresh", &obj.OptionsRefresh)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceCatalogOptions : The ReplaceCatalog options.
type ReplaceCatalogOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Unique ID.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// Display Name in the requested language.
	Label *string `json:"label,omitempty"`

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

	// Description in the requested language.
	ShortDescription *string `json:"short_description,omitempty"`

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

	// URL for an icon associated with this catalog.
	CatalogIconURL *string `json:"catalog_icon_url,omitempty"`

	// URL for a banner image for this catalog.
	CatalogBannerURL *string `json:"catalog_banner_url,omitempty"`

	// List of tags associated with this catalog.
	Tags []string `json:"tags,omitempty"`

	// List of features associated with this catalog.
	Features []Feature `json:"features,omitempty"`

	// Denotes whether a catalog is disabled.
	Disabled *bool `json:"disabled,omitempty"`

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

	// Catalog specific metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

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
func (_options *ReplaceCatalogOptions) SetCatalogIdentifier(catalogIdentifier string) *ReplaceCatalogOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetID : Allow user to set ID
func (_options *ReplaceCatalogOptions) SetID(id string) *ReplaceCatalogOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *ReplaceCatalogOptions) SetRev(rev string) *ReplaceCatalogOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *ReplaceCatalogOptions) SetLabel(label string) *ReplaceCatalogOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetLabelI18n : Allow user to set LabelI18n
func (_options *ReplaceCatalogOptions) SetLabelI18n(labelI18n map[string]string) *ReplaceCatalogOptions {
	_options.LabelI18n = labelI18n
	return _options
}

// SetShortDescription : Allow user to set ShortDescription
func (_options *ReplaceCatalogOptions) SetShortDescription(shortDescription string) *ReplaceCatalogOptions {
	_options.ShortDescription = core.StringPtr(shortDescription)
	return _options
}

// SetShortDescriptionI18n : Allow user to set ShortDescriptionI18n
func (_options *ReplaceCatalogOptions) SetShortDescriptionI18n(shortDescriptionI18n map[string]string) *ReplaceCatalogOptions {
	_options.ShortDescriptionI18n = shortDescriptionI18n
	return _options
}

// SetCatalogIconURL : Allow user to set CatalogIconURL
func (_options *ReplaceCatalogOptions) SetCatalogIconURL(catalogIconURL string) *ReplaceCatalogOptions {
	_options.CatalogIconURL = core.StringPtr(catalogIconURL)
	return _options
}

// SetCatalogBannerURL : Allow user to set CatalogBannerURL
func (_options *ReplaceCatalogOptions) SetCatalogBannerURL(catalogBannerURL string) *ReplaceCatalogOptions {
	_options.CatalogBannerURL = core.StringPtr(catalogBannerURL)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ReplaceCatalogOptions) SetTags(tags []string) *ReplaceCatalogOptions {
	_options.Tags = tags
	return _options
}

// SetFeatures : Allow user to set Features
func (_options *ReplaceCatalogOptions) SetFeatures(features []Feature) *ReplaceCatalogOptions {
	_options.Features = features
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *ReplaceCatalogOptions) SetDisabled(disabled bool) *ReplaceCatalogOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *ReplaceCatalogOptions) SetResourceGroupID(resourceGroupID string) *ReplaceCatalogOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetOwningAccount : Allow user to set OwningAccount
func (_options *ReplaceCatalogOptions) SetOwningAccount(owningAccount string) *ReplaceCatalogOptions {
	_options.OwningAccount = core.StringPtr(owningAccount)
	return _options
}

// SetCatalogFilters : Allow user to set CatalogFilters
func (_options *ReplaceCatalogOptions) SetCatalogFilters(catalogFilters *Filters) *ReplaceCatalogOptions {
	_options.CatalogFilters = catalogFilters
	return _options
}

// SetSyndicationSettings : Allow user to set SyndicationSettings
func (_options *ReplaceCatalogOptions) SetSyndicationSettings(syndicationSettings *SyndicationResource) *ReplaceCatalogOptions {
	_options.SyndicationSettings = syndicationSettings
	return _options
}

// SetKind : Allow user to set Kind
func (_options *ReplaceCatalogOptions) SetKind(kind string) *ReplaceCatalogOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *ReplaceCatalogOptions) SetMetadata(metadata map[string]interface{}) *ReplaceCatalogOptions {
	_options.Metadata = metadata
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceCatalogOptions) SetHeaders(param map[string]string) *ReplaceCatalogOptions {
	options.Headers = param
	return options
}

// ReplaceObjectOptions : The ReplaceObject options.
type ReplaceObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// unique id.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// The programmatic name of this object.
	Name *string `json:"name,omitempty"`

	// The crn for this specific object.
	CRN *string `json:"crn,omitempty"`

	// The url for this specific object.
	URL *string `json:"url,omitempty"`

	// The parent for this specific object.
	ParentID *string `json:"parent_id,omitempty"`

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

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

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

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
func (_options *ReplaceObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *ReplaceObjectOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *ReplaceObjectOptions) SetObjectIdentifier(objectIdentifier string) *ReplaceObjectOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetID : Allow user to set ID
func (_options *ReplaceObjectOptions) SetID(id string) *ReplaceObjectOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *ReplaceObjectOptions) SetRev(rev string) *ReplaceObjectOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceObjectOptions) SetName(name string) *ReplaceObjectOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetCRN : Allow user to set CRN
func (_options *ReplaceObjectOptions) SetCRN(crn string) *ReplaceObjectOptions {
	_options.CRN = core.StringPtr(crn)
	return _options
}

// SetURL : Allow user to set URL
func (_options *ReplaceObjectOptions) SetURL(url string) *ReplaceObjectOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetParentID : Allow user to set ParentID
func (_options *ReplaceObjectOptions) SetParentID(parentID string) *ReplaceObjectOptions {
	_options.ParentID = core.StringPtr(parentID)
	return _options
}

// SetLabelI18n : Allow user to set LabelI18n
func (_options *ReplaceObjectOptions) SetLabelI18n(labelI18n map[string]string) *ReplaceObjectOptions {
	_options.LabelI18n = labelI18n
	return _options
}

// SetLabel : Allow user to set Label
func (_options *ReplaceObjectOptions) SetLabel(label string) *ReplaceObjectOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ReplaceObjectOptions) SetTags(tags []string) *ReplaceObjectOptions {
	_options.Tags = tags
	return _options
}

// SetCreated : Allow user to set Created
func (_options *ReplaceObjectOptions) SetCreated(created *strfmt.DateTime) *ReplaceObjectOptions {
	_options.Created = created
	return _options
}

// SetUpdated : Allow user to set Updated
func (_options *ReplaceObjectOptions) SetUpdated(updated *strfmt.DateTime) *ReplaceObjectOptions {
	_options.Updated = updated
	return _options
}

// SetShortDescription : Allow user to set ShortDescription
func (_options *ReplaceObjectOptions) SetShortDescription(shortDescription string) *ReplaceObjectOptions {
	_options.ShortDescription = core.StringPtr(shortDescription)
	return _options
}

// SetShortDescriptionI18n : Allow user to set ShortDescriptionI18n
func (_options *ReplaceObjectOptions) SetShortDescriptionI18n(shortDescriptionI18n map[string]string) *ReplaceObjectOptions {
	_options.ShortDescriptionI18n = shortDescriptionI18n
	return _options
}

// SetKind : Allow user to set Kind
func (_options *ReplaceObjectOptions) SetKind(kind string) *ReplaceObjectOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetPublish : Allow user to set Publish
func (_options *ReplaceObjectOptions) SetPublish(publish *PublishObject) *ReplaceObjectOptions {
	_options.Publish = publish
	return _options
}

// SetState : Allow user to set State
func (_options *ReplaceObjectOptions) SetState(state *State) *ReplaceObjectOptions {
	_options.State = state
	return _options
}

// SetCatalogID : Allow user to set CatalogID
func (_options *ReplaceObjectOptions) SetCatalogID(catalogID string) *ReplaceObjectOptions {
	_options.CatalogID = core.StringPtr(catalogID)
	return _options
}

// SetCatalogName : Allow user to set CatalogName
func (_options *ReplaceObjectOptions) SetCatalogName(catalogName string) *ReplaceObjectOptions {
	_options.CatalogName = core.StringPtr(catalogName)
	return _options
}

// SetData : Allow user to set Data
func (_options *ReplaceObjectOptions) SetData(data map[string]interface{}) *ReplaceObjectOptions {
	_options.Data = data
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceObjectOptions) SetHeaders(param map[string]string) *ReplaceObjectOptions {
	options.Headers = param
	return options
}

// ReplaceOfferingOptions : The ReplaceOffering options.
type ReplaceOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

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

	// A map of translated strings, by language code.
	LabelI18n map[string]string `json:"label_i18n,omitempty"`

	// The programmatic name of this offering.
	Name *string `json:"name,omitempty"`

	// URL for an icon associated with this offering.
	OfferingIconURL *string `json:"offering_icon_url,omitempty"`

	// URL for an additional docs with this offering.
	OfferingDocsURL *string `json:"offering_docs_url,omitempty"`

	// [deprecated] - Use offering.support instead.  URL to be displayed in the Consumption UI for getting support on this
	// offering.
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

	// A map of translated strings, by language code.
	ShortDescriptionI18n map[string]string `json:"short_description_i18n,omitempty"`

	// Long description in the requested language.
	LongDescription *string `json:"long_description,omitempty"`

	// A map of translated strings, by language code.
	LongDescriptionI18n map[string]string `json:"long_description_i18n,omitempty"`

	// list of features associated with this offering.
	Features []Feature `json:"features,omitempty"`

	// Array of kind.
	Kinds []Kind `json:"kinds,omitempty"`

	// Offering is managed by Partner Center.
	PcManaged *bool `json:"pc_managed,omitempty"`

	// Offering has been approved to publish to permitted to IBM or Public Catalog.
	PublishApproved *bool `json:"publish_approved,omitempty"`

	// Denotes public availability of an Offering - if share_enabled is true.
	ShareWithAll *bool `json:"share_with_all,omitempty"`

	// Denotes IBM employee availability of an Offering - if share_enabled is true.
	ShareWithIBM *bool `json:"share_with_ibm,omitempty"`

	// Denotes sharing including access list availability of an Offering is enabled.
	ShareEnabled *bool `json:"share_enabled,omitempty"`

	// Is it permitted to request publishing to IBM or Public.
	// Deprecated: this field is deprecated and may be removed in a future release.
	PermitRequestIBMPublicPublish *bool `json:"permit_request_ibm_public_publish,omitempty"`

	// Indicates if this offering has been approved for use by all IBMers.
	// Deprecated: this field is deprecated and may be removed in a future release.
	IBMPublishApproved *bool `json:"ibm_publish_approved,omitempty"`

	// Indicates if this offering has been approved for use by all IBM Cloud users.
	// Deprecated: this field is deprecated and may be removed in a future release.
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

	// Deprecated - Provider of this offering.
	// Deprecated: this field is deprecated and may be removed in a future release.
	Provider *string `json:"provider,omitempty"`

	// Information on the provider for this offering, or omitted if no provider information is given.
	ProviderInfo *ProviderInfo `json:"provider_info,omitempty"`

	// Repository info for offerings.
	RepoInfo *RepoInfo `json:"repo_info,omitempty"`

	// Image pull keys for this offering.
	ImagePullKeys []ImagePullKey `json:"image_pull_keys,omitempty"`

	// Offering Support information.
	Support *Support `json:"support,omitempty"`

	// A list of media items related to this offering.
	Media []MediaItem `json:"media,omitempty"`

	// Deprecation information for an Offering.
	DeprecatePending *DeprecatePending `json:"deprecate_pending,omitempty"`

	// The product kind.  Valid values are module, solution, or empty string.
	ProductKind *string `json:"product_kind,omitempty"`

	// A list of badges for this offering.
	Badges []Badge `json:"badges,omitempty"`

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
func (_options *ReplaceOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *ReplaceOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *ReplaceOfferingOptions) SetOfferingID(offeringID string) *ReplaceOfferingOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ReplaceOfferingOptions) SetID(id string) *ReplaceOfferingOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *ReplaceOfferingOptions) SetRev(rev string) *ReplaceOfferingOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetURL : Allow user to set URL
func (_options *ReplaceOfferingOptions) SetURL(url string) *ReplaceOfferingOptions {
	_options.URL = core.StringPtr(url)
	return _options
}

// SetCRN : Allow user to set CRN
func (_options *ReplaceOfferingOptions) SetCRN(crn string) *ReplaceOfferingOptions {
	_options.CRN = core.StringPtr(crn)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *ReplaceOfferingOptions) SetLabel(label string) *ReplaceOfferingOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetLabelI18n : Allow user to set LabelI18n
func (_options *ReplaceOfferingOptions) SetLabelI18n(labelI18n map[string]string) *ReplaceOfferingOptions {
	_options.LabelI18n = labelI18n
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceOfferingOptions) SetName(name string) *ReplaceOfferingOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetOfferingIconURL : Allow user to set OfferingIconURL
func (_options *ReplaceOfferingOptions) SetOfferingIconURL(offeringIconURL string) *ReplaceOfferingOptions {
	_options.OfferingIconURL = core.StringPtr(offeringIconURL)
	return _options
}

// SetOfferingDocsURL : Allow user to set OfferingDocsURL
func (_options *ReplaceOfferingOptions) SetOfferingDocsURL(offeringDocsURL string) *ReplaceOfferingOptions {
	_options.OfferingDocsURL = core.StringPtr(offeringDocsURL)
	return _options
}

// SetOfferingSupportURL : Allow user to set OfferingSupportURL
func (_options *ReplaceOfferingOptions) SetOfferingSupportURL(offeringSupportURL string) *ReplaceOfferingOptions {
	_options.OfferingSupportURL = core.StringPtr(offeringSupportURL)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ReplaceOfferingOptions) SetTags(tags []string) *ReplaceOfferingOptions {
	_options.Tags = tags
	return _options
}

// SetKeywords : Allow user to set Keywords
func (_options *ReplaceOfferingOptions) SetKeywords(keywords []string) *ReplaceOfferingOptions {
	_options.Keywords = keywords
	return _options
}

// SetRating : Allow user to set Rating
func (_options *ReplaceOfferingOptions) SetRating(rating *Rating) *ReplaceOfferingOptions {
	_options.Rating = rating
	return _options
}

// SetCreated : Allow user to set Created
func (_options *ReplaceOfferingOptions) SetCreated(created *strfmt.DateTime) *ReplaceOfferingOptions {
	_options.Created = created
	return _options
}

// SetUpdated : Allow user to set Updated
func (_options *ReplaceOfferingOptions) SetUpdated(updated *strfmt.DateTime) *ReplaceOfferingOptions {
	_options.Updated = updated
	return _options
}

// SetShortDescription : Allow user to set ShortDescription
func (_options *ReplaceOfferingOptions) SetShortDescription(shortDescription string) *ReplaceOfferingOptions {
	_options.ShortDescription = core.StringPtr(shortDescription)
	return _options
}

// SetShortDescriptionI18n : Allow user to set ShortDescriptionI18n
func (_options *ReplaceOfferingOptions) SetShortDescriptionI18n(shortDescriptionI18n map[string]string) *ReplaceOfferingOptions {
	_options.ShortDescriptionI18n = shortDescriptionI18n
	return _options
}

// SetLongDescription : Allow user to set LongDescription
func (_options *ReplaceOfferingOptions) SetLongDescription(longDescription string) *ReplaceOfferingOptions {
	_options.LongDescription = core.StringPtr(longDescription)
	return _options
}

// SetLongDescriptionI18n : Allow user to set LongDescriptionI18n
func (_options *ReplaceOfferingOptions) SetLongDescriptionI18n(longDescriptionI18n map[string]string) *ReplaceOfferingOptions {
	_options.LongDescriptionI18n = longDescriptionI18n
	return _options
}

// SetFeatures : Allow user to set Features
func (_options *ReplaceOfferingOptions) SetFeatures(features []Feature) *ReplaceOfferingOptions {
	_options.Features = features
	return _options
}

// SetKinds : Allow user to set Kinds
func (_options *ReplaceOfferingOptions) SetKinds(kinds []Kind) *ReplaceOfferingOptions {
	_options.Kinds = kinds
	return _options
}

// SetPcManaged : Allow user to set PcManaged
func (_options *ReplaceOfferingOptions) SetPcManaged(pcManaged bool) *ReplaceOfferingOptions {
	_options.PcManaged = core.BoolPtr(pcManaged)
	return _options
}

// SetPublishApproved : Allow user to set PublishApproved
func (_options *ReplaceOfferingOptions) SetPublishApproved(publishApproved bool) *ReplaceOfferingOptions {
	_options.PublishApproved = core.BoolPtr(publishApproved)
	return _options
}

// SetShareWithAll : Allow user to set ShareWithAll
func (_options *ReplaceOfferingOptions) SetShareWithAll(shareWithAll bool) *ReplaceOfferingOptions {
	_options.ShareWithAll = core.BoolPtr(shareWithAll)
	return _options
}

// SetShareWithIBM : Allow user to set ShareWithIBM
func (_options *ReplaceOfferingOptions) SetShareWithIBM(shareWithIBM bool) *ReplaceOfferingOptions {
	_options.ShareWithIBM = core.BoolPtr(shareWithIBM)
	return _options
}

// SetShareEnabled : Allow user to set ShareEnabled
func (_options *ReplaceOfferingOptions) SetShareEnabled(shareEnabled bool) *ReplaceOfferingOptions {
	_options.ShareEnabled = core.BoolPtr(shareEnabled)
	return _options
}

// SetPermitRequestIBMPublicPublish : Allow user to set PermitRequestIBMPublicPublish
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ReplaceOfferingOptions) SetPermitRequestIBMPublicPublish(permitRequestIBMPublicPublish bool) *ReplaceOfferingOptions {
	_options.PermitRequestIBMPublicPublish = core.BoolPtr(permitRequestIBMPublicPublish)
	return _options
}

// SetIBMPublishApproved : Allow user to set IBMPublishApproved
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ReplaceOfferingOptions) SetIBMPublishApproved(ibmPublishApproved bool) *ReplaceOfferingOptions {
	_options.IBMPublishApproved = core.BoolPtr(ibmPublishApproved)
	return _options
}

// SetPublicPublishApproved : Allow user to set PublicPublishApproved
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ReplaceOfferingOptions) SetPublicPublishApproved(publicPublishApproved bool) *ReplaceOfferingOptions {
	_options.PublicPublishApproved = core.BoolPtr(publicPublishApproved)
	return _options
}

// SetPublicOriginalCRN : Allow user to set PublicOriginalCRN
func (_options *ReplaceOfferingOptions) SetPublicOriginalCRN(publicOriginalCRN string) *ReplaceOfferingOptions {
	_options.PublicOriginalCRN = core.StringPtr(publicOriginalCRN)
	return _options
}

// SetPublishPublicCRN : Allow user to set PublishPublicCRN
func (_options *ReplaceOfferingOptions) SetPublishPublicCRN(publishPublicCRN string) *ReplaceOfferingOptions {
	_options.PublishPublicCRN = core.StringPtr(publishPublicCRN)
	return _options
}

// SetPortalApprovalRecord : Allow user to set PortalApprovalRecord
func (_options *ReplaceOfferingOptions) SetPortalApprovalRecord(portalApprovalRecord string) *ReplaceOfferingOptions {
	_options.PortalApprovalRecord = core.StringPtr(portalApprovalRecord)
	return _options
}

// SetPortalUIURL : Allow user to set PortalUIURL
func (_options *ReplaceOfferingOptions) SetPortalUIURL(portalUIURL string) *ReplaceOfferingOptions {
	_options.PortalUIURL = core.StringPtr(portalUIURL)
	return _options
}

// SetCatalogID : Allow user to set CatalogID
func (_options *ReplaceOfferingOptions) SetCatalogID(catalogID string) *ReplaceOfferingOptions {
	_options.CatalogID = core.StringPtr(catalogID)
	return _options
}

// SetCatalogName : Allow user to set CatalogName
func (_options *ReplaceOfferingOptions) SetCatalogName(catalogName string) *ReplaceOfferingOptions {
	_options.CatalogName = core.StringPtr(catalogName)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *ReplaceOfferingOptions) SetMetadata(metadata map[string]interface{}) *ReplaceOfferingOptions {
	_options.Metadata = metadata
	return _options
}

// SetDisclaimer : Allow user to set Disclaimer
func (_options *ReplaceOfferingOptions) SetDisclaimer(disclaimer string) *ReplaceOfferingOptions {
	_options.Disclaimer = core.StringPtr(disclaimer)
	return _options
}

// SetHidden : Allow user to set Hidden
func (_options *ReplaceOfferingOptions) SetHidden(hidden bool) *ReplaceOfferingOptions {
	_options.Hidden = core.BoolPtr(hidden)
	return _options
}

// SetProvider : Allow user to set Provider
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ReplaceOfferingOptions) SetProvider(provider string) *ReplaceOfferingOptions {
	_options.Provider = core.StringPtr(provider)
	return _options
}

// SetProviderInfo : Allow user to set ProviderInfo
func (_options *ReplaceOfferingOptions) SetProviderInfo(providerInfo *ProviderInfo) *ReplaceOfferingOptions {
	_options.ProviderInfo = providerInfo
	return _options
}

// SetRepoInfo : Allow user to set RepoInfo
func (_options *ReplaceOfferingOptions) SetRepoInfo(repoInfo *RepoInfo) *ReplaceOfferingOptions {
	_options.RepoInfo = repoInfo
	return _options
}

// SetImagePullKeys : Allow user to set ImagePullKeys
func (_options *ReplaceOfferingOptions) SetImagePullKeys(imagePullKeys []ImagePullKey) *ReplaceOfferingOptions {
	_options.ImagePullKeys = imagePullKeys
	return _options
}

// SetSupport : Allow user to set Support
func (_options *ReplaceOfferingOptions) SetSupport(support *Support) *ReplaceOfferingOptions {
	_options.Support = support
	return _options
}

// SetMedia : Allow user to set Media
func (_options *ReplaceOfferingOptions) SetMedia(media []MediaItem) *ReplaceOfferingOptions {
	_options.Media = media
	return _options
}

// SetDeprecatePending : Allow user to set DeprecatePending
func (_options *ReplaceOfferingOptions) SetDeprecatePending(deprecatePending *DeprecatePending) *ReplaceOfferingOptions {
	_options.DeprecatePending = deprecatePending
	return _options
}

// SetProductKind : Allow user to set ProductKind
func (_options *ReplaceOfferingOptions) SetProductKind(productKind string) *ReplaceOfferingOptions {
	_options.ProductKind = core.StringPtr(productKind)
	return _options
}

// SetBadges : Allow user to set Badges
func (_options *ReplaceOfferingOptions) SetBadges(badges []Badge) *ReplaceOfferingOptions {
	_options.Badges = badges
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceOfferingOptions) SetHeaders(param map[string]string) *ReplaceOfferingOptions {
	options.Headers = param
	return options
}

// ReplaceOperatorsOptions : The ReplaceOperators options.
type ReplaceOperatorsOptions struct {
	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// Cluster ID.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster region.
	Region *string `json:"region,omitempty"`

	// Kube namespaces to deploy Operator(s) to.
	Namespaces []string `json:"namespaces,omitempty"`

	// Denotes whether to install Operator(s) globally.
	AllNamespaces *bool `json:"all_namespaces,omitempty"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `json:"version_locator_id,omitempty"`

	// Operator channel.
	Channel *string `json:"channel,omitempty"`

	// Plan.
	InstallPlan *string `json:"install_plan,omitempty"`

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
func (_options *ReplaceOperatorsOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *ReplaceOperatorsOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *ReplaceOperatorsOptions) SetClusterID(clusterID string) *ReplaceOperatorsOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *ReplaceOperatorsOptions) SetRegion(region string) *ReplaceOperatorsOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetNamespaces : Allow user to set Namespaces
func (_options *ReplaceOperatorsOptions) SetNamespaces(namespaces []string) *ReplaceOperatorsOptions {
	_options.Namespaces = namespaces
	return _options
}

// SetAllNamespaces : Allow user to set AllNamespaces
func (_options *ReplaceOperatorsOptions) SetAllNamespaces(allNamespaces bool) *ReplaceOperatorsOptions {
	_options.AllNamespaces = core.BoolPtr(allNamespaces)
	return _options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (_options *ReplaceOperatorsOptions) SetVersionLocatorID(versionLocatorID string) *ReplaceOperatorsOptions {
	_options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return _options
}

// SetChannel : Allow user to set Channel
func (_options *ReplaceOperatorsOptions) SetChannel(channel string) *ReplaceOperatorsOptions {
	_options.Channel = core.StringPtr(channel)
	return _options
}

// SetInstallPlan : Allow user to set InstallPlan
func (_options *ReplaceOperatorsOptions) SetInstallPlan(installPlan string) *ReplaceOperatorsOptions {
	_options.InstallPlan = core.StringPtr(installPlan)
	return _options
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

	// A map of translated strings, by language code.
	InstructionsI18n map[string]string `json:"instructions_i18n,omitempty"`

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
	err = core.UnmarshalPrimitive(m, "instructions_i18n", &obj.InstructionsI18n)
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
	Query *string `json:"query" validate:"required"`

	// The kind of the object. It will default to "vpe".
	Kind *string `json:"kind,omitempty"`

	// The maximum number of results to return.
	Limit *int64 `json:"limit,omitempty"`

	// The number of results to skip before returning values.
	Offset *int64 `json:"offset,omitempty"`

	// When true, hide private objects that correspond to public or IBM published objects.
	Collapse *bool `json:"collapse,omitempty"`

	// Display a digests of search results, has default value of true.
	Digest *bool `json:"digest,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the SearchObjectsOptions.Kind property.
// The kind of the object. It will default to "vpe".
const (
	SearchObjectsOptionsKindVpeConst = "vpe"
)

// NewSearchObjectsOptions : Instantiate SearchObjectsOptions
func (*CatalogManagementV1) NewSearchObjectsOptions(query string) *SearchObjectsOptions {
	return &SearchObjectsOptions{
		Query: core.StringPtr(query),
	}
}

// SetQuery : Allow user to set Query
func (_options *SearchObjectsOptions) SetQuery(query string) *SearchObjectsOptions {
	_options.Query = core.StringPtr(query)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *SearchObjectsOptions) SetKind(kind string) *SearchObjectsOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *SearchObjectsOptions) SetLimit(limit int64) *SearchObjectsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *SearchObjectsOptions) SetOffset(offset int64) *SearchObjectsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetCollapse : Allow user to set Collapse
func (_options *SearchObjectsOptions) SetCollapse(collapse bool) *SearchObjectsOptions {
	_options.Collapse = core.BoolPtr(collapse)
	return _options
}

// SetDigest : Allow user to set Digest
func (_options *SearchObjectsOptions) SetDigest(digest bool) *SearchObjectsOptions {
	_options.Digest = core.BoolPtr(digest)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SearchObjectsOptions) SetHeaders(param map[string]string) *SearchObjectsOptions {
	options.Headers = param
	return options
}

// SetDeprecateVersionOptions : The SetDeprecateVersion options.
type SetDeprecateVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// Set deprecation (true) or cancel deprecation (false).
	Setting *string `json:"setting" validate:"required,ne="`

	// Additional information that users can provide to be displayed in deprecation notification.
	Description *string `json:"description,omitempty"`

	// Specifies the amount of days until product is not available in catalog.
	DaysUntilDeprecate *int64 `json:"days_until_deprecate,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the SetDeprecateVersionOptions.Setting property.
// Set deprecation (true) or cancel deprecation (false).
const (
	SetDeprecateVersionOptionsSettingFalseConst = "false"
	SetDeprecateVersionOptionsSettingTrueConst = "true"
)

// NewSetDeprecateVersionOptions : Instantiate SetDeprecateVersionOptions
func (*CatalogManagementV1) NewSetDeprecateVersionOptions(versionLocID string, setting string) *SetDeprecateVersionOptions {
	return &SetDeprecateVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
		Setting: core.StringPtr(setting),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (_options *SetDeprecateVersionOptions) SetVersionLocID(versionLocID string) *SetDeprecateVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetSetting : Allow user to set Setting
func (_options *SetDeprecateVersionOptions) SetSetting(setting string) *SetDeprecateVersionOptions {
	_options.Setting = core.StringPtr(setting)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *SetDeprecateVersionOptions) SetDescription(description string) *SetDeprecateVersionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetDaysUntilDeprecate : Allow user to set DaysUntilDeprecate
func (_options *SetDeprecateVersionOptions) SetDaysUntilDeprecate(daysUntilDeprecate int64) *SetDeprecateVersionOptions {
	_options.DaysUntilDeprecate = core.Int64Ptr(daysUntilDeprecate)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetDeprecateVersionOptions) SetHeaders(param map[string]string) *SetDeprecateVersionOptions {
	options.Headers = param
	return options
}

// SetOfferingPublishOptions : The SetOfferingPublish options.
type SetOfferingPublishOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Type of approval.
	//  * `pc_managed` - Partner Center is managing this offering
	//  * `publish_approved` - Publishing approved, offering owners can now set who sees the offering in public catalog
	//  * `allow_request` - (deprecated)
	//  * `ibm` - (deprecated)
	//  * `public` - (deprecated).
	ApprovalType *string `json:"approval_type" validate:"required,ne="`

	// Approve (true) or disapprove (false).
	Approved *string `json:"approved" validate:"required,ne="`

	// Partner Center identifier for this offering.
	PortalRecord *string `json:"portal_record,omitempty"`

	// Partner Center url for this offering.
	PortalURL *string `json:"portal_url,omitempty"`

	// IAM token of partner center. Only needed when Partner Center accessing the private catalog offering. When accessing
	// the public offering Partner Center only needs to use their token in the authorization header.
	XApproverToken *string `json:"X-Approver-Token,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the SetOfferingPublishOptions.ApprovalType property.
// Type of approval.
//  * `pc_managed` - Partner Center is managing this offering
//  * `publish_approved` - Publishing approved, offering owners can now set who sees the offering in public catalog
//  * `allow_request` - (deprecated)
//  * `ibm` - (deprecated)
//  * `public` - (deprecated).
const (
	SetOfferingPublishOptionsApprovalTypePcManagedConst = "pc_managed"
	SetOfferingPublishOptionsApprovalTypePublishApprovedConst = "publish_approved"
)

// Constants associated with the SetOfferingPublishOptions.Approved property.
// Approve (true) or disapprove (false).
const (
	SetOfferingPublishOptionsApprovedFalseConst = "false"
	SetOfferingPublishOptionsApprovedTrueConst = "true"
)

// NewSetOfferingPublishOptions : Instantiate SetOfferingPublishOptions
func (*CatalogManagementV1) NewSetOfferingPublishOptions(catalogIdentifier string, offeringID string, approvalType string, approved string) *SetOfferingPublishOptions {
	return &SetOfferingPublishOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		ApprovalType: core.StringPtr(approvalType),
		Approved: core.StringPtr(approved),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *SetOfferingPublishOptions) SetCatalogIdentifier(catalogIdentifier string) *SetOfferingPublishOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *SetOfferingPublishOptions) SetOfferingID(offeringID string) *SetOfferingPublishOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetApprovalType : Allow user to set ApprovalType
func (_options *SetOfferingPublishOptions) SetApprovalType(approvalType string) *SetOfferingPublishOptions {
	_options.ApprovalType = core.StringPtr(approvalType)
	return _options
}

// SetApproved : Allow user to set Approved
func (_options *SetOfferingPublishOptions) SetApproved(approved string) *SetOfferingPublishOptions {
	_options.Approved = core.StringPtr(approved)
	return _options
}

// SetPortalRecord : Allow user to set PortalRecord
func (_options *SetOfferingPublishOptions) SetPortalRecord(portalRecord string) *SetOfferingPublishOptions {
	_options.PortalRecord = core.StringPtr(portalRecord)
	return _options
}

// SetPortalURL : Allow user to set PortalURL
func (_options *SetOfferingPublishOptions) SetPortalURL(portalURL string) *SetOfferingPublishOptions {
	_options.PortalURL = core.StringPtr(portalURL)
	return _options
}

// SetXApproverToken : Allow user to set XApproverToken
func (_options *SetOfferingPublishOptions) SetXApproverToken(xApproverToken string) *SetOfferingPublishOptions {
	_options.XApproverToken = core.StringPtr(xApproverToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetOfferingPublishOptions) SetHeaders(param map[string]string) *SetOfferingPublishOptions {
	options.Headers = param
	return options
}

// ShareObjectOptions : The ShareObject options.
type ShareObjectOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Object identifier.
	ObjectIdentifier *string `json:"object_identifier" validate:"required,ne="`

	// Visible to IBM employees.
	IBM *bool `json:"ibm,omitempty"`

	// Visible to everyone in the public catalog.
	Public *bool `json:"public,omitempty"`

	// Visible to access list.
	Enabled *bool `json:"enabled,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewShareObjectOptions : Instantiate ShareObjectOptions
func (*CatalogManagementV1) NewShareObjectOptions(catalogIdentifier string, objectIdentifier string) *ShareObjectOptions {
	return &ShareObjectOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		ObjectIdentifier: core.StringPtr(objectIdentifier),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *ShareObjectOptions) SetCatalogIdentifier(catalogIdentifier string) *ShareObjectOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetObjectIdentifier : Allow user to set ObjectIdentifier
func (_options *ShareObjectOptions) SetObjectIdentifier(objectIdentifier string) *ShareObjectOptions {
	_options.ObjectIdentifier = core.StringPtr(objectIdentifier)
	return _options
}

// SetIBM : Allow user to set IBM
func (_options *ShareObjectOptions) SetIBM(ibm bool) *ShareObjectOptions {
	_options.IBM = core.BoolPtr(ibm)
	return _options
}

// SetPublic : Allow user to set Public
func (_options *ShareObjectOptions) SetPublic(public bool) *ShareObjectOptions {
	_options.Public = core.BoolPtr(public)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *ShareObjectOptions) SetEnabled(enabled bool) *ShareObjectOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ShareObjectOptions) SetHeaders(param map[string]string) *ShareObjectOptions {
	options.Headers = param
	return options
}

// ShareOfferingOptions : The ShareOffering options.
type ShareOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Visible to IBM employees.
	IBM *bool `json:"ibm,omitempty"`

	// Visible to everyone in the public catalog.
	Public *bool `json:"public,omitempty"`

	// Visible to access list.
	Enabled *bool `json:"enabled,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewShareOfferingOptions : Instantiate ShareOfferingOptions
func (*CatalogManagementV1) NewShareOfferingOptions(catalogIdentifier string, offeringID string) *ShareOfferingOptions {
	return &ShareOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *ShareOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *ShareOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *ShareOfferingOptions) SetOfferingID(offeringID string) *ShareOfferingOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetIBM : Allow user to set IBM
func (_options *ShareOfferingOptions) SetIBM(ibm bool) *ShareOfferingOptions {
	_options.IBM = core.BoolPtr(ibm)
	return _options
}

// SetPublic : Allow user to set Public
func (_options *ShareOfferingOptions) SetPublic(public bool) *ShareOfferingOptions {
	_options.Public = core.BoolPtr(public)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *ShareOfferingOptions) SetEnabled(enabled bool) *ShareOfferingOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ShareOfferingOptions) SetHeaders(param map[string]string) *ShareOfferingOptions {
	options.Headers = param
	return options
}

// ShareSetting : Share setting information.
type ShareSetting struct {
	// Visible to IBM employees.
	IBM *bool `json:"ibm,omitempty"`

	// Visible to everyone in the public catalog.
	Public *bool `json:"public,omitempty"`

	// Visible to access list.
	Enabled *bool `json:"enabled,omitempty"`
}

// UnmarshalShareSetting unmarshals an instance of ShareSetting from the specified map of raw messages.
func UnmarshalShareSetting(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ShareSetting)
	err = core.UnmarshalPrimitive(m, "ibm", &obj.IBM)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "public", &obj.Public)
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

// SolutionInfo : Version Solution Information.  Only supported for Product kind Solution.
type SolutionInfo struct {
	// Architecture diagrams for this solution.
	ArchitectureDiagrams []ArchitectureDiagram `json:"architecture_diagrams,omitempty"`

	// Features - titles only.
	Features []Feature `json:"features,omitempty"`

	// Cost estimate definition.
	CostEstimate *CostEstimate `json:"cost_estimate,omitempty"`

	// Dependencies for this solution.
	Dependencies []OfferingReference `json:"dependencies,omitempty"`

	// The install type for this solution.
	InstallType *string `json:"install_type,omitempty"`
}

// UnmarshalSolutionInfo unmarshals an instance of SolutionInfo from the specified map of raw messages.
func UnmarshalSolutionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SolutionInfo)
	err = core.UnmarshalModel(m, "architecture_diagrams", &obj.ArchitectureDiagrams, UnmarshalArchitectureDiagram)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeature)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "cost_estimate", &obj.CostEstimate, UnmarshalCostEstimate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "dependencies", &obj.Dependencies, UnmarshalOfferingReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "install_type", &obj.InstallType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// State : Offering state.
type State struct {
	// one of: new, validated, consumable.
	Current *string `json:"current,omitempty"`

	// Date and time of current request.
	CurrentEntered *strfmt.DateTime `json:"current_entered,omitempty"`

	// one of: new, validated, consumable.
	Pending *string `json:"pending,omitempty"`

	// Date and time of pending request.
	PendingRequested *strfmt.DateTime `json:"pending_requested,omitempty"`

	// one of: new, validated, consumable.
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

// Support : Offering Support information.
type Support struct {
	// URL to be displayed in the Consumption UI for getting support on this offering.
	URL *string `json:"url,omitempty"`

	// Support process as provided by an ISV.
	Process *string `json:"process,omitempty"`

	// A map of translated strings, by language code.
	ProcessI18n map[string]string `json:"process_i18n,omitempty"`

	// A list of country codes indicating where support is provided.
	Locations []string `json:"locations,omitempty"`

	// A list of support options (e.g. email, phone, slack, other).
	SupportDetails []SupportDetail `json:"support_details,omitempty"`

	// Support escalation policy.
	SupportEscalation *SupportEscalation `json:"support_escalation,omitempty"`

	// Support type for this product.
	SupportType *string `json:"support_type,omitempty"`
}

// UnmarshalSupport unmarshals an instance of Support from the specified map of raw messages.
func UnmarshalSupport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Support)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "process", &obj.Process)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "process_i18n", &obj.ProcessI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locations", &obj.Locations)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "support_details", &obj.SupportDetails, UnmarshalSupportDetail)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "support_escalation", &obj.SupportEscalation, UnmarshalSupportEscalation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "support_type", &obj.SupportType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SupportAvailability : Times when support is available.
type SupportAvailability struct {
	// A list of support times.
	Times []SupportTime `json:"times,omitempty"`

	// Timezone (e.g. America/New_York).
	Timezone *string `json:"timezone,omitempty"`

	// Is this support always available.
	AlwaysAvailable *bool `json:"always_available,omitempty"`
}

// UnmarshalSupportAvailability unmarshals an instance of SupportAvailability from the specified map of raw messages.
func UnmarshalSupportAvailability(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportAvailability)
	err = core.UnmarshalModel(m, "times", &obj.Times, UnmarshalSupportTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "always_available", &obj.AlwaysAvailable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SupportDetail : A support option.
type SupportDetail struct {
	// Type of the current support detail.
	Type *string `json:"type,omitempty"`

	// Contact for the current support detail.
	Contact *string `json:"contact,omitempty"`

	// Time descriptor.
	ResponseWaitTime *SupportWaitTime `json:"response_wait_time,omitempty"`

	// Times when support is available.
	Availability *SupportAvailability `json:"availability,omitempty"`
}

// UnmarshalSupportDetail unmarshals an instance of SupportDetail from the specified map of raw messages.
func UnmarshalSupportDetail(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportDetail)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "contact", &obj.Contact)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "response_wait_time", &obj.ResponseWaitTime, UnmarshalSupportWaitTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "availability", &obj.Availability, UnmarshalSupportAvailability)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SupportEscalation : Support escalation policy.
type SupportEscalation struct {
	// Time descriptor.
	EscalationWaitTime *SupportWaitTime `json:"escalation_wait_time,omitempty"`

	// Time descriptor.
	ResponseWaitTime *SupportWaitTime `json:"response_wait_time,omitempty"`

	// Escalation contact.
	Contact *string `json:"contact,omitempty"`
}

// UnmarshalSupportEscalation unmarshals an instance of SupportEscalation from the specified map of raw messages.
func UnmarshalSupportEscalation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportEscalation)
	err = core.UnmarshalModel(m, "escalation_wait_time", &obj.EscalationWaitTime, UnmarshalSupportWaitTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "response_wait_time", &obj.ResponseWaitTime, UnmarshalSupportWaitTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "contact", &obj.Contact)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SupportTime : Time range for support on a given day.
type SupportTime struct {
	// The day of the week, represented as an integer.
	Day *int64 `json:"day,omitempty"`

	// HOURS:MINUTES:SECONDS using 24 hour time (e.g. 8:15:00).
	StartTime *string `json:"start_time,omitempty"`

	// HOURS:MINUTES:SECONDS using 24 hour time (e.g. 8:15:00).
	EndTime *string `json:"end_time,omitempty"`
}

// UnmarshalSupportTime unmarshals an instance of SupportTime from the specified map of raw messages.
func UnmarshalSupportTime(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportTime)
	err = core.UnmarshalPrimitive(m, "day", &obj.Day)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SupportWaitTime : Time descriptor.
type SupportWaitTime struct {
	// Amount of time to wait in unit 'type'.
	Value *int64 `json:"value,omitempty"`

	// Valid values are hour or day.
	Type *string `json:"type,omitempty"`
}

// UnmarshalSupportWaitTime unmarshals an instance of SupportWaitTime from the specified map of raw messages.
func UnmarshalSupportWaitTime(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportWaitTime)
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
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

// SuspendVersionOptions : The SuspendVersion options.
type SuspendVersionOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSuspendVersionOptions : Instantiate SuspendVersionOptions
func (*CatalogManagementV1) NewSuspendVersionOptions(versionLocID string) *SuspendVersionOptions {
	return &SuspendVersionOptions{
		VersionLocID: core.StringPtr(versionLocID),
	}
}

// SetVersionLocID : Allow user to set VersionLocID
func (_options *SuspendVersionOptions) SetVersionLocID(versionLocID string) *SuspendVersionOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SuspendVersionOptions) SetHeaders(param map[string]string) *SuspendVersionOptions {
	options.Headers = param
	return options
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

// URLProxy : Offering URL proxy information.
type URLProxy struct {
	// URL of the specified media item being proxied.
	URL *string `json:"url,omitempty"`

	// SHA256 fingerprint of image.
	Sha *string `json:"sha,omitempty"`
}

// UnmarshalURLProxy unmarshals an instance of URLProxy from the specified map of raw messages.
func UnmarshalURLProxy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(URLProxy)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sha", &obj.Sha)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateCatalogAccountOptions : The UpdateCatalogAccount options.
type UpdateCatalogAccountOptions struct {
	// Account identification.
	ID *string `json:"id,omitempty"`

	// Cloudant revision.
	Rev *string `json:"_rev,omitempty"`

	// Hide the public catalog in this account.
	HideIBMCloudCatalog *bool `json:"hide_IBM_cloud_catalog,omitempty"`

	// Filters for account and catalog filters.
	AccountFilters *Filters `json:"account_filters,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateCatalogAccountOptions : Instantiate UpdateCatalogAccountOptions
func (*CatalogManagementV1) NewUpdateCatalogAccountOptions() *UpdateCatalogAccountOptions {
	return &UpdateCatalogAccountOptions{}
}

// SetID : Allow user to set ID
func (_options *UpdateCatalogAccountOptions) SetID(id string) *UpdateCatalogAccountOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *UpdateCatalogAccountOptions) SetRev(rev string) *UpdateCatalogAccountOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHideIBMCloudCatalog : Allow user to set HideIBMCloudCatalog
func (_options *UpdateCatalogAccountOptions) SetHideIBMCloudCatalog(hideIBMCloudCatalog bool) *UpdateCatalogAccountOptions {
	_options.HideIBMCloudCatalog = core.BoolPtr(hideIBMCloudCatalog)
	return _options
}

// SetAccountFilters : Allow user to set AccountFilters
func (_options *UpdateCatalogAccountOptions) SetAccountFilters(accountFilters *Filters) *UpdateCatalogAccountOptions {
	_options.AccountFilters = accountFilters
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCatalogAccountOptions) SetHeaders(param map[string]string) *UpdateCatalogAccountOptions {
	options.Headers = param
	return options
}

// UpdateOfferingOptions : The UpdateOffering options.
type UpdateOfferingOptions struct {
	// Catalog identifier.
	CatalogIdentifier *string `json:"catalog_identifier" validate:"required,ne="`

	// Offering identification.
	OfferingID *string `json:"offering_id" validate:"required,ne="`

	// Offering etag contained in quotes.
	IfMatch *string `json:"If-Match" validate:"required"`

	Updates []JSONPatchOperation `json:"updates,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateOfferingOptions : Instantiate UpdateOfferingOptions
func (*CatalogManagementV1) NewUpdateOfferingOptions(catalogIdentifier string, offeringID string, ifMatch string) *UpdateOfferingOptions {
	return &UpdateOfferingOptions{
		CatalogIdentifier: core.StringPtr(catalogIdentifier),
		OfferingID: core.StringPtr(offeringID),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetCatalogIdentifier : Allow user to set CatalogIdentifier
func (_options *UpdateOfferingOptions) SetCatalogIdentifier(catalogIdentifier string) *UpdateOfferingOptions {
	_options.CatalogIdentifier = core.StringPtr(catalogIdentifier)
	return _options
}

// SetOfferingID : Allow user to set OfferingID
func (_options *UpdateOfferingOptions) SetOfferingID(offeringID string) *UpdateOfferingOptions {
	_options.OfferingID = core.StringPtr(offeringID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateOfferingOptions) SetIfMatch(ifMatch string) *UpdateOfferingOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetUpdates : Allow user to set Updates
func (_options *UpdateOfferingOptions) SetUpdates(updates []JSONPatchOperation) *UpdateOfferingOptions {
	_options.Updates = updates
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateOfferingOptions) SetHeaders(param map[string]string) *UpdateOfferingOptions {
	options.Headers = param
	return options
}

// ValidateInstallOptions : The ValidateInstall options.
type ValidateInstallOptions struct {
	// A dotted value of `catalogID`.`versionID`.
	VersionLocID *string `json:"version_loc_id" validate:"required,ne="`

	// IAM Refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token" validate:"required"`

	// Cluster ID.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster region.
	Region *string `json:"region,omitempty"`

	// Kube namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Validation override values. Required for virtual server image for VPC.
	OverrideValues *DeployRequestBodyOverrideValues `json:"override_values,omitempty"`

	// Schematics environment variables to use with this workspace.
	EnvironmentVariables []DeployRequestBodyEnvironmentVariablesItem `json:"environment_variables,omitempty"`

	// Entitlement API Key for this offering.
	EntitlementApikey *string `json:"entitlement_apikey,omitempty"`

	// Schematics workspace configuration.
	Schematics *DeployRequestBodySchematics `json:"schematics,omitempty"`

	// Script.
	Script *string `json:"script,omitempty"`

	// Script ID.
	ScriptID *string `json:"script_id,omitempty"`

	// A dotted value of `catalogID`.`versionID`.
	VersionLocatorID *string `json:"version_locator_id,omitempty"`

	// VCenter ID.
	VcenterID *string `json:"vcenter_id,omitempty"`

	// VCenter Location.
	VcenterLocation *string `json:"vcenter_location,omitempty"`

	// VCenter User.
	VcenterUser *string `json:"vcenter_user,omitempty"`

	// VCenter Password.
	VcenterPassword *string `json:"vcenter_password,omitempty"`

	// VCenter Datastore.
	VcenterDatastore *string `json:"vcenter_datastore,omitempty"`

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
func (_options *ValidateInstallOptions) SetVersionLocID(versionLocID string) *ValidateInstallOptions {
	_options.VersionLocID = core.StringPtr(versionLocID)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *ValidateInstallOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *ValidateInstallOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetClusterID : Allow user to set ClusterID
func (_options *ValidateInstallOptions) SetClusterID(clusterID string) *ValidateInstallOptions {
	_options.ClusterID = core.StringPtr(clusterID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *ValidateInstallOptions) SetRegion(region string) *ValidateInstallOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetNamespace : Allow user to set Namespace
func (_options *ValidateInstallOptions) SetNamespace(namespace string) *ValidateInstallOptions {
	_options.Namespace = core.StringPtr(namespace)
	return _options
}

// SetOverrideValues : Allow user to set OverrideValues
func (_options *ValidateInstallOptions) SetOverrideValues(overrideValues *DeployRequestBodyOverrideValues) *ValidateInstallOptions {
	_options.OverrideValues = overrideValues
	return _options
}

// SetEnvironmentVariables : Allow user to set EnvironmentVariables
func (_options *ValidateInstallOptions) SetEnvironmentVariables(environmentVariables []DeployRequestBodyEnvironmentVariablesItem) *ValidateInstallOptions {
	_options.EnvironmentVariables = environmentVariables
	return _options
}

// SetEntitlementApikey : Allow user to set EntitlementApikey
func (_options *ValidateInstallOptions) SetEntitlementApikey(entitlementApikey string) *ValidateInstallOptions {
	_options.EntitlementApikey = core.StringPtr(entitlementApikey)
	return _options
}

// SetSchematics : Allow user to set Schematics
func (_options *ValidateInstallOptions) SetSchematics(schematics *DeployRequestBodySchematics) *ValidateInstallOptions {
	_options.Schematics = schematics
	return _options
}

// SetScript : Allow user to set Script
func (_options *ValidateInstallOptions) SetScript(script string) *ValidateInstallOptions {
	_options.Script = core.StringPtr(script)
	return _options
}

// SetScriptID : Allow user to set ScriptID
func (_options *ValidateInstallOptions) SetScriptID(scriptID string) *ValidateInstallOptions {
	_options.ScriptID = core.StringPtr(scriptID)
	return _options
}

// SetVersionLocatorID : Allow user to set VersionLocatorID
func (_options *ValidateInstallOptions) SetVersionLocatorID(versionLocatorID string) *ValidateInstallOptions {
	_options.VersionLocatorID = core.StringPtr(versionLocatorID)
	return _options
}

// SetVcenterID : Allow user to set VcenterID
func (_options *ValidateInstallOptions) SetVcenterID(vcenterID string) *ValidateInstallOptions {
	_options.VcenterID = core.StringPtr(vcenterID)
	return _options
}

// SetVcenterLocation : Allow user to set VcenterLocation
func (_options *ValidateInstallOptions) SetVcenterLocation(vcenterLocation string) *ValidateInstallOptions {
	_options.VcenterLocation = core.StringPtr(vcenterLocation)
	return _options
}

// SetVcenterUser : Allow user to set VcenterUser
func (_options *ValidateInstallOptions) SetVcenterUser(vcenterUser string) *ValidateInstallOptions {
	_options.VcenterUser = core.StringPtr(vcenterUser)
	return _options
}

// SetVcenterPassword : Allow user to set VcenterPassword
func (_options *ValidateInstallOptions) SetVcenterPassword(vcenterPassword string) *ValidateInstallOptions {
	_options.VcenterPassword = core.StringPtr(vcenterPassword)
	return _options
}

// SetVcenterDatastore : Allow user to set VcenterDatastore
func (_options *ValidateInstallOptions) SetVcenterDatastore(vcenterDatastore string) *ValidateInstallOptions {
	_options.VcenterDatastore = core.StringPtr(vcenterDatastore)
	return _options
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

	// Any message needing to be conveyed as part of the validation job.
	Message *string `json:"message,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
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

	// Version Flavor Information.  Only supported for Product kind Solution.
	Flavor *Flavor `json:"flavor,omitempty"`

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

	// List of output values for this version.
	Outputs []Output `json:"outputs,omitempty"`

	// List of IAM permissions that are required to consume this version.
	IamPermissions []IamPermission `json:"iam_permissions,omitempty"`

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

	// Long description for version.
	LongDescription *string `json:"long_description,omitempty"`

	// A map of translated strings, by language code.
	LongDescriptionI18n map[string]string `json:"long_description_i18n,omitempty"`

	// Whitelisted accounts for version.
	WhitelistedAccounts []string `json:"whitelisted_accounts,omitempty"`

	// ID of the image pull key to use from Offering.ImagePullKeys.
	ImagePullKeyName *string `json:"image_pull_key_name,omitempty"`

	// Deprecation information for an Offering.
	DeprecatePending *DeprecatePending `json:"deprecate_pending,omitempty"`

	// Version Solution Information.  Only supported for Product kind Solution.
	SolutionInfo *SolutionInfo `json:"solution_info,omitempty"`

	// Is the version able to be shared.
	IsConsumable *bool `json:"is_consumable,omitempty"`

	// List of links to sec./compliance controls.
	Compliance []ComplianceControl `json:"compliance,omitempty"`
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
	err = core.UnmarshalModel(m, "flavor", &obj.Flavor, UnmarshalFlavor)
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
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalOutput)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "iam_permissions", &obj.IamPermissions, UnmarshalIamPermission)
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
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description_i18n", &obj.LongDescriptionI18n)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "whitelisted_accounts", &obj.WhitelistedAccounts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "image_pull_key_name", &obj.ImagePullKeyName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "deprecate_pending", &obj.DeprecatePending, UnmarshalDeprecatePending)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "solution_info", &obj.SolutionInfo, UnmarshalSolutionInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_consumable", &obj.IsConsumable)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance", &obj.Compliance, UnmarshalComplianceControl)
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

	// Version Flavor Information.  Only supported for Product kind Solution.
	Flavor *Flavor `json:"flavor,omitempty"`

	// Offering state.
	State *State `json:"state,omitempty"`

	// Resource requirments for installation.
	RequiredResources []Resource `json:"required_resources,omitempty"`

	// Version of package.
	PackageVersion *string `json:"package_version,omitempty"`

	// The SHA value of this version.
	Sha *string `json:"sha,omitempty"`

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
	err = core.UnmarshalModel(m, "flavor", &obj.Flavor, UnmarshalFlavor)
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
	err = core.UnmarshalPrimitive(m, "sha", &obj.Sha)
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

//
// CatalogAccountAuditsPager can be used to simplify the use of the "ListCatalogAccountAudits" method.
//
type CatalogAccountAuditsPager struct {
	hasNext bool
	options *ListCatalogAccountAuditsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *string
	}
}

// NewCatalogAccountAuditsPager returns a new CatalogAccountAuditsPager instance.
func (catalogManagement *CatalogManagementV1) NewCatalogAccountAuditsPager(options *ListCatalogAccountAuditsOptions) (pager *CatalogAccountAuditsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListCatalogAccountAuditsOptions = *options
	pager = &CatalogAccountAuditsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *CatalogAccountAuditsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *CatalogAccountAuditsPager) GetNextWithContext(ctx context.Context) (page []AuditLogDigest, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListCatalogAccountAuditsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Audits

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *CatalogAccountAuditsPager) GetAllWithContext(ctx context.Context) (allItems []AuditLogDigest, err error) {
	for pager.HasNext() {
		var nextPage []AuditLogDigest
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *CatalogAccountAuditsPager) GetNext() (page []AuditLogDigest, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *CatalogAccountAuditsPager) GetAll() (allItems []AuditLogDigest, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// CatalogAuditsPager can be used to simplify the use of the "ListCatalogAudits" method.
//
type CatalogAuditsPager struct {
	hasNext bool
	options *ListCatalogAuditsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *string
	}
}

// NewCatalogAuditsPager returns a new CatalogAuditsPager instance.
func (catalogManagement *CatalogManagementV1) NewCatalogAuditsPager(options *ListCatalogAuditsOptions) (pager *CatalogAuditsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListCatalogAuditsOptions = *options
	pager = &CatalogAuditsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *CatalogAuditsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *CatalogAuditsPager) GetNextWithContext(ctx context.Context) (page []AuditLogDigest, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListCatalogAuditsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Audits

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *CatalogAuditsPager) GetAllWithContext(ctx context.Context) (allItems []AuditLogDigest, err error) {
	for pager.HasNext() {
		var nextPage []AuditLogDigest
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *CatalogAuditsPager) GetNext() (page []AuditLogDigest, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *CatalogAuditsPager) GetAll() (allItems []AuditLogDigest, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// EnterpriseAuditsPager can be used to simplify the use of the "ListEnterpriseAudits" method.
//
type EnterpriseAuditsPager struct {
	hasNext bool
	options *ListEnterpriseAuditsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *string
	}
}

// NewEnterpriseAuditsPager returns a new EnterpriseAuditsPager instance.
func (catalogManagement *CatalogManagementV1) NewEnterpriseAuditsPager(options *ListEnterpriseAuditsOptions) (pager *EnterpriseAuditsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListEnterpriseAuditsOptions = *options
	pager = &EnterpriseAuditsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *EnterpriseAuditsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *EnterpriseAuditsPager) GetNextWithContext(ctx context.Context) (page []AuditLogDigest, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListEnterpriseAuditsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Audits

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *EnterpriseAuditsPager) GetAllWithContext(ctx context.Context) (allItems []AuditLogDigest, err error) {
	for pager.HasNext() {
		var nextPage []AuditLogDigest
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *EnterpriseAuditsPager) GetNext() (page []AuditLogDigest, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *EnterpriseAuditsPager) GetAll() (allItems []AuditLogDigest, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// GetConsumptionOfferingsPager can be used to simplify the use of the "GetConsumptionOfferings" method.
//
type GetConsumptionOfferingsPager struct {
	hasNext bool
	options *GetConsumptionOfferingsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *int64
	}
}

// NewGetConsumptionOfferingsPager returns a new GetConsumptionOfferingsPager instance.
func (catalogManagement *CatalogManagementV1) NewGetConsumptionOfferingsPager(options *GetConsumptionOfferingsOptions) (pager *GetConsumptionOfferingsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy GetConsumptionOfferingsOptions = *options
	pager = &GetConsumptionOfferingsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetConsumptionOfferingsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetConsumptionOfferingsPager) GetNextWithContext(ctx context.Context) (page []Offering, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.GetConsumptionOfferingsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetConsumptionOfferingsPager) GetAllWithContext(ctx context.Context) (allItems []Offering, err error) {
	for pager.HasNext() {
		var nextPage []Offering
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *GetConsumptionOfferingsPager) GetNext() (page []Offering, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetConsumptionOfferingsPager) GetAll() (allItems []Offering, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// OfferingsPager can be used to simplify the use of the "ListOfferings" method.
//
type OfferingsPager struct {
	hasNext bool
	options *ListOfferingsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *int64
	}
}

// NewOfferingsPager returns a new OfferingsPager instance.
func (catalogManagement *CatalogManagementV1) NewOfferingsPager(options *ListOfferingsOptions) (pager *OfferingsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListOfferingsOptions = *options
	pager = &OfferingsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *OfferingsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *OfferingsPager) GetNextWithContext(ctx context.Context) (page []Offering, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListOfferingsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *OfferingsPager) GetAllWithContext(ctx context.Context) (allItems []Offering, err error) {
	for pager.HasNext() {
		var nextPage []Offering
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *OfferingsPager) GetNext() (page []Offering, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *OfferingsPager) GetAll() (allItems []Offering, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// OfferingAuditsPager can be used to simplify the use of the "ListOfferingAudits" method.
//
type OfferingAuditsPager struct {
	hasNext bool
	options *ListOfferingAuditsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *string
	}
}

// NewOfferingAuditsPager returns a new OfferingAuditsPager instance.
func (catalogManagement *CatalogManagementV1) NewOfferingAuditsPager(options *ListOfferingAuditsOptions) (pager *OfferingAuditsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListOfferingAuditsOptions = *options
	pager = &OfferingAuditsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *OfferingAuditsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *OfferingAuditsPager) GetNextWithContext(ctx context.Context) (page []AuditLogDigest, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListOfferingAuditsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Audits

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *OfferingAuditsPager) GetAllWithContext(ctx context.Context) (allItems []AuditLogDigest, err error) {
	for pager.HasNext() {
		var nextPage []AuditLogDigest
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *OfferingAuditsPager) GetNext() (page []AuditLogDigest, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *OfferingAuditsPager) GetAll() (allItems []AuditLogDigest, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// GetOfferingAccessListPager can be used to simplify the use of the "GetOfferingAccessList" method.
//
type GetOfferingAccessListPager struct {
	hasNext bool
	options *GetOfferingAccessListOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *string
	}
}

// NewGetOfferingAccessListPager returns a new GetOfferingAccessListPager instance.
func (catalogManagement *CatalogManagementV1) NewGetOfferingAccessListPager(options *GetOfferingAccessListOptions) (pager *GetOfferingAccessListPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy GetOfferingAccessListOptions = *options
	pager = &GetOfferingAccessListPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetOfferingAccessListPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetOfferingAccessListPager) GetNextWithContext(ctx context.Context) (page []Access, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.GetOfferingAccessListWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetOfferingAccessListPager) GetAllWithContext(ctx context.Context) (allItems []Access, err error) {
	for pager.HasNext() {
		var nextPage []Access
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *GetOfferingAccessListPager) GetNext() (page []Access, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetOfferingAccessListPager) GetAll() (allItems []Access, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// GetNamespacesPager can be used to simplify the use of the "GetNamespaces" method.
//
type GetNamespacesPager struct {
	hasNext bool
	options *GetNamespacesOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *int64
	}
}

// NewGetNamespacesPager returns a new GetNamespacesPager instance.
func (catalogManagement *CatalogManagementV1) NewGetNamespacesPager(options *GetNamespacesOptions) (pager *GetNamespacesPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy GetNamespacesOptions = *options
	pager = &GetNamespacesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetNamespacesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetNamespacesPager) GetNextWithContext(ctx context.Context) (page []string, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.GetNamespacesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetNamespacesPager) GetAllWithContext(ctx context.Context) (allItems []string, err error) {
	for pager.HasNext() {
		var nextPage []string
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *GetNamespacesPager) GetNext() (page []string, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetNamespacesPager) GetAll() (allItems []string, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// SearchObjectsPager can be used to simplify the use of the "SearchObjects" method.
//
type SearchObjectsPager struct {
	hasNext bool
	options *SearchObjectsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *int64
	}
}

// NewSearchObjectsPager returns a new SearchObjectsPager instance.
func (catalogManagement *CatalogManagementV1) NewSearchObjectsPager(options *SearchObjectsOptions) (pager *SearchObjectsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy SearchObjectsOptions = *options
	pager = &SearchObjectsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SearchObjectsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SearchObjectsPager) GetNextWithContext(ctx context.Context) (page []CatalogObject, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.SearchObjectsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SearchObjectsPager) GetAllWithContext(ctx context.Context) (allItems []CatalogObject, err error) {
	for pager.HasNext() {
		var nextPage []CatalogObject
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *SearchObjectsPager) GetNext() (page []CatalogObject, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SearchObjectsPager) GetAll() (allItems []CatalogObject, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// ObjectsPager can be used to simplify the use of the "ListObjects" method.
//
type ObjectsPager struct {
	hasNext bool
	options *ListObjectsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *int64
	}
}

// NewObjectsPager returns a new ObjectsPager instance.
func (catalogManagement *CatalogManagementV1) NewObjectsPager(options *ListObjectsOptions) (pager *ObjectsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy ListObjectsOptions = *options
	pager = &ObjectsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ObjectsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ObjectsPager) GetNextWithContext(ctx context.Context) (page []CatalogObject, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListObjectsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ObjectsPager) GetAllWithContext(ctx context.Context) (allItems []CatalogObject, err error) {
	for pager.HasNext() {
		var nextPage []CatalogObject
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ObjectsPager) GetNext() (page []CatalogObject, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ObjectsPager) GetAll() (allItems []CatalogObject, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// ObjectAuditsPager can be used to simplify the use of the "ListObjectAudits" method.
//
type ObjectAuditsPager struct {
	hasNext bool
	options *ListObjectAuditsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *string
	}
}

// NewObjectAuditsPager returns a new ObjectAuditsPager instance.
func (catalogManagement *CatalogManagementV1) NewObjectAuditsPager(options *ListObjectAuditsOptions) (pager *ObjectAuditsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListObjectAuditsOptions = *options
	pager = &ObjectAuditsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ObjectAuditsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ObjectAuditsPager) GetNextWithContext(ctx context.Context) (page []AuditLogDigest, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListObjectAuditsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Audits

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ObjectAuditsPager) GetAllWithContext(ctx context.Context) (allItems []AuditLogDigest, err error) {
	for pager.HasNext() {
		var nextPage []AuditLogDigest
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ObjectAuditsPager) GetNext() (page []AuditLogDigest, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ObjectAuditsPager) GetAll() (allItems []AuditLogDigest, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// GetObjectAccessListPager can be used to simplify the use of the "GetObjectAccessList" method.
//
type GetObjectAccessListPager struct {
	hasNext bool
	options *GetObjectAccessListOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *string
	}
}

// NewGetObjectAccessListPager returns a new GetObjectAccessListPager instance.
func (catalogManagement *CatalogManagementV1) NewGetObjectAccessListPager(options *GetObjectAccessListOptions) (pager *GetObjectAccessListPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy GetObjectAccessListOptions = *options
	pager = &GetObjectAccessListPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetObjectAccessListPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetObjectAccessListPager) GetNextWithContext(ctx context.Context) (page []Access, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.GetObjectAccessListWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetObjectAccessListPager) GetAllWithContext(ctx context.Context) (allItems []Access, err error) {
	for pager.HasNext() {
		var nextPage []Access
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *GetObjectAccessListPager) GetNext() (page []Access, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetObjectAccessListPager) GetAll() (allItems []Access, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// GetObjectAccessListDeprecatedPager can be used to simplify the use of the "GetObjectAccessListDeprecated" method.
//
type GetObjectAccessListDeprecatedPager struct {
	hasNext bool
	options *GetObjectAccessListDeprecatedOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *int64
	}
}

// NewGetObjectAccessListDeprecatedPager returns a new GetObjectAccessListDeprecatedPager instance.
func (catalogManagement *CatalogManagementV1) NewGetObjectAccessListDeprecatedPager(options *GetObjectAccessListDeprecatedOptions) (pager *GetObjectAccessListDeprecatedPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = fmt.Errorf("the 'options.Offset' field should not be set")
		return
	}

	var optionsCopy GetObjectAccessListDeprecatedOptions = *options
	pager = &GetObjectAccessListDeprecatedPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetObjectAccessListDeprecatedPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetObjectAccessListDeprecatedPager) GetNextWithContext(ctx context.Context) (page []Access, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.GetObjectAccessListDeprecatedWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next, "offset")
		if err != nil {
			err = fmt.Errorf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next, err.Error())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetObjectAccessListDeprecatedPager) GetAllWithContext(ctx context.Context) (allItems []Access, err error) {
	for pager.HasNext() {
		var nextPage []Access
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *GetObjectAccessListDeprecatedPager) GetNext() (page []Access, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetObjectAccessListDeprecatedPager) GetAll() (allItems []Access, err error) {
	return pager.GetAllWithContext(context.Background())
}

//
// OfferingInstanceAuditsPager can be used to simplify the use of the "ListOfferingInstanceAudits" method.
//
type OfferingInstanceAuditsPager struct {
	hasNext bool
	options *ListOfferingInstanceAuditsOptions
	client  *CatalogManagementV1
	pageContext struct {
		next *string
	}
}

// NewOfferingInstanceAuditsPager returns a new OfferingInstanceAuditsPager instance.
func (catalogManagement *CatalogManagementV1) NewOfferingInstanceAuditsPager(options *ListOfferingInstanceAuditsOptions) (pager *OfferingInstanceAuditsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListOfferingInstanceAuditsOptions = *options
	pager = &OfferingInstanceAuditsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  catalogManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *OfferingInstanceAuditsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *OfferingInstanceAuditsPager) GetNextWithContext(ctx context.Context) (page []AuditLogDigest, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListOfferingInstanceAuditsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Audits

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *OfferingInstanceAuditsPager) GetAllWithContext(ctx context.Context) (allItems []AuditLogDigest, err error) {
	for pager.HasNext() {
		var nextPage []AuditLogDigest
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *OfferingInstanceAuditsPager) GetNext() (page []AuditLogDigest, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *OfferingInstanceAuditsPager) GetAll() (allItems []AuditLogDigest, err error) {
	return pager.GetAllWithContext(context.Background())
}
