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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

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

// AddonInstallationsClient is the client of the 'addon_installations' resource.
//
// Manages a collection of addon installations for a specific cluster
type AddonInstallationsClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddonInstallationsClient creates a new client for the 'addon_installations'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddonInstallationsClient(transport http.RoundTripper, path string) *AddonInstallationsClient {
	return &AddonInstallationsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Create a new addon status and add it to the collection of addons installation.
func (c *AddonInstallationsClient) Add() *AddonInstallationsAddRequest {
	return &AddonInstallationsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of addon installations for a cluster.
func (c *AddonInstallationsClient) List() *AddonInstallationsListRequest {
	return &AddonInstallationsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Addon returns the target 'addon_installation' resource for the given identifier.
//
// Returns a reference to the service that manages a specific addon installation.
func (c *AddonInstallationsClient) Addon(id string) *AddonInstallationClient {
	return NewAddonInstallationClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// AddonInstallationsAddRequest is the request for the 'add' method.
type AddonInstallationsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddonInstallation
}

// Parameter adds a query parameter.
func (r *AddonInstallationsAddRequest) Parameter(name string, value interface{}) *AddonInstallationsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonInstallationsAddRequest) Header(name string, value interface{}) *AddonInstallationsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonInstallationsAddRequest) Impersonate(user string) *AddonInstallationsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the addon installation.
func (r *AddonInstallationsAddRequest) Body(value *AddonInstallation) *AddonInstallationsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonInstallationsAddRequest) Send() (result *AddonInstallationsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonInstallationsAddRequest) SendContext(ctx context.Context) (result *AddonInstallationsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddonInstallationsAddRequest(r, buffer)
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
	result = &AddonInstallationsAddResponse{}
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
	err = readAddonInstallationsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonInstallationsAddResponse is the response for the 'add' method.
type AddonInstallationsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddonInstallation
}

// Status returns the response status code.
func (r *AddonInstallationsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonInstallationsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonInstallationsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the addon installation.
func (r *AddonInstallationsAddResponse) Body() *AddonInstallation {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the addon installation.
func (r *AddonInstallationsAddResponse) GetBody() (value *AddonInstallation, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddonInstallationsListRequest is the request for the 'list' method.
type AddonInstallationsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	order     *string
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *AddonInstallationsListRequest) Parameter(name string, value interface{}) *AddonInstallationsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonInstallationsListRequest) Header(name string, value interface{}) *AddonInstallationsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonInstallationsListRequest) Impersonate(user string) *AddonInstallationsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *AddonInstallationsListRequest) Order(value string) *AddonInstallationsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddonInstallationsListRequest) Page(value int) *AddonInstallationsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddonInstallationsListRequest) Size(value int) *AddonInstallationsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonInstallationsListRequest) Send() (result *AddonInstallationsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonInstallationsListRequest) SendContext(ctx context.Context) (result *AddonInstallationsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.order != nil {
		helpers.AddValue(&query, "order", *r.order)
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
	result = &AddonInstallationsListResponse{}
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
	err = readAddonInstallationsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonInstallationsListResponse is the response for the 'list' method.
type AddonInstallationsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *AddonInstallationList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *AddonInstallationsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonInstallationsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonInstallationsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of addon installations
func (r *AddonInstallationsListResponse) Items() *AddonInstallationList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of addon installations
func (r *AddonInstallationsListResponse) GetItems() (value *AddonInstallationList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddonInstallationsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddonInstallationsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddonInstallationsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddonInstallationsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection regardless of the size of the page.
func (r *AddonInstallationsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection regardless of the size of the page.
func (r *AddonInstallationsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
