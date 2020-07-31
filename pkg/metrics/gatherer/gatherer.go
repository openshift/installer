package gatherer

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/openshift/installer/pkg/metrics/builder"
	"github.com/openshift/installer/pkg/metrics/pushclient"
)

// Gatherer represents a metrics sink that has the push client and the metrics storage information and
// uses them to store and push metrics to the aggregation gateway.
type Gatherer struct {
	pushClient     pushclient.PushClient
	metricRegistry map[string]*builder.MetricBuilder
	enableMetrics  bool
}

// Initialize function allocates the storage requirement and the push client for the gatherer
func (g *Gatherer) Initialize() {
	g.enableMetrics = false
	if value, ok := os.LookupEnv("OPENSHIFT_INSTALL_METRICS_ENDPOINT"); ok {
		prometheusURL := value
		g.enableMetrics = true
		g.pushClient = pushclient.PushClient{URL: prometheusURL, Client: &http.Client{}, JobName: "openshift_installer_metrics"}
	}
	g.metricRegistry = make(map[string]*builder.MetricBuilder)

}

// AddLabelValue adds the label key and value to the storage of the Gatherer object.
func (g *Gatherer) AddLabelValue(metricName string, labelName string, labelValue string) {
	if !g.enableMetrics {
		return
	}
	if _, found := optsRegistry[metricName]; !found {
		return
	}
	if _, found := g.metricRegistry[metricName]; !found {
		opts := optsRegistry[metricName]
		metricBuilder, err := builder.NewMetricBuilder(*opts, 0, nil)
		if err == nil {
			g.metricRegistry[metricName] = metricBuilder
		}
	}
	g.metricRegistry[metricName].AddLabelValue(labelName, labelValue)
}

// SetValue sets the value of the metric that is stored in the Gatherer object.
func (g *Gatherer) SetValue(metricName string, value float64) {
	if !g.enableMetrics {
		return
	}
	if _, found := optsRegistry[metricName]; !found {
		return
	}
	if _, found := g.metricRegistry[metricName]; !found {
		opts := optsRegistry[metricName]
		metricBuilder, err := builder.NewMetricBuilder(*opts, 0, nil)
		if err == nil {
			g.metricRegistry[metricName] = metricBuilder
		}
	}
	g.metricRegistry[metricName].SetValue(value)
}

// Push is a driver function that sends the metrics created with all the information and pushes it
// to Prometheus through the Push Client.
func (g *Gatherer) Push(metricName string) {
	if !g.enableMetrics {
		return
	}
	if _, found := optsRegistry[metricName]; !found {
		return
	}
	if _, found := g.metricRegistry[metricName]; found {
		collector, err := g.metricRegistry[metricName].PromCollector()
		if err == nil {
			g.pushClient.Push(collector)
		}

	}
}

// PushAll function takes all the metrics whose label values are set and pushes them to Prometheus.
func (g *Gatherer) PushAll() {
	if !g.enableMetrics {
		return
	}
	var collectors []prometheus.Collector
	for _, value := range g.metricRegistry {
		collector, err := value.PromCollector()
		if err == nil {
			collectors = append(collectors, collector)
		}

	}
	g.pushClient.Push(collectors...)
}

var (
	// gatherer is the global Gatherer object that will be used to store all the information if this
	// package's global functions are called for collection of metrics.
	gatherer Gatherer

	// CurrentInvocationContext is used to store the current Invocation metric name used for building context.
	// Invocation data information is gathered from different parts of the installer and it gets hard for the
	// installer to keep track of the current command and target that is being run. Hence, the current invocation
	// metric name is stored here in this variable.
	CurrentInvocationContext string
)

// Initialize holds all the initial setup information like creating all the metrics that are allowed to
// be extracted, checking for user input about disabling metrics, URL and for creating the Push Client.
// Must be run before the metircs are collected.
func Initialize() {
	gatherer.Initialize()
}

// AddLabelValue adds a label key/value pair for a given metric. It keeps track of all the
// labels for all metrics if added through this function.
func AddLabelValue(metricName string, labelName string, labelValue string) {
	gatherer.AddLabelValue(metricName, labelName, labelValue)
}

// AddLabelValues adds all label key/value pair for a given metric. It keeps track of all the
// labels for all metrics if added through this function.
func AddLabelValues(metricName string, labelValuesMap map[string]string) {
	for labelName, labelValue := range labelValuesMap {
		gatherer.AddLabelValue(metricName, labelName, labelValue)
	}
}

// SetValue sets the value to the specific metric before sending to prometheus.
func SetValue(metricName string, value float64) {
	gatherer.SetValue(metricName, value)
}

// Push is a driver function that sends the metrics created with all the information and pushes it
// to Prometheus through the Push Client.
func Push(metricName string) {
	gatherer.Push(metricName)
}

// PushAll function takes all the metrics whose label values are set and pushes them to Prometheus.
func PushAll() {
	gatherer.PushAll()
}
