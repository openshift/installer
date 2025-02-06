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
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// subnetsInput handles subnets information gathered from metadata.
type subnetsInput struct {
	vpc            string
	privateSubnets aws.Subnets
	publicSubnets  aws.Subnets
	edgeSubnets    aws.Subnets
}

// zonesInput handles input parameters required to create managed and unmanaged
// Subnets to CAPI.
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
		return fmt.Errorf("failed to get edge subnets: %w", err)
	}
	if zin.Subnets.vpc, err = zin.InstallConfig.AWS.VPC(ctx); err != nil {
		return fmt.Errorf("failed to get VPC: %w", err)
	}
	return nil
}

// ZonesCAPI handles the discovered zones used to create subnets to CAPA.
// ZonesCAPI is scoped in this package, but exported to use complex scenarios
// with go-cmp on unit tests.
type ZonesCAPI struct {
	ControlPlaneZones sets.Set[string]
	ComputeZones      sets.Set[string]
	EdgeZones         sets.Set[string]
}

// GetAvailabilityZones returns a sorted union of Availability Zones defined
// in the zone attribute in the pools for control plane and compute.
func (zo *ZonesCAPI) GetAvailabilityZones() []string {
	return sets.List(zo.ControlPlaneZones.Union(zo.ComputeZones))
}

// GetEdgeZones returns a sorted union of Local Zones or Wavelength Zones
// defined in the zone attribute in the edge compute pool.
func (zo *ZonesCAPI) GetEdgeZones() []string {
	return sets.List(zo.EdgeZones)
}

// SetAvailabilityZones insert the zone to the given compute pool, and to
// the regular zone (zone type availability-zone) list.
func (zo *ZonesCAPI) SetAvailabilityZones(pool string, zones []string) {
	switch pool {
	case types.MachinePoolControlPlaneRoleName:
		zo.ControlPlaneZones.Insert(zones...)

	case types.MachinePoolComputeRoleName:
		zo.ComputeZones.Insert(zones...)
	}
}

// SetDefaultConfigZones evaluates if machine pools (control plane and workers) have been
// set the zones from install-config.yaml, if not sets the default from platform, when exists,
// otherwise set the default from the region discovered from AWS API.
func (zo *ZonesCAPI) SetDefaultConfigZones(pool string, defConfig []string, defRegion []string) {
	zones := []string{}
	switch pool {
	case types.MachinePoolControlPlaneRoleName:
		if len(zo.ControlPlaneZones) == 0 && len(defConfig) > 0 {
			zones = defConfig
		} else if len(zo.ControlPlaneZones) == 0 {
			zones = defRegion
		}
		zo.ControlPlaneZones.Insert(zones...)

	case types.MachinePoolComputeRoleName:
		if len(zo.ComputeZones) == 0 && len(defConfig) > 0 {
			zones = defConfig
		} else if len(zo.ComputeZones) == 0 {
			zones = defRegion
		}
		zo.ComputeZones.Insert(zones...)
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

	// BYO VPC ("unmanaged") deployments
	if len(in.InstallConfig.Config.AWS.DeprecatedSubnets) > 0 {
		if err := in.GatherSubnetsFromMetadata(ctx); err != nil {
			return fmt.Errorf("failed to get subnets from metadata: %w", err)
		}
		return setSubnetsBYOVPC(in)
	}

	// Managed VPC (fully automated) deployments
	if err := in.GatherZonesFromMetadata(ctx); err != nil {
		return fmt.Errorf("failed to get availability zones from metadata: %w", err)
	}
	return setSubnetsManagedVPC(in)
}

// setSubnetsBYOVPC creates the CAPI NetworkSpec.Subnets setting the
// desired subnets from install-config.yaml in the BYO VPC deployment.
// This function does not provide support for unit test to mock for AWS API,
// so all API calls must be done prior this execution.
// TODO: create support to mock AWS API calls in the unit tests, then the method
// GatherSubnetsFromMetadata() can be added in setSubnetsBYOVPC.
func setSubnetsBYOVPC(in *zonesInput) error {
	in.Cluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
		ID: in.Subnets.vpc,
	}

	// Skip adding private subnets if this is a public-only subnets install.
	// We need to skip because the Installer is tricked into thinking the public subnets are also private and we would
	// end up adding public subnets twice to the cluster manifest, causing a duplicate error.
	if !awstypes.IsPublicOnlySubnetsEnabled() {
		for _, subnet := range in.Subnets.privateSubnets {
			in.Cluster.Spec.NetworkSpec.Subnets = append(in.Cluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
				ID:               subnet.ID,
				CidrBlock:        subnet.CIDR,
				AvailabilityZone: subnet.Zone.Name,
				IsPublic:         subnet.Public,
			})
		}
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
// the AWS API, previously discovered.
// The CIDR blocks are calculated leaving free blocks to allow future expansions,
// in Day-2, when desired.
// This function does not have mock for AWS API, so all API calls must be added prior
// this execution.
// TODO: create support to mock AWS API calls in the unit tests, then the method
// GatherZonesFromMetadata() can be added in setSubnetsManagedVPC.
func setSubnetsManagedVPC(in *zonesInput) error {
	out, err := extractZonesFromInstallConfig(in)
	if err != nil {
		return fmt.Errorf("failed to get availability zones: %w", err)
	}

	isPublishingExternal := in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy
	allAvailabilityZones := out.GetAvailabilityZones()
	allEdgeZones := out.GetEdgeZones()

	mainCIDR := capiutils.CIDRFromInstallConfig(in.InstallConfig)
	in.Cluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
		CidrBlock: mainCIDR.String(),
	}

	// Base subnets count considering only private zones, leaving one free block to allow
	// future subnet expansions in Day-2.
	numSubnets := len(allAvailabilityZones) + 1

	// Public subnets consumes one range from private CIDR block.
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

	// no edge zones, nothing else to do
	if len(allEdgeZones) == 0 {
		return nil
	}

	// Create subnets from edge zone pool with type local-zone.

	// Slice the main CIDR (edgeCIDR) into N*zones for privates subnets,
	// and, when publish external, duplicate to create public subnets.
	numEdgeSubnets := len(allEdgeZones)
	if isPublishingExternal {
		numEdgeSubnets *= 2
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

	return nil
}

// extractZonesFromInstallConfig extracts zones defined in the install-config.
func extractZonesFromInstallConfig(in *zonesInput) (*ZonesCAPI, error) {
	out := ZonesCAPI{
		ControlPlaneZones: sets.New[string](),
		ComputeZones:      sets.New[string](),
		EdgeZones:         sets.New[string](),
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

	// set the zones in the compute/worker pool, when defined, otherwise use defaults.
	for _, pool := range cfg.Compute {
		if pool.Platform.AWS == nil {
			continue
		}
		// edge compute pools should have zones defined.
		if pool.Name == types.MachinePoolEdgeRoleName {
			if len(pool.Platform.AWS.Zones) == 0 {
				return nil, fmt.Errorf("expect one or more zones in the edge compute pool, got: %q", pool.Platform.AWS.Zones)
			}
			out.EdgeZones.Insert(pool.Platform.AWS.Zones...)
			continue
		}

		if len(pool.Platform.AWS.Zones) > 0 {
			out.SetAvailabilityZones(pool.Name, pool.Platform.AWS.Zones)
		}
		out.SetDefaultConfigZones(types.MachinePoolComputeRoleName, defaultZones, in.ZonesInRegion)
	}

	// set defaults for worker pool when not defined in config.
	if len(out.ComputeZones) == 0 {
		out.SetDefaultConfigZones(types.MachinePoolComputeRoleName, defaultZones, in.ZonesInRegion)
	}

	// should raise an error if no zones is available in the pools, default platform config, or metadata.
	if azs := out.GetAvailabilityZones(); len(azs) == 0 {
		return nil, fmt.Errorf("failed to set zones from config, got: %q", azs)
	}

	return &out, nil
}
