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

package tele

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracer struct {
	trace.Tracer
}

// Start creates a new context with a new Azure correlation ID, then
// creates a new trace.Span with that new context. This function then
// returns the new Context and Span.
func (t tracer) Start(
	ctx context.Context,
	op string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	ctx, corrID := ctxWithCorrID(ctx)
	opts = append(
		opts,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.String(
			string(CorrIDKeyVal),
			string(corrID),
		)),
	)
	return t.Tracer.Start(ctx, op, opts...)
}

// Tracer returns an OpenTelemetry Tracer implementation to be used
// to create spans. If you need access to the raw globally-registered
// tracer, use this function.
//
// Most people should not use this function directly, however.
// Instead, consider using StartSpanWithLogger, which uses
// this tracer to start a new span, configures logging, and
// more.
//
// Example usage:
//
//	ctx, span := tele.Tracer().Start(ctx, "myFunction")
//	defer span.End()
//	// use the span and context here
func Tracer() trace.Tracer {
	return tracer{
		Tracer: otel.Tracer("capz"),
	}
}
