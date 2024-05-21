package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"

	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	awsmanifest "github.com/openshift/installer/pkg/asset/manifests/aws"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

var (
	_ clusterapi.Provider           = (*Provider)(nil)
	_ clusterapi.PreProvider        = (*Provider)(nil)
	_ clusterapi.InfraReadyProvider = (*Provider)(nil)

	errNotFound = errors.New("not found")
)

// Provider implements AWS CAPI installation.
type Provider struct{}

// Name gives the name of the provider, AWS.
func (*Provider) Name() string { return awstypes.Name }

// PreProvision creates the IAM roles used by all nodes in the cluster.
func (*Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	if err := createIAMRoles(ctx, in.InfraID, in.InstallConfig); err != nil {
		return fmt.Errorf("failed to create IAM roles: %w", err)
	}

	amiID, err := copyAMIToRegion(ctx, in.InstallConfig, in.InfraID, string(*in.RhcosImage))
	if err != nil {
		return fmt.Errorf("failed to copy AMI: %w", err)
	}
	for i := range in.MachineManifests {
		if awsMachine, ok := in.MachineManifests[i].(*capa.AWSMachine); ok {
			awsMachine.Spec.AMI = capa.AMIReference{ID: ptr.To(amiID)}
		}
	}
	return nil
}

// InfraReady creates private hosted zone and DNS records.
func (*Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	awsCluster := &capa.AWSCluster{}
	key := k8sClient.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, awsCluster); err != nil {
		return fmt.Errorf("failed to get AWSCluster: %w", err)
	}

	awsSession, err := in.InstallConfig.AWS.Session(ctx)
	if err != nil {
		return fmt.Errorf("failed to get aws session: %w", err)
	}

	subnetIDs := make([]string, 0, len(awsCluster.Spec.NetworkSpec.Subnets))
	for _, s := range awsCluster.Spec.NetworkSpec.Subnets {
		subnetIDs = append(subnetIDs, s.ResourceID)
	}

	vpcID := awsCluster.Spec.NetworkSpec.VPC.ID
	if len(subnetIDs) > 0 && len(vpcID) == 0 {
		// All subnets belong to the same VPC, so we only need one
		vpcID, err = getVPCFromSubnets(ctx, awsSession, awsCluster.Spec.Region, subnetIDs[:1])
		if err != nil {
			return err
		}
	}

	tags := map[string]string{
		fmt.Sprintf("kubernetes.io/cluster/%s", in.InfraID): "owned",
	}
	for k, v := range awsCluster.Spec.AdditionalTags {
		tags[k] = v
	}

	client := awsconfig.NewClient(awsSession)

	phzID := in.InstallConfig.Config.AWS.HostedZone
	if len(phzID) == 0 {
		logrus.Infoln("Creating private Hosted Zone")
		res, err := client.CreateHostedZone(ctx, &awsconfig.HostedZoneInput{
			InfraID:  in.InfraID,
			VpcID:    vpcID,
			Region:   awsCluster.Spec.Region,
			Name:     in.InstallConfig.Config.ClusterDomain(),
			Role:     in.InstallConfig.Config.AWS.HostedZoneRole,
			UserTags: tags,
		})
		if err != nil {
			return fmt.Errorf("failed to create private hosted zone: %w", err)
		}
		phzID = aws.StringValue(res.Id)
	}

	logrus.Infoln("Creating Route53 records for control plane load balancer")
	aliasZoneID, err := getHostedZoneIDForNLB(ctx, awsSession, awsCluster.Spec.Region, awsCluster.Status.Network.APIServerELB.Name)
	if err != nil {
		return fmt.Errorf("failed to find HostedZone ID for NLB: %w", err)
	}
	apiHost := awsCluster.Status.Network.SecondaryAPIServerELB.DNSName
	if awsCluster.Status.Network.APIServerELB.Scheme == capa.ELBSchemeInternetFacing {
		apiHost = awsCluster.Status.Network.APIServerELB.DNSName
	}
	apiIntHost := awsCluster.Spec.ControlPlaneEndpoint.Host
	err = client.CreateOrUpdateRecord(ctx, in.InstallConfig.Config, apiHost, apiIntHost, phzID, aliasZoneID)
	if err != nil {
		return fmt.Errorf("failed to create route53 records: %w", err)
	}

	return nil
}

func getVPCFromSubnets(ctx context.Context, awsSession *session.Session, region string, subnetIDs []string) (string, error) {
	var vpcID string
	var lastError error
	client := ec2.New(awsSession, aws.NewConfig().WithRegion(region))
	err := client.DescribeSubnetsPagesWithContext(
		ctx,
		&ec2.DescribeSubnetsInput{SubnetIds: aws.StringSlice(subnetIDs)},
		func(results *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			for _, subnet := range results.Subnets {
				if subnet.SubnetId == nil {
					continue
				}
				if subnet.SubnetArn == nil {
					lastError = fmt.Errorf("%s has no ARN", *subnet.SubnetId)
					return false
				}
				if subnet.VpcId == nil {
					lastError = fmt.Errorf("%s has no VPC", *subnet.SubnetId)
					return false
				}
				if subnet.AvailabilityZone == nil {
					lastError = fmt.Errorf("%s has no availability zone", *subnet.SubnetId)
					return false
				}
				// All subnets belong to the same VPC
				vpcID = aws.StringValue(subnet.VpcId)
				lastError = nil
				return true
			}
			return !lastPage
		},
	)
	if err == nil {
		err = lastError
	}
	if err != nil {
		return "", fmt.Errorf("failed to get VPC from subnets: %w", err)
	}

	return vpcID, nil
}

// getHostedZoneIDForNLB returns the HostedZone ID for a region from a known table or queries it from the LB instead.
func getHostedZoneIDForNLB(ctx context.Context, awsSession *session.Session, region string, lbName string) (string, error) {
	if hzID, ok := awsconfig.HostedZoneIDPerRegionNLBMap[region]; ok {
		return hzID, nil
	}
	// If the HostedZoneID is not known, query from the LoadBalancer
	input := elbv2.DescribeLoadBalancersInput{
		Names: aws.StringSlice([]string{lbName}),
	}
	res, err := elbv2.New(awsSession).DescribeLoadBalancersWithContext(ctx, &input)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == elbv2.ErrCodeLoadBalancerNotFoundException {
			return "", errNotFound
		}
		return "", fmt.Errorf("failed to list load balancers: %w", err)
	}
	for _, lb := range res.LoadBalancers {
		return *lb.CanonicalHostedZoneId, nil
	}

	return "", errNotFound
}

// DestroyBootstrap removes aws bootstrap resources not handled
// by the deletion of the bootstrap machine by the capi controllers.
func (*Provider) DestroyBootstrap(ctx context.Context, in clusterapi.BootstrapDestroyInput) error {
	region := in.Metadata.AWS.Region
	awsSession, err := awsconfig.GetSessionWithOptions(awsconfig.WithRegion(region), awsconfig.WithServiceEndpoints(region, in.Metadata.AWS.ServiceEndpoints))
	if err != nil {
		return fmt.Errorf("failed to get aws session: %w", err)
	}
	client := ec2.New(awsSession)

	if err := removeSSHRule(ctx, in.Client, client, in.Metadata.InfraID); err != nil {
		return fmt.Errorf("failed to remove bootstrap SSH rule: %w", err)
	}
	return nil
}

// removeSSHRule removes the SSH rule for accessing the bootstrap node
// by removing the rule from the cluster spec and updating the object.
func removeSSHRule(ctx context.Context, cl k8sClient.Client, client *ec2.EC2, infraID string) error {
	awsCluster := &capa.AWSCluster{}
	key := k8sClient.ObjectKey{
		Name:      infraID,
		Namespace: capiutils.Namespace,
	}
	if err := cl.Get(ctx, key, awsCluster); err != nil {
		return fmt.Errorf("failed to get AWSCluster: %w", err)
	}

	sg, ok := awsCluster.Status.Network.SecurityGroups[capa.SecurityGroupControlPlane]
	if !ok {
		// Should never happen but if it does, we have an error msg for it.
		return fmt.Errorf("failed to get controlplane security group")
	}
	input := &ec2.RevokeSecurityGroupIngressInput{GroupId: aws.String(sg.ID)}
	for _, rule := range sg.IngressRules {
		if !strings.EqualFold(rule.Description, awsmanifest.BootstrapSSHDescription) {
			continue
		}
		rule := rule // TODO: remove with golang >= 1.22
		input.IpPermissions = append(input.IpPermissions, ingressRuleFromCAPAToSDKType(&rule))
	}

	if len(input.IpPermissions) == 0 {
		logrus.Infof("Bootstrap SSH rule not found or already revoked")
		return nil
	}

	var lastErr error
	timeout := 5 * time.Minute
	untilTime := time.Now().Add(timeout)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for bootstrap SSH rule to be destroyed...", timeout, untilTime.Format(time.Kitchen), timezone)
	if err := wait.ExponentialBackoffWithContext(ctx, wait.Backoff{
		Duration: time.Second * 10,
		Factor:   float64(1.5),
		Steps:    32,
		Cap:      timeout,
	}, func(ctx context.Context) (bool, error) {
		_, lastErr = client.RevokeSecurityGroupIngressWithContext(ctx, input)
		return lastErr == nil, nil
	}); err != nil {
		if wait.Interrupted(err) {
			return fmt.Errorf("bootstrap ssh rule was not removed within %v: %w", timeout, lastErr)
		}
		return fmt.Errorf("unable to remove bootstrap ssh rule: %w", lastErr)
	}
	logrus.Infof("Bootstrap SSH rule destroyed")

	return nil
}

// ingressRuleFromCAPAToSDKType converts a capa ingress rule to an aws sdk ec2.IpPermission.
// Adapted from https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/pkg/cloud/services/securitygroup/securitygroups.go#L766
func ingressRuleFromCAPAToSDKType(rule *capa.IngressRule) *ec2.IpPermission {
	perm := &ec2.IpPermission{
		IpProtocol: aws.String(string(rule.Protocol)),
		FromPort:   aws.Int64(rule.FromPort),
		ToPort:     aws.Int64(rule.ToPort),
	}

	for _, cidr := range rule.CidrBlocks {
		ipRange := &ec2.IpRange{
			CidrIp:      aws.String(cidr),
			Description: aws.String(rule.Description),
		}
		perm.IpRanges = append(perm.IpRanges, ipRange)
	}

	for _, cidr := range rule.IPv6CidrBlocks {
		ipV6Range := &ec2.Ipv6Range{
			CidrIpv6:    aws.String(cidr),
			Description: aws.String(rule.Description),
		}
		perm.Ipv6Ranges = append(perm.Ipv6Ranges, ipV6Range)
	}

	for _, groupID := range rule.SourceSecurityGroupIDs {
		groupPair := &ec2.UserIdGroupPair{
			GroupId:     aws.String(groupID),
			Description: aws.String(rule.Description),
		}
		perm.UserIdGroupPairs = append(perm.UserIdGroupPairs, groupPair)
	}

	return perm
}
