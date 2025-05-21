package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2v2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	iamv2 "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
)

// findEC2Instances returns the EC2 instances with tags that satisfy the filters.
// returns two lists, first one is the list of all resources that are not terminated and are not in shutdown
// stage and the second list is the list of resources that are not terminated.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func findEC2Instances(ctx context.Context, ec2Client *ec2v2.Client, deleted sets.Set[string], filters []Filter, logger logrus.FieldLogger) ([]string, []string, error) {
	if ec2Client.Options().Region == "" {
		return nil, nil, errors.New("EC2 client does not have region configured")
	}

	var resourcesRunning []string
	var resourcesNotTerminated []string
	for _, filter := range filters {
		logger.Debugf("search for instances by tag matching %#+v", filter)
		instanceFilters := make([]ec2v2types.Filter, 0, len(filter))
		for key, value := range filter {
			instanceFilters = append(instanceFilters, ec2v2types.Filter{
				Name:   aws.String("tag:" + key),
				Values: []string{value},
			})
		}

		paginator := ec2v2.NewDescribeInstancesPaginator(ec2Client, &ec2v2.DescribeInstancesInput{
			Filters: instanceFilters,
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to fetch EC2 instances: %w", err)
			}

			for _, reservation := range page.Reservations {
				if reservation.OwnerId == nil {
					continue
				}

				for _, instance := range reservation.Instances {
					if instance.InstanceId == nil || instance.State == nil {
						continue
					}

					instanceLogger := logger.WithField("instance", *instance.InstanceId)
					arn := fmt.Sprintf("arn:aws:ec2:%s:%s:instance/%s", ec2Client.Options().Region, *reservation.OwnerId, *instance.InstanceId)
					if instance.State.Name == "terminated" {
						if !deleted.Has(arn) {
							instanceLogger.Info("Terminated")
							deleted.Insert(arn)
						}
						continue
					}
					if instance.State.Name != "shutting-down" {
						resourcesRunning = append(resourcesRunning, arn)
					}
					resourcesNotTerminated = append(resourcesNotTerminated, arn)
				}
			}

		}
	}

	return resourcesRunning, resourcesNotTerminated, nil
}

// DeleteEC2Instances terminates all EC2 instances found.
func (o *ClusterUninstaller) DeleteEC2Instances(ctx context.Context, awsSession *session.Session, toDelete sets.Set[string], deleted sets.Set[string], tracker *ErrorTracker) error {
	lastTerminateTime := time.Now()
	err := wait.PollUntilContextCancel(
		ctx,
		time.Second*10,
		true,
		func(ctx context.Context) (bool, error) {
			instancesRunning, instancesNotTerminated, err := findEC2Instances(ctx, o.EC2Client, deleted, o.Filters, o.Logger)
			if err != nil {
				o.Logger.WithError(err).Info("error while finding EC2 instances to delete")
				return false, nil
			}
			if len(instancesNotTerminated) == 0 && len(instancesRunning) == 0 {
				return true, nil
			}
			instancesToDelete := instancesRunning
			if time.Since(lastTerminateTime) > 10*time.Minute {
				instancesToDelete = instancesNotTerminated
				lastTerminateTime = time.Now()
			}
			newlyDeleted, err := o.DeleteResources(ctx, awsSession, instancesToDelete, tracker)
			// Delete from the resources-to-delete set so that the current state of the resources to delete can be
			// returned if the context is completed.
			toDelete = toDelete.Difference(newlyDeleted)
			deleted = deleted.Union(newlyDeleted)
			if err != nil {
				o.Logger.WithError(err).Info("error while deleting EC2 instances")
			}
			return false, nil
		},
	)
	return err
}

func (o *ClusterUninstaller) deleteEC2(ctx context.Context, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id).WithField("resourceType", resourceType)

	switch resourceType {
	case "dhcp-options":
		return deleteEC2DHCPOptions(ctx, o.EC2Client, id, logger)
	case "elastic-ip":
		return deleteEC2ElasticIP(ctx, o.EC2Client, id, logger)
	case "image":
		return deleteEC2Image(ctx, o.EC2Client, id, logger)
	case "instance":
		return terminateEC2Instance(ctx, o.EC2Client, o.IAMClient, id, logger)
	case "internet-gateway":
		return deleteEC2InternetGateway(ctx, o.EC2Client, id, logger)
	case "carrier-gateway":
		return deleteEC2CarrierGateway(ctx, o.EC2Client, id, logger)
	case "natgateway":
		return deleteEC2NATGateway(ctx, o.EC2Client, id, logger)
	case "placement-group":
		return deleteEC2PlacementGroup(ctx, o.EC2Client, id, logger)
	case "route-table":
		return deleteEC2RouteTable(ctx, o.EC2Client, id, logger)
	case "security-group":
		return deleteEC2SecurityGroup(ctx, o.EC2Client, id, logger)
	case "snapshot":
		return deleteEC2Snapshot(ctx, o.EC2Client, id, logger)
	case "network-interface":
		return deleteEC2NetworkInterface(ctx, o.EC2Client, id, logger)
	case "subnet":
		return deleteEC2Subnet(ctx, o.EC2Client, id, logger)
	case "volume":
		return deleteEC2Volume(ctx, o.EC2Client, id, logger)
	case "vpc":
		return deleteEC2VPC(ctx, o.EC2Client, o.ELBClient, o.ELBV2Client, id, logger)
	case "vpc-endpoint":
		return deleteEC2VPCEndpoint(ctx, o.EC2Client, id, logger)
	case "vpc-peering-connection":
		return deleteEC2VPCPeeringConnection(ctx, o.EC2Client, id, logger)
	case "vpc-endpoint-service":
		return deleteEC2VPCEndpointService(ctx, o.EC2Client, id, logger)
	case "egress-only-internet-gateway":
		return deleteEgressOnlyInternetGateway(ctx, o.EC2Client, id, logger)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func deleteEC2DHCPOptions(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteDhcpOptions(ctx, &ec2v2.DeleteDhcpOptionsInput{
		DhcpOptionsId: &id,
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidDhcpOptions.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2Image(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	// tag the snapshots used by the AMI so that the snapshots are matched
	// by the filter and deleted
	response, err := client.DescribeImages(ctx, &ec2v2.DescribeImagesInput{
		ImageIds: []string{id},
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidAMI.NotFound" {
			return nil
		}
		return err
	}
	for _, image := range response.Images {
		// Early detection: Check if the AMI is managed by AWS Backup service.
		// AWS Backup-managed AMIs cannot be deleted via EC2 APIs and will fail with
		// AuthFailure. Detecting this early via tags avoids unnecessary API calls.
		for _, tag := range image.Tags {
			if aws.StringValue(tag.Key) == "aws:backup:source-resource" {
				logger.Warnf("Skipping AMI image %s deletion since it is managed by the AWS Backup service. To delete this image, please use the AWS Backup APIs, CLI, or console", id)
				return nil
			}
		}

		var snapshots []string
		for _, bdm := range image.BlockDeviceMappings {
			if bdm.Ebs != nil && bdm.Ebs.SnapshotId != nil {
				snapshots = append(snapshots, *bdm.Ebs.SnapshotId)
			}
		}
		if len(snapshots) != 0 {
			_, err = client.CreateTags(ctx, &ec2v2.CreateTagsInput{
				Resources: snapshots,
				Tags:      image.Tags,
			})
			if err != nil {
				return errors.Wrapf(err, "tagging snapshots for %s", id)
			}
		}
	}

	_, err = client.DeregisterImage(ctx, &ec2v2.DeregisterImageInput{
		ImageId: &id,
	})
	if err != nil {
		errCode := HandleErrorCode(err)
		switch errCode {
		case "InvalidAMI.NotFound":
			return nil
		case "AuthFailure":
			// AMIs, managed by AWS Backup service, cannot be deleted via EC2 APIs. When attempting to delete, the following error is returned
			//
			// AuthFailure: This image is managed by AWS Backup and cannot be deleted via EC2 APIs. To delete this image, please use the AWS Backup APIs, CLI, or console.
			//
			// The installer cannot handle this case automatically. Users must manually delete the AMI, following AWS instructions.
			// This is a fallback in case the tag check above didn't catch it (e.g., if AWS Backup tags were removed or not present).
			logger.WithError(err).Warnf("Skipping AMI image %s deletion", id)
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2ElasticIP(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.ReleaseAddress(ctx, &ec2v2.ReleaseAddressInput{
		AllocationId: aws.String(id),
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidAllocation.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Released")
	return nil
}

func terminateEC2Instance(ctx context.Context, ec2Client *ec2v2.Client, iamClient *iamv2.Client, id string, logger logrus.FieldLogger) error {
	response, err := ec2Client.DescribeInstances(ctx, &ec2v2.DescribeInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidInstance.NotFound" {
			return nil
		}
		return err
	}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			err = terminateEC2InstanceByInstance(ctx, ec2Client, iamClient, &instance, logger)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func terminateEC2InstanceByInstance(ctx context.Context, ec2Client *ec2v2.Client, iamClient *iamv2.Client, instance *ec2v2types.Instance, logger logrus.FieldLogger) error {
	// Ignore instances that are already terminated
	if instance.State == nil || instance.State.Name == "terminated" {
		return nil
	}

	_, err := ec2Client.TerminateInstances(ctx, &ec2v2.TerminateInstancesInput{
		InstanceIds: []string{*instance.InstanceId},
	})
	if err != nil {
		return err
	}

	logger.Debug("Terminating")
	return nil
}

func deleteEC2InternetGateway(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeInternetGateways(ctx, &ec2v2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []string{id},
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
			_, err := client.DetachInternetGateway(ctx, &ec2v2.DetachInternetGatewayInput{
				InternetGatewayId: gateway.InternetGatewayId,
				VpcId:             vpc.VpcId,
			})
			if err == nil {
				logger.WithField("vpc", *vpc.VpcId).Debug("Detached")
			} else if HandleErrorCode(err) == "Gateway.NotAttached" {
				return nil
			}
		}
	}

	_, err = client.DeleteInternetGateway(ctx, &ec2v2.DeleteInternetGatewayInput{
		InternetGatewayId: &id,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2CarrierGateway(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteCarrierGateway(ctx, &ec2v2.DeleteCarrierGatewayInput{
		CarrierGatewayId: &id,
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidCarrierGateway.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NATGateway(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteNatGateway(ctx, &ec2v2.DeleteNatGatewayInput{
		NatGatewayId: aws.String(id),
	})
	if err != nil {
		if HandleErrorCode(err) == "NatGateway.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NATGatewaysByVPC(ctx context.Context, client *ec2v2.Client, vpc string, failFast bool, logger logrus.FieldLogger) error {
	paginator := ec2v2.NewDescribeNatGatewaysPaginator(client, &ec2v2.DescribeNatGatewaysInput{
		Filter: []ec2v2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpc},
			},
		},
	})

	var lastError error

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("describing NAT gateways for VPC %s: %w", vpc, err)
		}

		for _, gateway := range page.NatGateways {
			err := deleteEC2NATGateway(ctx, client, *gateway.NatGatewayId, logger.WithField("NAT gateway", *gateway.NatGatewayId))
			if err != nil {
				if lastError != nil {
					logger.Debug(err)
				}
				lastError = fmt.Errorf("deleting EC2 NAT gateway %s: %w", *gateway.NatGatewayId, err)
				if failFast {
					break
				}
			}
		}
		if failFast && lastError != nil {
			break
		}
	}

	return lastError
}

func deleteEC2PlacementGroup(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribePlacementGroups(ctx, &ec2v2.DescribePlacementGroupsInput{
		GroupIds: []string{id},
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidPlacementGroup.NotFound" {
			return nil
		}
		return err
	}

	for _, placementGroup := range response.PlacementGroups {
		if _, err := client.DeletePlacementGroup(ctx, &ec2v2.DeletePlacementGroupInput{
			GroupName: placementGroup.GroupName,
		}); err != nil {
			return err
		}
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2RouteTable(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeRouteTables(ctx, &ec2v2.DescribeRouteTablesInput{
		RouteTableIds: []string{id},
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidRouteTableID.NotFound" {
			return nil
		}
		return err
	}

	for _, table := range response.RouteTables {
		err = deleteEC2RouteTableObject(ctx, client, &table, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteEC2RouteTableObject(ctx context.Context, client *ec2v2.Client, table *ec2v2types.RouteTable, logger logrus.FieldLogger) error {
	hasMain := false
	for _, association := range table.Associations {
		if *association.Main {
			// can't remove the 'Main' association
			hasMain = true
			continue
		}
		_, err := client.DisassociateRouteTable(ctx, &ec2v2.DisassociateRouteTableInput{
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

	_, err := client.DeleteRouteTable(ctx, &ec2v2.DeleteRouteTableInput{
		RouteTableId: table.RouteTableId,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2RouteTablesByVPC(ctx context.Context, client *ec2v2.Client, vpc string, failFast bool, logger logrus.FieldLogger) error {
	paginator := ec2v2.NewDescribeRouteTablesPaginator(client, &ec2v2.DescribeRouteTablesInput{
		Filters: []ec2v2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpc},
			},
		},
	})

	var lastError error
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("describing route tables for VPC %s: %w", vpc, err)
		}

		for _, table := range page.RouteTables {
			err := deleteEC2RouteTableObject(ctx, client, &table, logger.WithField("table", *table.RouteTableId))
			if err != nil {
				if lastError != nil {
					logger.Debug(err)
				}
				lastError = fmt.Errorf("deleting EC2 route table %s: %w", *table.RouteTableId, err)
				if failFast {
					break
				}
			}
		}
		if failFast && lastError != nil {
			break
		}
	}

	return lastError
}

func formatIPPermissions(perms []ec2v2types.IpPermission) []ec2v2types.IpPermission {
	formattedPerms := []ec2v2types.IpPermission{}

	for _, ipPermission := range perms {
		item := ec2v2types.IpPermission{
			FromPort:   ipPermission.FromPort,
			ToPort:     ipPermission.ToPort,
			IpProtocol: ipPermission.IpProtocol,
		}
		if len(ipPermission.IpRanges) > 0 {
			item.IpRanges = ipPermission.IpRanges
		}
		if len(ipPermission.Ipv6Ranges) > 0 {
			item.Ipv6Ranges = ipPermission.Ipv6Ranges
		}
		if len(ipPermission.PrefixListIds) > 0 {
			item.PrefixListIds = ipPermission.PrefixListIds
		}
		if len(ipPermission.UserIdGroupPairs) > 0 {
			item.UserIdGroupPairs = ipPermission.UserIdGroupPairs
		}
		formattedPerms = append(formattedPerms, item)
	}
	return formattedPerms
}

func formatSecurityGroup(sg ec2v2types.SecurityGroup) ec2v2types.SecurityGroup {
	return ec2v2types.SecurityGroup{
		GroupId:             sg.GroupId,
		GroupName:           sg.GroupName,
		IpPermissions:       formatIPPermissions(sg.IpPermissions),
		IpPermissionsEgress: formatIPPermissions(sg.IpPermissionsEgress),
		VpcId:               sg.VpcId,
		OwnerId:             sg.OwnerId,
		Description:         sg.Description,
	}
}

func deleteEC2SecurityGroup(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeSecurityGroups(ctx, &ec2v2.DescribeSecurityGroupsInput{
		GroupIds: []string{id},
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidGroup.NotFound" {
			return nil
		}
		return err
	}

	for _, group := range response.SecurityGroups {

		formattedGroup := formatSecurityGroup(group)
		err = deleteEC2SecurityGroupObject(ctx, client, &formattedGroup, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteEC2SecurityGroupObject(ctx context.Context, client *ec2v2.Client, group *ec2v2types.SecurityGroup, logger logrus.FieldLogger) error {
	if len(group.IpPermissions) > 0 {
		_, err := client.RevokeSecurityGroupIngress(ctx, &ec2v2.RevokeSecurityGroupIngressInput{
			GroupId:       group.GroupId,
			IpPermissions: group.IpPermissions,
		})
		if err != nil {
			return errors.Wrap(err, "revoking ingress permissions")
		}
		logger.Debug("Revoked ingress permissions")
	}

	if len(group.IpPermissionsEgress) > 0 {
		_, err := client.RevokeSecurityGroupEgress(ctx, &ec2v2.RevokeSecurityGroupEgressInput{
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

	_, err := client.DeleteSecurityGroup(ctx, &ec2v2.DeleteSecurityGroupInput{
		GroupId: group.GroupId,
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidGroup.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2SecurityGroupsByVPC(ctx context.Context, client *ec2v2.Client, vpc string, failFast bool, logger logrus.FieldLogger) error {
	paginator := ec2v2.NewDescribeSecurityGroupsPaginator(client, &ec2v2.DescribeSecurityGroupsInput{
		Filters: []ec2v2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpc},
			},
		},
	})

	var lastError error
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("describing security groups for VPC %s: %w", vpc, err)
		}

		for _, group := range page.SecurityGroups {
			formattedGroup := formatSecurityGroup(group)

			err := deleteEC2SecurityGroupObject(ctx, client, &formattedGroup, logger.WithField("security group", *group.GroupId))
			if err != nil {
				if lastError != nil {
					logger.Debug(err)
				}
				lastError = fmt.Errorf("deleting EC2 security group %s: %w", *group.GroupId, err)
				if failFast {
					break
				}
			}
		}

		if failFast && lastError != nil {
			break
		}
	}
	return lastError
}

func deleteEC2Snapshot(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteSnapshot(ctx, &ec2v2.DeleteSnapshotInput{
		SnapshotId: &id,
	})
	if err != nil {
		errCode := HandleErrorCode(err)
		switch errCode {
		case "InvalidSnapshot.NotFound":
			return nil
		case "InvalidParameterValue":
			// An InvalidParameterValue indicates the AWS request parameter is not valid, is unsupported, or cannot be used.
			// For example, snapshots, managed by the AWS Backup service, cannot be deleted via EC2 APIs. When attempting to delete, the following error is returned:
			//
			// InvalidParameterValue: This snapshot is managed by the AWS Backup service and cannot be deleted via EC2 APIs.
			// If you wish to delete this snapshot, please do so via the Backup console.
			//
			// The installer should not try to delete these backup snapshots, but leave it to the users to clean them up.
			logger.WithError(err).Warnf("Skipping snapshot %s deletion", id)
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NetworkInterface(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteNetworkInterface(ctx, &ec2v2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(id),
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidNetworkInterfaceID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NetworkInterfaceByVPC(ctx context.Context, client *ec2v2.Client, vpc string, failFast bool, logger logrus.FieldLogger) error {
	paginator := ec2v2.NewDescribeNetworkInterfacesPaginator(client, &ec2v2.DescribeNetworkInterfacesInput{
		Filters: []ec2v2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpc},
			},
		},
	})

	var lastError error
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("describing network interface for VPC %s: %w", vpc, err)
		}

		for _, networkInterface := range page.NetworkInterfaces {
			err := deleteEC2NetworkInterface(ctx, client, *networkInterface.NetworkInterfaceId, logger.WithField("network interface", *networkInterface.NetworkInterfaceId))
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("deleting EC2 network interface %s: %w", *networkInterface.NetworkInterfaceId, err)
				if failFast {
					break
				}
			}
		}
		if failFast && lastError != nil {
			break
		}
	}
	return lastError
}

func deleteEC2Subnet(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteSubnet(ctx, &ec2v2.DeleteSubnetInput{
		SubnetId: aws.String(id),
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidSubnetID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2SubnetsByVPC(ctx context.Context, client *ec2v2.Client, vpc string, failFast bool, logger logrus.FieldLogger) error {
	paginator := ec2v2.NewDescribeSubnetsPaginator(client, &ec2v2.DescribeSubnetsInput{
		Filters: []ec2v2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpc},
			},
		},
	})

	var lastError error
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("describing subnets for VPC %s: %w", vpc, err)
		}

		for _, subnet := range page.Subnets {
			err := deleteEC2Subnet(ctx, client, *subnet.SubnetId, logger.WithField("subnet", *subnet.SubnetId))
			if err != nil {
				err = errors.Wrapf(err, "deleting EC2 subnet %s", *subnet.SubnetId)
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = err
				if failFast {
					break
				}
			}
		}
		if failFast && lastError != nil {
			break
		}
	}
	return lastError
}

func deleteEC2Volume(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVolume(ctx, &ec2v2.DeleteVolumeInput{
		VolumeId: aws.String(id),
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidVolume.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPC(ctx context.Context, ec2Client *ec2v2.Client, elbClient *elb.Client, elbv2Client *elbv2.Client, id string, logger logrus.FieldLogger) error {
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
		helper   func(ctx context.Context, client *ec2v2.Client, vpc string, failFast bool, logger logrus.FieldLogger) error
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

	_, err := ec2Client.DeleteVpc(ctx, &ec2v2.DeleteVpcInput{
		VpcId: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpoint(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVpcEndpoints(ctx, &ec2v2.DeleteVpcEndpointsInput{
		VpcEndpointIds: []string{id},
	})
	if err != nil {
		return errors.Wrapf(err, "cannot delete VPC endpoint %s", id)
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpointsByVPC(ctx context.Context, client *ec2v2.Client, vpc string, failFast bool, logger logrus.FieldLogger) error {
	response, err := client.DescribeVpcEndpoints(ctx, &ec2v2.DescribeVpcEndpointsInput{
		Filters: []ec2v2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpc},
			},
		},
	})
	if err != nil {
		return err
	}

	for _, endpoint := range response.VpcEndpoints {
		err := deleteEC2VPCEndpoint(ctx, client, *endpoint.VpcEndpointId, logger.WithField("VPC endpoint", *endpoint.VpcEndpointId))
		if err != nil {
			if HandleErrorCode(err) == "InvalidVpcID.NotFound" {
				return nil
			}
			return err
		}
	}

	return nil
}

func deleteEC2VPCPeeringConnection(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVpcPeeringConnection(ctx, &ec2v2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: &id,
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidVpcPeeringConnection.NotFound" {
			return nil
		}
		return errors.Wrapf(err, "cannot delete VPC Peering Connection %s", id)
	}
	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpointService(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	output, err := client.DescribeVpcEndpointConnections(ctx, &ec2v2.DescribeVpcEndpointConnectionsInput{
		Filters: []ec2v2types.Filter{
			{
				Name:   aws.String("service-id"),
				Values: []string{id},
			},
		},
	})

	if err != nil {
		logger.Warn("Unable to get the list of VPC endpoint connections connected to service: ", err)
		logger.Warn("Attempting to delete the VPC Endpoint Service")
	} else {
		endpointList := make([]string, 0, len(output.VpcEndpointConnections))
		for _, endpoint := range output.VpcEndpointConnections {
			if endpoint.VpcEndpointState != "rejected" {
				endpointList = append(endpointList, *endpoint.VpcEndpointId)
			}
		}

		_, err = client.RejectVpcEndpointConnections(ctx, &ec2v2.RejectVpcEndpointConnectionsInput{
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

	_, err = client.DeleteVpcEndpointServiceConfigurations(ctx, &ec2v2.DeleteVpcEndpointServiceConfigurationsInput{
		ServiceIds: []string{id},
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidVpcEndpointService.NotFound" {
			return nil
		}
		return errors.Wrapf(err, "cannot delete VPC Endpoint Service %s", id)
	}
	logger.Info("Deleted")
	return nil
}

func deleteEgressOnlyInternetGateway(ctx context.Context, client *ec2v2.Client, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteEgressOnlyInternetGateway(ctx, &ec2v2.DeleteEgressOnlyInternetGatewayInput{
		EgressOnlyInternetGatewayId: aws.String(id),
	})
	if err != nil {
		if HandleErrorCode(err) == "InvalidEgressOnlyInternetGatewayId.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}
