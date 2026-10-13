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
	"path"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// SupportCasesClient is the client of the 'support_cases' resource.
//
// Manages the support cases endpoint
type SupportCasesClient struct {
	transport http.RoundTripper
	path      string
}

// NewSupportCasesClient creates a new client for the 'support_cases'
// resource using the given transport to send the requests and receive the
// responses.
func NewSupportCasesClient(transport http.RoundTripper, path string) *SupportCasesClient {
	return &SupportCasesClient{
		transport: transport,
		path:      path,
	}
}

// Post creates a request for the 'post' method.
//
// Create a support case related to Hydra
func (c *SupportCasesClient) Post() *SupportCasesPostRequest {
	return &SupportCasesPostRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// SupportCase returns the target 'support_case' resource for the given identifier.
//
// Reference to the service that manages a specific support case.
func (c *SupportCasesClient) SupportCase(id string) *SupportCaseClient {
	return NewSupportCaseClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// SupportCasesPostRequest is the request for the 'post' method.
type SupportCasesPostRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	request   *SupportCaseRequest
}

// Parameter adds a query parameter.
func (r *SupportCasesPostRequest) Parameter(name string, value interface{}) *SupportCasesPostRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *SupportCasesPostRequest) Header(name string, value interface{}) *SupportCasesPostRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *SupportCasesPostRequest) Impersonate(user string) *SupportCasesPostRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Request sets the value of the 'request' parameter.
func (r *SupportCasesPostRequest) Request(value *SupportCaseRequest) *SupportCasesPostRequest {
	r.request = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *SupportCasesPostRequest) Send() (result *SupportCasesPostResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *SupportCasesPostRequest) SendContext(ctx context.Context) (result *SupportCasesPostResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeSupportCasesPostRequest(r, buffer)
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
	result = &SupportCasesPostResponse{}
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
	err = readSupportCasesPostResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// SupportCasesPostResponse is the response for the 'post' method.
type SupportCasesPostResponse struct {
	status   int
	header   http.Header
	err      *errors.Error
	response *SupportCaseResponse
}

// Status returns the response status code.
func (r *SupportCasesPostResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *SupportCasesPostResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *SupportCasesPostResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Response returns the value of the 'response' parameter.
func (r *SupportCasesPostResponse) Response() *SupportCaseResponse {
	if r == nil {
		return nil
	}
	return r.response
}

// GetResponse returns the value of the 'response' parameter and
// a flag indicating if the parameter has a value.
func (r *SupportCasesPostResponse) GetResponse() (value *SupportCaseResponse, ok bool) {
	ok = r != nil && r.response != nil
	if ok {
		value = r.response
	}
	return
}
