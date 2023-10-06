package aws

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

func createVPCResources(logger *logrus.Logger, session *session.Session, vpcInput *CreateInfraOptions) error {
	ec2Client := ec2.New(session)

	vpcID, err := vpcInput.createVPC(logger, ec2Client)
	if err != nil {
		return err
	}
	vpcInput.vpcID = vpcID

	if err = vpcInput.CreateDHCPOptions(logger, ec2Client, vpcID); err != nil {
		return err
	}

	igwID, err := vpcInput.CreateInternetGateway(logger, ec2Client, vpcID)
	if err != nil {
		return err
	}

	_, err = vpcInput.CreateWorkerSecurityGroup(logger, ec2Client, vpcID)
	if err != nil {
		return err
	}

	// Per zone resources
	// TODO(alberto): Parameterize this.
	basePrivateSubnetCIDR := "10.0.128.0/20"
	basePublicSubnetCIDR := "10.0.0.0/20"
	var endpointRouteTableIds []*string
	var publicSubnetIDs []string
	var privateSubnetIDs []string
	_, privateNetwork, err := net.ParseCIDR(basePrivateSubnetCIDR)
	if err != nil {
		return err
	}
	_, publicNetwork, err := net.ParseCIDR(basePublicSubnetCIDR)
	if err != nil {
		return err
	}
	for _, zone := range vpcInput.Zones {
		privateSubnetID, err := vpcInput.CreatePrivateSubnet(logger, ec2Client, vpcID, zone, privateNetwork.String())
		if err != nil {
			return err
		}
		publicSubnetID, err := vpcInput.CreatePublicSubnet(logger, ec2Client, vpcID, zone, publicNetwork.String())
		if err != nil {
			return err
		}
		var natGatewayID string
		publicSubnetIDs = append(publicSubnetIDs, publicSubnetID)
		privateSubnetIDs = append(privateSubnetIDs, privateSubnetID)

		if !vpcInput.EnableProxy {
			natGatewayID, err = vpcInput.CreateNATGateway(logger, ec2Client, publicSubnetID, zone)
			if err != nil {
				return err
			}
		}
		privateRouteTable, err := vpcInput.CreatePrivateRouteTable(logger, ec2Client, vpcID, natGatewayID, privateSubnetID, zone)
		if err != nil {
			return err
		}
		endpointRouteTableIds = append(endpointRouteTableIds, aws.String(privateRouteTable))

		//result.Zones = append(result.Zones, &CreateInfraOutputZone{
		//	Name:     zone,
		//	SubnetID: privateSubnetID,
		//})

		// increment each subnet by /20
		privateNetwork.IP[2] = privateNetwork.IP[2] + 16
		publicNetwork.IP[2] = publicNetwork.IP[2] + 16
	}
	vpcInput.publicSubnetIDs = publicSubnetIDs
	vpcInput.privateSubnetIDs = privateSubnetIDs

	publicRouteTable, err := vpcInput.CreatePublicRouteTable(logger, ec2Client, vpcID, igwID, publicSubnetIDs)
	if err != nil {
		return err
	}
	endpointRouteTableIds = append(endpointRouteTableIds, aws.String(publicRouteTable))
	err = vpcInput.CreateVPCS3Endpoint(logger, ec2Client, vpcID, endpointRouteTableIds)
	if err != nil {
		return err
	}

	err = vpcInput.CreateLoadBalancers(logger, session, vpcID, privateSubnetIDs, publicSubnetIDs, vpcInput.public)
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateInfraOptions) createVPC(l *logrus.Logger, client ec2iface.EC2API) (string, error) {
	// TODO(alberto): pass this from input.
	defaultCIDRBlock := "10.0.0.0/16"
	vpcName := fmt.Sprintf("%s-vpc", o.InfraID)
	vpcID, err := o.existingVPC(client, vpcName)
	if err != nil {
		return "", err
	}
	if len(vpcID) == 0 {
		createResult, err := client.CreateVpc(&ec2.CreateVpcInput{
			CidrBlock:         aws.String(defaultCIDRBlock),
			TagSpecifications: o.ec2TagSpecifications("vpc", vpcName),
		})
		if err != nil {
			return "", fmt.Errorf("failed to create VPC: %w", err)
		}
		vpcID = aws.StringValue(createResult.Vpc.VpcId)
		l.Info("Created VPC", "id", vpcID)
	} else {
		l.Info("Found existing VPC", "id", vpcID)
	}
	_, err = client.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		VpcId:            aws.String(vpcID),
		EnableDnsSupport: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		return "", fmt.Errorf("failed to modify VPC attributes: %w", err)
	}
	l.Info("Enabled DNS support on VPC", "id", vpcID)
	_, err = client.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		VpcId:              aws.String(vpcID),
		EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		return "", fmt.Errorf("failed to modify VPC attributes: %w", err)
	}
	l.Info("Enabled DNS hostnames on VPC", "id", vpcID)
	return vpcID, nil
}

func (o *CreateInfraOptions) existingVPC(client ec2iface.EC2API, vpcName string) (string, error) {
	var vpcID string
	result, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{Filters: o.ec2Filters(vpcName)})
	if err != nil {
		return "", fmt.Errorf("cannot list vpcs: %w", err)
	}
	for _, vpc := range result.Vpcs {
		vpcID = aws.StringValue(vpc.VpcId)
		break
	}
	return vpcID, nil
}

func (o *CreateInfraOptions) ec2Filters(name string) []*ec2.Filter {
	filters := []*ec2.Filter{
		{
			Name:   aws.String(fmt.Sprintf("tag:%s", clusterTag(o.InfraID))),
			Values: []*string{aws.String(clusterTagValue)},
		},
	}
	if len(name) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String(name)},
		})
	}
	return filters
}

func (o *CreateInfraOptions) ec2TagSpecifications(resourceType, name string) []*ec2.TagSpecification {
	return []*ec2.TagSpecification{
		{
			ResourceType: aws.String(resourceType),
			Tags:         append(ec2Tags(o.InfraID, name), o.additionalEC2Tags...),
		},
	}
}

func clusterTag(infraID string) string {
	return fmt.Sprintf("kubernetes.io/cluster/%s", infraID)
}

func ec2Tags(infraID, name string) []*ec2.Tag {
	tags := []*ec2.Tag{
		{
			Key:   aws.String(clusterTag(infraID)),
			Value: aws.String(clusterTagValue),
		},
	}
	if len(name) > 0 {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String("Name"),
			Value: aws.String(name),
		})
	}
	return tags

}

func (o *CreateInfraOptions) CreateDHCPOptions(l *logrus.Logger, client ec2iface.EC2API, vpcID string) error {
	domainName := "ec2.internal"
	if o.Region != "us-east-1" {
		domainName = fmt.Sprintf("%s.compute.internal", o.Region)
	}
	optID, err := o.existingDHCPOptions(client)
	if err != nil {
		return err
	}
	if len(optID) == 0 {
		result, err := client.CreateDhcpOptions(&ec2.CreateDhcpOptionsInput{
			DhcpConfigurations: []*ec2.NewDhcpConfiguration{
				{
					Key:    aws.String("domain-name"),
					Values: []*string{aws.String(domainName)},
				},
				{
					Key:    aws.String("domain-name-servers"),
					Values: []*string{aws.String("AmazonProvidedDNS")},
				},
			},
			TagSpecifications: o.ec2TagSpecifications("dhcp-options", ""),
		})
		if err != nil {
			return fmt.Errorf("cannot create dhcp-options: %w", err)
		}
		optID = aws.StringValue(result.DhcpOptions.DhcpOptionsId)
		l.Info("Created DHCP options", "id", optID)
	} else {
		l.Info("Found existing DHCP options", "id", optID)
	}
	_, err = client.AssociateDhcpOptions(&ec2.AssociateDhcpOptionsInput{
		DhcpOptionsId: aws.String(optID),
		VpcId:         aws.String(vpcID),
	})
	if err != nil {
		return fmt.Errorf("cannot associate dhcp-options to VPC: %w", err)
	}
	l.Info("Associated DHCP options with VPC", "vpc", vpcID, "dhcp options", optID)
	return nil
}

func (o *CreateInfraOptions) existingDHCPOptions(client ec2iface.EC2API) (string, error) {
	var optID string
	result, err := client.DescribeDhcpOptions(&ec2.DescribeDhcpOptionsInput{Filters: o.ec2Filters("")})
	if err != nil {
		return "", fmt.Errorf("cannot list dhcp options: %w", err)
	}
	for _, opt := range result.DhcpOptions {
		optID = aws.StringValue(opt.DhcpOptionsId)
		break
	}
	return optID, nil
}

func (o *CreateInfraOptions) CreateInternetGateway(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (string, error) {
	gatewayName := fmt.Sprintf("%s-igw", o.InfraID)
	igw, err := o.existingInternetGateway(client, gatewayName)
	if err != nil {
		return "", err
	}
	if igw == nil {
		result, err := client.CreateInternetGateway(&ec2.CreateInternetGatewayInput{
			TagSpecifications: o.ec2TagSpecifications("internet-gateway", fmt.Sprintf("%s-igw", o.InfraID)),
		})
		if err != nil {
			return "", fmt.Errorf("cannot create internet gateway: %w", err)
		}
		igw = result.InternetGateway
		l.Info("Created internet gateway", "id", aws.StringValue(igw.InternetGatewayId))
	} else {
		l.Info("Found existing internet gateway", "id", aws.StringValue(igw.InternetGatewayId))
	}
	attached := false
	for _, attachment := range igw.Attachments {
		if aws.StringValue(attachment.VpcId) == vpcID {
			attached = true
			break
		}
	}
	if !attached {
		_, err = client.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
			InternetGatewayId: igw.InternetGatewayId,
			VpcId:             aws.String(vpcID),
		})
		if err != nil {
			return "", fmt.Errorf("cannot attach internet gateway to vpc: %w", err)
		}
		l.Info("Attached internet gateway to VPC", "internet gateway", aws.StringValue(igw.InternetGatewayId), "vpc", vpcID)
	}
	return aws.StringValue(igw.InternetGatewayId), nil
}

func (o *CreateInfraOptions) existingInternetGateway(client ec2iface.EC2API, name string) (*ec2.InternetGateway, error) {
	result, err := client.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{Filters: o.ec2Filters(name)})
	if err != nil {
		return nil, fmt.Errorf("cannot list internet gateways: %w", err)
	}
	for _, igw := range result.InternetGateways {
		return igw, nil
	}
	return nil, nil
}

func (o *CreateInfraOptions) CreateWorkerSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (string, error) {
	backoff := wait.Backoff{
		Steps:    10,
		Duration: 3 * time.Second,
		Factor:   1.0,
		Jitter:   0.1,
	}
	groupName := fmt.Sprintf("%s-worker-sg", o.InfraID)
	securityGroup, err := o.existingSecurityGroup(client, groupName)
	if err != nil {
		return "", err
	}
	if securityGroup == nil {
		result, err := client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
			GroupName:         aws.String(groupName),
			Description:       aws.String("worker security group"),
			VpcId:             aws.String(vpcID),
			TagSpecifications: o.ec2TagSpecifications("security-group", groupName),
		})
		if err != nil {
			return "", fmt.Errorf("cannot create worker security group: %w", err)
		}
		var sgResult *ec2.DescribeSecurityGroupsOutput
		err = retry.OnError(backoff, func(error) bool { return true }, func() error {
			var err error
			sgResult, err = client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
				GroupIds: []*string{result.GroupId},
			})
			if err != nil || len(sgResult.SecurityGroups) == 0 {
				return fmt.Errorf("not found yet")
			}
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("cannot find security group that was just created (%s)", aws.StringValue(result.GroupId))
		}
		securityGroup = sgResult.SecurityGroups[0]
		l.Info("Created security group", "name", groupName, "id", aws.StringValue(securityGroup.GroupId))
	} else {
		l.Info("Found existing security group", "name", groupName, "id", aws.StringValue(securityGroup.GroupId))
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	//sgUserID := aws.StringValue(securityGroup.OwnerId)
	egressPermissions := DefaultWorkerSGEgressRules()
	ingressPermissions := DefaultWorkerSGIngressRules()

	var egressToAuthorize []*ec2.IpPermission
	var ingressToAuthorize []*ec2.IpPermission

	for _, permission := range egressPermissions {
		if !includesPermission(securityGroup.IpPermissionsEgress, permission) {
			egressToAuthorize = append(egressToAuthorize, permission)
		}
	}

	for _, permission := range ingressPermissions {
		if !includesPermission(securityGroup.IpPermissions, permission) {
			ingressToAuthorize = append(ingressToAuthorize, permission)
		}
	}

	const duplicatePermissionErrorCode = "InvalidPermission.Duplicate"
	if len(egressToAuthorize) > 0 {
		err = retry.OnError(backoff, func(error) bool { return true }, func() error {
			_, err := client.AuthorizeSecurityGroupEgress(&ec2.AuthorizeSecurityGroupEgressInput{
				GroupId:       aws.String(securityGroupID),
				IpPermissions: egressToAuthorize,
			})
			var awsErr awserr.Error
			if err != nil {
				if errors.As(err, &awsErr) {
					// only return an error if the permission has not already been set
					if awsErr.Code() != duplicatePermissionErrorCode {
						return fmt.Errorf("cannot apply security group egress permissions: %w", err)
					}
				}
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		l.Info("Authorized egress rules on security group", "id", securityGroupID)
	}
	if len(ingressToAuthorize) > 0 {
		err = retry.OnError(backoff, func(error) bool { return true }, func() error {
			_, err := client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
				GroupId:       aws.String(securityGroupID),
				IpPermissions: ingressToAuthorize,
			})
			var awsErr awserr.Error
			if err != nil {
				if errors.As(err, &awsErr) {
					// only return an error if the permission has not already been set
					if awsErr.Code() != duplicatePermissionErrorCode {
						return fmt.Errorf("cannot apply security group ingress permissions: %w", err)
					}
				}
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		l.Info("Authorized ingress rules on security group", "id", securityGroupID)
	}

	o.allowAllSecurityGroupID = securityGroupID
	return securityGroupID, nil
}

func (o *CreateInfraOptions) existingSecurityGroup(client ec2iface.EC2API, name string) (*ec2.SecurityGroup, error) {
	result, err := client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{Filters: o.ec2Filters(name)})
	if err != nil {
		return nil, fmt.Errorf("cannot list security groups: %w", err)
	}
	for _, sg := range result.SecurityGroups {
		return sg, nil
	}
	return nil, nil
}

func includesPermission(list []*ec2.IpPermission, permission *ec2.IpPermission) bool {
	for _, p := range list {
		if samePermission(p, permission) {
			return true
		}
	}
	return false
}

func samePermission(a, b *ec2.IpPermission) bool {
	if a == nil || b == nil {
		return false
	}
	if a.String() == b.String() {
		return true
	}
	return false
}

func DefaultWorkerSGEgressRules() []*ec2.IpPermission {
	return []*ec2.IpPermission{
		{
			IpProtocol: aws.String("-1"),
			IpRanges: []*ec2.IpRange{
				{
					CidrIp: aws.String("0.0.0.0/0"),
				},
			},
		},
	}
}

// DefaultWorkerSGIngressRules
// TODO(alberto): scope this down to granular perms.
func DefaultWorkerSGIngressRules() []*ec2.IpPermission {
	return []*ec2.IpPermission{
		{
			IpProtocol: aws.String("-1"),
			IpRanges: []*ec2.IpRange{
				{
					CidrIp: aws.String("0.0.0.0/0"),
				},
			},
		},
	}
}

func (o *CreateInfraOptions) CreatePrivateSubnet(l *logrus.Logger, client ec2iface.EC2API, vpcID string, zone string, cidr string) (string, error) {
	return o.CreateSubnet(l, client, vpcID, zone, cidr, fmt.Sprintf("%s-private-%s", o.InfraID, zone), tagNameSubnetInternalELB)
}

func (o *CreateInfraOptions) CreatePublicSubnet(l *logrus.Logger, client ec2iface.EC2API, vpcID string, zone string, cidr string) (string, error) {
	return o.CreateSubnet(l, client, vpcID, zone, cidr, fmt.Sprintf("%s-public-%s", o.InfraID, zone), tagNameSubnetPublicELB)
}

func (o *CreateInfraOptions) CreateSubnet(l *logrus.Logger, client ec2iface.EC2API, vpcID, zone, cidr, name, scopeTag string) (string, error) {
	subnetID, err := o.existingSubnet(client, name)
	if err != nil {
		return "", err
	}
	if len(subnetID) > 0 {
		l.Info("Found existing subnet", "name", name, "id", subnetID)
		return subnetID, nil
	}
	tagSpec := o.ec2TagSpecifications("subnet", name)
	tagSpec[0].Tags = append(tagSpec[0].Tags, &ec2.Tag{
		Key:   aws.String(scopeTag),
		Value: aws.String("true"),
	})

	result, err := client.CreateSubnet(&ec2.CreateSubnetInput{
		AvailabilityZone:  aws.String(zone),
		VpcId:             aws.String(vpcID),
		CidrBlock:         aws.String(cidr),
		TagSpecifications: tagSpec,
	})
	if err != nil {
		return "", fmt.Errorf("cannot create public subnet: %w", err)
	}
	backoff := wait.Backoff{
		Steps:    10,
		Duration: 3 * time.Second,
		Factor:   1.0,
		Jitter:   0.1,
	}
	var subnetResult *ec2.DescribeSubnetsOutput
	err = retry.OnError(backoff, func(error) bool { return true }, func() error {
		var err error
		subnetResult, err = client.DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: []*string{result.Subnet.SubnetId},
		})
		if err != nil || len(subnetResult.Subnets) == 0 {
			return fmt.Errorf("not found yet")
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot find subnet that was just created (%s)", aws.StringValue(result.Subnet.SubnetId))
	}
	subnetID = aws.StringValue(result.Subnet.SubnetId)
	l.Info("Created subnet", "name", name, "id", subnetID)
	return subnetID, nil
}

func (o *CreateInfraOptions) existingSubnet(client ec2iface.EC2API, name string) (string, error) {
	var subnetID string
	result, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{Filters: o.ec2Filters(name)})
	if err != nil {
		return "", fmt.Errorf("cannot list subnets: %w", err)
	}
	for _, subnet := range result.Subnets {
		subnetID = aws.StringValue(subnet.SubnetId)
		break
	}
	return subnetID, nil
}

const (
	invalidNATGatewayError   = "InvalidNatGatewayID.NotFound"
	invalidRouteTableID      = "InvalidRouteTableId.NotFound"
	invalidElasticIPNotFound = "InvalidElasticIpID.NotFound"
	invalidSubnet            = "InvalidSubnet"

	// tagNameSubnetInternalELB is the tag name used on a subnet to designate that
	// it should be used for internal ELBs
	tagNameSubnetInternalELB = "kubernetes.io/role/internal-elb"

	// tagNameSubnetPublicELB is the tag name used on a subnet to designate that
	// it should be used for internet ELBs
	tagNameSubnetPublicELB = "kubernetes.io/role/elb"
)

var (
	retryBackoff = wait.Backoff{
		Steps:    5,
		Duration: 3 * time.Second,
		Factor:   3.0,
		Jitter:   0.1,
	}
)

func (o *CreateInfraOptions) CreateNATGateway(l *logrus.Logger, client ec2iface.EC2API, publicSubnetID, availabilityZone string) (string, error) {
	natGatewayName := fmt.Sprintf("%s-nat-%s", o.InfraID, availabilityZone)
	natGateway, _ := o.existingNATGateway(client, natGatewayName)
	if natGateway != nil {
		l.Info("Found existing NAT gateway", "id", aws.StringValue(natGateway.NatGatewayId))
		return *natGateway.NatGatewayId, nil
	}

	eipResult, err := client.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})
	if err != nil {
		return "", fmt.Errorf("cannot allocate EIP for NAT gateway: %w", err)
	}
	allocationID := aws.StringValue(eipResult.AllocationId)
	l.Info("Created elastic IP for NAT gateway", "id", allocationID)

	// NOTE: there's a potential to leak EIP addresses if the following tag operation fails, since we have no way of
	// recognizing the EIP as belonging to the cluster
	isRetriable := func(err error) bool {
		if awsErr, ok := err.(awserr.Error); ok {
			return strings.EqualFold(awsErr.Code(), invalidElasticIPNotFound)
		}
		return false
	}
	err = retry.OnError(retryBackoff, isRetriable, func() error {
		_, err = client.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{aws.String(allocationID)},
			Tags:      append(ec2Tags(o.InfraID, fmt.Sprintf("%s-eip-%s", o.InfraID, availabilityZone)), o.additionalEC2Tags...),
		})
		return err
	})
	if err != nil {
		return "", fmt.Errorf("cannot tag NAT gateway EIP: %w", err)
	}

	isNATGatewayRetriable := func(err error) bool {
		if awsErr, ok := err.(awserr.Error); ok {
			return strings.EqualFold(awsErr.Code(), invalidSubnet) ||
				strings.EqualFold(awsErr.Code(), invalidElasticIPNotFound)
		}
		return false
	}
	err = retry.OnError(retryBackoff, isNATGatewayRetriable, func() error {
		gatewayResult, err := client.CreateNatGateway(&ec2.CreateNatGatewayInput{
			AllocationId:      aws.String(allocationID),
			SubnetId:          aws.String(publicSubnetID),
			TagSpecifications: o.ec2TagSpecifications("natgateway", natGatewayName),
		})
		if err != nil {
			return err
		}
		natGateway = gatewayResult.NatGateway
		l.Info("Created NAT gateway", "id", aws.StringValue(natGateway.NatGatewayId))
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot create NAT gateway: %w", err)
	}

	natGatewayID := aws.StringValue(natGateway.NatGatewayId)
	return natGatewayID, nil
}

func (o *CreateInfraOptions) existingNATGateway(client ec2iface.EC2API, name string) (*ec2.NatGateway, error) {
	result, err := client.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{Filter: o.ec2Filters(name)})
	if err != nil {
		return nil, fmt.Errorf("cannot list NAT gateways: %w", err)
	}
	for _, gateway := range result.NatGateways {
		state := aws.StringValue(gateway.State)
		if state == "deleted" || state == "deleting" || state == "failed" {
			continue
		}
		return gateway, nil
	}
	return nil, nil
}

func (o *CreateInfraOptions) CreatePrivateRouteTable(l *logrus.Logger, client ec2iface.EC2API, vpcID, natGatewayID, subnetID, zone string) (string, error) {
	tableName := fmt.Sprintf("%s-private-%s", o.InfraID, zone)
	routeTable, err := o.existingRouteTable(l, client, tableName)
	if err != nil {
		return "", err
	}
	if routeTable == nil {
		routeTable, err = o.createRouteTable(l, client, vpcID, tableName)
		if err != nil {
			return "", err
		}
	}

	// Everything below this is only needed if direct internet access is used
	if o.EnableProxy {
		return aws.StringValue(routeTable.RouteTableId), nil
	}

	if !o.hasNATGatewayRoute(routeTable, natGatewayID) {
		isRetriable := func(err error) bool {
			if awsErr, ok := err.(awserr.Error); ok {
				return strings.EqualFold(awsErr.Code(), invalidNATGatewayError)
			}
			return false
		}
		err = retry.OnError(retryBackoff, isRetriable, func() error {
			_, err = client.CreateRoute(&ec2.CreateRouteInput{
				RouteTableId:         routeTable.RouteTableId,
				NatGatewayId:         aws.String(natGatewayID),
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
			})
			return err
		})
		if err != nil {
			return "", fmt.Errorf("cannot create nat gateway route in private route table: %w", err)
		}
		l.Info("Created route to NAT gateway", "route table", aws.StringValue(routeTable.RouteTableId), "nat gateway", natGatewayID)
	} else {
		l.Info("Found existing route to NAT gateway", "route table", aws.StringValue(routeTable.RouteTableId), "nat gateway", natGatewayID)
	}
	if !o.hasAssociatedSubnet(routeTable, subnetID) {
		_, err = client.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: routeTable.RouteTableId,
			SubnetId:     aws.String(subnetID),
		})
		if err != nil {
			return "", fmt.Errorf("cannot associate private route table with subnet: %w", err)
		}
		l.Info("Associated subnet with route table", "route table", aws.StringValue(routeTable.RouteTableId), "subnet", subnetID)
	} else {
		l.Info("Subnet already associated with route table", "route table", aws.StringValue(routeTable.RouteTableId), "subnet", subnetID)
	}
	return aws.StringValue(routeTable.RouteTableId), nil
}

func (o *CreateInfraOptions) CreatePublicRouteTable(l *logrus.Logger, client ec2iface.EC2API, vpcID, igwID string, subnetIDs []string) (string, error) {
	tableName := fmt.Sprintf("%s-public", o.InfraID)
	routeTable, err := o.existingRouteTable(l, client, tableName)
	if err != nil {
		return "", err
	}
	if routeTable == nil {
		routeTable, err = o.createRouteTable(l, client, vpcID, tableName)
		if err != nil {
			return "", err
		}
	}
	tableID := aws.StringValue(routeTable.RouteTableId)
	// Replace the VPC's main route table
	routeTableInfo, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
			{
				Name:   aws.String("association.main"),
				Values: []*string{aws.String("true")},
			},
		},
	})
	if err != nil {
		return "", err
	}
	if len(routeTableInfo.RouteTables) == 0 {
		return "", fmt.Errorf("no route tables associated with the vpc")
	}
	// Replace route table association only if it's not the associated route table already
	if aws.StringValue(routeTableInfo.RouteTables[0].RouteTableId) != tableID {
		var associationID string
		for _, assoc := range routeTableInfo.RouteTables[0].Associations {
			if aws.BoolValue(assoc.Main) {
				associationID = aws.StringValue(assoc.RouteTableAssociationId)
				break
			}
		}
		_, err = client.ReplaceRouteTableAssociation(&ec2.ReplaceRouteTableAssociationInput{
			RouteTableId:  aws.String(tableID),
			AssociationId: aws.String(associationID),
		})
		if err != nil {
			return "", fmt.Errorf("cannot set vpc main route table: %w", err)
		}
		l.Info("Set main VPC route table", "route table", tableID, "vpc", vpcID)
	}

	// Create route to internet gateway
	if !o.hasInternetGatewayRoute(routeTable, igwID) {
		_, err = client.CreateRoute(&ec2.CreateRouteInput{
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
			RouteTableId:         aws.String(tableID),
			GatewayId:            aws.String(igwID),
		})
		if err != nil {
			return "", fmt.Errorf("cannot create route to internet gateway: %w", err)
		}
		l.Info("Created route to internet gateway", "route table", tableID, "internet gateway", igwID)
	} else {
		l.Info("Found existing route to internet gateway", "route table", tableID, "internet gateway", igwID)
	}

	// Associate the route table with the public subnet ID
	for _, subnetID := range subnetIDs {
		if !o.hasAssociatedSubnet(routeTable, subnetID) {
			_, err = client.AssociateRouteTable(&ec2.AssociateRouteTableInput{
				RouteTableId: aws.String(tableID),
				SubnetId:     aws.String(subnetID),
			})
			if err != nil {
				return "", fmt.Errorf("cannot associate private route table with subnet: %w", err)
			}
			l.Info("Associated route table with subnet", "route table", tableID, "subnet", subnetID)
		} else {
			l.Info("Found existing association between route table and subnet", "route table", tableID, "subnet", subnetID)
		}
	}
	return tableID, nil
}

func (o *CreateInfraOptions) createRouteTable(l *logrus.Logger, client ec2iface.EC2API, vpcID, name string) (*ec2.RouteTable, error) {
	result, err := client.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId:             aws.String(vpcID),
		TagSpecifications: o.ec2TagSpecifications("route-table", name),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create route table: %w", err)
	}
	l.Info("Created route table", "name", name, "id", aws.StringValue(result.RouteTable.RouteTableId))
	return result.RouteTable, nil
}

func (o *CreateInfraOptions) existingRouteTable(l *logrus.Logger, client ec2iface.EC2API, name string) (*ec2.RouteTable, error) {
	result, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{Filters: o.ec2Filters(name)})
	if err != nil {
		return nil, fmt.Errorf("cannot list route tables: %w", err)
	}
	if len(result.RouteTables) > 0 {
		l.Info("Found existing route table", "name", name, "id", aws.StringValue(result.RouteTables[0].RouteTableId))
		return result.RouteTables[0], nil
	}
	return nil, nil
}

func (o *CreateInfraOptions) hasNATGatewayRoute(table *ec2.RouteTable, natGatewayID string) bool {
	for _, route := range table.Routes {
		if aws.StringValue(route.NatGatewayId) == natGatewayID &&
			aws.StringValue(route.DestinationCidrBlock) == "0.0.0.0/0" {
			return true
		}
	}
	return false
}

func (o *CreateInfraOptions) hasInternetGatewayRoute(table *ec2.RouteTable, igwID string) bool {
	for _, route := range table.Routes {
		if aws.StringValue(route.GatewayId) == igwID &&
			aws.StringValue(route.DestinationCidrBlock) == "0.0.0.0/0" {
			return true
		}
	}
	return false
}

func (o *CreateInfraOptions) hasAssociatedSubnet(table *ec2.RouteTable, subnetID string) bool {
	for _, assoc := range table.Associations {
		if aws.StringValue(assoc.RouteTableId) == subnetID {
			return true
		}
	}
	return false
}

// TODO(alberto): refactor into a createLB() func that receive input.
func (o *CreateInfraOptions) CreateLoadBalancers(l *logrus.Logger, session *session.Session, vpcID string, privateSubnets []string, publicSubnets []string, external bool) error {
	elbClient := elbv2.New(session)

	// Create internal LB.
	internalLBName := fmt.Sprintf("%s-int", o.InfraID)

	// Check if the internal load balancer already exists.
	describeLBInput := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{&internalLBName},
	}
	describeLBOutput, err := elbClient.DescribeLoadBalancers(describeLBInput)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != "LoadBalancerNotFound" {
				return fmt.Errorf("failed to describe lb: %w", awsErr)
			}
		} else {
			return fmt.Errorf("failed to describe lb: %w", err)
		}
	}
	var internalLBARN string
	if len(describeLBOutput.LoadBalancers) == 0 {
		// Convert subnets to *string as expected by aws sdk.
		privateSubnetsPointers := make([]*string, 0)
		for i := range privateSubnets {
			privateSubnetsPointers = append(privateSubnetsPointers, &privateSubnets[i])
		}

		createInternalLBInput := &elbv2.CreateLoadBalancerInput{
			CustomerOwnedIpv4Pool: nil,
			IpAddressType:         nil,
			Name:                  aws.String(internalLBName),
			Scheme:                aws.String("internal"),
			SecurityGroups:        nil,
			SubnetMappings:        nil,
			Subnets:               privateSubnetsPointers,
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String(clusterTag(o.InfraID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(internalLBName),
				},
			},
			Type: aws.String("network"),
		}

		internalLBOutput, err := elbClient.CreateLoadBalancer(createInternalLBInput)
		if err != nil {
			return fmt.Errorf("error creating internal load balancer: %w", err)
		}
		internalLBARN = *internalLBOutput.LoadBalancers[0].LoadBalancerArn
		attrInput := &elbv2.ModifyLoadBalancerAttributesInput{
			LoadBalancerArn: aws.String(internalLBARN),
			Attributes: []*elbv2.LoadBalancerAttribute{
				{
					Key:   aws.String("load_balancing.cross_zone.enabled"),
					Value: aws.String("true"),
				},
			},
		}
		_, err = elbClient.ModifyLoadBalancerAttributes(attrInput)
		if err != nil {
			return fmt.Errorf("error modifying load balancer attributes: %w", err)
		}
		o.LoadBalancers.Internal.ZoneID = *internalLBOutput.LoadBalancers[0].CanonicalHostedZoneId
		o.LoadBalancers.Internal.DNSName = *internalLBOutput.LoadBalancers[0].DNSName
		l.Infof("Internal Load Balancer created: %v", *internalLBOutput.LoadBalancers[0].DNSName)
	} else {
		internalLBARN = *describeLBOutput.LoadBalancers[0].LoadBalancerArn
		o.LoadBalancers.Internal.ZoneID = *describeLBOutput.LoadBalancers[0].CanonicalHostedZoneId
		o.LoadBalancers.Internal.DNSName = *describeLBOutput.LoadBalancers[0].DNSName
		l.Infof("Internal Load Balancer already exists: %v", *describeLBOutput.LoadBalancers[0].DNSName)
	}

	// Create InternalA TargetGroup.
	internalATGName := fmt.Sprintf("%s-aint", o.InfraID)

	// Check if the target group already exists
	describeInternalATGInput := &elbv2.DescribeTargetGroupsInput{
		Names: []*string{&internalATGName},
	}
	describeInternalATGOutput, err := elbClient.DescribeTargetGroups(describeInternalATGInput)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != "TargetGroupNotFound" {
				return fmt.Errorf("failed to describe lb: %w", awsErr)
			}
		} else {
			return fmt.Errorf("failed to describe lb: %w", err)
		}
	}
	var internalATGARN string
	if len(describeInternalATGOutput.TargetGroups) == 0 {
		createInternalTGInput := &elbv2.CreateTargetGroupInput{
			HealthCheckEnabled:         aws.Bool(true),
			HealthCheckIntervalSeconds: aws.Int64(10),
			HealthCheckPath:            aws.String("/readyz"),
			HealthCheckPort:            aws.String("6443"),
			HealthCheckProtocol:        aws.String("HTTPS"),
			HealthyThresholdCount:      aws.Int64(2),
			Name:                       aws.String(internalATGName),
			Port:                       aws.Int64(6443),
			Protocol:                   aws.String("TCP"),
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String(clusterTag(o.InfraID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(internalATGName),
				},
			},
			TargetType:              aws.String("ip"),
			UnhealthyThresholdCount: aws.Int64(2),
			VpcId:                   aws.String(vpcID),
		}
		internalTGOutput, err := elbClient.CreateTargetGroup(createInternalTGInput)
		if err != nil {
			return fmt.Errorf("error creating internal target group: %w", err)
		}
		internalATGARN = *internalTGOutput.TargetGroups[0].TargetGroupArn
		l.Infof("InternalA Target Group created: %v", *internalTGOutput.TargetGroups[0].TargetGroupArn)
	} else {
		internalATGARN = *describeInternalATGOutput.TargetGroups[0].TargetGroupArn
		l.Infof("InternalA Target Group already exists: %v", *describeInternalATGOutput.TargetGroups[0].TargetGroupArn)
	}

	// Create InternalA Listener.
	internalAListenerName := fmt.Sprintf("%s-aint", o.InfraID)
	createInternalAListenerInput := &elbv2.CreateListenerInput{
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: []*elbv2.TargetGroupTuple{
						{
							TargetGroupArn: &internalATGARN,
						},
					},
				},
			},
		},
		LoadBalancerArn: &internalLBARN,
		Port:            aws.Int64(6443),
		Protocol:        aws.String("TCP"),
		SslPolicy:       nil,
		Tags: []*elbv2.Tag{
			{
				Key:   aws.String(clusterTag(o.InfraID)),
				Value: aws.String(clusterTagValue),
			},
			{
				Key:   aws.String("Name"),
				Value: aws.String(internalAListenerName),
			},
		},
	}
	internalAListenerOutput, err := elbClient.CreateListener(createInternalAListenerInput)
	if err != nil {
		return fmt.Errorf("error creating internal listener: %w", err)
	}
	l.Infof("Internal Listener created: %v", *internalAListenerOutput.Listeners[0].ListenerArn)

	// Create InternalS TargetGroup.
	internalSTGName := fmt.Sprintf("%s-sint", o.InfraID)
	// Check if the target group already exists
	describeInternalSTGInput := &elbv2.DescribeTargetGroupsInput{
		Names: []*string{&internalSTGName},
	}
	describeInternalSTGOutput, err := elbClient.DescribeTargetGroups(describeInternalSTGInput)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != "TargetGroupNotFound" {
				return fmt.Errorf("failed to describe lb: %w", awsErr)
			}
		} else {
			return fmt.Errorf("failed to describe lb: %w", err)
		}
	}

	var internalSTGARN string
	// If the target group doesn't exist, create it
	if len(describeInternalSTGOutput.TargetGroups) == 0 {
		createSTGInput := &elbv2.CreateTargetGroupInput{
			Name:                       aws.String(internalSTGName),
			Protocol:                   aws.String("TCP"),
			Port:                       aws.Int64(int64(22623)),
			VpcId:                      aws.String(vpcID),
			TargetType:                 aws.String("ip"),
			HealthCheckProtocol:        aws.String("HTTPS"),
			HealthCheckPort:            aws.String("22623"),
			HealthCheckPath:            aws.String("/healthz"),
			HealthCheckIntervalSeconds: aws.Int64(10),
			HealthyThresholdCount:      aws.Int64(2),
			UnhealthyThresholdCount:    aws.Int64(2),
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String(clusterTag(o.InfraID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(internalSTGName),
				},
			},
		}

		createSTGOutput, err := elbClient.CreateTargetGroup(createSTGInput)
		if err != nil {
			fmt.Errorf("error creating target group: %w", err)
		}

		internalSTGARN = *createSTGOutput.TargetGroups[0].TargetGroupArn
		l.Infof("InternalS Target Group created: %v", internalSTGARN)
	} else {
		// Target group already exists
		internalSTGARN = *describeInternalSTGOutput.TargetGroups[0].TargetGroupArn
		l.Infof("InternalS Target Group already exists: %v", internalSTGARN)
	}

	internalSListenerName := fmt.Sprintf("%s-sint", o.InfraID)
	// Create a listener and associate it with the target group
	createInternalSListenerInput := &elbv2.CreateListenerInput{
		LoadBalancerArn: aws.String(internalLBARN),
		Protocol:        aws.String("TCP"),
		Port:            aws.Int64(22623),
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: []*elbv2.TargetGroupTuple{
						{
							TargetGroupArn: aws.String(internalSTGARN),
						},
					},
				},
			},
		},
		Tags: []*elbv2.Tag{
			{
				Key:   aws.String(clusterTag(o.InfraID)),
				Value: aws.String(clusterTagValue),
			},
			{
				Key:   aws.String("Name"),
				Value: aws.String(internalSListenerName),
			},
		},
	}

	_, err = elbClient.CreateListener(createInternalSListenerInput)
	if err != nil {
		fmt.Errorf("error creating listener: %v", err)
	}
	l.Infof("Internal Service Listener created: %v", *internalAListenerOutput.Listeners[0].ListenerArn)

	if !external {
		return nil
	}

	// Create external LB.
	externalLBName := fmt.Sprintf("%s-ext", o.InfraID)
	// Check if the internal load balancer already exists
	describeExternalLBInput := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{&externalLBName},
	}
	describeExternalLBOutput, err := elbClient.DescribeLoadBalancers(describeExternalLBInput)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != "LoadBalancerNotFound" {
				return fmt.Errorf("failed to describe lb: %w", awsErr)
			}
		} else {
			return fmt.Errorf("failed to describe lb: %w", err)
		}
	}

	var externalLBARN string
	if len(describeExternalLBOutput.LoadBalancers) == 0 {
		externalSubnetsPointers := make([]*string, 0)
		for i := range publicSubnets {
			externalSubnetsPointers = append(externalSubnetsPointers, &publicSubnets[i])
		}

		createExternalLBInput := &elbv2.CreateLoadBalancerInput{
			CustomerOwnedIpv4Pool: nil,
			IpAddressType:         nil,
			Name:                  aws.String(externalLBName),
			Scheme:                aws.String("internet-facing"),
			SecurityGroups:        nil,
			SubnetMappings:        nil,
			Subnets:               externalSubnetsPointers,
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String(clusterTag(o.InfraID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(externalLBName),
				},
			},
			Type: aws.String("network"),
		}

		externalLBOutput, err := elbClient.CreateLoadBalancer(createExternalLBInput)
		if err != nil {
			return fmt.Errorf("error creating external load balancer: %w", err)
		}
		externalLBARN = *externalLBOutput.LoadBalancers[0].LoadBalancerArn
		attrInput := &elbv2.ModifyLoadBalancerAttributesInput{
			LoadBalancerArn: aws.String(externalLBARN),
			Attributes: []*elbv2.LoadBalancerAttribute{
				{
					Key:   aws.String("load_balancing.cross_zone.enabled"),
					Value: aws.String("true"),
				},
			},
		}
		_, err = elbClient.ModifyLoadBalancerAttributes(attrInput)
		if err != nil {
			return fmt.Errorf("error modifying load balancer attributes: %w", err)
		}
		o.LoadBalancers.External.ZoneID = *externalLBOutput.LoadBalancers[0].CanonicalHostedZoneId
		o.LoadBalancers.External.DNSName = *externalLBOutput.LoadBalancers[0].DNSName
		l.Infof("External Load Balancer created: %v", *externalLBOutput.LoadBalancers[0].DNSName)
	} else {
		externalLBARN = *describeExternalLBOutput.LoadBalancers[0].LoadBalancerArn
		o.LoadBalancers.External.ZoneID = *describeExternalLBOutput.LoadBalancers[0].CanonicalHostedZoneId
		o.LoadBalancers.External.DNSName = *describeExternalLBOutput.LoadBalancers[0].DNSName
		l.Infof("External Load Balancer already exists: %v", *describeExternalLBOutput.LoadBalancers[0].DNSName)
	}

	// Create TargetGroup.
	externalTGName := fmt.Sprintf("%s-aext", o.InfraID)

	// Check if the target group already exists
	describeExternalTGInput := &elbv2.DescribeTargetGroupsInput{
		Names: []*string{&externalTGName},
	}
	describeExternalTGOutput, err := elbClient.DescribeTargetGroups(describeExternalTGInput)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != "TargetGroupNotFound" {
				return fmt.Errorf("failed to describe lb: %w", awsErr)
			}
		} else {
			return fmt.Errorf("failed to describe lb: %w", err)
		}
	}

	var externalTGARN string
	if len(describeExternalTGOutput.TargetGroups) == 0 {
		createExternalTGInput := &elbv2.CreateTargetGroupInput{
			HealthCheckEnabled:         aws.Bool(true),
			HealthCheckIntervalSeconds: aws.Int64(10),
			HealthCheckPath:            aws.String("/readyz"),
			HealthCheckPort:            aws.String("6443"),
			HealthCheckProtocol:        aws.String("HTTPS"),
			HealthyThresholdCount:      aws.Int64(2),
			Name:                       aws.String(externalTGName),
			Port:                       aws.Int64(6443),
			Protocol:                   aws.String("TCP"),
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String(clusterTag(o.InfraID)),
					Value: aws.String(clusterTagValue),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(externalTGName),
				},
			},
			TargetType:              aws.String("ip"),
			UnhealthyThresholdCount: aws.Int64(2),
			VpcId:                   aws.String(vpcID),
		}
		tgOutput, err := elbClient.CreateTargetGroup(createExternalTGInput)
		if err != nil {
			return fmt.Errorf("error creating internal target group: %w", err)
		}
		externalTGARN = *tgOutput.TargetGroups[0].TargetGroupArn
		l.Infof("External Target Group created: %v", *tgOutput.TargetGroups[0].TargetGroupArn)
	} else {
		externalTGARN = *describeExternalTGOutput.TargetGroups[0].TargetGroupArn
		l.Infof("External Target Group already exists: %v", *describeExternalTGOutput.TargetGroups[0].TargetGroupArn)
	}

	externalListenerName := fmt.Sprintf("%s-aext", o.InfraID)
	createExternalListenerInput := &elbv2.CreateListenerInput{
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: []*elbv2.TargetGroupTuple{
						{
							TargetGroupArn: &externalTGARN,
						},
					},
				},
			},
		},
		LoadBalancerArn: &externalLBARN,
		Port:            aws.Int64(6443),
		Protocol:        aws.String("TCP"),
		SslPolicy:       nil,
		Tags: []*elbv2.Tag{
			{
				Key:   aws.String(clusterTag(o.InfraID)),
				Value: aws.String(clusterTagValue),
			},
			{
				Key:   aws.String("Name"),
				Value: aws.String(externalListenerName),
			},
		},
	}
	externalListenerOutput, err := elbClient.CreateListener(createExternalListenerInput)
	if err != nil {
		return fmt.Errorf("error creating external listener: %w", err)
	}
	l.Infof("External Listener created: %v", *externalListenerOutput.Listeners[0].ListenerArn)

	o.targetGroupARNs = []string{externalTGARN, internalATGARN, internalSTGARN}
	return nil
}

func (o *CreateInfraOptions) existingVPCS3Endpoint(client ec2iface.EC2API) (string, error) {
	var endpointID string
	result, err := client.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{Filters: o.ec2Filters("")})
	if err != nil {
		return "", fmt.Errorf("cannot list vpc endpoints: %w", err)
	}
	for _, endpoint := range result.VpcEndpoints {
		endpointID = aws.StringValue(endpoint.VpcEndpointId)
	}
	return endpointID, nil
}

func (o *CreateInfraOptions) CreateVPCS3Endpoint(l *logrus.Logger, client ec2iface.EC2API, vpcID string, routeTableIds []*string) error {
	existingEndpoint, err := o.existingVPCS3Endpoint(client)
	if err != nil {
		return err
	}
	if len(existingEndpoint) > 0 {
		l.Info("Found existing s3 VPC endpoint", "id", existingEndpoint)
		return nil
	}
	isRetriable := func(err error) bool {
		if awsErr, ok := err.(awserr.Error); ok {
			return strings.EqualFold(awsErr.Code(), invalidRouteTableID)
		}
		return false
	}
	if err = retry.OnError(retryBackoff, isRetriable, func() error {
		result, err := client.CreateVpcEndpoint(&ec2.CreateVpcEndpointInput{
			VpcId:             aws.String(vpcID),
			ServiceName:       aws.String(fmt.Sprintf("com.amazonaws.%s.s3", o.Region)),
			RouteTableIds:     routeTableIds,
			TagSpecifications: o.ec2TagSpecifications("vpc-endpoint", ""),
		})
		if err == nil {
			l.Info("Created s3 VPC endpoint", "id", aws.StringValue(result.VpcEndpoint.VpcEndpointId))
		}
		return err
	}); err != nil {
		return fmt.Errorf("cannot create VPC S3 endpoint: %w", err)
	}
	return nil
}

// CreateInfraOptions
// TODO(alberto): this is brought form hypershift just to satisfy current functions.
// Feel free to model this as you best see fit.
type CreateInfraOptions struct {
	Region             string
	InfraID            string
	AWSCredentialsFile string
	AWSKey             string
	AWSSecretKey       string
	Name               string
	BaseDomain         string
	BaseDomainPrefix   string
	Zones              []string
	OutputFile         string
	AdditionalTags     []string
	EnableProxy        bool
	SSHKeyFile         string
	additionalEC2Tags  []*ec2.Tag

	// Additional output for consumers.
	vpcID            string
	public           bool
	publicSubnetIDs  []string
	privateSubnetIDs []string
	LoadBalancers    struct {
		Internal struct {
			DNSName string
			ZoneID  string
		}
		External struct {
			DNSName string
			ZoneID  string
		}
	}
	targetGroupARNs         []string
	allowAllSecurityGroupID string
}

// CreateInfraOutput
// TODO(alberto): this is brought form hypershift just to satisfy current functions.
// Feel free to model this as you best see fit. E.g. see dnsInput or boostrapInput.
//type CreateInfraOutput struct {
//	Region           string                   `json:"region"`
//	Zone             string                   `json:"zone"`
//	InfraID          string                   `json:"infraID"`
//	MachineCIDR      string                   `json:"machineCIDR"`
//	VPCID            string                   `json:"vpcID"`
//	Zones            []*CreateInfraOutputZone `json:"zones"`
//	SecurityGroupID  string                   `json:"securityGroupID"`
//	Name             string                   `json:"Name"`
//	BaseDomain       string                   `json:"baseDomain"`
//	BaseDomainPrefix string                   `json:"baseDomainPrefix"`
//	PublicZoneID     string                   `json:"publicZoneID"`
//	PrivateZoneID    string                   `json:"privateZoneID"`
//	LocalZoneID      string                   `json:"localZoneID"`
//	ProxyAddr        string                   `json:"proxyAddr"`
//}
//
//type CreateInfraOutputZone struct {
//	Name     string `json:"name"`
//	SubnetID string `json:"subnetID"`
//}
