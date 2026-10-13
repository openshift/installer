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

// PeerDependenciesClient is the client of the 'peer_dependencies' resource.
//
// Manages the collection of peer dependencies.
type PeerDependenciesClient struct {
	transport http.RoundTripper
	path      string
}

// NewPeerDependenciesClient creates a new client for the 'peer_dependencies'
// resource using the given transport to send the requests and receive the
// responses.
func NewPeerDependenciesClient(transport http.RoundTripper, path string) *PeerDependenciesClient {
	return &PeerDependenciesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
func (c *PeerDependenciesClient) Add() *PeerDependenciesAddRequest {
	return &PeerDependenciesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of peer dependencies.
func (c *PeerDependenciesClient) List() *PeerDependenciesListRequest {
	return &PeerDependenciesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// PeerDependency returns the target 'peer_dependency' resource for the given identifier.
func (c *PeerDependenciesClient) PeerDependency(id string) *PeerDependencyClient {
	return NewPeerDependencyClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// PeerDependenciesAddRequest is the request for the 'add' method.
type PeerDependenciesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *PeerDependency
}

// Parameter adds a query parameter.
func (r *PeerDependenciesAddRequest) Parameter(name string, value interface{}) *PeerDependenciesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PeerDependenciesAddRequest) Header(name string, value interface{}) *PeerDependenciesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PeerDependenciesAddRequest) Impersonate(user string) *PeerDependenciesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *PeerDependenciesAddRequest) Body(value *PeerDependency) *PeerDependenciesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PeerDependenciesAddRequest) Send() (result *PeerDependenciesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PeerDependenciesAddRequest) SendContext(ctx context.Context) (result *PeerDependenciesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writePeerDependenciesAddRequest(r, buffer)
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
	result = &PeerDependenciesAddResponse{}
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
	err = readPeerDependenciesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// PeerDependenciesAddResponse is the response for the 'add' method.
type PeerDependenciesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *PeerDependency
}

// Status returns the response status code.
func (r *PeerDependenciesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PeerDependenciesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PeerDependenciesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *PeerDependenciesAddResponse) Body() *PeerDependency {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *PeerDependenciesAddResponse) GetBody() (value *PeerDependency, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// PeerDependenciesListRequest is the request for the 'list' method.
type PeerDependenciesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	orderBy   *string
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *PeerDependenciesListRequest) Parameter(name string, value interface{}) *PeerDependenciesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PeerDependenciesListRequest) Header(name string, value interface{}) *PeerDependenciesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PeerDependenciesListRequest) Impersonate(user string) *PeerDependenciesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *PeerDependenciesListRequest) OrderBy(value string) *PeerDependenciesListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *PeerDependenciesListRequest) Page(value int) *PeerDependenciesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *PeerDependenciesListRequest) Size(value int) *PeerDependenciesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PeerDependenciesListRequest) Send() (result *PeerDependenciesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PeerDependenciesListRequest) SendContext(ctx context.Context) (result *PeerDependenciesListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.orderBy != nil {
		helpers.AddValue(&query, "order_by", *r.orderBy)
	}
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
	result = &PeerDependenciesListResponse{}
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
	err = readPeerDependenciesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// PeerDependenciesListResponse is the response for the 'list' method.
type PeerDependenciesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *PeerDependencyList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *PeerDependenciesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PeerDependenciesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PeerDependenciesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *PeerDependenciesListResponse) Items() *PeerDependencyList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *PeerDependenciesListResponse) GetItems() (value *PeerDependencyList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *PeerDependenciesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *PeerDependenciesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *PeerDependenciesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *PeerDependenciesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *PeerDependenciesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *PeerDependenciesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
