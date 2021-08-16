package gatherer

import (
	"os"
	"runtime"

	"github.com/openshift/installer/pkg/metrics/timer"
	"github.com/openshift/installer/pkg/version"
)

// InitializeInvocationMetrics is a utility function that adds the common information to create
// an invocation metrics.
func InitializeInvocationMetrics(metricName string) {
	if _, ok := os.LookupEnv("OPENSHIFT_INSTALL_INVOKER"); ok {
		AddLabelValue(metricName, "invoker", "non-user")
	} else {
		AddLabelValue(metricName, "invoker", "user")
	}
	AddLabelValue(metricName, "result", "Success")
	AddLabelValue(metricName, "os", runtime.GOOS)
	CurrentInvocationContext = metricName
}

// LogError sets the result and returnCode to the right values and pushes the information to the
// gateway. Must be called before the installer breaks execution.
func LogError(err string, metricName string) {
	AddLabelValue(metricName, "result", err)
	SendPrometheusInvocationData(metricName)
}

// SendPrometheusInvocationData gets the timer information for the duration metric and pushes it to
// the gateway.
func SendPrometheusInvocationData(metricName string) {
	duration := timer.StopTimer(timer.TotalTimeElapsed)
	SetValue(metricName, duration.Minutes())
	version, err := version.Version()
	if err == nil {
		AddLabelValue(CurrentInvocationContext, "version", version)
	}
	PushAll()
}

// UpdateDurationMetricsWithError sets all the duration metrics to error message.
func UpdateDurationMetricsWithError(err string) {
	listOfDurationMetrics := []string{
		DurationAPIMetricName,
		DurationBootstrapMetricName,
		DurationInfrastructureMetricName,
		DurationOperatorsMetricName,
		DurationProvisioningMetricName,
	}

	for _, item := range listOfDurationMetrics {
		AddLabelValue(item, "result", err)
	}
}

// UpdateDurationMetrics initializes the duration metrics to success.
func UpdateDurationMetrics(platform string) {
	listOfDurationMetrics := []string{
		DurationAPIMetricName,
		DurationBootstrapMetricName,
		DurationInfrastructureMetricName,
		DurationOperatorsMetricName,
		DurationProvisioningMetricName,
	}

	for _, item := range listOfDurationMetrics {
		AddLabelValue(item, "platform", platform)
		AddLabelValue(item, "result", "success")
		version, err := version.Version()
		if err == nil {
			AddLabelValue(item, "version", version)
		}

	}
}
