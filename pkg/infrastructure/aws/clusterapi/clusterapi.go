package clusterapi

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
)

type InfraHelper struct {
	clusterapi.CAPIInfraHelper
}

func (a InfraHelper) PreProvision(in clusterapi.PreProvisionInput) error {
	// TODO(padillon): skip if users bring their own roles
	if err := putIAMRoles(in.ClusterID, in.InstallConfig); err != nil {
		return fmt.Errorf("failed to create IAM roles: %w", err)
	}
	return nil
}

func (a InfraHelper) ControlPlaneAvailable(in clusterapi.ControlPlaneAvailableInput) error {
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

	subnets := awsCluster.Status.Network.APIServerELB.SubnetIDs
	var vpcID string
	var lastError error
	ec2Client := ec2.New(awsSession, aws.NewConfig().WithRegion(awsCluster.Spec.Region))
	err = ec2Client.DescribeSubnetsPagesWithContext(
		context.TODO(),
		&ec2.DescribeSubnetsInput{SubnetIds: []*string{aws.String(subnets[0])}}, //TODO ensure no segfault
		func(results *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			for _, subnet := range results.Subnets {
				if subnet.SubnetId == nil {
					continue
				}
				if subnet.SubnetArn == nil {
					lastError = errors.Errorf("%s has no ARN", *subnet.SubnetId)
					return false
				}
				if subnet.VpcId == nil {
					lastError = errors.Errorf("%s has no VPC", *subnet.SubnetId)
					return false
				}
				if subnet.AvailabilityZone == nil {
					lastError = errors.Errorf("%s has not availability zone", *subnet.SubnetId)
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

	elbClient := elbv2.New(awsSession)

	tags := map[string]string{
		fmt.Sprintf("kubernetes.io/cluster/%s", in.InfraID): "owned",
	}
	for k, v := range awsCluster.Spec.AdditionalTags {
		tags[k] = v
	}

	lb, err := createExtLB(elbClient, subnets, tags, in.InfraID, vpcID)
	if err != nil {
		return fmt.Errorf("error creating external LB: %w", err)
	}

	//TODO(padillon): support shared vpc (assume role client)
	r53Client := route53.New(awsSession)
	phz, err := createHostedZone(context.TODO(), r53Client, tags, in.InfraID, in.InstallConfig.Config.ClusterDomain(), vpcID, awsCluster.Spec.Region, true)

	if err := createDNSRecords(in.InstallConfig, *lb.DNSName, in.Cluster.Spec.ControlPlaneEndpoint.Host, *phz.Id); err != nil {
		return fmt.Errorf("failed to create DNS records: %w", err)
	}

	return nil
}
