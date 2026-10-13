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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

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

// IncidentsClient is the client of the 'incidents' resource.
//
// Manages the collection of incidents.
type IncidentsClient struct {
	transport http.RoundTripper
	path      string
}

// NewIncidentsClient creates a new client for the 'incidents'
// resource using the given transport to send the requests and receive the
// responses.
func NewIncidentsClient(transport http.RoundTripper, path string) *IncidentsClient {
	return &IncidentsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
func (c *IncidentsClient) Add() *IncidentsAddRequest {
	return &IncidentsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of incidents.
func (c *IncidentsClient) List() *IncidentsListRequest {
	return &IncidentsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Incident returns the target 'incident' resource for the given identifier.
func (c *IncidentsClient) Incident(id string) *IncidentClient {
	return NewIncidentClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// IncidentsAddRequest is the request for the 'add' method.
type IncidentsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Incident
}

// Parameter adds a query parameter.
func (r *IncidentsAddRequest) Parameter(name string, value interface{}) *IncidentsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *IncidentsAddRequest) Header(name string, value interface{}) *IncidentsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *IncidentsAddRequest) Impersonate(user string) *IncidentsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *IncidentsAddRequest) Body(value *Incident) *IncidentsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *IncidentsAddRequest) Send() (result *IncidentsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *IncidentsAddRequest) SendContext(ctx context.Context) (result *IncidentsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeIncidentsAddRequest(r, buffer)
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
	result = &IncidentsAddResponse{}
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
	err = readIncidentsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// IncidentsAddResponse is the response for the 'add' method.
type IncidentsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Incident
}

// Status returns the response status code.
func (r *IncidentsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *IncidentsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *IncidentsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *IncidentsAddResponse) Body() *Incident {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *IncidentsAddResponse) GetBody() (value *Incident, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// IncidentsListRequest is the request for the 'list' method.
type IncidentsListRequest struct {
	transport            http.RoundTripper
	path                 string
	query                url.Values
	header               http.Header
	creatorId            *string
	incidentCommanderId  *string
	incidentName         *string
	mine                 *bool
	onCallId             *string
	orderBy              *string
	page                 *int
	participantId        *string
	productId            *string
	publicId             *string
	responsibleManagerId *string
	size                 *int
	status_              *string
}

// Parameter adds a query parameter.
func (r *IncidentsListRequest) Parameter(name string, value interface{}) *IncidentsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *IncidentsListRequest) Header(name string, value interface{}) *IncidentsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *IncidentsListRequest) Impersonate(user string) *IncidentsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// CreatorId sets the value of the 'creator_id' parameter.
func (r *IncidentsListRequest) CreatorId(value string) *IncidentsListRequest {
	r.creatorId = &value
	return r
}

// IncidentCommanderId sets the value of the 'incident_commander_id' parameter.
func (r *IncidentsListRequest) IncidentCommanderId(value string) *IncidentsListRequest {
	r.incidentCommanderId = &value
	return r
}

// IncidentName sets the value of the 'incident_name' parameter.
func (r *IncidentsListRequest) IncidentName(value string) *IncidentsListRequest {
	r.incidentName = &value
	return r
}

// Mine sets the value of the 'mine' parameter.
func (r *IncidentsListRequest) Mine(value bool) *IncidentsListRequest {
	r.mine = &value
	return r
}

// OnCallId sets the value of the 'on_call_id' parameter.
func (r *IncidentsListRequest) OnCallId(value string) *IncidentsListRequest {
	r.onCallId = &value
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *IncidentsListRequest) OrderBy(value string) *IncidentsListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *IncidentsListRequest) Page(value int) *IncidentsListRequest {
	r.page = &value
	return r
}

// ParticipantId sets the value of the 'participant_id' parameter.
func (r *IncidentsListRequest) ParticipantId(value string) *IncidentsListRequest {
	r.participantId = &value
	return r
}

// ProductId sets the value of the 'product_id' parameter.
func (r *IncidentsListRequest) ProductId(value string) *IncidentsListRequest {
	r.productId = &value
	return r
}

// PublicId sets the value of the 'public_id' parameter.
func (r *IncidentsListRequest) PublicId(value string) *IncidentsListRequest {
	r.publicId = &value
	return r
}

// ResponsibleManagerId sets the value of the 'responsible_manager_id' parameter.
func (r *IncidentsListRequest) ResponsibleManagerId(value string) *IncidentsListRequest {
	r.responsibleManagerId = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *IncidentsListRequest) Size(value int) *IncidentsListRequest {
	r.size = &value
	return r
}

// Status sets the value of the 'status' parameter.
func (r *IncidentsListRequest) Status(value string) *IncidentsListRequest {
	r.status_ = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *IncidentsListRequest) Send() (result *IncidentsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *IncidentsListRequest) SendContext(ctx context.Context) (result *IncidentsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.creatorId != nil {
		helpers.AddValue(&query, "creator_id", *r.creatorId)
	}
	if r.incidentCommanderId != nil {
		helpers.AddValue(&query, "incident_commander_id", *r.incidentCommanderId)
	}
	if r.incidentName != nil {
		helpers.AddValue(&query, "incident_name", *r.incidentName)
	}
	if r.mine != nil {
		helpers.AddValue(&query, "mine", *r.mine)
	}
	if r.onCallId != nil {
		helpers.AddValue(&query, "on_call_id", *r.onCallId)
	}
	if r.orderBy != nil {
		helpers.AddValue(&query, "order_by", *r.orderBy)
	}
	if r.page != nil {
		helpers.AddValue(&query, "page", *r.page)
	}
	if r.participantId != nil {
		helpers.AddValue(&query, "participant_id", *r.participantId)
	}
	if r.productId != nil {
		helpers.AddValue(&query, "product_id", *r.productId)
	}
	if r.publicId != nil {
		helpers.AddValue(&query, "public_id", *r.publicId)
	}
	if r.responsibleManagerId != nil {
		helpers.AddValue(&query, "responsible_manager_id", *r.responsibleManagerId)
	}
	if r.size != nil {
		helpers.AddValue(&query, "size", *r.size)
	}
	if r.status_ != nil {
		helpers.AddValue(&query, "status", *r.status_)
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
	result = &IncidentsListResponse{}
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
	err = readIncidentsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// IncidentsListResponse is the response for the 'list' method.
type IncidentsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *IncidentList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *IncidentsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *IncidentsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *IncidentsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *IncidentsListResponse) Items() *IncidentList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *IncidentsListResponse) GetItems() (value *IncidentList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *IncidentsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *IncidentsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *IncidentsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *IncidentsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *IncidentsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *IncidentsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
