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

// UsersClient is the client of the 'users' resource.
//
// Manages the collection of users of a group.
type UsersClient struct {
	transport http.RoundTripper
	path      string
}

// NewUsersClient creates a new client for the 'users'
// resource using the given transport to send the requests and receive the
// responses.
func NewUsersClient(transport http.RoundTripper, path string) *UsersClient {
	return &UsersClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new user to the group.
func (c *UsersClient) Add() *UsersAddRequest {
	return &UsersAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of users.
func (c *UsersClient) List() *UsersListRequest {
	return &UsersListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// User returns the target 'user' resource for the given identifier.
//
// Reference to the service that manages an specific user.
func (c *UsersClient) User(id string) *UserClient {
	return NewUserClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// UsersAddRequest is the request for the 'add' method.
type UsersAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *User
}

// Parameter adds a query parameter.
func (r *UsersAddRequest) Parameter(name string, value interface{}) *UsersAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *UsersAddRequest) Header(name string, value interface{}) *UsersAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *UsersAddRequest) Impersonate(user string) *UsersAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the user.
func (r *UsersAddRequest) Body(value *User) *UsersAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *UsersAddRequest) Send() (result *UsersAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *UsersAddRequest) SendContext(ctx context.Context) (result *UsersAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeUsersAddRequest(r, buffer)
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
	result = &UsersAddResponse{}
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
	err = readUsersAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// UsersAddResponse is the response for the 'add' method.
type UsersAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *User
}

// Status returns the response status code.
func (r *UsersAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *UsersAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *UsersAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the user.
func (r *UsersAddResponse) Body() *User {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the user.
func (r *UsersAddResponse) GetBody() (value *User, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// UsersListRequest is the request for the 'list' method.
type UsersListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *UsersListRequest) Parameter(name string, value interface{}) *UsersListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *UsersListRequest) Header(name string, value interface{}) *UsersListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *UsersListRequest) Impersonate(user string) *UsersListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *UsersListRequest) Page(value int) *UsersListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *UsersListRequest) Size(value int) *UsersListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *UsersListRequest) Send() (result *UsersListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *UsersListRequest) SendContext(ctx context.Context) (result *UsersListResponse, err error) {
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
	result = &UsersListResponse{}
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
	err = readUsersListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// UsersListResponse is the response for the 'list' method.
type UsersListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *UserList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *UsersListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *UsersListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *UsersListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of users.
func (r *UsersListResponse) Items() *UserList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of users.
func (r *UsersListResponse) GetItems() (value *UserList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *UsersListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *UsersListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *UsersListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *UsersListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *UsersListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *UsersListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
