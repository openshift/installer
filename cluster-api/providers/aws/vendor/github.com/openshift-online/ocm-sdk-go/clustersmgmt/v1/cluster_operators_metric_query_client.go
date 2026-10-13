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
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ClusterOperatorsMetricQueryClient is the client of the 'cluster_operators_metric_query' resource.
//
// Provides detailed information about the Metrics of the Cluster
type ClusterOperatorsMetricQueryClient struct {
	transport http.RoundTripper
	path      string
}

// NewClusterOperatorsMetricQueryClient creates a new client for the 'cluster_operators_metric_query'
// resource using the given transport to send the requests and receive the
// responses.
func NewClusterOperatorsMetricQueryClient(transport http.RoundTripper, path string) *ClusterOperatorsMetricQueryClient {
	return &ClusterOperatorsMetricQueryClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
func (c *ClusterOperatorsMetricQueryClient) Get() *ClusterOperatorsMetricQueryGetRequest {
	return &ClusterOperatorsMetricQueryGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ClusterOperatorsMetricQueryPollRequest is the request for the Poll method.
type ClusterOperatorsMetricQueryPollRequest struct {
	request    *ClusterOperatorsMetricQueryGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ClusterOperatorsMetricQueryPollRequest) Parameter(name string, value interface{}) *ClusterOperatorsMetricQueryPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ClusterOperatorsMetricQueryPollRequest) Header(name string, value interface{}) *ClusterOperatorsMetricQueryPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ClusterOperatorsMetricQueryPollRequest) Interval(value time.Duration) *ClusterOperatorsMetricQueryPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ClusterOperatorsMetricQueryPollRequest) Status(value int) *ClusterOperatorsMetricQueryPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ClusterOperatorsMetricQueryPollRequest) Predicate(value func(*ClusterOperatorsMetricQueryGetResponse) bool) *ClusterOperatorsMetricQueryPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ClusterOperatorsMetricQueryGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ClusterOperatorsMetricQueryPollRequest) StartContext(ctx context.Context) (response *ClusterOperatorsMetricQueryPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ClusterOperatorsMetricQueryPollResponse{
			response: result.(*ClusterOperatorsMetricQueryGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ClusterOperatorsMetricQueryPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ClusterOperatorsMetricQueryPollResponse is the response for the Poll method.
type ClusterOperatorsMetricQueryPollResponse struct {
	response *ClusterOperatorsMetricQueryGetResponse
}

// Status returns the response status code.
func (r *ClusterOperatorsMetricQueryPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ClusterOperatorsMetricQueryPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ClusterOperatorsMetricQueryPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ClusterOperatorsMetricQueryPollResponse) Body() *ClusterOperatorsInfo {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ClusterOperatorsMetricQueryPollResponse) GetBody() (value *ClusterOperatorsInfo, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ClusterOperatorsMetricQueryClient) Poll() *ClusterOperatorsMetricQueryPollRequest {
	return &ClusterOperatorsMetricQueryPollRequest{
		request: c.Get(),
	}
}

// ClusterOperatorsMetricQueryGetRequest is the request for the 'get' method.
type ClusterOperatorsMetricQueryGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ClusterOperatorsMetricQueryGetRequest) Parameter(name string, value interface{}) *ClusterOperatorsMetricQueryGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterOperatorsMetricQueryGetRequest) Header(name string, value interface{}) *ClusterOperatorsMetricQueryGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterOperatorsMetricQueryGetRequest) Impersonate(user string) *ClusterOperatorsMetricQueryGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterOperatorsMetricQueryGetRequest) Send() (result *ClusterOperatorsMetricQueryGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterOperatorsMetricQueryGetRequest) SendContext(ctx context.Context) (result *ClusterOperatorsMetricQueryGetResponse, err error) {
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
	result = &ClusterOperatorsMetricQueryGetResponse{}
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
	err = readClusterOperatorsMetricQueryGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ClusterOperatorsMetricQueryGetResponse is the response for the 'get' method.
type ClusterOperatorsMetricQueryGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ClusterOperatorsInfo
}

// Status returns the response status code.
func (r *ClusterOperatorsMetricQueryGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterOperatorsMetricQueryGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterOperatorsMetricQueryGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ClusterOperatorsMetricQueryGetResponse) Body() *ClusterOperatorsInfo {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ClusterOperatorsMetricQueryGetResponse) GetBody() (value *ClusterOperatorsInfo, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
