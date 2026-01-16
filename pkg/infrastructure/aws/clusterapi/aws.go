package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/installconfig"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	awsmanifest "github.com/openshift/installer/pkg/asset/manifests/aws"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/dns"
)

var (
	_ clusterapi.Provider           = (*Provider)(nil)
	_ clusterapi.PreProvider        = (*Provider)(nil)
	_ clusterapi.IgnitionProvider   = (*Provider)(nil)
	_ clusterapi.InfraReadyProvider = (*Provider)(nil)
	_ clusterapi.BootstrapDestroyer = (*Provider)(nil)
	_ clusterapi.PostDestroyer      = (*Provider)(nil)

	errNotFound = errors.New("not found")
)

// Provider implements AWS CAPI installation.
type Provider struct {
	bestEffortDeleteIgnition bool
}

// Name gives the name of the provider, AWS.
func (*Provider) Name() string { return awstypes.Name }

// PublicGatherEndpoint indicates that machine ready checks should wait for an ExternalIP
// in the status and use that when gathering bootstrap log bundles.
func (*Provider) PublicGatherEndpoint() clusterapi.GatherEndpoint { return clusterapi.ExternalIP }

// PreProvision creates the IAM roles used by all nodes in the cluster.
func (*Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	if err := createIAMRoles(ctx, in.InfraID, in.InstallConfig); err != nil {
		return fmt.Errorf("failed to create IAM roles: %w", err)
	}

	// The AWSMachine manifests might already have the AMI ID set from the machine pool which takes into account the
	// ways in which the AMI can be specified: the default rhcos if already in the target region, a custom AMI ID set in
	// platform.aws.amiID, and a custom AMI ID specified in the controlPlane stanza. So we just get the value from the
	// first awsmachine manifest we find, instead of duplicating all the inheriting logic here.
	for i := range in.MachineManifests {
		if awsMachine, ok := in.MachineManifests[i].(*capa.AWSMachine); ok {
			// Default/custom AMI already in target region, nothing else to do
			if ptr.Deref(awsMachine.Spec.AMI.ID, "") != "" {
				return nil
			}
		}
	}

	// Notice that we have to use the default RHCOS value because we set the AMI.ID to empty if the default RHCOS is not
	// in the target region and it needs to be copied over. See pkg/asset/machines/clusterapi.go
	amiID, err := copyAMIToRegion(ctx, in.InstallConfig, in.InfraID, in.RhcosImage)
	if err != nil {
		return fmt.Errorf("failed to copy AMI: %w", err)
	}
	// Update manifests with the new ID
	for i := range in.MachineManifests {
		if awsMachine, ok := in.MachineManifests[i].(*capa.AWSMachine); ok {
			awsMachine.Spec.AMI.ID = ptr.To(amiID)
		}
	}
	return nil
}

// Ignition edits the ignition contents to add the public and private load balancer ip addresses to the
// infrastructure CR. The infrastructure CR is updated and added to the ignition files. CAPA creates a
// bucket for ignition, and this ignition data will be placed in the bucket.
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]*corev1.Secret, error) {
	ignOutput, err := editIgnition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to edit bootstrap master or worker ignition: %w", err)
	}

	ignSecrets := []*corev1.Secret{
		clusterapi.IgnitionSecret(ignOutput.UpdatedBootstrapIgn, in.InfraID, "bootstrap"),
		clusterapi.IgnitionSecret(ignOutput.UpdatedMasterIgn, in.InfraID, "master"),
		clusterapi.IgnitionSecret(ignOutput.UpdatedWorkerIgn, in.InfraID, "worker"),
	}
	return ignSecrets, nil
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
		vpcID, err = getVPCFromSubnets(ctx, in.InstallConfig, subnetIDs[:1])
		if err != nil {
			return err
		}
	}

	client := awsconfig.NewClient(awsSession)

	// The user has selected to provision their own DNS solution. Skip the creation of the
	// Hosted Zone(s) and the records for those zones.
	if in.InstallConfig.Config.AWS.UserProvisionedDNS == dns.UserProvisionedDNSEnabled {
		logrus.Debugf("User Provisioned DNS enabled, skipping dns record creation")
		return nil
	}

	logrus.Infoln("Creating Route53 records for control plane load balancer")

	phzID := in.InstallConfig.Config.AWS.HostedZone
	if len(phzID) == 0 {
		logrus.Debugln("Creating private Hosted Zone")

		res, err := client.CreateHostedZone(ctx, &awsconfig.HostedZoneInput{
			InfraID:  in.InfraID,
			VpcID:    vpcID,
			Region:   awsCluster.Spec.Region,
			Name:     in.InstallConfig.Config.ClusterDomain(),
			Role:     in.InstallConfig.Config.AWS.HostedZoneRole,
			UserTags: awsCluster.Spec.AdditionalTags,
		})
		if err != nil {
			return fmt.Errorf("failed to create private hosted zone: %w", err)
		}
		phzID = *res.Id
		logrus.Infoln("Created private Hosted Zone")
	}

	apiName := fmt.Sprintf("api.%s.", in.InstallConfig.Config.ClusterDomain())
	apiIntName := fmt.Sprintf("api-int.%s.", in.InstallConfig.Config.ClusterDomain())

	// Create api record in public zone
	if in.InstallConfig.Config.PublicAPI() {
		zone, err := client.GetBaseDomain(in.InstallConfig.Config.BaseDomain)
		if err != nil {
			return err
		}

		pubLB := awsCluster.Status.Network.SecondaryAPIServerELB
		aliasZoneID, err := getHostedZoneIDForNLB(ctx, in.InstallConfig, pubLB.Name)
		if err != nil {
			return fmt.Errorf("failed to find HostedZone ID for NLB: %w", err)
		}

		if err := client.CreateOrUpdateRecord(ctx, &awsconfig.CreateRecordInput{
			Name:           apiName,
			Region:         awsCluster.Spec.Region,
			DNSTarget:      pubLB.DNSName,
			ZoneID:         *zone.Id,
			AliasZoneID:    aliasZoneID,
			HostedZoneRole: "", // we dont want to assume role here
			EnableAAAA:     in.InstallConfig.Config.AWS.DualStackEnabled(),
		}); err != nil {
			return fmt.Errorf("failed to create records for api in public zone: %w", err)
		}
		logrus.Debugln("Created public API record in public zone")
	}

	aliasZoneID, err := getHostedZoneIDForNLB(ctx, in.InstallConfig, awsCluster.Status.Network.APIServerELB.Name)
	if err != nil {
		return fmt.Errorf("failed to find HostedZone ID for NLB: %w", err)
	}

	// Create api record in private zone
	if err := client.CreateOrUpdateRecord(ctx, &awsconfig.CreateRecordInput{
		Name:           apiName,
		Region:         awsCluster.Spec.Region,
		DNSTarget:      awsCluster.Spec.ControlPlaneEndpoint.Host,
		ZoneID:         phzID,
		AliasZoneID:    aliasZoneID,
		HostedZoneRole: in.InstallConfig.Config.AWS.HostedZoneRole,
		EnableAAAA:     in.InstallConfig.Config.AWS.DualStackEnabled(),
	}); err != nil {
		return fmt.Errorf("failed to create records for api in private zone: %w", err)
	}
	logrus.Debugln("Created public API record in private zone")

	// Create api-int record in private zone
	if err := client.CreateOrUpdateRecord(ctx, &awsconfig.CreateRecordInput{
		Name:           apiIntName,
		Region:         awsCluster.Spec.Region,
		DNSTarget:      awsCluster.Spec.ControlPlaneEndpoint.Host,
		ZoneID:         phzID,
		AliasZoneID:    aliasZoneID,
		HostedZoneRole: in.InstallConfig.Config.AWS.HostedZoneRole,
		EnableAAAA:     in.InstallConfig.Config.AWS.DualStackEnabled(),
	}); err != nil {
		return fmt.Errorf("failed to create records for api-int in private zone: %w", err)
	}
	logrus.Debugln("Created private API record in private zone")

	// FIXME: Provision the dualstack NLB and DNS record sets for default router.
	// Remove after CCM support dualstack NLB.
	if in.InstallConfig.Config.AWS.DualStackEnabled() {
		if err := provisionDefaultRouterResources(ctx, in, awsCluster, phzID); err != nil {
			return fmt.Errorf("failed to create resources for default router in dualstack mode: %w", err)
		}
		logrus.Debugln("Created resources for default router in dualstack mode")
	}

	return nil
}

func getVPCFromSubnets(ctx context.Context, ic *installconfig.InstallConfig, subnetIDs []string) (string, error) {
	var vpcID string

	client, err := ic.AWS.EC2Client(ctx)
	if err != nil {
		return "", err
	}

	paginator := ec2.NewDescribeSubnetsPaginator(client, &ec2.DescribeSubnetsInput{SubnetIds: subnetIDs})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get VPC from subnets: %w", err)
		}
		for _, subnet := range page.Subnets {
			if subnet.SubnetId == nil {
				continue
			}
			if subnet.SubnetArn == nil {
				return "", fmt.Errorf("%s has no ARN", *subnet.SubnetId)
			}
			if subnet.VpcId == nil {
				return "", fmt.Errorf("%s has no VPC", *subnet.SubnetId)
			}
			if subnet.AvailabilityZone == nil {
				return "", fmt.Errorf("%s has no availability zone", *subnet.SubnetId)
			}
			vpcID = *subnet.VpcId
			// All subnets belong to the same VPC
			break
		}
	}

	if vpcID == "" {
		return "", fmt.Errorf("no VPC found for subnets %v", subnetIDs)
	}

	return vpcID, nil
}

// getHostedZoneIDForNLB returns the HostedZone ID for a region from a known table or queries it from the LB instead.
func getHostedZoneIDForNLB(ctx context.Context, ic *installconfig.InstallConfig, lbName string) (string, error) {
	if hzID, ok := awsconfig.HostedZoneIDPerRegionNLBMap[ic.Config.AWS.Region]; ok {
		return hzID, nil
	}

	cfg, err := configv2.LoadDefaultConfig(ctx, configv2.WithRegion(ic.Config.Platform.AWS.Region))
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := elbv2.NewFromConfig(cfg, func(options *elbv2.Options) {
		options.Region = ic.Config.Platform.AWS.Region
		for _, endpoint := range ic.Config.AWS.ServiceEndpoints {
			if strings.EqualFold(endpoint.Name, "elasticloadbalancing") {
				options.BaseEndpoint = aws.String(endpoint.URL)
			}
		}
	})

	// If the HostedZoneID is not known, query from the LoadBalancer
	input := elbv2.DescribeLoadBalancersInput{
		Names: []string{lbName},
	}

	res, err := client.DescribeLoadBalancers(ctx, &input)
	if err != nil {
		var lbError *elbv2types.LoadBalancerNotFoundException
		if errors.As(err, &lbError) {
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
func (p *Provider) DestroyBootstrap(ctx context.Context, in clusterapi.BootstrapDestroyInput) error {
	awsCluster := &capa.AWSCluster{}
	key := k8sClient.ObjectKey{
		Name:      in.Metadata.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, awsCluster); err != nil {
		return fmt.Errorf("failed to get AWSCluster: %w", err)
	}

	// Save this value for use in the post-destroy hook since we don't have capi running anymore by that point.
	p.bestEffortDeleteIgnition = ptr.Deref(awsCluster.Spec.S3Bucket.BestEffortDeleteObjects, false)

	var sgID string
	if sg, ok := awsCluster.Status.Network.SecurityGroups[capa.SecurityGroupControlPlane]; ok && len(sg.ID) > 0 {
		sgID = sg.ID
	} else if ok {
		return fmt.Errorf("control plane security group id is not populated in awscluster status")
	} else {
		keys := make([]capa.SecurityGroupRole, 0, len(awsCluster.Status.Network.SecurityGroups))
		for sgr := range awsCluster.Status.Network.SecurityGroups {
			keys = append(keys, sgr)
		}
		return fmt.Errorf("controlplane not found in cluster security groups: %v", keys)
	}

	region := in.Metadata.ClusterPlatformMetadata.AWS.Region
	session, err := awsconfig.GetSessionWithOptions(
		awsconfig.WithRegion(region),
		awsconfig.WithServiceEndpoints(region, in.Metadata.ClusterPlatformMetadata.AWS.ServiceEndpoints),
	)
	if err != nil {
		return fmt.Errorf("failed to create aws session: %w", err)
	}

	timeout := 15 * time.Minute
	startTime := time.Now()
	untilTime := startTime.Add(timeout)
	timezone, _ := untilTime.Zone()
	logrus.Debugf("Waiting up to %v (until %v %s) for bootstrap SSH rule to be destroyed...", timeout, untilTime.Format(time.Kitchen), timezone)
	if err := wait.PollUntilContextTimeout(ctx, 15*time.Second, timeout, true,
		func(ctx context.Context) (bool, error) {
			if err := removeSSHRule(ctx, in.Client, in.Metadata.InfraID); err != nil {
				// If the cluster object has been modified between Get and Update, k8s client will refuse to update it.
				// In that case, we need to retry.
				if k8serrors.IsConflict(err) {
					logrus.Debugf("AWSCluster update conflict during SSH rule removal: %v", err)
					return false, nil
				}
				return true, fmt.Errorf("failed to remove bootstrap SSH rule: %w", err)
			}
			return isSSHRuleGone(ctx, session, region, sgID)
		},
	); err != nil {
		if wait.Interrupted(err) {
			return fmt.Errorf("bootstrap ssh rule was not removed within %v: %w", timeout, err)
		}
		return fmt.Errorf("unable to remove bootstrap ssh rule: %w", err)
	}
	logrus.Debugf("Completed removing bootstrap SSH rule after %v", time.Since(startTime))

	return nil
}

// removeSSHRule removes the SSH rule for accessing the bootstrap node
// by removing the rule from the cluster spec and updating the object.
func removeSSHRule(ctx context.Context, cl k8sClient.Client, infraID string) error {
	awsCluster := &capa.AWSCluster{}
	key := k8sClient.ObjectKey{
		Name:      infraID,
		Namespace: capiutils.Namespace,
	}
	if err := cl.Get(ctx, key, awsCluster); err != nil {
		return fmt.Errorf("failed to get AWSCluster: %w", err)
	}

	postBootstrapRules := []capa.IngressRule{}
	for _, rule := range awsCluster.Spec.NetworkSpec.AdditionalControlPlaneIngressRules {
		if strings.EqualFold(rule.Description, awsmanifest.BootstrapSSHDescription) {
			continue
		}
		postBootstrapRules = append(postBootstrapRules, rule)
	}

	// The spec has not been updated yet
	if len(postBootstrapRules) < len(awsCluster.Spec.NetworkSpec.AdditionalControlPlaneIngressRules) {
		awsCluster.Spec.NetworkSpec.AdditionalControlPlaneIngressRules = postBootstrapRules

		if err := cl.Update(ctx, awsCluster); err != nil {
			return fmt.Errorf("failed to update AWSCluster during bootstrap destroy: %w", err)
		}
		logrus.Debug("Updated AWSCluster to remove bootstrap SSH rule")
	}

	return nil
}

// isSSHRuleGone checks that the Public SSH rule has been removed from the security group.
func isSSHRuleGone(ctx context.Context, session *session.Session, region, sgID string) (bool, error) {
	sgs, err := awsconfig.DescribeSecurityGroups(ctx, session, []string{sgID}, region)
	if err != nil {
		return false, fmt.Errorf("error getting security group: %w", err)
	}

	if len(sgs) != 1 {
		ids := []string{}
		for _, sg := range sgs {
			ids = append(ids, *sg.GroupId)
		}
		return false, fmt.Errorf("expected exactly one security group with id %s, but got %v", sgID, ids)
	}

	sg := sgs[0]
	for _, rule := range sg.IpPermissions {
		if ptr.Deref(rule.ToPort, 0) != 22 {
			continue
		}
		for _, source := range rule.IpRanges {
			if source.CidrIp != nil && *source.CidrIp == "0.0.0.0/0" {
				ruleDesc := ptr.Deref(source.Description, "[no description]")
				logrus.Debugf("Found ingress rule %s with source cidr %s. Still waiting for deletion...", ruleDesc, *source.CidrIp)
				return false, nil
			}
		}
	}

	return true, nil
}

// PostDestroy deletes the ignition bucket after capi stopped running, so it won't try to reconcile the bucket.
func (p *Provider) PostDestroy(ctx context.Context, in clusterapi.PostDestroyerInput) error {
	bucketName := awsmanifest.GetIgnitionBucketName(in.Metadata.InfraID)
	if err := removeS3Bucket(ctx, in.Metadata.AWS.Region, bucketName, in.Metadata.AWS.ServiceEndpoints); err != nil {
		if p.bestEffortDeleteIgnition {
			logrus.Warnf("failed to delete ignition bucket %s: %v", bucketName, err)
			return nil
		}
		return fmt.Errorf("failed to delete ignition bucket %s: %w", bucketName, err)
	}

	return nil
}

// removeS3Bucket deletes an s3 bucket given its name.
func removeS3Bucket(ctx context.Context, region string, bucketName string, endpoints []awstypes.ServiceEndpoint) error {
	cfg, err := configv2.LoadDefaultConfig(ctx, configv2.WithRegion(region))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.Region = region
		for _, endpoint := range endpoints {
			if strings.EqualFold(endpoint.Name, "s3") {
				options.BaseEndpoint = aws.String(endpoint.URL)
			}
		}
	})

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list objects in bucket %s: %w", bucketName, err)
		}

		var objects []s3types.ObjectIdentifier
		for _, object := range page.Contents {
			objects = append(objects, s3types.ObjectIdentifier{
				Key: object.Key,
			})
		}

		if len(objects) > 0 {
			if _, err = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucketName),
				Delete: &s3types.Delete{
					Objects: objects,
				},
			}); err != nil {
				return fmt.Errorf("failed to delete objects in bucket %s: %w", bucketName, err)
			}
		}
	}

	if _, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucketName)}); err != nil {
		var noSuckBucket *s3types.NoSuchBucket
		if errors.As(err, &noSuckBucket) {
			logrus.Debugf("bucket %q already deleted", bucketName)
			return nil
		}
		return fmt.Errorf("failed to delete bucket %s: %w", bucketName, err)
	}

	logrus.Debugf("bucket %q emptied", bucketName)
	return nil
}

// provisionDefaultRouterResources is a HACK around AWS CCM limitation that does not support dual stack NLB.
// After the infra is ready, we provision the LB and create the necessary record sets for the default router.
// NOTE: When the cluster operators are progressing, we need to wait till worker nodes are ready to register all nodes
// to the target group of this ingress NLB.
func provisionDefaultRouterResources(ctx context.Context, in clusterapi.InfraReadyInput, infraCluster *capa.AWSCluster, privateZoneID string) error {
	ic := in.InstallConfig.Config
	region := ic.Platform.AWS.Region

	// Tags specifications
	clusterTagKey := aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", in.InfraID))
	clusterTagValue := aws.String("owned")

	// Create clients to call AWS API
	// FIXME: Let's ignore the custom endpoints for now
	cfg, err := configv2.LoadDefaultConfig(ctx, configv2.WithRegion(region))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	ec2Client := ec2.NewFromConfig(cfg)
	elbv2Client := elbv2.NewFromConfig(cfg)

	// CCM: Create security groups
	sgInput := &ec2.CreateSecurityGroupInput{
		VpcId:       aws.String(infraCluster.Spec.NetworkSpec.VPC.ID),
		GroupName:   aws.String(fmt.Sprintf("k8s-elb-ingress-%s", in.InfraID)),
		Description: aws.String("Security group for Kubernetes ELB (openshift-ingress/router-default)"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags:         []ec2types.Tag{{Key: clusterTagKey, Value: clusterTagValue}},
			},
		},
	}

	sgOut, err := ec2Client.CreateSecurityGroup(ctx, sgInput)
	if err != nil {
		return fmt.Errorf("failed to create security group for ingress elb: %w", err)
	}
	securityGroupID := sgOut.GroupId
	logrus.Infof("CCM: created security group for ingress LB: %s", *securityGroupID)

	// CCM: Authorize access to port 80 and 443 to anywhere IPv4 and IPv6.
	_, err = ec2Client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: securityGroupID,
		IpPermissions: []ec2types.IpPermission{
			{FromPort: aws.Int32(80), ToPort: aws.Int32(80), IpProtocol: aws.String("tcp"), IpRanges: []ec2types.IpRange{{CidrIp: aws.String(capiutils.AnyIPv4CidrBlock.String())}}},
			{FromPort: aws.Int32(80), ToPort: aws.Int32(80), IpProtocol: aws.String("tcp"), Ipv6Ranges: []ec2types.Ipv6Range{{CidrIpv6: aws.String(capiutils.AnyIPv6CidrBlock.String())}}},
			{FromPort: aws.Int32(443), ToPort: aws.Int32(443), IpProtocol: aws.String("tcp"), IpRanges: []ec2types.IpRange{{CidrIp: aws.String(capiutils.AnyIPv4CidrBlock.String())}}},
			{FromPort: aws.Int32(443), ToPort: aws.Int32(443), IpProtocol: aws.String("tcp"), Ipv6Ranges: []ec2types.Ipv6Range{{CidrIpv6: aws.String(capiutils.AnyIPv6CidrBlock.String())}}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to authorize security group ingress rules for ingress elb: %w", err)
	}
	logrus.Info("CCM: added ingress rules for 80/443 TCP for ingress LB security group")

	// CCM: Create the dual stack NLB for ingress
	// FIXME: Let's ignore the IPv4 IP pool.
	subnetIDs := infraCluster.Spec.NetworkSpec.Subnets.FilterPublic().IDs()
	scheme := elbv2types.LoadBalancerSchemeEnumInternetFacing
	if !in.InstallConfig.Config.PublicIngress() {
		subnetIDs = infraCluster.Spec.NetworkSpec.Subnets.FilterPrivate().IDs()
		scheme = elbv2types.LoadBalancerSchemeEnumInternal
	}

	lbInput := &elbv2.CreateLoadBalancerInput{
		Name:          aws.String(fmt.Sprintf("%.28s-ilb", in.InfraID)),
		Type:          elbv2types.LoadBalancerTypeEnumNetwork,
		Subnets:       subnetIDs,
		Scheme:        scheme,
		IpAddressType: elbv2types.IpAddressTypeDualstack,
		Tags: []elbv2types.Tag{
			{Key: clusterTagKey, Value: clusterTagValue},
		},
		SecurityGroups: []string{aws.ToString(securityGroupID)},
	}

	lbOut, err := elbv2Client.CreateLoadBalancer(ctx, lbInput)
	if err != nil {
		return fmt.Errorf("failed to create ingress dualstack nlb: %w", err)
	}
	lbSpec := lbOut.LoadBalancers[0]
	logrus.Infof("CCM: created dualstack ingress LB: %s. Waiting to be ready...", *lbSpec.LoadBalancerArn)

	// Wait till the load balancer becomes ready (5 minutes)
	waitDuration := 5 * time.Minute
	lbDescInput := &elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []string{aws.ToString(lbSpec.LoadBalancerArn)},
	}
	lbWaiter := elbv2.NewLoadBalancerAvailableWaiter(elbv2Client)
	if err := lbWaiter.Wait(ctx, lbDescInput, waitDuration); err != nil {
		return fmt.Errorf("failed to wait for ingress lb to become ready: %w", err)
	}
	logrus.Info("CCM: dualstack ingress LB is ready")

	// Track target groups by port to create the corresponding
	// listeners on the ingress nlb.
	targetGroupArns := map[int32]string{}

	// CCM: Create target group
	openPorts := []int32{443, 80}
	for _, port := range openPorts {
		name := fmt.Sprintf("%.20s-ingress-%d", in.InstallConfig.Config.GetName(), port)

		ipFamily := elbv2types.TargetGroupIpAddressTypeEnumIpv4
		if ic.AWS.IPFamily == awstypes.DualStackIPv6Primary {
			ipFamily = elbv2types.TargetGroupIpAddressTypeEnumIpv6
		}

		targetGroupInput := &elbv2.CreateTargetGroupInput{
			Name:          aws.String(name),
			Port:          aws.Int32(port),
			Protocol:      elbv2types.ProtocolEnumTcp,
			VpcId:         aws.String(infraCluster.Spec.NetworkSpec.VPC.ID),
			IpAddressType: ipFamily,
			Tags: []elbv2types.Tag{
				{Key: clusterTagKey, Value: clusterTagValue},
			},
		}

		tgOut, err := elbv2Client.CreateTargetGroup(ctx, targetGroupInput)
		if err != nil {
			return fmt.Errorf("failed to create target groups for ingress elb: %w", err)
		}
		logrus.Infof("CCM: created target group on port %d for ingress LB", port)

		// TODO: This hack requires manual registration of worker nodes
		// to the below target groups.
		//
		// FIXME: What happens to worker node scaling or recreating?
		// Do new nodes get registered automatically?
		//
		// HOWTO:
		// We need to get the port exposed on nodes for the ingress default controller.
		//
		// oc -n openshift-ingress get svc router-nodeport-default -o=wide
		//
		//
		// Note down the ports open on nodes, for example:
		// 80:30278/TCP
		// 443:31745/TCP
		// 1936:31350/TCP (health check)
		// Then, use these ports to configure the target port and health check port.

		targetGroupARn := tgOut.TargetGroups[0].TargetGroupArn
		targetGroupArns[port] = aws.ToString(targetGroupARn)

		// Set target group attribute
		attrModifyInput := &elbv2.ModifyTargetGroupAttributesInput{
			TargetGroupArn: targetGroupARn,
			Attributes: []elbv2types.TargetGroupAttribute{
				{
					Key:   aws.String("preserve_client_ip.enabled"),
					Value: aws.String("false"),
				},
				{
					Key:   aws.String("proxy_protocol_v2.enabled"),
					Value: aws.String("true"),
				},
				{
					Key:   aws.String("target_health_state.unhealthy.connection_termination.enabled"),
					Value: aws.String("false"),
				},
				{
					Key:   aws.String("target_health_state.unhealthy.draining_interval_seconds"),
					Value: aws.String("300"),
				},
			},
		}
		if _, err := elbv2Client.ModifyTargetGroupAttributes(ctx, attrModifyInput); err != nil {
			return fmt.Errorf("failed to modify target group attributes for ingress elb: %w", err)
		}
		logrus.Infof("CCM: modified attributes of target group on port %d for ingress LB", port)
	}

	// CCM: Create listeners for ingress nlb.
	for _, port := range openPorts {
		listenerInput := &elbv2.CreateListenerInput{
			DefaultActions: []elbv2types.Action{
				{
					TargetGroupArn: aws.String(targetGroupArns[port]),
					Type:           elbv2types.ActionTypeEnumForward,
				},
			},
			LoadBalancerArn: lbSpec.LoadBalancerArn,
			Port:            aws.Int32(port),
			Protocol:        elbv2types.ProtocolEnumTcp,
			Tags: []elbv2types.Tag{
				{Key: clusterTagKey, Value: clusterTagValue},
			},
		}

		if _, err := elbv2Client.CreateListener(ctx, listenerInput); err != nil {
			return fmt.Errorf("failed to create listeners for ingress elb: %w", err)
		}
		logrus.Infof("CCM: created listener on port %d for ingress LB", port)
	}

	// CIO: Create the wildcards DNS records for the ingress LB
	// in private zone and public zone (if public ingress)
	awsSession, err := in.InstallConfig.AWS.Session(ctx)
	if err != nil {
		return fmt.Errorf("failed to get aws session: %w", err)
	}
	client := awsconfig.NewClient(awsSession)

	appsName := fmt.Sprintf("*.apps.%s.", in.InstallConfig.Config.ClusterDomain())

	// Private zone
	if err := client.CreateOrUpdateRecord(ctx, &awsconfig.CreateRecordInput{
		Name:           appsName,
		Region:         infraCluster.Spec.Region,
		DNSTarget:      aws.ToString(lbSpec.DNSName),
		ZoneID:         privateZoneID,
		AliasZoneID:    aws.ToString(lbSpec.CanonicalHostedZoneId),
		HostedZoneRole: in.InstallConfig.Config.AWS.HostedZoneRole,
		EnableAAAA:     true,
	}); err != nil {
		return fmt.Errorf("failed to create records for ingress in private zone: %w", err)
	}
	logrus.Info("Created ingress record in private zone")

	// CIO: Public zone if cluster is public
	if ic.PublicIngress() {
		zone, err := client.GetBaseDomain(in.InstallConfig.Config.BaseDomain)
		if err != nil {
			return err
		}

		if err := client.CreateOrUpdateRecord(ctx, &awsconfig.CreateRecordInput{
			Name:           appsName,
			Region:         infraCluster.Spec.Region,
			DNSTarget:      aws.ToString(lbSpec.DNSName),
			ZoneID:         aws.ToString(zone.Id),
			AliasZoneID:    aws.ToString(lbSpec.CanonicalHostedZoneId),
			HostedZoneRole: "", // we dont want to assume role here
			EnableAAAA:     true,
		}); err != nil {
			return fmt.Errorf("failed to create records for ingress in public zone: %w", err)
		}
		logrus.Info("Created ingress record in public zone")
	}

	return nil
}
