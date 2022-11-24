package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./ec2client.go -destination=mock/ec2client_generated.go -package=mock

// EC2API represents the calls made to the AWS EC2 API
type EC2API interface {
	DeleteDhcpOptions(ctx context.Context, id string) error
	DeleteElasticIP(ctx context.Context, id string) error
	DeleteImage(ctx context.Context, id string) error
	DeleteInternetGateway(ctx context.Context, id string) error
	DeleteNATGateway(ctx context.Context, id string) error
	DeleteNetworkInterface(ctx context.Context, id string) error
	DeleteRouteTable(ctx context.Context, id string) error
	DeleteSecurityGroup(ctx context.Context, id string) error
	DeleteSnapshot(ctx context.Context, id string) error
	DeleteSubnet(ctx context.Context, id string) error
	DeleteVPC(ctx context.Context, id string) error
	DeleteVPCEndpoint(ctx context.Context, id string) error
	DeleteVPCEndpointService(ctx context.Context, id string) error
	DeleteVPCPeeringConnection(ctx context.Context, id string) error
	DeleteVolume(ctx context.Context, id string) error
	DescribeInstances(ctx context.Context, id string) (*ec2.DescribeInstancesOutput, error)
	DescribeInstancesPages(ctx context.Context, filters []*ec2.Filter, fn func(*ec2.DescribeInstancesOutput, bool) bool) error
	DescribeNatGatewaysPages(ctx context.Context, vpc string, fn func(*ec2.DescribeNatGatewaysOutput, bool) bool) error
	DescribeNetworkInterfacesPages(ctx context.Context, vpc string, fn func(*ec2.DescribeNetworkInterfacesOutput, bool) bool) error
	DescribeRouteTables(ctx context.Context, id string) (*ec2.DescribeRouteTablesOutput, error)
	DescribeRouteTablesPages(ctx context.Context, vpc string, fn func(*ec2.DescribeRouteTablesOutput, bool) bool) error
	DescribeSecurityGroupsPages(ctx context.Context, vpc string, fn func(*ec2.DescribeSecurityGroupsOutput, bool) bool) error
	DescribeSubnetsPages(ctx context.Context, vpc string, fn func(*ec2.DescribeSubnetsOutput, bool) bool) error
	DescribeVpcEndpoints(ctx context.Context, vpc string) (*ec2.DescribeVpcEndpointsOutput, error)
	DisassociateRouteTable(ctx context.Context, id string) error
	RevokeSecurityGroupEgress(ctx context.Context, group *ec2.SecurityGroup) error
	RevokeSecurityGroupIngress(ctx context.Context, group *ec2.SecurityGroup) error
	TerminateInstances(ctx context.Context, id *string) error
}

// EC2Client makes calls to the AWS EC2 API
type EC2Client struct {
	client *ec2.EC2
	logger logrus.FieldLogger
}

// NewEC2Client initializes a client
func NewEC2Client(awsSession *session.Session, logger logrus.FieldLogger) *EC2Client {
	return &EC2Client{client: ec2.New(awsSession), logger: logger}
}

// DeleteDhcpOptions deletes DhcpOptions with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteDhcpOptions(ctx context.Context, id string) error {
	logger := c.logger.WithField("id", id)
	_, err := c.client.DeleteDhcpOptionsWithContext(ctx, &ec2.DeleteDhcpOptionsInput{
		DhcpOptionsId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidDhcpOptionsID.NotFound" {
			return nil
		}
		return err
	}
	logger.Info("deleted")
	return nil
	// return err
}

// DeleteElasticIP deletes Elastic IP with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteElasticIP(ctx context.Context, id string) error {
	logger := c.logger.WithField("id", id)
	_, err := c.client.ReleaseAddressWithContext(ctx, &ec2.ReleaseAddressInput{
		AllocationId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAllocationID.NotFound" {
			return nil
		}
		return err
	}
	logger.Info("released")

	return nil
}

// DeleteImage deletes an EC2 Image with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteImage(ctx context.Context, id string) error {
	logger := c.logger.WithField("id", id)
	// tag the snapshots used by the AMI so that the snapshots are matched
	// by the filter and deleted
	response, err := c.client.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(id)},
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
			_, err = c.client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
				Resources: snapshots,
				Tags:      image.Tags,
			})
			if err != nil {
				return errors.Wrapf(err, "tagging snapshots for %s", id)
			}
		}
	}

	_, err = c.client.DeregisterImageWithContext(ctx, &ec2.DeregisterImageInput{
		ImageId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAMIID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("deleted")
	return nil
}

// DeleteInternetGateway deletes a Gateway with ARN=`id`. It ignores GWs without a VPC
func (c *EC2Client) DeleteInternetGateway(ctx context.Context, id string) error {
	logger := c.logger.WithField("id", id)
	response, err := c.client.DescribeInternetGatewaysWithContext(ctx, &ec2.DescribeInternetGatewaysInput{
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
			_, err = c.client.DetachInternetGatewayWithContext(ctx, &ec2.DetachInternetGatewayInput{
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

	_, err = c.client.DeleteInternetGatewayWithContext(ctx, &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: &id,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

// DeleteNATGateway deletes a NAT GW with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteNATGateway(ctx context.Context, id string) error {
	logger := c.logger.WithField("Nat Gateway", id)
	_, err := c.client.DeleteNatGatewayWithContext(ctx, &ec2.DeleteNatGatewayInput{
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

// DeleteNetworkInterface deletes a Network Interface with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteNetworkInterface(ctx context.Context, id string) error {
	logger := c.logger.WithField("Network Interface", id)
	_, err := c.client.DeleteNetworkInterfaceWithContext(ctx, &ec2.DeleteNetworkInterfaceInput{
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

// DeleteVolume deletes an EC2 Volume with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteVolume(ctx context.Context, id string) error {
	logger := c.logger.WithField("Volume", id)
	_, err := c.client.DeleteVolumeWithContext(ctx, &ec2.DeleteVolumeInput{
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

// DeleteVPC deletes a VPC with ARN=`id`
func (c *EC2Client) DeleteVPC(ctx context.Context, id string) error {
	logger := c.logger.WithField("VPC", id)
	_, err := c.client.DeleteVpcWithContext(ctx, &ec2.DeleteVpcInput{
		VpcId: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

// DeleteVPCEndpoint deletes a VPC Endpoint with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteVPCEndpoint(ctx context.Context, id string) error {
	logger := c.logger.WithField("VPC Endpoint", id)
	_, err := c.client.DeleteVpcEndpointsWithContext(ctx, &ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVpcID.NotFound" {
			return nil
		}
		return errors.Wrapf(err, "cannot delete VPC endpoint %s", id)
	}

	logger.Info("Deleted")
	return nil
}

// DeleteVPCEndpointService deletes a VPC Endpoint Service with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteVPCEndpointService(ctx context.Context, id string) error {
	logger := c.logger.WithField("VPC Endpoint Service", id)
	_, err := c.client.DeleteVpcEndpointServiceConfigurationsWithContext(
		ctx,
		&ec2.DeleteVpcEndpointServiceConfigurationsInput{
			ServiceIds: aws.StringSlice([]string{id}),
		},
	)
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVpcEndpointService.NotFound" {
			return nil
		}
		return errors.Wrapf(err, "cannot delete VPC Endpoint Service %s", id)
	}

	logger.Info("deleted")
	return nil
}

// DeleteVPCPeeringConnection deletes a VPC Peering Connection with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteVPCPeeringConnection(ctx context.Context, id string) error {
	logger := c.logger.WithField("id", id)
	_, err := c.client.DeleteVpcPeeringConnectionWithContext(ctx, &ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVpcPeeringConnectionID.NotFound" {
			return nil
		}
		return errors.Wrapf(err, "cannot delete VPC Peering Connection %s", id)
	}
	logger.Info("deleted")

	return nil
}

// DeletePlacementGroup deletes Placement Group with ARN=`id`. It ignores Unknown objects
func (c *EC2Client) DeletePlacementGroup(ctx context.Context, id string) error {
	logger := c.logger.WithField("Placement Group", id)
	response, err := c.client.DescribePlacementGroupsWithContext(ctx, &ec2.DescribePlacementGroupsInput{
		GroupIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidPlacementGroup.Unknown" {
			return nil
		}
		return err
	}

	for _, placementGroup := range response.PlacementGroups {
		if _, err := c.client.DeletePlacementGroupWithContext(ctx, &ec2.DeletePlacementGroupInput{
			GroupName: placementGroup.GroupName,
		}); err != nil {
			return err
		}
	}

	logger.Info("Deleted")
	return nil
}

// DeleteRouteTable deletes a Route Table with ARN=`id`.
func (c *EC2Client) DeleteRouteTable(ctx context.Context, id string) error {
	logger := c.logger.WithField("Route Table", id)
	_, err := c.client.DeleteRouteTableWithContext(ctx, &ec2.DeleteRouteTableInput{
		RouteTableId: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

// DeleteSecurityGroup deletes Security Group with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteSecurityGroup(ctx context.Context, id string) error {
	logger := c.logger.WithField("Security Group", id)
	_, err := c.client.DeleteSecurityGroupWithContext(ctx, &ec2.DeleteSecurityGroupInput{GroupId: aws.String(id)})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidGroup.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return err
}

// DeleteSnapshot deletes a Snapshot with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteSnapshot(ctx context.Context, id string) error {
	logger := c.logger.WithField("Snapshot", id)
	_, err := c.client.DeleteSnapshotWithContext(ctx, &ec2.DeleteSnapshotInput{
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

// DeleteSubnet deletes a Subnet with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DeleteSubnet(ctx context.Context, id string) error {
	logger := c.logger.WithField("Subnet", id)
	_, err := c.client.DeleteSubnetWithContext(ctx, &ec2.DeleteSubnetInput{
		SubnetId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidSubnetID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return err
}

// DescribeInstancesPages runs `fn` for each page of Instances
func (c *EC2Client) DescribeInstancesPages(ctx context.Context, filters []*ec2.Filter, fn func(*ec2.DescribeInstancesOutput, bool) bool) error {
	err := c.client.DescribeInstancesPagesWithContext(ctx, &ec2.DescribeInstancesInput{Filters: filters}, fn)
	if err != nil {
		return errors.Wrap(err, "get ec2 instances")
		// logger.Info(err)
	}
	return nil
}

// DescribeInstances returns EC2 Instance with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DescribeInstances(ctx context.Context, id string) (*ec2.DescribeInstancesOutput, error) {
	response, err := c.client.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidInstanceID.NotFound" {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// DescribeNatGatewaysPages runs `fn` for each page of Nat Gateways
func (c *EC2Client) DescribeNatGatewaysPages(ctx context.Context, vpc string, fn func(*ec2.DescribeNatGatewaysOutput, bool) bool) error {
	err := c.client.DescribeNatGatewaysPagesWithContext(
		ctx,
		&ec2.DescribeNatGatewaysInput{
			Filter: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{aws.String(vpc)},
				},
			},
		},
		fn,
	)
	return err
}

// DescribeRouteTables returns Route Tables with ARN=`id`. It ignores NotFound objects
func (c *EC2Client) DescribeRouteTables(ctx context.Context, id string) (*ec2.DescribeRouteTablesOutput, error) {
	response, err := c.client.DescribeRouteTablesWithContext(ctx, &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidRouteTableID.NotFound" {
			return nil, nil
		}
		return nil, err
	}

	return response, nil
}

// DescribeRouteTablesPages runs `fn` for each page of Route Tables
func (c *EC2Client) DescribeRouteTablesPages(ctx context.Context, vpc string, fn func(*ec2.DescribeRouteTablesOutput, bool) bool) error {
	err := c.client.DescribeRouteTablesPagesWithContext(
		ctx,
		&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{aws.String(vpc)},
				},
			},
		},
		fn,
	)
	return err
}

// DescribeSecurityGroupsPages runs `fn` for each page of Security Groups
func (c *EC2Client) DescribeSecurityGroupsPages(ctx context.Context, vpc string, fn func(*ec2.DescribeSecurityGroupsOutput, bool) bool) error {
	err := c.client.DescribeSecurityGroupsPagesWithContext(
		ctx,
		&ec2.DescribeSecurityGroupsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{aws.String(vpc)},
				},
			},
		},
		fn,
	)
	return err
}

// DescribeSubnetsPages runs `fn` for each page of Subnets
func (c *EC2Client) DescribeSubnetsPages(ctx context.Context, vpc string, fn func(*ec2.DescribeSubnetsOutput, bool) bool) error {
	err := c.client.DescribeSubnetsPagesWithContext(
		ctx,
		&ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{aws.String(vpc)},
				},
			},
		},
		fn,
	)
	return err
}

// DescribeNetworkInterfacesPages runs `fn` for each page of Network Interfaces
func (c *EC2Client) DescribeNetworkInterfacesPages(ctx context.Context, vpc string, fn func(*ec2.DescribeNetworkInterfacesOutput, bool) bool) error {
	err := c.client.DescribeNetworkInterfacesPagesWithContext(
		ctx,
		&ec2.DescribeNetworkInterfacesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{aws.String(vpc)},
				},
			},
		},
		fn,
	)
	return err
}

// DescribeVpcEndpoints returns VPC Endpoint with ARN=`id`
func (c *EC2Client) DescribeVpcEndpoints(ctx context.Context, vpc string) (*ec2.DescribeVpcEndpointsOutput, error) {
	response, err := c.client.DescribeVpcEndpointsWithContext(ctx, &ec2.DescribeVpcEndpointsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpc)},
			},
		},
	})
	return response, err
}

// DisassociateRouteTable disassociates a Route Table with ARN=`id`
func (c *EC2Client) DisassociateRouteTable(ctx context.Context, id string) error {
	_, err := c.client.DisassociateRouteTableWithContext(ctx, &ec2.DisassociateRouteTableInput{
		AssociationId: aws.String(id),
	})
	return err
}

// RevokeSecurityGroupEgress revokes Security Group Egress
func (c *EC2Client) RevokeSecurityGroupEgress(ctx context.Context, group *ec2.SecurityGroup) error {
	_, err := c.client.RevokeSecurityGroupEgressWithContext(ctx, &ec2.RevokeSecurityGroupEgressInput{
		GroupId:       group.GroupId,
		IpPermissions: group.IpPermissionsEgress,
	})
	return err
}

// RevokeSecurityGroupIngress revokes Security Group Ingress
func (c *EC2Client) RevokeSecurityGroupIngress(ctx context.Context, group *ec2.SecurityGroup) error {
	_, err := c.client.RevokeSecurityGroupIngressWithContext(ctx, &ec2.RevokeSecurityGroupIngressInput{
		GroupId:       group.GroupId,
		IpPermissions: group.IpPermissions,
	})
	return err
}

// TerminateInstances terminates Instance with ARN=`id`
func (c *EC2Client) TerminateInstances(ctx context.Context, id *string) error {
	_, err := c.client.TerminateInstancesWithContext(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []*string{id},
	})

	return err
}
