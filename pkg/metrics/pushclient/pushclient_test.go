package pushclient

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestPushMetricsToGateway(t *testing.T) {
	promCollector := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        "test",
			Help:        "test metrics",
			ConstLabels: map[string]string{"test1": "test1"},
		},
	)

	promCollector.Add(10)

	expectedOutput := `# HELP test test metrics
# TYPE test counter
test{test1="test1"} 10
`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, req.Body)
		assert.NoError(t, err)
		assert.EqualValues(t, expectedOutput, buf.String())
		fmt.Println("all done")
	}))
	defer testServer.Close()

	pushClient := PushClient{URL: testServer.URL, Client: &http.Client{}, JobName: "installer_metrics"}
	pushClient.Push(promCollector)

}

func TestPushMetricsToAllGateway(t *testing.T) {
	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        "test_counter",
			Help:        "test metrics",
			ConstLabels: map[string]string{"test1": "test1"},
		},
	)
	counter.Add(10)

	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        "test_histogram",
			Help:        "test metrics",
			ConstLabels: map[string]string{"test1": "test1"},
		},
	)
	histogram.Observe(10)

	promCollector := []prometheus.Collector{counter, histogram}

	expectedOutput := `# HELP test_counter test metrics
# TYPE test_counter counter
test_counter{test1="test1"} 10
# HELP test_histogram test metrics
# TYPE test_histogram histogram
test_histogram_bucket{test1="test1",le="0.005"} 0
test_histogram_bucket{test1="test1",le="0.01"} 0
test_histogram_bucket{test1="test1",le="0.025"} 0
test_histogram_bucket{test1="test1",le="0.05"} 0
test_histogram_bucket{test1="test1",le="0.1"} 0
test_histogram_bucket{test1="test1",le="0.25"} 0
test_histogram_bucket{test1="test1",le="0.5"} 0
test_histogram_bucket{test1="test1",le="1"} 0
test_histogram_bucket{test1="test1",le="2.5"} 0
test_histogram_bucket{test1="test1",le="5"} 0
test_histogram_bucket{test1="test1",le="10"} 1
test_histogram_bucket{test1="test1",le="+Inf"} 1
test_histogram_sum{test1="test1"} 10
test_histogram_count{test1="test1"} 1
`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, req.Body)
		assert.NoError(t, err)
		assert.EqualValues(t, expectedOutput, buf.String())
		fmt.Println("all done")
	}))
	defer testServer.Close()
	pushClient := PushClient{URL: testServer.URL, Client: &http.Client{}, JobName: "installer_metrics"}
	err := pushClient.Push(promCollector...)
	assert.NoError(t, err)
}

func TestPushMetricsBadServer(t *testing.T) {
	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        "test_counter",
			Help:        "test metrics",
			ConstLabels: map[string]string{"test1": "test1"},
		},
	)
	counter.Add(10)

	promCollector := []prometheus.Collector{counter}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "test error message", http.StatusForbidden)
	}))
	defer testServer.Close()
	pushClient := PushClient{URL: testServer.URL, Client: &http.Client{}, JobName: "installer_metrics"}
	err := pushClient.Push(promCollector...)

	errorMessage := fmt.Sprintf("failed to push metrics: unexpected status code 403 while pushing to %s/metrics/job/installer_metrics: test error message\n", testServer.URL)

	if assert.Error(t, err) {
		assert.EqualError(t, err, errorMessage)
	}
}
