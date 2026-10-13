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

// This file contains an implementation of the http.RoundTripper interface that sends to the log
// the details of the requests sent and the responses received.

package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	ordered "gitlab.com/c0b/go-ordered-json"
)

// RoundTripperBuilder contains the information an logic needed to build a new round tripper that
// sends to the log the details of the requests sent and the responses received. Don't create
// instances of this type directly; use the NewRoundTripper function instead.
type RoundTripperBuilder struct {
	logger *logrus.Logger
	redact map[string]bool
	next   http.RoundTripper
}

// RoundTripper is a round tripper that dumps the details of the requests and the responses to
// the log. Don't create instances of this type directly; use the NewRoundTripper function instead.
type RoundTripper struct {
	logger *logrus.Logger
	redact map[string]bool
	next   http.RoundTripper
}

// Make sure that we implement the http.RoundTripper interface:
var _ http.RoundTripper = &RoundTripper{}

// NewRoundTripper creates a builder that can then be used to create a round tripper that sends to
// the log the details of the requests sent and the responses received.
func NewRoundTripper() *RoundTripperBuilder {
	return &RoundTripperBuilder{}
}

// Logger sets the logger that the round tripper will use to send the details of request and
// responses to the log. This is mandatory.
func (b *RoundTripperBuilder) Logger(value *logrus.Logger) *RoundTripperBuilder {
	b.logger = value
	return b
}

// Redact specifies a field whose value should be removed from the messages sent to the log.
func (b *RoundTripperBuilder) Redact(value string) *RoundTripperBuilder {
	if b.redact == nil {
		b.redact = make(map[string]bool)
	}
	b.redact[value] = true
	return b
}

// Next sets the next round tripper. The details of the request will be sent to the log before
// calling it, and the details of the response will be sent to the log after calling it.
func (b *RoundTripperBuilder) Next(value http.RoundTripper) *RoundTripperBuilder {
	b.next = value
	return b
}

// Build uses the information stored in the builder to create a new round tripper that sends to the
// log the details of the requests sent and the responses received.
func (b *RoundTripperBuilder) Build() (result *RoundTripper, err error) {
	// Check parameters:
	if b.logger == nil {
		err = fmt.Errorf("Logger is mandatory")
		return
	}
	if b.next == nil {
		err = fmt.Errorf("Next handler is mandatory")
		return
	}

	// Copy the set of redactedReplacement fields:
	redact := make(map[string]bool)
	for key, value := range b.redact {
		redact[key] = value
	}

	// Create and populate the object:
	result = &RoundTripper{
		logger: b.logger,
		redact: redact,
		next:   b.next,
	}

	return
}

// RoundTrip is he implementation of the http.RoundTripper interface.
func (d *RoundTripper) RoundTrip(request *http.Request) (response *http.Response, err error) {
	// Read the complete body in memory, in order to send it to the log, and replace it with a
	// reader that reads it from memory:
	if request.Body != nil {
		var body []byte
		body, err = io.ReadAll(request.Body)
		if err != nil {
			return
		}
		err = request.Body.Close()
		if err != nil {
			return
		}
		d.dumpRequest(request, body)
		request.Body = io.NopCloser(bytes.NewBuffer(body))
	} else {
		d.dumpRequest(request, nil)
	}

	// Call the next round tripper:
	response, err = d.next.RoundTrip(request)
	if err != nil {
		return
	}

	// Read the complete response body in memory, in order to send it the log, and replace it
	// with a reader that reads it from memory:
	if response.Body != nil {
		var body []byte
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return
		}
		err = response.Body.Close()
		if err != nil {
			return
		}
		d.dumpResponse(response, body)
		response.Body = io.NopCloser(bytes.NewBuffer(body))
	} else {
		d.dumpResponse(response, nil)
	}

	return
}

// dumpRequest dumps to the log, in debug level, the details of the given HTTP request.
func (d *RoundTripper) dumpRequest(request *http.Request, body []byte) {
	d.logger.Debugf("Request method is %s", request.Method)
	d.logger.Debugf("Request URL is '%s'", request.URL)
	header := request.Header
	names := make([]string, len(header))
	i := 0
	for name := range header {
		names[i] = name
		i++
	}
	sort.Strings(names)
	for _, name := range names {
		values := header[name]
		for _, value := range values {
			if strings.ToLower(name) == "authorization" {
				d.logger.Debugf("Request header '%s' is omitted", name)
			} else {
				d.logger.Debugf("Request header '%s' is '%s'", name, value)
			}
		}
	}
	if body != nil {
		d.dumpBody("Request", header, body)
	}
}

// dumpResponse dumps to the log, in debug level, the details of the given HTTP response.
func (d *RoundTripper) dumpResponse(response *http.Response, body []byte) {
	d.logger.Debugf("Response status is '%s'", response.Status)
	header := response.Header
	names := make([]string, len(header))
	i := 0
	for name := range header {
		names[i] = name
		i++
	}
	sort.Strings(names)
	for _, name := range names {
		values := header[name]
		for _, value := range values {
			d.logger.Debugf("Response header '%s' is '%s'", name, value)
		}
	}
	if body != nil {
		d.dumpBody("Response", header, body)
	}
}

// dumpBody checks the content type used in the given header and then it dumps the given body in a
// format suitable for that content type.
func (d *RoundTripper) dumpBody(what string, header http.Header, body []byte) {
	// Try to parse the content type:
	var mediaType string
	contentType := header.Get("Content-Type")
	if contentType != "" {
		var err error
		mediaType, _, err = mime.ParseMediaType(contentType)
		if err != nil {
			d.logger.Errorf("Failed to parse content type '%s': %v", contentType, err)
		}
	} else {
		mediaType = contentType
	}

	// Dump the body according to the content type:
	switch mediaType {
	case "application/x-www-form-urlencoded":
		d.dumpForm(what, body)
	case "application/json", "application/x-amz-json-1.0", "application/x-amz-json-1.1":
		d.dumpJSON(what, body)
	default:
		d.dumpBytes(what, body)
	}
}

// dumpForm sends to the log the contents of the given form data, excluding security sensitive
// fields.
func (d *RoundTripper) dumpForm(what string, data []byte) {
	// Parse the form:
	form, err := url.ParseQuery(string(data))
	if err != nil {
		d.dumpBytes(what, data)
		return
	}

	// Redact values corresponding to security sensitive fields:
	for name, values := range form {
		if d.redact[name] {
			for i := range values {
				values[i] = redactedReplacement
			}
		}
	}

	// Get and sort the names of the fields of the form, so that the generated output will be
	// predictable:
	names := make([]string, 0, len(form))
	for name := range form {
		names = append(names, name)
	}
	sort.Strings(names)

	// Write the names and values to the buffer while redacting the sensitive fields:
	buffer := &bytes.Buffer{}
	for _, name := range names {
		key := url.QueryEscape(name)
		values := form[name]
		for _, value := range values {
			var redacted string
			if d.redact[name] {
				redacted = redactedReplacement
				d.logger.Debugf("%s field '%s' is redacted", what, name)
			} else {
				redacted = url.QueryEscape(value)
				d.logger.Debugf("%s field '%s' is '%s'", what, name, value)
			}
			if buffer.Len() > 0 {
				buffer.WriteByte('&') // #nosec G104
			}
			buffer.WriteString(key)      // #nosec G104
			buffer.WriteByte('=')        // #nosec G104
			buffer.WriteString(redacted) // #nosec G104
		}
	}

	// Send the redactedReplacement data to the log:
	d.dumpBytes(what, buffer.Bytes())
}

// dumpJSON tries to parse the given data as a JSON document. If that works, then it dumps it
// indented, otherwise dumps it as is.
func (d *RoundTripper) dumpJSON(what string, data []byte) {
	parsed := ordered.NewOrderedMap()
	err := json.Unmarshal(data, parsed)
	if err != nil {
		d.logger.Debugf("%s", data)
	} else {
		// remove sensitive information
		d.redactSensitive(parsed)

		indented, err := json.MarshalIndent(parsed, "", "  ")
		if err != nil {
			d.dumpBytes(what, data)
		} else {
			d.dumpBytes(what, indented)
		}
	}
}

// dumpBytes dump the given data as an array of bytes.
func (d *RoundTripper) dumpBytes(what string, data []byte) {
	size := len(data)
	if size > 0 {
		d.logger.Debugf("%s body follows", what)
		d.logger.Out.Write(data)
		last := data[size-1]
		if last != '\n' {
			d.logger.Out.Write([]byte("\n"))
		}
	}
}

// redactSensitive replaces sensitive fields within a response with redactionStr.
func (d *RoundTripper) redactSensitive(body *ordered.OrderedMap) {
	iterator := body.EntriesIter()
	for {
		pair, ok := iterator()
		if !ok {
			break
		}
		if d.redact[pair.Key] {
			body.Set(pair.Key, redactedReplacement)
		}
	}
}

// String that replaces redactedReplacement fields in messages sent to the log:
const redactedReplacement = "***"
