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

// ControlPlaneUpgradePoliciesClient is the client of the 'control_plane_upgrade_policies' resource.
//
// Manages the collection of upgrade policies for the control plane of a cluster.
type ControlPlaneUpgradePoliciesClient struct {
	transport http.RoundTripper
	path      string
}

// NewControlPlaneUpgradePoliciesClient creates a new client for the 'control_plane_upgrade_policies'
// resource using the given transport to send the requests and receive the
// responses.
func NewControlPlaneUpgradePoliciesClient(transport http.RoundTripper, path string) *ControlPlaneUpgradePoliciesClient {
	return &ControlPlaneUpgradePoliciesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new upgrade policy to the control plane of the cluster.
func (c *ControlPlaneUpgradePoliciesClient) Add() *ControlPlaneUpgradePoliciesAddRequest {
	return &ControlPlaneUpgradePoliciesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of upgrade policies for the control plane.
func (c *ControlPlaneUpgradePoliciesClient) List() *ControlPlaneUpgradePoliciesListRequest {
	return &ControlPlaneUpgradePoliciesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ControlPlaneUpgradePolicy returns the target 'control_plane_upgrade_policy' resource for the given identifier.
//
// Reference to the service that manages an specific upgrade policy for the control plane.
func (c *ControlPlaneUpgradePoliciesClient) ControlPlaneUpgradePolicy(id string) *ControlPlaneUpgradePolicyClient {
	return NewControlPlaneUpgradePolicyClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ControlPlaneUpgradePoliciesAddRequest is the request for the 'add' method.
type ControlPlaneUpgradePoliciesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ControlPlaneUpgradePolicy
}

// Parameter adds a query parameter.
func (r *ControlPlaneUpgradePoliciesAddRequest) Parameter(name string, value interface{}) *ControlPlaneUpgradePoliciesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ControlPlaneUpgradePoliciesAddRequest) Header(name string, value interface{}) *ControlPlaneUpgradePoliciesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ControlPlaneUpgradePoliciesAddRequest) Impersonate(user string) *ControlPlaneUpgradePoliciesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the upgrade policy.
func (r *ControlPlaneUpgradePoliciesAddRequest) Body(value *ControlPlaneUpgradePolicy) *ControlPlaneUpgradePoliciesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ControlPlaneUpgradePoliciesAddRequest) Send() (result *ControlPlaneUpgradePoliciesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ControlPlaneUpgradePoliciesAddRequest) SendContext(ctx context.Context) (result *ControlPlaneUpgradePoliciesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeControlPlaneUpgradePoliciesAddRequest(r, buffer)
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
	result = &ControlPlaneUpgradePoliciesAddResponse{}
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
	err = readControlPlaneUpgradePoliciesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ControlPlaneUpgradePoliciesAddResponse is the response for the 'add' method.
type ControlPlaneUpgradePoliciesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ControlPlaneUpgradePolicy
}

// Status returns the response status code.
func (r *ControlPlaneUpgradePoliciesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ControlPlaneUpgradePoliciesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ControlPlaneUpgradePoliciesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the upgrade policy.
func (r *ControlPlaneUpgradePoliciesAddResponse) Body() *ControlPlaneUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the upgrade policy.
func (r *ControlPlaneUpgradePoliciesAddResponse) GetBody() (value *ControlPlaneUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ControlPlaneUpgradePoliciesListRequest is the request for the 'list' method.
type ControlPlaneUpgradePoliciesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *ControlPlaneUpgradePoliciesListRequest) Parameter(name string, value interface{}) *ControlPlaneUpgradePoliciesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ControlPlaneUpgradePoliciesListRequest) Header(name string, value interface{}) *ControlPlaneUpgradePoliciesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ControlPlaneUpgradePoliciesListRequest) Impersonate(user string) *ControlPlaneUpgradePoliciesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ControlPlaneUpgradePoliciesListRequest) Page(value int) *ControlPlaneUpgradePoliciesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *ControlPlaneUpgradePoliciesListRequest) Size(value int) *ControlPlaneUpgradePoliciesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ControlPlaneUpgradePoliciesListRequest) Send() (result *ControlPlaneUpgradePoliciesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ControlPlaneUpgradePoliciesListRequest) SendContext(ctx context.Context) (result *ControlPlaneUpgradePoliciesListResponse, err error) {
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
	result = &ControlPlaneUpgradePoliciesListResponse{}
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
	err = readControlPlaneUpgradePoliciesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ControlPlaneUpgradePoliciesListResponse is the response for the 'list' method.
type ControlPlaneUpgradePoliciesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ControlPlaneUpgradePolicyList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ControlPlaneUpgradePoliciesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ControlPlaneUpgradePoliciesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ControlPlaneUpgradePoliciesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of upgrade policy.
func (r *ControlPlaneUpgradePoliciesListResponse) Items() *ControlPlaneUpgradePolicyList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of upgrade policy.
func (r *ControlPlaneUpgradePoliciesListResponse) GetItems() (value *ControlPlaneUpgradePolicyList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ControlPlaneUpgradePoliciesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ControlPlaneUpgradePoliciesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *ControlPlaneUpgradePoliciesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *ControlPlaneUpgradePoliciesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *ControlPlaneUpgradePoliciesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *ControlPlaneUpgradePoliciesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
