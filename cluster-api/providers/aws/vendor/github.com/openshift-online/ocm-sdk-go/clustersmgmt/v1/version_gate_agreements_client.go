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

// VersionGateAgreementsClient is the client of the 'version_gate_agreements' resource.
//
// Manages the collection of version gates agreements for a cluster.
type VersionGateAgreementsClient struct {
	transport http.RoundTripper
	path      string
}

// NewVersionGateAgreementsClient creates a new client for the 'version_gate_agreements'
// resource using the given transport to send the requests and receive the
// responses.
func NewVersionGateAgreementsClient(transport http.RoundTripper, path string) *VersionGateAgreementsClient {
	return &VersionGateAgreementsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new agreed version gate to the cluster.
func (c *VersionGateAgreementsClient) Add() *VersionGateAgreementsAddRequest {
	return &VersionGateAgreementsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of reasons.
func (c *VersionGateAgreementsClient) List() *VersionGateAgreementsListRequest {
	return &VersionGateAgreementsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// VersionGateAgreement returns the target 'version_gate_agreement' resource for the given identifier.
//
// Reference to the service that manages a specific version gate agreement.
func (c *VersionGateAgreementsClient) VersionGateAgreement(id string) *VersionGateAgreementClient {
	return NewVersionGateAgreementClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// VersionGateAgreementsAddRequest is the request for the 'add' method.
type VersionGateAgreementsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *VersionGateAgreement
}

// Parameter adds a query parameter.
func (r *VersionGateAgreementsAddRequest) Parameter(name string, value interface{}) *VersionGateAgreementsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *VersionGateAgreementsAddRequest) Header(name string, value interface{}) *VersionGateAgreementsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *VersionGateAgreementsAddRequest) Impersonate(user string) *VersionGateAgreementsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Details of the version gate agreement.
func (r *VersionGateAgreementsAddRequest) Body(value *VersionGateAgreement) *VersionGateAgreementsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *VersionGateAgreementsAddRequest) Send() (result *VersionGateAgreementsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *VersionGateAgreementsAddRequest) SendContext(ctx context.Context) (result *VersionGateAgreementsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeVersionGateAgreementsAddRequest(r, buffer)
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
	result = &VersionGateAgreementsAddResponse{}
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
	err = readVersionGateAgreementsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// VersionGateAgreementsAddResponse is the response for the 'add' method.
type VersionGateAgreementsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *VersionGateAgreement
}

// Status returns the response status code.
func (r *VersionGateAgreementsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *VersionGateAgreementsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *VersionGateAgreementsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Details of the version gate agreement.
func (r *VersionGateAgreementsAddResponse) Body() *VersionGateAgreement {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Details of the version gate agreement.
func (r *VersionGateAgreementsAddResponse) GetBody() (value *VersionGateAgreement, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// VersionGateAgreementsListRequest is the request for the 'list' method.
type VersionGateAgreementsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *VersionGateAgreementsListRequest) Parameter(name string, value interface{}) *VersionGateAgreementsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *VersionGateAgreementsListRequest) Header(name string, value interface{}) *VersionGateAgreementsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *VersionGateAgreementsListRequest) Impersonate(user string) *VersionGateAgreementsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *VersionGateAgreementsListRequest) Page(value int) *VersionGateAgreementsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *VersionGateAgreementsListRequest) Size(value int) *VersionGateAgreementsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *VersionGateAgreementsListRequest) Send() (result *VersionGateAgreementsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *VersionGateAgreementsListRequest) SendContext(ctx context.Context) (result *VersionGateAgreementsListResponse, err error) {
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
	result = &VersionGateAgreementsListResponse{}
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
	err = readVersionGateAgreementsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// VersionGateAgreementsListResponse is the response for the 'list' method.
type VersionGateAgreementsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *VersionGateAgreementList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *VersionGateAgreementsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *VersionGateAgreementsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *VersionGateAgreementsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of version gate agreement.
func (r *VersionGateAgreementsListResponse) Items() *VersionGateAgreementList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of version gate agreement.
func (r *VersionGateAgreementsListResponse) GetItems() (value *VersionGateAgreementList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *VersionGateAgreementsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *VersionGateAgreementsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *VersionGateAgreementsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *VersionGateAgreementsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *VersionGateAgreementsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *VersionGateAgreementsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
