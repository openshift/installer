/*
Copyright 2020 The Machine API Operator authors

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
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	DefaultHealthCheckMetricsAddress = ":8083"
)

var (
	// MachineHealthCheckNodesCovered is a Prometheus metric, which reports the number of nodes covered by MachineHealthChecks
	MachineHealthCheckNodesCovered = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mapi_machinehealthcheck_nodes_covered",
			Help: "Number of nodes covered by MachineHealthChecks",
		}, []string{"name", "namespace"},
	)

	// MachineHealthCheckRemediationSuccessTotal is a Prometheus metric, which reports the number of successful remediations by MachineHealthChecks
	MachineHealthCheckRemediationSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mapi_machinehealthcheck_remediation_success_total",
			Help: "Number of successful remediations performed by MachineHealthChecks",
		}, []string{"name", "namespace"},
	)

	// MachineHealthCheckShortCircuit is a Prometheus metric, which reports when the named MachineHealthCheck is currently short-circuited (0=no, 1=yes)
	MachineHealthCheckShortCircuit = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mapi_machinehealthcheck_short_circuit",
			Help: "Short circuit status for MachineHealthCheck (0=no, 1=yes)",
		}, []string{"name", "namespace"},
	)
)

func InitializeMachineHealthCheckMetrics() {
	metrics.Registry.MustRegister(
		MachineHealthCheckNodesCovered,
		MachineHealthCheckRemediationSuccessTotal,
		MachineHealthCheckShortCircuit,
	)
}

func DeleteMachineHealthCheckNodesCovered(name string, namespace string) {
	MachineHealthCheckNodesCovered.Delete(prometheus.Labels{
		"name":      name,
		"namespace": namespace,
	})
}

func ObserveMachineHealthCheckNodesCovered(name string, namespace string, count int) {
	MachineHealthCheckNodesCovered.With(prometheus.Labels{
		"name":      name,
		"namespace": namespace,
	}).Set(float64(count))
}

func ObserveMachineHealthCheckRemediationSuccess(name string, namespace string) {
	MachineHealthCheckRemediationSuccessTotal.With(prometheus.Labels{
		"name":      name,
		"namespace": namespace,
	}).Inc()
}

func ObserveMachineHealthCheckShortCircuitDisabled(name string, namespace string) {
	MachineHealthCheckShortCircuit.With(prometheus.Labels{
		"name":      name,
		"namespace": namespace,
	}).Set(0)
}

func ObserveMachineHealthCheckShortCircuitEnabled(name string, namespace string) {
	MachineHealthCheckShortCircuit.With(prometheus.Labels{
		"name":      name,
		"namespace": namespace,
	}).Set(1)
}
