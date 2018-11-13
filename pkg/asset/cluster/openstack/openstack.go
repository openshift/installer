// Package openstack extracts OpenStack metadata from install
// configurations.
package openstack

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// Metadata converts an install configuration to OpenStack metadata.
func Metadata(config *types.InstallConfig) *openstack.Metadata {
	return &openstack.Metadata{
		Region: config.Platform.OpenStack.Region,
		Identifier: map[string]string{
			"tectonicClusterID": config.ClusterID,
		},
	}
}
