package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

// findEC2Instances returns the EC2 instances with tags that satisfy the filters.
// returns two lists, first one is the list of all resources that are not terminated and are not in shutdown
// stage and the second list is the list of resources that are not terminated.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func findEC2Instances(ctx context.Context, ec2Client *EC2Client, deleted sets.String, filters []Filter, logger logrus.FieldLogger) ([]string, []string, error) {
	region := ec2Client.client.Config.Region
	if region == nil {
		return nil, nil, errors.New("EC2 client does not have region configured")
	}
	return findEC2InstancesWithRegion(ctx, ec2Client, deleted, filters, logger, *region)
}

// This is defined separately so it can be unit tested
func findEC2InstancesWithRegion(ctx context.Context, ec2Client EC2API, deleted sets.String, filters []Filter, logger logrus.FieldLogger, region string) ([]string, []string, error) {
	partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region)
	if !ok {
		return nil, nil, errors.Errorf("no partition found for region %q", region)
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
		err := ec2Client.DescribeInstancesPages(ctx, instanceFilters,
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
						arn := fmt.Sprintf("arn:%s:ec2:%s:%s:instance/%s", partition.ID(), region, *reservation.OwnerId, *instance.InstanceId)
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
	client := NewEC2Client(session, logger)

	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "dhcp-options":
		return client.DeleteDhcpOptions(ctx, id)
	case "elastic-ip":
		return client.DeleteElasticIP(ctx, id)
	case "image":
		return client.DeleteImage(ctx, id)
	case "instance":
		iamClient := NewIAMClient(session, logger)
		return terminateEC2Instance(ctx, client, iamClient, id, logger)
	case "internet-gateway":
		return client.DeleteInternetGateway(ctx, id)
	case "natgateway":
		return client.DeleteNATGateway(ctx, id)
	case "placement-group":
		return client.DeletePlacementGroup(ctx, id)
	case "route-table":
		return deleteEC2RouteTable(ctx, client, id, logger)
	case "security-group":
		return client.DeleteSecurityGroup(ctx, id)
	case "snapshot":
		return client.DeleteSnapshot(ctx, id)
	case "network-interface":
		return client.DeleteNetworkInterface(ctx, id)
	case "subnet":
		return client.DeleteSubnet(ctx, id)
	case "volume":
		return client.DeleteVolume(ctx, id)
	case "vpc":
		return deleteEC2VPC(ctx, client, session, id, logger)
	case "vpc-endpoint":
		return client.DeleteVPCEndpoint(ctx, id)
	case "vpc-peering-connection":
		return client.DeleteVPCPeeringConnection(ctx, id)
	case "vpc-endpoint-service":
		return client.DeleteVPCEndpointService(ctx, id)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func terminateEC2Instance(ctx context.Context, ec2Client EC2API, iamClient IAMAPI, id string, logger logrus.FieldLogger) error {
	response, err := ec2Client.DescribeInstances(ctx, id)
	if err != nil {
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

func terminateEC2InstanceByInstance(ctx context.Context, ec2Client EC2API, iamClient IAMAPI, instance *ec2.Instance, logger logrus.FieldLogger) error {
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

	err := ec2Client.TerminateInstances(ctx, instance.InstanceId)
	if err != nil {
		return err
	}

	logger.Debug("Terminating")
	return nil
}

func deleteEC2NATGatewaysByVPC(ctx context.Context, client EC2API, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeNatGatewaysPages(
		ctx,
		vpc,
		func(results *ec2.DescribeNatGatewaysOutput, lastPage bool) bool {
			for _, gateway := range results.NatGateways {
				err := client.DeleteNATGateway(ctx, *gateway.NatGatewayId)
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

func deleteEC2RouteTable(ctx context.Context, client EC2API, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeRouteTables(ctx, id)
	if err != nil {
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

func deleteEC2RouteTableObject(ctx context.Context, client EC2API, table *ec2.RouteTable, logger logrus.FieldLogger) error {
	hasMain := false
	for _, association := range table.Associations {
		if *association.Main {
			// can't remove the 'Main' association
			hasMain = true
			continue
		}
		err := client.DisassociateRouteTable(ctx, *association.RouteTableAssociationId)
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

	err := client.DeleteRouteTable(ctx, *table.RouteTableId)
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2RouteTablesByVPC(ctx context.Context, client EC2API, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeRouteTablesPages(ctx, vpc,
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

func deleteEC2SecurityGroupObject(ctx context.Context, client EC2API, group *ec2.SecurityGroup, logger logrus.FieldLogger) error {
	if len(group.IpPermissions) > 0 {
		err := client.RevokeSecurityGroupIngress(ctx, group)
		if err != nil {
			return errors.Wrap(err, "revoking ingress permissions")
		}
		logger.Debug("Revoked ingress permissions")
	}

	if len(group.IpPermissionsEgress) > 0 {
		err := client.RevokeSecurityGroupEgress(ctx, group)
		if err != nil {
			return errors.Wrap(err, "revoking egress permissions")
		}
		logger.Debug("Revoked egress permissions")
	}

	if group.GroupName != nil && *group.GroupName == "default" {
		logger.Debug("Skipping default security group")
		return nil
	}

	err := client.DeleteSecurityGroup(ctx, *group.GroupId)
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2SecurityGroupsByVPC(ctx context.Context, client EC2API, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeSecurityGroupsPages(ctx, vpc,
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

func deleteEC2NetworkInterfaceByVPC(ctx context.Context, client EC2API, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeNetworkInterfacesPages(ctx, vpc,
		func(results *ec2.DescribeNetworkInterfacesOutput, lastPage bool) bool {
			for _, networkInterface := range results.NetworkInterfaces {
				err := client.DeleteNetworkInterface(ctx, *networkInterface.NetworkInterfaceId)
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

func deleteEC2SubnetsByVPC(ctx context.Context, client EC2API, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeSubnetsPages(ctx, vpc,
		func(results *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			for _, subnet := range results.Subnets {
				err := client.DeleteSubnet(ctx, *subnet.SubnetId)
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

func deleteEC2VPC(ctx context.Context, ec2Client EC2API, session *session.Session, id string, logger logrus.FieldLogger) error {
	// first delete any Load Balancers under this VPC (not all of them are tagged)
	err := deleteElasticLoadBalancingByVPC(ctx, session, id, logger)
	if err != nil {
		return err
	}

	for _, child := range []struct {
		helper   func(ctx context.Context, client EC2API, vpc string, failFast bool, logger logrus.FieldLogger) error
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

	err = ec2Client.DeleteVPC(ctx, id)
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpointsByVPC(ctx context.Context, client EC2API, vpc string, failFast bool, logger logrus.FieldLogger) error {
	response, err := client.DescribeVpcEndpoints(ctx, vpc)
	if err != nil {
		return err
	}

	for _, endpoint := range response.VpcEndpoints {
		err := client.DeleteVPCEndpoint(ctx, *endpoint.VpcEndpointId)
		if err != nil {
			return err
		}
	}

	return nil
}
