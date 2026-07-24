package v3

import "strings"

const prismCentralService = "PRISM_CENTRAL"

// IsPrismCentral checks if the cluster is a prism central instance or not
// by checking if the service running on the cluster is PRISM_CENTRAL
func (cluster *ClusterIntentResponse) IsPrismCentral() bool {
	if cluster.Status == nil ||
		cluster.Status.Resources == nil ||
		cluster.Status.Resources.Config == nil ||
		cluster.Status.Resources.Config.ServiceList == nil ||
		len(cluster.Status.Resources.Config.ServiceList) == 0 {
		return false
	}

	for _, service := range cluster.Status.Resources.Config.ServiceList {
		if service != nil && strings.EqualFold(*service, prismCentralService) {
			return true
		}
	}

	return false
}

// GetPrismElements returns the prism elements from the cluster list response
func (clusters *ClusterListIntentResponse) GetPrismElements() []*ClusterIntentResponse {
	var prismElements []*ClusterIntentResponse
	for _, cluster := range clusters.Entities {
		if cluster.IsPrismCentral() {
			continue
		}

		prismElements = append(prismElements, cluster)
	}

	return prismElements
}
