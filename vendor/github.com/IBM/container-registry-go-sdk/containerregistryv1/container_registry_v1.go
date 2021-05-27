/**
 * (C) Copyright IBM Corp. 2020, 2021.
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
 * IBM OpenAPI SDK Code Generator Version: 3.29.1-b338fb38-20210313-010605
 */

// Package containerregistryv1 : Operations and models for the ContainerRegistryV1 service
package containerregistryv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	common "github.com/IBM/container-registry-go-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
)

// ContainerRegistryV1 : Management interface for IBM Cloud Container Registry
//
// Version: 1.1
type ContainerRegistryV1 struct {
	Service *core.BaseService

	// The unique ID for your IBM Cloud account.
	Account *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us.icr.io"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "container_registry"

// ContainerRegistryV1Options : Service options
type ContainerRegistryV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// The unique ID for your IBM Cloud account.
	Account *string `validate:"required"`
}

// NewContainerRegistryV1UsingExternalConfig : constructs an instance of ContainerRegistryV1 with passed in options and external configuration.
func NewContainerRegistryV1UsingExternalConfig(options *ContainerRegistryV1Options) (containerRegistry *ContainerRegistryV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	containerRegistry, err = NewContainerRegistryV1(options)
	if err != nil {
		return
	}

	err = containerRegistry.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = containerRegistry.Service.SetServiceURL(options.URL)
	}
	return
}

// NewContainerRegistryV1 : constructs an instance of ContainerRegistryV1 with passed in options.
func NewContainerRegistryV1(options *ContainerRegistryV1Options) (service *ContainerRegistryV1, err error) {
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

	service = &ContainerRegistryV1{
		Service: baseService,
		Account: options.Account,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south":   "https://us.icr.io",  // us-south
		"uk-south":   "https://uk.icr.io",  // uk-south
		"eu-gb":      "https://uk.icr.io",  // eu-gb
		"eu-central": "https://de.icr.io",  // eu-central
		"eu-de":      "https://de.icr.io",  // eu-de
		"ap-north":   "https://jp.icr.io",  // ap-north
		"jp-tok":     "https://jp.icr.io",  // jp-tok
		"ap-south":   "https://au.icr.io",  // ap-south
		"au-syd":     "https://au.icr.io",  // au-syd
		"global":     "https://icr.io",     // global
		"jp-osa":     "https://jp2.icr.io", // jp-osa
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}

// Clone makes a copy of "containerRegistry" suitable for processing requests.
func (containerRegistry *ContainerRegistryV1) Clone() *ContainerRegistryV1 {
	if core.IsNil(containerRegistry) {
		return nil
	}
	clone := *containerRegistry
	clone.Service = containerRegistry.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (containerRegistry *ContainerRegistryV1) SetServiceURL(url string) error {
	return containerRegistry.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (containerRegistry *ContainerRegistryV1) GetServiceURL() string {
	return containerRegistry.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (containerRegistry *ContainerRegistryV1) SetDefaultHeaders(headers http.Header) {
	containerRegistry.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (containerRegistry *ContainerRegistryV1) SetEnableGzipCompression(enableGzip bool) {
	containerRegistry.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (containerRegistry *ContainerRegistryV1) GetEnableGzipCompression() bool {
	return containerRegistry.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (containerRegistry *ContainerRegistryV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	containerRegistry.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (containerRegistry *ContainerRegistryV1) DisableRetries() {
	containerRegistry.Service.DisableRetries()
}

// GetAuth : Get authorization options
// Get authorization options for the targeted account.
func (containerRegistry *ContainerRegistryV1) GetAuth(getAuthOptions *GetAuthOptions) (result *AuthOptions, response *core.DetailedResponse, err error) {
	return containerRegistry.GetAuthWithContext(context.Background(), getAuthOptions)
}

// GetAuthWithContext is an alternate form of the GetAuth method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) GetAuthWithContext(ctx context.Context, getAuthOptions *GetAuthOptions) (result *AuthOptions, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getAuthOptions, "getAuthOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/auth`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAuthOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "GetAuth")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAuthOptions)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateAuth : Update authorization options
// Update authorization options for the targeted account.
func (containerRegistry *ContainerRegistryV1) UpdateAuth(updateAuthOptions *UpdateAuthOptions) (response *core.DetailedResponse, err error) {
	return containerRegistry.UpdateAuthWithContext(context.Background(), updateAuthOptions)
}

// UpdateAuthWithContext is an alternate form of the UpdateAuth method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) UpdateAuthWithContext(ctx context.Context, updateAuthOptions *UpdateAuthOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAuthOptions, "updateAuthOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAuthOptions, "updateAuthOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/auth`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateAuthOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "UpdateAuth")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	body := make(map[string]interface{})
	if updateAuthOptions.IamAuthz != nil {
		body["iam_authz"] = updateAuthOptions.IamAuthz
	}
	if updateAuthOptions.PrivateOnly != nil {
		body["private_only"] = updateAuthOptions.PrivateOnly
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, nil)

	return
}

// ListImages : List images
// List all images in namespaces in a targeted IBM Cloud account.
func (containerRegistry *ContainerRegistryV1) ListImages(listImagesOptions *ListImagesOptions) (result []RemoteAPIImage, response *core.DetailedResponse, err error) {
	return containerRegistry.ListImagesWithContext(context.Background(), listImagesOptions)
}

// ListImagesWithContext is an alternate form of the ListImages method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) ListImagesWithContext(ctx context.Context, listImagesOptions *ListImagesOptions) (result []RemoteAPIImage, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listImagesOptions, "listImagesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/images`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listImagesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "ListImages")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	if listImagesOptions.Namespace != nil {
		builder.AddQuery("namespace", fmt.Sprint(*listImagesOptions.Namespace))
	}
	if listImagesOptions.IncludeIBM != nil {
		builder.AddQuery("includeIBM", fmt.Sprint(*listImagesOptions.IncludeIBM))
	}
	if listImagesOptions.IncludePrivate != nil {
		builder.AddQuery("includePrivate", fmt.Sprint(*listImagesOptions.IncludePrivate))
	}
	if listImagesOptions.IncludeManifestLists != nil {
		builder.AddQuery("includeManifestLists", fmt.Sprint(*listImagesOptions.IncludeManifestLists))
	}
	if listImagesOptions.Vulnerabilities != nil {
		builder.AddQuery("vulnerabilities", fmt.Sprint(*listImagesOptions.Vulnerabilities))
	}
	if listImagesOptions.Repository != nil {
		builder.AddQuery("repository", fmt.Sprint(*listImagesOptions.Repository))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRemoteAPIImage)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// BulkDeleteImages : Bulk delete images
// Remove multiple container images from the registry.
func (containerRegistry *ContainerRegistryV1) BulkDeleteImages(bulkDeleteImagesOptions *BulkDeleteImagesOptions) (result *ImageBulkDeleteResult, response *core.DetailedResponse, err error) {
	return containerRegistry.BulkDeleteImagesWithContext(context.Background(), bulkDeleteImagesOptions)
}

// BulkDeleteImagesWithContext is an alternate form of the BulkDeleteImages method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) BulkDeleteImagesWithContext(ctx context.Context, bulkDeleteImagesOptions *BulkDeleteImagesOptions) (result *ImageBulkDeleteResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(bulkDeleteImagesOptions, "bulkDeleteImagesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(bulkDeleteImagesOptions, "bulkDeleteImagesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/images/bulkdelete`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range bulkDeleteImagesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "BulkDeleteImages")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	_, err = builder.SetBodyContentJSON(bulkDeleteImagesOptions.BulkDelete)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImageBulkDeleteResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListImageDigests : List images by digest
// List all images by digest in namespaces in a targeted IBM Cloud account.
func (containerRegistry *ContainerRegistryV1) ListImageDigests(listImageDigestsOptions *ListImageDigestsOptions) (result []ImageDigest, response *core.DetailedResponse, err error) {
	return containerRegistry.ListImageDigestsWithContext(context.Background(), listImageDigestsOptions)
}

// ListImageDigestsWithContext is an alternate form of the ListImageDigests method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) ListImageDigestsWithContext(ctx context.Context, listImageDigestsOptions *ListImageDigestsOptions) (result []ImageDigest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listImageDigestsOptions, "listImageDigestsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listImageDigestsOptions, "listImageDigestsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/images/digests`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listImageDigestsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "ListImageDigests")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	body := make(map[string]interface{})
	if listImageDigestsOptions.ExcludeTagged != nil {
		body["exclude_tagged"] = listImageDigestsOptions.ExcludeTagged
	}
	if listImageDigestsOptions.ExcludeVa != nil {
		body["exclude_va"] = listImageDigestsOptions.ExcludeVa
	}
	if listImageDigestsOptions.IncludeIBM != nil {
		body["include_ibm"] = listImageDigestsOptions.IncludeIBM
	}
	if listImageDigestsOptions.Repositories != nil {
		body["repositories"] = listImageDigestsOptions.Repositories
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
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImageDigest)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// TagImage : Create tag
// Create a new tag in a private registry that refers to an existing image in the same region.
func (containerRegistry *ContainerRegistryV1) TagImage(tagImageOptions *TagImageOptions) (response *core.DetailedResponse, err error) {
	return containerRegistry.TagImageWithContext(context.Background(), tagImageOptions)
}

// TagImageWithContext is an alternate form of the TagImage method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) TagImageWithContext(ctx context.Context, tagImageOptions *TagImageOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(tagImageOptions, "tagImageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(tagImageOptions, "tagImageOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/images/tags`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range tagImageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "TagImage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	builder.AddQuery("fromimage", fmt.Sprint(*tagImageOptions.Fromimage))
	builder.AddQuery("toimage", fmt.Sprint(*tagImageOptions.Toimage))

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, nil)

	return
}

// DeleteImage : Delete image
// Delete a container image from the registry.
func (containerRegistry *ContainerRegistryV1) DeleteImage(deleteImageOptions *DeleteImageOptions) (result *ImageDeleteResult, response *core.DetailedResponse, err error) {
	return containerRegistry.DeleteImageWithContext(context.Background(), deleteImageOptions)
}

// DeleteImageWithContext is an alternate form of the DeleteImage method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) DeleteImageWithContext(ctx context.Context, deleteImageOptions *DeleteImageOptions) (result *ImageDeleteResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteImageOptions, "deleteImageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteImageOptions, "deleteImageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"image": *deleteImageOptions.Image,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/images/{image}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteImageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "DeleteImage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImageDeleteResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// InspectImage : Inspect an image
// Inspect a container image in the private registry.
func (containerRegistry *ContainerRegistryV1) InspectImage(inspectImageOptions *InspectImageOptions) (result *ImageInspection, response *core.DetailedResponse, err error) {
	return containerRegistry.InspectImageWithContext(context.Background(), inspectImageOptions)
}

// InspectImageWithContext is an alternate form of the InspectImage method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) InspectImageWithContext(ctx context.Context, inspectImageOptions *InspectImageOptions) (result *ImageInspection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(inspectImageOptions, "inspectImageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(inspectImageOptions, "inspectImageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"image": *inspectImageOptions.Image,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/images/{image}/json`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range inspectImageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "InspectImage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImageInspection)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetImageManifest : Get image manifest
// Get the manifest for a container image in the private registry.
func (containerRegistry *ContainerRegistryV1) GetImageManifest(getImageManifestOptions *GetImageManifestOptions) (result map[string]interface{}, response *core.DetailedResponse, err error) {
	return containerRegistry.GetImageManifestWithContext(context.Background(), getImageManifestOptions)
}

// GetImageManifestWithContext is an alternate form of the GetImageManifest method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) GetImageManifestWithContext(ctx context.Context, getImageManifestOptions *GetImageManifestOptions) (result map[string]interface{}, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getImageManifestOptions, "getImageManifestOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getImageManifestOptions, "getImageManifestOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"image": *getImageManifestOptions.Image,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/images/{image}/manifest`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getImageManifestOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "GetImageManifest")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, &result)

	return
}

// GetMessages : Get messages
// Return any published system messages.
func (containerRegistry *ContainerRegistryV1) GetMessages(getMessagesOptions *GetMessagesOptions) (result *string, response *core.DetailedResponse, err error) {
	return containerRegistry.GetMessagesWithContext(context.Background(), getMessagesOptions)
}

// GetMessagesWithContext is an alternate form of the GetMessages method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) GetMessagesWithContext(ctx context.Context, getMessagesOptions *GetMessagesOptions) (result *string, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getMessagesOptions, "getMessagesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/messages`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getMessagesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "GetMessages")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, &result)

	return
}

// ListNamespaces : List namespaces
// List authorized namespaces in the targeted IBM Cloud account.
func (containerRegistry *ContainerRegistryV1) ListNamespaces(listNamespacesOptions *ListNamespacesOptions) (result []string, response *core.DetailedResponse, err error) {
	return containerRegistry.ListNamespacesWithContext(context.Background(), listNamespacesOptions)
}

// ListNamespacesWithContext is an alternate form of the ListNamespaces method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) ListNamespacesWithContext(ctx context.Context, listNamespacesOptions *ListNamespacesOptions) (result []string, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listNamespacesOptions, "listNamespacesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/namespaces`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listNamespacesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "ListNamespaces")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, &result)

	return
}

// ListNamespaceDetails : Detailed namespace list
// Retrieves details, such as resource group, for all your namespaces in the targeted registry.
func (containerRegistry *ContainerRegistryV1) ListNamespaceDetails(listNamespaceDetailsOptions *ListNamespaceDetailsOptions) (result []NamespaceDetails, response *core.DetailedResponse, err error) {
	return containerRegistry.ListNamespaceDetailsWithContext(context.Background(), listNamespaceDetailsOptions)
}

// ListNamespaceDetailsWithContext is an alternate form of the ListNamespaceDetails method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) ListNamespaceDetailsWithContext(ctx context.Context, listNamespaceDetailsOptions *ListNamespaceDetailsOptions) (result []NamespaceDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listNamespaceDetailsOptions, "listNamespaceDetailsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/namespaces/details`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listNamespaceDetailsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "ListNamespaceDetails")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNamespaceDetails)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateNamespace : Create namespace
// Add a namespace to the targeted IBM Cloud account.
func (containerRegistry *ContainerRegistryV1) CreateNamespace(createNamespaceOptions *CreateNamespaceOptions) (result *Namespace, response *core.DetailedResponse, err error) {
	return containerRegistry.CreateNamespaceWithContext(context.Background(), createNamespaceOptions)
}

// CreateNamespaceWithContext is an alternate form of the CreateNamespace method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) CreateNamespaceWithContext(ctx context.Context, createNamespaceOptions *CreateNamespaceOptions) (result *Namespace, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createNamespaceOptions, "createNamespaceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createNamespaceOptions, "createNamespaceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"name": *createNamespaceOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/namespaces/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createNamespaceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "CreateNamespace")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}
	if createNamespaceOptions.XAuthResourceGroup != nil {
		builder.AddHeader("X-Auth-Resource-Group", fmt.Sprint(*createNamespaceOptions.XAuthResourceGroup))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNamespace)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// AssignNamespace : Assign namespace
// Assign a namespace to the specified resource group in the targeted IBM Cloud account.
func (containerRegistry *ContainerRegistryV1) AssignNamespace(assignNamespaceOptions *AssignNamespaceOptions) (result *Namespace, response *core.DetailedResponse, err error) {
	return containerRegistry.AssignNamespaceWithContext(context.Background(), assignNamespaceOptions)
}

// AssignNamespaceWithContext is an alternate form of the AssignNamespace method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) AssignNamespaceWithContext(ctx context.Context, assignNamespaceOptions *AssignNamespaceOptions) (result *Namespace, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(assignNamespaceOptions, "assignNamespaceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(assignNamespaceOptions, "assignNamespaceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"name": *assignNamespaceOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/namespaces/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range assignNamespaceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "AssignNamespace")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}
	if assignNamespaceOptions.XAuthResourceGroup != nil {
		builder.AddHeader("X-Auth-Resource-Group", fmt.Sprint(*assignNamespaceOptions.XAuthResourceGroup))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNamespace)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteNamespace : Delete namespace
// Delete the IBM Cloud Container Registry namespace from the targeted IBM Cloud account, and removes all images that
// were in that namespace.
func (containerRegistry *ContainerRegistryV1) DeleteNamespace(deleteNamespaceOptions *DeleteNamespaceOptions) (response *core.DetailedResponse, err error) {
	return containerRegistry.DeleteNamespaceWithContext(context.Background(), deleteNamespaceOptions)
}

// DeleteNamespaceWithContext is an alternate form of the DeleteNamespace method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) DeleteNamespaceWithContext(ctx context.Context, deleteNamespaceOptions *DeleteNamespaceOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteNamespaceOptions, "deleteNamespaceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteNamespaceOptions, "deleteNamespaceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"name": *deleteNamespaceOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/namespaces/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteNamespaceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "DeleteNamespace")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, nil)

	return
}

// GetPlans : Get plans
// Get plans for the targeted account.
func (containerRegistry *ContainerRegistryV1) GetPlans(getPlansOptions *GetPlansOptions) (result *Plan, response *core.DetailedResponse, err error) {
	return containerRegistry.GetPlansWithContext(context.Background(), getPlansOptions)
}

// GetPlansWithContext is an alternate form of the GetPlans method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) GetPlansWithContext(ctx context.Context, getPlansOptions *GetPlansOptions) (result *Plan, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getPlansOptions, "getPlansOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/plans`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getPlansOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "GetPlans")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPlan)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdatePlans : Update plans
// Update plans for the targeted account.
func (containerRegistry *ContainerRegistryV1) UpdatePlans(updatePlansOptions *UpdatePlansOptions) (response *core.DetailedResponse, err error) {
	return containerRegistry.UpdatePlansWithContext(context.Background(), updatePlansOptions)
}

// UpdatePlansWithContext is an alternate form of the UpdatePlans method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) UpdatePlansWithContext(ctx context.Context, updatePlansOptions *UpdatePlansOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePlansOptions, "updatePlansOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updatePlansOptions, "updatePlansOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/plans`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range updatePlansOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "UpdatePlans")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	body := make(map[string]interface{})
	if updatePlansOptions.Plan != nil {
		body["plan"] = updatePlansOptions.Plan
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, nil)

	return
}

// GetQuota : Get quotas
// Get quotas for the targeted account.
func (containerRegistry *ContainerRegistryV1) GetQuota(getQuotaOptions *GetQuotaOptions) (result *Quota, response *core.DetailedResponse, err error) {
	return containerRegistry.GetQuotaWithContext(context.Background(), getQuotaOptions)
}

// GetQuotaWithContext is an alternate form of the GetQuota method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) GetQuotaWithContext(ctx context.Context, getQuotaOptions *GetQuotaOptions) (result *Quota, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getQuotaOptions, "getQuotaOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/quotas`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getQuotaOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "GetQuota")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQuota)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateQuota : Update quotas
// Update quotas for the targeted account.
func (containerRegistry *ContainerRegistryV1) UpdateQuota(updateQuotaOptions *UpdateQuotaOptions) (response *core.DetailedResponse, err error) {
	return containerRegistry.UpdateQuotaWithContext(context.Background(), updateQuotaOptions)
}

// UpdateQuotaWithContext is an alternate form of the UpdateQuota method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) UpdateQuotaWithContext(ctx context.Context, updateQuotaOptions *UpdateQuotaOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateQuotaOptions, "updateQuotaOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateQuotaOptions, "updateQuotaOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/quotas`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateQuotaOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "UpdateQuota")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	body := make(map[string]interface{})
	if updateQuotaOptions.StorageMegabytes != nil {
		body["storage_megabytes"] = updateQuotaOptions.StorageMegabytes
	}
	if updateQuotaOptions.TrafficMegabytes != nil {
		body["traffic_megabytes"] = updateQuotaOptions.TrafficMegabytes
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, nil)

	return
}

// ListRetentionPolicies : List retention policies
// List retention policies for all namespaces in the targeted IBM Cloud account.
func (containerRegistry *ContainerRegistryV1) ListRetentionPolicies(listRetentionPoliciesOptions *ListRetentionPoliciesOptions) (result map[string]RetentionPolicy, response *core.DetailedResponse, err error) {
	return containerRegistry.ListRetentionPoliciesWithContext(context.Background(), listRetentionPoliciesOptions)
}

// ListRetentionPoliciesWithContext is an alternate form of the ListRetentionPolicies method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) ListRetentionPoliciesWithContext(ctx context.Context, listRetentionPoliciesOptions *ListRetentionPoliciesOptions) (result map[string]RetentionPolicy, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listRetentionPoliciesOptions, "listRetentionPoliciesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/retentions`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRetentionPoliciesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "ListRetentionPolicies")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRetentionPolicy)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// SetRetentionPolicy : Set retention policy
// Set the retention policy for the specified namespace.
func (containerRegistry *ContainerRegistryV1) SetRetentionPolicy(setRetentionPolicyOptions *SetRetentionPolicyOptions) (response *core.DetailedResponse, err error) {
	return containerRegistry.SetRetentionPolicyWithContext(context.Background(), setRetentionPolicyOptions)
}

// SetRetentionPolicyWithContext is an alternate form of the SetRetentionPolicy method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) SetRetentionPolicyWithContext(ctx context.Context, setRetentionPolicyOptions *SetRetentionPolicyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setRetentionPolicyOptions, "setRetentionPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setRetentionPolicyOptions, "setRetentionPolicyOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/retentions`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range setRetentionPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "SetRetentionPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	body := make(map[string]interface{})
	if setRetentionPolicyOptions.Namespace != nil {
		body["namespace"] = setRetentionPolicyOptions.Namespace
	}
	if setRetentionPolicyOptions.ImagesPerRepo != nil {
		body["images_per_repo"] = setRetentionPolicyOptions.ImagesPerRepo
	}
	if setRetentionPolicyOptions.RetainUntagged != nil {
		body["retain_untagged"] = setRetentionPolicyOptions.RetainUntagged
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, nil)

	return
}

// AnalyzeRetentionPolicy : Analyze retention policy
// Analyze a retention policy, and get a list of what would be deleted by it.
func (containerRegistry *ContainerRegistryV1) AnalyzeRetentionPolicy(analyzeRetentionPolicyOptions *AnalyzeRetentionPolicyOptions) (result map[string][]string, response *core.DetailedResponse, err error) {
	return containerRegistry.AnalyzeRetentionPolicyWithContext(context.Background(), analyzeRetentionPolicyOptions)
}

// AnalyzeRetentionPolicyWithContext is an alternate form of the AnalyzeRetentionPolicy method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) AnalyzeRetentionPolicyWithContext(ctx context.Context, analyzeRetentionPolicyOptions *AnalyzeRetentionPolicyOptions) (result map[string][]string, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(analyzeRetentionPolicyOptions, "analyzeRetentionPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(analyzeRetentionPolicyOptions, "analyzeRetentionPolicyOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/retentions/analyze`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range analyzeRetentionPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "AnalyzeRetentionPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	body := make(map[string]interface{})
	if analyzeRetentionPolicyOptions.Namespace != nil {
		body["namespace"] = analyzeRetentionPolicyOptions.Namespace
	}
	if analyzeRetentionPolicyOptions.ImagesPerRepo != nil {
		body["images_per_repo"] = analyzeRetentionPolicyOptions.ImagesPerRepo
	}
	if analyzeRetentionPolicyOptions.RetainUntagged != nil {
		body["retain_untagged"] = analyzeRetentionPolicyOptions.RetainUntagged
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, &result)

	return
}

// GetRetentionPolicy : Get retention policy
// Get the retention policy for the specified namespace.
func (containerRegistry *ContainerRegistryV1) GetRetentionPolicy(getRetentionPolicyOptions *GetRetentionPolicyOptions) (result *RetentionPolicy, response *core.DetailedResponse, err error) {
	return containerRegistry.GetRetentionPolicyWithContext(context.Background(), getRetentionPolicyOptions)
}

// GetRetentionPolicyWithContext is an alternate form of the GetRetentionPolicy method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) GetRetentionPolicyWithContext(ctx context.Context, getRetentionPolicyOptions *GetRetentionPolicyOptions) (result *RetentionPolicy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRetentionPolicyOptions, "getRetentionPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRetentionPolicyOptions, "getRetentionPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"namespace": *getRetentionPolicyOptions.Namespace,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/retentions/{namespace}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRetentionPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "GetRetentionPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRetentionPolicy)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetSettings : Get account settings
// Get account settings for the targeted account.
func (containerRegistry *ContainerRegistryV1) GetSettings(getSettingsOptions *GetSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	return containerRegistry.GetSettingsWithContext(context.Background(), getSettingsOptions)
}

// GetSettingsWithContext is an alternate form of the GetSettings method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) GetSettingsWithContext(ctx context.Context, getSettingsOptions *GetSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSettingsOptions, "getSettingsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/settings`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "GetSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettings)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateSettings : Update account settings
// Update settings for the targeted account.
func (containerRegistry *ContainerRegistryV1) UpdateSettings(updateSettingsOptions *UpdateSettingsOptions) (response *core.DetailedResponse, err error) {
	return containerRegistry.UpdateSettingsWithContext(context.Background(), updateSettingsOptions)
}

// UpdateSettingsWithContext is an alternate form of the UpdateSettings method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) UpdateSettingsWithContext(ctx context.Context, updateSettingsOptions *UpdateSettingsOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSettingsOptions, "updateSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSettingsOptions, "updateSettingsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/settings`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "UpdateSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	body := make(map[string]interface{})
	if updateSettingsOptions.PlatformMetrics != nil {
		body["platform_metrics"] = updateSettingsOptions.PlatformMetrics
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, nil)

	return
}

// DeleteImageTag : Delete tag
// Untag a container image in the registry.
func (containerRegistry *ContainerRegistryV1) DeleteImageTag(deleteImageTagOptions *DeleteImageTagOptions) (result *ImageDeleteResult, response *core.DetailedResponse, err error) {
	return containerRegistry.DeleteImageTagWithContext(context.Background(), deleteImageTagOptions)
}

// DeleteImageTagWithContext is an alternate form of the DeleteImageTag method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) DeleteImageTagWithContext(ctx context.Context, deleteImageTagOptions *DeleteImageTagOptions) (result *ImageDeleteResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteImageTagOptions, "deleteImageTagOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteImageTagOptions, "deleteImageTagOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"image": *deleteImageTagOptions.Image,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/tags/{image}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteImageTagOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "DeleteImageTag")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImageDeleteResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListDeletedImages : List deleted images
// List all images that are in the trash can.
func (containerRegistry *ContainerRegistryV1) ListDeletedImages(listDeletedImagesOptions *ListDeletedImagesOptions) (result map[string]Trash, response *core.DetailedResponse, err error) {
	return containerRegistry.ListDeletedImagesWithContext(context.Background(), listDeletedImagesOptions)
}

// ListDeletedImagesWithContext is an alternate form of the ListDeletedImages method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) ListDeletedImagesWithContext(ctx context.Context, listDeletedImagesOptions *ListDeletedImagesOptions) (result map[string]Trash, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listDeletedImagesOptions, "listDeletedImagesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/trash`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listDeletedImagesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "ListDeletedImages")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	if listDeletedImagesOptions.Namespace != nil {
		builder.AddQuery("namespace", fmt.Sprint(*listDeletedImagesOptions.Namespace))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrash)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// RestoreTags : Restore a digest and all associated tags
// In the targeted region, restore a digest, and all of its tags in the same repository, from the trash.
func (containerRegistry *ContainerRegistryV1) RestoreTags(restoreTagsOptions *RestoreTagsOptions) (result *RestoreResult, response *core.DetailedResponse, err error) {
	return containerRegistry.RestoreTagsWithContext(context.Background(), restoreTagsOptions)
}

// RestoreTagsWithContext is an alternate form of the RestoreTags method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) RestoreTagsWithContext(ctx context.Context, restoreTagsOptions *RestoreTagsOptions) (result *RestoreResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(restoreTagsOptions, "restoreTagsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(restoreTagsOptions, "restoreTagsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"digest": *restoreTagsOptions.Digest,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/trash/{digest}/restoretags`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range restoreTagsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "RestoreTags")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = containerRegistry.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRestoreResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// RestoreImage : Restore deleted image
// Restore an image from the trash can.
func (containerRegistry *ContainerRegistryV1) RestoreImage(restoreImageOptions *RestoreImageOptions) (response *core.DetailedResponse, err error) {
	return containerRegistry.RestoreImageWithContext(context.Background(), restoreImageOptions)
}

// RestoreImageWithContext is an alternate form of the RestoreImage method which supports a Context parameter
func (containerRegistry *ContainerRegistryV1) RestoreImageWithContext(ctx context.Context, restoreImageOptions *RestoreImageOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(restoreImageOptions, "restoreImageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(restoreImageOptions, "restoreImageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"image": *restoreImageOptions.Image,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = containerRegistry.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(containerRegistry.Service.Options.URL, `/api/v1/trash/{image}/restore`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range restoreImageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("container_registry", "V1", "RestoreImage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if containerRegistry.Account != nil {
		builder.AddHeader("Account", fmt.Sprint(*containerRegistry.Account))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = containerRegistry.Service.Request(request, nil)

	return
}

// AccountSettings : Account settings for the targeted IBM Cloud account.
type AccountSettings struct {
	// Opt in to IBM Cloud Container Registry publishing platform metrics.
	PlatformMetrics *bool `json:"platform_metrics,omitempty"`
}

// UnmarshalAccountSettings unmarshals an instance of AccountSettings from the specified map of raw messages.
func UnmarshalAccountSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettings)
	err = core.UnmarshalPrimitive(m, "platform_metrics", &obj.PlatformMetrics)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AnalyzeRetentionPolicyOptions : The AnalyzeRetentionPolicy options.
type AnalyzeRetentionPolicyOptions struct {
	// The namespace to which the retention policy is attached.
	Namespace *string `validate:"required"`

	// Determines how many images will be retained for each repository when the retention policy is executed. The value -1
	// denotes 'Unlimited' (all images are retained).
	ImagesPerRepo *int64

	// Determines if untagged images are retained when executing the retention policy. This is false by default meaning
	// untagged images will be deleted when the policy is executed.
	RetainUntagged *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewAnalyzeRetentionPolicyOptions : Instantiate AnalyzeRetentionPolicyOptions
func (*ContainerRegistryV1) NewAnalyzeRetentionPolicyOptions(namespace string) *AnalyzeRetentionPolicyOptions {
	return &AnalyzeRetentionPolicyOptions{
		Namespace: core.StringPtr(namespace),
	}
}

// SetNamespace : Allow user to set Namespace
func (options *AnalyzeRetentionPolicyOptions) SetNamespace(namespace string) *AnalyzeRetentionPolicyOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetImagesPerRepo : Allow user to set ImagesPerRepo
func (options *AnalyzeRetentionPolicyOptions) SetImagesPerRepo(imagesPerRepo int64) *AnalyzeRetentionPolicyOptions {
	options.ImagesPerRepo = core.Int64Ptr(imagesPerRepo)
	return options
}

// SetRetainUntagged : Allow user to set RetainUntagged
func (options *AnalyzeRetentionPolicyOptions) SetRetainUntagged(retainUntagged bool) *AnalyzeRetentionPolicyOptions {
	options.RetainUntagged = core.BoolPtr(retainUntagged)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *AnalyzeRetentionPolicyOptions) SetHeaders(param map[string]string) *AnalyzeRetentionPolicyOptions {
	options.Headers = param
	return options
}

// AssignNamespaceOptions : The AssignNamespace options.
type AssignNamespaceOptions struct {
	// The ID of the resource group that the namespace will be created within.
	XAuthResourceGroup *string `validate:"required"`

	// The name of the namespace to be updated.
	Name *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewAssignNamespaceOptions : Instantiate AssignNamespaceOptions
func (*ContainerRegistryV1) NewAssignNamespaceOptions(xAuthResourceGroup string, name string) *AssignNamespaceOptions {
	return &AssignNamespaceOptions{
		XAuthResourceGroup: core.StringPtr(xAuthResourceGroup),
		Name:               core.StringPtr(name),
	}
}

// SetXAuthResourceGroup : Allow user to set XAuthResourceGroup
func (options *AssignNamespaceOptions) SetXAuthResourceGroup(xAuthResourceGroup string) *AssignNamespaceOptions {
	options.XAuthResourceGroup = core.StringPtr(xAuthResourceGroup)
	return options
}

// SetName : Allow user to set Name
func (options *AssignNamespaceOptions) SetName(name string) *AssignNamespaceOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *AssignNamespaceOptions) SetHeaders(param map[string]string) *AssignNamespaceOptions {
	options.Headers = param
	return options
}

// AuthOptions : The authorization options for the targeted IBM Cloud account.
type AuthOptions struct {
	// Enable role based authorization when authenticating with IBM Cloud IAM.
	IamAuthz *bool `json:"iam_authz,omitempty"`

	// Restrict account to only be able to push and pull images over private connections.
	PrivateOnly *bool `json:"private_only,omitempty"`
}

// UnmarshalAuthOptions unmarshals an instance of AuthOptions from the specified map of raw messages.
func UnmarshalAuthOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AuthOptions)
	err = core.UnmarshalPrimitive(m, "iam_authz", &obj.IamAuthz)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_only", &obj.PrivateOnly)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BulkDeleteImagesOptions : The BulkDeleteImages options.
type BulkDeleteImagesOptions struct {
	// The full IBM Cloud registry path to the images that you want to delete, including its digest. All tags for the
	// supplied digest are removed.
	BulkDelete []string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewBulkDeleteImagesOptions : Instantiate BulkDeleteImagesOptions
func (*ContainerRegistryV1) NewBulkDeleteImagesOptions(bulkDelete []string) *BulkDeleteImagesOptions {
	return &BulkDeleteImagesOptions{
		BulkDelete: bulkDelete,
	}
}

// SetBulkDelete : Allow user to set BulkDelete
func (options *BulkDeleteImagesOptions) SetBulkDelete(bulkDelete []string) *BulkDeleteImagesOptions {
	options.BulkDelete = bulkDelete
	return options
}

// SetHeaders : Allow user to set Headers
func (options *BulkDeleteImagesOptions) SetHeaders(param map[string]string) *BulkDeleteImagesOptions {
	options.Headers = param
	return options
}

// Config : The configuration data about a container.
type Config struct {
	// True if command is already escaped (Windows specific).
	ArgsEscaped *bool `json:"ArgsEscaped,omitempty"`

	// If true, standard error is attached.
	AttachStderr *bool `json:"AttachStderr,omitempty"`

	// If true, standard input is attached, which makes possible user interaction.
	AttachStdin *bool `json:"AttachStdin,omitempty"`

	// If true, standard output is attached.
	AttachStdout *bool `json:"AttachStdout,omitempty"`

	// Command that is run when starting the container.
	Cmd []string `json:"Cmd,omitempty"`

	// The FQDN for the container.
	Domainname *string `json:"Domainname,omitempty"`

	// Entrypoint to run when starting the container.
	Entrypoint []string `json:"Entrypoint,omitempty"`

	// List of environment variables to set in the container.
	Env []string `json:"Env,omitempty"`

	// A list of exposed ports in a format [123:{},456:{}].
	ExposedPorts map[string]interface{} `json:"ExposedPorts,omitempty"`

	Healthcheck *HealthConfig `json:"Healthcheck,omitempty"`

	// The host name of the container.
	Hostname *string `json:"Hostname,omitempty"`

	// Name of the image as it was passed by the operator (eg. could be symbolic).
	Image *string `json:"Image,omitempty"`

	// List of labels set to this container.
	Labels map[string]string `json:"Labels,omitempty"`

	// The MAC Address of the container.
	MacAddress *string `json:"MacAddress,omitempty"`

	// If true, containers are not given network access.
	NetworkDisabled *bool `json:"NetworkDisabled,omitempty"`

	// ONBUILD metadata that were defined on the image Dockerfile
	// https://docs.docker.com/engine/reference/builder/#onbuild.
	OnBuild []string `json:"OnBuild,omitempty"`

	// Open stdin.
	OpenStdin *bool `json:"OpenStdin,omitempty"`

	// Shell for shell-form of RUN, CMD, ENTRYPOINT.
	Shell []string `json:"Shell,omitempty"`

	// If true, close stdin after the 1 attached client disconnects.
	StdinOnce *bool `json:"StdinOnce,omitempty"`

	// Signal to stop a container.
	StopSignal *string `json:"StopSignal,omitempty"`

	// Timeout (in seconds) to stop a container.
	StopTimeout *int64 `json:"StopTimeout,omitempty"`

	// Attach standard streams to a tty, including stdin if it is not closed.
	Tty *bool `json:"Tty,omitempty"`

	// The user that will run the command(s) inside the container.
	User *string `json:"User,omitempty"`

	// List of volumes (mounts) used for the container.
	Volumes map[string]interface{} `json:"Volumes,omitempty"`

	// Current working directory (PWD) in the command will be launched.
	WorkingDir *string `json:"WorkingDir,omitempty"`
}

// UnmarshalConfig unmarshals an instance of Config from the specified map of raw messages.
func UnmarshalConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Config)
	err = core.UnmarshalPrimitive(m, "ArgsEscaped", &obj.ArgsEscaped)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "AttachStderr", &obj.AttachStderr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "AttachStdin", &obj.AttachStdin)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "AttachStdout", &obj.AttachStdout)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Cmd", &obj.Cmd)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Domainname", &obj.Domainname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Entrypoint", &obj.Entrypoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Env", &obj.Env)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ExposedPorts", &obj.ExposedPorts)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "Healthcheck", &obj.Healthcheck, UnmarshalHealthConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Hostname", &obj.Hostname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Image", &obj.Image)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "MacAddress", &obj.MacAddress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "NetworkDisabled", &obj.NetworkDisabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "OnBuild", &obj.OnBuild)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "OpenStdin", &obj.OpenStdin)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Shell", &obj.Shell)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "StdinOnce", &obj.StdinOnce)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "StopSignal", &obj.StopSignal)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "StopTimeout", &obj.StopTimeout)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Tty", &obj.Tty)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "User", &obj.User)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Volumes", &obj.Volumes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "WorkingDir", &obj.WorkingDir)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateNamespaceOptions : The CreateNamespace options.
type CreateNamespaceOptions struct {
	// The name of the namespace.
	Name *string `validate:"required,ne="`

	// The ID of the resource group that the namespace will be created within.
	XAuthResourceGroup *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateNamespaceOptions : Instantiate CreateNamespaceOptions
func (*ContainerRegistryV1) NewCreateNamespaceOptions(name string) *CreateNamespaceOptions {
	return &CreateNamespaceOptions{
		Name: core.StringPtr(name),
	}
}

// SetName : Allow user to set Name
func (options *CreateNamespaceOptions) SetName(name string) *CreateNamespaceOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetXAuthResourceGroup : Allow user to set XAuthResourceGroup
func (options *CreateNamespaceOptions) SetXAuthResourceGroup(xAuthResourceGroup string) *CreateNamespaceOptions {
	options.XAuthResourceGroup = core.StringPtr(xAuthResourceGroup)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateNamespaceOptions) SetHeaders(param map[string]string) *CreateNamespaceOptions {
	options.Headers = param
	return options
}

// DeleteImageOptions : The DeleteImage options.
type DeleteImageOptions struct {
	// The full IBM Cloud registry path to the image that you want to delete, including its tag. If you do not provide a
	// specific tag, the version with the `latest` tag is removed.
	Image *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteImageOptions : Instantiate DeleteImageOptions
func (*ContainerRegistryV1) NewDeleteImageOptions(image string) *DeleteImageOptions {
	return &DeleteImageOptions{
		Image: core.StringPtr(image),
	}
}

// SetImage : Allow user to set Image
func (options *DeleteImageOptions) SetImage(image string) *DeleteImageOptions {
	options.Image = core.StringPtr(image)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteImageOptions) SetHeaders(param map[string]string) *DeleteImageOptions {
	options.Headers = param
	return options
}

// DeleteImageTagOptions : The DeleteImageTag options.
type DeleteImageTagOptions struct {
	// The name of the image that you want to delete, in the format &lt;REPOSITORY&gt;:&lt;TAG&gt;.
	Image *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteImageTagOptions : Instantiate DeleteImageTagOptions
func (*ContainerRegistryV1) NewDeleteImageTagOptions(image string) *DeleteImageTagOptions {
	return &DeleteImageTagOptions{
		Image: core.StringPtr(image),
	}
}

// SetImage : Allow user to set Image
func (options *DeleteImageTagOptions) SetImage(image string) *DeleteImageTagOptions {
	options.Image = core.StringPtr(image)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteImageTagOptions) SetHeaders(param map[string]string) *DeleteImageTagOptions {
	options.Headers = param
	return options
}

// DeleteNamespaceOptions : The DeleteNamespace options.
type DeleteNamespaceOptions struct {
	// The name of the namespace that you want to delete.
	Name *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteNamespaceOptions : Instantiate DeleteNamespaceOptions
func (*ContainerRegistryV1) NewDeleteNamespaceOptions(name string) *DeleteNamespaceOptions {
	return &DeleteNamespaceOptions{
		Name: core.StringPtr(name),
	}
}

// SetName : Allow user to set Name
func (options *DeleteNamespaceOptions) SetName(name string) *DeleteNamespaceOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteNamespaceOptions) SetHeaders(param map[string]string) *DeleteNamespaceOptions {
	options.Headers = param
	return options
}

// GetAuthOptions : The GetAuth options.
type GetAuthOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAuthOptions : Instantiate GetAuthOptions
func (*ContainerRegistryV1) NewGetAuthOptions() *GetAuthOptions {
	return &GetAuthOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetAuthOptions) SetHeaders(param map[string]string) *GetAuthOptions {
	options.Headers = param
	return options
}

// GetImageManifestOptions : The GetImageManifest options.
type GetImageManifestOptions struct {
	// The full IBM Cloud registry path to the image that you want to inspect. Run `ibmcloud cr images` or call the `GET
	// /images/json` endpoint to review images that are in the registry.
	Image *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetImageManifestOptions : Instantiate GetImageManifestOptions
func (*ContainerRegistryV1) NewGetImageManifestOptions(image string) *GetImageManifestOptions {
	return &GetImageManifestOptions{
		Image: core.StringPtr(image),
	}
}

// SetImage : Allow user to set Image
func (options *GetImageManifestOptions) SetImage(image string) *GetImageManifestOptions {
	options.Image = core.StringPtr(image)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetImageManifestOptions) SetHeaders(param map[string]string) *GetImageManifestOptions {
	options.Headers = param
	return options
}

// GetMessagesOptions : The GetMessages options.
type GetMessagesOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetMessagesOptions : Instantiate GetMessagesOptions
func (*ContainerRegistryV1) NewGetMessagesOptions() *GetMessagesOptions {
	return &GetMessagesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetMessagesOptions) SetHeaders(param map[string]string) *GetMessagesOptions {
	options.Headers = param
	return options
}

// GetPlansOptions : The GetPlans options.
type GetPlansOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetPlansOptions : Instantiate GetPlansOptions
func (*ContainerRegistryV1) NewGetPlansOptions() *GetPlansOptions {
	return &GetPlansOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetPlansOptions) SetHeaders(param map[string]string) *GetPlansOptions {
	options.Headers = param
	return options
}

// GetQuotaOptions : The GetQuota options.
type GetQuotaOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetQuotaOptions : Instantiate GetQuotaOptions
func (*ContainerRegistryV1) NewGetQuotaOptions() *GetQuotaOptions {
	return &GetQuotaOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetQuotaOptions) SetHeaders(param map[string]string) *GetQuotaOptions {
	options.Headers = param
	return options
}

// GetRetentionPolicyOptions : The GetRetentionPolicy options.
type GetRetentionPolicyOptions struct {
	// Gets the retention policy for the specified namespace.
	Namespace *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRetentionPolicyOptions : Instantiate GetRetentionPolicyOptions
func (*ContainerRegistryV1) NewGetRetentionPolicyOptions(namespace string) *GetRetentionPolicyOptions {
	return &GetRetentionPolicyOptions{
		Namespace: core.StringPtr(namespace),
	}
}

// SetNamespace : Allow user to set Namespace
func (options *GetRetentionPolicyOptions) SetNamespace(namespace string) *GetRetentionPolicyOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetRetentionPolicyOptions) SetHeaders(param map[string]string) *GetRetentionPolicyOptions {
	options.Headers = param
	return options
}

// GetSettingsOptions : The GetSettings options.
type GetSettingsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSettingsOptions : Instantiate GetSettingsOptions
func (*ContainerRegistryV1) NewGetSettingsOptions() *GetSettingsOptions {
	return &GetSettingsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetSettingsOptions) SetHeaders(param map[string]string) *GetSettingsOptions {
	options.Headers = param
	return options
}

// HealthConfig : HealthConfig struct
type HealthConfig struct {
	// A Duration represents the elapsed time between two instants as an int64 nanosecond count.
	Interval *int64 `json:"Interval,omitempty"`

	// The number of consecutive failures needed to consider a container as unhealthy. Zero means inherit.
	Retries *int64 `json:"Retries,omitempty"`

	// The test to perform to check that the container is healthy. An empty slice means to inherit the default. The options
	// are:
	// {} : inherit healthcheck
	// {"NONE"} : disable healthcheck
	// {"CMD", args...} : exec arguments directly
	// {"CMD-SHELL", command} : run command with system's default shell.
	Test []string `json:"Test,omitempty"`

	// A Duration represents the elapsed time between two instants as an int64 nanosecond count.
	Timeout *int64 `json:"Timeout,omitempty"`
}

// UnmarshalHealthConfig unmarshals an instance of HealthConfig from the specified map of raw messages.
func UnmarshalHealthConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HealthConfig)
	err = core.UnmarshalPrimitive(m, "Interval", &obj.Interval)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Retries", &obj.Retries)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Test", &obj.Test)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Timeout", &obj.Timeout)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImageBulkDeleteError : Information about a failure to delete an image as part of a bulk delete.
type ImageBulkDeleteError struct {
	// An API error code.
	Code *string `json:"code,omitempty"`

	// The English text message associated with the error code.
	Message *string `json:"message,omitempty"`
}

// UnmarshalImageBulkDeleteError unmarshals an instance of ImageBulkDeleteError from the specified map of raw messages.
func UnmarshalImageBulkDeleteError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImageBulkDeleteError)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
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

// ImageBulkDeleteResult : The results of a bulk image delete request.
type ImageBulkDeleteResult struct {
	// A map of digests to the error object that explains the failure.
	Error map[string]ImageBulkDeleteError `json:"error,omitempty"`

	// A list of digests which were deleted successfully.
	Success []string `json:"success,omitempty"`
}

// UnmarshalImageBulkDeleteResult unmarshals an instance of ImageBulkDeleteResult from the specified map of raw messages.
func UnmarshalImageBulkDeleteResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImageBulkDeleteResult)
	err = core.UnmarshalModel(m, "error", &obj.Error, UnmarshalImageBulkDeleteError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImageDeleteResult : ImageDeleteResult struct
type ImageDeleteResult struct {
	Untagged *string `json:"Untagged,omitempty"`
}

// UnmarshalImageDeleteResult unmarshals an instance of ImageDeleteResult from the specified map of raw messages.
func UnmarshalImageDeleteResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImageDeleteResult)
	err = core.UnmarshalPrimitive(m, "Untagged", &obj.Untagged)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImageDigest : Important information about an image.
type ImageDigest struct {
	// The build date of the image.
	Created *int64 `json:"created,omitempty"`

	// The image digest.
	ID *string `json:"id,omitempty"`

	// The type of the image, such as 'Docker Image Manifest V2, Schema 2' or 'OCI Image Manifest v1'.
	ManifestType *string `json:"manifestType,omitempty"`

	// A map of image repositories to tags.
	RepoTags map[string]interface{} `json:"repoTags,omitempty"`

	// The size of the image in bytes.
	Size *int64 `json:"size,omitempty"`
}

// UnmarshalImageDigest unmarshals an instance of ImageDigest from the specified map of raw messages.
func UnmarshalImageDigest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImageDigest)
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "manifestType", &obj.ManifestType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "repoTags", &obj.RepoTags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImageInspection : An image JSON output consistent with the Docker Remote API.
type ImageInspection struct {
	// The processor architecture used to build this image, and required to run it.
	Architecture *string `json:"Architecture,omitempty"`

	// The author of the image.
	Author *string `json:"Author,omitempty"`

	// A plain text description of the image.
	Comment *string `json:"Comment,omitempty"`

	// The configuration data about a container.
	Config *Config `json:"Config,omitempty"`

	// The ID of the container which created this image.
	Container *string `json:"Container,omitempty"`

	// The configuration data about a container.
	ContainerConfig *Config `json:"ContainerConfig,omitempty"`

	// The unix timestamp for the date when the image was created.
	Created *string `json:"Created,omitempty"`

	// The Docker version used to build this image.
	DockerVersion *string `json:"DockerVersion,omitempty"`

	// The image ID.
	ID *string `json:"Id,omitempty"`

	// Media type of the manifest for the image.
	ManifestType *string `json:"ManifestType,omitempty"`

	// The operating system family used to build this image, and required to run it.
	Os *string `json:"Os,omitempty"`

	// The version of the operating system used to build this image.
	OsVersion *string `json:"OsVersion,omitempty"`

	// The ID of the base image for this image.
	Parent *string `json:"Parent,omitempty"`

	// RootFS contains information about the root filesystem of a container image.
	RootFs *RootFs `json:"RootFS,omitempty"`

	// The size of the image in bytes.
	Size *int64 `json:"Size,omitempty"`

	// The sum of the size of each layer in the image in bytes.
	VirtualSize *int64 `json:"VirtualSize,omitempty"`
}

// UnmarshalImageInspection unmarshals an instance of ImageInspection from the specified map of raw messages.
func UnmarshalImageInspection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImageInspection)
	err = core.UnmarshalPrimitive(m, "Architecture", &obj.Architecture)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Author", &obj.Author)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Comment", &obj.Comment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "Config", &obj.Config, UnmarshalConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Container", &obj.Container)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "ContainerConfig", &obj.ContainerConfig, UnmarshalConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "DockerVersion", &obj.DockerVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ManifestType", &obj.ManifestType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Os", &obj.Os)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "OsVersion", &obj.OsVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Parent", &obj.Parent)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "RootFS", &obj.RootFs, UnmarshalRootFs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Size", &obj.Size)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "VirtualSize", &obj.VirtualSize)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InspectImageOptions : The InspectImage options.
type InspectImageOptions struct {
	// The full IBM Cloud registry path to the image that you want to inspect. Run `ibmcloud cr images` or call the `GET
	// /images/json` endpoint to review images that are in the registry.
	Image *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewInspectImageOptions : Instantiate InspectImageOptions
func (*ContainerRegistryV1) NewInspectImageOptions(image string) *InspectImageOptions {
	return &InspectImageOptions{
		Image: core.StringPtr(image),
	}
}

// SetImage : Allow user to set Image
func (options *InspectImageOptions) SetImage(image string) *InspectImageOptions {
	options.Image = core.StringPtr(image)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *InspectImageOptions) SetHeaders(param map[string]string) *InspectImageOptions {
	options.Headers = param
	return options
}

// ListDeletedImagesOptions : The ListDeletedImages options.
type ListDeletedImagesOptions struct {
	// Limit results to trash can images in the given namespace only.
	Namespace *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListDeletedImagesOptions : Instantiate ListDeletedImagesOptions
func (*ContainerRegistryV1) NewListDeletedImagesOptions() *ListDeletedImagesOptions {
	return &ListDeletedImagesOptions{}
}

// SetNamespace : Allow user to set Namespace
func (options *ListDeletedImagesOptions) SetNamespace(namespace string) *ListDeletedImagesOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListDeletedImagesOptions) SetHeaders(param map[string]string) *ListDeletedImagesOptions {
	options.Headers = param
	return options
}

// ListImageDigestsOptions : The ListImageDigests options.
type ListImageDigestsOptions struct {
	// ExcludeTagged returns only untagged digests.
	ExcludeTagged *bool

	// ExcludeVA returns the digest list with no VA scan results.
	ExcludeVa *bool

	// When true, API will return the IBM public images if they exist in the targeted region.
	IncludeIBM *bool

	// Repositories in which to restrict the output. If left empty all images for the account will be returned.
	Repositories []string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListImageDigestsOptions : Instantiate ListImageDigestsOptions
func (*ContainerRegistryV1) NewListImageDigestsOptions() *ListImageDigestsOptions {
	return &ListImageDigestsOptions{}
}

// SetExcludeTagged : Allow user to set ExcludeTagged
func (options *ListImageDigestsOptions) SetExcludeTagged(excludeTagged bool) *ListImageDigestsOptions {
	options.ExcludeTagged = core.BoolPtr(excludeTagged)
	return options
}

// SetExcludeVa : Allow user to set ExcludeVa
func (options *ListImageDigestsOptions) SetExcludeVa(excludeVa bool) *ListImageDigestsOptions {
	options.ExcludeVa = core.BoolPtr(excludeVa)
	return options
}

// SetIncludeIBM : Allow user to set IncludeIBM
func (options *ListImageDigestsOptions) SetIncludeIBM(includeIBM bool) *ListImageDigestsOptions {
	options.IncludeIBM = core.BoolPtr(includeIBM)
	return options
}

// SetRepositories : Allow user to set Repositories
func (options *ListImageDigestsOptions) SetRepositories(repositories []string) *ListImageDigestsOptions {
	options.Repositories = repositories
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListImageDigestsOptions) SetHeaders(param map[string]string) *ListImageDigestsOptions {
	options.Headers = param
	return options
}

// ListImagesOptions : The ListImages options.
type ListImagesOptions struct {
	// Lists images that are stored in the specified namespace only. Query multiple namespaces by specifying this option
	// for each namespace. If this option is not specified, images from all namespaces in the specified IBM Cloud account
	// are listed.
	Namespace *string

	// Includes IBM-provided public images in the list of images. If this option is not specified, private images are
	// listed only. If this option is specified more than once, the last parsed setting is the setting that is used.
	IncludeIBM *bool

	// Includes private images in the list of images. If this option is not specified, private images are listed. If this
	// option is specified more than once, the last parsed setting is the setting that is used.
	IncludePrivate *bool

	// Includes tags that reference multi-architecture manifest lists in the image list. If this option is not specified,
	// tagged manifest lists are not shown in the list. If this option is specified more than once, the last parsed setting
	// is the setting that is used.
	IncludeManifestLists *bool

	// Displays Vulnerability Advisor status for the listed images. If this option is specified more than once, the last
	// parsed setting is the setting that is used.
	Vulnerabilities *bool

	// Lists images that are stored in the specified repository, under your namespaces. Query multiple repositories by
	// specifying this option for each repository. If this option is not specified, images from all repos are listed.
	Repository *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListImagesOptions : Instantiate ListImagesOptions
func (*ContainerRegistryV1) NewListImagesOptions() *ListImagesOptions {
	return &ListImagesOptions{}
}

// SetNamespace : Allow user to set Namespace
func (options *ListImagesOptions) SetNamespace(namespace string) *ListImagesOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetIncludeIBM : Allow user to set IncludeIBM
func (options *ListImagesOptions) SetIncludeIBM(includeIBM bool) *ListImagesOptions {
	options.IncludeIBM = core.BoolPtr(includeIBM)
	return options
}

// SetIncludePrivate : Allow user to set IncludePrivate
func (options *ListImagesOptions) SetIncludePrivate(includePrivate bool) *ListImagesOptions {
	options.IncludePrivate = core.BoolPtr(includePrivate)
	return options
}

// SetIncludeManifestLists : Allow user to set IncludeManifestLists
func (options *ListImagesOptions) SetIncludeManifestLists(includeManifestLists bool) *ListImagesOptions {
	options.IncludeManifestLists = core.BoolPtr(includeManifestLists)
	return options
}

// SetVulnerabilities : Allow user to set Vulnerabilities
func (options *ListImagesOptions) SetVulnerabilities(vulnerabilities bool) *ListImagesOptions {
	options.Vulnerabilities = core.BoolPtr(vulnerabilities)
	return options
}

// SetRepository : Allow user to set Repository
func (options *ListImagesOptions) SetRepository(repository string) *ListImagesOptions {
	options.Repository = core.StringPtr(repository)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListImagesOptions) SetHeaders(param map[string]string) *ListImagesOptions {
	options.Headers = param
	return options
}

// ListNamespaceDetailsOptions : The ListNamespaceDetails options.
type ListNamespaceDetailsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListNamespaceDetailsOptions : Instantiate ListNamespaceDetailsOptions
func (*ContainerRegistryV1) NewListNamespaceDetailsOptions() *ListNamespaceDetailsOptions {
	return &ListNamespaceDetailsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListNamespaceDetailsOptions) SetHeaders(param map[string]string) *ListNamespaceDetailsOptions {
	options.Headers = param
	return options
}

// ListNamespacesOptions : The ListNamespaces options.
type ListNamespacesOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListNamespacesOptions : Instantiate ListNamespacesOptions
func (*ContainerRegistryV1) NewListNamespacesOptions() *ListNamespacesOptions {
	return &ListNamespacesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListNamespacesOptions) SetHeaders(param map[string]string) *ListNamespacesOptions {
	options.Headers = param
	return options
}

// ListRetentionPoliciesOptions : The ListRetentionPolicies options.
type ListRetentionPoliciesOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListRetentionPoliciesOptions : Instantiate ListRetentionPoliciesOptions
func (*ContainerRegistryV1) NewListRetentionPoliciesOptions() *ListRetentionPoliciesOptions {
	return &ListRetentionPoliciesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListRetentionPoliciesOptions) SetHeaders(param map[string]string) *ListRetentionPoliciesOptions {
	options.Headers = param
	return options
}

// Namespace : Namespace struct
type Namespace struct {
	Namespace *string `json:"namespace,omitempty"`
}

// UnmarshalNamespace unmarshals an instance of Namespace from the specified map of raw messages.
func UnmarshalNamespace(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Namespace)
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NamespaceDetails : Details of a namespace.
type NamespaceDetails struct {
	// The IBM Cloud account that owns the namespace.
	Account *string `json:"account,omitempty"`

	// When the namespace was created.
	CreatedDate *string `json:"created_date,omitempty"`

	// If the namespace has been assigned to a resource group, this is the IBM Cloud CRN representing the namespace.
	CRN *string `json:"crn,omitempty"`

	Name *string `json:"name,omitempty"`

	// When the namespace was assigned to a resource group.
	ResourceCreatedDate *string `json:"resource_created_date,omitempty"`

	// The resource group that the namespace is assigned to.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// When the namespace was last updated.
	UpdatedDate *string `json:"updated_date,omitempty"`
}

// UnmarshalNamespaceDetails unmarshals an instance of NamespaceDetails from the specified map of raw messages.
func UnmarshalNamespaceDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NamespaceDetails)
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_date", &obj.CreatedDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_created_date", &obj.ResourceCreatedDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_date", &obj.UpdatedDate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Plan : The plan for the targeted IBM Cloud account.
type Plan struct {
	Plan *string `json:"plan,omitempty"`
}

// UnmarshalPlan unmarshals an instance of Plan from the specified map of raw messages.
func UnmarshalPlan(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Plan)
	err = core.UnmarshalPrimitive(m, "plan", &obj.Plan)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Quota : Current usage and limits for the targeted IBM Cloud account.
type Quota struct {
	Limit *QuotaDetails `json:"limit,omitempty"`

	Usage *QuotaDetails `json:"usage,omitempty"`
}

// UnmarshalQuota unmarshals an instance of Quota from the specified map of raw messages.
func UnmarshalQuota(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Quota)
	err = core.UnmarshalModel(m, "limit", &obj.Limit, UnmarshalQuotaDetails)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "usage", &obj.Usage, UnmarshalQuotaDetails)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QuotaDetails : QuotaDetails struct
type QuotaDetails struct {
	// Storage quota or usage in bytes. The value -1 denotes "Unlimited".
	StorageBytes *int64 `json:"storage_bytes,omitempty"`

	// Traffic quota or usage in bytes. The value -1 denotes "Unlimited".
	TrafficBytes *int64 `json:"traffic_bytes,omitempty"`
}

// UnmarshalQuotaDetails unmarshals an instance of QuotaDetails from the specified map of raw messages.
func UnmarshalQuotaDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QuotaDetails)
	err = core.UnmarshalPrimitive(m, "storage_bytes", &obj.StorageBytes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "traffic_bytes", &obj.TrafficBytes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RemoteAPIImage : Information about an image, in a format consistent with the Docker Remote API format.
type RemoteAPIImage struct {
	ConfigurationIssueCount *int64 `json:"ConfigurationIssueCount,omitempty"`

	Created *int64 `json:"Created,omitempty"`

	DigestTags map[string][]string `json:"DigestTags,omitempty"`

	ExemptIssueCount *int64 `json:"ExemptIssueCount,omitempty"`

	ID *string `json:"Id,omitempty"`

	IssueCount *int64 `json:"IssueCount,omitempty"`

	Labels map[string]string `json:"Labels,omitempty"`

	ManifestType *string `json:"ManifestType,omitempty"`

	ParentID *string `json:"ParentId,omitempty"`

	RepoDigests []string `json:"RepoDigests,omitempty"`

	RepoTags []string `json:"RepoTags,omitempty"`

	Size *int64 `json:"Size,omitempty"`

	VirtualSize *int64 `json:"VirtualSize,omitempty"`

	VulnerabilityCount *int64 `json:"VulnerabilityCount,omitempty"`

	Vulnerable *string `json:"Vulnerable,omitempty"`
}

// UnmarshalRemoteAPIImage unmarshals an instance of RemoteAPIImage from the specified map of raw messages.
func UnmarshalRemoteAPIImage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RemoteAPIImage)
	err = core.UnmarshalPrimitive(m, "ConfigurationIssueCount", &obj.ConfigurationIssueCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "DigestTags", &obj.DigestTags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ExemptIssueCount", &obj.ExemptIssueCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "IssueCount", &obj.IssueCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ManifestType", &obj.ManifestType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ParentId", &obj.ParentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "RepoDigests", &obj.RepoDigests)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "RepoTags", &obj.RepoTags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Size", &obj.Size)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "VirtualSize", &obj.VirtualSize)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "VulnerabilityCount", &obj.VulnerabilityCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Vulnerable", &obj.Vulnerable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RestoreImageOptions : The RestoreImage options.
type RestoreImageOptions struct {
	// The name of the image that you want to restore, in the format &lt;REPOSITORY&gt;:&lt;TAG&gt;. Run `ibmcloud cr
	// trash-list` or call the `GET /trash/json` endpoint to review images that are in the trash.
	Image *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewRestoreImageOptions : Instantiate RestoreImageOptions
func (*ContainerRegistryV1) NewRestoreImageOptions(image string) *RestoreImageOptions {
	return &RestoreImageOptions{
		Image: core.StringPtr(image),
	}
}

// SetImage : Allow user to set Image
func (options *RestoreImageOptions) SetImage(image string) *RestoreImageOptions {
	options.Image = core.StringPtr(image)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *RestoreImageOptions) SetHeaders(param map[string]string) *RestoreImageOptions {
	options.Headers = param
	return options
}

// RestoreResult : The result of restoring tags for a digest. In a successful request the digest is always restored, and zero or more of
// its tags may be restored.
type RestoreResult struct {
	// Successful is a list of tags that were restored.
	Successful []string `json:"successful,omitempty"`

	// Unsuccessful is a list of tags that were not restored because of a conflict.
	Unsuccessful []string `json:"unsuccessful,omitempty"`
}

// UnmarshalRestoreResult unmarshals an instance of RestoreResult from the specified map of raw messages.
func UnmarshalRestoreResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RestoreResult)
	err = core.UnmarshalPrimitive(m, "successful", &obj.Successful)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unsuccessful", &obj.Unsuccessful)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RestoreTagsOptions : The RestoreTags options.
type RestoreTagsOptions struct {
	// The full IBM Cloud registry digest reference for the digest that you want to restore such as
	// `icr.io/namespace/repo@sha256:a9be...`. Call the `GET /trash/json` endpoint to review digests that are in the trash
	// and their tags in the same repository.
	Digest *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewRestoreTagsOptions : Instantiate RestoreTagsOptions
func (*ContainerRegistryV1) NewRestoreTagsOptions(digest string) *RestoreTagsOptions {
	return &RestoreTagsOptions{
		Digest: core.StringPtr(digest),
	}
}

// SetDigest : Allow user to set Digest
func (options *RestoreTagsOptions) SetDigest(digest string) *RestoreTagsOptions {
	options.Digest = core.StringPtr(digest)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *RestoreTagsOptions) SetHeaders(param map[string]string) *RestoreTagsOptions {
	options.Headers = param
	return options
}

// RetentionPolicy : A document that contains the image retention settings for a namespace.
type RetentionPolicy struct {
	// Determines how many images will be retained for each repository when the retention policy is executed. The value -1
	// denotes 'Unlimited' (all images are retained).
	ImagesPerRepo *int64 `json:"images_per_repo,omitempty"`

	// The namespace to which the retention policy is attached.
	Namespace *string `json:"namespace" validate:"required"`

	// Determines if untagged images are retained when executing the retention policy. This is false by default meaning
	// untagged images will be deleted when the policy is executed.
	RetainUntagged *bool `json:"retain_untagged,omitempty"`
}

// NewRetentionPolicy : Instantiate RetentionPolicy (Generic Model Constructor)
func (*ContainerRegistryV1) NewRetentionPolicy(namespace string) (model *RetentionPolicy, err error) {
	model = &RetentionPolicy{
		Namespace: core.StringPtr(namespace),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRetentionPolicy unmarshals an instance of RetentionPolicy from the specified map of raw messages.
func UnmarshalRetentionPolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RetentionPolicy)
	err = core.UnmarshalPrimitive(m, "images_per_repo", &obj.ImagesPerRepo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "retain_untagged", &obj.RetainUntagged)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RootFs : RootFS contains information about the root filesystem of a container image.
type RootFs struct {
	// Descriptor for the base layer in the image.
	BaseLayer *string `json:"BaseLayer,omitempty"`

	// Descriptors for each layer in the image.
	Layers []string `json:"Layers,omitempty"`

	// The type of filesystem.
	Type *string `json:"Type,omitempty"`
}

// UnmarshalRootFs unmarshals an instance of RootFs from the specified map of raw messages.
func UnmarshalRootFs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RootFs)
	err = core.UnmarshalPrimitive(m, "BaseLayer", &obj.BaseLayer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Layers", &obj.Layers)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Type", &obj.Type)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SetRetentionPolicyOptions : The SetRetentionPolicy options.
type SetRetentionPolicyOptions struct {
	// The namespace to which the retention policy is attached.
	Namespace *string `validate:"required"`

	// Determines how many images will be retained for each repository when the retention policy is executed. The value -1
	// denotes 'Unlimited' (all images are retained).
	ImagesPerRepo *int64

	// Determines if untagged images are retained when executing the retention policy. This is false by default meaning
	// untagged images will be deleted when the policy is executed.
	RetainUntagged *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetRetentionPolicyOptions : Instantiate SetRetentionPolicyOptions
func (*ContainerRegistryV1) NewSetRetentionPolicyOptions(namespace string) *SetRetentionPolicyOptions {
	return &SetRetentionPolicyOptions{
		Namespace: core.StringPtr(namespace),
	}
}

// SetNamespace : Allow user to set Namespace
func (options *SetRetentionPolicyOptions) SetNamespace(namespace string) *SetRetentionPolicyOptions {
	options.Namespace = core.StringPtr(namespace)
	return options
}

// SetImagesPerRepo : Allow user to set ImagesPerRepo
func (options *SetRetentionPolicyOptions) SetImagesPerRepo(imagesPerRepo int64) *SetRetentionPolicyOptions {
	options.ImagesPerRepo = core.Int64Ptr(imagesPerRepo)
	return options
}

// SetRetainUntagged : Allow user to set RetainUntagged
func (options *SetRetentionPolicyOptions) SetRetainUntagged(retainUntagged bool) *SetRetentionPolicyOptions {
	options.RetainUntagged = core.BoolPtr(retainUntagged)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetRetentionPolicyOptions) SetHeaders(param map[string]string) *SetRetentionPolicyOptions {
	options.Headers = param
	return options
}

// TagImageOptions : The TagImage options.
type TagImageOptions struct {
	// The name of the image that you want to create a new tag for, in the format &lt;REPOSITORY&gt;:&lt;TAG&gt;. Run
	// `ibmcloud cr images` or call the `GET /images/json` endpoint to review images that are in the registry.
	Fromimage *string `validate:"required"`

	// The new tag for the image, in the format &lt;REPOSITORY&gt;:&lt;TAG&gt;.
	Toimage *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewTagImageOptions : Instantiate TagImageOptions
func (*ContainerRegistryV1) NewTagImageOptions(fromimage string, toimage string) *TagImageOptions {
	return &TagImageOptions{
		Fromimage: core.StringPtr(fromimage),
		Toimage:   core.StringPtr(toimage),
	}
}

// SetFromimage : Allow user to set Fromimage
func (options *TagImageOptions) SetFromimage(fromimage string) *TagImageOptions {
	options.Fromimage = core.StringPtr(fromimage)
	return options
}

// SetToimage : Allow user to set Toimage
func (options *TagImageOptions) SetToimage(toimage string) *TagImageOptions {
	options.Toimage = core.StringPtr(toimage)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *TagImageOptions) SetHeaders(param map[string]string) *TagImageOptions {
	options.Headers = param
	return options
}

// Trash : Details of the tags and days until expiry.
type Trash struct {
	DaysUntilExpiry *int64 `json:"daysUntilExpiry,omitempty"`

	Tags []string `json:"tags,omitempty"`
}

// UnmarshalTrash unmarshals an instance of Trash from the specified map of raw messages.
func UnmarshalTrash(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Trash)
	err = core.UnmarshalPrimitive(m, "daysUntilExpiry", &obj.DaysUntilExpiry)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateAuthOptions : The UpdateAuth options.
type UpdateAuthOptions struct {
	// Enable role based authorization when authenticating with IBM Cloud IAM.
	IamAuthz *bool

	// Restrict account to only be able to push and pull images over private connections.
	PrivateOnly *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateAuthOptions : Instantiate UpdateAuthOptions
func (*ContainerRegistryV1) NewUpdateAuthOptions() *UpdateAuthOptions {
	return &UpdateAuthOptions{}
}

// SetIamAuthz : Allow user to set IamAuthz
func (options *UpdateAuthOptions) SetIamAuthz(iamAuthz bool) *UpdateAuthOptions {
	options.IamAuthz = core.BoolPtr(iamAuthz)
	return options
}

// SetPrivateOnly : Allow user to set PrivateOnly
func (options *UpdateAuthOptions) SetPrivateOnly(privateOnly bool) *UpdateAuthOptions {
	options.PrivateOnly = core.BoolPtr(privateOnly)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAuthOptions) SetHeaders(param map[string]string) *UpdateAuthOptions {
	options.Headers = param
	return options
}

// UpdatePlansOptions : The UpdatePlans options.
type UpdatePlansOptions struct {
	Plan *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdatePlansOptions : Instantiate UpdatePlansOptions
func (*ContainerRegistryV1) NewUpdatePlansOptions() *UpdatePlansOptions {
	return &UpdatePlansOptions{}
}

// SetPlan : Allow user to set Plan
func (options *UpdatePlansOptions) SetPlan(plan string) *UpdatePlansOptions {
	options.Plan = core.StringPtr(plan)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdatePlansOptions) SetHeaders(param map[string]string) *UpdatePlansOptions {
	options.Headers = param
	return options
}

// UpdateQuotaOptions : The UpdateQuota options.
type UpdateQuotaOptions struct {
	// Storage quota in megabytes. The value -1 denotes "Unlimited".
	StorageMegabytes *int64

	// Traffic quota in megabytes. The value -1 denotes "Unlimited".
	TrafficMegabytes *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateQuotaOptions : Instantiate UpdateQuotaOptions
func (*ContainerRegistryV1) NewUpdateQuotaOptions() *UpdateQuotaOptions {
	return &UpdateQuotaOptions{}
}

// SetStorageMegabytes : Allow user to set StorageMegabytes
func (options *UpdateQuotaOptions) SetStorageMegabytes(storageMegabytes int64) *UpdateQuotaOptions {
	options.StorageMegabytes = core.Int64Ptr(storageMegabytes)
	return options
}

// SetTrafficMegabytes : Allow user to set TrafficMegabytes
func (options *UpdateQuotaOptions) SetTrafficMegabytes(trafficMegabytes int64) *UpdateQuotaOptions {
	options.TrafficMegabytes = core.Int64Ptr(trafficMegabytes)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateQuotaOptions) SetHeaders(param map[string]string) *UpdateQuotaOptions {
	options.Headers = param
	return options
}

// UpdateSettingsOptions : The UpdateSettings options.
type UpdateSettingsOptions struct {
	// Opt in to IBM Cloud Container Registry publishing platform metrics.
	PlatformMetrics *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSettingsOptions : Instantiate UpdateSettingsOptions
func (*ContainerRegistryV1) NewUpdateSettingsOptions() *UpdateSettingsOptions {
	return &UpdateSettingsOptions{}
}

// SetPlatformMetrics : Allow user to set PlatformMetrics
func (options *UpdateSettingsOptions) SetPlatformMetrics(platformMetrics bool) *UpdateSettingsOptions {
	options.PlatformMetrics = core.BoolPtr(platformMetrics)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSettingsOptions) SetHeaders(param map[string]string) *UpdateSettingsOptions {
	options.Headers = param
	return options
}
