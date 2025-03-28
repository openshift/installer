/*
Copyright 2021 The Kubernetes Authors.

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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// spanLogSink is a logr.LogSink implementation that writes log
// data to a span.
type spanLogSink struct {
	trace.Span
	name string
	vals []interface{}
}

func (*spanLogSink) Init(info logr.RuntimeInfo) {
}

func (s *spanLogSink) End(opts ...trace.SpanEndOption) {
	s.Span.End(opts...)
}

func (*spanLogSink) Enabled(v int) bool {
	return true
}

func (s *spanLogSink) kvsToAttrs(keysAndValues ...interface{}) []attribute.KeyValue {
	var ret []attribute.KeyValue
	for i := 0; i < len(keysAndValues); i += 2 {
		kv1 := fmt.Sprintf("%s", keysAndValues[i])
		kv2 := fmt.Sprintf("%s", keysAndValues[i+1])
		ret = append(ret, attribute.String(kv1, kv2))
	}
	for i := 0; i < len(s.vals); i += 2 {
		kv1 := fmt.Sprintf("%s", s.vals[i])
		kv2 := fmt.Sprintf("%s", s.vals[i+1])
		ret = append(ret, attribute.String(kv1, kv2))
	}
	return ret
}

func (s *spanLogSink) evtStr(evtType, msg string) string {
	return fmt.Sprintf(
		"[%s | %s] %s",
		evtType,
		s.name,
		msg,
	)
}

func (s *spanLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	attrs := s.kvsToAttrs(keysAndValues...)
	s.AddEvent(
		s.evtStr("INFO", msg),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(attrs...),
	)
}

func (s *spanLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	attrs := s.kvsToAttrs(keysAndValues...)
	s.AddEvent(
		s.evtStr("ERROR", fmt.Sprintf("%s (%s)", msg, err)),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(attrs...),
	)
}

func (s spanLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	// always create a new slice to avoid multiple loggers writing to the same backing array
	vals := make([]interface{}, len(s.vals)+len(keysAndValues))
	copy(vals, s.vals)
	copy(vals[len(s.vals):], keysAndValues)
	s.vals = vals
	return &s
}

func (s spanLogSink) WithName(name string) logr.LogSink {
	s.name = name
	return &s
}

// NewSpanLogSink is the main entry-point to this implementation.
func NewSpanLogSink(span trace.Span) logr.LogSink {
	return &spanLogSink{
		Span: span,
	}
}

// Config holds optional, arbitrary configuration information
// to be added to logs and telemetry data. Instances of
// Config get passed to StartSpanWithLogger via the KVP function.
type Config struct {
	KVPs map[string]string
}

func (c Config) teleKeyValues() []attribute.KeyValue {
	ret := make([]attribute.KeyValue, len(c.KVPs))
	i := 0
	for k, v := range c.KVPs {
		ret[i] = attribute.String(k, v)
		i++
	}
	return ret
}

// Option is the modifier function used to configure
// StartSpanWithLogger. Generally speaking, you should
// not create your own option function. Instead, use
// built-in functions (like KVP) that create them.
type Option func(*Config)

// KVP returns a new Option function that adds the given
// key-value pair.
func KVP(key, value string) Option {
	return func(cfg *Config) {
		cfg.KVPs[key] = value
	}
}

// StartSpanWithLogger starts a new span with the global
// tracer returned from Tracer(), then returns a new logger
// implementation that composes both the logger from the
// given ctx and a logger that logs to the newly created span.
//
// Callers should make sure to call the function in the 3rd return
// value to ensure that the span is ended properly. In many cases,
// that can be done with a defer:
//
//	ctx, lggr, done := StartSpanWithLogger(ctx, "my-span")
//	defer done()
func StartSpanWithLogger(
	ctx context.Context,
	spanName string,
	opts ...Option,
) (context.Context, logr.Logger, func()) {
	cfg := &Config{KVPs: make(map[string]string)}
	for _, opt := range opts {
		opt(cfg)
	}
	ctx, span := Tracer().Start(
		ctx,
		spanName,
		trace.WithAttributes(cfg.teleKeyValues()...),
	)
	endFn := func() {
		span.End()
	}

	kvs := make([]interface{}, 0, 2*len(cfg.KVPs))
	for k, v := range cfg.KVPs {
		kvs = append(kvs, k, v)
	}

	lggr := log.FromContext(ctx, kvs...).WithName(spanName)
	return ctx, NewCompositeLogger([]logr.LogSink{
		corrIDLogger(ctx, lggr).GetSink(),
		NewSpanLogSink(span),
	}), endFn
}
