/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genericarmclient

// NOTE: This file was copied from https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azcore/arm/runtime/policy_register_rp.go
// with modifications. The primary modification was to remove all retries and immediately report an error after registration.
// That error will be categorized and retried on by the operator.

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

	"github.com/Azure/azure-service-operator/v2/internal/version"
)

// NewRPRegistrationPolicy creates a policy object configured using the specified options.
// The policy controls whether an unregistered resource provider should automatically be
// registered. See https://aka.ms/rps-not-found for more information.
func NewRPRegistrationPolicy(cred azcore.TokenCredential, o *azpolicy.ClientOptions) (azpolicy.Policy, error) {
	conf, err := getConfiguration(o)
	if err != nil {
		return nil, err
	}
	authPolicy := armruntime.NewBearerTokenPolicy(cred, &armpolicy.BearerTokenOptions{Scopes: []string{conf.Audience + "/.default"}})
	p := &rpRegistrationPolicy{
		endpoint: conf.Endpoint,
		pipeline: runtime.NewPipeline("generic", version.BuildVersion, runtime.PipelineOptions{PerRetry: []azpolicy.Policy{authPolicy}}, o),
		options:  o,
	}
	return p, nil
}

type rpRegistrationPolicy struct {
	endpoint string
	pipeline runtime.Pipeline
	options  *azpolicy.ClientOptions
}

func (r *rpRegistrationPolicy) Do(req *azpolicy.Request) (*http.Response, error) {
	// The unique identifier returned by ARM when the required resource provider has not been registered
	const unregisteredRPCode = "MissingSubscriptionRegistration"
	var rp string
	var resp *http.Response
	var err error
	// make the original request
	resp, err = req.Next()
	// getting a 409 is the first indication that the RP might need to be registered, check error response
	if err != nil || resp.StatusCode != http.StatusConflict {
		return resp, err
	}
	var reqErr requestError
	if err = runtime.UnmarshalAsJSON(resp, &reqErr); err != nil {
		return resp, err
	}
	if reqErr.ServiceError == nil {
		return resp, errors.New("missing error information")
	}
	if !strings.EqualFold(reqErr.ServiceError.Code, unregisteredRPCode) {
		// not a 409 due to unregistered RP
		return resp, err
	}
	// RP needs to be registered.  start by getting the subscription ID from the original request
	subID, err := GetSubscription(req.Raw().URL.Path)
	if err != nil {
		return resp, err
	}
	// now get the RP from the error
	rp, err = getProvider(reqErr)
	if err != nil {
		return resp, err
	}
	// create client and make the registration request
	// we use the scheme and host from the original request
	rpOps := &providersOperations{
		p:     r.pipeline,
		u:     r.endpoint,
		subID: subID,
	}
	if _, err = rpOps.Register(req.Raw().Context(), rp); err != nil {
		return resp, err
	}
	return resp, fmt.Errorf("registering Resource Provider %s with subscription. Try again later", rp)
}

func getProvider(re requestError) (string, error) {
	if len(re.ServiceError.Details) > 0 {
		return re.ServiceError.Details[0].Target, nil
	}
	return "", errors.New("unexpected empty Details")
}

// minimal error definitions to simplify detection
type requestError struct {
	ServiceError *serviceError `json:"error"`
}

type serviceError struct {
	Code    string                `json:"code"`
	Details []serviceErrorDetails `json:"details"`
}

type serviceErrorDetails struct {
	Code   string `json:"code"`
	Target string `json:"target"`
}

///////////////////////////////////////////////////////////////////////////////////////////////
// the following code was copied from module armresources, providers.go and models.go
// only the minimum amount of code was copied to get this working and some edits were made.
///////////////////////////////////////////////////////////////////////////////////////////////

type providersOperations struct {
	p     runtime.Pipeline
	u     string
	subID string
}

// Register - Registers a subscription with a resource provider.
func (client *providersOperations) Register(ctx context.Context, resourceProviderNamespace string) (providerResponse, error) {
	req, err := client.registerCreateRequest(ctx, resourceProviderNamespace)
	if err != nil {
		return providerResponse{}, err
	}
	// The linter doesn't realize that the response is closed in the course of
	// the registerHandleResponse call below. Suppressing it as it is a false positive.
	// nolint:bodyclose
	resp, err := client.p.Do(req)
	if err != nil {
		return providerResponse{}, err
	}
	result, err := client.registerHandleResponse(resp)
	if err != nil {
		return providerResponse{}, err
	}
	return result, nil
}

// registerCreateRequest creates the Register request.
func (client *providersOperations) registerCreateRequest(ctx context.Context, resourceProviderNamespace string) (*azpolicy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}/register"
	urlPath = strings.ReplaceAll(urlPath, "{resourceProviderNamespace}", url.PathEscape(resourceProviderNamespace))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.Raw().URL.Query()
	query.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = query.Encode()
	return req, nil
}

// registerHandleResponse handles the Register response.
func (client *providersOperations) registerHandleResponse(resp *http.Response) (providerResponse, error) {
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return providerResponse{}, runtime.NewResponseError(resp)
	}
	result := providerResponse{RawResponse: resp}
	err := runtime.UnmarshalAsJSON(resp, &result.Provider)
	if err != nil {
		return providerResponse{}, err
	}
	return result, err
}

// ProviderResponse is the response envelope for operations that return a Provider type.
type providerResponse struct {
	// Resource provider information.
	Provider *provider

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// Provider - Resource provider information.
type provider struct {
	// The provider ID.
	ID *string `json:"id,omitempty"`

	// The namespace of the resource provider.
	Namespace *string `json:"namespace,omitempty"`

	// The registration policy of the resource provider.
	RegistrationPolicy *string `json:"registrationPolicy,omitempty"`

	// The registration state of the resource provider.
	RegistrationState *string `json:"registrationState,omitempty"`
}

func getConfiguration(o *azpolicy.ClientOptions) (cloud.ServiceConfiguration, error) {
	c := cloud.AzurePublic
	if !reflect.ValueOf(o.Cloud).IsZero() {
		c = o.Cloud
	}
	if conf, ok := c.Services[cloud.ResourceManager]; ok && conf.Endpoint != "" && conf.Audience != "" {
		return conf, nil
	} else {
		return conf, errors.New("provided Cloud field is missing Azure Resource Manager configuration")
	}
}
