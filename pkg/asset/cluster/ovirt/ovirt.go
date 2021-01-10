// Package ovirt extracts ovirt metadata from install configurations.
package ovirt

import (
	"os"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
)

// Metadata converts an install configuration to ovirt metadata.
func Metadata(config *types.InstallConfig) *ovirt.Metadata {
	_, ok := os.LookupEnv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE")
	m := ovirt.Metadata{
		ClusterID: config.Ovirt.ClusterID,
		// if we have a custom image, don't remove the template,
		// otherwise its a per deployment template, destroy it
		RemoveTemplate: !ok,
		APIVIP: config.Ovirt.APIVIP,
	}
	return &m
}
