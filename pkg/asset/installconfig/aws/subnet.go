package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Subnet holds metadata for a subnet.
type Subnet struct {
	// ARN is the subnet's Amazon Resource Name.
	ARN string

	// Zone is the subnet's availability zone.
	Zone string

	// CIDR is the subnet's CIDR block.
	CIDR string
}

// subnets retrieves metadata for the given subnet(s).
func subnets(ctx context.Context, session *session.Session, region string, ids []string) (vpc string, private map[string]Subnet, public map[string]Subnet, err error) {
	metas := make(map[string]Subnet, len(ids))
	private = map[string]Subnet{}
	public = map[string]Subnet{}
	var vpcFromSubnet string
	client := ec2.New(session, aws.NewConfig().WithRegion(region))

	idPointers := make([]*string, len(ids))
	for _, id := range ids {
		idPointers = append(idPointers, aws.String(id))
	}
	results, err := client.DescribeSubnetsWithContext( // FIXME: port to DescribeSubnetsPagesWithContext once we bump our vendored AWS package past v1.19.30
		ctx,
		&ec2.DescribeSubnetsInput{SubnetIds: idPointers},
	)
	if err != nil {
		return vpc, nil, nil, errors.Wrap(err, "describing subnets")
	}
	for _, subnet := range results.Subnets {
		if subnet.SubnetId == nil {
			continue
		}
		if subnet.SubnetArn == nil {
			return vpc, nil, nil, errors.Errorf("%s has no ARN", *subnet.SubnetId)
		}
		if subnet.VpcId == nil {
			return vpc, nil, nil, errors.Errorf("%s has no VPC", *subnet.SubnetId)
		}
		if subnet.AvailabilityZone == nil {
			return vpc, nil, nil, errors.Errorf("%s has no availability zone", *subnet.SubnetId)
		}

		if vpc == "" {
			vpc = *subnet.VpcId
			vpcFromSubnet = *subnet.SubnetId
		} else if *subnet.VpcId != vpc {
			return vpc, nil, nil, errors.Errorf("all subnets must belong to the same VPC: %s is from %s, but %s is from %s", *subnet.SubnetId, *subnet.VpcId, vpcFromSubnet, vpc)
		}

		metas[*subnet.SubnetId] = Subnet{
			ARN:  *subnet.SubnetArn,
			Zone: *subnet.AvailabilityZone,
			CIDR: *subnet.CidrBlock,
		}
	}

	var routeTables []*ec2.RouteTable
	err = client.DescribeRouteTablesPagesWithContext(
		ctx,
		&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpc)},
			}},
		},
		func(results *ec2.DescribeRouteTablesOutput, lastPage bool) bool {
			routeTables = append(routeTables, results.RouteTables...)
			return !lastPage
		},
	)
	if err != nil {
		return vpc, nil, nil, errors.Wrap(err, "describing route tables")
	}

	for _, id := range ids {
		meta, ok := metas[id]
		if !ok {
			return vpc, nil, nil, errors.Errorf("failed to find %s", id)
		}
		isPublic, err := isSubnetPublic(routeTables, id)
		if err != nil {
			return vpc, nil, nil, err
		}
		if isPublic {
			public[id] = meta
		} else {
			private[id] = meta
		}
	}

	return vpc, private, public, nil
}

// https://github.com/kubernetes/kubernetes/blob/9f036cd43d35a9c41d7ac4ca82398a6d0bef957b/staging/src/k8s.io/legacy-cloud-providers/aws/aws.go#L3376-L3419
func isSubnetPublic(rt []*ec2.RouteTable, subnetID string) (bool, error) {
	var subnetTable *ec2.RouteTable
	for _, table := range rt {
		for _, assoc := range table.Associations {
			if aws.StringValue(assoc.SubnetId) == subnetID {
				subnetTable = table
				break
			}
		}
	}

	if subnetTable == nil {
		// If there is no explicit association, the subnet will be implicitly
		// associated with the VPC's main routing table.
		for _, table := range rt {
			for _, assoc := range table.Associations {
				if aws.BoolValue(assoc.Main) {
					logrus.Debugf("Assuming implicit use of main routing table %s for %s",
						aws.StringValue(table.RouteTableId), subnetID)
					subnetTable = table
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
		if strings.HasPrefix(aws.StringValue(route.GatewayId), "igw") {
			return true, nil
		}
	}

	return false, nil
}
