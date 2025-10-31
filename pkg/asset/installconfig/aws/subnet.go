package aws

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// VPC holds metadata for a VPC.
type VPC struct {
	// ID is the VPC's Identifier.
	ID string

	// CIDR is the VPC's CIDR block.
	CIDR string

	// Tags is the map of the VPC's tags.
	Tags Tags
}

// Subnet holds metadata for a subnet.
type Subnet struct {
	// ID is the subnet's Identifier.
	ID string

	// ARN is the subnet's Amazon Resource Name.
	ARN string

	// Zone is the subnet's availability zone.
	Zone *Zone

	// CIDR is the subnet's CIDR block.
	CIDR string

	// IPv6CIDRAssociation is the subnet's IPv6 CIDR association.
	// It contains the allocated IPv6 CIDR and its association state.
	IPv6CIDRAssociation []ec2types.SubnetIpv6CidrBlockAssociation

	// Public is the flag to define the subnet public.
	Public bool

	// Tags is the map of the subnet's tags.
	Tags Tags

	// VPCID is the ID of the VPC containing the subnet.
	VPCID string
}

// Subnets is the map for the Subnet metadata indexed by subnetID.
type Subnets map[string]Subnet

// IDs returns the subnet IDs (i.e. map keys) in the Subnets.
func (sns Subnets) IDs() []string {
	subnetIDs := make([]string, 0)
	for id := range sns {
		subnetIDs = append(subnetIDs, id)
	}
	return subnetIDs
}

// SubnetsByZone is the map for the Subnet metadata indexed by zone.
type SubnetsByZone map[string]Subnet

// SubnetGroups is the group of subnets used by installer.
type SubnetGroups struct {
	Public  Subnets
	Private Subnets
	Edge    Subnets
	VpcID   string
}

// subnets retrieves metadata for the given subnet(s) or VPC.
func subnets(ctx context.Context, client *ec2.Client, subnetIDs []string, vpcID string) (SubnetGroups, error) {
	metas := make(Subnets, len(subnetIDs))
	zoneNames := make([]string, 0)
	availabilityZones := make(map[string]ec2types.AvailabilityZone, len(subnetIDs))
	subnetGroups := SubnetGroups{
		Public:  make(Subnets, len(subnetIDs)),
		Private: make(Subnets, len(subnetIDs)),
		Edge:    make(Subnets, len(subnetIDs)),
	}

	subnetInput := &ec2.DescribeSubnetsInput{}
	if len(vpcID) > 0 {
		subnetInput.Filters = append(subnetInput.Filters, ec2types.Filter{
			Name:   ptr.To("vpc-id"),
			Values: []string{vpcID},
		})
	}
	if len(subnetIDs) > 0 {
		subnetInput.SubnetIds = subnetIDs
	}

	err := describeSubnets(ctx, client, subnetInput, func(subnets []ec2types.Subnet) error {
		var vpcFromSubnet string
		for _, subnet := range subnets {
			if subnet.SubnetId == nil {
				continue
			}

			if len(ptr.Deref(subnet.SubnetArn, "")) == 0 {
				return fmt.Errorf("%s has no ARN", *subnet.SubnetId)
			}
			if len(ptr.Deref(subnet.VpcId, "")) == 0 {
				return fmt.Errorf("%s has no VPC", *subnet.SubnetId)
			}
			if len(ptr.Deref(subnet.AvailabilityZone, "")) == 0 {
				return fmt.Errorf("%s has no availability zone", *subnet.SubnetId)
			}
			if subnetGroups.VpcID == "" {
				subnetGroups.VpcID = *subnet.VpcId
				vpcFromSubnet = *subnet.SubnetId
			} else if *subnet.VpcId != subnetGroups.VpcID {
				return fmt.Errorf("all subnets must belong to the same VPC: %s is from %s, but %s is from %s", *subnet.SubnetId, *subnet.VpcId, vpcFromSubnet, subnetGroups.VpcID)
			}

			// At this point, we should be safe to dereference these fields.
			metas[*subnet.SubnetId] = Subnet{
				ID:                  *subnet.SubnetId,
				ARN:                 *subnet.SubnetArn,
				Zone:                &Zone{Name: *subnet.AvailabilityZone},
				CIDR:                ptr.Deref(subnet.CidrBlock, ""),
				Public:              false,
				Tags:                FromAWSTags(subnet.Tags),
				VPCID:               *subnet.VpcId,
				IPv6CIDRAssociation: subnet.Ipv6CidrBlockAssociationSet,
			}
			zoneNames = append(zoneNames, *subnet.AvailabilityZone)
		}
		return nil
	})
	if err != nil {
		return subnetGroups, err
	}

	var routeTables []ec2types.RouteTable
	err = describeRouteTables(ctx, client, &ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{{
			Name:   ptr.To("vpc-id"),
			Values: []string{subnetGroups.VpcID},
		}},
	}, func(rTables []ec2types.RouteTable) error {
		routeTables = append(routeTables, rTables...)
		return nil
	})
	if err != nil {
		return subnetGroups, err
	}

	azs, err := client.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{ZoneNames: zoneNames})
	if err != nil {
		return subnetGroups, fmt.Errorf("describing availability zones: %w", err)
	}
	for _, az := range azs.AvailabilityZones {
		availabilityZones[*az.ZoneName] = az
	}

	publicOnlySubnets := typesaws.IsPublicOnlySubnetsEnabled()

	var ids []string
	if len(vpcID) > 0 {
		ids = metas.IDs()
	}
	if len(subnetIDs) > 0 {
		ids = subnetIDs
	}

	for _, id := range ids {
		meta, ok := metas[id]
		if !ok {
			return subnetGroups, fmt.Errorf("failed to find %s", id)
		}

		isPublic, err := isSubnetPublic(routeTables, id)
		if err != nil {
			return subnetGroups, err
		}
		meta.Public = isPublic

		zoneName := meta.Zone.Name
		if _, ok := availabilityZones[zoneName]; !ok {
			return subnetGroups, fmt.Errorf("unable to read properties of zone name %s from the list %v: %w", zoneName, zoneNames, err)
		}
		zone := availabilityZones[zoneName]
		meta.Zone.Type = ptr.Deref(zone.ZoneType, "")
		meta.Zone.GroupName = ptr.Deref(zone.GroupName, "")
		if availabilityZones[zoneName].ParentZoneName != nil {
			meta.Zone.ParentZoneName = ptr.Deref(zone.ParentZoneName, "")
		}

		// AWS Local Zones are grouped as Edge subnets
		if meta.Zone.Type == typesaws.LocalZoneType ||
			meta.Zone.Type == typesaws.WavelengthZoneType {
			subnetGroups.Edge[id] = meta
			continue
		}
		if meta.Public {
			subnetGroups.Public[id] = meta

			// Let public subnets work as if they were private. This allows us to
			// have clusters with public-only subnets without having to introduce a
			// lot of changes in the installer. Such clusters can be used in a
			// NAT-less GW scenario, therefore decreasing costs in cases where node
			// security is not a concern (e.g, ephemeral clusters in CI)
			if publicOnlySubnets {
				subnetGroups.Private[id] = meta
			}
			continue
		}
		// Subnet is grouped by default as private
		subnetGroups.Private[id] = meta
	}
	return subnetGroups, nil
}

// https://github.com/kubernetes/kubernetes/blob/9f036cd43d35a9c41d7ac4ca82398a6d0bef957b/staging/src/k8s.io/legacy-cloud-providers/aws/aws.go#L3376-L3419
func isSubnetPublic(rt []ec2types.RouteTable, subnetID string) (bool, error) {
	var subnetTable *ec2types.RouteTable
	for _, table := range rt {
		for _, assoc := range table.Associations {
			if ptr.Equal(assoc.SubnetId, &subnetID) {
				subnetTable = &table
				break
			}
		}
	}

	if subnetTable == nil {
		// If there is no explicit association, the subnet will be implicitly
		// associated with the VPC's main routing table.
		for _, table := range rt {
			for _, assoc := range table.Associations {
				if ptr.Deref(assoc.Main, false) {
					logrus.Debugf("Assuming implicit use of main routing table %s for %s",
						ptr.Deref(table.RouteTableId, ""), subnetID)
					subnetTable = &table
					break
				}
			}
		}
	}

	if subnetTable == nil {
		return false, fmt.Errorf("could not locate routing table for %s", subnetID)
	}

	for _, route := range subnetTable.Routes {
		// There is no direct way in the AWS API to determine if a subnet is public or private.
		// A public subnet is one which has an internet gateway route
		// we look for the gatewayId and make sure it has the prefix of igw to differentiate
		// from the default in-subnet route which is called "local"
		// or other virtual gateway (starting with vgv)
		// or vpc peering connections (starting with pcx).
		if strings.HasPrefix(ptr.Deref(route.GatewayId, ""), "igw") {
			return true, nil
		}
		if strings.HasPrefix(ptr.Deref(route.CarrierGatewayId, ""), "cagw") {
			return true, nil
		}
	}

	return false, nil
}

// describeSubnets retrieves metadata for subnets with given filters.
func describeSubnets(ctx context.Context, client *ec2.Client, input *ec2.DescribeSubnetsInput, fn func(subnets []ec2types.Subnet) error) error {
	paginator := ec2.NewDescribeSubnetsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("describing subnets: %w", err)
		}

		// If the handler returns an error, we stop early to avoid extra API calls.
		if err = fn(page.Subnets); err != nil {
			return err
		}
	}
	return nil
}

// describeRouteTables retrieves metadata for route tables with given filters.
func describeRouteTables(ctx context.Context, client *ec2.Client, input *ec2.DescribeRouteTablesInput, fn func(subnets []ec2types.RouteTable) error) error {
	paginator := ec2.NewDescribeRouteTablesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("describing route tables: %w", err)
		}

		// If the handler returns an error, we stop early to avoid extra API calls.
		if err = fn(page.RouteTables); err != nil {
			return err
		}
	}
	return nil
}

// mergeSubnets merged two or more Subnets into a single one for convenience.
func mergeSubnets(groups ...Subnets) Subnets {
	subnets := make(Subnets)
	for _, group := range groups {
		maps.Copy(subnets, group)
	}
	return subnets
}

// vpc retrieves metadata for the given VPC ID.
func vpc(ctx context.Context, client *ec2.Client, vpcID string) (VPC, error) {
	vpcs := make([]VPC, 0)
	input := &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	}
	err := describeVPCs(ctx, client, input, func(awsVPCs []ec2types.Vpc) error {
		for _, vpc := range awsVPCs {
			vpcs = append(vpcs, VPC{
				ID:   aws.ToString(vpc.VpcId),
				CIDR: aws.ToString(vpc.CidrBlock),
				Tags: FromAWSTags(vpc.Tags),
			})
		}
		return nil
	})
	if err != nil {
		return VPC{}, err
	}
	if len(vpcs) == 0 {
		return VPC{}, fmt.Errorf("no vpc found for vpc id %s", vpcID)
	}
	if len(vpcs) > 1 {
		return VPC{}, fmt.Errorf("expected 1 vpc, but found %d: %v", len(vpcs), vpcs)
	}
	return vpcs[0], nil
}

// describeVPCs retrieves metadata for VPCs with given filters.
func describeVPCs(ctx context.Context, client *ec2.Client, input *ec2.DescribeVpcsInput, fn func(vpcs []ec2types.Vpc) error) error {
	vpcOutput, err := client.DescribeVpcs(ctx, input)
	if err != nil {
		return fmt.Errorf("describing vpcs: %w", err)
	}
	return fn(vpcOutput.Vpcs)
}
