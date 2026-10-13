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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ExternalAuthConfigClient is the client of the 'external_auth_config' resource.
//
// Manages the external authentication configuration for a ROSA HCP cluster.
type ExternalAuthConfigClient struct {
	transport http.RoundTripper
	path      string
}

// NewExternalAuthConfigClient creates a new client for the 'external_auth_config'
// resource using the given transport to send the requests and receive the
// responses.
func NewExternalAuthConfigClient(transport http.RoundTripper, path string) *ExternalAuthConfigClient {
	return &ExternalAuthConfigClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
func (c *ExternalAuthConfigClient) Get() *ExternalAuthConfigGetRequest {
	return &ExternalAuthConfigGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ExternalAuths returns the target 'external_auths' resource.
//
// Reference to the resource that manages the collection of ExternalAuths resources.
func (c *ExternalAuthConfigClient) ExternalAuths() *ExternalAuthsClient {
	return NewExternalAuthsClient(
		c.transport,
		path.Join(c.path, "external_auths"),
	)
}

// ExternalAuthConfigPollRequest is the request for the Poll method.
type ExternalAuthConfigPollRequest struct {
	request    *ExternalAuthConfigGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ExternalAuthConfigPollRequest) Parameter(name string, value interface{}) *ExternalAuthConfigPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ExternalAuthConfigPollRequest) Header(name string, value interface{}) *ExternalAuthConfigPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ExternalAuthConfigPollRequest) Interval(value time.Duration) *ExternalAuthConfigPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ExternalAuthConfigPollRequest) Status(value int) *ExternalAuthConfigPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ExternalAuthConfigPollRequest) Predicate(value func(*ExternalAuthConfigGetResponse) bool) *ExternalAuthConfigPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ExternalAuthConfigGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ExternalAuthConfigPollRequest) StartContext(ctx context.Context) (response *ExternalAuthConfigPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ExternalAuthConfigPollResponse{
			response: result.(*ExternalAuthConfigGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ExternalAuthConfigPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ExternalAuthConfigPollResponse is the response for the Poll method.
type ExternalAuthConfigPollResponse struct {
	response *ExternalAuthConfigGetResponse
}

// Status returns the response status code.
func (r *ExternalAuthConfigPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ExternalAuthConfigPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ExternalAuthConfigPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ExternalAuthConfigPollResponse) Body() *ExternalAuthConfig {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ExternalAuthConfigPollResponse) GetBody() (value *ExternalAuthConfig, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ExternalAuthConfigClient) Poll() *ExternalAuthConfigPollRequest {
	return &ExternalAuthConfigPollRequest{
		request: c.Get(),
	}
}

// ExternalAuthConfigGetRequest is the request for the 'get' method.
type ExternalAuthConfigGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ExternalAuthConfigGetRequest) Parameter(name string, value interface{}) *ExternalAuthConfigGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ExternalAuthConfigGetRequest) Header(name string, value interface{}) *ExternalAuthConfigGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ExternalAuthConfigGetRequest) Impersonate(user string) *ExternalAuthConfigGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ExternalAuthConfigGetRequest) Send() (result *ExternalAuthConfigGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ExternalAuthConfigGetRequest) SendContext(ctx context.Context) (result *ExternalAuthConfigGetResponse, err error) {
	query := helpers.CopyQuery(r.query)
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
	result = &ExternalAuthConfigGetResponse{}
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
	err = readExternalAuthConfigGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ExternalAuthConfigGetResponse is the response for the 'get' method.
type ExternalAuthConfigGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ExternalAuthConfig
}

// Status returns the response status code.
func (r *ExternalAuthConfigGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ExternalAuthConfigGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ExternalAuthConfigGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ExternalAuthConfigGetResponse) Body() *ExternalAuthConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ExternalAuthConfigGetResponse) GetBody() (value *ExternalAuthConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
