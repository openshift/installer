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
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
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
	_ clusterapi.BootstrapDestroyer = (*Provider)(nil)
	_ clusterapi.PostDestroyer      = (*Provider)(nil)

	errNotFound = errors.New("not found")
)

// Provider implements AWS CAPI installation.
type Provider struct {
	computeSGID              string
	bestEffortDeleteIgnition bool
}

// Name gives the name of the provider, AWS.
func (*Provider) Name() string { return awstypes.Name }

// BootstrapHasPublicIP indicates that machine ready checks
// should wait for an ExternalIP in the status.
func (*Provider) BootstrapHasPublicIP() bool { return true }

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

	client := awsconfig.NewClient(awsSession)

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
		phzID = aws.StringValue(res.Id)
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
		aliasZoneID, err := getHostedZoneIDForNLB(ctx, awsSession, awsCluster.Spec.Region, pubLB.DNSName)
		if err != nil {
			return fmt.Errorf("failed to find HostedZone ID for NLB: %w", err)
		}

		if err := client.CreateOrUpdateRecord(ctx, &awsconfig.CreateRecordInput{
			Name:           apiName,
			Region:         awsCluster.Spec.Region,
			DNSTarget:      pubLB.DNSName,
			ZoneID:         aws.StringValue(zone.Id),
			AliasZoneID:    aliasZoneID,
			HostedZoneRole: "", // we dont want to assume role here
		}); err != nil {
			return fmt.Errorf("failed to create records for api in public zone: %w", err)
		}
		logrus.Debugln("Created public API record in public zone")
	}

	aliasZoneID, err := getHostedZoneIDForNLB(ctx, awsSession, awsCluster.Spec.Region, awsCluster.Status.Network.APIServerELB.Name)
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
	}); err != nil {
		return fmt.Errorf("failed to create records for api-int in private zone: %w", err)
	}
	logrus.Debugln("Created private API record in private zone")

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

	// Save the node security group ID from CAPA for use later.
	if sg, ok := awsCluster.Status.Network.SecurityGroups[capa.SecurityGroupNode]; ok && len(sg.ID) > 0 {
		p.computeSGID = sg.ID
	} else if ok {
		return fmt.Errorf("compute security group id is not populated in awscluster status")
	} else {
		keys := make([]capa.SecurityGroupRole, 0, len(awsCluster.Status.Network.SecurityGroups))
		for sgr := range awsCluster.Status.Network.SecurityGroups {
			keys = append(keys, sgr)
		}
		return fmt.Errorf("compute not found in cluster security groups: %v", keys)
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
	region := in.Metadata.AWS.Region
	session, err := awsconfig.GetSessionWithOptions(
		awsconfig.WithRegion(region),
		awsconfig.WithServiceEndpoints(region, in.Metadata.AWS.ServiceEndpoints),
	)
	if err != nil {
		return fmt.Errorf("failed to create aws session: %w", err)
	}

	bucketName := awsmanifest.GetIgnitionBucketName(in.Metadata.InfraID)
	if err := removeS3Bucket(ctx, session, bucketName); err != nil {
		if p.bestEffortDeleteIgnition {
			logrus.Warnf("failed to delete ignition bucket %s: %v", bucketName, err)
			return nil
		}
		return fmt.Errorf("failed to delete ignition bucket %s: %w", bucketName, err)
	}

	// See https://issues.redhat.com/browse/OCPBUGS-43048
	if err := removeNodeServiceAllRule(ctx, session, p.computeSGID); err != nil {
		return fmt.Errorf("failed to delete compute Node Service superfluous ingress rules: %w", err)
	}

	return nil
}

// removeS3Bucket deletes an s3 bucket given its name.
func removeS3Bucket(ctx context.Context, session *session.Session, bucketName string) error {
	client := s3.New(session)

	iter := s3manager.NewDeleteListIterator(client, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	err := s3manager.NewBatchDeleteWithClient(client).Delete(ctx, iter)
	if err != nil && !isBucketNotFound(err) {
		return err
	}
	logrus.Debugf("bucket %q emptied", bucketName)

	if _, err := client.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucketName)}); err != nil {
		if isBucketNotFound(err) {
			logrus.Debugf("bucket %q already deleted", bucketName)
			return nil
		}
		return err
	}
	return nil
}

func isBucketNotFound(err interface{}) bool {
	switch s3Err := err.(type) {
	case awserr.Error:
		if s3Err.Code() == s3.ErrCodeNoSuchBucket {
			return true
		}
		origErr := s3Err.OrigErr()
		if origErr != nil {
			return isBucketNotFound(origErr)
		}
	case s3manager.Error:
		if s3Err.OrigErr != nil {
			return isBucketNotFound(s3Err.OrigErr)
		}
	case s3manager.Errors:
		if len(s3Err) == 1 {
			return isBucketNotFound(s3Err[0])
		}
	}
	return false
}

// removeNodeServiceAllRule removes a too-permissive capa-upstream ingress rule.
func removeNodeServiceAllRule(ctx context.Context, session *session.Session, sgID string) error {
	input := &ec2.RevokeSecurityGroupIngressInput{
		GroupId: aws.String(sgID),
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(30000),
				ToPort:     aws.Int64(32767),
				IpRanges: []*ec2.IpRange{
					{CidrIp: aws.String("0.0.0.0/0")},
				},
			},
		},
	}
	_, err := ec2.New(session).RevokeSecurityGroupIngressWithContext(ctx, input)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "InvalidPermission.NotFound" {
			logrus.Debugln("Superfluous Node Port Service ingress rules not found or already deleted")
			return nil
		}
		return err
	}
	logrus.Debugln("Removed superfluous Node Port Service ingress rules")

	return nil
}
