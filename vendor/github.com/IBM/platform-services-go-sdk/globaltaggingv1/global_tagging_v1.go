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
 * IBM OpenAPI SDK Code Generator Version: 99-SNAPSHOT-a8493a65-20210115-083246
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

// GlobalTaggingV1 : Manage your tags with the Tagging API in IBM Cloud. You can attach, detach, delete a tag or list
// all tags in your billing account with the Tagging API. The tag name must be unique within a billing account. You can
// create tags in two formats: `key:value` or `label`. The tagging API supports three types of tag: `user` `service`,
// and `access` tags. `service` tags cannot be attached to IMS resources. `service` tags must be in the form
// `service_prefix:tag_label` where `service_prefix` identifies the Service owning the tag. `access` tags cannot be
// attached to IMS and Cloud Foundry resources. They must be in the form `key:value`.
//
// Version: 1.2.0
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
			return
		}
	}

	globalTagging, err = NewGlobalTaggingV1(options)
	if err != nil {
		return
	}

	err = globalTagging.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = globalTagging.Service.SetServiceURL(options.URL)
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
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
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
	return "", fmt.Errorf("service does not support regional URLs")
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
	return globalTagging.Service.SetServiceURL(url)
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
// Lists all tags in a billing account. Use the `attached_to` parameter to return the list of tags attached to the
// specified resource.
func (globalTagging *GlobalTaggingV1) ListTags(listTagsOptions *ListTagsOptions) (result *TagList, response *core.DetailedResponse, err error) {
	return globalTagging.ListTagsWithContext(context.Background(), listTagsOptions)
}

// ListTagsWithContext is an alternate form of the ListTags method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) ListTagsWithContext(ctx context.Context, listTagsOptions *ListTagsOptions) (result *TagList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listTagsOptions, "listTagsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags`, nil)
	if err != nil {
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

	if listTagsOptions.ImpersonateUser != nil {
		builder.AddQuery("impersonate_user", fmt.Sprint(*listTagsOptions.ImpersonateUser))
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTagList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateTag : Create an access tag
// Create an access tag. To create an `access` tag, you must have the access listed in the [Granting users access to tag
// resources](https://cloud.ibm.com/docs/account?topic=account-access) documentation. `service` and `user` tags cannot
// be created upfront. They are created when they are attached for the first time to a resource.
func (globalTagging *GlobalTaggingV1) CreateTag(createTagOptions *CreateTagOptions) (result *CreateTagResults, response *core.DetailedResponse, err error) {
	return globalTagging.CreateTagWithContext(context.Background(), createTagOptions)
}

// CreateTagWithContext is an alternate form of the CreateTag method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) CreateTagWithContext(ctx context.Context, createTagOptions *CreateTagOptions) (result *CreateTagResults, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTagOptions, "createTagOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTagOptions, "createTagOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags`, nil)
	if err != nil {
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

	if createTagOptions.ImpersonateUser != nil {
		builder.AddQuery("impersonate_user", fmt.Sprint(*createTagOptions.ImpersonateUser))
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateTagResults)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteTagAll : Delete all unused tags
// Delete the tags that are not attached to any resource.
func (globalTagging *GlobalTaggingV1) DeleteTagAll(deleteTagAllOptions *DeleteTagAllOptions) (result *DeleteTagsResult, response *core.DetailedResponse, err error) {
	return globalTagging.DeleteTagAllWithContext(context.Background(), deleteTagAllOptions)
}

// DeleteTagAllWithContext is an alternate form of the DeleteTagAll method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) DeleteTagAllWithContext(ctx context.Context, deleteTagAllOptions *DeleteTagAllOptions) (result *DeleteTagsResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(deleteTagAllOptions, "deleteTagAllOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags`, nil)
	if err != nil {
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

	if deleteTagAllOptions.Providers != nil {
		builder.AddQuery("providers", fmt.Sprint(*deleteTagAllOptions.Providers))
	}
	if deleteTagAllOptions.ImpersonateUser != nil {
		builder.AddQuery("impersonate_user", fmt.Sprint(*deleteTagAllOptions.ImpersonateUser))
	}
	if deleteTagAllOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteTagAllOptions.AccountID))
	}
	if deleteTagAllOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*deleteTagAllOptions.TagType))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteTagsResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteTag : Delete an unused tag
// Delete an existing tag. A tag can be deleted only if it is not attached to any resource.
func (globalTagging *GlobalTaggingV1) DeleteTag(deleteTagOptions *DeleteTagOptions) (result *DeleteTagResults, response *core.DetailedResponse, err error) {
	return globalTagging.DeleteTagWithContext(context.Background(), deleteTagOptions)
}

// DeleteTagWithContext is an alternate form of the DeleteTag method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) DeleteTagWithContext(ctx context.Context, deleteTagOptions *DeleteTagOptions) (result *DeleteTagResults, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTagOptions, "deleteTagOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTagOptions, "deleteTagOptions")
	if err != nil {
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

	if deleteTagOptions.Providers != nil {
		builder.AddQuery("providers", strings.Join(deleteTagOptions.Providers, ","))
	}
	if deleteTagOptions.ImpersonateUser != nil {
		builder.AddQuery("impersonate_user", fmt.Sprint(*deleteTagOptions.ImpersonateUser))
	}
	if deleteTagOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteTagOptions.AccountID))
	}
	if deleteTagOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*deleteTagOptions.TagType))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteTagResults)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// AttachTag : Attach tags
// Attaches one or more tags to one or more resources. To attach a `user` tag on a resource, you must have the access
// listed in the [Granting users access to tag resources](https://cloud.ibm.com/docs/account?topic=account-access)
// documentation. To attach a `service` tag, you must be an authorized service. If that is the case, then you can attach
// a `service` tag with your registered `prefix` to any resource in any account. The account ID must be set through the
// `account_id` query parameter. To attach an `access` tag, you must be the resource administrator within the account.
// You can attach only `access` tags already existing.
func (globalTagging *GlobalTaggingV1) AttachTag(attachTagOptions *AttachTagOptions) (result *TagResults, response *core.DetailedResponse, err error) {
	return globalTagging.AttachTagWithContext(context.Background(), attachTagOptions)
}

// AttachTagWithContext is an alternate form of the AttachTag method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) AttachTagWithContext(ctx context.Context, attachTagOptions *AttachTagOptions) (result *TagResults, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(attachTagOptions, "attachTagOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(attachTagOptions, "attachTagOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags/attach`, nil)
	if err != nil {
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

	if attachTagOptions.ImpersonateUser != nil {
		builder.AddQuery("impersonate_user", fmt.Sprint(*attachTagOptions.ImpersonateUser))
	}
	if attachTagOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*attachTagOptions.AccountID))
	}
	if attachTagOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*attachTagOptions.TagType))
	}

	body := make(map[string]interface{})
	if attachTagOptions.Resources != nil {
		body["resources"] = attachTagOptions.Resources
	}
	if attachTagOptions.TagName != nil {
		body["tag_name"] = attachTagOptions.TagName
	}
	if attachTagOptions.TagNames != nil {
		body["tag_names"] = attachTagOptions.TagNames
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
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTagResults)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DetachTag : Detach tags
// Detaches one or more tags from one or more resources. To detach a `user` tag on a resource you must have the
// permissions listed in the [Granting users access to tag
// resources](https://cloud.ibm.com/docs/account?topic=account-access) documentation. To detach a `service` tag you must
// be an authorized Service. If that is the case, then you can detach a `service` tag with your registered `prefix` from
// any resource in any account. The account ID must be set through the `account_id` query parameter. To detach an
// `access` tag, you must be the resource administrator within the account.
func (globalTagging *GlobalTaggingV1) DetachTag(detachTagOptions *DetachTagOptions) (result *TagResults, response *core.DetailedResponse, err error) {
	return globalTagging.DetachTagWithContext(context.Background(), detachTagOptions)
}

// DetachTagWithContext is an alternate form of the DetachTag method which supports a Context parameter
func (globalTagging *GlobalTaggingV1) DetachTagWithContext(ctx context.Context, detachTagOptions *DetachTagOptions) (result *TagResults, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(detachTagOptions, "detachTagOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(detachTagOptions, "detachTagOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalTagging.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalTagging.Service.Options.URL, `/v3/tags/detach`, nil)
	if err != nil {
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

	if detachTagOptions.ImpersonateUser != nil {
		builder.AddQuery("impersonate_user", fmt.Sprint(*detachTagOptions.ImpersonateUser))
	}
	if detachTagOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*detachTagOptions.AccountID))
	}
	if detachTagOptions.TagType != nil {
		builder.AddQuery("tag_type", fmt.Sprint(*detachTagOptions.TagType))
	}

	body := make(map[string]interface{})
	if detachTagOptions.Resources != nil {
		body["resources"] = detachTagOptions.Resources
	}
	if detachTagOptions.TagName != nil {
		body["tag_name"] = detachTagOptions.TagName
	}
	if detachTagOptions.TagNames != nil {
		body["tag_names"] = detachTagOptions.TagNames
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
	response, err = globalTagging.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTagResults)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// AttachTagOptions : The AttachTag options.
type AttachTagOptions struct {
	// List of resources on which the tag or tags should be attached.
	Resources []Resource `validate:"required"`

	// The name of the tag to attach.
	TagName *string

	// An array of tag names to attach.
	TagNames []string

	// The user on whose behalf the attach operation must be performed (_for administrators only_).
	ImpersonateUser *string

	// The ID of the billing account where the resources to be tagged lives. It is a required parameter if `tag_type` is
	// set to `service`. Otherwise, it is inferred from the authorization IAM token.
	AccountID *string

	// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
	// for IMS resources.
	TagType *string

	// Allows users to set headers on API requests
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
func (*GlobalTaggingV1) NewAttachTagOptions(resources []Resource) *AttachTagOptions {
	return &AttachTagOptions{
		Resources: resources,
	}
}

// SetResources : Allow user to set Resources
func (options *AttachTagOptions) SetResources(resources []Resource) *AttachTagOptions {
	options.Resources = resources
	return options
}

// SetTagName : Allow user to set TagName
func (options *AttachTagOptions) SetTagName(tagName string) *AttachTagOptions {
	options.TagName = core.StringPtr(tagName)
	return options
}

// SetTagNames : Allow user to set TagNames
func (options *AttachTagOptions) SetTagNames(tagNames []string) *AttachTagOptions {
	options.TagNames = tagNames
	return options
}

// SetImpersonateUser : Allow user to set ImpersonateUser
func (options *AttachTagOptions) SetImpersonateUser(impersonateUser string) *AttachTagOptions {
	options.ImpersonateUser = core.StringPtr(impersonateUser)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *AttachTagOptions) SetAccountID(accountID string) *AttachTagOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetTagType : Allow user to set TagType
func (options *AttachTagOptions) SetTagType(tagType string) *AttachTagOptions {
	options.TagType = core.StringPtr(tagType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *AttachTagOptions) SetHeaders(param map[string]string) *AttachTagOptions {
	options.Headers = param
	return options
}

// CreateTagOptions : The CreateTag options.
type CreateTagOptions struct {
	// An array of tag names to create.
	TagNames []string `validate:"required"`

	// The user on whose behalf the create operation must be performed (_for administrators only_).
	ImpersonateUser *string

	// The ID of the billing account where the tag must be created. It is a required parameter if `impersonate_user` is
	// set.
	AccountID *string

	// The type of the tags you want to create. The only allowed value is `access`.
	TagType *string

	// Allows users to set headers on API requests
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
func (options *CreateTagOptions) SetTagNames(tagNames []string) *CreateTagOptions {
	options.TagNames = tagNames
	return options
}

// SetImpersonateUser : Allow user to set ImpersonateUser
func (options *CreateTagOptions) SetImpersonateUser(impersonateUser string) *CreateTagOptions {
	options.ImpersonateUser = core.StringPtr(impersonateUser)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *CreateTagOptions) SetAccountID(accountID string) *CreateTagOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetTagType : Allow user to set TagType
func (options *CreateTagOptions) SetTagType(tagType string) *CreateTagOptions {
	options.TagType = core.StringPtr(tagType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTagOptions) SetHeaders(param map[string]string) *CreateTagOptions {
	options.Headers = param
	return options
}

// CreateTagResults : Results of a create tag(s) request.
type CreateTagResults struct {
	// Array of results of an set_tags request.
	Results []CreateTagResultsResultsItem `json:"results,omitempty"`
}

// UnmarshalCreateTagResults unmarshals an instance of CreateTagResults from the specified map of raw messages.
func UnmarshalCreateTagResults(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateTagResults)
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalCreateTagResultsResultsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateTagResultsResultsItem : CreateTagResultsResultsItem struct
type CreateTagResultsResultsItem struct {
	// The name of the tag created.
	TagName *string `json:"tag_name,omitempty"`

	// true if the tag was not created.
	IsError *bool `json:"is_error,omitempty"`
}

// UnmarshalCreateTagResultsResultsItem unmarshals an instance of CreateTagResultsResultsItem from the specified map of raw messages.
func UnmarshalCreateTagResultsResultsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateTagResultsResultsItem)
	err = core.UnmarshalPrimitive(m, "tag_name", &obj.TagName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_error", &obj.IsError)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteTagAllOptions : The DeleteTagAll options.
type DeleteTagAllOptions struct {
	// Select a provider. Supported values are `ghost` and `ims`.
	Providers *string

	// The user on whose behalf the delete all operation must be performed (_for administrators only_).
	ImpersonateUser *string

	// The ID of the billing account to delete the tags for. If it is not set, then it is taken from the authorization
	// token. It is a required parameter if `tag_type` is set to `service`.
	AccountID *string

	// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
	// for IMS resources (`providers` parameter set to `ims`).
	TagType *string

	// Allows users to set headers on API requests
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

// SetProviders : Allow user to set Providers
func (options *DeleteTagAllOptions) SetProviders(providers string) *DeleteTagAllOptions {
	options.Providers = core.StringPtr(providers)
	return options
}

// SetImpersonateUser : Allow user to set ImpersonateUser
func (options *DeleteTagAllOptions) SetImpersonateUser(impersonateUser string) *DeleteTagAllOptions {
	options.ImpersonateUser = core.StringPtr(impersonateUser)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *DeleteTagAllOptions) SetAccountID(accountID string) *DeleteTagAllOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetTagType : Allow user to set TagType
func (options *DeleteTagAllOptions) SetTagType(tagType string) *DeleteTagAllOptions {
	options.TagType = core.StringPtr(tagType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTagAllOptions) SetHeaders(param map[string]string) *DeleteTagAllOptions {
	options.Headers = param
	return options
}

// DeleteTagOptions : The DeleteTag options.
type DeleteTagOptions struct {
	// The name of tag to be deleted.
	TagName *string `validate:"required,ne="`

	// Select a provider. Supported values are `ghost` and `ims`. To delete tag both in GhoST in IMS, use `ghost,ims`.
	Providers []string

	// The user on whose behalf the delete operation must be performed (_for administrators only_).
	ImpersonateUser *string

	// The ID of the billing account to delete the tag for. It is a required parameter if `tag_type` is set to `service`,
	// otherwise it is inferred from the authorization IAM token.
	AccountID *string

	// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
	// for IMS resources (`providers` parameter set to `ims`).
	TagType *string

	// Allows users to set headers on API requests
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
func (options *DeleteTagOptions) SetTagName(tagName string) *DeleteTagOptions {
	options.TagName = core.StringPtr(tagName)
	return options
}

// SetProviders : Allow user to set Providers
func (options *DeleteTagOptions) SetProviders(providers []string) *DeleteTagOptions {
	options.Providers = providers
	return options
}

// SetImpersonateUser : Allow user to set ImpersonateUser
func (options *DeleteTagOptions) SetImpersonateUser(impersonateUser string) *DeleteTagOptions {
	options.ImpersonateUser = core.StringPtr(impersonateUser)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *DeleteTagOptions) SetAccountID(accountID string) *DeleteTagOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetTagType : Allow user to set TagType
func (options *DeleteTagOptions) SetTagType(tagType string) *DeleteTagOptions {
	options.TagType = core.StringPtr(tagType)
	return options
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
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteTagResultsItem : Result of a delete_tag request.
type DeleteTagResultsItem struct {
	// The provider of the tag.
	Provider *string `json:"provider,omitempty"`

	// It is `true` if the operation exits with an error.
	IsError *bool `json:"is_error,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// Constants associated with the DeleteTagResultsItem.Provider property.
// The provider of the tag.
const (
	DeleteTagResultsItemProviderGhostConst = "ghost"
	DeleteTagResultsItemProviderImsConst   = "ims"
)

// SetProperty allows the user to set an arbitrary property on an instance of DeleteTagResultsItem
func (o *DeleteTagResultsItem) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of DeleteTagResultsItem
func (o *DeleteTagResultsItem) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of DeleteTagResultsItem
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
	return
}

// UnmarshalDeleteTagResultsItem unmarshals an instance of DeleteTagResultsItem from the specified map of raw messages.
func UnmarshalDeleteTagResultsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteTagResultsItem)
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		return
	}
	delete(m, "provider")
	err = core.UnmarshalPrimitive(m, "is_error", &obj.IsError)
	if err != nil {
		return
	}
	delete(m, "is_error")
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

// DeleteTagsResult : Results of a deleting unattatched tags.
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
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "items", &obj.Items, UnmarshalDeleteTagsResultItem)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "is_error", &obj.IsError)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DetachTagOptions : The DetachTag options.
type DetachTagOptions struct {
	// List of resources on which the tag or tags should be detached.
	Resources []Resource `validate:"required"`

	// The name of the tag to detach.
	TagName *string

	// An array of tag names to detach.
	TagNames []string

	// The user on whose behalf the detach operation must be performed (_for administrators only_).
	ImpersonateUser *string

	// The ID of the billing account where the resources to be un-tagged lives. It is a required parameter if `tag_type` is
	// set to `service`, otherwise it is inferred from the authorization IAM token.
	AccountID *string

	// The type of the tag. Supported values are `user`, `service` and `access`. `service` and `access` are not supported
	// for IMS resources.
	TagType *string

	// Allows users to set headers on API requests
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
func (*GlobalTaggingV1) NewDetachTagOptions(resources []Resource) *DetachTagOptions {
	return &DetachTagOptions{
		Resources: resources,
	}
}

// SetResources : Allow user to set Resources
func (options *DetachTagOptions) SetResources(resources []Resource) *DetachTagOptions {
	options.Resources = resources
	return options
}

// SetTagName : Allow user to set TagName
func (options *DetachTagOptions) SetTagName(tagName string) *DetachTagOptions {
	options.TagName = core.StringPtr(tagName)
	return options
}

// SetTagNames : Allow user to set TagNames
func (options *DetachTagOptions) SetTagNames(tagNames []string) *DetachTagOptions {
	options.TagNames = tagNames
	return options
}

// SetImpersonateUser : Allow user to set ImpersonateUser
func (options *DetachTagOptions) SetImpersonateUser(impersonateUser string) *DetachTagOptions {
	options.ImpersonateUser = core.StringPtr(impersonateUser)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *DetachTagOptions) SetAccountID(accountID string) *DetachTagOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetTagType : Allow user to set TagType
func (options *DetachTagOptions) SetTagType(tagType string) *DetachTagOptions {
	options.TagType = core.StringPtr(tagType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DetachTagOptions) SetHeaders(param map[string]string) *DetachTagOptions {
	options.Headers = param
	return options
}

// ListTagsOptions : The ListTags options.
type ListTagsOptions struct {
	// The user on whose behalf the get operation must be performed (_for administrators only_).
	ImpersonateUser *string

	// The ID of the billing account to list the tags for. If it is not set, then it is taken from the authorization token.
	// This parameter is required if `tag_type` is set to `service`.
	AccountID *string

	// The type of the tag you want to list. Supported values are `user`, `service` and `access`.
	TagType *string

	// If set to `true`, this query returns the provider, `ghost`, `ims` or `ghost,ims`, where the tag exists and the
	// number of attached resources.
	FullData *bool

	// Select a provider. Supported values are `ghost` and `ims`. To list GhoST tags and infrastructure tags use
	// `ghost,ims`. `service` and `access` tags can only be attached to GhoST onboarded resources, so you should not set
	// this parameter when listing them.
	Providers []string

	// If you want to return only the list of tags attached to a specified resource, pass the ID of the resource on this
	// parameter. For GhoST onboarded resources, the resource ID is the CRN; for IMS resources, it is the IMS ID. When
	// using this parameter, you must specify the appropriate provider (`ims` or `ghost`).
	AttachedTo *string

	// The offset is the index of the item from which you want to start returning data from.
	Offset *int64

	// The number of tags to return.
	Limit *int64

	// The search timeout bounds the search request to be executed within the specified time value. It returns the hits
	// accumulated until time runs out.
	Timeout *int64

	// Order the output by tag name.
	OrderByName *string

	// Filter on attached tags. If `true`, it returns only tags that are attached to one or more resources. If `false`, it
	// returns all tags.
	AttachedOnly *bool

	// Allows users to set headers on API requests
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

// SetImpersonateUser : Allow user to set ImpersonateUser
func (options *ListTagsOptions) SetImpersonateUser(impersonateUser string) *ListTagsOptions {
	options.ImpersonateUser = core.StringPtr(impersonateUser)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *ListTagsOptions) SetAccountID(accountID string) *ListTagsOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetTagType : Allow user to set TagType
func (options *ListTagsOptions) SetTagType(tagType string) *ListTagsOptions {
	options.TagType = core.StringPtr(tagType)
	return options
}

// SetFullData : Allow user to set FullData
func (options *ListTagsOptions) SetFullData(fullData bool) *ListTagsOptions {
	options.FullData = core.BoolPtr(fullData)
	return options
}

// SetProviders : Allow user to set Providers
func (options *ListTagsOptions) SetProviders(providers []string) *ListTagsOptions {
	options.Providers = providers
	return options
}

// SetAttachedTo : Allow user to set AttachedTo
func (options *ListTagsOptions) SetAttachedTo(attachedTo string) *ListTagsOptions {
	options.AttachedTo = core.StringPtr(attachedTo)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListTagsOptions) SetOffset(offset int64) *ListTagsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListTagsOptions) SetLimit(limit int64) *ListTagsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetTimeout : Allow user to set Timeout
func (options *ListTagsOptions) SetTimeout(timeout int64) *ListTagsOptions {
	options.Timeout = core.Int64Ptr(timeout)
	return options
}

// SetOrderByName : Allow user to set OrderByName
func (options *ListTagsOptions) SetOrderByName(orderByName string) *ListTagsOptions {
	options.OrderByName = core.StringPtr(orderByName)
	return options
}

// SetAttachedOnly : Allow user to set AttachedOnly
func (options *ListTagsOptions) SetAttachedOnly(attachedOnly bool) *ListTagsOptions {
	options.AttachedOnly = core.BoolPtr(attachedOnly)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListTagsOptions) SetHeaders(param map[string]string) *ListTagsOptions {
	options.Headers = param
	return options
}

// Resource : A resource that may have attached tags.
type Resource struct {
	// The CRN or IMS ID of the resource.
	ResourceID *string `json:"resource_id" validate:"required"`

	// The IMS resource type of the resource.
	ResourceType *string `json:"resource_type,omitempty"`
}

// NewResource : Instantiate Resource (Generic Model Constructor)
func (*GlobalTaggingV1) NewResource(resourceID string) (model *Resource, err error) {
	model = &Resource{
		ResourceID: core.StringPtr(resourceID),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalResource unmarshals an instance of Resource from the specified map of raw messages.
func UnmarshalResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Resource)
	err = core.UnmarshalPrimitive(m, "resource_id", &obj.ResourceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tag : A tag.
type Tag struct {
	// This is the name of the tag.
	Name *string `json:"name" validate:"required"`
}

// UnmarshalTag unmarshals an instance of Tag from the specified map of raw messages.
func UnmarshalTag(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tag)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TagList : A list of tags.
type TagList struct {
	// Set the occurrencies of the total tags associated to this account.
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
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "items", &obj.Items, UnmarshalTag)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "is_error", &obj.IsError)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
