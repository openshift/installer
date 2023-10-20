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

package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// AddValue creates the given set of query parameters if needed, an then adds
// the given parameter.
func AddValue(query *url.Values, name string, value interface{}) {
	if *query == nil {
		*query = make(url.Values)
	}
	query.Add(name, fmt.Sprintf("%v", value))
}

// CopyQuery creates a copy of the given set of query parameters.
func CopyQuery(query url.Values) url.Values {
	if query == nil {
		return nil
	}
	result := make(url.Values)
	for name, values := range query {
		result[name] = CopyValues(values)
	}
	return result
}

// AddHeader creates the given set of headers if needed, and then adds the given
// header:
func AddHeader(header *http.Header, name string, value interface{}) {
	if *header == nil {
		*header = make(http.Header)
	}
	header.Add(name, fmt.Sprintf("%v", value))
}

// CopyHeader creates a copy of the given set of headers.
func CopyHeader(header http.Header) http.Header {
	if header == nil {
		return nil
	}
	result := make(http.Header)
	for name, values := range header {
		result[name] = CopyValues(values)
	}
	return result
}

// CopyValues copies a slice of strings.
func CopyValues(values []string) []string {
	if values == nil {
		return nil
	}
	result := make([]string, len(values))
	copy(result, values)
	return result
}

// SetTimeout creates the given context, if needed, and sets the given timeout.
func SetTimeout(ctx *context.Context, cancel *context.CancelFunc, timeout time.Duration) {
	if *ctx == nil {
		*ctx = context.Background()
	}
	*ctx, *cancel = context.WithTimeout(*ctx, timeout)
}

// SetDeadline creates the given context, if needed, and sets the given deadline.
func SetDeadline(ctx *context.Context, cancel *context.CancelFunc, deadline time.Time) {
	if *ctx == nil {
		*ctx = context.Background()
	}
	*ctx, *cancel = context.WithDeadline(*ctx, deadline)
}
