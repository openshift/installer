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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

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

// ServicesClient is the client of the 'services' resource.
//
// Manages the collection of services.
type ServicesClient struct {
	transport http.RoundTripper
	path      string
}

// NewServicesClient creates a new client for the 'services'
// resource using the given transport to send the requests and receive the
// responses.
func NewServicesClient(transport http.RoundTripper, path string) *ServicesClient {
	return &ServicesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
func (c *ServicesClient) Add() *ServicesAddRequest {
	return &ServicesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of services.
func (c *ServicesClient) List() *ServicesListRequest {
	return &ServicesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Service returns the target 'service' resource for the given identifier.
func (c *ServicesClient) Service(id string) *ServiceClient {
	return NewServiceClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ServicesAddRequest is the request for the 'add' method.
type ServicesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Service
}

// Parameter adds a query parameter.
func (r *ServicesAddRequest) Parameter(name string, value interface{}) *ServicesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServicesAddRequest) Header(name string, value interface{}) *ServicesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServicesAddRequest) Impersonate(user string) *ServicesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ServicesAddRequest) Body(value *Service) *ServicesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServicesAddRequest) Send() (result *ServicesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServicesAddRequest) SendContext(ctx context.Context) (result *ServicesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeServicesAddRequest(r, buffer)
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
	result = &ServicesAddResponse{}
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
	err = readServicesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ServicesAddResponse is the response for the 'add' method.
type ServicesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Service
}

// Status returns the response status code.
func (r *ServicesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServicesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServicesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ServicesAddResponse) Body() *Service {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ServicesAddResponse) GetBody() (value *Service, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ServicesListRequest is the request for the 'list' method.
type ServicesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	fullname  *string
	mine      *bool
	orderBy   *string
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *ServicesListRequest) Parameter(name string, value interface{}) *ServicesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServicesListRequest) Header(name string, value interface{}) *ServicesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServicesListRequest) Impersonate(user string) *ServicesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Fullname sets the value of the 'fullname' parameter.
func (r *ServicesListRequest) Fullname(value string) *ServicesListRequest {
	r.fullname = &value
	return r
}

// Mine sets the value of the 'mine' parameter.
func (r *ServicesListRequest) Mine(value bool) *ServicesListRequest {
	r.mine = &value
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *ServicesListRequest) OrderBy(value string) *ServicesListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *ServicesListRequest) Page(value int) *ServicesListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
func (r *ServicesListRequest) Search(value string) *ServicesListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *ServicesListRequest) Size(value int) *ServicesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServicesListRequest) Send() (result *ServicesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServicesListRequest) SendContext(ctx context.Context) (result *ServicesListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.fullname != nil {
		helpers.AddValue(&query, "fullname", *r.fullname)
	}
	if r.mine != nil {
		helpers.AddValue(&query, "mine", *r.mine)
	}
	if r.orderBy != nil {
		helpers.AddValue(&query, "order_by", *r.orderBy)
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
	result = &ServicesListResponse{}
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
	err = readServicesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ServicesListResponse is the response for the 'list' method.
type ServicesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ServiceList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ServicesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServicesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServicesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *ServicesListResponse) Items() *ServiceList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *ServicesListResponse) GetItems() (value *ServiceList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *ServicesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *ServicesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *ServicesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *ServicesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *ServicesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *ServicesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
