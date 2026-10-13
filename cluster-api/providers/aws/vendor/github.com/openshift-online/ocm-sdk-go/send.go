/*
Copyright (c) 2018 Red Hat, Inc.

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

// This file contains the implementation of the methods of the connection that are used to send HTTP
// requests and receive HTTP responses.

package sdk

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/openshift-online/ocm-sdk-go/internal"
)

// RoundTrip is the implementation of the http.RoundTripper interface.
func (c *Connection) RoundTrip(request *http.Request) (response *http.Response, err error) {
	// Check if the connection is closed:
	err = c.checkClosed()
	if err != nil {
		return
	}

	// Get the context from the request:
	ctx := request.Context()

	// Check the request URL:
	if request.URL.Path == "" {
		err = fmt.Errorf("request path is mandatory")
		return
	}
	if request.URL.Scheme != "" || request.URL.Host != "" || !path.IsAbs(request.URL.Path) {
		err = fmt.Errorf("request URL '%s' isn't absolute", request.URL)
		return
	}

	// Select the target server add the base URL to the request URL:
	server, err := c.selectServer(ctx, request)
	if err != nil {
		return
	}
	request.URL = server.URL.ResolveReference(request.URL)

	// Check the request method and body:
	switch request.Method {
	case http.MethodGet, http.MethodDelete:
		if request.Body != nil {
			c.logger.Warn(ctx,
				"Request body is not allowed for the '%s' method",
				request.Method,
			)
		}
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		// POST and PATCH and PUT don't need to have a body. It is up to the server to decide if
		// this is acceptable.
	default:
		err = fmt.Errorf("method '%s' is not allowed", request.Method)
		return
	}

	// Add the default headers:
	if request.Header == nil {
		request.Header = make(http.Header)
	}
	if c.agent != "" {
		request.Header.Set("User-Agent", c.agent)
	}
	switch request.Method {
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		request.Header.Set("Content-Type", "application/json")
	}
	request.Header.Set("Accept", "application/json")

	// Select the client:
	client, err := c.clientSelector.Select(ctx, server)
	if err != nil {
		return
	}

	// Send the request:
	response, err = client.Do(request)
	if err != nil {
		return
	}

	// Check that the response content type is JSON:
	err = internal.CheckContentType(response)
	if err != nil {
		return
	}

	return
}

// selectServer selects the server that should be used for the given request, according its path and
// the alternative URLs configured when the connection was created.
func (c *Connection) selectServer(ctx context.Context,
	request *http.Request) (base *internal.ServerAddress, err error) {
	// Select the server corresponding to the longest matching prefix. Note that it is enough to
	// pick the first match because the entries have already been sorted by descending prefix
	// length when the connection was created.
	for _, entry := range c.urlTable {
		if entry.re.MatchString(request.URL.Path) {
			base = entry.url
			return
		}
	}
	if base == nil {
		err = fmt.Errorf(
			"can't find any matching URL for request path '%s'",
			request.URL.Path,
		)
	}
	return
}
