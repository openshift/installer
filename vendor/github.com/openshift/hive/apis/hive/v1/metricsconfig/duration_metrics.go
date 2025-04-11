package metricsconfig

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// MetricsWithDuration represents metrics that report time as values,like transition seconds.
// The purpose of these metrics should be to track outliers - ensure their duration is not set too low.
type MetricsWithDuration struct {
	// Name of the metric. It will correspond to an optional relevant metric in hive
	// +kubebuilder:validation:Enum=currentStopping;currentResuming;currentWaitingForCO;currentClusterSyncFailing;cumulativeHibernated;cumulativeResumed
	Name DurationMetricType `json:"name"`
	// Duration is the minimum time taken - the relevant metric will be logged only if the value reported by that metric
	// is more than the time mentioned here. For example, if a user opts-in for current clusters stopping and mentions
	// 1 hour here, only the clusters stopping for more than an hour will be reported.
	// This is a Duration value; see https://pkg.go.dev/time#ParseDuration for accepted formats.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	Duration *metav1.Duration `json:"duration"`
}

// DurationMetricType is a valid value for MetricsWithDuration.Name
type DurationMetricType string

const (
	// Metrics logged per cluster

	// CurrentStopping corresponds to hive_cluster_deployments_stopping_seconds
	CurrentStopping DurationMetricType = "currentStopping"
	// CurrentResuming corresponds to hive_cluster_deployments_resuming_seconds
	CurrentResuming DurationMetricType = "currentResuming"
	// CurrentWaitingForCO corresponds to hive_cluster_deployments_waiting_for_cluster_operators_seconds
	CurrentWaitingForCO DurationMetricType = "currentWaitingForCO"
	// CurrentClusterSyncFailing corresponds to hive_clustersync_failing_seconds
	CurrentClusterSyncFailing DurationMetricType = "currentClusterSyncFailing"

	// These metrics will not be cleared and can potentially blow up the cardinality

	// CumulativeHibernated corresponds to hive_cluster_deployment_hibernation_transition_seconds
	CumulativeHibernated DurationMetricType = "cumulativeHibernated"
	// CumulativeResumed corresponds to hive_cluster_deployment_running_transition_seconds
	CumulativeResumed DurationMetricType = "cumulativeResumed"
)
