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

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// NotifyClient is the client of the 'notify' resource.
//
// Manages the notify endpoint
type NotifyClient struct {
	transport http.RoundTripper
	path      string
}

// NewNotifyClient creates a new client for the 'notify'
// resource using the given transport to send the requests and receive the
// responses.
func NewNotifyClient(transport http.RoundTripper, path string) *NotifyClient {
	return &NotifyClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Notify user related to subscription/cluster via email
func (c *NotifyClient) Add() *NotifyAddRequest {
	return &NotifyAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// NotifyAddRequest is the request for the 'add' method.
type NotifyAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *SubscriptionNotify
}

// Parameter adds a query parameter.
func (r *NotifyAddRequest) Parameter(name string, value interface{}) *NotifyAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NotifyAddRequest) Header(name string, value interface{}) *NotifyAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NotifyAddRequest) Impersonate(user string) *NotifyAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *NotifyAddRequest) Body(value *SubscriptionNotify) *NotifyAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NotifyAddRequest) Send() (result *NotifyAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NotifyAddRequest) SendContext(ctx context.Context) (result *NotifyAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeNotifyAddRequest(r, buffer)
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
	result = &NotifyAddResponse{}
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
	err = readNotifyAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NotifyAddResponse is the response for the 'add' method.
type NotifyAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *SubscriptionNotify
}

// Status returns the response status code.
func (r *NotifyAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NotifyAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NotifyAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *NotifyAddResponse) Body() *SubscriptionNotify {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *NotifyAddResponse) GetBody() (value *SubscriptionNotify, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
