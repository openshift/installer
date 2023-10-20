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

// LimitedSupportReasonsClient is the client of the 'limited_support_reasons' resource.
//
// Manages the collection of limited support reason on a cluster.
type LimitedSupportReasonsClient struct {
	transport http.RoundTripper
	path      string
}

// NewLimitedSupportReasonsClient creates a new client for the 'limited_support_reasons'
// resource using the given transport to send the requests and receive the
// responses.
func NewLimitedSupportReasonsClient(transport http.RoundTripper, path string) *LimitedSupportReasonsClient {
	return &LimitedSupportReasonsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new reason to the cluster.
func (c *LimitedSupportReasonsClient) Add() *LimitedSupportReasonsAddRequest {
	return &LimitedSupportReasonsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of reasons.
func (c *LimitedSupportReasonsClient) List() *LimitedSupportReasonsListRequest {
	return &LimitedSupportReasonsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// LimitedSupportReason returns the target 'limited_support_reason' resource for the given identifier.
//
// Reference to the service that manages an specific reason.
func (c *LimitedSupportReasonsClient) LimitedSupportReason(id string) *LimitedSupportReasonClient {
	return NewLimitedSupportReasonClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// LimitedSupportReasonsAddRequest is the request for the 'add' method.
type LimitedSupportReasonsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *LimitedSupportReason
}

// Parameter adds a query parameter.
func (r *LimitedSupportReasonsAddRequest) Parameter(name string, value interface{}) *LimitedSupportReasonsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *LimitedSupportReasonsAddRequest) Header(name string, value interface{}) *LimitedSupportReasonsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *LimitedSupportReasonsAddRequest) Impersonate(user string) *LimitedSupportReasonsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the reason.
func (r *LimitedSupportReasonsAddRequest) Body(value *LimitedSupportReason) *LimitedSupportReasonsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *LimitedSupportReasonsAddRequest) Send() (result *LimitedSupportReasonsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *LimitedSupportReasonsAddRequest) SendContext(ctx context.Context) (result *LimitedSupportReasonsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeLimitedSupportReasonsAddRequest(r, buffer)
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
	result = &LimitedSupportReasonsAddResponse{}
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
	err = readLimitedSupportReasonsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// LimitedSupportReasonsAddResponse is the response for the 'add' method.
type LimitedSupportReasonsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *LimitedSupportReason
}

// Status returns the response status code.
func (r *LimitedSupportReasonsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *LimitedSupportReasonsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *LimitedSupportReasonsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the reason.
func (r *LimitedSupportReasonsAddResponse) Body() *LimitedSupportReason {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the reason.
func (r *LimitedSupportReasonsAddResponse) GetBody() (value *LimitedSupportReason, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// LimitedSupportReasonsListRequest is the request for the 'list' method.
type LimitedSupportReasonsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *LimitedSupportReasonsListRequest) Parameter(name string, value interface{}) *LimitedSupportReasonsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *LimitedSupportReasonsListRequest) Header(name string, value interface{}) *LimitedSupportReasonsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *LimitedSupportReasonsListRequest) Impersonate(user string) *LimitedSupportReasonsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *LimitedSupportReasonsListRequest) Page(value int) *LimitedSupportReasonsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *LimitedSupportReasonsListRequest) Size(value int) *LimitedSupportReasonsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *LimitedSupportReasonsListRequest) Send() (result *LimitedSupportReasonsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *LimitedSupportReasonsListRequest) SendContext(ctx context.Context) (result *LimitedSupportReasonsListResponse, err error) {
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
	result = &LimitedSupportReasonsListResponse{}
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
	err = readLimitedSupportReasonsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// LimitedSupportReasonsListResponse is the response for the 'list' method.
type LimitedSupportReasonsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *LimitedSupportReasonList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *LimitedSupportReasonsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *LimitedSupportReasonsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *LimitedSupportReasonsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of template.
func (r *LimitedSupportReasonsListResponse) Items() *LimitedSupportReasonList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of template.
func (r *LimitedSupportReasonsListResponse) GetItems() (value *LimitedSupportReasonList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *LimitedSupportReasonsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *LimitedSupportReasonsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *LimitedSupportReasonsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *LimitedSupportReasonsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *LimitedSupportReasonsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *LimitedSupportReasonsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
