package builder

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func getDescStruct(opts MetricOpts, labelKeyValues map[string]string) string {
	lpStrings := make([]string, 0, len(labelKeyValues))
	for key, value := range labelKeyValues {
		lpStrings = append(lpStrings, fmt.Sprintf("%s=%q", key, value))
	}
	sort.Strings(lpStrings)
	return fmt.Sprintf("Desc{fqName: %q, help: %q, constLabels: {%s}, variableLabels: %v}",
		opts.Name, opts.Desc, strings.Join(lpStrings, ","), "{}")
}

func getCollectorDescription(collector prometheus.Collector) string {
	switch reflect.TypeOf(collector).Elem().Name() {
	case "histogram":
		return collector.(prometheus.Histogram).Desc().String()
	case "counter":
		return collector.(prometheus.Counter).Desc().String()
	default:
		return ""
	}
}

// TestNewMetricBuilder tests the Metric Builder initializer.
func TestMetricBuilder(t *testing.T) {
	cases := []struct {
		name                  string
		expectedCollectorType string
		opts                  MetricOpts
		labelKeyValues        map[string]string
		value                 float64
		expectedErrorMessage  string
	}{
		{
			name:                  "Test histogram creation",
			expectedCollectorType: "histogram",
			opts: MetricOpts{
				Labels:     []string{"test1", "test2"},
				Desc:       "test histogram metrics",
				Name:       "test_histogram",
				Buckets:    []float64{10, 20, 30},
				MetricType: Histogram,
			},
			labelKeyValues:       map[string]string{"test1": "test1", "test2": "test2"},
			value:                0,
			expectedErrorMessage: "",
		},
		{
			name:                  "Test invalid labels",
			expectedCollectorType: "histogram",
			opts: MetricOpts{
				Labels:     nil,
				Desc:       "test metric",
				Name:       "test",
				MetricType: "Histogram",
			},
			labelKeyValues:       map[string]string{"test1": "test1", "test2": "test2"},
			value:                0,
			expectedErrorMessage: `labels cannot be empty`,
		},
		{
			name:                  "Test invalid name",
			expectedCollectorType: "histogram",
			opts: MetricOpts{
				Labels:     []string{"test1", "test2"},
				Desc:       "test metric",
				Name:       "",
				MetricType: "Histogram",
			},
			labelKeyValues:       map[string]string{"test1": "test1", "test2": "test2"},
			value:                0,
			expectedErrorMessage: `name cannot be empty`,
		},
		{
			name:                  "Test invalid metric type",
			expectedCollectorType: "histogram",
			opts: MetricOpts{
				Labels:     []string{"test1", "test2"},
				Desc:       "test metric",
				Name:       "test",
				MetricType: "",
			},
			labelKeyValues:       map[string]string{"test1": "test1", "test2": "test2"},
			value:                0,
			expectedErrorMessage: `metricType cannot be empty`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewMetricBuilder(tc.opts, tc.value, tc.labelKeyValues)
			if tc.expectedErrorMessage == "" {
				assert.NoError(t, err, "expected successful builder creation")
			} else {
				assert.Error(t, err, "expected builder creation failure but builder was created")
				assert.EqualError(t, err, tc.expectedErrorMessage)
			}
		})
	}
}

// TestMetricBuilderCollector tests the PromCollector function.
func TestMetricBuilderCollector(t *testing.T) {
	cases := []struct {
		name                  string
		expectedCollectorType string
		opts                  MetricOpts
		labelKeyValues        map[string]string
		value                 float64
		expectedErrorMessage  string
	}{
		{
			name:                  "Test histogram creation",
			expectedCollectorType: "histogram",
			opts: MetricOpts{
				Labels:     []string{"test1", "test2"},
				Desc:       "test histogram metrics",
				Name:       "test_histogram",
				Buckets:    []float64{10, 20, 30},
				MetricType: Histogram,
			},
			labelKeyValues:       map[string]string{"test1": "test1", "test2": "test2"},
			value:                0,
			expectedErrorMessage: "",
		},
		{
			name:                  "Test counter creation",
			expectedCollectorType: "counter",
			opts: MetricOpts{
				Labels:     []string{"test1", "test2"},
				Desc:       "test modification metric",
				Name:       "test_modification",
				MetricType: Counter,
			},
			labelKeyValues:       map[string]string{"test1": "test1", "test2": "test2"},
			expectedErrorMessage: "",
		},
		{
			name:                  "Test empty label values",
			expectedCollectorType: "counter",
			opts: MetricOpts{
				Labels:     []string{"test1", "test2"},
				Desc:       "test modification metric",
				Name:       "test_modification",
				MetricType: Counter,
			},
			labelKeyValues:       nil,
			value:                0,
			expectedErrorMessage: "",
		},
		{
			name:                  "Test invalid metricType",
			expectedCollectorType: "Linear",
			opts: MetricOpts{
				Labels:     []string{"test1", "test2"},
				Desc:       "test metric",
				Name:       "test",
				MetricType: "Linear",
			},
			labelKeyValues:       map[string]string{"test1": "test1", "test2": "test2"},
			value:                0,
			expectedErrorMessage: `invalid metric builder type "Linear". cannot create collector`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			metricBuilder, err := NewMetricBuilder(tc.opts, tc.value, tc.labelKeyValues)
			if !assert.NoError(t, err, "error constructing new metric builder") {
				return
			}
			collector, err := metricBuilder.PromCollector()
			if tc.expectedErrorMessage == "" {
				assert.NoError(t, err, "expected successful collector creation")
				assert.EqualValues(t, tc.expectedCollectorType, reflect.TypeOf(collector).Elem().Name())
				expectedString := getDescStruct(tc.opts, tc.labelKeyValues)
				assert.EqualValues(t, expectedString, getCollectorDescription(collector))
			} else {
				assert.Error(t, err, "expected collector creation failure but collector was created")
				assert.EqualError(t, err, tc.expectedErrorMessage)
			}
		})
	}
}

// TestNewMetricBuilder tests the Metric Builder initializer.
func TestNewMetricBuilder(t *testing.T) {
	opts := MetricOpts{
		Labels:     []string{"test1", "test2", "test3"},
		Desc:       "test metric",
		Name:       "test",
		Buckets:    []float64{10, 20, 30},
		MetricType: Histogram,
	}
	labelKeyValues := map[string]string{"test1": "test1", "test2": "test2"}
	metricBuilder, err := NewMetricBuilder(opts, 10, labelKeyValues)
	if err != nil {
		assert.Failf(t, err.Error(), "error creating metric builder")
	} else {
		assert.EqualValues(t, metricBuilder.value, 10)
		assert.EqualValues(t, metricBuilder.labels, opts.Labels)
		assert.EqualValues(t, metricBuilder.labelKeyValues, labelKeyValues)
		assert.EqualValues(t, metricBuilder.desc, opts.Desc)
		assert.EqualValues(t, metricBuilder.name, opts.Name)
		assert.EqualValues(t, metricBuilder.buckets, opts.Buckets)
		assert.EqualValues(t, metricBuilder.metricType, Histogram)
	}
	metricBuilder.SetValue(10)
	assert.EqualValues(t, 10, metricBuilder.value)

	err = metricBuilder.AddLabelValue("test", "metric labels")
	assert.Error(t, err, "expected error adding label")
	_, found := metricBuilder.labelKeyValues["test"]
	assert.False(t, found)

	err = metricBuilder.AddLabelValue("test3", "metric labels")
	if err != nil {
		assert.Failf(t, err.Error(), "error adding label value")
	}
	_, found = metricBuilder.labelKeyValues["test3"]
	assert.True(t, found)
}
