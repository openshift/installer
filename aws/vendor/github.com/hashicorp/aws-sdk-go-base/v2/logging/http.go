// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/aws-sdk-go-base/v2/internal/slices"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv/v1.17.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	maxRequestBodyLen = 1024

	MaxResponseBodyLen = 4096
)

func DecomposeHTTPRequest(req *http.Request) (map[string]any, error) {
	var attributes []attribute.KeyValue

	attributes = append(attributes, httpconv.ClientRequest(req)...)
	// Remove empty `http.flavor`
	attributes = slices.Filter(attributes, func(attr attribute.KeyValue) bool {
		return attr.Key != semconv.HTTPFlavorKey || attr.Value.Emit() != ""
	})

	attributes = append(attributes, decomposeRequestHeaders(req)...)

	bodyAttribute, err := decomposeRequestBody(req)
	if err != nil {
		return nil, err
	}
	attributes = append(attributes, bodyAttribute)

	result := make(map[string]any, len(attributes))
	for _, attribute := range attributes {
		result[string(attribute.Key)] = attribute.Value.AsInterface()
	}

	return result, nil
}

func decomposeRequestHeaders(req *http.Request) []attribute.KeyValue {
	header := req.Header.Clone()

	// Handled directly from the Request
	header.Del("Content-Length")
	header.Del("User-Agent")

	results := make([]attribute.KeyValue, 0, len(header)+1)

	attempt := header.Values("Amz-Sdk-Request")
	if len(attempt) > 0 {
		if resendAttribute, ok := resendCountAttribute(attempt[0]); ok {
			results = append(results, resendAttribute)
		}
	}

	auth := header.Values("Authorization")
	if len(auth) > 0 {
		if authHeader, ok := authorizationHeaderAttribute(auth[0]); ok {
			results = append(results, authHeader)
		}
	}
	header.Del("Authorization")

	securityToken := header.Values("X-Amz-Security-Token")
	if len(securityToken) > 0 {
		results = append(results, requestHeaderAttribute("X-Amz-Security-Token").String("*****"))
	}
	header.Del("X-Amz-Security-Token")

	results = append(results, httpconv.RequestHeader(header)...)

	results = cleanUpHeaderAttributes(results)

	return results
}

func decomposeRequestBody(req *http.Request) (kv attribute.KeyValue, err error) {
	reqBytes, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return kv, err
	}

	reader := textproto.NewReader(bufio.NewReader(bytes.NewReader(reqBytes)))

	if _, err = reader.ReadLine(); err != nil {
		return kv, err
	}

	if _, err = reader.ReadMIMEHeader(); err != nil {
		return kv, err
	}

	body, err := ReadTruncatedBody(reader, maxRequestBodyLen)
	if err != nil {
		return kv, err
	}

	return attribute.String("http.request.body", body), nil
}

func requestHeaderAttribute(k string) attribute.Key {
	return attribute.Key(requestHeaderAttributeName(k))
}

func requestHeaderAttributeName(k string) string {
	return fmt.Sprintf("http.request.header.%s", normalizeHeaderName(k))
}

func normalizeHeaderName(k string) string {
	canonical := http.CanonicalHeaderKey(k)
	lower := strings.ToLower(canonical)
	return strings.ReplaceAll(lower, "-", "_")
}

func authorizationHeaderAttribute(v string) (attribute.KeyValue, bool) {
	parts := regexp.MustCompile(`\s+`).Split(v, 2) //nolint:gomnd
	if len(parts) != 2 {                           //nolint:gomnd
		return attribute.KeyValue{}, false
	}
	scheme := parts[0]
	if scheme == "" {
		return attribute.KeyValue{}, false
	}
	params := parts[1]
	if params == "" {
		return attribute.KeyValue{}, false
	}

	key := requestHeaderAttribute("Authorization")
	if strings.HasPrefix(scheme, "AWS4-") {
		components := regexp.MustCompile(`,\s+`).Split(params, -1)
		var builder strings.Builder
		builder.Grow(len(params))
		for i, component := range components {
			parts := strings.SplitAfterN(component, "=", 2)
			name := parts[0]
			value := parts[1]
			if name != "SignedHeaders=" && name != "Credential=" {
				// "Signature" or an unknown field
				value = "*****"
			}
			builder.WriteString(name)
			builder.WriteString(value)
			if i < len(components)-1 {
				builder.WriteString(", ")
			}
		}
		return key.String(fmt.Sprintf("%s %s", scheme, MaskAWSAccessKey(builder.String()))), true
	} else {
		return key.String(fmt.Sprintf("%s %s", scheme, strings.Repeat("*", len(params)))), true
	}
}

func resendCountAttribute(v string) (kv attribute.KeyValue, ok bool) {
	re := regexp.MustCompile(`attempt=(\d+);`)
	match := re.FindStringSubmatch(v)
	if len(match) != 2 { //nolint:gomnd
		return
	}

	attempt, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	if attempt > 1 {
		return attribute.Int("http.resend_count", attempt), true
	}

	return
}

func DecomposeResponseHeaders(resp *http.Response) []attribute.KeyValue {
	header := resp.Header.Clone()

	// Handled directly from the Response
	header.Del("Content-Length")

	results := make([]attribute.KeyValue, 0, len(header))

	results = append(results, httpconv.ResponseHeader(header)...)

	results = cleanUpHeaderAttributes(results)

	return results
}

// cleanUpHeaderAttributes converts header attributes with a single element to a string
func cleanUpHeaderAttributes(attrs []attribute.KeyValue) []attribute.KeyValue {
	return slices.ApplyToAll(attrs, func(attr attribute.KeyValue) attribute.KeyValue {
		if l := attr.Value.AsStringSlice(); len(l) == 1 {
			return attr.Key.String(l[0])
		}
		return attr
	})
}

func ReadTruncatedBody(reader *textproto.Reader, len int) (string, error) {
	var builder strings.Builder
	for {
		line, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		fmt.Fprintln(&builder, line)
		if builder.Len() >= len {
			fmt.Fprint(&builder, "[truncated...]")
			break
		}
	}

	body := builder.String()
	body = MaskAWSAccessKey(body)

	return body, nil
}
