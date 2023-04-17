package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

// findEC2Instances returns the EC2 instances with tags that satisfy the filters.
// returns two lists, first one is the list of all resources that are not terminated and are not in shutdown
// stage and the second list is the list of resources that are not terminated.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func findEC2Instances(ctx context.Context, ec2Client *ec2.EC2, deleted sets.String, filters []Filter, logger logrus.FieldLogger) ([]string, []string, error) {
	if ec2Client.Config.Region == nil {
		return nil, nil, errors.New("EC2 client does not have region configured")
	}

	partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), *ec2Client.Config.Region)
	if !ok {
		return nil, nil, errors.Errorf("no partition found for region %q", *ec2Client.Config.Region)
	}

	var resourcesRunning []string
	var resourcesNotTerminated []string
	for _, filter := range filters {
		logger.Debugf("search for instances by tag matching %#+v", filter)
		instanceFilters := make([]*ec2.Filter, 0, len(filter))
		for key, value := range filter {
			instanceFilters = append(instanceFilters, &ec2.Filter{
				Name:   aws.String("tag:" + key),
				Values: []*string{aws.String(value)},
			})
		}
		err := ec2Client.DescribeInstancesPagesWithContext(
			ctx,
			&ec2.DescribeInstancesInput{Filters: instanceFilters},
			func(results *ec2.DescribeInstancesOutput, lastPage bool) bool {
				for _, reservation := range results.Reservations {
					if reservation.OwnerId == nil {
						continue
					}

					for _, instance := range reservation.Instances {
						if instance.InstanceId == nil || instance.State == nil {
							continue
						}

						instanceLogger := logger.WithField("instance", *instance.InstanceId)
						arn := fmt.Sprintf("arn:%s:ec2:%s:%s:instance/%s", partition.ID(), *ec2Client.Config.Region, *reservation.OwnerId, *instance.InstanceId)
						if *instance.State.Name == "terminated" {
							if !deleted.Has(arn) {
								instanceLogger.Info("Terminated")
								deleted.Insert(arn)
							}
							continue
						}
						if *instance.State.Name != "shutting-down" {
							resourcesRunning = append(resourcesRunning, arn)
						}
						resourcesNotTerminated = append(resourcesNotTerminated, arn)
					}
				}
				return !lastPage
			},
		)
		if err != nil {
			err = errors.Wrap(err, "get ec2 instances")
			logger.Info(err)
			return resourcesRunning, resourcesNotTerminated, err
		}
	}
	return resourcesRunning, resourcesNotTerminated, nil
}

func deleteEC2(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	client := ec2.New(session)

	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id).WithField("resourceType", resourceType)

	switch resourceType {
	case "dhcp-options":
		return deleteEC2DHCPOptions(ctx, client, id, logger)
	case "elastic-ip":
		return deleteEC2ElasticIP(ctx, client, id, logger)
	case "image":
		return deleteEC2Image(ctx, client, id, logger)
	case "instance":
		return terminateEC2Instance(ctx, client, iam.New(session), id, logger)
	case "internet-gateway":
		return deleteEC2InternetGateway(ctx, client, id, logger)
	case "natgateway":
		return deleteEC2NATGateway(ctx, client, id, logger)
	case "placement-group":
		return deleteEC2PlacementGroup(ctx, client, id, logger)
	case "route-table":
		return deleteEC2RouteTable(ctx, client, id, logger)
	case "security-group":
		return deleteEC2SecurityGroup(ctx, client, id, logger)
	case "snapshot":
		return deleteEC2Snapshot(ctx, client, id, logger)
	case "network-interface":
		return deleteEC2NetworkInterface(ctx, client, id, logger)
	case "subnet":
		return deleteEC2Subnet(ctx, client, id, logger)
	case "volume":
		return deleteEC2Volume(ctx, client, id, logger)
	case "vpc":
		return deleteEC2VPC(ctx, client, elb.New(session), elbv2.New(session), id, logger)
	case "vpc-endpoint":
		return deleteEC2VPCEndpoint(ctx, client, id, logger)
	case "vpc-peering-connection":
		return deleteEC2VPCPeeringConnection(ctx, client, id, logger)
	case "vpc-endpoint-service":
		return deleteEC2VPCEndpointService(ctx, client, id, logger)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func deleteEC2DHCPOptions(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteDhcpOptionsWithContext(ctx, &ec2.DeleteDhcpOptionsInput{
		DhcpOptionsId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidDhcpOptionsID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2Image(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	// tag the snapshots used by the AMI so that the snapshots are matched
	// by the filter and deleted
	response, err := client.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		ImageIds: []*string{&id},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAMIID.NotFound" {
			return nil
		}
		return err
	}
	for _, image := range response.Images {
		var snapshots []*string
		for _, bdm := range image.BlockDeviceMappings {
			if bdm.Ebs != nil && bdm.Ebs.SnapshotId != nil {
				snapshots = append(snapshots, bdm.Ebs.SnapshotId)
			}
		}
		if len(snapshots) != 0 {
			_, err = client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
				Resources: snapshots,
				Tags:      image.Tags,
			})
			if err != nil {
				return errors.Wrapf(err, "tagging snapshots for %s", id)
			}
		}
	}

	_, err = client.DeregisterImageWithContext(ctx, &ec2.DeregisterImageInput{
		ImageId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAMIID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2ElasticIP(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.ReleaseAddressWithContext(ctx, &ec2.ReleaseAddressInput{
		AllocationId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAllocationID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Released")
	return nil
}

func terminateEC2Instance(ctx context.Context, ec2Client *ec2.EC2, iamClient *iam.IAM, id string, logger logrus.FieldLogger) error {
	response, err := ec2Client.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidInstanceID.NotFound" {
			return nil
		}
		return err
	}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			err = terminateEC2InstanceByInstance(ctx, ec2Client, iamClient, instance, logger)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func terminateEC2InstanceByInstance(ctx context.Context, ec2Client *ec2.EC2, iamClient *iam.IAM, instance *ec2.Instance, logger logrus.FieldLogger) error {
	// Ignore instances that are already terminated
	if instance.State == nil || *instance.State.Name == "terminated" {
		return nil
	}

	if instance.IamInstanceProfile != nil {
		parsed, err := arn.Parse(*instance.IamInstanceProfile.Arn)
		if err != nil {
			return errors.Wrap(err, "parse ARN for IAM instance profile")
		}

		err = deleteIAMInstanceProfile(ctx, iamClient, parsed, logger)
		if err != nil {
			return errors.Wrapf(err, "deleting %s", parsed.String())
		}
	}

	_, err := ec2Client.TerminateInstancesWithContext(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []*string{instance.InstanceId},
	})
	if err != nil {
		return err
	}

	logger.Debug("Terminating")
	return nil
}

func deleteEC2InternetGateway(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeInternetGatewaysWithContext(ctx, &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{aws.String(id)},
	})
	if err != nil {
		return err
	}

	for _, gateway := range response.InternetGateways {
		for _, vpc := range gateway.Attachments {
			if vpc.VpcId == nil {
				logger.Warn("gateway does not have a VPC ID")
				continue
			}
			_, err := client.DetachInternetGatewayWithContext(ctx, &ec2.DetachInternetGatewayInput{
				InternetGatewayId: gateway.InternetGatewayId,
				VpcId:             vpc.VpcId,
			})
			if err == nil {
				logger.WithField("vpc", *vpc.VpcId).Debug("Detached")
			} else if err.(awserr.Error).Code() != "Gateway.NotAttached" {
				return errors.Wrapf(err, "detaching from %s", *vpc.VpcId)
			}
		}
	}

	_, err = client.DeleteInternetGatewayWithContext(ctx, &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: &id,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NATGateway(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteNatGatewayWithContext(ctx, &ec2.DeleteNatGatewayInput{
		NatGatewayId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "NatGatewayNotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NATGatewaysByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeNatGatewaysPagesWithContext(
		ctx,
		&ec2.DescribeNatGatewaysInput{
			Filter: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeNatGatewaysOutput, lastPage bool) bool {
			for _, gateway := range results.NatGateways {
				err := deleteEC2NATGateway(ctx, client, *gateway.NatGatewayId, logger.WithField("NAT gateway", *gateway.NatGatewayId))
				if err != nil {
					if lastError != nil {
						logger.Debug(err)
					}
					lastError = errors.Wrapf(err, "deleting EC2 NAT gateway %s", *gateway.NatGatewayId)
					if failFast {
						return false
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2PlacementGroup(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribePlacementGroupsWithContext(ctx, &ec2.DescribePlacementGroupsInput{
		GroupIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidPlacementGroup.Unknown" {
			return nil
		}
		return err
	}

	for _, placementGroup := range response.PlacementGroups {
		if _, err := client.DeletePlacementGroupWithContext(ctx, &ec2.DeletePlacementGroupInput{
			GroupName: placementGroup.GroupName,
		}); err != nil {
			return err
		}
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2RouteTable(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeRouteTablesWithContext(ctx, &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidRouteTableID.NotFound" {
			return nil
		}
		return err
	}

	for _, table := range response.RouteTables {
		err = deleteEC2RouteTableObject(ctx, client, table, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteEC2RouteTableObject(ctx context.Context, client *ec2.EC2, table *ec2.RouteTable, logger logrus.FieldLogger) error {
	hasMain := false
	for _, association := range table.Associations {
		if *association.Main {
			// can't remove the 'Main' association
			hasMain = true
			continue
		}
		_, err := client.DisassociateRouteTableWithContext(ctx, &ec2.DisassociateRouteTableInput{
			AssociationId: association.RouteTableAssociationId,
		})
		if err != nil {
			return errors.Wrapf(err, "dissociating %s", *association.RouteTableAssociationId)
		}
		logger.WithField("id", *association.RouteTableAssociationId).Info("Disassociated")
	}

	if hasMain {
		// can't delete route table with the 'Main' association
		// it will get cleaned up as part of deleting the VPC
		return nil
	}

	_, err := client.DeleteRouteTableWithContext(ctx, &ec2.DeleteRouteTableInput{
		RouteTableId: table.RouteTableId,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2RouteTablesByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeRouteTablesPagesWithContext(
		ctx,
		&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeRouteTablesOutput, lastPage bool) bool {
			for _, table := range results.RouteTables {
				err := deleteEC2RouteTableObject(ctx, client, table, logger.WithField("table", *table.RouteTableId))
				if err != nil {
					if lastError != nil {
						logger.Debug(err)
					}
					lastError = errors.Wrapf(err, "deleting EC2 route table %s", *table.RouteTableId)
					if failFast {
						return false
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2SecurityGroup(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidGroup.NotFound" {
			return nil
		}
		return err
	}

	for _, group := range response.SecurityGroups {
		err = deleteEC2SecurityGroupObject(ctx, client, group, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteEC2SecurityGroupObject(ctx context.Context, client *ec2.EC2, group *ec2.SecurityGroup, logger logrus.FieldLogger) error {
	if len(group.IpPermissions) > 0 {
		_, err := client.RevokeSecurityGroupIngressWithContext(ctx, &ec2.RevokeSecurityGroupIngressInput{
			GroupId:       group.GroupId,
			IpPermissions: group.IpPermissions,
		})
		if err != nil {
			return errors.Wrap(err, "revoking ingress permissions")
		}
		logger.Debug("Revoked ingress permissions")
	}

	if len(group.IpPermissionsEgress) > 0 {
		_, err := client.RevokeSecurityGroupEgressWithContext(ctx, &ec2.RevokeSecurityGroupEgressInput{
			GroupId:       group.GroupId,
			IpPermissions: group.IpPermissionsEgress,
		})
		if err != nil {
			return errors.Wrap(err, "revoking egress permissions")
		}
		logger.Debug("Revoked egress permissions")
	}

	if group.GroupName != nil && *group.GroupName == "default" {
		logger.Debug("Skipping default security group")
		return nil
	}

	_, err := client.DeleteSecurityGroupWithContext(ctx, &ec2.DeleteSecurityGroupInput{
		GroupId: group.GroupId,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidGroup.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2SecurityGroupsByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeSecurityGroupsPagesWithContext(
		ctx,
		&ec2.DescribeSecurityGroupsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			for _, group := range results.SecurityGroups {
				err := deleteEC2SecurityGroupObject(ctx, client, group, logger.WithField("security group", *group.GroupId))
				if err != nil {
					if lastError != nil {
						logger.Debug(err)
					}
					lastError = errors.Wrapf(err, "deleting EC2 security group %s", *group.GroupId)
					if failFast {
						return false
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2Snapshot(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteSnapshotWithContext(ctx, &ec2.DeleteSnapshotInput{
		SnapshotId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidSnapshot.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NetworkInterface(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteNetworkInterfaceWithContext(ctx, &ec2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidNetworkInterfaceID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NetworkInterfaceByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeNetworkInterfacesPagesWithContext(
		ctx,
		&ec2.DescribeNetworkInterfacesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeNetworkInterfacesOutput, lastPage bool) bool {
			for _, networkInterface := range results.NetworkInterfaces {
				err := deleteEC2NetworkInterface(ctx, client, *networkInterface.NetworkInterfaceId, logger.WithField("network interface", *networkInterface.NetworkInterfaceId))
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting EC2 network interface %s", *networkInterface.NetworkInterfaceId)
					if failFast {
						return false
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2Subnet(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteSubnetWithContext(ctx, &ec2.DeleteSubnetInput{
		SubnetId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidSubnetID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2SubnetsByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeSubnetsPagesWithContext(
		ctx,
		&ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			for _, subnet := range results.Subnets {
				err := deleteEC2Subnet(ctx, client, *subnet.SubnetId, logger.WithField("subnet", *subnet.SubnetId))
				if err != nil {
					err = errors.Wrapf(err, "deleting EC2 subnet %s", *subnet.SubnetId)
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = err
					if failFast {
						return false
					}
				}
			}
			return !lastPage
		},
	)
	if err != nil {
		return err
	}

	return lastError
}

func deleteEC2Volume(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVolumeWithContext(ctx, &ec2.DeleteVolumeInput{
		VolumeId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVolume.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPC(ctx context.Context, ec2Client *ec2.EC2, elbClient *elb.ELB, elbv2Client *elbv2.ELBV2, id string, logger logrus.FieldLogger) error {
	// first delete any Load Balancers under this VPC (not all of them are tagged)
	v1lbError := deleteElasticLoadBalancerClassicByVPC(ctx, elbClient, id, logger)
	v2lbError := deleteElasticLoadBalancerV2ByVPC(ctx, elbv2Client, id, logger)
	if v1lbError != nil {
		if v2lbError != nil {
			logger.Info(v2lbError)
		}
		return v1lbError
	} else if v2lbError != nil {
		return v2lbError
	}

	for _, child := range []struct {
		helper   func(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error
		failFast bool
	}{
		{helper: deleteEC2NATGatewaysByVPC, failFast: true},      // not always tagged
		{helper: deleteEC2NetworkInterfaceByVPC, failFast: true}, // not always tagged
		{helper: deleteEC2RouteTablesByVPC, failFast: true},      // not always tagged
		{helper: deleteEC2SecurityGroupsByVPC, failFast: false},  // not always tagged
		{helper: deleteEC2SubnetsByVPC, failFast: true},          // not always tagged
		{helper: deleteEC2VPCEndpointsByVPC, failFast: true},     // not taggable
	} {
		err := child.helper(ctx, ec2Client, id, child.failFast, logger)
		if err != nil {
			return err
		}
	}

	_, err := ec2Client.DeleteVpcWithContext(ctx, &ec2.DeleteVpcInput{
		VpcId: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpoint(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVpcEndpointsWithContext(ctx, &ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: []*string{aws.String(id)},
	})
	if err != nil {
		return errors.Wrapf(err, "cannot delete VPC endpoint %s", id)
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpointsByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	response, err := client.DescribeVpcEndpointsWithContext(ctx, &ec2.DescribeVpcEndpointsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpc)},
			},
		},
	})
	if err != nil {
		return err
	}

	for _, endpoint := range response.VpcEndpoints {
		err := deleteEC2VPCEndpoint(ctx, client, *endpoint.VpcEndpointId, logger.WithField("VPC endpoint", *endpoint.VpcEndpointId))
		if err != nil {
			if err.(awserr.Error).Code() == "InvalidVpcID.NotFound" {
				return nil
			}
			return err
		}
	}

	return nil
}

func deleteEC2VPCPeeringConnection(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVpcPeeringConnectionWithContext(ctx, &ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVpcPeeringConnectionID.NotFound" {
			return nil
		}
		return errors.Wrapf(err, "cannot delete VPC Peering Connection %s", id)
	}
	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpointService(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	output, err := client.DescribeVpcEndpointConnectionsWithContext(ctx, &ec2.DescribeVpcEndpointConnectionsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("service-id"),
				Values: aws.StringSlice([]string{id}),
			},
		},
	})

	if err != nil {
		logger.Warn("Unable to get the list of VPC endpoint connections connected to service: ", err)
		logger.Warn("Attempting to delete the VPC Endpoint Service")
	} else {
		endpointList := make([]*string, len(output.VpcEndpointConnections))
		for _, endpoint := range output.VpcEndpointConnections {
			if aws.StringValue(endpoint.VpcEndpointState) != "rejected" {
				endpointList = append(endpointList, endpoint.VpcEndpointId)
			}
		}

		_, err = client.RejectVpcEndpointConnectionsWithContext(ctx, &ec2.RejectVpcEndpointConnectionsInput{
			ServiceId:      &id,
			VpcEndpointIds: endpointList,
		})

		if err != nil {
			logger.Warn("Unable to reject VPC endpoint connections for service: ", err)
			logger.Warn("Attempting to delete the VPC Endpoint Service")
		} else {
			logger.WithField("resourceType", "VPC Endpoint Connection").Info("Rejected")
		}
	}

	_, err = client.DeleteVpcEndpointServiceConfigurationsWithContext(ctx, &ec2.DeleteVpcEndpointServiceConfigurationsInput{
		ServiceIds: aws.StringSlice([]string{id}),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVpcEndpointService.NotFound" {
			return nil
		}
		return errors.Wrapf(err, "cannot delete VPC Endpoint Service %s", id)
	}
	logger.Info("Deleted")
	return nil
}
