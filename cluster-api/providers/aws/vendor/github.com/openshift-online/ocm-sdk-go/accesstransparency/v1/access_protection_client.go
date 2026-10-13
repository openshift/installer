/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// AccessProtectionClient is the client of the 'access_protection' resource.
//
// Manages the Access Protection resource.
type AccessProtectionClient struct {
	transport http.RoundTripper
	path      string
}

// NewAccessProtectionClient creates a new client for the 'access_protection'
// resource using the given transport to send the requests and receive the
// responses.
func NewAccessProtectionClient(transport http.RoundTripper, path string) *AccessProtectionClient {
	return &AccessProtectionClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves an Access Protection by organization/cluster/subscription query param.
func (c *AccessProtectionClient) Get() *AccessProtectionGetRequest {
	return &AccessProtectionGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AccessProtectionPollRequest is the request for the Poll method.
type AccessProtectionPollRequest struct {
	request    *AccessProtectionGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AccessProtectionPollRequest) Parameter(name string, value interface{}) *AccessProtectionPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AccessProtectionPollRequest) Header(name string, value interface{}) *AccessProtectionPollRequest {
	r.request.Header(name, value)
	return r
}

// ClusterId sets the value of the 'cluster_id' parameter for all the requests that
// will be used to retrieve the object.
//
// Check status by Cluter.
func (r *AccessProtectionPollRequest) ClusterId(value string) *AccessProtectionPollRequest {
	r.request.ClusterId(value)
	return r
}

// OrganizationId sets the value of the 'organization_id' parameter for all the requests that
// will be used to retrieve the object.
//
// Check status by Organization.
func (r *AccessProtectionPollRequest) OrganizationId(value string) *AccessProtectionPollRequest {
	r.request.OrganizationId(value)
	return r
}

// SubscriptionId sets the value of the 'subscription_id' parameter for all the requests that
// will be used to retrieve the object.
//
// Check status by Subscription.
func (r *AccessProtectionPollRequest) SubscriptionId(value string) *AccessProtectionPollRequest {
	r.request.SubscriptionId(value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AccessProtectionPollRequest) Interval(value time.Duration) *AccessProtectionPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AccessProtectionPollRequest) Status(value int) *AccessProtectionPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AccessProtectionPollRequest) Predicate(value func(*AccessProtectionGetResponse) bool) *AccessProtectionPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AccessProtectionGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AccessProtectionPollRequest) StartContext(ctx context.Context) (response *AccessProtectionPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AccessProtectionPollResponse{
			response: result.(*AccessProtectionGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AccessProtectionPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AccessProtectionPollResponse is the response for the Poll method.
type AccessProtectionPollResponse struct {
	response *AccessProtectionGetResponse
}

// Status returns the response status code.
func (r *AccessProtectionPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AccessProtectionPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AccessProtectionPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
//
// AccessProtection status response.
func (r *AccessProtectionPollResponse) Body() *AccessProtection {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// AccessProtection status response.
func (r *AccessProtectionPollResponse) GetBody() (value *AccessProtection, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AccessProtectionClient) Poll() *AccessProtectionPollRequest {
	return &AccessProtectionPollRequest{
		request: c.Get(),
	}
}

// AccessProtectionGetRequest is the request for the 'get' method.
type AccessProtectionGetRequest struct {
	transport      http.RoundTripper
	path           string
	query          url.Values
	header         http.Header
	clusterId      *string
	organizationId *string
	subscriptionId *string
}

// Parameter adds a query parameter.
func (r *AccessProtectionGetRequest) Parameter(name string, value interface{}) *AccessProtectionGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccessProtectionGetRequest) Header(name string, value interface{}) *AccessProtectionGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccessProtectionGetRequest) Impersonate(user string) *AccessProtectionGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// ClusterId sets the value of the 'cluster_id' parameter.
//
// Check status by Cluter.
func (r *AccessProtectionGetRequest) ClusterId(value string) *AccessProtectionGetRequest {
	r.clusterId = &value
	return r
}

// OrganizationId sets the value of the 'organization_id' parameter.
//
// Check status by Organization.
func (r *AccessProtectionGetRequest) OrganizationId(value string) *AccessProtectionGetRequest {
	r.organizationId = &value
	return r
}

// SubscriptionId sets the value of the 'subscription_id' parameter.
//
// Check status by Subscription.
func (r *AccessProtectionGetRequest) SubscriptionId(value string) *AccessProtectionGetRequest {
	r.subscriptionId = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccessProtectionGetRequest) Send() (result *AccessProtectionGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccessProtectionGetRequest) SendContext(ctx context.Context) (result *AccessProtectionGetResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.clusterId != nil {
		helpers.AddValue(&query, "clusterId", *r.clusterId)
	}
	if r.organizationId != nil {
		helpers.AddValue(&query, "organizationId", *r.organizationId)
	}
	if r.subscriptionId != nil {
		helpers.AddValue(&query, "subscriptionId", *r.subscriptionId)
	}
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "GET",
		URL:    uri,
		Header: header,
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &AccessProtectionGetResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	err = readAccessProtectionGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AccessProtectionGetResponse is the response for the 'get' method.
type AccessProtectionGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AccessProtection
}

// Status returns the response status code.
func (r *AccessProtectionGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccessProtectionGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccessProtectionGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// AccessProtection status response.
func (r *AccessProtectionGetResponse) Body() *AccessProtection {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// AccessProtection status response.
func (r *AccessProtectionGetResponse) GetBody() (value *AccessProtection, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
