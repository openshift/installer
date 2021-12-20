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
 * IBM OpenAPI SDK Code Generator Version: 3.40.0-910cf8c2-20211006-154754
 */

// Package findingsv1 : Operations and models for the FindingsV1 service
package findingsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/scc-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// FindingsV1 : The Findings API is used to find and display occurrences of security issues in your IBM Cloud account by
// using the artifact metadata specification. Findings are summarized in cards in the Security and Compliance Center
// that allow you to see the security status of your account at a glance and start an investigation into any potential
// issues.
//
// API Version: 1.0.0
type FindingsV1 struct {
	Service *core.BaseService

	// Account ID.
	AccountID *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.secadvisor.cloud.ibm.com/findings"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "findings"

// FindingsV1Options : Service options
type FindingsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Account ID.
	AccountID *string `validate:"required"`
}

// NewFindingsV1UsingExternalConfig : constructs an instance of FindingsV1 with passed in options and external configuration.
func NewFindingsV1UsingExternalConfig(options *FindingsV1Options) (findings *FindingsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	findings, err = NewFindingsV1(options)
	if err != nil {
		return
	}

	err = findings.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = findings.Service.SetServiceURL(options.URL)
	}
	return
}

// NewFindingsV1 : constructs an instance of FindingsV1 with passed in options.
func NewFindingsV1(options *FindingsV1Options) (service *FindingsV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
	}

	err = core.ValidateStruct(options, "options")
	if err != nil {
		return
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

	service = &FindingsV1{
		Service:   baseService,
		AccountID: options.AccountID,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south": "https://us-south.secadvisor.cloud.ibm.com/findings",
		"us-east":  "https://us-south.secadvisor.cloud.ibm.com/findings",
		"eu-gb":    "https://eu-gb.secadvisor.cloud.ibm.com/findings",
		"eu-de":    "https://eu.compliance.cloud.ibm.com/si/findings",
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}

// Clone makes a copy of "findings" suitable for processing requests.
func (findings *FindingsV1) Clone() *FindingsV1 {
	if core.IsNil(findings) {
		return nil
	}
	clone := *findings
	clone.Service = findings.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (findings *FindingsV1) SetServiceURL(url string) error {
	return findings.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (findings *FindingsV1) GetServiceURL() string {
	return findings.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (findings *FindingsV1) SetDefaultHeaders(headers http.Header) {
	findings.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (findings *FindingsV1) SetEnableGzipCompression(enableGzip bool) {
	findings.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (findings *FindingsV1) GetEnableGzipCompression() bool {
	return findings.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (findings *FindingsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	findings.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (findings *FindingsV1) DisableRetries() {
	findings.Service.DisableRetries()
}

// PostGraph : Query findings
// Query findings by using the GraphQL query language. For more information about using GraphQL, see the [GraphQL
// documentation](https://graphql.org/learn/).
func (findings *FindingsV1) PostGraph(postGraphOptions *PostGraphOptions) (response *core.DetailedResponse, err error) {
	return findings.PostGraphWithContext(context.Background(), postGraphOptions)
}

// PostGraphWithContext is an alternate form of the PostGraph method which supports a Context parameter
func (findings *FindingsV1) PostGraphWithContext(ctx context.Context, postGraphOptions *PostGraphOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postGraphOptions, "postGraphOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postGraphOptions, "postGraphOptions")
	if err != nil {
		return
	}

	if postGraphOptions.Body != nil && postGraphOptions.ContentType == nil {
		postGraphOptions.SetContentType("application/json")
	}

	pathParamsMap := map[string]string{
		"account_id": *findings.AccountID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/graph`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postGraphOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "PostGraph")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if postGraphOptions.ContentType != nil {
		builder.AddHeader("Content-Type", fmt.Sprint(*postGraphOptions.ContentType))
	}
	if postGraphOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*postGraphOptions.TransactionID))
	}

	_, err = builder.SetBodyContent(core.StringNilMapper(postGraphOptions.ContentType), nil, nil, postGraphOptions.Body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = findings.Service.Request(request, nil)

	return
}

// ListProviders : List providers
// List all of the providers for a specified account.
func (findings *FindingsV1) ListProviders(listProvidersOptions *ListProvidersOptions) (result *APIListProvidersResponse, response *core.DetailedResponse, err error) {
	return findings.ListProvidersWithContext(context.Background(), listProvidersOptions)
}

// ListProvidersWithContext is an alternate form of the ListProviders method which supports a Context parameter
func (findings *FindingsV1) ListProvidersWithContext(ctx context.Context, listProvidersOptions *ListProvidersOptions) (result *APIListProvidersResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProvidersOptions, "listProvidersOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *findings.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProvidersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "ListProviders")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listProvidersOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listProvidersOptions.TransactionID))
	}

	if listProvidersOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProvidersOptions.Limit))
	}
	if listProvidersOptions.Skip != nil {
		builder.AddQuery("skip", fmt.Sprint(*listProvidersOptions.Skip))
	}
	if listProvidersOptions.StartProviderID != nil {
		builder.AddQuery("start_provider_id", fmt.Sprint(*listProvidersOptions.StartProviderID))
	}
	if listProvidersOptions.EndProviderID != nil {
		builder.AddQuery("end_provider_id", fmt.Sprint(*listProvidersOptions.EndProviderID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIListProvidersResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateNote : Create a note
// Register a new finding type with the Security and Compliance Center.
//
// A successful request creates a note with a high-level description of a particular type of finding. To learn more
// about creating notes to register findings, see [Custom
// findings](/docs/security-advisor?topic=security-advisor-setup_custom).
func (findings *FindingsV1) CreateNote(createNoteOptions *CreateNoteOptions) (result *APINote, response *core.DetailedResponse, err error) {
	return findings.CreateNoteWithContext(context.Background(), createNoteOptions)
}

// CreateNoteWithContext is an alternate form of the CreateNote method which supports a Context parameter
func (findings *FindingsV1) CreateNoteWithContext(ctx context.Context, createNoteOptions *CreateNoteOptions) (result *APINote, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createNoteOptions, "createNoteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createNoteOptions, "createNoteOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":  *findings.AccountID,
		"provider_id": *createNoteOptions.ProviderID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/notes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createNoteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "CreateNote")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createNoteOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createNoteOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if createNoteOptions.ShortDescription != nil {
		body["short_description"] = createNoteOptions.ShortDescription
	}
	if createNoteOptions.LongDescription != nil {
		body["long_description"] = createNoteOptions.LongDescription
	}
	if createNoteOptions.Kind != nil {
		body["kind"] = createNoteOptions.Kind
	}
	if createNoteOptions.ID != nil {
		body["id"] = createNoteOptions.ID
	}
	if createNoteOptions.ReportedBy != nil {
		body["reported_by"] = createNoteOptions.ReportedBy
	}
	if createNoteOptions.RelatedURL != nil {
		body["related_url"] = createNoteOptions.RelatedURL
	}
	if createNoteOptions.CreateTime != nil {
		body["create_time"] = createNoteOptions.CreateTime
	}
	if createNoteOptions.UpdateTime != nil {
		body["update_time"] = createNoteOptions.UpdateTime
	}
	if createNoteOptions.Shared != nil {
		body["shared"] = createNoteOptions.Shared
	}
	if createNoteOptions.Finding != nil {
		body["finding"] = createNoteOptions.Finding
	}
	if createNoteOptions.Kpi != nil {
		body["kpi"] = createNoteOptions.Kpi
	}
	if createNoteOptions.Card != nil {
		body["card"] = createNoteOptions.Card
	}
	if createNoteOptions.Section != nil {
		body["section"] = createNoteOptions.Section
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
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPINote)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListNotes : List notes
// List all of the available notes for a specific provider.
func (findings *FindingsV1) ListNotes(listNotesOptions *ListNotesOptions) (result *APIListNotesResponse, response *core.DetailedResponse, err error) {
	return findings.ListNotesWithContext(context.Background(), listNotesOptions)
}

// ListNotesWithContext is an alternate form of the ListNotes method which supports a Context parameter
func (findings *FindingsV1) ListNotesWithContext(ctx context.Context, listNotesOptions *ListNotesOptions) (result *APIListNotesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listNotesOptions, "listNotesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listNotesOptions, "listNotesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":  *findings.AccountID,
		"provider_id": *listNotesOptions.ProviderID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/notes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listNotesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "ListNotes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listNotesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listNotesOptions.TransactionID))
	}

	if listNotesOptions.PageSize != nil {
		builder.AddQuery("page_size", fmt.Sprint(*listNotesOptions.PageSize))
	}
	if listNotesOptions.PageToken != nil {
		builder.AddQuery("page_token", fmt.Sprint(*listNotesOptions.PageToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIListNotesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetNote : Get a note by provider
// Get the details of the note that is associated with a specified note ID and provider ID.
func (findings *FindingsV1) GetNote(getNoteOptions *GetNoteOptions) (result *APINote, response *core.DetailedResponse, err error) {
	return findings.GetNoteWithContext(context.Background(), getNoteOptions)
}

// GetNoteWithContext is an alternate form of the GetNote method which supports a Context parameter
func (findings *FindingsV1) GetNoteWithContext(ctx context.Context, getNoteOptions *GetNoteOptions) (result *APINote, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getNoteOptions, "getNoteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getNoteOptions, "getNoteOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":  *findings.AccountID,
		"provider_id": *getNoteOptions.ProviderID,
		"note_id":     *getNoteOptions.NoteID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/notes/{note_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getNoteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "GetNote")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getNoteOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getNoteOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPINote)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateNote : Update a note
// Update a note that already exists in your account.
func (findings *FindingsV1) UpdateNote(updateNoteOptions *UpdateNoteOptions) (result *APINote, response *core.DetailedResponse, err error) {
	return findings.UpdateNoteWithContext(context.Background(), updateNoteOptions)
}

// UpdateNoteWithContext is an alternate form of the UpdateNote method which supports a Context parameter
func (findings *FindingsV1) UpdateNoteWithContext(ctx context.Context, updateNoteOptions *UpdateNoteOptions) (result *APINote, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateNoteOptions, "updateNoteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateNoteOptions, "updateNoteOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":  *findings.AccountID,
		"provider_id": *updateNoteOptions.ProviderID,
		"note_id":     *updateNoteOptions.NoteID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/notes/{note_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateNoteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "UpdateNote")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateNoteOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateNoteOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if updateNoteOptions.ShortDescription != nil {
		body["short_description"] = updateNoteOptions.ShortDescription
	}
	if updateNoteOptions.LongDescription != nil {
		body["long_description"] = updateNoteOptions.LongDescription
	}
	if updateNoteOptions.Kind != nil {
		body["kind"] = updateNoteOptions.Kind
	}
	if updateNoteOptions.ID != nil {
		body["id"] = updateNoteOptions.ID
	}
	if updateNoteOptions.ReportedBy != nil {
		body["reported_by"] = updateNoteOptions.ReportedBy
	}
	if updateNoteOptions.RelatedURL != nil {
		body["related_url"] = updateNoteOptions.RelatedURL
	}
	if updateNoteOptions.CreateTime != nil {
		body["create_time"] = updateNoteOptions.CreateTime
	}
	if updateNoteOptions.UpdateTime != nil {
		body["update_time"] = updateNoteOptions.UpdateTime
	}
	if updateNoteOptions.Shared != nil {
		body["shared"] = updateNoteOptions.Shared
	}
	if updateNoteOptions.Finding != nil {
		body["finding"] = updateNoteOptions.Finding
	}
	if updateNoteOptions.Kpi != nil {
		body["kpi"] = updateNoteOptions.Kpi
	}
	if updateNoteOptions.Card != nil {
		body["card"] = updateNoteOptions.Card
	}
	if updateNoteOptions.Section != nil {
		body["section"] = updateNoteOptions.Section
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
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPINote)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteNote : Delete a note
// Delete a note with the ID and provider ID that you specify.
func (findings *FindingsV1) DeleteNote(deleteNoteOptions *DeleteNoteOptions) (response *core.DetailedResponse, err error) {
	return findings.DeleteNoteWithContext(context.Background(), deleteNoteOptions)
}

// DeleteNoteWithContext is an alternate form of the DeleteNote method which supports a Context parameter
func (findings *FindingsV1) DeleteNoteWithContext(ctx context.Context, deleteNoteOptions *DeleteNoteOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteNoteOptions, "deleteNoteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteNoteOptions, "deleteNoteOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":  *findings.AccountID,
		"provider_id": *deleteNoteOptions.ProviderID,
		"note_id":     *deleteNoteOptions.NoteID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/notes/{note_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteNoteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "DeleteNote")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteNoteOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteNoteOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = findings.Service.Request(request, nil)

	return
}

// GetOccurrenceNote : Get a note by occurrence
// Get a note that is associated with the occurrence ID that you specify.
func (findings *FindingsV1) GetOccurrenceNote(getOccurrenceNoteOptions *GetOccurrenceNoteOptions) (result *APINote, response *core.DetailedResponse, err error) {
	return findings.GetOccurrenceNoteWithContext(context.Background(), getOccurrenceNoteOptions)
}

// GetOccurrenceNoteWithContext is an alternate form of the GetOccurrenceNote method which supports a Context parameter
func (findings *FindingsV1) GetOccurrenceNoteWithContext(ctx context.Context, getOccurrenceNoteOptions *GetOccurrenceNoteOptions) (result *APINote, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOccurrenceNoteOptions, "getOccurrenceNoteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOccurrenceNoteOptions, "getOccurrenceNoteOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":    *findings.AccountID,
		"provider_id":   *getOccurrenceNoteOptions.ProviderID,
		"occurrence_id": *getOccurrenceNoteOptions.OccurrenceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/occurrences/{occurrence_id}/note`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOccurrenceNoteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "GetOccurrenceNote")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getOccurrenceNoteOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getOccurrenceNoteOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPINote)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateOccurrence : Create an occurrence
// Create an occurrence to denote the existence of a particular type of finding.
//
// An occurrence describes provider-specific details of a note and contains vulnerability details, remediation steps,
// and other general information.
func (findings *FindingsV1) CreateOccurrence(createOccurrenceOptions *CreateOccurrenceOptions) (result *APIOccurrence, response *core.DetailedResponse, err error) {
	return findings.CreateOccurrenceWithContext(context.Background(), createOccurrenceOptions)
}

// CreateOccurrenceWithContext is an alternate form of the CreateOccurrence method which supports a Context parameter
func (findings *FindingsV1) CreateOccurrenceWithContext(ctx context.Context, createOccurrenceOptions *CreateOccurrenceOptions) (result *APIOccurrence, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createOccurrenceOptions, "createOccurrenceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createOccurrenceOptions, "createOccurrenceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":  *findings.AccountID,
		"provider_id": *createOccurrenceOptions.ProviderID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/occurrences`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createOccurrenceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "CreateOccurrence")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createOccurrenceOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createOccurrenceOptions.TransactionID))
	}
	if createOccurrenceOptions.ReplaceIfExists != nil {
		builder.AddHeader("Replace-If-Exists", fmt.Sprint(*createOccurrenceOptions.ReplaceIfExists))
	}

	body := make(map[string]interface{})
	if createOccurrenceOptions.NoteName != nil {
		body["note_name"] = createOccurrenceOptions.NoteName
	}
	if createOccurrenceOptions.Kind != nil {
		body["kind"] = createOccurrenceOptions.Kind
	}
	if createOccurrenceOptions.ID != nil {
		body["id"] = createOccurrenceOptions.ID
	}
	if createOccurrenceOptions.ResourceURL != nil {
		body["resource_url"] = createOccurrenceOptions.ResourceURL
	}
	if createOccurrenceOptions.Remediation != nil {
		body["remediation"] = createOccurrenceOptions.Remediation
	}
	if createOccurrenceOptions.CreateTime != nil {
		body["create_time"] = createOccurrenceOptions.CreateTime
	}
	if createOccurrenceOptions.UpdateTime != nil {
		body["update_time"] = createOccurrenceOptions.UpdateTime
	}
	if createOccurrenceOptions.Context != nil {
		body["context"] = createOccurrenceOptions.Context
	}
	if createOccurrenceOptions.Finding != nil {
		body["finding"] = createOccurrenceOptions.Finding
	}
	if createOccurrenceOptions.Kpi != nil {
		body["kpi"] = createOccurrenceOptions.Kpi
	}
	if createOccurrenceOptions.ReferenceData != nil {
		body["reference_data"] = createOccurrenceOptions.ReferenceData
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
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIOccurrence)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListOccurrences : List occurrences
// List all of the occurrences that are associated with the provider ID that you specify.
func (findings *FindingsV1) ListOccurrences(listOccurrencesOptions *ListOccurrencesOptions) (result *APIListOccurrencesResponse, response *core.DetailedResponse, err error) {
	return findings.ListOccurrencesWithContext(context.Background(), listOccurrencesOptions)
}

// ListOccurrencesWithContext is an alternate form of the ListOccurrences method which supports a Context parameter
func (findings *FindingsV1) ListOccurrencesWithContext(ctx context.Context, listOccurrencesOptions *ListOccurrencesOptions) (result *APIListOccurrencesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listOccurrencesOptions, "listOccurrencesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listOccurrencesOptions, "listOccurrencesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":  *findings.AccountID,
		"provider_id": *listOccurrencesOptions.ProviderID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/occurrences`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listOccurrencesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "ListOccurrences")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listOccurrencesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listOccurrencesOptions.TransactionID))
	}

	if listOccurrencesOptions.PageSize != nil {
		builder.AddQuery("page_size", fmt.Sprint(*listOccurrencesOptions.PageSize))
	}
	if listOccurrencesOptions.PageToken != nil {
		builder.AddQuery("page_token", fmt.Sprint(*listOccurrencesOptions.PageToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIListOccurrencesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListNoteOccurrences : List occurrences by note
// Get a list of occurrences that are associated with a specific note.
func (findings *FindingsV1) ListNoteOccurrences(listNoteOccurrencesOptions *ListNoteOccurrencesOptions) (result *APIListNoteOccurrencesResponse, response *core.DetailedResponse, err error) {
	return findings.ListNoteOccurrencesWithContext(context.Background(), listNoteOccurrencesOptions)
}

// ListNoteOccurrencesWithContext is an alternate form of the ListNoteOccurrences method which supports a Context parameter
func (findings *FindingsV1) ListNoteOccurrencesWithContext(ctx context.Context, listNoteOccurrencesOptions *ListNoteOccurrencesOptions) (result *APIListNoteOccurrencesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listNoteOccurrencesOptions, "listNoteOccurrencesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listNoteOccurrencesOptions, "listNoteOccurrencesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":  *findings.AccountID,
		"provider_id": *listNoteOccurrencesOptions.ProviderID,
		"note_id":     *listNoteOccurrencesOptions.NoteID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/notes/{note_id}/occurrences`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listNoteOccurrencesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "ListNoteOccurrences")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listNoteOccurrencesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listNoteOccurrencesOptions.TransactionID))
	}

	if listNoteOccurrencesOptions.PageSize != nil {
		builder.AddQuery("page_size", fmt.Sprint(*listNoteOccurrencesOptions.PageSize))
	}
	if listNoteOccurrencesOptions.PageToken != nil {
		builder.AddQuery("page_token", fmt.Sprint(*listNoteOccurrencesOptions.PageToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIListNoteOccurrencesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOccurrence : Get a specific occurrence
// Get the details of a specific occurrence by specifying the ID and provider ID.
func (findings *FindingsV1) GetOccurrence(getOccurrenceOptions *GetOccurrenceOptions) (result *APIOccurrence, response *core.DetailedResponse, err error) {
	return findings.GetOccurrenceWithContext(context.Background(), getOccurrenceOptions)
}

// GetOccurrenceWithContext is an alternate form of the GetOccurrence method which supports a Context parameter
func (findings *FindingsV1) GetOccurrenceWithContext(ctx context.Context, getOccurrenceOptions *GetOccurrenceOptions) (result *APIOccurrence, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOccurrenceOptions, "getOccurrenceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOccurrenceOptions, "getOccurrenceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":    *findings.AccountID,
		"provider_id":   *getOccurrenceOptions.ProviderID,
		"occurrence_id": *getOccurrenceOptions.OccurrenceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/occurrences/{occurrence_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOccurrenceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "GetOccurrence")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getOccurrenceOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getOccurrenceOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIOccurrence)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateOccurrence : Update an occurrence
// Update an occurrence that already exists in your account.
func (findings *FindingsV1) UpdateOccurrence(updateOccurrenceOptions *UpdateOccurrenceOptions) (result *APIOccurrence, response *core.DetailedResponse, err error) {
	return findings.UpdateOccurrenceWithContext(context.Background(), updateOccurrenceOptions)
}

// UpdateOccurrenceWithContext is an alternate form of the UpdateOccurrence method which supports a Context parameter
func (findings *FindingsV1) UpdateOccurrenceWithContext(ctx context.Context, updateOccurrenceOptions *UpdateOccurrenceOptions) (result *APIOccurrence, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateOccurrenceOptions, "updateOccurrenceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateOccurrenceOptions, "updateOccurrenceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":    *findings.AccountID,
		"provider_id":   *updateOccurrenceOptions.ProviderID,
		"occurrence_id": *updateOccurrenceOptions.OccurrenceID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/occurrences/{occurrence_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateOccurrenceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "UpdateOccurrence")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateOccurrenceOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateOccurrenceOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if updateOccurrenceOptions.NoteName != nil {
		body["note_name"] = updateOccurrenceOptions.NoteName
	}
	if updateOccurrenceOptions.Kind != nil {
		body["kind"] = updateOccurrenceOptions.Kind
	}
	if updateOccurrenceOptions.ID != nil {
		body["id"] = updateOccurrenceOptions.ID
	}
	if updateOccurrenceOptions.ResourceURL != nil {
		body["resource_url"] = updateOccurrenceOptions.ResourceURL
	}
	if updateOccurrenceOptions.Remediation != nil {
		body["remediation"] = updateOccurrenceOptions.Remediation
	}
	if updateOccurrenceOptions.CreateTime != nil {
		body["create_time"] = updateOccurrenceOptions.CreateTime
	}
	if updateOccurrenceOptions.UpdateTime != nil {
		body["update_time"] = updateOccurrenceOptions.UpdateTime
	}
	if updateOccurrenceOptions.Context != nil {
		body["context"] = updateOccurrenceOptions.Context
	}
	if updateOccurrenceOptions.Finding != nil {
		body["finding"] = updateOccurrenceOptions.Finding
	}
	if updateOccurrenceOptions.Kpi != nil {
		body["kpi"] = updateOccurrenceOptions.Kpi
	}
	if updateOccurrenceOptions.ReferenceData != nil {
		body["reference_data"] = updateOccurrenceOptions.ReferenceData
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
	response, err = findings.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAPIOccurrence)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteOccurrence : Delete an occurrence
// Delete an occurrence by specifying the occurrence ID and provider ID.
func (findings *FindingsV1) DeleteOccurrence(deleteOccurrenceOptions *DeleteOccurrenceOptions) (response *core.DetailedResponse, err error) {
	return findings.DeleteOccurrenceWithContext(context.Background(), deleteOccurrenceOptions)
}

// DeleteOccurrenceWithContext is an alternate form of the DeleteOccurrence method which supports a Context parameter
func (findings *FindingsV1) DeleteOccurrenceWithContext(ctx context.Context, deleteOccurrenceOptions *DeleteOccurrenceOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteOccurrenceOptions, "deleteOccurrenceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteOccurrenceOptions, "deleteOccurrenceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"account_id":    *findings.AccountID,
		"provider_id":   *deleteOccurrenceOptions.ProviderID,
		"occurrence_id": *deleteOccurrenceOptions.OccurrenceID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = findings.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(findings.Service.Options.URL, `/v1/{account_id}/providers/{provider_id}/occurrences/{occurrence_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteOccurrenceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("findings", "V1", "DeleteOccurrence")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteOccurrenceOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteOccurrenceOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = findings.Service.Request(request, nil)

	return
}

// Card : Card provides details about a card kind of note.
type Card struct {
	// The section this card belongs to.
	Section *string `json:"section" validate:"required"`

	// The title of this card.
	Title *string `json:"title" validate:"required"`

	// The subtitle of this card.
	Subtitle *string `json:"subtitle" validate:"required"`

	// The order of the card in which it will appear on SA dashboard in the mentioned section.
	Order *int64 `json:"order,omitempty"`

	// The finding note names associated to this card.
	FindingNoteNames []string `json:"finding_note_names" validate:"required"`

	RequiresConfiguration *bool `json:"requires_configuration,omitempty"`

	// The text associated to the card's badge.
	BadgeText *string `json:"badge_text,omitempty"`

	// The base64 content of the image associated to the card's badge.
	BadgeImage *string `json:"badge_image,omitempty"`

	// The elements of this card.
	Elements []CardElementIntf `json:"elements" validate:"required"`
}

// NewCard : Instantiate Card (Generic Model Constructor)
func (*FindingsV1) NewCard(section string, title string, subtitle string, findingNoteNames []string, elements []CardElementIntf) (_model *Card, err error) {
	_model = &Card{
		Section:          core.StringPtr(section),
		Title:            core.StringPtr(title),
		Subtitle:         core.StringPtr(subtitle),
		FindingNoteNames: findingNoteNames,
		Elements:         elements,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalCard unmarshals an instance of Card from the specified map of raw messages.
func UnmarshalCard(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Card)
	err = core.UnmarshalPrimitive(m, "section", &obj.Section)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "title", &obj.Title)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "subtitle", &obj.Subtitle)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "order", &obj.Order)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "finding_note_names", &obj.FindingNoteNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "requires_configuration", &obj.RequiresConfiguration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "badge_text", &obj.BadgeText)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "badge_image", &obj.BadgeImage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "elements", &obj.Elements, UnmarshalCardElement)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CardElement : CardElement provides details about the elements of a Card.
// Models which "extend" this model:
// - CardElementNumericCardElement
// - CardElementBreakdownCardElement
// - CardElementTimeSeriesCardElement
type CardElement struct {
	// The text of this card element.
	Text *string `json:"text,omitempty"`

	// Kind of element
	// - NUMERIC&#58; Single numeric value
	// - BREAKDOWN&#58; Breakdown of numeric values
	// - TIME_SERIES&#58; Time-series of numeric values.
	Kind *string `json:"kind,omitempty"`

	// The default time range of this card element.
	DefaultTimeRange *string `json:"default_time_range,omitempty"`

	ValueType *NumericCardElementValueType `json:"value_type,omitempty"`

	// the value types associated to this card element.
	ValueTypes []ValueTypeIntf `json:"value_types,omitempty"`

	// The default interval of the time series.
	DefaultInterval *string `json:"default_interval,omitempty"`
}

// Constants associated with the CardElement.Kind property.
// Kind of element
// - NUMERIC&#58; Single numeric value
// - BREAKDOWN&#58; Breakdown of numeric values
// - TIME_SERIES&#58; Time-series of numeric values.
const (
	CardElementKindBreakdownConst  = "BREAKDOWN"
	CardElementKindNumericConst    = "NUMERIC"
	CardElementKindTimeSeriesConst = "TIME_SERIES"
)

func (*CardElement) isaCardElement() bool {
	return true
}

type CardElementIntf interface {
	isaCardElement() bool
}

// UnmarshalCardElement unmarshals an instance of CardElement from the specified map of raw messages.
func UnmarshalCardElement(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "kind", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'kind': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'kind' not found in JSON object")
		return
	}
	if discValue == "NUMERIC" {
		err = core.UnmarshalModel(m, "", result, UnmarshalCardElementNumericCardElement)
	} else if discValue == "BREAKDOWN" {
		err = core.UnmarshalModel(m, "", result, UnmarshalCardElementBreakdownCardElement)
	} else if discValue == "TIME_SERIES" {
		err = core.UnmarshalModel(m, "", result, UnmarshalCardElementTimeSeriesCardElement)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'kind': %s", discValue)
	}
	return
}

// Context : Context struct
type Context struct {
	// The IBM Cloud region.
	Region *string `json:"region,omitempty"`

	// The resource CRN (e.g. certificate CRN, image CRN).
	ResourceCRN *string `json:"resource_crn,omitempty"`

	// The resource ID, in case the CRN is not available.
	ResourceID *string `json:"resource_id,omitempty"`

	// The user-friendly resource name.
	ResourceName *string `json:"resource_name,omitempty"`

	// The resource type name (e.g. Pod, Cluster, Certificate, Image).
	ResourceType *string `json:"resource_type,omitempty"`

	// The service CRN (e.g. CertMgr Instance CRN).
	ServiceCRN *string `json:"service_crn,omitempty"`

	// The service name (e.g. CertMgr).
	ServiceName *string `json:"service_name,omitempty"`

	// The name of the environment the occurrence applies to.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// The name of the component the occurrence applies to.
	ComponentName *string `json:"component_name,omitempty"`

	// The id of the toolchain the occurrence applies to.
	ToolchainID *string `json:"toolchain_id,omitempty"`
}

// UnmarshalContext unmarshals an instance of Context from the specified map of raw messages.
func UnmarshalContext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Context)
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crn", &obj.ResourceCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_id", &obj.ResourceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_name", &obj.ComponentName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_id", &obj.ToolchainID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateNoteOptions : The CreateNote options.
type CreateNoteOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// A one sentence description of your note.
	ShortDescription *string `json:"short_description" validate:"required"`

	// A more detailed description of your note.
	LongDescription *string `json:"long_description" validate:"required"`

	// The type of note. Use this field to filter notes and occurences by kind.
	//  - FINDING&#58; The note and occurrence represent a finding.
	//  - KPI&#58; The note and occurrence represent a KPI value.
	//  - CARD&#58; The note represents a card showing findings and related metric values.
	//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
	//  - SECTION&#58; The note represents a section in a dashboard.
	Kind *string `json:"kind" validate:"required"`

	// The ID of the note.
	ID *string `json:"id" validate:"required"`

	// The entity reporting a note.
	ReportedBy *Reporter `json:"reported_by" validate:"required"`

	RelatedURL []APINoteRelatedURL `json:"related_url,omitempty"`

	// Output only. The time this note was created. This field can be used as a filter in list requests.
	CreateTime *strfmt.DateTime `json:"create_time,omitempty"`

	// Output only. The time this note was last updated. This field can be used as a filter in list requests.
	UpdateTime *strfmt.DateTime `json:"update_time,omitempty"`

	// True if this note can be shared by multiple accounts.
	Shared *bool `json:"shared,omitempty"`

	// FindingType provides details about a finding note.
	Finding *FindingType `json:"finding,omitempty"`

	// KpiType provides details about a KPI note.
	Kpi *KpiType `json:"kpi,omitempty"`

	// Card provides details about a card kind of note.
	Card *Card `json:"card,omitempty"`

	// Card provides details about a card kind of note.
	Section *Section `json:"section,omitempty"`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateNoteOptions.Kind property.
// The type of note. Use this field to filter notes and occurences by kind.
//  - FINDING&#58; The note and occurrence represent a finding.
//  - KPI&#58; The note and occurrence represent a KPI value.
//  - CARD&#58; The note represents a card showing findings and related metric values.
//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
//  - SECTION&#58; The note represents a section in a dashboard.
const (
	CreateNoteOptionsKindCardConst           = "CARD"
	CreateNoteOptionsKindCardConfiguredConst = "CARD_CONFIGURED"
	CreateNoteOptionsKindFindingConst        = "FINDING"
	CreateNoteOptionsKindKpiConst            = "KPI"
	CreateNoteOptionsKindSectionConst        = "SECTION"
)

// NewCreateNoteOptions : Instantiate CreateNoteOptions
func (*FindingsV1) NewCreateNoteOptions(providerID string, shortDescription string, longDescription string, kind string, id string, reportedBy *Reporter) *CreateNoteOptions {
	return &CreateNoteOptions{
		ProviderID:       core.StringPtr(providerID),
		ShortDescription: core.StringPtr(shortDescription),
		LongDescription:  core.StringPtr(longDescription),
		Kind:             core.StringPtr(kind),
		ID:               core.StringPtr(id),
		ReportedBy:       reportedBy,
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *CreateNoteOptions) SetProviderID(providerID string) *CreateNoteOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetShortDescription : Allow user to set ShortDescription
func (_options *CreateNoteOptions) SetShortDescription(shortDescription string) *CreateNoteOptions {
	_options.ShortDescription = core.StringPtr(shortDescription)
	return _options
}

// SetLongDescription : Allow user to set LongDescription
func (_options *CreateNoteOptions) SetLongDescription(longDescription string) *CreateNoteOptions {
	_options.LongDescription = core.StringPtr(longDescription)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *CreateNoteOptions) SetKind(kind string) *CreateNoteOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateNoteOptions) SetID(id string) *CreateNoteOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetReportedBy : Allow user to set ReportedBy
func (_options *CreateNoteOptions) SetReportedBy(reportedBy *Reporter) *CreateNoteOptions {
	_options.ReportedBy = reportedBy
	return _options
}

// SetRelatedURL : Allow user to set RelatedURL
func (_options *CreateNoteOptions) SetRelatedURL(relatedURL []APINoteRelatedURL) *CreateNoteOptions {
	_options.RelatedURL = relatedURL
	return _options
}

// SetCreateTime : Allow user to set CreateTime
func (_options *CreateNoteOptions) SetCreateTime(createTime *strfmt.DateTime) *CreateNoteOptions {
	_options.CreateTime = createTime
	return _options
}

// SetUpdateTime : Allow user to set UpdateTime
func (_options *CreateNoteOptions) SetUpdateTime(updateTime *strfmt.DateTime) *CreateNoteOptions {
	_options.UpdateTime = updateTime
	return _options
}

// SetShared : Allow user to set Shared
func (_options *CreateNoteOptions) SetShared(shared bool) *CreateNoteOptions {
	_options.Shared = core.BoolPtr(shared)
	return _options
}

// SetFinding : Allow user to set Finding
func (_options *CreateNoteOptions) SetFinding(finding *FindingType) *CreateNoteOptions {
	_options.Finding = finding
	return _options
}

// SetKpi : Allow user to set Kpi
func (_options *CreateNoteOptions) SetKpi(kpi *KpiType) *CreateNoteOptions {
	_options.Kpi = kpi
	return _options
}

// SetCard : Allow user to set Card
func (_options *CreateNoteOptions) SetCard(card *Card) *CreateNoteOptions {
	_options.Card = card
	return _options
}

// SetSection : Allow user to set Section
func (_options *CreateNoteOptions) SetSection(section *Section) *CreateNoteOptions {
	_options.Section = section
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateNoteOptions) SetTransactionID(transactionID string) *CreateNoteOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateNoteOptions) SetHeaders(param map[string]string) *CreateNoteOptions {
	options.Headers = param
	return options
}

// CreateOccurrenceOptions : The CreateOccurrence options.
type CreateOccurrenceOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// An analysis note associated with this image, in the form "{account_id}/providers/{provider_id}/notes/{note_id}" This
	// field can be used as a filter in list requests.
	NoteName *string `json:"note_name" validate:"required"`

	// The type of note. Use this field to filter notes and occurences by kind.
	//  - FINDING&#58; The note and occurrence represent a finding.
	//  - KPI&#58; The note and occurrence represent a KPI value.
	//  - CARD&#58; The note represents a card showing findings and related metric values.
	//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
	//  - SECTION&#58; The note represents a section in a dashboard.
	Kind *string `json:"kind" validate:"required"`

	// The id of the occurrence.
	ID *string `json:"id" validate:"required"`

	// The unique URL of the resource, image or the container, for which the `Occurrence` applies. For example,
	// https://gcr.io/provider/image@sha256:foo. This field can be used as a filter in list requests.
	ResourceURL *string `json:"resource_url,omitempty"`

	// A description of actions that can be taken to remedy the `Note`.
	Remediation *string `json:"remediation,omitempty"`

	// Output only. The time this `Occurrence` was created.
	CreateTime *strfmt.DateTime `json:"create_time,omitempty"`

	// Output only. The time this `Occurrence` was last updated.
	UpdateTime *strfmt.DateTime `json:"update_time,omitempty"`

	Context *Context `json:"context,omitempty"`

	// Finding provides details about a finding occurrence.
	Finding *Finding `json:"finding,omitempty"`

	// Kpi provides details about a KPI occurrence.
	Kpi *Kpi `json:"kpi,omitempty"`

	// Additional data for the finding, like AT event etc.
	ReferenceData interface{} `json:"reference_data,omitempty"`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// When set to true, an existing occurrence is replaced rather than duplicated.
	ReplaceIfExists *bool `json:"Replace-If-Exists,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateOccurrenceOptions.Kind property.
// The type of note. Use this field to filter notes and occurences by kind.
//  - FINDING&#58; The note and occurrence represent a finding.
//  - KPI&#58; The note and occurrence represent a KPI value.
//  - CARD&#58; The note represents a card showing findings and related metric values.
//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
//  - SECTION&#58; The note represents a section in a dashboard.
const (
	CreateOccurrenceOptionsKindCardConst           = "CARD"
	CreateOccurrenceOptionsKindCardConfiguredConst = "CARD_CONFIGURED"
	CreateOccurrenceOptionsKindFindingConst        = "FINDING"
	CreateOccurrenceOptionsKindKpiConst            = "KPI"
	CreateOccurrenceOptionsKindSectionConst        = "SECTION"
)

// NewCreateOccurrenceOptions : Instantiate CreateOccurrenceOptions
func (*FindingsV1) NewCreateOccurrenceOptions(providerID string, noteName string, kind string, id string) *CreateOccurrenceOptions {
	return &CreateOccurrenceOptions{
		ProviderID: core.StringPtr(providerID),
		NoteName:   core.StringPtr(noteName),
		Kind:       core.StringPtr(kind),
		ID:         core.StringPtr(id),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *CreateOccurrenceOptions) SetProviderID(providerID string) *CreateOccurrenceOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetNoteName : Allow user to set NoteName
func (_options *CreateOccurrenceOptions) SetNoteName(noteName string) *CreateOccurrenceOptions {
	_options.NoteName = core.StringPtr(noteName)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *CreateOccurrenceOptions) SetKind(kind string) *CreateOccurrenceOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateOccurrenceOptions) SetID(id string) *CreateOccurrenceOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetResourceURL : Allow user to set ResourceURL
func (_options *CreateOccurrenceOptions) SetResourceURL(resourceURL string) *CreateOccurrenceOptions {
	_options.ResourceURL = core.StringPtr(resourceURL)
	return _options
}

// SetRemediation : Allow user to set Remediation
func (_options *CreateOccurrenceOptions) SetRemediation(remediation string) *CreateOccurrenceOptions {
	_options.Remediation = core.StringPtr(remediation)
	return _options
}

// SetCreateTime : Allow user to set CreateTime
func (_options *CreateOccurrenceOptions) SetCreateTime(createTime *strfmt.DateTime) *CreateOccurrenceOptions {
	_options.CreateTime = createTime
	return _options
}

// SetUpdateTime : Allow user to set UpdateTime
func (_options *CreateOccurrenceOptions) SetUpdateTime(updateTime *strfmt.DateTime) *CreateOccurrenceOptions {
	_options.UpdateTime = updateTime
	return _options
}

// SetContext : Allow user to set Context
func (_options *CreateOccurrenceOptions) SetContext(context *Context) *CreateOccurrenceOptions {
	_options.Context = context
	return _options
}

// SetFinding : Allow user to set Finding
func (_options *CreateOccurrenceOptions) SetFinding(finding *Finding) *CreateOccurrenceOptions {
	_options.Finding = finding
	return _options
}

// SetKpi : Allow user to set Kpi
func (_options *CreateOccurrenceOptions) SetKpi(kpi *Kpi) *CreateOccurrenceOptions {
	_options.Kpi = kpi
	return _options
}

// SetReferenceData : Allow user to set ReferenceData
func (_options *CreateOccurrenceOptions) SetReferenceData(referenceData interface{}) *CreateOccurrenceOptions {
	_options.ReferenceData = referenceData
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateOccurrenceOptions) SetTransactionID(transactionID string) *CreateOccurrenceOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetReplaceIfExists : Allow user to set ReplaceIfExists
func (_options *CreateOccurrenceOptions) SetReplaceIfExists(replaceIfExists bool) *CreateOccurrenceOptions {
	_options.ReplaceIfExists = core.BoolPtr(replaceIfExists)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateOccurrenceOptions) SetHeaders(param map[string]string) *CreateOccurrenceOptions {
	options.Headers = param
	return options
}

// DataTransferred : It provides details about data transferred between clients and servers.
type DataTransferred struct {
	// The number of client bytes transferred.
	ClientBytes *int64 `json:"client_bytes,omitempty"`

	// The number of server bytes transferred.
	ServerBytes *int64 `json:"server_bytes,omitempty"`

	// The number of client packets transferred.
	ClientPackets *int64 `json:"client_packets,omitempty"`

	// The number of server packets transferred.
	ServerPackets *int64 `json:"server_packets,omitempty"`
}

// UnmarshalDataTransferred unmarshals an instance of DataTransferred from the specified map of raw messages.
func UnmarshalDataTransferred(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DataTransferred)
	err = core.UnmarshalPrimitive(m, "client_bytes", &obj.ClientBytes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_bytes", &obj.ServerBytes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_packets", &obj.ClientPackets)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_packets", &obj.ServerPackets)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteNoteOptions : The DeleteNote options.
type DeleteNoteOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// Second part of note `name`: providers/{provider_id}/notes/{note_id}.
	NoteID *string `json:"note_id" validate:"required,ne="`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteNoteOptions : Instantiate DeleteNoteOptions
func (*FindingsV1) NewDeleteNoteOptions(providerID string, noteID string) *DeleteNoteOptions {
	return &DeleteNoteOptions{
		ProviderID: core.StringPtr(providerID),
		NoteID:     core.StringPtr(noteID),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *DeleteNoteOptions) SetProviderID(providerID string) *DeleteNoteOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetNoteID : Allow user to set NoteID
func (_options *DeleteNoteOptions) SetNoteID(noteID string) *DeleteNoteOptions {
	_options.NoteID = core.StringPtr(noteID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteNoteOptions) SetTransactionID(transactionID string) *DeleteNoteOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteNoteOptions) SetHeaders(param map[string]string) *DeleteNoteOptions {
	options.Headers = param
	return options
}

// DeleteOccurrenceOptions : The DeleteOccurrence options.
type DeleteOccurrenceOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// Second part of occurrence `name`: providers/{provider_id}/occurrences/{occurrence_id}.
	OccurrenceID *string `json:"occurrence_id" validate:"required,ne="`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteOccurrenceOptions : Instantiate DeleteOccurrenceOptions
func (*FindingsV1) NewDeleteOccurrenceOptions(providerID string, occurrenceID string) *DeleteOccurrenceOptions {
	return &DeleteOccurrenceOptions{
		ProviderID:   core.StringPtr(providerID),
		OccurrenceID: core.StringPtr(occurrenceID),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *DeleteOccurrenceOptions) SetProviderID(providerID string) *DeleteOccurrenceOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetOccurrenceID : Allow user to set OccurrenceID
func (_options *DeleteOccurrenceOptions) SetOccurrenceID(occurrenceID string) *DeleteOccurrenceOptions {
	_options.OccurrenceID = core.StringPtr(occurrenceID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteOccurrenceOptions) SetTransactionID(transactionID string) *DeleteOccurrenceOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOccurrenceOptions) SetHeaders(param map[string]string) *DeleteOccurrenceOptions {
	options.Headers = param
	return options
}

// Finding : Finding provides details about a finding occurrence.
type Finding struct {
	// Note provider-assigned severity/impact ranking
	// - LOW&#58; Low Impact
	// - MEDIUM&#58; Medium Impact
	// - HIGH&#58; High Impact
	// - CRITICAL&#58; Critical Impact.
	Severity *string `json:"severity,omitempty"`

	// Note provider-assigned confidence on the validity of an occurrence
	// - LOW&#58; Low Certainty
	// - MEDIUM&#58; Medium Certainty
	// - HIGH&#58; High Certainty.
	Certainty *string `json:"certainty,omitempty"`

	// Remediation steps for the issues reported in this finding. They override the note's next steps.
	NextSteps []RemediationStep `json:"next_steps,omitempty"`

	// It provides details about a network connection.
	NetworkConnection *NetworkConnection `json:"network_connection,omitempty"`

	// It provides details about data transferred between clients and servers.
	DataTransferred *DataTransferred `json:"data_transferred,omitempty"`
}

// Constants associated with the Finding.Severity property.
// Note provider-assigned severity/impact ranking
// - LOW&#58; Low Impact
// - MEDIUM&#58; Medium Impact
// - HIGH&#58; High Impact
// - CRITICAL&#58; Critical Impact.
const (
	FindingSeverityCriticalConst = "CRITICAL"
	FindingSeverityHighConst     = "HIGH"
	FindingSeverityLowConst      = "LOW"
	FindingSeverityMediumConst   = "MEDIUM"
)

// Constants associated with the Finding.Certainty property.
// Note provider-assigned confidence on the validity of an occurrence
// - LOW&#58; Low Certainty
// - MEDIUM&#58; Medium Certainty
// - HIGH&#58; High Certainty.
const (
	FindingCertaintyHighConst   = "HIGH"
	FindingCertaintyLowConst    = "LOW"
	FindingCertaintyMediumConst = "MEDIUM"
)

// UnmarshalFinding unmarshals an instance of Finding from the specified map of raw messages.
func UnmarshalFinding(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Finding)
	err = core.UnmarshalPrimitive(m, "severity", &obj.Severity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certainty", &obj.Certainty)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next_steps", &obj.NextSteps, UnmarshalRemediationStep)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "network_connection", &obj.NetworkConnection, UnmarshalNetworkConnection)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "data_transferred", &obj.DataTransferred, UnmarshalDataTransferred)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FindingType : FindingType provides details about a finding note.
type FindingType struct {
	// Note provider-assigned severity/impact ranking
	// - LOW&#58; Low Impact
	// - MEDIUM&#58; Medium Impact
	// - HIGH&#58; High Impact
	// - CRITICAL&#58; Critical Impact.
	Severity *string `json:"severity" validate:"required"`

	// Common remediation steps for the finding of this type.
	NextSteps []RemediationStep `json:"next_steps,omitempty"`
}

// Constants associated with the FindingType.Severity property.
// Note provider-assigned severity/impact ranking
// - LOW&#58; Low Impact
// - MEDIUM&#58; Medium Impact
// - HIGH&#58; High Impact
// - CRITICAL&#58; Critical Impact.
const (
	FindingTypeSeverityCriticalConst = "CRITICAL"
	FindingTypeSeverityHighConst     = "HIGH"
	FindingTypeSeverityLowConst      = "LOW"
	FindingTypeSeverityMediumConst   = "MEDIUM"
)

// NewFindingType : Instantiate FindingType (Generic Model Constructor)
func (*FindingsV1) NewFindingType(severity string) (_model *FindingType, err error) {
	_model = &FindingType{
		Severity: core.StringPtr(severity),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalFindingType unmarshals an instance of FindingType from the specified map of raw messages.
func UnmarshalFindingType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FindingType)
	err = core.UnmarshalPrimitive(m, "severity", &obj.Severity)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next_steps", &obj.NextSteps, UnmarshalRemediationStep)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetNoteOptions : The GetNote options.
type GetNoteOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// Second part of note `name`: providers/{provider_id}/notes/{note_id}.
	NoteID *string `json:"note_id" validate:"required,ne="`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetNoteOptions : Instantiate GetNoteOptions
func (*FindingsV1) NewGetNoteOptions(providerID string, noteID string) *GetNoteOptions {
	return &GetNoteOptions{
		ProviderID: core.StringPtr(providerID),
		NoteID:     core.StringPtr(noteID),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *GetNoteOptions) SetProviderID(providerID string) *GetNoteOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetNoteID : Allow user to set NoteID
func (_options *GetNoteOptions) SetNoteID(noteID string) *GetNoteOptions {
	_options.NoteID = core.StringPtr(noteID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetNoteOptions) SetTransactionID(transactionID string) *GetNoteOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetNoteOptions) SetHeaders(param map[string]string) *GetNoteOptions {
	options.Headers = param
	return options
}

// GetOccurrenceNoteOptions : The GetOccurrenceNote options.
type GetOccurrenceNoteOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// Second part of occurrence `name`: providers/{provider_id}/occurrences/{occurrence_id}.
	OccurrenceID *string `json:"occurrence_id" validate:"required,ne="`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOccurrenceNoteOptions : Instantiate GetOccurrenceNoteOptions
func (*FindingsV1) NewGetOccurrenceNoteOptions(providerID string, occurrenceID string) *GetOccurrenceNoteOptions {
	return &GetOccurrenceNoteOptions{
		ProviderID:   core.StringPtr(providerID),
		OccurrenceID: core.StringPtr(occurrenceID),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *GetOccurrenceNoteOptions) SetProviderID(providerID string) *GetOccurrenceNoteOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetOccurrenceID : Allow user to set OccurrenceID
func (_options *GetOccurrenceNoteOptions) SetOccurrenceID(occurrenceID string) *GetOccurrenceNoteOptions {
	_options.OccurrenceID = core.StringPtr(occurrenceID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetOccurrenceNoteOptions) SetTransactionID(transactionID string) *GetOccurrenceNoteOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOccurrenceNoteOptions) SetHeaders(param map[string]string) *GetOccurrenceNoteOptions {
	options.Headers = param
	return options
}

// GetOccurrenceOptions : The GetOccurrence options.
type GetOccurrenceOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// Second part of occurrence `name`: providers/{provider_id}/occurrences/{occurrence_id}.
	OccurrenceID *string `json:"occurrence_id" validate:"required,ne="`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOccurrenceOptions : Instantiate GetOccurrenceOptions
func (*FindingsV1) NewGetOccurrenceOptions(providerID string, occurrenceID string) *GetOccurrenceOptions {
	return &GetOccurrenceOptions{
		ProviderID:   core.StringPtr(providerID),
		OccurrenceID: core.StringPtr(occurrenceID),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *GetOccurrenceOptions) SetProviderID(providerID string) *GetOccurrenceOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetOccurrenceID : Allow user to set OccurrenceID
func (_options *GetOccurrenceOptions) SetOccurrenceID(occurrenceID string) *GetOccurrenceOptions {
	_options.OccurrenceID = core.StringPtr(occurrenceID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetOccurrenceOptions) SetTransactionID(transactionID string) *GetOccurrenceOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOccurrenceOptions) SetHeaders(param map[string]string) *GetOccurrenceOptions {
	options.Headers = param
	return options
}

// Kpi : Kpi provides details about a KPI occurrence.
type Kpi struct {
	// The value of this KPI.
	Value *float64 `json:"value" validate:"required"`

	// The total value of this KPI.
	Total *float64 `json:"total,omitempty"`
}

// NewKpi : Instantiate Kpi (Generic Model Constructor)
func (*FindingsV1) NewKpi(value float64) (_model *Kpi, err error) {
	_model = &Kpi{
		Value: core.Float64Ptr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalKpi unmarshals an instance of Kpi from the specified map of raw messages.
func UnmarshalKpi(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Kpi)
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total", &obj.Total)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KpiType : KpiType provides details about a KPI note.
type KpiType struct {
	Severity *string `json:"Severity,omitempty"`

	// The aggregation type of the KPI values. - SUM&#58; A single-value metrics aggregation type that sums up numeric
	// values
	//   that are extracted from KPI occurrences.
	AggregationType *string `json:"aggregation_type" validate:"required"`
}

// Constants associated with the KpiType.Severity property.
const (
	KpiTypeSeverityCriticalConst = "CRITICAL"
	KpiTypeSeverityHighConst     = "HIGH"
	KpiTypeSeverityLowConst      = "LOW"
	KpiTypeSeverityMediumConst   = "MEDIUM"
)

// Constants associated with the KpiType.AggregationType property.
// The aggregation type of the KPI values. - SUM&#58; A single-value metrics aggregation type that sums up numeric
// values
//   that are extracted from KPI occurrences.
const (
	KpiTypeAggregationTypeSumConst = "SUM"
)

// NewKpiType : Instantiate KpiType (Generic Model Constructor)
func (*FindingsV1) NewKpiType(aggregationType string) (_model *KpiType, err error) {
	_model = &KpiType{
		AggregationType: core.StringPtr(aggregationType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalKpiType unmarshals an instance of KpiType from the specified map of raw messages.
func UnmarshalKpiType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KpiType)
	err = core.UnmarshalPrimitive(m, "Severity", &obj.Severity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aggregation_type", &obj.AggregationType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListNoteOccurrencesOptions : The ListNoteOccurrences options.
type ListNoteOccurrencesOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// Second part of note `name`: providers/{provider_id}/notes/{note_id}.
	NoteID *string `json:"note_id" validate:"required,ne="`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Number of notes to return in the list.
	PageSize *int64 `json:"page_size,omitempty"`

	// Token to provide to skip to a particular spot in the list.
	PageToken *string `json:"page_token,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListNoteOccurrencesOptions : Instantiate ListNoteOccurrencesOptions
func (*FindingsV1) NewListNoteOccurrencesOptions(providerID string, noteID string) *ListNoteOccurrencesOptions {
	return &ListNoteOccurrencesOptions{
		ProviderID: core.StringPtr(providerID),
		NoteID:     core.StringPtr(noteID),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *ListNoteOccurrencesOptions) SetProviderID(providerID string) *ListNoteOccurrencesOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetNoteID : Allow user to set NoteID
func (_options *ListNoteOccurrencesOptions) SetNoteID(noteID string) *ListNoteOccurrencesOptions {
	_options.NoteID = core.StringPtr(noteID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListNoteOccurrencesOptions) SetTransactionID(transactionID string) *ListNoteOccurrencesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetPageSize : Allow user to set PageSize
func (_options *ListNoteOccurrencesOptions) SetPageSize(pageSize int64) *ListNoteOccurrencesOptions {
	_options.PageSize = core.Int64Ptr(pageSize)
	return _options
}

// SetPageToken : Allow user to set PageToken
func (_options *ListNoteOccurrencesOptions) SetPageToken(pageToken string) *ListNoteOccurrencesOptions {
	_options.PageToken = core.StringPtr(pageToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListNoteOccurrencesOptions) SetHeaders(param map[string]string) *ListNoteOccurrencesOptions {
	options.Headers = param
	return options
}

// ListNotesOptions : The ListNotes options.
type ListNotesOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Number of notes to return in the list.
	PageSize *int64 `json:"page_size,omitempty"`

	// Token to provide to skip to a particular spot in the list.
	PageToken *string `json:"page_token,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListNotesOptions : Instantiate ListNotesOptions
func (*FindingsV1) NewListNotesOptions(providerID string) *ListNotesOptions {
	return &ListNotesOptions{
		ProviderID: core.StringPtr(providerID),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *ListNotesOptions) SetProviderID(providerID string) *ListNotesOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListNotesOptions) SetTransactionID(transactionID string) *ListNotesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetPageSize : Allow user to set PageSize
func (_options *ListNotesOptions) SetPageSize(pageSize int64) *ListNotesOptions {
	_options.PageSize = core.Int64Ptr(pageSize)
	return _options
}

// SetPageToken : Allow user to set PageToken
func (_options *ListNotesOptions) SetPageToken(pageToken string) *ListNotesOptions {
	_options.PageToken = core.StringPtr(pageToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListNotesOptions) SetHeaders(param map[string]string) *ListNotesOptions {
	options.Headers = param
	return options
}

// ListOccurrencesOptions : The ListOccurrences options.
type ListOccurrencesOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Number of notes to return in the list.
	PageSize *int64 `json:"page_size,omitempty"`

	// Token to provide to skip to a particular spot in the list.
	PageToken *string `json:"page_token,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListOccurrencesOptions : Instantiate ListOccurrencesOptions
func (*FindingsV1) NewListOccurrencesOptions(providerID string) *ListOccurrencesOptions {
	return &ListOccurrencesOptions{
		ProviderID: core.StringPtr(providerID),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *ListOccurrencesOptions) SetProviderID(providerID string) *ListOccurrencesOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListOccurrencesOptions) SetTransactionID(transactionID string) *ListOccurrencesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetPageSize : Allow user to set PageSize
func (_options *ListOccurrencesOptions) SetPageSize(pageSize int64) *ListOccurrencesOptions {
	_options.PageSize = core.Int64Ptr(pageSize)
	return _options
}

// SetPageToken : Allow user to set PageToken
func (_options *ListOccurrencesOptions) SetPageToken(pageToken string) *ListOccurrencesOptions {
	_options.PageToken = core.StringPtr(pageToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListOccurrencesOptions) SetHeaders(param map[string]string) *ListOccurrencesOptions {
	options.Headers = param
	return options
}

// ListProvidersOptions : The ListProviders options.
type ListProvidersOptions struct {
	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// The number of documents that you want to return.
	Limit *int64 `json:"limit,omitempty"`

	// The offset is the index of the item from which you want to start returning data from. Default is 0.
	Skip *int64 `json:"skip,omitempty"`

	// The first provider ID included in the result, sorted in ascending order. If not provided, this parameter is ignored.
	StartProviderID *string `json:"start_provider_id,omitempty"`

	// The last provider ID included in the result, sorted in ascending order. If not provided, this parameter is ignored.
	EndProviderID *string `json:"end_provider_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProvidersOptions : Instantiate ListProvidersOptions
func (*FindingsV1) NewListProvidersOptions() *ListProvidersOptions {
	return &ListProvidersOptions{}
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListProvidersOptions) SetTransactionID(transactionID string) *ListProvidersOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListProvidersOptions) SetLimit(limit int64) *ListProvidersOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *ListProvidersOptions) SetSkip(skip int64) *ListProvidersOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetStartProviderID : Allow user to set StartProviderID
func (_options *ListProvidersOptions) SetStartProviderID(startProviderID string) *ListProvidersOptions {
	_options.StartProviderID = core.StringPtr(startProviderID)
	return _options
}

// SetEndProviderID : Allow user to set EndProviderID
func (_options *ListProvidersOptions) SetEndProviderID(endProviderID string) *ListProvidersOptions {
	_options.EndProviderID = core.StringPtr(endProviderID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProvidersOptions) SetHeaders(param map[string]string) *ListProvidersOptions {
	options.Headers = param
	return options
}

// NetworkConnection : It provides details about a network connection.
type NetworkConnection struct {
	// The direction of this network connection.
	Direction *string `json:"direction,omitempty"`

	// The protocol of this network connection.
	Protocol *string `json:"protocol,omitempty"`

	// It provides details about a socket address.
	Client *SocketAddress `json:"client,omitempty"`

	// It provides details about a socket address.
	Server *SocketAddress `json:"server,omitempty"`
}

// UnmarshalNetworkConnection unmarshals an instance of NetworkConnection from the specified map of raw messages.
func UnmarshalNetworkConnection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NetworkConnection)
	err = core.UnmarshalPrimitive(m, "direction", &obj.Direction)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "protocol", &obj.Protocol)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "client", &obj.Client, UnmarshalSocketAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "server", &obj.Server, UnmarshalSocketAddress)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PostGraphOptions : The PostGraph options.
type PostGraphOptions struct {
	// Body for query findings.
	Body io.ReadCloser `json:"body,omitempty"`

	// The type of the input.
	ContentType *string `json:"Content-Type,omitempty"`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostGraphOptions : Instantiate PostGraphOptions
func (*FindingsV1) NewPostGraphOptions() *PostGraphOptions {
	return &PostGraphOptions{}
}

// SetBody : Allow user to set Body
func (_options *PostGraphOptions) SetBody(body io.ReadCloser) *PostGraphOptions {
	_options.Body = body
	return _options
}

// SetContentType : Allow user to set ContentType
func (_options *PostGraphOptions) SetContentType(contentType string) *PostGraphOptions {
	_options.ContentType = core.StringPtr(contentType)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *PostGraphOptions) SetTransactionID(transactionID string) *PostGraphOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostGraphOptions) SetHeaders(param map[string]string) *PostGraphOptions {
	options.Headers = param
	return options
}

// RemediationStep : A remediation step description and associated URL.
type RemediationStep struct {
	// Title of this next step.
	Title *string `json:"title,omitempty"`

	// The URL associated to this next steps.
	URL *string `json:"url,omitempty"`
}

// UnmarshalRemediationStep unmarshals an instance of RemediationStep from the specified map of raw messages.
func UnmarshalRemediationStep(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RemediationStep)
	err = core.UnmarshalPrimitive(m, "title", &obj.Title)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Reporter : The entity reporting a note.
type Reporter struct {
	// The id of this reporter.
	ID *string `json:"id" validate:"required"`

	// The title of this reporter.
	Title *string `json:"title" validate:"required"`

	// The url of this reporter.
	URL *string `json:"url,omitempty"`
}

// NewReporter : Instantiate Reporter (Generic Model Constructor)
func (*FindingsV1) NewReporter(id string, title string) (_model *Reporter, err error) {
	_model = &Reporter{
		ID:    core.StringPtr(id),
		Title: core.StringPtr(title),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalReporter unmarshals an instance of Reporter from the specified map of raw messages.
func UnmarshalReporter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Reporter)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "title", &obj.Title)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Section : Card provides details about a card kind of note.
type Section struct {
	// The title of this section.
	Title *string `json:"title" validate:"required"`

	// The image of this section.
	Image *string `json:"image" validate:"required"`
}

// NewSection : Instantiate Section (Generic Model Constructor)
func (*FindingsV1) NewSection(title string, image string) (_model *Section, err error) {
	_model = &Section{
		Title: core.StringPtr(title),
		Image: core.StringPtr(image),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSection unmarshals an instance of Section from the specified map of raw messages.
func UnmarshalSection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Section)
	err = core.UnmarshalPrimitive(m, "title", &obj.Title)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "image", &obj.Image)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SocketAddress : It provides details about a socket address.
type SocketAddress struct {
	// The IP address of this socket address.
	Address *string `json:"address" validate:"required"`

	// The port number of this socket address.
	Port *int64 `json:"port,omitempty"`
}

// NewSocketAddress : Instantiate SocketAddress (Generic Model Constructor)
func (*FindingsV1) NewSocketAddress(address string) (_model *SocketAddress, err error) {
	_model = &SocketAddress{
		Address: core.StringPtr(address),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSocketAddress unmarshals an instance of SocketAddress from the specified map of raw messages.
func UnmarshalSocketAddress(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SocketAddress)
	err = core.UnmarshalPrimitive(m, "address", &obj.Address)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateNoteOptions : The UpdateNote options.
type UpdateNoteOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// Second part of note `name`: providers/{provider_id}/notes/{note_id}.
	NoteID *string `json:"note_id" validate:"required,ne="`

	// A one sentence description of your note.
	ShortDescription *string `json:"short_description" validate:"required"`

	// A more detailed description of your note.
	LongDescription *string `json:"long_description" validate:"required"`

	// The type of note. Use this field to filter notes and occurences by kind.
	//  - FINDING&#58; The note and occurrence represent a finding.
	//  - KPI&#58; The note and occurrence represent a KPI value.
	//  - CARD&#58; The note represents a card showing findings and related metric values.
	//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
	//  - SECTION&#58; The note represents a section in a dashboard.
	Kind *string `json:"kind" validate:"required"`

	// The ID of the note.
	ID *string `json:"id" validate:"required"`

	// The entity reporting a note.
	ReportedBy *Reporter `json:"reported_by" validate:"required"`

	RelatedURL []APINoteRelatedURL `json:"related_url,omitempty"`

	// Output only. The time this note was created. This field can be used as a filter in list requests.
	CreateTime *strfmt.DateTime `json:"create_time,omitempty"`

	// Output only. The time this note was last updated. This field can be used as a filter in list requests.
	UpdateTime *strfmt.DateTime `json:"update_time,omitempty"`

	// True if this note can be shared by multiple accounts.
	Shared *bool `json:"shared,omitempty"`

	// FindingType provides details about a finding note.
	Finding *FindingType `json:"finding,omitempty"`

	// KpiType provides details about a KPI note.
	Kpi *KpiType `json:"kpi,omitempty"`

	// Card provides details about a card kind of note.
	Card *Card `json:"card,omitempty"`

	// Card provides details about a card kind of note.
	Section *Section `json:"section,omitempty"`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateNoteOptions.Kind property.
// The type of note. Use this field to filter notes and occurences by kind.
//  - FINDING&#58; The note and occurrence represent a finding.
//  - KPI&#58; The note and occurrence represent a KPI value.
//  - CARD&#58; The note represents a card showing findings and related metric values.
//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
//  - SECTION&#58; The note represents a section in a dashboard.
const (
	UpdateNoteOptionsKindCardConst           = "CARD"
	UpdateNoteOptionsKindCardConfiguredConst = "CARD_CONFIGURED"
	UpdateNoteOptionsKindFindingConst        = "FINDING"
	UpdateNoteOptionsKindKpiConst            = "KPI"
	UpdateNoteOptionsKindSectionConst        = "SECTION"
)

// NewUpdateNoteOptions : Instantiate UpdateNoteOptions
func (*FindingsV1) NewUpdateNoteOptions(providerID string, noteID string, shortDescription string, longDescription string, kind string, id string, reportedBy *Reporter) *UpdateNoteOptions {
	return &UpdateNoteOptions{
		ProviderID:       core.StringPtr(providerID),
		NoteID:           core.StringPtr(noteID),
		ShortDescription: core.StringPtr(shortDescription),
		LongDescription:  core.StringPtr(longDescription),
		Kind:             core.StringPtr(kind),
		ID:               core.StringPtr(id),
		ReportedBy:       reportedBy,
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *UpdateNoteOptions) SetProviderID(providerID string) *UpdateNoteOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetNoteID : Allow user to set NoteID
func (_options *UpdateNoteOptions) SetNoteID(noteID string) *UpdateNoteOptions {
	_options.NoteID = core.StringPtr(noteID)
	return _options
}

// SetShortDescription : Allow user to set ShortDescription
func (_options *UpdateNoteOptions) SetShortDescription(shortDescription string) *UpdateNoteOptions {
	_options.ShortDescription = core.StringPtr(shortDescription)
	return _options
}

// SetLongDescription : Allow user to set LongDescription
func (_options *UpdateNoteOptions) SetLongDescription(longDescription string) *UpdateNoteOptions {
	_options.LongDescription = core.StringPtr(longDescription)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *UpdateNoteOptions) SetKind(kind string) *UpdateNoteOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateNoteOptions) SetID(id string) *UpdateNoteOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetReportedBy : Allow user to set ReportedBy
func (_options *UpdateNoteOptions) SetReportedBy(reportedBy *Reporter) *UpdateNoteOptions {
	_options.ReportedBy = reportedBy
	return _options
}

// SetRelatedURL : Allow user to set RelatedURL
func (_options *UpdateNoteOptions) SetRelatedURL(relatedURL []APINoteRelatedURL) *UpdateNoteOptions {
	_options.RelatedURL = relatedURL
	return _options
}

// SetCreateTime : Allow user to set CreateTime
func (_options *UpdateNoteOptions) SetCreateTime(createTime *strfmt.DateTime) *UpdateNoteOptions {
	_options.CreateTime = createTime
	return _options
}

// SetUpdateTime : Allow user to set UpdateTime
func (_options *UpdateNoteOptions) SetUpdateTime(updateTime *strfmt.DateTime) *UpdateNoteOptions {
	_options.UpdateTime = updateTime
	return _options
}

// SetShared : Allow user to set Shared
func (_options *UpdateNoteOptions) SetShared(shared bool) *UpdateNoteOptions {
	_options.Shared = core.BoolPtr(shared)
	return _options
}

// SetFinding : Allow user to set Finding
func (_options *UpdateNoteOptions) SetFinding(finding *FindingType) *UpdateNoteOptions {
	_options.Finding = finding
	return _options
}

// SetKpi : Allow user to set Kpi
func (_options *UpdateNoteOptions) SetKpi(kpi *KpiType) *UpdateNoteOptions {
	_options.Kpi = kpi
	return _options
}

// SetCard : Allow user to set Card
func (_options *UpdateNoteOptions) SetCard(card *Card) *UpdateNoteOptions {
	_options.Card = card
	return _options
}

// SetSection : Allow user to set Section
func (_options *UpdateNoteOptions) SetSection(section *Section) *UpdateNoteOptions {
	_options.Section = section
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateNoteOptions) SetTransactionID(transactionID string) *UpdateNoteOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateNoteOptions) SetHeaders(param map[string]string) *UpdateNoteOptions {
	options.Headers = param
	return options
}

// UpdateOccurrenceOptions : The UpdateOccurrence options.
type UpdateOccurrenceOptions struct {
	// Part of the parent. This field contains the provider ID. For example: providers/{provider_id}.
	ProviderID *string `json:"provider_id" validate:"required,ne="`

	// Second part of occurrence `name`: providers/{provider_id}/occurrences/{occurrence_id}.
	OccurrenceID *string `json:"occurrence_id" validate:"required,ne="`

	// An analysis note associated with this image, in the form "{account_id}/providers/{provider_id}/notes/{note_id}" This
	// field can be used as a filter in list requests.
	NoteName *string `json:"note_name" validate:"required"`

	// The type of note. Use this field to filter notes and occurences by kind.
	//  - FINDING&#58; The note and occurrence represent a finding.
	//  - KPI&#58; The note and occurrence represent a KPI value.
	//  - CARD&#58; The note represents a card showing findings and related metric values.
	//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
	//  - SECTION&#58; The note represents a section in a dashboard.
	Kind *string `json:"kind" validate:"required"`

	// The id of the occurrence.
	ID *string `json:"id" validate:"required"`

	// The unique URL of the resource, image or the container, for which the `Occurrence` applies. For example,
	// https://gcr.io/provider/image@sha256:foo. This field can be used as a filter in list requests.
	ResourceURL *string `json:"resource_url,omitempty"`

	// A description of actions that can be taken to remedy the `Note`.
	Remediation *string `json:"remediation,omitempty"`

	// Output only. The time this `Occurrence` was created.
	CreateTime *strfmt.DateTime `json:"create_time,omitempty"`

	// Output only. The time this `Occurrence` was last updated.
	UpdateTime *strfmt.DateTime `json:"update_time,omitempty"`

	Context *Context `json:"context,omitempty"`

	// Finding provides details about a finding occurrence.
	Finding *Finding `json:"finding,omitempty"`

	// Kpi provides details about a KPI occurrence.
	Kpi *Kpi `json:"kpi,omitempty"`

	// Additional data for the finding, like AT event etc.
	ReferenceData interface{} `json:"reference_data,omitempty"`

	// The transaction ID for the request in UUID v4 format.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateOccurrenceOptions.Kind property.
// The type of note. Use this field to filter notes and occurences by kind.
//  - FINDING&#58; The note and occurrence represent a finding.
//  - KPI&#58; The note and occurrence represent a KPI value.
//  - CARD&#58; The note represents a card showing findings and related metric values.
//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
//  - SECTION&#58; The note represents a section in a dashboard.
const (
	UpdateOccurrenceOptionsKindCardConst           = "CARD"
	UpdateOccurrenceOptionsKindCardConfiguredConst = "CARD_CONFIGURED"
	UpdateOccurrenceOptionsKindFindingConst        = "FINDING"
	UpdateOccurrenceOptionsKindKpiConst            = "KPI"
	UpdateOccurrenceOptionsKindSectionConst        = "SECTION"
)

// NewUpdateOccurrenceOptions : Instantiate UpdateOccurrenceOptions
func (*FindingsV1) NewUpdateOccurrenceOptions(providerID string, occurrenceID string, noteName string, kind string, id string) *UpdateOccurrenceOptions {
	return &UpdateOccurrenceOptions{
		ProviderID:   core.StringPtr(providerID),
		OccurrenceID: core.StringPtr(occurrenceID),
		NoteName:     core.StringPtr(noteName),
		Kind:         core.StringPtr(kind),
		ID:           core.StringPtr(id),
	}
}

// SetProviderID : Allow user to set ProviderID
func (_options *UpdateOccurrenceOptions) SetProviderID(providerID string) *UpdateOccurrenceOptions {
	_options.ProviderID = core.StringPtr(providerID)
	return _options
}

// SetOccurrenceID : Allow user to set OccurrenceID
func (_options *UpdateOccurrenceOptions) SetOccurrenceID(occurrenceID string) *UpdateOccurrenceOptions {
	_options.OccurrenceID = core.StringPtr(occurrenceID)
	return _options
}

// SetNoteName : Allow user to set NoteName
func (_options *UpdateOccurrenceOptions) SetNoteName(noteName string) *UpdateOccurrenceOptions {
	_options.NoteName = core.StringPtr(noteName)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *UpdateOccurrenceOptions) SetKind(kind string) *UpdateOccurrenceOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateOccurrenceOptions) SetID(id string) *UpdateOccurrenceOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetResourceURL : Allow user to set ResourceURL
func (_options *UpdateOccurrenceOptions) SetResourceURL(resourceURL string) *UpdateOccurrenceOptions {
	_options.ResourceURL = core.StringPtr(resourceURL)
	return _options
}

// SetRemediation : Allow user to set Remediation
func (_options *UpdateOccurrenceOptions) SetRemediation(remediation string) *UpdateOccurrenceOptions {
	_options.Remediation = core.StringPtr(remediation)
	return _options
}

// SetCreateTime : Allow user to set CreateTime
func (_options *UpdateOccurrenceOptions) SetCreateTime(createTime *strfmt.DateTime) *UpdateOccurrenceOptions {
	_options.CreateTime = createTime
	return _options
}

// SetUpdateTime : Allow user to set UpdateTime
func (_options *UpdateOccurrenceOptions) SetUpdateTime(updateTime *strfmt.DateTime) *UpdateOccurrenceOptions {
	_options.UpdateTime = updateTime
	return _options
}

// SetContext : Allow user to set Context
func (_options *UpdateOccurrenceOptions) SetContext(context *Context) *UpdateOccurrenceOptions {
	_options.Context = context
	return _options
}

// SetFinding : Allow user to set Finding
func (_options *UpdateOccurrenceOptions) SetFinding(finding *Finding) *UpdateOccurrenceOptions {
	_options.Finding = finding
	return _options
}

// SetKpi : Allow user to set Kpi
func (_options *UpdateOccurrenceOptions) SetKpi(kpi *Kpi) *UpdateOccurrenceOptions {
	_options.Kpi = kpi
	return _options
}

// SetReferenceData : Allow user to set ReferenceData
func (_options *UpdateOccurrenceOptions) SetReferenceData(referenceData interface{}) *UpdateOccurrenceOptions {
	_options.ReferenceData = referenceData
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateOccurrenceOptions) SetTransactionID(transactionID string) *UpdateOccurrenceOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateOccurrenceOptions) SetHeaders(param map[string]string) *UpdateOccurrenceOptions {
	options.Headers = param
	return options
}

// ValueType : the value type of a card element.
// Models which "extend" this model:
// - ValueTypeKpiValueType
// - ValueTypeFindingCountValueType
type ValueType struct {
	// Kind of element
	// - KPI&#58; Kind of value derived from a KPI occurrence.
	Kind *string `json:"kind,omitempty"`

	// The name of the kpi note associated to the occurrence with the value for this card element value type.
	KpiNoteName *string `json:"kpi_note_name,omitempty"`

	// The text of this element type.
	Text *string `json:"text,omitempty"`

	// the names of the finding note associated that act as filters for counting the occurrences.
	FindingNoteNames []string `json:"finding_note_names,omitempty"`
}

// Constants associated with the ValueType.Kind property.
// Kind of element
// - KPI&#58; Kind of value derived from a KPI occurrence.
const (
	ValueTypeKindKpiConst = "KPI"
)

func (*ValueType) isaValueType() bool {
	return true
}

type ValueTypeIntf interface {
	isaValueType() bool
}

// UnmarshalValueType unmarshals an instance of ValueType from the specified map of raw messages.
func UnmarshalValueType(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "kind", &discValue)
	if err != nil {
		err = fmt.Errorf("error unmarshalling discriminator property 'kind': %s", err.Error())
		return
	}
	if discValue == "" {
		err = fmt.Errorf("required discriminator property 'kind' not found in JSON object")
		return
	}
	if discValue == "KPI" {
		err = core.UnmarshalModel(m, "", result, UnmarshalValueTypeKpiValueType)
	} else if discValue == "FINDING_COUNT" {
		err = core.UnmarshalModel(m, "", result, UnmarshalValueTypeFindingCountValueType)
	} else {
		err = fmt.Errorf("unrecognized value for discriminator property 'kind': %s", discValue)
	}
	return
}

// APIListNoteOccurrencesResponse : Response including listed occurrences for a note.
type APIListNoteOccurrencesResponse struct {
	// The occurrences attached to the specified note.
	Occurrences []APIOccurrence `json:"occurrences,omitempty"`

	// Token to receive the next page of notes.
	NextPageToken *string `json:"next_page_token,omitempty"`
}

// UnmarshalAPIListNoteOccurrencesResponse unmarshals an instance of APIListNoteOccurrencesResponse from the specified map of raw messages.
func UnmarshalAPIListNoteOccurrencesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIListNoteOccurrencesResponse)
	err = core.UnmarshalModel(m, "occurrences", &obj.Occurrences, UnmarshalAPIOccurrence)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_page_token", &obj.NextPageToken)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APIListNotesResponse : Response including listed notes.
type APIListNotesResponse struct {
	// The occurrences requested.
	Notes []APINote `json:"notes,omitempty"`

	// The next pagination token in the list response. It should be used as page_token for the following request. An empty
	// value means no more result.
	NextPageToken *string `json:"next_page_token,omitempty"`
}

// UnmarshalAPIListNotesResponse unmarshals an instance of APIListNotesResponse from the specified map of raw messages.
func UnmarshalAPIListNotesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIListNotesResponse)
	err = core.UnmarshalModel(m, "notes", &obj.Notes, UnmarshalAPINote)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_page_token", &obj.NextPageToken)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APIListOccurrencesResponse : Response including listed active occurrences.
type APIListOccurrencesResponse struct {
	// The occurrences requested.
	Occurrences []APIOccurrence `json:"occurrences,omitempty"`

	// The next pagination token in the list response. It should be used as
	// `page_token` for the following request. An empty value means no more results.
	NextPageToken *string `json:"next_page_token,omitempty"`
}

// UnmarshalAPIListOccurrencesResponse unmarshals an instance of APIListOccurrencesResponse from the specified map of raw messages.
func UnmarshalAPIListOccurrencesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIListOccurrencesResponse)
	err = core.UnmarshalModel(m, "occurrences", &obj.Occurrences, UnmarshalAPIOccurrence)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_page_token", &obj.NextPageToken)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APIListProvidersResponse : A list of providers is returned.
type APIListProvidersResponse struct {
	// The providers requested.
	Providers []APIProvider `json:"providers,omitempty"`

	// The number of elements returned in the current instance. The default is 200.
	Limit *int64 `json:"limit,omitempty"`

	// The offset is the index of the item from which you want to start returning data from. The default is 0.
	Skip *int64 `json:"skip,omitempty"`

	// The total number of providers available.
	TotalCount *int64 `json:"total_count,omitempty"`
}

// UnmarshalAPIListProvidersResponse unmarshals an instance of APIListProvidersResponse from the specified map of raw messages.
func UnmarshalAPIListProvidersResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIListProvidersResponse)
	err = core.UnmarshalModel(m, "providers", &obj.Providers, UnmarshalAPIProvider)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "skip", &obj.Skip)
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

// APINote : Provides a detailed description of a note.
type APINote struct {
	// A one sentence description of your note.
	ShortDescription *string `json:"short_description" validate:"required"`

	// A more detailed description of your note.
	LongDescription *string `json:"long_description" validate:"required"`

	// The type of note. Use this field to filter notes and occurences by kind.
	//  - FINDING&#58; The note and occurrence represent a finding.
	//  - KPI&#58; The note and occurrence represent a KPI value.
	//  - CARD&#58; The note represents a card showing findings and related metric values.
	//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
	//  - SECTION&#58; The note represents a section in a dashboard.
	Kind *string `json:"kind" validate:"required"`

	RelatedURL []APINoteRelatedURL `json:"related_url,omitempty"`

	// Output only. The time this note was created. This field can be used as a filter in list requests.
	CreateTime *strfmt.DateTime `json:"create_time,omitempty"`

	// Output only. The time this note was last updated. This field can be used as a filter in list requests.
	UpdateTime *strfmt.DateTime `json:"update_time,omitempty"`

	// The ID of the note.
	ID *string `json:"id" validate:"required"`

	// True if this note can be shared by multiple accounts.
	Shared *bool `json:"shared,omitempty"`

	// The entity reporting a note.
	ReportedBy *Reporter `json:"reported_by" validate:"required"`

	// FindingType provides details about a finding note.
	Finding *FindingType `json:"finding,omitempty"`

	// KpiType provides details about a KPI note.
	Kpi *KpiType `json:"kpi,omitempty"`

	// Card provides details about a card kind of note.
	Card *Card `json:"card,omitempty"`

	// Card provides details about a card kind of note.
	Section *Section `json:"section,omitempty"`
}

// Constants associated with the APINote.Kind property.
// The type of note. Use this field to filter notes and occurences by kind.
//  - FINDING&#58; The note and occurrence represent a finding.
//  - KPI&#58; The note and occurrence represent a KPI value.
//  - CARD&#58; The note represents a card showing findings and related metric values.
//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
//  - SECTION&#58; The note represents a section in a dashboard.
const (
	APINoteKindCardConst           = "CARD"
	APINoteKindCardConfiguredConst = "CARD_CONFIGURED"
	APINoteKindFindingConst        = "FINDING"
	APINoteKindKpiConst            = "KPI"
	APINoteKindSectionConst        = "SECTION"
)

// NewAPINote : Instantiate APINote (Generic Model Constructor)
func (*FindingsV1) NewAPINote(shortDescription string, longDescription string, kind string, id string, reportedBy *Reporter) (_model *APINote, err error) {
	_model = &APINote{
		ShortDescription: core.StringPtr(shortDescription),
		LongDescription:  core.StringPtr(longDescription),
		Kind:             core.StringPtr(kind),
		ID:               core.StringPtr(id),
		ReportedBy:       reportedBy,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalAPINote unmarshals an instance of APINote from the specified map of raw messages.
func UnmarshalAPINote(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APINote)
	err = core.UnmarshalPrimitive(m, "short_description", &obj.ShortDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "related_url", &obj.RelatedURL, UnmarshalAPINoteRelatedURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "create_time", &obj.CreateTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update_time", &obj.UpdateTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "shared", &obj.Shared)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "reported_by", &obj.ReportedBy, UnmarshalReporter)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "finding", &obj.Finding, UnmarshalFindingType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "kpi", &obj.Kpi, UnmarshalKpiType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "card", &obj.Card, UnmarshalCard)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "section", &obj.Section, UnmarshalSection)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APINoteRelatedURL : Metadata for any related URL information.
type APINoteRelatedURL struct {
	// Label to describe usage of the URL.
	Label *string `json:"label" validate:"required"`

	// The URL that you want to associate with the note.
	URL *string `json:"url" validate:"required"`
}

// NewAPINoteRelatedURL : Instantiate APINoteRelatedURL (Generic Model Constructor)
func (*FindingsV1) NewAPINoteRelatedURL(label string, url string) (_model *APINoteRelatedURL, err error) {
	_model = &APINoteRelatedURL{
		Label: core.StringPtr(label),
		URL:   core.StringPtr(url),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalAPINoteRelatedURL unmarshals an instance of APINoteRelatedURL from the specified map of raw messages.
func UnmarshalAPINoteRelatedURL(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APINoteRelatedURL)
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APIOccurrence : `Occurrence` includes information about analysis occurrences for an image.
type APIOccurrence struct {
	// The unique URL of the resource, image or the container, for which the `Occurrence` applies. For example,
	// https://gcr.io/provider/image@sha256:foo. This field can be used as a filter in list requests.
	ResourceURL *string `json:"resource_url,omitempty"`

	// An analysis note associated with this image, in the form "{account_id}/providers/{provider_id}/notes/{note_id}" This
	// field can be used as a filter in list requests.
	NoteName *string `json:"note_name" validate:"required"`

	// The type of note. Use this field to filter notes and occurences by kind.
	//  - FINDING&#58; The note and occurrence represent a finding.
	//  - KPI&#58; The note and occurrence represent a KPI value.
	//  - CARD&#58; The note represents a card showing findings and related metric values.
	//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
	//  - SECTION&#58; The note represents a section in a dashboard.
	Kind *string `json:"kind" validate:"required"`

	// A description of actions that can be taken to remedy the `Note`.
	Remediation *string `json:"remediation,omitempty"`

	// Output only. The time this `Occurrence` was created.
	CreateTime *strfmt.DateTime `json:"create_time,omitempty"`

	// Output only. The time this `Occurrence` was last updated.
	UpdateTime *strfmt.DateTime `json:"update_time,omitempty"`

	// The id of the occurrence.
	ID *string `json:"id" validate:"required"`

	Context *Context `json:"context,omitempty"`

	// Finding provides details about a finding occurrence.
	Finding *Finding `json:"finding,omitempty"`

	// Kpi provides details about a KPI occurrence.
	Kpi *Kpi `json:"kpi,omitempty"`

	// Additional data for the finding, like AT event etc.
	ReferenceData interface{} `json:"reference_data,omitempty"`
}

// Constants associated with the APIOccurrence.Kind property.
// The type of note. Use this field to filter notes and occurences by kind.
//  - FINDING&#58; The note and occurrence represent a finding.
//  - KPI&#58; The note and occurrence represent a KPI value.
//  - CARD&#58; The note represents a card showing findings and related metric values.
//  - CARD_CONFIGURED&#58; The note represents a card configured for a user account.
//  - SECTION&#58; The note represents a section in a dashboard.
const (
	APIOccurrenceKindCardConst           = "CARD"
	APIOccurrenceKindCardConfiguredConst = "CARD_CONFIGURED"
	APIOccurrenceKindFindingConst        = "FINDING"
	APIOccurrenceKindKpiConst            = "KPI"
	APIOccurrenceKindSectionConst        = "SECTION"
)

// NewAPIOccurrence : Instantiate APIOccurrence (Generic Model Constructor)
func (*FindingsV1) NewAPIOccurrence(noteName string, kind string, id string) (_model *APIOccurrence, err error) {
	_model = &APIOccurrence{
		NoteName: core.StringPtr(noteName),
		Kind:     core.StringPtr(kind),
		ID:       core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalAPIOccurrence unmarshals an instance of APIOccurrence from the specified map of raw messages.
func UnmarshalAPIOccurrence(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIOccurrence)
	err = core.UnmarshalPrimitive(m, "resource_url", &obj.ResourceURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "note_name", &obj.NoteName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "remediation", &obj.Remediation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "create_time", &obj.CreateTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update_time", &obj.UpdateTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "context", &obj.Context, UnmarshalContext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "finding", &obj.Finding, UnmarshalFinding)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "kpi", &obj.Kpi, UnmarshalKpi)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reference_data", &obj.ReferenceData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// APIProvider : Provides a detailed description of a provider.
type APIProvider struct {
	// The name of the provider in the form '{account_id}/providers/{provider_id}'.
	Name *string `json:"name" validate:"required"`

	// The ID of the provider.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalAPIProvider unmarshals an instance of APIProvider from the specified map of raw messages.
func UnmarshalAPIProvider(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIProvider)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CardElementBreakdownCardElement : A card element with a breakdown of values.
// This model "extends" CardElement
type CardElementBreakdownCardElement struct {
	// The text of this card element.
	Text *string `json:"text" validate:"required"`

	// Kind of element
	// - NUMERIC&#58; Single numeric value
	// - BREAKDOWN&#58; Breakdown of numeric values
	// - TIME_SERIES&#58; Time-series of numeric values.
	Kind *string `json:"kind" validate:"required"`

	// The default time range of this card element.
	DefaultTimeRange *string `json:"default_time_range,omitempty"`

	// the value types associated to this card element.
	ValueTypes []ValueTypeIntf `json:"value_types" validate:"required"`
}

// Constants associated with the CardElementBreakdownCardElement.Kind property.
// Kind of element
// - NUMERIC&#58; Single numeric value
// - BREAKDOWN&#58; Breakdown of numeric values
// - TIME_SERIES&#58; Time-series of numeric values.
const (
	CardElementBreakdownCardElementKindBreakdownConst  = "BREAKDOWN"
	CardElementBreakdownCardElementKindNumericConst    = "NUMERIC"
	CardElementBreakdownCardElementKindTimeSeriesConst = "TIME_SERIES"
)

// NewCardElementBreakdownCardElement : Instantiate CardElementBreakdownCardElement (Generic Model Constructor)
func (*FindingsV1) NewCardElementBreakdownCardElement(text string, kind string, valueTypes []ValueTypeIntf) (_model *CardElementBreakdownCardElement, err error) {
	_model = &CardElementBreakdownCardElement{
		Text:       core.StringPtr(text),
		Kind:       core.StringPtr(kind),
		ValueTypes: valueTypes,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*CardElementBreakdownCardElement) isaCardElement() bool {
	return true
}

// UnmarshalCardElementBreakdownCardElement unmarshals an instance of CardElementBreakdownCardElement from the specified map of raw messages.
func UnmarshalCardElementBreakdownCardElement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CardElementBreakdownCardElement)
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_time_range", &obj.DefaultTimeRange)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "value_types", &obj.ValueTypes, UnmarshalValueType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CardElementNumericCardElement : A card element with a single numeric value".
// This model "extends" CardElement
type CardElementNumericCardElement struct {
	// The text of this card element.
	Text *string `json:"text" validate:"required"`

	// Kind of element
	// - NUMERIC&#58; Single numeric value
	// - BREAKDOWN&#58; Breakdown of numeric values
	// - TIME_SERIES&#58; Time-series of numeric values.
	Kind *string `json:"kind" validate:"required"`

	// The default time range of this card element.
	DefaultTimeRange *string `json:"default_time_range,omitempty"`

	ValueType *NumericCardElementValueType `json:"value_type" validate:"required"`
}

// Constants associated with the CardElementNumericCardElement.Kind property.
// Kind of element
// - NUMERIC&#58; Single numeric value
// - BREAKDOWN&#58; Breakdown of numeric values
// - TIME_SERIES&#58; Time-series of numeric values.
const (
	CardElementNumericCardElementKindBreakdownConst  = "BREAKDOWN"
	CardElementNumericCardElementKindNumericConst    = "NUMERIC"
	CardElementNumericCardElementKindTimeSeriesConst = "TIME_SERIES"
)

// NewCardElementNumericCardElement : Instantiate CardElementNumericCardElement (Generic Model Constructor)
func (*FindingsV1) NewCardElementNumericCardElement(text string, kind string, valueType *NumericCardElementValueType) (_model *CardElementNumericCardElement, err error) {
	_model = &CardElementNumericCardElement{
		Text:      core.StringPtr(text),
		Kind:      core.StringPtr(kind),
		ValueType: valueType,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*CardElementNumericCardElement) isaCardElement() bool {
	return true
}

// UnmarshalCardElementNumericCardElement unmarshals an instance of CardElementNumericCardElement from the specified map of raw messages.
func UnmarshalCardElementNumericCardElement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CardElementNumericCardElement)
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_time_range", &obj.DefaultTimeRange)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "value_type", &obj.ValueType, UnmarshalNumericCardElementValueType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CardElementTimeSeriesCardElement : A card element with a time series chart.
// This model "extends" CardElement
type CardElementTimeSeriesCardElement struct {
	// The text of this card element.
	Text *string `json:"text" validate:"required"`

	// The default interval of the time series.
	DefaultInterval *string `json:"default_interval,omitempty"`

	// Kind of element
	// - NUMERIC&#58; Single numeric value
	// - BREAKDOWN&#58; Breakdown of numeric values
	// - TIME_SERIES&#58; Time-series of numeric values.
	Kind *string `json:"kind" validate:"required"`

	// The default time range of this card element.
	DefaultTimeRange *string `json:"default_time_range,omitempty"`

	// the value types associated to this card element.
	ValueTypes []ValueTypeIntf `json:"value_types" validate:"required"`
}

// Constants associated with the CardElementTimeSeriesCardElement.Kind property.
// Kind of element
// - NUMERIC&#58; Single numeric value
// - BREAKDOWN&#58; Breakdown of numeric values
// - TIME_SERIES&#58; Time-series of numeric values.
const (
	CardElementTimeSeriesCardElementKindBreakdownConst  = "BREAKDOWN"
	CardElementTimeSeriesCardElementKindNumericConst    = "NUMERIC"
	CardElementTimeSeriesCardElementKindTimeSeriesConst = "TIME_SERIES"
)

// NewCardElementTimeSeriesCardElement : Instantiate CardElementTimeSeriesCardElement (Generic Model Constructor)
func (*FindingsV1) NewCardElementTimeSeriesCardElement(text string, kind string, valueTypes []ValueTypeIntf) (_model *CardElementTimeSeriesCardElement, err error) {
	_model = &CardElementTimeSeriesCardElement{
		Text:       core.StringPtr(text),
		Kind:       core.StringPtr(kind),
		ValueTypes: valueTypes,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*CardElementTimeSeriesCardElement) isaCardElement() bool {
	return true
}

// UnmarshalCardElementTimeSeriesCardElement unmarshals an instance of CardElementTimeSeriesCardElement from the specified map of raw messages.
func UnmarshalCardElementTimeSeriesCardElement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CardElementTimeSeriesCardElement)
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_interval", &obj.DefaultInterval)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_time_range", &obj.DefaultTimeRange)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "value_types", &obj.ValueTypes, UnmarshalValueType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NumericCardElementValueType : NumericCardElementValueType struct
// This model "extends" ValueType
type NumericCardElementValueType struct {
	// Kind of element
	// - KPI&#58; Kind of value derived from a KPI occurrence.
	Kind *string `json:"kind,omitempty"`

	// The name of the kpi note associated to the occurrence with the value for this card element value type.
	KpiNoteName *string `json:"kpi_note_name,omitempty"`

	// The text of this element type.
	Text *string `json:"text,omitempty"`

	// the names of the finding note associated that act as filters for counting the occurrences.
	FindingNoteNames []string `json:"finding_note_names,omitempty"`
}

// Constants associated with the NumericCardElementValueType.Kind property.
// Kind of element
// - KPI&#58; Kind of value derived from a KPI occurrence.
const (
	NumericCardElementValueTypeKindKpiConst = "KPI"
)

func (*NumericCardElementValueType) isaValueType() bool {
	return true
}

// UnmarshalNumericCardElementValueType unmarshals an instance of NumericCardElementValueType from the specified map of raw messages.
func UnmarshalNumericCardElementValueType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NumericCardElementValueType)
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kpi_note_name", &obj.KpiNoteName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "finding_note_names", &obj.FindingNoteNames)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ValueTypeFindingCountValueType : ValueTypeFindingCountValueType struct
// This model "extends" ValueType
type ValueTypeFindingCountValueType struct {
	// Kind of element - FINDING_COUNT&#58; Kind of value derived from a count of finding occurrences.
	Kind *string `json:"kind" validate:"required"`

	// the names of the finding note associated that act as filters for counting the occurrences.
	FindingNoteNames []string `json:"finding_note_names" validate:"required"`

	// The text of this element type.
	Text *string `json:"text" validate:"required"`
}

// Constants associated with the ValueTypeFindingCountValueType.Kind property.
// Kind of element - FINDING_COUNT&#58; Kind of value derived from a count of finding occurrences.
const (
	ValueTypeFindingCountValueTypeKindFindingCountConst = "FINDING_COUNT"
)

// NewValueTypeFindingCountValueType : Instantiate ValueTypeFindingCountValueType (Generic Model Constructor)
func (*FindingsV1) NewValueTypeFindingCountValueType(kind string, findingNoteNames []string, text string) (_model *ValueTypeFindingCountValueType, err error) {
	_model = &ValueTypeFindingCountValueType{
		Kind:             core.StringPtr(kind),
		FindingNoteNames: findingNoteNames,
		Text:             core.StringPtr(text),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ValueTypeFindingCountValueType) isaValueType() bool {
	return true
}

// UnmarshalValueTypeFindingCountValueType unmarshals an instance of ValueTypeFindingCountValueType from the specified map of raw messages.
func UnmarshalValueTypeFindingCountValueType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ValueTypeFindingCountValueType)
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "finding_note_names", &obj.FindingNoteNames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ValueTypeKpiValueType : ValueTypeKpiValueType struct
// This model "extends" ValueType
type ValueTypeKpiValueType struct {
	// Kind of element
	// - KPI&#58; Kind of value derived from a KPI occurrence.
	Kind *string `json:"kind" validate:"required"`

	// The name of the kpi note associated to the occurrence with the value for this card element value type.
	KpiNoteName *string `json:"kpi_note_name" validate:"required"`

	// The text of this element type.
	Text *string `json:"text" validate:"required"`
}

// Constants associated with the ValueTypeKpiValueType.Kind property.
// Kind of element
// - KPI&#58; Kind of value derived from a KPI occurrence.
const (
	ValueTypeKpiValueTypeKindKpiConst = "KPI"
)

// NewValueTypeKpiValueType : Instantiate ValueTypeKpiValueType (Generic Model Constructor)
func (*FindingsV1) NewValueTypeKpiValueType(kind string, kpiNoteName string, text string) (_model *ValueTypeKpiValueType, err error) {
	_model = &ValueTypeKpiValueType{
		Kind:        core.StringPtr(kind),
		KpiNoteName: core.StringPtr(kpiNoteName),
		Text:        core.StringPtr(text),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ValueTypeKpiValueType) isaValueType() bool {
	return true
}

// UnmarshalValueTypeKpiValueType unmarshals an instance of ValueTypeKpiValueType from the specified map of raw messages.
func UnmarshalValueTypeKpiValueType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ValueTypeKpiValueType)
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kpi_note_name", &obj.KpiNoteName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
