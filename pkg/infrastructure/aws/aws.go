package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset"
	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/tfvars"
	awstfvars "github.com/openshift/installer/pkg/tfvars/aws"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

const (
	tfVarsFileName         = "terraform.tfvars.json"
	tfPlatformVarsFileName = "terraform.platform.auto.tfvars.json"

	ownedTagKey   = "kubernetes.io/cluster/%s"
	ownedTagValue = "owned"
)

// InfraProvider is the AWS SDK infra provider.
type InfraProvider struct{}

// InitializeProvider initializes the AWS SDK provider.
func InitializeProvider() infrastructure.Provider {
	return InfraProvider{}
}

// Provision creates the infrastructure resources for the stage.
// dir: the path of the install dir
// vars: cluster configuration input variables, such as terraform variables files
// returns a slice of File assets, which will be appended to the cluster asset file list.
func (a InfraProvider) Provision(dir string, vars []*asset.File) ([]*asset.File, error) {
	// Unmarshall input from tf variables, so we can use it along with
	// installConfig and other assets as the contractual input regardless of
	// the implementation.
	clusterConfig := &tfvars.Config{}
	clusterAWSConfig := &awstfvars.Config{}
	for _, file := range vars {
		switch file.Filename {
		case tfVarsFileName:
			if err := json.Unmarshal(file.Data, clusterConfig); err != nil {
				return nil, err
			}
		case tfPlatformVarsFileName:
			if err := json.Unmarshal(file.Data, clusterAWSConfig); err != nil {
				return nil, err
			}
		}
	}

	if clusterConfig == (&tfvars.Config{}) || clusterAWSConfig == (&awstfvars.Config{}) {
		return nil, fmt.Errorf("could not find terraform config files")
	}

	eps := []awstypes.ServiceEndpoint{}
	for k, v := range clusterAWSConfig.CustomEndpoints {
		eps = append(eps, awstypes.ServiceEndpoint{Name: k, URL: v})
	}

	awsSession, err := awssession.GetSessionWithOptions(
		awssession.WithRegion(clusterAWSConfig.Region),
		awssession.WithServiceEndpoints(clusterAWSConfig.Region, eps),
	)
	if err != nil {
		return nil, err
	}

	availabilityZones := sets.New(clusterAWSConfig.MasterSecurityGroups...)
	availabilityZones.Insert(clusterAWSConfig.WorkerAvailabilityZones...)

	tags := mergeTags(clusterAWSConfig.ExtraTags, map[string]string{
		clusterOwnedTag(clusterConfig.ClusterID): ownedTagValue,
	})

	logger := logrus.StandardLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	logger.Infoln("Creating VPC resources")
	ec2Client := ec2.New(awsSession)
	vpcInput := vpcInputOptions{
		infraID:          clusterConfig.ClusterID,
		region:           clusterAWSConfig.Region,
		vpcID:            clusterAWSConfig.VPC,
		cidrV4Block:      clusterConfig.MachineV4CIDRs[0],
		zones:            sets.List(availabilityZones),
		tags:             tags,
		privateSubnetIDs: clusterAWSConfig.PrivateSubnets,
	}
	if clusterAWSConfig.PublicSubnets != nil {
		vpcInput.publicSubnetIDs = *clusterAWSConfig.PublicSubnets
	}

	vpcOutput, err := createVPCResources(ctx, logger, ec2Client, &vpcInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC resources: %w", err)
	}

	logger.Infoln("Creating Load Balancer resources")
	elbClient := elbv2.New(awsSession)
	lbInput := lbInputOptions{
		infraID:          clusterConfig.ClusterID,
		vpcID:            vpcOutput.vpcID,
		privateSubnetIDs: vpcOutput.privateSubnetIDs,
		publicSubnetIDs:  vpcOutput.publicSubnetIDs,
		tags:             tags,
		isPrivateCluster: clusterAWSConfig.PublishStrategy != "External",
	}
	_, err = createLoadBalancers(ctx, logger, elbClient, &lbInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancers: %w", err)
	}

	return nil, fmt.Errorf("provision stage not implemented yet")
}

// DestroyBootstrap destroys the temporary bootstrap resources.
func (a InfraProvider) DestroyBootstrap(dir string) error {
	return nil
}

// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
func (a InfraProvider) ExtractHostAddresses(dir string, ic *types.InstallConfig, ha *infrastructure.HostAddresses) error {
	return nil
}

func mergeTags(lhsTags, rhsTags map[string]string) map[string]string {
	merged := make(map[string]string, len(lhsTags)+len(rhsTags))
	for k, v := range lhsTags {
		merged[k] = v
	}
	for k, v := range rhsTags {
		merged[k] = v
	}
	return merged
}

func clusterOwnedTag(infraID string) string {
	return fmt.Sprintf(ownedTagKey, infraID)
}
