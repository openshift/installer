/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package metrics

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/Azure/azure-service-operator/v2/internal/logging"
)

var _ policy.Policy = &MetricsPolicy{}

// MetricsPolicy is an Azure SDK Track 2 policy for logging metrics about HTTP status
type MetricsPolicy struct {
	metrics *ARMClientMetrics
}

func NewMetricsPolicy(metrics *ARMClientMetrics) policy.Policy {
	return &MetricsPolicy{
		metrics: metrics,
	}
}

func (m *MetricsPolicy) Do(req *policy.Request) (*http.Response, error) {
	raw := req.Raw()
	method := raw.Method
	path := raw.URL.Path

	id, err := arm.ParseResourceID(path)
	if err != nil {
		ctrl.Log.V(logging.Status).Error(err, "Error while parsing", "resourceID", path)
		// if error while parsing the resourceID, resourceType would be empty. We don't emit any metrics in this case
		return req.Next()
	}
	resourceType := id.ResourceType.String()

	requestStartTime := time.Now()
	resp, err := req.Next()

	m.metrics.RecordAzureRequestsTime(resourceType, time.Since(requestStartTime), method)
	if err != nil {
		m.metrics.RecordAzureFailedRequestsTotal(resourceType, method)
		return resp, err
	}

	m.metrics.RecordAzureSuccessRequestsTotal(resourceType, resp.StatusCode, method)

	return resp, err
}
