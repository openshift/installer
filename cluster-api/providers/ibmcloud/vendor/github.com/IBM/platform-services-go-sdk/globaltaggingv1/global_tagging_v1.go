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
 * IBM OpenAPI SDK Code Generator Version: 3.105.0-3c13b041-20250605-193116
 */

// Package globaltaggingv1 : Operations and models for the GlobalTaggingV1 service
package globaltaggingv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
)

// GlobalTaggingV1 : Manage your tags with the Tagging API in IBM Cloud. You can attach, detach, delete, or list all of
// the tags in your billing account with the Tagging API. The tag name must be unique within a billing account. You can
// create tags in two formats: `key:value` or `label`. The tagging API supports three types of tag: `user` `service`,
// and `access` tags. `service` tags cannot be attached to IMS resources. `service` tags must be in the form
// `service_prefix:tag_label` where `service_prefix` identifies the Service owning the tag. `access` tags cannot be
// attached to IMS resources. They must be in the form `key:value`. You can replace all resource's tags using the
// `replace` query parameter in the attach operation. You can update the `value` of a resource's tag in the format
// `key:value`, using the `update` query parameter in the attach operation.
//
// API Version: 1.2.0
type GlobalTaggingV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://tags.global-search-tagging.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "global_tagging"

// GlobalTaggingV1Options : Service options
type GlobalTaggingV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewGlobalTaggingV1UsingExternalConfig : constructs an instance of GlobalTaggingV1 with passed in options and external configuration.
func NewGlobalTaggingV1UsingExternalConfig(options *GlobalTaggingV1Options) (globalTagging *GlobalTaggingV1, err error) {
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

	globalTagging, err = NewGlobalTaggingV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = globalTagging.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = globalTagging.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewGlobalTaggingV1 : constructs an instance of GlobalTaggingV1 with passed in options.
func NewGlobalTaggingV1(options *GlobalTaggingV1Options) (service *GlobalTaggingV1, err error) {
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

	service = &GlobalTaggingV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "globalTagging" suitable for processing requests.
func (globalTagging *GlobalTaggingV1) Clone() *GlobalTaggingV1 {
	if core.IsNil(globalTagging) {
		return nil
	}
	clone := *globalTagging
	clone.Service = globalTagging.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (globalTagging *GlobalTaggingV1) SetServiceURL(url string) error {
	err := globalTagging.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (globalTagging *GlobalTaggingV1) GetServiceURL() string {
	return globalTagging.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (globalTagging *GlobalTaggingV1) SetDefaultHeaders(headers http.Header) {
	globalTagging.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (globalTagging *GlobalTaggingV1) SetEnableGzipCompression(enableGzip bool) {
	globalTagging.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (globalTagging *GlobalTaggingV1) GetEnableGzipCompression() bool {
	return globalTagging.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (globalTagging *GlobalTaggingV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	globalTagging.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (globalTagging *GlobalTaggingV1) DisableRetries() {
	globalTagging.Service.DisableRetries()
}

// ListTags : Get all tags
// Lists all tags that are in a billing account. Use the `attached_to` parameter to return the list of tags that are
// attached to the specified resource.
func (globalTagging *GlobalTaggingV1) ListTags(listTagsOptions *ListTagsOptions) (result *TagList, response *core.DetailedResponse, err error) {
	result, response, err = globalTagging.ListTagsWithContext(context.Background(), listTagsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTagsWithContext is an alternate form of the ListTags method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) ListTagsWithContext(ctx context.Context, listTagsOptions *ListTagsOptions) (result *TagList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listTagsOptions, "listTagsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTagsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_tagging", "V1", "ListTags")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listTagsOptions.XRequestID != nil {
		builder.AddHeader("x-request-id", fmt.Sprint(*listTagsOptions.XRequestID))
	}
	if listTagsOptions.XCorrelationID != nil {
		builder.AddHeader("x-correlation-id", fmt.Sprint(*listTagsOptions.XCorrelationID))
	}

	if listTagsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listTagsOptions.AccountID))
	}
	if listTagsOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*listTagsOptions.TagType))
	}
	if listTagsOptions.FullData != nil {
		builder.AddQuery("full_data", fmt.Sprint(*listTagsOptions.FullData))
	}
	if listTagsOptions.Providers != nil {
		builder.AddQuery("providers", strings.Join(listTagsOptions.Providers, ","))
	}
	if listTagsOptions.AttachedTo != nil {
		builder.AddQuery("attached_to", fmt.Sprint(*listTagsOptions.AttachedTo))
	}
	if listTagsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listTagsOptions.Offset))
	}
	if listTagsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listTagsOptions.Limit))
	}
	if listTagsOptions.Timeout != nil {
		builder.AddQuery("timeout", fmt.Sprint(*listTagsOptions.Timeout))
	}
	if listTagsOptions.OrderByName != nil {
		builder.AddQuery("order_by_name", fmt.Sprint(*listTagsOptions.OrderByName))
	}
	if listTagsOptions.AttachedOnly != nil {
		builder.AddQuery("attached_only", fmt.Sprint(*listTagsOptions.AttachedOnly))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_tags", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTagList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTag : Create an access management tag
// Create an access management tag. To create an `access` tag, you must have the access listed in the [Granting users
// access to tag resources](https://cloud.ibm.com/docs/account?topic=account-access) documentation. `service` and `user`
// tags cannot be created upfront. They are created when they are attached for the first time to a resource.
func (globalTagging *GlobalTaggingV1) CreateTag(createTagOptions *CreateTagOptions) (result *CreateTagResults, response *core.DetailedResponse, err error) {
	result, response, err = globalTagging.CreateTagWithContext(context.Background(), createTagOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTagWithContext is an alternate form of the CreateTag method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) CreateTagWithContext(ctx context.Context, createTagOptions *CreateTagOptions) (result *CreateTagResults, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTagOptions, "createTagOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTagOptions, "createTagOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTagOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_tagging", "V1", "CreateTag")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createTagOptions.XRequestID != nil {
		builder.AddHeader("x-request-id", fmt.Sprint(*createTagOptions.XRequestID))
	}
	if createTagOptions.XCorrelationID != nil {
		builder.AddHeader("x-correlation-id", fmt.Sprint(*createTagOptions.XCorrelationID))
	}

	if createTagOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createTagOptions.AccountID))
	}
	if createTagOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*createTagOptions.TagType))
	}

	body := make(map[string]interface{})
	if createTagOptions.TagNames != nil {
		body["tag_names"] = createTagOptions.TagNames
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
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_tag", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateTagResults)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTagAll : Delete all unused tags
// Delete the tags that are not attached to any resource.
func (globalTagging *GlobalTaggingV1) DeleteTagAll(deleteTagAllOptions *DeleteTagAllOptions) (result *DeleteTagsResult, response *core.DetailedResponse, err error) {
	result, response, err = globalTagging.DeleteTagAllWithContext(context.Background(), deleteTagAllOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTagAllWithContext is an alternate form of the DeleteTagAll method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) DeleteTagAllWithContext(ctx context.Context, deleteTagAllOptions *DeleteTagAllOptions) (result *DeleteTagsResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(deleteTagAllOptions, "deleteTagAllOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTagAllOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_tagging", "V1", "DeleteTagAll")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteTagAllOptions.XRequestID != nil {
		builder.AddHeader("x-request-id", fmt.Sprint(*deleteTagAllOptions.XRequestID))
	}
	if deleteTagAllOptions.XCorrelationID != nil {
		builder.AddHeader("x-correlation-id", fmt.Sprint(*deleteTagAllOptions.XCorrelationID))
	}

	if deleteTagAllOptions.Providers != nil {
		builder.AddQuery("providers", fmt.Sprint(*deleteTagAllOptions.Providers))
	}
	if deleteTagAllOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteTagAllOptions.AccountID))
	}
	if deleteTagAllOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*deleteTagAllOptions.TagType))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tag_all", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteTagsResult)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTag : Delete an unused tag
// Delete an existing tag. A tag can be deleted only if it is not attached to any resource.
func (globalTagging *GlobalTaggingV1) DeleteTag(deleteTagOptions *DeleteTagOptions) (result *DeleteTagResults, response *core.DetailedResponse, err error) {
	result, response, err = globalTagging.DeleteTagWithContext(context.Background(), deleteTagOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTagWithContext is an alternate form of the DeleteTag method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) DeleteTagWithContext(ctx context.Context, deleteTagOptions *DeleteTagOptions) (result *DeleteTagResults, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTagOptions, "deleteTagOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTagOptions, "deleteTagOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"tag_name": *deleteTagOptions.TagName,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags/{tag_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTagOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_tagging", "V1", "DeleteTag")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteTagOptions.XRequestID != nil {
		builder.AddHeader("x-request-id", fmt.Sprint(*deleteTagOptions.XRequestID))
	}
	if deleteTagOptions.XCorrelationID != nil {
		builder.AddHeader("x-correlation-id", fmt.Sprint(*deleteTagOptions.XCorrelationID))
	}

	if deleteTagOptions.Providers != nil {
		builder.AddQuery("providers", strings.Join(deleteTagOptions.Providers, ","))
	}
	if deleteTagOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteTagOptions.AccountID))
	}
	if deleteTagOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*deleteTagOptions.TagType))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tag", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteTagResults)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// AttachTag : Attach tags
// Attaches one or more tags to one or more resources. Each resource can have no more than 1000 tags per each 'user' and
// 'service' type, and no more than 250 'access' tags (which is the account limit).
func (globalTagging *GlobalTaggingV1) AttachTag(attachTagOptions *AttachTagOptions) (result *TagResults, response *core.DetailedResponse, err error) {
	result, response, err = globalTagging.AttachTagWithContext(context.Background(), attachTagOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// AttachTagWithContext is an alternate form of the AttachTag method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) AttachTagWithContext(ctx context.Context, attachTagOptions *AttachTagOptions) (result *TagResults, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(attachTagOptions, "attachTagOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(attachTagOptions, "attachTagOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags/attach`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range attachTagOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_tagging", "V1", "AttachTag")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if attachTagOptions.XRequestID != nil {
		builder.AddHeader("x-request-id", fmt.Sprint(*attachTagOptions.XRequestID))
	}
	if attachTagOptions.XCorrelationID != nil {
		builder.AddHeader("x-correlation-id", fmt.Sprint(*attachTagOptions.XCorrelationID))
	}

	if attachTagOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*attachTagOptions.AccountID))
	}
	if attachTagOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*attachTagOptions.TagType))
	}
	if attachTagOptions.Replace != nil {
		builder.AddQuery("replace", fmt.Sprint(*attachTagOptions.Replace))
	}
	if attachTagOptions.Update != nil {
		builder.AddQuery("update", fmt.Sprint(*attachTagOptions.Update))
	}

	body := make(map[string]interface{})
	if attachTagOptions.TagName != nil {
		body["tag_name"] = attachTagOptions.TagName
	}
	if attachTagOptions.TagNames != nil {
		body["tag_names"] = attachTagOptions.TagNames
	}
	if attachTagOptions.Resources != nil {
		body["resources"] = attachTagOptions.Resources
	}
	if attachTagOptions.Query != nil {
		body["query"] = attachTagOptions.Query
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
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "attach_tag", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTagResults)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DetachTag : Detach tags
// Detaches one or more tags from one or more resources.
func (globalTagging *GlobalTaggingV1) DetachTag(detachTagOptions *DetachTagOptions) (result *TagResults, response *core.DetailedResponse, err error) {
	result, response, err = globalTagging.DetachTagWithContext(context.Background(), detachTagOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DetachTagWithContext is an alternate form of the DetachTag method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) DetachTagWithContext(ctx context.Context, detachTagOptions *DetachTagOptions) (result *TagResults, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(detachTagOptions, "detachTagOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(detachTagOptions, "detachTagOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags/detach`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range detachTagOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_tagging", "V1", "DetachTag")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if detachTagOptions.XRequestID != nil {
		builder.AddHeader("x-request-id", fmt.Sprint(*detachTagOptions.XRequestID))
	}
	if detachTagOptions.XCorrelationID != nil {
		builder.AddHeader("x-correlation-id", fmt.Sprint(*detachTagOptions.XCorrelationID))
	}

	if detachTagOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*detachTagOptions.AccountID))
	}
	if detachTagOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*detachTagOptions.TagType))
	}

	body := make(map[string]interface{})
	if detachTagOptions.TagName != nil {
		body["tag_name"] = detachTagOptions.TagName
	}
	if detachTagOptions.TagNames != nil {
		body["tag_names"] = detachTagOptions.TagNames
	}
	if detachTagOptions.Resources != nil {
		body["resources"] = detachTagOptions.Resources
	}
	if detachTagOptions.Query != nil {
		body["query"] = detachTagOptions.Query
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
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "detach_tag", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTagResults)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.2.0")
}

// AttachTagOptions : The AttachTag options.
type AttachTagOptions struct {
	// The name of the tag to attach.
	TagName *string `json:"tag_name,omitempty"`

	// An array of tag names to attach.
	TagNames []string `json:"tag_names,omitempty"`

	// List of resources on which the tagging operation operates on.
	Resources []Resource `json:"resources,omitempty"`

	// A valid Global Search string.
	Query *QueryString `json:"query,omitempty"`

	// An alphanumeric string that is used to trace the request. The value  may include ASCII alphanumerics and any of
	// following segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024
	// bytes. The value is considered invalid and must be ignored if that value includes any other character or is longer
	// than 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XRequestID *string `json:"x-request-id,omitempty"`

	// An alphanumeric string that is used to trace the request as a part of a larger context: the same value is used for
	// downstream requests and retries of those requests. The value may include ASCII alphanumerics and any of following
	// segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024 bytes.
	// The value is considered invalid and must be ignored if that value includes any other character or is longer than
	// 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XCorrelationID *string `json:"x-correlation-id,omitempty"`

	// The ID of the billing account of the tagged resource. It is a required parameter if `tag_type` is set to `service`.
	// Otherwise, it is inferred from the authorization IAM token.
	AccountID *string `json:"account_id,omitempty"`

	// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
	// for IMS resources.
	TagType *string `json:"tag_type,omitempty"`

	// Flag to request replacement of all attached tags. Set `true` if you want to replace all tags attached to the
	// resource with the current ones. Default value is false.
	Replace *bool `json:"replace,omitempty"`

	// Flag to request update of attached tags in the format `key:value`. Here's how it works for each tag in the request
	// body: If the tag to attach is in the format `key:value`, the System will atomically detach all existing tags
	// starting with `key:` and attach the new `key:value` tag. If no such tags exist, a new `key:value` tag will be
	// attached. If the tag is not in the `key:value` format (e.g., a simple label), the System will attach the label as
	// usual. The update query parameter is available for user and access management tags, but not for service tags.
	Update *bool `json:"update,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the AttachTagOptions.TagType property.
// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
// for IMS resources.
const (
	AttachTagOptionsTagTypeAccessConst  = "access"
	AttachTagOptionsTagTypeServiceConst = "service"
	AttachTagOptionsTagTypeUserConst    = "user"
)

// NewAttachTagOptions : Instantiate AttachTagOptions
func (*GlobalTaggingV1) NewAttachTagOptions() *AttachTagOptions {
	return &AttachTagOptions{}
}

// SetTagName : Allow user to set TagName
func (_options *AttachTagOptions) SetTagName(tagName string) *AttachTagOptions {
	_options.TagName = core.StringPtr(tagName)
	return _options
}

// SetTagNames : Allow user to set TagNames
func (_options *AttachTagOptions) SetTagNames(tagNames []string) *AttachTagOptions {
	_options.TagNames = tagNames
	return _options
}

// SetResources : Allow user to set Resources
func (_options *AttachTagOptions) SetResources(resources []Resource) *AttachTagOptions {
	_options.Resources = resources
	return _options
}

// SetQuery : Allow user to set Query
func (_options *AttachTagOptions) SetQuery(query *QueryString) *AttachTagOptions {
	_options.Query = query
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *AttachTagOptions) SetXRequestID(xRequestID string) *AttachTagOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *AttachTagOptions) SetXCorrelationID(xCorrelationID string) *AttachTagOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *AttachTagOptions) SetAccountID(accountID string) *AttachTagOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTagType : Allow user to set TagType
func (_options *AttachTagOptions) SetTagType(tagType string) *AttachTagOptions {
	_options.TagType = core.StringPtr(tagType)
	return _options
}

// SetReplace : Allow user to set Replace
func (_options *AttachTagOptions) SetReplace(replace bool) *AttachTagOptions {
	_options.Replace = core.BoolPtr(replace)
	return _options
}

// SetUpdate : Allow user to set Update
func (_options *AttachTagOptions) SetUpdate(update bool) *AttachTagOptions {
	_options.Update = core.BoolPtr(update)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *AttachTagOptions) SetHeaders(param map[string]string) *AttachTagOptions {
	options.Headers = param
	return options
}

// CreateTagOptions : The CreateTag options.
type CreateTagOptions struct {
	// An array of tag names to create.
	TagNames []string `json:"tag_names" validate:"required"`

	// An alphanumeric string that is used to trace the request. The value  may include ASCII alphanumerics and any of
	// following segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024
	// bytes. The value is considered invalid and must be ignored if that value includes any other character or is longer
	// than 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XRequestID *string `json:"x-request-id,omitempty"`

	// An alphanumeric string that is used to trace the request as a part of a larger context: the same value is used for
	// downstream requests and retries of those requests. The value may include ASCII alphanumerics and any of following
	// segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024 bytes.
	// The value is considered invalid and must be ignored if that value includes any other character or is longer than
	// 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XCorrelationID *string `json:"x-correlation-id,omitempty"`

	// The ID of the billing account where the tag must be created.
	AccountID *string `json:"account_id,omitempty"`

	// The type of the tags you want to create. The only allowed value is `access`.
	TagType *string `json:"tag_type,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateTagOptions.TagType property.
// The type of the tags you want to create. The only allowed value is `access`.
const (
	CreateTagOptionsTagTypeAccessConst = "access"
)

// NewCreateTagOptions : Instantiate CreateTagOptions
func (*GlobalTaggingV1) NewCreateTagOptions(tagNames []string) *CreateTagOptions {
	return &CreateTagOptions{
		TagNames: tagNames,
	}
}

// SetTagNames : Allow user to set TagNames
func (_options *CreateTagOptions) SetTagNames(tagNames []string) *CreateTagOptions {
	_options.TagNames = tagNames
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *CreateTagOptions) SetXRequestID(xRequestID string) *CreateTagOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateTagOptions) SetXCorrelationID(xCorrelationID string) *CreateTagOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateTagOptions) SetAccountID(accountID string) *CreateTagOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTagType : Allow user to set TagType
func (_options *CreateTagOptions) SetTagType(tagType string) *CreateTagOptions {
	_options.TagType = core.StringPtr(tagType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTagOptions) SetHeaders(param map[string]string) *CreateTagOptions {
	options.Headers = param
	return options
}

// CreateTagResults : Results of a create tag(s) request.
type CreateTagResults struct {
	// Array of results of a create_tag request.
	Results []CreateTagResultsResultsItem `json:"results,omitempty"`
}

// UnmarshalCreateTagResults unmarshals an instance of CreateTagResults from the specified map of raw messages.
func UnmarshalCreateTagResults(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateTagResults)
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalCreateTagResultsResultsItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "results-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateTagResultsResultsItem : CreateTagResultsResultsItem struct
type CreateTagResultsResultsItem struct {
	// The name of the tag created.
	TagName *string `json:"tag_name,omitempty"`

	// true if the tag was not created (for example, the tag already exists).
	IsError *bool `json:"is_error,omitempty"`
}

// UnmarshalCreateTagResultsResultsItem unmarshals an instance of CreateTagResultsResultsItem from the specified map of raw messages.
func UnmarshalCreateTagResultsResultsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateTagResultsResultsItem)
	err = core.UnmarshalPrimitive(m, "tag_name", &obj.TagName)
	if err != nil {
		err = core.SDKErrorf(err, "", "tag_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_error", &obj.IsError)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_error-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteTagAllOptions : The DeleteTagAll options.
type DeleteTagAllOptions struct {
	// An alphanumeric string that is used to trace the request. The value  may include ASCII alphanumerics and any of
	// following segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024
	// bytes. The value is considered invalid and must be ignored if that value includes any other character or is longer
	// than 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XRequestID *string `json:"x-request-id,omitempty"`

	// An alphanumeric string that is used to trace the request as a part of a larger context: the same value is used for
	// downstream requests and retries of those requests. The value may include ASCII alphanumerics and any of following
	// segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024 bytes.
	// The value is considered invalid and must be ignored if that value includes any other character or is longer than
	// 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XCorrelationID *string `json:"x-correlation-id,omitempty"`

	// Select a provider. Supported values are `ghost` and `ims`.
	Providers *string `json:"providers,omitempty"`

	// The ID of the billing account to delete the tags for. If it is not set, then it is taken from the authorization
	// token. It is a required parameter if `tag_type` is set to `service`.
	AccountID *string `json:"account_id,omitempty"`

	// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
	// for IMS resources (`providers` parameter set to `ims`).
	TagType *string `json:"tag_type,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the DeleteTagAllOptions.Providers property.
// Select a provider. Supported values are `ghost` and `ims`.
const (
	DeleteTagAllOptionsProvidersGhostConst = "ghost"
	DeleteTagAllOptionsProvidersImsConst   = "ims"
)

// Constants associated with the DeleteTagAllOptions.TagType property.
// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
// for IMS resources (`providers` parameter set to `ims`).
const (
	DeleteTagAllOptionsTagTypeAccessConst  = "access"
	DeleteTagAllOptionsTagTypeServiceConst = "service"
	DeleteTagAllOptionsTagTypeUserConst    = "user"
)

// NewDeleteTagAllOptions : Instantiate DeleteTagAllOptions
func (*GlobalTaggingV1) NewDeleteTagAllOptions() *DeleteTagAllOptions {
	return &DeleteTagAllOptions{}
}

// SetXRequestID : Allow user to set XRequestID
func (_options *DeleteTagAllOptions) SetXRequestID(xRequestID string) *DeleteTagAllOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteTagAllOptions) SetXCorrelationID(xCorrelationID string) *DeleteTagAllOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetProviders : Allow user to set Providers
func (_options *DeleteTagAllOptions) SetProviders(providers string) *DeleteTagAllOptions {
	_options.Providers = core.StringPtr(providers)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteTagAllOptions) SetAccountID(accountID string) *DeleteTagAllOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTagType : Allow user to set TagType
func (_options *DeleteTagAllOptions) SetTagType(tagType string) *DeleteTagAllOptions {
	_options.TagType = core.StringPtr(tagType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTagAllOptions) SetHeaders(param map[string]string) *DeleteTagAllOptions {
	options.Headers = param
	return options
}

// DeleteTagOptions : The DeleteTag options.
type DeleteTagOptions struct {
	// The name of tag to be deleted.
	TagName *string `json:"tag_name" validate:"required,ne="`

	// An alphanumeric string that is used to trace the request. The value  may include ASCII alphanumerics and any of
	// following segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024
	// bytes. The value is considered invalid and must be ignored if that value includes any other character or is longer
	// than 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XRequestID *string `json:"x-request-id,omitempty"`

	// An alphanumeric string that is used to trace the request as a part of a larger context: the same value is used for
	// downstream requests and retries of those requests. The value may include ASCII alphanumerics and any of following
	// segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024 bytes.
	// The value is considered invalid and must be ignored if that value includes any other character or is longer than
	// 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XCorrelationID *string `json:"x-correlation-id,omitempty"`

	// Select a provider. Supported values are `ghost` and `ims`. To delete tags both in Global Search and Tagging and in
	// IMS, use `ghost,ims`.
	Providers []string `json:"providers,omitempty"`

	// The ID of the billing account to delete the tag for. It is a required parameter if `tag_type` is set to `service`,
	// otherwise it is inferred from the authorization IAM token.
	AccountID *string `json:"account_id,omitempty"`

	// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
	// for IMS resources (`providers` parameter set to `ims`).
	TagType *string `json:"tag_type,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the DeleteTagOptions.Providers property.
const (
	DeleteTagOptionsProvidersGhostConst = "ghost"
	DeleteTagOptionsProvidersImsConst   = "ims"
)

// Constants associated with the DeleteTagOptions.TagType property.
// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
// for IMS resources (`providers` parameter set to `ims`).
const (
	DeleteTagOptionsTagTypeAccessConst  = "access"
	DeleteTagOptionsTagTypeServiceConst = "service"
	DeleteTagOptionsTagTypeUserConst    = "user"
)

// NewDeleteTagOptions : Instantiate DeleteTagOptions
func (*GlobalTaggingV1) NewDeleteTagOptions(tagName string) *DeleteTagOptions {
	return &DeleteTagOptions{
		TagName: core.StringPtr(tagName),
	}
}

// SetTagName : Allow user to set TagName
func (_options *DeleteTagOptions) SetTagName(tagName string) *DeleteTagOptions {
	_options.TagName = core.StringPtr(tagName)
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *DeleteTagOptions) SetXRequestID(xRequestID string) *DeleteTagOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteTagOptions) SetXCorrelationID(xCorrelationID string) *DeleteTagOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetProviders : Allow user to set Providers
func (_options *DeleteTagOptions) SetProviders(providers []string) *DeleteTagOptions {
	_options.Providers = providers
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteTagOptions) SetAccountID(accountID string) *DeleteTagOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTagType : Allow user to set TagType
func (_options *DeleteTagOptions) SetTagType(tagType string) *DeleteTagOptions {
	_options.TagType = core.StringPtr(tagType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTagOptions) SetHeaders(param map[string]string) *DeleteTagOptions {
	options.Headers = param
	return options
}

// DeleteTagResults : Results of a delete_tag request.
type DeleteTagResults struct {
	// Array of results of a delete_tag request.
	Results []DeleteTagResultsItem `json:"results,omitempty"`
}

// UnmarshalDeleteTagResults unmarshals an instance of DeleteTagResults from the specified map of raw messages.
func UnmarshalDeleteTagResults(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteTagResults)
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalDeleteTagResultsItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "results-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteTagResultsItem : Result of a delete_tag request.
// This type supports additional properties of type interface{}.
type DeleteTagResultsItem struct {
	// The provider of the tag.
	Provider *string `json:"provider,omitempty"`

	// It is `true` if the operation exits with an error (for example, the tag does not exist).
	IsError *bool `json:"is_error,omitempty"`

	// Allows users to set arbitrary properties of type interface{}.
	additionalProperties map[string]interface{}
}

// Constants associated with the DeleteTagResultsItem.Provider property.
// The provider of the tag.
const (
	DeleteTagResultsItemProviderGhostConst = "ghost"
	DeleteTagResultsItemProviderImsConst   = "ims"
)

// SetProperty allows the user to set an arbitrary property on an instance of DeleteTagResultsItem.
func (o *DeleteTagResultsItem) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of DeleteTagResultsItem.
func (o *DeleteTagResultsItem) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of DeleteTagResultsItem.
func (o *DeleteTagResultsItem) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of DeleteTagResultsItem.
func (o *DeleteTagResultsItem) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of DeleteTagResultsItem
func (o *DeleteTagResultsItem) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Provider != nil {
		m["provider"] = o.Provider
	}
	if o.IsError != nil {
		m["is_error"] = o.IsError
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalDeleteTagResultsItem unmarshals an instance of DeleteTagResultsItem from the specified map of raw messages.
func UnmarshalDeleteTagResultsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteTagResultsItem)
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		err = core.SDKErrorf(err, "", "provider-error", common.GetComponentInfo())
		return
	}
	delete(m, "provider")
	err = core.UnmarshalPrimitive(m, "is_error", &obj.IsError)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_error-error", common.GetComponentInfo())
		return
	}
	delete(m, "is_error")
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteTagsResult : Results of deleting unattatched tags.
type DeleteTagsResult struct {
	// The number of tags that have been deleted.
	TotalCount *int64 `json:"total_count,omitempty"`

	// It is set to true if there is at least one tag operation in error.
	Errors *bool `json:"errors,omitempty"`

	// The list of tag operation results.
	Items []DeleteTagsResultItem `json:"items,omitempty"`
}

// UnmarshalDeleteTagsResult unmarshals an instance of DeleteTagsResult from the specified map of raw messages.
func UnmarshalDeleteTagsResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteTagsResult)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "items", &obj.Items, UnmarshalDeleteTagsResultItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "items-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteTagsResultItem : Result of a delete_tags request.
type DeleteTagsResultItem struct {
	// The name of the deleted tag.
	TagName *string `json:"tag_name,omitempty"`

	// true if the tag was not deleted.
	IsError *bool `json:"is_error,omitempty"`
}

// UnmarshalDeleteTagsResultItem unmarshals an instance of DeleteTagsResultItem from the specified map of raw messages.
func UnmarshalDeleteTagsResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteTagsResultItem)
	err = core.UnmarshalPrimitive(m, "tag_name", &obj.TagName)
	if err != nil {
		err = core.SDKErrorf(err, "", "tag_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_error", &obj.IsError)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_error-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DetachTagOptions : The DetachTag options.
type DetachTagOptions struct {
	// The name of the tag to detach.
	TagName *string `json:"tag_name,omitempty"`

	// An array of tag names to detach.
	TagNames []string `json:"tag_names,omitempty"`

	// List of resources on which the tagging operation operates on.
	Resources []Resource `json:"resources,omitempty"`

	// A valid Global Search string.
	Query *QueryString `json:"query,omitempty"`

	// An alphanumeric string that is used to trace the request. The value  may include ASCII alphanumerics and any of
	// following segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024
	// bytes. The value is considered invalid and must be ignored if that value includes any other character or is longer
	// than 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XRequestID *string `json:"x-request-id,omitempty"`

	// An alphanumeric string that is used to trace the request as a part of a larger context: the same value is used for
	// downstream requests and retries of those requests. The value may include ASCII alphanumerics and any of following
	// segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024 bytes.
	// The value is considered invalid and must be ignored if that value includes any other character or is longer than
	// 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XCorrelationID *string `json:"x-correlation-id,omitempty"`

	// The ID of the billing account of the untagged resource. It is a required parameter if `tag_type` is set to
	// `service`, otherwise it is inferred from the authorization IAM token.
	AccountID *string `json:"account_id,omitempty"`

	// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
	// for IMS resources.
	TagType *string `json:"tag_type,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the DetachTagOptions.TagType property.
// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
// for IMS resources.
const (
	DetachTagOptionsTagTypeAccessConst  = "access"
	DetachTagOptionsTagTypeServiceConst = "service"
	DetachTagOptionsTagTypeUserConst    = "user"
)

// NewDetachTagOptions : Instantiate DetachTagOptions
func (*GlobalTaggingV1) NewDetachTagOptions() *DetachTagOptions {
	return &DetachTagOptions{}
}

// SetTagName : Allow user to set TagName
func (_options *DetachTagOptions) SetTagName(tagName string) *DetachTagOptions {
	_options.TagName = core.StringPtr(tagName)
	return _options
}

// SetTagNames : Allow user to set TagNames
func (_options *DetachTagOptions) SetTagNames(tagNames []string) *DetachTagOptions {
	_options.TagNames = tagNames
	return _options
}

// SetResources : Allow user to set Resources
func (_options *DetachTagOptions) SetResources(resources []Resource) *DetachTagOptions {
	_options.Resources = resources
	return _options
}

// SetQuery : Allow user to set Query
func (_options *DetachTagOptions) SetQuery(query *QueryString) *DetachTagOptions {
	_options.Query = query
	return _options
}

// SetXRequestID : Allow user to set XRequestID
func (_options *DetachTagOptions) SetXRequestID(xRequestID string) *DetachTagOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DetachTagOptions) SetXCorrelationID(xCorrelationID string) *DetachTagOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DetachTagOptions) SetAccountID(accountID string) *DetachTagOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTagType : Allow user to set TagType
func (_options *DetachTagOptions) SetTagType(tagType string) *DetachTagOptions {
	_options.TagType = core.StringPtr(tagType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DetachTagOptions) SetHeaders(param map[string]string) *DetachTagOptions {
	options.Headers = param
	return options
}

// ListTagsOptions : The ListTags options.
type ListTagsOptions struct {
	// An alphanumeric string that is used to trace the request. The value  may include ASCII alphanumerics and any of
	// following segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024
	// bytes. The value is considered invalid and must be ignored if that value includes any other character or is longer
	// than 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XRequestID *string `json:"x-request-id,omitempty"`

	// An alphanumeric string that is used to trace the request as a part of a larger context: the same value is used for
	// downstream requests and retries of those requests. The value may include ASCII alphanumerics and any of following
	// segment separators: space ( ), comma (,), hyphen, (-), and underscore (_) and may have a length up to 1024 bytes.
	// The value is considered invalid and must be ignored if that value includes any other character or is longer than
	// 1024 bytes or is fewer than 8 characters. If not specified or invalid, it is automatically replaced by a random
	// (version 4) UUID.
	XCorrelationID *string `json:"x-correlation-id,omitempty"`

	// The ID of the billing account to list the tags for. If it is not set, then it is taken from the authorization token.
	// This parameter is required if `tag_type` is set to `service`.
	AccountID *string `json:"account_id,omitempty"`

	// The type of the tag you want to list. Supported values are `user`, `service` and `access`.
	TagType *string `json:"tag_type,omitempty"`

	// If set to `true`, this query returns the provider, `ghost`, `ims` or `ghost,ims`, where the tag exists and the
	// number of attached resources.
	FullData *bool `json:"full_data,omitempty"`

	// Select a provider. Supported values are `ghost` and `ims`. To list both Global Search and Tagging tags and
	// infrastructure tags, use `ghost,ims`. `service` and `access` tags can only be attached to resources that are
	// onboarded to Global Search and Tagging, so you should not set this parameter to list them.
	Providers []string `json:"providers,omitempty"`

	// If you want to return only the list of tags that are attached to a specified resource, pass the ID of the resource
	// on this parameter. For resources that are onboarded to Global Search and Tagging, the resource ID is the CRN; for
	// IMS resources, it is the IMS ID. When using this parameter, you must specify the appropriate provider (`ims` or
	// `ghost`).
	AttachedTo *string `json:"attached_to,omitempty"`

	// The offset is the index of the item from which you want to start returning data from.
	Offset *int64 `json:"offset,omitempty"`

	// The number of tags to return.
	Limit *int64 `json:"limit,omitempty"`

	// The timeout in milliseconds, bounds the request to run within the specified time value. It returns the accumulated
	// results until time runs out.
	Timeout *int64 `json:"timeout,omitempty"`

	// Order the output by tag name.
	OrderByName *string `json:"order_by_name,omitempty"`

	// Filter on attached tags. If `true`, it returns only tags that are attached to one or more resources. If `false`, it
	// returns all tags.
	AttachedOnly *bool `json:"attached_only,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListTagsOptions.TagType property.
// The type of the tag you want to list. Supported values are `user`, `service` and `access`.
const (
	ListTagsOptionsTagTypeAccessConst  = "access"
	ListTagsOptionsTagTypeServiceConst = "service"
	ListTagsOptionsTagTypeUserConst    = "user"
)

// Constants associated with the ListTagsOptions.Providers property.
const (
	ListTagsOptionsProvidersGhostConst = "ghost"
	ListTagsOptionsProvidersImsConst   = "ims"
)

// Constants associated with the ListTagsOptions.OrderByName property.
// Order the output by tag name.
const (
	ListTagsOptionsOrderByNameAscConst  = "asc"
	ListTagsOptionsOrderByNameDescConst = "desc"
)

// NewListTagsOptions : Instantiate ListTagsOptions
func (*GlobalTaggingV1) NewListTagsOptions() *ListTagsOptions {
	return &ListTagsOptions{}
}

// SetXRequestID : Allow user to set XRequestID
func (_options *ListTagsOptions) SetXRequestID(xRequestID string) *ListTagsOptions {
	_options.XRequestID = core.StringPtr(xRequestID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListTagsOptions) SetXCorrelationID(xCorrelationID string) *ListTagsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListTagsOptions) SetAccountID(accountID string) *ListTagsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTagType : Allow user to set TagType
func (_options *ListTagsOptions) SetTagType(tagType string) *ListTagsOptions {
	_options.TagType = core.StringPtr(tagType)
	return _options
}

// SetFullData : Allow user to set FullData
func (_options *ListTagsOptions) SetFullData(fullData bool) *ListTagsOptions {
	_options.FullData = core.BoolPtr(fullData)
	return _options
}

// SetProviders : Allow user to set Providers
func (_options *ListTagsOptions) SetProviders(providers []string) *ListTagsOptions {
	_options.Providers = providers
	return _options
}

// SetAttachedTo : Allow user to set AttachedTo
func (_options *ListTagsOptions) SetAttachedTo(attachedTo string) *ListTagsOptions {
	_options.AttachedTo = core.StringPtr(attachedTo)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListTagsOptions) SetOffset(offset int64) *ListTagsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListTagsOptions) SetLimit(limit int64) *ListTagsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetTimeout : Allow user to set Timeout
func (_options *ListTagsOptions) SetTimeout(timeout int64) *ListTagsOptions {
	_options.Timeout = core.Int64Ptr(timeout)
	return _options
}

// SetOrderByName : Allow user to set OrderByName
func (_options *ListTagsOptions) SetOrderByName(orderByName string) *ListTagsOptions {
	_options.OrderByName = core.StringPtr(orderByName)
	return _options
}

// SetAttachedOnly : Allow user to set AttachedOnly
func (_options *ListTagsOptions) SetAttachedOnly(attachedOnly bool) *ListTagsOptions {
	_options.AttachedOnly = core.BoolPtr(attachedOnly)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTagsOptions) SetHeaders(param map[string]string) *ListTagsOptions {
	options.Headers = param
	return options
}

// QueryString : A valid Global Search string.
type QueryString struct {
	// The Lucene-formatted query string.
	QueryString *string `json:"query_string" validate:"required"`
}

// NewQueryString : Instantiate QueryString (Generic Model Constructor)
func (*GlobalTaggingV1) NewQueryString(queryString string) (_model *QueryString, err error) {
	_model = &QueryString{
		QueryString: core.StringPtr(queryString),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalQueryString unmarshals an instance of QueryString from the specified map of raw messages.
func UnmarshalQueryString(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueryString)
	err = core.UnmarshalPrimitive(m, "query_string", &obj.QueryString)
	if err != nil {
		err = core.SDKErrorf(err, "", "query_string-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Resource : A resource that might have tags that are attached.
type Resource struct {
	// The CRN or IMS ID of the resource.
	ResourceID *string `json:"resource_id" validate:"required"`

	// The IMS resource type of the resource. It can be one of SoftLayer_Virtual_DedicatedHost, SoftLayer_Hardware,
	// SoftLayer_Hardware_Server, SoftLayer_Network_Application_Delivery_Controller, SoftLayer_Network_Vlan,
	// SoftLayer_Network_Vlan_Firewall, SoftLayer_Network_Component_Firewall, SoftLayer_Network_Firewall_Module_Context,
	// SoftLayer_Virtual_Guest.
	ResourceType *string `json:"resource_type,omitempty"`
}

// NewResource : Instantiate Resource (Generic Model Constructor)
func (*GlobalTaggingV1) NewResource(resourceID string) (_model *Resource, err error) {
	_model = &Resource{
		ResourceID: core.StringPtr(resourceID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalResource unmarshals an instance of Resource from the specified map of raw messages.
func UnmarshalResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Resource)
	err = core.UnmarshalPrimitive(m, "resource_id", &obj.ResourceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tag : A tag.
type Tag struct {
	// The name of the tag.
	Name *string `json:"name" validate:"required"`
}

// UnmarshalTag unmarshals an instance of Tag from the specified map of raw messages.
func UnmarshalTag(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tag)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TagList : A list of tags.
type TagList struct {
	// Set the occurrences of the total tags that are associated with this account.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The offset at which tags are returned.
	Offset *int64 `json:"offset,omitempty"`

	// The number of tags requested to be returned.
	Limit *int64 `json:"limit,omitempty"`

	// Array of output results.
	Items []Tag `json:"items,omitempty"`
}

// UnmarshalTagList unmarshals an instance of TagList from the specified map of raw messages.
func UnmarshalTagList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TagList)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalModel(m, "items", &obj.Items, UnmarshalTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "items-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TagResults : Results of an attach_tag or detach_tag request.
type TagResults struct {
	// Array of results of an attach_tag or detach_tag request.
	Results []TagResultsItem `json:"results,omitempty"`
}

// UnmarshalTagResults unmarshals an instance of TagResults from the specified map of raw messages.
func UnmarshalTagResults(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TagResults)
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalTagResultsItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "results-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TagResultsItem : Result of an attach_tag or detach_tag request for a tagged resource.
type TagResultsItem struct {
	// The CRN or IMS ID of the resource.
	ResourceID *string `json:"resource_id" validate:"required"`

	// It is `true` if the operation exits with an error.
	IsError *bool `json:"is_error,omitempty"`
}

// UnmarshalTagResultsItem unmarshals an instance of TagResultsItem from the specified map of raw messages.
func UnmarshalTagResultsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TagResultsItem)
	err = core.UnmarshalPrimitive(m, "resource_id", &obj.ResourceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_error", &obj.IsError)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_error-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
