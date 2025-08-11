// Package powervs extracts Power VS metadata from install configurations.
package powervs

import (
	"context"
	"time"

	icpowervs "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// Metadata converts an install configuration to PowerVS metadata.
func Metadata(config *types.InstallConfig, meta *icpowervs.Metadata) (*powervs.Metadata, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		err    error
	)

	// Update the saved session storage with the install config since the session
	// storage is used as the defaults.
	err = icpowervs.UpdateSessionStoreToAuthFile(&icpowervs.SessionStore{
		ID:                   config.PowerVS.UserID,
		DefaultRegion:        config.PowerVS.Region,
		DefaultZone:          config.PowerVS.Zone,
		PowerVSResourceGroup: config.PowerVS.PowerVSResourceGroup,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	cisCRN, err := meta.CISInstanceCRN(ctx)
	if err != nil {
		return nil, err
	}

	dnsCRN, err := meta.DNSInstanceCRN(ctx)
	if err != nil {
		return nil, err
	}

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
		overrides = meta.SetDefaultPrivateServiceEndpoints(ctx, overrides, cosRegion, vpcRegion)
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
		TransitGatewayName:   config.Platform.PowerVS.TGName,
	}, nil
}
