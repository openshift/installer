package cloudforms

import (
	"fmt"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
)

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
		return nil, maybeAwsErr(err)
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

// getHostedZoneName returns the domain name of a route53 hosted zone specified
// by its ID, without the final period.
func getHostedZoneName(sess *session.Session, hostedZoneID string) (string, error) {
	hostedZone, err := route53.New(sess).GetHostedZone(
		&route53.GetHostedZoneInput{
			Id: aws.String(hostedZoneID),
		},
	)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(aws.StringValue(hostedZone.HostedZone.Name), "."), nil
}
