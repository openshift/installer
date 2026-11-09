/*
Copyright 2026 The Kubernetes Authors.

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

// Package otel provides OpenTelemetry tracing initialization and helper.
package otel

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"sigs.k8s.io/cluster-api-provider-aws/v2/version"
)

// InitTracer initializes the OpenTelemetry tracer provider.
func InitTracer(ctx context.Context, traceSamplingRatio float64, tracingEndPoint string) (*sdktrace.TracerProvider, error) {
	if tracingEndPoint == "" {
		return nil, nil
	}
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(tracingEndPoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		slog.Error("failed to create OTLP trace exporter", "error", err)
		return nil, err
	}

	// Creating the resource for sdktrace
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("cluster-api-provider-aws"),
		semconv.ServiceVersionKey.String(version.Get().String()),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(
			sdktrace.ParentBased(
				sdktrace.TraceIDRatioBased(traceSamplingRatio),
			),
		),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

// StartSpan starts a new OpenTelemetry span for the given context.
func StartSpan(ctx context.Context, tracerName string, spanName string) (context.Context, trace.Span) {
	tr := otel.Tracer(tracerName)

	return tr.Start(ctx, spanName)
}
