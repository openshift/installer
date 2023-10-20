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

// NetworkVerificationsClient is the client of the 'network_verifications' resource.
//
// Manages the collection of subnet network verifications.
type NetworkVerificationsClient struct {
	transport http.RoundTripper
	path      string
}

// NewNetworkVerificationsClient creates a new client for the 'network_verifications'
// resource using the given transport to send the requests and receive the
// responses.
func NewNetworkVerificationsClient(transport http.RoundTripper, path string) *NetworkVerificationsClient {
	return &NetworkVerificationsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Creates an entry for a network verification for each subnet supplied setting then to initial state.
func (c *NetworkVerificationsClient) Add() *NetworkVerificationsAddRequest {
	return &NetworkVerificationsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// NetworkVerification returns the target 'network_verification' resource for the given identifier.
//
// Reference to the service that manages a specific network verification.
func (c *NetworkVerificationsClient) NetworkVerification(id string) *NetworkVerificationClient {
	return NewNetworkVerificationClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// NetworkVerificationsAddRequest is the request for the 'add' method.
type NetworkVerificationsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *NetworkVerification
}

// Parameter adds a query parameter.
func (r *NetworkVerificationsAddRequest) Parameter(name string, value interface{}) *NetworkVerificationsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NetworkVerificationsAddRequest) Header(name string, value interface{}) *NetworkVerificationsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NetworkVerificationsAddRequest) Impersonate(user string) *NetworkVerificationsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *NetworkVerificationsAddRequest) Body(value *NetworkVerification) *NetworkVerificationsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NetworkVerificationsAddRequest) Send() (result *NetworkVerificationsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NetworkVerificationsAddRequest) SendContext(ctx context.Context) (result *NetworkVerificationsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeNetworkVerificationsAddRequest(r, buffer)
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
	result = &NetworkVerificationsAddResponse{}
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
	err = readNetworkVerificationsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NetworkVerificationsAddResponse is the response for the 'add' method.
type NetworkVerificationsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *NetworkVerification
}

// Status returns the response status code.
func (r *NetworkVerificationsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NetworkVerificationsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NetworkVerificationsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *NetworkVerificationsAddResponse) Body() *NetworkVerification {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *NetworkVerificationsAddResponse) GetBody() (value *NetworkVerification, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
