/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file contains the implementations of a handler wrapper that generates Prometheus metrics.

package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// HandlerWrapperBuilder contains the data and logic needed to build a new metrics handler wrapper
// that creates HTTP handlers that generate the following Prometheus metrics:
//
//	<subsystem>_request_count - Number of API requests sent.
//	<subsystem>_request_duration_sum - Total time to send API requests, in seconds.
//	<subsystem>_request_duration_count - Total number of API requests measured.
//	<subsystem>_request_duration_bucket - Number of API requests organized in buckets.
//
// To set the subsystem prefix use the Subsystem method.
//
// The duration buckets metrics contain an `le` label that indicates the upper bound. For example if
// the `le` label is `1` then the value will be the number of requests that were processed in less
// than one second.
//
// The metrics will have the following labels:
//
//	method - Name of the HTTP method, for example GET or POST.
//	path - Request path, for example /api/clusters_mgmt/v1/clusters.
//	code - HTTP response code, for example 200 or 500.
//	apiservice - API service name, for example ocm-clusters-service.
//
// To calculate the average request duration during the last 10 minutes, for example, use a
// Prometheus expression like this:
//
//	rate(api_outbound_request_duration_sum[10m]) / rate(api_outbound_request_duration_count[10m])
//
// In order to reduce the cardinality of the metrics the path label is modified to remove the
// identifiers of the objects. For example, if the original path is .../clusters/123 then it will
// be replaced by .../clusters/-, and the values will be accumulated. The line returned by the
// metrics server will be like this:
//
//	     <subsystem>_request_count{code="200",method="GET",path="/api/clusters_mgmt/v1/clusters/-",
//			apiservice="ocm-clusters-service"} 56
//
// The meaning of that is that there were a total of 56 requests to get specific clusters,
// independently of the specific identifier of the cluster.
//
// The value of the `code` label will be zero when sending the request failed without a response
// code, for example if it wasn't possible to open the connection, or if there was a timeout waiting
// for the response.
//
// Note that setting this attribute is not enough to have metrics published, you also need to
// create and start a metrics server, as described in the documentation of the Prometheus library.
//
// Don't create objects of this type directly; use the NewHandlerWrapper function instead.
type HandlerWrapperBuilder struct {
	paths      []string
	subsystem  string
	registerer prometheus.Registerer
}

// HandlerWrapper contains the data and logic needed to wrap an HTTP handler with another one that
// generates Prometheus metrics.
type HandlerWrapper struct {
	paths           pathTree
	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

// handler is an HTTP handler that generates Prometheus metrics.
type handler struct {
	owner   *HandlerWrapper
	handler http.Handler
}

// Make sure that we implement the interface:
var _ http.Handler = (*handler)(nil)

// responseWriter is the HTTP response writer used to obtain the response code.
type responseWriter struct {
	code   int
	writer http.ResponseWriter
}

// Make sure that we implement the interface:
var _ http.ResponseWriter = (*responseWriter)(nil)

// NewHandlerWrapper creates a new builder that can then be used to configure and create a new
// metrics handler wrapper.
func NewHandlerWrapper() *HandlerWrapperBuilder {
	return &HandlerWrapperBuilder{
		registerer: prometheus.DefaultRegisterer,
	}
}

// Path adds a path that will be accepted as a value for the `path` label. By default all the paths
// of the API are already added. This is intended for additional pads, for example the path for
// token requests. If those paths aren't explicitly specified here then their metrics will be
// accumulated in the `/-` path.
func (b *HandlerWrapperBuilder) Path(value string) *HandlerWrapperBuilder {
	b.paths = append(b.paths, value)
	return b
}

// Subsystem sets the name of the subsystem that will be used by to register the metrics with
// Prometheus. For example, if the value is `api_inbound` then the following metrics will be
// registered:
//
//	api_inbound_request_count - Number of API requests sent.
//	api_inbound_request_duration_sum - Total time to send API requests, in seconds.
//	api_inbound_request_duration_count - Total number of API requests measured.
//	api_inbound_request_duration_bucket - Number of API requests organized in buckets.
//
// This is mandatory.
func (b *HandlerWrapperBuilder) Subsystem(value string) *HandlerWrapperBuilder {
	b.subsystem = value
	return b
}

// Registerer sets the Prometheus registerer that will be used to register the metrics. The default
// is to use the default Prometheus registerer and there is usually no need to change that. This is
// intended for unit tests, where it is convenient to have a registerer that doesn't interfere with
// the rest of the system.
func (b *HandlerWrapperBuilder) Registerer(value prometheus.Registerer) *HandlerWrapperBuilder {
	if value == nil {
		value = prometheus.DefaultRegisterer
	}
	b.registerer = value
	return b
}

// Build uses the information stored in the builder to create a new handler wrapper.
func (b *HandlerWrapperBuilder) Build() (result *HandlerWrapper, err error) {
	// Check parameters:
	if b.subsystem == "" {
		err = fmt.Errorf("subsystem is mandatory")
		return
	}

	// Register the request count metric:
	requestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: b.subsystem,
			Name:      "request_count",
			Help:      "Number of requests sent.",
		},
		requestLabelNames,
	)
	err = b.registerer.Register(requestCount)
	if err != nil {
		registered, ok := err.(prometheus.AlreadyRegisteredError)
		if ok {
			requestCount = registered.ExistingCollector.(*prometheus.CounterVec)
			err = nil //nolint:all
		} else {
			return
		}
	}

	// Create the path tree:
	paths := pathRoot.copy()
	for _, path := range b.paths {
		paths.add(path)
	}

	// Register the request duration metric:
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: b.subsystem,
			Name:      "request_duration",
			Help:      "Request duration in seconds.",
			Buckets: []float64{
				0.1,
				1.0,
				10.0,
				30.0,
			},
		},
		requestLabelNames,
	)
	err = b.registerer.Register(requestDuration)
	if err != nil {
		registered, ok := err.(prometheus.AlreadyRegisteredError)
		if ok {
			requestDuration = registered.ExistingCollector.(*prometheus.HistogramVec)
			err = nil
		} else {
			return
		}
	}

	// Create and populate the object:
	result = &HandlerWrapper{
		paths:           paths,
		requestCount:    requestCount,
		requestDuration: requestDuration,
	}

	return
}

// Wrap creates a new handler that wraps the given one and generates the Prometheus metrics.
func (w *HandlerWrapper) Wrap(h http.Handler) http.Handler {
	return &handler{
		owner:   w,
		handler: h,
	}
}

// ServeHTTP is the implementation of the HTTP handler interface.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We need to replace the response writer with a custom one that captures the response code
	// generated by the next handler:
	writer := responseWriter{
		code:   http.StatusOK,
		writer: w,
	}

	// Measure the time that it takes to process the request and send the response:
	start := time.Now()
	h.handler.ServeHTTP(&writer, r)
	elapsed := time.Since(start)

	// Update the metrics:
	path := r.URL.Path
	method := r.Method
	labels := prometheus.Labels{
		serviceLabelName: serviceLabel(path),
		methodLabelName:  methodLabel(method),
		pathLabelName:    pathLabel(h.owner.paths, path),
		codeLabelName:    codeLabel(writer.code),
	}
	h.owner.requestCount.With(labels).Inc()
	h.owner.requestDuration.With(labels).Observe(elapsed.Seconds())
}

// Header is part of the implementation of the http.ResponseWriter interface.
func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write is part of the implementation of the http.ResponseWriter interface.
func (w *responseWriter) Write(b []byte) (n int, err error) {
	n, err = w.writer.Write(b)
	return
}

// WriteHeader is part of the implementation of the http.ResponseWriter interface.
func (w *responseWriter) WriteHeader(code int) {
	w.code = code
	w.writer.WriteHeader(code)
}

// Flush is the implementation of the http.Flusher interface.
func (w *responseWriter) Flush() {
	flusher, ok := w.writer.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
