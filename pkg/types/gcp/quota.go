package gcp

import (
	"fmt"
	"strings"
)

const (
	// ServiceComputeEngineAPI is the GCE service URL
	ServiceComputeEngineAPI = "compute.googleapis.com"
	// ServiceIAMAPI is the IAM service URL
	ServiceIAMAPI = "iam.googleapis.com"
)

// Quota is a record of the quota in GCP consumed by a cluster
type Quota []QuotaUsage

// QuotaUsage identifies a quota metric and records the usage
type QuotaUsage struct {
	*Metric `json:",inline"`
	// Amount is the amount of the quota being used
	Amount int64 `json:"amount,omitempty"`
}

// String formats the quota usage
func (q *QuotaUsage) String() string {
	return fmt.Sprintf("%s:%d", q.Metric.String(), q.Amount)
}

// Metric identify a quota. Service/Label matches the Google Quota API names for quota metrics
type Metric struct {
	// Service is the Google Cloud Service to which this quota belongs (e.g. compute.googleapis.com)
	Service string `json:"service,omitempty"`
	// Limit is the name of the item that's limited (e.g. cpus)
	Limit string `json:"limit,omitempty"`
	// Dimensions are unique axes on which this Limit is applied (e.g. region: us-central-1)
	Dimensions map[string]string `json:"dimensions,omitempty"`
}

// String formats the metric
func (m *Metric) String() string {
	var dimensions []string
	for key, value := range m.Dimensions {
		dimensions = append(dimensions, fmt.Sprintf("%s=%s", key, value))
	}
	var suffix string
	if len(dimensions) > 0 {
		suffix = fmt.Sprintf("[%s]", strings.Join(dimensions, ","))
	}
	return fmt.Sprintf("%s/%s%s", m.Service, m.Limit, suffix)
}

// Matches determines if this metric matches the other
func (m *Metric) Matches(other *Metric) bool {
	if m.Service != other.Service {
		return false
	}
	if m.Limit != other.Limit {
		return false
	}

	if len(m.Dimensions) != len(other.Dimensions) {
		return false
	}
	for key, value := range m.Dimensions {
		otherValue, recorded := other.Dimensions[key]
		if !recorded {
			return false
		}
		if value != otherValue {
			return false
		}
	}
	for key, value := range other.Dimensions {
		ourValue, recorded := m.Dimensions[key]
		if !recorded {
			return false
		}
		if value != ourValue {
			return false
		}
	}
	return true
}
