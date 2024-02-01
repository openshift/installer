package clusterapi

import (
	"fmt"

	"github.com/openshift/installer/pkg/asset/installconfig"
)

func mergeLabels(ic *installconfig.InstallConfig, clusterID *installconfig.ClusterID) map[string]string {
	labels := map[string]string{}
	labels[fmt.Sprintf("kubernetes-io-cluster-%s", clusterID.InfraID)] = "owned"
	for _, label := range ic.Config.GCP.UserLabels {
		labels[label.Key] = label.Value
	}

	return labels
}
