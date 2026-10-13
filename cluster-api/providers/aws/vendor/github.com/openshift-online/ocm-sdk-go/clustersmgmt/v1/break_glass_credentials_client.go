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

// BreakGlassCredentialsClient is the client of the 'break_glass_credentials' resource.
//
// Manages the break glass credentials of a cluster.
type BreakGlassCredentialsClient struct {
	transport http.RoundTripper
	path      string
}

// NewBreakGlassCredentialsClient creates a new client for the 'break_glass_credentials'
// resource using the given transport to send the requests and receive the
// responses.
func NewBreakGlassCredentialsClient(transport http.RoundTripper, path string) *BreakGlassCredentialsClient {
	return &BreakGlassCredentialsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new break glass credential to the cluster.
func (c *BreakGlassCredentialsClient) Add() *BreakGlassCredentialsAddRequest {
	return &BreakGlassCredentialsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Revokes all the break glass certificates signed by a specific signer.
func (c *BreakGlassCredentialsClient) Delete() *BreakGlassCredentialsDeleteRequest {
	return &BreakGlassCredentialsDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of break glass credentials.
func (c *BreakGlassCredentialsClient) List() *BreakGlassCredentialsListRequest {
	return &BreakGlassCredentialsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// BreakGlassCredential returns the target 'break_glass_credential' resource for the given identifier.
//
// Reference to the service that manages a specific break glass credential.
func (c *BreakGlassCredentialsClient) BreakGlassCredential(id string) *BreakGlassCredentialClient {
	return NewBreakGlassCredentialClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// BreakGlassCredentialsAddRequest is the request for the 'add' method.
type BreakGlassCredentialsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *BreakGlassCredential
}

// Parameter adds a query parameter.
func (r *BreakGlassCredentialsAddRequest) Parameter(name string, value interface{}) *BreakGlassCredentialsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *BreakGlassCredentialsAddRequest) Header(name string, value interface{}) *BreakGlassCredentialsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *BreakGlassCredentialsAddRequest) Impersonate(user string) *BreakGlassCredentialsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the break glass credential.
func (r *BreakGlassCredentialsAddRequest) Body(value *BreakGlassCredential) *BreakGlassCredentialsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *BreakGlassCredentialsAddRequest) Send() (result *BreakGlassCredentialsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *BreakGlassCredentialsAddRequest) SendContext(ctx context.Context) (result *BreakGlassCredentialsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeBreakGlassCredentialsAddRequest(r, buffer)
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
	result = &BreakGlassCredentialsAddResponse{}
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
	err = readBreakGlassCredentialsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// BreakGlassCredentialsAddResponse is the response for the 'add' method.
type BreakGlassCredentialsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *BreakGlassCredential
}

// Status returns the response status code.
func (r *BreakGlassCredentialsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *BreakGlassCredentialsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *BreakGlassCredentialsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the break glass credential.
func (r *BreakGlassCredentialsAddResponse) Body() *BreakGlassCredential {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the break glass credential.
func (r *BreakGlassCredentialsAddResponse) GetBody() (value *BreakGlassCredential, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// BreakGlassCredentialsDeleteRequest is the request for the 'delete' method.
type BreakGlassCredentialsDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *BreakGlassCredentialsDeleteRequest) Parameter(name string, value interface{}) *BreakGlassCredentialsDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *BreakGlassCredentialsDeleteRequest) Header(name string, value interface{}) *BreakGlassCredentialsDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *BreakGlassCredentialsDeleteRequest) Impersonate(user string) *BreakGlassCredentialsDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *BreakGlassCredentialsDeleteRequest) Send() (result *BreakGlassCredentialsDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *BreakGlassCredentialsDeleteRequest) SendContext(ctx context.Context) (result *BreakGlassCredentialsDeleteResponse, err error) {
	query := helpers.CopyQuery(r.query)
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
	result = &BreakGlassCredentialsDeleteResponse{}
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

// BreakGlassCredentialsDeleteResponse is the response for the 'delete' method.
type BreakGlassCredentialsDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *BreakGlassCredentialsDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *BreakGlassCredentialsDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *BreakGlassCredentialsDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// BreakGlassCredentialsListRequest is the request for the 'list' method.
type BreakGlassCredentialsListRequest struct {
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
func (r *BreakGlassCredentialsListRequest) Parameter(name string, value interface{}) *BreakGlassCredentialsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *BreakGlassCredentialsListRequest) Header(name string, value interface{}) *BreakGlassCredentialsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *BreakGlassCredentialsListRequest) Impersonate(user string) *BreakGlassCredentialsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement, but using the names of the attributes of the break glass credentials
// instead of the the names of the columns of a table. For example, in order to sort the
// credentials descending by identifier the value should be:
//
// ```sql
// id desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *BreakGlassCredentialsListRequest) Order(value string) *BreakGlassCredentialsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *BreakGlassCredentialsListRequest) Page(value int) *BreakGlassCredentialsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of a
// SQL statement, but using the names of the attributes of the break glass credentials
// instead of the names of the columns of a table. For example, in order to retrieve all
// the credentials with a specific username and status the following is required:
//
// ```sql
// username='user1' AND status='expired'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// break glass credentials that the user has permission to see will be returned.
func (r *BreakGlassCredentialsListRequest) Search(value string) *BreakGlassCredentialsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *BreakGlassCredentialsListRequest) Size(value int) *BreakGlassCredentialsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *BreakGlassCredentialsListRequest) Send() (result *BreakGlassCredentialsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *BreakGlassCredentialsListRequest) SendContext(ctx context.Context) (result *BreakGlassCredentialsListResponse, err error) {
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
	result = &BreakGlassCredentialsListResponse{}
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
	err = readBreakGlassCredentialsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// BreakGlassCredentialsListResponse is the response for the 'list' method.
type BreakGlassCredentialsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *BreakGlassCredentialList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *BreakGlassCredentialsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *BreakGlassCredentialsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *BreakGlassCredentialsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of break glass credentials.
func (r *BreakGlassCredentialsListResponse) Items() *BreakGlassCredentialList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of break glass credentials.
func (r *BreakGlassCredentialsListResponse) GetItems() (value *BreakGlassCredentialList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *BreakGlassCredentialsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *BreakGlassCredentialsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *BreakGlassCredentialsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *BreakGlassCredentialsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *BreakGlassCredentialsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *BreakGlassCredentialsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
