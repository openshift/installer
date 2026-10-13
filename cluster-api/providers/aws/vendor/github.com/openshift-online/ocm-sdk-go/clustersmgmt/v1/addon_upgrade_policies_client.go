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

// AddonUpgradePoliciesClient is the client of the 'addon_upgrade_policies' resource.
//
// Manages the collection of addon upgrade policies of a cluster.
type AddonUpgradePoliciesClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddonUpgradePoliciesClient creates a new client for the 'addon_upgrade_policies'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddonUpgradePoliciesClient(transport http.RoundTripper, path string) *AddonUpgradePoliciesClient {
	return &AddonUpgradePoliciesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new addon upgrade policy to the cluster.
func (c *AddonUpgradePoliciesClient) Add() *AddonUpgradePoliciesAddRequest {
	return &AddonUpgradePoliciesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of addon upgrade policies.
func (c *AddonUpgradePoliciesClient) List() *AddonUpgradePoliciesListRequest {
	return &AddonUpgradePoliciesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AddonUpgradePolicy returns the target 'addon_upgrade_policy' resource for the given identifier.
//
// Reference to the service that manages an specific addon upgrade policy.
func (c *AddonUpgradePoliciesClient) AddonUpgradePolicy(id string) *AddonUpgradePolicyClient {
	return NewAddonUpgradePolicyClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// AddonUpgradePoliciesAddRequest is the request for the 'add' method.
type AddonUpgradePoliciesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddonUpgradePolicy
}

// Parameter adds a query parameter.
func (r *AddonUpgradePoliciesAddRequest) Parameter(name string, value interface{}) *AddonUpgradePoliciesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonUpgradePoliciesAddRequest) Header(name string, value interface{}) *AddonUpgradePoliciesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonUpgradePoliciesAddRequest) Impersonate(user string) *AddonUpgradePoliciesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the upgrade policy.
func (r *AddonUpgradePoliciesAddRequest) Body(value *AddonUpgradePolicy) *AddonUpgradePoliciesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonUpgradePoliciesAddRequest) Send() (result *AddonUpgradePoliciesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonUpgradePoliciesAddRequest) SendContext(ctx context.Context) (result *AddonUpgradePoliciesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddonUpgradePoliciesAddRequest(r, buffer)
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
	result = &AddonUpgradePoliciesAddResponse{}
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
	err = readAddonUpgradePoliciesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonUpgradePoliciesAddResponse is the response for the 'add' method.
type AddonUpgradePoliciesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddonUpgradePolicy
}

// Status returns the response status code.
func (r *AddonUpgradePoliciesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonUpgradePoliciesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonUpgradePoliciesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the upgrade policy.
func (r *AddonUpgradePoliciesAddResponse) Body() *AddonUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the upgrade policy.
func (r *AddonUpgradePoliciesAddResponse) GetBody() (value *AddonUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddonUpgradePoliciesListRequest is the request for the 'list' method.
type AddonUpgradePoliciesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *AddonUpgradePoliciesListRequest) Parameter(name string, value interface{}) *AddonUpgradePoliciesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonUpgradePoliciesListRequest) Header(name string, value interface{}) *AddonUpgradePoliciesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonUpgradePoliciesListRequest) Impersonate(user string) *AddonUpgradePoliciesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddonUpgradePoliciesListRequest) Page(value int) *AddonUpgradePoliciesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *AddonUpgradePoliciesListRequest) Size(value int) *AddonUpgradePoliciesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonUpgradePoliciesListRequest) Send() (result *AddonUpgradePoliciesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonUpgradePoliciesListRequest) SendContext(ctx context.Context) (result *AddonUpgradePoliciesListResponse, err error) {
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
	result = &AddonUpgradePoliciesListResponse{}
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
	err = readAddonUpgradePoliciesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonUpgradePoliciesListResponse is the response for the 'list' method.
type AddonUpgradePoliciesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *AddonUpgradePolicyList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *AddonUpgradePoliciesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonUpgradePoliciesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonUpgradePoliciesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of addon upgrade policy.
func (r *AddonUpgradePoliciesListResponse) Items() *AddonUpgradePolicyList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of addon upgrade policy.
func (r *AddonUpgradePoliciesListResponse) GetItems() (value *AddonUpgradePolicyList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddonUpgradePoliciesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddonUpgradePoliciesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *AddonUpgradePoliciesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *AddonUpgradePoliciesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *AddonUpgradePoliciesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *AddonUpgradePoliciesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
