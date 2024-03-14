package aws

import (
	"context"
	"fmt"
	"net"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	utilscidr "github.com/openshift/installer/pkg/asset/manifests/capiutils/cidr"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	"github.com/openshift/installer/pkg/types"
)

type zoneConfigInput struct {
	InstallConfig *installconfig.InstallConfig
	Config        *types.InstallConfig
	Meta          *aws.Metadata
	Cluster       *capa.AWSCluster
	ClusterID     *installconfig.ClusterID
}

// setZones creates the CAPI NetworkSpec structures for managed or
// BYO VPC deployments from install-config.yaml.
func setZones(in *zoneConfigInput) error {
	if len(in.Config.AWS.Subnets) > 0 {
		return setZonesBYOVPC(in)
	} else {
		return setZonesManagedVPC(in)
	}
}

// setZonesManagedVPC creates the CAPI NetworkSpec.Subnets setting the
// desired subnets from install-config.yaml in the BYO VPC deployment.
func setZonesBYOVPC(in *zoneConfigInput) error {
	privateSubnets, err := in.Meta.PrivateSubnets(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to get private subnets")
	}
	for _, subnet := range privateSubnets {
		in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
			ID:               subnet.ID,
			CidrBlock:        subnet.CIDR,
			AvailabilityZone: subnet.Zone.Name,
			IsPublic:         subnet.Public,
		})
	}

	publicSubnets, err := in.Meta.PublicSubnets(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to get public subnets")
	}
	for _, subnet := range publicSubnets {
		in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
			ID:               subnet.ID,
			CidrBlock:        subnet.CIDR,
			AvailabilityZone: subnet.Zone.Name,
			IsPublic:         subnet.Public,
		})
	}

	vpc, err := in.Meta.VPC(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to get VPC")
	}
	in.Cluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
		ID: vpc,
	}

	return nil
}

// setZonesManagedVPC creates the CAPI NetworkSpec.VPC setting the
// desired zones from install-config.yaml in the managed VPC deployment.
func setZonesManagedVPC(in *zoneConfigInput) error {

	zones, err := extractZonesFromInstallConfig(in)
	if err != nil {
		return errors.Wrap(err, "failed to get availability zones")
	}

	mainCIDR := capiutils.CIDRFromInstallConfig(in.InstallConfig)

	// Fallback to available zones in the region.
	if len(zones) == 0 {
		// Q? Do we need to use standard query or leave CAPA choose the zones automatically?
		// zonesMeta, err := in.Config.AWS.AvailabilityZones(context.TODO())
		// if err != nil {
		// 	return errors.Wrap(err, "failed to get availability zones")
		// }
		// for _, zoneMeta := range zonesMeta {
		// 	zones = append(zones, &aws.Zone{Name: zoneMeta})
		// }

		// Leaving CAPA to discover zones
		in.Cluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
			CidrBlock:                  mainCIDR.String(),
			AvailabilityZoneUsageLimit: ptr.To(len(zones)),
			AvailabilityZoneSelection:  &capa.AZSelectionSchemeOrdered,
		}
		return nil
	}

	in.Cluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
		CidrBlock: mainCIDR.String(),
	}

	// Base subnets considering only private zones, leaving one block free to allow
	// future subnet expansions in Day-2.
	numSubnets := len(zones) + 1

	// Public subnets consumes one range from base blocks.
	isPublishingExternal := in.Config.Publish == types.ExternalPublishingStrategy
	if isPublishingExternal {
		numSubnets++
	}

	subnetsCIDRs, err := utilscidr.SplitIntoSubnetsIPv4(mainCIDR.String(), numSubnets)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve CIDR blocks for all private subnets")
	}
	var publicSubnetsCIDRs []*net.IPNet
	if isPublishingExternal {
		publicSubnetsCIDRs, err = utilscidr.SplitIntoSubnetsIPv4(subnetsCIDRs[len(zones)].String(), len(zones))
		if err != nil {
			return errors.Wrap(err, "unable to retrieve CIDR blocks for all public subnets")
		}
	}

	idxCIDR := 0
	// Q: Can we use the standard terraform name (without 'subnet') and tell CAPA
	// to query it for Control Planes?
	subnetNamePrefix := fmt.Sprintf("%s-subnet", in.ClusterID.InfraID)
	for _, zone := range zones {
		if len(subnetsCIDRs) < idxCIDR {
			return errors.Wrap(err, "unable to define CIDR blocks for all private subnets")
		}
		cidr := subnetsCIDRs[idxCIDR]
		in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
			AvailabilityZone: zone.Name,
			CidrBlock:        cidr.String(),
			ID:               fmt.Sprintf("%s-private-%s", subnetNamePrefix, zone.Name),
			IsPublic:         false,
		})
		if isPublishingExternal {
			if len(publicSubnetsCIDRs) < idxCIDR {
				return errors.Wrap(err, "unable to define CIDR blocks for all public subnets")
			}
			cidr = publicSubnetsCIDRs[idxCIDR]
			in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
				AvailabilityZone: zone.Name,
				CidrBlock:        cidr.String(),
				ID:               fmt.Sprintf("%s-public-%s", subnetNamePrefix, zone.Name),
				IsPublic:         true,
			})
		}
		idxCIDR++
	}

	return nil
}

// extractZonesFromInstallConfig extract all zones defined in the install-config,
// otherwise discover it based in the AWS metadata when none is defined.
// TODO: Open Question: What is the expected behavior when only one pool defines the
// zones? Should the cluster be limited to those zones? Eg when worker defines single
// zone, and no controlPlane.platform.aws.zones is defined.
func extractZonesFromInstallConfig(in *zoneConfigInput) ([]*aws.Zone, error) {

	var zones []*aws.Zone
	zonesMap := make(map[string]struct{})

	if in.Config == nil {
		return nil, errors.New("unable to retreive Config")
	}

	cfg := in.Config
	if cfg.ControlPlane != nil && cfg.ControlPlane.Platform.AWS != nil &&
		len(cfg.ControlPlane.Platform.AWS.Zones) > 0 {
		for _, zone := range cfg.ControlPlane.Platform.AWS.Zones {
			if _, ok := zonesMap[zone]; !ok {
				zonesMap[zone] = struct{}{}
				zones = append(zones, &aws.Zone{Name: zone})
			}
		}
	}

	for _, compute := range cfg.Compute {
		if len(compute.Platform.AWS.Zones) > 0 {
			for _, zone := range compute.Platform.AWS.Zones {
				if _, ok := zonesMap[zone]; !ok {
					zonesMap[zone] = struct{}{}
					zones = append(zones, &aws.Zone{Name: zone})
				}
			}
		}
	}
	return zones, nil
}
