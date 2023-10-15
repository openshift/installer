package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

var (
	defaultBackoff = wait.Backoff{
		Steps:    10,
		Duration: 3 * time.Second,
		Factor:   1.0,
		Jitter:   0.1,
	}

	retryBackoff = wait.Backoff{
		Steps:    5,
		Duration: 3 * time.Second,
		Factor:   3.0,
		Jitter:   0.1,
	}
)

var iopsInputPermittedTypes = [...]string{"gp3", "io1", "io2"}

const (
	invalidNATGatewayError       = "InvalidNatGatewayID.NotFound"
	invalidRouteTableID          = "InvalidRouteTableId.NotFound"
	invalidElasticIPNotFound     = "InvalidElasticIpID.NotFound"
	invalidSubnet                = "InvalidSubnet"
	duplicatePermissionErrorCode = "InvalidPermission.Duplicate"
	targetGroupNotFound          = "TargetGroupNotFound"
	sharedCredsLoadError         = "SharedCredsLoad"
	notFoundError                = "NotFound"

	defaultDescription = "Created by Openshift Installer"
)

func iamCreateTags(tags map[string]string) []*iam.Tag {
	iamTags := make([]*iam.Tag, 0, len(tags))
	for k, v := range tags {
		k := k
		v := v
		iamTags = append(iamTags, &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return iamTags
}

func iamCreateRole(client iamiface.IAMAPI, roleName string, assumeRolePolicy string, iamTags []*iam.Tag) (*iam.Role, error) {
	_, err := client.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
		Tags: append(iamTags, &iam.Tag{
			Key:   aws.String("Name"),
			Value: aws.String(roleName),
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create role: %w", err)
	}
	return nil, nil
}

func iamGetRole(client iamiface.IAMAPI, roleName string) (*iam.Role, error) {
	result, err := client.GetRole(&iam.GetRoleInput{RoleName: aws.String(roleName)})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == iam.ErrCodeNoSuchEntityException {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get existing role: %w", err)
	}
	return result.Role, nil
}

func iamGetInstanceProfile(client iamiface.IAMAPI, profileName string) (*iam.InstanceProfile, error) {
	result, err := client.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == iam.ErrCodeNoSuchEntityException {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get existing instance profile: %w", err)
	}
	return result.InstanceProfile, nil
}

func iamCreateInstanceProfile(client iamiface.IAMAPI, profileName string, iamTags []*iam.Tag) (*iam.InstanceProfile, error) {
	result, err := client.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		Path:                aws.String("/"),
		Tags: append(iamTags, &iam.Tag{
			Key:   aws.String("Name"),
			Value: aws.String(profileName),
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create instance profile: %w", err)
	}

	waitContext, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var lastError error
	wait.Until(func() {
		instanceProfile, err := iamGetInstanceProfile(client, profileName)
		if err != nil {
			lastError = err
		}
		if instanceProfile != nil {
			lastError = nil
			cancel()
		}
	}, 2*time.Second, waitContext.Done())
	waitErr := waitContext.Err()
	if waitErr != nil {
		if errors.Is(waitErr, context.DeadlineExceeded) {
			return nil, fmt.Errorf("waiting for profile to exist process timed out: %w", lastError)
		}
	}

	return result.InstanceProfile, nil
}

func iamGetRolePolicy(client iamiface.IAMAPI, roleName, policyName string) (string, error) {
	result, err := client.GetRolePolicy(&iam.GetRolePolicyInput{
		RoleName:   aws.String(roleName),
		PolicyName: aws.String(policyName),
	})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == iam.ErrCodeNoSuchEntityException {
			return "", nil
		}
		return "", fmt.Errorf("cannot get existing role policy: %w", err)
	}
	return aws.StringValue(result.PolicyName), nil
}

func iamAddRoleToProfile(client iamiface.IAMAPI, instanceProfile *iam.InstanceProfile, roleName string) error {
	hasRole := false
	for _, role := range instanceProfile.Roles {
		if aws.StringValue(role.RoleName) == roleName {
			hasRole = true
		}
	}
	if hasRole {
		return nil
	}
	_, err := client.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: instanceProfile.InstanceProfileName,
		RoleName:            aws.String(roleName),
	})
	if err != nil {
		return fmt.Errorf("cannot add role to instance profile: %w", err)
	}
	return nil
}

func iamAddPolicyToRole(client iamiface.IAMAPI, roleName, policyName, policyDocument string) error {
	_, err := client.PutRolePolicy(&iam.PutRolePolicyInput{
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
		RoleName:       aws.String(roleName),
	})
	if err != nil {
		return fmt.Errorf("cannot create profile policy: %w", err)
	}
	return nil
}

type instanceOptions struct {
	name                     string
	amiID                    string
	instanceType             string
	subnetID                 string
	userData                 string
	securityGroupIDs         []string
	associatePublicIPAddress bool
	volumeType               string
	volumeSize               int64
	volumeIOPS               int64
	encrypted                bool
	kmsKeyID                 string
	iamInstanceProfileARN    string
	additionalEC2Tags        []*ec2.Tag
}

func ec2CreateInstance(ec2Client ec2iface.EC2API, options instanceOptions) (*ec2.Instance, error) {
	var err error
	kmsKeyID := options.kmsKeyID
	// Get default KMS key ID
	if kmsKeyID == "" {
		kmsKeyID, err = ec2GetDefaultKMSKeyID(ec2Client)
		if err != nil {
			return nil, err
		}
	}

	// InvalidParameterCombination: The parameter iops is not supported for gp2 volumes.
	var iops *int64
	if options.volumeIOPS > 0 && ec2IsIOPSPermitted(options.volumeType) {
		iops = aws.Int64(options.volumeIOPS)
	}

	// Create a new EC2 instance.
	runResult, err := ec2Client.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(options.amiID),
		InstanceType: aws.String(options.instanceType),
		//SubnetId:     aws.String(subnetID),
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				SubnetId:                 aws.String(options.subnetID),
				Groups:                   aws.StringSlice(options.securityGroupIDs),
				AssociatePublicIpAddress: aws.Bool(options.associatePublicIPAddress),
			},
		},
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(options.userData))),
		// InvalidParameterCombination: Network interfaces and an instance-level security groups may not be specified on the same request
		// SecurityGroupIds:  aws.StringSlice(options.securityGroupIDs),
		MinCount:          aws.Int64(1),
		MaxCount:          aws.Int64(1),
		TagSpecifications: ec2TagSpecifications("instance", options.name, options.additionalEC2Tags),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType: aws.String(options.volumeType),
					VolumeSize: aws.Int64(options.volumeSize),
					Encrypted:  aws.Bool(options.encrypted),
					KmsKeyId:   aws.String(kmsKeyID),
					Iops:       iops,
				},
			},
		},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Arn: aws.String(options.iamInstanceProfileARN),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(runResult.Instances) > 0 {
		return runResult.Instances[0], nil
	}

	return nil, fmt.Errorf("no instances were created")
}

func ec2GetInstance(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (*ec2.Instance, error) {
	existingInstances, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: ec2Filters,
	})
	if err != nil {
		return nil, err
	}
	// If an instance already exists, return its instance ID.
	if len(existingInstances.Reservations) > 0 && len(existingInstances.Reservations[0].Instances) > 0 {
		return existingInstances.Reservations[0].Instances[0], nil
	}
	return nil, nil
}

func ec2GetDefaultKMSKeyID(client ec2iface.EC2API) (string, error) {
	resp, err := client.GetEbsDefaultKmsKeyId(&ec2.GetEbsDefaultKmsKeyIdInput{})
	if err != nil {
		return "", fmt.Errorf("failed to get default KMS key: %w", err)
	}
	return aws.StringValue(resp.KmsKeyId), nil
}

func ec2IsIOPSPermitted(volumeType string) bool {
	for _, permitted := range iopsInputPermittedTypes {
		if volumeType == permitted {
			return true
		}
	}
	return false
}

func ec2CreateVPC(client ec2iface.EC2API, vpcName, cidrBlock string, ec2Tags []*ec2.Tag) (string, error) {
	createResult, err := client.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock:         aws.String(cidrBlock),
		TagSpecifications: ec2TagSpecifications("vpc", vpcName, ec2Tags),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create VPC: %w", err)
	}
	return aws.StringValue(createResult.Vpc.VpcId), nil
}

func ec2GetVPC(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (string, error) {
	result, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{Filters: ec2Filters})
	if err != nil {
		return "", fmt.Errorf("cannot list vpcs: %w", err)
	}
	for _, vpc := range result.Vpcs {
		return aws.StringValue(vpc.VpcId), nil
	}
	return "", nil
}

func ec2VPCEnableDNSSupport(client ec2iface.EC2API, vpcID string) error {
	_, err := client.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		VpcId:            aws.String(vpcID),
		EnableDnsSupport: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		return fmt.Errorf("failed to modify VPC attributes: %w", err)
	}
	return nil
}

func ec2VPCEnableDNSHostnames(client ec2iface.EC2API, vpcID string) error {
	_, err := client.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		VpcId:              aws.String(vpcID),
		EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		return fmt.Errorf("failed to modify VPC attributes: %w", err)
	}
	return nil
}

func ec2TagSpecifications(resourceType, name string, additionalTags []*ec2.Tag) []*ec2.TagSpecification {
	return []*ec2.TagSpecification{
		{
			ResourceType: aws.String(resourceType),
			Tags:         append(ec2Tags(name), additionalTags...),
		},
	}
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

func ec2CreateTags(tags map[string]string) []*ec2.Tag {
	ec2Tags := make([]*ec2.Tag, 0, len(tags))
	for k, v := range tags {
		k := k
		v := v
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return ec2Tags
}

func ec2CreateDHCPOptions(client ec2iface.EC2API, name, domainName string, ec2Tags []*ec2.Tag) (string, error) {
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
		TagSpecifications: ec2TagSpecifications("dhcp-options", name, ec2Tags),
	})
	if err != nil {
		return "", fmt.Errorf("cannot create dhcp-options: %w", err)
	}
	return aws.StringValue(result.DhcpOptions.DhcpOptionsId), nil
}

func ec2AssociateDHCPOptionsToVPC(client ec2iface.EC2API, optID, vpcID string) error {
	_, err := client.AssociateDhcpOptions(&ec2.AssociateDhcpOptionsInput{
		DhcpOptionsId: aws.String(optID),
		VpcId:         aws.String(vpcID),
	})
	if err != nil {
		return fmt.Errorf("cannot associate dhcp-options to VPC: %w", err)
	}
	return nil
}

func ec2GetDHCPOptions(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (string, error) {
	result, err := client.DescribeDhcpOptions(&ec2.DescribeDhcpOptionsInput{Filters: ec2Filters})
	if err != nil {
		return "", fmt.Errorf("cannot list dhcp options: %w", err)
	}
	for _, opt := range result.DhcpOptions {
		return aws.StringValue(opt.DhcpOptionsId), nil
	}
	return "", nil
}

func ec2CreateInternetGateway(client ec2iface.EC2API, gatewayName string, ec2Tags []*ec2.Tag) (*ec2.InternetGateway, error) {
	result, err := client.CreateInternetGateway(&ec2.CreateInternetGatewayInput{
		TagSpecifications: ec2TagSpecifications("internet-gateway", gatewayName, ec2Tags),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create internet gateway: %w", err)
	}
	return result.InternetGateway, nil
}

func ec2AttachInternetGatewayToVPC(client ec2iface.EC2API, igwID, vpcID string) error {
	_, err := client.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(igwID),
		VpcId:             aws.String(vpcID),
	})
	if err != nil {
		return fmt.Errorf("cannot attach internet gateway to vpc: %w", err)
	}
	return nil
}

func ec2GetInternetGateway(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (*ec2.InternetGateway, error) {
	result, err := client.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{Filters: ec2Filters})
	if err != nil {
		return nil, fmt.Errorf("cannot list internet gateways: %w", err)
	}
	for _, igw := range result.InternetGateways {
		return igw, nil
	}
	return nil, nil
}

func ec2CreateSecurityGroup(client ec2iface.EC2API, groupName, vpcID string, ec2Tags []*ec2.Tag) (*ec2.SecurityGroup, error) {
	result, err := client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:         aws.String(groupName),
		Description:       aws.String(defaultDescription),
		VpcId:             aws.String(vpcID),
		TagSpecifications: ec2TagSpecifications("security-group", groupName, ec2Tags),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create security group %s: %w", groupName, err)
	}
	var sgResult *ec2.DescribeSecurityGroupsOutput
	err = retry.OnError(defaultBackoff, func(error) bool { return true }, func() error {
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
	return sgResult.SecurityGroups[0], nil
}

func ec2GetSecurityGroup(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (*ec2.SecurityGroup, error) {
	result, err := client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{Filters: ec2Filters})
	if err != nil {
		return nil, fmt.Errorf("cannot list security groups: %w", err)
	}
	for _, sg := range result.SecurityGroups {
		return sg, nil
	}
	return nil, nil
}

func ec2AuthorizeEgressRules(client ec2iface.EC2API, securityGroupID string, egressPermissions []*ec2.IpPermission) error {
	return retry.OnError(defaultBackoff, func(error) bool { return true }, func() error {
		res, err := client.AuthorizeSecurityGroupEgress(&ec2.AuthorizeSecurityGroupEgressInput{
			GroupId:       aws.String(securityGroupID),
			IpPermissions: egressPermissions,
		})
		if err != nil {
			var awsErr awserr.Error
			// only return error if the permission has not already been set
			if errors.As(err, &awsErr) && awsErr.Code() == duplicatePermissionErrorCode {
				return nil
			}
			return err
		}
		if len(res.SecurityGroupRules) < len(egressPermissions) {
			//logger.Debugf("authorized %d egress rules out of %d", len(res.SecurityGroupRules), len(egressToAuthorize))
			return fmt.Errorf("authorized %d egress rules out of %d", len(res.SecurityGroupRules), len(egressPermissions))
		}
		return nil
	})
}

func ec2AuthorizeIngressRules(client ec2iface.EC2API, securityGroupID string, ingressPermissions []*ec2.IpPermission) error {
	return retry.OnError(defaultBackoff, func(error) bool { return true }, func() error {
		res, err := client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       aws.String(securityGroupID),
			IpPermissions: ingressPermissions,
		})
		if err != nil {
			var awsErr awserr.Error
			// only return error if the permission has not already been set
			if errors.As(err, &awsErr) && awsErr.Code() == duplicatePermissionErrorCode {
				return nil
			}
			return err
		}
		if len(res.SecurityGroupRules) < len(ingressPermissions) {
			//logger.Debugf("authorized %d ingress rules out of %d", len(res.SecurityGroupRules), len(ingressToAuthorize))
			return fmt.Errorf("authorized %d ingress rules out of %d", len(res.SecurityGroupRules), len(ingressPermissions))
		} else {
			return nil
		}
	})
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

func ec2CreateSGRule(securityGroupID string, input sgRuleInput) *ec2.IpPermission {
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
		description = defaultDescription
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

func ec2CreateSubnet(client ec2iface.EC2API, name, zone, vpcID, cidr string, ec2Tags []*ec2.Tag) (string, error) {
	result, err := client.CreateSubnet(&ec2.CreateSubnetInput{
		AvailabilityZone:  aws.String(zone),
		VpcId:             aws.String(vpcID),
		CidrBlock:         aws.String(cidr),
		TagSpecifications: ec2TagSpecifications("subnet", name, ec2Tags),
	})
	if err != nil {
		return "", fmt.Errorf("cannot create public subnet: %w", err)
	}
	var subnetResult *ec2.DescribeSubnetsOutput
	err = retry.OnError(defaultBackoff, func(error) bool { return true }, func() error {
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
	return aws.StringValue(result.Subnet.SubnetId), nil
}

func ec2GetSubnet(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (string, error) {
	result, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{Filters: ec2Filters})
	if err != nil {
		return "", fmt.Errorf("cannot list subnets: %w", err)
	}
	for _, subnet := range result.Subnets {
		return aws.StringValue(subnet.SubnetId), nil
	}
	return "", nil
}

func ec2GetSubnets(client ec2iface.EC2API, subnetIDs []string) ([]*ec2.Subnet, error) {
	result, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(subnetIDs),
	})
	if err != nil {
		return nil, err
	}
	return result.Subnets, nil
}

func ec2CreateRouteTable(client ec2iface.EC2API, name, vpcID string, ec2Tags []*ec2.Tag) (*ec2.RouteTable, error) {
	result, err := client.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId:             aws.String(vpcID),
		TagSpecifications: ec2TagSpecifications("route-table", name, ec2Tags),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create route table: %w", err)
	}
	var rtableResult *ec2.DescribeRouteTablesOutput
	err = retry.OnError(defaultBackoff, func(error) bool { return true }, func() error {
		var err error
		rtableResult, err = client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			RouteTableIds: []*string{result.RouteTable.RouteTableId},
		})
		if err != nil || len(rtableResult.RouteTables) == 0 {
			return fmt.Errorf("not found yet")
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot find route table that was just created (%s)", aws.StringValue(result.RouteTable.RouteTableId))
	}
	return result.RouteTable, nil
}

func ec2GetRouteTable(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (*ec2.RouteTable, error) {
	result, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{Filters: ec2Filters})
	if err != nil {
		return nil, fmt.Errorf("cannot list route tables: %w", err)
	}
	if len(result.RouteTables) > 0 {
		return result.RouteTables[0], nil
	}
	return nil, nil
}

func ec2CreateRoute(client ec2iface.EC2API, tableID, igwID string) error {
	_, err := client.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         aws.String(tableID),
		GatewayId:            aws.String(igwID),
	})
	if err != nil {
		return fmt.Errorf("cannot create route to internet gateway: %w", err)
	}
	return nil
}

func ec2AllocateEIPAddress(client ec2iface.EC2API, ec2Tags []*ec2.Tag) (string, error) {
	eipResult, err := client.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})
	if err != nil {
		return "", fmt.Errorf("cannot allocate EIP for NAT gateway: %w", err)
	}
	allocationID := aws.StringValue(eipResult.AllocationId)

	// NOTE: there's a potential to leak EIP addresses if the following tag operation fails, since we have no way of
	// recognizing the EIP as belonging to the cluster
	isRetriable := func(err error) bool {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			return strings.EqualFold(awsErr.Code(), invalidElasticIPNotFound)
		}
		return false
	}
	err = retry.OnError(retryBackoff, isRetriable, func() error {
		_, err = client.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{aws.String(allocationID)},
			Tags:      ec2Tags,
		})
		return err
	})
	if err != nil {
		return "", fmt.Errorf("cannot tag NAT gateway EIP: %w", err)
	}

	return allocationID, nil
}

func ec2CreateNatGateway(client ec2iface.EC2API, natGatewayName, allocationID, publicSubnetID string, ec2Tags []*ec2.Tag) (string, error) {
	isNATGatewayRetriable := func(err error) bool {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			return strings.EqualFold(awsErr.Code(), invalidSubnet) ||
				strings.EqualFold(awsErr.Code(), invalidElasticIPNotFound)
		}
		return false
	}
	var natGatewayID string
	err := retry.OnError(retryBackoff, isNATGatewayRetriable, func() error {
		gatewayResult, err := client.CreateNatGateway(&ec2.CreateNatGatewayInput{
			AllocationId:      aws.String(allocationID),
			SubnetId:          aws.String(publicSubnetID),
			TagSpecifications: ec2TagSpecifications("natgateway", natGatewayName, ec2Tags),
		})
		if err != nil {
			return err
		}
		natGatewayID = aws.StringValue(gatewayResult.NatGateway.NatGatewayId)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot create NAT gateway: %w", err)
	}
	return natGatewayID, err
}

func ec2GetNatGateway(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (*ec2.NatGateway, error) {
	result, err := client.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{Filter: ec2Filters})
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

func ec2CreateNatGatewayRoute(client ec2iface.EC2API, tableID, natGatewayID string) error {
	isRetriable := func(err error) bool {
		var awsErr awserr.Error
		return errors.As(err, &awsErr) && strings.EqualFold(awsErr.Code(), invalidNATGatewayError)
	}
	err := retry.OnError(retryBackoff, isRetriable, func() error {
		_, err := client.CreateRoute(&ec2.CreateRouteInput{
			RouteTableId:         aws.String(tableID),
			NatGatewayId:         aws.String(natGatewayID),
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
		})
		return err
	})
	if err != nil {
		return fmt.Errorf("cannot create nat gateway route in private route table: %w", err)
	}
	return nil
}

func ec2AssociateRouteTable(client ec2iface.EC2API, tableID, subnetID string) error {
	_, err := client.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(tableID),
		SubnetId:     aws.String(subnetID),
	})
	if err != nil {
		return fmt.Errorf("cannot associate private route table with subnet: %w", err)
	}
	return nil
}

func ec2ReplaceRouteTableAssociation(client ec2iface.EC2API, tableID, associationID string) error {
	_, err := client.ReplaceRouteTableAssociation(&ec2.ReplaceRouteTableAssociationInput{
		RouteTableId:  aws.String(tableID),
		AssociationId: aws.String(associationID),
	})
	if err != nil {
		return fmt.Errorf("cannot set vpc main route table: %w", err)
	}
	return nil
}

func ec2HasInternetGatewayRoute(table *ec2.RouteTable, igwID string) bool {
	for _, route := range table.Routes {
		if aws.StringValue(route.GatewayId) == igwID &&
			aws.StringValue(route.DestinationCidrBlock) == "0.0.0.0/0" {
			return true
		}
	}
	return false
}

func ec2HasAssociatedSubnet(table *ec2.RouteTable, subnetID string) bool {
	for _, assoc := range table.Associations {
		if aws.StringValue(assoc.RouteTableId) == subnetID {
			return true
		}
	}
	return false
}

func ec2HasNATGatewayRoute(table *ec2.RouteTable, natGatewayID string) bool {
	for _, route := range table.Routes {
		if aws.StringValue(route.NatGatewayId) == natGatewayID &&
			aws.StringValue(route.DestinationCidrBlock) == "0.0.0.0/0" {
			return true
		}
	}
	return false
}

func ec2CreateVPCS3Endpoint(client ec2iface.EC2API, endpointName, vpcID, serviceName string, routeTableIDs []*string, ec2Tags []*ec2.Tag) (*ec2.VpcEndpoint, error) {
	isRetriable := func(err error) bool {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			return strings.EqualFold(awsErr.Code(), invalidRouteTableID)
		}
		return false
	}
	var vpcEndpoint *ec2.VpcEndpoint
	if err := retry.OnError(retryBackoff, isRetriable, func() error {
		result, err := client.CreateVpcEndpoint(&ec2.CreateVpcEndpointInput{
			VpcId:             aws.String(vpcID),
			ServiceName:       aws.String(serviceName),
			RouteTableIds:     routeTableIDs,
			TagSpecifications: ec2TagSpecifications("vpc-endpoint", endpointName, ec2Tags),
		})
		if err == nil {
			vpcEndpoint = result.VpcEndpoint
		}
		return err
	}); err != nil {
		return nil, fmt.Errorf("cannot create VPC S3 endpoint: %w", err)
	}
	return vpcEndpoint, nil
}

func ec2GetVPCS3Endpoint(client ec2iface.EC2API, ec2Filters []*ec2.Filter) (string, error) {
	var endpointID string
	result, err := client.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{Filters: ec2Filters})
	if err != nil {
		return "", fmt.Errorf("cannot list vpc endpoints: %w", err)
	}
	for _, endpoint := range result.VpcEndpoints {
		endpointID = aws.StringValue(endpoint.VpcEndpointId)
	}
	return endpointID, nil
}

func ec2CreateFilter(name, value string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(name),
		Values: []*string{aws.String(value)},
	}
}

func r53CreateTags(tags map[string]string) []*route53.Tag {
	r53Tags := make([]*route53.Tag, 0, len(tags))
	for k, v := range tags {
		k := k
		v := v
		r53Tags = append(r53Tags, &route53.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return r53Tags
}

func isErrorRetryable(err error) bool {
	if aggregate, isAggregate := err.(utilerrors.Aggregate); isAggregate {
		if len(aggregate.Errors()) == 1 {
			err = aggregate.Errors()[0]
		} else {
			// We aggregate all errors, utilerrors.Aggregate does for safety reasons not support
			// errors.As (As it can't know what to do when there are multiple matches), so we
			// iterate and bail out if there are only credential load errors
			hasOnlyCredentialLoadErrors := true
			for _, err := range aggregate.Errors() {
				if !isCredentialLoadError(err) {
					hasOnlyCredentialLoadErrors = false
					break
				}
			}
			if hasOnlyCredentialLoadErrors {
				return false
			}

		}
	}

	if isCredentialLoadError(err) {
		return false
	}
	return true
}

func isCredentialLoadError(err error) bool {
	var awsErr awserr.Error
	return errors.As(err, &awsErr) && awsErr.Code() == sharedCredsLoadError
}

func retryRoute53WithBackoff(ctx context.Context, fn func() error) error {
	retriable := func(e error) bool {
		if !isErrorRetryable(e) {
			return false
		}
		select {
		case <-ctx.Done():
			return false
		default:
			return true
		}
	}
	// TODO: inspect the error for throttling details?
	return retry.OnError(defaultBackoff, retriable, fn)
}

func r53CleanZoneID(ID string) string {
	return strings.TrimPrefix(ID, "/hostedzone/")
}

func r53CleanRecordName(name string) string {
	str := name
	s, err := strconv.Unquote(`"` + str + `"`)
	if err != nil {
		return str
	}
	return s
}

func r53CreateRecord(client route53iface.Route53API, zoneID, domain, aliasDNSName, aliasZoneID string) (*route53.ChangeInfo, error) {
	// Create alias records for public endpoints
	createAliasRecordInput := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(r53CleanRecordName(domain)),
						Type: aws.String("A"),
						AliasTarget: &route53.AliasTarget{
							DNSName:              aws.String(aliasDNSName),
							HostedZoneId:         aws.String(aliasZoneID),
							EvaluateTargetHealth: aws.Bool(false),
						},
					},
				},
			},
		},
	}

	result, err := client.ChangeResourceRecordSets(createAliasRecordInput)
	if err != nil {
		return nil, fmt.Errorf("error creating alias record: %w", err)
	}
	return result.ChangeInfo, nil
}

func fqdn(name string) string {
	n := len(name)
	if n == 0 || name[n-1] == '.' {
		return name
	} else {
		return name + "."
	}
}

func r53GetRecord(ctx context.Context, client route53iface.Route53API, zoneID, name, recordType string) (*route53.ResourceRecordSet, error) {
	recordName := fqdn(strings.ToLower(name))
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		StartRecordName: aws.String(recordName),
		StartRecordType: aws.String(recordType),
		MaxItems:        aws.String("1"),
	}

	var record *route53.ResourceRecordSet
	err := client.ListResourceRecordSetsPagesWithContext(ctx, input, func(resp *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
		if len(resp.ResourceRecordSets) == 0 {
			return false
		}

		recordSet := resp.ResourceRecordSets[0]
		responseName := strings.ToLower(r53CleanRecordName(aws.StringValue(recordSet.Name)))
		responseType := strings.ToUpper(aws.StringValue(recordSet.Type))

		if recordName != responseName {
			return false
		}
		if recordType != responseType {
			return false
		}

		record = recordSet
		return false
	})
	return record, err
}

func r53CreateHostedZone(ctx context.Context, client route53iface.Route53API, name, vpcID, region string, isPrivate bool) (string, error) {
	var res *route53.CreateHostedZoneOutput
	if err := retryRoute53WithBackoff(ctx, func() error {
		callRef := fmt.Sprintf("%d", time.Now().Unix())
		if output, err := client.CreateHostedZoneWithContext(ctx, &route53.CreateHostedZoneInput{
			CallerReference: aws.String(callRef),
			Name:            aws.String(name),
			HostedZoneConfig: &route53.HostedZoneConfig{
				PrivateZone: aws.Bool(isPrivate),
				Comment:     aws.String(defaultDescription),
			},
			VPC: &route53.VPC{
				VPCId:     aws.String(vpcID),
				VPCRegion: aws.String(region),
			},
		}); err != nil {
			return err
		} else {
			res = output
			return nil
		}
	}); err != nil {
		return "", fmt.Errorf("failed to create hosted zone: %w", err)
	}
	if res == nil {
		return "", fmt.Errorf("unexpected output from hosted zone creation")
	}

	return r53CleanZoneID(aws.StringValue(res.HostedZone.Id)), nil
}

func r53FindHostedZone(ctx context.Context, client route53iface.Route53API, name string, isPrivate bool) (string, error) {
	var res *route53.HostedZone
	f := func(resp *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
		for idx, zone := range resp.HostedZones {
			if zone.Config != nil && isPrivate == aws.BoolValue(zone.Config.PrivateZone) && strings.TrimSuffix(aws.StringValue(zone.Name), ".") == strings.TrimSuffix(name, ".") {
				res = resp.HostedZones[idx]
				return false
			}
		}
		return !lastPage
	}
	if err := retryRoute53WithBackoff(ctx, func() error {
		if err := client.ListHostedZonesPagesWithContext(ctx, &route53.ListHostedZonesInput{}, f); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", fmt.Errorf("failed to list hosted zones: %w", err)
	}
	if res == nil {
		return "", nil
	}
	return r53CleanZoneID(aws.StringValue(res.Id)), nil
}

func r53TagHostedZone(ctx context.Context, client route53iface.Route53API, id string, r53Tags []*route53.Tag) error {
	if _, err := client.ChangeTagsForResourceWithContext(ctx, &route53.ChangeTagsForResourceInput{
		ResourceType: aws.String("hostedzone"),
		ResourceId:   aws.String(id),
		AddTags:      r53Tags,
	}); err != nil {
		return fmt.Errorf("failed to tag hosted zone: %w", err)
	}
	return nil
}

func r53HostedZoneChangeRecord(ctx context.Context, client route53iface.Route53API, id string, recordSet *route53.ResourceRecordSet) error {
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(id),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action:            aws.String("UPSERT"),
					ResourceRecordSet: recordSet,
				},
			},
		},
	}
	_, err := client.ChangeResourceRecordSetsWithContext(ctx, input)
	return err
}

func r53Tags(name string) []*route53.Tag {
	tags := []*route53.Tag{}
	if len(name) > 0 {
		tags = append(tags, &route53.Tag{
			Key:   aws.String("Name"),
			Value: aws.String(name),
		})
	}
	return tags
}

func elbGetLoadBalancer(client elbv2iface.ELBV2API, name string) (*elbv2.LoadBalancer, error) {
	describeLBInput := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{aws.String(name)},
	}
	describeLBOutput, err := client.DescribeLoadBalancers(describeLBInput)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			if awsErr.Code() != elbv2.ErrCodeLoadBalancerNotFoundException {
				return nil, fmt.Errorf("failed to describe lb: %w", awsErr)
			}
		} else {
			return nil, fmt.Errorf("failed to describe lb: %w", err)
		}
	}
	for _, lb := range describeLBOutput.LoadBalancers {
		return lb, nil
	}
	return nil, nil
}

func elbCreateLoadBalancer(client elbv2iface.ELBV2API, lbName string, subnets []string, internetFacing bool, elbTags []*elbv2.Tag) (*elbv2.LoadBalancer, error) {
	scheme := "internal"
	if internetFacing {
		scheme = "internet-facing"
	}

	lbOutput, err := client.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		CustomerOwnedIpv4Pool: nil,
		IpAddressType:         nil,
		Name:                  aws.String(lbName),
		Scheme:                aws.String(scheme),
		SecurityGroups:        nil,
		SubnetMappings:        nil,
		Subnets:               aws.StringSlice(subnets),
		Tags: append(elbTags, &elbv2.Tag{
			Key:   aws.String("Name"),
			Value: aws.String(lbName),
		}),
		Type: aws.String("network"),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating internal load balancer: %w", err)
	}
	lbARN := *lbOutput.LoadBalancers[0].LoadBalancerArn
	attrInput := &elbv2.ModifyLoadBalancerAttributesInput{
		LoadBalancerArn: aws.String(lbARN),
		Attributes: []*elbv2.LoadBalancerAttribute{
			{
				Key:   aws.String("load_balancing.cross_zone.enabled"),
				Value: aws.String("true"),
			},
		},
	}
	_, err = client.ModifyLoadBalancerAttributes(attrInput)
	if err != nil {
		return nil, fmt.Errorf("error modifying load balancer attributes: %w", err)
	}
	return lbOutput.LoadBalancers[0], nil
}

func elbGetTargetGroup(client elbv2iface.ELBV2API, targetName string) (*elbv2.TargetGroup, error) {
	describeInternalATGInput := &elbv2.DescribeTargetGroupsInput{
		Names: []*string{aws.String(targetName)},
	}
	describeInternalATGOutput, err := client.DescribeTargetGroups(describeInternalATGInput)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			if awsErr.Code() != targetGroupNotFound {
				return nil, fmt.Errorf("failed to describe lb: %w", awsErr)
			}
		} else {
			return nil, fmt.Errorf("failed to describe lb: %w", err)
		}
	}
	for _, tg := range describeInternalATGOutput.TargetGroups {
		return tg, nil
	}
	return nil, nil
}

func elbCreateTargetGroup(client elbv2iface.ELBV2API, targetName string, vpcID string, healthCheckPath string, port int64, elbTags []*elbv2.Tag) (*elbv2.TargetGroup, error) {
	createInternalTGInput := &elbv2.CreateTargetGroupInput{
		HealthCheckEnabled:         aws.Bool(true),
		HealthCheckPath:            aws.String(healthCheckPath),
		HealthCheckPort:            aws.String(strconv.FormatInt(port, 10)),
		HealthCheckProtocol:        aws.String("HTTPS"),
		HealthCheckIntervalSeconds: aws.Int64(10),
		HealthyThresholdCount:      aws.Int64(2),
		UnhealthyThresholdCount:    aws.Int64(2),
		Name:                       aws.String(targetName),
		Port:                       aws.Int64(port),
		Protocol:                   aws.String("TCP"),
		Tags: append(elbTags, &elbv2.Tag{
			Key:   aws.String("Name"),
			Value: aws.String(targetName),
		}),
		TargetType: aws.String("ip"),
		VpcId:      aws.String(vpcID),
	}
	internalTGOutput, err := client.CreateTargetGroup(createInternalTGInput)
	if err != nil {
		return nil, fmt.Errorf("error creating internal target group: %w", err)
	}
	return internalTGOutput.TargetGroups[0], nil
}

func elbCreateListener(client elbv2iface.ELBV2API, listenerName string, lbARN string, tgARN string, port int64, elbTags []*elbv2.Tag) (*elbv2.Listener, error) {
	createInternalSListenerInput := &elbv2.CreateListenerInput{
		LoadBalancerArn: aws.String(lbARN),
		Protocol:        aws.String("TCP"),
		Port:            aws.Int64(port),
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: []*elbv2.TargetGroupTuple{
						{
							TargetGroupArn: aws.String(tgARN),
						},
					},
				},
			},
		},
		Tags: append(elbTags, &elbv2.Tag{
			Key:   aws.String("Name"),
			Value: aws.String(listenerName),
		}),
	}

	listenerOutput, err := client.CreateListener(createInternalSListenerInput)
	if err != nil {
		return nil, fmt.Errorf("error creating listener: %w", err)
	}
	return listenerOutput.Listeners[0], nil
}

func elbRegisterTargetGroup(elbClient elbv2iface.ELBV2API, targetGroupARN string, ipAddressID string) error {
	_, err := elbClient.RegisterTargets(&elbv2.RegisterTargetsInput{
		TargetGroupArn: aws.String(targetGroupARN),
		Targets: []*elbv2.TargetDescription{
			{
				Id: aws.String(ipAddressID),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error registering target group: %w", err)
	}
	return nil
}

func elbCreateTags(tags map[string]string) []*elbv2.Tag {
	elbTags := make([]*elbv2.Tag, 0, len(tags))
	for k, v := range tags {
		k := k
		v := v
		elbTags = append(elbTags, &elbv2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return elbTags
}

func s3CreateTags(tags map[string]string) []*s3.Tag {
	s3Tags := make([]*s3.Tag, 0, len(tags))
	for k, v := range tags {
		k := k
		v := v
		s3Tags = append(s3Tags, &s3.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return s3Tags
}

func s3BucketExists(s3Client s3iface.S3API, name string) (bool, error) {
	_, err := s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			if awsErr.Code() == s3.ErrCodeNoSuchBucket || awsErr.Code() == notFoundError {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func s3CreateBucket(s3Client s3iface.S3API, name string) error {
	if _, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	}); err != nil {
		return fmt.Errorf("error creating S3 bucket: %w", err)
	}
	return nil
}

func s3TagBucket(s3Client s3iface.S3API, bucketName string, tags []*s3.Tag) error {
	_, err := s3Client.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3.Tagging{
			TagSet: append(tags, &s3.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(bucketName),
			}),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to tag s3 bucket: %w", err)
	}
	return nil
}

func s3BucketPutObject(s3Client s3iface.S3API, bucketName, key, content string) error {
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	})

	if err != nil {
		return fmt.Errorf("error uploading object to S3: %w", err)
	}
	return nil
}

func s3BucketTagObject(s3Client s3iface.S3API, bucketName string, key string, tags []*s3.Tag) error {
	_, err := s3Client.PutObjectTagging(&s3.PutObjectTaggingInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("bootstrap.ign"),
		Tagging: &s3.Tagging{
			TagSet: append(tags, &s3.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(bucketName),
			}),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to tag ignition s3 bucket object: %w", err)
	}
	return nil
}
