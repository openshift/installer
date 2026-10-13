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

package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
)

type OpenstackPrometheusMetrics struct {
	Duration *prometheus.HistogramVec
	Total    *prometheus.CounterVec
	Errors   *prometheus.CounterVec
}

// MetricPrometheusContext indicates the context for OpenStack metrics.
type MetricPrometheusContext struct {
	Start      time.Time
	Attributes []string
	Metrics    *OpenstackPrometheusMetrics
}

// NewMetricPrometheusContext creates a new MetricContext.
func NewMetricPrometheusContext(resource string, request string) *MetricPrometheusContext {
	return &MetricPrometheusContext{
		Start:      time.Now(),
		Attributes: []string{resource + "_" + request},
	}
}

// ObserveRequest records the request latency and counts the errors.
func (mc *MetricPrometheusContext) ObserveRequest(err error) error {
	return mc.Observe(apiRequestPrometheusMetrics, err)
}

// ObserveRequestIgnoreNotFound records the request latency and counts the errors if it's not IsNotFound.
func (mc *MetricPrometheusContext) ObserveRequestIgnoreNotFound(err error) error {
	if capoerrors.IsNotFound(err) {
		_ = mc.ObserveRequest(nil)
		return err
	}
	return mc.ObserveRequest(err)
}

// ObserveRequestIgnoreNotFoundorConflict records the request latency and counts the errors if it's not IsNotFound or IsConflict.
func (mc *MetricPrometheusContext) ObserveRequestIgnoreNotFoundorConflict(err error) error {
	if capoerrors.IsNotFound(err) {
		_ = mc.ObserveRequest(nil)
		return err
	}
	if capoerrors.IsConflict(err) {
		_ = mc.ObserveRequest(nil)
		return err
	}
	return mc.ObserveRequest(err)
}

// Observe records the request latency and counts the errors.
func (mc *MetricPrometheusContext) Observe(om *OpenstackPrometheusMetrics, err error) error {
	if om == nil {
		// mc.RequestMetrics not set, ignore this request
		return err
	}

	om.Duration.WithLabelValues(mc.Attributes...).Observe(
		time.Since(mc.Start).Seconds())
	om.Total.WithLabelValues(mc.Attributes...).Inc()
	if err != nil {
		om.Errors.WithLabelValues(mc.Attributes...).Inc()
	}
	return err
}

var apiRequestPrometheusMetrics = &OpenstackPrometheusMetrics{
	Duration: prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "capo",
			Name:      "openstack_api_request_duration_seconds",
			Help:      "Latency of an OpenStack API call",
		}, []string{"request"}),
	Total: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "capo",
			Name:      "openstack_api_requests_total",
			Help:      "Total number of OpenStack API calls",
		}, []string{"request"}),
	Errors: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "capo",
			Name:      "openstack_api_request_errors_total",
			Help:      "Total number of errors for an OpenStack API call",
		}, []string{"request"}),
}

var registerAPIPrometheusMetrics sync.Once

func RegisterAPIPrometheusMetrics() {
	registerAPIPrometheusMetrics.Do(func() {
		metrics.Registry.MustRegister(apiRequestPrometheusMetrics.Duration)
		metrics.Registry.MustRegister(apiRequestPrometheusMetrics.Total)
		metrics.Registry.MustRegister(apiRequestPrometheusMetrics.Errors)
	})
}
