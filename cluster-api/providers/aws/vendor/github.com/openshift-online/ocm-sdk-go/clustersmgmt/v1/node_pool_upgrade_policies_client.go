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
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// NodePoolUpgradePoliciesClient is the client of the 'node_pool_upgrade_policies' resource.
//
// Manages the collection of upgrade policies for the node pool of a cluster.
type NodePoolUpgradePoliciesClient struct {
	transport http.RoundTripper
	path      string
}

// NewNodePoolUpgradePoliciesClient creates a new client for the 'node_pool_upgrade_policies'
// resource using the given transport to send the requests and receive the
// responses.
func NewNodePoolUpgradePoliciesClient(transport http.RoundTripper, path string) *NodePoolUpgradePoliciesClient {
	return &NodePoolUpgradePoliciesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new upgrade policy to the node pool of the cluster.
func (c *NodePoolUpgradePoliciesClient) Add() *NodePoolUpgradePoliciesAddRequest {
	return &NodePoolUpgradePoliciesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of upgrade policies for the node pool.
func (c *NodePoolUpgradePoliciesClient) List() *NodePoolUpgradePoliciesListRequest {
	return &NodePoolUpgradePoliciesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// NodePoolUpgradePolicy returns the target 'node_pool_upgrade_policy' resource for the given identifier.
//
// Reference to the service that manages an specific upgrade policy for the node pool.
func (c *NodePoolUpgradePoliciesClient) NodePoolUpgradePolicy(id string) *NodePoolUpgradePolicyClient {
	return NewNodePoolUpgradePolicyClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// NodePoolUpgradePoliciesAddRequest is the request for the 'add' method.
type NodePoolUpgradePoliciesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *NodePoolUpgradePolicy
}

// Parameter adds a query parameter.
func (r *NodePoolUpgradePoliciesAddRequest) Parameter(name string, value interface{}) *NodePoolUpgradePoliciesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NodePoolUpgradePoliciesAddRequest) Header(name string, value interface{}) *NodePoolUpgradePoliciesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NodePoolUpgradePoliciesAddRequest) Impersonate(user string) *NodePoolUpgradePoliciesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the upgrade policy.
func (r *NodePoolUpgradePoliciesAddRequest) Body(value *NodePoolUpgradePolicy) *NodePoolUpgradePoliciesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NodePoolUpgradePoliciesAddRequest) Send() (result *NodePoolUpgradePoliciesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NodePoolUpgradePoliciesAddRequest) SendContext(ctx context.Context) (result *NodePoolUpgradePoliciesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeNodePoolUpgradePoliciesAddRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "POST",
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
	result = &NodePoolUpgradePoliciesAddResponse{}
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
	err = readNodePoolUpgradePoliciesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NodePoolUpgradePoliciesAddResponse is the response for the 'add' method.
type NodePoolUpgradePoliciesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *NodePoolUpgradePolicy
}

// Status returns the response status code.
func (r *NodePoolUpgradePoliciesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NodePoolUpgradePoliciesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NodePoolUpgradePoliciesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the upgrade policy.
func (r *NodePoolUpgradePoliciesAddResponse) Body() *NodePoolUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the upgrade policy.
func (r *NodePoolUpgradePoliciesAddResponse) GetBody() (value *NodePoolUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// NodePoolUpgradePoliciesListRequest is the request for the 'list' method.
type NodePoolUpgradePoliciesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *NodePoolUpgradePoliciesListRequest) Parameter(name string, value interface{}) *NodePoolUpgradePoliciesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NodePoolUpgradePoliciesListRequest) Header(name string, value interface{}) *NodePoolUpgradePoliciesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NodePoolUpgradePoliciesListRequest) Impersonate(user string) *NodePoolUpgradePoliciesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *NodePoolUpgradePoliciesListRequest) Page(value int) *NodePoolUpgradePoliciesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *NodePoolUpgradePoliciesListRequest) Size(value int) *NodePoolUpgradePoliciesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NodePoolUpgradePoliciesListRequest) Send() (result *NodePoolUpgradePoliciesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NodePoolUpgradePoliciesListRequest) SendContext(ctx context.Context) (result *NodePoolUpgradePoliciesListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.page != nil {
		helpers.AddValue(&query, "page", *r.page)
	}
	if r.size != nil {
		helpers.AddValue(&query, "size", *r.size)
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
	result = &NodePoolUpgradePoliciesListResponse{}
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
	err = readNodePoolUpgradePoliciesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NodePoolUpgradePoliciesListResponse is the response for the 'list' method.
type NodePoolUpgradePoliciesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *NodePoolUpgradePolicyList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *NodePoolUpgradePoliciesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NodePoolUpgradePoliciesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NodePoolUpgradePoliciesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of upgrade policy.
func (r *NodePoolUpgradePoliciesListResponse) Items() *NodePoolUpgradePolicyList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of upgrade policy.
func (r *NodePoolUpgradePoliciesListResponse) GetItems() (value *NodePoolUpgradePolicyList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *NodePoolUpgradePoliciesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *NodePoolUpgradePoliciesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *NodePoolUpgradePoliciesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *NodePoolUpgradePoliciesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *NodePoolUpgradePoliciesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *NodePoolUpgradePoliciesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
