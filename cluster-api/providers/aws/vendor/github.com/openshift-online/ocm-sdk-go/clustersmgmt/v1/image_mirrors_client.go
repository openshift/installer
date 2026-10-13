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

// ImageMirrorsClient is the client of the 'image_mirrors' resource.
//
// Manages the collection of image mirror configurations for a cluster.
// Image mirroring is only supported for ROSA HCP clusters in ready state (Day 2 operations only).
type ImageMirrorsClient struct {
	transport http.RoundTripper
	path      string
}

// NewImageMirrorsClient creates a new client for the 'image_mirrors'
// resource using the given transport to send the requests and receive the
// responses.
func NewImageMirrorsClient(transport http.RoundTripper, path string) *ImageMirrorsClient {
	return &ImageMirrorsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Creates a new image mirror configuration for the cluster.
// Cluster must be in ready state for this operation to succeed.
func (c *ImageMirrorsClient) Add() *ImageMirrorsAddRequest {
	return &ImageMirrorsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of image mirrors for the cluster.
func (c *ImageMirrorsClient) List() *ImageMirrorsListRequest {
	return &ImageMirrorsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ImageMirror returns the target 'image_mirror' resource for the given identifier.
//
// Reference to the service that manages a specific image mirror.
func (c *ImageMirrorsClient) ImageMirror(id string) *ImageMirrorClient {
	return NewImageMirrorClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ImageMirrorsAddRequest is the request for the 'add' method.
type ImageMirrorsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ImageMirror
}

// Parameter adds a query parameter.
func (r *ImageMirrorsAddRequest) Parameter(name string, value interface{}) *ImageMirrorsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ImageMirrorsAddRequest) Header(name string, value interface{}) *ImageMirrorsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ImageMirrorsAddRequest) Impersonate(user string) *ImageMirrorsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the image mirror configuration.
func (r *ImageMirrorsAddRequest) Body(value *ImageMirror) *ImageMirrorsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ImageMirrorsAddRequest) Send() (result *ImageMirrorsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ImageMirrorsAddRequest) SendContext(ctx context.Context) (result *ImageMirrorsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeImageMirrorsAddRequest(r, buffer)
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
	result = &ImageMirrorsAddResponse{}
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
	err = readImageMirrorsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ImageMirrorsAddResponse is the response for the 'add' method.
type ImageMirrorsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ImageMirror
}

// Status returns the response status code.
func (r *ImageMirrorsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ImageMirrorsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ImageMirrorsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the image mirror configuration.
func (r *ImageMirrorsAddResponse) Body() *ImageMirror {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the image mirror configuration.
func (r *ImageMirrorsAddResponse) GetBody() (value *ImageMirror, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ImageMirrorsListRequest is the request for the 'list' method.
type ImageMirrorsListRequest struct {
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
func (r *ImageMirrorsListRequest) Parameter(name string, value interface{}) *ImageMirrorsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ImageMirrorsListRequest) Header(name string, value interface{}) *ImageMirrorsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ImageMirrorsListRequest) Impersonate(user string) *ImageMirrorsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria for sorting results.
func (r *ImageMirrorsListRequest) Order(value string) *ImageMirrorsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ImageMirrorsListRequest) Page(value int) *ImageMirrorsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria for filtering results.
// Searchable fields: id, name, cluster_id, source, type
// All searchable fields can be ordered with asc/desc direction, default order by id
func (r *ImageMirrorsListRequest) Search(value string) *ImageMirrorsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *ImageMirrorsListRequest) Size(value int) *ImageMirrorsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ImageMirrorsListRequest) Send() (result *ImageMirrorsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ImageMirrorsListRequest) SendContext(ctx context.Context) (result *ImageMirrorsListResponse, err error) {
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
	result = &ImageMirrorsListResponse{}
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
	err = readImageMirrorsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ImageMirrorsListResponse is the response for the 'list' method.
type ImageMirrorsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ImageMirrorList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ImageMirrorsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ImageMirrorsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ImageMirrorsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of image mirrors.
func (r *ImageMirrorsListResponse) Items() *ImageMirrorList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of image mirrors.
func (r *ImageMirrorsListResponse) GetItems() (value *ImageMirrorList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ImageMirrorsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ImageMirrorsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *ImageMirrorsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *ImageMirrorsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *ImageMirrorsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *ImageMirrorsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
