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

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ExternalAuthsClient is the client of the 'external_auths' resource.
//
// Manages the collection of external authentication providers defined for an ARO HCP cluster.
type ExternalAuthsClient struct {
	transport http.RoundTripper
	path      string
}

// NewExternalAuthsClient creates a new client for the 'external_auths'
// resource using the given transport to send the requests and receive the
// responses.
func NewExternalAuthsClient(transport http.RoundTripper, path string) *ExternalAuthsClient {
	return &ExternalAuthsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'async_add' method.
//
// Adds a new external authentication provider to the cluster.
func (c *ExternalAuthsClient) Add() *ExternalAuthsAddRequest {
	return &ExternalAuthsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
func (c *ExternalAuthsClient) List() *ExternalAuthsListRequest {
	return &ExternalAuthsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ExternalAuth returns the target 'external_auth' resource for the given identifier.
//
// Reference to the resource that manages a specific external authentication provider.
func (c *ExternalAuthsClient) ExternalAuth(id string) *ExternalAuthClient {
	return NewExternalAuthClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ExternalAuthsAddRequest is the request for the 'async_add' method.
type ExternalAuthsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ExternalAuth
}

// Parameter adds a query parameter.
func (r *ExternalAuthsAddRequest) Parameter(name string, value interface{}) *ExternalAuthsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ExternalAuthsAddRequest) Header(name string, value interface{}) *ExternalAuthsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ExternalAuthsAddRequest) Impersonate(user string) *ExternalAuthsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the external authentication provider.
func (r *ExternalAuthsAddRequest) Body(value *ExternalAuth) *ExternalAuthsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ExternalAuthsAddRequest) Send() (result *ExternalAuthsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ExternalAuthsAddRequest) SendContext(ctx context.Context) (result *ExternalAuthsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeExternalAuthsAsyncAddRequest(r, buffer)
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
	result = &ExternalAuthsAddResponse{}
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
	err = readExternalAuthsAsyncAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ExternalAuthsAddResponse is the response for the 'async_add' method.
type ExternalAuthsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ExternalAuth
}

// Status returns the response status code.
func (r *ExternalAuthsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ExternalAuthsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ExternalAuthsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the external authentication provider.
func (r *ExternalAuthsAddResponse) Body() *ExternalAuth {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the external authentication provider.
func (r *ExternalAuthsAddResponse) GetBody() (value *ExternalAuth, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ExternalAuthsListRequest is the request for the 'list' method.
type ExternalAuthsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *ExternalAuthsListRequest) Parameter(name string, value interface{}) *ExternalAuthsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ExternalAuthsListRequest) Header(name string, value interface{}) *ExternalAuthsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ExternalAuthsListRequest) Impersonate(user string) *ExternalAuthsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ExternalAuthsListRequest) Page(value int) *ExternalAuthsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *ExternalAuthsListRequest) Size(value int) *ExternalAuthsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ExternalAuthsListRequest) Send() (result *ExternalAuthsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ExternalAuthsListRequest) SendContext(ctx context.Context) (result *ExternalAuthsListResponse, err error) {
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
	result = &ExternalAuthsListResponse{}
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
	err = readExternalAuthsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ExternalAuthsListResponse is the response for the 'list' method.
type ExternalAuthsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ExternalAuthList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ExternalAuthsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ExternalAuthsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ExternalAuthsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of external authentication providers.
func (r *ExternalAuthsListResponse) Items() *ExternalAuthList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of external authentication providers.
func (r *ExternalAuthsListResponse) GetItems() (value *ExternalAuthList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ExternalAuthsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ExternalAuthsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *ExternalAuthsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *ExternalAuthsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *ExternalAuthsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *ExternalAuthsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
