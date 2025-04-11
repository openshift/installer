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

// KubeletConfigsClient is the client of the 'kubelet_configs' resource.
//
// Manages the collection of KubeletConfigs for a cluster.
type KubeletConfigsClient struct {
	transport http.RoundTripper
	path      string
}

// NewKubeletConfigsClient creates a new client for the 'kubelet_configs'
// resource using the given transport to send the requests and receive the
// responses.
func NewKubeletConfigsClient(transport http.RoundTripper, path string) *KubeletConfigsClient {
	return &KubeletConfigsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new KubeletConfig to the cluster.
func (c *KubeletConfigsClient) Add() *KubeletConfigsAddRequest {
	return &KubeletConfigsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of KubeletConfigs for the cluster.
func (c *KubeletConfigsClient) List() *KubeletConfigsListRequest {
	return &KubeletConfigsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// KubeletConfig returns the target 'hcp_kubelet_config' resource for the given identifier.
//
// Reference to the service that manages a specific KubeletConfig.
func (c *KubeletConfigsClient) KubeletConfig(id string) *HcpKubeletConfigClient {
	return NewHcpKubeletConfigClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// KubeletConfigsAddRequest is the request for the 'add' method.
type KubeletConfigsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *KubeletConfig
}

// Parameter adds a query parameter.
func (r *KubeletConfigsAddRequest) Parameter(name string, value interface{}) *KubeletConfigsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *KubeletConfigsAddRequest) Header(name string, value interface{}) *KubeletConfigsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *KubeletConfigsAddRequest) Impersonate(user string) *KubeletConfigsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the KubeletConfig.
func (r *KubeletConfigsAddRequest) Body(value *KubeletConfig) *KubeletConfigsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *KubeletConfigsAddRequest) Send() (result *KubeletConfigsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *KubeletConfigsAddRequest) SendContext(ctx context.Context) (result *KubeletConfigsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeKubeletConfigsAddRequest(r, buffer)
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
	result = &KubeletConfigsAddResponse{}
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
	err = readKubeletConfigsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// KubeletConfigsAddResponse is the response for the 'add' method.
type KubeletConfigsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *KubeletConfig
}

// Status returns the response status code.
func (r *KubeletConfigsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *KubeletConfigsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *KubeletConfigsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the KubeletConfig.
func (r *KubeletConfigsAddResponse) Body() *KubeletConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the KubeletConfig.
func (r *KubeletConfigsAddResponse) GetBody() (value *KubeletConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// KubeletConfigsListRequest is the request for the 'list' method.
type KubeletConfigsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *KubeletConfigsListRequest) Parameter(name string, value interface{}) *KubeletConfigsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *KubeletConfigsListRequest) Header(name string, value interface{}) *KubeletConfigsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *KubeletConfigsListRequest) Impersonate(user string) *KubeletConfigsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *KubeletConfigsListRequest) Page(value int) *KubeletConfigsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *KubeletConfigsListRequest) Size(value int) *KubeletConfigsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *KubeletConfigsListRequest) Send() (result *KubeletConfigsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *KubeletConfigsListRequest) SendContext(ctx context.Context) (result *KubeletConfigsListResponse, err error) {
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
	result = &KubeletConfigsListResponse{}
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
	err = readKubeletConfigsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// KubeletConfigsListResponse is the response for the 'list' method.
type KubeletConfigsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *KubeletConfigList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *KubeletConfigsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *KubeletConfigsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *KubeletConfigsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of KubeletConfigs.
func (r *KubeletConfigsListResponse) Items() *KubeletConfigList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of KubeletConfigs.
func (r *KubeletConfigsListResponse) GetItems() (value *KubeletConfigList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *KubeletConfigsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *KubeletConfigsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *KubeletConfigsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *KubeletConfigsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *KubeletConfigsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *KubeletConfigsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
