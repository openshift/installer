package builder

import (
	"errors"
	"fmt"
	"sort"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricBuilder is the basic metric type that must be used by any potential metric added to
// the system. It allows user to create a prom collector which will be used to push.
type MetricBuilder struct {
	// labels contains the list of all label keys that the collector object will have.
	labels []string
	// labelKeyValues stores values for label keys that will be added to the metric
	// Prometheus metrics are designed as follows:
	// metric_name {metric_label_key1:value1, key2:value2, ....} metric_value
	// labelKeyValues stores the key/value pairs of the Prometheus metric.
	labelKeyValues map[string]string
	// desc describes the metric that is being created. This field will be shown on Prometheus
	// dashboard as a help text.
	desc string
	// name describes the name of the metric that needs to be pushed.
	name string
	// value stores the value of the metric that needs to be pushed. Collected from the installer.
	value float64
	// buckets keeps a list of bucket values that will be used during the creation of a Histogram
	// object.
	buckets []float64
	// metricType defines what type of a collector object should the PromCollector function return.
	metricType MetricType
}

// MetricOpts contains the properties that are required to create a MetricBuilder object.
type MetricOpts struct {
	// labels contains the list of all label keys that the collector object will have.
	Labels []string
	// desc describes the metric that is being created. This field will be shown on Prometheus
	// dashboard as a help text.
	Desc string
	// name describes the name of the metric that needs to be pushed.
	Name string
	// buckets keeps a list of bucket values that will be used during the creation of a Histogram
	// object.
	Buckets []float64
	// metricType defines what type of a collector object should the PromCollector function return.
	MetricType MetricType
}

// MetricType defines what types of metrics can be created. Restricted by the types of the Prometheus
// Collector types.
type MetricType string

const (
	// Histogram denotes that the type of the collector object should be a Prometheus Histogram.
	Histogram MetricType = "Histogram"
	// Counter denotes that the type of the collector object should be a Prometheus Counter.
	Counter MetricType = "Counter"
)

// PromCollector function creates the required prometheus collector object with the values
// it has in the MetricBuilder calling object.
func (m MetricBuilder) PromCollector() (prometheus.Collector, error) {
	switch m.metricType {
	case Counter:
		return m.buildCounter(), nil
	case Histogram:
		return m.buildHistogram(), nil
	default:
		return nil, fmt.Errorf(`invalid metric builder type "%s". cannot create collector`, m.metricType)
	}
}

// buildCounter returns a prometheus counter object with the value and labels set
// in the MetricBuilder object.
func (m *MetricBuilder) buildCounter() prometheus.Collector {
	collector := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        m.name,
			Help:        m.desc,
			ConstLabels: m.labelKeyValues,
		},
	)
	collector.Add(m.value)
	return collector
}

// buildHistogram returns a prometheus Histogram object with the value, labels and buckets set
// in the MetricBuilder object.
func (m *MetricBuilder) buildHistogram() prometheus.Collector {
	collector := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        m.name,
			Help:        m.desc,
			Buckets:     m.buckets,
			ConstLabels: m.labelKeyValues,
		},
	)
	collector.Observe(m.value)
	return collector
}

// NewMetricBuilder creates a new MetricBuilder object with the default values for the field.
func NewMetricBuilder(opts MetricOpts, value float64, labelKeyValues map[string]string) (*MetricBuilder, error) {
	if opts.Labels == nil {
		return nil, errors.New("labels cannot be empty")
	}
	if labelKeyValues == nil {
		labelKeyValues = make(map[string]string)
	}
	if opts.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if opts.MetricType == "" {
		return nil, errors.New("metricType cannot be empty")
	}
	sort.Strings(opts.Labels)
	return &MetricBuilder{
		labels:         opts.Labels,
		name:           opts.Name,
		metricType:     opts.MetricType,
		desc:           opts.Desc,
		buckets:        opts.Buckets,
		value:          value,
		labelKeyValues: labelKeyValues,
	}, nil
}

// SetValue is a setter function that assigns value to the metric builder.
func (m *MetricBuilder) SetValue(value float64) {
	m.value = value
}

// AddLabelValue takes in a key and value and sets it to the map in metric builder.
func (m *MetricBuilder) AddLabelValue(key string, value string) error {
	if i := sort.SearchStrings(m.labels, key); i < len(m.labels) && m.labels[i] == key {
		m.labelKeyValues[key] = value
		return nil
	}
	return fmt.Errorf("key %q not in metricBuilder labels %v", key, m.labels)
}
