package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

const errDuplicatePermission = "InvalidPermission.Duplicate"

type sgInputOptions struct {
	infraID          string
	vpcID            string
	cidrV4Blocks     []string
	isPrivateCluster bool
	tags             map[string]string
}

type sgOutput struct {
	bootstrap    string
	controlPlane string
	compute      string
}

func createSecurityGroups(ctx context.Context, logger logrus.FieldLogger, ec2Client ec2iface.EC2API, input *sgInputOptions) (*sgOutput, error) {
	bootstrapSGName := fmt.Sprintf("%s-bootstrap-sg", input.infraID)
	bootstrapSG, err := ensureSecurityGroup(ctx, logger, ec2Client, input.infraID, input.vpcID, bootstrapSGName, input.tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap security group: %w", err)
	}

	cidrs := input.cidrV4Blocks
	if !input.isPrivateCluster {
		cidrs = []string{"0.0.0.0/0"}
	}
	bootstrapIngress := defaultBootstrapSGIngressRules(bootstrapSG.GroupId, cidrs)
	err = authorizeIngressRules(ctx, ec2Client, bootstrapSG, bootstrapIngress)
	if err != nil {
		return nil, fmt.Errorf("failed to attach ingress rules to bootstrap security group: %w", err)
	}

	masterSGName := fmt.Sprintf("%s-master-sg", input.infraID)
	masterSG, err := ensureSecurityGroup(ctx, logger, ec2Client, input.infraID, input.vpcID, masterSGName, input.tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create control plane security group: %w", err)
	}

	workerSGName := fmt.Sprintf("%s-worker-sg", input.infraID)
	workerSG, err := ensureSecurityGroup(ctx, logger, ec2Client, input.infraID, input.vpcID, workerSGName, input.tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute security group: %w", err)
	}

	masterIngress := defaultMasterSGIngressRules(masterSG.GroupId, workerSG.GroupId, input.cidrV4Blocks)
	err = authorizeIngressRules(ctx, ec2Client, masterSG, masterIngress)
	if err != nil {
		return nil, fmt.Errorf("failed to attach ingress rules to master security group: %w", err)
	}
	workerIngress := defaultWorkerSGIngressRules(workerSG.GroupId, masterSG.GroupId, input.cidrV4Blocks)
	err = authorizeIngressRules(ctx, ec2Client, workerSG, workerIngress)
	if err != nil {
		return nil, fmt.Errorf("failed to attach ingress rules to worker security group: %w", err)
	}

	return &sgOutput{
		bootstrap:    aws.StringValue(bootstrapSG.GroupId),
		controlPlane: aws.StringValue(masterSG.GroupId),
		compute:      aws.StringValue(workerSG.GroupId),
	}, nil
}

func ensureSecurityGroup(ctx context.Context, logger logrus.FieldLogger, client ec2iface.EC2API, infraID, vpcID string, name string, tags map[string]string) (*ec2.SecurityGroup, error) {
	filters := ec2Filters(infraID, name)
	l := logger.WithField("name", name)
	createdOrFoundMsg := "Found existing security group"
	sg, err := existingSecurityGroup(ctx, client, filters)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		createdOrFoundMsg = "Created security group"
		gtags := mergeTags(tags, map[string]string{"Name": name})
		sg, err = createSecurityGroup(ctx, client, name, vpcID, gtags)
		if err != nil {
			return nil, err
		}
	}
	l.WithField("id", aws.StringValue(sg.GroupId)).Infoln(createdOrFoundMsg)

	egressPermissions := defaultEgressRules(sg.GroupId)
	// Apply egress rules that haven't been applied yet
	toAuthorize := make([]*ec2.IpPermission, 0, len(egressPermissions))
	for _, permission := range egressPermissions {
		permission := permission
		if !includesPermission(sg.IpPermissionsEgress, permission) {
			toAuthorize = append(toAuthorize, permission)
		}
	}
	if len(toAuthorize) == 0 {
		logger.Infoln("Egress rules already authorized")
		return sg, nil
	}

	err = wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			res, err := client.AuthorizeSecurityGroupEgressWithContext(ctx, &ec2.AuthorizeSecurityGroupEgressInput{
				GroupId:       sg.GroupId,
				IpPermissions: toAuthorize,
			})
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) && awsErr.Code() == errDuplicatePermission {
					return true, nil
				}
				return false, nil
			}
			if len(res.SecurityGroupRules) < len(toAuthorize) {
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to authorize egress rules for security group (%s): %w", name, err)
	}

	return sg, nil
}

func existingSecurityGroup(ctx context.Context, client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.SecurityGroup, error) {
	res, err := client.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		Filters: filters,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list security groups: %w", err)
	}
	for _, sg := range res.SecurityGroups {
		return sg, nil
	}

	return nil, errNotFound
}

func createSecurityGroup(ctx context.Context, client ec2iface.EC2API, name string, vpcID string, tags map[string]string) (*ec2.SecurityGroup, error) {
	res, err := client.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(name),
		Description: aws.String(defaultDescription),
		VpcId:       aws.String(vpcID),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("security-group"),
				Tags:         ec2Tags(tags),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create security group (%s): %w", name, err)
	}

	// Wait for SG to show up
	var out *ec2.DescribeSecurityGroupsOutput
	err = wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			out, err = client.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{GroupIds: []*string{res.GroupId}})
			if err != nil || len(out.SecurityGroups) == 0 {
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find security group (%s) that was just created: %w", name, err)
	}

	return out.SecurityGroups[0], nil
}

func defaultEgressRules(securityGroupID *string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		createSGRule(securityGroupID, "-1", []string{"0.0.0.0/0"}, nil, 0, 0, false, nil),
	}
}

func defaultBootstrapSGIngressRules(securityGroupID *string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// bootstrap ssh
		createSGRule(securityGroupID, "tcp", cidrBlocks, nil, 22, 22, false, nil),
		// bootstrap journald gateway
		createSGRule(securityGroupID, "tcp", cidrBlocks, nil, 19531, 19531, false, nil),
	}
}

func defaultMasterSGIngressRules(masterSGID *string, workerSGID *string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// master mcs
		createSGRule(masterSGID, "tcp", cidrBlocks, nil, 22623, 22623, false, nil),
		// master icmp
		createSGRule(masterSGID, "icmp", cidrBlocks, nil, -1, -1, false, nil),
		// master ssh
		createSGRule(masterSGID, "tcp", cidrBlocks, nil, 22, 22, false, nil),
		// master https
		createSGRule(masterSGID, "tcp", cidrBlocks, nil, 6443, 6443, false, nil),
		// master vxlan
		createSGRule(masterSGID, "udp", nil, nil, 4789, 4789, true, nil),
		// master geneve
		createSGRule(masterSGID, "udp", nil, nil, 6081, 6081, true, nil),
		// master ike
		createSGRule(masterSGID, "udp", nil, nil, 500, 500, true, nil),
		// master ike nat_t
		createSGRule(masterSGID, "udp", nil, nil, 4500, 4500, true, nil),
		// master esp
		createSGRule(masterSGID, "50", nil, nil, 0, 0, true, nil),
		// master ovndb
		createSGRule(masterSGID, "tcp", nil, nil, 6641, 6642, true, nil),
		// master internal
		createSGRule(masterSGID, "tcp", nil, nil, 9000, 9999, true, nil),
		// master internal udp
		createSGRule(masterSGID, "udp", nil, nil, 9000, 9999, true, nil),
		// master kube scheduler
		createSGRule(masterSGID, "tcp", nil, nil, 10259, 10259, true, nil),
		// master kube controller manager
		createSGRule(masterSGID, "tcp", nil, nil, 10257, 10257, true, nil),
		// master kubelet secure
		createSGRule(masterSGID, "tcp", nil, nil, 10250, 10250, true, nil),
		// master etcd
		createSGRule(masterSGID, "tcp", nil, nil, 2379, 2380, true, nil),
		// master services tcp
		createSGRule(masterSGID, "tcp", nil, nil, 30000, 32767, true, nil),
		// master services udp
		createSGRule(masterSGID, "udp", nil, nil, 30000, 32767, true, nil),
		// master vxlan from worker
		createSGRule(masterSGID, "udp", nil, nil, 4789, 4789, false, workerSGID),
		// master geneve from worker
		createSGRule(masterSGID, "udp", nil, nil, 6081, 6081, false, workerSGID),
		// master ike from worker
		createSGRule(masterSGID, "udp", nil, nil, 500, 500, false, workerSGID),
		// master ike nat_t from worker
		createSGRule(masterSGID, "udp", nil, nil, 4500, 4500, false, workerSGID),
		// master esp from worker
		createSGRule(masterSGID, "50", nil, nil, 0, 0, false, workerSGID),
		// master ovndb from worker
		createSGRule(masterSGID, "tcp", nil, nil, 6641, 6642, false, workerSGID),
		// master internal from worker
		createSGRule(masterSGID, "tcp", nil, nil, 9000, 9999, false, workerSGID),
		// master internal udp from worker
		createSGRule(masterSGID, "udp", nil, nil, 9000, 9999, false, workerSGID),
		// master kube scheduler from worker
		createSGRule(masterSGID, "tcp", nil, nil, 10259, 10259, false, workerSGID),
		// master kube controler manager from worker
		createSGRule(masterSGID, "tcp", nil, nil, 10257, 10257, false, workerSGID),
		// master kubelet secure from worker
		createSGRule(masterSGID, "tcp", nil, nil, 10250, 10250, false, workerSGID),
		// master services tcp from worker
		createSGRule(masterSGID, "tcp", nil, nil, 30000, 32767, false, workerSGID),
		// master services udp from worker
		createSGRule(masterSGID, "udp", nil, nil, 30000, 32767, false, workerSGID),
	}
}

func defaultWorkerSGIngressRules(workerSGID *string, masterSGID *string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// worker icmp
		createSGRule(workerSGID, "icmp", cidrBlocks, nil, -1, -1, false, nil),
		// worker vxlan
		createSGRule(workerSGID, "udp", nil, nil, 4789, 4789, true, nil),
		// worker geneve
		createSGRule(workerSGID, "udp", nil, nil, 6081, 6081, true, nil),
		// worker ike
		createSGRule(workerSGID, "udp", nil, nil, 500, 500, true, nil),
		// worker ike nat_t
		createSGRule(workerSGID, "udp", nil, nil, 4500, 4500, true, nil),
		// worker esp
		createSGRule(workerSGID, "50", nil, nil, 0, 0, true, nil),
		// worker internal
		createSGRule(workerSGID, "tcp", nil, nil, 9000, 9999, true, nil),
		// worker internal udp
		createSGRule(workerSGID, "udp", nil, nil, 9000, 9999, true, nil),
		// worker kubelet insecure
		createSGRule(workerSGID, "tcp", nil, nil, 10250, 10250, true, nil),
		// worker services tcp
		createSGRule(workerSGID, "tcp", nil, nil, 30000, 32767, true, nil),
		// worker services udp
		createSGRule(workerSGID, "udp", nil, nil, 30000, 32767, true, nil),
		// worker vxlan from master
		createSGRule(workerSGID, "udp", nil, nil, 4789, 4789, false, masterSGID),
		// worker geneve from master
		createSGRule(workerSGID, "udp", nil, nil, 6081, 6081, false, masterSGID),
		// worker ike from master
		createSGRule(workerSGID, "udp", nil, nil, 500, 500, false, masterSGID),
		// worker ike nat_t from master
		createSGRule(workerSGID, "udp", nil, nil, 4500, 4500, false, masterSGID),
		// worker esp from master
		createSGRule(workerSGID, "50", nil, nil, 0, 0, false, masterSGID),
		// worker internal from master
		createSGRule(workerSGID, "tcp", nil, nil, 9000, 9999, false, masterSGID),
		// master internal udp from worker
		createSGRule(workerSGID, "udp", nil, nil, 9000, 9999, false, masterSGID),
		// worker kubelet insecure from master
		createSGRule(workerSGID, "tcp", nil, nil, 10250, 10250, false, masterSGID),
		// worker services tcp from master
		createSGRule(workerSGID, "tcp", nil, nil, 30000, 32767, false, masterSGID),
		// worker services udp from master
		createSGRule(workerSGID, "udp", nil, nil, 30000, 32767, false, masterSGID),
	}
}

func createSGRule(sgID *string, protocol string, cidrV4Blocks, cidrV6Blocks []string, fromPort, toPort int64, self bool, sourceSGID *string) *ec2.IpPermission {
	rule := &ec2.IpPermission{
		IpProtocol: aws.String(protocol),
	}
	if protocol != "-1" {
		rule.FromPort = aws.Int64(fromPort)
		rule.ToPort = aws.Int64(toPort)
	}

	for _, v := range cidrV4Blocks {
		v := v
		rule.IpRanges = append(rule.IpRanges, &ec2.IpRange{CidrIp: aws.String(v)})
	}
	for _, v := range cidrV6Blocks {
		v := v
		rule.Ipv6Ranges = append(rule.Ipv6Ranges, &ec2.Ipv6Range{CidrIpv6: aws.String(v)})
	}

	if self {
		rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, &ec2.UserIdGroupPair{
			GroupId: sgID,
		})
	}

	if sourceSGID != nil && aws.StringValue(sourceSGID) != aws.StringValue(sgID) {
		// [OnwerID/]SecurityGroupID
		if parts := strings.Split(aws.StringValue(sourceSGID), "/"); len(parts) == 1 {
			rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, &ec2.UserIdGroupPair{
				GroupId: sourceSGID,
			})
		} else {
			rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, &ec2.UserIdGroupPair{
				GroupId: aws.String(parts[0]),
				UserId:  aws.String(parts[1]),
			})
		}
	}

	for _, v := range rule.IpRanges {
		v.Description = aws.String(defaultDescription)
	}
	for _, v := range rule.Ipv6Ranges {
		v.Description = aws.String(defaultDescription)
	}
	for _, v := range rule.PrefixListIds {
		v.Description = aws.String(defaultDescription)
	}
	for _, v := range rule.UserIdGroupPairs {
		v.Description = aws.String(defaultDescription)
	}

	return rule
}

func authorizeIngressRules(ctx context.Context, client ec2iface.EC2API, securityGroup *ec2.SecurityGroup, permissions []*ec2.IpPermission) error {
	toAuthorize := make([]*ec2.IpPermission, 0, len(permissions))
	for _, permission := range permissions {
		permission := permission
		if !includesPermission(securityGroup.IpPermissions, permission) {
			toAuthorize = append(toAuthorize, permission)
		}
	}
	if len(toAuthorize) == 0 {
		return nil
	}

	return wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (done bool, err error) {
			res, err := client.AuthorizeSecurityGroupIngressWithContext(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
				GroupId:       securityGroup.GroupId,
				IpPermissions: toAuthorize,
			})
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) && awsErr.Code() == errDuplicatePermission {
					return true, nil
				}
				return false, nil
			}
			if len(res.SecurityGroupRules) < len(toAuthorize) {
				return false, nil
			}
			return true, nil
		},
	)
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
	return a != nil && b != nil && a.String() == b.String()
}
