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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

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

// ResourceQuotasClient is the client of the 'resource_quotas' resource.
//
// Manages the collection of resource quotas for an organization.
type ResourceQuotasClient struct {
	transport http.RoundTripper
	path      string
}

// NewResourceQuotasClient creates a new client for the 'resource_quotas'
// resource using the given transport to send the requests and receive the
// responses.
func NewResourceQuotasClient(transport http.RoundTripper, path string) *ResourceQuotasClient {
	return &ResourceQuotasClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Creates a new resource quota.
func (c *ResourceQuotasClient) Add() *ResourceQuotasAddRequest {
	return &ResourceQuotasAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of resource quotas.
func (c *ResourceQuotasClient) List() *ResourceQuotasListRequest {
	return &ResourceQuotasListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ResourceQuota returns the target 'resource_quota' resource for the given identifier.
//
// Reference to the service that manages an specific resource quota.
func (c *ResourceQuotasClient) ResourceQuota(id string) *ResourceQuotaClient {
	return NewResourceQuotaClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ResourceQuotasAddRequest is the request for the 'add' method.
type ResourceQuotasAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ResourceQuota
}

// Parameter adds a query parameter.
func (r *ResourceQuotasAddRequest) Parameter(name string, value interface{}) *ResourceQuotasAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ResourceQuotasAddRequest) Header(name string, value interface{}) *ResourceQuotasAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ResourceQuotasAddRequest) Impersonate(user string) *ResourceQuotasAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Resource quota data.
func (r *ResourceQuotasAddRequest) Body(value *ResourceQuota) *ResourceQuotasAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ResourceQuotasAddRequest) Send() (result *ResourceQuotasAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ResourceQuotasAddRequest) SendContext(ctx context.Context) (result *ResourceQuotasAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeResourceQuotasAddRequest(r, buffer)
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
	result = &ResourceQuotasAddResponse{}
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
	err = readResourceQuotasAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ResourceQuotasAddResponse is the response for the 'add' method.
type ResourceQuotasAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ResourceQuota
}

// Status returns the response status code.
func (r *ResourceQuotasAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ResourceQuotasAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ResourceQuotasAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Resource quota data.
func (r *ResourceQuotasAddResponse) Body() *ResourceQuota {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Resource quota data.
func (r *ResourceQuotasAddResponse) GetBody() (value *ResourceQuota, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ResourceQuotasListRequest is the request for the 'list' method.
type ResourceQuotasListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *ResourceQuotasListRequest) Parameter(name string, value interface{}) *ResourceQuotasListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ResourceQuotasListRequest) Header(name string, value interface{}) *ResourceQuotasListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ResourceQuotasListRequest) Impersonate(user string) *ResourceQuotasListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ResourceQuotasListRequest) Page(value int) *ResourceQuotasListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause
// of an SQL statement, but using the names of the attributes of the account
// instead of the names of the columns of a table. For example, in order to
// retrieve resource quota with resource_type cluster.aws:
//
// ```sql
// resource_type = 'cluster.aws'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// items that the user has permission to see will be returned.
func (r *ResourceQuotasListRequest) Search(value string) *ResourceQuotasListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *ResourceQuotasListRequest) Size(value int) *ResourceQuotasListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ResourceQuotasListRequest) Send() (result *ResourceQuotasListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ResourceQuotasListRequest) SendContext(ctx context.Context) (result *ResourceQuotasListResponse, err error) {
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
	result = &ResourceQuotasListResponse{}
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
	err = readResourceQuotasListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ResourceQuotasListResponse is the response for the 'list' method.
type ResourceQuotasListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ResourceQuotaList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ResourceQuotasListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ResourceQuotasListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ResourceQuotasListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of resource quotas.
func (r *ResourceQuotasListResponse) Items() *ResourceQuotaList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of resource quotas.
func (r *ResourceQuotasListResponse) GetItems() (value *ResourceQuotaList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ResourceQuotasListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ResourceQuotasListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *ResourceQuotasListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *ResourceQuotasListResponse) GetSize() (value int, ok bool) {
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
func (r *ResourceQuotasListResponse) Total() int {
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
func (r *ResourceQuotasListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
