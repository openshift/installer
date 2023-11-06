package aws

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	errNotFound = errors.New("not found")

	defaultBackoff = wait.Backoff{
		Steps:    10,
		Duration: 3 * time.Second,
		Factor:   1.0,
		Jitter:   0.1,
	}
)

const (
	errVpcIDNotFound        = "InvalidVpcID.NotFound"
	errNatGatewayIDNotFound = "InvalidNatGatewayID.NotFound"
	errInvalidEIPNotFound   = "InvalidElasticIpID.NotFound"
	errRouteTableIDNotFound = "InvalidRouteTableId.NotFound"
	errInvalidSubnet        = "InvalidSubnet"
	errAuthFailure          = "AuthFailure"
	errIPAddrInUse          = "InvalidIPAddress.InUse"

	quadZeroRoute = "0.0.0.0/0"
)

type vpcInputOptions struct {
	infraID          string
	region           string
	vpcID            string
	cidrV4Block      string
	zones            []string
	edgeZones        []string
	publicSubnetIDs  []string
	privateSubnetIDs []string
	edgeParentMap    map[string]int
	tags             map[string]string
}

type vpcState struct {
	input *vpcInputOptions
	// created resources
	vpcID *string
	igwID *string
}

type vpcOutput struct {
	vpcID            string
	publicSubnetIDs  []string
	privateSubnetIDs []string
	zoneToSubnetMap  map[string]string
}

func createVPCResources(ctx context.Context, logger logrus.FieldLogger, ec2Client ec2iface.EC2API, vpcInput *vpcInputOptions) (*vpcOutput, error) {
	// User-supplied VPC. In this case we don't create any subnets
	if len(vpcInput.vpcID) > 0 {
		logger.WithField("id", vpcInput.vpcID).Infoln("Found user-supplied VPC")

		privateSubnetIDs := aws.StringSlice(vpcInput.privateSubnetIDs)
		privateSubnetZoneMap := make(map[string]*string, len(privateSubnetIDs))
		result, err := ec2Client.DescribeSubnetsWithContext(ctx, &ec2.DescribeSubnetsInput{
			SubnetIds: privateSubnetIDs,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user-supplied subnets: %w", err)
		}
		for _, subnet := range result.Subnets {
			privateSubnetZoneMap[aws.StringValue(subnet.AvailabilityZone)] = subnet.SubnetId
		}

		return &vpcOutput{
			vpcID:            vpcInput.vpcID,
			publicSubnetIDs:  vpcInput.publicSubnetIDs,
			privateSubnetIDs: vpcInput.privateSubnetIDs,
			zoneToSubnetMap:  aws.StringValueMap(privateSubnetZoneMap),
		}, nil
	}

	state := vpcState{input: vpcInput}
	endpointRouteTableIDs := make([]*string, 0, len(vpcInput.zones)+1)

	vpc, err := state.ensureVPC(ctx, logger, ec2Client)
	if err != nil {
		return nil, err
	}
	state.vpcID = vpc.VpcId

	_, err = state.ensureDHCPOptions(ctx, logger, ec2Client, vpc.VpcId)
	if err != nil {
		return nil, err
	}

	igw, err := state.ensureInternetGateway(ctx, logger, ec2Client, vpc.VpcId)
	if err != nil {
		return nil, err
	}
	state.igwID = igw.InternetGatewayId

	publicRouteTable, err := state.ensurePublicRouteTable(ctx, logger, ec2Client)
	if err != nil {
		return nil, err
	}
	endpointRouteTableIDs = append(endpointRouteTableIDs, publicRouteTable.RouteTableId)

	// Per-zone resources
	networks, err := splitNetworks(vpcInput.cidrV4Block, len(vpcInput.zones), len(vpcInput.edgeZones))
	if err != nil {
		return nil, err
	}

	publicSubnetIDs, publicSubnetZoneMap, err := state.ensurePublicSubnets(ctx, logger, ec2Client, networks.standard.public, publicRouteTable)
	if err != nil {
		return nil, err
	}

	privateSubnetIDs, privateSubnetZoneMap, err := state.ensurePrivateSubnets(ctx, logger, ec2Client, networks.standard.private)
	if err != nil {
		return nil, err
	}

	privateRouteTables := make([]*ec2.RouteTable, 0, len(vpcInput.zones))
	for _, zone := range vpcInput.zones {
		var natGwID *string
		if len(publicSubnetZoneMap) > 0 {
			natGw, err := state.ensureNatGateway(ctx, logger, ec2Client, publicSubnetZoneMap[zone], zone)
			if err != nil {
				return nil, fmt.Errorf("failed to create NAT gateway for zone (%s): %w", zone, err)
			}
			natGwID = natGw.NatGatewayId
		}

		privateRouteTable, err := state.ensurePrivateRouteTable(ctx, logger, ec2Client, natGwID, privateSubnetZoneMap[zone], zone)
		if err != nil {
			return nil, fmt.Errorf("failed to create private route table for zone (%s): %w", zone, err)
		}
		endpointRouteTableIDs = append(endpointRouteTableIDs, privateRouteTable.RouteTableId)
		privateRouteTables = append(privateRouteTables, privateRouteTable)
	}

	_, err = state.ensureVPCS3Endpoint(ctx, logger, ec2Client, endpointRouteTableIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC S3 endpoint: %w", err)
	}

	if len(vpcInput.edgeZones) > 0 {
		_, _, err := state.ensureEdgePublicSubnets(ctx, logger, ec2Client, networks.edge.public, publicRouteTable)
		if err != nil {
			return nil, err
		}

		_, subnetZoneMap, err := state.ensureEdgePrivateSubnets(ctx, logger, ec2Client, networks.edge.private)
		if err != nil {
			return nil, err
		}

		for _, zone := range vpcInput.edgeZones {
			// Lookup the index of the parent zone from a given local zone
			// name, getting the index for the route table id for that zone
			// (parent). When not found (parent zone's gateway does not exist),
			// the first route table will be used.
			// Example: edgeParentMap = {"us-east-1-nyc-1a": 0}
			tableIdx, found := vpcInput.edgeParentMap[zone]
			if !found {
				tableIdx = 0
			}
			table := privateRouteTables[tableIdx]
			subnetID := subnetZoneMap[zone]
			if err := addSubnetToRouteTable(ctx, logger, ec2Client, table, subnetID); err != nil {
				return nil, fmt.Errorf("failed to associate edge subnet (%s) to route table (%s): %w", aws.StringValue(subnetID), aws.StringValue(table.RouteTableId), err)
			}
		}
	}

	return &vpcOutput{
		vpcID:            aws.StringValue(state.vpcID),
		privateSubnetIDs: aws.StringValueSlice(privateSubnetIDs),
		publicSubnetIDs:  aws.StringValueSlice(publicSubnetIDs),
		zoneToSubnetMap:  aws.StringValueMap(privateSubnetZoneMap),
	}, nil
}

func (o *vpcState) ensureVPC(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API) (*ec2.Vpc, error) {
	vpcName := fmt.Sprintf("%s-vpc", o.input.infraID)
	l := logger.WithField("vpc", vpcName)

	createdOrFoundMsg := "Found existing VPC"
	vpc, err := existingVPC(ctx, client, ec2Filters(o.input.infraID, vpcName))
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		createdOrFoundMsg = "Created VPC"
		tags := mergeTags(o.input.tags, map[string]string{"Name": vpcName})
		var res *ec2.CreateVpcOutput
		res, err = client.CreateVpcWithContext(ctx,
			&ec2.CreateVpcInput{
				CidrBlock: aws.String(o.input.cidrV4Block),
				TagSpecifications: []*ec2.TagSpecification{
					{
						ResourceType: aws.String("vpc"),
						Tags:         ec2Tags(tags),
					},
				},
			})
		if err != nil {
			return nil, fmt.Errorf("failed to create VPC: %w", err)
		}
		vpc = res.Vpc
	}
	l = l.WithField("id", aws.StringValue(vpc.VpcId))
	l.Infoln(createdOrFoundMsg)

	// Enable DNS support
	_, err = client.ModifyVpcAttributeWithContext(ctx,
		&ec2.ModifyVpcAttributeInput{
			VpcId:            vpc.VpcId,
			EnableDnsSupport: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		})
	if err != nil {
		return nil, fmt.Errorf("failed to enable DNS support on VPC: %w", err)
	}
	l.Infoln("Enabled DNS support on VPC")

	// Enable DNS hostnames
	_, err = client.ModifyVpcAttributeWithContext(ctx,
		&ec2.ModifyVpcAttributeInput{
			VpcId:              vpc.VpcId,
			EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		})
	if err != nil {
		return nil, fmt.Errorf("failed to enable DNS hostnames on VPC: %w", err)
	}
	l.Infoln("Enabled DNS hostnames on VPC")

	return vpc, nil
}

func existingVPC(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.Vpc, error) {
	result, err := client.DescribeVpcsWithContext(ctx, &ec2.DescribeVpcsInput{Filters: filters})
	if err != nil {
		return nil, fmt.Errorf("failed to list VPCs: %w", err)
	}
	for _, vpc := range result.Vpcs {
		return vpc, nil
	}
	return nil, errNotFound
}

func (o *vpcState) ensureDHCPOptions(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, vpcID *string) (*ec2.DhcpOptions, error) {
	filters := ec2Filters(o.input.infraID, "")
	foundOrCreatedMsg := "Found existing DHCP options"
	opt, err := existingDHCPOptions(ctx, client, filters)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		var domainName string
		switch o.input.region {
		case "us-east-1":
			domainName = "ec2.internal"
		default:
			domainName = fmt.Sprintf("%s.compute.internal", o.input.region)
		}
		foundOrCreatedMsg = "Created DHCP options"
		res, err := client.CreateDhcpOptionsWithContext(ctx, &ec2.CreateDhcpOptionsInput{
			DhcpConfigurations: []*ec2.NewDhcpConfiguration{
				{
					Key:    aws.String("domain-name"),
					Values: aws.StringSlice([]string{domainName}),
				},
				{
					Key:    aws.String("domain-name-servers"),
					Values: aws.StringSlice([]string{"AmazonProvidedDNS"}),
				},
			},
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws.String("dhcp-options"),
					Tags:         ec2Tags(o.input.tags),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create DHCP options: %w", err)
		}
		opt = res.DhcpOptions
	}
	l := logger.WithField("id", aws.StringValue(opt.DhcpOptionsId))
	l.Infoln(foundOrCreatedMsg)

	_, err = client.AssociateDhcpOptionsWithContext(ctx, &ec2.AssociateDhcpOptionsInput{
		DhcpOptionsId: opt.DhcpOptionsId,
		VpcId:         vpcID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to associate DHCP options to VPC: %w", err)
	}
	l.WithField("vpcID", aws.StringValue(vpcID)).Infoln("Associated DHCP options with VPC")

	return opt, nil
}

func existingDHCPOptions(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.DhcpOptions, error) {
	res, err := client.DescribeDhcpOptionsWithContext(ctx, &ec2.DescribeDhcpOptionsInput{Filters: filters})
	if err != nil {
		return nil, fmt.Errorf("failed to list DHCP options: %w", err)
	}
	for _, opt := range res.DhcpOptions {
		return opt, nil
	}
	return nil, errNotFound
}

func (o *vpcState) ensureInternetGateway(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, vpcID *string) (*ec2.InternetGateway, error) {
	gatewayName := fmt.Sprintf("%s-igw", o.input.infraID)
	filters := ec2Filters(o.input.infraID, gatewayName)
	foundOrCreatedMsg := "Found existing Internet Gateway"
	igw, err := existingInternetGateway(ctx, client, filters)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		tags := mergeTags(o.input.tags, map[string]string{"Name": gatewayName})
		foundOrCreatedMsg = "Created Internet Gateway"
		res, err := client.CreateInternetGatewayWithContext(ctx, &ec2.CreateInternetGatewayInput{
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws.String("internet-gateway"),
					Tags:         ec2Tags(tags),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create internet gateway: %w", err)
		}
		igw = res.InternetGateway
	}
	l := logger.WithField("id", aws.StringValue(igw.InternetGatewayId))
	l.Infoln(foundOrCreatedMsg)

	l = l.WithField("vpc", aws.StringValue(vpcID))
	attached := false
	for _, attachment := range igw.Attachments {
		if aws.StringValue(attachment.VpcId) == aws.StringValue(vpcID) {
			attached = true
			break
		}
	}
	if !attached {
		_, err := client.AttachInternetGatewayWithContext(ctx, &ec2.AttachInternetGatewayInput{
			InternetGatewayId: igw.InternetGatewayId,
			VpcId:             vpcID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to attach internet gateway to VPC: %w", err)
		}
		l.Infoln("Attached Internet Gateway to VPC")
	} else {
		l.Infoln("Internet Gateway already attached to VPC")
	}

	return igw, nil
}

func existingInternetGateway(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.InternetGateway, error) {
	res, err := client.DescribeInternetGatewaysWithContext(ctx, &ec2.DescribeInternetGatewaysInput{Filters: filters})
	if err != nil {
		return nil, fmt.Errorf("failed to list internet gateways: %w", err)
	}
	for _, igw := range res.InternetGateways {
		return igw, nil
	}
	return nil, errNotFound
}

func (o *vpcState) ensurePublicRouteTable(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API) (*ec2.RouteTable, error) {
	tableName := fmt.Sprintf("%s-public", o.input.infraID)
	l := logger.WithField("route table", tableName)
	routeTable, err := o.ensureRouteTable(ctx, l, client, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to create public route table: %w", err)
	}
	l = l.WithField("id", aws.StringValue(routeTable.RouteTableId))

	// Replace the VPC's main route table
	l.Debugln("Replacing VPC's main route table association")
	filters := []*ec2.Filter{
		{Name: aws.String("vpc-id"), Values: []*string{o.vpcID}},
		{Name: aws.String("association.main"), Values: aws.StringSlice([]string{"true"})},
	}
	rTableInfo, err := existingRouteTable(ctx, client, filters)
	if err != nil {
		if errors.Is(err, errNotFound) {
			return nil, fmt.Errorf("no main route table associated with the VPC")
		}
		return nil, fmt.Errorf("failed to get main route table: %w", err)
	}
	ml := logger.WithField("id", aws.StringValue(rTableInfo.RouteTableId))
	ml.Debugln("Found main route table")

	// Replace route table association only if it's not the associated route table already
	if tID := rTableInfo.RouteTableId; tID != routeTable.RouteTableId {
		var associationID *string
		for _, assoc := range rTableInfo.Associations {
			associationID = assoc.RouteTableAssociationId
			break
		}
		if associationID == nil {
			return nil, fmt.Errorf("no associations found for main route table")
		}
		_, err := client.ReplaceRouteTableAssociationWithContext(ctx, &ec2.ReplaceRouteTableAssociationInput{
			AssociationId: associationID,
			RouteTableId:  routeTable.RouteTableId,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to replace vpc main route table: %w", err)
		}
		ml.Infoln("Associated public route table with VPC's main route table")
	} else {
		ml.Infoln("Public route table is already VPC's main route table association")
	}

	// Create route to Internet Gateway
	l.Debugln("Creating route to Internet Gateway")
	if !hasInternetGatewayRoute(routeTable, o.igwID) {
		_, err := client.CreateRouteWithContext(ctx, &ec2.CreateRouteInput{
			GatewayId:            o.igwID,
			RouteTableId:         routeTable.RouteTableId,
			DestinationCidrBlock: aws.String(quadZeroRoute),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create route to internet gateway: %w", err)
		}
		l.Infoln("Created route to Internet Gateway")
	} else {
		l.Infoln("Found existing route to Internet Gateway")
	}

	return routeTable, nil
}

func (o *vpcState) ensurePrivateRouteTable(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, natGwID *string, subnetID *string, zone string) (*ec2.RouteTable, error) {
	tableName := fmt.Sprintf("%s-private-%s", o.input.infraID, zone)
	l := logger.WithField("route table", tableName)
	table, err := o.ensureRouteTable(ctx, l, client, tableName)
	if err != nil {
		return nil, err
	}

	// Everything bellow is only needed if direct internet access is used
	if natGwID == nil {
		return table, nil
	}

	createdOrFoundMsg := "Found existing route to Nat gateway"
	if !hasNatGatewayRoute(table, natGwID) {
		createdOrFoundMsg = "Created route to Nat gateway"
		if err := createNatGatewayRoute(ctx, client, natGwID, table.RouteTableId); err != nil {
			return nil, fmt.Errorf("failed to create route to nat gateway (%s): %w", aws.StringValue(natGwID), err)
		}
	}
	l.WithField("nat gateway", aws.StringValue(natGwID)).Infoln(createdOrFoundMsg)

	if err = addSubnetToRouteTable(ctx, l, client, table, subnetID); err != nil {
		return nil, fmt.Errorf("failed to associate subnet (%s) to route table (%s): %w", aws.StringValue(subnetID), aws.StringValue(table.RouteTableId), err)
	}

	return table, nil
}

func createNatGatewayRoute(ctx context.Context, client ec2iface.EC2API, natGwID *string, tableID *string) error {
	return wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			_, err = client.CreateRouteWithContext(ctx, &ec2.CreateRouteInput{
				NatGatewayId:         natGwID,
				RouteTableId:         tableID,
				DestinationCidrBlock: aws.String(quadZeroRoute),
			})
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) && (strings.EqualFold(awsErr.Code(), errNatGatewayIDNotFound) || strings.EqualFold(awsErr.Code(), errRouteTableIDNotFound)) {
					return false, nil
				}
				return false, err
			}
			return true, nil
		},
	)
}

func (o *vpcState) ensureRouteTable(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, tableName string) (*ec2.RouteTable, error) {
	filters := ec2Filters(o.input.infraID, tableName)
	createdOrFoundMsg := "Found existing route table"
	routeTable, err := existingRouteTable(ctx, client, filters)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		tags := mergeTags(o.input.tags, map[string]string{"Name": tableName})
		createdOrFoundMsg = "Created route table"
		routeTable, err = createRouteTable(ctx, client, o.vpcID, tags)
		if err != nil {
			return nil, err
		}
	}
	logger.WithField("id", aws.StringValue(routeTable.RouteTableId)).Infoln(createdOrFoundMsg)

	return routeTable, nil
}

// creates a route table and waits until it shows up.
func createRouteTable(ctx context.Context, client ec2iface.EC2API, vpcID *string, tags map[string]string) (*ec2.RouteTable, error) {
	res, err := client.CreateRouteTableWithContext(ctx, &ec2.CreateRouteTableInput{
		VpcId: vpcID,
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("route-table"),
				Tags:         ec2Tags(tags),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	routeTable := res.RouteTable
	err = wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			res, err := client.DescribeRouteTablesWithContext(ctx, &ec2.DescribeRouteTablesInput{
				RouteTableIds: []*string{routeTable.RouteTableId},
			})
			if err != nil || len(res.RouteTables) == 0 {
				return false, nil
			}
			return true, nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to find route table (%s) that was just created: %w", aws.StringValue(routeTable.RouteTableId), err)
	}

	return routeTable, nil
}

func existingRouteTable(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.RouteTable, error) {
	res, err := client.DescribeRouteTablesWithContext(ctx, &ec2.DescribeRouteTablesInput{Filters: filters})
	if err != nil {
		return nil, fmt.Errorf("failed to list route tables: %w", err)
	}
	for _, rt := range res.RouteTables {
		return rt, nil
	}
	return nil, errNotFound
}

func hasInternetGatewayRoute(table *ec2.RouteTable, igwID *string) bool {
	for _, route := range table.Routes {
		if aws.StringValue(route.GatewayId) == aws.StringValue(igwID) && aws.StringValue(route.DestinationCidrBlock) == quadZeroRoute {
			return true
		}
	}
	return false
}

func hasNatGatewayRoute(table *ec2.RouteTable, natID *string) bool {
	for _, route := range table.Routes {
		if aws.StringValue(route.NatGatewayId) == aws.StringValue(natID) && aws.StringValue(route.DestinationCidrBlock) == quadZeroRoute {
			return true
		}
	}
	return false
}

func hasAssociatedSubnet(table *ec2.RouteTable, subnetID *string) bool {
	for _, assoc := range table.Associations {
		if aws.StringValue(assoc.SubnetId) == aws.StringValue(subnetID) {
			return true
		}
	}
	return false
}

func (o *vpcState) ensurePublicSubnets(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, network *net.IPNet, publicTable *ec2.RouteTable) ([]*string, map[string]*string, error) {
	// used to designate it should be used for internet ELBs
	const tagNameSubnetPublicELB = "kubernetes.io/role/elb"
	tags := mergeTags(o.input.tags, map[string]string{
		tagNameSubnetPublicELB: "true",
	})

	ids, zoneMap, err := o.ensureSubnets(ctx, logger, client, network, "public", o.input.zones, tags)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create public subnets: %w", err)
	}

	for _, subnetID := range ids {
		if err = addSubnetToRouteTable(ctx, logger, client, publicTable, subnetID); err != nil {
			return nil, nil, fmt.Errorf("failed to associate public subnet (%s) with public route table (%s): %w", aws.StringValue(subnetID), aws.StringValue(publicTable.RouteTableId), err)
		}
	}

	return ids, zoneMap, nil
}

func (o *vpcState) ensurePrivateSubnets(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, network *net.IPNet) ([]*string, map[string]*string, error) {
	// tag name used to designate the subnet should be used for internal ELBs
	const tagNameSubnetInternalELB = "kubernetes.io/role/internal-elb"
	tags := mergeTags(o.input.tags, map[string]string{
		tagNameSubnetInternalELB: "true",
	})

	ids, zoneMap, err := o.ensureSubnets(ctx, logger, client, network, "private", o.input.zones, tags)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create private subnets: %w", err)
	}
	return ids, zoneMap, nil
}

func (o *vpcState) ensureEdgePublicSubnets(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, network *net.IPNet, publicTable *ec2.RouteTable) ([]*string, map[string]*string, error) {
	ids, zoneMap, err := o.ensureSubnets(ctx, logger, client, network, "public", o.input.edgeZones, o.input.tags)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create public edge subnets: %w", err)
	}

	for _, subnetID := range ids {
		if err = addSubnetToRouteTable(ctx, logger, client, publicTable, subnetID); err != nil {
			return nil, nil, fmt.Errorf("failed to associate public edge subnet (%s) with public route table (%s): %w", aws.StringValue(subnetID), aws.StringValue(publicTable.RouteTableId), err)
		}
	}

	return ids, zoneMap, nil
}

func (o *vpcState) ensureEdgePrivateSubnets(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, network *net.IPNet) ([]*string, map[string]*string, error) {
	ids, zoneMap, err := o.ensureSubnets(ctx, logger, client, network, "private", o.input.edgeZones, o.input.tags)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create private edge subnets: %w", err)
	}
	return ids, zoneMap, nil
}

func (o *vpcState) ensureSubnets(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, network *net.IPNet, subnetType string, zones []string, tags map[string]string) ([]*string, map[string]*string, error) {
	subnetIDs := make([]*string, 0, len(zones))
	subnetZoneMap := make(map[string]*string, len(zones))
	newBits := int(math.Ceil(math.Log2(float64(len(zones)))))
	for i, zone := range zones {
		cidr, err := cidr.Subnet(network, newBits, i)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to subnet %s network: %w", subnetType, err)
		}

		name := fmt.Sprintf("%s-%s-%s", o.input.infraID, subnetType, zone)
		subnet, err := o.ensureSubnet(ctx, logger, client, zone, cidr.String(), name, tags)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create %s subnet (%s): %w", subnetType, name, err)
		}
		subnetIDs = append(subnetIDs, subnet.SubnetId)
		subnetZoneMap[zone] = subnet.SubnetId
	}

	return subnetIDs, subnetZoneMap, nil
}

func (o *vpcState) ensureSubnet(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, zone string, cidr string, name string, tags map[string]string) (*ec2.Subnet, error) {
	l := logger.WithField("subnet", name)
	filters := ec2Filters(o.input.infraID, name)
	createdOrFoundMsg := "Found existing subnet"
	subnet, err := existingSubnet(ctx, client, filters)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		tags := mergeTags(tags, map[string]string{"Name": name})
		createdOrFoundMsg = "Created subnet"
		res, err := client.CreateSubnetWithContext(ctx, &ec2.CreateSubnetInput{
			AvailabilityZone: aws.String(zone),
			VpcId:            o.vpcID,
			CidrBlock:        aws.String(cidr),
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws.String("subnet"),
					Tags:         ec2Tags(tags),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create subnet: %w", err)
		}
		subnet = res.Subnet
	}
	l.Infoln(createdOrFoundMsg)

	err = wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			res, err := client.DescribeSubnetsWithContext(ctx, &ec2.DescribeSubnetsInput{
				SubnetIds: []*string{subnet.SubnetId},
			})
			if err != nil || len(res.Subnets) == 0 {
				return false, nil
			}
			return true, nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet (%s) that was just created: %w", aws.StringValue(subnet.SubnetId), err)
	}

	return subnet, nil
}

func existingSubnet(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.Subnet, error) {
	res, err := client.DescribeSubnetsWithContext(ctx, &ec2.DescribeSubnetsInput{Filters: filters})
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets: %w", err)
	}
	for _, subnet := range res.Subnets {
		return subnet, nil
	}
	return nil, errNotFound
}

func (o *vpcState) ensureNatGateway(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, subnetID *string, zone string) (*ec2.NatGateway, error) {
	natName := fmt.Sprintf("%s-nat-%s", o.input.infraID, zone)
	l := logger.WithField("nat gateway", natName)
	filters := ec2Filters(o.input.infraID, natName)

	createdOrFoundMsg := "Found existing Nat gateway"
	natGW, err := existingNatGateway(ctx, client, filters)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}

		// First allocate an EIP address
		eipName := fmt.Sprintf("%s-eip-%s", o.input.infraID, zone)
		eipTags := mergeTags(o.input.tags, map[string]string{"Name": eipName})
		allocID, err := ensureEIP(ctx, client, eipTags)
		if err != nil {
			return nil, err
		}

		tags := mergeTags(o.input.tags, map[string]string{"Name": natName})
		createdOrFoundMsg = "Created Nat gateway"
		natGW, err = createNatGateway(ctx, client, allocID, subnetID, tags)
		if err != nil {
			return nil, fmt.Errorf("failed to create nat gateway: %w", err)
		}
	}
	l.WithField("id", aws.StringValue(natGW.NatGatewayId)).Infoln(createdOrFoundMsg)

	return natGW, nil
}

func ensureEIP(ctx context.Context, client ec2iface.EC2API, tags map[string]string) (*string, error) {
	res, err := client.AllocateAddressWithContext(ctx, &ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to allocate EIP: %w", err)
	}

	isErrorRetriable := func(err error) bool {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			return strings.EqualFold(awsErr.Code(), errInvalidEIPNotFound)
		}
		return false
	}

	ec2Tags := ec2Tags(tags)
	// NOTE: there is a potential to leak EIP addresses if the following tag
	// and release operations fail, since we have no way of recognizing the EIP
	// as belonging to the cluster
	err = wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			_, err = client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
				Resources: []*string{res.AllocationId},
				Tags:      ec2Tags,
			})
			if err != nil {
				if isErrorRetriable(err) {
					return false, nil
				}
				return true, err
			}
			return true, nil
		},
	)
	if err != nil {
		if rerr := releaseEIP(ctx, client, res.AllocationId); rerr != nil {
			return nil, rerr
		}
		return nil, fmt.Errorf("failed to tag EIP (%s): %w", aws.StringValue(res.AllocationId), err)
	}

	return res.AllocationId, nil
}

func releaseEIP(ctx context.Context, client ec2iface.EC2API, allocID *string) error {
	err := wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			_, err = client.DisassociateAddressWithContext(ctx, &ec2.DisassociateAddressInput{
				AssociationId: allocID,
			})
			if err != nil {
				return false, nil
			}

			_, err = client.ReleaseAddressWithContext(ctx, &ec2.ReleaseAddressInput{
				AllocationId: allocID,
			})
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) {
					// IP address already released: see https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_ReleaseAddress.html
					if strings.EqualFold(awsErr.Code(), errAuthFailure) {
						return true, nil
					}
					// Must be disassociated first
					if strings.EqualFold(awsErr.Code(), errIPAddrInUse) {
						return false, nil
					}
				}
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to release untagged EIP (%s): %w", aws.StringValue(allocID), err)
	}
	return nil
}

func createNatGateway(ctx context.Context, client ec2iface.EC2API, allocID *string, subnetID *string, tags map[string]string) (*ec2.NatGateway, error) {
	isErrorRetriable := func(err error) bool {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			return strings.EqualFold(awsErr.Code(), errInvalidSubnet) || strings.EqualFold(awsErr.Code(), errInvalidEIPNotFound)
		}
		return false
	}
	ec2Tags := ec2Tags(tags)
	var natGW *ec2.NatGateway
	err := wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			res, err := client.CreateNatGatewayWithContext(ctx, &ec2.CreateNatGatewayInput{
				AllocationId: allocID,
				SubnetId:     subnetID,
				TagSpecifications: []*ec2.TagSpecification{
					{
						ResourceType: aws.String("natgateway"),
						Tags:         ec2Tags,
					},
				},
			})
			if err != nil {
				if isErrorRetriable(err) {
					return false, nil
				}
				return true, err
			}
			natGW = res.NatGateway
			return true, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return natGW, nil
}

func existingNatGateway(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.NatGateway, error) {
	res, err := client.DescribeNatGatewaysWithContext(ctx, &ec2.DescribeNatGatewaysInput{Filter: filters})
	if err != nil {
		return nil, fmt.Errorf("failed to list Nat gateways: %w", err)
	}
	for _, nat := range res.NatGateways {
		state := aws.StringValue(nat.State)
		if state == "deleted" || state == "deleting" || state == "failed" {
			continue
		}
		return nat, nil
	}
	return nil, errNotFound
}

func (o *vpcState) ensureVPCS3Endpoint(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, routeTableIDs []*string) (*ec2.VpcEndpoint, error) {
	filters := ec2Filters(o.input.infraID, "")
	createdOrFoundMsg := "Found existing VPC S3 endpoint"
	endpoint, err := existingVPCEndpoint(ctx, client, filters)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		createdOrFoundMsg = "Created VPC S3 endpoint"
		serviceName := fmt.Sprintf("com.amazonaws.%s.s3", o.input.region)
		endpoint, err = createVPCEndpoint(ctx, client, serviceName, o.vpcID, routeTableIDs, o.input.tags)
		if err != nil {
			return nil, err
		}
	}
	logger.WithField("id", aws.StringValue(endpoint.VpcEndpointId)).Infoln(createdOrFoundMsg)

	return endpoint, nil
}

func createVPCEndpoint(ctx context.Context, client ec2iface.EC2API, serviceName string, vpcID *string, routeTableIDs []*string, tags map[string]string) (*ec2.VpcEndpoint, error) {
	ec2Tags := ec2Tags(tags)
	var vpcEndpoint *ec2.VpcEndpoint
	err := wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			res, err := client.CreateVpcEndpointWithContext(ctx, &ec2.CreateVpcEndpointInput{
				VpcId:         vpcID,
				ServiceName:   aws.String(serviceName),
				RouteTableIds: routeTableIDs,
				TagSpecifications: []*ec2.TagSpecification{
					{
						ResourceType: aws.String("vpc-endpoint"),
						Tags:         ec2Tags,
					},
				},
			})
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) {
					if strings.EqualFold(awsErr.Code(), errRouteTableIDNotFound) {
						return false, nil
					}
				}
				return true, err
			}
			vpcEndpoint = res.VpcEndpoint
			return true, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC endpoint: %w", err)
	}

	return vpcEndpoint, nil
}

func existingVPCEndpoint(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.VpcEndpoint, error) {
	res, err := client.DescribeVpcEndpointsWithContext(ctx, &ec2.DescribeVpcEndpointsInput{Filters: filters})
	if err != nil {
		return nil, fmt.Errorf("failed to list VPC endpoints: %w", err)
	}
	for _, endpoint := range res.VpcEndpoints {
		return endpoint, nil
	}
	return nil, errNotFound
}

type splitOutput struct {
	standard struct {
		private *net.IPNet
		public  *net.IPNet
	}
	edge struct {
		private *net.IPNet
		public  *net.IPNet
	}
}

func splitNetworks(cidrBlock string, numZones int, numEdgeZones int) (*splitOutput, error) {
	output := &splitOutput{}

	_, network, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IPv4 CIDR blocks: %w", err)
	}

	// CIDR blocks for default IPI installation
	privateNetwork, err := cidr.Subnet(network, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to determine private subnet: %w", err)
	}
	publicNetwork, err := cidr.Subnet(network, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to determine public subnet: %w", err)
	}

	// CIDR blocks used when creating subnets into edge zones.
	var edgePrivateNetwork *net.IPNet
	var edgePublicNetwork *net.IPNet
	if numEdgeZones > 0 {
		// The public CIDR is used to create the CIDR blocks for edge subnets.
		sharedPublicNetwork, err := cidr.Subnet(publicNetwork, 1, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to determine shared public subnet: %w", err)
		}
		sharedEdgeNetwork, err := cidr.Subnet(publicNetwork, 1, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to determine shared edge subnet: %w", err)
		}
		publicNetwork = sharedPublicNetwork

		// CIDR bloks for edge subnets
		edgePrivateNetwork, err = cidr.Subnet(sharedEdgeNetwork, 1, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to determine edge private subnet: %w", err)
		}
		edgePublicNetwork, err = cidr.Subnet(sharedEdgeNetwork, 1, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to determine edge public subnet: %w", err)
		}
	}

	// If a single-zone deployment, the available CIDR block is split into two
	// to allow for user expansion
	if numZones == 1 {
		privateNetwork, err = cidr.Subnet(privateNetwork, 1, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to split private subnet in single-zone deployment: %w", err)
		}
		publicNetwork, err = cidr.Subnet(publicNetwork, 1, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to split public subnet in single-zone deployment: %w", err)
		}
	}
	output.standard.private = privateNetwork
	output.standard.public = publicNetwork

	// If a single-zone deployment, the available CIDR block is split into two
	// to allow for user expansion
	if numEdgeZones == 1 {
		edgePrivateNetwork, err = cidr.Subnet(edgePrivateNetwork, 1, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to split private edge subnet in single-zone deployment: %w", err)
		}
		edgePublicNetwork, err = cidr.Subnet(edgePublicNetwork, 1, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to split public edge subnet in single-zone deployment: %w", err)
		}
	}
	output.edge.private = edgePrivateNetwork
	output.edge.public = edgePublicNetwork

	return output, nil
}

func addSubnetToRouteTable(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, table *ec2.RouteTable, subnetID *string) error {
	l := logger.WithFields(logrus.Fields{
		"subnet id": aws.StringValue(subnetID),
		"table id":  aws.StringValue(table.RouteTableId),
	})
	if hasAssociatedSubnet(table, subnetID) {
		l.Infoln("Subnet already associated with route table")
		return nil
	}
	err := wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (bool, error) {
			_, err := client.AssociateRouteTableWithContext(ctx, &ec2.AssociateRouteTableInput{
				RouteTableId: table.RouteTableId,
				SubnetId:     subnetID,
			})
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) && strings.EqualFold(awsErr.Code(), errRouteTableIDNotFound) {
					return false, nil
				}
				return false, err
			}
			return true, nil
		},
	)
	if err != nil {
		return err
	}
	l.Infoln("Associated subnet with route table")
	return nil
}

func ec2Filters(infraID, name string) []*ec2.Filter {
	filters := []*ec2.Filter{
		{
			Name:   aws.String(fmt.Sprintf("tag:%s", clusterOwnedTag(infraID))),
			Values: aws.StringSlice([]string{ownedTagValue}),
		},
	}
	if len(name) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: aws.StringSlice([]string{name}),
		})
	}
	return filters
}

func ec2Tags(tags map[string]string) []*ec2.Tag {
	ec2Tags := make([]*ec2.Tag, 0, len(tags))
	for k, v := range tags {
		k, v := k, v // needed because we use the addresses
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return ec2Tags
}
