package metrics

import (
	machinev1 "github.com/openshift/api/machine/v1beta1"
	machineinformers "github.com/openshift/client-go/machine/informers/externalversions/machine/v1beta1"
	machinelisters "github.com/openshift/client-go/machine/listers/machine/v1beta1"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	DefaultMachineSetMetricsAddress = ":8082"
	DefaultMachineMetricsAddress    = ":8081"
	DefaultMetal3MetricsAddress     = ":60000"
)

var (
	// MachineCountDesc is a metric about machine object count in the cluster
	MachineCountDesc = prometheus.NewDesc("mapi_machine_items", "Count of machine objects currently at the apiserver", nil, nil)
	// MachineSetCountDesc Count of machineset object count at the apiserver
	MachineSetCountDesc = prometheus.NewDesc("mapi_machineset_items", "Count of machinesets at the apiserver", nil, nil)
	// MachineInfoDesc is a metric about machine object info in the cluster
	MachineInfoDesc = prometheus.NewDesc("mapi_machine_created_timestamp_seconds", "Timestamp of the mapi managed Machine creation time", []string{"name", "namespace", "spec_provider_id", "node", "api_version", "phase"}, nil)
	// MachineSetInfoDesc is a metric about machine object info in the cluster
	MachineSetInfoDesc = prometheus.NewDesc("mapi_machineset_created_timestamp_seconds", "Timestamp of the mapi managed Machineset creation time", []string{"name", "namespace", "api_version"}, nil)

	// MachineSetStatusAvailableReplicasDesc is the information of the Machineset's status for available replicas.
	MachineSetStatusAvailableReplicasDesc = prometheus.NewDesc("mapi_machine_set_status_replicas_available", "Information of the mapi managed Machineset's status for available replicas", []string{"name", "namespace"}, nil)

	// MachineSetStatusReadyReplicasDesc is the information of the Machineset's status for ready replicas.
	MachineSetStatusReadyReplicasDesc = prometheus.NewDesc("mapi_machine_set_status_replicas_ready", "Information of the mapi managed Machineset's status for ready replicas", []string{"name", "namespace"}, nil)

	// MachineSetStatusReplicasDesc is the information of the Machineset's status for replicas.
	MachineSetStatusReplicasDesc = prometheus.NewDesc("mapi_machine_set_status_replicas", "Information of the mapi managed Machineset's status for replicas", []string{"name", "namespace"}, nil)

	// MachineCollectorUp is a Prometheus metric, which reports reflects successful collection and reporting of all the metrics
	MachineCollectorUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mapi_mao_collector_up",
		Help: "Machine API Operator metrics are being collected and reported successfully",
	}, []string{"kind"})

	failedInstanceCreateCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mapi_instance_create_failed",
			Help: "Number of times provider instance create has failed.",
		}, []string{"name", "namespace", "reason"},
	)

	failedInstanceUpdateCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mapi_instance_update_failed",
			Help: "Number of times provider instance update has failed.",
		}, []string{"name", "namespace", "reason"},
	)

	failedInstanceDeleteCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mapi_instance_delete_failed",
			Help: "Number of times provider instance delete has failed.",
		}, []string{"name", "namespace", "reason"},
	)
)

// Metrics for use in the Machine controller
var (
	// MachinePhaseTransitionSeconds is a metric to capute the time between a Machine being created and entering a particular phase
	MachinePhaseTransitionSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mapi_machine_phase_transition_seconds",
			Help:    "Number of seconds between Machine creation and Machine transition to a phase.",
			Buckets: []float64{5, 10, 20, 30, 60, 90, 120, 180, 240, 300, 360, 480, 600},
		}, []string{"phase"},
	)
)

func init() {
	prometheus.MustRegister(MachineCollectorUp)
	metrics.Registry.MustRegister(MachinePhaseTransitionSeconds)
	metrics.Registry.MustRegister(
		failedInstanceCreateCount,
		failedInstanceUpdateCount,
		failedInstanceDeleteCount,
	)
}

// MachineCollector is implementing prometheus.Collector interface.
type MachineCollector struct {
	machineLister    machinelisters.MachineLister
	machineSetLister machinelisters.MachineSetLister
	namespace        string
}

// MachineLabels is the group of labels that are applied to the machine metrics
type MachineLabels struct {
	Name      string
	Namespace string
	Reason    string
}

func NewMachineCollector(machineInformer machineinformers.MachineInformer, machinesetInformer machineinformers.MachineSetInformer, namespace string) *MachineCollector {
	return &MachineCollector{
		machineLister:    machineInformer.Lister(),
		machineSetLister: machinesetInformer.Lister(),
		namespace:        namespace,
	}
}

// Collect is method required to implement the prometheus.Collector(prometheus/client_golang/prometheus/collector.go) interface.
func (mc *MachineCollector) Collect(ch chan<- prometheus.Metric) {
	mc.collectMachineMetrics(ch)
	mc.collectMachineSetMetrics(ch)
}

// Describe implements the prometheus.Collector interface.
func (mc MachineCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- MachineCountDesc
	ch <- MachineSetCountDesc
}

// Collect implements the prometheus.Collector interface.
func (mc MachineCollector) collectMachineMetrics(ch chan<- prometheus.Metric) {
	machineList, err := mc.listMachines()
	if err != nil {
		MachineCollectorUp.With(prometheus.Labels{"kind": "mapi_machine_items"}).Set(float64(0))
		return
	}
	MachineCollectorUp.With(prometheus.Labels{"kind": "mapi_machine_items"}).Set(float64(1))

	for _, machine := range machineList {
		nodeName := ""
		if machine.Status.NodeRef != nil {
			nodeName = machine.Status.NodeRef.Name
		}
		// Only gather metrics for machines with a phase.  This indicates
		// That the machine-controller is running on this cluster.
		phase := ptr.Deref(machine.Status.Phase, "")
		if phase != "" {
			ch <- prometheus.MustNewConstMetric(
				MachineInfoDesc,
				prometheus.GaugeValue,
				float64(machine.ObjectMeta.GetCreationTimestamp().Time.Unix()),
				machine.ObjectMeta.Name,
				machine.ObjectMeta.Namespace,
				ptr.Deref(machine.Spec.ProviderID, ""),
				nodeName,
				machine.TypeMeta.APIVersion,
				phase,
			)
		}
	}

	ch <- prometheus.MustNewConstMetric(MachineCountDesc, prometheus.GaugeValue, float64(len(machineList)))
	klog.V(4).Infof("collectmachineMetrics exit")
}

// collectMachineSetMetrics is method to collect machineSet related metrics.
func (mc MachineCollector) collectMachineSetMetrics(ch chan<- prometheus.Metric) {
	machineSetList, err := mc.listMachineSets()
	if err != nil {
		MachineCollectorUp.With(prometheus.Labels{"kind": "mapi_machineset_items"}).Set(float64(0))
		return
	}
	MachineCollectorUp.With(prometheus.Labels{"kind": "mapi_machineset_items"}).Set(float64(1))
	ch <- prometheus.MustNewConstMetric(MachineSetCountDesc, prometheus.GaugeValue, float64(len(machineSetList)))

	for _, machineSet := range machineSetList {

		ch <- prometheus.MustNewConstMetric(
			MachineSetInfoDesc,
			prometheus.GaugeValue,
			float64(machineSet.GetCreationTimestamp().Time.Unix()),
			machineSet.Name, machineSet.Namespace, machineSet.TypeMeta.APIVersion,
		)
		ch <- prometheus.MustNewConstMetric(
			MachineSetStatusAvailableReplicasDesc,
			prometheus.GaugeValue,
			float64(machineSet.Status.AvailableReplicas),
			machineSet.Name, machineSet.Namespace,
		)
		ch <- prometheus.MustNewConstMetric(
			MachineSetStatusReadyReplicasDesc,
			prometheus.GaugeValue,
			float64(machineSet.Status.ReadyReplicas),
			machineSet.Name, machineSet.Namespace,
		)
		ch <- prometheus.MustNewConstMetric(
			MachineSetStatusReplicasDesc,
			prometheus.GaugeValue,
			float64(machineSet.Status.Replicas),
			machineSet.Name, machineSet.Namespace,
		)
	}
}

func (mc MachineCollector) listMachines() ([]*machinev1.Machine, error) {
	return mc.machineLister.Machines(mc.namespace).List(labels.Everything())
}

func (mc MachineCollector) listMachineSets() ([]*machinev1.MachineSet, error) {
	return mc.machineSetLister.MachineSets(mc.namespace).List(labels.Everything())
}

func RegisterFailedInstanceCreate(labels *MachineLabels) {
	failedInstanceCreateCount.With(prometheus.Labels{
		"name":      labels.Name,
		"namespace": labels.Namespace,
		"reason":    labels.Reason,
	}).Inc()
}

func RegisterFailedInstanceUpdate(labels *MachineLabels) {
	failedInstanceUpdateCount.With(prometheus.Labels{
		"name":      labels.Name,
		"namespace": labels.Namespace,
		"reason":    labels.Reason,
	}).Inc()
}

func RegisterFailedInstanceDelete(labels *MachineLabels) {
	failedInstanceDeleteCount.With(prometheus.Labels{
		"name":      labels.Name,
		"namespace": labels.Namespace,
		"reason":    labels.Reason,
	}).Inc()
}
