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

// SubscriptionNotifyClient is the client of the 'subscription_notify' resource.
//
// Manages the notify endpoint of a subscription resource
type SubscriptionNotifyClient struct {
	transport http.RoundTripper
	path      string
}

// NewSubscriptionNotifyClient creates a new client for the 'subscription_notify'
// resource using the given transport to send the requests and receive the
// responses.
func NewSubscriptionNotifyClient(transport http.RoundTripper, path string) *SubscriptionNotifyClient {
	return &SubscriptionNotifyClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Notify user related to subscription via email
func (c *SubscriptionNotifyClient) Add() *SubscriptionNotifyAddRequest {
	return &SubscriptionNotifyAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// SubscriptionNotifyAddRequest is the request for the 'add' method.
type SubscriptionNotifyAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *SubscriptionNotify
}

// Parameter adds a query parameter.
func (r *SubscriptionNotifyAddRequest) Parameter(name string, value interface{}) *SubscriptionNotifyAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *SubscriptionNotifyAddRequest) Header(name string, value interface{}) *SubscriptionNotifyAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *SubscriptionNotifyAddRequest) Impersonate(user string) *SubscriptionNotifyAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *SubscriptionNotifyAddRequest) Body(value *SubscriptionNotify) *SubscriptionNotifyAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *SubscriptionNotifyAddRequest) Send() (result *SubscriptionNotifyAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *SubscriptionNotifyAddRequest) SendContext(ctx context.Context) (result *SubscriptionNotifyAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeSubscriptionNotifyAddRequest(r, buffer)
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
	result = &SubscriptionNotifyAddResponse{}
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
	err = readSubscriptionNotifyAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// SubscriptionNotifyAddResponse is the response for the 'add' method.
type SubscriptionNotifyAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *SubscriptionNotify
}

// Status returns the response status code.
func (r *SubscriptionNotifyAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *SubscriptionNotifyAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *SubscriptionNotifyAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *SubscriptionNotifyAddResponse) Body() *SubscriptionNotify {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *SubscriptionNotifyAddResponse) GetBody() (value *SubscriptionNotify, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
