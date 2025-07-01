// Package powervs extracts Power VS metadata from install configurations.
package powervs

import (
	"context"

	icpowervs "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// Metadata converts an install configuration to PowerVS metadata.
func Metadata(config *types.InstallConfig, meta *icpowervs.Metadata) (*powervs.Metadata, error) {
	cisCRN, _ := meta.CISInstanceCRN(context.TODO())
	dnsCRN, _ := meta.DNSInstanceCRN(context.TODO())

	overrides := config.Platform.PowerVS.ServiceEndpoints
	if config.Publish == types.InternalPublishingStrategy &&
		(len(config.ImageDigestSources) > 0 || len(config.DeprecatedImageContentSources) > 0) {
		cosRegion, err := powervs.COSRegionForPowerVSRegion(config.PowerVS.Region)
		if err != nil {
			return nil, err
		}
		vpcRegion, err := powervs.VPCRegionForPowerVSRegion(config.PowerVS.Region)
		if err != nil {
			return nil, err
		}
		overrides = meta.SetDefaultPrivateServiceEndpoints(context.TODO(), overrides, cosRegion, vpcRegion)
	}

	return &powervs.Metadata{
		BaseDomain:           config.BaseDomain,
		PowerVSResourceGroup: config.Platform.PowerVS.PowerVSResourceGroup,
		CISInstanceCRN:       cisCRN,
		DNSInstanceCRN:       dnsCRN,
		Region:               config.Platform.PowerVS.Region,
		VPCRegion:            config.Platform.PowerVS.VPCRegion,
		Zone:                 config.Platform.PowerVS.Zone,
		ServiceInstanceGUID:  config.Platform.PowerVS.ServiceInstanceGUID,
		ServiceEndpoints:     overrides,
		TransitGateway:       config.Platform.PowerVS.TransitGateway,
		VPC:                  config.Platform.PowerVS.VPC,
	}, nil
}
