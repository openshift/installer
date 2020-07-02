// Package kubevirt extracts Kubevirt metadata from install configurations.
package kubevirt

import (
	kubevirtutils "github.com/openshift/cluster-api-provider-kubevirt/pkg/utils"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
)

// Metadata converts an install configuration to kubevirt metadata.
func Metadata(infraID string, config *types.InstallConfig) *kubevirt.Metadata {
	labels := kubevirtutils.BuildLabels(infraID)
	return &kubevirt.Metadata{
		Namespace: config.Kubevirt.Namespace,
		Labels:    labels,
	}
}
