/*
Copyright 2020 The Kubernetes Authors.

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

package ot

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type (
	// OpenTelemetryAutorestTracer implements the tracing interface for AutoRest.
	OpenTelemetryAutorestTracer struct {
		tracer trace.Tracer
	}
)

var _ tracing.Tracer = (*OpenTelemetryAutorestTracer)(nil)

// NewOpenTelemetryAutorestTracer creates a new Autorest tracing adapter for OpenTelemetry.
func NewOpenTelemetryAutorestTracer(tracer trace.Tracer) *OpenTelemetryAutorestTracer {
	return &OpenTelemetryAutorestTracer{
		tracer: tracer,
	}
}

// NewTransport creates a new http.RoundTripper which will augment the base http.RoundTripper with OpenTelemetry
// tracing and metrics.
func (ot *OpenTelemetryAutorestTracer) NewTransport(base *http.Transport) http.RoundTripper {
	return otelhttp.NewTransport(base)
}

// StartSpan creates a new span of a given name and adds it to the given context.Context.
func (ot *OpenTelemetryAutorestTracer) StartSpan(ctx context.Context, name string) context.Context {
	ctx, _ = ot.tracer.Start(ctx, name)
	return ctx
}

// EndSpan ends the current context span. It ignores the httpsStatusCode and error since they are recorded by
// the otelhttp.Transport.
func (ot *OpenTelemetryAutorestTracer) EndSpan(ctx context.Context, httpStatusCode int, err error) {
	span := trace.SpanFromContext(ctx)
	span.End()
}
