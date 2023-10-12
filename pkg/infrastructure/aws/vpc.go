package aws

import (
	"errors"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
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

const duplicatePermissionErrorCode = "InvalidPermission.Duplicate"

func createVPCResources(logger *logrus.Logger, session *session.Session, vpcInput *CreateInfraOptions) error {
	ec2Client := ec2.New(session)

	vpcInput.additionalEC2Tags = make([]*ec2.Tag, 0, len(vpcInput.AdditionalTags))
	for k, v := range vpcInput.AdditionalTags {
		k := k // workaround Go loopvar reuse
		v := v //
		vpcInput.additionalEC2Tags = append(vpcInput.additionalEC2Tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	var err error
	zoneToPrivateSubnetIDMap := make(map[string]string)
	vpcID := vpcInput.vpcID
	if vpcID == "" {
		vpcID, err = vpcInput.createVPC(logger, ec2Client)
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

		bootstrapSG, err := vpcInput.CreateBootstrapSecurityGroup(logger, ec2Client, vpcID)
		if err != nil {
			return err
		}
		// FIXME: only use 0.0.0.0/0 if using public endpoints
		var machineV4Cidrs []string
		if vpcInput.public {
			machineV4Cidrs = []string{"0.0.0.0/0"}
		} else {
			machineV4Cidrs = vpcInput.cidrV4Blocks
		}
		bootstrapIngressPermissions := DefaultBootstrapSGIngressRules(vpcInput.bootstrapSecurityGroupID, machineV4Cidrs)
		if err := vpcInput.AttachSecurityGroupIngressRules(logger, ec2Client, bootstrapSG, bootstrapIngressPermissions); err != nil {
			return err
		}

		// Per zone resources
		_, network, err := net.ParseCIDR(vpcInput.cidrV4Blocks[0])
		if err != nil {
			return err
		}
		privateNetwork, err := cidr.Subnet(network, 1, 0)
		if err != nil {
			return err
		}
		publicNetwork, err := cidr.Subnet(network, 1, 1)
		if err != nil {
			return err
		}

		// If a single-zone deployment, the available CIDR block will be split
		// into two to allow user expansion
		if len(vpcInput.Zones) == 1 {
			privateNetwork, err = cidr.Subnet(privateNetwork, 1, 0)
			if err != nil {
				return err
			}
			publicNetwork, err = cidr.Subnet(publicNetwork, 1, 0)
			if err != nil {
				return err
			}
		}

		var publicSubnetIDs []string
		var privateSubnetIDs []string
		var endpointRouteTableIds []*string
		newBits := int(math.Ceil(math.Log2(float64(len(vpcInput.Zones)))))
		for i, zone := range vpcInput.Zones {
			privateCIDR, err := cidr.Subnet(privateNetwork, newBits, i)
			if err != nil {
				return err
			}
			privateSubnetID, err := vpcInput.CreatePrivateSubnet(logger, ec2Client, vpcID, zone, privateCIDR.String())
			if err != nil {
				return err
			}
			zoneToPrivateSubnetIDMap[zone] = privateSubnetID

			publicCIDR, err := cidr.Subnet(publicNetwork, newBits, i)
			if err != nil {
				return err
			}
			publicSubnetID, err := vpcInput.CreatePublicSubnet(logger, ec2Client, vpcID, zone, publicCIDR.String())
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
	} else {
		logger.WithField("id", vpcID).Debugln("Using user-supplied VPC")

		result, err := ec2Client.DescribeSubnets(&ec2.DescribeSubnetsInput{SubnetIds: aws.StringSlice(vpcInput.privateSubnetIDs)})
		if err != nil {
			return err
		}
		for _, subnet := range result.Subnets {
			zoneToPrivateSubnetIDMap[aws.StringValue(subnet.AvailabilityZone)] = aws.StringValue(subnet.SubnetId)
		}
	}
	vpcInput.zoneToSubnetIDMap = zoneToPrivateSubnetIDMap

	masterSG, err := vpcInput.CreateMasterSecurityGroup(logger, ec2Client, vpcID)
	if err != nil {
		return err
	}
	//vpcInput.masterSecurityGroupID = aws.StringValue(masterSG.GroupId)

	workerSG, err := vpcInput.CreateWorkerSecurityGroup(logger, ec2Client, vpcID)
	if err != nil {
		return err
	}
	//vpcInput.workerSecurityGroupID = aws.StringValue(workerSG.GroupId)

	masterIngressPermissions := DefaultMasterSGIngressRules(vpcInput.masterSecurityGroupID, vpcInput.workerSecurityGroupID, vpcInput.cidrV4Blocks)
	// masterIngressPermissions := DefaultAllowAllSGIngressRules(vpcInput.masterSecurityGroupID, []string{})
	if err := vpcInput.AttachSecurityGroupIngressRules(logger, ec2Client, masterSG, masterIngressPermissions); err != nil {
		return err
	}

	workerIngressPermissions := DefaultWorkerSGIngressRules(vpcInput.workerSecurityGroupID, vpcInput.masterSecurityGroupID, vpcInput.cidrV4Blocks)
	// workerIngressPermissions := DefaultAllowAllSGIngressRules(vpcInput.workerSecurityGroupID, []string{})
	if err := vpcInput.AttachSecurityGroupIngressRules(logger, ec2Client, workerSG, workerIngressPermissions); err != nil {
		return err
	}

	err = vpcInput.CreateLoadBalancers(logger, session, vpcID, vpcInput.privateSubnetIDs, vpcInput.publicSubnetIDs, vpcInput.public)
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateInfraOptions) createVPC(l *logrus.Logger, client ec2iface.EC2API) (string, error) {
	defaultCIDRBlock := o.cidrV4Blocks[0]
	vpcName := fmt.Sprintf("%s-vpc", o.InfraID)
	vpcID, err := o.existingVPC(client, vpcName)
	if err != nil {
		return "", err
	}
	if len(vpcID) == 0 {
		createResult, err := client.CreateVpc(&ec2.CreateVpcInput{
			CidrBlock:         aws.String(defaultCIDRBlock),
			TagSpecifications: ec2TagSpecifications("vpc", vpcName, o.additionalEC2Tags),
		})
		if err != nil {
			return "", fmt.Errorf("failed to create VPC: %w", err)
		}
		vpcID = aws.StringValue(createResult.Vpc.VpcId)
		l.WithField("id", vpcID).Infoln("Created VPC")
	} else {
		l.WithField("id", vpcID).Infoln("Found existing VPC")
	}
	logger := l.WithField("id", vpcID)
	_, err = client.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		VpcId:            aws.String(vpcID),
		EnableDnsSupport: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		return "", fmt.Errorf("failed to modify VPC attributes: %w", err)
	}
	logger.Info("Enabled DNS support on VPC")
	_, err = client.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		VpcId:              aws.String(vpcID),
		EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		return "", fmt.Errorf("failed to modify VPC attributes: %w", err)
	}
	logger.Info("Enabled DNS hostnames on VPC")
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

func ec2TagSpecifications(resourceType, name string, additionalTags []*ec2.Tag) []*ec2.TagSpecification {
	return []*ec2.TagSpecification{
		{
			ResourceType: aws.String(resourceType),
			Tags:         append(ec2Tags(name), additionalTags...),
		},
	}
}

func clusterTag(infraID string) string {
	return fmt.Sprintf("kubernetes.io/cluster/%s", infraID)
}

func ec2Tags(name string) []*ec2.Tag {
	tags := []*ec2.Tag{}
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
			TagSpecifications: ec2TagSpecifications("dhcp-options", "", o.additionalEC2Tags),
		})
		if err != nil {
			return fmt.Errorf("cannot create dhcp-options: %w", err)
		}
		optID = aws.StringValue(result.DhcpOptions.DhcpOptionsId)
		l.WithField("id", optID).Info("Created DHCP options")
	} else {
		l.WithField("id", optID).Info("Found existing DHCP options")
	}
	_, err = client.AssociateDhcpOptions(&ec2.AssociateDhcpOptionsInput{
		DhcpOptionsId: aws.String(optID),
		VpcId:         aws.String(vpcID),
	})
	if err != nil {
		return fmt.Errorf("cannot associate dhcp-options to VPC: %w", err)
	}
	l.WithField("vpc", vpcID).WithField("dhcp option", optID).Infoln("Associated DHCP options with VPC")
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
			TagSpecifications: ec2TagSpecifications("internet-gateway", fmt.Sprintf("%s-igw", o.InfraID), o.additionalEC2Tags),
		})
		if err != nil {
			return "", fmt.Errorf("cannot create internet gateway: %w", err)
		}
		igw = result.InternetGateway
		l.WithField("id", aws.StringValue(igw.InternetGatewayId)).Infoln("Created internet gateway")
	} else {
		l.WithField("id", aws.StringValue(igw.InternetGatewayId)).Infoln("Found existing internet gateway")
	}
	igwId := aws.StringValue(igw.InternetGatewayId)
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
		l.WithField("internet gateway", igwId).WithField("vpc", vpcID).Infoln("Attached internet gateway to VPC")
	}
	return igwId, nil
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

func (o *CreateInfraOptions) createSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string, groupName string) (*ec2.SecurityGroup, error) {
	backoff := wait.Backoff{
		Steps:    10,
		Duration: 3 * time.Second,
		Factor:   1.0,
		Jitter:   0.1,
	}
	securityGroup, err := o.existingSecurityGroup(client, groupName)
	if err != nil {
		return nil, err
	}
	logger := l.WithField("name", groupName)
	if securityGroup == nil {
		result, err := client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
			GroupName:         aws.String(groupName),
			Description:       aws.String("Created by Openshift Installer"),
			VpcId:             aws.String(vpcID),
			TagSpecifications: ec2TagSpecifications("security-group", groupName, o.additionalEC2Tags),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create security group %s: %w", groupName, err)
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
			return nil, fmt.Errorf("cannot find security group that was just created (%s)", aws.StringValue(result.GroupId))
		}
		securityGroup = sgResult.SecurityGroups[0]
		logger.WithField("id", aws.StringValue(securityGroup.GroupId)).Infoln("Created security group")
	} else {
		logger.WithField("id", aws.StringValue(securityGroup.GroupId)).Infoln("Found existing security group")
	}
	return securityGroup, nil
}

func (o *CreateInfraOptions) CreateBootstrapSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (*ec2.SecurityGroup, error) {
	groupName := fmt.Sprintf("%s-bootstrap-sg", o.InfraID)
	securityGroup, err := o.createSecurityGroup(l, client, vpcID, groupName)
	if err != nil {
		return nil, err
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	//sgUserID := aws.StringValue(securityGroup.OwnerId)
	egressPermissions := DefaultSGEgressRules(securityGroupID)
	if err := o.AttachSecurityGroupEgressRules(l, client, securityGroup, egressPermissions); err != nil {
		return nil, err
	}
	o.bootstrapSecurityGroupID = securityGroupID
	return securityGroup, nil
}

func (o *CreateInfraOptions) CreateMasterSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (*ec2.SecurityGroup, error) {
	groupName := fmt.Sprintf("%s-master-sg", o.InfraID)
	securityGroup, err := o.createSecurityGroup(l, client, vpcID, groupName)
	if err != nil {
		return nil, err
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	//sgUserID := aws.StringValue(securityGroup.OwnerId)
	egressPermissions := DefaultSGEgressRules(securityGroupID)
	if err := o.AttachSecurityGroupEgressRules(l, client, securityGroup, egressPermissions); err != nil {
		return nil, err
	}
	o.masterSecurityGroupID = securityGroupID
	return securityGroup, nil
}

func (o *CreateInfraOptions) CreateWorkerSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (*ec2.SecurityGroup, error) {
	groupName := fmt.Sprintf("%s-worker-sg", o.InfraID)
	securityGroup, err := o.createSecurityGroup(l, client, vpcID, groupName)
	if err != nil {
		return nil, err
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	//sgUserID := aws.StringValue(securityGroup.OwnerId)
	egressPermissions := DefaultSGEgressRules(securityGroupID)
	if err := o.AttachSecurityGroupEgressRules(l, client, securityGroup, egressPermissions); err != nil {
		return nil, err
	}

	o.workerSecurityGroupID = securityGroupID
	return securityGroup, nil
}

func (o *CreateInfraOptions) AttachSecurityGroupEgressRules(l *logrus.Logger, client ec2iface.EC2API, securityGroup *ec2.SecurityGroup, egressPermissions []*ec2.IpPermission) error {
	backoff := wait.Backoff{
		Steps:    10,
		Duration: 3 * time.Second,
		Factor:   1.0,
		Jitter:   0.1,
	}
	var egressToAuthorize []*ec2.IpPermission
	for _, permission := range egressPermissions {
		if !includesPermission(securityGroup.IpPermissionsEgress, permission) {
			egressToAuthorize = append(egressToAuthorize, permission)
		}
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	logger := l.WithField("id", securityGroupID)
	logger.Infoln("Authorizing egress rules on security group")
	if len(egressToAuthorize) > 0 {
		err := retry.OnError(backoff, func(error) bool { return true }, func() error {
			res, err := client.AuthorizeSecurityGroupEgress(&ec2.AuthorizeSecurityGroupEgressInput{
				GroupId:       aws.String(securityGroupID),
				IpPermissions: egressToAuthorize,
			})
			if err != nil {
				var awsErr awserr.Error
				// only return error if the permission has not already been set
				if errors.As(err, &awsErr) && awsErr.Code() == duplicatePermissionErrorCode {
					return nil
				}
				return err
			}
			if len(res.SecurityGroupRules) < len(egressToAuthorize) {
				logger.Debugf("authorized %d egress rules out of %d", len(res.SecurityGroupRules), len(egressToAuthorize))
				return fmt.Errorf("not all security group ingress permissions applied")
			}
			return nil
		})
		if err != nil {
			return err
		}
		logger.Infoln("Authorized egress rules on security group")
	}
	return nil
}

func (o *CreateInfraOptions) AttachSecurityGroupIngressRules(l *logrus.Logger, client ec2iface.EC2API, securityGroup *ec2.SecurityGroup, ingressPermissions []*ec2.IpPermission) error {
	backoff := wait.Backoff{
		Steps:    10,
		Duration: 3 * time.Second,
		Factor:   1.0,
		Jitter:   0.1,
	}
	var ingressToAuthorize []*ec2.IpPermission

	for _, permission := range ingressPermissions {
		permission := permission
		if !includesPermission(securityGroup.IpPermissions, permission) {
			ingressToAuthorize = append(ingressToAuthorize, permission)
		}
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	logger := l.WithField("id", securityGroupID)
	logger.Infoln("Authorizing ingress rules on security group")
	if len(ingressToAuthorize) > 0 {
		err := retry.OnError(backoff, func(error) bool { return true }, func() error {
			res, err := client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
				GroupId:       aws.String(securityGroupID),
				IpPermissions: ingressToAuthorize,
			})
			if err != nil {
				var awsErr awserr.Error
				// only return error if the permission has not already been set
				if errors.As(err, &awsErr) && awsErr.Code() == duplicatePermissionErrorCode {
					return nil
				}
				return err
			}
			if len(res.SecurityGroupRules) < len(ingressToAuthorize) {
				logger.Debugf("authorized %d ingress rules out of %d", len(res.SecurityGroupRules), len(ingressToAuthorize))
				return fmt.Errorf("not all security group ingress permissions applied")
			} else {
				return nil
			}
		})
		if err != nil {
			return err
		}
	}
	logger.Infoln("Authorized ingress rules on security group")
	return nil
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

type sgRuleInput struct {
	protocol       string
	cidrBlocks     []string
	ipv6CidrBlocks []string
	fromPort       int64
	toPort         int64
	self           bool
	sourceSGID     string
	description    string
}

func createSGRule(securityGroupID string, input sgRuleInput) *ec2.IpPermission {
	rule := &ec2.IpPermission{
		IpProtocol: aws.String(input.protocol),
	}
	if input.protocol != "-1" {
		rule.FromPort = aws.Int64(input.fromPort)
		rule.ToPort = aws.Int64(input.toPort)
	}
	if len(input.cidrBlocks) > 0 {
		for _, v := range input.cidrBlocks {
			v := v
			rule.IpRanges = append(rule.IpRanges, &ec2.IpRange{CidrIp: aws.String(v)})
		}
	}
	if len(input.ipv6CidrBlocks) > 0 {
		for _, v := range input.ipv6CidrBlocks {
			v := v
			rule.Ipv6Ranges = append(rule.Ipv6Ranges, &ec2.Ipv6Range{CidrIpv6: aws.String(v)})
		}
	}
	if input.self {
		rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, &ec2.UserIdGroupPair{
			GroupId: aws.String(securityGroupID),
		})
	}

	if input.sourceSGID != "" && input.sourceSGID != securityGroupID {
		// [OnwerID/]SecurityGroupID
		if parts := strings.Split(input.sourceSGID, "/"); len(parts) == 1 {
			rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, &ec2.UserIdGroupPair{
				GroupId: aws.String(input.sourceSGID),
			})
		} else {
			rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, &ec2.UserIdGroupPair{
				GroupId: aws.String(parts[0]),
				UserId:  aws.String(parts[1]),
			})
		}
	}

	description := input.description
	if description == "" {
		description = "Created by Openshift Installer"
	}

	for _, v := range rule.IpRanges {
		v.Description = aws.String(description)
	}
	for _, v := range rule.Ipv6Ranges {
		v.Description = aws.String(description)
	}
	for _, v := range rule.PrefixListIds {
		v.Description = aws.String(description)
	}
	for _, v := range rule.UserIdGroupPairs {
		v.Description = aws.String(description)
	}

	return rule
}

func DefaultSGEgressRules(securityGroupID string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		createSGRule(securityGroupID, sgRuleInput{
			protocol:   "-1",
			cidrBlocks: []string{"0.0.0.0/0"},
			fromPort:   0,
			toPort:     0,
		}),
	}
}

func DefaultMasterSGIngressRules(masterSGID string, workerSGID string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// master mcs
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			cidrBlocks: cidrBlocks,
			fromPort:   22623,
			toPort:     22623,
		}),
		// master icmp
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "icmp",
			cidrBlocks: cidrBlocks,
			fromPort:   -1,
			toPort:     -1,
		}),
		// master ssh
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			cidrBlocks: cidrBlocks,
			fromPort:   22,
			toPort:     22,
		}),
		// master https
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			cidrBlocks: cidrBlocks,
			fromPort:   6443,
			toPort:     6443,
		}),
		// master vxlan
		createSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 4789,
			toPort:   4789,
			self:     true,
		}),
		// master geneve
		createSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 6081,
			toPort:   6081,
			self:     true,
		}),
		// master ike
		createSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 500,
			toPort:   500,
			self:     true,
		}),
		// master ike nat_t
		createSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 4500,
			toPort:   4500,
			self:     true,
		}),
		// master esp
		createSGRule(masterSGID, sgRuleInput{
			protocol: "50",
			fromPort: 0,
			toPort:   0,
			self:     true,
		}),
		// master ovndb
		createSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 6641,
			toPort:   6642,
			self:     true,
		}),
		// master internal
		createSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 9000,
			toPort:   9999,
			self:     true,
		}),
		// master internal udp
		createSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 9000,
			toPort:   9999,
			self:     true,
		}),
		// master kube scheduler
		createSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 10259,
			toPort:   10259,
			self:     true,
		}),
		// master kube controller manager
		createSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 10257,
			toPort:   10257,
			self:     true,
		}),
		// master kubelet secure
		createSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 10250,
			toPort:   10250,
			self:     true,
		}),
		// master etcd
		createSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 2379,
			toPort:   2380,
			self:     true,
		}),
		// master services tcp
		createSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 30000,
			toPort:   32767,
			self:     true,
		}),
		// master services udp
		createSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 30000,
			toPort:   32767,
			self:     true,
		}),
		// master vxlan from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   4789,
			toPort:     4789,
			sourceSGID: workerSGID,
		}),
		// master geneve from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   6081,
			toPort:     6081,
			sourceSGID: workerSGID,
		}),
		// master ike from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   500,
			toPort:     500,
			sourceSGID: workerSGID,
		}),
		// master ike nat_t from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   4500,
			toPort:     4500,
			sourceSGID: workerSGID,
		}),
		// master esp from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "50",
			fromPort:   0,
			toPort:     0,
			sourceSGID: workerSGID,
		}),
		// master ovndb from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   6641,
			toPort:     6642,
			sourceSGID: workerSGID,
		}),
		// master internal from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   9000,
			toPort:     9999,
			sourceSGID: workerSGID,
		}),
		// master internal udp from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   9000,
			toPort:     9999,
			sourceSGID: workerSGID,
		}),
		// master kube scheduler from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   10259,
			toPort:     10259,
			sourceSGID: workerSGID,
		}),
		// master kube controler manager from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   10257,
			toPort:     10257,
			sourceSGID: workerSGID,
		}),
		// master kubelet secure from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   10250,
			toPort:     10250,
			sourceSGID: workerSGID,
		}),
		// master services tcp from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   30000,
			toPort:     32767,
			sourceSGID: workerSGID,
		}),
		// master services udp from worker
		createSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   30000,
			toPort:     32767,
			sourceSGID: workerSGID,
		}),
	}
}

// DefaultWorkerSGIngressRules
func DefaultWorkerSGIngressRules(workerSGID string, masterSGID string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// worker icmp
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "icmp",
			cidrBlocks: cidrBlocks,
			fromPort:   -1,
			toPort:     -1,
		}),
		// worker vxlan
		createSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 4789,
			toPort:   4789,
			self:     true,
		}),
		// worker geneve
		createSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 6081,
			toPort:   6081,
			self:     true,
		}),
		// worker ike
		createSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 500,
			toPort:   500,
			self:     true,
		}),
		// worker ike nat_t
		createSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 4500,
			toPort:   4500,
			self:     true,
		}),
		// worker esp
		createSGRule(workerSGID, sgRuleInput{
			protocol: "50",
			fromPort: 0,
			toPort:   0,
			self:     true,
		}),
		// worker internal
		createSGRule(workerSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 9000,
			toPort:   9999,
			self:     true,
		}),
		// worker internal udp
		createSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 9000,
			toPort:   9999,
			self:     true,
		}),
		// worker kubelet insecure
		createSGRule(workerSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 10250,
			toPort:   10250,
			self:     true,
		}),
		// worker services tcp
		createSGRule(workerSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 30000,
			toPort:   32767,
			self:     true,
		}),
		// worker services udp
		createSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 30000,
			toPort:   32767,
			self:     true,
		}),
		// worker vxlan from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   4789,
			toPort:     4789,
			sourceSGID: masterSGID,
		}),
		// worker geneve from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   6081,
			toPort:     6081,
			sourceSGID: masterSGID,
		}),
		// worker ike from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   500,
			toPort:     500,
			sourceSGID: masterSGID,
		}),
		// worker ike nat_t from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   4500,
			toPort:     4500,
			sourceSGID: masterSGID,
		}),
		// worker esp from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "50",
			fromPort:   0,
			toPort:     0,
			sourceSGID: masterSGID,
		}),
		// worker internal from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   9000,
			toPort:     9999,
			sourceSGID: masterSGID,
		}),
		// master internal udp from worker
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   9000,
			toPort:     9999,
			sourceSGID: masterSGID,
		}),
		// worker kubelet insecure from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   10250,
			toPort:     10250,
			sourceSGID: masterSGID,
		}),
		// worker services tcp from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   30000,
			toPort:     32767,
			sourceSGID: masterSGID,
		}),
		// worker services udp from master
		createSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   30000,
			toPort:     32767,
			sourceSGID: masterSGID,
		}),
	}
}

func DefaultAllowAllSGIngressRules(securityGroupID string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		createSGRule(securityGroupID, sgRuleInput{
			protocol:   "-1",
			cidrBlocks: []string{"0.0.0.0/0"},
		}),
	}
}
func DefaultBootstrapSGIngressRules(securityGroupID string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// bootstrap ssh
		createSGRule(securityGroupID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   22,
			toPort:     22,
			cidrBlocks: cidrBlocks,
		}),
		// bootstrap journald gateway
		createSGRule(securityGroupID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   19531,
			toPort:     19531,
			cidrBlocks: cidrBlocks,
		}),
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
		l.WithField("name", name).WithField("id", subnetID).Infoln("Found existing subnet")
		return subnetID, nil
	}
	tagSpec := ec2TagSpecifications("subnet", name, o.additionalEC2Tags)
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
	l.WithField("name", name).WithField("id", subnetID).Infoln("Created subnet")
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
		l.WithField("id", aws.StringValue(natGateway.NatGatewayId)).Infoln("Found existing NAT gateway")
		return *natGateway.NatGatewayId, nil
	}

	eipResult, err := client.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})
	if err != nil {
		return "", fmt.Errorf("cannot allocate EIP for NAT gateway: %w", err)
	}
	allocationID := aws.StringValue(eipResult.AllocationId)
	l.WithField("id", allocationID).Infoln("Created elastic IP for NAT gateway")

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
			Tags:      append(ec2Tags(fmt.Sprintf("%s-eip-%s", o.InfraID, availabilityZone)), o.additionalEC2Tags...),
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
			TagSpecifications: ec2TagSpecifications("natgateway", natGatewayName, o.additionalEC2Tags),
		})
		if err != nil {
			return err
		}
		natGateway = gatewayResult.NatGateway
		l.WithField("id", aws.StringValue(natGateway.NatGatewayId)).Infoln("Created NAT gateway")
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
		l.WithField("route table", aws.StringValue(routeTable.RouteTableId)).WithField("nat gateway", natGatewayID).Infoln("Created route to NAT gateway")
	} else {
		l.WithField("route table", aws.StringValue(routeTable.RouteTableId)).WithField("nat gateway", natGatewayID).Infoln("Found existing route to NAT gateway")
	}
	if !o.hasAssociatedSubnet(routeTable, subnetID) {
		_, err = client.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: routeTable.RouteTableId,
			SubnetId:     aws.String(subnetID),
		})
		if err != nil {
			return "", fmt.Errorf("cannot associate private route table with subnet: %w", err)
		}
		l.WithField("route table", aws.StringValue(routeTable.RouteTableId)).WithField("subnet", subnetID).Infoln("Associated subnet with route table")
	} else {
		l.WithField("route table", aws.StringValue(routeTable.RouteTableId)).WithField("subnet", subnetID).Infoln("Subnet already associated with route table")
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
		l.WithField("route table", tableID).WithField("vpc", vpcID).Infoln("Set main VPC route table")
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
		l.WithField("route table", tableID).WithField("internet gateway", igwID).Infoln("Created route to internet gateway")
	} else {
		l.WithField("route table", tableID).WithField("internet gateway", igwID).Infoln("Found existing route to internet gateway")
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
			l.WithField("route table", tableID).WithField("subnet", subnetID).Infoln("Associated route table with subnet")
		} else {
			l.WithField("route table", tableID).WithField("subnet", subnetID).Infoln("Found existing association between route table and subnet")
		}
	}
	return tableID, nil
}

func (o *CreateInfraOptions) createRouteTable(l *logrus.Logger, client ec2iface.EC2API, vpcID, name string) (*ec2.RouteTable, error) {
	result, err := client.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId:             aws.String(vpcID),
		TagSpecifications: ec2TagSpecifications("route-table", name, o.additionalEC2Tags),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create route table: %w", err)
	}
	l.WithField("name", name).WithField("id", aws.StringValue(result.RouteTable.RouteTableId)).Infoln("Created route table")
	return result.RouteTable, nil
}

func (o *CreateInfraOptions) existingRouteTable(l *logrus.Logger, client ec2iface.EC2API, name string) (*ec2.RouteTable, error) {
	result, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{Filters: o.ec2Filters(name)})
	if err != nil {
		return nil, fmt.Errorf("cannot list route tables: %w", err)
	}
	if len(result.RouteTables) > 0 {
		l.WithField("name", name).WithField("id", aws.StringValue(result.RouteTables[0].RouteTableId)).Infoln("Found existing route table")
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
		l.WithField("name", aws.StringValue(internalLBOutput.LoadBalancers[0].DNSName)).Infoln("Internal Load Balancer created")
	} else {
		internalLBARN = *describeLBOutput.LoadBalancers[0].LoadBalancerArn
		o.LoadBalancers.Internal.ZoneID = *describeLBOutput.LoadBalancers[0].CanonicalHostedZoneId
		o.LoadBalancers.Internal.DNSName = *describeLBOutput.LoadBalancers[0].DNSName
		l.WithField("name", aws.StringValue(describeLBOutput.LoadBalancers[0].DNSName)).Infoln("Internal Load Balancer already exists")
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
		l.WithField("arn", aws.StringValue(internalTGOutput.TargetGroups[0].TargetGroupArn)).Infoln("InternalA Target Group created")
	} else {
		internalATGARN = *describeInternalATGOutput.TargetGroups[0].TargetGroupArn
		l.WithField("arn", aws.StringValue(describeInternalATGOutput.TargetGroups[0].TargetGroupArn)).Infoln("InternalA Target Group already exists")
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
	l.WithField("arn", aws.StringValue(internalAListenerOutput.Listeners[0].ListenerArn)).Infoln("Internal Listener created")

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
			return fmt.Errorf("error creating target group: %w", err)
		}

		internalSTGARN = aws.StringValue(createSTGOutput.TargetGroups[0].TargetGroupArn)
		l.WithField("arn", internalSTGARN).Infoln("InternalS Target Group created")
	} else {
		// Target group already exists
		internalSTGARN = aws.StringValue(describeInternalSTGOutput.TargetGroups[0].TargetGroupArn)
		l.WithField("arn", internalSTGARN).Infoln("InternalS Target Group already exists")
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
		return fmt.Errorf("error creating listener: %w", err)
	}
	l.WithField("arn", aws.StringValue(internalAListenerOutput.Listeners[0].ListenerArn)).Infoln("Internal Service Listener created")

	o.targetGroupARNs = []string{internalATGARN, internalSTGARN}

	if !external {
		l.Debugln("Skipping creation of a public LB because of private cluster")
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
		l.WithField("name", aws.StringValue(externalLBOutput.LoadBalancers[0].DNSName)).Infoln("External Load Balancer created")
	} else {
		externalLBARN = *describeExternalLBOutput.LoadBalancers[0].LoadBalancerArn
		o.LoadBalancers.External.ZoneID = *describeExternalLBOutput.LoadBalancers[0].CanonicalHostedZoneId
		o.LoadBalancers.External.DNSName = *describeExternalLBOutput.LoadBalancers[0].DNSName
		l.WithField("name", aws.StringValue(describeExternalLBOutput.LoadBalancers[0].DNSName)).Infoln("External Load Balancer already exists")
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
		externalTGARN = aws.StringValue(tgOutput.TargetGroups[0].TargetGroupArn)
		l.WithField("arn", externalTGARN).Infoln("External Target Group created")
	} else {
		externalTGARN = aws.StringValue(describeExternalTGOutput.TargetGroups[0].TargetGroupArn)
		l.WithField("arn", externalTGARN).Infoln("External Target Group already exists")
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
	l.WithField("arn", aws.StringValue(externalListenerOutput.Listeners[0].ListenerArn)).Infoln("External Listener created")

	o.targetGroupARNs = append(o.targetGroupARNs, externalTGARN)
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
		l.WithField("id", existingEndpoint).Infoln("Found existing s3 VPC endpoint")
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
			TagSpecifications: ec2TagSpecifications("vpc-endpoint", "", o.additionalEC2Tags),
		})
		if err == nil {
			l.WithField("id", aws.StringValue(result.VpcEndpoint.VpcEndpointId)).Infoln("Created s3 VPC endpoint")
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
	AdditionalTags     map[string]string
	EnableProxy        bool
	SSHKeyFile         string
	additionalEC2Tags  []*ec2.Tag
	cidrV4Blocks       []string
	cidrV6Blocks       []string

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
	targetGroupARNs          []string
	bootstrapSecurityGroupID string
	masterSecurityGroupID    string
	workerSecurityGroupID    string
	zoneToSubnetIDMap        map[string]string
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
