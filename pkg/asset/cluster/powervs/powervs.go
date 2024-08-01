// Package powervs extracts Power VS metadata from install configurations.
package powervs

import (
	"context"

	icpowervs "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// Metadata converts an install configuration to PowerVS metadata.
func Metadata(config *types.InstallConfig, meta *icpowervs.Metadata) *powervs.Metadata {
	cisCRN, _ := meta.CISInstanceCRN(context.TODO())
	dnsCRN, _ := meta.DNSInstanceCRN(context.TODO())

	return &powervs.Metadata{
		BaseDomain:           config.BaseDomain,
		PowerVSResourceGroup: config.Platform.PowerVS.PowerVSResourceGroup,
		CISInstanceCRN:       cisCRN,
		DNSInstanceCRN:       dnsCRN,
		Region:               config.Platform.PowerVS.Region,
		VPCRegion:            config.Platform.PowerVS.VPCRegion,
		Zone:                 config.Platform.PowerVS.Zone,
		ServiceInstanceGUID:  config.Platform.PowerVS.ServiceInstanceGUID,
		ServiceEndpoints:     config.Platform.PowerVS.ServiceEndpoints,
	}
}
