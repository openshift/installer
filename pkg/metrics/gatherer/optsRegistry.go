package gatherer

import (
	"math"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/openshift/installer/pkg/metrics/builder"
)

var (
	// INVOCATION METRICS

	// clusterMetricName is the variable that stores the metric name for the create cluster details.
	clusterMetricName = "cluster_installation_create"
	// ignitionMetricName is the variable that stores the metric name for the create ignition details.
	ignitionMetricName = "cluster_installation_ignition"
	// manifestsMetricName is the variable that stores the metric name for the create manifests details.
	manifestsMetricName = "cluster_installation_manifests"
	// WaitforMetricName is the variable that stores the metric name for the waitfor command details.
	WaitforMetricName = "cluster_installation_waitfor"
	// GatherMetricName is the variable that stores the metric name for the gather command configs details.
	GatherMetricName = "cluster_installation_gather"
	// DestroyMetricName is the variable that stores the metric name for the destroy command configs details.
	DestroyMetricName = "cluster_installation_destroy"

	// DURATION METRICS

	// DurationProvisioningMetricName is the variable that stores the metric name for the duration provisioning details.
	DurationProvisioningMetricName = "cluster_installation_duration_provisioning"
	// DurationInfrastructureMetricName is the variable that stores the metric name for the duration infrastructure details.
	DurationInfrastructureMetricName = "cluster_installation_duration_infrastructure"
	// DurationOperatorsMetricName is the variable that stores the metric name for the duration operators details.
	DurationOperatorsMetricName = "cluster_installation_duration_operators"
	// DurationBootstrapMetricName is the variable that stores the metric name for the duration bootstrap details.
	DurationBootstrapMetricName = "cluster_installation_duration_bootstrap"
	// DurationAPIMetricName is the variable that stores the metric name for the duration API details.
	DurationAPIMetricName = "cluster_installation_duration_api"

	// MODIFICATION METRICS

	// ModificationConfigMetricName is the variable that stores the metric name for the modification of config manifests details.
	ModificationConfigMetricName = "cluster_installation_modification_config_manifest"
	// ModificationDNSMetricName is the variable that stores the metric name for the modification of DNS manifests details.
	ModificationDNSMetricName = "cluster_installation_modification_dns_manifest"
	// ModificationNetworkMetricName is the variable that stores the metric name for the modification of network manifests details.
	ModificationNetworkMetricName = "cluster_installation_modification_network_manifest"
	// ModificationSchedulerMetricName is the variable that stores the metric name for the modification of scheduler manifests details.
	ModificationSchedulerMetricName = "cluster_installation_modification_scheduler_manifest"
	// ModificationCVOMetricName is the variable that stores the metric name for the modification of CVO manifests details.
	ModificationCVOMetricName = "cluster_installation_modification_cvo_manifest"
	// ModificationEtcdMetricName is the variable that stores the metric name for the modification of etcd manifests details.
	ModificationEtcdMetricName = "cluster_installation_modification_etcd_manifest"
	// ModificationMachineConfigMetricName is the variable that stores the metric name for the modification of machine config manifests details.
	ModificationMachineConfigMetricName = "cluster_installation_modification_machineconfig_manifest"
	// ModificationCaMetricName is the variable that stores the metric name for the modification of ca manifests details.
	ModificationCaMetricName = "cluster_installation_modification_ca_manifest"
	// ModificationPullSecretMetricName is the variable that stores the metric name for the modification of pull secret manifests details.
	ModificationPullSecretMetricName = "cluster_installation_modification_pullsecret_manifest"
	// ModificationBootstrapMetricName is the variable that stores the metric name for the modification of bootstrap manifests details.
	ModificationBootstrapMetricName = "cluster_installation_modification_bootstrap_ignition"
	// ModificationMasterMetricName is the variable that stores the metric name for the modification of master manifests details.
	ModificationMasterMetricName = "cluster_installation_modification_master_ignition"
	// ModificationWorkerMetricName is the variable that stores the metric name for the modification of worker manifests details.
	ModificationWorkerMetricName = "cluster_installation_modification_worker_manifest"
)

var (
	// optsRegistry holds all the information that every metric needs to create a metric builder.
	optsRegistry = map[string]*builder.MetricOpts{
		clusterMetricName: {
			Name:   clusterMetricName,
			Labels: []string{"result", "platform", "os", "version"},
			Desc: "This metric keeps track of the count of the number of times " +
				"the user ran create cluster command in the given OS " +
				"took the time that is lesser than or equal to the value in the duration label.",
			Buckets:    prometheus.LinearBuckets(15, 5, int(math.Ceil(60/5))+1),
			MetricType: builder.Histogram,
		},
		ignitionMetricName: {
			Name:   ignitionMetricName,
			Labels: []string{"result", "platform", "os", "version"},
			Desc: "This metric keeps track of the count of the number of times " +
				"the user ran create ignition-config command in the given OS " +
				"took the time that is lesser than or equal to the value in the duration label.",
			Buckets:    prometheus.LinearBuckets(15, 5, int(math.Ceil(60/5))+1),
			MetricType: builder.Histogram,
		},
		manifestsMetricName: {
			Name:   manifestsMetricName,
			Labels: []string{"result", "platform", "os", "version"},
			Desc: "This metric keeps track of the count of the number of times " +
				"the user ran create manifests command in the given OS " +
				"took the time that is lesser than or equal to the value in the duration label.",
			Buckets:    prometheus.LinearBuckets(15, 5, int(math.Ceil(60/5))+1),
			MetricType: builder.Histogram,
		},
		WaitforMetricName: {
			Name:   WaitforMetricName,
			Labels: []string{"result", "platform", "os", "version"},
			Desc: "This metric keeps track of the count of the number of times " +
				"the user ran wait-for command in the given OS " +
				"took the time that is lesser than or equal to the value in the duration label.",
			Buckets:    prometheus.LinearBuckets(5, 5, int(math.Ceil(30/5))+1),
			MetricType: builder.Histogram,
		},
		DestroyMetricName: {
			Name:   DestroyMetricName,
			Labels: []string{"result", "platform", "os", "version"},
			Desc: "This metric keeps track of the count of the number of times " +
				"the user ran destroy command in the given OS " +
				"took the time that is lesser than or equal to the value in the duration label.",
			Buckets:    prometheus.LinearBuckets(5, 5, int(math.Ceil(30/5))+1),
			MetricType: builder.Histogram,
		},
		GatherMetricName: {
			Name:   GatherMetricName,
			Labels: []string{"result", "platform", "os", "version"},
			Desc: "This metric keeps track of the count of the number of times " +
				"the user ran gather command in the given OS " +
				"took the time that is lesser than or equal to the value in the duration label.",
			Buckets:    prometheus.LinearBuckets(5, 5, int(math.Ceil(30/5))+1),
			MetricType: builder.Histogram,
		},
		DurationAPIMetricName: {
			Name:   DurationAPIMetricName,
			Labels: []string{"result", "version", "platform"},
			Desc: "This metric keeps track of the API stage" +
				"of the installer create command execution and the time it" +
				"took to complete the given stage. The values are kept as labels",
			Buckets:    prometheus.LinearBuckets(5, 5, int(math.Ceil(60/5))+1),
			MetricType: builder.Histogram,
		},
		DurationBootstrapMetricName: {
			Name:   DurationBootstrapMetricName,
			Labels: []string{"result", "version", "platform"},
			Desc: "This metric keeps track of the bootstrap stage" +
				"of the installer create command execution and the time it" +
				"took to complete the given stage. The values are kept as labels",
			Buckets:    prometheus.LinearBuckets(5, 5, int(math.Ceil(60/5))+1),
			MetricType: builder.Histogram,
		},
		DurationInfrastructureMetricName: {
			Name:   DurationInfrastructureMetricName,
			Labels: []string{"result", "version", "platform"},
			Desc: "This metric keeps track of the infrastructure stage" +
				"of the installer create command execution and the time it" +
				"took to complete the given stage. The values are kept as labels",
			Buckets:    prometheus.LinearBuckets(5, 5, int(math.Ceil(60/5))+1),
			MetricType: builder.Histogram,
		},
		DurationOperatorsMetricName: {
			Name:   DurationOperatorsMetricName,
			Labels: []string{"result", "version", "platform"},
			Desc: "This metric keeps track of the operator stage" +
				"of the installer create command execution and the time it" +
				"took to complete the given stage. The values are kept as labels",
			Buckets:    prometheus.LinearBuckets(5, 5, int(math.Ceil(60/5))+1),
			MetricType: builder.Histogram,
		},
		DurationProvisioningMetricName: {
			Name:   DurationProvisioningMetricName,
			Labels: []string{"result", "version", "platform"},
			Desc: "This metric keeps track of the provisioning stage" +
				"of the installer create command execution and the time it" +
				"took to complete the given stage. The values are kept as labels",
			Buckets:    prometheus.LinearBuckets(5, 5, int(math.Ceil(60/5))+1),
			MetricType: builder.Histogram,
		},
		ModificationBootstrapMetricName: {
			Name:   ModificationBootstrapMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the bootstrap ignition category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationCVOMetricName: {
			Name:   ModificationCVOMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the CVO category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationCaMetricName: {
			Name:   ModificationCaMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the CA category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationConfigMetricName: {
			Name:   ModificationConfigMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the config category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationDNSMetricName: {
			Name:   ModificationDNSMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the DNS category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationEtcdMetricName: {
			Name:   ModificationEtcdMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the etcd category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationMachineConfigMetricName: {
			Name:   ModificationMachineConfigMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the machine config category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationMasterMetricName: {
			Name:   ModificationMasterMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the master category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationNetworkMetricName: {
			Name:   ModificationNetworkMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the network category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationPullSecretMetricName: {
			Name:   ModificationPullSecretMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the pull secret category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationSchedulerMetricName: {
			Name:   ModificationSchedulerMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the scheduler category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
		ModificationWorkerMetricName: {
			Name:   ModificationWorkerMetricName,
			Labels: []string{"result"},
			Desc: "This metric keeps track of all the assets in the worker category that were modified " +
				"by the user before the invocation of the create command in the installer",
			Buckets:    nil,
			MetricType: builder.Counter,
		},
	}
)
