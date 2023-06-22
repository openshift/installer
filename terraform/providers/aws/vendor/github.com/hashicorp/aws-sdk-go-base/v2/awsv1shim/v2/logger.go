// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsv1shim

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/hashicorp/aws-sdk-go-base/v2/logging"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv/v1.17.0/httpconv"
)

const (
	responseBufferLen = logging.MaxResponseBodyLen + 1024
)

type debugLogger struct{}

func (l debugLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	s := strings.Join(tokens, " ")
	s = strings.ReplaceAll(s, "\r", "") // Works around https://github.com/jen20/teamcity-go-test/pull/2
	log.Printf("missing_context: %s aws.sdk=aws-sdk-go", s)
}

func setAWSFields(ctx context.Context, r *request.Request) context.Context {
	ctx = tflog.SetField(ctx, "aws.sdk", "aws-sdk-go")
	ctx = tflog.SetField(ctx, "aws.service", r.ClientInfo.ServiceID)
	ctx = tflog.SetField(ctx, "aws.operation", r.Operation.Name)

	region := aws.StringValue(r.Config.Region)
	ctx = tflog.SetField(ctx, "aws.region", region)

	if signingRegion := r.ClientInfo.SigningRegion; signingRegion != region {
		ctx = tflog.SetField(ctx, "aws.signing_region", signingRegion)
	}

	return ctx
}

type durationKeyT string

const durationKey durationKeyT = "request-duration"

// Replaces the built-in logging middleware from https://github.com/aws/aws-sdk-go/blob/main/aws/client/logger.go
// We want access to the request struct, and cannot get it from the built-in.
// The typical route of adding logging to the http.RoundTripper doesn't work for the AWS SDK for Go v1 without forcing us to manually implement
// configuration that the SDK handles for us.
var requestLogger = request.NamedHandler{
	Name: "TF_AWS_RequestLogger",
	Fn:   logRequest,
}

func logRequest(r *request.Request) {
	ctx := r.Context()

	ctx = setAWSFields(ctx, r)

	bodySeekable := aws.IsReaderSeekable(r.Body)

	requestFields, err := logging.DecomposeHTTPRequest(r.HTTPRequest)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("decomposing request: %s", err))
		return
	}

	if !bodySeekable {
		r.SetReaderBody(aws.ReadSeekCloser(r.HTTPRequest.Body))
	}
	// Reset the request body because dumpRequest will re-wrap the
	// r.HTTPRequest's Body as a NoOpCloser and will not be reset after
	// read by the HTTP client reader.
	if err := r.Error; err != nil {
		tflog.Error(ctx, fmt.Sprintf("decomposing request: %s", err))
		return
	}

	tflog.Debug(ctx, "HTTP Request Sent", requestFields)

	ctx = context.WithValue(ctx, durationKey, time.Now())

	r.SetContext(ctx)
}

// Replaces the built-in logging middleware from https://github.com/aws/aws-sdk-go/blob/main/aws/client/logger.go
// We want access to the response struct, and cannot get it from the built-in.
// The typical route of adding logging to the http.RoundTripper doesn't work for the AWS SDK for Go v1 without forcing us to manually implement
// configuration that the SDK handles for us.
var responseLogger = request.NamedHandler{
	Name: "TF_AWS_ResponseLogger",
	Fn:   logResponse,
}

func logResponse(r *request.Request) {
	ctx := r.Context()

	ctx = setAWSFields(ctx, r)

	if r.HTTPResponse == nil {
		tflog.Error(ctx, "HTTP response is nil")
		return
	}

	bodyBuffer := bytes.NewBuffer(nil)

	r.HTTPResponse.Body = &teeReaderCloser{
		Reader: io.TeeReader(r.HTTPResponse.Body, limitWriter(bodyBuffer, responseBufferLen)),
		Source: r.HTTPResponse.Body,
	}

	handlerFn := func(req *request.Request) {
		ctx := r.Context()

		var elapsed time.Duration
		if start, ok := ctx.Value(durationKey).(time.Time); ok {
			elapsed = time.Since(start)
		}

		ctx = setAWSFields(ctx, r)

		responseFields, err := decomposeHTTPResponse(r.HTTPResponse, bodyBuffer, elapsed)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("decomposing response: %s", err))
			return
		}
		tflog.Debug(ctx, "HTTP Response Received", responseFields)
	}

	const handlerName = "TF_AWS_ResponseBodyLogger"

	r.Handlers.Unmarshal.SetBackNamed(request.NamedHandler{
		Name: handlerName, Fn: handlerFn,
	})
	r.Handlers.UnmarshalError.SetBackNamed(request.NamedHandler{
		Name: handlerName, Fn: handlerFn,
	})
}

type teeReaderCloser struct {
	// io.Reader will be a tee reader that is used during logging.
	// This structure will read from a body and write the contents to a logger.
	io.Reader
	// Source is used just to close when we are done reading.
	Source io.ReadCloser
}

func (reader *teeReaderCloser) Close() error {
	return reader.Source.Close()
}

func decomposeHTTPResponse(resp *http.Response, body io.Reader, elapsed time.Duration) (map[string]any, error) {
	var attributes []attribute.KeyValue

	attributes = append(attributes, attribute.Int64("http.duration", elapsed.Milliseconds()))

	attributes = append(attributes, httpconv.ClientResponse(resp)...)

	attributes = append(attributes, logging.DecomposeResponseHeaders(resp)...)

	bodyAttribute, err := decomposeResponseBody(body)
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

func decomposeResponseBody(bodyReader io.Reader) (kv attribute.KeyValue, err error) {
	content, err := io.ReadAll(bodyReader)
	if err != nil {
		return kv, err
	}

	reader := textproto.NewReader(bufio.NewReader(bytes.NewReader(content)))

	body, err := logging.ReadTruncatedBody(reader, logging.MaxResponseBodyLen)
	if err != nil {
		return kv, err
	}

	return attribute.String("http.response.body", body), nil
}

func limitWriter(w io.Writer, n int64) io.Writer {
	return &limitedWriter{w, n}
}

type limitedWriter struct {
	W io.Writer // the underlying writer
	N int64     // max bytes remaining
}

// Write writes data into the wrapped Writer up to a limit of N bytes
// Silently stops writing and returns full size of p to allow use with io.TeeReader
func (w *limitedWriter) Write(p []byte) (int, error) {
	if w.N <= 0 {
		return len(p), nil
	}
	if int64(len(p)) > w.N {
		n, err := w.W.Write(p[0:w.N])
		w.N -= int64(n)
		return len(p), err
	} else {
		n, err := w.W.Write(p)
		w.N -= int64(n)
		return n, err
	}
}
