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

// DefaultCapabilitiesClient is the client of the 'default_capabilities' resource.
type DefaultCapabilitiesClient struct {
	transport http.RoundTripper
	path      string
}

// NewDefaultCapabilitiesClient creates a new client for the 'default_capabilities'
// resource using the given transport to send the requests and receive the
// responses.
func NewDefaultCapabilitiesClient(transport http.RoundTripper, path string) *DefaultCapabilitiesClient {
	return &DefaultCapabilitiesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Creates a new default capability.
func (c *DefaultCapabilitiesClient) Add() *DefaultCapabilitiesAddRequest {
	return &DefaultCapabilitiesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves a list of Dedfault Capabilities.
func (c *DefaultCapabilitiesClient) List() *DefaultCapabilitiesListRequest {
	return &DefaultCapabilitiesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// DefaultCapability returns the target 'default_capability' resource for the given identifier.
//
// Reference to the service that manages an specific default capability.
func (c *DefaultCapabilitiesClient) DefaultCapability(id string) *DefaultCapabilityClient {
	return NewDefaultCapabilityClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// DefaultCapabilitiesAddRequest is the request for the 'add' method.
type DefaultCapabilitiesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *DefaultCapability
}

// Parameter adds a query parameter.
func (r *DefaultCapabilitiesAddRequest) Parameter(name string, value interface{}) *DefaultCapabilitiesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DefaultCapabilitiesAddRequest) Header(name string, value interface{}) *DefaultCapabilitiesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DefaultCapabilitiesAddRequest) Impersonate(user string) *DefaultCapabilitiesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Default capability data.
func (r *DefaultCapabilitiesAddRequest) Body(value *DefaultCapability) *DefaultCapabilitiesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DefaultCapabilitiesAddRequest) Send() (result *DefaultCapabilitiesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DefaultCapabilitiesAddRequest) SendContext(ctx context.Context) (result *DefaultCapabilitiesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeDefaultCapabilitiesAddRequest(r, buffer)
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
	result = &DefaultCapabilitiesAddResponse{}
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
	err = readDefaultCapabilitiesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DefaultCapabilitiesAddResponse is the response for the 'add' method.
type DefaultCapabilitiesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *DefaultCapability
}

// Status returns the response status code.
func (r *DefaultCapabilitiesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DefaultCapabilitiesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DefaultCapabilitiesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Default capability data.
func (r *DefaultCapabilitiesAddResponse) Body() *DefaultCapability {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Default capability data.
func (r *DefaultCapabilitiesAddResponse) GetBody() (value *DefaultCapability, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// DefaultCapabilitiesListRequest is the request for the 'list' method.
type DefaultCapabilitiesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *DefaultCapabilitiesListRequest) Parameter(name string, value interface{}) *DefaultCapabilitiesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DefaultCapabilitiesListRequest) Header(name string, value interface{}) *DefaultCapabilitiesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DefaultCapabilitiesListRequest) Impersonate(user string) *DefaultCapabilitiesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *DefaultCapabilitiesListRequest) Page(value int) *DefaultCapabilitiesListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause
// of an SQL statement, but using the names of the attributes of the organization
// instead of the names of the columns of a table. For example, in order to
// retrieve organizations with name starting with my:
//
// ```sql
// name like 'my%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// items that the user has permission to see will be returned.
func (r *DefaultCapabilitiesListRequest) Search(value string) *DefaultCapabilitiesListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *DefaultCapabilitiesListRequest) Size(value int) *DefaultCapabilitiesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DefaultCapabilitiesListRequest) Send() (result *DefaultCapabilitiesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DefaultCapabilitiesListRequest) SendContext(ctx context.Context) (result *DefaultCapabilitiesListResponse, err error) {
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
	result = &DefaultCapabilitiesListResponse{}
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
	err = readDefaultCapabilitiesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DefaultCapabilitiesListResponse is the response for the 'list' method.
type DefaultCapabilitiesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *DefaultCapabilityList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *DefaultCapabilitiesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DefaultCapabilitiesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DefaultCapabilitiesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of Default Capabilities.
func (r *DefaultCapabilitiesListResponse) Items() *DefaultCapabilityList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of Default Capabilities.
func (r *DefaultCapabilitiesListResponse) GetItems() (value *DefaultCapabilityList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *DefaultCapabilitiesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *DefaultCapabilitiesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *DefaultCapabilitiesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *DefaultCapabilitiesListResponse) GetSize() (value int, ok bool) {
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
func (r *DefaultCapabilitiesListResponse) Total() int {
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
func (r *DefaultCapabilitiesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
