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

// Package cloudantv1 : Operations and models for the CloudantV1 service
package cloudantv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	common "github.com/IBM/cloudant-go-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
)

// CloudantV1 : NoSQL database based on Apache CouchDB
// See: https://cloud.ibm.com/docs/services/Cloudant/
type CloudantV1 struct {
	Service *common.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "http://localhost:5984"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "cloudant"

// CloudantV1Options : Service options
type CloudantV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewCloudantV1UsingExternalConfig : constructs an instance of CloudantV1 with passed in options and external configuration.
func NewCloudantV1UsingExternalConfig(options *CloudantV1Options) (cloudant *CloudantV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = common.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	cloudant, err = NewCloudantV1(options)
	if err != nil {
		return
	}

	err = cloudant.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = cloudant.Service.SetServiceURL(options.URL)
	}
	return
}

// NewCloudantV1 : constructs an instance of CloudantV1 with passed in options.
func NewCloudantV1(options *CloudantV1Options) (service *CloudantV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
		EnableGzipCompression: true,
	}

	baseService, err := common.NewBaseService(serviceOptions)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
			return
		}
	}

	service = &CloudantV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "cloudant" suitable for processing requests.
func (cloudant *CloudantV1) Clone() *CloudantV1 {
	if core.IsNil(cloudant) {
		return nil
	}
	clone := *cloudant
	clone.Service = cloudant.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (cloudant *CloudantV1) SetServiceURL(url string) error {
	return cloudant.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (cloudant *CloudantV1) GetServiceURL() string {
	return cloudant.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (cloudant *CloudantV1) SetDefaultHeaders(headers http.Header) {
	cloudant.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (cloudant *CloudantV1) SetEnableGzipCompression(enableGzip bool) {
	cloudant.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (cloudant *CloudantV1) GetEnableGzipCompression() bool {
	return cloudant.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (cloudant *CloudantV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	cloudant.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (cloudant *CloudantV1) DisableRetries() {
	cloudant.Service.DisableRetries()
}

// GetServerInformation : Retrieve server instance information
// When you access the root of an instance, IBM Cloudant returns meta-information about the instance. The response
// includes a JSON structure that contains information about the server, including a welcome message and the server's
// version.
func (cloudant *CloudantV1) GetServerInformation(getServerInformationOptions *GetServerInformationOptions) (result *ServerInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetServerInformationWithContext(context.Background(), getServerInformationOptions)
}

// GetServerInformationWithContext is an alternate form of the GetServerInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetServerInformationWithContext(ctx context.Context, getServerInformationOptions *GetServerInformationOptions) (result *ServerInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getServerInformationOptions, "getServerInformationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getServerInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetServerInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServerInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetMembershipInformation : Retrieve cluster membership information
// Displays the nodes that are part of the cluster as `cluster_nodes`. The field, `all_nodes`, displays all nodes this
// node knows about, including the ones that are part of the cluster. This endpoint is useful when you set up a cluster.
func (cloudant *CloudantV1) GetMembershipInformation(getMembershipInformationOptions *GetMembershipInformationOptions) (result *MembershipInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetMembershipInformationWithContext(context.Background(), getMembershipInformationOptions)
}

// GetMembershipInformationWithContext is an alternate form of the GetMembershipInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetMembershipInformationWithContext(ctx context.Context, getMembershipInformationOptions *GetMembershipInformationOptions) (result *MembershipInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getMembershipInformationOptions, "getMembershipInformationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_membership`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getMembershipInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetMembershipInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMembershipInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetUuids : Retrieve one or more UUIDs
// Requests one or more Universally Unique Identifiers (UUIDs) from the instance. The response is a JSON object that
// provides a list of UUIDs.
func (cloudant *CloudantV1) GetUuids(getUuidsOptions *GetUuidsOptions) (result *UuidsResult, response *core.DetailedResponse, err error) {
	return cloudant.GetUuidsWithContext(context.Background(), getUuidsOptions)
}

// GetUuidsWithContext is an alternate form of the GetUuids method which supports a Context parameter
func (cloudant *CloudantV1) GetUuidsWithContext(ctx context.Context, getUuidsOptions *GetUuidsOptions) (result *UuidsResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getUuidsOptions, "getUuidsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_uuids`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getUuidsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetUuids")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getUuidsOptions.Count != nil {
		builder.AddQuery("count", fmt.Sprint(*getUuidsOptions.Count))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUuidsResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCapacityThroughputInformation : Retrieve provisioned throughput capacity information
// View the amount of provisioned throughput capacity that is allocated to an IBM Cloudant instance and what is the
// target provisioned throughput capacity.
func (cloudant *CloudantV1) GetCapacityThroughputInformation(getCapacityThroughputInformationOptions *GetCapacityThroughputInformationOptions) (result *CapacityThroughputInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetCapacityThroughputInformationWithContext(context.Background(), getCapacityThroughputInformationOptions)
}

// GetCapacityThroughputInformationWithContext is an alternate form of the GetCapacityThroughputInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetCapacityThroughputInformationWithContext(ctx context.Context, getCapacityThroughputInformationOptions *GetCapacityThroughputInformationOptions) (result *CapacityThroughputInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getCapacityThroughputInformationOptions, "getCapacityThroughputInformationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/user/capacity/throughput`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCapacityThroughputInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetCapacityThroughputInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCapacityThroughputInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutCapacityThroughputConfiguration : Update the target provisioned throughput capacity
// Sets the target provisioned throughput capacity for an IBM Cloudant instance. When target capacity is changed, the
// current capacity asynchronously changes to meet the target capacity.
func (cloudant *CloudantV1) PutCapacityThroughputConfiguration(putCapacityThroughputConfigurationOptions *PutCapacityThroughputConfigurationOptions) (result *CapacityThroughputInformation, response *core.DetailedResponse, err error) {
	return cloudant.PutCapacityThroughputConfigurationWithContext(context.Background(), putCapacityThroughputConfigurationOptions)
}

// PutCapacityThroughputConfigurationWithContext is an alternate form of the PutCapacityThroughputConfiguration method which supports a Context parameter
func (cloudant *CloudantV1) PutCapacityThroughputConfigurationWithContext(ctx context.Context, putCapacityThroughputConfigurationOptions *PutCapacityThroughputConfigurationOptions) (result *CapacityThroughputInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putCapacityThroughputConfigurationOptions, "putCapacityThroughputConfigurationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putCapacityThroughputConfigurationOptions, "putCapacityThroughputConfigurationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/user/capacity/throughput`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range putCapacityThroughputConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutCapacityThroughputConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if putCapacityThroughputConfigurationOptions.Blocks != nil {
		body["blocks"] = putCapacityThroughputConfigurationOptions.Blocks
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCapacityThroughputInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetDbUpdates : Retrieve change events for all databases
// Lists changes to databases, like a global changes feed. Types of changes include updating the database and creating
// or deleting a database. Like the changes feed, the feed is not guaranteed to return changes in the correct order and
// might repeat changes. Polling modes for this method work like polling modes for the changes feed.
// **Note: This endpoint requires _admin or _db_updates role and is only available on dedicated clusters.**.
func (cloudant *CloudantV1) GetDbUpdates(getDbUpdatesOptions *GetDbUpdatesOptions) (result *DbUpdates, response *core.DetailedResponse, err error) {
	return cloudant.GetDbUpdatesWithContext(context.Background(), getDbUpdatesOptions)
}

// GetDbUpdatesWithContext is an alternate form of the GetDbUpdates method which supports a Context parameter
func (cloudant *CloudantV1) GetDbUpdatesWithContext(ctx context.Context, getDbUpdatesOptions *GetDbUpdatesOptions) (result *DbUpdates, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getDbUpdatesOptions, "getDbUpdatesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_db_updates`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDbUpdatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDbUpdates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getDbUpdatesOptions.Feed != nil {
		builder.AddQuery("feed", fmt.Sprint(*getDbUpdatesOptions.Feed))
	}
	if getDbUpdatesOptions.Heartbeat != nil {
		builder.AddQuery("heartbeat", fmt.Sprint(*getDbUpdatesOptions.Heartbeat))
	}
	if getDbUpdatesOptions.Timeout != nil {
		builder.AddQuery("timeout", fmt.Sprint(*getDbUpdatesOptions.Timeout))
	}
	if getDbUpdatesOptions.Since != nil {
		builder.AddQuery("since", fmt.Sprint(*getDbUpdatesOptions.Since))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDbUpdates)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostChanges : Query the database document changes feed
// Requests the database changes feed in the same way as `GET /{db}/_changes` does. It is widely used with the `filter`
// query parameter because it allows one to pass more information to the filter.
func (cloudant *CloudantV1) PostChanges(postChangesOptions *PostChangesOptions) (result *ChangesResult, response *core.DetailedResponse, err error) {
	return cloudant.PostChangesWithContext(context.Background(), postChangesOptions)
}

// PostChangesWithContext is an alternate form of the PostChanges method which supports a Context parameter
func (cloudant *CloudantV1) PostChangesWithContext(ctx context.Context, postChangesOptions *PostChangesOptions) (result *ChangesResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postChangesOptions, "postChangesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postChangesOptions, "postChangesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postChangesOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_changes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postChangesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostChanges")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if postChangesOptions.LastEventID != nil {
		builder.AddHeader("Last-Event-ID", fmt.Sprint(*postChangesOptions.LastEventID))
	}

	if postChangesOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*postChangesOptions.AttEncodingInfo))
	}
	if postChangesOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*postChangesOptions.Attachments))
	}
	if postChangesOptions.Conflicts != nil {
		builder.AddQuery("conflicts", fmt.Sprint(*postChangesOptions.Conflicts))
	}
	if postChangesOptions.Descending != nil {
		builder.AddQuery("descending", fmt.Sprint(*postChangesOptions.Descending))
	}
	if postChangesOptions.Feed != nil {
		builder.AddQuery("feed", fmt.Sprint(*postChangesOptions.Feed))
	}
	if postChangesOptions.Filter != nil {
		builder.AddQuery("filter", fmt.Sprint(*postChangesOptions.Filter))
	}
	if postChangesOptions.Heartbeat != nil {
		builder.AddQuery("heartbeat", fmt.Sprint(*postChangesOptions.Heartbeat))
	}
	if postChangesOptions.IncludeDocs != nil {
		builder.AddQuery("include_docs", fmt.Sprint(*postChangesOptions.IncludeDocs))
	}
	if postChangesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*postChangesOptions.Limit))
	}
	if postChangesOptions.SeqInterval != nil {
		builder.AddQuery("seq_interval", fmt.Sprint(*postChangesOptions.SeqInterval))
	}
	if postChangesOptions.Since != nil {
		builder.AddQuery("since", fmt.Sprint(*postChangesOptions.Since))
	}
	if postChangesOptions.Style != nil {
		builder.AddQuery("style", fmt.Sprint(*postChangesOptions.Style))
	}
	if postChangesOptions.Timeout != nil {
		builder.AddQuery("timeout", fmt.Sprint(*postChangesOptions.Timeout))
	}
	if postChangesOptions.View != nil {
		builder.AddQuery("view", fmt.Sprint(*postChangesOptions.View))
	}

	body := make(map[string]interface{})
	if postChangesOptions.DocIds != nil {
		body["doc_ids"] = postChangesOptions.DocIds
	}
	if postChangesOptions.Fields != nil {
		body["fields"] = postChangesOptions.Fields
	}
	if postChangesOptions.Selector != nil {
		body["selector"] = postChangesOptions.Selector
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalChangesResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostChangesAsStream : Query the database document changes feed as stream
// Requests the database changes feed in the same way as `GET /{db}/_changes` does. It is widely used with the `filter`
// query parameter because it allows one to pass more information to the filter.
func (cloudant *CloudantV1) PostChangesAsStream(postChangesOptions *PostChangesOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostChangesAsStreamWithContext(context.Background(), postChangesOptions)
}

// PostChangesAsStreamWithContext is an alternate form of the PostChangesAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostChangesAsStreamWithContext(ctx context.Context, postChangesOptions *PostChangesOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postChangesOptions, "postChangesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postChangesOptions, "postChangesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postChangesOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_changes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postChangesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostChangesAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if postChangesOptions.LastEventID != nil {
		builder.AddHeader("Last-Event-ID", fmt.Sprint(*postChangesOptions.LastEventID))
	}

	if postChangesOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*postChangesOptions.AttEncodingInfo))
	}
	if postChangesOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*postChangesOptions.Attachments))
	}
	if postChangesOptions.Conflicts != nil {
		builder.AddQuery("conflicts", fmt.Sprint(*postChangesOptions.Conflicts))
	}
	if postChangesOptions.Descending != nil {
		builder.AddQuery("descending", fmt.Sprint(*postChangesOptions.Descending))
	}
	if postChangesOptions.Feed != nil {
		builder.AddQuery("feed", fmt.Sprint(*postChangesOptions.Feed))
	}
	if postChangesOptions.Filter != nil {
		builder.AddQuery("filter", fmt.Sprint(*postChangesOptions.Filter))
	}
	if postChangesOptions.Heartbeat != nil {
		builder.AddQuery("heartbeat", fmt.Sprint(*postChangesOptions.Heartbeat))
	}
	if postChangesOptions.IncludeDocs != nil {
		builder.AddQuery("include_docs", fmt.Sprint(*postChangesOptions.IncludeDocs))
	}
	if postChangesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*postChangesOptions.Limit))
	}
	if postChangesOptions.SeqInterval != nil {
		builder.AddQuery("seq_interval", fmt.Sprint(*postChangesOptions.SeqInterval))
	}
	if postChangesOptions.Since != nil {
		builder.AddQuery("since", fmt.Sprint(*postChangesOptions.Since))
	}
	if postChangesOptions.Style != nil {
		builder.AddQuery("style", fmt.Sprint(*postChangesOptions.Style))
	}
	if postChangesOptions.Timeout != nil {
		builder.AddQuery("timeout", fmt.Sprint(*postChangesOptions.Timeout))
	}
	if postChangesOptions.View != nil {
		builder.AddQuery("view", fmt.Sprint(*postChangesOptions.View))
	}

	body := make(map[string]interface{})
	if postChangesOptions.DocIds != nil {
		body["doc_ids"] = postChangesOptions.DocIds
	}
	if postChangesOptions.Fields != nil {
		body["fields"] = postChangesOptions.Fields
	}
	if postChangesOptions.Selector != nil {
		body["selector"] = postChangesOptions.Selector
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// HeadDatabase : Retrieve the HTTP headers for a database
// Returns the HTTP headers that contain a minimal amount of information about the specified database. Since the
// response body is empty, using the HEAD method is a lightweight way to check if the database exists or not.
func (cloudant *CloudantV1) HeadDatabase(headDatabaseOptions *HeadDatabaseOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadDatabaseWithContext(context.Background(), headDatabaseOptions)
}

// HeadDatabaseWithContext is an alternate form of the HeadDatabase method which supports a Context parameter
func (cloudant *CloudantV1) HeadDatabaseWithContext(ctx context.Context, headDatabaseOptions *HeadDatabaseOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(headDatabaseOptions, "headDatabaseOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(headDatabaseOptions, "headDatabaseOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *headDatabaseOptions.Db,
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range headDatabaseOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadDatabase")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// GetAllDbs : Query a list of all database names in the instance
func (cloudant *CloudantV1) GetAllDbs(getAllDbsOptions *GetAllDbsOptions) (result []string, response *core.DetailedResponse, err error) {
	return cloudant.GetAllDbsWithContext(context.Background(), getAllDbsOptions)
}

// GetAllDbsWithContext is an alternate form of the GetAllDbs method which supports a Context parameter
func (cloudant *CloudantV1) GetAllDbsWithContext(ctx context.Context, getAllDbsOptions *GetAllDbsOptions) (result []string, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getAllDbsOptions, "getAllDbsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_all_dbs`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAllDbsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetAllDbs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getAllDbsOptions.Descending != nil {
		builder.AddQuery("descending", fmt.Sprint(*getAllDbsOptions.Descending))
	}
	if getAllDbsOptions.Endkey != nil {
		builder.AddQuery("endkey", fmt.Sprint(*getAllDbsOptions.Endkey))
	}
	if getAllDbsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getAllDbsOptions.Limit))
	}
	if getAllDbsOptions.Skip != nil {
		builder.AddQuery("skip", fmt.Sprint(*getAllDbsOptions.Skip))
	}
	if getAllDbsOptions.Startkey != nil {
		builder.AddQuery("startkey", fmt.Sprint(*getAllDbsOptions.Startkey))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostDbsInfo : Query information about multiple databases
// This operation enables you to request information about multiple databases in a single request, instead of issuing
// multiple `GET /{db}` requests. It returns a list that contains an information object for each database specified in
// the request.
func (cloudant *CloudantV1) PostDbsInfo(postDbsInfoOptions *PostDbsInfoOptions) (result []DbsInfoResult, response *core.DetailedResponse, err error) {
	return cloudant.PostDbsInfoWithContext(context.Background(), postDbsInfoOptions)
}

// PostDbsInfoWithContext is an alternate form of the PostDbsInfo method which supports a Context parameter
func (cloudant *CloudantV1) PostDbsInfoWithContext(ctx context.Context, postDbsInfoOptions *PostDbsInfoOptions) (result []DbsInfoResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDbsInfoOptions, "postDbsInfoOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postDbsInfoOptions, "postDbsInfoOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_dbs_info`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range postDbsInfoOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostDbsInfo")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postDbsInfoOptions.Keys != nil {
		body["keys"] = postDbsInfoOptions.Keys
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDbsInfoResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteDatabase : Delete a database
// Deletes the specified database and all documents and attachments contained within it. To avoid deleting a database,
// the server responds with a 400 HTTP status code when the request URL includes a `?rev=` parameter. This response
// suggests that a user wanted to delete a document but forgot to add the document ID to the URL.
func (cloudant *CloudantV1) DeleteDatabase(deleteDatabaseOptions *DeleteDatabaseOptions) (result *Ok, response *core.DetailedResponse, err error) {
	return cloudant.DeleteDatabaseWithContext(context.Background(), deleteDatabaseOptions)
}

// DeleteDatabaseWithContext is an alternate form of the DeleteDatabase method which supports a Context parameter
func (cloudant *CloudantV1) DeleteDatabaseWithContext(ctx context.Context, deleteDatabaseOptions *DeleteDatabaseOptions) (result *Ok, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDatabaseOptions, "deleteDatabaseOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteDatabaseOptions, "deleteDatabaseOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *deleteDatabaseOptions.Db,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteDatabaseOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "DeleteDatabase")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetDatabaseInformation : Retrieve information about a database
func (cloudant *CloudantV1) GetDatabaseInformation(getDatabaseInformationOptions *GetDatabaseInformationOptions) (result *DatabaseInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetDatabaseInformationWithContext(context.Background(), getDatabaseInformationOptions)
}

// GetDatabaseInformationWithContext is an alternate form of the GetDatabaseInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetDatabaseInformationWithContext(ctx context.Context, getDatabaseInformationOptions *GetDatabaseInformationOptions) (result *DatabaseInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDatabaseInformationOptions, "getDatabaseInformationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDatabaseInformationOptions, "getDatabaseInformationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getDatabaseInformationOptions.Db,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDatabaseInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDatabaseInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDatabaseInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutDatabase : Create a database
func (cloudant *CloudantV1) PutDatabase(putDatabaseOptions *PutDatabaseOptions) (result *Ok, response *core.DetailedResponse, err error) {
	return cloudant.PutDatabaseWithContext(context.Background(), putDatabaseOptions)
}

// PutDatabaseWithContext is an alternate form of the PutDatabase method which supports a Context parameter
func (cloudant *CloudantV1) PutDatabaseWithContext(ctx context.Context, putDatabaseOptions *PutDatabaseOptions) (result *Ok, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putDatabaseOptions, "putDatabaseOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putDatabaseOptions, "putDatabaseOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *putDatabaseOptions.Db,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putDatabaseOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutDatabase")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if putDatabaseOptions.Partitioned != nil {
		builder.AddQuery("partitioned", fmt.Sprint(*putDatabaseOptions.Partitioned))
	}
	if putDatabaseOptions.Q != nil {
		builder.AddQuery("q", fmt.Sprint(*putDatabaseOptions.Q))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// HeadDocument : Retrieve the HTTP headers for the document
// This method supports the same query arguments as the `GET /{db}/{docid}` method, but only the header information
// (including document size and the revision as an ETag) is returned. The ETag header shows the current revision for the
// requested document, and the Content-Length specifies the length of the data if the document was requested in full.
// Add any of the query arguments, then the resulting HTTP headers that correspond to it are returned.
func (cloudant *CloudantV1) HeadDocument(headDocumentOptions *HeadDocumentOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadDocumentWithContext(context.Background(), headDocumentOptions)
}

// HeadDocumentWithContext is an alternate form of the HeadDocument method which supports a Context parameter
func (cloudant *CloudantV1) HeadDocumentWithContext(ctx context.Context, headDocumentOptions *HeadDocumentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(headDocumentOptions, "headDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(headDocumentOptions, "headDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *headDocumentOptions.Db,
		"doc_id": *headDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range headDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if headDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*headDocumentOptions.IfNoneMatch))
	}

	if headDocumentOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*headDocumentOptions.Latest))
	}
	if headDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*headDocumentOptions.Rev))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// PostDocument : Create or modify a document in a database
// Creates or modifies a document in the specified database by using the supplied JSON document. If the JSON document
// doesn't specify an `_id` field, then the document is created with a new unique ID generated by the UUID algorithm
// that is configured for the server. If the document includes the `_id` field, then it is created with that `_id` or
// updated if the `_id` already exists, and an appropriate `_rev` is included in the JSON document. If the `_id`
// includes the `_local` or `_design` prefix, then this operation is used to create or modify local or design documents
// respectively.
func (cloudant *CloudantV1) PostDocument(postDocumentOptions *PostDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.PostDocumentWithContext(context.Background(), postDocumentOptions)
}

// PostDocumentWithContext is an alternate form of the PostDocument method which supports a Context parameter
func (cloudant *CloudantV1) PostDocumentWithContext(ctx context.Context, postDocumentOptions *PostDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDocumentOptions, "postDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postDocumentOptions, "postDocumentOptions")
	if err != nil {
		return
	}

	if postDocumentOptions.Document != nil && postDocumentOptions.ContentType == nil {
		postDocumentOptions.SetContentType("application/json")
	}

	pathParamsMap := map[string]string{
		"db": *postDocumentOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if postDocumentOptions.ContentType != nil {
		builder.AddHeader("Content-Type", fmt.Sprint(*postDocumentOptions.ContentType))
	}

	if postDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*postDocumentOptions.Batch))
	}

	_, err = builder.SetBodyContent(core.StringNilMapper(postDocumentOptions.ContentType), postDocumentOptions.Document, nil, postDocumentOptions.Body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostAllDocs : Query a list of all documents in a database
// Queries the primary index (all document IDs). The results that match the request body parameters are returned in a
// JSON object, including a list of matching documents with basic contents, such as the ID and revision. When no request
// body parameters are specified, results for all documents in the database are returned. Optionally, document content
// or additional metadata can be included in the response.
func (cloudant *CloudantV1) PostAllDocs(postAllDocsOptions *PostAllDocsOptions) (result *AllDocsResult, response *core.DetailedResponse, err error) {
	return cloudant.PostAllDocsWithContext(context.Background(), postAllDocsOptions)
}

// PostAllDocsWithContext is an alternate form of the PostAllDocs method which supports a Context parameter
func (cloudant *CloudantV1) PostAllDocsWithContext(ctx context.Context, postAllDocsOptions *PostAllDocsOptions) (result *AllDocsResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postAllDocsOptions, "postAllDocsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postAllDocsOptions, "postAllDocsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postAllDocsOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_all_docs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postAllDocsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostAllDocs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postAllDocsOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postAllDocsOptions.AttEncodingInfo
	}
	if postAllDocsOptions.Attachments != nil {
		body["attachments"] = postAllDocsOptions.Attachments
	}
	if postAllDocsOptions.Conflicts != nil {
		body["conflicts"] = postAllDocsOptions.Conflicts
	}
	if postAllDocsOptions.Descending != nil {
		body["descending"] = postAllDocsOptions.Descending
	}
	if postAllDocsOptions.IncludeDocs != nil {
		body["include_docs"] = postAllDocsOptions.IncludeDocs
	}
	if postAllDocsOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postAllDocsOptions.InclusiveEnd
	}
	if postAllDocsOptions.Limit != nil {
		body["limit"] = postAllDocsOptions.Limit
	}
	if postAllDocsOptions.Skip != nil {
		body["skip"] = postAllDocsOptions.Skip
	}
	if postAllDocsOptions.UpdateSeq != nil {
		body["update_seq"] = postAllDocsOptions.UpdateSeq
	}
	if postAllDocsOptions.Endkey != nil {
		body["endkey"] = postAllDocsOptions.Endkey
	}
	if postAllDocsOptions.Key != nil {
		body["key"] = postAllDocsOptions.Key
	}
	if postAllDocsOptions.Keys != nil {
		body["keys"] = postAllDocsOptions.Keys
	}
	if postAllDocsOptions.Startkey != nil {
		body["startkey"] = postAllDocsOptions.Startkey
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllDocsResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostAllDocsAsStream : Query a list of all documents in a database as stream
// Queries the primary index (all document IDs). The results that match the request body parameters are returned in a
// JSON object, including a list of matching documents with basic contents, such as the ID and revision. When no request
// body parameters are specified, results for all documents in the database are returned. Optionally, document content
// or additional metadata can be included in the response.
func (cloudant *CloudantV1) PostAllDocsAsStream(postAllDocsOptions *PostAllDocsOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostAllDocsAsStreamWithContext(context.Background(), postAllDocsOptions)
}

// PostAllDocsAsStreamWithContext is an alternate form of the PostAllDocsAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostAllDocsAsStreamWithContext(ctx context.Context, postAllDocsOptions *PostAllDocsOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postAllDocsOptions, "postAllDocsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postAllDocsOptions, "postAllDocsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postAllDocsOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_all_docs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postAllDocsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostAllDocsAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postAllDocsOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postAllDocsOptions.AttEncodingInfo
	}
	if postAllDocsOptions.Attachments != nil {
		body["attachments"] = postAllDocsOptions.Attachments
	}
	if postAllDocsOptions.Conflicts != nil {
		body["conflicts"] = postAllDocsOptions.Conflicts
	}
	if postAllDocsOptions.Descending != nil {
		body["descending"] = postAllDocsOptions.Descending
	}
	if postAllDocsOptions.IncludeDocs != nil {
		body["include_docs"] = postAllDocsOptions.IncludeDocs
	}
	if postAllDocsOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postAllDocsOptions.InclusiveEnd
	}
	if postAllDocsOptions.Limit != nil {
		body["limit"] = postAllDocsOptions.Limit
	}
	if postAllDocsOptions.Skip != nil {
		body["skip"] = postAllDocsOptions.Skip
	}
	if postAllDocsOptions.UpdateSeq != nil {
		body["update_seq"] = postAllDocsOptions.UpdateSeq
	}
	if postAllDocsOptions.Endkey != nil {
		body["endkey"] = postAllDocsOptions.Endkey
	}
	if postAllDocsOptions.Key != nil {
		body["key"] = postAllDocsOptions.Key
	}
	if postAllDocsOptions.Keys != nil {
		body["keys"] = postAllDocsOptions.Keys
	}
	if postAllDocsOptions.Startkey != nil {
		body["startkey"] = postAllDocsOptions.Startkey
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostAllDocsQueries : Multi-query the list of all documents in a database
// Runs multiple queries using the primary index (all document IDs). Returns a JSON object that contains a list of
// result objects, one for each query, with a structure equivalent to that of a single `_all_docs` request. This enables
// you to request multiple queries in a single request, in place of multiple `POST /{db}/_all_docs` requests.
func (cloudant *CloudantV1) PostAllDocsQueries(postAllDocsQueriesOptions *PostAllDocsQueriesOptions) (result *AllDocsQueriesResult, response *core.DetailedResponse, err error) {
	return cloudant.PostAllDocsQueriesWithContext(context.Background(), postAllDocsQueriesOptions)
}

// PostAllDocsQueriesWithContext is an alternate form of the PostAllDocsQueries method which supports a Context parameter
func (cloudant *CloudantV1) PostAllDocsQueriesWithContext(ctx context.Context, postAllDocsQueriesOptions *PostAllDocsQueriesOptions) (result *AllDocsQueriesResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postAllDocsQueriesOptions, "postAllDocsQueriesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postAllDocsQueriesOptions, "postAllDocsQueriesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postAllDocsQueriesOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_all_docs/queries`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postAllDocsQueriesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostAllDocsQueries")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postAllDocsQueriesOptions.Queries != nil {
		body["queries"] = postAllDocsQueriesOptions.Queries
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllDocsQueriesResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostAllDocsQueriesAsStream : Multi-query the list of all documents in a database as stream
// Runs multiple queries using the primary index (all document IDs). Returns a JSON object that contains a list of
// result objects, one for each query, with a structure equivalent to that of a single `_all_docs` request. This enables
// you to request multiple queries in a single request, in place of multiple `POST /{db}/_all_docs` requests.
func (cloudant *CloudantV1) PostAllDocsQueriesAsStream(postAllDocsQueriesOptions *PostAllDocsQueriesOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostAllDocsQueriesAsStreamWithContext(context.Background(), postAllDocsQueriesOptions)
}

// PostAllDocsQueriesAsStreamWithContext is an alternate form of the PostAllDocsQueriesAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostAllDocsQueriesAsStreamWithContext(ctx context.Context, postAllDocsQueriesOptions *PostAllDocsQueriesOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postAllDocsQueriesOptions, "postAllDocsQueriesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postAllDocsQueriesOptions, "postAllDocsQueriesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postAllDocsQueriesOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_all_docs/queries`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postAllDocsQueriesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostAllDocsQueriesAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postAllDocsQueriesOptions.Queries != nil {
		body["queries"] = postAllDocsQueriesOptions.Queries
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostBulkDocs : Bulk modify multiple documents in a database
// The bulk document API allows you to create and update multiple documents at the same time within a single request.
// The basic operation is similar to creating or updating a single document, except that you batch the document
// structure and information.
func (cloudant *CloudantV1) PostBulkDocs(postBulkDocsOptions *PostBulkDocsOptions) (result []DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.PostBulkDocsWithContext(context.Background(), postBulkDocsOptions)
}

// PostBulkDocsWithContext is an alternate form of the PostBulkDocs method which supports a Context parameter
func (cloudant *CloudantV1) PostBulkDocsWithContext(ctx context.Context, postBulkDocsOptions *PostBulkDocsOptions) (result []DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postBulkDocsOptions, "postBulkDocsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postBulkDocsOptions, "postBulkDocsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postBulkDocsOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_bulk_docs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postBulkDocsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostBulkDocs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContent("application/json", postBulkDocsOptions.BulkDocs, nil, postBulkDocsOptions.Body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostBulkGet : Bulk query revision information for multiple documents
// Fetch specific revisions or revision histories for multiple documents in bulk as replicators do.
func (cloudant *CloudantV1) PostBulkGet(postBulkGetOptions *PostBulkGetOptions) (result *BulkGetResult, response *core.DetailedResponse, err error) {
	return cloudant.PostBulkGetWithContext(context.Background(), postBulkGetOptions)
}

// PostBulkGetWithContext is an alternate form of the PostBulkGet method which supports a Context parameter
func (cloudant *CloudantV1) PostBulkGetWithContext(ctx context.Context, postBulkGetOptions *PostBulkGetOptions) (result *BulkGetResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postBulkGetOptions, "postBulkGetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postBulkGetOptions, "postBulkGetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postBulkGetOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_bulk_get`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postBulkGetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostBulkGet")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if postBulkGetOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*postBulkGetOptions.Attachments))
	}
	if postBulkGetOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*postBulkGetOptions.AttEncodingInfo))
	}
	if postBulkGetOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*postBulkGetOptions.Latest))
	}
	if postBulkGetOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*postBulkGetOptions.Revs))
	}

	body := make(map[string]interface{})
	if postBulkGetOptions.Docs != nil {
		body["docs"] = postBulkGetOptions.Docs
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBulkGetResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostBulkGetAsMixed : Bulk query revision information for multiple documents as mixed
// Fetch specific revisions or revision histories for multiple documents in bulk as replicators do.
func (cloudant *CloudantV1) PostBulkGetAsMixed(postBulkGetOptions *PostBulkGetOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostBulkGetAsMixedWithContext(context.Background(), postBulkGetOptions)
}

// PostBulkGetAsMixedWithContext is an alternate form of the PostBulkGetAsMixed method which supports a Context parameter
func (cloudant *CloudantV1) PostBulkGetAsMixedWithContext(ctx context.Context, postBulkGetOptions *PostBulkGetOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postBulkGetOptions, "postBulkGetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postBulkGetOptions, "postBulkGetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postBulkGetOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_bulk_get`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postBulkGetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostBulkGetAsMixed")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "multipart/mixed")
	builder.AddHeader("Content-Type", "application/json")

	if postBulkGetOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*postBulkGetOptions.Attachments))
	}
	if postBulkGetOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*postBulkGetOptions.AttEncodingInfo))
	}
	if postBulkGetOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*postBulkGetOptions.Latest))
	}
	if postBulkGetOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*postBulkGetOptions.Revs))
	}

	body := make(map[string]interface{})
	if postBulkGetOptions.Docs != nil {
		body["docs"] = postBulkGetOptions.Docs
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostBulkGetAsRelated : Bulk query revision information for multiple documents as related
// Fetch specific revisions or revision histories for multiple documents in bulk as replicators do.
func (cloudant *CloudantV1) PostBulkGetAsRelated(postBulkGetOptions *PostBulkGetOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostBulkGetAsRelatedWithContext(context.Background(), postBulkGetOptions)
}

// PostBulkGetAsRelatedWithContext is an alternate form of the PostBulkGetAsRelated method which supports a Context parameter
func (cloudant *CloudantV1) PostBulkGetAsRelatedWithContext(ctx context.Context, postBulkGetOptions *PostBulkGetOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postBulkGetOptions, "postBulkGetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postBulkGetOptions, "postBulkGetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postBulkGetOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_bulk_get`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postBulkGetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostBulkGetAsRelated")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "multipart/related")
	builder.AddHeader("Content-Type", "application/json")

	if postBulkGetOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*postBulkGetOptions.Attachments))
	}
	if postBulkGetOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*postBulkGetOptions.AttEncodingInfo))
	}
	if postBulkGetOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*postBulkGetOptions.Latest))
	}
	if postBulkGetOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*postBulkGetOptions.Revs))
	}

	body := make(map[string]interface{})
	if postBulkGetOptions.Docs != nil {
		body["docs"] = postBulkGetOptions.Docs
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostBulkGetAsStream : Bulk query revision information for multiple documents as stream
// Fetch specific revisions or revision histories for multiple documents in bulk as replicators do.
func (cloudant *CloudantV1) PostBulkGetAsStream(postBulkGetOptions *PostBulkGetOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostBulkGetAsStreamWithContext(context.Background(), postBulkGetOptions)
}

// PostBulkGetAsStreamWithContext is an alternate form of the PostBulkGetAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostBulkGetAsStreamWithContext(ctx context.Context, postBulkGetOptions *PostBulkGetOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postBulkGetOptions, "postBulkGetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postBulkGetOptions, "postBulkGetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postBulkGetOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_bulk_get`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postBulkGetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostBulkGetAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if postBulkGetOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*postBulkGetOptions.Attachments))
	}
	if postBulkGetOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*postBulkGetOptions.AttEncodingInfo))
	}
	if postBulkGetOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*postBulkGetOptions.Latest))
	}
	if postBulkGetOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*postBulkGetOptions.Revs))
	}

	body := make(map[string]interface{})
	if postBulkGetOptions.Docs != nil {
		body["docs"] = postBulkGetOptions.Docs
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// DeleteDocument : Delete a document
// Marks the specified document as deleted by adding a `_deleted` field with the value `true`. Documents with this field
// are not returned within requests anymore but stay in the database. You must supply the current (latest) revision,
// either by using the `rev` parameter or by using the `If-Match` header to specify the revision.
func (cloudant *CloudantV1) DeleteDocument(deleteDocumentOptions *DeleteDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.DeleteDocumentWithContext(context.Background(), deleteDocumentOptions)
}

// DeleteDocumentWithContext is an alternate form of the DeleteDocument method which supports a Context parameter
func (cloudant *CloudantV1) DeleteDocumentWithContext(ctx context.Context, deleteDocumentOptions *DeleteDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDocumentOptions, "deleteDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteDocumentOptions, "deleteDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *deleteDocumentOptions.Db,
		"doc_id": *deleteDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "DeleteDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteDocumentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deleteDocumentOptions.IfMatch))
	}

	if deleteDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*deleteDocumentOptions.Batch))
	}
	if deleteDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*deleteDocumentOptions.Rev))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetDocument : Retrieve a document
// Returns document with the specified `doc_id` from the specified database. Unless you request a specific revision, the
// latest revision of the document is always returned.
func (cloudant *CloudantV1) GetDocument(getDocumentOptions *GetDocumentOptions) (result *Document, response *core.DetailedResponse, err error) {
	return cloudant.GetDocumentWithContext(context.Background(), getDocumentOptions)
}

// GetDocumentWithContext is an alternate form of the GetDocument method which supports a Context parameter
func (cloudant *CloudantV1) GetDocumentWithContext(ctx context.Context, getDocumentOptions *GetDocumentOptions) (result *Document, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDocumentOptions, "getDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDocumentOptions, "getDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getDocumentOptions.Db,
		"doc_id": *getDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*getDocumentOptions.IfNoneMatch))
	}

	if getDocumentOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*getDocumentOptions.Attachments))
	}
	if getDocumentOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*getDocumentOptions.AttEncodingInfo))
	}
	if getDocumentOptions.Conflicts != nil {
		builder.AddQuery("conflicts", fmt.Sprint(*getDocumentOptions.Conflicts))
	}
	if getDocumentOptions.DeletedConflicts != nil {
		builder.AddQuery("deleted_conflicts", fmt.Sprint(*getDocumentOptions.DeletedConflicts))
	}
	if getDocumentOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*getDocumentOptions.Latest))
	}
	if getDocumentOptions.LocalSeq != nil {
		builder.AddQuery("local_seq", fmt.Sprint(*getDocumentOptions.LocalSeq))
	}
	if getDocumentOptions.Meta != nil {
		builder.AddQuery("meta", fmt.Sprint(*getDocumentOptions.Meta))
	}
	if getDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*getDocumentOptions.Rev))
	}
	if getDocumentOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*getDocumentOptions.Revs))
	}
	if getDocumentOptions.RevsInfo != nil {
		builder.AddQuery("revs_info", fmt.Sprint(*getDocumentOptions.RevsInfo))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocument)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetDocumentAsMixed : Retrieve a document as mixed
// Returns document with the specified `doc_id` from the specified database. Unless you request a specific revision, the
// latest revision of the document is always returned.
func (cloudant *CloudantV1) GetDocumentAsMixed(getDocumentOptions *GetDocumentOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.GetDocumentAsMixedWithContext(context.Background(), getDocumentOptions)
}

// GetDocumentAsMixedWithContext is an alternate form of the GetDocumentAsMixed method which supports a Context parameter
func (cloudant *CloudantV1) GetDocumentAsMixedWithContext(ctx context.Context, getDocumentOptions *GetDocumentOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDocumentOptions, "getDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDocumentOptions, "getDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getDocumentOptions.Db,
		"doc_id": *getDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDocumentAsMixed")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "multipart/mixed")
	if getDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*getDocumentOptions.IfNoneMatch))
	}

	if getDocumentOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*getDocumentOptions.Attachments))
	}
	if getDocumentOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*getDocumentOptions.AttEncodingInfo))
	}
	if getDocumentOptions.Conflicts != nil {
		builder.AddQuery("conflicts", fmt.Sprint(*getDocumentOptions.Conflicts))
	}
	if getDocumentOptions.DeletedConflicts != nil {
		builder.AddQuery("deleted_conflicts", fmt.Sprint(*getDocumentOptions.DeletedConflicts))
	}
	if getDocumentOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*getDocumentOptions.Latest))
	}
	if getDocumentOptions.LocalSeq != nil {
		builder.AddQuery("local_seq", fmt.Sprint(*getDocumentOptions.LocalSeq))
	}
	if getDocumentOptions.Meta != nil {
		builder.AddQuery("meta", fmt.Sprint(*getDocumentOptions.Meta))
	}
	if getDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*getDocumentOptions.Rev))
	}
	if getDocumentOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*getDocumentOptions.Revs))
	}
	if getDocumentOptions.RevsInfo != nil {
		builder.AddQuery("revs_info", fmt.Sprint(*getDocumentOptions.RevsInfo))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// GetDocumentAsRelated : Retrieve a document as related
// Returns document with the specified `doc_id` from the specified database. Unless you request a specific revision, the
// latest revision of the document is always returned.
func (cloudant *CloudantV1) GetDocumentAsRelated(getDocumentOptions *GetDocumentOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.GetDocumentAsRelatedWithContext(context.Background(), getDocumentOptions)
}

// GetDocumentAsRelatedWithContext is an alternate form of the GetDocumentAsRelated method which supports a Context parameter
func (cloudant *CloudantV1) GetDocumentAsRelatedWithContext(ctx context.Context, getDocumentOptions *GetDocumentOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDocumentOptions, "getDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDocumentOptions, "getDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getDocumentOptions.Db,
		"doc_id": *getDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDocumentAsRelated")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "multipart/related")
	if getDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*getDocumentOptions.IfNoneMatch))
	}

	if getDocumentOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*getDocumentOptions.Attachments))
	}
	if getDocumentOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*getDocumentOptions.AttEncodingInfo))
	}
	if getDocumentOptions.Conflicts != nil {
		builder.AddQuery("conflicts", fmt.Sprint(*getDocumentOptions.Conflicts))
	}
	if getDocumentOptions.DeletedConflicts != nil {
		builder.AddQuery("deleted_conflicts", fmt.Sprint(*getDocumentOptions.DeletedConflicts))
	}
	if getDocumentOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*getDocumentOptions.Latest))
	}
	if getDocumentOptions.LocalSeq != nil {
		builder.AddQuery("local_seq", fmt.Sprint(*getDocumentOptions.LocalSeq))
	}
	if getDocumentOptions.Meta != nil {
		builder.AddQuery("meta", fmt.Sprint(*getDocumentOptions.Meta))
	}
	if getDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*getDocumentOptions.Rev))
	}
	if getDocumentOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*getDocumentOptions.Revs))
	}
	if getDocumentOptions.RevsInfo != nil {
		builder.AddQuery("revs_info", fmt.Sprint(*getDocumentOptions.RevsInfo))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// GetDocumentAsStream : Retrieve a document as stream
// Returns document with the specified `doc_id` from the specified database. Unless you request a specific revision, the
// latest revision of the document is always returned.
func (cloudant *CloudantV1) GetDocumentAsStream(getDocumentOptions *GetDocumentOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.GetDocumentAsStreamWithContext(context.Background(), getDocumentOptions)
}

// GetDocumentAsStreamWithContext is an alternate form of the GetDocumentAsStream method which supports a Context parameter
func (cloudant *CloudantV1) GetDocumentAsStreamWithContext(ctx context.Context, getDocumentOptions *GetDocumentOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDocumentOptions, "getDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDocumentOptions, "getDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getDocumentOptions.Db,
		"doc_id": *getDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDocumentAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*getDocumentOptions.IfNoneMatch))
	}

	if getDocumentOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*getDocumentOptions.Attachments))
	}
	if getDocumentOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*getDocumentOptions.AttEncodingInfo))
	}
	if getDocumentOptions.Conflicts != nil {
		builder.AddQuery("conflicts", fmt.Sprint(*getDocumentOptions.Conflicts))
	}
	if getDocumentOptions.DeletedConflicts != nil {
		builder.AddQuery("deleted_conflicts", fmt.Sprint(*getDocumentOptions.DeletedConflicts))
	}
	if getDocumentOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*getDocumentOptions.Latest))
	}
	if getDocumentOptions.LocalSeq != nil {
		builder.AddQuery("local_seq", fmt.Sprint(*getDocumentOptions.LocalSeq))
	}
	if getDocumentOptions.Meta != nil {
		builder.AddQuery("meta", fmt.Sprint(*getDocumentOptions.Meta))
	}
	if getDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*getDocumentOptions.Rev))
	}
	if getDocumentOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*getDocumentOptions.Revs))
	}
	if getDocumentOptions.RevsInfo != nil {
		builder.AddQuery("revs_info", fmt.Sprint(*getDocumentOptions.RevsInfo))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PutDocument : Create or modify a document
// The PUT method creates a new named document, or creates a new revision of the existing document. Unlike the `POST
// /{db}` request, you must specify the document ID in the request URL.
func (cloudant *CloudantV1) PutDocument(putDocumentOptions *PutDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.PutDocumentWithContext(context.Background(), putDocumentOptions)
}

// PutDocumentWithContext is an alternate form of the PutDocument method which supports a Context parameter
func (cloudant *CloudantV1) PutDocumentWithContext(ctx context.Context, putDocumentOptions *PutDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putDocumentOptions, "putDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putDocumentOptions, "putDocumentOptions")
	if err != nil {
		return
	}

	if putDocumentOptions.Document != nil && putDocumentOptions.ContentType == nil {
		putDocumentOptions.SetContentType("application/json")
	}

	pathParamsMap := map[string]string{
		"db": *putDocumentOptions.Db,
		"doc_id": *putDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if putDocumentOptions.ContentType != nil {
		builder.AddHeader("Content-Type", fmt.Sprint(*putDocumentOptions.ContentType))
	}
	if putDocumentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*putDocumentOptions.IfMatch))
	}

	if putDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*putDocumentOptions.Batch))
	}
	if putDocumentOptions.NewEdits != nil {
		builder.AddQuery("new_edits", fmt.Sprint(*putDocumentOptions.NewEdits))
	}
	if putDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*putDocumentOptions.Rev))
	}

	_, err = builder.SetBodyContent(core.StringNilMapper(putDocumentOptions.ContentType), putDocumentOptions.Document, nil, putDocumentOptions.Body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// HeadDesignDocument : Retrieve the HTTP headers for a design document
// This method supports the same query arguments as the `GET /{db}/_design/{ddoc}` method, but the results include only
// the header information (including design document size, and the revision as an ETag). The ETag header shows the
// current revision for the requested design document, and if you requested the design document in full, the
// Content-Length specifies the length of the data. If you add any of the query arguments, then the resulting HTTP
// headers correspond to what is returned for the equivalent GET request.
func (cloudant *CloudantV1) HeadDesignDocument(headDesignDocumentOptions *HeadDesignDocumentOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadDesignDocumentWithContext(context.Background(), headDesignDocumentOptions)
}

// HeadDesignDocumentWithContext is an alternate form of the HeadDesignDocument method which supports a Context parameter
func (cloudant *CloudantV1) HeadDesignDocumentWithContext(ctx context.Context, headDesignDocumentOptions *HeadDesignDocumentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(headDesignDocumentOptions, "headDesignDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(headDesignDocumentOptions, "headDesignDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *headDesignDocumentOptions.Db,
		"ddoc": *headDesignDocumentOptions.Ddoc,
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range headDesignDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadDesignDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if headDesignDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*headDesignDocumentOptions.IfNoneMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// DeleteDesignDocument : Delete a design document
// Marks the specified design document as deleted by adding a `_deleted` field with the value `true`. Documents with
// this field are not returned with requests but stay in the database. You must supply the current (latest) revision,
// either by using the `rev` parameter or by using the `If-Match` header to specify the revision.
func (cloudant *CloudantV1) DeleteDesignDocument(deleteDesignDocumentOptions *DeleteDesignDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.DeleteDesignDocumentWithContext(context.Background(), deleteDesignDocumentOptions)
}

// DeleteDesignDocumentWithContext is an alternate form of the DeleteDesignDocument method which supports a Context parameter
func (cloudant *CloudantV1) DeleteDesignDocumentWithContext(ctx context.Context, deleteDesignDocumentOptions *DeleteDesignDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDesignDocumentOptions, "deleteDesignDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteDesignDocumentOptions, "deleteDesignDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *deleteDesignDocumentOptions.Db,
		"ddoc": *deleteDesignDocumentOptions.Ddoc,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteDesignDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "DeleteDesignDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteDesignDocumentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deleteDesignDocumentOptions.IfMatch))
	}

	if deleteDesignDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*deleteDesignDocumentOptions.Batch))
	}
	if deleteDesignDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*deleteDesignDocumentOptions.Rev))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetDesignDocument : Retrieve a design document
// Returns design document with the specified `doc_id` from the specified database. Unless you request a specific
// revision, the current revision of the design document is always returned.
func (cloudant *CloudantV1) GetDesignDocument(getDesignDocumentOptions *GetDesignDocumentOptions) (result *DesignDocument, response *core.DetailedResponse, err error) {
	return cloudant.GetDesignDocumentWithContext(context.Background(), getDesignDocumentOptions)
}

// GetDesignDocumentWithContext is an alternate form of the GetDesignDocument method which supports a Context parameter
func (cloudant *CloudantV1) GetDesignDocumentWithContext(ctx context.Context, getDesignDocumentOptions *GetDesignDocumentOptions) (result *DesignDocument, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDesignDocumentOptions, "getDesignDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDesignDocumentOptions, "getDesignDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getDesignDocumentOptions.Db,
		"ddoc": *getDesignDocumentOptions.Ddoc,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDesignDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDesignDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDesignDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*getDesignDocumentOptions.IfNoneMatch))
	}

	if getDesignDocumentOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*getDesignDocumentOptions.Attachments))
	}
	if getDesignDocumentOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*getDesignDocumentOptions.AttEncodingInfo))
	}
	if getDesignDocumentOptions.Conflicts != nil {
		builder.AddQuery("conflicts", fmt.Sprint(*getDesignDocumentOptions.Conflicts))
	}
	if getDesignDocumentOptions.DeletedConflicts != nil {
		builder.AddQuery("deleted_conflicts", fmt.Sprint(*getDesignDocumentOptions.DeletedConflicts))
	}
	if getDesignDocumentOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*getDesignDocumentOptions.Latest))
	}
	if getDesignDocumentOptions.LocalSeq != nil {
		builder.AddQuery("local_seq", fmt.Sprint(*getDesignDocumentOptions.LocalSeq))
	}
	if getDesignDocumentOptions.Meta != nil {
		builder.AddQuery("meta", fmt.Sprint(*getDesignDocumentOptions.Meta))
	}
	if getDesignDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*getDesignDocumentOptions.Rev))
	}
	if getDesignDocumentOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*getDesignDocumentOptions.Revs))
	}
	if getDesignDocumentOptions.RevsInfo != nil {
		builder.AddQuery("revs_info", fmt.Sprint(*getDesignDocumentOptions.RevsInfo))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDesignDocument)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutDesignDocument : Create or modify a design document
// The PUT method creates a new named design document, or creates a new revision of the existing design document.
func (cloudant *CloudantV1) PutDesignDocument(putDesignDocumentOptions *PutDesignDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.PutDesignDocumentWithContext(context.Background(), putDesignDocumentOptions)
}

// PutDesignDocumentWithContext is an alternate form of the PutDesignDocument method which supports a Context parameter
func (cloudant *CloudantV1) PutDesignDocumentWithContext(ctx context.Context, putDesignDocumentOptions *PutDesignDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putDesignDocumentOptions, "putDesignDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putDesignDocumentOptions, "putDesignDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *putDesignDocumentOptions.Db,
		"ddoc": *putDesignDocumentOptions.Ddoc,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putDesignDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutDesignDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if putDesignDocumentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*putDesignDocumentOptions.IfMatch))
	}

	if putDesignDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*putDesignDocumentOptions.Batch))
	}
	if putDesignDocumentOptions.NewEdits != nil {
		builder.AddQuery("new_edits", fmt.Sprint(*putDesignDocumentOptions.NewEdits))
	}
	if putDesignDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*putDesignDocumentOptions.Rev))
	}

	_, err = builder.SetBodyContentJSON(putDesignDocumentOptions.DesignDocument)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetDesignDocumentInformation : Retrieve information about a design document
// Retrieves information about the specified design document, including the index, index size, and current status of the
// design document and associated index information.
func (cloudant *CloudantV1) GetDesignDocumentInformation(getDesignDocumentInformationOptions *GetDesignDocumentInformationOptions) (result *DesignDocumentInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetDesignDocumentInformationWithContext(context.Background(), getDesignDocumentInformationOptions)
}

// GetDesignDocumentInformationWithContext is an alternate form of the GetDesignDocumentInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetDesignDocumentInformationWithContext(ctx context.Context, getDesignDocumentInformationOptions *GetDesignDocumentInformationOptions) (result *DesignDocumentInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDesignDocumentInformationOptions, "getDesignDocumentInformationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDesignDocumentInformationOptions, "getDesignDocumentInformationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getDesignDocumentInformationOptions.Db,
		"ddoc": *getDesignDocumentInformationOptions.Ddoc,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_info`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDesignDocumentInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDesignDocumentInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDesignDocumentInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostDesignDocs : Query a list of all design documents in a database
// Queries the index of all design document IDs. The results matching the request body parameters are returned in a JSON
// object, including a list of matching design documents with basic contents, such as the ID and revision. When no
// request body parameters are specified, results for all design documents in the database are returned. Optionally, the
// design document content or additional metadata can be included in the response.
func (cloudant *CloudantV1) PostDesignDocs(postDesignDocsOptions *PostDesignDocsOptions) (result *AllDocsResult, response *core.DetailedResponse, err error) {
	return cloudant.PostDesignDocsWithContext(context.Background(), postDesignDocsOptions)
}

// PostDesignDocsWithContext is an alternate form of the PostDesignDocs method which supports a Context parameter
func (cloudant *CloudantV1) PostDesignDocsWithContext(ctx context.Context, postDesignDocsOptions *PostDesignDocsOptions) (result *AllDocsResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDesignDocsOptions, "postDesignDocsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postDesignDocsOptions, "postDesignDocsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postDesignDocsOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design_docs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postDesignDocsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostDesignDocs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if postDesignDocsOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*postDesignDocsOptions.Accept))
	}

	body := make(map[string]interface{})
	if postDesignDocsOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postDesignDocsOptions.AttEncodingInfo
	}
	if postDesignDocsOptions.Attachments != nil {
		body["attachments"] = postDesignDocsOptions.Attachments
	}
	if postDesignDocsOptions.Conflicts != nil {
		body["conflicts"] = postDesignDocsOptions.Conflicts
	}
	if postDesignDocsOptions.Descending != nil {
		body["descending"] = postDesignDocsOptions.Descending
	}
	if postDesignDocsOptions.IncludeDocs != nil {
		body["include_docs"] = postDesignDocsOptions.IncludeDocs
	}
	if postDesignDocsOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postDesignDocsOptions.InclusiveEnd
	}
	if postDesignDocsOptions.Limit != nil {
		body["limit"] = postDesignDocsOptions.Limit
	}
	if postDesignDocsOptions.Skip != nil {
		body["skip"] = postDesignDocsOptions.Skip
	}
	if postDesignDocsOptions.UpdateSeq != nil {
		body["update_seq"] = postDesignDocsOptions.UpdateSeq
	}
	if postDesignDocsOptions.Endkey != nil {
		body["endkey"] = postDesignDocsOptions.Endkey
	}
	if postDesignDocsOptions.Key != nil {
		body["key"] = postDesignDocsOptions.Key
	}
	if postDesignDocsOptions.Keys != nil {
		body["keys"] = postDesignDocsOptions.Keys
	}
	if postDesignDocsOptions.Startkey != nil {
		body["startkey"] = postDesignDocsOptions.Startkey
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllDocsResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostDesignDocsQueries : Multi-query the list of all design documents
// This operation runs multiple view queries of all design documents in the database. This operation enables you to
// request numerous queries in a single request, in place of multiple POST `/{db}/_design_docs` requests.
func (cloudant *CloudantV1) PostDesignDocsQueries(postDesignDocsQueriesOptions *PostDesignDocsQueriesOptions) (result *AllDocsQueriesResult, response *core.DetailedResponse, err error) {
	return cloudant.PostDesignDocsQueriesWithContext(context.Background(), postDesignDocsQueriesOptions)
}

// PostDesignDocsQueriesWithContext is an alternate form of the PostDesignDocsQueries method which supports a Context parameter
func (cloudant *CloudantV1) PostDesignDocsQueriesWithContext(ctx context.Context, postDesignDocsQueriesOptions *PostDesignDocsQueriesOptions) (result *AllDocsQueriesResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDesignDocsQueriesOptions, "postDesignDocsQueriesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postDesignDocsQueriesOptions, "postDesignDocsQueriesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postDesignDocsQueriesOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design_docs/queries`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postDesignDocsQueriesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostDesignDocsQueries")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if postDesignDocsQueriesOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*postDesignDocsQueriesOptions.Accept))
	}

	body := make(map[string]interface{})
	if postDesignDocsQueriesOptions.Queries != nil {
		body["queries"] = postDesignDocsQueriesOptions.Queries
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllDocsQueriesResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostView : Query a MapReduce view
// This operation queries the specified MapReduce view of the specified design document. By default, the map and reduce
// functions of the view are run to update the view before returning the response. The advantage of using the HTTP
// `POST` method is that the query is submitted as a JSON object in the request body. This avoids the limitations of
// passing query options as URL query parameters of a `GET` request.
func (cloudant *CloudantV1) PostView(postViewOptions *PostViewOptions) (result *ViewResult, response *core.DetailedResponse, err error) {
	return cloudant.PostViewWithContext(context.Background(), postViewOptions)
}

// PostViewWithContext is an alternate form of the PostView method which supports a Context parameter
func (cloudant *CloudantV1) PostViewWithContext(ctx context.Context, postViewOptions *PostViewOptions) (result *ViewResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postViewOptions, "postViewOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postViewOptions, "postViewOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postViewOptions.Db,
		"ddoc": *postViewOptions.Ddoc,
		"view": *postViewOptions.View,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_view/{view}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postViewOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostView")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postViewOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postViewOptions.AttEncodingInfo
	}
	if postViewOptions.Attachments != nil {
		body["attachments"] = postViewOptions.Attachments
	}
	if postViewOptions.Conflicts != nil {
		body["conflicts"] = postViewOptions.Conflicts
	}
	if postViewOptions.Descending != nil {
		body["descending"] = postViewOptions.Descending
	}
	if postViewOptions.IncludeDocs != nil {
		body["include_docs"] = postViewOptions.IncludeDocs
	}
	if postViewOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postViewOptions.InclusiveEnd
	}
	if postViewOptions.Limit != nil {
		body["limit"] = postViewOptions.Limit
	}
	if postViewOptions.Skip != nil {
		body["skip"] = postViewOptions.Skip
	}
	if postViewOptions.UpdateSeq != nil {
		body["update_seq"] = postViewOptions.UpdateSeq
	}
	if postViewOptions.Endkey != nil {
		body["endkey"] = postViewOptions.Endkey
	}
	if postViewOptions.EndkeyDocid != nil {
		body["endkey_docid"] = postViewOptions.EndkeyDocid
	}
	if postViewOptions.Group != nil {
		body["group"] = postViewOptions.Group
	}
	if postViewOptions.GroupLevel != nil {
		body["group_level"] = postViewOptions.GroupLevel
	}
	if postViewOptions.Key != nil {
		body["key"] = postViewOptions.Key
	}
	if postViewOptions.Keys != nil {
		body["keys"] = postViewOptions.Keys
	}
	if postViewOptions.Reduce != nil {
		body["reduce"] = postViewOptions.Reduce
	}
	if postViewOptions.Stable != nil {
		body["stable"] = postViewOptions.Stable
	}
	if postViewOptions.Startkey != nil {
		body["startkey"] = postViewOptions.Startkey
	}
	if postViewOptions.StartkeyDocid != nil {
		body["startkey_docid"] = postViewOptions.StartkeyDocid
	}
	if postViewOptions.Update != nil {
		body["update"] = postViewOptions.Update
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalViewResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostViewAsStream : Query a MapReduce view as stream
// This operation queries the specified MapReduce view of the specified design document. By default, the map and reduce
// functions of the view are run to update the view before returning the response. The advantage of using the HTTP
// `POST` method is that the query is submitted as a JSON object in the request body. This avoids the limitations of
// passing query options as URL query parameters of a `GET` request.
func (cloudant *CloudantV1) PostViewAsStream(postViewOptions *PostViewOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostViewAsStreamWithContext(context.Background(), postViewOptions)
}

// PostViewAsStreamWithContext is an alternate form of the PostViewAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostViewAsStreamWithContext(ctx context.Context, postViewOptions *PostViewOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postViewOptions, "postViewOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postViewOptions, "postViewOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postViewOptions.Db,
		"ddoc": *postViewOptions.Ddoc,
		"view": *postViewOptions.View,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_view/{view}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postViewOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostViewAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postViewOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postViewOptions.AttEncodingInfo
	}
	if postViewOptions.Attachments != nil {
		body["attachments"] = postViewOptions.Attachments
	}
	if postViewOptions.Conflicts != nil {
		body["conflicts"] = postViewOptions.Conflicts
	}
	if postViewOptions.Descending != nil {
		body["descending"] = postViewOptions.Descending
	}
	if postViewOptions.IncludeDocs != nil {
		body["include_docs"] = postViewOptions.IncludeDocs
	}
	if postViewOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postViewOptions.InclusiveEnd
	}
	if postViewOptions.Limit != nil {
		body["limit"] = postViewOptions.Limit
	}
	if postViewOptions.Skip != nil {
		body["skip"] = postViewOptions.Skip
	}
	if postViewOptions.UpdateSeq != nil {
		body["update_seq"] = postViewOptions.UpdateSeq
	}
	if postViewOptions.Endkey != nil {
		body["endkey"] = postViewOptions.Endkey
	}
	if postViewOptions.EndkeyDocid != nil {
		body["endkey_docid"] = postViewOptions.EndkeyDocid
	}
	if postViewOptions.Group != nil {
		body["group"] = postViewOptions.Group
	}
	if postViewOptions.GroupLevel != nil {
		body["group_level"] = postViewOptions.GroupLevel
	}
	if postViewOptions.Key != nil {
		body["key"] = postViewOptions.Key
	}
	if postViewOptions.Keys != nil {
		body["keys"] = postViewOptions.Keys
	}
	if postViewOptions.Reduce != nil {
		body["reduce"] = postViewOptions.Reduce
	}
	if postViewOptions.Stable != nil {
		body["stable"] = postViewOptions.Stable
	}
	if postViewOptions.Startkey != nil {
		body["startkey"] = postViewOptions.Startkey
	}
	if postViewOptions.StartkeyDocid != nil {
		body["startkey_docid"] = postViewOptions.StartkeyDocid
	}
	if postViewOptions.Update != nil {
		body["update"] = postViewOptions.Update
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostViewQueries : Multi-query a MapReduce view
// This operation runs multiple specified view queries against the view function from the specified design document.
func (cloudant *CloudantV1) PostViewQueries(postViewQueriesOptions *PostViewQueriesOptions) (result *ViewQueriesResult, response *core.DetailedResponse, err error) {
	return cloudant.PostViewQueriesWithContext(context.Background(), postViewQueriesOptions)
}

// PostViewQueriesWithContext is an alternate form of the PostViewQueries method which supports a Context parameter
func (cloudant *CloudantV1) PostViewQueriesWithContext(ctx context.Context, postViewQueriesOptions *PostViewQueriesOptions) (result *ViewQueriesResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postViewQueriesOptions, "postViewQueriesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postViewQueriesOptions, "postViewQueriesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postViewQueriesOptions.Db,
		"ddoc": *postViewQueriesOptions.Ddoc,
		"view": *postViewQueriesOptions.View,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_view/{view}/queries`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postViewQueriesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostViewQueries")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postViewQueriesOptions.Queries != nil {
		body["queries"] = postViewQueriesOptions.Queries
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalViewQueriesResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostViewQueriesAsStream : Multi-query a MapReduce view as stream
// This operation runs multiple specified view queries against the view function from the specified design document.
func (cloudant *CloudantV1) PostViewQueriesAsStream(postViewQueriesOptions *PostViewQueriesOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostViewQueriesAsStreamWithContext(context.Background(), postViewQueriesOptions)
}

// PostViewQueriesAsStreamWithContext is an alternate form of the PostViewQueriesAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostViewQueriesAsStreamWithContext(ctx context.Context, postViewQueriesOptions *PostViewQueriesOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postViewQueriesOptions, "postViewQueriesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postViewQueriesOptions, "postViewQueriesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postViewQueriesOptions.Db,
		"ddoc": *postViewQueriesOptions.Ddoc,
		"view": *postViewQueriesOptions.View,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_view/{view}/queries`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postViewQueriesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostViewQueriesAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postViewQueriesOptions.Queries != nil {
		body["queries"] = postViewQueriesOptions.Queries
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// GetPartitionInformation : Retrieve information about a database partition
// Given a partition key, return the database name, sizes, partition, doc count, and doc delete count.
func (cloudant *CloudantV1) GetPartitionInformation(getPartitionInformationOptions *GetPartitionInformationOptions) (result *PartitionInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetPartitionInformationWithContext(context.Background(), getPartitionInformationOptions)
}

// GetPartitionInformationWithContext is an alternate form of the GetPartitionInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetPartitionInformationWithContext(ctx context.Context, getPartitionInformationOptions *GetPartitionInformationOptions) (result *PartitionInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPartitionInformationOptions, "getPartitionInformationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getPartitionInformationOptions, "getPartitionInformationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getPartitionInformationOptions.Db,
		"partition_key": *getPartitionInformationOptions.PartitionKey,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getPartitionInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetPartitionInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPartitionInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostPartitionAllDocs : Query a list of all documents in a database partition
// Queries the primary index (all document IDs). The results that match the query parameters are returned in a JSON
// object, including a list of matching documents with basic contents, such as the ID and revision. When no query
// parameters are specified, results for all documents in the database partition are returned. Optionally, document
// content or additional metadata can be included in the response.
func (cloudant *CloudantV1) PostPartitionAllDocs(postPartitionAllDocsOptions *PostPartitionAllDocsOptions) (result *AllDocsResult, response *core.DetailedResponse, err error) {
	return cloudant.PostPartitionAllDocsWithContext(context.Background(), postPartitionAllDocsOptions)
}

// PostPartitionAllDocsWithContext is an alternate form of the PostPartitionAllDocs method which supports a Context parameter
func (cloudant *CloudantV1) PostPartitionAllDocsWithContext(ctx context.Context, postPartitionAllDocsOptions *PostPartitionAllDocsOptions) (result *AllDocsResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postPartitionAllDocsOptions, "postPartitionAllDocsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postPartitionAllDocsOptions, "postPartitionAllDocsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postPartitionAllDocsOptions.Db,
		"partition_key": *postPartitionAllDocsOptions.PartitionKey,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}/_all_docs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postPartitionAllDocsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostPartitionAllDocs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postPartitionAllDocsOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postPartitionAllDocsOptions.AttEncodingInfo
	}
	if postPartitionAllDocsOptions.Attachments != nil {
		body["attachments"] = postPartitionAllDocsOptions.Attachments
	}
	if postPartitionAllDocsOptions.Conflicts != nil {
		body["conflicts"] = postPartitionAllDocsOptions.Conflicts
	}
	if postPartitionAllDocsOptions.Descending != nil {
		body["descending"] = postPartitionAllDocsOptions.Descending
	}
	if postPartitionAllDocsOptions.IncludeDocs != nil {
		body["include_docs"] = postPartitionAllDocsOptions.IncludeDocs
	}
	if postPartitionAllDocsOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postPartitionAllDocsOptions.InclusiveEnd
	}
	if postPartitionAllDocsOptions.Limit != nil {
		body["limit"] = postPartitionAllDocsOptions.Limit
	}
	if postPartitionAllDocsOptions.Skip != nil {
		body["skip"] = postPartitionAllDocsOptions.Skip
	}
	if postPartitionAllDocsOptions.UpdateSeq != nil {
		body["update_seq"] = postPartitionAllDocsOptions.UpdateSeq
	}
	if postPartitionAllDocsOptions.Endkey != nil {
		body["endkey"] = postPartitionAllDocsOptions.Endkey
	}
	if postPartitionAllDocsOptions.Key != nil {
		body["key"] = postPartitionAllDocsOptions.Key
	}
	if postPartitionAllDocsOptions.Keys != nil {
		body["keys"] = postPartitionAllDocsOptions.Keys
	}
	if postPartitionAllDocsOptions.Startkey != nil {
		body["startkey"] = postPartitionAllDocsOptions.Startkey
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllDocsResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostPartitionAllDocsAsStream : Query a list of all documents in a database partition as stream
// Queries the primary index (all document IDs). The results that match the query parameters are returned in a JSON
// object, including a list of matching documents with basic contents, such as the ID and revision. When no query
// parameters are specified, results for all documents in the database partition are returned. Optionally, document
// content or additional metadata can be included in the response.
func (cloudant *CloudantV1) PostPartitionAllDocsAsStream(postPartitionAllDocsOptions *PostPartitionAllDocsOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostPartitionAllDocsAsStreamWithContext(context.Background(), postPartitionAllDocsOptions)
}

// PostPartitionAllDocsAsStreamWithContext is an alternate form of the PostPartitionAllDocsAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostPartitionAllDocsAsStreamWithContext(ctx context.Context, postPartitionAllDocsOptions *PostPartitionAllDocsOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postPartitionAllDocsOptions, "postPartitionAllDocsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postPartitionAllDocsOptions, "postPartitionAllDocsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postPartitionAllDocsOptions.Db,
		"partition_key": *postPartitionAllDocsOptions.PartitionKey,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}/_all_docs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postPartitionAllDocsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostPartitionAllDocsAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postPartitionAllDocsOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postPartitionAllDocsOptions.AttEncodingInfo
	}
	if postPartitionAllDocsOptions.Attachments != nil {
		body["attachments"] = postPartitionAllDocsOptions.Attachments
	}
	if postPartitionAllDocsOptions.Conflicts != nil {
		body["conflicts"] = postPartitionAllDocsOptions.Conflicts
	}
	if postPartitionAllDocsOptions.Descending != nil {
		body["descending"] = postPartitionAllDocsOptions.Descending
	}
	if postPartitionAllDocsOptions.IncludeDocs != nil {
		body["include_docs"] = postPartitionAllDocsOptions.IncludeDocs
	}
	if postPartitionAllDocsOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postPartitionAllDocsOptions.InclusiveEnd
	}
	if postPartitionAllDocsOptions.Limit != nil {
		body["limit"] = postPartitionAllDocsOptions.Limit
	}
	if postPartitionAllDocsOptions.Skip != nil {
		body["skip"] = postPartitionAllDocsOptions.Skip
	}
	if postPartitionAllDocsOptions.UpdateSeq != nil {
		body["update_seq"] = postPartitionAllDocsOptions.UpdateSeq
	}
	if postPartitionAllDocsOptions.Endkey != nil {
		body["endkey"] = postPartitionAllDocsOptions.Endkey
	}
	if postPartitionAllDocsOptions.Key != nil {
		body["key"] = postPartitionAllDocsOptions.Key
	}
	if postPartitionAllDocsOptions.Keys != nil {
		body["keys"] = postPartitionAllDocsOptions.Keys
	}
	if postPartitionAllDocsOptions.Startkey != nil {
		body["startkey"] = postPartitionAllDocsOptions.Startkey
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostPartitionSearch : Query a database partition search index
// Partitioned Search indexes, which are defined in design documents, allow partition databases to be queried by using
// Lucene Query Parser Syntax. Search indexes are defined by an index function, similar to a map function in MapReduce
// views. The index function decides what data to index and store in the index.
func (cloudant *CloudantV1) PostPartitionSearch(postPartitionSearchOptions *PostPartitionSearchOptions) (result *SearchResult, response *core.DetailedResponse, err error) {
	return cloudant.PostPartitionSearchWithContext(context.Background(), postPartitionSearchOptions)
}

// PostPartitionSearchWithContext is an alternate form of the PostPartitionSearch method which supports a Context parameter
func (cloudant *CloudantV1) PostPartitionSearchWithContext(ctx context.Context, postPartitionSearchOptions *PostPartitionSearchOptions) (result *SearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postPartitionSearchOptions, "postPartitionSearchOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postPartitionSearchOptions, "postPartitionSearchOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postPartitionSearchOptions.Db,
		"partition_key": *postPartitionSearchOptions.PartitionKey,
		"ddoc": *postPartitionSearchOptions.Ddoc,
		"index": *postPartitionSearchOptions.Index,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}/_design/{ddoc}/_search/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postPartitionSearchOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostPartitionSearch")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postPartitionSearchOptions.Query != nil {
		body["query"] = postPartitionSearchOptions.Query
	}
	if postPartitionSearchOptions.Bookmark != nil {
		body["bookmark"] = postPartitionSearchOptions.Bookmark
	}
	if postPartitionSearchOptions.HighlightFields != nil {
		body["highlight_fields"] = postPartitionSearchOptions.HighlightFields
	}
	if postPartitionSearchOptions.HighlightNumber != nil {
		body["highlight_number"] = postPartitionSearchOptions.HighlightNumber
	}
	if postPartitionSearchOptions.HighlightPostTag != nil {
		body["highlight_post_tag"] = postPartitionSearchOptions.HighlightPostTag
	}
	if postPartitionSearchOptions.HighlightPreTag != nil {
		body["highlight_pre_tag"] = postPartitionSearchOptions.HighlightPreTag
	}
	if postPartitionSearchOptions.HighlightSize != nil {
		body["highlight_size"] = postPartitionSearchOptions.HighlightSize
	}
	if postPartitionSearchOptions.IncludeDocs != nil {
		body["include_docs"] = postPartitionSearchOptions.IncludeDocs
	}
	if postPartitionSearchOptions.IncludeFields != nil {
		body["include_fields"] = postPartitionSearchOptions.IncludeFields
	}
	if postPartitionSearchOptions.Limit != nil {
		body["limit"] = postPartitionSearchOptions.Limit
	}
	if postPartitionSearchOptions.Sort != nil {
		body["sort"] = postPartitionSearchOptions.Sort
	}
	if postPartitionSearchOptions.Stale != nil {
		body["stale"] = postPartitionSearchOptions.Stale
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSearchResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostPartitionSearchAsStream : Query a database partition search index as stream
// Partitioned Search indexes, which are defined in design documents, allow partition databases to be queried by using
// Lucene Query Parser Syntax. Search indexes are defined by an index function, similar to a map function in MapReduce
// views. The index function decides what data to index and store in the index.
func (cloudant *CloudantV1) PostPartitionSearchAsStream(postPartitionSearchOptions *PostPartitionSearchOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostPartitionSearchAsStreamWithContext(context.Background(), postPartitionSearchOptions)
}

// PostPartitionSearchAsStreamWithContext is an alternate form of the PostPartitionSearchAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostPartitionSearchAsStreamWithContext(ctx context.Context, postPartitionSearchOptions *PostPartitionSearchOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postPartitionSearchOptions, "postPartitionSearchOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postPartitionSearchOptions, "postPartitionSearchOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postPartitionSearchOptions.Db,
		"partition_key": *postPartitionSearchOptions.PartitionKey,
		"ddoc": *postPartitionSearchOptions.Ddoc,
		"index": *postPartitionSearchOptions.Index,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}/_design/{ddoc}/_search/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postPartitionSearchOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostPartitionSearchAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postPartitionSearchOptions.Query != nil {
		body["query"] = postPartitionSearchOptions.Query
	}
	if postPartitionSearchOptions.Bookmark != nil {
		body["bookmark"] = postPartitionSearchOptions.Bookmark
	}
	if postPartitionSearchOptions.HighlightFields != nil {
		body["highlight_fields"] = postPartitionSearchOptions.HighlightFields
	}
	if postPartitionSearchOptions.HighlightNumber != nil {
		body["highlight_number"] = postPartitionSearchOptions.HighlightNumber
	}
	if postPartitionSearchOptions.HighlightPostTag != nil {
		body["highlight_post_tag"] = postPartitionSearchOptions.HighlightPostTag
	}
	if postPartitionSearchOptions.HighlightPreTag != nil {
		body["highlight_pre_tag"] = postPartitionSearchOptions.HighlightPreTag
	}
	if postPartitionSearchOptions.HighlightSize != nil {
		body["highlight_size"] = postPartitionSearchOptions.HighlightSize
	}
	if postPartitionSearchOptions.IncludeDocs != nil {
		body["include_docs"] = postPartitionSearchOptions.IncludeDocs
	}
	if postPartitionSearchOptions.IncludeFields != nil {
		body["include_fields"] = postPartitionSearchOptions.IncludeFields
	}
	if postPartitionSearchOptions.Limit != nil {
		body["limit"] = postPartitionSearchOptions.Limit
	}
	if postPartitionSearchOptions.Sort != nil {
		body["sort"] = postPartitionSearchOptions.Sort
	}
	if postPartitionSearchOptions.Stale != nil {
		body["stale"] = postPartitionSearchOptions.Stale
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostPartitionView : Query a database partition MapReduce view function
// Runs the specified view function from the specified design document. Unlike `GET /{db}/_design/{ddoc}/_view/{view}`
// for accessing views, the POST method supports the specification of explicit keys to be retrieved from the view
// results. The remainder of the POST view functionality is identical to the `GET /{db}/_design/{ddoc}/_view/{view}`
// API.
func (cloudant *CloudantV1) PostPartitionView(postPartitionViewOptions *PostPartitionViewOptions) (result *ViewResult, response *core.DetailedResponse, err error) {
	return cloudant.PostPartitionViewWithContext(context.Background(), postPartitionViewOptions)
}

// PostPartitionViewWithContext is an alternate form of the PostPartitionView method which supports a Context parameter
func (cloudant *CloudantV1) PostPartitionViewWithContext(ctx context.Context, postPartitionViewOptions *PostPartitionViewOptions) (result *ViewResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postPartitionViewOptions, "postPartitionViewOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postPartitionViewOptions, "postPartitionViewOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postPartitionViewOptions.Db,
		"partition_key": *postPartitionViewOptions.PartitionKey,
		"ddoc": *postPartitionViewOptions.Ddoc,
		"view": *postPartitionViewOptions.View,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}/_design/{ddoc}/_view/{view}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postPartitionViewOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostPartitionView")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postPartitionViewOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postPartitionViewOptions.AttEncodingInfo
	}
	if postPartitionViewOptions.Attachments != nil {
		body["attachments"] = postPartitionViewOptions.Attachments
	}
	if postPartitionViewOptions.Conflicts != nil {
		body["conflicts"] = postPartitionViewOptions.Conflicts
	}
	if postPartitionViewOptions.Descending != nil {
		body["descending"] = postPartitionViewOptions.Descending
	}
	if postPartitionViewOptions.IncludeDocs != nil {
		body["include_docs"] = postPartitionViewOptions.IncludeDocs
	}
	if postPartitionViewOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postPartitionViewOptions.InclusiveEnd
	}
	if postPartitionViewOptions.Limit != nil {
		body["limit"] = postPartitionViewOptions.Limit
	}
	if postPartitionViewOptions.Skip != nil {
		body["skip"] = postPartitionViewOptions.Skip
	}
	if postPartitionViewOptions.UpdateSeq != nil {
		body["update_seq"] = postPartitionViewOptions.UpdateSeq
	}
	if postPartitionViewOptions.Endkey != nil {
		body["endkey"] = postPartitionViewOptions.Endkey
	}
	if postPartitionViewOptions.EndkeyDocid != nil {
		body["endkey_docid"] = postPartitionViewOptions.EndkeyDocid
	}
	if postPartitionViewOptions.Group != nil {
		body["group"] = postPartitionViewOptions.Group
	}
	if postPartitionViewOptions.GroupLevel != nil {
		body["group_level"] = postPartitionViewOptions.GroupLevel
	}
	if postPartitionViewOptions.Key != nil {
		body["key"] = postPartitionViewOptions.Key
	}
	if postPartitionViewOptions.Keys != nil {
		body["keys"] = postPartitionViewOptions.Keys
	}
	if postPartitionViewOptions.Reduce != nil {
		body["reduce"] = postPartitionViewOptions.Reduce
	}
	if postPartitionViewOptions.Stable != nil {
		body["stable"] = postPartitionViewOptions.Stable
	}
	if postPartitionViewOptions.Startkey != nil {
		body["startkey"] = postPartitionViewOptions.Startkey
	}
	if postPartitionViewOptions.StartkeyDocid != nil {
		body["startkey_docid"] = postPartitionViewOptions.StartkeyDocid
	}
	if postPartitionViewOptions.Update != nil {
		body["update"] = postPartitionViewOptions.Update
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalViewResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostPartitionViewAsStream : Query a database partition MapReduce view function as stream
// Runs the specified view function from the specified design document. Unlike `GET /{db}/_design/{ddoc}/_view/{view}`
// for accessing views, the POST method supports the specification of explicit keys to be retrieved from the view
// results. The remainder of the POST view functionality is identical to the `GET /{db}/_design/{ddoc}/_view/{view}`
// API.
func (cloudant *CloudantV1) PostPartitionViewAsStream(postPartitionViewOptions *PostPartitionViewOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostPartitionViewAsStreamWithContext(context.Background(), postPartitionViewOptions)
}

// PostPartitionViewAsStreamWithContext is an alternate form of the PostPartitionViewAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostPartitionViewAsStreamWithContext(ctx context.Context, postPartitionViewOptions *PostPartitionViewOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postPartitionViewOptions, "postPartitionViewOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postPartitionViewOptions, "postPartitionViewOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postPartitionViewOptions.Db,
		"partition_key": *postPartitionViewOptions.PartitionKey,
		"ddoc": *postPartitionViewOptions.Ddoc,
		"view": *postPartitionViewOptions.View,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}/_design/{ddoc}/_view/{view}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postPartitionViewOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostPartitionViewAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postPartitionViewOptions.AttEncodingInfo != nil {
		body["att_encoding_info"] = postPartitionViewOptions.AttEncodingInfo
	}
	if postPartitionViewOptions.Attachments != nil {
		body["attachments"] = postPartitionViewOptions.Attachments
	}
	if postPartitionViewOptions.Conflicts != nil {
		body["conflicts"] = postPartitionViewOptions.Conflicts
	}
	if postPartitionViewOptions.Descending != nil {
		body["descending"] = postPartitionViewOptions.Descending
	}
	if postPartitionViewOptions.IncludeDocs != nil {
		body["include_docs"] = postPartitionViewOptions.IncludeDocs
	}
	if postPartitionViewOptions.InclusiveEnd != nil {
		body["inclusive_end"] = postPartitionViewOptions.InclusiveEnd
	}
	if postPartitionViewOptions.Limit != nil {
		body["limit"] = postPartitionViewOptions.Limit
	}
	if postPartitionViewOptions.Skip != nil {
		body["skip"] = postPartitionViewOptions.Skip
	}
	if postPartitionViewOptions.UpdateSeq != nil {
		body["update_seq"] = postPartitionViewOptions.UpdateSeq
	}
	if postPartitionViewOptions.Endkey != nil {
		body["endkey"] = postPartitionViewOptions.Endkey
	}
	if postPartitionViewOptions.EndkeyDocid != nil {
		body["endkey_docid"] = postPartitionViewOptions.EndkeyDocid
	}
	if postPartitionViewOptions.Group != nil {
		body["group"] = postPartitionViewOptions.Group
	}
	if postPartitionViewOptions.GroupLevel != nil {
		body["group_level"] = postPartitionViewOptions.GroupLevel
	}
	if postPartitionViewOptions.Key != nil {
		body["key"] = postPartitionViewOptions.Key
	}
	if postPartitionViewOptions.Keys != nil {
		body["keys"] = postPartitionViewOptions.Keys
	}
	if postPartitionViewOptions.Reduce != nil {
		body["reduce"] = postPartitionViewOptions.Reduce
	}
	if postPartitionViewOptions.Stable != nil {
		body["stable"] = postPartitionViewOptions.Stable
	}
	if postPartitionViewOptions.Startkey != nil {
		body["startkey"] = postPartitionViewOptions.Startkey
	}
	if postPartitionViewOptions.StartkeyDocid != nil {
		body["startkey_docid"] = postPartitionViewOptions.StartkeyDocid
	}
	if postPartitionViewOptions.Update != nil {
		body["update"] = postPartitionViewOptions.Update
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostPartitionFind : Query a database partition index by using selector syntax (POST)
// Query documents by using a declarative JSON querying syntax. Queries can use the built-in `_all_docs` index or custom
// indices, specified by using the `_index` endpoint.
func (cloudant *CloudantV1) PostPartitionFind(postPartitionFindOptions *PostPartitionFindOptions) (result *FindResult, response *core.DetailedResponse, err error) {
	return cloudant.PostPartitionFindWithContext(context.Background(), postPartitionFindOptions)
}

// PostPartitionFindWithContext is an alternate form of the PostPartitionFind method which supports a Context parameter
func (cloudant *CloudantV1) PostPartitionFindWithContext(ctx context.Context, postPartitionFindOptions *PostPartitionFindOptions) (result *FindResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postPartitionFindOptions, "postPartitionFindOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postPartitionFindOptions, "postPartitionFindOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postPartitionFindOptions.Db,
		"partition_key": *postPartitionFindOptions.PartitionKey,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}/_find`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postPartitionFindOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostPartitionFind")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postPartitionFindOptions.Selector != nil {
		body["selector"] = postPartitionFindOptions.Selector
	}
	if postPartitionFindOptions.Bookmark != nil {
		body["bookmark"] = postPartitionFindOptions.Bookmark
	}
	if postPartitionFindOptions.Conflicts != nil {
		body["conflicts"] = postPartitionFindOptions.Conflicts
	}
	if postPartitionFindOptions.ExecutionStats != nil {
		body["execution_stats"] = postPartitionFindOptions.ExecutionStats
	}
	if postPartitionFindOptions.Fields != nil {
		body["fields"] = postPartitionFindOptions.Fields
	}
	if postPartitionFindOptions.Limit != nil {
		body["limit"] = postPartitionFindOptions.Limit
	}
	if postPartitionFindOptions.Skip != nil {
		body["skip"] = postPartitionFindOptions.Skip
	}
	if postPartitionFindOptions.Sort != nil {
		body["sort"] = postPartitionFindOptions.Sort
	}
	if postPartitionFindOptions.Stable != nil {
		body["stable"] = postPartitionFindOptions.Stable
	}
	if postPartitionFindOptions.Update != nil {
		body["update"] = postPartitionFindOptions.Update
	}
	if postPartitionFindOptions.UseIndex != nil {
		body["use_index"] = postPartitionFindOptions.UseIndex
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFindResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostPartitionFindAsStream : Query a database partition index by using selector syntax (POST) as stream
// Query documents by using a declarative JSON querying syntax. Queries can use the built-in `_all_docs` index or custom
// indices, specified by using the `_index` endpoint.
func (cloudant *CloudantV1) PostPartitionFindAsStream(postPartitionFindOptions *PostPartitionFindOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostPartitionFindAsStreamWithContext(context.Background(), postPartitionFindOptions)
}

// PostPartitionFindAsStreamWithContext is an alternate form of the PostPartitionFindAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostPartitionFindAsStreamWithContext(ctx context.Context, postPartitionFindOptions *PostPartitionFindOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postPartitionFindOptions, "postPartitionFindOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postPartitionFindOptions, "postPartitionFindOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postPartitionFindOptions.Db,
		"partition_key": *postPartitionFindOptions.PartitionKey,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_partition/{partition_key}/_find`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postPartitionFindOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostPartitionFindAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postPartitionFindOptions.Selector != nil {
		body["selector"] = postPartitionFindOptions.Selector
	}
	if postPartitionFindOptions.Bookmark != nil {
		body["bookmark"] = postPartitionFindOptions.Bookmark
	}
	if postPartitionFindOptions.Conflicts != nil {
		body["conflicts"] = postPartitionFindOptions.Conflicts
	}
	if postPartitionFindOptions.ExecutionStats != nil {
		body["execution_stats"] = postPartitionFindOptions.ExecutionStats
	}
	if postPartitionFindOptions.Fields != nil {
		body["fields"] = postPartitionFindOptions.Fields
	}
	if postPartitionFindOptions.Limit != nil {
		body["limit"] = postPartitionFindOptions.Limit
	}
	if postPartitionFindOptions.Skip != nil {
		body["skip"] = postPartitionFindOptions.Skip
	}
	if postPartitionFindOptions.Sort != nil {
		body["sort"] = postPartitionFindOptions.Sort
	}
	if postPartitionFindOptions.Stable != nil {
		body["stable"] = postPartitionFindOptions.Stable
	}
	if postPartitionFindOptions.Update != nil {
		body["update"] = postPartitionFindOptions.Update
	}
	if postPartitionFindOptions.UseIndex != nil {
		body["use_index"] = postPartitionFindOptions.UseIndex
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostExplain : Retrieve information about which index is used for a query
// Shows which index is being used by the query. Parameters are the same as the [`_find`
// endpoint](#query-an-index-by-using-selector-syntax).
func (cloudant *CloudantV1) PostExplain(postExplainOptions *PostExplainOptions) (result *ExplainResult, response *core.DetailedResponse, err error) {
	return cloudant.PostExplainWithContext(context.Background(), postExplainOptions)
}

// PostExplainWithContext is an alternate form of the PostExplain method which supports a Context parameter
func (cloudant *CloudantV1) PostExplainWithContext(ctx context.Context, postExplainOptions *PostExplainOptions) (result *ExplainResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postExplainOptions, "postExplainOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postExplainOptions, "postExplainOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postExplainOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_explain`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postExplainOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostExplain")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postExplainOptions.Selector != nil {
		body["selector"] = postExplainOptions.Selector
	}
	if postExplainOptions.Bookmark != nil {
		body["bookmark"] = postExplainOptions.Bookmark
	}
	if postExplainOptions.Conflicts != nil {
		body["conflicts"] = postExplainOptions.Conflicts
	}
	if postExplainOptions.ExecutionStats != nil {
		body["execution_stats"] = postExplainOptions.ExecutionStats
	}
	if postExplainOptions.Fields != nil {
		body["fields"] = postExplainOptions.Fields
	}
	if postExplainOptions.Limit != nil {
		body["limit"] = postExplainOptions.Limit
	}
	if postExplainOptions.Skip != nil {
		body["skip"] = postExplainOptions.Skip
	}
	if postExplainOptions.Sort != nil {
		body["sort"] = postExplainOptions.Sort
	}
	if postExplainOptions.Stable != nil {
		body["stable"] = postExplainOptions.Stable
	}
	if postExplainOptions.Update != nil {
		body["update"] = postExplainOptions.Update
	}
	if postExplainOptions.UseIndex != nil {
		body["use_index"] = postExplainOptions.UseIndex
	}
	if postExplainOptions.R != nil {
		body["r"] = postExplainOptions.R
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExplainResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostFind : Query an index by using selector syntax
// Query documents by using a declarative JSON querying syntax. Queries can use the built-in `_all_docs` index or custom
// indices, specified by using the `_index` endpoint.
func (cloudant *CloudantV1) PostFind(postFindOptions *PostFindOptions) (result *FindResult, response *core.DetailedResponse, err error) {
	return cloudant.PostFindWithContext(context.Background(), postFindOptions)
}

// PostFindWithContext is an alternate form of the PostFind method which supports a Context parameter
func (cloudant *CloudantV1) PostFindWithContext(ctx context.Context, postFindOptions *PostFindOptions) (result *FindResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postFindOptions, "postFindOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postFindOptions, "postFindOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postFindOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_find`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postFindOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostFind")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postFindOptions.Selector != nil {
		body["selector"] = postFindOptions.Selector
	}
	if postFindOptions.Bookmark != nil {
		body["bookmark"] = postFindOptions.Bookmark
	}
	if postFindOptions.Conflicts != nil {
		body["conflicts"] = postFindOptions.Conflicts
	}
	if postFindOptions.ExecutionStats != nil {
		body["execution_stats"] = postFindOptions.ExecutionStats
	}
	if postFindOptions.Fields != nil {
		body["fields"] = postFindOptions.Fields
	}
	if postFindOptions.Limit != nil {
		body["limit"] = postFindOptions.Limit
	}
	if postFindOptions.Skip != nil {
		body["skip"] = postFindOptions.Skip
	}
	if postFindOptions.Sort != nil {
		body["sort"] = postFindOptions.Sort
	}
	if postFindOptions.Stable != nil {
		body["stable"] = postFindOptions.Stable
	}
	if postFindOptions.Update != nil {
		body["update"] = postFindOptions.Update
	}
	if postFindOptions.UseIndex != nil {
		body["use_index"] = postFindOptions.UseIndex
	}
	if postFindOptions.R != nil {
		body["r"] = postFindOptions.R
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFindResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostFindAsStream : Query an index by using selector syntax as stream
// Query documents by using a declarative JSON querying syntax. Queries can use the built-in `_all_docs` index or custom
// indices, specified by using the `_index` endpoint.
func (cloudant *CloudantV1) PostFindAsStream(postFindOptions *PostFindOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostFindAsStreamWithContext(context.Background(), postFindOptions)
}

// PostFindAsStreamWithContext is an alternate form of the PostFindAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostFindAsStreamWithContext(ctx context.Context, postFindOptions *PostFindOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postFindOptions, "postFindOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postFindOptions, "postFindOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postFindOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_find`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postFindOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostFindAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postFindOptions.Selector != nil {
		body["selector"] = postFindOptions.Selector
	}
	if postFindOptions.Bookmark != nil {
		body["bookmark"] = postFindOptions.Bookmark
	}
	if postFindOptions.Conflicts != nil {
		body["conflicts"] = postFindOptions.Conflicts
	}
	if postFindOptions.ExecutionStats != nil {
		body["execution_stats"] = postFindOptions.ExecutionStats
	}
	if postFindOptions.Fields != nil {
		body["fields"] = postFindOptions.Fields
	}
	if postFindOptions.Limit != nil {
		body["limit"] = postFindOptions.Limit
	}
	if postFindOptions.Skip != nil {
		body["skip"] = postFindOptions.Skip
	}
	if postFindOptions.Sort != nil {
		body["sort"] = postFindOptions.Sort
	}
	if postFindOptions.Stable != nil {
		body["stable"] = postFindOptions.Stable
	}
	if postFindOptions.Update != nil {
		body["update"] = postFindOptions.Update
	}
	if postFindOptions.UseIndex != nil {
		body["use_index"] = postFindOptions.UseIndex
	}
	if postFindOptions.R != nil {
		body["r"] = postFindOptions.R
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// GetIndexesInformation : Retrieve information about all indexes
// When you make a GET request to `/db/_index`, you get a list of all indexes used by Cloudant Query in the database,
// including the primary index. In addition to the information available through this API, indexes are also stored in
// the `indexes` property of design documents.
func (cloudant *CloudantV1) GetIndexesInformation(getIndexesInformationOptions *GetIndexesInformationOptions) (result *IndexesInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetIndexesInformationWithContext(context.Background(), getIndexesInformationOptions)
}

// GetIndexesInformationWithContext is an alternate form of the GetIndexesInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetIndexesInformationWithContext(ctx context.Context, getIndexesInformationOptions *GetIndexesInformationOptions) (result *IndexesInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getIndexesInformationOptions, "getIndexesInformationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getIndexesInformationOptions, "getIndexesInformationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getIndexesInformationOptions.Db,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_index`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getIndexesInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetIndexesInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalIndexesInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostIndex : Create a new index on a database
// Create a new index on a database.
func (cloudant *CloudantV1) PostIndex(postIndexOptions *PostIndexOptions) (result *IndexResult, response *core.DetailedResponse, err error) {
	return cloudant.PostIndexWithContext(context.Background(), postIndexOptions)
}

// PostIndexWithContext is an alternate form of the PostIndex method which supports a Context parameter
func (cloudant *CloudantV1) PostIndexWithContext(ctx context.Context, postIndexOptions *PostIndexOptions) (result *IndexResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postIndexOptions, "postIndexOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postIndexOptions, "postIndexOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postIndexOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_index`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postIndexOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostIndex")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postIndexOptions.Index != nil {
		body["index"] = postIndexOptions.Index
	}
	if postIndexOptions.Ddoc != nil {
		body["ddoc"] = postIndexOptions.Ddoc
	}
	if postIndexOptions.Def != nil {
		body["def"] = postIndexOptions.Def
	}
	if postIndexOptions.Name != nil {
		body["name"] = postIndexOptions.Name
	}
	if postIndexOptions.Partitioned != nil {
		body["partitioned"] = postIndexOptions.Partitioned
	}
	if postIndexOptions.Type != nil {
		body["type"] = postIndexOptions.Type
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalIndexResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteIndex : Delete an index
func (cloudant *CloudantV1) DeleteIndex(deleteIndexOptions *DeleteIndexOptions) (result *Ok, response *core.DetailedResponse, err error) {
	return cloudant.DeleteIndexWithContext(context.Background(), deleteIndexOptions)
}

// DeleteIndexWithContext is an alternate form of the DeleteIndex method which supports a Context parameter
func (cloudant *CloudantV1) DeleteIndexWithContext(ctx context.Context, deleteIndexOptions *DeleteIndexOptions) (result *Ok, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteIndexOptions, "deleteIndexOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteIndexOptions, "deleteIndexOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *deleteIndexOptions.Db,
		"ddoc": *deleteIndexOptions.Ddoc,
		"type": *deleteIndexOptions.Type,
		"index": *deleteIndexOptions.Index,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_index/_design/{ddoc}/{type}/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteIndexOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "DeleteIndex")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostSearchAnalyze : Query tokenization of sample text
// Returns the results of analyzer tokenization of the provided sample text. This endpoint can be used for testing
// analyzer tokenization.
func (cloudant *CloudantV1) PostSearchAnalyze(postSearchAnalyzeOptions *PostSearchAnalyzeOptions) (result *SearchAnalyzeResult, response *core.DetailedResponse, err error) {
	return cloudant.PostSearchAnalyzeWithContext(context.Background(), postSearchAnalyzeOptions)
}

// PostSearchAnalyzeWithContext is an alternate form of the PostSearchAnalyze method which supports a Context parameter
func (cloudant *CloudantV1) PostSearchAnalyzeWithContext(ctx context.Context, postSearchAnalyzeOptions *PostSearchAnalyzeOptions) (result *SearchAnalyzeResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postSearchAnalyzeOptions, "postSearchAnalyzeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postSearchAnalyzeOptions, "postSearchAnalyzeOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_search_analyze`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range postSearchAnalyzeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostSearchAnalyze")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postSearchAnalyzeOptions.Analyzer != nil {
		body["analyzer"] = postSearchAnalyzeOptions.Analyzer
	}
	if postSearchAnalyzeOptions.Text != nil {
		body["text"] = postSearchAnalyzeOptions.Text
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSearchAnalyzeResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostSearch : Query a search index
// Search indexes, which are defined in design documents, allow databases to be queried by using Lucene Query Parser
// Syntax. An index function defines a search index, similar to a map function in MapReduce views. The index function
// decides what data to index and what data to store in the index. The advantage of using the HTTP `POST` method is that
// the query is submitted as a JSON object in the request body. This avoids the limitations of passing query options as
// URL query parameters of a `GET` request.
func (cloudant *CloudantV1) PostSearch(postSearchOptions *PostSearchOptions) (result *SearchResult, response *core.DetailedResponse, err error) {
	return cloudant.PostSearchWithContext(context.Background(), postSearchOptions)
}

// PostSearchWithContext is an alternate form of the PostSearch method which supports a Context parameter
func (cloudant *CloudantV1) PostSearchWithContext(ctx context.Context, postSearchOptions *PostSearchOptions) (result *SearchResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postSearchOptions, "postSearchOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postSearchOptions, "postSearchOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postSearchOptions.Db,
		"ddoc": *postSearchOptions.Ddoc,
		"index": *postSearchOptions.Index,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_search/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postSearchOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostSearch")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postSearchOptions.Query != nil {
		body["query"] = postSearchOptions.Query
	}
	if postSearchOptions.Bookmark != nil {
		body["bookmark"] = postSearchOptions.Bookmark
	}
	if postSearchOptions.HighlightFields != nil {
		body["highlight_fields"] = postSearchOptions.HighlightFields
	}
	if postSearchOptions.HighlightNumber != nil {
		body["highlight_number"] = postSearchOptions.HighlightNumber
	}
	if postSearchOptions.HighlightPostTag != nil {
		body["highlight_post_tag"] = postSearchOptions.HighlightPostTag
	}
	if postSearchOptions.HighlightPreTag != nil {
		body["highlight_pre_tag"] = postSearchOptions.HighlightPreTag
	}
	if postSearchOptions.HighlightSize != nil {
		body["highlight_size"] = postSearchOptions.HighlightSize
	}
	if postSearchOptions.IncludeDocs != nil {
		body["include_docs"] = postSearchOptions.IncludeDocs
	}
	if postSearchOptions.IncludeFields != nil {
		body["include_fields"] = postSearchOptions.IncludeFields
	}
	if postSearchOptions.Limit != nil {
		body["limit"] = postSearchOptions.Limit
	}
	if postSearchOptions.Sort != nil {
		body["sort"] = postSearchOptions.Sort
	}
	if postSearchOptions.Stale != nil {
		body["stale"] = postSearchOptions.Stale
	}
	if postSearchOptions.Counts != nil {
		body["counts"] = postSearchOptions.Counts
	}
	if postSearchOptions.Drilldown != nil {
		body["drilldown"] = postSearchOptions.Drilldown
	}
	if postSearchOptions.GroupField != nil {
		body["group_field"] = postSearchOptions.GroupField
	}
	if postSearchOptions.GroupLimit != nil {
		body["group_limit"] = postSearchOptions.GroupLimit
	}
	if postSearchOptions.GroupSort != nil {
		body["group_sort"] = postSearchOptions.GroupSort
	}
	if postSearchOptions.Ranges != nil {
		body["ranges"] = postSearchOptions.Ranges
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSearchResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostSearchAsStream : Query a search index as stream
// Search indexes, which are defined in design documents, allow databases to be queried by using Lucene Query Parser
// Syntax. An index function defines a search index, similar to a map function in MapReduce views. The index function
// decides what data to index and what data to store in the index. The advantage of using the HTTP `POST` method is that
// the query is submitted as a JSON object in the request body. This avoids the limitations of passing query options as
// URL query parameters of a `GET` request.
func (cloudant *CloudantV1) PostSearchAsStream(postSearchOptions *PostSearchOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.PostSearchAsStreamWithContext(context.Background(), postSearchOptions)
}

// PostSearchAsStreamWithContext is an alternate form of the PostSearchAsStream method which supports a Context parameter
func (cloudant *CloudantV1) PostSearchAsStreamWithContext(ctx context.Context, postSearchOptions *PostSearchOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postSearchOptions, "postSearchOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postSearchOptions, "postSearchOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postSearchOptions.Db,
		"ddoc": *postSearchOptions.Ddoc,
		"index": *postSearchOptions.Index,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_search/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postSearchOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostSearchAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postSearchOptions.Query != nil {
		body["query"] = postSearchOptions.Query
	}
	if postSearchOptions.Bookmark != nil {
		body["bookmark"] = postSearchOptions.Bookmark
	}
	if postSearchOptions.HighlightFields != nil {
		body["highlight_fields"] = postSearchOptions.HighlightFields
	}
	if postSearchOptions.HighlightNumber != nil {
		body["highlight_number"] = postSearchOptions.HighlightNumber
	}
	if postSearchOptions.HighlightPostTag != nil {
		body["highlight_post_tag"] = postSearchOptions.HighlightPostTag
	}
	if postSearchOptions.HighlightPreTag != nil {
		body["highlight_pre_tag"] = postSearchOptions.HighlightPreTag
	}
	if postSearchOptions.HighlightSize != nil {
		body["highlight_size"] = postSearchOptions.HighlightSize
	}
	if postSearchOptions.IncludeDocs != nil {
		body["include_docs"] = postSearchOptions.IncludeDocs
	}
	if postSearchOptions.IncludeFields != nil {
		body["include_fields"] = postSearchOptions.IncludeFields
	}
	if postSearchOptions.Limit != nil {
		body["limit"] = postSearchOptions.Limit
	}
	if postSearchOptions.Sort != nil {
		body["sort"] = postSearchOptions.Sort
	}
	if postSearchOptions.Stale != nil {
		body["stale"] = postSearchOptions.Stale
	}
	if postSearchOptions.Counts != nil {
		body["counts"] = postSearchOptions.Counts
	}
	if postSearchOptions.Drilldown != nil {
		body["drilldown"] = postSearchOptions.Drilldown
	}
	if postSearchOptions.GroupField != nil {
		body["group_field"] = postSearchOptions.GroupField
	}
	if postSearchOptions.GroupLimit != nil {
		body["group_limit"] = postSearchOptions.GroupLimit
	}
	if postSearchOptions.GroupSort != nil {
		body["group_sort"] = postSearchOptions.GroupSort
	}
	if postSearchOptions.Ranges != nil {
		body["ranges"] = postSearchOptions.Ranges
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// GetSearchInfo : Retrieve information about a search index
// Retrieve search index metadata information, such as the size of the index on disk.
func (cloudant *CloudantV1) GetSearchInfo(getSearchInfoOptions *GetSearchInfoOptions) (result *SearchInfoResult, response *core.DetailedResponse, err error) {
	return cloudant.GetSearchInfoWithContext(context.Background(), getSearchInfoOptions)
}

// GetSearchInfoWithContext is an alternate form of the GetSearchInfo method which supports a Context parameter
func (cloudant *CloudantV1) GetSearchInfoWithContext(ctx context.Context, getSearchInfoOptions *GetSearchInfoOptions) (result *SearchInfoResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSearchInfoOptions, "getSearchInfoOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSearchInfoOptions, "getSearchInfoOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getSearchInfoOptions.Db,
		"ddoc": *getSearchInfoOptions.Ddoc,
		"index": *getSearchInfoOptions.Index,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_search_info/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSearchInfoOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetSearchInfo")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSearchInfoResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetGeo : Query a geospatial index
// Executes a query against the requested geospatial index from the specified design document.
func (cloudant *CloudantV1) GetGeo(getGeoOptions *GetGeoOptions) (result *GeoResult, response *core.DetailedResponse, err error) {
	return cloudant.GetGeoWithContext(context.Background(), getGeoOptions)
}

// GetGeoWithContext is an alternate form of the GetGeo method which supports a Context parameter
func (cloudant *CloudantV1) GetGeoWithContext(ctx context.Context, getGeoOptions *GetGeoOptions) (result *GeoResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getGeoOptions, "getGeoOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getGeoOptions, "getGeoOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getGeoOptions.Db,
		"ddoc": *getGeoOptions.Ddoc,
		"index": *getGeoOptions.Index,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_geo/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getGeoOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetGeo")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getGeoOptions.Bbox != nil {
		builder.AddQuery("bbox", fmt.Sprint(*getGeoOptions.Bbox))
	}
	if getGeoOptions.Bookmark != nil {
		builder.AddQuery("bookmark", fmt.Sprint(*getGeoOptions.Bookmark))
	}
	if getGeoOptions.Format != nil {
		builder.AddQuery("format", fmt.Sprint(*getGeoOptions.Format))
	}
	if getGeoOptions.G != nil {
		builder.AddQuery("g", fmt.Sprint(*getGeoOptions.G))
	}
	if getGeoOptions.IncludeDocs != nil {
		builder.AddQuery("include_docs", fmt.Sprint(*getGeoOptions.IncludeDocs))
	}
	if getGeoOptions.Lat != nil {
		builder.AddQuery("lat", fmt.Sprint(*getGeoOptions.Lat))
	}
	if getGeoOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getGeoOptions.Limit))
	}
	if getGeoOptions.Lon != nil {
		builder.AddQuery("lon", fmt.Sprint(*getGeoOptions.Lon))
	}
	if getGeoOptions.Nearest != nil {
		builder.AddQuery("nearest", fmt.Sprint(*getGeoOptions.Nearest))
	}
	if getGeoOptions.Radius != nil {
		builder.AddQuery("radius", fmt.Sprint(*getGeoOptions.Radius))
	}
	if getGeoOptions.Rangex != nil {
		builder.AddQuery("rangex", fmt.Sprint(*getGeoOptions.Rangex))
	}
	if getGeoOptions.Rangey != nil {
		builder.AddQuery("rangey", fmt.Sprint(*getGeoOptions.Rangey))
	}
	if getGeoOptions.Relation != nil {
		builder.AddQuery("relation", fmt.Sprint(*getGeoOptions.Relation))
	}
	if getGeoOptions.Skip != nil {
		builder.AddQuery("skip", fmt.Sprint(*getGeoOptions.Skip))
	}
	if getGeoOptions.Stale != nil {
		builder.AddQuery("stale", fmt.Sprint(*getGeoOptions.Stale))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGeoResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetGeoAsStream : Query a geospatial index as stream
// Executes a query against the requested geospatial index from the specified design document.
func (cloudant *CloudantV1) GetGeoAsStream(getGeoOptions *GetGeoOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.GetGeoAsStreamWithContext(context.Background(), getGeoOptions)
}

// GetGeoAsStreamWithContext is an alternate form of the GetGeoAsStream method which supports a Context parameter
func (cloudant *CloudantV1) GetGeoAsStreamWithContext(ctx context.Context, getGeoOptions *GetGeoOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getGeoOptions, "getGeoOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getGeoOptions, "getGeoOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getGeoOptions.Db,
		"ddoc": *getGeoOptions.Ddoc,
		"index": *getGeoOptions.Index,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_geo/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getGeoOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetGeoAsStream")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getGeoOptions.Bbox != nil {
		builder.AddQuery("bbox", fmt.Sprint(*getGeoOptions.Bbox))
	}
	if getGeoOptions.Bookmark != nil {
		builder.AddQuery("bookmark", fmt.Sprint(*getGeoOptions.Bookmark))
	}
	if getGeoOptions.Format != nil {
		builder.AddQuery("format", fmt.Sprint(*getGeoOptions.Format))
	}
	if getGeoOptions.G != nil {
		builder.AddQuery("g", fmt.Sprint(*getGeoOptions.G))
	}
	if getGeoOptions.IncludeDocs != nil {
		builder.AddQuery("include_docs", fmt.Sprint(*getGeoOptions.IncludeDocs))
	}
	if getGeoOptions.Lat != nil {
		builder.AddQuery("lat", fmt.Sprint(*getGeoOptions.Lat))
	}
	if getGeoOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getGeoOptions.Limit))
	}
	if getGeoOptions.Lon != nil {
		builder.AddQuery("lon", fmt.Sprint(*getGeoOptions.Lon))
	}
	if getGeoOptions.Nearest != nil {
		builder.AddQuery("nearest", fmt.Sprint(*getGeoOptions.Nearest))
	}
	if getGeoOptions.Radius != nil {
		builder.AddQuery("radius", fmt.Sprint(*getGeoOptions.Radius))
	}
	if getGeoOptions.Rangex != nil {
		builder.AddQuery("rangex", fmt.Sprint(*getGeoOptions.Rangex))
	}
	if getGeoOptions.Rangey != nil {
		builder.AddQuery("rangey", fmt.Sprint(*getGeoOptions.Rangey))
	}
	if getGeoOptions.Relation != nil {
		builder.AddQuery("relation", fmt.Sprint(*getGeoOptions.Relation))
	}
	if getGeoOptions.Skip != nil {
		builder.AddQuery("skip", fmt.Sprint(*getGeoOptions.Skip))
	}
	if getGeoOptions.Stale != nil {
		builder.AddQuery("stale", fmt.Sprint(*getGeoOptions.Stale))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PostGeoCleanup : Cleanup old geospatial indexes
// Cleanup old geospatial indexes from disk that have been superseded by newer index builds.
func (cloudant *CloudantV1) PostGeoCleanup(postGeoCleanupOptions *PostGeoCleanupOptions) (result *Ok, response *core.DetailedResponse, err error) {
	return cloudant.PostGeoCleanupWithContext(context.Background(), postGeoCleanupOptions)
}

// PostGeoCleanupWithContext is an alternate form of the PostGeoCleanup method which supports a Context parameter
func (cloudant *CloudantV1) PostGeoCleanupWithContext(ctx context.Context, postGeoCleanupOptions *PostGeoCleanupOptions) (result *Ok, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postGeoCleanupOptions, "postGeoCleanupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postGeoCleanupOptions, "postGeoCleanupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postGeoCleanupOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_geo_cleanup`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postGeoCleanupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostGeoCleanup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetGeoIndexInformation : Retrieve information about a geospatial index
func (cloudant *CloudantV1) GetGeoIndexInformation(getGeoIndexInformationOptions *GetGeoIndexInformationOptions) (result *GeoIndexInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetGeoIndexInformationWithContext(context.Background(), getGeoIndexInformationOptions)
}

// GetGeoIndexInformationWithContext is an alternate form of the GetGeoIndexInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetGeoIndexInformationWithContext(ctx context.Context, getGeoIndexInformationOptions *GetGeoIndexInformationOptions) (result *GeoIndexInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getGeoIndexInformationOptions, "getGeoIndexInformationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getGeoIndexInformationOptions, "getGeoIndexInformationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getGeoIndexInformationOptions.Db,
		"ddoc": *getGeoIndexInformationOptions.Ddoc,
		"index": *getGeoIndexInformationOptions.Index,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_design/{ddoc}/_geo_info/{index}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getGeoIndexInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetGeoIndexInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGeoIndexInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// HeadReplicationDocument : Retrieve the HTTP headers for a replication document
// Retrieves the HTTP headers containing minimal amount of information about the specified replication document from the
// `_replicator` database.  The method supports the same query arguments as the `GET /_replicator/{doc_id}` method, but
// only headers like content length and the revision (ETag header) are returned.
func (cloudant *CloudantV1) HeadReplicationDocument(headReplicationDocumentOptions *HeadReplicationDocumentOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadReplicationDocumentWithContext(context.Background(), headReplicationDocumentOptions)
}

// HeadReplicationDocumentWithContext is an alternate form of the HeadReplicationDocument method which supports a Context parameter
func (cloudant *CloudantV1) HeadReplicationDocumentWithContext(ctx context.Context, headReplicationDocumentOptions *HeadReplicationDocumentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(headReplicationDocumentOptions, "headReplicationDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(headReplicationDocumentOptions, "headReplicationDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"doc_id": *headReplicationDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_replicator/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range headReplicationDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadReplicationDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if headReplicationDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*headReplicationDocumentOptions.IfNoneMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// HeadSchedulerDocument : Retrieve HTTP headers for a replication scheduler document
// Retrieves the HTTP headers containing minimal amount of information about the specified replication scheduler
// document.  Since the response body is empty, using the HEAD method is a lightweight way to check if the replication
// scheduler document exists or not.
func (cloudant *CloudantV1) HeadSchedulerDocument(headSchedulerDocumentOptions *HeadSchedulerDocumentOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadSchedulerDocumentWithContext(context.Background(), headSchedulerDocumentOptions)
}

// HeadSchedulerDocumentWithContext is an alternate form of the HeadSchedulerDocument method which supports a Context parameter
func (cloudant *CloudantV1) HeadSchedulerDocumentWithContext(ctx context.Context, headSchedulerDocumentOptions *HeadSchedulerDocumentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(headSchedulerDocumentOptions, "headSchedulerDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(headSchedulerDocumentOptions, "headSchedulerDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"doc_id": *headSchedulerDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_scheduler/docs/_replicator/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range headSchedulerDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadSchedulerDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// HeadSchedulerJob : Retrieve the HTTP headers for a replication scheduler job
// Returns the HTTP headers that contain a minimal amount of information about the specified replication task. Only the
// header information is returned.
func (cloudant *CloudantV1) HeadSchedulerJob(headSchedulerJobOptions *HeadSchedulerJobOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadSchedulerJobWithContext(context.Background(), headSchedulerJobOptions)
}

// HeadSchedulerJobWithContext is an alternate form of the HeadSchedulerJob method which supports a Context parameter
func (cloudant *CloudantV1) HeadSchedulerJobWithContext(ctx context.Context, headSchedulerJobOptions *HeadSchedulerJobOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(headSchedulerJobOptions, "headSchedulerJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(headSchedulerJobOptions, "headSchedulerJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"job_id": *headSchedulerJobOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_scheduler/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range headSchedulerJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadSchedulerJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// DeleteReplicationDocument : Cancel a replication
// Cancels a replication by deleting the document that describes it from the `_replicator` database.
func (cloudant *CloudantV1) DeleteReplicationDocument(deleteReplicationDocumentOptions *DeleteReplicationDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.DeleteReplicationDocumentWithContext(context.Background(), deleteReplicationDocumentOptions)
}

// DeleteReplicationDocumentWithContext is an alternate form of the DeleteReplicationDocument method which supports a Context parameter
func (cloudant *CloudantV1) DeleteReplicationDocumentWithContext(ctx context.Context, deleteReplicationDocumentOptions *DeleteReplicationDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteReplicationDocumentOptions, "deleteReplicationDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteReplicationDocumentOptions, "deleteReplicationDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"doc_id": *deleteReplicationDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_replicator/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteReplicationDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "DeleteReplicationDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteReplicationDocumentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deleteReplicationDocumentOptions.IfMatch))
	}

	if deleteReplicationDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*deleteReplicationDocumentOptions.Batch))
	}
	if deleteReplicationDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*deleteReplicationDocumentOptions.Rev))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReplicationDocument : Retrieve a replication document
// Retrieves a replication document from the `_replicator` database to view the configuration of the replication. The
// status of the replication is no longer recorded in the document but can be checked via the replication scheduler.
func (cloudant *CloudantV1) GetReplicationDocument(getReplicationDocumentOptions *GetReplicationDocumentOptions) (result *ReplicationDocument, response *core.DetailedResponse, err error) {
	return cloudant.GetReplicationDocumentWithContext(context.Background(), getReplicationDocumentOptions)
}

// GetReplicationDocumentWithContext is an alternate form of the GetReplicationDocument method which supports a Context parameter
func (cloudant *CloudantV1) GetReplicationDocumentWithContext(ctx context.Context, getReplicationDocumentOptions *GetReplicationDocumentOptions) (result *ReplicationDocument, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReplicationDocumentOptions, "getReplicationDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReplicationDocumentOptions, "getReplicationDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"doc_id": *getReplicationDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_replicator/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReplicationDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetReplicationDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getReplicationDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*getReplicationDocumentOptions.IfNoneMatch))
	}

	if getReplicationDocumentOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*getReplicationDocumentOptions.Attachments))
	}
	if getReplicationDocumentOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*getReplicationDocumentOptions.AttEncodingInfo))
	}
	if getReplicationDocumentOptions.Conflicts != nil {
		builder.AddQuery("conflicts", fmt.Sprint(*getReplicationDocumentOptions.Conflicts))
	}
	if getReplicationDocumentOptions.DeletedConflicts != nil {
		builder.AddQuery("deleted_conflicts", fmt.Sprint(*getReplicationDocumentOptions.DeletedConflicts))
	}
	if getReplicationDocumentOptions.Latest != nil {
		builder.AddQuery("latest", fmt.Sprint(*getReplicationDocumentOptions.Latest))
	}
	if getReplicationDocumentOptions.LocalSeq != nil {
		builder.AddQuery("local_seq", fmt.Sprint(*getReplicationDocumentOptions.LocalSeq))
	}
	if getReplicationDocumentOptions.Meta != nil {
		builder.AddQuery("meta", fmt.Sprint(*getReplicationDocumentOptions.Meta))
	}
	if getReplicationDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*getReplicationDocumentOptions.Rev))
	}
	if getReplicationDocumentOptions.Revs != nil {
		builder.AddQuery("revs", fmt.Sprint(*getReplicationDocumentOptions.Revs))
	}
	if getReplicationDocumentOptions.RevsInfo != nil {
		builder.AddQuery("revs_info", fmt.Sprint(*getReplicationDocumentOptions.RevsInfo))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReplicationDocument)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutReplicationDocument : Create or modify a replication using a replication document
// Creates or modifies a document in the `_replicator` database to start a new replication or to edit an existing
// replication.
func (cloudant *CloudantV1) PutReplicationDocument(putReplicationDocumentOptions *PutReplicationDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.PutReplicationDocumentWithContext(context.Background(), putReplicationDocumentOptions)
}

// PutReplicationDocumentWithContext is an alternate form of the PutReplicationDocument method which supports a Context parameter
func (cloudant *CloudantV1) PutReplicationDocumentWithContext(ctx context.Context, putReplicationDocumentOptions *PutReplicationDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putReplicationDocumentOptions, "putReplicationDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putReplicationDocumentOptions, "putReplicationDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"doc_id": *putReplicationDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_replicator/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putReplicationDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutReplicationDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if putReplicationDocumentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*putReplicationDocumentOptions.IfMatch))
	}

	if putReplicationDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*putReplicationDocumentOptions.Batch))
	}
	if putReplicationDocumentOptions.NewEdits != nil {
		builder.AddQuery("new_edits", fmt.Sprint(*putReplicationDocumentOptions.NewEdits))
	}
	if putReplicationDocumentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*putReplicationDocumentOptions.Rev))
	}

	_, err = builder.SetBodyContentJSON(putReplicationDocumentOptions.ReplicationDocument)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSchedulerDocs : Retrieve replication scheduler documents
// Lists replication documents, including information about all documents, even the ones in a completed or failed state.
// For each document, the endpoint returns the document ID, database, replication ID, source and target, and other
// information.
func (cloudant *CloudantV1) GetSchedulerDocs(getSchedulerDocsOptions *GetSchedulerDocsOptions) (result *SchedulerDocsResult, response *core.DetailedResponse, err error) {
	return cloudant.GetSchedulerDocsWithContext(context.Background(), getSchedulerDocsOptions)
}

// GetSchedulerDocsWithContext is an alternate form of the GetSchedulerDocs method which supports a Context parameter
func (cloudant *CloudantV1) GetSchedulerDocsWithContext(ctx context.Context, getSchedulerDocsOptions *GetSchedulerDocsOptions) (result *SchedulerDocsResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSchedulerDocsOptions, "getSchedulerDocsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_scheduler/docs`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSchedulerDocsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetSchedulerDocs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getSchedulerDocsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getSchedulerDocsOptions.Limit))
	}
	if getSchedulerDocsOptions.Skip != nil {
		builder.AddQuery("skip", fmt.Sprint(*getSchedulerDocsOptions.Skip))
	}
	if getSchedulerDocsOptions.States != nil {
		builder.AddQuery("states", strings.Join(getSchedulerDocsOptions.States, ","))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSchedulerDocsResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSchedulerDocument : Retrieve a replication scheduler document
// Retrieves information about a replication document from the replicator database. The endpoint returns the document
// ID, database, replication ID, source and target, and other information.
func (cloudant *CloudantV1) GetSchedulerDocument(getSchedulerDocumentOptions *GetSchedulerDocumentOptions) (result *SchedulerDocument, response *core.DetailedResponse, err error) {
	return cloudant.GetSchedulerDocumentWithContext(context.Background(), getSchedulerDocumentOptions)
}

// GetSchedulerDocumentWithContext is an alternate form of the GetSchedulerDocument method which supports a Context parameter
func (cloudant *CloudantV1) GetSchedulerDocumentWithContext(ctx context.Context, getSchedulerDocumentOptions *GetSchedulerDocumentOptions) (result *SchedulerDocument, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSchedulerDocumentOptions, "getSchedulerDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSchedulerDocumentOptions, "getSchedulerDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"doc_id": *getSchedulerDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_scheduler/docs/_replicator/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSchedulerDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetSchedulerDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSchedulerDocument)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSchedulerJobs : Retrieve replication scheduler jobs
// Retrieves information about replications that were created via `/_replicate` endpoint, as well as those created from
// replication documents. It doesn't include replications that completed or failed to start because replication
// documents were malformed. Each job description includes source and target information, replication ID, history of
// recent events, and other information.
func (cloudant *CloudantV1) GetSchedulerJobs(getSchedulerJobsOptions *GetSchedulerJobsOptions) (result *SchedulerJobsResult, response *core.DetailedResponse, err error) {
	return cloudant.GetSchedulerJobsWithContext(context.Background(), getSchedulerJobsOptions)
}

// GetSchedulerJobsWithContext is an alternate form of the GetSchedulerJobs method which supports a Context parameter
func (cloudant *CloudantV1) GetSchedulerJobsWithContext(ctx context.Context, getSchedulerJobsOptions *GetSchedulerJobsOptions) (result *SchedulerJobsResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSchedulerJobsOptions, "getSchedulerJobsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_scheduler/jobs`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSchedulerJobsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetSchedulerJobs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getSchedulerJobsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getSchedulerJobsOptions.Limit))
	}
	if getSchedulerJobsOptions.Skip != nil {
		builder.AddQuery("skip", fmt.Sprint(*getSchedulerJobsOptions.Skip))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSchedulerJobsResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSchedulerJob : Retrieve a replication scheduler job
// Retrieves the state of a single replication task based on its replication ID.
func (cloudant *CloudantV1) GetSchedulerJob(getSchedulerJobOptions *GetSchedulerJobOptions) (result *SchedulerJob, response *core.DetailedResponse, err error) {
	return cloudant.GetSchedulerJobWithContext(context.Background(), getSchedulerJobOptions)
}

// GetSchedulerJobWithContext is an alternate form of the GetSchedulerJob method which supports a Context parameter
func (cloudant *CloudantV1) GetSchedulerJobWithContext(ctx context.Context, getSchedulerJobOptions *GetSchedulerJobOptions) (result *SchedulerJob, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSchedulerJobOptions, "getSchedulerJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSchedulerJobOptions, "getSchedulerJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"job_id": *getSchedulerJobOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_scheduler/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSchedulerJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetSchedulerJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSchedulerJob)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSessionInformation : Retrieve current session cookie information
// Retrieves information about the authenticated user's session.
func (cloudant *CloudantV1) GetSessionInformation(getSessionInformationOptions *GetSessionInformationOptions) (result *SessionInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetSessionInformationWithContext(context.Background(), getSessionInformationOptions)
}

// GetSessionInformationWithContext is an alternate form of the GetSessionInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetSessionInformationWithContext(ctx context.Context, getSessionInformationOptions *GetSessionInformationOptions) (result *SessionInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSessionInformationOptions, "getSessionInformationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_session`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSessionInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetSessionInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSessionInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSecurity : Retrieve database permissions information
// See who has permission to read, write, and manage the database. The credentials you use to log in to the dashboard
// automatically include `_admin` permissions to all databases you create. Everyone and everything else, including users
// you share databases with and API keys you create, must be given a permission level explicitly.
func (cloudant *CloudantV1) GetSecurity(getSecurityOptions *GetSecurityOptions) (result *Security, response *core.DetailedResponse, err error) {
	return cloudant.GetSecurityWithContext(context.Background(), getSecurityOptions)
}

// GetSecurityWithContext is an alternate form of the GetSecurity method which supports a Context parameter
func (cloudant *CloudantV1) GetSecurityWithContext(ctx context.Context, getSecurityOptions *GetSecurityOptions) (result *Security, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecurityOptions, "getSecurityOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSecurityOptions, "getSecurityOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getSecurityOptions.Db,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_security`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSecurityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetSecurity")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecurity)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutSecurity : Modify database permissions
// Modify who has permission to read, write, or manage a database. This endpoint can be used to modify both Cloudant and
// CouchDB related permissions. Be careful: by removing a Cloudant API key, a member or an admin from the list of users
// that have access permissions, you remove it from the list of users that have access to the database.
//
// ### Note about nobody role
//
// The `nobody` username applies to all unauthenticated connection attempts. For example, if an application tries to
// read data from a database, but did not identify itself, the task can continue only if the `nobody` user has the role
// `_reader`.
func (cloudant *CloudantV1) PutSecurity(putSecurityOptions *PutSecurityOptions) (result *Ok, response *core.DetailedResponse, err error) {
	return cloudant.PutSecurityWithContext(context.Background(), putSecurityOptions)
}

// PutSecurityWithContext is an alternate form of the PutSecurity method which supports a Context parameter
func (cloudant *CloudantV1) PutSecurityWithContext(ctx context.Context, putSecurityOptions *PutSecurityOptions) (result *Ok, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putSecurityOptions, "putSecurityOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putSecurityOptions, "putSecurityOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *putSecurityOptions.Db,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_security`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putSecurityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutSecurity")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if putSecurityOptions.Admins != nil {
		body["admins"] = putSecurityOptions.Admins
	}
	if putSecurityOptions.Members != nil {
		body["members"] = putSecurityOptions.Members
	}
	if putSecurityOptions.Cloudant != nil {
		body["cloudant"] = putSecurityOptions.Cloudant
	}
	if putSecurityOptions.CouchdbAuthOnly != nil {
		body["couchdb_auth_only"] = putSecurityOptions.CouchdbAuthOnly
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostApiKeys : Generates API keys for apps or persons to enable database access
// Generates API keys to enable database access for a person or application, but without creating a new IBM Cloudant
// account for that person or application. An API key is a randomly generated username and password. The key is given
// the wanted access permissions for a database.
func (cloudant *CloudantV1) PostApiKeys(postApiKeysOptions *PostApiKeysOptions) (result *ApiKeysResult, response *core.DetailedResponse, err error) {
	return cloudant.PostApiKeysWithContext(context.Background(), postApiKeysOptions)
}

// PostApiKeysWithContext is an alternate form of the PostApiKeys method which supports a Context parameter
func (cloudant *CloudantV1) PostApiKeysWithContext(ctx context.Context, postApiKeysOptions *PostApiKeysOptions) (result *ApiKeysResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(postApiKeysOptions, "postApiKeysOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/api_keys`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range postApiKeysOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostApiKeys")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApiKeysResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutCloudantSecurityConfiguration : Modify only Cloudant related database permissions
// Modify only Cloudant related permissions to database. Be careful: by removing an API key from the list, you remove
// the API key from the list of users that have access to the database.
//
// ### Note about nobody role
//
// The `nobody` username applies to all unauthenticated connection attempts. For example, if an application tries to
// read data from a database, but did not identify itself, the task can continue only if the `nobody` user has the role
// `_reader`.
func (cloudant *CloudantV1) PutCloudantSecurityConfiguration(putCloudantSecurityConfigurationOptions *PutCloudantSecurityConfigurationOptions) (result *Ok, response *core.DetailedResponse, err error) {
	return cloudant.PutCloudantSecurityConfigurationWithContext(context.Background(), putCloudantSecurityConfigurationOptions)
}

// PutCloudantSecurityConfigurationWithContext is an alternate form of the PutCloudantSecurityConfiguration method which supports a Context parameter
func (cloudant *CloudantV1) PutCloudantSecurityConfigurationWithContext(ctx context.Context, putCloudantSecurityConfigurationOptions *PutCloudantSecurityConfigurationOptions) (result *Ok, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putCloudantSecurityConfigurationOptions, "putCloudantSecurityConfigurationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putCloudantSecurityConfigurationOptions, "putCloudantSecurityConfigurationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *putCloudantSecurityConfigurationOptions.Db,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/db/{db}/_security`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putCloudantSecurityConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutCloudantSecurityConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if putCloudantSecurityConfigurationOptions.Cloudant != nil {
		body["cloudant"] = putCloudantSecurityConfigurationOptions.Cloudant
	}
	if putCloudantSecurityConfigurationOptions.Admins != nil {
		body["admins"] = putCloudantSecurityConfigurationOptions.Admins
	}
	if putCloudantSecurityConfigurationOptions.Members != nil {
		body["members"] = putCloudantSecurityConfigurationOptions.Members
	}
	if putCloudantSecurityConfigurationOptions.CouchdbAuthOnly != nil {
		body["couchdb_auth_only"] = putCloudantSecurityConfigurationOptions.CouchdbAuthOnly
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCorsInformation : Retrieve CORS configuration information
// Lists all Cross-origin resource sharing (CORS) configuration. CORS defines a way in which the browser and the server
// interact to determine whether or not to allow the request.
func (cloudant *CloudantV1) GetCorsInformation(getCorsInformationOptions *GetCorsInformationOptions) (result *CorsInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetCorsInformationWithContext(context.Background(), getCorsInformationOptions)
}

// GetCorsInformationWithContext is an alternate form of the GetCorsInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetCorsInformationWithContext(ctx context.Context, getCorsInformationOptions *GetCorsInformationOptions) (result *CorsInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getCorsInformationOptions, "getCorsInformationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/user/config/cors`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCorsInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetCorsInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCorsInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutCorsConfiguration : Modify CORS configuration
// Sets the CORS configuration. The configuration applies to all databases and all account level endpoints in your
// account.
func (cloudant *CloudantV1) PutCorsConfiguration(putCorsConfigurationOptions *PutCorsConfigurationOptions) (result *Ok, response *core.DetailedResponse, err error) {
	return cloudant.PutCorsConfigurationWithContext(context.Background(), putCorsConfigurationOptions)
}

// PutCorsConfigurationWithContext is an alternate form of the PutCorsConfiguration method which supports a Context parameter
func (cloudant *CloudantV1) PutCorsConfigurationWithContext(ctx context.Context, putCorsConfigurationOptions *PutCorsConfigurationOptions) (result *Ok, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putCorsConfigurationOptions, "putCorsConfigurationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putCorsConfigurationOptions, "putCorsConfigurationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/user/config/cors`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range putCorsConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutCorsConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if putCorsConfigurationOptions.Origins != nil {
		body["origins"] = putCorsConfigurationOptions.Origins
	}
	if putCorsConfigurationOptions.AllowCredentials != nil {
		body["allow_credentials"] = putCorsConfigurationOptions.AllowCredentials
	}
	if putCorsConfigurationOptions.EnableCors != nil {
		body["enable_cors"] = putCorsConfigurationOptions.EnableCors
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// HeadAttachment : Retrieve the HTTP headers for an attachment
// Returns the HTTP headers that contain a minimal amount of information about the specified attachment. This method
// supports the same query arguments as the `GET /{db}/{doc_id}/{attachment_name}` method, but only the header
// information (including attachment size, encoding, and the MD5 hash as an ETag), is returned.
func (cloudant *CloudantV1) HeadAttachment(headAttachmentOptions *HeadAttachmentOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadAttachmentWithContext(context.Background(), headAttachmentOptions)
}

// HeadAttachmentWithContext is an alternate form of the HeadAttachment method which supports a Context parameter
func (cloudant *CloudantV1) HeadAttachmentWithContext(ctx context.Context, headAttachmentOptions *HeadAttachmentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(headAttachmentOptions, "headAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(headAttachmentOptions, "headAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *headAttachmentOptions.Db,
		"doc_id": *headAttachmentOptions.DocID,
		"attachment_name": *headAttachmentOptions.AttachmentName,
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}/{attachment_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range headAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if headAttachmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*headAttachmentOptions.IfMatch))
	}
	if headAttachmentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*headAttachmentOptions.IfNoneMatch))
	}

	if headAttachmentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*headAttachmentOptions.Rev))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// DeleteAttachment : Delete an attachment
// Deletes the attachment with the filename, `{attachment_name}`, from the specified doc. You must supply the `rev`
// query parameter or `If-Match` header with the current revision to delete the attachment.
func (cloudant *CloudantV1) DeleteAttachment(deleteAttachmentOptions *DeleteAttachmentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.DeleteAttachmentWithContext(context.Background(), deleteAttachmentOptions)
}

// DeleteAttachmentWithContext is an alternate form of the DeleteAttachment method which supports a Context parameter
func (cloudant *CloudantV1) DeleteAttachmentWithContext(ctx context.Context, deleteAttachmentOptions *DeleteAttachmentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAttachmentOptions, "deleteAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteAttachmentOptions, "deleteAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *deleteAttachmentOptions.Db,
		"doc_id": *deleteAttachmentOptions.DocID,
		"attachment_name": *deleteAttachmentOptions.AttachmentName,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}/{attachment_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "DeleteAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteAttachmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*deleteAttachmentOptions.IfMatch))
	}

	if deleteAttachmentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*deleteAttachmentOptions.Rev))
	}
	if deleteAttachmentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*deleteAttachmentOptions.Batch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAttachment : Retrieve an attachment
// Returns the file attachment that is associated with the document. The raw data of the associated attachment is
// returned, just as if you were accessing a static file. The returned Content-Type header is the same as the content
// type set when the document attachment was submitted to the database.
func (cloudant *CloudantV1) GetAttachment(getAttachmentOptions *GetAttachmentOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return cloudant.GetAttachmentWithContext(context.Background(), getAttachmentOptions)
}

// GetAttachmentWithContext is an alternate form of the GetAttachment method which supports a Context parameter
func (cloudant *CloudantV1) GetAttachmentWithContext(ctx context.Context, getAttachmentOptions *GetAttachmentOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAttachmentOptions, "getAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAttachmentOptions, "getAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getAttachmentOptions.Db,
		"doc_id": *getAttachmentOptions.DocID,
		"attachment_name": *getAttachmentOptions.AttachmentName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}/{attachment_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "*/*")
	if getAttachmentOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*getAttachmentOptions.Accept))
	}
	if getAttachmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*getAttachmentOptions.IfMatch))
	}
	if getAttachmentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*getAttachmentOptions.IfNoneMatch))
	}
	if getAttachmentOptions.Range != nil {
		builder.AddHeader("Range", fmt.Sprint(*getAttachmentOptions.Range))
	}

	if getAttachmentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*getAttachmentOptions.Rev))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, &result)

	return
}

// PutAttachment : Create or modify an attachment
// Uploads the supplied content as an attachment to the specified document. The attachment name that you provide must be
// a URL encoded string. You must supply the Content-Type header, and for an existing document, you must also supply
// either the `rev` query argument or the `If-Match` HTTP header. If you omit the revision, a new, otherwise empty,
// document is created with the provided attachment, or a conflict occurs. If the uploaded attachment uses an existing
// attachment name in the remote database, it updates the corresponding stored content of the database. Since you must
// supply the revision information to add an attachment to the document, this serves as validation to update the
// existing attachment.
func (cloudant *CloudantV1) PutAttachment(putAttachmentOptions *PutAttachmentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.PutAttachmentWithContext(context.Background(), putAttachmentOptions)
}

// PutAttachmentWithContext is an alternate form of the PutAttachment method which supports a Context parameter
func (cloudant *CloudantV1) PutAttachmentWithContext(ctx context.Context, putAttachmentOptions *PutAttachmentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putAttachmentOptions, "putAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putAttachmentOptions, "putAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *putAttachmentOptions.Db,
		"doc_id": *putAttachmentOptions.DocID,
		"attachment_name": *putAttachmentOptions.AttachmentName,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/{doc_id}/{attachment_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if putAttachmentOptions.ContentType != nil {
		builder.AddHeader("Content-Type", fmt.Sprint(*putAttachmentOptions.ContentType))
	}
	if putAttachmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*putAttachmentOptions.IfMatch))
	}

	if putAttachmentOptions.Rev != nil {
		builder.AddQuery("rev", fmt.Sprint(*putAttachmentOptions.Rev))
	}

	_, err = builder.SetBodyContent(core.StringNilMapper(putAttachmentOptions.ContentType), nil, nil, putAttachmentOptions.Attachment)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// HeadLocalDocument : Retrieve HTTP headers for a local document
// Retrieves the HTTP headers containing minimal amount of information about the specified local document. Since the
// response body is empty, using the HEAD method is a lightweight way to check if the local document exists or not.
func (cloudant *CloudantV1) HeadLocalDocument(headLocalDocumentOptions *HeadLocalDocumentOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadLocalDocumentWithContext(context.Background(), headLocalDocumentOptions)
}

// HeadLocalDocumentWithContext is an alternate form of the HeadLocalDocument method which supports a Context parameter
func (cloudant *CloudantV1) HeadLocalDocumentWithContext(ctx context.Context, headLocalDocumentOptions *HeadLocalDocumentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(headLocalDocumentOptions, "headLocalDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(headLocalDocumentOptions, "headLocalDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *headLocalDocumentOptions.Db,
		"doc_id": *headLocalDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_local/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range headLocalDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadLocalDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if headLocalDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*headLocalDocumentOptions.IfNoneMatch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// DeleteLocalDocument : Delete a local document
// Deletes the specified local document. The semantics are identical to deleting a standard document in the specified
// database, except that the document is not replicated.
func (cloudant *CloudantV1) DeleteLocalDocument(deleteLocalDocumentOptions *DeleteLocalDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.DeleteLocalDocumentWithContext(context.Background(), deleteLocalDocumentOptions)
}

// DeleteLocalDocumentWithContext is an alternate form of the DeleteLocalDocument method which supports a Context parameter
func (cloudant *CloudantV1) DeleteLocalDocumentWithContext(ctx context.Context, deleteLocalDocumentOptions *DeleteLocalDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLocalDocumentOptions, "deleteLocalDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLocalDocumentOptions, "deleteLocalDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *deleteLocalDocumentOptions.Db,
		"doc_id": *deleteLocalDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_local/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteLocalDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "DeleteLocalDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if deleteLocalDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*deleteLocalDocumentOptions.Batch))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLocalDocument : Retrieve a local document
// Retrieves the specified local document. The semantics are identical to accessing a standard document in the specified
// database, except that the document is not replicated.
func (cloudant *CloudantV1) GetLocalDocument(getLocalDocumentOptions *GetLocalDocumentOptions) (result *Document, response *core.DetailedResponse, err error) {
	return cloudant.GetLocalDocumentWithContext(context.Background(), getLocalDocumentOptions)
}

// GetLocalDocumentWithContext is an alternate form of the GetLocalDocument method which supports a Context parameter
func (cloudant *CloudantV1) GetLocalDocumentWithContext(ctx context.Context, getLocalDocumentOptions *GetLocalDocumentOptions) (result *Document, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLocalDocumentOptions, "getLocalDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLocalDocumentOptions, "getLocalDocumentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getLocalDocumentOptions.Db,
		"doc_id": *getLocalDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_local/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLocalDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetLocalDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getLocalDocumentOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*getLocalDocumentOptions.Accept))
	}
	if getLocalDocumentOptions.IfNoneMatch != nil {
		builder.AddHeader("If-None-Match", fmt.Sprint(*getLocalDocumentOptions.IfNoneMatch))
	}

	if getLocalDocumentOptions.Attachments != nil {
		builder.AddQuery("attachments", fmt.Sprint(*getLocalDocumentOptions.Attachments))
	}
	if getLocalDocumentOptions.AttEncodingInfo != nil {
		builder.AddQuery("att_encoding_info", fmt.Sprint(*getLocalDocumentOptions.AttEncodingInfo))
	}
	if getLocalDocumentOptions.LocalSeq != nil {
		builder.AddQuery("local_seq", fmt.Sprint(*getLocalDocumentOptions.LocalSeq))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocument)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutLocalDocument : Create or modify a local document
// Stores the specified local document. The semantics are identical to storing a standard document in the specified
// database, except that the document is not replicated.
func (cloudant *CloudantV1) PutLocalDocument(putLocalDocumentOptions *PutLocalDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	return cloudant.PutLocalDocumentWithContext(context.Background(), putLocalDocumentOptions)
}

// PutLocalDocumentWithContext is an alternate form of the PutLocalDocument method which supports a Context parameter
func (cloudant *CloudantV1) PutLocalDocumentWithContext(ctx context.Context, putLocalDocumentOptions *PutLocalDocumentOptions) (result *DocumentResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putLocalDocumentOptions, "putLocalDocumentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putLocalDocumentOptions, "putLocalDocumentOptions")
	if err != nil {
		return
	}

	if putLocalDocumentOptions.Document != nil && putLocalDocumentOptions.ContentType == nil {
		putLocalDocumentOptions.SetContentType("application/json")
	}

	pathParamsMap := map[string]string{
		"db": *putLocalDocumentOptions.Db,
		"doc_id": *putLocalDocumentOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_local/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putLocalDocumentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PutLocalDocument")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if putLocalDocumentOptions.ContentType != nil {
		builder.AddHeader("Content-Type", fmt.Sprint(*putLocalDocumentOptions.ContentType))
	}

	if putLocalDocumentOptions.Batch != nil {
		builder.AddQuery("batch", fmt.Sprint(*putLocalDocumentOptions.Batch))
	}

	_, err = builder.SetBodyContent(core.StringNilMapper(putLocalDocumentOptions.ContentType), putLocalDocumentOptions.Document, nil, putLocalDocumentOptions.Body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostRevsDiff : Query the document revisions and possible ancestors missing from the database
// The replicator is the primary user of this operation. After receiving a set of new revision IDs from the source
// database, the replicator sends this set to the destination database's `_revs_diff` to find out which of them already
// exists there. It can then avoid fetching and sending already-known document bodies.
func (cloudant *CloudantV1) PostRevsDiff(postRevsDiffOptions *PostRevsDiffOptions) (result map[string]RevsDiff, response *core.DetailedResponse, err error) {
	return cloudant.PostRevsDiffWithContext(context.Background(), postRevsDiffOptions)
}

// PostRevsDiffWithContext is an alternate form of the PostRevsDiff method which supports a Context parameter
func (cloudant *CloudantV1) PostRevsDiffWithContext(ctx context.Context, postRevsDiffOptions *PostRevsDiffOptions) (result map[string]RevsDiff, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postRevsDiffOptions, "postRevsDiffOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postRevsDiffOptions, "postRevsDiffOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *postRevsDiffOptions.Db,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_revs_diff`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postRevsDiffOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostRevsDiff")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(postRevsDiffOptions.DocumentRevisions)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRevsDiff)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetShardsInformation : Retrieve shard information
// List each shard range and the corresponding replicas for a specified database.
func (cloudant *CloudantV1) GetShardsInformation(getShardsInformationOptions *GetShardsInformationOptions) (result *ShardsInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetShardsInformationWithContext(context.Background(), getShardsInformationOptions)
}

// GetShardsInformationWithContext is an alternate form of the GetShardsInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetShardsInformationWithContext(ctx context.Context, getShardsInformationOptions *GetShardsInformationOptions) (result *ShardsInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getShardsInformationOptions, "getShardsInformationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getShardsInformationOptions, "getShardsInformationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getShardsInformationOptions.Db,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_shards`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getShardsInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetShardsInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalShardsInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetDocumentShardsInfo : Retrieve shard information for a specific document
// Retrieves information about a specific shard where a particular document is stored, along with information about the
// nodes where that shard has a replica.
func (cloudant *CloudantV1) GetDocumentShardsInfo(getDocumentShardsInfoOptions *GetDocumentShardsInfoOptions) (result *DocumentShardInfo, response *core.DetailedResponse, err error) {
	return cloudant.GetDocumentShardsInfoWithContext(context.Background(), getDocumentShardsInfoOptions)
}

// GetDocumentShardsInfoWithContext is an alternate form of the GetDocumentShardsInfo method which supports a Context parameter
func (cloudant *CloudantV1) GetDocumentShardsInfoWithContext(ctx context.Context, getDocumentShardsInfoOptions *GetDocumentShardsInfoOptions) (result *DocumentShardInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDocumentShardsInfoOptions, "getDocumentShardsInfoOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDocumentShardsInfoOptions, "getDocumentShardsInfoOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"db": *getDocumentShardsInfoOptions.Db,
		"doc_id": *getDocumentShardsInfoOptions.DocID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/{db}/_shards/{doc_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDocumentShardsInfoOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetDocumentShardsInfo")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDocumentShardInfo)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// HeadUpInformation : Retrieve HTTP headers about whether the server is up
// Retrieves the HTTP headers about whether the server is up.
func (cloudant *CloudantV1) HeadUpInformation(headUpInformationOptions *HeadUpInformationOptions) (response *core.DetailedResponse, err error) {
	return cloudant.HeadUpInformationWithContext(context.Background(), headUpInformationOptions)
}

// HeadUpInformationWithContext is an alternate form of the HeadUpInformation method which supports a Context parameter
func (cloudant *CloudantV1) HeadUpInformationWithContext(ctx context.Context, headUpInformationOptions *HeadUpInformationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(headUpInformationOptions, "headUpInformationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.HEAD)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_up`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range headUpInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "HeadUpInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cloudant.Service.Request(request, nil)

	return
}

// GetActiveTasks : Retrieve list of running tasks
// Lists running tasks, including the task type, name, status, and process ID. The result includes a JSON array of the
// currently running tasks, with each task described as a single object. Depending on the operation type, the set of
// response object fields might be different.
func (cloudant *CloudantV1) GetActiveTasks(getActiveTasksOptions *GetActiveTasksOptions) (result []ActiveTask, response *core.DetailedResponse, err error) {
	return cloudant.GetActiveTasksWithContext(context.Background(), getActiveTasksOptions)
}

// GetActiveTasksWithContext is an alternate form of the GetActiveTasks method which supports a Context parameter
func (cloudant *CloudantV1) GetActiveTasksWithContext(ctx context.Context, getActiveTasksOptions *GetActiveTasksOptions) (result []ActiveTask, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getActiveTasksOptions, "getActiveTasksOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_active_tasks`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getActiveTasksOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetActiveTasks")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalActiveTask)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetUpInformation : Retrieve information about whether the server is up
// Confirms that the server is up, running, and ready to respond to requests. If `maintenance_mode` is `true` or `nolb`,
// the endpoint returns a 404 response.
func (cloudant *CloudantV1) GetUpInformation(getUpInformationOptions *GetUpInformationOptions) (result *UpInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetUpInformationWithContext(context.Background(), getUpInformationOptions)
}

// GetUpInformationWithContext is an alternate form of the GetUpInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetUpInformationWithContext(ctx context.Context, getUpInformationOptions *GetUpInformationOptions) (result *UpInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getUpInformationOptions, "getUpInformationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_up`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getUpInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetUpInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUpInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetActivityTrackerEvents : Retrieve Activity Tracker events information
// Check event types that are being sent to IBM Cloud Activity Tracker for the IBM Cloudant instance.
func (cloudant *CloudantV1) GetActivityTrackerEvents(getActivityTrackerEventsOptions *GetActivityTrackerEventsOptions) (result *ActivityTrackerEvents, response *core.DetailedResponse, err error) {
	return cloudant.GetActivityTrackerEventsWithContext(context.Background(), getActivityTrackerEventsOptions)
}

// GetActivityTrackerEventsWithContext is an alternate form of the GetActivityTrackerEvents method which supports a Context parameter
func (cloudant *CloudantV1) GetActivityTrackerEventsWithContext(ctx context.Context, getActivityTrackerEventsOptions *GetActivityTrackerEventsOptions) (result *ActivityTrackerEvents, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getActivityTrackerEventsOptions, "getActivityTrackerEventsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/user/activity_tracker/events`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getActivityTrackerEventsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetActivityTrackerEvents")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalActivityTrackerEvents)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostActivityTrackerEvents : Modify Activity Tracker events configuration
// Configure event types that are being sent to IBM Cloud Activity Tracker for the IBM Cloudant instance.
func (cloudant *CloudantV1) PostActivityTrackerEvents(postActivityTrackerEventsOptions *PostActivityTrackerEventsOptions) (result *Ok, response *core.DetailedResponse, err error) {
	return cloudant.PostActivityTrackerEventsWithContext(context.Background(), postActivityTrackerEventsOptions)
}

// PostActivityTrackerEventsWithContext is an alternate form of the PostActivityTrackerEvents method which supports a Context parameter
func (cloudant *CloudantV1) PostActivityTrackerEventsWithContext(ctx context.Context, postActivityTrackerEventsOptions *PostActivityTrackerEventsOptions) (result *Ok, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postActivityTrackerEventsOptions, "postActivityTrackerEventsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postActivityTrackerEventsOptions, "postActivityTrackerEventsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/user/activity_tracker/events`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range postActivityTrackerEventsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "PostActivityTrackerEvents")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postActivityTrackerEventsOptions.Types != nil {
		body["types"] = postActivityTrackerEventsOptions.Types
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
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCurrentThroughputInformation : Retrieve the current provisioned throughput capacity consumption
// View the current consumption of provisioned throughput capacity for an IBM Cloudant instance. The current consumption
// shows the quantities of reads, writes, and global queries conducted against the instance for a given second.
func (cloudant *CloudantV1) GetCurrentThroughputInformation(getCurrentThroughputInformationOptions *GetCurrentThroughputInformationOptions) (result *CurrentThroughputInformation, response *core.DetailedResponse, err error) {
	return cloudant.GetCurrentThroughputInformationWithContext(context.Background(), getCurrentThroughputInformationOptions)
}

// GetCurrentThroughputInformationWithContext is an alternate form of the GetCurrentThroughputInformation method which supports a Context parameter
func (cloudant *CloudantV1) GetCurrentThroughputInformationWithContext(ctx context.Context, getCurrentThroughputInformationOptions *GetCurrentThroughputInformationOptions) (result *CurrentThroughputInformation, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getCurrentThroughputInformationOptions, "getCurrentThroughputInformationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cloudant.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cloudant.Service.Options.URL, `/_api/v2/user/current/throughput`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCurrentThroughputInformationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cloudant", "V1", "GetCurrentThroughputInformation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cloudant.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCurrentThroughputInformation)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ActiveTask : Schema for information about a running task.
type ActiveTask struct {
	// Processed changes.
	ChangesDone *int64 `json:"changes_done,omitempty"`

	// Source database.
	Database *string `json:"database" validate:"required"`

	// Cluster node where the task is running.
	Node *string `json:"node" validate:"required"`

	// Process ID.
	Pid *string `json:"pid" validate:"required"`

	// Current percentage progress.
	Progress *int64 `json:"progress,omitempty"`

	// Schema for a Unix epoch timestamp.
	StartedOn *int64 `json:"started_on" validate:"required"`

	// Task status message.
	Status *string `json:"status,omitempty"`

	// Task name.
	Task *string `json:"task,omitempty"`

	// Total changes to process.
	TotalChanges *int64 `json:"total_changes,omitempty"`

	// Operation type.
	Type *string `json:"type" validate:"required"`

	// Schema for a Unix epoch timestamp.
	UpdatedOn *int64 `json:"updated_on" validate:"required"`
}

// UnmarshalActiveTask unmarshals an instance of ActiveTask from the specified map of raw messages.
func UnmarshalActiveTask(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActiveTask)
	err = core.UnmarshalPrimitive(m, "changes_done", &obj.ChangesDone)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "database", &obj.Database)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "node", &obj.Node)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pid", &obj.Pid)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "progress", &obj.Progress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "started_on", &obj.StartedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task", &obj.Task)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_changes", &obj.TotalChanges)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActivityTrackerEvents : Schema for Activity Tracker events.
type ActivityTrackerEvents struct {
	// An array of event types that are being sent to IBM Cloud Activity Tracker for the IBM Cloudant instance.
	// "management" is a required element of this array.
	Types []string `json:"types" validate:"required"`
}

// Constants associated with the ActivityTrackerEvents.Types property.
const (
	ActivityTrackerEventsTypesDataConst = "data"
	ActivityTrackerEventsTypesManagementConst = "management"
)

// NewActivityTrackerEvents : Instantiate ActivityTrackerEvents (Generic Model Constructor)
func (*CloudantV1) NewActivityTrackerEvents(types []string) (_model *ActivityTrackerEvents, err error) {
	_model = &ActivityTrackerEvents{
		Types: types,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalActivityTrackerEvents unmarshals an instance of ActivityTrackerEvents from the specified map of raw messages.
func UnmarshalActivityTrackerEvents(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActivityTrackerEvents)
	err = core.UnmarshalPrimitive(m, "types", &obj.Types)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AllDocsQueriesResult : Schema for the result of an all documents queries operation.
type AllDocsQueriesResult struct {
	// An array of result objects - one for each query. Each result object contains the same fields as the response to a
	// regular `/_all_docs` request.
	Results []AllDocsResult `json:"results" validate:"required"`
}

// UnmarshalAllDocsQueriesResult unmarshals an instance of AllDocsQueriesResult from the specified map of raw messages.
func UnmarshalAllDocsQueriesResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllDocsQueriesResult)
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalAllDocsResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AllDocsQuery : Schema for an all documents query operation.
type AllDocsQuery struct {
	// Parameter to specify whether to include the encoding information in attachment stubs if the particular attachment is
	// compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the returned
	// document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Parameter to specify whether the specified end key should be included in the result.
	InclusiveEnd *bool `json:"inclusive_end,omitempty"`

	// Parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Parameter to specify whether to include in the response an update_seq value indicating the sequence id of the
	// database the view reflects.
	UpdateSeq *bool `json:"update_seq,omitempty"`

	// Schema for a document ID.
	Endkey *string `json:"endkey,omitempty"`

	// Schema for a document ID.
	Key *string `json:"key,omitempty"`

	// Schema for a list of document IDs.
	Keys []string `json:"keys,omitempty"`

	// Schema for a document ID.
	Startkey *string `json:"startkey,omitempty"`
}

// UnmarshalAllDocsQuery unmarshals an instance of AllDocsQuery from the specified map of raw messages.
func UnmarshalAllDocsQuery(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllDocsQuery)
	err = core.UnmarshalPrimitive(m, "att_encoding_info", &obj.AttEncodingInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachments", &obj.Attachments)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "conflicts", &obj.Conflicts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "descending", &obj.Descending)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "include_docs", &obj.IncludeDocs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inclusive_end", &obj.InclusiveEnd)
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
	err = core.UnmarshalPrimitive(m, "update_seq", &obj.UpdateSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "endkey", &obj.Endkey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "keys", &obj.Keys)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "startkey", &obj.Startkey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AllDocsResult : Schema for the result of an all documents operation.
type AllDocsResult struct {
	// Number of total rows.
	TotalRows *int64 `json:"total_rows" validate:"required"`

	// List of doc results.
	Rows []DocsResultRow `json:"rows" validate:"required"`

	// Current update sequence for the database.
	UpdateSeq *string `json:"update_seq,omitempty"`
}

// UnmarshalAllDocsResult unmarshals an instance of AllDocsResult from the specified map of raw messages.
func UnmarshalAllDocsResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllDocsResult)
	err = core.UnmarshalPrimitive(m, "total_rows", &obj.TotalRows)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rows", &obj.Rows, UnmarshalDocsResultRow)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update_seq", &obj.UpdateSeq)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Analyzer : Schema for a full text search analyzer.
type Analyzer struct {
	// Schema for the name of the Apache Lucene analyzer to use for text indexing. The default value varies depending on
	// the analyzer usage:
	// * For search indexes the default is `standard` * For query text indexes the default is `keyword` * For a query text
	// index default_field the default is `standard`.
	Name *string `json:"name,omitempty"`

	// Custom stopwords to use with the named analyzer.
	Stopwords []string `json:"stopwords,omitempty"`
}

// Constants associated with the Analyzer.Name property.
// Schema for the name of the Apache Lucene analyzer to use for text indexing. The default value varies depending on the
// analyzer usage:
// * For search indexes the default is `standard` * For query text indexes the default is `keyword` * For a query text
// index default_field the default is `standard`.
const (
	AnalyzerNameArabicConst = "arabic"
	AnalyzerNameArmenianConst = "armenian"
	AnalyzerNameBasqueConst = "basque"
	AnalyzerNameBrazilianConst = "brazilian"
	AnalyzerNameBulgarianConst = "bulgarian"
	AnalyzerNameCatalanConst = "catalan"
	AnalyzerNameChineseConst = "chinese"
	AnalyzerNameCjkConst = "cjk"
	AnalyzerNameClassicConst = "classic"
	AnalyzerNameCzechConst = "czech"
	AnalyzerNameDanishConst = "danish"
	AnalyzerNameDutchConst = "dutch"
	AnalyzerNameEmailConst = "email"
	AnalyzerNameEnglishConst = "english"
	AnalyzerNameFinnishConst = "finnish"
	AnalyzerNameFrenchConst = "french"
	AnalyzerNameGalicianConst = "galician"
	AnalyzerNameGermanConst = "german"
	AnalyzerNameGreekConst = "greek"
	AnalyzerNameHindiConst = "hindi"
	AnalyzerNameHungarianConst = "hungarian"
	AnalyzerNameIndonesianConst = "indonesian"
	AnalyzerNameIrishConst = "irish"
	AnalyzerNameItalianConst = "italian"
	AnalyzerNameJapaneseConst = "japanese"
	AnalyzerNameKeywordConst = "keyword"
	AnalyzerNameLatvianConst = "latvian"
	AnalyzerNameNorwegianConst = "norwegian"
	AnalyzerNamePerfieldConst = "perfield"
	AnalyzerNamePersianConst = "persian"
	AnalyzerNamePolishConst = "polish"
	AnalyzerNamePortugueseConst = "portuguese"
	AnalyzerNameRomanianConst = "romanian"
	AnalyzerNameRussianConst = "russian"
	AnalyzerNameSimpleConst = "simple"
	AnalyzerNameSpanishConst = "spanish"
	AnalyzerNameStandardConst = "standard"
	AnalyzerNameSwedishConst = "swedish"
	AnalyzerNameThaiConst = "thai"
	AnalyzerNameTurkishConst = "turkish"
	AnalyzerNameWhitespaceConst = "whitespace"
)

// UnmarshalAnalyzer unmarshals an instance of Analyzer from the specified map of raw messages.
func UnmarshalAnalyzer(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Analyzer)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "stopwords", &obj.Stopwords)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AnalyzerConfiguration : Schema for a search analyzer configuration.
type AnalyzerConfiguration struct {
	// Schema for the name of the Apache Lucene analyzer to use for text indexing. The default value varies depending on
	// the analyzer usage:
	// * For search indexes the default is `standard` * For query text indexes the default is `keyword` * For a query text
	// index default_field the default is `standard`.
	Name *string `json:"name,omitempty"`

	// Custom stopwords to use with the named analyzer.
	Stopwords []string `json:"stopwords,omitempty"`

	// Schema for mapping a field name to a per field analyzer.
	Fields map[string]Analyzer `json:"fields,omitempty"`
}

// Constants associated with the AnalyzerConfiguration.Name property.
// Schema for the name of the Apache Lucene analyzer to use for text indexing. The default value varies depending on the
// analyzer usage:
// * For search indexes the default is `standard` * For query text indexes the default is `keyword` * For a query text
// index default_field the default is `standard`.
const (
	AnalyzerConfigurationNameArabicConst = "arabic"
	AnalyzerConfigurationNameArmenianConst = "armenian"
	AnalyzerConfigurationNameBasqueConst = "basque"
	AnalyzerConfigurationNameBrazilianConst = "brazilian"
	AnalyzerConfigurationNameBulgarianConst = "bulgarian"
	AnalyzerConfigurationNameCatalanConst = "catalan"
	AnalyzerConfigurationNameChineseConst = "chinese"
	AnalyzerConfigurationNameCjkConst = "cjk"
	AnalyzerConfigurationNameClassicConst = "classic"
	AnalyzerConfigurationNameCzechConst = "czech"
	AnalyzerConfigurationNameDanishConst = "danish"
	AnalyzerConfigurationNameDutchConst = "dutch"
	AnalyzerConfigurationNameEmailConst = "email"
	AnalyzerConfigurationNameEnglishConst = "english"
	AnalyzerConfigurationNameFinnishConst = "finnish"
	AnalyzerConfigurationNameFrenchConst = "french"
	AnalyzerConfigurationNameGalicianConst = "galician"
	AnalyzerConfigurationNameGermanConst = "german"
	AnalyzerConfigurationNameGreekConst = "greek"
	AnalyzerConfigurationNameHindiConst = "hindi"
	AnalyzerConfigurationNameHungarianConst = "hungarian"
	AnalyzerConfigurationNameIndonesianConst = "indonesian"
	AnalyzerConfigurationNameIrishConst = "irish"
	AnalyzerConfigurationNameItalianConst = "italian"
	AnalyzerConfigurationNameJapaneseConst = "japanese"
	AnalyzerConfigurationNameKeywordConst = "keyword"
	AnalyzerConfigurationNameLatvianConst = "latvian"
	AnalyzerConfigurationNameNorwegianConst = "norwegian"
	AnalyzerConfigurationNamePerfieldConst = "perfield"
	AnalyzerConfigurationNamePersianConst = "persian"
	AnalyzerConfigurationNamePolishConst = "polish"
	AnalyzerConfigurationNamePortugueseConst = "portuguese"
	AnalyzerConfigurationNameRomanianConst = "romanian"
	AnalyzerConfigurationNameRussianConst = "russian"
	AnalyzerConfigurationNameSimpleConst = "simple"
	AnalyzerConfigurationNameSpanishConst = "spanish"
	AnalyzerConfigurationNameStandardConst = "standard"
	AnalyzerConfigurationNameSwedishConst = "swedish"
	AnalyzerConfigurationNameThaiConst = "thai"
	AnalyzerConfigurationNameTurkishConst = "turkish"
	AnalyzerConfigurationNameWhitespaceConst = "whitespace"
)

// UnmarshalAnalyzerConfiguration unmarshals an instance of AnalyzerConfiguration from the specified map of raw messages.
func UnmarshalAnalyzerConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AnalyzerConfiguration)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "stopwords", &obj.Stopwords)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "fields", &obj.Fields, UnmarshalAnalyzer)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApiKeysResult : Schema for api keys.
type ApiKeysResult struct {
	// ok.
	Ok *bool `json:"ok" validate:"required"`

	// The generated api key.
	Key *string `json:"key" validate:"required"`

	// The password associated with the api key.
	Password *string `json:"password" validate:"required"`
}

// UnmarshalApiKeysResult unmarshals an instance of ApiKeysResult from the specified map of raw messages.
func UnmarshalApiKeysResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApiKeysResult)
	err = core.UnmarshalPrimitive(m, "ok", &obj.Ok)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
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

// Attachment : Schema for an attachment.
type Attachment struct {
	// Attachment MIME type.
	ContentType *string `json:"content_type,omitempty"`

	// Base64-encoded content. Available if attachment content is requested by using the query parameters
	// `attachments=true` or `atts_since`. Note that when used with a view or changes feed `include_docs` must also be
	// `true`.
	Data *[]byte `json:"data,omitempty"`

	// Content hash digest. It starts with prefix which announce hash type (e.g. `md5-`) and continues with Base64-encoded
	// hash digest.
	Digest *string `json:"digest,omitempty"`

	// Compressed attachment size in bytes. Available if content_type was in list of compressible types when the attachment
	// was added and the query parameter `att_encoding_info` is `true`. Note that when used with a view or changes feed
	// `include_docs` must also be `true`.
	EncodedLength *int64 `json:"encoded_length,omitempty"`

	// Compression codec. Available if content_type was in list of compressible types when the attachment was added and the
	// and the query parameter `att_encoding_info` is `true`. Note that when used with a view or changes feed
	// `include_docs` must also be `true`.
	Encoding *string `json:"encoding,omitempty"`

	// True if the attachment follows in a multipart request or response.
	Follows *bool `json:"follows,omitempty"`

	// Real attachment size in bytes. Not available if inline attachment content requested.
	Length *int64 `json:"length,omitempty"`

	// Revision number when attachment was added.
	Revpos *int64 `json:"revpos,omitempty"`

	// Has `true` value if object contains stub info and no content. Otherwise omitted in response.
	Stub *bool `json:"stub,omitempty"`
}

// UnmarshalAttachment unmarshals an instance of Attachment from the specified map of raw messages.
func UnmarshalAttachment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Attachment)
	err = core.UnmarshalPrimitive(m, "content_type", &obj.ContentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "digest", &obj.Digest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "encoded_length", &obj.EncodedLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "encoding", &obj.Encoding)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "follows", &obj.Follows)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "length", &obj.Length)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revpos", &obj.Revpos)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "stub", &obj.Stub)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BulkDocs : Schema for submitting documents for bulk modifications.
type BulkDocs struct {
	// Array of documents.
	Docs []Document `json:"docs" validate:"required"`

	// If `false`, prevents the database from assigning them new revision IDs. Default is `true`.
	NewEdits *bool `json:"new_edits,omitempty"`
}

// NewBulkDocs : Instantiate BulkDocs (Generic Model Constructor)
func (*CloudantV1) NewBulkDocs(docs []Document) (_model *BulkDocs, err error) {
	_model = &BulkDocs{
		Docs: docs,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalBulkDocs unmarshals an instance of BulkDocs from the specified map of raw messages.
func UnmarshalBulkDocs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BulkDocs)
	err = core.UnmarshalModel(m, "docs", &obj.Docs, UnmarshalDocument)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "new_edits", &obj.NewEdits)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BulkGetQueryDocument : Schema for a document item in a bulk get query.
type BulkGetQueryDocument struct {
	// Includes attachments only since specified revisions.
	AttsSince []string `json:"atts_since,omitempty"`

	// Schema for a document ID.
	ID *string `json:"id" validate:"required"`

	// Schema for a document revision identifier.
	Rev *string `json:"rev,omitempty"`
}

// NewBulkGetQueryDocument : Instantiate BulkGetQueryDocument (Generic Model Constructor)
func (*CloudantV1) NewBulkGetQueryDocument(id string) (_model *BulkGetQueryDocument, err error) {
	_model = &BulkGetQueryDocument{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalBulkGetQueryDocument unmarshals an instance of BulkGetQueryDocument from the specified map of raw messages.
func UnmarshalBulkGetQueryDocument(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BulkGetQueryDocument)
	err = core.UnmarshalPrimitive(m, "atts_since", &obj.AttsSince)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rev", &obj.Rev)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BulkGetResult : Schema for the results object of a bulk get operation.
type BulkGetResult struct {
	// Results.
	Results []BulkGetResultItem `json:"results" validate:"required"`
}

// UnmarshalBulkGetResult unmarshals an instance of BulkGetResult from the specified map of raw messages.
func UnmarshalBulkGetResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BulkGetResult)
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalBulkGetResultItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BulkGetResultDocument : Schema for BulkGetResult object containing a successfully retrieved document or error information.
type BulkGetResultDocument struct {
	// Schema for the result of a document modification.
	Error *DocumentResult `json:"error,omitempty"`

	// Schema for a document.
	Ok *Document `json:"ok,omitempty"`
}

// UnmarshalBulkGetResultDocument unmarshals an instance of BulkGetResultDocument from the specified map of raw messages.
func UnmarshalBulkGetResultDocument(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BulkGetResultDocument)
	err = core.UnmarshalModel(m, "error", &obj.Error, UnmarshalDocumentResult)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "ok", &obj.Ok, UnmarshalDocument)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BulkGetResultItem : Schema for the document revisions information from a bulk get operation.
type BulkGetResultItem struct {
	// Array of document revisions or error information.
	Docs []BulkGetResultDocument `json:"docs" validate:"required"`

	// Schema for a document ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalBulkGetResultItem unmarshals an instance of BulkGetResultItem from the specified map of raw messages.
func UnmarshalBulkGetResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BulkGetResultItem)
	err = core.UnmarshalModel(m, "docs", &obj.Docs, UnmarshalBulkGetResultDocument)
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

// CapacityThroughputInformation : Schema for information about the currently provisioned and target throughput capacity.
type CapacityThroughputInformation struct {
	// Detailed information about provisioned throughput capacity.
	Current *CapacityThroughputInformationCurrent `json:"current" validate:"required"`

	// Detailed information about target throughput capacity.
	Target *CapacityThroughputInformationTarget `json:"target,omitempty"`
}

// UnmarshalCapacityThroughputInformation unmarshals an instance of CapacityThroughputInformation from the specified map of raw messages.
func UnmarshalCapacityThroughputInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CapacityThroughputInformation)
	err = core.UnmarshalModel(m, "current", &obj.Current, UnmarshalCapacityThroughputInformationCurrent)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalCapacityThroughputInformationTarget)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CapacityThroughputInformationCurrent : Detailed information about provisioned throughput capacity.
type CapacityThroughputInformationCurrent struct {
	// Schema for detailed information about throughput capacity with breakdown by specific throughput requests classes.
	Throughput *ThroughputInformation `json:"throughput" validate:"required"`
}

// UnmarshalCapacityThroughputInformationCurrent unmarshals an instance of CapacityThroughputInformationCurrent from the specified map of raw messages.
func UnmarshalCapacityThroughputInformationCurrent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CapacityThroughputInformationCurrent)
	err = core.UnmarshalModel(m, "throughput", &obj.Throughput, UnmarshalThroughputInformation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CapacityThroughputInformationTarget : Detailed information about target throughput capacity.
type CapacityThroughputInformationTarget struct {
	// Schema for detailed information about throughput capacity with breakdown by specific throughput requests classes.
	Throughput *ThroughputInformation `json:"throughput" validate:"required"`
}

// UnmarshalCapacityThroughputInformationTarget unmarshals an instance of CapacityThroughputInformationTarget from the specified map of raw messages.
func UnmarshalCapacityThroughputInformationTarget(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CapacityThroughputInformationTarget)
	err = core.UnmarshalModel(m, "throughput", &obj.Throughput, UnmarshalThroughputInformation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Change : Schema for a document leaf with single field rev.
type Change struct {
	// Schema for a document revision identifier.
	Rev *string `json:"rev" validate:"required"`
}

// UnmarshalChange unmarshals an instance of Change from the specified map of raw messages.
func UnmarshalChange(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Change)
	err = core.UnmarshalPrimitive(m, "rev", &obj.Rev)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ChangesResult : Schema for normal changes feed result.
type ChangesResult struct {
	// last_seq.
	LastSeq *string `json:"last_seq" validate:"required"`

	// pending.
	Pending *int64 `json:"pending" validate:"required"`

	// results.
	Results []ChangesResultItem `json:"results" validate:"required"`
}

// UnmarshalChangesResult unmarshals an instance of ChangesResult from the specified map of raw messages.
func UnmarshalChangesResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ChangesResult)
	err = core.UnmarshalPrimitive(m, "last_seq", &obj.LastSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pending", &obj.Pending)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalChangesResultItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ChangesResultItem : Schema for an item in the changes results array.
type ChangesResultItem struct {
	// List of document's leaves with single field rev.
	Changes []Change `json:"changes" validate:"required"`

	// if `true` then the document is deleted.
	Deleted *bool `json:"deleted,omitempty"`

	// Schema for a document.
	Doc *Document `json:"doc,omitempty"`

	// Schema for a document ID.
	ID *string `json:"id" validate:"required"`

	// Update sequence.
	Seq *string `json:"seq" validate:"required"`
}

// UnmarshalChangesResultItem unmarshals an instance of ChangesResultItem from the specified map of raw messages.
func UnmarshalChangesResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ChangesResultItem)
	err = core.UnmarshalModel(m, "changes", &obj.Changes, UnmarshalChange)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "deleted", &obj.Deleted)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "doc", &obj.Doc, UnmarshalDocument)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "seq", &obj.Seq)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ContentInformationSizes : Schema for size information of content.
type ContentInformationSizes struct {
	// The active size of the content, in bytes.
	Active *int64 `json:"active" validate:"required"`

	// The total uncompressed size of the content, in bytes.
	External *int64 `json:"external" validate:"required"`

	// The total size of the content as stored on disk, in bytes.
	File *int64 `json:"file" validate:"required"`
}

// UnmarshalContentInformationSizes unmarshals an instance of ContentInformationSizes from the specified map of raw messages.
func UnmarshalContentInformationSizes(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ContentInformationSizes)
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "external", &obj.External)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file", &obj.File)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CorsInformation : Schema for information about the CORS configuration.
type CorsInformation struct {
	// Boolean value to allow authentication credentials. If set to true, browser requests must be done by using
	// withCredentials = true.
	AllowCredentials *bool `json:"allow_credentials" validate:"required"`

	// Boolean value to turn CORS on and off.
	EnableCors *bool `json:"enable_cors" validate:"required"`

	// An array of strings that contain allowed origin domains. You have to specify the full URL including the protocol. It
	// is recommended that only the HTTPS protocol is used. Subdomains count as separate domains, so you have to specify
	// all subdomains used.
	Origins []string `json:"origins" validate:"required"`
}

// UnmarshalCorsInformation unmarshals an instance of CorsInformation from the specified map of raw messages.
func UnmarshalCorsInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CorsInformation)
	err = core.UnmarshalPrimitive(m, "allow_credentials", &obj.AllowCredentials)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_cors", &obj.EnableCors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "origins", &obj.Origins)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CurrentThroughputInformation : Schema for information about current consumption of a provisioned throughput capacity.
type CurrentThroughputInformation struct {
	// Detailed information about current consumption.
	Throughput *CurrentThroughputInformationThroughput `json:"throughput" validate:"required"`
}

// UnmarshalCurrentThroughputInformation unmarshals an instance of CurrentThroughputInformation from the specified map of raw messages.
func UnmarshalCurrentThroughputInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CurrentThroughputInformation)
	err = core.UnmarshalModel(m, "throughput", &obj.Throughput, UnmarshalCurrentThroughputInformationThroughput)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CurrentThroughputInformationThroughput : Detailed information about current consumption.
type CurrentThroughputInformationThroughput struct {
	// Number of global queries conducted against the instance for a given second.
	Query *int64 `json:"query" validate:"required"`

	// Number of reads conducted against the instance for a given second.
	Read *int64 `json:"read" validate:"required"`

	// Number of writes conducted against the instance for a given second.
	Write *int64 `json:"write" validate:"required"`
}

// UnmarshalCurrentThroughputInformationThroughput unmarshals an instance of CurrentThroughputInformationThroughput from the specified map of raw messages.
func UnmarshalCurrentThroughputInformationThroughput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CurrentThroughputInformationThroughput)
	err = core.UnmarshalPrimitive(m, "query", &obj.Query)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "read", &obj.Read)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "write", &obj.Write)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DatabaseInformation : Schema for information about a database.
type DatabaseInformation struct {
	// Schema for database cluster information.
	Cluster *DatabaseInformationCluster `json:"cluster" validate:"required"`

	// An opaque string that describes the committed state of the database.
	CommittedUpdateSeq *string `json:"committed_update_seq,omitempty"`

	// True if the database compaction routine is operating on this database.
	CompactRunning *bool `json:"compact_running" validate:"required"`

	// An opaque string that describes the compaction state of the database.
	CompactedSeq *string `json:"compacted_seq,omitempty"`

	// The name of the database.
	DbName *string `json:"db_name" validate:"required"`

	// The version of the physical format used for the data when it is stored on disk.
	DiskFormatVersion *int64 `json:"disk_format_version" validate:"required"`

	// A count of the documents in the specified database.
	DocCount *int64 `json:"doc_count" validate:"required"`

	// Number of deleted documents.
	DocDelCount *int64 `json:"doc_del_count" validate:"required"`

	// The engine used for the database.
	Engine *string `json:"engine,omitempty"`

	// Schema for database properties.
	Props *DatabaseInformationProps `json:"props" validate:"required"`

	// Schema for size information of content.
	Sizes *ContentInformationSizes `json:"sizes" validate:"required"`

	// An opaque string that describes the state of the database. Do not rely on this string for counting the number of
	// updates.
	UpdateSeq *string `json:"update_seq" validate:"required"`

	// The UUID of the database.
	UUID *string `json:"uuid,omitempty"`
}

// UnmarshalDatabaseInformation unmarshals an instance of DatabaseInformation from the specified map of raw messages.
func UnmarshalDatabaseInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DatabaseInformation)
	err = core.UnmarshalModel(m, "cluster", &obj.Cluster, UnmarshalDatabaseInformationCluster)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "committed_update_seq", &obj.CommittedUpdateSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compact_running", &obj.CompactRunning)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compacted_seq", &obj.CompactedSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "db_name", &obj.DbName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disk_format_version", &obj.DiskFormatVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_count", &obj.DocCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_del_count", &obj.DocDelCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "engine", &obj.Engine)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "props", &obj.Props, UnmarshalDatabaseInformationProps)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sizes", &obj.Sizes, UnmarshalContentInformationSizes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update_seq", &obj.UpdateSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uuid", &obj.UUID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DatabaseInformationCluster : Schema for database cluster information.
type DatabaseInformationCluster struct {
	// Schema for the number of replicas of a database in a cluster.
	N *int64 `json:"n" validate:"required"`

	// Schema for the number of shards in a database. Each shard is a partition of the hash value range.
	Q *int64 `json:"q" validate:"required"`

	// Read quorum. The number of consistent copies of a document that need to be read before a successful reply.
	R *int64 `json:"r" validate:"required"`

	// Write quorum. The number of copies of a document that need to be written before a successful reply.
	W *int64 `json:"w" validate:"required"`
}

// UnmarshalDatabaseInformationCluster unmarshals an instance of DatabaseInformationCluster from the specified map of raw messages.
func UnmarshalDatabaseInformationCluster(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DatabaseInformationCluster)
	err = core.UnmarshalPrimitive(m, "n", &obj.N)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "q", &obj.Q)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "r", &obj.R)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "w", &obj.W)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DatabaseInformationProps : Schema for database properties.
type DatabaseInformationProps struct {
	// The value is `true` for a partitioned database.
	Partitioned *bool `json:"partitioned,omitempty"`
}

// UnmarshalDatabaseInformationProps unmarshals an instance of DatabaseInformationProps from the specified map of raw messages.
func UnmarshalDatabaseInformationProps(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DatabaseInformationProps)
	err = core.UnmarshalPrimitive(m, "partitioned", &obj.Partitioned)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DbEvent : Schema for a database change event.
type DbEvent struct {
	// Account name.
	Account *string `json:"account,omitempty"`

	// Database name.
	DbName *string `json:"db_name" validate:"required"`

	// Sequence number.
	Seq *string `json:"seq" validate:"required"`

	// A database event.
	Type *string `json:"type" validate:"required"`
}

// Constants associated with the DbEvent.Type property.
// A database event.
const (
	DbEventTypeCreatedConst = "created"
	DbEventTypeDeletedConst = "deleted"
	DbEventTypeUpdatedConst = "updated"
)

// UnmarshalDbEvent unmarshals an instance of DbEvent from the specified map of raw messages.
func UnmarshalDbEvent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DbEvent)
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "db_name", &obj.DbName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "seq", &obj.Seq)
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

// DbUpdates : Schema for database updates.
type DbUpdates struct {
	// Last sequence number.
	LastSeq *string `json:"last_seq" validate:"required"`

	// results.
	Results []DbEvent `json:"results" validate:"required"`
}

// UnmarshalDbUpdates unmarshals an instance of DbUpdates from the specified map of raw messages.
func UnmarshalDbUpdates(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DbUpdates)
	err = core.UnmarshalPrimitive(m, "last_seq", &obj.LastSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalDbEvent)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DbsInfoResult : Schema for database information keyed by database name.
type DbsInfoResult struct {
	// The name of the error.
	Error *string `json:"error,omitempty"`

	// Schema for information about a database.
	Info *DatabaseInformation `json:"info,omitempty"`

	// Database name.
	Key *string `json:"key" validate:"required"`
}

// UnmarshalDbsInfoResult unmarshals an instance of DbsInfoResult from the specified map of raw messages.
func UnmarshalDbsInfoResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DbsInfoResult)
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "info", &obj.Info, UnmarshalDatabaseInformation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteAttachmentOptions : The DeleteAttachment options.
type DeleteAttachmentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Path parameter to specify the attachment name.
	AttachmentName *string `json:"attachment_name" validate:"required,ne="`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteAttachmentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	DeleteAttachmentOptionsBatchOkConst = "ok"
)

// NewDeleteAttachmentOptions : Instantiate DeleteAttachmentOptions
func (*CloudantV1) NewDeleteAttachmentOptions(db string, docID string, attachmentName string) *DeleteAttachmentOptions {
	return &DeleteAttachmentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
		AttachmentName: core.StringPtr(attachmentName),
	}
}

// SetDb : Allow user to set Db
func (_options *DeleteAttachmentOptions) SetDb(db string) *DeleteAttachmentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *DeleteAttachmentOptions) SetDocID(docID string) *DeleteAttachmentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetAttachmentName : Allow user to set AttachmentName
func (_options *DeleteAttachmentOptions) SetAttachmentName(attachmentName string) *DeleteAttachmentOptions {
	_options.AttachmentName = core.StringPtr(attachmentName)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeleteAttachmentOptions) SetIfMatch(ifMatch string) *DeleteAttachmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *DeleteAttachmentOptions) SetRev(rev string) *DeleteAttachmentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *DeleteAttachmentOptions) SetBatch(batch string) *DeleteAttachmentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAttachmentOptions) SetHeaders(param map[string]string) *DeleteAttachmentOptions {
	options.Headers = param
	return options
}

// DeleteDatabaseOptions : The DeleteDatabase options.
type DeleteDatabaseOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteDatabaseOptions : Instantiate DeleteDatabaseOptions
func (*CloudantV1) NewDeleteDatabaseOptions(db string) *DeleteDatabaseOptions {
	return &DeleteDatabaseOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *DeleteDatabaseOptions) SetDb(db string) *DeleteDatabaseOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDatabaseOptions) SetHeaders(param map[string]string) *DeleteDatabaseOptions {
	options.Headers = param
	return options
}

// DeleteDesignDocumentOptions : The DeleteDesignDocument options.
type DeleteDesignDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteDesignDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	DeleteDesignDocumentOptionsBatchOkConst = "ok"
)

// NewDeleteDesignDocumentOptions : Instantiate DeleteDesignDocumentOptions
func (*CloudantV1) NewDeleteDesignDocumentOptions(db string, ddoc string) *DeleteDesignDocumentOptions {
	return &DeleteDesignDocumentOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
	}
}

// SetDb : Allow user to set Db
func (_options *DeleteDesignDocumentOptions) SetDb(db string) *DeleteDesignDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *DeleteDesignDocumentOptions) SetDdoc(ddoc string) *DeleteDesignDocumentOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeleteDesignDocumentOptions) SetIfMatch(ifMatch string) *DeleteDesignDocumentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *DeleteDesignDocumentOptions) SetBatch(batch string) *DeleteDesignDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *DeleteDesignDocumentOptions) SetRev(rev string) *DeleteDesignDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDesignDocumentOptions) SetHeaders(param map[string]string) *DeleteDesignDocumentOptions {
	options.Headers = param
	return options
}

// DeleteDocumentOptions : The DeleteDocument options.
type DeleteDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	DeleteDocumentOptionsBatchOkConst = "ok"
)

// NewDeleteDocumentOptions : Instantiate DeleteDocumentOptions
func (*CloudantV1) NewDeleteDocumentOptions(db string, docID string) *DeleteDocumentOptions {
	return &DeleteDocumentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *DeleteDocumentOptions) SetDb(db string) *DeleteDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *DeleteDocumentOptions) SetDocID(docID string) *DeleteDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeleteDocumentOptions) SetIfMatch(ifMatch string) *DeleteDocumentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *DeleteDocumentOptions) SetBatch(batch string) *DeleteDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *DeleteDocumentOptions) SetRev(rev string) *DeleteDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDocumentOptions) SetHeaders(param map[string]string) *DeleteDocumentOptions {
	options.Headers = param
	return options
}

// DeleteIndexOptions : The DeleteIndex options.
type DeleteIndexOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the index type.
	Type *string `json:"type" validate:"required,ne="`

	// Path parameter to specify the index name.
	Index *string `json:"index" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteIndexOptions.Type property.
// Path parameter to specify the index type.
const (
	DeleteIndexOptionsTypeJSONConst = "json"
	DeleteIndexOptionsTypeSpecialConst = "special"
	DeleteIndexOptionsTypeTextConst = "text"
)

// NewDeleteIndexOptions : Instantiate DeleteIndexOptions
func (*CloudantV1) NewDeleteIndexOptions(db string, ddoc string, typeVar string, index string) *DeleteIndexOptions {
	return &DeleteIndexOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
		Type: core.StringPtr(typeVar),
		Index: core.StringPtr(index),
	}
}

// SetDb : Allow user to set Db
func (_options *DeleteIndexOptions) SetDb(db string) *DeleteIndexOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *DeleteIndexOptions) SetDdoc(ddoc string) *DeleteIndexOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetType : Allow user to set Type
func (_options *DeleteIndexOptions) SetType(typeVar string) *DeleteIndexOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetIndex : Allow user to set Index
func (_options *DeleteIndexOptions) SetIndex(index string) *DeleteIndexOptions {
	_options.Index = core.StringPtr(index)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteIndexOptions) SetHeaders(param map[string]string) *DeleteIndexOptions {
	options.Headers = param
	return options
}

// DeleteLocalDocumentOptions : The DeleteLocalDocument options.
type DeleteLocalDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteLocalDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	DeleteLocalDocumentOptionsBatchOkConst = "ok"
)

// NewDeleteLocalDocumentOptions : Instantiate DeleteLocalDocumentOptions
func (*CloudantV1) NewDeleteLocalDocumentOptions(db string, docID string) *DeleteLocalDocumentOptions {
	return &DeleteLocalDocumentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *DeleteLocalDocumentOptions) SetDb(db string) *DeleteLocalDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *DeleteLocalDocumentOptions) SetDocID(docID string) *DeleteLocalDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *DeleteLocalDocumentOptions) SetBatch(batch string) *DeleteLocalDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLocalDocumentOptions) SetHeaders(param map[string]string) *DeleteLocalDocumentOptions {
	options.Headers = param
	return options
}

// DeleteReplicationDocumentOptions : The DeleteReplicationDocument options.
type DeleteReplicationDocumentOptions struct {
	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteReplicationDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	DeleteReplicationDocumentOptionsBatchOkConst = "ok"
)

// NewDeleteReplicationDocumentOptions : Instantiate DeleteReplicationDocumentOptions
func (*CloudantV1) NewDeleteReplicationDocumentOptions(docID string) *DeleteReplicationDocumentOptions {
	return &DeleteReplicationDocumentOptions{
		DocID: core.StringPtr(docID),
	}
}

// SetDocID : Allow user to set DocID
func (_options *DeleteReplicationDocumentOptions) SetDocID(docID string) *DeleteReplicationDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *DeleteReplicationDocumentOptions) SetIfMatch(ifMatch string) *DeleteReplicationDocumentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *DeleteReplicationDocumentOptions) SetBatch(batch string) *DeleteReplicationDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *DeleteReplicationDocumentOptions) SetRev(rev string) *DeleteReplicationDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteReplicationDocumentOptions) SetHeaders(param map[string]string) *DeleteReplicationDocumentOptions {
	options.Headers = param
	return options
}

// DesignDocument : Schema for a design document.
type DesignDocument struct {
	// Schema for a map of attachment name to attachment metadata.
	Attachments map[string]Attachment `json:"_attachments,omitempty"`

	// Schema for a list of document revision identifiers.
	Conflicts []string `json:"_conflicts,omitempty"`

	// Deletion flag. Available if document was removed.
	Deleted *bool `json:"_deleted,omitempty"`

	// Schema for a list of document revision identifiers.
	DeletedConflicts []string `json:"_deleted_conflicts,omitempty"`

	// Document ID.
	ID *string `json:"_id,omitempty"`

	// Document's update sequence in current database. Available if requested with local_seq=true query parameter.
	LocalSeq *string `json:"_local_seq,omitempty"`

	// Schema for a document revision identifier.
	Rev *string `json:"_rev,omitempty"`

	// Schema for list of revision information.
	Revisions *Revisions `json:"_revisions,omitempty"`

	// Schema for a list of objects with information about local revisions and their status.
	RevsInfo []DocumentRevisionStatus `json:"_revs_info,omitempty"`

	// Indicates whether to automatically build indexes defined in this design document.
	Autoupdate *bool `json:"autoupdate,omitempty"`

	// Schema for filter functions definition. This schema is a map where keys are the names of the filter functions and
	// values are the function definition in string format.
	//
	// Filter function formats, or filters the changes feed that pass filter rules. The function takes 2 parameters:
	//
	//   * `doc`: The document that is being processed.
	//   * `req`: A Request JavaScript object with these properties:
	//
	//     * `body` - string, Request body data as string.
	//       If the request method is GET this field contains the value
	//       `"undefined"`.
	//       If the method is DELETE or HEAD the value is `""` (empty string).
	//     * `cookie` - Cookies object.
	//     * `form` - Form Data object, contains the decoded body as key-value
	//       pairs if the Content-Type header was
	//       application/x-www-form-urlencoded.
	//     * `headers` - Request Headers object.
	//     * `id` - string, requested document id if it was specified
	//       or null otherwise.
	//     * `info` - Database Information object,
	//       see `DatabaseInformation`.
	//     * `method` - string or an array of chars, request method.
	//       String value is a method as one of: HEAD, GET, POST, PUT,
	//       DELETE, OPTIONS, TRACE, COPY. For not supported methods
	//       it will be represented as an array of char codes e.g. for VIEW
	//       it will be 86,73,69,87.
	//     * `path` - array of strings, requested path sections.
	//     * `peer` - string, request source IP address.
	//     * `query` - string, URL query parameters object. Note that multiple
	//       keys are not supported and the last key value suppresses others.
	//     * `requested_path` - array of strings,
	//       actual requested path section.
	//     * `raw_path` - string, raw requested path.
	//     * `userCtx`: User Context Object, containing information about the
	//       user writing the document (if present), see the `UserContext`.
	//     * `secObj`: Security Object, with lists of database security roles,
	//       see the `SecurityObject`.
	//     * `uuid` - string, generated UUID by a specified algorithm in the
	//       config file.
	//
	// Filter functions must return true if a document passed all the rules.
	Filters map[string]string `json:"filters,omitempty"`

	// Search (text) index function definitions.
	Indexes map[string]SearchIndexDefinition `json:"indexes,omitempty"`

	// Defines Query Server key to process design document functions.
	Language *string `json:"language,omitempty"`

	// Schema for design document options.
	Options *DesignDocumentOptions `json:"options,omitempty"`

	// Validate document update function can be used to prevent invalid or unauthorized document update requests from being
	// stored. Validation functions typically examine the structure of the new document to ensure that required fields are
	// present and to verify that the requesting user should be allowed to make changes to the document properties. When a
	// write request is received for a given database, the validation function in each design document in that database is
	// called in an unspecified order. If any of the validation functions throw an error, the write will not succeed.
	//
	// The validation function can abort the pending document write by throwing one of two error objects:
	//
	// ```
	// // user is not authorized to make the change but may re-authenticate throw({ unauthorized: 'Error message here.' });
	//
	// // change is not allowed throw({ forbidden: 'Error message here.' });
	// ```
	//
	// The function takes 4 parameters:
	//
	//   * `newDoc` - New version of document that will be stored
	//     from the update request.
	//   * `oldDoc` - Previous version of document that is already stored.
	//   * `userCtx` - User Context Object, containing information about the
	//     user writing the document (if present), see the `UserContext`.
	//   * `secObj` - Security Object, with lists of database security roles,
	//     see the `SecurityObject`.
	ValidateDocUpdate *string `json:"validate_doc_update,omitempty"`

	// Schema for design document views.
	Views map[string]DesignDocumentViewsMapReduce `json:"views,omitempty"`

	// Schema for geospatial index function definitions.
	StIndexes map[string]GeoIndexDefinition `json:"st_indexes,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of DesignDocument
func (o *DesignDocument) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of DesignDocument
func (o *DesignDocument) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of DesignDocument
func (o *DesignDocument) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of DesignDocument
func (o *DesignDocument) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of DesignDocument
func (o *DesignDocument) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Attachments != nil {
		m["_attachments"] = o.Attachments
	}
	if o.Conflicts != nil {
		m["_conflicts"] = o.Conflicts
	}
	if o.Deleted != nil {
		m["_deleted"] = o.Deleted
	}
	if o.DeletedConflicts != nil {
		m["_deleted_conflicts"] = o.DeletedConflicts
	}
	if o.ID != nil {
		m["_id"] = o.ID
	}
	if o.LocalSeq != nil {
		m["_local_seq"] = o.LocalSeq
	}
	if o.Rev != nil {
		m["_rev"] = o.Rev
	}
	if o.Revisions != nil {
		m["_revisions"] = o.Revisions
	}
	if o.RevsInfo != nil {
		m["_revs_info"] = o.RevsInfo
	}
	if o.Autoupdate != nil {
		m["autoupdate"] = o.Autoupdate
	}
	if o.Filters != nil {
		m["filters"] = o.Filters
	}
	if o.Indexes != nil {
		m["indexes"] = o.Indexes
	}
	if o.Language != nil {
		m["language"] = o.Language
	}
	if o.Options != nil {
		m["options"] = o.Options
	}
	if o.ValidateDocUpdate != nil {
		m["validate_doc_update"] = o.ValidateDocUpdate
	}
	if o.Views != nil {
		m["views"] = o.Views
	}
	if o.StIndexes != nil {
		m["st_indexes"] = o.StIndexes
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalDesignDocument unmarshals an instance of DesignDocument from the specified map of raw messages.
func UnmarshalDesignDocument(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DesignDocument)
	err = core.UnmarshalModel(m, "_attachments", &obj.Attachments, UnmarshalAttachment)
	if err != nil {
		return
	}
	delete(m, "_attachments")
	err = core.UnmarshalPrimitive(m, "_conflicts", &obj.Conflicts)
	if err != nil {
		return
	}
	delete(m, "_conflicts")
	err = core.UnmarshalPrimitive(m, "_deleted", &obj.Deleted)
	if err != nil {
		return
	}
	delete(m, "_deleted")
	err = core.UnmarshalPrimitive(m, "_deleted_conflicts", &obj.DeletedConflicts)
	if err != nil {
		return
	}
	delete(m, "_deleted_conflicts")
	err = core.UnmarshalPrimitive(m, "_id", &obj.ID)
	if err != nil {
		return
	}
	delete(m, "_id")
	err = core.UnmarshalPrimitive(m, "_local_seq", &obj.LocalSeq)
	if err != nil {
		return
	}
	delete(m, "_local_seq")
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	delete(m, "_rev")
	err = core.UnmarshalModel(m, "_revisions", &obj.Revisions, UnmarshalRevisions)
	if err != nil {
		return
	}
	delete(m, "_revisions")
	err = core.UnmarshalModel(m, "_revs_info", &obj.RevsInfo, UnmarshalDocumentRevisionStatus)
	if err != nil {
		return
	}
	delete(m, "_revs_info")
	err = core.UnmarshalPrimitive(m, "autoupdate", &obj.Autoupdate)
	if err != nil {
		return
	}
	delete(m, "autoupdate")
	err = core.UnmarshalPrimitive(m, "filters", &obj.Filters)
	if err != nil {
		return
	}
	delete(m, "filters")
	err = core.UnmarshalModel(m, "indexes", &obj.Indexes, UnmarshalSearchIndexDefinition)
	if err != nil {
		return
	}
	delete(m, "indexes")
	err = core.UnmarshalPrimitive(m, "language", &obj.Language)
	if err != nil {
		return
	}
	delete(m, "language")
	err = core.UnmarshalModel(m, "options", &obj.Options, UnmarshalDesignDocumentOptions)
	if err != nil {
		return
	}
	delete(m, "options")
	err = core.UnmarshalPrimitive(m, "validate_doc_update", &obj.ValidateDocUpdate)
	if err != nil {
		return
	}
	delete(m, "validate_doc_update")
	err = core.UnmarshalModel(m, "views", &obj.Views, UnmarshalDesignDocumentViewsMapReduce)
	if err != nil {
		return
	}
	delete(m, "views")
	err = core.UnmarshalModel(m, "st_indexes", &obj.StIndexes, UnmarshalGeoIndexDefinition)
	if err != nil {
		return
	}
	delete(m, "st_indexes")
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

// DesignDocumentInformation : Schema for information about a design document.
type DesignDocumentInformation struct {
	// name.
	Name *string `json:"name" validate:"required"`

	// View index information.
	ViewIndex *DesignDocumentViewIndex `json:"view_index" validate:"required"`
}

// UnmarshalDesignDocumentInformation unmarshals an instance of DesignDocumentInformation from the specified map of raw messages.
func UnmarshalDesignDocumentInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DesignDocumentInformation)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "view_index", &obj.ViewIndex, UnmarshalDesignDocumentViewIndex)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DesignDocumentOptions : Schema for design document options.
type DesignDocumentOptions struct {
	// Whether this design document describes partitioned or global indexes.
	Partitioned *bool `json:"partitioned,omitempty"`
}

// UnmarshalDesignDocumentOptions unmarshals an instance of DesignDocumentOptions from the specified map of raw messages.
func UnmarshalDesignDocumentOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DesignDocumentOptions)
	err = core.UnmarshalPrimitive(m, "partitioned", &obj.Partitioned)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DesignDocumentViewIndex : View index information.
type DesignDocumentViewIndex struct {
	// Indicates whether a compaction routine is currently running on the view.
	CompactRunning *bool `json:"compact_running" validate:"required"`

	// Language for the defined views.
	Language *string `json:"language" validate:"required"`

	// MD5 signature of the views for the design document.
	Signature *string `json:"signature" validate:"required"`

	// Schema for size information of content.
	Sizes *ContentInformationSizes `json:"sizes" validate:"required"`

	// Indicates if the view is currently being updated.
	UpdaterRunning *bool `json:"updater_running" validate:"required"`

	// Number of clients waiting on views from this design document.
	WaitingClients *int64 `json:"waiting_clients" validate:"required"`

	// Indicates if there are outstanding commits to the underlying database that need to processed.
	WaitingCommit *bool `json:"waiting_commit" validate:"required"`
}

// UnmarshalDesignDocumentViewIndex unmarshals an instance of DesignDocumentViewIndex from the specified map of raw messages.
func UnmarshalDesignDocumentViewIndex(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DesignDocumentViewIndex)
	err = core.UnmarshalPrimitive(m, "compact_running", &obj.CompactRunning)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "language", &obj.Language)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signature", &obj.Signature)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sizes", &obj.Sizes, UnmarshalContentInformationSizes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updater_running", &obj.UpdaterRunning)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "waiting_clients", &obj.WaitingClients)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "waiting_commit", &obj.WaitingCommit)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DesignDocumentViewsMapReduce : Schema for view functions definition.
type DesignDocumentViewsMapReduce struct {
	// JavaScript map function as a string.
	Map *string `json:"map" validate:"required"`

	// JavaScript reduce function as a string.
	Reduce *string `json:"reduce,omitempty"`
}

// NewDesignDocumentViewsMapReduce : Instantiate DesignDocumentViewsMapReduce (Generic Model Constructor)
func (*CloudantV1) NewDesignDocumentViewsMapReduce(mapVar string) (_model *DesignDocumentViewsMapReduce, err error) {
	_model = &DesignDocumentViewsMapReduce{
		Map: core.StringPtr(mapVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalDesignDocumentViewsMapReduce unmarshals an instance of DesignDocumentViewsMapReduce from the specified map of raw messages.
func UnmarshalDesignDocumentViewsMapReduce(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DesignDocumentViewsMapReduce)
	err = core.UnmarshalPrimitive(m, "map", &obj.Map)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reduce", &obj.Reduce)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DocsResultRow : Schema for a row of document information in a DocsResult.
type DocsResultRow struct {
	// The cause of the error (if available).
	CausedBy *string `json:"caused_by,omitempty"`

	// The name of the error.
	Error *string `json:"error,omitempty"`

	// The reason the error occurred (if available).
	Reason *string `json:"reason,omitempty"`

	// Schema for a document.
	Doc *Document `json:"doc,omitempty"`

	// id.
	ID *string `json:"id,omitempty"`

	// Document ID.
	Key *string `json:"key" validate:"required"`

	// Value of built-in `/_all_docs` style view.
	Value *DocsResultRowValue `json:"value,omitempty"`
}

// UnmarshalDocsResultRow unmarshals an instance of DocsResultRow from the specified map of raw messages.
func UnmarshalDocsResultRow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DocsResultRow)
	err = core.UnmarshalPrimitive(m, "caused_by", &obj.CausedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "doc", &obj.Doc, UnmarshalDocument)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "value", &obj.Value, UnmarshalDocsResultRowValue)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DocsResultRowValue : Value of built-in `/_all_docs` style view.
type DocsResultRowValue struct {
	// Schema for a document revision identifier.
	Rev *string `json:"rev" validate:"required"`
}

// UnmarshalDocsResultRowValue unmarshals an instance of DocsResultRowValue from the specified map of raw messages.
func UnmarshalDocsResultRowValue(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DocsResultRowValue)
	err = core.UnmarshalPrimitive(m, "rev", &obj.Rev)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Document : Schema for a document.
type Document struct {
	// Schema for a map of attachment name to attachment metadata.
	Attachments map[string]Attachment `json:"_attachments,omitempty"`

	// Schema for a list of document revision identifiers.
	Conflicts []string `json:"_conflicts,omitempty"`

	// Deletion flag. Available if document was removed.
	Deleted *bool `json:"_deleted,omitempty"`

	// Schema for a list of document revision identifiers.
	DeletedConflicts []string `json:"_deleted_conflicts,omitempty"`

	// Document ID.
	ID *string `json:"_id,omitempty"`

	// Document's update sequence in current database. Available if requested with local_seq=true query parameter.
	LocalSeq *string `json:"_local_seq,omitempty"`

	// Schema for a document revision identifier.
	Rev *string `json:"_rev,omitempty"`

	// Schema for list of revision information.
	Revisions *Revisions `json:"_revisions,omitempty"`

	// Schema for a list of objects with information about local revisions and their status.
	RevsInfo []DocumentRevisionStatus `json:"_revs_info,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of Document
func (o *Document) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of Document
func (o *Document) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of Document
func (o *Document) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of Document
func (o *Document) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of Document
func (o *Document) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Attachments != nil {
		m["_attachments"] = o.Attachments
	}
	if o.Conflicts != nil {
		m["_conflicts"] = o.Conflicts
	}
	if o.Deleted != nil {
		m["_deleted"] = o.Deleted
	}
	if o.DeletedConflicts != nil {
		m["_deleted_conflicts"] = o.DeletedConflicts
	}
	if o.ID != nil {
		m["_id"] = o.ID
	}
	if o.LocalSeq != nil {
		m["_local_seq"] = o.LocalSeq
	}
	if o.Rev != nil {
		m["_rev"] = o.Rev
	}
	if o.Revisions != nil {
		m["_revisions"] = o.Revisions
	}
	if o.RevsInfo != nil {
		m["_revs_info"] = o.RevsInfo
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalDocument unmarshals an instance of Document from the specified map of raw messages.
func UnmarshalDocument(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Document)
	err = core.UnmarshalModel(m, "_attachments", &obj.Attachments, UnmarshalAttachment)
	if err != nil {
		return
	}
	delete(m, "_attachments")
	err = core.UnmarshalPrimitive(m, "_conflicts", &obj.Conflicts)
	if err != nil {
		return
	}
	delete(m, "_conflicts")
	err = core.UnmarshalPrimitive(m, "_deleted", &obj.Deleted)
	if err != nil {
		return
	}
	delete(m, "_deleted")
	err = core.UnmarshalPrimitive(m, "_deleted_conflicts", &obj.DeletedConflicts)
	if err != nil {
		return
	}
	delete(m, "_deleted_conflicts")
	err = core.UnmarshalPrimitive(m, "_id", &obj.ID)
	if err != nil {
		return
	}
	delete(m, "_id")
	err = core.UnmarshalPrimitive(m, "_local_seq", &obj.LocalSeq)
	if err != nil {
		return
	}
	delete(m, "_local_seq")
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	delete(m, "_rev")
	err = core.UnmarshalModel(m, "_revisions", &obj.Revisions, UnmarshalRevisions)
	if err != nil {
		return
	}
	delete(m, "_revisions")
	err = core.UnmarshalModel(m, "_revs_info", &obj.RevsInfo, UnmarshalDocumentRevisionStatus)
	if err != nil {
		return
	}
	delete(m, "_revs_info")
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

// DocumentResult : Schema for the result of a document modification.
type DocumentResult struct {
	// Schema for a document ID.
	ID *string `json:"id" validate:"required"`

	// Schema for a document revision identifier.
	Rev *string `json:"rev,omitempty"`

	// ok.
	Ok *bool `json:"ok,omitempty"`

	// The cause of the error (if available).
	CausedBy *string `json:"caused_by,omitempty"`

	// The name of the error.
	Error *string `json:"error,omitempty"`

	// The reason the error occurred (if available).
	Reason *string `json:"reason,omitempty"`
}

// UnmarshalDocumentResult unmarshals an instance of DocumentResult from the specified map of raw messages.
func UnmarshalDocumentResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DocumentResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rev", &obj.Rev)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ok", &obj.Ok)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "caused_by", &obj.CausedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DocumentRevisionStatus : Schema for information about revisions and their status.
type DocumentRevisionStatus struct {
	// Schema for a document revision identifier.
	Rev *string `json:"rev" validate:"required"`

	// Status of the revision. May be one of: - `available`: Revision is available for retrieving with rev query parameter
	// - `missing`: Revision is not available - `deleted`: Revision belongs to deleted document.
	Status *string `json:"status" validate:"required"`
}

// Constants associated with the DocumentRevisionStatus.Status property.
// Status of the revision. May be one of: - `available`: Revision is available for retrieving with rev query parameter -
// `missing`: Revision is not available - `deleted`: Revision belongs to deleted document.
const (
	DocumentRevisionStatusStatusAvailableConst = "available"
	DocumentRevisionStatusStatusDeletedConst = "deleted"
	DocumentRevisionStatusStatusMissingConst = "missing"
)

// NewDocumentRevisionStatus : Instantiate DocumentRevisionStatus (Generic Model Constructor)
func (*CloudantV1) NewDocumentRevisionStatus(rev string, status string) (_model *DocumentRevisionStatus, err error) {
	_model = &DocumentRevisionStatus{
		Rev: core.StringPtr(rev),
		Status: core.StringPtr(status),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalDocumentRevisionStatus unmarshals an instance of DocumentRevisionStatus from the specified map of raw messages.
func UnmarshalDocumentRevisionStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DocumentRevisionStatus)
	err = core.UnmarshalPrimitive(m, "rev", &obj.Rev)
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

// DocumentShardInfo : Schema for document shard information.
type DocumentShardInfo struct {
	// List of nodes serving a replica of the shard.
	Nodes []string `json:"nodes" validate:"required"`

	// The shard range in which the document is stored.
	Range *string `json:"range" validate:"required"`
}

// UnmarshalDocumentShardInfo unmarshals an instance of DocumentShardInfo from the specified map of raw messages.
func UnmarshalDocumentShardInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DocumentShardInfo)
	err = core.UnmarshalPrimitive(m, "nodes", &obj.Nodes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "range", &obj.Range)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExecutionStats : Schema for find query execution statistics.
type ExecutionStats struct {
	// Time to execute the query.
	ExecutionTimeMs *float64 `json:"execution_time_ms" validate:"required"`

	// Number of results returned.
	ResultsReturned *int64 `json:"results_returned" validate:"required"`

	// Number of documents fetched from the index.
	TotalDocsExamined *int64 `json:"total_docs_examined" validate:"required"`

	// Number of rows scanned in the index.
	TotalKeysExamined *int64 `json:"total_keys_examined" validate:"required"`

	// Number of documents fetched from the primary index with the specified read quorum.
	TotalQuorumDocsExamined *int64 `json:"total_quorum_docs_examined" validate:"required"`
}

// UnmarshalExecutionStats unmarshals an instance of ExecutionStats from the specified map of raw messages.
func UnmarshalExecutionStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExecutionStats)
	err = core.UnmarshalPrimitive(m, "execution_time_ms", &obj.ExecutionTimeMs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "results_returned", &obj.ResultsReturned)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_docs_examined", &obj.TotalDocsExamined)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_keys_examined", &obj.TotalKeysExamined)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_quorum_docs_examined", &obj.TotalQuorumDocsExamined)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExplainResult : Schema for information about the index used for a find query.
type ExplainResult struct {
	// dbname.
	Dbname *string `json:"dbname" validate:"required"`

	// fields.
	Fields []string `json:"fields" validate:"required"`

	// Schema for information about an index.
	Index *IndexInformation `json:"index" validate:"required"`

	// limit.
	Limit *int64 `json:"limit" validate:"required"`

	// opts.
	Opts map[string]interface{} `json:"opts" validate:"required"`

	// range.
	Range *ExplainResultRange `json:"range,omitempty"`

	// JSON object describing criteria used to select documents. The selector specifies fields in the document, and
	// provides an expression to evaluate with the field content or other data.
	//
	// The selector object must:
	//   * Be structured as valid JSON.
	//   * Contain a valid query expression.
	//
	// Using a selector is significantly more efficient than using a JavaScript filter function, and is the recommended
	// option if filtering on document attributes only.
	//
	// Elementary selector syntax requires you to specify one or more fields, and the corresponding values required for
	// those fields. You can create more complex selector expressions by combining operators.
	//
	// Operators are identified by the use of a dollar sign `$` prefix in the name field.
	//
	// There are two core types of operators in the selector syntax:
	// * Combination operators: applied at the topmost level of selection. They are used to combine selectors. In addition
	// to the common boolean operators (`$and`, `$or`, `$not`, `$nor`) there are three combination operators: `$all`,
	// `$elemMatch`, and `$allMatch`. A combination operator takes a single argument. The argument is either another
	// selector, or an array of selectors.
	// * Condition operators: are specific to a field, and are used to evaluate the value stored in that field. For
	// instance, the basic `$eq` operator matches when the specified field contains a value that is equal to the supplied
	// argument.
	Selector map[string]interface{} `json:"selector" validate:"required"`

	// skip.
	Skip *int64 `json:"skip" validate:"required"`
}

// UnmarshalExplainResult unmarshals an instance of ExplainResult from the specified map of raw messages.
func UnmarshalExplainResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExplainResult)
	err = core.UnmarshalPrimitive(m, "dbname", &obj.Dbname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "fields", &obj.Fields)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "index", &obj.Index, UnmarshalIndexInformation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "opts", &obj.Opts)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "range", &obj.Range, UnmarshalExplainResultRange)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "selector", &obj.Selector)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "skip", &obj.Skip)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExplainResultRange : range.
type ExplainResultRange struct {
	// end_key.
	EndKey []interface{} `json:"end_key,omitempty"`

	// start_key.
	StartKey []interface{} `json:"start_key,omitempty"`
}

// UnmarshalExplainResultRange unmarshals an instance of ExplainResultRange from the specified map of raw messages.
func UnmarshalExplainResultRange(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExplainResultRange)
	err = core.UnmarshalPrimitive(m, "end_key", &obj.EndKey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_key", &obj.StartKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FindResult : Schema for the result of a query find operation.
type FindResult struct {
	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark" validate:"required"`

	// Documents matching the selector.
	Docs []Document `json:"docs" validate:"required"`

	// Schema for find query execution statistics.
	ExecutionStats *ExecutionStats `json:"execution_stats,omitempty"`

	// warning.
	Warning *string `json:"warning,omitempty"`
}

// UnmarshalFindResult unmarshals an instance of FindResult from the specified map of raw messages.
func UnmarshalFindResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FindResult)
	err = core.UnmarshalPrimitive(m, "bookmark", &obj.Bookmark)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "docs", &obj.Docs, UnmarshalDocument)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "execution_stats", &obj.ExecutionStats, UnmarshalExecutionStats)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "warning", &obj.Warning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GeoIndexDefinition : Schema for a geospatial index definition.
type GeoIndexDefinition struct {
	// String form of a JavaScript function that is called for each document in the database. The function takes the
	// document as a parameter, extracts some geospatial data from it, and then calls the `st_index` function to index that
	// data. The `st_index` takes a GeoJSON geometry as a parameter.
	Index *string `json:"index" validate:"required"`
}

// NewGeoIndexDefinition : Instantiate GeoIndexDefinition (Generic Model Constructor)
func (*CloudantV1) NewGeoIndexDefinition(index string) (_model *GeoIndexDefinition, err error) {
	_model = &GeoIndexDefinition{
		Index: core.StringPtr(index),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalGeoIndexDefinition unmarshals an instance of GeoIndexDefinition from the specified map of raw messages.
func UnmarshalGeoIndexDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoIndexDefinition)
	err = core.UnmarshalPrimitive(m, "index", &obj.Index)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GeoIndexInformation : Schema for information about a geospatial index.
type GeoIndexInformation struct {
	// Schema for geospatial index statistics.
	GeoIndex *GeoIndexStats `json:"geo_index" validate:"required"`

	// The name of the geospatial index design document.
	Name *string `json:"name" validate:"required"`
}

// UnmarshalGeoIndexInformation unmarshals an instance of GeoIndexInformation from the specified map of raw messages.
func UnmarshalGeoIndexInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoIndexInformation)
	err = core.UnmarshalModel(m, "geo_index", &obj.GeoIndex, UnmarshalGeoIndexStats)
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

// GeoIndexStats : Schema for geospatial index statistics.
type GeoIndexStats struct {
	// The size of the geospatial index, in bytes.
	DataSize *int64 `json:"data_size" validate:"required"`

	// The size of the geospatial index, as stored on disk, in bytes.
	DiskSize *int64 `json:"disk_size" validate:"required"`

	// Number of documents in the geospatial index.
	DocCount *int64 `json:"doc_count" validate:"required"`
}

// UnmarshalGeoIndexStats unmarshals an instance of GeoIndexStats from the specified map of raw messages.
func UnmarshalGeoIndexStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoIndexStats)
	err = core.UnmarshalPrimitive(m, "data_size", &obj.DataSize)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disk_size", &obj.DiskSize)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_count", &obj.DocCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GeoJSONFeature : Schema for a GeoJSON feature object. Note that the content of the feature objects varies depending on the response
// format chosen and whether the `include_docs` parameter is `true`.
type GeoJSONFeature struct {
	// Schema for a document ID.
	ID *string `json:"_id,omitempty"`

	// Schema for a document revision identifier.
	Rev *string `json:"_rev,omitempty"`

	// Schema for a GeoJSON bounding box.
	Bbox []float64 `json:"bbox,omitempty"`

	// Schema for a GeoJSON geometry object.
	Geometry GeoJSONGeometryObjectIntf `json:"geometry" validate:"required"`

	// Schema for the properties of a GeoJSON feature object.
	Properties map[string]interface{} `json:"properties,omitempty"`

	// Declaration of the GeoJSON type: Feature Object.
	Type *string `json:"type" validate:"required"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// Constants associated with the GeoJSONFeature.Type property.
// Declaration of the GeoJSON type: Feature Object.
const (
	GeoJSONFeatureTypeFeatureConst = "Feature"
)

// SetProperty allows the user to set an arbitrary property on an instance of GeoJSONFeature
func (o *GeoJSONFeature) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of GeoJSONFeature
func (o *GeoJSONFeature) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of GeoJSONFeature
func (o *GeoJSONFeature) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of GeoJSONFeature
func (o *GeoJSONFeature) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of GeoJSONFeature
func (o *GeoJSONFeature) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.ID != nil {
		m["_id"] = o.ID
	}
	if o.Rev != nil {
		m["_rev"] = o.Rev
	}
	if o.Bbox != nil {
		m["bbox"] = o.Bbox
	}
	if o.Geometry != nil {
		m["geometry"] = o.Geometry
	}
	if o.Properties != nil {
		m["properties"] = o.Properties
	}
	if o.Type != nil {
		m["type"] = o.Type
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalGeoJSONFeature unmarshals an instance of GeoJSONFeature from the specified map of raw messages.
func UnmarshalGeoJSONFeature(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoJSONFeature)
	err = core.UnmarshalPrimitive(m, "_id", &obj.ID)
	if err != nil {
		return
	}
	delete(m, "_id")
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	delete(m, "_rev")
	err = core.UnmarshalPrimitive(m, "bbox", &obj.Bbox)
	if err != nil {
		return
	}
	delete(m, "bbox")
	err = core.UnmarshalModel(m, "geometry", &obj.Geometry, UnmarshalGeoJSONGeometryObject)
	if err != nil {
		return
	}
	delete(m, "geometry")
	err = core.UnmarshalPrimitive(m, "properties", &obj.Properties)
	if err != nil {
		return
	}
	delete(m, "properties")
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

// GeoJSONGeometryObject : Schema for a GeoJSON geometry object.
// Models which "extend" this model:
// - GeoJSONGeometry
// - GeoJSONGeometryCollection
type GeoJSONGeometryObject struct {
	// The type of GeoJSON Geometry.
	Type *string `json:"type,omitempty"`

	// Used for all geometry types except `GeometryCollection`. The structure of the elements in the array varies by
	// geometry type.
	Coordinates []interface{} `json:"coordinates,omitempty"`

	// Used for the `GeometryCollection` type.
	Geometries []GeoJSONGeometry `json:"geometries,omitempty"`
}

// Constants associated with the GeoJSONGeometryObject.Type property.
// The type of GeoJSON Geometry.
const (
	GeoJSONGeometryObjectTypeGeometrycollectionConst = "GeometryCollection"
	GeoJSONGeometryObjectTypeLinestringConst = "LineString"
	GeoJSONGeometryObjectTypeMultilinestringConst = "MultiLineString"
	GeoJSONGeometryObjectTypeMultipointConst = "MultiPoint"
	GeoJSONGeometryObjectTypeMultipolygonConst = "MultiPolygon"
	GeoJSONGeometryObjectTypePointConst = "Point"
	GeoJSONGeometryObjectTypePolygonConst = "Polygon"
)
func (*GeoJSONGeometryObject) isaGeoJSONGeometryObject() bool {
	return true
}

type GeoJSONGeometryObjectIntf interface {
	isaGeoJSONGeometryObject() bool
}

// UnmarshalGeoJSONGeometryObject unmarshals an instance of GeoJSONGeometryObject from the specified map of raw messages.
func UnmarshalGeoJSONGeometryObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoJSONGeometryObject)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "coordinates", &obj.Coordinates)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "geometries", &obj.Geometries, UnmarshalGeoJSONGeometry)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GeoResult : Schema for the result of a geospatial query operation. For the `legacy`, `geojson`, or `application/vnd.geo+json`
// format this is a GeoJson FeatureCollection with additional metadata in foreign members.
type GeoResult struct {
	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark" validate:"required"`

	// The array of GeoJSON Feature Objects matching the geospatial query.
	Features []GeoJSONFeature `json:"features,omitempty"`

	// The array of rows matching the geospatial query. Present only when using `view` format.
	Rows []GeoResultRow `json:"rows" validate:"required"`

	// Declaration of the GeoJSON type: FeatureCollection Object.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the GeoResult.Type property.
// Declaration of the GeoJSON type: FeatureCollection Object.
const (
	GeoResultTypeFeaturecollectionConst = "FeatureCollection"
)

// UnmarshalGeoResult unmarshals an instance of GeoResult from the specified map of raw messages.
func UnmarshalGeoResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoResult)
	err = core.UnmarshalPrimitive(m, "bookmark", &obj.Bookmark)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalGeoJSONFeature)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rows", &obj.Rows, UnmarshalGeoResultRow)
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

// GeoResultRow : Schema for a row of a geospatial result using view format.
type GeoResultRow struct {
	// Schema for a document.
	Doc *Document `json:"doc,omitempty"`

	// Schema for a GeoJSON geometry.
	Geometry *GeoJSONGeometry `json:"geometry,omitempty"`

	// Schema for a document ID.
	ID *string `json:"id,omitempty"`

	// Schema for a document revision identifier.
	Rev *string `json:"rev,omitempty"`
}

// UnmarshalGeoResultRow unmarshals an instance of GeoResultRow from the specified map of raw messages.
func UnmarshalGeoResultRow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoResultRow)
	err = core.UnmarshalModel(m, "doc", &obj.Doc, UnmarshalDocument)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "geometry", &obj.Geometry, UnmarshalGeoJSONGeometry)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rev", &obj.Rev)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetActiveTasksOptions : The GetActiveTasks options.
type GetActiveTasksOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetActiveTasksOptions : Instantiate GetActiveTasksOptions
func (*CloudantV1) NewGetActiveTasksOptions() *GetActiveTasksOptions {
	return &GetActiveTasksOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetActiveTasksOptions) SetHeaders(param map[string]string) *GetActiveTasksOptions {
	options.Headers = param
	return options
}

// GetActivityTrackerEventsOptions : The GetActivityTrackerEvents options.
type GetActivityTrackerEventsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetActivityTrackerEventsOptions : Instantiate GetActivityTrackerEventsOptions
func (*CloudantV1) NewGetActivityTrackerEventsOptions() *GetActivityTrackerEventsOptions {
	return &GetActivityTrackerEventsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetActivityTrackerEventsOptions) SetHeaders(param map[string]string) *GetActivityTrackerEventsOptions {
	options.Headers = param
	return options
}

// GetAllDbsOptions : The GetAllDbs options.
type GetAllDbsOptions struct {
	// Query parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Query parameter to specify to stop returning records when the specified key is reached. String representation of any
	// JSON type that matches the key type emitted by the view function.
	Endkey *string `json:"endkey,omitempty"`

	// Query parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Query parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Query parameter to specify to start returning records from the specified key. String representation of any JSON type
	// that matches the key type emitted by the view function.
	Startkey *string `json:"startkey,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAllDbsOptions : Instantiate GetAllDbsOptions
func (*CloudantV1) NewGetAllDbsOptions() *GetAllDbsOptions {
	return &GetAllDbsOptions{}
}

// SetDescending : Allow user to set Descending
func (_options *GetAllDbsOptions) SetDescending(descending bool) *GetAllDbsOptions {
	_options.Descending = core.BoolPtr(descending)
	return _options
}

// SetEndkey : Allow user to set Endkey
func (_options *GetAllDbsOptions) SetEndkey(endkey string) *GetAllDbsOptions {
	_options.Endkey = core.StringPtr(endkey)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetAllDbsOptions) SetLimit(limit int64) *GetAllDbsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *GetAllDbsOptions) SetSkip(skip int64) *GetAllDbsOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetStartkey : Allow user to set Startkey
func (_options *GetAllDbsOptions) SetStartkey(startkey string) *GetAllDbsOptions {
	_options.Startkey = core.StringPtr(startkey)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAllDbsOptions) SetHeaders(param map[string]string) *GetAllDbsOptions {
	options.Headers = param
	return options
}

// GetAttachmentOptions : The GetAttachment options.
type GetAttachmentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Path parameter to specify the attachment name.
	AttachmentName *string `json:"attachment_name" validate:"required,ne="`

	// The type of the response:  or *_/_*.
	Accept *string `json:"Accept,omitempty"`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Header parameter to specify the byte range for a request. This allows the implementation of resumable downloads and
	// skippable streams. This is available for all attachments inside CouchDB.
	Range *string `json:"Range,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAttachmentOptions : Instantiate GetAttachmentOptions
func (*CloudantV1) NewGetAttachmentOptions(db string, docID string, attachmentName string) *GetAttachmentOptions {
	return &GetAttachmentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
		AttachmentName: core.StringPtr(attachmentName),
	}
}

// SetDb : Allow user to set Db
func (_options *GetAttachmentOptions) SetDb(db string) *GetAttachmentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *GetAttachmentOptions) SetDocID(docID string) *GetAttachmentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetAttachmentName : Allow user to set AttachmentName
func (_options *GetAttachmentOptions) SetAttachmentName(attachmentName string) *GetAttachmentOptions {
	_options.AttachmentName = core.StringPtr(attachmentName)
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *GetAttachmentOptions) SetAccept(accept string) *GetAttachmentOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *GetAttachmentOptions) SetIfMatch(ifMatch string) *GetAttachmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *GetAttachmentOptions) SetIfNoneMatch(ifNoneMatch string) *GetAttachmentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetRange : Allow user to set Range
func (_options *GetAttachmentOptions) SetRange(rangeVar string) *GetAttachmentOptions {
	_options.Range = core.StringPtr(rangeVar)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *GetAttachmentOptions) SetRev(rev string) *GetAttachmentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAttachmentOptions) SetHeaders(param map[string]string) *GetAttachmentOptions {
	options.Headers = param
	return options
}

// GetCapacityThroughputInformationOptions : The GetCapacityThroughputInformation options.
type GetCapacityThroughputInformationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCapacityThroughputInformationOptions : Instantiate GetCapacityThroughputInformationOptions
func (*CloudantV1) NewGetCapacityThroughputInformationOptions() *GetCapacityThroughputInformationOptions {
	return &GetCapacityThroughputInformationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetCapacityThroughputInformationOptions) SetHeaders(param map[string]string) *GetCapacityThroughputInformationOptions {
	options.Headers = param
	return options
}

// GetCorsInformationOptions : The GetCorsInformation options.
type GetCorsInformationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCorsInformationOptions : Instantiate GetCorsInformationOptions
func (*CloudantV1) NewGetCorsInformationOptions() *GetCorsInformationOptions {
	return &GetCorsInformationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetCorsInformationOptions) SetHeaders(param map[string]string) *GetCorsInformationOptions {
	options.Headers = param
	return options
}

// GetCurrentThroughputInformationOptions : The GetCurrentThroughputInformation options.
type GetCurrentThroughputInformationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCurrentThroughputInformationOptions : Instantiate GetCurrentThroughputInformationOptions
func (*CloudantV1) NewGetCurrentThroughputInformationOptions() *GetCurrentThroughputInformationOptions {
	return &GetCurrentThroughputInformationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetCurrentThroughputInformationOptions) SetHeaders(param map[string]string) *GetCurrentThroughputInformationOptions {
	options.Headers = param
	return options
}

// GetDatabaseInformationOptions : The GetDatabaseInformation options.
type GetDatabaseInformationOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDatabaseInformationOptions : Instantiate GetDatabaseInformationOptions
func (*CloudantV1) NewGetDatabaseInformationOptions(db string) *GetDatabaseInformationOptions {
	return &GetDatabaseInformationOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *GetDatabaseInformationOptions) SetDb(db string) *GetDatabaseInformationOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDatabaseInformationOptions) SetHeaders(param map[string]string) *GetDatabaseInformationOptions {
	options.Headers = param
	return options
}

// GetDbUpdatesOptions : The GetDbUpdates options.
type GetDbUpdatesOptions struct {
	// Query parameter to specify the changes feed type.
	Feed *string `json:"feed,omitempty"`

	// Query parameter to specify the period in milliseconds after which an empty line is sent in the results. Only
	// applicable for longpoll, continuous, and eventsource feeds. Overrides any timeout to keep the feed alive
	// indefinitely. May also be `true` to use default value of 60000.
	Heartbeat *int64 `json:"heartbeat,omitempty"`

	// Query parameter to specify the maximum period in milliseconds to wait for a change before the response is sent, even
	// if there are no results. Only applicable for `longpoll` or `continuous` feeds. Default value is specified by
	// `httpd/changes_timeout` configuration option. Note that `60000` value is also the default maximum timeout to prevent
	// undetected dead connections.
	Timeout *int64 `json:"timeout,omitempty"`

	// Query parameter to specify to start the results from the change immediately after the given update sequence. Can be
	// a valid update sequence or `now` value. Default is `0` i.e. all changes.
	Since *string `json:"since,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetDbUpdatesOptions.Feed property.
// Query parameter to specify the changes feed type.
const (
	GetDbUpdatesOptionsFeedContinuousConst = "continuous"
	GetDbUpdatesOptionsFeedEventsourceConst = "eventsource"
	GetDbUpdatesOptionsFeedLongpollConst = "longpoll"
	GetDbUpdatesOptionsFeedNormalConst = "normal"
)

// NewGetDbUpdatesOptions : Instantiate GetDbUpdatesOptions
func (*CloudantV1) NewGetDbUpdatesOptions() *GetDbUpdatesOptions {
	return &GetDbUpdatesOptions{}
}

// SetFeed : Allow user to set Feed
func (_options *GetDbUpdatesOptions) SetFeed(feed string) *GetDbUpdatesOptions {
	_options.Feed = core.StringPtr(feed)
	return _options
}

// SetHeartbeat : Allow user to set Heartbeat
func (_options *GetDbUpdatesOptions) SetHeartbeat(heartbeat int64) *GetDbUpdatesOptions {
	_options.Heartbeat = core.Int64Ptr(heartbeat)
	return _options
}

// SetTimeout : Allow user to set Timeout
func (_options *GetDbUpdatesOptions) SetTimeout(timeout int64) *GetDbUpdatesOptions {
	_options.Timeout = core.Int64Ptr(timeout)
	return _options
}

// SetSince : Allow user to set Since
func (_options *GetDbUpdatesOptions) SetSince(since string) *GetDbUpdatesOptions {
	_options.Since = core.StringPtr(since)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDbUpdatesOptions) SetHeaders(param map[string]string) *GetDbUpdatesOptions {
	options.Headers = param
	return options
}

// GetDesignDocumentInformationOptions : The GetDesignDocumentInformation options.
type GetDesignDocumentInformationOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDesignDocumentInformationOptions : Instantiate GetDesignDocumentInformationOptions
func (*CloudantV1) NewGetDesignDocumentInformationOptions(db string, ddoc string) *GetDesignDocumentInformationOptions {
	return &GetDesignDocumentInformationOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
	}
}

// SetDb : Allow user to set Db
func (_options *GetDesignDocumentInformationOptions) SetDb(db string) *GetDesignDocumentInformationOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *GetDesignDocumentInformationOptions) SetDdoc(ddoc string) *GetDesignDocumentInformationOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDesignDocumentInformationOptions) SetHeaders(param map[string]string) *GetDesignDocumentInformationOptions {
	options.Headers = param
	return options
}

// GetDesignDocumentOptions : The GetDesignDocument options.
type GetDesignDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Query parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Query parameter to specify whether to include the encoding information in attachment stubs if the particular
	// attachment is compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Query parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the
	// returned document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Query parameter to specify whether to include a list of deleted conflicted revisions in the `_deleted_conflicts`
	// property of the returned document.
	DeletedConflicts *bool `json:"deleted_conflicts,omitempty"`

	// Query parameter to specify whether to force retrieving latest leaf revision, no matter what rev was requested.
	Latest *bool `json:"latest,omitempty"`

	// Query parameter to specify whether to include the last update sequence for the document.
	LocalSeq *bool `json:"local_seq,omitempty"`

	// Query parameter to specify whether to include document meta information. Acts the same as specifying all of the
	// conflicts, deleted_conflicts and open_revs query parameters.
	Meta *bool `json:"meta,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Query parameter to specify whether to include a list of all known document revisions.
	Revs *bool `json:"revs,omitempty"`

	// Query parameter to specify whether to includes detailed information for all known document revisions.
	RevsInfo *bool `json:"revs_info,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDesignDocumentOptions : Instantiate GetDesignDocumentOptions
func (*CloudantV1) NewGetDesignDocumentOptions(db string, ddoc string) *GetDesignDocumentOptions {
	return &GetDesignDocumentOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
	}
}

// SetDb : Allow user to set Db
func (_options *GetDesignDocumentOptions) SetDb(db string) *GetDesignDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *GetDesignDocumentOptions) SetDdoc(ddoc string) *GetDesignDocumentOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *GetDesignDocumentOptions) SetIfNoneMatch(ifNoneMatch string) *GetDesignDocumentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *GetDesignDocumentOptions) SetAttachments(attachments bool) *GetDesignDocumentOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *GetDesignDocumentOptions) SetAttEncodingInfo(attEncodingInfo bool) *GetDesignDocumentOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *GetDesignDocumentOptions) SetConflicts(conflicts bool) *GetDesignDocumentOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDeletedConflicts : Allow user to set DeletedConflicts
func (_options *GetDesignDocumentOptions) SetDeletedConflicts(deletedConflicts bool) *GetDesignDocumentOptions {
	_options.DeletedConflicts = core.BoolPtr(deletedConflicts)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *GetDesignDocumentOptions) SetLatest(latest bool) *GetDesignDocumentOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetLocalSeq : Allow user to set LocalSeq
func (_options *GetDesignDocumentOptions) SetLocalSeq(localSeq bool) *GetDesignDocumentOptions {
	_options.LocalSeq = core.BoolPtr(localSeq)
	return _options
}

// SetMeta : Allow user to set Meta
func (_options *GetDesignDocumentOptions) SetMeta(meta bool) *GetDesignDocumentOptions {
	_options.Meta = core.BoolPtr(meta)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *GetDesignDocumentOptions) SetRev(rev string) *GetDesignDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetRevs : Allow user to set Revs
func (_options *GetDesignDocumentOptions) SetRevs(revs bool) *GetDesignDocumentOptions {
	_options.Revs = core.BoolPtr(revs)
	return _options
}

// SetRevsInfo : Allow user to set RevsInfo
func (_options *GetDesignDocumentOptions) SetRevsInfo(revsInfo bool) *GetDesignDocumentOptions {
	_options.RevsInfo = core.BoolPtr(revsInfo)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDesignDocumentOptions) SetHeaders(param map[string]string) *GetDesignDocumentOptions {
	options.Headers = param
	return options
}

// GetDocumentOptions : The GetDocument options.
type GetDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Query parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Query parameter to specify whether to include the encoding information in attachment stubs if the particular
	// attachment is compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Query parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the
	// returned document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Query parameter to specify whether to include a list of deleted conflicted revisions in the `_deleted_conflicts`
	// property of the returned document.
	DeletedConflicts *bool `json:"deleted_conflicts,omitempty"`

	// Query parameter to specify whether to force retrieving latest leaf revision, no matter what rev was requested.
	Latest *bool `json:"latest,omitempty"`

	// Query parameter to specify whether to include the last update sequence for the document.
	LocalSeq *bool `json:"local_seq,omitempty"`

	// Query parameter to specify whether to include document meta information. Acts the same as specifying all of the
	// conflicts, deleted_conflicts and open_revs query parameters.
	Meta *bool `json:"meta,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Query parameter to specify whether to include a list of all known document revisions.
	Revs *bool `json:"revs,omitempty"`

	// Query parameter to specify whether to includes detailed information for all known document revisions.
	RevsInfo *bool `json:"revs_info,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDocumentOptions : Instantiate GetDocumentOptions
func (*CloudantV1) NewGetDocumentOptions(db string, docID string) *GetDocumentOptions {
	return &GetDocumentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *GetDocumentOptions) SetDb(db string) *GetDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *GetDocumentOptions) SetDocID(docID string) *GetDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *GetDocumentOptions) SetIfNoneMatch(ifNoneMatch string) *GetDocumentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *GetDocumentOptions) SetAttachments(attachments bool) *GetDocumentOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *GetDocumentOptions) SetAttEncodingInfo(attEncodingInfo bool) *GetDocumentOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *GetDocumentOptions) SetConflicts(conflicts bool) *GetDocumentOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDeletedConflicts : Allow user to set DeletedConflicts
func (_options *GetDocumentOptions) SetDeletedConflicts(deletedConflicts bool) *GetDocumentOptions {
	_options.DeletedConflicts = core.BoolPtr(deletedConflicts)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *GetDocumentOptions) SetLatest(latest bool) *GetDocumentOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetLocalSeq : Allow user to set LocalSeq
func (_options *GetDocumentOptions) SetLocalSeq(localSeq bool) *GetDocumentOptions {
	_options.LocalSeq = core.BoolPtr(localSeq)
	return _options
}

// SetMeta : Allow user to set Meta
func (_options *GetDocumentOptions) SetMeta(meta bool) *GetDocumentOptions {
	_options.Meta = core.BoolPtr(meta)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *GetDocumentOptions) SetRev(rev string) *GetDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetRevs : Allow user to set Revs
func (_options *GetDocumentOptions) SetRevs(revs bool) *GetDocumentOptions {
	_options.Revs = core.BoolPtr(revs)
	return _options
}

// SetRevsInfo : Allow user to set RevsInfo
func (_options *GetDocumentOptions) SetRevsInfo(revsInfo bool) *GetDocumentOptions {
	_options.RevsInfo = core.BoolPtr(revsInfo)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDocumentOptions) SetHeaders(param map[string]string) *GetDocumentOptions {
	options.Headers = param
	return options
}

// GetDocumentShardsInfoOptions : The GetDocumentShardsInfo options.
type GetDocumentShardsInfoOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDocumentShardsInfoOptions : Instantiate GetDocumentShardsInfoOptions
func (*CloudantV1) NewGetDocumentShardsInfoOptions(db string, docID string) *GetDocumentShardsInfoOptions {
	return &GetDocumentShardsInfoOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *GetDocumentShardsInfoOptions) SetDb(db string) *GetDocumentShardsInfoOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *GetDocumentShardsInfoOptions) SetDocID(docID string) *GetDocumentShardsInfoOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDocumentShardsInfoOptions) SetHeaders(param map[string]string) *GetDocumentShardsInfoOptions {
	options.Headers = param
	return options
}

// GetGeoIndexInformationOptions : The GetGeoIndexInformation options.
type GetGeoIndexInformationOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the index name.
	Index *string `json:"index" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetGeoIndexInformationOptions : Instantiate GetGeoIndexInformationOptions
func (*CloudantV1) NewGetGeoIndexInformationOptions(db string, ddoc string, index string) *GetGeoIndexInformationOptions {
	return &GetGeoIndexInformationOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
		Index: core.StringPtr(index),
	}
}

// SetDb : Allow user to set Db
func (_options *GetGeoIndexInformationOptions) SetDb(db string) *GetGeoIndexInformationOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *GetGeoIndexInformationOptions) SetDdoc(ddoc string) *GetGeoIndexInformationOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetIndex : Allow user to set Index
func (_options *GetGeoIndexInformationOptions) SetIndex(index string) *GetGeoIndexInformationOptions {
	_options.Index = core.StringPtr(index)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetGeoIndexInformationOptions) SetHeaders(param map[string]string) *GetGeoIndexInformationOptions {
	options.Headers = param
	return options
}

// GetGeoOptions : The GetGeo options.
type GetGeoOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the index name.
	Index *string `json:"index" validate:"required,ne="`

	// Query parameter to specify a geospatial query bounding box with two latitude,longitude coordinates for the
	// lower-left and upper-right corners. An example is `-11.05987446,12.28339928,-101.05987446,62.28339928`.
	Bbox *string `json:"bbox,omitempty"`

	// Query parameter to specify a bookmark that was received from a previous request. This parameter enables paging
	// through the results. If there are no more results after the bookmark, you get a response containing no further
	// results and the same bookmark, confirming the end of the result list.
	Bookmark *string `json:"bookmark,omitempty"`

	// Query parameter that causes the geospatial query output to be returned in the specified format.
	Format *string `json:"format,omitempty"`

	// Query parameter to specify a Well Known Text (WKT) representation of a geospatial query geometry. The valid values
	// for the WKT parameter include `Point`, `LineString`, `Polygon`, `MultiPoint`, `MultiLineString`, `MultiPolygon`, and
	// `GeometryCollection`.
	G *string `json:"g,omitempty"`

	// Query parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Query parameter to specify a latitude coordinate for use with radius or ellipse geospatial queries.
	Lat *float64 `json:"lat,omitempty"`

	// Query parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Query parameter to specify a longitude coordinate for use with radius or ellipse geospatial queries.
	Lon *float64 `json:"lon,omitempty"`

	// Query parameter to specify whether to perform a nearest neighbour (NN) search. If provided, the `nearest=true`
	// search returns all results by sorting their distances to the center of the query geometry. NN search can be used
	// alone or with any of the supported DE-9IM (Dimensionally Extended nine-Intersection Model) specification geometric
	// relations documented.
	Nearest *bool `json:"nearest,omitempty"`

	// Query parameter to specify the radius, in meters, to search from a lat,lon coordinate point in a circular geospatial
	// query.
	Radius *float64 `json:"radius,omitempty"`

	// Query parameter to specify the first radius, in meters, to search from a lat,lon coordinate point in an ellipse
	// geospatial query.
	Rangex *float64 `json:"rangex,omitempty"`

	// Query parameter to specify the second radius, in meters, to search from a lat,lon coordinate point in an ellipse
	// geospatial query.
	Rangey *float64 `json:"rangey,omitempty"`

	// Query parameter to specify the DE-9IM (Dimensionally Extended nine-Intersection Model)geospatial relationship
	// between the query geometry and the result documents.
	Relation *string `json:"relation,omitempty"`

	// Query parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Query parameter to specify to not wait for the index to finish building before returning results.
	Stale *string `json:"stale,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetGeoOptions.Format property.
// Query parameter that causes the geospatial query output to be returned in the specified format.
const (
	GetGeoOptionsFormatApplicationVndGeoJSONConst = "application/vnd.geo+json"
	GetGeoOptionsFormatGeojsonConst = "geojson"
	GetGeoOptionsFormatLegacyConst = "legacy"
	GetGeoOptionsFormatViewConst = "view"
)

// Constants associated with the GetGeoOptions.Relation property.
// Query parameter to specify the DE-9IM (Dimensionally Extended nine-Intersection Model)geospatial relationship between
// the query geometry and the result documents.
const (
	GetGeoOptionsRelationContainsConst = "contains"
	GetGeoOptionsRelationContainsProperlyConst = "contains_properly"
	GetGeoOptionsRelationCoveredByConst = "covered_by"
	GetGeoOptionsRelationCoversConst = "covers"
	GetGeoOptionsRelationCrossesConst = "crosses"
	GetGeoOptionsRelationDisjointConst = "disjoint"
	GetGeoOptionsRelationIntersectsConst = "intersects"
	GetGeoOptionsRelationOverlapsConst = "overlaps"
	GetGeoOptionsRelationTouchesConst = "touches"
	GetGeoOptionsRelationWithinConst = "within"
)

// Constants associated with the GetGeoOptions.Stale property.
// Query parameter to specify to not wait for the index to finish building before returning results.
const (
	GetGeoOptionsStaleOkConst = "ok"
)

// NewGetGeoOptions : Instantiate GetGeoOptions
func (*CloudantV1) NewGetGeoOptions(db string, ddoc string, index string) *GetGeoOptions {
	return &GetGeoOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
		Index: core.StringPtr(index),
	}
}

// SetDb : Allow user to set Db
func (_options *GetGeoOptions) SetDb(db string) *GetGeoOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *GetGeoOptions) SetDdoc(ddoc string) *GetGeoOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetIndex : Allow user to set Index
func (_options *GetGeoOptions) SetIndex(index string) *GetGeoOptions {
	_options.Index = core.StringPtr(index)
	return _options
}

// SetBbox : Allow user to set Bbox
func (_options *GetGeoOptions) SetBbox(bbox string) *GetGeoOptions {
	_options.Bbox = core.StringPtr(bbox)
	return _options
}

// SetBookmark : Allow user to set Bookmark
func (_options *GetGeoOptions) SetBookmark(bookmark string) *GetGeoOptions {
	_options.Bookmark = core.StringPtr(bookmark)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *GetGeoOptions) SetFormat(format string) *GetGeoOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetG : Allow user to set G
func (_options *GetGeoOptions) SetG(g string) *GetGeoOptions {
	_options.G = core.StringPtr(g)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *GetGeoOptions) SetIncludeDocs(includeDocs bool) *GetGeoOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetLat : Allow user to set Lat
func (_options *GetGeoOptions) SetLat(lat float64) *GetGeoOptions {
	_options.Lat = core.Float64Ptr(lat)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetGeoOptions) SetLimit(limit int64) *GetGeoOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetLon : Allow user to set Lon
func (_options *GetGeoOptions) SetLon(lon float64) *GetGeoOptions {
	_options.Lon = core.Float64Ptr(lon)
	return _options
}

// SetNearest : Allow user to set Nearest
func (_options *GetGeoOptions) SetNearest(nearest bool) *GetGeoOptions {
	_options.Nearest = core.BoolPtr(nearest)
	return _options
}

// SetRadius : Allow user to set Radius
func (_options *GetGeoOptions) SetRadius(radius float64) *GetGeoOptions {
	_options.Radius = core.Float64Ptr(radius)
	return _options
}

// SetRangex : Allow user to set Rangex
func (_options *GetGeoOptions) SetRangex(rangex float64) *GetGeoOptions {
	_options.Rangex = core.Float64Ptr(rangex)
	return _options
}

// SetRangey : Allow user to set Rangey
func (_options *GetGeoOptions) SetRangey(rangey float64) *GetGeoOptions {
	_options.Rangey = core.Float64Ptr(rangey)
	return _options
}

// SetRelation : Allow user to set Relation
func (_options *GetGeoOptions) SetRelation(relation string) *GetGeoOptions {
	_options.Relation = core.StringPtr(relation)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *GetGeoOptions) SetSkip(skip int64) *GetGeoOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetStale : Allow user to set Stale
func (_options *GetGeoOptions) SetStale(stale string) *GetGeoOptions {
	_options.Stale = core.StringPtr(stale)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetGeoOptions) SetHeaders(param map[string]string) *GetGeoOptions {
	options.Headers = param
	return options
}

// GetIndexesInformationOptions : The GetIndexesInformation options.
type GetIndexesInformationOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetIndexesInformationOptions : Instantiate GetIndexesInformationOptions
func (*CloudantV1) NewGetIndexesInformationOptions(db string) *GetIndexesInformationOptions {
	return &GetIndexesInformationOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *GetIndexesInformationOptions) SetDb(db string) *GetIndexesInformationOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetIndexesInformationOptions) SetHeaders(param map[string]string) *GetIndexesInformationOptions {
	options.Headers = param
	return options
}

// GetLocalDocumentOptions : The GetLocalDocument options.
type GetLocalDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// The type of the response: application/json, multipart/mixed, multipart/related, or application/octet-stream.
	Accept *string `json:"Accept,omitempty"`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Query parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Query parameter to specify whether to include the encoding information in attachment stubs if the particular
	// attachment is compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Query parameter to specify whether to include the last update sequence for the document.
	LocalSeq *bool `json:"local_seq,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLocalDocumentOptions : Instantiate GetLocalDocumentOptions
func (*CloudantV1) NewGetLocalDocumentOptions(db string, docID string) *GetLocalDocumentOptions {
	return &GetLocalDocumentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *GetLocalDocumentOptions) SetDb(db string) *GetLocalDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *GetLocalDocumentOptions) SetDocID(docID string) *GetLocalDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *GetLocalDocumentOptions) SetAccept(accept string) *GetLocalDocumentOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *GetLocalDocumentOptions) SetIfNoneMatch(ifNoneMatch string) *GetLocalDocumentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *GetLocalDocumentOptions) SetAttachments(attachments bool) *GetLocalDocumentOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *GetLocalDocumentOptions) SetAttEncodingInfo(attEncodingInfo bool) *GetLocalDocumentOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetLocalSeq : Allow user to set LocalSeq
func (_options *GetLocalDocumentOptions) SetLocalSeq(localSeq bool) *GetLocalDocumentOptions {
	_options.LocalSeq = core.BoolPtr(localSeq)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLocalDocumentOptions) SetHeaders(param map[string]string) *GetLocalDocumentOptions {
	options.Headers = param
	return options
}

// GetMembershipInformationOptions : The GetMembershipInformation options.
type GetMembershipInformationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetMembershipInformationOptions : Instantiate GetMembershipInformationOptions
func (*CloudantV1) NewGetMembershipInformationOptions() *GetMembershipInformationOptions {
	return &GetMembershipInformationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetMembershipInformationOptions) SetHeaders(param map[string]string) *GetMembershipInformationOptions {
	options.Headers = param
	return options
}

// GetPartitionInformationOptions : The GetPartitionInformation options.
type GetPartitionInformationOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the database partition key.
	PartitionKey *string `json:"partition_key" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetPartitionInformationOptions : Instantiate GetPartitionInformationOptions
func (*CloudantV1) NewGetPartitionInformationOptions(db string, partitionKey string) *GetPartitionInformationOptions {
	return &GetPartitionInformationOptions{
		Db: core.StringPtr(db),
		PartitionKey: core.StringPtr(partitionKey),
	}
}

// SetDb : Allow user to set Db
func (_options *GetPartitionInformationOptions) SetDb(db string) *GetPartitionInformationOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetPartitionKey : Allow user to set PartitionKey
func (_options *GetPartitionInformationOptions) SetPartitionKey(partitionKey string) *GetPartitionInformationOptions {
	_options.PartitionKey = core.StringPtr(partitionKey)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPartitionInformationOptions) SetHeaders(param map[string]string) *GetPartitionInformationOptions {
	options.Headers = param
	return options
}

// GetReplicationDocumentOptions : The GetReplicationDocument options.
type GetReplicationDocumentOptions struct {
	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Query parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Query parameter to specify whether to include the encoding information in attachment stubs if the particular
	// attachment is compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Query parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the
	// returned document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Query parameter to specify whether to include a list of deleted conflicted revisions in the `_deleted_conflicts`
	// property of the returned document.
	DeletedConflicts *bool `json:"deleted_conflicts,omitempty"`

	// Query parameter to specify whether to force retrieving latest leaf revision, no matter what rev was requested.
	Latest *bool `json:"latest,omitempty"`

	// Query parameter to specify whether to include the last update sequence for the document.
	LocalSeq *bool `json:"local_seq,omitempty"`

	// Query parameter to specify whether to include document meta information. Acts the same as specifying all of the
	// conflicts, deleted_conflicts and open_revs query parameters.
	Meta *bool `json:"meta,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Query parameter to specify whether to include a list of all known document revisions.
	Revs *bool `json:"revs,omitempty"`

	// Query parameter to specify whether to includes detailed information for all known document revisions.
	RevsInfo *bool `json:"revs_info,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReplicationDocumentOptions : Instantiate GetReplicationDocumentOptions
func (*CloudantV1) NewGetReplicationDocumentOptions(docID string) *GetReplicationDocumentOptions {
	return &GetReplicationDocumentOptions{
		DocID: core.StringPtr(docID),
	}
}

// SetDocID : Allow user to set DocID
func (_options *GetReplicationDocumentOptions) SetDocID(docID string) *GetReplicationDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *GetReplicationDocumentOptions) SetIfNoneMatch(ifNoneMatch string) *GetReplicationDocumentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *GetReplicationDocumentOptions) SetAttachments(attachments bool) *GetReplicationDocumentOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *GetReplicationDocumentOptions) SetAttEncodingInfo(attEncodingInfo bool) *GetReplicationDocumentOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *GetReplicationDocumentOptions) SetConflicts(conflicts bool) *GetReplicationDocumentOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDeletedConflicts : Allow user to set DeletedConflicts
func (_options *GetReplicationDocumentOptions) SetDeletedConflicts(deletedConflicts bool) *GetReplicationDocumentOptions {
	_options.DeletedConflicts = core.BoolPtr(deletedConflicts)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *GetReplicationDocumentOptions) SetLatest(latest bool) *GetReplicationDocumentOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetLocalSeq : Allow user to set LocalSeq
func (_options *GetReplicationDocumentOptions) SetLocalSeq(localSeq bool) *GetReplicationDocumentOptions {
	_options.LocalSeq = core.BoolPtr(localSeq)
	return _options
}

// SetMeta : Allow user to set Meta
func (_options *GetReplicationDocumentOptions) SetMeta(meta bool) *GetReplicationDocumentOptions {
	_options.Meta = core.BoolPtr(meta)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *GetReplicationDocumentOptions) SetRev(rev string) *GetReplicationDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetRevs : Allow user to set Revs
func (_options *GetReplicationDocumentOptions) SetRevs(revs bool) *GetReplicationDocumentOptions {
	_options.Revs = core.BoolPtr(revs)
	return _options
}

// SetRevsInfo : Allow user to set RevsInfo
func (_options *GetReplicationDocumentOptions) SetRevsInfo(revsInfo bool) *GetReplicationDocumentOptions {
	_options.RevsInfo = core.BoolPtr(revsInfo)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReplicationDocumentOptions) SetHeaders(param map[string]string) *GetReplicationDocumentOptions {
	options.Headers = param
	return options
}

// GetSchedulerDocsOptions : The GetSchedulerDocs options.
type GetSchedulerDocsOptions struct {
	// Query parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Query parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Query parameter to include only replication documents in the specified states. String must be a comma-delimited
	// string.
	States []string `json:"states,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSchedulerDocsOptions.States property.
// Schema for replication state.
const (
	GetSchedulerDocsOptionsStatesCompletedConst = "completed"
	GetSchedulerDocsOptionsStatesCrashingConst = "crashing"
	GetSchedulerDocsOptionsStatesErrorConst = "error"
	GetSchedulerDocsOptionsStatesFailedConst = "failed"
	GetSchedulerDocsOptionsStatesInitializingConst = "initializing"
	GetSchedulerDocsOptionsStatesPendingConst = "pending"
	GetSchedulerDocsOptionsStatesRunningConst = "running"
)

// NewGetSchedulerDocsOptions : Instantiate GetSchedulerDocsOptions
func (*CloudantV1) NewGetSchedulerDocsOptions() *GetSchedulerDocsOptions {
	return &GetSchedulerDocsOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *GetSchedulerDocsOptions) SetLimit(limit int64) *GetSchedulerDocsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *GetSchedulerDocsOptions) SetSkip(skip int64) *GetSchedulerDocsOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetStates : Allow user to set States
func (_options *GetSchedulerDocsOptions) SetStates(states []string) *GetSchedulerDocsOptions {
	_options.States = states
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSchedulerDocsOptions) SetHeaders(param map[string]string) *GetSchedulerDocsOptions {
	options.Headers = param
	return options
}

// GetSchedulerDocumentOptions : The GetSchedulerDocument options.
type GetSchedulerDocumentOptions struct {
	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSchedulerDocumentOptions : Instantiate GetSchedulerDocumentOptions
func (*CloudantV1) NewGetSchedulerDocumentOptions(docID string) *GetSchedulerDocumentOptions {
	return &GetSchedulerDocumentOptions{
		DocID: core.StringPtr(docID),
	}
}

// SetDocID : Allow user to set DocID
func (_options *GetSchedulerDocumentOptions) SetDocID(docID string) *GetSchedulerDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSchedulerDocumentOptions) SetHeaders(param map[string]string) *GetSchedulerDocumentOptions {
	options.Headers = param
	return options
}

// GetSchedulerJobOptions : The GetSchedulerJob options.
type GetSchedulerJobOptions struct {
	// Path parameter to specify the replication job id.
	JobID *string `json:"job_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSchedulerJobOptions : Instantiate GetSchedulerJobOptions
func (*CloudantV1) NewGetSchedulerJobOptions(jobID string) *GetSchedulerJobOptions {
	return &GetSchedulerJobOptions{
		JobID: core.StringPtr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *GetSchedulerJobOptions) SetJobID(jobID string) *GetSchedulerJobOptions {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSchedulerJobOptions) SetHeaders(param map[string]string) *GetSchedulerJobOptions {
	options.Headers = param
	return options
}

// GetSchedulerJobsOptions : The GetSchedulerJobs options.
type GetSchedulerJobsOptions struct {
	// Query parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Query parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSchedulerJobsOptions : Instantiate GetSchedulerJobsOptions
func (*CloudantV1) NewGetSchedulerJobsOptions() *GetSchedulerJobsOptions {
	return &GetSchedulerJobsOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *GetSchedulerJobsOptions) SetLimit(limit int64) *GetSchedulerJobsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *GetSchedulerJobsOptions) SetSkip(skip int64) *GetSchedulerJobsOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSchedulerJobsOptions) SetHeaders(param map[string]string) *GetSchedulerJobsOptions {
	options.Headers = param
	return options
}

// GetSearchInfoOptions : The GetSearchInfo options.
type GetSearchInfoOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the index name.
	Index *string `json:"index" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSearchInfoOptions : Instantiate GetSearchInfoOptions
func (*CloudantV1) NewGetSearchInfoOptions(db string, ddoc string, index string) *GetSearchInfoOptions {
	return &GetSearchInfoOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
		Index: core.StringPtr(index),
	}
}

// SetDb : Allow user to set Db
func (_options *GetSearchInfoOptions) SetDb(db string) *GetSearchInfoOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *GetSearchInfoOptions) SetDdoc(ddoc string) *GetSearchInfoOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetIndex : Allow user to set Index
func (_options *GetSearchInfoOptions) SetIndex(index string) *GetSearchInfoOptions {
	_options.Index = core.StringPtr(index)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSearchInfoOptions) SetHeaders(param map[string]string) *GetSearchInfoOptions {
	options.Headers = param
	return options
}

// GetSecurityOptions : The GetSecurity options.
type GetSecurityOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSecurityOptions : Instantiate GetSecurityOptions
func (*CloudantV1) NewGetSecurityOptions(db string) *GetSecurityOptions {
	return &GetSecurityOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *GetSecurityOptions) SetDb(db string) *GetSecurityOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecurityOptions) SetHeaders(param map[string]string) *GetSecurityOptions {
	options.Headers = param
	return options
}

// GetServerInformationOptions : The GetServerInformation options.
type GetServerInformationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetServerInformationOptions : Instantiate GetServerInformationOptions
func (*CloudantV1) NewGetServerInformationOptions() *GetServerInformationOptions {
	return &GetServerInformationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetServerInformationOptions) SetHeaders(param map[string]string) *GetServerInformationOptions {
	options.Headers = param
	return options
}

// GetSessionInformationOptions : The GetSessionInformation options.
type GetSessionInformationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSessionInformationOptions : Instantiate GetSessionInformationOptions
func (*CloudantV1) NewGetSessionInformationOptions() *GetSessionInformationOptions {
	return &GetSessionInformationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetSessionInformationOptions) SetHeaders(param map[string]string) *GetSessionInformationOptions {
	options.Headers = param
	return options
}

// GetShardsInformationOptions : The GetShardsInformation options.
type GetShardsInformationOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetShardsInformationOptions : Instantiate GetShardsInformationOptions
func (*CloudantV1) NewGetShardsInformationOptions(db string) *GetShardsInformationOptions {
	return &GetShardsInformationOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *GetShardsInformationOptions) SetDb(db string) *GetShardsInformationOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetShardsInformationOptions) SetHeaders(param map[string]string) *GetShardsInformationOptions {
	options.Headers = param
	return options
}

// GetUpInformationOptions : The GetUpInformation options.
type GetUpInformationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetUpInformationOptions : Instantiate GetUpInformationOptions
func (*CloudantV1) NewGetUpInformationOptions() *GetUpInformationOptions {
	return &GetUpInformationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetUpInformationOptions) SetHeaders(param map[string]string) *GetUpInformationOptions {
	options.Headers = param
	return options
}

// GetUuidsOptions : The GetUuids options.
type GetUuidsOptions struct {
	// Query parameter to specify the number of UUIDs to return.
	Count *int64 `json:"count,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetUuidsOptions : Instantiate GetUuidsOptions
func (*CloudantV1) NewGetUuidsOptions() *GetUuidsOptions {
	return &GetUuidsOptions{}
}

// SetCount : Allow user to set Count
func (_options *GetUuidsOptions) SetCount(count int64) *GetUuidsOptions {
	_options.Count = core.Int64Ptr(count)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetUuidsOptions) SetHeaders(param map[string]string) *GetUuidsOptions {
	options.Headers = param
	return options
}

// HeadAttachmentOptions : The HeadAttachment options.
type HeadAttachmentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Path parameter to specify the attachment name.
	AttachmentName *string `json:"attachment_name" validate:"required,ne="`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadAttachmentOptions : Instantiate HeadAttachmentOptions
func (*CloudantV1) NewHeadAttachmentOptions(db string, docID string, attachmentName string) *HeadAttachmentOptions {
	return &HeadAttachmentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
		AttachmentName: core.StringPtr(attachmentName),
	}
}

// SetDb : Allow user to set Db
func (_options *HeadAttachmentOptions) SetDb(db string) *HeadAttachmentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *HeadAttachmentOptions) SetDocID(docID string) *HeadAttachmentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetAttachmentName : Allow user to set AttachmentName
func (_options *HeadAttachmentOptions) SetAttachmentName(attachmentName string) *HeadAttachmentOptions {
	_options.AttachmentName = core.StringPtr(attachmentName)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *HeadAttachmentOptions) SetIfMatch(ifMatch string) *HeadAttachmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *HeadAttachmentOptions) SetIfNoneMatch(ifNoneMatch string) *HeadAttachmentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *HeadAttachmentOptions) SetRev(rev string) *HeadAttachmentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HeadAttachmentOptions) SetHeaders(param map[string]string) *HeadAttachmentOptions {
	options.Headers = param
	return options
}

// HeadDatabaseOptions : The HeadDatabase options.
type HeadDatabaseOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadDatabaseOptions : Instantiate HeadDatabaseOptions
func (*CloudantV1) NewHeadDatabaseOptions(db string) *HeadDatabaseOptions {
	return &HeadDatabaseOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *HeadDatabaseOptions) SetDb(db string) *HeadDatabaseOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HeadDatabaseOptions) SetHeaders(param map[string]string) *HeadDatabaseOptions {
	options.Headers = param
	return options
}

// HeadDesignDocumentOptions : The HeadDesignDocument options.
type HeadDesignDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadDesignDocumentOptions : Instantiate HeadDesignDocumentOptions
func (*CloudantV1) NewHeadDesignDocumentOptions(db string, ddoc string) *HeadDesignDocumentOptions {
	return &HeadDesignDocumentOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
	}
}

// SetDb : Allow user to set Db
func (_options *HeadDesignDocumentOptions) SetDb(db string) *HeadDesignDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *HeadDesignDocumentOptions) SetDdoc(ddoc string) *HeadDesignDocumentOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *HeadDesignDocumentOptions) SetIfNoneMatch(ifNoneMatch string) *HeadDesignDocumentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HeadDesignDocumentOptions) SetHeaders(param map[string]string) *HeadDesignDocumentOptions {
	options.Headers = param
	return options
}

// HeadDocumentOptions : The HeadDocument options.
type HeadDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Query parameter to specify whether to force retrieving latest leaf revision, no matter what rev was requested.
	Latest *bool `json:"latest,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadDocumentOptions : Instantiate HeadDocumentOptions
func (*CloudantV1) NewHeadDocumentOptions(db string, docID string) *HeadDocumentOptions {
	return &HeadDocumentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *HeadDocumentOptions) SetDb(db string) *HeadDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *HeadDocumentOptions) SetDocID(docID string) *HeadDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *HeadDocumentOptions) SetIfNoneMatch(ifNoneMatch string) *HeadDocumentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *HeadDocumentOptions) SetLatest(latest bool) *HeadDocumentOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *HeadDocumentOptions) SetRev(rev string) *HeadDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HeadDocumentOptions) SetHeaders(param map[string]string) *HeadDocumentOptions {
	options.Headers = param
	return options
}

// HeadLocalDocumentOptions : The HeadLocalDocument options.
type HeadLocalDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadLocalDocumentOptions : Instantiate HeadLocalDocumentOptions
func (*CloudantV1) NewHeadLocalDocumentOptions(db string, docID string) *HeadLocalDocumentOptions {
	return &HeadLocalDocumentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *HeadLocalDocumentOptions) SetDb(db string) *HeadLocalDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *HeadLocalDocumentOptions) SetDocID(docID string) *HeadLocalDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *HeadLocalDocumentOptions) SetIfNoneMatch(ifNoneMatch string) *HeadLocalDocumentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HeadLocalDocumentOptions) SetHeaders(param map[string]string) *HeadLocalDocumentOptions {
	options.Headers = param
	return options
}

// HeadReplicationDocumentOptions : The HeadReplicationDocument options.
type HeadReplicationDocumentOptions struct {
	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Header parameter to specify a double quoted document revision token for cache control.
	IfNoneMatch *string `json:"If-None-Match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadReplicationDocumentOptions : Instantiate HeadReplicationDocumentOptions
func (*CloudantV1) NewHeadReplicationDocumentOptions(docID string) *HeadReplicationDocumentOptions {
	return &HeadReplicationDocumentOptions{
		DocID: core.StringPtr(docID),
	}
}

// SetDocID : Allow user to set DocID
func (_options *HeadReplicationDocumentOptions) SetDocID(docID string) *HeadReplicationDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetIfNoneMatch : Allow user to set IfNoneMatch
func (_options *HeadReplicationDocumentOptions) SetIfNoneMatch(ifNoneMatch string) *HeadReplicationDocumentOptions {
	_options.IfNoneMatch = core.StringPtr(ifNoneMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HeadReplicationDocumentOptions) SetHeaders(param map[string]string) *HeadReplicationDocumentOptions {
	options.Headers = param
	return options
}

// HeadSchedulerDocumentOptions : The HeadSchedulerDocument options.
type HeadSchedulerDocumentOptions struct {
	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadSchedulerDocumentOptions : Instantiate HeadSchedulerDocumentOptions
func (*CloudantV1) NewHeadSchedulerDocumentOptions(docID string) *HeadSchedulerDocumentOptions {
	return &HeadSchedulerDocumentOptions{
		DocID: core.StringPtr(docID),
	}
}

// SetDocID : Allow user to set DocID
func (_options *HeadSchedulerDocumentOptions) SetDocID(docID string) *HeadSchedulerDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HeadSchedulerDocumentOptions) SetHeaders(param map[string]string) *HeadSchedulerDocumentOptions {
	options.Headers = param
	return options
}

// HeadSchedulerJobOptions : The HeadSchedulerJob options.
type HeadSchedulerJobOptions struct {
	// Path parameter to specify the replication job id.
	JobID *string `json:"job_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadSchedulerJobOptions : Instantiate HeadSchedulerJobOptions
func (*CloudantV1) NewHeadSchedulerJobOptions(jobID string) *HeadSchedulerJobOptions {
	return &HeadSchedulerJobOptions{
		JobID: core.StringPtr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *HeadSchedulerJobOptions) SetJobID(jobID string) *HeadSchedulerJobOptions {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *HeadSchedulerJobOptions) SetHeaders(param map[string]string) *HeadSchedulerJobOptions {
	options.Headers = param
	return options
}

// HeadUpInformationOptions : The HeadUpInformation options.
type HeadUpInformationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewHeadUpInformationOptions : Instantiate HeadUpInformationOptions
func (*CloudantV1) NewHeadUpInformationOptions() *HeadUpInformationOptions {
	return &HeadUpInformationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *HeadUpInformationOptions) SetHeaders(param map[string]string) *HeadUpInformationOptions {
	options.Headers = param
	return options
}

// IndexDefinition : Schema for a `json` or `text` query index definition. Indexes of type `text` have additional configuration properties
// that do not apply to `json` indexes, these are:
// * `default_analyzer` - the default text analyzer to use * `default_field` - whether to index the text in all document
// fields and what analyzer to use for that purpose.
type IndexDefinition struct {
	// Schema for a full text search analyzer.
	DefaultAnalyzer *Analyzer `json:"default_analyzer,omitempty"`

	// Schema for the text index default field configuration. The default field is used to index the text of all fields
	// within a document for use with the `$text` operator.
	DefaultField *IndexTextOperatorDefaultField `json:"default_field,omitempty"`

	// List of field objects to index.  Nested fields are also allowed, e.g. `person.name`.
	//
	// For "json" type indexes each object is a mapping of field name to sort direction (asc or desc).
	//
	// For "text" type indexes each object has a `name` property of the field name and a `type` property of the field type
	// (string, number, or boolean).
	Fields []IndexField `json:"fields,omitempty"`

	// Whether to scan every document for arrays and store the length for each array found. Set the index_array_lengths
	// field to false if:
	// * You do not need to know the length of an array. * You do not use the `$size` operator. * The documents in your
	// database are complex, or not completely under your control. As a result, it is difficult to estimate the impact of
	// the extra processing that is needed to determine and store the arrays lengths.
	IndexArrayLengths *bool `json:"index_array_lengths,omitempty"`

	// JSON object describing criteria used to select documents. The selector specifies fields in the document, and
	// provides an expression to evaluate with the field content or other data.
	//
	// The selector object must:
	//   * Be structured as valid JSON.
	//   * Contain a valid query expression.
	//
	// Using a selector is significantly more efficient than using a JavaScript filter function, and is the recommended
	// option if filtering on document attributes only.
	//
	// Elementary selector syntax requires you to specify one or more fields, and the corresponding values required for
	// those fields. You can create more complex selector expressions by combining operators.
	//
	// Operators are identified by the use of a dollar sign `$` prefix in the name field.
	//
	// There are two core types of operators in the selector syntax:
	// * Combination operators: applied at the topmost level of selection. They are used to combine selectors. In addition
	// to the common boolean operators (`$and`, `$or`, `$not`, `$nor`) there are three combination operators: `$all`,
	// `$elemMatch`, and `$allMatch`. A combination operator takes a single argument. The argument is either another
	// selector, or an array of selectors.
	// * Condition operators: are specific to a field, and are used to evaluate the value stored in that field. For
	// instance, the basic `$eq` operator matches when the specified field contains a value that is equal to the supplied
	// argument.
	PartialFilterSelector map[string]interface{} `json:"partial_filter_selector,omitempty"`
}

// UnmarshalIndexDefinition unmarshals an instance of IndexDefinition from the specified map of raw messages.
func UnmarshalIndexDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IndexDefinition)
	err = core.UnmarshalModel(m, "default_analyzer", &obj.DefaultAnalyzer, UnmarshalAnalyzer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "default_field", &obj.DefaultField, UnmarshalIndexTextOperatorDefaultField)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "fields", &obj.Fields, UnmarshalIndexField)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "index_array_lengths", &obj.IndexArrayLengths)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "partial_filter_selector", &obj.PartialFilterSelector)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IndexField : Schema for indexed fields for use with declarative JSON query.
type IndexField struct {
	// Name of the field.
	Name *string `json:"name,omitempty"`

	// The type of the named field.
	Type *string `json:"type,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]*string
}

// Constants associated with the IndexField.Type property.
// The type of the named field.
const (
	IndexFieldTypeBooleanConst = "boolean"
	IndexFieldTypeNumberConst = "number"
	IndexFieldTypeStringConst = "string"
)

// SetProperty allows the user to set an arbitrary property on an instance of IndexField
func (o *IndexField) SetProperty(key string, value *string) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]*string)
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of IndexField
func (o *IndexField) SetProperties(m map[string]*string) {
	o.additionalProperties = make(map[string]*string)
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of IndexField
func (o *IndexField) GetProperty(key string) *string {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of IndexField
func (o *IndexField) GetProperties() map[string]*string {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of IndexField
func (o *IndexField) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Name != nil {
		m["name"] = o.Name
	}
	if o.Type != nil {
		m["type"] = o.Type
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalIndexField unmarshals an instance of IndexField from the specified map of raw messages.
func UnmarshalIndexField(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IndexField)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	delete(m, "name")
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	delete(m, "type")
	for k := range m {
		var v *string
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

// IndexInformation : Schema for information about an index.
type IndexInformation struct {
	// Design document ID.
	Ddoc *string `json:"ddoc" validate:"required"`

	// Schema for a `json` or `text` query index definition. Indexes of type `text` have additional configuration
	// properties that do not apply to `json` indexes, these are:
	// * `default_analyzer` - the default text analyzer to use * `default_field` - whether to index the text in all
	// document fields and what analyzer to use for that purpose.
	Def *IndexDefinition `json:"def" validate:"required"`

	// Index name.
	Name *string `json:"name" validate:"required"`

	// Schema for the type of an index.
	Type *string `json:"type" validate:"required"`
}

// Constants associated with the IndexInformation.Type property.
// Schema for the type of an index.
const (
	IndexInformationTypeJSONConst = "json"
	IndexInformationTypeSpecialConst = "special"
	IndexInformationTypeTextConst = "text"
)

// UnmarshalIndexInformation unmarshals an instance of IndexInformation from the specified map of raw messages.
func UnmarshalIndexInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IndexInformation)
	err = core.UnmarshalPrimitive(m, "ddoc", &obj.Ddoc)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "def", &obj.Def, UnmarshalIndexDefinition)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IndexResult : Schema for the result of creating an index.
type IndexResult struct {
	// Id of the design document the index was created in.
	ID *string `json:"id" validate:"required"`

	// Name of the index created.
	Name *string `json:"name" validate:"required"`

	// Flag to show whether the index was created or one already exists.
	Result *string `json:"result" validate:"required"`
}

// Constants associated with the IndexResult.Result property.
// Flag to show whether the index was created or one already exists.
const (
	IndexResultResultCreatedConst = "created"
	IndexResultResultExistsConst = "exists"
)

// UnmarshalIndexResult unmarshals an instance of IndexResult from the specified map of raw messages.
func UnmarshalIndexResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IndexResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IndexTextOperatorDefaultField : Schema for the text index default field configuration. The default field is used to index the text of all fields
// within a document for use with the `$text` operator.
type IndexTextOperatorDefaultField struct {
	// Schema for a full text search analyzer.
	Analyzer *Analyzer `json:"analyzer,omitempty"`

	// Whether or not the default_field is enabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// UnmarshalIndexTextOperatorDefaultField unmarshals an instance of IndexTextOperatorDefaultField from the specified map of raw messages.
func UnmarshalIndexTextOperatorDefaultField(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IndexTextOperatorDefaultField)
	err = core.UnmarshalModel(m, "analyzer", &obj.Analyzer, UnmarshalAnalyzer)
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

// IndexesInformation : Schema for information about the indexes in a database.
type IndexesInformation struct {
	// Number of total rows.
	TotalRows *int64 `json:"total_rows" validate:"required"`

	// Indexes.
	Indexes []IndexInformation `json:"indexes" validate:"required"`
}

// UnmarshalIndexesInformation unmarshals an instance of IndexesInformation from the specified map of raw messages.
func UnmarshalIndexesInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IndexesInformation)
	err = core.UnmarshalPrimitive(m, "total_rows", &obj.TotalRows)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "indexes", &obj.Indexes, UnmarshalIndexInformation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// MembershipInformation : Schema for information about known nodes and cluster membership.
type MembershipInformation struct {
	// List of nodes this node knows about, including the ones that are part of the cluster.
	AllNodes []string `json:"all_nodes" validate:"required"`

	// All cluster nodes.
	ClusterNodes []string `json:"cluster_nodes" validate:"required"`
}

// UnmarshalMembershipInformation unmarshals an instance of MembershipInformation from the specified map of raw messages.
func UnmarshalMembershipInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MembershipInformation)
	err = core.UnmarshalPrimitive(m, "all_nodes", &obj.AllNodes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_nodes", &obj.ClusterNodes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Ok : Schema for an OK result.
type Ok struct {
	// ok.
	Ok *bool `json:"ok,omitempty"`
}

// UnmarshalOk unmarshals an instance of Ok from the specified map of raw messages.
func UnmarshalOk(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Ok)
	err = core.UnmarshalPrimitive(m, "ok", &obj.Ok)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PartitionInformation : Schema for information about a database partition.
type PartitionInformation struct {
	// The name of the database.
	DbName *string `json:"db_name" validate:"required"`

	// A count of the documents in the specified database partition.
	DocCount *int64 `json:"doc_count" validate:"required"`

	// Number of deleted documents.
	DocDelCount *int64 `json:"doc_del_count" validate:"required"`

	// The name of the partition in the database.
	Partition *string `json:"partition" validate:"required"`

	// Schema for information about the partition index count and limit in a database.
	PartitionedIndexes *PartitionInformationIndexes `json:"partitioned_indexes,omitempty"`

	// The size of active and external data, in bytes.
	Sizes *PartitionInformationSizes `json:"sizes" validate:"required"`
}

// UnmarshalPartitionInformation unmarshals an instance of PartitionInformation from the specified map of raw messages.
func UnmarshalPartitionInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PartitionInformation)
	err = core.UnmarshalPrimitive(m, "db_name", &obj.DbName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_count", &obj.DocCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_del_count", &obj.DocDelCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "partition", &obj.Partition)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "partitioned_indexes", &obj.PartitionedIndexes, UnmarshalPartitionInformationIndexes)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sizes", &obj.Sizes, UnmarshalPartitionInformationSizes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PartitionInformationIndexes : Schema for information about the partition index count and limit in a database.
type PartitionInformationIndexes struct {
	// Total count of the partitioned indexes.
	Count *int64 `json:"count,omitempty"`

	// The count breakdown of partitioned indexes.
	Indexes *PartitionInformationIndexesIndexes `json:"indexes,omitempty"`

	// The partitioned index limit.
	Limit *int64 `json:"limit,omitempty"`
}

// UnmarshalPartitionInformationIndexes unmarshals an instance of PartitionInformationIndexes from the specified map of raw messages.
func UnmarshalPartitionInformationIndexes(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PartitionInformationIndexes)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "indexes", &obj.Indexes, UnmarshalPartitionInformationIndexesIndexes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PartitionInformationIndexesIndexes : The count breakdown of partitioned indexes.
type PartitionInformationIndexesIndexes struct {
	// Number of partitioned search indexes.
	Search *int64 `json:"search,omitempty"`

	// Number of partitioned view indexes.
	View *int64 `json:"view,omitempty"`
}

// UnmarshalPartitionInformationIndexesIndexes unmarshals an instance of PartitionInformationIndexesIndexes from the specified map of raw messages.
func UnmarshalPartitionInformationIndexesIndexes(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PartitionInformationIndexesIndexes)
	err = core.UnmarshalPrimitive(m, "search", &obj.Search)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "view", &obj.View)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PartitionInformationSizes : The size of active and external data, in bytes.
type PartitionInformationSizes struct {
	// The size of live data inside the database, in bytes.
	Active *int64 `json:"active,omitempty"`

	// The uncompressed size of database contents in bytes.
	External *int64 `json:"external,omitempty"`
}

// UnmarshalPartitionInformationSizes unmarshals an instance of PartitionInformationSizes from the specified map of raw messages.
func UnmarshalPartitionInformationSizes(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PartitionInformationSizes)
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "external", &obj.External)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PostActivityTrackerEventsOptions : The PostActivityTrackerEvents options.
type PostActivityTrackerEventsOptions struct {
	// An array of event types that are being sent to IBM Cloud Activity Tracker for the IBM Cloudant instance.
	// "management" is a required element of this array.
	Types []string `json:"types" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostActivityTrackerEventsOptions.Types property.
const (
	PostActivityTrackerEventsOptionsTypesDataConst = "data"
	PostActivityTrackerEventsOptionsTypesManagementConst = "management"
)

// NewPostActivityTrackerEventsOptions : Instantiate PostActivityTrackerEventsOptions
func (*CloudantV1) NewPostActivityTrackerEventsOptions(types []string) *PostActivityTrackerEventsOptions {
	return &PostActivityTrackerEventsOptions{
		Types: types,
	}
}

// SetTypes : Allow user to set Types
func (_options *PostActivityTrackerEventsOptions) SetTypes(types []string) *PostActivityTrackerEventsOptions {
	_options.Types = types
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostActivityTrackerEventsOptions) SetHeaders(param map[string]string) *PostActivityTrackerEventsOptions {
	options.Headers = param
	return options
}

// PostAllDocsOptions : The PostAllDocs options.
type PostAllDocsOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Parameter to specify whether to include the encoding information in attachment stubs if the particular attachment is
	// compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the returned
	// document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Parameter to specify whether the specified end key should be included in the result.
	InclusiveEnd *bool `json:"inclusive_end,omitempty"`

	// Parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Parameter to specify whether to include in the response an update_seq value indicating the sequence id of the
	// database the view reflects.
	UpdateSeq *bool `json:"update_seq,omitempty"`

	// Schema for a document ID.
	Endkey *string `json:"endkey,omitempty"`

	// Schema for a document ID.
	Key *string `json:"key,omitempty"`

	// Schema for a list of document IDs.
	Keys []string `json:"keys,omitempty"`

	// Schema for a document ID.
	Startkey *string `json:"startkey,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostAllDocsOptions : Instantiate PostAllDocsOptions
func (*CloudantV1) NewPostAllDocsOptions(db string) *PostAllDocsOptions {
	return &PostAllDocsOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *PostAllDocsOptions) SetDb(db string) *PostAllDocsOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *PostAllDocsOptions) SetAttEncodingInfo(attEncodingInfo bool) *PostAllDocsOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *PostAllDocsOptions) SetAttachments(attachments bool) *PostAllDocsOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostAllDocsOptions) SetConflicts(conflicts bool) *PostAllDocsOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDescending : Allow user to set Descending
func (_options *PostAllDocsOptions) SetDescending(descending bool) *PostAllDocsOptions {
	_options.Descending = core.BoolPtr(descending)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *PostAllDocsOptions) SetIncludeDocs(includeDocs bool) *PostAllDocsOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetInclusiveEnd : Allow user to set InclusiveEnd
func (_options *PostAllDocsOptions) SetInclusiveEnd(inclusiveEnd bool) *PostAllDocsOptions {
	_options.InclusiveEnd = core.BoolPtr(inclusiveEnd)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostAllDocsOptions) SetLimit(limit int64) *PostAllDocsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *PostAllDocsOptions) SetSkip(skip int64) *PostAllDocsOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetUpdateSeq : Allow user to set UpdateSeq
func (_options *PostAllDocsOptions) SetUpdateSeq(updateSeq bool) *PostAllDocsOptions {
	_options.UpdateSeq = core.BoolPtr(updateSeq)
	return _options
}

// SetEndkey : Allow user to set Endkey
func (_options *PostAllDocsOptions) SetEndkey(endkey string) *PostAllDocsOptions {
	_options.Endkey = core.StringPtr(endkey)
	return _options
}

// SetKey : Allow user to set Key
func (_options *PostAllDocsOptions) SetKey(key string) *PostAllDocsOptions {
	_options.Key = core.StringPtr(key)
	return _options
}

// SetKeys : Allow user to set Keys
func (_options *PostAllDocsOptions) SetKeys(keys []string) *PostAllDocsOptions {
	_options.Keys = keys
	return _options
}

// SetStartkey : Allow user to set Startkey
func (_options *PostAllDocsOptions) SetStartkey(startkey string) *PostAllDocsOptions {
	_options.Startkey = core.StringPtr(startkey)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostAllDocsOptions) SetHeaders(param map[string]string) *PostAllDocsOptions {
	options.Headers = param
	return options
}

// PostAllDocsQueriesOptions : The PostAllDocsQueries options.
type PostAllDocsQueriesOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// An array of query objects with fields for the parameters of each individual view query to be executed. The field
	// names and their meaning are the same as the query parameters of a regular `/_all_docs` request.
	Queries []AllDocsQuery `json:"queries" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostAllDocsQueriesOptions : Instantiate PostAllDocsQueriesOptions
func (*CloudantV1) NewPostAllDocsQueriesOptions(db string, queries []AllDocsQuery) *PostAllDocsQueriesOptions {
	return &PostAllDocsQueriesOptions{
		Db: core.StringPtr(db),
		Queries: queries,
	}
}

// SetDb : Allow user to set Db
func (_options *PostAllDocsQueriesOptions) SetDb(db string) *PostAllDocsQueriesOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetQueries : Allow user to set Queries
func (_options *PostAllDocsQueriesOptions) SetQueries(queries []AllDocsQuery) *PostAllDocsQueriesOptions {
	_options.Queries = queries
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostAllDocsQueriesOptions) SetHeaders(param map[string]string) *PostAllDocsQueriesOptions {
	options.Headers = param
	return options
}

// PostApiKeysOptions : The PostApiKeys options.
type PostApiKeysOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostApiKeysOptions : Instantiate PostApiKeysOptions
func (*CloudantV1) NewPostApiKeysOptions() *PostApiKeysOptions {
	return &PostApiKeysOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *PostApiKeysOptions) SetHeaders(param map[string]string) *PostApiKeysOptions {
	options.Headers = param
	return options
}

// PostBulkDocsOptions : The PostBulkDocs options.
type PostBulkDocsOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// HTTP request body for postBulkDocs.
	BulkDocs *BulkDocs `json:"bulkDocs,omitempty"`

	// HTTP request body for postBulkDocs.
	Body io.ReadCloser `json:"body,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostBulkDocsOptions : Instantiate PostBulkDocsOptions
func (*CloudantV1) NewPostBulkDocsOptions(db string) *PostBulkDocsOptions {
	return &PostBulkDocsOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *PostBulkDocsOptions) SetDb(db string) *PostBulkDocsOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetBulkDocs : Allow user to set BulkDocs
func (_options *PostBulkDocsOptions) SetBulkDocs(bulkDocs *BulkDocs) *PostBulkDocsOptions {
	_options.BulkDocs = bulkDocs
	return _options
}

// SetBody : Allow user to set Body
func (_options *PostBulkDocsOptions) SetBody(body io.ReadCloser) *PostBulkDocsOptions {
	_options.Body = body
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostBulkDocsOptions) SetHeaders(param map[string]string) *PostBulkDocsOptions {
	options.Headers = param
	return options
}

// PostBulkGetOptions : The PostBulkGet options.
type PostBulkGetOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// List of document items to get in bulk.
	Docs []BulkGetQueryDocument `json:"docs" validate:"required"`

	// Query parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Query parameter to specify whether to include the encoding information in attachment stubs if the particular
	// attachment is compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Query parameter to specify whether to force retrieving latest leaf revision, no matter what rev was requested.
	Latest *bool `json:"latest,omitempty"`

	// Query parameter to specify whether to include a list of all known document revisions.
	Revs *bool `json:"revs,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostBulkGetOptions : Instantiate PostBulkGetOptions
func (*CloudantV1) NewPostBulkGetOptions(db string, docs []BulkGetQueryDocument) *PostBulkGetOptions {
	return &PostBulkGetOptions{
		Db: core.StringPtr(db),
		Docs: docs,
	}
}

// SetDb : Allow user to set Db
func (_options *PostBulkGetOptions) SetDb(db string) *PostBulkGetOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocs : Allow user to set Docs
func (_options *PostBulkGetOptions) SetDocs(docs []BulkGetQueryDocument) *PostBulkGetOptions {
	_options.Docs = docs
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *PostBulkGetOptions) SetAttachments(attachments bool) *PostBulkGetOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *PostBulkGetOptions) SetAttEncodingInfo(attEncodingInfo bool) *PostBulkGetOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *PostBulkGetOptions) SetLatest(latest bool) *PostBulkGetOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetRevs : Allow user to set Revs
func (_options *PostBulkGetOptions) SetRevs(revs bool) *PostBulkGetOptions {
	_options.Revs = core.BoolPtr(revs)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostBulkGetOptions) SetHeaders(param map[string]string) *PostBulkGetOptions {
	options.Headers = param
	return options
}

// PostChangesOptions : The PostChanges options.
type PostChangesOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Schema for a list of document IDs.
	DocIds []string `json:"doc_ids,omitempty"`

	// JSON array that uses the field syntax. Use this parameter to specify which fields of a document must be returned. If
	// it is omitted, the entire document is returned.
	Fields []string `json:"fields,omitempty"`

	// JSON object describing criteria used to select documents. The selector specifies fields in the document, and
	// provides an expression to evaluate with the field content or other data.
	//
	// The selector object must:
	//   * Be structured as valid JSON.
	//   * Contain a valid query expression.
	//
	// Using a selector is significantly more efficient than using a JavaScript filter function, and is the recommended
	// option if filtering on document attributes only.
	//
	// Elementary selector syntax requires you to specify one or more fields, and the corresponding values required for
	// those fields. You can create more complex selector expressions by combining operators.
	//
	// Operators are identified by the use of a dollar sign `$` prefix in the name field.
	//
	// There are two core types of operators in the selector syntax:
	// * Combination operators: applied at the topmost level of selection. They are used to combine selectors. In addition
	// to the common boolean operators (`$and`, `$or`, `$not`, `$nor`) there are three combination operators: `$all`,
	// `$elemMatch`, and `$allMatch`. A combination operator takes a single argument. The argument is either another
	// selector, or an array of selectors.
	// * Condition operators: are specific to a field, and are used to evaluate the value stored in that field. For
	// instance, the basic `$eq` operator matches when the specified field contains a value that is equal to the supplied
	// argument.
	Selector map[string]interface{} `json:"selector,omitempty"`

	// Header parameter to specify the ID of the last events received by the server on a previous connection. Overrides
	// `since` query parameter.
	LastEventID *string `json:"Last-Event-ID,omitempty"`

	// Query parameter to specify whether to include the encoding information in attachment stubs if the particular
	// attachment is compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Query parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Query parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the
	// returned document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Query parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Query parameter to specify the changes feed type.
	Feed *string `json:"feed,omitempty"`

	// Query parameter to specify a filter function from a design document that will filter the changes stream emitting
	// only filtered events. For example: `design_doc/filtername`.
	//
	// Additionally, some keywords are reserved for built-in filters:
	//
	//   * `_design` - Returns only changes to design documents.
	//   * `_doc_ids` - Returns changes for documents with an ID matching one specified in
	//       `doc_ids` request body parameter.
	//   * `_selector` - Returns changes for documents that match the `selector`
	//       request body parameter. The selector syntax is the same as used for
	//       `_find`.
	//   * `_view` - Returns changes for documents that match an existing map
	//       function in the view specified by the query parameter `view`.
	Filter *string `json:"filter,omitempty"`

	// Query parameter to specify the period in milliseconds after which an empty line is sent in the results. Only
	// applicable for longpoll, continuous, and eventsource feeds. Overrides any timeout to keep the feed alive
	// indefinitely. May also be `true` to use default value of 60000.
	Heartbeat *int64 `json:"heartbeat,omitempty"`

	// Query parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Query parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Query parameter to specify that the update seq should only be calculated with every Nth result returned. When
	// fetching changes in a batch, setting <code>seq_interval=&lt;batch size&gt;</code>, where &lt;batch size&gt; is the
	// number of results requested per batch, load can be reduced on the source database as computing the seq value across
	// many shards (especially in highly-sharded databases) is expensive.
	SeqInterval *int64 `json:"seq_interval,omitempty"`

	// Query parameter to specify to start the results from the change immediately after the given update sequence. Can be
	// a valid update sequence or `now` value. Default is `0` i.e. all changes.
	Since *string `json:"since,omitempty"`

	// Query parameter to specify how many revisions are returned in the changes array. The default, `main_only`, will only
	// return the current "winning" revision; all_docs will return all leaf revisions (including conflicts and deleted
	// former conflicts).
	Style *string `json:"style,omitempty"`

	// Query parameter to specify the maximum period in milliseconds to wait for a change before the response is sent, even
	// if there are no results. Only applicable for `longpoll` or `continuous` feeds. Default value is specified by
	// `httpd/changes_timeout` configuration option. Note that `60000` value is also the default maximum timeout to prevent
	// undetected dead connections.
	Timeout *int64 `json:"timeout,omitempty"`

	// Query parameter to specify a view function as a filter. Documents pass the filter if the view's map function emits
	// at least one record for them.
	View *string `json:"view,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostChangesOptions.Feed property.
// Query parameter to specify the changes feed type.
const (
	PostChangesOptionsFeedContinuousConst = "continuous"
	PostChangesOptionsFeedEventsourceConst = "eventsource"
	PostChangesOptionsFeedLongpollConst = "longpoll"
	PostChangesOptionsFeedNormalConst = "normal"
)

// NewPostChangesOptions : Instantiate PostChangesOptions
func (*CloudantV1) NewPostChangesOptions(db string) *PostChangesOptions {
	return &PostChangesOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *PostChangesOptions) SetDb(db string) *PostChangesOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocIds : Allow user to set DocIds
func (_options *PostChangesOptions) SetDocIds(docIds []string) *PostChangesOptions {
	_options.DocIds = docIds
	return _options
}

// SetFields : Allow user to set Fields
func (_options *PostChangesOptions) SetFields(fields []string) *PostChangesOptions {
	_options.Fields = fields
	return _options
}

// SetSelector : Allow user to set Selector
func (_options *PostChangesOptions) SetSelector(selector map[string]interface{}) *PostChangesOptions {
	_options.Selector = selector
	return _options
}

// SetLastEventID : Allow user to set LastEventID
func (_options *PostChangesOptions) SetLastEventID(lastEventID string) *PostChangesOptions {
	_options.LastEventID = core.StringPtr(lastEventID)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *PostChangesOptions) SetAttEncodingInfo(attEncodingInfo bool) *PostChangesOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *PostChangesOptions) SetAttachments(attachments bool) *PostChangesOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostChangesOptions) SetConflicts(conflicts bool) *PostChangesOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDescending : Allow user to set Descending
func (_options *PostChangesOptions) SetDescending(descending bool) *PostChangesOptions {
	_options.Descending = core.BoolPtr(descending)
	return _options
}

// SetFeed : Allow user to set Feed
func (_options *PostChangesOptions) SetFeed(feed string) *PostChangesOptions {
	_options.Feed = core.StringPtr(feed)
	return _options
}

// SetFilter : Allow user to set Filter
func (_options *PostChangesOptions) SetFilter(filter string) *PostChangesOptions {
	_options.Filter = core.StringPtr(filter)
	return _options
}

// SetHeartbeat : Allow user to set Heartbeat
func (_options *PostChangesOptions) SetHeartbeat(heartbeat int64) *PostChangesOptions {
	_options.Heartbeat = core.Int64Ptr(heartbeat)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *PostChangesOptions) SetIncludeDocs(includeDocs bool) *PostChangesOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostChangesOptions) SetLimit(limit int64) *PostChangesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSeqInterval : Allow user to set SeqInterval
func (_options *PostChangesOptions) SetSeqInterval(seqInterval int64) *PostChangesOptions {
	_options.SeqInterval = core.Int64Ptr(seqInterval)
	return _options
}

// SetSince : Allow user to set Since
func (_options *PostChangesOptions) SetSince(since string) *PostChangesOptions {
	_options.Since = core.StringPtr(since)
	return _options
}

// SetStyle : Allow user to set Style
func (_options *PostChangesOptions) SetStyle(style string) *PostChangesOptions {
	_options.Style = core.StringPtr(style)
	return _options
}

// SetTimeout : Allow user to set Timeout
func (_options *PostChangesOptions) SetTimeout(timeout int64) *PostChangesOptions {
	_options.Timeout = core.Int64Ptr(timeout)
	return _options
}

// SetView : Allow user to set View
func (_options *PostChangesOptions) SetView(view string) *PostChangesOptions {
	_options.View = core.StringPtr(view)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostChangesOptions) SetHeaders(param map[string]string) *PostChangesOptions {
	options.Headers = param
	return options
}

// PostDbsInfoOptions : The PostDbsInfo options.
type PostDbsInfoOptions struct {
	// A list of database names.
	Keys []string `json:"keys" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostDbsInfoOptions : Instantiate PostDbsInfoOptions
func (*CloudantV1) NewPostDbsInfoOptions(keys []string) *PostDbsInfoOptions {
	return &PostDbsInfoOptions{
		Keys: keys,
	}
}

// SetKeys : Allow user to set Keys
func (_options *PostDbsInfoOptions) SetKeys(keys []string) *PostDbsInfoOptions {
	_options.Keys = keys
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostDbsInfoOptions) SetHeaders(param map[string]string) *PostDbsInfoOptions {
	options.Headers = param
	return options
}

// PostDesignDocsOptions : The PostDesignDocs options.
type PostDesignDocsOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// The type of the response: application/json or application/octet-stream.
	Accept *string `json:"Accept,omitempty"`

	// Parameter to specify whether to include the encoding information in attachment stubs if the particular attachment is
	// compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the returned
	// document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Parameter to specify whether the specified end key should be included in the result.
	InclusiveEnd *bool `json:"inclusive_end,omitempty"`

	// Parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Parameter to specify whether to include in the response an update_seq value indicating the sequence id of the
	// database the view reflects.
	UpdateSeq *bool `json:"update_seq,omitempty"`

	// Schema for a document ID.
	Endkey *string `json:"endkey,omitempty"`

	// Schema for a document ID.
	Key *string `json:"key,omitempty"`

	// Schema for a list of document IDs.
	Keys []string `json:"keys,omitempty"`

	// Schema for a document ID.
	Startkey *string `json:"startkey,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostDesignDocsOptions : Instantiate PostDesignDocsOptions
func (*CloudantV1) NewPostDesignDocsOptions(db string) *PostDesignDocsOptions {
	return &PostDesignDocsOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *PostDesignDocsOptions) SetDb(db string) *PostDesignDocsOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *PostDesignDocsOptions) SetAccept(accept string) *PostDesignDocsOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *PostDesignDocsOptions) SetAttEncodingInfo(attEncodingInfo bool) *PostDesignDocsOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *PostDesignDocsOptions) SetAttachments(attachments bool) *PostDesignDocsOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostDesignDocsOptions) SetConflicts(conflicts bool) *PostDesignDocsOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDescending : Allow user to set Descending
func (_options *PostDesignDocsOptions) SetDescending(descending bool) *PostDesignDocsOptions {
	_options.Descending = core.BoolPtr(descending)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *PostDesignDocsOptions) SetIncludeDocs(includeDocs bool) *PostDesignDocsOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetInclusiveEnd : Allow user to set InclusiveEnd
func (_options *PostDesignDocsOptions) SetInclusiveEnd(inclusiveEnd bool) *PostDesignDocsOptions {
	_options.InclusiveEnd = core.BoolPtr(inclusiveEnd)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostDesignDocsOptions) SetLimit(limit int64) *PostDesignDocsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *PostDesignDocsOptions) SetSkip(skip int64) *PostDesignDocsOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetUpdateSeq : Allow user to set UpdateSeq
func (_options *PostDesignDocsOptions) SetUpdateSeq(updateSeq bool) *PostDesignDocsOptions {
	_options.UpdateSeq = core.BoolPtr(updateSeq)
	return _options
}

// SetEndkey : Allow user to set Endkey
func (_options *PostDesignDocsOptions) SetEndkey(endkey string) *PostDesignDocsOptions {
	_options.Endkey = core.StringPtr(endkey)
	return _options
}

// SetKey : Allow user to set Key
func (_options *PostDesignDocsOptions) SetKey(key string) *PostDesignDocsOptions {
	_options.Key = core.StringPtr(key)
	return _options
}

// SetKeys : Allow user to set Keys
func (_options *PostDesignDocsOptions) SetKeys(keys []string) *PostDesignDocsOptions {
	_options.Keys = keys
	return _options
}

// SetStartkey : Allow user to set Startkey
func (_options *PostDesignDocsOptions) SetStartkey(startkey string) *PostDesignDocsOptions {
	_options.Startkey = core.StringPtr(startkey)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostDesignDocsOptions) SetHeaders(param map[string]string) *PostDesignDocsOptions {
	options.Headers = param
	return options
}

// PostDesignDocsQueriesOptions : The PostDesignDocsQueries options.
type PostDesignDocsQueriesOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// An array of query objects with fields for the parameters of each individual view query to be executed. The field
	// names and their meaning are the same as the query parameters of a regular `/_all_docs` request.
	Queries []AllDocsQuery `json:"queries" validate:"required"`

	// The type of the response: application/json or application/octet-stream.
	Accept *string `json:"Accept,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostDesignDocsQueriesOptions : Instantiate PostDesignDocsQueriesOptions
func (*CloudantV1) NewPostDesignDocsQueriesOptions(db string, queries []AllDocsQuery) *PostDesignDocsQueriesOptions {
	return &PostDesignDocsQueriesOptions{
		Db: core.StringPtr(db),
		Queries: queries,
	}
}

// SetDb : Allow user to set Db
func (_options *PostDesignDocsQueriesOptions) SetDb(db string) *PostDesignDocsQueriesOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetQueries : Allow user to set Queries
func (_options *PostDesignDocsQueriesOptions) SetQueries(queries []AllDocsQuery) *PostDesignDocsQueriesOptions {
	_options.Queries = queries
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *PostDesignDocsQueriesOptions) SetAccept(accept string) *PostDesignDocsQueriesOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostDesignDocsQueriesOptions) SetHeaders(param map[string]string) *PostDesignDocsQueriesOptions {
	options.Headers = param
	return options
}

// PostDocumentOptions : The PostDocument options.
type PostDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// HTTP request body for Document operations.
	Document *Document `json:"document,omitempty"`

	// HTTP request body for Document operations.
	Body io.ReadCloser `json:"body,omitempty"`

	// The type of the input.
	ContentType *string `json:"Content-Type,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	PostDocumentOptionsBatchOkConst = "ok"
)

// NewPostDocumentOptions : Instantiate PostDocumentOptions
func (*CloudantV1) NewPostDocumentOptions(db string) *PostDocumentOptions {
	return &PostDocumentOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *PostDocumentOptions) SetDb(db string) *PostDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocument : Allow user to set Document
func (_options *PostDocumentOptions) SetDocument(document *Document) *PostDocumentOptions {
	_options.Document = document
	return _options
}

// SetBody : Allow user to set Body
func (_options *PostDocumentOptions) SetBody(body io.ReadCloser) *PostDocumentOptions {
	_options.Body = body
	return _options
}

// SetContentType : Allow user to set ContentType
func (_options *PostDocumentOptions) SetContentType(contentType string) *PostDocumentOptions {
	_options.ContentType = core.StringPtr(contentType)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *PostDocumentOptions) SetBatch(batch string) *PostDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostDocumentOptions) SetHeaders(param map[string]string) *PostDocumentOptions {
	options.Headers = param
	return options
}

// PostExplainOptions : The PostExplain options.
type PostExplainOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// JSON object describing criteria used to select documents. The selector specifies fields in the document, and
	// provides an expression to evaluate with the field content or other data.
	//
	// The selector object must:
	//   * Be structured as valid JSON.
	//   * Contain a valid query expression.
	//
	// Using a selector is significantly more efficient than using a JavaScript filter function, and is the recommended
	// option if filtering on document attributes only.
	//
	// Elementary selector syntax requires you to specify one or more fields, and the corresponding values required for
	// those fields. You can create more complex selector expressions by combining operators.
	//
	// Operators are identified by the use of a dollar sign `$` prefix in the name field.
	//
	// There are two core types of operators in the selector syntax:
	// * Combination operators: applied at the topmost level of selection. They are used to combine selectors. In addition
	// to the common boolean operators (`$and`, `$or`, `$not`, `$nor`) there are three combination operators: `$all`,
	// `$elemMatch`, and `$allMatch`. A combination operator takes a single argument. The argument is either another
	// selector, or an array of selectors.
	// * Condition operators: are specific to a field, and are used to evaluate the value stored in that field. For
	// instance, the basic `$eq` operator matches when the specified field contains a value that is equal to the supplied
	// argument.
	Selector map[string]interface{} `json:"selector" validate:"required"`

	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark,omitempty"`

	// A boolean value that indicates whether or not to include information about existing conflicts in the document.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Use this option to find information about the query that was run. This information includes total key lookups, total
	// document lookups (when `include_docs=true` is used), and total quorum document lookups (when each document replica
	// is fetched).
	ExecutionStats *bool `json:"execution_stats,omitempty"`

	// JSON array that uses the field syntax. Use this parameter to specify which fields of a document must be returned. If
	// it is omitted, the entire document is returned.
	Fields []string `json:"fields,omitempty"`

	// Maximum number of results returned. The `type: text` indexes are limited to 200 results when queried.
	Limit *int64 `json:"limit,omitempty"`

	// Skip the first 'n' results, where 'n' is the value that is specified.
	Skip *int64 `json:"skip,omitempty"`

	// The sort field contains a list of pairs, each mapping a field name to a sort direction (asc or desc). The first
	// field name and direction pair is the topmost level of sort. The second pair, if provided, is the next level of sort.
	// The field can be any field, using dotted notation if desired for sub-document fields.
	//
	// For example in JSON: `[{"fieldName1": "desc"}, {"fieldName2.subFieldName1": "desc"}]`
	//
	// When sorting with multiple fields they must use the same sort direction, either all ascending or all descending.
	Sort []map[string]string `json:"sort,omitempty"`

	// Whether or not the view results should be returned from a "stable" set of shards.
	Stable *bool `json:"stable,omitempty"`

	// Whether to update the index prior to returning the result.
	Update *string `json:"update,omitempty"`

	// Use this option to identify a specific index for query to run against, rather than by using the IBM Cloudant Query
	// algorithm to find the best index.
	UseIndex []string `json:"use_index,omitempty"`

	// The read quorum that is needed for the result. The value defaults to 1, in which case the document that was found in
	// the index is returned. If set to a higher value, each document is read from at least that many replicas before it is
	// returned in the results. The request will take more time than using only the document that is stored locally with
	// the index.
	R *int64 `json:"r,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostExplainOptions.Sort property.
// Schema for a mapping of field name to sort direction.
const (
	PostExplainOptionsSortAscConst = "asc"
	PostExplainOptionsSortDescConst = "desc"
)

// Constants associated with the PostExplainOptions.Update property.
// Whether to update the index prior to returning the result.
const (
	PostExplainOptionsUpdateFalseConst = "false"
	PostExplainOptionsUpdateLazyConst = "lazy"
	PostExplainOptionsUpdateTrueConst = "true"
)

// NewPostExplainOptions : Instantiate PostExplainOptions
func (*CloudantV1) NewPostExplainOptions(db string, selector map[string]interface{}) *PostExplainOptions {
	return &PostExplainOptions{
		Db: core.StringPtr(db),
		Selector: selector,
	}
}

// SetDb : Allow user to set Db
func (_options *PostExplainOptions) SetDb(db string) *PostExplainOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetSelector : Allow user to set Selector
func (_options *PostExplainOptions) SetSelector(selector map[string]interface{}) *PostExplainOptions {
	_options.Selector = selector
	return _options
}

// SetBookmark : Allow user to set Bookmark
func (_options *PostExplainOptions) SetBookmark(bookmark string) *PostExplainOptions {
	_options.Bookmark = core.StringPtr(bookmark)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostExplainOptions) SetConflicts(conflicts bool) *PostExplainOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetExecutionStats : Allow user to set ExecutionStats
func (_options *PostExplainOptions) SetExecutionStats(executionStats bool) *PostExplainOptions {
	_options.ExecutionStats = core.BoolPtr(executionStats)
	return _options
}

// SetFields : Allow user to set Fields
func (_options *PostExplainOptions) SetFields(fields []string) *PostExplainOptions {
	_options.Fields = fields
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostExplainOptions) SetLimit(limit int64) *PostExplainOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *PostExplainOptions) SetSkip(skip int64) *PostExplainOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *PostExplainOptions) SetSort(sort []map[string]string) *PostExplainOptions {
	_options.Sort = sort
	return _options
}

// SetStable : Allow user to set Stable
func (_options *PostExplainOptions) SetStable(stable bool) *PostExplainOptions {
	_options.Stable = core.BoolPtr(stable)
	return _options
}

// SetUpdate : Allow user to set Update
func (_options *PostExplainOptions) SetUpdate(update string) *PostExplainOptions {
	_options.Update = core.StringPtr(update)
	return _options
}

// SetUseIndex : Allow user to set UseIndex
func (_options *PostExplainOptions) SetUseIndex(useIndex []string) *PostExplainOptions {
	_options.UseIndex = useIndex
	return _options
}

// SetR : Allow user to set R
func (_options *PostExplainOptions) SetR(r int64) *PostExplainOptions {
	_options.R = core.Int64Ptr(r)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostExplainOptions) SetHeaders(param map[string]string) *PostExplainOptions {
	options.Headers = param
	return options
}

// PostFindOptions : The PostFind options.
type PostFindOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// JSON object describing criteria used to select documents. The selector specifies fields in the document, and
	// provides an expression to evaluate with the field content or other data.
	//
	// The selector object must:
	//   * Be structured as valid JSON.
	//   * Contain a valid query expression.
	//
	// Using a selector is significantly more efficient than using a JavaScript filter function, and is the recommended
	// option if filtering on document attributes only.
	//
	// Elementary selector syntax requires you to specify one or more fields, and the corresponding values required for
	// those fields. You can create more complex selector expressions by combining operators.
	//
	// Operators are identified by the use of a dollar sign `$` prefix in the name field.
	//
	// There are two core types of operators in the selector syntax:
	// * Combination operators: applied at the topmost level of selection. They are used to combine selectors. In addition
	// to the common boolean operators (`$and`, `$or`, `$not`, `$nor`) there are three combination operators: `$all`,
	// `$elemMatch`, and `$allMatch`. A combination operator takes a single argument. The argument is either another
	// selector, or an array of selectors.
	// * Condition operators: are specific to a field, and are used to evaluate the value stored in that field. For
	// instance, the basic `$eq` operator matches when the specified field contains a value that is equal to the supplied
	// argument.
	Selector map[string]interface{} `json:"selector" validate:"required"`

	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark,omitempty"`

	// A boolean value that indicates whether or not to include information about existing conflicts in the document.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Use this option to find information about the query that was run. This information includes total key lookups, total
	// document lookups (when `include_docs=true` is used), and total quorum document lookups (when each document replica
	// is fetched).
	ExecutionStats *bool `json:"execution_stats,omitempty"`

	// JSON array that uses the field syntax. Use this parameter to specify which fields of a document must be returned. If
	// it is omitted, the entire document is returned.
	Fields []string `json:"fields,omitempty"`

	// Maximum number of results returned. The `type: text` indexes are limited to 200 results when queried.
	Limit *int64 `json:"limit,omitempty"`

	// Skip the first 'n' results, where 'n' is the value that is specified.
	Skip *int64 `json:"skip,omitempty"`

	// The sort field contains a list of pairs, each mapping a field name to a sort direction (asc or desc). The first
	// field name and direction pair is the topmost level of sort. The second pair, if provided, is the next level of sort.
	// The field can be any field, using dotted notation if desired for sub-document fields.
	//
	// For example in JSON: `[{"fieldName1": "desc"}, {"fieldName2.subFieldName1": "desc"}]`
	//
	// When sorting with multiple fields they must use the same sort direction, either all ascending or all descending.
	Sort []map[string]string `json:"sort,omitempty"`

	// Whether or not the view results should be returned from a "stable" set of shards.
	Stable *bool `json:"stable,omitempty"`

	// Whether to update the index prior to returning the result.
	Update *string `json:"update,omitempty"`

	// Use this option to identify a specific index for query to run against, rather than by using the IBM Cloudant Query
	// algorithm to find the best index.
	UseIndex []string `json:"use_index,omitempty"`

	// The read quorum that is needed for the result. The value defaults to 1, in which case the document that was found in
	// the index is returned. If set to a higher value, each document is read from at least that many replicas before it is
	// returned in the results. The request will take more time than using only the document that is stored locally with
	// the index.
	R *int64 `json:"r,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostFindOptions.Sort property.
// Schema for a mapping of field name to sort direction.
const (
	PostFindOptionsSortAscConst = "asc"
	PostFindOptionsSortDescConst = "desc"
)

// Constants associated with the PostFindOptions.Update property.
// Whether to update the index prior to returning the result.
const (
	PostFindOptionsUpdateFalseConst = "false"
	PostFindOptionsUpdateLazyConst = "lazy"
	PostFindOptionsUpdateTrueConst = "true"
)

// NewPostFindOptions : Instantiate PostFindOptions
func (*CloudantV1) NewPostFindOptions(db string, selector map[string]interface{}) *PostFindOptions {
	return &PostFindOptions{
		Db: core.StringPtr(db),
		Selector: selector,
	}
}

// SetDb : Allow user to set Db
func (_options *PostFindOptions) SetDb(db string) *PostFindOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetSelector : Allow user to set Selector
func (_options *PostFindOptions) SetSelector(selector map[string]interface{}) *PostFindOptions {
	_options.Selector = selector
	return _options
}

// SetBookmark : Allow user to set Bookmark
func (_options *PostFindOptions) SetBookmark(bookmark string) *PostFindOptions {
	_options.Bookmark = core.StringPtr(bookmark)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostFindOptions) SetConflicts(conflicts bool) *PostFindOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetExecutionStats : Allow user to set ExecutionStats
func (_options *PostFindOptions) SetExecutionStats(executionStats bool) *PostFindOptions {
	_options.ExecutionStats = core.BoolPtr(executionStats)
	return _options
}

// SetFields : Allow user to set Fields
func (_options *PostFindOptions) SetFields(fields []string) *PostFindOptions {
	_options.Fields = fields
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostFindOptions) SetLimit(limit int64) *PostFindOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *PostFindOptions) SetSkip(skip int64) *PostFindOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *PostFindOptions) SetSort(sort []map[string]string) *PostFindOptions {
	_options.Sort = sort
	return _options
}

// SetStable : Allow user to set Stable
func (_options *PostFindOptions) SetStable(stable bool) *PostFindOptions {
	_options.Stable = core.BoolPtr(stable)
	return _options
}

// SetUpdate : Allow user to set Update
func (_options *PostFindOptions) SetUpdate(update string) *PostFindOptions {
	_options.Update = core.StringPtr(update)
	return _options
}

// SetUseIndex : Allow user to set UseIndex
func (_options *PostFindOptions) SetUseIndex(useIndex []string) *PostFindOptions {
	_options.UseIndex = useIndex
	return _options
}

// SetR : Allow user to set R
func (_options *PostFindOptions) SetR(r int64) *PostFindOptions {
	_options.R = core.Int64Ptr(r)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostFindOptions) SetHeaders(param map[string]string) *PostFindOptions {
	options.Headers = param
	return options
}

// PostGeoCleanupOptions : The PostGeoCleanup options.
type PostGeoCleanupOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostGeoCleanupOptions : Instantiate PostGeoCleanupOptions
func (*CloudantV1) NewPostGeoCleanupOptions(db string) *PostGeoCleanupOptions {
	return &PostGeoCleanupOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *PostGeoCleanupOptions) SetDb(db string) *PostGeoCleanupOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostGeoCleanupOptions) SetHeaders(param map[string]string) *PostGeoCleanupOptions {
	options.Headers = param
	return options
}

// PostIndexOptions : The PostIndex options.
type PostIndexOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Schema for a `json` or `text` query index definition. Indexes of type `text` have additional configuration
	// properties that do not apply to `json` indexes, these are:
	// * `default_analyzer` - the default text analyzer to use * `default_field` - whether to index the text in all
	// document fields and what analyzer to use for that purpose.
	Index *IndexDefinition `json:"index" validate:"required"`

	// Name of the design document in which the index will be created.
	Ddoc *string `json:"ddoc,omitempty"`

	// Schema for a `json` or `text` query index definition. Indexes of type `text` have additional configuration
	// properties that do not apply to `json` indexes, these are:
	// * `default_analyzer` - the default text analyzer to use * `default_field` - whether to index the text in all
	// document fields and what analyzer to use for that purpose.
	Def *IndexDefinition `json:"def,omitempty"`

	// name.
	Name *string `json:"name,omitempty"`

	// The default value is `true` for databases with `partitioned: true` and `false` otherwise. For databases with
	// `partitioned: false` if this option is specified the value must be `false`.
	Partitioned *bool `json:"partitioned,omitempty"`

	// Schema for the type of an index.
	Type *string `json:"type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostIndexOptions.Type property.
// Schema for the type of an index.
const (
	PostIndexOptionsTypeJSONConst = "json"
	PostIndexOptionsTypeSpecialConst = "special"
	PostIndexOptionsTypeTextConst = "text"
)

// NewPostIndexOptions : Instantiate PostIndexOptions
func (*CloudantV1) NewPostIndexOptions(db string, index *IndexDefinition) *PostIndexOptions {
	return &PostIndexOptions{
		Db: core.StringPtr(db),
		Index: index,
	}
}

// SetDb : Allow user to set Db
func (_options *PostIndexOptions) SetDb(db string) *PostIndexOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetIndex : Allow user to set Index
func (_options *PostIndexOptions) SetIndex(index *IndexDefinition) *PostIndexOptions {
	_options.Index = index
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *PostIndexOptions) SetDdoc(ddoc string) *PostIndexOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetDef : Allow user to set Def
func (_options *PostIndexOptions) SetDef(def *IndexDefinition) *PostIndexOptions {
	_options.Def = def
	return _options
}

// SetName : Allow user to set Name
func (_options *PostIndexOptions) SetName(name string) *PostIndexOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPartitioned : Allow user to set Partitioned
func (_options *PostIndexOptions) SetPartitioned(partitioned bool) *PostIndexOptions {
	_options.Partitioned = core.BoolPtr(partitioned)
	return _options
}

// SetType : Allow user to set Type
func (_options *PostIndexOptions) SetType(typeVar string) *PostIndexOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostIndexOptions) SetHeaders(param map[string]string) *PostIndexOptions {
	options.Headers = param
	return options
}

// PostPartitionAllDocsOptions : The PostPartitionAllDocs options.
type PostPartitionAllDocsOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the database partition key.
	PartitionKey *string `json:"partition_key" validate:"required,ne="`

	// Parameter to specify whether to include the encoding information in attachment stubs if the particular attachment is
	// compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the returned
	// document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Parameter to specify whether the specified end key should be included in the result.
	InclusiveEnd *bool `json:"inclusive_end,omitempty"`

	// Parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Parameter to specify whether to include in the response an update_seq value indicating the sequence id of the
	// database the view reflects.
	UpdateSeq *bool `json:"update_seq,omitempty"`

	// Schema for a document ID.
	Endkey *string `json:"endkey,omitempty"`

	// Schema for a document ID.
	Key *string `json:"key,omitempty"`

	// Schema for a list of document IDs.
	Keys []string `json:"keys,omitempty"`

	// Schema for a document ID.
	Startkey *string `json:"startkey,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostPartitionAllDocsOptions : Instantiate PostPartitionAllDocsOptions
func (*CloudantV1) NewPostPartitionAllDocsOptions(db string, partitionKey string) *PostPartitionAllDocsOptions {
	return &PostPartitionAllDocsOptions{
		Db: core.StringPtr(db),
		PartitionKey: core.StringPtr(partitionKey),
	}
}

// SetDb : Allow user to set Db
func (_options *PostPartitionAllDocsOptions) SetDb(db string) *PostPartitionAllDocsOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetPartitionKey : Allow user to set PartitionKey
func (_options *PostPartitionAllDocsOptions) SetPartitionKey(partitionKey string) *PostPartitionAllDocsOptions {
	_options.PartitionKey = core.StringPtr(partitionKey)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *PostPartitionAllDocsOptions) SetAttEncodingInfo(attEncodingInfo bool) *PostPartitionAllDocsOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *PostPartitionAllDocsOptions) SetAttachments(attachments bool) *PostPartitionAllDocsOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostPartitionAllDocsOptions) SetConflicts(conflicts bool) *PostPartitionAllDocsOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDescending : Allow user to set Descending
func (_options *PostPartitionAllDocsOptions) SetDescending(descending bool) *PostPartitionAllDocsOptions {
	_options.Descending = core.BoolPtr(descending)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *PostPartitionAllDocsOptions) SetIncludeDocs(includeDocs bool) *PostPartitionAllDocsOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetInclusiveEnd : Allow user to set InclusiveEnd
func (_options *PostPartitionAllDocsOptions) SetInclusiveEnd(inclusiveEnd bool) *PostPartitionAllDocsOptions {
	_options.InclusiveEnd = core.BoolPtr(inclusiveEnd)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostPartitionAllDocsOptions) SetLimit(limit int64) *PostPartitionAllDocsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *PostPartitionAllDocsOptions) SetSkip(skip int64) *PostPartitionAllDocsOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetUpdateSeq : Allow user to set UpdateSeq
func (_options *PostPartitionAllDocsOptions) SetUpdateSeq(updateSeq bool) *PostPartitionAllDocsOptions {
	_options.UpdateSeq = core.BoolPtr(updateSeq)
	return _options
}

// SetEndkey : Allow user to set Endkey
func (_options *PostPartitionAllDocsOptions) SetEndkey(endkey string) *PostPartitionAllDocsOptions {
	_options.Endkey = core.StringPtr(endkey)
	return _options
}

// SetKey : Allow user to set Key
func (_options *PostPartitionAllDocsOptions) SetKey(key string) *PostPartitionAllDocsOptions {
	_options.Key = core.StringPtr(key)
	return _options
}

// SetKeys : Allow user to set Keys
func (_options *PostPartitionAllDocsOptions) SetKeys(keys []string) *PostPartitionAllDocsOptions {
	_options.Keys = keys
	return _options
}

// SetStartkey : Allow user to set Startkey
func (_options *PostPartitionAllDocsOptions) SetStartkey(startkey string) *PostPartitionAllDocsOptions {
	_options.Startkey = core.StringPtr(startkey)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostPartitionAllDocsOptions) SetHeaders(param map[string]string) *PostPartitionAllDocsOptions {
	options.Headers = param
	return options
}

// PostPartitionFindOptions : The PostPartitionFind options.
type PostPartitionFindOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the database partition key.
	PartitionKey *string `json:"partition_key" validate:"required,ne="`

	// JSON object describing criteria used to select documents. The selector specifies fields in the document, and
	// provides an expression to evaluate with the field content or other data.
	//
	// The selector object must:
	//   * Be structured as valid JSON.
	//   * Contain a valid query expression.
	//
	// Using a selector is significantly more efficient than using a JavaScript filter function, and is the recommended
	// option if filtering on document attributes only.
	//
	// Elementary selector syntax requires you to specify one or more fields, and the corresponding values required for
	// those fields. You can create more complex selector expressions by combining operators.
	//
	// Operators are identified by the use of a dollar sign `$` prefix in the name field.
	//
	// There are two core types of operators in the selector syntax:
	// * Combination operators: applied at the topmost level of selection. They are used to combine selectors. In addition
	// to the common boolean operators (`$and`, `$or`, `$not`, `$nor`) there are three combination operators: `$all`,
	// `$elemMatch`, and `$allMatch`. A combination operator takes a single argument. The argument is either another
	// selector, or an array of selectors.
	// * Condition operators: are specific to a field, and are used to evaluate the value stored in that field. For
	// instance, the basic `$eq` operator matches when the specified field contains a value that is equal to the supplied
	// argument.
	Selector map[string]interface{} `json:"selector" validate:"required"`

	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark,omitempty"`

	// A boolean value that indicates whether or not to include information about existing conflicts in the document.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Use this option to find information about the query that was run. This information includes total key lookups, total
	// document lookups (when `include_docs=true` is used), and total quorum document lookups (when each document replica
	// is fetched).
	ExecutionStats *bool `json:"execution_stats,omitempty"`

	// JSON array that uses the field syntax. Use this parameter to specify which fields of a document must be returned. If
	// it is omitted, the entire document is returned.
	Fields []string `json:"fields,omitempty"`

	// Maximum number of results returned. The `type: text` indexes are limited to 200 results when queried.
	Limit *int64 `json:"limit,omitempty"`

	// Skip the first 'n' results, where 'n' is the value that is specified.
	Skip *int64 `json:"skip,omitempty"`

	// The sort field contains a list of pairs, each mapping a field name to a sort direction (asc or desc). The first
	// field name and direction pair is the topmost level of sort. The second pair, if provided, is the next level of sort.
	// The field can be any field, using dotted notation if desired for sub-document fields.
	//
	// For example in JSON: `[{"fieldName1": "desc"}, {"fieldName2.subFieldName1": "desc"}]`
	//
	// When sorting with multiple fields they must use the same sort direction, either all ascending or all descending.
	Sort []map[string]string `json:"sort,omitempty"`

	// Whether or not the view results should be returned from a "stable" set of shards.
	Stable *bool `json:"stable,omitempty"`

	// Whether to update the index prior to returning the result.
	Update *string `json:"update,omitempty"`

	// Use this option to identify a specific index for query to run against, rather than by using the IBM Cloudant Query
	// algorithm to find the best index.
	UseIndex []string `json:"use_index,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostPartitionFindOptions.Sort property.
// Schema for a mapping of field name to sort direction.
const (
	PostPartitionFindOptionsSortAscConst = "asc"
	PostPartitionFindOptionsSortDescConst = "desc"
)

// Constants associated with the PostPartitionFindOptions.Update property.
// Whether to update the index prior to returning the result.
const (
	PostPartitionFindOptionsUpdateFalseConst = "false"
	PostPartitionFindOptionsUpdateLazyConst = "lazy"
	PostPartitionFindOptionsUpdateTrueConst = "true"
)

// NewPostPartitionFindOptions : Instantiate PostPartitionFindOptions
func (*CloudantV1) NewPostPartitionFindOptions(db string, partitionKey string, selector map[string]interface{}) *PostPartitionFindOptions {
	return &PostPartitionFindOptions{
		Db: core.StringPtr(db),
		PartitionKey: core.StringPtr(partitionKey),
		Selector: selector,
	}
}

// SetDb : Allow user to set Db
func (_options *PostPartitionFindOptions) SetDb(db string) *PostPartitionFindOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetPartitionKey : Allow user to set PartitionKey
func (_options *PostPartitionFindOptions) SetPartitionKey(partitionKey string) *PostPartitionFindOptions {
	_options.PartitionKey = core.StringPtr(partitionKey)
	return _options
}

// SetSelector : Allow user to set Selector
func (_options *PostPartitionFindOptions) SetSelector(selector map[string]interface{}) *PostPartitionFindOptions {
	_options.Selector = selector
	return _options
}

// SetBookmark : Allow user to set Bookmark
func (_options *PostPartitionFindOptions) SetBookmark(bookmark string) *PostPartitionFindOptions {
	_options.Bookmark = core.StringPtr(bookmark)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostPartitionFindOptions) SetConflicts(conflicts bool) *PostPartitionFindOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetExecutionStats : Allow user to set ExecutionStats
func (_options *PostPartitionFindOptions) SetExecutionStats(executionStats bool) *PostPartitionFindOptions {
	_options.ExecutionStats = core.BoolPtr(executionStats)
	return _options
}

// SetFields : Allow user to set Fields
func (_options *PostPartitionFindOptions) SetFields(fields []string) *PostPartitionFindOptions {
	_options.Fields = fields
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostPartitionFindOptions) SetLimit(limit int64) *PostPartitionFindOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *PostPartitionFindOptions) SetSkip(skip int64) *PostPartitionFindOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *PostPartitionFindOptions) SetSort(sort []map[string]string) *PostPartitionFindOptions {
	_options.Sort = sort
	return _options
}

// SetStable : Allow user to set Stable
func (_options *PostPartitionFindOptions) SetStable(stable bool) *PostPartitionFindOptions {
	_options.Stable = core.BoolPtr(stable)
	return _options
}

// SetUpdate : Allow user to set Update
func (_options *PostPartitionFindOptions) SetUpdate(update string) *PostPartitionFindOptions {
	_options.Update = core.StringPtr(update)
	return _options
}

// SetUseIndex : Allow user to set UseIndex
func (_options *PostPartitionFindOptions) SetUseIndex(useIndex []string) *PostPartitionFindOptions {
	_options.UseIndex = useIndex
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostPartitionFindOptions) SetHeaders(param map[string]string) *PostPartitionFindOptions {
	options.Headers = param
	return options
}

// PostPartitionSearchOptions : The PostPartitionSearch options.
type PostPartitionSearchOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the database partition key.
	PartitionKey *string `json:"partition_key" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the index name.
	Index *string `json:"index" validate:"required,ne="`

	// The Lucene query to execute.
	Query *string `json:"query" validate:"required"`

	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark,omitempty"`

	// Specifies which fields to highlight. If specified, the result object contains a highlights field with an entry for
	// each specified field.
	HighlightFields []string `json:"highlight_fields,omitempty"`

	// Number of fragments that are returned in highlights. If the search term occurs less often than the number of
	// fragments that are specified, longer fragments are returned.
	HighlightNumber *int64 `json:"highlight_number,omitempty"`

	// A string that is inserted after the highlighted word in the highlights output.
	HighlightPostTag *string `json:"highlight_post_tag,omitempty"`

	// A string that is inserted before the highlighted word in the highlights output.
	HighlightPreTag *string `json:"highlight_pre_tag,omitempty"`

	// Number of characters in each fragment for highlights.
	HighlightSize *int64 `json:"highlight_size,omitempty"`

	// Include the full content of the documents in the return.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// A JSON array of field names to include in search results. Any fields that are included must be indexed with the
	// store:true option. The default is all fields.
	IncludeFields []string `json:"include_fields,omitempty"`

	// Limit the number of the returned documents to the specified number.
	Limit *int64 `json:"limit,omitempty"`

	// Specifies the sort order of the results. In a grouped search (when group_field is used), this parameter specifies
	// the sort order within a group. The default sort order is relevance.  A JSON string of the form
	// "fieldname&lt;type&gt;" or "-fieldname&lt;type&gt;" for descending order, where fieldname is the name of a string or
	// number field, and type is either a number, a string, or a JSON array of strings. The type part is optional, and
	// defaults to number. Some examples are "foo", "-foo", "bar&lt;string&gt;", "-foo&lt;number&gt;" and
	// ["-foo&lt;number&gt;", "bar&lt;string&gt;"]. String fields that are used for sorting must not be analyzed fields.
	// Fields that are used for sorting must be indexed by the same indexer that is used for the search query.
	Sort []string `json:"sort,omitempty"`

	// Do not wait for the index to finish building to return results.
	Stale *string `json:"stale,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostPartitionSearchOptions.Stale property.
// Do not wait for the index to finish building to return results.
const (
	PostPartitionSearchOptionsStaleOkConst = "ok"
)

// NewPostPartitionSearchOptions : Instantiate PostPartitionSearchOptions
func (*CloudantV1) NewPostPartitionSearchOptions(db string, partitionKey string, ddoc string, index string, query string) *PostPartitionSearchOptions {
	return &PostPartitionSearchOptions{
		Db: core.StringPtr(db),
		PartitionKey: core.StringPtr(partitionKey),
		Ddoc: core.StringPtr(ddoc),
		Index: core.StringPtr(index),
		Query: core.StringPtr(query),
	}
}

// SetDb : Allow user to set Db
func (_options *PostPartitionSearchOptions) SetDb(db string) *PostPartitionSearchOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetPartitionKey : Allow user to set PartitionKey
func (_options *PostPartitionSearchOptions) SetPartitionKey(partitionKey string) *PostPartitionSearchOptions {
	_options.PartitionKey = core.StringPtr(partitionKey)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *PostPartitionSearchOptions) SetDdoc(ddoc string) *PostPartitionSearchOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetIndex : Allow user to set Index
func (_options *PostPartitionSearchOptions) SetIndex(index string) *PostPartitionSearchOptions {
	_options.Index = core.StringPtr(index)
	return _options
}

// SetQuery : Allow user to set Query
func (_options *PostPartitionSearchOptions) SetQuery(query string) *PostPartitionSearchOptions {
	_options.Query = core.StringPtr(query)
	return _options
}

// SetBookmark : Allow user to set Bookmark
func (_options *PostPartitionSearchOptions) SetBookmark(bookmark string) *PostPartitionSearchOptions {
	_options.Bookmark = core.StringPtr(bookmark)
	return _options
}

// SetHighlightFields : Allow user to set HighlightFields
func (_options *PostPartitionSearchOptions) SetHighlightFields(highlightFields []string) *PostPartitionSearchOptions {
	_options.HighlightFields = highlightFields
	return _options
}

// SetHighlightNumber : Allow user to set HighlightNumber
func (_options *PostPartitionSearchOptions) SetHighlightNumber(highlightNumber int64) *PostPartitionSearchOptions {
	_options.HighlightNumber = core.Int64Ptr(highlightNumber)
	return _options
}

// SetHighlightPostTag : Allow user to set HighlightPostTag
func (_options *PostPartitionSearchOptions) SetHighlightPostTag(highlightPostTag string) *PostPartitionSearchOptions {
	_options.HighlightPostTag = core.StringPtr(highlightPostTag)
	return _options
}

// SetHighlightPreTag : Allow user to set HighlightPreTag
func (_options *PostPartitionSearchOptions) SetHighlightPreTag(highlightPreTag string) *PostPartitionSearchOptions {
	_options.HighlightPreTag = core.StringPtr(highlightPreTag)
	return _options
}

// SetHighlightSize : Allow user to set HighlightSize
func (_options *PostPartitionSearchOptions) SetHighlightSize(highlightSize int64) *PostPartitionSearchOptions {
	_options.HighlightSize = core.Int64Ptr(highlightSize)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *PostPartitionSearchOptions) SetIncludeDocs(includeDocs bool) *PostPartitionSearchOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetIncludeFields : Allow user to set IncludeFields
func (_options *PostPartitionSearchOptions) SetIncludeFields(includeFields []string) *PostPartitionSearchOptions {
	_options.IncludeFields = includeFields
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostPartitionSearchOptions) SetLimit(limit int64) *PostPartitionSearchOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *PostPartitionSearchOptions) SetSort(sort []string) *PostPartitionSearchOptions {
	_options.Sort = sort
	return _options
}

// SetStale : Allow user to set Stale
func (_options *PostPartitionSearchOptions) SetStale(stale string) *PostPartitionSearchOptions {
	_options.Stale = core.StringPtr(stale)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostPartitionSearchOptions) SetHeaders(param map[string]string) *PostPartitionSearchOptions {
	options.Headers = param
	return options
}

// PostPartitionViewOptions : The PostPartitionView options.
type PostPartitionViewOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the database partition key.
	PartitionKey *string `json:"partition_key" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the map reduce view function name.
	View *string `json:"view" validate:"required,ne="`

	// Parameter to specify whether to include the encoding information in attachment stubs if the particular attachment is
	// compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the returned
	// document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Parameter to specify whether the specified end key should be included in the result.
	InclusiveEnd *bool `json:"inclusive_end,omitempty"`

	// Parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Parameter to specify whether to include in the response an update_seq value indicating the sequence id of the
	// database the view reflects.
	UpdateSeq *bool `json:"update_seq,omitempty"`

	// Schema for any JSON type.
	Endkey interface{} `json:"endkey,omitempty"`

	// Schema for a document ID.
	EndkeyDocid *string `json:"endkey_docid,omitempty"`

	// Parameter to specify whether to group the results using the reduce function to a group rather than a single row.
	// Implies reduce is true and the maximum group_level.
	Group *bool `json:"group,omitempty"`

	// Parameter to specify the group level to be used. Implies group is true.
	GroupLevel *int64 `json:"group_level,omitempty"`

	// Schema for any JSON type.
	Key interface{} `json:"key,omitempty"`

	// Parameter to specify to return only documents that match the specified keys. String representation of a JSON array
	// containing elements that match the key type emitted by the view function.
	Keys []interface{} `json:"keys,omitempty"`

	// Parameter to specify whether to use the reduce function in a map-reduce view. Default is true when a reduce function
	// is defined.
	Reduce *bool `json:"reduce,omitempty"`

	// Parameter to specify whether view results should be returned from a stable set of shards.
	Stable *bool `json:"stable,omitempty"`

	// Schema for any JSON type.
	Startkey interface{} `json:"startkey,omitempty"`

	// Schema for a document ID.
	StartkeyDocid *string `json:"startkey_docid,omitempty"`

	// Parameter to specify whether or not the view in question should be updated prior to responding to the user.
	Update *string `json:"update,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostPartitionViewOptions.Update property.
// Parameter to specify whether or not the view in question should be updated prior to responding to the user.
const (
	PostPartitionViewOptionsUpdateFalseConst = "false"
	PostPartitionViewOptionsUpdateLazyConst = "lazy"
	PostPartitionViewOptionsUpdateTrueConst = "true"
)

// NewPostPartitionViewOptions : Instantiate PostPartitionViewOptions
func (*CloudantV1) NewPostPartitionViewOptions(db string, partitionKey string, ddoc string, view string) *PostPartitionViewOptions {
	return &PostPartitionViewOptions{
		Db: core.StringPtr(db),
		PartitionKey: core.StringPtr(partitionKey),
		Ddoc: core.StringPtr(ddoc),
		View: core.StringPtr(view),
	}
}

// SetDb : Allow user to set Db
func (_options *PostPartitionViewOptions) SetDb(db string) *PostPartitionViewOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetPartitionKey : Allow user to set PartitionKey
func (_options *PostPartitionViewOptions) SetPartitionKey(partitionKey string) *PostPartitionViewOptions {
	_options.PartitionKey = core.StringPtr(partitionKey)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *PostPartitionViewOptions) SetDdoc(ddoc string) *PostPartitionViewOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetView : Allow user to set View
func (_options *PostPartitionViewOptions) SetView(view string) *PostPartitionViewOptions {
	_options.View = core.StringPtr(view)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *PostPartitionViewOptions) SetAttEncodingInfo(attEncodingInfo bool) *PostPartitionViewOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *PostPartitionViewOptions) SetAttachments(attachments bool) *PostPartitionViewOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostPartitionViewOptions) SetConflicts(conflicts bool) *PostPartitionViewOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDescending : Allow user to set Descending
func (_options *PostPartitionViewOptions) SetDescending(descending bool) *PostPartitionViewOptions {
	_options.Descending = core.BoolPtr(descending)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *PostPartitionViewOptions) SetIncludeDocs(includeDocs bool) *PostPartitionViewOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetInclusiveEnd : Allow user to set InclusiveEnd
func (_options *PostPartitionViewOptions) SetInclusiveEnd(inclusiveEnd bool) *PostPartitionViewOptions {
	_options.InclusiveEnd = core.BoolPtr(inclusiveEnd)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostPartitionViewOptions) SetLimit(limit int64) *PostPartitionViewOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *PostPartitionViewOptions) SetSkip(skip int64) *PostPartitionViewOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetUpdateSeq : Allow user to set UpdateSeq
func (_options *PostPartitionViewOptions) SetUpdateSeq(updateSeq bool) *PostPartitionViewOptions {
	_options.UpdateSeq = core.BoolPtr(updateSeq)
	return _options
}

// SetEndkey : Allow user to set Endkey
func (_options *PostPartitionViewOptions) SetEndkey(endkey interface{}) *PostPartitionViewOptions {
	_options.Endkey = endkey
	return _options
}

// SetEndkeyDocid : Allow user to set EndkeyDocid
func (_options *PostPartitionViewOptions) SetEndkeyDocid(endkeyDocid string) *PostPartitionViewOptions {
	_options.EndkeyDocid = core.StringPtr(endkeyDocid)
	return _options
}

// SetGroup : Allow user to set Group
func (_options *PostPartitionViewOptions) SetGroup(group bool) *PostPartitionViewOptions {
	_options.Group = core.BoolPtr(group)
	return _options
}

// SetGroupLevel : Allow user to set GroupLevel
func (_options *PostPartitionViewOptions) SetGroupLevel(groupLevel int64) *PostPartitionViewOptions {
	_options.GroupLevel = core.Int64Ptr(groupLevel)
	return _options
}

// SetKey : Allow user to set Key
func (_options *PostPartitionViewOptions) SetKey(key interface{}) *PostPartitionViewOptions {
	_options.Key = key
	return _options
}

// SetKeys : Allow user to set Keys
func (_options *PostPartitionViewOptions) SetKeys(keys []interface{}) *PostPartitionViewOptions {
	_options.Keys = keys
	return _options
}

// SetReduce : Allow user to set Reduce
func (_options *PostPartitionViewOptions) SetReduce(reduce bool) *PostPartitionViewOptions {
	_options.Reduce = core.BoolPtr(reduce)
	return _options
}

// SetStable : Allow user to set Stable
func (_options *PostPartitionViewOptions) SetStable(stable bool) *PostPartitionViewOptions {
	_options.Stable = core.BoolPtr(stable)
	return _options
}

// SetStartkey : Allow user to set Startkey
func (_options *PostPartitionViewOptions) SetStartkey(startkey interface{}) *PostPartitionViewOptions {
	_options.Startkey = startkey
	return _options
}

// SetStartkeyDocid : Allow user to set StartkeyDocid
func (_options *PostPartitionViewOptions) SetStartkeyDocid(startkeyDocid string) *PostPartitionViewOptions {
	_options.StartkeyDocid = core.StringPtr(startkeyDocid)
	return _options
}

// SetUpdate : Allow user to set Update
func (_options *PostPartitionViewOptions) SetUpdate(update string) *PostPartitionViewOptions {
	_options.Update = core.StringPtr(update)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostPartitionViewOptions) SetHeaders(param map[string]string) *PostPartitionViewOptions {
	options.Headers = param
	return options
}

// PostRevsDiffOptions : The PostRevsDiff options.
type PostRevsDiffOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// HTTP request body for operations with Document revisions.
	DocumentRevisions map[string][]string `json:"documentRevisions" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostRevsDiffOptions : Instantiate PostRevsDiffOptions
func (*CloudantV1) NewPostRevsDiffOptions(db string, documentRevisions map[string][]string) *PostRevsDiffOptions {
	return &PostRevsDiffOptions{
		Db: core.StringPtr(db),
		DocumentRevisions: documentRevisions,
	}
}

// SetDb : Allow user to set Db
func (_options *PostRevsDiffOptions) SetDb(db string) *PostRevsDiffOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocumentRevisions : Allow user to set DocumentRevisions
func (_options *PostRevsDiffOptions) SetDocumentRevisions(documentRevisions map[string][]string) *PostRevsDiffOptions {
	_options.DocumentRevisions = documentRevisions
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostRevsDiffOptions) SetHeaders(param map[string]string) *PostRevsDiffOptions {
	options.Headers = param
	return options
}

// PostSearchAnalyzeOptions : The PostSearchAnalyze options.
type PostSearchAnalyzeOptions struct {
	// The analyzer type that is being used at the tokenization.
	Analyzer *string `json:"analyzer" validate:"required"`

	// The text to tokenize with the analyzer.
	Text *string `json:"text" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostSearchAnalyzeOptions.Analyzer property.
// The analyzer type that is being used at the tokenization.
const (
	PostSearchAnalyzeOptionsAnalyzerArabicConst = "arabic"
	PostSearchAnalyzeOptionsAnalyzerArmenianConst = "armenian"
	PostSearchAnalyzeOptionsAnalyzerBasqueConst = "basque"
	PostSearchAnalyzeOptionsAnalyzerBrazilianConst = "brazilian"
	PostSearchAnalyzeOptionsAnalyzerBulgarianConst = "bulgarian"
	PostSearchAnalyzeOptionsAnalyzerCatalanConst = "catalan"
	PostSearchAnalyzeOptionsAnalyzerChineseConst = "chinese"
	PostSearchAnalyzeOptionsAnalyzerCjkConst = "cjk"
	PostSearchAnalyzeOptionsAnalyzerClassicConst = "classic"
	PostSearchAnalyzeOptionsAnalyzerCzechConst = "czech"
	PostSearchAnalyzeOptionsAnalyzerDanishConst = "danish"
	PostSearchAnalyzeOptionsAnalyzerDutchConst = "dutch"
	PostSearchAnalyzeOptionsAnalyzerEmailConst = "email"
	PostSearchAnalyzeOptionsAnalyzerEnglishConst = "english"
	PostSearchAnalyzeOptionsAnalyzerFinnishConst = "finnish"
	PostSearchAnalyzeOptionsAnalyzerFrenchConst = "french"
	PostSearchAnalyzeOptionsAnalyzerGalicianConst = "galician"
	PostSearchAnalyzeOptionsAnalyzerGermanConst = "german"
	PostSearchAnalyzeOptionsAnalyzerGreekConst = "greek"
	PostSearchAnalyzeOptionsAnalyzerHindiConst = "hindi"
	PostSearchAnalyzeOptionsAnalyzerHungarianConst = "hungarian"
	PostSearchAnalyzeOptionsAnalyzerIndonesianConst = "indonesian"
	PostSearchAnalyzeOptionsAnalyzerIrishConst = "irish"
	PostSearchAnalyzeOptionsAnalyzerItalianConst = "italian"
	PostSearchAnalyzeOptionsAnalyzerJapaneseConst = "japanese"
	PostSearchAnalyzeOptionsAnalyzerKeywordConst = "keyword"
	PostSearchAnalyzeOptionsAnalyzerLatvianConst = "latvian"
	PostSearchAnalyzeOptionsAnalyzerNorwegianConst = "norwegian"
	PostSearchAnalyzeOptionsAnalyzerPersianConst = "persian"
	PostSearchAnalyzeOptionsAnalyzerPolishConst = "polish"
	PostSearchAnalyzeOptionsAnalyzerPortugueseConst = "portuguese"
	PostSearchAnalyzeOptionsAnalyzerRomanianConst = "romanian"
	PostSearchAnalyzeOptionsAnalyzerRussianConst = "russian"
	PostSearchAnalyzeOptionsAnalyzerSimpleConst = "simple"
	PostSearchAnalyzeOptionsAnalyzerSpanishConst = "spanish"
	PostSearchAnalyzeOptionsAnalyzerStandardConst = "standard"
	PostSearchAnalyzeOptionsAnalyzerSwedishConst = "swedish"
	PostSearchAnalyzeOptionsAnalyzerThaiConst = "thai"
	PostSearchAnalyzeOptionsAnalyzerTurkishConst = "turkish"
	PostSearchAnalyzeOptionsAnalyzerWhitespaceConst = "whitespace"
)

// NewPostSearchAnalyzeOptions : Instantiate PostSearchAnalyzeOptions
func (*CloudantV1) NewPostSearchAnalyzeOptions(analyzer string, text string) *PostSearchAnalyzeOptions {
	return &PostSearchAnalyzeOptions{
		Analyzer: core.StringPtr(analyzer),
		Text: core.StringPtr(text),
	}
}

// SetAnalyzer : Allow user to set Analyzer
func (_options *PostSearchAnalyzeOptions) SetAnalyzer(analyzer string) *PostSearchAnalyzeOptions {
	_options.Analyzer = core.StringPtr(analyzer)
	return _options
}

// SetText : Allow user to set Text
func (_options *PostSearchAnalyzeOptions) SetText(text string) *PostSearchAnalyzeOptions {
	_options.Text = core.StringPtr(text)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostSearchAnalyzeOptions) SetHeaders(param map[string]string) *PostSearchAnalyzeOptions {
	options.Headers = param
	return options
}

// PostSearchOptions : The PostSearch options.
type PostSearchOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the index name.
	Index *string `json:"index" validate:"required,ne="`

	// The Lucene query to execute.
	Query *string `json:"query" validate:"required"`

	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark,omitempty"`

	// Specifies which fields to highlight. If specified, the result object contains a highlights field with an entry for
	// each specified field.
	HighlightFields []string `json:"highlight_fields,omitempty"`

	// Number of fragments that are returned in highlights. If the search term occurs less often than the number of
	// fragments that are specified, longer fragments are returned.
	HighlightNumber *int64 `json:"highlight_number,omitempty"`

	// A string that is inserted after the highlighted word in the highlights output.
	HighlightPostTag *string `json:"highlight_post_tag,omitempty"`

	// A string that is inserted before the highlighted word in the highlights output.
	HighlightPreTag *string `json:"highlight_pre_tag,omitempty"`

	// Number of characters in each fragment for highlights.
	HighlightSize *int64 `json:"highlight_size,omitempty"`

	// Include the full content of the documents in the return.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// A JSON array of field names to include in search results. Any fields that are included must be indexed with the
	// store:true option. The default is all fields.
	IncludeFields []string `json:"include_fields,omitempty"`

	// Limit the number of the returned documents to the specified number.
	Limit *int64 `json:"limit,omitempty"`

	// Specifies the sort order of the results. In a grouped search (when group_field is used), this parameter specifies
	// the sort order within a group. The default sort order is relevance.  A JSON string of the form
	// "fieldname&lt;type&gt;" or "-fieldname&lt;type&gt;" for descending order, where fieldname is the name of a string or
	// number field, and type is either a number, a string, or a JSON array of strings. The type part is optional, and
	// defaults to number. Some examples are "foo", "-foo", "bar&lt;string&gt;", "-foo&lt;number&gt;" and
	// ["-foo&lt;number&gt;", "bar&lt;string&gt;"]. String fields that are used for sorting must not be analyzed fields.
	// Fields that are used for sorting must be indexed by the same indexer that is used for the search query.
	Sort []string `json:"sort,omitempty"`

	// Do not wait for the index to finish building to return results.
	Stale *string `json:"stale,omitempty"`

	// This field defines an array of names of string fields, for which counts are requested. The response contains counts
	// for each unique value of this field name among the documents that match the search query. Faceting must be enabled
	// for this parameter to function. This option is only available when making global queries.
	Counts []string `json:"counts,omitempty"`

	// Restrict results to documents with a dimension equal to the specified label(s). The search matches only documents
	// containing the value that was provided in the named field. It differs from using "fieldname:value" in the q
	// parameter only in that the values are not analyzed. Faceting must be enabled for this parameter to function.
	Drilldown [][]string `json:"drilldown,omitempty"`

	// Field by which to group search matches. A string that contains the name of a string field. Fields containing other
	// data such as numbers, objects, or arrays cannot be used. This option is only available when making global queries.
	GroupField *string `json:"group_field,omitempty"`

	// Maximum group count. This field can be used only if group_field is specified. This option is only available when
	// making global queries.
	GroupLimit *int64 `json:"group_limit,omitempty"`

	// This field defines the order of the groups in a search that uses group_field. The default sort order is relevance.
	// This field can have the same values as the sort field, so single fields and arrays of fields are supported. This
	// option is only available when making global queries.
	GroupSort []string `json:"group_sort,omitempty"`

	// This field defines ranges for faceted, numeric search fields. The value is a JSON object where the fields names are
	// faceted numeric search fields, and the values of the fields are JSON objects. The field names of the JSON objects
	// are names for ranges. The values are strings that describe the range, for example "[0 TO 10]". This option is only
	// available when making global queries.
	Ranges map[string]map[string]map[string]string `json:"ranges,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostSearchOptions.Stale property.
// Do not wait for the index to finish building to return results.
const (
	PostSearchOptionsStaleOkConst = "ok"
)

// NewPostSearchOptions : Instantiate PostSearchOptions
func (*CloudantV1) NewPostSearchOptions(db string, ddoc string, index string, query string) *PostSearchOptions {
	return &PostSearchOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
		Index: core.StringPtr(index),
		Query: core.StringPtr(query),
	}
}

// SetDb : Allow user to set Db
func (_options *PostSearchOptions) SetDb(db string) *PostSearchOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *PostSearchOptions) SetDdoc(ddoc string) *PostSearchOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetIndex : Allow user to set Index
func (_options *PostSearchOptions) SetIndex(index string) *PostSearchOptions {
	_options.Index = core.StringPtr(index)
	return _options
}

// SetQuery : Allow user to set Query
func (_options *PostSearchOptions) SetQuery(query string) *PostSearchOptions {
	_options.Query = core.StringPtr(query)
	return _options
}

// SetBookmark : Allow user to set Bookmark
func (_options *PostSearchOptions) SetBookmark(bookmark string) *PostSearchOptions {
	_options.Bookmark = core.StringPtr(bookmark)
	return _options
}

// SetHighlightFields : Allow user to set HighlightFields
func (_options *PostSearchOptions) SetHighlightFields(highlightFields []string) *PostSearchOptions {
	_options.HighlightFields = highlightFields
	return _options
}

// SetHighlightNumber : Allow user to set HighlightNumber
func (_options *PostSearchOptions) SetHighlightNumber(highlightNumber int64) *PostSearchOptions {
	_options.HighlightNumber = core.Int64Ptr(highlightNumber)
	return _options
}

// SetHighlightPostTag : Allow user to set HighlightPostTag
func (_options *PostSearchOptions) SetHighlightPostTag(highlightPostTag string) *PostSearchOptions {
	_options.HighlightPostTag = core.StringPtr(highlightPostTag)
	return _options
}

// SetHighlightPreTag : Allow user to set HighlightPreTag
func (_options *PostSearchOptions) SetHighlightPreTag(highlightPreTag string) *PostSearchOptions {
	_options.HighlightPreTag = core.StringPtr(highlightPreTag)
	return _options
}

// SetHighlightSize : Allow user to set HighlightSize
func (_options *PostSearchOptions) SetHighlightSize(highlightSize int64) *PostSearchOptions {
	_options.HighlightSize = core.Int64Ptr(highlightSize)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *PostSearchOptions) SetIncludeDocs(includeDocs bool) *PostSearchOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetIncludeFields : Allow user to set IncludeFields
func (_options *PostSearchOptions) SetIncludeFields(includeFields []string) *PostSearchOptions {
	_options.IncludeFields = includeFields
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostSearchOptions) SetLimit(limit int64) *PostSearchOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *PostSearchOptions) SetSort(sort []string) *PostSearchOptions {
	_options.Sort = sort
	return _options
}

// SetStale : Allow user to set Stale
func (_options *PostSearchOptions) SetStale(stale string) *PostSearchOptions {
	_options.Stale = core.StringPtr(stale)
	return _options
}

// SetCounts : Allow user to set Counts
func (_options *PostSearchOptions) SetCounts(counts []string) *PostSearchOptions {
	_options.Counts = counts
	return _options
}

// SetDrilldown : Allow user to set Drilldown
func (_options *PostSearchOptions) SetDrilldown(drilldown [][]string) *PostSearchOptions {
	_options.Drilldown = drilldown
	return _options
}

// SetGroupField : Allow user to set GroupField
func (_options *PostSearchOptions) SetGroupField(groupField string) *PostSearchOptions {
	_options.GroupField = core.StringPtr(groupField)
	return _options
}

// SetGroupLimit : Allow user to set GroupLimit
func (_options *PostSearchOptions) SetGroupLimit(groupLimit int64) *PostSearchOptions {
	_options.GroupLimit = core.Int64Ptr(groupLimit)
	return _options
}

// SetGroupSort : Allow user to set GroupSort
func (_options *PostSearchOptions) SetGroupSort(groupSort []string) *PostSearchOptions {
	_options.GroupSort = groupSort
	return _options
}

// SetRanges : Allow user to set Ranges
func (_options *PostSearchOptions) SetRanges(ranges map[string]map[string]map[string]string) *PostSearchOptions {
	_options.Ranges = ranges
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostSearchOptions) SetHeaders(param map[string]string) *PostSearchOptions {
	options.Headers = param
	return options
}

// PostViewOptions : The PostView options.
type PostViewOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the map reduce view function name.
	View *string `json:"view" validate:"required,ne="`

	// Parameter to specify whether to include the encoding information in attachment stubs if the particular attachment is
	// compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the returned
	// document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Parameter to specify whether the specified end key should be included in the result.
	InclusiveEnd *bool `json:"inclusive_end,omitempty"`

	// Parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Parameter to specify whether to include in the response an update_seq value indicating the sequence id of the
	// database the view reflects.
	UpdateSeq *bool `json:"update_seq,omitempty"`

	// Schema for any JSON type.
	Endkey interface{} `json:"endkey,omitempty"`

	// Schema for a document ID.
	EndkeyDocid *string `json:"endkey_docid,omitempty"`

	// Parameter to specify whether to group the results using the reduce function to a group rather than a single row.
	// Implies reduce is true and the maximum group_level.
	Group *bool `json:"group,omitempty"`

	// Parameter to specify the group level to be used. Implies group is true.
	GroupLevel *int64 `json:"group_level,omitempty"`

	// Schema for any JSON type.
	Key interface{} `json:"key,omitempty"`

	// Parameter to specify to return only documents that match the specified keys. String representation of a JSON array
	// containing elements that match the key type emitted by the view function.
	Keys []interface{} `json:"keys,omitempty"`

	// Parameter to specify whether to use the reduce function in a map-reduce view. Default is true when a reduce function
	// is defined.
	Reduce *bool `json:"reduce,omitempty"`

	// Parameter to specify whether view results should be returned from a stable set of shards.
	Stable *bool `json:"stable,omitempty"`

	// Schema for any JSON type.
	Startkey interface{} `json:"startkey,omitempty"`

	// Schema for a document ID.
	StartkeyDocid *string `json:"startkey_docid,omitempty"`

	// Parameter to specify whether or not the view in question should be updated prior to responding to the user.
	Update *string `json:"update,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostViewOptions.Update property.
// Parameter to specify whether or not the view in question should be updated prior to responding to the user.
const (
	PostViewOptionsUpdateFalseConst = "false"
	PostViewOptionsUpdateLazyConst = "lazy"
	PostViewOptionsUpdateTrueConst = "true"
)

// NewPostViewOptions : Instantiate PostViewOptions
func (*CloudantV1) NewPostViewOptions(db string, ddoc string, view string) *PostViewOptions {
	return &PostViewOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
		View: core.StringPtr(view),
	}
}

// SetDb : Allow user to set Db
func (_options *PostViewOptions) SetDb(db string) *PostViewOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *PostViewOptions) SetDdoc(ddoc string) *PostViewOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetView : Allow user to set View
func (_options *PostViewOptions) SetView(view string) *PostViewOptions {
	_options.View = core.StringPtr(view)
	return _options
}

// SetAttEncodingInfo : Allow user to set AttEncodingInfo
func (_options *PostViewOptions) SetAttEncodingInfo(attEncodingInfo bool) *PostViewOptions {
	_options.AttEncodingInfo = core.BoolPtr(attEncodingInfo)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *PostViewOptions) SetAttachments(attachments bool) *PostViewOptions {
	_options.Attachments = core.BoolPtr(attachments)
	return _options
}

// SetConflicts : Allow user to set Conflicts
func (_options *PostViewOptions) SetConflicts(conflicts bool) *PostViewOptions {
	_options.Conflicts = core.BoolPtr(conflicts)
	return _options
}

// SetDescending : Allow user to set Descending
func (_options *PostViewOptions) SetDescending(descending bool) *PostViewOptions {
	_options.Descending = core.BoolPtr(descending)
	return _options
}

// SetIncludeDocs : Allow user to set IncludeDocs
func (_options *PostViewOptions) SetIncludeDocs(includeDocs bool) *PostViewOptions {
	_options.IncludeDocs = core.BoolPtr(includeDocs)
	return _options
}

// SetInclusiveEnd : Allow user to set InclusiveEnd
func (_options *PostViewOptions) SetInclusiveEnd(inclusiveEnd bool) *PostViewOptions {
	_options.InclusiveEnd = core.BoolPtr(inclusiveEnd)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *PostViewOptions) SetLimit(limit int64) *PostViewOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSkip : Allow user to set Skip
func (_options *PostViewOptions) SetSkip(skip int64) *PostViewOptions {
	_options.Skip = core.Int64Ptr(skip)
	return _options
}

// SetUpdateSeq : Allow user to set UpdateSeq
func (_options *PostViewOptions) SetUpdateSeq(updateSeq bool) *PostViewOptions {
	_options.UpdateSeq = core.BoolPtr(updateSeq)
	return _options
}

// SetEndkey : Allow user to set Endkey
func (_options *PostViewOptions) SetEndkey(endkey interface{}) *PostViewOptions {
	_options.Endkey = endkey
	return _options
}

// SetEndkeyDocid : Allow user to set EndkeyDocid
func (_options *PostViewOptions) SetEndkeyDocid(endkeyDocid string) *PostViewOptions {
	_options.EndkeyDocid = core.StringPtr(endkeyDocid)
	return _options
}

// SetGroup : Allow user to set Group
func (_options *PostViewOptions) SetGroup(group bool) *PostViewOptions {
	_options.Group = core.BoolPtr(group)
	return _options
}

// SetGroupLevel : Allow user to set GroupLevel
func (_options *PostViewOptions) SetGroupLevel(groupLevel int64) *PostViewOptions {
	_options.GroupLevel = core.Int64Ptr(groupLevel)
	return _options
}

// SetKey : Allow user to set Key
func (_options *PostViewOptions) SetKey(key interface{}) *PostViewOptions {
	_options.Key = key
	return _options
}

// SetKeys : Allow user to set Keys
func (_options *PostViewOptions) SetKeys(keys []interface{}) *PostViewOptions {
	_options.Keys = keys
	return _options
}

// SetReduce : Allow user to set Reduce
func (_options *PostViewOptions) SetReduce(reduce bool) *PostViewOptions {
	_options.Reduce = core.BoolPtr(reduce)
	return _options
}

// SetStable : Allow user to set Stable
func (_options *PostViewOptions) SetStable(stable bool) *PostViewOptions {
	_options.Stable = core.BoolPtr(stable)
	return _options
}

// SetStartkey : Allow user to set Startkey
func (_options *PostViewOptions) SetStartkey(startkey interface{}) *PostViewOptions {
	_options.Startkey = startkey
	return _options
}

// SetStartkeyDocid : Allow user to set StartkeyDocid
func (_options *PostViewOptions) SetStartkeyDocid(startkeyDocid string) *PostViewOptions {
	_options.StartkeyDocid = core.StringPtr(startkeyDocid)
	return _options
}

// SetUpdate : Allow user to set Update
func (_options *PostViewOptions) SetUpdate(update string) *PostViewOptions {
	_options.Update = core.StringPtr(update)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostViewOptions) SetHeaders(param map[string]string) *PostViewOptions {
	options.Headers = param
	return options
}

// PostViewQueriesOptions : The PostViewQueries options.
type PostViewQueriesOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// Path parameter to specify the map reduce view function name.
	View *string `json:"view" validate:"required,ne="`

	// An array of query objects with fields for the parameters of each individual view query to be executed. The field
	// names and their meaning are the same as the query parameters of a regular view request.
	Queries []ViewQuery `json:"queries" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostViewQueriesOptions : Instantiate PostViewQueriesOptions
func (*CloudantV1) NewPostViewQueriesOptions(db string, ddoc string, view string, queries []ViewQuery) *PostViewQueriesOptions {
	return &PostViewQueriesOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
		View: core.StringPtr(view),
		Queries: queries,
	}
}

// SetDb : Allow user to set Db
func (_options *PostViewQueriesOptions) SetDb(db string) *PostViewQueriesOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *PostViewQueriesOptions) SetDdoc(ddoc string) *PostViewQueriesOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetView : Allow user to set View
func (_options *PostViewQueriesOptions) SetView(view string) *PostViewQueriesOptions {
	_options.View = core.StringPtr(view)
	return _options
}

// SetQueries : Allow user to set Queries
func (_options *PostViewQueriesOptions) SetQueries(queries []ViewQuery) *PostViewQueriesOptions {
	_options.Queries = queries
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostViewQueriesOptions) SetHeaders(param map[string]string) *PostViewQueriesOptions {
	options.Headers = param
	return options
}

// PutAttachmentOptions : The PutAttachment options.
type PutAttachmentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// Path parameter to specify the attachment name.
	AttachmentName *string `json:"attachment_name" validate:"required,ne="`

	// HTTP request body for attachment operations.
	Attachment io.ReadCloser `json:"attachment" validate:"required"`

	// Content-Type of the attachment.
	ContentType *string `json:"Content-Type" validate:"required"`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPutAttachmentOptions : Instantiate PutAttachmentOptions
func (*CloudantV1) NewPutAttachmentOptions(db string, docID string, attachmentName string, attachment io.ReadCloser, contentType string) *PutAttachmentOptions {
	return &PutAttachmentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
		AttachmentName: core.StringPtr(attachmentName),
		Attachment: attachment,
		ContentType: core.StringPtr(contentType),
	}
}

// SetDb : Allow user to set Db
func (_options *PutAttachmentOptions) SetDb(db string) *PutAttachmentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *PutAttachmentOptions) SetDocID(docID string) *PutAttachmentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetAttachmentName : Allow user to set AttachmentName
func (_options *PutAttachmentOptions) SetAttachmentName(attachmentName string) *PutAttachmentOptions {
	_options.AttachmentName = core.StringPtr(attachmentName)
	return _options
}

// SetAttachment : Allow user to set Attachment
func (_options *PutAttachmentOptions) SetAttachment(attachment io.ReadCloser) *PutAttachmentOptions {
	_options.Attachment = attachment
	return _options
}

// SetContentType : Allow user to set ContentType
func (_options *PutAttachmentOptions) SetContentType(contentType string) *PutAttachmentOptions {
	_options.ContentType = core.StringPtr(contentType)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *PutAttachmentOptions) SetIfMatch(ifMatch string) *PutAttachmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *PutAttachmentOptions) SetRev(rev string) *PutAttachmentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutAttachmentOptions) SetHeaders(param map[string]string) *PutAttachmentOptions {
	options.Headers = param
	return options
}

// PutCapacityThroughputConfigurationOptions : The PutCapacityThroughputConfiguration options.
type PutCapacityThroughputConfigurationOptions struct {
	// A number of blocks of throughput units. A block consists of 100 reads/sec, 50 writes/sec, and 5 global queries/sec
	// of provisioned throughput capacity.
	Blocks *int64 `json:"blocks" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPutCapacityThroughputConfigurationOptions : Instantiate PutCapacityThroughputConfigurationOptions
func (*CloudantV1) NewPutCapacityThroughputConfigurationOptions(blocks int64) *PutCapacityThroughputConfigurationOptions {
	return &PutCapacityThroughputConfigurationOptions{
		Blocks: core.Int64Ptr(blocks),
	}
}

// SetBlocks : Allow user to set Blocks
func (_options *PutCapacityThroughputConfigurationOptions) SetBlocks(blocks int64) *PutCapacityThroughputConfigurationOptions {
	_options.Blocks = core.Int64Ptr(blocks)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutCapacityThroughputConfigurationOptions) SetHeaders(param map[string]string) *PutCapacityThroughputConfigurationOptions {
	options.Headers = param
	return options
}

// PutCloudantSecurityConfigurationOptions : The PutCloudantSecurityConfiguration options.
type PutCloudantSecurityConfigurationOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Database permissions for Cloudant users and/or API keys.
	Cloudant map[string][]string `json:"cloudant" validate:"required"`

	// Schema for names and roles to map to a database permission.
	Admins *SecurityObject `json:"admins,omitempty"`

	// Schema for names and roles to map to a database permission.
	Members *SecurityObject `json:"members,omitempty"`

	// Manage permissions using the `_users` database only.
	CouchdbAuthOnly *bool `json:"couchdb_auth_only,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutCloudantSecurityConfigurationOptions.Cloudant property.
// Database permissions for Cloudant users and/or API keys.
const (
	PutCloudantSecurityConfigurationOptionsCloudantAdminConst = "_admin"
	PutCloudantSecurityConfigurationOptionsCloudantDbUpdatesConst = "_db_updates"
	PutCloudantSecurityConfigurationOptionsCloudantDesignConst = "_design"
	PutCloudantSecurityConfigurationOptionsCloudantReaderConst = "_reader"
	PutCloudantSecurityConfigurationOptionsCloudantReplicatorConst = "_replicator"
	PutCloudantSecurityConfigurationOptionsCloudantSecurityConst = "_security"
	PutCloudantSecurityConfigurationOptionsCloudantShardsConst = "_shards"
	PutCloudantSecurityConfigurationOptionsCloudantWriterConst = "_writer"
)

// NewPutCloudantSecurityConfigurationOptions : Instantiate PutCloudantSecurityConfigurationOptions
func (*CloudantV1) NewPutCloudantSecurityConfigurationOptions(db string, cloudant map[string][]string) *PutCloudantSecurityConfigurationOptions {
	return &PutCloudantSecurityConfigurationOptions{
		Db: core.StringPtr(db),
		Cloudant: cloudant,
	}
}

// SetDb : Allow user to set Db
func (_options *PutCloudantSecurityConfigurationOptions) SetDb(db string) *PutCloudantSecurityConfigurationOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetCloudant : Allow user to set Cloudant
func (_options *PutCloudantSecurityConfigurationOptions) SetCloudant(cloudant map[string][]string) *PutCloudantSecurityConfigurationOptions {
	_options.Cloudant = cloudant
	return _options
}

// SetAdmins : Allow user to set Admins
func (_options *PutCloudantSecurityConfigurationOptions) SetAdmins(admins *SecurityObject) *PutCloudantSecurityConfigurationOptions {
	_options.Admins = admins
	return _options
}

// SetMembers : Allow user to set Members
func (_options *PutCloudantSecurityConfigurationOptions) SetMembers(members *SecurityObject) *PutCloudantSecurityConfigurationOptions {
	_options.Members = members
	return _options
}

// SetCouchdbAuthOnly : Allow user to set CouchdbAuthOnly
func (_options *PutCloudantSecurityConfigurationOptions) SetCouchdbAuthOnly(couchdbAuthOnly bool) *PutCloudantSecurityConfigurationOptions {
	_options.CouchdbAuthOnly = core.BoolPtr(couchdbAuthOnly)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutCloudantSecurityConfigurationOptions) SetHeaders(param map[string]string) *PutCloudantSecurityConfigurationOptions {
	options.Headers = param
	return options
}

// PutCorsConfigurationOptions : The PutCorsConfiguration options.
type PutCorsConfigurationOptions struct {
	// An array of strings that contain allowed origin domains. You have to specify the full URL including the protocol. It
	// is recommended that only the HTTPS protocol is used. Subdomains count as separate domains, so you have to specify
	// all subdomains used.
	Origins []string `json:"origins" validate:"required"`

	// Boolean value to allow authentication credentials. If set to true, browser requests must be done by using
	// withCredentials = true.
	AllowCredentials *bool `json:"allow_credentials,omitempty"`

	// Boolean value to turn CORS on and off.
	EnableCors *bool `json:"enable_cors,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPutCorsConfigurationOptions : Instantiate PutCorsConfigurationOptions
func (*CloudantV1) NewPutCorsConfigurationOptions(origins []string) *PutCorsConfigurationOptions {
	return &PutCorsConfigurationOptions{
		Origins: origins,
	}
}

// SetOrigins : Allow user to set Origins
func (_options *PutCorsConfigurationOptions) SetOrigins(origins []string) *PutCorsConfigurationOptions {
	_options.Origins = origins
	return _options
}

// SetAllowCredentials : Allow user to set AllowCredentials
func (_options *PutCorsConfigurationOptions) SetAllowCredentials(allowCredentials bool) *PutCorsConfigurationOptions {
	_options.AllowCredentials = core.BoolPtr(allowCredentials)
	return _options
}

// SetEnableCors : Allow user to set EnableCors
func (_options *PutCorsConfigurationOptions) SetEnableCors(enableCors bool) *PutCorsConfigurationOptions {
	_options.EnableCors = core.BoolPtr(enableCors)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutCorsConfigurationOptions) SetHeaders(param map[string]string) *PutCorsConfigurationOptions {
	options.Headers = param
	return options
}

// PutDatabaseOptions : The PutDatabase options.
type PutDatabaseOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Query parameter to specify whether to enable database partitions when creating a database.
	Partitioned *bool `json:"partitioned,omitempty"`

	// The number of shards in the database. Each shard is a partition of the hash value range. Its value is set by the
	// service. For more information about modifying database configuration, contact IBM Cloudant support.
	Q *int64 `json:"q,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPutDatabaseOptions : Instantiate PutDatabaseOptions
func (*CloudantV1) NewPutDatabaseOptions(db string) *PutDatabaseOptions {
	return &PutDatabaseOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *PutDatabaseOptions) SetDb(db string) *PutDatabaseOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetPartitioned : Allow user to set Partitioned
func (_options *PutDatabaseOptions) SetPartitioned(partitioned bool) *PutDatabaseOptions {
	_options.Partitioned = core.BoolPtr(partitioned)
	return _options
}

// SetQ : Allow user to set Q
func (_options *PutDatabaseOptions) SetQ(q int64) *PutDatabaseOptions {
	_options.Q = core.Int64Ptr(q)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutDatabaseOptions) SetHeaders(param map[string]string) *PutDatabaseOptions {
	options.Headers = param
	return options
}

// PutDesignDocumentOptions : The PutDesignDocument options.
type PutDesignDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the design document name. The design document name is the design document ID excluding the
	// `_design/` prefix.
	Ddoc *string `json:"ddoc" validate:"required,ne="`

	// HTTP request body for DesignDocument operations.
	DesignDocument *DesignDocument `json:"designDocument" validate:"required"`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Query parameter to specify whether to prevent insertion of conflicting document revisions. If false, a well-formed
	// _rev must be included in the document. False is used by the replicator to insert documents into the target database
	// even if that leads to the creation of conflicts.
	NewEdits *bool `json:"new_edits,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutDesignDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	PutDesignDocumentOptionsBatchOkConst = "ok"
)

// NewPutDesignDocumentOptions : Instantiate PutDesignDocumentOptions
func (*CloudantV1) NewPutDesignDocumentOptions(db string, ddoc string, designDocument *DesignDocument) *PutDesignDocumentOptions {
	return &PutDesignDocumentOptions{
		Db: core.StringPtr(db),
		Ddoc: core.StringPtr(ddoc),
		DesignDocument: designDocument,
	}
}

// SetDb : Allow user to set Db
func (_options *PutDesignDocumentOptions) SetDb(db string) *PutDesignDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDdoc : Allow user to set Ddoc
func (_options *PutDesignDocumentOptions) SetDdoc(ddoc string) *PutDesignDocumentOptions {
	_options.Ddoc = core.StringPtr(ddoc)
	return _options
}

// SetDesignDocument : Allow user to set DesignDocument
func (_options *PutDesignDocumentOptions) SetDesignDocument(designDocument *DesignDocument) *PutDesignDocumentOptions {
	_options.DesignDocument = designDocument
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *PutDesignDocumentOptions) SetIfMatch(ifMatch string) *PutDesignDocumentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *PutDesignDocumentOptions) SetBatch(batch string) *PutDesignDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetNewEdits : Allow user to set NewEdits
func (_options *PutDesignDocumentOptions) SetNewEdits(newEdits bool) *PutDesignDocumentOptions {
	_options.NewEdits = core.BoolPtr(newEdits)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *PutDesignDocumentOptions) SetRev(rev string) *PutDesignDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutDesignDocumentOptions) SetHeaders(param map[string]string) *PutDesignDocumentOptions {
	options.Headers = param
	return options
}

// PutDocumentOptions : The PutDocument options.
type PutDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// HTTP request body for Document operations.
	Document *Document `json:"document,omitempty"`

	// HTTP request body for Document operations.
	Body io.ReadCloser `json:"body,omitempty"`

	// The type of the input.
	ContentType *string `json:"Content-Type,omitempty"`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Query parameter to specify whether to prevent insertion of conflicting document revisions. If false, a well-formed
	// _rev must be included in the document. False is used by the replicator to insert documents into the target database
	// even if that leads to the creation of conflicts.
	NewEdits *bool `json:"new_edits,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	PutDocumentOptionsBatchOkConst = "ok"
)

// NewPutDocumentOptions : Instantiate PutDocumentOptions
func (*CloudantV1) NewPutDocumentOptions(db string, docID string) *PutDocumentOptions {
	return &PutDocumentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *PutDocumentOptions) SetDb(db string) *PutDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *PutDocumentOptions) SetDocID(docID string) *PutDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetDocument : Allow user to set Document
func (_options *PutDocumentOptions) SetDocument(document *Document) *PutDocumentOptions {
	_options.Document = document
	return _options
}

// SetBody : Allow user to set Body
func (_options *PutDocumentOptions) SetBody(body io.ReadCloser) *PutDocumentOptions {
	_options.Body = body
	return _options
}

// SetContentType : Allow user to set ContentType
func (_options *PutDocumentOptions) SetContentType(contentType string) *PutDocumentOptions {
	_options.ContentType = core.StringPtr(contentType)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *PutDocumentOptions) SetIfMatch(ifMatch string) *PutDocumentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *PutDocumentOptions) SetBatch(batch string) *PutDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetNewEdits : Allow user to set NewEdits
func (_options *PutDocumentOptions) SetNewEdits(newEdits bool) *PutDocumentOptions {
	_options.NewEdits = core.BoolPtr(newEdits)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *PutDocumentOptions) SetRev(rev string) *PutDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutDocumentOptions) SetHeaders(param map[string]string) *PutDocumentOptions {
	options.Headers = param
	return options
}

// PutLocalDocumentOptions : The PutLocalDocument options.
type PutLocalDocumentOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// HTTP request body for Document operations.
	Document *Document `json:"document,omitempty"`

	// HTTP request body for Document operations.
	Body io.ReadCloser `json:"body,omitempty"`

	// The type of the input.
	ContentType *string `json:"Content-Type,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutLocalDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	PutLocalDocumentOptionsBatchOkConst = "ok"
)

// NewPutLocalDocumentOptions : Instantiate PutLocalDocumentOptions
func (*CloudantV1) NewPutLocalDocumentOptions(db string, docID string) *PutLocalDocumentOptions {
	return &PutLocalDocumentOptions{
		Db: core.StringPtr(db),
		DocID: core.StringPtr(docID),
	}
}

// SetDb : Allow user to set Db
func (_options *PutLocalDocumentOptions) SetDb(db string) *PutLocalDocumentOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetDocID : Allow user to set DocID
func (_options *PutLocalDocumentOptions) SetDocID(docID string) *PutLocalDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetDocument : Allow user to set Document
func (_options *PutLocalDocumentOptions) SetDocument(document *Document) *PutLocalDocumentOptions {
	_options.Document = document
	return _options
}

// SetBody : Allow user to set Body
func (_options *PutLocalDocumentOptions) SetBody(body io.ReadCloser) *PutLocalDocumentOptions {
	_options.Body = body
	return _options
}

// SetContentType : Allow user to set ContentType
func (_options *PutLocalDocumentOptions) SetContentType(contentType string) *PutLocalDocumentOptions {
	_options.ContentType = core.StringPtr(contentType)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *PutLocalDocumentOptions) SetBatch(batch string) *PutLocalDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutLocalDocumentOptions) SetHeaders(param map[string]string) *PutLocalDocumentOptions {
	options.Headers = param
	return options
}

// PutReplicationDocumentOptions : The PutReplicationDocument options.
type PutReplicationDocumentOptions struct {
	// Path parameter to specify the document ID.
	DocID *string `json:"doc_id" validate:"required,ne="`

	// HTTP request body for replication operations.
	ReplicationDocument *ReplicationDocument `json:"replicationDocument" validate:"required"`

	// Header parameter to specify the document revision. Alternative to rev query parameter.
	IfMatch *string `json:"If-Match,omitempty"`

	// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
	// code immediately.
	Batch *string `json:"batch,omitempty"`

	// Query parameter to specify whether to prevent insertion of conflicting document revisions. If false, a well-formed
	// _rev must be included in the document. False is used by the replicator to insert documents into the target database
	// even if that leads to the creation of conflicts.
	NewEdits *bool `json:"new_edits,omitempty"`

	// Query parameter to specify a document revision.
	Rev *string `json:"rev,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutReplicationDocumentOptions.Batch property.
// Query parameter to specify whether to store in batch mode. The server will respond with a HTTP 202 Accepted response
// code immediately.
const (
	PutReplicationDocumentOptionsBatchOkConst = "ok"
)

// NewPutReplicationDocumentOptions : Instantiate PutReplicationDocumentOptions
func (*CloudantV1) NewPutReplicationDocumentOptions(docID string, replicationDocument *ReplicationDocument) *PutReplicationDocumentOptions {
	return &PutReplicationDocumentOptions{
		DocID: core.StringPtr(docID),
		ReplicationDocument: replicationDocument,
	}
}

// SetDocID : Allow user to set DocID
func (_options *PutReplicationDocumentOptions) SetDocID(docID string) *PutReplicationDocumentOptions {
	_options.DocID = core.StringPtr(docID)
	return _options
}

// SetReplicationDocument : Allow user to set ReplicationDocument
func (_options *PutReplicationDocumentOptions) SetReplicationDocument(replicationDocument *ReplicationDocument) *PutReplicationDocumentOptions {
	_options.ReplicationDocument = replicationDocument
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *PutReplicationDocumentOptions) SetIfMatch(ifMatch string) *PutReplicationDocumentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetBatch : Allow user to set Batch
func (_options *PutReplicationDocumentOptions) SetBatch(batch string) *PutReplicationDocumentOptions {
	_options.Batch = core.StringPtr(batch)
	return _options
}

// SetNewEdits : Allow user to set NewEdits
func (_options *PutReplicationDocumentOptions) SetNewEdits(newEdits bool) *PutReplicationDocumentOptions {
	_options.NewEdits = core.BoolPtr(newEdits)
	return _options
}

// SetRev : Allow user to set Rev
func (_options *PutReplicationDocumentOptions) SetRev(rev string) *PutReplicationDocumentOptions {
	_options.Rev = core.StringPtr(rev)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutReplicationDocumentOptions) SetHeaders(param map[string]string) *PutReplicationDocumentOptions {
	options.Headers = param
	return options
}

// PutSecurityOptions : The PutSecurity options.
type PutSecurityOptions struct {
	// Path parameter to specify the database name.
	Db *string `json:"db" validate:"required,ne="`

	// Schema for names and roles to map to a database permission.
	Admins *SecurityObject `json:"admins,omitempty"`

	// Schema for names and roles to map to a database permission.
	Members *SecurityObject `json:"members,omitempty"`

	// Database permissions for Cloudant users and/or API keys.
	Cloudant map[string][]string `json:"cloudant,omitempty"`

	// Manage permissions using the `_users` database only.
	CouchdbAuthOnly *bool `json:"couchdb_auth_only,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PutSecurityOptions.Cloudant property.
// Database permissions for Cloudant users and/or API keys.
const (
	PutSecurityOptionsCloudantAdminConst = "_admin"
	PutSecurityOptionsCloudantDbUpdatesConst = "_db_updates"
	PutSecurityOptionsCloudantDesignConst = "_design"
	PutSecurityOptionsCloudantReaderConst = "_reader"
	PutSecurityOptionsCloudantReplicatorConst = "_replicator"
	PutSecurityOptionsCloudantSecurityConst = "_security"
	PutSecurityOptionsCloudantShardsConst = "_shards"
	PutSecurityOptionsCloudantWriterConst = "_writer"
)

// NewPutSecurityOptions : Instantiate PutSecurityOptions
func (*CloudantV1) NewPutSecurityOptions(db string) *PutSecurityOptions {
	return &PutSecurityOptions{
		Db: core.StringPtr(db),
	}
}

// SetDb : Allow user to set Db
func (_options *PutSecurityOptions) SetDb(db string) *PutSecurityOptions {
	_options.Db = core.StringPtr(db)
	return _options
}

// SetAdmins : Allow user to set Admins
func (_options *PutSecurityOptions) SetAdmins(admins *SecurityObject) *PutSecurityOptions {
	_options.Admins = admins
	return _options
}

// SetMembers : Allow user to set Members
func (_options *PutSecurityOptions) SetMembers(members *SecurityObject) *PutSecurityOptions {
	_options.Members = members
	return _options
}

// SetCloudant : Allow user to set Cloudant
func (_options *PutSecurityOptions) SetCloudant(cloudant map[string][]string) *PutSecurityOptions {
	_options.Cloudant = cloudant
	return _options
}

// SetCouchdbAuthOnly : Allow user to set CouchdbAuthOnly
func (_options *PutSecurityOptions) SetCouchdbAuthOnly(couchdbAuthOnly bool) *PutSecurityOptions {
	_options.CouchdbAuthOnly = core.BoolPtr(couchdbAuthOnly)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutSecurityOptions) SetHeaders(param map[string]string) *PutSecurityOptions {
	options.Headers = param
	return options
}

// ReplicationCreateTargetParameters : Request parameters to use during target database creation.
type ReplicationCreateTargetParameters struct {
	// Schema for the number of replicas of a database in a cluster.
	N *int64 `json:"n,omitempty"`

	// Parameter to specify whether to enable database partitions when creating the target database.
	Partitioned *bool `json:"partitioned,omitempty"`

	// Schema for the number of shards in a database. Each shard is a partition of the hash value range.
	Q *int64 `json:"q,omitempty"`
}

// UnmarshalReplicationCreateTargetParameters unmarshals an instance of ReplicationCreateTargetParameters from the specified map of raw messages.
func UnmarshalReplicationCreateTargetParameters(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReplicationCreateTargetParameters)
	err = core.UnmarshalPrimitive(m, "n", &obj.N)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "partitioned", &obj.Partitioned)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "q", &obj.Q)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplicationDatabase : Schema for a replication source or target database.
type ReplicationDatabase struct {
	// Schema for replication source or target database authentication.
	Auth *ReplicationDatabaseAuth `json:"auth,omitempty"`

	// Replication request headers.
	HeadersVar map[string]string `json:"headers,omitempty"`

	// Replication database URL.
	URL *string `json:"url" validate:"required"`
}

// NewReplicationDatabase : Instantiate ReplicationDatabase (Generic Model Constructor)
func (*CloudantV1) NewReplicationDatabase(url string) (_model *ReplicationDatabase, err error) {
	_model = &ReplicationDatabase{
		URL: core.StringPtr(url),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalReplicationDatabase unmarshals an instance of ReplicationDatabase from the specified map of raw messages.
func UnmarshalReplicationDatabase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReplicationDatabase)
	err = core.UnmarshalModel(m, "auth", &obj.Auth, UnmarshalReplicationDatabaseAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "headers", &obj.HeadersVar)
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

// ReplicationDatabaseAuth : Schema for replication source or target database authentication.
type ReplicationDatabaseAuth struct {
	// Schema for basic authentication of replication source or target database.
	Basic *ReplicationDatabaseAuthBasic `json:"basic,omitempty"`

	// Schema for an IAM API key for replication database authentication.
	Iam *ReplicationDatabaseAuthIam `json:"iam,omitempty"`
}

// UnmarshalReplicationDatabaseAuth unmarshals an instance of ReplicationDatabaseAuth from the specified map of raw messages.
func UnmarshalReplicationDatabaseAuth(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReplicationDatabaseAuth)
	err = core.UnmarshalModel(m, "basic", &obj.Basic, UnmarshalReplicationDatabaseAuthBasic)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "iam", &obj.Iam, UnmarshalReplicationDatabaseAuthIam)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplicationDatabaseAuthBasic : Schema for basic authentication of replication source or target database.
type ReplicationDatabaseAuthBasic struct {
	// The password associated with the username.
	Password *string `json:"password" validate:"required"`

	// The username.
	Username *string `json:"username" validate:"required"`
}

// NewReplicationDatabaseAuthBasic : Instantiate ReplicationDatabaseAuthBasic (Generic Model Constructor)
func (*CloudantV1) NewReplicationDatabaseAuthBasic(password string, username string) (_model *ReplicationDatabaseAuthBasic, err error) {
	_model = &ReplicationDatabaseAuthBasic{
		Password: core.StringPtr(password),
		Username: core.StringPtr(username),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalReplicationDatabaseAuthBasic unmarshals an instance of ReplicationDatabaseAuthBasic from the specified map of raw messages.
func UnmarshalReplicationDatabaseAuthBasic(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReplicationDatabaseAuthBasic)
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplicationDatabaseAuthIam : Schema for an IAM API key for replication database authentication.
type ReplicationDatabaseAuthIam struct {
	// IAM API key.
	ApiKey *string `json:"api_key" validate:"required"`
}

// NewReplicationDatabaseAuthIam : Instantiate ReplicationDatabaseAuthIam (Generic Model Constructor)
func (*CloudantV1) NewReplicationDatabaseAuthIam(apiKey string) (_model *ReplicationDatabaseAuthIam, err error) {
	_model = &ReplicationDatabaseAuthIam{
		ApiKey: core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalReplicationDatabaseAuthIam unmarshals an instance of ReplicationDatabaseAuthIam from the specified map of raw messages.
func UnmarshalReplicationDatabaseAuthIam(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReplicationDatabaseAuthIam)
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplicationDocument : Schema for a replication document. Note that `selector`, `doc_ids`, and `filter` are incompatible with each other.
type ReplicationDocument struct {
	// Schema for a map of attachment name to attachment metadata.
	Attachments map[string]Attachment `json:"_attachments,omitempty"`

	// Schema for a list of document revision identifiers.
	Conflicts []string `json:"_conflicts,omitempty"`

	// Deletion flag. Available if document was removed.
	Deleted *bool `json:"_deleted,omitempty"`

	// Schema for a list of document revision identifiers.
	DeletedConflicts []string `json:"_deleted_conflicts,omitempty"`

	// Document ID.
	ID *string `json:"_id,omitempty"`

	// Document's update sequence in current database. Available if requested with local_seq=true query parameter.
	LocalSeq *string `json:"_local_seq,omitempty"`

	// Schema for a document revision identifier.
	Rev *string `json:"_rev,omitempty"`

	// Schema for list of revision information.
	Revisions *Revisions `json:"_revisions,omitempty"`

	// Schema for a list of objects with information about local revisions and their status.
	RevsInfo []DocumentRevisionStatus `json:"_revs_info,omitempty"`

	// Cancels the replication.
	Cancel *bool `json:"cancel,omitempty"`

	// Defines replication checkpoint interval in milliseconds.
	CheckpointInterval *int64 `json:"checkpoint_interval,omitempty"`

	// HTTP connection timeout per replication. Even for very fast/reliable networks it might need to be increased if a
	// remote database is too busy.
	ConnectionTimeout *int64 `json:"connection_timeout,omitempty"`

	// Configure the replication to be continuous.
	Continuous *bool `json:"continuous,omitempty"`

	// Creates the target database. Requires administrator privileges on target server.
	CreateTarget *bool `json:"create_target,omitempty"`

	// Request parameters to use during target database creation.
	CreateTargetParams *ReplicationCreateTargetParameters `json:"create_target_params,omitempty"`

	// Schema for a list of document IDs.
	DocIds []string `json:"doc_ids,omitempty"`

	// The name of a filter function which is defined in a design document in the source database in {ddoc_id}/{filter}
	// format. It determines which documents get replicated. Using the selector option provides performance benefits when
	// compared with using the filter option. Use the selector option when possible.
	Filter *string `json:"filter,omitempty"`

	// Maximum number of HTTP connections per replication.
	HTTPConnections *int64 `json:"http_connections,omitempty"`

	// Schema for a map of string key value pairs, such as query parameters.
	QueryParams map[string]string `json:"query_params,omitempty"`

	// Number of times a replication request is retried. The requests are retried with a doubling exponential backoff
	// starting at 0.25 seconds, with a cap at 5 minutes.
	RetriesPerRequest *int64 `json:"retries_per_request,omitempty"`

	// JSON object describing criteria used to select documents. The selector specifies fields in the document, and
	// provides an expression to evaluate with the field content or other data.
	//
	// The selector object must:
	//   * Be structured as valid JSON.
	//   * Contain a valid query expression.
	//
	// Using a selector is significantly more efficient than using a JavaScript filter function, and is the recommended
	// option if filtering on document attributes only.
	//
	// Elementary selector syntax requires you to specify one or more fields, and the corresponding values required for
	// those fields. You can create more complex selector expressions by combining operators.
	//
	// Operators are identified by the use of a dollar sign `$` prefix in the name field.
	//
	// There are two core types of operators in the selector syntax:
	// * Combination operators: applied at the topmost level of selection. They are used to combine selectors. In addition
	// to the common boolean operators (`$and`, `$or`, `$not`, `$nor`) there are three combination operators: `$all`,
	// `$elemMatch`, and `$allMatch`. A combination operator takes a single argument. The argument is either another
	// selector, or an array of selectors.
	// * Condition operators: are specific to a field, and are used to evaluate the value stored in that field. For
	// instance, the basic `$eq` operator matches when the specified field contains a value that is equal to the supplied
	// argument.
	Selector map[string]interface{} `json:"selector,omitempty"`

	// Start the replication at a specific sequence value.
	SinceSeq *string `json:"since_seq,omitempty"`

	// Replication socket options.
	SocketOptions *string `json:"socket_options,omitempty"`

	// Schema for a replication source or target database.
	Source *ReplicationDatabase `json:"source" validate:"required"`

	// Address of a (http or socks5 protocol) proxy server through which replication with the source database should occur.
	SourceProxy *string `json:"source_proxy,omitempty"`

	// Schema for a replication source or target database.
	Target *ReplicationDatabase `json:"target" validate:"required"`

	// Address of a (http or socks5 protocol) proxy server through which replication with the target database should occur.
	TargetProxy *string `json:"target_proxy,omitempty"`

	// Specify if checkpoints should be saved during replication. Using checkpoints means a replication can be efficiently
	// resumed.
	UseCheckpoints *bool `json:"use_checkpoints,omitempty"`

	// Schema for the user context of a session.
	UserCtx *UserContext `json:"user_ctx,omitempty"`

	// Controls how many documents are processed. After each batch a checkpoint is written so this controls how frequently
	// checkpointing occurs.
	WorkerBatchSize *int64 `json:"worker_batch_size,omitempty"`

	// Controls how many separate processes will read from the changes manager and write to the target. A higher number can
	// improve throughput.
	WorkerProcesses *int64 `json:"worker_processes,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// NewReplicationDocument : Instantiate ReplicationDocument (Generic Model Constructor)
func (*CloudantV1) NewReplicationDocument(source *ReplicationDatabase, target *ReplicationDatabase) (_model *ReplicationDocument, err error) {
	_model = &ReplicationDocument{
		Source: source,
		Target: target,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// SetProperty allows the user to set an arbitrary property on an instance of ReplicationDocument
func (o *ReplicationDocument) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of ReplicationDocument
func (o *ReplicationDocument) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of ReplicationDocument
func (o *ReplicationDocument) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of ReplicationDocument
func (o *ReplicationDocument) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of ReplicationDocument
func (o *ReplicationDocument) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Attachments != nil {
		m["_attachments"] = o.Attachments
	}
	if o.Conflicts != nil {
		m["_conflicts"] = o.Conflicts
	}
	if o.Deleted != nil {
		m["_deleted"] = o.Deleted
	}
	if o.DeletedConflicts != nil {
		m["_deleted_conflicts"] = o.DeletedConflicts
	}
	if o.ID != nil {
		m["_id"] = o.ID
	}
	if o.LocalSeq != nil {
		m["_local_seq"] = o.LocalSeq
	}
	if o.Rev != nil {
		m["_rev"] = o.Rev
	}
	if o.Revisions != nil {
		m["_revisions"] = o.Revisions
	}
	if o.RevsInfo != nil {
		m["_revs_info"] = o.RevsInfo
	}
	if o.Cancel != nil {
		m["cancel"] = o.Cancel
	}
	if o.CheckpointInterval != nil {
		m["checkpoint_interval"] = o.CheckpointInterval
	}
	if o.ConnectionTimeout != nil {
		m["connection_timeout"] = o.ConnectionTimeout
	}
	if o.Continuous != nil {
		m["continuous"] = o.Continuous
	}
	if o.CreateTarget != nil {
		m["create_target"] = o.CreateTarget
	}
	if o.CreateTargetParams != nil {
		m["create_target_params"] = o.CreateTargetParams
	}
	if o.DocIds != nil {
		m["doc_ids"] = o.DocIds
	}
	if o.Filter != nil {
		m["filter"] = o.Filter
	}
	if o.HTTPConnections != nil {
		m["http_connections"] = o.HTTPConnections
	}
	if o.QueryParams != nil {
		m["query_params"] = o.QueryParams
	}
	if o.RetriesPerRequest != nil {
		m["retries_per_request"] = o.RetriesPerRequest
	}
	if o.Selector != nil {
		m["selector"] = o.Selector
	}
	if o.SinceSeq != nil {
		m["since_seq"] = o.SinceSeq
	}
	if o.SocketOptions != nil {
		m["socket_options"] = o.SocketOptions
	}
	if o.Source != nil {
		m["source"] = o.Source
	}
	if o.SourceProxy != nil {
		m["source_proxy"] = o.SourceProxy
	}
	if o.Target != nil {
		m["target"] = o.Target
	}
	if o.TargetProxy != nil {
		m["target_proxy"] = o.TargetProxy
	}
	if o.UseCheckpoints != nil {
		m["use_checkpoints"] = o.UseCheckpoints
	}
	if o.UserCtx != nil {
		m["user_ctx"] = o.UserCtx
	}
	if o.WorkerBatchSize != nil {
		m["worker_batch_size"] = o.WorkerBatchSize
	}
	if o.WorkerProcesses != nil {
		m["worker_processes"] = o.WorkerProcesses
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalReplicationDocument unmarshals an instance of ReplicationDocument from the specified map of raw messages.
func UnmarshalReplicationDocument(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReplicationDocument)
	err = core.UnmarshalModel(m, "_attachments", &obj.Attachments, UnmarshalAttachment)
	if err != nil {
		return
	}
	delete(m, "_attachments")
	err = core.UnmarshalPrimitive(m, "_conflicts", &obj.Conflicts)
	if err != nil {
		return
	}
	delete(m, "_conflicts")
	err = core.UnmarshalPrimitive(m, "_deleted", &obj.Deleted)
	if err != nil {
		return
	}
	delete(m, "_deleted")
	err = core.UnmarshalPrimitive(m, "_deleted_conflicts", &obj.DeletedConflicts)
	if err != nil {
		return
	}
	delete(m, "_deleted_conflicts")
	err = core.UnmarshalPrimitive(m, "_id", &obj.ID)
	if err != nil {
		return
	}
	delete(m, "_id")
	err = core.UnmarshalPrimitive(m, "_local_seq", &obj.LocalSeq)
	if err != nil {
		return
	}
	delete(m, "_local_seq")
	err = core.UnmarshalPrimitive(m, "_rev", &obj.Rev)
	if err != nil {
		return
	}
	delete(m, "_rev")
	err = core.UnmarshalModel(m, "_revisions", &obj.Revisions, UnmarshalRevisions)
	if err != nil {
		return
	}
	delete(m, "_revisions")
	err = core.UnmarshalModel(m, "_revs_info", &obj.RevsInfo, UnmarshalDocumentRevisionStatus)
	if err != nil {
		return
	}
	delete(m, "_revs_info")
	err = core.UnmarshalPrimitive(m, "cancel", &obj.Cancel)
	if err != nil {
		return
	}
	delete(m, "cancel")
	err = core.UnmarshalPrimitive(m, "checkpoint_interval", &obj.CheckpointInterval)
	if err != nil {
		return
	}
	delete(m, "checkpoint_interval")
	err = core.UnmarshalPrimitive(m, "connection_timeout", &obj.ConnectionTimeout)
	if err != nil {
		return
	}
	delete(m, "connection_timeout")
	err = core.UnmarshalPrimitive(m, "continuous", &obj.Continuous)
	if err != nil {
		return
	}
	delete(m, "continuous")
	err = core.UnmarshalPrimitive(m, "create_target", &obj.CreateTarget)
	if err != nil {
		return
	}
	delete(m, "create_target")
	err = core.UnmarshalModel(m, "create_target_params", &obj.CreateTargetParams, UnmarshalReplicationCreateTargetParameters)
	if err != nil {
		return
	}
	delete(m, "create_target_params")
	err = core.UnmarshalPrimitive(m, "doc_ids", &obj.DocIds)
	if err != nil {
		return
	}
	delete(m, "doc_ids")
	err = core.UnmarshalPrimitive(m, "filter", &obj.Filter)
	if err != nil {
		return
	}
	delete(m, "filter")
	err = core.UnmarshalPrimitive(m, "http_connections", &obj.HTTPConnections)
	if err != nil {
		return
	}
	delete(m, "http_connections")
	err = core.UnmarshalPrimitive(m, "query_params", &obj.QueryParams)
	if err != nil {
		return
	}
	delete(m, "query_params")
	err = core.UnmarshalPrimitive(m, "retries_per_request", &obj.RetriesPerRequest)
	if err != nil {
		return
	}
	delete(m, "retries_per_request")
	err = core.UnmarshalPrimitive(m, "selector", &obj.Selector)
	if err != nil {
		return
	}
	delete(m, "selector")
	err = core.UnmarshalPrimitive(m, "since_seq", &obj.SinceSeq)
	if err != nil {
		return
	}
	delete(m, "since_seq")
	err = core.UnmarshalPrimitive(m, "socket_options", &obj.SocketOptions)
	if err != nil {
		return
	}
	delete(m, "socket_options")
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalReplicationDatabase)
	if err != nil {
		return
	}
	delete(m, "source")
	err = core.UnmarshalPrimitive(m, "source_proxy", &obj.SourceProxy)
	if err != nil {
		return
	}
	delete(m, "source_proxy")
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalReplicationDatabase)
	if err != nil {
		return
	}
	delete(m, "target")
	err = core.UnmarshalPrimitive(m, "target_proxy", &obj.TargetProxy)
	if err != nil {
		return
	}
	delete(m, "target_proxy")
	err = core.UnmarshalPrimitive(m, "use_checkpoints", &obj.UseCheckpoints)
	if err != nil {
		return
	}
	delete(m, "use_checkpoints")
	err = core.UnmarshalModel(m, "user_ctx", &obj.UserCtx, UnmarshalUserContext)
	if err != nil {
		return
	}
	delete(m, "user_ctx")
	err = core.UnmarshalPrimitive(m, "worker_batch_size", &obj.WorkerBatchSize)
	if err != nil {
		return
	}
	delete(m, "worker_batch_size")
	err = core.UnmarshalPrimitive(m, "worker_processes", &obj.WorkerProcesses)
	if err != nil {
		return
	}
	delete(m, "worker_processes")
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

// Revisions : Schema for list of revision information.
type Revisions struct {
	// Array of valid revision IDs, in reverse order (latest first).
	Ids []string `json:"ids" validate:"required"`

	// Prefix number for the latest revision.
	Start *int64 `json:"start" validate:"required"`
}

// NewRevisions : Instantiate Revisions (Generic Model Constructor)
func (*CloudantV1) NewRevisions(ids []string, start int64) (_model *Revisions, err error) {
	_model = &Revisions{
		Ids: ids,
		Start: core.Int64Ptr(start),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRevisions unmarshals an instance of Revisions from the specified map of raw messages.
func UnmarshalRevisions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Revisions)
	err = core.UnmarshalPrimitive(m, "ids", &obj.Ids)
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

// RevsDiff : Schema for information about missing revs and possible ancestors.
type RevsDiff struct {
	// List of missing revisions.
	Missing []string `json:"missing,omitempty"`

	// List of possible ancestor revisions.
	PossibleAncestors []string `json:"possible_ancestors,omitempty"`
}

// UnmarshalRevsDiff unmarshals an instance of RevsDiff from the specified map of raw messages.
func UnmarshalRevsDiff(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RevsDiff)
	err = core.UnmarshalPrimitive(m, "missing", &obj.Missing)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "possible_ancestors", &obj.PossibleAncestors)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchedulerDocsResult : Schema for a listing of replication scheduler documents.
type SchedulerDocsResult struct {
	// Number of total rows.
	TotalRows *int64 `json:"total_rows" validate:"required"`

	// Array of replication scheduler doc objects.
	Docs []SchedulerDocument `json:"docs" validate:"required"`
}

// UnmarshalSchedulerDocsResult unmarshals an instance of SchedulerDocsResult from the specified map of raw messages.
func UnmarshalSchedulerDocsResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchedulerDocsResult)
	err = core.UnmarshalPrimitive(m, "total_rows", &obj.TotalRows)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "docs", &obj.Docs, UnmarshalSchedulerDocument)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchedulerDocument : Schema for a replication scheduler document.
type SchedulerDocument struct {
	// Database where replication document came from.
	Database *string `json:"database" validate:"required"`

	// Replication document ID.
	DocID *string `json:"doc_id" validate:"required"`

	// Consecutive errors count. Indicates how many times in a row this replication has crashed. Replication will be
	// retried with an exponential backoff based on this number. As soon as the replication succeeds this count is reset to
	// 0. To can be used to get an idea why a particular replication is not making progress.
	ErrorCount *int64 `json:"error_count" validate:"required"`

	// Replication ID, or null if state is completed or failed.
	ID *string `json:"id" validate:"required"`

	// Schema for scheduler document information. A JSON object that may contain additional information about the state.
	// For error states this will contain an error field and string value.
	Info *SchedulerInfo `json:"info" validate:"required"`

	// Timestamp of last state update.
	LastUpdated *strfmt.DateTime `json:"last_updated" validate:"required"`

	// Cluster node where the job is running.
	Node *string `json:"node,omitempty"`

	// Replication source.
	Source *string `json:"source,omitempty"`

	// Address of the (http or socks5 protocol) proxy server through which replication with the source database occurs.
	SourceProxy *string `json:"source_proxy,omitempty"`

	// Timestamp of when the replication was started.
	StartTime *strfmt.DateTime `json:"start_time" validate:"required"`

	// Schema for replication state.
	State *string `json:"state" validate:"required"`

	// Replication target.
	Target *string `json:"target,omitempty"`

	// Address of the (http or socks5 protocol) proxy server through which replication with the target database occurs.
	TargetProxy *string `json:"target_proxy,omitempty"`
}

// Constants associated with the SchedulerDocument.State property.
// Schema for replication state.
const (
	SchedulerDocumentStateCompletedConst = "completed"
	SchedulerDocumentStateCrashingConst = "crashing"
	SchedulerDocumentStateErrorConst = "error"
	SchedulerDocumentStateFailedConst = "failed"
	SchedulerDocumentStateInitializingConst = "initializing"
	SchedulerDocumentStatePendingConst = "pending"
	SchedulerDocumentStateRunningConst = "running"
)

// UnmarshalSchedulerDocument unmarshals an instance of SchedulerDocument from the specified map of raw messages.
func UnmarshalSchedulerDocument(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchedulerDocument)
	err = core.UnmarshalPrimitive(m, "database", &obj.Database)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_id", &obj.DocID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "info", &obj.Info, UnmarshalSchedulerInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_updated", &obj.LastUpdated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "node", &obj.Node)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_proxy", &obj.SourceProxy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_proxy", &obj.TargetProxy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchedulerInfo : Schema for scheduler document information. A JSON object that may contain additional information about the state. For
// error states this will contain an error field and string value.
type SchedulerInfo struct {
	// The count of changes not yet replicated.
	ChangesPending *int64 `json:"changes_pending,omitempty"`

	// The source sequence id which was last successfully replicated.
	CheckpointedSourceSeq *string `json:"checkpointed_source_seq,omitempty"`

	// The count of docs which failed to be written to the target.
	DocWriteFailures *int64 `json:"doc_write_failures,omitempty"`

	// The count of docs which have been read from the source.
	DocsRead *int64 `json:"docs_read,omitempty"`

	// The count of docs which have been written to the target.
	DocsWritten *int64 `json:"docs_written,omitempty"`

	// Replication error message.
	Error *string `json:"error,omitempty"`

	// The count of revisions which were found on the source, but missing from the target.
	MissingRevisionsFound *int64 `json:"missing_revisions_found,omitempty"`

	// The count of revisions which have been checked since this replication began.
	RevisionsChecked *int64 `json:"revisions_checked,omitempty"`

	// The last sequence number obtained from the source database changes feed.
	SourceSeq *string `json:"source_seq,omitempty"`

	// The last sequence number processed by the replicator.
	ThroughSeq *string `json:"through_seq,omitempty"`
}

// UnmarshalSchedulerInfo unmarshals an instance of SchedulerInfo from the specified map of raw messages.
func UnmarshalSchedulerInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchedulerInfo)
	err = core.UnmarshalPrimitive(m, "changes_pending", &obj.ChangesPending)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "checkpointed_source_seq", &obj.CheckpointedSourceSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_write_failures", &obj.DocWriteFailures)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "docs_read", &obj.DocsRead)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "docs_written", &obj.DocsWritten)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "missing_revisions_found", &obj.MissingRevisionsFound)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "revisions_checked", &obj.RevisionsChecked)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_seq", &obj.SourceSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "through_seq", &obj.ThroughSeq)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchedulerJob : Schema for a replication scheduler job.
type SchedulerJob struct {
	// Replication document database.
	Database *string `json:"database" validate:"required"`

	// Replication document ID.
	DocID *string `json:"doc_id" validate:"required"`

	// Timestamped history of events as a list of objects.
	History []SchedulerJobEvent `json:"history" validate:"required"`

	// Schema for a replication job id.
	ID *string `json:"id" validate:"required"`

	// Schema for scheduler document information. A JSON object that may contain additional information about the state.
	// For error states this will contain an error field and string value.
	Info *SchedulerInfo `json:"info" validate:"required"`

	// Cluster node where the job is running.
	Node *string `json:"node" validate:"required"`

	// Replication process ID.
	Pid *string `json:"pid" validate:"required"`

	// Replication source.
	Source *string `json:"source" validate:"required"`

	// Timestamp of when the replication was started.
	StartTime *strfmt.DateTime `json:"start_time" validate:"required"`

	// Replication target.
	Target *string `json:"target" validate:"required"`

	// Name of user running replication.
	User *string `json:"user" validate:"required"`
}

// UnmarshalSchedulerJob unmarshals an instance of SchedulerJob from the specified map of raw messages.
func UnmarshalSchedulerJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchedulerJob)
	err = core.UnmarshalPrimitive(m, "database", &obj.Database)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_id", &obj.DocID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalSchedulerJobEvent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "info", &obj.Info, UnmarshalSchedulerInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "node", &obj.Node)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pid", &obj.Pid)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user", &obj.User)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchedulerJobEvent : Schema for a replication scheduler job event.
type SchedulerJobEvent struct {
	// Reason for current state of event.
	Reason *string `json:"reason,omitempty"`

	// Timestamp of the event.
	Timestamp *strfmt.DateTime `json:"timestamp" validate:"required"`

	// Type of the event.
	Type *string `json:"type" validate:"required"`
}

// UnmarshalSchedulerJobEvent unmarshals an instance of SchedulerJobEvent from the specified map of raw messages.
func UnmarshalSchedulerJobEvent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchedulerJobEvent)
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timestamp", &obj.Timestamp)
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

// SchedulerJobsResult : Schema for a listing of replication scheduler jobs.
type SchedulerJobsResult struct {
	// Number of total rows.
	TotalRows *int64 `json:"total_rows" validate:"required"`

	// Array of replication job objects.
	Jobs []SchedulerJob `json:"jobs" validate:"required"`
}

// UnmarshalSchedulerJobsResult unmarshals an instance of SchedulerJobsResult from the specified map of raw messages.
func UnmarshalSchedulerJobsResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchedulerJobsResult)
	err = core.UnmarshalPrimitive(m, "total_rows", &obj.TotalRows)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "jobs", &obj.Jobs, UnmarshalSchedulerJob)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SearchAnalyzeResult : Schema for the output of testing search analyzer tokenization.
type SearchAnalyzeResult struct {
	// tokens.
	Tokens []string `json:"tokens" validate:"required"`
}

// UnmarshalSearchAnalyzeResult unmarshals an instance of SearchAnalyzeResult from the specified map of raw messages.
func UnmarshalSearchAnalyzeResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SearchAnalyzeResult)
	err = core.UnmarshalPrimitive(m, "tokens", &obj.Tokens)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SearchIndexDefinition : Schema for a search index definition.
type SearchIndexDefinition struct {
	// Schema for a search analyzer configuration.
	Analyzer *AnalyzerConfiguration `json:"analyzer,omitempty"`

	// String form of a JavaScript function that is called for each document in the database. The function takes the
	// document as a parameter, extracts some data from it, and then calls the `index` function to index that data. The
	// index function takes 2, or optionally 3, parameters.
	//
	// * The first parameter is the name of the field you intend to use when
	//   querying the index. If the special value `"default"` is used when you
	//   define the name, you do not have to specify a field name at query time.
	// * The second parameter is the data to be indexed. This data must be only a
	//   string, number, or boolean. Other types will cause an error to be thrown
	//   by the index function call.
	// * The optional third parameter is a JavaScript object with these
	//   properties:
	//
	//     * `facet` - boolean, default `false` - Creates a faceted index.
	//     * `index` - boolean, default `true` - If set to `false`, the data
	//       cannot be used for searches, but can still be retrieved from the
	//       index if `store` is set to `true`.
	//     * `store` - boolean, default `true` - If true, the value is returned
	//       in the search result; otherwise, the value is not returned.
	Index *string `json:"index" validate:"required"`
}

// NewSearchIndexDefinition : Instantiate SearchIndexDefinition (Generic Model Constructor)
func (*CloudantV1) NewSearchIndexDefinition(index string) (_model *SearchIndexDefinition, err error) {
	_model = &SearchIndexDefinition{
		Index: core.StringPtr(index),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSearchIndexDefinition unmarshals an instance of SearchIndexDefinition from the specified map of raw messages.
func UnmarshalSearchIndexDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SearchIndexDefinition)
	err = core.UnmarshalModel(m, "analyzer", &obj.Analyzer, UnmarshalAnalyzerConfiguration)
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

// SearchIndexInfo : Schema for metadata information about a search index.
type SearchIndexInfo struct {
	// The committed sequence identifier.
	CommittedSeq *int64 `json:"committed_seq" validate:"required"`

	// The size of the search index on disk.
	DiskSize *int64 `json:"disk_size" validate:"required"`

	// The count of the number of indexed documents.
	DocCount *int64 `json:"doc_count" validate:"required"`

	// The number of deleted documents.
	DocDelCount *int64 `json:"doc_del_count" validate:"required"`

	// The pending sequence identifier.
	PendingSeq *int64 `json:"pending_seq" validate:"required"`
}

// UnmarshalSearchIndexInfo unmarshals an instance of SearchIndexInfo from the specified map of raw messages.
func UnmarshalSearchIndexInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SearchIndexInfo)
	err = core.UnmarshalPrimitive(m, "committed_seq", &obj.CommittedSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disk_size", &obj.DiskSize)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_count", &obj.DocCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "doc_del_count", &obj.DocDelCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pending_seq", &obj.PendingSeq)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SearchInfoResult : Schema for search index information.
type SearchInfoResult struct {
	// The name of the search index prefixed by the design document ID where the index is stored.
	Name *string `json:"name" validate:"required"`

	// Schema for metadata information about a search index.
	SearchIndex *SearchIndexInfo `json:"search_index" validate:"required"`
}

// UnmarshalSearchInfoResult unmarshals an instance of SearchInfoResult from the specified map of raw messages.
func UnmarshalSearchInfoResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SearchInfoResult)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "search_index", &obj.SearchIndex, UnmarshalSearchIndexInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SearchResult : Schema for the result of a query search operation.
type SearchResult struct {
	// Number of total rows.
	TotalRows *int64 `json:"total_rows" validate:"required"`

	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark,omitempty"`

	// Grouped search matches.
	By *string `json:"by,omitempty"`

	// The counts facet syntax returns the number of query results for each unique value of each named field.
	Counts map[string]map[string]int64 `json:"counts,omitempty"`

	// The range facet syntax reuses the standard Lucene syntax for ranges to return counts of results that fit into each
	// specified category.
	Ranges map[string]map[string]int64 `json:"ranges,omitempty"`

	// Array of row objects.
	Rows []SearchResultRow `json:"rows,omitempty"`

	// Array of grouped search matches.
	Groups []SearchResultProperties `json:"groups,omitempty"`
}

// UnmarshalSearchResult unmarshals an instance of SearchResult from the specified map of raw messages.
func UnmarshalSearchResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SearchResult)
	err = core.UnmarshalPrimitive(m, "total_rows", &obj.TotalRows)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bookmark", &obj.Bookmark)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "by", &obj.By)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "counts", &obj.Counts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ranges", &obj.Ranges)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rows", &obj.Rows, UnmarshalSearchResultRow)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "groups", &obj.Groups, UnmarshalSearchResultProperties)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SearchResultProperties : Schema for the result of a query search operation.
type SearchResultProperties struct {
	// Number of total rows.
	TotalRows *int64 `json:"total_rows" validate:"required"`

	// Opaque bookmark token used when paginating results.
	Bookmark *string `json:"bookmark,omitempty"`

	// Grouped search matches.
	By *string `json:"by,omitempty"`

	// The counts facet syntax returns the number of query results for each unique value of each named field.
	Counts map[string]map[string]int64 `json:"counts,omitempty"`

	// The range facet syntax reuses the standard Lucene syntax for ranges to return counts of results that fit into each
	// specified category.
	Ranges map[string]map[string]int64 `json:"ranges,omitempty"`

	// Array of row objects.
	Rows []SearchResultRow `json:"rows,omitempty"`
}

// UnmarshalSearchResultProperties unmarshals an instance of SearchResultProperties from the specified map of raw messages.
func UnmarshalSearchResultProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SearchResultProperties)
	err = core.UnmarshalPrimitive(m, "total_rows", &obj.TotalRows)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bookmark", &obj.Bookmark)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "by", &obj.By)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "counts", &obj.Counts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ranges", &obj.Ranges)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rows", &obj.Rows, UnmarshalSearchResultRow)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SearchResultRow : Schema for a row of the result of a query search operation.
type SearchResultRow struct {
	// Schema for a document.
	Doc *Document `json:"doc,omitempty"`

	// Schema for the fields returned by a query search operation, a map of field name to value.
	Fields map[string]interface{} `json:"fields" validate:"required"`

	// Returns the context in which a search term was mentioned so that you can display more emphasized results to a user.
	Highlights map[string][]string `json:"highlights,omitempty"`

	// Schema for a document ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalSearchResultRow unmarshals an instance of SearchResultRow from the specified map of raw messages.
func UnmarshalSearchResultRow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SearchResultRow)
	err = core.UnmarshalModel(m, "doc", &obj.Doc, UnmarshalDocument)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "fields", &obj.Fields)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "highlights", &obj.Highlights)
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

// Security : Schema for a security document.
type Security struct {
	// Schema for names and roles to map to a database permission.
	Admins *SecurityObject `json:"admins,omitempty"`

	// Schema for names and roles to map to a database permission.
	Members *SecurityObject `json:"members,omitempty"`

	// Database permissions for Cloudant users and/or API keys.
	Cloudant map[string][]string `json:"cloudant,omitempty"`

	// Manage permissions using the `_users` database only.
	CouchdbAuthOnly *bool `json:"couchdb_auth_only,omitempty"`
}

// Constants associated with the Security.Cloudant property.
// Database permissions for Cloudant users and/or API keys.
const (
	SecurityCloudantAdminConst = "_admin"
	SecurityCloudantDbUpdatesConst = "_db_updates"
	SecurityCloudantDesignConst = "_design"
	SecurityCloudantReaderConst = "_reader"
	SecurityCloudantReplicatorConst = "_replicator"
	SecurityCloudantSecurityConst = "_security"
	SecurityCloudantShardsConst = "_shards"
	SecurityCloudantWriterConst = "_writer"
)

// UnmarshalSecurity unmarshals an instance of Security from the specified map of raw messages.
func UnmarshalSecurity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Security)
	err = core.UnmarshalModel(m, "admins", &obj.Admins, UnmarshalSecurityObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalSecurityObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloudant", &obj.Cloudant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "couchdb_auth_only", &obj.CouchdbAuthOnly)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecurityObject : Schema for names and roles to map to a database permission.
type SecurityObject struct {
	// List of usernames.
	Names []string `json:"names,omitempty"`

	// List of roles.
	Roles []string `json:"roles,omitempty"`
}

// UnmarshalSecurityObject unmarshals an instance of SecurityObject from the specified map of raw messages.
func UnmarshalSecurityObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecurityObject)
	err = core.UnmarshalPrimitive(m, "names", &obj.Names)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "roles", &obj.Roles)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServerInformation : Schema for information about the server instance.
type ServerInformation struct {
	// Welcome message.
	Couchdb *string `json:"couchdb" validate:"required"`

	// List of enabled optional features.
	Features []string `json:"features" validate:"required"`

	// Schema for server vendor information.
	Vendor *ServerVendor `json:"vendor" validate:"required"`

	// Apache CouchDB version.
	Version *string `json:"version" validate:"required"`

	// List of feature flags.
	FeaturesFlags []string `json:"features_flags" validate:"required"`
}

// UnmarshalServerInformation unmarshals an instance of ServerInformation from the specified map of raw messages.
func UnmarshalServerInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServerInformation)
	err = core.UnmarshalPrimitive(m, "couchdb", &obj.Couchdb)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "features", &obj.Features)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "vendor", &obj.Vendor, UnmarshalServerVendor)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "features_flags", &obj.FeaturesFlags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServerVendor : Schema for server vendor information.
type ServerVendor struct {
	// Vendor name.
	Name *string `json:"name" validate:"required"`

	// Vendor variant.
	Variant *string `json:"variant,omitempty"`

	// Vendor version.
	Version *string `json:"version,omitempty"`
}

// UnmarshalServerVendor unmarshals an instance of ServerVendor from the specified map of raw messages.
func UnmarshalServerVendor(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServerVendor)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "variant", &obj.Variant)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SessionAuthentication : Schema for session authentication information.
type SessionAuthentication struct {
	// authenticated.
	Authenticated *string `json:"authenticated,omitempty"`

	// authentication_db.
	AuthenticationDb *string `json:"authentication_db,omitempty"`

	// authentication_handlers.
	AuthenticationHandlers []string `json:"authentication_handlers" validate:"required"`
}

// UnmarshalSessionAuthentication unmarshals an instance of SessionAuthentication from the specified map of raw messages.
func UnmarshalSessionAuthentication(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SessionAuthentication)
	err = core.UnmarshalPrimitive(m, "authenticated", &obj.Authenticated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "authentication_db", &obj.AuthenticationDb)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "authentication_handlers", &obj.AuthenticationHandlers)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SessionInformation : Schema for information about a session.
type SessionInformation struct {
	// ok.
	Ok *bool `json:"ok" validate:"required"`

	// Schema for session authentication information.
	Info *SessionAuthentication `json:"info" validate:"required"`

	// Schema for the user context of a session.
	UserCtx *UserContext `json:"userCtx" validate:"required"`
}

// UnmarshalSessionInformation unmarshals an instance of SessionInformation from the specified map of raw messages.
func UnmarshalSessionInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SessionInformation)
	err = core.UnmarshalPrimitive(m, "ok", &obj.Ok)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "info", &obj.Info, UnmarshalSessionAuthentication)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "userCtx", &obj.UserCtx, UnmarshalUserContext)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ShardsInformation : Schema for a shards object that maps the hash value range for each shard to the array of nodes that contain a copy of
// that shard.
type ShardsInformation struct {
	// Mapping of shard hash value range to a list of nodes.
	Shards map[string][]string `json:"shards" validate:"required"`
}

// UnmarshalShardsInformation unmarshals an instance of ShardsInformation from the specified map of raw messages.
func UnmarshalShardsInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ShardsInformation)
	err = core.UnmarshalPrimitive(m, "shards", &obj.Shards)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ThroughputInformation : Schema for detailed information about throughput capacity with breakdown by specific throughput requests classes.
type ThroughputInformation struct {
	// A number of blocks of throughput units. A block consists of 100 reads/sec, 50 writes/sec, and 5 global queries/sec
	// of provisioned throughput capacity.
	Blocks *int64 `json:"blocks" validate:"required"`

	// Provisioned global queries capacity in operations per second.
	Query *int64 `json:"query" validate:"required"`

	// Provisioned reads capacity in operations per second.
	Read *int64 `json:"read" validate:"required"`

	// Provisioned writes capacity in operations per second.
	Write *int64 `json:"write" validate:"required"`
}

// UnmarshalThroughputInformation unmarshals an instance of ThroughputInformation from the specified map of raw messages.
func UnmarshalThroughputInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ThroughputInformation)
	err = core.UnmarshalPrimitive(m, "blocks", &obj.Blocks)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "query", &obj.Query)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "read", &obj.Read)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "write", &obj.Write)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpInformation : Schema for information about the up state of the server.
type UpInformation struct {
	// seeds.
	Seeds interface{} `json:"seeds" validate:"required"`

	// status.
	Status *string `json:"status" validate:"required"`
}

// Constants associated with the UpInformation.Status property.
// status.
const (
	UpInformationStatusMaintenanceModeConst = "maintenance_mode"
	UpInformationStatusNolbConst = "nolb"
	UpInformationStatusOkConst = "ok"
)

// UnmarshalUpInformation unmarshals an instance of UpInformation from the specified map of raw messages.
func UnmarshalUpInformation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpInformation)
	err = core.UnmarshalPrimitive(m, "seeds", &obj.Seeds)
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

// UserContext : Schema for the user context of a session.
type UserContext struct {
	// Database name in the context of the provided operation.
	Db *string `json:"db,omitempty"`

	// User name.
	Name *string `json:"name" validate:"required"`

	// List of user roles.
	Roles []string `json:"roles" validate:"required"`
}

// Constants associated with the UserContext.Roles property.
// Schema for a security role.
const (
	UserContextRolesAdminConst = "_admin"
	UserContextRolesDbUpdatesConst = "_db_updates"
	UserContextRolesDesignConst = "_design"
	UserContextRolesReaderConst = "_reader"
	UserContextRolesReplicatorConst = "_replicator"
	UserContextRolesSecurityConst = "_security"
	UserContextRolesShardsConst = "_shards"
	UserContextRolesWriterConst = "_writer"
)

// NewUserContext : Instantiate UserContext (Generic Model Constructor)
func (*CloudantV1) NewUserContext(name string, roles []string) (_model *UserContext, err error) {
	_model = &UserContext{
		Name: core.StringPtr(name),
		Roles: roles,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalUserContext unmarshals an instance of UserContext from the specified map of raw messages.
func UnmarshalUserContext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserContext)
	err = core.UnmarshalPrimitive(m, "db", &obj.Db)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "roles", &obj.Roles)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UuidsResult : Schema for a set of uuids generated by the server.
type UuidsResult struct {
	// uuids.
	Uuids []string `json:"uuids" validate:"required"`
}

// UnmarshalUuidsResult unmarshals an instance of UuidsResult from the specified map of raw messages.
func UnmarshalUuidsResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UuidsResult)
	err = core.UnmarshalPrimitive(m, "uuids", &obj.Uuids)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ViewQueriesResult : Schema for the results of a queries view operation.
type ViewQueriesResult struct {
	// An array of result objects - one for each query. Each result object contains the same fields as the response to a
	// regular view request.
	Results []ViewResult `json:"results" validate:"required"`
}

// UnmarshalViewQueriesResult unmarshals an instance of ViewQueriesResult from the specified map of raw messages.
func UnmarshalViewQueriesResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ViewQueriesResult)
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalViewResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ViewQuery : Schema for a query view operation.
type ViewQuery struct {
	// Parameter to specify whether to include the encoding information in attachment stubs if the particular attachment is
	// compressed.
	AttEncodingInfo *bool `json:"att_encoding_info,omitempty"`

	// Parameter to specify whether to include attachments bodies in a response.
	Attachments *bool `json:"attachments,omitempty"`

	// Parameter to specify whether to include a list of conflicted revisions in the `_conflicts` property of the returned
	// document. Ignored if `include_docs` isn't `true`.
	Conflicts *bool `json:"conflicts,omitempty"`

	// Parameter to specify whether to return the documents in descending by key order.
	Descending *bool `json:"descending,omitempty"`

	// Parameter to specify whether to include the full content of the documents in the response.
	IncludeDocs *bool `json:"include_docs,omitempty"`

	// Parameter to specify whether the specified end key should be included in the result.
	InclusiveEnd *bool `json:"inclusive_end,omitempty"`

	// Parameter to specify the number of returned documents to limit the result to.
	Limit *int64 `json:"limit,omitempty"`

	// Parameter to specify the number of records before starting to return the results.
	Skip *int64 `json:"skip,omitempty"`

	// Parameter to specify whether to include in the response an update_seq value indicating the sequence id of the
	// database the view reflects.
	UpdateSeq *bool `json:"update_seq,omitempty"`

	// Schema for any JSON type.
	Endkey interface{} `json:"endkey,omitempty"`

	// Schema for a document ID.
	EndkeyDocid *string `json:"endkey_docid,omitempty"`

	// Parameter to specify whether to group the results using the reduce function to a group rather than a single row.
	// Implies reduce is true and the maximum group_level.
	Group *bool `json:"group,omitempty"`

	// Parameter to specify the group level to be used. Implies group is true.
	GroupLevel *int64 `json:"group_level,omitempty"`

	// Schema for any JSON type.
	Key interface{} `json:"key,omitempty"`

	// Parameter to specify to return only documents that match the specified keys. String representation of a JSON array
	// containing elements that match the key type emitted by the view function.
	Keys []interface{} `json:"keys,omitempty"`

	// Parameter to specify whether to use the reduce function in a map-reduce view. Default is true when a reduce function
	// is defined.
	Reduce *bool `json:"reduce,omitempty"`

	// Parameter to specify whether view results should be returned from a stable set of shards.
	Stable *bool `json:"stable,omitempty"`

	// Schema for any JSON type.
	Startkey interface{} `json:"startkey,omitempty"`

	// Schema for a document ID.
	StartkeyDocid *string `json:"startkey_docid,omitempty"`

	// Parameter to specify whether or not the view in question should be updated prior to responding to the user.
	Update *string `json:"update,omitempty"`
}

// Constants associated with the ViewQuery.Update property.
// Parameter to specify whether or not the view in question should be updated prior to responding to the user.
const (
	ViewQueryUpdateFalseConst = "false"
	ViewQueryUpdateLazyConst = "lazy"
	ViewQueryUpdateTrueConst = "true"
)

// UnmarshalViewQuery unmarshals an instance of ViewQuery from the specified map of raw messages.
func UnmarshalViewQuery(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ViewQuery)
	err = core.UnmarshalPrimitive(m, "att_encoding_info", &obj.AttEncodingInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachments", &obj.Attachments)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "conflicts", &obj.Conflicts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "descending", &obj.Descending)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "include_docs", &obj.IncludeDocs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inclusive_end", &obj.InclusiveEnd)
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
	err = core.UnmarshalPrimitive(m, "update_seq", &obj.UpdateSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "endkey", &obj.Endkey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "endkey_docid", &obj.EndkeyDocid)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "group", &obj.Group)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "group_level", &obj.GroupLevel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "keys", &obj.Keys)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reduce", &obj.Reduce)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "stable", &obj.Stable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "startkey", &obj.Startkey)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "startkey_docid", &obj.StartkeyDocid)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update", &obj.Update)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ViewResult : Schema for the result of a query view operation.
type ViewResult struct {
	// Number of total rows.
	TotalRows *int64 `json:"total_rows,omitempty"`

	// Current update sequence for the database.
	UpdateSeq *string `json:"update_seq,omitempty"`

	// rows.
	Rows []ViewResultRow `json:"rows" validate:"required"`
}

// UnmarshalViewResult unmarshals an instance of ViewResult from the specified map of raw messages.
func UnmarshalViewResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ViewResult)
	err = core.UnmarshalPrimitive(m, "total_rows", &obj.TotalRows)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update_seq", &obj.UpdateSeq)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rows", &obj.Rows, UnmarshalViewResultRow)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ViewResultRow : Schema for a row of a view result.
type ViewResultRow struct {
	// The cause of the error (if available).
	CausedBy *string `json:"caused_by,omitempty"`

	// The name of the error.
	Error *string `json:"error,omitempty"`

	// The reason the error occurred (if available).
	Reason *string `json:"reason,omitempty"`

	// Schema for a document.
	Doc *Document `json:"doc,omitempty"`

	// Schema for a document ID.
	ID *string `json:"id,omitempty"`

	// Schema for any JSON type.
	Key interface{} `json:"key" validate:"required"`

	// Schema for any JSON type.
	Value interface{} `json:"value" validate:"required"`
}

// UnmarshalViewResultRow unmarshals an instance of ViewResultRow from the specified map of raw messages.
func UnmarshalViewResultRow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ViewResultRow)
	err = core.UnmarshalPrimitive(m, "caused_by", &obj.CausedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "doc", &obj.Doc, UnmarshalDocument)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
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

// GeoJSONGeometry : Schema for a GeoJSON geometry.
// This model "extends" GeoJSONGeometryObject
type GeoJSONGeometry struct {
	// The type of GeoJSON Geometry.
	Type *string `json:"type" validate:"required"`

	// Used for all geometry types except `GeometryCollection`. The structure of the elements in the array varies by
	// geometry type.
	Coordinates []interface{} `json:"coordinates" validate:"required"`
}

// Constants associated with the GeoJSONGeometry.Type property.
// The type of GeoJSON Geometry.
const (
	GeoJSONGeometryTypeGeometrycollectionConst = "GeometryCollection"
	GeoJSONGeometryTypeLinestringConst = "LineString"
	GeoJSONGeometryTypeMultilinestringConst = "MultiLineString"
	GeoJSONGeometryTypeMultipointConst = "MultiPoint"
	GeoJSONGeometryTypeMultipolygonConst = "MultiPolygon"
	GeoJSONGeometryTypePointConst = "Point"
	GeoJSONGeometryTypePolygonConst = "Polygon"
)

func (*GeoJSONGeometry) isaGeoJSONGeometryObject() bool {
	return true
}

// UnmarshalGeoJSONGeometry unmarshals an instance of GeoJSONGeometry from the specified map of raw messages.
func UnmarshalGeoJSONGeometry(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoJSONGeometry)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "coordinates", &obj.Coordinates)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GeoJSONGeometryCollection : Schema for a GeoJSON GeometryCollection type geometry.
// This model "extends" GeoJSONGeometryObject
type GeoJSONGeometryCollection struct {
	// The type of GeoJSON Geometry.
	Type *string `json:"type" validate:"required"`

	// Used for the `GeometryCollection` type.
	Geometries []GeoJSONGeometry `json:"geometries" validate:"required"`
}

// Constants associated with the GeoJSONGeometryCollection.Type property.
// The type of GeoJSON Geometry.
const (
	GeoJSONGeometryCollectionTypeGeometrycollectionConst = "GeometryCollection"
	GeoJSONGeometryCollectionTypeLinestringConst = "LineString"
	GeoJSONGeometryCollectionTypeMultilinestringConst = "MultiLineString"
	GeoJSONGeometryCollectionTypeMultipointConst = "MultiPoint"
	GeoJSONGeometryCollectionTypeMultipolygonConst = "MultiPolygon"
	GeoJSONGeometryCollectionTypePointConst = "Point"
	GeoJSONGeometryCollectionTypePolygonConst = "Polygon"
)

func (*GeoJSONGeometryCollection) isaGeoJSONGeometryObject() bool {
	return true
}

// UnmarshalGeoJSONGeometryCollection unmarshals an instance of GeoJSONGeometryCollection from the specified map of raw messages.
func UnmarshalGeoJSONGeometryCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GeoJSONGeometryCollection)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "geometries", &obj.Geometries, UnmarshalGeoJSONGeometry)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
