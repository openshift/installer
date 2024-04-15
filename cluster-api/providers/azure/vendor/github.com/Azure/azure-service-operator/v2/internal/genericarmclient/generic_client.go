/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genericarmclient

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	azcoreruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/internal/metrics"
	"github.com/Azure/azure-service-operator/v2/internal/version"
)

const CreatePollerID = "GenericClient.CreateOrUpdateByID"
const DeletePollerID = "GenericClient.DeleteByID"

// NOTE: All of these methods (and types) were adapted from
// https://github.com/Azure/azure-sdk-for-go/blob/sdk/resources/armresources/v0.3.0/sdk/resources/armresources/zz_generated_resources_client.go
// which was then moved to here: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/resourcemanager/resources/armresources/client.go

type GenericClient struct {
	endpoint string
	pl       runtime.Pipeline
	creds    azcore.TokenCredential
	opts     *arm.ClientOptions
}

// TODO: Need to do retryAfter detection in each call?

type GenericClientOptions struct {
	HttpClient *http.Client
	Metrics    *metrics.ARMClientMetrics
	UserAgent  string
}

// NewGenericClient creates a new instance of GenericClient
func NewGenericClient(
	cloudCfg cloud.Configuration,
	creds azcore.TokenCredential,
	options *GenericClientOptions,
) (*GenericClient, error) {
	rmConfig, ok := cloudCfg.Services[cloud.ResourceManager]
	if !ok {
		return nil, errors.Errorf("provided cloud missing %q entry", cloud.ResourceManager)
	}
	if rmConfig.Endpoint == "" {
		return nil, errors.New("provided cloud missing resourceManager.Endpoint entry")
	}

	if options == nil {
		options = &GenericClientOptions{}
	}

	ua := options.UserAgent
	if ua == "" {
		ua = userAgent
	}

	opts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloudCfg,
			Retry: policy.RetryOptions{
				MaxRetries: -1, // Have to use a value less than 0 means no retries (0 does NOT, 0 gets you 3...)
			},
			PerCallPolicies: []policy.Policy{
				NewUserAgentPolicy(ua),
			},
		},
		// Disabled here because we don't want the default configuration, it polls for 5+ minutes which is
		// far too long to block an operator.
		DisableRPRegistration: true,
	}

	// We assign this HTTPClient like this because if we actually set it to nil, due to the way
	// go interfaces wrap values, the subsequent if httpClient == nil check returns false (even though
	// the value IN the interface IS nil).
	if options.HttpClient != nil {
		opts.Transport = options.HttpClient
	} else {
		opts.Transport = defaultHttpClient
	}

	rpRegistrationPolicy, err := NewRPRegistrationPolicy(
		creds,
		&opts.ClientOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create rp registration policy")
	}

	opts.PerCallPolicies = append([]policy.Policy{rpRegistrationPolicy}, opts.PerCallPolicies...)
	if options.Metrics != nil {
		opts.PerCallPolicies = append(opts.PerCallPolicies, metrics.NewMetricsPolicy(options.Metrics))
	}
	pipeline, err := armruntime.NewPipeline("generic", version.BuildVersion, creds, runtime.PipelineOptions{}, opts)
	if err != nil {
		return nil, err
	}

	return &GenericClient{
		endpoint: rmConfig.Endpoint,
		pl:       pipeline,
		creds:    creds,
		opts:     opts,
	}, nil

}

// Creds returns the credentials used by this client
func (client *GenericClient) Creds() azcore.TokenCredential {
	return client.creds
}

// ClientOptions returns the arm.ClientOptions used by this client. These options include
// the HTTP pipeline. If these options are used to create a new client, it will share the configured
// HTTP pipeline.
func (client *GenericClient) ClientOptions() *arm.ClientOptions {
	return client.opts
}

func (client *GenericClient) BeginCreateOrUpdateByID(
	ctx context.Context,
	resourceID string,
	apiVersion string,
	resource interface{}) (*PollerResponse[GenericResource], error) {
	// The linter doesn't realize that the response is closed in the course of
	// the autorest.NewPoller call below. Suppressing it as it is a false positive.
	// nolint:bodyclose
	resp, err := client.createOrUpdateByID(ctx, resourceID, apiVersion, resource)
	if err != nil {
		return nil, err
	}
	result := PollerResponse[GenericResource]{
		RawResponse:  resp,
		ID:           CreatePollerID,
		ErrorHandler: client.handleError,
	}

	pt, err := azcoreruntime.NewPoller[GenericResource](resp, client.pl, nil)
	if err != nil {
		return nil, err
	}
	result.Poller = pt
	return &result, nil
}

func (client *GenericClient) createOrUpdateByID(
	ctx context.Context,
	resourceID string,
	apiVersion string,
	resource interface{}) (*http.Response, error) {

	req, err := client.createOrUpdateByIDCreateRequest(ctx, resourceID, apiVersion, resource)
	if err != nil {
		return nil, err
	}

	resp, err := client.pl.Do(req)
	if err != nil {
		return resp, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, client.handleError(resp)
	}

	return resp, nil
}

// createOrUpdateByIDCreateRequest creates the CreateOrUpdateByID request.
func (client *GenericClient) createOrUpdateByIDCreateRequest(
	ctx context.Context,
	resourceID string,
	apiVersion string,
	resource interface{}) (*policy.Request, error) {

	if resourceID == "" {
		return nil, errors.New("parameter resourceID cannot be empty")
	}

	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, resource)
}

// handleError handles the CreateOrUpdateByID error response.
func (client *GenericClient) handleError(resp *http.Response) error {
	errType := NewCloudError(runtime.NewResponseError(resp))
	if err := runtime.UnmarshalAsJSON(resp, errType); err != nil {
		return runtime.NewResponseError(resp)
	}

	return errType
}

// GetByID - Gets a resource by ID.
// If the operation fails it returns the *CloudError error type.
func (client *GenericClient) GetByID(
	ctx context.Context,
	resourceID string,
	apiVersion string,
	resource interface{},
) (time.Duration, error) {
	req, err := client.getByIDCreateRequest(ctx, resourceID, apiVersion)
	if err != nil {
		return zeroDuration, err
	}
	// The linter doesn't realize that the response is closed in the course of
	// the getByIDHandleResponse call below. Suppressing it as it is a false positive.
	// nolint:bodyclose
	resp, err := client.pl.Do(req)
	retryAfter := GetRetryAfter(resp)
	if err != nil {
		return retryAfter, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return retryAfter, runtime.NewResponseError(resp)
	}
	return zeroDuration, client.getByIDHandleResponse(resp, resource)
}

// getByIDCreateRequest creates the GetByID request.
func (client *GenericClient) getByIDCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := "/{resourceId}"
	if resourceID == "" {
		return nil, errors.New("parameter resourceID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getByIDHandleResponse handles the GetByID response.
func (client *GenericClient) getByIDHandleResponse(resp *http.Response, resource interface{}) error {
	if err := runtime.UnmarshalAsJSON(resp, resource); err != nil {
		return err
	}
	return nil
}

// CheckExistenceByID - Heads a resource by ID.
// If the operation fails it returns the *CloudError error type.
func (client *GenericClient) CheckExistenceByID(
	ctx context.Context,
	resourceID string,
	apiVersion string,
) (bool, time.Duration, error) {
	retryAfter, err := client.checkExistenceByIDImpl(ctx, resourceID, apiVersion)
	switch {
	case IsNotFoundError(err):
		return false, retryAfter, nil
	case err != nil:
		return false, retryAfter, err
	default:
		return true, retryAfter, nil
	}
}

func (client *GenericClient) checkExistenceByIDImpl(
	ctx context.Context,
	resourceID string,
	apiVersion string,
) (time.Duration, error) {
	req, err := client.checkExistenceByIDCreateRequest(ctx, resourceID, apiVersion)
	if err != nil {
		return zeroDuration, err
	}
	// The linter doesn't realize that the response is closed as part of the pipeline
	// nolint:bodyclose
	resp, err := client.pl.Do(req)
	retryAfter := GetRetryAfter(resp)
	if err != nil {
		return retryAfter, err
	}
	if !runtime.HasStatusCode(resp, http.StatusNoContent, http.StatusNotFound) {
		return retryAfter, runtime.NewResponseError(resp)
	}
	return zeroDuration, nil
}

func (client *GenericClient) checkExistenceByIDCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := "/{resourceId}"
	if resourceID == "" {
		return nil, errors.New("parameter resourceID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
	req, err := runtime.NewRequest(ctx, http.MethodHead, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

type listPageResponse[T any] struct {
	// Value - The list of resources.
	Value []T `json:"value,omitempty"`

	// NextLink - The URI to fetch the next page of resources.
	NextLink *string `json:"nextLink,omitempty"`
}

func (p *listPageResponse[T]) More() bool {
	return p.NextLink != nil && len(*p.NextLink) > 0
}

func (p *listPageResponse[T]) NextPage(
	ctx context.Context,
	client *GenericClient,
	containerID string,
	apiVersion string,
) (*listPageResponse[T], error) {
	var req *policy.Request
	var err error
	if p == nil {
		req, err = client.listByContainerIDCreateRequest(ctx, containerID, apiVersion)
	} else {
		req, err = runtime.NewRequest(ctx, http.MethodGet, *p.NextLink)
	}
	if err != nil {
		return nil, err
	}

	// The linter doesn't realize that the response is closed in the course of
	// the runtime.UnmarshalAsJSON() call below. Suppressing it as it is a false positive.
	// nolint:bodyclose
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, runtime.NewResponseError(resp)
	}

	newPage := listPageResponse[T]{}
	err = runtime.UnmarshalAsJSON(resp, &newPage)
	if err != nil {
		return nil, err
	}

	return &newPage, nil
}

// ListByContainerID returns all the resources of a given type under a specified parent.
// If the operation fails it returns the *CloudError error type.
// ctx is the context of the request.
// client is the GenericClient to use for the request (can't declare generic methods, so this is standalone).
// containerID is the unique ID of the container in which the resources are contained.
// apiVersion is the API version to use for the request.
// createResource is a function that returns a new instance of the resource type.
func ListByContainerID[T any](
	ctx context.Context,
	client *GenericClient,
	containerID string,
	apiVersion string,
) ([]T, error) {
	pager := runtime.NewPager(
		runtime.PagingHandler[listPageResponse[T]]{
			More: func(page listPageResponse[T]) bool {
				// We have more if we have a link to follow
				return page.More()
			},
			Fetcher: func(ctx context.Context, page *listPageResponse[T]) (listPageResponse[T], error) {
				nextPage, err := page.NextPage(ctx, client, containerID, apiVersion)
				if err != nil {
					return listPageResponse[T]{}, err
				}

				return *nextPage, nil
			},
		})

	var result []T
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		result = append(result, page.Value...)
	}

	return result, nil
}

// listByParentIDCreateRequest creates the ListByContainerID request.
// ctx is the context of the request.
// containerID is the unique ID of the container in which the resources are contained.
// apiVersion is the API version to use for the request.
func (client *GenericClient) listByContainerIDCreateRequest(
	ctx context.Context,
	containerID string,
	apiVersion string,
) (*policy.Request, error) {
	urlPath := "/{containerId}"
	if containerID == "" {
		return nil, errors.New("parameter containerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerId}", containerID)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// BeginDeleteByID - Deletes a resource by ID.
// If the operation fails it returns the *CloudError error type.
func (client *GenericClient) BeginDeleteByID(ctx context.Context, resourceID string, apiVersion string) (*PollerResponse[GenericDeleteResponse], error) {
	// The linter doesn't realize that the response is closed in the course of
	// the autorest.NewPoller call below. Suppressing it as it is a false positive.
	// nolint:bodyclose
	resp, err := client.deleteByID(ctx, resourceID, apiVersion)
	if err != nil {
		return nil, err
	}

	result := PollerResponse[GenericDeleteResponse]{
		RawResponse:  resp,
		ID:           DeletePollerID,
		ErrorHandler: client.handleError,
	}
	pt, err := azcoreruntime.NewPoller[GenericDeleteResponse](resp, client.pl, nil)
	if err != nil {
		return nil, err
	}
	result.Poller = pt
	return &result, nil
}

// DeleteByID - Deletes a resource by ID.
// If the operation fails it returns the *CloudError error type.
func (client *GenericClient) deleteByID(ctx context.Context, resourceID string, apiVersion string) (*http.Response, error) {
	req, err := client.deleteByIDCreateRequest(ctx, resourceID, apiVersion)
	if err != nil {
		return nil, err
	}

	resp, err := client.pl.Do(req)
	if err != nil {
		return resp, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteByIDCreateRequest creates the DeleteByID request.
func (client *GenericClient) deleteByIDCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := "/{resourceId}"
	if resourceID == "" {
		return nil, errors.New("parameter resourceID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

func (client *GenericClient) CheckExistenceWithGetByID(ctx context.Context, resourceID string, apiVersion string) (bool, time.Duration, error) {
	if resourceID == "" {
		return false, zeroDuration, errors.New("parameter resourceID cannot be empty")
	}

	ignored := struct{}{}
	retryAfter, err := client.GetByID(ctx, resourceID, apiVersion, &ignored)

	switch {
	case IsNotFoundError(err):
		return false, retryAfter, nil
	case err != nil:
		return false, retryAfter, err
	default:
		return true, retryAfter, nil
	}
}

func IsNotFoundError(err error) bool {
	var typedError *azcore.ResponseError
	if errors.As(err, &typedError) {
		if typedError.StatusCode == http.StatusNotFound {
			return true
		}
	}

	return false
}

func (client *GenericClient) ResumeDeletePoller(id string) *PollerResponse[GenericDeleteResponse] {
	return &PollerResponse[GenericDeleteResponse]{ID: id, ErrorHandler: client.handleError}
}

func (client *GenericClient) ResumeCreatePoller(id string) *PollerResponse[GenericResource] {
	return &PollerResponse[GenericResource]{ID: id, ErrorHandler: client.handleError}
}
