/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type ARMClientMetrics struct {
	azureSuccessfulRequestsTotal *prometheus.CounterVec
	azureFailedRequestsTotal     *prometheus.CounterVec
	azureRequestsTime            *prometheus.HistogramVec
}

var _ Metrics = &ARMClientMetrics{}

func NewARMClientMetrics() *ARMClientMetrics {

	azureSuccessfulRequestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "azure_successful_requests_total",
		Help: "Total number of successful requests to azure",
	}, []string{"resource", "requestType", "responseCode"})

	azureFailedRequestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "azure_failed_requests_total",
		Help: "Total number of failed requests to azure",
	}, []string{"resource", "requestType"})

	azureRequestsTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "azure_requests_time_seconds",
		Help: "Length of time per ARM request",
	}, []string{"resource", "requestType"})

	return &ARMClientMetrics{
		azureSuccessfulRequestsTotal: azureSuccessfulRequestsTotal,
		azureFailedRequestsTotal:     azureFailedRequestsTotal,
		azureRequestsTime:            azureRequestsTime,
	}
}

// RegisterMetrics registers the collectors with prometheus server.
func (a *ARMClientMetrics) RegisterMetrics() {
	metrics.Registry.MustRegister(a.azureRequestsTime, a.azureSuccessfulRequestsTotal, a.azureFailedRequestsTotal)
}

// RecordAzureSuccessRequestsTotal records the total successful number requests to ARM by increasing the counter.
func (a *ARMClientMetrics) RecordAzureSuccessRequestsTotal(resourceName string, statusCode int, method string) {
	a.azureSuccessfulRequestsTotal.WithLabelValues(resourceName, method, strconv.Itoa(statusCode)).Inc()
}

// RecordAzureFailedRequestsTotal records the number of failed requests to ARM.
func (a *ARMClientMetrics) RecordAzureFailedRequestsTotal(resourceName string, method string) {
	a.azureFailedRequestsTotal.WithLabelValues(resourceName, method).Inc()
}

// RecordAzureRequestsTime records the round-trip time taken by the request to ARM.
func (a ARMClientMetrics) RecordAzureRequestsTime(resourceName string, requestTime time.Duration, method string) {
	a.azureRequestsTime.WithLabelValues(resourceName, method).Observe(requestTime.Seconds())
}
