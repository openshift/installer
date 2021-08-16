package gatherer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsCollection(t *testing.T) {
	cases := []struct {
		enableMetrics  bool
		labelValueMap  map[string]string
		metricName     string
		testName       string
		value          float64
		expectedOutput string
	}{
		{
			enableMetrics: true,
			labelValueMap: map[string]string{"test1": "test1"},
			metricName:    clusterMetricName,
			testName:      "basic histogram metric test, incorrect labels",
			value:         10,
			expectedOutput: `# HELP cluster_installation_create This metric keeps track of the count of the number of times the user ran create cluster command in the given OS took the time that is lesser than or equal to the value in the duration label.
# TYPE cluster_installation_create histogram
cluster_installation_create_bucket{le="15"} 1
cluster_installation_create_bucket{le="20"} 1
cluster_installation_create_bucket{le="25"} 1
cluster_installation_create_bucket{le="30"} 1
cluster_installation_create_bucket{le="35"} 1
cluster_installation_create_bucket{le="40"} 1
cluster_installation_create_bucket{le="45"} 1
cluster_installation_create_bucket{le="50"} 1
cluster_installation_create_bucket{le="55"} 1
cluster_installation_create_bucket{le="60"} 1
cluster_installation_create_bucket{le="65"} 1
cluster_installation_create_bucket{le="70"} 1
cluster_installation_create_bucket{le="75"} 1
cluster_installation_create_bucket{le="+Inf"} 1
cluster_installation_create_sum 10
cluster_installation_create_count 1
`,
		},
		{
			enableMetrics: true,
			labelValueMap: map[string]string{"os": "linux"},
			metricName:    clusterMetricName,
			testName:      "basic histogram metric test, correct labels",
			value:         10,
			expectedOutput: `# HELP cluster_installation_create This metric keeps track of the count of the number of times the user ran create cluster command in the given OS took the time that is lesser than or equal to the value in the duration label.
# TYPE cluster_installation_create histogram
cluster_installation_create_bucket{os="linux",le="15"} 1
cluster_installation_create_bucket{os="linux",le="20"} 1
cluster_installation_create_bucket{os="linux",le="25"} 1
cluster_installation_create_bucket{os="linux",le="30"} 1
cluster_installation_create_bucket{os="linux",le="35"} 1
cluster_installation_create_bucket{os="linux",le="40"} 1
cluster_installation_create_bucket{os="linux",le="45"} 1
cluster_installation_create_bucket{os="linux",le="50"} 1
cluster_installation_create_bucket{os="linux",le="55"} 1
cluster_installation_create_bucket{os="linux",le="60"} 1
cluster_installation_create_bucket{os="linux",le="65"} 1
cluster_installation_create_bucket{os="linux",le="70"} 1
cluster_installation_create_bucket{os="linux",le="75"} 1
cluster_installation_create_bucket{os="linux",le="+Inf"} 1
cluster_installation_create_sum{os="linux"} 10
cluster_installation_create_count{os="linux"} 1
`,
		},
		{
			enableMetrics:  false,
			labelValueMap:  map[string]string{"test1": "test1"},
			metricName:     clusterMetricName,
			testName:       "basic disabled metrics test",
			value:          10,
			expectedOutput: "error should not work",
		},
		{
			enableMetrics: true,
			labelValueMap: map[string]string{"result": "success"},
			metricName:    ModificationBootstrapMetricName,
			testName:      "basic counter metric test",
			value:         10,
			expectedOutput: `# HELP cluster_installation_modification_bootstrap_ignition This metric keeps track of all the assets in the bootstrap ignition category that were modified by the user before the invocation of the create command in the installer
# TYPE cluster_installation_modification_bootstrap_ignition counter
cluster_installation_modification_bootstrap_ignition{result="success"} 10
`,
		},
		{
			enableMetrics: true,
			labelValueMap: nil,
			metricName:    ModificationBootstrapMetricName,
			testName:      "basic counter metric test, no labels",
			value:         10,
			expectedOutput: `# HELP cluster_installation_modification_bootstrap_ignition This metric keeps track of all the assets in the bootstrap ignition category that were modified by the user before the invocation of the create command in the installer
# TYPE cluster_installation_modification_bootstrap_ignition counter
cluster_installation_modification_bootstrap_ignition 10
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			if tc.enableMetrics {
				os.Setenv("OPENSHIFT_INSTALL_DISABLE_METRICS", "TRUE")
				testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					buf := new(strings.Builder)
					_, err := io.Copy(buf, req.Body)
					assert.NoError(t, err)
					assert.EqualValues(t, tc.expectedOutput, buf.String())
				}))
				defer testServer.Close()
				os.Setenv("OPENSHIFT_INSTALL_METRICS_ENDPOINT", testServer.URL)
			}

			Initialize()
			for key, value := range tc.labelValueMap {
				AddLabelValue(tc.metricName, key, value)
			}
			SetValue(tc.metricName, tc.value)
			Push(tc.metricName)
		})
	}
}

func TestPushAll(t *testing.T) {
	os.Setenv("OPENSHIFT_INSTALL_DISABLE_METRICS", "TRUE")
	expectedOutput := `# HELP cluster_installation_create This metric keeps track of the count of the number of times the user ran create cluster command in the given OS took the time that is lesser than or equal to the value in the duration label.
# TYPE cluster_installation_create histogram
cluster_installation_create_bucket{result="success",le="15"} 1
cluster_installation_create_bucket{result="success",le="20"} 1
cluster_installation_create_bucket{result="success",le="25"} 1
cluster_installation_create_bucket{result="success",le="30"} 1
cluster_installation_create_bucket{result="success",le="35"} 1
cluster_installation_create_bucket{result="success",le="40"} 1
cluster_installation_create_bucket{result="success",le="45"} 1
cluster_installation_create_bucket{result="success",le="50"} 1
cluster_installation_create_bucket{result="success",le="55"} 1
cluster_installation_create_bucket{result="success",le="60"} 1
cluster_installation_create_bucket{result="success",le="65"} 1
cluster_installation_create_bucket{result="success",le="70"} 1
cluster_installation_create_bucket{result="success",le="75"} 1
cluster_installation_create_bucket{result="success",le="+Inf"} 1
cluster_installation_create_sum{result="success"} 10
cluster_installation_create_count{result="success"} 1
# HELP cluster_installation_duration_api This metric keeps track of the API stageof the installer create command execution and the time ittook to complete the given stage. The values are kept as labels
# TYPE cluster_installation_duration_api histogram
cluster_installation_duration_api_bucket{result="success",le="5"} 0
cluster_installation_duration_api_bucket{result="success",le="10"} 1
cluster_installation_duration_api_bucket{result="success",le="15"} 1
cluster_installation_duration_api_bucket{result="success",le="20"} 1
cluster_installation_duration_api_bucket{result="success",le="25"} 1
cluster_installation_duration_api_bucket{result="success",le="30"} 1
cluster_installation_duration_api_bucket{result="success",le="35"} 1
cluster_installation_duration_api_bucket{result="success",le="40"} 1
cluster_installation_duration_api_bucket{result="success",le="45"} 1
cluster_installation_duration_api_bucket{result="success",le="50"} 1
cluster_installation_duration_api_bucket{result="success",le="55"} 1
cluster_installation_duration_api_bucket{result="success",le="60"} 1
cluster_installation_duration_api_bucket{result="success",le="65"} 1
cluster_installation_duration_api_bucket{result="success",le="+Inf"} 1
cluster_installation_duration_api_sum{result="success"} 10
cluster_installation_duration_api_count{result="success"} 1
# HELP cluster_installation_modification_bootstrap_ignition This metric keeps track of all the assets in the bootstrap ignition category that were modified by the user before the invocation of the create command in the installer
# TYPE cluster_installation_modification_bootstrap_ignition counter
cluster_installation_modification_bootstrap_ignition{result="success"} 10
`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, req.Body)
		assert.NoError(t, err)
		assert.EqualValues(t, expectedOutput, buf.String())
	}))
	os.Setenv("OPENSHIFT_INSTALL_METRICS_ENDPOINT", testServer.URL)

	Initialize()
	AddLabelValue(clusterMetricName, "result", "success")
	AddLabelValue(DurationAPIMetricName, "result", "success")
	AddLabelValue(ModificationBootstrapMetricName, "result", "success")

	SetValue(clusterMetricName, 10)
	SetValue(DurationAPIMetricName, 10)
	SetValue(ModificationBootstrapMetricName, 10)

	PushAll()
	testServer.Close()
	gatherer.enableMetrics = false
	PushAll()
}

func TestMultipleAddLabelValues(t *testing.T) {
	expectedOutput := `# HELP cluster_installation_create This metric keeps track of the count of the number of times the user ran create cluster command in the given OS took the time that is lesser than or equal to the value in the duration label.
# TYPE cluster_installation_create histogram
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="15"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="20"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="25"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="30"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="35"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="40"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="45"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="50"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="55"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="60"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="65"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="70"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="75"} 1
cluster_installation_create_bucket{os="Linux",result="Success",version="4.7",le="+Inf"} 1
cluster_installation_create_sum{os="Linux",result="Success",version="4.7"} 10
cluster_installation_create_count{os="Linux",result="Success",version="4.7"} 1
`
	metricName := clusterMetricName
	labelValueMap := map[string]string{"os": "Linux", "result": "Success", "version": "4.7"}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, req.Body)
		assert.NoError(t, err)
		assert.EqualValues(t, expectedOutput, buf.String())
	}))
	defer testServer.Close()
	os.Setenv("OPENSHIFT_INSTALL_METRICS_ENDPOINT", testServer.URL)

	Initialize()
	AddLabelValues(metricName, labelValueMap)
	SetValue(metricName, 10)
	Push(metricName)
}
