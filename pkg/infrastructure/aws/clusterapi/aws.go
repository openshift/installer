package clusterapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"

	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
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
	region := in.Metadata.ClusterPlatformMetadata.AWS.Region
	session, err := awsconfig.GetSessionWithOptions(
		awsconfig.WithRegion(region),
		awsconfig.WithServiceEndpoints(region, in.Metadata.ClusterPlatformMetadata.AWS.ServiceEndpoints),
	)
	if err != nil {
		return fmt.Errorf("failed to create aws session for bootstrap destroy: %w", err)
	}
	ec2Client := ec2.New(session)

	if err := removeSSHRule(ctx, ec2Client, in.Metadata.AWS.Identifier); err != nil {
		return fmt.Errorf("failed to remove bootstrap SSH rule: %w", err)
	}
	return nil
}

// removeSSHRule removes the SSH rule for accessing the bootstrap node
// by removing the rule from the cluster spec and updating the object.
func removeSSHRule(ctx context.Context, client *ec2.EC2, filters []map[string]string) error {
	logrus.Debug("Removing Bootstrap SSH security rule...")
	tagKey, err := getClusterOwnedKey(filters)
	if err != nil {
		return err
	}

	sg, err := getControlPlaneSecurityGroup(ctx, client, tagKey)
	if err != nil {
		return fmt.Errorf("unable to get controlplane security group: %w", err)
	}

	rule, err := getSSHRule(sg)
	if err != nil {
		return fmt.Errorf("unable to get bootstrap SSH rule: %w", err)
	}

	_, err = client.RevokeSecurityGroupIngressWithContext(ctx, &ec2.RevokeSecurityGroupIngressInput{
		GroupId:       sg.GroupId,
		IpPermissions: rule,
	})
	if err != nil {
		return fmt.Errorf("unable to revoke bootstrap SSH security rule: %w", err)
	}
	logrus.Debug("Removed Bootstrap SSH security rule")
	return nil
}

// getClusterOwnedKey returns a key for one of the cluster-owned filters
// from the cluster metadata.
func getClusterOwnedKey(filters []map[string]string) (string, error) {
	for _, filter := range filters {
		for k, v := range filter {
			// We want either the k8s owned tag or capi owned tag
			// doesn't matter which, we just don't want the cluster id.
			if v == "owned" {
				return k, nil
			}
		}
	}
	return "", errors.New("cluster owned filter was not found in metadata identifiers")
}

func getControlPlaneSecurityGroup(ctx context.Context, client *ec2.EC2, tagKey string) (*ec2.SecurityGroup, error) {
	res, err := client.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("description"),
				Values: []*string{aws.String(string(capa.SecurityGroupControlPlane))},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []*string{&tagKey},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(res.SecurityGroups) == 0 {
		return nil, errors.New("no matching controlplane security group found")
	}

	if len(res.SecurityGroups) > 1 {
		sgs := []string{}
		for _, sg := range res.SecurityGroups {
			logrus.Warnf("Found multiple controlplane security groups %s: %s - %s", *sg.GroupId, *sg.GroupName, *sg.Description)
			sgs = append(sgs, *sg.GroupId)
		}
		return nil, fmt.Errorf("found multiple controlplane security groups: %v", sgs)
	}

	return res.SecurityGroups[0], nil
}

func getSSHRule(sg *ec2.SecurityGroup) ([]*ec2.IpPermission, error) {
	rules := []*ec2.IpPermission{}
	for _, rule := range sg.IpPermissions {
		if rule.ToPort == aws.Int64(22) {
			rules = append(rules, rule)
		}
	}
	if len(rules) == 0 {
		return nil, errors.New("unable to find bootstrap SSH ingress rule")
	}
	return rules, nil
}
