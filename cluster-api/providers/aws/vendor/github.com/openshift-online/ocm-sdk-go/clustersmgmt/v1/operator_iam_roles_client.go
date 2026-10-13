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

// OperatorIAMRolesClient is the client of the 'operator_IAM_roles' resource.
//
// Manages the collection of operator roles.
type OperatorIAMRolesClient struct {
	transport http.RoundTripper
	path      string
}

// NewOperatorIAMRolesClient creates a new client for the 'operator_IAM_roles'
// resource using the given transport to send the requests and receive the
// responses.
func NewOperatorIAMRolesClient(transport http.RoundTripper, path string) *OperatorIAMRolesClient {
	return &OperatorIAMRolesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new operator role to the cluster.
func (c *OperatorIAMRolesClient) Add() *OperatorIAMRolesAddRequest {
	return &OperatorIAMRolesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of operator roles.
func (c *OperatorIAMRolesClient) List() *OperatorIAMRolesListRequest {
	return &OperatorIAMRolesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// OperatorIAMRole returns the target 'operator_IAM_role' resource for the given identifier.
//
// Returns a reference to the service that manages a specific operator role.
func (c *OperatorIAMRolesClient) OperatorIAMRole(id string) *OperatorIAMRoleClient {
	return NewOperatorIAMRoleClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// OperatorIAMRolesAddRequest is the request for the 'add' method.
type OperatorIAMRolesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *OperatorIAMRole
}

// Parameter adds a query parameter.
func (r *OperatorIAMRolesAddRequest) Parameter(name string, value interface{}) *OperatorIAMRolesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OperatorIAMRolesAddRequest) Header(name string, value interface{}) *OperatorIAMRolesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OperatorIAMRolesAddRequest) Impersonate(user string) *OperatorIAMRolesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the operator role.
func (r *OperatorIAMRolesAddRequest) Body(value *OperatorIAMRole) *OperatorIAMRolesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OperatorIAMRolesAddRequest) Send() (result *OperatorIAMRolesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OperatorIAMRolesAddRequest) SendContext(ctx context.Context) (result *OperatorIAMRolesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeOperatorIAMRolesAddRequest(r, buffer)
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
	result = &OperatorIAMRolesAddResponse{}
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
	err = readOperatorIAMRolesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// OperatorIAMRolesAddResponse is the response for the 'add' method.
type OperatorIAMRolesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *OperatorIAMRole
}

// Status returns the response status code.
func (r *OperatorIAMRolesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OperatorIAMRolesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OperatorIAMRolesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the operator role.
func (r *OperatorIAMRolesAddResponse) Body() *OperatorIAMRole {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the operator role.
func (r *OperatorIAMRolesAddResponse) GetBody() (value *OperatorIAMRole, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// OperatorIAMRolesListRequest is the request for the 'list' method.
type OperatorIAMRolesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *OperatorIAMRolesListRequest) Parameter(name string, value interface{}) *OperatorIAMRolesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OperatorIAMRolesListRequest) Header(name string, value interface{}) *OperatorIAMRolesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OperatorIAMRolesListRequest) Impersonate(user string) *OperatorIAMRolesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OperatorIAMRolesListRequest) Page(value int) *OperatorIAMRolesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items that will be contained in the returned page.
func (r *OperatorIAMRolesListRequest) Size(value int) *OperatorIAMRolesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OperatorIAMRolesListRequest) Send() (result *OperatorIAMRolesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OperatorIAMRolesListRequest) SendContext(ctx context.Context) (result *OperatorIAMRolesListResponse, err error) {
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
	result = &OperatorIAMRolesListResponse{}
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
	err = readOperatorIAMRolesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// OperatorIAMRolesListResponse is the response for the 'list' method.
type OperatorIAMRolesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *OperatorIAMRoleList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *OperatorIAMRolesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OperatorIAMRolesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OperatorIAMRolesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of operator roles.
func (r *OperatorIAMRolesListResponse) Items() *OperatorIAMRoleList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of operator roles.
func (r *OperatorIAMRolesListResponse) GetItems() (value *OperatorIAMRoleList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OperatorIAMRolesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OperatorIAMRolesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items that will be contained in the returned page.
func (r *OperatorIAMRolesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items that will be contained in the returned page.
func (r *OperatorIAMRolesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection that match the search criteria,
// regardless of the size of the page.
func (r *OperatorIAMRolesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection that match the search criteria,
// regardless of the size of the page.
func (r *OperatorIAMRolesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
