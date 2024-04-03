package clusterapi

import (
	"fmt"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
)

func mergeLabels(ic *installconfig.InstallConfig, clusterID string) map[string]string {
	labels := map[string]string{}
	labels[fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, clusterID)] = "owned"
	for _, label := range ic.Config.GCP.UserLabels {
		labels[label.Key] = label.Value
	}

	return labels
}
