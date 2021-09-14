package gcp

import "github.com/openshift/installer/pkg/types/gcp"

// mergeAllUsage merges the update into the usage we know about, matching on metrics
func mergeAllUsage(into []gcp.QuotaUsage, update []gcp.QuotaUsage) []gcp.QuotaUsage {
	for _, item := range update {
		updated := false
		for i := range into {
			if into[i].Metric.Matches(item.Metric) {
				into[i].Amount += item.Amount
				updated = true
			}
		}
		if !updated {
			into = append(into, item)
		}
	}
	return into
}
