// Package equinixmetal extracts equinixmetal metadata from install configurations.
package equinixmetal

import (
	"github.com/openshift/installer/pkg/types"
	equinixmetal "github.com/openshift/installer/pkg/types/equinixmetal"
)

// Metadata converts an install configuration to EquinixMetal metadata.
func Metadata(config *types.InstallConfig) *equinixmetal.Metadata {
	m := equinixmetal.Metadata{
		FacilityCode: config.Platform.EquinixMetal.FacilityCode,
		ProjectID:    config.Platform.EquinixMetal.ProjectID,
	}
	return &m
}
