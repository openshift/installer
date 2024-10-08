package metricsconfig

type MetricsConfig struct {
	// Optional metrics and their configurations
	// +optional
	MetricsWithDuration []MetricsWithDuration `json:"metricsWithDuration,omitempty"`
	// AdditionalClusterDeploymentLabels allows configuration of additional labels to be applied to certain metrics.
	// The keys can be any string value suitable for a metric label (see https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels).
	// The values can be any ClusterDeployment label key (from metadata.labels). When observing an affected metric,
	// hive will label it with the specified metric key, and copy the value from the specified ClusterDeployment label.
	// For example, including {"ocp_major_version": "hive.openshift.io/version-major"} will cause affected metrics to
	// include a label key ocp_major_version with the value from the hive.openshift.io/version-major ClusterDeployment
	// label -- e.g. "4".
	// NOTE: Avoid ClusterDeployment labels whose values are unbounded, such as those representing cluster names or IDs,
	// as these will cause your prometheus database to grow indefinitely.
	// Affected metrics are those whose type implements the metricsWithDynamicLabels interface found in
	// pkg/controller/metrics/metrics_with_dynamic_labels.go
	// +optional
	AdditionalClusterDeploymentLabels *map[string]string `json:"additionalClusterDeploymentLabels,omitempty"`
}
