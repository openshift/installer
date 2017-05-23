package aws

import (
	"fmt"
	"net"

	"strings"

	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/coreos/ipnets"
)

type VPCSubnet struct {
	// Identifier of the subnet if already existing
	ID string `json:"id"`
	// Logical name for this subnet
	// ignored if existing
	Name string `json:"name"`
	// Availability zone for this subnet
	// Max one subnet per availability zone
	AvailabilityZone string `json:"availabilityZone"`
	// CIDR for this subnet
	// must be disjoint from other subnets
	// must be contained by VPC CIDR
	InstanceCIDR string `json:"instanceCIDR"`
}

// GetVPCSubnets returns the lists of existing subnets in the given VPC, that
// are suitable for controllers and workers nodes.
func GetVPCSubnets(sess *session.Session, vpcID string) ([]VPCSubnet, []VPCSubnet, error) {
	ec2Svc := ec2.New(sess)

	// Get existing route tables for the given VPC.
	routeTables, err := ec2Svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	var publicSubnetIDs, privateSubnetIDs []*string
	for _, routeTable := range routeTables.RouteTables {
		// By looking at the default route in this route table, determine whether
		// its associated subnets are public (using an Internet Gateway) or private
		// (using a some other resource e.g. a NAT gateway, NAT instance, etc).
		//
		// Because there can be only one default route, the subnet can't be both
		// public and private at the same time.
		var defaultRoute *ec2.Route
		for _, route := range routeTable.Routes {
			if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == "0.0.0.0/0" {
				defaultRoute = route
				break
			}
		}
		if defaultRoute == nil {
			// No default route, skip the table.
			continue
		}

		isPublic := false
		if gwID := aws.StringValue(defaultRoute.GatewayId); strings.HasPrefix(gwID, "igw-") {
			// If the route table has an Internet gateway as its default route,
			// then assume it corresponds to a public subnet. If its default
			// route points to some other resource, then assume it corresponds
			// to a private subnet.
			isPublic = true
		}

		// Get the associated subnet IDs.
		for _, association := range routeTable.Associations {
			if aws.StringValue(association.SubnetId) == "" {
				continue
			}
			if isPublic {
				publicSubnetIDs = append(publicSubnetIDs, association.SubnetId)
			} else {
				privateSubnetIDs = append(privateSubnetIDs, association.SubnetId)
			}
		}
	}

	// Retrieve the actual VPCSubnet objects that are suitable for controllers
	// using the subnet IDs we collected.
	publicSubnets, err := getVPCSubnetsByIDs(sess, publicSubnetIDs)
	if err != nil {
		return nil, nil, err
	}

	// Retrieve the actual VPCSubnet objects that are suitable for workers
	// using the subnet IDs we collected.
	privateSubnets, err := getVPCSubnetsByIDs(sess, privateSubnetIDs)
	if err != nil {
		return nil, nil, err
	}

	return publicSubnets, privateSubnets, nil
}

// GetDefaultSubnets partitions a CIDR into subnets
func GetDefaultSubnets(sess *session.Session, vpcCIDR string) ([]VPCSubnet, []VPCSubnet, error) {
	zones, err := getAvailabilityZones(sess)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting availability zone %v", err)
	}

	_, vpcNet, err := net.ParseCIDR(vpcCIDR)
	if vpcNet == nil || err != nil {
		return nil, nil, fmt.Errorf("failed parsing VPC CIDR %v", err)
	}

	// Calculate subnetMultipler many times as many subnets as needed to
	// intentionally leave unused IPs for unspecified use. A multipler
	// of 1 divides the VPC among AZs, 2 leaves 50% of the VPC unallocated,
	// 4 leaves 75% unallocated, etc.
	cidrs, err := ipnets.SubnetInto(vpcNet, 2*2*len(zones))
	if err != nil {
		return nil, nil, fmt.Errorf("failed dividing VPC into subnets %v", err)
	}

	controllerSubnets := make([]VPCSubnet, len(zones))
	workerSubnets := make([]VPCSubnet, len(zones))

	// add generated multi-AZ subnets for controllers
	for i, zone := range zones {
		controllerSubnets[i] = VPCSubnet{
			AvailabilityZone: zone,
			InstanceCIDR:     cidrs[i].String(),
		}
	}

	// add generated multi-AZ subnets for workers
	for i, zone := range zones {
		workerSubnets[i] = VPCSubnet{
			AvailabilityZone: zone,
			InstanceCIDR:     cidrs[i+len(zones)].String(),
		}
	}

	return controllerSubnets, workerSubnets, nil
}

// GetHostedZoneRecords returns the records
func GetHostedZoneRecords(sess *session.Session, zoneID, recordName, recordType string, limit int) ([]*route53.ResourceRecordSet, error) {
	route53svc := route53.New(sess)

	r53input := route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		StartRecordName: aws.String(recordName),
		StartRecordType: aws.String(recordType),
		MaxItems:        aws.String(strconv.Itoa(limit)),
	}

	response, err := route53svc.ListResourceRecordSets(&r53input)
	if err != nil {
		return nil, err
	}

	return response.ResourceRecordSets, nil
}

// getAvailabilityZones lists zones in the region set in the client
// config.
func getAvailabilityZones(sess *session.Session) ([]string, error) {
	// create an EC2 service client
	ec2Svc := ec2.New(sess)

	// query for availability zone catalog
	output, err := ec2Svc.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return nil, err
	}

	// unpack DescribeAvailabilityZonesOutput
	zones := make([]string, len(output.AvailabilityZones))
	for i, zone := range output.AvailabilityZones {
		zones[i] = aws.StringValue(zone.ZoneName)
	}

	return zones, nil
}

// getInternetGateway returns the first available InternetGateway in the VPC.
func getInternetGateway(sess *session.Session, vpcID string) (*ec2.InternetGateway, error) {
	ec2Svc := ec2.New(sess)

	output, err := ec2Svc.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(output.InternetGateways) == 0 {
		return nil, fmt.Errorf("failed to find an Internet Gateway in VPC %s", vpcID)
	}

	gateway := output.InternetGateways[0]
	if len(gateway.Attachments) >= 1 {
		if aws.StringValue(gateway.Attachments[0].State) == "available" {
			return gateway, nil
		} else {
			return nil, fmt.Errorf("internet gateway %s is not 'available'", gateway.GoString())
		}
	}
	return nil, fmt.Errorf("internet gateway %s, has no attachments", gateway.GoString())
}

// getVPC returns the VPC corresponding to the given ID, or an error if it
// doesn't exist.
func getVPC(sess *session.Session, vpcID string) (*ec2.Vpc, error) {
	ec2Svc := ec2.New(sess)

	vpcOutput, err := ec2Svc.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(vpcID)},
	})
	if err != nil {
		return nil, fmt.Errorf("could not describe vpc %s: %s", vpcID, err)
	}
	if len(vpcOutput.Vpcs) == 0 {
		return nil, fmt.Errorf("could not find vpc %s in region", vpcID)
	}
	if len(vpcOutput.Vpcs) > 1 {
		// Theoretically this should never happen. If it does, we probably want to know.
		return nil, fmt.Errorf("found more than one vpc with id %s. this is NOT NORMAL", vpcID)
	}
	return vpcOutput.Vpcs[0], nil
}

// getVPCSubnetsByIDs returns the VPCSubnet objects for the given subnet IDs.
func getVPCSubnetsByIDs(sess *session.Session, subnetIDs []*string) ([]VPCSubnet, error) {
	ec2Svc := ec2.New(sess)

	output, err := ec2Svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("subnet-id"),
				Values: subnetIDs,
			},
		},
	})
	if err != nil {
		return []VPCSubnet{}, err
	}

	subnets := make([]VPCSubnet, 0, len(output.Subnets))
	for _, subnet := range output.Subnets {
		subnets = append(subnets, VPCSubnet{
			ID:               aws.StringValue(subnet.SubnetId),
			AvailabilityZone: aws.StringValue(subnet.AvailabilityZone),
			InstanceCIDR:     aws.StringValue(subnet.CidrBlock),
		})
	}
	return subnets, nil
}
