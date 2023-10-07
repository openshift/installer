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

// OidcConfigsClient is the client of the 'oidc_configs' resource.
//
// Manages the collection of oidc configurations.
type OidcConfigsClient struct {
	transport http.RoundTripper
	path      string
}

// NewOidcConfigsClient creates a new client for the 'oidc_configs'
// resource using the given transport to send the requests and receive the
// responses.
func NewOidcConfigsClient(transport http.RoundTripper, path string) *OidcConfigsClient {
	return &OidcConfigsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Creates a hosting under Red Hat's S3 bucket for byo oidc configuration.
func (c *OidcConfigsClient) Add() *OidcConfigsAddRequest {
	return &OidcConfigsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of oidc configs.
func (c *OidcConfigsClient) List() *OidcConfigsListRequest {
	return &OidcConfigsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// OidcConfig returns the target 'oidc_config' resource for the given identifier.
//
// Reference to the service that manages an specific identity provider.
func (c *OidcConfigsClient) OidcConfig(id string) *OidcConfigClient {
	return NewOidcConfigClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// OidcConfigsAddRequest is the request for the 'add' method.
type OidcConfigsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *OidcConfig
}

// Parameter adds a query parameter.
func (r *OidcConfigsAddRequest) Parameter(name string, value interface{}) *OidcConfigsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OidcConfigsAddRequest) Header(name string, value interface{}) *OidcConfigsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OidcConfigsAddRequest) Impersonate(user string) *OidcConfigsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *OidcConfigsAddRequest) Body(value *OidcConfig) *OidcConfigsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OidcConfigsAddRequest) Send() (result *OidcConfigsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OidcConfigsAddRequest) SendContext(ctx context.Context) (result *OidcConfigsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeOidcConfigsAddRequest(r, buffer)
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
	result = &OidcConfigsAddResponse{}
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
	err = readOidcConfigsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// OidcConfigsAddResponse is the response for the 'add' method.
type OidcConfigsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *OidcConfig
}

// Status returns the response status code.
func (r *OidcConfigsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OidcConfigsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OidcConfigsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *OidcConfigsAddResponse) Body() *OidcConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *OidcConfigsAddResponse) GetBody() (value *OidcConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// OidcConfigsListRequest is the request for the 'list' method.
type OidcConfigsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *OidcConfigsListRequest) Parameter(name string, value interface{}) *OidcConfigsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OidcConfigsListRequest) Header(name string, value interface{}) *OidcConfigsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OidcConfigsListRequest) Impersonate(user string) *OidcConfigsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OidcConfigsListRequest) Page(value int) *OidcConfigsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *OidcConfigsListRequest) Size(value int) *OidcConfigsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OidcConfigsListRequest) Send() (result *OidcConfigsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OidcConfigsListRequest) SendContext(ctx context.Context) (result *OidcConfigsListResponse, err error) {
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
	result = &OidcConfigsListResponse{}
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
	err = readOidcConfigsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// OidcConfigsListResponse is the response for the 'list' method.
type OidcConfigsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *OidcConfigList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *OidcConfigsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OidcConfigsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OidcConfigsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of identity providers.
func (r *OidcConfigsListResponse) Items() *OidcConfigList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of identity providers.
func (r *OidcConfigsListResponse) GetItems() (value *OidcConfigList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OidcConfigsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OidcConfigsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *OidcConfigsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *OidcConfigsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *OidcConfigsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *OidcConfigsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
