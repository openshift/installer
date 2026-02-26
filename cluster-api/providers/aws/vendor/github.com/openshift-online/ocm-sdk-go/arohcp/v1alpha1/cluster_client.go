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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ClusterClient is the client of the 'cluster' resource.
//
// Manages a specific cluster.
type ClusterClient struct {
	transport http.RoundTripper
	path      string
}

// NewClusterClient creates a new client for the 'cluster'
// resource using the given transport to send the requests and receive the
// responses.
func NewClusterClient(transport http.RoundTripper, path string) *ClusterClient {
	return &ClusterClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'async_delete' method.
//
// Deletes the cluster.
func (c *ClusterClient) Delete() *ClusterDeleteRequest {
	return &ClusterDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'async_update' method.
//
// Updates the cluster.
func (c *ClusterClient) Update() *ClusterUpdateRequest {
	return &ClusterUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the cluster.
func (c *ClusterClient) Get() *ClusterGetRequest {
	return &ClusterGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Autoscaler returns the target 'autoscaler' resource.
func (c *ClusterClient) Autoscaler() *AutoscalerClient {
	return NewAutoscalerClient(
		c.transport,
		path.Join(c.path, "autoscaler"),
	)
}

// ExternalAuthConfig returns the target 'external_auth_config' resource.
//
// Reference to the resource that manages the external authentication configuration.
func (c *ClusterClient) ExternalAuthConfig() *ExternalAuthConfigClient {
	return NewExternalAuthConfigClient(
		c.transport,
		path.Join(c.path, "external_auth_config"),
	)
}

// InflightChecks returns the target 'inflight_checks' resource.
//
// Reference to the resource that manages the collection of inflight checks.
func (c *ClusterClient) InflightChecks() *InflightChecksClient {
	return NewInflightChecksClient(
		c.transport,
		path.Join(c.path, "inflight_checks"),
	)
}

// NodePools returns the target 'node_pools' resource.
//
// Reference to the resource that manages the collection of node pool resources.
func (c *ClusterClient) NodePools() *NodePoolsClient {
	return NewNodePoolsClient(
		c.transport,
		path.Join(c.path, "node_pools"),
	)
}

// Status returns the target 'cluster_status' resource.
//
// Reference to the resource that manages the detailed status of the cluster.
func (c *ClusterClient) Status() *ClusterStatusClient {
	return NewClusterStatusClient(
		c.transport,
		path.Join(c.path, "status"),
	)
}

// ClusterPollRequest is the request for the Poll method.
type ClusterPollRequest struct {
	request    *ClusterGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ClusterPollRequest) Parameter(name string, value interface{}) *ClusterPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ClusterPollRequest) Header(name string, value interface{}) *ClusterPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ClusterPollRequest) Interval(value time.Duration) *ClusterPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ClusterPollRequest) Status(value int) *ClusterPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ClusterPollRequest) Predicate(value func(*ClusterGetResponse) bool) *ClusterPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ClusterGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ClusterPollRequest) StartContext(ctx context.Context) (response *ClusterPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ClusterPollResponse{
			response: result.(*ClusterGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ClusterPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ClusterPollResponse is the response for the Poll method.
type ClusterPollResponse struct {
	response *ClusterGetResponse
}

// Status returns the response status code.
func (r *ClusterPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ClusterPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ClusterPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ClusterPollResponse) Body() *Cluster {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ClusterPollResponse) GetBody() (value *Cluster, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ClusterClient) Poll() *ClusterPollRequest {
	return &ClusterPollRequest{
		request: c.Get(),
	}
}

// ClusterDeleteRequest is the request for the 'async_delete' method.
type ClusterDeleteRequest struct {
	transport   http.RoundTripper
	path        string
	query       url.Values
	header      http.Header
	bestEffort  *bool
	deprovision *bool
	dryRun      *bool
}

// Parameter adds a query parameter.
func (r *ClusterDeleteRequest) Parameter(name string, value interface{}) *ClusterDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterDeleteRequest) Header(name string, value interface{}) *ClusterDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterDeleteRequest) Impersonate(user string) *ClusterDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// BestEffort sets the value of the 'best_effort' parameter.
//
// BestEffort flag is used to check if the cluster deletion should be best-effort mode or not.
func (r *ClusterDeleteRequest) BestEffort(value bool) *ClusterDeleteRequest {
	r.bestEffort = &value
	return r
}

// Deprovision sets the value of the 'deprovision' parameter.
//
// If false it will only delete from OCM but not the actual cluster resources.
// false is only allowed for OCP clusters. true by default.
func (r *ClusterDeleteRequest) Deprovision(value bool) *ClusterDeleteRequest {
	r.deprovision = &value
	return r
}

// DryRun sets the value of the 'dry_run' parameter.
//
// Dry run flag is used to check if the operation can be completed, but won't delete.
func (r *ClusterDeleteRequest) DryRun(value bool) *ClusterDeleteRequest {
	r.dryRun = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterDeleteRequest) Send() (result *ClusterDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterDeleteRequest) SendContext(ctx context.Context) (result *ClusterDeleteResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.bestEffort != nil {
		helpers.AddValue(&query, "best_effort", *r.bestEffort)
	}
	if r.deprovision != nil {
		helpers.AddValue(&query, "deprovision", *r.deprovision)
	}
	if r.dryRun != nil {
		helpers.AddValue(&query, "dry_run", *r.dryRun)
	}
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "DELETE",
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
	result = &ClusterDeleteResponse{}
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
	return
}

// ClusterDeleteResponse is the response for the 'async_delete' method.
type ClusterDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ClusterDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ClusterUpdateRequest is the request for the 'async_update' method.
type ClusterUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Cluster
}

// Parameter adds a query parameter.
func (r *ClusterUpdateRequest) Parameter(name string, value interface{}) *ClusterUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterUpdateRequest) Header(name string, value interface{}) *ClusterUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterUpdateRequest) Impersonate(user string) *ClusterUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ClusterUpdateRequest) Body(value *Cluster) *ClusterUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterUpdateRequest) Send() (result *ClusterUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterUpdateRequest) SendContext(ctx context.Context) (result *ClusterUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeClusterAsyncUpdateRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "PATCH",
		URL:    uri,
		Header: header,
		Body:   io.NopCloser(buffer),
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &ClusterUpdateResponse{}
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
	err = readClusterAsyncUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ClusterUpdateResponse is the response for the 'async_update' method.
type ClusterUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Cluster
}

// Status returns the response status code.
func (r *ClusterUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ClusterUpdateResponse) Body() *Cluster {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ClusterUpdateResponse) GetBody() (value *Cluster, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ClusterGetRequest is the request for the 'get' method.
type ClusterGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ClusterGetRequest) Parameter(name string, value interface{}) *ClusterGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterGetRequest) Header(name string, value interface{}) *ClusterGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterGetRequest) Impersonate(user string) *ClusterGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterGetRequest) Send() (result *ClusterGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterGetRequest) SendContext(ctx context.Context) (result *ClusterGetResponse, err error) {
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
	result = &ClusterGetResponse{}
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
	err = readClusterGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ClusterGetResponse is the response for the 'get' method.
type ClusterGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Cluster
}

// Status returns the response status code.
func (r *ClusterGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ClusterGetResponse) Body() *Cluster {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ClusterGetResponse) GetBody() (value *Cluster, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
