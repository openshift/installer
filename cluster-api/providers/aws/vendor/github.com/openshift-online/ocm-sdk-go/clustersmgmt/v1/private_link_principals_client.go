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

// PrivateLinkPrincipalsClient is the client of the 'private_link_principals' resource.
//
// Contains a list of principals for the Private Link.
type PrivateLinkPrincipalsClient struct {
	transport http.RoundTripper
	path      string
}

// NewPrivateLinkPrincipalsClient creates a new client for the 'private_link_principals'
// resource using the given transport to send the requests and receive the
// responses.
func NewPrivateLinkPrincipalsClient(transport http.RoundTripper, path string) *PrivateLinkPrincipalsClient {
	return &PrivateLinkPrincipalsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new principal for the Private Link.
func (c *PrivateLinkPrincipalsClient) Add() *PrivateLinkPrincipalsAddRequest {
	return &PrivateLinkPrincipalsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of principals.
func (c *PrivateLinkPrincipalsClient) List() *PrivateLinkPrincipalsListRequest {
	return &PrivateLinkPrincipalsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Principal returns the target 'private_link_principal' resource for the given identifier.
func (c *PrivateLinkPrincipalsClient) Principal(id string) *PrivateLinkPrincipalClient {
	return NewPrivateLinkPrincipalClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// PrivateLinkPrincipalsAddRequest is the request for the 'add' method.
type PrivateLinkPrincipalsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *PrivateLinkPrincipal
}

// Parameter adds a query parameter.
func (r *PrivateLinkPrincipalsAddRequest) Parameter(name string, value interface{}) *PrivateLinkPrincipalsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PrivateLinkPrincipalsAddRequest) Header(name string, value interface{}) *PrivateLinkPrincipalsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PrivateLinkPrincipalsAddRequest) Impersonate(user string) *PrivateLinkPrincipalsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Properties of the principal.
func (r *PrivateLinkPrincipalsAddRequest) Body(value *PrivateLinkPrincipal) *PrivateLinkPrincipalsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PrivateLinkPrincipalsAddRequest) Send() (result *PrivateLinkPrincipalsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PrivateLinkPrincipalsAddRequest) SendContext(ctx context.Context) (result *PrivateLinkPrincipalsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writePrivateLinkPrincipalsAddRequest(r, buffer)
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
	result = &PrivateLinkPrincipalsAddResponse{}
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
	err = readPrivateLinkPrincipalsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// PrivateLinkPrincipalsAddResponse is the response for the 'add' method.
type PrivateLinkPrincipalsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *PrivateLinkPrincipal
}

// Status returns the response status code.
func (r *PrivateLinkPrincipalsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PrivateLinkPrincipalsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PrivateLinkPrincipalsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Properties of the principal.
func (r *PrivateLinkPrincipalsAddResponse) Body() *PrivateLinkPrincipal {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Properties of the principal.
func (r *PrivateLinkPrincipalsAddResponse) GetBody() (value *PrivateLinkPrincipal, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// PrivateLinkPrincipalsListRequest is the request for the 'list' method.
type PrivateLinkPrincipalsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *PrivateLinkPrincipalsListRequest) Parameter(name string, value interface{}) *PrivateLinkPrincipalsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PrivateLinkPrincipalsListRequest) Header(name string, value interface{}) *PrivateLinkPrincipalsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PrivateLinkPrincipalsListRequest) Impersonate(user string) *PrivateLinkPrincipalsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *PrivateLinkPrincipalsListRequest) Page(value int) *PrivateLinkPrincipalsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause
// of an SQL statement, but using the names of the attributes of the role binding
// instead of the names of the columns of a table. For example, in order to
// retrieve role bindings with role_id AuthenticatedUser:
//
// ```sql
// role_id = 'AuthenticatedUser'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// items that the user has permission to see will be returned.
func (r *PrivateLinkPrincipalsListRequest) Search(value string) *PrivateLinkPrincipalsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *PrivateLinkPrincipalsListRequest) Size(value int) *PrivateLinkPrincipalsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PrivateLinkPrincipalsListRequest) Send() (result *PrivateLinkPrincipalsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PrivateLinkPrincipalsListRequest) SendContext(ctx context.Context) (result *PrivateLinkPrincipalsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
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
	result = &PrivateLinkPrincipalsListResponse{}
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
	err = readPrivateLinkPrincipalsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// PrivateLinkPrincipalsListResponse is the response for the 'list' method.
type PrivateLinkPrincipalsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *PrivateLinkPrincipalList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *PrivateLinkPrincipalsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PrivateLinkPrincipalsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PrivateLinkPrincipalsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of principals.
func (r *PrivateLinkPrincipalsListResponse) Items() *PrivateLinkPrincipalList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of principals.
func (r *PrivateLinkPrincipalsListResponse) GetItems() (value *PrivateLinkPrincipalList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *PrivateLinkPrincipalsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *PrivateLinkPrincipalsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *PrivateLinkPrincipalsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *PrivateLinkPrincipalsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *PrivateLinkPrincipalsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *PrivateLinkPrincipalsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
