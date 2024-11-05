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

// RegistryAllowlistsClient is the client of the 'registry_allowlists' resource.
//
// Manages the registry allowlists.
type RegistryAllowlistsClient struct {
	transport http.RoundTripper
	path      string
}

// NewRegistryAllowlistsClient creates a new client for the 'registry_allowlists'
// resource using the given transport to send the requests and receive the
// responses.
func NewRegistryAllowlistsClient(transport http.RoundTripper, path string) *RegistryAllowlistsClient {
	return &RegistryAllowlistsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new break registry allowlist.
func (c *RegistryAllowlistsClient) Add() *RegistryAllowlistsAddRequest {
	return &RegistryAllowlistsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of registry allowlists.
func (c *RegistryAllowlistsClient) List() *RegistryAllowlistsListRequest {
	return &RegistryAllowlistsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// RegistryAllowlist returns the target 'registry_allowlist' resource for the given identifier.
//
// Reference to the service that manages a specific registry allowlist.
func (c *RegistryAllowlistsClient) RegistryAllowlist(id string) *RegistryAllowlistClient {
	return NewRegistryAllowlistClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// RegistryAllowlistsAddRequest is the request for the 'add' method.
type RegistryAllowlistsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *RegistryAllowlist
}

// Parameter adds a query parameter.
func (r *RegistryAllowlistsAddRequest) Parameter(name string, value interface{}) *RegistryAllowlistsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RegistryAllowlistsAddRequest) Header(name string, value interface{}) *RegistryAllowlistsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RegistryAllowlistsAddRequest) Impersonate(user string) *RegistryAllowlistsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Data of the new registry allowlist.
func (r *RegistryAllowlistsAddRequest) Body(value *RegistryAllowlist) *RegistryAllowlistsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RegistryAllowlistsAddRequest) Send() (result *RegistryAllowlistsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RegistryAllowlistsAddRequest) SendContext(ctx context.Context) (result *RegistryAllowlistsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeRegistryAllowlistsAddRequest(r, buffer)
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
	result = &RegistryAllowlistsAddResponse{}
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
	err = readRegistryAllowlistsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// RegistryAllowlistsAddResponse is the response for the 'add' method.
type RegistryAllowlistsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *RegistryAllowlist
}

// Status returns the response status code.
func (r *RegistryAllowlistsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RegistryAllowlistsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RegistryAllowlistsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Data of the new registry allowlist.
func (r *RegistryAllowlistsAddResponse) Body() *RegistryAllowlist {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Data of the new registry allowlist.
func (r *RegistryAllowlistsAddResponse) GetBody() (value *RegistryAllowlist, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// RegistryAllowlistsListRequest is the request for the 'list' method.
type RegistryAllowlistsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	order     *string
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *RegistryAllowlistsListRequest) Parameter(name string, value interface{}) *RegistryAllowlistsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RegistryAllowlistsListRequest) Header(name string, value interface{}) *RegistryAllowlistsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RegistryAllowlistsListRequest) Impersonate(user string) *RegistryAllowlistsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement, but using the names of the attributes of the registry allowlists
// instead of the the names of the columns of a table. For example, in order to sort the
// allowlists descending by identifier the value should be:
//
// ```sql
// creation_timestamp desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *RegistryAllowlistsListRequest) Order(value string) *RegistryAllowlistsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *RegistryAllowlistsListRequest) Page(value int) *RegistryAllowlistsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of a
// SQL statement, but using the names of the attributes of the registry allowlists
// instead of the names of the columns of a table. For example, in order to retrieve all
// the allowlists with a specific cloud provider and creation time the following is required:
//
// ```sql
// cloud_provider.id='aws' and creation_timestamp > '2023-03-01T00:00:00Z'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// registry allowlists that the user has permission to see will be returned.
func (r *RegistryAllowlistsListRequest) Search(value string) *RegistryAllowlistsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *RegistryAllowlistsListRequest) Size(value int) *RegistryAllowlistsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RegistryAllowlistsListRequest) Send() (result *RegistryAllowlistsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RegistryAllowlistsListRequest) SendContext(ctx context.Context) (result *RegistryAllowlistsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.order != nil {
		helpers.AddValue(&query, "order", *r.order)
	}
	if r.page != nil {
		helpers.AddValue(&query, "page", *r.page)
	}
	if r.search != nil {
		helpers.AddValue(&query, "search", *r.search)
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
	result = &RegistryAllowlistsListResponse{}
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
	err = readRegistryAllowlistsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// RegistryAllowlistsListResponse is the response for the 'list' method.
type RegistryAllowlistsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *RegistryAllowlistList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *RegistryAllowlistsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RegistryAllowlistsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RegistryAllowlistsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of registry allowlists.
func (r *RegistryAllowlistsListResponse) Items() *RegistryAllowlistList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of registry allowlists.
func (r *RegistryAllowlistsListResponse) GetItems() (value *RegistryAllowlistList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *RegistryAllowlistsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *RegistryAllowlistsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *RegistryAllowlistsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *RegistryAllowlistsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *RegistryAllowlistsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *RegistryAllowlistsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
