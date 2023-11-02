package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset"
	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/tfvars"
	awstfvars "github.com/openshift/installer/pkg/tfvars/aws"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/version"
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
	lbOutput, err := createLoadBalancers(ctx, logger, elbClient, &lbInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancers: %w", err)
	}

	logger.Infoln("Creating DNS resources")
	r53Client := route53.New(awsSession)
	dnsInput := dnsInputOptions{
		infraID:           clusterConfig.ClusterID,
		region:            clusterAWSConfig.Region,
		baseDomain:        clusterConfig.BaseDomain,
		clusterDomain:     clusterConfig.ClusterDomain,
		vpcID:             vpcOutput.vpcID,
		tags:              tags,
		lbExternalZoneID:  lbOutput.external.zoneID,
		lbExternalZoneDNS: lbOutput.external.dnsName,
		lbInternalZoneID:  lbOutput.internal.zoneID,
		lbInternalZoneDNS: lbOutput.internal.dnsName,
		isPrivateCluster:  clusterAWSConfig.PublishStrategy != "External",
	}
	err = createDNSResources(ctx, logger, r53Client, &dnsInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS rsources: %w", err)
	}

	logger.Infoln("Creating security groups")
	sgInput := sgInputOptions{
		infraID:          clusterConfig.ClusterID,
		vpcID:            vpcOutput.vpcID,
		cidrV4Blocks:     clusterConfig.MachineV4CIDRs,
		isPrivateCluster: clusterAWSConfig.PublishStrategy != "External",
		tags:             tags,
	}
	sgOutput, err := createSecurityGroups(ctx, logger, ec2Client, &sgInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create security groups: %w", err)
	}

	logger.Infoln("Creating bootstrap resources")
	bootstrapSubnet := vpcOutput.privateSubnetIDs[0]
	if clusterAWSConfig.PublishStrategy == "External" {
		bootstrapSubnet = vpcOutput.publicSubnetIDs[0]
	}
	bootstrapInput := bootstrapInputOptions{
		instanceInputOptions: instanceInputOptions{
			infraID:           clusterConfig.ClusterID,
			amiID:             clusterAWSConfig.AMI,
			instanceType:      clusterAWSConfig.MasterInstanceType,
			iamRole:           clusterAWSConfig.MasterIAMRoleName,
			volumeType:        "gp2",
			volumeSize:        30,
			volumeIOPS:        0,
			isEncrypted:       true,
			metadataAuth:      clusterAWSConfig.BootstrapMetadataAuthentication,
			kmsKeyID:          clusterAWSConfig.KMSKeyID,
			securityGroupIds:  []string{sgOutput.bootstrap, sgOutput.controlPlane},
			targetGroupARNs:   lbOutput.targetGroupArns,
			subnetID:          bootstrapSubnet,
			associatePublicIP: clusterAWSConfig.PublishStrategy == "External",
			userData:          clusterAWSConfig.BootstrapIgnitionStub,
			tags:              tags,
		},
		ignitionBucket:  clusterAWSConfig.IgnitionBucket,
		ignitionContent: clusterConfig.IgnitionBootstrap,
	}
	iamClient := iam.New(awsSession)
	s3Client := s3.New(awsSession)
	err = createBootstrapResources(ctx, logger, ec2Client, iamClient, s3Client, elbClient, &bootstrapInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap resources: %w", err)
	}

	logger.Infoln("Creating control plane resources")
	controlPlaneInput := controlPlaneInputOptions{
		instanceInputOptions: instanceInputOptions{
			infraID:           clusterConfig.ClusterID,
			amiID:             clusterAWSConfig.AMI,
			instanceType:      clusterAWSConfig.MasterInstanceType,
			iamRole:           clusterAWSConfig.MasterIAMRoleName,
			volumeType:        clusterAWSConfig.Type,
			volumeSize:        clusterAWSConfig.Size,
			volumeIOPS:        clusterAWSConfig.IOPS,
			isEncrypted:       clusterAWSConfig.Encrypted,
			kmsKeyID:          clusterAWSConfig.KMSKeyID,
			metadataAuth:      clusterAWSConfig.MasterMetadataAuthentication,
			securityGroupIds:  append(clusterAWSConfig.MasterSecurityGroups, sgOutput.controlPlane),
			targetGroupARNs:   lbOutput.targetGroupArns,
			associatePublicIP: false,
			userData:          clusterConfig.IgnitionMaster,
			tags:              tags,
		},
		nReplicas:         clusterConfig.Masters,
		privateSubnetIDs:  vpcOutput.privateSubnetIDs,
		zoneToSubnetMap:   vpcOutput.zoneToSubnetMap,
		availabilityZones: clusterAWSConfig.MasterAvailabilityZones,
	}
	err = createControlPlaneResources(ctx, logger, ec2Client, iamClient, elbClient, &controlPlaneInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create control plane resources: %w", err)
	}

	logger.Infoln("Creating compute resources")
	computeInput := computeInputOptions{
		infraID: clusterConfig.ClusterID,
		tags:    tags,
	}
	err = createComputeResources(ctx, logger, iamClient, &computeInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute resources: %w", err)
	}

	return nil, nil
}

// DestroyBootstrap destroys the temporary bootstrap resources.
func (a InfraProvider) DestroyBootstrap(dir string) error {
	// Unmarshall input from tf variables, so we can use it along with
	// installConfig and other assets as the contractual input regardless of
	// the implementation.
	clusterConfig := &tfvars.Config{}
	data, err := os.ReadFile(filepath.Join(dir, tfVarsFileName))
	if err == nil {
		err = json.Unmarshal(data, clusterConfig)
	}
	if err != nil {
		return fmt.Errorf("failed to load cluster terraform variables: %w", err)
	}
	clusterAWSConfig := &awstfvars.Config{}
	data, err = os.ReadFile(filepath.Join(dir, tfPlatformVarsFileName))
	if err == nil {
		err = json.Unmarshal(data, clusterAWSConfig)
	}
	if err != nil {
		return fmt.Errorf("failed to load AWS terraform variables: %w", err)
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
		return err
	}
	awsSession.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshiftInstaller.OpenshiftInstallerUserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Destroyer", version.Raw),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	logger := logrus.StandardLogger()
	input := &destroyInputOptions{
		infraID: clusterConfig.ClusterID,
		region:  clusterAWSConfig.Region,
	}
	err = destroyBootstrapResources(ctx, logger, awsSession, input)
	if err != nil {
		return fmt.Errorf("failed to delete bootstrap resources: %w", err)
	}

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
