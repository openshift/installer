package aws

import (
	"context"
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/util/sets"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	utilscidr "github.com/openshift/installer/pkg/asset/manifests/capiutils/cidr"
	"github.com/openshift/installer/pkg/types"
)

type subnetsInput struct {
	vpc            string
	privateSubnets aws.Subnets
	publicSubnets  aws.Subnets
	edgeSubnets    aws.Subnets
}

type zonesInput struct {
	InstallConfig *installconfig.InstallConfig
	Cluster       *capa.AWSCluster
	ClusterID     *installconfig.ClusterID
	ZonesInRegion []string
	Subnets       *subnetsInput
}

// GatherZonesFromMetadata retrieves zones from AWS API to be used
// when building the subnets to CAPA.
func (zin *zonesInput) GatherZonesFromMetadata(ctx context.Context) (err error) {
	zin.ZonesInRegion, err = zin.InstallConfig.AWS.AvailabilityZones(ctx)
	if err != nil {
		return fmt.Errorf("failed to get availability zones: %w", err)
	}
	return nil
}

// GatherSubnetsFromMetadata retrieves subnets from AWS API to be used
// when building the subnets to CAPA.
func (zin *zonesInput) GatherSubnetsFromMetadata(ctx context.Context) (err error) {
	zin.Subnets = &subnetsInput{}
	if zin.Subnets.privateSubnets, err = zin.InstallConfig.AWS.PrivateSubnets(ctx); err != nil {
		return fmt.Errorf("failed to get private subnets: %w", err)
	}
	if zin.Subnets.publicSubnets, err = zin.InstallConfig.AWS.PublicSubnets(ctx); err != nil {
		return fmt.Errorf("failed to get public subnets: %w", err)
	}
	if zin.Subnets.edgeSubnets, err = zin.InstallConfig.AWS.EdgeSubnets(ctx); err != nil {
		return fmt.Errorf("failed to get public subnets: %w", err)
	}
	if zin.Subnets.vpc, err = zin.InstallConfig.AWS.VPC(ctx); err != nil {
		return fmt.Errorf("failed to get VPC: %w", err)
	}
	return nil
}

type zonesCAPI struct {
	controlPlaneZones sets.Set[string]
	computeZones      sets.Set[string]
	edgeZones         sets.Set[string]
}

// AvailabilityZones returns a sorted union of Availability Zones defined
// in the zone attribute in the pools for control plane and compute.
func (zo *zonesCAPI) AvailabilityZones() []string {
	return sets.List(zo.controlPlaneZones.Union(zo.computeZones))
}

// EdgeZones returns a sorted union of Local Zones or Wavelength Zones
// defined in the zone attribute in the edge compute pool.
func (zo *zonesCAPI) EdgeZones() []string {
	return sets.List(zo.edgeZones)
}

// SetAvailabilityZones insert the zone to the given compute pool, and to
// the regular zone (zone type availability-zone) list.
func (zo *zonesCAPI) SetAvailabilityZones(pool string, zones []string) {
	switch pool {
	case types.MachinePoolControlPlaneRoleName:
		zo.controlPlaneZones.Insert(zones...)

	case types.MachinePoolComputeRoleName:
		zo.computeZones.Insert(zones...)
	}
}

// SetDefaultConfigZones evaluates if machine pools (control plane and workers) have been
// set the zones from install-config.yaml, if not sets the default from platform, when exists,
// otherwise set the default from the region discovered from AWS API.
func (zo *zonesCAPI) SetDefaultConfigZones(pool string, defConfig []string, defRegion []string) {
	zones := []string{}
	switch pool {
	case types.MachinePoolControlPlaneRoleName:
		if len(zo.controlPlaneZones) == 0 && len(defConfig) > 0 {
			zones = defConfig
		} else if len(zo.controlPlaneZones) == 0 {
			zones = defRegion
		}
		zo.controlPlaneZones.Insert(zones...)

	case types.MachinePoolComputeRoleName:
		if len(zo.computeZones) == 0 && len(defConfig) > 0 {
			zones = defConfig
		} else if len(zo.computeZones) == 0 {
			zones = defRegion
		}
		zo.computeZones.Insert(zones...)
	}
}

// setSubnets is the entrypoint to create the CAPI NetworkSpec structures
// for managed or BYO VPC deployments from install-config.yaml.
// The NetworkSpec.Subnets will be populated with the desired zones.
func setSubnets(ctx context.Context, in *zonesInput) error {
	if in.InstallConfig == nil {
		return fmt.Errorf("failed to get installConfig")
	}
	if in.InstallConfig.AWS == nil {
		return fmt.Errorf("failed to get AWS metadata")
	}
	if in.InstallConfig.Config == nil {
		return fmt.Errorf("unable to get Config")
	}
	if in.Cluster == nil {
		return fmt.Errorf("failed to get AWSCluster config")
	}
	if len(in.InstallConfig.Config.AWS.Subnets) > 0 {
		if err := in.GatherSubnetsFromMetadata(ctx); err != nil {
			return fmt.Errorf("failed to get subnets from metadata: %w", err)
		}
		return setSubnetsBYOVPC(in)
	}

	if err := in.GatherZonesFromMetadata(ctx); err != nil {
		return fmt.Errorf("failed to get availability zones from metadata: %w", err)
	}
	return setSubnetsManagedVPC(in)
}

// setSubnetsBYOVPC creates the CAPI NetworkSpec.Subnets setting the
// desired subnets from install-config.yaml in the BYO VPC deployment.
// This function does not have support for unit test to mock for AWS API,
// so all API calls must be done prior this execution.
// TODO: create support to mock AWS API calls in the unit tests, so we can merge
// the methods GatherSubnetsFromMetadata() into this.
func setSubnetsBYOVPC(in *zonesInput) error {
	in.Cluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
		ID: in.Subnets.vpc,
	}
	for _, subnet := range in.Subnets.privateSubnets {
		in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
			ID:               subnet.ID,
			CidrBlock:        subnet.CIDR,
			AvailabilityZone: subnet.Zone.Name,
			IsPublic:         subnet.Public,
		})
	}

	for _, subnet := range in.Subnets.publicSubnets {
		in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
			ID:               subnet.ID,
			CidrBlock:        subnet.CIDR,
			AvailabilityZone: subnet.Zone.Name,
			IsPublic:         subnet.Public,
		})
	}

	// edgeSubnets are subnet created on AWS Local Zones or Wavelength Zone,
	// discovered by ID and zone-type attribute.
	for _, subnet := range in.Subnets.edgeSubnets {
		in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
			ID:               subnet.ID,
			CidrBlock:        subnet.CIDR,
			AvailabilityZone: subnet.Zone.Name,
			IsPublic:         subnet.Public,
		})
	}

	return nil
}

// setSubnetsManagedVPC creates the CAPI NetworkSpec.VPC and the NetworkSpec.Subnets,
// setting the desired zones from install-config.yaml in the managed
// VPC deployment, when specified, otherwise default zones are set from
// the previously discovered from AWS API.
// This function does not have mock for AWS API, so all API calls must be done prior
// this execution.
// TODO: create support to mock AWS API calls in the unit tests, so we can merge
// the methods GatherZonesFromMetadata() into this.
// The CIDR blocks are calculated leaving free blocks to allow future expansions,
// in Day-2, when desired.
func setSubnetsManagedVPC(in *zonesInput) error {
	out, err := extractZonesFromInstallConfig(in)
	if err != nil {
		return fmt.Errorf("failed to get availability zones: %w", err)
	}

	isPublishingExternal := in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy
	allAvailabilityZones := out.AvailabilityZones()
	allEdgeZones := out.EdgeZones()

	mainCIDR := capiutils.CIDRFromInstallConfig(in.InstallConfig)
	in.Cluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
		CidrBlock: mainCIDR.String(),
	}

	// Base subnets considering only private zones, leaving one free block to allow
	// future subnet expansions in Day-2.
	numSubnets := len(allAvailabilityZones) + 1

	// Public subnets consumes one range from base blocks.
	if isPublishingExternal {
		numSubnets++
	}

	// Edge subnets consumes one CIDR block from private CIDR, slicing it
	// into smaller depending on the amount edge zones added to install config.
	if len(allEdgeZones) > 0 {
		numSubnets++
	}

	privateCIDRs, err := utilscidr.SplitIntoSubnetsIPv4(mainCIDR.String(), numSubnets)
	if err != nil {
		return fmt.Errorf("unable to generate CIDR blocks for all private subnets: %w", err)
	}

	publicCIDR := privateCIDRs[len(allAvailabilityZones)].String()
	var edgeCIDR string
	if len(allEdgeZones) > 0 {
		publicCIDR = privateCIDRs[len(allAvailabilityZones)].String()
		edgeCIDR = privateCIDRs[len(allAvailabilityZones)+1].String()
	}

	var publicCIDRs []*net.IPNet
	if isPublishingExternal {
		// The last num(zones) blocks are dedicated to the public subnets.
		publicCIDRs, err = utilscidr.SplitIntoSubnetsIPv4(publicCIDR, len(allAvailabilityZones))
		if err != nil {
			return fmt.Errorf("unable to generate CIDR blocks for all public subnets: %w", err)
		}
	}

	// Create subnets from zone pools (control plane and compute) with type availability-zone.
	if len(privateCIDRs) < len(allAvailabilityZones) {
		return fmt.Errorf("unable to define CIDR blocks to all zones for private subnets")
	}
	if isPublishingExternal && len(publicCIDRs) < len(allAvailabilityZones) {
		return fmt.Errorf("unable to define CIDR blocks to all zones for public subnets")
	}

	for idxCIDR, zone := range allAvailabilityZones {
		in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
			AvailabilityZone: zone,
			CidrBlock:        privateCIDRs[idxCIDR].String(),
			ID:               fmt.Sprintf("%s-subnet-private-%s", in.ClusterID.InfraID, zone),
			IsPublic:         false,
		})
		if isPublishingExternal {
			in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
				AvailabilityZone: zone,
				CidrBlock:        publicCIDRs[idxCIDR].String(),
				ID:               fmt.Sprintf("%s-subnet-public-%s", in.ClusterID.InfraID, zone),
				IsPublic:         true,
			})
		}
	}

	// Create subnets from edge zone pool with type local-zone.
	if len(allEdgeZones) > 0 {
		// Slice the main CIDR (edgeCIDR) into N*zones for privates subnets,
		// and, when publish external, duplicate to create public subnets.
		numEdgeSubnets := len(allEdgeZones)
		if isPublishingExternal {
			numEdgeSubnets = numEdgeSubnets * 2
		}

		// Allow one CIDR block for future expansion.
		numEdgeSubnets++

		// Slice the edgeCIDR into the amount of desired subnets.
		edgeCIDRs, err := utilscidr.SplitIntoSubnetsIPv4(edgeCIDR, numEdgeSubnets)
		if err != nil {
			return fmt.Errorf("unable to generate CIDR blocks for all edge subnets: %w", err)
		}
		if len(edgeCIDRs) < len(allEdgeZones) {
			return fmt.Errorf("unable to define CIDR blocks to all edge zones for private subnets")
		}
		if isPublishingExternal && (len(edgeCIDRs) < (len(allEdgeZones) * 2)) {
			return fmt.Errorf("unable to define CIDR blocks to all edge zones for public subnets")
		}

		// Create subnets from zone pool with type local-zone or wavelength-zone (edge zones)
		for idxCIDR, zone := range allEdgeZones {
			in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
				AvailabilityZone: zone,
				CidrBlock:        edgeCIDRs[idxCIDR].String(),
				ID:               fmt.Sprintf("%s-subnet-private-%s", in.ClusterID.InfraID, zone),
				IsPublic:         false,
			})
			if isPublishingExternal {
				in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
					AvailabilityZone: zone,
					CidrBlock:        edgeCIDRs[len(allEdgeZones)+idxCIDR].String(),
					ID:               fmt.Sprintf("%s-subnet-public-%s", in.ClusterID.InfraID, zone),
					IsPublic:         true,
				})
			}
		}
	}

	return nil
}

// extractZonesFromInstallConfig extracts zones defined in the install-config.
func extractZonesFromInstallConfig(in *zonesInput) (*zonesCAPI, error) {
	out := zonesCAPI{
		controlPlaneZones: sets.New[string](),
		computeZones:      sets.New[string](),
		edgeZones:         sets.New[string](),
	}

	cfg := in.InstallConfig.Config
	defaultZones := []string{}
	if cfg.AWS != nil && cfg.AWS.DefaultMachinePlatform != nil && len(cfg.AWS.DefaultMachinePlatform.Zones) > 0 {
		defaultZones = cfg.AWS.DefaultMachinePlatform.Zones
	}

	if cfg.ControlPlane != nil && cfg.ControlPlane.Platform.AWS != nil {
		out.SetAvailabilityZones(types.MachinePoolControlPlaneRoleName, cfg.ControlPlane.Platform.AWS.Zones)
	}
	out.SetDefaultConfigZones(types.MachinePoolControlPlaneRoleName, defaultZones, in.ZonesInRegion)

	for _, pool := range cfg.Compute {
		if pool.Platform.AWS == nil {
			continue
		}
		if len(pool.Platform.AWS.Zones) > 0 {
			out.SetAvailabilityZones(pool.Name, pool.Platform.AWS.Zones)
		}
		// Ignoring as edge pool is not yet supported by CAPA.
		// See https://github.com/openshift/installer/pull/8173
		if pool.Name == types.MachinePoolEdgeRoleName {
			out.edgeZones.Insert(pool.Platform.AWS.Zones...)
			continue
		}
		out.SetDefaultConfigZones(types.MachinePoolComputeRoleName, defaultZones, in.ZonesInRegion)
	}
	return &out, nil
}
