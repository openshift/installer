package clusterapi

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ clusterapi.Provider = (*Provider)(nil)

// Provider implements AWS CAPI installation.
type Provider struct{}

// Name gives the name of the provider, AWS.
func (*Provider) Name() string { return awstypes.Name }

func (*Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	awsCluster := &capa.AWSCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(context.Background(), key, awsCluster); err != nil {
		return fmt.Errorf("failed to get AWSCluster: %w", err)
	}
	awsSession, err := in.InstallConfig.AWS.Session(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to get session to create load balancer: %w", err)
	}

	subnetIDs := []string{}
	for _, s := range awsCluster.Spec.NetworkSpec.Subnets {
		if s.IsPublic {
			subnetIDs = append(subnetIDs, s.ResourceID)
		}
	}

	var vpcID string
	var lastError error
	ec2Client := ec2.New(awsSession, aws.NewConfig().WithRegion(awsCluster.Spec.Region))
	err = ec2Client.DescribeSubnetsPagesWithContext(
		context.TODO(),
		&ec2.DescribeSubnetsInput{SubnetIds: []*string{aws.String(subnetIDs[0])}}, //TODO ensure no segfault
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
					lastError = fmt.Errorf("%s has not availability zone", *subnet.SubnetId)
					return false
				}
				vpcID = *subnet.VpcId
			}
			return !lastPage
		},
	)
	if err == nil {
		err = lastError
	}
	if err != nil {
		return fmt.Errorf("error getting VPC ID: %w", err)
	}

	tags := map[string]string{
		fmt.Sprintf("kubernetes.io/cluster/%s", in.InfraID): "owned",
	}
	for k, v := range awsCluster.Spec.AdditionalTags {
		tags[k] = v
	}

	//TODO(padillon): support shared vpc (assume role client)
	r53Client := route53.New(awsSession)
	phz, err := createHostedZone(context.TODO(), r53Client, tags, in.InfraID, in.InstallConfig.Config.ClusterDomain(), vpcID, awsCluster.Spec.Region, true)
	if err != nil {
		return fmt.Errorf("failed to create private hosted zone: %w", err)
	}

	apiHost := awsCluster.Status.Network.SecondaryAPIServerELB.DNSName
	apiIntHost := awsCluster.Spec.ControlPlaneEndpoint.Host
	if err := createDNSRecords(in.InstallConfig, apiHost, apiIntHost, *phz.Id); err != nil {
		return fmt.Errorf("failed to create DNS records: %w", err)
	}
	return nil
}

func (*Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	if err := putIAMRoles(in.InfraID, in.InstallConfig); err != nil {
		return fmt.Errorf("failed to create IAM roles: %w", err)
	}
	return nil
}
