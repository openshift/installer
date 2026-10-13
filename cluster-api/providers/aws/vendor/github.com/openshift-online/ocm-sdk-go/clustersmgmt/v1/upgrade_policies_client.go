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

// UpgradePoliciesClient is the client of the 'upgrade_policies' resource.
//
// Manages the collection of upgrade policies of a cluster.
type UpgradePoliciesClient struct {
	transport http.RoundTripper
	path      string
}

// NewUpgradePoliciesClient creates a new client for the 'upgrade_policies'
// resource using the given transport to send the requests and receive the
// responses.
func NewUpgradePoliciesClient(transport http.RoundTripper, path string) *UpgradePoliciesClient {
	return &UpgradePoliciesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new upgrade policy to the cluster.
func (c *UpgradePoliciesClient) Add() *UpgradePoliciesAddRequest {
	return &UpgradePoliciesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of upgrade policies.
func (c *UpgradePoliciesClient) List() *UpgradePoliciesListRequest {
	return &UpgradePoliciesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// UpgradePolicy returns the target 'upgrade_policy' resource for the given identifier.
//
// Reference to the service that manages an specific upgrade policy.
func (c *UpgradePoliciesClient) UpgradePolicy(id string) *UpgradePolicyClient {
	return NewUpgradePolicyClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// UpgradePoliciesAddRequest is the request for the 'add' method.
type UpgradePoliciesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *UpgradePolicy
}

// Parameter adds a query parameter.
func (r *UpgradePoliciesAddRequest) Parameter(name string, value interface{}) *UpgradePoliciesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *UpgradePoliciesAddRequest) Header(name string, value interface{}) *UpgradePoliciesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *UpgradePoliciesAddRequest) Impersonate(user string) *UpgradePoliciesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the upgrade policy.
func (r *UpgradePoliciesAddRequest) Body(value *UpgradePolicy) *UpgradePoliciesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *UpgradePoliciesAddRequest) Send() (result *UpgradePoliciesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *UpgradePoliciesAddRequest) SendContext(ctx context.Context) (result *UpgradePoliciesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeUpgradePoliciesAddRequest(r, buffer)
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
	result = &UpgradePoliciesAddResponse{}
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
	err = readUpgradePoliciesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// UpgradePoliciesAddResponse is the response for the 'add' method.
type UpgradePoliciesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *UpgradePolicy
}

// Status returns the response status code.
func (r *UpgradePoliciesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *UpgradePoliciesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *UpgradePoliciesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the upgrade policy.
func (r *UpgradePoliciesAddResponse) Body() *UpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the upgrade policy.
func (r *UpgradePoliciesAddResponse) GetBody() (value *UpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// UpgradePoliciesListRequest is the request for the 'list' method.
type UpgradePoliciesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *UpgradePoliciesListRequest) Parameter(name string, value interface{}) *UpgradePoliciesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *UpgradePoliciesListRequest) Header(name string, value interface{}) *UpgradePoliciesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *UpgradePoliciesListRequest) Impersonate(user string) *UpgradePoliciesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *UpgradePoliciesListRequest) Page(value int) *UpgradePoliciesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *UpgradePoliciesListRequest) Size(value int) *UpgradePoliciesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *UpgradePoliciesListRequest) Send() (result *UpgradePoliciesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *UpgradePoliciesListRequest) SendContext(ctx context.Context) (result *UpgradePoliciesListResponse, err error) {
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
	result = &UpgradePoliciesListResponse{}
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
	err = readUpgradePoliciesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// UpgradePoliciesListResponse is the response for the 'list' method.
type UpgradePoliciesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *UpgradePolicyList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *UpgradePoliciesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *UpgradePoliciesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *UpgradePoliciesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of upgrade policy.
func (r *UpgradePoliciesListResponse) Items() *UpgradePolicyList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of upgrade policy.
func (r *UpgradePoliciesListResponse) GetItems() (value *UpgradePolicyList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *UpgradePoliciesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *UpgradePoliciesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *UpgradePoliciesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *UpgradePoliciesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *UpgradePoliciesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *UpgradePoliciesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
